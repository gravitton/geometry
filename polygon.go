package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"

	xslices "github.com/gravitton/x/slices"
)

// Polygon is a 2D polygon given by its vertices. The vertex count is not checked: a polygon with
// fewer than three vertices is degenerate but every method still answers for it.
//
// Points is shared, not copied: Pol keeps the slice it is given and every method that
// returns a Polygon allocates a new one. The polygon is immutable as long as its caller does
// not write into that slice. A nil Points stays nil through every method, so IsZero holds
// after Translate, Scale, Int or Float.
type Polygon[T Number] struct {
	Points []Point[T]
}

// Pol is shorthand for Polygon{vertices}.
func Pol[T Number](vertices []Point[T]) Polygon[T] {
	return Polygon[T]{vertices}
}

// minMax returns the minimum and maximum corner of the vertices, the corners of Bounds, exact
// for an integer T where Bounds places a center: the pair the intersection tests reject shapes
// by before examining any edge, without placing a rectangle. An empty polygon has no corners
// and returns two zero points.
func (p Polygon[T]) minMax() (Point[T], Point[T]) {
	if p.Empty() {
		return Point[T]{}, Point[T]{}
	}

	a, b := p.Points[0], p.Points[0]
	for _, v := range p.Points[1:] {
		a = Point[T]{min(a.X, v.X), min(a.Y, v.Y)}
		b = Point[T]{max(b.X, v.X), max(b.Y, v.Y)}
	}

	return a, b
}

// Edges iterates the polygon edges in vertex order, each from a vertex to the next and the
// last one closing back to the first, without allocating; collect them with slices.Collect
// where a slice is needed. A single vertex yields one zero-length edge and an empty polygon
// yields nothing. Every walk over the outline, containment, distance and the crossings of a
// segment, reads these edges, so the boundary they join is the one every test agrees on.
func (p Polygon[T]) Edges() iter.Seq[Line[T]] {
	return edgesOf(p.Points)
}

// Vertices iterates Points in order, the form every shape with an outline offers, so a
// polygon is drawn or measured by the same loop as a Rectangle or a RegularPolygon. Index
// Points directly where a position is needed.
func (p Polygon[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// edgesOf iterates the edges joining the vertices in order, the last one closing back to the
// first: the one iterator Polygon.Edges, Rectangle.Edges and Line.crossings walk, so the
// edges a segment crosses are the edges the walk reads. The slice is read as the edges are
// yielded and never retained, so a caller may pass a slice of a local array.
func edgesOf[T Number](vertices []Point[T]) iter.Seq[Line[T]] {
	return func(yield func(Line[T]) bool) {
		n := len(vertices)
		for i, vertex := range vertices {
			if !yield(Line[T]{vertex, vertices[(i+1)%n]}) {
				return
			}
		}
	}
}

// Center returns the polygon centroid: the center of the enclosed area, so a vertex added in
// the middle of an edge does not move it. A polygon that encloses no area, with fewer than
// three vertices or all of them collinear, has no such center and falls back to the average of
// its vertices; an empty polygon returns the zero point.
// For integer T the centroid is rounded like every other result stored into T; use float64
// when centroid accuracy matters.
//
// The sum is taken with the origin moved to the first vertex, so the products stay small and
// every edge at that vertex contributes exactly zero: a polygon of one or two vertices has an
// exact zero area whatever its coordinates.
func (p Polygon[T]) Center() Point[T] {
	if p.Empty() {
		return Point[T]{}
	}

	origin := p.Points[0].Float()
	offset := origin.Vector().Negate()

	var x, y, twiceArea float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		wedge := shifted.wedge()

		x += (shifted.Start.X + shifted.End.X) * wedge
		y += (shifted.Start.Y + shifted.End.Y) * wedge
		twiceArea += wedge
	}

	if twiceArea == 0 {
		return p.mean()
	}

	centroid := origin.AddXY(x/(3*twiceArea), y/(3*twiceArea))

	return Point[T]{Cast[T](centroid.X), Cast[T](centroid.Y)}
}

// Area returns the area enclosed by the polygon, by the shoelace formula, regardless of winding.
// It is a float64 even for an integer T, since a lattice polygon can enclose half a unit;
// a self-intersecting polygon has its lobes cancel where they wind the opposite way.
func (p Polygon[T]) Area() float64 {
	var twiceArea float64
	for edge := range p.Edges() {
		twiceArea += edge.wedge()
	}

	return math.Abs(twiceArea) / 2
}

// Perimeter returns the total length of the edges.
func (p Polygon[T]) Perimeter() float64 {
	var perimeter float64
	for edge := range p.Edges() {
		perimeter += edge.Length()
	}

	return perimeter
}

// Bounds returns the axis-aligned bounding rectangle of the vertices, or the zero rectangle
// for a polygon without vertices.
func (p Polygon[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(p.minMax())
}

// mean returns the average of the vertices, which Center falls back to when the
// polygon encloses no area.
func (p Polygon[T]) mean() Point[T] {
	var x, y float64
	for _, v := range p.Points {
		x, y = x+float64(v.X), y+float64(v.Y)
	}

	n := float64(len(p.Points))

	return Point[T]{Cast[T](x / n), Cast[T](y / n)}
}

// Translate creates a new Polygon translated by the given vector (applied to all vertices).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return point.Add(vector)
	})}
}

// MoveTo creates a new Polygon whose centroid is moved to point, preserving shape.
// For integer T the centroid is rounded, and a translation by a whole number of units preserves
// the fractional part of the centroid, so the moved centroid lands on point except when it sits
// exactly on a half and rounding away from zero flips side as the sign changes.
func (p Polygon[T]) MoveTo(point Point[T]) Polygon[T] {
	return p.Translate(point.Subtract(p.Center()))
}

// Scale creates a new Polygon uniformly scaled about its centroid by the factor.
func (p Polygon[T]) Scale(factor float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Unscale creates a new Polygon uniformly scaled about its centroid by the inverse factor, the
// inverse of Scale. Like Divide it panics for a zero factor.
func (p Polygon[T]) Unscale(factor float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Divide(factor))
	})}
}

// UnscaleXY creates a new Polygon scaled about its centroid by the inverse of the given
// factors, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (p Polygon[T]) UnscaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).DivideXY(factorX, factorY))
	})}
}

// Transform creates a new Polygon by applying the given matrix to every vertex, like Point.Transform.
func (p Polygon[T]) Transform[M Float](matrix Matrix[M]) Polygon[T] {
	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return point.Transform(matrix)
	})}
}

// Rotate creates a new Polygon rotated by the given angle (in radians) about its centroid, in
// the same sense as Vector.Rotate. For integer T the centroid and every rotated vertex are
// rounded; only multiples of 90° keep the shape exactly.
func (p Polygon[T]) Rotate(angle float64) Polygon[T] {
	pivot := p.Center()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return point.RotateAround(pivot, angle)
	})}
}

// Contains reports whether the given point lies within the polygon, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains. The interior follows the
// even-odd rule, so a self-intersecting polygon excludes the regions it winds around twice.
// A point outside the extent of the vertices is rejected before any edge is examined.
func (p Polygon[T]) Contains(point Point[T]) bool {
	if p.Empty() {
		return false
	}

	a, b := p.minMax()

	return p.encloses(point, a, b)
}

// DistanceTo returns the distance from the given point to the nearest point of the polygon:
// zero for a point within it, the same closed convention as Contains, and otherwise the
// distance to the nearest edge. An empty polygon is infinitely far from every point.
func (p Polygon[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(p.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. It is a float64 even for an integer T, like Line.DistanceSquaredTo, since the
// nearest point of an edge is not a lattice point in general.
func (p Polygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	return p.walk(point)
}

// walk returns the squared distance from the point to the polygon in one pass over the edges:
// zero for a point on an edge within Epsilon of T, snapped the way Line.DistanceSquaredTo
// snaps it with the tolerance computed once for the walk, or inside by the even-odd rule, the
// squared distance to the nearest edge otherwise, and infinity for an empty polygon. Contains
// and DistanceSquaredTo are both built on it, so the two agree by construction.
func (p Polygon[T]) walk(point Point[T]) float64 {
	inside, distance := false, math.Inf(1)

	for edge := range p.Edges() {
		distance = min(distance, edge.distanceSquaredTo(point))
		if lessOrEqualSquared[T](distance, 0) {
			return 0
		}
		if edge.crossesRay(point) {
			inside = !inside
		}
	}

	if inside {
		return 0
	}

	return distance
}

// encloses is Contains for a caller that already holds the extent of the vertices, so the
// intersection tests walk the vertices once for the box and reuse it for every point they test.
func (p Polygon[T]) encloses(point, a, b Point[T]) bool {
	return point.Between(a, b) && p.walk(point) == 0
}

// Intersects reports whether the polygons share a point: a vertex of one lies within the other,
// or an edge of one crosses an edge of the other. Touching polygons intersect, within Epsilon
// of T, the same closed convention as Contains, and an empty polygon intersects nothing.
// Polygons whose extents do not overlap are rejected before any edge pair is examined, and so
// is every edge whose extent lies outside the other polygon.
func (p Polygon[T]) Intersects(polygon Polygon[T]) bool {
	if p.Empty() || polygon.Empty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.encloses(p.Points[0], a2, b2) || p.encloses(polygon.Points[0], a1, b1) {
		return true
	}

	for edge := range p.Edges() {
		if polygon.crossesEdge(edge, a2, b2) {
			return true
		}
	}

	return false
}

// IntersectsLine reports whether the polygon and the segment share a point: an endpoint lies
// within the polygon, or the segment crosses one of its edges. Touching shapes intersect,
// within Epsilon of T. A segment whose extent lies outside the polygon is rejected before any
// edge is examined.
func (p Polygon[T]) IntersectsLine(line Line[T]) bool {
	if p.Empty() {
		return false
	}

	a, b := p.minMax()

	return p.encloses(line.Start, a, b) || p.crossesEdge(line, a, b)
}

// IntersectionLine returns the points where the segment crosses the polygon boundary, from
// the segment's Start to its End: the crossings with its edges by Line.Intersection, with a
// vertex hit by two edges counted once. A segment inside crosses no boundary and returns none
// while IntersectsLine still reports it, a segment along an edge is parallel to it and crosses
// only the edges at its ends, and an empty polygon has no boundary to cross.
func (p Polygon[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.crossings(p.Points)
}

// IntersectsRectangle reports whether the polygon and the rectangle share a point, the same
// answer as Intersects on the rectangle's Polygon, without building it: a corner of one lies
// within the other, or an edge of the rectangle crosses an edge of the polygon. The
// rectangle's extent rejects it before any edge is examined, whatever its angle.
func (p Polygon[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	if p.Empty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := rectangle.MinMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if rectangle.Contains(p.Points[0]) || p.encloses(rectangle.TopLeft(), a1, b1) {
		return true
	}

	for edge := range rectangle.Edges() {
		if p.crossesEdge(edge, a1, b1) {
			return true
		}
	}

	return false
}

// IntersectsCircle reports whether the polygon and the circle share a point: the center lies
// within the polygon, or an edge passes within the radius. Touching shapes intersect, within
// Epsilon of T, by the same comparison Circle.Contains makes on the squared distance Contains
// and DistanceSquaredTo measure. A circle whose Bounds lie outside the polygon is rejected
// before any edge is examined.
func (p Polygon[T]) IntersectsCircle(circle Circle[T]) bool {
	if p.Empty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := circle.Bounds().MinMax()

	return overlaps(a1, b1, a2, b2) && circle.reaches(p.walk(circle.Center))
}

// crossesEdge reports whether the segment intersects any edge of the polygon, given the extent
// of the vertices. The segment is rejected by its own extent against that box, and each edge by
// its extent against the segment, so the segment test runs only on edges that can share a point.
func (p Polygon[T]) crossesEdge(line Line[T], a, b Point[T]) bool {
	c, d := line.minMax()
	if !overlaps(a, b, c, d) {
		return false
	}

	for edge := range p.Edges() {
		if e, f := edge.minMax(); overlaps(c, d, e, f) && edge.Intersects(line) {
			return true
		}
	}

	return false
}

// Equal checks if two polygons have the same vertices. A nil and an empty Points are equal,
// having the same none; only IsZero tells them apart.
func (p Polygon[T]) Equal(polygon Polygon[T]) bool {
	if len(p.Points) != len(polygon.Points) {
		return false
	}

	for i, v := range p.Points {
		if !v.Equal(polygon.Points[i]) {
			return false
		}
	}

	return true
}

// IsZero checks if the vertices slice is nil.
func (p Polygon[T]) IsZero() bool {
	return p.Points == nil
}

// Empty checks if number of vertices is zero.
func (p Polygon[T]) Empty() bool {
	return len(p.Points) == 0
}

// Cast converts the polygon to a Polygon of another number type, rounding as Cast does.
func (p Polygon[T]) Cast[R Number]() Polygon[R] {
	return Polygon[R]{xslices.Map(p.Points, Point[T].Cast[R])}
}

// Int converts the polygon to a Polygon[int].
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{xslices.Map(p.Points, Point[T].Int)}
}

// Float converts the polygon to a Polygon[float64].
func (p Polygon[T]) Float() Polygon[float64] {
	return Polygon[float64]{xslices.Map(p.Points, Point[T].Float)}
}

// String returns the polygon in the form of its constructor: Pol((x,y);(x,y);...), the
// vertices separated as every other shape separates its fields.
func (p Polygon[T]) String() string {
	return fmt.Sprintf("Pol(%s)", strings.Join(xslices.Map(p.Points, Point[T].String), ";"))
}

// MarshalJSON implements json.Marshaler.
func (p Polygon[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON implements json.Unmarshaler. The vertices are decoded into a fresh slice, so a
// slice the polygon shared before decoding is left untouched.
func (p *Polygon[T]) UnmarshalJSON(bytes []byte) error {
	var vertices []Point[T]
	if err := json.Unmarshal(bytes, &vertices); err != nil {
		return err
	}

	p.Points = vertices

	return nil
}

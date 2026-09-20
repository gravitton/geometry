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

// Edges iterates the polygon edges in vertex order, each from a vertex to the next and the
// last one closing back to the first, without allocating; collect them with slices.Collect
// where a slice is needed. A single vertex yields one zero-length edge and an empty polygon
// yields nothing. Every walk over the outline, containment, distance and the crossings of a
// segment, reads these edges, so the boundary they join is the one every test agrees on.
func (p Polygon[T]) Edges() iter.Seq[Segment[T]] {
	return edgesOf(p.Points)
}

// Vertices iterates Points in order, the form every shape with an outline offers, so a
// polygon is drawn or measured by the same loop as a Rectangle or a RegularPolygon. Index
// Points directly where a position is needed.
func (p Polygon[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// Centroid returns the center of the enclosed area, so a vertex added in
// the middle of an edge does not move it. A polygon that encloses no area, with fewer than
// three vertices or all of them collinear, has no such center and falls back to the average of
// its vertices; an empty polygon returns the zero point.
// For integer T the centroid is rounded like every other result stored into T; use float64
// when centroid accuracy matters.
//
// The sum is taken with the origin moved to the first vertex, so the products stay small and
// every edge at that vertex contributes exactly zero: a polygon of one or two vertices has an
// exact zero area whatever its coordinates.
func (p Polygon[T]) Centroid() Point[T] {
	if p.IsEmpty() {
		return Point[T]{}
	}

	origin := p.Points[0].Float()
	offset := origin.Vector().Negate()

	var x, y, twiceArea float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		cross := shifted.cross()

		x += (shifted.Start.X + shifted.End.X) * cross
		y += (shifted.Start.Y + shifted.End.Y) * cross
		twiceArea += cross
	}

	if twiceArea == 0 {
		return p.mean()
	}

	centroid := origin.AddXY(x/(3*twiceArea), y/(3*twiceArea))

	return centroid.Cast[T]()
}

// Area returns the area enclosed by the polygon, by the shoelace formula, regardless of winding.
// It is a float64 even for an integer T, since a lattice polygon can enclose half a unit;
// a self-intersecting polygon has its lobes cancel where they wind the opposite way.
func (p Polygon[T]) Area() float64 {
	var twiceArea float64
	for edge := range p.Edges() {
		twiceArea += edge.cross()
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

// Inertia returns the polar second moment of area about the centroid, the rotational inertia
// of the enclosed area at unit density, summed per edge like Area and Centroid: about the
// first vertex, so the products stay small, then moved to the centroid by the parallel axis
// theorem. Winding does not matter, and a polygon that encloses no area has no moment.
func (p Polygon[T]) Inertia() float64 {
	if p.IsEmpty() {
		return 0
	}

	offset := p.Points[0].Vector().Float().Negate()

	var x, y, twiceArea, moment float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		a, b := shifted.Start.Vector(), shifted.End.Vector()
		cross := shifted.cross()

		x += (a.X + b.X) * cross
		y += (a.Y + b.Y) * cross
		twiceArea += cross
		moment += (a.Dot(a) + a.Dot(b) + b.Dot(b)) * cross
	}

	if twiceArea == 0 {
		return 0
	}

	centroid := Vector[float64]{x / (3 * twiceArea), y / (3 * twiceArea)}

	return math.Abs(moment)/12 - math.Abs(twiceArea)/2*centroid.LengthSquared()
}

// Bounds returns the axis-aligned bounding rectangle of the vertices, or the zero rectangle
// for a polygon without vertices.
func (p Polygon[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(p.minMax())
}

// mean returns the average of the vertices, which Centroid falls back to when the
// polygon encloses no area.
func (p Polygon[T]) mean() Point[T] {
	var x, y float64
	for _, vertex := range p.Points {
		x, y = x+float64(vertex.X), y+float64(vertex.Y)
	}

	n := float64(len(p.Points))

	return Point[T]{Cast[T](x / n), Cast[T](y / n)}
}

// minMax returns the minimum and maximum corner of the vertices, the corners of Bounds, as
// minMaxOf finds them; an empty polygon returns two zero points.
func (p Polygon[T]) minMax() (Point[T], Point[T]) {
	return minMaxOf(p.Points)
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
	return p.Translate(point.Subtract(p.Centroid()))
}

// Scale creates a new Polygon uniformly scaled about its centroid by the factor.
func (p Polygon[T]) Scale(factor float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Unscale creates a new Polygon uniformly scaled about its centroid by the inverse factor, the
// inverse of Scale. Like Divide it panics for a zero factor.
func (p Polygon[T]) Unscale(factor float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Divide(factor))
	})}
}

// UnscaleXY creates a new Polygon scaled about its centroid by the inverse of the given
// factors, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (p Polygon[T]) UnscaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

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
	pivot := p.Centroid()

	return Polygon[T]{xslices.Map(p.Points, func(point Point[T]) Point[T] {
		return point.RotateAround(pivot, angle)
	})}
}

// Contains reports whether the given point lies within the polygon, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains. The interior follows the
// even-odd rule, so a self-intersecting polygon excludes the regions it winds around twice.
// A point outside the extent of the vertices is rejected before any edge is examined.
func (p Polygon[T]) Contains(point Point[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a, b := p.minMax()

	return p.containsWithin(point, a, b)
}

// DistanceTo returns the distance from the given point to the nearest point of the polygon:
// zero for a point within it, the same closed convention as Contains, and otherwise the
// distance to the nearest edge. An empty polygon is infinitely far from every point.
func (p Polygon[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(p.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, in one pass over the edges, the edgeWalk every closed shape makes: zero for a
// point on an edge within Epsilon of T, snapped the way Segment.DistanceSquaredTo snaps it, or
// inside by the even-odd rule, the squared distance to the nearest edge otherwise, and
// infinity for an empty polygon. It is a float64 even for an integer T, since the nearest
// point of an edge is not a lattice point in general. Contains and IntersectsCircle are built
// on it, so the three agree by construction.
func (p Polygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range p.Edges() {
		if w.step(edge, point) {
			return 0
		}
	}

	return w.result()
}

// IntersectsCircle reports whether the polygon and the circle share a point, as
// Circle.IntersectsPolygon does.
func (p Polygon[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsPolygon(p)
}

// IntersectsSegment reports whether the polygon and the segment share a point, as
// Segment.IntersectsPolygon does.
func (p Polygon[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsPolygon(p)
}

// IntersectionSegment returns the points where the segment crosses the polygon boundary, as
// Segment.IntersectionPolygon does.
func (p Polygon[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionPolygon(p)
}

// IntersectsPolygon reports whether the polygons share a point: a vertex of one lies within the other,
// or an edge of one crosses an edge of the other. Touching polygons intersect, within Epsilon
// of T, the same closed convention as Contains, and an empty polygon intersects nothing.
// Polygons whose extents do not overlap are rejected before any edge pair is examined, and so
// is every edge whose extent lies outside the other polygon.
func (p Polygon[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if p.IsEmpty() || polygon.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(p.Points[0], a2, b2) || p.containsWithin(polygon.Points[0], a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range polygon.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectsRectangle reports whether the polygon and the rectangle share a point, the same
// answer as Intersects on the rectangle's Polygon, without building it: a corner of one lies
// within the other, or an edge of the rectangle crosses an edge of the polygon. The
// rectangle's extent rejects it before any edge is examined, whatever its angle.
func (p Polygon[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := rectangle.MinMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if rectangle.containsWithin(p.Points[0], a2, b2) || p.containsWithin(rectangle.TopLeft(), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range rectangle.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectsRegularPolygon reports whether the polygon and the regular polygon share a point,
// the same answer as IntersectsPolygon on the regular polygon's Polygon, without building it: a
// vertex of one lies within the other, or an edge of the regular polygon crosses an edge of
// the polygon. The two Bounds reject the pair before any edge is examined, and an empty
// polygon on either side intersects nothing.
func (p Polygon[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if p.IsEmpty() || polygon.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(p.Points[0], a2, b2) || p.containsWithin(polygon.vertex(0), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range polygon.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// containsWithin is Contains for a caller that already holds the extent of the vertices, so the
// intersection tests walk the vertices once for the box and reuse it for every point they test.
func (p Polygon[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && p.DistanceSquaredTo(point) == 0
}

// Equal checks if two polygons have the same vertices. A nil and an empty Points are equal,
// having the same none; only IsZero tells them apart.
func (p Polygon[T]) Equal(polygon Polygon[T]) bool {
	if len(p.Points) != len(polygon.Points) {
		return false
	}

	for i, vertex := range p.Points {
		if !vertex.Equal(polygon.Points[i]) {
			return false
		}
	}

	return true
}

// IsZero checks if the vertices slice is nil.
func (p Polygon[T]) IsZero() bool {
	return p.Points == nil
}

// IsEmpty checks if number of vertices is zero.
func (p Polygon[T]) IsEmpty() bool {
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

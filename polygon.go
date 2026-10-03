package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"
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

// Vertices iterates Points in order, the form every shape with an outline offers, so a
// polygon is drawn or measured by the same loop as a Rectangle or a RegularPolygon. Index
// Points directly where a position is needed.
func (p Polygon[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// Edges iterates the polygon edges in vertex order, each from a vertex to the next and the
// last one closing back to the first, without allocating; collect them with slices.Collect
// where a slice is needed. A single vertex yields one zero-length edge and an empty polygon
// yields nothing. Every walk over the outline, containment, distance and the crossings of a
// segment, reads these edges, so the boundary they join is the one every test agrees on.
func (p Polygon[T]) Edges() iter.Seq[Segment[T]] {
	return edgesOf(p.Points)
}

// Centroid returns the center of the enclosed area, so a vertex added in
// the middle of an edge does not move it. A polygon that encloses no area, with fewer than
// three vertices or all of them on one line within the tolerance, or with lobes that cancel
// exactly, has no such center and falls back to the average of its vertices; an empty polygon
// returns the zero point. The line is judged on the distance, as Contains judges a boundary, so
// vertices placed along a line in float64, each rounded a hair off it, are on it rather than
// around a sliver whose rounding residue of an area would put the centroid anywhere.
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

	var x, y, twiceArea, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		along, cross := shifted.Vector(), shifted.cross()

		x += float64((shifted.Start.X + shifted.End.X) * cross)
		y += float64((shifted.Start.Y + shifted.End.Y) * cross)
		twiceArea += cross
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if twiceArea == 0 || p.isFlatWithin(twiceArea, length) {
		return p.mean()
	}

	centroid := origin.AddXY(x/(3*twiceArea), y/(3*twiceArea))

	return centroid.Cast[T]()
}

// Area returns the area enclosed by the polygon, by the shoelace formula, regardless of winding.
// It is a float64 even for an integer T, since a lattice polygon can enclose half a unit;
// a self-intersecting polygon has its lobes cancel where they wind the opposite way, and a
// polygon whose vertices lie on one line within the tolerance encloses none.
//
// The sum is taken with the origin moved to the first vertex, as Centroid and Inertia take
// theirs, so the products stay small and a polygon far from the origin keeps its area rather
// than losing it between two large coordinates; an empty polygon encloses nothing.
func (p Polygon[T]) Area() float64 {
	return math.Abs(p.twiceArea()) / 2
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
// theorem. Winding does not matter, and a polygon that encloses no area, its vertices on one
// line within the tolerance as Centroid judges it, has no moment.
//
// The moment is the one a simple polygon has. Where the outline crosses itself the two sums
// cancel by different amounts, the way Area has its lobes cancel, and what is left is no
// longer a moment of area: it can come out negative, which no area about its own centroid has.
func (p Polygon[T]) Inertia() float64 {
	if p.IsEmpty() {
		return 0
	}

	offset := p.Points[0].Vector().Float().Negate()

	var x, y, twiceArea, moment, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		a, b := shifted.Start.Vector(), shifted.End.Vector()
		along, cross := shifted.Vector(), shifted.cross()

		x += float64((a.X + b.X) * cross)
		y += float64((a.Y + b.Y) * cross)
		twiceArea += cross
		moment += float64((a.Dot(a) + a.Dot(b) + b.Dot(b)) * cross)
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if twiceArea == 0 || p.isFlatWithin(twiceArea, length) {
		return 0
	}

	centroid := Vector[float64]{x / (3 * twiceArea), y / (3 * twiceArea)}

	return math.Abs(moment)/12 - float64(math.Abs(twiceArea)/2*centroid.LengthSquared())
}

// Winding returns the sense in which the vertices run around the area they enclose, from the
// sign of the sum Area takes the absolute value of: WindingClockwise for the winding of
// Rectangle.Vertices, WindingCounterClockwise for the reverse, and WindingNone for a polygon
// that encloses no area, its vertices on one line within the tolerance as Centroid judges it. A self-intersecting polygon winds the way its larger lobes do, and
// lobes of equal area either way give WindingNone.
func (p Polygon[T]) Winding() Winding {
	switch twiceArea := p.twiceArea(); {
	case twiceArea > 0:
		return WindingClockwise
	case twiceArea < 0:
		return WindingCounterClockwise
	default:
		return WindingNone
	}
}

// Bounds returns the axis-aligned bounding box of the vertices, or the zero box for a polygon
// without vertices.
func (p Polygon[T]) Bounds() Box[T] {
	a, b := p.minMax()

	return Box[T]{a, b}
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

// twiceArea returns the signed sum of the shoelace formula, twice the enclosed area, positive
// for a clockwise winding: Area takes its absolute value and Winding its sign. The sum is
// taken with the origin moved to the first vertex; a flat polygon, an empty one included, sums
// to zero.
func (p Polygon[T]) twiceArea() float64 {
	if p.IsEmpty() {
		return 0
	}

	offset := p.Points[0].Vector().Float().Negate()

	var twiceArea, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		along := shifted.Vector()

		twiceArea += shifted.cross()
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if p.isFlatWithin(twiceArea, length) {
		return 0
	}

	return twiceArea
}

// minMax returns the minimum and maximum corner of the vertices, the corners of Bounds, as
// minMaxOf finds them; an empty polygon returns two zero points.
func (p Polygon[T]) minMax() (Point[T], Point[T]) {
	return minMaxOf(p.Points)
}

// Translate creates a new Polygon translated by the given vector (applied to all vertices).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
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

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Unscale creates a new Polygon uniformly scaled about its centroid by the inverse factor, the
// inverse of Scale. Like Divide it panics for a zero factor.
func (p Polygon[T]) Unscale(factor float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Divide(factor))
	})}
}

// UnscaleXY creates a new Polygon scaled about its centroid by the inverse of the given
// factors, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (p Polygon[T]) UnscaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).DivideXY(factorX, factorY))
	})}
}

// Lerp creates a new Polygon in linear interpolation towards the given polygon, each vertex
// moving towards the vertex of the same index like Point.Lerp, extrapolating outside [0, 1]. The
// vertices pair by index, so two outlines that start at different vertices or run in opposite
// windings twist through each other on the way; align them first where that matters. Polygons
// with a different vertex count have no shape between them, so Lerp panics for them, as
// RegularPolygon.Lerp does for a different N.
func (p Polygon[T]) Lerp(polygon Polygon[T], t float64) Polygon[T] {
	if len(p.Points) != len(polygon.Points) {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", len(p.Points), len(polygon.Points)))
	}

	return Polygon[T]{p.mapPointsIndexed(func(i int, point Point[T]) Point[T] {
		return point.Lerp(polygon.Points[i], t)
	})}
}

// mapPoints applies fn to every vertex in order and returns the results in a new slice, as
// mapPointsIndexed does for a mapping that needs the index of the vertex.
func (p Polygon[T]) mapPoints[R any](fn func(Point[T]) R) []R {
	return p.mapPointsIndexed(func(_ int, point Point[T]) R {
		return fn(point)
	})
}

// mapPointsIndexed applies fn to every vertex and its index in order and returns the results
// in a new slice, the one allocation of every mapping of the polygon. A nil Points maps to nil,
// so IsZero holds through every mapping.
func (p Polygon[T]) mapPointsIndexed[R any](fn func(int, Point[T]) R) []R {
	if p.Points == nil {
		return nil
	}

	mapped := make([]R, len(p.Points))
	for i, point := range p.Points {
		mapped[i] = fn(i, point)
	}

	return mapped
}

// Transform creates a new Polygon by applying the given matrix to every vertex, like Point.Transform.
func (p Polygon[T]) Transform[M Float](matrix Matrix[M]) Polygon[T] {
	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return point.Transform(matrix)
	})}
}

// Rotate creates a new Polygon rotated by the given angle (in radians) about its centroid, in
// the same sense as Vector.Rotate. For integer T the centroid and every rotated vertex are
// rounded; only multiples of 90° keep the shape exactly.
func (p Polygon[T]) Rotate(angle float64) Polygon[T] {
	pivot := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return point.RotateAround(pivot, angle)
	})}
}

// ConvexHull returns the smallest convex polygon containing every vertex, wound clockwise
// like Rectangle.Vertices and starting at the least vertex by Point.Compare: the vertices it
// keeps are vertices of the polygon, and one lying on an edge of the hull or repeating
// another is dropped, so a polygon whose vertices are all collinear gives the two ends of
// the line and a single point gives itself. The turns are decided on exact signs, as
// IsConvex decides them, so the hull of a convex polygon is convex again and the hull of a
// hull is itself. An empty polygon gives an empty one.
//
// It sorts a copy of the vertices along the boundary, the side of the line between the least
// and the greatest vertex that the hull reaches first and then the other side back, and scans
// that copy once, compacting the hull into its front, so the copy is the one allocation.
func (p Polygon[T]) ConvexHull() Polygon[T] {
	return Polygon[T]{p.AppendConvexHull(nil)}
}

// AppendConvexHull appends the vertices of the polygon ConvexHull returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for every vertex
// of the polygon: the copy is sorted and compacted within the appended tail, and the points
// already in dst are left as they are. An empty polygon appends nothing.
func (p Polygon[T]) AppendConvexHull(dst []Point[T]) []Point[T] {
	if p.IsEmpty() {
		return dst
	}

	chord := Segment[T]{slices.MinFunc(p.Points, Point[T].Compare), slices.MaxFunc(p.Points, Point[T].Compare)}
	from := len(dst)
	dst = append(dst, p.Points...)
	points := dst[from:]
	slices.SortFunc(points, func(a, b Point[T]) int {
		return p.compareAround(chord, a, b)
	})

	n := 0
	for _, point := range points {
		if n > 0 && point == points[n-1] {
			continue
		}

		for n >= 2 && p.turn(points[n-2], points[n-1], point) <= 0 {
			n--
		}

		points[n] = point
		n++
	}

	for n >= 3 && p.turn(points[n-2], points[n-1], points[0]) <= 0 {
		n--
	}

	return dst[:from+n]
}

// Simplify returns the polygon without the vertices the outline does not need to stay within
// the tolerance, by the Douglas–Peucker algorithm: every dropped vertex lies within the
// tolerance of the edge that replaces it, by Segment.DistanceTo, boundary included within
// the tolerance. At a tolerance of zero it drops a vertex repeating a neighbour or on the
// straight line between the vertices kept either side, a fold running back along that line
// included, and keeps a spike reaching beyond it.
//
// The outline is taken as a chain from its least vertex by Point.Compare around and back to
// it, so the first split is at the vertex farthest from that one: both are extreme, so neither
// is kept that the tolerance would drop. Each chain keeps the vertex farthest from the edge
// joining its ends where it lies beyond the tolerance, and is split there; one lying within it
// drops every vertex between the ends. An outline within the tolerance of its least vertex
// keeps that vertex alone. The kept vertices are vertices of the polygon in its order, from
// the first kept one, so a polygon with nothing to drop is returned equal to itself and
// simplifying again changes nothing. A negative tolerance is taken absolute, and an empty
// polygon gives an empty one.
//
// The kept vertices are appended to one slice with room for all of them, the one allocation,
// and turned in place to start where the polygon does.
func (p Polygon[T]) Simplify(tolerance float64) Polygon[T] {
	return Polygon[T]{p.AppendSimplify(nil, tolerance)}
}

// AppendSimplify appends the vertices of the polygon Simplify returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for every vertex
// of the polygon: the kept vertices are turned within the appended tail, and the points already
// in dst are left as they are. An empty polygon appends nothing. The vertices are read while the
// tail is written, so dst must not share memory with Points.
func (p Polygon[T]) AppendSimplify(dst []Point[T], tolerance float64) []Point[T] {
	if p.IsEmpty() {
		return dst
	}

	n, least := len(p.Points), 0
	for i, vertex := range p.Points {
		if vertex.Compare(p.Points[least]) < 0 {
			least = i
		}
	}

	from := len(dst)
	if cap(dst)-from < n {
		dst = append(make([]Point[T], 0, from+n), dst...)
	}

	dst, wrapped := p.appendKept(append(dst, p.Points[least]), least, least+n, math.Abs(tolerance))

	kept := dst[from:]
	slices.Reverse(kept)
	slices.Reverse(kept[:wrapped])
	slices.Reverse(kept[wrapped:])

	return dst
}

// compareAround orders two points along the boundary of the hull whose least and greatest
// vertex the chord joins: the points on the side of the chord a clockwise walk from Start
// reaches first, or on the chord, by Point.Compare, then the points on the other side in the
// reverse order. It reads no field of the polygon, only the chord and the points, so the
// receiver is unnamed.
func (Polygon[T]) compareAround(chord Segment[T], a, b Point[T]) int {
	start, end := chord.Start.Float(), chord.End.Float()
	direction := end.Subtract(start)
	returning := func(point Point[T]) bool {
		return direction.Cross(point.Float().Subtract(start)) > 0
	}

	switch ra, rb := returning(a), returning(b); {
	case ra != rb && ra:
		return 1
	case ra != rb:
		return -1
	case ra:
		return b.Compare(a)
	default:
		return a.Compare(b)
	}
}

// turn returns the cross product of the edge from a to b and the edge from b to c, in
// float64: positive where the outline turns clockwise at b, the winding of Rectangle.Vertices,
// negative the other way and zero where the three are collinear. It reads no field of the
// polygon, only the points, so the receiver is unnamed.
func (Polygon[T]) turn(a, b, c Point[T]) float64 {
	return b.Float().Subtract(a.Float()).Cross(c.Float().Subtract(b.Float()))
}

// appendKept appends to dst the vertices Simplify keeps strictly between the vertices at from
// and to, indices counted on around the outline past the last one: the vertex farthest from
// the edge joining the two, the first of equals, where it lies beyond the tolerance, with the
// vertices kept between it and either end. It returns the extended slice and how many of the
// appended vertices lie past the last index, at the front of the polygon.
func (p Polygon[T]) appendKept(dst []Point[T], from, to int, tolerance float64) ([]Point[T], int) {
	n := len(p.Points)
	edge := Segment[T]{p.Points[from%n], p.Points[to%n]}

	farthest, distance := from, 0.0
	for i := from + 1; i < to; i++ {
		if d := edge.DistanceSquaredTo(p.Points[i%n]); d > distance {
			farthest, distance = i, d
		}
	}

	if lessOrEqualSquared(distance, tolerance, epsilonAt[T](max(edge.magnitude(), p.Points[farthest%n].magnitude()))) {
		return dst, 0
	}

	dst, before := p.appendKept(dst, from, farthest, tolerance)
	dst = append(dst, p.Points[farthest%n])
	dst, after := p.appendKept(dst, farthest, to, tolerance)

	return dst, before + farthest/n + after
}

// Contains reports whether the given point lies within the polygon, boundary included within
// the tolerance, the same closed convention as Rectangle.Contains. The interior follows the
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
// point on an edge within the tolerance, snapped the way Segment.DistanceSquaredTo snaps it, or
// inside by the even-odd rule, the squared distance to the nearest edge otherwise, and
// infinity for an empty polygon. It is a float64 even for an integer T, since the nearest
// point of an edge is not a lattice point in general. Contains and IntersectsCircle are built
// on it, so the three agree by construction.
func (p Polygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	return p.walk(point).result()
}

// Nearest returns the point of the polygon nearest to the given point: the point itself
// exactly where Contains holds, and otherwise the foot on the nearest edge, read off the same
// walk DistanceSquaredTo makes. For integer T the foot is rounded once and can land off the
// boundary, where Contains rejects it. An empty polygon returns the zero point.
func (p Polygon[T]) Nearest(point Point[T]) Point[T] {
	return p.walk(point).nearest(point)
}

// EnclosesCircle reports whether the circle lies within the polygon: its center is contained
// and every edge is at least the radius away, within the tolerance, read off the walk
// DistanceSquaredTo makes, so a circle touching an edge from inside is enclosed. A circle clear
// of every edge cannot reach outside, concave polygon or not. An empty polygon encloses nothing.
func (p Polygon[T]) EnclosesCircle(circle Circle[T]) bool {
	return p.walk(circle.Center).clears(circle)
}

// EnclosesSegment reports whether the segment lies within the polygon: both endpoints are
// contained, within the tolerance, and it leaves the polygon nowhere between them, which a
// concave polygon can let it do. It leaves by properly crossing an edge, decided on exact
// signs as Segment.IntersectsSegment decides a crossing, by passing a vertex into the outside,
// or by running from an endpoint on an edge toward the outer side of it, the side read from
// the Winding of the polygon. A crossing where an endpoint of either lies on the other within
// the tolerance is a touch, and a segment that runs beyond the line of an edge by no more than
// the tolerance runs along it. The polygon is taken to be simple: an edge crossing another edge
// of a self-intersecting polygon counts as leaving it, so such a polygon does not enclose
// itself. An empty polygon encloses nothing.
func (p Polygon[T]) EnclosesSegment(segment Segment[T]) bool {
	a, b := p.minMax()

	return p.containsWithin(segment.Start, a, b) && p.containsWithin(segment.End, a, b) && p.keeps(segment, p.twiceArea())
}

// EnclosesPolygon reports whether the given polygon lies within this one: every edge of it is
// enclosed, as EnclosesSegment decides, and an outline within a simple polygon holds its area
// within it too. An empty polygon on either side encloses nothing.
func (p Polygon[T]) EnclosesPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range polygon.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// EnclosesRectangle reports whether the rectangle lies within the polygon: every edge of it is
// enclosed, as EnclosesSegment decides, whatever its angle.
func (p Polygon[T]) EnclosesRectangle(rectangle Rectangle[T]) bool {
	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range rectangle.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// EnclosesRegularPolygon reports whether the regular polygon lies within this one: every edge
// of it is enclosed, as EnclosesSegment decides. An empty polygon on either side encloses
// nothing.
func (p Polygon[T]) EnclosesRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range polygon.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// EnclosesBox reports whether the box lies within the polygon: every edge of it, between its
// corners, is enclosed, as EnclosesSegment decides.
func (p Polygon[T]) EnclosesBox(box Box[T]) bool {
	a, b := p.minMax()
	twiceArea := p.twiceArea()

	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
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

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns the
// extended slice, as Segment.AppendIntersectionPolygon does.
func (p Polygon[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	return segment.AppendIntersectionPolygon(dst, p)
}

// IntersectsRay reports whether the polygon and the ray share a point, as
// Ray.IntersectsPolygon does.
func (p Polygon[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsPolygon(p)
}

// IntersectionRay returns the points where the ray crosses the polygon boundary, as
// Ray.IntersectionPolygon does.
func (p Polygon[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionPolygon(p)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice, as Ray.AppendIntersectionPolygon does.
func (p Polygon[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	return ray.AppendIntersectionPolygon(dst, p)
}

// IntersectsPolygon reports whether the polygons share a point: a vertex of one lies within the other,
// or an edge of one crosses an edge of the other. Touching polygons intersect, within the
// tolerance, the same closed convention as Contains, and an empty polygon intersects nothing.
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
// answer as IntersectsPolygon on the rectangle's Polygon, without building it: a corner of one lies
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

// IntersectsBox reports whether the polygon and the box share a point, as IntersectsRectangle
// decides it against a rectangle that is not rotated: a corner of one lies within the other,
// or an edge of the box, between its own corners, crosses an edge of the polygon. The box
// rejects it before any edge is examined, and an empty polygon intersects nothing.
func (p Polygon[T]) IntersectsBox(box Box[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a, b := p.minMax()

	if !overlaps(a, b, box.Min, box.Max) {
		return false
	}

	if box.Contains(p.Points[0]) || p.containsWithin(box.Min, a, b) {
		return true
	}

	probe := edgeProbe[T]{a: box.Min, b: box.Max}
	corners := box.corners()
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range edgesOf(corners[:]) {
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

// containsMidpoint is Contains at the midpoint of a and b, a point T cannot always hold: the
// step edgeWalk makes, the distance within the tolerance and the even-odd ray, taken on each
// edge in float64, so for an integer T the midpoint is not rounded onto either point.
// Segment.ClipPolygon judges each piece of a segment between two crossings by it.
func (p Polygon[T]) containsMidpoint(a, b Point[T]) bool {
	midpoint := a.Float().Midpoint(b.Float())

	inside := false
	for edge := range p.Edges() {
		if lessOrEqualSquared(edge.Float().distanceSquaredTo(midpoint), 0, epsilonAt[T](max(edge.magnitude(), midpoint.magnitude()))) {
			return true
		}

		if edge.Float().crossesRay(midpoint) {
			inside = !inside
		}
	}

	return inside
}

// walk folds every edge into the edgeWalk DistanceSquaredTo, Nearest and EnclosesCircle read, stopping
// at an edge the point lies on within the tolerance.
func (p Polygon[T]) walk(point Point[T]) edgeWalk[T] {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range p.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	return w
}

// keeps reports whether the segment stays within the polygon between its endpoints, which the
// caller has found contained. A segment leaves a simple polygon in one of three ways, each
// judged on exact signs in the sense twiceArea gives: it properly crosses an edge, short of a
// touch; it runs through a vertex into the outside; or it runs from an endpoint inside an edge
// toward the outer side of that edge, a boundary point where the outline runs straight. The
// last two are asked of leavesThrough, a vertex or an endpoint is found on the boundary within
// the tolerance, and an endpoint at a vertex is left to the vertex, which judges both its edges.
// A vertex repeating the one before it is skipped, as IsConvex skips it, so a vertex is judged
// by the edges toward its distinct neighbours, never by the zero direction of a repeat, which
// would leave a reflex vertex only one of its sides.
func (p Polygon[T]) keeps(segment Segment[T], twiceArea float64) bool {
	previous := p.Points[0]
	for _, vertex := range slices.Backward(p.Points) {
		if vertex != p.Points[0] {
			previous = vertex

			break
		}
	}

	for edge := range p.Edges() {
		along := edge.Float().Vector()
		if !along.hasDirection() {
			continue
		}

		if segment.crosses(edge) {
			if _, touching := segment.touch(edge); !touching {
				return false
			}
		}

		if vertex := edge.Start; segment.Contains(vertex) {
			if p.leavesThrough(vertex, previous.Float().Subtract(vertex.Float()), along, segment, twiceArea) {
				return false
			}
		}

		for _, endpoint := range [2]Point[T]{segment.Start, segment.End} {
			if !edge.Contains(endpoint) || endpoint.coincides(edge.Start) || endpoint.coincides(edge.End) {
				continue
			}

			if p.leavesThrough(endpoint, along.Negate(), along, segment, twiceArea) {
				return false
			}
		}

		previous = edge.Start
	}

	return true
}

// leavesThrough reports whether the segment, passing a point of the boundary where the outline
// turns from the direction toward the previous vertex to the one toward the next, runs from
// there into the outside toward either endpoint, as admits judges the turn in the sense
// twiceArea gives.
func (p Polygon[T]) leavesThrough(point Point[T], toPrevious, toNext Vector[float64], segment Segment[T], twiceArea float64) bool {
	if twiceArea < 0 {
		toPrevious, toNext = toNext, toPrevious
	}

	start, end, origin := segment.Start.Float(), segment.End.Float(), point.Float()
	epsilon := epsilonAt[T](max(point.magnitude(), segment.magnitude()))

	return !p.admits(toPrevious, toNext, start.Subtract(origin), epsilon) ||
		!p.admits(toPrevious, toNext, end.Subtract(origin), epsilon)
}

// admits reports whether the offset from a point of the boundary stays within the polygon,
// where the outline turns there from the direction toward the previous vertex to the one
// toward the next with the area between them: on the inner side of both edges where the turn
// is convex, straight or a spike, and of either where it is reflex, as stays judges each
// side. The turn is decided on the exact sign of the cross product, as crosses decides.
func (p Polygon[T]) admits(toPrevious, toNext, offset Vector[float64], epsilon float64) bool {
	insideNext, insidePrevious := p.stays(toNext.Cross(offset), toNext, epsilon), p.stays(offset.Cross(toPrevious), toPrevious, epsilon)
	if toNext.Cross(toPrevious) >= 0 {
		return insideNext && insidePrevious
	}

	return insideNext || insidePrevious
}

// stays reports whether an offset lies on the inner side of an edge, given their cross
// product, positive on that side, and the direction of the edge: on it, or beyond the line of
// the edge by no more than the tolerance, on the squared gap Segment.distanceSquaredTo measures
// beside a segment. An offset within the tolerance of the point stays on every side, so a
// segment ending there is not judged by a direction it has not got. It reads no field of the
// polygon, so the receiver is unnamed.
func (Polygon[T]) stays(cross float64, edge Vector[float64], epsilon float64) bool {
	return cross >= 0 || lessOrEqualSquared(cross*cross/edge.LengthSquared(), 0, epsilon)
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

// IsConvex reports whether the polygon is convex and simple: it encloses an area, turns the
// same way at every vertex, in either winding, and goes around once, so a star that turns one
// way throughout is not convex. A vertex repeating the one before it is skipped and a vertex
// on the straight line between its neighbours is allowed, while an edge doubling back along
// the one before it is not. The turns are decided on the exact signs of cross products, with
// no tolerance, so a float vertex a rounding error inside the line of its neighbours makes the
// polygon concave.
func (p Polygon[T]) IsConvex() bool {
	c := edgeConvexity{}
	for edge := range p.Edges() {
		if !c.step(edge.Float().Vector()) {
			return false
		}
	}

	return c.result()
}

// isFlat reports whether the polygon encloses no area within the tolerance: it has fewer than
// three vertices, or every vertex lies on the segment between the least and the greatest by
// Point.Compare, as Segment.Contains judges a point on it, which for vertices on one line are
// its two ends. The tolerance is on the distance of each vertex from that line, never on the
// area, so an outline is flat exactly where none of its vertices stands off the line by more.
// The sums ask it through isFlatWithin, which skips it for every polygon enclosing more.
func (p Polygon[T]) isFlat() bool {
	if len(p.Points) < 3 {
		return true
	}

	chord := Segment[T]{slices.MinFunc(p.Points, Point[T].Compare), slices.MaxFunc(p.Points, Point[T].Compare)}
	for _, vertex := range p.Points {
		if !chord.Contains(vertex) {
			return false
		}
	}

	return true
}

// isFlatWithin is isFlat for a caller that has already summed the shoelace terms and the lengths
// the edges run along the two axes: an outline within the tolerance ε of a line sums to no more
// than 2ε times that length, and the rounding of the sum to less than the second term of the
// bound, so only a sum within it is asked of isFlat, and every polygon enclosing more skips the
// walk over its vertices. No vertex lies farther from the first than that length, so the
// tolerance is taken at the first vertex's coordinates widened by it, never below the one any
// vertex is judged at.
func (p Polygon[T]) isFlatWithin(twiceArea, length float64) bool {
	epsilon := epsilonAt[T](p.Points[0].magnitude() + length)
	bound := float64(4*epsilon) + float64(float64(len(p.Points))*0x1p-48*length)

	return math.Abs(twiceArea) <= float64(bound*length) && p.isFlat()
}

// Cast converts the polygon to a Polygon of another number type, rounding as Cast does.
func (p Polygon[T]) Cast[R Number]() Polygon[R] {
	return Polygon[R]{p.mapPoints(Point[T].Cast[R])}
}

// Int converts the polygon to a Polygon[int].
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{p.mapPoints(Point[T].Int)}
}

// Float converts the polygon to a Polygon[float64].
func (p Polygon[T]) Float() Polygon[float64] {
	return Polygon[float64]{p.mapPoints(Point[T].Float)}
}

// String returns the polygon in the form of its constructor: Pol((x,y);(x,y);...), the
// vertices separated as every other shape separates its fields.
func (p Polygon[T]) String() string {
	return fmt.Sprintf("Pol(%s)", strings.Join(p.mapPoints(Point[T].String), ";"))
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

package geom

import (
	"cmp"
	"fmt"
	"iter"
	"math"
	"slices"
)

// Segment is a 2D line segment, bounded by Start and End.
type Segment[T Number] struct {
	Start Point[T] `json:"s"`
	End   Point[T] `json:"e"`
}

// Seg is shorthand for Segment{start, end}.
func Seg[T Number](start, end Point[T]) Segment[T] {
	return Segment[T]{start, end}
}

// Vector returns the segment as a vector, from start to end.
func (s Segment[T]) Vector() Vector[T] {
	return s.End.Subtract(s.Start)
}

// Length returns the length of the segment.
func (s Segment[T]) Length() float64 {
	return s.Vector().Length()
}

// Angle returns the angle of the segment in radians, the angle of the vector from Start to End.
// A zero-length segment has no direction and gives 0, the angle Vector.Angle gives it.
func (s Segment[T]) Angle() float64 {
	return s.Vector().Angle()
}

// Direction returns the direction nearest to the segment, from Start to End, or DirectionNone
// for a zero-length segment, as Vector.Direction judges it.
func (s Segment[T]) Direction() Direction {
	return s.Vector().Direction()
}

// Vertices iterates the start and end points, in that order, without allocating; collect
// them with slices.Collect where a slice is needed.
func (s Segment[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		_ = yield(s.Start) && yield(s.End)
	}
}

// Edges iterates the one edge of the segment, itself: a segment is an open outline, so
// unlike the closed shapes no edge returns to the first vertex.
func (s Segment[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		yield(s)
	}
}

// Midpoint returns the midpoint of the segment, Lerp(0.5).
func (s Segment[T]) Midpoint() Point[T] {
	return s.Start.Midpoint(s.End)
}

// Bounds returns the axis-aligned bounding rectangle.
func (s Segment[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(s.minMax())
}

// minMax returns the minimum and maximum corner of the segment, the corners of Bounds, exact
// for an integer T where Bounds places a center: the pair the intersection tests reject shapes
// by before examining any edge, without placing a rectangle.
func (s Segment[T]) minMax() (Point[T], Point[T]) {
	return Point[T]{min(s.Start.X, s.End.X), min(s.Start.Y, s.End.Y)}, Point[T]{max(s.Start.X, s.End.X), max(s.Start.Y, s.End.Y)}
}

// cross returns Start × End in float64, the term the shoelace formula sums per edge. Cross is
// exact for parallel vectors, so a degenerate edge contributes exactly zero.
func (s Segment[T]) cross() float64 {
	return s.Start.Vector().Float().Cross(s.End.Vector().Float())
}

// Translate creates a new Segment translated by the given vector.
func (s Segment[T]) Translate(vector Vector[T]) Segment[T] {
	return Segment[T]{s.Start.Add(vector), s.End.Add(vector)}
}

// MoveTo creates a new Segment with its midpoint moved to point and the same length and
// direction, the center every shape places with MoveTo and the pivot Scale, Resize and Rotate
// turn about. For integer T the midpoint is rounded, so an odd span lands within half a unit
// of the point, on the side Midpoint rounds to.
func (s Segment[T]) MoveTo(point Point[T]) Segment[T] {
	return s.Translate(point.Subtract(s.Midpoint()))
}

// Scale creates a new Segment uniformly scaled about its midpoint by the factor: the midpoint and
// direction stay, the length multiplies. A zero factor collapses the segment onto its midpoint.
// For integer T the midpoint and both scaled points are rounded, so an odd span scaled by one
// is not exactly the same segment.
func (s Segment[T]) Scale(factor float64) Segment[T] {
	return s.ScaleXY(factor, factor)
}

// ScaleXY creates a new Segment scaled about its midpoint by the factors along X and Y, which
// changes the direction unless the factors are equal.
func (s Segment[T]) ScaleXY(factorX, factorY float64) Segment[T] {
	pivot := s.Midpoint()

	return Segment[T]{pivot.Add(s.Start.Subtract(pivot).MultiplyXY(factorX, factorY)), pivot.Add(s.End.Subtract(pivot).MultiplyXY(factorX, factorY))}
}

// Unscale creates a new Segment uniformly scaled about its midpoint by the inverse factor, the
// inverse of Scale. Like Divide it panics for a zero factor.
func (s Segment[T]) Unscale(factor float64) Segment[T] {
	return s.UnscaleXY(factor, factor)
}

// UnscaleXY creates a new Segment scaled about its midpoint by the inverse of the given factors,
// the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (s Segment[T]) UnscaleXY(factorX, factorY float64) Segment[T] {
	pivot := s.Midpoint()

	return Segment[T]{pivot.Add(s.Start.Subtract(pivot).DivideXY(factorX, factorY)), pivot.Add(s.End.Subtract(pivot).DivideXY(factorX, factorY))}
}

// Resize creates a new Segment of the given length about its midpoint, where Scale multiplies the
// length it has: the midpoint and direction stay and both ends move to half the length either
// side. A zero-length segment has no direction and resizes along +X, the convention
// Vector.Resize follows, and a negative length flips the ends, giving the reverse of the
// segment of that absolute length.
// The length is a float64 like the one Length returns, so an integer segment can be resized to
// a length no integer expresses; for integer T the midpoint and both ends are rounded, and the
// actual length may differ from the requested value.
func (s Segment[T]) Resize(length float64) Segment[T] {
	pivot := s.Midpoint().Float()
	half := s.Vector().Float().Resize(length / 2)

	start, end := pivot.Add(half.Negate()), pivot.Add(half)

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Reverse creates a new Segment with the start and end points swapped.
func (s Segment[T]) Reverse() Segment[T] {
	return Segment[T]{s.End, s.Start}
}

// Lerp returns the point at the fraction t of the way from Start to End, extrapolating along
// the segment outside [0, 1] like Point.Lerp. Midpoint is Lerp(0.5).
func (s Segment[T]) Lerp(t float64) Point[T] {
	return s.Start.Lerp(s.End, t)
}

// Transform creates a new Segment by applying the given matrix to both points, like Point.Transform.
func (s Segment[T]) Transform[M Float](matrix Matrix[M]) Segment[T] {
	return Segment[T]{s.Start.Transform(matrix), s.End.Transform(matrix)}
}

// Rotate creates a new Segment rotated by the given angle (in radians) about its midpoint, in the
// same sense as Vector.Rotate. For integer T the midpoint and both rotated points are rounded;
// only multiples of 90° keep the length exactly.
func (s Segment[T]) Rotate(angle float64) Segment[T] {
	pivot := s.Midpoint()

	return Segment[T]{s.Start.RotateAround(pivot, angle), s.End.RotateAround(pivot, angle)}
}

// Normal returns the perpendicular of the segment, Vector.Normal of its vector: the direction
// from Start to End turned a quarter turn in the sense of Vector.Rotate, clockwise as drawn on
// a screen with Y pointing down, with the length of the segment. On an edge of a Rectangle, or
// of any polygon wound the same way, it points inward. A zero-length segment has no normal and
// gives the zero vector.
func (s Segment[T]) Normal() Vector[T] {
	return s.Vector().Normal()
}

// Contains reports whether the given point lies on the segment, within Epsilon of T, the same
// closed convention as Rectangle.Contains: it holds exactly where DistanceTo is zero.
func (s Segment[T]) Contains(point Point[T]) bool {
	return s.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the segment:
// zero exactly where Contains holds, so a point within Epsilon of T of the segment is at
// distance zero rather than at the rounding error that put it there. A lattice point on a
// lattice segment is at zero without any tolerance, since the perpendicular distance comes
// from a cross product that is exact in float64 rather than from a projection.
func (s Segment[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(s.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. It is a float64 even for an integer T, unlike Point.DistanceSquaredTo, since the
// nearest point of a segment is not a lattice point in general. A distance within Epsilon of T
// is snapped to zero, which is where Contains reads it, so a polygon boundary is walked without
// a square root per edge.
func (s Segment[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := s.distanceSquaredTo(point)
	if lessOrEqualSquared[T](distance, 0) {
		return 0
	}

	return distance
}

// DistanceToSegment returns the distance between the nearest points of the two segments: zero
// exactly where Intersects holds, and otherwise the smallest distance from an endpoint of one
// to the other.
func (s Segment[T]) DistanceToSegment(segment Segment[T]) float64 {
	return math.Sqrt(s.DistanceSquaredToSegment(segment))
}

// DistanceSquaredToSegment returns the squared distance DistanceToSegment takes the root of, faster
// for comparisons: zero where the segments properly cross or an endpoint of one lies on the
// other within Epsilon of T, as DistanceSquaredTo snaps it.
func (s Segment[T]) DistanceSquaredToSegment(segment Segment[T]) float64 {
	if s.crosses(segment) {
		return 0
	}

	return min(
		s.DistanceSquaredTo(segment.Start), s.DistanceSquaredTo(segment.End),
		segment.DistanceSquaredTo(s.Start), segment.DistanceSquaredTo(s.End),
	)
}

// IntersectsCircle reports whether the segment and the circle share a point, as
// Circle.IntersectsSegment does.
func (s Segment[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsSegment(s)
}

// IntersectionCircle returns the points where the segment crosses the circle boundary, as
// Circle.IntersectionSegment does.
func (s Segment[T]) IntersectionCircle(circle Circle[T]) []Point[T] {
	return circle.IntersectionSegment(s)
}

// IntersectsSegment reports whether the segments share a point, within Epsilon of T, the same closed
// convention as Contains: segments that touch at an endpoint or overlap collinearly intersect.
// It holds exactly where DistanceToSegment is zero.
func (s Segment[T]) IntersectsSegment(segment Segment[T]) bool {
	return s.DistanceSquaredToSegment(segment) == 0
}

// IntersectionSegment returns the point where the segments cross, and false when they do not. It
// answers exactly where Intersects holds, less the parallel case: parallel segments have no
// single crossing point and return false even where they overlap, which Intersects still
// reports. A zero-length segment is a point, parallel to nothing, and is the answer wherever it
// lies on the other segment. For integer T the crossing is rounded like every other result
// stored into T.
//
// The segments cross where Intersects says they do: either they properly cross, decided on
// cross products that are exact for an integer T, and the point is where the lines through
// them meet; or an endpoint of one lies on the other within Epsilon of T, and that endpoint is
// the point. The touch is judged on the endpoint's distance, like Contains, rather than on the
// fraction along the segment, which for a shallow crossing can put the same endpoint far
// outside the other segment.
func (s Segment[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	if s.crosses(segment) {
		return s.crossing(segment), true
	}

	if s.parallel(segment) {
		return Point[T]{}, false
	}

	return s.touch(segment)
}

// IntersectsPolygon reports whether the segment and the polygon share a point: the start lies
// within the polygon, or the segment crosses one of its edges. Touching shapes intersect,
// within Epsilon of T. A segment whose extent lies outside the polygon is rejected before any
// edge is examined, and an empty polygon intersects nothing.
func (s Segment[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if polygon.Empty() {
		return false
	}

	a, b := polygon.minMax()
	probe := edgeProbe[T]{a: a, b: b}

	if !probe.aim(s) {
		return false
	}

	if polygon.containsWithin(s.Start, a, b) {
		return true
	}

	for edge := range polygon.Edges() {
		if probe.meets(edge) {
			return true
		}
	}

	return false
}

// IntersectionPolygon returns the points where the segment crosses the polygon boundary, from
// Start to End: the crossings with its edges by IntersectionSegment, with a vertex hit by two
// edges counted once. A segment inside crosses no boundary and returns none while
// IntersectsPolygon still reports it, a segment along an edge is parallel to it and crosses
// only the edges at its ends, and an empty polygon has no boundary to cross.
func (s Segment[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s}
	for edge := range polygon.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectsRectangle reports whether the segment and the rectangle share a point: the start
// lies within the rectangle, or the segment crosses one of its edges. Touching shapes intersect,
// within Epsilon of T. A segment whose extent lies outside the rectangle is rejected before any
// edge is examined.
func (s Segment[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a, b := rectangle.MinMax()
	probe := edgeProbe[T]{a: a, b: b}

	if !probe.aim(s) {
		return false
	}

	if rectangle.containsWithin(s.Start, a, b) {
		return true
	}

	for edge := range rectangle.Edges() {
		if probe.meets(edge) {
			return true
		}
	}

	return false
}

// IntersectionRectangle returns the points where the segment crosses the rectangle boundary,
// from Start to End: the crossings with its edges by Intersection, with a corner hit by two
// edges counted once. A segment inside crosses no boundary and returns none while
// IntersectsRectangle still reports it, and a segment along an edge is parallel to it and
// crosses only the edges at its ends, if it reaches them.
func (s Segment[T]) IntersectionRectangle(rectangle Rectangle[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s}
	for edge := range rectangle.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectsRegularPolygon reports whether the segment and the regular polygon share a point:
// the start lies within the polygon, or the segment crosses one of its edges, the answer
// IntersectsPolygon gives on the polygon's Polygon form, without building it. Touching shapes
// intersect, within Epsilon of T. A segment whose extent lies outside the polygon's Bounds is
// rejected before any edge is examined, and an empty polygon intersects nothing.
func (s Segment[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.Empty() {
		return false
	}

	a, b := polygon.minMax()
	probe := edgeProbe[T]{a: a, b: b}

	if !probe.aim(s) {
		return false
	}

	if polygon.containsWithin(s.Start, a, b) {
		return true
	}

	for edge := range polygon.Edges() {
		if probe.meets(edge) {
			return true
		}
	}

	return false
}

// IntersectionRegularPolygon returns the points where the segment crosses the regular polygon
// boundary, from Start to End, the points IntersectionPolygon returns on the polygon's Polygon
// form, collected over the edges Edges iterates without building the vertices: the crossings
// by IntersectionSegment, with a vertex hit by two edges counted once. A segment inside crosses
// no boundary and returns none while IntersectsRegularPolygon still reports it, and an empty
// polygon has no boundary to cross.
func (s Segment[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s}
	for edge := range polygon.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// distanceSquaredTo returns the squared distance to the point with no tolerance applied, which
// DistanceSquaredTo snaps to zero within Epsilon of T.
func (s Segment[T]) distanceSquaredTo(point Point[T]) float64 {
	direction, offset := s.Vector().Float(), point.Subtract(s.Start).Float()

	along := offset.Dot(direction)
	if along <= 0 {
		return offset.LengthSquared()
	}

	lengthSquared := direction.LengthSquared()
	if along >= lengthSquared {
		return point.Subtract(s.End).Float().LengthSquared()
	}

	cross := offset.Cross(direction)

	return cross * cross / lengthSquared
}

// crosses reports whether the segments properly cross: each has its endpoints on opposite sides
// of the other. Touching and collinear segments do not cross and are left to the endpoint
// distances, which cover them within the tolerance of the caller.
func (s Segment[T]) crosses(segment Segment[T]) bool {
	return s.separates(segment) && segment.separates(s)
}

// separates reports whether the endpoints of the given segment lie strictly on opposite sides
// of the line through this one.
func (s Segment[T]) separates(segment Segment[T]) bool {
	direction := s.Vector().Float()
	start := direction.Cross(segment.Start.Subtract(s.Start).Float())
	end := direction.Cross(segment.End.Subtract(s.Start).Float())

	return (start > 0 && end < 0) || (start < 0 && end > 0)
}

// crossing returns the point where the lines through two properly crossing segments meet, the
// fraction of the way along this segment from the same cross products crosses decided on.
func (s Segment[T]) crossing(segment Segment[T]) Point[T] {
	a, b := s.Float(), segment.Float()

	t := b.Start.Subtract(a.Start).Cross(b.Vector()) / a.Vector().Cross(b.Vector())
	point := a.Lerp(t)

	return point.Cast[T]()
}

// parallel reports whether the segments run along the same direction. A zero-length segment
// has no direction and is parallel to nothing, so its point can still be found on the other.
func (s Segment[T]) parallel(segment Segment[T]) bool {
	a, b := s.Vector(), segment.Vector()

	return a.hasDirection() && b.hasDirection() && a.Float().Cross(b.Float()) == 0
}

// touch returns the endpoint of either segment that lies on the other, within Epsilon of T
// as Contains judges it, and false when there is none. Non-parallel segments that do not
// properly cross can share a point only this way.
func (s Segment[T]) touch(segment Segment[T]) (Point[T], bool) {
	switch {
	case s.Contains(segment.Start):
		return segment.Start, true
	case s.Contains(segment.End):
		return segment.End, true
	case segment.Contains(s.Start):
		return s.Start, true
	case segment.Contains(s.End):
		return s.End, true
	default:
		return Point[T]{}, false
	}
}

// chord returns the fractions along the segment where the line through it enters and leaves
// the circle, in that order and not clamped to the segment, and false where the line misses
// or the segment has no direction. Whether the line reaches the circle is decided on the
// squared distance of the line, the expression DistanceSquaredTo evaluates for a point beside
// the segment, by the same comparison IntersectsCircle makes, so the two agree to the last
// bit. A chord whose ends lie within Epsilon of T of each other is a tangent and both
// fractions are its midpoint: the tolerance collapses two crossings only where they would
// compare Equal, never a chord that merely grazes the boundary within the tolerance.
func (s Segment[T]) chord(circle Circle[T]) (float64, float64, bool) {
	direction, offset := s.Vector().Float(), circle.Center.Subtract(s.Start).Float()
	if !direction.hasDirection() {
		return 0, 0, false
	}

	lengthSquared := direction.LengthSquared()
	cross := offset.Cross(direction)
	gapSquared := cross * cross / lengthSquared

	if !circle.containsSquared(gapSquared) {
		return 0, 0, false
	}

	along := offset.Dot(direction) / lengthSquared
	radius := float64(circle.Radius)
	halfChord := math.Sqrt(max(radius*radius-gapSquared, 0))
	if 2*halfChord <= Epsilon[T]() {
		return along, along, true
	}

	half := halfChord / math.Sqrt(lengthSquared)

	return along - half, along + half, true
}

// snapToEndpoint replaces whichever of the two chord fractions lies nearer the given endpoint,
// 0 for Start and 1 for End, with the endpoint itself: an endpoint on the boundary is the
// crossing nearest to it, not a third crossing beside it. It reads no field of the segment,
// only the fractions along it, so the receiver is unnamed.
func (Segment[T]) snapToEndpoint(entry, exit, endpoint float64) (float64, float64) {
	if math.Abs(entry-endpoint) <= math.Abs(exit-endpoint) {
		return endpoint, exit
	}

	return entry, endpoint
}

// pointsAt returns the points at the given fractions along the segment that lie within it,
// the endpoints included within Epsilon of T scaled to the length, so the tolerance is the same
// distance the Intersects methods apply, with points that compare Equal counted once, as
// crossings counts them. The fractions must be in increasing order.
func (s Segment[T]) pointsAt(fractions ...float64) []Point[T] {
	var points []Point[T]
	for _, t := range fractions {
		if !s.containsAt(t) {
			continue
		}

		lerped := s.Float().Lerp(Clamp(t, 0, 1))
		point := lerped.Cast[T]()
		if slices.ContainsFunc(points, point.Equal) {
			continue
		}

		points = append(points, point)
	}

	return points
}

// containsAt reports whether the fraction t of the way along the segment lies within it, the
// endpoints included within Epsilon of T scaled to the length, so that the tolerance is the
// same distance Contains and Intersects apply.
func (s Segment[T]) containsAt(t float64) bool {
	epsilon := ratio(Epsilon[T](), s.Length())

	return LessOrEqualDelta(0, t, epsilon) && LessOrEqualDelta(t, 1, epsilon)
}

// compareDistance orders two points by their distance from Start, the order every boundary crossing
// method returns its points in.
func (s Segment[T]) compareDistance(a, b Point[T]) int {
	return cmp.Compare(s.Start.DistanceSquaredTo(a), s.Start.DistanceSquaredTo(b))
}

// crossesRay reports whether a ray cast from the point along +X crosses the segment, counting
// an endpoint on the ray only when it is the lower one, so a ray through a vertex is counted
// once by the two edges that share it. It is the step of the even-odd rule Polygon.Contains uses.
//
// The ray crosses when the point lies on the side of the segment facing -X: the left side of a
// segment running toward +Y, the right side of one running toward -Y. The side comes from the
// sign of a cross product, which needs no division by the segment's Y span.
func (s Segment[T]) crossesRay(point Point[T]) bool {
	start, end, p := s.Start.Float(), s.End.Float(), point.Float()

	if (start.Y > p.Y) == (end.Y > p.Y) {
		return false
	}

	upward := end.Y > start.Y
	left := end.Subtract(start).Cross(p.Subtract(start)) > 0

	return left == upward
}

// Equal checks if the start and end points of the segments are equal.
func (s Segment[T]) Equal(segment Segment[T]) bool {
	return s.Start.Equal(segment.Start) && s.End.Equal(segment.End)
}

// IsZero checks if start and end points are zero.
func (s Segment[T]) IsZero() bool {
	return s.Equal(Segment[T]{})
}

// Cast converts the segment to a Segment of another number type, rounding as Cast does.
func (s Segment[T]) Cast[R Number]() Segment[R] {
	return Segment[R]{s.Start.Cast[R](), s.End.Cast[R]()}
}

// Int converts the segment to a Segment[int].
func (s Segment[T]) Int() Segment[int] {
	return Segment[int]{s.Start.Int(), s.End.Int()}
}

// Float converts the segment to a Segment[float64].
func (s Segment[T]) Float() Segment[float64] {
	return Segment[float64]{s.Start.Float(), s.End.Float()}
}

// String returns the segment in the form of its constructor: Seg((x,y);(x,y)).
func (s Segment[T]) String() string {
	return fmt.Sprintf("Seg(%s;%s)", s.Start.String(), s.End.String())
}

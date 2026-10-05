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
	return s.Float().Vector().Length()
}

// Angle returns the angle of the segment in radians, the angle of the vector from Start to End.
// A zero-length segment has no direction and gives 0, the angle Vector.Angle gives it.
func (s Segment[T]) Angle() float64 {
	return s.Float().Vector().Angle()
}

// Direction returns the direction nearest to the segment, from Start to End, or DirectionNone
// for a zero-length segment, as Vector.Direction judges it.
func (s Segment[T]) Direction() Direction {
	return s.Float().Vector().Direction()
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

// Midpoint returns the midpoint of the segment, PointAt(0.5).
func (s Segment[T]) Midpoint() Point[T] {
	return s.Start.Midpoint(s.End)
}

// Bounds returns the axis-aligned bounding box.
func (s Segment[T]) Bounds() Box[T] {
	a, b := s.minMax()

	return Box[T]{a, b}
}

// minMax returns the minimum and maximum corner of the segment, the corners of Bounds: the pair
// the intersection tests reject shapes by before examining any edge.
func (s Segment[T]) minMax() (Point[T], Point[T]) {
	return Point[T]{min(s.Start.X, s.End.X), min(s.Start.Y, s.End.Y)}, Point[T]{max(s.Start.X, s.End.X), max(s.Start.Y, s.End.Y)}
}

// magnitude returns the largest absolute coordinate of the segment, that of an endpoint: the
// size epsilonAt widens the tolerance by for a comparison that reads the segment.
func (s Segment[T]) magnitude() float64 {
	return max(s.Start.magnitude(), s.End.magnitude())
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
	pivot := s.Midpoint().Float()
	start := pivot.Add(s.Start.Float().Subtract(pivot).MultiplyXY(factorX, factorY))
	end := pivot.Add(s.End.Float().Subtract(pivot).MultiplyXY(factorX, factorY))

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Unscale creates a new Segment uniformly scaled about its midpoint by the inverse factor, the
// inverse of Scale. Like Divide it panics for a zero factor.
func (s Segment[T]) Unscale(factor float64) Segment[T] {
	return s.UnscaleXY(factor, factor)
}

// UnscaleXY creates a new Segment scaled about its midpoint by the inverse of the given factors,
// the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (s Segment[T]) UnscaleXY(factorX, factorY float64) Segment[T] {
	pivot := s.Midpoint().Float()
	start := pivot.Add(s.Start.Float().Subtract(pivot).DivideXY(factorX, factorY))
	end := pivot.Add(s.End.Float().Subtract(pivot).DivideXY(factorX, factorY))

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
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
	half := s.Float().Vector().Resize(length / 2)

	start, end := pivot.Add(half.Negate()), pivot.Add(half)

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Reverse creates a new Segment with the start and end points swapped.
func (s Segment[T]) Reverse() Segment[T] {
	return Segment[T]{s.End, s.Start}
}

// PointAt returns the point at the fraction t of the way from Start to End, extrapolating along
// the segment outside [0, 1] like Point.Lerp. Midpoint is PointAt(0.5).
func (s Segment[T]) PointAt(t float64) Point[T] {
	return s.Start.Lerp(s.End, t)
}

// Transform creates a new Segment by applying the given matrix to both points, like
// Point.Transform.
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

// Contains reports whether the given point lies on the segment, within the tolerance: it holds
// exactly where DistanceTo is zero.
func (s Segment[T]) Contains(point Point[T]) bool {
	return s.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the segment:
// zero exactly where Contains holds, so a point within the tolerance of the segment is at
// distance zero rather than at the rounding error that put it there. A lattice point on a
// lattice segment is at zero without any tolerance, since the perpendicular distance comes
// from a cross product that is exact in float64 rather than from a projection.
func (s Segment[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(s.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. It is a float64 even for an integer T, unlike Point.DistanceSquaredTo, since the
// nearest point of a segment is not a lattice point in general. A distance within the tolerance
// is snapped to zero, which is where Contains reads it, so a polygon boundary is walked on one
// comparison per edge.
func (s Segment[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := s.distanceSquaredTo(point)
	if lessOrEqualSquared(distance, 0, epsilonAt[T](max(s.magnitude(), point.magnitude()))) {
		return 0
	}

	return distance
}

// Nearest returns the point of the segment nearest to the given point: the point itself
// exactly where Contains holds, and otherwise the foot of the perpendicular, or the endpoint
// where the foot falls beyond it, decided on the projection DistanceSquaredTo makes. For
// integer T the foot is rounded once and can land off the segment, where Contains rejects it.
func (s Segment[T]) Nearest(point Point[T]) Point[T] {
	if s.DistanceSquaredTo(point) == 0 {
		return point
	}

	return s.foot(point)
}

// DistanceToSegment returns the distance between the nearest points of the two segments: zero
// exactly where Intersects holds, and otherwise the smallest distance from an endpoint of one
// to the other.
func (s Segment[T]) DistanceToSegment(segment Segment[T]) float64 {
	return math.Sqrt(s.DistanceSquaredToSegment(segment))
}

// DistanceSquaredToSegment returns the squared distance DistanceToSegment takes the root of, faster
// for comparisons: zero where the segments properly cross or an endpoint of one lies on the
// other within the tolerance, as DistanceSquaredTo snaps it.
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

// AppendIntersectionCircle appends the points IntersectionCircle returns to dst and returns the
// extended slice, as Circle.AppendIntersectionSegment does.
func (s Segment[T]) AppendIntersectionCircle(dst []Point[T], circle Circle[T]) []Point[T] {
	return circle.AppendIntersectionSegment(dst, s)
}

// IntersectsSegment reports whether the segments share a point, within the tolerance, the same
// closed convention as Contains: segments that touch at an endpoint or overlap collinearly
// intersect. It holds exactly where DistanceToSegment is zero.
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
// them meet, within both; or an endpoint of one lies on the other within the tolerance, and that
// endpoint is the point, which is also the answer for nearly collinear float segments whose
// crossing rounding would place off them. The touch is judged on the endpoint's distance, like
// Contains, rather than on the fraction along the segment, which for a shallow crossing can put
// the same endpoint far outside the other segment. An endpoint within the tolerance replaces a
// proper crossing too, so a segment passing a vertex within the tolerance meets both edges there
// at that one vertex, rather than at the vertex on one and a crossing beside it on the other.
// That covers segments running along each other within the tolerance, sharing a stretch rather
// than a point: an endpoint of one lies on the other, so it is the answer even where rounding
// gives them a proper crossing, which would land anywhere along the stretch. Such segments
// share more than one point, and the endpoint returned is one of the given segment before one
// of this one, so the two orders can name different ends of the stretch.
func (s Segment[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	if point, ok := s.crossing(segment); ok {
		if endpoint, ok := s.touch(segment); ok {
			return endpoint, true
		}

		return point.Cast[T](), true
	}

	if s.parallel(segment) {
		return Point[T]{}, false
	}

	return s.touch(segment)
}

// IntersectsRay reports whether the segment and the ray share a point, as IntersectsSegment
// decides it on the reach of the ray past the segment. Touching shapes intersect, within
// the tolerance.
func (s Segment[T]) IntersectsRay(ray Ray[T]) bool {
	return s.IntersectsSegment(ray.reach(s.minMax()))
}

// IntersectionRay returns the point where the segment and the ray cross, and false when they do
// not, as IntersectionSegment finds it on the reach of the ray past the segment: exactly where
// IntersectsRay holds, less the parallel case.
func (s Segment[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	return s.IntersectionSegment(ray.reach(s.minMax()))
}

// IntersectsPolygon reports whether the segment and the polygon share a point: the start lies
// within the polygon, or the segment crosses one of its edges. Touching shapes intersect,
// within the tolerance. A segment whose extent lies outside the polygon is rejected before any
// edge is examined, and an empty polygon intersects nothing.
func (s Segment[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
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
	return s.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points IntersectionPolygon returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room. The points already
// in dst are kept as they are: a crossing equal to one of them is still appended, and only the
// appended ones are ordered from Start.
func (s Segment[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	for edge := range polygon.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectsRectangle reports whether the segment and the rectangle share a point: the start
// lies within the rectangle, or the segment crosses one of its edges. Touching shapes intersect,
// within the tolerance. A segment whose extent lies outside the rectangle is rejected before any
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
// from Start to End: the crossings with its edges by IntersectionSegment, with a corner hit by two
// edges counted once. A segment inside crosses no boundary and returns none while
// IntersectsRectangle still reports it, and a segment along an edge is parallel to it and
// crosses only the edges at its ends, if it reaches them.
func (s Segment[T]) IntersectionRectangle(rectangle Rectangle[T]) []Point[T] {
	return s.AppendIntersectionRectangle(nil, rectangle)
}

// AppendIntersectionRectangle appends the points IntersectionRectangle returns to dst and returns
// the extended slice, so a caller reusing dst allocates nothing once it has room. The points
// already in dst are kept as they are: a crossing equal to one of them is still appended, and only
// the appended ones are ordered from Start.
func (s Segment[T]) AppendIntersectionRectangle(dst []Point[T], rectangle Rectangle[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	for edge := range rectangle.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectsRegularPolygon reports whether the segment and the regular polygon share a point:
// the start lies within the polygon, or the segment crosses one of its edges, the answer
// IntersectsPolygon gives on the polygon's Polygon form, without building it. Touching shapes
// intersect, within the tolerance. A segment whose extent lies outside the polygon's Bounds is
// rejected before any edge is examined, and an empty polygon intersects nothing.
func (s Segment[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
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
	return s.AppendIntersectionRegularPolygon(nil, polygon)
}

// AppendIntersectionRegularPolygon appends the points IntersectionRegularPolygon returns to dst
// and returns the extended slice, so a caller reusing dst allocates nothing once it has room. The
// points already in dst are kept as they are: a crossing equal to one of them is still appended,
// and only the appended ones are ordered from Start.
func (s Segment[T]) AppendIntersectionRegularPolygon(dst []Point[T], polygon RegularPolygon[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	for edge := range polygon.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectsBox reports whether the segment and the box share a point, as IntersectsRectangle
// decides it on a rectangle that is not rotated: the start lies within the box, or the segment
// crosses one of its edges, the edges between its own corners. Touching shapes intersect,
// within the tolerance. A segment whose extent lies outside the box is rejected before any edge
// is examined.
func (s Segment[T]) IntersectsBox(box Box[T]) bool {
	probe := edgeProbe[T]{a: box.Min, b: box.Max}

	if !probe.aim(s) {
		return false
	}

	if box.Contains(s.Start) {
		return true
	}

	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if probe.meets(edge) {
			return true
		}
	}

	return false
}

// IntersectionBox returns the points where the segment crosses the box boundary, from Start to
// End, as IntersectionRectangle finds them on a rectangle that is not rotated, over the edges
// between the corners of the box.
func (s Segment[T]) IntersectionBox(box Box[T]) []Point[T] {
	return s.AppendIntersectionBox(nil, box)
}

// AppendIntersectionBox appends the points IntersectionBox returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room. The points already
// in dst are kept as they are, as AppendIntersectionRectangle keeps them.
func (s Segment[T]) AppendIntersectionBox(dst []Point[T], box Box[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		e.add(edge)
	}

	return e.sorted()
}

// ClipCircle returns the part of the segment inside the circle, boundary included within
// the tolerance, and false where they share no point: from Start where the circle contains it,
// or else from the first point IntersectionCircle returns, to End where the circle contains it,
// or else to the last. A segment tangent to the circle, or touching it with an endpoint from
// outside, is clipped to that one point, so the part exists exactly where IntersectsCircle
// holds. It allocates nothing.
func (s Segment[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	var buffer [2]Point[T]
	crossings := circle.AppendIntersectionSegment(buffer[:0], s)

	return s.clipConvex(crossings, circle.Contains(s.Start), circle.Contains(s.End))
}

// ClipPolygon returns the parts of the segment inside the polygon, boundary included within the
// tolerance, from Start to End. The points IntersectionPolygon returns cut the segment into pieces,
// each wholly inside or outside, and each piece is judged at its midpoint by the walk Contains
// makes, Start and End by Contains itself. A Start that a crossing compares Equal to is that
// crossing, as an End is, so an entry an integer T rounds onto a Start outside the polygon still
// opens the part there. A piece is inside only where both its ends are too: a segment running along
// an edge at a shallow angle leaves the tolerance without crossing it, so a piece from a Start or
// to an End the polygon does not contain is outside whatever its midpoint, and the part ends at a
// point the polygon holds. Pieces inside run together across a point where the segment touches the
// boundary from inside, such as a reflex vertex, and a point where it touches the boundary from
// outside is a part of zero length, so there are parts exactly where IntersectsPolygon holds, less
// a flat polygon, its vertices on one line: a segment along that line with both ends beyond the
// polygon is parallel to every edge, crosses none and has no part, while IntersectsPolygon still
// reports it. The midpoint is taken in float64, so for an integer T a gap between two crossings is
// judged where it is rather than at a rounded point on either side; the crossings themselves are
// rounded as IntersectionPolygon rounds them. The polygon follows the even-odd rule of Contains,
// and an empty polygon clips everything away.
//
// The crossings are held in an array, the next few from the point reached, and the edges are
// swept again only where a segment crosses more often than it holds, so the result is the one
// allocation, made on the first part with room for the one a convex polygon gives; only a
// concave polygon grows it.
func (s Segment[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return s.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts ClipPolygon returns to dst and returns the extended slice,
// so a caller reusing dst allocates nothing once it has room. The sweep ends at End once it has
// stepped there, even where End compares Equal to nothing, as a NaN or an infinite coordinate
// does, so every segment is clipped in a bounded number of sweeps.
func (s Segment[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	var from Point[T]

	point, open, ended := s.Start, false, false
	sweep := s.sweep(polygon, point)
	contained := polygon.Contains(point) || sweep.touchedAt(point)
	for {
		next, found := sweep.next(point)
		if !found && sweep.spent() {
			sweep = s.sweep(polygon, point)
			next, found = sweep.next(point)
		}

		ahead := found && !ended
		nextContained := ahead
		if !ahead && !ended && !point.Equal(s.End) {
			next, ahead, ended = s.End, true, true
			nextContained = polygon.Contains(s.End)
		}

		inside := ahead && contained && nextContained && polygon.containsMidpoint(point, next)
		if !open && contained {
			from, open = point, true
		}

		if open && !inside {
			dst, open = append(dst, Segment[T]{from, point}), false
		}

		if !ahead {
			return dst
		}

		point, contained = next, nextContained
	}
}

// ClipRectangle returns the part of the segment inside the rectangle, boundary included within
// the tolerance, whatever its angle, and false where they share no point: from Start where the
// rectangle contains it, or else from the first point IntersectionRectangle returns, to End
// where the rectangle contains it, or else to the last. A segment along an edge is clipped to
// the part of the edge it covers, and one touching a corner from outside to that corner, so
// the part exists exactly where IntersectsRectangle holds. It allocates nothing.
func (s Segment[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	for edge := range rectangle.Edges() {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), rectangle.Contains(s.Start), rectangle.Contains(s.End))
}

// ClipRegularPolygon returns the part of the segment inside the regular polygon, boundary
// included within the tolerance, and false where they share no point: from Start where the
// polygon contains it, or else from the first point IntersectionRegularPolygon returns, to End
// where the polygon contains it, or else to the last, the part ClipPolygon returns on the
// polygon's Polygon form, without building it. A segment touching a vertex from outside is
// clipped to that vertex, so the part exists exactly where IntersectsRegularPolygon holds, and
// an empty polygon clips everything away. A flat polygon, of two vertices or a zero semi-axis,
// is the exception: a segment along its line with an end beyond it is parallel to every edge,
// crosses none and has no part, while IntersectsRegularPolygon still reports it. It allocates
// nothing.
func (s Segment[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	for edge := range polygon.Edges() {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), polygon.Contains(s.Start), polygon.Contains(s.End))
}

// ClipBox returns the part of the segment inside the box, as ClipRectangle clips it to a
// rectangle that is not rotated, over the edges between the corners of the box: the segment as
// a viewport shows it. It allocates nothing.
func (s Segment[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), box.Contains(s.Start), box.Contains(s.End))
}

// distanceSquaredTo returns the squared distance to the point with no tolerance applied, which
// DistanceSquaredTo snaps to zero within the tolerance. Beside the segment it is the square of
// the gap, the cross product divided by the length, rather than the squared cross product
// divided by the squared length: for an integer T the cross product is exact and so is the
// length wherever it is a whole number, the only lengths a lattice circle can be tangent on, so
// a tangent is at exactly the radius across the supported range, where the squared cross
// product rounds once the radius times the length passes 2^26.
func (s Segment[T]) distanceSquaredTo(point Point[T]) float64 {
	start, end, p := s.Start.Float(), s.End.Float(), point.Float()
	direction, offset := end.Subtract(start), p.Subtract(start)

	along := offset.Dot(direction)
	if along <= 0 {
		return offset.LengthSquared()
	}

	lengthSquared := direction.LengthSquared()
	if along >= lengthSquared {
		return p.Subtract(end).LengthSquared()
	}

	gap := offset.Cross(direction) / math.Sqrt(lengthSquared)

	return gap * gap
}

// foot returns the point of the segment nearest to the given point with no tolerance applied:
// Start or End where the projection distanceSquaredTo makes falls beyond them, on the same
// comparisons, and the point at that fraction along the segment otherwise.
func (s Segment[T]) foot(point Point[T]) Point[T] {
	start, end, p := s.Start.Float(), s.End.Float(), point.Float()
	direction, offset := end.Subtract(start), p.Subtract(start)

	along := offset.Dot(direction)
	if along <= 0 {
		return s.Start
	}

	lengthSquared := direction.LengthSquared()
	if along >= lengthSquared {
		return s.End
	}

	return s.PointAt(along / lengthSquared)
}

// crosses reports whether the segments properly cross, as crossing decides it: each has its
// endpoints on opposite sides of the other, and the lines through them meet within both.
// Touching and collinear segments do not cross and are left to the endpoint distances, which
// cover them within the tolerance of the caller.
func (s Segment[T]) crosses(segment Segment[T]) bool {
	_, ok := s.crossing(segment)

	return ok
}

// crossing returns the point where the lines through two properly crossing segments meet, in
// float64, and false where they do not properly cross: the fraction along each segment lies
// strictly between its ends, decided on the signs and the order of the cross products it is the
// quotient of, which are exact for an integer T, with no division. For nearly collinear float
// segments rounding can pass both fractions, and the point, divided out of two cross products
// near zero, lands anywhere on the line; a point whose projection onto the other segment falls
// outside it is therefore no crossing, and the pair is left to the endpoints, on the
// comparisons foot makes.
//
// The point is placed from Start by each component of the direction times the numerator of the
// fraction, divided by its denominator once, rather than by the rounded fraction: for an integer
// T the product and the cross products are exact, so a crossing on a half unit is exactly there
// and rounds the same way whatever lengths the two segments have, while that product stays
// within the integers float64 holds exactly. Past it the product rounds, so the pair is taken
// in one order whichever segment asks and whichever way each runs: each from its lesser
// endpoint by Point.Compare, then the lesser segment by Start and then End first, and a
// crossing on a half unit still rounds the same way from either side and for a Reverse.
func (s Segment[T]) crossing(segment Segment[T]) (Point[float64], bool) {
	if s.End.Compare(s.Start) < 0 {
		s = s.Reverse()
	}

	if segment.End.Compare(segment.Start) < 0 {
		segment = segment.Reverse()
	}

	if cmp.Or(segment.Start.Compare(s.Start), segment.End.Compare(s.End)) < 0 {
		s, segment = segment, s
	}

	a, b := s.Float(), segment.Float()
	along, direction, offset := a.Vector(), b.Vector(), b.Start.Subtract(a.Start)
	denominator, numerator, other := along.Cross(direction), offset.Cross(direction), offset.Cross(along)
	if denominator < 0 {
		denominator, numerator, other = -denominator, -numerator, -other
	}

	if !(0 < numerator && numerator < denominator && 0 < other && other < denominator) {
		return Point[float64]{}, false
	}

	point := Point[float64]{a.Start.X + float64(numerator*along.X)/denominator, a.Start.Y + float64(numerator*along.Y)/denominator}
	projection := point.Subtract(b.Start).Dot(direction)

	return point, 0 <= projection && projection <= direction.LengthSquared()
}

// parallel reports whether the segments run along the same direction. A zero-length segment
// has no direction and is parallel to nothing, so its point can still be found on the other.
func (s Segment[T]) parallel(segment Segment[T]) bool {
	a, b := s.Float().Vector(), segment.Float().Vector()

	return a.hasDirection() && b.hasDirection() && a.Cross(b) == 0
}

// touch returns the endpoint of either segment that lies on the other, within the tolerance
// as Contains judges it, and false when there is none. Non-parallel segments that do not
// properly cross can share a point only this way. An endpoint of the given segment comes
// first, so a segment ending at a vertex within the tolerance meets both edges there at the
// vertex itself, the one point the two edges share.
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
// squared distance of the line, the squared gap DistanceSquaredTo evaluates for a point beside
// the segment, by the same comparison IntersectsCircle makes, so the two agree to the last
// bit. A chord whose ends lie within the tolerance of each other is a tangent and both
// fractions are its midpoint, judged on the squared chord by the same comparison, before any
// root is taken: the tolerance collapses two crossings only where they would compare Equal,
// never a chord that merely grazes the boundary within the tolerance.
func (s Segment[T]) chord(circle Circle[T]) (float64, float64, bool) {
	start, end, center := s.Start.Float(), s.End.Float(), circle.Center.Float()
	direction, offset := end.Subtract(start), center.Subtract(start)
	if !direction.hasDirection() {
		return 0, 0, false
	}

	lengthSquared := direction.LengthSquared()
	gap := offset.Cross(direction) / math.Sqrt(lengthSquared)
	gapSquared := float64(gap * gap)

	if !circle.containsSquared(gapSquared, s.magnitude()) {
		return 0, 0, false
	}

	along := offset.Dot(direction) / lengthSquared
	radius := float64(circle.Radius)
	halfChordSquared := max(float64(radius*radius)-gapSquared, 0)
	if lessOrEqualSquared(4*halfChordSquared, 0, circle.epsilonWith(s.magnitude())) {
		return along, along, true
	}

	half := math.Sqrt(halfChordSquared / lengthSquared)

	return along - half, along + half, true
}

// snapToEndpoint replaces one of the two chord fractions with the given endpoint, 0 for Start
// and 1 for End: an endpoint on the boundary is one of the two crossings, not a third beside
// them. Start is the entry and End the exit, the crossing each is of a segment running through
// the chord, so an endpoint within the chord keeps the crossing ahead of it and gives the one
// behind it up, whichever way the segment runs; the nearer fraction would be decided by rounding
// for an endpoint in the middle of the chord, and the reversed segment would lose the crossing
// this one keeps. Only where the chord lies wholly past the endpoint, outside the segment or
// ending at the endpoint within the tolerance scaled to the length as containsAt scales it, is
// the endpoint the end of the chord it meets, the exit for Start and the entry for End.
func (s Segment[T]) snapToEndpoint(entry, exit, endpoint float64) (float64, float64) {
	epsilon := ratio(epsilonAt[T](s.magnitude()), s.Length())

	if endpoint == 0 {
		if entry < 0 && exit <= epsilon {
			return entry, endpoint
		}

		return endpoint, exit
	}

	if exit > 1 && entry >= 1-epsilon {
		return endpoint, exit
	}

	return entry, endpoint
}

// appendPointsAt appends to dst the points at the given fractions along the segment that lie
// within it, the endpoints included within the tolerance scaled to the length, so the tolerance
// is the same distance the Intersects methods apply, with points that compare Equal counted
// once, as edgeIntersections counts them. The fractions must be in increasing order. Only the
// appended points are deduplicated, and a nil dst is allocated once on the first point with
// room for the two a segment can have.
func (s Segment[T]) appendPointsAt(dst []Point[T], fractions ...float64) []Point[T] {
	n := len(dst)
	for _, t := range fractions {
		if !s.containsAt(t) {
			continue
		}

		lerped := s.Float().PointAt(Clamp(t, 0, 1))
		point := lerped.Cast[T]()
		if slices.ContainsFunc(dst[n:], point.Equal) {
			continue
		}

		if dst == nil {
			dst = make([]Point[T], 0, 2)
		}

		dst = append(dst, point)
	}

	return dst
}

// containsAt reports whether the fraction t of the way along the segment lies within it, the
// endpoints included within the tolerance scaled to the length, so that the tolerance is the
// same distance Contains and Intersects apply.
func (s Segment[T]) containsAt(t float64) bool {
	epsilon := ratio(epsilonAt[T](s.magnitude()), s.Length())

	return LessOrEqualDelta(0, t, epsilon) && LessOrEqualDelta(t, 1, epsilon)
}

// compareDistance orders two points by their distance from Start, the order every boundary crossing
// method returns its points in. The distances are measured in float64, where the squared
// distance of an integer T in T would overflow past the coordinate differences doc.go names.
func (s Segment[T]) compareDistance(a, b Point[T]) int {
	start := s.Start.Float()

	return cmp.Compare(start.DistanceSquaredTo(a.Float()), start.DistanceSquaredTo(b.Float()))
}

// clipConvex returns the part of the segment inside a convex shape, given the crossings of its
// boundary in order from Start and whether the shape contains Start and End: the whole segment
// where it contains both, since a convex shape holds every point between two of its own, and
// otherwise the span from the first of its points on the segment to the last, a contained
// endpoint standing in for the crossing beside it. It is false where the shape contains
// neither endpoint and the segment crosses no boundary.
func (s Segment[T]) clipConvex(crossings []Point[T], start, end bool) (Segment[T], bool) {
	if start && end {
		return s, true
	}

	if len(crossings) == 0 {
		return Segment[T]{}, false
	}

	clipped := Segment[T]{crossings[0], crossings[len(crossings)-1]}
	if start {
		clipped.Start = s.Start
	}
	if end {
		clipped.End = s.End
	}

	return clipped, true
}

// sweep returns the next points where the segment crosses the polygon boundary beyond the one
// reached, as edgeSweep gathers them over the polygon's Edges.
func (s Segment[T]) sweep(polygon Polygon[T], reached Point[T]) edgeSweep[T] {
	sweep := edgeSweep[T]{segment: s, reached: reached}
	for edge := range polygon.Edges() {
		sweep.add(edge)
	}

	return sweep
}

// crossesRay reports whether a ray cast from the point along +X crosses the segment, counting
// an endpoint on the ray only when it is the lower one, so a ray through a vertex is counted
// once by the two edges that share it. It is the step of the even-odd rule Polygon.Contains uses.
//
// The ray crosses when the point lies on the side of the segment facing -X: the left side of a
// segment running toward +Y, the right side of one running toward -Y. The side comes from the
// sign of a cross product, which needs no division by the segment's Y span, and each side is
// asked for by its own sign, so a NaN coordinate, which has neither, crosses nothing.
func (s Segment[T]) crossesRay(point Point[T]) bool {
	start, end, p := s.Start.Float(), s.End.Float(), point.Float()

	if (start.Y > p.Y) == (end.Y > p.Y) {
		return false
	}

	cross := end.Subtract(start).Cross(p.Subtract(start))
	if end.Y > start.Y {
		return cross > 0
	}

	return cross < 0
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

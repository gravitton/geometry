package geom

import (
	"cmp"
	"iter"
	"math"
)

// Segment is a line segment from Start to End. The order is its direction, so a segment and its
// Reverse are different values over the same points, and a segment of zero length is the point both
// ends are.
//
// It is a Shape with no inside, and the edge every outline is walked by. It moves, turns and scales
// about its Midpoint.
type Segment[T Number] struct {
	Start Point[T] `json:"s"`
	End   Point[T] `json:"e"`
}

// Seg is shorthand for Segment{start, end}.
func Seg[T Number](start, end Point[T]) Segment[T] {
	return Segment[T]{start, end}
}

// Vector returns the vector from Start to End.
func (s Segment[T]) Vector() Vector[T] {
	return s.End.Subtract(s.Start)
}

// Length returns the length.
func (s Segment[T]) Length() float64 {
	return s.Float().Vector().Length()
}

// Angle returns the angle of Vector in radians, zero for a segment of zero length.
func (s Segment[T]) Angle() float64 {
	return s.Float().Vector().Angle()
}

// Direction returns the direction nearest to Vector, or DirectionNone for a segment of zero length.
func (s Segment[T]) Direction() Direction {
	return s.Float().Vector().Direction()
}

// Vertices returns Start and End, in that order.
func (s Segment[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		_ = yield(s.Start) && yield(s.End)
	}
}

// Edges returns the one edge, the segment itself: the outline is open, so no edge returns to Start.
func (s Segment[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		yield(s)
	}
}

// Midpoint returns PointAt(0.5).
func (s Segment[T]) Midpoint() Point[T] {
	return s.Start.Midpoint(s.End)
}

// Bounds returns the axis-aligned bounding box.
func (s Segment[T]) Bounds() Box[T] {
	a, b := s.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle about Midpoint through the endpoints. For integer T the radius
// is rounded up, so the circle still contains both.
func (s Segment[T]) BoundingCircle() Circle[T] {
	center := s.Midpoint()
	midpoint := center.Float()
	reach := max(midpoint.DistanceSquaredTo(s.Start.Float()), midpoint.DistanceSquaredTo(s.End.Float()))

	return Circle[T]{center, castUp[T](math.Sqrt(reach))}
}

// minMax returns the corners of Bounds, the prefilter of the intersection tests.
func (s Segment[T]) minMax() (Point[T], Point[T]) {
	return Point[T]{min(s.Start.X, s.End.X), min(s.Start.Y, s.End.Y)}, Point[T]{max(s.Start.X, s.End.X), max(s.Start.Y, s.End.Y)}
}

// magnitude returns the largest absolute coordinate of an endpoint, the magnitude the tolerance of
// a comparison with the segment is taken from.
func (s Segment[T]) magnitude() float64 {
	return max(s.Start.magnitude(), s.End.magnitude())
}

// cross returns Start × End, the shoelace term of the edge. It is exact for parallel vectors, so a
// degenerate edge adds exactly zero.
func (s Segment[T]) cross() float64 {
	return s.Start.Vector().Float().Cross(s.End.Vector().Float())
}

// Translate returns the segment moved by the vector.
func (s Segment[T]) Translate(vector Vector[T]) Segment[T] {
	return Segment[T]{s.Start.Add(vector), s.End.Add(vector)}
}

// MoveTo returns the segment moved so that its midpoint lies at the point. For integer T an odd
// span is placed within half a unit of it, so Midpoint of the result can be a unit off.
func (s Segment[T]) MoveTo(point Point[T]) Segment[T] {
	return s.Translate(point.Subtract(s.Midpoint()))
}

// Rotate returns the segment turned by the angle in radians about its Midpoint.
func (s Segment[T]) Rotate(angle float64) Segment[T] {
	pivot := s.Midpoint()

	return Segment[T]{s.Start.RotateAround(pivot, angle), s.End.RotateAround(pivot, angle)}
}

// Scale returns the segment scaled by the factor about its Midpoint. A zero factor collapses it
// onto the midpoint.
func (s Segment[T]) Scale(factor float64) Segment[T] {
	return s.ScaleXY(factor, factor)
}

// ScaleXY returns the segment scaled by the factors along X and Y about its Midpoint, which can
// turn it where the factors differ.
func (s Segment[T]) ScaleXY(factorX, factorY float64) Segment[T] {
	pivot := s.Midpoint().Float()
	start := pivot.Add(s.Start.Float().Subtract(pivot).MultiplyXY(factorX, factorY))
	end := pivot.Add(s.End.Float().Subtract(pivot).MultiplyXY(factorX, factorY))

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Unscale returns the segment scaled by the inverse of the factor about its Midpoint. It panics for
// a zero factor.
func (s Segment[T]) Unscale(factor float64) Segment[T] {
	return s.UnscaleXY(factor, factor)
}

// UnscaleXY returns the segment scaled by the inverses of the factors about its Midpoint. It panics
// for a zero factor.
func (s Segment[T]) UnscaleXY(factorX, factorY float64) Segment[T] {
	pivot := s.Midpoint().Float()
	start := pivot.Add(s.Start.Float().Subtract(pivot).DivideXY(factorX, factorY))
	end := pivot.Add(s.End.Float().Subtract(pivot).DivideXY(factorX, factorY))

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Resize returns the segment of the length about its Midpoint, in the same direction. A segment of
// zero length resizes along X, and a negative length gives the Reverse of the segment of its
// absolute length.
func (s Segment[T]) Resize(length float64) Segment[T] {
	pivot := s.Midpoint().Float()
	half := s.Float().Vector().Resize(length / 2)

	start, end := pivot.Add(half.Negate()), pivot.Add(half)

	return Segment[T]{start.Cast[T](), end.Cast[T]()}
}

// Transform returns the segment with the matrix applied to both ends.
func (s Segment[T]) Transform[M Float](matrix Matrix[M]) Segment[T] {
	return Segment[T]{s.Start.Transform(matrix), s.End.Transform(matrix)}
}

// Lerp returns the segment interpolated toward the given one, both ends together, and extrapolated
// outside [0, 1].
func (s Segment[T]) Lerp(segment Segment[T], t float64) Segment[T] {
	return Segment[T]{s.Start.Lerp(segment.Start, t), s.End.Lerp(segment.End, t)}
}

// Reverse returns the segment with Start and End swapped.
func (s Segment[T]) Reverse() Segment[T] {
	return Segment[T]{s.End, s.Start}
}

// PointAt returns the point at the fraction t of the way from Start to End, extrapolated along the
// segment outside [0, 1].
func (s Segment[T]) PointAt(t float64) Point[T] {
	return s.Start.Lerp(s.End, t)
}

// Split returns the segment cut at PointAt(t), the first part from Start and the second to End.
// Outside [0, 1] one part runs past an end and the other back over it.
func (s Segment[T]) Split(t float64) (Segment[T], Segment[T]) {
	point := s.PointAt(t)

	return Segment[T]{s.Start, point}, Segment[T]{point, s.End}
}

// Normal returns Vector turned a quarter turn in the sense of Vector.Rotate, of the same length. On
// an edge of a Rectangle, or of a polygon wound the same way, it points inward. A segment of zero
// length gives the zero vector.
func (s Segment[T]) Normal() Vector[T] {
	return s.Vector().Normal()
}

// DistanceTo returns the distance from the point to the segment, zero on it.
func (s Segment[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(s.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (s Segment[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := s.distanceSquaredTo(point)
	if lessOrEqualSquared(distance, 0, epsilonAt[T](max(s.magnitude(), point.magnitude()))) {
		return 0
	}

	return distance
}

// DistanceToSegment returns the distance between the nearest points of the segments, zero where
// they share a point.
func (s Segment[T]) DistanceToSegment(segment Segment[T]) float64 {
	return math.Sqrt(s.DistanceSquaredToSegment(segment))
}

// DistanceSquaredToSegment returns the square of DistanceToSegment.
func (s Segment[T]) DistanceSquaredToSegment(segment Segment[T]) float64 {
	if s.crosses(segment) {
		return 0
	}

	return min(
		s.DistanceSquaredTo(segment.Start), s.DistanceSquaredTo(segment.End),
		segment.DistanceSquaredTo(s.Start), segment.DistanceSquaredTo(s.End),
	)
}

// Nearest returns the point of the segment nearest to the given point, the point itself on it.
func (s Segment[T]) Nearest(point Point[T]) Point[T] {
	if s.DistanceSquaredTo(point) == 0 {
		return point
	}

	return s.foot(point)
}

// NearestFraction returns the fraction at which PointAt gives the point of the segment nearest to
// the given point, held to [0, 1]. A segment of zero length gives zero.
func (s Segment[T]) NearestFraction(point Point[T]) float64 {
	start, end, p := s.Start.Float(), s.End.Float(), point.Float()
	direction, offset := end.Subtract(start), p.Subtract(start)

	along := offset.Dot(direction)
	if along <= 0 {
		return 0
	}

	lengthSquared := direction.LengthSquared()
	if along >= lengthSquared {
		return 1
	}

	return along / lengthSquared
}

// Support returns the point of the segment farthest along the direction: End where the segment runs
// along it, and Start where it runs against or across it.
func (s Segment[T]) Support(direction Vector[float64]) Point[float64] {
	start, end := s.Start.Float(), s.End.Float()
	if end.Subtract(start).Dot(direction) > 0 {
		return end
	}

	return start
}

// Contains reports whether the point lies on the segment.
func (s Segment[T]) Contains(point Point[T]) bool {
	return s.DistanceSquaredTo(point) == 0
}

// IntersectsCircle reports whether the segment and the circle share a point.
func (s Segment[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsSegment(s)
}

// IntersectsCapsule reports whether the segment and the capsule share a point.
func (s Segment[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsSegment(s)
}

// IntersectsSector reports whether the segment and the sector share a point.
func (s Segment[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsSegment(s)
}

// IntersectsSegment reports whether the segments share a point, collinear overlap included.
func (s Segment[T]) IntersectsSegment(segment Segment[T]) bool {
	return s.DistanceSquaredToSegment(segment) == 0
}

// IntersectsRay reports whether the segment and the ray share a point.
func (s Segment[T]) IntersectsRay(ray Ray[T]) bool {
	return s.IntersectsSegment(ray.reach(s.minMax()))
}

// IntersectsLine reports whether the segment and the line share a point.
func (s Segment[T]) IntersectsLine(line Line[T]) bool {
	return s.IntersectsSegment(line.span(s.minMax()))
}

// IntersectsPolygon reports whether the segment and the polygon share a point.
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

// IntersectsRectangle reports whether the segment and the rectangle share a point.
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

// IntersectsRegularPolygon reports whether the segment and the regular polygon share a point.
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

// IntersectsBox reports whether the segment and the box share a point.
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

// IntersectionCircle returns the points where the segment crosses the boundary of the circle, from
// Start to End.
func (s Segment[T]) IntersectionCircle(circle Circle[T]) Pair[Point[T]] {
	return circle.IntersectionSegment(s)
}

// IntersectionCapsule returns the points where the segment crosses the boundary of the capsule,
// from Start to End.
func (s Segment[T]) IntersectionCapsule(capsule Capsule[T]) Pair[Point[T]] {
	return capsule.IntersectionSegment(s)
}

// IntersectionSector returns the points where the segment crosses the boundary of the sector, from
// Start to End: at most two for a sector of up to half a turn, and up to four past it for a segment
// through the notch between its radii. A segment entirely inside gives none while IntersectsSector
// holds. An endpoint on the boundary is one of the points, so a segment along a radius gives the
// ends of the stretch they share.
func (s Segment[T]) IntersectionSector(sector Sector[T]) []Point[T] {
	return s.AppendIntersectionSector(nil, sector)
}

// AppendIntersectionSector appends the points of IntersectionSector(sector) to dst and returns the
// extended slice.
func (s Segment[T]) AppendIntersectionSector(dst []Point[T], sector Sector[T]) []Point[T] {
	return sector.appendIntersectionSegment(dst, s)
}

// IntersectionArc returns the points where the segment crosses the arc, from Start to End.
func (s Segment[T]) IntersectionArc(arc Arc[T]) Pair[Point[T]] {
	return arc.IntersectionSegment(s)
}

// IntersectionSegment returns the point where the segments cross, and false where they do not.
// Parallel segments give false even where they overlap, which IntersectsSegment reports. A segment
// of zero length is parallel to nothing and is the answer wherever it lies on the other.
//
// An endpoint of one segment on the other is the answer, even where rounding finds a proper
// crossing beside it, so a segment passing a vertex meets both edges at that vertex. Segments
// running along each other share a stretch, and the endpoint returned is one of the given segment
// before one of this one, so the two orders can name different ends of it.
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

// IntersectionRay returns the point where the segment and the ray cross, and false where they do
// not or are parallel.
func (s Segment[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	return s.IntersectionSegment(ray.reach(s.minMax()))
}

// IntersectionLine returns the point where the segment and the line cross, and false where they do
// not or are parallel.
func (s Segment[T]) IntersectionLine(line Line[T]) (Point[T], bool) {
	return s.IntersectionSegment(line.span(s.minMax()))
}

// IntersectionPolyline returns the points where the segment crosses the path, from Start to End,
// and none for an empty path.
func (s Segment[T]) IntersectionPolyline(polyline Polyline[T]) []Point[T] {
	return s.AppendIntersectionPolyline(nil, polyline)
}

// AppendIntersectionPolyline appends the points of IntersectionPolyline(polyline) to dst and
// returns the extended slice.
func (s Segment[T]) AppendIntersectionPolyline(dst []Point[T], polyline Polyline[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	for edge := range polyline.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectionQuadraticBezier returns the points where the segment crosses the curve, from Start to
// End.
func (s Segment[T]) IntersectionQuadraticBezier(curve QuadraticBezier[T]) Pair[Point[T]] {
	return curve.IntersectionSegment(s)
}

// IntersectionCubicBezier returns the points where the segment crosses the curve, from Start to
// End, three at most, each within Delta of the curve. A segment along a straight curve gives an end
// of the stretch they share or none, and one touching a bend over a stretch without crossing it can
// give two points where there is one touch.
func (s Segment[T]) IntersectionCubicBezier(curve CubicBezier[T]) []Point[T] {
	return s.AppendIntersectionCubicBezier(nil, curve)
}

// AppendIntersectionCubicBezier appends the points of IntersectionCubicBezier(curve) to dst and
// returns the extended slice.
func (s Segment[T]) AppendIntersectionCubicBezier(dst []Point[T], curve CubicBezier[T]) []Point[T] {
	return curve.appendIntersectionSegment(dst, s)
}

// IntersectionPolygon returns the points where the segment crosses the boundary of the polygon,
// from Start to End. A segment entirely inside gives none while IntersectsPolygon holds, a segment
// along an edge crosses only the edges at its ends, and an empty polygon gives none. Each edge
// takes the tolerance of its own coordinates, so an endpoint near a corner can give a point of a
// long edge beside the crossing of its shorter neighbour.
func (s Segment[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return s.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points of IntersectionPolygon(polygon) to dst and returns
// the extended slice. The points already in dst are left as they are: none is compared with the
// appended ones or ordered with them.
func (s Segment[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	e := edgeIntersections[T]{segment: s, points: dst, from: len(dst)}
	for edge := range polygon.Edges() {
		e.add(edge)
	}

	return e.sorted()
}

// IntersectionRectangle returns the points where the segment crosses the boundary of the rectangle,
// from Start to End: the first and the last crossing, the ends ClipRectangle clips to. A segment
// entirely inside gives none while IntersectsRectangle holds, and a segment along an edge crosses
// only the edges at its ends.
func (s Segment[T]) IntersectionRectangle(rectangle Rectangle[T]) Pair[Point[T]] {
	span := edgeSpan[T]{segment: s}
	for edge := range rectangle.Edges() {
		span.add(edge)
	}

	return span.pair()
}

// IntersectionRegularPolygon returns the points where the segment crosses the boundary of the
// regular polygon, from Start to End: the first and the last crossing, the ends ClipRegularPolygon
// clips to. A segment entirely inside gives none while IntersectsRegularPolygon holds, and an empty
// polygon gives none.
func (s Segment[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) Pair[Point[T]] {
	span := edgeSpan[T]{segment: s}
	for edge := range polygon.Edges() {
		span.add(edge)
	}

	return span.pair()
}

// IntersectionBox returns the points where the segment crosses the boundary of the box, from Start
// to End: the first and the last crossing.
func (s Segment[T]) IntersectionBox(box Box[T]) Pair[Point[T]] {
	span := edgeSpan[T]{segment: s}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		span.add(edge)
	}

	return span.pair()
}

// ClipCircle returns the part of the segment in the circle, and false where they share no point. A
// segment touching the circle is clipped to that one point.
func (s Segment[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	crossings := circle.IntersectionSegment(s)

	return s.clipConvex(crossings.Items[:crossings.Count], circle.Contains(s.Start), circle.Contains(s.End))
}

// ClipCapsule returns the part of the segment in the capsule, and false where they share no point.
// A segment touching the capsule is clipped to that one point.
func (s Segment[T]) ClipCapsule(capsule Capsule[T]) (Segment[T], bool) {
	crossings := capsule.IntersectionSegment(s)

	return s.clipConvex(crossings.Items[:crossings.Count], capsule.Contains(s.Start), capsule.Contains(s.End))
}

// ClipSector returns the parts of the segment in the sector, from Start to End: one for a sector of
// up to half a turn, and two past it where the segment passes through the notch between its radii.
// Parts run together across a touch from inside, and a touch from outside is a part of zero length.
func (s Segment[T]) ClipSector(sector Sector[T]) []Segment[T] {
	return s.AppendClipSector(nil, sector)
}

// AppendClipSector appends the parts of ClipSector(sector) to dst and returns the extended slice.
func (s Segment[T]) AppendClipSector(dst []Segment[T], sector Sector[T]) []Segment[T] {
	cast := s.Float()
	crossed, n := sector.crossings(s)

	fractions := [11]float64{}
	copy(fractions[1:], crossed[:n])
	fractions[n+1] = 1

	var from Point[T]

	open := false
	for i, t := range fractions[:n+2] {
		point, last := cast.PointAt(t).Cast[T](), i == n+1
		switch {
		case i == 0:
			point = s.Start
		case last:
			point = s.End
		}

		if !open && sector.holds(cast.PointAt(t)) {
			from, open = point, true
		}

		if open && (last || !sector.holds(cast.PointAt((t+fractions[i+1])/2))) {
			dst, open = append(dst, Segment[T]{from, point}), false
		}
	}

	return dst
}

// ClipPolygon returns the parts of the segment in the polygon, from Start to End, by the even-odd
// rule of Contains, and none for an empty polygon. Parts run together across a touch from inside,
// such as at a reflex vertex, and a touch from outside is a part of zero length.
//
// A part never starts or ends at a point the polygon does not contain, so a segment leaving the
// boundary at a shallow angle without crossing it ends there. There are parts exactly where
// IntersectsPolygon holds, except for a flat polygon, its vertices on one line: a segment along
// that line with both ends beyond the polygon has none.
func (s Segment[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return s.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts of ClipPolygon(polygon) to dst and returns the extended
// slice.
func (s Segment[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	var from Point[T]

	point, open, ended := s.Start, false, false
	pass := s.pass(polygon, point)
	contained := polygon.Contains(point) || pass.touchedAt(point)
	for {
		next, found := pass.next(point)
		if !found && pass.spent() {
			pass = s.pass(polygon, point)
			next, found = pass.next(point)
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

// ClipRectangle returns the part of the segment in the rectangle, and false where they share no
// point. A segment along an edge is clipped to the stretch of the edge it covers, and one touching
// a corner from outside to that corner.
func (s Segment[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	for edge := range rectangle.Edges() {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), rectangle.Contains(s.Start), rectangle.Contains(s.End))
}

// ClipRegularPolygon returns the part of the segment in the regular polygon, and false where they
// share no point or the polygon is empty. A segment touching a vertex from outside is clipped to
// that vertex. A flat polygon, of two vertices or a zero semi-axis, is the exception: a segment
// along its line with an end beyond it has no part, while IntersectsRegularPolygon holds.
func (s Segment[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	for edge := range polygon.Edges() {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), polygon.Contains(s.Start), polygon.Contains(s.End))
}

// ClipBox returns the part of the segment in the box, and false where they share no point.
func (s Segment[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	span := edgeSpan[T]{segment: s}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		span.add(edge)
	}

	return s.clipConvex(span.crossings(), box.Contains(s.Start), box.Contains(s.End))
}

// HitCircle returns where the segment first reaches the circle from outside, and false where it
// does not or Start lies in it. The normal points from the center, and back along the segment for a
// circle of zero radius.
func (s Segment[T]) HitCircle(circle Circle[T]) (Hit[T], bool) {
	if circle.Contains(s.Start) {
		return Hit[T]{}, false
	}

	fractions, n := circle.fractions(s)
	for _, t := range fractions[:n] {
		if !s.containsAt(t) {
			continue
		}

		point := s.Float().PointAt(Clamp(t, 0, 1))
		normal := point.Subtract(circle.Center.Float())
		if circle.Radius == 0 {
			normal = Vector[float64]{}
		}

		return s.hitAt(point.Cast[T](), normal), true
	}

	return Hit[T]{}, false
}

// HitCapsule returns where the segment first reaches the capsule from outside, and false where it
// does not or Start lies in it. The normal points from the nearest point of the capsule's segment.
// A capsule of zero radius is hit as its segment, and one of zero length as a circle.
func (s Segment[T]) HitCapsule(capsule Capsule[T]) (Hit[T], bool) {
	switch {
	case capsule.Radius == 0:
		return s.HitSegment(capsule.Segment)
	case capsule.Segment.IsPoint():
		return s.HitCircle(Circle[T]{capsule.Segment.Start, capsule.Radius})
	case capsule.Contains(s.Start):
		return Hit[T]{}, false
	}

	fractions, n := capsule.fractions(s)
	if n == 0 {
		return Hit[T]{}, false
	}

	point := s.Float().PointAt(fractions[0])

	return s.hitAt(point.Cast[T](), point.Subtract(capsule.Segment.Float().foot(point))), true
}

// HitSector returns where the segment first reaches the sector from outside, and false where it
// does not or Start lies in it. The normal points from Center on the arc and is the perpendicular
// of a radius on it, the one facing Start where the radii meet at Center. A Start in the notch of a
// sector past half a turn reaches a radius.
func (s Segment[T]) HitSector(sector Sector[T]) (Hit[T], bool) {
	if sector.Contains(s.Start) {
		return Hit[T]{}, false
	}

	fractions, n := sector.crossings(s)
	if n == 0 {
		return Hit[T]{}, false
	}

	cast := s.Float()
	point := cast.PointAt(fractions[0])

	return s.hitAt(point.Cast[T](), sector.normal(point, cast.Start)), true
}

// HitArc returns where the segment first reaches the arc, a curved wall, and false where it does
// not or Start lies on it. The normal is the direction from the center, turned to face the segment.
func (s Segment[T]) HitArc(arc Arc[T]) (Hit[T], bool) {
	if arc.Contains(s.Start) {
		return Hit[T]{}, false
	}

	fractions, n := arc.crossings(s)
	if n == 0 {
		return Hit[T]{}, false
	}

	point := s.Float().PointAt(fractions[0])

	return s.hitAt(point.Cast[T](), point.Subtract(arc.Center.Float())), true
}

// HitSegment returns where the segment first reaches the given one, a wall, and false where it does
// not or Start lies on it. The normal is the perpendicular of the wall, turned to face the segment.
// A wall running along the segment is hit at the end it reaches first, as is a wall of zero length,
// and the normal faces back along the segment.
func (s Segment[T]) HitSegment(segment Segment[T]) (Hit[T], bool) {
	if segment.Contains(s.Start) {
		return Hit[T]{}, false
	}

	point, crossed := s.IntersectionSegment(segment)
	start, end := s.Contains(segment.Start), s.Contains(segment.End)

	switch {
	case start && end:
		point = segment.Start
		if s.compareDistance(segment.End, point) < 0 {
			point = segment.End
		}

		return s.hitAt(point, Vector[float64]{}), true
	case crossed:
		return s.hitAt(point, segment.Float().Normal()), true
	case start:
		return s.hitAt(segment.Start, Vector[float64]{}), true
	case end:
		return s.hitAt(segment.End, Vector[float64]{}), true
	default:
		return Hit[T]{}, false
	}
}

// HitRay returns where the segment first reaches the ray, a wall, and false where it does not or
// Start lies on it.
func (s Segment[T]) HitRay(ray Ray[T]) (Hit[T], bool) {
	return s.HitSegment(ray.reach(s.minMax()))
}

// HitLine returns where the segment first reaches the line, a wall, and false where it does not or
// Start lies on it.
func (s Segment[T]) HitLine(line Line[T]) (Hit[T], bool) {
	return s.HitSegment(line.span(s.minMax()))
}

// HitPolyline returns where the segment first reaches the path, a wall of several edges, and false
// where it does not or Start lies on it. A point two edges share takes the edge that faces Start,
// and the earlier in Edges where both do or neither does.
func (s Segment[T]) HitPolyline(polyline Polyline[T]) (Hit[T], bool) {
	if polyline.Contains(s.Start) {
		return Hit[T]{}, false
	}

	var first Hit[T]
	var wall Segment[T]

	found := false
	for edge := range polyline.Edges() {
		hit, ok := s.HitSegment(edge)

		switch {
		case !ok:
		case !found || hit.Fraction < first.Fraction:
			first, wall, found = hit, edge, true
		case hit.Point == first.Point && edge.faces(s.Start, wall) && !wall.faces(s.Start, edge):
			first, wall = hit, edge
		}
	}

	return first, found
}

// HitQuadraticBezier returns where the segment first reaches the curve, a curved wall, and false
// where it does not or Start lies on it. The normal is the perpendicular of the curve, turned to
// face the segment.
func (s Segment[T]) HitQuadraticBezier(curve QuadraticBezier[T]) (Hit[T], bool) {
	return s.hitBezier(curve.Float().CubicBezier())
}

// HitCubicBezier returns where the segment first reaches the curve, a curved wall, and false where
// it does not or Start lies on it. The normal is the perpendicular of the curve, turned to face the
// segment.
func (s Segment[T]) HitCubicBezier(curve CubicBezier[T]) (Hit[T], bool) {
	return s.hitBezier(curve.Float())
}

// HitPolygon returns where the segment first reaches the polygon from outside, and false where it
// does not, Start lies in it or the polygon is empty. The normal is that of the edge hit, turned to
// face the segment, so the polygon may be wound either way. A vertex takes the edge that faces
// Start, and the earlier in Edges where both do or, in a notch, neither does.
func (s Segment[T]) HitPolygon(polygon Polygon[T]) (Hit[T], bool) {
	if polygon.Contains(s.Start) {
		return Hit[T]{}, false
	}

	span := edgeSpan[T]{segment: s}
	for edge := range polygon.Edges() {
		span.add(edge)
	}

	return span.hit()
}

// HitRectangle returns where the segment first reaches the rectangle from outside, and false where
// it does not or Start lies in it. The normal is the outward one of the edge hit. A corner takes
// the edge that faces Start, so a segment touching it in passing gets the edge it passes, and the
// earlier in Edges where both do; a segment along an edge gets the edge it meets at the corner.
func (s Segment[T]) HitRectangle(rectangle Rectangle[T]) (Hit[T], bool) {
	if rectangle.Contains(s.Start) {
		return Hit[T]{}, false
	}

	span := edgeSpan[T]{segment: s}
	for edge := range rectangle.Edges() {
		span.add(edge)
	}

	return span.hit()
}

// HitRegularPolygon returns where the segment first reaches the regular polygon from outside, and
// false where it does not, Start lies in it or the polygon is empty. The normal is that of the edge
// hit, turned to face the segment. A vertex takes the edge that faces Start, and the earlier in
// Edges where both do.
func (s Segment[T]) HitRegularPolygon(polygon RegularPolygon[T]) (Hit[T], bool) {
	if polygon.Contains(s.Start) {
		return Hit[T]{}, false
	}

	span := edgeSpan[T]{segment: s}
	for edge := range polygon.Edges() {
		span.add(edge)
	}

	return span.hit()
}

// HitBox returns where the segment first reaches the box from outside, and false where it does not
// or Start lies in it.
func (s Segment[T]) HitBox(box Box[T]) (Hit[T], bool) {
	if box.Contains(s.Start) {
		return Hit[T]{}, false
	}

	span := edgeSpan[T]{segment: s}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		span.add(edge)
	}

	return span.hit()
}

// HitEllipse returns where the segment first reaches the ellipse from outside, and false where it
// does not or Start lies in it. The normal is the perpendicular of the boundary. A segment touching
// the ellipse in passing is decided with no tolerance and may answer either way. An ellipse with a
// zero semi-axis is hit as the segment its other axis spans.
func (s Segment[T]) HitEllipse(ellipse Ellipse[T]) (Hit[T], bool) {
	w, h := ellipse.Size.Float().XY()
	if w == 0 || h == 0 {
		axis := Vector[float64]{w, h}

		return s.HitSegment(Segment[T]{ellipse.worldPoint(axis.Negate()), ellipse.worldPoint(axis)})
	}

	if ellipse.Contains(s.Start) {
		return Hit[T]{}, false
	}

	offset := ellipse.localOffset(s.Start).DivideXY(w, h)
	direction := s.Float().Vector().Rotate(-ellipse.Angle).DivideXY(w, h)

	a, b, c := direction.LengthSquared(), offset.Dot(direction), offset.LengthSquared()-1
	discriminant := float64(b*b) - float64(a*c)
	if a == 0 || !(discriminant >= 0) {
		return Hit[T]{}, false
	}

	t := (-b - math.Sqrt(discriminant)) / a
	if !s.containsAt(t) {
		return Hit[T]{}, false
	}

	t = Clamp(t, 0, 1)
	normal := offset.Add(direction.Multiply(t)).DivideXY(w, h).Rotate(ellipse.Angle)

	return s.hitAt(s.Float().PointAt(t).Cast[T](), normal), true
}

// distanceSquaredTo returns the squared distance to the point with no tolerance. Beside the segment
// it squares cross/length rather than dividing cross² by length², so a lattice tangent at a whole
// length lands exactly on the radius.
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

// foot returns the point at NearestFraction with no tolerance, Start or End themselves at the ends.
func (s Segment[T]) foot(point Point[T]) Point[T] {
	switch t := s.NearestFraction(point); t {
	case 0:
		return s.Start
	case 1:
		return s.End
	default:
		return s.PointAt(t)
	}
}

// crosses reports whether the segments properly cross, each with its ends on opposite sides of the
// other. Touching and collinear segments do not; endpoint distances cover them.
func (s Segment[T]) crosses(segment Segment[T]) bool {
	_, ok := s.crossing(segment)

	return ok
}

// crossing returns the point where two properly crossing segments meet, and false where they do
// not: each fraction lies strictly between the ends, decided on the signs of exact cross products
// with no division. A nearly collinear float pair can pass that test with a point anywhere on the
// line, so a point that projects outside the other segment is no crossing and is left to the
// endpoints.
//
// The point is placed by the numerator times each component, divided once, so an integer crossing
// on a half unit is exact. Past the exact range the product rounds, so the pair is first put in one
// canonical order, each from its lesser endpoint by Point.Compare and the lesser segment first, and
// the answer is the same from either side and for a Reverse.
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

// parallel reports whether the segments run in the same direction. A segment of zero length is
// parallel to nothing, so its point can still be found on the other.
func (s Segment[T]) parallel(segment Segment[T]) bool {
	a, b := s.Float().Vector(), segment.Float().Vector()

	return a.hasDirection() && b.hasDirection() && a.Cross(b) == 0
}

// touch returns the endpoint of either segment that lies on the other, and false where there is
// none, the only way segments that do not properly cross share a point. The given segment's
// endpoints come first, so a segment ending at a vertex meets both edges at the vertex itself.
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

// chord returns the fractions along the line of the segment where it enters and leaves the circle,
// not clamped to the segment, and false where it misses or the segment has no direction. The miss
// is decided on the squared gap DistanceSquaredTo computes, so it agrees with IntersectsCircle. A
// chord whose ends would compare Equal is a tangent and gives its midpoint twice; a chord that
// merely grazes within the tolerance is not collapsed.
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

// snapToEndpoint replaces one chord fraction with the endpoint, 0 for Start and 1 for End, so an
// endpoint on the boundary is one of the two crossings rather than a third. Start replaces the
// entry and End the exit, so a segment and its Reverse keep the same crossing, unless the chord
// lies wholly past the endpoint, where it replaces the end of the chord it meets.
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

// pointsAt returns the points at the fractions, at most two and ascending, that lie on the segment
// by containsAt, two that compare Equal given once.
func (s Segment[T]) pointsAt(fractions ...float64) Pair[Point[T]] {
	var points Pair[Point[T]]
	for _, t := range fractions {
		if !s.containsAt(t) {
			continue
		}

		point := s.Float().PointAt(Clamp(t, 0, 1)).Cast[T]()
		if points.Count > 0 && point.Equal(points.Items[0]) {
			continue
		}

		points.Items[points.Count] = point
		points.Count++
	}

	return points
}

// containsAt reports whether the fraction t lies on the segment, the ends widened by the tolerance
// scaled to the length, so it is the same distance Contains applies.
func (s Segment[T]) containsAt(t float64) bool {
	epsilon := ratio(epsilonAt[T](s.magnitude()), s.Length())

	return LessOrEqualDelta(0, t, epsilon) && LessOrEqualDelta(t, 1, epsilon)
}

// compareDistance orders two points by their distance from Start, the order of every boundary
// crossing result. It measures in float64, where an integer squared distance cannot overflow.
func (s Segment[T]) compareDistance(a, b Point[T]) int {
	start := s.Start.Float()

	return cmp.Compare(start.DistanceSquaredTo(a.Float()), start.DistanceSquaredTo(b.Float()))
}

// clipConvex returns the part of the segment in a convex shape, from the crossings of its boundary
// in order from Start and whether the shape contains Start and End: the whole segment where it
// contains both, otherwise the first to the last crossing with a contained endpoint replacing its
// neighbour, and false where it contains neither and nothing crosses.
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

// pass returns the next crossings of the polygon boundary beyond the point reached.
func (s Segment[T]) pass(polygon Polygon[T], reached Point[T]) edgePass[T] {
	pass := edgePass[T]{segment: s, reached: reached}
	for edge := range polygon.Edges() {
		pass.add(edge)
	}

	return pass
}

// hitAt returns the Hit at the point with the normal, of any length and sense, turned against the
// segment and made unit. A normal with no direction becomes the direction back along the segment.
func (s Segment[T]) hitAt(point Point[T], normal Vector[float64]) Hit[T] {
	direction := s.Float().Vector()

	switch {
	case !normal.hasDirection():
		normal = direction.Negate()
	case normal.Dot(direction) > 0:
		normal = normal.Negate()
	}

	return Hit[T]{point, normal.Normalize(), s.fraction(point)}
}

// faces reports whether the edge is turned toward the point at the vertex it shares with the
// neighbour: the point lies strictly on the side of its line the neighbour does not run to, on
// exact signs, so no winding is needed. At a convex corner that is the outer side; at a notch
// neither edge faces.
func (s Segment[T]) faces(point Point[T], neighbour Segment[T]) bool {
	line := s.Line()

	return float64(line.side(point)*(line.side(neighbour.Start)+line.side(neighbour.End))) < 0
}

// hitBezier returns the first hit of the segment on the curve, and false where there is none or
// Start lies on the curve. Start is tested first as a point, since a segment leaving along the
// curve crosses none of its chords there.
func (s Segment[T]) hitBezier(curve CubicBezier[float64]) (Hit[T], bool) {
	start := s.Start.Float()

	origin := bezierCast{segment: Segment[float64]{start, start}, curve: curve}
	origin.step(curve, 0, 1, bezierDepth)

	if origin.count > 0 {
		return Hit[T]{}, false
	}

	cast := bezierCast{segment: s.Float(), curve: curve}
	cast.step(curve, 0, 1, bezierDepth)

	if cast.count == 0 {
		return Hit[T]{}, false
	}

	return s.hitAt(cast.hits[0].Point.Cast[T](), cast.hits[0].Normal), true
}

// fraction returns the fraction at which the point projects onto the line of the segment, not
// clamped, and zero for a segment of zero length.
func (s Segment[T]) fraction(point Point[T]) float64 {
	start, direction := s.Start.Float(), s.Float().Vector()

	return ratio(point.Float().Subtract(start).Dot(direction), direction.LengthSquared())
}

// crossesRay reports whether a ray from the point along +X crosses the segment, the step of the
// even-odd rule of Polygon.Contains. An endpoint on the ray counts only as the lower one, so a
// vertex is counted once. Each side is tested by its own sign of a cross product, so a NaN crosses
// nothing.
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

// Equal reports whether the segments have equal Start and End.
func (s Segment[T]) Equal(segment Segment[T]) bool {
	return s.Start.Equal(segment.Start) && s.End.Equal(segment.End)
}

// IsZero reports whether the segment is the zero value.
func (s Segment[T]) IsZero() bool {
	return s.Equal(Segment[T]{})
}

// IsPoint reports whether Start and End are equal, as in the part a Clip gives where a cast touches
// a shape in passing.
func (s Segment[T]) IsPoint() bool {
	return s.Start.Equal(s.End)
}

// Line returns the line through the segment, from Start as its Origin along Vector. A segment of
// zero length gives a line that is its origin alone.
func (s Segment[T]) Line() Line[T] {
	return Line[T]{s.Start, s.Vector()}
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
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Seg("...)
	buf = s.Start.appendString(buf)
	buf = append(buf, ';')
	buf = s.End.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

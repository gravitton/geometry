package geom

import (
	"math"
	"slices"
)

// Sector is a slice of a disc: the area between the arc about Center from Angle through Sweep and
// the two radii to its ends. The angles are in radians and the sweep is signed; a sweep of a full
// turn or more is the whole disc, and a sweep of zero the single radius at Angle.
//
// The radius is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative radius written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it.
//
// The ends of the arc are no lattice points, so for integer T every distance to the sector is
// measured in float64 and compared at no less than the float64 tolerance.
//
// Past half a turn the sector is not convex: its Support is that of its hull, and the package
// functions over two shapes take it as its two convex halves.
type Sector[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius T        `json:"r"`
	Angle  float64  `json:"a,omitzero"`
	Sweep  float64  `json:"s"`
}

// Sect is shorthand for Sector{center, radius, angle, sweep}, with the radius taken absolute.
func Sect[T Number](center Point[T], radius T, angle, sweep float64) Sector[T] {
	return Sector[T]{center, Abs(radius), angle, sweep}
}

// Centroid returns the center of the area, on the bisector of the sweep, Equal to Center for the
// whole disc.
func (s Sector[T]) Centroid() Point[T] {
	return s.Center.Float().Add(VectorFromAngle(s.Angle+float64(s.turn()/2), s.reach())).Cast[T]()
}

// Area returns the area (r²·|sweep|/2), the sweep held to a full turn.
func (s Sector[T]) Area() float64 {
	radius := float64(s.Radius)

	return radius * radius * math.Abs(s.turn()) / 2
}

// Perimeter returns the length of the arc and the two radii, or the circumference alone for the
// whole disc.
func (s Sector[T]) Perimeter() float64 {
	radius, turn := float64(s.Radius), math.Abs(s.turn())
	if s.isWhole() {
		return turn * radius
	}

	return float64(turn*radius) + 2*radius
}

// Inertia returns the polar second moment of area about the Centroid, the rotational inertia of the
// sector at unit density.
func (s Sector[T]) Inertia() float64 {
	radius, reach := float64(s.Radius), s.reach()

	return float64(radius*radius*radius*radius*math.Abs(s.turn())/4) - float64(s.Area()*reach*reach)
}

// Bounds returns the axis-aligned bounding box. For integer T the ends of the arc are rounded
// outward, except one within the tolerance of a whole coordinate, which takes it.
func (s Sector[T]) Bounds() Box[T] {
	a, b := s.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle the sector is a slice of.
func (s Sector[T]) BoundingCircle() Circle[T] {
	return Circle[T]{s.Center, s.Radius}
}

// minMax returns the corners of Bounds, those of the arc widened to hold the center.
func (s Sector[T]) minMax() (Point[T], Point[T]) {
	a, b := s.Arc().minMax()

	return Point[T]{min(a.X, s.Center.X), min(a.Y, s.Center.Y)}, Point[T]{max(b.X, s.Center.X), max(b.Y, s.Center.Y)}
}

// magnitude returns the largest absolute coordinate of the circle of the sector, the magnitude the
// tolerance of a comparison with the sector is taken from.
func (s Sector[T]) magnitude() float64 {
	return s.Center.magnitude() + float64(s.Radius)
}

// turn returns the sweep held to a full turn either way.
func (s Sector[T]) turn() float64 {
	return Clamp(s.Sweep, -2*math.Pi, 2*math.Pi)
}

// reach returns the distance of the Centroid from Center: two thirds of the radius times the ratio
// of the chord of the sweep to its arc.
func (s Sector[T]) reach() float64 {
	reach, half := 2*float64(s.Radius)/3, math.Abs(s.turn())/2
	if half == 0 {
		return reach
	}

	return reach * math.Sin(half) / half
}

// radii returns the two radii in float64, from Center to each end of the arc, the one to Start
// first. A whole disc has them too, inside it.
func (s Sector[T]) radii() [2]Segment[float64] {
	arc, center := s.Arc(), s.Center.Float()

	return [2]Segment[float64]{{center, arc.pointAt(0)}, {center, arc.pointAt(1)}}
}

// corners returns Center and the two ends of the arc in float64, Start first.
func (s Sector[T]) corners() [3]Point[float64] {
	arc := s.Arc()

	return [3]Point[float64]{s.Center.Float(), arc.pointAt(0), arc.pointAt(1)}
}

// halves returns the two sectors of half the sweep the sector is made of, the first from Angle,
// both convex.
func (s Sector[T]) halves() [2]Sector[T] {
	half := s.turn() / 2

	return [2]Sector[T]{{s.Center, s.Radius, s.Angle, half}, {s.Center, s.Radius, s.Angle + float64(half), half}}
}

// notch returns the sector of the rest of the turn, convex where this one is past half a turn and a
// single radius where it is the whole disc.
func (s Sector[T]) notch() Sector[T] {
	turn := s.turn()

	return Sector[T]{s.Center, s.Radius, s.Angle + turn, math.Copysign(2*math.Pi, turn) - turn}
}

// Translate returns the sector moved by the vector.
func (s Sector[T]) Translate(vector Vector[T]) Sector[T] {
	return Sector[T]{s.Center.Add(vector), s.Radius, s.Angle, s.Sweep}
}

// MoveTo returns the sector with its center at the point.
func (s Sector[T]) MoveTo(point Point[T]) Sector[T] {
	return Sector[T]{point, s.Radius, s.Angle, s.Sweep}
}

// Rotate returns the sector turned by the angle in radians about its center, the sweep unchanged
// and the angle normalized to [0, 2π).
func (s Sector[T]) Rotate(angle float64) Sector[T] {
	return Sector[T]{s.Center, s.Radius, NormalizeAngle(s.Angle + angle), s.Sweep}
}

// Scale returns the sector with its radius scaled by the factor, taken absolute.
func (s Sector[T]) Scale(factor float64) Sector[T] {
	return Sector[T]{s.Center, Abs(Multiply(s.Radius, factor)), s.Angle, s.Sweep}
}

// Unscale returns the sector with its radius divided by the factor, taken absolute. It panics for a
// zero factor.
func (s Sector[T]) Unscale(factor float64) Sector[T] {
	return Sector[T]{s.Center, Abs(Divide(s.Radius, factor)), s.Angle, s.Sweep}
}

// Resize returns the sector with the radius, taken absolute.
func (s Sector[T]) Resize(radius T) Sector[T] {
	return Sector[T]{s.Center, Abs(radius), s.Angle, s.Sweep}
}

// Grow returns the sector with its radius increased by the amount, clamped at zero.
func (s Sector[T]) Grow(amount T) Sector[T] {
	return Sector[T]{s.Center, max(s.Radius+amount, 0), s.Angle, s.Sweep}
}

// Shrink returns the sector with its radius decreased by the amount, clamped at zero.
func (s Sector[T]) Shrink(amount T) Sector[T] {
	return Sector[T]{s.Center, max(s.Radius-amount, 0), s.Angle, s.Sweep}
}

// Lerp returns the sector interpolated toward the given one, center, radius and sweep together and
// the angle along the shorter way, normalized to [0, 2π), and extrapolated outside [0, 1]. The
// radius is taken absolute.
func (s Sector[T]) Lerp(sector Sector[T], t float64) Sector[T] {
	return Sector[T]{s.Center.Lerp(sector.Center, t), Abs(Lerp(s.Radius, sector.Radius, t)), NormalizeAngle(LerpAngle(s.Angle, sector.Angle, t)), Lerp(s.Sweep, sector.Sweep, t)}
}

// Canonical returns the sector with its radius taken absolute, its angle normalized to [0, 2π) and
// its sweep held to a full turn. An angle within Delta of zero or a full turn becomes zero.
func (s Sector[T]) Canonical() Sector[T] {
	return Sector[T]{s.Center, Abs(s.Radius), snapAngle(s.Angle), s.turn()}
}

// DistanceTo returns the distance from the point to the sector, zero inside.
func (s Sector[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(s.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (s Sector[T]) DistanceSquaredTo(point Point[T]) float64 {
	_, distanceSquared := s.nearest(point.Float())
	if s.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return distanceSquared
}

// DistanceToSector returns the distance between the nearest points of the sectors, zero where they
// share a point.
func (s Sector[T]) DistanceToSector(sector Sector[T]) float64 {
	return math.Sqrt(s.DistanceSquaredToSector(sector))
}

// DistanceSquaredToSector returns the square of DistanceToSector.
func (s Sector[T]) DistanceSquaredToSector(sector Sector[T]) float64 {
	distanceSquared := s.distanceSquaredToSector(sector)
	if s.containsSquared(distanceSquared, sector.magnitude()) {
		return 0
	}

	return distanceSquared
}

// Nearest returns the point of the sector nearest to the given point, the point itself inside or
// with a NaN coordinate.
func (s Sector[T]) Nearest(point Point[T]) Point[T] {
	nearest, distanceSquared := s.nearest(point.Float())
	if s.containsSquared(distanceSquared, point.magnitude()) || math.IsNaN(distanceSquared) {
		return point
	}

	return nearest.Cast[T]()
}

// Support returns the point of the sector farthest along the direction, whatever its length, and
// Center for a zero direction. Past half a turn it is the Support of the hull.
func (s Sector[T]) Support(direction Vector[float64]) Point[float64] {
	arc, center := s.Arc(), s.Center.Float()
	if direction.hasDirection() && arc.covers(direction.Angle()) {
		return center.Add(direction.Normalize().Multiply(float64(s.Radius)))
	}

	return supportOf([]Point[float64]{center, arc.pointAt(0), arc.pointAt(1)}, direction)
}

// Contains reports whether the point lies in the sector. The tolerance is on the distance, never on
// an angle, so it does not widen far from the center.
func (s Sector[T]) Contains(point Point[T]) bool {
	return s.DistanceSquaredTo(point) == 0
}

// ContainsCircle reports whether the circle lies entirely in the sector.
func (s Sector[T]) ContainsCircle(circle Circle[T]) bool {
	_, toBoundary := s.boundary(circle.Center.Float())

	return s.Contains(circle.Center) && greaterOrEqualSquared(toBoundary, float64(circle.Radius), s.epsilon(circle.magnitude()))
}

// ContainsCapsule reports whether the capsule lies entirely in the sector.
func (s Sector[T]) ContainsCapsule(capsule Capsule[T]) bool {
	segment := capsule.Segment.Float()
	if !s.containsSegment(segment) {
		return false
	}

	_, toBoundary := s.boundary(segment.Start)
	_, toEnd := s.boundary(segment.End)
	toBoundary = min(toBoundary, toEnd)

	if !s.isWhole() {
		for _, side := range s.radii() {
			toBoundary = min(toBoundary, side.DistanceSquaredToSegment(segment))
		}
	}

	return greaterOrEqualSquared(toBoundary, float64(capsule.Radius), s.epsilon(capsule.magnitude()))
}

// ContainsSector reports whether the given sector lies entirely in this one.
func (s Sector[T]) ContainsSector(sector Sector[T]) bool {
	for _, side := range sector.radii() {
		if !s.containsSegment(side) {
			return false
		}
	}

	return s.containsArc(sector.Arc().Float())
}

// ContainsSegment reports whether the segment lies entirely in the sector, not leaving it through
// the notch of a sector past half a turn.
func (s Sector[T]) ContainsSegment(segment Segment[T]) bool {
	return s.containsSegment(segment.Float())
}

// ContainsPolygon reports whether the polygon lies entirely in the sector.
func (s Sector[T]) ContainsPolygon(polygon Polygon[T]) bool {
	for edge := range polygon.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return !polygon.IsEmpty()
}

// ContainsRectangle reports whether the rectangle lies entirely in the sector.
func (s Sector[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for edge := range rectangle.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the sector.
func (s Sector[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	for edge := range polygon.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return !polygon.IsEmpty()
}

// ContainsBox reports whether the box lies entirely in the sector.
func (s Sector[T]) ContainsBox(box Box[T]) bool {
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the sector and the circle share a point.
func (s Sector[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsSector(s)
}

// IntersectsCapsule reports whether the sector and the capsule share a point.
func (s Sector[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsSector(s)
}

// IntersectsSector reports whether the sectors share a point.
func (s Sector[T]) IntersectsSector(sector Sector[T]) bool {
	return s.DistanceSquaredToSector(sector) == 0
}

// IntersectsSegment reports whether the sector and the segment share a point.
func (s Sector[T]) IntersectsSegment(segment Segment[T]) bool {
	return s.containsSquared(s.distanceSquaredToSegment(segment.Float()), segment.magnitude())
}

// IntersectsRay reports whether the sector and the ray share a point.
func (s Sector[T]) IntersectsRay(ray Ray[T]) bool {
	return s.IntersectsSegment(ray.reach(s.minMax()))
}

// IntersectsLine reports whether the sector and the line share a point.
func (s Sector[T]) IntersectsLine(line Line[T]) bool {
	return s.IntersectsSegment(line.span(s.minMax()))
}

// IntersectsPolygon reports whether the sector and the polygon share a point.
func (s Sector[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	a1, b1 := s.minMax()
	a2, b2 := polygon.minMax()

	if polygon.IsEmpty() || !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.Contains(s.Center) {
		return true
	}

	for edge := range polygon.Edges() {
		if s.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsRectangle reports whether the sector and the rectangle share a point.
func (s Sector[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a1, b1 := s.minMax()
	a2, b2 := rectangle.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if rectangle.Contains(s.Center) {
		return true
	}

	for edge := range rectangle.Edges() {
		if s.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsRegularPolygon reports whether the sector and the regular polygon share a point.
func (s Sector[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	a1, b1 := s.minMax()
	a2, b2 := polygon.minMax()

	if polygon.IsEmpty() || !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.Contains(s.Center) {
		return true
	}

	for edge := range polygon.Edges() {
		if s.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsBox reports whether the sector and the box share a point.
func (s Sector[T]) IntersectsBox(box Box[T]) bool {
	a, b := s.minMax()
	if !overlaps(a, b, box.Min, box.Max) {
		return false
	}

	if box.Contains(s.Center) {
		return true
	}

	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if s.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (s Sector[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionSector(s)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (s Sector[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionSector(s)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (s Sector[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionSector(s)
}

// appendIntersectionSegment appends the points at the fractions crossings returns to dst: once where
// two are Equal before or after rounding, and none inside a stretch the segment shares with a
// radius. Segment.AppendIntersectionSector and the ray and line forms call it.
func (s Sector[T]) appendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	cast := segment.Float()
	fractions, n := s.crossings(segment)

	distinct := 0
	for _, t := range fractions[:n] {
		if distinct > 0 && cast.PointAt(t).Equal(cast.PointAt(fractions[distinct-1])) {
			continue
		}

		fractions[distinct] = t
		distinct++
	}

	held := len(dst)
	for i, t := range fractions[:distinct] {
		if i > 0 && i < distinct-1 && s.along(cast, fractions[i-1], t) && s.along(cast, t, fractions[i+1]) {
			continue
		}

		point := cast.PointAt(t).Cast[T]()
		if len(dst) > held && point.Equal(dst[len(dst)-1]) {
			continue
		}

		dst = append(dst, point)
	}

	return dst
}

// along reports whether the segment runs on the boundary between the fractions a and b, read at
// their midpoint, so a crossing inside a stretch it shares with a radius is no crossing.
func (s Sector[T]) along(segment Segment[float64], a, b float64) bool {
	return s.touches(segment.PointAt((a + b) / 2))
}

// nearest returns the point of the sector nearest to the given point and the squared distance to
// it, with no tolerance: the point itself inside, else the nearest point of the boundary.
func (s Sector[T]) nearest(point Point[float64]) (Point[float64], float64) {
	arc, center, radius := s.Arc(), s.Center.Float(), float64(s.Radius)

	if offset := point.Subtract(center); arc.covers(offset.Angle()) && offset.LengthSquared() <= radius*radius {
		return point, 0
	}

	return s.boundary(point)
}

// boundary returns the point of the boundary nearest to the given point and the squared distance to
// it, with no tolerance, the arc first of equals and the arc alone for the whole disc.
func (s Sector[T]) boundary(point Point[float64]) (Point[float64], float64) {
	nearest, distanceSquared := s.Arc().nearest(point)
	if s.isWhole() {
		return nearest, distanceSquared
	}

	for _, side := range s.radii() {
		foot := side.foot(point)
		if toFoot := foot.DistanceSquaredTo(point); toFoot < distanceSquared {
			nearest, distanceSquared = foot, toFoot
		}
	}

	return nearest, distanceSquared
}

// holds reports whether the sector contains the float64 point, the test Contains makes before
// rounding into T.
func (s Sector[T]) holds(point Point[float64]) bool {
	_, distanceSquared := s.nearest(point)

	return s.containsSquared(distanceSquared, point.magnitude())
}

// containsSegment reports whether the segment lies in the sector: its ends and every probe are
// held.
func (s Sector[T]) containsSegment(segment Segment[float64]) bool {
	if !s.holds(segment.Start) || !s.holds(segment.End) {
		return false
	}

	probes, n := s.probes(segment)
	for _, probe := range probes[:n] {
		if !s.holds(probe) {
			return false
		}
	}

	return true
}

// containsArc reports whether the arc lies in the sector: the middle of each part between the
// places it meets the circle or a radius is held.
func (s Sector[T]) containsArc(arc Arc[float64]) bool {
	sector := s.Float()

	fractions := [8]float64{0, 1}

	n := 2
	for point := range sector.Circle().IntersectionCircle(arc.Circle()).All() {
		fractions[n] = arc.fraction(point)
		n++
	}
	for _, side := range sector.radii() {
		crossed, count := arc.Circle().fractions(side)
		for _, t := range crossed[:count] {
			fractions[n] = arc.fraction(side.PointAt(Clamp(t, 0, 1)))
			n++
		}
	}

	slices.Sort(fractions[:n])
	for i := 1; i < n; i++ {
		if !s.holds(arc.pointAt((fractions[i-1] + fractions[i]) / 2)) {
			return false
		}
	}

	return true
}

// containsSupport reports whether a sector past half a turn holds a shape no Contains method names:
// its circle holds the shape, and the notch reaches no deeper into it than the tolerance.
func (s Sector[T]) containsSupport(shape Shape[T]) bool {
	if !containsSupport[T](s.Circle(), shape) {
		return false
	}

	if s.isWhole() {
		return true
	}

	notch := minkowski[T]{a: s.notch(), b: shape}
	if !notch.intersects() {
		return true
	}

	_, depth := notch.depth()

	return depth <= notch.tolerance()
}

// liesLeftOf reports whether the sector lies on the left of the line through the edge, the inside of
// an outline wound as Vector.Normal turns: its support away from that side is on the line within
// the tolerance or short of it. An edge with no direction has no side and holds anything.
func (s Sector[T]) liesLeftOf(edge Segment[float64]) bool {
	along := edge.Vector()
	cross := along.Cross(s.Support(Vector[float64]{along.Y, -along.X}).Subtract(edge.Start))
	if cross >= 0 {
		return true
	}

	gap := cross / math.Sqrt(along.LengthSquared())

	return lessOrEqualSquared(float64(gap*gap), 0, s.epsilon(edge.magnitude()))
}

// enters reports whether the segment passes through the inside of the sector, farther than the
// tolerance from the boundary, so an edge running along it does not enter.
func (s Sector[T]) enters(segment Segment[float64]) bool {
	probes, n := s.probes(segment)
	for _, probe := range probes[:n] {
		_, distanceSquared := s.nearest(probe)
		_, toBoundary := s.boundary(probe)

		if distanceSquared == 0 && !s.containsSquared(toBoundary, probe.magnitude()) {
			return true
		}
	}

	return false
}

// probes returns the middle of each part of the segment between the places it can meet the
// boundary, and how many there are. No part crosses the boundary, so the segment lies where its
// ends and probes do.
func (s Sector[T]) probes(segment Segment[float64]) ([6]Point[float64], int) {
	sector := s.Float()

	fractions := [7]float64{0, 1, segment.fraction(segment.foot(sector.Center))}

	n := 3
	crossed, count := sector.Circle().fractions(segment)
	for _, t := range crossed[:count] {
		fractions[n] = Clamp(t, 0, 1)
		n++
	}
	for _, side := range sector.radii() {
		if point, ok := segment.IntersectionSegment(side); ok {
			fractions[n] = segment.fraction(point)
			n++
		}
	}

	slices.Sort(fractions[:n])

	var probes [6]Point[float64]
	for i := 1; i < n; i++ {
		probes[i-1] = segment.PointAt((fractions[i-1] + fractions[i]) / 2)
	}

	return probes, n - 1
}

// crossings returns the fractions along the segment, ascending and clamped to it, at which it meets
// the boundary, and how many there are, a place repeated where two meet. The candidates are its
// ends, its feet from the three corners, its crossings of the circle and of the radii; those on the
// boundary are kept. IntersectionSegment, Segment.ClipSector and Segment.HitSector place points
// from them.
func (s Sector[T]) crossings(segment Segment[T]) ([9]float64, int) {
	sector, cast := s.Float(), segment.Float()

	places, n := [9]float64{0, 1}, 2
	for _, corner := range sector.corners() {
		places[n] = cast.fraction(cast.foot(corner))
		n++
	}

	crossed, count := s.Circle().fractions(segment)
	n += copy(places[n:], crossed[:count])

	for _, side := range sector.radii() {
		if point, ok := cast.IntersectionSegment(side); ok {
			places[n] = cast.fraction(point)
			n++
		}
	}

	kept := 0
	for _, t := range places[:n] {
		if !segment.containsAt(t) {
			continue
		}

		t = Clamp(t, 0, 1)
		if point := cast.PointAt(t); s.touches(point) {
			places[kept] = t
			kept++
		}
	}

	slices.Sort(places[:kept])

	return places, kept
}

// touches reports whether the float64 point lies on the boundary, within the tolerance.
func (s Sector[T]) touches(point Point[float64]) bool {
	_, toBoundary := s.boundary(point)

	return s.containsSquared(toBoundary, point.magnitude())
}

// normal returns the direction the boundary faces at the point, of any length and either sense, for
// a cast from the given start: radial on the arc, reached only from outside the circle, and else
// the normal of the nearer radius. At Center it is that of the radius facing the start, the radius
// to Start where both or neither do. Arc and radii are told apart at the tolerance.
func (s Sector[T]) normal(point, from Point[float64]) Vector[float64] {
	center, radius := s.Center.Float(), float64(s.Radius)

	_, toArc := s.Arc().nearest(point)
	if outside := from.DistanceSquaredTo(center) > radius*radius; s.isWhole() || outside && s.containsSquared(toArc, point.magnitude()) {
		return point.Subtract(center)
	}

	sides := s.radii()
	first, second := sides[0].distanceSquaredTo(point), sides[1].distanceSquaredTo(point)

	switch {
	case s.containsSquared(first, point.magnitude()) && s.containsSquared(second, point.magnitude()):
		if sides[1].faces(from, sides[0]) && !sides[0].faces(from, sides[1]) {
			return sides[1].Normal()
		}

		return sides[0].Normal()
	case second < first:
		return sides[1].Normal()
	default:
		return sides[0].Normal()
	}
}

// distanceSquaredToSegment returns the squared distance from the sector to the segment, with no
// tolerance but that of two segments meeting: the least over its ends, its foot from Center and the
// two radii. A segment reaching into the sector gives zero through one of them.
func (s Sector[T]) distanceSquaredToSegment(segment Segment[float64]) float64 {
	distanceSquared := math.Inf(1)
	for _, point := range [3]Point[float64]{segment.Start, segment.End, segment.foot(s.Center.Float())} {
		_, toPoint := s.nearest(point)
		distanceSquared = min(distanceSquared, toPoint)
	}
	for _, side := range s.radii() {
		distanceSquared = min(distanceSquared, side.DistanceSquaredToSegment(segment))
	}

	return distanceSquared
}

// distanceSquaredToSector returns the squared distance between the sectors, with no tolerance but
// that of two segments meeting: the least over each one's radii against the other, and the arcs at
// their crossings and on the line through the centers, each measured from a point of its own
// sector.
func (s Sector[T]) distanceSquaredToSector(sector Sector[T]) float64 {
	a, b := s.Float(), sector.Float()

	distanceSquared := math.Inf(1)
	for _, side := range b.radii() {
		distanceSquared = min(distanceSquared, a.distanceSquaredToSegment(side))
	}
	for _, side := range a.radii() {
		distanceSquared = min(distanceSquared, b.distanceSquaredToSegment(side))
	}

	along := b.Center.Subtract(a.Center).Normalize()
	places := [6]Point[float64]{
		a.Center.Add(along.Multiply(a.Radius)), a.Center.Add(along.Multiply(-a.Radius)),
		b.Center.Add(along.Multiply(b.Radius)), b.Center.Add(along.Multiply(-b.Radius)),
	}

	n := 4
	for point := range a.Circle().IntersectionCircle(b.Circle()).All() {
		places[n] = point
		n++
	}

	for _, place := range places[:n] {
		own, _ := a.Arc().nearest(place)
		_, toOwn := b.nearest(own)
		other, _ := b.Arc().nearest(place)
		_, toOther := a.nearest(other)
		distanceSquared = min(distanceSquared, toOwn, toOther)
	}

	return distanceSquared
}

// containsSquared reports whether a point at the squared distance from the sector lies in it,
// within the tolerance.
func (s Sector[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, 0, s.epsilon(magnitude))
}

// epsilon returns the tolerance of a comparison between the sector and geometry of the given
// magnitude, never less than that of float64, since the ends of the arc are no lattice points.
func (s Sector[T]) epsilon(magnitude float64) float64 {
	return s.Arc().epsilon(magnitude)
}

// Equal reports whether the sectors have equal center, radius, angle and sweep. The angles are
// compared with EqualAngle, the sweeps as numbers.
func (s Sector[T]) Equal(sector Sector[T]) bool {
	return s.Center.Equal(sector.Center) && Equal(s.Radius, sector.Radius) && EqualAngle(s.Angle, sector.Angle) && Equal(s.Sweep, sector.Sweep)
}

// IsZero reports whether the sector is the zero value.
func (s Sector[T]) IsZero() bool {
	return s.Equal(Sector[T]{})
}

// isWhole reports whether the sector sweeps a full turn, whose radii are no part of its boundary.
func (s Sector[T]) isWhole() bool {
	return math.Abs(s.turn()) == 2*math.Pi
}

// isReflex reports whether the sector sweeps more than half a turn, where it is not convex.
func (s Sector[T]) isReflex() bool {
	return math.Abs(s.turn()) > math.Pi
}

// Circle returns the circle the sector is a slice of.
func (s Sector[T]) Circle() Circle[T] {
	return Circle[T]{s.Center, s.Radius}
}

// Arc returns the arc of the sector, from one radius to the other.
func (s Sector[T]) Arc() Arc[T] {
	return Arc[T](s)
}

// Polygon returns the sector as a polygon within the tolerance of it: Center and the arc as few
// equal chords as keep within it, at most 1024 and at least three per full turn, wound like
// Rectangle.Vertices whichever way the sector sweeps. The whole disc leaves out Center. A negative
// tolerance is taken absolute.
func (s Sector[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{s.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of Polygon(tolerance) to dst and returns the extended slice.
func (s Sector[T]) AppendPolygon(dst []Point[T], tolerance float64) []Point[T] {
	arc := s.Arc()
	if arc.turn() < 0 {
		arc = arc.Reverse()
	}

	pieces := max(arc.pieces(tolerance), math.Ceil(arc.turn()/(2*math.Pi)*3))
	from := len(dst)

	dst = slices.Grow(dst, int(pieces)+2)
	if !s.isWhole() {
		dst = append(dst, s.Center)
	}

	dst = arc.appendPieces(dst, pieces)
	dst = dst[:from+len(slices.Compact(dst[from:]))]

	if last := len(dst) - 1; last > from && dst[last].Equal(dst[from]) {
		dst = dst[:last]
	}

	return dst
}

// Cast converts the sector to a Sector of another number type, rounding as Cast does and keeping
// the angle and the sweep.
func (s Sector[T]) Cast[R Number]() Sector[R] {
	return Sector[R]{s.Center.Cast[R](), Cast[R](float64(s.Radius)), s.Angle, s.Sweep}
}

// Int converts the sector to a Sector[int].
func (s Sector[T]) Int() Sector[int] {
	return Sector[int]{s.Center.Int(), Int(s.Radius), s.Angle, s.Sweep}
}

// Float converts the sector to a Sector[float64].
func (s Sector[T]) Float() Sector[float64] {
	return Sector[float64]{s.Center.Float(), float64(s.Radius), s.Angle, s.Sweep}
}

// String returns the sector in the form of its constructor: Sect((x,y);r;a;s).
func (s Sector[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Sect("...)
	buf = s.Center.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, s.Radius)
	buf = append(buf, ';')
	buf = appendString(buf, s.Angle)
	buf = append(buf, ';')
	buf = appendString(buf, s.Sweep)
	buf = append(buf, ')')

	return string(buf)
}

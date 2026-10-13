package geom

import (
	"math"
	"slices"
)

// Sector is a 2D circular sector, a slice of a disc: the area between the arc of the circle
// about Center with the given Radius, from Angle through Sweep, and the two radii to its ends.
// It is a vision cone, a pie slice or the swept area of a turning arm. The angles are in
// radians and the sweep is signed, as on Arc; a sweep of a full turn or more is the whole disc
// and a sweep of zero the single radius at Angle.
//
// The radius is never negative: Sect takes it absolute and Scale takes a negative factor
// absolute. A negative radius can only be written as a struct literal or decoded from JSON, and
// Contains, DistanceTo and Bounds give no meaningful answer for it; Canonical repairs it.
//
// The angles place the ends of the arc off the lattice, so for an integer T every distance to
// the sector is measured in float64 and compared at the tolerance of float64, where the other
// shapes compare an integer exactly.
//
// The sector is a Shape and a Body, the area a spell or a field of view covers, and one of the
// named kinds: it has an Intersects and a Contains method for every other and each has one for
// it, decided on the distance to its arc and its two radii, so a sector of any sweep is tested
// as what it is, and a segment, a ray and a line are cast at it by their Hit, Intersection and
// Clip methods naming it. Past half a turn it is not convex, where its Support is that of its
// hull; the functions over two shapes take such a sector as its two halves, each convex.
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

// Centroid returns the center of the enclosed area: the point on the bisector of the sweep two
// thirds of the radius from Center for a thin sector, nearer as it widens, and Center itself
// for the whole disc. For integer T it is rounded once.
func (s Sector[T]) Centroid() Point[T] {
	return s.Center.Float().Add(VectorFromAngle(s.Angle+float64(s.turn()/2), s.reach())).Cast[T]()
}

// Area returns the enclosed area, half the squared radius times the sweep. It is a float64
// even for an integer T, like the area of a Circle.
func (s Sector[T]) Area() float64 {
	radius := float64(s.Radius)

	return radius * radius * math.Abs(s.turn()) / 2
}

// Perimeter returns the length of the boundary: the arc and the two radii, or the
// circumference alone for the whole disc, which has no radius on its boundary.
func (s Sector[T]) Perimeter() float64 {
	radius, turn := float64(s.Radius), math.Abs(s.turn())
	if s.isWhole() {
		return turn * radius
	}

	return float64(turn*radius) + 2*radius
}

// Inertia returns the polar second moment of area about the Centroid, the rotational inertia
// of the enclosed area at unit density: the moment about Center, a quarter of the fourth power
// of the radius times the sweep, moved to the centroid by the parallel axis theorem.
func (s Sector[T]) Inertia() float64 {
	radius, reach := float64(s.Radius), s.reach()

	return float64(radius*radius*radius*radius*math.Abs(s.turn())/4) - float64(s.Area()*reach*reach)
}

// Bounds returns the axis-aligned bounding box: the box around Center, the two ends of the arc
// and each of the four points of the circle along the axes that the sweep passes. For integer T
// the ends are rounded outward, as Arc.Bounds rounds them, so the box holds the sector.
func (s Sector[T]) Bounds() Box[T] {
	a, b := s.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle around the sector, the Circle it is a slice of.
func (s Sector[T]) BoundingCircle() Circle[T] {
	return Circle[T]{s.Center, s.Radius}
}

// minMax returns the minimum and maximum corner of the sector, the corners of Bounds: those of
// its arc, moved out to hold the center.
func (s Sector[T]) minMax() (Point[T], Point[T]) {
	a, b := s.Arc().minMax()

	return Point[T]{min(a.X, s.Center.X), min(a.Y, s.Center.Y)}, Point[T]{max(b.X, s.Center.X), max(b.Y, s.Center.Y)}
}

// magnitude returns the largest absolute coordinate of the circle the sector is part of: the
// size epsilonAt widens the tolerance by for a comparison that reads the sector.
func (s Sector[T]) magnitude() float64 {
	return s.Center.magnitude() + float64(s.Radius)
}

// turn returns the sweep of the sector held to a full turn either way, as Arc holds its own.
func (s Sector[T]) turn() float64 {
	return Clamp(s.Sweep, -2*math.Pi, 2*math.Pi)
}

// reach returns how far the Centroid lies from Center along the bisector of the sweep: two
// thirds of the radius times the ratio of the chord of the sweep to its arc, which is one for
// a sweep of zero and zero for a full turn.
func (s Sector[T]) reach() float64 {
	reach, half := 2*float64(s.Radius)/3, math.Abs(s.turn())/2
	if half == 0 {
		return reach
	}

	return reach * math.Sin(half) / half
}

// radii returns the two radii on the boundary of the sector in float64, each from Center to an
// end of the arc, the first to its Start: the two straight sides every distance to the sector
// is measured to. A whole disc has them too, lying within it.
func (s Sector[T]) radii() [2]Segment[float64] {
	arc, center := s.Arc(), s.Center.Float()

	return [2]Segment[float64]{{center, arc.pointAt(0)}, {center, arc.pointAt(1)}}
}

// corners returns the three points the straight sides of the sector meet at, in float64: Center
// and the two ends of the arc, its Start first.
func (s Sector[T]) corners() [3]Point[float64] {
	arc := s.Arc()

	return [3]Point[float64]{s.Center.Float(), arc.pointAt(0), arc.pointAt(1)}
}

// halves returns the two sectors of half the sweep each that the sector is made of, the first
// from Angle: both convex, whatever the sweep, and together every point of the sector.
func (s Sector[T]) halves() [2]Sector[T] {
	half := s.turn() / 2

	return [2]Sector[T]{{s.Center, s.Radius, s.Angle, half}, {s.Center, s.Radius, s.Angle + float64(half), half}}
}

// Translate creates a new Sector translated by the given vector.
func (s Sector[T]) Translate(vector Vector[T]) Sector[T] {
	return Sector[T]{s.Center.Add(vector), s.Radius, s.Angle, s.Sweep}
}

// MoveTo creates a new Sector with the same radius, angle and sweep and the center, the point
// the two radii meet at, set to point.
func (s Sector[T]) MoveTo(point Point[T]) Sector[T] {
	return Sector[T]{point, s.Radius, s.Angle, s.Sweep}
}

// Rotate creates a new Sector turned by the given angle (in radians) about its center, in the
// same sense as Vector.Rotate: the angle it starts at moves and the sweep stays. The stored
// angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (s Sector[T]) Rotate(angle float64) Sector[T] {
	return Sector[T]{s.Center, s.Radius, NormalizeAngle(s.Angle + angle), s.Sweep}
}

// Scale creates a new Sector with radius scaled by the given factor, about its center. A
// negative factor scales by its absolute value.
func (s Sector[T]) Scale(factor float64) Sector[T] {
	return Sector[T]{s.Center, Abs(Multiply(s.Radius, factor)), s.Angle, s.Sweep}
}

// Unscale creates a new Sector with radius scaled by the inverse factor, the inverse of Scale,
// a negative factor by its absolute value like it. Like Divide it panics for a zero factor.
func (s Sector[T]) Unscale(factor float64) Sector[T] {
	return Sector[T]{s.Center, Abs(Divide(s.Radius, factor)), s.Angle, s.Sweep}
}

// Resize creates a new Sector with the given radius, taken absolute like Sect.
func (s Sector[T]) Resize(radius T) Sector[T] {
	return Sector[T]{s.Center, Abs(radius), s.Angle, s.Sweep}
}

// Grow creates a new Sector with radius increased by amount, clamped to zero: the arc moves out
// and the two radii lengthen, the center staying where it is.
func (s Sector[T]) Grow(amount T) Sector[T] {
	return Sector[T]{s.Center, max(s.Radius+amount, 0), s.Angle, s.Sweep}
}

// Shrink creates a new Sector with radius decreased by amount, clamped to zero.
func (s Sector[T]) Shrink(amount T) Sector[T] {
	return Sector[T]{s.Center, max(s.Radius-amount, 0), s.Angle, s.Sweep}
}

// Lerp creates a new Sector in linear interpolation towards the given sector, moving the
// center, the radius and the sweep together and turning the angle along the shorter arc with
// LerpAngle, normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like Point.Lerp,
// with the radius taken absolute.
func (s Sector[T]) Lerp(sector Sector[T], t float64) Sector[T] {
	return Sector[T]{s.Center.Lerp(sector.Center, t), Abs(Lerp(s.Radius, sector.Radius, t)), NormalizeAngle(LerpAngle(s.Angle, sector.Angle, t)), Lerp(s.Sweep, sector.Sweep, t)}
}

// Canonical creates a new Sector in the form Sect and Rotate build, with the radius taken
// absolute, the angle normalized to [0, 2π) and a sweep past a full turn held to it: a
// well-formed sector is returned as it is, up to the full turns Equal already ignores in the
// angle. An angle within Delta of zero or of a full turn becomes exactly zero.
func (s Sector[T]) Canonical() Sector[T] {
	return Sector[T]{s.Center, Abs(s.Radius), snapAngle(s.Angle), s.turn()}
}

// DistanceTo returns the distance from the given point to the nearest point of the sector:
// zero exactly where Contains holds, and otherwise the distance to the nearest of its arc and
// its two radii.
func (s Sector[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(s.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons: zero for a point within the sector, or within the tolerance of its boundary,
// and the squared distance to the nearest of the arc and the two radii otherwise. It is a
// float64 even for an integer T, like every shape with an edge. Contains and IntersectsCircle
// are built on it.
func (s Sector[T]) DistanceSquaredTo(point Point[T]) float64 {
	_, distanceSquared := s.nearest(point.Float())
	if s.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return distanceSquared
}

// DistanceToSector returns the distance between the nearest points of the two sectors: zero
// exactly where IntersectsSector holds, and otherwise the least between the arc and the two
// radii of one and those of the other, whatever the two sweeps.
func (s Sector[T]) DistanceToSector(sector Sector[T]) float64 {
	return math.Sqrt(s.DistanceSquaredToSector(sector))
}

// DistanceSquaredToSector returns the squared distance DistanceToSector takes the root of,
// faster for comparisons. A distance within the tolerance is snapped to zero, which is where
// IntersectsSector reads it.
func (s Sector[T]) DistanceSquaredToSector(sector Sector[T]) float64 {
	distanceSquared := s.distanceSquaredToSector(sector)
	if s.containsSquared(distanceSquared, sector.magnitude()) {
		return 0
	}

	return distanceSquared
}

// Nearest returns the point of the sector nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the nearest point of its boundary, on the arc or on a
// radius. For integer T it is rounded once and can land off the sector, where Contains rejects
// it. A point with a NaN coordinate has no nearest point and is returned as it is.
func (s Sector[T]) Nearest(point Point[T]) Point[T] {
	nearest, distanceSquared := s.nearest(point.Float())
	if s.containsSquared(distanceSquared, point.magnitude()) || math.IsNaN(distanceSquared) {
		return point
	}

	return nearest.Cast[T]()
}

// Support returns the point of the sector farthest along the given direction, in float64: the
// point of the arc that way where the sweep passes it, whatever the length of the direction,
// and otherwise the farthest of Center and the two ends of the arc, Center first where several
// reach alike, as for a direction of no length. Past half a turn it is the Support of the hull,
// which the notch between the two radii is filled in by.
func (s Sector[T]) Support(direction Vector[float64]) Point[float64] {
	arc, center := s.Arc(), s.Center.Float()
	if direction.hasDirection() && arc.covers(direction.Angle()) {
		return center.Add(direction.Normalize().Multiply(float64(s.Radius)))
	}

	return supportOf([]Point[float64]{center, arc.pointAt(0), arc.pointAt(1)}, direction)
}

// Contains reports whether the given point lies within the sector, boundary included within
// the tolerance: it holds exactly where DistanceTo is zero. The tolerance is on the distance to
// the arc and to the two radii, never on an angle, so a point a rounding error outside a radius
// is contained however far from the center it lies.
func (s Sector[T]) Contains(point Point[T]) bool {
	return s.DistanceSquaredTo(point) == 0
}

// ContainsCircle reports whether the given circle lies within the sector: its center does, and
// the arc and the two radii pass no nearer the center than its radius, within the tolerance,
// so a circle touching the boundary from inside is enclosed.
func (s Sector[T]) ContainsCircle(circle Circle[T]) bool {
	_, toBoundary := s.boundary(circle.Center.Float())

	return s.Contains(circle.Center) && greaterOrEqualSquared(toBoundary, float64(circle.Radius), s.epsilon(circle.magnitude()))
}

// ContainsCapsule reports whether the given capsule lies within the sector: its segment does,
// as ContainsSegment decides it, and the arc and the two radii pass no nearer that segment than
// its radius, within the tolerance, so a capsule touching the boundary from inside is enclosed.
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

// ContainsSector reports whether the given sector lies within this one: its two radii do, as
// ContainsSegment decides it, and its arc does not leave between its ends, each part of it
// between two places it meets the boundary of this one being contained, so a sector touching
// the boundary from inside is enclosed and every sector encloses itself.
func (s Sector[T]) ContainsSector(sector Sector[T]) bool {
	for _, side := range sector.radii() {
		if !s.containsSegment(side) {
			return false
		}
	}

	return s.containsArc(sector.Arc().Float())
}

// ContainsSegment reports whether the segment lies within the sector: both endpoints are
// contained, within the tolerance, and it does not leave between them, through the notch of a
// sector past half a turn, each part of it between two places it meets the boundary being
// contained.
func (s Sector[T]) ContainsSegment(segment Segment[T]) bool {
	return s.containsSegment(segment.Float())
}

// ContainsPolygon reports whether the polygon lies within the sector: every edge does, as
// ContainsSegment decides it. An empty polygon is enclosed by nothing.
func (s Sector[T]) ContainsPolygon(polygon Polygon[T]) bool {
	for edge := range polygon.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return !polygon.IsEmpty()
}

// ContainsRectangle reports whether the rectangle lies within the sector: every edge does, as
// ContainsSegment decides it.
func (s Sector[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for edge := range rectangle.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies within the sector: every
// edge does, as ContainsSegment decides it. An empty polygon is enclosed by nothing.
func (s Sector[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	for edge := range polygon.Edges() {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return !polygon.IsEmpty()
}

// ContainsBox reports whether the box lies within the sector: every edge does, as
// ContainsSegment decides it.
func (s Sector[T]) ContainsBox(box Box[T]) bool {
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if !s.containsSegment(edge.Float()) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the sector and the circle share a point, as
// Circle.IntersectsSector does.
func (s Sector[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsSector(s)
}

// IntersectsCapsule reports whether the sector and the capsule share a point: the segment of
// the capsule passes within its radius of the sector. Touching shapes intersect, within the
// tolerance, by the comparison Circle.IntersectsSector makes, on the squared distance
// IntersectsSegment measures.
func (s Sector[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return lessOrEqualSquared(s.distanceSquaredToSegment(capsule.Segment.Float()), float64(capsule.Radius), s.epsilon(capsule.magnitude()))
}

// IntersectsSector reports whether the sectors overlap: a radius of one reaches the other, one
// lying within the other included, or their arcs meet. Touching sectors intersect, within the
// tolerance: it holds exactly where DistanceToSector is zero, the distance being measured in
// float64 for every T, the ends of an arc being no lattice points.
func (s Sector[T]) IntersectsSector(sector Sector[T]) bool {
	return s.DistanceSquaredToSector(sector) == 0
}

// IntersectsSegment reports whether the sector and the segment share a point: an end of the
// segment or its point nearest Center lies within the sector, or it reaches one of the two
// radii. Touching shapes intersect, within the tolerance, by the same comparison Contains
// makes, on the squared distance between the two, which is measured in float64 for every T.
func (s Sector[T]) IntersectsSegment(segment Segment[T]) bool {
	return s.containsSquared(s.distanceSquaredToSegment(segment.Float()), segment.magnitude())
}

// IntersectsRay reports whether the sector and the ray share a point, as IntersectsSegment
// decides it on the reach of the ray past the sector.
func (s Sector[T]) IntersectsRay(ray Ray[T]) bool {
	return s.IntersectsSegment(ray.reach(s.minMax()))
}

// IntersectsLine reports whether the sector and the line share a point, as IntersectsSegment
// decides it on the span of the line past the sector.
func (s Sector[T]) IntersectsLine(line Line[T]) bool {
	return s.IntersectsSegment(line.span(s.minMax()))
}

// IntersectsPolygon reports whether the sector and the polygon share a point: the polygon
// contains Center, so the sector lies within it or crosses out of it, or an edge of the
// polygon meets the sector, as IntersectsSegment decides it. An empty polygon intersects
// nothing.
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

// IntersectsRectangle reports whether the sector and the rectangle overlap: the rectangle
// contains Center, or one of its edges meets the sector, as IntersectsSegment decides it.
func (s Sector[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a1, b1 := s.minMax()
	a2, b2 := rectangle.MinMax()

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

// IntersectsRegularPolygon reports whether the sector and the regular polygon share a point:
// the polygon contains Center, or one of its edges meets the sector, as IntersectsSegment
// decides it. An empty polygon intersects nothing.
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

// IntersectsBox reports whether the sector and the box overlap: the box contains Center, or
// one of its edges meets the sector, as IntersectsSegment decides it.
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

// IntersectionSegment returns the points where the segment crosses the sector boundary, from
// its Start to its End: where it meets the arc and the two radii, a point two of them share
// held once. A sector of no more than half a turn is crossed at most twice, and one past it up
// to four times, by a segment through the notch between its radii. A segment inside crosses
// no boundary and returns none, which IntersectsSegment still reports. A point admitted by the
// tolerance is placed on the boundary: an end of the segment within it of the boundary is a
// crossing, so a segment running along a radius gives the ends of the stretch the two share.
// They are found in float64 and rounded once for an integer T.
func (s Sector[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return s.AppendIntersectionSegment(nil, segment)
}

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns
// the extended slice, so a caller reusing dst allocates nothing once it has room for four.
func (s Sector[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	cast := segment.Float()
	fractions, n := s.crossings(segment)

	held := len(dst)
	for _, t := range fractions[:n] {
		point := cast.PointAt(t).Cast[T]()
		if len(dst) > held && point.Equal(dst[len(dst)-1]) {
			continue
		}

		dst = append(dst, point)
	}

	return dst
}

// IntersectionRay returns the points where the ray crosses the sector boundary, from its Origin
// on, as IntersectionSegment finds them on the reach of the ray past the sector.
func (s Sector[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return s.AppendIntersectionRay(nil, ray)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice.
func (s Sector[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	return s.AppendIntersectionSegment(dst, ray.reach(s.minMax()))
}

// IntersectionLine returns the points where the line crosses the sector boundary, in order
// along its Direction, as IntersectionSegment finds them on the span of the line past the
// sector.
func (s Sector[T]) IntersectionLine(line Line[T]) []Point[T] {
	return s.AppendIntersectionLine(nil, line)
}

// AppendIntersectionLine appends the points IntersectionLine returns to dst and returns the
// extended slice.
func (s Sector[T]) AppendIntersectionLine(dst []Point[T], line Line[T]) []Point[T] {
	return s.AppendIntersectionSegment(dst, line.span(s.minMax()))
}

// nearest returns the point of the sector nearest to the given point, in float64, with the
// squared distance to it and no tolerance applied: the point itself at zero where the sweep
// covers its direction from the center and it lies no farther than the radius, and otherwise
// the nearest of the points its arc and its two radii are nearest at, the arc first of equals.
func (s Sector[T]) nearest(point Point[float64]) (Point[float64], float64) {
	arc, center, radius := s.Arc(), s.Center.Float(), float64(s.Radius)

	if offset := point.Subtract(center); arc.covers(offset.Angle()) && offset.LengthSquared() <= radius*radius {
		return point, 0
	}

	return s.boundary(point)
}

// boundary returns the point of the boundary of the sector nearest to the given point, in
// float64, with the squared distance to it and no tolerance applied: the nearest of the points
// its arc and its two radii are nearest at, the arc first of equals, and of the arc alone for
// the whole disc, whose radii lie within it.
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

// holds reports whether the sector contains the given point in float64, boundary included
// within the tolerance: the comparison Contains makes, on a point not yet rounded into T.
func (s Sector[T]) holds(point Point[float64]) bool {
	_, distanceSquared := s.nearest(point)

	return s.containsSquared(distanceSquared, point.magnitude())
}

// containsSegment reports whether the given segment lies within the sector, as ContainsSegment
// reports it: both ends are held, and so is every probe, the middle of each part of it between
// two places it meets the boundary.
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

// containsArc reports whether the given arc lies within the sector: the middle of each part of
// it between two places it meets the boundary is held, those places being where its circle
// crosses the circle of the sector and each of the two radii.
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

// enters reports whether the given segment passes through the inside of the sector, off its
// boundary by more than the tolerance: a probe of it lies within and not on the arc or a
// radius. An edge of a shape that encloses the sector may touch it and run along it, and never
// enters it.
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

// probes returns the middle of each part of the given segment between two places it meets the
// boundary of the sector, and how many of the six there are: its ends, its point nearest
// Center, and where it crosses the circle and the two radii. No part crosses the boundary, so
// the segment lies where its probes and its ends do.
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

// crossings returns the fractions along the given segment at which it crosses the boundary of
// the sector, in increasing order and held to the segment, and how many of the nine there are,
// one place given more than once where two of them meet. The places are its own two ends and
// its points nearest the three corners, so an end on the boundary within the tolerance is a
// crossing and so is a pass within it of a corner, the fractions Circle.fractions gives for the
// circle, at the tolerance of T, and its crossings of the two radii; those kept lie on the
// boundary within the tolerance, the arc alone for the whole disc. IntersectionSegment,
// Segment.ClipSector and Segment.HitSector place their points from them.
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

// touches reports whether the given point in float64 lies on the boundary of the sector within
// the tolerance, the arc or a radius: the test of a place a cast crosses the boundary at.
func (s Sector[T]) touches(point Point[float64]) bool {
	_, toBoundary := s.boundary(point)

	return s.containsSquared(toBoundary, point.magnitude())
}

// normal returns the direction the boundary of the sector faces at the given point of it, of
// any length and either sense, for a cast that reaches it from the given start: that from
// Center on the arc, which is reached from outside the circle alone, and otherwise the
// perpendicular of the nearer radius. At Center, where the two radii meet, it is that of the
// one that faces the start, as a corner of a polygon is read, and of the radius to Start where
// both do or neither does. The arc and the radii are told apart at the tolerance, never by the
// rounding of two distances as near.
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

// distanceSquaredToSegment returns the squared distance from the sector to the given segment,
// in float64 and with no tolerance applied but that of two segments meeting: the least of the
// distances to the two ends of the segment and to its point nearest Center, and of those from
// the segment to the two radii. A segment that reaches into the sector has one of the three
// points within it or meets a radius, so the distance is zero there with no crossing looked
// for.
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

// distanceSquaredToSector returns the squared distance between the two sectors, in float64 and
// with no tolerance applied but that of two segments meeting: the least of the distances from
// each to the radii of the other, by distanceSquaredToSegment, and between the two arcs, which
// are nearest where they cross or on the line through the two centers. Each arc is measured
// from its point nearest one of those places, a point of its sector wherever the sweep ends,
// so no distance is to a point off the two.
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

// containsSquared reports whether a point at the given squared distance from the sector lies
// within it at the tolerance of the sector and of geometry whose largest absolute coordinate is
// the given magnitude: the snap DistanceSquaredTo and Nearest apply.
func (s Sector[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, 0, s.epsilon(magnitude))
}

// epsilon returns the tolerance of a comparison that reads the sector and geometry whose
// largest absolute coordinate is the given magnitude: that of T, and no less than that of
// float64, which an integer T is held to as well, since the ends of the arc are no lattice
// points and every distance to the sector is measured in float64.
func (s Sector[T]) epsilon(magnitude float64) float64 {
	return s.Arc().epsilon(magnitude)
}

// Equal checks for equal center, radius, angle and sweep with given sector. The angles are
// compared with EqualAngle, so a full turn does not matter, and the sweeps as numbers.
func (s Sector[T]) Equal(sector Sector[T]) bool {
	return s.Center.Equal(sector.Center) && Equal(s.Radius, sector.Radius) && EqualAngle(s.Angle, sector.Angle) && Equal(s.Sweep, sector.Sweep)
}

// IsZero checks if center, radius, angle and sweep are zero.
func (s Sector[T]) IsZero() bool {
	return s.Equal(Sector[T]{})
}

// isWhole reports whether the sector sweeps a full turn, the whole disc, whose two radii lie
// within it and are no part of its boundary.
func (s Sector[T]) isWhole() bool {
	return math.Abs(s.turn()) == 2*math.Pi
}

// isReflex reports whether the sector sweeps more than half a turn, where it is not convex and
// its Support is that of its hull.
func (s Sector[T]) isReflex() bool {
	return math.Abs(s.turn()) > math.Pi
}

// Arc converts the sector to the Arc on its boundary, from one radius to the other.
func (s Sector[T]) Arc() Arc[T] {
	return Arc[T](s)
}

// Circle converts the sector to the Circle it is a slice of.
func (s Sector[T]) Circle() Circle[T] {
	return Circle[T]{s.Center, s.Radius}
}

// Polygon converts the sector to a Polygon within the given tolerance of it, the outline that
// draws or walks one: Center and the arc as equal chords, as few as keep it within the
// tolerance of each, with its vertices on the boundary and wound as the corners of a Rectangle
// are, whichever way the sector sweeps. The whole disc has no Center among them. A negative
// tolerance is taken absolute, and the arc has at most a thousand chords, which a tolerance of
// zero gives. For integer T the vertices are rounded and one rounding onto the vertex before
// it is left out.
func (s Sector[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{s.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of the polygon Polygon returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them.
func (s Sector[T]) AppendPolygon(dst []Point[T], tolerance float64) []Point[T] {
	arc := s.Arc()
	if arc.turn() < 0 {
		arc = arc.Reverse()
	}

	from := len(dst)

	dst = slices.Grow(dst, int(arc.pieces(tolerance))+2)
	if !s.isWhole() {
		dst = append(dst, s.Center)
	}

	dst = arc.AppendPolyline(dst, tolerance)
	dst = dst[:from+len(slices.Compact(dst[from:]))]

	if last := len(dst) - 1; last > from && dst[last] == dst[from] {
		dst = dst[:last]
	}

	return dst
}

// Cast converts the sector to a Sector of another number type, rounding as Cast does and
// keeping the angle and the sweep.
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

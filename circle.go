package geom

import (
	"math"
	"slices"
)

// Circle is a disc of a center and a radius.
//
// The radius is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative radius written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it.
type Circle[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius T        `json:"r"`
}

// Circ is shorthand for Circle{center, radius}, with the radius taken absolute.
func Circ[T Number](center Point[T], radius T) Circle[T] {
	return Circle[T]{center, Abs(radius)}
}

// CircleBetween returns the circle with a and b at the ends of a diameter.
func CircleBetween[T Number](a, b Point[T]) Circle[T] {
	return Circle[T]{a.Midpoint(b), Cast[T](a.Float().DistanceTo(b.Float()) / 2)}
}

// CircleThrough returns the circle through the three points, and false where they lie on one line,
// two or all of them being the same point included. The radius must fit T, which bounds how nearly
// collinear the points may be.
func CircleThrough[T Number](a, b, c Point[T]) (Circle[T], bool) {
	points := [3]Point[T]{a, b, c}
	if (Polygon[T]{points[:]}).isFlat() {
		return Circle[T]{}, false
	}

	origin := a.Float()
	toB, toC := b.Float().Subtract(origin), c.Float().Subtract(origin)
	weighted := toB.Multiply(toC.LengthSquared()).Subtract(toC.Multiply(toB.LengthSquared()))
	center := origin.Add(weighted.Normal().Multiply(1 / (2 * toB.Cross(toC)))).Cast[T]()

	return Circle[T]{center, Cast[T](center.Float().DistanceTo(origin))}, true
}

// Diameter returns twice the radius.
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Anchor returns the point of the boundary in the direction from the center, or the center for
// DirectionNone. For integer T a diagonal anchor is rounded like Direction.Vector and can lie just
// outside a small circle.
func (c Circle[T]) Anchor(direction Direction) Point[T] {
	return c.Center.Add(direction.Vector(c.Radius))
}

// Centroid returns the center.
func (c Circle[T]) Centroid() Point[T] {
	return c.Center
}

// Area returns the area (π·r²).
func (c Circle[T]) Area() float64 {
	radius := float64(c.Radius)

	return Pi * radius * radius
}

// Perimeter returns the circumference (2π·r).
func (c Circle[T]) Perimeter() float64 {
	return 2 * Pi * float64(c.Radius)
}

// Inertia returns the polar second moment of area about the center (π·r⁴/2), the rotational inertia
// of the disc at unit density.
func (c Circle[T]) Inertia() float64 {
	radius := float64(c.Radius)

	return Pi * radius * radius * radius * radius / 2
}

// Bounds returns the axis-aligned bounding box.
func (c Circle[T]) Bounds() Box[T] {
	a, b := c.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle itself.
func (c Circle[T]) BoundingCircle() Circle[T] {
	return c
}

// minMax returns the corners of Bounds, the prefilter of the intersection tests.
func (c Circle[T]) minMax() (Point[T], Point[T]) {
	return c.Center.AddXY(-c.Radius, -c.Radius), c.Center.AddXY(c.Radius, c.Radius)
}

// magnitude returns the largest absolute coordinate of Bounds, the magnitude the tolerance of a
// comparison with the circle is taken from.
func (c Circle[T]) magnitude() float64 {
	return c.Center.magnitude() + float64(c.Radius)
}

// Translate returns the circle moved by the vector.
func (c Circle[T]) Translate(vector Vector[T]) Circle[T] {
	return Circle[T]{c.Center.Add(vector), c.Radius}
}

// MoveTo returns the circle with its center at the point.
func (c Circle[T]) MoveTo(point Point[T]) Circle[T] {
	return Circle[T]{point, c.Radius}
}

// AlignTo returns the circle moved so that its Anchor in the direction lies at the point.
func (c Circle[T]) AlignTo(direction Direction, point Point[T]) Circle[T] {
	return c.Translate(point.Subtract(c.Anchor(direction)))
}

// Rotate returns the circle unchanged: no angle moves a circle about its center.
func (c Circle[T]) Rotate(_ float64) Circle[T] {
	return c
}

// Scale returns the circle with its radius scaled by the factor, taken absolute.
func (c Circle[T]) Scale(factor float64) Circle[T] {
	return Circle[T]{c.Center, Abs(Multiply(c.Radius, factor))}
}

// Unscale returns the circle with its radius divided by the factor, taken absolute. It panics for a
// zero factor.
func (c Circle[T]) Unscale(factor float64) Circle[T] {
	return Circle[T]{c.Center, Abs(Divide(c.Radius, factor))}
}

// Resize returns the circle with the radius, taken absolute.
func (c Circle[T]) Resize(radius T) Circle[T] {
	return Circle[T]{c.Center, Abs(radius)}
}

// Grow returns the circle with its radius increased by the amount, clamped at zero.
func (c Circle[T]) Grow(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius+amount, 0)}
}

// Shrink returns the circle with its radius decreased by the amount, clamped at zero.
func (c Circle[T]) Shrink(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius-amount, 0)}
}

// Lerp returns the circle interpolated toward the given one, center and radius together, and
// extrapolated outside [0, 1]. The radius is taken absolute, so extrapolating past zero grows the
// circle again.
func (c Circle[T]) Lerp(circle Circle[T], t float64) Circle[T] {
	return Circle[T]{c.Center.Lerp(circle.Center, t), Abs(Lerp(c.Radius, circle.Radius, t))}
}

// Canonical returns the circle with its radius taken absolute, the receiver when it already is.
func (c Circle[T]) Canonical() Circle[T] {
	return Circle[T]{c.Center, Abs(c.Radius)}
}

// DistanceTo returns the distance from the point to the circle, zero inside.
func (c Circle[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.centerDistanceSquared(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return math.Sqrt(distanceSquared) - float64(c.Radius)
}

// DistanceSquaredTo returns the square of DistanceTo.
func (c Circle[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := c.DistanceTo(point)

	return distance * distance
}

// DistanceToCircle returns the distance between the nearest points of the circles, zero where they
// share a point.
func (c Circle[T]) DistanceToCircle(circle Circle[T]) float64 {
	if c.IntersectsCircle(circle) {
		return 0
	}

	return math.Sqrt(c.centerDistanceSquared(circle.Center)) - float64(c.Radius) - float64(circle.Radius)
}

// DistanceSquaredToCircle returns the square of DistanceToCircle.
func (c Circle[T]) DistanceSquaredToCircle(circle Circle[T]) float64 {
	distance := c.DistanceToCircle(circle)

	return distance * distance
}

// Nearest returns the point of the circle nearest to the given point, the point itself inside.
func (c Circle[T]) Nearest(point Point[T]) Point[T] {
	distanceSquared := c.centerDistanceSquared(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return point
	}

	center := c.Center.Float()
	offset := point.Float().Subtract(center).Multiply(float64(c.Radius) / math.Sqrt(distanceSquared))

	return center.Add(offset).Cast[T]()
}

// Support returns the point of the circle farthest along the direction, whatever its length, and
// along X for a zero direction.
func (c Circle[T]) Support(direction Vector[float64]) Point[float64] {
	return c.Center.Float().Add(direction.Normalize().Multiply(float64(c.Radius)))
}

// Contains reports whether the point lies in the circle.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.centerDistanceSquared(point), point.magnitude())
}

// ContainsCircle reports whether the given circle lies entirely in this one.
func (c Circle[T]) ContainsCircle(circle Circle[T]) bool {
	far := math.Sqrt(c.centerDistanceSquared(circle.Center)) + float64(circle.Radius)

	return c.containsSquared(far*far, circle.magnitude())
}

// ContainsCapsule reports whether the capsule lies entirely in the circle.
func (c Circle[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return c.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && c.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the sector lies entirely in the circle.
func (c Circle[T]) ContainsSector(sector Sector[T]) bool {
	outer, inner := c.Float(), sector.Float()
	for _, corner := range inner.corners() {
		if !outer.Contains(corner) {
			return false
		}
	}

	return outer.Contains(inner.Support(inner.Center.Subtract(outer.Center)))
}

// ContainsSegment reports whether the segment lies entirely in the circle.
func (c Circle[T]) ContainsSegment(segment Segment[T]) bool {
	return c.Contains(segment.Start) && c.Contains(segment.End)
}

// ContainsPolygon reports whether the polygon lies entirely in the circle.
func (c Circle[T]) ContainsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRectangle reports whether the rectangle lies entirely in the circle.
func (c Circle[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the circle.
func (c Circle[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsBox reports whether the box lies entirely in the circle.
func (c Circle[T]) ContainsBox(box Box[T]) bool {
	for _, corner := range box.corners() {
		if !c.Contains(corner) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the circles share a point.
func (c Circle[T]) IntersectsCircle(circle Circle[T]) bool {
	return lessOrEqualSquared(c.centerDistanceSquared(circle.Center), float64(c.Radius)+float64(circle.Radius), c.epsilonWith(circle.magnitude()))
}

// IntersectsCapsule reports whether the circle and the capsule share a point.
func (c Circle[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return lessOrEqualSquared(capsule.Segment.DistanceSquaredTo(c.Center), float64(c.Radius)+float64(capsule.Radius), c.epsilonWith(capsule.magnitude()))
}

// IntersectsSector reports whether the circle and the sector share a point.
func (c Circle[T]) IntersectsSector(sector Sector[T]) bool {
	return lessOrEqualSquared(sector.DistanceSquaredTo(c.Center), float64(c.Radius), sector.epsilon(c.magnitude()))
}

// IntersectsSegment reports whether the circle and the segment share a point.
func (c Circle[T]) IntersectsSegment(segment Segment[T]) bool {
	return c.containsSquared(segment.DistanceSquaredTo(c.Center), segment.magnitude())
}

// IntersectsRay reports whether the circle and the ray share a point.
func (c Circle[T]) IntersectsRay(ray Ray[T]) bool {
	return c.IntersectsSegment(ray.reach(c.minMax()))
}

// IntersectsLine reports whether the circle and the line share a point.
func (c Circle[T]) IntersectsLine(line Line[T]) bool {
	return c.IntersectsSegment(line.span(c.minMax()))
}

// IntersectsPolygon reports whether the circle and the polygon share a point.
func (c Circle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && polygon.walk(c.Center).reaches(c)
}

// IntersectsRectangle reports whether the circle and the rectangle share a point.
func (c Circle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.walk(c.Center).reaches(c)
}

// IntersectsRegularPolygon reports whether the circle and the regular polygon share a point.
func (c Circle[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && polygon.walk(c.Center).reaches(c)
}

// IntersectsBox reports whether the circle and the box share a point.
func (c Circle[T]) IntersectsBox(box Box[T]) bool {
	return c.containsSquared(box.DistanceSquaredTo(c.Center), box.magnitude())
}

// IntersectionCircle returns the points where the boundaries cross: two for overlapping circles,
// one for tangent ones, and none for circles apart, nested or coincident, centers within the
// tolerance of each other counting as coincident. The second is the mirror of the first across the
// line of centers.
func (c Circle[T]) IntersectionCircle(circle Circle[T]) Pair[Point[T]] {
	center := c.Center.Float()
	direction := circle.Center.Float().Subtract(center)
	distanceSquared := c.centerDistanceSquared(circle.Center)
	r1, r2 := float64(c.Radius), float64(circle.Radius)
	outer, inner := r1+r2, math.Abs(r1-r2)
	epsilon := c.epsilonWith(circle.magnitude())

	if lessOrEqualSquared(distanceSquared, 0, epsilon) || !lessOrEqualSquared(distanceSquared, outer, epsilon) || !greaterOrEqualSquared(distanceSquared, inner, epsilon) {
		return Pair[Point[T]]{}
	}

	distance := math.Sqrt(distanceSquared)

	if external, internal := equalSquared(distanceSquared, outer, epsilon), equalSquared(distanceSquared, inner, epsilon); external || internal {
		point := center.Add(direction.Resize(c.tangent(circle, distance, external)))

		return Pair[Point[T]]{[2]Point[T]{point.Cast[T]()}, 1}
	}

	return c.chordPoints(direction, (float64(r1*r1)-float64(r2*r2)+distanceSquared)/(2*distance))
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End:
// two where it passes through, one where it is tangent or ends inside, and none where it misses. A
// segment entirely inside crosses no boundary, so it gives none while IntersectsSegment holds. An
// endpoint on the boundary is one of the points, so a segment and its Reverse give the same points.
func (c Circle[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	fractions, n := c.fractions(segment)

	return segment.pointsAt(fractions[:n]...)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on: two
// where it passes through, one where it is tangent or starts inside, and none where it misses.
func (c Circle[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return c.IntersectionSegment(ray.reach(c.minMax()))
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector: two where it passes through, one where it is tangent, and none where it misses.
func (c Circle[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return c.IntersectionSegment(line.span(c.minMax()))
}

// Tangents returns the points where the lines through the point touch the circle: two for a point
// outside, the point itself on the boundary, and none inside. The second is the mirror of the first
// across the line from the center to the point.
func (c Circle[T]) Tangents(point Point[T]) Pair[Point[T]] {
	distanceSquared := c.centerDistanceSquared(point)
	magnitude := point.magnitude()

	if c.touchesSquared(distanceSquared, magnitude) {
		return Pair[Point[T]]{[2]Point[T]{point}, 1}
	}
	if c.containsSquared(distanceSquared, magnitude) {
		return Pair[Point[T]]{}
	}

	radius := float64(c.Radius)

	return c.chordPoints(point.Float().Subtract(c.Center.Float()), radius*radius/math.Sqrt(distanceSquared))
}

// OuterTangents returns the segments that touch both circles on the same side, each from its point
// on this circle to its point on the given one: two unless one circle lies within the other or the
// centers coincide, which give none. The second is the mirror of the first across the line of
// centers. A circle of zero radius is touched at its center, and two of them are joined by one
// segment.
func (c Circle[T]) OuterTangents(circle Circle[T]) Pair[Segment[T]] {
	return c.tangents(circle, float64(circle.Radius))
}

// InnerTangents returns the segments that touch the circles on opposite sides and cross between
// them, each from its point on this circle to its point on the given one: two for circles that
// share no point, and none otherwise, touching circles included. The second is the mirror of the
// first across the line of centers. Two circles of zero radius are joined by one segment.
func (c Circle[T]) InnerTangents(circle Circle[T]) Pair[Segment[T]] {
	return c.tangents(circle, -float64(circle.Radius))
}

// Union returns the smallest circle containing both, the one that contains the other where there is
// one, the receiver first. For integer T the radius is rounded up, so the union still contains
// both.
func (c Circle[T]) Union(circle Circle[T]) Circle[T] {
	switch {
	case c.ContainsCircle(circle):
		return c
	case circle.ContainsCircle(c):
		return circle
	}

	a, b := c, circle
	if b.Center.Compare(a.Center) < 0 {
		a, b = b, a
	}

	from, to := a.Center.Float(), b.Center.Float()
	offset := to.Subtract(from)
	distance := offset.Length()
	along := (distance + float64(b.Radius) - float64(a.Radius)) / 2
	center := Point[float64]{from.X + float64(offset.X*along)/distance, from.Y + float64(offset.Y*along)/distance}.Cast[T]()
	reach := max(center.Float().DistanceTo(from)+float64(a.Radius), center.Float().DistanceTo(to)+float64(b.Radius))

	return Circle[T]{center, castUp[T](reach)}
}

// fractions returns the fractions along the segment, ascending and not yet clamped to it, at which
// it crosses the boundary, and how many there are. An endpoint on the boundary stands in for the
// crossing it is, so IntersectionSegment and Segment.HitCircle place points from the same values.
func (c Circle[T]) fractions(segment Segment[T]) ([2]float64, int) {
	entry, exit, ok := segment.chord(c)
	start := c.touchesSquared(c.centerDistanceSquared(segment.Start), segment.magnitude())
	end := c.touchesSquared(c.centerDistanceSquared(segment.End), segment.magnitude())

	switch {
	case start && end && segment.Vector().hasDirection():
		return [2]float64{0, 1}, 2
	case ok && entry < exit:
		if start {
			entry, exit = segment.snapToEndpoint(entry, exit, 0)
		}
		if end {
			entry, exit = segment.snapToEndpoint(entry, exit, 1)
		}

		return [2]float64{entry, exit}, 2
	case start:
		return [2]float64{0}, 1
	case end:
		return [2]float64{1}, 1
	case ok:
		return [2]float64{entry}, 1
	default:
		return [2]float64{}, 0
	}
}

// centerDistanceSquared returns the squared distance from the center to the point, the one value
// every test of the circle compares with its radius.
func (c Circle[T]) centerDistanceSquared(point Point[T]) float64 {
	return c.Center.Float().DistanceSquaredTo(point.Float())
}

// containsSquared reports whether a point at the squared distance from the center lies in the
// circle, within the tolerance.
func (c Circle[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, float64(c.Radius), c.epsilonWith(magnitude))
}

// touchesSquared reports whether a point at the squared distance from the center lies on the
// boundary, within the tolerance.
func (c Circle[T]) touchesSquared(distanceSquared, magnitude float64) bool {
	return equalSquared(distanceSquared, float64(c.Radius), c.epsilonWith(magnitude))
}

// epsilonWith returns the tolerance of a comparison between the circle and geometry of the given
// magnitude, so both sides of a pair read the same one.
func (c Circle[T]) epsilonWith(magnitude float64) float64 {
	return epsilonAt[T](max(c.magnitude(), magnitude))
}

// tangent returns how far along the line of centers the boundaries of two tangent circles meet,
// halfway between their crossings of that line. The exact crossing is ill-conditioned at a tangent,
// and the midpoint stays within half the tolerance of both boundaries.
func (c Circle[T]) tangent(circle Circle[T], distance float64, external bool) float64 {
	r1, r2 := float64(c.Radius), float64(circle.Radius)

	switch {
	case external:
		return (r1 + distance - r2) / 2
	case r1 >= r2:
		return (r1 + distance + r2) / 2
	default:
		return (distance - r1 - r2) / 2
	}
}

// tangents returns the segments touching this circle and the given one, whose radius is passed
// signed: positive for the outer tangents, negative for the inner. There are none where the
// distance between the centers is at most the difference of the radii.
func (c Circle[T]) tangents(circle Circle[T], radius float64) Pair[Segment[T]] {
	distanceSquared := c.centerDistanceSquared(circle.Center)
	reach := float64(c.Radius) - radius

	if lessOrEqualSquared(distanceSquared, math.Abs(reach), c.epsilonWith(circle.magnitude())) {
		return Pair[Segment[T]]{}
	}

	direction := circle.Center.Float().Subtract(c.Center.Float())
	ratio := reach / math.Sqrt(distanceSquared)
	start1, start2 := c.chordEnds(direction, float64(c.Radius)*ratio)
	end1, end2 := circle.chordEnds(direction, radius*ratio)
	if radius < 0 {
		end1, end2 = end2, end1
	}

	first, second := Segment[T]{start1, end1}, Segment[T]{start2, end2}
	if second.Equal(first) {
		return Pair[Segment[T]]{[2]Segment[T]{first}, 1}
	}

	return Pair[Segment[T]]{[2]Segment[T]{first, second}, 2}
}

// chordPoints returns chordEnds stored into T, once where both round to the same point.
func (c Circle[T]) chordPoints(direction Vector[float64], along float64) Pair[Point[T]] {
	first, second := c.chordEnds(direction, along)
	if second.Equal(first) {
		return Pair[Point[T]]{[2]Point[T]{first}, 1}
	}

	return Pair[Point[T]]{[2]Point[T]{first, second}, 2}
}

// chordEnds returns the ends of the chord perpendicular to the direction at the given distance
// along it from the center, negative behind it: the first a quarter turn from the direction, the
// second its mirror.
func (c Circle[T]) chordEnds(direction Vector[float64], along float64) (Point[T], Point[T]) {
	radius := float64(c.Radius)
	middle := c.Center.Float().Add(direction.Resize(along))
	normal := direction.Normal().Resize(math.Sqrt(max(float64(radius*radius)-float64(along*along), 0)))

	return middle.Add(normal).Cast[T](), middle.Add(normal.Negate()).Cast[T]()
}

// Equal reports whether the circles have equal center and radius.
func (c Circle[T]) Equal(circle Circle[T]) bool {
	return c.Center.Equal(circle.Center) && Equal(c.Radius, circle.Radius)
}

// IsZero reports whether the circle is the zero value.
func (c Circle[T]) IsZero() bool {
	return c.Equal(Circle[T]{})
}

// Sector returns the slice of the circle from the angle through the sweep, both in radians.
func (c Circle[T]) Sector(angle, sweep float64) Sector[T] {
	return Sector[T]{c.Center, c.Radius, angle, sweep}
}

// Arc returns the part of the boundary from the angle through the sweep, both in radians. A sweep
// of a full turn is the whole boundary without the inside.
func (c Circle[T]) Arc(angle, sweep float64) Arc[T] {
	return Arc[T]{c.Center, c.Radius, angle, sweep}
}

// Polygon returns the circle as a polygon within the tolerance of it: as few equal chords as keep
// within it, at least three and at most 1024, from the point at angle zero and wound like
// Rectangle.Vertices. A negative tolerance is taken absolute, and a circle of zero radius gives its
// center alone. RegularPolygon chooses the vertex count instead.
func (c Circle[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{c.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of Polygon(tolerance) to dst and returns the extended slice.
func (c Circle[T]) AppendPolygon(dst []Point[T], tolerance float64) []Point[T] {
	arc := c.Arc(0, 2*math.Pi)
	pieces := max(arc.pieces(tolerance), 3)

	from := len(dst)
	dst = slices.Grow(dst, int(pieces))

	for i := 0.0; i < pieces; i++ {
		if point := arc.PointAt(i / pieces); len(dst) == from || point != dst[len(dst)-1] {
			dst = append(dst, point)
		}
	}

	if last := len(dst) - 1; last > from && dst[last].Equal(dst[from]) {
		dst = dst[:last]
	}

	return dst
}

// RegularPolygon returns the regular polygon of n vertices inscribed in the circle, its first
// vertex placed by RegularPolygonOrientationPhase. It panics for an orientation that is neither
// OrientationFlatTop nor OrientationPointyTop.
func (c Circle[T]) RegularPolygon(n int, orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{c.Center, SzU(c.Radius), n, 0, RegularPolygonOrientationPhase(n, orientation)}
}

// Ellipse returns the circle as an Ellipse of equal semi-axes and no angle. A circle has no
// Transform, since a matrix takes it to an ellipse: transform Ellipse() instead.
func (c Circle[T]) Ellipse() Ellipse[T] {
	return Ellipse[T]{c.Center, SzU(c.Radius), 0}
}

// Cast converts the circle to a Circle of another number type, rounding as Cast does.
func (c Circle[T]) Cast[R Number]() Circle[R] {
	return Circle[R]{c.Center.Cast[R](), Cast[R](float64(c.Radius))}
}

// Int converts the circle to a Circle[int].
func (c Circle[T]) Int() Circle[int] {
	return Circle[int]{c.Center.Int(), Int(c.Radius)}
}

// Float converts the circle to a Circle[float64].
func (c Circle[T]) Float() Circle[float64] {
	return Circle[float64]{c.Center.Float(), float64(c.Radius)}
}

// String returns the circle in the form of its constructor: Circ((x,y);r).
func (c Circle[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Circ("...)
	buf = c.Center.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, c.Radius)
	buf = append(buf, ')')

	return string(buf)
}

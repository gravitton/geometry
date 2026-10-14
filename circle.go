package geom

import (
	"math"
	"slices"
)

// Circle is a 2D circle.
//
// The radius is never negative: Circ, Resize and Lerp take it absolute, Scale takes a negative
// factor absolute, and Grow and Shrink clamp at zero, since a circle mirrored about its center
// is the same circle. A negative radius can only be written as a struct literal or decoded from
// JSON, and Contains, DistanceTo and the Intersects methods give no meaningful answer for it;
// Canonical repairs it.
type Circle[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius T        `json:"r"`
}

// Circ is shorthand for Circle{center, radius}, with the radius taken absolute.
func Circ[T Number](center Point[T], radius T) Circle[T] {
	return Circle[T]{center, Abs(radius)}
}

// CircleBetween creates the Circle with the two points at the ends of a diameter: its center
// is their midpoint and its radius half the distance between them. For integer T both are
// rounded like every other result stored into T, so the points can lie off the boundary.
func CircleBetween[T Number](a, b Point[T]) Circle[T] {
	return Circle[T]{a.Midpoint(b), Cast[T](a.Float().DistanceTo(b.Float()) / 2)}
}

// CircleThrough creates the Circle through three points, the circumscribed circle of their
// triangle, and reports false where there is none: the points lie on one line within the
// tolerance, two or all of them being the same point included. The center is measured from
// the first point, so the products stay small far from the origin, and the radius from the
// center as it is stored into T to that point. For integer T both are rounded, so the points
// can lie off the boundary. The radius grows without bound as the points near a line, and
// must fit T.
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

// Anchor returns the point on the circle boundary in the given direction from its center,
// or the center itself for DirectionNone. For integer T a diagonal anchor is rounded like
// Direction.Vector and only approximates the boundary: at a small radius the rounding can
// carry it outside the circle, where Contains rejects it.
func (c Circle[T]) Anchor(direction Direction) Point[T] {
	return c.Center.Add(direction.Vector(c.Radius))
}

// Centroid returns the center of the enclosed area, the Center of the circle.
func (c Circle[T]) Centroid() Point[T] {
	return c.Center
}

// Area returns the circle area (π * radius^2). It is a float64 even for an integer T, since
// the factor π leaves no radius with an area T could express.
func (c Circle[T]) Area() float64 {
	radius := float64(c.Radius)

	return Pi * radius * radius
}

// Perimeter returns the circle circumference (2 * π * radius).
func (c Circle[T]) Perimeter() float64 {
	return 2 * Pi * float64(c.Radius)
}

// Inertia returns the polar second moment of area about the center (π * radius^4 / 2), the
// rotational inertia of the disc at unit density.
func (c Circle[T]) Inertia() float64 {
	radius := float64(c.Radius)

	return Pi * radius * radius * radius * radius / 2
}

// Diameter returns the circle diameter (2 * radius). It stays in T like a sum.
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Bounds returns the axis-aligned bounding box: the square of side Diameter centered on the
// circle.
func (c Circle[T]) Bounds() Box[T] {
	a, b := c.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the smallest circle around the circle, which is the circle itself. It
// is here so every shape answers the same way.
func (c Circle[T]) BoundingCircle() Circle[T] {
	return c
}

// minMax returns the minimum and maximum corner of the circle, the corners of Bounds: the pair
// the intersection tests reject shapes by before examining any edge.
func (c Circle[T]) minMax() (Point[T], Point[T]) {
	return c.Center.AddXY(-c.Radius, -c.Radius), c.Center.AddXY(c.Radius, c.Radius)
}

// magnitude returns the largest absolute coordinate of the circle, that of a corner of its
// Bounds: the size epsilonAt widens the tolerance by for a comparison that reads the circle.
func (c Circle[T]) magnitude() float64 {
	return c.Center.magnitude() + float64(c.Radius)
}

// Translate creates a new Circle translated by the given vector.
func (c Circle[T]) Translate(vector Vector[T]) Circle[T] {
	return Circle[T]{c.Center.Add(vector), c.Radius}
}

// MoveTo creates a new Circle with the same radius and the center set to point.
func (c Circle[T]) MoveTo(point Point[T]) Circle[T] {
	return Circle[T]{point, c.Radius}
}

// AlignTo creates a new Circle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (c Circle[T]) AlignTo(direction Direction, point Point[T]) Circle[T] {
	return c.Translate(point.Subtract(c.Anchor(direction)))
}

// Rotate creates a new Circle turned by the given angle (in radians) about its center, in the
// same sense as Vector.Rotate, which is the circle itself: a circle is symmetric about its
// center, so no angle moves it. It is here so every shape turns the same way.
func (c Circle[T]) Rotate(_ float64) Circle[T] {
	return c
}

// Scale creates a new Circle with radius scaled by the given factor. A negative factor scales
// by its absolute value, since a circle mirrored about its center is the same circle.
func (c Circle[T]) Scale(factor float64) Circle[T] {
	return Circle[T]{c.Center, Abs(Multiply(c.Radius, factor))}
}

// Unscale creates a new Circle with radius scaled by the inverse factor, the inverse of Scale,
// a negative factor by its absolute value like it. Like Divide it panics for a zero factor.
func (c Circle[T]) Unscale(factor float64) Circle[T] {
	return Circle[T]{c.Center, Abs(Divide(c.Radius, factor))}
}

// Resize creates a new Circle with the given radius, taken absolute like Circ.
func (c Circle[T]) Resize(radius T) Circle[T] {
	return Circle[T]{c.Center, Abs(radius)}
}

// Grow creates a new Circle with radius increased by amount, clamped to zero.
func (c Circle[T]) Grow(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius+amount, 0)}
}

// Shrink creates a new Circle with radius decreased by amount, clamped to zero.
func (c Circle[T]) Shrink(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius-amount, 0)}
}

// Lerp creates a new Circle in linear interpolation towards the given circle, moving the center
// and the radius together, and extrapolating outside [0, 1] like Point.Lerp, with the radius
// taken absolute like Circ: an extrapolation past a zero radius grows the circle again.
func (c Circle[T]) Lerp(circle Circle[T], t float64) Circle[T] {
	return Circle[T]{c.Center.Lerp(circle.Center, t), Abs(Lerp(c.Radius, circle.Radius, t))}
}

// Canonical creates a new Circle in the form Circ builds, with the radius taken absolute: a
// well-formed circle is returned as it is. It repairs a negative radius written as a struct
// literal or decoded from JSON before Contains, DistanceTo or the Intersects methods read it.
func (c Circle[T]) Canonical() Circle[T] {
	return Circle[T]{c.Center, Abs(c.Radius)}
}

// DistanceTo returns the distance from the given point to the nearest point of the circle:
// zero exactly where Contains holds, so a point within the tolerance of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the center less the radius.
func (c Circle[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.centerDistanceSquared(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return math.Sqrt(distanceSquared) - float64(c.Radius)
}

// DistanceSquaredTo returns the square of DistanceTo, so every shape offers the same pair. It
// is a float64 even for an integer T and saves no square root, since the distance to a circle
// already needs the root of the distance to its center.
func (c Circle[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := c.DistanceTo(point)

	return distance * distance
}

// DistanceToCircle returns the distance between the nearest points of the two circles: zero
// exactly where IntersectsCircle holds, by the comparison it makes, and otherwise the distance
// between the centers less the two radii.
func (c Circle[T]) DistanceToCircle(circle Circle[T]) float64 {
	if c.IntersectsCircle(circle) {
		return 0
	}

	return math.Sqrt(c.centerDistanceSquared(circle.Center)) - float64(c.Radius) - float64(circle.Radius)
}

// DistanceSquaredToCircle returns the square of DistanceToCircle, so the pair is offered as
// every distance is. It saves no square root, since the distance between two circles already
// needs the root of the distance between their centers.
func (c Circle[T]) DistanceSquaredToCircle(circle Circle[T]) float64 {
	distance := c.DistanceToCircle(circle)

	return distance * distance
}

// Nearest returns the point of the circle nearest to the given point: the point itself
// exactly where Contains holds, and otherwise the point of the boundary toward it from the
// center, on the squared distance Contains compares. For integer T it is rounded once and can
// land off the boundary, where Contains rejects it.
func (c Circle[T]) Nearest(point Point[T]) Point[T] {
	distanceSquared := c.centerDistanceSquared(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return point
	}

	center := c.Center.Float()
	offset := point.Float().Subtract(center).Multiply(float64(c.Radius) / math.Sqrt(distanceSquared))

	return center.Add(offset).Cast[T]()
}

// Support returns the point of the circle farthest along the given direction: the center moved
// by the radius that way, whatever the length of the direction, and along X for a direction of
// no length. It is a float64 point even for an integer T, since it is no lattice point in
// general.
func (c Circle[T]) Support(direction Vector[float64]) Point[float64] {
	return c.Center.Float().Add(direction.Normalize().Multiply(float64(c.Radius)))
}

// Tangents returns the points where the lines through the given point touch the circle: two
// for a point outside, the point itself for one on the boundary within the tolerance, by the
// comparison Contains makes, and none for a point inside. The second point is the mirror of
// the first across the line from the center to the point, the order IntersectionCircle gives
// its own. For integer T the points are rounded like every other result stored into T, and two
// rounding to the same point give it once; a circle of zero radius is touched at its center.
//
// The points are a Pair, which allocates nothing.
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

// OuterTangents returns the segments of the lines that touch both circles from the same side,
// the straight runs of a belt around two pulleys, each from its point on this circle to its
// point on the given one: two for circles apart, touching from outside or overlapping, and
// none where one lies within the other, touching it from inside included, or the centers are
// within the tolerance of each other. The second is the mirror of the first across the line of
// centers, the order Tangents gives its own, and a circle of zero radius is touched at its
// center, so against one the points on this circle are the Tangents of that center. Two
// circles of zero radius are joined by one segment, and for integer T the points are rounded
// like every other result stored into T, two segments rounding to the same one given once.
//
// The segments are a Pair, which allocates nothing.
func (c Circle[T]) OuterTangents(circle Circle[T]) Pair[Segment[T]] {
	return c.tangents(circle, float64(circle.Radius))
}

// InnerTangents returns the segments of the lines that touch the circles from opposite sides
// and cross between them, the straight runs of a crossed belt, each from its point on this
// circle to its point on the given one: two for circles apart, and none for circles that
// intersect, by the comparison IntersectsCircle makes, so touching circles have none. The
// second is the mirror of the first across the line of centers. Two circles of zero radius are
// joined by one segment, and for integer T the points are rounded like every other result
// stored into T, two segments rounding to the same one given once.
//
// The segments are a Pair, which allocates nothing.
func (c Circle[T]) InnerTangents(circle Circle[T]) Pair[Segment[T]] {
	return c.tangents(circle, -float64(circle.Radius))
}

// Contains reports whether the given point lies within the circle, boundary included within
// the tolerance: a float point a rounding error outside the radius, such as an Anchor, is
// still contained.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.centerDistanceSquared(point), point.magnitude())
}

// ContainsCircle reports whether the given circle lies within this one: its far point, the
// point of it farthest from this center, is contained within the tolerance by the comparison
// Contains makes, so a circle touching the boundary from inside is enclosed.
func (c Circle[T]) ContainsCircle(circle Circle[T]) bool {
	far := math.Sqrt(c.centerDistanceSquared(circle.Center)) + float64(circle.Radius)

	return c.containsSquared(far*far, circle.magnitude())
}

// ContainsCapsule reports whether the given capsule lies within the circle: the circles of its
// radius about both ends of its segment do, as ContainsCircle decides it, and a circle holds every point between points it contains.
func (c Circle[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return c.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && c.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the given sector lies within the circle: its center, the two
// ends of its arc and the point of it farthest from this center are contained, within the
// tolerance, and a circle holds every point between points it contains.
func (c Circle[T]) ContainsSector(sector Sector[T]) bool {
	outer, inner := c.Float(), sector.Float()
	for _, corner := range inner.corners() {
		if !outer.Contains(corner) {
			return false
		}
	}

	return outer.Contains(inner.Support(inner.Center.Subtract(outer.Center)))
}

// ContainsSegment reports whether the segment lies within the circle: both endpoints are
// contained, within the tolerance, and a circle holds every point between two it contains.
func (c Circle[T]) ContainsSegment(segment Segment[T]) bool {
	return c.Contains(segment.Start) && c.Contains(segment.End)
}

// ContainsPolygon reports whether the polygon lies within the circle: every vertex is
// contained, within the tolerance, and a circle holds every point between points it contains.
// An empty polygon is enclosed by nothing.
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

// ContainsRectangle reports whether the rectangle lies within the circle: every corner is
// contained, within the tolerance, whatever the rectangle's angle.
func (c Circle[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies within the circle: every
// vertex is contained, within the tolerance. An empty polygon is enclosed by nothing.
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

// ContainsBox reports whether the box lies within the circle: every corner is contained, within
// the tolerance.
func (c Circle[T]) ContainsBox(box Box[T]) bool {
	for _, corner := range box.corners() {
		if !c.Contains(corner) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the circles overlap: the distance between the centers is at
// most the sum of the radii. Touching circles intersect, within the tolerance, by the same
// comparison Contains makes, on the squared distance.
func (c Circle[T]) IntersectsCircle(circle Circle[T]) bool {
	return lessOrEqualSquared(c.centerDistanceSquared(circle.Center), float64(c.Radius)+float64(circle.Radius), c.epsilonWith(circle.magnitude()))
}

// IntersectsCapsule reports whether the circle and the capsule share a point: the segment of
// the capsule passes within the sum of the two radii of the center. Touching shapes intersect,
// within the tolerance, by the same comparison IntersectsCircle makes on the squared distance.
func (c Circle[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return lessOrEqualSquared(capsule.Segment.DistanceSquaredTo(c.Center), float64(c.Radius)+float64(capsule.Radius), c.epsilonWith(capsule.magnitude()))
}

// IntersectsSector reports whether the circle and the sector share a point: the center lies
// within the sector, or its arc or one of its radii passes within the radius. Touching shapes
// intersect, within the tolerance, by the same comparison Contains makes on the squared
// distance the sector's DistanceSquaredTo measures.
func (c Circle[T]) IntersectsSector(sector Sector[T]) bool {
	return c.containsSquared(sector.DistanceSquaredTo(c.Center), sector.magnitude())
}

// IntersectsSegment reports whether the circle and the segment share a point: the point of the
// segment closest to the center lies within the radius. Touching shapes intersect, within
// the tolerance, by the same comparison Contains makes, on the squared distance.
func (c Circle[T]) IntersectsSegment(segment Segment[T]) bool {
	return c.containsSquared(segment.DistanceSquaredTo(c.Center), segment.magnitude())
}

// IntersectsRay reports whether the circle and the ray share a point, as IntersectsSegment
// decides it on the reach of the ray past the circle: the point of the ray closest to the center
// lies within the radius. Touching shapes intersect, within the tolerance.
func (c Circle[T]) IntersectsRay(ray Ray[T]) bool {
	return c.IntersectsSegment(ray.reach(c.minMax()))
}

// IntersectsLine reports whether the circle and the line share a point, as IntersectsSegment
// decides it on the span of the line past the circle: the point of the line closest to the
// center lies within the radius. Touching shapes intersect, within the tolerance.
func (c Circle[T]) IntersectsLine(line Line[T]) bool {
	return c.IntersectsSegment(line.span(c.minMax()))
}

// IntersectsPolygon reports whether the circle and the polygon share a point: the center lies
// within the polygon, or an edge passes within the radius. Touching shapes intersect, within
// the tolerance, by the same comparison Contains makes on the squared distance the polygon's
// DistanceSquaredTo measures. A circle whose Bounds lie outside the polygon is rejected before
// any edge is examined, and an empty polygon intersects nothing.
func (c Circle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && polygon.walk(c.Center).reaches(c)
}

// IntersectsRectangle reports whether the circle and the rectangle overlap: the center lies
// within the rectangle, or an edge passes within the radius. Touching shapes intersect, within
// the tolerance, by the same comparison Contains makes on the squared distance the rectangle's
// DistanceSquaredTo measures.
func (c Circle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.walk(c.Center).reaches(c)
}

// IntersectsRegularPolygon reports whether the circle and the regular polygon share a point:
// the center lies within the polygon, or an edge passes within the radius, by the same
// comparison Contains makes on the squared distance the polygon's DistanceSquaredTo measures.
// A circle whose Bounds lie outside the polygon's is rejected before any edge is examined, and
// an empty polygon intersects nothing.
func (c Circle[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && polygon.walk(c.Center).reaches(c)
}

// IntersectsBox reports whether the circle and the box overlap: the center lies within the box,
// or an edge passes within the radius. Touching shapes intersect, within the tolerance, by the
// same comparison Contains makes on the squared distance the box's DistanceSquaredTo measures,
// the gap beyond it on the two axes with no edge to walk.
func (c Circle[T]) IntersectsBox(box Box[T]) bool {
	return c.containsSquared(box.DistanceSquaredTo(c.Center), box.magnitude())
}

// IntersectionCircle returns the points where the circles cross: two for overlapping circles, one
// for tangent ones, within the tolerance like IntersectsCircle, and none for circles apart or
// nested.
// Coincident circles share every point and also return none, and circles with centers within
// the tolerance of each other count as coincident. The second point is the mirror of the first
// across the line of centers. Tangency is judged on the sum and the difference of the radii by
// the same comparison IntersectsCircle makes, so a tangent a rounding error outside its reach gives
// one point rather than the same point twice; that point is placed halfway between the two
// boundaries where they meet, since the crossing formula amplifies the tolerance there. For
// integer T the points are rounded like every other result stored into T, and two crossings
// rounding to the same point give it once.
//
// The points are a Pair, which allocates nothing.
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

// IntersectionSegment returns the points where the segment crosses the circle boundary, from the
// segment's Start to its End: two where it passes through, one where it is tangent or ends inside,
// within the tolerance like IntersectsSegment, and none where it misses or lies entirely inside. A
// segment inside crosses no boundary, so it returns none while IntersectsSegment still reports it.
// An endpoint within the tolerance of the boundary is one of the crossings, judged by the same
// comparison IntersectsSegment makes, so a shallow touch is not lost to the fraction along the
// chord and the two agree to the last bit: Start is where the segment enters and End where it
// leaves, unless the chord lies wholly past that endpoint, and two such endpoints are the two
// crossings, so a segment and its Reverse cross at the same points; where the chord is a tangent
// the endpoint replaces it. A tangent is placed at the foot on the segment, a fraction along it
// like every crossing, so it lies within the tolerance of the circle rather than halfway between
// the two boundaries, and a float32 tangent far from the origin can round just past the tolerance.
// For integer T the points are rounded like every other result stored into T.
//
// The points are a Pair, which allocates nothing.
func (c Circle[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	fractions, n := c.fractions(segment)

	return segment.pointsAt(fractions[:n]...)
}

// IntersectionRay returns the points where the ray crosses the circle boundary, from its Origin
// on, as IntersectionSegment finds them on the reach of the ray past the circle, so the two
// agree with IntersectsRay to the last bit: two where it passes through, one where it is tangent
// or starts inside, and none where it misses.
func (c Circle[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return c.IntersectionSegment(ray.reach(c.minMax()))
}

// IntersectionLine returns the points where the line crosses the circle boundary, in order along
// its Direction, as IntersectionSegment finds them on the span of the line past the circle, so
// the two agree with IntersectsLine to the last bit: two where it passes through, one where it
// is tangent, and none where it misses.
func (c Circle[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return c.IntersectionSegment(line.span(c.minMax()))
}

// fractions returns the fractions along the segment at which it crosses the circle boundary, in
// increasing order, and how many of the two there are: the ends of the chord, an endpoint on the
// boundary within the tolerance standing in for the end it is by snapToEndpoint, or the one
// fraction of a tangent. They are not yet held to the segment, which pointsAt does, so
// IntersectionSegment and Segment.HitCircle place their points from the same fractions.
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

// centerDistanceSquared returns the squared distance from the center to the point, in
// float64: the value every test of the circle compares against its radius, so Contains,
// DistanceTo, IntersectsCircle, IntersectionCircle and IntersectionSegment agree to the last bit.
func (c Circle[T]) centerDistanceSquared(point Point[T]) float64 {
	return c.Center.Float().DistanceSquaredTo(point.Float())
}

// containsSquared reports whether a point at the given squared distance from the center lies within
// the circle, boundary included within the tolerance: lessOrEqualSquared on the radius, the
// comparison Contains, DistanceTo and every Intersects method of the circle make.
func (c Circle[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, float64(c.Radius), c.epsilonWith(magnitude))
}

// touchesSquared reports whether a point at the given squared distance from the center lies on the
// boundary within the tolerance: equalSquared on the radius, the endpoint test of
// IntersectionSegment.
func (c Circle[T]) touchesSquared(distanceSquared, magnitude float64) bool {
	return equalSquared(distanceSquared, float64(c.Radius), c.epsilonWith(magnitude))
}

// epsilonWith returns the tolerance of a comparison between the circle and geometry whose
// largest absolute coordinate is the given magnitude, as epsilonAt gives it for the larger of
// the two, so every test of a pair reads the same one.
func (c Circle[T]) epsilonWith(magnitude float64) float64 {
	return epsilonAt[T](max(c.magnitude(), magnitude))
}

// tangent returns how far along the line of centers, from this center toward the other at the
// given distance, the boundaries of two tangent circles meet. Each boundary crosses that line
// at a known place, this one at ±Radius and the other at the distance ± its radius, and the
// point is taken halfway between the two crossings that meet: externally, from outside; or
// internally, on the far side of the smaller circle. The exact crossing formula is ill-conditioned
// at a tangent, moving twice as far as the center distance it is given, so the midpoint keeps the
// point within half the tolerance of both boundaries.
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

// tangents returns the segments of the lines that touch this circle and the given one,
// whose radius is passed signed: as it is for the outer tangents, and negative for the inner
// ones, which touch it on the far side. The line lies at a right angle to one direction from
// both centers, a radius from each, so along the line of centers the two radii differ by the
// distance between the centers times the cosine of that direction: the ratio both points are
// placed by, each an end of a chord of its own circle. There is no such direction where the
// distance is at most the difference of the radii, within the tolerance, which is the test
// IntersectsCircle makes on their sum for the inner tangents.
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

// chordPoints returns the two points chordEnds returns, the points IntersectionCircle and
// Tangents place, held once where both are stored into T as the same point.
func (c Circle[T]) chordPoints(direction Vector[float64], along float64) Pair[Point[T]] {
	first, second := c.chordEnds(direction, along)
	if second.Equal(first) {
		return Pair[Point[T]]{[2]Point[T]{first}, 1}
	}

	return Pair[Point[T]]{[2]Point[T]{first, second}, 2}
}

// chordEnds returns the two ends of the chord that crosses the given direction from the center
// at a right angle, the given length along it, negative for a chord behind the center: the
// first a quarter turn from the direction by Normal and the second its mirror.
func (c Circle[T]) chordEnds(direction Vector[float64], along float64) (Point[T], Point[T]) {
	radius := float64(c.Radius)
	middle := c.Center.Float().Add(direction.Resize(along))
	normal := direction.Normal().Resize(math.Sqrt(max(float64(radius*radius)-float64(along*along), 0)))

	return middle.Add(normal).Cast[T](), middle.Add(normal.Negate()).Cast[T]()
}

// Equal checks for equal center and radius with given circle.
func (c Circle[T]) Equal(circle Circle[T]) bool {
	return c.Center.Equal(circle.Center) && Equal(c.Radius, circle.Radius)
}

// IsZero checks if center point and radius are zero.
func (c Circle[T]) IsZero() bool {
	return c.Equal(Circle[T]{})
}

// Ellipse converts the circle into an Ellipse of equal semi-axes, with no angle, and the one
// Ellipse.Circle gives back unchanged: the circle has no Transform of its own, since an affine
// matrix takes it to an ellipse, so a transform of a circle goes Ellipse().Transform(matrix),
// which says plainly that the circle is left behind.
func (c Circle[T]) Ellipse() Ellipse[T] {
	return Ellipse[T]{c.Center, SzU(c.Radius), 0}
}

// RegularPolygon converts the circle into the RegularPolygon of n vertices inscribed in it, with
// the given orientation, as Ellipse.RegularPolygon does: every vertex lies on the boundary, and the
// orientation places the first of them by RegularPolygonOrientationPhase, so OrientationPointyTop
// puts a vertex at the top and OrientationFlatTop the midpoint of an edge. It is the outline a
// circle does not have, so its Vertices and Edges are what draws or walks one.
//
// Like RegularPolygonOrientationPhase it panics for an orientation that is neither
// OrientationFlatTop nor OrientationPointyTop.
func (c Circle[T]) RegularPolygon(n int, orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{c.Center, SzU(c.Radius), n, 0, RegularPolygonOrientationPhase(n, orientation)}
}

// Arc returns the part of the circle boundary that starts at the given angle and runs through
// the given sweep, both in radians: the Arc of this center and radius. A sweep of a full turn is
// the whole boundary, the line of the circle without its inside, which contains the points on
// the boundary alone and is measured to from within.
func (c Circle[T]) Arc(angle, sweep float64) Arc[T] {
	return Arc[T]{c.Center, c.Radius, angle, sweep}
}

// Sector returns the slice of the circle between the radius at the given angle and the one the
// given sweep away, both in radians: the Sector of this center and radius. A sweep of a full
// turn is the whole disc, which the circle itself answers for faster.
func (c Circle[T]) Sector(angle, sweep float64) Sector[T] {
	return Sector[T]{c.Center, c.Radius, angle, sweep}
}

// Polygon converts the circle to a Polygon within the given tolerance of it, the outline that
// draws or walks one: equal chords from the point at the angle zero, as few as keep the circle
// within the tolerance of each and never fewer than three, with its vertices on the boundary
// and wound as the corners of a Rectangle are. A negative tolerance is taken absolute, and the
// circle has at most a thousand chords, which a tolerance of zero gives. A circle of no radius
// gives its center alone. For integer T the vertices are rounded and one rounding onto the
// vertex before it is left out. RegularPolygon is the outline of a chosen vertex count.
func (c Circle[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{c.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of the polygon Polygon returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them.
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

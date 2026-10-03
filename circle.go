package geom

import (
	"fmt"
	"math"
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

// Diameter returns the circle diameter (2 * radius). It stays in T like a sum, since doubling
// has no intermediate to overflow: only a diameter beyond the range of a narrow integer T is
// lost, as with every other result outside it.
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Bounds returns the axis-aligned bounding box: the square of side Diameter centered on the
// circle, clamped into the range of a narrow integer T where it reaches past it.
func (c Circle[T]) Bounds() Box[T] {
	a, b := c.minMax()

	return Box[T]{a, b}
}

// minMax returns the minimum and maximum corner of the circle, the corners of Bounds: the pair
// the intersection tests reject shapes by before examining any edge. The corners are taken in
// float64 and clamped into the range of an integer T, so a narrow T whose circle reaches past
// its range bounds it at the end of the range, which every shape it is tested against lies
// within, rather than at a corner wrapped to the other side.
func (c Circle[T]) minMax() (Point[T], Point[T]) {
	center, radius := c.Center.Float(), float64(c.Radius)
	a := Point[T]{castClamped[T](center.X - radius), castClamped[T](center.Y - radius)}
	b := Point[T]{castClamped[T](center.X + radius), castClamped[T](center.Y + radius)}

	return a, b
}

// magnitude returns the largest absolute coordinate of the circle, that of a corner of its
// Bounds, summed in float64 so a narrow integer T cannot overflow: the size epsilonAt widens
// the tolerance by for a comparison that reads the circle.
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

// Canonical creates a new Circle in the form Circ builds, with the radius taken absolute: a
// well-formed circle is returned as it is. It repairs a negative radius written as a struct
// literal or decoded from JSON before Contains, DistanceTo or the Intersects methods read it.
func (c Circle[T]) Canonical() Circle[T] {
	return Circle[T]{c.Center, Abs(c.Radius)}
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

// Rotate creates a new Circle turned by the given angle (in radians) about its center, in the
// same sense as Vector.Rotate, which is the circle itself: a circle is symmetric about its
// center, so no angle moves it. It is here so every shape turns the same way.
func (c Circle[T]) Rotate(_ float64) Circle[T] {
	return c
}

// AlignTo creates a new Circle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (c Circle[T]) AlignTo(direction Direction, point Point[T]) Circle[T] {
	return c.Translate(point.Subtract(c.Anchor(direction)))
}

// Contains reports whether the given point lies within the circle, boundary included within
// the tolerance: a float point a rounding error outside the radius, such as an Anchor, is
// still contained.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.centerDistanceSquared(point), point.magnitude())
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

// EnclosesCircle reports whether the given circle lies within this one: its far point, the
// point of it farthest from this center, is contained within the tolerance by the comparison
// Contains makes, so a circle touching the boundary from inside is enclosed.
func (c Circle[T]) EnclosesCircle(circle Circle[T]) bool {
	far := math.Sqrt(c.centerDistanceSquared(circle.Center)) + float64(circle.Radius)

	return c.containsSquared(far*far, circle.magnitude())
}

// EnclosesSegment reports whether the segment lies within the circle: both endpoints are
// contained, within the tolerance, and a circle holds every point between two it contains.
func (c Circle[T]) EnclosesSegment(segment Segment[T]) bool {
	return c.Contains(segment.Start) && c.Contains(segment.End)
}

// EnclosesPolygon reports whether the polygon lies within the circle: every vertex is
// contained, within the tolerance, and a circle holds every point between points it contains.
// An empty polygon is enclosed by nothing.
func (c Circle[T]) EnclosesPolygon(polygon Polygon[T]) bool {
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

// EnclosesRectangle reports whether the rectangle lies within the circle: every corner is
// contained, within the tolerance, whatever the rectangle's angle.
func (c Circle[T]) EnclosesRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// EnclosesRegularPolygon reports whether the regular polygon lies within the circle: every
// vertex is contained, within the tolerance. An empty polygon is enclosed by nothing.
func (c Circle[T]) EnclosesRegularPolygon(polygon RegularPolygon[T]) bool {
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

// EnclosesBox reports whether the box lies within the circle: every corner is contained, within
// the tolerance.
func (c Circle[T]) EnclosesBox(box Box[T]) bool {
	for _, corner := range box.corners() {
		if !c.Contains(corner) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the circles overlap: the distance between the centers is at
// most the sum of the radii. Touching circles intersect, within the tolerance, by the same comparison
// Contains makes, on the squared distance. The radii are summed in float64, so a narrow integer
// T cannot overflow the threshold.
func (c Circle[T]) IntersectsCircle(circle Circle[T]) bool {
	return lessOrEqualSquared(c.centerDistanceSquared(circle.Center), float64(c.Radius)+float64(circle.Radius), c.epsilonWith(circle.magnitude()))
}

// IntersectionCircle returns the points where the circles cross: two for overlapping circles, one
// for tangent ones, within the tolerance like Intersects, and none for circles apart or nested.
// Coincident circles share every point and also return none, and circles with centers within
// the tolerance of each other count as coincident. The second point is the mirror of the first
// across the line of centers. Tangency is judged on the sum and the difference of the radii by
// the same comparison Intersects makes, so a tangent a rounding error outside its reach gives
// one point rather than the same point twice; that point is placed halfway between the two
// boundaries where they meet, since the crossing formula amplifies the tolerance there. For
// integer T the points are rounded like every other result stored into T, and two crossings
// rounding to the same point give it once.
func (c Circle[T]) IntersectionCircle(circle Circle[T]) []Point[T] {
	return c.AppendIntersectionCircle(nil, circle)
}

// AppendIntersectionCircle appends the points IntersectionCircle returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room.
func (c Circle[T]) AppendIntersectionCircle(dst []Point[T], circle Circle[T]) []Point[T] {
	center := c.Center.Float()
	direction := circle.Center.Float().Subtract(center)
	distanceSquared := c.centerDistanceSquared(circle.Center)
	r1, r2 := float64(c.Radius), float64(circle.Radius)
	outer, inner := r1+r2, math.Abs(r1-r2)
	epsilon := c.epsilonWith(circle.magnitude())

	if lessOrEqualSquared(distanceSquared, 0, epsilon) || !lessOrEqualSquared(distanceSquared, outer, epsilon) || !greaterOrEqualSquared(distanceSquared, inner, epsilon) {
		return dst
	}

	distance := math.Sqrt(distanceSquared)

	if external, internal := equalSquared(distanceSquared, outer, epsilon), equalSquared(distanceSquared, inner, epsilon); external || internal {
		point := center.Add(direction.Resize(c.tangent(circle, distance, external)))

		return append(dst, point.Cast[T]())
	}

	along := (float64(r1*r1) - float64(r2*r2) + distanceSquared) / (2 * distance)
	middle := center.Add(direction.Resize(along))
	normal := direction.Normal().Resize(math.Sqrt(max(float64(r1*r1)-float64(along*along), 0)))
	first, second := middle.Add(normal).Cast[T](), middle.Add(normal.Negate()).Cast[T]()
	if second.Equal(first) {
		return append(dst, first)
	}

	return append(dst, first, second)
}

// IntersectsSegment reports whether the circle and the segment share a point: the point of the
// segment closest to the center lies within the radius. Touching shapes intersect, within
// the tolerance, by the same comparison Contains makes, on the squared distance.
func (c Circle[T]) IntersectsSegment(segment Segment[T]) bool {
	return c.containsSquared(segment.DistanceSquaredTo(c.Center), segment.magnitude())
}

// IntersectionSegment returns the points where the segment crosses the circle boundary, from the
// segment's Start to its End: two where it passes through, one where it is tangent or ends
// inside, within the tolerance like IntersectsSegment, and none where it misses or lies entirely
// inside. A segment inside crosses no boundary, so it returns none while IntersectsSegment still
// reports it. An endpoint within the tolerance of the boundary is the crossing nearest to it, judged
// by the same comparison IntersectsSegment makes, so a shallow touch is not lost to the fraction
// along the chord and the two agree to the last bit; where the chord is a tangent the endpoint
// replaces it. A tangent is placed at the foot on the segment, a fraction along it like every
// crossing, so it lies within the tolerance of the circle rather than halfway between the two
// boundaries, and a float32 tangent far from the origin can round just past the tolerance. For
// integer T the points are rounded like every other result stored into T.
func (c Circle[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return c.AppendIntersectionSegment(nil, segment)
}

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room. The points already
// in dst are kept as they are: a crossing equal to one of them is still appended, and only the
// appended ones are ordered from Start. A nil dst is allocated on the first point with room for
// two, and Segment.ClipCircle passes a buffer of its own.
func (c Circle[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	entry, exit, ok := segment.chord(c)
	start := c.touchesSquared(c.centerDistanceSquared(segment.Start), segment.magnitude())
	end := c.touchesSquared(c.centerDistanceSquared(segment.End), segment.magnitude())

	switch {
	case ok && entry < exit:
		if start {
			entry, exit = segment.snapToEndpoint(entry, exit, 0)
		}
		if end {
			entry, exit = segment.snapToEndpoint(entry, exit, 1)
		}

		return segment.appendPointsAt(dst, entry, exit)
	case start && end && segment.Vector().hasDirection():
		return segment.appendPointsAt(dst, 0, 1)
	case start:
		return segment.appendPointsAt(dst, 0)
	case end:
		return segment.appendPointsAt(dst, 1)
	case ok:
		return segment.appendPointsAt(dst, entry)
	default:
		return dst
	}
}

// IntersectsRay reports whether the circle and the ray share a point, as IntersectsSegment
// decides it on the reach of the ray past the circle: the point of the ray closest to the center
// lies within the radius. Touching shapes intersect, within the tolerance.
func (c Circle[T]) IntersectsRay(ray Ray[T]) bool {
	return c.IntersectsSegment(ray.reach(c.minMax()))
}

// IntersectionRay returns the points where the ray crosses the circle boundary, from its Origin
// on, as IntersectionSegment finds them on the reach of the ray past the circle, so the two
// agree with IntersectsRay to the last bit: two where it passes through, one where it is tangent
// or starts inside, and none where it misses.
func (c Circle[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return c.AppendIntersectionRay(nil, ray)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice, as AppendIntersectionSegment does on the reach of the ray past the circle.
func (c Circle[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	return c.AppendIntersectionSegment(dst, ray.reach(c.minMax()))
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

// RegularPolygon converts the circle into the RegularPolygon of n vertices inscribed in it,
// with the given orientation, as Ellipse.RegularPolygon does: every vertex lies on
// the boundary, and the orientation places the first of them by RegularPolygonOrientationPhase,
// so OrientationPointyTop puts a vertex at the top and OrientationFlatTop the midpoint of an edge. It is the outline a
// circle does not have, so its Vertices and Edges are what draws or walks one.
//
// Like RegularPolygonOrientationPhase it panics for an orientation that is neither
// OrientationFlatTop nor OrientationPointyTop.
func (c Circle[T]) RegularPolygon(n int, orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{c.Center, SzU(c.Radius), n, 0, RegularPolygonOrientationPhase(n, orientation)}
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
	return fmt.Sprintf("Circ(%s;%s)", c.Center.String(), String(c.Radius))
}

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

// Diameter returns the circle diameter (2 * radius). It stays in T like a sum, since doubling
// has no intermediate to overflow: only a diameter beyond the range of a narrow integer T is
// lost, as with every other result outside it.
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Bounds returns the axis-aligned bounding rectangle: the square of side Diameter
// centered on the circle.
func (c Circle[T]) Bounds() Rectangle[T] {
	side := c.Diameter()

	return Rectangle[T]{c.Center, Size[T]{side, side}, 0}
}

// Anchor returns the point on the circle boundary in the given direction from its center,
// or the center itself for DirectionNone. For integer T a diagonal anchor is rounded like
// Direction.Vector and only approximates the boundary: at a small radius it can land outside
// the circle, so Circ(Pt(0, 0), 1).Anchor(BottomRight) is (1,1), which Contains rejects.
func (c Circle[T]) Anchor(direction Direction) Point[T] {
	return c.Center.Add(direction.Vector(c.Radius))
}

// minMax returns the minimum and maximum corner of the circle, the corners of Bounds, exact
// for an integer T: the pair the intersection tests reject shapes by before examining any
// edge, without placing a rectangle.
func (c Circle[T]) minMax() (Point[T], Point[T]) {
	return Point[T]{c.Center.X - c.Radius, c.Center.Y - c.Radius}, Point[T]{c.Center.X + c.Radius, c.Center.Y + c.Radius}
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

// AlignTo creates a new Circle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (c Circle[T]) AlignTo(direction Direction, point Point[T]) Circle[T] {
	return c.Translate(point.Subtract(c.Anchor(direction)))
}

// Contains reports whether the given point lies within the circle, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains: a float point a rounding
// error outside the radius, such as an Anchor, is still contained.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.centerDistanceSquared(point))
}

// DistanceTo returns the distance from the given point to the nearest point of the circle:
// zero exactly where Contains holds, so a point within Epsilon of T of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the center less the radius.
func (c Circle[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.centerDistanceSquared(point)
	if c.containsSquared(distanceSquared) {
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

// IntersectsCircle reports whether the circles overlap: the center of one lies within the sum of the
// radii of the other. Touching circles intersect, within Epsilon of T, by the same comparison
// Contains makes, on the squared distance. The radii are summed in float64, so a narrow integer
// T cannot overflow the threshold.
func (c Circle[T]) IntersectsCircle(circle Circle[T]) bool {
	return lessOrEqualSquared[T](c.centerDistanceSquared(circle.Center), float64(c.Radius)+float64(circle.Radius))
}

// IntersectionCircle returns the points where the circles cross: two for overlapping circles, one
// for tangent ones, within Epsilon of T like Intersects, and none for circles apart or nested.
// Coincident circles share every point and also return none, and circles with centers within
// Epsilon of T of each other count as coincident. The second point is the mirror of the first
// across the line of centers. Tangency is judged on the sum and the difference of the radii by
// the same comparison Intersects makes, so a tangent a rounding error outside its reach gives
// one point rather than the same point twice; that point is placed halfway between the two
// boundaries where they meet, since the crossing formula amplifies the tolerance there. For
// integer T the points are rounded like every other result stored into T.
func (c Circle[T]) IntersectionCircle(circle Circle[T]) []Point[T] {
	direction := circle.Center.Subtract(c.Center).Float()
	distanceSquared := c.centerDistanceSquared(circle.Center)
	r1, r2 := float64(c.Radius), float64(circle.Radius)
	outer, inner := r1+r2, math.Abs(r1-r2)

	if lessOrEqualSquared[T](distanceSquared, 0) || !lessOrEqualSquared[T](distanceSquared, outer) || !greaterOrEqualSquared[T](distanceSquared, inner) {
		return nil
	}

	distance := math.Sqrt(distanceSquared)

	if external, internal := equalSquared[T](distanceSquared, outer), equalSquared[T](distanceSquared, inner); external || internal {
		point := c.Center.Float().Add(direction.Resize(c.tangent(circle, distance, external)))

		return []Point[T]{{Cast[T](point.X), Cast[T](point.Y)}}
	}

	along := (r1*r1 - r2*r2 + distanceSquared) / (2 * distance)
	middle := c.Center.Float().Add(direction.Resize(along))
	normal := direction.Normal().Resize(math.Sqrt(max(r1*r1-along*along, 0)))
	first, second := middle.Add(normal), middle.Add(normal.Negate())

	return []Point[T]{{Cast[T](first.X), Cast[T](first.Y)}, {Cast[T](second.X), Cast[T](second.Y)}}
}

// IntersectsLine reports whether the circle and the segment share a point: the point of the
// segment closest to the center lies within the radius. Touching shapes intersect, within
// Epsilon of T, by the same comparison Contains makes, on the squared distance.
func (c Circle[T]) IntersectsLine(line Line[T]) bool {
	return c.containsSquared(line.DistanceSquaredTo(c.Center))
}

// IntersectionLine returns the points where the segment crosses the circle boundary, from the
// segment's Start to its End: two where it passes through, one where it is tangent or ends
// inside, within Epsilon of T like IntersectsLine, and none where it misses or lies entirely
// inside. A segment inside crosses no boundary, so it returns none while IntersectsLine still
// reports it. An endpoint within Epsilon of the boundary is the crossing nearest to it, judged
// by the same comparison IntersectsLine makes, so a shallow touch is not lost to the fraction
// along the chord and the two agree to the last bit; where the chord is a tangent the endpoint
// replaces it. For integer T the points are rounded like every other result stored into T.
func (c Circle[T]) IntersectionLine(line Line[T]) []Point[T] {
	entry, exit, ok := line.chord(c)
	start := c.touchesSquared(c.centerDistanceSquared(line.Start))
	end := c.touchesSquared(c.centerDistanceSquared(line.End))

	switch {
	case ok && entry < exit:
		if start {
			entry, exit = snap(entry, exit, 0)
		}
		if end {
			entry, exit = snap(entry, exit, 1)
		}

		return line.pointsAt(entry, exit)
	case start && end && line.Vector().hasDirection():
		return line.pointsAt(0, 1)
	case start:
		return line.pointsAt(0)
	case end:
		return line.pointsAt(1)
	case ok:
		return line.pointsAt(entry)
	default:
		return nil
	}
}

// IntersectsPolygon reports whether the circle and the polygon share a point: the center lies
// within the polygon, or an edge passes within the radius. Touching shapes intersect, within
// Epsilon of T, by the same comparison Contains makes on the squared distance the polygon's
// DistanceSquaredTo measures. A circle whose Bounds lie outside the polygon is rejected before
// any edge is examined, and an empty polygon intersects nothing.
func (c Circle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if polygon.Empty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && c.containsSquared(polygon.DistanceSquaredTo(c.Center))
}

// IntersectsRectangle reports whether the circle and the rectangle overlap: the center lies
// within the rectangle, or an edge passes within the radius. Touching shapes intersect, within
// Epsilon of T, by the same comparison Contains makes on the squared distance the rectangle's
// DistanceSquaredTo measures.
func (c Circle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return c.containsSquared(rectangle.DistanceSquaredTo(c.Center))
}

// IntersectsRegularPolygon reports whether the circle and the regular polygon share a point:
// the center lies within the polygon, or an edge passes within the radius, by the same
// comparison Contains makes on the squared distance the polygon's DistanceSquaredTo measures.
// A circle whose Bounds lie outside the polygon's is rejected before any edge is examined, and
// an empty polygon intersects nothing.
func (c Circle[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.Empty() {
		return false
	}

	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	return overlaps(a1, b1, a2, b2) && c.containsSquared(polygon.DistanceSquaredTo(c.Center))
}

// centerDistanceSquared returns the squared distance from the center to the point, in
// float64: the value every test of the circle compares against its radius, so Contains,
// DistanceTo, IntersectsCircle, IntersectionCircle and IntersectionLine agree to the last bit.
func (c Circle[T]) centerDistanceSquared(point Point[T]) float64 {
	return c.Center.Float().DistanceSquaredTo(point.Float())
}

// containsSquared reports whether a point at the given squared distance from the center lies within
// the circle, boundary included within Epsilon of T: lessOrEqualSquared on the radius, the
// comparison Contains, DistanceTo and every Intersects method of the circle make.
func (c Circle[T]) containsSquared(distanceSquared float64) bool {
	return lessOrEqualSquared[T](distanceSquared, float64(c.Radius))
}

// touchesSquared reports whether a point at the given squared distance from the center lies on the
// boundary within Epsilon of T: equalSquared on the radius, the endpoint test of
// IntersectionLine.
func (c Circle[T]) touchesSquared(distanceSquared float64) bool {
	return equalSquared[T](distanceSquared, float64(c.Radius))
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

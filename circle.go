package geom

import (
	"fmt"
	"math"
)

// Circle is a 2D circle.
type Circle[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius T        `json:"r"`
}

// Circ is shorthand for Circle{center, radius}.
func Circ[T Number](center Point[T], radius T) Circle[T] {
	return Circle[T]{center, radius}
}

// Area returns the circle area (π * radius^2).
func (c Circle[T]) Area() float64 {
	radius := float64(c.Radius)

	return Pi * radius * radius
}

// Circumference returns the circle circumference (2 * π * radius).
func (c Circle[T]) Circumference() float64 {
	return 2 * Pi * float64(c.Radius)
}

// Diameter returns the circle diameter (2 * radius). It stays in T like a sum, since doubling
// has no intermediate to overflow: only a diameter beyond the range of a narrow integer T is
// lost, as with every other result outside it.
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Bounds returns the axis-aligned bounding rectangle: the square of side Diameter
// centered on the circle. A negative radius contains nothing and bounds the zero-size
// rectangle at the center, never a rectangle with a negative size.
func (c Circle[T]) Bounds() Rectangle[T] {
	side := max(c.Diameter(), 0)

	return Rectangle[T]{c.Center, Size[T]{side, side}}
}

// Anchor returns the point on the circle boundary in the given direction from its center,
// or the center itself for DirectionNone. A circle with a negative radius has no boundary,
// as it contains and intersects nothing, and anchors everything at its center rather than
// on the far side, where a negative length would put it.
// For integer T a diagonal anchor is rounded like
// Direction.Vector and only approximates the boundary: at a small radius it can land outside
// the circle, so Circ(Pt(0, 0), 1).Anchor(BottomRight) is (1,1), which Contains rejects.
func (c Circle[T]) Anchor(direction Direction) Point[T] {
	if c.Radius < 0 {
		return c.Center
	}

	return c.Center.Add(direction.Vector(c.Radius))
}

// Translate creates a new Circle translated by the given vector.
func (c Circle[T]) Translate(vector Vector[T]) Circle[T] {
	return Circle[T]{c.Center.Add(vector), c.Radius}
}

// MoveTo creates a new Circle with the same radius and the center set to point.
func (c Circle[T]) MoveTo(point Point[T]) Circle[T] {
	return Circle[T]{point, c.Radius}
}

// Scale creates a new Circle with radius scaled by the given factor.
func (c Circle[T]) Scale(factor float64) Circle[T] {
	return Circle[T]{c.Center, Multiply(c.Radius, factor)}
}

// Resize creates a new Circle with the given radius.
func (c Circle[T]) Resize(radius T) Circle[T] {
	return Circle[T]{c.Center, radius}
}

// Grow creates a new Circle with radius increased by amount, clamped to zero.
func (c Circle[T]) Grow(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius+amount, 0)}
}

// Shrink creates a new Circle with radius decreased by amount, clamped to zero.
func (c Circle[T]) Shrink(amount T) Circle[T] {
	return Circle[T]{c.Center, max(c.Radius-amount, 0)}
}

// Contains reports whether the given point lies within the circle, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains: a float point a rounding
// error outside the radius, such as an Anchor, is still contained.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.reaches(c.Center.Float().DistanceSquaredTo(point.Float()))
}

// DistanceTo returns the distance from the given point to the nearest point of the circle:
// zero exactly where Contains holds, so a point within Epsilon of T of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the center less the radius.
func (c Circle[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.Center.Float().DistanceSquaredTo(point.Float())
	if c.reaches(distanceSquared) {
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

// Intersects reports whether the circles overlap. Touching circles intersect, within Epsilon
// of T, the same closed convention as Contains, and a circle with a negative radius intersects
// nothing. The radii are summed in float64, so a narrow integer T cannot overflow the threshold.
func (c Circle[T]) Intersects(circle Circle[T]) bool {
	if c.Radius < 0 || circle.Radius < 0 {
		return false
	}

	distance := c.Center.Subtract(circle.Center).Length()
	threshold := float64(c.Radius) + float64(circle.Radius)

	return LessOrEqualDelta(distance, threshold, Epsilon[T]())
}

// Intersection returns the points where the circles cross: two for overlapping circles, one
// for tangent ones, within Epsilon of T like Intersects, and none for circles apart, nested, or
// with a negative radius. Coincident circles share every point and also return none. The
// second point is the mirror of the first across the line of centers. Circles with centers
// within Epsilon of T of each other count as coincident. For integer T the points are rounded
// like every other result stored into T.
func (c Circle[T]) Intersection(circle Circle[T]) []Point[T] {
	if c.Radius < 0 || circle.Radius < 0 {
		return nil
	}

	direction := circle.Center.Subtract(c.Center).Float()
	distance := direction.Length()
	r1, r2 := float64(c.Radius), float64(circle.Radius)
	epsilon := Epsilon[T]()

	if LessOrEqualDelta(distance, 0, epsilon) || !LessOrEqualDelta(distance, r1+r2, epsilon) || !LessOrEqualDelta(math.Abs(r1-r2), distance, epsilon) {
		return nil
	}

	along := (r1*r1 - r2*r2 + distance*distance) / (2 * distance)
	middle := c.Center.Float().Add(direction.Resize(along))

	if EqualDelta(distance, r1+r2, epsilon) || EqualDelta(distance, math.Abs(r1-r2), epsilon) {
		return []Point[T]{{Cast[T](middle.X), Cast[T](middle.Y)}}
	}

	normal := direction.Normal().Resize(math.Sqrt(max(r1*r1-along*along, 0)))
	first, second := middle.Add(normal), middle.Add(normal.Negate())

	return []Point[T]{{Cast[T](first.X), Cast[T](first.Y)}, {Cast[T](second.X), Cast[T](second.Y)}}
}

// IntersectsRectangle reports whether the circle and the rectangle overlap, as
// Rectangle.IntersectsCircle does.
func (c Circle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.IntersectsCircle(c)
}

// IntersectsLine reports whether the circle and the segment share a point, as
// Line.IntersectsCircle does.
func (c Circle[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsCircle(c)
}

// IntersectionLine returns the points where the segment crosses the circle boundary, as
// Line.IntersectionCircle does.
func (c Circle[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionCircle(c)
}

// IntersectsPolygon reports whether the circle and the polygon share a point, as
// Polygon.IntersectsCircle does.
func (c Circle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsCircle(c)
}

// reaches reports whether a point at the given squared distance from the center lies within
// the circle, boundary included within Epsilon of T. It is the one comparison against the
// radius that Contains, DistanceTo and every IntersectsCircle make, on the squared distance so
// that no test pays a square root and all of them round alike at the boundary. A negative
// radius reaches nothing.
func (c Circle[T]) reaches(distanceSquared float64) bool {
	reach := float64(c.Radius) + Epsilon[T]()

	return c.Radius >= 0 && distanceSquared <= reach*reach
}

// touches reports whether a point at the given squared distance from the center lies on the
// boundary within Epsilon of T: reached, and no deeper inside than the tolerance. It is the
// endpoint test of Line.IntersectionCircle, built on reaches so it agrees with IntersectsCircle.
// Both take the squared distance from the center as Point.DistanceSquaredTo gives it in float64.
func (c Circle[T]) touches(distanceSquared float64) bool {
	inner := max(float64(c.Radius)-Epsilon[T](), 0)

	return c.reaches(distanceSquared) && distanceSquared >= inner*inner
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

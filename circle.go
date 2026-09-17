package geom

import (
	"fmt"
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
// centered on the circle.
func (c Circle[T]) Bounds() Rectangle[T] {
	return Rectangle[T]{c.Center, Size[T]{c.Diameter(), c.Diameter()}}
}

// Anchor returns the point on the circle boundary in the given direction from its center,
// or the center itself for DirectionNone.
// For integer T a diagonal anchor is rounded like
// Direction.Vector and only approximates the boundary: at a small radius it can land outside
// the circle, so Circ(Pt(0, 0), 1).Anchor(BottomRight) is (1,1), which Contains rejects.
func (c Circle[T]) Anchor(direction Direction) Point[T] {
	return c.Center.Add(direction.Vector(c.Radius))
}

// Equal checks for equal center and radius with given circle.
func (c Circle[T]) Equal(circle Circle[T]) bool {
	return c.Center.Equal(circle.Center) && Equal(c.Radius, circle.Radius)
}

// IsZero checks if center point and radius are zero.
func (c Circle[T]) IsZero() bool {
	return c.Center.IsZero() && Equal(c.Radius, 0)
}

// Contains reports whether the given point lies within the circle, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains: a float point a rounding
// error outside the radius, such as an Anchor, is still contained.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.Center.Subtract(point).LessOrEqual(c.Radius)
}

// Intersects reports whether the circles overlap. Touching circles intersect, within Epsilon
// of T, the same closed convention as Contains. The radii are summed in float64, so a narrow
// integer T cannot overflow the threshold.
func (c Circle[T]) Intersects(circle Circle[T]) bool {
	distance := c.Center.Subtract(circle.Center).Length()
	threshold := float64(c.Radius) + float64(circle.Radius)

	return LessOrEqualDelta(distance, threshold, Epsilon[T]())
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

// IntersectsPolygon reports whether the circle and the polygon share a point, as
// Polygon.IntersectsCircle does.
func (c Circle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsCircle(c)
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

package geom

import (
	"fmt"
	"math"
)

// Circle is a 2D circle.
type Circle[T Number] struct {
	Center Point[T] `json:",inline"`
	Radius T        `json:"r"`
}

// C is shorthand for Circle{center, radius}.
func C[T Number](center Point[T], radius T) Circle[T] {
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

// Multiple creates a new Circle with radius scaled by the given factor.
func (c Circle[T]) Scale(scale float64) Circle[T] {
	return Circle[T]{c.Center, Multiple(c.Radius, scale)}
}

// Resize creates a new Circle with the given radius.
func (c Circle[T]) Resize(radius T) Circle[T] {
	return Circle[T]{c.Center, radius}
}

// Expand creates a new Circle with radius increased by amount.
func (c Circle[T]) Expand(amount T) Circle[T] {
	return Circle[T]{c.Center, c.Radius + amount}
}

// Shrunk creates a new Circle with radius decreased by amount.
func (c Circle[T]) Shrunk(amount T) Circle[T] {
	return Circle[T]{c.Center, c.Radius - amount}
}

// Area returns the circle area: (π * radius^2)
func (c Circle[T]) Area() float64 {
	return math.Pi * float64(c.Radius*c.Radius)
}

// Circumference returns the circle circumference (2 * π * radius).
func (c Circle[T]) Circumference() float64 {
	return 2 * math.Pi * float64(c.Radius)
}

// Diameter returns the circle diameter (2 * radius).
func (c Circle[T]) Diameter() T {
	return c.Radius * 2
}

// Bounds returns the axis-aligned bounding rectangle.
func (c Circle[T]) Bounds() Rectangle[T] {
	return Rectangle[T]{c.Center, Size[T]{c.Radius, c.Radius}}
}

// Equal checks for equal center and radius with given circle.
func (c Circle[T]) Equal(circle Circle[T]) bool {
	return c.Center.Equal(circle.Center) && Equal(c.Radius, circle.Radius)
}

// Contains checks if the given point lies inside the circle.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.Center.Subtract(point).Less(c.Radius)
}

// String returns a string representation of the Circle.
func (c Circle[T]) String() string {
	return fmt.Sprintf("C(%s;%s)", c.Center.String(), ToString(c.Radius))
}

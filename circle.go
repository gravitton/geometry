package geom

import (
	"fmt"
)

// Circle is a 2D circle.
type Circle[T Number] struct {
	Center Point[T] `json:"c"`
	Radius T        `json:"r"`
}

// C is shorthand for Circle[T]{Point[T]{x, y}, radius}.
func C[T Number](x, y, radius T) Circle[T] {
	return Circle[T]{Point[T]{x, y}, radius}
}

// Translate creates a new Circle by adding the given vector to circle center point.
func (c Circle[T]) Translate(change Vector[T]) Circle[T] {
	return Circle[T]{c.Center.Add(change), c.Radius}
}

// Move creates a new Circle with center in given point.
func (c Circle[T]) MoveTo(point Point[T]) Circle[T] {
	return Circle[T]{point, c.Radius}
}

// Scale returns a new Circle with scaled size.
func (c Circle[T]) Scale(scale float64) Circle[T] {
	return Circle[T]{c.Center, Scale(c.Radius, scale)}
}

// Scale returns a new Circle with given radius.
func (c Circle[T]) Resize(radius T) Circle[T] {
	return Circle[T]{c.Center, radius}
}

// Equal checks for equal values with given circle.
func (c Circle[T]) Equal(other Circle[T]) bool {
	return c.Center.Equal(other.Center) && Equal(c.Radius, other.Radius)
}

// Contains checks if the given point is in the circle.
func (c Circle[T]) Contains(point Point[T]) bool {
	return c.Center.Sub(point).Less(c.Radius)
}

// String makes a string representation of the Circle.
func (c Circle[T]) String() string {
	return fmt.Sprintf("C(%s;%s)", c.Center.String(), ToString(c.Radius))
}

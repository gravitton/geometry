package geom

import (
	"fmt"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// P is shorthand for Point{x, y}.
func P[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// Add creates a new Point by adding the given vector to the current point.
func (p Point[T]) Add(vector Vector[T]) Point[T] {
	return Point[T]{p.X + vector.X, p.Y + vector.Y}
}

// AddXY creates a new Point by adding the given values to the current point.
func (p Point[T]) AddXY(deltaX, deltaY T) Point[T] {
	return Point[T]{p.X + deltaX, p.Y + deltaY}
}

// Subtract creates a new Vector from given point to current point.
func (p Point[T]) Subtract(point Point[T]) Vector[T] {
	return Vector[T]{p.X - point.X, p.Y - point.Y}
}

// Multiply creates a new Point by multiplying the given value to the current point.
func (p Point[T]) Multiply(scale float64) Point[T] {
	return Point[T]{Multiple(p.X, scale), Multiple(p.Y, scale)}
}

// MultiplyXY creates a new Point by multiplying the given values to the current point.
func (p Point[T]) MultiplyXY(scaleX, scaleY float64) Point[T] {
	return Point[T]{Multiple(p.X, scaleX), Multiple(p.Y, scaleY)}
}

// Divide creates a new Point by dividing the given value to the current point.
func (p Point[T]) Divide(scale float64) Point[T] {
	return Point[T]{Divide(p.X, scale), Divide(p.Y, scale)}
}

// DivideXY creates a new Point by dividing the given values to the current point.
func (p Point[T]) DivideXY(scaleX, scaleY float64) Point[T] {
	return Point[T]{Divide(p.X, scaleX), Divide(p.Y, scaleY)}
}

// DistanceTo return euclidean distance from the current point to the given point.
func (p Point[T]) DistanceTo(point Point[T]) float64 {
	return point.Subtract(p).Length()
}

// DistanceToSquared return euclidean distance squared (for faster comparison) from the current point to the given point.
func (p Point[T]) DistanceToSquared(point Point[T]) T {
	return point.Subtract(p).LengthSquared()
}

// Midpoint creates a new Point between current and given points.
func (p Point[T]) Midpoint(point Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, point.X), Midpoint(p.Y, point.Y)}
}

// Lerp creates a new Point in linear interpolation towards given point.
func (p Point[T]) Lerp(point Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, point.X, t), Lerp(p.Y, point.Y, t)}
}

// AngleTo return angle between the current point to the given point.
func (p Point[T]) AngleTo(point Point[T]) float64 {
	return point.Subtract(p).Angle()
}

// Equal checks for equal X and Y values with given point.
func (p Point[T]) Equal(point Point[T]) bool {
	return Equal(p.X, point.X) && Equal(p.Y, point.Y)
}

// Zero checks if X and Y values are 0.
func (p Point[T]) Zero() bool {
	return p.Equal(Point[T]{})
}

// XY returns the point X, Y values in standard order.
func (p Point[T]) XY() (T, T) {
	return p.X, p.Y
}

// String returns a string representing the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%s,%s)", ToString(p.X), ToString(p.Y))
}

// ZeroPoint creates a new Point with zero values (0,0).
func ZeroPoint[T Number]() Point[T] {
	return Point[T]{}
}

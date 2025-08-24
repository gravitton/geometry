package geom

import (
	"fmt"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// P is shorthand for Point[T]{x, y}.
func P[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// XY returns the point X, Y values in standard order.
func (p Point[T]) XY() (T, T) {
	return p.X, p.Y
}

// Add creates a new Point by adding the given vector to the current point.
func (p Point[T]) Add(v Vector[T]) Point[T] {
	return Point[T]{p.X + v.X, p.Y + v.Y}
}

// AddXY creates a new Point by adding the given values to the current point.
func (p Point[T]) AddXY(dx, dy T) Point[T] {
	return Point[T]{p.X + dx, p.Y + dy}
}

// Subtract creates a new Vector from given point to current point.
func (p Point[T]) Subtract(o Point[T]) Vector[T] {
	return Vector[T]{p.X - o.X, p.Y - o.Y}
}

// Multiply creates a new Point by multiplying the given value to the current point.
func (p Point[T]) Multiply(s float64) Point[T] {
	return Point[T]{Multiple(p.X, s), Multiple(p.Y, s)}
}

// MultiplyXY creates a new Point by multiplying the given values to the current point.
func (p Point[T]) MultiplyXY(sx, sy float64) Point[T] {
	return Point[T]{Multiple(p.X, sx), Multiple(p.Y, sy)}
}

// Divide creates a new Point by dividing the given value to the current point.
func (p Point[T]) Divide(s float64) Point[T] {
	return Point[T]{Divide(p.X, s), Divide(p.Y, s)}
}

// DivideXY creates a new Point by dividing the given values to the current point.
func (p Point[T]) DivideXY(sx, sy float64) Point[T] {
	return Point[T]{Divide(p.X, sx), Divide(p.Y, sy)}
}

// Midpoint creates a new Point between current and given points.
func (p Point[T]) Midpoint(other Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, other.X), Midpoint(p.Y, other.Y)}
}

// Lerp creates a new Point in linear interpolation towards given point.
func (p Point[T]) Lerp(other Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, other.X, t), Lerp(p.Y, other.Y, t)}
}

// DistanceTo return euclidean distance from the current point to the given point.
func (p Point[T]) DistanceTo(other Point[T]) float64 {
	return other.Subtract(p).Length()
}

// DistanceToSquared return euclidean distance squared (for faster comparison) from the current point to the given point.
func (p Point[T]) DistanceToSquared(other Point[T]) T {
	return other.Subtract(p).LengthSquared()
}

// AngleTo return angle between the current point to the given point.
func (p Point[T]) AngleTo(other Point[T]) float64 {
	return other.Subtract(p).Angle()
}

// Equal checks for equal X and Y values with given point.
func (p Point[T]) Equal(other Point[T]) bool {
	return Equal(p.X, other.X) && Equal(p.Y, other.Y)
}

// Zero checks if X and Y values are 0.
func (p Point[T]) Zero() bool {
	return p.Equal(Point[T]{})
}

// String returns a string representing the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%s,%s)", ToString(p.X), ToString(p.Y))
}

// ZeroPoint creates a new Point with zero values (0,0).
func ZeroPoint[T Number]() Point[T] {
	return Point[T]{}
}

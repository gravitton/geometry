package geom

import (
	"fmt"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Pt is shorthand for Point{x, y}.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// Transform creates a new Point by applying the given matrix to the current point.
// For integer T, the float64 result of each component is rounded; rotations and
// non-integer scales lose precision.
func (p Point[T]) Transform(matrix Matrix[float64]) Point[T] {
	return Point[T]{Cast[T](matrix.A*float64(p.X) + matrix.B*float64(p.Y) + matrix.C), Cast[T](matrix.D*float64(p.X) + matrix.E*float64(p.Y) + matrix.F)}
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
func (p Point[T]) Multiply(factor float64) Point[T] {
	return Point[T]{Multiply(p.X, factor), Multiply(p.Y, factor)}
}

// MultiplyXY creates a new Point by multiplying the given values to the current point.
func (p Point[T]) MultiplyXY(factorX, factorY float64) Point[T] {
	return Point[T]{Multiply(p.X, factorX), Multiply(p.Y, factorY)}
}

// Divide creates a new Point by dividing the given value to the current point.
func (p Point[T]) Divide(factor float64) Point[T] {
	return Point[T]{Divide(p.X, factor), Divide(p.Y, factor)}
}

// DivideXY creates a new Point by dividing the given values to the current point.
func (p Point[T]) DivideXY(factorX, factorY float64) Point[T] {
	return Point[T]{Divide(p.X, factorX), Divide(p.Y, factorY)}
}

// Abs creates a new Point with absolute X and Y.
func (p Point[T]) Abs() Point[T] {
	return Point[T]{Abs(p.X), Abs(p.Y)}
}

// Round creates a new Point by rounding X, Y values to the nearest integer.
func (p Point[T]) Round() Point[T] {
	return Point[T]{Round(p.X), Round(p.Y)}
}

// Floor creates a new Point by rounding down X, Y values to the nearest integer.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{Floor(p.X), Floor(p.Y)}
}

// Ceil creates a new Point by rounding up X, Y values to the nearest integer.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{Ceil(p.X), Ceil(p.Y)}
}

// DistanceTo returns the Euclidean distance from the current point to the given point.
func (p Point[T]) DistanceTo(point Point[T]) float64 {
	return point.Subtract(p).Length()
}

// DistanceSquaredTo returns the squared Euclidean distance to the given point (faster than DistanceTo for comparisons).
func (p Point[T]) DistanceSquaredTo(point Point[T]) T {
	return point.Subtract(p).LengthSquared()
}

// ManhattanDistanceTo returns the Manhattan (taxicab) distance from the current point to the given point.
func (p Point[T]) ManhattanDistanceTo(point Point[T]) T {
	return Abs(point.X-p.X) + Abs(point.Y-p.Y)
}

// ChebyshevDistanceTo returns the Chebyshev distance (chessboard distance) from the current point to the given point.
// It is the maximum of the absolute differences of the coordinates: max(|dx|, |dy|).
func (p Point[T]) ChebyshevDistanceTo(point Point[T]) T {
	return max(Abs(point.X-p.X), Abs(point.Y-p.Y))
}

// OctileDistanceTo returns the Octile distance from the current point to the given point.
// Used in grid-based pathfinding where cardinal moves cost 1 and diagonal moves cost √2.
// It is: max(dx, dy) + (√2 - 1) * min(dx, dy).
func (p Point[T]) OctileDistanceTo(point Point[T]) float64 {
	dx := Abs(float64(point.X - p.X))
	dy := Abs(float64(point.Y - p.Y))

	return max(dx, dy) + (Sqrt2-1)*min(dx, dy)
}

// Midpoint creates a new Point between current and given points.
func (p Point[T]) Midpoint(point Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, point.X), Midpoint(p.Y, point.Y)}
}

// Lerp creates a new Point in linear interpolation towards given point.
func (p Point[T]) Lerp(point Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, point.X, t), Lerp(p.Y, point.Y, t)}
}

// AngleTo returns the angle in radians from the current point to the given point.
func (p Point[T]) AngleTo(point Point[T]) float64 {
	return point.Subtract(p).Angle()
}

// Equal checks for equal X and Y values with given point.
func (p Point[T]) Equal(point Point[T]) bool {
	return Equal(p.X, point.X) && Equal(p.Y, point.Y)
}

// IsZero checks if X and Y values are zero.
func (p Point[T]) IsZero() bool {
	return p.Equal(Point[T]{})
}

// XY returns the point X, Y values in standard order.
func (p Point[T]) XY() (T, T) {
	return p.X, p.Y
}

// Vector converts the point to a Vector.
func (p Point[T]) Vector() Vector[T] {
	return Vector[T](p)
}

// Int converts the point to a Point[int].
func (p Point[T]) Int() Point[int] {
	return Point[int]{Cast[int](float64(p.X)), Cast[int](float64(p.Y))}
}

// Float converts the point to a Point[float64].
func (p Point[T]) Float() Point[float64] {
	return Point[float64]{float64(p.X), float64(p.Y)}
}

// String returns a string representing the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%s,%s)", String(p.X), String(p.Y))
}

// ZeroPoint creates a new Point with zero values (0,0).
func ZeroPoint[T Number]() Point[T] {
	return Point[T]{}
}

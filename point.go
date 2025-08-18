package geom

import (
	"fmt"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Pt is shorthand for Point[T]{x, y}.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// Add creates a new Point by adding the given vector to the current point.
func (p Point[T]) Add(change Vector[T]) Point[T] {
	return Point[T]{p.X + change.X, p.Y + change.Y}
}

// Sub creates a new Vector by subtracting the given point from the current point.
func (p Point[T]) Sub(other Point[T]) Vector[T] {
	return Vector[T]{p.X - other.X, p.Y - other.Y}
}

// Midpoint creates a new Point between current and given points
func (p Point[T]) Midpoint(other Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, other.X), Midpoint(p.Y, other.Y)}
}

// Lerp creates a new Point in linear interpolation towards given point
func (p Point[T]) Lerp(other Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, other.X, t), Lerp(p.Y, other.Y, t)}
}

// DistanceTo return euclidean distance from the current point to the given point.
func (p Point[T]) DistanceTo(other Point[T]) float64 {
	return other.Sub(p).Length()
}

// DistanceToSquared return euclidean distance squared (for faster comparison) from the current point to the given point.
func (p Point[T]) DistanceToSquared(other Point[T]) T {
	return other.Sub(p).LengthSquared()
}

// AngleTo return angle between the current point to the given point.
func (p Point[T]) AngleTo(other Point[T]) float64 {
	return other.Sub(p).Angle()
}

// Equal checks for equal X and Y values with given point.
func (p Point[T]) Equal(other Point[T]) bool {
	return Equal(p.X, other.X) && Equal(p.Y, other.Y)
}

// Zero checks if X and Y values are 0.
func (p Point[T]) Zero() bool {
	// return Equal(p.X, 0) && Equal(p.Y, 0)
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

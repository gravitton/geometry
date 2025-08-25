package geom

import (
	"fmt"
	"math"
)

// Vector is a 2D Vector.
type Vector[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// V is shorthand for Vector{x, y}.
func V[T Number](x, y T) Vector[T] {
	return Vector[T]{x, y}
}

// VecFromAngle is shorthand for V(1,0).Rotate(angle)
func VecFromAngle[T Number](angle float64, length T) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(length) * cos), Cast[T](float64(length) * sin)}
}

// Add creates a new Vector by adding the given vector to the current vector.
func (v Vector[T]) Add(vector Vector[T]) Vector[T] {
	return Vector[T]{v.X + vector.X, v.Y + vector.Y}
}

// AddXY creates a new Vector by adding the given values to the current vector.
func (v Vector[T]) AddXY(deltaX, deltaY T) Vector[T] {
	return Vector[T]{v.X + deltaX, v.Y + deltaY}
}

// Subtract creates a new Vector by subtracting the given vector from the current vector.
func (v Vector[T]) Subtract(vector Vector[T]) Vector[T] {
	return Vector[T]{v.X - vector.X, v.Y - vector.Y}
}

// SubtractXY creates a new Vector by subtracting the given values from the current vector.
func (v Vector[T]) SubtractXY(deltaX, deltaY T) Vector[T] {
	return Vector[T]{v.X - deltaX, v.Y - deltaY}
}

// Multiply creates a new Vector by multiplying the given value to the current vector.
func (v Vector[T]) Multiply(factor float64) Vector[T] {
	return Vector[T]{Multiple(v.X, factor), Multiple(v.Y, factor)}
}

// MultiplyXY creates a new Vector by multiplying the given values to the current vector.
func (v Vector[T]) MultiplyXY(factorX, factorY float64) Vector[T] {
	return Vector[T]{Multiple(v.X, factorX), Multiple(v.Y, factorY)}
}

// Divide creates a new Vector by dividing the given value to the current vector.
func (v Vector[T]) Divide(factor float64) Vector[T] {
	return Vector[T]{Divide(v.X, factor), Divide(v.Y, factor)}
}

// DivideXY creates a new Vector by dividing the given values to the current vector.
func (v Vector[T]) DivideXY(factorX, factorY float64) Vector[T] {
	return Vector[T]{Divide(v.X, factorX), Divide(v.Y, factorY)}
}

// Negate creates a new Vector with opposite direction.
func (v Vector[T]) Negate() Vector[T] {
	return Vector[T]{-v.X, -v.Y}
}

// Rotate creates a new Vector rotated by the given angle (in radians).
func (v Vector[T]) Rotate(angle float64) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(v.X)*cos - float64(v.Y)*sin), Cast[T](float64(v.X)*sin + float64(v.Y)*cos)}
}

// Resize creates a new Vector resized to the given length.
func (v Vector[T]) Resize(length float64) Vector[T] {
	return v.Multiply(length / v.Length())
}

// Normalize creates a new Vector resized to a length of 1.
func (v Vector[T]) Normalize() Vector[T] {
	if v.Zero() {
		return Vector[T]{1, 0}
	}

	return v.Resize(1)
}

// Abs creates a new Vector with absolute X and Y.
func (v Vector[T]) Abs() Vector[T] {
	return Vector[T]{Abs(v.X), Abs(v.Y)}
}

// Dot returns dot (scalar) product of two vectors.
func (v Vector[T]) Dot(vector Vector[T]) T {
	return v.X*vector.X + v.Y*vector.Y
}

// Cross returns cross product of two vectors.
func (v Vector[T]) Cross(vector Vector[T]) T {
	return v.X*vector.Y - v.Y*vector.X
}

// Normal creates a new Vector as normal to current vector. Faster equivalent to Rotate(math.Pi/2).
func (v Vector[T]) Normal() Vector[T] {
	return Vector[T]{-v.Y, v.X}
}

// Length returns the Vector's length (magnitude).
func (v Vector[T]) Length() float64 {
	return math.Hypot(float64(v.X), float64(v.Y))
}

// LengthSquared returns the Vector's length (magnitude) squared (for faster comparison).
func (v Vector[T]) LengthSquared() T {
	return v.X*v.X + v.Y*v.Y
}

// Angle returns Vector's angle (in radian).
func (v Vector[T]) Angle() float64 {
	return math.Atan2(float64(v.Y), float64(v.X))
}

// Equal checks for equal X and Y values with given vector.
func (v Vector[T]) Equal(vector Vector[T]) bool {
	return Equal(v.X, vector.X) && Equal(v.Y, vector.Y)
}

// Zero checks if X and Y values are 0.
func (v Vector[T]) Zero() bool {
	return v.Equal(Vector[T]{})
}

// Unit checks if Vector is normalized.
func (v Vector[T]) Unit() bool {
	return Equal(v.LengthSquared(), 1.0)
}

// Less checks if Vector length is less than given value.
func (v Vector[T]) Less(value T) bool {
	return v.LengthSquared() < value*value
}

// XY returns the point X, Y values in standard order.
func (v Vector[T]) XY() (T, T) {
	return v.X, v.Y
}

// String returns a string representing the vector.
func (v Vector[T]) String() string {
	return fmt.Sprintf("⟨%s,%s⟩", ToString(v.X), ToString(v.Y))
}

// ZeroVector creates a new Vector with zero values (0,0).
func ZeroVector[T Number]() Vector[T] {
	return Vector[T]{}
}

// IdentityVector creates a new Vector with identity values (+1,+1).
func IdentityVector[T Number]() Vector[T] {
	return Vector[T]{1, 1}
}

// UpVector creates a new Vector with up (+y) direction (0,+1)
func UpVector[T Number]() Vector[T] {
	return Vector[T]{0, 1}
}

// DownVector creates a new Vector with down (-y) direction (0,-1)
func DownVector[T Number]() Vector[T] {
	return Vector[T]{0, -1}
}

// RightVector creates a new Vector with right (+x) direction (+1,0)
func RightVector[T Number]() Vector[T] {
	return Vector[T]{1, 0}
}

// LeftVector creates a new Vector with left (-x) direction (-1,0)
func LeftVector[T Number]() Vector[T] {
	return Vector[T]{-1, 0}
}

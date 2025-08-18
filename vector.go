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

// Vec is shorthand for Vector[T]{x, y}.
func Vec[T Number](x, y T) Vector[T] {
	return Vector[T]{x, y}
}

// Add creates a new Vector by adding the given vector to the current vector.
func (p Vector[T]) Add(other Vector[T]) Vector[T] {
	return Vector[T]{p.X + other.X, p.Y + other.Y}
}

// Sub creates a new Vector by subtracting the given vector from the current vector.
func (p Vector[T]) Sub(other Vector[T]) Vector[T] {
	return Vector[T]{p.X - other.X, p.Y - other.Y}
}

// Scale creates a new Vector with uniform scaled X and Y.
func (v Vector[T]) Scale(scale float64) Vector[T] {
	return Vector[T]{Scale(v.X, scale), Scale(v.Y, scale)}
}

// Scale creates a new Vector with scaled X and Y.
func (v Vector[T]) Stretch(scaleX, scaleY float64) Vector[T] {
	return Vector[T]{Scale(v.X, scaleX), Scale(v.Y, scaleY)}
}

// Dot returns dot (scalar) product of two vectors.
func (v Vector[T]) Dot(other Vector[T]) T {
	return v.X*other.X + v.Y*other.Y
}

// Cross returns cross product of two vectors.
func (v Vector[T]) Cross(other Vector[T]) T {
	return v.X*other.Y - v.Y*other.X
}

// Length returns the Vector's length (magnitude).
func (v Vector[T]) Length() float64 {
	return math.Hypot(float64(v.X), float64(v.Y))
}

// LengthSquared returns the Vector's length (magnitude) squared (for faster comparison).
func (v Vector[T]) LengthSquared() T {
	return v.X*v.X + v.Y*v.Y
}

// Negate creates a new Vector with opposite direction.
func (v Vector[T]) Negate() Vector[T] {
	return v.Scale(-1)
}

// Resize creates a new Vector resized to the given length.
func (v Vector[T]) Resize(length float64) Vector[T] {
	return v.Scale(length / v.Length())
}

// Unit creates a new Vector resized to a length of 1.
func (v Vector[T]) Unit() Vector[T] {
	return v.Resize(1)
}

// Abs creates a new Vector with absolute X and Y.
func (v Vector[T]) Abs() Vector[T] {
	return Vector[T]{Abs(v.X), Abs(v.Y)}
}

// Angle returns Vector's angle (in radian).
func (v Vector[T]) Angle() float64 {
	return math.Atan2(float64(v.Y), float64(v.X))
}

// Rotate creates a new Vector rotated by the given angle (in radians).
func (v Vector[T]) Rotate(angle float64) Vector[T] {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return Vector[T]{T(float64(v.X)*cosA - float64(v.Y)*sinA), T(float64(v.X)*sinA + float64(v.Y)*cosA)}
}

// Equal checks for equal X and Y values with given vector.
func (p Vector[T]) Equal(other Vector[T]) bool {
	return Equal(p.X, other.X) && Equal(p.Y, other.Y)
}

// Zero checks if X and Y values are 0.
func (p Vector[T]) Zero() bool {
	// return Equal(p.X, 0) && Equal(p.Y, 0)
	return p.Equal(Vector[T]{})
}

// Less checks the Vector length is less than given value.
func (v Vector[T]) Less(value T) bool {
	return v.LengthSquared() < value*value
}

// String returns a string representing the vector.
func (p Vector[T]) String() string {
	return fmt.Sprintf("⟨%s,%s⟩", ToString(p.X), ToString(p.Y))
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

// UpVector creates a new Vector with down (-y) direction (0,-1)
func DownVector[T Number]() Vector[T] {
	return Vector[T]{0, -1}
}

// UpVector creates a new Vector with right (+x) direction (+1,0)
func RightVector[T Number]() Vector[T] {
	return Vector[T]{1, 0}
}

// UpVector creates a new Vector with left (-x) direction (-1,0)
func LeftVector[T Number]() Vector[T] {
	return Vector[T]{-1, 0}
}

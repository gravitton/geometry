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

// Vec is shorthand for Vector{x, y}.
func Vec[T Number](x, y T) Vector[T] {
	return Vector[T]{x, y}
}

// Transform creates a new Vector by applying the given matrix to the current vector.
// For integer T, the float64 result of each component is rounded; rotations and
// non-integer scales lose precision.
func (v Vector[T]) Transform(matrix Matrix[float64]) Vector[T] {
	return Vector[T]{Cast[T](matrix.A*float64(v.X) + matrix.B*float64(v.Y)), Cast[T](matrix.D*float64(v.X) + matrix.E*float64(v.Y))}
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
	return Vector[T]{Multiply(v.X, factor), Multiply(v.Y, factor)}
}

// MultiplyXY creates a new Vector by multiplying the given values to the current vector.
func (v Vector[T]) MultiplyXY(factorX, factorY float64) Vector[T] {
	return Vector[T]{Multiply(v.X, factorX), Multiply(v.Y, factorY)}
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
// For integer T, sin/cos components are rounded; only multiples of 90° give exact results.
func (v Vector[T]) Rotate(angle float64) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(v.X)*cos - float64(v.Y)*sin), Cast[T](float64(v.X)*sin + float64(v.Y)*cos)}
}

// Resize creates a new Vector resized to the given length.
// For integer T, the result is rounded and the actual length may differ from the requested value.
func (v Vector[T]) Resize(length float64) Vector[T] {
	return v.Multiply(length / v.Length())
}

// Normalize creates a new Vector resized to a length of 1.
// For integer T, the result is one of the four axis-aligned unit vectors (±1,0)/(0,±1);
// the zero vector returns (1,0) by convention.
func (v Vector[T]) Normalize() Vector[T] {
	if v.IsZero() {
		return Vector[T]{1, 0}
	}

	unit := v.Resize(1)

	if isIntType[T]() && unit.X == unit.Y {
		if v.X > v.Y {
			unit.X = 1
			unit.Y = 0
		} else {
			unit.X = 0
			unit.Y = 1
		}
	}

	return unit
}

// Abs creates a new Vector with absolute X and Y.
func (v Vector[T]) Abs() Vector[T] {
	return Vector[T]{Abs(v.X), Abs(v.Y)}
}

// Round creates a new Vector by rounding X, Y values to the nearest integer.
func (v Vector[T]) Round() Vector[T] {
	return Vector[T]{Round(v.X), Round(v.Y)}
}

// Floor creates a new Vector by rounding down X, Y values to the nearest integer.
func (v Vector[T]) Floor() Vector[T] {
	return Vector[T]{Floor(v.X), Floor(v.Y)}
}

// Ceil creates a new Vector by rounding up X, Y values to the nearest integer.
func (v Vector[T]) Ceil() Vector[T] {
	return Vector[T]{Ceil(v.X), Ceil(v.Y)}
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

// Angle returns the vector's angle in radians.
func (v Vector[T]) Angle() float64 {
	return math.Atan2(float64(v.Y), float64(v.X))
}

// Lerp creates a new Vector in linear interpolation towards given vector.
func (v Vector[T]) Lerp(vector Vector[T], t float64) Vector[T] {
	return Vector[T]{Lerp(v.X, vector.X, t), Lerp(v.Y, vector.Y, t)}
}

// Equal checks for equal X and Y values with given vector.
func (v Vector[T]) Equal(vector Vector[T]) bool {
	return Equal(v.X, vector.X) && Equal(v.Y, vector.Y)
}

// IsZero checks if X and Y values are zero.
func (v Vector[T]) IsZero() bool {
	return v.Equal(Vector[T]{})
}

// IsOne checks if X and Y values are (1,1).
func (v Vector[T]) IsOne() bool {
	return v.Equal(Vector[T]{1, 1})
}

// IsUp checks whether the vector points upward (-Y).
func (v Vector[T]) IsUp() bool {
	return v.Y < 0
}

// IsDown checks whether the vector points downward (+Y).
func (v Vector[T]) IsDown() bool {
	return v.Y > 0
}

// IsLeft checks whether the vector points leftward (-X).
func (v Vector[T]) IsLeft() bool {
	return v.X < 0
}

// IsRight checks whether the vector points rightward (+X).
func (v Vector[T]) IsRight() bool {
	return v.X > 0
}

// IsNormalized checks if Vector is normalized.
func (v Vector[T]) IsNormalized() bool {
	return Equal(v.LengthSquared(), 1.0)
}

// Less checks if Vector length is less than given value.
func (v Vector[T]) Less(value T) bool {
	return v.LengthSquared() < value*value
}

// XY returns the vector X, Y values in standard order.
func (v Vector[T]) XY() (T, T) {
	return v.X, v.Y
}

// Point converts the vector to a Point.
func (v Vector[T]) Point() Point[T] {
	return Point[T](v)
}

// Size converts the vector to a Size (using absolute component values).
func (v Vector[T]) Size() Size[T] {
	return Size[T]{Abs(v.X), Abs(v.Y)}
}

// Int converts the vector to a [int] vector.
func (v Vector[T]) Int() Vector[int] {
	return Vector[int]{Cast[int](float64(v.X)), Cast[int](float64(v.Y))}
}

// Float converts the vector to a [float64] vector.
func (v Vector[T]) Float() Vector[float64] {
	return Vector[float64]{float64(v.X), float64(v.Y)}
}

// String returns a string representing the vector.
func (v Vector[T]) String() string {
	return fmt.Sprintf("⟨%s,%s⟩", String(v.X), String(v.Y))
}

// VectorFromAngle returns a vector of the given length pointing in the direction of angle (in radians).
// For integer T, each component is independently rounded after multiplying by length; angles
// whose sin or cos falls near ±0.5 may round to 0 instead of ±1 due to float64 precision.
// Use float64 when directional accuracy matters.
func VectorFromAngle[T Number](angle float64, length T) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(length) * cos), Cast[T](float64(length) * sin)}
}

// ZeroVector creates a new Vector with zero values (0,0).
func ZeroVector[T Number]() Vector[T] {
	return Vector[T]{}
}

// OneVector creates a new Vector with identity values (+1,+1).
func OneVector[T Number]() Vector[T] {
	return Vector[T]{1, 1}
}

// UpVector returns the unit vector pointing up: (0, -1).
func UpVector[T Number]() Vector[T] {
	return Vector[T]{0, -1}
}

// DownVector returns the unit vector pointing down: (0, +1).
func DownVector[T Number]() Vector[T] {
	return Vector[T]{0, 1}
}

// LeftVector returns the unit vector pointing left: (-1, 0).
func LeftVector[T Number]() Vector[T] {
	return Vector[T]{-1, 0}
}

// RightVector returns the unit vector pointing right: (+1, 0).
func RightVector[T Number]() Vector[T] {
	return Vector[T]{1, 0}
}

// UpLeftVector creates a new unit Vector with up (-y) and left (-x) directions.
func UpLeftVector() Vector[float64] {
	return Vector[float64]{-OneOverSqrt2, -OneOverSqrt2}
}

// UpRightVector creates a new unit Vector with up (-y) and right (+x) directions.
func UpRightVector() Vector[float64] {
	return Vector[float64]{OneOverSqrt2, -OneOverSqrt2}
}

// DownLeftVector creates a new unit Vector with down (+y) and left (-x) directions.
func DownLeftVector() Vector[float64] {
	return Vector[float64]{-OneOverSqrt2, OneOverSqrt2}
}

// DownRightVector creates a new unit Vector with down (+y) and right (+x) directions.
func DownRightVector() Vector[float64] {
	return Vector[float64]{OneOverSqrt2, OneOverSqrt2}
}

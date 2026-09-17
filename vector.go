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

// ZeroVector creates a new Vector with zero values (0,0).
func ZeroVector[T Number]() Vector[T] {
	return Vector[T]{}
}

// OneVector creates a new Vector with identity values (+1,+1).
func OneVector[T Number]() Vector[T] {
	return Vector[T]{1, 1}
}

// VectorFromAngle creates a new Vector of the given length pointing at the given angle (in
// radians), measured in the standard math convention where Y grows upward.
// For integer T, both components are rounded; only multiples of 90° give exact results.
func VectorFromAngle[T Number](angle float64, length T) Vector[T] {
	return VectorFromAngleSize(angle, SzU(length))
}

// VectorFromAngleSize creates a new Vector at the given angle (in radians) on an ellipse
// with the given size as its semi-axes: (width*cos(angle), height*sin(angle)).
// For integer T, both components are rounded; only multiples of 90° give exact results.
func VectorFromAngleSize[T Number](angle float64, size Size[T]) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(size.Width) * cos), Cast[T](float64(size.Height) * sin)}
}

// Transform creates a new Vector by applying the given matrix to the current vector.
// The matrix is float-only, like an angle: convert an integer matrix with Matrix.Float first.
// For integer T, the float64 result of each component is rounded; rotations and non-integer scales lose precision.
func (v Vector[T]) Transform[M Float](matrix Matrix[M]) Vector[T] {
	x, y := float64(v.X), float64(v.Y)
	m := matrix.Float()

	return Vector[T]{Cast[T](m.A*x + m.B*y), Cast[T](m.D*x + m.E*y)}
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

// Rotate creates a new Vector rotated by the given angle (in radians), in the standard math
// convention where Y grows upward, so a positive angle appears clockwise on a screen with Y
// pointing down. This is the same sense as a positive Direction.Rotate step.
// For integer T, sin/cos components are rounded; only multiples of 90° give exact results.
func (v Vector[T]) Rotate(angle float64) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(v.X)*cos - float64(v.Y)*sin), Cast[T](float64(v.X)*sin + float64(v.Y)*cos)}
}

// Resize creates a new Vector resized to the given length. The zero vector has no direction
// and resizes along +X to (length,0), the same convention Normalize follows. Only the exact
// zero vector is treated this way: a vector shorter than Epsilon keeps its direction.
// For integer T, the result is rounded and the actual length may differ from the requested value.
func (v Vector[T]) Resize(length float64) Vector[T] {
	if !v.hasDirection() {
		return Vector[T]{Cast[T](length), 0}
	}

	return v.Multiply(length / v.Length())
}

// Normalize creates a new Vector resized to a length of 1.
// For integer T, the only vectors of length 1 are the four axis-aligned unit vectors, so the
// result snaps to the longer axis and keeps its sign, X winning a tie; the exact zero vector
// returns (1,0) by convention.
func (v Vector[T]) Normalize() Vector[T] {
	if !v.hasDirection() {
		return Vector[T]{1, 0}
	}

	if isIntType[T]() {
		if Abs(v.X) >= Abs(v.Y) {
			return Vector[T]{Sign(v.X), 0}
		}

		return Vector[T]{0, Sign(v.Y)}
	}

	return v.Resize(1)
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

// Direction returns the direction nearest to the vector, or DirectionNone for the exact zero vector.
func (v Vector[T]) Direction() Direction {
	if !v.hasDirection() {
		return DirectionNone
	}

	return DirectionFromAngle(v.Angle())
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

// hasDirection reports whether the vector points somewhere: only the exact zero vector does not.
// It deliberately ignores Epsilon, since a vector shorter than the tolerance still has a direction.
func (v Vector[T]) hasDirection() bool {
	return v.X != 0 || v.Y != 0
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

// Less reports whether the vector is shorter than the given length.
// No vector is shorter than a non-positive length.
func (v Vector[T]) Less(length T) bool {
	return length > 0 && v.LengthSquared() < length*length
}

// LessOrEqual reports whether the vector is at most the given length.
// No vector is at most a negative length.
func (v Vector[T]) LessOrEqual(length T) bool {
	return length >= 0 && v.LengthSquared() <= length*length
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

// Int converts the vector to a Vector[int].
func (v Vector[T]) Int() Vector[int] {
	return Vector[int]{Int(v.X), Int(v.Y)}
}

// Float converts the vector to a Vector[float64].
func (v Vector[T]) Float() Vector[float64] {
	return Vector[float64]{float64(v.X), float64(v.Y)}
}

// String returns a string representing the vector.
func (v Vector[T]) String() string {
	return fmt.Sprintf("⟨%s,%s⟩", String(v.X), String(v.Y))
}

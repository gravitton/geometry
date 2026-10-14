package geom

import (
	"math"
)

// Vector is a displacement, or a direction with a length.
type Vector[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Vec is shorthand for Vector{x, y}.
func Vec[T Number](x, y T) Vector[T] {
	return Vector[T]{x, y}
}

// ZeroVector returns the zero vector.
func ZeroVector[T Number]() Vector[T] {
	return Vector[T]{}
}

// OneVector returns the vector of both components one, the factors that scale nothing.
func OneVector[T Number]() Vector[T] {
	return Vector[T]{1, 1}
}

// VectorFromAngle returns the vector of the length at the angle in radians.
func VectorFromAngle[T Number](angle float64, length T) Vector[T] {
	return VectorFromAngleSize(angle, SzU(length))
}

// VectorFromAngleSize returns the vector at the angle in radians on the ellipse with the size as
// its semi-axes (width·cos, height·sin).
func VectorFromAngleSize[T Number](angle float64, size Size[T]) Vector[T] {
	sin, cos := math.Sincos(angle)

	return Vector[T]{Cast[T](float64(size.Width) * cos), Cast[T](float64(size.Height) * sin)}
}

// XY returns X and Y.
func (v Vector[T]) XY() (T, T) {
	return v.X, v.Y
}

// Length returns the length of the vector.
func (v Vector[T]) Length() float64 {
	return math.Hypot(float64(v.X), float64(v.Y))
}

// LengthSquared returns the square of Length.
func (v Vector[T]) LengthSquared() T {
	x, y := float64(v.X), float64(v.Y)

	return Cast[T](float64(x*x) + float64(y*y))
}

// Angle returns the angle of the vector in radians.
func (v Vector[T]) Angle() float64 {
	return math.Atan2(float64(v.Y), float64(v.X))
}

// Direction returns the direction nearest to the vector, or DirectionNone for the exact zero
// vector.
func (v Vector[T]) Direction() Direction {
	if !v.hasDirection() {
		return DirectionNone
	}

	return DirectionFromAngle(v.Angle())
}

// Add returns the sum of the vectors.
func (v Vector[T]) Add(vector Vector[T]) Vector[T] {
	return Vector[T]{v.X + vector.X, v.Y + vector.Y}
}

// AddXY returns the vector with the deltas added.
func (v Vector[T]) AddXY(deltaX, deltaY T) Vector[T] {
	return Vector[T]{v.X + deltaX, v.Y + deltaY}
}

// Subtract returns the difference of the vectors.
func (v Vector[T]) Subtract(vector Vector[T]) Vector[T] {
	return Vector[T]{v.X - vector.X, v.Y - vector.Y}
}

// SubtractXY returns the vector with the deltas subtracted.
func (v Vector[T]) SubtractXY(deltaX, deltaY T) Vector[T] {
	return Vector[T]{v.X - deltaX, v.Y - deltaY}
}

// Multiply returns the vector with both components multiplied by the factor.
func (v Vector[T]) Multiply(factor float64) Vector[T] {
	return Vector[T]{Multiply(v.X, factor), Multiply(v.Y, factor)}
}

// MultiplyXY returns the vector with each component multiplied by its own factor.
func (v Vector[T]) MultiplyXY(factorX, factorY float64) Vector[T] {
	return Vector[T]{Multiply(v.X, factorX), Multiply(v.Y, factorY)}
}

// Divide returns the vector with both components divided by the factor. It panics for a zero
// factor.
func (v Vector[T]) Divide(factor float64) Vector[T] {
	return Vector[T]{Divide(v.X, factor), Divide(v.Y, factor)}
}

// DivideXY returns the vector with each component divided by its own factor. It panics for a zero
// factor.
func (v Vector[T]) DivideXY(factorX, factorY float64) Vector[T] {
	return Vector[T]{Divide(v.X, factorX), Divide(v.Y, factorY)}
}

// Negate returns the vector pointing the opposite way.
func (v Vector[T]) Negate() Vector[T] {
	return Vector[T]{-v.X, -v.Y}
}

// Abs returns the vector with absolute components.
func (v Vector[T]) Abs() Vector[T] {
	return Vector[T]{Abs(v.X), Abs(v.Y)}
}

// Round returns the vector with components rounded to the nearest integer.
func (v Vector[T]) Round() Vector[T] {
	return Vector[T]{Round(v.X), Round(v.Y)}
}

// Floor returns the vector with components rounded down.
func (v Vector[T]) Floor() Vector[T] {
	return Vector[T]{Floor(v.X), Floor(v.Y)}
}

// Ceil returns the vector with components rounded up.
func (v Vector[T]) Ceil() Vector[T] {
	return Vector[T]{Ceil(v.X), Ceil(v.Y)}
}

// Lerp returns the vector interpolated toward the given one, extrapolated outside [0, 1].
func (v Vector[T]) Lerp(vector Vector[T], t float64) Vector[T] {
	return Vector[T]{Lerp(v.X, vector.X, t), Lerp(v.Y, vector.Y, t)}
}

// Slerp returns the vector interpolated toward the given one, the direction along the shorter arc
// as LerpAngle turns it and the length linearly, extrapolated outside [0, 1]. The zero vector takes
// the direction of the other, so the result runs as Lerp does, and opposite vectors turn toward
// increasing angle. A length taken past zero flips the direction, and a t of zero returns the
// receiver exactly.
func (v Vector[T]) Slerp(vector Vector[T], t float64) Vector[T] {
	from, to := v.Float(), vector.Float()

	switch {
	case !from.hasDirection():
		return to.Multiply(t).Cast[T]()
	case !to.hasDirection():
		return from.Multiply(1 - t).Cast[T]()
	}

	length := from.Length()
	turned := from.Rotate(LerpAngle(0, to.Angle()-from.Angle(), t))

	return turned.Multiply(Lerp(length, to.Length(), t) / length).Cast[T]()
}

// AtLeast returns the vector resized to the length where Less holds for it, and the vector itself
// otherwise, so a non-positive length changes nothing. The zero vector is resized along +X.
func (v Vector[T]) AtLeast(length T) Vector[T] {
	if !v.Less(length) {
		return v
	}

	return v.Resize(float64(length))
}

// AtMost returns the vector resized to the length where LessOrEqual fails for it, and the vector
// itself otherwise. A non-positive length gives the zero vector.
func (v Vector[T]) AtMost(length T) Vector[T] {
	if v.LessOrEqual(length) {
		return v
	}

	return v.Resize(float64(max(length, 0)))
}

// Transform returns the vector mapped by the matrix, without its translation. The matrix is
// float-only: convert an integer one with Matrix.Float first.
func (v Vector[T]) Transform[M Float](matrix Matrix[M]) Vector[T] {
	x, y := float64(v.X), float64(v.Y)
	a, b := float64(matrix.A), float64(matrix.B)
	d, e := float64(matrix.D), float64(matrix.E)

	return Vector[T]{Cast[T](float64(a*x) + float64(b*y)), Cast[T](float64(d*x) + float64(e*y))}
}

// Rotate returns the vector rotated by the angle in radians.
func (v Vector[T]) Rotate(angle float64) Vector[T] {
	sin, cos := math.Sincos(angle)
	x, y := float64(v.X), float64(v.Y)

	return Vector[T]{Cast[T](float64(x*cos) - float64(y*sin)), Cast[T](float64(x*sin) + float64(y*cos))}
}

// RotateToward returns the vector turned toward the angle in radians by at most the step along the
// shorter arc, keeping its length and never passing the angle. A vector exactly opposite turns
// toward increasing angle, a negative step turns away unclamped, and a zero step or the zero vector
// returns the receiver. For integer T a small step can leave the vector in place.
func (v Vector[T]) RotateToward(angle, step float64) Vector[T] {
	if !v.hasDirection() {
		return v
	}

	return v.Rotate(MoveTowardAngle(0, NormalizeAngle(angle-v.Angle()), step))
}

// MoveToward returns the vector moved toward the target by at most the step, direction and length
// together, never past it: where the distance left is at most the step, the result is the target
// exactly. A negative step moves away and is not clamped, and a vector at the target stays there.
// For integer T a step under half a unit can leave the vector in place.
func (v Vector[T]) MoveToward(target Vector[T], step float64) Vector[T] {
	from := v.Float()
	offset := target.Float().Subtract(from)

	distance := offset.Length()
	if distance == 0 || step >= distance {
		return target
	}

	return from.Add(offset.Multiply(step / distance)).Cast[T]()
}

// Resize returns the vector of the length in its own direction, the zero vector resized along +X. A
// vector of any nonzero length keeps its direction, however short. A negative length flips the
// direction. For integer T the length of the result can differ from the one asked.
func (v Vector[T]) Resize(length float64) Vector[T] {
	if !v.hasDirection() {
		return Vector[T]{Cast[T](length), 0}
	}

	current := v.Length()

	return Vector[T]{Cast[T](float64(v.X) / current * length), Cast[T](float64(v.Y) / current * length)}
}

// Normalize returns the vector of length one in its own direction, and (1,0) for the zero vector.
// For integer T it snaps to the axis of the longer component, keeping its sign, X on a tie.
func (v Vector[T]) Normalize() Vector[T] {
	if !v.hasDirection() {
		return Vector[T]{1, 0}
	}

	if isInt[T]() {
		if Abs(v.X) >= Abs(v.Y) {
			return Vector[T]{Sign(v.X), 0}
		}

		return Vector[T]{0, Sign(v.Y)}
	}

	return v.Resize(1)
}

// Normal returns the vector a quarter turn in the sense of a positive angle, exactly and of the
// same length. Normalize it for a unit normal, and Negate it for the other side.
func (v Vector[T]) Normal() Vector[T] {
	return Vector[T]{-v.Y, v.X}
}

// Project returns the component of the vector along the given one, and the zero vector for a zero
// one.
func (v Vector[T]) Project(vector Vector[T]) Vector[T] {
	a, b := v.Float(), vector.Float()

	length := b.LengthSquared()
	if length == 0 {
		return Vector[T]{}
	}

	return vector.Multiply(a.Dot(b) / length)
}

// Reject returns the component of the vector perpendicular to the given one, so Project and Reject
// sum to the vector. A zero vector rejects nothing and gives the vector itself.
func (v Vector[T]) Reject(vector Vector[T]) Vector[T] {
	return v.Subtract(v.Project(vector))
}

// Reflect returns the vector mirrored across the surface with the given normal, which need not be
// unit: the component along the normal flips. A zero normal returns the vector itself.
func (v Vector[T]) Reflect(normal Vector[T]) Vector[T] {
	a, b := v.Float(), normal.Float()

	length := b.LengthSquared()
	if length == 0 {
		return v
	}

	return v.Subtract(normal.Multiply(2 * a.Dot(b) / length))
}

// Dot returns the dot product of the vectors, exactly zero for perpendicular ones.
func (v Vector[T]) Dot(vector Vector[T]) T {
	a, b := v.Float(), vector.Float()

	return Cast[T](float64(a.X*b.X) + float64(a.Y*b.Y))
}

// Cross returns the cross product of the vectors, exactly zero for parallel ones.
func (v Vector[T]) Cross(vector Vector[T]) T {
	a, b := v.Float(), vector.Float()

	return Cast[T](float64(a.X*b.Y) - float64(a.Y*b.X))
}

// AngleBetween returns the unsigned angle between the vectors in radians, in [0, π], and zero where
// either is the zero vector.
func (v Vector[T]) AngleBetween(vector Vector[T]) float64 {
	a, b := v.Float(), vector.Float()

	return math.Atan2(math.Abs(a.Cross(b)), a.Dot(b))
}

// Equal reports whether the vectors have equal components.
func (v Vector[T]) Equal(vector Vector[T]) bool {
	return Equal(v.X, vector.X) && Equal(v.Y, vector.Y)
}

// IsZero reports whether the vector is the zero vector.
func (v Vector[T]) IsZero() bool {
	return v.Equal(Vector[T]{})
}

// IsFinite reports whether both components are neither NaN nor infinite, which Cast and Int need
// for an integer type. An integer vector always is.
func (v Vector[T]) IsFinite() bool {
	return IsFinite(v.X) && IsFinite(v.Y)
}

// IsOne reports whether both components are one.
func (v Vector[T]) IsOne() bool {
	return v.Equal(Vector[T]{1, 1})
}

// IsUp reports whether the vector runs up on screen, its Y negative, exactly. A diagonal is up and
// to a side at once; Direction names the nearest one.
func (v Vector[T]) IsUp() bool {
	return v.Y < 0
}

// IsDown reports whether the vector runs down on screen, its Y positive.
func (v Vector[T]) IsDown() bool {
	return v.Y > 0
}

// IsLeft reports whether the vector runs left, its X negative.
func (v Vector[T]) IsLeft() bool {
	return v.X < 0
}

// IsRight reports whether the vector runs right, its X positive.
func (v Vector[T]) IsRight() bool {
	return v.X > 0
}

// IsNormalized reports whether the vector has length one.
func (v Vector[T]) IsNormalized() bool {
	return equalSquared(v.Float().LengthSquared(), 1, Epsilon[T]())
}

// Less reports whether the vector is strictly shorter than the length, with no tolerance. No vector
// is shorter than a non-positive length.
func (v Vector[T]) Less(length T) bool {
	reach := float64(length)

	return reach > 0 && v.Float().LengthSquared() < reach*reach
}

// LessOrEqual reports whether the vector is at most the length, so its end lies in the circle of
// that radius about the origin. No vector is at most a negative length: the tolerance widens the
// boundary, never the sign.
func (v Vector[T]) LessOrEqual(length T) bool {
	return length >= 0 && lessOrEqualSquared(v.Float().LengthSquared(), float64(length), epsilonAt[T](float64(length)))
}

// hasDirection reports whether the vector is not the exact zero vector, ignoring the tolerance,
// since a vector shorter than it still has a direction.
func (v Vector[T]) hasDirection() bool {
	return v.X != 0 || v.Y != 0
}

// Point converts the vector to a Point.
func (v Vector[T]) Point() Point[T] {
	return Point[T](v)
}

// Size converts the vector to a Size, the inverse of Size.Vector. The sign of each component is
// kept: take Abs where an extent must not be negative.
func (v Vector[T]) Size() Size[T] {
	return Size[T]{v.X, v.Y}
}

// Cast converts the vector to a Vector of another number type, rounding as Cast does.
func (v Vector[T]) Cast[R Number]() Vector[R] {
	return Vector[R]{Cast[R](float64(v.X)), Cast[R](float64(v.Y))}
}

// Int converts the vector to a Vector[int].
func (v Vector[T]) Int() Vector[int] {
	return Vector[int]{Int(v.X), Int(v.Y)}
}

// Float converts the vector to a Vector[float64].
func (v Vector[T]) Float() Vector[float64] {
	return Vector[float64]{float64(v.X), float64(v.Y)}
}

// String returns the vector as its components in angle brackets, with no type name: ⟨x,y⟩.
func (v Vector[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = v.appendString(buf)

	return string(buf)
}

// appendString appends the vector as String prints it to buf and returns the extended buffer.
func (v Vector[T]) appendString(buf []byte) []byte {
	buf = append(buf, "⟨"...)
	buf = appendString(buf, v.X)
	buf = append(buf, ',')
	buf = appendString(buf, v.Y)
	buf = append(buf, "⟩"...)

	return buf
}

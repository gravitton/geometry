package geom

import (
	"math"
)

// Matrix is a 2D affine transform, mapping x' = A·x + B·y + C and y' = D·x + E·y + F, the third row
// [0 0 1] implied.
//
// Translate, Rotate, RotateAround, Shear and Reflect compose the matrix with the constructor of the
// same name on the right, applied first, and their Pre forms on the left, applied last. Scale,
// ScaleAround and Unscale take float64 factors, where ScaleMatrix and ScaleMatrixAround take T, so an
// integer matrix can be scaled by a half.
//
// For integer T, translation, integer scaling, reflection and quarter turns are exact. Any other
// rotation rounds into a different matrix, and Inverse and Unscale do unless the result is integer.
type Matrix[T Number] struct {
	A T `json:"a"` // scale X
	B T `json:"b"` // shear X (contribution of y to x')
	C T `json:"c"` // translate X
	D T `json:"d"` // shear Y (contribution of x to y')
	E T `json:"e"` // scale Y
	F T `json:"f"` // translate Y
}

// Mat is shorthand for Matrix{a, b, c, d, e, f}.
func Mat[T Number](a, b, c, d, e, f T) Matrix[T] {
	return Matrix[T]{a, b, c, d, e, f}
}

// IdentityMatrix returns the matrix that leaves every point in place.
func IdentityMatrix[T Number]() Matrix[T] {
	return Matrix[T]{
		1, 0, 0,
		0, 1, 0,
	}
}

// TranslationMatrix returns the matrix that moves by the deltas.
func TranslationMatrix[T Number](deltaX, deltaY T) Matrix[T] {
	return Matrix[T]{
		1, 0, deltaX,
		0, 1, deltaY,
	}
}

// RotationMatrix returns the matrix that rotates by the angle in radians about the origin. For
// integer T any angle but a multiple of a quarter turn rounds into a matrix that scales and shears
// as well.
func RotationMatrix[T Number](angle float64) Matrix[T] {
	sin, cos := math.Sincos(angle)

	return Matrix[T]{
		Cast[T](cos), Cast[T](-sin), 0,
		Cast[T](sin), Cast[T](cos), 0,
	}
}

// RotationMatrixAround returns the matrix that rotates by the angle in radians about the pivot,
// equal to the bit to the chain of translations and RotationMatrix. For integer T the pivot maps
// exactly onto itself.
func RotationMatrixAround[T Number](pivotX, pivotY T, angle float64) Matrix[T] {
	return RotationMatrix[T](angle).around(pivotX, pivotY)
}

// ScaleMatrix returns the matrix that scales by the factors about the origin.
func ScaleMatrix[T Number](factorX, factorY T) Matrix[T] {
	return Matrix[T]{
		factorX, 0, 0,
		0, factorY, 0,
	}
}

// ScaleMatrixAround returns the matrix that scales by the factors about the pivot, equal to the bit
// to the chain of translations and ScaleMatrix. The pivot maps onto itself.
func ScaleMatrixAround[T Number](pivotX, pivotY, factorX, factorY T) Matrix[T] {
	return ScaleMatrix(factorX, factorY).around(pivotX, pivotY)
}

// TRSMatrix returns the matrix that scales, then rotates by the angle in radians, then translates
// (TranslationMatrix·RotationMatrix·ScaleMatrix). Translation, Angle and Scaling read back parts
// that compose the same matrix: a negative X factor comes back as a half turn and a negative Y
// factor. For integer T each entry is rounded once.
func TRSMatrix[T Number](deltaX, deltaY T, angle float64, factorX, factorY T) Matrix[T] {
	sin, cos := math.Sincos(angle)

	return Matrix[T]{
		Multiply(factorX, cos), Multiply(factorY, -sin), deltaX,
		Multiply(factorX, sin), Multiply(factorY, cos), deltaY,
	}
}

// ShearMatrix returns the shear matrix: shearX is the contribution of y to x', shearY that of x to
// y'.
func ShearMatrix[T Number](shearX, shearY T) Matrix[T] {
	return Matrix[T]{
		1, shearX, 0,
		shearY, 1, 0,
	}
}

// ReflectionMatrix returns the matrix that mirrors across the axis: AxisHorizontal flips Y,
// AxisVertical flips X, and AxisNone gives the identity.
func ReflectionMatrix[T Number](axis Axis) Matrix[T] {
	switch axis {
	case AxisHorizontal:
		return ScaleMatrix[T](1, -1)
	case AxisVertical:
		return ScaleMatrix[T](-1, 1)
	default:
		return IdentityMatrix[T]()
	}
}

func (m Matrix[T]) around(pivotX, pivotY T) Matrix[T] {
	f := m.Float()
	x, y := float64(pivotX), float64(pivotY)

	return Matrix[T]{
		m.A, m.B, Cast[T](float64(f.A*-x) + float64(f.B*-y) + x),
		m.D, m.E, Cast[T](float64(f.D*-x) + float64(f.E*-y) + y),
	}
}

// Translation returns the translation the matrix applies, its C and F.
func (m Matrix[T]) Translation() Vector[T] {
	return Vector[T]{m.C, m.F}
}

// Angle returns the angle in radians of the transformed X axis, the rotation of a matrix built by
// TRSMatrix. A matrix that collapses the X axis gives the angle of its Y axis a quarter turn back,
// and the zero matrix gives zero.
func (m Matrix[T]) Angle() float64 {
	if m.A == 0 && m.D == 0 {
		return math.Atan2(-float64(m.B), float64(m.E))
	}

	return math.Atan2(float64(m.D), float64(m.A))
}

// Scaling returns the length of the transformed X axis and the determinant divided by it, negative
// for a reflection: the factors of a matrix built by TRSMatrix, and for a sheared one those that
// keep the X axis and the area. A matrix that collapses the X axis gives zero and the length of its
// Y axis.
func (m Matrix[T]) Scaling() Vector[T] {
	f := m.Float()

	x := math.Hypot(f.A, f.D)
	if x == 0 {
		return Vector[T]{0, Cast[T](math.Hypot(f.B, f.E))}
	}

	return Vector[T]{Cast[T](x), Cast[T](f.determinant() / x)}
}

// Determinant returns the determinant of the linear part.
func (m Matrix[T]) Determinant() T {
	return Cast[T](m.determinant())
}

// turnedAngle returns the angle a shape at the given angle takes under the matrix, normalized to
// [0, 2π): its own angle plus the angle of the matrix, mirrored about it where the given Scaling
// has a negative Y. The caller passes the factors it already scaled the shape by, so the matrix is
// read once.
func (m Matrix[T]) turnedAngle(angle float64, scaling Vector[T]) float64 {
	if scaling.Y < 0 {
		return NormalizeAngle(m.Angle() - angle)
	}

	return NormalizeAngle(angle + m.Angle())
}

// mapping returns the linear part of the matrix applied to the frame of an oriented shape, the unit
// circle stretched by the size and turned by the angle, so Ellipse.Transform and
// RegularPolygon.Transform decompose the one expression.
func (m Matrix[T]) mapping(angle float64, size Size[float64]) Matrix[float64] {
	f := m.Float()
	sin, cos := math.Sincos(angle)
	w, h := size.XY()

	return Matrix[float64]{f.A, f.B, 0, f.D, f.E, 0}.Multiply(Matrix[float64]{cos * w, -sin * h, 0, sin * w, cos * h, 0})
}

// decomposition returns the linear part as R(after)·diag(stretch)·R(before), the singular value
// decomposition, and whether a mirror of Y stands before the stretch, so the stretch is never
// negative. The turn before is kept within an eighth of a turn of zero, swapping the stretch for a
// quarter turn, so a turn and an axis-aligned scale report none, and it is zero for a uniform
// stretch. The lesser stretch is the determinant divided by the greater, which keeps it on a long
// thin shape where the difference of the halves cancels.
func (m Matrix[T]) decomposition() (float64, Vector[float64], float64, bool) {
	f := m.Float()
	e, g := float64((f.A+f.E)/2), float64((f.A-f.E)/2)
	h, k := float64((f.D+f.B)/2), float64((f.D-f.B)/2)

	q, r := math.Hypot(e, k), math.Hypot(g, h)
	anisotropy, rotation := math.Atan2(h, g), math.Atan2(k, e)
	if r == 0 {
		anisotropy = rotation
	}
	if q == 0 {
		rotation = anisotropy
	}

	after, before := float64((rotation+anisotropy)/2), float64((rotation-anisotropy)/2)
	stretch, reflected := Vector[float64]{q + r, 0}, q < r
	if q+r > 0 {
		stretch.Y = math.Abs(f.determinant()) / (q + r)
	}

	quarters := math.Round(before / (Pi / 2))
	before -= float64(quarters * (Pi / 2))
	if reflected {
		after -= float64(quarters * (Pi / 2))
	} else {
		after += float64(quarters * (Pi / 2))
	}
	if math.Mod(quarters, 2) != 0 {
		stretch = Vector[float64]{stretch.Y, stretch.X}
	}

	return after, stretch, before, reflected
}

// determinant returns the determinant in float64, so Inverse and IsInvertible judge an integer
// matrix on the exact value, and a singular matrix gives exactly zero on every platform.
func (m Matrix[T]) determinant() float64 {
	f := m.Float()

	return float64(f.A*f.E) - float64(f.B*f.D)
}

// Translate returns the matrix composed with TranslationMatrix on the right (m·T).
func (m Matrix[T]) Translate(deltaX, deltaY T) Matrix[T] {
	return m.Multiply(TranslationMatrix(deltaX, deltaY))
}

// Untranslate returns the matrix composed with TranslationMatrix of the negated deltas on the
// right.
func (m Matrix[T]) Untranslate(deltaX, deltaY T) Matrix[T] {
	return m.Multiply(TranslationMatrix(-deltaX, -deltaY))
}

// PreTranslate returns the matrix composed with TranslationMatrix on the left (T·m).
func (m Matrix[T]) PreTranslate(deltaX, deltaY T) Matrix[T] {
	return TranslationMatrix(deltaX, deltaY).Multiply(m)
}

// Rotate returns the matrix composed with RotationMatrix on the right (m·R).
func (m Matrix[T]) Rotate(angle float64) Matrix[T] {
	return m.Multiply(RotationMatrix[T](angle))
}

// RotateAround returns the matrix composed with RotationMatrixAround on the right (m·R).
func (m Matrix[T]) RotateAround(pivotX, pivotY T, angle float64) Matrix[T] {
	return m.Multiply(RotationMatrixAround(pivotX, pivotY, angle))
}

// PreRotate returns the matrix composed with RotationMatrix on the left (R·m).
func (m Matrix[T]) PreRotate(angle float64) Matrix[T] {
	return RotationMatrix[T](angle).Multiply(m)
}

// PreRotateAround returns the matrix composed with RotationMatrixAround on the left (R·m).
func (m Matrix[T]) PreRotateAround(pivotX, pivotY T, angle float64) Matrix[T] {
	return RotationMatrixAround(pivotX, pivotY, angle).Multiply(m)
}

// Scale returns the matrix composed with a scale by the factors on the right (m·S), each column
// multiplied by its factor.
func (m Matrix[T]) Scale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorY), m.C,
		Multiply(m.D, factorX), Multiply(m.E, factorY), m.F,
	}
}

// ScaleAround returns the matrix composed with a scale by the factors about the pivot on the right
// (m·S), so the pivot keeps its place before the current transform. Each entry is rounded once, and
// for float64 it equals the product with ScaleMatrixAround to the bit.
func (m Matrix[T]) ScaleAround(pivotX, pivotY T, factorX, factorY float64) Matrix[T] {
	f := m.Float()
	x, y := pivotOffset(pivotX, factorX), pivotOffset(pivotY, factorY)

	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorY), Cast[T](float64(f.A*x) + float64(f.B*y) + f.C),
		Multiply(m.D, factorX), Multiply(m.E, factorY), Cast[T](float64(f.D*x) + float64(f.E*y) + f.F),
	}
}

// Unscale returns the matrix composed with a scale by the inverse factors on the right, each column
// divided by its factor. It panics for a zero factor.
func (m Matrix[T]) Unscale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Divide(m.A, factorX), Divide(m.B, factorY), m.C,
		Divide(m.D, factorX), Divide(m.E, factorY), m.F,
	}
}

// PreScale returns the matrix composed with a scale by the factors on the left (S·m), each row
// multiplied by its factor.
func (m Matrix[T]) PreScale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorX), Multiply(m.C, factorX),
		Multiply(m.D, factorY), Multiply(m.E, factorY), Multiply(m.F, factorY),
	}
}

// PreScaleAround returns the matrix composed with a scale by the factors about the pivot on the
// left (S·m), so the pivot keeps its place after the current transform. Each entry is rounded once.
func (m Matrix[T]) PreScaleAround(pivotX, pivotY T, factorX, factorY float64) Matrix[T] {
	f := m.Float()
	x, y := pivotOffset(pivotX, factorX), pivotOffset(pivotY, factorY)

	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorX), Cast[T](float64(f.C*factorX) + x),
		Multiply(m.D, factorY), Multiply(m.E, factorY), Cast[T](float64(f.F*factorY) + y),
	}
}

// Shear returns the matrix composed with ShearMatrix on the right (m·H).
func (m Matrix[T]) Shear(shearX, shearY T) Matrix[T] {
	return m.Multiply(ShearMatrix(shearX, shearY))
}

// PreShear returns the matrix composed with ShearMatrix on the left (H·m).
func (m Matrix[T]) PreShear(shearX, shearY T) Matrix[T] {
	return ShearMatrix(shearX, shearY).Multiply(m)
}

// Reflect returns the matrix composed with ReflectionMatrix on the right (m·F), unchanged for
// AxisNone.
func (m Matrix[T]) Reflect(axis Axis) Matrix[T] {
	return m.Multiply(ReflectionMatrix[T](axis))
}

// PreReflect returns the matrix composed with ReflectionMatrix on the left (F·m).
func (m Matrix[T]) PreReflect(axis Axis) Matrix[T] {
	return ReflectionMatrix[T](axis).Multiply(m)
}

// Multiply returns the product of this matrix and the given one, the given one applied first.
func (m Matrix[T]) Multiply(matrix Matrix[T]) Matrix[T] {
	l, r := m.Float(), matrix.Float()

	return Matrix[T]{
		Cast[T](float64(l.A*r.A) + float64(l.B*r.D)),
		Cast[T](float64(l.A*r.B) + float64(l.B*r.E)),
		Cast[T](float64(l.A*r.C) + float64(l.B*r.F) + l.C),
		Cast[T](float64(l.D*r.A) + float64(l.E*r.D)),
		Cast[T](float64(l.D*r.B) + float64(l.E*r.E)),
		Cast[T](float64(l.D*r.C) + float64(l.E*r.F) + l.F),
	}
}

// Inverse returns the matrix that undoes this one. It panics for a singular matrix; check
// IsInvertible first. For integer T only a determinant of ±1 inverts exactly.
func (m Matrix[T]) Inverse() Matrix[T] {
	det := m.determinant()
	if det == 0 {
		panic(ErrSingularMatrix)
	}

	f := m.Float()

	return Matrix[T]{
		Cast[T](f.E / det),
		Cast[T](-f.B / det),
		Cast[T]((float64(f.B*f.F) - float64(f.C*f.E)) / det),
		Cast[T](-f.D / det),
		Cast[T](f.A / det),
		Cast[T]((float64(f.C*f.D) - float64(f.A*f.F)) / det),
	}
}

// pivotOffset returns the translation of a scale about the pivot along one axis, the pivot less its
// scaled self, in the order the product with ScaleMatrixAround rounds it.
func pivotOffset[T Number](pivot T, factor float64) float64 {
	p := float64(pivot)

	return float64(factor*-p) + p
}

// Equal reports whether the matrices have equal components.
func (m Matrix[T]) Equal(matrix Matrix[T]) bool {
	return Equal(m.A, matrix.A) && Equal(m.B, matrix.B) && Equal(m.C, matrix.C) && Equal(m.D, matrix.D) && Equal(m.E, matrix.E) && Equal(m.F, matrix.F)
}

// IsZero reports whether every component is zero.
func (m Matrix[T]) IsZero() bool {
	return m.Equal(Matrix[T]{})
}

// IsIdentity reports whether the matrix equals IdentityMatrix.
func (m Matrix[T]) IsIdentity() bool {
	return m.Equal(IdentityMatrix[T]())
}

// IsInvertible reports whether the determinant is not exactly zero. No tolerance applies: a tiny
// determinant means a large inverse, not a missing one.
func (m Matrix[T]) IsInvertible() bool {
	return m.determinant() != 0
}

// Cast converts the matrix to a Matrix of another number type, rounding as Cast does.
func (m Matrix[T]) Cast[R Number]() Matrix[R] {
	return Matrix[R]{
		Cast[R](float64(m.A)), Cast[R](float64(m.B)), Cast[R](float64(m.C)),
		Cast[R](float64(m.D)), Cast[R](float64(m.E)), Cast[R](float64(m.F)),
	}
}

// Int converts the matrix to a Matrix[int].
func (m Matrix[T]) Int() Matrix[int] {
	return Matrix[int]{
		Int(m.A), Int(m.B), Int(m.C),
		Int(m.D), Int(m.E), Int(m.F),
	}
}

// Float converts the matrix to a Matrix[float64].
func (m Matrix[T]) Float() Matrix[float64] {
	return Matrix[float64]{
		float64(m.A), float64(m.B), float64(m.C),
		float64(m.D), float64(m.E), float64(m.F),
	}
}

// String returns the matrix as its two rows, with no type name: [[a, b, c], [d, e, f]].
func (m Matrix[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "[["...)
	buf = appendString(buf, m.A)
	buf = append(buf, ", "...)
	buf = appendString(buf, m.B)
	buf = append(buf, ", "...)
	buf = appendString(buf, m.C)
	buf = append(buf, "], ["...)
	buf = appendString(buf, m.D)
	buf = append(buf, ", "...)
	buf = appendString(buf, m.E)
	buf = append(buf, ", "...)
	buf = appendString(buf, m.F)
	buf = append(buf, "]]"...)

	return string(buf)
}

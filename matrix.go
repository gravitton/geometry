package geom

import (
	"fmt"
	"math"
)

// Matrix is a 2D affine matrix.
//
// Translate, Rotate, Scale, Shear and Reflect compose the matrix with the constructor of the
// same name, TranslationMatrix, RotationMatrix, ScaleMatrix, ShearMatrix and ReflectionMatrix,
// and the Pre variants compose on the other side. Translate, Rotate, Shear and Reflect take the
// arguments of their constructor. Scale and Unscale take float64 factors like every other
// Scale in the package, where ScaleMatrix takes them in T, so an integer matrix can be scaled
// by a half; each scaled component is rounded, following Multiply.
//
// For integer T, the matrix stores and composes lattice transforms – translation, integer scaling,
// reflection, quarter turns – exactly. The rest are not closed over the integers and round into a
// different matrix: rotation by anything but a multiple of 90° (see RotationMatrix), Inverse, and
// Unscale.
//
// Multiply, Determinant and Inverse compute in float64 and round the result into T, so a
// narrow integer T does not overflow mid-computation; only a result outside its range is lost.
type Matrix[T Number] struct {
	A T `json:"a"` // scale X
	B T `json:"b"` // shear X (contribution of y to x')
	C T `json:"c"` // translate X
	D T `json:"d"` // shear Y (contribution of x to y')
	E T `json:"e"` // scale Y
	F T `json:"f"` // translate Y
	// [0 0 1] implicit third row: x' = A*x + B*y + C, y' = D*x + E*y + F
}

// Mat is shorthand for Matrix{a, b, c, d, e, f}.
func Mat[T Number](a, b, c, d, e, f T) Matrix[T] {
	return Matrix[T]{a, b, c, d, e, f}
}

// IdentityMatrix creates a new identity matrix.
func IdentityMatrix[T Number]() Matrix[T] {
	return Matrix[T]{
		1, 0, 0,
		0, 1, 0,
	}
}

// TranslationMatrix creates a new translation matrix.
func TranslationMatrix[T Number](deltaX, deltaY T) Matrix[T] {
	return Matrix[T]{
		1, 0, deltaX,
		0, 1, deltaY,
	}
}

// RotationMatrix creates a new rotation matrix (angle in radians).
// For integer T, sin/cos components are rounded; only multiples of 90° give exact results.
// Any other angle is not a rotation at all: the rounded entries scale and shear as well as turn.
func RotationMatrix[T Number](angle float64) Matrix[T] {
	sin, cos := math.Sincos(angle)

	return Matrix[T]{
		Cast[T](cos), Cast[T](-sin), 0,
		Cast[T](sin), Cast[T](cos), 0,
	}
}

// ScaleMatrix creates a new scale matrix.
func ScaleMatrix[T Number](factorX, factorY T) Matrix[T] {
	return Matrix[T]{
		factorX, 0, 0,
		0, factorY, 0,
	}
}

// ShearMatrix creates a new shear matrix: shearX is the contribution of y to x', shearY the
// contribution of x to y'.
func ShearMatrix[T Number](shearX, shearY T) Matrix[T] {
	return Matrix[T]{
		1, shearX, 0,
		shearY, 1, 0,
	}
}

// ReflectionMatrix creates a new matrix reflecting across the given axis: across the
// horizontal axis flips Y, across the vertical axis flips X. AxisNone reflects across nothing
// and gives the identity.
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

// Determinant calculates the determinant of the 2x2 matrix.
func (m Matrix[T]) Determinant() T {
	return Cast[T](m.determinant())
}

// Translation returns the translation the matrix applies, its C and F components.
func (m Matrix[T]) Translation() Vector[T] {
	return Vector[T]{m.C, m.F}
}

// Angle returns the rotation the matrix applies, in radians: the angle of the transformed X
// axis, so a matrix built from a rotation, a scale and a translation gives the rotation back.
// A sheared matrix has no single rotation and gives the angle of its X axis. The zero matrix
// gives 0.
func (m Matrix[T]) Angle() float64 {
	return math.Atan2(float64(m.D), float64(m.A))
}

// Scaling returns the scale factors the matrix applies along its rotated X and Y axes: the
// length of the transformed X axis, and the signed length of the Y axis, negative for a
// reflection. A matrix built from a rotation, a scale and a translation gives the scale back;
// a sheared matrix gives the length of its X axis and the determinant divided by it, the
// factors that keep that axis and the area. For integer T the factors are rounded; the zero
// matrix gives the zero vector.
func (m Matrix[T]) Scaling() Vector[T] {
	f := m.Float()

	x := math.Hypot(f.A, f.D)
	if x == 0 {
		return Vector[T]{}
	}

	return Vector[T]{Cast[T](x), Cast[T](f.determinant() / x)}
}

// determinant calculates the determinant in float64, the form Inverse and IsInvertible use
// so that an integer matrix is judged on its exact determinant, not a rounded one.
// The two products are rounded separately, which keeps a fused multiply-add from turning the
// determinant of a singular matrix into a rounding residue: it is exactly zero on every platform.
func (m Matrix[T]) determinant() float64 {
	f := m.Float()

	return float64(f.A*f.E) - float64(f.B*f.D)
}

// turnedAngle returns the angle a shape at the given angle takes under the matrix: its own
// angle plus the angle of the matrix, or mirrored about that angle where the matrix reflects,
// which the negative Y factor of Scaling names. It is normalized to [0, 2π) like Rotate, and
// is the angle Rectangle.Transform places, the nearest rectangle where the matrix shears it.
// The factors are the ones the caller already scaled the shape by, rather than read from the
// matrix a second time.
func (m Matrix[T]) turnedAngle(angle float64, scaling Vector[T]) float64 {
	if scaling.Y < 0 {
		return NormalizeAngle(m.Angle() - angle)
	}

	return NormalizeAngle(angle + m.Angle())
}

// decomposition returns the linear part of the matrix as a turn after, a stretch along the axes
// and a turn before, its singular value decomposition R(after) · diag(stretch) · R(before), and
// whether a mirror of Y stands between the stretch and the turn before, which a reflection
// needs, so the stretch is never negative. The decomposition is not unique, and the turn before
// is taken within an eighth of a turn of zero, a quarter turn moved into the turn after with
// the stretch swapped, so a matrix that turns and scales along the axes reports none; where
// the stretch is the same on both axes, any turn before would do, and it is zero. The lesser
// stretch is the determinant divided by the greater, its product with it, rather than the
// difference of the two halves, which loses the lesser to cancellation on a long thin shape.
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

// mapping returns the linear part of the matrix applied to the frame of an oriented shape, the
// unit circle stretched by the size and turned by the angle: the map whose decomposition gives
// the turn, the semi-axes and, for a regular polygon, the phase the shape takes under the matrix,
// so Ellipse.Transform and RegularPolygon.Transform read the one expression.
func (m Matrix[T]) mapping(angle float64, size Size[float64]) Matrix[float64] {
	f := m.Float()
	sin, cos := math.Sincos(angle)
	w, h := size.XY()

	return Matrix[float64]{f.A, f.B, 0, f.D, f.E, 0}.Multiply(Matrix[float64]{cos * w, -sin * h, 0, sin * w, cos * h, 0})
}

// Multiply creates a new matrix by multiplying the current matrix with given matrix.
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

// Inverse creates a new inverse affine matrix. A singular matrix has no inverse and Inverse
// panics for one, the same convention Divide follows for a zero factor; check IsInvertible first
// when the matrix may be singular.
// For integer T, all six components are rounded; only |det| = 1 gives exact results, which covers
// translations, reflections and quarter turns. Otherwise the inverse does not undo the matrix:
// the inverse of an integer scale rounds its fractional factor to a whole one.
func (m Matrix[T]) Inverse() Matrix[T] {
	det := m.determinant()
	if det == 0 {
		panic("geom: inverse of a singular matrix")
	}

	f := m.Float()
	invDet := 1 / det

	return Matrix[T]{
		Cast[T](f.E * invDet),
		Cast[T](-f.B * invDet),
		Cast[T]((float64(f.B*f.F) - float64(f.C*f.E)) * invDet),
		Cast[T](-f.D * invDet),
		Cast[T](f.A * invDet),
		Cast[T]((float64(f.C*f.D) - float64(f.A*f.F)) * invDet),
	}
}

// Translate creates a new matrix by right-multiplying a translation matrix.
// Composition order: result = m * m_T(deltaX,deltaY).
func (m Matrix[T]) Translate(deltaX, deltaY T) Matrix[T] {
	return m.Multiply(TranslationMatrix(deltaX, deltaY))
}

// Untranslate creates a new matrix by right-multiplying a translation matrix with negative deltas.
// Composition order: result = m * m_T(-deltaX,-deltaY).
func (m Matrix[T]) Untranslate(deltaX, deltaY T) Matrix[T] {
	return m.Multiply(TranslationMatrix(-deltaX, -deltaY))
}

// PreTranslate creates a new matrix by left-multiplying a translation matrix.
// Composition order: result = m_T(deltaX,deltaY) * m.
func (m Matrix[T]) PreTranslate(deltaX, deltaY T) Matrix[T] {
	return TranslationMatrix(deltaX, deltaY).Multiply(m)
}

// Rotate creates a new matrix by right-multiplying a rotation matrix (angle in radians).
// Composition order: result = m * m_R(angle).
// For integer T, see RotationMatrix: only multiples of 90° give exact results.
func (m Matrix[T]) Rotate(angle float64) Matrix[T] {
	return m.Multiply(RotationMatrix[T](angle))
}

// PreRotate creates a new matrix by left-multiplying a rotation matrix (angle in radians).
// Composition order: result = m_R(angle) * m.
// For integer T, see RotationMatrix: only multiples of 90° give exact results.
func (m Matrix[T]) PreRotate(angle float64) Matrix[T] {
	return RotationMatrix[T](angle).Multiply(m)
}

// Scale creates a new matrix equal to right-multiplying a scale matrix.
// Composition order: result = m * m_S(factorX,factorY).
// The factors are float64 like every other Scale in the package; for integer T each scaled
// component is rounded, following Multiply.
func (m Matrix[T]) Scale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorY), m.C,
		Multiply(m.D, factorX), Multiply(m.E, factorY), m.F,
	}
}

// Unscale creates a new matrix equal to right-multiplying a scale matrix with inverse factors.
// Composition order: result = m * m_S(1/factorX,1/factorY).
// Each column is divided on its own, following Divide, which panics for a zero factor.
// For integer T, each component is rounded, so the result is exact when the components of the
// column are multiples of its factor.
func (m Matrix[T]) Unscale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Divide(m.A, factorX), Divide(m.B, factorY), m.C,
		Divide(m.D, factorX), Divide(m.E, factorY), m.F,
	}
}

// PreScale creates a new matrix equal to left-multiplying a scale matrix.
// Composition order: result = m_S(factorX,factorY) * m.
// The factors are float64 like Scale; for integer T each scaled component is rounded.
func (m Matrix[T]) PreScale(factorX, factorY float64) Matrix[T] {
	return Matrix[T]{
		Multiply(m.A, factorX), Multiply(m.B, factorX), Multiply(m.C, factorX),
		Multiply(m.D, factorY), Multiply(m.E, factorY), Multiply(m.F, factorY),
	}
}

// Shear creates a new matrix by right-multiplying a shear matrix.
// Composition order: result = m * m_H(shearX,shearY).
func (m Matrix[T]) Shear(shearX, shearY T) Matrix[T] {
	return m.Multiply(ShearMatrix(shearX, shearY))
}

// PreShear creates a new matrix by left-multiplying a shear matrix.
// Composition order: result = m_H(shearX,shearY) * m.
func (m Matrix[T]) PreShear(shearX, shearY T) Matrix[T] {
	return ShearMatrix(shearX, shearY).Multiply(m)
}

// Reflect creates a new matrix by right-multiplying a reflection matrix.
// Composition order: result = m * m_F(axis). Reflecting across AxisNone leaves the matrix
// unchanged, since ReflectionMatrix gives the identity for it.
func (m Matrix[T]) Reflect(axis Axis) Matrix[T] {
	return m.Multiply(ReflectionMatrix[T](axis))
}

// PreReflect creates a new matrix by left-multiplying a reflection matrix.
// Composition order: result = m_F(axis) * m.
func (m Matrix[T]) PreReflect(axis Axis) Matrix[T] {
	return ReflectionMatrix[T](axis).Multiply(m)
}

// Equal checks for equal values.
func (m Matrix[T]) Equal(matrix Matrix[T]) bool {
	return Equal(m.A, matrix.A) && Equal(m.B, matrix.B) && Equal(m.C, matrix.C) && Equal(m.D, matrix.D) && Equal(m.E, matrix.E) && Equal(m.F, matrix.F)
}

// IsZero checks if values are zero.
func (m Matrix[T]) IsZero() bool {
	return m.Equal(Matrix[T]{})
}

// IsIdentity checks whether the matrix is the identity, the transform that leaves every point
// where it is, comparing like Equal.
func (m Matrix[T]) IsIdentity() bool {
	return m.Equal(IdentityMatrix[T]())
}

// IsInvertible reports whether the matrix has an inverse: its determinant is not exactly zero.
// No tolerance is applied, since a determinant scales with the square of the matrix and a small
// one only means a large inverse, not a missing one: a tiny uniform scale is invertible.
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

// Int converts the matrix to a Matrix[int], rounding each component.
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

// String returns a string representation of the Matrix[T].
func (m Matrix[T]) String() string {
	return fmt.Sprintf("[[%s, %s, %s], [%s, %s, %s]]", String(m.A), String(m.B), String(m.C), String(m.D), String(m.E), String(m.F))
}

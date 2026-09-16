package geom

import (
	"fmt"
	"math"
)

// Matrix is a 2D affine matrix.
//
// For integer T, the matrix stores and composes lattice transforms – translation, integer scaling,
// reflection, quarter turns – exactly. The rest are not closed over the integers and round into a
// different matrix: rotation by anything but a multiple of 90° (see RotationMatrix), Inverse, and
// Unscale.
//
// Multiply, Determinant and Inverse compute in T. They are meant for int and int64; a narrow
// integer T such as int8 or int16 overflows in ordinary products and is not supported.
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
// Any other angle is not a rotation at all: π/6 rounds to [[1, -1], [1, 1]], scaling by √2 and shearing.
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

// Multiply creates a new matrix by multiplying the current matrix with given matrix.
func (m Matrix[T]) Multiply(matrix Matrix[T]) Matrix[T] {
	return Matrix[T]{
		m.A*matrix.A + m.B*matrix.D,
		m.A*matrix.B + m.B*matrix.E,
		m.A*matrix.C + m.B*matrix.F + m.C,
		m.D*matrix.A + m.E*matrix.D,
		m.D*matrix.B + m.E*matrix.E,
		m.D*matrix.C + m.E*matrix.F + m.F,
	}
}

// Inverse creates a new inverse affine matrix. A singular matrix has no inverse and is returned
// unchanged, the same convention Divide follows for a zero scale; check IsInvertible first when
// that matters.
// For integer T, all six components are rounded; only |det| = 1 gives exact results, which covers
// translations, reflections and quarter turns. Otherwise the inverse does not undo the matrix:
// ScaleMatrix(2, 2).Inverse() rounds 0.5 back up to identity.
func (m Matrix[T]) Inverse() Matrix[T] {
	if !m.IsInvertible() {
		return m
	}

	det := m.Determinant()

	a, b, c := float64(m.A), float64(m.B), float64(m.C)
	d, e, f := float64(m.D), float64(m.E), float64(m.F)
	invDet := 1.0 / float64(det)

	return Matrix[T]{
		Cast[T](e * invDet),
		Cast[T](-b * invDet),
		Cast[T]((b*f - c*e) * invDet),
		Cast[T](-d * invDet),
		Cast[T](a * invDet),
		Cast[T]((c*d - a*f) * invDet),
	}
}

// IsInvertible reports whether the matrix has an inverse: its determinant is not zero
// (within Epsilon of T).
func (m Matrix[T]) IsInvertible() bool {
	return !Equal(m.Determinant(), 0.0)
}

// Determinant calculates the determinant of the 2x2 matrix.
func (m Matrix[T]) Determinant() T {
	return m.A*m.E - m.B*m.D
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

// Scale creates a new matrix by right-multiplying a scale matrix.
// Composition order: result = m * m_S(factorX,factorY).
func (m Matrix[T]) Scale(factorX, factorY T) Matrix[T] {
	return m.Multiply(ScaleMatrix(factorX, factorY))
}

// Unscale creates a new matrix by right-multiplying a scale matrix with inverse factors.
// Composition order: result = m * m_S(1/factorX,1/factorY).
// Each axis is inverted on its own, following Divide: a zero factor leaves that axis unchanged.
// For integer T, each inverse factor is rounded; only ±1 gives exact results.
func (m Matrix[T]) Unscale(factorX, factorY T) Matrix[T] {
	return m.Multiply(ScaleMatrix(Divide[T](1, float64(factorX)), Divide[T](1, float64(factorY))))
}

// PreScale creates a new matrix by left-multiplying a scale matrix.
// Composition order: result = m_S(factorX,factorY) * m.
func (m Matrix[T]) PreScale(factorX, factorY T) Matrix[T] {
	return ScaleMatrix(factorX, factorY).Multiply(m)
}

// Equal checks for equal values.
func (m Matrix[T]) Equal(matrix Matrix[T]) bool {
	return Equal(m.A, matrix.A) && Equal(m.B, matrix.B) && Equal(m.C, matrix.C) && Equal(m.D, matrix.D) && Equal(m.E, matrix.E) && Equal(m.F, matrix.F)
}

// IsZero checks if values are zero.
func (m Matrix[T]) IsZero() bool {
	return m.Equal(Matrix[T]{})
}

// Int converts the matrix to a Matrix[int], rounding each component.
func (m Matrix[T]) Int() Matrix[int] {
	return Matrix[int]{
		Cast[int](float64(m.A)), Cast[int](float64(m.B)), Cast[int](float64(m.C)),
		Cast[int](float64(m.D)), Cast[int](float64(m.E)), Cast[int](float64(m.F)),
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

package geom

import (
	"fmt"
	"math"
)

// Matrix is a 2D matrix.
type Matrix[T Number] struct {
	A T `json:"a"` // scale X
	B T `json:"b"` // shear Y
	C T `json:"c"` // translate X
	D T `json:"d"` // shear X
	E T `json:"e"` // scale Y
	F T `json:"f"` // translate Y
	// [0 0 1] implicit third row
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

// RotationMatrix creates a new rotation matrix.
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

// Inverse creates a new inverse affine matrix. If non-invertible (det ~ 0), returns the same matrix.
func (m Matrix[T]) Inverse() Matrix[T] {
	det := m.Determinant()
	if Equal(det, 0.0) {
		return m
	}

	invDet := 1.0 / float64(det)
	return Matrix[T]{
		Cast[T](float64(m.E) * invDet),
		Cast[T](float64(-m.B) * invDet),
		Cast[T](float64(m.B*m.F-m.C*m.E) * invDet),
		Cast[T](float64(-m.D) * invDet),
		Cast[T](float64(m.A) * invDet),
		Cast[T](float64(m.C*m.D-m.A*m.F) * invDet),
	}
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
func (m Matrix[T]) Rotate(angle float64) Matrix[T] {
	return m.Multiply(RotationMatrix[T](angle))
}

// PreRotate creates a new matrix by left-multiplying a rotation matrix (angle in radians).
// Composition order: result = m_R(angle) * m.
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
func (m Matrix[T]) Unscale(factorX, factorY T) Matrix[T] {
	if factorX == 0 && factorY == 0 {
		return m
	}

	return m.Multiply(ScaleMatrix[T](1.0/factorX, 1.0/factorY))
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

// Int converts the matrix to a Matrix[int].
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

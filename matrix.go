package geom

import (
	"math"
)

// Matrix is a 2D matrix.
type Matrix struct {
	A, B, C float64 // scale X, shear Y, translate X
	D, E, F float64 // shear X, scale Y, translate Y
	// [0 0 1] implicit third row
}

// M is shorthand for Matrix{a, b, c, d, e, f}.
func M(a, b, c, d, e, f float64) Matrix {
	return Matrix{a, b, c, d, e, f}
}

// Multiply creates a new matrix by multiplying the current matrix with given matrix.
func (m Matrix) Multiply(matrix Matrix) Matrix {
	return Matrix{
		m.A*matrix.A + m.B*matrix.D,
		m.A*matrix.B + m.B*matrix.E,
		m.A*matrix.C + m.B*matrix.F + m.C,
		m.D*matrix.A + m.E*matrix.D,
		m.D*matrix.B + m.E*matrix.E,
		m.D*matrix.C + m.E*matrix.F + m.F,
	}
}

// Inverse creates a new inverse affine matrix. If non-invertible (det ~ 0), returns the same matrix.
func (m Matrix) Inverse() Matrix {
	det := m.Determinant()
	if Equal(det, 0.0) {
		return m
	}

	invDet := 1.0 / det
	return Matrix{
		m.E * invDet,
		-m.B * invDet,
		(m.B*m.F - m.C*m.E) * invDet,
		-m.D * invDet,
		m.A * invDet,
		(m.C*m.D - m.A*m.F) * invDet,
	}
}

// Determinant calculates the determinant of the 2x2 matrix.
func (m Matrix) Determinant() float64 {
	return m.A*m.E - m.B*m.D
}

// Translate creates a new matrix by translating the current matrix.
func (m Matrix) Translate(deltaX, deltaY float64) Matrix {
	return m.Multiply(TranslationMatrix(deltaX, deltaY))
}

// Rotate creates a new matrix by rotating the current matrix.
func (m Matrix) Rotate(angle float64) Matrix {
	return m.Multiply(RotationMatrix(angle))
}

// Scale creates a new matrix by scaling the current matrix.
func (m Matrix) Scale(factorX, factorY float64) Matrix {
	return m.Multiply(ScaleMatrix(factorX, factorY))
}

// Equal checks for equal values.
func (m Matrix) Equal(other Matrix) bool {
	return Equal(m.A, other.A) && Equal(m.B, other.B) && Equal(m.C, other.C) && Equal(m.D, other.D) && Equal(m.E, other.E) && Equal(m.F, other.F)
}

// IsZero checks if values are zero.
func (m Matrix) IsZero() bool { return m.Equal(Matrix{}) }

// IdentityMatrix creates a new identity matrix.
func IdentityMatrix() Matrix {
	return Matrix{
		1, 0, 0,
		0, 1, 0,
	}
}

// TranslationMatrix creates a new translation matrix.
func TranslationMatrix(deltaX, deltaY float64) Matrix {
	return Matrix{
		1, 0, deltaX,
		0, 1, deltaY,
	}
}

// RotationMatrix creates a new rotation matrix.
func RotationMatrix(angle float64) Matrix {
	sin, cos := math.Sincos(angle)
	return Matrix{
		cos, -sin, 0,
		sin, cos, 0,
	}
}

// ScaleMatrix creates a new scale matrix.
func ScaleMatrix(factorX, factorY float64) Matrix {
	return Matrix{
		factorX, 0, 0,
		0, factorY, 0,
	}
}

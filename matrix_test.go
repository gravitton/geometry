package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestMatrix_New(t *testing.T) {
	AssertMatrix(t, IdentityMatrix(), Mat(1, 0, 0, 0, 1, 0))
	AssertMatrix(t, TranslationMatrix(5, 3), Mat(1, 0, 5, 0, 1, 3))
	AssertMatrix(t, ScaleMatrix(2, 3), Mat(2, 0, 0, 0, 3, 0))
	AssertMatrix(t, RotationMatrix(0), Mat(1, 0, 0, 0, 1, 0))

	rot90 := RotationMatrix(Pi / 2)
	assert.EqualDelta(t, rot90.A, 0.0, Delta)
	assert.EqualDelta(t, rot90.B, -1.0, Delta)
	assert.EqualDelta(t, rot90.D, 1.0, Delta)
	assert.EqualDelta(t, rot90.E, 0.0, Delta)
}

func TestMatrix_Multiply(t *testing.T) {
	// identity is neutral
	AssertMatrix(t, IdentityMatrix().Multiply(TranslationMatrix(5, 3)), TranslationMatrix(5, 3))
	AssertMatrix(t, TranslationMatrix(5, 3).Multiply(IdentityMatrix()), TranslationMatrix(5, 3))

	// two translations compose additively
	AssertMatrix(t, TranslationMatrix(5, 3).Multiply(TranslationMatrix(2, 1)), TranslationMatrix(7, 4))

	// two scales compose multiplicatively
	AssertMatrix(t, ScaleMatrix(2, 3).Multiply(ScaleMatrix(4, 2)), ScaleMatrix(8, 6))
}

func TestMatrix_Determinant(t *testing.T) {
	assert.EqualDelta(t, IdentityMatrix().Determinant(), 1.0, Delta)
	assert.EqualDelta(t, ScaleMatrix(2, 3).Determinant(), 6.0, Delta)
	assert.EqualDelta(t, RotationMatrix(Pi/4).Determinant(), 1.0, Delta)
	assert.EqualDelta(t, Mat(1, 2, 0, 3, 4, 0).Determinant(), -2.0, Delta) // 1*4 - 2*3
}

func TestMatrix_Inverse(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().Inverse(), IdentityMatrix())
	AssertMatrix(t, TranslationMatrix(5, 3).Inverse(), TranslationMatrix(-5, -3))
	AssertMatrix(t, ScaleMatrix(2, 4).Inverse(), ScaleMatrix(0.5, 0.25))

	// non-invertible matrix returns itself
	zero := Mat(0, 0, 0, 0, 0, 0)
	AssertMatrix(t, zero.Inverse(), zero)
}

func TestMatrix_Translate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().Translate(5, 3), TranslationMatrix(5, 3))
	AssertMatrix(t, IdentityMatrix().Translate(2, 1).Translate(3, 2), TranslationMatrix(5, 3))
}

func TestMatrix_Untranslate(t *testing.T) {
	AssertMatrix(t, TranslationMatrix(5, 3).Untranslate(5, 3), IdentityMatrix())
	AssertMatrix(t, TranslationMatrix(5, 3).Untranslate(2, 1), TranslationMatrix(3, 2))
}

func TestMatrix_PreTranslate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().PreTranslate(5, 3), TranslationMatrix(5, 3))

	// S * PreTranslate(tx,ty) = T(tx,ty) * S  (translation applied before scale)
	got := ScaleMatrix(2, 2).PreTranslate(5, 3)
	AssertMatrix(t, got, Mat(2, 0, 5, 0, 2, 3))
}

func TestMatrix_Rotate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().Rotate(0), IdentityMatrix())

	// point (1,0) rotated 90° → (0,1)
	p := Pt(1.0, 0.0).Transform(IdentityMatrix().Rotate(Pi / 2))
	assert.EqualDelta(t, p.X, 0.0, Delta)
	assert.EqualDelta(t, p.Y, 1.0, Delta)
}

func TestMatrix_PreRotate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().PreRotate(0), IdentityMatrix())

	// PreRotate(θ) * T differs from T * Rotate(θ)
	m1 := TranslationMatrix(5, 0).PreRotate(Pi / 2)
	m2 := TranslationMatrix(5, 0).Rotate(Pi / 2)
	assert.False(t, m1.Equal(m2))
}

func TestMatrix_Scale(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().Scale(2, 3), ScaleMatrix(2, 3))
	AssertMatrix(t, ScaleMatrix(2, 2).Scale(3, 3), ScaleMatrix(6, 6))
}

func TestMatrix_Unscale(t *testing.T) {
	AssertMatrix(t, ScaleMatrix(2, 3).Unscale(2, 3), IdentityMatrix())
	AssertMatrix(t, ScaleMatrix(4, 6).Unscale(2, 3), ScaleMatrix(2, 2))

	// zero guard: no change
	m := ScaleMatrix(2, 3)
	AssertMatrix(t, m.Unscale(0, 0), m)
}

func TestMatrix_PreScale(t *testing.T) {
	AssertMatrix(t, IdentityMatrix().PreScale(2, 3), ScaleMatrix(2, 3))

	// S(2,2) * T(5,5) → point(1,2) → (12,14)
	m := TranslationMatrix(5, 5).PreScale(2, 2)
	AssertMatrix(t, m, Mat(2, 0, 10, 0, 2, 10))

	p := Pt(1.0, 2.0).Transform(m)
	assert.EqualDelta(t, p.X, 12.0, Delta)
	assert.EqualDelta(t, p.Y, 14.0, Delta)
}

func TestMatrix_Equal(t *testing.T) {
	assert.True(t, IdentityMatrix().Equal(IdentityMatrix()))
	assert.False(t, IdentityMatrix().Equal(ScaleMatrix(2, 2)))
	assert.True(t, Mat(1, 2, 3, 4, 5, 6).Equal(Mat(1, 2, 3, 4, 5, 6)))
}

func TestMatrix_IsZero(t *testing.T) {
	assert.True(t, Matrix{}.IsZero())
	assert.False(t, IdentityMatrix().IsZero())
}

func TestMatrix_String(t *testing.T) {
	assert.Equal(t, IdentityMatrix().String(), "[[1.00, 0.00, 0.00], [0.00, 1.00, 0.00]]")
	assert.Equal(t, TranslationMatrix(5, 3).String(), "[[1.00, 0.00, 5.00], [0.00, 1.00, 3.00]]")
}

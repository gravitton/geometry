package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestMatrix_Constructor(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64](), Mat(1.0, 0.0, 0.0, 0.0, 1.0, 0.0))
	AssertMatrix(t, IdentityMatrix[float32](), Mat(float32(1), float32(0), float32(0), float32(0), float32(1), float32(0)))
	AssertMatrix(t, IdentityMatrix[int](), Mat(1, 0, 0, 0, 1, 0))
	AssertMatrix(t, TranslationMatrix(5.0, 3.0), Mat(1.0, 0.0, 5.0, 0.0, 1.0, 3.0))
	AssertMatrix(t, TranslationMatrix(5, 3), Mat(1, 0, 5, 0, 1, 3))
	AssertMatrix(t, ScaleMatrix(2.0, 3.0), Mat(2.0, 0.0, 0.0, 0.0, 3.0, 0.0))
	AssertMatrix(t, ScaleMatrix(2, 3), Mat(2, 0, 0, 0, 3, 0))
	AssertMatrix(t, RotationMatrix[float64](0.0), Mat(1.0, 0.0, 0.0, 0.0, 1.0, 0.0))

	rot90 := RotationMatrix[float64](Pi / 2.0)
	assert.EqualDelta(t, rot90.A, 0.0, Delta)
	assert.EqualDelta(t, rot90.B, -1.0, Delta)
	assert.EqualDelta(t, rot90.D, 1.0, Delta)
	assert.EqualDelta(t, rot90.E, 0.0, Delta)

	// cos(π/2) ≈ 6e-17 rounds to 0; sin(π/2) = 1 exactly
	AssertMatrix(t, RotationMatrix[int](Pi/2), Mat(0, -1, 0, 1, 0, 0))
}

func TestMatrix_Multiply(t *testing.T) {
	// identity is neutral
	AssertMatrix(t, IdentityMatrix[float64]().Multiply(TranslationMatrix(5.0, 3.0)), TranslationMatrix(5.0, 3.0))
	AssertMatrix(t, TranslationMatrix(5.0, 3.0).Multiply(IdentityMatrix[float64]()), TranslationMatrix(5.0, 3.0))

	// two translations compose additively
	AssertMatrix(t, TranslationMatrix(5.0, 3.0).Multiply(TranslationMatrix(2.0, 1.0)), TranslationMatrix(7.0, 4.0))

	// two scales compose multiplicatively
	AssertMatrix(t, ScaleMatrix(2.0, 3.0).Multiply(ScaleMatrix(4.0, 2.0)), ScaleMatrix(8.0, 6.0))

	// integer types
	AssertMatrix(t, IdentityMatrix[int]().Multiply(TranslationMatrix(5, 3)), TranslationMatrix(5, 3))
	AssertMatrix(t, TranslationMatrix(5, 3).Multiply(TranslationMatrix(2, 1)), TranslationMatrix(7, 4))
	AssertMatrix(t, ScaleMatrix(2, 3).Multiply(ScaleMatrix(4, 2)), ScaleMatrix(8, 6))
}

func TestMatrix_Determinant(t *testing.T) {
	assert.EqualDelta(t, IdentityMatrix[float64]().Determinant(), 1.0, Delta)
	assert.EqualDelta(t, ScaleMatrix(2.0, 3.0).Determinant(), 6.0, Delta)
	assert.EqualDelta(t, RotationMatrix[float64](Pi/4).Determinant(), 1.0, Delta)
	assert.EqualDelta(t, Mat(1.0, 2.0, 0.0, 3.0, 4.0, 0.0).Determinant(), -2.0, Delta) // 1*4 - 2*3

	assert.Equal(t, IdentityMatrix[int]().Determinant(), 1)
	assert.Equal(t, ScaleMatrix(2, 3).Determinant(), 6)
	assert.Equal(t, Mat(1, 2, 0, 3, 4, 0).Determinant(), -2) // 1*4 - 2*3
}

func TestMatrix_Inverse(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().Inverse(), IdentityMatrix[float64]())
	AssertMatrix(t, TranslationMatrix(5.0, 3.0).Inverse(), TranslationMatrix(-5.0, -3.0))
	AssertMatrix(t, ScaleMatrix(2.0, 4.0).Inverse(), ScaleMatrix(0.5, 0.25))

	// non-invertible matrix returns itself
	zero := Mat(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
	AssertMatrix(t, zero.Inverse(), zero)

	// integer: inverse uses Cast[int] which rounds; det=1 → exact
	AssertMatrix(t, IdentityMatrix[int]().Inverse(), IdentityMatrix[int]())
	AssertMatrix(t, TranslationMatrix(5, 3).Inverse(), TranslationMatrix(-5, -3))

	// det=4, invDet=0.25 → scale component Cast[int](2*0.25)=Cast[int](0.5)=1 (rounds)
	AssertMatrix(t, ScaleMatrix(2, 2).Inverse(), IdentityMatrix[int]())
}

func TestMatrix_Translate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().Translate(5.0, 3.0), TranslationMatrix(5.0, 3.0))
	AssertMatrix(t, IdentityMatrix[float64]().Translate(2.0, 1.0).Translate(3.0, 2.0), TranslationMatrix(5.0, 3.0))

	AssertMatrix(t, IdentityMatrix[int]().Translate(5, 3), TranslationMatrix(5, 3))
	AssertMatrix(t, IdentityMatrix[int]().Translate(2, 1).Translate(3, 2), TranslationMatrix(5, 3))
}

func TestMatrix_Untranslate(t *testing.T) {
	AssertMatrix(t, TranslationMatrix(5.0, 3.0).Untranslate(5.0, 3.0), IdentityMatrix[float64]())
	AssertMatrix(t, TranslationMatrix(5.0, 3.0).Untranslate(2.0, 1.0), TranslationMatrix(3.0, 2.0))

	AssertMatrix(t, TranslationMatrix(5, 3).Untranslate(5, 3), IdentityMatrix[int]())
	AssertMatrix(t, TranslationMatrix(5, 3).Untranslate(2, 1), TranslationMatrix(3, 2))
}

func TestMatrix_PreTranslate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().PreTranslate(5.0, 3.0), TranslationMatrix(5.0, 3.0))

	// S * PreTranslate(tx,ty) = T(tx,ty) * S  (translation applied before scale)
	got := ScaleMatrix(2.0, 2.0).PreTranslate(5.0, 3.0)
	AssertMatrix(t, got, Mat(2.0, 0.0, 5.0, 0.0, 2.0, 3.0))

	AssertMatrix(t, IdentityMatrix[int]().PreTranslate(5, 3), TranslationMatrix(5, 3))
	AssertMatrix(t, ScaleMatrix(2, 2).PreTranslate(5, 3), Mat(2, 0, 5, 0, 2, 3))
}

func TestMatrix_Rotate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().Rotate(0.0), IdentityMatrix[float64]())

	// point (1.0,0) rotated 90° → (0.0,1)
	p := Pt(1.0, 0.0).Transform(IdentityMatrix[float64]().Rotate(Pi / 2.0))
	assert.EqualDelta(t, p.X, 0.0, Delta)
	assert.EqualDelta(t, p.Y, 1.0, Delta)

	// integer: only multiples of 90° are lossless since cos/sin are then exact 0 or ±1
	AssertMatrix(t, IdentityMatrix[int]().Rotate(0), IdentityMatrix[int]())
	AssertMatrix(t, IdentityMatrix[int]().Rotate(Pi/2), Mat(0, -1, 0, 1, 0, 0))
	AssertMatrix(t, IdentityMatrix[int]().Rotate(Pi), Mat(-1, 0, 0, 0, -1, 0))
}

func TestMatrix_PreRotate(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().PreRotate(0.0), IdentityMatrix[float64]())

	// PreRotate(θ) * T differs from T * Rotate(θ)
	m1 := TranslationMatrix(5.0, 0.0).PreRotate(Pi / 2.0)
	m2 := TranslationMatrix(5.0, 0.0).Rotate(Pi / 2.0)
	assert.False(t, m1.Equal(m2))
}

func TestMatrix_Scale(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().Scale(2.0, 3.0), ScaleMatrix(2.0, 3.0))
	AssertMatrix(t, ScaleMatrix(2.0, 2.0).Scale(3.0, 3.0), ScaleMatrix(6.0, 6.0))

	AssertMatrix(t, IdentityMatrix[int]().Scale(2, 3), ScaleMatrix(2, 3))
	AssertMatrix(t, ScaleMatrix(2, 2).Scale(3, 3), ScaleMatrix(6, 6))
}

func TestMatrix_Unscale(t *testing.T) {
	AssertMatrix(t, ScaleMatrix(2.0, 3.0).Unscale(2.0, 3.0), IdentityMatrix[float64]())
	AssertMatrix(t, ScaleMatrix(4.0, 6.0).Unscale(2.0, 3.0), ScaleMatrix(2.0, 2.0))

	// zero guard: no change
	m := ScaleMatrix(2.0, 3.0)
	AssertMatrix(t, m.Unscale(0.0, 0.0), m)

	// integer: 1/1=1 is exact; useful for unscaling scale-1 matrices
	AssertMatrix(t, ScaleMatrix(1, 1).Unscale(1, 1), IdentityMatrix[int]())
}

func TestMatrix_PreScale(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().PreScale(2.0, 3.0), ScaleMatrix(2.0, 3.0))

	// S(2.0,2) * T(5.0,5) → point(1.0,2) → (12,14)
	m := TranslationMatrix(5.0, 5.0).PreScale(2.0, 2.0)
	AssertMatrix(t, m, Mat(2.0, 0.0, 10, 0.0, 2.0, 10))

	p := Pt(1.0, 2.0).Transform(m)
	assert.EqualDelta(t, p.X, 12.0, Delta)
	assert.EqualDelta(t, p.Y, 14.0, Delta)

	AssertMatrix(t, IdentityMatrix[int]().PreScale(2, 3), ScaleMatrix(2, 3))
	AssertMatrix(t, TranslationMatrix(5, 5).PreScale(2, 2), Mat(2, 0, 10, 0, 2, 10))
}

func TestMatrix_Equal(t *testing.T) {
	assert.True(t, IdentityMatrix[float64]().Equal(IdentityMatrix[float64]()))
	assert.False(t, IdentityMatrix[float64]().Equal(ScaleMatrix(2.0, 2.0)))
	assert.True(t, Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0).Equal(Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)))

	assert.True(t, IdentityMatrix[int]().Equal(IdentityMatrix[int]()))
	assert.False(t, IdentityMatrix[int]().Equal(ScaleMatrix(2, 2)))
	assert.True(t, Mat(1, 2, 3, 4, 5, 6).Equal(Mat(1, 2, 3, 4, 5, 6)))
}

func TestMatrix_IsZero(t *testing.T) {
	assert.True(t, Matrix[float64]{}.IsZero())
	assert.False(t, IdentityMatrix[float64]().IsZero())
	assert.False(t, IdentityMatrix[float32]().IsZero())
	assert.True(t, Matrix[int]{}.IsZero())
	assert.False(t, IdentityMatrix[int]().IsZero())
}

func TestMatrix_Int(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[float64]().Int(), IdentityMatrix[int]())
	AssertMatrix(t, TranslationMatrix(5, 3).Int(), TranslationMatrix(5, 3))
	// truncates toward zero
	AssertMatrix(t, TranslationMatrix(5.9, 3.1).Int(), TranslationMatrix(5, 3))
}

func TestMatrix_Float(t *testing.T) {
	AssertMatrix(t, IdentityMatrix[int]().Float(), IdentityMatrix[float64]())
	AssertMatrix(t, TranslationMatrix(5, 3).Float(), TranslationMatrix(5.0, 3.0))
}

func TestMatrix_String(t *testing.T) {
	assert.Equal(t, IdentityMatrix[float64]().String(), "[[1, 0, 0], [0, 1, 0]]")
	assert.Equal(t, TranslationMatrix(5.1, 3).String(), "[[1, 0, 5.10], [0, 1, 3]]")

	assert.Equal(t, IdentityMatrix[int]().String(), "[[1, 0, 0], [0, 1, 0]]")
	assert.Equal(t, TranslationMatrix(5, 3).String(), "[[1, 0, 5], [0, 1, 3]]")
}

func TestMatrix_Marshall(t *testing.T) {
	assert.JSON(t, Mat[float64](1.0, 2.1, 3.2, 4.0, 5.3, 6.4), `{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`)
	assert.JSON(t, Mat[int](1, 2, 3, 4, 5, 6), `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`)
}

func TestMatrix_Unmarshall(t *testing.T) {
	var m1 Matrix[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`), &m1))
	AssertMatrix(t, m1, Matrix[float64]{1.0, 2.1, 3.2, 4.0, 5.3, 6.4})

	var m2 Matrix[float32]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`), &m2))
	AssertMatrix(t, m2, Matrix[float32]{1.0, 2.1, 3.2, 4.0, 5.3, 6.4})

	var m3 Matrix[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`), &m3))
	AssertMatrix(t, m3, Matrix[int]{1, 2, 3, 4, 5, 6})
}

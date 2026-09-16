package geom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestMatrix_Constructor(t *testing.T) {
	t.Run("from components", func(t *testing.T) {
		AssertMatrix(t, Mat(1, 2, 3, 4, 5, 6), Matrix[int]{A: 1, B: 2, C: 3, D: 4, E: 5, F: 6})
		AssertMatrix(t, Mat(1.0, 2.1, 3.2, 4.0, 5.3, 6.4), Matrix[float64]{A: 1.0, B: 2.1, C: 3.2, D: 4.0, E: 5.3, F: 6.4})
	})
	t.Run("identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64](), Mat(1.0, 0.0, 0.0, 0.0, 1.0, 0.0))
		AssertMatrix(t, IdentityMatrix[float32](), Mat[float32](1, 0, 0, 0, 1, 0))
		AssertMatrix(t, IdentityMatrix[int](), Mat(1, 0, 0, 0, 1, 0))
	})
	t.Run("translation", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5.0, 3.0), Mat(1.0, 0.0, 5.0, 0.0, 1.0, 3.0))
		AssertMatrix(t, TranslationMatrix[float32](5, 3), Mat[float32](1, 0, 5, 0, 1, 3))
		AssertMatrix(t, TranslationMatrix(5, 3), Mat(1, 0, 5, 0, 1, 3))
	})
	t.Run("scale", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(2.0, 3.0), Mat(2.0, 0.0, 0.0, 0.0, 3.0, 0.0))
		AssertMatrix(t, ScaleMatrix[float32](2, 3), Mat[float32](2, 0, 0, 0, 3, 0))
		AssertMatrix(t, ScaleMatrix(2, 3), Mat(2, 0, 0, 0, 3, 0))
	})
	t.Run("rotation", func(t *testing.T) {
		AssertMatrix(t, RotationMatrix[float64](0.0), IdentityMatrix[float64]())
		AssertMatrix(t, RotationMatrix[float64](Pi/2), Mat(0.0, -1.0, 0.0, 1.0, 0.0, 0.0))
		AssertMatrix(t, RotationMatrix[float64](Pi), Mat(-1.0, 0.0, 0.0, 0.0, -1.0, 0.0))

		// the angle is always computed in float64 and narrowed afterwards
		AssertMatrix(t, RotationMatrix[float32](Pi/2), Mat[float32](0, -1, 0, 1, 0, 0))
	})
	t.Run("integer rotation is exact only for quarter turns", func(t *testing.T) {
		// cos(π/2) ≈ 6e-17 rounds to 0; sin(π/2) = 1 exactly
		AssertMatrix(t, RotationMatrix[int](Pi/2), Mat(0, -1, 0, 1, 0, 0))
		AssertMatrix(t, RotationMatrix[int](Pi), Mat(-1, 0, 0, 0, -1, 0))

		// anything else is not a rotation: π/6 rounds to a √2 scale plus shear
		AssertMatrix(t, RotationMatrix[int](Pi/6), Mat(1, -1, 0, 1, 1, 0))
	})
}

func TestMatrix_Multiply(t *testing.T) {
	t.Run("identity is neutral", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Multiply(TranslationMatrix(5.0, 3.0)), TranslationMatrix(5.0, 3.0))
		AssertMatrix(t, TranslationMatrix(5.0, 3.0).Multiply(IdentityMatrix[float64]()), TranslationMatrix(5.0, 3.0))
		AssertMatrix(t, IdentityMatrix[int]().Multiply(TranslationMatrix(5, 3)), TranslationMatrix(5, 3))
	})
	t.Run("translations compose additively", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5.0, 3.0).Multiply(TranslationMatrix(2.0, 1.0)), TranslationMatrix(7.0, 4.0))
		AssertMatrix(t, TranslationMatrix[float32](5, 3).Multiply(TranslationMatrix[float32](2, 1)), TranslationMatrix[float32](7, 4))
		AssertMatrix(t, TranslationMatrix(5, 3).Multiply(TranslationMatrix(2, 1)), TranslationMatrix(7, 4))
	})
	t.Run("scales compose multiplicatively", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(2.0, 3.0).Multiply(ScaleMatrix(4.0, 2.0)), ScaleMatrix(8.0, 6.0))
		AssertMatrix(t, ScaleMatrix[float32](2, 3).Multiply(ScaleMatrix[float32](4, 2)), ScaleMatrix[float32](8, 6))
		AssertMatrix(t, ScaleMatrix(2, 3).Multiply(ScaleMatrix(4, 2)), ScaleMatrix(8, 6))
	})
}

func TestMatrix_Inverse(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Inverse(), IdentityMatrix[float64]())
		AssertMatrix(t, TranslationMatrix(5.0, 3.0).Inverse(), TranslationMatrix(-5.0, -3.0))
		AssertMatrix(t, ScaleMatrix(2.0, 4.0).Inverse(), ScaleMatrix(0.5, 0.25))
		AssertMatrix(t, ScaleMatrix(0.0005, 0.0005).Inverse(), ScaleMatrix(2000.0, 2000.0))

		// float32 inverts exactly too — the fractions are not rounded away
		AssertMatrix(t, ScaleMatrix[float32](2, 4).Inverse(), ScaleMatrix[float32](0.5, 0.25))
		AssertMatrix(t, ScaleMatrix[float32](0.005, 0.005).Inverse(), ScaleMatrix[float32](200, 200))
		AssertMatrix(t, TranslationMatrix[float32](5, 3).Inverse(), TranslationMatrix[float32](-5, -3))
	})
	t.Run("singular matrix returns itself", func(t *testing.T) {
		zero := Mat(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
		AssertMatrix(t, zero.Inverse(), zero)
	})
	t.Run("large translations do not overflow", func(t *testing.T) {
		m := Mat[int64](1, 0, 1<<40, 0, 1, 1<<40)
		AssertMatrix(t, m.Inverse(), Mat[int64](1, 0, -(1<<40), 0, 1, -(1<<40)))
	})
	t.Run("integer is exact only for unit determinant", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[int]().Inverse(), IdentityMatrix[int]())
		AssertMatrix(t, TranslationMatrix(5, 3).Inverse(), TranslationMatrix(-5, -3))
		AssertMatrix(t, RotationMatrix[int](Pi/2).Inverse(), RotationMatrix[int](-Pi/2))

		// det=4, invDet=0.25 → Cast[int](2*0.25) = Cast[int](0.5) = 1, so the inverse does not undo it
		AssertMatrix(t, ScaleMatrix(2, 2).Inverse(), IdentityMatrix[int]())
	})
}

func TestMatrix_IsInvertible(t *testing.T) {
	t.Run("non-zero determinant", func(t *testing.T) {
		assert.True(t, IdentityMatrix[int]().IsInvertible())
		assert.True(t, ScaleMatrix(2.0, 0.5).IsInvertible())
		assert.True(t, ScaleMatrix(0.0005, 0.0005).IsInvertible())
		assert.True(t, ScaleMatrix[float32](0.005, 0.005).IsInvertible())
	})
	t.Run("zero determinant", func(t *testing.T) {
		assert.False(t, Matrix[int]{}.IsInvertible())
		assert.False(t, ScaleMatrix(1.0, 0.0).IsInvertible())
		assert.False(t, Mat(1.0, 2.0, 0.0, 2.0, 4.0, 0.0).IsInvertible())
	})
}

func TestMatrix_Determinant(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, IdentityMatrix[float64]().Determinant(), 1.0)
		AssertNumber(t, ScaleMatrix(2.0, 3.0).Determinant(), 6.0)
		AssertNumber(t, Mat(1.0, 2.0, 0.0, 3.0, 4.0, 0.0).Determinant(), -2.0) // 1*4 - 2*3
		AssertNumber(t, ScaleMatrix[float32](2, 3).Determinant(), float32(6))
	})
	t.Run("rotation preserves area", func(t *testing.T) {
		AssertNumber(t, RotationMatrix[float64](Pi/4).Determinant(), 1.0)
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, IdentityMatrix[int]().Determinant(), 1)
		AssertNumber(t, ScaleMatrix(2, 3).Determinant(), 6)
		AssertNumber(t, Mat(1, 2, 0, 3, 4, 0).Determinant(), -2) // 1*4 - 2*3
	})
}

func TestMatrix_Translate(t *testing.T) {
	t.Run("from identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Translate(5.0, 3.0), TranslationMatrix(5.0, 3.0))
	})
	t.Run("accumulates", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Translate(2.0, 1.0).Translate(3.0, 2.0), TranslationMatrix(5.0, 3.0))
		AssertMatrix(t, IdentityMatrix[float32]().Translate(2, 1).Translate(3, 2), TranslationMatrix[float32](5, 3))
		AssertMatrix(t, IdentityMatrix[int]().Translate(2, 1).Translate(3, 2), TranslationMatrix(5, 3))
	})
}

func TestMatrix_Untranslate(t *testing.T) {
	t.Run("undoes a translation", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5.0, 3.0).Untranslate(5.0, 3.0), IdentityMatrix[float64]())
	})
	t.Run("partial", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5.0, 3.0).Untranslate(2.0, 1.0), TranslationMatrix(3.0, 2.0))
		AssertMatrix(t, TranslationMatrix[float32](5, 3).Untranslate(2, 1), TranslationMatrix[float32](3, 2))
		AssertMatrix(t, TranslationMatrix(5, 3).Untranslate(2, 1), TranslationMatrix(3, 2))
	})
}

func TestMatrix_PreTranslate(t *testing.T) {
	t.Run("from identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().PreTranslate(5.0, 3.0), TranslationMatrix(5.0, 3.0))
	})
	t.Run("applies before the current transform", func(t *testing.T) {
		// S * PreTranslate(tx,ty) = T(tx,ty) * S
		AssertMatrix(t, ScaleMatrix(2.0, 2.0).PreTranslate(5.0, 3.0), Mat(2.0, 0.0, 5.0, 0.0, 2.0, 3.0))
		AssertMatrix(t, ScaleMatrix[float32](2, 2).PreTranslate(5, 3), Mat[float32](2, 0, 5, 0, 2, 3))
		AssertMatrix(t, ScaleMatrix(2, 2).PreTranslate(5, 3), Mat(2, 0, 5, 0, 2, 3))
	})
}

func TestMatrix_Rotate(t *testing.T) {
	t.Run("zero angle is identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Rotate(0.0), IdentityMatrix[float64]())
		AssertMatrix(t, IdentityMatrix[int]().Rotate(0), IdentityMatrix[int]())
	})
	t.Run("quarter turn", func(t *testing.T) {
		AssertPoint(t, Pt(1.0, 0.0).Transform(IdentityMatrix[float64]().Rotate(Pi/2)), Pt(0.0, 1.0))
		AssertMatrix(t, IdentityMatrix[float32]().Rotate(Pi), Mat[float32](-1, 0, 0, 0, -1, 0))
	})
	t.Run("float matrix rotates an int point", func(t *testing.T) {
		// the point keeps its integer coordinates; the matrix carries the fractions
		AssertPoint(t, Pt(3, 0).Transform(IdentityMatrix[float64]().Rotate(Pi/6)), Pt(3, 2))
	})
	t.Run("integer is exact only for quarter turns", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[int]().Rotate(Pi/2), Mat(0, -1, 0, 1, 0, 0))
		AssertMatrix(t, IdentityMatrix[int]().Rotate(Pi), Mat(-1, 0, 0, 0, -1, 0))
	})
}

func TestMatrix_PreRotate(t *testing.T) {
	t.Run("zero angle is identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().PreRotate(0.0), IdentityMatrix[float64]())
	})
	t.Run("differs from rotating after", func(t *testing.T) {
		before := TranslationMatrix(5.0, 0.0).PreRotate(Pi / 2)
		after := TranslationMatrix(5.0, 0.0).Rotate(Pi / 2)

		assert.False(t, before.Equal(after))
	})
}

func TestMatrix_Scale(t *testing.T) {
	t.Run("from identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Scale(2.0, 3.0), ScaleMatrix(2.0, 3.0))
	})
	t.Run("accumulates", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(2.0, 2.0).Scale(3.0, 3.0), ScaleMatrix(6.0, 6.0))
		AssertMatrix(t, ScaleMatrix[float32](2, 2).Scale(3, 3), ScaleMatrix[float32](6, 6))
		AssertMatrix(t, ScaleMatrix(2, 2).Scale(3, 3), ScaleMatrix(6, 6))
	})
}

func TestMatrix_Unscale(t *testing.T) {
	t.Run("undoes a scale", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(2.0, 3.0).Unscale(2.0, 3.0), IdentityMatrix[float64]())
	})
	t.Run("partial", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(4.0, 6.0).Unscale(2.0, 3.0), ScaleMatrix(2.0, 2.0))
		AssertMatrix(t, ScaleMatrix[float32](4, 6).Unscale(2, 3), ScaleMatrix[float32](2, 2))
	})
	t.Run("zero factor leaves the axis unchanged", func(t *testing.T) {
		m := ScaleMatrix(2.0, 3.0)

		AssertMatrix(t, m.Unscale(0.0, 0.0), m)
		AssertMatrix(t, m.Unscale(0.0, 3.0), ScaleMatrix(2.0, 1.0))
		AssertMatrix(t, m.Unscale(2.0, 0.0), ScaleMatrix(1.0, 3.0))
	})
	t.Run("integer divides each column", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(4, 6).Unscale(1, 1), ScaleMatrix(4, 6))
		AssertMatrix(t, ScaleMatrix(4, 6).Unscale(2, 3), ScaleMatrix(2, 2))
		AssertMatrix(t, ScaleMatrix(2, 3).Unscale(2, 3), IdentityMatrix[int]())
		AssertMatrix(t, Mat(4, 6, 8, 2, 9, 5).Unscale(2, 3), Mat(2, 2, 8, 1, 3, 5))
	})
	t.Run("integer rounds a component that does not divide", func(t *testing.T) {
		AssertMatrix(t, ScaleMatrix(3, 5).Unscale(2, 2), ScaleMatrix(2, 3))
	})
}

func TestMatrix_PreScale(t *testing.T) {
	t.Run("from identity", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().PreScale(2.0, 3.0), ScaleMatrix(2.0, 3.0))
	})
	t.Run("applies before the current transform", func(t *testing.T) {
		// S(2,2) * T(5,5) maps (1,2) to (12,14)
		m := TranslationMatrix(5.0, 5.0).PreScale(2.0, 2.0)

		AssertMatrix(t, m, Mat(2.0, 0.0, 10.0, 0.0, 2.0, 10.0))
		AssertPoint(t, Pt(1.0, 2.0).Transform(m), Pt(12.0, 14.0))

		AssertMatrix(t, TranslationMatrix[float32](5, 5).PreScale(2, 2), Mat[float32](2, 0, 10, 0, 2, 10))
		AssertMatrix(t, TranslationMatrix(5, 5).PreScale(2, 2), Mat(2, 0, 10, 0, 2, 10))
	})
}

func TestMatrix_Equal(t *testing.T) {
	t.Run("same matrix", func(t *testing.T) {
		assert.True(t, IdentityMatrix[float64]().Equal(IdentityMatrix[float64]()))
		assert.True(t, IdentityMatrix[float32]().Equal(IdentityMatrix[float32]()))
		assert.True(t, Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0).Equal(Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)))
		assert.True(t, Mat(1, 2, 3, 4, 5, 6).Equal(Mat(1, 2, 3, 4, 5, 6)))
	})
	t.Run("different matrix", func(t *testing.T) {
		assert.False(t, IdentityMatrix[float64]().Equal(ScaleMatrix(2.0, 2.0)))
		assert.False(t, IdentityMatrix[float32]().Equal(ScaleMatrix[float32](2, 2)))
		assert.False(t, IdentityMatrix[int]().Equal(ScaleMatrix(2, 2)))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, IdentityMatrix[float64]().Equal(Mat(1.000001, 0.0, 0.0, 0.0, 1.0, 0.0)))
	})
}

func TestMatrix_IsZero(t *testing.T) {
	t.Run("zero matrix", func(t *testing.T) {
		assert.True(t, Matrix[float64]{}.IsZero())
		assert.True(t, Matrix[float32]{}.IsZero())
		assert.True(t, Matrix[int]{}.IsZero())
	})
	t.Run("non-zero matrix", func(t *testing.T) {
		assert.False(t, IdentityMatrix[float64]().IsZero())
		assert.False(t, IdentityMatrix[float32]().IsZero())
		assert.False(t, IdentityMatrix[int]().IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Mat(0.0, 0.0, 0.000001, 0.0, 0.0, 0.0).IsZero())
	})
}

func TestMatrix_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5, 3).Int(), TranslationMatrix(5, 3))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float64]().Int(), IdentityMatrix[int]())
		AssertMatrix(t, TranslationMatrix(5.9, 3.1).Int(), TranslationMatrix(6, 3))
		AssertMatrix(t, TranslationMatrix(-5.9, -3.1).Int(), TranslationMatrix(-6, -3))
	})
}

func TestMatrix_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[int]().Float(), IdentityMatrix[float64]())
	})
	t.Run("float32 widens", func(t *testing.T) {
		AssertMatrix(t, IdentityMatrix[float32]().Float(), IdentityMatrix[float64]())
		AssertMatrix(t, TranslationMatrix[float32](5, 3).Float(), TranslationMatrix(5.0, 3.0))
	})
	t.Run("float64 is a no-op", func(t *testing.T) {
		AssertMatrix(t, TranslationMatrix(5.1, 3.0).Float(), TranslationMatrix(5.1, 3.0))
	})
}

func TestMatrix_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, IdentityMatrix[int]().String(), "[[1, 0, 0], [0, 1, 0]]")
		assert.Equal(t, TranslationMatrix(5, 3).String(), "[[1, 0, 5], [0, 1, 3]]")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, IdentityMatrix[float64]().String(), "[[1.00, 0.00, 0.00], [0.00, 1.00, 0.00]]")
		assert.Equal(t, TranslationMatrix(5.1, 3.0).String(), "[[1.00, 0.00, 5.10], [0.00, 1.00, 3.00]]")

		assert.Equal(t, IdentityMatrix[float32]().String(), "[[1.00, 0.00, 0.00], [0.00, 1.00, 0.00]]")
		assert.Equal(t, TranslationMatrix[float32](5, 3).String(), "[[1.00, 0.00, 5.00], [0.00, 1.00, 3.00]]")
	})
}

func TestMatrix_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Mat(1, 2, 3, 4, 5, 6), `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`)

		var m Matrix[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`), &m))
		AssertMatrix(t, m, Mat(1, 2, 3, 4, 5, 6))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Mat(1.0, 2.1, 3.2, 4.0, 5.3, 6.4), `{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`)
		assert.JSON(t, Mat[float32](1, 2, 3, 4, 5, 6), `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`)

		var m64 Matrix[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`), &m64))
		AssertMatrix(t, m64, Mat(1.0, 2.1, 3.2, 4.0, 5.3, 6.4))

		var m32 Matrix[float32]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":1,"b":2.1,"c":3.2,"d":4,"e":5.3,"f":6.4}`), &m32))
		AssertMatrix(t, m32, Mat[float32](1, 2.1, 3.2, 4, 5.3, 6.4))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, matrix := range matrixFixtures {
			data, err := json.Marshal(matrix)
			assert.NoError(t, err)

			var decoded Matrix[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, matrix)
		}
	})
}

func TestMatrix_Properties(t *testing.T) {
	t.Run("identity is neutral", func(t *testing.T) {
		identity := IdentityMatrix[float64]()

		for _, matrix := range matrixFixtures {
			assert.True(t, matrix.Multiply(identity).Equal(matrix), fmt.Sprintf("%s: ", matrix))
			assert.True(t, identity.Multiply(matrix).Equal(matrix), fmt.Sprintf("%s: ", matrix))
		}
	})
	t.Run("multiply is associative", func(t *testing.T) {
		for _, a := range matrixFixtures {
			for _, b := range matrixFixtures {
				for _, c := range matrixFixtures {
					left, right := a.Multiply(b).Multiply(c), a.Multiply(b.Multiply(c))
					assert.True(t, left.Equal(right), fmt.Sprintf("%s → %s → %s: ", a, b, c))
				}
			}
		}
	})
	t.Run("inverse undoes the transform", func(t *testing.T) {
		for _, matrix := range matrixFixtures {
			if !matrix.IsInvertible() {
				continue
			}

			assert.True(t, matrix.Multiply(matrix.Inverse()).Equal(IdentityMatrix[float64]()), fmt.Sprintf("%s: ", matrix))
		}
	})
	t.Run("determinant is multiplicative", func(t *testing.T) {
		for _, a := range matrixFixtures {
			for _, b := range matrixFixtures {
				product := a.Multiply(b).Determinant()
				AssertNumber(t, product, a.Determinant()*b.Determinant(), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("translate and untranslate are inverse", func(t *testing.T) {
		for _, matrix := range matrixFixtures {
			assert.True(t, matrix.Translate(5, 3).Untranslate(5, 3).Equal(matrix), fmt.Sprintf("%s: ", matrix))
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, matrix := range matrixFixtures {
			assert.True(t, matrix.Scale(2, 4).Unscale(2, 4).Equal(matrix), fmt.Sprintf("%s: ", matrix))
		}
	})
	t.Run("multiplication composes point transforms", func(t *testing.T) {
		for _, a := range matrixFixtures {
			for _, b := range matrixFixtures {
				for _, point := range pointFixtures {
					composed := point.Transform(a.Multiply(b))
					stepwise := point.Transform(b).Transform(a)

					assert.True(t, composed.Equal(stepwise), fmt.Sprintf("%s → %s on %s: ", a, b, point))
				}
			}
		}
	})
	t.Run("rotation preserves area", func(t *testing.T) {
		for _, angle := range []float64{0, Pi / 6, Pi / 2, 2, -1.5} {
			AssertNumber(t, RotationMatrix[float64](angle).Determinant(), 1.0, fmt.Sprintf("%v rad: ", angle))
		}
	})
}

func TestMatrix_Immutable(t *testing.T) {
	m := Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)

	m.Multiply(IdentityMatrix[float64]())
	m.Inverse()
	m.Translate(5, 3)
	m.Untranslate(5, 3)
	m.PreTranslate(5, 3)
	m.Rotate(Pi / 3)
	m.PreRotate(Pi / 3)
	m.Scale(2, 4)
	m.Unscale(2, 4)
	m.PreScale(2, 4)

	AssertMatrix(t, m, Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0))
}

// matrixFixtures span identity, translation, scale, rotation, and a general affine transform.
var matrixFixtures = []Matrix[float64]{
	IdentityMatrix[float64](),
	TranslationMatrix(5.0, 3.0),
	TranslationMatrix(-2.5, 0.75),
	ScaleMatrix(2.0, 3.0),
	ScaleMatrix(-1.0, 0.5),
	RotationMatrix[float64](Pi / 6),
	RotationMatrix[float64](-Pi / 2),
	Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6),
}

func ExampleMat() {
	fmt.Println(Mat(1, 0, 5, 0, 1, 3))
	// Output: [[1, 0, 5], [0, 1, 3]]
}

func ExampleIdentityMatrix() {
	fmt.Println(IdentityMatrix[float64]())
	// Output: [[1.00, 0.00, 0.00], [0.00, 1.00, 0.00]]
}

package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestVector_Constructor(t *testing.T) {
	t.Run("from components", func(t *testing.T) {
		AssertVector(t, Vec(10, 16), Vector[int]{X: 10, Y: 16})
		AssertVector(t, Vec[float64](0.16, 204), Vector[float64]{X: 0.16, Y: 204})
	})
	t.Run("zero", func(t *testing.T) {
		AssertVector(t, ZeroVector[int](), Vector[int]{})
		AssertVector(t, ZeroVector[float64](), Vector[float64]{})
	})
	t.Run("one", func(t *testing.T) {
		AssertVector(t, OneVector[int](), Vec(1, 1))
		AssertVector(t, OneVector[float64](), Vec(1.0, 1.0))
	})
}

func TestVectorFromAngle(t *testing.T) {
	t.Run("cardinal angles", func(t *testing.T) {
		AssertVector(t, VectorFromAngle(0, 5.0), Vec(5.0, 0.0))
		AssertVector(t, VectorFromAngle(Pi/2, 1.0), Vec(0.0, 1.0))
		AssertVector(t, VectorFromAngle(Pi, 1.0), Vec(-1.0, 0.0))
		AssertVector(t, VectorFromAngle(-Pi/2, 1.0), Vec(0.0, -1.0))
	})
	t.Run("diagonal angle", func(t *testing.T) {
		AssertVector(t, VectorFromAngle(Pi/4, 1.0), Vec(OneOverSqrt2, OneOverSqrt2))
	})
	t.Run("integer rounds the components", func(t *testing.T) {
		// cos(π/2) ≈ 6e-17 rounds to 0, sin(π/2) = 1
		AssertVector(t, VectorFromAngle(Pi/2, 4), Vec(0, 4))
	})
}

func TestVectorFromAngleSize(t *testing.T) {
	t.Run("traces the ellipse", func(t *testing.T) {
		AssertVector(t, VectorFromAngleSize(0, Sz(3.0, 5.0)), Vec(3.0, 0.0))
		AssertVector(t, VectorFromAngleSize(Pi/2, Sz(3.0, 5.0)), Vec(0.0, 5.0))
		AssertVector(t, VectorFromAngleSize(Pi, Sz(3.0, 5.0)), Vec(-3.0, 0.0))
	})
	t.Run("square size matches the angle constructor", func(t *testing.T) {
		for _, angle := range []float64{0, Pi / 6, Pi / 4, 2, -1.5} {
			ellipse, circle := VectorFromAngleSize(angle, SzU(7.0)), VectorFromAngle(angle, 7.0)
			assert.True(t, ellipse.Equal(circle), fmt.Sprintf("%v rad: ", angle))
		}
	})
	t.Run("integer rounds after scaling", func(t *testing.T) {
		// (3*cos(π/6), 3*sin(π/6)) = (2.598, 1.5) → (3, 2).
		// Rounding the unit vector first would have given (1, 1) scaled to (3, 3).
		AssertVector(t, VectorFromAngleSize(Pi/6, Sz(3, 3)), Vec(3, 2))
	})
}

func TestVector_Transform(t *testing.T) {
	t.Run("float64 matrix", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), Vec(48, 132))
		AssertVector(t, Vec(0.6, -0.25).Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), Vec(0.085, 1.265))
	})
	t.Run("float32 matrix", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Transform(Mat[float32](1, 2, 3, 4, 5, 6)), Vec(42, 120))
		AssertVector(t, Vec(0.6, -0.25).Transform(Mat[float32](1, 2, 3, 4, 5, 6)), Vec(0.1, 1.15))
	})
	t.Run("integer matrix converted to float", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Transform(Mat(1, 2, 3, 4, 5, 6).Float()), Vec(42, 120))
		AssertVector(t, Vec(0.6, -0.25).Transform(Mat(1, 2, 3, 4, 5, 6).Float()), Vec(0.1, 1.15))
	})
	t.Run("translation is ignored", func(t *testing.T) {
		// a vector is a displacement, so only the linear part applies
		AssertVector(t, Vec(1.0, 2.0).Transform(TranslationMatrix(5.0, 7.0)), Vec(1.0, 2.0))
	})
}

func TestVector_Add(t *testing.T) {
	t.Run("vector", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Add(Vec(3, -2)), Vec(13, 14))
		AssertVector(t, Vec(0.6, -0.25).Add(Vec(100.1, -0.1)), Vec(100.7, -0.35))
	})
	t.Run("components", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).AddXY(3, -2), Vec(13, 14))
		AssertVector(t, Vec(0.6, -0.25).AddXY(100.1, -0.1), Vec(100.7, -0.35))
	})
}

func TestVector_Subtract(t *testing.T) {
	t.Run("vector", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Subtract(Vec(3, -3)), Vec(7, 19))
		AssertVector(t, Vec(0.6, -0.25).Subtract(Vec(100.1, -0.1)), Vec(-99.5, -0.15))
	})
	t.Run("components", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).SubtractXY(3, -3), Vec(7, 19))
		AssertVector(t, Vec(0.6, -0.25).SubtractXY(100.1, -0.1), Vec(-99.5, -0.15))
	})
}

func TestVector_Multiply(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Multiply(3), Vec(30, 48))
		AssertVector(t, Vec(0.6, -0.25).Multiply(-1.5), Vec(-0.9, 0.375))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).MultiplyXY(3, 2), Vec(30, 32))
		AssertVector(t, Vec(0.6, -0.25).MultiplyXY(-1.5, 2), Vec(-0.9, -0.5))
	})
}

func TestVector_Divide(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Divide(2), Vec(5, 8))
		AssertVector(t, Vec(0.6, -0.25).Divide(-2), Vec(-0.3, 0.125))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).DivideXY(3, 2), Vec(3, 8)) // int: 3.33 rounds to 3
		AssertVector(t, Vec(0.6, -0.25).DivideXY(-4, 0.5), Vec(-0.15, -0.5))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Vec(10, 16).Divide(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Vec(10, 16).DivideXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestVector_Negate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Negate(), Vec(-10, -16))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Negate(), Vec(-0.6, 0.25))
	})
}

func TestVector_Rotate(t *testing.T) {
	t.Run("quarter turn", func(t *testing.T) {
		AssertVector(t, Vec(1, 0).Rotate(ToRadians(90)), Vec(0, 1))
		AssertVector(t, Vec(1.0, 0.0).Rotate(ToRadians(-90)), Vec(0.0, -1.0))
	})
	t.Run("half turn negates", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Rotate(Pi), Vec(-0.6, 0.25))
	})
	t.Run("full turn is identity", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Rotate(2*Pi), Vec(0.6, -0.25))
	})
}

func TestVector_Resize(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Resize(5), Vec(3, 4))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Resize(5), Vec(4.615384, -1.923076))
	})
	t.Run("the zero vector resizes along +X", func(t *testing.T) {
		AssertVector(t, ZeroVector[float64]().Resize(5), Vec(5.0, 0.0))
		AssertVector(t, ZeroVector[int]().Resize(2.4), Vec(2, 0))
	})
	t.Run("a vector shorter than Epsilon keeps its direction", func(t *testing.T) {
		AssertVector(t, Vec(1e-7, 1e-7).Resize(5), Vec(5*OneOverSqrt2, 5*OneOverSqrt2))
		AssertVector(t, Vec[float32](0, -1e-5).Resize(2), Vec[float32](0, -2))
	})
	t.Run("a subnormal vector does not overflow", func(t *testing.T) {
		AssertVector(t, Vec(5e-324, 0).Resize(3), Vec(3.0, 0.0))
		AssertVector(t, Vec(0, -5e-324).Normalize(), Vec(0.0, -1.0))
	})
}

func TestVector_Normalize(t *testing.T) {
	t.Run("float keeps the direction", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Normalize(), Vec(0.923076, -0.384615))
		AssertNumber(t, Vec(0.6, -0.25).Normalize().Length(), 1.0)
		AssertVector(t, Vec(1e-7, 1e-7).Normalize(), Vec(OneOverSqrt2, OneOverSqrt2))
	})
	t.Run("int snaps to the longer axis", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Normalize(), Vec(0, 1))
		AssertVector(t, Vec(-10, 16).Normalize(), Vec(0, 1))
		AssertVector(t, Vec(10, -16).Normalize(), Vec(0, -1))
		AssertVector(t, Vec(-10, -16).Normalize(), Vec(0, -1))

		AssertVector(t, Vec(16, 10).Normalize(), Vec(1, 0))
		AssertVector(t, Vec(-16, 10).Normalize(), Vec(-1, 0))
		AssertVector(t, Vec(16, -10).Normalize(), Vec(1, 0))
		AssertVector(t, Vec(-16, -10).Normalize(), Vec(-1, 0))
	})
	t.Run("int is already normalized on an axis", func(t *testing.T) {
		AssertVector(t, Vec(3, 0).Normalize(), Vec(1, 0))
		AssertVector(t, Vec(0, -5).Normalize(), Vec(0, -1))
	})
	t.Run("int breaks a diagonal tie towards X", func(t *testing.T) {
		AssertVector(t, Vec(1, 1).Normalize(), Vec(1, 0))
		AssertVector(t, Vec(-4, 4).Normalize(), Vec(-1, 0))
	})
	t.Run("zero vector is right by convention", func(t *testing.T) {
		AssertVector(t, Vec(0, 0).Normalize(), Vec(1, 0))
		AssertVector(t, ZeroVector[float64]().Normalize(), Vec(1.0, 0.0))
	})
}

func TestVector_Abs(t *testing.T) {
	t.Run("negative components", func(t *testing.T) {
		AssertVector(t, Vec(-1, -3).Abs(), Vec(1, 3))
		AssertVector(t, Vec(0.6, -0.25).Abs(), Vec(0.6, 0.25))
	})
	t.Run("non-negative components unchanged", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Abs(), Vec(10, 16))
		AssertVector(t, Vec(0.0, 0.0).Abs(), Vec(0.0, 0.0))
	})
}

func TestVector_Round(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(1.4, -1.5).Round(), Vec(1.0, -2.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertVector(t, Vec(3, -2).Round(), Vec(3, -2))
	})
}

func TestVector_Floor(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(1.9, -1.1).Floor(), Vec(1.0, -2.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertVector(t, Vec(3, -2).Floor(), Vec(3, -2))
	})
}

func TestVector_Ceil(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(1.1, -1.9).Ceil(), Vec(2.0, -1.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertVector(t, Vec(3, -2).Ceil(), Vec(3, -2))
	})
}

func TestVector_Dot(t *testing.T) {
	t.Run("general vectors", func(t *testing.T) {
		AssertNumber(t, Vec(10, 16).Dot(Vec(3, -3)), -18)
		AssertNumber(t, Vec(0.6, -0.25).Dot(Vec(100.1, -0.1)), 60.085)
	})
	t.Run("perpendicular vectors are zero", func(t *testing.T) {
		AssertNumber(t, Vec(1, 0).Dot(Vec(0, 1)), 0)
	})
	t.Run("a float vector and its normal are exactly zero", func(t *testing.T) {
		for _, v := range []Vector[float64]{Vec(0.1, 0.3), Vec(1.1, 3.3), Vec(2.5, 1e5+0.1)} {
			assert.Equal(t, v.Dot(v.Normal()), 0.0, v.String())
		}
		for _, v := range []Vector[float32]{Vec[float32](0.1, 0.3), Vec[float32](1.1, 3.3)} {
			assert.Equal(t, v.Dot(v.Normal()), float32(0), v.String())
		}
	})
	t.Run("narrow integers do not overflow mid-computation", func(t *testing.T) {
		AssertNumber(t, Vec[int8](100, 50).Dot(Vec[int8](2, -2)), 100)
	})
}

func TestVector_Cross(t *testing.T) {
	t.Run("general vectors", func(t *testing.T) {
		AssertNumber(t, Vec(10, 16).Cross(Vec(3, -3)), -78)
		AssertNumber(t, Vec(0.6, -0.25).Cross(Vec(100.1, -0.1)), 24.965)
	})
	t.Run("parallel vectors are zero", func(t *testing.T) {
		AssertNumber(t, Vec(2, 4).Cross(Vec(1, 2)), 0)
	})
	t.Run("a float vector with itself is exactly zero", func(t *testing.T) {
		for _, v := range []Vector[float64]{Vec(0.1, 0.3), Vec(1.1, 3.3), Vec(2.5, 1e5+0.1)} {
			assert.Equal(t, v.Cross(v), 0.0, v.String())
		}
		for _, v := range []Vector[float32]{Vec[float32](0.1, 0.3), Vec[float32](1.1, 3.3)} {
			assert.Equal(t, v.Cross(v), float32(0), v.String())
		}
	})
}

func TestVector_Normal(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Normal(), Vec(-16, 10))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Normal(), Vec(0.25, 0.6))
	})
}

func TestVector_Length(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Vec(10, 16).Length(), math.Sqrt(356))
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Vec(0.6, -0.25).Length(), 0.65)
	})
}

func TestVector_LengthSquared(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Vec(10, 16).LengthSquared(), 356)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Vec(0.6, -0.25).LengthSquared(), 0.4225)
	})
}

func TestVector_Angle(t *testing.T) {
	t.Run("diagonal", func(t *testing.T) {
		AssertNumber(t, Vec(2, 2).Angle(), ToRadians(45))
		AssertNumber(t, Vec(0.6, -0.25).Angle(), -0.39479111)
	})
	t.Run("zero vector", func(t *testing.T) {
		AssertNumber(t, Vec(0, 0).Angle(), 0)
	})
}

func TestVector_Direction(t *testing.T) {
	t.Run("snaps to the nearest direction", func(t *testing.T) {
		assert.Equal(t, Vec(3, 0).Direction(), DirectionRight)
		assert.Equal(t, Vec(2.0, -2.0).Direction(), DirectionUpRight)
		assert.Equal(t, Vec(10.0, 1.0).Direction(), DirectionRight)
	})
	t.Run("zero vector has no direction", func(t *testing.T) {
		assert.Equal(t, Vec(0, 0).Direction(), DirectionNone)
	})
	t.Run("a vector shorter than Epsilon still has one", func(t *testing.T) {
		assert.Equal(t, Vec(-1e-7, 0.0).Direction(), DirectionLeft)
	})
	t.Run("NaN has no direction", func(t *testing.T) {
		assert.Equal(t, Vec(math.NaN(), 1.0).Direction(), DirectionNone)
	})
}

func TestVector_Lerp(t *testing.T) {
	t.Run("between vectors", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Lerp(Vec(20, 20), 0.25), Vec(13, 17))
		AssertVector(t, Vec(0.6, -0.25).Lerp(Vec(-10.0, 10.0), 0.25), Vec(-2.05, 2.3125))
	})
	t.Run("endpoints", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Lerp(Vec(20, 20), 0), Vec(10, 16))
		AssertVector(t, Vec(10, 16).Lerp(Vec(20, 20), 1), Vec(20, 20))
	})
	t.Run("extrapolates outside the unit range", func(t *testing.T) {
		AssertVector(t, Vec(0.0, 0.0).Lerp(Vec(2.0, 2.0), 2), Vec(4.0, 4.0))
		AssertVector(t, Vec(0.0, 0.0).Lerp(Vec(2.0, 2.0), -1), Vec(-2.0, -2.0))
	})
}

func TestVector_Equal(t *testing.T) {
	t.Run("same vector", func(t *testing.T) {
		assert.True(t, Vec(10, 16).Equal(Vec(10, 16)))
		assert.True(t, Vec(0.6, -0.25).Equal(Vec(0.6, -0.25)))
	})
	t.Run("different vector", func(t *testing.T) {
		assert.False(t, Vec(10, 16).Equal(Vec(3, -3)))
		assert.False(t, Vec(0.6, -0.25).Equal(Vec(100.1, -0.1)))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Vec(0.6, -0.25).Equal(Vec(0.6, -0.250001)))
	})
}

func TestVector_IsZero(t *testing.T) {
	t.Run("zero vector", func(t *testing.T) {
		assert.True(t, Vec(0, 0).IsZero())
		assert.True(t, ZeroVector[int]().IsZero())

		assert.True(t, Vec(0.0, 0.0).IsZero())
		assert.True(t, ZeroVector[float64]().IsZero())
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.True(t, Vec(-0, -0).IsZero())
		assert.True(t, Vec(negativeZero, negativeZero).IsZero())
	})
	t.Run("non-zero vector", func(t *testing.T) {
		assert.False(t, Vec(10, 16).IsZero())
		assert.False(t, Vec(0.6, -0.25).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Vec(0.0, 0.000001).IsZero())
	})
}

func TestVector_IsOne(t *testing.T) {
	t.Run("one vector", func(t *testing.T) {
		assert.True(t, OneVector[int]().IsOne())
		assert.True(t, OneVector[float64]().IsOne())
	})
	t.Run("other vector", func(t *testing.T) {
		assert.False(t, Vec(10, 16).IsOne())
		assert.False(t, Vec(0.6, -0.25).IsOne())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Vec(1.0, 1.000001).IsOne())
	})
}

func TestVector_IsUp(t *testing.T) {
	t.Run("upward directions", func(t *testing.T) {
		assert.True(t, DirectionUp.Unit[int]().IsUp())
		assert.True(t, DirectionUpLeft.Unit[float64]().IsUp())
		assert.True(t, DirectionUpRight.Unit[float64]().IsUp())
	})
	t.Run("other directions", func(t *testing.T) {
		assert.False(t, DirectionDown.Unit[int]().IsUp())
		assert.False(t, DirectionLeft.Unit[int]().IsUp())
		assert.False(t, DirectionRight.Unit[int]().IsUp())
		assert.False(t, DirectionDownLeft.Unit[float64]().IsUp())
		assert.False(t, DirectionDownRight.Unit[float64]().IsUp())
	})
	t.Run("degenerate vectors", func(t *testing.T) {
		assert.False(t, ZeroVector[int]().IsUp())
		assert.False(t, OneVector[int]().IsUp())
	})
}

func TestVector_IsDown(t *testing.T) {
	t.Run("downward directions", func(t *testing.T) {
		assert.True(t, DirectionDown.Unit[int]().IsDown())
		assert.True(t, DirectionDownLeft.Unit[float64]().IsDown())
		assert.True(t, DirectionDownRight.Unit[float64]().IsDown())
	})
	t.Run("other directions", func(t *testing.T) {
		assert.False(t, DirectionUp.Unit[int]().IsDown())
		assert.False(t, DirectionLeft.Unit[int]().IsDown())
		assert.False(t, DirectionRight.Unit[int]().IsDown())
		assert.False(t, DirectionUpLeft.Unit[float64]().IsDown())
		assert.False(t, DirectionUpRight.Unit[float64]().IsDown())
	})
	t.Run("degenerate vectors", func(t *testing.T) {
		assert.False(t, ZeroVector[int]().IsDown())
		assert.True(t, OneVector[int]().IsDown())
	})
}

func TestVector_IsLeft(t *testing.T) {
	t.Run("leftward directions", func(t *testing.T) {
		assert.True(t, DirectionLeft.Unit[int]().IsLeft())
		assert.True(t, DirectionUpLeft.Unit[float64]().IsLeft())
		assert.True(t, DirectionDownLeft.Unit[float64]().IsLeft())
	})
	t.Run("other directions", func(t *testing.T) {
		assert.False(t, DirectionUp.Unit[int]().IsLeft())
		assert.False(t, DirectionDown.Unit[int]().IsLeft())
		assert.False(t, DirectionRight.Unit[int]().IsLeft())
		assert.False(t, DirectionUpRight.Unit[float64]().IsLeft())
		assert.False(t, DirectionDownRight.Unit[float64]().IsLeft())
	})
	t.Run("degenerate vectors", func(t *testing.T) {
		assert.False(t, ZeroVector[int]().IsLeft())
		assert.False(t, OneVector[int]().IsLeft())
	})
}

func TestVector_IsRight(t *testing.T) {
	t.Run("rightward directions", func(t *testing.T) {
		assert.True(t, DirectionRight.Unit[int]().IsRight())
		assert.True(t, DirectionUpRight.Unit[float64]().IsRight())
		assert.True(t, DirectionDownRight.Unit[float64]().IsRight())
	})
	t.Run("other directions", func(t *testing.T) {
		assert.False(t, DirectionUp.Unit[int]().IsRight())
		assert.False(t, DirectionDown.Unit[int]().IsRight())
		assert.False(t, DirectionLeft.Unit[int]().IsRight())
		assert.False(t, DirectionUpLeft.Unit[float64]().IsRight())
		assert.False(t, DirectionDownLeft.Unit[float64]().IsRight())
	})
	t.Run("degenerate vectors", func(t *testing.T) {
		assert.False(t, ZeroVector[int]().IsRight())
		assert.True(t, OneVector[int]().IsRight())
	})
}

func TestVector_IsNormalized(t *testing.T) {
	t.Run("unit vectors", func(t *testing.T) {
		assert.True(t, Vec(1, 0).IsNormalized())
		assert.True(t, Vec(OneOverSqrt2, OneOverSqrt2).IsNormalized())
	})
	t.Run("after normalize", func(t *testing.T) {
		assert.True(t, Vec(10, 16).Normalize().IsNormalized())
		assert.True(t, Vec(1.1, 2.1).Normalize().IsNormalized())
	})
	t.Run("every direction unit", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.True(t, direction.Unit[float64]().IsNormalized(), direction.String()+": ")
		}
	})
	t.Run("other vectors", func(t *testing.T) {
		assert.False(t, Vec(10, 16).IsNormalized())
		assert.False(t, Vec(0, 0).IsNormalized())
	})
}

func TestVector_Less(t *testing.T) {
	t.Run("shorter than the value", func(t *testing.T) {
		assert.True(t, Vec(10, 16).Less(19))
		assert.True(t, Vec(0.6, -0.25).Less(0.7))
	})
	t.Run("longer than the value", func(t *testing.T) {
		assert.False(t, Vec(10, 16).Less(18))
		assert.False(t, Vec(0.6, -0.25).Less(0.1))
	})
	t.Run("nothing is shorter than a non-positive length", func(t *testing.T) {
		assert.False(t, Vec(10, 16).Less(-19))
		assert.False(t, ZeroVector[int]().Less(0))
	})
	t.Run("strict without tolerance", func(t *testing.T) {
		assert.False(t, Vec(3.0, 4.0).Less(5))
		assert.False(t, Vec(3.0, 4.0).Less(5-Delta/2))
	})
}

func TestVector_LessOrEqual(t *testing.T) {
	t.Run("shorter or equal", func(t *testing.T) {
		assert.True(t, Vec(3, 4).LessOrEqual(5))
		assert.True(t, Vec(3, 4).LessOrEqual(6))
		assert.True(t, ZeroVector[float64]().LessOrEqual(0))
	})
	t.Run("longer", func(t *testing.T) {
		assert.False(t, Vec(3, 4).LessOrEqual(4))
		assert.False(t, Vec(0.6, -0.25).LessOrEqual(0.1))
	})
	t.Run("nothing is at most a negative length", func(t *testing.T) {
		assert.False(t, ZeroVector[int]().LessOrEqual(-1))
	})
	t.Run("float allows Delta past the length", func(t *testing.T) {
		assert.True(t, Vec(3.0, 4.0).LessOrEqual(5-Delta/2))
		assert.False(t, Vec(3.0, 4.0).LessOrEqual(5-2*Delta))
	})
}

func TestVector_XY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		x, y := Vec(10, 16).XY()
		AssertNumber(t, x, 10)
		AssertNumber(t, y, 16)
	})
	t.Run("float", func(t *testing.T) {
		x, y := Vec(0.6, -0.25).XY()
		AssertNumber(t, x, 0.6)
		AssertNumber(t, y, -0.25)
	})
}

func TestVector_Point(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertPoint(t, Vec(10, 16).Point(), Pt(10, 16))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Vec(0.6, -0.25).Point(), Pt(0.6, -0.25))
	})
}

func TestVector_Size(t *testing.T) {
	t.Run("positive components", func(t *testing.T) {
		AssertSize(t, Vec(10, 16).Size(), Sz(10, 16))
	})
	t.Run("negative components become absolute", func(t *testing.T) {
		AssertSize(t, Vec(0.6, -0.25).Size(), Sz(0.6, 0.25))
	})
}

func TestVector_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Int(), Vec(10, 16))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Int(), Vec(1, 0))
		AssertVector(t, Vec(-1.5, 2.5).Int(), Vec(-2, 3))
	})
}

func TestVector_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertVector(t, Vec(10, 16).Float(), Vec(10.0, 16.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertVector(t, Vec(0.6, -0.25).Float(), Vec(0.6, -0.25))
	})
}

func TestVector_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Vec(10, 16).String(), "⟨10,16⟩")
		assert.Equal(t, Vec(-4, 0).String(), "⟨-4,0⟩")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Vec(100, -34.0000115).String(), "⟨100.00,-34.00⟩")
		assert.Equal(t, Vec(1.5, -0.25).String(), "⟨1.50,-0.25⟩")
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.Equal(t, Vec(-0, 0).String(), "⟨0,0⟩")
		assert.Equal(t, Vec(negativeZero, 0.0).String(), "⟨0.00,0.00⟩")
	})
}

func TestVector_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Vec(10, 16), `{"x":10,"y":16}`)

		var v Vector[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &v))
		AssertVector(t, v, Vec(10, 16))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Vec(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)

		var v Vector[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &v))
		AssertVector(t, v, Vec(10.1, -34.0000115))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			data, err := json.Marshal(vector)
			assert.NoError(t, err)

			var decoded Vector[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, vector)
		}
	})
}

func TestVector_Properties(t *testing.T) {
	t.Run("add and subtract are inverse", func(t *testing.T) {
		for _, a := range vectorFixtures {
			for _, b := range vectorFixtures {
				assert.True(t, a.Add(b).Subtract(b).Equal(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("negate is its own inverse", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			assert.True(t, vector.Negate().Negate().Equal(vector), fmt.Sprintf("%s: ", vector))
		}
	})
	t.Run("dot with itself is the squared length", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			AssertNumber(t, vector.Dot(vector), vector.LengthSquared(), fmt.Sprintf("%s: ", vector))
			AssertNumber(t, vector.Length()*vector.Length(), vector.LengthSquared(), fmt.Sprintf("%s: ", vector))
		}
	})
	t.Run("normal is perpendicular and keeps the length", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			normal := vector.Normal()

			AssertNumber(t, vector.Dot(normal), 0, fmt.Sprintf("%s: ", vector))
			AssertNumber(t, normal.Length(), vector.Length(), fmt.Sprintf("%s: ", vector))
			assert.True(t, normal.Equal(vector.Rotate(Pi/2)), fmt.Sprintf("%s: ", vector))
		}
	})
	t.Run("rotate preserves the length", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			for _, angle := range []float64{0, Pi / 6, Pi / 2, 2, -1.5} {
				AssertNumber(t, vector.Rotate(angle).Length(), vector.Length(), fmt.Sprintf("%s ∠%v: ", vector, angle))
			}
		}
	})
	t.Run("an int normal is an axis unit pointing forwards", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			integer := vector.Int()
			if integer.IsZero() {
				continue // the zero vector normalizes to (1,0) by convention
			}

			unit := integer.Normalize()

			assert.True(t, unit.IsNormalized(), fmt.Sprintf("%s: ", integer))
			assert.True(t, unit.Dot(integer) > 0, fmt.Sprintf("%s: ", integer))
		}
	})
	t.Run("normalize keeps the angle at unit length", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			if vector.IsZero() {
				continue // the zero vector normalizes to (1,0) by convention
			}

			assert.True(t, vector.Normalize().IsNormalized(), fmt.Sprintf("%s: ", vector))
			AssertNumber(t, vector.Normalize().Angle(), vector.Angle(), fmt.Sprintf("%s: ", vector))
		}
	})
	t.Run("angle and length reconstruct the vector", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			rebuilt := VectorFromAngle(vector.Angle(), vector.Length())
			assert.True(t, rebuilt.Equal(vector), fmt.Sprintf("%s: ", vector))
		}
	})
	t.Run("cross is antisymmetric", func(t *testing.T) {
		for _, a := range vectorFixtures {
			for _, b := range vectorFixtures {
				AssertNumber(t, a.Cross(b), -b.Cross(a), fmt.Sprintf("%s → %s: ", a, b))
				AssertNumber(t, a.Dot(b), b.Dot(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("identity matrix leaves the vector untouched", func(t *testing.T) {
		for _, vector := range vectorFixtures {
			assert.True(t, vector.Transform(IdentityMatrix[float64]()).Equal(vector), fmt.Sprintf("%s: ", vector))
		}
	})
}

func TestVector_Immutable(t *testing.T) {
	v1 := Vec(10, 16)
	v2 := Vec(3, -3)

	v1.Transform(IdentityMatrix[float64]())
	v1.Add(Vec(3, -2))
	v1.AddXY(3, -2)
	v1.Subtract(v2)
	v1.SubtractXY(3, -3)
	v1.Multiply(2)
	v1.MultiplyXY(3, 4)
	v1.Divide(5)
	v1.DivideXY(10, 100)
	v1.Negate()
	v1.Rotate(Pi / 3)
	v1.Resize(4)
	v1.Normalize()
	v1.Abs()
	v1.Round()
	v1.Floor()
	v1.Ceil()
	v1.Normal()
	v1.Lerp(v2, 0.1)

	AssertVector(t, v1, Vec(10, 16))
	AssertVector(t, v2, Vec(3, -3))
}

// vectorFixtures span the quadrants, the axes, and the diagonal.
var vectorFixtures = []Vector[float64]{
	Vec(0.0, 0.0),
	Vec(1.0, 2.0),
	Vec(-3.5, 0.25),
	Vec(12.5, -0.1),
	Vec(-1.0, -1.0),
	Vec(0.6, -0.25),
	Vec(7.0, 7.0),
	Vec(-0.5, 12.75),
}

func ExampleVec() {
	fmt.Println(Vec(3, 4))
	fmt.Println(Vec(1.5, -0.25))
	// Output:
	// ⟨3,4⟩
	// ⟨1.50,-0.25⟩
}

func ExampleVector_Length() {
	fmt.Printf("%.2f\n", Vec(3.0, 4.0).Length())
	// Output: 5.00
}

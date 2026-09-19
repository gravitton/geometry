package geom

import (
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

// negativeZero is -0.0; writing it as a literal would fold to +0.0 at compile time.
var negativeZero = math.Copysign(0, -1)

// sinkPoints is a shared slice variable used to store results in allocation tests, so the result escapes as it does for any caller.
var sinkPoints []Point[int]

// sinkBool is a shared boolean variable used to store results in benchmarking tests.
var sinkBool bool

func TestMultiply(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		AssertNumber(t, Multiply(5, 2.0), 10)
		AssertNumber(t, Multiply(3, 0.5), 2) // int: 1.5 rounds to 2
		AssertNumber(t, Multiply(-4, 1.5), -6)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Multiply(4.0, 2.5), 10.0)
		AssertNumber(t, Multiply(0.4, 0.5), 0.2)
	})
}

func TestDivide(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		AssertNumber(t, Divide(10, 2.0), 5)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Divide(7.5, 2.5), 3.0)
		AssertNumber(t, Divide(1.0, 3.0), 1.0/3.0)
	})
	t.Run("zero scale panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Divide(10, 0.0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Divide(2.5, 0.0)
		}, "geom: division by zero")
	})
}

func TestAbs(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		AssertNumber(t, Abs(-5), 5)
		AssertNumber(t, Abs(-3.14), 3.14)
	})
	t.Run("non-negative is unchanged", func(t *testing.T) {
		AssertNumber(t, Abs(5), 5)
		AssertNumber(t, Abs(0), 0)
		AssertNumber(t, Abs(3.14), 3.14)
	})
	t.Run("wide integers stay exact", func(t *testing.T) {
		AssertNumber(t, Abs(int64(-(1<<53 + 1))), int64(1<<53+1))
		AssertNumber(t, Abs(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestSign(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Sign(5), 1)
		AssertNumber(t, Sign(-5), -1)
		AssertNumber(t, Sign(0), 0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Sign(3.14), 1.0)
		AssertNumber(t, Sign(-3.14), -1.0)
		AssertNumber(t, Sign(0.0), 0.0)
	})
}

func TestRound(t *testing.T) {
	t.Run("float rounds half away from zero", func(t *testing.T) {
		AssertNumber(t, Round(1.4), 1.0)
		AssertNumber(t, Round(1.5), 2.0)
		AssertNumber(t, Round(-1.5), -2.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertNumber(t, Round(3), 3)
		AssertNumber(t, Round(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestFloor(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Floor(1.9), 1.0)
		AssertNumber(t, Floor(-1.1), -2.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertNumber(t, Floor(3), 3)
		AssertNumber(t, Floor(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestCeil(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ceil(1.1), 2.0)
		AssertNumber(t, Ceil(-1.9), -1.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertNumber(t, Ceil(3), 3)
		AssertNumber(t, Ceil(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestMod(t *testing.T) {
	t.Run("non-negative operands", func(t *testing.T) {
		AssertNumber(t, Mod(5, 8), 5)
		AssertNumber(t, Mod(8, 8), 0)
		AssertNumber(t, Mod(9, 8), 1)
		AssertNumber(t, Mod(0, 8), 0)
		AssertNumber(t, Mod(7, 6), 1)
	})
	t.Run("negative operands stay in range", func(t *testing.T) {
		AssertNumber(t, Mod(-1, 8), 7)
		AssertNumber(t, Mod(-8, 8), 0)
		AssertNumber(t, Mod(-9, 8), 7)
		AssertNumber(t, Mod(-1, 6), 5)
	})
	t.Run("a negative modulus wraps into (m, 0]", func(t *testing.T) {
		AssertNumber(t, Mod(3, -5), -2)
		AssertNumber(t, Mod(-3, -5), -3)
		AssertNumber(t, Mod(5, -5), 0)
	})
	t.Run("defined integer types", func(t *testing.T) {
		AssertNumber(t, Mod(namedInt(12), 8), namedInt(4))
	})
}

func TestClamp(t *testing.T) {
	t.Run("inside the range", func(t *testing.T) {
		AssertNumber(t, Clamp(5, 0, 10), 5)
		AssertNumber(t, Clamp(0.5, 0.0, 1.0), 0.5)
	})
	t.Run("outside the range", func(t *testing.T) {
		AssertNumber(t, Clamp(-5, 0, 10), 0)
		AssertNumber(t, Clamp(15, 0, 10), 10)
		AssertNumber(t, Clamp(-0.5, 0.0, 1.0), 0.0)
		AssertNumber(t, Clamp(1.5, 0.0, 1.0), 1.0)
	})
	t.Run("the bounds are inclusive", func(t *testing.T) {
		AssertNumber(t, Clamp(0, 0, 10), 0)
		AssertNumber(t, Clamp(10, 0, 10), 10)
	})
}

func TestSum(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Sum([]int{1, 2, 3}), 6)
		assert.Equal(t, Sum([]int{-4, 1}), -3)
	})
	t.Run("narrow integers do not overflow mid-sum", func(t *testing.T) {
		assert.Equal(t, Sum([]int8{100, 100, -100}), int8(100))
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Sum([]float64{0.5, 0.25, 0.125}), 0.875)
		AssertNumber(t, Sum([]float32{0.5, 0.25}), float32(0.75))
	})
	t.Run("empty and nil are zero", func(t *testing.T) {
		assert.Equal(t, Sum([]int{}), 0)
		assert.Equal(t, Sum[float64](nil), 0.0)
	})
}

func TestLerp(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		AssertNumber(t, Lerp(1, 2, 0.25), 1)
		AssertNumber(t, Lerp(1, 3, 0.25), 2)
		AssertNumber(t, Lerp(1, 5, 0.25), 2)
		AssertNumber(t, Lerp(1, 7, 0.25), 3)
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		AssertNumber(t, Lerp[int8](-100, 100, 0.5), 0)
		AssertNumber(t, Lerp[int8](-128, 127, 1), 127)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Lerp(1.0, 6.0, 0.25), 2.25)
		AssertNumber(t, Lerp(1.0, 6.0, 0.75), 4.75)
	})
	t.Run("endpoints", func(t *testing.T) {
		AssertNumber(t, Lerp(1.0, 6.0, 0), 1.0)
		AssertNumber(t, Lerp(1.0, 6.0, 1), 6.0)
	})
}

func TestMidpoint(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		AssertNumber(t, Midpoint(1, 3), 2)
		AssertNumber(t, Midpoint(1, 4), 3)
		AssertNumber(t, Midpoint(1, 5), 3)
		AssertNumber(t, Midpoint(1, 6), 4)
		AssertNumber(t, Midpoint(1, 7), 4)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Midpoint(1.0, 6.0), 3.5)
	})
	t.Run("is lerp at one half", func(t *testing.T) {
		for _, pair := range [][2]float64{{1, 6}, {-3.5, 0.25}, {7, 7}, {12.5, -0.1}} {
			AssertNumber(t, Midpoint(pair[0], pair[1]), Lerp(pair[0], pair[1], 0.5))
		}
	})
}

func TestEqual(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, Equal(1, 1))
		assert.False(t, Equal(1, 2))
	})
	t.Run("float64 within Delta", func(t *testing.T) {
		assert.True(t, Equal(1.0000005, 1.0))
		assert.False(t, Equal(1.0000015, 1.0))
		assert.True(t, Equal(0.0, 0.0))
	})
	t.Run("float32 within Delta32", func(t *testing.T) {
		assert.True(t, Equal[float32](1.00005, 1.0))
		assert.False(t, Equal[float32](1.0005, 1.0))
	})
	t.Run("absolute, so it degenerates far from zero", func(t *testing.T) {
		// one float32 ulp at 1e4 is 9.8e-4, above Delta32, so Equal asks for bit-exact
		// equality from here up. This is the documented cost of the fast path.
		const far float32 = 1e4
		next := math.Nextafter32(far, math.MaxFloat32)

		assert.False(t, Equal(far, next))
		assert.True(t, EqualRelative(far, next))
	})
}

func BenchmarkEqual_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(float64(i), float64(i))
	}
}

func BenchmarkEqual_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(float32(i), float32(i))
	}
}

func TestEqualDelta(t *testing.T) {
	t.Run("within the delta", func(t *testing.T) {
		assert.True(t, EqualDelta(1, 2, 1.5))
		assert.True(t, EqualDelta(1.0, 1.001, 0.01))
	})
	t.Run("outside the delta", func(t *testing.T) {
		assert.False(t, EqualDelta(1, 3, 1.5))
		assert.False(t, EqualDelta(1.0, 1.02, 0.01))
	})
	t.Run("a zero delta asks for exact equality", func(t *testing.T) {
		assert.True(t, EqualDelta(5, 5, 0.0))
		assert.False(t, EqualDelta(5, 6, 0.0))
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		assert.False(t, EqualDelta[int8](127, -128, 1))
		assert.True(t, EqualDelta[int8](127, -128, 255))
	})
}

func BenchmarkEqualDelta_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = EqualDelta(float64(i), float64(i), Delta)
	}
}

func TestEqualRelative(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, EqualRelative(1, 1))
		assert.False(t, EqualRelative(1, 2))
	})
	t.Run("near zero it matches Equal", func(t *testing.T) {
		assert.True(t, EqualRelative(1.0000005, 1.0))
		assert.False(t, EqualRelative(1.0000015, 1.0))
		assert.True(t, EqualRelative(0.0, 0.0))
	})
	t.Run("scales with magnitude", func(t *testing.T) {
		assert.True(t, EqualRelative(1e6, 1e6+0.5))
		assert.False(t, EqualRelative(1e6, 1e6+2.0))
		assert.True(t, EqualRelative[float32](1e4, 1e4+0.05))
	})
	t.Run("a full turn in float32 returns to its start", func(t *testing.T) {
		v := Vec[float32](1000, 0)
		for range 8 {
			v = v.Rotate(Pi / 4)
		}

		// one ulp out at magnitude 1e3, which Delta32 still covers
		assert.True(t, EqualRelative(v.X, 1000))
		assert.True(t, Equal(v.X, 1000))
	})
}

func BenchmarkEqualRelative_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = EqualRelative(float64(i), float64(i))
	}
}

func TestLessOrEqual(t *testing.T) {
	t.Run("int compares exactly", func(t *testing.T) {
		assert.True(t, LessOrEqual(1, 2))
		assert.True(t, LessOrEqual(2, 2))
		assert.False(t, LessOrEqual(3, 2))
	})
	t.Run("int64 beyond 2^53 stays exact", func(t *testing.T) {
		big := int64(1) << 53

		assert.False(t, LessOrEqual(big+1, big))
		assert.True(t, LessOrEqual(big, big+1))
	})
	t.Run("float64 allows Delta past the bound", func(t *testing.T) {
		assert.True(t, LessOrEqual(1.0, 2.0))
		assert.True(t, LessOrEqual(2.0+Delta/2, 2.0))
		assert.False(t, LessOrEqual(2.0+2*Delta, 2.0))
	})
	t.Run("float32 allows Delta32 past the bound", func(t *testing.T) {
		assert.True(t, LessOrEqual(float32(2+Delta32/2), 2))
		assert.False(t, LessOrEqual(float32(2+2*Delta32), 2))
	})
	t.Run("a recomputed float corner is still inside", func(t *testing.T) {
		low := 0.1
		recomputed := (low + 0.6/2) - 0.6/2

		assert.False(t, low >= recomputed)
		assert.True(t, LessOrEqual(recomputed, low))
	})
}

func TestLessOrEqualDelta(t *testing.T) {
	t.Run("within the delta", func(t *testing.T) {
		assert.True(t, LessOrEqualDelta(1, 2, 0.0))
		assert.True(t, LessOrEqualDelta(3, 2, 1.5))
		assert.True(t, LessOrEqualDelta(1.001, 1.0, 0.01))
	})
	t.Run("outside the delta", func(t *testing.T) {
		assert.False(t, LessOrEqualDelta(4, 2, 1.5))
		assert.False(t, LessOrEqualDelta(1.02, 1.0, 0.01))
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		assert.True(t, LessOrEqualDelta[int8](-128, 127, 0))
		assert.False(t, LessOrEqualDelta[int8](127, -128, 254))
		assert.True(t, LessOrEqualDelta[int8](127, -128, 255))
	})
}

func TestEpsilon(t *testing.T) {
	t.Run("int compares exactly", func(t *testing.T) {
		AssertNumber(t, Epsilon[int](), 0.0)
		AssertNumber(t, Epsilon[int8](), 0.0)
		AssertNumber(t, Epsilon[namedInt](), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Epsilon[float32](), Delta32)
		AssertNumber(t, Epsilon[float64](), Delta)
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		AssertNumber(t, Epsilon[namedFloat32](), Delta32)
		AssertNumber(t, Epsilon[namedFloat64](), Delta)
	})
}

func TestEpsilonRelative(t *testing.T) {
	t.Run("floored at Epsilon near zero", func(t *testing.T) {
		AssertNumber(t, EpsilonRelative(0.0, 0.5), Delta)
	})
	t.Run("scales with the larger magnitude", func(t *testing.T) {
		AssertNumber(t, EpsilonRelative(100.0, 0.0), 100*Delta)
		AssertNumber(t, EpsilonRelative(0.0, -100.0), 100*Delta)
		AssertNumber(t, EpsilonRelative[float32](100, 0), 100*Delta32)
	})
	t.Run("int stays exact at any magnitude", func(t *testing.T) {
		AssertNumber(t, EpsilonRelative(5, 5), 0.0)
		AssertNumber(t, EpsilonRelative(1000000, 0), 0.0)
	})
}

func TestToRadians(t *testing.T) {
	AssertNumber(t, ToRadians(0), 0.0)
	AssertNumber(t, ToRadians(90), Pi/2)
	AssertNumber(t, ToRadians(180), Pi)
	AssertNumber(t, ToRadians(360), 2*Pi)
}

func TestToDegrees(t *testing.T) {
	t.Run("converts", func(t *testing.T) {
		AssertNumber(t, ToDegrees(0), 0.0)
		AssertNumber(t, ToDegrees(Pi/2), 90.0)
		AssertNumber(t, ToDegrees(Pi), 180.0)
		AssertNumber(t, ToDegrees(2*Pi), 360.0)
	})
	t.Run("inverts ToRadians", func(t *testing.T) {
		for _, degrees := range []float64{0, 30, 45, 90, 179.5, -270} {
			AssertNumber(t, ToDegrees(ToRadians(degrees)), degrees)
		}
	})
}

func TestNormalizeAngle(t *testing.T) {
	t.Run("inside the unit turn", func(t *testing.T) {
		AssertNumber(t, NormalizeAngle(0), 0.0)
		AssertNumber(t, NormalizeAngle(Pi), Pi)
	})
	t.Run("wraps above a full turn", func(t *testing.T) {
		AssertNumber(t, NormalizeAngle(3*Pi), Pi)
		AssertNumber(t, NormalizeAngle(2*Pi), 0.0)
	})
	t.Run("negative angles come back positive", func(t *testing.T) {
		AssertNumber(t, NormalizeAngle(-Pi/2), 3*Pi/2)
	})
	t.Run("a tiny negative angle stays below a full turn", func(t *testing.T) {
		assert.Less(t, NormalizeAngle(-1e-17), 2*Pi)
		assert.GreaterOrEqual(t, NormalizeAngle(-1e-17), 0.0)
	})
	t.Run("NaN and Inf have no normal form", func(t *testing.T) {
		assert.True(t, math.IsNaN(NormalizeAngle(math.NaN())))
		assert.True(t, math.IsNaN(NormalizeAngle(math.Inf(1))))
		assert.True(t, math.IsNaN(NormalizeAngle(math.Inf(-1))))
	})
}

func TestAngleDistance(t *testing.T) {
	t.Run("same angle", func(t *testing.T) {
		AssertNumber(t, AngleDistance(1, 1), 0.0)
		AssertNumber(t, AngleDistance(1, 1+2*Pi), 0.0)
	})
	t.Run("shortest way around", func(t *testing.T) {
		AssertNumber(t, AngleDistance(0, Pi/2), Pi/2)
		AssertNumber(t, AngleDistance(Pi/2, 0), Pi/2)
		AssertNumber(t, AngleDistance(0, 3*Pi/2), Pi/2)
		AssertNumber(t, AngleDistance(0, Pi), Pi)
	})
	t.Run("across the seam", func(t *testing.T) {
		AssertNumber(t, AngleDistance(-0.1, 0.1), 0.2)
		AssertNumber(t, AngleDistance(2*Pi-0.1, 0.1), 0.2)
	})
}

func TestLerpAngle(t *testing.T) {
	t.Run("turns along the shorter arc", func(t *testing.T) {
		AssertNumber(t, LerpAngle(0, Pi/2, 0.5), Pi/4)
		AssertNumber(t, LerpAngle(Pi/2, 0, 0.5), Pi/4)
		AssertNumber(t, LerpAngle(0, 3*Pi/2, 0.5), -Pi/4)
	})
	t.Run("crosses the seam the short way", func(t *testing.T) {
		AssertNumber(t, LerpAngle(ToRadians(350), ToRadians(10), 0.5), ToRadians(360))
		AssertNumber(t, LerpAngle(ToRadians(10), ToRadians(350), 0.5), ToRadians(0))
	})
	t.Run("the ends are the angles themselves, up to a turn", func(t *testing.T) {
		assert.True(t, EqualAngle(LerpAngle(1, 4, 0), 1))
		assert.True(t, EqualAngle(LerpAngle(1, 4, 1), 4))
		assert.True(t, EqualAngle(LerpAngle(1, 4+2*Pi, 1), 4))
	})
	t.Run("extrapolates along the same arc", func(t *testing.T) {
		AssertNumber(t, LerpAngle(0, Pi/2, 2), Pi)
		AssertNumber(t, LerpAngle(0, Pi/2, -1), -Pi/2)
	})
	t.Run("half a turn apart turns by increasing angle", func(t *testing.T) {
		AssertNumber(t, LerpAngle(0, Pi, 0.5), Pi/2)
	})
	t.Run("never turns more than half a turn", func(t *testing.T) {
		for _, a := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
			for _, b := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
				assert.True(t, math.Abs(LerpAngle(a, b, 1)-a) <= Pi+Delta, fmt.Sprintf("%v → %v: ", a, b))
				assert.True(t, EqualAngle(LerpAngle(a, b, 1), b), fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
}

func TestEqualAngle(t *testing.T) {
	t.Run("same angle", func(t *testing.T) {
		assert.True(t, EqualAngle(0.5, 0.5))
		assert.True(t, EqualAngle(0.5, 0.5+2*Pi))
		assert.True(t, EqualAngle(-Pi/2, 3*Pi/2))
	})
	t.Run("across the seam", func(t *testing.T) {
		assert.True(t, EqualAngle(0, -1e-9))
		assert.True(t, EqualAngle(0, 2*Pi-1e-9))
		assert.True(t, EqualAngle(-1e-9, 1e-9))
	})
	t.Run("different angle", func(t *testing.T) {
		assert.False(t, EqualAngle(0, 1e-3))
		assert.False(t, EqualAngle(0, Pi))
		assert.False(t, EqualAngle(0, math.NaN()))
	})
}

func TestParse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v, err := Parse[int]("42")
		assert.NoError(t, err)
		AssertNumber(t, v, 42)
	})
	t.Run("every integer width", func(t *testing.T) {
		v16, err := Parse[int16]("32000")
		assert.NoError(t, err)
		AssertNumber(t, v16, int16(32000))

		v32, err := Parse[int32]("2147483647")
		assert.NoError(t, err)
		AssertNumber(t, v32, int32(2147483647))

		v64, err := Parse[int64]("9223372036854775807")
		assert.NoError(t, err)
		AssertNumber(t, v64, int64(9223372036854775807))
	})
	t.Run("float", func(t *testing.T) {
		v32, err := Parse[float32]("3.14")
		assert.NoError(t, err)
		AssertNumber(t, v32, float32(3.14))

		v64, err := Parse[float64]("23.0")
		assert.NoError(t, err)
		AssertNumber(t, v64, 23.0)
	})
	t.Run("defined types parse like their underlying type", func(t *testing.T) {
		direction, err := Parse[Direction]("3")
		assert.NoError(t, err)
		assert.Equal(t, direction, DirectionDownLeft)

		named, err := Parse[namedFloat32]("3.14")
		assert.NoError(t, err)
		AssertNumber(t, named, namedFloat32(3.14))
	})
	t.Run("int rejects float strings", func(t *testing.T) {
		_, err := Parse[int]("3.14")
		assert.Error(t, err)
	})
	t.Run("the range is checked against T", func(t *testing.T) {
		_, err := Parse[int8]("200")
		assert.Error(t, err)

		_, err = Parse[namedInt8]("200")
		assert.Error(t, err)

		_, err = Parse[float32]("1e39")
		assert.Error(t, err)

		v, err := Parse[float64]("1e39")
		assert.NoError(t, err)
		AssertNumber(t, v, 1e39)
	})
	t.Run("non-numeric input", func(t *testing.T) {
		_, err := Parse[int]("abc")
		assert.Error(t, err)

		_, err = Parse[float64]("abc")
		assert.Error(t, err)
	})
}

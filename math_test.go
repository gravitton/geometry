package geom

import (
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestNormalizeAngle(t *testing.T) {
	assert.EqualDelta(t, NormalizeAngle(0), 0.0, Delta)
	assert.EqualDelta(t, NormalizeAngle(Pi), Pi, Delta)
	assert.EqualDelta(t, NormalizeAngle(3*Pi), Pi, Delta)      // wraps: 3π → π
	assert.EqualDelta(t, NormalizeAngle(-Pi/2), 3*Pi/2, Delta) // negative → 3π/2
	assert.EqualDelta(t, NormalizeAngle(2*Pi), 0.0, Delta)     // exactly 2π → 0
}

func TestToRadians(t *testing.T) {
	assert.EqualDelta(t, ToRadians(0), 0.0, Delta)
	assert.EqualDelta(t, ToRadians(90), Pi/2, Delta)
	assert.EqualDelta(t, ToRadians(180), Pi, Delta)
	assert.EqualDelta(t, ToRadians(360), 2*Pi, Delta)
}

func TestToDegrees(t *testing.T) {
	assert.EqualDelta(t, ToDegrees(0), 0.0, Delta)
	assert.EqualDelta(t, ToDegrees(Pi/2), 90.0, Delta)
	assert.EqualDelta(t, ToDegrees(Pi), 180.0, Delta)
	assert.EqualDelta(t, ToDegrees(2*Pi), 360.0, Delta)
}

func TestMultiply(t *testing.T) {
	assert.Equal(t, Multiply(5, 2.0), 10)
	assert.Equal(t, Multiply(3, 0.5), 2) // int: 1.5 rounds to 2
	assert.Equal(t, Multiply(-4, 1.5), -6)
	assert.EqualDelta(t, Multiply(4.0, 2.5), 10.0, Delta)
	assert.EqualDelta(t, Multiply(0.4, 0.5), 0.2, Delta)
}

func TestDivide(t *testing.T) {
	assert.Equal(t, Divide(10, 2.0), 5)
	assert.Equal(t, Divide(10, 0.0), 10) // zero guard: returns original
	assert.EqualDelta(t, Divide(7.5, 2.5), 3.0, Delta)
	assert.EqualDelta(t, Divide(1.0, 3.0), 0.333333, Delta)
}

func TestAbs(t *testing.T) {
	assert.Equal(t, Abs(-5), 5)
	assert.Equal(t, Abs(5), 5)
	assert.Equal(t, Abs(0), 0)
	assert.EqualDelta(t, Abs(-3.14), 3.14, Delta)
	assert.EqualDelta(t, Abs(3.14), 3.14, Delta)
}

func TestRound(t *testing.T) {
	assert.Equal(t, Round(1.4), 1.0)
	assert.Equal(t, Round(1.5), 2.0)
	assert.Equal(t, Round(-1.5), -2.0)
	assert.Equal(t, Round(3), 3) // int: no-op
}

func TestFloor(t *testing.T) {
	assert.Equal(t, Floor(1.9), 1.0)
	assert.Equal(t, Floor(-1.1), -2.0)
	assert.Equal(t, Floor(3), 3) // int: no-op
}

func TestCeil(t *testing.T) {
	assert.Equal(t, Ceil(1.1), 2.0)
	assert.Equal(t, Ceil(-1.9), -1.0)
	assert.Equal(t, Ceil(3), 3) // int: no-op
}

func TestMod(t *testing.T) {
	assert.Equal(t, Mod(5, 8), 5)
	assert.Equal(t, Mod(8, 8), 0)
	assert.Equal(t, Mod(9, 8), 1)
	assert.Equal(t, Mod(-1, 8), 7)
	assert.Equal(t, Mod(-8, 8), 0)
	assert.Equal(t, Mod(-9, 8), 7)
	assert.Equal(t, Mod(0, 8), 0)

	assert.Equal(t, Mod(7, 6), 1)
	assert.Equal(t, Mod(-1, 6), 5)

	type testInt int32
	assert.Equal(t, Mod(testInt(12), 8), testInt(4))
}

func TestMidpoint(t *testing.T) {
	assert.Equal(t, Midpoint(1, 3), 2)
	assert.Equal(t, Midpoint(1, 4), 3)
	assert.Equal(t, Midpoint(1, 5), 3)
	assert.Equal(t, Midpoint(1, 6), 4)
	assert.Equal(t, Midpoint(1, 7), 4)

	assert.Equal(t, Midpoint(1.0, 6.0), 3.5)
}

func TestLerp(t *testing.T) {
	assert.Equal(t, Lerp(1, 2, 0.25), 1)
	assert.Equal(t, Lerp(1, 3, 0.25), 2)
	assert.Equal(t, Lerp(1, 4, 0.25), 2)
	assert.Equal(t, Lerp(1, 5, 0.25), 2)
	assert.Equal(t, Lerp(1, 6, 0.25), 2)
	assert.Equal(t, Lerp(1, 7, 0.25), 3)

	assert.Equal(t, Lerp(1.0, 6.0, 0.25), 2.25)
	assert.Equal(t, Lerp(1.0, 6.0, 0.75), 4.75)
}

func TestSign(t *testing.T) {
	assert.Equal(t, Sign(5), 1)
	assert.Equal(t, Sign(-5), -1)
	assert.Equal(t, Sign(0), 0)
	assert.EqualDelta(t, Sign(3.14), 1.0, Delta)
	assert.EqualDelta(t, Sign(-3.14), -1.0, Delta)
	assert.EqualDelta(t, Sign(0.0), 0.0, Delta)
}

func TestClamp(t *testing.T) {
	assert.Equal(t, Clamp(5, 0, 10), 5)
	assert.Equal(t, Clamp(-5, 0, 10), 0)
	assert.Equal(t, Clamp(15, 0, 10), 10)
	assert.Equal(t, Clamp(0, 0, 10), 0)
	assert.Equal(t, Clamp(10, 0, 10), 10)
	assert.EqualDelta(t, Clamp(0.5, 0.0, 1.0), 0.5, Delta)
	assert.EqualDelta(t, Clamp(-0.5, 0.0, 1.0), 0.0, Delta)
	assert.EqualDelta(t, Clamp(1.5, 0.0, 1.0), 1.0, Delta)
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

func TestEpsilon(t *testing.T) {
	assert.Equal(t, Epsilon[int](), 0.0)
	assert.Equal(t, Epsilon[int8](), 0.0)
	assert.Equal(t, Epsilon[float32](), Delta32)
	assert.Equal(t, Epsilon[float64](), Delta)

	assert.Equal(t, Epsilon[namedInt](), 0.0)
	assert.Equal(t, Epsilon[namedFloat32](), Delta32)
	assert.Equal(t, Epsilon[namedFloat64](), Delta)
}

func TestEqualRelative(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, EqualRelative(1, 1))
		assert.False(t, EqualRelative(1, 2))
	})
	t.Run("near zero matches Equal", func(t *testing.T) {
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

func TestRelativeEpsilon(t *testing.T) {
	assert.Equal(t, EpsilonRelative(0.0, 0.5), Delta) // floored at Epsilon near zero
	assert.Equal(t, EpsilonRelative(100.0, 0.0), 100*Delta)
	assert.Equal(t, EpsilonRelative(0.0, -100.0), 100*Delta) // larger magnitude wins
	assert.Equal(t, EpsilonRelative[float32](100, 0), 100*Delta32)
	assert.Equal(t, EpsilonRelative(5, 5), 0.0) // int
}

func TestEqualDelta(t *testing.T) {
	assert.True(t, EqualDelta(1, 2, 1.5))
	assert.False(t, EqualDelta(1, 3, 1.5))
	assert.True(t, EqualDelta(1.0, 1.001, 0.01))
	assert.False(t, EqualDelta(1.0, 1.02, 0.01))
	assert.True(t, EqualDelta(5, 5, 0.0))
}

func TestParseNumber(t *testing.T) {
	// int — decimal only
	v1, err := Parse[int]("42")
	assert.NoError(t, err)
	assert.Equal(t, v1, 42)

	// int8 — respects 8-bit overflow
	_, err = Parse[int8]("200")
	assert.Error(t, err)

	// int16
	v2, err := Parse[int16]("32000")
	assert.NoError(t, err)
	assert.Equal(t, v2, int16(32000))

	// int32
	v3, err := Parse[int32]("2147483647")
	assert.NoError(t, err)
	assert.Equal(t, v3, int32(2147483647))

	// int64
	v4, err := Parse[int64]("9223372036854775807")
	assert.NoError(t, err)
	assert.Equal(t, v4, int64(9223372036854775807))

	// int rejects float strings
	_, err = Parse[int]("3.14")
	assert.Error(t, err)

	// float32
	v5, err := Parse[float32]("3.14")
	assert.NoError(t, err)
	assert.EqualDelta(t, float64(v5), 3.14, 1e-5)

	// float64
	v6, err := Parse[float64]("23.0")
	assert.NoError(t, err)
	assert.EqualDelta(t, v6, 23.0, Delta)

	// error: non-numeric
	_, err = Parse[int]("abc")
	assert.Error(t, err)

	_, err = Parse[float64]("abc")
	assert.Error(t, err)

	// defined types over int/float parse like their underlying type
	v7, err := Parse[Direction]("3")
	assert.NoError(t, err)
	assert.Equal(t, v7, DirectionDownLeft)

	v8, err := Parse[namedFloat32]("3.14")
	assert.NoError(t, err)
	assert.EqualDelta(t, float64(v8), 3.14, 1e-5)

	// range is checked against T, not int64/float64
	_, err = Parse[namedInt8]("200")
	assert.Error(t, err)

	_, err = Parse[float32]("1e39")
	assert.Error(t, err)

	v9, err := Parse[float64]("1e39")
	assert.NoError(t, err)
	assert.Equal(t, v9, 1e39)
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

func BenchmarkEqualRelative_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = EqualRelative(float64(i), float64(i))
	}
}

var sinkBool bool

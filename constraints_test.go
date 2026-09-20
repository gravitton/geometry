package geom

import (
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

// Named types over the constraint's underlying types: the ~ in Number admits them,
// so every type-level predicate has to recognize them as their underlying kind.
type (
	namedInt     int
	namedInt8    int8
	namedFloat32 float32
	namedFloat64 float64
)

func TestCast(t *testing.T) {
	t.Run("int rounds half away from zero", func(t *testing.T) {
		assertCast[int](t, 1, 1)
		assertCast[int](t, 1.2, 1)
		assertCast[int](t, 1.6, 2)
		assertCast[int](t, -0.3, 0)
		assertCast[int](t, -0.51, -1)
	})
	t.Run("float keeps the fraction", func(t *testing.T) {
		assertCast[float64](t, 1.0, 1.0)
		assertCast[float64](t, 1.6, 1.6)
		assertCast[float64](t, -15.68, -15.68)
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		assertCast[namedInt](t, 1.6, 2)
		assertCast[namedInt](t, -0.51, -1)
		assertCast[namedInt8](t, 1.2, 1)
		assertCast[namedFloat64](t, 1.6, 1.6)
		assertCast[namedFloat32](t, -15.5, -15.5)
	})
	t.Run("int panics on a non-finite value", func(t *testing.T) {
		for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			assert.Panics(t, func() {
				Cast[int](value)
			}, fmt.Sprintf("%v: ", value))
			assert.Panics(t, func() {
				Cast[namedInt8](value)
			}, fmt.Sprintf("%v: ", value))
		}
	})
	t.Run("float keeps a non-finite value", func(t *testing.T) {
		assert.True(t, math.IsNaN(float64(Cast[float64](math.NaN()))))
		assert.True(t, math.IsInf(float64(Cast[float32](math.Inf(1))), 1))
		assert.True(t, math.IsInf(float64(Cast[namedFloat64](math.Inf(-1))), -1))
	})
}

func BenchmarkCast_Int(b *testing.B) {
	for b.Loop() {
		Cast[int](1.51)
	}
}

func BenchmarkCast_Float64(b *testing.B) {
	for b.Loop() {
		Cast[float64](1.51)
	}
}

func TestInt(t *testing.T) {
	t.Run("integer converts directly", func(t *testing.T) {
		assert.Equal(t, Int(3), 3)
		assert.Equal(t, Int(namedInt8(-2)), -2)
		assert.Equal(t, Int(int64(1<<53+1)), 1<<53+1)
	})
	t.Run("float rounds half away from zero", func(t *testing.T) {
		assert.Equal(t, Int(2.5), 3)
		assert.Equal(t, Int(-2.5), -3)
		assert.Equal(t, Int(float32(1.4)), 1)
	})
}

func TestString(t *testing.T) {
	t.Run("int has no decimal point", func(t *testing.T) {
		assertString(t, 3, "3")
		assertString(t, -2, "-2")
		assertString(t, 8, "8")
	})
	t.Run("float keeps two decimals", func(t *testing.T) {
		assertString(t, 29.59, "29.59")
		assertString(t, 1.001, "1.00")
		assertString(t, 1.009, "1.01")
		assertString(t, -1.011, "-1.01")
	})
	t.Run("the format follows T, not the value", func(t *testing.T) {
		// a whole float still prints with decimals
		assertString(t, 8.0, "8.00")
		assertString(t, 8.1, "8.10")
		assertString(t, 0.0000, "0.00")
		assertString(t, math.Copysign(0, -1), "0.00")
		assertString(t, -0.004, "0.00")
		assertString(t, float32(-0.004), "0.00")
		assertString(t, -0.005, "-0.01")
		assertString(t, float32(1), "1.00")
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		assertString(t, namedInt(3), "3")
		assertString(t, namedInt8(-2), "-2")
		assertString(t, namedFloat64(29.59), "29.59")
		assertString(t, namedFloat32(1.00), "1.00")
	})
}

func TestIsIntType(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		assertIsInt[int](t, true)
		assertIsInt[int8](t, true)
		assertIsInt[int16](t, true)
		assertIsInt[int32](t, true)
		assertIsInt[int64](t, true)
	})
	t.Run("floats", func(t *testing.T) {
		assertIsInt[float32](t, false)
		assertIsInt[float64](t, false)
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		assertIsInt[namedInt](t, true)
		assertIsInt[namedInt8](t, true)
		assertIsInt[namedFloat32](t, false)
		assertIsInt[namedFloat64](t, false)
	})
}

func BenchmarkIsIntType_Int(b *testing.B) {
	for b.Loop() {
		isInt[int]()
	}
}

func BenchmarkIsIntType_Float64(b *testing.B) {
	for b.Loop() {
		isInt[float64]()
	}
}

func TestIsFloat32(t *testing.T) {
	t.Run("float32", func(t *testing.T) {
		assertIsFloat32[float32](t, true)
	})
	t.Run("anything else", func(t *testing.T) {
		assertIsFloat32[float64](t, false)
		assertIsFloat32[int](t, false)
		assertIsFloat32[int32](t, false) // same width, but an integer
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		assertIsFloat32[namedFloat32](t, true)
		assertIsFloat32[namedFloat64](t, false)
		assertIsFloat32[namedInt](t, false)
	})
}

func assertCast[T Number](t *testing.T, value float64, expected T) {
	t.Helper()

	AssertNumber(t, Cast[T](value), expected)
}

func assertString[T Number](t *testing.T, value T, expected string) {
	t.Helper()

	assert.Equal(t, String(value), expected)
}

func assertIsInt[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isInt[T](), expected)
}

func assertIsFloat32[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isFloat32[T](), expected)
}

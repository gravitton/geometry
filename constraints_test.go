package geom

import (
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
		assertIsIntType[int](t, true)
		assertIsIntType[int8](t, true)
		assertIsIntType[int16](t, true)
		assertIsIntType[int32](t, true)
		assertIsIntType[int64](t, true)
	})
	t.Run("floats", func(t *testing.T) {
		assertIsIntType[float32](t, false)
		assertIsIntType[float64](t, false)
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		assertIsIntType[namedInt](t, true)
		assertIsIntType[namedInt8](t, true)
		assertIsIntType[namedFloat32](t, false)
		assertIsIntType[namedFloat64](t, false)
	})
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

func assertIsIntType[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isIntType[T](), expected)
}

func assertIsFloat32[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isFloat32[T](), expected)
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

func BenchmarkIsIntType_Int(b *testing.B) {
	for b.Loop() {
		isIntType[int]()
	}
}

func BenchmarkIsIntType_Float64(b *testing.B) {
	for b.Loop() {
		isIntType[float64]()
	}
}

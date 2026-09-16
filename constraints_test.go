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

func TestIsIntType(t *testing.T) {
	testIsIntType[int](t, true)
	testIsIntType[int8](t, true)
	testIsIntType[int16](t, true)
	testIsIntType[int32](t, true)
	testIsIntType[int64](t, true)
	testIsIntType[float32](t, false)
	testIsIntType[float64](t, false)

	testIsIntType[namedInt](t, true)
	testIsIntType[namedInt8](t, true)
	testIsIntType[namedFloat32](t, false)
	testIsIntType[namedFloat64](t, false)
}

func TestCast(t *testing.T) {
	testCast[int](t, 1, 1)
	testCast[int](t, 1.2, 1)
	testCast[int](t, 1.6, 2)
	testCast[int](t, -0.3, 0)
	testCast[int](t, -0.51, -1)
	testCast[float64](t, 1.0, 1.0)
	testCast[float64](t, 1.6, 1.6)
	testCast[float64](t, -15.68, -15.68)

	testCast[namedInt](t, 1.6, 2)
	testCast[namedInt](t, -0.51, -1)
	testCast[namedInt8](t, 1.2, 1)
	testCast[namedFloat64](t, 1.6, 1.6)
	testCast[namedFloat32](t, -15.5, -15.5)
}

func TestToString(t *testing.T) {
	testToString(t, 3, "3")
	testToString(t, -2, "-2")
	testToString(t, 29.59, "29.59")
	testToString(t, 1.001, "1.00")
	testToString(t, 1.009, "1.01")
	testToString(t, -1.011, "-1.01")

	// the format follows T, not the value: a whole float still prints with decimals
	testToString(t, 8.0, "8.00")
	testToString(t, 8.1, "8.10")
	testToString(t, 0.0000, "0.00")
	testToString(t, float32(1), "1.00")
	testToString(t, 8, "8")

	testToString(t, namedInt(3), "3")
	testToString(t, namedInt8(-2), "-2")
	testToString(t, namedFloat64(29.59), "29.59")
	testToString(t, namedFloat32(1.00), "1.00")
}

func TestIsFloat32(t *testing.T) {
	testIsFloat32[float32](t, true)
	testIsFloat32[float64](t, false)
	testIsFloat32[int](t, false)
	testIsFloat32[int32](t, false) // same width, but an integer

	testIsFloat32[namedFloat32](t, true)
	testIsFloat32[namedFloat64](t, false)
	testIsFloat32[namedInt](t, false)
}

func testIsFloat32[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isFloat32[T](), expected)
}

func testIsIntType[T Number](t *testing.T, expected bool) {
	t.Helper()

	assert.Equal(t, isIntType[T](), expected)
}

func testCast[T Number](t *testing.T, value float64, expected T) {
	t.Helper()

	assert.Equal(t, Cast[T](value), expected)
}

func testToString[T Number](t *testing.T, value T, expected string) {
	t.Helper()

	assert.Equal(t, String(value), expected)
}

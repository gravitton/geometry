package geom

import (
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func BenchmarkIsIntValue_Int(b *testing.B) {
	for b.Loop() {
		isIntValue(23)
	}
}

func BenchmarkIsIntValue_Float64(b *testing.B) {
	for b.Loop() {
		isIntValue(0.5)
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

func TestIsIntValue(t *testing.T) {
	testIsIntValue(t, 1, true)
	testIsIntValue(t, int32(1), true)
	testIsIntValue(t, -2, true)
	testIsIntValue(t, 986, true)
	testIsIntValue(t, float32(1.000), true)
	testIsIntValue(t, -2.0, true)
	testIsIntValue(t, 189.2, false)
	testIsIntValue(t, float32(-9.3333), false)
	testIsIntValue(t, 1.000001, false)
	testIsIntValue(t, 23.0000000, true)
	testIsIntValue(t, math.NaN(), false)
	testIsIntValue(t, math.Inf(1), false)
}

func TestIsIntType(t *testing.T) {
	testIsIntType[int](t, true)
	testIsIntType[int8](t, true)
	testIsIntType[int16](t, true)
	testIsIntType[int32](t, true)
	testIsIntType[int64](t, true)
	testIsIntType[float32](t, false)
	testIsIntType[float64](t, false)
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
}

func TestToString(t *testing.T) {
	testToString(t, 3, "3")
	testToString(t, -2, "-2")
	testToString(t, 0.0000, "0")
	testToString(t, 1.00, "1")
	testToString(t, 29.59, "29.59")
	testToString(t, 1.001, "1.00")
	testToString(t, 1.009, "1.01")
	testToString(t, -1.011, "-1.01")
}

func testIsIntValue[T Number](t *testing.T, value T, expected bool) {
	t.Helper()

	assert.Equal(t, isIntValue(value), expected)
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

	assert.Equal(t, ToString(value), expected)
}

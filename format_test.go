package geom

import (
	"github.com/gravitton/assert"
	"math"
	"testing"
)

func TestIsInt(t *testing.T) {
	testIsInt(t, 1, true)
	testIsInt(t, -2, true)
	testIsInt(t, 986, true)
	testIsInt(t, 1.0, true)
	testIsInt(t, -2.0, true)
	testIsInt(t, 189.2, false)
	testIsInt(t, -9.3333, false)
	testIsInt(t, 1.000001, false)
	testIsInt(t, 23.0000000, true)
	testIsInt(t, math.NaN(), false)
	testIsInt(t, math.Inf(1), false)
}

func TestToString(t *testing.T) {
	testToString(t, 3, "+3")
	testToString(t, -2, "-2")
	testToString(t, 0.0000, "+0")
	testToString(t, 1.00, "+1")
	testToString(t, 29.59, "+29.59")
	testToString(t, 1.001, "+1.00")
	testToString(t, 1.009, "+1.01")
	testToString(t, -1.011, "-1.01")
}

func testIsInt[T Number](t *testing.T, value T, expected bool) {
	t.Helper()

	assert.Equal(t, IsInt(value), expected)
}

func testToString[T Number](t *testing.T, value T, expected string) {
	t.Helper()

	assert.Equal(t, ToString(value), expected)
}

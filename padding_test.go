package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestPadding_New(t *testing.T) {
	testPadding(t, Pad(2, 4, 3, 5), 2, 4, 3, 5)
	testPadding(t, PadU(20.0), 20.0, 20.0, 20.0, 20.0)
	testPadding(t, Pad(2.1, 4.0, 3.3, 5.4), 2.1, 4.0, 3.3, 5.4)
}

func TestPadding_Width(t *testing.T) {
	assert.Equal(t, Pad(2, 4, 3, 5).Width(), 9)
}

func TestPadding_Height(t *testing.T) {
	assert.Equal(t, Pad(2, 4, 3, 5).Height(), 5)
}

func TestPadding_XY(t *testing.T) {
	w1, h1 := Pad(2, 4, 3, 5).XY()
	assert.Equal(t, w1, 9)
	assert.Equal(t, h1, 5)
}

func testPadding[T Number](t *testing.T, p Padding[T], top, right, bottom, left T) {
	t.Helper()

	assert.EqualDelta(t, float64(p.Top), float64(top), Delta)
	assert.EqualDelta(t, float64(p.Right), float64(right), Delta)
	assert.EqualDelta(t, float64(p.Bottom), float64(bottom), Delta)
	assert.EqualDelta(t, float64(p.Left), float64(left), Delta)
}

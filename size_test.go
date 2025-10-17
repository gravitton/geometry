package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

var (
	sizeInt   = Size[int]{2, 3}
	sizeFloat = Size[float64]{1.2, 3.6}
)

func TestSize_New(t *testing.T) {
	assertSize(t, Sz(10, 16), 10, 16)
	assertSize(t, Sz(0.16, 204), 0.16, 204.0)
	assertSize(t, SzU(0.2), 0.2, 0.2)
}

func TestSize_Scale(t *testing.T) {
	assertSize(t, Sz(2, 3).Scale(2.5), 5, 8)
	assertSize(t, Sz(2, 3).ScaleXY(2, 3), 4, 9)
	assertSize(t, Sz(0.4, -0.25).Scale(2.5), 1.0, -0.625)
	assertSize(t, Sz(0.4, -0.25).ScaleXY(-1.5, 2), -0.6, -0.5)
}

func TestSize_Grow(t *testing.T) {
	assertSize(t, Sz(2, 3).Grow(2), 4, 5)
	assertSize(t, Sz(2, 3).GrowXY(2, 3), 4, 6)
	assertSize(t, Sz(0.4, 0.25).Grow(0.1), 0.5, 0.35)
	assertSize(t, Sz(0.4, 0.25).GrowXY(0.1, 0.2), 0.5, 0.45)
}

func TestSize_Shrink(t *testing.T) {
	assertSize(t, Sz(2, 3).Shrink(1), 1, 2)
	assertSize(t, Sz(2, 3).ShrinkXY(1, 2), 1, 1)
	assertSize(t, Sz(0.4, 0.25).Shrink(0.1), 0.3, 0.15)
	assertSize(t, Sz(0.4, 0.25).ShrinkXY(0.1, 0.2), 0.3, 0.05)
}

func TestSize_Area(t *testing.T) {
	assert.Equal(t, Sz(5, 3).Area(), 15)
	assert.EqualDelta(t, Sz(0.4, 0.25).Area(), 0.1, Delta)
}

func TestSize_Perimeter(t *testing.T) {
	assert.Equal(t, Sz(5, 3).Perimeter(), 16)
	assert.EqualDelta(t, Sz(0.4, 0.25).Perimeter(), 1.3, Delta)
}

func TestSize_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, Sz(5, 3).AspectRatio(), 1.666666, Delta)
	assert.EqualDelta(t, Sz(0.4, 0.25).AspectRatio(), 1.6, Delta)
}

func TestSize_Equal(t *testing.T) {
	assert.False(t, Sz(1, 2).Equal(Sz(3, -3)))
	assert.True(t, Sz(1, 2).Equal(Sz(1, 2)))
	assert.False(t, Sz(0.4, -0.25).Equal(Sz(100.1, -0.1)))
	assert.True(t, Sz(0.4, -0.25).Equal(Sz(0.4, -0.25)))
	assert.True(t, Sz(0.4, -0.25).Equal(Sz(0.4, -0.250001)))
}

func TestSize_IsZero(t *testing.T) {
	assert.False(t, Sz(1, 2).IsZero())
	assert.True(t, Sz(0, 0).IsZero())
	assert.True(t, Sz[float64](0.0, 0.000001).IsZero())
}

func TestSize_XY(t *testing.T) {
	w1, h1 := Sz(10, 16).XY()
	assert.Equal(t, w1, 10)
	assert.Equal(t, h1, 16)

	w2, h2 := Sz(0.4, -0.25).XY()
	assert.Equal(t, w2, 0.4)
	assert.Equal(t, h2, -0.25)
}

func TestSize_Int(t *testing.T) {
	assertSize(t, sizeInt.Int(), 2, 3)
	assertSize(t, sizeFloat.Int(), 1, 4)
}

func TestSize_Float(t *testing.T) {
	assertSize(t, sizeInt.Float(), 2.0, 3.0)
	assertSize(t, sizeFloat.Float(), 1.2, 3.6)
}

func TestSize_String(t *testing.T) {
	assert.Equal(t, Sz(10, 16).String(), "10x16")
	assert.Equal(t, Sz(100, -34.0000115).String(), "100x-34.00")
}

func TestSize_Marshall(t *testing.T) {
	assert.JSON(t, Sz(10, 16), `{"w":10,"h":16}`)
	assert.JSON(t, Sz(100, -34.0000115), `{"w":100.0,"h":-34.0000115}`)
}

func TestSize_Unmarshall(t *testing.T) {
	var s1 Size[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"w":10,"h":16}`), &s1))
	assertSize(t, s1, 10, 16)

	var s2 Size[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"w":10.1,"h":34.0000115}`), &s2))
	assertSize(t, s2, 10.1, 34.0000115)
}

func TestSize_Immutable(t *testing.T) {
	s := Sz(2, 3)

	s.Scale(2)
	s.ScaleXY(2, 3)
	s.Grow(1)
	s.GrowXY(1, 2)
	s.Shrink(1)
	s.ShrinkXY(1, 2)

	assertSize(t, s, 2, 3)
}

func assertSize[T Number](t *testing.T, s Size[T], w, h T, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(s.Width), float64(w), Delta, append(messages, "Width: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(s.Height), float64(h), Delta, append(messages, "Height: ")...) {
		ok = false
	}

	return ok
}

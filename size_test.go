package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestSize_New(t *testing.T) {
	testSize(t, S(10, 16), 10, 16)
	testSize(t, S[float64](0.16, 204), 0.16, 204.0)
}

func TestSize_Scale(t *testing.T) {
	testSize(t, S(2, 3).Scale(2.5), 5, 8)
	testSize(t, S(2, 3).ScaleXY(2, 3), 4, 9)
	testSize(t, S(0.4, -0.25).Scale(2.5), 1.0, -0.625)
	testSize(t, S(0.4, -0.25).ScaleXY(-1.5, 2), -0.6, -0.5)
}

func TestSize_Expand(t *testing.T) {
	testSize(t, S(2, 3).Expand(2), 4, 5)
	testSize(t, S(2, 3).ExpandXY(2, 3), 4, 6)
	testSize(t, S(0.4, 0.25).Expand(0.1), 0.5, 0.35)
	testSize(t, S(0.4, 0.25).ExpandXY(0.1, 0.2), 0.5, 0.45)
}

func TestSize_Shrunk(t *testing.T) {
	testSize(t, S(2, 3).Shrunk(1), 1, 2)
	testSize(t, S(2, 3).ShrunkXY(1, 2), 1, 1)
	testSize(t, S(0.4, 0.25).Shrunk(0.1), 0.3, 0.15)
	testSize(t, S(0.4, 0.25).ShrunkXY(0.1, 0.2), 0.3, 0.05)
}

func TestSize_Area(t *testing.T) {
	assert.Equal(t, S(5, 3).Area(), 15)
	assert.EqualDelta(t, S(0.4, 0.25).Area(), 0.1, Delta)
}

func TestSize_Perimeter(t *testing.T) {
	assert.Equal(t, S(5, 3).Perimeter(), 16)
	assert.EqualDelta(t, S(0.4, 0.25).Perimeter(), 1.3, Delta)
}

func TestSize_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, S(5, 3).AspectRatio(), 1.666666, Delta)
	assert.EqualDelta(t, S(0.4, 0.25).AspectRatio(), 1.6, Delta)
}

func TestSize_Equal(t *testing.T) {
	assert.False(t, S(1, 2).Equal(S(3, -3)))
	assert.True(t, S(1, 2).Equal(S(1, 2)))
	assert.False(t, S(0.4, -0.25).Equal(S(100.1, -0.1)))
	assert.True(t, S(0.4, -0.25).Equal(S(0.4, -0.25)))
	assert.True(t, S(0.4, -0.25).Equal(S(0.4, -0.250001)))
}

func TestSize_Zero(t *testing.T) {
	assert.False(t, S(1, 2).Zero())
	assert.True(t, S(0, 0).Zero())
	assert.True(t, S[float64](0.0, 0.000001).Zero())
}

func TestSize_XY(t *testing.T) {
	w1, h1 := S(10, 16).XY()
	assert.Equal(t, w1, 10)
	assert.Equal(t, h1, 16)

	w2, h2 := S(0.4, -0.25).XY()
	assert.Equal(t, w2, 0.4)
	assert.Equal(t, h2, -0.25)
}

func TestSize_String(t *testing.T) {
	assert.Equal(t, S(10, 16).String(), "10x16")
	assert.Equal(t, S(100, -34.0000115).String(), "100x-34.00")
}

func TestSize_Marshall(t *testing.T) {
	assert.JSON(t, S(10, 16), `{"w":10,"h":16}`)
	assert.JSON(t, S(100, -34.0000115), `{"w":100.0,"h":-34.0000115}`)
}

func TestSize_Unmarshall(t *testing.T) {
	var s1 Size[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"w":10,"h":16}`), &s1))
	testSize(t, s1, 10, 16)

	var s2 Size[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"w":10.1,"h":34.0000115}`), &s2))
	testSize(t, s2, 10.1, 34.0000115)
}

func TestSize_Immutable(t *testing.T) {
	s := S(2, 3)

	s.Scale(2)
	s.ScaleXY(2, 3)
	s.Expand(1)
	s.ExpandXY(1, 2)
	s.Shrunk(1)
	s.ShrunkXY(1, 2)

	testSize(t, s, 2, 3)
}

func testSize[T Number](t *testing.T, s Size[T], w, h T) {
	t.Helper()

	assert.EqualDelta(t, float64(s.Width), float64(w), Delta)
	assert.EqualDelta(t, float64(s.Height), float64(h), Delta)
}

package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestPadding_Constructor(t *testing.T) {
	AssertPadding(t, Pad(2, 4, 3, 5), Padding[int]{Top: 2, Right: 4, Bottom: 3, Left: 5})
	AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4), Padding[float64]{Top: 0.1, Right: 3.0, Bottom: 0.6, Left: 2.4})

	AssertPadding(t, PadU(20.0), Pad(20.0, 20.0, 20.0, 20.0))
	AssertPadding(t, PadXY(10.0, 20.0), Pad(10.0, 20.0, 10.0, 20.0))
}

func TestPadding_Width(t *testing.T) {
	assert.Equal(t, Pad(2, 4, 3, 5).Width(), 9)
	assert.Equal(t, Pad(0.1, 3.0, 0.6, 2.4).Width(), 5.4)
}

func TestPadding_Height(t *testing.T) {
	assert.Equal(t, Pad(2, 4, 3, 5).Height(), 5)
	assert.Equal(t, Pad(0.1, 3.0, 0.6, 2.4).Height(), 0.7)
}

func TestPadding_XY(t *testing.T) {
	w1, h1 := Pad(2, 4, 3, 5).XY()
	assert.Equal(t, w1, 9)
	assert.Equal(t, h1, 5)

	w2, h2 := Pad(0.1, 3.0, 0.6, 2.4).XY()
	assert.Equal(t, w2, 5.4)
	assert.Equal(t, h2, 0.7)
}

func TestPadding_Size(t *testing.T) {
	AssertSize(t, Pad(2, 4, 3, 5).Size(), Sz(9, 5))
	AssertSize(t, Pad(0.1, 3.0, 0.6, 2.4).Size(), Sz(5.4, 0.7))
}

func TestPadding_Int(t *testing.T) {
	AssertPadding(t, Pad(2, 4, 3, 5).Int(), Pad(2, 4, 3, 5))
	AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4).Int(), Pad(0, 3, 1, 2))
}

func TestPadding_Float(t *testing.T) {
	AssertPadding(t, Pad(2, 4, 3, 5).Float(), Pad(2.0, 4.0, 3.0, 5.0))
	AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4).Float(), Pad(0.1, 3.0, 0.6, 2.4))
}

func TestPadding_String(t *testing.T) {
	assert.Equal(t, Pad(2, 4, 3, 5).String(), "Pad(2;4;3;5)")
	assert.Equal(t, Pad(0.1, 3.0, 0.6, 2.4).String(), "Pad(0.10;3.00;0.60;2.40)")
}

func TestPadding_Marshall(t *testing.T) {
	assert.JSON(t, Pad(2, 4, 3, 5), `{"t":2,"r":4,"b":3,"l":5}`)
	assert.JSON(t, Pad(0.1, 3.0, 0.6, 2.4), `{"t":0.10,"r":3,"b":0.60,"l":2.40}`)
}

func TestPadding_Unmarshall(t *testing.T) {
	var p1 Padding[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"t":2,"r":4,"b":3,"l":5}`), &p1))
	AssertPadding(t, p1, Pad(2, 4, 3, 5))

	var p2 Padding[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"t":0.10,"r":3,"b":0.60,"l":2.40}`), &p2))
	AssertPadding(t, p2, Pad(0.1, 3.0, 0.6, 2.4))
}

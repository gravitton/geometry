package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

var (
	paddingInt   = Padding[int]{2, 4, 3, 5}
	paddingFloat = Padding[float64]{0.1, 3.0, 0.6, 2.4}
)

func TestPadding_New(t *testing.T) {
	AssertPadding(t, paddingInt, 2, 4, 3, 5)
	AssertPadding(t, paddingFloat, 0.1, 3.0, 0.6, 2.4)

	AssertPadding(t, PadU(20.0), 20.0, 20.0, 20.0, 20.0)
	AssertPadding(t, PadXY(10.0, 20.0), 10.0, 20.0, 10.0, 20.0)
}

func TestPadding_Width(t *testing.T) {
	assert.Equal(t, paddingInt.Width(), 9)
	assert.Equal(t, paddingFloat.Width(), 5.4)
}

func TestPadding_Height(t *testing.T) {
	assert.Equal(t, paddingInt.Height(), 5)
	assert.Equal(t, paddingFloat.Height(), 0.7)
}

func TestPadding_XY(t *testing.T) {
	w1, h1 := paddingInt.XY()
	assert.Equal(t, w1, 9)
	assert.Equal(t, h1, 5)

	w2, h2 := paddingFloat.XY()
	assert.Equal(t, w2, 5.4)
	assert.Equal(t, h2, 0.7)
}

func TestPadding_Size(t *testing.T) {
	AssertSize(t, paddingInt.Size(), 9, 5)
	AssertSize(t, paddingFloat.Size(), 5.4, 0.7)
}

func TestPadding_Int(t *testing.T) {
	AssertPadding(t, paddingInt.Int(), 2, 4, 3, 5)
	AssertPadding(t, paddingFloat.Int(), 0, 3, 1, 2)
}

func TestPadding_Float(t *testing.T) {
	AssertPadding(t, paddingInt.Float(), 2.0, 4.0, 3.0, 5.0)
	AssertPadding(t, paddingFloat.Float(), 0.1, 3.0, 0.6, 2.4)
}

func TestPadding_String(t *testing.T) {
	assert.Equal(t, paddingInt.String(), "Pad(2;4;3;5)")
	assert.Equal(t, paddingFloat.String(), "Pad(0.10;3;0.60;2.40)")
}

func TestPadding_Marshall(t *testing.T) {
	assert.JSON(t, paddingInt, `{"t":2,"r":4,"b":3,"l":5}`)
	assert.JSON(t, paddingFloat, `{"t":0.10,"r":3,"b":0.60,"l":2.40}`)
}

func TestPadding_Unmarshall(t *testing.T) {
	var p1 Padding[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"t":2,"r":4,"b":3,"l":5}`), &p1))
	AssertPadding(t, p1, 2, 4, 3, 5)

	var p2 Padding[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"t":0.10,"r":3,"b":0.60,"l":2.40}`), &p2))
	AssertPadding(t, p2, 0.1, 3.0, 0.6, 2.4)
}

package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

var (
	lineInt   = Line[int]{Point[int]{1, 2}, Point[int]{3, 5}}
	lineFloat = Line[float64]{Point[float64]{0.6, -0.25}, Point[float64]{1.2, 3.4}}
)

func TestLine_New(t *testing.T) {
	AssertLine(t, Ln(Pt(1, -1), Pt(2, 0)), 1, -1, 2, 0)
	AssertLine(t, Ln(Pt(0.5, -1.25), Pt(2.5, 3.75)), 0.5, -1.25, 2.5, 3.75)
}

func TestLine_Translate(t *testing.T) {
	AssertLine(t, lineInt.Translate(Vec(3, -2)), 4, 0, 6, 3)
	AssertLine(t, lineFloat.Translate(Vec(100.1, -0.1)), 100.7, -0.35, 101.3, 3.3)
}

func TestLine_MoveTo(t *testing.T) {
	AssertLine(t, lineInt.MoveTo(Pt(3, -2)), 3, -2, 5, 1)
	AssertLine(t, lineFloat.MoveTo(Pt(100.1, -0.1)), 100.1, -0.1, 100.7, 3.55)
}

func TestLine_Reversed(t *testing.T) {
	AssertLine(t, lineInt.Reversed(), 3, 5, 1, 2)
	AssertLine(t, lineFloat.Reversed(), 1.2, 3.4, 0.6, -0.25)
}

func TestLine_Midpoint(t *testing.T) {
	AssertPoint(t, lineInt.Midpoint(), 2, 4)
	AssertPoint(t, lineFloat.Midpoint(), 0.9, 1.575)
}

func TestLine_Direction(t *testing.T) {
	AssertVector(t, lineInt.Direction(), 2, 3)
	AssertVector(t, lineFloat.Direction(), 0.6, 3.65)
}

func TestLine_Length(t *testing.T) {
	assert.EqualDelta(t, lineInt.Length(), math.Sqrt(13), Delta)
	assert.EqualDelta(t, lineFloat.Length(), math.Sqrt(13.6825), Delta)
}

func TestLine_Bounds(t *testing.T) {
	AssertRect(t, lineInt.Bounds(), 2, 3, 2, 3)
	assert.Equal(t, lineInt.Start, lineInt.Bounds().Min())
	assert.Equal(t, lineInt.End, lineInt.Bounds().Max())

	AssertRect(t, lineFloat.Bounds(), 0.9, 1.575, 0.6, 3.65)
}

func TestLine_Equal(t *testing.T) {
	assert.False(t, lineInt.Equal(Ln(Pt(1, 2), Pt(3, 4))))
	assert.True(t, lineInt.Equal(lineInt))

	assert.False(t, lineFloat.Equal(Ln(Pt(0.5, -0.25), Pt(1.2, 3.4))))
	assert.True(t, lineFloat.Equal(lineFloat))
}

func TestLine_IsZero(t *testing.T) {
	assert.True(t, Line[int]{}.IsZero())
	assert.True(t, Ln(Pt(0, 0), Pt(0, 0)).IsZero())
	assert.False(t, Ln(Pt(1, 0), Pt(0, 0)).IsZero())
	assert.False(t, Ln(Pt(0, 0), Pt(0, 1)).IsZero())
	assert.False(t, lineInt.IsZero())

	assert.True(t, Line[float64]{}.IsZero())
	assert.False(t, lineFloat.IsZero())
}

func TestLine_Int(t *testing.T) {
	AssertLine(t, lineInt.Int(), 1, 2, 3, 5)
	AssertLine(t, lineFloat.Int(), 1, 0, 1, 3)
}

func TestLine_Float(t *testing.T) {
	AssertLine(t, lineInt.Float(), 1.0, 2.0, 3.0, 5.0)
	AssertLine(t, lineFloat.Float(), 0.6, -0.25, 1.2, 3.4)
}

func TestLine_String(t *testing.T) {
	assert.Equal(t, Ln(Pt(10, 16), Pt(1, 2)).String(), "L((10,16);(1,2))")
	assert.Equal(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)).String(), "L((100,-34.00);(0.20,0.40))")
}

func TestLine_Marshal(t *testing.T) {
	assert.JSON(t, Ln(Pt(10, 16), Pt(1, 2)), `{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`)
	assert.JSON(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)), `{"a":{"x":100.0,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`)
}

func TestLine_Unmarshal(t *testing.T) {
	var l1 Line[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`), &l1))
	assert.True(t, l1.Equal(Ln(Pt(10, 16), Pt(1, 2))))

	var l2 Line[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10.1,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`), &l2))
	assert.True(t, l2.Equal(Ln(Pt(10.1, -34.0000115), Pt(0.2, 0.4))))
}

func TestLine_Immutable(t *testing.T) {
	l := lineInt

	l.Translate(Vec(3, -2))
	l.MoveTo(Pt(4, 3))
	l.Reversed()

	assert.True(t, l.Equal(lineInt))
}

func AssertLine[T Number](t *testing.T, l Line[T], sx, sy, ex, ey T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, l.Start, sx, sy, append(messages, "Start.")...) {
		ok = false
	}
	if !AssertPoint(t, l.End, ex, ey, append(messages, "End.")...) {
		ok = false
	}

	return ok
}

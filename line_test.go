package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestLine_Constructor(t *testing.T) {
	AssertLine(t, Ln(Pt(1, -1), Pt(2, 0)), Line[int]{Start: Pt(1, -1), End: Pt(2, 0)})
	AssertLine(t, Ln(Pt(0.5, -1.25), Pt(2.5, 3.75)), Line[float64]{Start: Pt(0.5, -1.25), End: Pt(2.5, 3.75)})
}

func TestLine_Translate(t *testing.T) {
	AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Translate(Vec(3, -2)), Ln(Pt(4, 0), Pt(6, 3)))
	AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Translate(Vec(100.1, -0.1)), Ln(Pt(100.7, -0.35), Pt(101.3, 3.3)))
}

func TestLine_MoveTo(t *testing.T) {
	AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).MoveTo(Pt(3, -2)), Ln(Pt(3, -2), Pt(5, 1)))
	AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).MoveTo(Pt(100.1, -0.1)), Ln(Pt(100.1, -0.1), Pt(100.7, 3.55)))
}

func TestLine_Reverse(t *testing.T) {
	AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Reverse(), Ln(Pt(3, 5), Pt(1, 2)))
	AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Reverse(), Ln(Pt(1.2, 3.4), Pt(0.6, -0.25)))
}

func TestLine_Midpoint(t *testing.T) {
	AssertPoint(t, Ln(Pt(1, 2), Pt(3, 5)).Midpoint(), Pt(2, 4))
	AssertPoint(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Midpoint(), Pt(0.9, 1.575))
}

func TestLine_Vector(t *testing.T) {
	AssertVector(t, Ln(Pt(1, 2), Pt(3, 5)).Vector(), Vec(2, 3))
	AssertVector(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vector(), Vec(0.6, 3.65))
}

func TestLine_Length(t *testing.T) {
	assert.EqualDelta(t, Ln(Pt(1, 2), Pt(3, 5)).Length(), math.Sqrt(13), Delta)
	assert.EqualDelta(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Length(), math.Sqrt(13.6825), Delta)
}

func TestLine_Vertices(t *testing.T) {
	AssertVertices(t, Ln(Pt(1, 2), Pt(3, 5)).Vertices(), []Point[int]{{1, 2}, {3, 5}})
	AssertVertices(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vertices(), []Point[float64]{{0.6, -0.25}, {1.2, 3.4}})
}

func TestLine_Bounds(t *testing.T) {
	l := Ln(Pt(1, 2), Pt(3, 5))

	AssertRect(t, l.Bounds(), Rect(Pt(2, 3), Sz(2, 3)))
	assert.Equal(t, l.Start, l.Bounds().Min())
	assert.Equal(t, l.End, l.Bounds().Max())

	AssertRect(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Bounds(), Rect(Pt(0.9, 1.575), Sz(0.6, 3.65)))
}

func TestLine_Equal(t *testing.T) {
	assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 4))))
	assert.True(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 5))))

	assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.5, -0.25), Pt(1.2, 3.4))))
	assert.True(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.6, -0.25), Pt(1.2, 3.4))))
}

func TestLine_IsZero(t *testing.T) {
	assert.True(t, Line[int]{}.IsZero())
	assert.True(t, Ln(Pt(0, 0), Pt(0, 0)).IsZero())
	assert.False(t, Ln(Pt(1, 0), Pt(0, 0)).IsZero())
	assert.False(t, Ln(Pt(0, 0), Pt(0, 1)).IsZero())
	assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).IsZero())

	assert.True(t, Line[float64]{}.IsZero())
	assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).IsZero())
}

func TestLine_Int(t *testing.T) {
	AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Int(), Ln(Pt(1, 2), Pt(3, 5)))
	AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Int(), Ln(Pt(1, 0), Pt(1, 3)))
}

func TestLine_Float(t *testing.T) {
	AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Float(), Ln(Pt(1.0, 2.0), Pt(3.0, 5.0)))
	AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Float(), Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)))
}

func TestLine_String(t *testing.T) {
	assert.Equal(t, Ln(Pt(10, 16), Pt(1, 2)).String(), "L((10,16);(1,2))")
	assert.Equal(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)).String(), "L((100.00,-34.00);(0.20,0.40))")
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
	l := Ln(Pt(1, 2), Pt(3, 5))

	l.Translate(Vec(3, -2))
	l.MoveTo(Pt(4, 3))
	l.Reverse()

	assert.True(t, l.Equal(Ln(Pt(1, 2), Pt(3, 5))))
}

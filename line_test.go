package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestLine_New(t *testing.T) {
	testLine(t, L(P(1, -1), P(2, 0)), 1, -1, 2, 0)
	testLine(t, L(P(0.5, -1.25), P(2.5, 3.75)), 0.5, -1.25, 2.5, 3.75)
}

func TestLine_Translate(t *testing.T) {
	testLine(t, L(P(1, 2), P(3, 4)).Translate(V(3, -2)), 4, 0, 6, 2)
	testLine(t, L(P(0.4, -0.25), P(1.2, 3.4)).Translate(V(100.1, -0.1)), 100.5, -0.35, 101.3, 3.3)
}

func TestLine_MoveTo(t *testing.T) {
	testLine(t, L(P(1, 2), P(3, 4)).MoveTo(P(3, -2)), 3, -2, 5, 0)
	testLine(t, L(P(0.4, -0.25), P(1.2, 3.4)).MoveTo(P(100.1, -0.1)), 100.1, -0.1, 100.9, 3.55)
}

func TestLine_Reversed(t *testing.T) {
	testLine(t, L(P(1, 2), P(3, 4)).Reversed(), 3, 4, 1, 2)
	testLine(t, L(P(0.4, -0.25), P(1.2, 3.4)).Reversed(), 1.2, 3.4, 0.4, -0.25)
}

func TestLine_Midpoint(t *testing.T) {
	testPoint(t, L(P(1, 2), P(3, 4)).Midpoint(), 2, 3)
	testPoint(t, L(P(0.4, -0.25), P(1.2, 3.4)).Midpoint(), 0.8, 1.575)
}

func TestLine_Direction(t *testing.T) {
	testVector(t, L(P(1, 2), P(3, 4)).Direction(), 2, 2)
	testVector(t, L(P(0.4, -0.25), P(1.2, 3.4)).Direction(), 0.8, 3.65)
}

func TestLine_Length(t *testing.T) {
	assert.EqualDelta(t, L(P(1, 2), P(3, 5)).Length(), math.Sqrt(13), Delta)
	assert.EqualDelta(t, L(P(0.4, -0.25), P(1.2, 3.4)).Length(), math.Sqrt(13.9625), Delta)
}

func TestLine_Equal(t *testing.T) {
	assert.False(t, L(P(1, 2), P(3, 4)).Equal(L(P(1, 2), P(3, 5))))
	assert.True(t, L(P(1, 2), P(3, 4)).Equal(L(P(1, 2), P(3, 4))))

	assert.False(t, L(P(0.4, -0.25), P(1.2, 3.4)).Equal(L(P(0.5, -0.25), P(1.2, 3.4))))
	assert.True(t, L(P(0.4, -0.25), P(1.2, 3.4)).Equal(L(P(0.4, -0.25), P(1.2, 3.4))))
}

func TestLine_String(t *testing.T) {
	assert.Equal(t, L(P(10, 16), P(1, 2)).String(), "L((10,16);(1,2))")
	assert.Equal(t, L(P(100, -34.0000115), P(0.2, 0.4)).String(), "L((100,-34.00);(0.20,0.40))")
}

func TestLine_Marshal_Unmarshal(t *testing.T) {
	assert.JSON(t, L(P(10, 16), P(1, 2)), `{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`)
	assert.JSON(t, L(P(100, -34.0000115), P(0.2, 0.4)), `{"a":{"x":100.0,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`)
}

func TestLine_Unmarshal(t *testing.T) {
	var l1 Line[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`), &l1))
	assert.True(t, l1.Equal(L(P(10, 16), P(1, 2))))

	var l2 Line[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10.1,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`), &l2))
	assert.True(t, l2.Equal(L(P(10.1, -34.0000115), P(0.2, 0.4))))
}

func TestLine_Immutable(t *testing.T) {
	l := L(P(1, 2), P(3, 4))

	l.Translate(V(3, -2))
	l.MoveTo(P(4, 3))
	l.Reversed()

	assert.True(t, l.Equal(L(P(1, 2), P(3, 4))))
}

func testLine[T Number](t *testing.T, l Line[T], sx, sy, ex, ey T) {
	t.Helper()

	testPoint(t, l.Start, sx, sy)
	testPoint(t, l.End, ex, ey)
}

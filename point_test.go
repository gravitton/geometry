package geom

import (
	"encoding/json"
	"github.com/gravitton/assert"
	"testing"
)

func TestPoint_New(t *testing.T) {
	testPoint(t, Pt(10, 16), 10, 16)
	testPoint(t, Pt[float64](0.16, 204), 0.16, 204.0)

	testPoint(t, ZeroPoint[int](), 0, 0)
}

func TestPoint_Add(t *testing.T) {
	testPoint(t, Pt(1, 2).Add(Vec(3, -2)), 4, 0)
	testPoint(t, Pt(0.4, -0.25).Add(Vec(100.1, -0.1)), 100.5, -0.35)
}

func TestPoint_Sub(t *testing.T) {
	testVector(t, Pt(1, 2).Sub(Pt(3, -3)), -2, 5)
	testVector(t, Pt(0.4, -0.25).Sub(Pt(100.1, -0.1)), -99.7, -0.15)
}

func TestPoint_Midpoint(t *testing.T) {
	testPoint(t, Pt(1, 2).Midpoint(Pt(3, -3)), 2, 0)
	testPoint(t, Pt(0.4, -0.25).Midpoint(Pt(100.1, -0.1)), 50.25, -0.175)
}

func TestPoint_Lerp(t *testing.T) {
	testPoint(t, Pt(1, 2).Lerp(Pt(3, -3), 0.3), 1, 0)
	testPoint(t, Pt(0.4, -0.25).Lerp(Pt(100.1, -0.1), 0.1), 10.37, -0.235)
}

func TestPoint_Equal(t *testing.T) {
	assert.False(t, Pt(1, 2).Equal(Pt(3, -3)))
	assert.True(t, Pt(1, 2).Equal(Pt(1, 2)))

	assert.False(t, Pt(0.4, -0.25).Equal(Pt(100.1, -0.1)))
	assert.True(t, Pt(0.4, -0.25).Equal(Pt(0.4, -0.25)))
	assert.True(t, Pt(0.4, -0.25).Equal(Pt(0.4, -0.250001)))
}

func TestPoint_Zero(t *testing.T) {
	assert.False(t, Pt(1, 2).Zero())
	assert.True(t, Pt(0, -0).Zero())
	assert.True(t, ZeroPoint[int]().Zero())

	assert.False(t, Pt(0.4, -0.25).Zero())
	assert.True(t, Pt(0.0, -0.0).Zero())
	assert.True(t, Pt(0.0, 0.000001).Zero())
	assert.True(t, ZeroPoint[float64]().Zero())
}

func TestPoint_String(t *testing.T) {
	assert.Equal(t, Pt(10, 16).String(), "(+10,+16)")
	assert.Equal(t, Pt(100, -34.0000115).String(), "(+100,-34.00)")
}

func TestPoint_Marshall(t *testing.T) {
	assert.JSON(t, Pt(10, 16), `{"x":10,"y":16}`)
	assert.JSON(t, Pt(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)
}

func TestPoint_Unmarshall(t *testing.T) {
	var p1 Point[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &p1))
	testPoint(t, p1, 10, 16)

	var p2 Point[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &p2))
	testPoint(t, p2, 10.1, -34.0000115)
}

func TestPoint_Immutable(t *testing.T) {
	p1 := Pt(1, 2)
	p2 := Pt(3, -3)

	p1.Add(Vec(3, -2))
	p1.Sub(p2)
	p1.Midpoint(p2)
	p1.Lerp(p2, 0.1)

	testPoint(t, p1, 1, 2)
	testPoint(t, p2, 3, -3)
}

func testPoint[T Number](t *testing.T, p Point[T], x, y T) {
	t.Helper()

	assert.EqualDelta(t, float64(p.X), float64(x), Delta)
	assert.EqualDelta(t, float64(p.Y), float64(y), Delta)
}

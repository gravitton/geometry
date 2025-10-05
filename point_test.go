package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestPoint_New(t *testing.T) {
	testPoint(t, Pt(10, 16), 10, 16)
	testPoint(t, Pt[float64](0.16, 204), 0.16, 204.0)

	testPoint(t, ZeroPoint[int](), 0, 0)
}

func TestPoint_Add(t *testing.T) {
	testPoint(t, Pt(1, 2).Add(Vec(3, -2)), 4, 0)
	testPoint(t, Pt(1, 2).AddXY(3, -2), 4, 0)
	testPoint(t, Pt(0.4, -0.25).Add(Vec(100.1, -0.1)), 100.5, -0.35)
	testPoint(t, Pt(0.4, -0.25).AddXY(100.1, -0.1), 100.5, -0.35)
}

func TestPoint_Subtract(t *testing.T) {
	testVector(t, Pt(1, 2).Subtract(Pt(3, -3)), -2, 5)
	testVector(t, Pt(0.4, -0.25).Subtract(Pt(100.1, -0.1)), -99.7, -0.15)
}

func TestPoint_Multiply(t *testing.T) {
	testPoint(t, Pt(1, 2).Multiply(3), 3, 6)
	testPoint(t, Pt(1, 2).MultiplyXY(3, 4), 3, 8)
	testPoint(t, Pt(0.4, -0.25).Multiply(-1.5), -0.6, 0.375)
	testPoint(t, Pt(0.4, -0.25).MultiplyXY(-1.5, 2), -0.6, -0.5)
}

func TestPoint_Divide(t *testing.T) {
	testPoint(t, Pt(5, 10).Divide(2), 3, 5)
	testPoint(t, Pt(5, 10).DivideXY(3, 2), 2, 5)
	testPoint(t, Pt(0.4, -0.25).Divide(-2), -0.2, 0.125)
	testPoint(t, Pt(0.4, -0.25).DivideXY(-4, 0.5), -0.1, -0.5)
}

func TestPoint_Midpoint(t *testing.T) {
	testPoint(t, Pt(1, 2).Midpoint(Pt(3, -3)), 2, -1)
	testPoint(t, Pt(0.4, -0.25).Midpoint(Pt(100.1, -0.1)), 50.25, -0.175)
}

func TestPoint_Lerp(t *testing.T) {
	testPoint(t, Pt(1, 2).Lerp(Pt(3, -3), 0.3), 2, 1)
	testPoint(t, Pt(0.4, -0.25).Lerp(Pt(100.1, -0.1), 0.1), 10.37, -0.235)
}

func TestPoint_DistanceTo(t *testing.T) {
	assert.EqualDelta(t, Pt(1, 2).DistanceTo(Pt(2, 3)), math.Sqrt(2), Delta)
	assert.EqualDelta(t, Pt(0.4, -0.25).DistanceTo(Pt(0.5, -0.35)), math.Sqrt(0.02), Delta)
}

func TestPoint_DistanceSquaredTo(t *testing.T) {
	assert.Equal(t, Pt(1, 2).DistanceSquaredTo(Pt(2, 3)), 2)
	assert.EqualDelta(t, Pt(0.4, -0.25).DistanceSquaredTo(Pt(0.5, -0.35)), 0.02, Delta)
}

func TestPoint_AngleTo(t *testing.T) {
	assert.EqualDelta(t, Pt(2, 2).AngleTo(Pt(3, 2)), ToRadians(0), Delta)
	assert.EqualDelta(t, Pt(0.4, -0.25).AngleTo(Pt(0.5, -0.35)), ToRadians(-45), Delta)
}

func TestPoint_Equal(t *testing.T) {
	assert.False(t, Pt(1, 2).Equal(Pt(3, -3)))
	assert.True(t, Pt(1, 2).Equal(Pt(1, 2)))

	assert.False(t, Pt(0.4, -0.25).Equal(Pt(100.1, -0.1)))
	assert.True(t, Pt(0.4, -0.25).Equal(Pt(0.4, -0.25)))
	assert.True(t, Pt(0.4, -0.25).Equal(Pt(0.4, -0.250001)))
}

func TestPoint_IsZero(t *testing.T) {
	assert.False(t, Pt(1, 2).IsZero())
	assert.True(t, Pt(0, -0).IsZero())
	assert.True(t, ZeroPoint[int]().IsZero())

	assert.False(t, Pt(0.4, -0.25).IsZero())
	assert.True(t, Pt(0.0, 0.0).IsZero())
	assert.True(t, Pt(0.0, 0.000001).IsZero())
	assert.True(t, ZeroPoint[float64]().IsZero())
}

func TestPoint_XY(t *testing.T) {
	x1, y1 := Pt(10, 16).XY()
	assert.Equal(t, x1, 10)
	assert.Equal(t, y1, 16)

	x2, y2 := Pt(0.4, -0.25).XY()
	assert.Equal(t, x2, x2)
	assert.Equal(t, y2, y2)
}

func TestPoint_String(t *testing.T) {
	assert.Equal(t, Pt(10, 16).String(), "(10,16)")
	assert.Equal(t, Pt(100, -34.0000115).String(), "(100,-34.00)")
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
	p1.AddXY(2, 3)
	p1.Subtract(p2)
	p1.Multiply(2)
	p1.MultiplyXY(3, 4)
	p1.Divide(5)
	p1.DivideXY(10, 100)
	p1.Lerp(p2, 0.1)

	testPoint(t, p1, 1, 2)
	testPoint(t, p2, 3, -3)
}

func testPoint[T Number](t *testing.T, p Point[T], x, y T) {
	t.Helper()

	assert.EqualDelta(t, float64(p.X), float64(x), Delta)
	assert.EqualDelta(t, float64(p.Y), float64(y), Delta)
}

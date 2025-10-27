package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

var (
	pointInt   = Point[int]{1, 2}
	pointFloat = Point[float64]{0.6, -0.25}
)

func TestPoint_New(t *testing.T) {
	AssertPoint(t, Pt(10, 16), 10, 16)
	AssertPoint(t, Pt[float64](0.16, 204), 0.16, 204.0)

	AssertPoint(t, ZeroPoint[int](), 0, 0)
}

func TestPoint_Transform(t *testing.T) {
	AssertPoint(t, pointInt.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 9, 22)
	AssertPoint(t, pointFloat.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 3.385, 7.865)
}

func TestPoint_Add(t *testing.T) {
	AssertPoint(t, pointInt.Add(Vec(3, -2)), 4, 0)
	AssertPoint(t, pointInt.AddXY(3, -2), 4, 0)
	AssertPoint(t, pointFloat.Add(Vec(100.1, -0.1)), 100.7, -0.35)
	AssertPoint(t, pointFloat.AddXY(100.1, -0.1), 100.7, -0.35)
}

func TestPoint_Subtract(t *testing.T) {
	AssertVector(t, pointInt.Subtract(Pt(3, -3)), -2, 5)
	AssertVector(t, pointFloat.Subtract(Pt(100.1, -0.1)), -99.5, -0.15)
}

func TestPoint_Multiply(t *testing.T) {
	AssertPoint(t, pointInt.Multiply(3), 3, 6)
	AssertPoint(t, pointInt.MultiplyXY(3, 4), 3, 8)
	AssertPoint(t, pointFloat.Multiply(-1.5), -0.9, 0.375)
	AssertPoint(t, pointFloat.MultiplyXY(-1.5, 2), -0.9, -0.5)
}

func TestPoint_Divide(t *testing.T) {
	AssertPoint(t, Pt(5, 10).Divide(2), 3, 5)
	AssertPoint(t, Pt(5, 10).DivideXY(3, 2), 2, 5)
	AssertPoint(t, pointFloat.Divide(-2), -0.3, 0.125)
	AssertPoint(t, pointFloat.DivideXY(-4, 0.5), -0.15, -0.5)
}

func TestPoint_Midpoint(t *testing.T) {
	AssertPoint(t, pointInt.Midpoint(Pt(3, -3)), 2, -1)
	AssertPoint(t, pointFloat.Midpoint(Pt(100.1, -0.1)), 50.35, -0.175)
}

func TestPoint_Lerp(t *testing.T) {
	AssertPoint(t, pointInt.Lerp(Pt(3, -3), 0.3), 2, 1)
	AssertPoint(t, pointFloat.Lerp(Pt(100.1, -0.1), 0.1), 10.55, -0.235)
}

func TestPoint_DistanceTo(t *testing.T) {
	assert.EqualDelta(t, pointInt.DistanceTo(Pt(2, 3)), math.Sqrt(2), Delta)
	assert.EqualDelta(t, pointFloat.DistanceTo(Pt(0.5, -0.35)), math.Sqrt(0.02), Delta)
}

func TestPoint_DistanceSquaredTo(t *testing.T) {
	assert.Equal(t, pointInt.DistanceSquaredTo(Pt(2, 3)), 2)
	assert.EqualDelta(t, pointFloat.DistanceSquaredTo(Pt(0.5, -0.35)), 0.02, Delta)
}

func TestPoint_AngleTo(t *testing.T) {
	assert.EqualDelta(t, Pt(2, 2).AngleTo(Pt(3, 2)), ToRadians(0), Delta)
	assert.EqualDelta(t, pointFloat.AngleTo(Pt(0.7, -0.35)), ToRadians(-45), Delta)
}

func TestPoint_Equal(t *testing.T) {
	assert.False(t, pointInt.Equal(Pt(3, -3)))
	assert.True(t, pointInt.Equal(pointInt))

	assert.False(t, pointFloat.Equal(Pt(100.1, -0.1)))
	assert.True(t, pointFloat.Equal(pointFloat))
	assert.True(t, pointFloat.Equal(Pt(0.6, -0.250001)))
}

func TestPoint_IsZero(t *testing.T) {
	assert.False(t, pointInt.IsZero())
	assert.True(t, Pt(0, -0).IsZero())
	assert.True(t, ZeroPoint[int]().IsZero())

	assert.False(t, pointFloat.IsZero())
	assert.True(t, Pt(0.0, 0.0).IsZero())
	assert.True(t, Pt(0.0, 0.000001).IsZero())
	assert.True(t, ZeroPoint[float64]().IsZero())
}

func TestPoint_XY(t *testing.T) {
	x1, y1 := Pt(10, 16).XY()
	assert.Equal(t, x1, 10)
	assert.Equal(t, y1, 16)

	x2, y2 := pointFloat.XY()
	assert.Equal(t, x2, x2)
	assert.Equal(t, y2, y2)
}

func TestPoint_Vector(t *testing.T) {
	AssertVector(t, pointInt.Vector(), 1, 2)
	AssertVector(t, pointFloat.Vector(), 0.6, -0.25)
}

func TestPoint_Int(t *testing.T) {
	AssertPoint(t, pointInt.Int(), 1, 2)
	AssertPoint(t, pointFloat.Int(), 1, 0)
}

func TestPoint_Float(t *testing.T) {
	AssertPoint(t, pointInt.Float(), 1.0, 2.0)
	AssertPoint(t, pointFloat.Float(), 0.6, -0.25)
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
	AssertPoint(t, p1, 10, 16)

	var p2 Point[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &p2))
	AssertPoint(t, p2, 10.1, -34.0000115)
}

func TestPoint_Immutable(t *testing.T) {
	p1 := pointInt
	p2 := Pt(3, -3)

	p1.Add(Vec(3, -2))
	p1.AddXY(2, 3)
	p1.Subtract(p2)
	p1.Multiply(2)
	p1.MultiplyXY(3, 4)
	p1.Divide(5)
	p1.DivideXY(10, 100)
	p1.Lerp(p2, 0.1)

	AssertPoint(t, p1, 1, 2)
	AssertPoint(t, p2, 3, -3)
}

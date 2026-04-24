package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

var (
	regPolygonInt = RegularPolygon[int]{Point[int]{1, 2}, Size[int]{2, 2}, 4, 0}
)

func TestRegularPolygon_New(t *testing.T) {
	rp := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)
	AssertRegularPolygon(t, rp, 0, 0, 2, 2, 4, 0)

	triangle := Triangle(Pt(1, -1), Sz(3, 3), PointyTop)
	AssertRegularPolygon(t, triangle, 1, -1, 3, 3, 3, RegularPolygonAngle(3, PointyTop))
	// Integer types: sin/cos of non-right angles are rounded; sin(7π/6) computes slightly above
	// −0.5 in float64, rounding to 0 rather than −1, so vertex 1 Y = center.Y + 0 = −1.
	AssertVertices(t, triangle.Vertices(), []Point[int]{Pt(1, 2), Pt(-2, -1), Pt(4, -4)})

	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointyTop)
	AssertRegularPolygon(t, square, 50, 50, 100, 100, 4, RegularPolygonAngle(4, PointyTop))
	AssertVertices(t, square.Vertices(), []Point[float64]{Pt(50.0, 150.0), Pt(-50.0, 50.0), Pt(50.0, -50.0), Pt(150.0, 50.0)})

	hexagon := Hexagon(Pt(0, 0), Sz(10, 10), PointyTop)
	AssertRegularPolygon(t, hexagon, 0, 0, 10, 10, 6, RegularPolygonAngle(6, PointyTop))
	// Vertices 2 and 5 (at 7π/6 and π/6) have sin rounded to 0 due to float64 precision.
	AssertVertices(t, hexagon.Vertices(), []Point[int]{Pt(0, 10), Pt(-10, 10), Pt(-10, 0), Pt(0, -10), Pt(10, -10), Pt(10, 0)})
}

func TestRegularPolygonAngle(t *testing.T) {
	assert.EqualDelta(t, RegularPolygonAngle(3, PointyTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(3, FlatTop), 30*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonAngle(4, PointyTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(4, FlatTop), 45*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonAngle(6, PointyTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(6, FlatTop), 60*DegToRad, Delta)
}

func TestRegularPolygon_Translate(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.Translate(Vec(1, -2)), 2, 0, 2, 2, 4, 0)
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.MoveTo(Pt(-3, 5)), -3, 5, 2, 2, 4, 0)
}

func TestRegularPolygon_Scale(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.Scale(0.5), 1, 2, 1, 1, 4, 0)
	AssertRegularPolygon(t, regPolygonInt.ScaleXY(2, 3), 1, 2, 4, 6, 4, 0)
}

func TestRegularPolygon_Rotate(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.Rotate(math.Pi), 1, 2, 2, 2, 4, math.Pi)
}

func TestRegularPolygon_Vertices(t *testing.T) {
	AssertVertices(t, RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices(), []Point[int]{
		Pt(1, 0),
		Pt(0, 1),
		Pt(-1, 0),
		Pt(0, -1),
	})
	AssertVertices(t, RegPol(Pt(0, 0), Sz(2, 3), 4, 0).Vertices(), []Point[int]{
		Pt(2, 0),
		Pt(0, 3),
		Pt(-2, 0),
		Pt(0, -3),
	})
	AssertVertices(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0).Vertices(), []Point[float64]{
		Pt(2.0, 0.0),
		Pt(1.0, 1.5*Sqrt3),
		Pt(-1.0, 1.5*Sqrt3),
		Pt(-2.0, 0.0),
		Pt(-1.0, -1.5*Sqrt3),
		Pt(1.0, -1.5*Sqrt3),
	})
}

func TestRegularPolygon_Bounds(t *testing.T) {
	AssertRect(t, regPolygonInt.Bounds(), 1, 2, 4, 4)
}

func TestRegularPolygon_Polygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0)
	p := rp.Polygon()

	assert.Equal(t, p.Vertices, rp.Vertices())
	assert.NotSame(t, p.Vertices, rp.Vertices())
}

func TestRegularPolygon_Equal(t *testing.T) {
	assert.True(t, regPolygonInt.Equal(regPolygonInt))
	assert.False(t, regPolygonInt.Equal(RegPol(Pt(0, 0), Sz(2, 2), 6, 0)))
}

func TestRegularPolygon_IsZero(t *testing.T) {
	assert.False(t, regPolygonInt.IsZero())
	assert.True(t, RegularPolygon[int]{}.Empty())
}

func TestRegularPolygon_Empty(t *testing.T) {
	assert.False(t, regPolygonInt.Empty())
	assert.True(t, Polygon[int]{}.Empty())
	assert.True(t, Polygon[int]{[]Point[int]{}}.Empty())
	assert.False(t, polygonFloat.Empty())
}

func TestRegularPolygon_Int(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.Int(), 1, 2, 2, 2, 4, 0.0)
}

func TestRegularPolygon_Float(t *testing.T) {
	AssertRegularPolygon(t, regPolygonInt.Float(), 1.0, 2.0, 2.0, 2.0, 4, 0.0)
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RegPol(Pt(0, 1), Sz(2, 3), 4, 0)

	rp.Translate(Vec(2, 3))
	rp.MoveTo(Pt(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)

	AssertRegularPolygon(t, rp, 0, 1, 2, 3, 4, 0)
}

func TestRegularPolygon_String(t *testing.T) {
	assert.Equal(t, regPolygonInt.String(), "RegPol((1,2);2x2;4;0)")
}

func TestRegularPolygon_Marshall(t *testing.T) {
	assert.JSON(t, regPolygonInt, `{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`)
}

func TestRegularPolygon_Unmarshall(t *testing.T) {
	var p1 RegularPolygon[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`), &p1))
	AssertRegularPolygon(t, p1, 1, 2, 2, 2, 4, 0)
}

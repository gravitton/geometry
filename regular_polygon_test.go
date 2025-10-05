package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestRegularPolygon_New(t *testing.T) {
	rp := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)
	testRegularPolygon(t, rp, 0, 0, 2, 2, 4, 0)

	triangle := Triangle(Pt(1, -1), Sz(3, 3), PointTop)
	testRegularPolygon(t, triangle, 1, -1, 3, 3, 3, RegularPolygonAngle(3, PointTop))
	//testVertices(t, triangle.Vertices(), []Point[float64]{{0, 0}})

	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointTop)
	testRegularPolygon(t, square, 50, 50, 100, 100, 4, RegularPolygonAngle(4, PointTop))
	//testVertices(t, square.Vertices(), []Point[float64]{{0, 0}})

	hexagon := Hexagon(Pt(0, 0), Sz(10, 10), PointTop)
	testRegularPolygon(t, hexagon, 0, 0, 10, 10, 6, RegularPolygonAngle(6, PointTop))
	//testVertices(t, hexagon.Vertices(), []Point[float64]{{0, 0}})
}

func TestRegularPolygonAngle(t *testing.T) {
	assert.EqualDelta(t, RegularPolygonAngle(3, PointTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(3, FlatTop), 30*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonAngle(4, PointTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(4, FlatTop), 45*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonAngle(6, PointTop), 90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonAngle(6, FlatTop), 60*DegToRad, Delta)
}

func TestRegularPolygon_Translate(t *testing.T) {
	testRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Translate(Vec(1, -2)), 2, 0, 2, 2, 4, 0)
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	testRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).MoveTo(Pt(-3, 5)), -3, 5, 2, 2, 4, 0)
}

func TestRegularPolygon_Scale(t *testing.T) {
	testRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 3), 4, 0).Scale(0.5), 1, 2, 1, 2, 4, 0)
	testRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 3), 4, 0).ScaleXY(2, 3), 1, 2, 4, 9, 4, 0)
}

func TestRegularPolygon_Vertices(t *testing.T) {
	testVertices(t, RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices(), []Point[int]{
		Pt(1, 0),
		Pt(0, 1),
		Pt(-1, 0),
		Pt(0, -1),
	})
	testVertices(t, RegPol(Pt(0, 0), Sz(2, 3), 4, 0).Vertices(), []Point[int]{
		Pt(2, 0),
		Pt(0, 3),
		Pt(-2, 0),
		Pt(0, -3),
	})

	testVertices(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0).Vertices(), []Point[float64]{
		Pt(2.0, 0.0),
		Pt(1.0, 1.5*Sqrt3),
		Pt(-1.0, 1.5*Sqrt3),
		Pt(-2.0, 0.0),
		Pt(-1.0, -1.5*Sqrt3),
		Pt(1.0, -1.5*Sqrt3),
	})
}

func TestRegularPolygon_ToPolygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0)
	p := rp.ToPolygon()

	assert.Equal(t, p.Vertices, rp.Vertices())
	assert.NotSame(t, p.Vertices, rp.Vertices())
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RegPol(Pt(0, 1), Sz(2, 3), 4, 0)

	rp.Translate(Vec(2, 3))
	rp.MoveTo(Pt(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)

	testRegularPolygon(t, rp, 0, 1, 2, 3, 4, 0)
}

func testRegularPolygon[T Number](t *testing.T, p RegularPolygon[T], x, y, w, h T, n int, angle float64) {
	t.Helper()

	testPoint(t, p.Center, x, y)
	testSize(t, p.Size, w, h)
	assert.Equal(t, p.N, n)
	assert.EqualDelta(t, p.Angle, angle, Delta)
}

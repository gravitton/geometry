package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestRegularPolygon_New(t *testing.T) {
	rp := RP(P(0, 0), S(2, 2), 4)
	testRegularPolygon(t, rp, 0, 0, 2, 2, 4)

	triangle := Triangle(P(1, -1), S(3, 3))
	testRegularPolygon(t, triangle, 1, -1, 3, 3, 3)

	square := Square(P(-1, 2), S(1, 2))
	testRegularPolygon(t, square, -1, 2, 1, 2, 4)

	hexagon := Hexagon(P(0, 0), S(2, 3))
	testRegularPolygon(t, hexagon, 0, 0, 2, 3, 6)
}

func TestRegularPolygon_Translate(t *testing.T) {
	testRegularPolygon(t, RP(P(1, 2), S(2, 2), 4).Translate(V(1, -2)), 2, 0, 2, 2, 4)
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	testRegularPolygon(t, RP(P(1, 2), S(2, 2), 4).MoveTo(P(-3, 5)), -3, 5, 2, 2, 4)
}

func TestRegularPolygon_Scale(t *testing.T) {
	testRegularPolygon(t, RP(P(1, 2), S(2, 3), 4).Scale(0.5), 1, 2, 1, 2, 4)
	testRegularPolygon(t, RP(P(1, 2), S(2, 3), 4).ScaleXY(2, 3), 1, 2, 4, 9, 4)
}

func TestRegularPolygon_Vertices(t *testing.T) {
	testVertices(t, RP(P(0, 0), S(1, 1), 4).Vertices(), []Point[int]{
		P(1, 0),
		P(0, 1),
		P(-1, 0),
		P(0, -1),
	})
	testVertices(t, RP(P(0, 0), S(2, 3), 4).Vertices(), []Point[int]{
		P(2, 0),
		P(0, 3),
		P(-2, 0),
		P(0, -3),
	})

	testVertices(t, RP(P(0.0, 0.0), S(2.0, 3.0), 6).Vertices(), []Point[float64]{
		P(2.0, 0.0),
		P(1.0, 1.5*Sqrt3),
		P(-1.0, 1.5*Sqrt3),
		P(-2.0, 0.0),
		P(-1.0, -1.5*Sqrt3),
		P(1.0, -1.5*Sqrt3),
	})
}

func TestRegularPolygon_ToPolygon(t *testing.T) {
	rp := RP(P(0.0, 0.0), S(1.0, 1.0), 5)
	p := rp.ToPolygon()

	assert.Equal(t, p.Vertices, rp.Vertices())
	assert.NotSame(t, p.Vertices, rp.Vertices())
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RP(P(0, 1), S(2, 3), 4)

	rp.Translate(V(2, 3))
	rp.MoveTo(P(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)

	testRegularPolygon(t, rp, 0, 1, 2, 3, 4)
}

func testRegularPolygon[T Number](t *testing.T, p RegularPolygon[T], x, y, w, h T, n int) {
	t.Helper()

	testPoint(t, p.Center, x, y)
	testSize(t, p.Size, w, h)
	assert.Equal(t, p.N, n)
}

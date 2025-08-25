package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestPolygon_Center(t *testing.T) {
	p := Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0), P(2, 2), P(0, 2)}}
	testPoint(t, p.Center(), 1, 1)

	pf := Polygon[float64]{Vertices: []Point[float64]{P(0.0, 0.0), P(2.0, 0.0), P(2.0, 1.0)}}
	cx, cy := pf.Center().XY()
	assert.EqualDelta(t, cx, 4.0/3.0, Delta)
	assert.EqualDelta(t, cy, 1.0/3.0, Delta)
}

func TestPolygon_Translate_MoveTo(t *testing.T) {
	p := Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0)}}
	p2 := p.Translate(V(1, -1))
	testPoint(t, p2.Vertices[0], 1, -1)
	testPoint(t, p2.Vertices[1], 3, -1)

	p3 := p.MoveTo(P(10, 10))
	testPoint(t, p3.Center(), 10, 10)
}

func TestPolygon_Scale(t *testing.T) {
	p := Polygon[float64]{Vertices: []Point[float64]{P(1.0, 0.0), P(0.0, 1.0), P(-1.0, 0.0), P(0.0, -1.0)}}
	center := p.Center()
	assert.True(t, center.Equal(P(0.0, 0.0)))

	p2 := p.Scale(2)
	// Every point should double its distance from center (0,0).
	assert.EqualDelta(t, p2.Vertices[0].X, 2.0, Delta)
	assert.EqualDelta(t, p2.Vertices[0].Y, 0.0, Delta)
	assert.EqualDelta(t, p2.Vertices[1].X, 0.0, Delta)
	assert.EqualDelta(t, p2.Vertices[1].Y, 2.0, Delta)

	p3 := p.ScaleXY(2, 0.5)
	assert.EqualDelta(t, p3.Vertices[0].X, 2.0, Delta)
	assert.EqualDelta(t, p3.Vertices[0].Y, 0.0, Delta)
	assert.EqualDelta(t, p3.Vertices[1].X, 0.0, Delta)
	assert.EqualDelta(t, p3.Vertices[1].Y, 0.5, Delta)
}

func TestPolygon_Immutable(t *testing.T) {
	p := Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0)}}

	p.Translate(V(1, -1))
	p.MoveTo(P(10, 10))
	p.Scale(2)
	p.ScaleXY(2, 3)

	// unchanged
	testPoint(t, p.Vertices[0], 0, 0)
	testPoint(t, p.Vertices[1], 2, 0)
}

func testPolygon[T Number](t *testing.T, p Polygon[T], vertices []Point[T]) {
	t.Helper()

	testVertices(t, p.Vertices, vertices)
}

func testVertices[T Number](t *testing.T, vertices []Point[T], points []Point[T]) {
	t.Helper()

	assert.Equal(t, len(vertices), len(points))
	for i := 0; i < len(vertices); i++ {
		testPoint(t, vertices[i], points[i].X, points[i].Y)
	}
}

package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestPolygon_New(t *testing.T) {
	vertices := []Point[int]{P(0, 0), P(2, 0), P(2, 2), P(0, 2)}

	testPolygon(t, Polygon[int]{vertices}, vertices)
}

func TestPolygon_Center(t *testing.T) {
	testPoint(t, Polygon[int]{[]Point[int]{P(0, 0), P(2, 0), P(2, 2), P(0, 2)}}.Center(), 1, 1)
	testPoint(t, Polygon[float64]{[]Point[float64]{P(0.0, 0.0), P(2.0, 0.0), P(2.0, 1.0)}}.Center(), 4.0/3.0, 1.0/3.0)
}

func TestPolygon_Translate(t *testing.T) {
	testPolygon(t, Polygon[int]{[]Point[int]{P(0, 0), P(2, 0)}}.Translate(V(1, -1)), []Point[int]{
		P(1, -1),
		P(3, -1),
	})
}

func TestPolygon_MoveTo(t *testing.T) {
	testPolygon(t, Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0)}}.MoveTo(P(10, 10)), []Point[int]{
		P(9, 10),
		P(11, 10),
	})
}

func TestPolygon_Scale(t *testing.T) {
	testPolygon(t, Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0)}}.Scale(2), []Point[int]{
		P(-1, 0),
		P(3, 0),
	})
	testPolygon(t, Polygon[float64]{Vertices: []Point[float64]{P(0.0, 0.0), P(2.0, 1.0)}}.ScaleXY(0.5, 2.5), []Point[float64]{
		P(0.5, -0.75),
		P(1.5, 1.75),
	})
}

func TestPolygon_Immutable(t *testing.T) {
	p := Polygon[int]{Vertices: []Point[int]{P(0, 0), P(2, 0)}}

	p.Translate(V(1, -1))
	p.MoveTo(P(10, 10))
	p.Scale(2)
	p.ScaleXY(2, 3)

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

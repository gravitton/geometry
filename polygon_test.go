package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestPolygon_New(t *testing.T) {
	vertices := []Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}

	testPolygon(t, Pol(vertices), vertices)
}

func TestPolygon_Center(t *testing.T) {
	testPoint(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}).Center(), 1, 1)
	testPoint(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 1.0)}).Center(), 4.0/3.0, 1.0/3.0)
}

func TestPolygon_Translate(t *testing.T) {
	testPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0)}).Translate(Vec(1, -1)), []Point[int]{
		Pt(1, -1),
		Pt(3, -1),
	})
}

func TestPolygon_MoveTo(t *testing.T) {
	testPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0)}).MoveTo(Pt(10, 10)), []Point[int]{
		Pt(9, 10),
		Pt(11, 10),
	})
}

func TestPolygon_Scale(t *testing.T) {
	testPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0)}).Scale(2), []Point[int]{
		Pt(-1, 0),
		Pt(3, 0),
	})
	testPolygon(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 1.0)}).ScaleXY(0.5, 2.5), []Point[float64]{
		Pt(0.5, -0.75),
		Pt(1.5, 1.75),
	})
}

func TestPolygon_Immutable(t *testing.T) {
	p := Pol([]Point[int]{Pt(0, 0), Pt(2, 0)})

	p.Translate(Vec(1, -1))
	p.MoveTo(Pt(10, 10))
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

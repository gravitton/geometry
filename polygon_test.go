package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

var (
	polygonInt   = Polygon[int]{[]Point[int]{{0, 0}, {2, 0}, {2, 2}, {0, 2}}}
	polygonFloat = Polygon[float64]{[]Point[float64]{{0.0, 0.0}, {2.5, 0.5}, {2.0, 1.0}}}
)

func TestPolygon_New(t *testing.T) {
	AssertPolygon(t, polygonInt, polygonInt.Vertices)
	AssertPolygon(t, polygonFloat, polygonFloat.Vertices)
}

func TestPolygon_Center(t *testing.T) {
	AssertPoint(t, polygonInt.Center(), 1, 1)
	AssertPoint(t, polygonFloat.Center(), 1.5, 0.5)
}

func TestPolygon_Translate(t *testing.T) {
	AssertPolygon(t, polygonInt.Translate(Vec(1, -1)), []Point[int]{
		Pt(1, -1),
		Pt(3, -1),
		Pt(3, 1),
		Pt(1, 1),
	})
}

func TestPolygon_MoveTo(t *testing.T) {
	AssertPolygon(t, polygonInt.MoveTo(Pt(10, 10)), []Point[int]{
		Pt(9, 9),
		Pt(11, 9),
		Pt(11, 11),
		Pt(9, 11),
	})
}

func TestPolygon_Scale(t *testing.T) {
	AssertPolygon(t, polygonInt.Scale(2), []Point[int]{
		Pt(-1, -1),
		Pt(3, -1),
		Pt(3, 3),
		Pt(-1, 3),
	})
	AssertPolygon(t, polygonFloat.ScaleXY(0.5, 2.5), []Point[float64]{
		Pt(0.75, -0.75),
		Pt(2, 0.5),
		Pt(1.75, 1.75),
	})
}

func TestPolygon_Equal(t *testing.T) {
	assert.True(t, polygonInt.Equal(polygonInt))
	assert.False(t, polygonInt.Equal(Pol([]Point[int]{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {3, 2}})))
	assert.False(t, polygonInt.Equal(Pol([]Point[int]{{0, 0}, {2, 0}, {3, 2}, {0, 2}})))
	assert.True(t, polygonFloat.Equal(polygonFloat))
}

func TestPolygon_IsZeo(t *testing.T) {
	assert.False(t, polygonInt.IsZero())
	assert.True(t, Polygon[int]{}.IsZero())
	assert.False(t, Polygon[int]{[]Point[int]{}}.IsZero())
}

func TestPolygon_Empty(t *testing.T) {
	assert.False(t, polygonInt.Empty())
	assert.True(t, Polygon[int]{}.Empty())
	assert.True(t, Polygon[int]{[]Point[int]{}}.Empty())
	assert.False(t, polygonFloat.Empty())
}

func TestPolygon_Int(t *testing.T) {
	AssertPolygon(t, polygonInt.Int(), []Point[int]{
		Pt(0, 0),
		Pt(2, 0),
		Pt(2, 2),
		Pt(0, 2),
	})
	AssertPolygon(t, polygonFloat.Int(), []Point[int]{
		Pt(0, 0),
		Pt(3, 1),
		Pt(2, 1),
	})
}

func TestPolygon_Float(t *testing.T) {
	AssertPolygon(t, polygonInt.Float(), []Point[float64]{
		Pt(0.0, 0.0),
		Pt(2.0, 0.0),
		Pt(2.0, 2.0),
		Pt(0.0, 2.0),
	})
	AssertPolygon(t, polygonFloat.Float(), []Point[float64]{
		Pt(0.0, 0.0),
		Pt(2.5, 0.5),
		Pt(2.0, 1.0),
	})
}

func TestPolygon_Immutable(t *testing.T) {
	p := Pol([]Point[int]{Pt(0, 0), Pt(2, 0)})

	p.Translate(Vec(1, -1))
	p.MoveTo(Pt(10, 10))
	p.Scale(2)
	p.ScaleXY(2, 3)

	AssertPoint(t, p.Vertices[0], 0, 0)
	AssertPoint(t, p.Vertices[1], 2, 0)
}

func TestPolygon_String(t *testing.T) {
	assert.Equal(t, polygonInt.String(), "Pol((0,0), (2,0), (2,2), (0,2))")
	assert.Equal(t, polygonFloat.String(), "Pol((0,0), (2.50,0.50), (2,1))")
}

func TestPolygon_Marshall(t *testing.T) {
	assert.JSON(t, polygonInt, `[{"x":0,"y":0},{"x":2,"y":0},{"x":2,"y":2},{"x":0,"y":2}]`)
	assert.JSON(t, polygonFloat, `[{"x":0,"y":0},{"x":2.5,"y":0.5},{"x":2,"y":1}]`)
}

func TestPolygon_Unmarshall(t *testing.T) {
	var p1 Polygon[int]
	assert.NoError(t, json.Unmarshal([]byte(`[{"x":0,"y":0},{"x":2,"y":0},{"x":2,"y":2},{"x":0,"y":2}]`), &p1))
	AssertPolygon(t, p1, nil)

	var p2 Polygon[float64]
	assert.NoError(t, json.Unmarshal([]byte(`[{"x":0,"y":0},{"x":2.5,"y":0.5},{"x":2,"y":1}]`), &p2))
	AssertPolygon(t, p2, nil)
}

package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

var (
	// (0,1)-(2,4)
	rectInt = Rectangle[int]{Point[int]{1, 2}, Size[int]{2, 3}}
	// (0.0,-2.05)-(1.2,1.55)
	rectFloat = Rectangle[float64]{Point[float64]{0.6, -0.25}, Size[float64]{1.2, 3.6}}
)

func TestRectangle_New(t *testing.T) {
	assertRect(t, Rect(Pt(10, 16), Sz(3, 4)), 10, 16, 3, 4)
	assertRect(t, Rect[float64](Pt(0.5, -1.25), Sz(2.5, 3.75)), 0.5, -1.25, 2.5, 3.75)

	assertRect(t, RectFromMin(Pt(0, 0), Sz(4, 2)), 2, 1, 4, 2)
	assertRect(t, RectFromMin(Pt(0.0, 0.0), Sz(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	assertRect(t, RectFromMinMax(Pt(0, 0), Pt(4, 2)), 2, 1, 4, 2)
	assertRect(t, RectFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	assertRect(t, RectFromSize(Sz(4, 2)), 2, 1, 4, 2)
	assertRect(t, RectFromSize(Sz(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	// int(1x1) rectangle
	assertRect(t, Rect(Pt(0, 0), Sz(1, 1)), 0, 0, 1, 1)
	assertRect(t, RectFromMin(Pt(0, 0), Sz(1, 1)), 0, 0, 1, 1)
}

func TestRectangle_Translate(t *testing.T) {
	assertRect(t, rectInt.Translate(Vec(3, -2)), 4, 0, 2, 3)
	assertRect(t, rectFloat.Translate(Vec(100.1, -0.1)), 100.7, -0.35, 1.2, 3.6)
}

func TestRectangle_MoveTo(t *testing.T) {
	assertRect(t, rectInt.MoveTo(Pt(3, -2)), 3, -2, 2, 3)
	assertRect(t, rectFloat.MoveTo(Pt(100.1, -0.1)), 100.1, -0.1, 1.2, 3.6)
}

func TestRectangle_Scale(t *testing.T) {
	assertRect(t, rectInt.Scale(2.5), 1, 2, 5, 8)
	assertRect(t, rectInt.ScaleXY(2, 3), 1, 2, 4, 9)
	assertRect(t, rectFloat.Scale(2.5), 0.6, -0.25, 3.0, 9)
	assertRect(t, rectFloat.ScaleXY(-1.5, 2), 0.6, -0.25, -1.8, 7.2)
}

func TestRectangle_Resize(t *testing.T) {
	assertRect(t, rectInt.Resize(Sz(8, 9)), 1, 2, 8, 9)
	assertRect(t, rectFloat.Resize(Sz(3.1, 0.2)), 0.6, -0.25, 3.1, 0.2)
}

func TestRectangle_Grow(t *testing.T) {
	assertRect(t, rectInt.Grow(2), 1, 2, 4, 5)
	assertRect(t, rectInt.GrowXY(2, 3), 1, 2, 4, 6)
	assertRect(t, rectFloat.Grow(0.1), 0.6, -0.25, 1.3, 3.7)
	assertRect(t, rectFloat.GrowXY(0.1, 0.2), 0.6, -0.25, 1.3, 3.8)
}

func TestRectangle_Shrink(t *testing.T) {
	assertRect(t, rectInt.Shrink(1), 1, 2, 1, 2)
	assertRect(t, rectInt.ShrinkXY(1, 2), 1, 2, 1, 1)
	assertRect(t, rectFloat.Shrink(0.1), 0.6, -0.25, 1.1, 3.5)
	assertRect(t, rectFloat.ShrinkXY(0.1, 0.2), 0.6, -0.25, 1.1, 3.4)
}

func TestRectangle_Inset(t *testing.T) {
	assertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(1, 1, 1, 1)), 0, 0, 8, 8)
	assertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(3, 1, 1, 5)), 2, -1, 4, 6)
	assertRect(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), 1.5, -0.75, 11.0, 8.5)
}

func TestRectangle_Width(t *testing.T) {
	assert.Equal(t, rectInt.Width(), 2)
	assert.EqualDelta(t, rectFloat.Width(), 1.2, Delta)
}

func TestRectangle_Height(t *testing.T) {
	assert.Equal(t, rectInt.Height(), 3)
	assert.EqualDelta(t, rectFloat.Height(), 3.6, Delta)
}

func TestRectangle_Min(t *testing.T) {
	assertPoint(t, rectInt.Min(), 0, 1)
	assertPoint(t, rectFloat.Min(), 0.0, -2.05)

	assert.Equal(t, rectInt.Min(), rectInt.TopLeft())
}

func TestRectangle_Max(t *testing.T) {
	assertPoint(t, rectInt.Max(), 2, 4)
	assertPoint(t, rectFloat.Max(), 1.2, 1.55)

	assert.Equal(t, rectInt.Max(), rectInt.BottomRight())
}

func TestRectangle_Vertices(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	vertices := r.Vertices()

	assert.Equal(t, len(vertices), 4)
	assertPoint(t, vertices[0], -1, -1)
	assertPoint(t, vertices[1], -1, 1)
	assertPoint(t, vertices[2], 1, 1)
	assertPoint(t, vertices[3], 1, -1)

	assert.Equal(t, vertices[0], r.TopLeft())
	assert.Equal(t, vertices[1], r.BottomLeft())
	assert.Equal(t, vertices[2], r.BottomRight())
	assert.Equal(t, vertices[3], r.TopRight())
}

func TestRectangle_Edges(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	edges := r.Edges()

	assert.Equal(t, len(edges), 4)
	assertLine(t, edges[0], -1, -1, -1, 1)
	assertLine(t, edges[1], -1, 1, 1, 1)
	assertLine(t, edges[2], 1, 1, 1, -1)
	assertLine(t, edges[3], 1, -1, -1, -1)

	assert.Equal(t, edges[0].Start, r.TopLeft())
	assert.Equal(t, edges[1].Start, r.BottomLeft())
	assert.Equal(t, edges[2].Start, r.BottomRight())
	assert.Equal(t, edges[3].Start, r.TopRight())
}

func TestRectangle_Area(t *testing.T) {
	assert.Equal(t, rectInt.Area(), 6)
	assert.EqualDelta(t, rectFloat.Area(), 4.32, Delta)
}

func TestRectangle_Perimeter(t *testing.T) {
	assert.Equal(t, rectInt.Perimeter(), 10)
	assert.EqualDelta(t, rectFloat.Perimeter(), 9.6, Delta)
}

func TestRectangle_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, rectInt.AspectRatio(), 2.0/3.0, Delta)
	assert.EqualDelta(t, rectFloat.AspectRatio(), 1.0/3.0, Delta)
}

func TestRectangle_Bounds(t *testing.T) {
	assertRect(t, rectInt.Bounds(), 1, 2, 2, 3)
	assertRect(t, rectFloat.Bounds(), 0.6, -0.25, 1.2, 3.6)
}

func TestRectangle_Clamp(t *testing.T) {
	assertPoint(t, rectInt.Clamp(Pt(2, 2)), 2, 2)
	assertPoint(t, rectInt.Clamp(Pt(10, 10)), 2, 4)
	assertPoint(t, rectFloat.Clamp(Pt(-1.0, 1.2)), 0.0, 1.2)
}

func TestRectangle_Equal(t *testing.T) {
	assert.False(t, rectInt.Equal(Rect(Pt(3, -3), Sz(3, 4))))
	assert.True(t, rectInt.Equal(rectInt))

	assert.False(t, rectFloat.Equal(Rect(Pt(100.1, -0.1), Sz(1.2, 3.4))))
	assert.True(t, rectFloat.Equal(rectFloat))
}

func TestRectangle_IsZero(t *testing.T) {
	assert.False(t, rectInt.IsZero())
	assert.True(t, Rectangle[int]{}.IsZero())

	assert.False(t, rectFloat.IsZero())
	assert.True(t, Rectangle[float64]{}.IsZero())
}

func TestRectangle_Contains(t *testing.T) {
	r := rectInt
	assert.False(t, r.Contains(Pt(3, 0)))
	assert.True(t, r.Contains(Pt(1, 1)))

	r2 := rectFloat
	assert.False(t, r2.Contains(Pt(-0.1, 0)))
	assert.True(t, r2.Contains(Pt(0.25, -0.75)))
}

func TestRectangle_Polygon(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	p := r.Polygon()

	assert.Equal(t, p.Vertices, r.Vertices())
	assert.NotSame(t, p.Vertices, r.Vertices())
}

func TestRectangle_Int(t *testing.T) {
	assertRect(t, rectInt.Int(), 1, 2, 2, 3)
	assertRect(t, rectFloat.Int(), 1, 0, 1, 4)
}

func TestRectangle_Float(t *testing.T) {
	assertRect(t, rectInt.Float(), 1.0, 2.0, 2.0, 3.0)
	assertRect(t, rectFloat.Float(), 0.6, -0.25, 1.2, 3.6)
}

func TestRectangle_String(t *testing.T) {
	assert.Equal(t, rectInt.String(), "(0,1)-(2,4)")
	assert.Equal(t, rectFloat.String(), "(0,-2.05)-(1.20,1.55)")
}

func TestRectangle_Marshall(t *testing.T) {
	assert.JSON(t, rectInt, `{"x":1,"y":2,"w":2,"h":3}`)
	assert.JSON(t, rectFloat, `{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`)
}

func TestRectangle_Unmarshall(t *testing.T) {
	var r1 Rectangle[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":3}`), &r1))
	assertRect(t, r1, 1, 2, 2, 3)

	var r2 Rectangle[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`), &r2))
	assertRect(t, r2, 0.6, -0.25, 1.2, 3.6)
}

func TestRectangle_Immutable(t *testing.T) {
	r := rectInt

	r.Translate(Vec(3, -2))
	r.MoveTo(Pt(4, 3))
	r.Scale(2)
	r.ScaleXY(3, 4)
	r.Resize(Sz(5, 6))
	r.Grow(1)
	r.GrowXY(1, 2)
	r.Shrink(2)
	r.ShrinkXY(2, 3)

	assertRect(t, r, 1, 2, 2, 3)
}

func assertRect[T Number](t *testing.T, r Rectangle[T], cx, cy, w, h T) bool {
	t.Helper()

	ok := true

	if !assertPoint(t, r.Center, cx, cy, "Center.") {
		ok = false
	}
	if !assertSize(t, r.Size, w, h, "Size.") {
		ok = false
	}

	return ok
}

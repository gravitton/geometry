package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestRectangle_New(t *testing.T) {
	testRect(t, Rect(Pt(10, 16), Sz(3, 4)), 10, 16, 3, 4)
	testRect(t, Rect[float64](Pt(0.5, -1.25), Sz(2.5, 3.75)), 0.5, -1.25, 2.5, 3.75)

	testRect(t, RectFromMin(Pt(0, 0), Sz(4, 2)), 2, 1, 4, 2)
	testRect(t, RectFromMin(Pt(0.0, 0.0), Sz(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	testRect(t, RectFromMinMax(Pt(0, 0), Pt(4, 2)), 2, 1, 4, 2)
	testRect(t, RectFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	testRect(t, RectFromSize(Sz(4, 2)), 2, 1, 4, 2)
	testRect(t, RectFromSize(Sz(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	// int(1x1) rectangle
	testRect(t, Rect(Pt(0, 0), Sz(1, 1)), 0, 0, 1, 1)
	testRect(t, RectFromMin(Pt(0, 0), Sz(1, 1)), 0, 0, 1, 1)
}

func TestRectangle_Translate(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).Translate(Vec(3, -2)), 4, 0, 3, 4)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Translate(Vec(100.1, -0.1)), 100.5, -0.35, 1.2, 3.4)
}

func TestRectangle_MoveTo(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).MoveTo(Pt(3, -2)), 3, -2, 3, 4)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).MoveTo(Pt(100.1, -0.1)), 100.1, -0.1, 1.2, 3.4)
}

func TestRectangle_Scale(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).Scale(2.5), 1, 2, 8, 10)
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).ScaleXY(2, 3), 1, 2, 6, 12)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Scale(2.5), 0.4, -0.25, 3.0, 8.5)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).ScaleXY(-1.5, 2), 0.4, -0.25, -1.8, 6.8)
}

func TestRectangle_Resize(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).Resize(Sz(8, 9)), 1, 2, 8, 9)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Resize(Sz(3.1, 0.2)), 0.4, -0.25, 3.1, 0.2)
}

func TestRectangle_Grow(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).Grow(2), 1, 2, 5, 6)
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).GrowXY(2, 3), 1, 2, 5, 7)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Grow(0.1), 0.4, -0.25, 1.3, 3.5)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).GrowXY(0.1, 0.2), 0.4, -0.25, 1.3, 3.6)
}

func TestRectangle_Shrink(t *testing.T) {
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).Shrink(1), 1, 2, 2, 3)
	testRect(t, Rect(Pt(1, 2), Sz(3, 4)).ShrinkXY(1, 2), 1, 2, 2, 2)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Shrink(0.1), 0.4, -0.25, 1.1, 3.3)
	testRect(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).ShrinkXY(0.1, 0.2), 0.4, -0.25, 1.1, 3.2)
}

func TestRectangle_Inset(t *testing.T) {
	testRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(1, 1, 1, 1)), 0, 0, 8, 8)
	testRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(3, 1, 1, 5)), 2, -1, 4, 6)
	testRect(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), 1.5, -0.75, 11.0, 8.5)
}

func TestRectangle_Width(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(3, 4)).Width(), 3)
	assert.EqualDelta(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Width(), 1.2, Delta)
}

func TestRectangle_Height(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(3, 4)).Height(), 4)
	assert.EqualDelta(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Height(), 3.4, Delta)
}

func TestRectangle_Min(t *testing.T) {
	testPoint(t, Rect(Pt(1, 2), Sz(4, 2)).Min(), -1, 1)
	testPoint(t, Rect(Pt(0.5, -1.5), Sz(1.0, 3.0)).Min(), 0.0, -3.0)
}

func TestRectangle_Max(t *testing.T) {
	testPoint(t, Rect(Pt(1, 2), Sz(4, 2)).Max(), 3, 3)
	testPoint(t, Rect(Pt(0.5, -1.5), Sz(1.0, 3.0)).Max(), 1.0, 0.0)
}

func TestRectangle_Vertices(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	vertices := r.Vertices()

	assert.Equal(t, len(vertices), 4)
	testPoint(t, vertices[0], -1, -1)
	testPoint(t, vertices[1], -1, 1)
	testPoint(t, vertices[2], 1, 1)
	testPoint(t, vertices[3], 1, -1)

	assert.Equal(t, vertices[0], r.TopLeft())
	assert.Equal(t, vertices[1], r.BottomLeft())
	assert.Equal(t, vertices[2], r.BottomRight())
	assert.Equal(t, vertices[3], r.TopRight())
}

func TestRectangle_Area(t *testing.T) {
	assert.Equal(t, Rect(Pt(0, 0), Sz(4, 2)).Area(), 8)
	assert.EqualDelta(t, Rect(Pt(0.0, 0.0), Sz(1.5, 0.5)).Area(), 0.75, Delta)
}

func TestRectangle_Perimeter(t *testing.T) {
	assert.Equal(t, Rect(Pt(0, 0), Sz(4, 2)).Perimeter(), 12)
	assert.EqualDelta(t, Rect(Pt(0.0, 0.0), Sz(1.5, 0.5)).Perimeter(), 4.0, Delta)
}

func TestRectangle_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, Rect(Pt(0, 0), Sz(4, 2)).AspectRatio(), 2.0, Delta)
	assert.EqualDelta(t, Rect(Pt(0.0, 0.0), Sz(1.5, 0.5)).AspectRatio(), 3.0, Delta)
}

func TestRectangle_Equal(t *testing.T) {
	assert.False(t, Rect(Pt(1, 2), Sz(3, 4)).Equal(Rect(Pt(3, -3), Sz(3, 4))))
	assert.True(t, Rect(Pt(1, 2), Sz(3, 4)).Equal(Rect(Pt(1, 2), Sz(3, 4))))

	assert.False(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Equal(Rect(Pt(100.1, -0.1), Sz(1.2, 3.4))))
	assert.True(t, Rect(Pt(0.4, -0.25), Sz(1.2, 3.4)).Equal(Rect(Pt(0.4, -0.25), Sz(1.2, 3.4))))
}

func TestRectangle_Contains(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(4, 2))
	assert.False(t, r.Contains(Pt(3, 0)))
	assert.True(t, r.Contains(Pt(1, 0)))

	r2 := Rect(Pt(0.5, -1.5), Sz(1.0, 3.0))
	assert.False(t, r2.Contains(Pt(-0.1, 0)))
	assert.True(t, r2.Contains(Pt(0.25, -0.75)))
}

func TestRectangle_ToPolygon(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	p := r.ToPolygon()

	assert.Equal(t, p.Vertices, r.Vertices())
	assert.NotSame(t, p.Vertices, r.Vertices())
}

func TestRectangle_String(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(4, 2)).String(), "(-1,1)-(3,3)")
	assert.Equal(t, Rect(Pt(100, -34.0000115), Sz(0.2, 0.4)).String(), "(99.90,-34.20)-(100.10,-33.80)")
}

func TestRectangle_Immutable(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(3, 4))

	r.Translate(Vec(3, -2))
	r.MoveTo(Pt(4, 3))
	r.Scale(2)
	r.ScaleXY(3, 4)
	r.Resize(Sz(5, 6))
	r.Grow(1)
	r.GrowXY(1, 2)
	r.Shrink(2)
	r.ShrinkXY(2, 3)

	testRect(t, r, 1, 2, 3, 4)
}

func testRect[T Number](t *testing.T, r Rectangle[T], cx, cy, w, h T) {
	t.Helper()

	testPoint(t, r.Center, cx, cy)
	testSize(t, r.Size, w, h)
}

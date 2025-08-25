package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestRectangle_New(t *testing.T) {
	testRect(t, R(P(10, 16), S(3, 4)), 10, 16, 3, 4)
	testRect(t, R[float64](P(0.5, -1.25), S(2.5, 3.75)), 0.5, -1.25, 2.5, 3.75)

	testRect(t, RectangleFromMin(P(0, 0), S(4, 2)), 2, 1, 4, 2)
	testRect(t, RectangleFromMin(P(0.0, 0.0), S(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	testRect(t, RectangleFromMinMax(P(0, 0), P(4, 2)), 2, 1, 4, 2)
	testRect(t, RectangleFromMinMax(P(0.0, 0.0), P(1.0, 3.0)), 0.5, 1.5, 1.0, 3.0)

	testRect(t, RectangleFromSize(S(4, 2)), 0, 0, 4, 2)
	testRect(t, RectangleFromSize(S(1.0, 3.0)), 0.0, 0.0, 1.0, 3.0)
}

func TestRectangle_Translate(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).Translate(V(3, -2)), 4, 0, 3, 4)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).Translate(V(100.1, -0.1)), 100.5, -0.35, 1.2, 3.4)
}

func TestRectangle_MoveTo(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).MoveTo(P(3, -2)), 3, -2, 3, 4)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).MoveTo(P(100.1, -0.1)), 100.1, -0.1, 1.2, 3.4)
}

func TestRectangle_Scale(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).Scale(2.5), 1, 2, 8, 10)
	testRect(t, R(P(1, 2), S(3, 4)).ScaleXY(2, 3), 1, 2, 6, 12)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).Scale(2.5), 0.4, -0.25, 3.0, 8.5)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).ScaleXY(-1.5, 2), 0.4, -0.25, -1.8, 6.8)
}

func TestRectangle_Resize(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).Resize(S(8, 9)), 1, 2, 8, 9)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).Resize(S(3.1, 0.2)), 0.4, -0.25, 3.1, 0.2)
}

func TestRectangle_Expand(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).Grow(2), 1, 2, 5, 6)
	testRect(t, R(P(1, 2), S(3, 4)).GrowXY(2, 3), 1, 2, 5, 7)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).Grow(0.1), 0.4, -0.25, 1.3, 3.5)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).GrowXY(0.1, 0.2), 0.4, -0.25, 1.3, 3.6)
}

func TestRectangle_Shrunk(t *testing.T) {
	testRect(t, R(P(1, 2), S(3, 4)).Shrink(1), 1, 2, 2, 3)
	testRect(t, R(P(1, 2), S(3, 4)).ShrinkXY(1, 2), 1, 2, 2, 2)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).Shrink(0.1), 0.4, -0.25, 1.1, 3.3)
	testRect(t, R(P(0.4, -0.25), S(1.2, 3.4)).ShrinkXY(0.1, 0.2), 0.4, -0.25, 1.1, 3.2)
}

func TestRectangle_Width(t *testing.T) {
	assert.Equal(t, R(P(1, 2), S(3, 4)).Width(), 3)
	assert.EqualDelta(t, R(P(0.4, -0.25), S(1.2, 3.4)).Width(), 1.2, Delta)
}

func TestRectangle_Height(t *testing.T) {
	assert.Equal(t, R(P(1, 2), S(3, 4)).Height(), 4)
	assert.EqualDelta(t, R(P(0.4, -0.25), S(1.2, 3.4)).Height(), 3.4, Delta)
}

func TestRectangle_Min(t *testing.T) {
	testPoint(t, R(P(1, 2), S(4, 2)).Min(), -1, 1)
	testPoint(t, R(P(0.5, -1.5), S(1.0, 3.0)).Min(), 0.0, -3.0)
}

func TestRectangle_Max(t *testing.T) {
	testPoint(t, R(P(1, 2), S(4, 2)).Max(), 3, 3)
	testPoint(t, R(P(0.5, -1.5), S(1.0, 3.0)).Max(), 1.0, 0.0)
}

func TestRectangle_Vertices(t *testing.T) {
	r := R(P(0, 0), S(2, 2))
	vertices := r.Vertices()

	assert.Equal(t, len(vertices), 4)
	testPoint(t, vertices[0], -1, -1)
	testPoint(t, vertices[1], 1, -1)
	testPoint(t, vertices[2], 1, 1)
	testPoint(t, vertices[3], -1, 1)

	assert.Equal(t, vertices[0], r.BottomLeft())
	assert.Equal(t, vertices[1], r.BottomRight())
	assert.Equal(t, vertices[2], r.TopRight())
	assert.Equal(t, vertices[3], r.TopLeft())
}

func TestRectangle_Area(t *testing.T) {
	assert.Equal(t, R(P(0, 0), S(4, 2)).Area(), 8)
	assert.EqualDelta(t, R(P(0.0, 0.0), S(1.5, 0.5)).Area(), 0.75, Delta)
}

func TestRectangle_Perimeter(t *testing.T) {
	assert.Equal(t, R(P(0, 0), S(4, 2)).Perimeter(), 12)
	assert.EqualDelta(t, R(P(0.0, 0.0), S(1.5, 0.5)).Perimeter(), 4.0, Delta)
}

func TestRectangle_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, R(P(0, 0), S(4, 2)).AspectRatio(), 2.0, Delta)
	assert.EqualDelta(t, R(P(0.0, 0.0), S(1.5, 0.5)).AspectRatio(), 3.0, Delta)
}

func TestRectangle_Equal(t *testing.T) {
	assert.False(t, R(P(1, 2), S(3, 4)).Equal(R(P(3, -3), S(3, 4))))
	assert.True(t, R(P(1, 2), S(3, 4)).Equal(R(P(1, 2), S(3, 4))))

	assert.False(t, R(P(0.4, -0.25), S(1.2, 3.4)).Equal(R(P(100.1, -0.1), S(1.2, 3.4))))
	assert.True(t, R(P(0.4, -0.25), S(1.2, 3.4)).Equal(R(P(0.4, -0.25), S(1.2, 3.4))))
}

func TestRectangle_Contains(t *testing.T) {
	r := R(P(0, 0), S(4, 2))
	assert.False(t, r.Contains(P(3, 0)))
	assert.True(t, r.Contains(P(1, 0)))

	r2 := R(P(0.5, -1.5), S(1.0, 3.0))
	assert.False(t, r2.Contains(P(-0.1, 0)))
	assert.True(t, r2.Contains(P(0.25, -0.75)))
}

func TestRectangle_ToPolygon(t *testing.T) {
	r := R(P(0, 0), S(2, 2))
	p := r.ToPolygon()

	assert.Equal(t, p.Vertices, r.Vertices())
	assert.NotSame(t, p.Vertices, r.Vertices())
}

func TestRectangle_String(t *testing.T) {
	assert.Equal(t, R(P(1, 2), S(4, 2)).String(), "(-1,1)-(3,3)")
	assert.Equal(t, R(P(100, -34.0000115), S(0.2, 0.4)).String(), "(99.90,-34.20)-(100.10,-33.80)")
}

func TestRectangle_Immutable(t *testing.T) {
	r := R(P(1, 2), S(3, 4))

	r.Translate(V(3, -2))
	r.MoveTo(P(4, 3))
	r.Scale(2)
	r.ScaleXY(3, 4)
	r.Resize(S(5, 6))
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

package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestRectangle_Constructor(t *testing.T) {
	AssertRect(t, Rect(Pt(10, 16), Sz(3, 4)), Rectangle[int]{Center: Pt(10, 16), Size: Sz(3, 4)})
	AssertRect(t, Rect[float64](Pt(0.5, -1.25), Sz(2.5, 3.75)), Rectangle[float64]{Center: Pt(0.5, -1.25), Size: Sz(2.5, 3.75)})

	AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
	AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
	AssertRect(t, RectangleFromMin(Pt(0.0, 0.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))

	AssertRect(t, RectangleFromMax(Pt(4, 2), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
	AssertRect(t, RectangleFromMax(Pt(5, 3), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
	AssertRect(t, RectangleFromMax(Pt(1.0, 3.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))

	AssertRect(t, RectangleFromMinMax(Pt(0, 0), Pt(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
	AssertRect(t, RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))

	AssertRect(t, RectangleFromSize(Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
	AssertRect(t, RectangleFromSize(Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))

	// int(1x1) rectangle
	AssertRect(t, Rect(Pt(0, 0), Sz(1, 1)), Rect(Pt(0, 0), Sz(1, 1)))
	AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(1, 1)), Rect(Pt(0, 0), Sz(1, 1)))
}

func TestRectangle_Translate(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Translate(Vec(3, -2)), Rect(Pt(4, 0), Sz(2, 3)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Translate(Vec(100.1, -0.1)), Rect(Pt(100.7, -0.35), Sz(1.2, 3.6)))
}

func TestRectangle_MoveTo(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).MoveTo(Pt(3, -2)), Rect(Pt(3, -2), Sz(2, 3)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).MoveTo(Pt(100.1, -0.1)), Rect(Pt(100.1, -0.1), Sz(1.2, 3.6)))
}

func TestRectangle_Scale(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Scale(2.5), Rect(Pt(1, 2), Sz(5, 8)))
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).ScaleXY(2, 3), Rect(Pt(1, 2), Sz(4, 9)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Scale(2.5), Rect(Pt(0.6, -0.25), Sz(3.0, 9)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ScaleXY(-1.5, 2), Rect(Pt(0.6, -0.25), Sz(-1.8, 7.2)))
}

func TestRectangle_Resize(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Resize(Sz(8, 9)), Rect(Pt(1, 2), Sz(8, 9)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Resize(Sz(3.1, 0.2)), Rect(Pt(0.6, -0.25), Sz(3.1, 0.2)))
}

func TestRectangle_Grow(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Grow(2), Rect(Pt(1, 2), Sz(4, 5)))
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).GrowXY(2, 3), Rect(Pt(1, 2), Sz(4, 6)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Grow(0.1), Rect(Pt(0.6, -0.25), Sz(1.3, 3.7)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).GrowXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.3, 3.8)))
}

func TestRectangle_Shrink(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(1), Rect(Pt(1, 2), Sz(1, 2)))
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).ShrinkXY(1, 2), Rect(Pt(1, 2), Sz(1, 1)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Shrink(0.1), Rect(Pt(0.6, -0.25), Sz(1.1, 3.5)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ShrinkXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.1, 3.4)))

	// clamped to zero — never negative size
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(100), Rect(Pt(1, 2), Sz(0, 0)))
}

func TestRectangle_Inset(t *testing.T) {
	AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(1, 1, 1, 1)), Rect(Pt(0, 0), Sz(8, 8)))
	AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(3, 1, 1, 5)), Rect(Pt(2, -1), Sz(4, 6)))
	AssertRect(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), Rect(Pt(1.5, -0.75), Sz(11.0, 8.5)))
}

func TestRectangle_Width(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Width(), 2)
	assert.EqualDelta(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Width(), 1.2, Delta)
}

func TestRectangle_Height(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Height(), 3)
	assert.EqualDelta(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Height(), 3.6, Delta)
}

func TestRectangle_Min(t *testing.T) {
	AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Min(), Pt(0, 1))
	AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Min(), Pt(0.0, -2.05))

	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Min(), Rect(Pt(1, 2), Sz(2, 3)).TopLeft())
}

func TestRectangle_Max(t *testing.T) {
	AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Max(), Pt(2, 4))
	AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Max(), Pt(1.2, 1.55))

	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Max(), Rect(Pt(1, 2), Sz(2, 3)).BottomRight())
}

func TestRectangle_Vertices(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	vertices := r.Vertices()

	assert.Equal(t, len(vertices), 4)
	AssertPoint(t, vertices[0], Pt(-1, -1))
	AssertPoint(t, vertices[1], Pt(1, -1))
	AssertPoint(t, vertices[2], Pt(1, 1))
	AssertPoint(t, vertices[3], Pt(-1, 1))

	assert.Equal(t, vertices[0], r.TopLeft())
	assert.Equal(t, vertices[1], r.TopRight())
	assert.Equal(t, vertices[2], r.BottomRight())
	assert.Equal(t, vertices[3], r.BottomLeft())
}

func TestRectangle_Edges(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	edges := r.Edges()

	assert.Equal(t, len(edges), 4)
	AssertLine(t, edges[0], Ln(Pt(-1, -1), Pt(1, -1)))
	AssertLine(t, edges[1], Ln(Pt(1, -1), Pt(1, 1)))
	AssertLine(t, edges[2], Ln(Pt(1, 1), Pt(-1, 1)))
	AssertLine(t, edges[3], Ln(Pt(-1, 1), Pt(-1, -1)))

	assert.Equal(t, edges[0].Start, r.TopLeft())
	assert.Equal(t, edges[1].Start, r.TopRight())
	assert.Equal(t, edges[2].Start, r.BottomRight())
	assert.Equal(t, edges[3].Start, r.BottomLeft())

	// Edges agrees with the dedicated accessors, in the same order
	assert.Equal(t, edges[0], r.TopEdge())
	assert.Equal(t, edges[1], r.RightEdge())
	assert.Equal(t, edges[2], r.BottomEdge())
	assert.Equal(t, edges[3], r.LeftEdge())

	// the edges form a closed chain: each starts where the previous ends
	for i, edge := range edges {
		assert.Equal(t, edge.Start, edges[(i+3)%4].End)
	}

	// vertices and edge starts agree
	for i, vertex := range r.Vertices() {
		assert.Equal(t, edges[i].Start, vertex)
	}
}

func TestRectangle_Anchor(t *testing.T) {
	r := RectangleFromMin(Pt(0, 0), Sz(10, 20))

	AssertPoint(t, r.Anchor(TopLeft), Pt(0, 0))
	AssertPoint(t, r.Anchor(TopRight), Pt(10, 0))
	AssertPoint(t, r.Anchor(BottomLeft), Pt(0, 20))
	AssertPoint(t, r.Anchor(BottomRight), Pt(10, 20))

	AssertPoint(t, r.Anchor(Top), Pt(5, 0))
	AssertPoint(t, r.Anchor(Bottom), Pt(5, 20))
	AssertPoint(t, r.Anchor(DirectionLeft), Pt(0, 10))
	AssertPoint(t, r.Anchor(DirectionRight), Pt(10, 10))

	AssertPoint(t, r.Anchor(DirectionNone), Pt(5, 10))

	// garbage input still wraps to a meaningful direction
	AssertPoint(t, r.Anchor(Direction(99)), Pt(r.BottomLeft().X, r.BottomLeft().Y))

	// corners and edge midpoints agree with the dedicated accessors
	assert.True(t, r.Anchor(TopLeft).Equal(r.TopLeft()))
	assert.True(t, r.Anchor(BottomRight).Equal(r.BottomRight()))
	assert.True(t, r.Anchor(Top).Equal(r.Top()))
	assert.True(t, r.Anchor(Bottom).Equal(r.Bottom()))
	assert.True(t, r.Anchor(DirectionLeft).Equal(r.Left()))
	assert.True(t, r.Anchor(DirectionRight).Equal(r.Right()))
	assert.True(t, r.Anchor(DirectionNone).Equal(r.Center))

	// odd integer extents split the way Min and Max do
	odd := RectangleFromMin(Pt(0, 0), Sz(3, 3))
	AssertPoint(t, odd.Anchor(TopLeft), Pt(0, 0))
	AssertPoint(t, odd.Anchor(BottomRight), Pt(3, 3))
	AssertPoint(t, odd.Top(), Pt(1, 0))
	AssertPoint(t, odd.Bottom(), Pt(1, 3))
	AssertPoint(t, odd.Left(), Pt(0, 1))
	AssertPoint(t, odd.Right(), Pt(3, 1))
}

func TestRectangle_Area(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Area(), 6)
	assert.EqualDelta(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Area(), 4.32, Delta)
}

func TestRectangle_Perimeter(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Perimeter(), 10)
	assert.EqualDelta(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Perimeter(), 9.6, Delta)
}

func TestRectangle_AspectRatio(t *testing.T) {
	assert.EqualDelta(t, Rect(Pt(1, 2), Sz(2, 3)).AspectRatio(), 2.0/3.0, Delta)
	assert.EqualDelta(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).AspectRatio(), 1.0/3.0, Delta)
}

func TestRectangle_Bounds(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Bounds(), Rect(Pt(1, 2), Sz(2, 3)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Bounds(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
}

func TestRectangle_Clamp(t *testing.T) {
	AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Clamp(Pt(2, 2)), Pt(2, 2))
	AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Clamp(Pt(10, 10)), Pt(2, 4))
	AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Clamp(Pt(-1.0, 1.2)), Pt(0.0, 1.2))
}

func TestRectangle_Equal(t *testing.T) {
	assert.False(t, Rect(Pt(1, 2), Sz(2, 3)).Equal(Rect(Pt(3, -3), Sz(3, 4))))
	assert.True(t, Rect(Pt(1, 2), Sz(2, 3)).Equal(Rect(Pt(1, 2), Sz(2, 3))))

	assert.False(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Equal(Rect(Pt(100.1, -0.1), Sz(1.2, 3.4))))
	assert.True(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Equal(Rect(Pt(0.6, -0.25), Sz(1.2, 3.6))))
}

func TestRectangle_IsZero(t *testing.T) {
	assert.False(t, Rect(Pt(1, 2), Sz(2, 3)).IsZero())
	assert.True(t, Rectangle[int]{}.IsZero())

	assert.False(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).IsZero())
	assert.True(t, Rectangle[float64]{}.IsZero())
}

func TestRectangle_Contains(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(2, 3))
	assert.False(t, r.Contains(Pt(3, 0)))
	assert.True(t, r.Contains(Pt(1, 1)))

	r2 := Rect(Pt(0.6, -0.25), Sz(1.2, 3.6))
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
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Int(), Rect(Pt(1, 2), Sz(2, 3)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Int(), Rect(Pt(1, 0), Sz(1, 4)))
}

func TestRectangle_Float(t *testing.T) {
	AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Float(), Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)))
	AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Float(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
}

func TestRectangle_String(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).String(), "(0,1)-(2,4)")
	assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).String(), "(0.00,-2.05)-(1.20,1.55)")
}

func TestRectangle_Marshall(t *testing.T) {
	assert.JSON(t, Rect(Pt(1, 2), Sz(2, 3)), `{"x":1,"y":2,"w":2,"h":3}`)
	assert.JSON(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)), `{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`)
}

func TestRectangle_Unmarshall(t *testing.T) {
	var r1 Rectangle[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":3}`), &r1))
	AssertRect(t, r1, Rect(Pt(1, 2), Sz(2, 3)))

	var r2 Rectangle[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`), &r2))
	AssertRect(t, r2, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
}

func TestRectangle_Immutable(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(2, 3))

	r.Translate(Vec(3, -2))
	r.MoveTo(Pt(4, 3))
	r.Scale(2)
	r.ScaleXY(3, 4)
	r.Resize(Sz(5, 6))
	r.Grow(1)
	r.GrowXY(1, 2)
	r.Shrink(2)
	r.ShrinkXY(2, 3)

	AssertRect(t, r, Rect(Pt(1, 2), Sz(2, 3)))
}

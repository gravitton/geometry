package geom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestRectangle_Constructor(t *testing.T) {
	t.Run("from center and size", func(t *testing.T) {
		AssertRect(t, Rect(Pt(10, 16), Sz(3, 4)), Rectangle[int]{Center: Pt(10, 16), Size: Sz(3, 4)})
		AssertRect(t, Rect(Pt(0.5, -1.25), Sz(2.5, 3.75)), Rectangle[float64]{Center: Pt(0.5, -1.25), Size: Sz(2.5, 3.75)})
	})
	t.Run("from min", func(t *testing.T) {
		AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
		AssertRect(t, RectangleFromMin(Pt(0.0, 0.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from max", func(t *testing.T) {
		AssertRect(t, RectangleFromMax(Pt(4, 2), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		AssertRect(t, RectangleFromMax(Pt(5, 3), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
		AssertRect(t, RectangleFromMax(Pt(1.0, 3.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from min and max", func(t *testing.T) {
		AssertRect(t, RectangleFromMinMax(Pt(0, 0), Pt(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		AssertRect(t, RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from size at the origin", func(t *testing.T) {
		AssertRect(t, RectangleFromSize(Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		AssertRect(t, RectangleFromSize(Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("unit integer rectangle", func(t *testing.T) {
		AssertRect(t, RectangleFromMin(Pt(0, 0), Sz(1, 1)), Rect(Pt(0, 0), Sz(1, 1)))
	})
}

func TestRectangle_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Translate(Vec(3, -2)), Rect(Pt(4, 0), Sz(2, 3)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Translate(Vec(100.1, -0.1)), Rect(Pt(100.7, -0.35), Sz(1.2, 3.6)))
	})
}

func TestRectangle_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).MoveTo(Pt(3, -2)), Rect(Pt(3, -2), Sz(2, 3)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).MoveTo(Pt(100.1, -0.1)), Rect(Pt(100.1, -0.1), Sz(1.2, 3.6)))
	})
}

func TestRectangle_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Scale(2.5), Rect(Pt(1, 2), Sz(5, 8)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Scale(2.5), Rect(Pt(0.6, -0.25), Sz(3.0, 9.0)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).ScaleXY(2, 3), Rect(Pt(1, 2), Sz(4, 9)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ScaleXY(-1.5, 2), Rect(Pt(0.6, -0.25), Sz(-1.8, 7.2)))
	})
}

func TestRectangle_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Resize(Sz(8, 9)), Rect(Pt(1, 2), Sz(8, 9)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Resize(Sz(3.1, 0.2)), Rect(Pt(0.6, -0.25), Sz(3.1, 0.2)))
	})
}

func TestRectangle_Grow(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Grow(2), Rect(Pt(1, 2), Sz(4, 5)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Grow(0.1), Rect(Pt(0.6, -0.25), Sz(1.3, 3.7)))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).GrowXY(2, 3), Rect(Pt(1, 2), Sz(4, 6)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).GrowXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.3, 3.8)))
	})
}

func TestRectangle_Shrink(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(1), Rect(Pt(1, 2), Sz(1, 2)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Shrink(0.1), Rect(Pt(0.6, -0.25), Sz(1.1, 3.5)))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).ShrinkXY(1, 2), Rect(Pt(1, 2), Sz(1, 1)))
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ShrinkXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.1, 3.4)))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(100), Rect(Pt(1, 2), Sz(0, 0)))
	})
}

func TestRectangle_Outset(t *testing.T) {
	t.Run("uniform padding keeps the center", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Outset(PadU(1)), Rect(Pt(0, 0), Sz(12, 12)))
	})
	t.Run("asymmetric padding shifts the center", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Outset(Pad(3, 1, 1, 5)), Rect(Pt(-2, -1), Sz(16, 14)))
	})
	t.Run("undoes an inset", func(t *testing.T) {
		r := RectangleFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		padding := Pad(1.0, 2.0, 3.0, 4.0)

		AssertRect(t, r.Inset(padding).Outset(padding), r)
		AssertRect(t, r.Outset(padding).Inset(padding), r)
	})
}

func TestRectangle_Inset(t *testing.T) {
	t.Run("uniform padding keeps the center", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(PadU(1)), Rect(Pt(0, 0), Sz(8, 8)))
	})
	t.Run("asymmetric padding shifts the center", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(3, 1, 1, 5)), Rect(Pt(2, 1), Sz(4, 6)))
	})
	t.Run("odd integer padding stays inside the rectangle", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(10, 10))

		AssertPoint(t, r.Inset(Pad(0, 0, 0, 1)).Min(), r.Min().AddXY(1, 0))
		AssertPoint(t, r.Inset(Pad(0, 0, 0, 1)).Max(), r.Max())
		AssertPoint(t, r.Inset(Pad(0, 3, 1, 0)).Min(), r.Min())
		AssertPoint(t, r.Inset(Pad(0, 3, 1, 0)).Max(), r.Max().AddXY(-3, -1))
	})
	t.Run("padding beyond the size collapses inside the rectangle", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(10, 10))

		AssertRect(t, r.Inset(Pad(0, 0, 0, 20)), Rect(Pt(5, 0), Sz(0, 10)))
		AssertRect(t, r.Inset(Pad(0, 20, 0, 0)), Rect(Pt(-5, 0), Sz(0, 10)))
		AssertRect(t, r.Inset(Pad(20, 0, 20, 0)), Rect(Pt(0, 5), Sz(10, 0)))
		AssertRect(t, r.Inset(PadU(20)), Rect(Pt(5, 5), Sz(0, 0)))
		AssertRect(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(0.0, 0.0, 0.0, 12.5)), Rect(Pt(5.0, 0.0), Sz(0.0, 10.0)))
	})
	t.Run("a collapsed rectangle is contained by the original", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(10, 10))

		for _, padding := range []Padding[int]{Pad(0, 0, 0, 20), Pad(0, 20, 0, 0), Pad(20, 0, 0, 0), Pad(0, 0, 20, 0), PadU(20), Pad(7, 7, 7, 7)} {
			inset := r.Inset(padding)

			assert.True(t, r.Contains(inset.Min()), fmt.Sprintf("%s min: ", padding))
			assert.True(t, r.Contains(inset.Max()), fmt.Sprintf("%s max: ", padding))
		}
	})
	t.Run("negative padding grows the edge", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), Rect(Pt(1.5, 0.75), Sz(11.0, 8.5)))
	})
	t.Run("each edge moves inward by its own padding", func(t *testing.T) {
		r := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
		padding := Pad(3.0, 1.0, 1.0, 5.0)
		inset := r.Inset(padding)

		AssertPoint(t, inset.Min(), r.Min().AddXY(padding.Left, padding.Top))
		AssertPoint(t, inset.Max(), r.Max().AddXY(-padding.Right, -padding.Bottom))
	})
}

func TestRectangle_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Width(), 2)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Width(), 1.2)
	})
}

func TestRectangle_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Height(), 3)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Height(), 3.6)
	})
}

func TestRectangle_Min(t *testing.T) {
	t.Run("top-left corner", func(t *testing.T) {
		AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Min(), Pt(0, 1))
		AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Min(), Pt(0.0, -2.05))
	})
	t.Run("agrees with the corner accessor", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))
		AssertPoint(t, r.Min(), r.TopLeft())
	})
}

func TestRectangle_Max(t *testing.T) {
	t.Run("bottom-right corner", func(t *testing.T) {
		AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Max(), Pt(2, 4))
		AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Max(), Pt(1.2, 1.55))
	})
	t.Run("agrees with the corner accessor", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))
		AssertPoint(t, r.Max(), r.BottomRight())
	})
}

func TestRectangle_MinMax(t *testing.T) {
	t.Run("both corners at once", func(t *testing.T) {
		for _, r := range rectFixtures {
			a, b := r.MinMax()

			AssertPoint(t, a, r.Min(), fmt.Sprintf("%s min: ", r))
			AssertPoint(t, b, r.Max(), fmt.Sprintf("%s max: ", r))
		}
	})
}

func TestRectangle_Anchor(t *testing.T) {
	r := RectangleFromMin(Pt(0, 0), Sz(10, 20))

	t.Run("corners", func(t *testing.T) {
		AssertPoint(t, r.Anchor(TopLeft), Pt(0, 0))
		AssertPoint(t, r.Anchor(TopRight), Pt(10, 0))
		AssertPoint(t, r.Anchor(BottomLeft), Pt(0, 20))
		AssertPoint(t, r.Anchor(BottomRight), Pt(10, 20))
	})
	t.Run("edge midpoints", func(t *testing.T) {
		AssertPoint(t, r.Anchor(Top), Pt(5, 0))
		AssertPoint(t, r.Anchor(Bottom), Pt(5, 20))
		AssertPoint(t, r.Anchor(DirectionLeft), Pt(0, 10))
		AssertPoint(t, r.Anchor(DirectionRight), Pt(10, 10))
	})
	t.Run("none is the center", func(t *testing.T) {
		AssertPoint(t, r.Anchor(DirectionNone), Pt(5, 10))
		AssertPoint(t, r.Anchor(DirectionNone), r.Center)
	})
	t.Run("out-of-range direction wraps", func(t *testing.T) {
		AssertPoint(t, r.Anchor(Direction(99)), r.BottomLeft())
	})
	t.Run("agrees with the dedicated accessors", func(t *testing.T) {
		AssertPoint(t, r.Anchor(TopLeft), r.TopLeft())
		AssertPoint(t, r.Anchor(TopRight), r.TopRight())
		AssertPoint(t, r.Anchor(BottomLeft), r.BottomLeft())
		AssertPoint(t, r.Anchor(BottomRight), r.BottomRight())
		AssertPoint(t, r.Anchor(Top), r.Top())
		AssertPoint(t, r.Anchor(Bottom), r.Bottom())
		AssertPoint(t, r.Anchor(DirectionLeft), r.Left())
		AssertPoint(t, r.Anchor(DirectionRight), r.Right())
	})
	t.Run("odd integer extents split like min and max", func(t *testing.T) {
		odd := RectangleFromMin(Pt(0, 0), Sz(3, 3))

		AssertPoint(t, odd.Anchor(TopLeft), Pt(0, 0))
		AssertPoint(t, odd.Anchor(BottomRight), Pt(3, 3))
		AssertPoint(t, odd.Top(), Pt(1, 0))
		AssertPoint(t, odd.Bottom(), Pt(1, 3))
		AssertPoint(t, odd.Left(), Pt(0, 1))
		AssertPoint(t, odd.Right(), Pt(3, 1))
	})
}

func TestRectangle_Edges(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	edges := r.Edges()

	t.Run("clockwise from the top", func(t *testing.T) {
		assert.Equal(t, len(edges), 4)
		AssertLine(t, edges[0], Ln(Pt(-1, -1), Pt(1, -1)))
		AssertLine(t, edges[1], Ln(Pt(1, -1), Pt(1, 1)))
		AssertLine(t, edges[2], Ln(Pt(1, 1), Pt(-1, 1)))
		AssertLine(t, edges[3], Ln(Pt(-1, 1), Pt(-1, -1)))
	})
	t.Run("agrees with the dedicated accessors", func(t *testing.T) {
		AssertLine(t, edges[0], r.TopEdge())
		AssertLine(t, edges[1], r.RightEdge())
		AssertLine(t, edges[2], r.BottomEdge())
		AssertLine(t, edges[3], r.LeftEdge())
	})
	t.Run("form a closed chain", func(t *testing.T) {
		for i, edge := range edges {
			AssertPoint(t, edge.Start, edges[(i+3)%4].End)
		}
	})
}

func TestRectangle_Vertices(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	vertices := r.Vertices()

	t.Run("clockwise from the top left", func(t *testing.T) {
		AssertVertices(t, vertices, []Point[int]{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}})
	})
	t.Run("agrees with the corner accessors", func(t *testing.T) {
		AssertPoint(t, vertices[0], r.TopLeft())
		AssertPoint(t, vertices[1], r.TopRight())
		AssertPoint(t, vertices[2], r.BottomRight())
		AssertPoint(t, vertices[3], r.BottomLeft())
	})
	t.Run("agrees with the edge starts", func(t *testing.T) {
		for i, edge := range r.Edges() {
			AssertPoint(t, edge.Start, vertices[i])
		}
	})
}

func TestRectangle_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Area(), 6)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Area(), 4.32)
	})
}

func TestRectangle_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Perimeter(), 10)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Perimeter(), 9.6)
	})
}

func TestRectangle_AspectRatio(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).AspectRatio(), 2.0/3.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).AspectRatio(), 1.0/3.0)
	})
}

func TestRectangle_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Bounds(), Rect(Pt(1, 2), Sz(2, 3)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Bounds(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
}

func TestRectangle_Clamp(t *testing.T) {
	t.Run("inside point is unchanged", func(t *testing.T) {
		AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Clamp(Pt(2, 2)), Pt(2, 2))
	})
	t.Run("outside point snaps to the edge", func(t *testing.T) {
		AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Clamp(Pt(10, 10)), Pt(2, 4))
		AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Clamp(Pt(-1.0, 1.2)), Pt(0.0, 1.2))
	})
}

func TestRectangle_Equal(t *testing.T) {
	t.Run("same rectangle", func(t *testing.T) {
		assert.True(t, Rect(Pt(1, 2), Sz(2, 3)).Equal(Rect(Pt(1, 2), Sz(2, 3))))
		assert.True(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Equal(Rect(Pt(0.6, -0.25), Sz(1.2, 3.6))))
	})
	t.Run("different rectangle", func(t *testing.T) {
		assert.False(t, Rect(Pt(1, 2), Sz(2, 3)).Equal(Rect(Pt(3, -3), Sz(3, 4))))
		assert.False(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Equal(Rect(Pt(100.1, -0.1), Sz(1.2, 3.4))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Equal(Rect(Pt(0.6, -0.250001), Sz(1.2, 3.6))))
	})
}

func TestRectangle_IsZero(t *testing.T) {
	t.Run("zero rectangle", func(t *testing.T) {
		assert.True(t, Rectangle[int]{}.IsZero())
		assert.True(t, Rectangle[float64]{}.IsZero())
	})
	t.Run("non-zero rectangle", func(t *testing.T) {
		assert.False(t, Rect(Pt(1, 2), Sz(2, 3)).IsZero())
		assert.False(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Rect(Pt(0.0, 0.000001), Sz(0.0, 0.0)).IsZero())
	})
}

func TestRectangle_Contains(t *testing.T) {
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Rect(Pt(1, 2), Sz(2, 3)).Contains(Pt(1, 1)))
		assert.True(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Contains(Pt(0.25, -0.75)))
	})
	t.Run("outside", func(t *testing.T) {
		assert.False(t, Rect(Pt(1, 2), Sz(2, 3)).Contains(Pt(3, 0)))
		assert.False(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Contains(Pt(-0.1, 0)))
	})
	t.Run("the boundary counts as inside", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))

		assert.True(t, r.Contains(r.Min()))
		assert.True(t, r.Contains(r.Max()))
	})
	t.Run("a float rectangle contains the corners it was built from", func(t *testing.T) {
		for _, corners := range [][2]float64{{0.1, 0.7}, {0.2, 0.9}, {5.7, 9.1}, {0.15, 0.45}} {
			low, high := Pt(corners[0], corners[0]), Pt(corners[1], corners[1])
			r := RectangleFromMinMax(low, high)

			assert.True(t, r.Contains(low), fmt.Sprintf("%v: min", corners))
			assert.True(t, r.Contains(high), fmt.Sprintf("%v: max", corners))
		}
	})
	t.Run("float beyond Delta is outside", func(t *testing.T) {
		r := RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, r.Contains(Pt(1+Delta/2, 0.5)))
		assert.False(t, r.Contains(Pt(1+2*Delta, 0.5)))
	})
}

func TestRectangle_Intersects(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, rectangle.Intersects(Rect(Pt(100.0, -50.0), Sz(200.0, 50.0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, rectangle.Intersects(Rect(Pt(100.0, 350.0), Sz(200.0, 450.0))))
	})
	t.Run("a shared edge counts as an intersection", func(t *testing.T) {
		assert.True(t, rectangle.Intersects(Rect(Pt(200.0, 0.0), Sz(200.0, 100.0))))
		assert.False(t, rectangle.Intersects(Rect(Pt(201.0, 0.0), Sz(200.0, 100.0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, rectangle.Intersects(Rect(Pt(0.0, 0.0), Sz(50.0, 50.0))))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				assert.Equal(t, a.Intersects(b), b.Intersects(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRectangle_IntersectsCircle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(150.0, 0.0), 60.0)))
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(110.0, 80.0), 60.0)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, rectangle.IntersectsCircle(Circ(Pt(150.0, 0.0), 40.0)))
	})
	t.Run("the circle center on the boundary counts", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(100.0, 0.0), 1.0)))
	})
	t.Run("touching the edge from outside counts", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(200.0, 0.0), 100.0)))
		assert.False(t, rectangle.IntersectsCircle(Circ(Pt(201.0, 0.0), 100.0)))
	})
	t.Run("touching the corner from outside counts", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(103.0, 54.0), 5.0)))
		assert.False(t, rectangle.IntersectsCircle(Circ(Pt(104.0, 54.0), 5.0)))
	})
	t.Run("odd integer sizes keep exact half extents", func(t *testing.T) {
		odd := Rect(Pt(0, 0), Sz(3, 3))

		assert.True(t, odd.IntersectsCircle(Circ(Pt(3, 0), 1)))
		assert.False(t, odd.IntersectsCircle(Circ(Pt(4, 0), 1)))
	})
	t.Run("circle fully inside the rectangle", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsCircle(Circ(Pt(0.0, 0.0), 10.0)))
	})
	t.Run("a circle intersects its own bounds", func(t *testing.T) {
		for _, c := range circleFixtures {
			if c.Radius == 0 {
				continue // a degenerate circle touches nothing
			}

			assert.True(t, c.Bounds().IntersectsCircle(c), fmt.Sprintf("%s: ", c))
		}
	})
}

func TestRectangle_IntersectsLine(t *testing.T) {
	t.Run("mirrors Line.IntersectsRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, l := range lineFixtures {
				assert.Equal(t, r.IntersectsLine(l), l.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", r, l))
			}
		}
	})
}

func TestRectangle_IntersectsPolygon(t *testing.T) {
	t.Run("mirrors Polygon.IntersectsRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, r.IntersectsPolygon(p), p.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRectangle_Polygon(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	p := r.Polygon()

	t.Run("carries the vertices", func(t *testing.T) {
		AssertVertices(t, p.Vertices, r.Vertices())
	})
	t.Run("owns its slice", func(t *testing.T) {
		assert.NotSame(t, p.Vertices, r.Vertices())
	})
}

func TestRectangle_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Int(), Rect(Pt(1, 2), Sz(2, 3)))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Int(), Rect(Pt(1, 0), Sz(1, 4)))
	})
}

func TestRectangle_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertRect(t, Rect(Pt(1, 2), Sz(2, 3)).Float(), Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertRect(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Float(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
}

func TestRectangle_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).String(), "Rect((1,2);2x3)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).String(), "Rect((0.60,-0.25);1.20x3.60)")
	})
}

func TestRectangle_MinMaxString(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).MinMaxString(), "(0,1)-(2,4)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).MinMaxString(), "(0.00,-2.05)-(1.20,1.55)")
	})
}

func TestRectangle_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Rect(Pt(1, 2), Sz(2, 3)), `{"x":1,"y":2,"w":2,"h":3}`)

		var r Rectangle[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":3}`), &r))
		AssertRect(t, r, Rect(Pt(1, 2), Sz(2, 3)))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)), `{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`)

		var r Rectangle[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`), &r))
		AssertRect(t, r, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, rectangle := range rectFixtures {
			data, err := json.Marshal(rectangle)
			assert.NoError(t, err)

			var decoded Rectangle[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, rectangle)
		}
	})
}

func TestRectangle_Properties(t *testing.T) {
	t.Run("min and max bracket the center", func(t *testing.T) {
		for _, r := range rectFixtures {
			assert.True(t, r.Min().Midpoint(r.Max()).Equal(r.Center), fmt.Sprintf("%s: ", r))
			AssertNumber(t, r.Max().X-r.Min().X, r.Width(), fmt.Sprintf("%s: ", r))
			AssertNumber(t, r.Max().Y-r.Min().Y, r.Height(), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("area and perimeter follow the size", func(t *testing.T) {
		for _, r := range rectFixtures {
			AssertNumber(t, r.Area(), r.Size.Area(), fmt.Sprintf("%s: ", r))
			AssertNumber(t, r.Perimeter(), r.Size.Perimeter(), fmt.Sprintf("%s: ", r))
			AssertNumber(t, r.AspectRatio(), r.Size.AspectRatio(), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("bounds is the rectangle itself", func(t *testing.T) {
		for _, r := range rectFixtures {
			assert.True(t, r.Bounds().Equal(r), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("translate keeps the size", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, vector := range vectorFixtures {
				moved := r.Translate(vector)

				assert.True(t, moved.Size.Equal(r.Size), fmt.Sprintf("%s → %s: ", r, vector))
				assert.True(t, moved.Center.Equal(r.Center.Add(vector)), fmt.Sprintf("%s → %s: ", r, vector))
			}
		}
	})
	t.Run("move to keeps the size and centers where asked", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, point := range pointFixtures {
				moved := r.MoveTo(point)

				assert.True(t, moved.Size.Equal(r.Size), fmt.Sprintf("%s → %s: ", r, point))
				assert.True(t, moved.Center.Equal(point), fmt.Sprintf("%s → %s: ", r, point))
			}
		}
	})
	t.Run("grow and shrink are inverse above zero", func(t *testing.T) {
		for _, r := range rectFixtures {
			if r.Width() < 1 || r.Height() < 1 {
				continue // shrink clamps to zero, so the growth is not recoverable
			}

			assert.True(t, r.Grow(1).Shrink(1).Equal(r), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("clamp lands inside and fixes inner points", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, point := range pointFixtures {
				clamped := r.Clamp(point)

				assert.True(t, clamped.X >= r.Min().X-Delta && clamped.X <= r.Max().X+Delta, fmt.Sprintf("%s → %s: ", r, point))
				assert.True(t, clamped.Y >= r.Min().Y-Delta && clamped.Y <= r.Max().Y+Delta, fmt.Sprintf("%s → %s: ", r, point))
				assert.True(t, r.Clamp(clamped).Equal(clamped), fmt.Sprintf("%s → %s: ", r, point))
			}
		}
	})
	t.Run("vertices and edges describe the same outline", func(t *testing.T) {
		for _, r := range rectFixtures {
			vertices, edges := r.Vertices(), r.Edges()

			assert.Equal(t, len(vertices), 4, fmt.Sprintf("%s: ", r))
			assert.Equal(t, len(edges), 4, fmt.Sprintf("%s: ", r))

			for i, edge := range edges {
				assert.True(t, edge.Start.Equal(vertices[i]), fmt.Sprintf("%s: ", r))
				assert.True(t, edge.End.Equal(vertices[(i+1)%4]), fmt.Sprintf("%s: ", r))
			}
		}
	})
	t.Run("anchors lie on the rectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, direction := range Directions() {
				anchor := r.Anchor(direction)

				assert.True(t, anchor.X >= r.Min().X-Delta && anchor.X <= r.Max().X+Delta, fmt.Sprintf("%s → %s: ", r, direction))
				assert.True(t, anchor.Y >= r.Min().Y-Delta && anchor.Y <= r.Max().Y+Delta, fmt.Sprintf("%s → %s: ", r, direction))
			}
		}
	})
	t.Run("polygon carries the vertices", func(t *testing.T) {
		for _, r := range rectFixtures {
			AssertVertices(t, r.Polygon().Vertices, r.Vertices())
		}
	})
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
	r.Inset(PadU(1))

	AssertRect(t, r, Rect(Pt(1, 2), Sz(2, 3)))
}

// rectFixtures span square, portrait, landscape, degenerate, and off-origin rectangles.
var rectFixtures = []Rectangle[float64]{
	Rect(Pt(0.0, 0.0), Sz(0.0, 0.0)),
	Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)),
	Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)),
	Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)),
	Rect(Pt(-3.5, 0.25), Sz(7.0, 1.0)),
	RectangleFromMin(Pt(0.0, 0.0), Sz(10.0, 20.0)),
	RectangleFromMinMax(Pt(-2.0, -4.0), Pt(6.0, 2.0)),
}

func ExampleRect() {
	fmt.Println(Rect(Pt(1, 2), Sz(2, 3)))
	// Output: Rect((1,2);2x3)
}

func ExampleRectangle_MinMaxString() {
	fmt.Println(Rect(Pt(1, 2), Sz(2, 3)).MinMaxString())
	// Output: (0,1)-(2,4)
}

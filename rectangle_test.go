package geom_test

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestRectangle_Constructor(t *testing.T) {
	t.Run("from center and size", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(10, 16), Sz(3, 4)), Rectangle[int]{Center: Pt(10, 16), Size: Sz(3, 4)})
		geomtest.AssertRectangle(t, Rect(Pt(0.5, -1.25), Sz(2.5, 3.75)), Rectangle[float64]{Center: Pt(0.5, -1.25), Size: Sz(2.5, 3.75)})
	})
	t.Run("from min", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMin(Pt(0, 0), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromMin(Pt(0, 0), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
		geomtest.AssertRectangle(t, RectangleFromMin(Pt(0.0, 0.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from max", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMax(Pt(4, 2), Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromMax(Pt(5, 3), Sz(5, 3)), Rect(Pt(2, 1), Sz(5, 3)))
		geomtest.AssertRectangle(t, RectangleFromMax(Pt(1.0, 3.0), Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from min and max", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMinMax(Pt(0, 0), Pt(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("from size at the origin", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromSize(Sz(4, 2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromSize(Sz(1.0, 3.0)), Rect(Pt(0.5, 1.5), Sz(1.0, 3.0)))
	})
	t.Run("unit integer rectangle", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMin(Pt(0, 0), Sz(1, 1)), Rect(Pt(0, 0), Sz(1, 1)))
	})
	t.Run("a negative size is taken absolute", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(2, 1), Sz(-4, -2)), Rect(Pt(2, 1), Sz(4, 2)))
	})
	t.Run("min and max round-trip for odd and even integer sizes", func(t *testing.T) {
		for w := 0; w < 6; w++ {
			for h := 0; h < 6; h++ {
				r := Rect(Pt(-3, 7), Sz(w, h))

				geomtest.AssertRectangle(t, RectangleFromMin(r.Min(), r.Size), r, r.String())
				geomtest.AssertRectangle(t, RectangleFromMax(r.Max(), r.Size), r, r.String())
				geomtest.AssertRectangle(t, RectangleFromMinMax(r.MinMax()), r, r.String())
			}
		}
	})
	t.Run("a negative extent measures the other way from the corner", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMin(Pt(4, 2), Sz(-4, -2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromMax(Pt(0, 0), Sz(-4, -2)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromSize(Sz(-4, 2)), RectangleFromMinMax(Pt(-4, 0), Pt(0, 2)))
	})
	t.Run("corners in either order", func(t *testing.T) {
		geomtest.AssertRectangle(t, RectangleFromMinMax(Pt(4, 2), Pt(0, 0)), Rect(Pt(2, 1), Sz(4, 2)))
		geomtest.AssertRectangle(t, RectangleFromMinMax(Pt(4, 0), Pt(0, 2)), Rect(Pt(2, 1), Sz(4, 2)))
	})
	t.Run("every constructor builds a rectangle that is not rotated", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(2, 1), Sz(4, 2)).Angle, 0.0)
		assert.Equal(t, RectangleFromMinMax(Pt(0, 0), Pt(4, 2)).Angle, 0.0)
		assert.Equal(t, Rectangle[int]{Pt(2, 1), Sz(4, 2), Pi}.Angle, Pi)
	})
}

func TestRectangle_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Width(), 2)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Width(), 1.2)
	})
}

func TestRectangle_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Height(), 3)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Height(), 3.6)
	})
}

func TestRectangle_Min(t *testing.T) {
	t.Run("top-left corner", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Min(), Pt(0, 1))
		geomtest.AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Min(), Pt(0.0, -2.05))
	})
	t.Run("agrees with the corner accessor", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))
		geomtest.AssertPoint(t, r.Min(), r.TopLeft())
	})
	t.Run("rotated is the minimum of the vertices", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi/2).Min(), Pt(-1, -2))
		geomtest.AssertPoint(t, Rect(Pt(1.0, 2.0), Sz(2.0, 4.0)).Rotate(Pi/2).Min(), Pt(-1.0, 1.0))
	})
}

func TestRectangle_Max(t *testing.T) {
	t.Run("bottom-right corner", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Max(), Pt(2, 4))
		geomtest.AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Max(), Pt(1.2, 1.55))
	})
	t.Run("agrees with the corner accessor", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))
		geomtest.AssertPoint(t, r.Max(), r.BottomRight())
	})
	t.Run("rotated is the maximum of the vertices", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi/2).Max(), Pt(1, 2))
		geomtest.AssertPoint(t, Rect(Pt(1.0, 2.0), Sz(2.0, 4.0)).Rotate(Pi/2).Max(), Pt(3.0, 3.0))
	})
}

func TestRectangle_MinMax(t *testing.T) {
	t.Run("both corners at once", func(t *testing.T) {
		for _, r := range rectFixtures {
			a, b := r.MinMax()

			geomtest.AssertPoint(t, a, r.Min(), fmt.Sprintf("%s min: ", r))
			geomtest.AssertPoint(t, b, r.Max(), fmt.Sprintf("%s max: ", r))
		}
	})
}

func TestRectangle_Corners(t *testing.T) {
	r := RectangleFromMin(Pt(0, 0), Sz(10, 20))

	t.Run("clockwise from the minimum", func(t *testing.T) {
		geomtest.AssertPoint(t, r.TopLeft(), Pt(0, 0))
		geomtest.AssertPoint(t, r.TopRight(), Pt(10, 0))
		geomtest.AssertPoint(t, r.BottomRight(), Pt(10, 20))
		geomtest.AssertPoint(t, r.BottomLeft(), Pt(0, 20))
	})
	t.Run("the diagonal corners are min and max", func(t *testing.T) {
		geomtest.AssertPoint(t, r.TopLeft(), r.Min())
		geomtest.AssertPoint(t, r.BottomRight(), r.Max())
	})
	t.Run("float", func(t *testing.T) {
		f := RectangleFromMin(Pt(0.5, 1.25), Sz(2.0, 3.0))

		geomtest.AssertPoint(t, f.TopLeft(), Pt(0.5, 1.25))
		geomtest.AssertPoint(t, f.TopRight(), Pt(2.5, 1.25))
		geomtest.AssertPoint(t, f.BottomRight(), Pt(2.5, 4.25))
		geomtest.AssertPoint(t, f.BottomLeft(), Pt(0.5, 4.25))
	})
	t.Run("rotated corners keep their names and turn about the center", func(t *testing.T) {
		turned := Rect(Pt(1.0, 2.0), Sz(2.0, 4.0)).Rotate(Pi / 2)

		geomtest.AssertPoint(t, turned.TopLeft(), Pt(3.0, 1.0))
		geomtest.AssertPoint(t, turned.TopRight(), Pt(3.0, 3.0))
		geomtest.AssertPoint(t, turned.BottomRight(), Pt(-1.0, 3.0))
		geomtest.AssertPoint(t, turned.BottomLeft(), Pt(-1.0, 1.0))
	})
	t.Run("rotated integer corners are rounded onto the lattice", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		geomtest.AssertPoint(t, turned.TopLeft(), Pt(1, -2))
		geomtest.AssertPoint(t, turned.TopRight(), Pt(1, 2))
		geomtest.AssertPoint(t, turned.BottomRight(), Pt(-1, 2))
		geomtest.AssertPoint(t, turned.BottomLeft(), Pt(-1, -2))
	})
}

func TestRectangle_Anchor(t *testing.T) {
	r := RectangleFromMin(Pt(0, 0), Sz(10, 20))

	t.Run("corners", func(t *testing.T) {
		geomtest.AssertPoint(t, r.Anchor(TopLeft), Pt(0, 0))
		geomtest.AssertPoint(t, r.Anchor(TopRight), Pt(10, 0))
		geomtest.AssertPoint(t, r.Anchor(BottomLeft), Pt(0, 20))
		geomtest.AssertPoint(t, r.Anchor(BottomRight), Pt(10, 20))
	})
	t.Run("edge midpoints halve the edge they sit on", func(t *testing.T) {
		geomtest.AssertPoint(t, r.Anchor(Top), Pt(5, 0))
		geomtest.AssertPoint(t, r.Anchor(Bottom), Pt(5, 20))
		geomtest.AssertPoint(t, r.Anchor(DirectionLeft), Pt(0, 10))
		geomtest.AssertPoint(t, r.Anchor(DirectionRight), Pt(10, 10))
	})
	t.Run("edge midpoints share a coordinate with the center", func(t *testing.T) {
		geomtest.AssertNumber(t, r.Anchor(Top).X, r.Center.X)
		geomtest.AssertNumber(t, r.Anchor(Bottom).X, r.Center.X)
		geomtest.AssertNumber(t, r.Anchor(DirectionLeft).Y, r.Center.Y)
		geomtest.AssertNumber(t, r.Anchor(DirectionRight).Y, r.Center.Y)
	})
	t.Run("float", func(t *testing.T) {
		f := RectangleFromMin(Pt(0.5, 1.25), Sz(2.0, 3.0))

		geomtest.AssertPoint(t, f.Anchor(Top), Pt(1.5, 1.25))
		geomtest.AssertPoint(t, f.Anchor(DirectionRight), Pt(2.5, 2.75))
		geomtest.AssertPoint(t, f.Anchor(Bottom), Pt(1.5, 4.25))
		geomtest.AssertPoint(t, f.Anchor(DirectionLeft), Pt(0.5, 2.75))
	})
	t.Run("none is the center", func(t *testing.T) {
		geomtest.AssertPoint(t, r.Anchor(DirectionNone), Pt(5, 10))
		geomtest.AssertPoint(t, r.Anchor(DirectionNone), r.Center)
	})
	t.Run("out-of-range direction wraps", func(t *testing.T) {
		geomtest.AssertPoint(t, r.Anchor(Direction(99)), r.BottomLeft())
	})
	t.Run("agrees with the dedicated accessors", func(t *testing.T) {
		geomtest.AssertPoint(t, r.Anchor(TopLeft), r.TopLeft())
		geomtest.AssertPoint(t, r.Anchor(TopRight), r.TopRight())
		geomtest.AssertPoint(t, r.Anchor(BottomLeft), r.BottomLeft())
		geomtest.AssertPoint(t, r.Anchor(BottomRight), r.BottomRight())
	})
	t.Run("rotated anchors are named before the turn", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		geomtest.AssertPoint(t, turned.Anchor(Top), Pt(1, 0))
		geomtest.AssertPoint(t, turned.Anchor(DirectionRight), Pt(0, 2))
		geomtest.AssertPoint(t, turned.Anchor(Bottom), Pt(-1, 0))
		geomtest.AssertPoint(t, turned.Anchor(DirectionLeft), Pt(0, -2))
		geomtest.AssertPoint(t, turned.Anchor(TopLeft), Pt(1, -2))
		geomtest.AssertPoint(t, turned.Anchor(DirectionNone), Pt(0, 0))
	})
	t.Run("odd integer extents split like min and max", func(t *testing.T) {
		odd := RectangleFromMin(Pt(0, 0), Sz(3, 3))

		geomtest.AssertPoint(t, odd.Anchor(TopLeft), Pt(0, 0))
		geomtest.AssertPoint(t, odd.Anchor(BottomRight), Pt(3, 3))
		geomtest.AssertPoint(t, odd.Anchor(Top), Pt(1, 0))
		geomtest.AssertPoint(t, odd.Anchor(Bottom), Pt(1, 3))
		geomtest.AssertPoint(t, odd.Anchor(DirectionLeft), Pt(0, 1))
		geomtest.AssertPoint(t, odd.Anchor(DirectionRight), Pt(3, 1))
	})
}

func TestRectangle_Vertices(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	vertices := slices.Collect(r.Vertices())

	t.Run("clockwise from the top left", func(t *testing.T) {
		geomtest.AssertVertices(t, vertices, []Point[int]{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}})
	})
	t.Run("agrees with the corner accessors", func(t *testing.T) {
		geomtest.AssertPoint(t, vertices[0], r.TopLeft())
		geomtest.AssertPoint(t, vertices[1], r.TopRight())
		geomtest.AssertPoint(t, vertices[2], r.BottomRight())
		geomtest.AssertPoint(t, vertices[3], r.BottomLeft())
	})
	t.Run("agrees with the edge starts", func(t *testing.T) {
		edges := slices.Collect(r.Edges())

		for i, edge := range edges {
			geomtest.AssertPoint(t, edge.Start, vertices[i])
		}
	})
	t.Run("rotated vertices are the corners turned about the center", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertVertices(t, slices.Collect(diamond.Vertices()), []Point[float64]{{0, -Sqrt2}, {Sqrt2, 0}, {0, Sqrt2}, {-Sqrt2, 0}})
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for vertex := range r.Vertices() {
			geomtest.AssertPoint(t, vertex, r.TopLeft())

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
				for vertex := range r.Vertices() {
					sinkBool = vertex.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", r))
		}
	})
}

func TestRectangle_Edges(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	edges := slices.Collect(r.Edges())

	t.Run("clockwise from the top", func(t *testing.T) {
		assert.Length(t, edges, 4)
		geomtest.AssertSegment(t, edges[0], Seg(Pt(-1, -1), Pt(1, -1)))
		geomtest.AssertSegment(t, edges[1], Seg(Pt(1, -1), Pt(1, 1)))
		geomtest.AssertSegment(t, edges[2], Seg(Pt(1, 1), Pt(-1, 1)))
		geomtest.AssertSegment(t, edges[3], Seg(Pt(-1, 1), Pt(-1, -1)))
	})
	t.Run("run clockwise from their first corner", func(t *testing.T) {
		box := RectangleFromMin(Pt(0, 0), Sz(10, 20))
		boxEdges := slices.Collect(box.Edges())

		geomtest.AssertSegment(t, boxEdges[0], Seg(Pt(0, 0), Pt(10, 0)))
		geomtest.AssertSegment(t, boxEdges[1], Seg(Pt(10, 0), Pt(10, 20)))
		geomtest.AssertSegment(t, boxEdges[2], Seg(Pt(10, 20), Pt(0, 20)))
		geomtest.AssertSegment(t, boxEdges[3], Seg(Pt(0, 20), Pt(0, 0)))
	})
	t.Run("their midpoints are the cardinal anchors", func(t *testing.T) {
		geomtest.AssertPoint(t, edges[0].Midpoint(), r.Anchor(Top))
		geomtest.AssertPoint(t, edges[1].Midpoint(), r.Anchor(DirectionRight))
		geomtest.AssertPoint(t, edges[2].Midpoint(), r.Anchor(Bottom))
		geomtest.AssertPoint(t, edges[3].Midpoint(), r.Anchor(DirectionLeft))
	})
	t.Run("their lengths are the extents", func(t *testing.T) {
		box := RectangleFromMin(Pt(0, 0), Sz(10, 20))
		boxEdges := slices.Collect(box.Edges())

		geomtest.AssertNumber(t, boxEdges[0].Length(), float64(box.Width()))
		geomtest.AssertNumber(t, boxEdges[2].Length(), float64(box.Width()))
		geomtest.AssertNumber(t, boxEdges[1].Length(), float64(box.Height()))
		geomtest.AssertNumber(t, boxEdges[3].Length(), float64(box.Height()))
	})
	t.Run("form a closed chain", func(t *testing.T) {
		for i, edge := range edges {
			geomtest.AssertPoint(t, edge.Start, edges[(i+3)%4].End)
		}
	})
	t.Run("rotated edges join the turned corners", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		geomtest.AssertSegment(t, slices.Collect(turned.Edges())[0], Seg(Pt(1, -2), Pt(1, 2)))
		geomtest.AssertSegment(t, slices.Collect(turned.Edges())[2], Seg(Pt(-1, 2), Pt(-1, -2)))
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for edge := range r.Edges() {
			geomtest.AssertSegment(t, edge, edges[0])

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
				for edge := range r.Edges() {
					sinkBool = edge.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", r))
		}
	})
}

func TestRectangle_Centroid(t *testing.T) {
	geomtest.AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Rotate(1).Centroid(), Pt(1, 2))
}

func TestRectangle_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Area(), 6.0)
	})
	t.Run("large integer sides do not overflow", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0, 0), SzU(3037000500)).Area(), 3037000500.0*3037000500)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Area(), 4.32)
	})
}

func TestRectangle_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Perimeter(), 10)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Perimeter(), 9.6)
	})
}

func TestRectangle_Inertia(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Inertia(), 6.5)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Inertia(), 4.32*14.4/12)
	})
	t.Run("the turn does not change it", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).Rotate(1).Inertia(), 6.5)
	})
	t.Run("agrees with the polygon at any angle", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertNumber(t, r.Inertia(), r.Polygon().Inertia(), fmt.Sprintf("%s: ", r))
		}
	})
}

func TestRectangle_AspectRatio(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(1, 2), Sz(2, 3)).AspectRatio(), 2.0/3.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).AspectRatio(), 1.0/3.0)
	})
}

func TestRectangle_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertBox(t, Rect(Pt(1, 2), Sz(2, 3)).Bounds(), BoxFromMinMax(Pt(0, 1), Pt(2, 4)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Bounds(), BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)))
	})
	t.Run("rotated is the box around the vertices, not rotated", func(t *testing.T) {
		geomtest.AssertBox(t, Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi/2).Bounds(), BoxFromMinMax(Pt(-1, -2), Pt(1, 2)))
		geomtest.AssertBox(t, Rect(Pt(1.0, 2.0), Sz(2.0, 4.0)).Rotate(Pi/2).Bounds(), BoxFromMinMax(Pt(-1.0, 1.0), Pt(3.0, 3.0)))
		geomtest.AssertBox(t, Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi/4).Bounds(), BoxFromMinMax(Pt(-Sqrt2, -Sqrt2), Pt(Sqrt2, Sqrt2)))
	})
}

func TestRectangle_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Translate(Vec(3, -2)), Rect(Pt(4, 0), Sz(2, 3)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Translate(Vec(100.1, -0.1)), Rect(Pt(100.7, -0.35), Sz(1.2, 3.6)))
	})
}

func TestRectangle_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).MoveTo(Pt(3, -2)), Rect(Pt(3, -2), Sz(2, 3)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).MoveTo(Pt(100.1, -0.1)), Rect(Pt(100.1, -0.1), Sz(1.2, 3.6)))
	})
}

func TestRectangle_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Scale(2.5), Rect(Pt(1, 2), Sz(5, 8)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Scale(2.5), Rect(Pt(0.6, -0.25), Sz(3.0, 9.0)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).ScaleXY(2, 3), Rect(Pt(1, 2), Sz(4, 9)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ScaleXY(-1.5, 2), Rect(Pt(0.6, -0.25), Sz(1.8, 7.2)))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Scale(-2), Rect(Pt(1, 2), Sz(4, 6)))
	})
}

func TestRectangle_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(10, 20)).Unscale(2.5), Rect(Pt(1, 2), Sz(4, 8)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(3.0, 9.0)).Unscale(2.5), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(10, 20)).UnscaleXY(2, 4), Rect(Pt(1, 2), Sz(5, 5)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.8, 7.2)).UnscaleXY(-1.5, 2), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Rect(Pt(1, 2), Sz(10, 20)).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Rect(Pt(1, 2), Sz(10, 20)).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestRectangle_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Resize(Sz(8, 9)), Rect(Pt(1, 2), Sz(8, 9)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Resize(Sz(3.1, 0.2)), Rect(Pt(0.6, -0.25), Sz(3.1, 0.2)))
	})
	t.Run("a negative size is taken absolute", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Resize(Sz(-8, 9)), Rect(Pt(1, 2), Sz(8, 9)))
	})
}

func TestRectangle_Canonical(t *testing.T) {
	t.Run("takes a literal negative size absolute and keeps the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rectangle[int]{Pt(1, 2), Sz(-8, 9), 0}.Canonical(), Rect(Pt(1, 2), Sz(8, 9)))
		geomtest.AssertRectangle(t, Rectangle[float64]{Pt(0.6, -0.25), Sz(-1.2, -3.6), 0}.Canonical(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
	t.Run("is the rectangle Rect builds", func(t *testing.T) {
		r := Rectangle[int]{Pt(1, 2), Sz(-8, 9), 0}

		geomtest.AssertRectangle(t, r.Canonical(), Rect(r.Center, r.Size))
		assert.True(t, r.Canonical().Contains(r.Canonical().Min()))
	})
	t.Run("repairs decoded JSON", func(t *testing.T) {
		var r Rectangle[int]

		assert.Nil(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":-8,"h":9}`), &r))
		geomtest.AssertRectangle(t, r.Canonical(), Rect(Pt(1, 2), Sz(8, 9)))
	})
	t.Run("normalizes the angle", func(t *testing.T) {
		r := Rectangle[float64]{Pt(1.0, 2.0), Sz(2.0, 3.0), 5 * Pi}.Canonical()

		geomtest.AssertRectangle(t, r, Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)).Rotate(Pi))
		assert.Equal(t, r.Angle, Pi)
	})
	t.Run("snaps a residue of turning back to exactly zero, where Rotate does not", func(t *testing.T) {
		drifted := Rect(Pt(1, 2), Sz(2, 3)).Rotate(0.1).Rotate(0.2).Rotate(-0.3)

		assert.False(t, drifted.IsAligned())
		assert.True(t, drifted.Canonical().IsAligned())
		assert.True(t, Rectangle[int]{Pt(1, 2), Sz(2, 3), 2*Pi - Delta/2}.Canonical().IsAligned())
		assert.False(t, Rectangle[int]{Pt(1, 2), Sz(2, 3), 2 * Delta}.Canonical().IsAligned())
	})
	t.Run("a well-formed rectangle is unchanged", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertRectangle(t, r.Canonical(), r, r.String())
		}
	})
}

func TestRectangle_Grow(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Grow(2), Rect(Pt(1, 2), Sz(4, 5)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Grow(0.1), Rect(Pt(0.6, -0.25), Sz(1.3, 3.7)))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).GrowXY(2, 3), Rect(Pt(1, 2), Sz(4, 6)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).GrowXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.3, 3.8)))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Grow(-5), Rect(Pt(1, 2), Sz(0, 0)))
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).GrowXY(-5, 1), Rect(Pt(1, 2), Sz(0, 4)))
	})
}

func TestRectangle_Shrink(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(1), Rect(Pt(1, 2), Sz(1, 2)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Shrink(0.1), Rect(Pt(0.6, -0.25), Sz(1.1, 3.5)))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).ShrinkXY(1, 2), Rect(Pt(1, 2), Sz(1, 1)))
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).ShrinkXY(0.1, 0.2), Rect(Pt(0.6, -0.25), Sz(1.1, 3.4)))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Shrink(100), Rect(Pt(1, 2), Sz(0, 0)))
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).ShrinkXY(5, 1), Rect(Pt(1, 2), Sz(0, 2)))
	})
}

func TestRectangle_Inset(t *testing.T) {
	t.Run("uniform padding keeps the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(PadU(1)), Rect(Pt(0, 0), Sz(8, 8)))
	})
	t.Run("asymmetric padding shifts the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(3, 1, 1, 5)), Rect(Pt(2, 1), Sz(4, 6)))
	})
	t.Run("odd integer padding stays inside the rectangle", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(10, 10))

		geomtest.AssertPoint(t, r.Inset(Pad(0, 0, 0, 1)).Min(), r.Min().AddXY(1, 0))
		geomtest.AssertPoint(t, r.Inset(Pad(0, 0, 0, 1)).Max(), r.Max())
		geomtest.AssertPoint(t, r.Inset(Pad(0, 3, 1, 0)).Min(), r.Min())
		geomtest.AssertPoint(t, r.Inset(Pad(0, 3, 1, 0)).Max(), r.Max().AddXY(-3, -1))
	})
	t.Run("padding beyond the size collapses inside the rectangle", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(10, 10))

		geomtest.AssertRectangle(t, r.Inset(Pad(0, 0, 0, 20)), Rect(Pt(5, 0), Sz(0, 10)))
		geomtest.AssertRectangle(t, r.Inset(Pad(0, 20, 0, 0)), Rect(Pt(-5, 0), Sz(0, 10)))
		geomtest.AssertRectangle(t, r.Inset(Pad(20, 0, 20, 0)), Rect(Pt(0, 5), Sz(10, 0)))
		geomtest.AssertRectangle(t, r.Inset(PadU(20)), Rect(Pt(5, 5), Sz(0, 0)))
		geomtest.AssertRectangle(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(0.0, 0.0, 0.0, 12.5)), Rect(Pt(5.0, 0.0), Sz(0.0, 10.0)))
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
		geomtest.AssertRectangle(t, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), Rect(Pt(1.5, 0.75), Sz(11.0, 8.5)))
	})
	t.Run("each edge moves inward by its own padding", func(t *testing.T) {
		r := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
		padding := Pad(3.0, 1.0, 1.0, 5.0)
		inset := r.Inset(padding)

		geomtest.AssertPoint(t, inset.Min(), r.Min().AddXY(padding.Left, padding.Top))
		geomtest.AssertPoint(t, inset.Max(), r.Max().AddXY(-padding.Right, -padding.Bottom))
	})
	t.Run("rotated moves the edges in the frame before the turn", func(t *testing.T) {
		turned := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Rotate(Pi / 2)
		inset := turned.Inset(Pad(0.0, 0.0, 0.0, 4.0))

		geomtest.AssertRectangle(t, inset, Rect(Pt(0.0, 2.0), Sz(6.0, 10.0)).Rotate(Pi/2))
		geomtest.AssertPoint(t, inset.Anchor(DirectionLeft), Pt(0.0, -1.0))
		geomtest.AssertPoint(t, inset.Anchor(DirectionRight), turned.Anchor(DirectionRight))
	})
}

func TestRectangle_Outset(t *testing.T) {
	t.Run("uniform padding keeps the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0, 0), Sz(10, 10)).Outset(PadU(1)), Rect(Pt(0, 0), Sz(12, 12)))
	})
	t.Run("asymmetric padding shifts the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0, 0), Sz(10, 10)).Outset(Pad(3, 1, 1, 5)), Rect(Pt(-2, -1), Sz(16, 14)))
	})
	t.Run("undoes an inset", func(t *testing.T) {
		r := RectangleFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		padding := Pad(1.0, 2.0, 3.0, 4.0)

		geomtest.AssertRectangle(t, r.Inset(padding).Outset(padding), r)
		geomtest.AssertRectangle(t, r.Outset(padding).Inset(padding), r)
	})
}

func TestRectangle_Lerp(t *testing.T) {
	a := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
	b := Rect(Pt(10.0, 20.0), Sz(20.0, 30.0))

	t.Run("moves the center and the size together", func(t *testing.T) {
		geomtest.AssertRectangle(t, a.Lerp(b, 0.5), Rect(Pt(5.0, 10.0), Sz(15.0, 20.0)))
	})
	t.Run("the ends are the rectangles themselves", func(t *testing.T) {
		geomtest.AssertRectangle(t, a.Lerp(b, 0), a)
		geomtest.AssertRectangle(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates outside the unit range", func(t *testing.T) {
		geomtest.AssertRectangle(t, a.Lerp(b, 2), Rect(Pt(20.0, 40.0), Sz(30.0, 50.0)))
	})
	t.Run("the size stays absolute", func(t *testing.T) {
		shrinking := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)).Lerp(Rect(Pt(0.0, 0.0), Sz(0.0, 0.0)), 2)

		geomtest.AssertRectangle(t, shrinking, Rect(Pt(0.0, 0.0), Sz(10.0, 10.0)))
	})
	t.Run("int rounds like every other interpolation", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0, 0), Sz(0, 10)).Lerp(Rect(Pt(1, 1), Sz(1, 11)), 0.5), Rect(Pt(1, 1), Sz(1, 11)))
	})
	t.Run("turns the angle along the shorter arc", func(t *testing.T) {
		geomtest.AssertRectangle(t, a.Lerp(b.Rotate(Pi/2), 0.5), Rect(Pt(5.0, 10.0), Sz(15.0, 20.0)).Rotate(Pi/4))
		geomtest.AssertRectangle(t, a.Rotate(-Pi/4).Lerp(a.Rotate(Pi/4), 0.5), a)
	})
}

func TestRectangle_Transform(t *testing.T) {
	rectangle := Rect(Pt(2.0, 3.0), Sz(4.0, 2.0))

	t.Run("the identity keeps the rectangle", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(IdentityMatrix[float64]()), rectangle)
	})
	t.Run("a translation moves the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(TranslationMatrix(1.0, -1.0)), Rect(Pt(3.0, 2.0), Sz(4.0, 2.0)))
	})
	t.Run("a uniform scale scales the size", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(ScaleMatrix(2.0, 2.0)), Rect(Pt(4.0, 6.0), Sz(8.0, 4.0)))
	})
	t.Run("axes scaled by different factors stay a rectangle while it is aligned", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(ScaleMatrix(2.0, 3.0)), Rect(Pt(4.0, 9.0), Sz(8.0, 6.0)))
	})
	t.Run("a rotation turns the angle", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(RotationMatrix[float64](Pi/2)), Rect(Pt(-3.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi/2))
	})
	t.Run("a reflection mirrors the angle about the axis of the matrix", func(t *testing.T) {
		turned := rectangle.Rotate(Pi / 6)

		geomtest.AssertRectangle(t, turned.Transform(ReflectionMatrix[float64](AxisHorizontal)), Rect(Pt(2.0, -3.0), Sz(4.0, 2.0)).Rotate(-Pi/6))
	})
	t.Run("a shear gives the nearest rectangle", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(ShearMatrix(1.0, 0.0)), Rect(Pt(5.0, 3.0), Sz(4.0, 2.0)))
	})
	t.Run("a matrix that collapses the plane gives the zero size", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Transform(ScaleMatrix(0.0, 0.0)), Rect(Pt(0.0, 0.0), Sz(0.0, 0.0)))
	})
	t.Run("int rounds the center and the size once", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 1), Sz(3, 3)).Transform(ScaleMatrix(1.5, 1.5)), Rect(Pt(2, 2), Sz(5, 5)))
	})
	t.Run("matches the polygon of the rectangle wherever the matrix keeps one", func(t *testing.T) {
		matrices := []Matrix[float64]{
			IdentityMatrix[float64](),
			TranslationMatrix(3.0, -2.0),
			ScaleMatrix(2.0, 2.0),
			RotationMatrix[float64](Pi / 3),
			RotationMatrix[float64](Pi / 3).Multiply(ScaleMatrix(2.0, 2.0)),
		}

		for _, r := range rectFixtures {
			for _, m := range matrices {
				geomtest.AssertPolygon(t, r.Transform(m).Polygon(), r.Polygon().Transform(m), fmt.Sprintf("%s → %s: ", r, m))
			}
		}
	})
}

func TestRectangle_Rotate(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(2, 3))

	t.Run("adds to the angle about the center", func(t *testing.T) {
		turned := r.Rotate(Pi / 2)

		assert.Equal(t, turned.Angle, Pi/2)
		geomtest.AssertPoint(t, turned.Center, r.Center)
		geomtest.AssertSize(t, turned.Size, r.Size)
		assert.Equal(t, turned.Rotate(Pi/2).Angle, Pi)
	})
	t.Run("normalizes into a turn", func(t *testing.T) {
		assert.Equal(t, r.Rotate(-Pi/2).Angle, 3*Pi/2)
		assert.Equal(t, r.Rotate(2*Pi).Angle, 0.0)
	})
	t.Run("a full turn is the rectangle itself, corners exact", func(t *testing.T) {
		assert.True(t, r.Rotate(2*Pi).IsAligned())
		geomtest.AssertVertices(t, slices.Collect(r.Rotate(2*Pi).Vertices()), slices.Collect(r.Vertices()))
	})
	t.Run("turning back undoes the turn", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertRectangle(t, r.Rotate(Pi/3).Rotate(-Pi/3), r, r.String())
		}
	})
}

func TestRectangle_AlignTo(t *testing.T) {
	r := RectangleFromMin(Pt(0, 0), Sz(10, 20))

	t.Run("a corner lands on the point", func(t *testing.T) {
		geomtest.AssertRectangle(t, r.AlignTo(TopLeft, Pt(100, 50)), RectangleFromMin(Pt(100, 50), Sz(10, 20)))
		geomtest.AssertRectangle(t, r.AlignTo(BottomRight, Pt(100, 50)), RectangleFromMax(Pt(100, 50), Sz(10, 20)))
	})
	t.Run("an edge midpoint lands on the point", func(t *testing.T) {
		geomtest.AssertPoint(t, r.AlignTo(Top, Pt(100, 50)).Anchor(Top), Pt(100, 50))
		geomtest.AssertPoint(t, r.AlignTo(DirectionLeft, Pt(100, 50)).Anchor(DirectionLeft), Pt(100, 50))
	})
	t.Run("none aligns the center", func(t *testing.T) {
		geomtest.AssertRectangle(t, r.AlignTo(DirectionNone, Pt(100, 50)), r.MoveTo(Pt(100, 50)))
	})
	t.Run("rotated lands the turned anchor on the point", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		geomtest.AssertRectangle(t, turned.AlignTo(TopLeft, Pt(10, 10)), Rect(Pt(9, 12), Sz(4, 2)).Rotate(Pi/2))
	})
	t.Run("keeps the size and inverts Anchor", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, direction := range Directions() {
				aligned := rect.AlignTo(direction, Pt(3.5, -2.25))

				geomtest.AssertSize(t, aligned.Size, rect.Size, fmt.Sprintf("%s %s: ", rect, direction))
				geomtest.AssertPoint(t, aligned.Anchor(direction), Pt(3.5, -2.25), fmt.Sprintf("%s %s: ", rect, direction))
			}
		}
	})
}

func TestRectangle_Clamp(t *testing.T) {
	box := RectangleFromMinMax(Pt(0, 0), Pt(10, 10))

	t.Run("a rectangle within is unchanged", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(5, 5), Sz(2, 2)).Clamp(box), Rect(Pt(5, 5), Sz(2, 2)))
	})
	t.Run("int rectangles that are not rotated are clamped on their exact corners", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(8, -7), Sz(9, 9)).Clamp(Rect(Pt(6, -6), Sz(8, 9))), Rect(Pt(6, -6), Sz(9, 9)))
	})
	t.Run("a rectangle outside moves in by the least", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(12, 5), Sz(4, 2)).Clamp(box), Rect(Pt(8, 5), Sz(4, 2)))
		geomtest.AssertRectangle(t, Rect(Pt(-3, -3), Sz(2, 2)).Clamp(box), Rect(Pt(1, 1), Sz(2, 2)))
	})
	t.Run("an axis larger than the other is centered", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(20, 3), Sz(14, 2)).Clamp(box), Rect(Pt(5, 3), Sz(14, 2)))
		geomtest.AssertRectangle(t, Rect(Pt(-7, 30), Sz(12, 16)).Clamp(box), Rect(Pt(5, 5), Sz(12, 16)))
	})
	t.Run("the same angle clamps in the shared frame", func(t *testing.T) {
		container := Rect(Pt(0.0, 0.0), Sz(10.0, 4.0)).Rotate(Pi / 2)

		geomtest.AssertRectangle(t, Rect(Pt(0.0, 9.0), Sz(2.0, 2.0)).Rotate(Pi/2).Clamp(container), Rect(Pt(0.0, 4.0), Sz(2.0, 2.0)).Rotate(Pi/2))
	})
	t.Run("a rotated rectangle is clamped by its corners", func(t *testing.T) {
		diamond := Rect(Pt(10.0, 5.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertRectangle(t, diamond.Clamp(box.Float()), Rect(Pt(10.0-Sqrt2, 5.0), Sz(2.0, 2.0)).Rotate(Pi/4))
	})
	t.Run("different angles clamp within the turned edges, not the bounds", func(t *testing.T) {
		container := Rect(Pt(0.0, 0.0), Sz(10.0, 2.0)).Rotate(Pi / 4)
		clamped := Rect(Pt(10.0, 10.0), Sz(1.0, 1.0)).Clamp(container)

		geomtest.AssertRectangle(t, clamped, Rect(Pt(5*OneOverSqrt2-0.5, 5*OneOverSqrt2-0.5), Sz(1.0, 1.0)))
		assert.True(t, within(clamped, container))
	})
	t.Run("int rounds the move once", func(t *testing.T) {
		container := Rect(Pt(0, 0), Sz(10, 4)).Rotate(Pi / 2)

		geomtest.AssertRectangle(t, Rect(Pt(0, 9), Sz(2, 2)).Clamp(container), Rect(Pt(0, 4), Sz(2, 2)))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, container := range rectFixtures {
				clamped := r.Clamp(container)
				message := fmt.Sprintf("%s → %s: ", r, container)

				geomtest.AssertSize(t, clamped.Size, r.Size, message)
				geomtest.AssertAngle(t, clamped.Angle, r.Angle, message)
				assert.True(t, clamped.Clamp(container).Equal(clamped), message)
				if within(r, container) {
					assert.True(t, clamped.Equal(r), message)
				}
				if EqualAngle(r.Angle, container.Angle) && r.Width() <= container.Width() && r.Height() <= container.Height() {
					assert.True(t, within(clamped, container), message)
				}
				if math.Hypot(r.Width(), r.Height()) <= min(container.Width(), container.Height()) {
					assert.True(t, within(clamped, container), message)
				}
			}
		}
	})
	t.Run("an int rectangle within one of an odd extent is not moved by the half unit of its frame", func(t *testing.T) {
		container := Rect(Pt(0, 0), Sz(3, 3))
		r := Rect(Pt(2, 2), Sz(0, 0)).Rotate(Pi)

		assert.True(t, container.EnclosesRectangle(r))
		geomtest.AssertRectangle(t, r.Clamp(container), r)
	})
	t.Run("an int rectangle within is kept at every quarter turn of either", func(t *testing.T) {
		for x := -3; x <= 3; x++ {
			for y := -3; y <= 3; y++ {
				for size := range 16 {
					for turn := range 4 {
						container := Rect(Pt(0, 0), Sz(5, 3)).Rotate(float64(turn) * Pi / 2)
						r := Rect(Pt(x, y), Sz(size%4, size/4)).Rotate(Pi)

						if container.EnclosesRectangle(r) {
							geomtest.AssertRectangle(t, r.Clamp(container), r, fmt.Sprintf("%s → %s: ", r, container))
						}
					}
				}
			}
		}
	})
}

// within reports whether every vertex of the rectangle lies within the other.
func within[T Number](r, container Rectangle[T]) bool {
	for vertex := range r.Vertices() {
		if !container.Contains(vertex) {
			return false
		}
	}

	return true
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
	t.Run("rotated excludes the corners of its bounds", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		assert.True(t, diamond.Contains(Pt(0.9, 0.0)))
		assert.True(t, diamond.Contains(Pt(Sqrt2, 0.0)))
		assert.False(t, diamond.Contains(Pt(0.9, 0.9)))
		assert.True(t, diamond.Bounds().Contains(Pt(0.9, 0.9)))
	})
	t.Run("rotated int contains its rounded corners", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		assert.True(t, turned.Contains(Pt(1, 2)))
		assert.True(t, turned.Contains(Pt(0, 2)))
		assert.False(t, turned.Contains(Pt(2, 0)))
	})
	t.Run("a NaN coordinate is contained by no rectangle", func(t *testing.T) {
		r := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))

		assert.False(t, r.Contains(Pt(math.NaN(), 0.0)))
		assert.False(t, r.Rotate(0.3).Contains(Pt(math.NaN(), 0.0)))
		assert.False(t, r.Polygon().Contains(Pt(0.0, math.NaN())))
	})
}

func TestRectangle_DistanceTo(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("beside an edge measures to the edge", func(t *testing.T) {
		geomtest.AssertNumber(t, rectangle.DistanceTo(Pt(5, 0)), 3.0)
		geomtest.AssertNumber(t, rectangle.DistanceTo(Pt(0, -6)), 4.0)
	})
	t.Run("beyond a corner measures to the corner", func(t *testing.T) {
		geomtest.AssertNumber(t, rectangle.DistanceTo(Pt(5, 6)), 5.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, rectangle.DistanceTo(Pt(1, -1)), 0.0)
		assert.Equal(t, rectangle.DistanceTo(Pt(2, 2)), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).DistanceTo(Pt(2.0, 2.0)), Sqrt2)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		r := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))

		assert.Equal(t, r.DistanceTo(Pt(1.0+Delta/2, 0.0)), 0.0)
		geomtest.AssertNumber(t, r.DistanceTo(Pt(1.0+2*Delta, 0.0)), 2*Delta)
	})
	t.Run("rotated measures to the turned edge", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertNumber(t, diamond.DistanceTo(Pt(0.9, 0.9)), (1.8-Sqrt2)/Sqrt2)
		geomtest.AssertNumber(t, Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi/2).DistanceTo(Pt(4, 0)), 3.0)
	})
	t.Run("zero exactly where Contains holds", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, r.DistanceTo(p) == 0, r.Contains(p), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRectangle_DistanceSquaredTo(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("is the square of DistanceTo, exact for an aligned int rectangle", func(t *testing.T) {
		assert.Equal(t, rectangle.DistanceSquaredTo(Pt(5, 6)), 25.0)
		assert.Equal(t, rectangle.DistanceSquaredTo(Pt(5, 0)), 9.0)
		assert.Equal(t, rectangle.DistanceSquaredTo(Pt(1, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, r.DistanceSquaredTo(p), r.DistanceTo(p)*r.DistanceTo(p), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
	t.Run("rotated int is not rounded, since the nearest point is off the lattice", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 4)).Rotate(Pi / 4)

		geomtest.AssertNumber(t, turned.DistanceTo(Pt(3, 3)), 3/Sqrt2)
		assert.Equal(t, turned.DistanceSquaredTo(Pt(3, 3)), 4.5)
	})
}

func TestRectangle_Nearest(t *testing.T) {
	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Nearest(Pt(2, 2)), Pt(2, 2))
	})
	t.Run("the foot on the nearest edge", func(t *testing.T) {
		geomtest.AssertPoint(t, Rect(Pt(1, 2), Sz(2, 3)).Nearest(Pt(10, 10)), Pt(2, 4))
		geomtest.AssertPoint(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Nearest(Pt(-1.0, 1.2)), Pt(0.0, 1.2))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Nearest(Pt(1.0+Delta/2, 0.5)), Pt(1.0+Delta/2, 0.5))
	})
	t.Run("rotated gives the turned edge, not the bounds", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertPoint(t, diamond.Nearest(Pt(2.0, 2.0)), Pt(OneOverSqrt2, OneOverSqrt2))
		geomtest.AssertPoint(t, diamond.Nearest(Pt(0.1, 0.2)), Pt(0.1, 0.2))
	})
	t.Run("rotated int rounds once", func(t *testing.T) {
		turned := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		geomtest.AssertPoint(t, turned.Nearest(Pt(5, 0)), Pt(1, 0))
		geomtest.AssertPoint(t, turned.Nearest(Pt(0, 7)), Pt(0, 2))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range pointFixtures {
				assertNearest[float64](t, r, p)
			}
		}
	})
}

func TestRectangle_EnclosesCircle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(10, 6))

	t.Run("inside", func(t *testing.T) {
		assert.True(t, rectangle.EnclosesCircle(Circ(Pt(1, 0), 2)))
	})
	t.Run("touching a side from inside counts", func(t *testing.T) {
		assert.True(t, rectangle.EnclosesCircle(Circ(Pt(0, 0), 3)))
		assert.False(t, rectangle.EnclosesCircle(Circ(Pt(0, 1), 3)))
	})
	t.Run("the center outside", func(t *testing.T) {
		assert.False(t, rectangle.EnclosesCircle(Circ(Pt(6, 0), 0)))
	})
	t.Run("rotated measures to the turned sides", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		assert.True(t, diamond.EnclosesCircle(Circ(Pt(0.0, 0.0), 1.0)))
		assert.False(t, diamond.EnclosesCircle(Circ(Pt(0.0, 0.0), 1.1)))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, r.EnclosesCircle(c), r.Polygon().EnclosesCircle(c), fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRectangle_EnclosesSegment(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("inside and along a side", func(t *testing.T) {
		assert.True(t, rectangle.EnclosesSegment(Seg(Pt(-1, -1), Pt(1, 2))))
		assert.True(t, rectangle.EnclosesSegment(Seg(Pt(-2, 2), Pt(2, 2))))
	})
	t.Run("an endpoint outside", func(t *testing.T) {
		assert.False(t, rectangle.EnclosesSegment(Seg(Pt(0, 0), Pt(3, 0))))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, r.EnclosesSegment(s), r.Polygon().EnclosesSegment(s), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestRectangle_EnclosesPolygon(t *testing.T) {
	t.Run("inside and outside", func(t *testing.T) {
		assert.True(t, Rect(Pt(1, 1), Sz(2, 2)).EnclosesPolygon(Pol(squareVertices())))
		assert.False(t, Rect(Pt(1, 1), Sz(2, 1)).EnclosesPolygon(Pol(squareVertices())))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesPolygon(Pol[int](nil)))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, r.EnclosesPolygon(p), r.Polygon().EnclosesPolygon(p), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRectangle_EnclosesRectangle(t *testing.T) {
	t.Run("inside and sharing a side", func(t *testing.T) {
		assert.True(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesRectangle(Rect(Pt(1, 0), Sz(2, 4))))
		assert.False(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesRectangle(Rect(Pt(1, 0), Sz(2, 6))))
	})
	t.Run("a turned square fits only where its corners do", func(t *testing.T) {
		square := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))

		assert.False(t, square.EnclosesRectangle(square.Rotate(Pi/4)))
		assert.True(t, square.Rotate(Pi/4).EnclosesRectangle(square.Unscale(2)))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				assert.Equal(t, a.EnclosesRectangle(b), a.Polygon().EnclosesRectangle(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRectangle_EnclosesRegularPolygon(t *testing.T) {
	t.Run("inside and outside", func(t *testing.T) {
		assert.True(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 6, 0, 0)))
		assert.False(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesRegularPolygon(RegPol(Pt(1, 0), Sz(2, 2), 6, 0, 0)))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, Rect(Pt(0, 0), Sz(4, 4)).EnclosesRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, r.EnclosesRegularPolygon(rp), r.Polygon().EnclosesRegularPolygon(rp), fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRectangle_EnclosesBox(t *testing.T) {
	t.Run("the bounds of a rotated rectangle are not enclosed by it", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		assert.False(t, diamond.EnclosesBox(diamond.Bounds()))
		assert.True(t, diamond.Bounds().Rectangle().EnclosesRectangle(diamond))
	})
	t.Run("matches the polygon of the corners", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, r.EnclosesBox(b), r.Polygon().EnclosesBox(b), fmt.Sprintf("%s → %s: ", r, b))
			}
		}
	})
}

func TestRectangle_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, r.IntersectsCircle(c), c.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRectangle_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, r.IntersectsSegment(s), s.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestRectangle_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, r.IntersectionSegment(s), s.IntersectionRectangle(r), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestRectangle_AppendIntersectionSegment(t *testing.T) {
	t.Run("matches Segment.AppendIntersectionRectangle", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, rect.AppendIntersectionSegment(bufferWith(prefixPoint), s), s.AppendIntersectionRectangle(bufferWith(prefixPoint), rect), fmt.Sprintf("%s → %s: ", rect, s))
			}
		}
	})
}

func TestRectangle_IntersectsRay(t *testing.T) {
	t.Run("mirrors Ray.IntersectsRectangle", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, rect.IntersectsRay(r), r.IntersectsRectangle(rect), fmt.Sprintf("%s → %s: ", rect, r))
			}
		}
	})
}

func TestRectangle_IntersectionRay(t *testing.T) {
	t.Run("matches Ray.IntersectionRectangle", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, rect.IntersectionRay(r), r.IntersectionRectangle(rect), fmt.Sprintf("%s → %s: ", rect, r))
			}
		}
	})
}

func TestRectangle_AppendIntersectionRay(t *testing.T) {
	t.Run("matches Ray.AppendIntersectionRectangle", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, rect.AppendIntersectionRay(bufferWith(prefixPoint), r), r.AppendIntersectionRectangle(bufferWith(prefixPoint), rect), fmt.Sprintf("%s → %s: ", rect, r))
			}
		}
	})
}

func TestRectangle_IntersectsPolygon(t *testing.T) {
	t.Run("rotated is tested on its turned edges", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)
		corner := Pol([]Point[float64]{{1, 1}, {3, 1}, {3, 3}})

		assert.False(t, diamond.IntersectsPolygon(corner))
		assert.True(t, diamond.Bounds().Rectangle().IntersectsPolygon(corner))
	})
	t.Run("mirrors Polygon.IntersectsRectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, r.IntersectsPolygon(p), p.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRectangle_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsRectangle(Rect(Pt(100.0, -50.0), Sz(200.0, 50.0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, rectangle.IntersectsRectangle(Rect(Pt(100.0, 350.0), Sz(200.0, 450.0))))
	})
	t.Run("rotated int rectangles of one angle are judged in their shared frame", func(t *testing.T) {
		assert.True(t, Rect(Pt(-7, 1), Sz(3, 7)).Rotate(3.8).IntersectsRectangle(Rect(Pt(-1, -3), Sz(2, 6)).Rotate(3.8)))
		assert.False(t, Rect(Pt(5, -6), Sz(4, 12)).Rotate(5.5).IntersectsRectangle(Rect(Pt(-3, -9), Sz(2, 5)).Rotate(5.5)))
	})
	t.Run("a shared edge counts as an intersection", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsRectangle(Rect(Pt(200.0, 0.0), Sz(200.0, 100.0))))
		assert.False(t, rectangle.IntersectsRectangle(Rect(Pt(201.0, 0.0), Sz(200.0, 100.0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, rectangle.IntersectsRectangle(Rect(Pt(0.0, 0.0), Sz(50.0, 50.0))))
	})
	t.Run("a corner gap within the tolerance on both axes is decided on the distance", func(t *testing.T) {
		r := RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, r.IntersectsRectangle(RectangleFromMinMax(Pt(1.0+Delta/2, 1.0), Pt(2.0, 2.0))))
		assert.False(t, r.IntersectsRectangle(RectangleFromMinMax(Pt(1.0+0.9*Delta, 1.0+0.9*Delta), Pt(2.0, 2.0))))
		assert.True(t, r.IntersectsRectangle(RectangleFromMinMax(Pt(1.0+0.7*Delta, 1.0+0.7*Delta), Pt(2.0, 2.0))))
	})
	t.Run("rotated against aligned is decided on the edges, not the bounds", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		assert.False(t, diamond.IntersectsRectangle(RectangleFromMinMax(Pt(1.0, 1.0), Pt(3.0, 3.0))))
		assert.True(t, diamond.Bounds().IntersectsBox(BoxFromMinMax(Pt(1.0, 1.0), Pt(3.0, 3.0))))
		assert.True(t, diamond.IntersectsRectangle(Rect(Pt(1.2, 0.0), Sz(1.0, 1.0))))
		assert.True(t, diamond.IntersectsRectangle(Rect(Pt(0.0, 0.0), Sz(0.5, 0.5))))
		assert.True(t, Rect(Pt(0.0, 0.0), Sz(0.5, 0.5)).IntersectsRectangle(diamond))
	})
	t.Run("rotated of the same angle is tested in the shared frame", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2)

		assert.True(t, a.IntersectsRectangle(Rect(Pt(0.0, 1.0), Sz(4.0, 2.0)).Rotate(Pi/2)))
		assert.True(t, a.IntersectsRectangle(Rect(Pt(2.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi/2)))
		assert.False(t, a.IntersectsRectangle(Rect(Pt(3.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi/2)))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				assert.Equal(t, a.IntersectsRectangle(b), b.IntersectsRectangle(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("a turn too small to judge on the angle still moves a far corner of a long rectangle", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(1e7, 10.0)).Rotate(0.1)
		b := Rect(Pt(0.0, 0.0).Add(Vec(0.0, 13.0).Rotate(0.1)), Sz(1e7, 10.0)).Rotate(0.1 + 9e-7)

		assert.True(t, a.IntersectsRectangle(b))
		assert.True(t, b.IntersectsRectangle(a))
	})
	t.Run("long rectangles turned a hair apart answer as their polygons do", func(t *testing.T) {
		for _, r := range rectFixtures {
			a := r.Scale(1e6)
			for _, gap := range []float64{0.1, 1.0, 10.0} {
				b := a.Translate(Vec(0.0, a.Height()+gap).Rotate(a.Angle)).Rotate(9e-7)

				assert.Equal(t, a.IntersectsRectangle(b), a.Polygon().IntersectsPolygon(b.Polygon()), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func BenchmarkRectangle_IntersectsRectangle(b *testing.B) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
	other := Rect(Pt(5.0, 5.0), Sz(10.0, 10.0))

	for b.Loop() {
		sinkBool = rectangle.IntersectsRectangle(other)
	}
}

func TestRectangle_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		overlap, ok := rectangle.IntersectionRectangle(Rect(Pt(2, 2), Sz(4, 4)))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, RectangleFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := rectangle.IntersectionRectangle(Rect(Pt(10, 10), Sz(4, 4)))

		assert.False(t, ok)
	})
	t.Run("rotated int rectangles of one angle overlap only where their shared frame holds one", func(t *testing.T) {
		_, ok := Rect(Pt(5, -6), Sz(4, 12)).Rotate(5.5).IntersectionRectangle(Rect(Pt(-3, -9), Sz(2, 5)).Rotate(5.5))
		assert.False(t, ok)

		overlap, ok := Rect(Pt(-7, 1), Sz(3, 7)).Rotate(3.8).IntersectionRectangle(Rect(Pt(-1, -3), Sz(2, 6)).Rotate(3.8))
		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, Rect(Pt(-4, -2), Sz(0, 0)).Rotate(3.8))
	})
	t.Run("half a turn apart overlap in a rectangle of the receiver's angle", func(t *testing.T) {
		turned := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3)
		overlap, ok := turned.IntersectionRectangle(turned.Rotate(Pi))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, turned)
	})
	t.Run("touching gives a zero extent", func(t *testing.T) {
		overlap, ok := rectangle.IntersectionRectangle(Rect(Pt(4, 0), Sz(4, 4)))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, RectangleFromMinMax(Pt(2, -2), Pt(2, 2)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		inner := Rect(Pt(0, 0), Sz(2, 2))

		overlap, ok := rectangle.IntersectionRectangle(inner)

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, inner)
	})
	t.Run("float", func(t *testing.T) {
		overlap, ok := Rect(Pt(0.0, 0.0), Sz(3.0, 3.0)).IntersectionRectangle(Rect(Pt(1.0, 1.0), Sz(3.0, 3.0)))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, RectangleFromMinMax(Pt(-0.5, -0.5), Pt(1.5, 1.5)))
	})
	t.Run("rotated of the same angle overlap in a rectangle of that angle", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2)

		overlap, ok := a.IntersectionRectangle(Rect(Pt(0.0, 1.0), Sz(4.0, 2.0)).Rotate(Pi / 2))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, Rect(Pt(0.0, 0.5), Sz(3.0, 2.0)).Rotate(Pi/2))
	})
	t.Run("rotated int rounds the shared frame", func(t *testing.T) {
		a := Rect(Pt(0, 0), Sz(4, 2)).Rotate(Pi / 2)

		overlap, ok := a.IntersectionRectangle(Rect(Pt(0, 1), Sz(4, 2)).Rotate(Pi / 2))

		assert.True(t, ok)
		geomtest.AssertRectangle(t, overlap, Rect(Pt(0, 0), Sz(3, 2)).Rotate(Pi/2))
	})
	t.Run("rotated of the same angle apart", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2)

		_, ok := a.IntersectionRectangle(Rect(Pt(3.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2))

		assert.False(t, ok)
	})
	t.Run("different angles have no rectangle in common", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)
		square := Rect(Pt(1.2, 0.0), Sz(1.0, 1.0))

		_, ok := diamond.IntersectionRectangle(square)

		assert.False(t, ok)
		assert.True(t, diamond.IntersectsRectangle(square))
	})
	t.Run("rotated int answers exactly where Intersects holds", func(t *testing.T) {
		// the world extent of the rounded corners rejects these pairs while the rounded
		// offset in the shared frame still finds a touch
		a, b := Rect(Pt(37, 39), Sz(10, 3)).Rotate(3.22), Rect(Pt(29, 34), Sz(9, 5)).Rotate(3.22)

		assert.False(t, a.IntersectsRectangle(b))
		_, ok := a.IntersectionRectangle(b)
		assert.False(t, ok)

		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				for _, angle := range []float64{Pi / 5, 3.22, -1.16} {
					ra, rb := a.Int().Rotate(angle), b.Int().Rotate(angle)
					_, ok := ra.IntersectionRectangle(rb)

					assert.Equal(t, ok, ra.IntersectsRectangle(rb) && EqualAngle(ra.Angle, rb.Angle), fmt.Sprintf("%s → %s: ", ra, rb))
				}
			}
		}
	})
	t.Run("symmetric and contained by both where the angles agree", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				ab, okAB := a.IntersectionRectangle(b)
				ba, okBA := b.IntersectionRectangle(a)

				assert.Equal(t, okAB, okBA, fmt.Sprintf("%s → %s: ", a, b))
				assert.Equal(t, okAB, a.IntersectsRectangle(b) && EqualAngle(a.Angle, b.Angle), fmt.Sprintf("%s → %s: ", a, b))
				if !okAB {
					continue
				}

				geomtest.AssertRectangle(t, ab, ba, fmt.Sprintf("%s → %s: ", a, b))
				for corner := range ab.Vertices() {
					assert.True(t, a.Contains(corner), fmt.Sprintf("%s → %s: %s in a: ", a, b, corner))
					assert.True(t, b.Contains(corner), fmt.Sprintf("%s → %s: %s in b: ", a, b, corner))
				}
			}
		}
	})
	t.Run("long rectangles turned a hair apart share no frame", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(1e7, 10.0)).Rotate(0.1)

		_, ok := a.IntersectionRectangle(a.Rotate(9e-7))
		assert.False(t, ok)
	})
}

func TestRectangle_Union(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Union(Rect(Pt(2, 2), Sz(4, 4))), RectangleFromMinMax(Pt(-2, -2), Pt(4, 4)))
	})
	t.Run("apart spans the gap", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Union(Rect(Pt(10, 0), Sz(4, 4))), RectangleFromMinMax(Pt(-2, -2), Pt(12, 2)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		geomtest.AssertRectangle(t, rectangle.Union(Rect(Pt(0, 0), Sz(2, 2))), rectangle)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.0, 0.0), Sz(3.0, 3.0)).Union(Rect(Pt(1.0, 1.0), Sz(3.0, 3.0))), RectangleFromMinMax(Pt(-1.5, -1.5), Pt(2.5, 2.5)))
	})
	t.Run("rotated of the same angle unite in a rectangle of that angle", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2)

		geomtest.AssertRectangle(t, a.Union(Rect(Pt(0.0, 1.0), Sz(4.0, 2.0)).Rotate(Pi/2)), Rect(Pt(0.0, 0.5), Sz(5.0, 2.0)).Rotate(Pi/2))
	})
	t.Run("half a turn apart unite in a rectangle of the receiver's angle", func(t *testing.T) {
		turned := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3)

		geomtest.AssertRectangle(t, turned.Union(turned.Rotate(Pi)), turned)
	})
	t.Run("different angles unite in the box around both", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertRectangle(t, diamond.Union(Rect(Pt(1.2, 0.0), Sz(1.0, 1.0))), RectangleFromMinMax(Pt(-Sqrt2, -Sqrt2), Pt(1.7, Sqrt2)))
	})
	t.Run("symmetric and contains both", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				union := a.Union(b)

				geomtest.AssertRectangle(t, union, b.Union(a), fmt.Sprintf("%s → %s: ", a, b))
				for _, r := range [2]Rectangle[float64]{a, b} {
					for corner := range r.Vertices() {
						assert.True(t, union.Contains(corner), fmt.Sprintf("%s → %s: %s: ", a, b, corner))
					}
				}
			}
		}
	})
	t.Run("long rectangles turned a hair apart unite in their bounds", func(t *testing.T) {
		a := Rect(Pt(0.0, 0.0), Sz(1e7, 10.0)).Rotate(0.1)
		b := a.Rotate(9e-7)

		geomtest.AssertRectangle(t, a.Union(b), a.Bounds().Union(b.Bounds()).Rectangle())
	})
}

func TestRectangle_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, Rect(Pt(1, 1), Sz(2, 2)).IntersectsRegularPolygon(diamond))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Rect(Pt(4, 4), Sz(2, 2)).IntersectsRegularPolygon(diamond))
		assert.False(t, Rect(Pt(2, 2), Sz(1, 1)).IntersectsRegularPolygon(diamond))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, Rect(Pt(0, 0), Sz(10, 10)).IntersectsRegularPolygon(diamond))
		assert.True(t, Rect(Pt(0, 0), Sz(1, 1)).IntersectsRegularPolygon(diamond))
	})
	t.Run("edges crossing without a corner inside", func(t *testing.T) {
		assert.True(t, Rect(Pt(0, 0), Sz(10, 1)).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 20), 4, 0, 0)))
	})
	t.Run("rotated is tested on its turned edges", func(t *testing.T) {
		turned := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)
		corner := RegPol(Pt(1.5, 1.5), Sz(0.5, 0.5), 3, 0, 0)

		assert.False(t, turned.IntersectsRegularPolygon(corner))
		assert.True(t, turned.Bounds().Rectangle().IntersectsRegularPolygon(corner))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Rect(Pt(0, 0), Sz(2, 2)).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, r.IntersectsRegularPolygon(rp), r.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRectangle_IntersectsBox(t *testing.T) {
	t.Run("a shared edge counts as an intersection", func(t *testing.T) {
		assert.True(t, Rect(Pt(0, 0), Sz(4, 4)).IntersectsBox(BoxFromMinMax(Pt(2, -1), Pt(5, 1))))
		assert.False(t, Rect(Pt(0, 0), Sz(4, 4)).IntersectsBox(BoxFromMinMax(Pt(3, -1), Pt(5, 1))))
	})
	t.Run("rotated is tested on its turned edges", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)
		corner := BoxFromMinMax(Pt(1.0, 1.0), Pt(3.0, 3.0))

		assert.False(t, diamond.IntersectsBox(corner))
		assert.True(t, diamond.Bounds().IntersectsBox(corner))
	})
	t.Run("not rotated agrees with its bounds", func(t *testing.T) {
		for _, r := range rectFixtures {
			if !r.IsAligned() {
				continue
			}

			for _, b := range boxFixtures {
				assert.Equal(t, r.IntersectsBox(b), r.Bounds().IntersectsBox(b), fmt.Sprintf("%s → %s: ", r, b))
			}
		}
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
	t.Run("the angle counts, up to a full turn", func(t *testing.T) {
		r := Rect(Pt(1, 2), Sz(2, 3))

		assert.False(t, r.Equal(r.Rotate(Pi/2)))
		assert.True(t, r.Equal(r.Rotate(2*Pi)))
		assert.True(t, r.Rotate(Pi/2).Equal(Rectangle[int]{r.Center, r.Size, -3 * Pi / 2}))
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
	t.Run("a full turn is zero", func(t *testing.T) {
		assert.True(t, Rectangle[int]{Angle: 2 * Pi}.IsZero())
		assert.False(t, Rectangle[int]{Angle: Pi}.IsZero())
	})
}

func TestRectangle_IsAligned(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(2, 3))

	t.Run("not rotated", func(t *testing.T) {
		assert.True(t, r.IsAligned())
		assert.True(t, RectangleFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)).IsAligned())
		assert.True(t, Rectangle[int]{}.IsAligned())
	})
	t.Run("rotated", func(t *testing.T) {
		assert.False(t, r.Rotate(Pi/2).IsAligned())
		assert.False(t, Rectangle[int]{r.Center, r.Size, Delta / 2}.IsAligned())
	})
	t.Run("a full turn is aligned again", func(t *testing.T) {
		assert.True(t, r.Rotate(2*Pi).IsAligned())
		assert.True(t, Rectangle[int]{r.Center, r.Size, 2 * Pi}.Canonical().IsAligned())
	})
}

func TestRectangle_Polygon(t *testing.T) {
	r := Rect(Pt(0, 0), Sz(2, 2))
	p := r.Polygon()

	t.Run("carries the vertices", func(t *testing.T) {
		geomtest.AssertVertices(t, p.Points, slices.Collect(r.Vertices()))
	})
	t.Run("owns its slice", func(t *testing.T) {
		assert.NotSame(t, p.Points, r.Polygon().Points)
	})
	t.Run("rotated carries the turned vertices", func(t *testing.T) {
		turned := r.Rotate(Pi / 2)

		geomtest.AssertVertices(t, turned.Polygon().Points, slices.Collect(turned.Vertices()))
		assert.True(t, turned.Polygon().Contains(Pt(1, 1)))
	})
}

func TestRectangle_Cast(t *testing.T) {
	r := Rect(Pt(1.5, -2.5), Sz(3.5, 4.5)).Rotate(Pi / 6)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertRectangle(t, r.Cast[int](), r.Int())
		geomtest.AssertRectangle(t, r.Cast[float64](), r.Float())
	})
	t.Run("a type the other conversions cannot name, and the angle is kept", func(t *testing.T) {
		geomtest.AssertRectangle(t, r.Cast[int8](), Rect(Pt[int8](2, -3), Sz[int8](4, 5)).Rotate(Pi/6))
	})
}

func TestRectangle_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Int(), Rect(Pt(1, 2), Sz(2, 3)))
	})
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Int(), Rect(Pt(1, 0), Sz(1, 4)))
	})
	t.Run("keeps the size where the corners would not", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.5, 0.5), SzU(16.0)).Int(), Rect(Pt(1, 1), SzU(16)))
	})
	t.Run("keeps the angle", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Rotate(Pi/2).Int(), Rect(Pt(1, 0), Sz(1, 4)).Rotate(Pi/2))
	})
}

func TestRectangle_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(1, 2), Sz(2, 3)).Float(), Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertRectangle(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Float(), Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
}

func TestRectangle_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).String(), "Rect((1,2);2x3)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).String(), "Rect((0.60,-0.25);1.20x3.60)")
	})
	t.Run("rotated appends the angle", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Rotate(Pi/2).String(), "Rect((1,2);2x3;1.57)")
	})
}

func TestRectangle_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Rect(Pt(1, 2), Sz(2, 3)), `{"x":1,"y":2,"w":2,"h":3}`)

		var r Rectangle[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":3}`), &r))
		geomtest.AssertRectangle(t, r, Rect(Pt(1, 2), Sz(2, 3)))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)), `{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`)

		var r Rectangle[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":0.60,"y":-0.25,"w":1.20,"h":3.60}`), &r))
		geomtest.AssertRectangle(t, r, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)))
	})
	t.Run("rotated carries the angle, absent for a rectangle that is not rotated", func(t *testing.T) {
		assert.JSON(t, Rect(Pt(1, 2), Sz(2, 3)).Rotate(Pi/2), `{"x":1,"y":2,"w":2,"h":3,"a":1.5707963267948966}`)

		var r Rectangle[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":3,"a":1.5707963267948966}`), &r))
		geomtest.AssertRectangle(t, r, Rect(Pt(1, 2), Sz(2, 3)).Rotate(Pi/2))
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
	t.Run("min and max bracket the center and span the size before the turn", func(t *testing.T) {
		for _, r := range rectFixtures {
			assert.True(t, r.Min().Midpoint(r.Max()).Equal(r.Center), fmt.Sprintf("%s: ", r))
			if !r.IsAligned() {
				continue
			}

			geomtest.AssertNumber(t, r.Max().X-r.Min().X, r.Width(), fmt.Sprintf("%s: ", r))
			geomtest.AssertNumber(t, r.Max().Y-r.Min().Y, r.Height(), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("area and perimeter follow the size", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertNumber(t, r.Area(), r.Size.Area(), fmt.Sprintf("%s: ", r))
			geomtest.AssertNumber(t, r.Perimeter(), r.Size.Perimeter(), fmt.Sprintf("%s: ", r))
			geomtest.AssertNumber(t, r.AspectRatio(), r.Size.AspectRatio(), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("corners are the corners before the turn, turned about the center", func(t *testing.T) {
		for _, r := range rectFixtures {
			flat := Rectangle[float64]{r.Center, r.Size, 0}
			a, b := flat.MinMax()

			geomtest.AssertPoint(t, r.TopLeft(), a.RotateAround(r.Center, r.Angle), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, r.BottomRight(), b.RotateAround(r.Center, r.Angle), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, r.TopRight(), Pt(b.X, a.Y).RotateAround(r.Center, r.Angle), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, r.BottomLeft(), Pt(a.X, b.Y).RotateAround(r.Center, r.Angle), fmt.Sprintf("%s: ", r))

			edges := slices.Collect(r.Edges())

			geomtest.AssertPoint(t, edges[0].Midpoint(), r.Anchor(Top), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, edges[1].Midpoint(), r.Anchor(DirectionRight), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, edges[2].Midpoint(), r.Anchor(Bottom), fmt.Sprintf("%s: ", r))
			geomtest.AssertPoint(t, edges[3].Midpoint(), r.Anchor(DirectionLeft), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("lerp ends on the two rectangles", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, r.Scale(factor).Unscale(factor).Equal(r), fmt.Sprintf("%s ×%v: ", r, factor))
				assert.True(t, r.ScaleXY(factor, 2).UnscaleXY(factor, 2).Equal(r), fmt.Sprintf("%s ×%v: ", r, factor))
			}
		}
	})
	t.Run("bounds is the box around the vertices, the rectangle itself before a turn", func(t *testing.T) {
		for _, r := range rectFixtures {
			bounds := r.Bounds()

			assert.True(t, bounds.Equal(r.Polygon().Bounds()), fmt.Sprintf("%s: ", r))
			assert.Equal(t, bounds.Rectangle().Equal(r), r.IsAligned(), fmt.Sprintf("%s: ", r))
		}
	})
	t.Run("rotate keeps the center and the size", func(t *testing.T) {
		for _, r := range rectFixtures {
			turned := r.Rotate(Pi / 3)

			assert.True(t, turned.Center.Equal(r.Center), fmt.Sprintf("%s: ", r))
			assert.True(t, turned.Size.Equal(r.Size), fmt.Sprintf("%s: ", r))
			assert.True(t, EqualAngle(turned.Angle, r.Angle+Pi/3), fmt.Sprintf("%s: ", r))
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

	t.Run("vertices and edges describe the same outline", func(t *testing.T) {
		for _, r := range rectFixtures {
			vertices, edges := slices.Collect(r.Vertices()), slices.Collect(r.Edges())

			assert.Length(t, vertices, 4, fmt.Sprintf("%s: ", r))
			assert.Length(t, edges, 4, fmt.Sprintf("%s: ", r))

			for i, edge := range edges {
				assert.True(t, edge.Start.Equal(vertices[i]), fmt.Sprintf("%s: ", r))
				assert.True(t, edge.End.Equal(vertices[(i+1)%4]), fmt.Sprintf("%s: ", r))
			}
		}
	})
	t.Run("anchors lie on the rectangle", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, direction := range Directions() {
				assert.True(t, r.Contains(r.Anchor(direction)), fmt.Sprintf("%s → %s: ", r, direction))
			}
		}
	})
	t.Run("polygon carries the vertices", func(t *testing.T) {
		for _, r := range rectFixtures {
			geomtest.AssertVertices(t, r.Polygon().Points, slices.Collect(r.Vertices()))
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

	geomtest.AssertRectangle(t, r, Rect(Pt(1, 2), Sz(2, 3)))
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
	Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 2),
	Rect(Pt(1.0, 2.0), Sz(2.0, 3.0)).Rotate(Pi / 6),
	Rect(Pt(-3.5, 0.25), Sz(7.0, 1.0)).Rotate(-Pi / 4),
}

func ExampleRect() {
	fmt.Println(Rect(Pt(1, 2), Sz(2, 3)))
	// Output: Rect((1,2);2x3)
}

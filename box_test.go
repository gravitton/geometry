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

func TestBox_Constructor(t *testing.T) {
	t.Run("shorthand", func(t *testing.T) {
		geomtest.AssertBox(t, Bx(Pt(0, 0), Pt(4, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
		geomtest.AssertBox(t, Bx(Pt(1.0, 3.0), Pt(0.0, 0.0)), Box[float64]{Pt(0.0, 0.0), Pt(1.0, 3.0)})
	})
	t.Run("from min", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMin(Pt(1, 2), Sz(4, 3)), Box[int]{Pt(1, 2), Pt(5, 5)})
		geomtest.AssertBox(t, BoxFromMin(Pt(0.5, -1.0), Sz(1.5, 2.5)), Box[float64]{Pt(0.5, -1.0), Pt(2.0, 1.5)})
	})
	t.Run("from max", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMax(Pt(5, 5), Sz(4, 3)), Box[int]{Pt(1, 2), Pt(5, 5)})
		geomtest.AssertBox(t, BoxFromMax(Pt(2.0, 1.5), Sz(1.5, 2.5)), Box[float64]{Pt(0.5, -1.0), Pt(2.0, 1.5)})
	})
	t.Run("from min and max", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0, 0), Pt(4, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), Box[float64]{Pt(0.0, 0.0), Pt(1.0, 3.0)})
	})
	t.Run("from size at the origin", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromSize(Sz(4, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
	})
	t.Run("corners in either order", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(4, 2), Pt(0, 0)), Box[int]{Pt(0, 0), Pt(4, 2)})
		geomtest.AssertBox(t, BoxFromMinMax(Pt(4, 0), Pt(0, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
	})
	t.Run("a negative extent measures the other way from the corner", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMin(Pt(4, 2), Sz(-4, -2)), Box[int]{Pt(0, 0), Pt(4, 2)})
		geomtest.AssertBox(t, BoxFromSize(Sz(-4, 2)), Box[int]{Pt(-4, 0), Pt(0, 2)})
		geomtest.AssertBox(t, BoxFromMax(Pt(0, 0), Sz(-4, -2)), Box[int]{Pt(0, 0), Pt(4, 2)})
	})
}

func TestBox_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Width(), 3)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Width(), 1.2)
	})
}

func TestBox_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Height(), 5)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Height(), 3.6)
	})
}

func TestBox_Size(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSize(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Size(), Sz(3, 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSize(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Size(), Sz(1.2, 3.6))
	})
}

func TestBox_Anchor(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("a diagonal gives a corner, Top the side of Min.Y", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Anchor(TopLeft), Pt(0, 0))
		geomtest.AssertPoint(t, box.Anchor(TopRight), Pt(4, 0))
		geomtest.AssertPoint(t, box.Anchor(BottomRight), Pt(4, 2))
		geomtest.AssertPoint(t, box.Anchor(BottomLeft), Pt(0, 2))
	})
	t.Run("a cardinal gives the midpoint of an edge", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Anchor(Top), Pt(2, 0))
		geomtest.AssertPoint(t, box.Anchor(Right), Pt(4, 1))
		geomtest.AssertPoint(t, box.Anchor(Bottom), Pt(2, 2))
		geomtest.AssertPoint(t, box.Anchor(Left), Pt(0, 1))
	})
	t.Run("none gives the center", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Anchor(DirectionNone), box.Center())
	})
	t.Run("an int midpoint of an odd extent keeps nearer Min", func(t *testing.T) {
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(0, 0), Pt(3, 5)).Anchor(Bottom), Pt(1, 5))
	})
	t.Run("matches Anchor on the box's Rectangle", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, direction := range Directions() {
				geomtest.AssertPoint(t, b.Anchor(direction), b.Rectangle().Anchor(direction), fmt.Sprintf("%s %s: ", b, direction))
			}
		}
	})
}

func TestBox_Vertices(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	t.Run("the corners from Min, clockwise on screen", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(box.Vertices()), []Point[int]{Pt(1, 2), Pt(4, 2), Pt(4, 7), Pt(1, 7)})
	})
	t.Run("are the vertices of its rectangle", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(box.Vertices()), slices.Collect(box.Rectangle().Vertices()))
	})
	t.Run("stops early", func(t *testing.T) {
		for vertex := range box.Vertices() {
			geomtest.AssertPoint(t, vertex, Pt(1, 2))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			for vertex := range box.Vertices() {
				sinkBool = vertex.IsZero()
			}
		}), 0)
	})
}

func TestBox_Edges(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	t.Run("each from a corner to the next, closing back to Min", func(t *testing.T) {
		geomtest.AssertSegments(t, slices.Collect(box.Edges()), []Segment[int]{
			Seg(Pt(1, 2), Pt(4, 2)), Seg(Pt(4, 2), Pt(4, 7)), Seg(Pt(4, 7), Pt(1, 7)), Seg(Pt(1, 7), Pt(1, 2)),
		})
	})
	t.Run("stops early", func(t *testing.T) {
		for edge := range box.Edges() {
			geomtest.AssertSegment(t, edge, Seg(Pt(1, 2), Pt(4, 2)))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			for edge := range box.Edges() {
				sinkBool = edge.IsZero()
			}
		}), 0)
	})
}

func TestBox_Center(t *testing.T) {
	t.Run("int truncates toward min", func(t *testing.T) {
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Center(), Pt(2, 4))
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(-4, -7), Pt(-1, -2)).Center(), Pt(-3, -5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Center(), Pt(1.2, 1.55))
	})
}

func TestBox_Centroid(t *testing.T) {
	t.Run("is the center", func(t *testing.T) {
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Centroid(), Pt(2, 4))

		for _, box := range boxFixtures {
			geomtest.AssertPoint(t, box.Centroid(), box.Center(), box.String())
		}
	})
}

func TestBox_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Area(), 15.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(0.5, 1.0), Pt(2.0, 3.0)).Area(), 3.0)
	})
	t.Run("is the area of its rectangle", func(t *testing.T) {
		for _, box := range boxFixtures {
			geomtest.AssertNumber(t, box.Area(), box.Rectangle().Area(), box.String())
		}
	})
}

func TestBox_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Perimeter(), 16)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(0.5, 1.0), Pt(2.0, 3.0)).Perimeter(), 7.0)
	})
}

func TestBox_Inertia(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, BoxFromMinMax(Pt(0, 0), Pt(4, 2)).Inertia(), 40.0/3)
	})
	t.Run("is the inertia of its rectangle", func(t *testing.T) {
		for _, box := range boxFixtures {
			geomtest.AssertNumber(t, box.Inertia(), box.Rectangle().Inertia(), box.String())
		}
	})
	t.Run("the products are rounded before they are added", func(t *testing.T) {
		w, h := tenth*3, tenth*7
		assert.Equal(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(w, h)).Inertia(), w*h*(float64(w*w)+float64(h*h))/12)
	})
}

func TestBox_Bounds(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	geomtest.AssertBox(t, box.Bounds(), box)
}

func TestBox_BoundingCircle(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, BoxFromMinMax(Pt(0, 0), Pt(6, 8)).BoundingCircle(), Circ(Pt(3, 4), 5))
	})
	t.Run("int reaches the max corner from the truncated center", func(t *testing.T) {
		geomtest.AssertCircle(t, BoxFromMinMax(Pt(0, 0), Pt(3, 3)).BoundingCircle(), Circ(Pt(1, 1), 3))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, BoxFromMinMax(Pt(0.5, -1.0), Pt(3.5, 3.0)).BoundingCircle(), Circ(Pt(2.0, 1.0), 2.5))
	})
	t.Run("contains every corner, rounded to int too", func(t *testing.T) {
		for _, b := range boxFixtures {
			assertBoundingCircle(t, b.BoundingCircle(), b.Rectangle().Vertices(), b.String())
			assertBoundingCircle(t, b.Int().BoundingCircle(), b.Int().Rectangle().Vertices(), b.String())
		}
	})
}

func TestBox_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Translate(Vec(3, -2)), BoxFromMinMax(Pt(4, 0), Pt(7, 5)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Translate(Vec(100.1, -0.1)), BoxFromMinMax(Pt(100.7, -0.35), Pt(101.9, 3.25)))
	})
}

func TestBox_MoveTo(t *testing.T) {
	t.Run("int keeps the size and puts the center at the point", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).MoveTo(Pt(10, 10)), BoxFromMinMax(Pt(9, 8), Pt(12, 13)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(2.0, 4.0)).MoveTo(Pt(-1.0, 0.5)), BoxFromMinMax(Pt(-2.0, -1.5), Pt(0.0, 2.5)))
	})
	t.Run("the center lands on the point and the size stays", func(t *testing.T) {
		for _, box := range boxFixtures {
			for _, point := range pointFixtures {
				moved := box.MoveTo(point)

				geomtest.AssertPoint(t, moved.Center(), point, box.String())
				geomtest.AssertSize(t, moved.Size(), box.Size(), box.String())
			}
		}
	})
}

func TestBox_AlignTo(t *testing.T) {
	box := BoxFromMin(Pt(0, 0), Sz(10, 20))

	t.Run("a corner lands on the point", func(t *testing.T) {
		geomtest.AssertBox(t, box.AlignTo(TopLeft, Pt(100, 50)), BoxFromMin(Pt(100, 50), Sz(10, 20)))
		geomtest.AssertBox(t, box.AlignTo(BottomRight, Pt(100, 50)), BoxFromMax(Pt(100, 50), Sz(10, 20)))
	})
	t.Run("none aligns the center", func(t *testing.T) {
		geomtest.AssertBox(t, box.AlignTo(DirectionNone, Pt(100, 50)), box.MoveTo(Pt(100, 50)))
	})
	t.Run("keeps the size and inverts Anchor", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, direction := range Directions() {
				aligned := b.AlignTo(direction, Pt(3.5, -2.25))

				geomtest.AssertSize(t, aligned.Size(), b.Size(), fmt.Sprintf("%s %s: ", b, direction))
				geomtest.AssertPoint(t, aligned.Anchor(direction), Pt(3.5, -2.25), fmt.Sprintf("%s %s: ", b, direction))
			}
		}
	})
}

func TestBox_Clamp(t *testing.T) {
	container := BoxFromMinMax(Pt(0, 0), Pt(10, 10))

	t.Run("a box within is unchanged", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(4, 4), Pt(6, 6)).Clamp(container), BoxFromMinMax(Pt(4, 4), Pt(6, 6)))
	})
	t.Run("a box outside moves in by the least", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(10, 4), Pt(14, 6)).Clamp(container), BoxFromMinMax(Pt(6, 4), Pt(10, 6)))
		geomtest.AssertBox(t, BoxFromMinMax(Pt(-4, -4), Pt(-2, -2)).Clamp(container), BoxFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("an axis larger than the other is centered", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(13, 2), Pt(27, 4)).Clamp(container), BoxFromMinMax(Pt(-2, 2), Pt(12, 4)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(9.5, 0.25), Pt(10.5, 1.25)).Clamp(container.Float()), BoxFromMinMax(Pt(9.0, 0.25), Pt(10.0, 1.25)))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, box := range boxFixtures {
			for _, other := range boxFixtures {
				clamped := box.Clamp(other)
				message := fmt.Sprintf("%s → %s: ", box, other)

				geomtest.AssertSize(t, clamped.Size(), box.Size(), message)
				geomtest.AssertBox(t, clamped.Clamp(other), clamped, message)
				if other.Contains(box.Min) && other.Contains(box.Max) {
					geomtest.AssertBox(t, clamped, box, message)
				}
				if box.Width() <= other.Width() && box.Height() <= other.Height() {
					assert.True(t, other.Contains(clamped.Min) && other.Contains(clamped.Max), message)
				}
			}
		}
	})
}

func TestBox_Inset(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(10, 10))

	t.Run("each edge moves in by its own padding", func(t *testing.T) {
		geomtest.AssertBox(t, box.Inset(PadU(1)), BoxFromMinMax(Pt(1, 1), Pt(9, 9)))
		geomtest.AssertBox(t, box.Inset(Pad(3, 1, 2, 5)), BoxFromMinMax(Pt(5, 3), Pt(9, 8)))
	})
	t.Run("padding beyond the size collapses at the opposite edge", func(t *testing.T) {
		geomtest.AssertBox(t, box.Inset(Pad(0, 0, 0, 20)), BoxFromMinMax(Pt(10, 0), Pt(10, 10)))
		geomtest.AssertBox(t, box.Inset(Pad(0, 20, 0, 0)), BoxFromMinMax(Pt(0, 0), Pt(0, 10)))
		geomtest.AssertBox(t, box.Inset(Pad(20, 0, 0, 0)), BoxFromMinMax(Pt(0, 10), Pt(10, 10)))
	})
	t.Run("both paddings over-running leaves the collapse at left and top", func(t *testing.T) {
		geomtest.AssertBox(t, box.Inset(PadU(20)), BoxFromMinMax(Pt(10, 10), Pt(10, 10)))
	})
	t.Run("negative padding grows the edge", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), BoxFromMinMax(Pt(1.0, 1.5), Pt(12.0, 10.0)))
	})
}

func TestBox_Outset(t *testing.T) {
	t.Run("each edge moves out by its own padding", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0, 0), Pt(10, 10)).Outset(Pad(3, 1, 2, 5)), BoxFromMinMax(Pt(-5, -3), Pt(11, 12)))
	})
	t.Run("undoes an inset", func(t *testing.T) {
		box := BoxFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		padding := Pad(1.0, 2.0, 3.0, 4.0)

		geomtest.AssertBox(t, box.Inset(padding).Outset(padding), box)
		geomtest.AssertBox(t, box.Outset(padding).Inset(padding), box)
	})
}

func TestBox_Lerp(t *testing.T) {
	a, b := BoxFromMinMax(Pt(0, 0), Pt(4, 2)), BoxFromMinMax(Pt(10, 10), Pt(12, 20))

	t.Run("the ends and the middle", func(t *testing.T) {
		geomtest.AssertBox(t, a.Lerp(b, 0), a)
		geomtest.AssertBox(t, a.Lerp(b, 1), b)
		geomtest.AssertBox(t, a.Lerp(b, 0.5), BoxFromMinMax(Pt(5, 5), Pt(8, 11)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, a.Float().Lerp(b.Float(), 0.25), BoxFromMinMax(Pt(2.5, 2.5), Pt(6.0, 6.5)))
	})
	t.Run("an extrapolation that turns the box inside out stays a box", func(t *testing.T) {
		geomtest.AssertBox(t, a.Lerp(b, 3), BoxFromMinMax(Pt(28, 30), Pt(30, 56)))
	})
}

func TestBox_Canonical(t *testing.T) {
	t.Run("orders the corners", func(t *testing.T) {
		geomtest.AssertBox(t, Box[int]{Pt(4, 0), Pt(0, 2)}.Canonical(), Box[int]{Pt(0, 0), Pt(4, 2)})
		geomtest.AssertBox(t, Box[float64]{Pt(4.5, 2.5), Pt(0.5, 0.5)}.Canonical(), Box[float64]{Pt(0.5, 0.5), Pt(4.5, 2.5)})
	})
	t.Run("a well-formed box is kept", func(t *testing.T) {
		for _, box := range boxFixtures {
			geomtest.AssertBox(t, box.Canonical(), box, box.String())
		}
	})
	t.Run("repairs a decoded box before Contains reads it", func(t *testing.T) {
		var box Box[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":4,"y":4},"b":{"x":0,"y":0}}`), &box))

		assert.False(t, box.Contains(Pt(2, 2)))
		assert.True(t, box.Canonical().Contains(Pt(2, 2)))
	})
}

func TestBox_DistanceTo(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("beside an edge measures to the edge", func(t *testing.T) {
		geomtest.AssertNumber(t, box.DistanceTo(Pt(5, 0)), 3.0)
		geomtest.AssertNumber(t, box.DistanceTo(Pt(0, -6)), 4.0)
	})
	t.Run("beyond a corner measures to the corner", func(t *testing.T) {
		geomtest.AssertNumber(t, box.DistanceTo(Pt(5, 6)), 5.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, box.DistanceTo(Pt(1, -1)), 0.0)
		assert.Equal(t, box.DistanceTo(Pt(2, 2)), 0.0)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(-1.0, -1.0), Pt(1.0, 1.0))

		assert.Equal(t, unit.DistanceTo(Pt(1.0+Delta/2, 0.0)), 0.0)
		geomtest.AssertNumber(t, unit.DistanceTo(Pt(1.0+2*Delta, 0.0)), 2*Delta)
	})
	t.Run("agrees with the rectangle of the box", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, b.DistanceTo(p), b.Rectangle().DistanceTo(p), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
	t.Run("a narrow integer box measures a gap wider than its range", func(t *testing.T) {
		geomtest.AssertNumber(t, Bx(Pt[int8](-100, -100), Pt[int8](-90, -90)).DistanceTo(Pt[int8](100, -95)), 190.0)
	})
}

func TestBox_DistanceSquaredTo(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("exact for an int box", func(t *testing.T) {
		assert.Equal(t, box.DistanceSquaredTo(Pt(5, 6)), 25.0)
		assert.Equal(t, box.DistanceSquaredTo(Pt(5, 0)), 9.0)
		assert.Equal(t, box.DistanceSquaredTo(Pt(1, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, b.DistanceSquaredTo(p), b.DistanceTo(p)*b.DistanceTo(p), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
}

func TestBox_DistanceToBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(2, 2))

	t.Run("the gap along an axis", func(t *testing.T) {
		geomtest.AssertNumber(t, box.DistanceToBox(BoxFromMinMax(Pt(5, 1), Pt(6, 3))), 3.0)
	})
	t.Run("the gap across a corner", func(t *testing.T) {
		geomtest.AssertNumber(t, box.DistanceToBox(BoxFromMinMax(Pt(5, 6), Pt(7, 8))), 5.0)
	})
	t.Run("zero for boxes that overlap or touch", func(t *testing.T) {
		geomtest.AssertNumber(t, box.DistanceToBox(BoxFromMinMax(Pt(1, 1), Pt(4, 4))), 0.0)
		geomtest.AssertNumber(t, box.DistanceToBox(BoxFromMinMax(Pt(2, 0), Pt(4, 2))), 0.0)
	})
	t.Run("zero exactly where the boxes intersect, and the same from either side", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, a.DistanceToBox(b) == 0, a.IntersectsBox(b), message)
				assert.Equal(t, a.DistanceToBox(b), b.DistanceToBox(a), message)
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		other := BoxFromMinMax(Pt(5, 6), Pt(7, 8))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = box.DistanceToBox(other) > 0
		}), 0)
	})
}

func TestBox_DistanceSquaredToBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(2, 2))

	geomtest.AssertNumber(t, box.DistanceSquaredToBox(BoxFromMinMax(Pt(5, 6), Pt(7, 8))), 25.0)
	geomtest.AssertNumber(t, box.DistanceSquaredToBox(BoxFromMinMax(Pt(1, 1), Pt(4, 4))), 0.0)
	assert.True(t, math.IsNaN(BoxFromMinMax(Pt(0.0, 0.0), Pt(2.0, 2.0)).DistanceSquaredToBox(BoxFromMinMax(Pt(math.NaN(), 0.0), Pt(7.0, 8.0)))))
}

func TestBox_Nearest(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 1), Pt(2, 4))

	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Nearest(Pt(1, 2)), Pt(1, 2))
	})
	t.Run("the point clamped onto the box", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Nearest(Pt(10, 10)), Pt(2, 4))
		geomtest.AssertPoint(t, box.Nearest(Pt(-3, 2)), Pt(0, 2))
		geomtest.AssertPoint(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).Nearest(Pt(-1.0, 1.2)), Pt(0.0, 1.2))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		assert.Equal(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)).Nearest(Pt(1.0+Delta/2, 0.5)), Pt(1.0+Delta/2, 0.5))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range pointFixtures {
				assertNearest[float64](t, b, p)
			}
		}
	})
}

func TestBox_Support(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("the corner reaching farthest", func(t *testing.T) {
		geomtest.AssertPoint(t, box.Support(Vec(1.0, 1.0)), Pt(4.0, 2.0))
		geomtest.AssertPoint(t, box.Support(Vec(-1.0, -1.0)), Pt(0.0, 0.0))
		geomtest.AssertPoint(t, box.Support(Vec(-1.0, 3.0)), Pt(0.0, 2.0))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		direction := Vec(tenth, 0.7)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = box.Support(direction).X > 0
		}), 0)
	})
}

func TestBox_Contains(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 1), Pt(2, 4))

	t.Run("inside", func(t *testing.T) {
		assert.True(t, box.Contains(Pt(1, 2)))
		assert.True(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).Contains(Pt(0.25, -0.75)))
	})
	t.Run("outside", func(t *testing.T) {
		assert.False(t, box.Contains(Pt(3, 0)))
		assert.False(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).Contains(Pt(-0.1, 0.0)))
	})
	t.Run("the boundary counts as inside", func(t *testing.T) {
		assert.True(t, box.Contains(box.Min))
		assert.True(t, box.Contains(box.Max))
		assert.True(t, box.Contains(Pt(2, 3)))
	})
	t.Run("float within Delta is inside, beyond it outside", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, unit.Contains(Pt(1+Delta/2, 0.5)))
		assert.False(t, unit.Contains(Pt(1+2*Delta, 0.5)))
	})
	t.Run("the tolerance is on the distance, not on each coordinate", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, unit.Contains(Pt(1+Delta/2, 1+Delta/2)))
		assert.False(t, unit.Contains(Pt(1+0.8*Delta, 1+0.8*Delta)))
	})
	t.Run("a NaN coordinate is contained by no box", func(t *testing.T) {
		assert.False(t, box.Float().Contains(Pt(math.NaN(), 2.0)))
		assert.False(t, box.Float().Contains(Pt(1.0, math.NaN())))
	})
	t.Run("agrees with the rectangle of the box", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, b.Contains(p), b.Rectangle().Contains(p), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
	t.Run("a narrow integer box does not contain a point a wrapped gap away", func(t *testing.T) {
		assert.False(t, Bx(Pt[int8](-100, 0), Pt[int8](-100, 0)).Contains(Pt[int8](28, 0)))
	})
	t.Run("every test of a point on the box agrees with it at the boundary", func(t *testing.T) {
		boxes := []Box[float32]{Bx(Pt[float32](3.5497787, 1.5453138), Pt[float32](7.1348896, 4.0049815)), Bx(Pt[float32](8.56645, 9.768216), Pt[float32](9.503749, 19.746952))}
		points := []Point[float32]{Pt[float32](3.5496886, 1.5452706), Pt[float32](8.566437, 9.768117)}

		for i, b := range boxes {
			p, contains := points[i], boxes[i].Contains(points[i])
			message := fmt.Sprintf("%s → %s: ", p, b)

			assert.Equal(t, b.IntersectsBox(Bx(p, p)), contains, message)
			assert.Equal(t, b.IntersectsSegment(Seg(p, p)), contains, message)
			assert.Equal(t, b.IntersectsRectangle(Rect(p, Sz[float32](0, 0), 0)), contains, message)
			assert.Equal(t, b.IntersectsPolygon(Pol([]Point[float32]{p})), contains, message)
			assert.Equal(t, b.ContainsSegment(Seg(p, p)), contains, message)
		}
	})
}

func TestBox_ContainsCircle(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(10, 6))

	t.Run("touching a side from inside counts", func(t *testing.T) {
		assert.True(t, box.ContainsCircle(Circ(Pt(3, 3), 3)))
		assert.False(t, box.ContainsCircle(Circ(Pt(3, 2), 3)))
	})
	t.Run("encloses exactly where it encloses the bounds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, b.ContainsCircle(c), b.ContainsBox(c.Bounds()), fmt.Sprintf("%s → %s: ", b, c))
			}
		}
	})
	t.Run("a center outside is not enclosed, a point on the boundary only as a point", func(t *testing.T) {
		assert.False(t, box.ContainsCircle(Circ(Pt(-1, 3), 0)))
		assert.True(t, box.ContainsCircle(Circ(Pt(0, 3), 0)))
		assert.False(t, box.ContainsCircle(Circ(Pt(0, 3), 1)))
	})
}

func TestBox_ContainsCapsule(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(10, 6))

	t.Run("touching the sides from inside counts", func(t *testing.T) {
		assert.True(t, box.ContainsCapsule(Cap(Pt(3, 3), Pt(7, 3), 3)))
		assert.False(t, box.ContainsCapsule(Cap(Pt(3, 3), Pt(8, 3), 3)))
	})
}

func TestBox_ContainsSector(t *testing.T) {
	t.Run("a sector within it, touching its sides from inside", func(t *testing.T) {
		assert.True(t, Bx(Pt(0.0, 0.0), Pt(5.0, 5.0)).ContainsSector(quarterSector()))
		assert.True(t, Bx(Pt(-4.0, 0.0), Pt(4.0, 5.0)).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
	})
	t.Run("a sector with a corner outside, or its arc alone", func(t *testing.T) {
		assert.False(t, Bx(Pt(0.0, 0.0), Pt(4.9, 5.0)).ContainsSector(quarterSector()))
		assert.False(t, Bx(Pt(-4.0, 0.0), Pt(4.0, 4.5)).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Bx(Pt(0, 0), Pt(5, 5)).ContainsSector(Sect(Pt(0, 0), 5, 0, Pi/2)))
		assert.False(t, Bx(Pt(0, 0), Pt(4, 5)).ContainsSector(Sect(Pt(0, 0), 5, 0, Pi/2)))
	})
}

func TestBox_ContainsSegment(t *testing.T) {
	t.Run("encloses exactly where it encloses the bounds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, b.ContainsSegment(s), b.ContainsBox(s.Bounds()), fmt.Sprintf("%s → %s: ", b, s))
			}
		}
	})
}

func TestBox_ContainsPolygon(t *testing.T) {
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, BoxFromMinMax(Pt(0, 0), Pt(4, 4)).ContainsPolygon(Pol[int](nil)))
	})
	t.Run("encloses exactly where it encloses the bounds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, b.ContainsPolygon(p), b.ContainsBox(p.Bounds()), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
}

func TestBox_ContainsRectangle(t *testing.T) {
	t.Run("encloses exactly where it encloses the bounds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, r := range rectFixtures {
				assert.Equal(t, b.ContainsRectangle(r), b.ContainsBox(r.Bounds()), fmt.Sprintf("%s → %s: ", b, r))
			}
		}
	})
}

func TestBox_ContainsRegularPolygon(t *testing.T) {
	t.Run("encloses exactly where it encloses the bounds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, b.ContainsRegularPolygon(rp), b.ContainsBox(rp.Bounds()), fmt.Sprintf("%s → %s: ", b, rp))
			}
		}
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, BoxFromMinMax(Pt(0, 0), Pt(4, 4)).ContainsRegularPolygon(RegPol(Pt(2, 2), SzU(1), 0, 0, 0)))
	})
	t.Run("inside and touching a side", func(t *testing.T) {
		box := BoxFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))

		assert.True(t, box.ContainsRegularPolygon(RegPol(Pt(5.0, 5.0), SzU(2.0), 6, 0, 0)))
		assert.True(t, box.ContainsRegularPolygon(RegPol(Pt(8.0, 5.0), SzU(2.0), 4, 0, 0)))
		assert.False(t, box.ContainsRegularPolygon(RegPol(Pt(9.0, 5.0), SzU(2.0), 4, 0, 0)))
	})
}

func TestBox_ContainsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 4))

	t.Run("inside and sharing an edge", func(t *testing.T) {
		assert.True(t, box.ContainsBox(BoxFromMinMax(Pt(1, 1), Pt(2, 2))))
		assert.True(t, box.ContainsBox(BoxFromMinMax(Pt(0, 1), Pt(4, 3))))
	})
	t.Run("overlapping and apart", func(t *testing.T) {
		assert.False(t, box.ContainsBox(BoxFromMinMax(Pt(2, 2), Pt(6, 6))))
		assert.False(t, box.ContainsBox(BoxFromMinMax(Pt(5, 5), Pt(6, 6))))
	})
	t.Run("float within the tolerance of a corner counts, beyond it not", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, unit.ContainsBox(BoxFromMinMax(Pt(0.5, 0.5), Pt(1+Delta/2, 1+Delta/2))))
		assert.False(t, unit.ContainsBox(BoxFromMinMax(Pt(0.5, 0.5), Pt(1+0.8*Delta, 1+0.8*Delta))))
	})
}

func TestBox_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, b.IntersectsCircle(c), c.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, c))
			}
		}
	})
}

func TestBox_IntersectsCapsule(t *testing.T) {
	t.Run("mirrors Capsule.IntersectsBox", func(t *testing.T) {
		for _, shape := range boxFixtures {
			for _, c := range capsuleFixtures {
				assert.Equal(t, shape.IntersectsCapsule(c), c.IntersectsBox(shape), fmt.Sprintf("%s → %s: ", shape, c))
			}
		}
	})
}

func TestBox_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, b.IntersectsSegment(s), s.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, s))
			}
		}
	})
}

func TestBox_IntersectsRay(t *testing.T) {
	t.Run("mirrors Ray.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, b.IntersectsRay(r), r.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, r))
			}
		}
	})
}

func TestBox_IntersectsLine(t *testing.T) {
	t.Run("mirrors Line.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, l := range lineFixtures {
				assert.Equal(t, b.IntersectsLine(l), l.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, l))
			}
		}
	})
}

func TestBox_IntersectsPolygon(t *testing.T) {
	t.Run("mirrors Polygon.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, b.IntersectsPolygon(p), p.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
}

func TestBox_IntersectsRectangle(t *testing.T) {
	t.Run("mirrors Rectangle.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, r := range rectFixtures {
				assert.Equal(t, b.IntersectsRectangle(r), r.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, r))
			}
		}
	})
}

func TestBox_IntersectsRegularPolygon(t *testing.T) {
	t.Run("mirrors RegularPolygon.IntersectsBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, b.IntersectsRegularPolygon(rp), rp.IntersectsBox(b), fmt.Sprintf("%s → %s: ", b, rp))
			}
		}
	})
}

func TestBox_IntersectsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, box.IntersectsBox(BoxFromMinMax(Pt(2, 2), Pt(6, 6))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, box.IntersectsBox(BoxFromMinMax(Pt(1, 1), Pt(2, 2))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, box.IntersectsBox(BoxFromMinMax(Pt(5, 0), Pt(8, 4))))
		assert.False(t, box.IntersectsBox(BoxFromMinMax(Pt(5, 5), Pt(8, 8))))
	})
	t.Run("a shared edge or corner counts as an intersection", func(t *testing.T) {
		assert.True(t, box.IntersectsBox(BoxFromMinMax(Pt(4, 1), Pt(8, 3))))
		assert.True(t, box.IntersectsBox(BoxFromMinMax(Pt(4, 4), Pt(8, 8))))
	})
	t.Run("float within the tolerance of a corner intersects, beyond it not", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, unit.IntersectsBox(BoxFromMinMax(Pt(1+Delta/2, 1+Delta/2), Pt(2.0, 2.0))))
		assert.False(t, unit.IntersectsBox(BoxFromMinMax(Pt(1+0.8*Delta, 1+0.8*Delta), Pt(2.0, 2.0))))
	})
	t.Run("agrees with the rectangles of the boxes", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, a.IntersectsBox(b), a.Rectangle().IntersectsRectangle(b.Rectangle()), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestBox_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, b.IntersectionSegment(s).AppendTo(nil), s.IntersectionBox(b).AppendTo(nil), fmt.Sprintf("%s → %s: ", b, s))
			}
		}
	})
}

func TestBox_IntersectionRay(t *testing.T) {
	t.Run("matches Ray.IntersectionBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, b.IntersectionRay(r).AppendTo(nil), r.IntersectionBox(b).AppendTo(nil), fmt.Sprintf("%s → %s: ", b, r))
			}
		}
	})
}

func TestBox_IntersectionLine(t *testing.T) {
	t.Run("matches Line.IntersectionBox", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, l := range lineFixtures {
				geomtest.AssertVertices(t, b.IntersectionLine(l).AppendTo(nil), l.IntersectionBox(b).AppendTo(nil), fmt.Sprintf("%s → %s: ", b, l))
			}
		}
	})
}

func TestBox_ClipBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		overlap, ok := box.ClipBox(BoxFromMinMax(Pt(2, 1), Pt(6, 3)))

		assert.True(t, ok)
		geomtest.AssertBox(t, overlap, BoxFromMinMax(Pt(2, 1), Pt(4, 3)))
	})
	t.Run("apart is false", func(t *testing.T) {
		_, ok := box.ClipBox(BoxFromMinMax(Pt(5, 0), Pt(8, 4)))

		assert.False(t, ok)
	})
	t.Run("touching intersect in a zero extent", func(t *testing.T) {
		overlap, ok := box.ClipBox(BoxFromMinMax(Pt(4, 1), Pt(8, 3)))

		assert.True(t, ok)
		geomtest.AssertBox(t, overlap, BoxFromMinMax(Pt(4, 1), Pt(4, 3)))
	})
	t.Run("meeting at the origin alone is the zero box, and true", func(t *testing.T) {
		overlap, ok := BoxFromMinMax(Pt(-2, -2), Pt(0, 0)).ClipBox(box)

		assert.True(t, ok)
		assert.Zero(t, overlap)
	})
	t.Run("a corner admitted by the tolerance is placed on the boundary", func(t *testing.T) {
		overlap, ok := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)).ClipBox(BoxFromMinMax(Pt(1+Delta/2, 0.0), Pt(2.0, 1.0)))

		assert.True(t, ok)
		geomtest.AssertBox(t, overlap, BoxFromMinMax(Pt(1+Delta/2, 0.0), Pt(1+Delta/2, 1.0)))
		assert.Equal(t, overlap.Width(), 0.0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				overlap, ok := a.ClipBox(b)
				reversed, reversedOk := b.ClipBox(a)
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, ok, a.IntersectsBox(b), message)
				assert.Equal(t, reversedOk, ok, message)
				if !ok {
					continue
				}

				geomtest.AssertBox(t, overlap, reversed, message)
				for _, corner := range []Point[float64]{overlap.Min, overlap.Max} {
					assert.True(t, a.Contains(corner) && b.Contains(corner), message)
				}
			}
		}
	})
}

func TestBox_Union(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		geomtest.AssertBox(t, box.Union(BoxFromMinMax(Pt(2, 2), Pt(6, 6))), BoxFromMinMax(Pt(0, 0), Pt(6, 6)))
	})
	t.Run("apart spans the gap", func(t *testing.T) {
		geomtest.AssertBox(t, box.Union(BoxFromMinMax(Pt(8, -2), Pt(10, 2))), BoxFromMinMax(Pt(0, -2), Pt(10, 4)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		geomtest.AssertBox(t, box.Union(BoxFromMinMax(Pt(1, 1), Pt(2, 2))), box)
	})
	t.Run("symmetric and contains both", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				union := a.Union(b)
				message := fmt.Sprintf("%s → %s: ", a, b)

				geomtest.AssertBox(t, union, b.Union(a), message)
				for _, corner := range []Point[float64]{a.Min, a.Max, b.Min, b.Max} {
					assert.True(t, union.Contains(corner), message)
				}
			}
		}
	})
}

func TestBox_Equal(t *testing.T) {
	box := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 2.0))

	t.Run("equal within the tolerance", func(t *testing.T) {
		assert.True(t, box.Equal(BoxFromMinMax(Pt(Delta/2, 0.0), Pt(1.0, 2.0))))
	})
	t.Run("a different corner", func(t *testing.T) {
		assert.False(t, box.Equal(BoxFromMinMax(Pt(0.5, 0.0), Pt(1.0, 2.0))))
		assert.False(t, box.Equal(BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 2.5))))
	})
}

func TestBox_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.True(t, Box[int]{}.IsZero())
		assert.True(t, BoxFromSize(Sz(0.0, 0.0)).IsZero())
	})
	t.Run("not zero", func(t *testing.T) {
		assert.False(t, BoxFromSize(Sz(1, 0)).IsZero())
		assert.False(t, BoxFromMinMax(Pt(1, 1), Pt(1, 1)).IsZero())
	})
}

func TestBox_Polygon(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	t.Run("the corners in the order of Vertices", func(t *testing.T) {
		geomtest.AssertPolygon(t, box.Polygon(), Pol([]Point[int]{Pt(1, 2), Pt(4, 2), Pt(4, 7), Pt(1, 7)}))
		geomtest.AssertPolygon(t, box.Polygon(), box.Rectangle().Polygon())
	})
	t.Run("allocates once", func(t *testing.T) {
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = box.Polygon().Points
		}), 1)
	})
}

func TestBox_Rectangle(t *testing.T) {
	t.Run("int keeps the corners", func(t *testing.T) {
		r := BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Rectangle()

		geomtest.AssertRectangle(t, r, Rect(Pt(2, 4), Sz(3, 5), 0))
		geomtest.AssertPoint(t, r.Min(), Pt(1, 2))
		geomtest.AssertPoint(t, r.Max(), Pt(4, 7))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRectangle(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Rectangle(), Rect(Pt(1.2, 1.55), Sz(1.2, 3.6), 0))
	})
	t.Run("is the inverse of Rectangle.Bounds", func(t *testing.T) {
		for _, box := range boxFixtures {
			geomtest.AssertBox(t, box.Rectangle().Bounds(), box, box.String())
		}
	})
	t.Run("an int box maps to the image rectangle of the same corners", func(t *testing.T) {
		box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))
		r := box.Rectangle().Rectangle()

		geomtest.AssertPoint(t, PointFromImage[int](r.Min), box.Min)
		geomtest.AssertPoint(t, PointFromImage[int](r.Max), box.Max)
	})
}

func TestBox_Cast(t *testing.T) {
	geomtest.AssertBox(t, BoxFromMinMax(Pt(0.4, -1.5), Pt(2.5, 3.6)).Cast[int](), BoxFromMinMax(Pt(0, -2), Pt(3, 4)))
}

func TestBox_Int(t *testing.T) {
	t.Run("rounds each corner", func(t *testing.T) {
		geomtest.AssertBox(t, BoxFromMinMax(Pt(0.4, -1.5), Pt(2.5, 3.6)).Int(), BoxFromMinMax(Pt(0, -2), Pt(3, 4)))
	})
	t.Run("the size can change as the box moves", func(t *testing.T) {
		box := BoxFromMin(Pt(0.3, 0.0), Sz(2.5, 1.0))

		assert.Equal(t, box.Int().Width(), 3)
		assert.Equal(t, box.Translate(Vec(0.3, 0.0)).Int().Width(), 2)
	})
}

func TestBox_Float(t *testing.T) {
	geomtest.AssertBox(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Float(), BoxFromMinMax(Pt(1.0, 2.0), Pt(4.0, 7.0)))
}

func TestBox_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, BoxFromMinMax(Pt(0, 1), Pt(2, 4)).String(), "(0,1)-(2,4)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).String(), "(0.00,-2.05)-(1.20,1.55)")
	})
	t.Run("prints the extent of a rectangle", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3), 0).Bounds().String(), "(0,1)-(2,4)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).String)
	})
}

func TestBox_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, BoxFromMinMax(Pt(0, 1), Pt(2, 4)), `{"a":{"x":0,"y":1},"b":{"x":2,"y":4}}`)

		var box Box[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":0,"y":1},"b":{"x":2,"y":4}}`), &box))
		geomtest.AssertBox(t, box, BoxFromMinMax(Pt(0, 1), Pt(2, 4)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, box := range boxFixtures {
			data, err := json.Marshal(box)
			assert.NoError(t, err)

			var decoded Box[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, box)
		}
	})
}

func TestBox_Properties(t *testing.T) {
	t.Run("the size spans the corners and the center lies between them", func(t *testing.T) {
		for _, box := range boxFixtures {
			geomtest.AssertPoint(t, box.Min.Add(box.Size().Vector()), box.Max, box.String())
			assert.True(t, box.Contains(box.Center()), box.String())
		}
	})
	t.Run("translate keeps the size", func(t *testing.T) {
		for _, box := range boxFixtures {
			for _, vector := range vectorFixtures {
				moved := box.Translate(vector)

				geomtest.AssertSize(t, moved.Size(), box.Size(), fmt.Sprintf("%s → %s: ", box, vector))
				geomtest.AssertPoint(t, moved.Min, box.Min.Add(vector), fmt.Sprintf("%s → %s: ", box, vector))
			}
		}
	})
}

func TestBox_Immutable(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	box.Translate(Vec(3, -2))
	box.Canonical()
	box.Inset(PadU(1))
	box.Outset(PadU(1))
	box.Clamp(BoxFromSize(Sz(2, 2)))

	geomtest.AssertBox(t, box, BoxFromMinMax(Pt(1, 2), Pt(4, 7)))
}

// boxFixtures span square, portrait, landscape, degenerate, and off-origin boxes.
var boxFixtures = []Box[float64]{
	{},
	BoxFromMinMax(Pt(-1.0, -1.0), Pt(1.0, 1.0)),
	BoxFromMinMax(Pt(0.0, 0.5), Pt(2.0, 3.5)),
	BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)),
	BoxFromMinMax(Pt(-7.0, -0.25), Pt(0.0, 0.75)),
	BoxFromMin(Pt(0.0, 0.0), Sz(10.0, 20.0)),
	BoxFromMinMax(Pt(-2.0, -4.0), Pt(6.0, 2.0)),
	BoxFromMinMax(Pt(3.0, 3.0), Pt(3.0, 5.0)),
}

func ExampleBx() {
	fmt.Println(Bx(Pt(0, 2), Pt(4, 0)))
	// Output: (0,0)-(4,2)
}

func ExampleBoxFromMinMax() {
	fmt.Println(BoxFromMinMax(Pt(4, 2), Pt(0, 0)))
	// Output: (0,0)-(4,2)
}

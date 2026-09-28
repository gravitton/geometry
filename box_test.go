package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestBox_Constructor(t *testing.T) {
	t.Run("from min", func(t *testing.T) {
		AssertBox(t, BoxFromMin(Pt(1, 2), Sz(4, 3)), Box[int]{Pt(1, 2), Pt(5, 5)})
		AssertBox(t, BoxFromMin(Pt(0.5, -1.0), Sz(1.5, 2.5)), Box[float64]{Pt(0.5, -1.0), Pt(2.0, 1.5)})
	})
	t.Run("from min and max", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(0, 0), Pt(4, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
		AssertBox(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 3.0)), Box[float64]{Pt(0.0, 0.0), Pt(1.0, 3.0)})
	})
	t.Run("from size at the origin", func(t *testing.T) {
		AssertBox(t, BoxFromSize(Sz(4, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
	})
	t.Run("corners in either order", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(4, 2), Pt(0, 0)), Box[int]{Pt(0, 0), Pt(4, 2)})
		AssertBox(t, BoxFromMinMax(Pt(4, 0), Pt(0, 2)), Box[int]{Pt(0, 0), Pt(4, 2)})
	})
	t.Run("a negative extent measures the other way from the corner", func(t *testing.T) {
		AssertBox(t, BoxFromMin(Pt(4, 2), Sz(-4, -2)), Box[int]{Pt(0, 0), Pt(4, 2)})
		AssertBox(t, BoxFromSize(Sz(-4, 2)), Box[int]{Pt(-4, 0), Pt(0, 2)})
	})
}

func TestBox_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Width(), 3)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Width(), 1.2)
	})
}

func TestBox_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Height(), 5)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Height(), 3.6)
	})
}

func TestBox_Size(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertSize(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Size(), Sz(3, 5))
	})
	t.Run("float", func(t *testing.T) {
		AssertSize(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Size(), Sz(1.2, 3.6))
	})
}

func TestBox_Center(t *testing.T) {
	t.Run("int truncates toward min", func(t *testing.T) {
		AssertPoint(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Center(), Pt(2, 4))
		AssertPoint(t, BoxFromMinMax(Pt(-4, -7), Pt(-1, -2)).Center(), Pt(-3, -5))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Center(), Pt(1.2, 1.55))
	})
}

func TestBox_Bounds(t *testing.T) {
	box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))

	AssertBox(t, box.Bounds(), box)
}

func TestBox_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Translate(Vec(3, -2)), BoxFromMinMax(Pt(4, 0), Pt(7, 5)))
	})
	t.Run("float", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Translate(Vec(100.1, -0.1)), BoxFromMinMax(Pt(100.7, -0.35), Pt(101.9, 3.25)))
	})
}

func TestBox_Canonical(t *testing.T) {
	t.Run("orders the corners", func(t *testing.T) {
		AssertBox(t, Box[int]{Pt(4, 0), Pt(0, 2)}.Canonical(), Box[int]{Pt(0, 0), Pt(4, 2)})
		AssertBox(t, Box[float64]{Pt(4.5, 2.5), Pt(0.5, 0.5)}.Canonical(), Box[float64]{Pt(0.5, 0.5), Pt(4.5, 2.5)})
	})
	t.Run("a well-formed box is kept", func(t *testing.T) {
		for _, box := range boxFixtures {
			AssertBox(t, box.Canonical(), box, box.String())
		}
	})
	t.Run("repairs a decoded box before Contains reads it", func(t *testing.T) {
		var box Box[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":4,"y":4},"b":{"x":0,"y":0}}`), &box))

		assert.False(t, box.Contains(Pt(2, 2)))
		assert.True(t, box.Canonical().Contains(Pt(2, 2)))
	})
}

func TestBox_Inset(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(10, 10))

	t.Run("each edge moves in by its own padding", func(t *testing.T) {
		AssertBox(t, box.Inset(PadU(1)), BoxFromMinMax(Pt(1, 1), Pt(9, 9)))
		AssertBox(t, box.Inset(Pad(3, 1, 2, 5)), BoxFromMinMax(Pt(5, 3), Pt(9, 8)))
	})
	t.Run("padding beyond the size collapses at the opposite edge", func(t *testing.T) {
		AssertBox(t, box.Inset(Pad(0, 0, 0, 20)), BoxFromMinMax(Pt(10, 0), Pt(10, 10)))
		AssertBox(t, box.Inset(Pad(0, 20, 0, 0)), BoxFromMinMax(Pt(0, 0), Pt(0, 10)))
		AssertBox(t, box.Inset(Pad(20, 0, 0, 0)), BoxFromMinMax(Pt(0, 10), Pt(10, 10)))
	})
	t.Run("both paddings over-running leaves the collapse at left and top", func(t *testing.T) {
		AssertBox(t, box.Inset(PadU(20)), BoxFromMinMax(Pt(10, 10), Pt(10, 10)))
	})
	t.Run("negative padding grows the edge", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0)).Inset(Pad(1.5, -2.0, 0.0, 1.0)), BoxFromMinMax(Pt(1.0, 1.5), Pt(12.0, 10.0)))
	})
}

func TestBox_Outset(t *testing.T) {
	t.Run("each edge moves out by its own padding", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(0, 0), Pt(10, 10)).Outset(Pad(3, 1, 2, 5)), BoxFromMinMax(Pt(-5, -3), Pt(11, 12)))
	})
	t.Run("undoes an inset", func(t *testing.T) {
		box := BoxFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		padding := Pad(1.0, 2.0, 3.0, 4.0)

		AssertBox(t, box.Inset(padding).Outset(padding), box)
		AssertBox(t, box.Outset(padding).Inset(padding), box)
	})
}

func TestBox_Clamp(t *testing.T) {
	container := BoxFromMinMax(Pt(0, 0), Pt(10, 10))

	t.Run("a box within is unchanged", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(4, 4), Pt(6, 6)).Clamp(container), BoxFromMinMax(Pt(4, 4), Pt(6, 6)))
	})
	t.Run("a box outside moves in by the least", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(10, 4), Pt(14, 6)).Clamp(container), BoxFromMinMax(Pt(6, 4), Pt(10, 6)))
		AssertBox(t, BoxFromMinMax(Pt(-4, -4), Pt(-2, -2)).Clamp(container), BoxFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("an axis larger than the other is centered", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(13, 2), Pt(27, 4)).Clamp(container), BoxFromMinMax(Pt(-2, 2), Pt(12, 4)))
	})
	t.Run("float", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(9.5, 0.25), Pt(10.5, 1.25)).Clamp(container.Float()), BoxFromMinMax(Pt(9.0, 0.25), Pt(10.0, 1.25)))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, box := range boxFixtures {
			for _, other := range boxFixtures {
				clamped := box.Clamp(other)
				message := fmt.Sprintf("%s → %s: ", box, other)

				AssertSize(t, clamped.Size(), box.Size(), message)
				AssertBox(t, clamped.Clamp(other), clamped, message)
				if other.Contains(box.Min) && other.Contains(box.Max) {
					AssertBox(t, clamped, box, message)
				}
				if box.Width() <= other.Width() && box.Height() <= other.Height() {
					assert.True(t, other.Contains(clamped.Min) && other.Contains(clamped.Max), message)
				}
			}
		}
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
}

func TestBox_DistanceTo(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("beside an edge measures to the edge", func(t *testing.T) {
		AssertNumber(t, box.DistanceTo(Pt(5, 0)), 3.0)
		AssertNumber(t, box.DistanceTo(Pt(0, -6)), 4.0)
	})
	t.Run("beyond a corner measures to the corner", func(t *testing.T) {
		AssertNumber(t, box.DistanceTo(Pt(5, 6)), 5.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, box.DistanceTo(Pt(1, -1)), 0.0)
		assert.Equal(t, box.DistanceTo(Pt(2, 2)), 0.0)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		unit := BoxFromMinMax(Pt(-1.0, -1.0), Pt(1.0, 1.0))

		assert.Equal(t, unit.DistanceTo(Pt(1.0+Delta/2, 0.0)), 0.0)
		AssertNumber(t, unit.DistanceTo(Pt(1.0+2*Delta, 0.0)), 2*Delta)
	})
	t.Run("agrees with the rectangle of the box", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, p := range pointFixtures {
				AssertNumber(t, b.DistanceTo(p), b.Rectangle().DistanceTo(p), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
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
				AssertNumber(t, b.DistanceSquaredTo(p), b.DistanceTo(p)*b.DistanceTo(p), fmt.Sprintf("%s → %s: ", b, p))
			}
		}
	})
}

func TestBox_Nearest(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 1), Pt(2, 4))

	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		AssertPoint(t, box.Nearest(Pt(1, 2)), Pt(1, 2))
	})
	t.Run("the point clamped onto the box", func(t *testing.T) {
		AssertPoint(t, box.Nearest(Pt(10, 10)), Pt(2, 4))
		AssertPoint(t, box.Nearest(Pt(-3, 2)), Pt(0, 2))
		AssertPoint(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).Nearest(Pt(-1.0, 1.2)), Pt(0.0, 1.2))
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

func TestBox_IntersectionBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 4))

	t.Run("overlapping", func(t *testing.T) {
		AssertBox(t, box.IntersectionBox(BoxFromMinMax(Pt(2, 1), Pt(6, 3))), BoxFromMinMax(Pt(2, 1), Pt(4, 3)))
	})
	t.Run("apart is the zero box", func(t *testing.T) {
		AssertBox(t, box.IntersectionBox(BoxFromMinMax(Pt(5, 0), Pt(8, 4))), Box[int]{})
	})
	t.Run("touching intersect in a zero extent", func(t *testing.T) {
		AssertBox(t, box.IntersectionBox(BoxFromMinMax(Pt(4, 1), Pt(8, 3))), BoxFromMinMax(Pt(4, 1), Pt(4, 3)))
	})
	t.Run("a corner admitted by the tolerance is placed on the boundary", func(t *testing.T) {
		overlap := BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)).IntersectionBox(BoxFromMinMax(Pt(1+Delta/2, 0.0), Pt(2.0, 1.0)))

		AssertBox(t, overlap, BoxFromMinMax(Pt(1+Delta/2, 0.0), Pt(1+Delta/2, 1.0)))
		assert.Equal(t, overlap.Width(), 0.0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				overlap := a.IntersectionBox(b)
				message := fmt.Sprintf("%s → %s: ", a, b)

				AssertBox(t, overlap, b.IntersectionBox(a), message)
				if !a.IntersectsBox(b) {
					assert.True(t, overlap.IsZero(), message)

					continue
				}

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
		AssertBox(t, box.Union(BoxFromMinMax(Pt(2, 2), Pt(6, 6))), BoxFromMinMax(Pt(0, 0), Pt(6, 6)))
	})
	t.Run("apart spans the gap", func(t *testing.T) {
		AssertBox(t, box.Union(BoxFromMinMax(Pt(8, -2), Pt(10, 2))), BoxFromMinMax(Pt(0, -2), Pt(10, 4)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		AssertBox(t, box.Union(BoxFromMinMax(Pt(1, 1), Pt(2, 2))), box)
	})
	t.Run("symmetric and contains both", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				union := a.Union(b)
				message := fmt.Sprintf("%s → %s: ", a, b)

				AssertBox(t, union, b.Union(a), message)
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

func TestBox_Rectangle(t *testing.T) {
	t.Run("int keeps the corners", func(t *testing.T) {
		r := BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Rectangle()

		AssertRectangle(t, r, Rect(Pt(2, 4), Sz(3, 5)))
		AssertPoint(t, r.Min(), Pt(1, 2))
		AssertPoint(t, r.Max(), Pt(4, 7))
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, BoxFromMinMax(Pt(0.6, -0.25), Pt(1.8, 3.35)).Rectangle(), Rect(Pt(1.2, 1.55), Sz(1.2, 3.6)))
	})
	t.Run("is the inverse of Rectangle.Bounds", func(t *testing.T) {
		for _, box := range boxFixtures {
			AssertBox(t, box.Rectangle().Bounds(), box, box.String())
		}
	})
	t.Run("an int box maps to the image rectangle of the same corners", func(t *testing.T) {
		box := BoxFromMinMax(Pt(1, 2), Pt(4, 7))
		r := box.Rectangle().Rectangle()

		AssertPoint(t, PointFromImage[int](r.Min), box.Min)
		AssertPoint(t, PointFromImage[int](r.Max), box.Max)
	})
}

func TestBox_Cast(t *testing.T) {
	AssertBox(t, BoxFromMinMax(Pt(0.4, -1.5), Pt(2.5, 3.6)).Cast[int](), BoxFromMinMax(Pt(0, -2), Pt(3, 4)))
}

func TestBox_Int(t *testing.T) {
	t.Run("rounds each corner", func(t *testing.T) {
		AssertBox(t, BoxFromMinMax(Pt(0.4, -1.5), Pt(2.5, 3.6)).Int(), BoxFromMinMax(Pt(0, -2), Pt(3, 4)))
	})
	t.Run("the size can change as the box moves", func(t *testing.T) {
		box := BoxFromMin(Pt(0.3, 0.0), Sz(2.5, 1.0))

		assert.Equal(t, box.Int().Width(), 3)
		assert.Equal(t, box.Translate(Vec(0.3, 0.0)).Int().Width(), 2)
	})
}

func TestBox_Float(t *testing.T) {
	AssertBox(t, BoxFromMinMax(Pt(1, 2), Pt(4, 7)).Float(), BoxFromMinMax(Pt(1.0, 2.0), Pt(4.0, 7.0)))
}

func TestBox_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, BoxFromMinMax(Pt(0, 1), Pt(2, 4)).String(), "(0,1)-(2,4)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, BoxFromMinMax(Pt(0.0, -2.05), Pt(1.2, 1.55)).String(), "(0.00,-2.05)-(1.20,1.55)")
	})
	t.Run("prints the extent of a rectangle", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Bounds().String(), "(0,1)-(2,4)")
	})
}

func TestBox_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, BoxFromMinMax(Pt(0, 1), Pt(2, 4)), `{"a":{"x":0,"y":1},"b":{"x":2,"y":4}}`)

		var box Box[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":0,"y":1},"b":{"x":2,"y":4}}`), &box))
		AssertBox(t, box, BoxFromMinMax(Pt(0, 1), Pt(2, 4)))
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
			AssertPoint(t, box.Min.Add(box.Size().Vector()), box.Max, box.String())
			assert.True(t, box.Contains(box.Center()), box.String())
		}
	})
	t.Run("translate keeps the size", func(t *testing.T) {
		for _, box := range boxFixtures {
			for _, vector := range vectorFixtures {
				moved := box.Translate(vector)

				AssertSize(t, moved.Size(), box.Size(), fmt.Sprintf("%s → %s: ", box, vector))
				AssertPoint(t, moved.Min, box.Min.Add(vector), fmt.Sprintf("%s → %s: ", box, vector))
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

	AssertBox(t, box, BoxFromMinMax(Pt(1, 2), Pt(4, 7)))
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

func ExampleBoxFromMinMax() {
	fmt.Println(BoxFromMinMax(Pt(4, 2), Pt(0, 0)))
	// Output: (0,0)-(4,2)
}

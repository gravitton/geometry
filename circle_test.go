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

func TestCircle_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(10, 16), 12), Circle[int]{Center: Pt(10, 16), Radius: 12})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.16, 204), 5.1), Circle[float64]{Center: Pt(0.16, 204.0), Radius: 5.1})
	})
	t.Run("a negative radius is taken absolute", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(10, 16), -12), Circ(Pt(10, 16), 12))
	})
}

func TestCircle_Anchor(t *testing.T) {
	c := Circ(Pt(10.0, 10.0), 5.0)

	t.Run("cardinal directions", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Anchor(Right), Pt(15.0, 10.0))
		geomtest.AssertPoint(t, c.Anchor(Left), Pt(5.0, 10.0))
		geomtest.AssertPoint(t, c.Anchor(Top), Pt(10.0, 5.0))
		geomtest.AssertPoint(t, c.Anchor(Bottom), Pt(10.0, 15.0))
	})
	t.Run("diagonals land on the boundary", func(t *testing.T) {
		// unlike Rectangle, whose diagonals reach the corners
		geomtest.AssertNumber(t, c.Center.DistanceTo(c.Anchor(DirectionUpRight)), c.Radius)
	})
	t.Run("none is the center", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Anchor(DirectionNone), Pt(10.0, 10.0))
	})
}

func TestCircle_Centroid(t *testing.T) {
	geomtest.AssertPoint(t, Circ(Pt(1, 2), 10).Centroid(), Pt(1, 2))
}

func TestCircle_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(1, 2), 10).Area(), Pi*100.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Area(), Pi*1.44)
	})
	t.Run("large integer radius does not overflow", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0, 0), 3037000500).Area(), Pi*3037000500*3037000500)
	})
}

func TestCircle_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(1, 2), 10).Perimeter(), Pi*20.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Perimeter(), Pi*2.4)
	})
}

func TestCircle_Inertia(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(1, 2), 10).Inertia(), Pi*5000.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Inertia(), Pi*1.0368)
	})
	t.Run("agrees with the ellipse and is approached by the polygon", func(t *testing.T) {
		for _, c := range circleFixtures {
			geomtest.AssertNumber(t, c.Inertia(), c.Ellipse().Inertia(), c.String())
			assert.EqualDelta(t, c.Inertia(), c.RegularPolygon(360, OrientationFlatTop).Inertia(), c.Inertia()*1e-3, c.String())
		}
	})
}

func TestCircle_Diameter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(1, 2), 10).Diameter(), 20)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Diameter(), 2.4)
	})
}

func TestCircle_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertBox(t, Circ(Pt(1, 2), 10).Bounds(), BoxFromMinMax(Pt(-9, -8), Pt(11, 12)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, Circ(Pt(0.6, -0.25), 1.2).Bounds(), BoxFromMinMax(Pt(-0.6, -1.45), Pt(1.8, 0.95)))
	})
}

func TestCircle_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Translate(Vec(3, -2)), Circ(Pt(4, 0), 10))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Translate(Vec(100.1, -0.1)), Circ(Pt(100.7, -0.35), 1.2))
	})
}

func TestCircle_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).MoveTo(Pt(3, -2)), Circ(Pt(3, -2), 10))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).MoveTo(Pt(100.1, -0.1)), Circ(Pt(100.1, -0.1), 1.2))
	})
}

func TestCircle_Scale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Scale(2.5), Circ(Pt(1, 2), 25))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Scale(2.5), Circ(Pt(0.6, -0.25), 3.0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Scale(-2), Circ(Pt(1, 2), 20))
	})
}

func TestCircle_Unscale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 25).Unscale(2.5), Circ(Pt(1, 2), 10))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 3.0).Unscale(2.5), Circ(Pt(0.6, -0.25), 1.2))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 20).Unscale(-2), Circ(Pt(1, 2), 10))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Circ(Pt(1, 2), 10).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestCircle_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Resize(8), Circ(Pt(1, 2), 8))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Resize(3.1), Circ(Pt(0.6, -0.25), 3.1))
	})
	t.Run("a negative radius is taken absolute", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Resize(-8), Circ(Pt(1, 2), 8))
	})
}

func TestCircle_Canonical(t *testing.T) {
	t.Run("takes a literal negative radius absolute and keeps the center", func(t *testing.T) {
		geomtest.AssertCircle(t, Circle[int]{Pt(1, 2), -8}.Canonical(), Circ(Pt(1, 2), 8))
		geomtest.AssertCircle(t, Circle[float64]{Pt(0.6, -0.25), -1.2}.Canonical(), Circ(Pt(0.6, -0.25), 1.2))
	})
	t.Run("is the circle Circ builds", func(t *testing.T) {
		c := Circle[int]{Pt(1, 2), -8}

		geomtest.AssertCircle(t, c.Canonical(), Circ(c.Center, c.Radius))
		assert.True(t, c.Canonical().Contains(c.Canonical().Anchor(Right)))
	})
	t.Run("repairs decoded JSON", func(t *testing.T) {
		var c Circle[int]

		assert.Nil(t, json.Unmarshal([]byte(`{"x":1,"y":2,"r":-8}`), &c))
		geomtest.AssertCircle(t, c.Canonical(), Circ(Pt(1, 2), 8))
	})
	t.Run("a well-formed circle is unchanged", func(t *testing.T) {
		for _, c := range circleFixtures {
			geomtest.AssertCircle(t, c.Canonical(), c, c.String())
		}
	})
}

func TestCircle_Grow(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Grow(8), Circ(Pt(1, 2), 18))
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Grow(-12), Circ(Pt(1, 2), 0))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Grow(3.1), Circ(Pt(0.6, -0.25), 4.3))
	})
}

func TestCircle_Shrink(t *testing.T) {
	t.Run("reduces the radius", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(8), Circ(Pt(1, 2), 2))
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Shrink(0.3), Circ(Pt(0.6, -0.25), 0.9))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(100), Circ(Pt(1, 2), 0))
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Shrink(5.0), Circ(Pt(0.6, -0.25), 0.0))
	})
}

func TestCircle_Lerp(t *testing.T) {
	a, b := Circ(Pt(0.0, 0.0), 2.0), Circ(Pt(10.0, 20.0), 8.0)

	t.Run("moves the center and the radius together", func(t *testing.T) {
		geomtest.AssertCircle(t, a.Lerp(b, 0.5), Circ(Pt(5.0, 10.0), 5.0))
	})
	t.Run("the ends are the circles themselves", func(t *testing.T) {
		geomtest.AssertCircle(t, a.Lerp(b, 0), a)
		geomtest.AssertCircle(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates and keeps the radius absolute", func(t *testing.T) {
		geomtest.AssertCircle(t, a.Lerp(b, 2), Circ(Pt(20.0, 40.0), 14.0))
		geomtest.AssertCircle(t, Circ(Pt(0.0, 0.0), 2.0).Lerp(Circ(Pt(0.0, 0.0), 0.0), 2), Circ(Pt(0.0, 0.0), 2.0))
	})
}

func TestCircle_Rotate(t *testing.T) {
	circle := Circ(Pt(1.0, 2.0), 3.0)

	t.Run("no angle moves a circle", func(t *testing.T) {
		geomtest.AssertCircle(t, circle.Rotate(Pi/3), circle)
		geomtest.AssertCircle(t, circle.Rotate(-Pi), circle)
	})
	t.Run("the bounds turn with nothing", func(t *testing.T) {
		for _, c := range circleFixtures {
			geomtest.AssertCircle(t, c.Rotate(Pi/7), c, c.String())
		}
	})
}

func TestCircle_AlignTo(t *testing.T) {
	c := Circ(Pt(10.0, 10.0), 5.0)

	t.Run("moves the anchor onto the point", func(t *testing.T) {
		geomtest.AssertCircle(t, c.AlignTo(Bottom, Pt(0.0, 0.0)), Circ(Pt(0.0, -5.0), 5.0))
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 3).AlignTo(Left, Pt(0, 0)), Circ(Pt(3, 0), 3))
	})
	t.Run("none aligns the center like MoveTo", func(t *testing.T) {
		geomtest.AssertCircle(t, c.AlignTo(DirectionNone, Pt(1.0, 2.0)), c.MoveTo(Pt(1.0, 2.0)))
	})
	t.Run("is the inverse of Anchor", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, direction := range Directions() {
				geomtest.AssertPoint(t, c.AlignTo(direction, Pt(1.5, -2.5)).Anchor(direction), Pt(1.5, -2.5), fmt.Sprintf("%s %s: ", c, direction))
			}
		}
	})
}

func TestCircle_Contains(t *testing.T) {
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Circ(Pt(1, 2), 10).Contains(Pt(4, 4)))
		assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Contains(Pt(0.1, 0.8)))
	})
	t.Run("outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(1, 2), 10).Contains(Pt(1, 13)))
		assert.False(t, Circ(Pt(0.6, -0.25), 1.2).Contains(Pt(0.0, 1.7)))
	})
	t.Run("the center and the boundary are inside", func(t *testing.T) {
		c := Circ(Pt(1, 2), 10)

		assert.True(t, c.Contains(c.Center))
		assert.True(t, c.Contains(c.Anchor(Right)))
		assert.False(t, c.Contains(c.Anchor(Right).AddXY(1, 0)))
	})
	t.Run("agrees with a zero-length segment at the point", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, c.Contains(p), Seg(p, p).IntersectsCircle(c), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
	})
	t.Run("a float anchor is inside despite rounding", func(t *testing.T) {
		c := Circ(Pt(0.1, 0.2), 0.7)

		for _, direction := range Directions() {
			assert.True(t, c.Contains(c.Anchor(direction)), direction.String())
		}
		assert.False(t, c.Contains(c.Anchor(Right).AddXY(2*Delta, 0)))
	})
}

func TestCircle_DistanceTo(t *testing.T) {
	circle := Circ(Pt(0, 0), 3)

	t.Run("outside measures to the boundary", func(t *testing.T) {
		geomtest.AssertNumber(t, circle.DistanceTo(Pt(7, 0)), 4.0)
		geomtest.AssertNumber(t, circle.DistanceTo(Pt(3, 4)), 2.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, circle.DistanceTo(Pt(1, 1)), 0.0)
		assert.Equal(t, circle.DistanceTo(Pt(3, 0)), 0.0)
		assert.Equal(t, circle.DistanceTo(circle.Center), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0.0, 0.0), 1.0).DistanceTo(Pt(1.0, 1.0)), Sqrt2-1)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		c := Circ(Pt(0.0, 0.0), 1.0)

		assert.Equal(t, c.DistanceTo(Pt(1.0+Delta/2, 0.0)), 0.0)
		geomtest.AssertNumber(t, c.DistanceTo(Pt(1.0+2*Delta, 0.0)), 2*Delta)
	})
	t.Run("zero exactly where Contains holds", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, c.DistanceTo(p) == 0, c.Contains(p), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
	})
}

func TestCircle_DistanceSquaredTo(t *testing.T) {
	circle := Circ(Pt(0, 0), 3)

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		geomtest.AssertNumber(t, circle.DistanceSquaredTo(Pt(7, 0)), 16.0)
		geomtest.AssertNumber(t, circle.DistanceSquaredTo(Pt(3, 4)), 4.0)
		assert.Equal(t, circle.DistanceSquaredTo(Pt(1, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, c.DistanceSquaredTo(p), c.DistanceTo(p)*c.DistanceTo(p), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
	})
}

func TestCircle_Nearest(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 3.0)

	t.Run("the boundary toward the point", func(t *testing.T) {
		geomtest.AssertPoint(t, circle.Nearest(Pt(7.0, 0.0)), Pt(3.0, 0.0))
		geomtest.AssertPoint(t, circle.Nearest(Pt(3.0, 4.0)), Pt(1.8, 2.4))
	})
	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		geomtest.AssertPoint(t, circle.Nearest(Pt(1.0, 1.0)), Pt(1.0, 1.0))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		assert.Equal(t, circle.Nearest(Pt(3.0+Delta/2, 0.0)), Pt(3.0+Delta/2, 0.0))
	})
	t.Run("int rounds once", func(t *testing.T) {
		geomtest.AssertPoint(t, Circ(Pt(0, 0), 2).Nearest(Pt(5, 1)), Pt(2, 0))
	})
	t.Run("a zero radius is its center", func(t *testing.T) {
		geomtest.AssertPoint(t, Circ(Pt(1.0, 1.0), 0.0).Nearest(Pt(5.0, 4.0)), Pt(1.0, 1.0))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				assertNearest[float64](t, c, p)
			}
		}
	})
}

func TestCircle_EnclosesCircle(t *testing.T) {
	circle := Circ(Pt(0, 0), 10)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.EnclosesCircle(Circ(Pt(2, 0), 5)))
		assert.True(t, circle.EnclosesCircle(Circ(Pt(0, 0), 0)))
	})
	t.Run("touching the boundary from inside counts", func(t *testing.T) {
		assert.True(t, circle.EnclosesCircle(Circ(Pt(5, 0), 5)))
		assert.False(t, circle.EnclosesCircle(Circ(Pt(6, 0), 5)))
	})
	t.Run("larger or apart", func(t *testing.T) {
		assert.False(t, circle.EnclosesCircle(Circ(Pt(0, 0), 11)))
		assert.False(t, circle.EnclosesCircle(Circ(Pt(30, 0), 1)))
	})
	t.Run("float within the tolerance of the boundary counts, beyond it not", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.0), 10.0).EnclosesCircle(Circ(Pt(5.0+Delta/2, 0.0), 5.0)))
		assert.False(t, Circ(Pt(0.0, 0.0), 10.0).EnclosesCircle(Circ(Pt(5.0+2*Delta, 0.0), 5.0)))
	})
	t.Run("encloses exactly where it encloses the far point", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				far := b.Center.AddXY(b.Radius, 0)
				if direction := b.Center.Subtract(a.Center); direction != (Vector[float64]{}) {
					far = b.Center.Add(direction.Resize(b.Radius))
				}

				assert.Equal(t, a.EnclosesCircle(b), a.Contains(far), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestCircle_EnclosesSegment(t *testing.T) {
	circle := Circ(Pt(0, 0), 10)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.EnclosesSegment(Seg(Pt(-5, 0), Pt(5, 5))))
	})
	t.Run("a chord counts", func(t *testing.T) {
		assert.True(t, circle.EnclosesSegment(Seg(Pt(-10, 0), Pt(10, 0))))
		assert.True(t, circle.EnclosesSegment(Seg(Pt(-6, 8), Pt(8, -6))))
	})
	t.Run("an endpoint outside", func(t *testing.T) {
		assert.False(t, circle.EnclosesSegment(Seg(Pt(0, 0), Pt(11, 0))))
	})
}

func TestCircle_EnclosesPolygon(t *testing.T) {
	circle := Circ(Pt(1, 1), 2)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.EnclosesPolygon(Pol(squareVertices())))
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(1, 1), 1).EnclosesPolygon(Pol(squareVertices())))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, circle.EnclosesPolygon(Pol[int](nil)))
	})
}

func TestCircle_EnclosesRectangle(t *testing.T) {
	t.Run("corners on the boundary count", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 5).EnclosesRectangle(Rect(Pt(0, 0), Sz(6, 8))))
		assert.False(t, Circ(Pt(0, 0), 5).EnclosesRectangle(Rect(Pt(0, 0), Sz(6, 9))))
	})
	t.Run("rotated is tested on its turned corners", func(t *testing.T) {
		long := Rect(Pt(0.0, 0.0), Sz(9.0, 1.0))

		assert.False(t, Circ(Pt(0.0, 0.0), 4.0).EnclosesRectangle(long))
		assert.False(t, Circ(Pt(0.0, 0.0), 4.0).EnclosesRectangle(long.Rotate(Pi/4)))
		assert.True(t, Circ(Pt(0.0, 0.0), 5.0).EnclosesRectangle(long.Rotate(Pi/4)))
	})
}

func TestCircle_EnclosesRegularPolygon(t *testing.T) {
	t.Run("a circle encloses its inscribed polygon", func(t *testing.T) {
		for _, c := range circleFixtures {
			assert.True(t, c.EnclosesRegularPolygon(c.RegularPolygon(7, OrientationFlatTop)), c.String())
		}
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 2).EnclosesRegularPolygon(RegPol(Pt(1, 0), Sz(2, 2), 4, 0, 0)))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 2).EnclosesRegularPolygon(RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0)))
	})
}

func TestCircle_EnclosesBox(t *testing.T) {
	t.Run("an inscribed box counts, the bounds of the circle not", func(t *testing.T) {
		c := Circ(Pt(0, 0), 5)

		assert.True(t, c.EnclosesBox(BoxFromMinMax(Pt(-3, -4), Pt(3, 4))))
		assert.False(t, c.EnclosesBox(c.Bounds()))
	})
}

func TestCircle_IntersectsCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 100.0)

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, circle.IntersectsCircle(Circ(Pt(199.0, 0.0), 100.0)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, circle.IntersectsCircle(Circ(Pt(210.0, 0.0), 100.0)))
	})
	t.Run("exactly touching counts as an intersection", func(t *testing.T) {
		assert.True(t, circle.IntersectsCircle(Circ(Pt(200.0, 0.0), 100.0)))
		assert.False(t, circle.IntersectsCircle(Circ(Pt(201.0, 0.0), 100.0)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, circle.IntersectsCircle(Circ(Pt(0.0, 0.0), 50.0)))
	})
	t.Run("narrow integers do not overflow the radii sum", func(t *testing.T) {
		assert.True(t, Circ(Pt[int8](0, 0), 100).IntersectsCircle(Circ(Pt[int8](0, 50), 100)))
		assert.False(t, Circ(Pt[int8](-20, 0), 50).IntersectsCircle(Circ(Pt[int8](100, 0), 50)))
	})
	t.Run("holds wherever Intersection finds a point", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				if len(a.IntersectionCircle(b)) > 0 {
					assert.True(t, a.IntersectsCircle(b), fmt.Sprintf("%s → %s: ", a, b))
				}
			}
		}
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.Equal(t, a.IntersectsCircle(b), b.IntersectsCircle(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestCircle_IntersectionCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 5.0)

	t.Run("overlapping gives two points mirrored across the centers", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(6.0, 0.0), 5.0)), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, -4.0)})
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(0.0, 6.0), 5.0)), []Point[float64]{Pt(-4.0, 3.0), Pt(4.0, 3.0)})
	})
	t.Run("a larger circle behind the first", func(t *testing.T) {
		height := math.Sqrt(1 - 0.125*0.125)

		geomtest.AssertVertices(t, Circ(Pt(0.0, 0.0), 1.0).IntersectionCircle(Circ(Pt(1.0, 0.0), 1.5)), []Point[float64]{Pt(-0.125, height), Pt(-0.125, -height)})
	})
	t.Run("tangent from outside gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(8.0, 0.0), 3.0)), []Point[float64]{Pt(5.0, 0.0)})
	})
	t.Run("tangent from inside gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(2.0, 0.0), 3.0)), []Point[float64]{Pt(5.0, 0.0)})
		geomtest.AssertVertices(t, Circ(Pt(0.0, 0.0), 3.0).IntersectionCircle(Circ(Pt(2.0, 0.0), 5.0)), []Point[float64]{Pt(-3.0, 0.0)})
	})
	t.Run("apart, nested and concentric give none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(20.0, 0.0), 5.0)))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(1.0, 0.0), 1.0)))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(0.0, 0.0), 3.0)))
		assert.Nil(t, circle.IntersectionCircle(circle))
	})
	t.Run("float treats centers within Epsilon as coincident", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(1e-7, 0.0), 5.0)))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(0.0, -1e-7), 3.0)))
	})
	t.Run("float is tolerant at a tangent", func(t *testing.T) {
		assert.Length(t, circle.IntersectionCircle(Circ(Pt(8.0+Delta/2, 0.0), 3.0)), 1)
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(8.0+2*Delta, 0.0), 3.0)))
	})
	t.Run("a tangent exactly Delta outside is judged like Intersects, on both sides", func(t *testing.T) {
		for _, other := range []Circle[float64]{Circ(Pt(8.0+Delta, 0.0), 3.0), Circ(Pt(2.0+Delta, 0.0), 3.0)} {
			assert.True(t, circle.IntersectsCircle(other))
			geomtest.AssertVertices(t, circle.IntersectionCircle(other), []Point[float64]{Pt(5.0, 0.0)})
		}
	})
	t.Run("a tangent within the tolerance lands on both boundaries", func(t *testing.T) {
		a, b := Circ(Pt(0.0, 0.0), 6.0), Circ(Pt(2.0000005, 0.0), 4.0)
		points := a.IntersectionCircle(b)

		assert.Length(t, points, 1)
		for _, p := range points {
			assert.True(t, touches(a, p), "on a")
			assert.True(t, touches(b, p), "on b")
		}
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)), []Point[int]{Pt(3, 4), Pt(3, -4)})
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 2).IntersectionCircle(Circ(Pt(3, 0), 2)), []Point[int]{Pt(2, 1), Pt(2, -1)})
	})
	t.Run("int gives two crossings rounding to one point once", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(11, -20), 15).IntersectionCircle(Circ(Pt(-3, -19), 1)), []Point[int]{Pt(-4, -19)})
	})
	t.Run("points lie on both circles and mirror Intersects", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				points := a.IntersectionCircle(b)

				assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two points: ", a, b))
				if len(points) > 0 {
					assert.True(t, a.IntersectsCircle(b), fmt.Sprintf("%s → %s: ", a, b))
				}
				for _, p := range points {
					assert.True(t, touches(a, p), fmt.Sprintf("%s → %s: %s on a: ", a, b, p))
					assert.True(t, touches(b, p), fmt.Sprintf("%s → %s: %s on b: ", a, b, p))
				}
			}
		}
	})
	t.Run("far from the origin float32 points lie on both circles", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, a := range circleFixtures {
				for _, b := range circleFixtures {
					a, b := a.Cast[float32]().Translate(offset), b.Cast[float32]().Translate(offset)

					for _, p := range a.IntersectionCircle(b) {
						assert.True(t, touches(a, p), fmt.Sprintf("%s → %s: %s on a: ", a, b, p))
						assert.True(t, touches(b, p), fmt.Sprintf("%s → %s: %s on b: ", a, b, p))
					}
				}
			}
		}
	})
}

func BenchmarkCircle_IntersectionCircle(b *testing.B) {
	circle, other := Circ(Pt(0.0, 0.0), 5.0), Circ(Pt(6.0, 0.0), 5.0)

	for b.Loop() {
		_ = circle.IntersectionCircle(other)
	}
}

func TestCircle_AppendIntersectionCircle(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionCircle([]Point[int]{Pt(3, 4)}, Circ(Pt(6, 0), 5)), []Point[int]{Pt(3, 4), Pt(3, 4), Pt(3, -4)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionCircle([]Point[int]{Pt(9, 9)}, Circ(Pt(20, 0), 5)), []Point[int]{Pt(9, 9)})
		assert.Nil(t, circle.AppendIntersectionCircle(nil, Circ(Pt(20, 0), 5)))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.AppendIntersectionCircle(buffer[:0], Circ(Pt(6, 0), 5))
		}), 0)
	})
	t.Run("matches IntersectionCircle after the points in dst", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, other := range circleFixtures {
				geomtest.AssertVertices(t, c.AppendIntersectionCircle(bufferWith(prefixPoint), other), append([]Point[float64]{prefixPoint}, c.IntersectionCircle(other)...), fmt.Sprintf("%s → %s: ", c, other))
			}
		}
	})
}

func FuzzCircle_IntersectionCircle(f *testing.F) {
	f.Add(0.0, 0.0, 1.0, 2.0, 0.0, 1.0)
	f.Add(0.0, 0.0, 1.0, 2.0+Delta/2, 0.0, 1.0)
	f.Add(0.0, 0.0, 2.0, 0.0, 0.0, 1.0)
	f.Add(0.0, 0.0, 6.0, 2.0000005, 0.0, 4.0)
	f.Add(11.0, -20.0, 15.0, -3.0, -19.0, 1.0)

	f.Fuzz(func(t *testing.T, x1, y1, r1, x2, y2, r2 float64) {
		for _, v := range []float64{x1, y1, r1, x2, y2, r2} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		a, b := Circ(Pt(x1, y1), r1), Circ(Pt(x2, y2), r2)
		points := a.IntersectionCircle(b)

		assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings, got %d: ", a, b, len(points)))

		for _, p := range points {
			assert.True(t, touches(a, p), fmt.Sprintf("%s → %s: %s on the first: ", a, b, p))
			assert.True(t, touches(b, p), fmt.Sprintf("%s → %s: %s on the second: ", a, b, p))
		}

		if len(points) == 2 {
			assert.True(t, !points[0].Equal(points[1]), fmt.Sprintf("%s → %s: two distinct points: ", a, b))
		}

		if len(points) > 0 {
			assert.True(t, a.IntersectsCircle(b), fmt.Sprintf("%s → %s: points imply intersects: ", a, b))
		}

		if rounded := a.Int().IntersectionCircle(b.Int()); len(rounded) == 2 {
			assert.True(t, !rounded[0].Equal(rounded[1]), fmt.Sprintf("%s → %s: two distinct rounded points: ", a.Int(), b.Int()))
		}
	})
}

func TestCircle_IntersectsSegment(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, circle.IntersectsSegment(Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, circle.IntersectsSegment(Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0))))
		assert.False(t, circle.IntersectsSegment(Seg(Pt(2.0, 0.0), Pt(3.0, 0.0))))
	})
	t.Run("tangent counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsSegment(Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0))))
		assert.False(t, circle.IntersectsSegment(Seg(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta))))
	})
	t.Run("an endpoint inside counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsSegment(Seg(Pt(0.5, 0.0), Pt(5.0, 0.0))))
	})
	t.Run("a segment inside counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsSegment(Seg(Pt(-0.5, 0.0), Pt(0.5, 0.0))))
	})
	t.Run("an int tangent counts across the supported range", func(t *testing.T) {
		for radius := 100; radius < 60000; radius += 211 {
			for length := 100; length < 60000; length += 911 {
				c := Circ(Pt(0, 0), radius)
				wall := Seg(Pt(radius, -length), Pt(radius, length/3))
				floor := Bx(Pt(-length, radius), Pt(length/3, radius+40))
				message := fmt.Sprintf("%s → %s: ", c, wall)

				assert.True(t, c.IntersectsSegment(wall), message)
				assert.Length(t, c.IntersectionSegment(wall), 1, message)
				assert.True(t, c.IntersectsBox(floor), message)
				assert.True(t, c.IntersectsRectangle(floor.Rectangle()), message)
				assert.True(t, c.IntersectsPolygon(floor.Rectangle().Polygon()), message)
			}
		}
		for k := 1; k < 4000000; k += 9973 {
			c := Circ(Pt(0, 0), 5*k)
			slope := Seg(Pt(-k, 7*k), Pt(7*k, k))

			assert.True(t, c.IntersectsSegment(slope), fmt.Sprintf("%s → %s: ", c, slope))
		}
	})
}

func TestCircle_IntersectionSegment(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(2.0, 0.0), Pt(-2.0, 0.0))), []Point[float64]{Pt(1.0, 0.0), Pt(-1.0, 0.0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(0.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(0.5, 0.0), Pt(5.0, 0.0))), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("allocates once for the result and not at all for none", func(t *testing.T) {
		unit := Circ(Pt(0, 0), 5)
		through, apart := Seg(Pt(-10, 0), Pt(10, 0)), Seg(Pt(-10, 9), Pt(10, 9))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = unit.IntersectionSegment(through)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = unit.IntersectionSegment(apart)
		}), 0)
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0))), []Point[float64]{Pt(0.0, 1.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0+Delta/2), Pt(2.0, 1.0+Delta/2))), []Point[float64]{Pt(0.0, 1.0+Delta/2)})
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta))))
	})
	t.Run("a chord within the tolerance band of the gap keeps both ends", func(t *testing.T) {
		height := 1.0 - Delta/2
		half := math.Sqrt(1 - height*height)

		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, height), Pt(1.0, height))), []Point[float64]{Pt(-half, height), Pt(half, height)})
	})
	t.Run("a chord shorter than Delta is a tangent", func(t *testing.T) {
		height := math.Sqrt(1 - Delta*Delta/16)

		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, height), Pt(1.0, height))), []Point[float64]{Pt(0.0, height)})
	})
	t.Run("a segment within the tolerance with both ends on the boundary gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(1.0, 0.0), Pt(1.0, Delta/10))), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("a tangent beyond the segment is missed", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(1.0, 1.0), Pt(2.0, 1.0))))
	})
	t.Run("a shallow touch is decided on the endpoint distance, like IntersectsSegment", func(t *testing.T) {
		shallow := Seg(Pt(-1.0, 1.0+Delta/2), Pt(-0.0005, 1.0+Delta/2))

		assert.True(t, circle.IntersectsSegment(shallow))
		geomtest.AssertVertices(t, circle.IntersectionSegment(shallow), []Point[float64]{shallow.End})
		geomtest.AssertVertices(t, circle.IntersectionSegment(shallow.Reverse()), []Point[float64]{shallow.End})
	})
	t.Run("an interior graze exactly Delta outside is judged like IntersectsSegment", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 1.000001), Pt(1.0, 1.000001)), Circ(Pt(0.0, 0.0), 1.0)

		assert.Equal(t, len(c.IntersectionSegment(s)) > 0, c.IntersectsSegment(s))
		for _, y := range []float64{1.0000009, 1.0000011, 1.0000015} {
			s := Seg(Pt(-1.0, y), Pt(1.0, y))

			assert.Equal(t, len(c.IntersectionSegment(s)) > 0, c.IntersectsSegment(s), s.String())
		}
	})
	t.Run("an endpoint exactly Delta outside is judged like IntersectsSegment", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 2.000001), Pt(-882.0315, 56.11111116666667)), Circ(Pt(-1.0, 0.0), 2.0)

		assert.Equal(t, len(c.IntersectionSegment(s)) > 0, c.IntersectsSegment(s))
	})
	t.Run("an endpoint on the boundary is counted once with its crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, 0.0), Pt(2.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(1.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
	})
	t.Run("an endpoint on the boundary before the chord is where the segment enters", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 2.000001), Pt(-877.0315, 0.11111116666666668)), Circ(Pt(-1.0, 0.0), 2.0)
		points := c.IntersectionSegment(s)

		assert.Length(t, points, 2)
		geomtest.AssertPoint(t, points[0], s.Start)
	})
	t.Run("an endpoint in the middle of a chord keeps the crossing ahead of it, whichever way the segment runs", func(t *testing.T) {
		height := 1.0 - Delta/10
		half := math.Sqrt(1 - height*height)
		s := Seg(Pt(0.0, height), Pt(5.0, height))

		geomtest.AssertVertices(t, circle.IntersectionSegment(s), []Point[float64]{s.Start, Pt(half, height)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(s.Reverse()), []Point[float64]{Pt(half, height), s.Start})
	})
	t.Run("an endpoint on the boundary of a segment leaving the circle is the one crossing", func(t *testing.T) {
		leaving := Seg(Pt(1.0, 0.0), Pt(2.0, 0.0))

		geomtest.AssertVertices(t, circle.IntersectionSegment(leaving), []Point[float64]{leaving.Start})
		geomtest.AssertVertices(t, circle.IntersectionSegment(leaving.Reverse()), []Point[float64]{leaving.Start})
	})
	t.Run("a tangent segment with both ends on the boundary gives its ends", func(t *testing.T) {
		grazing := Seg(Pt(-0.0003, 1.0), Pt(0.0003, 1.0))

		geomtest.AssertVertices(t, circle.IntersectionSegment(grazing), []Point[float64]{grazing.Start, grazing.End})
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0))))
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(2.0, 0.0), Pt(3.0, 0.0))))
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-0.5, 0.0), Pt(0.5, 0.0))))
	})
	t.Run("a degenerate segment is a point on the boundary or nothing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(1.0, 0.0), Pt(1.0, 0.0))), []Point[float64]{Pt(1.0, 0.0)})
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(0.5, 0.0), Pt(0.5, 0.0))))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionSegment(Seg(Pt(-5, -5), Pt(5, 5))), []Point[int]{Pt(-4, -4), Pt(4, 4)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsSegment holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				points := c.IntersectionSegment(s)

				assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings: ", s, c))
				for _, p := range points {
					assert.True(t, s.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", s, c, p))
					assert.True(t, touches(c, p), fmt.Sprintf("%s → %s: %s on the boundary: ", s, c, p))
				}
				if len(points) > 0 {
					assert.True(t, c.IntersectsSegment(s), fmt.Sprintf("%s → %s: ", s, c))
				}
			}
		}
	})
	t.Run("a segment and its Reverse cross at the same points", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range append(grazingSegments(c), segmentFixtures...) {
				assertReversed(t, c.IntersectionSegment(s), c.IntersectionSegment(s.Reverse()), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("far from the origin a float32 segment and its Reverse cross as often", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, c := range circleFixtures {
				for _, s := range append(grazingSegments(c), segmentFixtures...) {
					s, c := s.Cast[float32]().Translate(offset), c.Cast[float32]().Translate(offset)

					assert.Length(t, c.IntersectionSegment(s.Reverse()), len(c.IntersectionSegment(s)), fmt.Sprintf("%s → %s: ", s, c))
				}
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, c := range circleFixtures {
					s, c := s.Cast[float32]().Translate(offset), c.Cast[float32]().Translate(offset)
					points := c.IntersectionSegment(s)

					assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings: ", s, c))
					for _, p := range points {
						assert.True(t, s.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", s, c, p))
						assert.True(t, touches(c, p), fmt.Sprintf("%s → %s: %s on the boundary: ", s, c, p))
					}
					assert.True(t, len(points) == 0 || c.IntersectsSegment(s), fmt.Sprintf("%s → %s: ", s, c))
				}
			}
		}
	})
}

func TestCircle_AppendIntersectionSegment(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)
	through := Seg(Pt(-10, 0), Pt(10, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionSegment([]Point[int]{Pt(5, 0)}, through), []Point[int]{Pt(5, 0), Pt(-5, 0), Pt(5, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionSegment([]Point[int]{Pt(9, 9)}, Seg(Pt(-10, 9), Pt(10, 9))), []Point[int]{Pt(9, 9)})
		assert.Nil(t, circle.AppendIntersectionSegment(nil, Seg(Pt(-10, 9), Pt(10, 9))))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.AppendIntersectionSegment(buffer[:0], through)
		}), 0)
	})
	t.Run("matches IntersectionSegment after the points in dst", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, c.AppendIntersectionSegment(bufferWith(prefixPoint), s), append([]Point[float64]{prefixPoint}, c.IntersectionSegment(s)...), fmt.Sprintf("%s → %s: ", c, s))
			}
		}
	})
}

func FuzzCircle_IntersectionSegment(f *testing.F) {
	f.Add(-2.0, 0.0, 2.0, 0.0, 0.0, 0.0, 1.0)
	f.Add(-1.0, 1.0+Delta/2, -0.0005, 1.0+Delta/2, 0.0, 0.0, 1.0)
	f.Add(-0.5, 0.0, 0.5, 0.0, 0.0, 0.0, 1.0)
	f.Add(1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0)

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, cx, cy, r float64) {
		for _, v := range []float64{x1, y1, x2, y2, cx, cy, r} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		s, c := Seg(Pt(x1, y1), Pt(x2, y2)), Circ(Pt(cx, cy), r)
		points := c.IntersectionSegment(s)

		assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings, got %d: ", s, c, len(points)))

		for _, p := range points {
			assert.True(t, s.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", s, c, p))
			assert.True(t, touches(c, p), fmt.Sprintf("%s → %s: %s on the boundary: ", s, c, p))
		}

		if inside := c.Contains(s.Start) && c.Contains(s.End); !inside {
			assert.Equal(t, len(points) > 0, c.IntersectsSegment(s), fmt.Sprintf("%s → %s: ", s, c))
		}
	})
}

func TestCircle_IntersectsRay(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, circle.IntersectsRay(Ry(Pt(-2.0, 0.0), Vec(1.0, 0.0))))
	})
	t.Run("pointing away", func(t *testing.T) {
		assert.False(t, circle.IntersectsRay(Ry(Pt(-2.0, 0.0), Vec(-1.0, 0.0))))
	})
	t.Run("tangent counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsRay(Ry(Pt(-2.0, 1.0), Vec(1.0, 0.0))))
		assert.False(t, circle.IntersectsRay(Ry(Pt(-2.0, 1.0+2*Delta), Vec(1.0, 0.0))))
	})
	t.Run("a tangent behind the origin is missed", func(t *testing.T) {
		assert.False(t, circle.IntersectsRay(Ry(Pt(1.5, 1.0), Vec(1.0, 0.0))))
	})
	t.Run("the origin inside counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsRay(Ry(Pt(0.5, 0.0), Vec(0.0, -3.0))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, circle.IntersectsRay(Ry(Pt(0.5, 0.0), Vec(0.0, 0.0))))
		assert.False(t, circle.IntersectsRay(Ry(Pt(2.0, 0.0), Vec(0.0, 0.0))))
	})
	t.Run("matches a segment reaching past the circle along the ray", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, c.IntersectsRay(r), c.IntersectsSegment(far(r)), fmt.Sprintf("%s → %s: ", c, r))
			}
		}
	})
}

func TestCircle_IntersectionRay(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(-10, 0), Vec(1, 0))), []Point[int]{Pt(-5, 0), Pt(5, 0)})
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(10, 0), Vec(-2, 0))), []Point[int]{Pt(5, 0), Pt(-5, 0)})
	})
	t.Run("starting inside gives the exit", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(0, 0), Vec(0, 1))), []Point[int]{Pt(0, 5)})
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(-10, 5), Vec(1, 0))), []Point[int]{Pt(0, 5)})
	})
	t.Run("pointing away gives none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionRay(Ry(Pt(-10, 0), Vec(-1, 0))))
	})
	t.Run("allocates once for the result and not at all for none", func(t *testing.T) {
		through, away := Ry(Pt(-10, 0), Vec(1, 0)), Ry(Pt(-10, 0), Vec(-1, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.IntersectionRay(through)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.IntersectionRay(away)
		}), 0)
	})
	t.Run("every point lies on the ray and the boundary, and exists exactly where IntersectsRay holds for a ray with a direction", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, r := range rayFixtures {
				points := c.IntersectionRay(r)
				message := fmt.Sprintf("%s → %s: ", c, r)

				for _, point := range points {
					assert.True(t, r.Contains(point), message+point.String()+" on the ray: ")
					assert.True(t, touches(c, point), message+point.String()+" on the boundary: ")
				}
				assert.True(t, len(points) <= 2, message)
				if r.Direction != (Vector[float64]{}) {
					assert.Equal(t, len(points) > 0, c.IntersectsRay(r), message)
				}
			}
		}
	})
}

func TestCircle_AppendIntersectionRay(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)
	through := Ry(Pt(-10, 0), Vec(1, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionRay([]Point[int]{Pt(5, 0)}, through), []Point[int]{Pt(5, 0), Pt(-5, 0), Pt(5, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.AppendIntersectionRay([]Point[int]{Pt(9, 9)}, Ry(Pt(-10, 0), Vec(-1, 0))), []Point[int]{Pt(9, 9)})
		assert.Nil(t, circle.AppendIntersectionRay(nil, Ry(Pt(-10, 0), Vec(-1, 0))))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.AppendIntersectionRay(buffer[:0], through)
		}), 0)
	})
	t.Run("matches IntersectionRay after the points in dst", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, c.AppendIntersectionRay(bufferWith(prefixPoint), r), append([]Point[float64]{prefixPoint}, c.IntersectionRay(r)...), fmt.Sprintf("%s → %s: ", c, r))
			}
		}
	})
}

func TestCircle_IntersectsPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("center inside", func(t *testing.T) {
		assert.True(t, Circ(Pt(1, 1), 5).IntersectsPolygon(square))
	})
	t.Run("an edge within the radius", func(t *testing.T) {
		assert.True(t, Circ(Pt(3, 1), 1).IntersectsPolygon(square))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(4, 1), 1).IntersectsPolygon(square))
		assert.False(t, Circ(Pt(3, 3), 1).IntersectsPolygon(square))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 1).IntersectsPolygon(Pol[int](nil)))
	})
	t.Run("matches the circle test on the polygon edges", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, c := range circleFixtures {
				expected := p.Contains(c.Center) || slices.ContainsFunc(slices.Collect(p.Edges()), func(edge Segment[float64]) bool {
					return c.IntersectsSegment(edge)
				})

				assert.Equal(t, c.IntersectsPolygon(p), expected, fmt.Sprintf("%s → %s: ", p, c))
			}
		}
	})
}

func TestCircle_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, Circ(Pt(150.0, 0.0), 60.0).IntersectsRectangle(rectangle))
		assert.True(t, Circ(Pt(110.0, 80.0), 60.0).IntersectsRectangle(rectangle))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(150.0, 0.0), 40.0).IntersectsRectangle(rectangle))
	})
	t.Run("the circle center on the boundary counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(100.0, 0.0), 1.0).IntersectsRectangle(rectangle))
	})
	t.Run("touching the edge from outside counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(200.0, 0.0), 100.0).IntersectsRectangle(rectangle))
		assert.False(t, Circ(Pt(201.0, 0.0), 100.0).IntersectsRectangle(rectangle))
	})
	t.Run("touching the corner from outside counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(103.0, 54.0), 5.0).IntersectsRectangle(rectangle))
		assert.False(t, Circ(Pt(104.0, 54.0), 5.0).IntersectsRectangle(rectangle))
	})
	t.Run("odd integer sizes keep exact half extents", func(t *testing.T) {
		odd := Rect(Pt(0, 0), Sz(3, 3))

		assert.True(t, Circ(Pt(3, 0), 1).IntersectsRectangle(odd))
		assert.False(t, Circ(Pt(4, 0), 1).IntersectsRectangle(odd))
	})
	t.Run("circle fully inside the rectangle", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.0), 10.0).IntersectsRectangle(rectangle))
	})
	t.Run("rotated measures to the turned edge", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		assert.False(t, Circ(Pt(1.5, 1.5), 0.5).IntersectsRectangle(diamond))
		assert.True(t, Circ(Pt(1.5, 1.5), 0.5).IntersectsRectangle(diamond.Bounds().Rectangle()))
		assert.True(t, Circ(Pt(1.0, 1.0), 0.5).IntersectsRectangle(diamond))
	})
	t.Run("a circle intersects its own bounds", func(t *testing.T) {
		for _, c := range circleFixtures {
			if c.Radius == 0 {
				continue // a degenerate circle touches nothing
			}

			assert.True(t, c.IntersectsRectangle(c.Bounds().Rectangle()), fmt.Sprintf("%s: ", c))
		}
	})
}

func TestCircle_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("center inside", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 1).IntersectsRegularPolygon(diamond))
	})
	t.Run("touching a vertex", func(t *testing.T) {
		assert.True(t, Circ(Pt(3, 0), 1).IntersectsRegularPolygon(diamond))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(4, 0), 1).IntersectsRegularPolygon(diamond))
		assert.False(t, Circ(Pt(2, 2), 1).IntersectsRegularPolygon(diamond))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 1).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, c.IntersectsRegularPolygon(rp), c.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", rp, c))
			}
		}
	})
}

func TestCircle_IntersectsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-100.0, -50.0), Pt(100.0, 50.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, Circ(Pt(150.0, 0.0), 60.0).IntersectsBox(box))
		assert.True(t, Circ(Pt(110.0, 80.0), 60.0).IntersectsBox(box))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(150.0, 0.0), 40.0).IntersectsBox(box))
	})
	t.Run("touching the edge from outside counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(200.0, 0.0), 100.0).IntersectsBox(box))
		assert.False(t, Circ(Pt(201.0, 0.0), 100.0).IntersectsBox(box))
	})
	t.Run("touching the corner from outside counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(103.0, 54.0), 5.0).IntersectsBox(box))
		assert.False(t, Circ(Pt(104.0, 54.0), 5.0).IntersectsBox(box))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Circ(Pt(4, 0), 1).IntersectsBox(BoxFromMinMax(Pt(0, 0), Pt(3, 3))))
		assert.False(t, Circ(Pt(5, 0), 1).IntersectsBox(BoxFromMinMax(Pt(0, 0), Pt(3, 3))))
	})
	t.Run("circle fully inside the box", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.0), 10.0).IntersectsBox(box))
	})
	t.Run("agrees with the rectangle of the box", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, c.IntersectsBox(b), c.IntersectsRectangle(b.Rectangle()), fmt.Sprintf("%s → %s: ", c, b))
			}
		}
	})
}

func TestCircle_Equal(t *testing.T) {
	t.Run("same circle", func(t *testing.T) {
		assert.True(t, Circ(Pt(1, 2), 10).Equal(Circ(Pt(1, 2), 10)))
		assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(0.6, -0.25), 1.2)))
	})
	t.Run("different circle", func(t *testing.T) {
		assert.False(t, Circ(Pt(1, 2), 10).Equal(Circ(Pt(3, -3), 10)))
		assert.False(t, Circ(Pt(1, 2), 10).Equal(Circ(Pt(1, 2), 11)))
		assert.False(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(100.1, -0.1), 1.2)))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(0.6, -0.250001), 1.2)))
	})
}

func TestCircle_IsZero(t *testing.T) {
	t.Run("zero circle", func(t *testing.T) {
		assert.True(t, Circle[int]{}.IsZero())
		assert.True(t, Circ(Pt(0, 0), 0).IsZero())
		assert.True(t, Circle[float64]{}.IsZero())
	})
	t.Run("only the center or only the radius is zero", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 10).IsZero())
		assert.False(t, Circ(Pt(2, 1), 0).IsZero())
		assert.False(t, Circ(Pt(0.0, 0.0), 10.0).IsZero())
		assert.False(t, Circ(Pt(2.0, 1.0), 0.0).IsZero())
	})
	t.Run("non-zero circle", func(t *testing.T) {
		assert.False(t, Circ(Pt(1, 2), 10).IsZero())
		assert.False(t, Circ(Pt(1.0, 2.0), 10.0).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.000001), 0.0).IsZero())
	})
}

func TestCircle_Ellipse(t *testing.T) {
	t.Run("equal semi-axes and no angle", func(t *testing.T) {
		geomtest.AssertEllipse(t, Circ(Pt(1, 2), 10).Ellipse(), Ell(Pt(1, 2), SzU(10), 0))
	})
	t.Run("round-trips through Circle", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Ellipse().Circle(), Circ(Pt(1, 2), 10))
	})
	t.Run("it is what a matrix transforms", func(t *testing.T) {
		c := Circ(Pt(2.0, 3.0), 4.0)

		geomtest.AssertEllipse(t, c.Ellipse().Transform(ScaleMatrix(3.0, 1.0)), Ell(Pt(6.0, 3.0), Sz(12.0, 4.0), 0))
	})
	t.Run("a uniform transform converts back", func(t *testing.T) {
		c, m := Circ(Pt(2.0, 3.0), 4.0), RotationMatrix[float64](Pi/3)
		turned := c.Ellipse().Transform(m)

		assert.True(t, turned.IsCircle())
		geomtest.AssertCircle(t, turned.Circle(), c.MoveTo(c.Center.Transform(m)))
	})
}

func TestCircle_RegularPolygon(t *testing.T) {
	c := Circ(Pt(0.0, 0.0), 10.0)

	t.Run("pointy top places a vertex at the top", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, c.RegularPolygon(6, OrientationPointyTop), RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationPointyTop))
		assertOrientation(t, c.RegularPolygon(6, OrientationPointyTop), OrientationPointyTop, "")
	})
	t.Run("flat top places an edge at the top", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, c.RegularPolygon(6, OrientationFlatTop), RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop))
	})
	t.Run("every vertex lies on the boundary", func(t *testing.T) {
		for vertex := range c.RegularPolygon(7, OrientationPointyTop).Vertices() {
			geomtest.AssertNumber(t, c.Center.DistanceTo(vertex), c.Radius, fmt.Sprintf("%s: ", vertex))
		}
	})
	t.Run("an orientation with no meaning panics", func(t *testing.T) {
		assert.Panics(t, func() {
			c.RegularPolygon(6, OrientationNone)
		})
	})
}

func TestCircle_Cast(t *testing.T) {
	c := Circ(Pt(1.5, -2.5), 3.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertCircle(t, c.Cast[int](), c.Int())
		geomtest.AssertCircle(t, c.Cast[float64](), c.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertCircle(t, c.Cast[int8](), Circ(Pt[int8](2, -3), 4))
	})
}

func TestCircle_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Int(), Circ(Pt(1, 2), 10))
	})
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Int(), Circ(Pt(1, 0), 1))
	})
}

func TestCircle_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).Float(), Circ(Pt(1.0, 2.0), 10.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Float(), Circ(Pt(0.6, -0.25), 1.2))
	})
}

func TestCircle_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(10, 16), 5).String(), "Circ((10,16);5)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(100, -34.0000115), 0.2).String(), "Circ((100.00,-34.00);0.20)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Circ(Pt(100, -34.0000115), 0.2).String)
	})
}

func TestCircle_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Circ(Pt(10, 16), 12), `{"x":10,"y":16,"r":12}`)

		var c Circle[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16,"r":12}`), &c))
		geomtest.AssertCircle(t, c, Circ(Pt(10, 16), 12))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Circ(Pt(100, -34.0000115), 0.2), `{"x":100.0,"y":-34.0000115,"r":0.2}`)

		var c Circle[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115,"r":0.2}`), &c))
		geomtest.AssertCircle(t, c, Circ(Pt(10.1, -34.0000115), 0.2))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, circle := range circleFixtures {
			data, err := json.Marshal(circle)
			assert.NoError(t, err)

			var decoded Circle[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, circle)
		}
	})
}

func TestCircle_Properties(t *testing.T) {
	t.Run("area and circumference follow the radius", func(t *testing.T) {
		for _, c := range circleFixtures {
			geomtest.AssertNumber(t, c.Area(), Pi*c.Radius*c.Radius, fmt.Sprintf("%s: ", c))
			geomtest.AssertNumber(t, c.Perimeter(), 2*Pi*c.Radius, fmt.Sprintf("%s: ", c))
			geomtest.AssertNumber(t, c.Diameter(), 2*c.Radius, fmt.Sprintf("%s: ", c))
		}
	})
	t.Run("lerp ends on the two circles", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, c.Scale(factor).Unscale(factor).Equal(c), fmt.Sprintf("%s ×%v: ", c, factor))
			}
		}
	})
	t.Run("scale keeps the sign of the radius", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, factor := range []float64{0.5, 2.5, -3} {
				assert.Equal(t, Sign(c.Scale(factor).Radius), Sign(c.Radius), fmt.Sprintf("%s ×%v: ", c, factor))
			}
		}
	})
	t.Run("bounds is the circumscribing square", func(t *testing.T) {
		for _, c := range circleFixtures {
			bounds := c.Bounds()

			assert.True(t, bounds.Center().Equal(c.Center), fmt.Sprintf("%s: ", c))
			assert.True(t, bounds.Size().Equal(SzU(c.Diameter())), fmt.Sprintf("%s: ", c))
		}
	})
	t.Run("every anchor is on the boundary", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, direction := range Directions() {
				anchor := c.Anchor(direction)

				geomtest.AssertNumber(t, c.Center.DistanceTo(anchor), c.Radius, fmt.Sprintf("%s → %s: ", c, direction))
				assert.True(t, c.Bounds().Contains(anchor), fmt.Sprintf("%s → %s: ", c, direction))
			}
		}
	})
	t.Run("contains matches the distance to the center", func(t *testing.T) {
		// the boundary is included, like Rectangle.Contains
		for _, c := range circleFixtures {
			for _, point := range pointFixtures {
				inside := c.Center.DistanceTo(point) <= c.Radius

				assert.Equal(t, c.Contains(point), inside, fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("translate and move to keep the radius", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, point := range pointFixtures {
				geomtest.AssertNumber(t, c.MoveTo(point).Radius, c.Radius, fmt.Sprintf("%s → %s: ", c, point))
				assert.True(t, c.MoveTo(point).Center.Equal(point), fmt.Sprintf("%s → %s: ", c, point))
				geomtest.AssertNumber(t, c.Translate(point.Vector()).Radius, c.Radius, fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("grow and shrink are inverse above zero", func(t *testing.T) {
		for _, c := range circleFixtures {
			if c.Radius < 1 {
				continue // shrink clamps to zero, so the growth is not recoverable
			}

			assert.True(t, c.Grow(1).Shrink(1).Equal(c), fmt.Sprintf("%s: ", c))
		}
	})
}

func TestCircle_Immutable(t *testing.T) {
	c := Circ(Pt(1, 2), 10)

	c.Translate(Vec(3, -2))
	c.MoveTo(Pt(4, 3))
	c.Scale(2)
	c.Resize(15)
	c.Grow(1)
	c.Shrink(2)

	geomtest.AssertCircle(t, c, Circ(Pt(1, 2), 10))
}

// circleFixtures span the degenerate, unit, and off-origin cases.

var circleFixtures = []Circle[float64]{
	Circ(Pt(0.0, 0.0), 0.0),
	Circ(Pt(0.0, 0.0), 1.0),
	Circ(Pt(1.0, 2.0), 10.0),
	Circ(Pt(0.6, -0.25), 1.2),
	Circ(Pt(-3.5, 0.25), 7.0),
	Circ(Pt(12.5, -0.1), 0.5),
}

func ExampleCirc() {
	fmt.Println(Circ(Pt(10, 16), 5))
	// Output: Circ((10,16);5)
}

// grazingSegments returns the segments a crossing of the circle is decided on within the
// tolerance: each starts a tenth of Delta inside the boundary, in one of the eight directions
// from the center, and runs along the tangent there, which puts the start in the middle of a
// short chord, or to one of the point fixtures.
func grazingSegments(circle Circle[float64]) []Segment[float64] {
	var segments []Segment[float64]
	for _, direction := range Directions() {
		radius := VectorFromAngle(direction.Angle(), max(circle.Radius-Delta/10, 0))
		start := circle.Center.Add(radius)

		segments = append(segments, Seg(start, start.Add(radius.Normal())))
		for _, p := range pointFixtures {
			segments = append(segments, Seg(start, p))
		}
	}

	return segments
}

// touches reports whether the point lies on the boundary of the circle within the tolerance, by
// the comparison the package makes for an endpoint on the boundary: a zero-length segment at the
// point crosses the boundary exactly there.
func touches[T Number](circle Circle[T], point Point[T]) bool {
	return len(circle.IntersectionSegment(Seg(point, point))) == 1
}

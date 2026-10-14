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

func TestCircleBetween(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, CircleBetween(Pt(0, 0), Pt(6, 8)), Circ(Pt(3, 4), 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, CircleBetween(Pt(0.0, 0.0), Pt(3.0, 4.0)), Circ(Pt(1.5, 2.0), 2.5))
	})
	t.Run("int rounds the center and the radius", func(t *testing.T) {
		geomtest.AssertCircle(t, CircleBetween(Pt(0, 0), Pt(1, 0)), Circ(Pt(1, 0), 1))
		geomtest.AssertCircle(t, CircleBetween(Pt(0, 0), Pt(3, 3)), Circ(Pt(2, 2), 2))
	})
	t.Run("one point twice is a circle of no radius", func(t *testing.T) {
		geomtest.AssertCircle(t, CircleBetween(Pt(1.5, 2.0), Pt(1.5, 2.0)), Circ(Pt(1.5, 2.0), 0.0))
	})
	t.Run("both points lie on the boundary, in either order", func(t *testing.T) {
		for _, s := range segmentFixtures {
			circle := CircleBetween(s.Start, s.End)

			assert.True(t, touches(circle, s.Start), s.String())
			assert.True(t, touches(circle, s.End), s.String())
			geomtest.AssertCircle(t, CircleBetween(s.End, s.Start), circle, s.String())
		}
	})
}

func TestCircleThrough(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		circle, ok := CircleThrough(Pt(0, 0), Pt(6, 0), Pt(0, 8))

		assert.True(t, ok)
		geomtest.AssertCircle(t, circle, Circ(Pt(3, 4), 5))
	})
	t.Run("float", func(t *testing.T) {
		circle, ok := CircleThrough(Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(2.0, 4.0))

		assert.True(t, ok)
		geomtest.AssertCircle(t, circle, Circ(Pt(2.0, 1.5), 2.5))
	})
	t.Run("int far from the origin", func(t *testing.T) {
		circle, ok := CircleThrough(Pt(1<<26, -1<<25), Pt(1<<26+6, -1<<25), Pt(1<<26, -1<<25+8))

		assert.True(t, ok)
		geomtest.AssertCircle(t, circle, Circ(Pt(1<<26+3, -1<<25+4), 5))
	})
	t.Run("rounds each product, the same bits on every architecture", func(t *testing.T) {
		circle, _ := CircleThrough(Pt(tenth, 0.7), Pt(3.3, 0.3), Pt(1.9, 2.6))

		assert.Equal(t, circle.Center.X, 1.7544117647058823)
	})
	t.Run("points on one line have no circle", func(t *testing.T) {
		_, ok := CircleThrough(Pt(0, 0), Pt(1, 1), Pt(3, 3))
		assert.False(t, ok)

		_, ok = CircleThrough(alongLine()[0], alongLine()[1], alongLine()[2])
		assert.False(t, ok)
	})
	t.Run("a point given twice has no circle", func(t *testing.T) {
		_, ok := CircleThrough(Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 0.0))
		assert.False(t, ok)

		_, ok = CircleThrough(Pt(2, 3), Pt(2, 3), Pt(2, 3))
		assert.False(t, ok)
	})
	t.Run("allocates nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = CircleThrough(Pt(0, 0), Pt(6, 0), Pt(0, 8))
		}), 0)
	})
	t.Run("three anchors of a circle give it back, in any order", func(t *testing.T) {
		for _, c := range circleFixtures[1:] {
			a, b, d := c.Anchor(Top), c.Anchor(BottomRight), c.Anchor(Left)

			for _, points := range [][3]Point[float64]{{a, b, d}, {b, d, a}, {d, a, b}, {a, d, b}} {
				circle, ok := CircleThrough(points[0], points[1], points[2])

				assert.True(t, ok, c.String())
				geomtest.AssertCircle(t, circle, c, c.String())
			}
		}
	})
	t.Run("far from the origin the float32 points lie on the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, c := range circleFixtures[1:] {
				c := c.Cast[float32]().Translate(offset)
				a, b, d := c.Anchor(Top), c.Anchor(BottomRight), c.Anchor(Left)
				circle, ok := CircleThrough(a, b, d)

				assert.True(t, ok, c.String())
				for _, p := range []Point[float32]{a, b, d} {
					assert.True(t, touches(circle, p), fmt.Sprintf("%s: %s: ", c, p))
				}
			}
		}
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

func TestCircle_BoundingCircle(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(1, 2), 10).BoundingCircle(), Circ(Pt(1, 2), 10))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).BoundingCircle(), Circ(Pt(0.6, -0.25), 1.2))
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

func TestCircle_DistanceToCircle(t *testing.T) {
	t.Run("the distance between the centers less the radii", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0, 0), 1).DistanceToCircle(Circ(Pt(3, 4), 2)), 2.0)
		geomtest.AssertNumber(t, Circ(Pt(0.0, 0.0), 1.0).DistanceToCircle(Circ(Pt(5.0, 0.0), 1.5)), 2.5)
	})
	t.Run("zero for circles that overlap, touch or lie one within the other", func(t *testing.T) {
		geomtest.AssertNumber(t, Circ(Pt(0, 0), 2).DistanceToCircle(Circ(Pt(3, 0), 2)), 0.0)
		geomtest.AssertNumber(t, Circ(Pt(0, 0), 2).DistanceToCircle(Circ(Pt(4, 0), 2)), 0.0)
		geomtest.AssertNumber(t, Circ(Pt(0, 0), 5).DistanceToCircle(Circ(Pt(1, 0), 1)), 0.0)
	})
	t.Run("zero exactly where the circles intersect, and the same from either side", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, a.DistanceToCircle(b) == 0, a.IntersectsCircle(b), message)
				geomtest.AssertNumber(t, a.DistanceToCircle(b), b.DistanceToCircle(a), message)
				assert.True(t, a.DistanceToCircle(b) >= 0, message)
			}
		}
	})
	t.Run("a NaN coordinate is at NaN distance", func(t *testing.T) {
		assert.True(t, math.IsNaN(Circ(Pt(math.NaN(), 0.0), 1.0).DistanceToCircle(Circ(Pt(5.0, 0.0), 1.5))))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		circle, other := Circ(Pt(0, 0), 1), Circ(Pt(3, 4), 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.DistanceToCircle(other) > 0
		}), 0)
	})
}

func TestCircle_DistanceSquaredToCircle(t *testing.T) {
	geomtest.AssertNumber(t, Circ(Pt(0, 0), 1).DistanceSquaredToCircle(Circ(Pt(3, 4), 2)), 4.0)
	geomtest.AssertNumber(t, Circ(Pt(0, 0), 2).DistanceSquaredToCircle(Circ(Pt(3, 0), 2)), 0.0)
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

func TestCircle_Support(t *testing.T) {
	t.Run("the center moved by the radius along the direction", func(t *testing.T) {
		geomtest.AssertPoint(t, Circ(Pt(1.0, 2.0), 5.0).Support(Vec(3.0, 4.0)), Pt(4.0, 6.0))
		geomtest.AssertPoint(t, Circ(Pt(1.0, 2.0), 5.0).Support(Vec(-0.3, -0.4)), Pt(-2.0, -2.0))
	})
	t.Run("int is not rounded", func(t *testing.T) {
		geomtest.AssertPoint(t, Circ(Pt(0, 0), 1).Support(Vec(1.0, 1.0)), Pt(math.Sqrt2/2, math.Sqrt2/2))
	})
	t.Run("a direction of no length points along X", func(t *testing.T) {
		geomtest.AssertPoint(t, Circ(Pt(1.0, 2.0), 5.0).Support(Vec(0.0, 0.0)), Pt(6.0, 2.0))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		circle, direction := Circ(Pt(1, 2), 5), Vec(tenth, 0.7)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.Support(direction).X > 0
		}), 0)
	})
}

func TestCircle_Tangents(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 3.0)

	t.Run("a point outside gives two, mirrored across the line to the center", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.Tangents(Pt(5.0, 0.0)).AppendTo(nil), []Point[float64]{Pt(1.8, 2.4), Pt(1.8, -2.4)})
		geomtest.AssertVertices(t, circle.Tangents(Pt(0.0, -5.0)).AppendTo(nil), []Point[float64]{Pt(2.4, -1.8), Pt(-2.4, -1.8)})
	})
	t.Run("a point on the boundary is the one tangent point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.Tangents(Pt(3.0, 0.0)).AppendTo(nil), []Point[float64]{Pt(3.0, 0.0)})
	})
	t.Run("a point within the tolerance of the boundary is kept as it is", func(t *testing.T) {
		assert.Equal(t, circle.Tangents(Pt(3.0+Delta/2, 0.0)).AppendTo(nil), []Point[float64]{Pt(3.0+Delta/2, 0.0)})
		assert.Equal(t, circle.Tangents(Pt(3.0-Delta/2, 0.0)).AppendTo(nil), []Point[float64]{Pt(3.0-Delta/2, 0.0)})
	})
	t.Run("a point inside gives none", func(t *testing.T) {
		assert.Nil(t, circle.Tangents(Pt(1.0, 1.0)).AppendTo(nil))
		assert.Nil(t, circle.Tangents(Pt(0.0, 0.0)).AppendTo(nil))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).Tangents(Pt(13, 0)).AppendTo(nil), []Point[int]{Pt(2, 5), Pt(2, -5)})
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 1).Tangents(Pt(1, 1)).AppendTo(nil), []Point[int]{Pt(0, 1), Pt(1, 0)})
	})
	t.Run("a zero radius is touched at its center", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(1.0, 1.0), 0.0).Tangents(Pt(5.0, 4.0)).AppendTo(nil), []Point[float64]{Pt(1.0, 1.0)})
		geomtest.AssertVertices(t, Circ(Pt(1, 1), 0).Tangents(Pt(1, 1)).AppendTo(nil), []Point[int]{Pt(1, 1)})
	})
	t.Run("points lie on the boundary, at a right angle to the radius, and mirror Contains", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				points := c.Tangents(p).AppendTo(nil)
				message := fmt.Sprintf("%s → %s: ", c, p)

				switch {
				case touches(c, p):
					assert.Equal(t, points, []Point[float64]{p}, message)
				case c.Contains(p):
					assert.Nil(t, points, message)
				case c.Radius == 0:
					geomtest.AssertVertices(t, points, []Point[float64]{c.Center}, message)
				default:
					assert.Length(t, points, 2, message)
					for _, tangent := range points {
						assert.True(t, touches(c, tangent), message)
						geomtest.AssertNumber(t, tangent.Subtract(c.Center).Normalize().Dot(p.Subtract(tangent).Normalize()), 0, message)
					}
				}
			}
		}
	})
	t.Run("far from the origin float32 points lie on the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, c := range circleFixtures {
				for _, p := range pointFixtures {
					c, p := c.Cast[float32]().Translate(offset), p.Cast[float32]().Add(offset)

					for _, tangent := range c.Tangents(p).AppendTo(nil) {
						assert.True(t, touches(c, tangent), fmt.Sprintf("%s → %s: %s: ", c, p, tangent))
					}
				}
			}
		}
	})
}

func BenchmarkCircle_Tangents(b *testing.B) {
	circle, point := Circ(Pt(0.0, 0.0), 3.0), Pt(5.0, 0.0)

	for b.Loop() {
		_ = circle.Tangents(point)
	}
}

func TestCircle_OuterTangents(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 3.0)

	t.Run("equal circles are joined along both sides", func(t *testing.T) {
		geomtest.AssertSegments(t, circle.OuterTangents(Circ(Pt(10.0, 0.0), 3.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(0.0, 3.0), Pt(10.0, 3.0)), Seg(Pt(0.0, -3.0), Pt(10.0, -3.0))})
	})
	t.Run("unequal circles are joined on a slant", func(t *testing.T) {
		geomtest.AssertSegments(t, Circ(Pt(0.0, 0.0), 8.0).OuterTangents(Circ(Pt(10.0, 0.0), 2.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(4.8, 6.4), Pt(11.2, 1.6)), Seg(Pt(4.8, -6.4), Pt(11.2, -1.6))})
	})
	t.Run("touching from outside and overlapping circles keep both", func(t *testing.T) {
		geomtest.AssertSegments(t, circle.OuterTangents(Circ(Pt(6.0, 0.0), 3.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(0.0, 3.0), Pt(6.0, 3.0)), Seg(Pt(0.0, -3.0), Pt(6.0, -3.0))})
		geomtest.AssertSegments(t, circle.OuterTangents(Circ(Pt(4.0, 0.0), 3.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(0.0, 3.0), Pt(4.0, 3.0)), Seg(Pt(0.0, -3.0), Pt(4.0, -3.0))})
	})
	t.Run("a circle within the other, touching or not, has none", func(t *testing.T) {
		assert.Nil(t, circle.OuterTangents(Circ(Pt(1.0, 0.0), 2.0)).AppendTo(nil))
		assert.Nil(t, circle.OuterTangents(Circ(Pt(1.0+Delta/2, 0.0), 2.0)).AppendTo(nil))
		assert.Nil(t, circle.OuterTangents(Circ(Pt(0.5, 0.0), 1.0)).AppendTo(nil))
		assert.Nil(t, Circ(Pt(0.5, 0.0), 1.0).OuterTangents(circle).AppendTo(nil))
		assert.Nil(t, circle.OuterTangents(circle).AppendTo(nil))
	})
	t.Run("a zero radius is touched at its center, from the Tangents of it", func(t *testing.T) {
		geomtest.AssertSegments(t, circle.OuterTangents(Circ(Pt(5.0, 0.0), 0.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.8, 2.4), Pt(5.0, 0.0)), Seg(Pt(1.8, -2.4), Pt(5.0, 0.0))})
	})
	t.Run("two zero radii are joined once", func(t *testing.T) {
		geomtest.AssertSegments(t, Circ(Pt(1.0, 1.0), 0.0).OuterTangents(Circ(Pt(5.0, 4.0), 0.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(5.0, 4.0))})
		assert.Nil(t, Circ(Pt(1.0, 1.0), 0.0).OuterTangents(Circ(Pt(1.0, 1.0), 0.0)).AppendTo(nil))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertSegments(t, Circ(Pt(0, 0), 3).OuterTangents(Circ(Pt(10, 0), 3)).AppendTo(nil), []Segment[int]{Seg(Pt(0, 3), Pt(10, 3)), Seg(Pt(0, -3), Pt(10, -3))})
		geomtest.AssertSegments(t, Circ(Pt(0, 0), 8).OuterTangents(Circ(Pt(10, 0), 2)).AppendTo(nil), []Segment[int]{Seg(Pt(5, 6), Pt(11, 2)), Seg(Pt(5, -6), Pt(11, -2))})
	})
	t.Run("ends lie on their circles, at a right angle to the radius, from either side", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assertTangents(t, a, b, a.OuterTangents(b).AppendTo(nil), b.OuterTangents(a).AppendTo(nil))
			}
		}
	})
	t.Run("far from the origin float32 ends lie on their circles", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, a := range circleFixtures {
				for _, b := range circleFixtures {
					a, b := a.Cast[float32]().Translate(offset), b.Cast[float32]().Translate(offset)

					for _, s := range a.OuterTangents(b).AppendTo(nil) {
						assert.True(t, touches(a, s.Start), fmt.Sprintf("%s → %s: %s on a: ", a, b, s))
						assert.True(t, touches(b, s.End), fmt.Sprintf("%s → %s: %s on b: ", a, b, s))
					}
				}
			}
		}
	})
}

func BenchmarkCircle_OuterTangents(b *testing.B) {
	circle, other := Circ(Pt(0.0, 0.0), 8.0), Circ(Pt(10.0, 0.0), 2.0)

	for b.Loop() {
		_ = circle.OuterTangents(other)
	}
}

func TestCircle_InnerTangents(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 3.0)

	t.Run("circles apart are joined across the line of centers", func(t *testing.T) {
		geomtest.AssertSegments(t, circle.InnerTangents(Circ(Pt(10.0, 0.0), 3.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.8, 2.4), Pt(8.2, -2.4)), Seg(Pt(1.8, -2.4), Pt(8.2, 2.4))})
		geomtest.AssertSegments(t, Circ(Pt(0.0, 0.0), 2.0).InnerTangents(Circ(Pt(10.0, 0.0), 4.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.2, 1.6), Pt(7.6, -3.2)), Seg(Pt(1.2, -1.6), Pt(7.6, 3.2))})
	})
	t.Run("touching and overlapping circles have none", func(t *testing.T) {
		assert.Nil(t, circle.InnerTangents(Circ(Pt(6.0, 0.0), 3.0)).AppendTo(nil))
		assert.Nil(t, circle.InnerTangents(Circ(Pt(6.0+Delta/2, 0.0), 3.0)).AppendTo(nil))
		assert.Nil(t, circle.InnerTangents(Circ(Pt(4.0, 0.0), 3.0)).AppendTo(nil))
		assert.Nil(t, circle.InnerTangents(Circ(Pt(1.0, 0.0), 1.0)).AppendTo(nil))
		assert.Nil(t, circle.InnerTangents(circle).AppendTo(nil))
	})
	t.Run("a zero radius is touched at its center, from the Tangents of it", func(t *testing.T) {
		geomtest.AssertSegments(t, circle.InnerTangents(Circ(Pt(5.0, 0.0), 0.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.8, 2.4), Pt(5.0, 0.0)), Seg(Pt(1.8, -2.4), Pt(5.0, 0.0))})
	})
	t.Run("two zero radii are joined once", func(t *testing.T) {
		geomtest.AssertSegments(t, Circ(Pt(1.0, 1.0), 0.0).InnerTangents(Circ(Pt(5.0, 4.0), 0.0)).AppendTo(nil), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(5.0, 4.0))})
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertSegments(t, Circ(Pt(0, 0), 3).InnerTangents(Circ(Pt(10, 0), 3)).AppendTo(nil), []Segment[int]{Seg(Pt(2, 2), Pt(8, -2)), Seg(Pt(2, -2), Pt(8, 2))})
	})
	t.Run("exist exactly where the circles do not intersect", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.Equal(t, len(a.InnerTangents(b).AppendTo(nil)) > 0, !a.IntersectsCircle(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("ends lie on their circles, at a right angle to the radius, from either side", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assertTangents(t, a, b, a.InnerTangents(b).AppendTo(nil), b.InnerTangents(a).AppendTo(nil))
			}
		}
	})
	t.Run("far from the origin float32 ends lie on their circles", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, a := range circleFixtures {
				for _, b := range circleFixtures {
					a, b := a.Cast[float32]().Translate(offset), b.Cast[float32]().Translate(offset)

					for _, s := range a.InnerTangents(b).AppendTo(nil) {
						assert.True(t, touches(a, s.Start), fmt.Sprintf("%s → %s: %s on a: ", a, b, s))
						assert.True(t, touches(b, s.End), fmt.Sprintf("%s → %s: %s on b: ", a, b, s))
					}
				}
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

func TestCircle_ContainsCircle(t *testing.T) {
	circle := Circ(Pt(0, 0), 10)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.ContainsCircle(Circ(Pt(2, 0), 5)))
		assert.True(t, circle.ContainsCircle(Circ(Pt(0, 0), 0)))
	})
	t.Run("touching the boundary from inside counts", func(t *testing.T) {
		assert.True(t, circle.ContainsCircle(Circ(Pt(5, 0), 5)))
		assert.False(t, circle.ContainsCircle(Circ(Pt(6, 0), 5)))
	})
	t.Run("larger or apart", func(t *testing.T) {
		assert.False(t, circle.ContainsCircle(Circ(Pt(0, 0), 11)))
		assert.False(t, circle.ContainsCircle(Circ(Pt(30, 0), 1)))
	})
	t.Run("float within the tolerance of the boundary counts, beyond it not", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.0), 10.0).ContainsCircle(Circ(Pt(5.0+Delta/2, 0.0), 5.0)))
		assert.False(t, Circ(Pt(0.0, 0.0), 10.0).ContainsCircle(Circ(Pt(5.0+2*Delta, 0.0), 5.0)))
	})
	t.Run("encloses exactly where it encloses the far point", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				far := b.Center.AddXY(b.Radius, 0)
				if direction := b.Center.Subtract(a.Center); direction != (Vector[float64]{}) {
					far = b.Center.Add(direction.Resize(b.Radius))
				}

				assert.Equal(t, a.ContainsCircle(b), a.Contains(far), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestCircle_ContainsCapsule(t *testing.T) {
	t.Run("touching the boundary from inside counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 10).ContainsCapsule(Cap(Pt(-5, 0), Pt(5, 0), 5)))
		assert.False(t, Circ(Pt(0, 0), 10).ContainsCapsule(Cap(Pt(-5, 0), Pt(5, 0), 6)))
	})
	t.Run("an end outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 10).ContainsCapsule(Cap(Pt(0, 0), Pt(12, 0), 1)))
	})
}

func TestCircle_ContainsSector(t *testing.T) {
	t.Run("a sector of the circle or within it", func(t *testing.T) {
		assert.True(t, Circ(Pt(0.0, 0.0), 5.0).ContainsSector(quarterSector()))
		assert.True(t, Circ(Pt(1.0, 1.0), 5.0).ContainsSector(quarterSector()))
		assert.True(t, Circ(Pt(0.0, -1.0), 6.0).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
	})
	t.Run("a sector with a corner outside, or its arc alone", func(t *testing.T) {
		assert.False(t, Circ(Pt(0.0, 0.0), 4.9).ContainsSector(quarterSector()))
		assert.False(t, Circ(Pt(0.0, -1.0), 5.9).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 5).ContainsSector(Sect(Pt(0, 0), 5, 0, Pi/2)))
		assert.False(t, Circ(Pt(0, 0), 5).ContainsSector(Sect(Pt(1, 0), 5, 0, Pi/2)))
	})
}

func TestCircle_ContainsSegment(t *testing.T) {
	circle := Circ(Pt(0, 0), 10)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.ContainsSegment(Seg(Pt(-5, 0), Pt(5, 5))))
	})
	t.Run("a chord counts", func(t *testing.T) {
		assert.True(t, circle.ContainsSegment(Seg(Pt(-10, 0), Pt(10, 0))))
		assert.True(t, circle.ContainsSegment(Seg(Pt(-6, 8), Pt(8, -6))))
	})
	t.Run("an endpoint outside", func(t *testing.T) {
		assert.False(t, circle.ContainsSegment(Seg(Pt(0, 0), Pt(11, 0))))
	})
}

func TestCircle_ContainsPolygon(t *testing.T) {
	circle := Circ(Pt(1, 1), 2)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, circle.ContainsPolygon(Pol(squareVertices())))
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(1, 1), 1).ContainsPolygon(Pol(squareVertices())))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, circle.ContainsPolygon(Pol[int](nil)))
	})
}

func TestCircle_ContainsRectangle(t *testing.T) {
	t.Run("corners on the boundary count", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 5).ContainsRectangle(Rect(Pt(0, 0), Sz(6, 8))))
		assert.False(t, Circ(Pt(0, 0), 5).ContainsRectangle(Rect(Pt(0, 0), Sz(6, 9))))
	})
	t.Run("rotated is tested on its turned corners", func(t *testing.T) {
		long := Rect(Pt(0.0, 0.0), Sz(9.0, 1.0))

		assert.False(t, Circ(Pt(0.0, 0.0), 4.0).ContainsRectangle(long))
		assert.False(t, Circ(Pt(0.0, 0.0), 4.0).ContainsRectangle(long.Rotate(Pi/4)))
		assert.True(t, Circ(Pt(0.0, 0.0), 5.0).ContainsRectangle(long.Rotate(Pi/4)))
	})
}

func TestCircle_ContainsRegularPolygon(t *testing.T) {
	t.Run("a circle encloses its inscribed polygon", func(t *testing.T) {
		for _, c := range circleFixtures {
			assert.True(t, c.ContainsRegularPolygon(c.RegularPolygon(7, OrientationFlatTop)), c.String())
		}
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 2).ContainsRegularPolygon(RegPol(Pt(1, 0), Sz(2, 2), 4, 0, 0)))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, Circ(Pt(0, 0), 2).ContainsRegularPolygon(RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0)))
	})
}

func TestCircle_ContainsBox(t *testing.T) {
	t.Run("an inscribed box counts, the bounds of the circle not", func(t *testing.T) {
		c := Circ(Pt(0, 0), 5)

		assert.True(t, c.ContainsBox(BoxFromMinMax(Pt(-3, -4), Pt(3, 4))))
		assert.False(t, c.ContainsBox(c.Bounds()))
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
				if len(a.IntersectionCircle(b).AppendTo(nil)) > 0 {
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

func TestCircle_IntersectsCapsule(t *testing.T) {
	capsule := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("touching a side or a round end counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(5, 5), 3).IntersectsCapsule(capsule))
		assert.True(t, Circ(Pt(15, 0), 3).IntersectsCapsule(capsule))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(5, 6), 3).IntersectsCapsule(capsule))
		assert.False(t, Circ(Pt(16, 0), 3).IntersectsCapsule(capsule))
	})
	t.Run("float within the tolerance counts, beyond it not", func(t *testing.T) {
		assert.True(t, Circ(Pt(5.0, 5.0+Delta/2), 3.0).IntersectsCapsule(capsule.Float()))
		assert.False(t, Circ(Pt(5.0, 5.0+2*Delta), 3.0).IntersectsCapsule(capsule.Float()))
	})
	t.Run("int is exact near the end of the supported range", func(t *testing.T) {
		long := Cap(Pt(0, 0), Pt(1<<25, 0), 3)

		assert.True(t, Circ(Pt(1<<24, 5), 2).IntersectsCapsule(long))
		assert.False(t, Circ(Pt(1<<24, 5), 1).IntersectsCapsule(long))
	})
	t.Run("a circle of no radius is its center, a capsule of no length a circle", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, Circ(p, 0.0).IntersectsCapsule(c), c.Contains(p), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.Equal(t, a.IntersectsCapsule(Cap(b.Center, b.Center, b.Radius)), a.IntersectsCircle(b), fmt.Sprintf("%s → %s: ", a, b))
			}
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
				assert.Length(t, c.IntersectionSegment(wall).AppendTo(nil), 1, message)
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

func TestCircle_IntersectsLine(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through, whichever way it runs", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(-2.0, 0.0), Vec(1.0, 0.0))))
		assert.True(t, circle.IntersectsLine(Ln(Pt(-2.0, 0.0), Vec(-1.0, 0.0))))
	})
	t.Run("tangent counts, on either side of the origin", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(-2.0, 1.0), Vec(1.0, 0.0))))
		assert.True(t, circle.IntersectsLine(Ln(Pt(1.5, 1.0), Vec(1.0, 0.0))))
		assert.False(t, circle.IntersectsLine(Ln(Pt(-2.0, 1.0+2*Delta), Vec(1.0, 0.0))))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, circle.IntersectsLine(Ln(Pt(-2.0, 0.0), Vec(0.0, 1.0))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(0.5, 0.0), Vec(0.0, 0.0))))
		assert.False(t, circle.IntersectsLine(Ln(Pt(2.0, 0.0), Vec(0.0, 0.0))))
	})
	t.Run("matches a segment reaching past the circle along the line both ways", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, l := range lineFixtures {
				assert.Equal(t, c.IntersectsLine(l), c.IntersectsSegment(wide(l)), fmt.Sprintf("%s → %s: ", c, l))
			}
		}
	})
	t.Run("holds where either ray from the origin of the line meets the circle", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, l := range lineFixtures {
				ahead, behind := Ry(l.Origin, l.Direction), Ry(l.Origin, l.Direction.Negate())

				assert.True(t, !(c.IntersectsRay(ahead) || c.IntersectsRay(behind)) || c.IntersectsLine(l), fmt.Sprintf("%s → %s: ", c, l))
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

func TestCircle_IntersectsSector(t *testing.T) {
	sector := Sect(Pt(0.0, 0.0), 5.0, 0, Pi/2)

	t.Run("touching the arc or a radius counts", func(t *testing.T) {
		assert.True(t, Circ(Pt(6.0, 8.0), 5.0).IntersectsSector(sector))
		assert.True(t, Circ(Pt(-1.0, 1.0), 1.0).IntersectsSector(sector))
	})
	t.Run("a center within the sector", func(t *testing.T) {
		assert.True(t, Circ(Pt(2.0, 2.0), 0.0).IntersectsSector(sector))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Circ(Pt(6.0, 8.0), 4.9).IntersectsSector(sector))
		assert.False(t, Circ(Pt(2.0, -2.0), 1.0).IntersectsSector(sector))
	})
	t.Run("float within the tolerance counts, beyond it not", func(t *testing.T) {
		assert.True(t, Circ(Pt(2.0, -1.0-Delta/2), 1.0).IntersectsSector(sector))
		assert.False(t, Circ(Pt(2.0, -1.0-2*Delta), 1.0).IntersectsSector(sector))
	})
}

func TestCircle_IntersectionCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 5.0)

	t.Run("overlapping gives two points mirrored across the centers", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(6.0, 0.0), 5.0)).AppendTo(nil), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, -4.0)})
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(0.0, 6.0), 5.0)).AppendTo(nil), []Point[float64]{Pt(-4.0, 3.0), Pt(4.0, 3.0)})
	})
	t.Run("a larger circle behind the first", func(t *testing.T) {
		height := math.Sqrt(1 - 0.125*0.125)

		geomtest.AssertVertices(t, Circ(Pt(0.0, 0.0), 1.0).IntersectionCircle(Circ(Pt(1.0, 0.0), 1.5)).AppendTo(nil), []Point[float64]{Pt(-0.125, height), Pt(-0.125, -height)})
	})
	t.Run("tangent from outside gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(8.0, 0.0), 3.0)).AppendTo(nil), []Point[float64]{Pt(5.0, 0.0)})
	})
	t.Run("tangent from inside gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionCircle(Circ(Pt(2.0, 0.0), 3.0)).AppendTo(nil), []Point[float64]{Pt(5.0, 0.0)})
		geomtest.AssertVertices(t, Circ(Pt(0.0, 0.0), 3.0).IntersectionCircle(Circ(Pt(2.0, 0.0), 5.0)).AppendTo(nil), []Point[float64]{Pt(-3.0, 0.0)})
	})
	t.Run("apart, nested and concentric give none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(20.0, 0.0), 5.0)).AppendTo(nil))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(1.0, 0.0), 1.0)).AppendTo(nil))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(0.0, 0.0), 3.0)).AppendTo(nil))
		assert.Nil(t, circle.IntersectionCircle(circle).AppendTo(nil))
	})
	t.Run("float treats centers within Epsilon as coincident", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(1e-7, 0.0), 5.0)).AppendTo(nil))
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(0.0, -1e-7), 3.0)).AppendTo(nil))
	})
	t.Run("float is tolerant at a tangent", func(t *testing.T) {
		assert.Length(t, circle.IntersectionCircle(Circ(Pt(8.0+Delta/2, 0.0), 3.0)).AppendTo(nil), 1)
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(8.0+2*Delta, 0.0), 3.0)).AppendTo(nil))
	})
	t.Run("a tangent exactly Delta outside is judged like Intersects, on both sides", func(t *testing.T) {
		for _, other := range []Circle[float64]{Circ(Pt(8.0+Delta, 0.0), 3.0), Circ(Pt(2.0+Delta, 0.0), 3.0)} {
			assert.True(t, circle.IntersectsCircle(other))
			geomtest.AssertVertices(t, circle.IntersectionCircle(other).AppendTo(nil), []Point[float64]{Pt(5.0, 0.0)})
		}
	})
	t.Run("a tangent within the tolerance lands on both boundaries", func(t *testing.T) {
		a, b := Circ(Pt(0.0, 0.0), 6.0), Circ(Pt(2.0000005, 0.0), 4.0)
		points := a.IntersectionCircle(b).AppendTo(nil)

		assert.Length(t, points, 1)
		for _, p := range points {
			assert.True(t, touches(a, p), "on a")
			assert.True(t, touches(b, p), "on b")
		}
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)).AppendTo(nil), []Point[int]{Pt(3, 4), Pt(3, -4)})
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 2).IntersectionCircle(Circ(Pt(3, 0), 2)).AppendTo(nil), []Point[int]{Pt(2, 1), Pt(2, -1)})
	})
	t.Run("int gives two crossings rounding to one point once", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(11, -20), 15).IntersectionCircle(Circ(Pt(-3, -19), 1)).AppendTo(nil), []Point[int]{Pt(-4, -19)})
	})
	t.Run("points lie on both circles and mirror Intersects", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				points := a.IntersectionCircle(b).AppendTo(nil)

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

					for _, p := range a.IntersectionCircle(b).AppendTo(nil) {
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
		points := a.IntersectionCircle(b).AppendTo(nil)

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

		if rounded := a.Int().IntersectionCircle(b.Int()).AppendTo(nil); len(rounded) == 2 {
			assert.True(t, !rounded[0].Equal(rounded[1]), fmt.Sprintf("%s → %s: two distinct rounded points: ", a.Int(), b.Int()))
		}
	})
}

func TestCircle_IntersectionSegment(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(2.0, 0.0), Pt(-2.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(1.0, 0.0), Pt(-1.0, 0.0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(0.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(-1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(0.5, 0.0), Pt(5.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		unit := Circ(Pt(0, 0), 5)
		through, apart := Seg(Pt(-10, 0), Pt(10, 0)), Seg(Pt(-10, 9), Pt(10, 9))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = unit.IntersectionSegment(through).Count > 0
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = unit.IntersectionSegment(apart).Count > 0
		}), 0)
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0))).AppendTo(nil), []Point[float64]{Pt(0.0, 1.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0+Delta/2), Pt(2.0, 1.0+Delta/2))).AppendTo(nil), []Point[float64]{Pt(0.0, 1.0+Delta/2)})
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta))).AppendTo(nil))
	})
	t.Run("a chord within the tolerance band of the gap keeps both ends", func(t *testing.T) {
		height := 1.0 - Delta/2
		half := math.Sqrt(1 - height*height)

		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, height), Pt(1.0, height))).AppendTo(nil), []Point[float64]{Pt(-half, height), Pt(half, height)})
	})
	t.Run("a chord shorter than Delta is a tangent", func(t *testing.T) {
		height := math.Sqrt(1 - Delta*Delta/16)

		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, height), Pt(1.0, height))).AppendTo(nil), []Point[float64]{Pt(0.0, height)})
	})
	t.Run("a segment within the tolerance with both ends on the boundary gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(1.0, 0.0), Pt(1.0, Delta/10))).AppendTo(nil), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("a tangent beyond the segment is missed", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(1.0, 1.0), Pt(2.0, 1.0))).AppendTo(nil))
	})
	t.Run("a shallow touch is decided on the endpoint distance, like IntersectsSegment", func(t *testing.T) {
		shallow := Seg(Pt(-1.0, 1.0+Delta/2), Pt(-0.0005, 1.0+Delta/2))

		assert.True(t, circle.IntersectsSegment(shallow))
		geomtest.AssertVertices(t, circle.IntersectionSegment(shallow).AppendTo(nil), []Point[float64]{shallow.End})
		geomtest.AssertVertices(t, circle.IntersectionSegment(shallow.Reverse()).AppendTo(nil), []Point[float64]{shallow.End})
	})
	t.Run("an interior graze exactly Delta outside is judged like IntersectsSegment", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 1.000001), Pt(1.0, 1.000001)), Circ(Pt(0.0, 0.0), 1.0)

		assert.Equal(t, len(c.IntersectionSegment(s).AppendTo(nil)) > 0, c.IntersectsSegment(s))
		for _, y := range []float64{1.0000009, 1.0000011, 1.0000015} {
			s := Seg(Pt(-1.0, y), Pt(1.0, y))

			assert.Equal(t, len(c.IntersectionSegment(s).AppendTo(nil)) > 0, c.IntersectsSegment(s), s.String())
		}
	})
	t.Run("an endpoint exactly Delta outside is judged like IntersectsSegment", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 2.000001), Pt(-882.0315, 56.11111116666667)), Circ(Pt(-1.0, 0.0), 2.0)

		assert.Equal(t, len(c.IntersectionSegment(s).AppendTo(nil)) > 0, c.IntersectsSegment(s))
	})
	t.Run("an endpoint on the boundary is counted once with its crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-1.0, 0.0), Pt(2.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(-2.0, 0.0), Pt(1.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
	})
	t.Run("an endpoint on the boundary before the chord is where the segment enters", func(t *testing.T) {
		s, c := Seg(Pt(-1.0, 2.000001), Pt(-877.0315, 0.11111116666666668)), Circ(Pt(-1.0, 0.0), 2.0)
		points := c.IntersectionSegment(s).AppendTo(nil)

		assert.Length(t, points, 2)
		geomtest.AssertPoint(t, points[0], s.Start)
	})
	t.Run("an endpoint in the middle of a chord keeps the crossing ahead of it, whichever way the segment runs", func(t *testing.T) {
		height := 1.0 - Delta/10
		half := math.Sqrt(1 - height*height)
		s := Seg(Pt(0.0, height), Pt(5.0, height))

		geomtest.AssertVertices(t, circle.IntersectionSegment(s).AppendTo(nil), []Point[float64]{s.Start, Pt(half, height)})
		geomtest.AssertVertices(t, circle.IntersectionSegment(s.Reverse()).AppendTo(nil), []Point[float64]{Pt(half, height), s.Start})
	})
	t.Run("an endpoint on the boundary of a segment leaving the circle is the one crossing", func(t *testing.T) {
		leaving := Seg(Pt(1.0, 0.0), Pt(2.0, 0.0))

		geomtest.AssertVertices(t, circle.IntersectionSegment(leaving).AppendTo(nil), []Point[float64]{leaving.Start})
		geomtest.AssertVertices(t, circle.IntersectionSegment(leaving.Reverse()).AppendTo(nil), []Point[float64]{leaving.Start})
	})
	t.Run("a tangent segment with both ends on the boundary gives its ends", func(t *testing.T) {
		grazing := Seg(Pt(-0.0003, 1.0), Pt(0.0003, 1.0))

		geomtest.AssertVertices(t, circle.IntersectionSegment(grazing).AppendTo(nil), []Point[float64]{grazing.Start, grazing.End})
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0))).AppendTo(nil))
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(2.0, 0.0), Pt(3.0, 0.0))).AppendTo(nil))
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(-0.5, 0.0), Pt(0.5, 0.0))).AppendTo(nil))
	})
	t.Run("a degenerate segment is a point on the boundary or nothing", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionSegment(Seg(Pt(1.0, 0.0), Pt(1.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(1.0, 0.0)})
		assert.Nil(t, circle.IntersectionSegment(Seg(Pt(0.5, 0.0), Pt(0.5, 0.0))).AppendTo(nil))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionSegment(Seg(Pt(-5, -5), Pt(5, 5))).AppendTo(nil), []Point[int]{Pt(-4, -4), Pt(4, 4)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsSegment holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				points := c.IntersectionSegment(s).AppendTo(nil)

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
				assertReversed(t, c.IntersectionSegment(s).AppendTo(nil), c.IntersectionSegment(s.Reverse()).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("far from the origin a float32 segment and its Reverse cross as often", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, c := range circleFixtures {
				for _, s := range append(grazingSegments(c), segmentFixtures...) {
					s, c := s.Cast[float32]().Translate(offset), c.Cast[float32]().Translate(offset)

					assert.Length(t, c.IntersectionSegment(s.Reverse()).AppendTo(nil), len(c.IntersectionSegment(s).AppendTo(nil)), fmt.Sprintf("%s → %s: ", s, c))
				}
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, c := range circleFixtures {
					s, c := s.Cast[float32]().Translate(offset), c.Cast[float32]().Translate(offset)
					points := c.IntersectionSegment(s).AppendTo(nil)

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
		points := c.IntersectionSegment(s).AppendTo(nil)

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

func TestCircle_IntersectionRay(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(-10, 0), Vec(1, 0))).AppendTo(nil), []Point[int]{Pt(-5, 0), Pt(5, 0)})
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(10, 0), Vec(-2, 0))).AppendTo(nil), []Point[int]{Pt(5, 0), Pt(-5, 0)})
	})
	t.Run("starting inside gives the exit", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(0, 0), Vec(0, 1))).AppendTo(nil), []Point[int]{Pt(0, 5)})
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionRay(Ry(Pt(-10, 5), Vec(1, 0))).AppendTo(nil), []Point[int]{Pt(0, 5)})
	})
	t.Run("pointing away gives none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionRay(Ry(Pt(-10, 0), Vec(-1, 0))).AppendTo(nil))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through, away := Ry(Pt(-10, 0), Vec(1, 0)), Ry(Pt(-10, 0), Vec(-1, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.IntersectionRay(through).Count > 0
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.IntersectionRay(away).Count > 0
		}), 0)
	})
	t.Run("every point lies on the ray and the boundary, and exists exactly where IntersectsRay holds for a ray with a direction", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, r := range rayFixtures {
				points := c.IntersectionRay(r).AppendTo(nil)
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

func TestCircle_IntersectionLine(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionLine(Ln(Pt(-10, 0), Vec(1, 0))).AppendTo(nil), []Point[int]{Pt(-5, 0), Pt(5, 0)})
		geomtest.AssertVertices(t, circle.IntersectionLine(Ln(Pt(-10, 0), Vec(-2, 0))).AppendTo(nil), []Point[int]{Pt(5, 0), Pt(-5, 0)})
	})
	t.Run("an origin inside gives both", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionLine(Ln(Pt(0, 0), Vec(0, 1))).AppendTo(nil), []Point[int]{Pt(0, -5), Pt(0, 5)})
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		geomtest.AssertVertices(t, circle.IntersectionLine(Ln(Pt(10, 5), Vec(1, 0))).AppendTo(nil), []Point[int]{Pt(0, 5)})
	})
	t.Run("passing by gives none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(-10, 0), Vec(0, 1))).AppendTo(nil))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through, apart := Ln(Pt(-10, 0), Vec(1, 0)), Ln(Pt(-10, 0), Vec(0, 1))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.IntersectionLine(through).Count > 0
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = circle.IntersectionLine(apart).Count > 0
		}), 0)
	})
	t.Run("every point lies on the line and the boundary, and exists exactly where IntersectsLine holds for a line with a direction", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, l := range lineFixtures {
				points := c.IntersectionLine(l).AppendTo(nil)
				message := fmt.Sprintf("%s → %s: ", c, l)

				for _, point := range points {
					assert.True(t, l.Contains(point), message+point.String()+" on the line: ")
					assert.True(t, touches(c, point), message+point.String()+" on the boundary: ")
				}
				assert.True(t, len(points) <= 2, message)
				if l.Direction != (Vector[float64]{}) {
					assert.Equal(t, len(points) > 0, c.IntersectsLine(l), message)
				}
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

func TestCircle_Arc(t *testing.T) {
	c := Circ(Pt(1.0, 2.0), 5)

	t.Run("a part of the boundary", func(t *testing.T) {
		geomtest.AssertArc(t, c.Arc(Pi/2, -Pi), Ac(Pt(1.0, 2.0), 5, Pi/2, -Pi))
		geomtest.AssertCircle(t, c.Arc(Pi/2, -Pi).Circle(), c)
	})
	t.Run("the whole boundary", func(t *testing.T) {
		geomtest.AssertNumber(t, c.Arc(0, 2*Pi).Length(), c.Perimeter())
	})
	t.Run("contains the boundary and not the inside", func(t *testing.T) {
		assert.True(t, c.Arc(0, 2*Pi).Contains(Pt(1.0, 7.0)))
		assert.False(t, c.Arc(0, 2*Pi).Contains(c.Center))
		geomtest.AssertNumber(t, c.Arc(0, 2*Pi).DistanceTo(c.Center), 5)
	})
}

func TestCircle_Sector(t *testing.T) {
	c := Circ(Pt(1.0, 2.0), 5)

	t.Run("a slice of the disc", func(t *testing.T) {
		geomtest.AssertSector(t, c.Sector(Pi/2, -Pi), Sect(Pt(1.0, 2.0), 5, Pi/2, -Pi))
		geomtest.AssertCircle(t, c.Sector(Pi/2, -Pi).Circle(), c)
		geomtest.AssertArc(t, c.Sector(Pi/2, -Pi).Arc(), c.Arc(Pi/2, -Pi))
	})
	t.Run("the whole disc", func(t *testing.T) {
		geomtest.AssertNumber(t, c.Sector(0, 2*Pi).Area(), c.Area())
	})
}

func TestCircle_Polygon(t *testing.T) {
	c := Circ(Pt(0.0, 0.0), 5.0)

	t.Run("equal chords from the angle zero", func(t *testing.T) {
		geomtest.AssertPolygon(t, c.Polygon(2), Pol([]Point[float64]{Pt(5.0, 0.0), Pt(0.0, 5.0), Pt(-5.0, 0.0), Pt(0.0, -5.0)}))
		geomtest.AssertPolygon(t, c.Polygon(-2), c.Polygon(2))
	})
	t.Run("never fewer than three", func(t *testing.T) {
		assert.Length(t, c.Polygon(100).Points, 3)
	})
	t.Run("a thousand chords at most", func(t *testing.T) {
		assert.Length(t, c.Polygon(0).Points, 1024)
		assert.Length(t, c.Polygon(math.NaN()).Points, 1024)
	})
	t.Run("int rounds the vertices", func(t *testing.T) {
		geomtest.AssertPolygon(t, Circ(Pt(1, 2), 5).Polygon(2), Pol([]Point[int]{Pt(6, 2), Pt(1, 7), Pt(-4, 2), Pt(1, -3)}))
		geomtest.AssertPolygon(t, Circ(Pt(1, 2), 1).Polygon(0), Pol([]Point[int]{Pt(2, 2), Pt(2, 3), Pt(1, 3), Pt(0, 3), Pt(0, 2), Pt(0, 1), Pt(1, 1), Pt(2, 1)}))
	})
	t.Run("no radius gives the center", func(t *testing.T) {
		geomtest.AssertPolygon(t, Circ(Pt(1, 2), 0).Polygon(1), Pol([]Point[int]{Pt(1, 2)}))
		geomtest.AssertPolygon(t, Circ(Pt(1.0, 2.0), 0.0).Polygon(1), Pol([]Point[float64]{Pt(1.0, 2.0)}))
	})
	t.Run("the vertices lie on the boundary, the edges within the tolerance of it, wound as a rectangle", func(t *testing.T) {
		for _, circle := range circleFixtures {
			for _, tolerance := range []float64{0.01, 0.25} {
				polygon := circle.Polygon(tolerance)
				message := fmt.Sprintf("%s ±%v: ", circle, tolerance)

				for vertex := range polygon.Vertices() {
					assert.True(t, circle.Arc(0, 2*Pi).Contains(vertex), message+vertex.String())
				}

				if circle.Radius == 0 {
					continue
				}

				for edge := range polygon.Edges() {
					assert.True(t, LessOrEqual(circle.Arc(0, 2*Pi).DistanceTo(edge.Midpoint()), tolerance), message+edge.String())
				}

				assert.True(t, polygon.IsConvex(), message)
				assert.Equal(t, polygon.Winding(), WindingClockwise, message)
			}
		}
	})
}

func TestCircle_AppendPolygon(t *testing.T) {
	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, Circ(Pt(0, 0), 5).AppendPolygon([]Point[int]{Pt(9, 9)}, 2), []Point[int]{Pt(9, 9), Pt(5, 0), Pt(0, 5), Pt(-5, 0), Pt(0, -5)})
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		circle := Circ(Pt(0, 0), 50)
		buffer := make([]Point[int], 0, 64)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.Polygon(1).Points
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = circle.AppendPolygon(buffer[:0], 1)
		}), 0)
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

// assertTangents asserts the properties the tangents between two circles have, outer and inner
// alike: at most two, each from the boundary of a to the boundary of b and at a right angle to
// the radius at both ends, where the circle has one, and the same segments reversed when asked
// from b.
func assertTangents(t *testing.T, a, b Circle[float64], tangents, reversed []Segment[float64]) {
	t.Helper()

	message := fmt.Sprintf("%s → %s: ", a, b)

	assert.True(t, len(tangents) <= 2, message)
	assert.Length(t, reversed, len(tangents), message)

	for i, s := range tangents {
		assert.True(t, touches(a, s.Start), message)
		assert.True(t, touches(b, s.End), message)
		assert.True(t, slices.ContainsFunc(reversed, s.Reverse().Equal), message)

		along := s.Vector().Normalize()
		if a.Radius > 0 {
			geomtest.AssertNumber(t, s.Start.Subtract(a.Center).Normalize().Dot(along), 0, fmt.Sprintf("%s#%d at a: ", message, i))
		}
		if b.Radius > 0 {
			geomtest.AssertNumber(t, s.End.Subtract(b.Center).Normalize().Dot(along), 0, fmt.Sprintf("%s#%d at b: ", message, i))
		}
	}
}

// touches reports whether the point lies on the boundary of the circle within the tolerance, by
// the comparison the package makes for an endpoint on the boundary: a zero-length segment at the
// point crosses the boundary exactly there.
func touches[T Number](circle Circle[T], point Point[T]) bool {
	return circle.IntersectionSegment(Seg(point, point)).Count == 1
}

var circleFixtures = []Circle[float64]{
	Circ(Pt(0.0, 0.0), 0.0),
	Circ(Pt(0.0, 0.0), 1.0),
	Circ(Pt(1.0, 2.0), 10.0),
	Circ(Pt(0.6, -0.25), 1.2),
	Circ(Pt(-3.5, 0.25), 7.0),
	Circ(Pt(12.5, -0.1), 0.5),
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

func ExampleCirc() {
	fmt.Println(Circ(Pt(10, 16), 5))
	// Output: Circ((10,16);5)
}

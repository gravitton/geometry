package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestCircle_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(10, 16), 12), Circle[int]{Center: Pt(10, 16), Radius: 12})
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.16, 204), 5.1), Circle[float64]{Center: Pt(0.16, 204.0), Radius: 5.1})
	})
	t.Run("a negative radius is taken absolute", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(10, 16), -12), Circ(Pt(10, 16), 12))
	})
}

func TestCircle_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(1, 2), 10).Area(), Pi*100.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Area(), Pi*1.44)
	})
	t.Run("large integer radius does not overflow", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(0, 0), 3037000500).Area(), Pi*3037000500*3037000500)
	})
}

func TestCircle_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(1, 2), 10).Perimeter(), Pi*20.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Perimeter(), Pi*2.4)
	})
}

func TestCircle_Diameter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(1, 2), 10).Diameter(), 20)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(0.6, -0.25), 1.2).Diameter(), 2.4)
	})
}

func TestCircle_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRectangle(t, Circ(Pt(1, 2), 10).Bounds(), Rect(Pt(1, 2), Sz(20, 20)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, Circ(Pt(0.6, -0.25), 1.2).Bounds(), Rect(Pt(0.6, -0.25), Sz(2.4, 2.4)))
	})
}

func TestCircle_Anchor(t *testing.T) {
	c := Circ(Pt(10.0, 10.0), 5.0)

	t.Run("cardinal directions", func(t *testing.T) {
		AssertPoint(t, c.Anchor(Right), Pt(15.0, 10.0))
		AssertPoint(t, c.Anchor(Left), Pt(5.0, 10.0))
		AssertPoint(t, c.Anchor(Top), Pt(10.0, 5.0))
		AssertPoint(t, c.Anchor(Bottom), Pt(10.0, 15.0))
	})
	t.Run("diagonals land on the boundary", func(t *testing.T) {
		// unlike Rectangle, whose diagonals reach the corners
		AssertNumber(t, c.Center.DistanceTo(c.Anchor(DirectionUpRight)), c.Radius)
	})
	t.Run("none is the center", func(t *testing.T) {
		AssertPoint(t, c.Anchor(DirectionNone), Pt(10.0, 10.0))
	})
}

func TestCircle_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Translate(Vec(3, -2)), Circ(Pt(4, 0), 10))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Translate(Vec(100.1, -0.1)), Circ(Pt(100.7, -0.35), 1.2))
	})
}

func TestCircle_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).MoveTo(Pt(3, -2)), Circ(Pt(3, -2), 10))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).MoveTo(Pt(100.1, -0.1)), Circ(Pt(100.1, -0.1), 1.2))
	})
}

func TestCircle_Scale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Scale(2.5), Circ(Pt(1, 2), 25))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Scale(2.5), Circ(Pt(0.6, -0.25), 3.0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Scale(-2), Circ(Pt(1, 2), 20))
	})
}

func TestCircle_Unscale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 25).Unscale(2.5), Circ(Pt(1, 2), 10))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 3.0).Unscale(2.5), Circ(Pt(0.6, -0.25), 1.2))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 20).Unscale(-2), Circ(Pt(1, 2), 10))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Circ(Pt(1, 2), 10).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestCircle_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Resize(8), Circ(Pt(1, 2), 8))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Resize(3.1), Circ(Pt(0.6, -0.25), 3.1))
	})
	t.Run("a negative radius is taken absolute", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Resize(-8), Circ(Pt(1, 2), 8))
	})
}

func TestCircle_Canonical(t *testing.T) {
	t.Run("takes a literal negative radius absolute and keeps the center", func(t *testing.T) {
		AssertCircle(t, Circle[int]{Pt(1, 2), -8}.Canonical(), Circ(Pt(1, 2), 8))
		AssertCircle(t, Circle[float64]{Pt(0.6, -0.25), -1.2}.Canonical(), Circ(Pt(0.6, -0.25), 1.2))
	})
	t.Run("is the circle Circ builds", func(t *testing.T) {
		c := Circle[int]{Pt(1, 2), -8}

		AssertCircle(t, c.Canonical(), Circ(c.Center, c.Radius))
		assert.True(t, c.Canonical().Contains(c.Canonical().Anchor(Right)))
	})
	t.Run("repairs decoded JSON", func(t *testing.T) {
		var c Circle[int]

		assert.Nil(t, json.Unmarshal([]byte(`{"x":1,"y":2,"r":-8}`), &c))
		AssertCircle(t, c.Canonical(), Circ(Pt(1, 2), 8))
	})
	t.Run("a well-formed circle is unchanged", func(t *testing.T) {
		for _, c := range circleFixtures {
			AssertCircle(t, c.Canonical(), c, c.String())
		}
	})
}

func TestCircle_Grow(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Grow(8), Circ(Pt(1, 2), 18))
		AssertCircle(t, Circ(Pt(1, 2), 10).Grow(-12), Circ(Pt(1, 2), 0))
	})
	t.Run("float", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Grow(3.1), Circ(Pt(0.6, -0.25), 4.3))
	})
}

func TestCircle_Shrink(t *testing.T) {
	t.Run("reduces the radius", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(8), Circ(Pt(1, 2), 2))
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Shrink(0.3), Circ(Pt(0.6, -0.25), 0.9))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(100), Circ(Pt(1, 2), 0))
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Shrink(5.0), Circ(Pt(0.6, -0.25), 0.0))
	})
}

func TestCircle_Lerp(t *testing.T) {
	a, b := Circ(Pt(0.0, 0.0), 2.0), Circ(Pt(10.0, 20.0), 8.0)

	t.Run("moves the center and the radius together", func(t *testing.T) {
		AssertCircle(t, a.Lerp(b, 0.5), Circ(Pt(5.0, 10.0), 5.0))
	})
	t.Run("the ends are the circles themselves", func(t *testing.T) {
		AssertCircle(t, a.Lerp(b, 0), a)
		AssertCircle(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates and keeps the radius absolute", func(t *testing.T) {
		AssertCircle(t, a.Lerp(b, 2), Circ(Pt(20.0, 40.0), 14.0))
		AssertCircle(t, Circ(Pt(0.0, 0.0), 2.0).Lerp(Circ(Pt(0.0, 0.0), 0.0), 2), Circ(Pt(0.0, 0.0), 2.0))
	})
}

func TestCircle_Rotate(t *testing.T) {
	circle := Circ(Pt(1.0, 2.0), 3.0)

	t.Run("no angle moves a circle", func(t *testing.T) {
		AssertCircle(t, circle.Rotate(Pi/3), circle)
		AssertCircle(t, circle.Rotate(-Pi), circle)
	})
	t.Run("the bounds turn with nothing", func(t *testing.T) {
		for _, c := range circleFixtures {
			AssertCircle(t, c.Rotate(Pi/7), c, c.String())
		}
	})
}

func TestCircle_AlignTo(t *testing.T) {
	c := Circ(Pt(10.0, 10.0), 5.0)

	t.Run("moves the anchor onto the point", func(t *testing.T) {
		AssertCircle(t, c.AlignTo(Bottom, Pt(0.0, 0.0)), Circ(Pt(0.0, -5.0), 5.0))
		AssertCircle(t, Circ(Pt(1, 2), 3).AlignTo(Left, Pt(0, 0)), Circ(Pt(3, 0), 3))
	})
	t.Run("none aligns the center like MoveTo", func(t *testing.T) {
		AssertCircle(t, c.AlignTo(DirectionNone, Pt(1.0, 2.0)), c.MoveTo(Pt(1.0, 2.0)))
	})
	t.Run("is the inverse of Anchor", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, direction := range Directions() {
				AssertPoint(t, c.AlignTo(direction, Pt(1.5, -2.5)).Anchor(direction), Pt(1.5, -2.5), fmt.Sprintf("%s %s: ", c, direction))
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
				assert.Equal(t, c.Contains(p), Ln(p, p).IntersectsCircle(c), fmt.Sprintf("%s → %s: ", c, p))
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
		AssertNumber(t, circle.DistanceTo(Pt(7, 0)), 4.0)
		AssertNumber(t, circle.DistanceTo(Pt(3, 4)), 2.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, circle.DistanceTo(Pt(1, 1)), 0.0)
		assert.Equal(t, circle.DistanceTo(Pt(3, 0)), 0.0)
		assert.Equal(t, circle.DistanceTo(circle.Center), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Circ(Pt(0.0, 0.0), 1.0).DistanceTo(Pt(1.0, 1.0)), Sqrt2-1)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		c := Circ(Pt(0.0, 0.0), 1.0)

		assert.Equal(t, c.DistanceTo(Pt(1.0+Delta/2, 0.0)), 0.0)
		AssertNumber(t, c.DistanceTo(Pt(1.0+2*Delta, 0.0)), 2*Delta)
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
		AssertNumber(t, circle.DistanceSquaredTo(Pt(7, 0)), 16.0)
		AssertNumber(t, circle.DistanceSquaredTo(Pt(3, 4)), 4.0)
		assert.Equal(t, circle.DistanceSquaredTo(Pt(1, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, p := range pointFixtures {
				AssertNumber(t, c.DistanceSquaredTo(p), c.DistanceTo(p)*c.DistanceTo(p), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
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
		AssertVertices(t, circle.IntersectionCircle(Circ(Pt(6.0, 0.0), 5.0)), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, -4.0)})
		AssertVertices(t, circle.IntersectionCircle(Circ(Pt(0.0, 6.0), 5.0)), []Point[float64]{Pt(-4.0, 3.0), Pt(4.0, 3.0)})
	})
	t.Run("a larger circle behind the first", func(t *testing.T) {
		height := math.Sqrt(1 - 0.125*0.125)

		AssertVertices(t, Circ(Pt(0.0, 0.0), 1.0).IntersectionCircle(Circ(Pt(1.0, 0.0), 1.5)), []Point[float64]{Pt(-0.125, height), Pt(-0.125, -height)})
	})
	t.Run("tangent from outside gives one point", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionCircle(Circ(Pt(8.0, 0.0), 3.0)), []Point[float64]{Pt(5.0, 0.0)})
	})
	t.Run("tangent from inside gives one point", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionCircle(Circ(Pt(2.0, 0.0), 3.0)), []Point[float64]{Pt(5.0, 0.0)})
		AssertVertices(t, Circ(Pt(0.0, 0.0), 3.0).IntersectionCircle(Circ(Pt(2.0, 0.0), 5.0)), []Point[float64]{Pt(-3.0, 0.0)})
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
		assert.Equal(t, len(circle.IntersectionCircle(Circ(Pt(8.0+Delta/2, 0.0), 3.0))), 1)
		assert.Nil(t, circle.IntersectionCircle(Circ(Pt(8.0+2*Delta, 0.0), 3.0)))
	})
	t.Run("a tangent exactly Delta outside is judged like Intersects, on both sides", func(t *testing.T) {
		for _, other := range []Circle[float64]{Circ(Pt(8.0+Delta, 0.0), 3.0), Circ(Pt(2.0+Delta, 0.0), 3.0)} {
			assert.True(t, circle.IntersectsCircle(other))
			AssertVertices(t, circle.IntersectionCircle(other), []Point[float64]{Pt(5.0, 0.0)})
		}
	})
	t.Run("a tangent within the tolerance lands on both boundaries", func(t *testing.T) {
		a, b := Circ(Pt(0.0, 0.0), 6.0), Circ(Pt(2.0000005, 0.0), 4.0)
		points := a.IntersectionCircle(b)

		assert.Equal(t, len(points), 1)
		for _, p := range points {
			assert.True(t, a.touchesSquared(a.Center.DistanceSquaredTo(p)), "on a")
			assert.True(t, b.touchesSquared(b.Center.DistanceSquaredTo(p)), "on b")
		}
	})
	t.Run("int rounds the points", func(t *testing.T) {
		AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)), []Point[int]{Pt(3, 4), Pt(3, -4)})
		AssertVertices(t, Circ(Pt(0, 0), 2).IntersectionCircle(Circ(Pt(3, 0), 2)), []Point[int]{Pt(2, 1), Pt(2, -1)})
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
					assert.True(t, a.touchesSquared(a.Center.DistanceSquaredTo(p)), fmt.Sprintf("%s → %s: %s on a: ", a, b, p))
					assert.True(t, b.touchesSquared(b.Center.DistanceSquaredTo(p)), fmt.Sprintf("%s → %s: %s on b: ", a, b, p))
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
			assert.True(t, a.touchesSquared(a.Center.Float().DistanceSquaredTo(p.Float())), fmt.Sprintf("%s → %s: %s on the first: ", a, b, p))
			assert.True(t, b.touchesSquared(b.Center.Float().DistanceSquaredTo(p.Float())), fmt.Sprintf("%s → %s: %s on the second: ", a, b, p))
		}

		if len(points) == 2 {
			assert.True(t, !points[0].Equal(points[1]), fmt.Sprintf("%s → %s: two distinct points: ", a, b))
		}

		if len(points) > 0 {
			assert.True(t, a.IntersectsCircle(b), fmt.Sprintf("%s → %s: points imply intersects: ", a, b))
		}
	})
}

func TestCircle_IntersectsLine(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, circle.IntersectsLine(Ln(Pt(-2.0, 2.0), Pt(2.0, 2.0))))
		assert.False(t, circle.IntersectsLine(Ln(Pt(2.0, 0.0), Pt(3.0, 0.0))))
	})
	t.Run("tangent counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(-2.0, 1.0), Pt(2.0, 1.0))))
		assert.False(t, circle.IntersectsLine(Ln(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta))))
	})
	t.Run("an endpoint inside counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(0.5, 0.0), Pt(5.0, 0.0))))
	})
	t.Run("a segment inside counts", func(t *testing.T) {
		assert.True(t, circle.IntersectsLine(Ln(Pt(-0.5, 0.0), Pt(0.5, 0.0))))
	})
}

func TestCircle_IntersectionLine(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(2.0, 0.0), Pt(-2.0, 0.0))), []Point[float64]{Pt(1.0, 0.0), Pt(-1.0, 0.0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-2.0, 0.0), Pt(0.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0)})
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(0.5, 0.0), Pt(5.0, 0.0))), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-2.0, 1.0), Pt(2.0, 1.0))), []Point[float64]{Pt(0.0, 1.0)})
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-2.0, 1.0+Delta/2), Pt(2.0, 1.0+Delta/2))), []Point[float64]{Pt(0.0, 1.0+Delta/2)})
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta))))
	})
	t.Run("a chord within the tolerance band of the gap keeps both ends", func(t *testing.T) {
		height := 1.0 - Delta/2
		half := math.Sqrt(1 - height*height)

		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-1.0, height), Pt(1.0, height))), []Point[float64]{Pt(-half, height), Pt(half, height)})
	})
	t.Run("a chord shorter than Delta is a tangent", func(t *testing.T) {
		height := math.Sqrt(1 - Delta*Delta/16)

		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-1.0, height), Pt(1.0, height))), []Point[float64]{Pt(0.0, height)})
	})
	t.Run("a segment within the tolerance with both ends on the boundary gives one point", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(1.0, 0.0), Pt(1.0, Delta/10))), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("a tangent beyond the segment is missed", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(1.0, 1.0), Pt(2.0, 1.0))))
	})
	t.Run("a shallow touch is decided on the endpoint distance, like IntersectsLine", func(t *testing.T) {
		shallow := Ln(Pt(-1.0, 1.0+Delta/2), Pt(-0.0005, 1.0+Delta/2))

		assert.True(t, circle.IntersectsLine(shallow))
		AssertVertices(t, circle.IntersectionLine(shallow), []Point[float64]{shallow.End})
		AssertVertices(t, circle.IntersectionLine(shallow.Reverse()), []Point[float64]{shallow.End})
	})
	t.Run("an interior graze exactly Delta outside is judged like IntersectsLine", func(t *testing.T) {
		l, c := Ln(Pt(-1.0, 1.000001), Pt(1.0, 1.000001)), Circ(Pt(0.0, 0.0), 1.0)

		assert.Equal(t, len(c.IntersectionLine(l)) > 0, c.IntersectsLine(l))
		for _, y := range []float64{1.0000009, 1.0000011, 1.0000015} {
			l := Ln(Pt(-1.0, y), Pt(1.0, y))

			assert.Equal(t, len(c.IntersectionLine(l)) > 0, c.IntersectsLine(l), l.String())
		}
	})
	t.Run("an endpoint exactly Delta outside is judged like IntersectsLine", func(t *testing.T) {
		l, c := Ln(Pt(-1.0, 2.000001), Pt(-882.0315, 56.11111116666667)), Circ(Pt(-1.0, 0.0), 2.0)

		assert.Equal(t, len(c.IntersectionLine(l)) > 0, c.IntersectsLine(l))
	})
	t.Run("an endpoint on the boundary is counted once with its crossing", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-1.0, 0.0), Pt(2.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(-2.0, 0.0), Pt(1.0, 0.0))), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
	})
	t.Run("an endpoint on the boundary replaces the crossing nearest to it", func(t *testing.T) {
		l, c := Ln(Pt(-1.0, 2.000001), Pt(-877.0315, 0.11111116666666668)), Circ(Pt(-1.0, 0.0), 2.0)
		points := c.IntersectionLine(l)

		assert.Equal(t, len(points), 2)
		AssertPoint(t, points[0], l.Start)
	})
	t.Run("a tangent segment with both ends on the boundary gives its ends", func(t *testing.T) {
		grazing := Ln(Pt(-0.0003, 1.0), Pt(0.0003, 1.0))

		AssertVertices(t, circle.IntersectionLine(grazing), []Point[float64]{grazing.Start, grazing.End})
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(-2.0, 2.0), Pt(2.0, 2.0))))
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(2.0, 0.0), Pt(3.0, 0.0))))
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(-0.5, 0.0), Pt(0.5, 0.0))))
	})
	t.Run("a degenerate segment is a point on the boundary or nothing", func(t *testing.T) {
		AssertVertices(t, circle.IntersectionLine(Ln(Pt(1.0, 0.0), Pt(1.0, 0.0))), []Point[float64]{Pt(1.0, 0.0)})
		assert.Nil(t, circle.IntersectionLine(Ln(Pt(0.5, 0.0), Pt(0.5, 0.0))))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		AssertVertices(t, Circ(Pt(0, 0), 5).IntersectionLine(Ln(Pt(-5, -5), Pt(5, 5))), []Point[int]{Pt(-4, -4), Pt(4, 4)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsLine holds", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				points := c.IntersectionLine(l)

				assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings: ", l, c))
				for _, p := range points {
					assert.True(t, l.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", l, c, p))
					assert.True(t, c.touchesSquared(c.Center.DistanceSquaredTo(p)), fmt.Sprintf("%s → %s: %s on the boundary: ", l, c, p))
				}
				if len(points) > 0 {
					assert.True(t, c.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, c))
				}
			}
		}
	})
}

func FuzzCircle_IntersectionLine(f *testing.F) {
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

		l, c := Ln(Pt(x1, y1), Pt(x2, y2)), Circ(Pt(cx, cy), r)
		points := c.IntersectionLine(l)

		assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings, got %d: ", l, c, len(points)))

		for _, p := range points {
			assert.True(t, l.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", l, c, p))
			assert.True(t, c.touchesSquared(c.Center.Float().DistanceSquaredTo(p.Float())), fmt.Sprintf("%s → %s: %s on the boundary: ", l, c, p))
		}

		if inside := c.Contains(l.Start) && c.Contains(l.End); !inside {
			assert.Equal(t, len(points) > 0, c.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, c))
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
				expected := p.Contains(c.Center) || slices.ContainsFunc(slices.Collect(p.Edges()), func(edge Line[float64]) bool {
					return c.IntersectsLine(edge)
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
		assert.True(t, Circ(Pt(1.5, 1.5), 0.5).IntersectsRectangle(diamond.Bounds()))
		assert.True(t, Circ(Pt(1.0, 1.0), 0.5).IntersectsRectangle(diamond))
	})
	t.Run("a circle intersects its own bounds", func(t *testing.T) {
		for _, c := range circleFixtures {
			if c.Radius == 0 {
				continue // a degenerate circle touches nothing
			}

			assert.True(t, c.IntersectsRectangle(c.Bounds()), fmt.Sprintf("%s: ", c))
		}
	})
}

func TestCircle_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)

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
		assert.False(t, Circ(Pt(0, 0), 1).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, c.IntersectsRegularPolygon(rp), c.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", rp, c))
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

func TestCircle_Cast(t *testing.T) {
	c := Circ(Pt(1.5, -2.5), 3.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		AssertCircle(t, c.Cast[int](), c.Int())
		AssertCircle(t, c.Cast[float64](), c.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		AssertCircle(t, c.Cast[int8](), Circ(Pt[int8](2, -3), 4))
	})
}

func TestCircle_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Int(), Circ(Pt(1, 2), 10))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Int(), Circ(Pt(1, 0), 1))
	})
}

func TestCircle_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(1, 2), 10).Float(), Circ(Pt(1.0, 2.0), 10.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Float(), Circ(Pt(0.6, -0.25), 1.2))
	})
}

func TestCircle_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(10, 16), 5).String(), "Circ((10,16);5)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(100, -34.0000115), 0.2).String(), "Circ((100.00,-34.00);0.20)")
	})
}

func TestCircle_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Circ(Pt(10, 16), 12), `{"x":10,"y":16,"r":12}`)

		var c Circle[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16,"r":12}`), &c))
		AssertCircle(t, c, Circ(Pt(10, 16), 12))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Circ(Pt(100, -34.0000115), 0.2), `{"x":100.0,"y":-34.0000115,"r":0.2}`)

		var c Circle[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115,"r":0.2}`), &c))
		AssertCircle(t, c, Circ(Pt(10.1, -34.0000115), 0.2))
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
			AssertNumber(t, c.Area(), Pi*c.Radius*c.Radius, fmt.Sprintf("%s: ", c))
			AssertNumber(t, c.Perimeter(), 2*Pi*c.Radius, fmt.Sprintf("%s: ", c))
			AssertNumber(t, c.Diameter(), 2*c.Radius, fmt.Sprintf("%s: ", c))
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

			assert.True(t, bounds.Center.Equal(c.Center), fmt.Sprintf("%s: ", c))
			assert.True(t, bounds.Size.Equal(SzU(c.Diameter())), fmt.Sprintf("%s: ", c))
		}
	})
	t.Run("every anchor is on the boundary", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, direction := range Directions() {
				anchor := c.Anchor(direction)

				AssertNumber(t, c.Center.DistanceTo(anchor), c.Radius, fmt.Sprintf("%s → %s: ", c, direction))
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
				AssertNumber(t, c.MoveTo(point).Radius, c.Radius, fmt.Sprintf("%s → %s: ", c, point))
				assert.True(t, c.MoveTo(point).Center.Equal(point), fmt.Sprintf("%s → %s: ", c, point))
				AssertNumber(t, c.Translate(point.Vector()).Radius, c.Radius, fmt.Sprintf("%s → %s: ", c, point))
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

	AssertCircle(t, c, Circ(Pt(1, 2), 10))
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

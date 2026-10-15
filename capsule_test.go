package geom_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestCapsule_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2), Capsule[int]{Segment: Seg(Pt(0, 0), Pt(10, 0)), Radius: 2})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0.16, 204), Pt(5.1, 1.2), 0.5), Capsule[float64]{Segment: Seg(Pt(0.16, 204.0), Pt(5.1, 1.2)), Radius: 0.5})
	})
	t.Run("a negative radius is taken absolute", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), -2), Cap(Pt(0, 0), Pt(10, 0), 2))
	})
}

func TestCapsule_Centroid(t *testing.T) {
	t.Run("the midpoint of the segment", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(0, 0), Pt(10, 4), 2).Centroid(), Pt(5, 2))
	})
	t.Run("int rounds as the midpoint does", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(0, 0), Pt(3, 0), 2).Centroid(), Seg(Pt(0, 0), Pt(3, 0)).Midpoint())
	})
}

func TestCapsule_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 2).Area(), 4*Pi+40)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(1.0, 2.0), Pt(4.0, 6.0), 0.5).Area(), 0.25*Pi+5)
	})
	t.Run("a radius of zero encloses nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 0).Area(), 0.0)
	})
}

func TestCapsule_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 2).Perimeter(), 4*Pi+20)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(1.0, 2.0), Pt(4.0, 6.0), 0.5).Perimeter(), Pi+10)
	})
}

func TestCapsule_Inertia(t *testing.T) {
	t.Run("a segment of no length is the disc", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(1, 2), Pt(1, 2), 3).Inertia(), Circ(Pt(1, 2), 3).Inertia())
	})
	t.Run("a radius of zero encloses nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 0).Inertia(), 0.0)
	})
	t.Run("the turn does not change it", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0.0, 0.0), Pt(3.0, 4.0), 2.0).Inertia(), Cap(Pt(0.0, 0.0), Pt(5.0, 0.0), 2.0).Inertia())
	})
	t.Run("that of the outline walked as a polygon", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			outline := outlineOf(c, 4096)

			assert.True(t, math.Abs(c.Area()-outline.Area()) <= 1e-5*c.Area(), fmt.Sprintf("%s: %v ≠ %v", c, c.Area(), outline.Area()))
			assert.True(t, math.Abs(c.Inertia()-outline.Inertia()) <= 1e-5*c.Inertia(), fmt.Sprintf("%s: %v ≠ %v", c, c.Inertia(), outline.Inertia()))
		}
	})
}

// outlineOf returns the polygon of n points of the capsule boundary, its Support in n
// directions around: the outline the mass properties are checked against.
func outlineOf(c Capsule[float64], n int) Polygon[float64] {
	points := make([]Point[float64], n)
	for i := range points {
		points[i] = c.Support(VectorFromAngle(2*Pi*float64(i)/float64(n), 1.0))
	}

	return Pol(points)
}

func TestCapsule_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertBox(t, Cap(Pt(0, 0), Pt(10, 0), 2).Bounds(), Bx(Pt(-2, -2), Pt(12, 2)))
	})
	t.Run("float runs either way", func(t *testing.T) {
		geomtest.AssertBox(t, Cap(Pt(4.0, 6.0), Pt(1.0, 2.0), 0.5).Bounds(), Bx(Pt(0.5, 1.5), Pt(4.5, 6.5)))
	})
}

func TestCapsule_BoundingCircle(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Cap(Pt(0, 0), Pt(10, 0), 2).BoundingCircle(), Circ(Pt(5, 0), 7))
	})
	t.Run("int rounds the radius up from the rounded midpoint", func(t *testing.T) {
		geomtest.AssertCircle(t, Cap(Pt(0, 0), Pt(3, 0), 2).BoundingCircle(), Circ(Pt(2, 0), 4))
	})
	t.Run("contains the boundary", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			circle := c.BoundingCircle()

			for _, direction := range vectorFixtures {
				assert.True(t, circle.Contains(c.Support(direction)), fmt.Sprintf("%s → %s: ", c, direction))
			}
		}
	})
}

func TestCapsule_Translate(t *testing.T) {
	geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Translate(Vec(3, -2)), Cap(Pt(3, -2), Pt(13, -2), 2))
}

func TestCapsule_MoveTo(t *testing.T) {
	geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).MoveTo(Pt(0, 3)), Cap(Pt(-5, 3), Pt(5, 3), 2))
}

func TestCapsule_Rotate(t *testing.T) {
	t.Run("turns about the midpoint", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(-2.0, 0.0), Pt(2.0, 0.0), 1.0).Rotate(Pi/2), Cap(Pt(0.0, -2.0), Pt(0.0, 2.0), 1.0))
	})
	t.Run("a circle stays where it is", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(1.0, 2.0), Pt(1.0, 2.0), 3.0).Rotate(1.0), Cap(Pt(1.0, 2.0), Pt(1.0, 2.0), 3.0))
	})
}

func TestCapsule_Scale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(-2, 0), Pt(2, 0), 2).Scale(2.5), Cap(Pt(-5, 0), Pt(5, 0), 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0.0, 1.0), Pt(4.0, 1.0), 0.5).Scale(0.5), Cap(Pt(1.0, 1.0), Pt(3.0, 1.0), 0.25))
	})
	t.Run("a negative factor swaps the ends and keeps the radius positive", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(-2, 0), Pt(2, 0), 2).Scale(-2), Cap(Pt(4, 0), Pt(-4, 0), 4))
	})
}

func TestCapsule_Unscale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(-5, 0), Pt(5, 0), 5).Unscale(2.5), Cap(Pt(-2, 0), Pt(2, 0), 2))
	})
	t.Run("a negative factor swaps the ends and keeps the radius positive", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(-4, 0), Pt(4, 0), 4).Unscale(-2), Cap(Pt(2, 0), Pt(-2, 0), 2))
	})
	t.Run("a zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Cap(Pt(-4, 0), Pt(4, 0), 4).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestCapsule_Grow(t *testing.T) {
	t.Run("the radius grows and the segment stays", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Grow(3), Cap(Pt(0, 0), Pt(10, 0), 5))
	})
	t.Run("a negative amount clamps at zero", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Grow(-3), Cap(Pt(0, 0), Pt(10, 0), 0))
	})
}

func TestCapsule_Shrink(t *testing.T) {
	t.Run("the radius shrinks and the segment stays", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 5).Shrink(3), Cap(Pt(0, 0), Pt(10, 0), 2))
	})
	t.Run("clamps at zero, the segment itself", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Shrink(3), Cap(Pt(0, 0), Pt(10, 0), 0))
	})
}

func TestCapsule_Lerp(t *testing.T) {
	a, b := Cap(Pt(0, 0), Pt(10, 0), 2), Cap(Pt(0, 10), Pt(10, 20), 4)

	t.Run("the ends and the radius move together", func(t *testing.T) {
		geomtest.AssertCapsule(t, a.Lerp(b, 0.5), Cap(Pt(0, 5), Pt(10, 10), 3))
	})
	t.Run("ends on the two capsules", func(t *testing.T) {
		for _, a := range capsuleFixtures {
			for _, b := range capsuleFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
	t.Run("an extrapolation past a zero radius grows again", func(t *testing.T) {
		geomtest.AssertNumber(t, a.Lerp(b, -2).Radius, 2)
	})
}

func TestCapsule_Canonical(t *testing.T) {
	t.Run("a well-formed capsule is unchanged", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Canonical(), Cap(Pt(0, 0), Pt(10, 0), 2))
	})
	t.Run("a negative radius is repaired", func(t *testing.T) {
		geomtest.AssertCapsule(t, Capsule[int]{Seg(Pt(0, 0), Pt(10, 0)), -2}.Canonical(), Cap(Pt(0, 0), Pt(10, 0), 2))
	})
}

func TestCapsule_DistanceTo(t *testing.T) {
	c := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("zero inside", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceTo(Pt(5.0, 1.0)), 0.0)
		geomtest.AssertNumber(t, c.DistanceTo(Pt(11.0, 1.0)), 0.0)
	})
	t.Run("beside the segment", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceTo(Pt(5.0, 6.0)), 4.0)
	})
	t.Run("past an end it is measured to the round end", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceTo(Pt(15.0, 0.0)), 3.0)
		geomtest.AssertNumber(t, c.DistanceTo(Pt(14.0, 3.0)), 3.0)
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 2).DistanceTo(Pt(14, 3)), 3.0)
	})
}

func TestCapsule_DistanceSquaredTo(t *testing.T) {
	c := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("the square of the distance", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceSquaredTo(Pt(5.0, 6.0)), 16.0)
	})
	t.Run("a point within the tolerance of the boundary is on it", func(t *testing.T) {
		assert.Equal(t, c.DistanceSquaredTo(Pt(5.0, 2.0+Delta/2)), 0.0)
	})
	t.Run("int has no tolerance", func(t *testing.T) {
		assert.True(t, Cap(Pt(0, 0), Pt(10, 0), 2).DistanceSquaredTo(Pt(5, 3)) > 0)
	})
}

func TestCapsule_DistanceToCapsule(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("the distance between the segments less the radii", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceToCapsule(Cap(Pt(0, 7), Pt(10, 7), 1)), 4.0)
		geomtest.AssertNumber(t, c.DistanceToCapsule(Cap(Pt(13, 4), Pt(20, 4), 1)), 2.0)
	})
	t.Run("zero for capsules that overlap, touch or lie one within the other", func(t *testing.T) {
		geomtest.AssertNumber(t, c.DistanceToCapsule(Cap(Pt(5, -9), Pt(5, 9), 1)), 0.0)
		geomtest.AssertNumber(t, c.DistanceToCapsule(Cap(Pt(0, 4), Pt(10, 4), 2)), 0.0)
		geomtest.AssertNumber(t, c.DistanceToCapsule(Cap(Pt(4, 0), Pt(6, 0), 1)), 0.0)
	})
	t.Run("zero exactly where the capsules intersect, and the same from either side", func(t *testing.T) {
		for _, a := range capsuleFixtures {
			for _, b := range capsuleFixtures {
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, a.DistanceToCapsule(b) == 0, a.IntersectsCapsule(b), message)
				geomtest.AssertNumber(t, a.DistanceToCapsule(b), b.DistanceToCapsule(a), message)
				assert.True(t, a.DistanceToCapsule(b) >= 0, message)
			}
		}
	})
	t.Run("capsules of no length are measured as their circles", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				geomtest.AssertNumber(t, Cap(a.Center, a.Center, a.Radius).DistanceToCapsule(Cap(b.Center, b.Center, b.Radius)), a.DistanceToCircle(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("a NaN coordinate is at NaN distance", func(t *testing.T) {
		assert.True(t, math.IsNaN(Cap(Pt(math.NaN(), 0.0), Pt(1.0, 0.0), 1.0).DistanceToCapsule(Cap(Pt(5.0, 5.0), Pt(9.0, 5.0), 1.5))))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		other := Cap(Pt(0, 7), Pt(10, 7), 1)

		assertAllocations(t, func() {
			sinkBool = c.DistanceToCapsule(other) > 0
		}, 0)
	})
}

func TestCapsule_DistanceSquaredToCapsule(t *testing.T) {
	geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 2).DistanceSquaredToCapsule(Cap(Pt(0, 7), Pt(10, 7), 1)), 16.0)
	geomtest.AssertNumber(t, Cap(Pt(0, 0), Pt(10, 0), 2).DistanceSquaredToCapsule(Cap(Pt(0, 4), Pt(10, 4), 2)), 0.0)
}

func TestCapsule_Nearest(t *testing.T) {
	c := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("a contained point is itself", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Nearest(Pt(5.0, 1.0)), Pt(5.0, 1.0))
	})
	t.Run("beside the segment it is the point of the side", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Nearest(Pt(5.0, 6.0)), Pt(5.0, 2.0))
	})
	t.Run("past an end it is the point of the round end", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Nearest(Pt(14.0, 3.0)), Pt(11.6, 1.2))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(0, 0), Pt(10, 0), 2).Nearest(Pt(14, 3)), Pt(12, 1))
	})
	t.Run("a radius of zero gives the foot on the segment", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 0.0).Nearest(Pt(5.0, 6.0)), Pt(5.0, 0.0))
	})
}

func TestCapsule_Support(t *testing.T) {
	c := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("the end of the segment along the direction, moved by the radius", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Support(Vec(3.0, 0.0)), Pt(12.0, 0.0))
		geomtest.AssertPoint(t, c.Support(Vec(-1.0, 0.0)), Pt(-2.0, 0.0))
		geomtest.AssertPoint(t, c.Support(Vec(1.0, 1.0)), Pt(10.0+math.Sqrt2, math.Sqrt2))
	})
	t.Run("across the segment it is a point of a side", func(t *testing.T) {
		geomtest.AssertPoint(t, c.Support(Vec(0.0, -3.0)), Pt(0.0, -2.0))
	})
	t.Run("a segment of no length answers as the circle", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(1.0, 2.0), Pt(1.0, 2.0), 5.0).Support(Vec(3.0, 4.0)), Circ(Pt(1.0, 2.0), 5.0).Support(Vec(3.0, 4.0)))
	})
	t.Run("int is not rounded", func(t *testing.T) {
		geomtest.AssertPoint(t, Cap(Pt(0, 0), Pt(10, 0), 1).Support(Vec(1.0, 1.0)), Pt(10+math.Sqrt2/2, math.Sqrt2/2))
	})
	t.Run("no point of the boundary reaches farther", func(t *testing.T) {
		for _, capsule := range capsuleFixtures {
			for _, direction := range vectorFixtures {
				support := capsule.Support(direction)
				message := fmt.Sprintf("%s → %s: ", capsule, direction)

				assert.True(t, capsule.Contains(support), message)
				for _, other := range vectorFixtures {
					assert.True(t, LessOrEqual(capsule.Support(other).Vector().Dot(direction), support.Vector().Dot(direction)), message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		capsule, direction := Cap(Pt(1, 2), Pt(4, 2), 2), Vec(tenth, 0.7)

		assertAllocations(t, func() {
			sinkBool = capsule.Support(direction).X > 0
		}, 0)
	})
}

func TestCapsule_Contains(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, c.Contains(Pt(5, 1)))
		assert.True(t, c.Contains(Pt(11, 1)))
	})
	t.Run("the boundary is inside", func(t *testing.T) {
		assert.True(t, c.Contains(Pt(5, 2)))
		assert.True(t, c.Contains(Pt(12, 0)))
		assert.True(t, c.Contains(Pt(-2, 0)))
	})
	t.Run("outside", func(t *testing.T) {
		assert.False(t, c.Contains(Pt(5, 3)))
		assert.False(t, c.Contains(Pt(13, 0)))
	})
	t.Run("the corner of the box around a round end is outside", func(t *testing.T) {
		assert.False(t, c.Contains(Pt(12, 2)))
	})
	t.Run("float within the tolerance of the boundary counts, beyond it not", func(t *testing.T) {
		assert.True(t, c.Float().Contains(Pt(5.0, 2.0+Delta/2)))
		assert.False(t, c.Float().Contains(Pt(5.0, 2.0+2*Delta)))
	})
	t.Run("agrees with the circle of its radius about the point meeting its segment", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, c.Contains(p), Circ(p, c.Radius).IntersectsSegment(c.Segment), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
	})
}

func BenchmarkCapsule_Contains(b *testing.B) {
	c := Cap(Pt(0.0, 0.0), Pt(50.0, 30.0), 5.0)
	inside, outside := Pt(10.0, 8.0), Pt(99.0, 99.0)

	b.Run("inside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = c.Contains(inside)
		}
	})
	b.Run("outside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = c.Contains(outside)
		}
	})
}

func TestCapsule_ContainsCircle(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 5)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, c.ContainsCircle(Circ(Pt(5, 0), 3)))
		assert.True(t, c.ContainsCircle(Circ(Pt(12, 0), 0)))
	})
	t.Run("touching the boundary from inside counts", func(t *testing.T) {
		assert.True(t, c.ContainsCircle(Circ(Pt(5, 2), 3)))
		assert.True(t, c.ContainsCircle(Circ(Pt(12, 0), 3)))
		assert.False(t, c.ContainsCircle(Circ(Pt(5, 3), 3)))
	})
	t.Run("larger or apart", func(t *testing.T) {
		assert.False(t, c.ContainsCircle(Circ(Pt(5, 0), 6)))
		assert.False(t, c.ContainsCircle(Circ(Pt(30, 0), 1)))
	})
	t.Run("float within the tolerance of the boundary counts, beyond it not", func(t *testing.T) {
		assert.True(t, c.Float().ContainsCircle(Circ(Pt(5.0, 2.0+Delta/2), 3.0)))
		assert.False(t, c.Float().ContainsCircle(Circ(Pt(5.0, 2.0+2*Delta), 3.0)))
	})
}

func TestCapsule_ContainsCapsule(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 5)

	t.Run("inside, itself included", func(t *testing.T) {
		assert.True(t, c.ContainsCapsule(Cap(Pt(2, 0), Pt(8, 0), 3)))
		assert.True(t, c.ContainsCapsule(c))
	})
	t.Run("touching the boundary from inside counts", func(t *testing.T) {
		assert.True(t, c.ContainsCapsule(Cap(Pt(0, 3), Pt(10, 3), 2)))
		assert.False(t, c.ContainsCapsule(Cap(Pt(0, 3), Pt(10, 3), 3)))
	})
	t.Run("larger or longer", func(t *testing.T) {
		assert.False(t, c.ContainsCapsule(Cap(Pt(0, 0), Pt(10, 0), 6)))
		assert.False(t, c.ContainsCapsule(Cap(Pt(0, 0), Pt(16, 0), 1)))
	})
}

func TestCapsule_ContainsSector(t *testing.T) {
	t.Run("a sector within it", func(t *testing.T) {
		assert.True(t, Cap(Pt(1.0, 1.0), Pt(3.0, 3.0), 4.0).ContainsSector(quarterSector()))
		assert.True(t, Cap(Pt(-3.0, 0.0), Pt(3.0, 0.0), 5.0).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
	})
	t.Run("a sector with a corner outside, or its arc alone", func(t *testing.T) {
		assert.False(t, Cap(Pt(1.0, 1.0), Pt(3.0, 3.0), 3.0).ContainsSector(quarterSector()))
		assert.False(t, Cap(Pt(-3.0, 0.0), Pt(3.0, 0.0), 4.9).ContainsSector(Sect(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2)))
		assert.False(t, Cap(Pt(-9.0, 3.0), Pt(-4.0, 3.0), 4.9).ContainsSector(Sect(Pt(-4.0, 3.0), 5.0, -Pi/4, Pi/2)))
	})
}

func TestCapsule_ContainsSegment(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("inside, its own segment included", func(t *testing.T) {
		assert.True(t, c.ContainsSegment(Seg(Pt(1, 1), Pt(9, -1))))
		assert.True(t, c.ContainsSegment(c.Segment))
	})
	t.Run("from end to end", func(t *testing.T) {
		assert.True(t, c.ContainsSegment(Seg(Pt(-2, 0), Pt(12, 0))))
	})
	t.Run("an endpoint outside", func(t *testing.T) {
		assert.False(t, c.ContainsSegment(Seg(Pt(5, 0), Pt(5, 3))))
	})
}

func TestCapsule_ContainsPolygon(t *testing.T) {
	c := Cap(Pt(0, 1), Pt(2, 1), 2)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, c.ContainsPolygon(Pol(squareVertices())))
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, Cap(Pt(0, 1), Pt(1, 1), 1).ContainsPolygon(Pol(squareVertices())))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, c.ContainsPolygon(Pol[int](nil)))
	})
}

func TestCapsule_ContainsRectangle(t *testing.T) {
	c := Cap(Pt(-4.0, 0.0), Pt(4.0, 0.0), 2.0)

	t.Run("corners on the boundary count", func(t *testing.T) {
		assert.True(t, c.ContainsRectangle(Rect(Pt(0.0, 0.0), Sz(8.0, 4.0), 0)))
		assert.False(t, c.ContainsRectangle(Rect(Pt(0.0, 0.0), Sz(8.0, 5.0), 0)))
	})
	t.Run("rotated is tested on its turned corners", func(t *testing.T) {
		assert.True(t, c.ContainsRectangle(Rect(Pt(0.0, 0.0), Sz(2.0, 1.0), 0).Rotate(Pi/4)))
		assert.False(t, c.ContainsRectangle(Rect(Pt(0.0, 0.0), Sz(8.0, 4.0), 0).Rotate(Pi/4)))
	})
}

func TestCapsule_ContainsRegularPolygon(t *testing.T) {
	c := Cap(Pt(-4, 0), Pt(4, 0), 2)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, c.ContainsRegularPolygon(RegPol(Pt(3, 0), Sz(2, 2), 6, 0, 0)))
	})
	t.Run("a vertex outside", func(t *testing.T) {
		assert.False(t, c.ContainsRegularPolygon(RegPol(Pt(5, 0), Sz(2, 2), 4, 0, 0)))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, c.ContainsRegularPolygon(RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0)))
	})
}

func TestCapsule_ContainsBox(t *testing.T) {
	c := Cap(Pt(-4, 0), Pt(4, 0), 2)

	t.Run("corners on the boundary count", func(t *testing.T) {
		assert.True(t, c.ContainsBox(Bx(Pt(-4, -2), Pt(4, 2))))
	})
	t.Run("a corner outside", func(t *testing.T) {
		assert.False(t, c.ContainsBox(Bx(Pt(-5, -2), Pt(4, 2))))
	})
}

func TestCapsule_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsCapsule", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, circle := range circleFixtures {
				assert.Equal(t, c.IntersectsCircle(circle), circle.IntersectsCapsule(c), fmt.Sprintf("%s → %s: ", c, circle))
			}
		}
	})
}

func TestCapsule_IntersectsCapsule(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("touching side to side or end to end counts", func(t *testing.T) {
		assert.True(t, c.IntersectsCapsule(Cap(Pt(0, 4), Pt(10, 4), 2)))
		assert.True(t, c.IntersectsCapsule(Cap(Pt(14, 0), Pt(20, 0), 2)))
	})
	t.Run("crossing", func(t *testing.T) {
		assert.True(t, c.IntersectsCapsule(Cap(Pt(5, -9), Pt(5, 9), 1)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, c.IntersectsCapsule(Cap(Pt(0, 5), Pt(10, 5), 2)))
		assert.False(t, c.IntersectsCapsule(Cap(Pt(15, 0), Pt(20, 0), 2)))
	})
	t.Run("float within the tolerance counts, beyond it not", func(t *testing.T) {
		assert.True(t, c.Float().IntersectsCapsule(Cap(Pt(0.0, 4.0+Delta/2), Pt(10.0, 4.0+Delta/2), 2.0)))
		assert.False(t, c.Float().IntersectsCapsule(Cap(Pt(0.0, 4.0+2*Delta), Pt(10.0, 4.0+2*Delta), 2.0)))
	})
	t.Run("either side gives the same answer, the one their segments' distance gives", func(t *testing.T) {
		for _, a := range capsuleFixtures {
			for _, b := range capsuleFixtures {
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, a.IntersectsCapsule(b), b.IntersectsCapsule(a), message)
				assert.Equal(t, a.IntersectsCapsule(b), LessOrEqual(a.Segment.DistanceToSegment(b.Segment), a.Radius+b.Radius), message)
			}
		}
	})
}

func TestCapsule_IntersectsSector(t *testing.T) {
	t.Run("mirrors Sector.IntersectsCapsule", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, s := range sectorFixtures() {
				assert.Equal(t, c.IntersectsSector(s), s.IntersectsCapsule(c), fmt.Sprintf("%s → %s: ", c, s))
			}
		}
	})
}

func TestCapsule_IntersectsSegment(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("touching a side or a round end counts", func(t *testing.T) {
		assert.True(t, c.IntersectsSegment(Seg(Pt(0, 2), Pt(10, 2))))
		assert.True(t, c.IntersectsSegment(Seg(Pt(12, -5), Pt(12, 5))))
	})
	t.Run("crossing or inside", func(t *testing.T) {
		assert.True(t, c.IntersectsSegment(Seg(Pt(5, -9), Pt(5, 9))))
		assert.True(t, c.IntersectsSegment(Seg(Pt(4, 1), Pt(6, -1))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, c.IntersectsSegment(Seg(Pt(0, 3), Pt(10, 3))))
		assert.False(t, c.IntersectsSegment(Seg(Pt(12, 2), Pt(14, 4))))
	})
	t.Run("int is exact near the end of the supported range", func(t *testing.T) {
		long := Cap(Pt(0, 0), Pt(1<<25, 0), 3)

		assert.True(t, long.IntersectsSegment(Seg(Pt(0, 3), Pt(1<<25, 3))))
		assert.False(t, long.IntersectsSegment(Seg(Pt(0, 4), Pt(1<<25, 4))))
	})
	t.Run("a segment of no length is the point", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, c.IntersectsSegment(Seg(p, p)), c.Contains(p), fmt.Sprintf("%s → %s: ", c, p))
			}
		}
	})
}

func TestCapsule_IntersectsRay(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("toward it, or touching in passing", func(t *testing.T) {
		assert.True(t, c.IntersectsRay(Ry(Pt(5, 10), Vec(0, -1))))
		assert.True(t, c.IntersectsRay(Ry(Pt(-9, 2), Vec(1, 0))))
	})
	t.Run("away from it or beside it", func(t *testing.T) {
		assert.False(t, c.IntersectsRay(Ry(Pt(5, 10), Vec(0, 1))))
		assert.False(t, c.IntersectsRay(Ry(Pt(-9, 3), Vec(1, 0))))
	})
}

func TestCapsule_IntersectsLine(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	t.Run("through it either way, or touching", func(t *testing.T) {
		assert.True(t, c.IntersectsLine(Ln(Pt(5, 10), Vec(0, 1))))
		assert.True(t, c.IntersectsLine(Ln(Pt(-9, 2), Vec(-1, 0))))
	})
	t.Run("beside it", func(t *testing.T) {
		assert.False(t, c.IntersectsLine(Ln(Pt(-9, 3), Vec(1, 0))))
	})
}

func TestCapsule_IntersectsPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("touching an edge or a vertex counts", func(t *testing.T) {
		assert.True(t, Cap(Pt(4, 1), Pt(8, 1), 2).IntersectsPolygon(square))
		assert.True(t, Cap(Pt(5, 6), Pt(9, 6), 5).IntersectsPolygon(square))
	})
	t.Run("one within the other", func(t *testing.T) {
		assert.True(t, Cap(Pt(1, 1), Pt(1, 1), 0).IntersectsPolygon(square))
		assert.True(t, Cap(Pt(-9, 1), Pt(9, 1), 9).IntersectsPolygon(square))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Cap(Pt(5, 1), Pt(8, 1), 2).IntersectsPolygon(square))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Cap(Pt(1, 1), Pt(1, 1), 9).IntersectsPolygon(Pol[int](nil)))
	})
}

func TestCapsule_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("touching a side counts", func(t *testing.T) {
		assert.True(t, Cap(Pt(0.0, 3.0), Pt(5.0, 3.0), 2.0).IntersectsRectangle(rectangle))
		assert.False(t, Cap(Pt(0.0, 4.0), Pt(5.0, 4.0), 2.0).IntersectsRectangle(rectangle))
	})
	t.Run("one within the other", func(t *testing.T) {
		assert.True(t, Cap(Pt(-1.0, 0.0), Pt(1.0, 0.0), 0.5).IntersectsRectangle(rectangle))
		assert.True(t, Cap(Pt(-9.0, 0.0), Pt(9.0, 0.0), 9.0).IntersectsRectangle(rectangle))
	})
	t.Run("rotated is tested on its turned edges", func(t *testing.T) {
		assert.False(t, Cap(Pt(0.0, 4.0), Pt(5.0, 4.0), 2.0).IntersectsRectangle(rectangle))
		assert.True(t, Cap(Pt(0.0, 4.0), Pt(5.0, 4.0), 2.0).IntersectsRectangle(rectangle.Rotate(Pi/2)))
	})
}

func TestCapsule_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("touching a vertex counts", func(t *testing.T) {
		assert.True(t, Cap(Pt(4, 0), Pt(8, 0), 2).IntersectsRegularPolygon(diamond))
		assert.False(t, Cap(Pt(5, 0), Pt(8, 0), 2).IntersectsRegularPolygon(diamond))
	})
	t.Run("one within the other", func(t *testing.T) {
		assert.True(t, Cap(Pt(0, 0), Pt(0, 0), 0).IntersectsRegularPolygon(diamond))
		assert.True(t, Cap(Pt(-9, 0), Pt(9, 0), 9).IntersectsRegularPolygon(diamond))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Cap(Pt(0, 0), Pt(0, 0), 9).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
}

func TestCapsule_IntersectsBox(t *testing.T) {
	box := Bx(Pt(0, 0), Pt(4, 2))

	t.Run("touching a side counts", func(t *testing.T) {
		assert.True(t, Cap(Pt(6, 1), Pt(9, 1), 2).IntersectsBox(box))
		assert.False(t, Cap(Pt(7, 1), Pt(9, 1), 2).IntersectsBox(box))
	})
	t.Run("one within the other", func(t *testing.T) {
		assert.True(t, Cap(Pt(1, 1), Pt(3, 1), 0).IntersectsBox(box))
		assert.True(t, Cap(Pt(-9, 1), Pt(9, 1), 9).IntersectsBox(box))
	})
	t.Run("a round end clears a corner its bounds overlap", func(t *testing.T) {
		assert.False(t, Cap(Pt(6, 4), Pt(9, 4), 2).IntersectsBox(box))
		assert.True(t, Cap(Pt(6, 4), Pt(9, 4), 2).Bounds().IntersectsBox(box))
	})
}

func TestCapsule_IntersectionSegment(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("passing through gives both crossings from Start", func(t *testing.T) {
		geomtest.AssertVertices(t, capsule.IntersectionSegment(Seg(Pt(5.0, 10.0), Pt(5.0, -10.0))).AppendTo(nil), []Point[float64]{Pt(5.0, 2.0), Pt(5.0, -2.0)})
		geomtest.AssertVertices(t, capsule.IntersectionSegment(Seg(Pt(-5.0, 0.0), Pt(15.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(-2.0, 0.0), Pt(12.0, 0.0)})
	})
	t.Run("an end inside gives one, and a tangent its point", func(t *testing.T) {
		geomtest.AssertVertices(t, capsule.IntersectionSegment(Seg(Pt(5.0, 0.0), Pt(5.0, 10.0))).AppendTo(nil), []Point[float64]{Pt(5.0, 2.0)})
		geomtest.AssertVertices(t, capsule.IntersectionSegment(Seg(Pt(12.0, -5.0), Pt(12.0, 5.0))).AppendTo(nil), []Point[float64]{Pt(12.0, 0.0)})
	})
	t.Run("inside or apart gives none", func(t *testing.T) {
		assert.Nil(t, capsule.IntersectionSegment(Seg(Pt(1.0, 1.0), Pt(9.0, -1.0))).AppendTo(nil))
		assert.Nil(t, capsule.IntersectionSegment(Seg(Pt(-5.0, 3.0), Pt(15.0, 3.0))).AppendTo(nil))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Cap(Pt(0, 0), Pt(10, 0), 2).IntersectionSegment(Seg(Pt(5, 10), Pt(5, -10))).AppendTo(nil), []Point[int]{Pt(5, 2), Pt(5, -2)})
	})
	t.Run("a capsule of no length is crossed as its circle", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, Cap(c.Center, c.Center, c.Radius).IntersectionSegment(s).AppendTo(nil), c.IntersectionSegment(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", c, s))
			}
		}
	})
	t.Run("every point lies on the segment and the boundary, where the two intersect", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					points := c.IntersectionSegment(s).AppendTo(nil)
					message := fmt.Sprintf("%s → %s: ", c, s)

					for _, point := range points {
						assert.True(t, s.Contains(point) && c.Contains(point), message+point.String())
						geomtest.AssertNumber(t, c.Segment.DistanceTo(point), c.Radius, message+point.String())
					}
					assertOrdered(t, s, points, message)
					assert.True(t, len(points) == 0 || c.IntersectsSegment(s), message)
					assert.True(t, len(points) > 0 || !c.IntersectsSegment(s) || c.Contains(start) || c.Contains(end), message)
				}
			}
		}
	})
	t.Run("a graze within the tolerance is a crossing, far from the origin in float32 too", func(t *testing.T) {
		assertGrazes(t, Vec(0.0, 0.0))
		assertGrazes(t, Vec[float32](0, 0))
		for _, offset := range farOffsets {
			assertGrazes(t, offset)
		}
	})
	t.Run("int crosses, clips and hits at one point", func(t *testing.T) {
		assertGrazes(t, Vec(0, 0))
		assertGrazes(t, Vec(1<<25, -1<<25))

		skewed, cast := Cap(Pt(-5, 4), Pt(-1, 1), 2), Seg(Pt(0, 5), Pt(-2, 2))
		hit, ok := cast.HitCapsule(skewed)

		assert.True(t, ok)
		geomtest.AssertVertices(t, skewed.IntersectionSegment(cast).AppendTo(nil), []Point[int]{hit.Point})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(5.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			sinkBool = capsule.IntersectionSegment(through).Count > 0
		}, 0)
	})
}

// assertGrazes asserts, over the capsules and the segments of a small lattice moved by the
// offset, where tangents and ends on the boundary are common, that a segment crosses a capsule
// where the two intersect, unless it lies within, has a part inside exactly there, and hits it
// from outside at the first of the crossings.
func assertGrazes[T Number](t *testing.T, offset Vector[T]) {
	t.Helper()

	var lattice []Point[T]
	for _, x := range []T{-2, 0, 1, 3} {
		for _, y := range []T{-2, 0, 1, 3} {
			lattice = append(lattice, Pt(x, y).Add(offset))
		}
	}

	for _, radius := range []T{0, 1, 2} {
		for i, a := range lattice {
			for _, b := range lattice[i:] {
				c := Cap(a, b, radius)

				for _, start := range lattice {
					for _, end := range lattice {
						s := Seg(start, end)
						points := c.IntersectionSegment(s)
						intersects := c.IntersectsSegment(s)
						message := fmt.Sprintf("%s → %s: ", c, s)

						assert.True(t, points.Count == 0 || intersects, message)
						assert.True(t, points.Count > 0 || !intersects || c.Contains(start) && c.Contains(end), message)

						_, clipped := s.ClipCapsule(c)
						assert.Equal(t, clipped, intersects, message)

						hit, ok := s.HitCapsule(c)
						assert.Equal(t, ok, !c.Contains(start) && intersects, message)
						assert.True(t, !ok || hit.Point == points.Items[0], message)
					}
				}
			}
		}
	}
}

func TestCapsule_IntersectionRay(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	geomtest.AssertVertices(t, capsule.IntersectionRay(Ry(Pt(5.0, 10.0), Vec(0.0, -1.0))).AppendTo(nil), []Point[float64]{Pt(5.0, 2.0), Pt(5.0, -2.0)})
	geomtest.AssertVertices(t, capsule.IntersectionRay(Ry(Pt(5.0, 0.0), Vec(0.0, -1.0))).AppendTo(nil), []Point[float64]{Pt(5.0, -2.0)})
	assert.Nil(t, capsule.IntersectionRay(Ry(Pt(5.0, 10.0), Vec(0.0, 1.0))).AppendTo(nil))
}

func TestCapsule_IntersectionLine(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	geomtest.AssertVertices(t, capsule.IntersectionLine(Ln(Pt(5.0, 0.0), Vec(0.0, -1.0))).AppendTo(nil), []Point[float64]{Pt(5.0, 2.0), Pt(5.0, -2.0)})
	assert.Nil(t, capsule.IntersectionLine(Ln(Pt(5.0, 10.0), Vec(1.0, 0.0))).AppendTo(nil))
}

func TestCapsule_Equal(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assert.True(t, Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.5).Equal(Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.5)))
	})
	t.Run("a different radius", func(t *testing.T) {
		assert.False(t, Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.5).Equal(Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.6)))
	})
	t.Run("the reversed segment is another value", func(t *testing.T) {
		assert.False(t, Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.5).Equal(Cap(Pt(3.0, 5.0), Pt(1.0, 2.0), 0.5)))
	})
}

func TestCapsule_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.True(t, Capsule[int]{}.IsZero())
	})
	t.Run("a radius makes it non-zero", func(t *testing.T) {
		assert.False(t, Cap(Pt(0, 0), Pt(0, 0), 1).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Cap(Pt(0.0, 0.0000001), Pt(0.0, 0.0), 0.0).IsZero())
	})
}

func TestCapsule_Polygon(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("a tolerance of the radius leaves the corners of the sides", func(t *testing.T) {
		geomtest.AssertPolygon(t, capsule.Polygon(2), Pol([]Point[float64]{Pt(10.0, -2.0), Pt(10.0, 2.0), Pt(0.0, 2.0), Pt(0.0, -2.0)}))
		geomtest.AssertPolygon(t, capsule.Polygon(-2), capsule.Polygon(2))
	})
	t.Run("a finer tolerance rounds the ends", func(t *testing.T) {
		geomtest.AssertPolygon(t, capsule.Polygon(0.9), Pol([]Point[float64]{Pt(10.0, -2.0), Pt(12.0, 0.0), Pt(10.0, 2.0), Pt(0.0, 2.0), Pt(-2.0, 0.0), Pt(0.0, -2.0)}))
		assert.Length(t, capsule.Polygon(0).Points, 2*1025)
	})
	t.Run("int rounds the vertices", func(t *testing.T) {
		geomtest.AssertPolygon(t, Cap(Pt(0, 0), Pt(10, 0), 2).Polygon(0.9), Pol([]Point[int]{Pt(10, -2), Pt(12, 0), Pt(10, 2), Pt(0, 2), Pt(-2, 0), Pt(0, -2)}))
	})
	t.Run("no radius gives the two ends, and no length the circle about the point", func(t *testing.T) {
		geomtest.AssertPolygon(t, Cap(Pt(0, 0), Pt(10, 0), 0).Polygon(1), Pol([]Point[int]{Pt(10, 0), Pt(0, 0)}))
		geomtest.AssertPolygon(t, Cap(Pt(1, 2), Pt(1, 2), 0).Polygon(1), Pol([]Point[int]{Pt(1, 2)}))
		geomtest.AssertPolygon(t, Cap(Pt(1, 2), Pt(1, 2), 2).Polygon(0.9), Pol([]Point[int]{Pt(1, 0), Pt(3, 2), Pt(1, 4), Pt(-1, 2)}))
		geomtest.AssertPolygon(t, Cap(Pt(1.0, 2.0), Pt(1.0, 2.0), 2.0).Polygon(0.9), Pol([]Point[float64]{Pt(1.0, 0.0), Pt(3.0, 2.0), Pt(1.0, 4.0), Pt(-1.0, 2.0)}))
		geomtest.AssertPolygon(t, Cap(Pt(1.0, 2.0), Pt(1.0, 2.0), 2.0).Polygon(9), Pol([]Point[float64]{Pt(1.0, 0.0), Pt(3.0, 2.0), Pt(1.0, 4.0), Pt(-1.0, 2.0)}))
	})
	t.Run("the vertices lie on the boundary, the edges within the tolerance of it, wound as a rectangle", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, tolerance := range []float64{0.01, 0.25, 5} {
				polygon := c.Polygon(tolerance)
				message := fmt.Sprintf("%s ±%v: ", c, tolerance)

				for vertex := range polygon.Vertices() {
					geomtest.AssertNumber(t, c.Segment.DistanceTo(vertex), c.Radius, message+vertex.String())
				}
				for edge := range polygon.Edges() {
					assert.True(t, c.Contains(edge.Midpoint()), message+edge.String())
					assert.True(t, LessOrEqual(c.Radius-c.Segment.DistanceTo(edge.Midpoint()), tolerance), message+edge.String())
				}

				assert.True(t, polygon.Area() == 0 || polygon.Winding() == WindingClockwise && polygon.IsConvex(), message)
				assert.True(t, polygon.Area() <= c.Area()+Delta, message)
			}
		}
	})
	t.Run("a fine outline has the area and the perimeter", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			polygon := c.Polygon(0)

			assert.EqualDelta(t, polygon.Area(), c.Area(), 1e-4, c.String())
			assert.EqualDelta(t, polygon.Perimeter(), c.Perimeter(), 1e-4, c.String())
		}
	})
}

func TestCapsule_AppendPolygon(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, capsule.AppendPolygon([]Point[float64]{Pt(9.0, 9.0)}, 2), []Point[float64]{Pt(9.0, 9.0), Pt(10.0, -2.0), Pt(10.0, 2.0), Pt(0.0, 2.0), Pt(0.0, -2.0)})
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		whole := Cap(Pt(0, 0), Pt(100, 0), 20)
		buffer := make([]Point[int], 0, 64)

		assertAllocations(t, func() {
			sinkPoints = whole.Polygon(tenth).Points
		}, 1)
		assertAllocations(t, func() {
			sinkPoints = whole.AppendPolygon(buffer[:0], tenth)
		}, 0)
	})
}

func TestCapsule_Cast(t *testing.T) {
	c := Cap(Pt(1.5, -2.5), Pt(3.5, 1.5), 1.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertCapsule(t, c.Cast[int](), c.Int())
		geomtest.AssertCapsule(t, c.Cast[float64](), c.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertCapsule(t, c.Cast[int8](), Cap(Pt[int8](2, -3), Pt[int8](4, 2), 2))
	})
}

func TestCapsule_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Int(), Cap(Pt(0, 0), Pt(10, 0), 2))
	})
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0.6, -0.25), Pt(1.2, 3.4), 0.5).Int(), Cap(Pt(1, 0), Pt(1, 3), 1))
	})
}

func TestCapsule_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0, 0), Pt(10, 0), 2).Float(), Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertCapsule(t, Cap(Pt(0.6, -0.25), Pt(1.2, 3.4), 0.5).Float(), Cap(Pt(0.6, -0.25), Pt(1.2, 3.4), 0.5))
	})
}

func TestCapsule_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Cap(Pt(0, 0), Pt(10, 0), 2).String(), "Cap((0,0);(10,0);2)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Cap(Pt(100, -34.0000115), Pt(0.2, 0.1), 0.5).String(), "Cap((100.00,-34.00);(0.20,0.10);0.50)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Cap(Pt(0, 0), Pt(10, 0), 2).String)
	})
}

func TestCapsule_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Cap(Pt(0, 1), Pt(10, 0), 2), `{"s":{"x":0,"y":1},"e":{"x":10,"y":0},"r":2}`)

		var c Capsule[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":0,"y":1},"e":{"x":10,"y":0},"r":2}`), &c))
		geomtest.AssertCapsule(t, c, Cap(Pt(0, 1), Pt(10, 0), 2))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, capsule := range capsuleFixtures {
			data, err := json.Marshal(capsule)
			assert.NoError(t, err)

			var decoded Capsule[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, capsule)
		}
	})
}

func TestCapsule_Properties(t *testing.T) {
	t.Run("area and perimeter follow the radius and the length", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			length := c.Segment.Length()

			geomtest.AssertNumber(t, c.Area(), Pi*c.Radius*c.Radius+2*c.Radius*length, fmt.Sprintf("%s: ", c))
			geomtest.AssertNumber(t, c.Perimeter(), 2*Pi*c.Radius+2*length, fmt.Sprintf("%s: ", c))
		}
	})
	t.Run("a segment of no length answers as the circle of the same radius", func(t *testing.T) {
		for _, circle := range circleFixtures {
			c := Cap(circle.Center, circle.Center, circle.Radius)

			geomtest.AssertNumber(t, c.Area(), circle.Area(), c.String())
			geomtest.AssertNumber(t, c.Perimeter(), circle.Perimeter(), c.String())
			geomtest.AssertNumber(t, c.Inertia(), circle.Inertia(), c.String())
			geomtest.AssertBox(t, c.Bounds(), circle.Bounds(), c.String())

			for _, point := range pointFixtures {
				assert.Equal(t, c.Contains(point), circle.Contains(point), fmt.Sprintf("%s → %s: ", c, point))
				geomtest.AssertNumber(t, c.DistanceTo(point), circle.DistanceTo(point), fmt.Sprintf("%s → %s: ", c, point))
				geomtest.AssertPoint(t, c.Nearest(point), circle.Nearest(point), fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("a radius of zero answers as the segment", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			c := Cap(segment.Start, segment.End, 0.0)

			geomtest.AssertBox(t, c.Bounds(), segment.Bounds(), c.String())

			for _, point := range pointFixtures {
				assert.Equal(t, c.Contains(point), segment.Contains(point), fmt.Sprintf("%s → %s: ", c, point))
				geomtest.AssertNumber(t, c.DistanceTo(point), segment.DistanceTo(point), fmt.Sprintf("%s → %s: ", c, point))
				geomtest.AssertPoint(t, c.Nearest(point), segment.Nearest(point), fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("the boundary is contained and its bounds hold it", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, direction := range vectorFixtures {
				point := c.Support(direction)

				assert.True(t, c.Contains(point), fmt.Sprintf("%s → %s: ", c, point))
				assert.True(t, c.Bounds().Contains(point), fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("the turn does not change the distance to a point turned with it", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, point := range pointFixtures {
				turned := point.RotateAround(c.Centroid(), Pi/3)

				geomtest.AssertNumber(t, c.Rotate(Pi/3).DistanceTo(turned), c.DistanceTo(point), fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, c.Scale(factor).Unscale(factor).Equal(c), fmt.Sprintf("%s ×%v: ", c, factor))
			}
		}
	})
	t.Run("grow and shrink are inverse", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			assert.True(t, c.Grow(1).Shrink(1).Equal(c), fmt.Sprintf("%s: ", c))
		}
	})
	t.Run("a grown capsule contains the points within the amount of it", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, point := range pointFixtures {
				assert.Equal(t, c.Grow(1).Contains(point), LessOrEqual(c.DistanceTo(point), 1), fmt.Sprintf("%s → %s: ", c, point))
			}
		}
	})
	t.Run("an enclosed shape intersects it and has its vertices within", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, circle := range circleFixtures {
				assert.True(t, !c.ContainsCircle(circle) || Intersects[float64](c, circle), fmt.Sprintf("%s → %s: ", c, circle))
				assert.True(t, !c.ContainsCircle(circle) || c.Contains(circle.Center), fmt.Sprintf("%s → %s: ", c, circle))
			}
			for _, s := range segmentFixtures {
				assert.True(t, !c.ContainsSegment(s) || Intersects[float64](c, s), fmt.Sprintf("%s → %s: ", c, s))
			}
			for _, p := range polygonFixtures() {
				assert.True(t, !c.ContainsPolygon(p) || Intersects[float64](c, p), fmt.Sprintf("%s → %s: ", c, p))
			}
			for _, r := range rectFixtures {
				assert.True(t, !c.ContainsRectangle(r) || Intersects[float64](c, r), fmt.Sprintf("%s → %s: ", c, r))
			}
			for _, rp := range regularPolygonFixtures {
				assert.True(t, !c.ContainsRegularPolygon(rp) || Intersects[float64](c, rp), fmt.Sprintf("%s → %s: ", c, rp))
			}
			for _, b := range boxFixtures {
				assert.True(t, !c.ContainsBox(b) || Intersects[float64](c, b), fmt.Sprintf("%s → %s: ", c, b))
			}
		}
	})
	t.Run("a capsule encloses its own segment and the circles at its ends", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			assert.True(t, c.ContainsSegment(c.Segment), c.String())
			assert.True(t, c.ContainsCircle(Circ(c.Segment.Start, c.Radius)), c.String())
			assert.True(t, c.ContainsCircle(Circ(c.Segment.End, c.Radius)), c.String())
		}
	})
}

func TestCapsule_Immutable(t *testing.T) {
	c := Cap(Pt(0, 0), Pt(10, 0), 2)

	c.Translate(Vec(3, -2))
	c.MoveTo(Pt(4, 3))
	c.Rotate(1)
	c.Scale(2)
	c.Grow(1)
	c.Shrink(2)

	geomtest.AssertCapsule(t, c, Cap(Pt(0, 0), Pt(10, 0), 2))
}

// capsuleFixtures span the point, the circle, the segment, and upright, lying and slanted
// capsules off the origin.
var capsuleFixtures = []Capsule[float64]{
	Cap(Pt(0.0, 0.0), Pt(0.0, 0.0), 0.0),
	Cap(Pt(0.0, 0.0), Pt(0.0, 0.0), 1.0),
	Cap(Pt(1.0, 2.0), Pt(3.0, 5.0), 0.0),
	Cap(Pt(1.0, 2.0), Pt(5.0, 2.0), 1.5),
	Cap(Pt(-3.5, 0.25), Pt(-3.5, 7.0), 2.0),
	Cap(Pt(12.5, -0.1), Pt(-0.5, 12.75), 0.5),
}

func ExampleCap() {
	fmt.Println(Cap(Pt(0, 0), Pt(10, 0), 2))
	// Output: Cap((0,0);(10,0);2)
}

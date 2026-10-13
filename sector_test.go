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

func TestSector_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, 1.5), Sector[int]{Center: Pt(1, 2), Radius: 5, Angle: 0.5, Sweep: 1.5})
	})
	t.Run("a negative radius is taken absolute and a negative sweep kept", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1.0, 2.0), -5.0, 0.5, -1.5), Sector[float64]{Center: Pt(1.0, 2.0), Radius: 5, Angle: 0.5, Sweep: -1.5})
	})
}

func TestSector_Centroid(t *testing.T) {
	t.Run("on the bisector of a quarter disc", func(t *testing.T) {
		geomtest.AssertPoint(t, quarterSector().Centroid(), Pt(20/(3*Pi), 20/(3*Pi)))
	})
	t.Run("the center of the whole disc", func(t *testing.T) {
		geomtest.AssertPoint(t, Sect(Pt(1.0, 2.0), 5.0, 1, 2*Pi).Centroid(), Pt(1.0, 2.0))
	})
	t.Run("two thirds along a sweep of zero", func(t *testing.T) {
		geomtest.AssertPoint(t, Sect(Pt(0.0, 0.0), 6.0, 0, 0).Centroid(), Pt(4.0, 0.0))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, Sect(Pt(0, 0), 5, 0, Pi/2).Centroid(), Pt(2, 2))
	})
}

func TestSector_Area(t *testing.T) {
	t.Run("half the squared radius times the sweep, either way", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, Pi/2).Area(), 6.25*Pi)
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, -Pi/2).Area(), 6.25*Pi)
	})
	t.Run("no more than the disc", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, 9).Area(), Circ(Pt(0, 0), 5).Area())
	})
}

func TestSector_Perimeter(t *testing.T) {
	t.Run("the arc and the two radii", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, Pi/2).Perimeter(), 2.5*Pi+10)
	})
	t.Run("the whole disc has no radius on its boundary", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, -9).Perimeter(), 10*Pi)
	})
}

func TestSector_Inertia(t *testing.T) {
	t.Run("the whole disc", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(1, 2), 3, 1, 2*Pi).Inertia(), Circ(Pt(1, 2), 3).Inertia())
	})
	t.Run("half a disc about its own centroid", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0.0, 0.0), 3.0, 0, Pi).Inertia(), 81*Pi/4-4.5*Pi*(4/Pi)*(4/Pi))
	})
	t.Run("the turn does not change it", func(t *testing.T) {
		geomtest.AssertNumber(t, quarterSector().Rotate(1).Inertia(), quarterSector().Inertia())
	})
}

func TestSector_Bounds(t *testing.T) {
	t.Run("the center and the ends", func(t *testing.T) {
		geomtest.AssertBox(t, quarterSector().Bounds(), Bx(Pt(0.0, 0.0), Pt(5.0, 5.0)))
	})
	t.Run("out to the circle along an axis the sweep passes, back to the center", func(t *testing.T) {
		geomtest.AssertBox(t, Sect(Pt(0.0, 0.0), 5.0, 3*Pi/4, Pi/2).Bounds(), Bx(Pt(-5.0, -5/math.Sqrt2), Pt(0.0, 5/math.Sqrt2)))
	})
	t.Run("the whole disc", func(t *testing.T) {
		geomtest.AssertBox(t, Sect(Pt(1.0, 2.0), 5.0, 1, 2*Pi).Bounds(), Circ(Pt(1.0, 2.0), 5.0).Bounds())
	})
}

func TestSector_BoundingCircle(t *testing.T) {
	geomtest.AssertCircle(t, Sect(Pt(1, 2), 5, 0.5, 1.5).BoundingCircle(), Circ(Pt(1, 2), 5))
}

func TestSector_Translate(t *testing.T) {
	geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, 1.5).Translate(Vec(3, -2)), Sect(Pt(4, 0), 5, 0.5, 1.5))
}

func TestSector_MoveTo(t *testing.T) {
	geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, 1.5).MoveTo(Pt(3, -2)), Sect(Pt(3, -2), 5, 0.5, 1.5))
}

func TestSector_Rotate(t *testing.T) {
	t.Run("the angle moves and the sweep stays", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, 1.5).Rotate(1), Sect(Pt(1, 2), 5, 1.5, 1.5))
	})
	t.Run("normalized to a single turn", func(t *testing.T) {
		assert.Equal(t, Sect(Pt(1, 2), 5, Pi, 1.5).Rotate(3*Pi/2).Angle, Pi/2)
	})
}

func TestSector_Scale(t *testing.T) {
	t.Run("the radius", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Scale(2.5), Sect(Pt(1, 2), 10, 0.5, 1.5))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Scale(-2), Sect(Pt(1, 2), 8, 0.5, 1.5))
	})
}

func TestSector_Unscale(t *testing.T) {
	t.Run("the radius", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 10, 0.5, 1.5).Unscale(-2.5), Sect(Pt(1, 2), 4, 0.5, 1.5))
	})
	t.Run("a zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Sect(Pt(1, 2), 10, 0.5, 1.5).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestSector_Resize(t *testing.T) {
	geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Resize(-9), Sect(Pt(1, 2), 9, 0.5, 1.5))
}

func TestSector_Grow(t *testing.T) {
	t.Run("the radius grows and the center stays", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Grow(3), Sect(Pt(1, 2), 7, 0.5, 1.5))
	})
	t.Run("a negative amount clamps at zero", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Grow(-9), Sect(Pt(1, 2), 0, 0.5, 1.5))
	})
}

func TestSector_Shrink(t *testing.T) {
	t.Run("the radius shrinks and the center stays", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Shrink(3), Sect(Pt(1, 2), 1, 0.5, 1.5))
	})
	t.Run("clamps at zero", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 4, 0.5, 1.5).Shrink(9), Sect(Pt(1, 2), 0, 0.5, 1.5))
	})
}

func TestSector_Lerp(t *testing.T) {
	t.Run("the center, the radius, the angle and the sweep move together", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(0, 0), 4, 0, 1).Lerp(Sect(Pt(10, 0), 8, 2, 3), 0.5), Sect(Pt(5, 0), 6, 1, 2))
	})
	t.Run("ends on the two sectors", func(t *testing.T) {
		for _, a := range sectorFixtures() {
			for _, b := range sectorFixtures() {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
}

func TestSector_Canonical(t *testing.T) {
	t.Run("a well-formed sector is unchanged", func(t *testing.T) {
		geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, -1.5).Canonical(), Sect(Pt(1, 2), 5, 0.5, -1.5))
	})
	t.Run("the radius, the angle and a sweep past a full turn are repaired", func(t *testing.T) {
		canonical := Sector[int]{Pt(1, 2), -5, -Pi / 2, 3 * Pi}.Canonical()

		assert.Equal(t, canonical, Sector[int]{Pt(1, 2), 5, 3 * Pi / 2, 2 * Pi})
	})
}

func TestSector_DistanceTo(t *testing.T) {
	s := quarterSector()

	t.Run("zero inside", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(1.0, 1.0)), 0.0)
	})
	t.Run("to the arc", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(6.0, 8.0)), 5.0)
	})
	t.Run("to a radius", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(-3.0, 4.0)), 3.0)
		geomtest.AssertNumber(t, s.DistanceTo(Pt(2.0, -1.0)), 1.0)
	})
	t.Run("to the center behind it", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(-3.0, -4.0)), 5.0)
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Sect(Pt(0, 0), 5, 0, Pi/2).DistanceTo(Pt(-3, 4)), 3.0)
	})
}

func TestSector_DistanceSquaredTo(t *testing.T) {
	t.Run("the square of the distance", func(t *testing.T) {
		geomtest.AssertNumber(t, quarterSector().DistanceSquaredTo(Pt(-3.0, -4.0)), 25.0)
	})
	t.Run("a point within the tolerance of the boundary is on it", func(t *testing.T) {
		assert.Equal(t, quarterSector().DistanceSquaredTo(Pt(3.0, 4.0+Delta/2)), 0.0)
		assert.Equal(t, quarterSector().DistanceSquaredTo(Pt(2.0, -Delta/2)), 0.0)
	})
}

func TestSector_DistanceToSector(t *testing.T) {
	s := quarterSector()

	t.Run("between the arcs, and from an arc to a radius", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceToSector(Sect(Pt(10.0, 0.0), 2.0, Pi/2, Pi)), 3.0)
		geomtest.AssertNumber(t, s.DistanceToSector(Sect(Pt(-3.0, 1.0), 2.0, -Pi/2, Pi)), 1.0)
	})
	t.Run("zero for sectors that overlap, touch or lie one within the other", func(t *testing.T) {
		assert.Equal(t, s.DistanceToSector(Sect(Pt(4.0, 0.0), 2.0, 0, 2*Pi)), 0.0)
		assert.Equal(t, s.DistanceToSector(Sect(Pt(7.0, 0.0), 2.0, 0, 2*Pi)), 0.0)
		assert.Equal(t, s.DistanceToSector(Sect(Pt(1.0, 1.0), 1.0, 0, Pi/2)), 0.0)
	})
	t.Run("zero exactly where the sectors intersect, the same from either side, no farther than their outlines", func(t *testing.T) {
		for _, a := range sectorFixtures() {
			for _, b := range sectorFixtures() {
				message := fmt.Sprintf("%s → %s: ", a, b)
				distance := a.DistanceToSector(b)

				assert.Equal(t, distance == 0, a.IntersectsSector(b), message)
				geomtest.AssertNumber(t, distance, b.DistanceToSector(a), message)
				geomtest.AssertNumber(t, distance, Distance[float64](a, b), message)

				outlines := Distance[float64](sectorOutline(a, 257), sectorOutline(b, 257))
				assert.True(t, distance <= outlines+Delta && outlines <= distance+sagitta(a, 257)+sagitta(b, 257)+Delta, message)
			}
		}
	})
	t.Run("a NaN coordinate is at NaN distance", func(t *testing.T) {
		assert.True(t, math.IsNaN(s.DistanceToSector(Sect(Pt(math.NaN(), 0.0), 2.0, 0, 1))))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		other := Sect(Pt(10.0, tenth), 2.0, Pi/2, Pi)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = s.DistanceToSector(other) > 0
		}), 0)
	})
}

func TestSector_DistanceSquaredToSector(t *testing.T) {
	geomtest.AssertNumber(t, quarterSector().DistanceSquaredToSector(Sect(Pt(10.0, 0.0), 2.0, Pi/2, Pi)), 9.0)
	assert.Equal(t, quarterSector().DistanceSquaredToSector(Sect(Pt(7.0, 0.0), 2.0, 0, 2*Pi)), 0.0)
}

func TestSector_Nearest(t *testing.T) {
	s := quarterSector()

	t.Run("a contained point is itself", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(1.0, 1.0)), Pt(1.0, 1.0))
	})
	t.Run("the point of the arc toward it", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(6.0, 8.0)), Pt(3.0, 4.0))
	})
	t.Run("the foot on a radius", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(-3.0, 4.0)), Pt(0.0, 4.0))
		geomtest.AssertPoint(t, s.Nearest(Pt(2.0, -1.0)), Pt(2.0, 0.0))
	})
	t.Run("the center behind it", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(-3.0, -4.0)), Pt(0.0, 0.0))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, Sect(Pt(0, 0), 5, 0, Pi/2).Nearest(Pt(7, 7)), Pt(4, 4))
	})
	t.Run("the nearest point of a NaN coordinate is NaN", func(t *testing.T) {
		for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
			nearest := s.Nearest(p)

			assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y), p.String())
		}
	})
}

func TestSector_Support(t *testing.T) {
	s := quarterSector()

	t.Run("the point of the arc where the sweep passes the direction, whatever its length", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Support(Vec(3.0, 3.0)), Pt(5/math.Sqrt2, 5/math.Sqrt2))
		geomtest.AssertPoint(t, s.Support(Vec(0.6, 0.8)), Pt(3.0, 4.0))
	})
	t.Run("an end of the arc or the center past the sweep", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Support(Vec(1.0, -1.0)), Pt(5.0, 0.0))
		geomtest.AssertPoint(t, s.Support(Vec(-1.0, 0.5)), Pt(0.0, 5.0))
		geomtest.AssertPoint(t, s.Support(Vec(-1.0, -1.0)), Pt(0.0, 0.0))
	})
	t.Run("no direction gives the center", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Support(Vec(0.0, 0.0)), Pt(0.0, 0.0))
	})
	t.Run("past half a turn it is that of the hull", func(t *testing.T) {
		geomtest.AssertPoint(t, notchedSector().Support(Vec(1.0, 0.0)), Pt(5/math.Sqrt2, 5/math.Sqrt2))
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertPoint(t, Sect(Pt(1, 2), 5, 0, Pi/2).Support(Vec(0.6, 0.8)), Pt(4.0, 6.0))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		direction := Vec(tenth, 1.0)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = s.Support(direction).X > 0
		}), 0)
	})
}

func TestSector_Contains(t *testing.T) {
	s := quarterSector()

	t.Run("inside, the center, the radii and the arc included", func(t *testing.T) {
		for _, point := range []Point[float64]{Pt(1.0, 1.0), Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(0.0, 5.0), Pt(3.0, 4.0)} {
			assert.True(t, s.Contains(point), point.String())
		}
	})
	t.Run("past the arc or outside the sweep", func(t *testing.T) {
		for _, point := range []Point[float64]{Pt(4.0, 4.0), Pt(-1.0, 1.0), Pt(2.0, -1.0), Pt(-1.0, -1.0)} {
			assert.False(t, s.Contains(point), point.String())
		}
	})
	t.Run("the tolerance is on the distance to a radius, however far along it", func(t *testing.T) {
		assert.True(t, s.Contains(Pt(4.9, -Delta/2)))
		assert.True(t, s.Contains(Pt(0.1, -Delta/2)))
		assert.False(t, s.Contains(Pt(4.9, -2*Delta)))
	})
	t.Run("more than half a turn holds what the rest does not", func(t *testing.T) {
		wide := Sect(Pt(0.0, 0.0), 5.0, 0, 3*Pi/2)

		assert.True(t, wide.Contains(Pt(-1.0, 1.0)))
		assert.True(t, wide.Contains(Pt(-1.0, -1.0)))
		assert.False(t, wide.Contains(Pt(1.0, -1.0)))
	})
	t.Run("a negative sweep runs the other way from the angle", func(t *testing.T) {
		assert.True(t, Sect(Pt(0.0, 0.0), 5.0, 0, -Pi/2).Contains(Pt(1.0, -1.0)))
		assert.False(t, Sect(Pt(0.0, 0.0), 5.0, 0, -Pi/2).Contains(Pt(1.0, 1.0)))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Sect(Pt(0, 0), 5, 0, Pi/2).Contains(Pt(3, 4)))
		assert.False(t, Sect(Pt(0, 0), 5, 0, Pi/2).Contains(Pt(4, 4)))
	})
	t.Run("a NaN coordinate is in no sector", func(t *testing.T) {
		assert.False(t, s.Contains(Pt(math.NaN(), 0.0)))
		assert.True(t, math.IsNaN(s.DistanceTo(Pt(0.0, math.NaN()))))
	})
}

func BenchmarkSector_Contains(b *testing.B) {
	s := Sect(Pt(0.0, 0.0), 50.0, Pi/6, Pi/3)
	inside, outside := Pt(10.0, 20.0), Pt(-99.0, 99.0)

	b.Run("inside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = s.Contains(inside)
		}
	})
	b.Run("outside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = s.Contains(outside)
		}
	})
}

func TestSector_ContainsCircle(t *testing.T) {
	s := quarterSector()

	t.Run("a circle within it, touching the radii from inside", func(t *testing.T) {
		assert.True(t, s.ContainsCircle(Circ(Pt(2.0, 2.0), 1.0)))
		assert.True(t, s.ContainsCircle(Circ(Pt(2.0, 2.0), 2.0)))
	})
	t.Run("a circle past a radius or the arc, or with its center outside", func(t *testing.T) {
		assert.False(t, s.ContainsCircle(Circ(Pt(2.0, 2.0), 2.1)))
		assert.False(t, s.ContainsCircle(Circ(Pt(3.0, 3.0), 1.0)))
		assert.False(t, s.ContainsCircle(Circ(Pt(-1.0, -1.0), 0.5)))
	})
	t.Run("a sector past half a turn holds a circle up to its center and none in its notch", func(t *testing.T) {
		assert.True(t, notchedSector().ContainsCircle(Circ(Pt(-2.0, 0.0), 2.0)))
		assert.False(t, notchedSector().ContainsCircle(Circ(Pt(-2.0, 0.0), 2.1)))
		assert.False(t, notchedSector().ContainsCircle(Circ(Pt(3.0, 0.0), 1.0)))
	})
	t.Run("the whole disc has no radius to clear", func(t *testing.T) {
		assert.True(t, Sect(Pt(0.0, 0.0), 5.0, 1, 2*Pi).ContainsCircle(Circ(Pt(1.0, 0.0), 4.0)))
		assert.False(t, Sect(Pt(0.0, 0.0), 5.0, 1, 2*Pi).ContainsCircle(Circ(Pt(1.0, 0.0), 4.1)))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Sect(Pt(0, 0), 5, 0, Pi/2).ContainsCircle(Circ(Pt(2, 2), 2)))
		assert.False(t, Sect(Pt(0, 0), 5, 0, Pi/2).ContainsCircle(Circ(Pt(2, 2), 3)))
	})
}

func TestSector_ContainsCapsule(t *testing.T) {
	s := quarterSector()

	t.Run("a capsule within it", func(t *testing.T) {
		assert.True(t, s.ContainsCapsule(Cap(Pt(1.5, 1.5), Pt(2.5, 2.5), 1.0)))
		assert.True(t, notchedSector().ContainsCapsule(Cap(Pt(-3.0, -1.0), Pt(-3.0, 1.0), 1.0)))
		assert.True(t, Sect(Pt(0.0, 0.0), 5.0, 1, 2*Pi).ContainsCapsule(Cap(Pt(-2.0, 0.0), Pt(2.0, 0.0), 3.0)))
	})
	t.Run("a capsule reaching past the arc or a radius, or leaving through the notch", func(t *testing.T) {
		assert.False(t, s.ContainsCapsule(Cap(Pt(1.5, 1.5), Pt(3.0, 3.0), 1.0)))
		assert.False(t, s.ContainsCapsule(Cap(Pt(1.5, 0.5), Pt(2.5, 2.5), 1.0)))
		assert.False(t, s.ContainsCapsule(Cap(Pt(0.5, 2.0), Pt(3.0, 0.5), 1.0)))
		assert.False(t, notchedSector().ContainsCapsule(Cap(Pt(2.0, 3.0), Pt(2.0, -3.0), 0.2)))
	})
}

func TestSector_ContainsSector(t *testing.T) {
	s := quarterSector()

	t.Run("a sector within it, and itself", func(t *testing.T) {
		assert.True(t, s.ContainsSector(Sect(Pt(1.0, 1.0), 2.0, 0, Pi/2)))
		assert.True(t, s.ContainsSector(s))
		assert.True(t, notchedSector().ContainsSector(Sect(Pt(-1.0, 0.0), 2.0, Pi/2, Pi)))
	})
	t.Run("a sector with a radius reaching out, or its arc alone", func(t *testing.T) {
		assert.False(t, s.ContainsSector(Sect(Pt(1.0, 1.0), 4.0, 0, Pi/2)))
		assert.False(t, s.ContainsSector(Sect(Pt(2.0, 2.0), 2.5, -Pi/12, 2*Pi/3)))
	})
	t.Run("an arc passing through the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, notchedSector().ContainsSector(Sect(Pt(-1.0, 0.0), 3.0, Pi/3, -2*Pi/3)))
		assert.True(t, notchedSector().ContainsSector(Sect(Pt(-1.0, 0.0), 3.0, Pi/3, 4*Pi/3)))
	})
}

func TestSector_ContainsSegment(t *testing.T) {
	s := quarterSector()

	t.Run("a segment within it or along a radius", func(t *testing.T) {
		assert.True(t, s.ContainsSegment(Seg(Pt(1.0, 1.0), Pt(3.0, 3.0))))
		assert.True(t, s.ContainsSegment(Seg(Pt(0.0, 0.0), Pt(5.0, 0.0))))
		assert.True(t, notchedSector().ContainsSegment(Seg(Pt(-2.0, 3.0), Pt(-2.0, -3.0))))
	})
	t.Run("a segment with an end outside", func(t *testing.T) {
		assert.False(t, s.ContainsSegment(Seg(Pt(1.0, 1.0), Pt(4.0, 4.0))))
		assert.False(t, s.ContainsSegment(Seg(Pt(-1.0, 1.0), Pt(1.0, 1.0))))
	})
	t.Run("a segment leaving through the notch of a sector past half a turn between ends within", func(t *testing.T) {
		assert.False(t, notchedSector().ContainsSegment(Seg(Pt(2.0, 3.0), Pt(2.0, -3.0))))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Sect(Pt(0, 0), 5, 0, Pi/2).ContainsSegment(Seg(Pt(1, 1), Pt(3, 4))))
		assert.False(t, Sect(Pt(0, 0), 5, 0, Pi/2).ContainsSegment(Seg(Pt(1, 1), Pt(4, 4))))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		segment := Seg(Pt(tenth, 1.0), Pt(3.0, 3.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = s.ContainsSegment(segment)
		}), 0)
	})
}

func TestSector_ContainsPolygon(t *testing.T) {
	s := quarterSector()

	assert.True(t, s.ContainsPolygon(Pol([]Point[float64]{Pt(1.0, 1.0), Pt(3.0, 1.0), Pt(1.0, 3.0)})))
	assert.False(t, s.ContainsPolygon(Pol([]Point[float64]{Pt(1.0, 1.0), Pt(4.0, 4.0), Pt(1.0, 3.0)})))
	assert.False(t, s.ContainsPolygon(Pol([]Point[float64]{})))
	assert.False(t, notchedSector().ContainsPolygon(Pol([]Point[float64]{Pt(2.0, 3.0), Pt(2.0, -3.0), Pt(-2.0, 0.0)})))
}

func TestSector_ContainsRectangle(t *testing.T) {
	s := quarterSector()

	assert.True(t, s.ContainsRectangle(Rect(Pt(2.0, 2.0), Sz(2.0, 2.0))))
	assert.True(t, s.ContainsRectangle(Rect(Pt(2.0, 2.0), Sz(2.0, 2.0)).Rotate(Pi/4)))
	assert.False(t, s.ContainsRectangle(Rect(Pt(2.0, 2.0), Sz(5.0, 5.0))))
}

func TestSector_ContainsRegularPolygon(t *testing.T) {
	s := quarterSector()

	assert.True(t, s.ContainsRegularPolygon(RegPol(Pt(2.0, 2.0), SzU(1.0), 6, 0, 0)))
	assert.False(t, s.ContainsRegularPolygon(RegPol(Pt(2.0, 2.0), SzU(3.0), 6, 0, 0)))
	assert.False(t, s.ContainsRegularPolygon(RegPol(Pt(2.0, 2.0), SzU(1.0), 0, 0, 0)))
}

func TestSector_ContainsBox(t *testing.T) {
	s := quarterSector()

	assert.True(t, s.ContainsBox(Bx(Pt(1.0, 1.0), Pt(3.0, 3.0))))
	assert.True(t, s.ContainsBox(Bx(Pt(0.0, 0.0), Pt(3.0, 4.0))))
	assert.False(t, s.ContainsBox(Bx(Pt(1.0, 1.0), Pt(4.0, 4.0))))
	assert.False(t, notchedSector().ContainsBox(Bx(Pt(1.0, 2.0), Pt(3.0, 3.0))))
}

func TestSector_IntersectsCircle(t *testing.T) {
	s := quarterSector()

	t.Run("a circle reaching the arc, a radius or the center", func(t *testing.T) {
		assert.True(t, s.IntersectsCircle(Circ(Pt(6.0, 8.0), 5.0)))
		assert.True(t, s.IntersectsCircle(Circ(Pt(-1.0, 1.0), 1.0)))
		assert.True(t, s.IntersectsCircle(Circ(Pt(-3.0, -4.0), 5.0)))
	})
	t.Run("a circle within it", func(t *testing.T) {
		assert.True(t, s.IntersectsCircle(Circ(Pt(2.0, 2.0), 0.5)))
	})
	t.Run("a circle short of it", func(t *testing.T) {
		assert.False(t, s.IntersectsCircle(Circ(Pt(6.0, 8.0), 4.9)))
		assert.False(t, s.IntersectsCircle(Circ(Pt(-2.0, -2.0), 1.0)))
	})
	t.Run("either side gives the same answer, the one the distance to its center gives", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, circle := range circleFixtures {
				message := fmt.Sprintf("%s → %s: ", sector, circle)

				assert.Equal(t, sector.IntersectsCircle(circle), circle.IntersectsSector(sector), message)
				assert.True(t, !sector.Contains(circle.Center) || sector.IntersectsCircle(circle), message)
				assert.True(t, sector.IntersectsCircle(circle) || sector.DistanceTo(circle.Center) > circle.Radius, message)
			}
		}
	})
}

func TestSector_IntersectsCapsule(t *testing.T) {
	s := quarterSector()

	t.Run("a capsule whose segment passes within its radius of the sector", func(t *testing.T) {
		assert.True(t, s.IntersectsCapsule(Cap(Pt(4.0, 4.0), Pt(6.0, 6.0), 1.0)))
		assert.True(t, s.IntersectsCapsule(Cap(Pt(-1.0, 1.0), Pt(-1.0, 3.0), 1.0)))
		assert.True(t, s.IntersectsCapsule(Cap(Pt(1.0, 1.0), Pt(2.0, 2.0), 0.5)))
	})
	t.Run("a capsule short of it", func(t *testing.T) {
		assert.False(t, s.IntersectsCapsule(Cap(Pt(7.0, 7.0), Pt(9.0, 9.0), 1.0)))
		assert.False(t, s.IntersectsCapsule(Cap(Pt(-1.5, 1.0), Pt(-1.5, 3.0), 1.0)))
	})
	t.Run("a capsule in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, notchedSector().IntersectsCapsule(Cap(Pt(3.0, 0.0), Pt(4.0, 0.0), 0.5)))
		assert.True(t, notchedSector().IntersectsCapsule(Cap(Pt(3.0, 0.0), Pt(4.0, 0.0), 2.5)))
	})
	t.Run("either side gives the same answer, the one the segment gives where the radius is none", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, capsule := range capsuleFixtures {
				message := fmt.Sprintf("%s → %s: ", sector, capsule)

				assert.Equal(t, sector.IntersectsCapsule(capsule), capsule.IntersectsSector(sector), message)
				assert.True(t, !sector.IntersectsSegment(capsule.Segment) || sector.IntersectsCapsule(capsule), message)
				assert.Equal(t, sector.IntersectsCapsule(Cap(capsule.Segment.Start, capsule.Segment.End, 0.0)), sector.IntersectsSegment(capsule.Segment), message)
			}
		}
	})
}

func TestSector_IntersectsSector(t *testing.T) {
	s := quarterSector()

	t.Run("a radius of one reaching the other", func(t *testing.T) {
		assert.True(t, s.IntersectsSector(Sect(Pt(6.0, 6.0), 5.0, Pi, Pi/2)))
		assert.True(t, s.IntersectsSector(Sect(Pt(0.0, 0.0), 3.0, Pi, Pi/2)))
	})
	t.Run("one within the other", func(t *testing.T) {
		assert.True(t, s.IntersectsSector(Sect(Pt(2.0, 2.0), 1.0, 0, 2*Pi)))
		assert.True(t, Sect(Pt(2.0, 2.0), 1.0, 0, 2*Pi).IntersectsSector(s))
	})
	t.Run("arcs that touch, and not a tolerance apart", func(t *testing.T) {
		assert.True(t, s.IntersectsSector(Sect(Pt(10.0, 0.0), 5.0, Pi-0.5, 1)))
		assert.False(t, s.IntersectsSector(Sect(Pt(10.0+1e-4, 0.0), 5.0, Pi-0.5, 1)))
	})
	t.Run("sectors turned away from each other", func(t *testing.T) {
		assert.False(t, s.IntersectsSector(Sect(Pt(6.0, 6.0), 5.0, 0, Pi/2)))
		assert.False(t, s.IntersectsSector(Sect(Pt(-1.0, -1.0), 5.0, Pi, Pi/2)))
	})
	t.Run("a sector in the notch of one past half a turn", func(t *testing.T) {
		assert.False(t, notchedSector().IntersectsSector(Sect(Pt(3.0, 0.0), 1.0, 0, 2*Pi)))
		assert.True(t, notchedSector().IntersectsSector(Sect(Pt(3.0, 0.0), 2.5, 0, 2*Pi)))
	})
	t.Run("either side gives the same answer and every sector meets itself", func(t *testing.T) {
		for _, a := range sectorFixtures() {
			assert.True(t, a.IntersectsSector(a), a.String())

			for _, b := range sectorFixtures() {
				assert.Equal(t, a.IntersectsSector(b), b.IntersectsSector(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("a sector meets another where its outline does, and no farther than the outline is from the arc", func(t *testing.T) {
		for _, a := range sectorFixtures() {
			for _, b := range sectorFixtures() {
				for _, offset := range vectorFixtures {
					moved := b.Translate(offset)
					outline := sectorOutline(moved, 257)
					message := fmt.Sprintf("%s → %s: ", a, moved)

					assert.True(t, !a.IntersectsPolygon(outline) || a.IntersectsSector(moved), message)
					assert.True(t, !a.IntersectsSector(moved) || a.Grow(sagitta(moved, 257)).IntersectsPolygon(outline), message)
				}
			}
		}
	})
}

func TestSector_IntersectsSegment(t *testing.T) {
	s := quarterSector()

	t.Run("a segment ending within it or passing through", func(t *testing.T) {
		assert.True(t, s.IntersectsSegment(Seg(Pt(3.0, 3.0), Pt(6.0, 6.0))))
		assert.True(t, s.IntersectsSegment(Seg(Pt(-1.0, 2.0), Pt(2.0, -1.0))))
		assert.True(t, s.IntersectsSegment(Seg(Pt(0.0, 6.0), Pt(6.0, 0.0))))
		assert.True(t, s.IntersectsSegment(Seg(Pt(1.0, 1.0), Pt(2.0, 2.0))))
	})
	t.Run("a segment touching the arc, a radius or the center", func(t *testing.T) {
		assert.True(t, s.IntersectsSegment(Seg(Pt(5.0, -1.0), Pt(5.0, 1.0))))
		assert.True(t, s.IntersectsSegment(Seg(Pt(2.0, -1.0), Pt(2.0, 0.0))))
		assert.True(t, s.IntersectsSegment(Seg(Pt(-2.0, 0.0), Pt(0.0, 0.0))))
	})
	t.Run("a segment short of it or beside it", func(t *testing.T) {
		assert.False(t, s.IntersectsSegment(Seg(Pt(4.0, 4.0), Pt(6.0, 6.0))))
		assert.False(t, s.IntersectsSegment(Seg(Pt(-1.0, -1.0), Pt(-1.0, 5.0))))
		assert.False(t, s.IntersectsSegment(Seg(Pt(5.0, -1.0), Pt(5.0, -1e-4))))
	})
	t.Run("a segment in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, notchedSector().IntersectsSegment(Seg(Pt(6.0, 1.0), Pt(1.0, 0.0))))
		assert.True(t, notchedSector().IntersectsSegment(Seg(Pt(6.0, 1.0), Pt(-1.0, 0.0))))
		assert.True(t, notchedSector().IntersectsSegment(Seg(Pt(3.0, 4.0), Pt(3.0, -4.0))))
	})
	t.Run("int", func(t *testing.T) {
		assert.True(t, Sect(Pt(0, 0), 5, 0, Pi/2).IntersectsSegment(Seg(Pt(3, 4), Pt(6, 8))))
		assert.False(t, Sect(Pt(0, 0), 5, 0, Pi/2).IntersectsSegment(Seg(Pt(4, 4), Pt(6, 8))))
	})
	t.Run("either side gives the same answer, and a segment with an end within meets it", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, segment := range segmentFixtures {
				message := fmt.Sprintf("%s → %s: ", sector, segment)

				assert.Equal(t, sector.IntersectsSegment(segment), segment.IntersectsSector(sector), message)
				assert.True(t, !(sector.Contains(segment.Start) || sector.Contains(segment.End)) || sector.IntersectsSegment(segment), message)
			}
		}
	})
	t.Run("a segment meets the sector where it meets its outline, and no farther than the outline is from the arc", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			outline := sectorOutline(sector, 257)

			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					segment := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", sector, segment)

					assert.True(t, !outline.IntersectsSegment(segment) || sector.IntersectsSegment(segment), message)
					assert.True(t, !sector.IntersectsSegment(segment) || Distance[float64](outline, segment) <= sagitta(sector, 257)+Delta, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		segment := Seg(Pt(tenth, 6.0), Pt(6.0, 0.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = s.IntersectsSegment(segment)
		}), 0)
	})
}

func TestSector_IntersectsRay(t *testing.T) {
	s := quarterSector()

	t.Run("a ray toward it and one away from it", func(t *testing.T) {
		assert.True(t, s.IntersectsRay(Ry(Pt(-1.0, 2.0), Vec(1.0, 0.0))))
		assert.False(t, s.IntersectsRay(Ry(Pt(-1.0, 2.0), Vec(-1.0, 0.0))))
	})
	t.Run("a ray out through the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, notchedSector().IntersectsRay(Ry(Pt(1.0, 0.0), Vec(1.0, 0.0))))
		assert.True(t, notchedSector().IntersectsRay(Ry(Pt(1.0, 0.0), Vec(-1.0, 0.0))))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, ray := range rayFixtures {
				assert.Equal(t, sector.IntersectsRay(ray), ray.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, ray))
			}
		}
	})
}

func TestSector_IntersectsLine(t *testing.T) {
	s := quarterSector()

	t.Run("a line through it and one beside it", func(t *testing.T) {
		assert.True(t, s.IntersectsLine(LineThrough(Pt(-1.0, 2.0), Pt(-2.0, 2.0))))
		assert.True(t, s.IntersectsLine(LineThrough(Pt(-1.0, 0.0), Pt(-2.0, 0.0))))
		assert.False(t, s.IntersectsLine(LineThrough(Pt(0.0, -1.0), Pt(1.0, -1.0))))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, line := range lineFixtures {
				assert.Equal(t, sector.IntersectsLine(line), line.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, line))
			}
		}
	})
}

func TestSector_IntersectsPolygon(t *testing.T) {
	s := quarterSector()

	t.Run("a polygon crossing it, around it or within it", func(t *testing.T) {
		assert.True(t, s.IntersectsPolygon(Pol([]Point[float64]{Pt(3.0, 3.0), Pt(6.0, 3.0), Pt(6.0, 6.0)})))
		assert.True(t, s.IntersectsPolygon(Pol([]Point[float64]{Pt(-10.0, -10.0), Pt(20.0, -10.0), Pt(-10.0, 20.0)})))
		assert.True(t, s.IntersectsPolygon(Pol([]Point[float64]{Pt(1.0, 1.0), Pt(2.0, 1.0), Pt(1.0, 2.0)})))
	})
	t.Run("a polygon short of it, an empty one, and one in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, s.IntersectsPolygon(Pol([]Point[float64]{Pt(4.0, 4.0), Pt(6.0, 4.0), Pt(6.0, 6.0)})))
		assert.False(t, s.IntersectsPolygon(Pol([]Point[float64]{})))
		assert.False(t, notchedSector().IntersectsPolygon(Pol([]Point[float64]{Pt(2.0, 0.5), Pt(2.0, -0.5), Pt(4.0, 0.0)})))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, polygon := range polygonFixtures() {
				assert.Equal(t, sector.IntersectsPolygon(polygon), polygon.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, polygon))
			}
		}
	})
}

func TestSector_IntersectsRectangle(t *testing.T) {
	s := quarterSector()

	t.Run("a rectangle crossing it, around it or within it", func(t *testing.T) {
		assert.True(t, s.IntersectsRectangle(Rect(Pt(4.0, 4.0), Sz(2.0, 2.0))))
		assert.True(t, s.IntersectsRectangle(Rect(Pt(0.0, 0.0), Sz(20.0, 20.0)).Rotate(Pi/6)))
		assert.True(t, s.IntersectsRectangle(Rect(Pt(2.0, 2.0), Sz(1.0, 1.0))))
	})
	t.Run("a rectangle short of it, and one in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, s.IntersectsRectangle(Rect(Pt(5.0, 5.0), Sz(2.0, 2.0))))
		assert.False(t, notchedSector().IntersectsRectangle(Rect(Pt(2.5, 0.0), Sz(1.0, 1.0))))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, rectangle := range rectFixtures {
				assert.Equal(t, sector.IntersectsRectangle(rectangle), rectangle.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, rectangle))
			}
		}
	})
}

func TestSector_IntersectsRegularPolygon(t *testing.T) {
	s := quarterSector()

	t.Run("a polygon crossing it, around it or within it", func(t *testing.T) {
		assert.True(t, s.IntersectsRegularPolygon(RegPol(Pt(5.0, 5.0), SzU(3.0), 6, 0, 0)))
		assert.True(t, s.IntersectsRegularPolygon(RegPol(Pt(0.0, 0.0), SzU(20.0), 6, 0, 0)))
		assert.True(t, s.IntersectsRegularPolygon(RegPol(Pt(2.0, 2.0), SzU(0.5), 6, 0, 0)))
	})
	t.Run("a polygon short of it, an empty one, and one in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, s.IntersectsRegularPolygon(RegPol(Pt(6.0, 6.0), SzU(1.0), 6, 0, 0)))
		assert.False(t, s.IntersectsRegularPolygon(RegPol(Pt(2.0, 2.0), SzU(0.5), 0, 0, 0)))
		assert.False(t, notchedSector().IntersectsRegularPolygon(RegPol(Pt(3.0, 0.0), SzU(0.5), 6, 0, 0)))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, polygon := range regularPolygonFixtures {
				assert.Equal(t, sector.IntersectsRegularPolygon(polygon), polygon.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, polygon))
			}
		}
	})
}

func TestSector_IntersectsBox(t *testing.T) {
	s := quarterSector()

	t.Run("a box crossing it, around it or within it", func(t *testing.T) {
		assert.True(t, s.IntersectsBox(Bx(Pt(3.0, 3.0), Pt(6.0, 6.0))))
		assert.True(t, s.IntersectsBox(Bx(Pt(-1.0, -1.0), Pt(1.0, 1.0))))
		assert.True(t, s.IntersectsBox(Bx(Pt(-10.0, -10.0), Pt(10.0, 10.0))))
		assert.True(t, s.IntersectsBox(Bx(Pt(1.0, 1.0), Pt(2.0, 2.0))))
	})
	t.Run("a box short of it, and one in the notch of a sector past half a turn", func(t *testing.T) {
		assert.False(t, s.IntersectsBox(Bx(Pt(4.0, 4.0), Pt(6.0, 6.0))))
		assert.False(t, notchedSector().IntersectsBox(Bx(Pt(2.0, -0.5), Pt(3.0, 0.5))))
	})
	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, box := range boxFixtures {
				assert.Equal(t, sector.IntersectsBox(box), box.IntersectsSector(sector), fmt.Sprintf("%s → %s: ", sector, box))
			}
		}
	})
}

func TestSector_IntersectionSegment(t *testing.T) {
	s := quarterSector()

	t.Run("the arc and a radius, from Start", func(t *testing.T) {
		geomtest.AssertVertices(t, s.IntersectionSegment(Seg(Pt(3.0, 10.0), Pt(3.0, -10.0))), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, 0.0)})
		geomtest.AssertVertices(t, s.IntersectionSegment(Seg(Pt(3.0, -10.0), Pt(3.0, 10.0))), []Point[float64]{Pt(3.0, 0.0), Pt(3.0, 4.0)})
	})
	t.Run("the center is held once", func(t *testing.T) {
		geomtest.AssertVertices(t, s.IntersectionSegment(Seg(Pt(-1.0, -1.0), Pt(1.0, 1.0))), []Point[float64]{Pt(0.0, 0.0)})
	})
	t.Run("through the notch of a sector past half a turn there are four", func(t *testing.T) {
		geomtest.AssertVertices(t, notchedSector().IntersectionSegment(Seg(Pt(3.0, 10.0), Pt(3.0, -10.0))), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, 3.0), Pt(3.0, -3.0), Pt(3.0, -4.0)})
	})
	t.Run("the whole disc is crossed as its circle", func(t *testing.T) {
		geomtest.AssertVertices(t, Sect(Pt(0.0, 0.0), 5.0, 1, 2*Pi).IntersectionSegment(Seg(Pt(-10.0, 0.0), Pt(10.0, 0.0))), []Point[float64]{Pt(-5.0, 0.0), Pt(5.0, 0.0)})
	})
	t.Run("inside or apart gives none", func(t *testing.T) {
		assert.Nil(t, s.IntersectionSegment(Seg(Pt(1.0, 1.0), Pt(2.0, 3.0))))
		assert.Nil(t, s.IntersectionSegment(Seg(Pt(-1.0, -1.0), Pt(-1.0, 9.0))))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Sect(Pt(0, 0), 5, 0, Pi/2).IntersectionSegment(Seg(Pt(3, 10), Pt(3, -10))), []Point[int]{Pt(3, 4), Pt(3, 0)})
	})
	t.Run("every point lies on the segment and the sector, where the two intersect", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					segment := Seg(start, end)
					points := sector.IntersectionSegment(segment)
					message := fmt.Sprintf("%s → %s: ", sector, segment)

					for _, point := range points {
						assert.True(t, segment.Contains(point) && sector.Contains(point), message+point.String())
					}
					assertOrdered(t, segment, points, message)
					assert.True(t, len(points) <= 4, message)
					assert.True(t, len(points) == 0 || sector.IntersectsSegment(segment), message)
				}
			}
		}
	})
	t.Run("a graze within the tolerance is a crossing, far from the origin in float32 and in int too", func(t *testing.T) {
		assertSectorGrazes(t, Vec(0.0, 0.0))
		assertSectorGrazes(t, Vec(0, 0))
		assertSectorGrazes(t, Vec(1<<25, -1<<25))
		for _, offset := range farOffsets {
			assertSectorGrazes(t, offset)
		}
	})
	t.Run("appends to a buffer with room without allocating", func(t *testing.T) {
		buffer, through := make([]Point[float64], 0, 4), Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = len(s.AppendIntersectionSegment(buffer[:0], through)) > 0
		}), 0)
		geomtest.AssertVertices(t, s.AppendIntersectionSegment(bufferWith(prefixPoint), through)[:1], []Point[float64]{prefixPoint})
	})
}

// assertSectorGrazes asserts, over the sectors and the segments of a small lattice moved by the
// offset, where tangents, corners and ends on the boundary are common, that a segment crosses a
// sector where the two intersect, unless it lies within, has a part inside exactly there, and
// hits it from outside at the first of the crossings.
func assertSectorGrazes[T Number](t *testing.T, offset Vector[T]) {
	t.Helper()

	var lattice []Point[T]
	for _, x := range []T{-2, 0, 1, 3} {
		for _, y := range []T{-2, 0, 1, 3} {
			lattice = append(lattice, Pt(x, y).Add(offset))
		}
	}

	for _, radius := range []T{0, 1, 2} {
		for _, sweep := range []float64{0, Pi / 2, -3 * Pi / 4, 3 * Pi / 2, 2 * Pi} {
			for i, center := range lattice {
				sector := Sect(center, radius, float64(i)*Pi/4, sweep)

				for _, start := range lattice {
					for _, end := range lattice {
						s := Seg(start, end)
						points := sector.IntersectionSegment(s)
						intersects := sector.IntersectsSegment(s)
						message := fmt.Sprintf("%s → %s: ", sector, s)

						assert.True(t, len(points) == 0 || intersects, message)
						assert.True(t, len(points) > 0 || !intersects || sector.Contains(start) && sector.Contains(end), message)
						assert.Equal(t, len(s.ClipSector(sector)) > 0, intersects, message)

						hit, ok := s.HitSector(sector)
						assert.Equal(t, ok, !sector.Contains(start) && intersects, message)
						assert.True(t, !ok || hit.Point == points[0], message)
					}
				}
			}
		}
	}
}

func TestSector_IntersectionRay(t *testing.T) {
	s := quarterSector()

	geomtest.AssertVertices(t, s.IntersectionRay(Ry(Pt(3.0, 10.0), Vec(0.0, -1.0))), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, 0.0)})
	geomtest.AssertVertices(t, s.IntersectionRay(Ry(Pt(3.0, 1.0), Vec(0.0, -1.0))), []Point[float64]{Pt(3.0, 0.0)})
	assert.Nil(t, s.IntersectionRay(Ry(Pt(3.0, 10.0), Vec(0.0, 1.0))))
}

func TestSector_IntersectionLine(t *testing.T) {
	s := quarterSector()

	geomtest.AssertVertices(t, s.IntersectionLine(Ln(Pt(3.0, 1.0), Vec(0.0, -1.0))), []Point[float64]{Pt(3.0, 4.0), Pt(3.0, 0.0)})
	assert.Nil(t, s.IntersectionLine(Ln(Pt(3.0, -1.0), Vec(1.0, 0.0))))
}

func TestSector_Equal(t *testing.T) {
	t.Run("equal, a full turn of the angle aside", func(t *testing.T) {
		assert.True(t, Sect(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Sect(Pt(1.0, 2.0), 5.0, 0.5+2*Pi, 1.5)))
	})
	t.Run("a different radius or sweep", func(t *testing.T) {
		assert.False(t, Sect(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Sect(Pt(1.0, 2.0), 6.0, 0.5, 1.5)))
		assert.False(t, Sect(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Sect(Pt(1.0, 2.0), 5.0, 0.5, -1.5)))
	})
}

func TestSector_IsZero(t *testing.T) {
	assert.True(t, Sector[int]{}.IsZero())
	assert.False(t, Sect(Pt(0, 0), 1, 0, 0).IsZero())
}

func TestSector_Arc(t *testing.T) {
	geomtest.AssertArc(t, Sect(Pt(1, 2), 5, 0.5, 1.5).Arc(), Ac(Pt(1, 2), 5, 0.5, 1.5))
}

func TestSector_Circle(t *testing.T) {
	geomtest.AssertCircle(t, Sect(Pt(1, 2), 5, 0.5, 1.5).Circle(), Circ(Pt(1, 2), 5))
}

func TestSector_Polygon(t *testing.T) {
	t.Run("the center and the arc as chords", func(t *testing.T) {
		geomtest.AssertPolygon(t, quarterSector().Polygon(2), Pol([]Point[float64]{Pt(0.0, 0.0), Pt(5.0, 0.0), Pt(0.0, 5.0)}))
		geomtest.AssertPolygon(t, quarterSector().Polygon(-0.5), Pol([]Point[float64]{Pt(0.0, 0.0), Pt(5.0, 0.0), Pt(5/math.Sqrt2, 5/math.Sqrt2), Pt(0.0, 5.0)}))
	})
	t.Run("a sector swept the other way is wound the same", func(t *testing.T) {
		geomtest.AssertPolygon(t, Sect(Pt(0.0, 0.0), 5.0, Pi/2, -Pi/2).Polygon(2), quarterSector().Polygon(2))
	})
	t.Run("the whole disc has no center among its vertices", func(t *testing.T) {
		geomtest.AssertPolygon(t, Sect(Pt(0, 0), 5, 0, 2*Pi).Polygon(2), Pol([]Point[int]{Pt(5, 0), Pt(0, 5), Pt(-5, 0), Pt(0, -5)}))
	})
	t.Run("no sweep gives the radius, and no radius the center", func(t *testing.T) {
		geomtest.AssertPolygon(t, Sect(Pt(1, 2), 5, 0, 0).Polygon(1), Pol([]Point[int]{Pt(1, 2), Pt(6, 2)}))
		geomtest.AssertPolygon(t, Sect(Pt(1, 2), 0, 0, 1).Polygon(1), Pol([]Point[int]{Pt(1, 2)}))
	})
	t.Run("the vertices lie on the boundary, the edges within the tolerance of it, wound as a rectangle", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			for _, tolerance := range []float64{0.01, 0.25} {
				polygon := s.Polygon(tolerance)
				message := fmt.Sprintf("%s ±%v: ", s, tolerance)

				for vertex := range polygon.Vertices() {
					assert.True(t, s.Contains(vertex), message+vertex.String())
				}
				for edge := range polygon.Edges() {
					assert.True(t, LessOrEqual(s.Arc().DistanceTo(edge.Midpoint()), tolerance) || s.Contains(edge.Midpoint()), message+edge.String())
				}

				assert.True(t, polygon.Area() < Delta || polygon.Winding() == WindingClockwise, message)
			}
		}
	})
	t.Run("a fine outline has the area and the perimeter", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			polygon := s.Polygon(0)

			assert.EqualDelta(t, polygon.Area(), s.Area(), 1e-3, s.String())
			assert.EqualDelta(t, polygon.Perimeter(), s.Perimeter(), 1e-3, s.String())
		}
	})
}

func TestSector_AppendPolygon(t *testing.T) {
	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, quarterSector().AppendPolygon([]Point[float64]{Pt(9.0, 9.0)}, 2), []Point[float64]{Pt(9.0, 9.0), Pt(0.0, 0.0), Pt(5.0, 0.0), Pt(0.0, 5.0)})
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		sector := Sect(Pt(0, 0), 50, 0, 1.5)
		buffer := make([]Point[int], 0, 64)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = sector.Polygon(tenth).Points
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = sector.AppendPolygon(buffer[:0], tenth)
		}), 0)
	})
}

func TestSector_Cast(t *testing.T) {
	s := Sect(Pt(1.5, -2.5), 3.5, 0.5, 1.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertSector(t, s.Cast[int](), s.Int())
		geomtest.AssertSector(t, s.Cast[float64](), s.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertSector(t, s.Cast[int8](), Sect(Pt[int8](2, -3), int8(4), 0.5, 1.5))
	})
}

func TestSector_Int(t *testing.T) {
	geomtest.AssertSector(t, Sect(Pt(0.6, -0.25), 1.2, 0.5, 1.5).Int(), Sect(Pt(1, 0), 1, 0.5, 1.5))
}

func TestSector_Float(t *testing.T) {
	geomtest.AssertSector(t, Sect(Pt(1, 2), 5, 0.5, 1.5).Float(), Sect(Pt(1.0, 2.0), 5.0, 0.5, 1.5))
}

func TestSector_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Sect(Pt(1, 2), 5, 0.5, -1.5).String(), "Sect((1,2);5;0.50;-1.50)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Sect(Pt(1, 2), 5, 0.5, -1.5).String)
	})
}

func TestSector_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Sect(Pt(1, 2), 5, 0.5, -1.5), `{"x":1,"y":2,"r":5,"a":0.5,"s":-1.5}`)

		var s Sector[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"r":5,"a":0.5,"s":-1.5}`), &s))
		geomtest.AssertSector(t, s, Sect(Pt(1, 2), 5, 0.5, -1.5))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			data, err := json.Marshal(sector)
			assert.NoError(t, err)

			var decoded Sector[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, sector)
		}
	})
}

func TestSector_Properties(t *testing.T) {
	t.Run("distance is zero exactly where contains holds, and the nearest point is at it", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			for _, point := range pointFixtures {
				nearest := s.Nearest(point)
				message := fmt.Sprintf("%s → %s: ", s, point)

				assert.Equal(t, s.DistanceTo(point) == 0, s.Contains(point), message)
				assert.Equal(t, nearest.Equal(point), s.Contains(point), message)
				assert.True(t, s.Contains(nearest), message)
				assert.True(t, !s.Contains(point) || s.Bounds().Contains(point), message)
				geomtest.AssertNumber(t, point.DistanceTo(nearest), s.DistanceTo(point), message)
			}
		}
	})
	t.Run("the mass properties and the distances are those of the outline walked as a polygon", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			if s.Area() == 0 {
				continue
			}

			outline := sectorOutline(s, 4096)

			assert.True(t, math.Abs(s.Area()-outline.Area()) <= 1e-5*s.Area(), fmt.Sprintf("%s: %v ≠ %v", s, s.Area(), outline.Area()))
			assert.True(t, math.Abs(s.Inertia()-outline.Inertia()) <= 1e-5*s.Inertia(), fmt.Sprintf("%s: %v ≠ %v", s, s.Inertia(), outline.Inertia()))
			assert.True(t, s.Centroid().DistanceTo(outline.Centroid()) <= 1e-5*s.Radius, fmt.Sprintf("%s: %v ≠ %v", s, s.Centroid(), outline.Centroid()))
			assert.True(t, s.Contains(s.Centroid()), s.String())

			for _, point := range pointFixtures {
				assert.True(t, math.Abs(s.DistanceTo(point)-outline.DistanceTo(point)) <= 1e-5*s.Radius, fmt.Sprintf("%s → %s: %v ≠ %v", s, point, s.DistanceTo(point), outline.DistanceTo(point)))
				assert.True(t, !outline.Contains(point) || s.Contains(point), fmt.Sprintf("%s → %s: ", s, point))
			}
		}
	})
	t.Run("the boundary is contained and its bounds hold it", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			for _, fraction := range []float64{0, 0.1, 0.35, 0.5, 0.8, 1} {
				point := s.Arc().PointAt(fraction)
				message := fmt.Sprintf("%s @%v: ", s, fraction)

				assert.True(t, s.Contains(point), message)
				assert.True(t, s.Contains(s.Center.Midpoint(point)), message)
				assert.True(t, s.Bounds().Contains(point), message)
			}

			assert.True(t, s.Contains(s.Center), s.String())
			assert.True(t, s.Bounds().Contains(s.Center), s.String())
		}
	})
	t.Run("the whole disc answers as the circle", func(t *testing.T) {
		for _, c := range circleFixtures {
			s := Sect(c.Center, c.Radius, 1, -2*Pi)

			geomtest.AssertNumber(t, s.Area(), c.Area(), s.String())
			geomtest.AssertNumber(t, s.Perimeter(), c.Perimeter(), s.String())

			for _, point := range pointFixtures {
				assert.Equal(t, s.Contains(point), c.Contains(point), fmt.Sprintf("%s → %s: ", s, point))
				geomtest.AssertNumber(t, s.DistanceTo(point), c.DistanceTo(point), fmt.Sprintf("%s → %s: ", s, point))
			}
		}
	})
	t.Run("a turned or moved sector is as far from a point turned or moved with it", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			for _, point := range pointFixtures {
				message := fmt.Sprintf("%s → %s: ", s, point)

				geomtest.AssertNumber(t, s.Rotate(Pi/3).DistanceTo(point.RotateAround(s.Center, Pi/3)), s.DistanceTo(point), message)
				geomtest.AssertNumber(t, s.Translate(Vec(3.0, -2.0)).DistanceTo(point.Add(Vec(3.0, -2.0))), s.DistanceTo(point), message)
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, s := range sectorFixtures() {
			for _, factor := range []float64{0.5, 2.5, -3} {
				assert.True(t, s.Scale(factor).Unscale(factor).Equal(s), fmt.Sprintf("%s ×%v: ", s, factor))
			}
		}
	})
	t.Run("far from the origin a float32 sector contains its boundary and its nearest points", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, fixture := range sectorFixtures() {
				s := fixture.Cast[float32]().Translate(offset)

				for _, fraction := range []float64{0, 0.35, 1} {
					assert.True(t, s.Contains(s.Arc().PointAt(fraction)), fmt.Sprintf("%s @%v: ", s, fraction))
				}
				for _, p := range pointFixtures {
					point := p.Cast[float32]().Add(offset)

					assert.True(t, s.Contains(s.Nearest(point)), fmt.Sprintf("%s → %s: ", s, point))
				}
			}
		}
	})
}

func TestSector_Immutable(t *testing.T) {
	s := Sect(Pt(1, 2), 5, 0.5, 1.5)

	s.Translate(Vec(3, -2))
	s.MoveTo(Pt(4, 3))
	s.Rotate(1)
	s.Scale(2)
	s.Resize(9)
	s.Grow(1)
	s.Shrink(1)
	s.Polygon(1)

	geomtest.AssertSector(t, s, Sect(Pt(1, 2), 5, 0.5, 1.5))
}

// quarterSector returns the sector of radius five about the origin from the X axis a quarter
// turn on.
func quarterSector() Sector[float64] {
	return Sect(Pt(0.0, 0.0), 5.0, 0, Pi/2)
}

// notchedSector returns the sector of radius five about the origin that leaves out the quarter
// turn about the X axis: a sector past half a turn, with a notch between its two radii.
func notchedSector() Sector[float64] {
	return Sect(Pt(0.0, 0.0), 5.0, Pi/4, 3*Pi/2)
}

// sectorFixtures returns the sectors the arc fixtures bound.
func sectorFixtures() []Sector[float64] {
	sectors := make([]Sector[float64], len(arcFixtures))
	for i, arc := range arcFixtures {
		sectors[i] = arc.Sector()
	}

	return sectors
}

// sectorOutline returns the polygon of the center and n points along the arc of the sector,
// the center left out for the whole disc: the outline the sector is checked against.
func sectorOutline(s Sector[float64], n int) Polygon[float64] {
	var points []Point[float64]
	if math.Abs(s.Sweep) < 2*Pi {
		points = append(points, s.Center)
	}
	for i := range n {
		points = append(points, s.Arc().PointAt(float64(i)/float64(n-1)))
	}

	return Pol(points)
}

// sagitta returns how far the outline of n points along the arc of the sector lies within the
// arc at most, the height of the arc over one of its chords.
func sagitta(s Sector[float64], n int) float64 {
	return s.Radius * (1 - math.Cos(math.Min(math.Abs(s.Sweep), 2*Pi)/float64(2*(n-1))))
}

func ExampleSect() {
	fmt.Println(Sect(Pt(1, 2), 5, 0, 1.5))
	// Output: Sect((1,2);5;0.00;1.50)
}

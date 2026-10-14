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

func TestArc_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, 1.5), Arc[int]{Center: Pt(1, 2), Radius: 5, Angle: 0.5, Sweep: 1.5})
	})
	t.Run("a negative radius is taken absolute and a negative sweep kept", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1.0, 2.0), -5.0, 0.5, -1.5), Arc[float64]{Center: Pt(1.0, 2.0), Radius: 5, Angle: 0.5, Sweep: -1.5})
	})
}

func TestArcThrough(t *testing.T) {
	t.Run("from the first point through the second to the third", func(t *testing.T) {
		arc, ok := ArcThrough(Pt(5.0, 0.0), Pt(0.0, 5.0), Pt(-5.0, 0.0))

		assert.True(t, ok)
		geomtest.AssertArc(t, arc, Ac(Pt(0.0, 0.0), 5.0, 0, Pi))
	})
	t.Run("swept the way that passes the second point", func(t *testing.T) {
		arc, ok := ArcThrough(Pt(5.0, 0.0), Pt(0.0, -5.0), Pt(-5.0, 0.0))

		assert.True(t, ok)
		geomtest.AssertArc(t, arc, Ac(Pt(0.0, 0.0), 5.0, 0, -Pi))
	})
	t.Run("int rounds the center and the radius", func(t *testing.T) {
		arc, ok := ArcThrough(Pt(5, 0), Pt(0, 5), Pt(-5, 0))

		assert.True(t, ok)
		geomtest.AssertArc(t, arc, Ac(Pt(0, 0), 5, 0, Pi))
	})
	t.Run("points on one line have no arc", func(t *testing.T) {
		_, ok := ArcThrough(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(2.0, 2.0))
		assert.False(t, ok)

		_, ok = ArcThrough(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(0.0, 0.0))
		assert.False(t, ok)
	})
	t.Run("the arc through three points of an arc is that arc", func(t *testing.T) {
		for _, a := range arcFixtures {
			if a.Length() < Delta || math.Abs(a.Sweep) >= 2*Pi {
				continue
			}

			arc, ok := ArcThrough(a.Start(), a.Midpoint(), a.End())

			assert.True(t, ok, a.String())
			geomtest.AssertPoint(t, arc.Start(), a.Start(), a.String())
			geomtest.AssertPoint(t, arc.Midpoint(), a.Midpoint(), a.String())
			geomtest.AssertPoint(t, arc.End(), a.End(), a.String())
		}
	})
}

func TestArc_Length(t *testing.T) {
	t.Run("the radius times the sweep, either way", func(t *testing.T) {
		geomtest.AssertNumber(t, Ac(Pt(0, 0), 5, 0, Pi/2).Length(), 2.5*Pi)
		geomtest.AssertNumber(t, Ac(Pt(0, 0), 5, 1, -Pi/2).Length(), 2.5*Pi)
	})
	t.Run("no more than the circumference", func(t *testing.T) {
		geomtest.AssertNumber(t, Ac(Pt(0, 0), 5, 0, 7*Pi).Length(), 10*Pi)
	})
}

func TestArc_Start(t *testing.T) {
	geomtest.AssertPoint(t, quarterArc().Start(), Pt(5.0, 0.0))
}

func TestArc_End(t *testing.T) {
	t.Run("at the angle plus the sweep", func(t *testing.T) {
		geomtest.AssertPoint(t, quarterArc().End(), Pt(0.0, 5.0))
		geomtest.AssertPoint(t, Ac(Pt(0.0, 0.0), 5.0, 0, -Pi/2).End(), Pt(0.0, -5.0))
	})
	t.Run("the whole circle ends where it starts", func(t *testing.T) {
		geomtest.AssertPoint(t, Ac(Pt(0.0, 0.0), 5.0, 1, 9).End(), Ac(Pt(0.0, 0.0), 5.0, 1, 9).Start())
	})
}

func TestArc_Midpoint(t *testing.T) {
	geomtest.AssertPoint(t, quarterArc().Midpoint(), Pt(5/math.Sqrt2, 5/math.Sqrt2))
}

func TestArc_Chord(t *testing.T) {
	t.Run("from the start to the end", func(t *testing.T) {
		geomtest.AssertSegment(t, quarterArc().Chord(), Seg(Pt(5.0, 0.0), Pt(0.0, 5.0)))
	})
	t.Run("the whole circle has a chord of no length", func(t *testing.T) {
		assert.True(t, Ac(Pt(1.0, 2.0), 5.0, 1, 2*Pi).Chord().IsPoint())
	})
}

func TestArc_Bounds(t *testing.T) {
	t.Run("the ends alone where the sweep passes no axis", func(t *testing.T) {
		geomtest.AssertBox(t, quarterArc().Bounds(), Bx(Pt(0.0, 0.0), Pt(5.0, 5.0)))
	})
	t.Run("out to the circle along each axis the sweep passes", func(t *testing.T) {
		geomtest.AssertBox(t, Ac(Pt(0.0, 0.0), 5.0, Pi/4, Pi/2).Bounds(), Bx(Pt(-5/math.Sqrt2, 5/math.Sqrt2), Pt(5/math.Sqrt2, 5.0)))
		geomtest.AssertBox(t, Ac(Pt(0.0, 0.0), 5.0, Pi/4, -Pi/2).Bounds(), Bx(Pt(5/math.Sqrt2, -5/math.Sqrt2), Pt(5.0, 5/math.Sqrt2)))
	})
	t.Run("the whole circle", func(t *testing.T) {
		geomtest.AssertBox(t, Ac(Pt(1.0, 2.0), 5.0, 1, 2*Pi).Bounds(), Circ(Pt(1.0, 2.0), 5.0).Bounds())
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertBox(t, Ac(Pt(1, 2), 5, 0, Pi).Bounds(), Bx(Pt(-4, 2), Pt(6, 7)))
		geomtest.AssertBox(t, Ac(Pt(0, 0), 5, Pi/2, Pi/2).Bounds(), Bx(Pt(-5, 0), Pt(0, 5)))
	})
	t.Run("int rounds the ends outward, so the box holds the arc", func(t *testing.T) {
		arc := Ac(Pt(0, 0), 10, 0.3, 0.5)

		geomtest.AssertBox(t, arc.Bounds(), Bx(Pt(6, 2), Pt(10, 8)))
		assert.True(t, arc.Bounds().Float().ContainsBox(arc.Float().Bounds()))
	})
}

func TestArc_Translate(t *testing.T) {
	geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, 1.5).Translate(Vec(3, -2)), Ac(Pt(4, 0), 5, 0.5, 1.5))
}

func TestArc_MoveTo(t *testing.T) {
	geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, 1.5).MoveTo(Pt(3, -2)), Ac(Pt(3, -2), 5, 0.5, 1.5))
}

func TestArc_Rotate(t *testing.T) {
	t.Run("the angle moves and the sweep stays", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, 1.5).Rotate(1), Ac(Pt(1, 2), 5, 1.5, 1.5))
	})
	t.Run("normalized to a single turn", func(t *testing.T) {
		assert.Equal(t, Ac(Pt(1, 2), 5, Pi, 1.5).Rotate(3*Pi/2).Angle, Pi/2)
	})
}

func TestArc_Scale(t *testing.T) {
	t.Run("the radius", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Scale(2.5), Ac(Pt(1, 2), 10, 0.5, 1.5))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Scale(-2), Ac(Pt(1, 2), 8, 0.5, 1.5))
	})
}

func TestArc_Unscale(t *testing.T) {
	t.Run("the radius", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 10, 0.5, 1.5).Unscale(-2.5), Ac(Pt(1, 2), 4, 0.5, 1.5))
	})
	t.Run("a zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Ac(Pt(1, 2), 10, 0.5, 1.5).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestArc_Resize(t *testing.T) {
	geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Resize(-9), Ac(Pt(1, 2), 9, 0.5, 1.5))
}

func TestArc_Grow(t *testing.T) {
	t.Run("the radius alone", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Grow(3), Ac(Pt(1, 2), 7, 0.5, 1.5))
	})
	t.Run("clamps at zero", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Grow(-9), Ac(Pt(1, 2), 0, 0.5, 1.5))
	})
}

func TestArc_Shrink(t *testing.T) {
	t.Run("the radius alone", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Shrink(3), Ac(Pt(1, 2), 1, 0.5, 1.5))
	})
	t.Run("clamps at zero", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 4, 0.5, 1.5).Shrink(9), Ac(Pt(1, 2), 0, 0.5, 1.5))
	})
}

func TestArc_Lerp(t *testing.T) {
	t.Run("the center, the radius, the angle and the sweep move together", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(0, 0), 4, 0, 1).Lerp(Ac(Pt(10, 0), 8, 2, 3), 0.5), Ac(Pt(5, 0), 6, 1, 2))
	})
	t.Run("the angle turns along the shorter arc", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(0, 0), 4, 0.5, 1).Lerp(Ac(Pt(0, 0), 4, 2*Pi-0.5, 1), 0.5), Ac(Pt(0, 0), 4, 0, 1))
	})
	t.Run("ends on the two arcs", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, b := range arcFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
}

func TestArc_Canonical(t *testing.T) {
	t.Run("a well-formed arc is unchanged", func(t *testing.T) {
		geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, -1.5).Canonical(), Ac(Pt(1, 2), 5, 0.5, -1.5))
	})
	t.Run("the radius, the angle and a sweep past a full turn are repaired", func(t *testing.T) {
		canonical := Arc[int]{Pt(1, 2), -5, -Pi / 2, -3 * Pi}.Canonical()

		assert.Equal(t, canonical, Arc[int]{Pt(1, 2), 5, 3 * Pi / 2, -2 * Pi})
	})
}

func TestArc_Reverse(t *testing.T) {
	t.Run("starts where it ended and sweeps back", func(t *testing.T) {
		geomtest.AssertArc(t, quarterArc().Reverse(), Ac(Pt(0.0, 0.0), 5.0, Pi/2, -Pi/2))
	})
	t.Run("the same points the other way round", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, fraction := range []float64{0, 0.25, 0.5, 1} {
				geomtest.AssertPoint(t, a.Reverse().PointAt(1-fraction), a.PointAt(fraction), fmt.Sprintf("%s @%v: ", a, fraction))
			}
		}
	})
}

func TestArc_PointAt(t *testing.T) {
	t.Run("a fraction of the sweep", func(t *testing.T) {
		geomtest.AssertPoint(t, quarterArc().PointAt(0.5), Pt(5/math.Sqrt2, 5/math.Sqrt2))
		geomtest.AssertPoint(t, Ac(Pt(1.0, 2.0), 5.0, Pi, -Pi/2).PointAt(1), Pt(1.0, 7.0))
	})
	t.Run("outside the arc it runs on around the circle", func(t *testing.T) {
		geomtest.AssertPoint(t, quarterArc().PointAt(2), Pt(-5.0, 0.0))
		geomtest.AssertPoint(t, quarterArc().PointAt(-1), Pt(0.0, -5.0))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, Ac(Pt(0, 0), 5, 0, Pi/2).PointAt(0.5), Pt(4, 4))
	})
}

func TestArc_TangentAt(t *testing.T) {
	t.Run("a quarter turn from the radius, the way the arc is drawn", func(t *testing.T) {
		geomtest.AssertVector(t, quarterArc().TangentAt(0), Vec(0.0, 2.5*Pi))
		geomtest.AssertVector(t, quarterArc().TangentAt(1), Vec(-2.5*Pi, 0.0))
		geomtest.AssertVector(t, quarterArc().Reverse().TangentAt(0), Vec(2.5*Pi, 0.0))
	})
	t.Run("as long as the arc, in the direction it moves in", func(t *testing.T) {
		const step = 1e-6

		for _, a := range arcFixtures {
			for _, fraction := range []float64{0, 0.1, 0.35, 0.5, 0.8, 1} {
				moved := a.PointAt(fraction + step).Subtract(a.PointAt(fraction - step)).Divide(2 * step)
				message := fmt.Sprintf("%s @%v: %v ≠ %v", a, fraction, a.TangentAt(fraction), moved)

				geomtest.AssertNumber(t, a.TangentAt(fraction).Length(), a.Length(), message)
				assert.True(t, a.TangentAt(fraction).Subtract(moved).LessOrEqual(1e-4*(1+a.Length())), message)
			}
		}
	})
}

func TestArc_Split(t *testing.T) {
	t.Run("two arcs meeting at the fraction", func(t *testing.T) {
		first, second := quarterArc().Split(0.5)

		geomtest.AssertArc(t, first, Ac(Pt(0.0, 0.0), 5.0, 0, Pi/4))
		geomtest.AssertArc(t, second, Ac(Pt(0.0, 0.0), 5.0, Pi/4, Pi/4))
	})
	t.Run("a sweep past a full turn is held to it", func(t *testing.T) {
		first, second := Ac(Pt(0.0, 0.0), 5.0, 0, 9).Split(0.25)

		geomtest.AssertArc(t, first, Ac(Pt(0.0, 0.0), 5.0, 0, Pi/2))
		geomtest.AssertArc(t, second, Ac(Pt(0.0, 0.0), 5.0, Pi/2, 3*Pi/2))
	})
	t.Run("the two trace the same points the same way, and are as long together", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, at := range []float64{0, 0.3, 1, 1.5} {
				if math.Abs(a.Sweep)*at > 2*Pi {
					continue
				}

				first, second := a.Split(at)

				for _, fraction := range []float64{0, 0.25, 0.5, 1} {
					message := fmt.Sprintf("%s /%v @%v: ", a, at, fraction)

					geomtest.AssertPoint(t, first.PointAt(fraction), a.PointAt(at*fraction), message)
					geomtest.AssertPoint(t, second.PointAt(fraction), a.PointAt(at+(1-at)*fraction), message)
				}
			}

			first, second := a.Split(0.3)
			geomtest.AssertNumber(t, first.Length()+second.Length(), a.Length(), a.String())
		}
	})
}

func TestArc_DistanceTo(t *testing.T) {
	a := quarterArc()

	t.Run("zero on the arc", func(t *testing.T) {
		geomtest.AssertNumber(t, a.DistanceTo(Pt(3.0, 4.0)), 0.0)
	})
	t.Run("to the circle where the sweep covers the direction", func(t *testing.T) {
		geomtest.AssertNumber(t, a.DistanceTo(Pt(6.0, 8.0)), 5.0)
		geomtest.AssertNumber(t, a.DistanceTo(Pt(1.0, 1.0)), 5-math.Sqrt2)
	})
	t.Run("to the nearer end where it does not", func(t *testing.T) {
		geomtest.AssertNumber(t, a.DistanceTo(Pt(-5.0, 0.0)), math.Sqrt(50))
		geomtest.AssertNumber(t, a.DistanceTo(Pt(0.0, -3.0)), math.Sqrt(34))
	})
	t.Run("the center is a radius away", func(t *testing.T) {
		geomtest.AssertNumber(t, a.DistanceTo(Pt(0.0, 0.0)), 5.0)
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Ac(Pt(0, 0), 5, 0, Pi/2).DistanceTo(Pt(6, 8)), 5.0)
	})
}

func TestArc_DistanceSquaredTo(t *testing.T) {
	t.Run("the square of the distance", func(t *testing.T) {
		geomtest.AssertNumber(t, quarterArc().DistanceSquaredTo(Pt(6.0, 8.0)), 25.0)
	})
	t.Run("a point within the tolerance of the arc is on it", func(t *testing.T) {
		assert.Equal(t, quarterArc().DistanceSquaredTo(Pt(3.0, 4.0+Delta/2)), 0.0)
	})
}

func TestArc_Nearest(t *testing.T) {
	a := quarterArc()

	t.Run("a point of the arc is itself", func(t *testing.T) {
		geomtest.AssertPoint(t, a.Nearest(Pt(3.0, 4.0)), Pt(3.0, 4.0))
	})
	t.Run("the point of the circle toward it", func(t *testing.T) {
		geomtest.AssertPoint(t, a.Nearest(Pt(6.0, 8.0)), Pt(3.0, 4.0))
		geomtest.AssertPoint(t, a.Nearest(Pt(0.3, 0.4)), Pt(3.0, 4.0))
	})
	t.Run("the nearer end outside the sweep", func(t *testing.T) {
		geomtest.AssertPoint(t, a.Nearest(Pt(-5.0, 0.0)), Pt(0.0, 5.0))
		geomtest.AssertPoint(t, a.Nearest(Pt(0.0, -3.0)), Pt(5.0, 0.0))
	})
	t.Run("the center is given the start", func(t *testing.T) {
		geomtest.AssertPoint(t, a.Nearest(Pt(0.0, 0.0)), Pt(5.0, 0.0))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, Ac(Pt(0, 0), 5, 0, Pi/2).Nearest(Pt(7, 7)), Pt(4, 4))
	})
	t.Run("the nearest point of a NaN coordinate is NaN", func(t *testing.T) {
		for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
			nearest := a.Nearest(p)

			assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y), p.String())
		}
	})
}

func TestArc_NearestFraction(t *testing.T) {
	a := quarterArc()

	t.Run("how far the direction to the point has swept", func(t *testing.T) {
		geomtest.AssertNumber(t, a.NearestFraction(Pt(7.0, 7.0)), 0.5)
		geomtest.AssertNumber(t, a.NearestFraction(Pt(0.0, 9.0)), 1.0)
	})
	t.Run("the nearer end outside the sweep, and the start for the center", func(t *testing.T) {
		geomtest.AssertNumber(t, a.NearestFraction(Pt(-5.0, 0.0)), 1.0)
		geomtest.AssertNumber(t, a.NearestFraction(Pt(0.0, -3.0)), 0.0)
		geomtest.AssertNumber(t, a.NearestFraction(Pt(0.0, 0.0)), 0.0)
	})
	t.Run("a sweep the other way counts from its own start", func(t *testing.T) {
		geomtest.AssertNumber(t, a.Reverse().NearestFraction(Pt(0.0, 9.0)), 0.0)
		geomtest.AssertNumber(t, a.Reverse().NearestFraction(Pt(7.0, 7.0)), 0.5)
	})
	t.Run("the point at it is the nearest point", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, point := range pointFixtures {
				fraction := a.NearestFraction(point)
				message := fmt.Sprintf("%s → %s: ", a, point)

				assert.True(t, fraction >= 0 && fraction <= 1, message)
				assert.True(t, a.PointAt(fraction).Equal(a.Nearest(point)), message)
			}
		}
	})
}

func TestArc_Contains(t *testing.T) {
	a := quarterArc()

	t.Run("a point of the arc, its ends included", func(t *testing.T) {
		assert.True(t, a.Contains(Pt(3.0, 4.0)))
		assert.True(t, a.Contains(Pt(5.0, 0.0)))
		assert.True(t, a.Contains(Pt(0.0, 5.0)))
	})
	t.Run("off the circle or past an end", func(t *testing.T) {
		assert.False(t, a.Contains(Pt(4.0, 4.0)))
		assert.False(t, a.Contains(Pt(0.0, -5.0)))
		assert.False(t, a.Contains(Pt(0.0, 0.0)))
	})
	t.Run("the tolerance is on the distance, past an end as beside the arc", func(t *testing.T) {
		assert.True(t, a.Contains(Pt(5.0, -Delta/2)))
		assert.True(t, a.Contains(Pt(5.0+Delta/2, 1.0e-3)))
		assert.False(t, a.Contains(Pt(5.0, -2*Delta)))
	})
	t.Run("int contains the lattice points of the arc, at the tolerance of float64", func(t *testing.T) {
		half := Ac(Pt(0, 0), 5, 0, Pi)

		for _, point := range []Point[int]{Pt(5, 0), Pt(4, 3), Pt(3, 4), Pt(0, 5), Pt(-3, 4), Pt(-5, 0), half.Start(), half.Midpoint(), half.End()} {
			assert.True(t, half.Contains(point), point.String())
			assert.Equal(t, half.DistanceTo(point), 0.0, point.String())
		}
		for _, point := range []Point[int]{Pt(3, 3), Pt(0, -5), Pt(0, 0), Pt(5, -1)} {
			assert.False(t, half.Contains(point), point.String())
		}
	})
	t.Run("a NaN coordinate is on no arc", func(t *testing.T) {
		assert.False(t, a.Contains(Pt(math.NaN(), 0.0)))
		assert.True(t, math.IsNaN(a.DistanceTo(Pt(0.0, math.NaN()))))
	})
}

func TestArc_IntersectionSegment(t *testing.T) {
	a := quarterArc()

	t.Run("the points of the circle that lie on the arc, from Start", func(t *testing.T) {
		geomtest.AssertVertices(t, a.IntersectionSegment(Seg(Pt(3.0, 10.0), Pt(3.0, -10.0))).AppendTo(nil), []Point[float64]{Pt(3.0, 4.0)})
		geomtest.AssertVertices(t, a.IntersectionSegment(Seg(Pt(-10.0, 3.0), Pt(10.0, 3.0))).AppendTo(nil), []Point[float64]{Pt(4.0, 3.0)})
		assert.Equal(t, a.IntersectionSegment(Seg(Pt(0.0, 6.0), Pt(6.0, 0.0))).Count, 2)
	})
	t.Run("a tangent is its point", func(t *testing.T) {
		geomtest.AssertVertices(t, Ac(Pt(0.0, 0.0), 5.0, 0, Pi).IntersectionSegment(Seg(Pt(-5.0, 5.0), Pt(5.0, 5.0))).AppendTo(nil), []Point[float64]{Pt(0.0, 5.0)})
	})
	t.Run("past the sweep, short of it or inside the circle gives none", func(t *testing.T) {
		assert.Nil(t, a.IntersectionSegment(Seg(Pt(-3.0, 10.0), Pt(-3.0, -10.0))).AppendTo(nil))
		assert.Nil(t, a.IntersectionSegment(Seg(Pt(3.0, -10.0), Pt(3.0, -1.0))).AppendTo(nil))
		assert.Nil(t, a.IntersectionSegment(Seg(Pt(1.0, 1.0), Pt(2.0, 3.0))).AppendTo(nil))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Ac(Pt(0, 0), 5, 0, Pi/2).IntersectionSegment(Seg(Pt(3, 10), Pt(3, -10))).AppendTo(nil), []Point[int]{Pt(3, 4)})
	})
	t.Run("every point lies on the segment and the arc", func(t *testing.T) {
		for _, arc := range arcFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					points := arc.IntersectionSegment(s).AppendTo(nil)
					message := fmt.Sprintf("%s → %s: ", arc, s)

					for _, point := range points {
						assert.True(t, s.Contains(point) && arc.Contains(point), message+point.String())
					}
					assertOrdered(t, s, points, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = a.IntersectionSegment(through).Count > 0
		}), 0)
	})
}

func TestArc_IntersectionRay(t *testing.T) {
	a := quarterArc()

	geomtest.AssertVertices(t, a.IntersectionRay(Ry(Pt(3.0, 10.0), Vec(0.0, -1.0))).AppendTo(nil), []Point[float64]{Pt(3.0, 4.0)})
	assert.Nil(t, a.IntersectionRay(Ry(Pt(3.0, 10.0), Vec(0.0, 1.0))).AppendTo(nil))
}

func TestArc_IntersectionLine(t *testing.T) {
	a := quarterArc()

	geomtest.AssertVertices(t, a.IntersectionLine(Ln(Pt(3.0, 10.0), Vec(0.0, 1.0))).AppendTo(nil), []Point[float64]{Pt(3.0, 4.0)})
	assert.Nil(t, a.IntersectionLine(Ln(Pt(-3.0, 10.0), Vec(0.0, 1.0))).AppendTo(nil))
}

func TestArc_Equal(t *testing.T) {
	t.Run("equal, a full turn of the angle aside", func(t *testing.T) {
		assert.True(t, Ac(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Ac(Pt(1.0, 2.0), 5.0, 0.5+2*Pi, 1.5)))
	})
	t.Run("a different radius or sweep", func(t *testing.T) {
		assert.False(t, Ac(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Ac(Pt(1.0, 2.0), 6.0, 0.5, 1.5)))
		assert.False(t, Ac(Pt(1.0, 2.0), 5.0, 0.5, 1.5).Equal(Ac(Pt(1.0, 2.0), 5.0, 0.5, 1.6)))
	})
	t.Run("the reversed arc is another value", func(t *testing.T) {
		assert.False(t, quarterArc().Equal(quarterArc().Reverse()))
	})
}

func TestArc_IsZero(t *testing.T) {
	assert.True(t, Arc[int]{}.IsZero())
	assert.False(t, Ac(Pt(0, 0), 0, 0, 1).IsZero())
}

func TestArc_Circle(t *testing.T) {
	geomtest.AssertCircle(t, Ac(Pt(1, 2), 5, 0.5, 1.5).Circle(), Circ(Pt(1, 2), 5))
}

func TestArc_Polyline(t *testing.T) {
	t.Run("equal chords from the start to the end", func(t *testing.T) {
		geomtest.AssertPolyline(t, quarterArc().Polyline(0.5), PolLn([]Point[float64]{Pt(5.0, 0.0), Pt(5/math.Sqrt2, 5/math.Sqrt2), Pt(0.0, 5.0)}))
	})
	t.Run("a finer tolerance gives more points, either sign", func(t *testing.T) {
		assert.True(t, len(quarterArc().Polyline(0.01).Points) > len(quarterArc().Polyline(0.5).Points))
		assert.Equal(t, len(quarterArc().Polyline(-0.5).Points), 3)
	})
	t.Run("a tolerance that cannot be met stops at a thousand edges", func(t *testing.T) {
		assert.Length(t, quarterArc().Polyline(0).Points, 1025)
		assert.Length(t, quarterArc().Polyline(math.NaN()).Points, 1025)
	})
	t.Run("a tolerance past the arc gives its chord, and an arc of no length its point", func(t *testing.T) {
		geomtest.AssertPolyline(t, quarterArc().Polyline(100), PolLn([]Point[float64]{Pt(5.0, 0.0), Pt(0.0, 5.0)}))
		geomtest.AssertPolyline(t, Ac(Pt(1.0, 2.0), 5.0, 0, 0).Polyline(0.1), PolLn([]Point[float64]{Pt(6.0, 2.0)}))
		geomtest.AssertPolyline(t, Ac(Pt(1.0, 2.0), 0.0, 0, 3).Polyline(0), PolLn([]Point[float64]{Pt(1.0, 2.0)}))
	})
	t.Run("int leaves out a point rounding onto the one before", func(t *testing.T) {
		path := Ac(Pt(0, 0), 2, 0, Pi).Polyline(0.01)

		geomtest.AssertPoint(t, path.Points[0], Pt(2, 0))
		geomtest.AssertPoint(t, path.Points[len(path.Points)-1], Pt(-2, 0))
		for i := 1; i < len(path.Points); i++ {
			assert.True(t, path.Points[i] != path.Points[i-1], path.String())
		}
	})
	t.Run("the path lies on the arc and the arc within the tolerance of it", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, tolerance := range []float64{0.5, 0.05, 0.001} {
				path := a.Polyline(tolerance)
				message := fmt.Sprintf("%s ±%v: ", a, tolerance)

				assert.True(t, path.Points[0].Equal(a.Start()) && path.Points[len(path.Points)-1].Equal(a.End()), message)
				for _, point := range path.Points {
					assert.True(t, a.Contains(point), message+point.String())
				}
				for i := range 257 {
					point := a.PointAt(float64(i) / 256)

					assert.True(t, LessOrEqual(path.DistanceTo(point), tolerance), message+point.String())
				}
			}
		}
	})
}

func TestArc_AppendPolyline(t *testing.T) {
	t.Run("appends after the points already there", func(t *testing.T) {
		points := quarterArc().AppendPolyline([]Point[float64]{Pt(9.0, 9.0)}, 0.5)

		geomtest.AssertPoint(t, points[0], Pt(9.0, 9.0))
		geomtest.AssertVertices(t, points[1:], quarterArc().Polyline(0.5).Points)
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		arc, buffer := Ac(Pt(0, 0), 50, 0, Pi), make([]Point[int], 0, 64)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = arc.Polyline(tenth).Points
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = arc.AppendPolyline(buffer[:0], 1)
		}), 0)
	})
}

func TestArc_Sector(t *testing.T) {
	geomtest.AssertSector(t, Ac(Pt(1, 2), 5, 0.5, 1.5).Sector(), Sect(Pt(1, 2), 5, 0.5, 1.5))
}

func TestArc_Cast(t *testing.T) {
	a := Ac(Pt(1.5, -2.5), 3.5, 0.5, 1.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertArc(t, a.Cast[int](), a.Int())
		geomtest.AssertArc(t, a.Cast[float64](), a.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertArc(t, a.Cast[int8](), Ac(Pt[int8](2, -3), int8(4), 0.5, 1.5))
	})
}

func TestArc_Int(t *testing.T) {
	geomtest.AssertArc(t, Ac(Pt(0.6, -0.25), 1.2, 0.5, 1.5).Int(), Ac(Pt(1, 0), 1, 0.5, 1.5))
}

func TestArc_Float(t *testing.T) {
	geomtest.AssertArc(t, Ac(Pt(1, 2), 5, 0.5, 1.5).Float(), Ac(Pt(1.0, 2.0), 5.0, 0.5, 1.5))
}

func TestArc_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Ac(Pt(1, 2), 5, 0.5, -1.5).String(), "Arc((1,2);5;0.50;-1.50)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Ac(Pt(1.5, 2), 0.25, 0, 3).String(), "Arc((1.50,2.00);0.25;0.00;3.00)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Ac(Pt(1, 2), 5, 0.5, -1.5).String)
	})
}

func TestArc_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Ac(Pt(1, 2), 5, 0.5, -1.5), `{"x":1,"y":2,"r":5,"a":0.5,"s":-1.5}`)

		var a Arc[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"r":5,"a":0.5,"s":-1.5}`), &a))
		geomtest.AssertArc(t, a, Ac(Pt(1, 2), 5, 0.5, -1.5))
	})
	t.Run("an angle of zero is left out", func(t *testing.T) {
		assert.JSON(t, Ac(Pt(1, 2), 5, 0, 1.5), `{"x":1,"y":2,"r":5,"s":1.5}`)
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, arc := range arcFixtures {
			data, err := json.Marshal(arc)
			assert.NoError(t, err)

			var decoded Arc[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, arc)
		}
	})
}

func TestArc_Properties(t *testing.T) {
	fractions := []float64{0, 0.1, 0.35, 0.5, 0.8, 1}

	t.Run("a point at a fraction lies on the arc, on its circle and in its bounds", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, fraction := range fractions {
				point := a.PointAt(fraction)
				message := fmt.Sprintf("%s @%v: ", a, fraction)

				assert.True(t, a.Contains(point), message)
				assert.True(t, a.Bounds().Contains(point), message)
				geomtest.AssertNumber(t, a.Center.DistanceTo(point), a.Radius, message)
			}
		}
	})
	t.Run("the length is that of the path through many of its points", func(t *testing.T) {
		for _, a := range arcFixtures {
			points := make([]Point[float64], 2049)
			for i := range points {
				points[i] = a.PointAt(float64(i) / 2048)
			}

			assert.True(t, math.Abs(a.Length()-PolLn(points).Length()) <= 1e-5*(1+a.Length()), a.String())
		}
	})
	t.Run("the nearest point is on the arc and at the distance, no point of it nearer", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, point := range pointFixtures {
				nearest := a.Nearest(point)
				message := fmt.Sprintf("%s → %s: ", a, point)

				assert.Equal(t, a.DistanceTo(point) == 0, a.Contains(point), message)
				assert.Equal(t, nearest.Equal(point), a.Contains(point), message)
				assert.True(t, a.Contains(nearest), message)
				geomtest.AssertNumber(t, point.DistanceTo(nearest), a.DistanceTo(point), message)

				for i := range 257 {
					assert.True(t, LessOrEqual(a.DistanceTo(point), point.DistanceTo(a.PointAt(float64(i)/256))), message)
				}
			}
		}
	})
	t.Run("the whole circle is as far as the boundary of its circle", func(t *testing.T) {
		for _, c := range circleFixtures {
			a := Ac(c.Center, c.Radius, 1, -2*Pi)

			for _, point := range pointFixtures {
				geomtest.AssertNumber(t, a.DistanceTo(point), math.Abs(c.Center.DistanceTo(point)-c.Radius), fmt.Sprintf("%s → %s: ", a, point))
			}
		}
	})
	t.Run("a reversed, turned or moved arc is as long and as far", func(t *testing.T) {
		for _, a := range arcFixtures {
			geomtest.AssertNumber(t, a.Reverse().Length(), a.Length(), a.String())

			for _, point := range pointFixtures {
				message := fmt.Sprintf("%s → %s: ", a, point)

				geomtest.AssertNumber(t, a.Reverse().DistanceTo(point), a.DistanceTo(point), message)
				geomtest.AssertNumber(t, a.Rotate(Pi/3).DistanceTo(point.RotateAround(a.Center, Pi/3)), a.DistanceTo(point), message)
				geomtest.AssertNumber(t, a.Translate(Vec(3.0, -2.0)).DistanceTo(point.Add(Vec(3.0, -2.0))), a.DistanceTo(point), message)
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, factor := range []float64{0.5, 2.5, -3} {
				assert.True(t, a.Scale(factor).Unscale(factor).Equal(a), fmt.Sprintf("%s ×%v: ", a, factor))
			}
		}
	})
	t.Run("far from the origin a float32 point at a fraction lies on the arc", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, fixture := range arcFixtures {
				a := fixture.Cast[float32]().Translate(offset)

				for _, fraction := range fractions {
					assert.True(t, a.Contains(a.PointAt(fraction)), fmt.Sprintf("%s @%v: ", a, fraction))
				}
			}
		}
	})
}

func TestArc_Immutable(t *testing.T) {
	a := Ac(Pt(1, 2), 5, 0.5, 1.5)

	a.Translate(Vec(3, -2))
	a.MoveTo(Pt(4, 3))
	a.Rotate(1)
	a.Scale(2)
	a.Resize(9)
	a.Grow(1)
	a.Shrink(1)
	a.Reverse()
	a.Split(0.5)

	geomtest.AssertArc(t, a, Ac(Pt(1, 2), 5, 0.5, 1.5))
}

// quarterArc returns the arc of radius five about the origin from the X axis a quarter turn on.
func quarterArc() Arc[float64] {
	return Ac(Pt(0.0, 0.0), 5.0, 0, Pi/2)
}

// arcFixtures span a point, a quarter, a half and more than a half turn either way, a sweep of
// zero and the whole circle, on and off the origin.
var arcFixtures = []Arc[float64]{
	Ac(Pt(0.0, 0.0), 0.0, 0, 0),
	Ac(Pt(0.0, 0.0), 1.0, 0, Pi/2),
	Ac(Pt(1.0, 2.0), 10.0, Pi/6, Pi),
	Ac(Pt(-3.5, 0.25), 7.0, 5, -4),
	Ac(Pt(0.6, -0.25), 1.2, 1, 2*Pi),
	Ac(Pt(12.5, -0.1), 0.5, 3, 0),
	Ac(Pt(1.0, 2.0), 3.0, 0, -Pi/3),
}

func ExampleAc() {
	fmt.Println(Ac(Pt(1, 2), 5, 0, 1.5))
	// Output: Arc((1,2);5;0.00;1.50)
}

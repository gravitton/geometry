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

func TestLine_Constructor(t *testing.T) {
	t.Run("shorthand", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, -1), Vec(3, 4)), Line[int]{Origin: Pt(1, -1), Vector: Vec(3, 4)})
		geomtest.AssertLine(t, Ln(Pt(0.5, -1.25), Vec(2.0, 5.0)), Line[float64]{Origin: Pt(0.5, -1.25), Vector: Vec(2.0, 5.0)})
	})
	t.Run("through two points", func(t *testing.T) {
		geomtest.AssertLine(t, LineThrough(Pt(1, -1), Pt(4, 3)), Ln(Pt(1, -1), Vec(3, 4)))
		geomtest.AssertLine(t, LineThrough(Pt(0.5, -1.25), Pt(2.5, 3.75)), Ln(Pt(0.5, -1.25), Vec(2.0, 5.0)))
	})
	t.Run("through one point twice runs nowhere", func(t *testing.T) {
		geomtest.AssertLine(t, LineThrough(Pt(2, 3), Pt(2, 3)), Line[int]{Origin: Pt(2, 3)})
	})
}

func TestLine_Angle(t *testing.T) {
	t.Run("the angle of the direction", func(t *testing.T) {
		geomtest.AssertNumber(t, Ln(Pt(1, 1), Vec(0, 2)).Angle(), Pi/2)
		geomtest.AssertNumber(t, Ln(Pt(1.0, 1.0), Vec(-3.0, 0.0)).Angle(), Pi)
	})
	t.Run("a zero direction gives zero", func(t *testing.T) {
		geomtest.AssertNumber(t, Ln(Pt(1, 1), Vec(0, 0)).Angle(), 0)
	})
}

func TestLine_Direction(t *testing.T) {
	t.Run("the nearest direction of the vector", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(1, 1), Vec(5, 0)).Direction(), DirectionRight)
		assert.Equal(t, Ln(Pt(1, 1), Vec(0, -5)).Direction(), DirectionUp)
		assert.Equal(t, Ln(Pt(1.0, 1.0), Vec(4.0, 5.0)).Direction(), DirectionDownRight)
	})
	t.Run("a zero vector has no direction", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(3, 4), Vec(0, 0)).Direction(), DirectionNone)
	})
	t.Run("is the direction of the vector over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			assert.Equal(t, l.Direction(), l.Vector.Direction(), l.String())
		}
	})
}

func TestLine_Normal(t *testing.T) {
	t.Run("a quarter turn of the vector, with its length", func(t *testing.T) {
		geomtest.AssertVector(t, Ln(Pt(1, 1), Vec(3, 0)).Normal(), Vec(0, 3))
		geomtest.AssertVector(t, Ln(Pt(0.0, 0.0), Vec(0.0, 2.5)).Normal(), Vec(-2.5, 0.0))
	})
	t.Run("a zero vector has none", func(t *testing.T) {
		geomtest.AssertVector(t, Ln(Pt(3, 4), Vec(0, 0)).Normal(), ZeroVector[int]())
	})
	t.Run("is perpendicular to the vector over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			geomtest.AssertNumber(t, l.Normal().Dot(l.Vector), 0.0, l.String())
			geomtest.AssertNumber(t, l.Normal().Length(), l.Vector.Length(), l.String())
		}
	})
	t.Run("points to the side Side gives 1 for", func(t *testing.T) {
		for _, l := range lineFixtures[1:] {
			assert.Equal(t, l.Side(l.Origin.Add(l.Normal())), 1, l.String())
		}
	})
}

func TestLine_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Translate(Vec(-2, 5)), Ln(Pt(-1, 7), Vec(3, 4)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(0.5, 2.0), Vec(3.0, 4.0)).Translate(Vec(0.25, -1.5)), Ln(Pt(0.75, 0.5), Vec(3.0, 4.0)))
	})
}

func TestLine_MoveTo(t *testing.T) {
	t.Run("passes through the point with the same direction", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).MoveTo(Pt(-5, 0)), Ln(Pt(-5, 0), Vec(3, 4)))
	})
}

func TestLine_Rotate(t *testing.T) {
	t.Run("quarter turn about the origin", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 1), Vec(2, 0)).Rotate(Pi/2), Ln(Pt(1, 1), Vec(0, 2)))
	})
	t.Run("a full turn is identity", func(t *testing.T) {
		for _, l := range lineFixtures {
			geomtest.AssertLine(t, l.Rotate(2*Pi), l, l.String()+": ")
		}
	})
}

func TestLine_Scale(t *testing.T) {
	t.Run("scales the direction about the origin", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Scale(2), Ln(Pt(1, 2), Vec(6, 8)))
		geomtest.AssertLine(t, Ln(Pt(1.0, 2.0), Vec(3.0, 4.0)).Scale(0.5), Ln(Pt(1.0, 2.0), Vec(1.5, 2.0)))
	})
	t.Run("per-axis factor changes the direction", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).ScaleXY(2, -1), Ln(Pt(1, 2), Vec(6, -4)))
	})
	t.Run("a negative factor turns the direction about", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Scale(-1), Ln(Pt(1, 2), Vec(-3, -4)))
	})
	t.Run("zero collapses onto the origin", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Scale(0), Ln(Pt(1, 2), Vec(0, 0)))
	})
	t.Run("int rounds the direction", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Scale(0.5), Ln(Pt(1, 2), Vec(2, 2)))
	})
	t.Run("keeps the points and steps the scaled length", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, factor := range []float64{0.5, -2.5} {
				scaled := l.Scale(factor)

				assert.True(t, l.Contains(scaled.PointAt(3)), fmt.Sprintf("%s x%v: ", l, factor))
				geomtest.AssertNumber(t, scaled.Vector.Length(), l.Vector.Length()*math.Abs(factor), fmt.Sprintf("%s x%v: ", l, factor))
			}
		}
	})
}

func TestLine_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(6, 8)).Unscale(2), Ln(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(6, -4)).UnscaleXY(2, -1), Ln(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("undoes scale", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				geomtest.AssertLine(t, l.Scale(factor).Unscale(factor), l, fmt.Sprintf("%s x%v: ", l, factor))
			}
		}
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			Ln(Pt(1, 2), Vec(3, 4)).Unscale(0)
		}, ErrDivisionByZero)
		assert.PanicsWith(t, func() {
			Ln(Pt(1, 2), Vec(3, 4)).UnscaleXY(1, 0)
		}, ErrDivisionByZero)
	})
}

func TestLine_Transform(t *testing.T) {
	t.Run("a translation moves the origin alone", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Transform(Mat(1.0, 0.0, 5.0, 0.0, 1.0, -1.0)), Ln(Pt(6, 1), Vec(3, 4)))
	})
	t.Run("a scale applies to both", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Transform(Mat(2.0, 0.0, 0.0, 0.0, -1.0, 0.0)), Ln(Pt(2, -2), Vec(6, -4)))
	})
	t.Run("keeps every point on the transformed line", func(t *testing.T) {
		matrix := IdentityMatrix[float64]().Rotate(Pi/5).Translate(3, -2)
		for _, l := range lineFixtures {
			transformed := l.Transform(matrix)

			for _, t0 := range []float64{-3, 0, 0.5, 3} {
				assert.True(t, transformed.Contains(l.PointAt(t0).Transform(matrix)), fmt.Sprintf("%s at %v: ", l, t0))
			}
		}
	})
}

func TestLine_Lerp(t *testing.T) {
	a := Ln(Pt(0.0, 0.0), Vec(4.0, 0.0))
	b := Ln(Pt(10.0, 2.0), Vec(0.0, 8.0))

	t.Run("the origin and the vector move together", func(t *testing.T) {
		geomtest.AssertLine(t, a.Lerp(b, 0.5), Ln(Pt(5.0, 1.0), Vec(2.0, 4.0)))
	})
	t.Run("the ends are the two lines", func(t *testing.T) {
		geomtest.AssertLine(t, a.Lerp(b, 0), a)
		geomtest.AssertLine(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates outside [0, 1]", func(t *testing.T) {
		geomtest.AssertLine(t, a.Lerp(b, 2), Ln(Pt(20.0, 4.0), Vec(-4.0, 16.0)))
	})
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(0, 0), Vec(4, 0)).Lerp(Ln(Pt(1, 1), Vec(1, 3)), 0.5), Ln(Pt(1, 1), Vec(3, 2)))
	})
	t.Run("towards a translation is a part of it over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			geomtest.AssertLine(t, l.Lerp(l.Translate(Vec(4.0, -2.0)), 0.25), l.Translate(Vec(1.0, -0.5)), l.String())
		}
	})
}

func TestLine_Reverse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, -4)).Reverse(), Ln(Pt(1, 2), Vec(-3, 4)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(0.6, -0.25), Vec(1.2, 3.4)).Reverse(), Ln(Pt(0.6, -0.25), Vec(-1.2, -3.4)))
	})
	t.Run("holds the same points and swaps the sides over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				message := fmt.Sprintf("%s → %s: ", l, p)

				assert.Equal(t, l.Reverse().Contains(p), l.Contains(p), message)
				assert.Equal(t, l.Reverse().Side(p), -l.Side(p), message)
			}
		}
	})
}

func TestLine_PointAt(t *testing.T) {
	line := Ln(Pt(1.0, 2.0), Vec(2.0, -1.0))

	t.Run("the origin at 0 and one direction away at 1", func(t *testing.T) {
		geomtest.AssertPoint(t, line.PointAt(0), Pt(1.0, 2.0))
		geomtest.AssertPoint(t, line.PointAt(1), Pt(3.0, 1.0))
	})
	t.Run("along the line", func(t *testing.T) {
		geomtest.AssertPoint(t, line.PointAt(2.5), Pt(6.0, -0.5))
	})
	t.Run("on the other side of the origin for a negative fraction", func(t *testing.T) {
		geomtest.AssertPoint(t, line.PointAt(-1), Pt(-1.0, 3.0))
	})
	t.Run("int rounds once, on the sum", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(1, 0), Vec(-1, 0)).PointAt(0.5), Pt(1, 0))
		geomtest.AssertPoint(t, Ln(Pt(1, 1), Vec(1, 1)).PointAt(0.5), Pt(2, 2))
	})
}

func TestLine_Side(t *testing.T) {
	line := Ln(Pt(1, 1), Vec(2, 0))

	t.Run("positive on the side the normal points to", func(t *testing.T) {
		assert.Equal(t, line.Side(Pt(0, 4)), 1)
		assert.Equal(t, line.Side(line.Origin.Add(line.Vector.Normal())), 1)
	})
	t.Run("negative on the other", func(t *testing.T) {
		assert.Equal(t, line.Side(Pt(9, -4)), -1)
	})
	t.Run("zero on the line, either way from the origin", func(t *testing.T) {
		assert.Equal(t, line.Side(Pt(7, 1)), 0)
		assert.Equal(t, line.Side(Pt(-7, 1)), 0)
	})
	t.Run("exact where contains is tolerant", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.True(t, l.Contains(Pt(3.0, Delta/2)))
		assert.Equal(t, l.Side(Pt(3.0, Delta/2)), 1)
		assert.Equal(t, l.Side(Pt(3.0, -Delta/2)), -1)
	})
	t.Run("a zero direction has no sides", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(1, 1), Vec(0, 0)).Side(Pt(4, 5)), 0)
	})
	t.Run("a NaN coordinate has no sign", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(0.0, 0.0), Vec(1.0, 0.0)).Side(Pt(0.0, math.NaN())), 0)
	})
	t.Run("the reversed line swaps the sides", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, l.Scale(-1).Side(p), -l.Side(p), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
	t.Run("the sign of the distance along the normal", func(t *testing.T) {
		for _, l := range lineFixtures {
			normal := l.Vector.Normal()

			for _, p := range pointFixtures {
				along := p.Subtract(l.Origin).Dot(normal)

				assert.Equal(t, l.Side(p), int(Sign(along)), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_DistanceTo(t *testing.T) {
	line := Ln(Pt(0, 0), Vec(1, 0))

	t.Run("the length of the perpendicular, either way from the origin", func(t *testing.T) {
		geomtest.AssertNumber(t, line.DistanceTo(Pt(100, 3)), 3)
		geomtest.AssertNumber(t, line.DistanceTo(Pt(1, -4)), 4)
		geomtest.AssertNumber(t, line.DistanceTo(Pt(-3, 4)), 4)
	})
	t.Run("on the line is exactly zero", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(1, 1), Vec(3, 7)).DistanceTo(Pt(7, 15)), 0.0)
		assert.Equal(t, Ln(Pt(1, 1), Vec(3, 7)).DistanceTo(Pt(-5, -13)), 0.0)
	})
	t.Run("a zero direction measures to the origin", func(t *testing.T) {
		geomtest.AssertNumber(t, Ln(Pt(1, 1), Vec(0, 0)).DistanceTo(Pt(4, 5)), 5)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.Equal(t, l.DistanceTo(Pt(-3.0, Delta/2)), 0.0)
		geomtest.AssertNumber(t, l.DistanceTo(Pt(-3.0, 2*Delta)), 2*Delta)
	})
}

func TestLine_DistanceSquaredTo(t *testing.T) {
	line := Ln(Pt(0, 0), Vec(2, 1))

	t.Run("stays fractional for an integer T", func(t *testing.T) {
		geomtest.AssertNumber(t, line.DistanceSquaredTo(Pt(0, 1)), 0.8)
		geomtest.AssertNumber(t, line.DistanceSquaredTo(Pt(0, -1)), 0.8)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, l.DistanceSquaredTo(p), l.DistanceTo(p)*l.DistanceTo(p), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
	t.Run("the same from the reversed line", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, l.Scale(-1).DistanceSquaredTo(p), l.DistanceSquaredTo(p), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_Nearest(t *testing.T) {
	line := Ln(Pt(0.0, 0.0), Vec(2.0, 0.0))

	t.Run("the foot of the perpendicular, either way from the origin", func(t *testing.T) {
		geomtest.AssertPoint(t, line.Nearest(Pt(7.0, 3.0)), Pt(7.0, 0.0))
		geomtest.AssertPoint(t, line.Nearest(Pt(-2.0, 3.0)), Pt(-2.0, 0.0))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		geomtest.AssertPoint(t, line.Nearest(Pt(-3.0, Delta/2)), Pt(-3.0, Delta/2))
	})
	t.Run("int rounds once and can land off the line", func(t *testing.T) {
		l := Ln(Pt(0, 0), Vec(2, 1))

		geomtest.AssertPoint(t, l.Nearest(Pt(3, 3)), Pt(4, 2))
		geomtest.AssertPoint(t, l.Nearest(Pt(-3, -3)), Pt(-4, -2))
		geomtest.AssertPoint(t, l.Nearest(Pt(0, 2)), Pt(1, 0))
		assert.False(t, l.Contains(Pt(1, 0)))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(1, 1), Vec(0, 0)).Nearest(Pt(4, 5)), Pt(1, 1))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				nearest := l.Nearest(p)
				message := fmt.Sprintf("%s → %s: ", l, p)

				assert.Equal(t, nearest.Equal(p), l.Contains(p), message)
				assert.True(t, l.Contains(nearest), message)
				assert.True(t, l.Nearest(nearest).Equal(nearest), message)
				geomtest.AssertNumber(t, p.DistanceSquaredTo(nearest), l.DistanceSquaredTo(p), message)
			}
		}
	})
	t.Run("far from the origin the float32 nearest point stays on the line", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, l := range lineFixtures {
				line := l.Cast[float32]().Translate(offset)

				for _, p := range pointFixtures {
					point := p.Cast[float32]().Add(offset)

					assert.True(t, line.Contains(line.Nearest(point)), fmt.Sprintf("%s → %s: ", line, point))
				}
			}
		}
	})
	t.Run("the nearest point of a NaN coordinate is NaN", func(t *testing.T) {
		nearest := line.Nearest(Pt(math.NaN(), 0.0))

		assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y))
	})
}

func TestLine_NearestFraction(t *testing.T) {
	line := Ln(Pt(0.0, 0.0), Vec(2.0, 0.0))

	t.Run("in units of the direction", func(t *testing.T) {
		geomtest.AssertNumber(t, line.NearestFraction(Pt(7.0, 3.0)), 3.5)
	})
	t.Run("negative on the other side of the origin", func(t *testing.T) {
		geomtest.AssertNumber(t, line.NearestFraction(Pt(-2.0, 3.0)), -1.0)
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Ln(Pt(0, 0), Vec(2, 1)).NearestFraction(Pt(3, 3)), 1.8)
	})
	t.Run("a zero direction gives zero", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(1, 1), Vec(0, 0)).NearestFraction(Pt(4, 5)), 0.0)
	})
	t.Run("a NaN coordinate is nearest to nothing", func(t *testing.T) {
		assert.True(t, math.IsNaN(line.NearestFraction(Pt(math.NaN(), 0.0))))
	})
	t.Run("PointAt takes it back to the nearest point over the fixtures", func(t *testing.T) {
		for _, fixture := range lineFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertPoint(t, fixture.PointAt(fixture.NearestFraction(p)), fixture.Nearest(p), fmt.Sprintf("%s → %s: ", fixture, p))
			}
		}
	})
}

func TestLine_Contains(t *testing.T) {
	line := Ln(Pt(0, 0), Vec(2, 1))

	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, line.Contains(Pt(0, 0)))
		assert.True(t, line.Contains(Pt(4, 2)))
		assert.True(t, line.Contains(Pt(200, 100)))
		assert.True(t, line.Contains(Pt(-2, -1)))
		assert.False(t, line.Contains(Pt(1, 1)))
		assert.False(t, line.Contains(Pt(-1, -1)))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.True(t, l.Contains(Pt(5.0, Delta/2)))
		assert.True(t, l.Contains(Pt(-5.0, -Delta/2)))
		assert.False(t, l.Contains(Pt(5.0, 2*Delta)))
		assert.False(t, l.Contains(Pt(-5.0, 2*Delta)))
	})
	t.Run("a zero direction contains its origin alone", func(t *testing.T) {
		l := Ln(Pt(1, 1), Vec(0, 0))

		assert.True(t, l.Contains(Pt(1, 1)))
		assert.False(t, l.Contains(Pt(2, 1)))
	})
	t.Run("contains every point a ray along it does, either way", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range pointFixtures {
				message := fmt.Sprintf("%s → %s: ", r, p)

				assert.True(t, !r.Contains(p) || r.Line().Contains(p), message)
				assert.True(t, !r.Scale(-1).Contains(p) || r.Line().Contains(p), message)
			}
		}
	})
	t.Run("holds exactly where DistanceTo is zero", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, l.Contains(p), l.DistanceTo(p) == 0, fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
	t.Run("a NaN coordinate is contained by nothing", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
				assert.False(t, l.Contains(p), fmt.Sprintf("%s → %s: ", l, p))
				assert.True(t, math.IsNaN(l.DistanceTo(p)), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, l.IntersectsCircle(c), c.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, c))
			}
		}
	})
}

func TestLine_IntersectsCapsule(t *testing.T) {
	t.Run("mirrors Capsule.IntersectsLine", func(t *testing.T) {
		for _, shape := range lineFixtures {
			for _, c := range capsuleFixtures {
				assert.Equal(t, shape.IntersectsCapsule(c), c.IntersectsLine(shape), fmt.Sprintf("%s → %s: ", shape, c))
			}
		}
	})
}

func TestLine_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, l.IntersectsSegment(s), s.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, s))
			}
		}
	})
}

func TestLine_IntersectsRay(t *testing.T) {
	t.Run("mirrors Ray.IntersectsLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, l.IntersectsRay(r), r.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, r))
			}
		}
	})
}

func TestLine_IntersectsLine(t *testing.T) {
	line := Ln(Pt(0, 0), Vec(2, 1))

	t.Run("lines that are not parallel cross somewhere", func(t *testing.T) {
		assert.True(t, line.IntersectsLine(Ln(Pt(100, -7), Vec(1, 1))))
	})
	t.Run("parallel lines apart share nothing", func(t *testing.T) {
		assert.False(t, line.IntersectsLine(Ln(Pt(0, 3), Vec(-4, -2))))
	})
	t.Run("lines along each other share every point", func(t *testing.T) {
		assert.True(t, line.IntersectsLine(Ln(Pt(4, 2), Vec(-2, -1))))
		assert.True(t, line.IntersectsLine(line))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, line.IntersectsLine(Ln(Pt(4, 2), Vec(0, 0))))
		assert.True(t, Ln(Pt(4, 2), Vec(0, 0)).IntersectsLine(line))
		assert.False(t, line.IntersectsLine(Ln(Pt(4, 3), Vec(0, 0))))
	})
	t.Run("either side gives the same answer, and there is a point where the lines are not parallel", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				message := fmt.Sprintf("%s × %s: ", a, b)
				_, ok := a.IntersectionLine(b)

				assert.Equal(t, a.IntersectsLine(b), b.IntersectsLine(a), message)
				assert.True(t, !ok || a.IntersectsLine(b), message)
				if a.Vector.Cross(b.Vector) != 0 {
					assert.True(t, ok, message)
				}
			}
		}
	})
}

func TestLine_IntersectsPolygon(t *testing.T) {
	square := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(0, 4)})

	t.Run("passing through, whichever way it runs", func(t *testing.T) {
		assert.True(t, Ln(Pt(-2, 2), Vec(1, 0)).IntersectsPolygon(square))
		assert.True(t, Ln(Pt(-2, 2), Vec(-1, 0)).IntersectsPolygon(square))
	})
	t.Run("an origin inside", func(t *testing.T) {
		assert.True(t, Ln(Pt(2, 2), Vec(-1, 3)).IntersectsPolygon(square))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, Ln(Pt(-2, 2), Vec(0, 1)).IntersectsPolygon(square))
	})
	t.Run("touching a vertex counts", func(t *testing.T) {
		assert.True(t, Ln(Pt(2, 6), Vec(-1, 1)).IntersectsPolygon(square))
	})
	t.Run("reaches a polygon however far", func(t *testing.T) {
		assert.True(t, Ln(Pt(-2, 2), Vec(-1, 0)).IntersectsPolygon(square.Translate(Vec(100000, 0))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Ln(Pt(0, 0), Vec(1, 0)).IntersectsPolygon(Polygon[int]{}))
	})
	t.Run("matches a segment reaching past the polygon along the line both ways", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, l.IntersectsPolygon(p), wide(l).IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-2, 2), Vec(1, 0))

		assertAllocations(t, func() {
			sinkBool = through.IntersectsPolygon(square)
		}, 0)
	})
}

func TestLine_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through, whichever way it runs", func(t *testing.T) {
		assert.True(t, Ln(Pt(-5.0, 0.0), Vec(1.0, 0.0)).IntersectsRectangle(rectangle))
		assert.True(t, Ln(Pt(-5.0, 0.0), Vec(-1.0, 0.0)).IntersectsRectangle(rectangle))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, Ln(Pt(-5.0, 0.0), Vec(0.0, 1.0)).IntersectsRectangle(rectangle))
	})
	t.Run("a rotated rectangle is met on its turned edges", func(t *testing.T) {
		turned := rectangle.Rotate(Pi / 4)

		assert.True(t, Ln(Pt(5.0, 1.6), Vec(1.0, 0.0)).IntersectsRectangle(turned))
		assert.False(t, Ln(Pt(5.0, 1.6), Vec(1.0, 0.0)).IntersectsRectangle(rectangle))
	})
	t.Run("matches a segment reaching past the rectangle along the line both ways", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, rect := range rectFixtures {
				assert.Equal(t, l.IntersectsRectangle(rect), wide(l).IntersectsRectangle(rect), fmt.Sprintf("%s → %s: ", l, rect))
			}
		}
	})
}

func TestLine_IntersectsRegularPolygon(t *testing.T) {
	hexagon := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Ln(Pt(20.0, 0.0), Vec(1.0, 0.0)).IntersectsRegularPolygon(hexagon))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, Ln(Pt(-20.0, 0.0), Vec(0.0, 1.0)).IntersectsRegularPolygon(hexagon))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Ln(Pt(0.0, 0.0), Vec(1.0, 0.0)).IntersectsRegularPolygon(RegularPolygon[float64]{}))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, l.IntersectsRegularPolygon(rp), l.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", l, rp))
			}
		}
	})
}

func TestLine_IntersectsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("passing through, whichever way it runs", func(t *testing.T) {
		assert.True(t, Ln(Pt(-1, 1), Vec(1, 0)).IntersectsBox(box))
		assert.True(t, Ln(Pt(-1, 1), Vec(-1, 0)).IntersectsBox(box))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, Ln(Pt(-1, 1), Vec(0, 1)).IntersectsBox(box))
	})
	t.Run("matches IntersectsRectangle on the box's Rectangle", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, l.IntersectsBox(b), l.IntersectsRectangle(b.Rectangle()), fmt.Sprintf("%s → %s: ", l, b))
			}
		}
	})
}

func TestLine_IntersectsEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through, whichever way it runs", func(t *testing.T) {
		assert.True(t, Ln(Pt(-8.0, 0.0), Vec(1.0, 0.0)).IntersectsEllipse(ellipse))
		assert.True(t, Ln(Pt(-8.0, 0.0), Vec(-1.0, 0.0)).IntersectsEllipse(ellipse))
	})
	t.Run("passing by", func(t *testing.T) {
		assert.False(t, Ln(Pt(-8.0, 3.0), Vec(1.0, 0.0)).IntersectsEllipse(ellipse))
	})
	t.Run("matches IntersectsCircle on the circle's Ellipse", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, l.IntersectsEllipse(c.Ellipse()), l.IntersectsCircle(c), fmt.Sprintf("%s → %s: ", l, c))
			}
		}
	})
}

func TestLine_IntersectionCircle(t *testing.T) {
	t.Run("matches Circle.IntersectionLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, l.IntersectionCircle(c).AppendTo(nil), c.IntersectionLine(l).AppendTo(nil), fmt.Sprintf("%s → %s: ", l, c))
			}
		}
	})
}

func TestLine_IntersectionCapsule(t *testing.T) {
	t.Run("matches Capsule.IntersectionLine", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, c := range capsuleFixtures {
				geomtest.AssertVertices(t, s.IntersectionCapsule(c).AppendTo(nil), c.IntersectionLine(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestLine_IntersectionSector(t *testing.T) {
	t.Run("matches Sector.IntersectionLine", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, sector := range sectorFixtures() {
				geomtest.AssertVertices(t, s.IntersectionSector(sector), sector.IntersectionLine(s), fmt.Sprintf("%s → %s: ", s, sector))
			}
		}
	})
}

func TestLine_AppendIntersectionSector(t *testing.T) {
	sector, through := quarterSector(), LineThrough(Pt(-1.0, 1.0), Pt(5.0, 1.0+tenth))

	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionSector(bufferWith(prefixPoint), sector), append([]Point[float64]{prefixPoint}, through.IntersectionSector(sector)...))
	})
	t.Run("appends to a buffer with room without allocating", func(t *testing.T) {
		buffer := make([]Point[float64], 0, 4)

		assertAllocations(t, func() {
			sinkBool = len(through.AppendIntersectionSector(buffer[:0], sector)) > 0
		}, 0)
	})
}

func TestLine_IntersectionArc(t *testing.T) {
	t.Run("matches Arc.IntersectionLine", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, a := range arcFixtures {
				geomtest.AssertVertices(t, s.IntersectionArc(a).AppendTo(nil), a.IntersectionLine(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, a))
			}
		}
	})
}

func TestLine_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, s := range segmentFixtures {
				actual, ok := l.IntersectionSegment(s)
				expected, expectedOk := s.IntersectionLine(l)

				assert.Equal(t, ok, expectedOk, fmt.Sprintf("%s → %s: ", l, s))
				geomtest.AssertPoint(t, actual, expected, fmt.Sprintf("%s → %s: ", l, s))
			}
		}
	})
}

func TestLine_IntersectionRay(t *testing.T) {
	t.Run("matches Ray.IntersectionLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, r := range rayFixtures {
				actual, ok := l.IntersectionRay(r)
				expected, expectedOk := r.IntersectionLine(l)

				assert.Equal(t, ok, expectedOk, fmt.Sprintf("%s → %s: ", l, r))
				geomtest.AssertPoint(t, actual, expected, fmt.Sprintf("%s → %s: ", l, r))
			}
		}
	})
}

func TestLine_IntersectionLine(t *testing.T) {
	t.Run("crossing lines", func(t *testing.T) {
		point, ok := Ln(Pt(0, 0), Vec(1, 1)).IntersectionLine(Ln(Pt(4, 0), Vec(-1, 1)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("on the far side of both origins", func(t *testing.T) {
		point, ok := Ln(Pt(0.0, 0.0), Vec(1.0, 0.0)).IntersectionLine(Ln(Pt(-3.0, 5.0), Vec(0.0, 2.0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(-3.0, 0.0))
	})
	t.Run("parallel lines have no crossing", func(t *testing.T) {
		_, ok := Ln(Pt(0, 0), Vec(2, 1)).IntersectionLine(Ln(Pt(0, 3), Vec(-4, -2)))

		assert.False(t, ok)
	})
	t.Run("a line has no crossing with itself", func(t *testing.T) {
		_, ok := Ln(Pt(0, 0), Vec(2, 1)).IntersectionLine(Ln(Pt(4, 2), Vec(2, 1)))

		assert.False(t, ok)
	})
	t.Run("int rounds a crossing on a half unit the same way from either side", func(t *testing.T) {
		a, b := Ln(Pt(0, 0), Vec(1, 1)), Ln(Pt(1, 0), Vec(-1, 1))
		point, ok := a.IntersectionLine(b)
		other, _ := b.IntersectionLine(a)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(1, 1))
		geomtest.AssertPoint(t, other, Pt(1, 1))
	})
	t.Run("a zero direction is its origin, the answer where it lies on the other", func(t *testing.T) {
		line, dot := Ln(Pt(0, 0), Vec(2, 1)), Ln(Pt(4, 2), Vec(0, 0))

		point, ok := line.IntersectionLine(dot)
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 2))

		point, ok = dot.IntersectionLine(line)
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 2))

		_, ok = line.IntersectionLine(Ln(Pt(4, 3), Vec(0, 0)))
		assert.False(t, ok)

		point, ok = dot.IntersectionLine(dot)
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 2))
	})
	t.Run("a zero direction within the tolerance of the other origin is the answer from either side", func(t *testing.T) {
		line, dot, near := Ln(Pt(0.0, 0.0), Vec(1.0, 0.0)), Ln(Pt(1e-12, 0.0), Vec(0.0, 0.0)), Ln(Pt(0.0, 0.0), Vec(0.0, 0.0))

		for _, pair := range [][2]Line[float64]{{line, dot}, {dot, line}} {
			point, ok := pair[0].IntersectionLine(pair[1])
			assert.True(t, ok)
			assert.Equal(t, point, dot.Origin)
		}
		for _, pair := range [][2]Line[float64]{{near, dot}, {dot, near}} {
			point, ok := pair[0].IntersectionLine(pair[1])
			assert.True(t, ok)
			assert.Equal(t, point, near.Origin)
		}
	})
	t.Run("the same point whichever line asks, on both lines", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				point, ok := a.IntersectionLine(b)
				other, otherOk := b.IntersectionLine(a)
				message := fmt.Sprintf("%s × %s: ", a, b)

				assert.Equal(t, ok, otherOk, message)
				if !ok {
					continue
				}

				assert.True(t, a.Contains(point) && b.Contains(point), message)
				assert.Equal(t, point, other, message)
			}
		}
	})
	t.Run("int near the end of the supported range answers alike from either side", func(t *testing.T) {
		a, b := Ln(Pt(1<<25, 3), Vec(1<<20+1, 7)), Ln(Pt(-5, 1<<24), Vec(3, -(1<<21)+1))
		point, ok := a.IntersectionLine(b)
		other, otherOk := b.IntersectionLine(a)

		assert.True(t, ok && otherOk)
		assert.Equal(t, point, other)
	})
}

func TestLine_IntersectionPolyline(t *testing.T) {
	path := PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 3.0)})

	t.Run("the crossings of its edges in order, a shared point once", func(t *testing.T) {
		geomtest.AssertVertices(t, LineThrough(Pt(-1.0, -1.0), Pt(5.0, 2.0)).IntersectionPolyline(path), []Point[float64]{Pt(1.0, 0.0), Pt(4.0, 1.5)})
		geomtest.AssertVertices(t, LineThrough(Pt(0.0, 2.0), Pt(6.0, -1.0)).IntersectionPolyline(path), []Point[float64]{Pt(4.0, 0.0)})
		assert.Nil(t, LineThrough(Pt(0.0, 5.0), Pt(-3.0, 5.0)).IntersectionPolyline(path))
		assert.Nil(t, LineThrough(Pt(0.0, 1.0), Pt(3.0, 1.0)).IntersectionPolyline(Polyline[float64]{}))
	})
	t.Run("a closed path is crossed as the polygon of its points", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, s.IntersectionPolyline(p.Polyline()), s.IntersectionPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("appends after what the buffer holds", func(t *testing.T) {
		geomtest.AssertVertices(t, LineThrough(Pt(-1.0, -1.0), Pt(5.0, 2.0)).AppendIntersectionPolyline(bufferWith(prefixPoint), path), []Point[float64]{prefixPoint, Pt(1.0, 0.0), Pt(4.0, 1.5)})
	})
}

func TestLine_IntersectionQuadraticBezier(t *testing.T) {
	t.Run("matches QuadraticBezier.IntersectionLine", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, q := range quadraticFixtures {
				geomtest.AssertVertices(t, s.IntersectionQuadraticBezier(q).AppendTo(nil), q.IntersectionLine(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, q))
			}
		}
	})
}

func TestLine_IntersectionCubicBezier(t *testing.T) {
	t.Run("matches CubicBezier.IntersectionLine", func(t *testing.T) {
		for _, s := range lineFixtures {
			for _, c := range cubicFixtures {
				geomtest.AssertVertices(t, s.IntersectionCubicBezier(c), c.IntersectionLine(s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestLine_AppendIntersectionCubicBezier(t *testing.T) {
	curve, through := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0)), LineThrough(Pt(-1.0, 1.0), Pt(5.0, 1.0+tenth))

	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionCubicBezier(bufferWith(prefixPoint), curve), append([]Point[float64]{prefixPoint}, through.IntersectionCubicBezier(curve)...))
	})
	t.Run("appends to a buffer with room without allocating", func(t *testing.T) {
		buffer := make([]Point[float64], 0, 4)

		assertAllocations(t, func() {
			sinkBool = len(through.AppendIntersectionCubicBezier(buffer[:0], curve)) > 0
		}, 0)
	})
}

func TestLine_IntersectionPolygon(t *testing.T) {
	square := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(0, 4)})
	notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 1.0), Pt(0.0, 4.0)})

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(-2, 2), Vec(1, 0)).IntersectionPolygon(square), []Point[int]{Pt(0, 2), Pt(4, 2)})
		geomtest.AssertVertices(t, Ln(Pt(-2, 2), Vec(-3, 0)).IntersectionPolygon(square), []Point[int]{Pt(4, 2), Pt(0, 2)})
	})
	t.Run("an origin inside gives both", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(2, 2), Vec(0, 1)).IntersectionPolygon(square), []Point[int]{Pt(2, 0), Pt(2, 4)})
	})
	t.Run("through a vertex counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(6, 6), Vec(1, 1)).IntersectionPolygon(square), []Point[int]{Pt(0, 0), Pt(4, 4)})
		geomtest.AssertVertices(t, Ln(Pt(2, 6), Vec(1, -1)).IntersectionPolygon(square), []Point[int]{Pt(4, 4)})
	})
	t.Run("a concave polygon is crossed more than twice", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(9.0, 3.0), Vec(1.0, 0.0)).IntersectionPolygon(notched), []Point[float64]{Pt(0.0, 3.0), Pt(2.0/3, 3.0), Pt(10.0/3, 3.0), Pt(4.0, 3.0)})
	})
	t.Run("passing by and empty give none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-2, 2), Vec(0, 1)).IntersectionPolygon(square))
		assert.Nil(t, Ln(Pt(-2, 2), Vec(1, 0)).IntersectionPolygon(Polygon[int]{}))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, apart := Ln(Pt(-2, 2), Vec(1, 0)), Ln(Pt(-2, 2), Vec(0, 1))

		assertAllocations(t, func() {
			sinkPoints = through.IntersectionPolygon(square)
		}, 1)
		assertAllocations(t, func() {
			sinkPoints = apart.IntersectionPolygon(square)
		}, 0)
	})
	t.Run("every point lies on the line and an edge, and exists where IntersectsPolygon holds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, l := range lineFixtures {
				assertCrossings(t, l, l.IntersectionPolygon(p), l.IntersectsPolygon(p), p.Edges, fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_AppendIntersectionPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := Ln(Pt(-1, 1), Vec(1, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionPolygon([]Point[int]{Pt(2, 1), Pt(9, 9)}, square), []Point[int]{Pt(2, 1), Pt(9, 9), Pt(0, 1), Pt(2, 1)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(-1, 1), Vec(0, 1)).AppendIntersectionPolygon([]Point[int]{Pt(9, 9)}, square), []Point[int]{Pt(9, 9)})
		assert.Nil(t, Ln(Pt(-1, 1), Vec(0, 1)).AppendIntersectionPolygon(nil, square))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		assertAllocations(t, func() {
			sinkPoints = through.AppendIntersectionPolygon(buffer[:0], square)
		}, 0)
	})
	t.Run("matches IntersectionPolygon after the points in dst", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, l.AppendIntersectionPolygon(bufferWith(prefixPoint), p), append([]Point[float64]{prefixPoint}, l.IntersectionPolygon(p)...), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(5.0, 0.0), Vec(1.0, 0.0)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[float64]{Pt(-2.0, 0.0), Pt(2.0, 0.0)})
		geomtest.AssertVertices(t, Ln(Pt(5.0, 0.0), Vec(1.0, 0.0)).IntersectionRectangle(rectangle.Rotate(Pi/2)).AppendTo(nil), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
	})
	t.Run("passing by gives none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-5.0, 0.0), Vec(0.0, 1.0)).IntersectionRectangle(rectangle).AppendTo(nil))
	})
	t.Run("int is exact", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(5, 1), Vec(3, 1)).IntersectionRectangle(Rect(Pt(0, 0), Sz(4, 2), 0)).AppendTo(nil), []Point[int]{Pt(-1, -1), Pt(2, 0)})
	})
	t.Run("every point lies on the line and the boundary, and exists where IntersectsRectangle holds", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, l := range lineFixtures {
				assertCrossings(t, l, l.IntersectionRectangle(rect).AppendTo(nil), l.IntersectsRectangle(rect), rect.Edges, fmt.Sprintf("%s → %s: ", l, rect))
			}
		}
	})
}

func TestLine_IntersectionRegularPolygon(t *testing.T) {
	hexagon := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop)

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(20.0, 0.0), Vec(1.0, 0.0)).IntersectionRegularPolygon(hexagon).AppendTo(nil), []Point[float64]{Pt(-10.0, 0.0), Pt(10.0, 0.0)})
	})
	t.Run("passing by and empty give none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-20.0, 0.0), Vec(0.0, 1.0)).IntersectionRegularPolygon(hexagon).AppendTo(nil))
		assert.Nil(t, Ln(Pt(-20.0, 0.0), Vec(1.0, 0.0)).IntersectionRegularPolygon(RegularPolygon[float64]{}).AppendTo(nil))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, l.IntersectionRegularPolygon(rp).AppendTo(nil), l.IntersectionPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", l, rp))
			}
		}
	})
}

func TestLine_IntersectionBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(9, 1), Vec(1, 0)).IntersectionBox(box).AppendTo(nil), []Point[int]{Pt(0, 1), Pt(4, 1)})
	})
	t.Run("every point lies on the line and the boundary, and exists where IntersectsBox holds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, l := range lineFixtures {
				assertCrossings(t, l, l.IntersectionBox(b).AppendTo(nil), l.IntersectsBox(b), b.Rectangle().Edges, fmt.Sprintf("%s → %s: ", l, b))
			}
		}
	})
}

func TestLine_IntersectionEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives both crossings along the direction", func(t *testing.T) {
		geomtest.AssertVertices(t, Ln(Pt(-8.0, 0.0), Vec(-1.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0), Pt(-4.0, 0.0)})
		geomtest.AssertVertices(t, Ln(Pt(9.0, 2.0), Vec(1.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(0.0, 2.0)})
	})
	t.Run("matches IntersectionCircle on the circle's Ellipse", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, l.IntersectionEllipse(c.Ellipse()).AppendTo(nil), l.IntersectionCircle(c).AppendTo(nil), fmt.Sprintf("%s → %s: ", l, c))
			}
		}
	})
	t.Run("every point lies on the line and the boundary, and exists where IntersectsEllipse holds", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, l := range lineFixtures {
				assertEllipseCrossings(t, l, l.IntersectionEllipse(e).AppendTo(nil), l.IntersectsEllipse(e), e, fmt.Sprintf("%s → %s: ", l, e))
			}
		}
	})
}

func TestLine_ClipCircle(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("passing through gives the chord along the direction", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(-10, 0), Vec(1, 0)).ClipCircle(circle)), []Segment[int]{Seg(Pt(-5, 0), Pt(5, 0))})
		geomtest.AssertSegments(t, partsOf(Ln(Pt(-10, 0), Vec(-1, 0)).ClipCircle(circle)), []Segment[int]{Seg(Pt(5, 0), Pt(-5, 0))})
	})
	t.Run("an origin inside gives the whole chord", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(0, 0), Vec(0, 1)).ClipCircle(circle)), []Segment[int]{Seg(Pt(0, -5), Pt(0, 5))})
	})
	t.Run("a tangent is the point of contact", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(10, 5), Vec(1, 0)).ClipCircle(circle)), []Segment[int]{Seg(Pt(0, 5), Pt(0, 5))})
	})
	t.Run("passing by gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Ln(Pt(-10, 0), Vec(0, 1)).ClipCircle(circle)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-10, 0), Vec(1, 0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipCircle(circle)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				assertClipped(t, l, partsOf(l.ClipCircle(c)), l.IntersectsCircle(c), c.ContainsSegment, fmt.Sprintf("%s → %s: ", l, c))
			}
		}
	})
}

func TestLine_ClipCapsule(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(5.0, 10.0), Vec(0.0, -1.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(5.0, 2.0), Pt(5.0, -2.0))})
		assert.Nil(t, partsOf(Ln(Pt(5.0, 10.0), Vec(1.0, 0.0)).ClipCapsule(capsule)))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range lineFixtures {
			for _, c := range capsuleFixtures {
				assertClipped(t, r, partsOf(r.ClipCapsule(c)), r.IntersectsCapsule(c), c.ContainsSegment, fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestLine_ClipSector(t *testing.T) {
	sector := Sect(Pt(0.0, 0.0), 5.0, Pi/4, 3*Pi/2)

	t.Run("through the notch of a sector past half a turn there are two parts", func(t *testing.T) {
		geomtest.AssertSegments(t, Ln(Pt(3.0, 10.0), Vec(0.0, -1.0)).ClipSector(sector), []Segment[float64]{Seg(Pt(3.0, 4.0), Pt(3.0, 3.0)), Seg(Pt(3.0, -3.0), Pt(3.0, -4.0))})
		assert.Nil(t, Ln(Pt(9.0, 10.0), Vec(0.0, -1.0)).ClipSector(sector))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range lineFixtures {
			for _, s := range sectorFixtures() {
				assertClipped(t, r, r.ClipSector(s), r.IntersectsSector(s), s.ContainsSegment, fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestLine_ClipPolygon(t *testing.T) {
	notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 1.0), Pt(0.0, 4.0)})

	t.Run("a concave polygon cuts the line into parts along the direction", func(t *testing.T) {
		geomtest.AssertSegments(t, Ln(Pt(9.0, 3.0), Vec(1.0, 0.0)).ClipPolygon(notched), []Segment[float64]{
			Seg(Pt(0.0, 3.0), Pt(2.0/3, 3.0)), Seg(Pt(10.0/3, 3.0), Pt(4.0, 3.0)),
		})
	})
	t.Run("an origin inside is passed through", func(t *testing.T) {
		geomtest.AssertSegments(t, Ln(Pt(1.0, 0.5), Vec(1.0, 0.0)).ClipPolygon(notched), []Segment[float64]{Seg(Pt(0.0, 0.5), Pt(4.0, 0.5))})
	})
	t.Run("passing by and empty give none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-1.0, 3.0), Vec(0.0, 1.0)).ClipPolygon(notched))
		assert.Nil(t, Ln(Pt(-1.0, 3.0), Vec(1.0, 0.0)).ClipPolygon(Polygon[float64]{}))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range outlineFixtures() {
				assertClipped(t, l, l.ClipPolygon(p), l.IntersectsPolygon(p), p.ContainsSegment, fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_AppendClipPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := Ln(Pt(-1, 1), Vec(1, 0))
	prefix := Seg(Pt(-7.5, 3.25), Pt(1.0, 1.0))

	t.Run("appends the parts after the segments in dst", func(t *testing.T) {
		geomtest.AssertSegments(t, through.AppendClipPolygon([]Segment[int]{Seg(Pt(9, 9), Pt(9, 9))}, square), []Segment[int]{Seg(Pt(9, 9), Pt(9, 9)), Seg(Pt(0, 1), Pt(2, 1))})
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Segment[int], 0, 1)

		assertAllocations(t, func() {
			sinkSegments = through.AppendClipPolygon(buffer[:0], square)
		}, 0)
	})
	t.Run("matches ClipPolygon after the segments in dst", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range outlineFixtures() {
				geomtest.AssertSegments(t, l.AppendClipPolygon(bufferWith(prefix), p), append([]Segment[float64]{prefix}, l.ClipPolygon(p)...), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_ClipRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(5.0, 0.0), Vec(1.0, 0.0)).ClipRectangle(rectangle)), []Segment[float64]{Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))})
	})
	t.Run("a rotated rectangle clips on its turned edges", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(5.0, 0.0), Vec(1.0, 0.0)).ClipRectangle(rectangle.Rotate(Pi/2))), []Segment[float64]{Seg(Pt(-1.0, 0.0), Pt(1.0, 0.0))})
	})
	t.Run("passing by gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Ln(Pt(-5.0, 0.0), Vec(0.0, 1.0)).ClipRectangle(rectangle)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-5.0, 0.0), Vec(1.0, 0.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipRectangle(rectangle)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, rect := range rectFixtures {
				assertClipped(t, l, partsOf(l.ClipRectangle(rect)), l.IntersectsRectangle(rect), rect.ContainsSegment, fmt.Sprintf("%s → %s: ", l, rect))
			}
		}
	})
}

func TestLine_ClipRegularPolygon(t *testing.T) {
	hexagon := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(20.0, 0.0), Vec(1.0, 0.0)).ClipRegularPolygon(hexagon)), []Segment[float64]{Seg(Pt(-10.0, 0.0), Pt(10.0, 0.0))})
	})
	t.Run("passing by and empty give none", func(t *testing.T) {
		assert.Nil(t, partsOf(Ln(Pt(-20.0, 0.0), Vec(0.0, 1.0)).ClipRegularPolygon(hexagon)))
		assert.Nil(t, partsOf(Ln(Pt(-20.0, 0.0), Vec(1.0, 0.0)).ClipRegularPolygon(RegularPolygon[float64]{})))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-20.0, 0.0), Vec(1.0, 0.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipRegularPolygon(hexagon)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, rp := range regularPolygonFixtures {
				assertClipped(t, l, partsOf(l.ClipRegularPolygon(rp)), l.IntersectsRegularPolygon(rp), rp.ContainsSegment, fmt.Sprintf("%s → %s: ", l, rp))
			}
		}
	})
}

func TestLine_ClipBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("the line as the box shows it, along the direction", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(9, 1), Vec(1, 0)).ClipBox(box)), []Segment[int]{Seg(Pt(0, 1), Pt(4, 1))})
		geomtest.AssertSegments(t, partsOf(Ln(Pt(2, 1), Vec(0, -3)).ClipBox(box)), []Segment[int]{Seg(Pt(2, 2), Pt(2, 0))})
	})
	t.Run("matches ClipRectangle on the box's Rectangle", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, b := range boxFixtures {
				geomtest.AssertSegments(t, partsOf(l.ClipBox(b)), partsOf(l.ClipRectangle(b.Rectangle())), fmt.Sprintf("%s → %s: ", l, b))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-5.0, 0.5), Vec(1.0, 0.0))
		box := Bx(Pt(-2.0, -2.0), Pt(2.0, 2.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipBox(box)
		}, 0)
	})
}

func TestLine_ClipEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives the chord along the direction", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Ln(Pt(0.0, 0.0), Vec(0.0, -1.0)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(0.0, 2.0), Pt(0.0, -2.0))})
	})
	t.Run("passing by gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Ln(Pt(-8.0, 3.0), Vec(1.0, 0.0)).ClipEllipse(ellipse)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Ln(Pt(-8.0, 0.0), Vec(1.0, 0.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipEllipse(ellipse)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, l := range lineFixtures {
				assertClipped(t, l, partsOf(l.ClipEllipse(e)), l.IntersectsEllipse(e), func(part Segment[float64]) bool {
					return e.Contains(part.Start) && e.Contains(part.End)
				}, fmt.Sprintf("%s → %s: ", l, e))
			}
		}
	})
}

func TestLine_Reflect(t *testing.T) {
	t.Run("across an axis-parallel line", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(0, 2), Vec(5, 0)).Reflect(Pt(3, 7)), Pt(3, -3))
		geomtest.AssertPoint(t, Ln(Pt(-1, 0), Vec(0, -3)).Reflect(Pt(3, 7)), Pt(-5, 7))
	})
	t.Run("across a diagonal swaps the coordinates", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(0, 0), Vec(1, 1)).Reflect(Pt(3, 1)), Pt(1, 3))
		geomtest.AssertPoint(t, Ln(Pt(0.0, 0.0), Vec(-2.0, -2.0)).Reflect(Pt(0.5, -1.25)), Pt(-1.25, 0.5))
	})
	t.Run("a point on the line is its own image", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(1, 1), Vec(3, 7)).Reflect(Pt(-5, -13)), Pt(-5, -13))
	})
	t.Run("int rounds once", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(0, 0), Vec(2, 1)).Reflect(Pt(0, 1)), Pt(1, -1))
	})
	t.Run("a zero direction mirrors through the origin", func(t *testing.T) {
		geomtest.AssertPoint(t, Ln(Pt(1, 1), Vec(0, 0)).Reflect(Pt(4, 5)), Pt(-2, -3))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				image := l.Reflect(p)
				message := fmt.Sprintf("%s → %s: ", l, p)

				geomtest.AssertPoint(t, l.Reflect(image), p, message+"twice is the point: ")
				geomtest.AssertPoint(t, p.Midpoint(image), l.Float().Nearest(p), message+"halfway is the foot: ")
				geomtest.AssertNumber(t, l.DistanceTo(image), l.DistanceTo(p), message+"as far from the line: ")
				assert.True(t, l.Contains(p) || l.Side(image) == -l.Side(p), message+"on the other side: ")
			}
		}
	})
}

func TestLine_Equal(t *testing.T) {
	t.Run("same line", func(t *testing.T) {
		assert.True(t, Ln(Pt(1, 2), Vec(3, 4)).Equal(Ln(Pt(1, 2), Vec(3, 4))))
	})
	t.Run("different line", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Vec(3, 4)).Equal(Ln(Pt(1, 3), Vec(3, 4))))
	})
	t.Run("the same points from another origin or direction are another value", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Vec(3, 4)).Equal(Ln(Pt(4, 6), Vec(3, 4))))
		assert.False(t, Ln(Pt(1, 2), Vec(3, 4)).Equal(Ln(Pt(1, 2), Vec(6, 8))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(1.0, 2.0), Vec(3.0, 4.0)).Equal(Ln(Pt(1.0+Delta/2, 2.0), Vec(3.0, 4.0-Delta/2))))
	})
}

func TestLine_IsZero(t *testing.T) {
	t.Run("zero line", func(t *testing.T) {
		assert.True(t, Line[int]{}.IsZero())
	})
	t.Run("a direction alone is not zero", func(t *testing.T) {
		assert.False(t, Line[int]{Vector: Vec(1, 0)}.IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(Delta/2, 0.0), Vec(0.0, -Delta/2)).IsZero())
	})
}

func TestLine_Cast(t *testing.T) {
	line := Ln(Pt(1.4, -2.6), Vec(0.4, 3.5))

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertLine(t, line.Cast[int](), line.Int())
		geomtest.AssertLine(t, line.Int().Cast[float64](), line.Int().Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertLine(t, line.Cast[int8](), Ln(Pt[int8](1, -3), Vec[int8](0, 4)))
	})
}

func TestLine_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Int(), Ln(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("float rounds, and a short direction rounds to zero", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1.5, -2.5), Vec(0.4, 0.3)).Int(), Ln(Pt(2, -3), Vec(0, 0)))
	})
}

func TestLine_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1, 2), Vec(3, 4)).Float(), Ln(Pt(1.0, 2.0), Vec(3.0, 4.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertLine(t, Ln(Pt(1.5, 2.0), Vec(3.0, 4.25)).Float(), Ln(Pt(1.5, 2.0), Vec(3.0, 4.25)))
	})
}

func TestLine_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(10, 16), Vec(1, -2)).String(), "Ln((10,16);⟨1,-2⟩)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(100, -34.0000115), Vec(0.2, 0.4)).String(), "Ln((100.00,-34.00);⟨0.20,0.40⟩)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Ln(Pt(100, -34.0000115), Vec(0.2, 0.4)).String)
	})
}

func TestLine_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Ln(Pt(10, 16), Vec(1, -2)), `{"o":{"x":10,"y":16},"d":{"x":1,"y":-2}}`)

		var l Line[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"o":{"x":10,"y":16},"d":{"x":1,"y":-2}}`), &l))
		geomtest.AssertLine(t, l, Ln(Pt(10, 16), Vec(1, -2)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, line := range lineFixtures {
			data, err := json.Marshal(line)
			assert.NoError(t, err)

			var decoded Line[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, line)
		}
	})
}

func TestLine_Properties(t *testing.T) {
	t.Run("the angle is the direction's", func(t *testing.T) {
		for _, l := range lineFixtures {
			geomtest.AssertNumber(t, l.Angle(), l.Vector.Angle(), l.String()+": ")
		}
	})
	t.Run("every point along the direction lies on the line, either way, on no side", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, t0 := range []float64{-40, -1, 0, 0.25, 1, 40} {
				assert.True(t, l.Contains(l.PointAt(t0)), fmt.Sprintf("%s at %v: ", l, t0))
			}
		}
	})
	t.Run("translate and move to keep the direction", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				assert.True(t, l.MoveTo(p).Vector.Equal(l.Vector), fmt.Sprintf("%s → %s: ", l, p))
				assert.True(t, l.Translate(p.Vector()).Vector.Equal(l.Vector), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
	t.Run("a line through two points contains both", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				line := LineThrough(a, b)

				assert.True(t, line.Contains(a) && line.Contains(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestLine_Immutable(t *testing.T) {
	l := Ln(Pt(1, 2), Vec(3, 5))

	l.Translate(Vec(3, -2))
	l.MoveTo(Pt(4, 3))
	l.Scale(2)
	l.Rotate(Pi / 2)
	l.Lerp(Ln(Pt(0, 0), Vec(1, 1)), 0.5)
	l.Reverse()

	geomtest.AssertLine(t, l, Ln(Pt(1, 2), Vec(3, 5)))
}

// wide returns a segment along the line reaching well past every fixture both ways, the line as
// a caller without one would test it.
func wide[T Number](l Line[T]) Segment[T] {
	return Segment[T]{l.PointAt(-10000), l.PointAt(10000)}
}

// lineFixtures span a zero direction, the axes both ways, diagonals, and short and long
// directions off the origin.
var lineFixtures = []Line[float64]{
	Ln(Pt(0.0, 0.0), Vec(0.0, 0.0)),
	Ln(Pt(0.0, 0.0), Vec(1.0, 0.0)),
	Ln(Pt(-5.0, -5.0), Vec(1.0, 1.0)),
	Ln(Pt(1.0, 2.0), Vec(-0.5, 2.0)),
	Ln(Pt(10.0, 3.0), Vec(-3.0, 0.0)),
	Ln(Pt(0.6, -0.25), Vec(0.0, -1.0)),
	Ln(Pt(-3.5, 7.0), Vec(0.3, -0.2)),
	Ln(Pt(12.5, -0.1), Vec(-13.0, 12.85)),
}

func ExampleLn() {
	fmt.Println(Ln(Pt(1, 2), Vec(3, 4)))
	// Output: Ln((1,2);⟨3,4⟩)
}

func ExampleLineThrough() {
	fmt.Println(LineThrough(Pt(1, 2), Pt(4, 6)))
	// Output: Ln((1,2);⟨3,4⟩)
}

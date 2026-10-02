package geom_test

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestRay_Constructor(t *testing.T) {
	t.Run("along a direction", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, -1), Vec(3, 4)), Ray[int]{Origin: Pt(1, -1), Direction: Vec(3, 4)})
		geomtest.AssertRay(t, RayAlong(Pt(0.5, -1.25), Vec(2.0, 5.0)), Ray[float64]{Origin: Pt(0.5, -1.25), Direction: Vec(2.0, 5.0)})
	})
	t.Run("through a point", func(t *testing.T) {
		geomtest.AssertRay(t, RayThrough(Pt(1, -1), Pt(4, 3)), RayAlong(Pt(1, -1), Vec(3, 4)))
		geomtest.AssertRay(t, RayThrough(Pt(0.5, -1.25), Pt(2.5, 3.75)), RayAlong(Pt(0.5, -1.25), Vec(2.0, 5.0)))
	})
	t.Run("through the origin itself runs nowhere", func(t *testing.T) {
		geomtest.AssertRay(t, RayThrough(Pt(2, 3), Pt(2, 3)), Ray[int]{Origin: Pt(2, 3)})
	})
}

func TestRay_Angle(t *testing.T) {
	t.Run("the angle of the direction", func(t *testing.T) {
		geomtest.AssertNumber(t, RayAlong(Pt(1, 1), Vec(0, 2)).Angle(), Pi/2)
		geomtest.AssertNumber(t, RayAlong(Pt(1.0, 1.0), Vec(-3.0, 0.0)).Angle(), Pi)
	})
	t.Run("a zero direction gives zero", func(t *testing.T) {
		geomtest.AssertNumber(t, RayAlong(Pt(1, 1), Vec(0, 0)).Angle(), 0)
	})
}

func TestRay_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Translate(Vec(-2, 5)), RayAlong(Pt(-1, 7), Vec(3, 4)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(0.5, 2.0), Vec(3.0, 4.0)).Translate(Vec(0.25, -1.5)), RayAlong(Pt(0.75, 0.5), Vec(3.0, 4.0)))
	})
}

func TestRay_MoveTo(t *testing.T) {
	t.Run("starts at the point with the same direction", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).MoveTo(Pt(-5, 0)), RayAlong(Pt(-5, 0), Vec(3, 4)))
	})
}

func TestRay_Scale(t *testing.T) {
	t.Run("scales the direction about the origin", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Scale(2), RayAlong(Pt(1, 2), Vec(6, 8)))
		geomtest.AssertRay(t, RayAlong(Pt(1.0, 2.0), Vec(3.0, 4.0)).Scale(0.5), RayAlong(Pt(1.0, 2.0), Vec(1.5, 2.0)))
	})
	t.Run("per-axis factor changes the direction", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).ScaleXY(2, -1), RayAlong(Pt(1, 2), Vec(6, -4)))
	})
	t.Run("a negative factor turns it to the other side of its origin", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Scale(-1), RayAlong(Pt(1, 2), Vec(-3, -4)))
	})
	t.Run("zero collapses onto the origin", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Scale(0), RayAlong(Pt(1, 2), Vec(0, 0)))
	})
	t.Run("int rounds the direction", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Scale(0.5), RayAlong(Pt(1, 2), Vec(2, 2)))
	})
	t.Run("keeps the points and steps the scaled length", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, factor := range []float64{0.5, 2.5} {
				scaled := r.Scale(factor)

				assert.True(t, r.Contains(scaled.PointAt(3)), fmt.Sprintf("%s x%v: ", r, factor))
				geomtest.AssertNumber(t, scaled.Direction.Length(), r.Direction.Length()*factor, fmt.Sprintf("%s x%v: ", r, factor))
			}
		}
	})
}

func TestRay_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(6, 8)).Unscale(2), RayAlong(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(6, -4)).UnscaleXY(2, -1), RayAlong(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("undoes scale", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				geomtest.AssertRay(t, r.Scale(factor).Unscale(factor), r, fmt.Sprintf("%s x%v: ", r, factor))
			}
		}
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			RayAlong(Pt(1, 2), Vec(3, 4)).Unscale(0)
		})
		assert.Panics(t, func() {
			RayAlong(Pt(1, 2), Vec(3, 4)).UnscaleXY(1, 0)
		})
	})
}

func TestRay_PointAt(t *testing.T) {
	ray := RayAlong(Pt(1.0, 2.0), Vec(2.0, -1.0))

	t.Run("the origin at 0 and one direction away at 1", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.PointAt(0), Pt(1.0, 2.0))
		geomtest.AssertPoint(t, ray.PointAt(1), Pt(3.0, 1.0))
	})
	t.Run("along the ray", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.PointAt(2.5), Pt(6.0, -0.5))
	})
	t.Run("behind the origin for a negative fraction", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.PointAt(-1), Pt(-1.0, 3.0))
	})
	t.Run("int rounds once, on the sum", func(t *testing.T) {
		geomtest.AssertPoint(t, RayAlong(Pt(1, 0), Vec(-1, 0)).PointAt(0.5), Pt(1, 0))
		geomtest.AssertPoint(t, RayAlong(Pt(1, 1), Vec(1, 1)).PointAt(0.5), Pt(2, 2))
	})
}

func TestRay_Transform(t *testing.T) {
	t.Run("a translation moves the origin alone", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Transform(Mat(1.0, 0.0, 5.0, 0.0, 1.0, -1.0)), RayAlong(Pt(6, 1), Vec(3, 4)))
	})
	t.Run("a scale applies to both", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Transform(Mat(2.0, 0.0, 0.0, 0.0, -1.0, 0.0)), RayAlong(Pt(2, -2), Vec(6, -4)))
	})
	t.Run("keeps every point on the transformed ray", func(t *testing.T) {
		matrix := IdentityMatrix[float64]().Rotate(Pi/5).Translate(3, -2)
		for _, r := range rayFixtures {
			transformed := r.Transform(matrix)

			for _, t0 := range []float64{0, 0.5, 3} {
				assert.True(t, transformed.Contains(r.PointAt(t0).Transform(matrix)), fmt.Sprintf("%s at %v: ", r, t0))
			}
		}
	})
}

func TestRay_Rotate(t *testing.T) {
	t.Run("quarter turn about the origin", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 1), Vec(2, 0)).Rotate(Pi/2), RayAlong(Pt(1, 1), Vec(0, 2)))
	})
	t.Run("a full turn is identity", func(t *testing.T) {
		for _, r := range rayFixtures {
			geomtest.AssertRay(t, r.Rotate(2*Pi), r, r.String()+": ")
		}
	})
}

func TestRay_Contains(t *testing.T) {
	ray := RayAlong(Pt(0, 0), Vec(2, 1))

	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, ray.Contains(Pt(0, 0)))
		assert.True(t, ray.Contains(Pt(4, 2)))
		assert.True(t, ray.Contains(Pt(200, 100)))
		assert.False(t, ray.Contains(Pt(-2, -1)))
		assert.False(t, ray.Contains(Pt(1, 1)))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		r := RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.True(t, r.Contains(Pt(5.0, Delta/2)))
		assert.True(t, r.Contains(Pt(-Delta/2, 0.0)))
		assert.False(t, r.Contains(Pt(5.0, 2*Delta)))
		assert.False(t, r.Contains(Pt(-2*Delta, 0.0)))
	})
	t.Run("a zero direction contains its origin alone", func(t *testing.T) {
		r := RayAlong(Pt(1, 1), Vec(0, 0))

		assert.True(t, r.Contains(Pt(1, 1)))
		assert.False(t, r.Contains(Pt(2, 1)))
	})
	t.Run("holds exactly where DistanceTo is zero", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, r.Contains(p), r.DistanceTo(p) == 0, fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
	t.Run("a NaN coordinate is contained by nothing", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
				assert.False(t, r.Contains(p), fmt.Sprintf("%s → %s: ", r, p))
				assert.True(t, math.IsNaN(r.DistanceTo(p)), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_DistanceTo(t *testing.T) {
	ray := RayAlong(Pt(0, 0), Vec(1, 0))

	t.Run("ahead of the origin measures to the line", func(t *testing.T) {
		geomtest.AssertNumber(t, ray.DistanceTo(Pt(100, 3)), 3)
		geomtest.AssertNumber(t, ray.DistanceTo(Pt(1, -4)), 4)
	})
	t.Run("behind the origin measures to the origin", func(t *testing.T) {
		geomtest.AssertNumber(t, ray.DistanceTo(Pt(-3, 4)), 5)
	})
	t.Run("on the ray is exactly zero", func(t *testing.T) {
		assert.Equal(t, RayAlong(Pt(1, 1), Vec(3, 7)).DistanceTo(Pt(7, 15)), 0.0)
	})
	t.Run("a zero direction measures to the origin", func(t *testing.T) {
		geomtest.AssertNumber(t, RayAlong(Pt(1, 1), Vec(0, 0)).DistanceTo(Pt(4, 5)), 5)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		r := RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.Equal(t, r.DistanceTo(Pt(3.0, Delta/2)), 0.0)
		geomtest.AssertNumber(t, r.DistanceTo(Pt(3.0, 2*Delta)), 2*Delta)
	})
}

func TestRay_DistanceSquaredTo(t *testing.T) {
	ray := RayAlong(Pt(0, 0), Vec(2, 1))

	t.Run("stays fractional for an integer T", func(t *testing.T) {
		geomtest.AssertNumber(t, ray.DistanceSquaredTo(Pt(0, 1)), 0.8)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, r.DistanceSquaredTo(p), r.DistanceTo(p)*r.DistanceTo(p), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_Nearest(t *testing.T) {
	ray := RayAlong(Pt(0.0, 0.0), Vec(2.0, 0.0))

	t.Run("the foot of the perpendicular", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.Nearest(Pt(7.0, 3.0)), Pt(7.0, 0.0))
	})
	t.Run("the origin where the foot falls behind it", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.Nearest(Pt(-2.0, 3.0)), Pt(0.0, 0.0))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		geomtest.AssertPoint(t, ray.Nearest(Pt(3.0, Delta/2)), Pt(3.0, Delta/2))
	})
	t.Run("int rounds once and can land off the ray", func(t *testing.T) {
		r := RayAlong(Pt(0, 0), Vec(2, 1))
		nearest := r.Nearest(Pt(0, 1))

		geomtest.AssertPoint(t, nearest, Pt(0, 0))
		geomtest.AssertPoint(t, r.Nearest(Pt(3, 3)), Pt(4, 2))
		geomtest.AssertPoint(t, r.Nearest(Pt(0, 2)), Pt(1, 0))
		assert.False(t, r.Contains(Pt(1, 0)))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		geomtest.AssertPoint(t, RayAlong(Pt(1, 1), Vec(0, 0)).Nearest(Pt(4, 5)), Pt(1, 1))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range pointFixtures {
				nearest := r.Nearest(p)
				message := fmt.Sprintf("%s → %s: ", r, p)

				assert.Equal(t, nearest.Equal(p), r.Contains(p), message)
				assert.True(t, r.Contains(nearest), message)
				assert.True(t, r.Nearest(nearest).Equal(nearest), message)
				geomtest.AssertNumber(t, p.DistanceSquaredTo(nearest), r.DistanceSquaredTo(p), message)
			}
		}
	})
	t.Run("the nearest point of a NaN coordinate is NaN", func(t *testing.T) {
		nearest := ray.Nearest(Pt(math.NaN(), 0.0))

		assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y))
	})
}

func TestRay_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsRay", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, r.IntersectsCircle(c), c.IntersectsRay(r), fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRay_IntersectionCircle(t *testing.T) {
	t.Run("matches Circle.IntersectionRay", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, r.IntersectionCircle(c), c.IntersectionRay(r), fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRay_AppendIntersectionCircle(t *testing.T) {
	t.Run("matches Circle.AppendIntersectionRay", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, r.AppendIntersectionCircle(bufferWith(prefixPoint), c), c.AppendIntersectionRay(bufferWith(prefixPoint), r), fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRay_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsRay", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, r.IntersectsSegment(s), s.IntersectsRay(r), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestRay_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionRay", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, s := range segmentFixtures {
				actual, ok := r.IntersectionSegment(s)
				expected, expectedOk := s.IntersectionRay(r)

				assert.Equal(t, ok, expectedOk, fmt.Sprintf("%s → %s: ", r, s))
				geomtest.AssertPoint(t, actual, expected, fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
}

func TestRay_IntersectsRay(t *testing.T) {
	diagonal := RayAlong(Pt(0, 0), Vec(1, 1))

	t.Run("crossing ahead of both origins", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(4, 0), Vec(-1, 1))))
	})
	t.Run("the lines cross behind an origin", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(4, 0), Vec(1, -1))))
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(-4, 0), Vec(1, -1))))
	})
	t.Run("an origin on the other ray counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(3, 3), Vec(1, 0))))
		assert.True(t, RayAlong(Pt(3, 3), Vec(1, 0)).IntersectsRay(diagonal))
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(0, 0), Vec(-1, 0))))
	})
	t.Run("collinear rays facing each other overlap, facing away they do not", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(5, 5), Vec(-1, -1))))
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(5, 5), Vec(2, 2))))
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(-1, -1), Vec(-1, -1))))
	})
	t.Run("parallel rays do not cross", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(0, 1), Vec(1, 1))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(2, 2), Vec(0, 0))))
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(2, 3), Vec(0, 0))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		r := RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0))

		assert.True(t, r.IntersectsRay(RayAlong(Pt(0.5, Delta/2), Vec(0.0, 1.0))))
		assert.False(t, r.IntersectsRay(RayAlong(Pt(0.5, 2*Delta), Vec(0.0, 1.0))))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range rayFixtures {
			for _, b := range rayFixtures {
				assert.Equal(t, a.IntersectsRay(b), b.IntersectsRay(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRay_IntersectionRay(t *testing.T) {
	diagonal := RayAlong(Pt(0, 0), Vec(1, 1))

	t.Run("crossing", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(RayAlong(Pt(4, 0), Vec(-1, 1)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("the lines cross behind an origin", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(RayAlong(Pt(4, 0), Vec(1, -1)))

		assert.False(t, ok)
	})
	t.Run("an origin on the other ray is the point", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(RayAlong(Pt(3, 3), Vec(1, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(3, 3))

		point, ok = RayAlong(Pt(3, 3), Vec(1, 0)).IntersectionRay(diagonal)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(3, 3))
	})
	t.Run("parallel and collinear rays have no single point", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(RayAlong(Pt(0, 1), Vec(1, 1)))
		assert.False(t, ok)

		_, ok = diagonal.IntersectionRay(RayAlong(Pt(5, 5), Vec(-1, -1)))
		assert.False(t, ok)
	})
	t.Run("a zero direction is the origin where it lies on the other", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(RayAlong(Pt(2, 2), Vec(0, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))

		_, ok = diagonal.IntersectionRay(RayAlong(Pt(2, 3), Vec(0, 0)))
		assert.False(t, ok)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := RayAlong(Pt(0, 0), Vec(1, 0)).IntersectionRay(RayAlong(Pt(1, -1), Vec(1, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 0))
	})
	t.Run("float keeps the crossing", func(t *testing.T) {
		point, ok := RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0)).IntersectionRay(RayAlong(Pt(1.0, -1.0), Vec(1.0, 2.0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(1.5, 0.0))
	})
	t.Run("nearly collinear float rays fall back to the origin on the other", func(t *testing.T) {
		a := RayAlong(Pt(1.6, 2.0), Vec(0.8059999999999999, 23.009999999999998))
		b := RayAlong(Pt(8.854, 209.08999999999997), Vec(2.0149999999999997, 57.52499999999999))
		point, ok := a.IntersectionRay(b)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, b.Origin)
		assert.True(t, a.IntersectsRay(b))
	})
	t.Run("agrees with Intersects on non-parallel fixtures, on both rays", func(t *testing.T) {
		for _, a := range rayFixtures {
			for _, b := range rayFixtures {
				if a.Direction.Float().Cross(b.Direction.Float()) == 0 && a.Direction != (Vector[float64]{}) && b.Direction != (Vector[float64]{}) {
					continue
				}

				point, ok := a.IntersectionRay(b)
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, ok, a.IntersectsRay(b), message)
				if ok {
					assert.True(t, a.Contains(point) && b.Contains(point), message)
				}
			}
		}
	})
}

func TestRay_IntersectsPolygon(t *testing.T) {
	square := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(0, 4)})

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(-2, 2), Vec(1, 0)).IntersectsPolygon(square))
	})
	t.Run("starting inside", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(2, 2), Vec(-1, 3)).IntersectsPolygon(square))
	})
	t.Run("pointing away", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(-2, 2), Vec(-1, 0)).IntersectsPolygon(square))
		assert.False(t, RayAlong(Pt(-2, 2), Vec(0, 1)).IntersectsPolygon(square))
	})
	t.Run("touching a vertex counts", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(2, 6), Vec(1, -1)).IntersectsPolygon(square))
	})
	t.Run("reaches a polygon however far", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(-2, 2), Vec(1, 0)).IntersectsPolygon(square.Translate(Vec(100000, 0))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(0, 0), Vec(1, 0)).IntersectsPolygon(Polygon[int]{}))
	})
	t.Run("matches a segment reaching past the polygon along the ray", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, r.IntersectsPolygon(p), far(r).IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := RayAlong(Pt(-2, 2), Vec(1, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = through.IntersectsPolygon(square)
		}), 0)
	})
}

func TestRay_IntersectionPolygon(t *testing.T) {
	square := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(0, 4)})
	notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 1.0), Pt(0.0, 4.0)})

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-2, 2), Vec(1, 0)).IntersectionPolygon(square), []Point[int]{Pt(0, 2), Pt(4, 2)})
		geomtest.AssertVertices(t, RayAlong(Pt(6, 2), Vec(-3, 0)).IntersectionPolygon(square), []Point[int]{Pt(4, 2), Pt(0, 2)})
	})
	t.Run("starting inside gives the exit", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(2, 2), Vec(0, 1)).IntersectionPolygon(square), []Point[int]{Pt(2, 4)})
	})
	t.Run("through a vertex counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-2, -2), Vec(1, 1)).IntersectionPolygon(square), []Point[int]{Pt(0, 0), Pt(4, 4)})
		geomtest.AssertVertices(t, RayAlong(Pt(2, 6), Vec(1, -1)).IntersectionPolygon(square), []Point[int]{Pt(4, 4)})
	})
	t.Run("a concave polygon is crossed more than twice", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-1.0, 3.0), Vec(1.0, 0.0)).IntersectionPolygon(notched), []Point[float64]{Pt(0.0, 3.0), Pt(2.0/3, 3.0), Pt(10.0/3, 3.0), Pt(4.0, 3.0)})
	})
	t.Run("pointing away and empty give none", func(t *testing.T) {
		assert.Nil(t, RayAlong(Pt(-2, 2), Vec(-1, 0)).IntersectionPolygon(square))
		assert.Nil(t, RayAlong(Pt(-2, 2), Vec(1, 0)).IntersectionPolygon(Polygon[int]{}))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, away := RayAlong(Pt(-2, 2), Vec(1, 0)), RayAlong(Pt(-2, 2), Vec(-1, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.IntersectionPolygon(square)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = away.IntersectionPolygon(square)
		}), 0)
	})
	t.Run("every point lies on the ray and an edge, and exists where IntersectsPolygon holds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, r := range rayFixtures {
				assertRayCrossings(t, r, r.IntersectionPolygon(p), r.IntersectsPolygon(p), p.Edges, fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_AppendIntersectionPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := RayAlong(Pt(-1, 1), Vec(1, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionPolygon([]Point[int]{Pt(2, 1), Pt(9, 9)}, square), []Point[int]{Pt(2, 1), Pt(9, 9), Pt(0, 1), Pt(2, 1)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-1, 1), Vec(-1, 0)).AppendIntersectionPolygon([]Point[int]{Pt(9, 9)}, square), []Point[int]{Pt(9, 9)})
		assert.Nil(t, RayAlong(Pt(-1, 1), Vec(-1, 0)).AppendIntersectionPolygon(nil, square))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionPolygon(buffer[:0], square)
		}), 0)
	})
	t.Run("matches IntersectionPolygon after the points in dst", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, r.AppendIntersectionPolygon(bufferWith(prefixPoint), p), append([]Point[float64]{prefixPoint}, r.IntersectionPolygon(p)...), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0))

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0)).IntersectsRectangle(rectangle))
	})
	t.Run("starting inside", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(0.5, 0.5), Vec(1.0, 7.0)).IntersectsRectangle(rectangle))
	})
	t.Run("pointing away", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(-5.0, 0.0), Vec(-1.0, 0.0)).IntersectsRectangle(rectangle))
	})
	t.Run("a rotated rectangle is met on its turned edges", func(t *testing.T) {
		turned := rectangle.Rotate(Pi / 4)

		assert.True(t, RayAlong(Pt(-5.0, 1.6), Vec(1.0, 0.0)).IntersectsRectangle(turned))
		assert.False(t, RayAlong(Pt(-5.0, 1.6), Vec(1.0, 0.0)).IntersectsRectangle(rectangle))
	})
	t.Run("matches a segment reaching past the rectangle along the ray", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rect := range rectFixtures {
				assert.Equal(t, r.IntersectsRectangle(rect), far(r).IntersectsRectangle(rect), fmt.Sprintf("%s → %s: ", r, rect))
			}
		}
	})
}

func TestRay_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0))

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0)).IntersectionRectangle(rectangle), []Point[float64]{Pt(-2.0, 0.0), Pt(2.0, 0.0)})
		geomtest.AssertVertices(t, RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0)).IntersectionRectangle(rectangle.Rotate(Pi/2)), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
	})
	t.Run("pointing away gives none", func(t *testing.T) {
		assert.Nil(t, RayAlong(Pt(-5.0, 0.0), Vec(-1.0, 0.0)).IntersectionRectangle(rectangle))
	})
	t.Run("int is exact", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-7, -3), Vec(3, 1)).IntersectionRectangle(Rect(Pt(0, 0), Sz(4, 2))), []Point[int]{Pt(-1, -1), Pt(2, 0)})
	})
	t.Run("every point lies on the ray and the boundary, and exists where IntersectsRectangle holds", func(t *testing.T) {
		for _, rect := range rectFixtures {
			for _, r := range rayFixtures {
				assertRayCrossings(t, r, r.IntersectionRectangle(rect), r.IntersectsRectangle(rect), rect.Edges, fmt.Sprintf("%s → %s: ", r, rect))
			}
		}
	})
}

func TestRay_AppendIntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 2))
	through := RayAlong(Pt(-5, 0), Vec(1, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionRectangle([]Point[int]{Pt(2, 0), Pt(9, 9)}, rectangle), []Point[int]{Pt(2, 0), Pt(9, 9), Pt(-2, 0), Pt(2, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-5, 0), Vec(-1, 0)).AppendIntersectionRectangle([]Point[int]{Pt(9, 9)}, rectangle), []Point[int]{Pt(9, 9)})
		assert.Nil(t, RayAlong(Pt(-5, 0), Vec(-1, 0)).AppendIntersectionRectangle(nil, rectangle))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionRectangle(buffer[:0], rectangle)
		}), 0)
	})
	t.Run("matches IntersectionRectangle after the points in dst", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rect := range rectFixtures {
				geomtest.AssertVertices(t, r.AppendIntersectionRectangle(bufferWith(prefixPoint), rect), append([]Point[float64]{prefixPoint}, r.IntersectionRectangle(rect)...), fmt.Sprintf("%s → %s: ", r, rect))
			}
		}
	})
}

func TestRay_IntersectsRegularPolygon(t *testing.T) {
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0)).IntersectsRegularPolygon(hexagon))
	})
	t.Run("pointing away", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(-20.0, 0.0), Vec(0.0, 1.0)).IntersectsRegularPolygon(hexagon))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0)).IntersectsRegularPolygon(RegularPolygon[float64]{}))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, r.IntersectsRegularPolygon(rp), r.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRay_IntersectionRegularPolygon(t *testing.T) {
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0)).IntersectionRegularPolygon(hexagon), []Point[float64]{Pt(-10.0, 0.0), Pt(10.0, 0.0)})
	})
	t.Run("pointing away and empty give none", func(t *testing.T) {
		assert.Nil(t, RayAlong(Pt(-20.0, 0.0), Vec(-1.0, 0.0)).IntersectionRegularPolygon(hexagon))
		assert.Nil(t, RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0)).IntersectionRegularPolygon(RegularPolygon[float64]{}))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, r.IntersectionRegularPolygon(rp), r.IntersectionPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRay_AppendIntersectionRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)
	through := RayAlong(Pt(-3, 0), Vec(1, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionRegularPolygon([]Point[int]{Pt(2, 0), Pt(9, 9)}, diamond), []Point[int]{Pt(2, 0), Pt(9, 9), Pt(-2, 0), Pt(2, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-3, 0), Vec(-1, 0)).AppendIntersectionRegularPolygon([]Point[int]{Pt(9, 9)}, diamond), []Point[int]{Pt(9, 9)})
		assert.Nil(t, RayAlong(Pt(-3, 0), Vec(-1, 0)).AppendIntersectionRegularPolygon(nil, diamond))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionRegularPolygon(buffer[:0], diamond)
		}), 0)
	})
	t.Run("matches IntersectionRegularPolygon after the points in dst", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, r.AppendIntersectionRegularPolygon(bufferWith(prefixPoint), rp), append([]Point[float64]{prefixPoint}, r.IntersectionRegularPolygon(rp)...), fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRay_IntersectsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(-1, 1), Vec(1, 0)).IntersectsBox(box))
	})
	t.Run("pointing away", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(-1, 1), Vec(-1, 0)).IntersectsBox(box))
	})
	t.Run("matches IntersectsRectangle on the box's Rectangle", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, r.IntersectsBox(b), r.IntersectsRectangle(b.Rectangle()), fmt.Sprintf("%s → %s: ", r, b))
			}
		}
	})
}

func TestRay_IntersectionBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("passing through gives both crossings from Origin on", func(t *testing.T) {
		geomtest.AssertVertices(t, RayAlong(Pt(-1, 1), Vec(1, 0)).IntersectionBox(box), []Point[int]{Pt(0, 1), Pt(4, 1)})
	})
	t.Run("every point lies on the ray and the boundary, and exists where IntersectsBox holds", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, r := range rayFixtures {
				assertRayCrossings(t, r, r.IntersectionBox(b), r.IntersectsBox(b), b.Rectangle().Edges, fmt.Sprintf("%s → %s: ", r, b))
			}
		}
	})
}

func TestRay_AppendIntersectionBox(t *testing.T) {
	t.Run("matches AppendIntersectionRectangle on the box Rectangle", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, b := range boxFixtures {
				geomtest.AssertVertices(t, r.AppendIntersectionBox(bufferWith(prefixPoint), b), r.AppendIntersectionRectangle(bufferWith(prefixPoint), b.Rectangle()), fmt.Sprintf("%s → %s: ", r, b))
			}
		}
	})
}

func TestRay_ClipCircle(t *testing.T) {
	circle := Circ(Pt(0, 0), 5)

	t.Run("passing through gives the chord", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-10, 0), Vec(1, 0)).ClipCircle(circle)), []Segment[int]{Seg(Pt(-5, 0), Pt(5, 0))})
	})
	t.Run("the origin inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(0, 0), Vec(0, 1)).ClipCircle(circle)), []Segment[int]{Seg(Pt(0, 0), Pt(0, 5))})
	})
	t.Run("a tangent is the point of contact", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-10, 5), Vec(1, 0)).ClipCircle(circle)), []Segment[int]{Seg(Pt(0, 5), Pt(0, 5))})
	})
	t.Run("pointing away gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(RayAlong(Pt(-10, 0), Vec(-1, 0)).ClipCircle(circle)))
		assert.Nil(t, partsOf(RayAlong(Pt(1, 5), Vec(1, 0)).ClipCircle(circle)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := RayAlong(Pt(-10, 0), Vec(1, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipCircle(circle)
		}), 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, c := range circleFixtures {
				assertRayClipped(t, r, partsOf(r.ClipCircle(c)), r.IntersectsCircle(c), c.EnclosesSegment, fmt.Sprintf("%s → %s: ", r, c))
			}
		}
	})
}

func TestRay_ClipPolygon(t *testing.T) {
	notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 1.0), Pt(0.0, 4.0)})

	t.Run("a concave polygon cuts the ray into parts from Origin on", func(t *testing.T) {
		geomtest.AssertSegments(t, RayAlong(Pt(-1.0, 3.0), Vec(1.0, 0.0)).ClipPolygon(notched), []Segment[float64]{
			Seg(Pt(0.0, 3.0), Pt(2.0/3, 3.0)), Seg(Pt(10.0/3, 3.0), Pt(4.0, 3.0)),
		})
	})
	t.Run("the origin inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, RayAlong(Pt(1.0, 0.5), Vec(1.0, 0.0)).ClipPolygon(notched), []Segment[float64]{Seg(Pt(1.0, 0.5), Pt(4.0, 0.5))})
	})
	t.Run("pointing away and empty give none", func(t *testing.T) {
		assert.Nil(t, RayAlong(Pt(-1.0, 3.0), Vec(-1.0, 0.0)).ClipPolygon(notched))
		assert.Nil(t, RayAlong(Pt(-1.0, 3.0), Vec(1.0, 0.0)).ClipPolygon(Polygon[float64]{}))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range outlineFixtures() {
				assertRayClipped(t, r, r.ClipPolygon(p), r.IntersectsPolygon(p), p.EnclosesSegment, fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_AppendClipPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := RayAlong(Pt(-1, 1), Vec(1, 0))
	prefix := Seg(Pt(-7.5, 3.25), Pt(1.0, 1.0))

	t.Run("appends the parts after the segments in dst", func(t *testing.T) {
		geomtest.AssertSegments(t, through.AppendClipPolygon([]Segment[int]{Seg(Pt(9, 9), Pt(9, 9))}, square), []Segment[int]{Seg(Pt(9, 9), Pt(9, 9)), Seg(Pt(0, 1), Pt(2, 1))})
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Segment[int], 0, 1)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkSegments = through.AppendClipPolygon(buffer[:0], square)
		}), 0)
	})
	t.Run("matches ClipPolygon after the segments in dst", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range outlineFixtures() {
				geomtest.AssertSegments(t, r.AppendClipPolygon(bufferWith(prefix), p), append([]Segment[float64]{prefix}, r.ClipPolygon(p)...), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
}

func TestRay_ClipRectangle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0))

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0)).ClipRectangle(rectangle)), []Segment[float64]{Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))})
	})
	t.Run("a rotated rectangle clips on its turned edges", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0)).ClipRectangle(rectangle.Rotate(Pi/2))), []Segment[float64]{Seg(Pt(-1.0, 0.0), Pt(1.0, 0.0))})
	})
	t.Run("pointing away gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(RayAlong(Pt(-5.0, 0.0), Vec(-1.0, 0.0)).ClipRectangle(rectangle)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := RayAlong(Pt(-5.0, 0.0), Vec(1.0, 0.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipRectangle(rectangle)
		}), 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rect := range rectFixtures {
				assertRayClipped(t, r, partsOf(r.ClipRectangle(rect)), r.IntersectsRectangle(rect), rect.EnclosesSegment, fmt.Sprintf("%s → %s: ", r, rect))
			}
		}
	})
}

func TestRay_ClipRegularPolygon(t *testing.T) {
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0)).ClipRegularPolygon(hexagon)), []Segment[float64]{Seg(Pt(-10.0, 0.0), Pt(10.0, 0.0))})
	})
	t.Run("pointing away and empty give none", func(t *testing.T) {
		assert.Nil(t, partsOf(RayAlong(Pt(-20.0, 0.0), Vec(-1.0, 0.0)).ClipRegularPolygon(hexagon)))
		assert.Nil(t, partsOf(RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0)).ClipRegularPolygon(RegularPolygon[float64]{})))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := RayAlong(Pt(-20.0, 0.0), Vec(1.0, 0.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipRegularPolygon(hexagon)
		}), 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, rp := range regularPolygonFixtures {
				assertRayClipped(t, r, partsOf(r.ClipRegularPolygon(rp)), r.IntersectsRegularPolygon(rp), rp.EnclosesSegment, fmt.Sprintf("%s → %s: ", r, rp))
			}
		}
	})
}

func TestRay_ClipBox(t *testing.T) {
	box := BoxFromMinMax(Pt(0, 0), Pt(4, 2))

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(RayAlong(Pt(-1, 1), Vec(1, 0)).ClipBox(box)), []Segment[int]{Seg(Pt(0, 1), Pt(4, 1))})
	})
	t.Run("matches ClipRectangle on the box's Rectangle", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, b := range boxFixtures {
				geomtest.AssertSegments(t, partsOf(r.ClipBox(b)), partsOf(r.ClipRectangle(b.Rectangle())), fmt.Sprintf("%s → %s: ", r, b))
			}
		}
	})
}

func TestRay_Equal(t *testing.T) {
	t.Run("same ray", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(1, 2), Vec(3, 4)).Equal(RayAlong(Pt(1, 2), Vec(3, 4))))
	})
	t.Run("different ray", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(1, 2), Vec(3, 4)).Equal(RayAlong(Pt(1, 3), Vec(3, 4))))
	})
	t.Run("a direction of another length is another value", func(t *testing.T) {
		assert.False(t, RayAlong(Pt(1, 2), Vec(3, 4)).Equal(RayAlong(Pt(1, 2), Vec(6, 8))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(1.0, 2.0), Vec(3.0, 4.0)).Equal(RayAlong(Pt(1.0+Delta/2, 2.0), Vec(3.0, 4.0-Delta/2))))
	})
}

func TestRay_IsZero(t *testing.T) {
	t.Run("zero ray", func(t *testing.T) {
		assert.True(t, Ray[int]{}.IsZero())
	})
	t.Run("a direction alone is not zero", func(t *testing.T) {
		assert.False(t, Ray[int]{Direction: Vec(1, 0)}.IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, RayAlong(Pt(Delta/2, 0.0), Vec(0.0, -Delta/2)).IsZero())
	})
}

func TestRay_Cast(t *testing.T) {
	ray := RayAlong(Pt(1.4, -2.6), Vec(0.4, 3.5))

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertRay(t, ray.Cast[int](), ray.Int())
		geomtest.AssertRay(t, ray.Int().Cast[float64](), ray.Int().Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertRay(t, ray.Cast[int8](), RayAlong(Pt[int8](1, -3), Vec[int8](0, 4)))
	})
}

func TestRay_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Int(), RayAlong(Pt(1, 2), Vec(3, 4)))
	})
	t.Run("float rounds, and a short direction rounds to zero", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1.5, -2.5), Vec(0.4, 0.3)).Int(), RayAlong(Pt(2, -3), Vec(0, 0)))
	})
}

func TestRay_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1, 2), Vec(3, 4)).Float(), RayAlong(Pt(1.0, 2.0), Vec(3.0, 4.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertRay(t, RayAlong(Pt(1.5, 2.0), Vec(3.0, 4.25)).Float(), RayAlong(Pt(1.5, 2.0), Vec(3.0, 4.25)))
	})
}

func TestRay_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, RayAlong(Pt(10, 16), Vec(1, -2)).String(), "Ray((10,16);⟨1,-2⟩)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, RayAlong(Pt(100, -34.0000115), Vec(0.2, 0.4)).String(), "Ray((100.00,-34.00);⟨0.20,0.40⟩)")
	})
}

func TestRay_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, RayAlong(Pt(10, 16), Vec(1, -2)), `{"o":{"x":10,"y":16},"d":{"x":1,"y":-2}}`)

		var r Ray[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"o":{"x":10,"y":16},"d":{"x":1,"y":-2}}`), &r))
		geomtest.AssertRay(t, r, RayAlong(Pt(10, 16), Vec(1, -2)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, ray := range rayFixtures {
			data, err := json.Marshal(ray)
			assert.NoError(t, err)

			var decoded Ray[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, ray)
		}
	})
}

func TestRay_Properties(t *testing.T) {
	t.Run("the angle is the direction's", func(t *testing.T) {
		for _, r := range rayFixtures {
			geomtest.AssertNumber(t, r.Angle(), r.Direction.Angle(), r.String()+": ")
		}
	})
	t.Run("every point ahead of the origin lies on the ray, none behind it", func(t *testing.T) {
		for _, r := range rayFixtures {
			if r.Direction == (Vector[float64]{}) {
				continue
			}

			for _, t0 := range []float64{0, 0.25, 1, 40} {
				assert.True(t, r.Contains(r.PointAt(t0)), fmt.Sprintf("%s at %v: ", r, t0))
			}
			assert.False(t, r.Contains(r.PointAt(-1)), r.String()+": ")
		}
	})
	t.Run("translate and move to keep the direction", func(t *testing.T) {
		for _, r := range rayFixtures {
			for _, p := range pointFixtures {
				assert.True(t, r.MoveTo(p).Direction.Equal(r.Direction), fmt.Sprintf("%s → %s: ", r, p))
				assert.True(t, r.Translate(p.Vector()).Direction.Equal(r.Direction), fmt.Sprintf("%s → %s: ", r, p))
			}
		}
	})
	t.Run("a ray through a point contains it", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				assert.True(t, RayThrough(a, b).Contains(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRay_Immutable(t *testing.T) {
	r := RayAlong(Pt(1, 2), Vec(3, 5))

	r.Translate(Vec(3, -2))
	r.MoveTo(Pt(4, 3))
	r.Scale(2)
	r.Rotate(Pi / 2)

	geomtest.AssertRay(t, r, RayAlong(Pt(1, 2), Vec(3, 5)))
}

// far returns a segment along the ray reaching well past every fixture, the ray as a caller
// without one would cast it.
func far[T Number](r Ray[T]) Segment[T] {
	return Segment[T]{r.Origin, r.PointAt(10000)}
}

// assertRayCrossings checks the boundary crossings of a ray over the fixtures: each lies on the
// ray and on an edge, they follow one another from Origin, and there are some only where the ray
// intersects the shape.
func assertRayCrossings[T Number](t *testing.T, r Ray[T], points []Point[T], intersects bool, edges func() iter.Seq[Segment[T]], message string) {
	t.Helper()

	reached := 0.0
	for _, point := range points {
		distance := r.Origin.DistanceTo(point)

		assert.True(t, r.Contains(point), message+point.String()+" on the ray: ")
		assert.True(t, slices.ContainsFunc(slices.Collect(edges()), func(edge Segment[T]) bool {
			return edge.Contains(point)
		}), message+point.String()+" on the boundary: ")
		assert.True(t, LessOrEqual(reached, distance), message+point.String()+" in order: ")

		reached = distance
	}
	if len(points) > 0 {
		assert.True(t, intersects, message)
	}
}

// assertRayClipped checks the parts a Clip method returns for a ray over the fixtures, as
// assertClipped checks them for a segment: they exist exactly where the ray intersects the
// shape, each lies on the ray and within the shape, and they follow one another from Origin
// without overlapping.
func assertRayClipped[T Number](t *testing.T, r Ray[T], parts []Segment[T], intersects bool, encloses func(Segment[T]) bool, message string) {
	t.Helper()

	assert.Equal(t, len(parts) > 0, intersects, message)

	reached := 0.0
	for _, part := range parts {
		from, to := r.Origin.DistanceTo(part.Start), r.Origin.DistanceTo(part.End)

		assert.True(t, r.Contains(part.Start) && r.Contains(part.End), message+part.String()+" on the ray: ")
		assert.True(t, encloses(part), message+part.String()+" within the shape: ")
		assert.True(t, LessOrEqual(reached, from) && LessOrEqual(from, to), message+part.String()+" in order: ")

		reached = to
	}
}

// rayFixtures span a zero direction, the axes both ways, diagonals, and short and long
// directions off the origin.
var rayFixtures = []Ray[float64]{
	RayAlong(Pt(0.0, 0.0), Vec(0.0, 0.0)),
	RayAlong(Pt(0.0, 0.0), Vec(1.0, 0.0)),
	RayAlong(Pt(-5.0, -5.0), Vec(1.0, 1.0)),
	RayAlong(Pt(1.0, 2.0), Vec(-0.5, 2.0)),
	RayAlong(Pt(10.0, 3.0), Vec(-3.0, 0.0)),
	RayAlong(Pt(0.6, -0.25), Vec(0.0, -1.0)),
	RayAlong(Pt(-3.5, 7.0), Vec(0.3, -0.2)),
	RayAlong(Pt(12.5, -0.1), Vec(-13.0, 12.85)),
}

func ExampleRayAlong() {
	fmt.Println(RayAlong(Pt(1, 2), Vec(3, 4)))
	// Output: Ray((1,2);⟨3,4⟩)
}

func ExampleRayThrough() {
	fmt.Println(RayThrough(Pt(1, 2), Pt(4, 6)))
	// Output: Ray((1,2);⟨3,4⟩)
}

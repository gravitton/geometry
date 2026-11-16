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

func TestCubicBezier_Constructor(t *testing.T) {
	geomtest.AssertCubicBezier(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)), CubicBezier[int]{Start: Pt(0, 0), Control1: Pt(0, 4), Control2: Pt(8, 4), End: Pt(8, 0)})
}

func TestCubicBezier_Bounds(t *testing.T) {
	t.Run("the curve, not its control points", func(t *testing.T) {
		geomtest.AssertBox(t, archCurve().Bounds(), Bx(Pt(0.0, 0.0), Pt(8.0, 3.0)))
	})
	t.Run("a curve turning back along both axes", func(t *testing.T) {
		bounds := CubBez(Pt(0.0, 0.0), Pt(20.0, 0.0), Pt(-14.0, 6.0), Pt(6.0, 6.0)).Bounds()

		assert.True(t, bounds.Min.X < 0 && bounds.Max.X > 6, bounds.String())
		geomtest.AssertNumber(t, bounds.Min.Y, 0.0)
		geomtest.AssertNumber(t, bounds.Max.Y, 6.0)
	})
	t.Run("a straight curve is the box of its ends", func(t *testing.T) {
		geomtest.AssertBox(t, CubBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(2.0, 2.0), Pt(3.0, 3.0)).Bounds(), Bx(Pt(0.0, 0.0), Pt(3.0, 3.0)))
	})
	t.Run("int is rounded outward", func(t *testing.T) {
		geomtest.AssertBox(t, CubBez(Pt(0, 0), Pt(0, 3), Pt(8, 3), Pt(8, 0)).Bounds(), Bx(Pt(0, 0), Pt(8, 3)))
		geomtest.AssertBox(t, CubBez(Pt(0, 0), Pt(0, -3), Pt(8, -3), Pt(8, 0)).Bounds(), Bx(Pt(0, -3), Pt(8, 0)))
	})
}

func TestCubicBezier_Translate(t *testing.T) {
	geomtest.AssertCubicBezier(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).Translate(Vec(1, -2)), CubBez(Pt(1, -2), Pt(1, 2), Pt(9, 2), Pt(9, -2)))
}

func TestCubicBezier_Transform(t *testing.T) {
	t.Run("the four points go through the matrix", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).Transform(ScaleMatrix(2.0, 0.5)), CubBez(Pt(0, 0), Pt(0, 2), Pt(16, 2), Pt(16, 0)))
	})
	t.Run("a point of the transformed curve is the transformed point", func(t *testing.T) {
		matrix := TRSMatrix(3.0, -2.0, 0.7, 2.0, 0.5)

		for _, curve := range cubicFixtures {
			for _, fraction := range curveFractions {
				geomtest.AssertPoint(t, curve.Transform(matrix).PointAt(fraction), curve.PointAt(fraction).Transform(matrix), fmt.Sprintf("%s @%v: ", curve, fraction))
			}
		}
	})
}

func TestCubicBezier_Lerp(t *testing.T) {
	a, b := CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)), CubBez(Pt(2, 2), Pt(4, 0), Pt(10, 8), Pt(12, 4))

	t.Run("every point moves towards the point of the other", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, a.Lerp(b, 0.5), CubBez(Pt(1, 1), Pt(2, 2), Pt(9, 6), Pt(10, 2)))
	})
	t.Run("the ends are the two curves", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, a.Lerp(b, 0), a)
		geomtest.AssertCubicBezier(t, a.Lerp(b, 1), b)
	})
	t.Run("a point of the curve between is between the points of the two", func(t *testing.T) {
		for i, curve := range cubicFixtures {
			other := cubicFixtures[(i+1)%len(cubicFixtures)]

			for _, toward := range []float64{-0.5, 0.25, 1.5} {
				for _, fraction := range curveFractions {
					geomtest.AssertPoint(t, curve.Lerp(other, toward).PointAt(fraction), curve.PointAt(fraction).Lerp(other.PointAt(fraction), toward), fmt.Sprintf("%s → %s @%v: ", curve, other, fraction))
				}
			}
		}
	})
}

func TestCubicBezier_Reverse(t *testing.T) {
	t.Run("the points in the opposite order", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 5), Pt(8, 0)).Reverse(), CubBez(Pt(8, 0), Pt(8, 5), Pt(0, 4), Pt(0, 0)))
	})
	t.Run("the same points the other way round", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			for _, fraction := range curveFractions {
				geomtest.AssertPoint(t, curve.Reverse().PointAt(1-fraction), curve.PointAt(fraction), fmt.Sprintf("%s @%v: ", curve, fraction))
			}
		}
	})
}

func TestCubicBezier_PointAt(t *testing.T) {
	t.Run("the ends", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			geomtest.AssertPoint(t, curve.PointAt(0), curve.Start, curve.String())
			geomtest.AssertPoint(t, curve.PointAt(1), curve.End, curve.String())
		}
	})
	t.Run("the middle of an arch", func(t *testing.T) {
		geomtest.AssertPoint(t, archCurve().PointAt(0.5), Pt(4.0, 3.0))
	})
	t.Run("the Bernstein form", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			for _, t1 := range []float64{-0.5, 0.1, 0.35, 0.5, 0.8, 1.5} {
				u := 1 - t1
				expected := Pt(
					u*u*u*curve.Start.X+3*u*u*t1*curve.Control1.X+3*u*t1*t1*curve.Control2.X+t1*t1*t1*curve.End.X,
					u*u*u*curve.Start.Y+3*u*u*t1*curve.Control1.Y+3*u*t1*t1*curve.Control2.Y+t1*t1*t1*curve.End.Y,
				)

				geomtest.AssertPoint(t, curve.PointAt(t1), expected, fmt.Sprintf("%s @%v: ", curve, t1))
			}
		}
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, CubBez(Pt(0, 0), Pt(0, 3), Pt(7, 3), Pt(7, 0)).PointAt(0.5), Pt(4, 2))
	})
}

func TestCubicBezier_TangentAt(t *testing.T) {
	t.Run("toward the first control point at the start and from the second at the end", func(t *testing.T) {
		geomtest.AssertVector(t, archCurve().TangentAt(0), Vec(0.0, 12.0))
		geomtest.AssertVector(t, archCurve().TangentAt(1), Vec(0.0, -12.0))
		geomtest.AssertVector(t, archCurve().TangentAt(0.5), Vec(12.0, 0.0))
	})
	t.Run("no direction where the control point lies on the end", func(t *testing.T) {
		assert.True(t, CubBez(Pt(1.0, 1.0), Pt(1.0, 1.0), Pt(8.0, 4.0), Pt(8.0, 0.0)).TangentAt(0).IsZero())
	})
	t.Run("the direction the curve moves in", func(t *testing.T) {
		const step = 1e-6

		for _, curve := range cubicFixtures {
			for _, fraction := range curveFractions {
				moved := curve.PointAt(fraction + step).Subtract(curve.PointAt(fraction - step)).Divide(2 * step)

				assert.True(t, curve.TangentAt(fraction).Subtract(moved).LessOrEqual(1e-4), fmt.Sprintf("%s @%v: %v ≠ %v", curve, fraction, curve.TangentAt(fraction), moved))
			}
		}
	})
}

func TestCubicBezier_Split(t *testing.T) {
	t.Run("the halves of an arch", func(t *testing.T) {
		first, second := archCurve().Split(0.5)

		geomtest.AssertCubicBezier(t, first, CubBez(Pt(0.0, 0.0), Pt(0.0, 2.0), Pt(2.0, 3.0), Pt(4.0, 3.0)))
		geomtest.AssertCubicBezier(t, second, CubBez(Pt(4.0, 3.0), Pt(6.0, 3.0), Pt(8.0, 2.0), Pt(8.0, 0.0)))
	})
	t.Run("the two trace the curve", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			for _, at := range []float64{0.25, 0.5, 0.9} {
				first, second := curve.Split(at)

				for _, fraction := range curveFractions {
					message := fmt.Sprintf("%s /%v @%v: ", curve, at, fraction)

					geomtest.AssertPoint(t, first.PointAt(fraction), curve.PointAt(fraction*at), message)
					geomtest.AssertPoint(t, second.PointAt(fraction), curve.PointAt(at+fraction*(1-at)), message)
				}
			}
		}
	})
	t.Run("int keeps the ends and rounds the rest", func(t *testing.T) {
		first, second := CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).Split(0.5)

		geomtest.AssertCubicBezier(t, first, CubBez(Pt(0, 0), Pt(0, 2), Pt(2, 3), Pt(4, 3)))
		geomtest.AssertCubicBezier(t, second, CubBez(Pt(4, 3), Pt(6, 3), Pt(8, 2), Pt(8, 0)))
	})
}

func TestCubicBezier_IntersectionSegment(t *testing.T) {
	wave := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0))

	t.Run("the three crossings of a line through a wave, from Start", func(t *testing.T) {
		geomtest.AssertVertices(t, wave.IntersectionSegment(Seg(Pt(-1.0, 0.0), Pt(7.0, 0.0))), []Point[float64]{Pt(0.0, 0.0), Pt(3.0, 0.0), Pt(6.0, 0.0)})
		geomtest.AssertVertices(t, wave.IntersectionSegment(Seg(Pt(7.0, 0.0), Pt(-1.0, 0.0))), []Point[float64]{Pt(6.0, 0.0), Pt(3.0, 0.0), Pt(0.0, 0.0)})
	})
	t.Run("one crossing, and none", func(t *testing.T) {
		geomtest.AssertVertices(t, wave.IntersectionSegment(Seg(Pt(3.0, 5.0), Pt(3.0, -5.0))), []Point[float64]{Pt(3.0, 0.0)})
		assert.Nil(t, wave.IntersectionSegment(Seg(Pt(-1.0, 3.0), Pt(7.0, 3.0))))
		assert.Nil(t, wave.IntersectionSegment(Seg(Pt(3.0, 5.0), Pt(3.0, 1.0))))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, CubBez(Pt(0, 0), Pt(20, 60), Pt(40, -60), Pt(60, 0)).IntersectionSegment(Seg(Pt(-10, 0), Pt(70, 0))), []Point[int]{Pt(0, 0), Pt(30, 0), Pt(60, 0)})
	})
	t.Run("every point lies on the segment and the curve, three at most", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			path := curve.Polyline(1e-6)

			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					points := curve.IntersectionSegment(s)
					message := fmt.Sprintf("%s → %s: ", curve, s)

					for _, point := range points {
						assert.True(t, s.Contains(point), message+point.String())
						assert.True(t, path.DistanceTo(point) <= 1e-5, message+point.String())
					}
					assertOrdered(t, s, points, message)
					assert.True(t, len(points) <= 3, message)
				}
			}
		}
	})
	t.Run("a curve passing the segment within the tolerance touches it once", func(t *testing.T) {
		curve := CubBez(Pt(-0.9589141136356627, 7.897593105795719), Pt(-1.8683972828661322, -3.059377200831941), Pt(-2.8179811580914427, 1.4316269269607185), Pt(5.050910566843976, -0.827377732113157))
		s := Seg(Pt(-2.030576000305963, 0.14222560539036552), Pt(1.7657750606089917, -0.1459449398840639))

		geomtest.AssertVertices(t, curve.IntersectionSegment(s), []Point[float64]{Pt(-0.21238622718481381, 0.0042126716946197336), Pt(1.5998483513727035, -0.13334990072435851)})
	})
	t.Run("far from the origin a straight curve along the segment is met at its Start alone", func(t *testing.T) {
		for _, offset := range []float64{1e10, 1e12, 1e15} {
			start := Pt(offset, offset)
			curve := CubBez(start, start.Add(Vec(333.3, 111.1)), start.Add(Vec(666.7, 222.2)), start.Add(Vec(1000.0, 333.3)))

			geomtest.AssertVertices(t, curve.IntersectionSegment(Seg(start.Add(Vec(-1.0, -1.0/3)), start.Add(Vec(1001.0, 1001.0/3)))), []Point[float64]{start})
		}
	})
	t.Run("far from the origin a curve is crossed as often as near it", func(t *testing.T) {
		far := Vec(1e10, -1e10)

		for _, curve := range cubicFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)

					assert.Length(t, curve.Translate(far).IntersectionSegment(s.Translate(far)), len(curve.IntersectionSegment(s)), fmt.Sprintf("%s → %s", curve, s))
				}
			}
		}
	})
	t.Run("a segment that crosses the flattened curve well within it crosses the curve", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			path := curve.Polyline(1e-6)

			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					if path.DistanceTo(start) < 1e-3 || path.DistanceTo(end) < 1e-3 || s.IsPoint() {
						continue
					}

					for _, crossing := range s.IntersectionPolyline(path) {
						assert.True(t, slices.ContainsFunc(curve.IntersectionSegment(s), func(point Point[float64]) bool {
							return point.DistanceTo(crossing) <= 1e-3
						}), fmt.Sprintf("%s → %s: %s", curve, s, crossing))
					}
				}
			}
		}
	})
}

func TestCubicBezier_IntersectionRay(t *testing.T) {
	wave := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0))

	geomtest.AssertVertices(t, wave.IntersectionRay(Ry(Pt(1.0, 0.0), Vec(1.0, 0.0))), []Point[float64]{Pt(3.0, 0.0), Pt(6.0, 0.0)})
	assert.Nil(t, wave.IntersectionRay(Ry(Pt(-1.0, 0.0), Vec(-1.0, 0.0))))
}

func TestCubicBezier_IntersectionLine(t *testing.T) {
	wave := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0))

	geomtest.AssertVertices(t, wave.IntersectionLine(Ln(Pt(1.0, 0.0), Vec(1.0, 0.0))), []Point[float64]{Pt(0.0, 0.0), Pt(3.0, 0.0), Pt(6.0, 0.0)})
	assert.Nil(t, wave.IntersectionLine(Ln(Pt(1.0, 3.0), Vec(1.0, 0.0))))
}

func TestCubicBezier_Equal(t *testing.T) {
	t.Run("the same points", func(t *testing.T) {
		assert.True(t, archCurve().Equal(archCurve()))
	})
	t.Run("a different point", func(t *testing.T) {
		assert.False(t, archCurve().Equal(archCurve().Translate(Vec(1.0, 0.0))))
		assert.False(t, archCurve().Equal(CubBez(Pt(0.0, 0.0), Pt(0.0, 4.0), Pt(8.0, 5.0), Pt(8.0, 0.0))))
	})
	t.Run("the reversed curve is another value", func(t *testing.T) {
		assert.False(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 5), Pt(9, 0)).Equal(CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 5), Pt(9, 0)).Reverse()))
	})
}

func TestCubicBezier_IsZero(t *testing.T) {
	assert.True(t, CubicBezier[int]{}.IsZero())
	assert.False(t, CubBez(Pt(0, 0), Pt(0, 0), Pt(0, 1), Pt(0, 0)).IsZero())
}

func TestCubicBezier_Polyline(t *testing.T) {
	t.Run("a straight curve is its two ends", func(t *testing.T) {
		geomtest.AssertPolyline(t, CubBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(2.0, 2.0), Pt(3.0, 3.0)).Polyline(0.01), PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(3.0, 3.0)}))
	})
	t.Run("a point is itself", func(t *testing.T) {
		geomtest.AssertPolyline(t, CubBez(Pt(1, 2), Pt(1, 2), Pt(1, 2), Pt(1, 2)).Polyline(0.01), PolLn([]Point[int]{Pt(1, 2)}))
	})
	t.Run("a finer tolerance gives more points", func(t *testing.T) {
		assert.True(t, len(archCurve().Polyline(0.01).Points) > len(archCurve().Polyline(0.5).Points))
		assert.Equal(t, len(archCurve().Polyline(-0.5).Points), len(archCurve().Polyline(0.5).Points))
	})
	t.Run("a tolerance that cannot be met stops at the depth", func(t *testing.T) {
		assert.Length(t, archCurve().Polyline(0).Points, 1025)
		assert.Length(t, archCurve().Polyline(math.NaN()).Points, 1025)
	})
	t.Run("int halves the curve as it is and leaves out a point rounding onto the one before", func(t *testing.T) {
		path := CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).Polyline(0.25)

		geomtest.AssertPoint(t, path.Points[0], Pt(0, 0))
		geomtest.AssertPoint(t, path.Points[len(path.Points)-1], Pt(8, 0))
		for i := 1; i < len(path.Points); i++ {
			assert.True(t, path.Points[i] != path.Points[i-1], path.String())
			assert.True(t, LessOrEqual(archCurve().Polyline(0.25).DistanceTo(path.Points[i].Float()), 1), path.String())
		}
	})
	t.Run("the path runs from end to end within the tolerance of the curve", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			for _, tolerance := range []float64{0.5, 0.05, 0.001} {
				path := curve.Polyline(tolerance)
				message := fmt.Sprintf("%s ±%v: ", curve, tolerance)

				assert.Equal(t, path.Points[0], curve.Start, message)
				assert.True(t, path.Points[len(path.Points)-1].Equal(curve.End), message)
				for i := range 257 {
					point := curve.PointAt(float64(i) / 256)

					assert.True(t, LessOrEqual(path.DistanceTo(point), tolerance), message+point.String())
				}
				for _, point := range path.Points {
					assert.True(t, curve.Bounds().Contains(point), message+point.String())
				}
			}
		}
	})
	t.Run("the length of the path approaches that of the curve", func(t *testing.T) {
		// a quarter circle of radius 4 as a cubic, the kappa approximation
		const kappa = 0.5522847498

		quarter := CubBez(Pt(4.0, 0.0), Pt(4.0, 4*kappa), Pt(4*kappa, 4.0), Pt(0.0, 4.0))

		assert.True(t, math.Abs(quarter.Polyline(1e-4).Length()-2*Pi) < 1e-2)
	})
}

func TestCubicBezier_AppendPolyline(t *testing.T) {
	t.Run("appends after the points already there", func(t *testing.T) {
		points := archCurve().AppendPolyline([]Point[float64]{Pt(9.0, 9.0)}, 0.5)

		geomtest.AssertPoint(t, points[0], Pt(9.0, 9.0))
		geomtest.AssertVertices(t, points[1:], archCurve().Polyline(0.5).Points)
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		curve, buffer := CubBez(Pt(0, 0), Pt(0, 40), Pt(80, 40), Pt(80, 0)), make([]Point[int], 0, 64)

		assertAllocations(t, func() {
			sinkPoints = curve.Polyline(tenth).Points
		}, 1)
		assertAllocations(t, func() {
			sinkPoints = curve.AppendPolyline(buffer[:0], 1)
		}, 0)
	})
}

func BenchmarkCubicBezier_AppendPolyline(b *testing.B) {
	curve, buffer := CubBez(Pt(0, 0), Pt(0, 40), Pt(80, 40), Pt(80, 0)), make([]Point[int], 0, 64)

	for b.Loop() {
		sinkPoints = curve.AppendPolyline(buffer[:0], 1)
	}
}

func TestCubicBezier_Cast(t *testing.T) {
	curve := CubBez(Pt(1.5, -2.5), Pt(3.5, 1.5), Pt(0.4, 0.6), Pt(7.5, 8.5))

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, curve.Cast[int](), curve.Int())
		geomtest.AssertCubicBezier(t, curve.Cast[float64](), curve.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, curve.Cast[int8](), CubBez(Pt[int8](2, -3), Pt[int8](4, 2), Pt[int8](0, 1), Pt[int8](8, 9)))
	})
}

func TestCubicBezier_Int(t *testing.T) {
	geomtest.AssertCubicBezier(t, CubBez(Pt(1.5, -2.5), Pt(3.5, 1.5), Pt(0.4, 0.6), Pt(7.5, 8.5)).Int(), CubBez(Pt(2, -3), Pt(4, 2), Pt(0, 1), Pt(8, 9)))
}

func TestCubicBezier_Float(t *testing.T) {
	geomtest.AssertCubicBezier(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).Float(), archCurve())
}

func TestCubicBezier_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).String(), "CubBez((0,0);(0,4);(8,4);(8,0))")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)).String)
	})
}

func TestCubicBezier_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)), `{"s":{"x":0,"y":0},"c":{"x":0,"y":4},"d":{"x":8,"y":4},"e":{"x":8,"y":0}}`)

		var curve CubicBezier[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":0,"y":0},"c":{"x":0,"y":4},"d":{"x":8,"y":4},"e":{"x":8,"y":0}}`), &curve))
		geomtest.AssertCubicBezier(t, curve, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			data, err := json.Marshal(curve)
			assert.NoError(t, err)

			var decoded CubicBezier[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, curve)
		}
	})
}

func TestCubicBezier_Properties(t *testing.T) {
	t.Run("the bounds hold every point of the curve and are no larger than those of its control points", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			bounds, hull := curve.Bounds(), Pol([]Point[float64]{curve.Start, curve.Control1, curve.Control2, curve.End}).Bounds()

			assert.True(t, hull.ContainsBox(bounds), curve.String())
			for i := range 257 {
				assert.True(t, bounds.Contains(curve.PointAt(float64(i)/256)), curve.String())
			}
		}
	})
	t.Run("the bounds are reached on every side", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			bounds, reached := curve.Bounds(), Pol(curve.Polyline(1e-6).Points).Bounds()

			assert.True(t, bounds.Min.DistanceTo(reached.Min) < 1e-4 && bounds.Max.DistanceTo(reached.Max) < 1e-4, fmt.Sprintf("%s: %s ≠ %s", curve, bounds, reached))
		}
	})
	t.Run("int bounds hold the curve", func(t *testing.T) {
		for _, fixture := range cubicFixtures {
			curve := fixture.Transform(ScaleMatrix(3.0, 3.0)).Int()
			bounds := curve.Bounds()

			for i := range 65 {
				assert.True(t, bounds.Float().Contains(curve.Float().PointAt(float64(i)/64)), curve.String())
			}
		}
	})
}

func TestCubicBezier_Immutable(t *testing.T) {
	curve := CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0))

	curve.Translate(Vec(3, -2))
	curve.Transform(ScaleMatrix(2.0, 2.0))
	curve.Lerp(CubBez(Pt(2, 2), Pt(4, 0), Pt(10, 8), Pt(12, 4)), 0.5)
	curve.Reverse()
	curve.Split(0.5)

	geomtest.AssertCubicBezier(t, curve, CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)))
}

// archCurve returns the curve that rises from the origin and comes back down eight along X,
// its control points above its ends.
func archCurve() CubicBezier[float64] {
	return CubBez(Pt(0.0, 0.0), Pt(0.0, 4.0), Pt(8.0, 4.0), Pt(8.0, 0.0))
}

// curveFractions are the parameters the properties of a curve are checked at.
var curveFractions = []float64{0, 0.1, 0.35, 0.5, 0.8, 1}

// cubicFixtures span a point, a straight curve, an arch, an S, a loop crossing itself, a curve
// closing on its start and one off the origin.
var cubicFixtures = []CubicBezier[float64]{
	CubBez(Pt(1.0, 2.0), Pt(1.0, 2.0), Pt(1.0, 2.0), Pt(1.0, 2.0)),
	CubBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(2.0, 2.0), Pt(3.0, 3.0)),
	CubBez(Pt(0.0, 0.0), Pt(0.0, 4.0), Pt(8.0, 4.0), Pt(8.0, 0.0)),
	CubBez(Pt(-3.5, 0.25), Pt(2.0, 6.0), Pt(1.0, -6.0), Pt(7.0, 1.0)),
	CubBez(Pt(0.0, 0.0), Pt(10.0, 6.0), Pt(-4.0, 6.0), Pt(6.0, 0.0)),
	CubBez(Pt(1.0, 1.0), Pt(5.0, 1.0), Pt(1.0, 5.0), Pt(1.0, 1.0)),
	CubBez(Pt(12.5, -0.1), Pt(12.0, 3.0), Pt(-0.5, 12.75), Pt(0.6, -0.25)),
}

func ExampleCubBez() {
	fmt.Println(CubBez(Pt(0, 0), Pt(0, 4), Pt(8, 4), Pt(8, 0)))
	// Output: CubBez((0,0);(0,4);(8,4);(8,0))
}

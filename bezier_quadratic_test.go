package geom_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestQuadraticBezier_Constructor(t *testing.T) {
	geomtest.AssertQuadraticBezier(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)), QuadraticBezier[int]{Start: Pt(0, 0), Control: Pt(4, 4), End: Pt(8, 0)})
}

func TestQuadraticBezier_Bounds(t *testing.T) {
	t.Run("the curve, not its control point", func(t *testing.T) {
		geomtest.AssertBox(t, cornerCurve().Bounds(), Bx(Pt(0.0, 0.0), Pt(8.0, 2.0)))
	})
	t.Run("a curve turning back along X", func(t *testing.T) {
		geomtest.AssertBox(t, QuadBez(Pt(0.0, 0.0), Pt(6.0, 2.0), Pt(0.0, 4.0)).Bounds(), Bx(Pt(0.0, 0.0), Pt(3.0, 4.0)))
	})
	t.Run("a straight curve is the box of its ends", func(t *testing.T) {
		geomtest.AssertBox(t, QuadBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)).Bounds(), Bx(Pt(0.0, 0.0), Pt(3.0, 3.0)))
	})
	t.Run("int is rounded outward", func(t *testing.T) {
		geomtest.AssertBox(t, QuadBez(Pt(0, 0), Pt(4, 3), Pt(8, 0)).Bounds(), Bx(Pt(0, 0), Pt(8, 2)))
		geomtest.AssertBox(t, QuadBez(Pt(0, 0), Pt(4, -3), Pt(8, 0)).Bounds(), Bx(Pt(0, -2), Pt(8, 0)))
	})
}

func TestQuadraticBezier_Translate(t *testing.T) {
	geomtest.AssertQuadraticBezier(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).Translate(Vec(1, -2)), QuadBez(Pt(1, -2), Pt(5, 2), Pt(9, -2)))
}

func TestQuadraticBezier_Transform(t *testing.T) {
	t.Run("the three points go through the matrix", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).Transform(ScaleMatrix(2.0, 0.5)), QuadBez(Pt(0, 0), Pt(8, 2), Pt(16, 0)))
	})
	t.Run("a point of the transformed curve is the transformed point", func(t *testing.T) {
		matrix := TRSMatrix(3.0, -2.0, 0.7, 2.0, 0.5)

		for _, curve := range quadraticFixtures {
			for _, fraction := range curveFractions {
				geomtest.AssertPoint(t, curve.Transform(matrix).PointAt(fraction), curve.PointAt(fraction).Transform(matrix), fmt.Sprintf("%s @%v: ", curve, fraction))
			}
		}
	})
}

func TestQuadraticBezier_Lerp(t *testing.T) {
	a, b := QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)), QuadBez(Pt(2, 2), Pt(6, 0), Pt(12, 4))

	t.Run("every point moves towards the point of the other", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, a.Lerp(b, 0.5), QuadBez(Pt(1, 1), Pt(5, 2), Pt(10, 2)))
	})
	t.Run("the ends are the two curves", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, a.Lerp(b, 0), a)
		geomtest.AssertQuadraticBezier(t, a.Lerp(b, 1), b)
	})
	t.Run("a point of the curve between is between the points of the two", func(t *testing.T) {
		for i, curve := range quadraticFixtures {
			other := quadraticFixtures[(i+1)%len(quadraticFixtures)]

			for _, toward := range []float64{-0.5, 0.25, 1.5} {
				for _, fraction := range curveFractions {
					geomtest.AssertPoint(t, curve.Lerp(other, toward).PointAt(fraction), curve.PointAt(fraction).Lerp(other.PointAt(fraction), toward), fmt.Sprintf("%s → %s @%v: ", curve, other, fraction))
				}
			}
		}
	})
}

func TestQuadraticBezier_Reverse(t *testing.T) {
	t.Run("the points in the opposite order", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 1)).Reverse(), QuadBez(Pt(8, 1), Pt(4, 4), Pt(0, 0)))
	})
	t.Run("the same points the other way round", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			for _, fraction := range curveFractions {
				geomtest.AssertPoint(t, curve.Reverse().PointAt(1-fraction), curve.PointAt(fraction), fmt.Sprintf("%s @%v: ", curve, fraction))
			}
		}
	})
}

func TestQuadraticBezier_PointAt(t *testing.T) {
	t.Run("the ends", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			geomtest.AssertPoint(t, curve.PointAt(0), curve.Start, curve.String())
			geomtest.AssertPoint(t, curve.PointAt(1), curve.End, curve.String())
		}
	})
	t.Run("the middle, halfway to the control point", func(t *testing.T) {
		geomtest.AssertPoint(t, cornerCurve().PointAt(0.5), Pt(4.0, 2.0))
	})
	t.Run("the Bernstein form", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			for _, t1 := range []float64{-0.5, 0.1, 0.35, 0.5, 0.8, 1.5} {
				u := 1 - t1
				expected := Pt(
					u*u*curve.Start.X+2*u*t1*curve.Control.X+t1*t1*curve.End.X,
					u*u*curve.Start.Y+2*u*t1*curve.Control.Y+t1*t1*curve.End.Y,
				)

				geomtest.AssertPoint(t, curve.PointAt(t1), expected, fmt.Sprintf("%s @%v: ", curve, t1))
			}
		}
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, QuadBez(Pt(0, 0), Pt(4, 3), Pt(7, 0)).PointAt(0.5), Pt(4, 2))
	})
}

func TestQuadraticBezier_TangentAt(t *testing.T) {
	t.Run("toward the control point at the start and from it at the end", func(t *testing.T) {
		geomtest.AssertVector(t, cornerCurve().TangentAt(0), Vec(8.0, 8.0))
		geomtest.AssertVector(t, cornerCurve().TangentAt(1), Vec(8.0, -8.0))
		geomtest.AssertVector(t, cornerCurve().TangentAt(0.5), Vec(8.0, 0.0))
	})
	t.Run("no direction where the control point lies on the end", func(t *testing.T) {
		assert.True(t, QuadBez(Pt(1.0, 1.0), Pt(1.0, 1.0), Pt(8.0, 0.0)).TangentAt(0).IsZero())
	})
	t.Run("the direction the curve moves in", func(t *testing.T) {
		const step = 1e-6

		for _, curve := range quadraticFixtures {
			for _, fraction := range curveFractions {
				moved := curve.PointAt(fraction + step).Subtract(curve.PointAt(fraction - step)).Divide(2 * step)

				assert.True(t, curve.TangentAt(fraction).Subtract(moved).LessOrEqual(1e-4), fmt.Sprintf("%s @%v: %v ≠ %v", curve, fraction, curve.TangentAt(fraction), moved))
			}
		}
	})
}

func TestQuadraticBezier_Split(t *testing.T) {
	t.Run("the halves of a corner", func(t *testing.T) {
		first, second := cornerCurve().Split(0.5)

		geomtest.AssertQuadraticBezier(t, first, QuadBez(Pt(0.0, 0.0), Pt(2.0, 2.0), Pt(4.0, 2.0)))
		geomtest.AssertQuadraticBezier(t, second, QuadBez(Pt(4.0, 2.0), Pt(6.0, 2.0), Pt(8.0, 0.0)))
	})
	t.Run("the two trace the curve", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
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
}

func TestQuadraticBezier_IntersectionSegment(t *testing.T) {
	curve := QuadBez(Pt(0.0, 0.0), Pt(2.0, 4.0), Pt(4.0, 0.0))

	t.Run("the two crossings of a line through the arch, from Start", func(t *testing.T) {
		geomtest.AssertVertices(t, curve.IntersectionSegment(Seg(Pt(-1.0, 1.5), Pt(5.0, 1.5))).AppendTo(nil), []Point[float64]{Pt(1.0, 1.5), Pt(3.0, 1.5)})
		geomtest.AssertVertices(t, curve.IntersectionSegment(Seg(Pt(5.0, 1.5), Pt(-1.0, 1.5))).AppendTo(nil), []Point[float64]{Pt(3.0, 1.5), Pt(1.0, 1.5)})
	})
	t.Run("one crossing, a touch at the top, and none", func(t *testing.T) {
		geomtest.AssertVertices(t, curve.IntersectionSegment(Seg(Pt(2.0, 5.0), Pt(2.0, -1.0))).AppendTo(nil), []Point[float64]{Pt(2.0, 2.0)})
		geomtest.AssertVertices(t, curve.IntersectionSegment(Seg(Pt(-1.0, 2.0), Pt(5.0, 2.0))).AppendTo(nil), []Point[float64]{Pt(2.0, 2.0)})
		assert.Nil(t, curve.IntersectionSegment(Seg(Pt(-1.0, 3.0), Pt(5.0, 3.0))).AppendTo(nil))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		geomtest.AssertVertices(t, QuadBez(Pt(0, 0), Pt(20, 40), Pt(40, 0)).IntersectionSegment(Seg(Pt(-10, 15), Pt(50, 15))).AppendTo(nil), []Point[int]{Pt(10, 15), Pt(30, 15)})
	})
	t.Run("the points are those of the curve raised to a cubic", func(t *testing.T) {
		for _, q := range quadraticFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, q.IntersectionSegment(s).AppendTo(nil), q.Cubic().IntersectionSegment(s), fmt.Sprintf("%s → %s: ", q, s))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-1.0, 1.5), Pt(5.0, 1.5+tenth))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = curve.IntersectionSegment(through).Count > 0
		}), 0)
	})
}

func TestQuadraticBezier_IntersectionRay(t *testing.T) {
	curve := QuadBez(Pt(0.0, 0.0), Pt(2.0, 4.0), Pt(4.0, 0.0))

	geomtest.AssertVertices(t, curve.IntersectionRay(Ry(Pt(-1.0, 1.5), Vec(1.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(1.0, 1.5), Pt(3.0, 1.5)})
	geomtest.AssertVertices(t, curve.IntersectionRay(Ry(Pt(2.0, 1.5), Vec(1.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(3.0, 1.5)})
	assert.Nil(t, curve.IntersectionRay(Ry(Pt(-1.0, 1.5), Vec(-1.0, 0.0))).AppendTo(nil))
}

func TestQuadraticBezier_IntersectionLine(t *testing.T) {
	curve := QuadBez(Pt(0.0, 0.0), Pt(2.0, 4.0), Pt(4.0, 0.0))

	geomtest.AssertVertices(t, curve.IntersectionLine(Ln(Pt(2.0, 1.5), Vec(1.0, 0.0))).AppendTo(nil), []Point[float64]{Pt(1.0, 1.5), Pt(3.0, 1.5)})
	assert.Nil(t, curve.IntersectionLine(Ln(Pt(2.0, 3.0), Vec(1.0, 0.0))).AppendTo(nil))
}

func TestQuadraticBezier_Equal(t *testing.T) {
	t.Run("the same points", func(t *testing.T) {
		assert.True(t, cornerCurve().Equal(cornerCurve()))
	})
	t.Run("a different point", func(t *testing.T) {
		assert.False(t, cornerCurve().Equal(cornerCurve().Translate(Vec(1.0, 0.0))))
		assert.False(t, cornerCurve().Equal(QuadBez(Pt(0.0, 0.0), Pt(4.0, 5.0), Pt(8.0, 0.0))))
		assert.False(t, cornerCurve().Equal(QuadBez(Pt(0.0, 0.0), Pt(4.0, 4.0), Pt(8.0, 1.0))))
	})
}

func TestQuadraticBezier_IsZero(t *testing.T) {
	assert.True(t, QuadraticBezier[int]{}.IsZero())
	assert.False(t, QuadBez(Pt(0, 0), Pt(0, 1), Pt(0, 0)).IsZero())
}

func TestQuadraticBezier_Cubic(t *testing.T) {
	t.Run("the control points two thirds of the way to the one", func(t *testing.T) {
		geomtest.AssertCubicBezier(t, QuadBez(Pt(0.0, 0.0), Pt(3.0, 6.0), Pt(9.0, 0.0)).Cubic(), CubBez(Pt(0.0, 0.0), Pt(2.0, 4.0), Pt(5.0, 4.0), Pt(9.0, 0.0)))
	})
	t.Run("the cubic traces the curve", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			for _, fraction := range curveFractions {
				geomtest.AssertPoint(t, curve.Cubic().PointAt(fraction), curve.PointAt(fraction), fmt.Sprintf("%s @%v: ", curve, fraction))
			}
		}
	})
}

func TestQuadraticBezier_Polyline(t *testing.T) {
	t.Run("a straight curve is its two ends", func(t *testing.T) {
		geomtest.AssertPolyline(t, QuadBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)).Polyline(0.01), PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(3.0, 3.0)}))
	})
	t.Run("a finer tolerance gives more points", func(t *testing.T) {
		assert.True(t, len(cornerCurve().Polyline(0.01).Points) > len(cornerCurve().Polyline(0.5).Points))
	})
	t.Run("int halves the curve as it is, not as its rounded cubic", func(t *testing.T) {
		path := QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).Polyline(0.1)

		for _, point := range path.Points {
			assert.True(t, LessOrEqual(cornerCurve().Polyline(0.1).DistanceTo(point.Float()), 0.75), point.String())
		}
	})
	t.Run("the path runs from end to end within the tolerance of the curve", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			for _, tolerance := range []float64{0.5, 0.05, 0.001} {
				path := curve.Polyline(tolerance)
				message := fmt.Sprintf("%s ±%v: ", curve, tolerance)

				assert.Equal(t, path.Points[0], curve.Start, message)
				assert.True(t, path.Points[len(path.Points)-1].Equal(curve.End), message)
				for i := range 257 {
					point := curve.PointAt(float64(i) / 256)

					assert.True(t, LessOrEqual(path.DistanceTo(point), tolerance), message+point.String())
				}
			}
		}
	})
}

func TestQuadraticBezier_AppendPolyline(t *testing.T) {
	t.Run("appends after the points already there", func(t *testing.T) {
		points := cornerCurve().AppendPolyline([]Point[float64]{Pt(9.0, 9.0)}, 0.5)

		geomtest.AssertPoint(t, points[0], Pt(9.0, 9.0))
		geomtest.AssertVertices(t, points[1:], cornerCurve().Polyline(0.5).Points)
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		curve, buffer := QuadBez(Pt(0, 0), Pt(40, 40), Pt(80, 0)), make([]Point[int], 0, 64)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = curve.Polyline(tenth).Points
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = curve.AppendPolyline(buffer[:0], 1)
		}), 0)
	})
}

func TestQuadraticBezier_Cast(t *testing.T) {
	curve := QuadBez(Pt(1.5, -2.5), Pt(3.5, 1.5), Pt(7.5, 8.5))

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, curve.Cast[int](), curve.Int())
		geomtest.AssertQuadraticBezier(t, curve.Cast[float64](), curve.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertQuadraticBezier(t, curve.Cast[int8](), QuadBez(Pt[int8](2, -3), Pt[int8](4, 2), Pt[int8](8, 9)))
	})
}

func TestQuadraticBezier_Int(t *testing.T) {
	geomtest.AssertQuadraticBezier(t, QuadBez(Pt(1.5, -2.5), Pt(3.5, 1.5), Pt(7.5, 8.5)).Int(), QuadBez(Pt(2, -3), Pt(4, 2), Pt(8, 9)))
}

func TestQuadraticBezier_Float(t *testing.T) {
	geomtest.AssertQuadraticBezier(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).Float(), cornerCurve())
}

func TestQuadraticBezier_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).String(), "QuadBez((0,0);(4,4);(8,0))")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)).String)
	})
}

func TestQuadraticBezier_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)), `{"s":{"x":0,"y":0},"c":{"x":4,"y":4},"e":{"x":8,"y":0}}`)

		var curve QuadraticBezier[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":0,"y":0},"c":{"x":4,"y":4},"e":{"x":8,"y":0}}`), &curve))
		geomtest.AssertQuadraticBezier(t, curve, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			data, err := json.Marshal(curve)
			assert.NoError(t, err)

			var decoded QuadraticBezier[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, curve)
		}
	})
}

func TestQuadraticBezier_Properties(t *testing.T) {
	t.Run("the bounds hold every point of the curve and are no larger than those of its control points", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			bounds, hull := curve.Bounds(), Pol([]Point[float64]{curve.Start, curve.Control, curve.End}).Bounds()

			assert.True(t, hull.ContainsBox(bounds), curve.String())
			for i := range 257 {
				assert.True(t, bounds.Contains(curve.PointAt(float64(i)/256)), curve.String())
			}
		}
	})
	t.Run("the bounds are those of the cubic over the same points", func(t *testing.T) {
		for _, curve := range quadraticFixtures {
			geomtest.AssertBox(t, curve.Bounds(), curve.Cubic().Bounds(), curve.String())
		}
	})
}

func TestQuadraticBezier_Immutable(t *testing.T) {
	curve := QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0))

	curve.Translate(Vec(3, -2))
	curve.Transform(ScaleMatrix(2.0, 2.0))
	curve.Lerp(QuadBez(Pt(2, 2), Pt(6, 0), Pt(12, 4)), 0.5)
	curve.Reverse()
	curve.Split(0.5)

	geomtest.AssertQuadraticBezier(t, curve, QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)))
}

// cornerCurve returns the curve that rounds the corner above the middle of its two ends, eight
// apart along X.
func cornerCurve() QuadraticBezier[float64] {
	return QuadBez(Pt(0.0, 0.0), Pt(4.0, 4.0), Pt(8.0, 0.0))
}

// quadraticFixtures span a point, a straight curve, a corner, a curve turning back along X, a
// curve closing on its start and one off the origin.
var quadraticFixtures = []QuadraticBezier[float64]{
	QuadBez(Pt(1.0, 2.0), Pt(1.0, 2.0), Pt(1.0, 2.0)),
	QuadBez(Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)),
	QuadBez(Pt(0.0, 0.0), Pt(4.0, 4.0), Pt(8.0, 0.0)),
	QuadBez(Pt(0.0, 0.0), Pt(6.0, 2.0), Pt(0.0, 4.0)),
	QuadBez(Pt(1.0, 1.0), Pt(5.0, 3.0), Pt(1.0, 1.0)),
	QuadBez(Pt(12.5, -0.1), Pt(-0.5, 12.75), Pt(0.6, -0.25)),
}

func ExampleQuadBez() {
	fmt.Println(QuadBez(Pt(0, 0), Pt(4, 4), Pt(8, 0)))
	// Output: QuadBez((0,0);(4,4);(8,0))
}

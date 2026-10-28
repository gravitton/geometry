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

func TestPolyline_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()), Polyline[int]{Points: []Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 3)}})
	})
	t.Run("the slice is shared", func(t *testing.T) {
		points := cornerPath()

		assert.Equal(t, &PolLn(points).Points[0], &points[0])
	})
}

func TestPolyline_Vertices(t *testing.T) {
	geomtest.AssertVertices(t, slices.Collect(PolLn(cornerPath()).Vertices()), cornerPath())
}

func TestPolyline_Edges(t *testing.T) {
	t.Run("one fewer than the points, none returning to the first", func(t *testing.T) {
		geomtest.AssertSegments(t, slices.Collect(PolLn(cornerPath()).Edges()), []Segment[int]{Seg(Pt(0, 0), Pt(4, 0)), Seg(Pt(4, 0), Pt(4, 3))})
	})
	t.Run("a point is one edge of no length", func(t *testing.T) {
		geomtest.AssertSegments(t, slices.Collect(PolLn([]Point[int]{Pt(1, 2)}).Edges()), []Segment[int]{Seg(Pt(1, 2), Pt(1, 2))})
	})
	t.Run("an empty path has none", func(t *testing.T) {
		assert.Length(t, slices.Collect(Polyline[int]{}.Edges()), 0)
	})
	t.Run("stops when the loop does", func(t *testing.T) {
		for _, path := range []Polyline[int]{PolLn(cornerPath()), PolLn([]Point[int]{Pt(1, 2)})} {
			count := 0
			for range path.Edges() {
				count++

				break
			}

			assert.Equal(t, count, 1)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		path := PolLn(cornerPath())

		assertAllocations(t, func() {
			for edge := range path.Edges() {
				sinkBool = edge.Start.X > 0
			}
		}, 0)
	})
}

func TestPolyline_Length(t *testing.T) {
	t.Run("the sum of the edges", func(t *testing.T) {
		geomtest.AssertNumber(t, PolLn(cornerPath()).Length(), 7.0)
	})
	t.Run("a point and an empty path have none", func(t *testing.T) {
		geomtest.AssertNumber(t, PolLn([]Point[int]{Pt(1, 2)}).Length(), 0.0)
		geomtest.AssertNumber(t, Polyline[int]{}.Length(), 0.0)
	})
	t.Run("a closed path is as long as the perimeter of its polygon", func(t *testing.T) {
		geomtest.AssertNumber(t, Pol(squareVertices()).Polyline().Length(), 8.0)
	})
}

func TestPolyline_Bounds(t *testing.T) {
	t.Run("the box of the points", func(t *testing.T) {
		geomtest.AssertBox(t, PolLn(cornerPath()).Bounds(), Bx(Pt(0, 0), Pt(4, 3)))
	})
	t.Run("an empty path has the zero box", func(t *testing.T) {
		geomtest.AssertBox(t, Polyline[int]{}.Bounds(), Box[int]{})
	})
}

func TestPolyline_Translate(t *testing.T) {
	t.Run("every point moves", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Translate(Vec(1, -2)), PolLn([]Point[int]{Pt(1, -2), Pt(5, -2), Pt(5, 1)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Polyline[int]{}.Translate(Vec(1, -2)).IsZero())
	})
}

func TestPolyline_Transform(t *testing.T) {
	t.Run("every point goes through the matrix", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Transform(ScaleMatrix(2.0, 3.0)), PolLn([]Point[int]{Pt(0, 0), Pt(8, 0), Pt(8, 9)}))
	})
	t.Run("a turn about a pivot", func(t *testing.T) {
		turn := TranslationMatrix(4.0, 0.0).Rotate(Pi/2).Translate(-4.0, 0.0)

		geomtest.AssertPolyline(t, PolLn(cornerPath()).Float().Transform(turn), PolLn([]Point[float64]{Pt(4.0, -4.0), Pt(4.0, 0.0), Pt(1.0, 0.0)}))
	})
}

func TestPolyline_Lerp(t *testing.T) {
	a := PolLn(cornerPath()).Float()
	b := PolLn([]Point[float64]{Pt(10.0, 0.0), Pt(12.0, 0.0), Pt(12.0, 7.0)})

	t.Run("each point moves towards the point of its index", func(t *testing.T) {
		geomtest.AssertPolyline(t, a.Lerp(b, 0.5), PolLn([]Point[float64]{Pt(5.0, 0.0), Pt(8.0, 0.0), Pt(8.0, 5.0)}))
	})
	t.Run("the ends are the two paths", func(t *testing.T) {
		geomtest.AssertPolyline(t, a.Lerp(b, 0), a)
		geomtest.AssertPolyline(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates outside [0, 1]", func(t *testing.T) {
		geomtest.AssertPolyline(t, a.Lerp(b, 2), PolLn([]Point[float64]{Pt(20.0, 0.0), Pt(20.0, 0.0), Pt(20.0, 11.0)}))
	})
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Lerp(PolLn([]Point[int]{Pt(1, 1), Pt(5, 1), Pt(5, 4)}), 0.5), PolLn([]Point[int]{Pt(1, 1), Pt(5, 1), Pt(5, 4)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Polyline[int]{}.Lerp(Polyline[int]{}, 0.5).IsZero())
	})
	t.Run("a different point count panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			a.Lerp(PolLn([]Point[float64]{Pt(1.0, 2.0)}), 0.5)
		}, ErrCountMismatch)
	})
	t.Run("towards a translation is a part of it over the fixtures", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			moved := path.Translate(Vec(3.0, -2.0))

			assert.True(t, path.Lerp(moved, 0.25).Equal(path.Translate(Vec(0.75, -0.5))), path.String())
		}
	})
}

func TestPolyline_Reverse(t *testing.T) {
	t.Run("the points in the opposite order", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Reverse(), PolLn([]Point[int]{Pt(4, 3), Pt(4, 0), Pt(0, 0)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Polyline[int]{}.Reverse().IsZero())
		assert.False(t, PolLn([]Point[int]{}).Reverse().IsZero())
	})
}

func TestPolyline_PointAt(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("the ends", func(t *testing.T) {
		geomtest.AssertPoint(t, path.PointAt(0), Pt(0.0, 0.0))
		geomtest.AssertPoint(t, path.PointAt(1), Pt(4.0, 3.0))
	})
	t.Run("a fraction of the length, around the corner", func(t *testing.T) {
		geomtest.AssertPoint(t, path.PointAt(0.5), Pt(3.5, 0.0))
		geomtest.AssertPoint(t, path.PointAt(4.0/7), Pt(4.0, 0.0))
		geomtest.AssertPoint(t, path.PointAt(6.0/7), Pt(4.0, 2.0))
	})
	t.Run("outside the path it runs on along the first edge and the last", func(t *testing.T) {
		geomtest.AssertPoint(t, path.PointAt(-1.0/7), Pt(-1.0, 0.0))
		geomtest.AssertPoint(t, path.PointAt(8.0/7), Pt(4.0, 4.0))
	})
	t.Run("a repeated point is passed over", func(t *testing.T) {
		repeated := PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 0.0)})

		geomtest.AssertPoint(t, repeated.PointAt(-0.25), Pt(-1.0, 0.0))
		geomtest.AssertPoint(t, repeated.PointAt(0.5), Pt(2.0, 0.0))
		geomtest.AssertPoint(t, repeated.PointAt(1), Pt(4.0, 0.0))
		geomtest.AssertPoint(t, repeated.PointAt(1.25), Pt(5.0, 0.0))
	})
	t.Run("points all alike are that point", func(t *testing.T) {
		geomtest.AssertPoint(t, PolLn([]Point[int]{Pt(1, 2), Pt(1, 2), Pt(1, 2)}).PointAt(1.5), Pt(1, 2))
	})
	t.Run("a path of two points answers as the segment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, fraction := range []float64{-0.5, 0, 0.25, 1, 1.5} {
				geomtest.AssertPoint(t, PolLn([]Point[float64]{s.Start, s.End}).PointAt(fraction), s.PointAt(fraction), fmt.Sprintf("%s @%v: ", s, fraction))
			}
		}
	})
	t.Run("a point is itself and an empty path the zero point", func(t *testing.T) {
		geomtest.AssertPoint(t, PolLn([]Point[int]{Pt(1, 2)}).PointAt(0.5), Pt(1, 2))
		geomtest.AssertPoint(t, Polyline[int]{}.PointAt(0.5), Pt(0, 0))
	})
	t.Run("int is rounded once", func(t *testing.T) {
		geomtest.AssertPoint(t, PolLn(cornerPath()).PointAt(0.5), Pt(4, 0))
	})
}

func TestPolyline_Every(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("a point every step from the first, around the corner", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(path.Every(2.5)), []Point[float64]{Pt(0.0, 0.0), Pt(2.5, 0.0), Pt(4.0, 1.0)})
	})
	t.Run("the last point where the steps reach it", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(path.Every(3.5)), []Point[float64]{Pt(0.0, 0.0), Pt(3.5, 0.0), Pt(4.0, 3.0)})
		assert.Length(t, slices.Collect(path.Every(1)), 8)
	})
	t.Run("each point is at its count of steps, the rounding of one not carried to the next", func(t *testing.T) {
		points := slices.Collect(PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 0.0)}).Every(tenth))

		assert.Length(t, points, 11)
		assert.Equal(t, points[10], Pt(1.0, 0.0))
	})
	t.Run("a step that is not positive yields the first point alone", func(t *testing.T) {
		for _, step := range []float64{0, -1, math.NaN()} {
			geomtest.AssertVertices(t, slices.Collect(path.Every(step)), []Point[float64]{Pt(0.0, 0.0)})
		}
	})
	t.Run("a point yields itself and an empty path nothing", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(PolLn([]Point[int]{Pt(1, 2)}).Every(1)), []Point[int]{Pt(1, 2)})
		assert.Length(t, slices.Collect(Polyline[int]{}.Every(1)), 0)
	})
	t.Run("the walk ends where the length of the path stops being finite", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 0.0), Pt(math.Inf(1), 0.0)}).Every(0.5)), []Point[float64]{Pt(0.0, 0.0), Pt(0.5, 0.0), Pt(1.0, 0.0)})
		assert.Length(t, slices.Collect(PolLn([]Point[float64]{Pt(-1e308, 0.0), Pt(1e308, 0.0)}).Every(1e307)), 0)
	})
	t.Run("stops when the loop does", func(t *testing.T) {
		count := 0
		for range path.Every(1) {
			count++
			if count == 3 {
				break
			}
		}

		assert.Equal(t, count, 3)
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assertAllocations(t, func() {
			for point := range path.Every(tenth) {
				sinkBool = point.X > 0
			}
		}, 0)
	})
}

func TestPolyline_Simplify(t *testing.T) {
	t.Run("drops a point on the line between its neighbours and a repeated one", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(4, 0), Pt(4, 0), Pt(4, 3)}).Simplify(0), PolLn(cornerPath()))
	})
	t.Run("drops a point within the tolerance", func(t *testing.T) {
		bumped := PolLn([]Point[int]{Pt(0, 0), Pt(5, 1), Pt(10, 0), Pt(10, 10)})

		geomtest.AssertPolyline(t, bumped.Simplify(1), PolLn([]Point[int]{Pt(0, 0), Pt(10, 0), Pt(10, 10)}))
		geomtest.AssertPolyline(t, bumped.Simplify(-1), PolLn([]Point[int]{Pt(0, 0), Pt(10, 0), Pt(10, 10)}))
		geomtest.AssertPolyline(t, bumped.Simplify(0.5), bumped)
	})
	t.Run("the ends are kept whatever the tolerance", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Simplify(100), PolLn([]Point[int]{Pt(0, 0), Pt(4, 3)}))
	})
	t.Run("a closed path keeps the point farthest from where it closes", func(t *testing.T) {
		geomtest.AssertPolyline(t, Pol(squareVertices()).Polyline().Simplify(1.5), PolLn([]Point[int]{Pt(0, 0), Pt(2, 2), Pt(0, 0)}))
	})
	t.Run("a point and an empty path are returned as they are", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn([]Point[int]{Pt(1, 2)}).Simplify(1), PolLn([]Point[int]{Pt(1, 2)}))
		assert.True(t, Polyline[int]{}.Simplify(1).IsZero())
	})
	t.Run("the kept points are points of the path, and simplifying again changes nothing", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			for _, tolerance := range []float64{0, 0.5, 2} {
				simplified := path.Simplify(tolerance)
				message := fmt.Sprintf("%s ±%v: ", path, tolerance)

				geomtest.AssertPolyline(t, simplified.Simplify(tolerance), simplified, message)
				for _, point := range simplified.Points {
					assert.True(t, slices.Contains(path.Points, point), message+point.String())
				}
				for _, point := range path.Points {
					assert.True(t, LessOrEqual(simplified.DistanceTo(point), tolerance), message+point.String())
				}
			}
		}
	})
}

func TestPolyline_AppendSimplify(t *testing.T) {
	path := PolLn([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(4, 0), Pt(4, 3)})

	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, path.AppendSimplify([]Point[int]{Pt(9, 9)}, 0), []Point[int]{Pt(9, 9), Pt(0, 0), Pt(4, 0), Pt(4, 3)})
	})
	t.Run("an empty path appends nothing", func(t *testing.T) {
		geomtest.AssertVertices(t, Polyline[int]{}.AppendSimplify([]Point[int]{Pt(9, 9)}, 1), []Point[int]{Pt(9, 9)})
		assert.Nil(t, Polyline[int]{}.AppendSimplify(nil, 1))
	})
	t.Run("allocates once, and nothing into a buffer with room", func(t *testing.T) {
		buffer := make([]Point[int], 0, 8)

		assertAllocations(t, func() {
			sinkPoints = path.Simplify(0).Points
		}, 1)
		assertAllocations(t, func() {
			sinkPoints = path.AppendSimplify(buffer[:0], 0)
		}, 0)
	})
}

func TestPolyline_DistanceTo(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("zero on the path", func(t *testing.T) {
		geomtest.AssertNumber(t, path.DistanceTo(Pt(2.0, 0.0)), 0.0)
		geomtest.AssertNumber(t, path.DistanceTo(Pt(4.0, 3.0)), 0.0)
	})
	t.Run("to the nearest edge", func(t *testing.T) {
		geomtest.AssertNumber(t, path.DistanceTo(Pt(2.0, 1.5)), 1.5)
		geomtest.AssertNumber(t, path.DistanceTo(Pt(3.0, 2.0)), 1.0)
	})
	t.Run("past an end, to the end", func(t *testing.T) {
		geomtest.AssertNumber(t, path.DistanceTo(Pt(-3.0, 4.0)), 5.0)
	})
	t.Run("an empty path is infinitely far", func(t *testing.T) {
		assert.True(t, math.IsInf(Polyline[float64]{}.DistanceTo(Pt(1.0, 2.0)), 1))
	})
}

func TestPolyline_DistanceSquaredTo(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("the square of the distance", func(t *testing.T) {
		geomtest.AssertNumber(t, path.DistanceSquaredTo(Pt(-3.0, 4.0)), 25.0)
	})
	t.Run("a point within the tolerance of an edge is on it", func(t *testing.T) {
		assert.Equal(t, path.DistanceSquaredTo(Pt(2.0, Delta/2)), 0.0)
	})
	t.Run("int has no tolerance", func(t *testing.T) {
		geomtest.AssertNumber(t, PolLn(cornerPath()).DistanceSquaredTo(Pt(2, 1)), 1.0)
	})
}

func TestPolyline_Nearest(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("a point of the path is itself", func(t *testing.T) {
		geomtest.AssertPoint(t, path.Nearest(Pt(2.0, 0.0)), Pt(2.0, 0.0))
	})
	t.Run("the foot on the nearest edge", func(t *testing.T) {
		geomtest.AssertPoint(t, path.Nearest(Pt(2.0, 1.5)), Pt(2.0, 0.0))
		geomtest.AssertPoint(t, path.Nearest(Pt(3.0, 2.0)), Pt(4.0, 2.0))
	})
	t.Run("past an end, the end", func(t *testing.T) {
		geomtest.AssertPoint(t, path.Nearest(Pt(9.0, 9.0)), Pt(4.0, 3.0))
	})
	t.Run("the first edge of two as near", func(t *testing.T) {
		geomtest.AssertPoint(t, path.Nearest(Pt(3.0, 1.0)), Pt(3.0, 0.0))
	})
	t.Run("an empty path gives the zero point", func(t *testing.T) {
		geomtest.AssertPoint(t, Polyline[float64]{}.Nearest(Pt(1.0, 2.0)), Pt(0.0, 0.0))
	})
	t.Run("a NaN coordinate is on no path, and its nearest point is NaN", func(t *testing.T) {
		for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
			nearest := path.Nearest(p)

			assert.False(t, path.Contains(p), p.String())
			assert.True(t, math.IsNaN(path.DistanceTo(p)) || math.IsInf(path.DistanceTo(p), 1), p.String())
			assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y), p.String())
		}
	})
}

func TestPolyline_NearestFraction(t *testing.T) {
	path := PolLn(cornerPath()).Float()

	t.Run("how far along the nearest point lies", func(t *testing.T) {
		geomtest.AssertNumber(t, path.NearestFraction(Pt(2.0, 1.5)), 2.0/7)
		geomtest.AssertNumber(t, path.NearestFraction(Pt(6.0, 2.0)), 6.0/7)
	})
	t.Run("zero before the start and one past the end", func(t *testing.T) {
		geomtest.AssertNumber(t, path.NearestFraction(Pt(-3.0, 4.0)), 0.0)
		geomtest.AssertNumber(t, path.NearestFraction(Pt(9.0, 9.0)), 1.0)
	})
	t.Run("the first edge of two as near", func(t *testing.T) {
		geomtest.AssertNumber(t, path.NearestFraction(Pt(3.0, 1.0)), 3.0/7)
	})
	t.Run("a path that retraces itself answers on the way out, as Nearest does", func(t *testing.T) {
		retraced := PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(10.0, 3.0), Pt(0.0, 0.0)})

		for _, fraction := range []float64{0, 0.00075, 0.1, 0.25, 0.4, 0.5} {
			geomtest.AssertNumber(t, retraced.NearestFraction(retraced.PointAt(fraction)), fraction, fmt.Sprintf("@%v: ", fraction))
		}
	})
	t.Run("a path of no length and an empty one give zero", func(t *testing.T) {
		geomtest.AssertNumber(t, PolLn([]Point[int]{Pt(1, 2)}).NearestFraction(Pt(5, 5)), 0.0)
		geomtest.AssertNumber(t, Polyline[int]{}.NearestFraction(Pt(5, 5)), 0.0)
	})
	t.Run("a NaN coordinate is nearest to nothing", func(t *testing.T) {
		for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
			assert.True(t, math.IsNaN(path.NearestFraction(p)), p.String())
		}
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, PolLn(cornerPath()).NearestFraction(Pt(6, 2)), 6.0/7)
	})
}

func TestPolyline_Contains(t *testing.T) {
	t.Run("a point of an edge", func(t *testing.T) {
		assert.True(t, PolLn(cornerPath()).Contains(Pt(2, 0)))
		assert.True(t, PolLn(cornerPath()).Contains(Pt(4, 3)))
	})
	t.Run("a point off the path", func(t *testing.T) {
		assert.False(t, PolLn(cornerPath()).Contains(Pt(2, 1)))
		assert.False(t, PolLn(cornerPath()).Contains(Pt(5, 0)))
	})
	t.Run("a closed path has no inside", func(t *testing.T) {
		assert.False(t, Pol(squareVertices()).Polyline().Contains(Pt(1, 1)))
		assert.True(t, Pol(squareVertices()).Polyline().Contains(Pt(0, 1)))
	})
	t.Run("an empty path contains nothing", func(t *testing.T) {
		assert.False(t, Polyline[int]{}.Contains(Pt(0, 0)))
	})
}

func TestPolyline_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionPolyline", func(t *testing.T) {
		for _, p := range polylineFixtures() {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, p.IntersectionSegment(s), s.IntersectionPolyline(p), fmt.Sprintf("%s → %s: ", p, s))
			}
		}
	})
}

func TestPolyline_IntersectionRay(t *testing.T) {
	t.Run("matches Ray.IntersectionPolyline", func(t *testing.T) {
		for _, p := range polylineFixtures() {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, p.IntersectionRay(r), r.IntersectionPolyline(p), fmt.Sprintf("%s → %s: ", p, r))
			}
		}
	})
}

func TestPolyline_IntersectionLine(t *testing.T) {
	t.Run("matches Line.IntersectionPolyline", func(t *testing.T) {
		for _, p := range polylineFixtures() {
			for _, l := range lineFixtures {
				geomtest.AssertVertices(t, p.IntersectionLine(l), l.IntersectionPolyline(p), fmt.Sprintf("%s → %s: ", p, l))
			}
		}
	})
}

func TestPolyline_Equal(t *testing.T) {
	t.Run("the same points", func(t *testing.T) {
		assert.True(t, PolLn(cornerPath()).Equal(PolLn(cornerPath())))
	})
	t.Run("a different point or count", func(t *testing.T) {
		assert.False(t, PolLn(cornerPath()).Equal(PolLn([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4)})))
		assert.False(t, PolLn(cornerPath()).Equal(PolLn([]Point[int]{Pt(0, 0), Pt(4, 0)})))
	})
	t.Run("the reversed path is another value", func(t *testing.T) {
		assert.False(t, PolLn(cornerPath()).Equal(PolLn(cornerPath()).Reverse()))
	})
	t.Run("nil and empty are equal", func(t *testing.T) {
		assert.True(t, Polyline[int]{}.Equal(PolLn([]Point[int]{})))
	})
}

func TestPolyline_IsZero(t *testing.T) {
	assert.True(t, Polyline[int]{}.IsZero())
	assert.False(t, PolLn([]Point[int]{}).IsZero())
}

func TestPolyline_IsEmpty(t *testing.T) {
	assert.True(t, Polyline[int]{}.IsEmpty())
	assert.True(t, PolLn([]Point[int]{}).IsEmpty())
	assert.False(t, PolLn(cornerPath()).IsEmpty())
}

func TestPolyline_IsClosed(t *testing.T) {
	t.Run("the last point is the first", func(t *testing.T) {
		assert.True(t, PolLn([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(0, 0)}).IsClosed())
	})
	t.Run("an open path", func(t *testing.T) {
		assert.False(t, PolLn(cornerPath()).IsClosed())
	})
	t.Run("fewer than two points have no end to return from", func(t *testing.T) {
		assert.False(t, PolLn([]Point[int]{Pt(1, 2)}).IsClosed())
		assert.False(t, Polyline[int]{}.IsClosed())
	})
}

func TestPolyline_Polygon(t *testing.T) {
	t.Run("an open path closes back to its first point", func(t *testing.T) {
		geomtest.AssertPolygon(t, PolLn(cornerPath()).Polygon(), Pol(cornerPath()))
	})
	t.Run("a closed path gives up the repeated point", func(t *testing.T) {
		geomtest.AssertPolygon(t, PolLn([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2), Pt(0, 0)}).Polygon(), Pol(squareVertices()))
	})
	t.Run("the points are shared", func(t *testing.T) {
		points := cornerPath()

		assert.Equal(t, &PolLn(points).Polygon().Points[0], &points[0])
	})
	t.Run("appending to the polygon of a closed path leaves the path as it is", func(t *testing.T) {
		points := []Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 3), Pt(0, 0)}
		polygon := PolLn(points).Polygon()

		assert.Length(t, append(polygon.Points, Pt(9, 9)), 4)
		geomtest.AssertVertices(t, points, []Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 3), Pt(0, 0)})
	})
	t.Run("a polygon comes back from its path", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			geomtest.AssertPolygon(t, polygon.Polyline().Polygon(), polygon, polygon.String())
		}
	})
}

func TestPolyline_Cast(t *testing.T) {
	path := PolLn([]Point[float64]{Pt(1.5, -2.5), Pt(3.5, 1.5)})

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertPolyline(t, path.Cast[int](), path.Int())
		geomtest.AssertPolyline(t, path.Cast[float64](), path.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertPolyline(t, path.Cast[int8](), PolLn([]Point[int8]{Pt[int8](2, -3), Pt[int8](4, 2)}))
	})
}

func TestPolyline_Int(t *testing.T) {
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn([]Point[float64]{Pt(0.6, -0.25), Pt(1.2, 3.4)}).Int(), PolLn([]Point[int]{Pt(1, 0), Pt(1, 3)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Polyline[float64]{}.Int().IsZero())
	})
}

func TestPolyline_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertPolyline(t, PolLn(cornerPath()).Float(), PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 3.0)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Polyline[int]{}.Float().IsZero())
	})
}

func TestPolyline_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, PolLn(cornerPath()).String(), "PolLn((0,0);(4,0);(4,3))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, PolLn([]Point[float64]{Pt(0.5, -1.25)}).String(), "PolLn((0.50,-1.25))")
	})
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, Polyline[int]{}.String(), "PolLn()")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, PolLn(cornerPath()).String)
	})
}

func TestPolyline_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, PolLn(cornerPath()), `[{"x":0,"y":0},{"x":4,"y":0},{"x":4,"y":3}]`)

		var p Polyline[int]
		assert.NoError(t, json.Unmarshal([]byte(`[{"x":0,"y":0},{"x":4,"y":0},{"x":4,"y":3}]`), &p))
		geomtest.AssertPolyline(t, p, PolLn(cornerPath()))
	})
	t.Run("decoding leaves a shared slice untouched", func(t *testing.T) {
		shared := cornerPath()
		p := PolLn(shared)

		assert.NoError(t, json.Unmarshal([]byte(`[{"x":9,"y":9}]`), &p))
		geomtest.AssertVertices(t, shared, cornerPath())
		geomtest.AssertPolyline(t, p, PolLn([]Point[int]{{9, 9}}))
	})
	t.Run("invalid input leaves the path untouched", func(t *testing.T) {
		p := PolLn(cornerPath())

		assert.Error(t, json.Unmarshal([]byte(`[{"x":"nine"}]`), &p))
		geomtest.AssertPolyline(t, p, PolLn(cornerPath()))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			data, err := json.Marshal(path)
			assert.NoError(t, err)

			var decoded Polyline[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.True(t, decoded.Equal(path), path.String())
		}
	})
}

func TestPolyline_Properties(t *testing.T) {
	t.Run("the distance is the least to an edge, zero exactly where contains holds", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			for _, point := range pointFixtures {
				least := math.Inf(1)
				for edge := range path.Edges() {
					least = min(least, edge.DistanceTo(point))
				}

				message := fmt.Sprintf("%s → %s: ", path, point)

				geomtest.AssertNumber(t, path.DistanceTo(point), least, message)
				assert.Equal(t, path.DistanceTo(point) == 0, path.Contains(point), message)
			}
		}
	})
	t.Run("the nearest point is on the path, at the distance and at the nearest fraction", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			if path.IsEmpty() {
				continue
			}

			for _, point := range pointFixtures {
				nearest := path.Nearest(point)
				message := fmt.Sprintf("%s → %s: ", path, point)

				assert.Equal(t, nearest.Equal(point), path.Contains(point), message)
				assert.True(t, path.Contains(nearest), message)
				geomtest.AssertNumber(t, point.DistanceTo(nearest), path.DistanceTo(point), message)
				assert.True(t, path.PointAt(path.NearestFraction(point)).Equal(nearest), message)
			}
		}
	})
	t.Run("a point at a fraction lies on the path, at the fraction from the other end of the reversed one", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			for _, fraction := range []float64{-0.5, 0, 0.1, 0.35, 0.5, 0.8, 1, 1.5} {
				point := path.PointAt(fraction)
				message := fmt.Sprintf("%s @%v: ", path, fraction)

				assert.True(t, path.IsEmpty() || fraction < 0 || fraction > 1 || path.Contains(point), message)
				geomtest.AssertPoint(t, path.Reverse().PointAt(1-fraction), point, message)
			}
		}
	})
	t.Run("a repeated point at either end changes no point at a fraction", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			if path.IsEmpty() {
				continue
			}

			n := len(path.Points)
			repeated := PolLn(slices.Concat(path.Points[:1], path.Points, path.Points[n-1:]))

			for _, fraction := range []float64{-0.5, 0, 0.35, 1, 1.5} {
				geomtest.AssertPoint(t, repeated.PointAt(fraction), path.PointAt(fraction), fmt.Sprintf("%s @%v: ", path, fraction))
			}
		}
	})
	t.Run("the fraction of a point at a fraction is that fraction, where the path does not meet itself", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			if path.Length() == 0 || path.IsClosed() {
				continue
			}

			for _, fraction := range []float64{0, 0.1, 0.35, 0.5, 0.8, 1} {
				geomtest.AssertNumber(t, path.NearestFraction(path.PointAt(fraction)), fraction, fmt.Sprintf("%s @%v: ", path, fraction))
			}
		}
	})
	t.Run("the points a step apart lie on the path, as many as the steps that fit", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			if path.IsEmpty() {
				continue
			}

			for _, step := range []float64{0.3, 1.7, 100} {
				points := slices.Collect(path.Every(step))
				message := fmt.Sprintf("%s +%v: ", path, step)

				assert.Length(t, points, int(path.Length()/step)+1)
				for i, point := range points {
					assert.True(t, path.Contains(point), message+point.String())
					assert.True(t, point.Equal(path.PointAt(ratioOf(float64(i)*step, path.Length()))), message+point.String())
				}
			}
		}
	})
	t.Run("a reversed or moved path is as long and as far", func(t *testing.T) {
		for _, path := range polylineFixtures() {
			geomtest.AssertNumber(t, path.Reverse().Length(), path.Length(), path.String())
			geomtest.AssertNumber(t, path.Translate(Vec(3.0, -2.0)).Length(), path.Length(), path.String())

			for _, point := range pointFixtures {
				geomtest.AssertNumber(t, path.Reverse().DistanceTo(point), path.DistanceTo(point), path.String())
				geomtest.AssertNumber(t, path.Translate(Vec(3.0, -2.0)).DistanceTo(point.Add(Vec(3.0, -2.0))), path.DistanceTo(point), path.String())
			}
		}
	})
	t.Run("the path around a polygon has its edges, its perimeter and its bounds", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			path := polygon.Polyline()

			assert.True(t, path.IsClosed(), polygon.String())
			geomtest.AssertSegments(t, slices.Collect(path.Edges()), slices.Collect(polygon.Edges()), polygon.String())
			geomtest.AssertNumber(t, path.Length(), polygon.Perimeter(), polygon.String())
			geomtest.AssertBox(t, path.Bounds(), polygon.Bounds(), polygon.String())
		}
	})
	t.Run("far from the origin a float32 point at a fraction lies on the path", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, fixture := range polylineFixtures() {
				path := fixture.Cast[float32]().Translate(offset)

				for _, fraction := range []float64{0, 0.1, 0.35, 0.5, 0.8, 1} {
					assert.True(t, path.IsEmpty() || path.Contains(path.PointAt(fraction)), fmt.Sprintf("%s @%v: ", path, fraction))
				}
			}
		}
	})
	t.Run("far from the origin a float32 nearest point lies on the path, at the nearest fraction", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, fixture := range polylineFixtures() {
				path := fixture.Cast[float32]().Translate(offset)
				if path.IsEmpty() {
					continue
				}

				for _, p := range pointFixtures {
					point := p.Cast[float32]().Add(offset)
					message := fmt.Sprintf("%s → %s: ", path, point)

					assert.True(t, path.Contains(path.Nearest(point)), message)
					assert.True(t, path.Contains(path.PointAt(path.NearestFraction(point))), message)
				}
			}
		}
	})
	t.Run("int near the end of the supported range", func(t *testing.T) {
		const far = 1 << 25

		path := PolLn([]Point[int]{Pt(-far, -far), Pt(far, -far), Pt(far, far)})

		geomtest.AssertNumber(t, path.Length(), 4.0*far)
		assert.True(t, path.Contains(Pt(1, -far)))
		assert.False(t, path.Contains(Pt(1, 1-far)))
		geomtest.AssertNumber(t, path.DistanceSquaredTo(Pt(1, 1-far)), 1.0)
		geomtest.AssertPoint(t, path.Nearest(Pt(far-1, 1)), Pt(far, 1))
		geomtest.AssertNumber(t, path.NearestFraction(Pt(far-1, 0)), 0.75)
		geomtest.AssertPoint(t, path.PointAt(0.75), Pt(far, 0))
	})
}

// ratioOf returns a divided by b, or zero for a zero b, the fraction of a length a path of no
// length has.
func ratioOf(a, b float64) float64 {
	if b == 0 {
		return 0
	}

	return a / b
}

func TestPolyline_Immutable(t *testing.T) {
	points := cornerPath()
	path := PolLn(points)

	path.Translate(Vec(3, -2))
	path.Transform(ScaleMatrix(2.0, 2.0))
	path.Lerp(PolLn([]Point[int]{Pt(1, 1), Pt(5, 1), Pt(5, 4)}), 0.5)
	path.Reverse()
	path.Simplify(10)
	path.AppendSimplify(nil, 10)

	geomtest.AssertVertices(t, points, cornerPath())
}

// cornerPath returns the points of a path along X and then along Y, of whole lengths.
func cornerPath() []Point[int] {
	return []Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 3)}
}

// polylineFixtures span an empty path, a point, a segment, a corner, a path with a repeated
// point and a straight run, and a closed one.
func polylineFixtures() []Polyline[float64] {
	return []Polyline[float64]{
		{},
		PolLn([]Point[float64]{Pt(1.0, 2.0)}),
		PolLn([]Point[float64]{Pt(1.0, 2.0), Pt(3.0, 5.0)}),
		PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 3.0)}),
		PolLn([]Point[float64]{Pt(-3.5, 0.25), Pt(-1.0, 2.0), Pt(-1.0, 2.0), Pt(2.0, 2.0), Pt(5.0, 2.0), Pt(7.0, -1.0)}),
		PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0), Pt(0.0, 0.0)}),
	}
}

func ExamplePolLn() {
	fmt.Println(PolLn([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 3)}))
	// Output: PolLn((0,0);(4,0);(4,3))
}

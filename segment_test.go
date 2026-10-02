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

func TestSegment_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, -1), Pt(2, 0)), Segment[int]{Start: Pt(1, -1), End: Pt(2, 0)})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.5, -1.25), Pt(2.5, 3.75)), Segment[float64]{Start: Pt(0.5, -1.25), End: Pt(2.5, 3.75)})
	})
}

func TestSegment_Vector(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertVector(t, Seg(Pt(1, 2), Pt(3, 5)).Vector(), Vec(2, 3))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertVector(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vector(), Vec(0.6, 3.65))
	})
}

func TestSegment_Length(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(1, 2), Pt(3, 5)).Length(), math.Sqrt(13))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Length(), math.Sqrt(13.6825))
	})
	t.Run("a narrow integer span wider than its range", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt[int8](-100, 0), Pt[int8](100, 0)).Length(), 200.0)
	})
}

func TestSegment_Angle(t *testing.T) {
	t.Run("measures from start to end", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(1, 0)).Angle(), 0.0)
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(0, 1)).Angle(), Pi/2)
		geomtest.AssertNumber(t, Seg(Pt(0.0, 0.0), Pt(-1.0, -1.0)).Angle(), -3*Pi/4)
	})
	t.Run("reverse turns it half a turn", func(t *testing.T) {
		s := Seg(Pt(1.0, 2.0), Pt(4.0, 6.0))

		assert.True(t, EqualAngle(s.Reverse().Angle(), s.Angle()+Pi), s.String()+": ")
	})
	t.Run("a zero-length segment has no direction", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(3, 4), Pt(3, 4)).Angle(), 0.0)
	})
	t.Run("a narrow integer span wider than its range", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt[int8](-100, 0), Pt[int8](100, 0)).Angle(), 0.0)
	})
}

func TestSegment_Direction(t *testing.T) {
	t.Run("the nearest direction from start to end", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(0, 0), Pt(5, 0)).Direction(), DirectionRight)
		assert.Equal(t, Seg(Pt(0, 0), Pt(0, -5)).Direction(), DirectionUp)
		assert.Equal(t, Seg(Pt(0, 0), Pt(4, 5)).Direction(), DirectionDownRight)
	})
	t.Run("reverse gives the opposite", func(t *testing.T) {
		s := Seg(Pt(1, 2), Pt(4, 6))

		assert.Equal(t, s.Reverse().Direction(), s.Direction().Opposite())
	})
	t.Run("a zero-length segment has no direction", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(3, 4), Pt(3, 4)).Direction(), DirectionNone)
	})
	t.Run("a narrow integer span wider than its range", func(t *testing.T) {
		assert.Equal(t, Seg(Pt[int8](-100, 0), Pt[int8](100, 0)).Direction(), DirectionRight)
	})
}

func TestSegment_Vertices(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(Seg(Pt(1, 2), Pt(3, 5)).Vertices()), []Point[int]{{1, 2}, {3, 5}})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vertices()), []Point[float64]{{0.6, -0.25}, {1.2, 3.4}})
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for vertex := range Seg(Pt(1, 2), Pt(3, 5)).Vertices() {
			geomtest.AssertPoint(t, vertex, Pt(1, 2))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		s := Seg(Pt(1, 2), Pt(3, 5))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			for vertex := range s.Vertices() {
				sinkBool = vertex.IsZero()
			}
		}), 0)
	})
}

func TestSegment_Edges(t *testing.T) {
	s := Seg(Pt(1, 2), Pt(3, 5))

	t.Run("is the segment itself", func(t *testing.T) {
		edges := slices.Collect(s.Edges())

		assert.Length(t, edges, 1)
		geomtest.AssertSegment(t, edges[0], s)
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			for edge := range s.Edges() {
				sinkBool = edge.IsZero()
			}
		}), 0)
	})
}

func TestSegment_Midpoint(t *testing.T) {
	t.Run("int rounds the half away from zero", func(t *testing.T) {
		geomtest.AssertPoint(t, Seg(Pt(1, 2), Pt(3, 5)).Midpoint(), Pt(2, 4))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPoint(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Midpoint(), Pt(0.9, 1.575))
	})
}

func TestSegment_Bounds(t *testing.T) {
	t.Run("spans the endpoints", func(t *testing.T) {
		s := Seg(Pt(1, 2), Pt(3, 5))

		geomtest.AssertBox(t, s.Bounds(), BoxFromMinMax(s.Start, s.End))
		geomtest.AssertPoint(t, s.Bounds().Center(), Pt(2, 3))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertBox(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Bounds(), BoxFromMinMax(Pt(0.6, -0.25), Pt(1.2, 3.4)))
	})
}

func TestSegment_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).Translate(Vec(3, -2)), Seg(Pt(4, 0), Pt(6, 3)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Translate(Vec(100.1, -0.1)), Seg(Pt(100.7, -0.35), Pt(101.3, 3.3)))
	})
}

func TestSegment_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 6)).MoveTo(Pt(3, -2)), Seg(Pt(2, -4), Pt(4, 0)))
	})
	t.Run("int odd span rounds the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).MoveTo(Pt(3, -2)), Seg(Pt(2, -4), Pt(4, -1)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).MoveTo(Pt(100.1, -0.1)), Seg(Pt(99.8, -1.925), Pt(100.4, 1.725)))
	})
}

func TestSegment_Scale(t *testing.T) {
	t.Run("uniform factor about the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(4, 2)).Scale(2), Seg(Pt(-2, -1), Pt(6, 3)))
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(4.0, 2.0)).Scale(0.5), Seg(Pt(1.0, 0.5), Pt(3.0, 1.5)))
	})
	t.Run("per-axis factor changes the direction", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(4.0, 2.0)).ScaleXY(1, 3), Seg(Pt(0.0, -2.0), Pt(4.0, 4.0)))
	})
	t.Run("zero collapses onto the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(4.0, 2.0)).Scale(0), Seg(Pt(2.0, 1.0), Pt(2.0, 1.0)))
	})
	t.Run("int rounds the midpoint, so an odd span drifts by one", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(3, 0)).Scale(1), Seg(Pt(0, 0), Pt(3, 0)))
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(3, 0)).Scale(2), Seg(Pt(-2, 0), Pt(4, 0)))
	})
	t.Run("keeps the midpoint and scales the length", func(t *testing.T) {
		for _, s := range segmentFixtures {
			scaled := s.Scale(2.5)

			geomtest.AssertPoint(t, scaled.Midpoint(), s.Midpoint(), fmt.Sprintf("%s: ", s))
			geomtest.AssertNumber(t, scaled.Length(), s.Length()*2.5, fmt.Sprintf("%s: ", s))
		}
	})
}

func TestSegment_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).Unscale(2), Seg(Pt(2.0, 1.0), Pt(6.0, 3.0)))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).UnscaleXY(2, 4), Seg(Pt(2.0, 1.5), Pt(6.0, 2.5)))
	})
	t.Run("undoes scale", func(t *testing.T) {
		s := Seg(Pt(1.0, 2.0), Pt(9.0, 6.0))

		geomtest.AssertSegment(t, s.Scale(2.5).Unscale(2.5), s)
		geomtest.AssertSegment(t, s.ScaleXY(2.0, 4.0).UnscaleXY(2.0, 4.0), s)
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestSegment_Resize(t *testing.T) {
	t.Run("keeps the midpoint and the direction", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(6.0, 8.0)) // length 10, midpoint (3,4)
		resized := s.Resize(20)

		geomtest.AssertSegment(t, resized, Seg(Pt(-3.0, -4.0), Pt(9.0, 12.0)))
		geomtest.AssertNumber(t, resized.Length(), 20.0)
		geomtest.AssertPoint(t, resized.Midpoint(), s.Midpoint())
	})
	t.Run("shortens as readily as it lengthens", func(t *testing.T) {
		resized := Seg(Pt(0.0, 0.0), Pt(6.0, 8.0)).Resize(5)

		geomtest.AssertSegment(t, resized, Seg(Pt(1.5, 2.0), Pt(4.5, 6.0)))
		geomtest.AssertNumber(t, resized.Length(), 5.0)
	})
	t.Run("a zero length collapses onto the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(6.0, 8.0)).Resize(0), Seg(Pt(3.0, 4.0), Pt(3.0, 4.0)))
	})
	t.Run("a negative length flips the ends", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(6.0, 8.0))

		geomtest.AssertSegment(t, s.Resize(-10), s.Reverse())
		geomtest.AssertNumber(t, s.Resize(-10).Length(), 10.0)
	})
	t.Run("a zero-length segment resizes along +X", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(3.0, 4.0), Pt(3.0, 4.0)).Resize(10), Seg(Pt(-2.0, 4.0), Pt(8.0, 4.0)))
	})
	t.Run("int rounds both ends", func(t *testing.T) {
		resized := Seg(Pt(0, 0), Pt(6, 8)).Resize(5)

		geomtest.AssertSegment(t, resized, Seg(Pt(2, 2), Pt(5, 6)))
	})
	t.Run("a narrow integer span wider than its range keeps its direction", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt[int8](-100, 0), Pt[int8](100, 0)).Resize(100), Seg(Pt[int8](-50, 0), Pt[int8](50, 0)))
	})
}

func TestSegment_Reverse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).Reverse(), Seg(Pt(3, 5), Pt(1, 2)))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Reverse(), Seg(Pt(1.2, 3.4), Pt(0.6, -0.25)))
	})
}

func TestSegment_PointAt(t *testing.T) {
	segment := Seg(Pt(0.0, 0.0), Pt(4.0, 2.0))

	t.Run("endpoints at 0 and 1", func(t *testing.T) {
		geomtest.AssertPoint(t, segment.PointAt(0), segment.Start)
		geomtest.AssertPoint(t, segment.PointAt(1), segment.End)
	})
	t.Run("along the segment", func(t *testing.T) {
		geomtest.AssertPoint(t, segment.PointAt(0.25), Pt(1.0, 0.5))
		geomtest.AssertPoint(t, segment.PointAt(0.5), segment.Midpoint())
	})
	t.Run("extrapolates beyond the segment", func(t *testing.T) {
		geomtest.AssertPoint(t, segment.PointAt(-0.5), Pt(-2.0, -1.0))
		geomtest.AssertPoint(t, segment.PointAt(1.5), Pt(6.0, 3.0))
	})
	t.Run("int rounds the half away from zero", func(t *testing.T) {
		geomtest.AssertPoint(t, Seg(Pt(0, 0), Pt(3, 3)).PointAt(0.5), Pt(2, 2))
		geomtest.AssertPoint(t, Seg(Pt(0, 0), Pt(-3, -3)).PointAt(0.5), Pt(-2, -2))
	})
	t.Run("lands on the segment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, fraction := range []float64{0, 0.25, 0.5, 0.75, 1} {
				assert.True(t, s.Contains(s.PointAt(fraction)), fmt.Sprintf("%s at %v: ", s, fraction))
			}
		}
	})
}

func TestSegment_Transform(t *testing.T) {
	t.Run("applies the matrix to both points", func(t *testing.T) {
		matrix := Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)

		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 4)).Transform(matrix), Seg(Pt(1, 2).Transform(matrix), Pt(3, 4).Transform(matrix)))
	})
	t.Run("float32 matrix", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1.0, 2.0), Pt(0.0, 0.0)).Transform(Mat[float32](1, 0, 1, 0, 1, 1)), Seg(Pt(2.0, 3.0), Pt(1.0, 1.0)))
	})
}

func TestSegment_Rotate(t *testing.T) {
	t.Run("quarter turn about the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(4, 0)).Rotate(Pi/2), Seg(Pt(2, -2), Pt(2, 2)))
	})
	t.Run("half turn reverses the segment", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 6)).Rotate(Pi), Seg(Pt(3, 6), Pt(1, 2)))
	})
	t.Run("int rounds the midpoint, so an odd span drifts by one", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).Rotate(Pi), Seg(Pt(3, 6), Pt(1, 3)))
	})
	t.Run("float keeps the midpoint and length", func(t *testing.T) {
		s := Seg(Pt(0.6, -0.25), Pt(1.2, 3.4))
		rotated := s.Rotate(0.7)

		geomtest.AssertPoint(t, rotated.Midpoint(), s.Midpoint())
		geomtest.AssertNumber(t, rotated.Length(), s.Length())
		geomtest.AssertNumber(t, rotated.Vector().AngleBetween(s.Vector()), 0.7)
	})
	t.Run("a full turn is identity", func(t *testing.T) {
		for _, s := range segmentFixtures {
			geomtest.AssertSegment(t, s.Rotate(2*Pi), s, fmt.Sprintf("%s: ", s))
		}
	})
}

func TestSegment_Normal(t *testing.T) {
	t.Run("a quarter turn of the vector, with its length", func(t *testing.T) {
		geomtest.AssertVector(t, Seg(Pt(1, 1), Pt(4, 1)).Normal(), Vec(0, 3))
		geomtest.AssertVector(t, Seg(Pt(0.0, 0.0), Pt(0.0, 2.5)).Normal(), Vec(-2.5, 0.0))
	})
	t.Run("points inward on a rectangle edge", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(4, 4))

		for edge := range r.Edges() {
			assert.True(t, r.Contains(edge.Midpoint().Add(edge.Normal().Resize(1))), edge.String())
		}
	})
	t.Run("a zero-length segment has none", func(t *testing.T) {
		geomtest.AssertVector(t, Seg(Pt(3, 4), Pt(3, 4)).Normal(), ZeroVector[int]())
	})
	t.Run("is perpendicular to the segment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			geomtest.AssertNumber(t, s.Normal().Dot(s.Vector()), 0.0, s.String())
			geomtest.AssertNumber(t, s.Normal().Length(), s.Length(), s.String())
		}
	})
}

func TestSegment_Contains(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		s := Seg(Pt(0, 0), Pt(6, 3))

		assert.True(t, s.Contains(Pt(2, 1)))
		assert.True(t, s.Contains(s.Start))
		assert.True(t, s.Contains(s.End))
		assert.False(t, s.Contains(Pt(2, 2)))
		assert.False(t, s.Contains(Pt(8, 4)))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, s.Contains(Pt(0.5, 0.5+Delta/2)))
		assert.False(t, s.Contains(Pt(0.5, 0.5+2*Delta)))
	})
	t.Run("float32 is exact on coordinates far from the origin", func(t *testing.T) {
		s := Seg(Pt[float32](0.0625, 0.1875), Pt[float32](1<<21+1, 3*(1<<21+1)))

		assert.True(t, s.Contains(Pt[float32](1<<20+1, 3*(1<<20+1))))
	})
	t.Run("a narrow integer does not overflow", func(t *testing.T) {
		s := Seg(Pt[int8](-100, -100), Pt[int8](100, 100))

		assert.True(t, s.Contains(Pt[int8](0, 0)))
	})
	t.Run("holds exactly where DistanceTo is zero", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, s.Contains(p), s.DistanceTo(p) == 0, fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_DistanceTo(t *testing.T) {
	s := Seg(Pt(0, 0), Pt(4, 0))

	t.Run("perpendicular to the segment", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(2, 3)), 3.0)
		geomtest.AssertNumber(t, s.DistanceTo(Pt(1, -2)), 2.0)
	})
	t.Run("beyond the start measures to the start", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(-3, 4)), 5.0)
	})
	t.Run("beyond the end measures to the end", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceTo(Pt(7, 4)), 5.0)
	})
	t.Run("on the segment is exactly zero", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(0, 0), Pt(3, 3)).DistanceTo(Pt(1, 1)), 0.0)
		assert.Equal(t, s.DistanceTo(Pt(4, 0)), 0.0)
	})
	t.Run("degenerate segment measures to the point", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(1, 1), Pt(1, 1)).DistanceTo(Pt(4, 5)), 5.0)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.Equal(t, s.DistanceTo(Pt(0.5, Delta/2)), 0.0)
		geomtest.AssertNumber(t, s.DistanceTo(Pt(0.5, 2*Delta)), 2*Delta)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0.0, 0.0), Pt(1.0, 1.0)).DistanceTo(Pt(1.0, 0.0)), OneOverSqrt2)
	})
}

func TestSegment_DistanceSquaredTo(t *testing.T) {
	s := Seg(Pt(0, 0), Pt(4, 0))

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		geomtest.AssertNumber(t, s.DistanceSquaredTo(Pt(2, 3)), 9.0)
		geomtest.AssertNumber(t, s.DistanceSquaredTo(Pt(-3, 4)), 25.0)
		geomtest.AssertNumber(t, s.DistanceSquaredTo(Pt(7, 4)), 25.0)
	})
	t.Run("stays fractional for an integer T", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(2, 1)).DistanceSquaredTo(Pt(0, 1)), 0.8)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, s.DistanceSquaredTo(p), s.DistanceTo(p)*s.DistanceTo(p), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_Nearest(t *testing.T) {
	s := Seg(Pt(0, 0), Pt(4, 0))

	t.Run("the foot of the perpendicular", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(2, 3)), Pt(2, 0))
	})
	t.Run("the endpoint where the foot falls beyond it", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(-3, 4)), Pt(0, 0))
		geomtest.AssertPoint(t, s.Nearest(Pt(7, 4)), Pt(4, 0))
	})
	t.Run("a point on the segment is its own nearest point", func(t *testing.T) {
		geomtest.AssertPoint(t, s.Nearest(Pt(1, 0)), Pt(1, 0))
	})
	t.Run("a point within the tolerance is kept as it is", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(0.0, 0.0), Pt(4.0, 0.0)).Nearest(Pt(2.0, Delta/2)), Pt(2.0, Delta/2))
	})
	t.Run("int rounds once and can land off the segment", func(t *testing.T) {
		skewed := Seg(Pt(0, 0), Pt(4, 2))

		geomtest.AssertPoint(t, skewed.Nearest(Pt(0, 3)), Pt(1, 1))
		assert.False(t, skewed.Contains(Pt(1, 1)))
	})
	t.Run("a zero-length segment is its point", func(t *testing.T) {
		geomtest.AssertPoint(t, Seg(Pt(1, 1), Pt(1, 1)).Nearest(Pt(5, 4)), Pt(1, 1))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range pointFixtures {
				assertNearest[float64](t, s, p)
			}
		}
	})
}

func TestSegment_DistanceToSegment(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing segments are at zero", func(t *testing.T) {
		assert.Equal(t, diagonal.DistanceToSegment(Seg(Pt(0, 4), Pt(4, 0))), 0.0)
	})
	t.Run("touching segments are at zero", func(t *testing.T) {
		assert.Equal(t, diagonal.DistanceToSegment(Seg(Pt(4, 4), Pt(8, 0))), 0.0)
	})
	t.Run("parallel segments measure the gap", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(4, 0)).DistanceToSegment(Seg(Pt(1, 3), Pt(3, 3))), 3.0)
	})
	t.Run("apart measures between the nearest endpoints", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(4, 0)).DistanceToSegment(Seg(Pt(7, 4), Pt(9, 4))), 5.0)
	})
	t.Run("an endpoint nearest an interior point", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(4, 0)).DistanceToSegment(Seg(Pt(2, 2), Pt(2, 5))), 2.0)
	})
	t.Run("symmetric and zero exactly where Intersects holds", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				geomtest.AssertNumber(t, a.DistanceToSegment(b), b.DistanceToSegment(a), fmt.Sprintf("%s → %s: ", a, b))
				assert.Equal(t, a.DistanceToSegment(b) <= Delta, a.IntersectsSegment(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestSegment_DistanceSquaredToSegment(t *testing.T) {
	t.Run("is the square of DistanceToSegment", func(t *testing.T) {
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(4, 0)).DistanceSquaredToSegment(Seg(Pt(7, 4), Pt(9, 4))), 25.0)
		geomtest.AssertNumber(t, Seg(Pt(0, 0), Pt(4, 0)).DistanceSquaredToSegment(Seg(Pt(1, 3), Pt(3, 3))), 9.0)
		assert.Equal(t, Seg(Pt(0, 0), Pt(4, 4)).DistanceSquaredToSegment(Seg(Pt(0, 4), Pt(4, 0))), 0.0)
	})
	t.Run("agrees with DistanceToSegment", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				geomtest.AssertNumber(t, a.DistanceSquaredToSegment(b), a.DistanceToSegment(b)*a.DistanceToSegment(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestSegment_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, s.IntersectsCircle(c), c.IntersectsSegment(s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_IntersectionCircle(t *testing.T) {
	t.Run("matches Circle.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, s.IntersectionCircle(c), c.IntersectionSegment(s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_AppendIntersectionCircle(t *testing.T) {
	t.Run("matches Circle.AppendIntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, s.AppendIntersectionCircle(bufferWith(prefixPoint), c), c.AppendIntersectionSegment(bufferWith(prefixPoint), s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_IntersectsSegment(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsSegment(Seg(Pt(0, 4), Pt(4, 0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsSegment(Seg(Pt(5, 0), Pt(5, 4))))
		assert.False(t, diagonal.IntersectsSegment(Seg(Pt(0, 1), Pt(3, 4))))
	})
	t.Run("touching at an endpoint counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsSegment(Seg(Pt(4, 4), Pt(8, 0))))
		assert.True(t, diagonal.IntersectsSegment(Seg(Pt(2, 2), Pt(2, 8))))
	})
	t.Run("collinear overlap counts and a collinear gap does not", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsSegment(Seg(Pt(2, 2), Pt(6, 6))))
		assert.False(t, diagonal.IntersectsSegment(Seg(Pt(5, 5), Pt(6, 6))))
	})
	t.Run("parallel segments do not cross", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsSegment(Seg(Pt(0, 1), Pt(4, 5))))
	})
	t.Run("a degenerate segment is a point", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsSegment(Seg(Pt(1, 1), Pt(1, 1))))
		assert.False(t, diagonal.IntersectsSegment(Seg(Pt(1, 2), Pt(1, 2))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.True(t, s.IntersectsSegment(Seg(Pt(0.5, Delta/2), Pt(0.5, 1.0))))
		assert.False(t, s.IntersectsSegment(Seg(Pt(0.5, 2*Delta), Pt(0.5, 1.0))))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				assert.Equal(t, a.IntersectsSegment(b), b.IntersectsSegment(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func BenchmarkSegment_IntersectsSegment(b *testing.B) {
	segment := Seg(Pt(0.0, 0.0), Pt(10.0, 10.0))
	crossing, apart := Seg(Pt(0.0, 10.0), Pt(10.0, 0.0)), Seg(Pt(20.0, 0.0), Pt(20.0, 10.0))

	b.Run("crossing", func(b *testing.B) {
		for b.Loop() {
			sinkBool = segment.IntersectsSegment(crossing)
		}
	})
	b.Run("apart", func(b *testing.B) {
		for b.Loop() {
			sinkBool = segment.IntersectsSegment(apart)
		}
	})
}

func TestSegment_IntersectionSegment(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		point, ok := diagonal.IntersectionSegment(Seg(Pt(0, 4), Pt(4, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := diagonal.IntersectionSegment(Seg(Pt(5, 0), Pt(5, 4)))

		assert.False(t, ok)
	})
	t.Run("touching at an endpoint counts", func(t *testing.T) {
		point, ok := diagonal.IntersectionSegment(Seg(Pt(4, 4), Pt(8, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 4))

		point, ok = Seg(Pt(2, 2), Pt(2, 8)).IntersectionSegment(diagonal)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("every endpoint can be the touching one", func(t *testing.T) {
		for _, segment := range []Segment[int]{Seg(Pt(4, 4), Pt(8, 0)), Seg(Pt(8, 0), Pt(4, 4))} {
			point, ok := diagonal.IntersectionSegment(segment)

			assert.True(t, ok, segment.String())
			geomtest.AssertPoint(t, point, Pt(4, 4), segment.String())
		}

		point, ok := diagonal.IntersectionSegment(Seg(Pt(2, 6), Pt(6, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 4))

		point, ok = Seg(Pt(4, 4), Pt(0, 0)).IntersectionSegment(Seg(Pt(2, 6), Pt(6, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(4, 4))
	})
	t.Run("parallel and collinear segments have no single point", func(t *testing.T) {
		_, ok := diagonal.IntersectionSegment(Seg(Pt(0, 1), Pt(4, 5)))
		assert.False(t, ok)

		_, ok = diagonal.IntersectionSegment(Seg(Pt(2, 2), Pt(6, 6)))
		assert.False(t, ok)
	})
	t.Run("the lines cross beyond a segment", func(t *testing.T) {
		_, ok := diagonal.IntersectionSegment(Seg(Pt(5, 0), Pt(6, 4)))

		assert.False(t, ok)
	})
	t.Run("a degenerate segment is the point where it lies on the other", func(t *testing.T) {
		point, ok := diagonal.IntersectionSegment(Seg(Pt(1, 1), Pt(1, 1)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(1, 1))

		point, ok = Seg(Pt(1, 1), Pt(1, 1)).IntersectionSegment(diagonal)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(1, 1))

		_, ok = diagonal.IntersectionSegment(Seg(Pt(1, 2), Pt(1, 2)))
		assert.False(t, ok)

		_, ok = Seg(Pt(1, 1), Pt(1, 1)).IntersectionSegment(Seg(Pt(1, 2), Pt(1, 2)))
		assert.False(t, ok)
	})
	t.Run("nearly collinear float segments fall back to the endpoint on the other", func(t *testing.T) {
		s, other := Seg(Pt(2.37, 35.51), Pt(3.1499999999999995, 57.85)), Seg(Pt(1.2, 2.0), Pt(3.3449999999999998, 63.435))
		point, ok := s.IntersectionSegment(other)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, s.Start)
		assert.True(t, s.IntersectsSegment(other))
	})
	t.Run("a shallow touch is decided on the endpoint distance, like Intersects", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(100.0, 0.0))
		shallow := Seg(Pt(50.0, Delta/2), Pt(150.0, 1e-3))

		point, ok := s.IntersectionSegment(shallow)

		assert.True(t, s.IntersectsSegment(shallow))
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, shallow.Start)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0, 0), Pt(3, 3)).IntersectionSegment(Seg(Pt(0, 3), Pt(3, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("float keeps the crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0.0, 0.0), Pt(3.0, 3.0)).IntersectionSegment(Seg(Pt(0.0, 3.0), Pt(3.0, 0.0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(1.5, 1.5))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		_, ok := s.IntersectionSegment(Seg(Pt(0.5, Delta/2), Pt(0.5, 1.0)))
		assert.True(t, ok)

		_, ok = s.IntersectionSegment(Seg(Pt(0.5, 2*Delta), Pt(0.5, 1.0)))
		assert.False(t, ok)
	})
	t.Run("int decides exactly where the crossing is not a lattice point", func(t *testing.T) {
		a, b := Seg(Pt(-6, -6), Pt(-5, -5)), Seg(Pt(-6, -5), Pt(-4, -6))

		point, ok := a.IntersectionSegment(b)
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(-5, -5))

		point, ok = b.IntersectionSegment(a)
		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(-5, -5))
	})
	t.Run("agrees with Intersects on non-parallel fixtures", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				point, ok := a.IntersectionSegment(b)
				parallel := !a.Vector().IsZero() && !b.Vector().IsZero() && a.Vector().Cross(b.Vector()) == 0

				assert.Equal(t, ok, a.IntersectsSegment(b) && !parallel, fmt.Sprintf("%s → %s: ", a, b))
				if !ok {
					continue
				}

				assert.True(t, a.Contains(point), fmt.Sprintf("%s → %s on a: ", a, b))
				assert.True(t, b.Contains(point), fmt.Sprintf("%s → %s on b: ", a, b))
			}
		}
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				_, ok := a.IntersectionSegment(b)
				_, reverse := b.IntersectionSegment(a)

				assert.Equal(t, ok, reverse, fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("segments running along each other within the tolerance answer with an endpoint", func(t *testing.T) {
		r := Rectangle[float64]{Pt(-9.555794994732457, -7.1447161806712565), Sz(1.9672661593065617, 2.772250161572721), 2.3562338764659856}
		corners := slices.Collect(r.Vertices())
		edge, s := Seg(corners[0], corners[1]), Seg(Pt(-7.1845731257763275, -7.555551584844245), Pt(-9.966817184817401, -4.773526681623244))

		point, ok := s.IntersectionSegment(edge)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, edge.Start)
	})
	t.Run("a crossing on a half unit rounds alike whatever the lengths", func(t *testing.T) {
		edge := Seg(Pt(9, -9), Pt(-1, 1))

		for _, end := range []Point[int]{Pt(1, -12), Pt(-986, -999)} {
			point, ok := Seg(Pt(14, 1), end).IntersectionSegment(edge)

			assert.True(t, ok)
			geomtest.AssertPoint(t, point, Pt(7, -7), fmt.Sprintf("to %s: ", end))
		}
	})
}

func FuzzSegment_IntersectionSegment(f *testing.F) {
	f.Add(0.0, 0.0, 4.0, 4.0, 0.0, 4.0, 4.0, 0.0)
	f.Add(0.0, 0.0, 100.0, 0.0, 50.0, Delta/2, 150.0, 1e-3)
	f.Add(0.0, 0.0, 4.0, 4.0, 1.0, 1.0, 1.0, 1.0)
	f.Add(0.0, 0.0, 4.0, 4.0, 2.0, 2.0, 6.0, 6.0)

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, x3, y3, x4, y4 float64) {
		for _, v := range []float64{x1, y1, x2, y2, x3, y3, x4, y4} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		a, b := Seg(Pt(x1, y1), Pt(x2, y2)), Seg(Pt(x3, y3), Pt(x4, y4))
		parallel := !a.Vector().IsZero() && !b.Vector().IsZero() && a.Vector().Cross(b.Vector()) == 0

		point, ok := a.IntersectionSegment(b)
		_, reverse := b.IntersectionSegment(a)

		assert.Equal(t, ok, a.IntersectsSegment(b) && !parallel, fmt.Sprintf("%s → %s: ", a, b))
		assert.Equal(t, ok, reverse, fmt.Sprintf("%s → %s: symmetric: ", a, b))
		if !ok {
			return
		}

		assert.True(t, a.Contains(point), fmt.Sprintf("%s → %s on a: ", a, b))
		assert.True(t, b.Contains(point), fmt.Sprintf("%s → %s on b: ", a, b))
	})
}

func TestSegment_IntersectsRay(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(0, 4), Vec(1, -1))))
	})
	t.Run("the lines cross behind the origin", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(3, 1), Vec(1, -1))))
	})
	t.Run("the lines cross beyond the segment", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(10, 0), Vec(0, 1))))
	})
	t.Run("the origin on the segment counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(2, 2), Vec(1, 0))))
	})
	t.Run("an endpoint on the ray counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(0, 4), Vec(1, 0))))
	})
	t.Run("collinear toward the segment counts, away from it does not", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(9, 9), Vec(-1, -1))))
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(5, 5), Vec(1, 1))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(RayAlong(Pt(1, 1), Vec(0, 0))))
		assert.False(t, diagonal.IntersectsRay(RayAlong(Pt(1, 2), Vec(0, 0))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.True(t, s.IntersectsRay(RayAlong(Pt(0.5, Delta/2), Vec(0.0, 1.0))))
		assert.False(t, s.IntersectsRay(RayAlong(Pt(0.5, 2*Delta), Vec(0.0, 1.0))))
	})
	t.Run("matches a segment reaching past the segment along the ray", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, s.IntersectsRay(r), s.IntersectsSegment(far(r)), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
}

func TestSegment_IntersectionRay(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(RayAlong(Pt(0, 4), Vec(1, -1)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(RayAlong(Pt(3, 1), Vec(1, -1)))

		assert.False(t, ok)
	})
	t.Run("the origin on the segment is the point", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(RayAlong(Pt(2, 2), Vec(1, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("collinear has no single point", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(RayAlong(Pt(9, 9), Vec(-1, -1)))

		assert.False(t, ok)
	})
	t.Run("a nearly collinear float ray falls back to the endpoint on it", func(t *testing.T) {
		s, r := Seg(Pt(42.400000000000006, 126.41), Pt(70.0, 209.35)), RayAlong(Pt(1.0, 2.0), Vec(13.8, 41.47))
		point, ok := s.IntersectionRay(r)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, s.End)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0, 0), Pt(3, 0)).IntersectionRay(RayAlong(Pt(1, -1), Vec(1, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 0))
	})
	t.Run("agrees with Intersects on non-parallel fixtures, on both", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rayFixtures {
				if s.Vector().Float().Cross(r.Direction.Float()) == 0 && s.Vector() != (Vector[float64]{}) && r.Direction != (Vector[float64]{}) {
					continue
				}

				point, ok := s.IntersectionRay(r)
				message := fmt.Sprintf("%s → %s: ", s, r)

				assert.Equal(t, ok, s.IntersectsRay(r), message)
				if ok {
					assert.True(t, s.Contains(point) && r.Contains(point), message)
				}
			}
		}
	})
}

func TestSegment_IntersectsPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("crossing an edge", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 1), Pt(5, 1)).IntersectsPolygon(square))
		assert.True(t, Seg(Pt(5, 1), Pt(1, 1)).IntersectsPolygon(square))
	})
	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Seg(Pt(-1, 1), Pt(5, 1)).IntersectsPolygon(square))
	})
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 1), Pt(1, 1)).IntersectsPolygon(square))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Seg(Pt(3, -1), Pt(3, 3)).IntersectsPolygon(square))
	})
	t.Run("touching a vertex counts", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectsPolygon(square))
	})
	t.Run("outside the extent is rejected", func(t *testing.T) {
		assert.False(t, Seg(Pt(3, 3), Pt(5, 5)).IntersectsPolygon(square))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Seg(Pt(0, 0), Pt(1, 1)).IntersectsPolygon(Pol[int](nil)))
	})
	t.Run("matches IntersectsPolygon on the segment as a polygon", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, s.IntersectsPolygon(p), p.IntersectsPolygon(Pol(slices.Collect(s.Vertices()))), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_IntersectionPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-1, 1), Pt(5, 1)).IntersectionPolygon(square), []Point[int]{Pt(0, 1), Pt(2, 1)})
		geomtest.AssertVertices(t, Seg(Pt(5, 1), Pt(-1, 1)).IntersectionPolygon(square), []Point[int]{Pt(2, 1), Pt(0, 1)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 1), Pt(5, 1)).IntersectionPolygon(square), []Point[int]{Pt(2, 1)})
	})
	t.Run("through a vertex counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectionPolygon(square), []Point[int]{Pt(2, 2)})
	})
	t.Run("a concave polygon is crossed more than twice", func(t *testing.T) {
		notched := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(2, 1), Pt(0, 4)})

		geomtest.AssertVertices(t, Seg(Pt(-1, 3), Pt(5, 3)).IntersectionPolygon(notched), []Point[int]{Pt(0, 3), Pt(1, 3), Pt(3, 3), Pt(4, 3)})
	})
	t.Run("inside, apart and empty give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(1, 1), Pt(1, 1)).IntersectionPolygon(square))
		assert.Nil(t, Seg(Pt(3, -1), Pt(3, 3)).IntersectionPolygon(square))
		assert.Nil(t, Seg(Pt(-1, 1), Pt(5, 1)).IntersectionPolygon(Pol[int](nil)))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, apart := Seg(Pt(-1, 1), Pt(5, 1)), Seg(Pt(3, -1), Pt(3, 3))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.IntersectionPolygon(square)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = apart.IntersectionPolygon(square)
		}), 0)
	})
	t.Run("matches the rectangle crossings on the rectangle as a polygon", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, s.IntersectionPolygon(r.Polygon()), s.IntersectionRectangle(r), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
	t.Run("every point lies on the segment and an edge, and exists where IntersectsPolygon holds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				assertCrossings(t, s, p, s.IntersectionPolygon(p), s.IntersectsPolygon(p))
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and an edge", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, p := range polygonFixtures() {
				for _, s := range segmentFixtures {
					p, s := p.Cast[float32]().Translate(offset), s.Cast[float32]().Translate(offset)

					assertCrossings(t, s, p, s.IntersectionPolygon(p), s.IntersectsPolygon(p))
				}
			}
		}
	})
	t.Run("int16 spans whose square leaves its range are ordered from Start", func(t *testing.T) {
		s := Seg(Pt[int16](0, 0), Pt[int16](400, 0))
		square := Pol([]Point[int16]{Pt[int16](150, -10), Pt[int16](200, -10), Pt[int16](200, 10), Pt[int16](150, 10)})

		geomtest.AssertVertices(t, s.IntersectionPolygon(square), []Point[int16]{Pt[int16](150, 0), Pt[int16](200, 0)})
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the points follow from Start", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				p, s := p.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				points := s.IntersectionPolygon(p)

				assert.True(t, len(points) == 0 || s.IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
				assertOrdered(t, s, points, fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_AppendIntersectionPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := Seg(Pt(-1, 1), Pt(5, 1))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionPolygon([]Point[int]{Pt(2, 1), Pt(9, 9)}, square), []Point[int]{Pt(2, 1), Pt(9, 9), Pt(0, 1), Pt(2, 1)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(3, -1), Pt(3, 3)).AppendIntersectionPolygon([]Point[int]{Pt(9, 9)}, square), []Point[int]{Pt(9, 9)})
		assert.Nil(t, Seg(Pt(3, -1), Pt(3, 3)).AppendIntersectionPolygon(nil, square))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionPolygon(buffer[:0], square)
		}), 0)
	})
	t.Run("matches IntersectionPolygon after the points in dst", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, s.AppendIntersectionPolygon(bufferWith(prefixPoint), p), append([]Point[float64]{prefixPoint}, s.IntersectionPolygon(p)...), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("crossing an edge", func(t *testing.T) {
		assert.True(t, Seg(Pt(0, 0), Pt(5, 0)).IntersectsRectangle(rectangle))
	})
	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectsRectangle(rectangle))
	})
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectsRectangle(rectangle))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectsRectangle(rectangle))
		assert.False(t, Seg(Pt(3, 3), Pt(5, 5)).IntersectsRectangle(rectangle))
	})
	t.Run("touching a corner counts", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectsRectangle(rectangle))
		assert.False(t, Seg(Pt(2, 4), Pt(4, 2)).IntersectsRectangle(rectangle))
	})
}

func BenchmarkSegment_IntersectsRectangle(b *testing.B) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
	through, apart := Seg(Pt(-20.0, 0.0), Pt(20.0, 0.0)), Seg(Pt(-20.0, 20.0), Pt(20.0, 20.0))

	b.Run("through", func(b *testing.B) {
		for b.Loop() {
			sinkBool = through.IntersectsRectangle(rectangle)
		}
	})
	b.Run("apart", func(b *testing.B) {
		for b.Loop() {
			sinkBool = apart.IntersectsRectangle(rectangle)
		}
	})
}

func TestSegment_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(5, 0), Pt(-5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(0, 0), Pt(5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 0)})
	})
	t.Run("through a corner counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 2)})
		geomtest.AssertVertices(t, Seg(Pt(-4, -4), Pt(4, 4)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, -2), Pt(2, 2)})
	})
	t.Run("an endpoint exactly Delta outside is judged like IntersectsRectangle", func(t *testing.T) {
		s, r := Seg(Pt(-2.0, 2.000001), Pt(2.0, 68.0000005)), Rect(Pt(0.0, 0.0), Sz(26.0, 4.0))

		assert.Equal(t, len(s.IntersectionRectangle(r)) > 0, s.IntersectsRectangle(r))
	})
	t.Run("along an edge crosses the edges at its ends", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, -2), Pt(5, -2)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, -2), Pt(2, -2)})
		geomtest.AssertVertices(t, Seg(Pt(-1, 2), Pt(1, 2)).IntersectionRectangle(rectangle), nil)
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectionRectangle(rectangle))
		assert.Nil(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectionRectangle(rectangle))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, apart := Seg(Pt(-5, 1), Pt(5, 1)), Seg(Pt(3, -5), Pt(3, 5))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.IntersectionRectangle(rectangle)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.IntersectionRectangle(rectangle.Rotate(Pi / 5))
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = apart.IntersectionRectangle(rectangle)
		}), 0)
	})
	t.Run("float keeps the crossings", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5.0, 1.0), Pt(5.0, 1.0)).IntersectionRectangle(Rect(Pt(0.0, 0.0), Sz(3.0, 3.0))), []Point[float64]{Pt(-1.5, 1.0), Pt(1.5, 1.0)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsRectangle holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rectFixtures {
				assertCrossings(t, s, r, s.IntersectionRectangle(r), s.IntersectsRectangle(r))
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, r := range rectFixtures {
					s, r := s.Cast[float32]().Translate(offset), r.Cast[float32]().Translate(offset)

					assertCrossings(t, s, r, s.IntersectionRectangle(r), s.IntersectsRectangle(r))
				}
			}
		}
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the points follow from Start", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				r, s := r.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				points := s.IntersectionRectangle(r)

				assert.True(t, len(points) == 0 || s.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", s, r))
				assertOrdered(t, s, points, fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
	t.Run("a segment a rounding error off an edge crosses only the edges at its ends", func(t *testing.T) {
		r := Rectangle[float64]{Pt(-9.555794994732457, -7.1447161806712565), Sz(1.9672661593065617, 2.772250161572721), 2.3562338764659856}
		corners := slices.Collect(r.Vertices())
		s := Seg(Pt(-7.1845731257763275, -7.555551584844245), Pt(-9.966817184817401, -4.773526681623244))

		geomtest.AssertVertices(t, s.IntersectionRectangle(r), corners[:2])
		geomtest.AssertVertices(t, s.IntersectionPolygon(r.Polygon()), corners[:2])
	})
	t.Run("a segment along an edge of the fixtures, extended past both ends, crosses at most twice", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, turn := range []float64{0, 0.3, 2.3562338764659856} {
				r := r.Rotate(turn)

				for edge := range r.Edges() {
					s := Seg(edge.PointAt(-0.5), edge.PointAt(1.5))

					assert.True(t, len(s.IntersectionRectangle(r)) <= 2, fmt.Sprintf("%s → %s: ", s, r))
				}
			}
		}
	})
}

func TestSegment_AppendIntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))
	through := Seg(Pt(-5, 0), Pt(5, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionRectangle([]Point[int]{Pt(2, 0), Pt(9, 9)}, rectangle), []Point[int]{Pt(2, 0), Pt(9, 9), Pt(-2, 0), Pt(2, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(5, -5), Pt(5, 5)).AppendIntersectionRectangle([]Point[int]{Pt(9, 9)}, rectangle), []Point[int]{Pt(9, 9)})
		assert.Nil(t, Seg(Pt(5, -5), Pt(5, 5)).AppendIntersectionRectangle(nil, rectangle))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionRectangle(buffer[:0], rectangle)
		}), 0)
	})
	t.Run("matches IntersectionRectangle after the points in dst", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rect := range rectFixtures {
				geomtest.AssertVertices(t, s.AppendIntersectionRectangle(bufferWith(prefixPoint), rect), append([]Point[float64]{prefixPoint}, s.IntersectionRectangle(rect)...), fmt.Sprintf("%s → %s: ", s, rect))
			}
		}
	})
}

func FuzzSegment_IntersectionRectangle(f *testing.F) {
	f.Add(-5.0, 0.0, 5.0, 0.0, 0.0, 0.0, 4.0, 4.0)
	f.Add(-5.0, -2.0, 5.0, -2.0, 0.0, 0.0, 4.0, 4.0)
	f.Add(-1.0, -1.0, 1.0, 1.0, 0.0, 0.0, 4.0, 4.0)
	f.Add(-2.0, 1.0+Delta/2, 2.0, 1.0+Delta/2, 0.0, 0.0, 2.0, 2.0)
	f.Add(-2.0, 2.000001, 2.0, 68.0000005, 0.0, 0.0, -26.0, 4.0)

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, cx, cy, w, h float64) {
		for _, v := range []float64{x1, y1, x2, y2, cx, cy, w, h} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		s, r := Seg(Pt(x1, y1), Pt(x2, y2)), Rect(Pt(cx, cy), Sz(w, h))
		points := s.IntersectionRectangle(r)

		assert.True(t, len(points) <= 4, fmt.Sprintf("%s → %s: at most four crossings, got %d: ", s, r, len(points)))

		for _, p := range points {
			assert.True(t, s.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", s, r, p))
			assert.True(t, slices.ContainsFunc(slices.Collect(r.Edges()), func(edge Segment[float64]) bool {
				return edge.Contains(p)
			}), fmt.Sprintf("%s → %s: %s on the boundary: ", s, r, p))
		}

		if inside := r.Contains(s.Start) && r.Contains(s.End); !inside {
			assert.Equal(t, len(points) > 0, s.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", s, r))
		}
	})
}

func TestSegment_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectsRegularPolygon(diamond))
	})
	t.Run("starting inside", func(t *testing.T) {
		assert.True(t, Seg(Pt(0, 0), Pt(5, 5)).IntersectsRegularPolygon(diamond))
	})
	t.Run("touching a vertex", func(t *testing.T) {
		assert.True(t, Seg(Pt(2, -2), Pt(2, 2)).IntersectsRegularPolygon(diamond))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Seg(Pt(2, 2), Pt(3, 1)).IntersectsRegularPolygon(diamond))
		assert.False(t, Seg(Pt(5, 0), Pt(6, 0)).IntersectsRegularPolygon(diamond))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, s.IntersectsRegularPolygon(rp), s.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_IntersectionRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectionRegularPolygon(diamond), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(3, 0), Pt(-3, 0)).IntersectionRegularPolygon(diamond), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("a vertex hit by two edges is counted once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(2, -2), Pt(2, 2)).IntersectionRegularPolygon(diamond), []Point[int]{Pt(2, 0)})
	})
	t.Run("inside, apart and empty give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(0, 0), Pt(1, 0)).IntersectionRegularPolygon(diamond))
		assert.Nil(t, Seg(Pt(5, 0), Pt(6, 0)).IntersectionRegularPolygon(diamond))
		assert.Nil(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectionRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, apart := Seg(Pt(-3, 0), Pt(3, 0)), Seg(Pt(5, 0), Pt(6, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.IntersectionRegularPolygon(diamond)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = apart.IntersectionRegularPolygon(diamond)
		}), 0)
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, s.IntersectionRegularPolygon(rp), s.IntersectionPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the points follow from Start", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				rp, s := rp.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				points := s.IntersectionRegularPolygon(rp)

				assert.True(t, len(points) == 0 || s.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", s, rp))
				assertOrdered(t, s, points, fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_AppendIntersectionRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)
	through := Seg(Pt(-3, 0), Pt(3, 0))

	t.Run("appends after the points in dst, comparing and ordering only its own", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionRegularPolygon([]Point[int]{Pt(2, 0), Pt(9, 9)}, diamond), []Point[int]{Pt(2, 0), Pt(9, 9), Pt(-2, 0), Pt(2, 0)})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(5, 0), Pt(6, 0)).AppendIntersectionRegularPolygon([]Point[int]{Pt(9, 9)}, diamond), []Point[int]{Pt(9, 9)})
		assert.Nil(t, Seg(Pt(5, 0), Pt(6, 0)).AppendIntersectionRegularPolygon(nil, diamond))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Point[int], 0, 2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = through.AppendIntersectionRegularPolygon(buffer[:0], diamond)
		}), 0)
	})
	t.Run("matches IntersectionRegularPolygon after the points in dst", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, s.AppendIntersectionRegularPolygon(bufferWith(prefixPoint), rp), append([]Point[float64]{prefixPoint}, s.IntersectionRegularPolygon(rp)...), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_IntersectsBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectsBox(box))
	})
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectsBox(box))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectsBox(box))
		assert.False(t, Seg(Pt(3, 3), Pt(5, 5)).IntersectsBox(box))
	})
	t.Run("touching a corner counts", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectsBox(box))
		assert.False(t, Seg(Pt(2, 4), Pt(4, 2)).IntersectsBox(box))
	})
	t.Run("a segment with an endpoint in the box intersects it", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, b := range boxFixtures {
				if b.Contains(s.Start) || b.Contains(s.End) {
					assert.True(t, s.IntersectsBox(b), fmt.Sprintf("%s → %s: ", s, b))
				}
			}
		}
	})
	t.Run("decides on the corners of the box, where Contains does", func(t *testing.T) {
		boxes := []Box[float32]{Bx(Pt[float32](3.5497787, 1.5453138), Pt[float32](7.1348896, 4.0049815)), Bx(Pt[float32](8.56645, 9.768216), Pt[float32](9.503749, 19.746952))}
		points := []Point[float32]{Pt[float32](3.5496886, 1.5452706), Pt[float32](8.566437, 9.768117)}

		for i, b := range boxes {
			p := points[i]

			assert.Equal(t, Seg(p, p).IntersectsBox(b), b.Contains(p), fmt.Sprintf("%s → %s: ", p, b))
		}
	})
}

func TestSegment_IntersectionBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectionBox(box), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(5, 0), Pt(-5, 0)).IntersectionBox(box), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("through a corner counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectionBox(box), []Point[int]{Pt(2, 2)})
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectionBox(box))
		assert.Nil(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectionBox(box))
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsBox holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, b := range boxFixtures {
				assertCrossings(t, s, b.Rectangle(), s.IntersectionBox(b), s.IntersectsBox(b))
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, b := range boxFixtures {
					s, b := s.Cast[float32]().Translate(offset), b.Cast[float32]().Translate(offset)

					assertCrossings(t, s, b.Rectangle(), s.IntersectionBox(b), s.IntersectsBox(b))
				}
			}
		}
	})
}

func TestSegment_AppendIntersectionBox(t *testing.T) {
	t.Run("matches AppendIntersectionRectangle on the box Rectangle", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, b := range boxFixtures {
				geomtest.AssertVertices(t, s.AppendIntersectionBox(bufferWith(prefixPoint), b), s.AppendIntersectionRectangle(bufferWith(prefixPoint), b.Rectangle()), fmt.Sprintf("%s → %s: ", s, b))
			}
		}
	})
}

func TestSegment_ClipCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through gives the chord from Start to End", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(-1.0, 0.0), Pt(1.0, 0.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(2.0, 0.0), Pt(-2.0, 0.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(1.0, 0.0), Pt(-1.0, 0.0))})
	})
	t.Run("an endpoint inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(0.5, 0.0), Pt(5.0, 0.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(0.5, 0.0), Pt(1.0, 0.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-2.0, 0.0), Pt(0.0, 0.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(-1.0, 0.0), Pt(0.0, 0.0))})
	})
	t.Run("inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-0.5, 0.0), Pt(0.5, 0.5)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(-0.5, 0.0), Pt(0.5, 0.5))})
	})
	t.Run("a tangent is the point of contact", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(0.0, 1.0), Pt(0.0, 1.0))})
	})
	t.Run("an endpoint touching from outside is that point", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(1.0, 0.0), Pt(2.0, 0.0)).ClipCircle(circle)), []Segment[float64]{Seg(Pt(1.0, 0.0), Pt(1.0, 0.0))})
	})
	t.Run("apart gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0)).ClipCircle(circle)))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-10, 1), Pt(10, 1)).ClipCircle(Circ(Pt(0, 0), 5))), []Segment[int]{Seg(Pt(-5, 1), Pt(5, 1))})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		unit, through := Circ(Pt(0, 0), 5), Seg(Pt(-10, 0), Pt(10, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipCircle(unit)
		}), 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				assertClipped(t, s, partsOf(s.ClipCircle(c)), s.IntersectsCircle(c), c.EnclosesSegment, fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_ClipPolygon(t *testing.T) {
	square := Pol(squareVertices())
	notched := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(2, 1), Pt(0, 4)})

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(-1, 1), Pt(5, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(0, 1), Pt(2, 1))})
		geomtest.AssertSegments(t, Seg(Pt(5, 1), Pt(-1, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(2, 1), Pt(0, 1))})
	})
	t.Run("an endpoint inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(1, 1), Pt(5, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(1, 1), Pt(2, 1))})
	})
	t.Run("inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(1, 1), Pt(1, 2)).ClipPolygon(square), []Segment[int]{Seg(Pt(1, 1), Pt(1, 2))})
		geomtest.AssertSegments(t, Seg(Pt(1, 1), Pt(1, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(1, 1), Pt(1, 1))})
	})
	t.Run("a concave polygon cuts the segment into parts from Start to End", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(-1, 3), Pt(5, 3)).ClipPolygon(notched), []Segment[int]{Seg(Pt(0, 3), Pt(1, 3)), Seg(Pt(3, 3), Pt(4, 3))})
		geomtest.AssertSegments(t, Seg(Pt(5, 3), Pt(-1, 3)).ClipPolygon(notched), []Segment[int]{Seg(Pt(4, 3), Pt(3, 3)), Seg(Pt(1, 3), Pt(0, 3))})
	})
	t.Run("a gap of one unit between two crossings is judged where it is", func(t *testing.T) {
		u := Pol([]Point[int]{Pt(0, 0), Pt(5, 0), Pt(5, 4), Pt(3, 4), Pt(3, 1), Pt(2, 1), Pt(2, 4), Pt(0, 4)})

		geomtest.AssertSegments(t, Seg(Pt(-1, 2), Pt(6, 2)).ClipPolygon(u), []Segment[int]{Seg(Pt(0, 2), Pt(2, 2)), Seg(Pt(3, 2), Pt(5, 2))})
	})
	t.Run("reflex vertices with the outside between them bound two parts", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(-1, 4), Pt(9, 4)).ClipPolygon(Pol(notchedVertices())), []Segment[int]{Seg(Pt(0, 4), Pt(2, 4)), Seg(Pt(6, 4), Pt(8, 4))})
	})
	t.Run("a vertex touched from inside joins the parts either side", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(-1, 1), Pt(5, 1)).ClipPolygon(notched), []Segment[int]{Seg(Pt(0, 1), Pt(4, 1))})
	})
	t.Run("a vertex touched from outside is a part of zero length", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(1, 3), Pt(3, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(2, 2), Pt(2, 2))})
	})
	t.Run("along an edge gives the part of the edge it covers", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(-1, 0), Pt(3, 0)).ClipPolygon(square), []Segment[int]{Seg(Pt(0, 0), Pt(2, 0))})
	})
	t.Run("apart and empty give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(3, -1), Pt(3, 3)).ClipPolygon(square))
		assert.Nil(t, Seg(Pt(-1, 1), Pt(5, 1)).ClipPolygon(Pol[int](nil)))
	})
	t.Run("allocates the result alone", func(t *testing.T) {
		through, apart := Seg(Pt(-1, 1), Pt(5, 1)), Seg(Pt(3, -1), Pt(3, 3))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkSegments = through.ClipPolygon(square)
		}), 1)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkSegments = apart.ClipPolygon(square)
		}), 0)
	})
	t.Run("matches ClipRectangle on the rectangle as a polygon", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertSegments(t, s.ClipPolygon(r.Polygon()), partsOf(s.ClipRectangle(r)), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, p := range outlineFixtures() {
			for _, s := range segmentFixtures {
				assertClipped(t, s, s.ClipPolygon(p), s.IntersectsPolygon(p), p.EnclosesSegment, fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("int16 spans whose square leaves its range are clipped from Start", func(t *testing.T) {
		s := Seg(Pt[int16](0, 0), Pt[int16](400, 0))
		square := Pol([]Point[int16]{Pt[int16](150, -10), Pt[int16](200, -10), Pt[int16](200, 10), Pt[int16](150, 10)})

		geomtest.AssertSegments(t, s.ClipPolygon(square), []Segment[int16]{Seg(Pt[int16](150, 0), Pt[int16](200, 0))})
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the parts follow from Start", func(t *testing.T) {
		for _, p := range outlineFixtures() {
			for _, s := range segmentFixtures {
				p, s := p.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				parts := s.ClipPolygon(p)

				assert.Equal(t, len(parts) > 0, s.IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
				assertOrdered(t, s, endsOf(parts), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("a NaN or infinite End ends the sweep", func(t *testing.T) {
		square := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(10.0, 0.0), Pt(10.0, 10.0), Pt(0.0, 10.0)})

		for _, end := range []Point[float64]{Pt(math.NaN(), 5.0), Pt(math.Inf(1), 5.0)} {
			geomtest.AssertSegments(t, Seg(Pt(1.0, 1.0), end).ClipPolygon(square), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(1.0, 1.0))}, fmt.Sprintf("to %s: ", end))
		}
	})
}

func TestSegment_AppendClipPolygon(t *testing.T) {
	square := Pol(squareVertices())
	through := Seg(Pt(-1, 1), Pt(5, 1))
	prefix := Seg(Pt(-7.5, 3.25), Pt(1.0, 1.0))

	t.Run("appends the parts after the segments in dst", func(t *testing.T) {
		geomtest.AssertSegments(t, through.AppendClipPolygon([]Segment[int]{Seg(Pt(9, 9), Pt(9, 9))}, square), []Segment[int]{Seg(Pt(9, 9), Pt(9, 9)), Seg(Pt(0, 1), Pt(2, 1))})
	})
	t.Run("none leaves dst as it is", func(t *testing.T) {
		apart := Seg(Pt(3, -1), Pt(3, 3))

		geomtest.AssertSegments(t, apart.AppendClipPolygon([]Segment[int]{Seg(Pt(9, 9), Pt(9, 9))}, square), []Segment[int]{Seg(Pt(9, 9), Pt(9, 9))})
		assert.Nil(t, apart.AppendClipPolygon(nil, square))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		buffer := make([]Segment[int], 0, 1)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkSegments = through.AppendClipPolygon(buffer[:0], square)
		}), 0)
	})
	t.Run("matches ClipPolygon after the segments in dst", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range outlineFixtures() {
				geomtest.AssertSegments(t, s.AppendClipPolygon(bufferWith(prefix), p), append([]Segment[float64]{prefix}, s.ClipPolygon(p)...), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_ClipRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-5, 0), Pt(5, 0)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(-2, 0), Pt(2, 0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(5, 0), Pt(-5, 0)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(2, 0), Pt(-2, 0))})
	})
	t.Run("an endpoint inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(0, 0), Pt(5, 0)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(0, 0), Pt(2, 0))})
	})
	t.Run("inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-1, -1), Pt(1, 1)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(-1, -1), Pt(1, 1))})
	})
	t.Run("a corner touched from outside is a part of zero length", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(1, 3), Pt(3, 1)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(2, 2), Pt(2, 2))})
	})
	t.Run("along an edge gives the part of the edge it covers", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-5, -2), Pt(5, -2)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(-2, -2), Pt(2, -2))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(0, -2), Pt(5, -2)).ClipRectangle(rectangle)), []Segment[int]{Seg(Pt(0, -2), Pt(2, -2))})
	})
	t.Run("apart gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Seg(Pt(3, -5), Pt(3, 5)).ClipRectangle(rectangle)))
		assert.Nil(t, partsOf(Seg(Pt(2, 4), Pt(4, 2)).ClipRectangle(rectangle)))
	})
	t.Run("a rotated rectangle clips on its turned edges", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0)).Rotate(Pi / 4)

		geomtest.AssertSegments(t, partsOf(Seg(Pt(-3.0, 0.0), Pt(3.0, 0.0)).ClipRectangle(diamond)), []Segment[float64]{Seg(Pt(-Sqrt2, 0.0), Pt(Sqrt2, 0.0))})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-5, 1), Pt(5, 1))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipRectangle(rectangle)
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipRectangle(rectangle.Rotate(Pi / 5))
		}), 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				assertClipped(t, s, partsOf(s.ClipRectangle(r)), s.IntersectsRectangle(r), r.EnclosesSegment, fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
	t.Run("int16 spans whose square leaves its range are clipped from Start", func(t *testing.T) {
		s := Seg(Pt[int16](0, 0), Pt[int16](400, 0))
		rectangle := RectangleFromMinMax(Pt[int16](150, -10), Pt[int16](200, 10))

		geomtest.AssertSegments(t, partsOf(s.ClipRectangle(rectangle)), []Segment[int16]{Seg(Pt[int16](150, 0), Pt[int16](200, 0))})
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the parts follow from Start", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				r, s := r.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				parts := partsOf(s.ClipRectangle(r))

				assert.Equal(t, len(parts) > 0, s.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", s, r))
				assertOrdered(t, s, endsOf(parts), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
}

func TestSegment_ClipRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-3, 0), Pt(3, 0)).ClipRegularPolygon(diamond)), []Segment[int]{Seg(Pt(-2, 0), Pt(2, 0))})
	})
	t.Run("an endpoint inside is kept", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(0, 0), Pt(5, 0)).ClipRegularPolygon(diamond)), []Segment[int]{Seg(Pt(0, 0), Pt(2, 0))})
	})
	t.Run("inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(0, 0), Pt(1, 0)).ClipRegularPolygon(diamond)), []Segment[int]{Seg(Pt(0, 0), Pt(1, 0))})
	})
	t.Run("a vertex touched from outside is a part of zero length", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(2, -2), Pt(2, 2)).ClipRegularPolygon(diamond)), []Segment[int]{Seg(Pt(2, 0), Pt(2, 0))})
	})
	t.Run("apart and empty give none", func(t *testing.T) {
		assert.Nil(t, partsOf(Seg(Pt(5, 0), Pt(6, 0)).ClipRegularPolygon(diamond)))
		assert.Nil(t, partsOf(Seg(Pt(-3, 0), Pt(3, 0)).ClipRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0))))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-3, 0), Pt(3, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipRegularPolygon(diamond)
		}), 0)
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertSegments(t, partsOf(s.ClipRegularPolygon(rp)), s.ClipPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				assertClipped(t, s, partsOf(s.ClipRegularPolygon(rp)), s.IntersectsRegularPolygon(rp), rp.EnclosesSegment, fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the parts follow from Start", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				rp, s := rp.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				parts := partsOf(s.ClipRegularPolygon(rp))

				assert.Equal(t, len(parts) > 0, s.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", s, rp))
				assertOrdered(t, s, endsOf(parts), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_ClipBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-5, 0), Pt(5, 0)).ClipBox(box)), []Segment[int]{Seg(Pt(-2, 0), Pt(2, 0))})
	})
	t.Run("apart gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Seg(Pt(3, -5), Pt(3, 5)).ClipBox(box)))
	})
	t.Run("matches ClipRectangle on the box's Rectangle", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertSegments(t, partsOf(s.ClipBox(b)), partsOf(s.ClipRectangle(b.Rectangle())), fmt.Sprintf("%s → %s: ", s, b))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-5, 0), Pt(5, 0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = through.ClipBox(box)
		}), 0)
	})
}

func TestSegment_Equal(t *testing.T) {
	t.Run("same segment", func(t *testing.T) {
		assert.True(t, Seg(Pt(1, 2), Pt(3, 5)).Equal(Seg(Pt(1, 2), Pt(3, 5))))
		assert.True(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Seg(Pt(0.6, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("different segment", func(t *testing.T) {
		assert.False(t, Seg(Pt(1, 2), Pt(3, 5)).Equal(Seg(Pt(1, 2), Pt(3, 4))))
		assert.False(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Seg(Pt(0.5, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("orientation matters", func(t *testing.T) {
		assert.False(t, Seg(Pt(1, 2), Pt(3, 5)).Equal(Seg(Pt(3, 5), Pt(1, 2))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Seg(Pt(0.6, -0.250001), Pt(1.2, 3.4))))
	})
}

func TestSegment_IsZero(t *testing.T) {
	t.Run("zero segment", func(t *testing.T) {
		assert.True(t, Segment[int]{}.IsZero())
		assert.True(t, Seg(Pt(0, 0), Pt(0, 0)).IsZero())
		assert.True(t, Segment[float64]{}.IsZero())
	})
	t.Run("one endpoint off the origin", func(t *testing.T) {
		assert.False(t, Seg(Pt(1, 0), Pt(0, 0)).IsZero())
		assert.False(t, Seg(Pt(0, 0), Pt(0, 1)).IsZero())
	})
	t.Run("non-zero segment", func(t *testing.T) {
		assert.False(t, Seg(Pt(1, 2), Pt(3, 5)).IsZero())
		assert.False(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Seg(Pt(0.0, 0.000001), Pt(0.0, 0.0)).IsZero())
	})
}

func TestSegment_Cast(t *testing.T) {
	s := Seg(Pt(1.5, -2.5), Pt(3.5, 4.5))

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertSegment(t, s.Cast[int](), s.Int())
		geomtest.AssertSegment(t, s.Cast[float64](), s.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertSegment(t, s.Cast[int8](), Seg(Pt[int8](2, -3), Pt[int8](4, 5)))
	})
}

func TestSegment_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).Int(), Seg(Pt(1, 2), Pt(3, 5)))
	})
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Int(), Seg(Pt(1, 0), Pt(1, 3)))
	})
}

func TestSegment_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 5)).Float(), Seg(Pt(1.0, 2.0), Pt(3.0, 5.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)).Float(), Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)))
	})
}

func TestSegment_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(10, 16), Pt(1, 2)).String(), "Seg((10,16);(1,2))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(100, -34.0000115), Pt(0.2, 0.4)).String(), "Seg((100.00,-34.00);(0.20,0.40))")
	})
}

func TestSegment_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Seg(Pt(10, 16), Pt(1, 2)), `{"s":{"x":10,"y":16},"e":{"x":1,"y":2}}`)

		var s Segment[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":10,"y":16},"e":{"x":1,"y":2}}`), &s))
		geomtest.AssertSegment(t, s, Seg(Pt(10, 16), Pt(1, 2)))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Seg(Pt(100, -34.0000115), Pt(0.2, 0.4)), `{"s":{"x":100.0,"y":-34.0000115},"e":{"x":0.2,"y":0.4}}`)

		var s Segment[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":10.1,"y":-34.0000115},"e":{"x":0.2,"y":0.4}}`), &s))
		geomtest.AssertSegment(t, s, Seg(Pt(10.1, -34.0000115), Pt(0.2, 0.4)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			data, err := json.Marshal(segment)
			assert.NoError(t, err)

			var decoded Segment[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, segment)
		}
	})
}

func TestSegment_Properties(t *testing.T) {
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, s.Scale(factor).Unscale(factor).Equal(s), fmt.Sprintf("%s x%v: ", s, factor))
			}
		}
	})
	t.Run("resize keeps the midpoint and reaches the length", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, length := range []float64{0, 1, 12.5} {
				resized := s.Resize(length)

				geomtest.AssertPoint(t, resized.Midpoint(), s.Midpoint(), fmt.Sprintf("%s ->%v: ", s, length))
				geomtest.AssertNumber(t, resized.Length(), length, fmt.Sprintf("%s ->%v: ", s, length))
			}
		}
	})
	t.Run("angle and direction follow the vector", func(t *testing.T) {
		for _, s := range segmentFixtures {
			geomtest.AssertNumber(t, s.Angle(), s.Vector().Angle(), s.String()+": ")
			assert.Equal(t, s.Direction(), s.Vector().Direction(), s.String()+": ")
		}
	})
	t.Run("reverse is its own inverse", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			assert.True(t, segment.Reverse().Reverse().Equal(segment), fmt.Sprintf("%s: ", segment))
		}
	})
	t.Run("reverse keeps the length and midpoint", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			geomtest.AssertNumber(t, segment.Reverse().Length(), segment.Length(), fmt.Sprintf("%s: ", segment))
			assert.True(t, segment.Reverse().Midpoint().Equal(segment.Midpoint()), fmt.Sprintf("%s: ", segment))
		}
	})
	t.Run("length is the vector length", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			geomtest.AssertNumber(t, segment.Length(), segment.Vector().Length(), fmt.Sprintf("%s: ", segment))
			geomtest.AssertNumber(t, segment.Length(), segment.Start.DistanceTo(segment.End), fmt.Sprintf("%s: ", segment))
		}
	})
	t.Run("midpoint is equidistant from both ends", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			midpoint := segment.Midpoint()

			geomtest.AssertNumber(t, midpoint.DistanceTo(segment.Start), midpoint.DistanceTo(segment.End), fmt.Sprintf("%s: ", segment))
		}
	})
	t.Run("translate keeps the vector", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			for _, vector := range vectorFixtures {
				assert.True(t, segment.Translate(vector).Vector().Equal(segment.Vector()), fmt.Sprintf("%s → %s: ", segment, vector))
			}
		}
	})
	t.Run("move to keeps the vector and centers where asked", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			for _, point := range pointFixtures {
				moved := segment.MoveTo(point)

				assert.True(t, moved.Midpoint().Equal(point), fmt.Sprintf("%s → %s: ", segment, point))
				assert.True(t, moved.Vector().Equal(segment.Vector()), fmt.Sprintf("%s → %s: ", segment, point))
			}
		}
	})
	t.Run("bounds span the endpoints", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			bounds := segment.Bounds()
			start, end := segment.Start, segment.End

			assert.True(t, bounds.Min.Equal(Pt(min(start.X, end.X), min(start.Y, end.Y))), fmt.Sprintf("%s: ", segment))
			assert.True(t, bounds.Max.Equal(Pt(max(start.X, end.X), max(start.Y, end.Y))), fmt.Sprintf("%s: ", segment))
		}
	})
	t.Run("vertices are the endpoints", func(t *testing.T) {
		for _, segment := range segmentFixtures {
			geomtest.AssertVertices(t, slices.Collect(segment.Vertices()), []Point[float64]{segment.Start, segment.End})
		}
	})
}

func TestSegment_Immutable(t *testing.T) {
	s := Seg(Pt(1, 2), Pt(3, 5))

	s.Translate(Vec(3, -2))
	s.MoveTo(Pt(4, 3))
	s.Reverse()

	geomtest.AssertSegment(t, s, Seg(Pt(1, 2), Pt(3, 5)))
}

// partsOf gives the part a convex Clip method returns as the parts ClipPolygon would, so both
// are checked alike: the part alone where it exists, and nil where it does not.
func partsOf[T Number](part Segment[T], ok bool) []Segment[T] {
	if !ok {
		return nil
	}

	return []Segment[T]{part}
}

// assertOrdered checks that the points follow one another from the segment's Start, the order
// every crossing and every clipped part is returned in, measured in float64. An integer T
// rounds them onto the lattice and off the segment, so their order is all that is asserted.
func assertOrdered[T Number](t *testing.T, s Segment[T], points []Point[T], message string) {
	t.Helper()

	reached := 0.0
	for _, p := range points {
		distance := s.Start.DistanceTo(p)

		assert.True(t, reached <= distance, message+p.String()+" in order: ")

		reached = distance
	}
}

// endsOf returns the ends of the parts in order, the points a clip visits from Start.
func endsOf[T Number](parts []Segment[T]) []Point[T] {
	var ends []Point[T]
	for _, part := range parts {
		ends = append(ends, part.Start, part.End)
	}

	return ends
}

// assertClipped checks the parts a Clip method returns for a segment over the fixtures: they
// exist exactly where the segment intersects the shape, each lies on the segment and within the
// shape, and they follow one another from Start without overlapping.
func assertClipped[T Number](t *testing.T, s Segment[T], parts []Segment[T], intersects bool, encloses func(Segment[T]) bool, message string) {
	t.Helper()

	assert.Equal(t, len(parts) > 0, intersects, message)

	reached := 0.0
	for _, part := range parts {
		from, to := s.Start.DistanceTo(part.Start), s.Start.DistanceTo(part.End)

		assert.True(t, s.Contains(part.Start) && s.Contains(part.End), message+part.String()+" on the segment: ")
		assert.True(t, encloses(part), message+part.String()+" within the shape: ")
		assert.True(t, LessOrEqual(reached, from) && LessOrEqual(from, to), message+part.String()+" in order: ")

		reached = to
	}
}

// segmentFixtures span axis-aligned, diagonal, degenerate, and backwards segments.
var segmentFixtures = []Segment[float64]{
	Seg(Pt(0.0, 0.0), Pt(0.0, 0.0)),
	Seg(Pt(1.0, 2.0), Pt(3.0, 5.0)),
	Seg(Pt(0.6, -0.25), Pt(1.2, 3.4)),
	Seg(Pt(-3.5, 0.25), Pt(-3.5, 7.0)),
	Seg(Pt(2.0, 2.0), Pt(-4.0, 2.0)),
	Seg(Pt(12.5, -0.1), Pt(-0.5, 12.75)),
}

func ExampleSeg() {
	fmt.Println(Seg(Pt(1, 2), Pt(3, 5)))
	// Output: Seg((1,2);(3,5))
}

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

func TestSegment_Normal(t *testing.T) {
	t.Run("a quarter turn of the vector, with its length", func(t *testing.T) {
		geomtest.AssertVector(t, Seg(Pt(1, 1), Pt(4, 1)).Normal(), Vec(0, 3))
		geomtest.AssertVector(t, Seg(Pt(0.0, 0.0), Pt(0.0, 2.5)).Normal(), Vec(-2.5, 0.0))
	})
	t.Run("points inward on a rectangle edge", func(t *testing.T) {
		r := Rect(Pt(0, 0), Sz(4, 4), 0)

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

		assertAllocations(t, func() {
			for vertex := range s.Vertices() {
				sinkBool = vertex.IsZero()
			}
		}, 0)
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
		assertAllocations(t, func() {
			for edge := range s.Edges() {
				sinkBool = edge.IsZero()
			}
		}, 0)
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

func TestSegment_BoundingCircle(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertCircle(t, Seg(Pt(0, 0), Pt(6, 8)).BoundingCircle(), Circ(Pt(3, 4), 5))
	})
	t.Run("int reaches the farther endpoint from the rounded midpoint", func(t *testing.T) {
		geomtest.AssertCircle(t, Seg(Pt(0, 0), Pt(1, 0)).BoundingCircle(), Circ(Pt(1, 0), 1))
		geomtest.AssertCircle(t, Seg(Pt(0, 0), Pt(3, 3)).BoundingCircle(), Circ(Pt(2, 2), 3))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertCircle(t, Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)).BoundingCircle(), Circ(Pt(1.5, 2.0), 2.5))
	})
	t.Run("a zero-length segment is a circle of no radius", func(t *testing.T) {
		geomtest.AssertCircle(t, Seg(Pt(1.5, 2.0), Pt(1.5, 2.0)).BoundingCircle(), Circ(Pt(1.5, 2.0), 0.0))
	})
	t.Run("contains both endpoints, rounded to int too", func(t *testing.T) {
		for _, s := range segmentFixtures {
			assertBoundingCircle(t, s.BoundingCircle(), s.Vertices(), s.String())
			assertBoundingCircle(t, s.Int().BoundingCircle(), s.Int().Vertices(), s.String())
		}
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

func TestSegment_Scale(t *testing.T) {
	t.Run("uniform factor about the midpoint", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(4, 2)).Scale(2), Seg(Pt(-2, -1), Pt(6, 3)))
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(4.0, 2.0)).Scale(0.5), Seg(Pt(1.0, 0.5), Pt(3.0, 1.5)))
	})
	t.Run("per-axis factor changes the direction", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0.0, 0.0), Pt(4.0, 2.0)).ScaleXY(1, 3), Seg(Pt(0.0, -2.0), Pt(4.0, 4.0)))
	})
	t.Run("the offsets from the midpoint are taken in float64", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt[int8](-128, 0), Pt[int8](127, 0)).Scale(0.5), Seg(Pt[int8](-65, 0), Pt[int8](63, 0)))
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
		assert.PanicsWith(t, func() {
			Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).Unscale(0)
		}, ErrDivisionByZero)
		assert.PanicsWith(t, func() {
			Seg(Pt(0.0, 0.0), Pt(8.0, 4.0)).UnscaleXY(0, 2)
		}, ErrDivisionByZero)
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

func TestSegment_Transform(t *testing.T) {
	t.Run("applies the matrix to both points", func(t *testing.T) {
		matrix := Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)

		geomtest.AssertSegment(t, Seg(Pt(1, 2), Pt(3, 4)).Transform(matrix), Seg(Pt(1, 2).Transform(matrix), Pt(3, 4).Transform(matrix)))
	})
	t.Run("float32 matrix", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(1.0, 2.0), Pt(0.0, 0.0)).Transform(Mat[float32](1, 0, 1, 0, 1, 1)), Seg(Pt(2.0, 3.0), Pt(1.0, 1.0)))
	})
}

func TestSegment_Lerp(t *testing.T) {
	a := Seg(Pt(0.0, 0.0), Pt(4.0, 0.0))
	b := Seg(Pt(10.0, 2.0), Pt(12.0, 8.0))

	t.Run("both ends move towards the ends of the other", func(t *testing.T) {
		geomtest.AssertSegment(t, a.Lerp(b, 0.5), Seg(Pt(5.0, 1.0), Pt(8.0, 4.0)))
	})
	t.Run("the ends are the two segments", func(t *testing.T) {
		geomtest.AssertSegment(t, a.Lerp(b, 0), a)
		geomtest.AssertSegment(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates outside [0, 1]", func(t *testing.T) {
		geomtest.AssertSegment(t, a.Lerp(b, 2), Seg(Pt(20.0, 4.0), Pt(20.0, 16.0)))
	})
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertSegment(t, Seg(Pt(0, 0), Pt(4, 0)).Lerp(Seg(Pt(1, 1), Pt(5, 3)), 0.5), Seg(Pt(1, 1), Pt(5, 2)))
	})
	t.Run("towards a translation is a part of it over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			geomtest.AssertSegment(t, s.Lerp(s.Translate(Vec(4.0, -2.0)), 0.25), s.Translate(Vec(1.0, -0.5)), s.String())
		}
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

func TestSegment_Split(t *testing.T) {
	t.Run("two segments meeting at the fraction", func(t *testing.T) {
		first, second := Seg(Pt(0, 0), Pt(8, 4)).Split(0.25)

		geomtest.AssertSegment(t, first, Seg(Pt(0, 0), Pt(2, 1)))
		geomtest.AssertSegment(t, second, Seg(Pt(2, 1), Pt(8, 4)))
	})
	t.Run("outside the ends one runs past and the other back", func(t *testing.T) {
		first, second := Seg(Pt(0.0, 0.0), Pt(4.0, 0.0)).Split(1.5)

		geomtest.AssertSegment(t, first, Seg(Pt(0.0, 0.0), Pt(6.0, 0.0)))
		geomtest.AssertSegment(t, second, Seg(Pt(6.0, 0.0), Pt(4.0, 0.0)))
	})
	t.Run("the two meet at PointAt and keep the ends and the length", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, fraction := range []float64{0, 0.3, 0.5, 1} {
				first, second := s.Split(fraction)
				message := fmt.Sprintf("%s @%v: ", s, fraction)

				geomtest.AssertPoint(t, first.Start, s.Start, message)
				geomtest.AssertPoint(t, first.End, s.PointAt(fraction), message)
				geomtest.AssertPoint(t, second.Start, first.End, message)
				geomtest.AssertPoint(t, second.End, s.End, message)
				geomtest.AssertNumber(t, first.Length()+second.Length(), s.Length(), message)
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

func TestSegment_NearestFraction(t *testing.T) {
	s := Seg(Pt(0, 0), Pt(4, 0))

	t.Run("how far along the foot of the perpendicular lies", func(t *testing.T) {
		geomtest.AssertNumber(t, s.NearestFraction(Pt(1, 3)), 0.25)
		geomtest.AssertNumber(t, s.NearestFraction(Pt(3, 0)), 0.75)
	})
	t.Run("zero before the start and one past the end", func(t *testing.T) {
		assert.Equal(t, s.NearestFraction(Pt(-3, 4)), 0.0)
		assert.Equal(t, s.NearestFraction(Pt(7, 4)), 1.0)
		assert.Equal(t, s.NearestFraction(Pt(0, 0)), 0.0)
		assert.Equal(t, s.NearestFraction(Pt(4, 0)), 1.0)
	})
	t.Run("a zero-length segment gives zero", func(t *testing.T) {
		assert.Equal(t, Seg(Pt(1, 1), Pt(1, 1)).NearestFraction(Pt(5, 4)), 0.0)
	})
	t.Run("a NaN coordinate is nearest to nothing", func(t *testing.T) {
		assert.True(t, math.IsNaN(s.Float().NearestFraction(Pt(math.NaN(), 0.0))))
	})
	t.Run("PointAt takes it back to the nearest point, and a point at a fraction has that fraction", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range pointFixtures {
				fraction := s.NearestFraction(p)
				message := fmt.Sprintf("%s → %s: ", s, p)

				assert.True(t, fraction >= 0 && fraction <= 1, message)
				assert.True(t, s.PointAt(fraction).Equal(s.Nearest(p)), message)
			}

			if s.IsPoint() {
				continue
			}

			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, s.Reverse().NearestFraction(p), 1-s.NearestFraction(p), fmt.Sprintf("%s → %s: ", s, p))
			}

			for _, fraction := range []float64{0, 0.1, 0.5, 0.8, 1} {
				geomtest.AssertNumber(t, s.NearestFraction(s.PointAt(fraction)), fraction, fmt.Sprintf("%s @%v: ", s, fraction))
			}
		}
	})
	t.Run("far from the origin a float32 point at the fraction lies on the segment", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, fixture := range segmentFixtures {
				s := fixture.Cast[float32]().Translate(offset)

				for _, p := range pointFixtures {
					assert.True(t, s.Contains(s.PointAt(s.NearestFraction(p.Cast[float32]().Add(offset)))), fmt.Sprintf("%s → %s: ", s, p))
				}
			}
		}
	})
}

func TestSegment_Support(t *testing.T) {
	segment := Seg(Pt(0, 0), Pt(4, 2))

	t.Run("the end the direction runs toward", func(t *testing.T) {
		geomtest.AssertPoint(t, segment.Support(Vec(1.0, 0.0)), Pt(4.0, 2.0))
		geomtest.AssertPoint(t, segment.Support(Vec(-1.0, 0.0)), Pt(0.0, 0.0))
	})
	t.Run("the start for a direction across it or of no length", func(t *testing.T) {
		geomtest.AssertPoint(t, segment.Support(Vec(-1.0, 2.0)), Pt(0.0, 0.0))
		geomtest.AssertPoint(t, segment.Support(Vec(0.0, 0.0)), Pt(0.0, 0.0))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		direction := Vec(tenth, 0.7)

		assertAllocations(t, func() {
			sinkBool = segment.Support(direction).X > 0
		}, 0)
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

func TestSegment_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, s.IntersectsCircle(c), c.IntersectsSegment(s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_IntersectsCapsule(t *testing.T) {
	t.Run("mirrors Capsule.IntersectsSegment", func(t *testing.T) {
		for _, shape := range segmentFixtures {
			for _, c := range capsuleFixtures {
				assert.Equal(t, shape.IntersectsCapsule(c), c.IntersectsSegment(shape), fmt.Sprintf("%s → %s: ", shape, c))
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

func TestSegment_IntersectsRay(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(Ry(Pt(0, 4), Vec(1, -1))))
	})
	t.Run("the lines cross behind the origin", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(Ry(Pt(3, 1), Vec(1, -1))))
	})
	t.Run("the lines cross beyond the segment", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsRay(Ry(Pt(10, 0), Vec(0, 1))))
	})
	t.Run("the origin on the segment counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(Ry(Pt(2, 2), Vec(1, 0))))
	})
	t.Run("an endpoint on the ray counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(Ry(Pt(0, 4), Vec(1, 0))))
	})
	t.Run("collinear toward the segment counts, away from it does not", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(Ry(Pt(9, 9), Vec(-1, -1))))
		assert.False(t, diagonal.IntersectsRay(Ry(Pt(5, 5), Vec(1, 1))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsRay(Ry(Pt(1, 1), Vec(0, 0))))
		assert.False(t, diagonal.IntersectsRay(Ry(Pt(1, 2), Vec(0, 0))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.True(t, s.IntersectsRay(Ry(Pt(0.5, Delta/2), Vec(0.0, 1.0))))
		assert.False(t, s.IntersectsRay(Ry(Pt(0.5, 2*Delta), Vec(0.0, 1.0))))
	})
	t.Run("matches a segment reaching past the segment along the ray", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, s.IntersectsRay(r), s.IntersectsSegment(far(r)), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
}

func TestSegment_IntersectsLine(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing, on either side of the origin of the line", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsLine(Ln(Pt(0, 4), Vec(1, -1))))
		assert.True(t, diagonal.IntersectsLine(Ln(Pt(3, 1), Vec(1, -1))))
	})
	t.Run("the lines cross beyond the segment", func(t *testing.T) {
		assert.False(t, diagonal.IntersectsLine(Ln(Pt(10, 0), Vec(0, 1))))
	})
	t.Run("an endpoint on the line counts", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsLine(Ln(Pt(9, 4), Vec(1, 0))))
	})
	t.Run("along the segment counts, beside it does not", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsLine(Ln(Pt(9, 9), Vec(1, 1))))
		assert.False(t, diagonal.IntersectsLine(Ln(Pt(9, 8), Vec(1, 1))))
	})
	t.Run("a zero direction is its origin", func(t *testing.T) {
		assert.True(t, diagonal.IntersectsLine(Ln(Pt(1, 1), Vec(0, 0))))
		assert.False(t, diagonal.IntersectsLine(Ln(Pt(1, 2), Vec(0, 0))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		s := Seg(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.True(t, s.IntersectsLine(Ln(Pt(5.0, Delta/2), Vec(1.0, 0.0))))
		assert.False(t, s.IntersectsLine(Ln(Pt(5.0, 2*Delta), Vec(1.0, 0.0))))
	})
	t.Run("matches a segment reaching past the segment along the line both ways", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, l := range lineFixtures {
				assert.Equal(t, s.IntersectsLine(l), s.IntersectsSegment(wide(l)), fmt.Sprintf("%s → %s: ", s, l))
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

func TestSegment_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4), 0)

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
	rectangle := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0), 0)
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

func TestSegment_IntersectsEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through, ending inside or inside", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(-8.0, 0.0), Pt(8.0, 0.0)), Seg(Pt(-8.0, 1.0), Pt(0.0, 1.0)), Seg(Pt(-1.0, 0.5), Pt(1.0, -0.5))} {
			assert.True(t, s.IntersectsEllipse(ellipse), s.String())
		}
	})
	t.Run("apart or short of the ellipse", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(-8.0, 3.0), Pt(8.0, 3.0)), Seg(Pt(-8.0, 0.0), Pt(-5.0, 0.0)), Seg(Pt(4.0, 2.0), Pt(6.0, 0.0))} {
			assert.False(t, s.IntersectsEllipse(ellipse), s.String())
		}
	})
	t.Run("a tangent touches, within the tolerance and not past it", func(t *testing.T) {
		assert.True(t, Seg(Pt(-5.0, 2.0), Pt(5.0, 2.0)).IntersectsEllipse(ellipse))
		assert.True(t, Seg(Pt(-5.0, 2.0+1e-12), Pt(5.0, 2.0+1e-12)).IntersectsEllipse(ellipse))
		assert.False(t, Seg(Pt(-5.0, 2.0+1e-4), Pt(5.0, 2.0+1e-4)).IntersectsEllipse(ellipse))
	})
	t.Run("a degenerate ellipse is the segment of its other axis", func(t *testing.T) {
		flat := Ell(Pt(0.0, 0.0), Sz(4.0, 0.0), 0)

		assert.True(t, Seg(Pt(0.0, -3.0), Pt(0.0, 3.0)).IntersectsEllipse(flat))
		assert.True(t, Seg(Pt(-6.0, 0.0), Pt(6.0, 0.0)).IntersectsEllipse(flat))
		assert.True(t, Seg(Pt(4.0+1e-12, -3.0), Pt(4.0+1e-12, 3.0)).IntersectsEllipse(flat))
		assert.False(t, Seg(Pt(5.0, -3.0), Pt(5.0, 3.0)).IntersectsEllipse(flat))
		assert.True(t, Seg(Pt(0.0, 0.0), Pt(2.0, 2.0)).IntersectsEllipse(Ell(Pt(1.0, 1.0), Sz(0.0, 0.0), 0)))
	})
	t.Run("int touches on the lattice", func(t *testing.T) {
		assert.True(t, Seg(Pt(-10, 2), Pt(10, 2)).IntersectsEllipse(Ell(Pt(0, 0), Sz(4, 2), 0)))
		assert.False(t, Seg(Pt(-10, 3), Pt(10, 3)).IntersectsEllipse(Ell(Pt(0, 0), Sz(4, 2), 0)))
	})
	t.Run("int touches on every exact tangent", func(t *testing.T) {
		for w := 1; w <= 20; w++ {
			for h := 1; h <= 20; h++ {
				for m := -10; m <= 10; m++ {
					c := int(math.Sqrt(float64(w*w*m*m + h*h)))
					if c*c != w*w*m*m+h*h {
						continue
					}

					s, e := Seg(Pt(-2*w, -2*w*m+c), Pt(2*w, 2*w*m+c)), Ell(Pt(0, 0), Sz(w, h), 0)
					assert.True(t, s.IntersectsEllipse(e), fmt.Sprintf("%s → %s: ", s, e))
				}
			}
		}
	})
	t.Run("int touches the tip of a flat ellipse", func(t *testing.T) {
		assert.True(t, Seg(Pt(6, -3), Pt(-3, 6)).IntersectsEllipse(Ell(Pt(0, 0), Sz(3, 0), 0)))
		assert.True(t, Seg(Pt(-8, -1), Pt(3, 10)).IntersectsEllipse(Ell(Pt(-4, 0), Sz(3, 0), 0)))
	})
	t.Run("int turned by a quarter turn matches the transposed ellipse", func(t *testing.T) {
		for _, s := range segmentFixtures {
			s := s.Int()
			for _, e := range ellipseFixtures {
				e := e.Int()
				transposed := Ell(e.Center, Sz(e.Radius.Height, e.Radius.Width), 0)

				assert.Equal(t, s.IntersectsEllipse(Ell(e.Center, e.Radius, Pi/2)), s.IntersectsEllipse(transposed), fmt.Sprintf("%s → %s: ", s, e))
				assert.Equal(t, s.IntersectsEllipse(Ell(e.Center, e.Radius, Pi)), s.IntersectsEllipse(Ell(e.Center, e.Radius, 0)), fmt.Sprintf("%s → %s: ", s, e))
			}
		}
		assert.True(t, Seg(Pt(4, 6), Pt(4, -4)).IntersectsEllipse(Ell(Pt(1, -2), Sz(3, 3), Pi)))
		assert.True(t, Seg(Pt(9, 3), Pt(-4, 3)).IntersectsEllipse(Ell(Pt(-3, -1), Sz(4, 1), Pi/2)))
	})
	t.Run("a NaN coordinate intersects nothing", func(t *testing.T) {
		assert.False(t, Seg(Pt(math.NaN(), 0.0), Pt(8.0, 0.0)).IntersectsEllipse(ellipse))
	})
	t.Run("matches IntersectsCircle on the circle's Ellipse", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, s.IntersectsEllipse(c.Ellipse()), s.IntersectsCircle(c), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("holds wherever a point along the segment lies in the ellipse, either way round", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, s := range ellipseCasts() {
				message := fmt.Sprintf("%s → %s: ", s, e)
				intersects := s.IntersectsEllipse(e)

				assert.Equal(t, s.Reverse().IntersectsEllipse(e), intersects, message)
				for i := range 33 {
					assert.True(t, !e.Contains(s.PointAt(float64(i)/32)) || intersects, message)
				}
			}
		}
	})
}

func TestSegment_IntersectionCircle(t *testing.T) {
	t.Run("matches Circle.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, s.IntersectionCircle(c).AppendTo(nil), c.IntersectionSegment(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_IntersectionCapsule(t *testing.T) {
	t.Run("matches Capsule.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range capsuleFixtures {
				geomtest.AssertVertices(t, s.IntersectionCapsule(c).AppendTo(nil), c.IntersectionSegment(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_IntersectionSector(t *testing.T) {
	t.Run("matches Sector.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, sector := range sectorFixtures() {
				geomtest.AssertVertices(t, s.IntersectionSector(sector), sector.IntersectionSegment(s), fmt.Sprintf("%s → %s: ", s, sector))
			}
		}
	})
}

func TestSegment_AppendIntersectionSector(t *testing.T) {
	sector, through := quarterSector(), Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

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

func TestSegment_IntersectionArc(t *testing.T) {
	t.Run("matches Arc.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, a := range arcFixtures {
				geomtest.AssertVertices(t, s.IntersectionArc(a).AppendTo(nil), a.IntersectionSegment(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, a))
			}
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
	t.Run("float segments apart on one line do not cross by rounding", func(t *testing.T) {
		pairs := [][2]Segment[float64]{
			{Seg(Pt(435.20098342749014, 249.93526920361418), Pt(269.47972773956184, 599.8089399849782)), Seg(Pt(442.2171109693957, 235.12269559881773), Pt(580.3540543853841, -56.51448312544401))},
			{Seg(Pt(71.83665157236865, 3.930612186721294), Pt(16.415664262469555, 35.56353823021138)), Seg(Pt(11.32413283789975, 38.469658182584624), Pt(-38.713033109816024, 67.02963242121056))},
		}

		for _, pair := range pairs {
			_, ok := pair[0].IntersectionSegment(pair[1])

			assert.False(t, ok, fmt.Sprintf("%s → %s: ", pair[0], pair[1]))
			assert.False(t, pair[0].IntersectsSegment(pair[1]), fmt.Sprintf("%s → %s: ", pair[0], pair[1]))
		}
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
	t.Run("nearly parallel segments with no endpoint on the other keep their crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0.0, 0.0), Pt(10.0, 0.0)).IntersectionSegment(Seg(Pt(0.0, 1.5*Delta), Pt(10.0, -1.5*Delta)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(5.0, 0.0))
	})
	t.Run("the same point in either order", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				point, ok := a.IntersectionSegment(b)
				reversed, reversedOk := b.IntersectionSegment(a)
				message := fmt.Sprintf("%s → %s: ", a, b)

				assert.Equal(t, ok, reversedOk, message)
				assert.Equal(t, point, reversed, message)
			}
		}
	})
	t.Run("a long int crossing on a half unit rounds the same way in either order", func(t *testing.T) {
		pairs := [][2]Segment[int]{
			{Seg(Pt(-67116, -108610), Pt(266866, 416662)), Seg(Pt(-290417, 594344), Pt(160846, -330808))},
			{Seg(Pt(524931, 492740), Pt(-108509, -103996)), Seg(Pt(123825, -454723), Pt(-103886, 395159))},
			{Seg(Pt(-14395822, -8911196), Pt(4330846, 2707652)), Seg(Pt(960221, -6190890), Pt(-864568, 5484762))},
			{Seg(Pt(28146470, 21943757), Pt(-24376732, -18988783)), Seg(Pt(31427391, -8045225), Pt(-63955104, 16459085))},
		}
		for _, pair := range pairs {
			point, ok := pair[0].IntersectionSegment(pair[1])
			reversed, reversedOk := pair[1].IntersectionSegment(pair[0])
			message := fmt.Sprintf("%s → %s: ", pair[0], pair[1])

			assert.True(t, ok, message)
			assert.True(t, reversedOk, message)
			assert.Equal(t, point, reversed, message)
		}
	})
	t.Run("a long int crossing on a half unit rounds the same way for a Reverse", func(t *testing.T) {
		s, edge := Seg(Pt(-2097166, 0), Pt(5242915, 2097166)), Seg(Pt(-1048583, 1048583), Pt(4194332, 1048583))

		for _, pair := range [][2]Segment[int]{{s, edge}, {s.Reverse(), edge}, {s, edge.Reverse()}, {edge.Reverse(), s.Reverse()}} {
			point, ok := pair[0].IntersectionSegment(pair[1])
			message := fmt.Sprintf("%s → %s: ", pair[0], pair[1])

			assert.True(t, ok, message)
			geomtest.AssertPoint(t, point, Pt(1572874, 1048583), message)
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

func TestSegment_IntersectionRay(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(Ry(Pt(0, 4), Vec(1, -1)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(Ry(Pt(3, 1), Vec(1, -1)))

		assert.False(t, ok)
	})
	t.Run("the origin on the segment is the point", func(t *testing.T) {
		point, ok := diagonal.IntersectionRay(Ry(Pt(2, 2), Vec(1, 0)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("collinear has no single point", func(t *testing.T) {
		_, ok := diagonal.IntersectionRay(Ry(Pt(9, 9), Vec(-1, -1)))

		assert.False(t, ok)
	})
	t.Run("a nearly collinear float ray falls back to the endpoint on it", func(t *testing.T) {
		s, r := Seg(Pt(42.400000000000006, 126.41), Pt(70.0, 209.35)), Ry(Pt(1.0, 2.0), Vec(13.8, 41.47))
		point, ok := s.IntersectionRay(r)

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, s.Start)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0, 0), Pt(3, 0)).IntersectionRay(Ry(Pt(1, -1), Vec(1, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 0))
	})
	t.Run("agrees with Intersects on non-parallel fixtures, on both", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rayFixtures {
				if s.Vector().Float().Cross(r.Vector.Float()) == 0 && s.Vector() != (Vector[float64]{}) && r.Vector != (Vector[float64]{}) {
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

func TestSegment_IntersectionLine(t *testing.T) {
	diagonal := Seg(Pt(0, 0), Pt(4, 4))

	t.Run("crossing, on either side of the origin of the line", func(t *testing.T) {
		for _, line := range []Line[int]{Ln(Pt(0, 4), Vec(1, -1)), Ln(Pt(3, 1), Vec(1, -1))} {
			point, ok := diagonal.IntersectionLine(line)

			assert.True(t, ok)
			geomtest.AssertPoint(t, point, Pt(2, 2))
		}
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := diagonal.IntersectionLine(Ln(Pt(10, 0), Vec(0, 1)))

		assert.False(t, ok)
	})
	t.Run("along the segment has no single point", func(t *testing.T) {
		_, ok := diagonal.IntersectionLine(Ln(Pt(9, 9), Vec(-1, -1)))

		assert.False(t, ok)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := Seg(Pt(0, 0), Pt(3, 0)).IntersectionLine(Ln(Pt(3, 3), Vec(1, 2)))

		assert.True(t, ok)
		geomtest.AssertPoint(t, point, Pt(2, 0))
	})
	t.Run("agrees with Intersects on non-parallel fixtures, on both", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, l := range lineFixtures {
				if s.Vector().Float().Cross(l.Vector.Float()) == 0 && s.Vector() != (Vector[float64]{}) && l.Vector != (Vector[float64]{}) {
					continue
				}

				point, ok := s.IntersectionLine(l)
				message := fmt.Sprintf("%s → %s: ", s, l)

				assert.Equal(t, ok, s.IntersectsLine(l), message)
				if ok {
					assert.True(t, s.Contains(point) && l.Contains(point), message)
				}
			}
		}
	})
}

func TestSegment_IntersectionPolyline(t *testing.T) {
	path := PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 3.0)})

	t.Run("the crossings of its edges in order, a shared point once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-1.0, -1.0), Pt(5.0, 2.0)).IntersectionPolyline(path), []Point[float64]{Pt(1.0, 0.0), Pt(4.0, 1.5)})
		geomtest.AssertVertices(t, Seg(Pt(0.0, 2.0), Pt(6.0, -1.0)).IntersectionPolyline(path), []Point[float64]{Pt(4.0, 0.0)})
		assert.Nil(t, Seg(Pt(0.0, 5.0), Pt(-3.0, 5.0)).IntersectionPolyline(path))
		assert.Nil(t, Seg(Pt(0.0, 1.0), Pt(3.0, 1.0)).IntersectionPolyline(Polyline[float64]{}))
	})
	t.Run("a closed path is crossed as the polygon of its points", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, s.IntersectionPolyline(p.Polyline()), s.IntersectionPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("appends after what the buffer holds", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-1.0, -1.0), Pt(5.0, 2.0)).AppendIntersectionPolyline(bufferWith(prefixPoint), path), []Point[float64]{prefixPoint, Pt(1.0, 0.0), Pt(4.0, 1.5)})
	})
}

func TestSegment_IntersectionQuadraticBezier(t *testing.T) {
	t.Run("matches QuadraticBezier.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, q := range quadraticFixtures {
				geomtest.AssertVertices(t, s.IntersectionQuadraticBezier(q).AppendTo(nil), q.IntersectionSegment(s).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, q))
			}
		}
	})
}

func TestSegment_IntersectionCubicBezier(t *testing.T) {
	t.Run("matches CubicBezier.IntersectionSegment", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range cubicFixtures {
				geomtest.AssertVertices(t, s.IntersectionCubicBezier(c), c.IntersectionSegment(s), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_AppendIntersectionCubicBezier(t *testing.T) {
	wave, through := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0)), Seg(Pt(-1.0, tenth), Pt(7.0, 0.0))

	t.Run("appends after the points already there", func(t *testing.T) {
		geomtest.AssertVertices(t, through.AppendIntersectionCubicBezier(bufferWith(prefixPoint), wave), append([]Point[float64]{prefixPoint}, through.IntersectionCubicBezier(wave)...))
	})
	t.Run("appends to a buffer with room without allocating", func(t *testing.T) {
		buffer := make([]Point[float64], 0, 4)

		assertAllocations(t, func() {
			sinkBool = len(through.AppendIntersectionCubicBezier(buffer[:0], wave)) > 0
		}, 0)
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
	t.Run("passing a vertex within the tolerance meets it once", func(t *testing.T) {
		triangle := Pol([]Point[float64]{Pt(10.0, 0.0), Pt(0.0, 0.0), Pt(0.0, 1.0)})

		assert.Length(t, Seg(Pt(-1.0, 4e-7), Pt(12.0, 6e-7)).IntersectionPolygon(triangle), 2)
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

		assertAllocations(t, func() {
			sinkPoints = through.IntersectionPolygon(square)
		}, 1)
		assertAllocations(t, func() {
			sinkPoints = apart.IntersectionPolygon(square)
		}, 0)
	})
	t.Run("matches the rectangle crossings on the rectangle as a polygon", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, s.IntersectionPolygon(r.Polygon()), s.IntersectionRectangle(r).AppendTo(nil), fmt.Sprintf("%s → %s: ", r, s))
			}
		}
	})
	t.Run("every point lies on the segment and an edge, and exists where IntersectsPolygon holds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				assertCrossings(t, s, s.IntersectionPolygon(p), s.IntersectsPolygon(p), p.Edges, fmt.Sprintf("%v → %v: ", s, p))
			}
		}
	})
	t.Run("a segment and its Reverse cross at the same points", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				assertReversed(t, s.IntersectionPolygon(p), s.Reverse().IntersectionPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and an edge", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, p := range polygonFixtures() {
				for _, s := range segmentFixtures {
					p, s := p.Cast[float32]().Translate(offset), s.Cast[float32]().Translate(offset)

					assertCrossings(t, s, s.IntersectionPolygon(p), s.IntersectsPolygon(p), p.Edges, fmt.Sprintf("%v → %v: ", s, p))
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

		assertAllocations(t, func() {
			sinkPoints = through.AppendIntersectionPolygon(buffer[:0], square)
		}, 0)
	})
	t.Run("matches IntersectionPolygon after the points in dst", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, p := range polygonFixtures() {
				geomtest.AssertVertices(t, s.AppendIntersectionPolygon(bufferWith(prefixPoint), p), append([]Point[float64]{prefixPoint}, s.IntersectionPolygon(p)...), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4), 0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(5, 0), Pt(-5, 0)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(0, 0), Pt(5, 0)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(2, 0)})
	})
	t.Run("through a corner counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(2, 2)})
		geomtest.AssertVertices(t, Seg(Pt(-4, -4), Pt(4, 4)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(-2, -2), Pt(2, 2)})
	})
	t.Run("an end within the band of a long edge just past a corner is the exit, the end ClipRectangle clips to", func(t *testing.T) {
		rect := Rectangle[float32]{Center: Pt[float32](-50001.85, 10000.37), Size: Sz[float32](20000.74, 60002.22), Angle: 1.5707958391575285}
		s := Seg(Pt[float32](-60002.22, 40001.48), Pt[float32](-20000.73, -0.0073))
		clipped, ok := s.ClipRectangle(rect)

		assert.True(t, ok)
		geomtest.AssertVertices(t, s.IntersectionRectangle(rect).AppendTo(nil), []Point[float32]{clipped.Start, clipped.End})
	})
	t.Run("an endpoint exactly Delta outside is judged like IntersectsRectangle", func(t *testing.T) {
		s, r := Seg(Pt(-2.0, 2.000001), Pt(2.0, 68.0000005)), Rect(Pt(0.0, 0.0), Sz(26.0, 4.0), 0)

		assert.Equal(t, len(s.IntersectionRectangle(r).AppendTo(nil)) > 0, s.IntersectsRectangle(r))
	})
	t.Run("along an edge crosses the edges at its ends", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, -2), Pt(5, -2)).IntersectionRectangle(rectangle).AppendTo(nil), []Point[int]{Pt(-2, -2), Pt(2, -2)})
		geomtest.AssertVertices(t, Seg(Pt(-1, 2), Pt(1, 2)).IntersectionRectangle(rectangle).AppendTo(nil), nil)
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectionRectangle(rectangle).AppendTo(nil))
		assert.Nil(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectionRectangle(rectangle).AppendTo(nil))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through, apart := Seg(Pt(-5, 1), Pt(5, 1)), Seg(Pt(3, -5), Pt(3, 5))

		assertAllocations(t, func() {
			sinkBool = through.IntersectionRectangle(rectangle).Count > 0
		}, 0)
		assertAllocations(t, func() {
			sinkBool = through.IntersectionRectangle(rectangle.Rotate(Pi/5)).Count > 0
		}, 0)
		assertAllocations(t, func() {
			sinkBool = apart.IntersectionRectangle(rectangle).Count > 0
		}, 0)
	})
	t.Run("float keeps the crossings", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5.0, 1.0), Pt(5.0, 1.0)).IntersectionRectangle(Rect(Pt(0.0, 0.0), Sz(3.0, 3.0), 0)).AppendTo(nil), []Point[float64]{Pt(-1.5, 1.0), Pt(1.5, 1.0)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsRectangle holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rectFixtures {
				assertCrossings(t, s, s.IntersectionRectangle(r).AppendTo(nil), s.IntersectsRectangle(r), r.Edges, fmt.Sprintf("%v → %v: ", s, r))
			}
		}
	})
	t.Run("a segment and its Reverse cross at the same points", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, r := range rectFixtures {
				assertReversed(t, s.IntersectionRectangle(r).AppendTo(nil), s.Reverse().IntersectionRectangle(r).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, r := range rectFixtures {
					s, r := s.Cast[float32]().Translate(offset), r.Cast[float32]().Translate(offset)

					assertCrossings(t, s, s.IntersectionRectangle(r).AppendTo(nil), s.IntersectsRectangle(r), r.Edges, fmt.Sprintf("%v → %v: ", s, r))
				}
			}
		}
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the points follow from Start", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				r, s := r.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				points := s.IntersectionRectangle(r).AppendTo(nil)

				assert.True(t, len(points) == 0 || s.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", s, r))
				assertOrdered(t, s, points, fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
	t.Run("a segment a rounding error off an edge crosses only the edges at its ends", func(t *testing.T) {
		r := Rectangle[float64]{Pt(-9.555794994732457, -7.1447161806712565), Sz(1.9672661593065617, 2.772250161572721), 2.3562338764659856}
		corners := slices.Collect(r.Vertices())
		s := Seg(Pt(-7.1845731257763275, -7.555551584844245), Pt(-9.966817184817401, -4.773526681623244))

		geomtest.AssertVertices(t, s.IntersectionRectangle(r).AppendTo(nil), corners[:2])
		geomtest.AssertVertices(t, s.IntersectionPolygon(r.Polygon()), corners[:2])
	})
	t.Run("a segment along an edge of the fixtures, extended past both ends, crosses at most twice", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, turn := range []float64{0, 0.3, 2.3562338764659856} {
				r := r.Rotate(turn)

				for edge := range r.Edges() {
					s := Seg(edge.PointAt(-0.5), edge.PointAt(1.5))

					assert.True(t, len(s.IntersectionRectangle(r).AppendTo(nil)) <= 2, fmt.Sprintf("%s → %s: ", s, r))
				}
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
	f.Add(-12.0, -2.0+5e-7, 12.0, -2.0+7e-7, 0.0, 0.0, 20.0, 4.0)

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, cx, cy, w, h float64) {
		for _, v := range []float64{x1, y1, x2, y2, cx, cy, w, h} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		s, r := Seg(Pt(x1, y1), Pt(x2, y2)), Rect(Pt(cx, cy), Sz(w, h), 0)
		points := s.IntersectionRectangle(r).AppendTo(nil)

		assert.True(t, len(points) <= 2, fmt.Sprintf("%s → %s: at most two crossings, got %d: ", s, r, len(points)))

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

func TestSegment_IntersectionRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectionRegularPolygon(diamond).AppendTo(nil), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(3, 0), Pt(-3, 0)).IntersectionRegularPolygon(diamond).AppendTo(nil), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("a vertex hit by two edges is counted once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(2, -2), Pt(2, 2)).IntersectionRegularPolygon(diamond).AppendTo(nil), []Point[int]{Pt(2, 0)})
	})
	t.Run("a segment passing vertices within the tolerance meets each once", func(t *testing.T) {
		points := Seg(Pt(-12.0, 5e-7), Pt(12.0, 7e-7)).IntersectionRegularPolygon(RegPol(Pt(0.0, 0.0), Sz(10.0, 1.0), 4, 0, 0)).AppendTo(nil)

		assert.Length(t, points, 2)
	})
	t.Run("inside, apart and empty give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(0, 0), Pt(1, 0)).IntersectionRegularPolygon(diamond).AppendTo(nil))
		assert.Nil(t, Seg(Pt(5, 0), Pt(6, 0)).IntersectionRegularPolygon(diamond).AppendTo(nil))
		assert.Nil(t, Seg(Pt(-3, 0), Pt(3, 0)).IntersectionRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)).AppendTo(nil))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through, apart := Seg(Pt(-3, 0), Pt(3, 0)), Seg(Pt(5, 0), Pt(6, 0))

		assertAllocations(t, func() {
			sinkBool = through.IntersectionRegularPolygon(diamond).Count > 0
		}, 0)
		assertAllocations(t, func() {
			sinkBool = apart.IntersectionRegularPolygon(diamond).Count > 0
		}, 0)
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				geomtest.AssertVertices(t, s.IntersectionRegularPolygon(rp).AppendTo(nil), s.IntersectionPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("a segment and its Reverse cross at the same points", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, rp := range regularPolygonFixtures {
				assertReversed(t, s.IntersectionRegularPolygon(rp).AppendTo(nil), s.Reverse().IntersectionRegularPolygon(rp).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("over the int16 fixtures, spanning past the square root of its range, the points follow from Start", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				rp, s := rp.Transform(narrowMatrix).Cast[int16](), s.Transform(narrowMatrix).Cast[int16]()
				points := s.IntersectionRegularPolygon(rp).AppendTo(nil)

				assert.True(t, len(points) == 0 || s.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", s, rp))
				assertOrdered(t, s, points, fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_IntersectionBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5, 0), Pt(5, 0)).IntersectionBox(box).AppendTo(nil), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		geomtest.AssertVertices(t, Seg(Pt(5, 0), Pt(-5, 0)).IntersectionBox(box).AppendTo(nil), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("through a corner counts it once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(1, 3), Pt(3, 1)).IntersectionBox(box).AppendTo(nil), []Point[int]{Pt(2, 2)})
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Seg(Pt(3, -5), Pt(3, 5)).IntersectionBox(box).AppendTo(nil))
		assert.Nil(t, Seg(Pt(-1, -1), Pt(1, 1)).IntersectionBox(box).AppendTo(nil))
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsBox holds", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, b := range boxFixtures {
				assertCrossings(t, s, s.IntersectionBox(b).AppendTo(nil), s.IntersectsBox(b), b.Rectangle().Edges, fmt.Sprintf("%v → %v: ", s, b))
			}
		}
	})
	t.Run("a segment and its Reverse cross at the same points", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, b := range boxFixtures {
				assertReversed(t, s.IntersectionBox(b).AppendTo(nil), s.Reverse().IntersectionBox(b).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, b))
			}
		}
	})
	t.Run("a long int segment and its Reverse round a half unit alike", func(t *testing.T) {
		s, b := Seg(Pt(-2097166, 0), Pt(5242915, 2097166)), BoxFromMinMax(Pt(-1048583, 0), Pt(4194332, 1048583))

		assertReversed(t, s.IntersectionBox(b).AppendTo(nil), s.Reverse().IntersectionBox(b).AppendTo(nil), "")
	})
	t.Run("far from the origin every float32 point lies on the segment and the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, s := range segmentFixtures {
				for _, b := range boxFixtures {
					s, b := s.Cast[float32]().Translate(offset), b.Cast[float32]().Translate(offset)

					assertCrossings(t, s, s.IntersectionBox(b).AppendTo(nil), s.IntersectsBox(b), b.Rectangle().Edges, fmt.Sprintf("%v → %v: ", s, b))
				}
			}
		}
	})
}

func TestSegment_IntersectionEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-8.0, 0.0), Pt(8.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(-4.0, 0.0), Pt(4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(8.0, 1.6), Pt(-8.0, 1.6)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(2.4, 1.6), Pt(-2.4, 1.6)})
	})
	t.Run("a turned ellipse turns its crossings with it", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5.0, 2.0), Pt(7.0, 2.0)).IntersectionEllipse(Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/2)).AppendTo(nil), []Point[float64]{Pt(-1.0, 2.0), Pt(3.0, 2.0)})
	})
	t.Run("an end inside gives the one crossing, and inside none", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-8.0, 0.0), Pt(0.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(-4.0, 0.0)})
		assert.True(t, Seg(Pt(-1.0, 0.5), Pt(1.0, -0.5)).IntersectionEllipse(ellipse).IsEmpty())
	})
	t.Run("an endpoint on the boundary is one of the points", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(4.0, 0.0), Pt(9.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(0.0, 0.0), Pt(4.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(4.0, 0.0), Pt(-4.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0), Pt(-4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(-2.0, 0.0), Pt(-4.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(-4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(-8.0, 0.0), Pt(-4.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(-4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(4.0, 0.0), Pt(8.0, 0.0)).Reverse().IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0)})
	})
	t.Run("a tangent is the point of contact, placed within the tolerance of the boundary", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-5.0, 2.0), Pt(5.0, 2.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(0.0, 2.0)})
		geomtest.AssertVertices(t, Seg(Pt(-5.0, 2.0+1e-12), Pt(5.0, 2.0+1e-12)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(0.0, 2.0)})

		turned := Ell(Pt(0.0, 0.0), Sz(4.0, 1.0), Pi/4)
		touch := turned.Support(Vec(-1.0, 1.0))
		for _, offset := range []float64{0, 1e-12} {
			edge := Pt(touch.X-offset, touch.Y+offset)
			points := Seg(edge.Add(Vec(-3.0, -3.0)), edge.Add(Vec(3.0, 3.0))).IntersectionEllipse(turned)

			assert.Equal(t, points.Count, 1)
			assert.True(t, turned.Contains(points.Items[0]), points.Items[0].String())
		}
	})
	t.Run("apart or of no length crosses nothing", func(t *testing.T) {
		assert.True(t, Seg(Pt(-8.0, 3.0), Pt(8.0, 3.0)).IntersectionEllipse(ellipse).IsEmpty())
		assert.True(t, Seg(Pt(-8.0, 0.0), Pt(-8.0, 0.0)).IntersectionEllipse(ellipse).IsEmpty())
		geomtest.AssertVertices(t, Seg(Pt(4.0, 0.0), Pt(4.0, 0.0)).IntersectionEllipse(ellipse).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0)})
	})
	t.Run("a degenerate ellipse is crossed as the segment of its other axis", func(t *testing.T) {
		flat := Ell(Pt(0.0, 0.0), Sz(4.0, 0.0), 0)

		geomtest.AssertVertices(t, Seg(Pt(0.0, -3.0), Pt(0.0, 3.0)).IntersectionEllipse(flat).AppendTo(nil), []Point[float64]{Pt(0.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(-6.0, 0.0), Pt(6.0, 0.0)).IntersectionEllipse(flat).AppendTo(nil), []Point[float64]{Pt(-4.0, 0.0), Pt(4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(4.0+1e-12, -3.0), Pt(4.0+1e-12, 3.0)).IntersectionEllipse(flat).AppendTo(nil), []Point[float64]{Pt(4.0, 0.0)})
		geomtest.AssertVertices(t, Seg(Pt(0.0, 0.0), Pt(2.0, 2.0)).IntersectionEllipse(Ell(Pt(1.0, 1.0), Sz(0.0, 0.0), 0)).AppendTo(nil), []Point[float64]{Pt(1.0, 1.0)})
		geomtest.AssertVertices(t, Seg(Pt(-1.0, 1.0), Pt(3.0, 1.0)).IntersectionEllipse(Ell(Pt(1.0, 1.0), Sz(0.0, 0.0), 0)).AppendTo(nil), []Point[float64]{Pt(1.0, 1.0)})
	})
	t.Run("a NaN coordinate crosses nothing", func(t *testing.T) {
		assert.True(t, Seg(Pt(math.NaN(), 0.0), Pt(8.0, 0.0)).IntersectionEllipse(ellipse).IsEmpty())
	})
	t.Run("int rounds the crossings once", func(t *testing.T) {
		geomtest.AssertVertices(t, Seg(Pt(-10, 1), Pt(10, 1)).IntersectionEllipse(Ell(Pt(0, 0), Sz(4, 2), 0)).AppendTo(nil), []Point[int]{Pt(-3, 1), Pt(3, 1)})
		geomtest.AssertVertices(t, Seg(Pt(-10, 2), Pt(10, 2)).IntersectionEllipse(Ell(Pt(0, 0), Sz(4, 2), 0)).AppendTo(nil), []Point[int]{Pt(0, 2)})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-8.0, 1.0), Pt(8.0, 1.0))

		assertAllocations(t, func() {
			sinkBool = through.IntersectionEllipse(ellipse).IsEmpty()
		}, 0)
	})
	t.Run("matches IntersectionCircle on the circle's Ellipse", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertVertices(t, s.IntersectionEllipse(c.Ellipse()).AppendTo(nil), s.IntersectionCircle(c).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("every point lies on the segment and the boundary, in order and the same reversed", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, s := range ellipseCasts() {
				message := fmt.Sprintf("%s → %s: ", s, e)
				points := s.IntersectionEllipse(e).AppendTo(nil)

				assertEllipseCrossings(t, s, points, s.IntersectsEllipse(e), e, message)
				assertReversed(t, points, s.Reverse().IntersectionEllipse(e).AppendTo(nil), message)
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

		assertAllocations(t, func() {
			_, sinkBool = through.ClipCircle(unit)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				assertClipped(t, s, partsOf(s.ClipCircle(c)), s.IntersectsCircle(c), c.ContainsSegment, fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("a segment and its Reverse are clipped to the same part", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range append(grazingSegments(c), segmentFixtures...) {
				parts := partsOf(s.Reverse().ClipCircle(c))
				for i, part := range parts {
					parts[i] = part.Reverse()
				}

				geomtest.AssertSegments(t, parts, partsOf(s.ClipCircle(c)), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_ClipCapsule(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(5.0, 10.0), Pt(5.0, -10.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(5.0, 2.0), Pt(5.0, -2.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-5.0, 0.0), Pt(15.0, 0.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(-2.0, 0.0), Pt(12.0, 0.0))})
	})
	t.Run("an endpoint inside is kept, and inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(5.0, 0.0), Pt(5.0, 10.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(5.0, 0.0), Pt(5.0, 2.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(1.0, 1.0), Pt(9.0, -1.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(9.0, -1.0))})
	})
	t.Run("a tangent is the point of contact, and apart gives none", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(12.0, -5.0), Pt(12.0, 5.0)).ClipCapsule(capsule)), []Segment[float64]{Seg(Pt(12.0, 0.0), Pt(12.0, 0.0))})
		assert.Nil(t, partsOf(Seg(Pt(-5.0, 3.0), Pt(15.0, 3.0)).ClipCapsule(capsule)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(5.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipCapsule(capsule)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)

					assertClipped(t, s, partsOf(s.ClipCapsule(c)), s.IntersectsCapsule(c), c.ContainsSegment, fmt.Sprintf("%s → %s: ", s, c))
				}
			}
		}
	})
}

func TestSegment_ClipSector(t *testing.T) {
	sector := Sect(Pt(0.0, 0.0), 5.0, 0, Pi/2)
	notched := Sect(Pt(0.0, 0.0), 5.0, Pi/4, 3*Pi/2)

	t.Run("passing through gives the part between the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(3.0, 10.0), Pt(3.0, -10.0)).ClipSector(sector), []Segment[float64]{Seg(Pt(3.0, 4.0), Pt(3.0, 0.0))})
	})
	t.Run("an endpoint inside is kept, and inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(3.0, 1.0), Pt(3.0, 10.0)).ClipSector(sector), []Segment[float64]{Seg(Pt(3.0, 1.0), Pt(3.0, 4.0))})
		geomtest.AssertSegments(t, Seg(Pt(1.0, 1.0), Pt(2.0, 3.0)).ClipSector(sector), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(2.0, 3.0))})
	})
	t.Run("through the notch of a sector past half a turn there are two parts", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(3.0, 10.0), Pt(3.0, -10.0)).ClipSector(notched), []Segment[float64]{Seg(Pt(3.0, 4.0), Pt(3.0, 3.0)), Seg(Pt(3.0, -3.0), Pt(3.0, -4.0))})
	})
	t.Run("touching from outside is that point, and apart gives none", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(5.0, -1.0), Pt(5.0, 1.0)).ClipSector(sector), []Segment[float64]{Seg(Pt(5.0, 0.0), Pt(5.0, 0.0))})
		assert.Nil(t, Seg(Pt(-1.0, -1.0), Pt(-1.0, 9.0)).ClipSector(sector))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(3, 10), Pt(3, -10)).ClipSector(Sect(Pt(0, 0), 5, 0, Pi/2)), []Segment[int]{Seg(Pt(3, 4), Pt(3, 0))})
	})
	t.Run("appends to a buffer with room without allocating", func(t *testing.T) {
		buffer, through := make([]Segment[float64], 0, 4), Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			sinkBool = len(through.AppendClipSector(buffer[:0], notched)) > 0
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, fixture := range sectorFixtures() {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)

					assertClipped(t, s, s.ClipSector(fixture), s.IntersectsSector(fixture), fixture.ContainsSegment, fmt.Sprintf("%s → %s: ", s, fixture))
				}
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
	t.Run("an entry rounded onto a Start outside the polygon opens the part there", func(t *testing.T) {
		wedge := Pol([]Point[int]{Pt(0, 0), Pt(10, 0), Pt(0, 4)})
		segment := Seg(Pt(6, 2), Pt(6, -1))

		assert.False(t, wedge.Contains(segment.Start))
		geomtest.AssertSegments(t, segment.ClipPolygon(wedge), []Segment[int]{Seg(Pt(6, 2), Pt(6, 0))})
		geomtest.AssertSegments(t, segment.Reverse().ClipPolygon(wedge), []Segment[int]{Seg(Pt(6, 0), Pt(6, 2))})
	})
	t.Run("an entry that compares Equal to a float Start outside the tolerance opens the part there", func(t *testing.T) {
		diamond := Pol([]Point[float64]{Pt(0.0, -10.0), Pt(10.0, 0.0), Pt(0.0, 10.0), Pt(-10.0, 0.0)})
		segment := Seg(Pt(-5-0.85e-6, -5-0.85e-6), Pt(20.0, 20.0))

		assert.False(t, diamond.Contains(segment.Start))
		geomtest.AssertSegments(t, segment.ClipPolygon(diamond), []Segment[float64]{Seg(Pt(-5.0, -5.0), Pt(5.0, 5.0))})
		geomtest.AssertSegments(t, segment.Reverse().ClipPolygon(diamond), []Segment[float64]{Seg(Pt(5.0, 5.0), Pt(-5.0, -5.0))})
	})
	t.Run("inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, Seg(Pt(1, 1), Pt(1, 2)).ClipPolygon(square), []Segment[int]{Seg(Pt(1, 1), Pt(1, 2))})
		geomtest.AssertSegments(t, Seg(Pt(1, 1), Pt(1, 1)).ClipPolygon(square), []Segment[int]{Seg(Pt(1, 1), Pt(1, 1))})
	})
	t.Run("leaving an edge at a shallow angle ends the part at a point the polygon holds", func(t *testing.T) {
		room := RectangleFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		segment := Seg(Pt(0.0, 0.0), Pt(8.0, -1.5e-6))
		clipped, _ := segment.ClipRectangle(room)

		geomtest.AssertSegments(t, segment.ClipPolygon(room.Polygon()), []Segment[float64]{clipped})
		geomtest.AssertSegments(t, segment.Reverse().ClipPolygon(room.Polygon()), []Segment[float64]{clipped})
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

		assertAllocations(t, func() {
			sinkSegments = through.ClipPolygon(square)
		}, 1)
		assertAllocations(t, func() {
			sinkSegments = apart.ClipPolygon(square)
		}, 0)
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
				assertClipped(t, s, s.ClipPolygon(p), s.IntersectsPolygon(p), p.ContainsSegment, fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
	t.Run("a rounded segment and its Reverse are clipped to the same parts, near the origin and near 2^26", func(t *testing.T) {
		for _, offset := range []Vector[int]{{}, Vec(1<<25, -(1 << 25))} {
			for _, p := range outlineFixtures() {
				for _, s := range segmentFixtures {
					p, s := p.Int().Translate(offset), s.Int().Translate(offset)

					geomtest.AssertSegments(t, reversedParts(s.Reverse().ClipPolygon(p)), s.ClipPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
				}
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
	t.Run("a NaN or infinite End ends the pass", func(t *testing.T) {
		square := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(10.0, 0.0), Pt(10.0, 10.0), Pt(0.0, 10.0)})

		for _, end := range []Point[float64]{Pt(math.NaN(), 5.0), Pt(math.Inf(1), 5.0)} {
			geomtest.AssertSegments(t, Seg(Pt(1.0, 1.0), end).ClipPolygon(square), []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(1.0, 1.0))}, fmt.Sprintf("to %s: ", end))
		}
	})
}

// reversedParts returns the parts of a clipped segment as its Reverse would give them: in the
// opposite order, each running the other way.
func reversedParts[T Number](parts []Segment[T]) []Segment[T] {
	reversed := make([]Segment[T], len(parts))
	for i, part := range parts {
		reversed[len(parts)-1-i] = part.Reverse()
	}

	return reversed
}

func FuzzSegment_ClipPolygon(f *testing.F) {
	f.Add(0.0, 0.0, 8.0, -1.5e-6)
	f.Add(-1.0, 3.0, 5.0, 3.0)
	f.Add(2.0, 1.0, 2.0, 5.0)
	f.Add(4.0, 0.0, 4.0+1e-7, 4.0)

	notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 1.0), Pt(0.0, 4.0)})

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2 float64) {
		for _, v := range []float64{x1, y1, x2, y2} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		s := Seg(Pt(x1, y1), Pt(x2, y2))
		parts := s.ClipPolygon(notched)

		assert.Equal(t, len(parts) > 0, s.IntersectsPolygon(notched), fmt.Sprintf("%s: ", s))
		for _, part := range parts {
			message := fmt.Sprintf("%s → %s: ", s, part)

			assert.True(t, notched.Contains(part.Start) && notched.Contains(part.End), message+"ends in the polygon: ")
			assert.True(t, s.Contains(part.Start) && s.Contains(part.End), message+"ends on the segment: ")
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

		assertAllocations(t, func() {
			sinkSegments = through.AppendClipPolygon(buffer[:0], square)
		}, 0)
	})
	t.Run("more crossings than one pass holds are found by passing again, without allocating", func(t *testing.T) {
		const teeth = 12

		comb := []Point[int]{Pt(0, 0), Pt(4*teeth-2, 0)}
		parts := make([]Segment[int], 0, teeth)
		for tooth := teeth - 1; tooth >= 0; tooth-- {
			left := 4 * tooth

			comb = append(comb, Pt(left+2, 2), Pt(left+2, 10), Pt(left, 10), Pt(left, 2))
			parts = append(parts, Seg(Pt(left+2, 5), Pt(left, 5)))
		}

		across := Seg(Pt(4*teeth, 5), Pt(-2, 5))
		buffer := make([]Segment[int], 0, teeth)

		geomtest.AssertSegments(t, across.AppendClipPolygon(buffer[:0], Pol(comb)), parts)

		slices.Reverse(parts)
		for i, part := range parts {
			parts[i] = part.Reverse()
		}

		geomtest.AssertSegments(t, across.Reverse().AppendClipPolygon(buffer[:0], Pol(comb)), parts)
		assertAllocations(t, func() {
			sinkSegments = across.AppendClipPolygon(buffer[:0], Pol(comb))
		}, 0)
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
	rectangle := Rect(Pt(0, 0), Sz(4, 4), 0)

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
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0), 0).Rotate(Pi / 4)

		geomtest.AssertSegments(t, partsOf(Seg(Pt(-3.0, 0.0), Pt(3.0, 0.0)).ClipRectangle(diamond)), []Segment[float64]{Seg(Pt(-Sqrt2, 0.0), Pt(Sqrt2, 0.0))})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-5, 1), Pt(5, 1))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipRectangle(rectangle)
		}, 0)
		assertAllocations(t, func() {
			_, sinkBool = through.ClipRectangle(rectangle.Rotate(Pi / 5))
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				assertClipped(t, s, partsOf(s.ClipRectangle(r)), s.IntersectsRectangle(r), r.ContainsSegment, fmt.Sprintf("%s → %s: ", s, r))
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

		assertAllocations(t, func() {
			_, sinkBool = through.ClipRegularPolygon(diamond)
		}, 0)
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
				assertClipped(t, s, partsOf(s.ClipRegularPolygon(rp)), s.IntersectsRegularPolygon(rp), rp.ContainsSegment, fmt.Sprintf("%s → %s: ", s, rp))
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

		assertAllocations(t, func() {
			_, sinkBool = through.ClipBox(box)
		}, 0)
	})
}

func TestSegment_ClipEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("passing through gives the chord from Start to End", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-8.0, 0.0), Pt(8.0, 0.0)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(-4.0, 0.0), Pt(4.0, 0.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(8.0, 1.6), Pt(-8.0, 1.6)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(2.4, 1.6), Pt(-2.4, 1.6))})
	})
	t.Run("an endpoint inside is kept, and inside is the segment itself", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(1.0, 0.0), Pt(8.0, 0.0)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(1.0, 0.0), Pt(4.0, 0.0))})
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-1.0, 0.5), Pt(1.0, -0.5)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(-1.0, 0.5), Pt(1.0, -0.5))})
	})
	t.Run("a tangent is the point of contact", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-5.0, 2.0), Pt(5.0, 2.0)).ClipEllipse(ellipse)), []Segment[float64]{Seg(Pt(0.0, 2.0), Pt(0.0, 2.0))})
	})
	t.Run("apart gives none", func(t *testing.T) {
		assert.Nil(t, partsOf(Seg(Pt(-8.0, 3.0), Pt(8.0, 3.0)).ClipEllipse(ellipse)))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		geomtest.AssertSegments(t, partsOf(Seg(Pt(-10, 1), Pt(10, 1)).ClipEllipse(Ell(Pt(0, 0), Sz(4, 2), 0))), []Segment[int]{Seg(Pt(-3, 1), Pt(3, 1))})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-8.0, 1.0), Pt(8.0, 1.0))

		assertAllocations(t, func() {
			_, sinkBool = through.ClipEllipse(ellipse)
		}, 0)
	})
	t.Run("matches ClipCircle on the circle's Ellipse", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, c := range circleFixtures {
				geomtest.AssertSegments(t, partsOf(s.ClipEllipse(c.Ellipse())), partsOf(s.ClipCircle(c)), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, s := range ellipseCasts() {
				assertClipped(t, s, partsOf(s.ClipEllipse(e)), s.IntersectsEllipse(e), func(part Segment[float64]) bool {
					return e.Contains(part.Start) && e.Contains(part.End)
				}, fmt.Sprintf("%s → %s: ", s, e))
			}
		}
	})
}

func TestSegment_HitCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("from outside is the first crossing, with the normal from the center", func(t *testing.T) {
		hit, ok := Seg(Pt(-3.0, 0.0), Pt(1.0, 0.0)).HitCircle(circle)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(-1.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})

		hit, ok = Seg(Pt(-3.0, 0.6), Pt(3.0, 0.6)).HitCircle(circle)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(-0.8, 0.6), Normal: Vec(-0.8, 0.6), Fraction: 2.2 / 6})
	})
	t.Run("a tangent is the point of contact", func(t *testing.T) {
		hit, ok := Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0)).HitCircle(circle)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 1.0), Normal: Vec(0.0, 1.0), Fraction: 0.5})
	})
	t.Run("an End touching from outside is hit there", func(t *testing.T) {
		hit, ok := Seg(Pt(3.0, 0.0), Pt(1.0, 0.0)).HitCircle(circle)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(1.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 1})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(0.5, 0.0), Pt(5.0, 0.0)).HitCircle(circle)
		assert.False(t, ok)

		_, ok = Seg(Pt(1.0, 0.0), Pt(5.0, 0.0)).HitCircle(circle)
		assert.False(t, ok)

		_, ok = Seg(Pt(1.0, 0.0), Pt(-5.0, 0.0)).HitCircle(circle)
		assert.False(t, ok)
	})
	t.Run("apart or short of the circle hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0)).HitCircle(circle)
		assert.False(t, ok)

		_, ok = Seg(Pt(-3.0, 0.0), Pt(-2.0, 0.0)).HitCircle(circle)
		assert.False(t, ok)
	})
	t.Run("a circle of zero radius faces back along the segment", func(t *testing.T) {
		hit, ok := Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0)).HitCircle(Circ(Pt(0.0, 0.0), 0.0))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(-10, 1), Pt(10, 1)).HitCircle(Circ(Pt(0, 0), 5))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-5, 1), Normal: Vec(-math.Sqrt(24)/5, 0.2), Fraction: 0.25})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		unit, through := Circ(Pt(0, 0), 5), Seg(Pt(-10, 0), Pt(10, 0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitCircle(unit)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range append(grazingSegments(c), segmentFixtures...) {
				hit, ok := s.HitCircle(c)
				assertShapeHit(t, s, hit, ok, c.Contains(s.Start), s.IntersectionCircle(c).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
}

func TestSegment_HitCapsule(t *testing.T) {
	capsule := Cap(Pt(0.0, 0.0), Pt(10.0, 0.0), 2.0)

	t.Run("a straight side, with its perpendicular", func(t *testing.T) {
		hit, ok := Seg(Pt(5.0, 10.0), Pt(5.0, 0.0)).HitCapsule(capsule)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(5.0, 2.0), Normal: Vec(0.0, 1.0), Fraction: 0.8})

		hit, ok = Seg(Pt(5.0, -10.0), Pt(5.0, 10.0)).HitCapsule(capsule)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(5.0, -2.0), Normal: Vec(0.0, -1.0), Fraction: 0.4})
	})
	t.Run("a round end, with the direction from the end of the segment", func(t *testing.T) {
		hit, ok := Seg(Pt(20.0, 0.0), Pt(0.0, 0.0)).HitCapsule(capsule)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(12.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 0.4})

		hit, ok = Seg(Pt(11.2, 10.0), Pt(11.2, 0.0)).HitCapsule(capsule)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(11.2, 1.6), Normal: Vec(0.6, 0.8), Fraction: 0.84})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(5.0, 1.0), Pt(5.0, 9.0)), Seg(Pt(5.0, 2.0), Pt(5.0, 9.0)), Seg(Pt(12.0, 0.0), Pt(0.0, 0.0))} {
			_, ok := s.HitCapsule(capsule)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("short of the capsule or beside it hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(5.0, 10.0), Pt(5.0, 3.0)), Seg(Pt(20.0, 5.0), Pt(-10.0, 5.0)), Seg(Pt(13.0, 3.0), Pt(13.0, -3.0))} {
			_, ok := s.HitCapsule(capsule)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(5, 10), Pt(5, 0)).HitCapsule(Cap(Pt(0, 0), Pt(10, 0), 2))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(5, 2), Normal: Vec(0.0, 1.0), Fraction: 0.8})
	})
	t.Run("a capsule of no radius is hit as its segment and one of no length as a circle", func(t *testing.T) {
		for _, cast := range segmentFixtures {
			for _, wall := range segmentFixtures {
				hit, ok := cast.HitCapsule(Cap(wall.Start, wall.End, 0.0))
				expected, found := cast.HitSegment(wall)

				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%s → %s: ", cast, wall))
			}
			for _, c := range circleFixtures {
				hit, ok := cast.HitCapsule(Cap(c.Center, c.Center, c.Radius))
				expected, found := cast.HitCircle(c)

				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%s → %s: ", cast, c))
			}
		}
	})
	t.Run("the hit lies on the segment and the boundary, facing the segment", func(t *testing.T) {
		for _, c := range capsuleFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", s, c)

					hit, ok := s.HitCapsule(c)
					assert.Equal(t, ok, !c.Contains(start) && c.IntersectsSegment(s), message)
					if !ok {
						assert.Equal(t, hit, Hit[float64]{}, message)

						continue
					}

					assert.True(t, s.Contains(hit.Point) && c.Contains(hit.Point), message+hit.Point.String())
					assert.False(t, c.Contains(s.PointAt(hit.Fraction/2)), message+hit.Point.String())
					assertFacing(t, s, hit, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(5.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitCapsule(capsule)
		}, 0)
	})
}

func TestSegment_HitSector(t *testing.T) {
	sector := Sect(Pt(0.0, 0.0), 5.0, 0, Pi/2)
	notched := Sect(Pt(0.0, 0.0), 5.0, Pi/4, 3*Pi/2)

	t.Run("the arc, with the direction from the center", func(t *testing.T) {
		hit, ok := Seg(Pt(3.0, 10.0), Pt(3.0, 0.0)).HitSector(sector)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(3.0, 4.0), Normal: Vec(0.6, 0.8), Fraction: 0.6})
	})
	t.Run("a radius, with its perpendicular", func(t *testing.T) {
		hit, ok := Seg(Pt(-2.0, 2.0), Pt(2.0, 2.0)).HitSector(sector)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 2.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})

		hit, ok = Seg(Pt(2.0, -2.0), Pt(2.0, 2.0)).HitSector(sector)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.0, 0.0), Normal: Vec(0.0, -1.0), Fraction: 0.5})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(1.0, 1.0), Pt(9.0, 9.0)), Seg(Pt(3.0, 4.0), Pt(9.0, 9.0)), Seg(Pt(0.0, 2.0), Pt(-5.0, 2.0))} {
			_, ok := s.HitSector(sector)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("short of the sector or beside it hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(9.0, 9.0), Pt(4.0, 4.0)), Seg(Pt(-1.0, -1.0), Pt(-1.0, 9.0)), Seg(Pt(3.0, -4.0), Pt(-4.0, 3.0))} {
			_, ok := s.HitSector(sector)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("from the notch of a sector past half a turn it reaches a radius, and passes out of the notch", func(t *testing.T) {
		hit, ok := Seg(Pt(3.0, 0.0), Pt(3.0, 4.0)).HitSector(notched)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(3.0, 3.0), Normal: Vec(math.Sqrt2/2, -math.Sqrt2/2), Fraction: 0.75})

		hit, ok = Seg(Pt(8.0, 0.0), Pt(-2.0, 0.0)).HitSector(notched)
		assert.True(t, ok)
		geomtest.AssertPoint(t, hit.Point, Pt(0.0, 0.0))
		geomtest.AssertNumber(t, hit.Fraction, 0.8)

		_, ok = Seg(Pt(1.0, 0.2), Pt(8.0, 1.0)).HitSector(notched)
		assert.False(t, ok)
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(3, 10), Pt(3, 0)).HitSector(Sect(Pt(0, 0), 5, 0, Pi/2))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(3, 4), Normal: Vec(0.6, 0.8), Fraction: 0.6})
	})
	t.Run("the whole disc is hit as its circle", func(t *testing.T) {
		for _, cast := range segmentFixtures {
			for _, c := range circleFixtures {
				hit, ok := cast.HitSector(Sect(c.Center, c.Radius, 1, 2*Pi))
				expected, found := cast.HitCircle(c)
				message := fmt.Sprintf("%s → %s: ", cast, c)

				assert.Equal(t, ok, found, message)
				geomtest.AssertHit(t, hit, expected, message)
			}
		}
	})
	t.Run("the hit lies on the segment and the boundary, facing the segment", func(t *testing.T) {
		for _, sector := range sectorFixtures() {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", s, sector)

					hit, ok := s.HitSector(sector)
					assert.Equal(t, ok, !sector.Contains(start) && sector.IntersectsSegment(s), message)
					if !ok {
						assert.Equal(t, hit, Hit[float64]{}, message)

						continue
					}

					assert.True(t, s.Contains(hit.Point) && sector.Contains(hit.Point), message+hit.Point.String())
					assert.False(t, sector.Contains(s.PointAt(hit.Fraction/2)), message+hit.Point.String())
					assertFacing(t, s, hit, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitSector(sector)
		}, 0)
	})
}

func TestSegment_HitArc(t *testing.T) {
	arc := Ac(Pt(0.0, 0.0), 5.0, 0, Pi/2)

	t.Run("from outside the circle, with the direction from the center", func(t *testing.T) {
		hit, ok := Seg(Pt(3.0, 10.0), Pt(3.0, 0.0)).HitArc(arc)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(3.0, 4.0), Normal: Vec(0.6, 0.8), Fraction: 0.6})
	})
	t.Run("from inside the circle, with the direction toward the center", func(t *testing.T) {
		hit, ok := Seg(Pt(3.0, 0.0), Pt(3.0, 10.0)).HitArc(arc)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(3.0, 4.0), Normal: Vec(-0.6, -0.8), Fraction: 0.4})
	})
	t.Run("a Start on the arc, a cast past the sweep or short of it hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(3.0, 4.0), Pt(3.0, 10.0)), Seg(Pt(-3.0, 10.0), Pt(-3.0, -10.0)), Seg(Pt(3.0, 10.0), Pt(3.0, 5.0))} {
			_, ok := s.HitArc(arc)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(3, 10), Pt(3, 0)).HitArc(Ac(Pt(0, 0), 5, 0, Pi/2))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(3, 4), Normal: Vec(0.6, 0.8), Fraction: 0.6})
	})
	t.Run("the hit is the first crossing, from a Start off the arc", func(t *testing.T) {
		for _, a := range arcFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", s, a)

					hit, ok := s.HitArc(a)
					crossings := s.IntersectionArc(a).AppendTo(nil)
					assert.Equal(t, ok, !a.Contains(start) && len(crossings) > 0, message)
					if ok {
						geomtest.AssertPoint(t, hit.Point, crossings[0], message)
						assertFacing(t, s, hit, message)
					}
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(3.0, 10.0), Pt(tenth, -10.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitArc(arc)
		}, 0)
	})
}

func TestSegment_HitSegment(t *testing.T) {
	through := Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))

	t.Run("crossing a wall is the crossing, with the perpendicular facing the segment", func(t *testing.T) {
		hit, ok := through.HitSegment(Seg(Pt(0.0, -1.0), Pt(0.0, 1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})

		hit, ok = through.HitSegment(Seg(Pt(0.0, 1.0), Pt(0.0, -1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})

		hit, ok = through.Reverse().HitSegment(Seg(Pt(1.0, -1.0), Pt(1.0, 1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(1.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 0.25})
	})
	t.Run("a wall ending on the segment is hit at that end", func(t *testing.T) {
		hit, ok := through.HitSegment(Seg(Pt(0.0, 0.0), Pt(0.0, 3.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})
	})
	t.Run("an End touching a wall is hit there", func(t *testing.T) {
		hit, ok := through.HitSegment(Seg(Pt(2.0, -1.0), Pt(2.0, 1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 1})
	})
	t.Run("a wall along the segment is met end on at the end reached first", func(t *testing.T) {
		expected := Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5}

		hit, ok := through.HitSegment(Seg(Pt(0.0, 0.0), Pt(1.0, 0.0)))
		assertHit(t, hit, ok, expected)

		hit, ok = through.HitSegment(Seg(Pt(1.0, 0.0), Pt(0.0, 0.0)))
		assertHit(t, hit, ok, expected)

		hit, ok = through.HitSegment(Seg(Pt(0.0, 0.0), Pt(5.0, 0.0)))
		assertHit(t, hit, ok, expected)

		hit, ok = through.HitSegment(Seg(Pt(5.0, 0.0), Pt(0.0, 0.0)))
		assertHit(t, hit, ok, expected)
	})
	t.Run("a wall of zero length faces back along the segment", func(t *testing.T) {
		hit, ok := through.HitSegment(Seg(Pt(1.0, 0.0), Pt(1.0, 0.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(1.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.75})
	})
	t.Run("a Start on the wall hits nothing", func(t *testing.T) {
		_, ok := through.HitSegment(Seg(Pt(-2.0, -1.0), Pt(-2.0, 1.0)))
		assert.False(t, ok)

		_, ok = through.HitSegment(Seg(Pt(-3.0, 0.0), Pt(1.0, 0.0)))
		assert.False(t, ok)
	})
	t.Run("apart or short of the wall hits nothing", func(t *testing.T) {
		hit, ok := through.HitSegment(Seg(Pt(3.0, -1.0), Pt(3.0, 1.0)))
		assert.False(t, ok)
		assert.Equal(t, hit, Hit[float64]{})

		_, ok = through.HitSegment(Seg(Pt(-2.0, 1.0), Pt(2.0, 1.0)))
		assert.False(t, ok)
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(0, 0), Pt(3, 4)).HitSegment(Seg(Pt(0, 2), Pt(4, 2)))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(2, 2), Normal: Vec(0.0, -1.0), Fraction: 0.56})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		wall := Seg(Pt(0.0, -1.0), Pt(0.0, 1.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitSegment(wall)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, wall := range segmentFixtures {
				hit, ok := s.HitSegment(wall)
				assertWallHit(t, s, wall, hit, ok, s.IntersectsSegment(wall), fmt.Sprintf("%s → %s: ", s, wall))
			}
		}
	})
}

func TestSegment_HitRay(t *testing.T) {
	through := Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))

	t.Run("crossing the ray is the crossing, with the perpendicular facing the segment", func(t *testing.T) {
		hit, ok := through.HitRay(Ry(Pt(0.0, -1.0), Vec(0.0, 1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})
	})
	t.Run("a ray running away or a Start on it hits nothing", func(t *testing.T) {
		_, ok := through.HitRay(Ry(Pt(0.0, 1.0), Vec(0.0, 1.0)))
		assert.False(t, ok)

		_, ok = through.HitRay(Ry(Pt(-2.0, -1.0), Vec(0.0, 1.0)))
		assert.False(t, ok)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, wall := range rayFixtures {
				hit, ok := s.HitRay(wall)
				assertWallHit(t, s, wall, hit, ok, s.IntersectsRay(wall), fmt.Sprintf("%s → %s: ", s, wall))
			}
		}
	})
}

func TestSegment_HitLine(t *testing.T) {
	through := Seg(Pt(-2.0, 0.0), Pt(2.0, 0.0))

	t.Run("crossing the line is the crossing, with the perpendicular facing the segment", func(t *testing.T) {
		hit, ok := through.HitLine(Ln(Pt(1.0, 5.0), Vec(0.0, 1.0)))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(1.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.75})
	})
	t.Run("short of the line or a Start on it hits nothing", func(t *testing.T) {
		_, ok := through.HitLine(Ln(Pt(3.0, 5.0), Vec(0.0, 1.0)))
		assert.False(t, ok)

		_, ok = through.HitLine(Ln(Pt(-2.0, 5.0), Vec(0.0, 1.0)))
		assert.False(t, ok)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, s := range segmentFixtures {
			for _, wall := range lineFixtures {
				hit, ok := s.HitLine(wall)
				assertWallHit(t, s, wall, hit, ok, s.IntersectsLine(wall), fmt.Sprintf("%s → %s: ", s, wall))
			}
		}
	})
}

func TestSegment_HitPolyline(t *testing.T) {
	path := PolLn([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 3.0)})

	t.Run("the first edge reached, with its perpendicular", func(t *testing.T) {
		hit, ok := Seg(Pt(2.0, 5.0), Pt(2.0, -5.0)).HitPolyline(path)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.0, 0.0), Normal: Vec(0.0, 1.0), Fraction: 0.5})

		hit, ok = Seg(Pt(6.0, 1.0), Pt(-2.0, 1.0)).HitPolyline(path)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(4.0, 1.0), Normal: Vec(1.0, 0.0), Fraction: 0.25})
	})
	t.Run("a corner touched in passing takes the edge that faces Start", func(t *testing.T) {
		hit, ok := Seg(Pt(5.0, 1.0), Pt(3.0, -1.0)).HitPolyline(path)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(4.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 0.5})

		hit, ok = Seg(Pt(3.0, -1.0), Pt(5.0, 1.0)).HitPolyline(path)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(4.0, 0.0), Normal: Vec(0.0, -1.0), Fraction: 0.5})

		hit, ok = Seg(Pt(5.0, 1.0), Pt(3.0, -1.0)).HitPolyline(path.Reverse())
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(4.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 0.5})
	})
	t.Run("a Start on the path, a cast beside it or an empty path hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(2.0, 0.0), Pt(2.0, 5.0)).HitPolyline(path)
		assert.False(t, ok)

		_, ok = Seg(Pt(0.0, 1.0), Pt(3.0, 1.0)).HitPolyline(path)
		assert.False(t, ok)

		_, ok = Seg(Pt(0.0, 1.0), Pt(3.0, 1.0)).HitPolyline(Polyline[float64]{})
		assert.False(t, ok)
	})
	t.Run("a closed path is hit from outside where the polygon of its points is", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					if p.Contains(start) {
						continue
					}

					hit, ok := s.HitPolyline(p.Polyline())
					expected, found := s.HitPolygon(p)
					message := fmt.Sprintf("%s → %s: ", s, p)

					assert.Equal(t, ok, found, message)
					geomtest.AssertHit(t, hit, expected, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(2.0, 5.0), Pt(tenth, -5.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitPolyline(path)
		}, 0)
	})
}

func TestSegment_HitQuadraticBezier(t *testing.T) {
	curve := QuadBez(Pt(0.0, 0.0), Pt(2.0, 4.0), Pt(4.0, 0.0))

	t.Run("the first crossing, with the perpendicular of the curve there", func(t *testing.T) {
		hit, ok := Seg(Pt(2.0, 5.0), Pt(2.0, -1.0)).HitQuadraticBezier(curve)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.0, 2.0), Normal: Vec(0.0, 1.0), Fraction: 0.5})

		hit, ok = Seg(Pt(2.0, -1.0), Pt(2.0, 5.0)).HitQuadraticBezier(curve)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.0, 2.0), Normal: Vec(0.0, -1.0), Fraction: 0.5})
	})
	t.Run("a Start on the curve or a cast short of it hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(2.0, 2.0), Pt(2.0, 5.0)).HitQuadraticBezier(curve)
		assert.False(t, ok)

		_, ok = Seg(Pt(2.0, 5.0), Pt(2.0, 3.0)).HitQuadraticBezier(curve)
		assert.False(t, ok)
	})
	t.Run("matches the hit of the curve raised to a cubic", func(t *testing.T) {
		for _, q := range quadraticFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitQuadraticBezier(q)
				expected, found := s.HitCubicBezier(q.CubicBezier())

				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%s → %s: ", s, q))
			}
		}
	})
}

func TestSegment_HitCubicBezier(t *testing.T) {
	wave := CubBez(Pt(0.0, 0.0), Pt(2.0, 6.0), Pt(4.0, -6.0), Pt(6.0, 0.0))

	t.Run("the first crossing, with the perpendicular of the curve there", func(t *testing.T) {
		hit, ok := Seg(Pt(8.0, 0.0), Pt(-2.0, 0.0)).HitCubicBezier(wave)
		assert.True(t, ok)
		geomtest.AssertPoint(t, hit.Point, Pt(6.0, 0.0))
		geomtest.AssertVector(t, hit.Normal, Vec(3.0, -1.0).Normalize(), fmt.Sprint(hit))
		geomtest.AssertNumber(t, hit.Fraction, 0.2)
	})
	t.Run("a Start on the curve or a cast beside it hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(3.0, 0.0), Pt(3.0, 5.0)).HitCubicBezier(wave)
		assert.False(t, ok)

		_, ok = Seg(Pt(-1.0, 3.0), Pt(7.0, 3.0)).HitCubicBezier(wave)
		assert.False(t, ok)
	})
	t.Run("a Start on the curve hits nothing where the cast leaves along it", func(t *testing.T) {
		curve := CubBez(Pt(-3.0, -3.5), Pt(1.0, -2.0), Pt(1.5, -2.0), Pt(-3.0, -2.0))

		_, ok := Seg(curve.End, Pt(-0.5, -2.0)).HitCubicBezier(curve)
		assert.False(t, ok)

		for _, c := range cubicFixtures {
			for _, fraction := range curveFractions {
				start := c.PointAt(fraction)

				_, ok := Seg(start, start.Add(c.TangentAt(fraction))).HitCubicBezier(c)
				assert.False(t, ok, fmt.Sprintf("%s @%v: ", c, fraction))
			}
		}
	})
	t.Run("int rounds the point", func(t *testing.T) {
		hit, ok := Seg(Pt(30, 50), Pt(30, -50)).HitCubicBezier(CubBez(Pt(0, 0), Pt(20, 60), Pt(40, -60), Pt(60, 0)))
		assert.True(t, ok)
		geomtest.AssertPoint(t, hit.Point, Pt(30, 0))
	})
	t.Run("the hit is the first crossing, facing the segment", func(t *testing.T) {
		for _, curve := range cubicFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", s, curve)

					hit, ok := s.HitCubicBezier(curve)
					crossings := s.IntersectionCubicBezier(curve)
					if !ok {
						assert.True(t, len(crossings) == 0 || crossings[0].Equal(start), message)

						continue
					}

					geomtest.AssertPoint(t, hit.Point, crossings[0], message)
					assertFacing(t, s, hit, message)
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(8.0, tenth), Pt(-2.0, 0.0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitCubicBezier(wave)
		}, 0)
	})
}

func TestSegment_HitPolygon(t *testing.T) {
	square := Pol(squareVertices())
	notched := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(2, 1), Pt(0, 4)})

	t.Run("from outside is the first crossing, with the normal of its edge", func(t *testing.T) {
		hit, ok := Seg(Pt(-1, 1), Pt(5, 1)).HitPolygon(square)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(0, 1), Normal: Vec(-1.0, 0.0), Fraction: 1.0 / 6})

		hit, ok = Seg(Pt(5, 1), Pt(-1, 1)).HitPolygon(square)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(2, 1), Normal: Vec(1.0, 0.0), Fraction: 0.5})
	})
	t.Run("the normal faces the segment whichever way the polygon is wound", func(t *testing.T) {
		vertices := squareVertices()
		slices.Reverse(vertices)

		hit, ok := Seg(Pt(-1, 1), Pt(5, 1)).HitPolygon(Pol(vertices))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(0, 1), Normal: Vec(-1.0, 0.0), Fraction: 1.0 / 6})
	})
	t.Run("a concave polygon is hit on the first edge reached", func(t *testing.T) {
		hit, ok := Seg(Pt(2, 6), Pt(2, -2)).HitPolygon(notched)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(2, 1), Normal: Vec(1/math.Sqrt(13)*-3, 2/math.Sqrt(13)), Fraction: 0.625})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(1, 1), Pt(5, 1)).HitPolygon(square)
		assert.False(t, ok)

		_, ok = Seg(Pt(0, 1), Pt(-5, 1)).HitPolygon(square)
		assert.False(t, ok)
	})
	t.Run("apart and empty hit nothing", func(t *testing.T) {
		_, ok := Seg(Pt(-1, 5), Pt(5, 5)).HitPolygon(square)
		assert.False(t, ok)

		_, ok = Seg(Pt(-1, 1), Pt(5, 1)).HitPolygon(Polygon[int]{})
		assert.False(t, ok)
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-1, 1), Pt(5, 1))

		assertAllocations(t, func() {
			_, sinkBool = through.HitPolygon(notched)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				hit, ok := s.HitPolygon(p)
				assertShapeHit(t, s, hit, ok, p.Contains(s.Start), s.IntersectionPolygon(p), fmt.Sprintf("%s → %s: ", s, p))
			}
		}
	})
}

func TestSegment_HitRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4), 0)

	t.Run("from outside is the first crossing, with the normal of its edge", func(t *testing.T) {
		hit, ok := Seg(Pt(-5, 0), Pt(5, 0)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3})

		hit, ok = Seg(Pt(1, 6), Pt(1, -6)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(1, 2), Normal: Vec(0.0, 1.0), Fraction: 1.0 / 3})
	})
	t.Run("a corner met head on takes the edge first in Edges", func(t *testing.T) {
		hit, ok := Seg(Pt(-4, -4), Pt(4, 4)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-2, -2), Normal: Vec(0.0, -1.0), Fraction: 0.25})
	})
	t.Run("a corner touched in passing takes the edge facing Start", func(t *testing.T) {
		hit, ok := Seg(Pt(1, 3), Pt(3, 1)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(2, 2), Normal: Vec(0.0, 1.0), Fraction: 0.5})

		hit, ok = Seg(Pt(3, 1), Pt(1, 3)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(2, 2), Normal: Vec(1.0, 0.0), Fraction: 0.5})
	})
	t.Run("along an edge is the corner it reaches, on the edge it meets there", func(t *testing.T) {
		hit, ok := Seg(Pt(-5, -2), Pt(5, -2)).HitRectangle(rectangle)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-2, -2), Normal: Vec(-1.0, 0.0), Fraction: 0.3})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(0, 0), Pt(5, 0)).HitRectangle(rectangle)
		assert.False(t, ok)

		_, ok = Seg(Pt(2, 0), Pt(5, 0)).HitRectangle(rectangle)
		assert.False(t, ok)
	})
	t.Run("apart hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(3, -5), Pt(3, 5)).HitRectangle(rectangle)
		assert.False(t, ok)
	})
	t.Run("a rotated rectangle is hit on its turned edges", func(t *testing.T) {
		diamond := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0), 0).Rotate(Pi / 4)

		hit, ok := Seg(Pt(-3.0, 0.5), Pt(3.0, 0.5)).HitRectangle(diamond)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.5-Sqrt2, 0.5), Normal: Vec(-Sqrt2/2, Sqrt2/2), Fraction: (3.5 - Sqrt2) / 6})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-5, 1), Pt(5, 1))

		assertAllocations(t, func() {
			_, sinkBool = through.HitRectangle(rectangle)
		}, 0)
		assertAllocations(t, func() {
			_, sinkBool = through.HitRectangle(rectangle.Rotate(Pi / 5))
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitRectangle(r)
				assertShapeHit(t, s, hit, ok, r.Contains(s.Start), s.IntersectionRectangle(r).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, r))
			}
		}
	})
}

func TestSegment_HitRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("from outside is the first crossing, with the normal of its edge", func(t *testing.T) {
		hit, ok := Seg(Pt(-3, 1), Pt(3, 1)).HitRegularPolygon(diamond)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-1, 1), Normal: Vec(-Sqrt2/2, Sqrt2/2), Fraction: 1.0 / 3})
	})
	t.Run("a Start inside hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(0, 0), Pt(5, 0)).HitRegularPolygon(diamond)
		assert.False(t, ok)
	})
	t.Run("apart and empty hit nothing", func(t *testing.T) {
		_, ok := Seg(Pt(-3, 3), Pt(3, 3)).HitRegularPolygon(diamond)
		assert.False(t, ok)

		_, ok = Seg(Pt(-3, 1), Pt(3, 1)).HitRegularPolygon(RegularPolygon[int]{})
		assert.False(t, ok)
	})
	t.Run("matches HitPolygon on the Polygon form", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitRegularPolygon(rp)
				expected, expectedOK := s.HitPolygon(rp.Polygon())

				assert.Equal(t, ok, expectedOK, fmt.Sprintf("%s → %s: ", s, rp))
				assert.Equal(t, hit, expected, fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-3, 1), Pt(3, 1))

		assertAllocations(t, func() {
			_, sinkBool = through.HitRegularPolygon(diamond)
		}, 0)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitRegularPolygon(rp)
				assertShapeHit(t, s, hit, ok, rp.Contains(s.Start), s.IntersectionRegularPolygon(rp).AppendTo(nil), fmt.Sprintf("%s → %s: ", s, rp))
			}
		}
	})
}

func TestSegment_HitBox(t *testing.T) {
	box := BoxFromMinMax(Pt(-2, -2), Pt(2, 2))

	t.Run("from outside is the first crossing, with the normal of its edge", func(t *testing.T) {
		hit, ok := Seg(Pt(-5, 0), Pt(5, 0)).HitBox(box)
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3})
	})
	t.Run("a Start inside hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(0, 0), Pt(5, 0)).HitBox(box)
		assert.False(t, ok)
	})
	t.Run("apart hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(3, -5), Pt(3, 5)).HitBox(box)
		assert.False(t, ok)
	})
	t.Run("matches HitRectangle on the box's Rectangle", func(t *testing.T) {
		for _, b := range boxFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitBox(b)
				expected, expectedOK := s.HitRectangle(b.Rectangle())

				assert.Equal(t, ok, expectedOK, fmt.Sprintf("%s → %s: ", s, b))
				geomtest.AssertHit(t, hit, expected, fmt.Sprintf("%s → %s: ", s, b))
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through := Seg(Pt(-5, 0), Pt(5, 0))

		assertAllocations(t, func() {
			_, sinkBool = through.HitBox(box)
		}, 0)
	})
}

func TestSegment_HitEllipse(t *testing.T) {
	ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

	t.Run("from outside is the first crossing, with the perpendicular of the boundary", func(t *testing.T) {
		hit, ok := Seg(Pt(-8.0, 0.0), Pt(0.0, 0.0)).HitEllipse(ellipse)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(-4.0, 0.0), Normal: Vec(-1.0, 0.0), Fraction: 0.5})

		hit, ok = Seg(Pt(0.0, -6.0), Pt(0.0, 2.0)).HitEllipse(ellipse)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, -2.0), Normal: Vec(0.0, -1.0), Fraction: 0.5})
	})
	t.Run("off the axes the normal is not the direction from the center", func(t *testing.T) {
		hit, ok := Seg(Pt(8.0, 1.6), Pt(-8.0, 1.6)).HitEllipse(ellipse)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(2.4, 1.6), Normal: Vec(2.4/16, 1.6/4).Normalize(), Fraction: 0.35})
	})
	t.Run("a turned ellipse turns the point and the normal with it", func(t *testing.T) {
		hit, ok := Seg(Pt(-5.0, 2.0), Pt(1.0, 2.0)).HitEllipse(Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/2))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(-1.0, 2.0), Normal: Vec(-1.0, 0.0), Fraction: 4.0 / 6})
	})
	t.Run("a tangent on exact numbers is the point of contact", func(t *testing.T) {
		hit, ok := Seg(Pt(-5.0, 2.0), Pt(5.0, 2.0)).HitEllipse(ellipse)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 2.0), Normal: Vec(0.0, 1.0), Fraction: 0.5})
	})
	t.Run("an End touching from outside is hit there", func(t *testing.T) {
		hit, ok := Seg(Pt(8.0, 0.0), Pt(4.0, 0.0)).HitEllipse(ellipse)
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(4.0, 0.0), Normal: Vec(1.0, 0.0), Fraction: 1})
	})
	t.Run("a Start inside or on the boundary hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(1.0, 0.0), Pt(9.0, 0.0)), Seg(Pt(4.0, 0.0), Pt(9.0, 0.0)), Seg(Pt(4.0, 0.0), Pt(-9.0, 0.0))} {
			_, ok := s.HitEllipse(ellipse)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("apart, short of the ellipse or of no length hits nothing", func(t *testing.T) {
		for _, s := range []Segment[float64]{Seg(Pt(-5.0, 3.0), Pt(5.0, 3.0)), Seg(Pt(-8.0, 0.0), Pt(-5.0, 0.0)), Seg(Pt(-8.0, 0.0), Pt(-12.0, 0.0)), Seg(Pt(-8.0, 0.0), Pt(-8.0, 0.0))} {
			_, ok := s.HitEllipse(ellipse)
			assert.False(t, ok, s.String())
		}
	})
	t.Run("a degenerate ellipse is hit as the wall it is", func(t *testing.T) {
		hit, ok := Seg(Pt(0.0, -3.0), Pt(0.0, 3.0)).HitEllipse(Ell(Pt(0.0, 0.0), Sz(4.0, 0.0), 0))
		assertHit(t, hit, ok, Hit[float64]{Point: Pt(0.0, 0.0), Normal: Vec(0.0, -1.0), Fraction: 0.5})

		_, ok = Seg(Pt(5.0, -3.0), Pt(5.0, 3.0)).HitEllipse(Ell(Pt(0.0, 0.0), Sz(4.0, 0.0), 0))
		assert.False(t, ok)
	})
	t.Run("a NaN coordinate hits nothing", func(t *testing.T) {
		_, ok := Seg(Pt(math.NaN(), 0.0), Pt(0.0, 0.0)).HitEllipse(ellipse)
		assert.False(t, ok)

		_, ok = Seg(Pt(-8.0, 0.0), Pt(math.NaN(), 0.0)).HitEllipse(ellipse)
		assert.False(t, ok)
	})
	t.Run("int rounds the point and keeps the normal", func(t *testing.T) {
		hit, ok := Seg(Pt(-10, 0), Pt(10, 0)).HitEllipse(Ell(Pt(0, 0), Sz(4, 2), 0))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-4, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3})

		hit, ok = Seg(Pt(-10, 1), Pt(10, 1)).HitEllipse(Ell(Pt(0, 0), Sz(4, 2), 0))
		assertHit(t, hit, ok, Hit[int]{Point: Pt(-3, 1), Normal: Vec(-math.Sqrt(12)/16, 0.25).Normalize(), Fraction: 0.35})
	})
	t.Run("a circle is hit as the circle of the same radius", func(t *testing.T) {
		for _, c := range circleFixtures {
			for _, s := range segmentFixtures {
				hit, ok := s.HitEllipse(c.Ellipse())
				expected, found := s.HitCircle(c)

				assert.Equal(t, ok, found, fmt.Sprintf("%s → %s: ", s, c))
				geomtest.AssertHit(t, hit, expected, fmt.Sprintf("%s → %s: ", s, c))
			}
		}
	})
	t.Run("the hit lies on the segment and the boundary, facing the segment", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, start := range pointFixtures {
				for _, end := range pointFixtures {
					s := Seg(start, end)
					message := fmt.Sprintf("%s → %s: ", s, e)

					hit, ok := s.HitEllipse(e)
					assert.Equal(t, ok, !e.Contains(start) && Intersects[float64](s, e), message)
					if !ok {
						assert.Equal(t, hit, Hit[float64]{}, message)

						continue
					}

					assert.True(t, s.Contains(hit.Point) && e.Contains(hit.Point), message+hit.Point.String())
					assert.False(t, e.Contains(s.PointAt(hit.Fraction/2)), message+hit.Point.String())
					assertFacing(t, s, hit, message)
				}
			}
		}
	})
	t.Run("far from the origin the float32 hit lies on the boundary", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, e := range ellipseFixtures {
				for _, s := range segmentFixtures {
					ellipse, segment := e.Cast[float32]().Translate(offset), s.Cast[float32]().Translate(offset)

					if hit, ok := segment.HitEllipse(ellipse); ok {
						assert.True(t, ellipse.Contains(hit.Point), fmt.Sprintf("%s → %s: ", segment, ellipse))
					}
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		through, turned := Seg(Pt(-10.0, tenth), Pt(10.0, 0.0)), Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0.3)

		assertAllocations(t, func() {
			_, sinkBool = through.HitEllipse(turned)
		}, 0)
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

func TestSegment_IsPoint(t *testing.T) {
	t.Run("equal ends", func(t *testing.T) {
		assert.True(t, Segment[int]{}.IsPoint())
		assert.True(t, Seg(Pt(3, -2), Pt(3, -2)).IsPoint())
	})
	t.Run("different ends", func(t *testing.T) {
		assert.False(t, Seg(Pt(3, -2), Pt(3, -1)).IsPoint())
		assert.False(t, Seg(Pt(0.0, 0.0), Pt(1.5, 0.0)).IsPoint())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Seg(Pt(1.0, 2.0), Pt(1.0+Delta/2, 2.0)).IsPoint())
		assert.False(t, Seg(Pt(1.0, 2.0), Pt(1.0+2*Delta, 2.0)).IsPoint())
	})
	t.Run("the part a tangent segment is clipped to", func(t *testing.T) {
		part, ok := Seg(Pt(-5.0, 1.0), Pt(5.0, 1.0)).ClipCircle(Circ(Pt(0.0, 0.0), 1.0))

		assert.True(t, ok)
		assert.True(t, part.IsPoint())
	})
}

func TestSegment_Line(t *testing.T) {
	t.Run("from the start along the vector", func(t *testing.T) {
		geomtest.AssertLine(t, Seg(Pt(1, -1), Pt(4, 3)).Line(), Ln(Pt(1, -1), Vec(3, 4)))
	})
	t.Run("a zero-length segment is a line with no direction", func(t *testing.T) {
		geomtest.AssertLine(t, Seg(Pt(2, 3), Pt(2, 3)).Line(), Line[int]{Origin: Pt(2, 3)})
	})
	t.Run("the line contains both ends and the points beyond them", func(t *testing.T) {
		for _, s := range segmentFixtures {
			line := s.Line()

			for _, t0 := range []float64{-2, 0, 0.5, 1, 3} {
				assert.True(t, line.Contains(s.PointAt(t0)), fmt.Sprintf("%s at %v: ", s, t0))
			}
		}
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
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Seg(Pt(100, -34.0000115), Pt(0.2, 0.4)).String)
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
	s.Lerp(Seg(Pt(5, 5), Pt(7, 9)), 0.5)
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

// assertHit checks that a Hit method found the hit expected of it.
func assertHit[T Number](t *testing.T, hit Hit[T], ok bool, expected Hit[T]) {
	t.Helper()

	assert.True(t, ok)
	geomtest.AssertHit(t, hit, expected)
}

// ellipseCasts returns the segments the ellipse fixtures are crossed by: every one between two
// point fixtures, and the segment fixtures.
func ellipseCasts() []Segment[float64] {
	casts := slices.Clone(segmentFixtures)
	for _, start := range pointFixtures {
		for _, end := range pointFixtures {
			casts = append(casts, Seg(start, end))
		}
	}

	return casts
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

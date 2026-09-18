package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestLine_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, -1), Pt(2, 0)), Line[int]{Start: Pt(1, -1), End: Pt(2, 0)})
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.5, -1.25), Pt(2.5, 3.75)), Line[float64]{Start: Pt(0.5, -1.25), End: Pt(2.5, 3.75)})
	})
}

func TestLine_Vector(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVector(t, Ln(Pt(1, 2), Pt(3, 5)).Vector(), Vec(2, 3))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vector(), Vec(0.6, 3.65))
	})
}

func TestLine_Length(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(1, 2), Pt(3, 5)).Length(), math.Sqrt(13))
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Length(), math.Sqrt(13.6825))
	})
}

func TestLine_Midpoint(t *testing.T) {
	t.Run("int rounds the half away from zero", func(t *testing.T) {
		AssertPoint(t, Ln(Pt(1, 2), Pt(3, 5)).Midpoint(), Pt(2, 4))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Midpoint(), Pt(0.9, 1.575))
	})
}

func TestLine_Vertices(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(1, 2), Pt(3, 5)).Vertices(), []Point[int]{{1, 2}, {3, 5}})
	})
	t.Run("float", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vertices(), []Point[float64]{{0.6, -0.25}, {1.2, 3.4}})
	})
}

func TestLine_MinMax(t *testing.T) {
	t.Run("orders the corners", func(t *testing.T) {
		a, b := Ln(Pt(4, 1), Pt(0, 3)).MinMax()

		AssertPoint(t, a, Pt(0, 1))
		AssertPoint(t, b, Pt(4, 3))
	})
	t.Run("matches the corners of Bounds", func(t *testing.T) {
		for _, l := range lineFixtures {
			a, b := l.MinMax()
			c, d := l.Bounds().MinMax()

			AssertPoint(t, a, c, l.String())
			AssertPoint(t, b, d, l.String())
		}
	})
}

func TestLine_Bounds(t *testing.T) {
	t.Run("spans the endpoints", func(t *testing.T) {
		l := Ln(Pt(1, 2), Pt(3, 5))

		AssertRectangle(t, l.Bounds(), Rect(Pt(2, 3), Sz(2, 3)))
		AssertPoint(t, l.Bounds().Min(), l.Start)
		AssertPoint(t, l.Bounds().Max(), l.End)
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Bounds(), Rect(Pt(0.9, 1.575), Sz(0.6, 3.65)))
	})
}

func TestLine_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Translate(Vec(3, -2)), Ln(Pt(4, 0), Pt(6, 3)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Translate(Vec(100.1, -0.1)), Ln(Pt(100.7, -0.35), Pt(101.3, 3.3)))
	})
}

func TestLine_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).MoveTo(Pt(3, -2)), Ln(Pt(3, -2), Pt(5, 1)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).MoveTo(Pt(100.1, -0.1)), Ln(Pt(100.1, -0.1), Pt(100.7, 3.55)))
	})
}

func TestLine_Scale(t *testing.T) {
	t.Run("uniform factor about the midpoint", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0, 0), Pt(4, 2)).Scale(2), Ln(Pt(-2, -1), Pt(6, 3)))
		AssertLine(t, Ln(Pt(0.0, 0.0), Pt(4.0, 2.0)).Scale(0.5), Ln(Pt(1.0, 0.5), Pt(3.0, 1.5)))
	})
	t.Run("per-axis factor changes the direction", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.0, 0.0), Pt(4.0, 2.0)).ScaleXY(1, 3), Ln(Pt(0.0, -2.0), Pt(4.0, 4.0)))
	})
	t.Run("zero collapses onto the midpoint", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.0, 0.0), Pt(4.0, 2.0)).Scale(0), Ln(Pt(2.0, 1.0), Pt(2.0, 1.0)))
	})
	t.Run("int rounds the midpoint, so an odd span drifts by one", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0, 0), Pt(3, 0)).Scale(1), Ln(Pt(0, 0), Pt(3, 0)))
		AssertLine(t, Ln(Pt(0, 0), Pt(3, 0)).Scale(2), Ln(Pt(-2, 0), Pt(4, 0)))
	})
	t.Run("keeps the midpoint and scales the length", func(t *testing.T) {
		for _, l := range lineFixtures {
			scaled := l.Scale(2.5)

			AssertPoint(t, scaled.Midpoint(), l.Midpoint(), fmt.Sprintf("%s: ", l))
			AssertNumber(t, scaled.Length(), l.Length()*2.5, fmt.Sprintf("%s: ", l))
		}
	})
}

func TestLine_Reverse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Reverse(), Ln(Pt(3, 5), Pt(1, 2)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Reverse(), Ln(Pt(1.2, 3.4), Pt(0.6, -0.25)))
	})
}

func TestLine_Lerp(t *testing.T) {
	line := Ln(Pt(0.0, 0.0), Pt(4.0, 2.0))

	t.Run("endpoints at 0 and 1", func(t *testing.T) {
		AssertPoint(t, line.Lerp(0), line.Start)
		AssertPoint(t, line.Lerp(1), line.End)
	})
	t.Run("along the segment", func(t *testing.T) {
		AssertPoint(t, line.Lerp(0.25), Pt(1.0, 0.5))
		AssertPoint(t, line.Lerp(0.5), line.Midpoint())
	})
	t.Run("extrapolates beyond the segment", func(t *testing.T) {
		AssertPoint(t, line.Lerp(-0.5), Pt(-2.0, -1.0))
		AssertPoint(t, line.Lerp(1.5), Pt(6.0, 3.0))
	})
	t.Run("int rounds the half away from zero", func(t *testing.T) {
		AssertPoint(t, Ln(Pt(0, 0), Pt(3, 3)).Lerp(0.5), Pt(2, 2))
		AssertPoint(t, Ln(Pt(0, 0), Pt(-3, -3)).Lerp(0.5), Pt(-2, -2))
	})
	t.Run("lands on the segment", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, fraction := range []float64{0, 0.25, 0.5, 0.75, 1} {
				assert.True(t, l.Contains(l.Lerp(fraction)), fmt.Sprintf("%s at %v: ", l, fraction))
			}
		}
	})
}

func TestLine_Transform(t *testing.T) {
	t.Run("applies the matrix to both points", func(t *testing.T) {
		matrix := Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)

		AssertLine(t, Ln(Pt(1, 2), Pt(3, 4)).Transform(matrix), Ln(Pt(1, 2).Transform(matrix), Pt(3, 4).Transform(matrix)))
	})
	t.Run("float32 matrix", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1.0, 2.0), Pt(0.0, 0.0)).Transform(Mat[float32](1, 0, 1, 0, 1, 1)), Ln(Pt(2.0, 3.0), Pt(1.0, 1.0)))
	})
}

func TestLine_Rotate(t *testing.T) {
	t.Run("quarter turn about the midpoint", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0, 0), Pt(4, 0)).Rotate(Pi/2), Ln(Pt(2, -2), Pt(2, 2)))
	})
	t.Run("half turn reverses the line", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 6)).Rotate(Pi), Ln(Pt(3, 6), Pt(1, 2)))
	})
	t.Run("int rounds the midpoint, so an odd span drifts by one", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Rotate(Pi), Ln(Pt(3, 6), Pt(1, 3)))
	})
	t.Run("float keeps the midpoint and length", func(t *testing.T) {
		l := Ln(Pt(0.6, -0.25), Pt(1.2, 3.4))
		rotated := l.Rotate(0.7)

		AssertPoint(t, rotated.Midpoint(), l.Midpoint())
		AssertNumber(t, rotated.Length(), l.Length())
		AssertNumber(t, rotated.Vector().AngleBetween(l.Vector()), 0.7)
	})
	t.Run("a full turn is identity", func(t *testing.T) {
		for _, l := range lineFixtures {
			AssertLine(t, l.Rotate(2*Pi), l, fmt.Sprintf("%s: ", l))
		}
	})
}

func TestLine_Contains(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		l := Ln(Pt(0, 0), Pt(6, 3))

		assert.True(t, l.Contains(Pt(2, 1)))
		assert.True(t, l.Contains(l.Start))
		assert.True(t, l.Contains(l.End))
		assert.False(t, l.Contains(Pt(2, 2)))
		assert.False(t, l.Contains(Pt(8, 4)))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Pt(1.0, 1.0))

		assert.True(t, l.Contains(Pt(0.5, 0.5+Delta/2)))
		assert.False(t, l.Contains(Pt(0.5, 0.5+2*Delta)))
	})
	t.Run("holds exactly where DistanceTo is zero", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, l.Contains(p), l.DistanceTo(p) == 0, fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_DistanceTo(t *testing.T) {
	l := Ln(Pt(0, 0), Pt(4, 0))

	t.Run("perpendicular to the segment", func(t *testing.T) {
		AssertNumber(t, l.DistanceTo(Pt(2, 3)), 3.0)
		AssertNumber(t, l.DistanceTo(Pt(1, -2)), 2.0)
	})
	t.Run("beyond the start measures to the start", func(t *testing.T) {
		AssertNumber(t, l.DistanceTo(Pt(-3, 4)), 5.0)
	})
	t.Run("beyond the end measures to the end", func(t *testing.T) {
		AssertNumber(t, l.DistanceTo(Pt(7, 4)), 5.0)
	})
	t.Run("on the segment is exactly zero", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(0, 0), Pt(3, 3)).DistanceTo(Pt(1, 1)), 0.0)
		assert.Equal(t, l.DistanceTo(Pt(4, 0)), 0.0)
	})
	t.Run("degenerate segment measures to the point", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(1, 1), Pt(1, 1)).DistanceTo(Pt(4, 5)), 5.0)
	})
	t.Run("float within the tolerance is zero, beyond it is measured", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.Equal(t, l.DistanceTo(Pt(0.5, Delta/2)), 0.0)
		AssertNumber(t, l.DistanceTo(Pt(0.5, 2*Delta)), 2*Delta)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0.0, 0.0), Pt(1.0, 1.0)).DistanceTo(Pt(1.0, 0.0)), OneOverSqrt2)
	})
}

func TestLine_DistanceSquaredTo(t *testing.T) {
	l := Ln(Pt(0, 0), Pt(4, 0))

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		AssertNumber(t, l.DistanceSquaredTo(Pt(2, 3)), 9.0)
		AssertNumber(t, l.DistanceSquaredTo(Pt(-3, 4)), 25.0)
		AssertNumber(t, l.DistanceSquaredTo(Pt(7, 4)), 25.0)
	})
	t.Run("stays fractional for an integer T", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0, 0), Pt(2, 1)).DistanceSquaredTo(Pt(0, 1)), 0.8)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range pointFixtures {
				AssertNumber(t, l.DistanceSquaredTo(p), l.DistanceTo(p)*l.DistanceTo(p), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_DistanceToLine(t *testing.T) {
	diagonal := Ln(Pt(0, 0), Pt(4, 4))

	t.Run("crossing segments are at zero", func(t *testing.T) {
		assert.Equal(t, diagonal.DistanceToLine(Ln(Pt(0, 4), Pt(4, 0))), 0.0)
	})
	t.Run("touching segments are at zero", func(t *testing.T) {
		assert.Equal(t, diagonal.DistanceToLine(Ln(Pt(4, 4), Pt(8, 0))), 0.0)
	})
	t.Run("parallel segments measure the gap", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0, 0), Pt(4, 0)).DistanceToLine(Ln(Pt(1, 3), Pt(3, 3))), 3.0)
	})
	t.Run("apart measures between the nearest endpoints", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0, 0), Pt(4, 0)).DistanceToLine(Ln(Pt(7, 4), Pt(9, 4))), 5.0)
	})
	t.Run("an endpoint nearest an interior point", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0, 0), Pt(4, 0)).DistanceToLine(Ln(Pt(2, 2), Pt(2, 5))), 2.0)
	})
	t.Run("symmetric and zero exactly where Intersects holds", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				AssertNumber(t, a.DistanceToLine(b), b.DistanceToLine(a), fmt.Sprintf("%s → %s: ", a, b))
				assert.Equal(t, a.DistanceToLine(b) <= Delta, a.Intersects(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestLine_DistanceSquaredToLine(t *testing.T) {
	t.Run("is the square of DistanceToLine", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0, 0), Pt(4, 0)).DistanceSquaredToLine(Ln(Pt(7, 4), Pt(9, 4))), 25.0)
		AssertNumber(t, Ln(Pt(0, 0), Pt(4, 0)).DistanceSquaredToLine(Ln(Pt(1, 3), Pt(3, 3))), 9.0)
		assert.Equal(t, Ln(Pt(0, 0), Pt(4, 4)).DistanceSquaredToLine(Ln(Pt(0, 4), Pt(4, 0))), 0.0)
	})
	t.Run("agrees with DistanceToLine", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				AssertNumber(t, a.DistanceSquaredToLine(b), a.DistanceToLine(b)*a.DistanceToLine(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestLine_Intersects(t *testing.T) {
	diagonal := Ln(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		assert.True(t, diagonal.Intersects(Ln(Pt(0, 4), Pt(4, 0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, diagonal.Intersects(Ln(Pt(5, 0), Pt(5, 4))))
		assert.False(t, diagonal.Intersects(Ln(Pt(0, 1), Pt(3, 4))))
	})
	t.Run("touching at an endpoint counts", func(t *testing.T) {
		assert.True(t, diagonal.Intersects(Ln(Pt(4, 4), Pt(8, 0))))
		assert.True(t, diagonal.Intersects(Ln(Pt(2, 2), Pt(2, 8))))
	})
	t.Run("collinear overlap counts and a collinear gap does not", func(t *testing.T) {
		assert.True(t, diagonal.Intersects(Ln(Pt(2, 2), Pt(6, 6))))
		assert.False(t, diagonal.Intersects(Ln(Pt(5, 5), Pt(6, 6))))
	})
	t.Run("parallel segments do not cross", func(t *testing.T) {
		assert.False(t, diagonal.Intersects(Ln(Pt(0, 1), Pt(4, 5))))
	})
	t.Run("a degenerate segment is a point", func(t *testing.T) {
		assert.True(t, diagonal.Intersects(Ln(Pt(1, 1), Pt(1, 1))))
		assert.False(t, diagonal.Intersects(Ln(Pt(1, 2), Pt(1, 2))))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Pt(1.0, 0.0))

		assert.True(t, l.Intersects(Ln(Pt(0.5, Delta/2), Pt(0.5, 1.0))))
		assert.False(t, l.Intersects(Ln(Pt(0.5, 2*Delta), Pt(0.5, 1.0))))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				assert.Equal(t, a.Intersects(b), b.Intersects(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func BenchmarkLine_Intersects(b *testing.B) {
	line := Ln(Pt(0.0, 0.0), Pt(10.0, 10.0))
	crossing, apart := Ln(Pt(0.0, 10.0), Pt(10.0, 0.0)), Ln(Pt(20.0, 0.0), Pt(20.0, 10.0))

	b.Run("crossing", func(b *testing.B) {
		for b.Loop() {
			sinkBool = line.Intersects(crossing)
		}
	})
	b.Run("apart", func(b *testing.B) {
		for b.Loop() {
			sinkBool = line.Intersects(apart)
		}
	})
}

func TestLine_Intersection(t *testing.T) {
	diagonal := Ln(Pt(0, 0), Pt(4, 4))

	t.Run("crossing", func(t *testing.T) {
		point, ok := diagonal.Intersection(Ln(Pt(0, 4), Pt(4, 0)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("apart", func(t *testing.T) {
		_, ok := diagonal.Intersection(Ln(Pt(5, 0), Pt(5, 4)))

		assert.False(t, ok)
	})
	t.Run("touching at an endpoint counts", func(t *testing.T) {
		point, ok := diagonal.Intersection(Ln(Pt(4, 4), Pt(8, 0)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(4, 4))

		point, ok = Ln(Pt(2, 2), Pt(2, 8)).Intersection(diagonal)

		assert.True(t, ok)
		AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("every endpoint can be the touching one", func(t *testing.T) {
		for _, line := range []Line[int]{Ln(Pt(4, 4), Pt(8, 0)), Ln(Pt(8, 0), Pt(4, 4))} {
			point, ok := diagonal.Intersection(line)

			assert.True(t, ok, line.String())
			AssertPoint(t, point, Pt(4, 4), line.String())
		}

		point, ok := diagonal.Intersection(Ln(Pt(2, 6), Pt(6, 2)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(4, 4))

		point, ok = Ln(Pt(4, 4), Pt(0, 0)).Intersection(Ln(Pt(2, 6), Pt(6, 2)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(4, 4))
	})
	t.Run("parallel and collinear segments have no single point", func(t *testing.T) {
		_, ok := diagonal.Intersection(Ln(Pt(0, 1), Pt(4, 5)))
		assert.False(t, ok)

		_, ok = diagonal.Intersection(Ln(Pt(2, 2), Pt(6, 6)))
		assert.False(t, ok)
	})
	t.Run("the lines cross beyond a segment", func(t *testing.T) {
		_, ok := diagonal.Intersection(Ln(Pt(5, 0), Pt(6, 4)))

		assert.False(t, ok)
	})
	t.Run("a degenerate segment is the point where it lies on the other", func(t *testing.T) {
		point, ok := diagonal.Intersection(Ln(Pt(1, 1), Pt(1, 1)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(1, 1))

		point, ok = Ln(Pt(1, 1), Pt(1, 1)).Intersection(diagonal)

		assert.True(t, ok)
		AssertPoint(t, point, Pt(1, 1))

		_, ok = diagonal.Intersection(Ln(Pt(1, 2), Pt(1, 2)))
		assert.False(t, ok)

		_, ok = Ln(Pt(1, 1), Pt(1, 1)).Intersection(Ln(Pt(1, 2), Pt(1, 2)))
		assert.False(t, ok)
	})
	t.Run("a shallow touch is decided on the endpoint distance, like Intersects", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Pt(100.0, 0.0))
		shallow := Ln(Pt(50.0, Delta/2), Pt(150.0, 1e-3))

		point, ok := l.Intersection(shallow)

		assert.True(t, l.Intersects(shallow))
		assert.True(t, ok)
		AssertPoint(t, point, shallow.Start)
	})
	t.Run("int rounds the crossing", func(t *testing.T) {
		point, ok := Ln(Pt(0, 0), Pt(3, 3)).Intersection(Ln(Pt(0, 3), Pt(3, 0)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(2, 2))
	})
	t.Run("float keeps the crossing", func(t *testing.T) {
		point, ok := Ln(Pt(0.0, 0.0), Pt(3.0, 3.0)).Intersection(Ln(Pt(0.0, 3.0), Pt(3.0, 0.0)))

		assert.True(t, ok)
		AssertPoint(t, point, Pt(1.5, 1.5))
	})
	t.Run("float is tolerant", func(t *testing.T) {
		l := Ln(Pt(0.0, 0.0), Pt(1.0, 0.0))

		_, ok := l.Intersection(Ln(Pt(0.5, Delta/2), Pt(0.5, 1.0)))
		assert.True(t, ok)

		_, ok = l.Intersection(Ln(Pt(0.5, 2*Delta), Pt(0.5, 1.0)))
		assert.False(t, ok)
	})
	t.Run("int decides exactly where the crossing is not a lattice point", func(t *testing.T) {
		a, b := Ln(Pt(-6, -6), Pt(-5, -5)), Ln(Pt(-6, -5), Pt(-4, -6))

		point, ok := a.Intersection(b)
		assert.True(t, ok)
		AssertPoint(t, point, Pt(-5, -5))

		point, ok = b.Intersection(a)
		assert.True(t, ok)
		AssertPoint(t, point, Pt(-5, -5))
	})
	t.Run("agrees with Intersects on non-parallel fixtures", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				point, ok := a.Intersection(b)
				parallel := !a.Vector().IsZero() && !b.Vector().IsZero() && a.Vector().Cross(b.Vector()) == 0

				assert.Equal(t, ok, a.Intersects(b) && !parallel, fmt.Sprintf("%s → %s: ", a, b))
				if !ok {
					continue
				}

				assert.True(t, a.Contains(point), fmt.Sprintf("%s → %s on a: ", a, b))
				assert.True(t, b.Contains(point), fmt.Sprintf("%s → %s on b: ", a, b))
			}
		}
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range lineFixtures {
			for _, b := range lineFixtures {
				_, ok := a.Intersection(b)
				_, reverse := b.Intersection(a)

				assert.Equal(t, ok, reverse, fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func FuzzLine_Intersection(f *testing.F) {
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

		a, b := Ln(Pt(x1, y1), Pt(x2, y2)), Ln(Pt(x3, y3), Pt(x4, y4))
		parallel := !a.Vector().IsZero() && !b.Vector().IsZero() && a.Vector().Cross(b.Vector()) == 0

		point, ok := a.Intersection(b)
		_, reverse := b.Intersection(a)

		assert.Equal(t, ok, a.Intersects(b) && !parallel, fmt.Sprintf("%s → %s: ", a, b))
		assert.Equal(t, ok, reverse, fmt.Sprintf("%s → %s: symmetric: ", a, b))
		if !ok {
			return
		}

		assert.True(t, a.Contains(point), fmt.Sprintf("%s → %s on a: ", a, b))
		assert.True(t, b.Contains(point), fmt.Sprintf("%s → %s on b: ", a, b))
	})
}

func TestLine_IntersectionCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0)).IntersectionCircle(circle), []Point[float64]{Pt(-1.0, 0.0), Pt(1.0, 0.0)})
		AssertVertices(t, Ln(Pt(2.0, 0.0), Pt(-2.0, 0.0)).IntersectionCircle(circle), []Point[float64]{Pt(1.0, 0.0), Pt(-1.0, 0.0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-2.0, 0.0), Pt(0.0, 0.0)).IntersectionCircle(circle), []Point[float64]{Pt(-1.0, 0.0)})
		AssertVertices(t, Ln(Pt(0.5, 0.0), Pt(5.0, 0.0)).IntersectionCircle(circle), []Point[float64]{Pt(1.0, 0.0)})
	})
	t.Run("tangent gives one point", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-2.0, 1.0), Pt(2.0, 1.0)).IntersectionCircle(circle), []Point[float64]{Pt(0.0, 1.0)})
		AssertVertices(t, Ln(Pt(-2.0, 1.0+Delta/2), Pt(2.0, 1.0+Delta/2)).IntersectionCircle(circle), []Point[float64]{Pt(0.0, 1.0+Delta/2)})
		assert.Nil(t, Ln(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta)).IntersectionCircle(circle))
	})
	t.Run("a tangent beyond the segment is missed", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(1.0, 1.0), Pt(2.0, 1.0)).IntersectionCircle(circle))
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-2.0, 2.0), Pt(2.0, 2.0)).IntersectionCircle(circle))
		assert.Nil(t, Ln(Pt(2.0, 0.0), Pt(3.0, 0.0)).IntersectionCircle(circle))
		assert.Nil(t, Ln(Pt(-0.5, 0.0), Pt(0.5, 0.0)).IntersectionCircle(circle))
	})
	t.Run("a degenerate segment is a point on the boundary or nothing", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(1.0, 0.0), Pt(1.0, 0.0)).IntersectionCircle(circle), []Point[float64]{Pt(1.0, 0.0)})
		assert.Nil(t, Ln(Pt(0.5, 0.0), Pt(0.5, 0.0)).IntersectionCircle(circle))
	})
	t.Run("a negative radius gives none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0)).IntersectionCircle(Circ(Pt(0.0, 0.0), -1.0)))
	})
	t.Run("int rounds the crossings", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-5, -5), Pt(5, 5)).IntersectionCircle(Circ(Pt(0, 0), 5)), []Point[int]{Pt(-4, -4), Pt(4, 4)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsCircle holds", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, c := range circleFixtures {
				points := l.IntersectionCircle(c)

				for _, p := range points {
					assert.True(t, l.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", l, c, p))
					assert.True(t, EqualDelta(c.Center.DistanceTo(p), c.Radius, Delta), fmt.Sprintf("%s → %s: %s on the boundary: ", l, c, p))
				}
				if len(points) > 0 {
					assert.True(t, l.IntersectsCircle(c), fmt.Sprintf("%s → %s: ", l, c))
				}
			}
		}
	})
}

func TestLine_IntersectionRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-5, 0), Pt(5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, 0), Pt(2, 0)})
		AssertVertices(t, Ln(Pt(5, 0), Pt(-5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 0), Pt(-2, 0)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(0, 0), Pt(5, 0)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 0)})
	})
	t.Run("through a corner counts it once", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(1, 3), Pt(3, 1)).IntersectionRectangle(rectangle), []Point[int]{Pt(2, 2)})
		AssertVertices(t, Ln(Pt(-4, -4), Pt(4, 4)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, -2), Pt(2, 2)})
	})
	t.Run("along an edge crosses the edges at its ends", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-5, -2), Pt(5, -2)).IntersectionRectangle(rectangle), []Point[int]{Pt(-2, -2), Pt(2, -2)})
		AssertVertices(t, Ln(Pt(-1, 2), Pt(1, 2)).IntersectionRectangle(rectangle), nil)
	})
	t.Run("apart and inside give none", func(t *testing.T) {
		assert.Nil(t, Ln(Pt(3, -5), Pt(3, 5)).IntersectionRectangle(rectangle))
		assert.Nil(t, Ln(Pt(-1, -1), Pt(1, 1)).IntersectionRectangle(rectangle))
	})
	t.Run("float keeps the crossings", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(-5.0, 1.0), Pt(5.0, 1.0)).IntersectionRectangle(Rect(Pt(0.0, 0.0), Sz(3.0, 3.0))), []Point[float64]{Pt(-1.5, 1.0), Pt(1.5, 1.0)})
	})
	t.Run("every point lies on the segment and the boundary, and exists where IntersectsRectangle holds", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, r := range rectFixtures {
				points := l.IntersectionRectangle(r)

				for _, p := range points {
					assert.True(t, l.Contains(p), fmt.Sprintf("%s → %s: %s on the segment: ", l, r, p))
					assert.True(t, slices.ContainsFunc(r.Edges(), func(edge Line[float64]) bool {
						return edge.Contains(p)
					}), fmt.Sprintf("%s → %s: %s on the boundary: ", l, r, p))
				}
				if len(points) > 0 {
					assert.True(t, l.IntersectsRectangle(r), fmt.Sprintf("%s → %s: ", l, r))
				}
			}
		}
	})
}

func TestLine_IntersectsCircle(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 1.0)

	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0)).IntersectsCircle(circle))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Ln(Pt(-2.0, 2.0), Pt(2.0, 2.0)).IntersectsCircle(circle))
		assert.False(t, Ln(Pt(2.0, 0.0), Pt(3.0, 0.0)).IntersectsCircle(circle))
	})
	t.Run("tangent counts", func(t *testing.T) {
		assert.True(t, Ln(Pt(-2.0, 1.0), Pt(2.0, 1.0)).IntersectsCircle(circle))
		assert.False(t, Ln(Pt(-2.0, 1.0+2*Delta), Pt(2.0, 1.0+2*Delta)).IntersectsCircle(circle))
	})
	t.Run("an endpoint inside counts", func(t *testing.T) {
		assert.True(t, Ln(Pt(0.5, 0.0), Pt(5.0, 0.0)).IntersectsCircle(circle))
	})
	t.Run("a segment inside counts", func(t *testing.T) {
		assert.True(t, Ln(Pt(-0.5, 0.0), Pt(0.5, 0.0)).IntersectsCircle(circle))
	})
	t.Run("a negative radius intersects nothing", func(t *testing.T) {
		assert.False(t, Ln(Pt(-2.0, 0.0), Pt(2.0, 0.0)).IntersectsCircle(Circ(Pt(0.0, 0.0), -1.0)))
	})
}

func TestLine_IntersectsRectangle(t *testing.T) {
	rectangle := Rect(Pt(0, 0), Sz(4, 4))

	t.Run("crossing an edge", func(t *testing.T) {
		assert.True(t, Ln(Pt(0, 0), Pt(5, 0)).IntersectsRectangle(rectangle))
	})
	t.Run("passing through", func(t *testing.T) {
		assert.True(t, Ln(Pt(-5, 0), Pt(5, 0)).IntersectsRectangle(rectangle))
	})
	t.Run("inside", func(t *testing.T) {
		assert.True(t, Ln(Pt(-1, -1), Pt(1, 1)).IntersectsRectangle(rectangle))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, Ln(Pt(3, -5), Pt(3, 5)).IntersectsRectangle(rectangle))
		assert.False(t, Ln(Pt(3, 3), Pt(5, 5)).IntersectsRectangle(rectangle))
	})
	t.Run("touching a corner counts", func(t *testing.T) {
		assert.True(t, Ln(Pt(1, 3), Pt(3, 1)).IntersectsRectangle(rectangle))
		assert.False(t, Ln(Pt(2, 4), Pt(4, 2)).IntersectsRectangle(rectangle))
	})
}

func BenchmarkLine_IntersectsRectangle(b *testing.B) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(10.0, 10.0))
	through, apart := Ln(Pt(-20.0, 0.0), Pt(20.0, 0.0)), Ln(Pt(-20.0, 20.0), Pt(20.0, 20.0))

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

func TestLine_IntersectionPolygon(t *testing.T) {
	t.Run("matches Polygon.IntersectionLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range polygonFixtures() {
				AssertVertices(t, l.IntersectionPolygon(p), p.IntersectionLine(l), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_IntersectsPolygon(t *testing.T) {
	t.Run("mirrors Polygon.IntersectsLine", func(t *testing.T) {
		for _, l := range lineFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, l.IntersectsPolygon(p), p.IntersectsLine(l), fmt.Sprintf("%s → %s: ", l, p))
			}
		}
	})
}

func TestLine_Equal(t *testing.T) {
	t.Run("same line", func(t *testing.T) {
		assert.True(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 5))))
		assert.True(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.6, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("different line", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 4))))
		assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.5, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("orientation matters", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(3, 5), Pt(1, 2))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.6, -0.250001), Pt(1.2, 3.4))))
	})
}

func TestLine_IsZero(t *testing.T) {
	t.Run("zero line", func(t *testing.T) {
		assert.True(t, Line[int]{}.IsZero())
		assert.True(t, Ln(Pt(0, 0), Pt(0, 0)).IsZero())
		assert.True(t, Line[float64]{}.IsZero())
	})
	t.Run("one endpoint off the origin", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 0), Pt(0, 0)).IsZero())
		assert.False(t, Ln(Pt(0, 0), Pt(0, 1)).IsZero())
	})
	t.Run("non-zero line", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).IsZero())
		assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(0.0, 0.000001), Pt(0.0, 0.0)).IsZero())
	})
}

func TestLine_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Int(), Ln(Pt(1, 2), Pt(3, 5)))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Int(), Ln(Pt(1, 0), Pt(1, 3)))
	})
}

func TestLine_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Float(), Ln(Pt(1.0, 2.0), Pt(3.0, 5.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Float(), Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)))
	})
}

func TestLine_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(10, 16), Pt(1, 2)).String(), "Ln((10,16);(1,2))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)).String(), "Ln((100.00,-34.00);(0.20,0.40))")
	})
}

func TestLine_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Ln(Pt(10, 16), Pt(1, 2)), `{"s":{"x":10,"y":16},"e":{"x":1,"y":2}}`)

		var l Line[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":10,"y":16},"e":{"x":1,"y":2}}`), &l))
		AssertLine(t, l, Ln(Pt(10, 16), Pt(1, 2)))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)), `{"s":{"x":100.0,"y":-34.0000115},"e":{"x":0.2,"y":0.4}}`)

		var l Line[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"s":{"x":10.1,"y":-34.0000115},"e":{"x":0.2,"y":0.4}}`), &l))
		AssertLine(t, l, Ln(Pt(10.1, -34.0000115), Pt(0.2, 0.4)))
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
	t.Run("reverse is its own inverse", func(t *testing.T) {
		for _, line := range lineFixtures {
			assert.True(t, line.Reverse().Reverse().Equal(line), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("reverse keeps the length and midpoint", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertNumber(t, line.Reverse().Length(), line.Length(), fmt.Sprintf("%s: ", line))
			assert.True(t, line.Reverse().Midpoint().Equal(line.Midpoint()), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("length is the vector length", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertNumber(t, line.Length(), line.Vector().Length(), fmt.Sprintf("%s: ", line))
			AssertNumber(t, line.Length(), line.Start.DistanceTo(line.End), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("midpoint is equidistant from both ends", func(t *testing.T) {
		for _, line := range lineFixtures {
			midpoint := line.Midpoint()

			AssertNumber(t, midpoint.DistanceTo(line.Start), midpoint.DistanceTo(line.End), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("translate keeps the vector", func(t *testing.T) {
		for _, line := range lineFixtures {
			for _, vector := range vectorFixtures {
				assert.True(t, line.Translate(vector).Vector().Equal(line.Vector()), fmt.Sprintf("%s → %s: ", line, vector))
			}
		}
	})
	t.Run("move to keeps the vector and starts where asked", func(t *testing.T) {
		for _, line := range lineFixtures {
			for _, point := range pointFixtures {
				moved := line.MoveTo(point)

				assert.True(t, moved.Start.Equal(point), fmt.Sprintf("%s → %s: ", line, point))
				assert.True(t, moved.Vector().Equal(line.Vector()), fmt.Sprintf("%s → %s: ", line, point))
			}
		}
	})
	t.Run("bounds span the endpoints", func(t *testing.T) {
		for _, line := range lineFixtures {
			bounds := line.Bounds()
			start, end := line.Start, line.End

			assert.True(t, bounds.Min().Equal(Pt(min(start.X, end.X), min(start.Y, end.Y))), fmt.Sprintf("%s: ", line))
			assert.True(t, bounds.Max().Equal(Pt(max(start.X, end.X), max(start.Y, end.Y))), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("vertices are the endpoints", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertVertices(t, line.Vertices(), []Point[float64]{line.Start, line.End})
		}
	})
}

func TestLine_Immutable(t *testing.T) {
	l := Ln(Pt(1, 2), Pt(3, 5))

	l.Translate(Vec(3, -2))
	l.MoveTo(Pt(4, 3))
	l.Reverse()

	AssertLine(t, l, Ln(Pt(1, 2), Pt(3, 5)))
}

// lineFixtures span axis-aligned, diagonal, degenerate, and backwards lines.
var lineFixtures = []Line[float64]{
	Ln(Pt(0.0, 0.0), Pt(0.0, 0.0)),
	Ln(Pt(1.0, 2.0), Pt(3.0, 5.0)),
	Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)),
	Ln(Pt(-3.5, 0.25), Pt(-3.5, 7.0)),
	Ln(Pt(2.0, 2.0), Pt(-4.0, 2.0)),
	Ln(Pt(12.5, -0.1), Pt(-0.5, 12.75)),
}

func ExampleLn() {
	fmt.Println(Ln(Pt(1, 2), Pt(3, 5)))
	// Output: Ln((1,2);(3,5))
}

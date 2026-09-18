package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestPolygon_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()), Polygon[int]{Vertices: squareVertices()})
	})
	t.Run("float", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()), Polygon[float64]{Vertices: triangleVertices()})
	})
}

func TestPolygon_MinMax(t *testing.T) {
	t.Run("spans the vertices", func(t *testing.T) {
		a, b := Pol([]Point[int]{Pt(3, 1), Pt(-2, 4), Pt(0, 0)}).MinMax()

		AssertPoint(t, a, Pt(-2, 0))
		AssertPoint(t, b, Pt(3, 4))
	})
	t.Run("an empty polygon has zero corners", func(t *testing.T) {
		a, b := Pol[int](nil).MinMax()

		AssertPoint(t, a, Pt(0, 0))
		AssertPoint(t, b, Pt(0, 0))
	})
	t.Run("matches the corners of Bounds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			a, b := p.MinMax()
			c, d := p.Bounds().MinMax()

			AssertPoint(t, a, c, p.String())
			AssertPoint(t, b, d, p.String())
		}
	})
}

func TestPolygon_Edges(t *testing.T) {
	t.Run("closes back to the first vertex", func(t *testing.T) {
		edges := Pol(squareVertices()).Edges()

		assert.Equal(t, len(edges), 4)
		AssertLine(t, edges[0], Ln(Pt(0, 0), Pt(2, 0)))
		AssertLine(t, edges[3], Ln(Pt(0, 2), Pt(0, 0)))
	})
	t.Run("single vertex is one zero-length edge", func(t *testing.T) {
		edges := Pol([]Point[int]{Pt(1, 1)}).Edges()

		assert.Equal(t, len(edges), 1)
		AssertLine(t, edges[0], Ln(Pt(1, 1), Pt(1, 1)))
	})
	t.Run("nil stays nil and empty stays empty", func(t *testing.T) {
		assert.Nil(t, Pol[int](nil).Edges())
		assert.Equal(t, len(Pol([]Point[int]{}).Edges()), 0)
	})
}

func TestPolygon_Center(t *testing.T) {
	t.Run("int rounds the average", func(t *testing.T) {
		AssertPoint(t, Pol(squareVertices()).Center(), Pt(1, 1))
		AssertPoint(t, Pol([]Point[int]{Pt(-1, -2), Pt(0, 0), Pt(0, 0)}).Center(), Pt(0, -1))
	})
	t.Run("narrow integers do not overflow the sum", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[int8]{Pt[int8](100, 100), Pt[int8](100, 100), Pt[int8](100, 100)}).Center(), Pt[int8](100, 100))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Pol(triangleVertices()).Center(), Pt(1.5, 0.5))
	})
	t.Run("a vertex on an edge does not move the centroid", func(t *testing.T) {
		subdivided := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)})

		AssertPoint(t, subdivided.Center(), Pt(1.0, 1.0))
	})
	t.Run("winding does not matter", func(t *testing.T) {
		reversed := Pol([]Point[float64]{Pt(0.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 0.0), Pt(1.0, 0.0), Pt(0.0, 0.0)})

		AssertPoint(t, reversed.Center(), Pt(1.0, 1.0))
	})
	t.Run("degenerate float vertices are their own center exactly", func(t *testing.T) {
		repeated := Pol([]Point[float64]{Pt(13.5, 1.9), Pt(13.5, 1.9)})

		assert.Equal(t, repeated.Center(), Pt(13.5, 1.9))
	})
	t.Run("collinear falls back to the vertex average", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)}).Center(), Pt(4.0/3, 4.0/3))
		AssertPoint(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0)}).Center(), Pt(2, 0))
	})
	t.Run("empty is the zero point", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[int]{}).Center(), Pt(0, 0))
		AssertPoint(t, Polygon[float64]{}.Center(), Pt(0.0, 0.0))
	})
}

func TestPolygon_Area(t *testing.T) {
	t.Run("square", func(t *testing.T) {
		AssertNumber(t, Pol(squareVertices()).Area(), 4.0)
	})
	t.Run("winding does not matter", func(t *testing.T) {
		reversed := Pol([]Point[int]{Pt(0, 2), Pt(2, 2), Pt(2, 0), Pt(0, 0)})

		AssertNumber(t, reversed.Area(), 4.0)
	})
	t.Run("lattice triangle encloses half units", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 0), Pt(0, 1)}).Area(), 0.5)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Pol(triangleVertices()).Area(), 0.75)
	})
	t.Run("degenerate is zero", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{}).Area(), 0.0)
		AssertNumber(t, Pol([]Point[int]{Pt(1, 1), Pt(4, 4)}).Area(), 0.0)
	})
	t.Run("degenerate float is exactly zero", func(t *testing.T) {
		assert.Equal(t, Pol([]Point[float64]{Pt(13.5, 1.9), Pt(13.5, 1.9)}).Area(), 0.0)
		assert.Equal(t, Pol([]Point[float64]{Pt(0.1, 0.2), Pt(0.3, 0.6), Pt(0.1, 0.2)}).Area(), 0.0)
	})
}

func TestPolygon_Perimeter(t *testing.T) {
	t.Run("square", func(t *testing.T) {
		AssertNumber(t, Pol(squareVertices()).Perimeter(), 8.0)
	})
	t.Run("right triangle", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{Pt(0, 0), Pt(3, 0), Pt(0, 4)}).Perimeter(), 12.0)
	})
	t.Run("two vertices count the segment twice", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{Pt(0, 0), Pt(3, 4)}).Perimeter(), 10.0)
	})
	t.Run("empty is zero", func(t *testing.T) {
		AssertNumber(t, Polygon[float64]{}.Perimeter(), 0.0)
	})
}

func TestPolygon_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRectangle(t, Pol(squareVertices()).Bounds(), RectangleFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, Pol(triangleVertices()).Bounds(), RectangleFromMinMax(Pt(0.0, 0.0), Pt(2.5, 1.0)))
	})
	t.Run("vertex order does not matter", func(t *testing.T) {
		AssertRectangle(t, Pol([]Point[int]{Pt(3, -1), Pt(-2, 4), Pt(0, 0)}).Bounds(), RectangleFromMinMax(Pt(-2, -1), Pt(3, 4)))
	})
	t.Run("empty is the zero rectangle", func(t *testing.T) {
		AssertRectangle(t, Polygon[int]{}.Bounds(), Rectangle[int]{})
	})
}

func TestPolygon_Translate(t *testing.T) {
	AssertPolygon(t, Pol(squareVertices()).Translate(Vec(1, -1)), Pol([]Point[int]{
		Pt(1, -1),
		Pt(3, -1),
		Pt(3, 1),
		Pt(1, 1),
	}))
}

func TestPolygon_MoveTo(t *testing.T) {
	AssertPolygon(t, Pol(squareVertices()).MoveTo(Pt(10, 10)), Pol([]Point[int]{
		Pt(9, 9),
		Pt(11, 9),
		Pt(11, 11),
		Pt(9, 11),
	}))

	t.Run("int lands on the point when the average crosses zero", func(t *testing.T) {
		moved := Pol([]Point[int]{Pt(-1, 0), Pt(0, 0), Pt(0, 0)}).MoveTo(Pt(1, 0))

		AssertPoint(t, moved.Center(), Pt(1, 0))
	})
	t.Run("int misses by one when a half average changes sign", func(t *testing.T) {
		moved := Pol([]Point[int]{Pt(-1, 0), Pt(0, 0)}).MoveTo(Pt(0, 0))

		AssertPoint(t, moved.Center(), Pt(1, 0))
	})
}

func TestPolygon_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()).Scale(2), Pol([]Point[int]{
			Pt(-1, -1),
			Pt(3, -1),
			Pt(3, 3),
			Pt(-1, 3),
		}))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()).ScaleXY(0.5, 2.5), Pol([]Point[float64]{
			Pt(0.75, -0.75),
			Pt(2.0, 0.5),
			Pt(1.75, 1.75),
		}))
	})
}

func TestPolygon_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()).Scale(2).Unscale(2), Pol(squareVertices()))
		AssertPolygon(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(0.0, 4.0)}).Unscale(2), Pol([]Point[float64]{Pt(1.0, 1.0), Pt(3.0, 1.0), Pt(3.0, 3.0), Pt(1.0, 3.0)}))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()).ScaleXY(0.5, 2.5).UnscaleXY(0.5, 2.5), Pol(triangleVertices()))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Pol[int](nil).Unscale(2).IsZero())
		assert.True(t, Pol[int](nil).UnscaleXY(2, 3).IsZero())
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Pol(squareVertices()).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Pol(squareVertices()).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestPolygon_Transform(t *testing.T) {
	t.Run("applies the matrix to every vertex", func(t *testing.T) {
		matrix := Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)
		square := Pol(squareVertices())

		AssertPolygon(t, square.Transform(matrix), Pol([]Point[int]{Pt(3, 7), Pt(6, 15), Pt(10, 26), Pt(8, 18)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Pol[int](nil).Transform(IdentityMatrix[float64]()).IsZero())
	})
}

func TestPolygon_Rotate(t *testing.T) {
	t.Run("quarter turn about the centroid", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 2), Pt(0, 2)}).Rotate(Pi/2), Pol([]Point[int]{
			Pt(3, -1),
			Pt(3, 3),
			Pt(1, 3),
			Pt(1, -1),
		}))
	})
	t.Run("float keeps the centroid, area and perimeter", func(t *testing.T) {
		p := Pol(triangleVertices())
		rotated := p.Rotate(0.7)

		AssertPoint(t, rotated.Center(), p.Center())
		AssertNumber(t, rotated.Area(), p.Area())
		AssertNumber(t, rotated.Perimeter(), p.Perimeter())
	})
	t.Run("a full turn is identity", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			AssertPolygon(t, p.Rotate(2*Pi), p, fmt.Sprintf("%s: ", p))
		}
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Pol[int](nil).Rotate(1).IsZero())
	})
}

func TestPolygon_Contains(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("inside", func(t *testing.T) {
		assert.True(t, square.Contains(Pt(1, 1)))
	})
	t.Run("outside on every side", func(t *testing.T) {
		assert.False(t, square.Contains(Pt(-1, 1)))
		assert.False(t, square.Contains(Pt(3, 1)))
		assert.False(t, square.Contains(Pt(1, -1)))
		assert.False(t, square.Contains(Pt(1, 3)))
	})
	t.Run("boundary and vertices are included", func(t *testing.T) {
		assert.True(t, square.Contains(Pt(2, 1)))
		assert.True(t, square.Contains(Pt(0, 0)))
		assert.True(t, square.Contains(Pt(2, 2)))
	})
	t.Run("outside the extent is rejected on the point's own row", func(t *testing.T) {
		assert.False(t, square.Contains(Pt(-1, 0)))
		assert.False(t, square.Contains(Pt(3, 2)))
	})
	t.Run("edges running down and up count alike", func(t *testing.T) {
		clockwise := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)})
		counter := Pol([]Point[float64]{Pt(0.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 0.0), Pt(0.0, 0.0)})

		for _, point := range []Point[float64]{Pt(1.0, 1.0), Pt(0.5, 1.5), Pt(1.9, 0.1)} {
			assert.True(t, clockwise.Contains(point), point.String())
			assert.True(t, counter.Contains(point), point.String())
		}
	})
	t.Run("ray through a vertex is counted once", func(t *testing.T) {
		diamond := Pol([]Point[int]{Pt(0, -2), Pt(2, 0), Pt(0, 2), Pt(-2, 0)})

		assert.True(t, diamond.Contains(Pt(-1, 0)))
		assert.False(t, diamond.Contains(Pt(-3, 0)))
		assert.False(t, diamond.Contains(Pt(3, 0)))
	})
	t.Run("concave notch is outside", func(t *testing.T) {
		notched := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(2, 1), Pt(0, 4)})

		assert.True(t, notched.Contains(Pt(1, 1)))
		assert.False(t, notched.Contains(Pt(2, 3)))
	})
	t.Run("float is tolerant at the boundary", func(t *testing.T) {
		triangle := Pol(triangleVertices())

		assert.True(t, triangle.Contains(Pt(2.0, 0.5)))
		assert.True(t, triangle.Contains(Pt(1.25, 0.25-Delta/2)))
		assert.False(t, triangle.Contains(Pt(1.25, 0.25-2*Delta)))
	})
	t.Run("degenerate polygons contain only their points", func(t *testing.T) {
		assert.False(t, Pol([]Point[int]{}).Contains(Pt(0, 0)))
		assert.True(t, Pol([]Point[int]{Pt(1, 1)}).Contains(Pt(1, 1)))
		assert.True(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0)}).Contains(Pt(2, 0)))
		assert.False(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0)}).Contains(Pt(2, 1)))
	})
}

func BenchmarkPolygon_Contains(b *testing.B) {
	polygon := benchPolygon()
	inside, outside, far := Pt(10.0, 20.0), Pt(99.0, 99.0), Pt(500.0, 500.0)

	b.Run("inside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Contains(inside)
		}
	})
	b.Run("outside within the extent", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Contains(outside)
		}
	})
	b.Run("outside the extent", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Contains(far)
		}
	})
}

func TestPolygon_DistanceTo(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("beside an edge measures to the edge", func(t *testing.T) {
		AssertNumber(t, square.DistanceTo(Pt(5, 1)), 3.0)
		AssertNumber(t, square.DistanceTo(Pt(1, -2)), 2.0)
	})
	t.Run("beyond a vertex measures to the vertex", func(t *testing.T) {
		AssertNumber(t, square.DistanceTo(Pt(5, 6)), 5.0)
	})
	t.Run("inside and on the boundary are zero", func(t *testing.T) {
		assert.Equal(t, square.DistanceTo(Pt(1, 1)), 0.0)
		assert.Equal(t, square.DistanceTo(Pt(2, 1)), 0.0)
	})
	t.Run("a concave notch measures to the notch edges", func(t *testing.T) {
		notched := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(4.0, 0.0), Pt(4.0, 4.0), Pt(2.0, 2.0), Pt(0.0, 4.0)})

		AssertNumber(t, notched.DistanceTo(Pt(2.0, 4.0)), Sqrt2)
	})
	t.Run("an empty polygon is infinitely far", func(t *testing.T) {
		assert.True(t, math.IsInf(Pol[int](nil).DistanceTo(Pt(0, 0)), 1))
	})
	t.Run("zero exactly where Contains holds", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, p := range pointFixtures {
				assert.Equal(t, polygon.DistanceTo(p) == 0, polygon.Contains(p), fmt.Sprintf("%s → %s: ", polygon, p))
			}
		}
	})
}

func BenchmarkPolygon_DistanceTo(b *testing.B) {
	polygon := benchPolygon()
	point := Pt(150.0, 150.0)

	for b.Loop() {
		_ = polygon.DistanceTo(point)
	}
}

func TestPolygon_DistanceSquaredTo(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		AssertNumber(t, square.DistanceSquaredTo(Pt(5, 6)), 25.0)
		AssertNumber(t, square.DistanceSquaredTo(Pt(5, 1)), 9.0)
		assert.Equal(t, square.DistanceSquaredTo(Pt(1, 1)), 0.0)
	})
	t.Run("stays fractional for an integer T", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 1)}).DistanceSquaredTo(Pt(0, 1)), 0.8)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, p := range pointFixtures {
				AssertNumber(t, polygon.DistanceSquaredTo(p), polygon.DistanceTo(p)*polygon.DistanceTo(p), fmt.Sprintf("%s → %s: ", polygon, p))
			}
		}
	})
}

func TestPolygon_Intersects(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, square.Intersects(square.Translate(Vec(1, 1))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.Intersects(square.Translate(Vec(3, 0))))
	})
	t.Run("a shared edge counts", func(t *testing.T) {
		assert.True(t, square.Intersects(square.Translate(Vec(2, 0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, square.Intersects(Pol([]Point[int]{Pt(1, 1), Pt(1, 1)})))
		assert.True(t, Pol([]Point[int]{Pt(1, 1), Pt(1, 1)}).Intersects(square))
	})
	t.Run("edges crossing without a vertex inside", func(t *testing.T) {
		cross := Pol([]Point[int]{Pt(-1, 1), Pt(3, 1), Pt(3, 1), Pt(-1, 1)})
		plus := Pol([]Point[int]{Pt(1, -1), Pt(1, 3), Pt(1, 3), Pt(1, -1)})

		assert.True(t, cross.Intersects(plus))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, square.Intersects(Pol[int](nil)))
		assert.False(t, Pol[int](nil).Intersects(square))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range polygonFixtures() {
			for _, b := range polygonFixtures() {
				assert.Equal(t, a.Intersects(b), b.Intersects(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func BenchmarkPolygon_Intersects(b *testing.B) {
	polygon := benchPolygon()
	overlapping := polygon.Translate(Vec(150.0, 0.0))
	apartWithinBounds := polygon.Translate(Vec(150.0, 150.0))
	apart := polygon.Translate(Vec(300.0, 0.0))

	b.Run("overlapping", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Intersects(overlapping)
		}
	})
	b.Run("apart within overlapping bounds", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Intersects(apartWithinBounds)
		}
	})
	b.Run("apart", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.Intersects(apart)
		}
	})
}

func TestPolygon_IntersectsLine(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("crossing an edge", func(t *testing.T) {
		assert.True(t, square.IntersectsLine(Ln(Pt(1, 1), Pt(5, 1))))
		assert.True(t, square.IntersectsLine(Ln(Pt(5, 1), Pt(1, 1))))
	})
	t.Run("passing through", func(t *testing.T) {
		assert.True(t, square.IntersectsLine(Ln(Pt(-1, 1), Pt(5, 1))))
	})
	t.Run("inside", func(t *testing.T) {
		assert.True(t, square.IntersectsLine(Ln(Pt(1, 1), Pt(1, 1))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsLine(Ln(Pt(3, -1), Pt(3, 3))))
	})
	t.Run("touching a vertex counts", func(t *testing.T) {
		assert.True(t, square.IntersectsLine(Ln(Pt(1, 3), Pt(3, 1))))
	})
	t.Run("outside the extent is rejected", func(t *testing.T) {
		assert.False(t, square.IntersectsLine(Ln(Pt(3, 3), Pt(5, 5))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).IntersectsLine(Ln(Pt(0, 0), Pt(1, 1))))
	})
	t.Run("matches Intersects on the segment as a polygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, l := range lineFixtures {
				assert.Equal(t, p.IntersectsLine(l), p.Intersects(Pol(l.Vertices())), fmt.Sprintf("%s → %s: ", p, l))
			}
		}
	})
}

func TestPolygon_IntersectionLine(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("passing through gives both crossings from Start to End", func(t *testing.T) {
		AssertVertices(t, square.IntersectionLine(Ln(Pt(-1, 1), Pt(5, 1))), []Point[int]{Pt(0, 1), Pt(2, 1)})
		AssertVertices(t, square.IntersectionLine(Ln(Pt(5, 1), Pt(-1, 1))), []Point[int]{Pt(2, 1), Pt(0, 1)})
	})
	t.Run("ending inside gives one crossing", func(t *testing.T) {
		AssertVertices(t, square.IntersectionLine(Ln(Pt(1, 1), Pt(5, 1))), []Point[int]{Pt(2, 1)})
	})
	t.Run("through a vertex counts it once", func(t *testing.T) {
		AssertVertices(t, square.IntersectionLine(Ln(Pt(1, 3), Pt(3, 1))), []Point[int]{Pt(2, 2)})
	})
	t.Run("a concave polygon is crossed more than twice", func(t *testing.T) {
		notched := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(2, 1), Pt(0, 4)})

		AssertVertices(t, notched.IntersectionLine(Ln(Pt(-1, 3), Pt(5, 3))), []Point[int]{Pt(0, 3), Pt(1, 3), Pt(3, 3), Pt(4, 3)})
	})
	t.Run("inside, apart and empty give none", func(t *testing.T) {
		assert.Nil(t, square.IntersectionLine(Ln(Pt(1, 1), Pt(1, 1))))
		assert.Nil(t, square.IntersectionLine(Ln(Pt(3, -1), Pt(3, 3))))
		assert.Nil(t, Pol[int](nil).IntersectionLine(Ln(Pt(-1, 1), Pt(5, 1))))
	})
	t.Run("matches the rectangle crossings on the rectangle as a polygon", func(t *testing.T) {
		for _, r := range rectFixtures {
			for _, l := range lineFixtures {
				AssertVertices(t, r.Polygon().IntersectionLine(l), l.IntersectionRectangle(r), fmt.Sprintf("%s → %s: ", r, l))
			}
		}
	})
	t.Run("every point lies on the segment and an edge, and exists where IntersectsLine holds", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, l := range lineFixtures {
				points := p.IntersectionLine(l)

				for _, point := range points {
					assert.True(t, l.Contains(point), fmt.Sprintf("%s → %s: %s on the segment: ", p, l, point))
					assert.True(t, slices.ContainsFunc(p.Edges(), func(edge Line[float64]) bool {
						return edge.Contains(point)
					}), fmt.Sprintf("%s → %s: %s on the boundary: ", p, l, point))
				}
				if len(points) > 0 {
					assert.True(t, p.IntersectsLine(l), fmt.Sprintf("%s → %s: ", p, l))
				}
			}
		}
	})
}

func TestPolygon_IntersectsRectangle(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, square.IntersectsRectangle(Rect(Pt(2, 2), Sz(2, 2))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsRectangle(Rect(Pt(4, 4), Sz(2, 2))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, square.IntersectsRectangle(Rect(Pt(1, 1), Sz(20, 20))))
		assert.True(t, square.IntersectsRectangle(Rect(Pt(1, 1), Sz(1, 1))))
	})
	t.Run("edges crossing without a corner inside", func(t *testing.T) {
		assert.True(t, square.IntersectsRectangle(Rect(Pt(1, 1), Sz(6, 1))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).IntersectsRectangle(Rect(Pt(1, 1), Sz(2, 2))))
	})
	t.Run("matches the rectangle as a polygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, r := range rectFixtures {
				assert.Equal(t, p.IntersectsRectangle(r), p.Intersects(r.Polygon()), fmt.Sprintf("%s → %s: ", p, r))
			}
		}
	})
}

func TestPolygon_IntersectsCircle(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("center inside", func(t *testing.T) {
		assert.True(t, square.IntersectsCircle(Circ(Pt(1, 1), 5)))
	})
	t.Run("an edge within the radius", func(t *testing.T) {
		assert.True(t, square.IntersectsCircle(Circ(Pt(3, 1), 1)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsCircle(Circ(Pt(4, 1), 1)))
		assert.False(t, square.IntersectsCircle(Circ(Pt(3, 3), 1)))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).IntersectsCircle(Circ(Pt(0, 0), 1)))
	})
	t.Run("matches the circle test on the polygon edges", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, c := range circleFixtures {
				expected := p.Contains(c.Center) || slices.ContainsFunc(p.Edges(), func(edge Line[float64]) bool {
					return edge.IntersectsCircle(c)
				})

				assert.Equal(t, p.IntersectsCircle(c), expected, fmt.Sprintf("%s → %s: ", p, c))
			}
		}
	})
}

func TestPolygon_Equal(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("same polygon", func(t *testing.T) {
		assert.True(t, square.Equal(Pol(squareVertices())))
		assert.True(t, Pol(triangleVertices()).Equal(Pol(triangleVertices())))
	})
	t.Run("different vertex count", func(t *testing.T) {
		assert.False(t, square.Equal(Pol([]Point[int]{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {3, 2}})))
	})
	t.Run("different vertex", func(t *testing.T) {
		assert.False(t, square.Equal(Pol([]Point[int]{{0, 0}, {2, 0}, {3, 2}, {0, 2}})))
	})
}

func TestPolygon_IsZero(t *testing.T) {
	t.Run("nil vertices", func(t *testing.T) {
		assert.True(t, Polygon[int]{}.IsZero())
		assert.True(t, Polygon[float64]{}.IsZero())
	})
	t.Run("an empty slice is not nil", func(t *testing.T) {
		assert.False(t, Polygon[int]{[]Point[int]{}}.IsZero())
	})
	t.Run("nil survives every mapping", func(t *testing.T) {
		var polygon Polygon[float64]

		assert.True(t, polygon.Translate(Vec(1.0, 1.0)).IsZero())
		assert.True(t, polygon.Scale(2).IsZero())
		assert.True(t, polygon.ScaleXY(2, 3).IsZero())
		assert.True(t, polygon.Int().IsZero())
		assert.True(t, polygon.Float().IsZero())
	})
	t.Run("non-zero polygon", func(t *testing.T) {
		assert.False(t, Pol(squareVertices()).IsZero())
	})
}

func TestPolygon_Empty(t *testing.T) {
	t.Run("no vertices", func(t *testing.T) {
		assert.True(t, Polygon[int]{}.Empty())
		assert.True(t, Polygon[int]{[]Point[int]{}}.Empty())
	})
	t.Run("with vertices", func(t *testing.T) {
		assert.False(t, Pol(squareVertices()).Empty())
		assert.False(t, Pol(triangleVertices()).Empty())
	})
}

func TestPolygon_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()).Int(), Pol(squareVertices()))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()).Int(), Pol([]Point[int]{
			Pt(0, 0),
			Pt(3, 1),
			Pt(2, 1),
		}))
	})
}

func TestPolygon_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()).Float(), Pol([]Point[float64]{
			Pt(0.0, 0.0),
			Pt(2.0, 0.0),
			Pt(2.0, 2.0),
			Pt(0.0, 2.0),
		}))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()).Float(), Pol(triangleVertices()))
	})
}

func TestPolygon_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Pol(squareVertices()).String(), "Pol((0,0);(2,0);(2,2);(0,2))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Pol(triangleVertices()).String(), "Pol((0.00,0.00);(2.50,0.50);(2.00,1.00))")
	})
}

func TestPolygon_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Pol(squareVertices()), `[{"x":0,"y":0},{"x":2,"y":0},{"x":2,"y":2},{"x":0,"y":2}]`)

		var p Polygon[int]
		assert.NoError(t, json.Unmarshal([]byte(`[{"x":0,"y":0},{"x":2,"y":0},{"x":2,"y":2},{"x":0,"y":2}]`), &p))
		AssertPolygon(t, p, Pol(squareVertices()))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Pol(triangleVertices()), `[{"x":0,"y":0},{"x":2.5,"y":0.5},{"x":2,"y":1}]`)

		var p Polygon[float64]
		assert.NoError(t, json.Unmarshal([]byte(`[{"x":0,"y":0},{"x":2.5,"y":0.5},{"x":2,"y":1}]`), &p))
		AssertPolygon(t, p, Pol(triangleVertices()))
	})
	t.Run("decoding leaves a shared slice untouched", func(t *testing.T) {
		shared := squareVertices()
		p := Pol(shared)

		assert.NoError(t, json.Unmarshal([]byte(`[{"x":9,"y":9}]`), &p))
		AssertVertices(t, shared, squareVertices())
		AssertPolygon(t, p, Pol([]Point[int]{{9, 9}}))
	})
	t.Run("invalid input leaves the polygon untouched", func(t *testing.T) {
		p := Pol(squareVertices())

		assert.Error(t, json.Unmarshal([]byte(`[{"x":"nine"}]`), &p))
		AssertPolygon(t, p, Pol(squareVertices()))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			data, err := json.Marshal(polygon)
			assert.NoError(t, err)

			var decoded Polygon[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.True(t, decoded.Equal(polygon), fmt.Sprintf("%s: ", polygon))
		}
	})
}

func TestPolygon_Properties(t *testing.T) {
	t.Run("translate moves every vertex and the centroid", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, vector := range vectorFixtures {
				moved := polygon.Translate(vector)

				assert.True(t, moved.Center().Equal(polygon.Center().Add(vector)), fmt.Sprintf("%s → %s: ", polygon, vector))
				for i, vertex := range moved.Vertices {
					assert.True(t, vertex.Equal(polygon.Vertices[i].Add(vector)), fmt.Sprintf("%s → %s: ", polygon, vector))
				}
			}
		}
	})
	t.Run("move to centers where asked", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, point := range pointFixtures {
				assert.True(t, polygon.MoveTo(point).Center().Equal(point), fmt.Sprintf("%s → %s: ", polygon, point))
			}
		}
	})
	t.Run("scale keeps the centroid", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				scaled := polygon.Scale(factor)

				assert.True(t, scaled.Center().Equal(polygon.Center()), fmt.Sprintf("%s ×%v: ", polygon, factor))
			}
		}
	})
	t.Run("scale by one is the identity", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			assert.True(t, polygon.Scale(1).Equal(polygon), fmt.Sprintf("%s: ", polygon))
			assert.True(t, polygon.ScaleXY(1, 1).Equal(polygon), fmt.Sprintf("%s: ", polygon))
		}
	})
	t.Run("centroid of a triangle is the average of the vertices", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			if len(polygon.Vertices) != 3 {
				continue
			}

			var sum Point[float64]
			for _, vertex := range polygon.Vertices {
				sum = sum.AddXY(vertex.XY())
			}

			expected := sum.Divide(3)
			assert.True(t, polygon.Center().Equal(expected), fmt.Sprintf("%s: ", polygon))
		}
	})
	t.Run("centroid lies within the bounds", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			assert.True(t, polygon.Bounds().Contains(polygon.Center()), fmt.Sprintf("%s: ", polygon))
		}
	})
}

func TestPolygon_Immutable(t *testing.T) {
	p := Pol([]Point[int]{Pt(0, 0), Pt(2, 0)})

	p.Translate(Vec(1, -1))
	p.MoveTo(Pt(10, 10))
	p.Scale(2)
	p.ScaleXY(2, 3)

	AssertVertices(t, p.Vertices, []Point[int]{Pt(0, 0), Pt(2, 0)})
}

// squareVertices and triangleVertices build fresh slices, so a test that mutates one
// cannot leak into the next.
func squareVertices() []Point[int] {
	return []Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}
}

func triangleVertices() []Point[float64] {
	return []Point[float64]{Pt(0.0, 0.0), Pt(2.5, 0.5), Pt(2.0, 1.0)}
}

// polygonFixtures span a triangle, a square, and a degenerate two-vertex polygon.
func polygonFixtures() []Polygon[float64] {
	return []Polygon[float64]{
		Pol(triangleVertices()),
		Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)}),
		Pol([]Point[float64]{Pt(-3.5, 0.25), Pt(12.5, -0.1), Pt(-0.5, 12.75)}),
		Pol([]Point[float64]{Pt(1.0, 2.0), Pt(1.0, 2.0)}),
	}
}

func ExamplePol() {
	fmt.Println(Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2)}))
	// Output: Pol((0,0);(2,0);(2,2))
}

// benchPolygon is a 64-gon, large enough for the edge walk to dominate.
func benchPolygon() Polygon[float64] {
	return RegPol(Pt(0.0, 0.0), SzU(100.0), 64, 0).Polygon()
}

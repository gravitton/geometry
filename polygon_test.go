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
		AssertPolygon(t, Pol(squareVertices()), Polygon[int]{Points: squareVertices()})
	})
	t.Run("float", func(t *testing.T) {
		AssertPolygon(t, Pol(triangleVertices()), Polygon[float64]{Points: triangleVertices()})
	})
}

func TestPolygon_Vertices(t *testing.T) {
	t.Run("iterates the points in order", func(t *testing.T) {
		AssertVertices(t, slices.Collect(Pol(squareVertices()).Vertices()), squareVertices())
	})
	t.Run("nil and empty yield nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(Pol[int](nil).Vertices()))
		assert.Nil(t, slices.Collect(Pol([]Point[int]{}).Vertices()))
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			AssertNumber(t, testing.AllocsPerRun(100, func() {
				for vertex := range p.Vertices() {
					sinkBool = vertex.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", p))
		}
	})
}

func TestPolygon_Edges(t *testing.T) {
	t.Run("closes back to the first vertex", func(t *testing.T) {
		edges := slices.Collect(Pol(squareVertices()).Edges())

		assert.Equal(t, len(edges), 4)
		AssertSegment(t, edges[0], Seg(Pt(0, 0), Pt(2, 0)))
		AssertSegment(t, edges[3], Seg(Pt(0, 2), Pt(0, 0)))
	})
	t.Run("single vertex is one zero-length edge", func(t *testing.T) {
		edges := slices.Collect(Pol([]Point[int]{Pt(1, 1)}).Edges())

		assert.Equal(t, len(edges), 1)
		AssertSegment(t, edges[0], Seg(Pt(1, 1), Pt(1, 1)))
	})
	t.Run("nil and empty yield nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(Pol[int](nil).Edges()))
		assert.Nil(t, slices.Collect(Pol([]Point[int]{}).Edges()))
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for edge := range Pol(squareVertices()).Edges() {
			AssertSegment(t, edge, Seg(Pt(0, 0), Pt(2, 0)))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			AssertNumber(t, testing.AllocsPerRun(100, func() {
				for edge := range p.Edges() {
					sinkBool = edge.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", p))
		}
	})
}

func TestPolygon_Centroid(t *testing.T) {
	t.Run("int rounds the average", func(t *testing.T) {
		AssertPoint(t, Pol(squareVertices()).Centroid(), Pt(1, 1))
		AssertPoint(t, Pol([]Point[int]{Pt(-1, -2), Pt(0, 0), Pt(0, 0)}).Centroid(), Pt(0, -1))
	})
	t.Run("narrow integers do not overflow the sum", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[int8]{Pt[int8](100, 100), Pt[int8](100, 100), Pt[int8](100, 100)}).Centroid(), Pt[int8](100, 100))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Pol(triangleVertices()).Centroid(), Pt(1.5, 0.5))
	})
	t.Run("a vertex on an edge does not move the centroid", func(t *testing.T) {
		subdivided := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)})

		AssertPoint(t, subdivided.Centroid(), Pt(1.0, 1.0))
	})
	t.Run("winding does not matter", func(t *testing.T) {
		reversed := Pol([]Point[float64]{Pt(0.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 0.0), Pt(1.0, 0.0), Pt(0.0, 0.0)})

		AssertPoint(t, reversed.Centroid(), Pt(1.0, 1.0))
	})
	t.Run("degenerate float vertices are their own center exactly", func(t *testing.T) {
		repeated := Pol([]Point[float64]{Pt(13.5, 1.9), Pt(13.5, 1.9)})

		assert.Equal(t, repeated.Centroid(), Pt(13.5, 1.9))
	})
	t.Run("collinear falls back to the vertex average", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)}).Centroid(), Pt(4.0/3, 4.0/3))
		AssertPoint(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0)}).Centroid(), Pt(2, 0))
	})
	t.Run("empty is the zero point", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[int]{}).Centroid(), Pt(0, 0))
		AssertPoint(t, Polygon[float64]{}.Centroid(), Pt(0.0, 0.0))
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
	t.Run("far from the origin keeps the area", func(t *testing.T) {
		square := Pol(squareVertices()).Float()

		AssertNumber(t, square.Translate(Vec(1e9, 1e9)).Area(), square.Area())
		AssertNumber(t, Pol(triangleVertices()).Translate(Vec(1e9, -1e9)).Area(), Pol(triangleVertices()).Area())
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

func TestPolygon_Inertia(t *testing.T) {
	t.Run("square", func(t *testing.T) {
		AssertNumber(t, Pol(squareVertices()).Inertia(), 8.0/3)
	})
	t.Run("winding does not matter", func(t *testing.T) {
		reversed := Pol([]Point[int]{Pt(0, 2), Pt(2, 2), Pt(2, 0), Pt(0, 0)})

		AssertNumber(t, reversed.Inertia(), 8.0/3)
	})
	t.Run("is taken about the centroid wherever the polygon lies", func(t *testing.T) {
		square := Pol(squareVertices())

		AssertNumber(t, square.Translate(Vec(100, -250)).Inertia(), square.Inertia())
		AssertNumber(t, Pol(triangleVertices()).Translate(Vec(100.0, -250.0)).Inertia(), Pol(triangleVertices()).Inertia())
	})
	t.Run("a triangle has the moment of its sides, A(a²+b²+c²)/36", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			if len(p.Points) != 3 {
				continue
			}

			var sides float64
			for edge := range p.Edges() {
				sides += edge.Vector().LengthSquared()
			}

			AssertNumber(t, p.Inertia(), p.Area()*sides/36, p.String())
		}
	})
	t.Run("degenerate is zero", func(t *testing.T) {
		AssertNumber(t, Pol([]Point[int]{}).Inertia(), 0.0)
		AssertNumber(t, Pol([]Point[int]{Pt(1, 1), Pt(4, 4)}).Inertia(), 0.0)
		AssertNumber(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, 1.0), Pt(3.0, 3.0)}).Inertia(), 0.0)
	})
}

func TestPolygon_Winding(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("the winding of a rectangle is clockwise", func(t *testing.T) {
		assert.Equal(t, square.Winding(), WindingClockwise)
	})
	t.Run("the reverse is counterclockwise", func(t *testing.T) {
		assert.Equal(t, Pol([]Point[int]{Pt(0, 2), Pt(2, 2), Pt(2, 0), Pt(0, 0)}).Winding(), WindingCounterClockwise)
	})
	t.Run("no area has no winding", func(t *testing.T) {
		assert.Equal(t, Polygon[int]{}.Winding(), WindingNone)
		assert.Equal(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 1), Pt(2, 2)}).Winding(), WindingNone)
	})
	t.Run("lobes of equal area either way have no winding", func(t *testing.T) {
		assert.Equal(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 2), Pt(2, 0), Pt(0, 2)}).Winding(), WindingNone)
	})
	t.Run("allocates nothing", func(t *testing.T) {
		AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = square.Winding().IsNone()
		}), 0)
	})
}

func TestPolygon_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertBox(t, Pol(squareVertices()).Bounds(), BoxFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("float", func(t *testing.T) {
		AssertBox(t, Pol(triangleVertices()).Bounds(), BoxFromMinMax(Pt(0.0, 0.0), Pt(2.5, 1.0)))
	})
	t.Run("vertex order does not matter", func(t *testing.T) {
		AssertBox(t, Pol([]Point[int]{Pt(3, -1), Pt(-2, 4), Pt(0, 0)}).Bounds(), BoxFromMinMax(Pt(-2, -1), Pt(3, 4)))
	})
	t.Run("empty is the zero rectangle", func(t *testing.T) {
		AssertBox(t, Polygon[int]{}.Bounds(), Box[int]{})
	})
}

func TestPolygon_minMax(t *testing.T) {
	t.Run("spans the vertices", func(t *testing.T) {
		a, b := Pol([]Point[int]{Pt(3, 1), Pt(-2, 4), Pt(0, 0)}).minMax()

		AssertPoint(t, a, Pt(-2, 0))
		AssertPoint(t, b, Pt(3, 4))
	})
	t.Run("an empty polygon has zero corners", func(t *testing.T) {
		a, b := Pol[int](nil).minMax()

		AssertPoint(t, a, Pt(0, 0))
		AssertPoint(t, b, Pt(0, 0))
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

		AssertPoint(t, moved.Centroid(), Pt(1, 0))
	})
	t.Run("int misses by one when a half average changes sign", func(t *testing.T) {
		moved := Pol([]Point[int]{Pt(-1, 0), Pt(0, 0)}).MoveTo(Pt(0, 0))

		AssertPoint(t, moved.Centroid(), Pt(1, 0))
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

func TestPolygon_Lerp(t *testing.T) {
	a := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)})
	b := Pol([]Point[float64]{Pt(10.0, 0.0), Pt(14.0, 0.0), Pt(14.0, 8.0), Pt(10.0, 8.0)})

	t.Run("moves every vertex towards its pair", func(t *testing.T) {
		AssertPolygon(t, a.Lerp(b, 0.5), Pol([]Point[float64]{Pt(5.0, 0.0), Pt(8.0, 0.0), Pt(8.0, 5.0), Pt(5.0, 5.0)}))
	})
	t.Run("the ends are the polygons themselves", func(t *testing.T) {
		AssertPolygon(t, a.Lerp(b, 0), a)
		AssertPolygon(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates outside the unit range", func(t *testing.T) {
		AssertPolygon(t, a.Lerp(b, 2), Pol([]Point[float64]{Pt(20.0, 0.0), Pt(26.0, 0.0), Pt(26.0, 14.0), Pt(20.0, 14.0)}))
	})
	t.Run("int rounds", func(t *testing.T) {
		AssertPolygon(t, Pol(squareVertices()).Lerp(Pol([]Point[int]{Pt(5, 5), Pt(7, 5), Pt(7, 7), Pt(5, 7)}), 0.5), Pol([]Point[int]{Pt(3, 3), Pt(5, 3), Pt(5, 5), Pt(3, 5)}))
	})
	t.Run("nil stays nil", func(t *testing.T) {
		assert.True(t, Pol[int](nil).Lerp(Pol[int](nil), 0.5).IsZero())
	})
	t.Run("a different vertex count panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			a.Lerp(Pol(triangleVertices()), 0.5)
		}, "geom: lerp between polygons of 4 and 3 vertices")
	})
	t.Run("towards a translation is a part of it over the fixtures", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, vector := range vectorFixtures {
				moved := polygon.Translate(vector)

				assert.True(t, polygon.Lerp(moved, 0).Equal(polygon), fmt.Sprintf("%s → %s: ", polygon, vector))
				assert.True(t, polygon.Lerp(moved, 1).Equal(moved), fmt.Sprintf("%s → %s: ", polygon, vector))
				assert.True(t, polygon.Lerp(moved, 0.25).Equal(polygon.Translate(vector.Multiply(0.25))), fmt.Sprintf("%s → %s: ", polygon, vector))
			}
		}
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

		AssertPoint(t, rotated.Centroid(), p.Centroid())
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

func TestPolygon_ConvexHull(t *testing.T) {
	dart := Pol([]Point[int]{Pt(0, 0), Pt(4, 2), Pt(0, 4), Pt(1, 2)})

	t.Run("drops the vertices inside and on the edges", func(t *testing.T) {
		scattered := Pol([]Point[int]{Pt(1, 1), Pt(0, 0), Pt(0, 1), Pt(2, 2), Pt(2, 0), Pt(0, 2), Pt(1, 0)})

		AssertPolygon(t, scattered.ConvexHull(), Pol(squareVertices()))
	})
	t.Run("drops the vertices of a concavity", func(t *testing.T) {
		AssertPolygon(t, dart.ConvexHull(), Pol([]Point[int]{Pt(0, 0), Pt(4, 2), Pt(0, 4)}))
	})
	t.Run("winds clockwise from the least vertex whatever the input winding", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 2), Pt(2, 2), Pt(2, 0), Pt(0, 0)}).ConvexHull(), Pol(squareVertices()))
	})
	t.Run("collinear vertices give the two ends", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 2), Pt(1, 1), Pt(3, 3)}).ConvexHull(), Pol([]Point[int]{Pt(0, 0), Pt(3, 3)}))
	})
	t.Run("repeated vertices count once", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(1, 1), Pt(1, 1), Pt(1, 1)}).ConvexHull(), Pol([]Point[int]{Pt(1, 1)}))
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 0), Pt(0, 2), Pt(0, 0)}).ConvexHull(), Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(0, 2)}))
	})
	t.Run("an empty polygon is returned as it is", func(t *testing.T) {
		assert.True(t, Polygon[int]{}.ConvexHull().IsZero())
	})
	t.Run("allocates once", func(t *testing.T) {
		AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = dart.ConvexHull().Points
		}), 1)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, p := range outlineFixtures() {
			hull := p.ConvexHull()
			message := fmt.Sprintf("%s → %s: ", p, hull)

			for _, vertex := range p.Points {
				assert.True(t, hull.Contains(vertex), message)
			}
			for _, vertex := range hull.Points {
				assert.True(t, slices.Contains(p.Points, vertex), message)
			}
			assert.True(t, hull.ConvexHull().Equal(hull), message)
			if len(hull.Points) >= 3 {
				assert.True(t, hull.IsConvex(), message)
				assert.Equal(t, hull.Winding(), WindingClockwise, message)
			}
		}
	})
}

func BenchmarkPolygon_ConvexHull(b *testing.B) {
	polygon := benchPolygon()

	for b.Loop() {
		sinkBool = polygon.ConvexHull().IsEmpty()
	}
}

func TestPolygon_Simplify(t *testing.T) {
	t.Run("drops a vertex on the line between its neighbours", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}).Simplify(0), Pol(squareVertices()))
	})
	t.Run("drops a repeated vertex", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}).Simplify(0), Pol(squareVertices()))
	})
	t.Run("keeps a spike reaching beyond the edge that would replace it", func(t *testing.T) {
		spike := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(2, 0), Pt(2, 2)})

		AssertPolygon(t, spike.Simplify(0), spike)
	})
	t.Run("drops a fold running back along the edge that replaces it", func(t *testing.T) {
		folded := Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(2, 0), Pt(6, 0), Pt(6, 6), Pt(0, 6)})

		AssertPolygon(t, folded.Simplify(0), Pol([]Point[int]{Pt(0, 0), Pt(6, 0), Pt(6, 6), Pt(0, 6)}))
	})
	t.Run("drops a vertex within the tolerance", func(t *testing.T) {
		bumped := Pol([]Point[int]{Pt(0, 0), Pt(5, 1), Pt(10, 0), Pt(10, 10), Pt(0, 10)})

		AssertPolygon(t, bumped.Simplify(1), Pol([]Point[int]{Pt(0, 0), Pt(10, 0), Pt(10, 10), Pt(0, 10)}))
		AssertPolygon(t, bumped.Simplify(0.5), bumped)
	})
	t.Run("drops a run of vertices within the tolerance of the edge replacing them", func(t *testing.T) {
		wavy := Pol([]Point[int]{Pt(0, 0), Pt(10, 2), Pt(14, -1), Pt(20, 0), Pt(20, 20), Pt(0, 20)})

		AssertPolygon(t, wavy.Simplify(2.2), Pol([]Point[int]{Pt(0, 0), Pt(20, 0), Pt(20, 20), Pt(0, 20)}))
	})
	t.Run("drops the vertices where the outline closes", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(1, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2), Pt(0, 0)}).Simplify(0), Pol([]Point[int]{Pt(2, 0), Pt(2, 2), Pt(0, 2), Pt(0, 0)}))
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2), Pt(0, 1)}).Simplify(0), Pol(squareVertices()))
	})
	t.Run("collinear vertices give the two ends", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 1), Pt(2, 2), Pt(3, 3)}).Simplify(0), Pol([]Point[int]{Pt(0, 0), Pt(3, 3)}))
	})
	t.Run("an outline within the tolerance of one vertex keeps it alone", func(t *testing.T) {
		AssertPolygon(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 0), Pt(0, 1)}).Simplify(2), Pol([]Point[int]{Pt(0, 0)}))
		AssertPolygon(t, Pol([]Point[int]{Pt(1, 1), Pt(1, 1), Pt(1, 1)}).Simplify(0), Pol([]Point[int]{Pt(1, 1)}))
	})
	t.Run("a negative tolerance is taken absolute", func(t *testing.T) {
		bumped := Pol([]Point[int]{Pt(0, 0), Pt(5, 1), Pt(10, 0), Pt(10, 10), Pt(0, 10)})

		AssertPolygon(t, bumped.Simplify(-1), bumped.Simplify(1))
	})
	t.Run("float drops within Epsilon at zero tolerance", func(t *testing.T) {
		nearly := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(1.0, Delta/2), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)})

		AssertPolygon(t, nearly.Simplify(0), Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.0, 0.0), Pt(2.0, 2.0), Pt(0.0, 2.0)}))
	})
	t.Run("keeps the vertex a curve strays farthest by, and drops those within the edges to it", func(t *testing.T) {
		arc := Pol([]Point[int]{Pt(0, 0), Pt(4, 1), Pt(8, 2), Pt(12, 2), Pt(16, 1), Pt(20, 0), Pt(20, 10), Pt(0, 10)})

		AssertPolygon(t, arc.Simplify(1.5), Pol([]Point[int]{Pt(0, 0), Pt(8, 2), Pt(20, 0), Pt(20, 10), Pt(0, 10)}))
		AssertPolygon(t, arc.Simplify(2), Pol([]Point[int]{Pt(0, 0), Pt(20, 0), Pt(20, 10), Pt(0, 10)}))
	})
	t.Run("keeps the polygon's order, from the first kept vertex", func(t *testing.T) {
		turned := Pol([]Point[int]{Pt(2, 2), Pt(0, 2), Pt(0, 1), Pt(0, 0), Pt(2, 0)})

		AssertPolygon(t, turned.Simplify(0), Pol([]Point[int]{Pt(2, 2), Pt(0, 2), Pt(0, 0), Pt(2, 0)}))
	})
	t.Run("an empty polygon is returned as it is", func(t *testing.T) {
		assert.True(t, Polygon[int]{}.Simplify(1).IsZero())
	})
	t.Run("allocates once", func(t *testing.T) {
		square := Pol([]Point[int]{Pt(0, 0), Pt(1, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)})

		AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkPoints = square.Simplify(0).Points
		}), 1)
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, p := range outlineFixtures() {
			for _, tolerance := range []float64{0, 0.5, 2} {
				simplified := p.Simplify(tolerance)
				message := fmt.Sprintf("%s ~%v → %s: ", p, tolerance, simplified)

				assert.True(t, simplified.Simplify(tolerance).Equal(simplified), message)
				for _, vertex := range p.Points {
					assert.True(t, slices.ContainsFunc(slices.Collect(simplified.Edges()), func(edge Segment[float64]) bool {
						return lessOrEqualSquared[float64](edge.DistanceSquaredTo(vertex), tolerance)
					}), message+vertex.String()+" within the tolerance: ")
				}
				assert.Equal(t, simplified.IsEmpty(), p.IsEmpty(), message)

				i := 0
				for _, vertex := range simplified.Points {
					for i < len(p.Points) && p.Points[i] != vertex {
						i++
					}

					assert.True(t, i < len(p.Points), message+vertex.String()+" in order: ")
					i++
				}
			}

			AssertNumber(t, p.Simplify(0).Area(), p.Area(), fmt.Sprintf("%s: ", p))
		}
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

func TestPolygon_Nearest(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("the foot on the nearest edge", func(t *testing.T) {
		AssertPoint(t, square.Nearest(Pt(5, 1)), Pt(2, 1))
		AssertPoint(t, square.Nearest(Pt(5, 6)), Pt(2, 2))
	})
	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		AssertPoint(t, square.Nearest(Pt(1, 1)), Pt(1, 1))
	})
	t.Run("an empty polygon returns the zero point", func(t *testing.T) {
		AssertPoint(t, Polygon[int]{}.Nearest(Pt(3, 4)), Pt(0, 0))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, p := range pointFixtures {
				assertNearest[float64](t, polygon, p)
			}
		}
	})
}

func TestPolygon_EnclosesCircle(t *testing.T) {
	notched := Pol(notchedVertices())

	t.Run("touching the sides and a reflex vertex from inside counts", func(t *testing.T) {
		assert.True(t, notched.EnclosesCircle(Circ(Pt(4, 1), 1)))
		assert.True(t, notched.EnclosesCircle(Circ(Pt(1, 4), 1)))
	})
	t.Run("reaching into the notch or past a side", func(t *testing.T) {
		assert.False(t, notched.Float().EnclosesCircle(Circ(Pt(4.0, 1.0), 1.1)))
		assert.False(t, notched.EnclosesCircle(Circ(Pt(4, 1), 2)))
	})
	t.Run("an empty polygon encloses nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).EnclosesCircle(Circ(Pt(0, 0), 0)))
	})
}

func TestPolygon_EnclosesSegment(t *testing.T) {
	backward := notchedVertices()
	slices.Reverse(backward)
	notched, reversed := Pol(notchedVertices()), Pol(backward)

	t.Run("inside and along an edge", func(t *testing.T) {
		assert.True(t, notched.EnclosesSegment(Seg(Pt(1, 1), Pt(7, 1))))
		assert.True(t, notched.EnclosesSegment(Seg(Pt(0, 0), Pt(8, 0))))
	})
	t.Run("an endpoint in the notch", func(t *testing.T) {
		assert.False(t, notched.EnclosesSegment(Seg(Pt(1, 7), Pt(4, 7))))
	})
	t.Run("crossing the edges of the notch", func(t *testing.T) {
		assert.False(t, notched.EnclosesSegment(Seg(Pt(1, 5), Pt(7, 5))))
	})
	t.Run("leaving through two reflex vertices without crossing an edge", func(t *testing.T) {
		assert.False(t, notched.EnclosesSegment(Seg(Pt(1, 4), Pt(7, 4))))
		assert.False(t, reversed.EnclosesSegment(Seg(Pt(1, 4), Pt(7, 4))))
	})
	t.Run("a chord across the notch, touching an edge at each end", func(t *testing.T) {
		assert.False(t, notched.EnclosesSegment(Seg(Pt(3, 3), Pt(5, 3))))
		assert.False(t, reversed.EnclosesSegment(Seg(Pt(3, 3), Pt(5, 3))))
	})
	t.Run("passing a reflex vertex without leaving", func(t *testing.T) {
		assert.True(t, notched.EnclosesSegment(Seg(Pt(1, 2), Pt(7, 2))))
		assert.True(t, reversed.EnclosesSegment(Seg(Pt(1, 2), Pt(7, 2))))
	})
	t.Run("ending on a reflex vertex", func(t *testing.T) {
		assert.True(t, notched.EnclosesSegment(Seg(Pt(1, 4), Pt(2, 4))))
		assert.True(t, notched.EnclosesSegment(Seg(Pt(4, 0), Pt(4, 2))))
	})
	t.Run("an empty polygon encloses nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).EnclosesSegment(Seg(Pt(0, 0), Pt(0, 0))))
	})
	t.Run("every point of an enclosed segment is contained", func(t *testing.T) {
		for _, p := range append(outlineFixtures(), notched.Float()) {
			for _, s := range segmentFixtures {
				if !p.EnclosesSegment(s) {
					continue
				}

				for i := range 9 {
					assert.True(t, p.Contains(s.PointAt(float64(i)/8)), fmt.Sprintf("%s → %s: ", p, s))
				}
			}
		}
	})
}

func TestPolygon_EnclosesPolygon(t *testing.T) {
	notched := Pol(notchedVertices())

	t.Run("inside and across the notch", func(t *testing.T) {
		assert.True(t, notched.EnclosesPolygon(Pol([]Point[int]{Pt(1, 1), Pt(7, 1), Pt(4, 0)})))
		assert.False(t, notched.EnclosesPolygon(Pol([]Point[int]{Pt(1, 3), Pt(7, 3), Pt(4, 1)})))
	})
	t.Run("an empty polygon on either side encloses nothing", func(t *testing.T) {
		assert.False(t, notched.EnclosesPolygon(Pol[int](nil)))
		assert.False(t, Pol[int](nil).EnclosesPolygon(notched))
	})
}

func TestPolygon_EnclosesRectangle(t *testing.T) {
	notched := Pol(notchedVertices())

	t.Run("an edge along the bottom vertex of the notch stays inside", func(t *testing.T) {
		assert.True(t, notched.EnclosesRectangle(Rect(Pt(4, 1), Sz(6, 2))))
		assert.False(t, notched.EnclosesRectangle(Rect(Pt(4, 2), Sz(6, 2))))
	})
}

func TestPolygon_EnclosesRegularPolygon(t *testing.T) {
	notched := Pol(notchedVertices())

	t.Run("beside the notch and reaching into it", func(t *testing.T) {
		assert.True(t, notched.EnclosesRegularPolygon(RegPol(Pt(2, 2), Sz(1, 1), 4, 0)))
		assert.False(t, notched.EnclosesRegularPolygon(RegPol(Pt(4, 2), Sz(2, 2), 4, 0)))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, notched.EnclosesRegularPolygon(RegPol(Pt(2, 2), Sz(1, 1), 0, 0)))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		inside := RegPol(Pt(2, 2), Sz(1, 1), 6, Pi/5)

		AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = notched.EnclosesRegularPolygon(inside)
		}), 0)
	})
}

func TestPolygon_EnclosesBox(t *testing.T) {
	notched := Pol(notchedVertices())

	t.Run("below the notch and across it", func(t *testing.T) {
		assert.True(t, notched.EnclosesBox(BoxFromMinMax(Pt(1, 0), Pt(7, 2))))
		assert.False(t, notched.EnclosesBox(BoxFromMinMax(Pt(1, 1), Pt(7, 3))))
	})
}

func TestPolygon_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsPolygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, c := range circleFixtures {
				assert.Equal(t, p.IntersectsCircle(c), c.IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", p, c))
			}
		}
	})
}

func TestPolygon_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsPolygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				assert.Equal(t, p.IntersectsSegment(s), s.IntersectsPolygon(p), fmt.Sprintf("%s → %s: ", p, s))
			}
		}
	})
}

func TestPolygon_IntersectionSegment(t *testing.T) {
	t.Run("matches Segment.IntersectionPolygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, s := range segmentFixtures {
				AssertVertices(t, p.IntersectionSegment(s), s.IntersectionPolygon(p), fmt.Sprintf("%s → %s: ", p, s))
			}
		}
	})
}

func TestPolygon_IntersectsPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, square.IntersectsPolygon(square.Translate(Vec(1, 1))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsPolygon(square.Translate(Vec(3, 0))))
	})
	t.Run("a shared edge counts", func(t *testing.T) {
		assert.True(t, square.IntersectsPolygon(square.Translate(Vec(2, 0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, square.IntersectsPolygon(Pol([]Point[int]{Pt(1, 1), Pt(1, 1)})))
		assert.True(t, Pol([]Point[int]{Pt(1, 1), Pt(1, 1)}).IntersectsPolygon(square))
	})
	t.Run("edges crossing without a vertex inside", func(t *testing.T) {
		cross := Pol([]Point[int]{Pt(-1, 1), Pt(3, 1), Pt(3, 1), Pt(-1, 1)})
		plus := Pol([]Point[int]{Pt(1, -1), Pt(1, 3), Pt(1, 3), Pt(1, -1)})

		assert.True(t, cross.IntersectsPolygon(plus))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, square.IntersectsPolygon(Pol[int](nil)))
		assert.False(t, Pol[int](nil).IntersectsPolygon(square))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range polygonFixtures() {
			for _, b := range polygonFixtures() {
				assert.Equal(t, a.IntersectsPolygon(b), b.IntersectsPolygon(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func BenchmarkPolygon_IntersectsPolygon(b *testing.B) {
	polygon := benchPolygon()
	overlapping := polygon.Translate(Vec(150.0, 0.0))
	apartWithinBounds := polygon.Translate(Vec(150.0, 150.0))
	apart := polygon.Translate(Vec(300.0, 0.0))

	b.Run("overlapping", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.IntersectsPolygon(overlapping)
		}
	})
	b.Run("apart within overlapping bounds", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.IntersectsPolygon(apartWithinBounds)
		}
	})
	b.Run("apart", func(b *testing.B) {
		for b.Loop() {
			sinkBool = polygon.IntersectsPolygon(apart)
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
				assert.Equal(t, p.IntersectsRectangle(r), p.IntersectsPolygon(r.Polygon()), fmt.Sprintf("%s → %s: ", p, r))
			}
		}
	})
}

func TestPolygon_IntersectsRegularPolygon(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, square.IntersectsRegularPolygon(RegPol(Pt(2, 2), Sz(2, 2), 4, 0)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsRegularPolygon(RegPol(Pt(5, 5), Sz(2, 2), 4, 0)))
	})
	t.Run("apart within overlapping bounds", func(t *testing.T) {
		assert.False(t, square.IntersectsRegularPolygon(RegPol(Pt(4, 4), Sz(2, 2), 4, 0)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, square.IntersectsRegularPolygon(RegPol(Pt(1, 1), Sz(20, 20), 6, 0)))
		assert.True(t, square.IntersectsRegularPolygon(RegPol(Pt(1, 1), Sz(1, 1), 3, 0)))
	})
	t.Run("edges crossing without a vertex inside", func(t *testing.T) {
		assert.True(t, square.IntersectsRegularPolygon(RegPol(Pt(1, 1), Sz(1, 20), 4, 0)))
	})
	t.Run("an empty polygon on either side intersects nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).IntersectsRegularPolygon(RegPol(Pt(1, 1), Sz(2, 2), 4, 0)))
		assert.False(t, square.IntersectsRegularPolygon(RegPol(Pt(1, 1), Sz(2, 2), 0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, rp := range regularPolygonFixtures {
				assert.Equal(t, p.IntersectsRegularPolygon(rp), p.IntersectsPolygon(rp.Polygon()), fmt.Sprintf("%s → %s: ", p, rp))
			}
		}
	})
}

func TestPolygon_IntersectsBox(t *testing.T) {
	square := Pol(squareVertices())

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, square.IntersectsBox(BoxFromMinMax(Pt(1, 1), Pt(3, 3))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, square.IntersectsBox(BoxFromMinMax(Pt(3, 3), Pt(5, 5))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, square.IntersectsBox(BoxFromMinMax(Pt(-9, -9), Pt(11, 11))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, Pol[int](nil).IntersectsBox(BoxFromMinMax(Pt(0, 0), Pt(2, 2))))
	})
	t.Run("matches the box as a polygon", func(t *testing.T) {
		for _, p := range polygonFixtures() {
			for _, b := range boxFixtures {
				assert.Equal(t, p.IntersectsBox(b), p.IntersectsPolygon(b.Rectangle().Polygon()), fmt.Sprintf("%s → %s: ", p, b))
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

func TestPolygon_IsEmpty(t *testing.T) {
	t.Run("no vertices", func(t *testing.T) {
		assert.True(t, Polygon[int]{}.IsEmpty())
		assert.True(t, Polygon[int]{[]Point[int]{}}.IsEmpty())
	})
	t.Run("with vertices", func(t *testing.T) {
		assert.False(t, Pol(squareVertices()).IsEmpty())
		assert.False(t, Pol(triangleVertices()).IsEmpty())
	})
}

func TestPolygon_IsConvex(t *testing.T) {
	t.Run("a triangle and a square in either winding", func(t *testing.T) {
		assert.True(t, Pol(triangleVertices()).IsConvex())
		assert.True(t, Pol(squareVertices()).IsConvex())
		assert.True(t, Pol([]Point[int]{Pt(0, 2), Pt(2, 2), Pt(2, 0), Pt(0, 0)}).IsConvex())
	})
	t.Run("a concavity", func(t *testing.T) {
		assert.False(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 2), Pt(0, 4), Pt(1, 2)}).IsConvex())
	})
	t.Run("a star turning one way throughout goes around twice", func(t *testing.T) {
		pentagon := slices.Collect(RegPol(Pt(0.0, 0.0), SzU(10.0), 5, 0).Vertices())
		star := []Point[float64]{pentagon[0], pentagon[2], pentagon[4], pentagon[1], pentagon[3]}

		assert.False(t, Pol(star).IsConvex())
	})
	t.Run("a repeated vertex and a vertex on an edge are allowed", func(t *testing.T) {
		assert.True(t, Pol([]Point[int]{Pt(0, 0), Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}).IsConvex())
		assert.True(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}).IsConvex())
	})
	t.Run("an edge doubling back is not", func(t *testing.T) {
		assert.False(t, Pol([]Point[int]{Pt(0, 0), Pt(4, 0), Pt(4, 4), Pt(4, 2), Pt(0, 2)}).IsConvex())
	})
	t.Run("no area is not convex", func(t *testing.T) {
		assert.False(t, Polygon[int]{}.IsConvex())
		assert.False(t, Pol([]Point[int]{Pt(1, 1)}).IsConvex())
		assert.False(t, Pol([]Point[int]{Pt(0, 0), Pt(2, 2)}).IsConvex())
		assert.False(t, Pol([]Point[int]{Pt(0, 0), Pt(1, 1), Pt(2, 2)}).IsConvex())
	})
	t.Run("a NaN vertex is not convex", func(t *testing.T) {
		assert.False(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(math.NaN(), 0.0), Pt(0.0, 2.0)}).IsConvex())
	})
	t.Run("allocates nothing", func(t *testing.T) {
		square := Pol(squareVertices())

		AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = square.IsConvex()
		}), 0)
	})
	t.Run("every rectangle and regular polygon with an area is convex", func(t *testing.T) {
		for _, r := range rectFixtures {
			assert.Equal(t, r.Polygon().IsConvex(), r.Area() > 0, r.String()+": ")
		}
		for _, rp := range regularPolygonFixtures {
			assert.True(t, rp.Polygon().IsConvex(), rp.String()+": ")
		}
	})
}

func TestPolygon_Cast(t *testing.T) {
	p := Pol([]Point[float64]{Pt(1.5, -2.5), Pt(3.5, 4.5), Pt(-1.5, 2.5)})

	t.Run("matches Int and Float", func(t *testing.T) {
		AssertPolygon(t, p.Cast[int](), p.Int())
		AssertPolygon(t, p.Cast[float64](), p.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		AssertPolygon(t, p.Cast[int8](), Pol([]Point[int8]{Pt[int8](2, -3), Pt[int8](4, 5), Pt[int8](-2, 3)}))
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

				assert.True(t, moved.Centroid().Equal(polygon.Centroid().Add(vector)), fmt.Sprintf("%s → %s: ", polygon, vector))
				for i, vertex := range moved.Points {
					assert.True(t, vertex.Equal(polygon.Points[i].Add(vector)), fmt.Sprintf("%s → %s: ", polygon, vector))
				}
			}
		}
	})
	t.Run("move to centers where asked", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, point := range pointFixtures {
				assert.True(t, polygon.MoveTo(point).Centroid().Equal(point), fmt.Sprintf("%s → %s: ", polygon, point))
			}
		}
	})
	t.Run("scale keeps the centroid", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				scaled := polygon.Scale(factor)

				assert.True(t, scaled.Centroid().Equal(polygon.Centroid()), fmt.Sprintf("%s ×%v: ", polygon, factor))
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
			if len(polygon.Points) != 3 {
				continue
			}

			var sum Point[float64]
			for _, vertex := range polygon.Points {
				sum = sum.AddXY(vertex.XY())
			}

			expected := sum.Divide(3)
			assert.True(t, polygon.Centroid().Equal(expected), fmt.Sprintf("%s: ", polygon))
		}
	})
	t.Run("centroid lies within the bounds", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			assert.True(t, polygon.Bounds().Contains(polygon.Centroid()), fmt.Sprintf("%s: ", polygon))
		}
	})
}

func TestPolygon_Immutable(t *testing.T) {
	p := Pol([]Point[int]{Pt(0, 0), Pt(2, 0)})

	p.Translate(Vec(1, -1))
	p.MoveTo(Pt(10, 10))
	p.Scale(2)
	p.ScaleXY(2, 3)
	p.Lerp(Pol([]Point[int]{Pt(4, 4), Pt(6, 4)}), 0.5)
	p.ConvexHull()

	AssertVertices(t, p.Points, []Point[int]{Pt(0, 0), Pt(2, 0)})
}

// squareVertices and triangleVertices build fresh slices, so a test that mutates one
// cannot leak into the next.
func squareVertices() []Point[int] {
	return []Point[int]{Pt(0, 0), Pt(2, 0), Pt(2, 2), Pt(0, 2)}
}

// notchedVertices is a square with a notch cut in from the top that widens below into a diamond,
// whose side corners are reflex vertices on one line with the outside between them.
func notchedVertices() []Point[int] {
	return []Point[int]{
		Pt(0, 0), Pt(8, 0), Pt(8, 8), Pt(5, 8), Pt(5, 6), Pt(6, 4), Pt(4, 2), Pt(2, 4), Pt(3, 6), Pt(3, 8), Pt(0, 8),
	}
}

// outlineFixtures returns the polygon fixtures with the polygons of the rectangle and the
// regular polygon fixtures, every outline the polygon methods are checked over.
func outlineFixtures() []Polygon[float64] {
	polygons := polygonFixtures()
	for _, r := range rectFixtures {
		polygons = append(polygons, r.Polygon())
	}
	for _, rp := range regularPolygonFixtures {
		polygons = append(polygons, rp.Polygon())
	}

	return polygons
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

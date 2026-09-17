package geom

import (
	"encoding/json"
	"fmt"
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
	t.Run("empty is the zero point", func(t *testing.T) {
		AssertPoint(t, Pol([]Point[int]{}).Center(), Pt(0, 0))
		AssertPoint(t, Polygon[float64]{}.Center(), Pt(0.0, 0.0))
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

func TestPolygon_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRect(t, Pol(squareVertices()).Bounds(), RectangleFromMinMax(Pt(0, 0), Pt(2, 2)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Pol(triangleVertices()).Bounds(), RectangleFromMinMax(Pt(0.0, 0.0), Pt(2.5, 1.0)))
	})
	t.Run("vertex order does not matter", func(t *testing.T) {
		AssertRect(t, Pol([]Point[int]{Pt(3, -1), Pt(-2, 4), Pt(0, 0)}).Bounds(), RectangleFromMinMax(Pt(-2, -1), Pt(3, 4)))
	})
	t.Run("empty is the zero rectangle", func(t *testing.T) {
		AssertRect(t, Polygon[int]{}.Bounds(), Rectangle[int]{})
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
		assert.Equal(t, Pol(squareVertices()).String(), "Pol((0,0), (2,0), (2,2), (0,2))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Pol(triangleVertices()).String(), "Pol((0.00,0.00), (2.50,0.50), (2.00,1.00))")
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
	t.Run("centroid is the average of the vertices", func(t *testing.T) {
		for _, polygon := range polygonFixtures() {
			var sum Point[float64]
			for _, vertex := range polygon.Vertices {
				sum = sum.AddXY(vertex.XY())
			}

			expected := sum.Divide(float64(len(polygon.Vertices)))
			assert.True(t, polygon.Center().Equal(expected), fmt.Sprintf("%s: ", polygon))
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
	// Output: Pol((0,0), (2,0), (2,2))
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

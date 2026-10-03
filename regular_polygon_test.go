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

func TestRegularPolygon_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0), RegularPolygon[int]{Center: Pt(0, 0), Size: Sz(2, 2), N: 4, Angle: 0})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, Pi/3, 0), RegularPolygon[float64]{Center: Pt(0.5, -1.25), Size: Sz(2.5, 3.75), N: 6, Angle: Pi / 3})
	})
	t.Run("a negative size is taken absolute", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0, 0), Sz(-2, 3), 4, 0, 0), RegPol(Pt(0, 0), Sz(2, 3), 4, 0, 0))
		geomtest.AssertRegularPolygon(t, Square(Pt(0.0, 0.0), Sz(-2.0, -3.0), OrientationPointyTop), Square(Pt(0.0, 0.0), Sz(2.0, 3.0), OrientationPointyTop))
	})
}

func TestRegularPolygonWithOrientation(t *testing.T) {
	center, size := Pt(0, 0), Sz(10, 10)
	pointy := RegularPolygonWithOrientation(center, size, 6, OrientationPointyTop)
	flat := RegularPolygonWithOrientation(center, size, 6, OrientationFlatTop)

	t.Run("places the first vertex by the orientation phase and does not turn", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, pointy, RegPol(center, size, 6, 0, RegularPolygonOrientationPhase(6, OrientationPointyTop)))
		geomtest.AssertRegularPolygon(t, flat, RegPol(center, size, 6, 0, RegularPolygonOrientationPhase(6, OrientationFlatTop)))

		assert.NotEqual(t, pointy.Phase, flat.Phase)
	})
	t.Run("pointy top has a vertex at the top and flat top an edge", func(t *testing.T) {
		assertOrientation(t, pointy.Float(), OrientationPointyTop, "")
		assertOrientation(t, flat.Float(), OrientationFlatTop, "")
	})
	t.Run("unequal semi-axes are stored as given", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, Hexagon(Pt(0.0, 0.0), Sz(20.0, 10.0), OrientationPointyTop), RegPol(Pt(0.0, 0.0), Sz(20.0, 10.0), 6, 0, 3*Pi/2))
	})
	t.Run("unequal semi-axes span Width across and Height up", func(t *testing.T) {
		for _, orientation := range Orientations() {
			for n := 3; n <= 9; n++ {
				polygon := RegularPolygonWithOrientation(Pt(0.0, 0.0), Sz(20.0, 10.0), n, orientation)
				message := fmt.Sprintf("%s n=%d: ", orientation, n)

				assertOrientation(t, polygon, orientation, message)
				assertSteppedFromTop(t, polygon, Sz(20.0, 10.0), orientation, message)
			}
		}
	})
	t.Run("a flat-top square of unequal semi-axes is the rectangle across its edges", func(t *testing.T) {
		square := Square(Pt(0.0, 0.0), Sz(20.0, 10.0), OrientationFlatTop)

		assertSteppedFromTop(t, square, Sz(20.0, 10.0), OrientationFlatTop, "")
		geomtest.AssertBox(t, square.Bounds(), Bx(Pt(-10*Sqrt2, -5*Sqrt2), Pt(10*Sqrt2, 5*Sqrt2)))
	})
}

// assertOrientation checks the top of a polygon an orientation constructor made: a single
// vertex above the center for OrientationPointyTop, and two vertices level with each other
// either side of it for OrientationFlatTop.
func assertOrientation(t *testing.T, polygon RegularPolygon[float64], orientation Orientation, message string) {
	t.Helper()

	top := math.Inf(1)
	for vertex := range polygon.Vertices() {
		top = min(top, vertex.Y)
	}

	var highest []Point[float64]
	for vertex := range polygon.Vertices() {
		if Equal(vertex.Y, top) {
			highest = append(highest, vertex)
		}
	}

	if orientation == OrientationPointyTop {
		if assert.Equal(t, len(highest), 1, message+"one vertex at the top: ") {
			geomtest.AssertNumber(t, highest[0].X, polygon.Center.X, message)
		}

		return
	}

	if assert.Equal(t, len(highest), 2, message+"an edge at the top: ") {
		geomtest.AssertNumber(t, highest[0].Midpoint(highest[1]).X, polygon.Center.X, message)
	}
}

// assertSteppedFromTop checks that the polygon has the vertices of the polygon stepped around
// the axis-aligned ellipse of the size from the angle that places its top, a vertex there for
// OrientationPointyTop and half a step before it for OrientationFlatTop, in whatever order.
func assertSteppedFromTop(t *testing.T, polygon RegularPolygon[float64], size Size[float64], orientation Orientation, message string) {
	t.Helper()

	start := 3 * Pi / 2
	if orientation == OrientationFlatTop {
		start -= Pi / float64(polygon.N)
	}

	vertices := slices.Collect(polygon.Vertices())
	assert.Equal(t, len(vertices), polygon.N, message)

	for i := range polygon.N {
		expected := polygon.Center.Add(VectorFromAngleSize(start+float64(i)*2*Pi/float64(polygon.N), size))

		assert.True(t, slices.ContainsFunc(vertices, expected.Equal), message+expected.String()+": ")
	}
}

func TestTriangle(t *testing.T) {
	triangle := Triangle(Pt(1, -1), Sz(3, 3), OrientationPointyTop)

	t.Run("has three sides at the orientation angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, triangle, RegPol(Pt(1, -1), Sz(3, 3), 3, 0, RegularPolygonOrientationPhase(3, OrientationPointyTop)))
	})
	t.Run("integer vertices round after scaling", func(t *testing.T) {
		// vertices 1 and 2 land within one unit of the exact (3.598, 0.5) and (-1.598, 0.5);
		// their Y of 0.5 falls just below the .5 tie in float64 and rounds down to 0
		geomtest.AssertVertices(t, slices.Collect(triangle.Vertices()), []Point[int]{Pt(1, -4), Pt(4, 0), Pt(-2, 0)})
	})
}

func TestSquare(t *testing.T) {
	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), OrientationPointyTop)

	t.Run("has four sides at the orientation angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, square, RegPol(Pt(50.0, 50.0), Sz(100.0, 100.0), 4, 0, RegularPolygonOrientationPhase(4, OrientationPointyTop)))
	})
	t.Run("pointy top starts at the visual top and winds to the right", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(square.Vertices()), []Point[float64]{Pt(50.0, -50.0), Pt(150.0, 50.0), Pt(50.0, 150.0), Pt(-50.0, 50.0)})
	})
}

func TestHexagon(t *testing.T) {
	// float64 avoids the int-rounding collapse where sin(±π/6) ≈ 0.4999 would truncate to 0
	hexagon := Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), OrientationPointyTop)

	t.Run("has six sides at the orientation angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hexagon, RegPol(Pt(0.0, 0.0), Sz(10.0, 10.0), 6, 0, RegularPolygonOrientationPhase(6, OrientationPointyTop)))
	})
	t.Run("pointy top has vertices at top and bottom and flat sides left and right", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(hexagon.Vertices()), []Point[float64]{
			Pt(0.0, -10.0),
			Pt(5*Sqrt3, -5.0),
			Pt(5*Sqrt3, 5.0),
			Pt(0.0, 10.0),
			Pt(-5*Sqrt3, 5.0),
			Pt(-5*Sqrt3, -5.0),
		})
	})
}

func TestRegularPolygon_Anchor(t *testing.T) {
	t.Run("a flat-top hexagon anchors the top edge midpoint", func(t *testing.T) {
		hex := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)

		geomtest.AssertPoint(t, hex.Anchor(Top), Pt(0.0, -10*math.Cos(Pi/6)))
		geomtest.AssertPoint(t, hex.Anchor(Bottom), Pt(0.0, 10*math.Cos(Pi/6)))
		geomtest.AssertPoint(t, hex.Anchor(Right), Pt(10.0, 0.0))
	})
	t.Run("a pointy-top hexagon anchors the top vertex", func(t *testing.T) {
		hex := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationPointyTop)

		geomtest.AssertPoint(t, hex.Anchor(Top), Pt(0.0, -10.0))
		geomtest.AssertPoint(t, hex.Anchor(Right), Pt(10*math.Cos(Pi/6), 0.0))
	})
	t.Run("a diagonal leaves through the edge on the diagonal", func(t *testing.T) {
		diamond := RegPol(Pt(0.0, 0.0), Sz(10.0, 4.0), 4, 0, 0)
		reach := 1 / (1/10.0 + 1/4.0)

		geomtest.AssertPoint(t, diamond.Anchor(DirectionDownRight), Pt(reach, reach))
		geomtest.AssertPoint(t, diamond.Anchor(DirectionUpLeft), Pt(-reach, -reach))
	})
	t.Run("the direction is taken in the world, so a turn moves the anchor along the boundary", func(t *testing.T) {
		square := RegPol(Pt(0.0, 0.0), SzU(10.0), 4, 0, 0)

		geomtest.AssertPoint(t, square.Anchor(Right), Pt(10.0, 0.0))
		geomtest.AssertPoint(t, square.Rotate(Pi/4).Anchor(Right), Pt(10*OneOverSqrt2, 0.0))
		geomtest.AssertPoint(t, square.Rotate(Pi/2).Anchor(Right), Pt(10.0, 0.0))
	})
	t.Run("int lands on the rounded edge", func(t *testing.T) {
		hex := Hexagon(Pt(0, 0), SzU(10), OrientationPointyTop)

		geomtest.AssertPoint(t, hex.Anchor(Top), Pt(0, -10))
		assert.True(t, hex.Contains(hex.Anchor(DirectionDownRight)))
	})
	t.Run("none is the center", func(t *testing.T) {
		geomtest.AssertPoint(t, Hexagon(Pt(1.0, 2.0), SzU(10.0), OrientationFlatTop).Anchor(DirectionNone), Pt(1.0, 2.0))
	})
	t.Run("fewer than three vertices anchor at the center", func(t *testing.T) {
		geomtest.AssertPoint(t, RegPol(Pt(1.0, 2.0), SzU(10.0), 2, 0, 0).Anchor(Right), Pt(1.0, 2.0))
		geomtest.AssertPoint(t, RegPol(Pt(1.0, 2.0), SzU(10.0), 0, 0, 0).Anchor(Right), Pt(1.0, 2.0))
	})
	t.Run("a semi-axis along the ray anchors at the center", func(t *testing.T) {
		geomtest.AssertPoint(t, RegPol(Pt(1.0, 2.0), Sz(10.0, 0.0), 4, 0, 0).Anchor(Right), Pt(1.0, 2.0))
		geomtest.AssertPoint(t, RegPol(Pt(1.0, 2.0), Sz(10.0, 0.0), 4, 0, 0).Anchor(Top), Pt(1.0, 2.0))
	})
}

func TestRegularPolygon_Vertices(t *testing.T) {
	t.Run("fewer than one side yields nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0).Vertices()))
		assert.Nil(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0, 0).Vertices()))
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 4, 0, 0).Vertices()), []Point[int]{
			Pt(1, 0),
			Pt(0, 1),
			Pt(-1, 0),
			Pt(0, -1),
		})
	})
	t.Run("non-square size stretches the ellipse", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(RegPol(Pt(0, 0), Sz(2, 3), 4, 0, 0).Vertices()), []Point[int]{
			Pt(2, 0),
			Pt(0, 3),
			Pt(-2, 0),
			Pt(0, -3),
		})
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0, 0).Vertices()), []Point[float64]{
			Pt(2.0, 0.0),
			Pt(1.0, 1.5*Sqrt3),
			Pt(-1.0, 1.5*Sqrt3),
			Pt(-2.0, 0.0),
			Pt(-1.0, -1.5*Sqrt3),
			Pt(1.0, -1.5*Sqrt3),
		})
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for vertex := range RegPol(Pt(0, 0), Sz(1, 1), 4, 0, 0).Vertices() {
			geomtest.AssertPoint(t, vertex, Pt(1, 0))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
				for vertex := range rp.Vertices() {
					sinkBool = vertex.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", rp))
		}
	})
}

func TestRegularPolygon_Edges(t *testing.T) {
	t.Run("closes back to the first vertex", func(t *testing.T) {
		edges := slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 4, 0, 0).Edges())

		assert.Length(t, edges, 4)
		geomtest.AssertSegment(t, edges[0], Seg(Pt(1, 0), Pt(0, 1)))
		geomtest.AssertSegment(t, edges[3], Seg(Pt(0, -1), Pt(1, 0)))
	})
	t.Run("single vertex is one zero-length edge", func(t *testing.T) {
		edges := slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 1, 0, 0).Edges())

		assert.Length(t, edges, 1)
		geomtest.AssertSegment(t, edges[0], Seg(Pt(1.0, 0.0), Pt(1.0, 0.0)))
	})
	t.Run("fewer than one side yields nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0).Edges()))
		assert.Nil(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0, 0).Edges()))
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for edge := range RegPol(Pt(0, 0), Sz(1, 1), 4, 0, 0).Edges() {
			geomtest.AssertSegment(t, edge, Seg(Pt(1, 0), Pt(0, 1)))

			break
		}
	})
	t.Run("are the edges of the polygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			edges, expected := slices.Collect(rp.Edges()), slices.Collect(rp.Polygon().Edges())

			assert.Length(t, edges, len(expected), fmt.Sprintf("%s: ", rp))
			for i := range edges {
				geomtest.AssertSegment(t, edges[i], expected[i], fmt.Sprintf("%s #%d: ", rp, i))
			}
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
				for edge := range rp.Edges() {
					sinkBool = edge.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", rp))
		}
	})
}

func TestRegularPolygon_Centroid(t *testing.T) {
	geomtest.AssertPoint(t, Hexagon(Pt(1, 2), SzU(10), OrientationFlatTop).Centroid(), Pt(1, 2))
}

func TestRegularPolygon_Area(t *testing.T) {
	t.Run("agrees with the polygon", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(2.0), OrientationFlatTop)

		geomtest.AssertNumber(t, hexagon.Area(), hexagon.Polygon().Area())
		geomtest.AssertNumber(t, hexagon.Area(), 3*Sqrt3/2*4) // 3√3/2 · r²
	})
	t.Run("a square of semi-axis r encloses 2r²", func(t *testing.T) {
		geomtest.AssertNumber(t, Square(Pt(0.0, 0.0), SzU(3.0), OrientationPointyTop).Area(), 18.0)
	})
	t.Run("an ellipse scales the area by both semi-axes", func(t *testing.T) {
		geomtest.AssertNumber(t, Square(Pt(0.0, 0.0), Sz(2.0, 5.0), OrientationPointyTop).Area(), 20.0)
	})
	t.Run("fewer than three vertices enclose nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 0, 0, 0).Area(), 0.0)
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 2, 0, 0).Area(), 0.0)
	})
	t.Run("agrees with the polygon at any angle and size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertNumber(t, rp.Area(), rp.Polygon().Area(), rp.String())
		}
	})
}

func TestRegularPolygon_Perimeter(t *testing.T) {
	t.Run("agrees with the polygon", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(2.0), OrientationFlatTop)

		geomtest.AssertNumber(t, hexagon.Perimeter(), hexagon.Polygon().Perimeter())
		geomtest.AssertNumber(t, hexagon.Perimeter(), 12.0) // six edges of length r
	})
	t.Run("a square of semi-axis r has edges of r√2", func(t *testing.T) {
		geomtest.AssertNumber(t, Square(Pt(0.0, 0.0), SzU(3.0), OrientationPointyTop).Perimeter(), 12*Sqrt2)
	})
	t.Run("an ellipse sums chords of different lengths", func(t *testing.T) {
		square := Square(Pt(0.0, 0.0), Sz(2.0, 5.0), OrientationPointyTop)

		geomtest.AssertNumber(t, square.Perimeter(), 4*math.Hypot(2, 5))
		geomtest.AssertNumber(t, square.Perimeter(), square.Polygon().Perimeter())
	})
	t.Run("fewer than two vertices have no edge", func(t *testing.T) {
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 0, 0, 0).Perimeter(), 0.0)
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 1, 0, 0).Perimeter(), 0.0)
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 2, 0, 0).Perimeter(), 40.0)
	})
	t.Run("int measures the exact polygon, not the rounded vertices", func(t *testing.T) {
		hexagon := Hexagon(Pt(0, 0), SzU(1), OrientationPointyTop)

		geomtest.AssertNumber(t, hexagon.Perimeter(), 6.0)
		assert.NotEqual(t, hexagon.Polygon().Perimeter(), 6.0)
	})
	t.Run("agrees with the polygon at any angle and size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertNumber(t, rp.Perimeter(), rp.Polygon().Perimeter(), rp.String())
		}
	})
}

func TestRegularPolygon_Inertia(t *testing.T) {
	t.Run("agrees with the polygon", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(2.0), OrientationFlatTop)

		geomtest.AssertNumber(t, hexagon.Inertia(), hexagon.Polygon().Inertia())
		geomtest.AssertNumber(t, hexagon.Inertia(), 5*Sqrt3/8*16) // 5√3/8 · r⁴
	})
	t.Run("a square of semi-axis r has the moment 2r⁴/3", func(t *testing.T) {
		geomtest.AssertNumber(t, Square(Pt(0.0, 0.0), SzU(3.0), OrientationPointyTop).Inertia(), 54.0)
	})
	t.Run("an ellipse scales each axis by the cube of one semi-axis", func(t *testing.T) {
		square := Square(Pt(0.0, 0.0), Sz(2.0, 5.0), OrientationPointyTop)

		geomtest.AssertNumber(t, square.Inertia(), square.Polygon().Inertia())
		geomtest.AssertNumber(t, square.Inertia(), 4*10*(16+100)/48.0) // a rhombus of diagonals p, q: pq(p²+q²)/48
	})
	t.Run("fewer than three vertices enclose nothing", func(t *testing.T) {
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 0, 0, 0).Inertia(), 0.0)
		geomtest.AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 2, 0, 0).Inertia(), 0.0)
	})
	t.Run("agrees with the polygon at any angle and size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertNumber(t, rp.Inertia(), rp.Polygon().Inertia(), rp.String())
		}
	})
}

func TestRegularPolygon_Bounds(t *testing.T) {
	t.Run("vertices on the axes", func(t *testing.T) {
		geomtest.AssertBox(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Bounds(), BoxFromMinMax(Pt(-1, 0), Pt(3, 4)))
	})
	t.Run("hexagon is tight around its vertices", func(t *testing.T) {
		// width = 2r, height = √3 r
		bounds := Hexagon(Pt(0.0, 0.0), Sz(2.0, 2.0), OrientationFlatTop).Bounds()

		geomtest.AssertNumber(t, bounds.Width(), 4.0)
		geomtest.AssertNumber(t, bounds.Height(), 2.0*Sqrt3)
	})
	t.Run("one and two vertices reach only where they lie", func(t *testing.T) {
		geomtest.AssertBox(t, RegPol(Pt(0.0, 0.0), SzU(2.0), 1, Pi/2, 0).Bounds(), BoxFromMinMax(Pt(0.0, 2.0), Pt(0.0, 2.0)))
		geomtest.AssertBox(t, RegPol(Pt(0.0, 0.0), SzU(2.0), 2, Pi/4, 0).Bounds(), BoxFromMinMax(Pt(-Sqrt2, -Sqrt2), Pt(Sqrt2, Sqrt2)))
	})
	t.Run("no vertices is the zero rectangle", func(t *testing.T) {
		geomtest.AssertBox(t, RegPol(Pt(3, 4), Sz(10, 10), 0, 0, 0).Bounds(), Box[int]{})
	})
	t.Run("is exactly the box around the vertices, rounded alike for int", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, angle := range []float64{0, 0.3, Pi / 3, 2, -1} {
				turned := rp.Rotate(angle)

				geomtest.AssertBox(t, turned.Bounds(), turned.Polygon().Bounds(), turned.String())
				assert.Equal(t, turned.Int().Bounds(), turned.Int().Polygon().Bounds(), turned.String())
			}
		}
	})
	t.Run("two vertices reaching equally far are both read, for int whichever rounds farther", func(t *testing.T) {
		rp := RegPol(Pt(0, 8), Sz(3, 5), 3, 3*Pi/2, 0)

		geomtest.AssertBox(t, rp.Bounds(), Bx(Pt(-4, 5), Pt(4, 10)))
		assert.True(t, rp.Contains(Pt(-4, 10)))
	})
	t.Run("is exactly the box around the vertices of small int polygons turned and phased in steps", func(t *testing.T) {
		for n := 3; n <= 6; n++ {
			for size := range 25 {
				for turn := range 16 {
					for phase := range 24 {
						rp := RegPol(Pt(0, 8), Sz(size%5+1, size/5+1), n, float64(turn)*Pi/8, float64(phase)*Pi/12)

						assert.Equal(t, rp.Bounds(), rp.Polygon().Bounds(), rp.String())
					}
				}
			}
		}
	})
}

func TestRegularPolygon_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Translate(Vec(1, -2)), RegPol(Pt(2, 0), Sz(2, 2), 4, 0, 0))
	})
	t.Run("an int polygon of equal semi-axes moves every vertex by the vector", func(t *testing.T) {
		for _, polygon := range []RegularPolygon[int]{RegPol(Pt(0, 0), Sz(9, 9), 9, Pi/6, 0), RegPol(Pt(0, 0), Sz(1, 1), 3, 0, 2*Pi/3)} {
			for _, vector := range []Vector[int]{Vec(-1, 7), Vec(-9, -3)} {
				moved := slices.Collect(polygon.Translate(vector).Vertices())
				expected := slices.Collect(polygon.Polygon().Translate(vector).Vertices())

				geomtest.AssertVertices(t, moved, expected, fmt.Sprintf("%s by %s: ", polygon, vector))
			}
		}
	})
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).MoveTo(Pt(-3, 5)), RegPol(Pt(-3, 5), Sz(2, 2), 4, 0, 0))
}

func TestRegularPolygon_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Scale(0.5), RegPol(Pt(1, 2), Sz(1, 1), 4, 0, 0))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).ScaleXY(2, 3), RegPol(Pt(1, 2), Sz(4, 6), 4, 0, 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Scale(-2), RegPol(Pt(1, 2), Sz(4, 4), 4, 0, 0))
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).ScaleXY(-2, 3), RegPol(Pt(1, 2), Sz(4, 6), 4, 0, 0))
	})
}

func TestRegularPolygon_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 4), 4, 0, 0).Unscale(2), RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 6), 4, 0, 0).UnscaleXY(2, 3), RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 4), 4, 0, 0).Unscale(-2), RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			RegPol(Pt(1, 2), Sz(4, 4), 4, 0, 0).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			RegPol(Pt(1, 2), Sz(4, 4), 4, 0, 0).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestRegularPolygon_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Resize(Sz(3, 7)), RegPol(Pt(1, 2), Sz(3, 7), 6, 0, 0))
	})
	t.Run("a negative semi-axis is taken absolute", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Resize(Sz(-3, 7)), RegPol(Pt(1, 2), Sz(3, 7), 6, 0, 0))
	})
	t.Run("keeps the vertex count, the angle and the phase", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			resized := rp.Resize(Sz(3.0, 7.0))

			geomtest.AssertRegularPolygon(t, resized, RegPol(rp.Center, Sz(3.0, 7.0), rp.N, rp.Angle, rp.Phase), rp.String())
			geomtest.AssertRegularPolygon(t, resized.Resize(rp.Size), rp, rp.String())
		}
	})
}

func TestRegularPolygon_Canonical(t *testing.T) {
	t.Run("takes a literal negative size absolute and keeps the rest", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegularPolygon[int]{Pt(1, 2), Sz(-2, 3), 6, Pi / 3, 0}.Canonical(), RegPol(Pt(1, 2), Sz(2, 3), 6, Pi/3, 0))
	})
	t.Run("is the polygon RegPol builds", func(t *testing.T) {
		rp := RegularPolygon[float64]{Pt(0.5, -1.25), Sz(-2.5, -3.75), 5, 1, 0}

		geomtest.AssertRegularPolygon(t, rp.Canonical(), RegPol(rp.Center, rp.Size, rp.N, rp.Angle, 0))
		geomtest.AssertVertices(t, slices.Collect(rp.Canonical().Vertices()), slices.Collect(RegPol(rp.Center, rp.Size, rp.N, rp.Angle, 0).Vertices()))
	})
	t.Run("normalizes the angle the way Rotate stores it", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, 7, 0}.Canonical(), RegPol(Pt(0.0, 0.0), SzU(2.0), 4, 7-2*Pi, 0))
		geomtest.AssertNumber(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, -Pi / 2, 0}.Canonical().Angle, 3*Pi/2)
	})
	t.Run("snaps a residue of turning back to exactly zero, where Rotate does not", func(t *testing.T) {
		drifted := RegPol(Pt(0.0, 0.0), SzU(2.0), 4, 0, 0).Rotate(0.1).Rotate(0.2).Rotate(-0.3)

		assert.True(t, drifted.Angle != 0)
		assert.Equal(t, drifted.Canonical().Angle, 0.0)
		assert.Equal(t, RegularPolygon[int]{Pt(0, 0), SzU(2), 4, -Delta / 2, 0}.Canonical().Angle, 0.0)
		assert.Equal(t, RegularPolygon[int]{Pt(0, 0), SzU(2), 4, 2 * Delta, 0}.Canonical().Angle, 2*Delta)
	})
	t.Run("a well-formed polygon is unchanged", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertRegularPolygon(t, rp.Canonical(), rp, rp.String())
			geomtest.AssertRegularPolygon(t, rp.Rotate(0).Canonical(), rp.Rotate(0), rp.String())
		}
	})
	t.Run("normalizes and snaps the phase like the angle, without reducing it to a step", func(t *testing.T) {
		geomtest.AssertNumber(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, 0, -Pi / 2}.Canonical().Phase, 3*Pi/2)
		assert.Equal(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, 0, -Delta / 2}.Canonical().Phase, 0.0)
		geomtest.AssertNumber(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, 0, 3 * Pi / 2}.Canonical().Phase, 3*Pi/2)
	})
}

func TestRegularPolygon_Grow(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Grow(2), RegPol(Pt(1, 2), Sz(12, 6), 6, 0, 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Grow(-6), RegPol(Pt(1, 2), Sz(4, 0), 6, 0, 0))
	})
	t.Run("shrink undoes it", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertRegularPolygon(t, rp.Grow(1.5).Shrink(1.5), rp, rp.String())
		}
	})
}

func TestRegularPolygon_GrowXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).GrowXY(2, 3), RegPol(Pt(1, 2), Sz(12, 7), 6, 0, 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).GrowXY(0, -6), RegPol(Pt(1, 2), Sz(10, 0), 6, 0, 0))
	})
}

func TestRegularPolygon_Shrink(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Shrink(2), RegPol(Pt(1, 2), Sz(8, 2), 6, 0, 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Shrink(6), RegPol(Pt(1, 2), Sz(4, 0), 6, 0, 0))
	})
}

func TestRegularPolygon_ShrinkXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).ShrinkXY(2, 3), RegPol(Pt(1, 2), Sz(8, 1), 6, 0, 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).ShrinkXY(0, 6), RegPol(Pt(1, 2), Sz(10, 0), 6, 0, 0))
	})
}

func TestRegularPolygon_Lerp(t *testing.T) {
	a, b := RegPol(Pt(0.0, 0.0), SzU(2.0), 6, 0, 0), RegPol(Pt(10.0, 20.0), Sz(8.0, 4.0), 6, Pi/2, 0)

	t.Run("moves the center, the size and the angle together", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, a.Lerp(b, 0.5), RegPol(Pt(5.0, 10.0), Sz(5.0, 3.0), 6, Pi/4, 0))
	})
	t.Run("turns along the shorter arc across the seam", func(t *testing.T) {
		from, to := RegPol(Pt(0, 0), SzU(2), 4, ToRadians(350), 0), RegPol(Pt(0, 0), SzU(2), 4, ToRadians(10), 0)

		geomtest.AssertRegularPolygon(t, from.Lerp(to, 0.5), RegPol(Pt(0, 0), SzU(2), 4, 0, 0))
		geomtest.AssertRegularPolygon(t, to.Lerp(from, 0.5), RegPol(Pt(0, 0), SzU(2), 4, 0, 0))
	})
	t.Run("the ends are the polygons themselves", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, a.Lerp(b, 0), a)
		geomtest.AssertRegularPolygon(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates and keeps the size absolute", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, a.Lerp(b, 2), RegPol(Pt(20.0, 40.0), Sz(14.0, 6.0), 6, Pi, 0))
		geomtest.AssertRegularPolygon(t, b.Lerp(a, 2), RegPol(Pt(-10.0, -20.0), Sz(4.0, 0.0), 6, 3*Pi/2, 0))
	})
	t.Run("int rounds the center and the size", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0, 0), SzU(2), 4, 0, 0).Lerp(RegPol(Pt(5, 5), SzU(5), 4, 0, 0), 0.5), RegPol(Pt(3, 3), SzU(4), 4, 0, 0))
	})
	t.Run("a different vertex count panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			a.Lerp(RegPol(Pt(10.0, 20.0), Sz(8.0, 4.0), 5, Pi/2, 0), 0.5)
		}, "geom: lerp between polygons of 6 and 5 vertices")
	})
	t.Run("the ends hold over the fixtures and the angle never turns more than half a turn", func(t *testing.T) {
		for _, from := range regularPolygonFixtures {
			for _, to := range regularPolygonFixtures {
				if from.N != to.N {
					continue
				}

				geomtest.AssertRegularPolygon(t, from.Lerp(to, 0), from, fmt.Sprintf("%s → %s: ", from, to))
				geomtest.AssertRegularPolygon(t, from.Lerp(to, 1), to, fmt.Sprintf("%s → %s: ", from, to))
				assert.True(t, AngleDistance(from.Lerp(to, 0.5).Angle, from.Angle) <= Pi/2+Delta, fmt.Sprintf("%s → %s: ", from, to))
			}
		}
	})
	t.Run("moves the phase along the shorter arc", func(t *testing.T) {
		from, to := RegPol(Pt(0.0, 0.0), Sz(2.0, 1.0), 6, 0, ToRadians(350)), RegPol(Pt(0.0, 0.0), Sz(2.0, 1.0), 6, 0, ToRadians(30))

		geomtest.AssertNumber(t, from.Lerp(to, 0.5).Phase, ToRadians(10))
	})
}

func TestRegularPolygon_Transform(t *testing.T) {
	hexagon := Hexagon(Pt(2.0, 3.0), SzU(4.0), OrientationFlatTop)

	t.Run("the identity keeps the polygon", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hexagon.Transform(IdentityMatrix[float64]()), hexagon)
	})
	t.Run("a translation moves the center and keeps the vertex count", func(t *testing.T) {
		moved := hexagon.Transform(TranslationMatrix(1.0, -1.0))

		geomtest.AssertRegularPolygon(t, moved, hexagon.Translate(Vec(1.0, -1.0)))
		assert.Equal(t, moved.N, 6)
	})
	t.Run("a uniform scale scales the semi-axes", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hexagon.Transform(ScaleMatrix(2.0, 2.0)), Hexagon(Pt(4.0, 6.0), SzU(8.0), OrientationFlatTop))
	})
	t.Run("a rotation turns the angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hexagon.Transform(RotationMatrix[float64](Pi/2)), Hexagon(Pt(-3.0, 2.0), SzU(4.0), OrientationFlatTop).Rotate(Pi/2))
	})
	t.Run("a reflection mirrors the angle and the phase about the axis of the matrix", func(t *testing.T) {
		turned := hexagon.Rotate(Pi / 6)

		geomtest.AssertRegularPolygon(t, turned.Transform(ReflectionMatrix[float64](AxisHorizontal)), RegPol(Pt(2.0, -3.0), SzU(4.0), 6, -turned.Angle, -turned.Phase))
	})
	t.Run("a shear is exact", func(t *testing.T) {
		shear := ShearMatrix(1.0, 0.0)

		geomtest.AssertPolygon(t, hexagon.Transform(shear).Polygon(), hexagon.Polygon().Transform(shear))
	})
	t.Run("an empty polygon stays empty", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1.0, 1.0), SzU(2.0), 0, 0, 0).Transform(ScaleMatrix(2.0, 2.0)).IsEmpty())
	})
	t.Run("int rounds the center and the semi-axes once", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 1), SzU(3), 6, 0, 0).Transform(ScaleMatrix(1.5, 1.5)), RegPol(Pt(2, 2), SzU(5), 6, 0, 0))
	})
	t.Run("a turn of unequal semi-axes is exact too", func(t *testing.T) {
		ellipse := RegPol(Pt(0.0, 0.0), Sz(2.0, 4.0), 5, 0, 0)
		turn := RotationMatrix[float64](Pi / 2)

		geomtest.AssertRegularPolygon(t, ellipse.Transform(turn), RegPol(Pt(0.0, 0.0), Sz(2.0, 4.0), 5, Pi/2, 0))
		geomtest.AssertPolygon(t, ellipse.Transform(turn).Polygon(), ellipse.Polygon().Transform(turn))
	})
	t.Run("axes scaled by different factors of a turned polygon are exact", func(t *testing.T) {
		turned := RegPol(Pt(0.0, 0.0), Sz(4.0, 2.0), 5, Pi/4, Pi/7)
		squeeze := ScaleMatrix(2.0, 3.0)

		geomtest.AssertPolygon(t, turned.Transform(squeeze).Polygon(), turned.Polygon().Transform(squeeze))
	})
	t.Run("a turn and a scale along the axes keep the phase", func(t *testing.T) {
		rp := RegPol(Pt(0.0, 0.0), Sz(4.0, 2.0), 5, 0, Pi/7)

		geomtest.AssertNumber(t, rp.Transform(RotationMatrix[float64](Pi/3)).Phase, Pi/7)
		geomtest.AssertNumber(t, rp.Transform(ScaleMatrix(2.0, 3.0)).Phase, Pi/7)
	})
	t.Run("matches the polygon of the vertices for every matrix, vertex by vertex", func(t *testing.T) {
		matrices := []Matrix[float64]{
			IdentityMatrix[float64](),
			TranslationMatrix(3.0, -2.0),
			ScaleMatrix(2.0, 2.0),
			ScaleMatrix(2.0, 3.0),
			RotationMatrix[float64](Pi / 3),
			RotationMatrix[float64](Pi / 3).Multiply(ScaleMatrix(2.0, 0.5)),
			ShearMatrix(0.5, -0.25),
			ReflectionMatrix[float64](AxisVertical),
			Mat(1.5, 0.5, 3.0, -0.25, 0.75, -2.0),
			Mat(-1.0, 2.0, 0.0, 0.5, 1.0, 1.0),
		}

		for _, rp := range regularPolygonFixtures {
			for _, m := range matrices {
				actual, expected := slices.Collect(rp.Transform(m).Vertices()), rp.Polygon().Transform(m).Points

				assert.Equal(t, len(actual), len(expected), fmt.Sprintf("%s → %s: ", rp, m))
				for i, vertex := range expected {
					j := i
					if m.Determinant() < 0 {
						j = (len(expected) - i) % len(expected)
					}

					geomtest.AssertPoint(t, actual[j], vertex, fmt.Sprintf("%s → %s: vertex %d: ", rp, m, i))
				}
			}
		}
	})
	t.Run("a long thin polygon keeps its lesser semi-axis", func(t *testing.T) {
		thin := RegPol(Pt(0.0, 0.0), Sz(1e6, 1e-3), 6, 0, 0)

		assert.EqualDelta(t, thin.Transform(IdentityMatrix[float64]()).Size.Height, 1e-3, 1e-18)
	})
}

func TestRegularPolygon_Rotate(t *testing.T) {
	t.Run("adds to the stored angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Rotate(Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi, 0))
	})
	t.Run("normalizes to the unit turn", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Rotate(3*Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi, 0))      // 0 + 3π → π
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Rotate(-Pi/2), RegPol(Pt(1, 2), Sz(2, 2), 4, 3*Pi/2, 0)) // 0 − π/2 → 3π/2
	})
}

func TestRegularPolygon_AlignTo(t *testing.T) {
	hex := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationPointyTop)

	t.Run("an anchor lands on the point", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hex.AlignTo(Top, Pt(5.0, 5.0)), Hexagon(Pt(5.0, 15.0), SzU(10.0), OrientationPointyTop))
	})
	t.Run("none aligns the center", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, hex.AlignTo(DirectionNone, Pt(5.0, 5.0)), Hexagon(Pt(5.0, 5.0), SzU(10.0), OrientationPointyTop))
	})
}

func TestRegularPolygon_Contains(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, diamond.Contains(Pt(0, 0)))
		assert.True(t, diamond.Contains(Pt(1, 0)))
	})
	t.Run("on an edge and on a vertex", func(t *testing.T) {
		assert.True(t, diamond.Contains(Pt(1, 1)))
		assert.True(t, diamond.Contains(Pt(2, 0)))
	})
	t.Run("outside", func(t *testing.T) {
		assert.False(t, diamond.Contains(Pt(2, 2)))
		assert.False(t, diamond.Contains(Pt(3, 0)))
	})
	t.Run("within the tolerance of an edge", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0, 0).Contains(Pt(1.0, 1.0+Delta/2)))
		assert.False(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0, 0).Contains(Pt(1.0, 1.0+2*Delta)))
	})
	t.Run("an empty polygon contains nothing", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).Contains(Pt(0, 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, rp.Contains(p), rp.Polygon().Contains(p), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
	t.Run("int matches the polygon of the rounded vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				assert.Equal(t, rp.Int().Contains(p.Int()), rp.Int().Polygon().Contains(p.Int()), fmt.Sprintf("%s → %s: ", rp.Int(), p.Int()))
			}
		}
	})
}

func BenchmarkRegularPolygon_Contains(b *testing.B) {
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)
	inside, outside := Pt(1.0, 2.0), Pt(9.0, 9.0)

	b.Run("inside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = hexagon.Contains(inside)
		}
	})
	b.Run("outside within the extent", func(b *testing.B) {
		for b.Loop() {
			sinkBool = hexagon.Contains(outside)
		}
	})
}

func TestRegularPolygon_DistanceTo(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("zero within the polygon", func(t *testing.T) {
		assert.Equal(t, diamond.DistanceTo(Pt(0, 0)), 0.0)
		assert.Equal(t, diamond.DistanceTo(Pt(1, 1)), 0.0)
	})
	t.Run("to the nearest edge", func(t *testing.T) {
		geomtest.AssertNumber(t, diamond.DistanceTo(Pt(2, 2)), Sqrt2)
	})
	t.Run("to the nearest vertex", func(t *testing.T) {
		geomtest.AssertNumber(t, diamond.DistanceTo(Pt(3, 0)), 1.0)
	})
	t.Run("an empty polygon is infinitely far", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).DistanceTo(Pt(0, 0)), math.Inf(1))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, rp.DistanceTo(p), rp.Polygon().DistanceTo(p), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
}

func TestRegularPolygon_DistanceSquaredTo(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		geomtest.AssertNumber(t, diamond.DistanceSquaredTo(Pt(2, 2)), 2.0)
		geomtest.AssertNumber(t, diamond.DistanceSquaredTo(Pt(3, 0)), 1.0)
		assert.Equal(t, diamond.DistanceSquaredTo(Pt(0, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				geomtest.AssertNumber(t, rp.DistanceSquaredTo(p), rp.DistanceTo(p)*rp.DistanceTo(p), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
}

func TestRegularPolygon_Nearest(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("the foot on the nearest edge", func(t *testing.T) {
		geomtest.AssertPoint(t, diamond.Nearest(Pt(2, 2)), Pt(1, 1))
		geomtest.AssertPoint(t, diamond.Nearest(Pt(3, 0)), Pt(2, 0))
	})
	t.Run("a point inside is its own nearest point", func(t *testing.T) {
		geomtest.AssertPoint(t, diamond.Nearest(Pt(0, 1)), Pt(0, 1))
	})
	t.Run("an empty polygon returns the zero point", func(t *testing.T) {
		geomtest.AssertPoint(t, RegPol(Pt(1, 1), Sz(2, 2), 0, 0, 0).Nearest(Pt(3, 4)), Pt(0, 0))
	})
	t.Run("over the fixtures", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				assertNearest[float64](t, rp, p)
			}
		}
	})
}

func TestRegularPolygon_EnclosesCircle(t *testing.T) {
	square := RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, Pi/4, 0)

	t.Run("the inscribed circle touches every side from inside", func(t *testing.T) {
		assert.True(t, square.EnclosesCircle(Circ(Pt(0.0, 0.0), math.Sqrt2)))
		assert.False(t, square.EnclosesCircle(Circ(Pt(0.0, 0.0), 1.5)))
	})
	t.Run("an empty polygon encloses nothing", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).EnclosesCircle(Circ(Pt(0, 0), 0)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, rp.EnclosesCircle(c), rp.Polygon().EnclosesCircle(c), fmt.Sprintf("%s → %s: ", rp, c))
			}
		}
	})
}

func TestRegularPolygon_EnclosesSegment(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("a diagonal counts, a segment past a vertex not", func(t *testing.T) {
		assert.True(t, diamond.EnclosesSegment(Seg(Pt(-2, 0), Pt(2, 0))))
		assert.False(t, diamond.EnclosesSegment(Seg(Pt(-2, 0), Pt(3, 0))))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, rp.EnclosesSegment(s), rp.Polygon().EnclosesSegment(s), fmt.Sprintf("%s → %s: ", rp, s))
			}
		}
	})
}

func TestRegularPolygon_EnclosesPolygon(t *testing.T) {
	t.Run("inside, inscribed and outside", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1, 1), Sz(2, 2), 4, Pi/4, 0).EnclosesPolygon(Pol(squareVertices())))
		assert.True(t, RegPol(Pt(1, 1), Sz(2, 2), 4, 0, 0).EnclosesPolygon(Pol(squareVertices())))
		assert.False(t, RegPol(Pt(1, 1), Sz(1, 1), 4, 0, 0).EnclosesPolygon(Pol(squareVertices())))
	})
	t.Run("an empty polygon is enclosed by nothing", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0).EnclosesPolygon(Pol[int](nil)))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, rp.EnclosesPolygon(p), rp.Polygon().EnclosesPolygon(p), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
}

func TestRegularPolygon_EnclosesRectangle(t *testing.T) {
	t.Run("a square rotated into a diamond", func(t *testing.T) {
		diamond := RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0, 0)

		assert.True(t, diamond.EnclosesRectangle(Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))))
		assert.False(t, diamond.EnclosesRectangle(Rect(Pt(0.0, 0.0), Sz(2.0, 2.2))))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, r := range rectFixtures {
				assert.Equal(t, rp.EnclosesRectangle(r), rp.Polygon().EnclosesRectangle(r), fmt.Sprintf("%s → %s: ", rp, r))
			}
		}
	})
}

func TestRegularPolygon_EnclosesRegularPolygon(t *testing.T) {
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)

	t.Run("a smaller one and a turned one of the same size", func(t *testing.T) {
		assert.True(t, hexagon.EnclosesRegularPolygon(hexagon.Unscale(2)))
		assert.False(t, hexagon.EnclosesRegularPolygon(hexagon.Rotate(Pi/6)))
	})
	t.Run("an empty polygon on either side encloses nothing", func(t *testing.T) {
		assert.False(t, hexagon.EnclosesRegularPolygon(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 0, 0, 0)))
		assert.False(t, RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 0, 0, 0).EnclosesRegularPolygon(hexagon))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, a := range regularPolygonFixtures {
			for _, b := range regularPolygonFixtures {
				assert.Equal(t, a.EnclosesRegularPolygon(b), a.Polygon().EnclosesRegularPolygon(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRegularPolygon_EnclosesBox(t *testing.T) {
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, rp.EnclosesBox(b), rp.Polygon().EnclosesBox(b), fmt.Sprintf("%s → %s: ", rp, b))
			}
		}
	})
}

func TestRegularPolygon_IntersectsCircle(t *testing.T) {
	t.Run("mirrors Circle.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, c := range circleFixtures {
				assert.Equal(t, rp.IntersectsCircle(c), c.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, c))
			}
		}
	})
}

func TestRegularPolygon_IntersectsSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				assert.Equal(t, rp.IntersectsSegment(s), s.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, s))
			}
		}
	})
}

func TestRegularPolygon_IntersectionSegment(t *testing.T) {
	t.Run("mirrors Segment.IntersectionRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, rp.IntersectionSegment(s), s.IntersectionRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, s))
			}
		}
	})
}

func TestRegularPolygon_AppendIntersectionSegment(t *testing.T) {
	t.Run("matches Segment.AppendIntersectionRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, s := range segmentFixtures {
				geomtest.AssertVertices(t, rp.AppendIntersectionSegment(bufferWith(prefixPoint), s), s.AppendIntersectionRegularPolygon(bufferWith(prefixPoint), rp), fmt.Sprintf("%s → %s: ", rp, s))
			}
		}
	})
}

func TestRegularPolygon_IntersectsRay(t *testing.T) {
	t.Run("mirrors Ray.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, r := range rayFixtures {
				assert.Equal(t, rp.IntersectsRay(r), r.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, r))
			}
		}
	})
}

func TestRegularPolygon_IntersectionRay(t *testing.T) {
	t.Run("matches Ray.IntersectionRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, rp.IntersectionRay(r), r.IntersectionRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, r))
			}
		}
	})
}

func TestRegularPolygon_AppendIntersectionRay(t *testing.T) {
	t.Run("matches Ray.AppendIntersectionRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, r := range rayFixtures {
				geomtest.AssertVertices(t, rp.AppendIntersectionRay(bufferWith(prefixPoint), r), r.AppendIntersectionRegularPolygon(bufferWith(prefixPoint), rp), fmt.Sprintf("%s → %s: ", rp, r))
			}
		}
	})
}

func TestRegularPolygon_IntersectsPolygon(t *testing.T) {
	t.Run("mirrors Polygon.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range polygonFixtures() {
				assert.Equal(t, rp.IntersectsPolygon(p), p.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
}

func TestRegularPolygon_IntersectsRectangle(t *testing.T) {
	t.Run("mirrors Rectangle.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, r := range rectFixtures {
				assert.Equal(t, rp.IntersectsRectangle(r), r.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, r))
			}
		}
	})
}

func TestRegularPolygon_IntersectsRegularPolygon(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, diamond.IntersectsRegularPolygon(diamond.Translate(Vec(1, 1))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, diamond.IntersectsRegularPolygon(diamond.Translate(Vec(5, 0))))
	})
	t.Run("a shared vertex counts", func(t *testing.T) {
		assert.True(t, diamond.IntersectsRegularPolygon(diamond.Translate(Vec(4, 0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, diamond.IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(1, 1), 3, 0, 0)))
		assert.True(t, RegPol(Pt(0, 0), Sz(1, 1), 3, 0, 0).IntersectsRegularPolygon(diamond))
	})
	t.Run("edges crossing without a vertex inside", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0, 0).IntersectsRegularPolygon(RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, Pi/4, 0)))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, diamond.IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0)))
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).IntersectsRegularPolygon(diamond))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range regularPolygonFixtures {
			for _, b := range regularPolygonFixtures {
				assert.Equal(t, a.IntersectsRegularPolygon(b), b.IntersectsRegularPolygon(a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("matches the polygons of the vertices", func(t *testing.T) {
		for _, a := range regularPolygonFixtures {
			for _, b := range regularPolygonFixtures {
				assert.Equal(t, a.IntersectsRegularPolygon(b), a.Polygon().IntersectsPolygon(b.Polygon()), fmt.Sprintf("%s → %s: ", a, b))
				assert.Equal(t, a.IntersectsRegularPolygon(b.Translate(Vec(3.0, -2.0))), a.Polygon().IntersectsPolygon(b.Translate(Vec(3.0, -2.0)).Polygon()), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestRegularPolygon_IntersectsBox(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0, 0)

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, diamond.IntersectsBox(BoxFromMinMax(Pt(0, 0), Pt(2, 2))))
	})
	t.Run("apart within overlapping bounds", func(t *testing.T) {
		assert.False(t, diamond.IntersectsBox(BoxFromMinMax(Pt(2, 2), Pt(3, 3))))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).IntersectsBox(BoxFromMinMax(Pt(-1, -1), Pt(1, 1))))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, b := range boxFixtures {
				assert.Equal(t, rp.IntersectsBox(b), rp.Polygon().IntersectsBox(b), fmt.Sprintf("%s → %s: ", rp, b))
			}
		}
	})
}

func TestRegularPolygon_Equal(t *testing.T) {
	t.Run("same polygon", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0)))
	})
	t.Run("different polygon", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Equal(RegPol(Pt(0, 0), Sz(2, 2), 6, 0, 0)))
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Equal(RegPol(Pt(1, 2), Sz(3, 3), 4, 0, 0)))
	})
	t.Run("the angle is compared", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, Pi, 0)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0.5, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 0.5, 0)))
	})
	t.Run("the angle is compared normalized", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationPointyTop)

		assert.True(t, hexagon.Equal(hexagon.Rotate(0)))
		assert.True(t, hexagon.Equal(hexagon.Rotate(2*Pi)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, -Pi/2, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 3*Pi/2, 0)))
	})
	t.Run("the angle is compared across the seam", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationPointyTop)

		assert.True(t, hexagon.Equal(hexagon.Rotate(-1e-9)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 2*Pi-1e-9, 0)))
	})
	t.Run("the phase is compared as a value, not as the shape it draws", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1.0, 2.0), Sz(2.0, 1.0), 4, 0, 0).Equal(RegPol(Pt(1.0, 2.0), Sz(2.0, 1.0), 4, 0, Pi/4)))
		assert.False(t, RegPol(Pt(1.0, 2.0), SzU(2.0), 4, 0, Pi/2).Equal(RegPol(Pt(1.0, 2.0), SzU(2.0), 4, 0, 0)))
		assert.True(t, RegPol(Pt(1.0, 2.0), SzU(2.0), 4, 0, Pi/4).Equal(RegPol(Pt(1.0, 2.0), SzU(2.0), 4, 0, Pi/4+2*Pi)))
	})
}

func TestRegularPolygon_IsZero(t *testing.T) {
	t.Run("zero polygon", func(t *testing.T) {
		assert.True(t, RegularPolygon[int]{}.IsZero())
		assert.True(t, RegularPolygon[float64]{}.IsZero())
	})
	t.Run("a full turn is a zero angle, like Equal", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0, 0), Sz(0, 0), 0, 2*Pi, 0).IsZero())
	})
	t.Run("only one field is zero", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(0, 0), 4, 0, 0).IsZero())
		assert.False(t, RegPol(Pt(1, 2), Sz(0, 0), 0, 0, 0).IsZero())
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0, 0).IsZero())
		assert.False(t, RegPol(Pt(0, 0), Sz(0, 0), 0, 1, 0).IsZero())
	})
	t.Run("non-zero polygon", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).IsZero())
	})
}

func TestRegularPolygon_IsEmpty(t *testing.T) {
	t.Run("no sides", func(t *testing.T) {
		assert.True(t, RegularPolygon[int]{}.IsEmpty())
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 0, 0, 0).IsEmpty())
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), -1, 0, 0).IsEmpty())
	})
	t.Run("with sides", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).IsEmpty())
	})
}

func TestRegularPolygon_IsAligned(t *testing.T) {
	t.Run("no turn", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1, 2), Sz(3, 4), 6, 0, 0).IsAligned())
	})

	t.Run("turned", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(3, 4), 6, Pi/6, 0).IsAligned())
	})

	t.Run("a full turn is exactly zero again", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1.0, 2.0), Sz(3.0, 4.0), 6, 0, 0).Rotate(2*Pi).IsAligned())
	})

	t.Run("no tolerance", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1.0, 2.0), Sz(3.0, 4.0), 6, 1e-9, 0).IsAligned())
		assert.True(t, RegPol(Pt(1.0, 2.0), Sz(3.0, 4.0), 6, 1e-9, 0).Canonical().IsAligned())
	})
}

func TestRegularPolygon_Polygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0, 0)
	p := rp.Polygon()

	t.Run("carries the vertices", func(t *testing.T) {
		geomtest.AssertVertices(t, p.Points, slices.Collect(rp.Vertices()))
	})
	t.Run("owns its slice", func(t *testing.T) {
		assert.NotSame(t, p.Points, rp.Polygon().Points)
	})
	t.Run("fewer than one side is the zero polygon", func(t *testing.T) {
		assert.Zero(t, RegPol(Pt(0, 0), Sz(1, 1), 0, 0, 0).Polygon())
		assert.Zero(t, RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0, 0).Polygon())
	})
}

func TestRegularPolygon_Ellipse(t *testing.T) {
	t.Run("the same center, semi-axes and angle", func(t *testing.T) {
		geomtest.AssertEllipse(t, RegPol(Pt(1, 2), Sz(10, 4), 6, 0, 0).Ellipse(), Ell(Pt(1, 2), Sz(10, 4), 0))
		geomtest.AssertEllipse(t, RegPol(Pt(1.0, 2.0), Sz(10.0, 4.0), 6, 1.0, 0).Ellipse(), Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 1.0))
	})
	t.Run("every vertex lies on it", func(t *testing.T) {
		rp := RegPol(Pt(-3.5, 0.25), Sz(2.0, 3.0), 5, Pi/7, 0)

		for vertex := range rp.Vertices() {
			geomtest.AssertNumber(t, rp.Ellipse().DistanceTo(vertex), 0.0, fmt.Sprintf("%s: ", vertex))
		}
	})
	t.Run("an empty polygon still names its ellipse", func(t *testing.T) {
		geomtest.AssertEllipse(t, RegPol(Pt(1, 2), Sz(10, 4), 0, 0, 0).Ellipse(), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
}

func TestRegularPolygon_Circle(t *testing.T) {
	t.Run("the circumscribed circle of a polygon of equal semi-axes", func(t *testing.T) {
		geomtest.AssertCircle(t, Hexagon(Pt(1, 2), SzU(10), OrientationFlatTop).Circle(), Circ(Pt(1, 2), 10))
	})
	t.Run("the circle around an elliptical polygon", func(t *testing.T) {
		geomtest.AssertCircle(t, RegPol(Pt(1, 2), Sz(4, 10), 5, 0, 0).Circle(), Circ(Pt(1, 2), 10))
	})
}

func TestRegularPolygon_Cast(t *testing.T) {
	rp := RegPol(Pt(1.5, -2.5), Sz(3.5, 4.5), 6, Pi/6, 0)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, rp.Cast[int](), rp.Int())
		geomtest.AssertRegularPolygon(t, rp.Cast[float64](), rp.Float())
	})
	t.Run("a type the other conversions cannot name, and N and the angle are kept", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, rp.Cast[int8](), RegPol(Pt[int8](2, -3), Sz[int8](4, 5), 6, Pi/6, 0))
	})
}

func TestRegularPolygon_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Int(), RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	})
	t.Run("float rounds but keeps the angle", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3, 0).Int(), RegPol(Pt(1, 0), Sz(1, 4), 6, Pi/3, 0))
	})
}

func TestRegularPolygon_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).Float(), RegPol(Pt(1.0, 2.0), Sz(2.0, 2.0), 4, 0, 0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertRegularPolygon(t, RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3, 0).Float(), RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3, 0))
	})
}

func TestRegularPolygon_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0).String(), "RegPol((1,2);2x2;4)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, Pi, 0).String(), "RegPol((0.50,-1.25);2.50x3.75;6;3.14)")
	})
	t.Run("a phase is printed after the angle, the angle with it", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(1.0, 2.0), Sz(2.0, 1.0), 4, 0, 0.5).String(), "RegPol((1.00,2.00);2.00x1.00;4;0.00;0.50)")
	})
}

func TestRegularPolygon_JSON(t *testing.T) {
	t.Run("int wire format omits a zero angle", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0), `{"x":1,"y":2,"w":2,"h":2,"n":4}`)

		var rp RegularPolygon[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4}`), &rp))
		geomtest.AssertRegularPolygon(t, rp, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`), &rp))
		geomtest.AssertRegularPolygon(t, rp, RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, 0.5, 0), `{"x":0.50,"y":-1.25,"w":2.50,"h":3.75,"n":6,"a":0.5}`)

		var rp RegularPolygon[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":0.50,"y":-1.25,"w":2.50,"h":3.75,"n":6,"a":0.5}`), &rp))
		geomtest.AssertRegularPolygon(t, rp, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, 0.5, 0))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			data, err := json.Marshal(rp)
			assert.NoError(t, err)

			var decoded RegularPolygon[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, rp)
		}
	})
	t.Run("a phase is carried as p, and only where it is not zero", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(1.0, 2.0), Sz(2.0, 1.0), 4, 0, 0.5), `{"x":1,"y":2,"w":2,"h":1,"n":4,"p":0.5}`)

		var rp RegularPolygon[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":1,"n":4,"p":0.5}`), &rp))
		geomtest.AssertRegularPolygon(t, rp, RegPol(Pt(1.0, 2.0), Sz(2.0, 1.0), 4, 0, 0.5))
	})
}

func TestRegularPolygonOrientationPhase(t *testing.T) {
	t.Run("pointy top puts the first vertex at -Y", func(t *testing.T) {
		// in +Y-down screen coordinates, "top" means minimum Y, so the first vertex
		// has to point in the -Y direction: -π/2, stored normalized as 3π/2
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(3, OrientationPointyTop), 270*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(4, OrientationPointyTop), 270*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(6, OrientationPointyTop), 270*DegToRad)
	})
	t.Run("flat top sits half a step before the top", func(t *testing.T) {
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(3, OrientationFlatTop), 210*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(4, OrientationFlatTop), 225*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(5, OrientationFlatTop), 234*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(6, OrientationFlatTop), 240*DegToRad)
	})
	t.Run("every n places its top", func(t *testing.T) {
		for _, orientation := range Orientations() {
			for n := 3; n <= 9; n++ {
				assertOrientation(t, RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), n, orientation), orientation, fmt.Sprintf("%s n=%d: ", orientation, n))
			}
		}
	})
	t.Run("an orientation that is neither panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			RegularPolygonOrientationPhase(6, Orientation(99))
		}, "geom: unknown orientation 99")

		assert.PanicsWith(t, func() {
			RegularPolygonOrientationPhase(6, OrientationNone)
		}, "geom: unknown orientation -1")
	})
	t.Run("no vertices give the top phase instead of dividing by n", func(t *testing.T) {
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(0, OrientationFlatTop), 270*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(-1, OrientationFlatTop), 270*DegToRad)
		geomtest.AssertNumber(t, RegularPolygonOrientationPhase(0, OrientationPointyTop), 270*DegToRad)

		empty := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 0, OrientationFlatTop)
		assert.True(t, empty.Equal(empty))
	})
}

func TestRegularPolygon_Properties(t *testing.T) {
	t.Run("vertex count matches the side count", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			assert.Length(t, slices.Collect(rp.Vertices()), rp.N, fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("vertices lie on the ellipse of the size, in the frame before the turn", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for vertex := range rp.Vertices() {
				offset := vertex.Subtract(rp.Center).Rotate(-rp.Angle)
				normalized := offset.X*offset.X/(rp.Size.Width*rp.Size.Width) + offset.Y*offset.Y/(rp.Size.Height*rp.Size.Height)

				geomtest.AssertNumber(t, normalized, 1.0, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("every anchor is on the boundary in its direction", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, direction := range Directions() {
				anchor := rp.Anchor(direction)

				geomtest.AssertNumber(t, rp.DistanceTo(anchor), 0.0, fmt.Sprintf("%s → %s: ", rp, direction))
				geomtest.AssertAngle(t, rp.Center.AngleTo(anchor), direction.Angle(), fmt.Sprintf("%s → %s: ", rp, direction))
				geomtest.AssertRegularPolygon(t, rp.AlignTo(direction, anchor), rp, fmt.Sprintf("%s → %s: ", rp, direction))
			}
		}
	})
	t.Run("a full turn restores the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertVertices(t, slices.Collect(rp.Rotate(2*Pi).Vertices()), slices.Collect(rp.Vertices()), fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("rotating by one step permutes the vertices of equal semi-axes", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			if rp.Size.Width != rp.Size.Height {
				continue // a turned ellipse is a different shape, not the same ring one step on
			}

			rotated := slices.Collect(rp.Rotate(2 * Pi / float64(rp.N)).Vertices())
			vertices := slices.Collect(rp.Vertices())

			for i := range vertices {
				assert.True(t, rotated[i].Equal(vertices[(i+1)%rp.N]), fmt.Sprintf("%s #%d: ", rp, i))
			}
		}
	})
	t.Run("translate keeps the shape", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, vector := range vectorFixtures {
				moved := rp.Translate(vector)

				assert.True(t, moved.Center.Equal(rp.Center.Add(vector)), fmt.Sprintf("%s → %s: ", rp, vector))
				assert.True(t, moved.Size.Equal(rp.Size), fmt.Sprintf("%s → %s: ", rp, vector))
				geomtest.AssertNumber(t, moved.Angle, rp.Angle, fmt.Sprintf("%s → %s: ", rp, vector))
			}
		}
	})
	t.Run("bounds enclose every vertex", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			bounds := rp.Bounds()

			for vertex := range rp.Vertices() {
				assert.True(t, vertex.X >= bounds.Min.X-Delta && vertex.X <= bounds.Max.X+Delta, fmt.Sprintf("%s: ", rp))
				assert.True(t, vertex.Y >= bounds.Min.Y-Delta && vertex.Y <= bounds.Max.Y+Delta, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("polygon carries the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			geomtest.AssertVertices(t, rp.Polygon().Points, slices.Collect(rp.Vertices()), fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("every vertex lies on the ellipse", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for vertex := range rp.Vertices() {
				geomtest.AssertNumber(t, rp.Ellipse().DistanceTo(vertex), 0.0, fmt.Sprintf("%s → %s: ", rp, vertex))
			}
		}
	})
	t.Run("the ellipse round-trips through the polygon of the phase of an orientation", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, orientation := range []Orientation{OrientationFlatTop, OrientationPointyTop} {
				oriented := RegPol(rp.Center, rp.Size, rp.N, rp.Angle, RegularPolygonOrientationPhase(rp.N, orientation))

				geomtest.AssertRegularPolygon(t, rp.Ellipse().RegularPolygon(rp.N, orientation), oriented, fmt.Sprintf("%s %s: ", rp, orientation))
			}
		}
	})
	t.Run("the circle around holds every vertex", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for vertex := range rp.Vertices() {
				assert.True(t, rp.Circle().Contains(vertex), fmt.Sprintf("%s → %s: ", rp, vertex))
			}
		}
	})
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RegPol(Pt(0, 1), Sz(2, 3), 4, 0, 0)

	rp.Translate(Vec(2, 3))
	rp.MoveTo(Pt(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)
	rp.Rotate(Pi / 3)

	geomtest.AssertRegularPolygon(t, rp, RegPol(Pt(0, 1), Sz(2, 3), 4, 0, 0))
}

// regularPolygonFixtures span the triangle, square, and hexagon at both orientations, and
// unequal semi-axes at a phase, turned and not.
var regularPolygonFixtures = []RegularPolygon[float64]{
	Triangle(Pt(0.0, 0.0), Sz(3.0, 3.0), OrientationPointyTop),
	Square(Pt(50.0, 50.0), Sz(100.0, 100.0), OrientationPointyTop),
	Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), OrientationFlatTop),
	RegPol(Pt(-3.5, 0.25), Sz(2.0, 3.0), 5, Pi/7, 0),
	RegPol(Pt(12.5, -0.1), Sz(0.5, 0.5), 8, 0, 0),
	Hexagon(Pt(0.0, 0.0), Sz(20.0, 10.0), OrientationFlatTop),
	Square(Pt(1.0, 2.0), Sz(4.0, 2.0), OrientationFlatTop),
	RegPol(Pt(1.0, -2.0), Sz(6.0, 3.0), 7, Pi/5, Pi/9),
	RegPol(Pt(-2.0, 4.0), Sz(5.0, 1.5), 6, 2*Pi/3, Pi/11),
}

func ExampleRegPol() {
	fmt.Println(RegPol(Pt(1, 2), Sz(2, 2), 4, 0, 0))
	// Output: RegPol((1,2);2x2;4)
}

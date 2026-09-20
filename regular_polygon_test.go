package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestRegularPolygon_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0, 0), Sz(2, 2), 4, 0), RegularPolygon[int]{Center: Pt(0, 0), Size: Sz(2, 2), N: 4, Angle: 0})
	})
	t.Run("float", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, Pi/3), RegularPolygon[float64]{Center: Pt(0.5, -1.25), Size: Sz(2.5, 3.75), N: 6, Angle: Pi / 3})
	})
	t.Run("a negative size is taken absolute", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0, 0), Sz(-2, 3), 4, 0), RegPol(Pt(0, 0), Sz(2, 3), 4, 0))
		AssertRegularPolygon(t, Square(Pt(0.0, 0.0), Sz(-2.0, -3.0), PointyTop), Square(Pt(0.0, 0.0), Sz(2.0, 3.0), PointyTop))
	})
}

func TestRegularPolygonWithOrientation(t *testing.T) {
	center, size := Pt(0, 0), Sz(10, 10)
	pointy := RegularPolygonWithOrientation(center, size, 6, PointyTop)
	flat := RegularPolygonWithOrientation(center, size, 6, FlatTop)

	t.Run("carries the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, pointy, RegPol(center, size, 6, RegularPolygonOrientationAngle(6, PointyTop)))
		AssertRegularPolygon(t, flat, RegPol(center, size, 6, RegularPolygonOrientationAngle(6, FlatTop)))

		assert.NotEqual(t, pointy.Angle, flat.Angle)
	})
	t.Run("pointy top starts above the center", func(t *testing.T) {
		assert.True(t, slices.Collect(pointy.Vertices())[0].Y < center.Y)
	})
	t.Run("flat top starts left of and above the center", func(t *testing.T) {
		first := slices.Collect(flat.Vertices())[0]

		assert.True(t, first.X < center.X)
		assert.True(t, first.Y < center.Y)
	})
}

func TestTriangle(t *testing.T) {
	triangle := Triangle(Pt(1, -1), Sz(3, 3), PointyTop)

	t.Run("has three sides at the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, triangle, RegPol(Pt(1, -1), Sz(3, 3), 3, RegularPolygonOrientationAngle(3, PointyTop)))
	})
	t.Run("integer vertices round after scaling", func(t *testing.T) {
		// vertices 1 and 2 land within one unit of the exact (3.598, 0.5) and (-1.598, 0.5);
		// their Y of 0.5 falls just below the .5 tie in float64 and rounds down to 0
		AssertVertices(t, slices.Collect(triangle.Vertices()), []Point[int]{Pt(1, -4), Pt(4, 0), Pt(-2, 0)})
	})
}

func TestSquare(t *testing.T) {
	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointyTop)

	t.Run("has four sides at the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, square, RegPol(Pt(50.0, 50.0), Sz(100.0, 100.0), 4, RegularPolygonOrientationAngle(4, PointyTop)))
	})
	t.Run("pointy top starts at the visual top and winds to the right", func(t *testing.T) {
		AssertVertices(t, slices.Collect(square.Vertices()), []Point[float64]{Pt(50.0, -50.0), Pt(150.0, 50.0), Pt(50.0, 150.0), Pt(-50.0, 50.0)})
	})
}

func TestHexagon(t *testing.T) {
	// float64 avoids the int-rounding collapse where sin(±π/6) ≈ 0.4999 would truncate to 0
	hexagon := Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), PointyTop)

	t.Run("has six sides at the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon, RegPol(Pt(0.0, 0.0), Sz(10.0, 10.0), 6, RegularPolygonOrientationAngle(6, PointyTop)))
	})
	t.Run("pointy top has vertices at top and bottom and flat sides left and right", func(t *testing.T) {
		AssertVertices(t, slices.Collect(hexagon.Vertices()), []Point[float64]{
			Pt(0.0, -10.0),
			Pt(5*Sqrt3, -5.0),
			Pt(5*Sqrt3, 5.0),
			Pt(0.0, 10.0),
			Pt(-5*Sqrt3, 5.0),
			Pt(-5*Sqrt3, -5.0),
		})
	})
}

func TestRegularPolygonOrientationAngle(t *testing.T) {
	t.Run("pointy top points at -Y", func(t *testing.T) {
		// in +Y-down screen coordinates, "top" means minimum Y, so the first vertex
		// has to point in the -Y direction: angle = -π/2, stored normalized as 3π/2
		AssertNumber(t, RegularPolygonOrientationAngle(3, PointyTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(4, PointyTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(6, PointyTop), 270*DegToRad)
	})
	t.Run("flat top sits half a step before the top", func(t *testing.T) {
		AssertNumber(t, RegularPolygonOrientationAngle(3, FlatTop), 210*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(4, FlatTop), 225*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(5, FlatTop), 234*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(6, FlatTop), 240*DegToRad)
	})
	t.Run("flat top puts an edge midpoint at the top for any n", func(t *testing.T) {
		for n := 3; n <= 9; n++ {
			vertices := slices.Collect(RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), n, FlatTop).Vertices())

			AssertPoint(t, vertices[0].Midpoint(vertices[1]), Pt(0.0, -10*math.Cos(Pi/float64(n))), fmt.Sprintf("n=%d: ", n))
		}
	})
	t.Run("an orientation that is neither panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			RegularPolygonOrientationAngle(6, Orientation(99))
		}, "geom: unknown orientation 99")

		assert.PanicsWith(t, func() {
			RegularPolygonOrientationAngle(6, OrientationNone)
		}, "geom: unknown orientation -1")
	})
	t.Run("no vertices give the top angle instead of dividing by n", func(t *testing.T) {
		AssertNumber(t, RegularPolygonOrientationAngle(0, FlatTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(-1, FlatTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(0, PointyTop), 270*DegToRad)

		empty := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 0, FlatTop)
		assert.True(t, empty.Equal(empty))
	})
}

func TestRegularPolygon_Vertices(t *testing.T) {
	t.Run("fewer than one side yields nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 0, 0).Vertices()))
		assert.Nil(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0).Vertices()))
	})
	t.Run("int", func(t *testing.T) {
		AssertVertices(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices()), []Point[int]{
			Pt(1, 0),
			Pt(0, 1),
			Pt(-1, 0),
			Pt(0, -1),
		})
	})
	t.Run("non-square size stretches the ellipse", func(t *testing.T) {
		AssertVertices(t, slices.Collect(RegPol(Pt(0, 0), Sz(2, 3), 4, 0).Vertices()), []Point[int]{
			Pt(2, 0),
			Pt(0, 3),
			Pt(-2, 0),
			Pt(0, -3),
		})
	})
	t.Run("float", func(t *testing.T) {
		AssertVertices(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0).Vertices()), []Point[float64]{
			Pt(2.0, 0.0),
			Pt(1.0, 1.5*Sqrt3),
			Pt(-1.0, 1.5*Sqrt3),
			Pt(-2.0, 0.0),
			Pt(-1.0, -1.5*Sqrt3),
			Pt(1.0, -1.5*Sqrt3),
		})
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for vertex := range RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices() {
			AssertPoint(t, vertex, Pt(1, 0))

			break
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertNumber(t, testing.AllocsPerRun(100, func() {
				for vertex := range rp.Vertices() {
					sinkBool = vertex.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", rp))
		}
	})
}

func TestRegularPolygon_Edges(t *testing.T) {
	t.Run("closes back to the first vertex", func(t *testing.T) {
		edges := slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Edges())

		assert.Equal(t, len(edges), 4)
		AssertLine(t, edges[0], Ln(Pt(1, 0), Pt(0, 1)))
		AssertLine(t, edges[3], Ln(Pt(0, -1), Pt(1, 0)))
	})
	t.Run("single vertex is one zero-length edge", func(t *testing.T) {
		edges := slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 1, 0).Edges())

		assert.Equal(t, len(edges), 1)
		AssertLine(t, edges[0], Ln(Pt(1.0, 0.0), Pt(1.0, 0.0)))
	})
	t.Run("fewer than one side yields nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(RegPol(Pt(0, 0), Sz(1, 1), 0, 0).Edges()))
		assert.Nil(t, slices.Collect(RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0).Edges()))
	})
	t.Run("stops where the caller breaks", func(t *testing.T) {
		for edge := range RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Edges() {
			AssertLine(t, edge, Ln(Pt(1, 0), Pt(0, 1)))

			break
		}
	})
	t.Run("are the edges of the polygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			edges, expected := slices.Collect(rp.Edges()), slices.Collect(rp.Polygon().Edges())

			assert.Equal(t, len(edges), len(expected), fmt.Sprintf("%s: ", rp))
			for i := range edges {
				AssertLine(t, edges[i], expected[i], fmt.Sprintf("%s #%d: ", rp, i))
			}
		}
	})
	t.Run("ranging allocates nothing", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertNumber(t, testing.AllocsPerRun(100, func() {
				for edge := range rp.Edges() {
					sinkBool = edge.IsZero()
				}
			}), 0, fmt.Sprintf("%s: ", rp))
		}
	})
}

func TestRegularPolygon_Area(t *testing.T) {
	t.Run("agrees with the polygon", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(2.0), FlatTop)

		AssertNumber(t, hexagon.Area(), hexagon.Polygon().Area())
		AssertNumber(t, hexagon.Area(), 3*Sqrt3/2*4) // 3√3/2 · r²
	})
	t.Run("a square of semi-axis r encloses 2r²", func(t *testing.T) {
		AssertNumber(t, Square(Pt(0.0, 0.0), SzU(3.0), PointyTop).Area(), 18.0)
	})
	t.Run("an ellipse scales the area by both semi-axes", func(t *testing.T) {
		AssertNumber(t, Square(Pt(0.0, 0.0), Sz(2.0, 5.0), PointyTop).Area(), 20.0)
	})
	t.Run("fewer than three vertices enclose nothing", func(t *testing.T) {
		AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 0, 0).Area(), 0.0)
		AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 2, 0).Area(), 0.0)
	})
	t.Run("agrees with the polygon at any angle and size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertNumber(t, rp.Area(), rp.Polygon().Area(), rp.String())
		}
	})
}

func TestRegularPolygon_Perimeter(t *testing.T) {
	t.Run("agrees with the polygon", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(2.0), FlatTop)

		AssertNumber(t, hexagon.Perimeter(), hexagon.Polygon().Perimeter())
		AssertNumber(t, hexagon.Perimeter(), 12.0) // six edges of length r
	})
	t.Run("a square of semi-axis r has edges of r√2", func(t *testing.T) {
		AssertNumber(t, Square(Pt(0.0, 0.0), SzU(3.0), PointyTop).Perimeter(), 12*Sqrt2)
	})
	t.Run("an ellipse sums chords of different lengths", func(t *testing.T) {
		square := Square(Pt(0.0, 0.0), Sz(2.0, 5.0), PointyTop)

		AssertNumber(t, square.Perimeter(), 4*math.Hypot(2, 5))
		AssertNumber(t, square.Perimeter(), square.Polygon().Perimeter())
	})
	t.Run("fewer than two vertices have no edge", func(t *testing.T) {
		AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 0, 0).Perimeter(), 0.0)
		AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 1, 0).Perimeter(), 0.0)
		AssertNumber(t, RegPol(Pt(3.0, 4.0), SzU(10.0), 2, 0).Perimeter(), 40.0)
	})
	t.Run("int measures the exact polygon, not the rounded vertices", func(t *testing.T) {
		hexagon := Hexagon(Pt(0, 0), SzU(1), PointyTop)

		AssertNumber(t, hexagon.Perimeter(), 6.0)
		assert.NotEqual(t, hexagon.Polygon().Perimeter(), 6.0)
	})
	t.Run("agrees with the polygon at any angle and size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertNumber(t, rp.Perimeter(), rp.Polygon().Perimeter(), rp.String())
		}
	})
}

func TestRegularPolygon_Bounds(t *testing.T) {
	t.Run("vertices on the axes", func(t *testing.T) {
		AssertRectangle(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Bounds(), Rect(Pt(1, 2), Sz(4, 4)))
	})
	t.Run("hexagon is tight around its vertices", func(t *testing.T) {
		// width = 2r, height = √3 r
		bounds := Hexagon(Pt(0.0, 0.0), Sz(2.0, 2.0), FlatTop).Bounds()

		AssertNumber(t, bounds.Size.Width, 4.0)
		AssertNumber(t, bounds.Size.Height, 2.0*Sqrt3)
	})
	t.Run("one and two vertices reach only where they lie", func(t *testing.T) {
		AssertRectangle(t, RegPol(Pt(0.0, 0.0), SzU(2.0), 1, Pi/2).Bounds(), RectangleFromMinMax(Pt(0.0, 2.0), Pt(0.0, 2.0)))
		AssertRectangle(t, RegPol(Pt(0.0, 0.0), SzU(2.0), 2, Pi/4).Bounds(), RectangleFromMinMax(Pt(-Sqrt2, -Sqrt2), Pt(Sqrt2, Sqrt2)))
	})
	t.Run("no vertices is the zero rectangle", func(t *testing.T) {
		AssertRectangle(t, RegPol(Pt(3, 4), Sz(10, 10), 0, 0).Bounds(), Rectangle[int]{})
	})
	t.Run("is exactly the box around the vertices, rounded alike for int", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, angle := range []float64{0, 0.3, Pi / 3, 2, -1} {
				turned := rp.Rotate(angle)

				AssertRectangle(t, turned.Bounds(), turned.Polygon().Bounds(), turned.String())
				assert.Equal(t, turned.Int().Bounds(), turned.Int().Polygon().Bounds(), turned.String())
			}
		}
	})
}

func TestRegularPolygon_Translate(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Translate(Vec(1, -2)), RegPol(Pt(2, 0), Sz(2, 2), 4, 0))
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).MoveTo(Pt(-3, 5)), RegPol(Pt(-3, 5), Sz(2, 2), 4, 0))
}

func TestRegularPolygon_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Scale(0.5), RegPol(Pt(1, 2), Sz(1, 1), 4, 0))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).ScaleXY(2, 3), RegPol(Pt(1, 2), Sz(4, 6), 4, 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Scale(-2), RegPol(Pt(1, 2), Sz(4, 4), 4, 0))
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).ScaleXY(-2, 3), RegPol(Pt(1, 2), Sz(4, 6), 4, 0))
	})
}

func TestRegularPolygon_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 4), 4, 0).Unscale(2), RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 6), 4, 0).UnscaleXY(2, 3), RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(4, 4), 4, 0).Unscale(-2), RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			RegPol(Pt(1, 2), Sz(4, 4), 4, 0).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			RegPol(Pt(1, 2), Sz(4, 4), 4, 0).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestRegularPolygon_Canonical(t *testing.T) {
	t.Run("takes a literal negative size absolute and keeps the rest", func(t *testing.T) {
		AssertRegularPolygon(t, RegularPolygon[int]{Pt(1, 2), Sz(-2, 3), 6, Pi / 3}.Canonical(), RegPol(Pt(1, 2), Sz(2, 3), 6, Pi/3))
	})
	t.Run("is the polygon RegPol builds", func(t *testing.T) {
		rp := RegularPolygon[float64]{Pt(0.5, -1.25), Sz(-2.5, -3.75), 5, 1}

		AssertRegularPolygon(t, rp.Canonical(), RegPol(rp.Center, rp.Size, rp.N, rp.Angle))
		AssertVertices(t, slices.Collect(rp.Canonical().Vertices()), slices.Collect(RegPol(rp.Center, rp.Size, rp.N, rp.Angle).Vertices()))
	})
	t.Run("normalizes the angle the way Rotate stores it", func(t *testing.T) {
		AssertRegularPolygon(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, 7}.Canonical(), RegPol(Pt(0.0, 0.0), SzU(2.0), 4, 7-2*Pi))
		AssertNumber(t, RegularPolygon[float64]{Pt(0.0, 0.0), SzU(2.0), 4, -Pi / 2}.Canonical().Angle, 3*Pi/2)
	})
	t.Run("snaps a residue of turning back to exactly zero, where Rotate does not", func(t *testing.T) {
		drifted := RegPol(Pt(0.0, 0.0), SzU(2.0), 4, 0).Rotate(0.1).Rotate(0.2).Rotate(-0.3)

		assert.True(t, drifted.Angle != 0)
		assert.Equal(t, drifted.Canonical().Angle, 0.0)
		assert.Equal(t, RegularPolygon[int]{Pt(0, 0), SzU(2), 4, -Delta / 2}.Canonical().Angle, 0.0)
		assert.Equal(t, RegularPolygon[int]{Pt(0, 0), SzU(2), 4, 2 * Delta}.Canonical().Angle, 2*Delta)
	})
	t.Run("a well-formed polygon is unchanged", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertRegularPolygon(t, rp.Canonical(), rp, rp.String())
			AssertRegularPolygon(t, rp.Rotate(0).Canonical(), rp.Rotate(0), rp.String())
		}
	})
}

func TestRegularPolygon_Lerp(t *testing.T) {
	a, b := RegPol(Pt(0.0, 0.0), SzU(2.0), 6, 0), RegPol(Pt(10.0, 20.0), Sz(8.0, 4.0), 6, Pi/2)

	t.Run("moves the center, the size and the angle together", func(t *testing.T) {
		AssertRegularPolygon(t, a.Lerp(b, 0.5), RegPol(Pt(5.0, 10.0), Sz(5.0, 3.0), 6, Pi/4))
	})
	t.Run("turns along the shorter arc across the seam", func(t *testing.T) {
		from, to := RegPol(Pt(0, 0), SzU(2), 4, ToRadians(350)), RegPol(Pt(0, 0), SzU(2), 4, ToRadians(10))

		AssertRegularPolygon(t, from.Lerp(to, 0.5), RegPol(Pt(0, 0), SzU(2), 4, 0))
		AssertRegularPolygon(t, to.Lerp(from, 0.5), RegPol(Pt(0, 0), SzU(2), 4, 0))
	})
	t.Run("the ends are the polygons themselves", func(t *testing.T) {
		AssertRegularPolygon(t, a.Lerp(b, 0), a)
		AssertRegularPolygon(t, a.Lerp(b, 1), b)
	})
	t.Run("extrapolates and keeps the size absolute", func(t *testing.T) {
		AssertRegularPolygon(t, a.Lerp(b, 2), RegPol(Pt(20.0, 40.0), Sz(14.0, 6.0), 6, Pi))
		AssertRegularPolygon(t, b.Lerp(a, 2), RegPol(Pt(-10.0, -20.0), Sz(4.0, 0.0), 6, 3*Pi/2))
	})
	t.Run("int rounds the center and the size", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0, 0), SzU(2), 4, 0).Lerp(RegPol(Pt(5, 5), SzU(5), 4, 0), 0.5), RegPol(Pt(3, 3), SzU(4), 4, 0))
	})
	t.Run("a different vertex count panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			a.Lerp(RegPol(Pt(10.0, 20.0), Sz(8.0, 4.0), 5, Pi/2), 0.5)
		}, "geom: lerp between polygons of 6 and 5 vertices")
	})
	t.Run("the ends hold over the fixtures and the angle never turns more than half a turn", func(t *testing.T) {
		for _, from := range regularPolygonFixtures {
			for _, to := range regularPolygonFixtures {
				if from.N != to.N {
					continue
				}

				AssertRegularPolygon(t, from.Lerp(to, 0), from, fmt.Sprintf("%s → %s: ", from, to))
				AssertRegularPolygon(t, from.Lerp(to, 1), to, fmt.Sprintf("%s → %s: ", from, to))
				assert.True(t, AngleDistance(from.Lerp(to, 0.5).Angle, from.Angle) <= Pi/2+Delta, fmt.Sprintf("%s → %s: ", from, to))
			}
		}
	})
}

func TestRegularPolygon_Transform(t *testing.T) {
	hexagon := Hexagon(Pt(2.0, 3.0), SzU(4.0), FlatTop)

	t.Run("the identity keeps the polygon", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon.Transform(IdentityMatrix[float64]()), hexagon)
	})
	t.Run("a translation moves the center and keeps the vertex count", func(t *testing.T) {
		moved := hexagon.Transform(TranslationMatrix(1.0, -1.0))

		AssertRegularPolygon(t, moved, hexagon.Translate(Vec(1.0, -1.0)))
		assert.Equal(t, moved.N, 6)
	})
	t.Run("a uniform scale scales the semi-axes", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon.Transform(ScaleMatrix(2.0, 2.0)), Hexagon(Pt(4.0, 6.0), SzU(8.0), FlatTop))
	})
	t.Run("a rotation turns the angle", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon.Transform(RotationMatrix[float64](Pi/2)), Hexagon(Pt(-3.0, 2.0), SzU(4.0), FlatTop).Rotate(Pi/2))
	})
	t.Run("a reflection mirrors the angle about the axis of the matrix", func(t *testing.T) {
		turned := hexagon.Rotate(Pi / 6)

		AssertRegularPolygon(t, turned.Transform(ReflectionMatrix[float64](AxisHorizontal)), RegPol(Pt(2.0, -3.0), SzU(4.0), 6, -turned.Angle))
	})
	t.Run("a shear gives the nearest polygon", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon.Transform(ShearMatrix(1.0, 0.0)), Hexagon(Pt(5.0, 3.0), SzU(4.0), FlatTop))
	})
	t.Run("an empty polygon stays empty", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1.0, 1.0), SzU(2.0), 0, 0).Transform(ScaleMatrix(2.0, 2.0)).Empty())
	})
	t.Run("int rounds the center and the semi-axes once", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 1), SzU(3), 6, 0).Transform(ScaleMatrix(1.5, 1.5)), RegPol(Pt(2, 2), SzU(5), 6, 0))
	})
	t.Run("a turn of unequal semi-axes is exact too", func(t *testing.T) {
		ellipse := RegPol(Pt(0.0, 0.0), Sz(2.0, 4.0), 5, 0)
		turn := RotationMatrix[float64](Pi / 2)

		AssertRegularPolygon(t, ellipse.Transform(turn), RegPol(Pt(0.0, 0.0), Sz(2.0, 4.0), 5, Pi/2))
		AssertPolygon(t, ellipse.Transform(turn).Polygon(), ellipse.Polygon().Transform(turn))
	})
	t.Run("axes scaled by different factors turn the polygon into one it cannot hold", func(t *testing.T) {
		turned := RegPol(Pt(0.0, 0.0), SzU(4.0), 5, Pi/4)
		squeeze := ScaleMatrix(2.0, 3.0)

		assert.False(t, turned.Transform(squeeze).Polygon().Equal(turned.Polygon().Transform(squeeze)))
	})
	t.Run("matches the polygon of the vertices wherever the matrix keeps one", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			matrices := []Matrix[float64]{
				IdentityMatrix[float64](),
				TranslationMatrix(3.0, -2.0),
				ScaleMatrix(2.0, 2.0),
				RotationMatrix[float64](Pi / 3),
				RotationMatrix[float64](Pi / 3).Multiply(ScaleMatrix(2.0, 2.0)),
			}

			if rp.Angle == 0 {
				matrices = append(matrices, ScaleMatrix(2.0, 3.0))
			}

			for _, m := range matrices {
				AssertPolygon(t, rp.Transform(m).Polygon(), rp.Polygon().Transform(m), fmt.Sprintf("%s → %s: ", rp, m))
			}
		}
	})
}

func TestRegularPolygon_Rotate(t *testing.T) {
	t.Run("adds to the stored angle", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi))
	})
	t.Run("normalizes to the unit turn", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(3*Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi))      // 0 + 3π → π
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(-Pi/2), RegPol(Pt(1, 2), Sz(2, 2), 4, 3*Pi/2)) // 0 − π/2 → 3π/2
	})
}

func TestRegularPolygon_Contains(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)

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
		assert.True(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0).Contains(Pt(1.0, 1.0+Delta/2)))
		assert.False(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0).Contains(Pt(1.0, 1.0+2*Delta)))
	})
	t.Run("an empty polygon contains nothing", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0).Contains(Pt(0, 0)))
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
	hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop)
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
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)

	t.Run("zero within the polygon", func(t *testing.T) {
		assert.Equal(t, diamond.DistanceTo(Pt(0, 0)), 0.0)
		assert.Equal(t, diamond.DistanceTo(Pt(1, 1)), 0.0)
	})
	t.Run("to the nearest edge", func(t *testing.T) {
		AssertNumber(t, diamond.DistanceTo(Pt(2, 2)), Sqrt2)
	})
	t.Run("to the nearest vertex", func(t *testing.T) {
		AssertNumber(t, diamond.DistanceTo(Pt(3, 0)), 1.0)
	})
	t.Run("an empty polygon is infinitely far", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0).DistanceTo(Pt(0, 0)), math.Inf(1))
	})
	t.Run("matches the polygon of the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				AssertNumber(t, rp.DistanceTo(p), rp.Polygon().DistanceTo(p), fmt.Sprintf("%s → %s: ", rp, p))
			}
		}
	})
}

func TestRegularPolygon_DistanceSquaredTo(t *testing.T) {
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)

	t.Run("is the square of DistanceTo", func(t *testing.T) {
		AssertNumber(t, diamond.DistanceSquaredTo(Pt(2, 2)), 2.0)
		AssertNumber(t, diamond.DistanceSquaredTo(Pt(3, 0)), 1.0)
		assert.Equal(t, diamond.DistanceSquaredTo(Pt(0, 1)), 0.0)
	})
	t.Run("agrees with DistanceTo", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, p := range pointFixtures {
				AssertNumber(t, rp.DistanceSquaredTo(p), rp.DistanceTo(p)*rp.DistanceTo(p), fmt.Sprintf("%s → %s: ", rp, p))
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

func TestRegularPolygon_IntersectsLine(t *testing.T) {
	t.Run("mirrors Line.IntersectsRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, l := range lineFixtures {
				assert.Equal(t, rp.IntersectsLine(l), l.IntersectsRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, l))
			}
		}
	})
}

func TestRegularPolygon_IntersectionLine(t *testing.T) {
	t.Run("mirrors Line.IntersectionRegularPolygon", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, l := range lineFixtures {
				AssertVertices(t, rp.IntersectionLine(l), l.IntersectionRegularPolygon(rp), fmt.Sprintf("%s → %s: ", rp, l))
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
	diamond := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)

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
		assert.True(t, diamond.IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(1, 1), 3, 0)))
		assert.True(t, RegPol(Pt(0, 0), Sz(1, 1), 3, 0).IntersectsRegularPolygon(diamond))
	})
	t.Run("edges crossing without a vertex inside", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0).IntersectsRegularPolygon(RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, Pi/4)))
	})
	t.Run("an empty polygon intersects nothing", func(t *testing.T) {
		assert.False(t, diamond.IntersectsRegularPolygon(RegPol(Pt(0, 0), Sz(2, 2), 0, 0)))
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0).IntersectsRegularPolygon(diamond))
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

func TestRegularPolygon_Equal(t *testing.T) {
	t.Run("same polygon", func(t *testing.T) {
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 0)))
	})
	t.Run("different polygon", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(0, 0), Sz(2, 2), 6, 0)))
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(1, 2), Sz(3, 3), 4, 0)))
	})
	t.Run("the angle is compared", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, Pi)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0.5).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 0.5)))
	})
	t.Run("the angle is compared normalized", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), PointyTop)

		assert.True(t, hexagon.Equal(hexagon.Rotate(0)))
		assert.True(t, hexagon.Equal(hexagon.Rotate(2*Pi)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, -Pi/2).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 3*Pi/2)))
	})
	t.Run("the angle is compared across the seam", func(t *testing.T) {
		hexagon := Hexagon(Pt(0.0, 0.0), SzU(10.0), PointyTop)

		assert.True(t, hexagon.Equal(hexagon.Rotate(-1e-9)))
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 2*Pi-1e-9)))
	})
}

func TestRegularPolygon_IsZero(t *testing.T) {
	t.Run("zero polygon", func(t *testing.T) {
		assert.True(t, RegularPolygon[int]{}.IsZero())
		assert.True(t, RegularPolygon[float64]{}.IsZero())
	})
	t.Run("a full turn is a zero angle, like Equal", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0, 0), Sz(0, 0), 0, 2*Pi).IsZero())
	})
	t.Run("only one field is zero", func(t *testing.T) {
		assert.False(t, RegPol(Pt(0, 0), Sz(0, 0), 4, 0).IsZero())
		assert.False(t, RegPol(Pt(1, 2), Sz(0, 0), 0, 0).IsZero())
		assert.False(t, RegPol(Pt(0, 0), Sz(2, 2), 0, 0).IsZero())
		assert.False(t, RegPol(Pt(0, 0), Sz(0, 0), 0, 1).IsZero())
	})
	t.Run("non-zero polygon", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).IsZero())
	})
}

func TestRegularPolygon_Empty(t *testing.T) {
	t.Run("no sides", func(t *testing.T) {
		assert.True(t, RegularPolygon[int]{}.Empty())
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 0, 0).Empty())
		assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), -1, 0).Empty())
	})
	t.Run("with sides", func(t *testing.T) {
		assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Empty())
	})
}

func TestRegularPolygon_Polygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0)
	p := rp.Polygon()

	t.Run("carries the vertices", func(t *testing.T) {
		AssertVertices(t, p.Points, slices.Collect(rp.Vertices()))
	})
	t.Run("owns its slice", func(t *testing.T) {
		assert.NotSame(t, p.Points, rp.Polygon().Points)
	})
	t.Run("fewer than one side is the zero polygon", func(t *testing.T) {
		assert.True(t, RegPol(Pt(0, 0), Sz(1, 1), 0, 0).Polygon().IsZero())
		assert.True(t, RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0).Polygon().IsZero())
	})
}

func TestRegularPolygon_Cast(t *testing.T) {
	rp := RegPol(Pt(1.5, -2.5), Sz(3.5, 4.5), 6, Pi/6)

	t.Run("matches Int and Float", func(t *testing.T) {
		AssertRegularPolygon(t, rp.Cast[int](), rp.Int())
		AssertRegularPolygon(t, rp.Cast[float64](), rp.Float())
	})
	t.Run("a type the other conversions cannot name, and N and the angle are kept", func(t *testing.T) {
		AssertRegularPolygon(t, rp.Cast[int8](), RegPol(Pt[int8](2, -3), Sz[int8](4, 5), 6, Pi/6))
	})
}

func TestRegularPolygon_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Int(), RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	})
	t.Run("float rounds but keeps the angle", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3).Int(), RegPol(Pt(1, 0), Sz(1, 4), 6, Pi/3))
	})
}

func TestRegularPolygon_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Float(), RegPol(Pt(1.0, 2.0), Sz(2.0, 2.0), 4, 0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertRegularPolygon(t, RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3).Float(), RegPol(Pt(0.6, -0.25), Sz(1.2, 3.6), 6, Pi/3))
	})
}

func TestRegularPolygon_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).String(), "RegPol((1,2);2x2;4)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, Pi).String(), "RegPol((0.50,-1.25);2.50x3.75;6;3.14)")
	})
}

func TestRegularPolygon_JSON(t *testing.T) {
	t.Run("int wire format omits a zero angle", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0), `{"x":1,"y":2,"w":2,"h":2,"n":4}`)

		var rp RegularPolygon[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4}`), &rp))
		AssertRegularPolygon(t, rp, RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
		assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`), &rp))
		AssertRegularPolygon(t, rp, RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, 0.5), `{"x":0.50,"y":-1.25,"w":2.50,"h":3.75,"n":6,"a":0.5}`)

		var rp RegularPolygon[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":0.50,"y":-1.25,"w":2.50,"h":3.75,"n":6,"a":0.5}`), &rp))
		AssertRegularPolygon(t, rp, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, 0.5))
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
}

func TestRegularPolygon_Properties(t *testing.T) {
	t.Run("vertex count matches the side count", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			assert.Equal(t, len(slices.Collect(rp.Vertices())), rp.N, fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("vertices lie on the ellipse of the size, in the frame before the turn", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for vertex := range rp.Vertices() {
				offset := vertex.Subtract(rp.Center).Rotate(-rp.Angle)
				normalized := offset.X*offset.X/(rp.Size.Width*rp.Size.Width) + offset.Y*offset.Y/(rp.Size.Height*rp.Size.Height)

				AssertNumber(t, normalized, 1.0, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("a full turn restores the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertVertices(t, slices.Collect(rp.Rotate(2*Pi).Vertices()), slices.Collect(rp.Vertices()), fmt.Sprintf("%s: ", rp))
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
				AssertNumber(t, moved.Angle, rp.Angle, fmt.Sprintf("%s → %s: ", rp, vector))
			}
		}
	})
	t.Run("bounds enclose every vertex", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			bounds := rp.Bounds()

			for vertex := range rp.Vertices() {
				assert.True(t, vertex.X >= bounds.Min().X-Delta && vertex.X <= bounds.Max().X+Delta, fmt.Sprintf("%s: ", rp))
				assert.True(t, vertex.Y >= bounds.Min().Y-Delta && vertex.Y <= bounds.Max().Y+Delta, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("polygon carries the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertVertices(t, rp.Polygon().Points, slices.Collect(rp.Vertices()), fmt.Sprintf("%s: ", rp))
		}
	})
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RegPol(Pt(0, 1), Sz(2, 3), 4, 0)

	rp.Translate(Vec(2, 3))
	rp.MoveTo(Pt(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)
	rp.Rotate(Pi / 3)

	AssertRegularPolygon(t, rp, RegPol(Pt(0, 1), Sz(2, 3), 4, 0))
}

// regularPolygonFixtures span the triangle, square, and hexagon at both orientations.
var regularPolygonFixtures = []RegularPolygon[float64]{
	Triangle(Pt(0.0, 0.0), Sz(3.0, 3.0), PointyTop),
	Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointyTop),
	Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), FlatTop),
	RegPol(Pt(-3.5, 0.25), Sz(2.0, 3.0), 5, Pi/7),
	RegPol(Pt(12.5, -0.1), Sz(0.5, 0.5), 8, 0),
}

func ExampleRegPol() {
	fmt.Println(RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
	// Output: RegPol((1,2);2x2;4)
}

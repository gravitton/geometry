package geom

import (
	"encoding/json"
	"fmt"
	"math"
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
			vertices := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), n, FlatTop).Vertices()

			AssertPoint(t, vertices[0].Midpoint(vertices[1]), Pt(0.0, -10*math.Cos(Pi/float64(n))), fmt.Sprintf("n=%d: ", n))
		}
	})
	t.Run("unknown orientation panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			RegularPolygonOrientationAngle(6, Orientation(99))
		}, "geom: unknown orientation 99")
	})
	t.Run("no vertices give the top angle instead of dividing by n", func(t *testing.T) {
		AssertNumber(t, RegularPolygonOrientationAngle(0, FlatTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(-1, FlatTop), 270*DegToRad)
		AssertNumber(t, RegularPolygonOrientationAngle(0, PointyTop), 270*DegToRad)

		empty := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 0, FlatTop)
		assert.True(t, empty.Equal(empty))
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
		assert.True(t, pointy.Vertices()[0].Y < center.Y)
	})
	t.Run("flat top starts left of and above the center", func(t *testing.T) {
		first := flat.Vertices()[0]

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
		AssertVertices(t, triangle.Vertices(), []Point[int]{Pt(1, -4), Pt(4, 0), Pt(-2, 0)})
	})
}

func TestSquare(t *testing.T) {
	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointyTop)

	t.Run("has four sides at the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, square, RegPol(Pt(50.0, 50.0), Sz(100.0, 100.0), 4, RegularPolygonOrientationAngle(4, PointyTop)))
	})
	t.Run("pointy top starts at the visual top and winds to the right", func(t *testing.T) {
		AssertVertices(t, square.Vertices(), []Point[float64]{Pt(50.0, -50.0), Pt(150.0, 50.0), Pt(50.0, 150.0), Pt(-50.0, 50.0)})
	})
}

func TestHexagon(t *testing.T) {
	// float64 avoids the int-rounding collapse where sin(±π/6) ≈ 0.4999 would truncate to 0
	hexagon := Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), PointyTop)

	t.Run("has six sides at the orientation angle", func(t *testing.T) {
		AssertRegularPolygon(t, hexagon, RegPol(Pt(0.0, 0.0), Sz(10.0, 10.0), 6, RegularPolygonOrientationAngle(6, PointyTop)))
	})
	t.Run("pointy top has vertices at top and bottom and flat sides left and right", func(t *testing.T) {
		AssertVertices(t, hexagon.Vertices(), []Point[float64]{
			Pt(0.0, -10.0),
			Pt(5*Sqrt3, -5.0),
			Pt(5*Sqrt3, 5.0),
			Pt(0.0, 10.0),
			Pt(-5*Sqrt3, 5.0),
			Pt(-5*Sqrt3, -5.0),
		})
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

func TestRegularPolygon_Vertices(t *testing.T) {
	t.Run("fewer than one side has nil vertices", func(t *testing.T) {
		assert.Nil(t, RegPol(Pt(0, 0), Sz(1, 1), 0, 0).Vertices())
		assert.Nil(t, RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), -3, 0).Vertices())
		assert.True(t, RegPol(Pt(0, 0), Sz(1, 1), 0, 0).Polygon().IsZero())
	})
	t.Run("int", func(t *testing.T) {
		AssertVertices(t, RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices(), []Point[int]{
			Pt(1, 0),
			Pt(0, 1),
			Pt(-1, 0),
			Pt(0, -1),
		})
	})
	t.Run("non-square size stretches the ellipse", func(t *testing.T) {
		AssertVertices(t, RegPol(Pt(0, 0), Sz(2, 3), 4, 0).Vertices(), []Point[int]{
			Pt(2, 0),
			Pt(0, 3),
			Pt(-2, 0),
			Pt(0, -3),
		})
	})
	t.Run("float", func(t *testing.T) {
		AssertVertices(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0).Vertices(), []Point[float64]{
			Pt(2.0, 0.0),
			Pt(1.0, 1.5*Sqrt3),
			Pt(-1.0, 1.5*Sqrt3),
			Pt(-2.0, 0.0),
			Pt(-1.0, -1.5*Sqrt3),
			Pt(1.0, -1.5*Sqrt3),
		})
	})
	t.Run("no sides means no vertices", func(t *testing.T) {
		AssertVertices(t, RegPol(Pt(3, 4), Sz(10, 10), 0, 0).Vertices(), []Point[int]{})
	})
}

func TestRegularPolygon_Bounds(t *testing.T) {
	t.Run("vertices on the axes", func(t *testing.T) {
		AssertRect(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Bounds(), Rect(Pt(1, 2), Sz(4, 4)))
	})
	t.Run("hexagon is tight around its vertices", func(t *testing.T) {
		// width = 2r, height = √3 r
		bounds := Hexagon(Pt(0.0, 0.0), Sz(2.0, 2.0), FlatTop).Bounds()

		AssertNumber(t, bounds.Size.Width, 4.0)
		AssertNumber(t, bounds.Size.Height, 2.0*Sqrt3)
	})
	t.Run("no vertices is the zero rectangle", func(t *testing.T) {
		AssertRect(t, RegPol(Pt(3, 4), Sz(10, 10), 0, 0).Bounds(), Rectangle[int]{})
	})
}

func TestRegularPolygon_Polygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0)
	p := rp.Polygon()

	t.Run("carries the vertices", func(t *testing.T) {
		AssertVertices(t, p.Vertices, rp.Vertices())
	})
	t.Run("owns its slice", func(t *testing.T) {
		assert.NotSame(t, p.Vertices, rp.Vertices())
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
		assert.Equal(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).String(), "RegPol((1,2);2x2;4;0.00)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, RegPol(Pt(0.5, -1.25), Sz(2.5, 3.75), 6, Pi).String(), "RegPol((0.50,-1.25);2.50x3.75;6;3.14)")
	})
}

func TestRegularPolygon_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0), `{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`)

		var rp RegularPolygon[int]
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
			assert.Equal(t, len(rp.Vertices()), rp.N, fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("vertices lie on the ellipse of the size", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			for _, vertex := range rp.Vertices() {
				offset := vertex.Subtract(rp.Center)
				normalized := offset.X*offset.X/(rp.Size.Width*rp.Size.Width) + offset.Y*offset.Y/(rp.Size.Height*rp.Size.Height)

				AssertNumber(t, normalized, 1.0, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("a full turn restores the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertVertices(t, rp.Rotate(2*Pi).Vertices(), rp.Vertices(), fmt.Sprintf("%s: ", rp))
		}
	})
	t.Run("rotating by one step permutes the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			rotated := rp.Rotate(2 * Pi / float64(rp.N)).Vertices()
			vertices := rp.Vertices()

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

			for _, vertex := range rp.Vertices() {
				assert.True(t, vertex.X >= bounds.Min().X-Delta && vertex.X <= bounds.Max().X+Delta, fmt.Sprintf("%s: ", rp))
				assert.True(t, vertex.Y >= bounds.Min().Y-Delta && vertex.Y <= bounds.Max().Y+Delta, fmt.Sprintf("%s: ", rp))
			}
		}
	})
	t.Run("polygon carries the vertices", func(t *testing.T) {
		for _, rp := range regularPolygonFixtures {
			AssertVertices(t, rp.Polygon().Vertices, rp.Vertices(), fmt.Sprintf("%s: ", rp))
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
	// Output: RegPol((1,2);2x2;4;0.00)
}

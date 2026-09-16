package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestRegularPolygon_Constructor(t *testing.T) {
	rp := RegPol(Pt(0, 0), Sz(2, 2), 4, 0)
	AssertRegularPolygon(t, rp, RegularPolygon[int]{Center: Pt(0, 0), Size: Sz(2, 2), N: 4, Angle: 0})

	triangle := Triangle(Pt(1, -1), Sz(3, 3), PointyTop)
	AssertRegularPolygon(t, triangle, RegPol(Pt(1, -1), Sz(3, 3), 3, RegularPolygonOrientationAngle(3, PointyTop)))
	// Integer rounding happens after scaling by the size, so vertices 1 and 2 land within one
	// unit of the exact (3.598, 0.5) and (-1.598, 0.5); their Y of 0.5 falls just below the
	// .5 tie in float64 and rounds down to 0 and up to 1 respectively.
	AssertVertices(t, triangle.Vertices(), []Point[int]{Pt(1, -4), Pt(4, 0), Pt(-2, 1)})

	square := Square(Pt(50.0, 50.0), Sz(100.0, 100.0), PointyTop)
	AssertRegularPolygon(t, square, RegPol(Pt(50.0, 50.0), Sz(100.0, 100.0), 4, RegularPolygonOrientationAngle(4, PointyTop)))
	// PointyTop: first vertex at visual top (Y = center.Y − size = −50), last at visual right.
	AssertVertices(t, square.Vertices(), []Point[float64]{Pt(50.0, -50.0), Pt(150.0, 50.0), Pt(50.0, 150.0), Pt(-50.0, 50.0)})

	// float64 avoids the int-rounding collapse where sin(±π/6) ≈ 0.4999 would truncate to 0.
	hexagon := Hexagon(Pt(0.0, 0.0), Sz(10.0, 10.0), PointyTop)
	AssertRegularPolygon(t, hexagon, RegPol(Pt(0.0, 0.0), Sz(10.0, 10.0), 6, RegularPolygonOrientationAngle(6, PointyTop)))
	// PointyTop: V0 at top (0,−10), V3 at bottom (0,10); flat sides left and right.
	AssertVertices(t, hexagon.Vertices(), []Point[float64]{
		Pt(0.0, -10.0),
		Pt(5*Sqrt3, -5.0),
		Pt(5*Sqrt3, 5.0),
		Pt(0.0, 10.0),
		Pt(-5*Sqrt3, 5.0),
		Pt(-5*Sqrt3, -5.0),
	})
}

func TestRegularPolygon_OrientationAngle(t *testing.T) {
	// In +Y-down screen coordinates, "top" means minimum Y, so PointyTop requires the first
	// vertex to point in the -Y direction: angle = -π/2.
	assert.EqualDelta(t, RegularPolygonOrientationAngle(3, PointyTop), -90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonOrientationAngle(3, FlatTop), 30*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonOrientationAngle(4, PointyTop), -90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonOrientationAngle(4, FlatTop), 45*DegToRad, Delta)

	assert.EqualDelta(t, RegularPolygonOrientationAngle(6, PointyTop), -90*DegToRad, Delta)
	assert.EqualDelta(t, RegularPolygonOrientationAngle(6, FlatTop), 60*DegToRad, Delta)

	// an orientation outside the two constants has no initial angle
	assert.EqualDelta(t, RegularPolygonOrientationAngle(6, Orientation(99)), 0.0, Delta)
}

func TestRegularPolygon_WithOrientation(t *testing.T) {
	center := Pt(0, 0)
	size := Sz(10, 10)

	rp := RegularPolygonWithOrientation(center, size, 6, PointyTop)
	AssertRegularPolygon(t, rp, RegPol(Pt(0, 0), Sz(10, 10), 6, RegularPolygonOrientationAngle(6, PointyTop)))

	rpFlat := RegularPolygonWithOrientation(center, size, 6, FlatTop)
	AssertRegularPolygon(t, rpFlat, RegPol(Pt(0, 0), Sz(10, 10), 6, RegularPolygonOrientationAngle(6, FlatTop)))

	// angle differs between orientations
	assert.NotEqual(t, rp.Angle, rpFlat.Angle)

	// In +Y-down coordinates, PointyTop places the first vertex above the center (Y < center.Y).
	assert.True(t, rp.Vertices()[0].Y < center.Y)

	// FlatTop first vertex is to the lower-right of center (Y > center.Y).
	assert.True(t, rpFlat.Vertices()[0].Y > center.Y)
}

func TestRegularPolygon_Translate(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Translate(Vec(1, -2)), RegPol(Pt(2, 0), Sz(2, 2), 4, 0))
}

func TestRegularPolygon_MoveTo(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).MoveTo(Pt(-3, 5)), RegPol(Pt(-3, 5), Sz(2, 2), 4, 0))
}

func TestRegularPolygon_Scale(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Scale(0.5), RegPol(Pt(1, 2), Sz(1, 1), 4, 0))
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).ScaleXY(2, 3), RegPol(Pt(1, 2), Sz(4, 6), 4, 0))
}

func TestRegularPolygon_Rotate(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi))
	// Normalization: angle wraps to [0, 2π).
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(3*Pi), RegPol(Pt(1, 2), Sz(2, 2), 4, Pi))      // 0 + 3π → π
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Rotate(-Pi/2), RegPol(Pt(1, 2), Sz(2, 2), 4, 3*Pi/2)) // 0 − π/2 → 3π/2
}

func TestRegularPolygon_Vertices(t *testing.T) {
	AssertVertices(t, RegPol(Pt(0, 0), Sz(1, 1), 4, 0).Vertices(), []Point[int]{
		Pt(1, 0),
		Pt(0, 1),
		Pt(-1, 0),
		Pt(0, -1),
	})
	AssertVertices(t, RegPol(Pt(0, 0), Sz(2, 3), 4, 0).Vertices(), []Point[int]{
		Pt(2, 0),
		Pt(0, 3),
		Pt(-2, 0),
		Pt(0, -3),
	})
	AssertVertices(t, RegPol(Pt(0.0, 0.0), Sz(2.0, 3.0), 6, 0).Vertices(), []Point[float64]{
		Pt(2.0, 0.0),
		Pt(1.0, 1.5*Sqrt3),
		Pt(-1.0, 1.5*Sqrt3),
		Pt(-2.0, 0.0),
		Pt(-1.0, -1.5*Sqrt3),
		Pt(1.0, -1.5*Sqrt3),
	})
}

func TestRegularPolygon_Bounds(t *testing.T) {
	// Square-like (N=4, angle=0): vertices on axes, so bounds == 2×size centered.
	AssertRect(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Bounds(), Rect(Pt(1, 2), Sz(4, 4)))

	// Hexagon (float64): tight bounds from actual vertices; width = 2r, height = √3 * r.
	hex := Hexagon(Pt(0.0, 0.0), Sz(2.0, 2.0), FlatTop)
	b := hex.Bounds()
	assert.EqualDelta(t, float64(b.Size.Width), 4.0, Delta)
	assert.EqualDelta(t, float64(b.Size.Height), 2.0*Sqrt3, Delta)

	// no vertices: the bounds collapse to the center with zero size
	AssertRect(t, RegPol(Pt(3, 4), Sz(10, 10), 0, 0).Bounds(), Rect(Pt(3, 4), Sz(0, 0)))
}

func TestRegularPolygon_Polygon(t *testing.T) {
	rp := RegPol(Pt(0.0, 0.0), Sz(1.0, 1.0), 5, 0)
	p := rp.Polygon()

	assert.Equal(t, p.Vertices, rp.Vertices())
	assert.NotSame(t, p.Vertices, rp.Vertices())
}

func TestRegularPolygon_Equal(t *testing.T) {
	assert.True(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(1, 2), Sz(2, 2), 4, 0)))
	assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Equal(RegPol(Pt(0, 0), Sz(2, 2), 6, 0)))
}

func TestRegularPolygon_IsZero(t *testing.T) {
	assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).IsZero())
	assert.True(t, RegularPolygon[int]{}.Empty())
}

func TestRegularPolygon_Empty(t *testing.T) {
	assert.False(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Empty())
	assert.True(t, Polygon[int]{}.Empty())
	assert.True(t, Polygon[int]{[]Point[int]{}}.Empty())
	assert.False(t, Pol([]Point[float64]{Pt(0.0, 0.0), Pt(2.5, 0.5), Pt(2.0, 1.0)}).Empty())
}

func TestRegularPolygon_Int(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Int(), RegPol(Pt(1, 2), Sz(2, 2), 4, 0.0))
}

func TestRegularPolygon_Float(t *testing.T) {
	AssertRegularPolygon(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).Float(), RegPol(Pt(1.0, 2.0), Sz(2.0, 2.0), 4, 0.0))
}

func TestRegularPolygon_Immutable(t *testing.T) {
	rp := RegPol(Pt(0, 1), Sz(2, 3), 4, 0)

	rp.Translate(Vec(2, 3))
	rp.MoveTo(Pt(1, 1))
	rp.Scale(2)
	rp.ScaleXY(2, 3)

	AssertRegularPolygon(t, rp, RegPol(Pt(0, 1), Sz(2, 3), 4, 0))
}

func TestRegularPolygon_String(t *testing.T) {
	assert.Equal(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0).String(), "RegPol((1,2);2x2;4;0.00)")
}

func TestRegularPolygon_Marshall(t *testing.T) {
	assert.JSON(t, RegPol(Pt(1, 2), Sz(2, 2), 4, 0), `{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`)
}

func TestRegularPolygon_Unmarshall(t *testing.T) {
	var p1 RegularPolygon[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":1,"y":2,"w":2,"h":2,"n":4,"a":0}`), &p1))
	AssertRegularPolygon(t, p1, RegPol(Pt(1, 2), Sz(2, 2), 4, 0))
}

package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestEllipse_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(10, 16), Sz(12, 4), 0), Ellipse[int]{Center: Pt(10, 16), Size: Sz(12, 4)})
	})
	t.Run("float", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(0.16, 204), Sz(5.1, 1.2), 0.5), Ellipse[float64]{Center: Pt(0.16, 204.0), Size: Sz(5.1, 1.2), Angle: 0.5})
	})
	t.Run("a negative semi-axis is taken absolute", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(10, 16), Sz(-12, -4), 0), Ell(Pt(10, 16), Sz(12, 4), 0))
	})
}

func TestEllipse_SemiMajor(t *testing.T) {
	t.Run("the wider axis", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0, 0), Sz(5, 3), 0).SemiMajor(), 5)
	})
	t.Run("the taller axis", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0, 0), Sz(3, 5), 0).SemiMajor(), 5)
	})
}

func TestEllipse_SemiMinor(t *testing.T) {
	t.Run("the wider axis", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0, 0), Sz(5, 3), 0).SemiMinor(), 3)
	})
	t.Run("the taller axis", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0, 0), Sz(3, 5), 0).SemiMinor(), 3)
	})
}

func TestEllipse_Eccentricity(t *testing.T) {
	t.Run("a circle is not eccentric", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.0, 0.0), SzU(4.0), 0).Eccentricity(), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0).Eccentricity(), 0.8)
	})
	t.Run("a degenerate ellipse is fully eccentric", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.0, 0.0), Sz(5.0, 0.0), 0).Eccentricity(), 1.0)
	})
	t.Run("an ellipse of no extent is a point", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.0, 0.0), Sz(0.0, 0.0), 0).Eccentricity(), 0.0)
	})
}

func TestEllipse_Foci(t *testing.T) {
	t.Run("on the horizontal major axis", func(t *testing.T) {
		a, b := Ell(Pt(1.0, 2.0), Sz(5.0, 3.0), 0).Foci()

		AssertPoint(t, a, Pt(-3.0, 2.0))
		AssertPoint(t, b, Pt(5.0, 2.0))
	})
	t.Run("on the vertical major axis", func(t *testing.T) {
		a, b := Ell(Pt(1.0, 2.0), Sz(3.0, 5.0), 0).Foci()

		AssertPoint(t, a, Pt(1.0, -2.0))
		AssertPoint(t, b, Pt(1.0, 6.0))
	})
	t.Run("a circle has both at the center", func(t *testing.T) {
		a, b := Ell(Pt(1.0, 2.0), SzU(3.0), 0).Foci()

		AssertPoint(t, a, Pt(1.0, 2.0))
		AssertPoint(t, b, Pt(1.0, 2.0))
	})
	t.Run("turned with the ellipse", func(t *testing.T) {
		a, b := Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), Pi/2).Foci()

		AssertPoint(t, a, Pt(0.0, -4.0))
		AssertPoint(t, b, Pt(0.0, 4.0))
	})
}

func TestEllipse_Anchor(t *testing.T) {
	e := Ell(Pt(10.0, 10.0), Sz(5.0, 3.0), 0)

	t.Run("cardinal directions end the semi-axes", func(t *testing.T) {
		AssertPoint(t, e.Anchor(Right), Pt(15.0, 10.0))
		AssertPoint(t, e.Anchor(Left), Pt(5.0, 10.0))
		AssertPoint(t, e.Anchor(Top), Pt(10.0, 7.0))
		AssertPoint(t, e.Anchor(Bottom), Pt(10.0, 13.0))
	})
	t.Run("diagonals land on the boundary in their direction", func(t *testing.T) {
		reach := 15 / math.Sqrt(34)

		AssertPoint(t, e.Anchor(DirectionDownRight), Pt(10+reach, 10+reach))
		AssertPoint(t, e.Anchor(DirectionUpLeft), Pt(10-reach, 10-reach))
	})
	t.Run("a circle anchors like Circle", func(t *testing.T) {
		AssertPoint(t, Ell(Pt(10.0, 10.0), SzU(5.0), 0).Anchor(DirectionDownRight), Circ(Pt(10.0, 10.0), 5.0).Anchor(DirectionDownRight))
	})
	t.Run("a degenerate ellipse anchors on its segment", func(t *testing.T) {
		AssertPoint(t, Ell(Pt(10.0, 10.0), Sz(0.0, 3.0), 0).Anchor(Top), Pt(10.0, 7.0))
		AssertPoint(t, Ell(Pt(10.0, 10.0), Sz(0.0, 3.0), 0).Anchor(DirectionDownRight), Pt(10.0, 10+3*OneOverSqrt2))
	})
	t.Run("none is the center", func(t *testing.T) {
		AssertPoint(t, e.Anchor(DirectionNone), Pt(10.0, 10.0))
	})
	t.Run("named in the frame before the turn", func(t *testing.T) {
		AssertPoint(t, e.Rotate(Pi/2).Anchor(Right), Pt(10.0, 15.0))
	})
}

func TestEllipse_Centroid(t *testing.T) {
	AssertPoint(t, Ell(Pt(1, 2), Sz(10, 4), 0).Centroid(), Pt(1, 2))
}

func TestEllipse_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1, 2), Sz(10, 4), 0).Area(), Pi*40.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0).Area(), Pi*0.6)
	})
	t.Run("the turn does not change it", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1, 2), Sz(10, 4), 1.0).Area(), Pi*40.0)
	})
}

func TestEllipse_Inertia(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1, 2), Sz(10, 4), 0).Inertia(), Pi*1160.0)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0).Inertia(), Pi*0.6*1.69/4)
	})
	t.Run("the turn does not change it", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1, 2), Sz(10, 4), 1.0).Inertia(), Pi*1160.0)
	})
	t.Run("a circle agrees with Circle", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1.0, 2.0), SzU(10.0), 0).Inertia(), Circ(Pt(1.0, 2.0), 10.0).Inertia())
	})
	t.Run("is approached by the polygon", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			assert.EqualDelta(t, e.Inertia(), e.RegularPolygon(360).Inertia(), e.Inertia()*1e-3, e.String())
		}
	})
}

func TestEllipse_Perimeter(t *testing.T) {
	t.Run("a circle is exact", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1.0, 2.0), SzU(10.0), 0).Perimeter(), 2*Pi*10)
	})
	t.Run("float", func(t *testing.T) {
		// Ramanujan's second approximation, within a part in 1e9 of the true 25.527
		assert.EqualDelta(t, Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0).Perimeter(), 25.526998, 1e-5)
	})
	t.Run("a degenerate ellipse is twice its segment", func(t *testing.T) {
		assert.EqualDelta(t, Ell(Pt(0.0, 0.0), Sz(5.0, 0.0), 0).Perimeter(), 20.0, 1e-2)
	})
	t.Run("an ellipse of no extent has no boundary", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(1.0, 2.0), Sz(0.0, 0.0), 0).Perimeter(), 0.0)
	})
}

func TestEllipse_Bounds(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRectangle(t, Ell(Pt(1, 2), Sz(10, 4), 0).Bounds(), Rect(Pt(1, 2), Sz(20, 8)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0).Bounds(), Rect(Pt(0.6, -0.25), Sz(2.4, 1.0)))
	})
	t.Run("a quarter turn transposes it", func(t *testing.T) {
		AssertRectangle(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), Pi/2).Bounds(), Rect(Pt(1.0, 2.0), Sz(8.0, 20.0)))
	})
	t.Run("a turned ellipse reaches where its tangent is axis-aligned", func(t *testing.T) {
		// not the corner of the turned box: √(w²cos² + h²sin²) either side
		AssertRectangle(t, Ell(Pt(0.0, 0.0), Sz(10.0, 4.0), Pi/6).Bounds(), Rect(Pt(0.0, 0.0), Sz(2*math.Hypot(10*math.Cos(Pi/6), 4*math.Sin(Pi/6)), 2*math.Hypot(10*math.Sin(Pi/6), 4*math.Cos(Pi/6)))))
	})
}

func TestEllipse_minMax(t *testing.T) {
	t.Run("spans the extent", func(t *testing.T) {
		a, b := Ell(Pt(1, 2), Sz(10, 4), 0).minMax()

		AssertPoint(t, a, Pt(-9, -2))
		AssertPoint(t, b, Pt(11, 6))
	})
	t.Run("matches the corners of Bounds", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			a, b := e.minMax()
			c, d := e.Bounds().MinMax()

			AssertPoint(t, a, c, e.String())
			AssertPoint(t, b, d, e.String())
		}
	})
}

func TestEllipse_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Translate(Vec(3, -2)), Ell(Pt(4, 0), Sz(10, 4), 0))
	})
	t.Run("the turn is kept", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5).Translate(Vec(3.0, -2.0)), Ell(Pt(4.0, 0.0), Sz(10.0, 4.0), 0.5))
	})
}

func TestEllipse_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).MoveTo(Pt(3, -2)), Ell(Pt(3, -2), Sz(10, 4), 0))
	})
	t.Run("the turn is kept", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5).MoveTo(Pt(3.0, -2.0)), Ell(Pt(3.0, -2.0), Sz(10.0, 4.0), 0.5))
	})
}

func TestEllipse_Scale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Scale(2.5), Ell(Pt(1, 2), Sz(25, 10), 0))
	})
	t.Run("float", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0).Scale(2.5), Ell(Pt(0.6, -0.25), Sz(3.0, 1.25), 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Scale(-2), Ell(Pt(1, 2), Sz(20, 8), 0))
	})
}

func TestEllipse_ScaleXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).ScaleXY(2, 0.5), Ell(Pt(1, 2), Sz(20, 2), 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).ScaleXY(-2, 1), Ell(Pt(1, 2), Sz(20, 4), 0))
	})
}

func TestEllipse_Unscale(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(25, 10), 0).Unscale(2.5), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
	t.Run("a negative factor scales by its absolute value", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(20, 8), 0).Unscale(-2), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
	t.Run("a zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Ell(Pt(1, 2), Sz(10, 4), 0).Unscale(0)
		}, "geom: division by zero")
	})
}

func TestEllipse_UnscaleXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(20, 2), 0).UnscaleXY(2, 0.5), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
	t.Run("a zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Ell(Pt(1, 2), Sz(10, 4), 0).UnscaleXY(1, 0)
		}, "geom: division by zero")
	})
}

func TestEllipse_Resize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Resize(Sz(3, 7)), Ell(Pt(1, 2), Sz(3, 7), 0))
	})
	t.Run("a negative semi-axis is taken absolute", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Resize(Sz(-3, 7)), Ell(Pt(1, 2), Sz(3, 7), 0))
	})
}

func TestEllipse_Canonical(t *testing.T) {
	t.Run("a well-formed ellipse is unchanged", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 1.0).Canonical(), Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 1.0))
	})
	t.Run("a negative semi-axis is repaired", func(t *testing.T) {
		AssertEllipse(t, Ellipse[int]{Pt(1, 2), Sz(-10, 4), 0}.Canonical(), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
	t.Run("the angle is normalized", func(t *testing.T) {
		assert.Equal(t, Ellipse[int]{Pt(1, 2), Sz(10, 4), -Pi / 2}.Canonical().Angle, 3*Pi/2)
	})
	t.Run("a residue of a turn snaps to zero", func(t *testing.T) {
		assert.True(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 1e-9).Canonical().IsAligned())
	})
}

func TestEllipse_Grow(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Grow(2), Ell(Pt(1, 2), Sz(12, 6), 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Grow(-6), Ell(Pt(1, 2), Sz(4, 0), 0))
	})
}

func TestEllipse_GrowXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).GrowXY(2, 3), Ell(Pt(1, 2), Sz(12, 7), 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).GrowXY(0, -6), Ell(Pt(1, 2), Sz(10, 0), 0))
	})
}

func TestEllipse_Shrink(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Shrink(2), Ell(Pt(1, 2), Sz(8, 2), 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Shrink(6), Ell(Pt(1, 2), Sz(4, 0), 0))
	})
}

func TestEllipse_ShrinkXY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).ShrinkXY(2, 3), Ell(Pt(1, 2), Sz(8, 1), 0))
	})
	t.Run("clamped at zero", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).ShrinkXY(0, 6), Ell(Pt(1, 2), Sz(10, 0), 0))
	})
}

func TestEllipse_Lerp(t *testing.T) {
	a, b := Ell(Pt(0.0, 0.0), Sz(10.0, 4.0), 0), Ell(Pt(4.0, 8.0), Sz(2.0, 8.0), Pi/2)

	t.Run("half way", func(t *testing.T) {
		AssertEllipse(t, a.Lerp(b, 0.5), Ell(Pt(2.0, 4.0), Sz(6.0, 6.0), Pi/4))
	})
	t.Run("turns along the shorter arc", func(t *testing.T) {
		assert.Equal(t, Ell(Pt(0.0, 0.0), SzU(1.0), 2*Pi-0.2).Lerp(Ell(Pt(0.0, 0.0), SzU(1.0), 0.2), 0.5).Angle, 0.0)
	})
	t.Run("extrapolates past a zero semi-axis", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(0.0, 0.0), Sz(2.0, 2.0), 0).Lerp(Ell(Pt(0.0, 0.0), Sz(1.0, 1.0), 0), 3), Ell(Pt(0.0, 0.0), Sz(1.0, 1.0), 0))
	})
}

func TestEllipse_Transform(t *testing.T) {
	e := Ell(Pt(2.0, 3.0), Sz(4.0, 2.0), 0)

	t.Run("a move", func(t *testing.T) {
		AssertEllipse(t, e.Transform(TranslationMatrix(1.0, -1.0)), Ell(Pt(3.0, 2.0), Sz(4.0, 2.0), 0))
	})
	t.Run("a turn", func(t *testing.T) {
		AssertEllipse(t, e.Transform(RotationMatrix[float64](Pi/2)), Ell(Pt(-3.0, 2.0), Sz(4.0, 2.0), Pi/2))
	})
	t.Run("a scale of the axes while aligned", func(t *testing.T) {
		AssertEllipse(t, e.Transform(ScaleMatrix(2.0, 3.0)), Ell(Pt(4.0, 9.0), Sz(8.0, 6.0), 0))
	})
	t.Run("a reflection mirrors the angle", func(t *testing.T) {
		AssertEllipse(t, e.Rotate(Pi/6).Transform(ReflectionMatrix[float64](AxisHorizontal)), Ell(Pt(2.0, -3.0), Sz(4.0, 2.0), -Pi/6))
	})
	t.Run("a shear gives the nearest ellipse", func(t *testing.T) {
		// the true image is a turned ellipse; Scaling and Angle name the nearest one
		AssertEllipse(t, e.Transform(ShearMatrix(1.0, 0.0)), Ell(Pt(5.0, 3.0), Sz(4.0, 2.0), 0))
	})
}

func TestEllipse_Rotate(t *testing.T) {
	t.Run("turns about the center", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0).Rotate(Pi/2), Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), Pi/2))
	})
	t.Run("normalized to a single turn", func(t *testing.T) {
		assert.Equal(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), Pi).Rotate(3*Pi/2).Angle, Pi/2)
	})
}

func TestEllipse_AlignTo(t *testing.T) {
	e := Ell(Pt(10.0, 10.0), Sz(5.0, 3.0), 0)

	t.Run("an anchor lands on the point", func(t *testing.T) {
		AssertEllipse(t, e.AlignTo(Left, Pt(0.0, 0.0)), Ell(Pt(5.0, 0.0), Sz(5.0, 3.0), 0))
	})
	t.Run("none aligns the center", func(t *testing.T) {
		AssertEllipse(t, e.AlignTo(DirectionNone, Pt(0.0, 0.0)), Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0))
	})
}

func TestEllipse_Contains(t *testing.T) {
	e := Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0)

	t.Run("inside", func(t *testing.T) {
		assert.True(t, e.Contains(Pt(4.0, 0.0)))
		assert.True(t, e.Contains(Pt(0.0, 2.9)))
	})
	t.Run("outside the shorter axis", func(t *testing.T) {
		assert.False(t, e.Contains(Pt(0.0, 4.0)))
		assert.False(t, e.Contains(Pt(4.0, 2.0)))
	})
	t.Run("the boundary is included", func(t *testing.T) {
		assert.True(t, e.Contains(Pt(5.0, 0.0)))
		assert.True(t, e.Contains(Pt(0.0, 3.0)))
		assert.True(t, e.Contains(Pt(5.0+Delta/2, 0.0)))
	})
	t.Run("the turn is applied", func(t *testing.T) {
		assert.True(t, e.Rotate(Pi/2).Contains(Pt(0.0, 4.0)))
		assert.False(t, e.Rotate(Pi/2).Contains(Pt(4.0, 0.0)))
	})
	t.Run("a degenerate ellipse is its segment", func(t *testing.T) {
		d := Ell(Pt(0.0, 0.0), Sz(5.0, 0.0), 0)

		assert.True(t, d.Contains(Pt(3.0, 0.0)))
		assert.False(t, d.Contains(Pt(3.0, 1.0)))
		assert.False(t, d.Contains(Pt(6.0, 0.0)))
	})
}

func BenchmarkEllipse_Contains(b *testing.B) {
	e := Ell(Pt(0.0, 0.0), Sz(50.0, 30.0), Pi/6)
	inside, outside := Pt(10.0, 20.0), Pt(99.0, 99.0)

	b.Run("inside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = e.Contains(inside)
		}
	})
	b.Run("outside", func(b *testing.B) {
		for b.Loop() {
			sinkBool = e.Contains(outside)
		}
	})
}

func TestEllipse_DistanceTo(t *testing.T) {
	e := Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0)

	t.Run("zero inside", func(t *testing.T) {
		AssertNumber(t, e.DistanceTo(Pt(1.0, 1.0)), 0.0)
	})
	t.Run("along an axis", func(t *testing.T) {
		AssertNumber(t, e.DistanceTo(Pt(10.0, 0.0)), 5.0)
		AssertNumber(t, e.DistanceTo(Pt(0.0, 10.0)), 7.0)
	})
	t.Run("off both axes", func(t *testing.T) {
		// the nearest point of the boundary, not the point at the same angle
		AssertNumber(t, e.DistanceTo(Pt(8.0, 6.0)), Pt(8.0, 6.0).DistanceTo(e.worldPoint(e.nearestOffset(Vec(8.0, 6.0)))))
	})
	t.Run("the turn is applied", func(t *testing.T) {
		AssertNumber(t, e.Rotate(Pi/2).DistanceTo(Pt(0.0, 10.0)), 5.0)
	})
	t.Run("a degenerate ellipse measures to its segment", func(t *testing.T) {
		AssertNumber(t, Ell(Pt(0.0, 0.0), Sz(5.0, 0.0), 0).DistanceTo(Pt(9.0, 3.0)), 5.0)
	})
}

func TestEllipse_DistanceSquaredTo(t *testing.T) {
	e := Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0)

	t.Run("the square of the distance", func(t *testing.T) {
		AssertNumber(t, e.DistanceSquaredTo(Pt(0.0, 10.0)), 49.0)
	})
	t.Run("a point within the tolerance of the boundary is on it", func(t *testing.T) {
		assert.Equal(t, e.DistanceSquaredTo(Pt(0.0, 3.0+Delta/2)), 0.0)
	})
	t.Run("int has no tolerance", func(t *testing.T) {
		assert.True(t, Ell(Pt(0, 0), Sz(5, 3), 0).DistanceSquaredTo(Pt(0, 4)) > 0)
	})
}

func TestEllipse_nearestOffset(t *testing.T) {
	e := Ell(Pt(0.0, 0.0), Sz(5.0, 3.0), 0)

	t.Run("on the major axis beyond the evolute", func(t *testing.T) {
		AssertVector(t, e.nearestOffset(Vec(10.0, 0.0)), Vec(5.0, 0.0))
	})
	t.Run("on the major axis within the evolute", func(t *testing.T) {
		// the foot leaves the axis: the center of curvature at (5,0) is at (16/5, 0)
		ratio := 5 * 1.0 / (5*5 - 3*3)

		AssertVector(t, e.nearestOffset(Vec(1.0, 0.0)), Vec(5*ratio, 3*math.Sqrt(1-ratio*ratio)))
	})
	t.Run("on the minor axis", func(t *testing.T) {
		AssertVector(t, e.nearestOffset(Vec(0.0, 10.0)), Vec(0.0, 3.0))
	})
	t.Run("inside, off both axes", func(t *testing.T) {
		// the bisection brackets the root from the other side, and still lands on the boundary
		assert.True(t, EqualDelta(e.form(e.nearestOffset(Vec(1.0, 1.0))), 1, Delta))
	})
	t.Run("the signs of the offset are kept", func(t *testing.T) {
		AssertVector(t, e.nearestOffset(Vec(-10.0, 0.0)), Vec(-5.0, 0.0))
		AssertVector(t, e.nearestOffset(Vec(0.0, -10.0)), Vec(0.0, -3.0))
	})
	t.Run("a taller ellipse swaps the axes", func(t *testing.T) {
		AssertVector(t, Ell(Pt(0.0, 0.0), Sz(3.0, 5.0), 0).nearestOffset(Vec(0.0, 10.0)), Vec(0.0, 5.0))
		AssertVector(t, Ell(Pt(0.0, 0.0), Sz(3.0, 5.0), 0).nearestOffset(Vec(10.0, 0.0)), Vec(3.0, 0.0))
	})
	t.Run("a degenerate ellipse clamps to its segment", func(t *testing.T) {
		AssertVector(t, Ell(Pt(0.0, 0.0), Sz(5.0, 0.0), 0).nearestOffset(Vec(9.0, 3.0)), Vec(5.0, 0.0))
		AssertVector(t, Ell(Pt(0.0, 0.0), Sz(0.0, 5.0), 0).nearestOffset(Vec(9.0, 3.0)), Vec(0.0, 3.0))
	})
	t.Run("it is the nearest point of a fine sampling", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, point := range pointFixtures {
				local := e.localOffset(point)
				nearest := local.Subtract(e.nearestOffset(local)).LengthSquared()

				for i := range 720 {
					sample := VectorFromAngleSize(float64(i)*Pi/360, e.Size.Float())

					assert.True(t, nearest <= local.Subtract(sample).LengthSquared()+Delta, fmt.Sprintf("%s → %s: ", e, point))
				}
			}
		}
	})
}

func TestEllipse_Equal(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assert.True(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5).Equal(Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5)))
	})
	t.Run("a full turn does not matter", func(t *testing.T) {
		assert.True(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5).Equal(Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0.5+2*Pi)))
	})
	t.Run("two circles of different angles are different values", func(t *testing.T) {
		assert.False(t, Ell(Pt(1.0, 2.0), SzU(4.0), 0).Equal(Ell(Pt(1.0, 2.0), SzU(4.0), 1.0)))
	})
	t.Run("a different size", func(t *testing.T) {
		assert.False(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0).Equal(Ell(Pt(1.0, 2.0), Sz(10.0, 5.0), 0)))
	})
}

func TestEllipse_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.True(t, Ellipse[int]{}.IsZero())
	})
	t.Run("a size makes it non-zero", func(t *testing.T) {
		assert.False(t, Ell(Pt(0, 0), Sz(1, 0), 0).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ell(Pt(0.0, 0.0000001), Sz(0.0, 0.0), 0).IsZero())
	})
}

func TestEllipse_IsAligned(t *testing.T) {
	t.Run("aligned", func(t *testing.T) {
		assert.True(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0).IsAligned())
	})
	t.Run("turned", func(t *testing.T) {
		assert.False(t, Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 1e-9).IsAligned())
	})
}

func TestEllipse_IsCircle(t *testing.T) {
	t.Run("equal semi-axes", func(t *testing.T) {
		assert.True(t, Ell(Pt(1.0, 2.0), SzU(4.0), 1.0).IsCircle())
	})
	t.Run("different semi-axes", func(t *testing.T) {
		assert.False(t, Ell(Pt(1.0, 2.0), Sz(4.0, 3.0), 0).IsCircle())
	})
}

func TestEllipse_Circle(t *testing.T) {
	t.Run("a circle converts exactly", func(t *testing.T) {
		AssertCircle(t, Ell(Pt(1, 2), SzU(4), 1.0).Circle(), Circ(Pt(1, 2), 4))
	})
	t.Run("an ellipse gives the circle around it", func(t *testing.T) {
		AssertCircle(t, Ell(Pt(1, 2), Sz(4, 3), 0).Circle(), Circ(Pt(1, 2), 4))
		AssertCircle(t, Ell(Pt(1, 2), Sz(3, 4), 1.0).Circle(), Circ(Pt(1, 2), 4))
	})
}

func TestEllipse_RegularPolygon(t *testing.T) {
	e := Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), Pi/6)

	t.Run("the same center, semi-axes and angle", func(t *testing.T) {
		AssertRegularPolygon(t, e.RegularPolygon(6), RegPol(Pt(1.0, 2.0), Sz(10.0, 4.0), 6, Pi/6))
	})
	t.Run("every vertex lies on the boundary", func(t *testing.T) {
		for vertex := range e.RegularPolygon(7).Vertices() {
			AssertNumber(t, e.DistanceTo(vertex), 0.0, fmt.Sprintf("%s: ", vertex))
		}
	})
	t.Run("the polygon is inside the ellipse", func(t *testing.T) {
		for edge := range e.RegularPolygon(7).Edges() {
			assert.True(t, e.Contains(edge.Midpoint()), fmt.Sprintf("%s: ", edge))
		}
	})
	t.Run("fewer than one vertex is empty", func(t *testing.T) {
		assert.True(t, e.RegularPolygon(0).IsEmpty())
	})
	t.Run("it round-trips through Ellipse", func(t *testing.T) {
		AssertEllipse(t, e.RegularPolygon(6).Ellipse(), e)
	})
}

func TestEllipse_Cast(t *testing.T) {
	e := Ell(Pt(1.5, -2.5), Sz(3.5, 1.5), 1.0)

	t.Run("matches Int and Float", func(t *testing.T) {
		AssertEllipse(t, e.Cast[int](), e.Int())
		AssertEllipse(t, e.Cast[float64](), e.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		AssertEllipse(t, e.Cast[int8](), Ell(Pt[int8](2, -3), Sz[int8](4, 2), 1.0))
	})
}

func TestEllipse_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Int(), Ell(Pt(1, 2), Sz(10, 4), 0))
	})
	t.Run("float rounds and keeps the angle", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 1.0).Int(), Ell(Pt(1, 0), Sz(1, 1), 1.0))
	})
}

func TestEllipse_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(1, 2), Sz(10, 4), 0).Float(), Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertEllipse(t, Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0).Float(), Ell(Pt(0.6, -0.25), Sz(1.2, 0.5), 0))
	})
}

func TestEllipse_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Ell(Pt(10, 16), Sz(5, 2), 0).String(), "Ell((10,16);5x2)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Ell(Pt(100, -34.0000115), Sz(0.2, 0.1), 0).String(), "Ell((100.00,-34.00);0.20x0.10)")
	})
	t.Run("a turned ellipse carries its angle", func(t *testing.T) {
		assert.Equal(t, Ell(Pt(10, 16), Sz(5, 2), 1.5).String(), "Ell((10,16);5x2;1.50)")
	})
}

func TestEllipse_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Ell(Pt(10, 16), Sz(12, 4), 0), `{"x":10,"y":16,"w":12,"h":4}`)

		var e Ellipse[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16,"w":12,"h":4}`), &e))
		AssertEllipse(t, e, Ell(Pt(10, 16), Sz(12, 4), 0))
	})
	t.Run("a turned ellipse carries its angle", func(t *testing.T) {
		assert.JSON(t, Ell(Pt(10, 16), Sz(12, 4), 1.5), `{"x":10,"y":16,"w":12,"h":4,"a":1.5}`)
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, ellipse := range ellipseFixtures {
			data, err := json.Marshal(ellipse)
			assert.NoError(t, err)

			var decoded Ellipse[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, ellipse)
		}
	})
}

func TestEllipse_Properties(t *testing.T) {
	t.Run("area follows the semi-axes", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			AssertNumber(t, e.Area(), Pi*e.Size.Width*e.Size.Height, fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("the perimeter is between the axes and the circle around them", func(t *testing.T) {
		// the true bounds, up to the part in a thousand Ramanujan's approximation loses
		// on the degenerate ellipse, where it is loosest
		for _, e := range ellipseFixtures {
			major := float64(e.SemiMajor())

			assert.True(t, e.Perimeter() >= 4*major*(1-1e-3), fmt.Sprintf("%s: ", e))
			assert.True(t, e.Perimeter() <= 2*Pi*major+Delta, fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("every anchor is on the boundary in its direction", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, direction := range Directions() {
				anchor := e.Anchor(direction)

				AssertNumber(t, e.DistanceTo(anchor), 0.0, fmt.Sprintf("%s → %s: ", e, direction))
				assert.True(t, e.Bounds().Contains(anchor), fmt.Sprintf("%s → %s: ", e, direction))
				if e.SemiMinor() > 0 {
					AssertAngle(t, e.Center.AngleTo(anchor), e.Angle+direction.Angle(), fmt.Sprintf("%s → %s: ", e, direction))
				}
			}
		}
	})
	t.Run("the boundary is contained and its bounds hold it", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for i := range 16 {
				point := e.worldPoint(VectorFromAngleSize(float64(i)*Pi/8, e.Size.Float()))

				assert.True(t, e.Contains(point), fmt.Sprintf("%s → %s: ", e, point))
				assert.True(t, e.Bounds().Contains(point), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
	t.Run("the focal distances of a boundary point sum to the major axis", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			a, b := e.Foci()

			for i := range 16 {
				point := e.worldPoint(VectorFromAngleSize(float64(i)*Pi/8, e.Size.Float()))

				AssertNumber(t, point.DistanceTo(a)+point.DistanceTo(b), 2*float64(e.SemiMajor()), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
	t.Run("distance is zero exactly where contains holds", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, point := range pointFixtures {
				assert.Equal(t, e.DistanceTo(point) == 0, e.Contains(point), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
	t.Run("the turn does not change the distance to a point turned with it", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, point := range pointFixtures {
				turned := point.RotateAround(e.Center, Pi/3)

				AssertNumber(t, e.Rotate(Pi/3).DistanceTo(turned), e.DistanceTo(point), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
	t.Run("lerp ends on the two ellipses", func(t *testing.T) {
		for _, a := range ellipseFixtures {
			for _, b := range ellipseFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, e.Scale(factor).Unscale(factor).Equal(e), fmt.Sprintf("%s ×%v: ", e, factor))
			}
		}
	})
	t.Run("a full turn is the same ellipse", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			assert.True(t, e.Rotate(2*Pi).Equal(e), fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("the identity transform leaves the ellipse", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			AssertEllipse(t, e.Transform(IdentityMatrix[float64]()), e.Canonical(), fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("grow and shrink are inverse above zero", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			if e.SemiMinor() < 1 {
				continue // shrink clamps to zero, so the growth is not recoverable
			}

			assert.True(t, e.Grow(1).Shrink(1).Equal(e), fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("the circle around holds the boundary", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			around := e.Circle()

			for i := range 16 {
				point := e.worldPoint(VectorFromAngleSize(float64(i)*Pi/8, e.Size.Float()))

				assert.True(t, around.Contains(point), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
	t.Run("the polygon lies within the ellipse, on its boundary", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, n := range []int{3, 5, 8, 32} {
				polygon := e.RegularPolygon(n)

				for vertex := range polygon.Vertices() {
					AssertNumber(t, e.DistanceTo(vertex), 0.0, fmt.Sprintf("%s ×%d → %s: ", e, n, vertex))
				}

				for edge := range polygon.Edges() {
					assert.True(t, e.Contains(edge.Midpoint()), fmt.Sprintf("%s ×%d → %s: ", e, n, edge))
				}
			}
		}
	})
	t.Run("scaling by the apothem ratio puts the polygon about the ellipse", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			for _, n := range []int{3, 5, 8, 32} {
				around := e.Scale(1 / math.Cos(Pi/float64(n))).RegularPolygon(n)

				for edge := range around.Edges() {
					AssertNumber(t, e.DistanceTo(edge.Midpoint()), 0.0, fmt.Sprintf("%s ×%d → %s: ", e, n, edge))
				}
			}
		}
	})
	t.Run("the polygon round-trips through the ellipse", func(t *testing.T) {
		for _, e := range ellipseFixtures {
			AssertEllipse(t, e.RegularPolygon(6).Ellipse(), e, fmt.Sprintf("%s: ", e))
		}
	})
	t.Run("a circle answers as the circle of the same radius", func(t *testing.T) {
		for _, c := range circleFixtures {
			e := c.Ellipse()

			for _, point := range pointFixtures {
				assert.Equal(t, e.Contains(point), c.Contains(point), fmt.Sprintf("%s → %s: ", e, point))
				AssertNumber(t, e.DistanceTo(point), c.DistanceTo(point), fmt.Sprintf("%s → %s: ", e, point))
			}
		}
	})
}

func TestEllipse_Immutable(t *testing.T) {
	e := Ell(Pt(1, 2), Sz(10, 4), 0)

	e.Translate(Vec(3, -2))
	e.MoveTo(Pt(4, 3))
	e.Scale(2)
	e.Resize(Sz(1, 1))
	e.Grow(1)
	e.Shrink(2)
	e.Rotate(1)

	AssertEllipse(t, e, Ell(Pt(1, 2), Sz(10, 4), 0))
}

// ellipseFixtures span the degenerate, circular, turned and off-origin cases.

var ellipseFixtures = []Ellipse[float64]{
	Ell(Pt(0.0, 0.0), Sz(0.0, 0.0), 0),
	Ell(Pt(0.0, 0.0), SzU(1.0), 0),
	Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), 0),
	Ell(Pt(1.0, 2.0), Sz(10.0, 4.0), Pi/6),
	Ell(Pt(-3.5, 0.25), Sz(2.0, 7.0), Pi/2),
	Ell(Pt(12.5, -0.1), Sz(3.0, 0.0), 1.0),
}

func ExampleEll() {
	fmt.Println(Ell(Pt(10, 16), Sz(5, 2), 0))
	// Output: Ell((10,16);5x2)
}

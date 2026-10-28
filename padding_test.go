package geom_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestPadding_Constructor(t *testing.T) {
	t.Run("from edges", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(5, 4, 2, 3), Padding[int]{Left: 5, Right: 4, Top: 2, Bottom: 3})
		geomtest.AssertPadding(t, Pad(2.4, 3.0, 0.1, 0.6), Padding[float64]{Left: 2.4, Right: 3.0, Top: 0.1, Bottom: 0.6})
	})
	t.Run("uniform", func(t *testing.T) {
		geomtest.AssertPadding(t, PadU(20), Pad(20, 20, 20, 20))
		geomtest.AssertPadding(t, PadU(20.0), Pad(20.0, 20.0, 20.0, 20.0))
	})
	t.Run("per-axis", func(t *testing.T) {
		geomtest.AssertPadding(t, PadXY(20, 10), Pad(20, 20, 10, 10))
		geomtest.AssertPadding(t, PadXY(20.0, 10.0), Pad(20.0, 20.0, 10.0, 10.0))
	})
}

func TestPadding_XY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		width, height := Pad(5, 4, 2, 3).XY()
		geomtest.AssertNumber(t, width, 9)
		geomtest.AssertNumber(t, height, 5)
	})
	t.Run("float", func(t *testing.T) {
		width, height := Pad(2.4, 3.0, 0.1, 0.6).XY()
		geomtest.AssertNumber(t, width, 5.4)
		geomtest.AssertNumber(t, height, 0.7)
	})
}

func TestPadding_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Pad(5, 4, 2, 3).Width(), 9)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Pad(2.4, 3.0, 0.1, 0.6).Width(), 5.4)
	})
}

func TestPadding_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Pad(5, 4, 2, 3).Height(), 5)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Pad(2.4, 3.0, 0.1, 0.6).Height(), 0.7)
	})
}

func TestPadding_Add(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(4, 2, 1, 3).Add(Pad(40, 20, 10, 30)), Pad(44, 22, 11, 33))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(0.4, 0.2, 0.1, 0.3).Add(PadU(0.5)), Pad(0.9, 0.7, 0.6, 0.8))
	})
}

func TestPadding_Subtract(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(44, 22, 11, 33).Subtract(Pad(4, 2, 1, 3)), Pad(40, 20, 10, 30))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(0.9, 0.7, 0.6, 0.8).Subtract(PadU(0.5)), Pad(0.4, 0.2, 0.1, 0.3))
	})
	t.Run("undoes add", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(4, 2, 1, 3).Add(Pad(8, 6, 5, 7)).Subtract(Pad(8, 6, 5, 7)), Pad(4, 2, 1, 3))
	})
}

func TestPadding_Negate(t *testing.T) {
	t.Run("negates every edge", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(-4, -2, 1, 3).Negate(), Pad(4, 2, -1, -3))
	})
	t.Run("undoes an inset", func(t *testing.T) {
		r := RectangleFromMinMax(Pt(0.0, 0.0), Pt(10.0, 10.0))
		padding := Pad(4.0, 2.0, 1.0, 3.0)

		geomtest.AssertRectangle(t, r.Inset(padding).Inset(padding.Negate()), r)
		geomtest.AssertRectangle(t, r.Outset(padding).Inset(padding.Negate()), r.Outset(padding).Outset(padding))
	})
}

func TestPadding_Scale(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(4, 2, 1, 3).Scale(1.5), Pad(6, 3, 2, 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(0.4, 0.2, 0.1, 0.3).Scale(2), Pad(0.8, 0.4, 0.2, 0.6))
	})
}

func TestPadding_Unscale(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(12, 6, 3, 9).Unscale(2), Pad(6, 3, 2, 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(0.8, 0.4, 0.2, 0.6).Unscale(2), Pad(0.4, 0.2, 0.1, 0.3))
	})
	t.Run("undoes scale", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(4.5, 2.5, 1.5, 3.5).Scale(4).Unscale(4), Pad(4.5, 2.5, 1.5, 3.5))
	})
	t.Run("zero factor panics", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			Pad(4, 2, 1, 3).Unscale(0)
		}, ErrDivisionByZero)
	})
}

func TestPadding_Lerp(t *testing.T) {
	t.Run("interpolates edge by edge", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(0.0, 0.0, 0.0, 0.0).Lerp(Pad(8.0, 4.0, 2.0, 6.0), 0.5), Pad(4.0, 2.0, 1.0, 3.0))
	})
	t.Run("the ends are the paddings themselves", func(t *testing.T) {
		a, b := Pad(4, 2, 1, 3), Pad(40, 20, 10, 30)

		geomtest.AssertPadding(t, a.Lerp(b, 0), a)
		geomtest.AssertPadding(t, a.Lerp(b, 1), b)
	})
}

func TestPadding_Equal(t *testing.T) {
	t.Run("same padding", func(t *testing.T) {
		assert.True(t, Pad(4, 2, 1, 3).Equal(Pad(4, 2, 1, 3)))
		assert.True(t, Pad(0.4, 0.2, 0.1, 0.3).Equal(Pad(0.4000001, 0.2, 0.1, 0.3)))
	})
	t.Run("different padding", func(t *testing.T) {
		assert.False(t, Pad(4, 2, 1, 3).Equal(Pad(1, 3, 4, 2)))
		assert.False(t, Pad(0.4, 0.2, 0.1, 0.3).Equal(Pad(0.5, 0.2, 0.1, 0.3)))
	})
}

func TestPadding_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.True(t, PadU(0).IsZero())
		assert.True(t, Padding[float64]{}.IsZero())
	})
	t.Run("non-zero", func(t *testing.T) {
		assert.False(t, Pad(1, 0, 0, 0).IsZero())
		assert.False(t, PadU(0.5).IsZero())
	})
}

func TestPadding_Size(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertSize(t, Pad(5, 4, 2, 3).Size(), Sz(9, 5))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertSize(t, Pad(2.4, 3.0, 0.1, 0.6).Size(), Sz(5.4, 0.7))
	})
}

func TestPadding_Cast(t *testing.T) {
	p := Pad(-2.5, 2.5, 1.5, -1.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		geomtest.AssertPadding(t, p.Cast[int](), p.Int())
		geomtest.AssertPadding(t, p.Cast[float64](), p.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		geomtest.AssertPadding(t, p.Cast[int8](), Pad[int8](-3, 3, 2, -2))
	})
}

func TestPadding_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(5, 4, 2, 3).Int(), Pad(5, 4, 2, 3))
	})
	t.Run("float rounds", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(2.4, 3.0, 0.1, 0.6).Int(), Pad(2, 3, 0, 1))
	})
}

func TestPadding_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(5, 4, 2, 3).Float(), Pad(5.0, 4.0, 2.0, 3.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		geomtest.AssertPadding(t, Pad(2.4, 3.0, 0.1, 0.6).Float(), Pad(2.4, 3.0, 0.1, 0.6))
	})
}

func TestPadding_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Pad(5, 4, 2, 3).String(), "Pad(5;4;2;3)")
		assert.Equal(t, Pad(-3, 0, -1, 2).String(), "Pad(-3;0;-1;2)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Pad(2.4, 3.0, 0.1, 0.6).String(), "Pad(2.40;3.00;0.10;0.60)")
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.Equal(t, Pad(0, 0, -0, 0).String(), "Pad(0;0;0;0)")
		assert.Equal(t, Pad(0.0, 0.0, negativeZero, 0.0).String(), "Pad(0.00;0.00;0.00;0.00)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assertStringAllocation(t, Pad(2.4, 3.0, 0.1, 0.6).String)
	})
}

func TestPadding_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Pad(5, 4, 2, 3), `{"l":5,"r":4,"t":2,"b":3}`)

		var p Padding[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"l":5,"r":4,"t":2,"b":3}`), &p))
		geomtest.AssertPadding(t, p, Pad(5, 4, 2, 3))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Pad(2.4, 3.0, 0.1, 0.6), `{"l":2.40,"r":3,"t":0.10,"b":0.60}`)

		var p Padding[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"l":2.40,"r":3,"t":0.10,"b":0.60}`), &p))
		geomtest.AssertPadding(t, p, Pad(2.4, 3.0, 0.1, 0.6))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, padding := range paddingFixtures {
			data, err := json.Marshal(padding)
			assert.NoError(t, err)

			var decoded Padding[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, padding)
		}
	})
}

func TestPadding_Properties(t *testing.T) {
	t.Run("width and height sum the opposite edges", func(t *testing.T) {
		for _, padding := range paddingFixtures {
			geomtest.AssertNumber(t, padding.Width(), padding.Left+padding.Right, fmt.Sprintf("%s: ", padding))
			geomtest.AssertNumber(t, padding.Height(), padding.Top+padding.Bottom, fmt.Sprintf("%s: ", padding))
		}
	})
	t.Run("xy and size agree with width and height", func(t *testing.T) {
		for _, padding := range paddingFixtures {
			width, height := padding.XY()

			geomtest.AssertNumber(t, width, padding.Width(), fmt.Sprintf("%s: ", padding))
			geomtest.AssertNumber(t, height, padding.Height(), fmt.Sprintf("%s: ", padding))
			assert.True(t, padding.Size().Equal(Sz(padding.Width(), padding.Height())), fmt.Sprintf("%s: ", padding))
		}
	})
	t.Run("lerp ends on the two paddings", func(t *testing.T) {
		for _, a := range paddingFixtures {
			for _, b := range paddingFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s -> %s: ", a, b))
				assert.True(t, a.Lerp(b, 0.5).Equal(b.Lerp(a, 0.5)), fmt.Sprintf("%s -> %s: ", a, b))
			}
		}
	})
	t.Run("float is the inverse of int on whole values", func(t *testing.T) {
		for _, padding := range []Padding[int]{Pad(5, 4, 2, 3), PadU(0), PadXY(-1, 7)} {
			geomtest.AssertPadding(t, padding.Float().Int(), padding)
		}
	})
}

// paddingFixtures span uniform, per-axis, zero, and negative edges.
var paddingFixtures = []Padding[float64]{
	Pad(0.0, 0.0, 0.0, 0.0),
	Pad(2.4, 3.0, 0.1, 0.6),
	Pad(5.0, 4.0, 2.0, 3.0),
	PadU(20.0),
	PadXY(20.0, 10.0),
	Pad(12.75, 0.25, -1.5, -0.5),
}

func ExamplePad() {
	fmt.Println(Pad(5, 4, 2, 3))
	fmt.Println(PadU(1.5))
	// Output:
	// Pad(5;4;2;3)
	// Pad(1.50;1.50;1.50;1.50)
}

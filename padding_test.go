package geom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestPadding_Constructor(t *testing.T) {
	t.Run("from edges", func(t *testing.T) {
		AssertPadding(t, Pad(2, 4, 3, 5), Padding[int]{Top: 2, Right: 4, Bottom: 3, Left: 5})
		AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4), Padding[float64]{Top: 0.1, Right: 3.0, Bottom: 0.6, Left: 2.4})
	})
	t.Run("uniform", func(t *testing.T) {
		AssertPadding(t, PadU(20), Pad(20, 20, 20, 20))
		AssertPadding(t, PadU(20.0), Pad(20.0, 20.0, 20.0, 20.0))
	})
	t.Run("per-axis", func(t *testing.T) {
		AssertPadding(t, PadXY(10, 20), Pad(10, 20, 10, 20))
		AssertPadding(t, PadXY(10.0, 20.0), Pad(10.0, 20.0, 10.0, 20.0))
	})
}

func TestPadding_Width(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pad(2, 4, 3, 5).Width(), 9)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Pad(0.1, 3.0, 0.6, 2.4).Width(), 5.4)
	})
}

func TestPadding_Height(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pad(2, 4, 3, 5).Height(), 5)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Pad(0.1, 3.0, 0.6, 2.4).Height(), 0.7)
	})
}

func TestPadding_XY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		width, height := Pad(2, 4, 3, 5).XY()
		AssertNumber(t, width, 9)
		AssertNumber(t, height, 5)
	})
	t.Run("float", func(t *testing.T) {
		width, height := Pad(0.1, 3.0, 0.6, 2.4).XY()
		AssertNumber(t, width, 5.4)
		AssertNumber(t, height, 0.7)
	})
}

func TestPadding_Size(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertSize(t, Pad(2, 4, 3, 5).Size(), Sz(9, 5))
	})
	t.Run("float", func(t *testing.T) {
		AssertSize(t, Pad(0.1, 3.0, 0.6, 2.4).Size(), Sz(5.4, 0.7))
	})
}

func TestPadding_Equal(t *testing.T) {
	t.Run("same padding", func(t *testing.T) {
		assert.True(t, Pad(1, 2, 3, 4).Equal(Pad(1, 2, 3, 4)))
		assert.True(t, Pad(0.1, 0.2, 0.3, 0.4).Equal(Pad(0.1, 0.2, 0.3, 0.4000001)))
	})
	t.Run("different padding", func(t *testing.T) {
		assert.False(t, Pad(1, 2, 3, 4).Equal(Pad(4, 3, 2, 1)))
		assert.False(t, Pad(0.1, 0.2, 0.3, 0.4).Equal(Pad(0.1, 0.2, 0.3, 0.5)))
	})
}

func TestPadding_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.True(t, PadU(0).IsZero())
		assert.True(t, Padding[float64]{}.IsZero())
	})
	t.Run("non-zero", func(t *testing.T) {
		assert.False(t, Pad(0, 0, 0, 1).IsZero())
		assert.False(t, PadU(0.5).IsZero())
	})
}

func TestPadding_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPadding(t, Pad(2, 4, 3, 5).Int(), Pad(2, 4, 3, 5))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4).Int(), Pad(0, 3, 1, 2))
	})
}

func TestPadding_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertPadding(t, Pad(2, 4, 3, 5).Float(), Pad(2.0, 4.0, 3.0, 5.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertPadding(t, Pad(0.1, 3.0, 0.6, 2.4).Float(), Pad(0.1, 3.0, 0.6, 2.4))
	})
}

func TestPadding_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Pad(2, 4, 3, 5).String(), "Pad(2;4;3;5)")
		assert.Equal(t, Pad(-1, 0, 2, -3).String(), "Pad(-1;0;2;-3)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Pad(0.1, 3.0, 0.6, 2.4).String(), "Pad(0.10;3.00;0.60;2.40)")
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.Equal(t, Pad(-0, 0, 0, 0).String(), "Pad(0;0;0;0)")
		assert.Equal(t, Pad(negativeZero, 0.0, 0.0, 0.0).String(), "Pad(0.00;0.00;0.00;0.00)")
	})
}

func TestPadding_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Pad(2, 4, 3, 5), `{"t":2,"r":4,"b":3,"l":5}`)

		var p Padding[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"t":2,"r":4,"b":3,"l":5}`), &p))
		AssertPadding(t, p, Pad(2, 4, 3, 5))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Pad(0.1, 3.0, 0.6, 2.4), `{"t":0.10,"r":3,"b":0.60,"l":2.40}`)

		var p Padding[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"t":0.10,"r":3,"b":0.60,"l":2.40}`), &p))
		AssertPadding(t, p, Pad(0.1, 3.0, 0.6, 2.4))
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
			AssertNumber(t, padding.Width(), padding.Left+padding.Right, fmt.Sprintf("%s: ", padding))
			AssertNumber(t, padding.Height(), padding.Top+padding.Bottom, fmt.Sprintf("%s: ", padding))
		}
	})
	t.Run("xy and size agree with width and height", func(t *testing.T) {
		for _, padding := range paddingFixtures {
			width, height := padding.XY()

			AssertNumber(t, width, padding.Width(), fmt.Sprintf("%s: ", padding))
			AssertNumber(t, height, padding.Height(), fmt.Sprintf("%s: ", padding))
			assert.True(t, padding.Size().Equal(Sz(padding.Width(), padding.Height())), fmt.Sprintf("%s: ", padding))
		}
	})
	t.Run("float is the inverse of int on whole values", func(t *testing.T) {
		for _, padding := range []Padding[int]{Pad(2, 4, 3, 5), PadU(0), PadXY(-1, 7)} {
			AssertPadding(t, padding.Float().Int(), padding)
		}
	})
}

// paddingFixtures span uniform, per-axis, zero, and negative edges.
var paddingFixtures = []Padding[float64]{
	Pad(0.0, 0.0, 0.0, 0.0),
	Pad(0.1, 3.0, 0.6, 2.4),
	Pad(2.0, 4.0, 3.0, 5.0),
	PadU(20.0),
	PadXY(10.0, 20.0),
	Pad(-1.5, 0.25, -0.5, 12.75),
}

func ExamplePad() {
	fmt.Println(Pad(2, 4, 3, 5))
	fmt.Println(PadU(1.5))
	// Output:
	// Pad(2;4;3;5)
	// Pad(1.50;1.50;1.50;1.50)
}

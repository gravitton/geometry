package geom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestSize_Constructor(t *testing.T) {
	t.Run("from dimensions", func(t *testing.T) {
		AssertSize(t, Sz(10, 16), Size[int]{Width: 10, Height: 16})
		AssertSize(t, Sz[float64](0.16, 204), Size[float64]{Width: 0.16, Height: 204})
	})
	t.Run("uniform", func(t *testing.T) {
		AssertSize(t, SzU(2), Sz(2, 2))
		AssertSize(t, SzU(0.2), Sz(0.2, 0.2))
	})
}

func TestParseSize(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		size, err := ParseSize[int]("16x32")
		assert.NoError(t, err)
		AssertSize(t, size, Sz(16, 32))
	})
	t.Run("float", func(t *testing.T) {
		size, err := ParseSize[float64]("23.0x12.1")
		assert.NoError(t, err)
		AssertSize(t, size, Sz(23.0, 12.1))

		size, err = ParseSize[float64]("16x32")
		assert.NoError(t, err)
		AssertSize(t, size, Sz(16.0, 32.0))
	})
	t.Run("int rejects fractional values", func(t *testing.T) {
		_, err := ParseSize[int]("23.5x12.4")
		assert.Error(t, err)
	})
	t.Run("malformed input", func(t *testing.T) {
		_, err := ParseSize[int]("16")
		assert.Error(t, err, "missing separator: ")

		_, err = ParseSize[int]("axb")
		assert.Error(t, err, "invalid width: ")

		_, err = ParseSize[float64]("1.0xb")
		assert.Error(t, err, "invalid height: ")
	})
}

func TestSize_Scale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Scale(2.5), Sz(5, 8)) // int: 7.5 rounds to 8
		AssertSize(t, Sz(0.4, -0.25).Scale(2.5), Sz(1.0, -0.625))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).ScaleXY(2, 3), Sz(4, 9))
		AssertSize(t, Sz(0.4, -0.25).ScaleXY(-1.5, 2), Sz(-0.6, -0.5))
	})
}

func TestSize_Unscale(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertSize(t, Sz(10, 20).Unscale(2.0), Sz(5, 10))
		AssertSize(t, Sz(0.5, 2.5).Unscale(2.5), Sz(0.2, 1.0))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertSize(t, Sz(10, 20).UnscaleXY(2.0, 4.0), Sz(5, 5))
		AssertSize(t, Sz(0.5, 2.5).UnscaleXY(2.5, 0.5), Sz(0.2, 5.0))
	})
	t.Run("zero factor leaves the axis unchanged", func(t *testing.T) {
		AssertSize(t, Sz(10, 20).Unscale(0), Sz(10, 20))

		// the guard is per axis, so the other one still divides
		AssertSize(t, Sz(10, 20).UnscaleXY(0, 2), Sz(10, 10))
	})
}

func TestSize_Grow(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Grow(2), Sz(4, 5))
		AssertSize(t, Sz(0.4, 0.25).Grow(0.1), Sz(0.5, 0.35))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).GrowXY(2, 3), Sz(4, 6))
		AssertSize(t, Sz(0.4, 0.25).GrowXY(0.1, 0.2), Sz(0.5, 0.45))
	})
}

func TestSize_Shrink(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Shrink(1), Sz(1, 2))
		AssertSize(t, Sz(0.4, 0.25).Shrink(0.1), Sz(0.3, 0.15))
	})
	t.Run("per-axis amount", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).ShrinkXY(1, 2), Sz(1, 1))
		AssertSize(t, Sz(0.4, 0.25).ShrinkXY(0.1, 0.2), Sz(0.3, 0.05))
	})
	t.Run("clamps to zero", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Shrink(5), Sz(0, 0))
		AssertSize(t, Sz(2, 3).ShrinkXY(5, 1), Sz(0, 2))
		AssertSize(t, Sz(0.4, 0.25).Shrink(1.0), Sz(0.0, 0.0))
	})
}

func TestSize_Area(t *testing.T) {
	AssertNumber(t, Sz(5, 3).Area(), 15)
	AssertNumber(t, Sz(0.4, 0.25).Area(), 0.1)
}

func TestSize_Perimeter(t *testing.T) {
	AssertNumber(t, Sz(5, 3).Perimeter(), 16)
	AssertNumber(t, Sz(0.4, 0.25).Perimeter(), 1.3)
}

func TestSize_AspectRatio(t *testing.T) {
	t.Run("wider than tall", func(t *testing.T) {
		AssertNumber(t, Sz(5, 3).AspectRatio(), 5.0/3.0)
		AssertNumber(t, Sz(0.4, 0.25).AspectRatio(), 1.6)
	})
	t.Run("square", func(t *testing.T) {
		AssertNumber(t, SzU(7).AspectRatio(), 1.0)
	})
	t.Run("zero height", func(t *testing.T) {
		AssertNumber(t, Sz(16, 0).AspectRatio(), 0.0)
	})
}

func TestSize_AtLeast(t *testing.T) {
	t.Run("raises the smaller dimension", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).AtLeast(Sz(1, 5)), Sz(2, 5))
		AssertSize(t, Sz(2, 3).AtLeast(Sz(4, 1)), Sz(4, 3))
		AssertSize(t, Sz(0.5, 1.5).AtLeast(Sz(1.0, 1.0)), Sz(1.0, 1.5))
	})
	t.Run("same size is unchanged", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).AtLeast(Sz(2, 3)), Sz(2, 3))
	})
}

func TestSize_AtMost(t *testing.T) {
	t.Run("lowers the larger dimension", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).AtMost(Sz(1, 5)), Sz(1, 3))
		AssertSize(t, Sz(2, 3).AtMost(Sz(4, 1)), Sz(2, 1))
		AssertSize(t, Sz(0.5, 1.5).AtMost(Sz(1.0, 1.0)), Sz(0.5, 1.0))
	})
	t.Run("same size is unchanged", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).AtMost(Sz(2, 3)), Sz(2, 3))
	})
}

func TestSize_Equal(t *testing.T) {
	t.Run("same size", func(t *testing.T) {
		assert.True(t, Sz(1, 2).Equal(Sz(1, 2)))
		assert.True(t, Sz(0.4, -0.25).Equal(Sz(0.4, -0.25)))
	})
	t.Run("different size", func(t *testing.T) {
		assert.False(t, Sz(1, 2).Equal(Sz(3, -3)))
		assert.False(t, Sz(0.4, -0.25).Equal(Sz(100.1, -0.1)))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Sz(0.4, -0.25).Equal(Sz(0.4, -0.250001)))
	})
}

func TestSize_IsZero(t *testing.T) {
	t.Run("zero size", func(t *testing.T) {
		assert.True(t, Sz(0, 0).IsZero())
		assert.True(t, Sz(0.0, 0.0).IsZero())
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.True(t, Sz(-0, -0).IsZero())
		assert.True(t, Sz(negativeZero, negativeZero).IsZero())
	})
	t.Run("non-zero size", func(t *testing.T) {
		assert.False(t, Sz(1, 2).IsZero())
		assert.False(t, Sz(0.4, -0.25).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Sz(0.0, 0.000001).IsZero())
	})
}

func TestSize_XY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		width, height := Sz(10, 16).XY()
		AssertNumber(t, width, 10)
		AssertNumber(t, height, 16)
	})
	t.Run("float", func(t *testing.T) {
		width, height := Sz(0.4, -0.25).XY()
		AssertNumber(t, width, 0.4)
		AssertNumber(t, height, -0.25)
	})
}

func TestSize_Vector(t *testing.T) {
	AssertVector(t, Sz(10, 16).Vector(), Vec(10, 16))
	AssertVector(t, Sz(1.5, -2.5).Vector(), Vec(1.5, -2.5))
}

func TestSize_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Int(), Sz(2, 3))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertSize(t, Sz(1.2, 3.6).Int(), Sz(1, 4))
		AssertSize(t, Sz(-1.5, 2.5).Int(), Sz(-2, 3))
	})
}

func TestSize_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Float(), Sz(2.0, 3.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertSize(t, Sz(1.2, 3.6).Float(), Sz(1.2, 3.6))
	})
}

func TestSize_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Sz(10, 16).String(), "10x16")
		assert.Equal(t, Sz(-4, 0).String(), "-4x0")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Sz(100, -34.0000115).String(), "100.00x-34.00")
		assert.Equal(t, Sz(1.5, -0.25).String(), "1.50x-0.25")
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.Equal(t, Sz(-0, 0).String(), "0x0")
		assert.Equal(t, Sz(negativeZero, 0.0).String(), "-0.00x0.00")
	})
}

func TestSize_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Sz(10, 16), `{"w":10,"h":16}`)

		var s Size[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"w":10,"h":16}`), &s))
		AssertSize(t, s, Sz(10, 16))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Sz(100, -34.0000115), `{"w":100.0,"h":-34.0000115}`)

		var s Size[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"w":10.1,"h":34.0000115}`), &s))
		AssertSize(t, s, Sz(10.1, 34.0000115))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, size := range sizeFixtures {
			data, err := json.Marshal(size)
			assert.NoError(t, err)

			var decoded Size[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, size)
		}
	})
}

func TestSize_Properties(t *testing.T) {
	t.Run("area and perimeter follow the dimensions", func(t *testing.T) {
		for _, size := range sizeFixtures {
			AssertNumber(t, size.Area(), size.Width*size.Height, fmt.Sprintf("%s: ", size))
			AssertNumber(t, size.Perimeter(), 2*(size.Width+size.Height), fmt.Sprintf("%s: ", size))
		}
	})
	t.Run("scale and unscale are inverse", func(t *testing.T) {
		for _, size := range sizeFixtures {
			for _, factor := range []float64{0.5, 1, 2.5, -3} {
				assert.True(t, size.Scale(factor).Unscale(factor).Equal(size), fmt.Sprintf("%s ×%v: ", size, factor))
			}
		}
	})
	t.Run("grow and shrink are inverse above zero", func(t *testing.T) {
		for _, size := range sizeFixtures {
			if size.Width < 1 || size.Height < 1 {
				continue // shrink clamps to zero, so the growth is not recoverable
			}

			assert.True(t, size.Grow(1).Shrink(1).Equal(size), fmt.Sprintf("%s: ", size))
		}
	})
	t.Run("at least and at most bound each other", func(t *testing.T) {
		for _, a := range sizeFixtures {
			for _, b := range sizeFixtures {
				assert.True(t, a.AtLeast(b).Equal(b.AtLeast(a)), fmt.Sprintf("%s → %s: ", a, b))
				assert.True(t, a.AtMost(b).Equal(b.AtMost(a)), fmt.Sprintf("%s → %s: ", a, b))

				atLeast, atMost := a.AtLeast(b), a.AtMost(b)
				assert.True(t, atMost.Width <= atLeast.Width, fmt.Sprintf("%s → %s: ", a, b))
				assert.True(t, atMost.Height <= atLeast.Height, fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("string round-trips through parse", func(t *testing.T) {
		assertSizeRoundTrip(t, Sz(16, 32))
		assertSizeRoundTrip(t, Sz(0, -34))
		assertSizeRoundTrip(t, Sz[int32](7, 9))
		assertSizeRoundTrip(t, Sz(1.2, 3.6))
		assertSizeRoundTrip(t, Sz(100.0, -34.25))
		assertSizeRoundTrip(t, Sz[float32](0.5, -0.75))
		assertSizeRoundTrip(t, Sz[namedInt](5, 7))

		// String keeps two decimals, so only sizes on that grid survive: 1.005 formats as "1.00"
		size, err := ParseSize[float64](Sz(1.005, -34.0000115).String())
		assert.NoError(t, err)
		AssertSize(t, size, Sz(1.0, -34.0))
	})
}

func assertSizeRoundTrip[T Number](t *testing.T, size Size[T]) {
	t.Helper()

	parsed, err := ParseSize[T](size.String())
	if assert.NoError(t, err) {
		assert.True(t, parsed.Equal(size), fmt.Sprintf("%s: ", size))
	}
}

func TestSize_Immutable(t *testing.T) {
	s := Sz(2, 3)

	s.Scale(2)
	s.ScaleXY(2, 3)
	s.Unscale(2)
	s.UnscaleXY(2, 3)
	s.Grow(1)
	s.GrowXY(1, 2)
	s.Shrink(1)
	s.ShrinkXY(1, 2)
	s.AtLeast(Sz(9, 9))
	s.AtMost(Sz(0, 0))

	AssertSize(t, s, Sz(2, 3))
}

// sizeFixtures span square, portrait, landscape, zero, and negative dimensions.
var sizeFixtures = []Size[float64]{
	Sz(0.0, 0.0),
	Sz(1.0, 2.0),
	Sz(7.0, 7.0),
	Sz(0.4, -0.25),
	Sz(12.5, -0.1),
	Sz(-3.5, 0.25),
	Sz(100.0, -34.25),
}

func ExampleSz() {
	fmt.Println(Sz(16, 32))
	fmt.Println(Sz(1.5, -0.25))
	// Output:
	// 16x32
	// 1.50x-0.25
}

func ExampleParseSize() {
	size, err := ParseSize[int]("16x32")
	fmt.Println(size, err)
	// Output: 16x32 <nil>
}

package geom

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	t.Run("the parse error is wrapped", func(t *testing.T) {
		_, err := ParseSize[int8]("300x1")
		assert.True(t, errors.Is(err, strconv.ErrRange))

		_, err = ParseSize[int]("1xb")
		assert.True(t, errors.Is(err, strconv.ErrSyntax))
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

func TestSize_Area(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Sz(5, 3).Area(), 15)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Sz(0.4, 0.25).Area(), 0.1)
	})
}

func TestSize_Perimeter(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Sz(5, 3).Perimeter(), 16)
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Sz(0.4, 0.25).Perimeter(), 1.3)
	})
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
	t.Run("zero factor panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Sz(10, 20).Unscale(0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Sz(10, 20).UnscaleXY(0, 2)
		}, "geom: division by zero")
	})
}

func TestSize_Abs(t *testing.T) {
	AssertSize(t, Sz(-2, 3).Abs(), Sz(2, 3))
	AssertSize(t, Sz(-0.4, -0.25).Abs(), Sz(0.4, 0.25))
	AssertSize(t, Sz(2, 3).Abs(), Sz(2, 3))
}

func TestSize_Transpose(t *testing.T) {
	t.Run("swaps width and height", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Transpose(), Sz(3, 2))
		AssertSize(t, Sz(0.4, -0.25).Transpose(), Sz(-0.25, 0.4))
	})
	t.Run("twice is the identity and matches Axis.Size", func(t *testing.T) {
		for _, s := range sizeFixtures {
			AssertSize(t, s.Transpose().Transpose(), s, s.String())
			AssertSize(t, s.Transpose(), AxisVertical.Size(s.Width, s.Height), s.String())
		}
	})
}

func TestSize_Round(t *testing.T) {
	AssertSize(t, Sz(1.4, 2.5).Round(), Sz(1.0, 3.0))
	AssertSize(t, Sz(-1.4, -2.5).Round(), Sz(-1.0, -3.0))
	AssertSize(t, Sz(1, 2).Round(), Sz(1, 2))
}

func TestSize_Floor(t *testing.T) {
	AssertSize(t, Sz(1.4, 2.5).Floor(), Sz(1.0, 2.0))
	AssertSize(t, Sz(-1.4, -2.5).Floor(), Sz(-2.0, -3.0))
	AssertSize(t, Sz(1, 2).Floor(), Sz(1, 2))
}

func TestSize_Ceil(t *testing.T) {
	AssertSize(t, Sz(1.4, 2.5).Ceil(), Sz(2.0, 3.0))
	AssertSize(t, Sz(-1.4, -2.5).Ceil(), Sz(-1.0, -2.0))
	AssertSize(t, Sz(1, 2).Ceil(), Sz(1, 2))
}

func TestSize_Lerp(t *testing.T) {
	t.Run("interpolates between the sizes", func(t *testing.T) {
		AssertSize(t, Sz(0.0, 10.0).Lerp(Sz(10.0, 20.0), 0.25), Sz(2.5, 12.5))
		AssertSize(t, Sz(0, 10).Lerp(Sz(10, 20), 0.25), Sz(3, 13)) // int: 2.5 rounds away from zero
	})
	t.Run("the ends are the sizes themselves", func(t *testing.T) {
		AssertSize(t, Sz(2.0, 3.0).Lerp(Sz(8.0, 9.0), 0), Sz(2.0, 3.0))
		AssertSize(t, Sz(2.0, 3.0).Lerp(Sz(8.0, 9.0), 1), Sz(8.0, 9.0))
	})
	t.Run("extrapolates outside the unit range", func(t *testing.T) {
		AssertSize(t, Sz(2.0, 3.0).Lerp(Sz(4.0, 5.0), 2), Sz(6.0, 7.0))
		AssertSize(t, Sz(2.0, 3.0).Lerp(Sz(4.0, 5.0), -1), Sz(0.0, 1.0))
	})
}

func TestSize_Grow(t *testing.T) {
	t.Run("uniform amount", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Grow(2), Sz(4, 5))
		AssertSize(t, Sz(0.4, 0.25).Grow(0.1), Sz(0.5, 0.35))
	})
	t.Run("signed, so nothing is clamped", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Grow(-5), Sz(-3, -2))
		AssertSize(t, Sz(-10, 5).GrowXY(2, -1), Sz(-8, 4))
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
	t.Run("signed, so nothing is clamped", func(t *testing.T) {
		AssertSize(t, Sz(2, 3).Shrink(5), Sz(-3, -2))
		AssertSize(t, Sz(2, 3).ShrinkXY(5, 1), Sz(-3, 2))
		AssertSize(t, Sz(-10, 5).Shrink(2), Sz(-12, 3))
		AssertSize(t, Sz(0.4, 0.25).Shrink(1.0), Sz(-0.6, -0.75))
	})
}

func TestSize_Fit(t *testing.T) {
	t.Run("wide into square is limited by width", func(t *testing.T) {
		AssertSize(t, Sz(1920, 1080).Fit(SzU(960)), Sz(960, 540))
	})
	t.Run("tall into square is limited by height", func(t *testing.T) {
		AssertSize(t, Sz(300, 600).Fit(SzU(200)), Sz(100, 200))
	})
	t.Run("scales up as well as down", func(t *testing.T) {
		AssertSize(t, Sz(4.0, 2.0).Fit(Sz(10.0, 10.0)), Sz(10.0, 5.0))
	})
	t.Run("a zero extent fits as the zero size", func(t *testing.T) {
		AssertSize(t, Sz(0, 5).Fit(SzU(10)), Sz(0, 0))
		AssertSize(t, Sz(5, 0).Fit(SzU(10)), Sz(0, 0))
	})
	t.Run("int rounds", func(t *testing.T) {
		AssertSize(t, Sz(3, 2).Fit(SzU(4)), Sz(4, 3))
	})
	t.Run("a negative extent is out of contract and stays negative", func(t *testing.T) {
		AssertSize(t, Sz(-4.0, 2.0).Fit(SzU(10.0)), Sz(10.0, -5.0))
		AssertSize(t, Sz(-4.0, 2.0).Abs().Fit(SzU(10.0)), Sz(10.0, 5.0))
	})
	t.Run("fits within the target and keeps the ratio", func(t *testing.T) {
		for _, s := range positiveSizeFixtures() {
			for _, target := range positiveSizeFixtures() {
				fitted := s.Fit(target)

				assert.True(t, LessOrEqual(fitted.Width, target.Width) && LessOrEqual(fitted.Height, target.Height), fmt.Sprintf("%s into %s: %s", s, target, fitted))
				if !s.IsZero() && !fitted.IsZero() {
					AssertNumber(t, fitted.AspectRatio(), s.AspectRatio(), fmt.Sprintf("%s into %s: ", s, target))
				}
			}
		}
	})
}

func TestSize_Fill(t *testing.T) {
	t.Run("wide over square is limited by height", func(t *testing.T) {
		AssertSize(t, Sz(1920, 1080).Fill(SzU(540)), Sz(960, 540))
	})
	t.Run("tall over square is limited by width", func(t *testing.T) {
		AssertSize(t, Sz(300, 600).Fill(SzU(200)), Sz(200, 400))
	})
	t.Run("scales down as well as up", func(t *testing.T) {
		AssertSize(t, Sz(40.0, 20.0).Fill(Sz(10.0, 10.0)), Sz(20.0, 10.0))
	})
	t.Run("a zero extent fills as the zero size", func(t *testing.T) {
		AssertSize(t, Sz(0, 5).Fill(SzU(10)), Sz(0, 0))
	})
	t.Run("a negative extent is out of contract and stays negative", func(t *testing.T) {
		AssertSize(t, Sz(-4.0, 2.0).Fill(SzU(10.0)), Sz(-20.0, 10.0))
		AssertSize(t, Sz(-4.0, 2.0).Abs().Fill(SzU(10.0)), Sz(20.0, 10.0))
	})
	t.Run("covers the target and keeps the ratio", func(t *testing.T) {
		for _, s := range positiveSizeFixtures() {
			for _, target := range positiveSizeFixtures() {
				if s.Width == 0 || s.Height == 0 {
					continue
				}

				filled := s.Fill(target)

				assert.True(t, LessOrEqual(target.Width, filled.Width) && LessOrEqual(target.Height, filled.Height), fmt.Sprintf("%s over %s: %s", s, target, filled))
				if !filled.IsZero() {
					AssertNumber(t, filled.AspectRatio(), s.AspectRatio(), fmt.Sprintf("%s over %s: ", s, target))
				}
			}
		}
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

func TestSize_Vector(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVector(t, Sz(10, 16).Vector(), Vec(10, 16))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Sz(1.5, -2.5).Vector(), Vec(1.5, -2.5))
	})
}

func TestSize_Cast(t *testing.T) {
	s := Sz(1.5, 2.5)

	t.Run("matches Int and Float", func(t *testing.T) {
		AssertSize(t, s.Cast[int](), s.Int())
		AssertSize(t, s.Cast[float64](), s.Float())
	})
	t.Run("a type the other conversions cannot name", func(t *testing.T) {
		AssertSize(t, s.Cast[int8](), Sz[int8](2, 3))
	})
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
		assert.Equal(t, Sz(negativeZero, 0.0).String(), "0.00x0.00")
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
	t.Run("round, floor and ceil bracket the size", func(t *testing.T) {
		for _, size := range sizeFixtures {
			floor, ceil := size.Floor(), size.Ceil()

			assert.True(t, floor.Width <= size.Width && size.Width <= ceil.Width, fmt.Sprintf("%s: ", size))
			assert.True(t, floor.Height <= size.Height && size.Height <= ceil.Height, fmt.Sprintf("%s: ", size))
			round := size.Round()
			assert.True(t, Equal(round.Width, floor.Width) || Equal(round.Width, ceil.Width), fmt.Sprintf("%s: ", size))
			assert.True(t, Equal(round.Height, floor.Height) || Equal(round.Height, ceil.Height), fmt.Sprintf("%s: ", size))
		}
	})
	t.Run("lerp ends on the two sizes", func(t *testing.T) {
		for _, a := range sizeFixtures {
			for _, b := range sizeFixtures {
				assert.True(t, a.Lerp(b, 0).Equal(a), fmt.Sprintf("%s → %s: ", a, b))
				assert.True(t, a.Lerp(b, 1).Equal(b), fmt.Sprintf("%s → %s: ", a, b))
				assert.True(t, a.Lerp(b, 0.5).Equal(b.Lerp(a, 0.5)), fmt.Sprintf("%s → %s: ", a, b))
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

// positiveSizeFixtures are the fixtures with no negative extent, the sizes Fit and Fill are defined for.
func positiveSizeFixtures() []Size[float64] {
	var sizes []Size[float64]
	for _, s := range sizeFixtures {
		if s.Width >= 0 && s.Height >= 0 {
			sizes = append(sizes, s)
		}
	}

	return sizes
}

package geom_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

// negativeZero is -0.0; writing it as a literal would fold to +0.0 at compile time.
var negativeZero = math.Copysign(0, -1)

// tenth is 0.1 read at run time: a test pinning exact bits takes an operand from it, since the
// compiler folds arithmetic on literals into a constant, computed without the fused multiply-add
// the test guards against.
var tenth = 0.1

// sinkPoints is a shared slice variable used to store results in allocation tests, so the result
// escapes as it does for any caller.
var sinkPoints []Point[int]

// sinkSegments is the sinkPoints of the methods returning segments.
var sinkSegments []Segment[int]

// sinkInts is the sinkPoints of the methods returning indices.
var sinkInts []int

// sinkBool is a shared boolean variable used to store results in benchmarking tests.
var sinkBool bool

// sinkString is the sinkPoints of the String methods.
var sinkString string

// assertStringAllocation asserts that the String method allocates the string it returns and
// nothing else.
func assertStringAllocation(t *testing.T, stringer func() string) {
	t.Helper()

	geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
		sinkString = stringer()
	}), 1)
}

func TestMultiply(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertNumber(t, Multiply(5, 2.0), 10)
		geomtest.AssertNumber(t, Multiply(3, 0.5), 2) // int: 1.5 rounds to 2
		geomtest.AssertNumber(t, Multiply(-4, 1.5), -6)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Multiply(4.0, 2.5), 10.0)
		geomtest.AssertNumber(t, Multiply(0.4, 0.5), 0.2)
	})
}

func TestDivide(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertNumber(t, Divide(10, 2.0), 5)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Divide(7.5, 2.5), 3.0)
		geomtest.AssertNumber(t, Divide(1.0, 3.0), 1.0/3.0)
	})
	t.Run("zero scale panics", func(t *testing.T) {
		assert.Panics(t, func() {
			Divide(10, 0.0)
		}, "geom: division by zero")
		assert.Panics(t, func() {
			Divide(2.5, 0.0)
		}, "geom: division by zero")
	})
}

func TestAbs(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		geomtest.AssertNumber(t, Abs(-5), 5)
		geomtest.AssertNumber(t, Abs(-3.14), 3.14)
	})
	t.Run("non-negative is unchanged", func(t *testing.T) {
		geomtest.AssertNumber(t, Abs(5), 5)
		geomtest.AssertNumber(t, Abs(0), 0)
		geomtest.AssertNumber(t, Abs(3.14), 3.14)
	})
	t.Run("wide integers stay exact", func(t *testing.T) {
		geomtest.AssertNumber(t, Abs(int64(-(1<<53 + 1))), int64(1<<53+1))
		geomtest.AssertNumber(t, Abs(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestSign(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Sign(5), 1)
		geomtest.AssertNumber(t, Sign(-5), -1)
		geomtest.AssertNumber(t, Sign(0), 0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Sign(3.14), 1.0)
		geomtest.AssertNumber(t, Sign(-3.14), -1.0)
		geomtest.AssertNumber(t, Sign(0.0), 0.0)
	})
}

func TestRound(t *testing.T) {
	t.Run("float rounds half away from zero", func(t *testing.T) {
		geomtest.AssertNumber(t, Round(1.4), 1.0)
		geomtest.AssertNumber(t, Round(1.5), 2.0)
		geomtest.AssertNumber(t, Round(-1.5), -2.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertNumber(t, Round(3), 3)
		geomtest.AssertNumber(t, Round(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestRoundHalfUp(t *testing.T) {
	t.Run("float rounds to the nearest integer", func(t *testing.T) {
		geomtest.AssertNumber(t, RoundHalfUp(1.4), 1.0)
		geomtest.AssertNumber(t, RoundHalfUp(1.6), 2.0)
		geomtest.AssertNumber(t, RoundHalfUp(-1.4), -1.0)
		geomtest.AssertNumber(t, RoundHalfUp(-1.6), -2.0)
	})
	t.Run("a tie rounds toward positive infinity", func(t *testing.T) {
		geomtest.AssertNumber(t, RoundHalfUp(0.5), 1.0)
		geomtest.AssertNumber(t, RoundHalfUp(1.5), 2.0)
		geomtest.AssertNumber(t, RoundHalfUp(-0.5), 0.0)
		geomtest.AssertNumber(t, RoundHalfUp(-1.5), -1.0)
	})
	t.Run("the value just below a tie rounds down", func(t *testing.T) {
		geomtest.AssertNumber(t, RoundHalfUp(math.Nextafter(0.5, 0)), 0.0)
		geomtest.AssertNumber(t, RoundHalfUp(math.Nextafter(-0.5, -1)), -1.0)
	})
	t.Run("shifts with an integer", func(t *testing.T) {
		for _, x := range []float64{0, 0.25, 0.5, 0.75} {
			for shift := -8.0; shift <= 8; shift++ {
				geomtest.AssertNumber(t, RoundHalfUp(x+shift), RoundHalfUp(x)+shift)
			}
		}
	})
	t.Run("float32", func(t *testing.T) {
		geomtest.AssertNumber(t, RoundHalfUp(float32(-2.5)), float32(-2))
	})
	t.Run("a non-finite float is returned as it is", func(t *testing.T) {
		assert.True(t, math.IsNaN(RoundHalfUp(math.NaN())))
		assert.True(t, math.IsInf(RoundHalfUp(math.Inf(-1)), -1))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertNumber(t, RoundHalfUp(3), 3)
		geomtest.AssertNumber(t, RoundHalfUp(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestFloor(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Floor(1.9), 1.0)
		geomtest.AssertNumber(t, Floor(-1.1), -2.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertNumber(t, Floor(3), 3)
		geomtest.AssertNumber(t, Floor(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestCeil(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Ceil(1.1), 2.0)
		geomtest.AssertNumber(t, Ceil(-1.9), -1.0)
	})
	t.Run("int is a no-op", func(t *testing.T) {
		geomtest.AssertNumber(t, Ceil(3), 3)
		geomtest.AssertNumber(t, Ceil(int64(1<<53+1)), int64(1<<53+1))
	})
}

func TestSnap(t *testing.T) {
	t.Run("float moves to the nearest multiple", func(t *testing.T) {
		geomtest.AssertNumber(t, Snap(7.4, 2.5), 7.5)
		geomtest.AssertNumber(t, Snap(-3.6, 2.5), -2.5)
		geomtest.AssertNumber(t, Snap(0.2, 0.25), 0.25)
	})
	t.Run("a tie rounds toward positive infinity", func(t *testing.T) {
		geomtest.AssertNumber(t, Snap(5.0, 2.0), 6.0)
		geomtest.AssertNumber(t, Snap(-5.0, 2.0), -4.0)
		geomtest.AssertNumber(t, Snap(1.0, 2.0), 2.0)
		geomtest.AssertNumber(t, Snap(-1.0, 2.0), 0.0)
	})
	t.Run("shifts with a multiple of the step", func(t *testing.T) {
		for _, x := range []float64{0, 0.75, 1.5, 2.25} {
			for shift := -8.0; shift <= 8; shift++ {
				geomtest.AssertNumber(t, Snap(x+shift*3, 3), Snap(x, 3)+shift*3)
			}
		}
	})
	t.Run("a negative step gives the same lattice", func(t *testing.T) {
		for _, x := range []float64{-5, -3.6, -1, 0, 1, 2.4, 5, 7.4} {
			geomtest.AssertNumber(t, Snap(x, -2.5), Snap(x, 2.5))
		}
		geomtest.AssertNumber(t, Snap(5, -2), 6)
	})
	t.Run("a zero step leaves the value", func(t *testing.T) {
		geomtest.AssertNumber(t, Snap(7.4, 0), 7.4)
		geomtest.AssertNumber(t, Snap(7, 0), 7)
	})
	t.Run("int", func(t *testing.T) {
		geomtest.AssertNumber(t, Snap(17, 8), 16)
		geomtest.AssertNumber(t, Snap(20, 8), 24)
		geomtest.AssertNumber(t, Snap(-20, 8), -16)
		geomtest.AssertNumber(t, Snap(-21, 8), -24)
	})
	t.Run("float32", func(t *testing.T) {
		geomtest.AssertNumber(t, Snap(float32(7.4), 2.5), float32(7.5))
	})
	t.Run("a non-finite float is returned as it is", func(t *testing.T) {
		assert.True(t, math.IsNaN(Snap(math.NaN(), 2)))
		assert.True(t, math.IsInf(Snap(math.Inf(-1), 2), -1))
	})
}

func ExampleSnap() {
	fmt.Println(Snap(37, 16))
	fmt.Println(Snap(7.4, 2.5))
	// Output:
	// 32
	// 7.5
}

func TestMod(t *testing.T) {
	t.Run("non-negative operands", func(t *testing.T) {
		geomtest.AssertNumber(t, Mod(5, 8), 5)
		geomtest.AssertNumber(t, Mod(8, 8), 0)
		geomtest.AssertNumber(t, Mod(9, 8), 1)
		geomtest.AssertNumber(t, Mod(0, 8), 0)
		geomtest.AssertNumber(t, Mod(7, 6), 1)
	})
	t.Run("negative operands stay in range", func(t *testing.T) {
		geomtest.AssertNumber(t, Mod(-1, 8), 7)
		geomtest.AssertNumber(t, Mod(-8, 8), 0)
		geomtest.AssertNumber(t, Mod(-9, 8), 7)
		geomtest.AssertNumber(t, Mod(-1, 6), 5)
	})
	t.Run("a negative modulus wraps into (m, 0]", func(t *testing.T) {
		geomtest.AssertNumber(t, Mod(3, -5), -2)
		geomtest.AssertNumber(t, Mod(-3, -5), -3)
		geomtest.AssertNumber(t, Mod(5, -5), 0)
	})
	t.Run("defined integer types", func(t *testing.T) {
		geomtest.AssertNumber(t, Mod(namedInt(12), 8), namedInt(4))
	})
	t.Run("a narrow T near the end of its range does not overflow", func(t *testing.T) {
		geomtest.AssertNumber(t, Mod[int8](99, 100), 99)
		geomtest.AssertNumber(t, Mod[int8](-1, 100), 99)
		geomtest.AssertNumber(t, Mod[int8](-128, 127), 126)
		geomtest.AssertNumber(t, Mod[int8](100, -128), -28)
	})
}

func TestIsFinite(t *testing.T) {
	t.Run("a number", func(t *testing.T) {
		assert.True(t, IsFinite(1.5))
		assert.True(t, IsFinite(negativeZero))
		assert.True(t, IsFinite(-math.MaxFloat64))
		assert.True(t, IsFinite(math.SmallestNonzeroFloat64))
	})
	t.Run("NaN and the infinities are not", func(t *testing.T) {
		assert.False(t, IsFinite(math.NaN()))
		assert.False(t, IsFinite(math.Inf(1)))
		assert.False(t, IsFinite(math.Inf(-1)))
	})
	t.Run("float32", func(t *testing.T) {
		assert.True(t, IsFinite[float32](math.MaxFloat32))
		assert.False(t, IsFinite(float32(math.NaN())))
		assert.False(t, IsFinite(float32(math.Inf(-1))))
	})
	t.Run("an integer always is", func(t *testing.T) {
		assert.True(t, IsFinite(math.MaxInt))
		assert.True(t, IsFinite[int8](math.MinInt8))
	})
	t.Run("a product is not fused into the check", func(t *testing.T) {
		assert.True(t, IsFinite(tenth*3))
		assert.False(t, IsFinite(tenth*math.Inf(1)))
	})
}

func TestClamp(t *testing.T) {
	t.Run("inside the range", func(t *testing.T) {
		geomtest.AssertNumber(t, Clamp(5, 0, 10), 5)
		geomtest.AssertNumber(t, Clamp(0.5, 0.0, 1.0), 0.5)
	})
	t.Run("outside the range", func(t *testing.T) {
		geomtest.AssertNumber(t, Clamp(-5, 0, 10), 0)
		geomtest.AssertNumber(t, Clamp(15, 0, 10), 10)
		geomtest.AssertNumber(t, Clamp(-0.5, 0.0, 1.0), 0.0)
		geomtest.AssertNumber(t, Clamp(1.5, 0.0, 1.0), 1.0)
	})
	t.Run("the bounds are inclusive", func(t *testing.T) {
		geomtest.AssertNumber(t, Clamp(0, 0, 10), 0)
		geomtest.AssertNumber(t, Clamp(10, 0, 10), 10)
	})
}

func TestSum(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Sum([]int{1, 2, 3}), 6)
		assert.Equal(t, Sum([]int{-4, 1}), -3)
	})
	t.Run("narrow integers do not overflow mid-sum", func(t *testing.T) {
		assert.Equal(t, Sum([]int8{100, 100, -100}), int8(100))
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Sum([]float64{0.5, 0.25, 0.125}), 0.875)
		geomtest.AssertNumber(t, Sum([]float32{0.5, 0.25}), float32(0.75))
	})
	t.Run("empty and nil are zero", func(t *testing.T) {
		assert.Equal(t, Sum([]int{}), 0)
		assert.Equal(t, Sum[float64](nil), 0.0)
	})
}

func TestLerp(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertNumber(t, Lerp(1, 2, 0.25), 1)
		geomtest.AssertNumber(t, Lerp(1, 3, 0.25), 2)
		geomtest.AssertNumber(t, Lerp(1, 5, 0.25), 2)
		geomtest.AssertNumber(t, Lerp(1, 7, 0.25), 3)
	})
	t.Run("rounds each product, the same bits on every architecture", func(t *testing.T) {
		assert.Equal(t, Lerp(tenth, 0.2, 0.1), 0.11000000000000001)
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		geomtest.AssertNumber(t, Lerp[int8](-100, 100, 0.5), 0)
		geomtest.AssertNumber(t, Lerp[int8](-128, 127, 1), 127)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Lerp(1.0, 6.0, 0.25), 2.25)
		geomtest.AssertNumber(t, Lerp(1.0, 6.0, 0.75), 4.75)
	})
	t.Run("endpoints", func(t *testing.T) {
		geomtest.AssertNumber(t, Lerp(1.0, 6.0, 0), 1.0)
		geomtest.AssertNumber(t, Lerp(1.0, 6.0, 1), 6.0)
	})
}

func TestInverseLerp(t *testing.T) {
	t.Run("the ratio between the ends", func(t *testing.T) {
		geomtest.AssertNumber(t, InverseLerp(10.0, 20.0, 12.5), 0.25)
		geomtest.AssertNumber(t, InverseLerp(0, 120, 30), 0.25)
		geomtest.AssertNumber(t, InverseLerp[float32](-1, 1, 0.5), 0.75)
	})
	t.Run("the ends are zero and one", func(t *testing.T) {
		assert.Equal(t, InverseLerp(tenth, 3.3, tenth), 0.0)
		assert.Equal(t, InverseLerp(tenth, 3.3, 3.3), 1.0)
	})
	t.Run("a range running downward", func(t *testing.T) {
		geomtest.AssertNumber(t, InverseLerp(20.0, 10.0, 12.5), 0.75)
	})
	t.Run("is not clamped outside the range", func(t *testing.T) {
		geomtest.AssertNumber(t, InverseLerp(10.0, 20.0, 40.0), 3.0)
		geomtest.AssertNumber(t, InverseLerp(10.0, 20.0, 0.0), -1.0)
	})
	t.Run("an empty range gives zero", func(t *testing.T) {
		assert.Equal(t, InverseLerp(5.0, 5.0, 5.0), 0.0)
		assert.Equal(t, InverseLerp(5.0, 5.0, 9.0), 0.0)
		assert.Equal(t, InverseLerp(5, 5, 9), 0.0)
	})
	t.Run("the differences are taken in float64", func(t *testing.T) {
		geomtest.AssertNumber(t, InverseLerp[int8](-100, 100, 50), 0.75)
	})
	t.Run("undoes Lerp", func(t *testing.T) {
		for _, a := range []float64{0, 1, -2.5, 100} {
			for _, b := range []float64{3, -7, 0.5, 1e3} {
				for _, ratio := range []float64{0, 0.25, 0.5, 1, -1, 2} {
					geomtest.AssertNumber(t, InverseLerp(a, b, Lerp(a, b, ratio)), ratio, fmt.Sprintf("%v → %v at %v: ", a, b, ratio))
				}
			}
		}
	})
}

func TestRemap(t *testing.T) {
	t.Run("keeps the ratio between the ranges", func(t *testing.T) {
		geomtest.AssertNumber(t, Remap(12.5, 10.0, 20.0, 0.0, 100.0), 25.0)
		geomtest.AssertNumber(t, Remap(0.5, -1.0, 1.0, -30.0, 30.0), 15.0)
		geomtest.AssertNumber(t, Remap[float32](5, 0, 10, 1, 2), 1.5)
	})
	t.Run("the ends map to the ends", func(t *testing.T) {
		assert.Equal(t, Remap(tenth, tenth, 3.3, 7.0, 9.0), 7.0)
		assert.Equal(t, Remap(3.3, tenth, 3.3, 7.0, 9.0), 9.0)
	})
	t.Run("a range running downward", func(t *testing.T) {
		geomtest.AssertNumber(t, Remap(400.0, 200.0, 1000.0, 1.0, 0.0), 0.75)
		geomtest.AssertNumber(t, Remap(12.5, 20.0, 10.0, 0.0, 100.0), 75.0)
	})
	t.Run("is not clamped outside the range", func(t *testing.T) {
		geomtest.AssertNumber(t, Remap(40.0, 10.0, 20.0, 0.0, 100.0), 300.0)
		geomtest.AssertNumber(t, Remap(0.0, 10.0, 20.0, 0.0, 100.0), -100.0)
	})
	t.Run("an empty first range gives the start of the second", func(t *testing.T) {
		assert.Equal(t, Remap(9.0, 5.0, 5.0, 7.0, 8.0), 7.0)
	})
	t.Run("int rounds once", func(t *testing.T) {
		geomtest.AssertNumber(t, Remap(1, 0, 3, 0, 10), 3)
		geomtest.AssertNumber(t, Remap(2, 0, 3, 0, 10), 7)
	})
	t.Run("back again is the value", func(t *testing.T) {
		for _, value := range []float64{0, 1, -2.5, 12, 100} {
			geomtest.AssertNumber(t, Remap(Remap(value, 1.0, 9.0, -4.0, 0.5), -4.0, 0.5, 1.0, 9.0), value, fmt.Sprintf("%v: ", value))
		}
	})
}

func TestMidpoint(t *testing.T) {
	t.Run("int rounds", func(t *testing.T) {
		geomtest.AssertNumber(t, Midpoint(1, 3), 2)
		geomtest.AssertNumber(t, Midpoint(1, 4), 3)
		geomtest.AssertNumber(t, Midpoint(1, 5), 3)
		geomtest.AssertNumber(t, Midpoint(1, 6), 4)
		geomtest.AssertNumber(t, Midpoint(1, 7), 4)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Midpoint(1.0, 6.0), 3.5)
	})
	t.Run("is lerp at one half", func(t *testing.T) {
		for _, pair := range [][2]float64{{1, 6}, {-3.5, 0.25}, {7, 7}, {12.5, -0.1}} {
			geomtest.AssertNumber(t, Midpoint(pair[0], pair[1]), Lerp(pair[0], pair[1], 0.5))
		}
	})
}

func TestMoveToward(t *testing.T) {
	t.Run("moves by the step", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveToward(1.0, 4.0, 0.5), 1.5)
		geomtest.AssertNumber(t, MoveToward(4.0, 1.0, 0.5), 3.5)
		geomtest.AssertNumber(t, MoveToward(1, 10, 3), 4)
		geomtest.AssertNumber(t, MoveToward[float32](1, -4, 0.5), 0.5)
	})
	t.Run("a step reaching the target gives the target itself", func(t *testing.T) {
		assert.Equal(t, MoveToward(tenth, 3.3, 3.2), 3.3)
		assert.Equal(t, MoveToward(tenth, 3.3, 100), 3.3)
		assert.Equal(t, MoveToward(10, 1, 9), 1)
	})
	t.Run("a step of zero stays", func(t *testing.T) {
		assert.Equal(t, MoveToward(tenth, 3.3, 0), tenth)
	})
	t.Run("a negative step moves away and is not clamped", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveToward(1.0, 4.0, -5), -4.0)
		geomtest.AssertNumber(t, MoveToward(4.0, 1.0, -5), 9.0)
	})
	t.Run("a value at the target stays for any step", func(t *testing.T) {
		assert.Equal(t, MoveToward(1.0, 1.0, 3), 1.0)
		assert.Equal(t, MoveToward(1.0, 1.0, -3), 1.0)
	})
	t.Run("int rounds, and a step below half a unit does not move", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveToward(0, 10, 2.5), 3)
		geomtest.AssertNumber(t, MoveToward(0, -10, 2.5), -3)
		geomtest.AssertNumber(t, MoveToward(0, 10, 0.4), 0)
	})
	t.Run("the difference is taken in float64", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveToward[int8](-100, 100, 50), int8(-50))
	})
	t.Run("a NaN gives NaN", func(t *testing.T) {
		assert.True(t, math.IsNaN(MoveToward(0.0, math.NaN(), 1)))
		assert.True(t, math.IsNaN(MoveToward(0.0, 1.0, math.NaN())))
	})
}

func TestEqual(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, Equal(1, 1))
		assert.False(t, Equal(1, 2))
	})
	t.Run("float64 within Delta", func(t *testing.T) {
		assert.True(t, Equal(1.0000005, 1.0))
		assert.False(t, Equal(1.0000015, 1.0))
		assert.True(t, Equal(0.0, 0.0))
	})
	t.Run("float32 within Delta32", func(t *testing.T) {
		assert.True(t, Equal[float32](1.00005, 1.0))
		assert.False(t, Equal[float32](1.0005, 1.0))
	})
	t.Run("absolute, so it degenerates far from zero", func(t *testing.T) {
		// one float32 ulp at 1e4 is 9.8e-4, above Delta32, so Equal asks for bit-exact
		// equality from here up. This is the documented cost of the fast path.
		const far float32 = 1e4
		next := math.Nextafter32(far, math.MaxFloat32)

		assert.False(t, Equal(far, next))
		assert.True(t, EqualRelative(far, next))
	})
}

func BenchmarkEqual_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(float64(i), float64(i))
	}
}

func BenchmarkEqual_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(float32(i), float32(i))
	}
}

func TestEqualDelta(t *testing.T) {
	t.Run("within the delta", func(t *testing.T) {
		assert.True(t, EqualDelta(1, 2, 1.5))
		assert.True(t, EqualDelta(1.0, 1.001, 0.01))
	})
	t.Run("outside the delta", func(t *testing.T) {
		assert.False(t, EqualDelta(1, 3, 1.5))
		assert.False(t, EqualDelta(1.0, 1.02, 0.01))
	})
	t.Run("a zero delta asks for exact equality", func(t *testing.T) {
		assert.True(t, EqualDelta(5, 5, 0.0))
		assert.False(t, EqualDelta(5, 6, 0.0))
	})
	t.Run("an infinity equals itself", func(t *testing.T) {
		assert.True(t, EqualDelta(math.Inf(1), math.Inf(1), Delta))
		assert.False(t, EqualDelta(math.Inf(1), math.Inf(-1), Delta))
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		assert.False(t, EqualDelta[int8](127, -128, 1))
		assert.True(t, EqualDelta[int8](127, -128, 255))
	})
}

func BenchmarkEqualDelta_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = EqualDelta(float64(i), float64(i), Delta)
	}
}

func TestEqualRelative(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assert.True(t, EqualRelative(1, 1))
		assert.False(t, EqualRelative(1, 2))
	})
	t.Run("near zero it matches Equal", func(t *testing.T) {
		assert.True(t, EqualRelative(1.0000005, 1.0))
		assert.False(t, EqualRelative(1.0000015, 1.0))
		assert.True(t, EqualRelative(0.0, 0.0))
	})
	t.Run("scales with magnitude", func(t *testing.T) {
		assert.True(t, EqualRelative(1e6, 1e6+0.5))
		assert.False(t, EqualRelative(1e6, 1e6+2.0))
		assert.True(t, EqualRelative[float32](1e4, 1e4+0.05))
	})
	t.Run("a full turn in float32 returns to its start", func(t *testing.T) {
		v := Vec[float32](1000, 0)
		for range 8 {
			v = v.Rotate(Pi / 4)
		}

		// one ulp out at magnitude 1e3, which Delta32 still covers
		assert.True(t, EqualRelative(v.X, 1000))
		assert.True(t, Equal(v.X, 1000))
	})
}

func BenchmarkEqualRelative_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = EqualRelative(float64(i), float64(i))
	}
}

func TestLessOrEqual(t *testing.T) {
	t.Run("int compares exactly", func(t *testing.T) {
		assert.True(t, LessOrEqual(1, 2))
		assert.True(t, LessOrEqual(2, 2))
		assert.False(t, LessOrEqual(3, 2))
	})
	t.Run("int64 beyond 2^53 stays exact", func(t *testing.T) {
		big := int64(1) << 53

		assert.False(t, LessOrEqual(big+1, big))
		assert.True(t, LessOrEqual(big, big+1))
	})
	t.Run("float64 allows Delta past the bound", func(t *testing.T) {
		assert.True(t, LessOrEqual(1.0, 2.0))
		assert.True(t, LessOrEqual(2.0+Delta/2, 2.0))
		assert.False(t, LessOrEqual(2.0+2*Delta, 2.0))
	})
	t.Run("float32 allows Delta32 past the bound", func(t *testing.T) {
		assert.True(t, LessOrEqual(float32(2+Delta32/2), 2))
		assert.False(t, LessOrEqual(float32(2+2*Delta32), 2))
	})
	t.Run("a recomputed float corner is still inside", func(t *testing.T) {
		low := 0.1
		recomputed := (low + 0.6/2) - 0.6/2

		assert.False(t, low >= recomputed)
		assert.True(t, LessOrEqual(recomputed, low))
	})
}

func TestLessOrEqualDelta(t *testing.T) {
	t.Run("within the delta", func(t *testing.T) {
		assert.True(t, LessOrEqualDelta(1, 2, 0.0))
		assert.True(t, LessOrEqualDelta(3, 2, 1.5))
		assert.True(t, LessOrEqualDelta(1.001, 1.0, 0.01))
	})
	t.Run("outside the delta", func(t *testing.T) {
		assert.False(t, LessOrEqualDelta(4, 2, 1.5))
		assert.False(t, LessOrEqualDelta(1.02, 1.0, 0.01))
	})
	t.Run("narrow integers do not overflow", func(t *testing.T) {
		assert.True(t, LessOrEqualDelta[int8](-128, 127, 0))
		assert.False(t, LessOrEqualDelta[int8](127, -128, 254))
		assert.True(t, LessOrEqualDelta[int8](127, -128, 255))
	})
}

func TestEpsilon(t *testing.T) {
	t.Run("int compares exactly", func(t *testing.T) {
		geomtest.AssertNumber(t, Epsilon[int](), 0.0)
		geomtest.AssertNumber(t, Epsilon[int8](), 0.0)
		geomtest.AssertNumber(t, Epsilon[int16](), 0.0)
		geomtest.AssertNumber(t, Epsilon[int32](), 0.0)
		geomtest.AssertNumber(t, Epsilon[int64](), 0.0)
		geomtest.AssertNumber(t, Epsilon[namedInt](), 0.0)
		geomtest.AssertNumber(t, Epsilon[namedInt8](), 0.0)
	})
	t.Run("float", func(t *testing.T) {
		geomtest.AssertNumber(t, Epsilon[float32](), Delta32)
		geomtest.AssertNumber(t, Epsilon[float64](), Delta)
	})
	t.Run("defined types follow their underlying kind", func(t *testing.T) {
		geomtest.AssertNumber(t, Epsilon[namedFloat32](), Delta32)
		geomtest.AssertNumber(t, Epsilon[namedFloat64](), Delta)
	})
}

func TestEpsilonRelative(t *testing.T) {
	t.Run("floored at Epsilon near zero", func(t *testing.T) {
		geomtest.AssertNumber(t, EpsilonRelative(0.0, 0.5), Delta)
	})
	t.Run("scales with the larger magnitude", func(t *testing.T) {
		geomtest.AssertNumber(t, EpsilonRelative(100.0, 0.0), 100*Delta)
		geomtest.AssertNumber(t, EpsilonRelative(0.0, -100.0), 100*Delta)
		geomtest.AssertNumber(t, EpsilonRelative[float32](100, 0), 100*Delta32)
	})
	t.Run("int stays exact at any magnitude", func(t *testing.T) {
		geomtest.AssertNumber(t, EpsilonRelative(5, 5), 0.0)
		geomtest.AssertNumber(t, EpsilonRelative(1000000, 0), 0.0)
	})
}

func TestToRadians(t *testing.T) {
	geomtest.AssertNumber(t, ToRadians(0), 0.0)
	geomtest.AssertNumber(t, ToRadians(90), Pi/2)
	geomtest.AssertNumber(t, ToRadians(180), Pi)
	geomtest.AssertNumber(t, ToRadians(360), 2*Pi)
}

func TestToDegrees(t *testing.T) {
	t.Run("converts", func(t *testing.T) {
		geomtest.AssertNumber(t, ToDegrees(0), 0.0)
		geomtest.AssertNumber(t, ToDegrees(Pi/2), 90.0)
		geomtest.AssertNumber(t, ToDegrees(Pi), 180.0)
		geomtest.AssertNumber(t, ToDegrees(2*Pi), 360.0)
	})
	t.Run("inverts ToRadians", func(t *testing.T) {
		for _, degrees := range []float64{0, 30, 45, 90, 179.5, -270} {
			geomtest.AssertNumber(t, ToDegrees(ToRadians(degrees)), degrees)
		}
	})
}

func TestNormalizeAngle(t *testing.T) {
	t.Run("inside the unit turn", func(t *testing.T) {
		geomtest.AssertNumber(t, NormalizeAngle(0), 0.0)
		geomtest.AssertNumber(t, NormalizeAngle(Pi), Pi)
	})
	t.Run("wraps above a full turn", func(t *testing.T) {
		geomtest.AssertNumber(t, NormalizeAngle(3*Pi), Pi)
		geomtest.AssertNumber(t, NormalizeAngle(2*Pi), 0.0)
	})
	t.Run("negative angles come back positive", func(t *testing.T) {
		geomtest.AssertNumber(t, NormalizeAngle(-Pi/2), 3*Pi/2)
	})
	t.Run("a tiny negative angle stays below a full turn", func(t *testing.T) {
		assert.Less(t, NormalizeAngle(-1e-17), 2*Pi)
		assert.GreaterOrEqual(t, NormalizeAngle(-1e-17), 0.0)
	})
	t.Run("NaN and Inf have no normal form", func(t *testing.T) {
		assert.True(t, math.IsNaN(NormalizeAngle(math.NaN())))
		assert.True(t, math.IsNaN(NormalizeAngle(math.Inf(1))))
		assert.True(t, math.IsNaN(NormalizeAngle(math.Inf(-1))))
	})
}

func TestAngleDistance(t *testing.T) {
	t.Run("same angle", func(t *testing.T) {
		geomtest.AssertNumber(t, AngleDistance(1, 1), 0.0)
		geomtest.AssertNumber(t, AngleDistance(1, 1+2*Pi), 0.0)
	})
	t.Run("shortest way around", func(t *testing.T) {
		geomtest.AssertNumber(t, AngleDistance(0, Pi/2), Pi/2)
		geomtest.AssertNumber(t, AngleDistance(Pi/2, 0), Pi/2)
		geomtest.AssertNumber(t, AngleDistance(0, 3*Pi/2), Pi/2)
		geomtest.AssertNumber(t, AngleDistance(0, Pi), Pi)
	})
	t.Run("across the seam", func(t *testing.T) {
		geomtest.AssertNumber(t, AngleDistance(-0.1, 0.1), 0.2)
		geomtest.AssertNumber(t, AngleDistance(2*Pi-0.1, 0.1), 0.2)
	})
}

func TestLerpAngle(t *testing.T) {
	t.Run("turns along the shorter arc", func(t *testing.T) {
		geomtest.AssertNumber(t, LerpAngle(0, Pi/2, 0.5), Pi/4)
		geomtest.AssertNumber(t, LerpAngle(Pi/2, 0, 0.5), Pi/4)
		geomtest.AssertNumber(t, LerpAngle(0, 3*Pi/2, 0.5), -Pi/4)
	})
	t.Run("crosses the seam the short way", func(t *testing.T) {
		geomtest.AssertNumber(t, LerpAngle(ToRadians(350), ToRadians(10), 0.5), ToRadians(360))
		geomtest.AssertNumber(t, LerpAngle(ToRadians(10), ToRadians(350), 0.5), ToRadians(0))
	})
	t.Run("the ends are the angles themselves, up to a turn", func(t *testing.T) {
		assert.True(t, EqualAngle(LerpAngle(1, 4, 0), 1))
		assert.True(t, EqualAngle(LerpAngle(1, 4, 1), 4))
		assert.True(t, EqualAngle(LerpAngle(1, 4+2*Pi, 1), 4))
	})
	t.Run("extrapolates along the same arc", func(t *testing.T) {
		geomtest.AssertNumber(t, LerpAngle(0, Pi/2, 2), Pi)
		geomtest.AssertNumber(t, LerpAngle(0, Pi/2, -1), -Pi/2)
	})
	t.Run("half a turn apart turns by increasing angle", func(t *testing.T) {
		geomtest.AssertNumber(t, LerpAngle(0, Pi, 0.5), Pi/2)
	})
	t.Run("never turns more than half a turn", func(t *testing.T) {
		for _, a := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
			for _, b := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
				assert.True(t, math.Abs(LerpAngle(a, b, 1)-a) <= Pi+Delta, fmt.Sprintf("%v → %v: ", a, b))
				assert.True(t, EqualAngle(LerpAngle(a, b, 1), b), fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
}

func TestMoveTowardAngle(t *testing.T) {
	t.Run("turns by the step along the shorter arc", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveTowardAngle(0, Pi/2, 0.25), 0.25)
		geomtest.AssertNumber(t, MoveTowardAngle(Pi/2, 0, 0.25), Pi/2-0.25)
		geomtest.AssertNumber(t, MoveTowardAngle(0, 3*Pi/2, 0.25), -0.25)
	})
	t.Run("crosses the seam the short way", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveTowardAngle(ToRadians(350), ToRadians(10), ToRadians(15)), ToRadians(365))
		geomtest.AssertNumber(t, MoveTowardAngle(ToRadians(10), ToRadians(350), ToRadians(15)), ToRadians(-5))
	})
	t.Run("a step reaching the target gives the target itself", func(t *testing.T) {
		assert.Equal(t, MoveTowardAngle(1, 4, 3), 4.0)
		assert.Equal(t, MoveTowardAngle(1, 4, 100), 4.0)
		assert.Equal(t, MoveTowardAngle(tenth, 4+2*Pi, 5), 4+2*Pi)
	})
	t.Run("a step of zero stays", func(t *testing.T) {
		assert.Equal(t, MoveTowardAngle(1, 4, 0), 1.0)
	})
	t.Run("half a turn apart turns by increasing angle", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveTowardAngle(0, Pi, 0.25), 0.25)
		geomtest.AssertNumber(t, MoveTowardAngle(Pi, 0, 0.25), Pi+0.25)
	})
	t.Run("a negative step turns away and is not clamped", func(t *testing.T) {
		geomtest.AssertNumber(t, MoveTowardAngle(0, Pi/2, -0.25), -0.25)
		geomtest.AssertNumber(t, MoveTowardAngle(Pi/2, 0, -0.25), Pi/2+0.25)
		geomtest.AssertNumber(t, MoveTowardAngle(0, Pi/2, -Pi), -Pi)
	})
	t.Run("an angle at the target stays for any step", func(t *testing.T) {
		assert.Equal(t, MoveTowardAngle(1, 1, -0.25), 1.0)
		assert.Equal(t, MoveTowardAngle(1, 1+2*Pi, -0.25), 1+2*Pi)
	})
	t.Run("a NaN gives NaN", func(t *testing.T) {
		assert.True(t, math.IsNaN(MoveTowardAngle(0, math.NaN(), 1)))
		assert.True(t, math.IsNaN(MoveTowardAngle(0, 1, math.NaN())))
	})
	t.Run("never passes the target and arrives in steps", func(t *testing.T) {
		for _, a := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
			for _, b := range []float64{0, 1, -2, Pi, 5, 2 * Pi, 100} {
				angle, distance := a, AngleDistance(a, b)
				for range 8 {
					angle = MoveTowardAngle(angle, b, 0.5)

					assert.True(t, AngleDistance(angle, b) <= max(distance-0.5, 0)+Delta, fmt.Sprintf("%v → %v: ", a, b))
					distance = AngleDistance(angle, b)
				}

				assert.Equal(t, angle, b, fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
}

func TestEqualAngle(t *testing.T) {
	t.Run("same angle", func(t *testing.T) {
		assert.True(t, EqualAngle(0.5, 0.5))
		assert.True(t, EqualAngle(0.5, 0.5+2*Pi))
		assert.True(t, EqualAngle(-Pi/2, 3*Pi/2))
	})
	t.Run("across the seam", func(t *testing.T) {
		assert.True(t, EqualAngle(0, -1e-9))
		assert.True(t, EqualAngle(0, 2*Pi-1e-9))
		assert.True(t, EqualAngle(-1e-9, 1e-9))
	})
	t.Run("different angle", func(t *testing.T) {
		assert.False(t, EqualAngle(0, 1e-3))
		assert.False(t, EqualAngle(0, Pi))
		assert.False(t, EqualAngle(0, math.NaN()))
	})
}

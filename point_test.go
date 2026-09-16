package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestPoint_Constructor(t *testing.T) {
	t.Run("from coordinates", func(t *testing.T) {
		AssertPoint(t, Pt(10, 16), Point[int]{X: 10, Y: 16})
		AssertPoint(t, Pt[float64](0.16, 204), Point[float64]{X: 0.16, Y: 204})
	})
	t.Run("zero", func(t *testing.T) {
		AssertPoint(t, ZeroPoint[int](), Point[int]{})
		AssertPoint(t, ZeroPoint[float64](), Point[float64]{})
	})
}

func TestPoint_Transform(t *testing.T) {
	t.Run("float64 matrix", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), Pt(9, 22))
		AssertPoint(t, Pt(0.6, -0.25).Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), Pt(3.385, 7.865))
	})
	t.Run("float32 matrix", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Transform(Mat[float32](1, 2, 3, 4, 5, 6)), Pt(8, 20))
		AssertPoint(t, Pt(0.6, -0.25).Transform(Mat[float32](1, 2, 3, 4, 5, 6)), Pt(3.1, 7.15))
	})
	t.Run("integer matrix converted to float", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Transform(Mat(1, 2, 3, 4, 5, 6).Float()), Pt(8, 20))
		AssertPoint(t, Pt(0.6, -0.25).Transform(Mat(1, 2, 3, 4, 5, 6).Float()), Pt(3.1, 7.15))
	})
}

func TestPoint_Add(t *testing.T) {
	t.Run("vector", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Add(Vec(3, -2)), Pt(4, 0))
		AssertPoint(t, Pt(0.6, -0.25).Add(Vec(100.1, -0.1)), Pt(100.7, -0.35))
	})
	t.Run("coordinates", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).AddXY(3, -2), Pt(4, 0))
		AssertPoint(t, Pt(0.6, -0.25).AddXY(100.1, -0.1), Pt(100.7, -0.35))
	})
}

func TestPoint_Subtract(t *testing.T) {
	AssertVector(t, Pt(1, 2).Subtract(Pt(3, -3)), Vec(-2, 5))
	AssertVector(t, Pt(0.6, -0.25).Subtract(Pt(100.1, -0.1)), Vec(-99.5, -0.15))
}

func TestPoint_Multiply(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Multiply(3), Pt(3, 6))
		AssertPoint(t, Pt(0.6, -0.25).Multiply(-1.5), Pt(-0.9, 0.375))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).MultiplyXY(3, 4), Pt(3, 8))
		AssertPoint(t, Pt(0.6, -0.25).MultiplyXY(-1.5, 2), Pt(-0.9, -0.5))
	})
}

func TestPoint_Divide(t *testing.T) {
	t.Run("uniform factor", func(t *testing.T) {
		AssertPoint(t, Pt(5, 10).Divide(2), Pt(3, 5)) // int: 2.5 rounds to 3
		AssertPoint(t, Pt(0.6, -0.25).Divide(-2), Pt(-0.3, 0.125))
	})
	t.Run("per-axis factor", func(t *testing.T) {
		AssertPoint(t, Pt(5, 10).DivideXY(3, 2), Pt(2, 5)) // int: 1.66 rounds to 2
		AssertPoint(t, Pt(0.6, -0.25).DivideXY(-4, 0.5), Pt(-0.15, -0.5))
	})
	t.Run("zero factor leaves the axis unchanged", func(t *testing.T) {
		AssertPoint(t, Pt(5, 10).Divide(0), Pt(5, 10))
		AssertPoint(t, Pt(0.6, -0.25).Divide(0), Pt(0.6, -0.25))

		// the guard is per axis, so the other one still divides
		AssertPoint(t, Pt(4, 8).DivideXY(0, 2), Pt(4, 4))
		AssertPoint(t, Pt(0.6, -0.25).DivideXY(2, 0), Pt(0.3, -0.25))
	})
}

func TestPoint_Abs(t *testing.T) {
	t.Run("negative coordinates", func(t *testing.T) {
		AssertPoint(t, Pt(-3, -4).Abs(), Pt(3, 4))
		AssertPoint(t, Pt(-1.5, 2.5).Abs(), Pt(1.5, 2.5))
	})
	t.Run("non-negative coordinates unchanged", func(t *testing.T) {
		AssertPoint(t, Pt(3, 4).Abs(), Pt(3, 4))
		AssertPoint(t, Pt(0, 0).Abs(), Pt(0, 0))
	})
}

func TestPoint_Round(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Pt(1.4, 2.5).Round(), Pt(1.0, 3.0))
		AssertPoint(t, Pt(-1.5, -2.4).Round(), Pt(-2.0, -2.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPoint(t, Pt(3, 4).Round(), Pt(3, 4))
	})
}

func TestPoint_Floor(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Pt(1.9, 2.1).Floor(), Pt(1.0, 2.0))
		AssertPoint(t, Pt(-1.1, -2.9).Floor(), Pt(-2.0, -3.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPoint(t, Pt(3, 4).Floor(), Pt(3, 4))
	})
}

func TestPoint_Ceil(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Pt(1.1, 2.9).Ceil(), Pt(2.0, 3.0))
		AssertPoint(t, Pt(-1.9, -2.1).Ceil(), Pt(-1.0, -2.0))
	})
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPoint(t, Pt(3, 4).Ceil(), Pt(3, 4))
	})
}

func TestPoint_DistanceTo(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		for _, test := range distanceFixtures {
			t.Run(test.name, func(t *testing.T) {
				AssertNumber(t, test.a.DistanceTo(test.b), test.euclidean)
			})
		}
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pt(1, 2).DistanceTo(Pt(2, 3)), Sqrt2)
	})
}

func TestPoint_DistanceSquaredTo(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		for _, test := range distanceFixtures {
			t.Run(test.name, func(t *testing.T) {
				AssertNumber(t, test.a.DistanceSquaredTo(test.b), test.squared)
			})
		}
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pt(1, 2).DistanceSquaredTo(Pt(2, 3)), 2)
	})
}

func TestPoint_ManhattanDistanceTo(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		for _, test := range distanceFixtures {
			t.Run(test.name, func(t *testing.T) {
				AssertNumber(t, test.a.ManhattanDistanceTo(test.b), test.manhattan)
			})
		}
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pt(1, 2).ManhattanDistanceTo(Pt(2, 3)), 2)
	})
}

func TestPoint_ChebyshevDistanceTo(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		for _, test := range distanceFixtures {
			t.Run(test.name, func(t *testing.T) {
				AssertNumber(t, test.a.ChebyshevDistanceTo(test.b), test.chebyshev)
			})
		}
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pt(1, 2).ChebyshevDistanceTo(Pt(2, 3)), 1)
	})
}

func TestPoint_OctileDistanceTo(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		for _, test := range distanceFixtures {
			t.Run(test.name, func(t *testing.T) {
				AssertNumber(t, test.a.OctileDistanceTo(test.b), test.octile)
			})
		}
	})
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Pt(1, 2).OctileDistanceTo(Pt(2, 3)), Sqrt2)
	})
}

// distanceFixtures feeds every DistanceTo variant, which all measure the same geometry.
var distanceFixtures = []struct {
	name                                             string
	a, b                                             Point[float64]
	euclidean, squared, manhattan, chebyshev, octile float64
}{
	{"diagonal step", Pt(1.0, 2.0), Pt(2.0, 3.0), Sqrt2, 2, 2, 1, Sqrt2},
	{"cardinal step", Pt(1.0, 2.0), Pt(4.0, 2.0), 3, 9, 3, 3, 3},
	{"knight move", Pt(0.0, 0.0), Pt(1.0, 2.0), math.Sqrt(5), 5, 3, 2, 1 + Sqrt2},
	{"fractional", Pt(0.6, -0.25), Pt(0.5, -0.35), math.Sqrt(0.02), 0.02, 0.2, 0.1, 0.1 * Sqrt2},
	{"negative direction", Pt(2.0, 3.0), Pt(1.0, 2.0), Sqrt2, 2, 2, 1, Sqrt2},
	{"same point", Pt(0.6, -0.25), Pt(0.6, -0.25), 0, 0, 0, 0, 0},
}

func TestPoint_Midpoint(t *testing.T) {
	AssertPoint(t, Pt(1, 2).Midpoint(Pt(3, -3)), Pt(2, -1)) // int: -0.5 rounds to -1
	AssertPoint(t, Pt(0.6, -0.25).Midpoint(Pt(100.1, -0.1)), Pt(50.35, -0.175))
}

func TestPoint_Lerp(t *testing.T) {
	t.Run("between points", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Lerp(Pt(3, -3), 0.3), Pt(2, 1))
		AssertPoint(t, Pt(0.6, -0.25).Lerp(Pt(100.1, -0.1), 0.1), Pt(10.55, -0.235))
	})
	t.Run("endpoints", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Lerp(Pt(3, -3), 0), Pt(1, 2))
		AssertPoint(t, Pt(1, 2).Lerp(Pt(3, -3), 1), Pt(3, -3))
	})
	t.Run("extrapolates outside the unit range", func(t *testing.T) {
		AssertPoint(t, Pt(0.0, 0.0).Lerp(Pt(2.0, 2.0), 2), Pt(4.0, 4.0))
		AssertPoint(t, Pt(0.0, 0.0).Lerp(Pt(2.0, 2.0), -1), Pt(-2.0, -2.0))
	})
}

func TestPoint_AngleTo(t *testing.T) {
	t.Run("cardinal directions", func(t *testing.T) {
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(1, 0)), ToRadians(0))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(0, 1)), ToRadians(90))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(-1, 0)), ToRadians(180))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(0, -1)), ToRadians(-90))
	})
	t.Run("diagonals", func(t *testing.T) {
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(1, 1)), ToRadians(45))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(-1, 1)), ToRadians(135))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(-1, -1)), ToRadians(-135))
		AssertNumber(t, Pt(0, 0).AngleTo(Pt(1, -1)), ToRadians(-45))
	})
	t.Run("off the origin", func(t *testing.T) {
		AssertNumber(t, Pt(2, 2).AngleTo(Pt(3, 2)), ToRadians(0))
		AssertNumber(t, Pt(0.6, -0.25).AngleTo(Pt(0.7, -0.35)), ToRadians(-45))
	})
}

func TestPoint_Equal(t *testing.T) {
	t.Run("same point", func(t *testing.T) {
		assert.True(t, Pt(1, 2).Equal(Pt(1, 2)))
		assert.True(t, Pt(0.6, -0.25).Equal(Pt(0.6, -0.25)))
	})
	t.Run("different point", func(t *testing.T) {
		assert.False(t, Pt(1, 2).Equal(Pt(3, -3)))
		assert.False(t, Pt(0.6, -0.25).Equal(Pt(100.1, -0.1)))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Pt(0.6, -0.25).Equal(Pt(0.6, -0.250001)))
	})
}

func TestPoint_Compare(t *testing.T) {
	t.Run("orders by x first", func(t *testing.T) {
		assert.Equal(t, Pt(1, 9).Compare(Pt(2, 0)), -1)
		assert.Equal(t, Pt(2, 0).Compare(Pt(1, 9)), 1)
	})
	t.Run("falls back to y", func(t *testing.T) {
		assert.Equal(t, Pt(1, 1).Compare(Pt(1, 2)), -1)
		assert.Equal(t, Pt(1, 2).Compare(Pt(1, 1)), 1)
	})
	t.Run("equal", func(t *testing.T) {
		assert.Equal(t, Pt(1, 2).Compare(Pt(1, 2)), 0)
		assert.Equal(t, Pt(0.6, -0.25).Compare(Pt(0.6, -0.25)), 0)
	})
	t.Run("negative coordinates", func(t *testing.T) {
		assert.Equal(t, Pt(-2, 0).Compare(Pt(-1, 0)), -1)
		assert.Equal(t, Pt(0, -2).Compare(Pt(0, -1)), -1)
	})
	t.Run("antisymmetric and transitive", func(t *testing.T) {
		points := []Point[int]{Pt(0, 0), Pt(0, 1), Pt(1, -1), Pt(1, 0), Pt(-1, 5)}
		for _, a := range points {
			for _, b := range points {
				assert.Equal(t, a.Compare(b), -b.Compare(a))

				for _, c := range points {
					if a.Compare(b) < 0 && b.Compare(c) < 0 {
						assert.True(t, a.Compare(c) < 0)
					}
				}
			}
		}
	})
	t.Run("sorts", func(t *testing.T) {
		points := []Point[int]{Pt(1, 2), Pt(-1, 0), Pt(1, -3), Pt(0, 7)}
		slices.SortFunc(points, Point[int].Compare)

		assert.Equal(t, points, []Point[int]{Pt(-1, 0), Pt(0, 7), Pt(1, -3), Pt(1, 2)})

		index, found := slices.BinarySearchFunc(points, Pt(1, -3), Point[int].Compare)
		assert.True(t, found)
		assert.Equal(t, index, 2)
	})
	t.Run("exact, unlike Equal", func(t *testing.T) {
		// within Delta, so Equal reports true while Compare still orders them
		a, b := Pt(0.0, 0.0), Pt(0.0, Delta/2)

		assert.True(t, a.Equal(b))
		assert.Equal(t, a.Compare(b), -1)
	})
}

func TestPoint_IsZero(t *testing.T) {
	t.Run("zero point", func(t *testing.T) {
		assert.True(t, Pt(0, 0).IsZero())
		assert.True(t, ZeroPoint[int]().IsZero())

		assert.True(t, Pt(0.0, 0.0).IsZero())
		assert.True(t, ZeroPoint[float64]().IsZero())
	})
	t.Run("negative zero", func(t *testing.T) {
		// an integer -0 is the same value as 0; a float one carries a sign bit
		assert.True(t, Pt(-0, 0).IsZero())
		assert.True(t, Pt(-0, -0).IsZero())

		assert.True(t, Pt(negativeZero, 0.0).IsZero())
		assert.True(t, Pt(negativeZero, negativeZero).IsZero())
	})
	t.Run("non-zero point", func(t *testing.T) {
		assert.False(t, Pt(1, 2).IsZero())
		assert.False(t, Pt(0.6, -0.25).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Pt(0.0, 0.000001).IsZero())
	})
}

func TestPoint_XY(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		x, y := Pt(10, 16).XY()
		AssertNumber(t, x, 10)
		AssertNumber(t, y, 16)
	})
	t.Run("float", func(t *testing.T) {
		x, y := Pt(0.6, -0.25).XY()
		AssertNumber(t, x, 0.6)
		AssertNumber(t, y, -0.25)
	})
}

func TestPoint_Vector(t *testing.T) {
	AssertVector(t, Pt(1, 2).Vector(), Vec(1, 2))
	AssertVector(t, Pt(0.6, -0.25).Vector(), Vec(0.6, -0.25))
}

func TestPoint_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Int(), Pt(1, 2))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertPoint(t, Pt(0.6, -0.25).Int(), Pt(1, 0))
		AssertPoint(t, Pt(-1.5, 2.5).Int(), Pt(-2, 3))
	})
}

func TestPoint_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertPoint(t, Pt(1, 2).Float(), Pt(1.0, 2.0))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertPoint(t, Pt(0.6, -0.25).Float(), Pt(0.6, -0.25))
	})
}

func TestPoint_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Pt(10, 16).String(), "(10,16)")
		assert.Equal(t, Pt(-4, 0).String(), "(-4,0)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Pt(100, -34.0000115).String(), "(100.00,-34.00)")
		assert.Equal(t, Pt(1.5, -0.25).String(), "(1.50,-0.25)")
	})
	t.Run("negative zero", func(t *testing.T) {
		assert.Equal(t, Pt(-0, 0).String(), "(0,0)")
		assert.Equal(t, Pt(negativeZero, 0.0).String(), "(-0.00,0.00)")
	})
}

func TestPoint_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Pt(10, 16), `{"x":10,"y":16}`)

		var p Point[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &p))
		AssertPoint(t, p, Pt(10, 16))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Pt(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)

		var p Point[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"x":100.0,"y":-34.0000115}`), &p))
		AssertPoint(t, p, Pt(100.0, -34.0000115))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, point := range pointFixtures {
			data, err := json.Marshal(point)
			assert.NoError(t, err)

			var decoded Point[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, point)
		}
	})
}

func TestPoint_Properties(t *testing.T) {
	t.Run("add and subtract are inverse", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				assert.True(t, a.Add(b.Subtract(a)).Equal(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("midpoint is lerp at one half", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				assert.True(t, a.Midpoint(b).Equal(a.Lerp(b, 0.5)), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("distance squared is the squared distance", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				distance := a.DistanceTo(b)
				AssertNumber(t, distance*distance, a.DistanceSquaredTo(b), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("chebyshev ≤ euclidean ≤ octile ≤ manhattan", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				chebyshev, euclidean := a.ChebyshevDistanceTo(b), a.DistanceTo(b)
				octile, manhattan := a.OctileDistanceTo(b), a.ManhattanDistanceTo(b)

				assert.True(t, chebyshev <= euclidean+Delta, fmt.Sprintf("%s → %s: chebyshev ≤ euclidean: ", a, b))
				assert.True(t, euclidean <= octile+Delta, fmt.Sprintf("%s → %s: euclidean ≤ octile: ", a, b))
				assert.True(t, octile <= manhattan+Delta, fmt.Sprintf("%s → %s: octile ≤ manhattan: ", a, b))
			}
		}
	})
	t.Run("opposite angles differ by pi", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				if a.Equal(b) {
					continue
				}

				AssertNumber(t, Abs(a.AngleTo(b)-b.AngleTo(a)), Pi, fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("compare zero implies equal", func(t *testing.T) {
		for _, a := range pointFixtures {
			for _, b := range pointFixtures {
				if a.Compare(b) == 0 {
					assert.True(t, a.Equal(b), fmt.Sprintf("%s → %s: ", a, b))
				}
			}
		}
	})
	t.Run("identity matrix leaves the point untouched", func(t *testing.T) {
		for _, point := range pointFixtures {
			assert.True(t, point.Transform(IdentityMatrix[float64]()).Equal(point), fmt.Sprintf("%s: ", point))
		}
	})
}

func TestPoint_Immutable(t *testing.T) {
	p1 := Pt(1, 2)
	p2 := Pt(3, -3)

	p1.Transform(IdentityMatrix[float64]())
	p1.Add(Vec(3, -2))
	p1.AddXY(2, 3)
	p1.Subtract(p2)
	p1.Multiply(2)
	p1.MultiplyXY(3, 4)
	p1.Divide(5)
	p1.DivideXY(10, 100)
	p1.Abs()
	p1.Round()
	p1.Floor()
	p1.Ceil()
	p1.Midpoint(p2)
	p1.Lerp(p2, 0.1)

	AssertPoint(t, p1, Pt(1, 2))
	AssertPoint(t, p2, Pt(3, -3))
}

// pointFixtures span the quadrants, the axes, and the diagonal.
var pointFixtures = []Point[float64]{
	Pt(0.0, 0.0),
	Pt(1.0, 2.0),
	Pt(-3.5, 0.25),
	Pt(12.5, -0.1),
	Pt(-1.0, -1.0),
	Pt(0.6, -0.25),
	Pt(7.0, 7.0),
	Pt(-0.5, 12.75),
}

func ExamplePt() {
	fmt.Println(Pt(3, 4))
	fmt.Println(Pt(1.5, -0.25))
	// Output:
	// (3,4)
	// (1.50,-0.25)
}

func ExamplePoint_Compare() {
	points := []Point[int]{Pt(1, 2), Pt(-1, 0), Pt(1, -3)}
	slices.SortFunc(points, Point[int].Compare)

	fmt.Println(points)
	// Output: [(-1,0) (1,-3) (1,2)]
}

func ExamplePoint_DistanceTo() {
	fmt.Printf("%.2f\n", Pt(0.0, 0.0).DistanceTo(Pt(3.0, 4.0)))
	// Output: 5.00
}

package geom

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestLine_Constructor(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, -1), Pt(2, 0)), Line[int]{Start: Pt(1, -1), End: Pt(2, 0)})
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.5, -1.25), Pt(2.5, 3.75)), Line[float64]{Start: Pt(0.5, -1.25), End: Pt(2.5, 3.75)})
	})
}

func TestLine_Translate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Translate(Vec(3, -2)), Ln(Pt(4, 0), Pt(6, 3)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Translate(Vec(100.1, -0.1)), Ln(Pt(100.7, -0.35), Pt(101.3, 3.3)))
	})
}

func TestLine_MoveTo(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).MoveTo(Pt(3, -2)), Ln(Pt(3, -2), Pt(5, 1)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).MoveTo(Pt(100.1, -0.1)), Ln(Pt(100.1, -0.1), Pt(100.7, 3.55)))
	})
}

func TestLine_Reverse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Reverse(), Ln(Pt(3, 5), Pt(1, 2)))
	})
	t.Run("float", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Reverse(), Ln(Pt(1.2, 3.4), Pt(0.6, -0.25)))
	})
}

func TestLine_Midpoint(t *testing.T) {
	t.Run("int rounds the half away from zero", func(t *testing.T) {
		AssertPoint(t, Ln(Pt(1, 2), Pt(3, 5)).Midpoint(), Pt(2, 4))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Midpoint(), Pt(0.9, 1.575))
	})
}

func TestLine_Vector(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVector(t, Ln(Pt(1, 2), Pt(3, 5)).Vector(), Vec(2, 3))
	})
	t.Run("float", func(t *testing.T) {
		AssertVector(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vector(), Vec(0.6, 3.65))
	})
}

func TestLine_Length(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(1, 2), Pt(3, 5)).Length(), math.Sqrt(13))
	})
	t.Run("float", func(t *testing.T) {
		AssertNumber(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Length(), math.Sqrt(13.6825))
	})
}

func TestLine_Vertices(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(1, 2), Pt(3, 5)).Vertices(), []Point[int]{{1, 2}, {3, 5}})
	})
	t.Run("float", func(t *testing.T) {
		AssertVertices(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Vertices(), []Point[float64]{{0.6, -0.25}, {1.2, 3.4}})
	})
}

func TestLine_Bounds(t *testing.T) {
	t.Run("spans the endpoints", func(t *testing.T) {
		l := Ln(Pt(1, 2), Pt(3, 5))

		AssertRect(t, l.Bounds(), Rect(Pt(2, 3), Sz(2, 3)))
		AssertPoint(t, l.Bounds().Min(), l.Start)
		AssertPoint(t, l.Bounds().Max(), l.End)
	})
	t.Run("float", func(t *testing.T) {
		AssertRect(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Bounds(), Rect(Pt(0.9, 1.575), Sz(0.6, 3.65)))
	})
}

func TestLine_Equal(t *testing.T) {
	t.Run("same line", func(t *testing.T) {
		assert.True(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 5))))
		assert.True(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.6, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("different line", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(1, 2), Pt(3, 4))))
		assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.5, -0.25), Pt(1.2, 3.4))))
	})
	t.Run("orientation matters", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).Equal(Ln(Pt(3, 5), Pt(1, 2))))
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Equal(Ln(Pt(0.6, -0.250001), Pt(1.2, 3.4))))
	})
}

func TestLine_IsZero(t *testing.T) {
	t.Run("zero line", func(t *testing.T) {
		assert.True(t, Line[int]{}.IsZero())
		assert.True(t, Ln(Pt(0, 0), Pt(0, 0)).IsZero())
		assert.True(t, Line[float64]{}.IsZero())
	})
	t.Run("one endpoint off the origin", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 0), Pt(0, 0)).IsZero())
		assert.False(t, Ln(Pt(0, 0), Pt(0, 1)).IsZero())
	})
	t.Run("non-zero line", func(t *testing.T) {
		assert.False(t, Ln(Pt(1, 2), Pt(3, 5)).IsZero())
		assert.False(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).IsZero())
	})
	t.Run("within delta", func(t *testing.T) {
		assert.True(t, Ln(Pt(0.0, 0.000001), Pt(0.0, 0.0)).IsZero())
	})
}

func TestLine_Int(t *testing.T) {
	t.Run("int is a no-op", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Int(), Ln(Pt(1, 2), Pt(3, 5)))
	})
	t.Run("float rounds", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Int(), Ln(Pt(1, 0), Pt(1, 3)))
	})
}

func TestLine_Float(t *testing.T) {
	t.Run("int widens", func(t *testing.T) {
		AssertLine(t, Ln(Pt(1, 2), Pt(3, 5)).Float(), Ln(Pt(1.0, 2.0), Pt(3.0, 5.0)))
	})
	t.Run("float is a no-op", func(t *testing.T) {
		AssertLine(t, Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)).Float(), Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)))
	})
}

func TestLine_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(10, 16), Pt(1, 2)).String(), "L((10,16);(1,2))")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)).String(), "L((100.00,-34.00);(0.20,0.40))")
	})
}

func TestLine_JSON(t *testing.T) {
	t.Run("int wire format", func(t *testing.T) {
		assert.JSON(t, Ln(Pt(10, 16), Pt(1, 2)), `{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`)

		var l Line[int]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10,"y":16},"b":{"x":1,"y":2}}`), &l))
		AssertLine(t, l, Ln(Pt(10, 16), Pt(1, 2)))
	})
	t.Run("float wire format", func(t *testing.T) {
		assert.JSON(t, Ln(Pt(100, -34.0000115), Pt(0.2, 0.4)), `{"a":{"x":100.0,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`)

		var l Line[float64]
		assert.NoError(t, json.Unmarshal([]byte(`{"a":{"x":10.1,"y":-34.0000115},"b":{"x":0.2,"y":0.4}}`), &l))
		AssertLine(t, l, Ln(Pt(10.1, -34.0000115), Pt(0.2, 0.4)))
	})
	t.Run("round-trip", func(t *testing.T) {
		for _, line := range lineFixtures {
			data, err := json.Marshal(line)
			assert.NoError(t, err)

			var decoded Line[float64]
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, line)
		}
	})
}

func TestLine_Properties(t *testing.T) {
	t.Run("reverse is its own inverse", func(t *testing.T) {
		for _, line := range lineFixtures {
			assert.True(t, line.Reverse().Reverse().Equal(line), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("reverse keeps the length and midpoint", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertNumber(t, line.Reverse().Length(), line.Length(), fmt.Sprintf("%s: ", line))
			assert.True(t, line.Reverse().Midpoint().Equal(line.Midpoint()), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("length is the vector length", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertNumber(t, line.Length(), line.Vector().Length(), fmt.Sprintf("%s: ", line))
			AssertNumber(t, line.Length(), line.Start.DistanceTo(line.End), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("midpoint is equidistant from both ends", func(t *testing.T) {
		for _, line := range lineFixtures {
			midpoint := line.Midpoint()

			AssertNumber(t, midpoint.DistanceTo(line.Start), midpoint.DistanceTo(line.End), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("translate keeps the vector", func(t *testing.T) {
		for _, line := range lineFixtures {
			for _, vector := range vectorFixtures {
				assert.True(t, line.Translate(vector).Vector().Equal(line.Vector()), fmt.Sprintf("%s → %s: ", line, vector))
			}
		}
	})
	t.Run("move to keeps the vector and starts where asked", func(t *testing.T) {
		for _, line := range lineFixtures {
			for _, point := range pointFixtures {
				moved := line.MoveTo(point)

				assert.True(t, moved.Start.Equal(point), fmt.Sprintf("%s → %s: ", line, point))
				assert.True(t, moved.Vector().Equal(line.Vector()), fmt.Sprintf("%s → %s: ", line, point))
			}
		}
	})
	t.Run("bounds span the endpoints", func(t *testing.T) {
		for _, line := range lineFixtures {
			bounds := line.Bounds()
			start, end := line.Start, line.End

			assert.True(t, bounds.Min().Equal(Pt(min(start.X, end.X), min(start.Y, end.Y))), fmt.Sprintf("%s: ", line))
			assert.True(t, bounds.Max().Equal(Pt(max(start.X, end.X), max(start.Y, end.Y))), fmt.Sprintf("%s: ", line))
		}
	})
	t.Run("vertices are the endpoints", func(t *testing.T) {
		for _, line := range lineFixtures {
			AssertVertices(t, line.Vertices(), []Point[float64]{line.Start, line.End})
		}
	})
}

func TestLine_Immutable(t *testing.T) {
	l := Ln(Pt(1, 2), Pt(3, 5))

	l.Translate(Vec(3, -2))
	l.MoveTo(Pt(4, 3))
	l.Reverse()

	AssertLine(t, l, Ln(Pt(1, 2), Pt(3, 5)))
}

// lineFixtures span axis-aligned, diagonal, degenerate, and backwards lines.
var lineFixtures = []Line[float64]{
	Ln(Pt(0.0, 0.0), Pt(0.0, 0.0)),
	Ln(Pt(1.0, 2.0), Pt(3.0, 5.0)),
	Ln(Pt(0.6, -0.25), Pt(1.2, 3.4)),
	Ln(Pt(-3.5, 0.25), Pt(-3.5, 7.0)),
	Ln(Pt(2.0, 2.0), Pt(-4.0, 2.0)),
	Ln(Pt(12.5, -0.1), Pt(-0.5, 12.75)),
}

func ExampleLn() {
	fmt.Println(Ln(Pt(1, 2), Pt(3, 5)))
	// Output: L((1,2);(3,5))
}

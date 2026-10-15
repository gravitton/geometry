package geom_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestPair_All(t *testing.T) {
	a, b := Pt(1, 2), Pt(3, 4)

	t.Run("yields the items it holds, in order", func(t *testing.T) {
		geomtest.AssertVertices(t, slices.Collect(Pair[Point[int]]{[2]Point[int]{a, b}, 2}.All()), []Point[int]{a, b})
		geomtest.AssertVertices(t, slices.Collect(Pair[Point[int]]{[2]Point[int]{a}, 1}.All()), []Point[int]{a})
	})
	t.Run("an empty pair yields nothing", func(t *testing.T) {
		assert.Nil(t, slices.Collect(Pair[Point[int]]{}.All()))
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		count := 0
		for range (Pair[Point[int]]{[2]Point[int]{a, b}, 2}).All() {
			count++

			break
		}

		assert.Equal(t, count, 1)
	})
	t.Run("the items are the slice Items[:Count]", func(t *testing.T) {
		crossings := Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5))

		geomtest.AssertVertices(t, crossings.Items[:crossings.Count], slices.Collect(crossings.All()))
	})
	t.Run("pairs of equal items compare equal", func(t *testing.T) {
		assert.True(t, Circ(Pt(0, 0), 5).Tangents(Pt(5, 0)) == Pair[Point[int]]{[2]Point[int]{Pt(5, 0)}, 1})
		assert.True(t, Circ(Pt(0, 0), 5).Tangents(Pt(1, 0)) == Pair[Point[int]]{})
	})
	t.Run("ranging a result allocates nothing", func(t *testing.T) {
		circle, other := Circ(Pt(0.0, 0.0), 5.0), Circ(Pt(6.0, 0.0), 5.0)

		assertAllocations(t, func() {
			for p := range circle.IntersectionCircle(other).All() {
				sinkBool = p.X > 0
			}
		}, 0)
	})
}

func TestPair_Equal(t *testing.T) {
	a, b := Pt(1, 2), Pt(3, 4)

	t.Run("the same items in the same order", func(t *testing.T) {
		assert.True(t, Pair[Point[int]]{[2]Point[int]{a, b}, 2}.Equal(Pair[Point[int]]{[2]Point[int]{a, b}, 2}))
		assert.True(t, Pair[Point[int]]{}.Equal(Pair[Point[int]]{}))
	})
	t.Run("another order, count or item differs", func(t *testing.T) {
		assert.False(t, Pair[Point[int]]{[2]Point[int]{a, b}, 2}.Equal(Pair[Point[int]]{[2]Point[int]{b, a}, 2}))
		assert.False(t, Pair[Point[int]]{[2]Point[int]{a}, 1}.Equal(Pair[Point[int]]{}))
		assert.False(t, Pair[Point[int]]{[2]Point[int]{a}, 1}.Equal(Pair[Point[int]]{[2]Point[int]{b}, 1}))
	})
	t.Run("items equal within the tolerance are equal, where == tells them apart", func(t *testing.T) {
		near, exact := Pair[Point[float64]]{[2]Point[float64]{Pt(1.0, 2.0+Delta/2)}, 1}, Pair[Point[float64]]{[2]Point[float64]{Pt(1.0, 2.0)}, 1}

		assert.True(t, near.Equal(exact))
		assert.False(t, near == exact)
	})
	t.Run("segments are compared by their own Equal", func(t *testing.T) {
		tangents := Circ(Pt(0.0, 0.0), 2.0).OuterTangents(Circ(Pt(10.0, 0.0), 2.0))

		assert.True(t, tangents.Equal(tangents))
		assert.False(t, tangents.Equal(Circ(Pt(0.0, 0.0), 2.0).OuterTangents(Circ(Pt(10.0, 0.0), 3.0))))
	})
}

func TestPair_IsZero(t *testing.T) {
	t.Run("a pair of no items", func(t *testing.T) {
		assert.True(t, Pair[Point[int]]{}.IsZero())
		assert.False(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)).IsZero())
	})
	t.Run("omitzero leaves an empty pair out", func(t *testing.T) {
		type result struct {
			Crossings Pair[Point[int]] `json:"c,omitzero"`
		}

		assert.JSON(t, result{}, `{}`)
		assert.JSON(t, result{Circ(Pt(0, 0), 5).Tangents(Pt(5, 0))}, `{"c":[{"x":5,"y":0}]}`)
	})
}

func TestPair_IsEmpty(t *testing.T) {
	assert.True(t, Pair[Point[int]]{}.IsEmpty())
	assert.True(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(20, 0), 5)).IsEmpty())
	assert.False(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)).IsEmpty())
}

func TestPair_AppendTo(t *testing.T) {
	a, b := Pt(1, 2), Pt(3, 4)

	t.Run("appends the items after the ones in dst", func(t *testing.T) {
		geomtest.AssertVertices(t, Pair[Point[int]]{[2]Point[int]{a, b}, 2}.AppendTo([]Point[int]{Pt(9, 9)}), []Point[int]{Pt(9, 9), a, b})
		geomtest.AssertVertices(t, Pair[Point[int]]{[2]Point[int]{a}, 1}.AppendTo(nil), []Point[int]{a})
	})
	t.Run("an empty pair leaves dst as it is", func(t *testing.T) {
		assert.Nil(t, Pair[Point[int]]{}.AppendTo(nil))
	})
	t.Run("a buffer with room allocates nothing", func(t *testing.T) {
		circle, other := Circ(Pt(0, 0), 5), Circ(Pt(6, 0), 5)
		buffer := make([]Point[int], 0, 2)

		assertAllocations(t, func() {
			sinkPoints = circle.IntersectionCircle(other).AppendTo(buffer[:0])
		}, 0)
	})
}

func TestPair_String(t *testing.T) {
	t.Run("points", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(0, 0), 5).IntersectionCircle(Circ(Pt(6, 0), 5)).String(), "Pair((3,4);(3,-4))")
	})
	t.Run("a single item and none", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(0, 0), 5).Tangents(Pt(5, 0)).String(), "Pair((5,0))")
		assert.Equal(t, Pair[Point[int]]{}.String(), "Pair()")
	})
	t.Run("segments", func(t *testing.T) {
		assert.Equal(t, Circ(Pt(0, 0), 2).OuterTangents(Circ(Pt(10, 0), 2)).String(), "Pair(Seg((0,2);(10,2));Seg((0,-2);(10,-2)))")
	})
}

func TestPair_JSON(t *testing.T) {
	a, b := Pt(1, 2), Pt(3, 4)

	t.Run("wire format is the array of the items it holds", func(t *testing.T) {
		assert.JSON(t, Pair[Point[int]]{[2]Point[int]{a, b}, 2}, `[{"x":1,"y":2},{"x":3,"y":4}]`)
		assert.JSON(t, Pair[Point[int]]{[2]Point[int]{a}, 1}, `[{"x":1,"y":2}]`)
		assert.JSON(t, Pair[Point[int]]{}, `[]`)
	})
	t.Run("decodes an array of at most two items", func(t *testing.T) {
		pair := Pair[Point[int]]{[2]Point[int]{b, a}, 2}

		assert.NoError(t, json.Unmarshal([]byte(`[{"x":1,"y":2}]`), &pair))
		assert.True(t, pair == Pair[Point[int]]{[2]Point[int]{a}, 1})

		assert.NoError(t, json.Unmarshal([]byte(`[]`), &pair))
		assert.True(t, pair.IsEmpty())
	})
	t.Run("a longer array or invalid input leaves the pair untouched", func(t *testing.T) {
		pair := Pair[Point[int]]{[2]Point[int]{a, b}, 2}

		assert.Error(t, json.Unmarshal([]byte(`[{"x":1,"y":2},{"x":3,"y":4},{"x":5,"y":6}]`), &pair))
		assert.Error(t, json.Unmarshal([]byte(`[{"x":"one"}]`), &pair))
		assert.True(t, pair == Pair[Point[int]]{[2]Point[int]{a, b}, 2})
	})
	t.Run("round-trip of segments", func(t *testing.T) {
		tangents := Circ(Pt(0.0, 0.0), 2.0).OuterTangents(Circ(Pt(10.0, 0.0), 2.0))

		bytes, err := json.Marshal(tangents)
		assert.NoError(t, err)

		var decoded Pair[Segment[float64]]
		assert.NoError(t, json.Unmarshal(bytes, &decoded))
		assert.True(t, decoded.Equal(tangents))
	})
}

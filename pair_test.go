package geom_test

import (
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

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			for p := range circle.IntersectionCircle(other).All() {
				sinkBool = p.X > 0
			}
		}), 0)
	})
}

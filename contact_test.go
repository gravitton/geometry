package geom_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestContact_Separation(t *testing.T) {
	t.Run("two circles part along the line of centers", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(1.5, 0.0), 1.0))
		separation := contact.Separation()

		assert.True(t, ok)
		assert.True(t, EqualDelta(separation.X, -0.5, 1e-5) && EqualDelta(separation.Y, 0, 1e-3), separation.String())
	})
	t.Run("two boxes part along the axis of the least overlap", func(t *testing.T) {
		contact, ok := Collide[float64](BoxFromMinMax(Pt(0.0, 0.0), Pt(2.0, 2.0)), BoxFromMinMax(Pt(1.5, 1.0), Pt(4.0, 4.0)))
		separation := contact.Separation()

		assert.True(t, ok)
		geomtest.AssertVector(t, separation, Vec(-0.5, 0.0))
	})
	t.Run("int", func(t *testing.T) {
		contact, ok := Collide[int](BoxFromMinMax(Pt(0, 0), Pt(4, 4)), BoxFromMinMax(Pt(3, -8), Pt(9, 8)))
		separation := contact.Separation()

		assert.True(t, ok)
		geomtest.AssertVector(t, separation, Vec(-1.0, 0.0))
	})
	t.Run("shapes apart have none", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(3.0, 0.0), 1.0))
		separation := contact.Separation()

		assert.False(t, ok)
		assert.True(t, separation.IsZero())
	})
	t.Run("touching shapes and shapes of no area are parted by nothing", func(t *testing.T) {
		for _, pair := range [][2]Shape[float64]{
			{BoxFromMinMax(Pt(0.0, 0.0), Pt(2.0, 2.0)), BoxFromMinMax(Pt(2.0, 0.0), Pt(4.0, 2.0))},
			{Seg(Pt(0.0, 0.0), Pt(4.0, 0.0)), Seg(Pt(2.0, 0.0), Pt(6.0, 0.0))},
			{Circ(Pt(1.0, 1.0), 0.0), Circ(Pt(1.0, 1.0), 0.0)},
		} {
			contact, ok := Collide(pair[0], pair[1])
			separation := contact.Separation()

			assert.True(t, ok, fmt.Sprintf("%v → %v: ", pair[0], pair[1]))
			assert.True(t, separation.LessOrEqual(Delta), fmt.Sprintf("%v → %v: %s", pair[0], pair[1], separation))
		}
	})
	t.Run("a shape in the notch of a polygon is parted from a side of the notch, not from its hull", func(t *testing.T) {
		polygon := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(6.0, 0.0), Pt(6.0, 6.0), Pt(4.0, 6.0), Pt(4.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 6.0), Pt(0.0, 6.0)})

		contact, ok := Collide[float64](Circ(Pt(3.0, 3.0), 1.5), polygon)
		separation := contact.Separation()

		assert.True(t, ok)
		assert.True(t, EqualDelta(separation.Length(), 0.5, 1e-5), separation.String())
	})
	t.Run("moved by it the shapes touch, and no longer intersect a little past it", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				message := fmt.Sprintf("%v → %v: ", a, b)

				contact, ok := Collide(a, b)
				separation := contact.Separation()
				assert.Equal(t, ok, Intersects[float64](a, b), message)
				if separation.LessOrEqual(1e-3) {
					continue
				}

				step := separation.Resize(1e-3)

				assert.False(t, Intersects[float64](translated(a, separation.Add(step)), b), message+separation.String())
				assert.True(t, Intersects[float64](translated(a, separation.Subtract(step)), b), message+separation.String())
				assert.False(t, Intersects[float64](a, translated(b, separation.Add(step).Negate())), message+separation.String())
			}
		}
	})
	t.Run("a NaN coordinate has none", func(t *testing.T) {
		_, ok := Collide[float64](Circ(Pt(math.NaN(), 0.0), 1.0), Circ(Pt(1.5, 0.0), 1.0))

		assert.False(t, ok)
	})
}

func TestContact_JSON(t *testing.T) {
	contact := Contact[int]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 1), Pt(4, 1)}, 2}}

	data, err := json.Marshal(contact)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"n":{"x":0,"y":-1},"d":1,"p":[{"x":0,"y":1},{"x":4,"y":1}]}`)

	var decoded Contact[int]
	assert.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, decoded, contact)
}

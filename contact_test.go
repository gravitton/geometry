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

func TestContact_Equal(t *testing.T) {
	contact := Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.0, 1.0), Pt(4.0, 1.0)}, 2}}

	t.Run("the same contact, within the tolerance", func(t *testing.T) {
		assert.True(t, contact.Equal(contact))
		assert.True(t, contact.Equal(Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 1 + Delta/2, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.0, 1.0+Delta/2), Pt(4.0, 1.0)}, 2}}))
	})
	t.Run("another normal, depth or points", func(t *testing.T) {
		assert.False(t, contact.Equal(Contact[float64]{Normal: Vec(0.0, 1.0), Depth: 1, Points: contact.Points}))
		assert.False(t, contact.Equal(Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 2, Points: contact.Points}))
		assert.False(t, contact.Equal(Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.0, 1.0)}, 1}}))
	})
	t.Run("the points in their places", func(t *testing.T) {
		assert.False(t, contact.Equal(Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(4.0, 1.0), Pt(0.0, 1.0)}, 2}}))
	})
}

func TestContact_IsZero(t *testing.T) {
	t.Run("two shapes apart", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(5.0, 0.0), 1.0))

		assert.False(t, ok)
		assert.True(t, contact.IsZero())
	})
	t.Run("a contact", func(t *testing.T) {
		assert.False(t, Contact[int]{Normal: Vec(0.0, -1.0), Depth: 1}.IsZero())
		assert.False(t, Contact[int]{Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 0)}, 1}}.IsZero())
	})
}

func TestContact_Cast(t *testing.T) {
	contact := Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 0.5, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.4, 1.0), Pt(3.6, 1.0)}, 2}}

	t.Run("rounds the points and keeps the normal and the depth", func(t *testing.T) {
		geomtest.AssertContact(t, Contact[float64].Cast[int32](contact), Contact[int32]{Normal: Vec(0.0, -1.0), Depth: 0.5, Points: Pair[Point[int32]]{[2]Point[int32]{Pt[int32](0, 1), Pt[int32](4, 1)}, 2}})
	})
	t.Run("two points rounding to the same one are held once", func(t *testing.T) {
		near := Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 0.5, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.4, 1.0), Pt(-0.4, 1.0)}, 2}}

		assert.Equal(t, Contact[float64].Cast[int32](near).Points, Pair[Point[int32]]{[2]Point[int32]{Pt[int32](0, 1)}, 1})
	})
}

func TestContact_Int(t *testing.T) {
	t.Run("rounds the points", func(t *testing.T) {
		contact := Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 0.5, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.4, 1.0), Pt(3.6, 1.0)}, 2}}

		geomtest.AssertContact(t, contact.Int(), Contact[int]{Normal: Vec(0.0, -1.0), Depth: 0.5, Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 1), Pt(4, 1)}, 2}})
	})
	t.Run("two points rounding to the same one are held once", func(t *testing.T) {
		contact := Contact[float64]{Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.4, 1.0), Pt(-0.4, 1.0)}, 2}}

		assert.Equal(t, contact.Int().Points, Pair[Point[int]]{[2]Point[int]{Pt(0, 1)}, 1})
	})
	t.Run("one point and none", func(t *testing.T) {
		assert.Equal(t, Contact[float64]{Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.4, 1.0)}, 1}}.Int().Points, Pair[Point[int]]{[2]Point[int]{Pt(0, 1)}, 1})
		assert.True(t, Contact[float64]{}.Int().Points.IsEmpty())
	})
}

func TestContact_Float(t *testing.T) {
	contact := Contact[int]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 1), Pt(4, 1)}, 2}}

	geomtest.AssertContact(t, contact.Float(), Contact[float64]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(0.0, 1.0), Pt(4.0, 1.0)}, 2}})
}

func TestContact_String(t *testing.T) {
	t.Run("two points", func(t *testing.T) {
		contact := Contact[int]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 1), Pt(4, 1)}, 2}}

		assert.Equal(t, contact.String(), "Contact(⟨0.00,-1.00⟩;1.00;(0,1);(4,1))")
	})
	t.Run("one point and none", func(t *testing.T) {
		assert.Equal(t, Contact[float64]{Normal: Vec(0.6, 0.8), Depth: 0.25, Points: Pair[Point[float64]]{[2]Point[float64]{Pt(1.5, 2.0)}, 1}}.String(), "Contact(⟨0.60,0.80⟩;0.25;(1.50,2.00))")
		assert.Equal(t, Contact[int]{}.String(), "Contact(⟨0.00,0.00⟩;0.00)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		contact := Contact[int]{Normal: Vec(0.0, -1.0), Depth: 1, Points: Pair[Point[int]]{[2]Point[int]{Pt(0, 1), Pt(4, 1)}, 2}}

		assertStringAllocation(t, contact.String)
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

package geom

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
)

func TestCircle_Constructor(t *testing.T) {
	AssertCircle(t, Circ(Pt(10, 16), 12), Circle[int]{Center: Pt(10, 16), Radius: 12})
	AssertCircle(t, Circ(Pt(0.16, 204), 5.1), Circle[float64]{Center: Pt(0.16, 204.0), Radius: 5.1})
}

func TestCircle_Translate(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Translate(Vec(3, -2)), Circ(Pt(4, 0), 10))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Translate(Vec(100.1, -0.1)), Circ(Pt(100.7, -0.35), 1.2))
}

func TestCircle_MoveTo(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).MoveTo(Pt(3, -2)), Circ(Pt(3, -2), 10))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).MoveTo(Pt(100.1, -0.1)), Circ(Pt(100.1, -0.1), 1.2))
}

func TestCircle_Scale(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Scale(2.5), Circ(Pt(1, 2), 25))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Scale(2.5), Circ(Pt(0.6, -0.25), 3.0))
}

func TestCircle_Resize(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Resize(8), Circ(Pt(1, 2), 8))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Resize(3.1), Circ(Pt(0.6, -0.25), 3.1))
}

func TestCircle_Grow(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Grow(8), Circ(Pt(1, 2), 18))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Grow(3.1), Circ(Pt(0.6, -0.25), 4.3))
}

func TestCircle_Shrink(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(8), Circ(Pt(1, 2), 2))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Shrink(0.3), Circ(Pt(0.6, -0.25), 0.9))

	// clamped to zero — never negative radius
	AssertCircle(t, Circ(Pt(1, 2), 10).Shrink(100), Circ(Pt(1, 2), 0))
}

func TestCircle_Area(t *testing.T) {
	assert.EqualDelta(t, Circ(Pt(1, 2), 10).Area(), Pi*100.0, Delta)
	assert.EqualDelta(t, Circ(Pt(0.6, -0.25), 1.2).Area(), Pi*1.44, Delta)
}

func TestCircle_Circumference(t *testing.T) {
	assert.EqualDelta(t, Circ(Pt(1, 2), 10).Circumference(), Pi*20.0, Delta)
	assert.EqualDelta(t, Circ(Pt(0.6, -0.25), 1.2).Circumference(), Pi*2.4, Delta)
}

func TestCircle_Diameter(t *testing.T) {
	assert.Equal(t, Circ(Pt(1, 2), 10).Diameter(), 20)
	assert.EqualDelta(t, Circ(Pt(0.6, -0.25), 1.2).Diameter(), 2.4, Delta)
}

func TestCircle_Bounds(t *testing.T) {
	AssertRect(t, Circ(Pt(1, 2), 10).Bounds(), Rect(Pt(1, 2), Sz(10, 10)))
	AssertRect(t, Circ(Pt(0.6, -0.25), 1.2).Bounds(), Rect(Pt(0.6, -0.25), Sz(1.2, 1.2)))
}

func TestCircle_Anchor(t *testing.T) {
	c := Circ(Pt(10.0, 10.0), 5.0)

	AssertPoint(t, c.Anchor(Right), Pt(15.0, 10.0))
	AssertPoint(t, c.Anchor(Left), Pt(5.0, 10.0))
	AssertPoint(t, c.Anchor(Top), Pt(10.0, 5.0))
	AssertPoint(t, c.Anchor(Bottom), Pt(10.0, 15.0))
	AssertPoint(t, c.Anchor(DirectionNone), Pt(10.0, 10.0))

	// diagonals land on the boundary too, unlike Rectangle's corners
	assert.EqualDelta(t, c.Center.DistanceTo(c.Anchor(DirectionUpRight)), c.Radius, Delta)
}

func TestCircle_Equal(t *testing.T) {
	assert.False(t, Circ(Pt(1, 2), 10).Equal(Circ(Pt(3, -3), 10)))
	assert.True(t, Circ(Pt(1, 2), 10).Equal(Circ(Pt(1, 2), 10)))

	assert.False(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(100.1, -0.1), 1.2)))
	assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(0.6, -0.25), 1.2)))
	assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Equal(Circ(Pt(0.6, -0.250001), 1.2)))
}

func TestCircle_IsZero(t *testing.T) {
	assert.True(t, Circle[int]{}.IsZero())
	assert.True(t, Circ(Pt(0, 0), 0).IsZero())
	assert.False(t, Circ(Pt(0, 0), 10).IsZero())
	assert.False(t, Circ(Pt(2, 1), 0).IsZero())
	assert.False(t, Circ(Pt(1, 2), 10).IsZero())

	assert.True(t, Circle[float64]{}.IsZero())
	assert.True(t, Circ(Pt(0.0, 0.000001), 0.0).IsZero())
	assert.False(t, Circ(Pt(0.0, 0.0), 10).IsZero())
	assert.False(t, Circ(Pt(2.0, 1.0), 0.0).IsZero())
	assert.False(t, Circ(Pt(1.0, 2.0), 10.0).IsZero())
}

func TestCircle_Contains(t *testing.T) {
	assert.False(t, Circ(Pt(1, 2), 10).Contains(Pt(1, 12)))
	assert.True(t, Circ(Pt(1, 2), 10).Contains(Pt(4, 4)))

	assert.False(t, Circ(Pt(0.6, -0.25), 1.2).Contains(Pt(0.0, 1.7)))
	assert.True(t, Circ(Pt(0.6, -0.25), 1.2).Contains(Pt(0.1, 0.8)))
}

func TestCircle_Int(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Int(), Circ(Pt(1, 2), 10))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Int(), Circ(Pt(1, 0), 1))
}

func TestCircle_Float(t *testing.T) {
	AssertCircle(t, Circ(Pt(1, 2), 10).Float(), Circ(Pt(1.0, 2.0), 10.0))
	AssertCircle(t, Circ(Pt(0.6, -0.25), 1.2).Float(), Circ(Pt(0.6, -0.25), 1.2))
}

func TestCircle_String(t *testing.T) {
	assert.Equal(t, Circ(Pt(10, 16), 5).String(), "C((10,16);5)")
	assert.Equal(t, Circ(Pt(100, -34.0000115), 0.2).String(), "C((100.00,-34.00);0.20)")
}

func TestCircle_Marshall(t *testing.T) {
	assert.JSON(t, Circ(Pt(10, 16), 12), `{"x":10,"y":16,"r":12}`)
	assert.JSON(t, Circ(Pt(100, -34.0000115), 0.2), `{"x":100.0,"y":-34.0000115,"r":0.2}`)
}

func TestCircle_Unmarshall(t *testing.T) {
	var p1 Circle[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16,"r":12}`), &p1))
	AssertCircle(t, p1, Circ(Pt(10, 16), 12))

	var p2 Circle[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115,"r":0.2}`), &p2))
	AssertCircle(t, p2, Circ(Pt(10.1, -34.0000115), 0.2))
}

func TestCircle_Immutable(t *testing.T) {
	c1 := Circ(Pt(1, 2), 10)

	c1.Translate(Vec(3, -2))
	c1.MoveTo(Pt(4, 3))
	c1.Scale(2)
	c1.Resize(15)
	c1.Grow(1)
	c1.Shrink(2)

	AssertCircle(t, c1, Circ(Pt(1, 2), 10))
}

package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

var (
	circleInt   = Circle[int]{Point[int]{1, 2}, 10}
	circleFloat = Circle[float64]{Point[float64]{0.6, -0.25}, 1.2}
)

func TestCircle_New(t *testing.T) {
	AssertCircle(t, Circ(Pt(10, 16), 12), 10, 16, 12)
	AssertCircle(t, Circ(Pt(0.16, 204), 5.1), 0.16, 204.0, 5.1)
}

func TestCircle_Translate(t *testing.T) {
	AssertCircle(t, circleInt.Translate(Vec(3, -2)), 4, 0, 10)
	AssertCircle(t, circleFloat.Translate(Vec(100.1, -0.1)), 100.7, -0.35, 1.2)
}

func TestCircle_MoveTo(t *testing.T) {
	AssertCircle(t, circleInt.MoveTo(Pt(3, -2)), 3, -2, 10)
	AssertCircle(t, circleFloat.MoveTo(Pt(100.1, -0.1)), 100.1, -0.1, 1.2)
}

func TestCircle_Scale(t *testing.T) {
	AssertCircle(t, circleInt.Scale(2.5), 1, 2, 25)
	AssertCircle(t, circleFloat.Scale(2.5), 0.6, -0.25, 3.0)
}

func TestCircle_Resize(t *testing.T) {
	AssertCircle(t, circleInt.Resize(8), 1, 2, 8)
	AssertCircle(t, circleFloat.Resize(3.1), 0.6, -0.25, 3.1)
}

func TestCircle_Grow(t *testing.T) {
	AssertCircle(t, circleInt.Grow(8), 1, 2, 18)
	AssertCircle(t, circleFloat.Grow(3.1), 0.6, -0.25, 4.3)
}

func TestCircle_Shrink(t *testing.T) {
	AssertCircle(t, circleInt.Shrink(8), 1, 2, 2)
	AssertCircle(t, circleFloat.Shrink(0.3), 0.6, -0.25, 0.9)
}

func TestCircle_Area(t *testing.T) {
	assert.EqualDelta(t, circleInt.Area(), math.Pi*100.0, Delta)
	assert.EqualDelta(t, circleFloat.Area(), math.Pi*1.44, Delta)
}

func TestCircle_Circumference(t *testing.T) {
	assert.EqualDelta(t, circleInt.Circumference(), math.Pi*20.0, Delta)
	assert.EqualDelta(t, circleFloat.Circumference(), math.Pi*2.4, Delta)
}

func TestCircle_Diameter(t *testing.T) {
	assert.Equal(t, circleInt.Diameter(), 20)
	assert.EqualDelta(t, circleFloat.Diameter(), 2.4, Delta)
}

func TestCircle_Bounds(t *testing.T) {
	AssertRect(t, circleInt.Bounds(), 1, 2, 10, 10)
	AssertRect(t, circleFloat.Bounds(), 0.6, -0.25, 1.2, 1.2)
}

func TestCircle_Equal(t *testing.T) {
	assert.False(t, circleInt.Equal(Circ(Pt(3, -3), 10)))
	assert.True(t, circleInt.Equal(circleInt))

	assert.False(t, circleFloat.Equal(Circ(Pt(100.1, -0.1), 1.2)))
	assert.True(t, circleFloat.Equal(Circ(Pt(0.6, -0.25), 1.2)))
	assert.True(t, circleFloat.Equal(Circ(Pt(0.6, -0.250001), 1.2)))
}

func TestCircle_IsZero(t *testing.T) {
	assert.True(t, Circle[int]{}.IsZero())
	assert.True(t, Circ(Pt(0, 0), 0).IsZero())
	assert.False(t, Circ(Pt(0, 0), 10).IsZero())
	assert.False(t, Circ(Pt(2, 1), 0).IsZero())
	assert.False(t, circleInt.IsZero())

	assert.True(t, Circle[float64]{}.IsZero())
	assert.True(t, Circ(Pt(0.0, 0.000001), 0.0).IsZero())
	assert.False(t, Circ(Pt(0.0, 0.0), 10).IsZero())
	assert.False(t, Circ(Pt(2.0, 1.0), 0.0).IsZero())
	assert.False(t, Circ(Pt(1.0, 2.0), 10.0).IsZero())
}

func TestCircle_Contains(t *testing.T) {
	assert.False(t, circleInt.Contains(Pt(1, 12)))
	assert.True(t, circleInt.Contains(Pt(4, 4)))

	assert.False(t, circleFloat.Contains(Pt(0.0, 1.7)))
	assert.True(t, circleFloat.Contains(Pt(0.1, 0.8)))
}

func TestCircle_Int(t *testing.T) {
	AssertCircle(t, circleInt.Int(), 1, 2, 10)
	AssertCircle(t, circleFloat.Int(), 1, 0, 1)
}

func TestCircle_Float(t *testing.T) {
	AssertCircle(t, circleInt.Float(), 1.0, 2.0, 10.0)
	AssertCircle(t, circleFloat.Float(), 0.6, -0.25, 1.2)
}

func TestCircle_String(t *testing.T) {
	assert.Equal(t, Circ(Pt(10, 16), 5).String(), "C((10,16);5)")
	assert.Equal(t, Circ(Pt(100, -34.0000115), 0.2).String(), "C((100,-34.00);0.20)")
}

func TestCircle_Marshall(t *testing.T) {
	assert.JSON(t, Circ(Pt(10, 16), 12), `{"x":10,"y":16,"r":12}`)
	assert.JSON(t, Circ(Pt(100, -34.0000115), 0.2), `{"x":100.0,"y":-34.0000115,"r":0.2}`)
}

func TestCircle_Unmarshall(t *testing.T) {
	var p1 Circle[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16,"r":12}`), &p1))
	AssertCircle(t, p1, 10, 16, 12)

	var p2 Circle[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115,"r":0.2}`), &p2))
	AssertCircle(t, p2, 10.1, -34.0000115, 0.2)
}

func TestCircle_Immutable(t *testing.T) {
	c1 := circleInt

	c1.Translate(Vec(3, -2))
	c1.MoveTo(Pt(4, 3))
	c1.Scale(2)
	c1.Resize(15)
	c1.Grow(1)
	c1.Shrink(2)

	AssertCircle(t, c1, 1, 2, 10)
}

func AssertCircle[T Number](t *testing.T, c Circle[T], x, y, radius T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, c.Center, x, y, append(messages, "Center.")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(c.Radius), float64(radius), Delta, append(messages, "Radius: ")...) {
		ok = false
	}

	return ok
}

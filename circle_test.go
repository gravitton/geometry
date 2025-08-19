package geom

import (
	"encoding/json"
	"github.com/gravitton/assert"
	"testing"
)

func TestCircle_New(t *testing.T) {
	testCircle(t, C(10, 16, 12), 10, 16, 12)
	testCircle(t, C[float64](0.16, 204, 5.1), 0.16, 204.0, 5.1)
}

func TestCircle_Translate(t *testing.T) {
	testCircle(t, C(1, 2, 10).Translate(Vec(3, -2)), 4, 0, 10)
	testCircle(t, C(0.4, -0.25, 1.2).Translate(Vec(100.1, -0.1)), 100.5, -0.35, 1.2)
}

func TestCircle_MoveTo(t *testing.T) {
	testCircle(t, C(1, 2, 10).MoveTo(Pt(3, -2)), 3, -2, 10)
	testCircle(t, C(0.4, -0.25, 1.2).MoveTo(Pt(100.1, -0.1)), 100.1, -0.1, 1.2)
}

func TestCircle_Scale(t *testing.T) {
	testCircle(t, C(1, 2, 10).Scale(2.5), 1, 2, 25)
	testCircle(t, C(0.4, -0.25, 1.2).Scale(2.5), 0.4, -0.25, 3.0)
}

func TestCircle_Resize(t *testing.T) {
	testCircle(t, C(1, 2, 10).Resize(8), 1, 2, 8)
	testCircle(t, C(0.4, -0.25, 1.2).Resize(3.1), 0.4, -0.25, 3.1)
}

func TestCircle_Equal(t *testing.T) {
	assert.False(t, C(1, 2, 10).Equal(C(3, -3, 10)))
	assert.True(t, C(1, 2, 10).Equal(C(1, 2, 10)))

	assert.False(t, C(0.4, -0.25, 1.2).Equal(C(100.1, -0.1, 1.2)))
	assert.True(t, C(0.4, -0.25, 1.2).Equal(C(0.4, -0.25, 1.2)))
	assert.True(t, C(0.4, -0.25, 1.2).Equal(C(0.4, -0.250001, 1.2)))
}

func TestCircle_Contains(t *testing.T) {
	assert.False(t, C(1, 2, 10).Contains(Pt(1, 12)))
	assert.True(t, C(1, 2, 10).Contains(Pt(4, 4)))

	assert.False(t, C(0.4, -0.25, 1.2).Contains(Pt(0.0, 1.7)))
	assert.True(t, C(0.4, -0.25, 1.2).Contains(Pt(0.1, 0.8)))
}

func TestCircle_String(t *testing.T) {
	assert.Equal(t, C(10, 16, 5).String(), "C((+10,+16);+5)")
	assert.Equal(t, C(100, -34.0000115, 0.2).String(), "C((+100,-34.00);+0.20)")
}

func TestCircle_Marshall(t *testing.T) {
	assert.JSON(t, C(10, 16, 12), `{"c":{"x":10,"y":16},"r":12}`)
	assert.JSON(t, C(100, -34.0000115, 0.2), `{"c":{"x":100.0,"y":-34.0000115},"r":0.2}`)
}

func TestCircle_Unmarshall(t *testing.T) {
	var p1 Circle[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"c":{"x":10,"y":16},"r":12}`), &p1))
	testCircle(t, p1, 10, 16, 12)

	var p2 Circle[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"c":{"x":10.1,"y":-34.0000115},"r":0.2}`), &p2))
	testCircle(t, p2, 10.1, -34.0000115, 0.2)
}

func TestCircle_Immutable(t *testing.T) {
	c1 := C(1, 2, 10)

	c1.Translate(Vec(3, -2))

	testCircle(t, c1, 1, 2, 10)
}

func testCircle[T Number](t *testing.T, c Circle[T], x, y, radius T) {
	t.Helper()

	testPoint(t, c.Center, x, y)
	assert.EqualDelta(t, float64(c.Radius), float64(radius), Delta)
}

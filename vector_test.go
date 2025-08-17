package geom

import (
	"encoding/json"
	"github.com/gravitton/assert"
	"testing"
)

func TestVector_New(t *testing.T) {
	testVector(t, Vec(10, 16), 10, 16)
	testVector(t, Vec[float64](0.16, 204), 0.16, 204.0)

	testVector(t, ZeroVector[int](), 0, 0)
	testVector(t, IdentityVector[float64](), 1, 1)
	testVector(t, UpVector[float64](), 0, 1)
	testVector(t, DownVector[float64](), 0, -1)
	testVector(t, RightVector[float64](), 1, 0)
	testVector(t, LeftVector[float64](), -1, 0)
}

func TestVector_Add(t *testing.T) {
	testVector(t, Vec(1, 2).Add(Vec(3, -2)), 4, 0)
	testVector(t, Vec(0.4, -0.25).Add(Vec(100.1, -0.1)), 100.5, -0.35)
}

func TestVector_Sub(t *testing.T) {
	testVector(t, Vec(1, 2).Sub(Vec(3, -3)), -2, 5)
	testVector(t, Vec(0.4, -0.25).Sub(Vec(100.1, -0.1)), -99.7, -0.15)
}

func TestVector_Equal(t *testing.T) {
	assert.False(t, Vec(1, 2).Equal(Vec(3, -3)))
	assert.True(t, Vec(1, 2).Equal(Vec(1, 2)))

	assert.False(t, Vec(0.4, -0.25).Equal(Vec(100.1, -0.1)))
	assert.True(t, Vec(0.4, -0.25).Equal(Vec(0.4, -0.25)))
	assert.True(t, Vec(0.4, -0.25).Equal(Vec(0.4, -0.250001)))
}

func TestVector_Zero(t *testing.T) {
	assert.False(t, Vec(1, 2).Zero())
	assert.True(t, Vec(0, -0).Zero())
	assert.True(t, ZeroVector[int]().Zero())

	assert.False(t, Vec(0.4, -0.25).Zero())
	assert.True(t, Vec(0.0, -0.0).Zero())
	assert.True(t, Vec(0.0, 0.000001).Zero())
	assert.True(t, ZeroVector[float64]().Zero())
}

func TestVector_String(t *testing.T) {
	assert.Equal(t, Vec(10, 16).String(), "⟨+10,+16⟩")
	assert.Equal(t, Vec(100, -34.0000115).String(), "⟨+100,-34.00⟩")
}

func TestVector_Marshall(t *testing.T) {
	assert.JSON(t, Vec(10, 16), `{"x":10,"y":16}`)
	assert.JSON(t, Vec(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)
}

func TestVector_Unmarshall(t *testing.T) {
	var p1 Vector[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &p1))
	testVector(t, p1, 10, 16)

	var p2 Vector[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &p2))
	testVector(t, p2, 10.1, -34.0000115)
}

func TestVector_Immutable(t *testing.T) {
	p1 := Vec(1, 2)
	p2 := Vec(3, -3)

	p1.Add(Vec(3, -2))
	p1.Sub(p2)

	testVector(t, p1, 1, 2)
	testVector(t, p2, 3, -3)
}

func testVector[T Number](t *testing.T, p Vector[T], x, y T) {
	t.Helper()

	assert.EqualDelta(t, float64(p.X), float64(x), Delta)
	assert.EqualDelta(t, float64(p.Y), float64(y), Delta)
}

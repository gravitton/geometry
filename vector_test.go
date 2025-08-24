package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestVector_New(t *testing.T) {
	testVector(t, V(10, 16), 10, 16)
	testVector(t, V[float64](0.16, 204), 0.16, 204.0)

	testVector(t, ZeroVector[int](), 0, 0)
	testVector(t, IdentityVector[float64](), 1, 1)
	testVector(t, UpVector[float64](), 0, 1)
	testVector(t, DownVector[float64](), 0, -1)
	testVector(t, RightVector[float64](), 1, 0)
	testVector(t, LeftVector[float64](), -1, 0)
}

func TestVector_Add(t *testing.T) {
	testVector(t, V(1, 2).Add(V(3, -2)), 4, 0)
	testVector(t, V(1, 2).AddXY(3, -2), 4, 0)
	testVector(t, V(0.4, -0.25).Add(V(100.1, -0.1)), 100.5, -0.35)
	testVector(t, V(0.4, -0.25).AddXY(100.1, -0.1), 100.5, -0.35)
}

func TestVector_Sub(t *testing.T) {
	testVector(t, V(1, 2).Subtract(V(3, -3)), -2, 5)
	testVector(t, V(1, 2).SubtractXY(3, -3), -2, 5)
	testVector(t, V(0.4, -0.25).Subtract(V(100.1, -0.1)), -99.7, -0.15)
	testVector(t, V(0.4, -0.25).SubtractXY(100.1, -0.1), -99.7, -0.15)
}

func TestVector_Multiply(t *testing.T) {
	testVector(t, V(1, 2).Multiply(3), 3, 6)
	testVector(t, V(1, 2).MultiplyXY(3, 2), 3, 4)
	testVector(t, V(0.4, -0.25).Multiply(-1.5), -0.6, 0.375)
	testVector(t, V(0.4, -0.25).MultiplyXY(-1.5, 2), -0.6, -0.5)
}

func TestVector_Divide(t *testing.T) {
	testVector(t, V(5, 10).Divide(2), 3, 5)
	testVector(t, V(5, 10).DivideXY(3, 2), 2, 5)
	testVector(t, V(0.4, -0.25).Divide(-2), -0.2, 0.125)
	testVector(t, V(0.4, -0.25).DivideXY(-4, 0.5), -0.1, -0.5)
}

func TestVector_Negate(t *testing.T) {
	testVector(t, V(1, 2).Negate(), -1, -2)
	testVector(t, V(0.4, -0.25).Negate(), -0.4, 0.25)
}

func TestVector_Rotate(t *testing.T) {
	testVector(t, V(1, 0).Rotate(ToRadians(90)), 0, 1)
	testVector(t, V(0.4, -0.25).Rotate(math.Pi), -0.4, 0.25)
}

func TestVector_Resize(t *testing.T) {
	testVector(t, V(1, 2).Resize(5), 2, 4)
	testVector(t, V(0.4, -0.25).Resize(5), 4.239991, -2.649994)
}

func TestVector_Normalize(t *testing.T) {
	testVector(t, V(1, 2).Normalize(), 0, 1)
	testVector(t, V(0, 0).Normalize(), 1, 0)
	testVector(t, V(0.4, -0.25).Normalize(), 0.847998, -0.529998)
}

func TestVector_Abs(t *testing.T) {
	testVector(t, V(1, 2).Abs(), 1, 2)
	testVector(t, V(-1, 3).Abs(), 1, 3)
	testVector(t, V(-0.4, -0.25).Abs(), 0.4, 0.25)
}

func TestVector_Dot(t *testing.T) {
	assert.Equal(t, V(1, 2).Dot(V(3, -3)), -3)
	assert.EqualDelta(t, V(0.4, -0.25).Dot(V(100.1, -0.1)), 40.065, Delta)
}

func TestVector_Cross(t *testing.T) {
	assert.Equal(t, V(1, 2).Cross(V(3, -3)), -9)
	assert.EqualDelta(t, V(0.4, -0.25).Cross(V(100.1, -0.1)), 24.985, Delta)
}

func TestVector_Normal(t *testing.T) {
	testVector(t, V(5, 10).Normal(), -10, 5)
	testVector(t, V(0.4, -0.25).Normal(), 0.25, 0.4)
}

func TestVector_Length(t *testing.T) {
	assert.EqualDelta(t, V(1, 2).Length(), math.Sqrt(5), Delta)
	assert.EqualDelta(t, V(0.4, -0.25).Length(), 0.4716990566028302, Delta)
}

func TestVector_LengthSquared(t *testing.T) {
	assert.Equal(t, V(1, 2).LengthSquared(), 5)
	assert.EqualDelta(t, V(0.4, -0.25).LengthSquared(), 0.2225, Delta)
}

func TestVector_Angle(t *testing.T) {
	assert.EqualDelta(t, V(2, 2).Angle(), ToRadians(45), Delta)
	assert.EqualDelta(t, V(0.4, -0.25).Angle(), ToRadians(-32.005383), Delta)
}

func TestVector_Equal(t *testing.T) {
	assert.False(t, V(1, 2).Equal(V(3, -3)))
	assert.True(t, V(1, 2).Equal(V(1, 2)))

	assert.False(t, V(0.4, -0.25).Equal(V(100.1, -0.1)))
	assert.True(t, V(0.4, -0.25).Equal(V(0.4, -0.25)))
	assert.True(t, V(0.4, -0.25).Equal(V(0.4, -0.250001)))
}

func TestVector_Zero(t *testing.T) {
	assert.False(t, V(1, 2).Zero())
	assert.True(t, V(0, 0).Zero())
	assert.True(t, ZeroVector[int]().Zero())

	assert.False(t, V(0.4, -0.25).Zero())
	assert.True(t, V(0.0, 0.0).Zero())
	assert.True(t, V(0.0, 0.000001).Zero())
	assert.True(t, ZeroVector[float64]().Zero())
}

func TestVector_Unit(t *testing.T) {
	assert.False(t, V(1, 2).Unit())
	assert.True(t, V(1, 2).Normalize().Unit())
	assert.True(t, V(1.1, 2.1).Normalize().Unit())
	assert.False(t, V(0, 0).Unit())
	assert.True(t, V(1, 0).Unit())
	assert.True(t, V(1/math.Sqrt2, 1/math.Sqrt2).Unit())

	assert.True(t, UpVector[float64]().Unit())
	assert.True(t, DownVector[float64]().Unit())
	assert.True(t, LeftVector[float64]().Unit())
	assert.True(t, RightVector[float64]().Unit())
}

func TestVector_Less(t *testing.T) {
	assert.False(t, V(1, 2).Less(2))
	assert.True(t, V(1, 2).Less(3))

	assert.False(t, V(0.4, -0.25).Less(0.1))
	assert.True(t, V(0.4, -0.25).Less(0.49))
}

func TestVector_String(t *testing.T) {
	assert.Equal(t, V(10, 16).String(), "⟨10,16⟩")
	assert.Equal(t, V(100, -34.0000115).String(), "⟨100,-34.00⟩")
}

func TestVector_Marshall(t *testing.T) {
	assert.JSON(t, V(10, 16), `{"x":10,"y":16}`)
	assert.JSON(t, V(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)
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
	p1 := V(1, 2)
	p2 := V(3, -3)

	p1.Add(V(3, -2))
	p1.Subtract(p2)
	p1.Multiply(2)
	p1.MultiplyXY(3, 4)

	testVector(t, p1, 1, 2)
	testVector(t, p2, 3, -3)
}

func testVector[T Number](t *testing.T, p Vector[T], x, y T) {
	t.Helper()

	assert.EqualDelta(t, float64(p.X), float64(x), Delta)
	assert.EqualDelta(t, float64(p.Y), float64(y), Delta)
}

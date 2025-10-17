package geom

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gravitton/assert"
)

var (
	vectorInt   = Vector[int]{10, 16}
	vectorFloat = Vector[float64]{0.6, -0.25}
)

func TestVector_New(t *testing.T) {
	assertVector(t, vectorInt, 10, 16)
	assertVector(t, Vec[float64](0.16, 204), 0.16, 204.0)

	assertVector(t, ZeroVector[int](), 0, 0)
	assertVector(t, OneVector[float64](), 1, 1)
	assertVector(t, UpVector[float64](), 0, -1)
	assertVector(t, DownVector[float64](), 0, 1)
	assertVector(t, RightVector[float64](), 1, 0)
	assertVector(t, LeftVector[float64](), -1, 0)
}

func TestVector_Transform(t *testing.T) {
	assertVector(t, vectorInt.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 48, 132)
	assertVector(t, vectorFloat.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 0.085, 1.265)
}

func TestVector_Add(t *testing.T) {
	assertVector(t, vectorInt.Add(Vec(3, -2)), 13, 14)
	assertVector(t, vectorInt.AddXY(3, -2), 13, 14)
	assertVector(t, vectorFloat.Add(Vec(100.1, -0.1)), 100.7, -0.35)
	assertVector(t, vectorFloat.AddXY(100.1, -0.1), 100.7, -0.35)
}

func TestVector_Sub(t *testing.T) {
	assertVector(t, vectorInt.Subtract(Vec(3, -3)), 7, 19)
	assertVector(t, vectorInt.SubtractXY(3, -3), 7, 19)
	assertVector(t, vectorFloat.Subtract(Vec(100.1, -0.1)), -99.5, -0.15)
	assertVector(t, vectorFloat.SubtractXY(100.1, -0.1), -99.5, -0.15)
}

func TestVector_Multiply(t *testing.T) {
	assertVector(t, vectorInt.Multiply(3), 30, 48)
	assertVector(t, vectorInt.MultiplyXY(3, 2), 30, 32)
	assertVector(t, vectorFloat.Multiply(-1.5), -0.9, 0.375)
	assertVector(t, vectorFloat.MultiplyXY(-1.5, 2), -0.9, -0.5)
}

func TestVector_Divide(t *testing.T) {
	assertVector(t, vectorInt.Divide(2), 5, 8)
	assertVector(t, vectorInt.DivideXY(3, 2), 3, 8)
	assertVector(t, vectorFloat.Divide(-2), -0.3, 0.125)
	assertVector(t, vectorFloat.DivideXY(-4, 0.5), -0.15, -0.5)
}

func TestVector_Negate(t *testing.T) {
	assertVector(t, vectorInt.Negate(), -10, -16)
	assertVector(t, vectorFloat.Negate(), -0.6, 0.25)
}

func TestVector_Rotate(t *testing.T) {
	assertVector(t, Vec(1, 0).Rotate(ToRadians(90)), 0, 1)
	assertVector(t, vectorFloat.Rotate(math.Pi), -0.6, 0.25)
}

func TestVector_Resize(t *testing.T) {
	assertVector(t, vectorInt.Resize(5), 3, 4)
	assertVector(t, vectorFloat.Resize(5), 4.615384, -1.923076)
}

func TestVector_Normalize(t *testing.T) {
	assertVector(t, vectorInt.Normalize(), 0, 1)
	assert.Equal(t, vectorInt.Normalize().Length(), 1)
	assertVector(t, vectorInt.Negate().Normalize(), 1, 0)
	assert.Equal(t, vectorInt.Negate().Normalize().Length(), 1)
	assertVector(t, Vec(0, 0).Normalize(), 1, 0)
	assertVector(t, vectorFloat.Normalize(), 0.923076, -0.384615)
	assert.Equal(t, vectorFloat.Normalize().Length(), 1)
}

func TestVector_Abs(t *testing.T) {
	assertVector(t, vectorInt.Abs(), 10, 16)
	assertVector(t, Vec(-1, -3).Abs(), 1, 3)
	assertVector(t, vectorFloat.Abs(), 0.6, 0.25)
}

func TestVector_Dot(t *testing.T) {
	assert.Equal(t, vectorInt.Dot(Vec(3, -3)), -18)
	assert.EqualDelta(t, vectorFloat.Dot(Vec(100.1, -0.1)), 60.085, Delta)
}

func TestVector_Cross(t *testing.T) {
	assert.Equal(t, vectorInt.Cross(Vec(3, -3)), -78)
	assert.EqualDelta(t, vectorFloat.Cross(Vec(100.1, -0.1)), 24.965, Delta)
}

func TestVector_Normal(t *testing.T) {
	assertVector(t, vectorInt.Normal(), -16, 10)
	assertVector(t, vectorFloat.Normal(), 0.25, 0.6)
}

func TestVector_Length(t *testing.T) {
	assert.EqualDelta(t, vectorInt.Length(), math.Sqrt(356), Delta)
	assert.EqualDelta(t, vectorFloat.Length(), 0.65, Delta)
}

func TestVector_LengthSquared(t *testing.T) {
	assert.Equal(t, vectorInt.LengthSquared(), 356)
	assert.EqualDelta(t, vectorFloat.LengthSquared(), 0.4225, Delta)
}

func TestVector_Angle(t *testing.T) {
	assert.EqualDelta(t, Vec(2, 2).Angle(), ToRadians(45), Delta)
	assert.EqualDelta(t, vectorFloat.Angle(), -0.39479111, Delta)
}

func TestVector_Lerp(t *testing.T) {
	assertVector(t, vectorInt.Lerp(Vec(20, 20), 0.25), 13, 17)
	assertVector(t, vectorFloat.Lerp(Vec(-10.0, 10.0), 0.25), -2.05, 2.3125)
}

func TestVector_Equal(t *testing.T) {
	assert.True(t, vectorInt.Equal(vectorInt))
	assert.False(t, vectorInt.Equal(Vec(3, -3)))

	assert.True(t, vectorFloat.Equal(vectorFloat))
	assert.False(t, vectorFloat.Equal(Vec(100.1, -0.1)))
	assert.True(t, vectorFloat.Equal(Vec(0.6, -0.250001)))
}

func TestVector_IsZero(t *testing.T) {
	assert.False(t, vectorInt.IsZero())
	assert.True(t, Vec(0, 0).IsZero())
	assert.True(t, ZeroVector[int]().IsZero())

	assert.False(t, vectorFloat.IsZero())
	assert.True(t, Vec(0.0, 0.0).IsZero())
	assert.True(t, Vec(0.0, 0.000001).IsZero())
	assert.True(t, ZeroVector[float64]().IsZero())
}

func TestVector_IsOne(t *testing.T) {
	assert.False(t, vectorInt.IsOne())
	assert.True(t, OneVector[int]().IsOne())

	assert.False(t, vectorFloat.IsOne())
	assert.True(t, Vec(1.0, 1.000001).IsOne())
	assert.True(t, OneVector[float64]().IsOne())
}

func TestVector_IsUp(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsUp())
	assert.False(t, OneVector[int]().IsUp())
	assert.True(t, UpVector[int]().IsUp())
	assert.False(t, DownVector[int]().IsUp())
	assert.False(t, LeftVector[int]().IsUp())
	assert.False(t, RightVector[int]().IsUp())
	assert.True(t, UpLeftVector().IsUp())
	assert.True(t, UpRightVector().IsUp())
	assert.False(t, DownLeftVector().IsUp())
	assert.False(t, DownRightVector().IsUp())
}

func TestVector_IsDown(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsDown())
	assert.True(t, OneVector[int]().IsDown())
	assert.False(t, UpVector[int]().IsDown())
	assert.True(t, DownVector[int]().IsDown())
	assert.False(t, LeftVector[int]().IsDown())
	assert.False(t, RightVector[int]().IsDown())
	assert.False(t, UpLeftVector().IsDown())
	assert.False(t, UpRightVector().IsDown())
	assert.True(t, DownLeftVector().IsDown())
	assert.True(t, DownRightVector().IsDown())
}

func TestVector_IsLeft(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsLeft())
	assert.False(t, OneVector[int]().IsLeft())
	assert.False(t, UpVector[int]().IsLeft())
	assert.False(t, DownVector[int]().IsLeft())
	assert.True(t, LeftVector[int]().IsLeft())
	assert.False(t, RightVector[int]().IsLeft())
	assert.True(t, UpLeftVector().IsLeft())
	assert.False(t, UpRightVector().IsLeft())
	assert.True(t, DownLeftVector().IsLeft())
	assert.False(t, DownRightVector().IsLeft())
}
func TestVector_IsRight(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsRight())
	assert.True(t, OneVector[int]().IsRight())
	assert.False(t, UpVector[int]().IsRight())
	assert.False(t, DownVector[int]().IsRight())
	assert.False(t, LeftVector[int]().IsRight())
	assert.True(t, RightVector[int]().IsRight())
	assert.False(t, UpLeftVector().IsRight())
	assert.True(t, UpRightVector().IsRight())
	assert.False(t, DownLeftVector().IsRight())
	assert.True(t, DownRightVector().IsRight())
}

func TestVector_IsNormalized(t *testing.T) {
	assert.False(t, vectorInt.IsNormalized())
	assert.True(t, vectorInt.Normalize().IsNormalized())
	assert.True(t, Vec(1.1, 2.1).Normalize().IsNormalized())
	assert.False(t, Vec(0, 0).IsNormalized())
	assert.True(t, Vec(1, 0).IsNormalized())
	assert.True(t, Vec(OneOverSqrt2, OneOverSqrt2).IsNormalized())

	assert.True(t, UpVector[float64]().IsNormalized())
	assert.True(t, DownVector[float64]().IsNormalized())
	assert.True(t, LeftVector[float64]().IsNormalized())
	assert.True(t, RightVector[float64]().IsNormalized())
	assert.True(t, UpLeftVector().IsNormalized())
	assert.True(t, UpRightVector().IsNormalized())
	assert.True(t, DownLeftVector().IsNormalized())
	assert.True(t, DownRightVector().IsNormalized())
}
func TestVector_Less(t *testing.T) {
	assert.False(t, vectorInt.Less(18))
	assert.True(t, vectorInt.Less(19))

	assert.False(t, vectorFloat.Less(0.1))
	assert.True(t, vectorFloat.Less(0.7))
}

func TestVector_XY(t *testing.T) {
	x1, y1 := vectorInt.XY()
	assert.Equal(t, x1, 10)
	assert.Equal(t, y1, 16)

	x2, y2 := vectorFloat.XY()
	assert.Equal(t, x2, x2)
	assert.Equal(t, y2, y2)
}

func TestVector_Point(t *testing.T) {
	assertPoint(t, vectorInt.Point(), 10, 16)
	assertPoint(t, vectorFloat.Point(), 0.6, -0.25)
}

func TestVector_Size(t *testing.T) {
	assertSize(t, vectorInt.Size(), 10, 16)
	assertSize(t, vectorFloat.Size(), 0.6, 0.25)
}

func TestVector_Int(t *testing.T) {
	assertVector(t, vectorInt.Int(), 10, 16)
	assertVector(t, vectorFloat.Int(), 1, 0)
}

func TestVector_Float(t *testing.T) {
	assertVector(t, vectorInt.Float(), 10.0, 16.0)
	assertVector(t, vectorFloat.Float(), 0.6, -0.25)
}

func TestVector_String(t *testing.T) {
	assert.Equal(t, vectorInt.String(), "⟨10,16⟩")
	assert.Equal(t, Vec(100, -34.0000115).String(), "⟨100,-34.00⟩")
}

func TestVector_Marshall(t *testing.T) {
	assert.JSON(t, vectorInt, `{"x":10,"y":16}`)
	assert.JSON(t, Vec(100, -34.0000115), `{"x":100.0,"y":-34.0000115}`)
}

func TestVector_Unmarshall(t *testing.T) {
	var p1 Vector[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &p1))
	assertVector(t, p1, 10, 16)

	var p2 Vector[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &p2))
	assertVector(t, p2, 10.1, -34.0000115)
}

func TestVector_Immutable(t *testing.T) {
	p1 := vectorInt
	p2 := Vec(3, -3)

	p1.Add(Vec(3, -2))
	p1.Subtract(p2)
	p1.Multiply(2)
	p1.MultiplyXY(3, 4)

	assertVector(t, p1, 10, 16)
	assertVector(t, p2, 3, -3)
}

func assertVector[T Number](t *testing.T, p Vector[T], x, y T) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(p.X), float64(x), Delta, "X: ") {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Y), float64(y), Delta, "Y: ") {
		ok = false
	}

	return ok
}

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

func TestVector_Constructor(t *testing.T) {
	AssertVector(t, vectorInt, 10, 16)
	AssertVector(t, Vec[float64](0.16, 204), 0.16, 204.0)

	AssertVector(t, ZeroVector[int](), 0, 0)
	AssertVector(t, OneVector[float64](), 1, 1)
}

func TestVector_FromAngle(t *testing.T) {
	// angle 0 → (length, 0)
	AssertVector(t, VectorFromAngle(0, 5.0), 5.0, 0.0)
	// angle π/2 → (0, length)
	v := VectorFromAngle(Pi/2, 1.0)
	assert.EqualDelta(t, v.X, 0.0, Delta)
	assert.EqualDelta(t, v.Y, 1.0, Delta)
	// angle π → (-length, 0)
	v2 := VectorFromAngle(Pi, 1.0)
	assert.EqualDelta(t, v2.X, -1.0, Delta)
	assert.EqualDelta(t, v2.Y, 0.0, Delta)
	// integer: cos(π/2)≈6e-17→0, sin(π/2)=1→1
	AssertVector(t, VectorFromAngle(Pi/2, 4), 0, 4)
}

func TestVector_Transform(t *testing.T) {
	AssertVector(t, vectorInt.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 48, 132)
	AssertVector(t, vectorFloat.Transform(Mat(1.1, 2.3, 3.3, 4.4, 5.5, 6.6)), 0.085, 1.265)
}

func TestVector_Add(t *testing.T) {
	AssertVector(t, vectorInt.Add(Vec(3, -2)), 13, 14)
	AssertVector(t, vectorInt.AddXY(3, -2), 13, 14)
	AssertVector(t, vectorFloat.Add(Vec(100.1, -0.1)), 100.7, -0.35)
	AssertVector(t, vectorFloat.AddXY(100.1, -0.1), 100.7, -0.35)
}

func TestVector_Sub(t *testing.T) {
	AssertVector(t, vectorInt.Subtract(Vec(3, -3)), 7, 19)
	AssertVector(t, vectorInt.SubtractXY(3, -3), 7, 19)
	AssertVector(t, vectorFloat.Subtract(Vec(100.1, -0.1)), -99.5, -0.15)
	AssertVector(t, vectorFloat.SubtractXY(100.1, -0.1), -99.5, -0.15)
}

func TestVector_Multiply(t *testing.T) {
	AssertVector(t, vectorInt.Multiply(3), 30, 48)
	AssertVector(t, vectorInt.MultiplyXY(3, 2), 30, 32)
	AssertVector(t, vectorFloat.Multiply(-1.5), -0.9, 0.375)
	AssertVector(t, vectorFloat.MultiplyXY(-1.5, 2), -0.9, -0.5)
}

func TestVector_Divide(t *testing.T) {
	AssertVector(t, vectorInt.Divide(2), 5, 8)
	AssertVector(t, vectorInt.DivideXY(3, 2), 3, 8)
	AssertVector(t, vectorFloat.Divide(-2), -0.3, 0.125)
	AssertVector(t, vectorFloat.DivideXY(-4, 0.5), -0.15, -0.5)
}

func TestVector_Negate(t *testing.T) {
	AssertVector(t, vectorInt.Negate(), -10, -16)
	AssertVector(t, vectorFloat.Negate(), -0.6, 0.25)
}

func TestVector_Rotate(t *testing.T) {
	AssertVector(t, Vec(1, 0).Rotate(ToRadians(90)), 0, 1)
	AssertVector(t, vectorFloat.Rotate(Pi), -0.6, 0.25)
}

func TestVector_Resize(t *testing.T) {
	AssertVector(t, vectorInt.Resize(5), 3, 4)
	AssertVector(t, vectorFloat.Resize(5), 4.615384, -1.923076)
}

func TestVector_Normalize(t *testing.T) {
	AssertVector(t, vectorInt.Normalize(), 0, 1)
	assert.Equal(t, vectorInt.Normalize().Length(), 1)
	AssertVector(t, vectorInt.Negate().Normalize(), 1, 0)
	assert.Equal(t, vectorInt.Negate().Normalize().Length(), 1)
	AssertVector(t, Vec(0, 0).Normalize(), 1, 0)
	AssertVector(t, vectorFloat.Normalize(), 0.923076, -0.384615)
	assert.Equal(t, vectorFloat.Normalize().Length(), 1)
}

func TestVector_Abs(t *testing.T) {
	AssertVector(t, vectorInt.Abs(), 10, 16)
	AssertVector(t, Vec(-1, -3).Abs(), 1, 3)
	AssertVector(t, vectorFloat.Abs(), 0.6, 0.25)
}

func TestVector_Round(t *testing.T) {
	AssertVector(t, Vec(1.4, -1.5).Round(), 1.0, -2.0)
	AssertVector(t, Vec(3, -2).Round(), 3, -2)
}

func TestVector_Floor(t *testing.T) {
	AssertVector(t, Vec(1.9, -1.1).Floor(), 1.0, -2.0)
	AssertVector(t, Vec(3, -2).Floor(), 3, -2)
}

func TestVector_Ceil(t *testing.T) {
	AssertVector(t, Vec(1.1, -1.9).Ceil(), 2.0, -1.0)
	AssertVector(t, Vec(3, -2).Ceil(), 3, -2)
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
	AssertVector(t, vectorInt.Normal(), -16, 10)
	AssertVector(t, vectorFloat.Normal(), 0.25, 0.6)
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
	assert.EqualDelta(t, Vec(0, 0).Angle(), 0, Delta)
	assert.EqualDelta(t, Vec(2, 2).Angle(), ToRadians(45), Delta)
	assert.EqualDelta(t, vectorFloat.Angle(), -0.39479111, Delta)
}

func TestVector_Direction(t *testing.T) {
	assert.Equal(t, Vec(3, 0).Direction(), DirectionRight)
	assert.Equal(t, Vec(2.0, -2.0).Direction(), DirectionUpRight)

	// a nearly-horizontal vector still snaps to Right
	assert.Equal(t, Vec(10.0, 1.0).Direction(), DirectionRight)
	assert.Equal(t, Vec(0, 0).Direction(), DirectionNone)
}

func TestVector_Lerp(t *testing.T) {
	AssertVector(t, vectorInt.Lerp(Vec(20, 20), 0.25), 13, 17)
	AssertVector(t, vectorFloat.Lerp(Vec(-10.0, 10.0), 0.25), -2.05, 2.3125)
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
	assert.True(t, DirectionUp.Unit[int]().IsUp())
	assert.False(t, DirectionDown.Unit[int]().IsUp())
	assert.False(t, DirectionLeft.Unit[int]().IsUp())
	assert.False(t, DirectionRight.Unit[int]().IsUp())
	assert.True(t, DirectionUpLeft.Unit[float64]().IsUp())
	assert.True(t, DirectionUpRight.Unit[float64]().IsUp())
	assert.False(t, DirectionDownLeft.Unit[float64]().IsUp())
	assert.False(t, DirectionDownRight.Unit[float64]().IsUp())
}

func TestVector_IsDown(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsDown())
	assert.True(t, OneVector[int]().IsDown())
	assert.False(t, DirectionUp.Unit[int]().IsDown())
	assert.True(t, DirectionDown.Unit[int]().IsDown())
	assert.False(t, DirectionLeft.Unit[int]().IsDown())
	assert.False(t, DirectionRight.Unit[int]().IsDown())
	assert.False(t, DirectionUpLeft.Unit[float64]().IsDown())
	assert.False(t, DirectionUpRight.Unit[float64]().IsDown())
	assert.True(t, DirectionDownLeft.Unit[float64]().IsDown())
	assert.True(t, DirectionDownRight.Unit[float64]().IsDown())
}

func TestVector_IsLeft(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsLeft())
	assert.False(t, OneVector[int]().IsLeft())
	assert.False(t, DirectionUp.Unit[int]().IsLeft())
	assert.False(t, DirectionDown.Unit[int]().IsLeft())
	assert.True(t, DirectionLeft.Unit[int]().IsLeft())
	assert.False(t, DirectionRight.Unit[int]().IsLeft())
	assert.True(t, DirectionUpLeft.Unit[float64]().IsLeft())
	assert.False(t, DirectionUpRight.Unit[float64]().IsLeft())
	assert.True(t, DirectionDownLeft.Unit[float64]().IsLeft())
	assert.False(t, DirectionDownRight.Unit[float64]().IsLeft())
}
func TestVector_IsRight(t *testing.T) {
	assert.False(t, ZeroVector[int]().IsRight())
	assert.True(t, OneVector[int]().IsRight())
	assert.False(t, DirectionUp.Unit[int]().IsRight())
	assert.False(t, DirectionDown.Unit[int]().IsRight())
	assert.False(t, DirectionLeft.Unit[int]().IsRight())
	assert.True(t, DirectionRight.Unit[int]().IsRight())
	assert.False(t, DirectionUpLeft.Unit[float64]().IsRight())
	assert.True(t, DirectionUpRight.Unit[float64]().IsRight())
	assert.False(t, DirectionDownLeft.Unit[float64]().IsRight())
	assert.True(t, DirectionDownRight.Unit[float64]().IsRight())
}

func TestVector_IsNormalized(t *testing.T) {
	assert.False(t, vectorInt.IsNormalized())
	assert.True(t, vectorInt.Normalize().IsNormalized())
	assert.True(t, Vec(1.1, 2.1).Normalize().IsNormalized())
	assert.False(t, Vec(0, 0).IsNormalized())
	assert.True(t, Vec(1, 0).IsNormalized())
	assert.True(t, Vec(OneOverSqrt2, OneOverSqrt2).IsNormalized())

	for _, direction := range Directions {
		assert.True(t, direction.Unit[float64]().IsNormalized(), direction.String())
	}
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
	AssertPoint(t, vectorInt.Point(), 10, 16)
	AssertPoint(t, vectorFloat.Point(), 0.6, -0.25)
}

func TestVector_Size(t *testing.T) {
	AssertSize(t, vectorInt.Size(), 10, 16)
	AssertSize(t, vectorFloat.Size(), 0.6, 0.25)
}

func TestVector_Int(t *testing.T) {
	AssertVector(t, vectorInt.Int(), 10, 16)
	AssertVector(t, vectorFloat.Int(), 1, 0)
}

func TestVector_Float(t *testing.T) {
	AssertVector(t, vectorInt.Float(), 10.0, 16.0)
	AssertVector(t, vectorFloat.Float(), 0.6, -0.25)
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
	var v1 Vector[int]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10,"y":16}`), &v1))
	AssertVector(t, v1, 10, 16)

	var v2 Vector[float64]
	assert.NoError(t, json.Unmarshal([]byte(`{"x":10.1,"y":-34.0000115}`), &v2))
	AssertVector(t, v2, 10.1, -34.0000115)
}

func TestVector_Immutable(t *testing.T) {
	v1 := vectorInt
	v2 := Vec(3, -3)

	v1.Add(Vec(3, -2))
	v1.Subtract(v2)
	v1.Multiply(2)
	v1.MultiplyXY(3, 4)

	AssertVector(t, v1, 10, 16)
	AssertVector(t, v2, 3, -3)
}

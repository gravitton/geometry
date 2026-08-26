package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestDirection_Constructor(t *testing.T) {
	assert.Equal(t, DirectionFromAngle(0), DirectionRight)
	assert.Equal(t, DirectionFromAngle(-Pi/2), DirectionUp)
	assert.Equal(t, DirectionFromAngle(Pi/2), DirectionDown)
	assert.Equal(t, DirectionFromAngle(Pi), DirectionLeft)
	assert.Equal(t, DirectionFromAngle(-Pi), DirectionLeft)

	// diagonals
	assert.Equal(t, DirectionFromAngle(Pi/4), DirectionDownRight)
	assert.Equal(t, DirectionFromAngle(-Pi/4), DirectionUpRight)
	assert.Equal(t, DirectionFromAngle(3*Pi/4), DirectionDownLeft)
	assert.Equal(t, DirectionFromAngle(-3*Pi/4), DirectionUpLeft)

	// snaps to the nearest eighth turn, symmetrically from either side
	assert.Equal(t, DirectionFromAngle(-Pi/2+0.1), DirectionUp)
	assert.Equal(t, DirectionFromAngle(-Pi/2-0.1), DirectionUp)
	assert.Equal(t, DirectionFromAngle(Pi/4+0.1), DirectionDownRight)
	assert.Equal(t, DirectionFromAngle(Pi/4-0.1), DirectionDownRight)

	assert.Equal(t, DirectionFromAxes(true, false, false, false), DirectionUp)
	assert.Equal(t, DirectionFromAxes(false, true, false, false), DirectionDown)
	assert.Equal(t, DirectionFromAxes(false, false, true, false), DirectionLeft)
	assert.Equal(t, DirectionFromAxes(false, false, false, true), DirectionRight)

	assert.Equal(t, DirectionFromAxes(true, false, true, false), DirectionUpLeft)
	assert.Equal(t, DirectionFromAxes(true, false, false, true), DirectionUpRight)
	assert.Equal(t, DirectionFromAxes(false, true, true, false), DirectionDownLeft)
	assert.Equal(t, DirectionFromAxes(false, true, false, true), DirectionDownRight)

	// opposite inputs cancel out
	assert.Equal(t, DirectionFromAxes(true, true, false, false), DirectionNone)
	assert.Equal(t, DirectionFromAxes(false, false, true, true), DirectionNone)
	assert.Equal(t, DirectionFromAxes(true, true, true, true), DirectionNone)
	assert.Equal(t, DirectionFromAxes(false, false, false, false), DirectionNone)
}

func TestDirection_Aliases(t *testing.T) {
	assert.Equal(t, East, DirectionRight)
	assert.Equal(t, North, DirectionUp)
	assert.Equal(t, West, DirectionLeft)
	assert.Equal(t, South, DirectionDown)
	assert.Equal(t, NorthEast, DirectionUpRight)
	assert.Equal(t, SouthWest, DirectionDownLeft)

	assert.Equal(t, Top, DirectionUp)
	assert.Equal(t, Bottom, DirectionDown)
	assert.Equal(t, TopLeft, DirectionUpLeft)
	assert.Equal(t, TopRight, DirectionUpRight)
	assert.Equal(t, BottomLeft, DirectionDownLeft)
	assert.Equal(t, BottomRight, DirectionDownRight)
}

func TestDirection_Order(t *testing.T) {
	assert.Equal(t, len(Directions), 8)

	for i, direction := range Directions {
		assert.Equal(t, int(direction), i)
	}

	// each step is an eighth turn counterclockwise, so the offset rotates accordingly
	for i, direction := range Directions {
		AssertVector(t, direction.Offset[int](), directionOffsets[i].X, directionOffsets[i].Y, direction.String())
	}
}

func TestDirection_Opposite(t *testing.T) {
	for _, direction := range Directions {
		assert.Equal(t, direction.Opposite(), direction.Rotate(4), direction.String())
		assert.Equal(t, direction.Opposite().Opposite(), direction, direction.String())
		AssertVector(t, direction.Offset[int]().Add(direction.Opposite().Offset[int]()), 0, 0, direction.String())
	}

	assert.Equal(t, DirectionRight.Opposite(), DirectionLeft)
	assert.Equal(t, DirectionUp.Opposite(), DirectionDown)
	assert.Equal(t, DirectionUpRight.Opposite(), DirectionDownLeft)
}

func TestDirection_Rotate(t *testing.T) {
	assert.Equal(t, DirectionRight.Rotate(1), DirectionUpRight)
	assert.Equal(t, DirectionRight.Rotate(2), DirectionUp)
	assert.Equal(t, DirectionRight.Rotate(8), DirectionRight)
	assert.Equal(t, DirectionRight.Rotate(-1), DirectionDownRight)
	assert.Equal(t, DirectionRight.Rotate(-9), DirectionDownRight)
	assert.Equal(t, DirectionUp.Rotate(-2), DirectionRight)
}

func TestDirection_Axis(t *testing.T) {
	assert.Equal(t, DirectionRight.Axis(), AxisHorizontal)
	assert.Equal(t, DirectionDown.Axis(), AxisVertical)
	assert.Equal(t, DirectionUpLeft.Axis(), AxisNone)
}

func TestDirection_Offset(t *testing.T) {
	AssertVector(t, DirectionRight.Offset[int](), 1, 0)
	AssertVector(t, DirectionUp.Offset[int](), 0, -1)
	AssertVector(t, DirectionDownRight.Offset[int](), 1, 1)

	// offsets stay exact for integers; diagonals are not normalized
	AssertVector(t, DirectionUpLeft.Offset[float64](), -1, -1)
}

func TestDirection_Unit(t *testing.T) {
	AssertVector(t, DirectionRight.Unit[float64](), 1, 0)
	AssertVector(t, DirectionUp.Unit[float64](), 0, -1)
	AssertVector(t, DirectionUpRight.Unit[float64](), OneOverSqrt2, -OneOverSqrt2)

	for _, direction := range Directions {
		assert.EqualDelta(t, direction.Unit[float64]().Length(), 1.0, Delta, direction.String())
	}

	AssertVector(t, DirectionRight.Unit[float32](), 1, 0)

	// integer T: cardinals stay exact, diagonals snap to an axis-aligned unit vector
	AssertVector(t, DirectionRight.Unit[int](), 1, 0)
	AssertVector(t, DirectionUpRight.Unit[int](), 1, -1)
}

func TestDirection_Vector(t *testing.T) {
	AssertVector(t, DirectionRight.Vector(5.0), 5, 0)
	AssertVector(t, DirectionDown.Vector(5.0), 0, 5)
	AssertVector(t, DirectionLeft.Vector(3), -3, 0)

	assert.EqualDelta(t, DirectionDownLeft.Vector(4.0).Length(), 4.0, Delta)
}

func TestDirection_Angle(t *testing.T) {
	assert.EqualDelta(t, DirectionRight.Angle(), 0.0, Delta)
	assert.EqualDelta(t, DirectionUp.Angle(), -Pi/2, Delta)
	assert.EqualDelta(t, DirectionUpRight.Angle(), -Pi/4, Delta)
	assert.EqualDelta(t, DirectionDown.Angle(), Pi/2, Delta)

	for _, direction := range Directions {
		assert.Equal(t, DirectionFromAngle(direction.Angle()), direction, direction.String())
	}
}

func TestDirection_IsNone(t *testing.T) {
	assert.True(t, DirectionNone.IsNone())
	assert.False(t, DirectionRight.IsNone())

	assert.Equal(t, DirectionNone.Opposite(), DirectionNone)
	assert.Equal(t, DirectionNone.Rotate(3), DirectionNone)
	assert.False(t, DirectionNone.IsCardinal())
	assert.False(t, DirectionNone.IsDiagonal())
	assert.False(t, DirectionNone.IsPositive())
	assert.Equal(t, DirectionNone.Axis(), AxisNone)

	AssertVector(t, DirectionNone.Offset[int](), 0, 0)
	AssertVector(t, DirectionNone.Unit[float64](), 0, 0)
	AssertVector(t, DirectionNone.Vector(5.0), 0, 0)
	assert.EqualDelta(t, DirectionNone.Angle(), 0.0, Delta)

	// the zero value is Right, so a zero-valued field means "rightward", not "unset"
	assert.Equal(t, Direction(0), DirectionRight)
}

func TestDirection_IsCardinalDiagonal(t *testing.T) {
	for _, direction := range CardinalDirections {
		assert.True(t, direction.IsCardinal(), direction.String())
		assert.False(t, direction.IsDiagonal(), direction.String())

		assert.Equal(t, direction.Axis(), Axis(int(direction)/2%2), direction.String())
	}

	for _, direction := range DiagonalDirections {
		assert.False(t, direction.IsCardinal(), direction.String())
		assert.True(t, direction.IsDiagonal(), direction.String())

		assert.Equal(t, direction.Axis(), AxisNone, direction.String())
	}
}

func TestDirection_IsPositive(t *testing.T) {
	assert.True(t, DirectionRight.IsPositive())
	assert.True(t, DirectionDown.IsPositive())
	assert.False(t, DirectionLeft.IsPositive())
	assert.False(t, DirectionUp.IsPositive())

	for _, direction := range DiagonalDirections {
		assert.False(t, direction.IsPositive(), direction.String())
	}
}

func TestDirection_String(t *testing.T) {
	assert.Equal(t, DirectionRight.String(), "Right")
	assert.Equal(t, DirectionUpRight.String(), "UpRight")
	assert.Equal(t, DirectionUp.String(), "Up")
	assert.Equal(t, DirectionUpLeft.String(), "UpLeft")
	assert.Equal(t, DirectionLeft.String(), "Left")
	assert.Equal(t, DirectionDownLeft.String(), "DownLeft")
	assert.Equal(t, DirectionDown.String(), "Down")
	assert.Equal(t, DirectionDownRight.String(), "DownRight")
	assert.Equal(t, DirectionNone.String(), "None")
	assert.Equal(t, Direction(9).String(), "UpRight")
}

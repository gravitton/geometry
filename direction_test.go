package geom

import (
	"math"
	"testing"

	"github.com/gravitton/assert"
)

func TestDirectionFromAngle(t *testing.T) {
	t.Run("cardinal angles", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(0), DirectionRight)
		assert.Equal(t, DirectionFromAngle(-Pi/2), DirectionUp)
		assert.Equal(t, DirectionFromAngle(Pi/2), DirectionDown)
		assert.Equal(t, DirectionFromAngle(Pi), DirectionLeft)
		assert.Equal(t, DirectionFromAngle(-Pi), DirectionLeft)
	})
	t.Run("diagonal angles", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(Pi/4), DirectionDownRight)
		assert.Equal(t, DirectionFromAngle(-Pi/4), DirectionUpRight)
		assert.Equal(t, DirectionFromAngle(3*Pi/4), DirectionDownLeft)
		assert.Equal(t, DirectionFromAngle(-3*Pi/4), DirectionUpLeft)
	})
	t.Run("snaps to the nearest eighth turn", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(-Pi/2+0.1), DirectionUp)
		assert.Equal(t, DirectionFromAngle(-Pi/2-0.1), DirectionUp)
		assert.Equal(t, DirectionFromAngle(Pi/4+0.1), DirectionDownRight)
		assert.Equal(t, DirectionFromAngle(Pi/4-0.1), DirectionDownRight)
	})
	t.Run("NaN has no direction", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(math.NaN()), DirectionNone)
	})
}

func TestDirectionFromAxes(t *testing.T) {
	t.Run("one axis gives a cardinal", func(t *testing.T) {
		assert.Equal(t, DirectionFromAxes(true, false, false, false), DirectionUp)
		assert.Equal(t, DirectionFromAxes(false, true, false, false), DirectionDown)
		assert.Equal(t, DirectionFromAxes(false, false, true, false), DirectionLeft)
		assert.Equal(t, DirectionFromAxes(false, false, false, true), DirectionRight)
	})
	t.Run("two axes give a diagonal", func(t *testing.T) {
		assert.Equal(t, DirectionFromAxes(true, false, true, false), DirectionUpLeft)
		assert.Equal(t, DirectionFromAxes(true, false, false, true), DirectionUpRight)
		assert.Equal(t, DirectionFromAxes(false, true, true, false), DirectionDownLeft)
		assert.Equal(t, DirectionFromAxes(false, true, false, true), DirectionDownRight)
	})
	t.Run("opposite inputs cancel out", func(t *testing.T) {
		assert.Equal(t, DirectionFromAxes(true, true, false, false), DirectionNone)
		assert.Equal(t, DirectionFromAxes(false, false, true, true), DirectionNone)
		assert.Equal(t, DirectionFromAxes(true, true, true, true), DirectionNone)
		assert.Equal(t, DirectionFromAxes(false, false, false, false), DirectionNone)
	})
}

func TestDirection_Opposite(t *testing.T) {
	t.Run("turns by a half turn", func(t *testing.T) {
		assert.Equal(t, DirectionRight.Opposite(), DirectionLeft)
		assert.Equal(t, DirectionUp.Opposite(), DirectionDown)
		assert.Equal(t, DirectionUpRight.Opposite(), DirectionDownLeft)
	})
	t.Run("none has no opposite", func(t *testing.T) {
		assert.Equal(t, DirectionNone.Opposite(), DirectionNone)
	})
}

func TestDirection_Rotate(t *testing.T) {
	t.Run("positive steps increase the angle", func(t *testing.T) {
		// clockwise as drawn on screen, counterclockwise in math coordinates
		assert.Equal(t, DirectionRight.Rotate(1), DirectionDownRight)
		assert.Equal(t, DirectionRight.Rotate(2), DirectionDown)
		assert.Equal(t, DirectionUp.Rotate(2), DirectionRight)
	})
	t.Run("negative steps decrease it", func(t *testing.T) {
		assert.Equal(t, DirectionRight.Rotate(-1), DirectionUpRight)
		assert.Equal(t, DirectionRight.Rotate(-9), DirectionUpRight)
	})
	t.Run("a full turn is the identity", func(t *testing.T) {
		assert.Equal(t, DirectionRight.Rotate(8), DirectionRight)
	})
	t.Run("none does not rotate", func(t *testing.T) {
		assert.Equal(t, DirectionNone.Rotate(3), DirectionNone)
	})
}

func TestDirection_Axis(t *testing.T) {
	t.Run("cardinals have an axis", func(t *testing.T) {
		assert.Equal(t, DirectionRight.Axis(), AxisHorizontal)
		assert.Equal(t, DirectionLeft.Axis(), AxisHorizontal)
		assert.Equal(t, DirectionDown.Axis(), AxisVertical)
		assert.Equal(t, DirectionUp.Axis(), AxisVertical)
	})
	t.Run("diagonals and none do not", func(t *testing.T) {
		assert.Equal(t, DirectionUpLeft.Axis(), AxisNone)
		assert.Equal(t, DirectionNone.Axis(), AxisNone)
	})
}

func TestDirection_Offset(t *testing.T) {
	t.Run("unit steps on the lattice", func(t *testing.T) {
		AssertVector(t, DirectionRight.Offset[int](), Vec(1, 0))
		AssertVector(t, DirectionUp.Offset[int](), Vec(0, -1))
		AssertVector(t, DirectionDownRight.Offset[int](), Vec(1, 1))
	})
	t.Run("diagonals are not normalized", func(t *testing.T) {
		AssertVector(t, DirectionUpLeft.Offset[float64](), Vec(-1.0, -1.0))
	})
	t.Run("none is the zero vector", func(t *testing.T) {
		AssertVector(t, DirectionNone.Offset[int](), Vec(0, 0))
	})
}

func TestDirection_Unit(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		AssertVector(t, DirectionRight.Unit[float64](), Vec(1.0, 0.0))
		AssertVector(t, DirectionUp.Unit[float64](), Vec(0.0, -1.0))
		AssertVector(t, DirectionUpRight.Unit[float64](), Vec(OneOverSqrt2, -OneOverSqrt2))

		AssertVector(t, DirectionRight.Unit[float32](), Vec[float32](1, 0))
	})
	t.Run("int keeps the cardinals", func(t *testing.T) {
		AssertVector(t, DirectionRight.Unit[int](), Vec(1, 0))
		AssertVector(t, DirectionUp.Unit[int](), Vec(0, -1))
	})
	t.Run("int collapses the diagonals to an axis", func(t *testing.T) {
		// no integer diagonal has length 1, so a diagonal cannot survive normalization;
		// Offset keeps the lattice step (±1,±1) for callers that need the direction
		AssertVector(t, DirectionUpRight.Unit[int](), Vec(1, 0))
		AssertVector(t, DirectionUpRight.Offset[int](), Vec(1, -1))
	})
	t.Run("none is the zero vector", func(t *testing.T) {
		AssertVector(t, DirectionNone.Unit[float64](), Vec(0.0, 0.0))
	})
}

func TestDirection_Vector(t *testing.T) {
	t.Run("scales the unit vector", func(t *testing.T) {
		AssertVector(t, DirectionRight.Vector(5.0), Vec(5.0, 0.0))
		AssertVector(t, DirectionDown.Vector(5.0), Vec(0.0, 5.0))
		AssertVector(t, DirectionLeft.Vector(3), Vec(-3, 0))
	})
	t.Run("a diagonal keeps the requested length", func(t *testing.T) {
		AssertNumber(t, DirectionDownLeft.Vector(4.0).Length(), 4.0)
	})
	t.Run("none is the zero vector", func(t *testing.T) {
		AssertVector(t, DirectionNone.Vector(5.0), Vec(0.0, 0.0))
	})
}

func TestDirection_Angle(t *testing.T) {
	t.Run("measured from the positive X axis", func(t *testing.T) {
		AssertNumber(t, DirectionRight.Angle(), 0.0)
		AssertNumber(t, DirectionDown.Angle(), Pi/2)
		AssertNumber(t, DirectionUp.Angle(), -Pi/2)
		AssertNumber(t, DirectionUpRight.Angle(), -Pi/4)
	})
	t.Run("none has no angle", func(t *testing.T) {
		AssertNumber(t, DirectionNone.Angle(), 0.0)
	})
}

func TestDirection_IsNone(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assert.True(t, DirectionNone.IsNone())
	})
	t.Run("any real direction", func(t *testing.T) {
		for _, direction := range Directions {
			assert.False(t, direction.IsNone(), direction.String()+": ")
		}
	})
	t.Run("the zero value is right, not unset", func(t *testing.T) {
		assert.Equal(t, Direction(0), DirectionRight)
		assert.False(t, Direction(0).IsNone())
	})
}

func TestDirection_IsCardinal(t *testing.T) {
	t.Run("cardinals", func(t *testing.T) {
		for _, direction := range CardinalDirections {
			assert.True(t, direction.IsCardinal(), direction.String()+": ")
		}
	})
	t.Run("diagonals and none", func(t *testing.T) {
		for _, direction := range DiagonalDirections {
			assert.False(t, direction.IsCardinal(), direction.String()+": ")
		}

		assert.False(t, DirectionNone.IsCardinal())
	})
}

func TestDirection_IsDiagonal(t *testing.T) {
	t.Run("diagonals", func(t *testing.T) {
		for _, direction := range DiagonalDirections {
			assert.True(t, direction.IsDiagonal(), direction.String()+": ")
		}
	})
	t.Run("cardinals and none", func(t *testing.T) {
		for _, direction := range CardinalDirections {
			assert.False(t, direction.IsDiagonal(), direction.String()+": ")
		}

		assert.False(t, DirectionNone.IsDiagonal())
	})
}

func TestDirection_IsPositive(t *testing.T) {
	t.Run("right and down grow the coordinate", func(t *testing.T) {
		assert.True(t, DirectionRight.IsPositive())
		assert.True(t, DirectionDown.IsPositive())
	})
	t.Run("left and up shrink it", func(t *testing.T) {
		assert.False(t, DirectionLeft.IsPositive())
		assert.False(t, DirectionUp.IsPositive())
	})
	t.Run("diagonals and none have no sign", func(t *testing.T) {
		for _, direction := range DiagonalDirections {
			assert.False(t, direction.IsPositive(), direction.String()+": ")
		}

		assert.False(t, DirectionNone.IsPositive())
	})
}

func TestDirection_String(t *testing.T) {
	t.Run("every direction", func(t *testing.T) {
		assert.Equal(t, DirectionRight.String(), "Right")
		assert.Equal(t, DirectionDownRight.String(), "DownRight")
		assert.Equal(t, DirectionDown.String(), "Down")
		assert.Equal(t, DirectionDownLeft.String(), "DownLeft")
		assert.Equal(t, DirectionLeft.String(), "Left")
		assert.Equal(t, DirectionUpLeft.String(), "UpLeft")
		assert.Equal(t, DirectionUp.String(), "Up")
		assert.Equal(t, DirectionUpRight.String(), "UpRight")
	})
	t.Run("none", func(t *testing.T) {
		assert.Equal(t, DirectionNone.String(), "None")
	})
	t.Run("out-of-range values wrap", func(t *testing.T) {
		assert.Equal(t, Direction(9).String(), "DownRight")
	})
}

func TestDirection_Properties(t *testing.T) {
	t.Run("the compass aliases name the same directions", func(t *testing.T) {
		assert.Equal(t, East, DirectionRight)
		assert.Equal(t, North, DirectionUp)
		assert.Equal(t, West, DirectionLeft)
		assert.Equal(t, South, DirectionDown)
		assert.Equal(t, NorthEast, DirectionUpRight)
		assert.Equal(t, SouthWest, DirectionDownLeft)
	})
	t.Run("the screen aliases name the same directions", func(t *testing.T) {
		assert.Equal(t, Top, DirectionUp)
		assert.Equal(t, Bottom, DirectionDown)
		assert.Equal(t, TopLeft, DirectionUpLeft)
		assert.Equal(t, TopRight, DirectionUpRight)
		assert.Equal(t, BottomLeft, DirectionDownLeft)
		assert.Equal(t, BottomRight, DirectionDownRight)
	})
	t.Run("directions are ordered by increasing angle", func(t *testing.T) {
		assert.Equal(t, len(Directions), 8)

		for i, direction := range Directions {
			assert.Equal(t, int(direction), i, direction.String()+": ")
			AssertVector(t, direction.Offset[int](), directionOffsets[i], direction.String()+": ")
		}
	})
	t.Run("opposite is four steps and its own inverse", func(t *testing.T) {
		for _, direction := range Directions {
			assert.Equal(t, direction.Opposite(), direction.Rotate(4), direction.String()+": ")
			assert.Equal(t, direction.Opposite().Opposite(), direction, direction.String()+": ")

			AssertVector(t, direction.Offset[int]().Add(direction.Opposite().Offset[int]()), Vec(0, 0), direction.String()+": ")
		}
	})
	t.Run("angle round-trips through the constructor", func(t *testing.T) {
		for _, direction := range Directions {
			assert.Equal(t, DirectionFromAngle(direction.Angle()), direction, direction.String()+": ")
		}
	})
	t.Run("one step is one eighth turn of increasing angle", func(t *testing.T) {
		for _, direction := range Directions {
			for steps := -8; steps <= 8; steps++ {
				rotated := DirectionFromAngle(direction.Angle() + float64(steps)*Pi/4)

				assert.Equal(t, rotated, direction.Rotate(steps), direction.String()+": ")
			}
		}
	})
	t.Run("a step turns the same way as vector rotation", func(t *testing.T) {
		for _, direction := range Directions {
			rotated := direction.Offset[float64]().Rotate(Pi / 2)

			assert.Equal(t, rotated.Direction(), direction.Rotate(2), direction.String()+": ")
		}
	})
	t.Run("every unit vector has length one", func(t *testing.T) {
		for _, direction := range Directions {
			AssertNumber(t, direction.Unit[float64]().Length(), 1.0, direction.String()+": ")
		}
	})
	t.Run("cardinal and diagonal partition the directions", func(t *testing.T) {
		for _, direction := range Directions {
			assert.NotEqual(t, direction.IsCardinal(), direction.IsDiagonal(), direction.String()+": ")
		}
	})
	t.Run("a cardinal axis follows its position in the order", func(t *testing.T) {
		for _, direction := range CardinalDirections {
			assert.Equal(t, direction.Axis(), Axis(int(direction)/2%2), direction.String()+": ")
		}
	})
}

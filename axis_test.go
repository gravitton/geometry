package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestAxis_Cross(t *testing.T) {
	assert.Equal(t, AxisHorizontal.Cross(), AxisVertical)
	assert.Equal(t, AxisVertical.Cross(), AxisHorizontal)

	for _, axis := range Axes {
		assert.Equal(t, axis.Cross().Cross(), axis, axis.String())
	}

	assert.Equal(t, AxisNone.Cross(), AxisNone)
}

func TestAxis_Direction(t *testing.T) {
	assert.Equal(t, AxisHorizontal.Direction(true), DirectionRight)
	assert.Equal(t, AxisHorizontal.Direction(false), DirectionLeft)
	assert.Equal(t, AxisVertical.Direction(true), DirectionDown)
	assert.Equal(t, AxisVertical.Direction(false), DirectionUp)
	assert.Equal(t, AxisNone.Direction(true), DirectionNone)

	for _, axis := range Axes {
		for _, positive := range []bool{true, false} {
			direction := axis.Direction(positive)

			assert.Equal(t, direction.Axis(), axis, direction.String())
			assert.Equal(t, direction.IsPositive(), positive, direction.String())
		}
	}
}

func TestAxis_AlongAcross(t *testing.T) {
	size := Sz(10, 20)

	assert.Equal(t, AxisHorizontal.Along(size), 10)
	assert.Equal(t, AxisHorizontal.Across(size), 20)
	assert.Equal(t, AxisVertical.Along(size), 20)
	assert.Equal(t, AxisVertical.Across(size), 10)

	// Size and Along are inverses on both axes
	for _, axis := range Axes {
		built := axis.Size(1.0, 2.0)

		assert.EqualDelta(t, axis.Along(built), 1.0, Delta, axis.String())
		assert.EqualDelta(t, axis.Across(built), 2.0, Delta, axis.String())
	}
}

func TestAxis_Project(t *testing.T) {
	assert.Equal(t, AxisHorizontal.Project(Vec(3, 4)), 3)
	assert.Equal(t, AxisVertical.Project(Vec(3, 4)), 4)
	assert.Equal(t, AxisVertical.Project(DirectionRight.Offset[int]()), 0)
	assert.Equal(t, AxisVertical.Project(DirectionDown.Offset[int]()), 1)
}

func TestAxis_ScaleAlong(t *testing.T) {
	AssertSize(t, AxisHorizontal.ScaleAlong(Sz(10, 20), 0.5), 5, 20)
	AssertSize(t, AxisVertical.ScaleAlong(Sz(10, 20), 0.5), 10, 10)
	AssertSize(t, AxisHorizontal.ScaleAlong(Sz(1.0, 2.0), 3), 3.0, 2.0)
}

func TestAxis_Vector(t *testing.T) {
	AssertVector(t, AxisHorizontal.Vector(3, 4), 3, 4)
	AssertVector(t, AxisVertical.Vector(3, 4), 4, 3)
	AssertVector(t, AxisVertical.Vector(1.5, 2.5), 2.5, 1.5)
}

func TestAxis_Size(t *testing.T) {
	AssertSize(t, AxisHorizontal.Size(3, 4), 3, 4)
	AssertSize(t, AxisVertical.Size(3, 4), 4, 3)
	AssertSize(t, AxisVertical.Size(1.5, 2.5), 2.5, 1.5)
}

func TestAxis_IsNone(t *testing.T) {
	assert.False(t, AxisHorizontal.IsNone())
	assert.False(t, AxisVertical.IsNone())
	assert.True(t, AxisNone.IsNone())
	assert.True(t, Axis(2).IsNone())

	assert.Equal(t, Axis(0), AxisHorizontal)
}

func TestAxis_String(t *testing.T) {
	assert.Equal(t, AxisHorizontal.String(), "Horizontal")
	assert.Equal(t, AxisVertical.String(), "Vertical")
	assert.Equal(t, AxisNone.String(), "None")
	assert.Equal(t, Axis(2).String(), "None")
}

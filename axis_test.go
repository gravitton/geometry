package geom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestParseAxis(t *testing.T) {
	t.Run("every name round-trips", func(t *testing.T) {
		for _, axis := range Axes() {
			parsed, err := ParseAxis(axis.String())

			assert.NoError(t, err, axis.String()+": ")
			assert.Equal(t, parsed, axis, axis.String()+": ")
		}
	})
	t.Run("none", func(t *testing.T) {
		parsed, err := ParseAxis("None")

		assert.NoError(t, err)
		assert.Equal(t, parsed, AxisNone)
	})
	t.Run("unknown is an error", func(t *testing.T) {
		_, err := ParseAxis("Diagonal")

		assert.Error(t, err)
	})
}

func TestAxis_Cross(t *testing.T) {
	t.Run("swaps the two axes", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Cross(), AxisVertical)
		assert.Equal(t, AxisVertical.Cross(), AxisHorizontal)
	})
	t.Run("none has no cross", func(t *testing.T) {
		assert.Equal(t, AxisNone.Cross(), AxisNone)
		assert.Equal(t, Axis(2).Cross(), AxisNone)
	})
}

func TestAxis_Direction(t *testing.T) {
	t.Run("positive points right and down", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Direction(true), DirectionRight)
		assert.Equal(t, AxisVertical.Direction(true), DirectionDown)
	})
	t.Run("negative points left and up", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Direction(false), DirectionLeft)
		assert.Equal(t, AxisVertical.Direction(false), DirectionUp)
	})
	t.Run("none has no direction", func(t *testing.T) {
		assert.Equal(t, AxisNone.Direction(true), DirectionNone)
		assert.Equal(t, AxisNone.Direction(false), DirectionNone)
	})
}

func TestAxis_Along(t *testing.T) {
	size := Sz(10, 20)

	t.Run("picks the extent on the axis", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Along(size), 10)
		assert.Equal(t, AxisVertical.Along(size), 20)
	})
	t.Run("none has no extent", func(t *testing.T) {
		assert.Equal(t, AxisNone.Along(size), 0)
	})
}

func TestAxis_Across(t *testing.T) {
	size := Sz(10, 20)

	t.Run("picks the extent on the axis", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Across(size), 20)
		assert.Equal(t, AxisVertical.Across(size), 10)
	})
	t.Run("none has no extent", func(t *testing.T) {
		assert.Equal(t, AxisNone.Across(size), 0)
	})
}

func TestAxis_Project(t *testing.T) {
	t.Run("picks the component on the main axis", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Project(Vec(3, 4)), 3)
		assert.Equal(t, AxisVertical.Project(Vec(3, 4)), 4)
	})
	t.Run("a perpendicular direction projects to zero", func(t *testing.T) {
		assert.Equal(t, AxisVertical.Project(DirectionRight.Offset[int]()), 0)
		assert.Equal(t, AxisVertical.Project(DirectionDown.Offset[int]()), 1)
	})
	t.Run("keeps the sign", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.Project(Vec(-3, 2)), -3)
		assert.Equal(t, AxisVertical.Project(Vec(3, -4.5)), -4.5)
	})
	t.Run("none projects to zero", func(t *testing.T) {
		assert.Equal(t, AxisNone.Project(Vec(3, 4)), 0)
	})
}

func TestAxis_ScaleAlong(t *testing.T) {
	t.Run("leaves the cross axis alone", func(t *testing.T) {
		AssertSize(t, AxisHorizontal.ScaleAlong(Sz(10, 20), 0.5), Sz(5, 20))
		AssertSize(t, AxisVertical.ScaleAlong(Sz(10, 20), 0.5), Sz(10, 10))
		AssertSize(t, AxisHorizontal.ScaleAlong(Sz(1.0, 2.0), 3), Sz(3.0, 2.0))
	})
	t.Run("none leaves the size unchanged", func(t *testing.T) {
		AssertSize(t, AxisNone.ScaleAlong(Sz(10, 20), 0.5), Sz(10, 20))
	})
}

func TestAxis_Vector(t *testing.T) {
	t.Run("horizontal keeps the order", func(t *testing.T) {
		AssertVector(t, AxisHorizontal.Vector(3, 4), Vec(3, 4))
	})
	t.Run("vertical swaps it", func(t *testing.T) {
		AssertVector(t, AxisVertical.Vector(3, 4), Vec(4, 3))
		AssertVector(t, AxisVertical.Vector(1.5, 2.5), Vec(2.5, 1.5))
	})
	t.Run("none is the zero vector", func(t *testing.T) {
		AssertVector(t, AxisNone.Vector(3, 4), Vec(0, 0))
	})
}

func TestAxis_Size(t *testing.T) {
	t.Run("horizontal keeps the order", func(t *testing.T) {
		AssertSize(t, AxisHorizontal.Size(3, 4), Sz(3, 4))
	})
	t.Run("vertical swaps it", func(t *testing.T) {
		AssertSize(t, AxisVertical.Size(3, 4), Sz(4, 3))
		AssertSize(t, AxisVertical.Size(1.5, 2.5), Sz(2.5, 1.5))
	})
	t.Run("keeps the sign, like Sz", func(t *testing.T) {
		AssertSize(t, AxisHorizontal.Size(-3, 4), Sz(-3, 4))
		AssertSize(t, AxisHorizontal.ScaleAlong(Sz(10, 4), -1), Sz(10, 4).ScaleXY(-1, 1))
	})
	t.Run("none has no extent", func(t *testing.T) {
		AssertSize(t, AxisNone.Size(3, 4), Sz(0, 0))
	})
}

func TestAxis_IsNone(t *testing.T) {
	t.Run("the two axes", func(t *testing.T) {
		assert.False(t, AxisHorizontal.IsNone())
		assert.False(t, AxisVertical.IsNone())
		assert.Equal(t, Axis(0), AxisHorizontal)
	})
	t.Run("anything else", func(t *testing.T) {
		assert.True(t, AxisNone.IsNone())
		assert.True(t, Axis(2).IsNone())
	})
}

func TestAxis_String(t *testing.T) {
	t.Run("named axes", func(t *testing.T) {
		assert.Equal(t, AxisHorizontal.String(), "Horizontal")
		assert.Equal(t, AxisVertical.String(), "Vertical")
	})
	t.Run("none and out of range", func(t *testing.T) {
		assert.Equal(t, AxisNone.String(), "None")
		assert.Equal(t, Axis(2).String(), "None")
	})
}

func TestAxis_JSON(t *testing.T) {
	t.Run("wire format is the name", func(t *testing.T) {
		assert.JSON(t, AxisVertical, `"Vertical"`)
		assert.JSON(t, AxisNone, `"None"`)
	})
	t.Run("round-trip", func(t *testing.T) {
		axes := Axes()
		for _, axis := range append(axes[:], AxisNone) {
			data, err := json.Marshal(axis)
			assert.NoError(t, err)

			var decoded Axis
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, axis, axis.String()+": ")
		}
	})
	t.Run("an unknown name is an error", func(t *testing.T) {
		var decoded Axis

		assert.Error(t, json.Unmarshal([]byte(`"Diagonal"`), &decoded))
	})
}

func TestAxis_Properties(t *testing.T) {
	t.Run("cross is its own inverse", func(t *testing.T) {
		for _, axis := range Axes() {
			assert.Equal(t, axis.Cross().Cross(), axis, axis.String()+": ")
		}
	})
	t.Run("along and across swap on the cross axis", func(t *testing.T) {
		size := Sz(10.0, 20.0)

		for _, axis := range Axes() {
			AssertNumber(t, axis.Cross().Along(size), axis.Across(size), axis.String()+": ")
			AssertNumber(t, axis.Cross().Across(size), axis.Along(size), axis.String()+": ")
		}
	})
	t.Run("size and along are inverse", func(t *testing.T) {
		for _, axis := range Axes() {
			built := axis.Size(1.0, 2.0)

			AssertNumber(t, axis.Along(built), 1.0, axis.String()+": ")
			AssertNumber(t, axis.Across(built), 2.0, axis.String()+": ")
		}
	})
	t.Run("direction round-trips through axis and sign", func(t *testing.T) {
		for _, axis := range Axes() {
			for _, positive := range []bool{true, false} {
				direction := axis.Direction(positive)

				assert.Equal(t, direction.Axis(), axis, direction.String()+": ")
				assert.Equal(t, direction.IsPositive(), positive, direction.String()+": ")
			}
		}
	})
	t.Run("project is the component of the vector", func(t *testing.T) {
		for _, axis := range Axes() {
			for _, vector := range vectorFixtures {
				AssertNumber(t, axis.Project(vector), axis.Along(Sz(vector.X, vector.Y)), fmt.Sprintf("%s → %s: ", axis, vector))
			}
		}
	})
}

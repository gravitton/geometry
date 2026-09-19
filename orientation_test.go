package geom

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

func TestOrientations(t *testing.T) {
	t.Run("lists both orientations in order", func(t *testing.T) {
		assert.Equal(t, Orientations(), [2]Orientation{FlatTop, PointyTop})
	})
	t.Run("returns a fresh array", func(t *testing.T) {
		orientations := Orientations()
		orientations[0] = PointyTop

		assert.Equal(t, Orientations()[0], FlatTop)
	})
}

func TestParseOrientation(t *testing.T) {
	t.Run("every name round-trips", func(t *testing.T) {
		for _, orientation := range Orientations() {
			parsed, err := ParseOrientation(orientation.String())

			assert.NoError(t, err, orientation.String()+": ")
			assert.Equal(t, parsed, orientation, orientation.String()+": ")
		}
	})
	t.Run("none", func(t *testing.T) {
		parsed, err := ParseOrientation("None")

		assert.NoError(t, err)
		assert.Equal(t, parsed, OrientationNone)
	})
	t.Run("unknown is an error", func(t *testing.T) {
		parsed, err := ParseOrientation("SharpTop")

		assert.Error(t, err)
		assert.Equal(t, parsed, OrientationNone)
	})
}

func TestOrientation_IsNone(t *testing.T) {
	t.Run("the two constants are an orientation", func(t *testing.T) {
		assert.False(t, FlatTop.IsNone())
		assert.False(t, PointyTop.IsNone())
	})
	t.Run("none and out of range", func(t *testing.T) {
		assert.True(t, OrientationNone.IsNone())
		assert.True(t, Orientation(99).IsNone())
	})
}

func TestOrientation_String(t *testing.T) {
	t.Run("named orientations", func(t *testing.T) {
		assert.Equal(t, FlatTop.String(), "FlatTop")
		assert.Equal(t, PointyTop.String(), "PointyTop")
	})
	t.Run("none and out of range", func(t *testing.T) {
		assert.Equal(t, OrientationNone.String(), "None")
		assert.Equal(t, Orientation(99).String(), "None")
	})
}

func TestOrientation_JSON(t *testing.T) {
	t.Run("wire format is the name", func(t *testing.T) {
		assert.JSON(t, PointyTop, `"PointyTop"`)
		assert.JSON(t, OrientationNone, `"None"`)
	})
	t.Run("round-trip", func(t *testing.T) {
		orientations := Orientations()
		for _, orientation := range append(orientations[:], OrientationNone) {
			data, err := json.Marshal(orientation)
			assert.NoError(t, err)

			var decoded Orientation
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, orientation, orientation.String()+": ")
		}
	})
	t.Run("an unknown name is an error", func(t *testing.T) {
		var decoded Orientation

		assert.Error(t, json.Unmarshal([]byte(`"SharpTop"`), &decoded))
	})
}

func TestOrientation_Properties(t *testing.T) {
	t.Run("every orientation marshals to the name it prints", func(t *testing.T) {
		orientations := Orientations()
		for _, orientation := range append(orientations[:], OrientationNone) {
			text, err := orientation.MarshalText()

			assert.NoError(t, err, orientation.String()+": ")
			assert.Equal(t, string(text), orientation.String(), orientation.String()+": ")
		}
	})
	t.Run("every orientation places a vertex or an edge at the top", func(t *testing.T) {
		for _, orientation := range Orientations() {
			for n := 3; n <= 8; n++ {
				polygon := RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), n, orientation)
				vertices := slices.Collect(polygon.Vertices())
				top := vertices[0]

				if orientation == PointyTop {
					AssertPoint(t, top, Pt(0.0, -10.0), orientation.String()+": ")
				} else {
					AssertVector(t, top.Midpoint(vertices[1]).Vector().Normalize(), Vec(0.0, -1.0), orientation.String()+": ")
				}
			}
		}
	})
}

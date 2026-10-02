package geom_test

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
)

func TestWindings(t *testing.T) {
	t.Run("lists both windings in order", func(t *testing.T) {
		assert.Equal(t, Windings(), [2]Winding{WindingClockwise, WindingCounterClockwise})
	})
	t.Run("returns a fresh array", func(t *testing.T) {
		windings := Windings()
		windings[0] = WindingCounterClockwise

		assert.Equal(t, Windings()[0], WindingClockwise)
	})
}

func TestParseWinding(t *testing.T) {
	t.Run("every name round-trips", func(t *testing.T) {
		for _, winding := range Windings() {
			parsed, err := ParseWinding(winding.String())

			assert.NoError(t, err, winding.String()+": ")
			assert.Equal(t, parsed, winding, winding.String()+": ")
		}
	})
	t.Run("none", func(t *testing.T) {
		parsed, err := ParseWinding("None")

		assert.NoError(t, err)
		assert.Equal(t, parsed, WindingNone)
	})
	t.Run("unknown is an error", func(t *testing.T) {
		parsed, err := ParseWinding("Widdershins")

		assert.Error(t, err)
		assert.Equal(t, parsed, WindingNone)
	})
}

func TestWinding_IsNone(t *testing.T) {
	t.Run("the two constants are a winding", func(t *testing.T) {
		assert.False(t, WindingClockwise.IsNone())
		assert.False(t, WindingCounterClockwise.IsNone())
	})
	t.Run("none and out of range", func(t *testing.T) {
		assert.True(t, WindingNone.IsNone())
		assert.True(t, Winding(99).IsNone())
	})
}

func TestWinding_String(t *testing.T) {
	t.Run("named windings", func(t *testing.T) {
		assert.Equal(t, WindingClockwise.String(), "Clockwise")
		assert.Equal(t, WindingCounterClockwise.String(), "CounterClockwise")
	})
	t.Run("none and out of range", func(t *testing.T) {
		assert.Equal(t, WindingNone.String(), "None")
		assert.Equal(t, Winding(99).String(), "None")
	})
}

func TestWinding_JSON(t *testing.T) {
	t.Run("wire format is the name", func(t *testing.T) {
		assert.JSON(t, WindingCounterClockwise, `"CounterClockwise"`)
		assert.JSON(t, WindingNone, `"None"`)
	})
	t.Run("round-trip", func(t *testing.T) {
		windings := Windings()
		for _, winding := range append(windings[:], WindingNone) {
			data, err := json.Marshal(winding)
			assert.NoError(t, err)

			var decoded Winding
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, winding, winding.String()+": ")
		}
	})
	t.Run("an unknown name is an error", func(t *testing.T) {
		var decoded Winding

		assert.Error(t, json.Unmarshal([]byte(`"Widdershins"`), &decoded))
	})
}

func TestWinding_Properties(t *testing.T) {
	t.Run("every winding marshals to the name it prints", func(t *testing.T) {
		windings := Windings()
		for _, winding := range append(windings[:], WindingNone) {
			text, err := winding.MarshalText()

			assert.NoError(t, err, winding.String()+": ")
			assert.Equal(t, string(text), winding.String(), winding.String()+": ")
		}
	})
	t.Run("every rectangle with an area winds clockwise", func(t *testing.T) {
		for _, r := range rectFixtures {
			if r.Area() == 0 {
				continue
			}

			assert.Equal(t, r.Polygon().Winding(), WindingClockwise, r.String()+": ")
		}
	})
}

package geom

import (
	"fmt"
)

// Winding is the sense in which an outline runs around the area it encloses, as drawn on a
// screen with Y pointing down, the sense Vector.Rotate turns in: the corners of a Rectangle
// run clockwise.
type Winding int

const (
	// WindingClockwise runs clockwise as drawn on a screen, the winding of Rectangle.Vertices.
	WindingClockwise Winding = iota
	// WindingCounterClockwise runs the other way round.
	WindingCounterClockwise

	// WindingNone is the winding of an outline that encloses no area, or whose lobes enclose
	// as much area one way as the other.
	WindingNone Winding = -1
)

// Windings lists both windings in order. It returns a fresh array, so a caller cannot alter
// the list.
func Windings() [2]Winding {
	return [2]Winding{WindingClockwise, WindingCounterClockwise}
}

// ParseWinding returns the winding with the given name, as String prints it, and an error for
// any other string. "None" parses to WindingNone.
func ParseWinding(name string) (Winding, error) {
	if name == WindingNone.String() {
		return WindingNone, nil
	}

	for _, winding := range Windings() {
		if winding.String() == name {
			return winding, nil
		}
	}

	return WindingNone, fmt.Errorf("geom: unknown winding %q", name)
}

// IsNone reports whether the winding is neither WindingClockwise nor WindingCounterClockwise.
// Like Orientation, a Winding outside the constants is not normalized, so every such value
// counts as WindingNone.
func (w Winding) IsNone() bool {
	return w != WindingClockwise && w != WindingCounterClockwise
}

// String returns the name of the winding constant.
func (w Winding) String() string {
	switch w {
	case WindingClockwise:
		return "Clockwise"
	case WindingCounterClockwise:
		return "CounterClockwise"
	default:
		return "None"
	}
}

// MarshalText implements encoding.TextMarshaler with the name String prints, so a winding is
// stored as "Clockwise" in JSON and as a map key rather than as its number.
func (w Winding) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, the inverse of MarshalText through
// ParseWinding.
func (w *Winding) UnmarshalText(text []byte) error {
	winding, err := ParseWinding(string(text))
	if err != nil {
		return err
	}

	*w = winding

	return nil
}

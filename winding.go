package geom

import (
	"fmt"
)

// Winding is the sense in which an outline runs around the area it encloses, as drawn on a screen
// with Y pointing down.
type Winding int

const (
	// WindingClockwise runs clockwise on screen, the winding of Rectangle.Vertices.
	WindingClockwise Winding = iota
	// WindingCounterClockwise runs counterclockwise on screen.
	WindingCounterClockwise

	// WindingNone is the winding of an outline that encloses no area, or as much one way as the
	// other.
	WindingNone Winding = -1
)

// Windings returns both windings in order.
func Windings() [2]Winding {
	return [2]Winding{WindingClockwise, WindingCounterClockwise}
}

// ParseWinding returns the winding with the name String prints, and an error for any other string.
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

// IsNone reports whether the winding is neither WindingClockwise nor WindingCounterClockwise, so
// every other value counts as WindingNone.
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

// MarshalText implements encoding.TextMarshaler with the name String prints.
func (w Winding) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler through ParseWinding.
func (w *Winding) UnmarshalText(text []byte) error {
	winding, err := ParseWinding(string(text))
	if err != nil {
		return err
	}

	*w = winding

	return nil
}

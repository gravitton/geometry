package geom

import (
	"fmt"
)

// Orientation is the rotational alignment of a regular polygon.
type Orientation int

const (
	// OrientationFlatTop places a flat edge at the top of the polygon.
	OrientationFlatTop Orientation = iota
	// OrientationPointyTop places a vertex at the top of the polygon.
	OrientationPointyTop

	// OrientationNone is the absence of an orientation, for which RegularPolygonOrientationPhase
	// panics.
	OrientationNone Orientation = -1
)

// Orientations returns both orientations in order.
func Orientations() [2]Orientation {
	return [2]Orientation{OrientationFlatTop, OrientationPointyTop}
}

// ParseOrientation returns the orientation with the name String prints, and an error for any other
// string.
func ParseOrientation(name string) (Orientation, error) {
	if name == OrientationNone.String() {
		return OrientationNone, nil
	}

	for _, orientation := range Orientations() {
		if orientation.String() == name {
			return orientation, nil
		}
	}

	return OrientationNone, fmt.Errorf("geom: unknown orientation %q", name)
}

// IsNone reports whether the orientation is neither OrientationFlatTop nor OrientationPointyTop, so
// every other value counts as OrientationNone.
func (o Orientation) IsNone() bool {
	return o != OrientationFlatTop && o != OrientationPointyTop
}

// String returns the name of the orientation constant.
func (o Orientation) String() string {
	switch o {
	case OrientationFlatTop:
		return "FlatTop"
	case OrientationPointyTop:
		return "PointyTop"
	default:
		return "None"
	}
}

// MarshalText implements encoding.TextMarshaler with the name String prints.
func (o Orientation) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler through ParseOrientation.
func (o *Orientation) UnmarshalText(text []byte) error {
	orientation, err := ParseOrientation(string(text))
	if err != nil {
		return err
	}

	*o = orientation

	return nil
}

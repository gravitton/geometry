package geom

import (
	"fmt"
)

// Axis is one of the two main coordinate axes.
type Axis int

const (
	// AxisHorizontal is the X axis.
	AxisHorizontal Axis = iota
	// AxisVertical is the Y axis.
	AxisVertical

	// AxisNone is the absence of an axis, with no extent on any method.
	AxisNone Axis = -1
)

// Axes returns both axes in order.
func Axes() [2]Axis {
	return [2]Axis{AxisHorizontal, AxisVertical}
}

// ParseAxis returns the axis with the name String prints, and an error for any other string.
func ParseAxis(name string) (Axis, error) {
	if name == AxisNone.String() {
		return AxisNone, nil
	}

	for _, axis := range Axes() {
		if axis.String() == name {
			return axis, nil
		}
	}

	return AxisNone, fmt.Errorf("geom: unknown axis %q", name)
}

// Perpendicular returns the perpendicular axis.
func (a Axis) Perpendicular() Axis {
	switch a {
	case AxisHorizontal:
		return AxisVertical
	case AxisVertical:
		return AxisHorizontal
	default:
		return AxisNone
	}
}

// Direction returns the direction along the axis with the sign: positive is DirectionRight on the
// horizontal axis and DirectionDown on the vertical one.
func (a Axis) Direction(positive bool) Direction {
	direction := a.direction()
	if !positive {
		return direction.Opposite()
	}

	return direction
}

// direction returns the direction of growing coordinates on the axis.
func (a Axis) direction() Direction {
	switch a {
	case AxisHorizontal:
		return DirectionRight
	case AxisVertical:
		return DirectionDown
	default:
		return DirectionNone
	}
}

// Along returns the extent of the size on the axis.
func (a Axis) Along[T Number](size Size[T]) T {
	switch a {
	case AxisHorizontal:
		return size.Width
	case AxisVertical:
		return size.Height
	default:
		return 0
	}
}

// Across returns the extent of the size on the perpendicular axis.
func (a Axis) Across[T Number](size Size[T]) T {
	return a.Perpendicular().Along(size)
}

// Project returns the signed component of the vector on the axis.
func (a Axis) Project[T Number](vector Vector[T]) T {
	switch a {
	case AxisHorizontal:
		return vector.X
	case AxisVertical:
		return vector.Y
	default:
		return 0
	}
}

// ScaleAlong returns the size scaled by the factor on the axis only, unchanged for AxisNone.
func (a Axis) ScaleAlong[T Number](size Size[T], factor float64) Size[T] {
	if a.IsNone() {
		return size
	}

	return a.Size(Multiply(a.Along(size), factor), a.Across(size))
}

// ScaleAcross returns the size scaled by the factor on the perpendicular axis only, unchanged for
// AxisNone.
func (a Axis) ScaleAcross[T Number](size Size[T], factor float64) Size[T] {
	return a.Perpendicular().ScaleAlong(size, factor)
}

// IsNone reports whether the axis is neither AxisHorizontal nor AxisVertical, so every other value
// counts as AxisNone.
func (a Axis) IsNone() bool {
	return a != AxisHorizontal && a != AxisVertical
}

// Vector returns the vector of along on the axis and across on the perpendicular one.
func (a Axis) Vector[T Number](along, across T) Vector[T] {
	switch a {
	case AxisHorizontal:
		return Vector[T]{along, across}
	case AxisVertical:
		return Vector[T]{across, along}
	default:
		return Vector[T]{}
	}
}

// Size returns the size of along on the axis and across on the perpendicular one, stored as given.
func (a Axis) Size[T Number](along, across T) Size[T] {
	switch a {
	case AxisHorizontal:
		return Size[T]{along, across}
	case AxisVertical:
		return Size[T]{across, along}
	default:
		return Size[T]{}
	}
}

// String returns the name of the axis constant.
func (a Axis) String() string {
	switch a {
	case AxisHorizontal:
		return "Horizontal"
	case AxisVertical:
		return "Vertical"
	default:
		return "None"
	}
}

// MarshalText implements encoding.TextMarshaler with the name String prints.
func (a Axis) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler through ParseAxis.
func (a *Axis) UnmarshalText(text []byte) error {
	axis, err := ParseAxis(string(text))
	if err != nil {
		return err
	}

	*a = axis

	return nil
}

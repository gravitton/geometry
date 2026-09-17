package geom

// Axis is one of the two main coordinate axes.
type Axis int

const (
	// AxisHorizontal is the X axis.
	AxisHorizontal Axis = iota
	// AxisVertical is the Y axis.
	AxisVertical

	// AxisNone is the absence of an axis. Every method treats it as having no extent.
	AxisNone Axis = -1
)

// Axes lists both axes in order. It returns a fresh array, so a caller cannot alter the list.
func Axes() [2]Axis {
	return [2]Axis{AxisHorizontal, AxisVertical}
}

// Cross returns the perpendicular axis.
func (a Axis) Cross() Axis {
	switch a {
	case AxisHorizontal:
		return AxisVertical
	case AxisVertical:
		return AxisHorizontal
	default:
		return AxisNone
	}
}

// Direction returns the axis direction with the given sign: positive is DirectionRight on the
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

// Along returns the extent of the given size on the main axis.
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

// Across returns the extent of the given size on the cross axis.
func (a Axis) Across[T Number](size Size[T]) T {
	return a.Cross().Along(size)
}

// Project returns the signed component of the given vector on the main axis.
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

// ScaleAlong creates a new Size scaled by the given factor on the main axis only.
// AxisNone has no axis to scale along and returns the size unchanged.
func (a Axis) ScaleAlong[T Number](size Size[T], factor float64) Size[T] {
	if a.IsNone() {
		return size
	}

	return a.Size(Multiply(a.Along(size), factor), a.Across(size))
}

// Vector creates a new Vector displaced by along on the main axis and across on the cross axis.
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

// Size creates a new Size measuring along on the main axis and across on the cross axis.
// The values are stored as given, like Sz; a negative one is not made absolute.
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

// IsNone reports whether the axis is neither AxisHorizontal nor AxisVertical. Unlike
// Direction, which wraps any value into the eight directions, an Axis outside the two
// constants is not normalized, so every such value counts as AxisNone.
func (a Axis) IsNone() bool {
	return a != AxisHorizontal && a != AxisVertical
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

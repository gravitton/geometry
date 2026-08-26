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

// Axes lists both axes in order.
var Axes = [2]Axis{AxisHorizontal, AxisVertical}

// Cross returns the perpendicular axis.
func (a Axis) Cross() Axis {
	var cross Axis

	switch a {
	case AxisHorizontal:
		cross = AxisVertical
	case AxisVertical:
		cross = AxisHorizontal
	default:
		cross = AxisNone
	}

	return cross
}

// Direction returns the axis direction with the given sign: positive is DirectionRight on the
// horizontal axis and DirectionDown on the vertical one.
func (a Axis) Direction(positive bool) Direction {
	var direction Direction

	switch a {
	case AxisHorizontal:
		direction = DirectionRight
	case AxisVertical:
		direction = DirectionDown
	default:
		direction = DirectionNone
	}

	if !positive {
		direction = direction.Opposite()
	}

	return direction
}

// Along returns the extent of the given size on the main axis.
func (a Axis) Along[T Number](size Size[T]) T {
	var along T

	switch a {
	case AxisHorizontal:
		along = size.Width
	case AxisVertical:
		along = size.Height
	}

	return along
}

// Across returns the extent of the given size on the cross axis.
func (a Axis) Across[T Number](size Size[T]) T {
	return a.Cross().Along(size)
}

// Project returns the component of the given vector on the main axis.
func (a Axis) Project[T Number](vector Vector[T]) T {
	return a.Along(vector.Size())
}

// ScaleAlong creates a new Size scaled by the given factor on the main axis only.
func (a Axis) ScaleAlong[T Number](size Size[T], factor float64) Size[T] {
	return a.Size(Multiply(a.Along(size), factor), a.Across(size))
}

// Vector creates a new Vector displaced by along on the main axis and across on the cross axis.
func (a Axis) Vector[T Number](along, across T) Vector[T] {
	var vector Vector[T]

	switch a {
	case AxisHorizontal:
		vector = Vector[T]{along, across}
	case AxisVertical:
		vector = Vector[T]{across, along}
	}

	return vector
}

// Size creates a new Size measuring along on the main axis and across on the cross axis.
func (a Axis) Size[T Number](along, across T) Size[T] {
	return a.Vector(along, across).Size()
}

// IsNone reports whether the axis is AxisNone.
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

package geom

import (
	"math"
)

// Direction is one of the eight neighbor directions on a square lattice.
type Direction int

const (
	DirectionRight Direction = iota
	DirectionUpRight
	DirectionUp
	DirectionUpLeft
	DirectionLeft
	DirectionDownLeft
	DirectionDown
	DirectionDownRight

	// DirectionNone is the absence of a direction.
	DirectionNone Direction = -1
)

// Direction aliases for lattice and map code.
const (
	East      = DirectionRight
	NorthEast = DirectionUpRight
	North     = DirectionUp
	NorthWest = DirectionUpLeft
	West      = DirectionLeft
	SouthWest = DirectionDownLeft
	South     = DirectionDown
	SouthEast = DirectionDownRight
)

// Direction aliases for the edges and corners of an axis-aligned rectangle.
const (
	Top         = DirectionUp
	Bottom      = DirectionDown
	Left        = DirectionLeft
	Right       = DirectionRight
	TopRight    = DirectionUpRight
	TopLeft     = DirectionUpLeft
	BottomLeft  = DirectionDownLeft
	BottomRight = DirectionDownRight
)

// Directions lists all eight directions ordered counterclockwise from DirectionRight.
var Directions = [8]Direction{DirectionRight, DirectionUpRight, DirectionUp, DirectionUpLeft, DirectionLeft, DirectionDownLeft, DirectionDown, DirectionDownRight}

// CardinalDirections lists the four cardinal directions ordered counterclockwise from DirectionRight.
var CardinalDirections = [4]Direction{DirectionRight, DirectionUp, DirectionLeft, DirectionDown}

// DiagonalDirections lists the four diagonal directions ordered counterclockwise from DirectionUpRight.
var DiagonalDirections = [4]Direction{DirectionUpRight, DirectionUpLeft, DirectionDownLeft, DirectionDownRight}

// directionOffsets lists the lattice step of each direction, indexed by direction.
var directionOffsets = [8]Vector[int]{
	{1, 0},   // Right
	{1, -1},  // UpRight
	{0, -1},  // Up
	{-1, -1}, // UpLeft
	{-1, 0},  // Left
	{-1, 1},  // DownLeft
	{0, 1},   // Down
	{1, 1},   // DownRight
}

// DirectionFromAngle returns the direction nearest to the given angle in radians.
func DirectionFromAngle(angle float64) Direction {
	return Mod(Direction(math.Round(-angle/(Pi/4))), 8)
}

// DirectionFromAxes returns the direction for the given axis inputs, canceling opposite
// directions (e.g. left+right = DirectionNone).
func DirectionFromAxes(up, down, left, right bool) Direction {
	if left && right {
		left, right = false, false
	}
	if up && down {
		up, down = false, false
	}

	switch {
	case up && left:
		return DirectionUpLeft
	case up && right:
		return DirectionUpRight
	case down && left:
		return DirectionDownLeft
	case down && right:
		return DirectionDownRight
	case up:
		return DirectionUp
	case down:
		return DirectionDown
	case left:
		return DirectionLeft
	case right:
		return DirectionRight
	default:
		return DirectionNone
	}
}

// Opposite returns the opposite direction, rotated 180°.
func (d Direction) Opposite() Direction {
	return d.Rotate(4)
}

// Rotate advances the direction by steps eighth-turns counterclockwise.
// On screens with Y pointing down, positive steps appear clockwise. Negative steps go the other way.
func (d Direction) Rotate(steps int) Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d+Direction(steps), 8)
}

// Axis returns the axis the direction runs on, or AxisNone for diagonals and DirectionNone.
func (d Direction) Axis() Axis {
	switch d.normalized() {
	case Right, Left:
		return AxisHorizontal
	case Top, Bottom:
		return AxisVertical
	default:
		return AxisNone
	}
}

// Offset creates a new non-normalized Vector with the lattice step of the direction, and the zero vector for DirectionNone.
func (d Direction) Offset[T Number]() Vector[T] {
	offset := d.offset()

	return Vector[T]{T(offset.X), T(offset.Y)}
}

// Unit creates a new normalized Vector in the direction, and the zero vector for DirectionNone.
func (d Direction) Unit[T Number]() Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Normalize()
}

// Vector creates a new Vector of the given length pointing in the direction, and the zero vector for DirectionNone.
func (d Direction) Vector[T Number](length T) Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Resize(float64(length))
}

// Angle returns the angle of the direction in radians.
func (d Direction) Angle() float64 {
	return d.Offset[float64]().Angle()
}

// IsNone reports whether the direction is DirectionNone.
func (d Direction) IsNone() bool {
	return d == DirectionNone
}

// normalized returns d wrapped into [DirectionRight, DirectionDownRight], preserving DirectionNone.
func (d Direction) normalized() Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d, 8)
}

// IsCardinal reports whether the direction is one of DirectionRight, DirectionUp, DirectionLeft, or DirectionDown.
func (d Direction) IsCardinal() bool {
	switch d.normalized() {
	case Right, Top, Left, Bottom:
		return true
	default:
		return false
	}
}

// IsDiagonal reports whether the direction is one of DirectionUpRight, DirectionUpLeft, DirectionDownLeft, or DirectionDownRight.
func (d Direction) IsDiagonal() bool {
	switch d.normalized() {
	case DirectionUpRight, DirectionUpLeft, DirectionDownLeft, DirectionDownRight:
		return true
	default:
		return false
	}
}

// IsPositive reports whether the direction runs toward growing coordinates on its axis:
// DirectionRight and DirectionDown, since Y grows downward.
// It is false for DirectionNone and for diagonals, which run on neither axis.
func (d Direction) IsPositive() bool {
	switch d.normalized() {
	case Right, Bottom:
		return true
	default:
		return false
	}
}

// String returns the name of the direction constant.
func (d Direction) String() string {
	switch d.normalized() {
	case DirectionRight:
		return "Right"
	case DirectionUpRight:
		return "UpRight"
	case DirectionUp:
		return "Up"
	case DirectionUpLeft:
		return "UpLeft"
	case DirectionLeft:
		return "Left"
	case DirectionDownLeft:
		return "DownLeft"
	case DirectionDown:
		return "Down"
	case DirectionDownRight:
		return "DownRight"
	default:
		return "None"
	}
}

// offset returns the lattice step of the direction, or a zero step for DirectionNone.
func (d Direction) offset() Vector[int] {
	if d.IsNone() {
		return Vector[int]{}
	}

	return directionOffsets[d.normalized()]
}

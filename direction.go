package geom

import (
	"math"
)

// Direction is one of the eight neighbor directions on a square lattice.
//
// Directions are numbered by increasing angle, matching Angle and Vector.Angle: counterclockwise
// in the standard math convention where Y grows upward, which appears clockwise as drawn on a
// screen with Y pointing down. A negative step or angle is therefore counterclockwise on screen.
type Direction int

const (
	DirectionRight Direction = iota
	DirectionDownRight
	DirectionDown
	DirectionDownLeft
	DirectionLeft
	DirectionUpLeft
	DirectionUp
	DirectionUpRight

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

// Directions lists all eight directions ordered by increasing angle from DirectionRight.
var Directions = [8]Direction{DirectionRight, DirectionDownRight, DirectionDown, DirectionDownLeft, DirectionLeft, DirectionUpLeft, DirectionUp, DirectionUpRight}

// CardinalDirections lists the four cardinal directions ordered by increasing angle from DirectionRight.
var CardinalDirections = [4]Direction{DirectionRight, DirectionDown, DirectionLeft, DirectionUp}

// DiagonalDirections lists the four diagonal directions ordered by increasing angle from DirectionDownRight.
var DiagonalDirections = [4]Direction{DirectionDownRight, DirectionDownLeft, DirectionUpLeft, DirectionUpRight}

// directionOffsets lists the lattice step of each direction, indexed by direction.
var directionOffsets = [8]Vector[int]{
	{1, 0},   // Right
	{1, 1},   // DownRight
	{0, 1},   // Down
	{-1, 1},  // DownLeft
	{-1, 0},  // Left
	{-1, -1}, // UpLeft
	{0, -1},  // Up
	{1, -1},  // UpRight
}

// DirectionFromAngle returns the direction nearest to the given angle in radians.
func DirectionFromAngle(angle float64) Direction {
	return Mod(Direction(math.Round(angle/(Pi/4))), 8)
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

// Rotate advances the direction by steps eighth-turns of increasing angle, the same sense as a
// positive Vector.Rotate angle: counterclockwise in math coordinates, clockwise as drawn on a
// screen with Y pointing down. Negative steps go the other way.
func (d Direction) Rotate(steps int) Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d+Direction(steps), 8)
}

// Axis returns the axis the direction runs on, or AxisNone for diagonals and DirectionNone.
func (d Direction) Axis() Axis {
	switch d.normalize() {
	case DirectionRight, DirectionLeft:
		return AxisHorizontal
	case DirectionUp, DirectionDown:
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

// Angle returns the angle of the direction in radians, measured in the standard math
// convention where Y grows upward — so DirectionUp is -Pi/2, not +Pi/2. Direction ordering
// follows this angle, so DirectionFromAngle and Angle round-trip for every direction.
func (d Direction) Angle() float64 {
	return d.Offset[float64]().Angle()
}

// IsNone reports whether the direction is DirectionNone.
func (d Direction) IsNone() bool {
	return d == DirectionNone
}

// IsCardinal reports whether the direction is one of DirectionRight, DirectionUp, DirectionLeft, or DirectionDown.
func (d Direction) IsCardinal() bool {
	switch d.normalize() {
	case DirectionRight, DirectionUp, DirectionLeft, DirectionDown:
		return true
	default:
		return false
	}
}

// IsDiagonal reports whether the direction is one of DirectionUpRight, DirectionUpLeft, DirectionDownLeft, or DirectionDownRight.
func (d Direction) IsDiagonal() bool {
	switch d.normalize() {
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
	switch d.normalize() {
	case DirectionRight, DirectionDown:
		return true
	default:
		return false
	}
}

// String returns the name of the direction constant.
func (d Direction) String() string {
	switch d.normalize() {
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

// normalize returns the direction wrapped into [DirectionRight, DirectionDownRight],
// preserving DirectionNone.
func (d Direction) normalize() Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d, 8)
}

// offset returns the lattice step of the direction, or a zero step for DirectionNone.
func (d Direction) offset() Vector[int] {
	if d.IsNone() {
		return Vector[int]{}
	}

	return directionOffsets[d.normalize()]
}

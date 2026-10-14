package geom

import (
	"fmt"
	"math"
)

// Direction is one of the eight neighbor directions on a square lattice.
//
// Directions are numbered by increasing angle, counterclockwise in math coordinates, which is
// clockwise as drawn on a screen with Y pointing down.
type Direction int

const (
	// DirectionRight points along +X, at angle 0.
	DirectionRight Direction = iota
	// DirectionDownRight points along +X and +Y, at angle π/4.
	DirectionDownRight
	// DirectionDown points along +Y, at angle π/2.
	DirectionDown
	// DirectionDownLeft points along -X and +Y, at angle 3π/4.
	DirectionDownLeft
	// DirectionLeft points along -X, at angle π.
	DirectionLeft
	// DirectionUpLeft points along -X and -Y, at angle -3π/4.
	DirectionUpLeft
	// DirectionUp points along -Y, at angle -π/2.
	DirectionUp
	// DirectionUpRight points along +X and -Y, at angle -π/4.
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

// Directions returns all eight directions by increasing angle from DirectionRight.
func Directions() [8]Direction {
	return [8]Direction{DirectionRight, DirectionDownRight, DirectionDown, DirectionDownLeft, DirectionLeft, DirectionUpLeft, DirectionUp, DirectionUpRight}
}

// CardinalDirections returns the four cardinal directions by increasing angle from DirectionRight.
func CardinalDirections() [4]Direction {
	return [4]Direction{DirectionRight, DirectionDown, DirectionLeft, DirectionUp}
}

// DiagonalDirections returns the four diagonal directions by increasing angle from
// DirectionDownRight.
func DiagonalDirections() [4]Direction {
	return [4]Direction{DirectionDownRight, DirectionDownLeft, DirectionUpLeft, DirectionUpRight}
}

// DirectionOffsets returns the Offset of each direction in the order of Directions.
func DirectionOffsets[T Number]() [8]Vector[T] {
	var offsets [8]Vector[T]
	for i, direction := range Directions() {
		offsets[i] = direction.Offset[T]()
	}

	return offsets
}

// CardinalOffsets returns the Offset of each direction in the order of CardinalDirections.
func CardinalOffsets[T Number]() [4]Vector[T] {
	var offsets [4]Vector[T]
	for i, direction := range CardinalDirections() {
		offsets[i] = direction.Offset[T]()
	}

	return offsets
}

// DiagonalOffsets returns the Offset of each direction in the order of DiagonalDirections.
func DiagonalOffsets[T Number]() [4]Vector[T] {
	var offsets [4]Vector[T]
	for i, direction := range DiagonalDirections() {
		offsets[i] = direction.Offset[T]()
	}

	return offsets
}

// DirectionFromAngle returns the direction nearest to the angle in radians, or DirectionNone for a
// non-finite angle.
func DirectionFromAngle(angle float64) Direction {
	if !IsFinite(angle) {
		return DirectionNone
	}

	return Mod(Direction(math.Round(NormalizeAngle(angle)/(Pi/4))), 8)
}

// DirectionFromAxes returns the direction of the pressed axis inputs, opposite inputs canceling
// each other.
func DirectionFromAxes(left, right, up, down bool) Direction {
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

// ParseDirection returns the direction with the name String prints, and an error for any other
// string.
func ParseDirection(name string) (Direction, error) {
	if name == DirectionNone.String() {
		return DirectionNone, nil
	}

	for _, direction := range Directions() {
		if direction.String() == name {
			return direction, nil
		}
	}

	return DirectionNone, fmt.Errorf("geom: unknown direction %q", name)
}

// Angle returns the angle of the direction in radians, in math coordinates, so DirectionUp is -π/2,
// and NaN for DirectionNone. DirectionFromAngle inverts it for every direction, DirectionNone
// included.
func (d Direction) Angle() float64 {
	if d.IsNone() {
		return math.NaN()
	}

	return d.Offset[float64]().Angle()
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

// normalize returns the direction wrapped into the eight constants, DirectionNone kept.
func (d Direction) normalize() Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d, 8)
}

// Opposite returns the direction turned by half a turn.
func (d Direction) Opposite() Direction {
	return d.Turn(4)
}

// Turn returns the direction advanced by steps eighth-turns of increasing angle, the sense of a
// positive Vector.Rotate; negative steps go the other way.
func (d Direction) Turn(steps int) Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return Mod(d+Direction(steps), 8)
}

// IsNone reports whether the direction is DirectionNone.
func (d Direction) IsNone() bool {
	return d == DirectionNone
}

// IsCardinal reports whether the direction runs along an axis.
func (d Direction) IsCardinal() bool {
	switch d.normalize() {
	case DirectionRight, DirectionUp, DirectionLeft, DirectionDown:
		return true
	default:
		return false
	}
}

// IsDiagonal reports whether the direction runs between the axes.
func (d Direction) IsDiagonal() bool {
	switch d.normalize() {
	case DirectionUpRight, DirectionUpLeft, DirectionDownLeft, DirectionDownRight:
		return true
	default:
		return false
	}
}

// IsPositive reports whether the direction runs toward growing coordinates on its axis:
// DirectionRight and DirectionDown. It is false for diagonals and DirectionNone.
func (d Direction) IsPositive() bool {
	switch d.normalize() {
	case DirectionRight, DirectionDown:
		return true
	default:
		return false
	}
}

// Offset returns the lattice step of the direction, with components of -1, 0 or 1, and the zero
// vector for DirectionNone.
func (d Direction) Offset[T Number]() Vector[T] {
	offset := d.offset()

	return Vector[T]{T(offset.X), T(offset.Y)}
}

// Unit returns the unit vector in the direction, and the zero vector for DirectionNone. For integer
// T a diagonal collapses onto an axis; Offset keeps it.
func (d Direction) Unit[T Number]() Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Normalize()
}

// Vector returns the vector of the length in the direction, and the zero vector for DirectionNone.
// A negative length points the other way. For integer T a diagonal is rounded per component, so its
// length is approximate.
func (d Direction) Vector[T Number](length T) Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Resize(float64(length))
}

// offset returns the lattice step of the direction, or a zero step for DirectionNone.
func (d Direction) offset() Vector[int] {
	if d.IsNone() {
		return Vector[int]{}
	}

	return directionOffsets[d.normalize()]
}

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

// MarshalText implements encoding.TextMarshaler with the name String prints.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler through ParseDirection.
func (d *Direction) UnmarshalText(text []byte) error {
	direction, err := ParseDirection(string(text))
	if err != nil {
		return err
	}

	*d = direction

	return nil
}

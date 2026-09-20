package geom

import (
	"fmt"
	"math"
)

// Direction is one of the eight neighbor directions on a square lattice.
//
// Directions are numbered by increasing normalized angle, matching Angle and Vector.Angle: counterclockwise
// in the standard math convention where Y grows upward, which appears clockwise as drawn on a
// screen with Y pointing down. A negative step or angle is therefore counterclockwise on screen.
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

// Directions lists all eight directions ordered by increasing angle from DirectionRight.
// It returns a fresh array, so a caller cannot alter the list.
func Directions() [8]Direction {
	return [8]Direction{DirectionRight, DirectionDownRight, DirectionDown, DirectionDownLeft, DirectionLeft, DirectionUpLeft, DirectionUp, DirectionUpRight}
}

// CardinalDirections lists the four cardinal directions ordered by increasing angle from DirectionRight.
// It returns a fresh array, so a caller cannot alter the list.
func CardinalDirections() [4]Direction {
	return [4]Direction{DirectionRight, DirectionDown, DirectionLeft, DirectionUp}
}

// DiagonalDirections lists the four diagonal directions ordered by increasing angle from DirectionDownRight.
// It returns a fresh array, so a caller cannot alter the list.
func DiagonalDirections() [4]Direction {
	return [4]Direction{DirectionDownRight, DirectionDownLeft, DirectionUpLeft, DirectionUpRight}
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

// DirectionFromAngle returns the direction nearest to the given angle in radians,
// or DirectionNone for NaN and ±Inf.
func DirectionFromAngle(angle float64) Direction {
	if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return DirectionNone
	}

	return Mod(Direction(math.Round(NormalizeAngle(angle)/(Pi/4))), 8)
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

// ParseDirection returns the direction with the given name, as String prints it, and an error
// for any other string. "None" parses to DirectionNone.
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

// Opposite returns the opposite direction, rotated 180°.
func (d Direction) Opposite() Direction {
	return d.Turn(4)
}

// Turn advances the direction by steps eighth-turns of increasing angle, the same sense as a
// positive Vector.Rotate angle: counterclockwise in math coordinates, clockwise as drawn on a
// screen with Y pointing down. Negative steps go the other way.
func (d Direction) Turn(steps int) Direction {
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
// For integer T, only the four axis-aligned vectors have length 1, so a diagonal collapses onto
// one of them; use Offset for the lattice step (±1,±1) that keeps the diagonal.
func (d Direction) Unit[T Number]() Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Normalize()
}

// Vector creates a new Vector of the given length pointing in the direction, and the zero vector for DirectionNone.
// A negative length points the other way, as Vector.Resize gives it, so the result runs along
// the opposite direction with the absolute length.
// For integer T, a diagonal has both components rounded, so its actual length is only
// approximate. The diagonal is kept, unlike Unit, which snaps to an axis: a diagonal of unit
// length is the lattice step of the direction, not the axis vector Unit gives.
func (d Direction) Vector[T Number](length T) Vector[T] {
	if d.IsNone() {
		return Vector[T]{}
	}

	return d.Offset[T]().Resize(float64(length))
}

// Angle returns the angle of the direction in radians, measured in the standard math
// convention where Y grows upward — so DirectionUp is -Pi/2, not +Pi/2, and NaN for
// DirectionNone, which has no angle. Direction ordering follows this angle normalized to
// [0, 2π), so DirectionFromAngle and Angle round-trip for every direction, DirectionNone included.
func (d Direction) Angle() float64 {
	if d.IsNone() {
		return math.NaN()
	}

	return d.Offset[float64]().Angle()
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

// MarshalText implements encoding.TextMarshaler with the name String prints, so a direction
// is stored as "UpRight" in JSON and as a map key rather than as its number.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, the inverse of MarshalText through ParseDirection.
func (d *Direction) UnmarshalText(text []byte) error {
	direction, err := ParseDirection(string(text))
	if err != nil {
		return err
	}

	*d = direction

	return nil
}

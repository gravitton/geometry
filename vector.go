package geom

import (
	"fmt"
)

// Vector is a 2D Vector.
type Vector[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Vec is shorthand for Vector[T]{x, y}.
func Vec[T Number](x, y T) Vector[T] {
	return Vector[T]{x, y}
}

// Add creates a new Vector by adding the given vector to the current vector.
func (p Vector[T]) Add(other Vector[T]) Vector[T] {
	return Vector[T]{p.X + other.X, p.Y + other.Y}
}

// Sub creates a new Vector by subtracting the given vector from the current vector.
func (p Vector[T]) Sub(other Vector[T]) Vector[T] {
	return Vector[T]{p.X - other.X, p.Y - other.Y}
}

// Equal checks for equal X and Y values with given vector.
func (p Vector[T]) Equal(other Vector[T]) bool {
	return Equal(p.X, other.X) && Equal(p.Y, other.Y)
}

// Zero checks if X and Y values are 0.
func (p Vector[T]) Zero() bool {
	// return Equal(p.X, 0) && Equal(p.Y, 0)
	return p.Equal(Vector[T]{})
}

// String returns a string representing the vector.
func (p Vector[T]) String() string {
	return fmt.Sprintf("⟨%s,%s⟩", ToString(p.X), ToString(p.Y))
}

// ZeroVector creates a new Vector with zero values (0,0).
func ZeroVector[T Number]() Vector[T] {
	return Vector[T]{}
}

// IdentityVector creates a new Vector with identity values (+1,+1).
func IdentityVector[T Number]() Vector[T] {
	return Vector[T]{1, 1}
}

// UpVector creates a new Vector with up (+y) direction (0,+1)
func UpVector[T Number]() Vector[T] {
	return Vector[T]{0, 1}
}

// UpVector creates a new Vector with down (-y) direction (0,-1)
func DownVector[T Number]() Vector[T] {
	return Vector[T]{0, -1}
}

// UpVector creates a new Vector with right (+x) direction (+1,0)
func RightVector[T Number]() Vector[T] {
	return Vector[T]{1, 0}
}

// UpVector creates a new Vector with left (-x) direction (-1,0)
func LeftVector[T Number]() Vector[T] {
	return Vector[T]{-1, 0}
}

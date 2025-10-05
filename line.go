package geom

import (
	"fmt"
)

// Line is a 2D line.
type Line[T Number] struct {
	Start Point[T] `json:"a"`
	End   Point[T] `json:"b"`
}

// Ln is shorthand for Line{start, end}.
func Ln[T Number](start, end Point[T]) Line[T] {
	return Line[T]{start, end}
}

// Translate creates a new Line translated by the given vector.
func (l Line[T]) Translate(vector Vector[T]) Line[T] {
	return Line[T]{l.Start.Add(vector), l.End.Add(vector)}
}

// MoveTo creates a new Line with the start point moved to point and same length and direction.
func (l Line[T]) MoveTo(point Point[T]) Line[T] {
	return Line[T]{point, l.End.Add(point.Subtract(l.Start))}
}

// Reversed creates a new Line with the start and end points swapped.
func (l Line[T]) Reversed() Line[T] {
	return Line[T]{l.End, l.Start}
}

// Midpoint returns the midpoint of the line.
func (l Line[T]) Midpoint() Point[T] {
	return l.Start.Midpoint(l.End)
}

// Direction returns the direction vector of the line.
func (l Line[T]) Direction() Vector[T] {
	return l.End.Subtract(l.Start)
}

// Length returns the length of the line.
func (l Line[T]) Length() float64 {
	return l.Direction().Length()
}

// Bounds returns the axis-aligned bounding rectangle.
func (l Line[T]) Bounds() Rectangle[T] {
	return Rectangle[T]{Center: l.Midpoint(), Size: Sz(l.Direction().Abs().XY())}
}

// Equal checks if the start and end points of the lines are equal.
func (l Line[T]) Equal(other Line[T]) bool {
	return l.Start.Equal(other.Start) && l.End.Equal(other.End)
}

// IsZero checks if start and end points are zero.
func (l Line[T]) IsZero() bool {
	return l.Start.IsZero() && l.End.IsZero()
}

// String returns a string representation of the Line.
func (l Line[T]) String() string {
	return fmt.Sprintf("L(%s;%s)", l.Start.String(), l.End.String())
}

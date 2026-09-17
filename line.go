package geom

import (
	"fmt"
	"math"
)

// Line is a 2D line.
type Line[T Number] struct {
	Start Point[T] `json:"s"`
	End   Point[T] `json:"e"`
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

// Reverse creates a new Line with the start and end points swapped.
func (l Line[T]) Reverse() Line[T] {
	return Line[T]{l.End, l.Start}
}

// Midpoint returns the midpoint of the line.
func (l Line[T]) Midpoint() Point[T] {
	return l.Start.Midpoint(l.End)
}

// Vector returns the line as a vector, from start to end.
func (l Line[T]) Vector() Vector[T] {
	return l.End.Subtract(l.Start)
}

// Length returns the length of the line.
func (l Line[T]) Length() float64 {
	return l.Vector().Length()
}

// DistanceTo returns the distance from the given point to the nearest point of the segment.
// A point on the segment is at distance 0 exactly, for an integer T, since the perpendicular
// distance comes from a cross product that is exact in float64 rather than from a projection.
func (l Line[T]) DistanceTo(point Point[T]) float64 {
	direction, offset := l.Vector().Float(), point.Subtract(l.Start).Float()

	along := offset.Dot(direction)
	if along <= 0 {
		return offset.Length()
	}
	if along >= direction.LengthSquared() {
		return point.Subtract(l.End).Length()
	}

	return math.Abs(offset.Cross(direction)) / direction.Length()
}

// Contains reports whether the given point lies on the segment, within Epsilon of T, the same
// closed convention as Rectangle.Contains.
func (l Line[T]) Contains(point Point[T]) bool {
	return LessOrEqualDelta(l.DistanceTo(point), 0, Epsilon[T]())
}

// wedge returns Start × End in float64, the term the shoelace formula sums per edge.
// The two products are rounded separately, which keeps a fused multiply-add from turning the
// wedge of two equal points into a rounding error: a degenerate edge contributes exactly zero.
func (l Line[T]) wedge() float64 {
	start, end := l.Start.Float(), l.End.Float()

	return float64(start.X*end.Y) - float64(start.Y*end.X)
}

// crossesRay reports whether a ray cast from the point along +X crosses the segment, counting
// an endpoint on the ray only when it is the lower one, so a ray through a vertex is counted
// once by the two edges that share it. It is the step of the even-odd rule Polygon.Contains uses.
func (l Line[T]) crossesRay(point Point[T]) bool {
	start, end, p := l.Start.Float(), l.End.Float(), point.Float()

	if (start.Y > p.Y) == (end.Y > p.Y) {
		return false
	}

	return p.X < start.X+(p.Y-start.Y)*(end.X-start.X)/(end.Y-start.Y)
}

// Vertices returns the start and end points as a slice.
func (l Line[T]) Vertices() []Point[T] {
	return []Point[T]{l.Start, l.End}
}

// Bounds returns the axis-aligned bounding rectangle.
func (l Line[T]) Bounds() Rectangle[T] {
	minPoint := Point[T]{min(l.Start.X, l.End.X), min(l.Start.Y, l.End.Y)}

	return RectangleFromMin(minPoint, l.Vector().Size())
}

// Equal checks if the start and end points of the lines are equal.
func (l Line[T]) Equal(line Line[T]) bool {
	return l.Start.Equal(line.Start) && l.End.Equal(line.End)
}

// IsZero checks if start and end points are zero.
func (l Line[T]) IsZero() bool {
	return l.Start.IsZero() && l.End.IsZero()
}

// Int converts the line to a Line[int].
func (l Line[T]) Int() Line[int] {
	return Line[int]{l.Start.Int(), l.End.Int()}
}

// Float converts the line to a Line[float64].
func (l Line[T]) Float() Line[float64] {
	return Line[float64]{l.Start.Float(), l.End.Float()}
}

// String returns the line in the form of its constructor: Ln((x,y);(x,y)).
func (l Line[T]) String() string {
	return fmt.Sprintf("Ln(%s;%s)", l.Start.String(), l.End.String())
}

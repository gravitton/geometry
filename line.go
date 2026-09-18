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

// Transform creates a new Line by applying the given matrix to both points, like Point.Transform.
func (l Line[T]) Transform[M Float](matrix Matrix[M]) Line[T] {
	return Line[T]{l.Start.Transform(matrix), l.End.Transform(matrix)}
}

// Translate creates a new Line translated by the given vector.
func (l Line[T]) Translate(vector Vector[T]) Line[T] {
	return Line[T]{l.Start.Add(vector), l.End.Add(vector)}
}

// MoveTo creates a new Line with the start point moved to point and same length and direction.
func (l Line[T]) MoveTo(point Point[T]) Line[T] {
	return Line[T]{point, l.End.Add(point.Subtract(l.Start))}
}

// Rotate creates a new Line rotated by the given angle (in radians) about its midpoint, in the
// same sense as Vector.Rotate. For integer T the midpoint and both rotated points are rounded;
// only multiples of 90° keep the length exactly.
func (l Line[T]) Rotate(angle float64) Line[T] {
	pivot := l.Midpoint()

	return Line[T]{l.Start.RotateAround(pivot, angle), l.End.RotateAround(pivot, angle)}
}

// Reverse creates a new Line with the start and end points swapped.
func (l Line[T]) Reverse() Line[T] {
	return Line[T]{l.End, l.Start}
}

// Lerp returns the point at the fraction t of the way from Start to End, extrapolating along
// the line outside [0, 1] like Point.Lerp. Midpoint is Lerp(0.5).
func (l Line[T]) Lerp(t float64) Point[T] {
	return l.Start.Lerp(l.End, t)
}

// Midpoint returns the midpoint of the line, Lerp(0.5).
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

// Intersects reports whether the segments share a point, within Epsilon of T, the same closed
// convention as Contains: segments that touch at an endpoint or overlap collinearly intersect.
func (l Line[T]) Intersects(line Line[T]) bool {
	return LessOrEqualDelta(l.distanceToLine(line), 0, Epsilon[T]())
}

// Intersection returns the point where the segments cross, and false when they do not.
// Parallel segments have no single crossing point and return false even where they overlap,
// which Intersects still reports. Touching at an endpoint counts, within Epsilon of T, the
// same closed convention as Intersects. For integer T the crossing is rounded like every
// other result stored into T.
func (l Line[T]) Intersection(line Line[T]) (Point[T], bool) {
	a, b := l.Float(), line.Float()

	denominator := a.Vector().Cross(b.Vector())
	if denominator == 0 {
		return Point[T]{}, false
	}

	t := b.Start.Subtract(a.Start).Cross(b.Vector()) / denominator
	point := a.Start.Add(a.Vector().Multiply(Clamp(t, 0, 1)))

	if !LessOrEqualDelta(b.DistanceTo(point), 0, Epsilon[T]()) {
		return Point[T]{}, false
	}

	return Point[T]{Cast[T](point.X), Cast[T](point.Y)}, true
}

// IntersectsCircle reports whether the segment and the circle share a point: the point of the
// segment closest to the center lies within the radius. Touching shapes intersect, within
// Epsilon of T.
func (l Line[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.Radius >= 0 && LessOrEqualDelta(l.DistanceTo(circle.Center), float64(circle.Radius), Epsilon[T]())
}

// IntersectsRectangle reports whether the segment and the rectangle share a point: the start
// lies within the rectangle, or the segment crosses one of its edges. Touching shapes intersect,
// within Epsilon of T.
func (l Line[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	if rectangle.Contains(l.Start) {
		return true
	}

	for edge := range rectangle.edges() {
		if l.Intersects(edge) {
			return true
		}
	}

	return false
}

// IntersectsPolygon reports whether the segment and the polygon share a point, as
// Polygon.IntersectsLine does.
func (l Line[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsLine(l)
}

// distanceToLine returns the distance between the nearest points of the two segments: zero
// when they cross, and otherwise the smallest distance from an endpoint of one to the other.
func (l Line[T]) distanceToLine(line Line[T]) float64 {
	if l.crosses(line) {
		return 0
	}

	return min(
		l.DistanceTo(line.Start), l.DistanceTo(line.End),
		line.DistanceTo(l.Start), line.DistanceTo(l.End),
	)
}

// crosses reports whether the segments properly cross: each has its endpoints on opposite sides
// of the other. Touching and collinear segments do not cross and are left to the endpoint
// distances, which cover them within the tolerance of the caller.
func (l Line[T]) crosses(line Line[T]) bool {
	return l.separates(line) && line.separates(l)
}

// separates reports whether the endpoints of the given segment lie strictly on opposite sides
// of the line through this one.
func (l Line[T]) separates(line Line[T]) bool {
	direction := l.Vector().Float()
	start := direction.Cross(line.Start.Subtract(l.Start).Float())
	end := direction.Cross(line.End.Subtract(l.Start).Float())

	return (start > 0 && end < 0) || (start < 0 && end > 0)
}

// wedge returns Start × End in float64, the term the shoelace formula sums per edge. Cross is
// exact for parallel vectors, so a degenerate edge contributes exactly zero.
func (l Line[T]) wedge() float64 {
	return l.Start.Vector().Float().Cross(l.End.Vector().Float())
}

// crossesRay reports whether a ray cast from the point along +X crosses the segment, counting
// an endpoint on the ray only when it is the lower one, so a ray through a vertex is counted
// once by the two edges that share it. It is the step of the even-odd rule Polygon.Contains uses.
//
// The ray crosses when the point lies on the side of the segment facing -X: the left side of a
// segment running toward +Y, the right side of one running toward -Y. The side comes from the
// sign of a cross product, which needs no division by the segment's Y span.
func (l Line[T]) crossesRay(point Point[T]) bool {
	start, end, p := l.Start.Float(), l.End.Float(), point.Float()

	if (start.Y > p.Y) == (end.Y > p.Y) {
		return false
	}

	upward := end.Y > start.Y
	left := end.Subtract(start).Cross(p.Subtract(start)) > 0

	return left == upward
}

// Vertices returns the start and end points as a slice.
func (l Line[T]) Vertices() []Point[T] {
	return []Point[T]{l.Start, l.End}
}

// Bounds returns the axis-aligned bounding rectangle.
func (l Line[T]) Bounds() Rectangle[T] {
	a := Point[T]{min(l.Start.X, l.End.X), min(l.Start.Y, l.End.Y)}

	return RectangleFromMin(a, l.Vector().Size())
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

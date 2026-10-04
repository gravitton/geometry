package geom

import (
	"cmp"
	"fmt"
	"math"
	"strings"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Pt is shorthand for Point{x, y}.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// ZeroPoint creates a new Point with zero values (0,0).
func ZeroPoint[T Number]() Point[T] {
	return Point[T]{}
}

// ParsePoint parses a point string in the form "(x,y)", the form String prints, each coordinate a
// number Parse accepts. For integer T, only integer strings parse; a fractional value is an
// error, not a rounded point.
func ParsePoint[T Number](s string) (Point[T], error) {
	coordinates, opened := strings.CutPrefix(s, "(")
	coordinates, closed := strings.CutSuffix(coordinates, ")")
	before, after, separated := strings.Cut(coordinates, ",")
	if !opened || !closed || !separated {
		return Point[T]{}, fmt.Errorf("geom: invalid point format %q", s)
	}

	x, err := Parse[T](before)
	if err != nil {
		return Point[T]{}, fmt.Errorf("geom: invalid x value: %w", err)
	}
	y, err := Parse[T](after)
	if err != nil {
		return Point[T]{}, fmt.Errorf("geom: invalid y value: %w", err)
	}

	return Point[T]{x, y}, nil
}

// XY returns the point X, Y values in standard order.
func (p Point[T]) XY() (T, T) {
	return p.X, p.Y
}

// magnitude returns the larger absolute coordinate of the point, in float64: the size epsilonAt
// widens the tolerance by for a comparison that reads it. An infinite coordinate has no ulp to
// widen it by, and gives a NaN, as a NaN coordinate does, so a point at infinity lies on no
// boundary rather than within an infinite tolerance of every one.
func (p Point[T]) magnitude() float64 {
	magnitude := max(math.Abs(float64(p.X)), math.Abs(float64(p.Y)))

	return magnitude + (magnitude - magnitude)
}

// Add creates a new Point by adding the given vector to the current point.
func (p Point[T]) Add(vector Vector[T]) Point[T] {
	return Point[T]{p.X + vector.X, p.Y + vector.Y}
}

// AddXY creates a new Point by adding the given values to the current point.
func (p Point[T]) AddXY(deltaX, deltaY T) Point[T] {
	return Point[T]{p.X + deltaX, p.Y + deltaY}
}

// Subtract creates a new Vector from a given point to the current point.
func (p Point[T]) Subtract(point Point[T]) Vector[T] {
	return Vector[T]{p.X - point.X, p.Y - point.Y}
}

// Multiply creates a new Point by multiplying the given value to the current point.
func (p Point[T]) Multiply(factor float64) Point[T] {
	return Point[T]{Multiply(p.X, factor), Multiply(p.Y, factor)}
}

// MultiplyXY creates a new Point by multiplying the given values to the current point.
func (p Point[T]) MultiplyXY(factorX, factorY float64) Point[T] {
	return Point[T]{Multiply(p.X, factorX), Multiply(p.Y, factorY)}
}

// Divide creates a new Point by dividing the given value to the current point.
func (p Point[T]) Divide(factor float64) Point[T] {
	return Point[T]{Divide(p.X, factor), Divide(p.Y, factor)}
}

// DivideXY creates a new Point by dividing the given values to the current point.
func (p Point[T]) DivideXY(factorX, factorY float64) Point[T] {
	return Point[T]{Divide(p.X, factorX), Divide(p.Y, factorY)}
}

// Abs creates a new Point with absolute X and Y.
func (p Point[T]) Abs() Point[T] {
	return Point[T]{Abs(p.X), Abs(p.Y)}
}

// Round creates a new Point by rounding X, Y values to the nearest integer.
func (p Point[T]) Round() Point[T] {
	return Point[T]{Round(p.X), Round(p.Y)}
}

// Floor creates a new Point by rounding down X, Y values to the nearest integer.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{Floor(p.X), Floor(p.Y)}
}

// Ceil creates a new Point by rounding up X, Y values to the nearest integer.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{Ceil(p.X), Ceil(p.Y)}
}

// Lerp creates a new Point in linear interpolation towards given point.
func (p Point[T]) Lerp(point Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, point.X, t), Lerp(p.Y, point.Y, t)}
}

// Midpoint creates a new Point between current and given points.
func (p Point[T]) Midpoint(point Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, point.X), Midpoint(p.Y, point.Y)}
}

// Transform creates a new Point by applying the given matrix to the current point.
// The matrix is float-only, like an angle: convert an integer matrix with Matrix.Float first.
// For integer T, the float64 result of each component is rounded; rotations and non-integer scales lose precision.
func (p Point[T]) Transform[M Float](matrix Matrix[M]) Point[T] {
	x, y := float64(p.X), float64(p.Y)
	a, b, c := float64(matrix.A), float64(matrix.B), float64(matrix.C)
	d, e, f := float64(matrix.D), float64(matrix.E), float64(matrix.F)

	return Point[T]{Cast[T](float64(a*x) + float64(b*y) + c), Cast[T](float64(d*x) + float64(e*y) + f)}
}

// RotateAround creates a new Point rotated by the given angle (in radians) about the pivot, in
// the same sense as Vector.Rotate. For integer T the result is rounded; only multiples of 90°
// give exact results.
func (p Point[T]) RotateAround(pivot Point[T], angle float64) Point[T] {
	center := pivot.Float()

	return center.Add(p.Float().Subtract(center).Rotate(angle)).Cast[T]()
}

// AngleTo returns the angle in radians from the current point to the given point.
func (p Point[T]) AngleTo(point Point[T]) float64 {
	return point.Float().Subtract(p.Float()).Angle()
}

// Between reports whether the point lies within the box from corner a to corner b, boundary
// included within the tolerance. It is the extent check Polygon.Contains makes before walking
// the edges, and a must be the lesser corner on each axis: a box given the other way round
// contains nothing, since the corners are not reordered.
func (p Point[T]) Between(a, b Point[T]) bool {
	epsilon := epsilonAt[T](max(p.magnitude(), a.magnitude(), b.magnitude()))

	return LessOrEqualDelta(a.X, p.X, epsilon) && LessOrEqualDelta(p.X, b.X, epsilon) &&
		LessOrEqualDelta(a.Y, p.Y, epsilon) && LessOrEqualDelta(p.Y, b.Y, epsilon)
}

// DistanceTo returns the Euclidean distance from the current point to the given point.
func (p Point[T]) DistanceTo(point Point[T]) float64 {
	return point.Float().Subtract(p.Float()).Length()
}

// DistanceSquaredTo returns the squared Euclidean distance to the given point (faster than DistanceTo for comparisons).
// The differences are taken in float64, as DistanceTo takes them, so a float32 difference is not
// rounded before it is squared and the two agree; the result is stored back through Cast.
func (p Point[T]) DistanceSquaredTo(point Point[T]) T {
	return Cast[T](point.Float().Subtract(p.Float()).LengthSquared())
}

// ManhattanDistanceTo returns the Manhattan (taxicab) distance from the current point to the given point.
func (p Point[T]) ManhattanDistanceTo(point Point[T]) T {
	dx, dy := p.deltas(point)

	return Cast[T](dx + dy)
}

// ChebyshevDistanceTo returns the Chebyshev distance (chessboard distance) from the current point to the given point.
// It is the maximum of the absolute differences of the coordinates: max(|dx|, |dy|).
func (p Point[T]) ChebyshevDistanceTo(point Point[T]) T {
	dx, dy := p.deltas(point)

	return Cast[T](max(dx, dy))
}

// OctileDistanceTo returns the Octile distance from the current point to the given point.
// Used in grid-based pathfinding where cardinal moves cost 1 and diagonal moves cost √2.
// It is: max(dx, dy) + (√2 - 1) * min(dx, dy).
func (p Point[T]) OctileDistanceTo(point Point[T]) float64 {
	dx, dy := p.deltas(point)

	return max(dx, dy) + float64((Sqrt2-1)*min(dx, dy))
}

// deltas returns the absolute coordinate differences to the given point in float64.
func (p Point[T]) deltas(point Point[T]) (float64, float64) {
	return math.Abs(float64(point.X) - float64(p.X)), math.Abs(float64(point.Y) - float64(p.Y))
}

// coincides reports whether the points lie within the tolerance of each other, judged on the
// squared distance as a point on a boundary is, where Equal compares each coordinate on its own.
func (p Point[T]) coincides(point Point[T]) bool {
	return lessOrEqualSquared(p.Float().DistanceSquaredTo(point.Float()), 0, epsilonAt[T](max(p.magnitude(), point.magnitude())))
}

// Equal checks for equal X and Y values with given point.
func (p Point[T]) Equal(point Point[T]) bool {
	return Equal(p.X, point.X) && Equal(p.Y, point.Y)
}

// Compare returns -1, 0, or +1 as p sorts before, with, or after point, ordering by
// X and then by Y. It follows the [cmp.Compare] convention and applies no tolerance, unlike
// Equal, so two float points that Equal considers the same can still order apart. A NaN
// coordinate sorts before every other value and with itself, as [cmp.Compare] orders it, so
// points carrying one still sort into a total order rather than breaking the sort.
func (p Point[T]) Compare(point Point[T]) int {
	if c := cmp.Compare(p.X, point.X); c != 0 {
		return c
	}

	return cmp.Compare(p.Y, point.Y)
}

// IsZero checks if X and Y values are zero.
func (p Point[T]) IsZero() bool {
	return p.Equal(Point[T]{})
}

// Vector converts the point to a Vector.
func (p Point[T]) Vector() Vector[T] {
	return Vector[T](p)
}

// Cast converts the point to a Point of another number type, rounding as Cast does.
func (p Point[T]) Cast[R Number]() Point[R] {
	return Point[R]{Cast[R](float64(p.X)), Cast[R](float64(p.Y))}
}

// Int converts the point to a Point[int].
func (p Point[T]) Int() Point[int] {
	return Point[int]{Int(p.X), Int(p.Y)}
}

// Float converts the point to a Point[float64].
func (p Point[T]) Float() Point[float64] {
	return Point[float64]{float64(p.X), float64(p.Y)}
}

// String returns a string representing the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%s,%s)", String(p.X), String(p.Y))
}

// minMaxOf returns the minimum and maximum corner of the vertices, the corners of their
// Bounds: the pair the intersection tests reject shapes by before examining any edge. No
// vertices have no corners and return two zero points.
func minMaxOf[T Number](vertices []Point[T]) (Point[T], Point[T]) {
	if len(vertices) == 0 {
		return Point[T]{}, Point[T]{}
	}

	a, b := vertices[0], vertices[0]
	for _, vertex := range vertices[1:] {
		a = Point[T]{min(a.X, vertex.X), min(a.Y, vertex.Y)}
		b = Point[T]{max(b.X, vertex.X), max(b.Y, vertex.Y)}
	}

	return a, b
}

// overlaps reports whether the box from a1 to b1 and the box from a2 to b2 share a point,
// boundary included within the tolerance. It is the check Rectangle.IntersectsRectangle makes on its
// corners and the rejection every other intersection test makes before examining edges, and
// like Between it expects each a to be the lesser corner on each axis.
func overlaps[T Number](a1, b1, a2, b2 Point[T]) bool {
	epsilon := epsilonAt[T](max(a1.magnitude(), b1.magnitude(), a2.magnitude(), b2.magnitude()))

	return LessOrEqualDelta(a1.X, b2.X, epsilon) && LessOrEqualDelta(a2.X, b1.X, epsilon) &&
		LessOrEqualDelta(a1.Y, b2.Y, epsilon) && LessOrEqualDelta(a2.Y, b1.Y, epsilon)
}

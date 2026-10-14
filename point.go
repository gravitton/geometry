package geom

import (
	"cmp"
	"fmt"
	"math"
	"strings"
)

// Point is a position in the plane.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Pt is shorthand for Point{x, y}.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{x, y}
}

// ZeroPoint returns the origin.
func ZeroPoint[T Number]() Point[T] {
	return Point[T]{}
}

// ParsePoint returns the point in the form String prints, and an error for any other string. Each
// coordinate is a number Parse accepts; for integer T a fractional one is an error, not rounded.
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

// XY returns X and Y.
func (p Point[T]) XY() (T, T) {
	return p.X, p.Y
}

// magnitude returns the larger absolute coordinate, the magnitude epsilonAt widens the tolerance
// by. An infinite coordinate gives NaN, so a point at infinity lies on no boundary rather than
// within an infinite tolerance of every one.
func (p Point[T]) magnitude() float64 {
	magnitude := max(math.Abs(float64(p.X)), math.Abs(float64(p.Y)))

	return magnitude + (magnitude - magnitude)
}

// Add returns the point moved by the vector.
func (p Point[T]) Add(vector Vector[T]) Point[T] {
	return Point[T]{p.X + vector.X, p.Y + vector.Y}
}

// AddXY returns the point moved by the deltas.
func (p Point[T]) AddXY(deltaX, deltaY T) Point[T] {
	return Point[T]{p.X + deltaX, p.Y + deltaY}
}

// Subtract returns the vector from the given point to this one.
func (p Point[T]) Subtract(point Point[T]) Vector[T] {
	return Vector[T]{p.X - point.X, p.Y - point.Y}
}

// Multiply returns the point with both coordinates multiplied by the factor.
func (p Point[T]) Multiply(factor float64) Point[T] {
	return Point[T]{Multiply(p.X, factor), Multiply(p.Y, factor)}
}

// MultiplyXY returns the point with each coordinate multiplied by its own factor.
func (p Point[T]) MultiplyXY(factorX, factorY float64) Point[T] {
	return Point[T]{Multiply(p.X, factorX), Multiply(p.Y, factorY)}
}

// Divide returns the point with both coordinates divided by the factor. It panics for a zero
// factor.
func (p Point[T]) Divide(factor float64) Point[T] {
	return Point[T]{Divide(p.X, factor), Divide(p.Y, factor)}
}

// DivideXY returns the point with each coordinate divided by its own factor. It panics for a zero
// factor.
func (p Point[T]) DivideXY(factorX, factorY float64) Point[T] {
	return Point[T]{Divide(p.X, factorX), Divide(p.Y, factorY)}
}

// Abs returns the point with absolute coordinates.
func (p Point[T]) Abs() Point[T] {
	return Point[T]{Abs(p.X), Abs(p.Y)}
}

// Round returns the point with coordinates rounded to the nearest integer.
func (p Point[T]) Round() Point[T] {
	return Point[T]{Round(p.X), Round(p.Y)}
}

// Floor returns the point with coordinates rounded down.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{Floor(p.X), Floor(p.Y)}
}

// Ceil returns the point with coordinates rounded up.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{Ceil(p.X), Ceil(p.Y)}
}

// Lerp returns the point interpolated toward the given one, extrapolated outside [0, 1].
func (p Point[T]) Lerp(point Point[T], t float64) Point[T] {
	return Point[T]{Lerp(p.X, point.X, t), Lerp(p.Y, point.Y, t)}
}

// Midpoint returns the point halfway to the given one.
func (p Point[T]) Midpoint(point Point[T]) Point[T] {
	return Point[T]{Midpoint(p.X, point.X), Midpoint(p.Y, point.Y)}
}

// Transform returns the point mapped by the matrix. The matrix is float-only: convert an integer
// one with Matrix.Float first.
func (p Point[T]) Transform[M Float](matrix Matrix[M]) Point[T] {
	x, y := float64(p.X), float64(p.Y)
	a, b, c := float64(matrix.A), float64(matrix.B), float64(matrix.C)
	d, e, f := float64(matrix.D), float64(matrix.E), float64(matrix.F)

	return Point[T]{Cast[T](float64(a*x) + float64(b*y) + c), Cast[T](float64(d*x) + float64(e*y) + f)}
}

// RotateAround returns the point rotated by the angle in radians about the pivot, in the sense of
// Vector.Rotate.
func (p Point[T]) RotateAround(pivot Point[T], angle float64) Point[T] {
	center := pivot.Float()

	return center.Add(p.Float().Subtract(center).Rotate(angle)).Cast[T]()
}

// MoveToward returns the point moved toward the target by at most the step, never past it: where
// the distance left is at most the step, the result is the target exactly. A negative step moves
// away and is not clamped, and a point at the target stays there. For integer T a step under half a
// unit can leave the point in place.
func (p Point[T]) MoveToward(target Point[T], step float64) Point[T] {
	return p.Vector().MoveToward(target.Vector(), step).Point()
}

// AngleTo returns the angle in radians of the vector from this point to the given one.
func (p Point[T]) AngleTo(point Point[T]) float64 {
	return point.Float().Subtract(p.Float()).Angle()
}

// DistanceTo returns the distance to the given point.
func (p Point[T]) DistanceTo(point Point[T]) float64 {
	return point.Float().Subtract(p.Float()).Length()
}

// DistanceSquaredTo returns the square of DistanceTo.
func (p Point[T]) DistanceSquaredTo(point Point[T]) T {
	return Cast[T](point.Float().Subtract(p.Float()).LengthSquared())
}

// ManhattanDistanceTo returns the taxicab distance to the given point (|dx|+|dy|).
func (p Point[T]) ManhattanDistanceTo(point Point[T]) T {
	dx, dy := p.deltas(point)

	return Cast[T](dx + dy)
}

// ChebyshevDistanceTo returns the chessboard distance to the given point (max(|dx|, |dy|)).
func (p Point[T]) ChebyshevDistanceTo(point Point[T]) T {
	dx, dy := p.deltas(point)

	return Cast[T](max(dx, dy))
}

// OctileDistanceTo returns the grid distance to the given point where a cardinal move costs one and
// a diagonal one √2 (max + (√2−1)·min of |dx| and |dy|).
func (p Point[T]) OctileDistanceTo(point Point[T]) float64 {
	dx, dy := p.deltas(point)

	return max(dx, dy) + float64((Sqrt2-1)*min(dx, dy))
}

// Between reports whether the point lies in the axis-aligned box from a to b. The corners are not
// reordered, so a box with a greater than b on either axis contains nothing.
func (p Point[T]) Between(a, b Point[T]) bool {
	epsilon := epsilonAt[T](max(p.magnitude(), a.magnitude(), b.magnitude()))

	return LessOrEqualDelta(a.X, p.X, epsilon) && LessOrEqualDelta(p.X, b.X, epsilon) &&
		LessOrEqualDelta(a.Y, p.Y, epsilon) && LessOrEqualDelta(p.Y, b.Y, epsilon)
}

// deltas returns the absolute coordinate differences to the given point.
func (p Point[T]) deltas(point Point[T]) (float64, float64) {
	return math.Abs(float64(point.X) - float64(p.X)), math.Abs(float64(point.Y) - float64(p.Y))
}

// coincides reports whether the points lie within the tolerance of each other on the squared
// distance, as a point on a boundary is judged, where Equal compares each coordinate on its own.
func (p Point[T]) coincides(point Point[T]) bool {
	return lessOrEqualSquared(p.Float().DistanceSquaredTo(point.Float()), 0, epsilonAt[T](max(p.magnitude(), point.magnitude())))
}

// Equal reports whether the points have equal coordinates.
func (p Point[T]) Equal(point Point[T]) bool {
	return Equal(p.X, point.X) && Equal(p.Y, point.Y)
}

// Compare returns -1, 0 or +1 as the point sorts before, with or after the given one, by X and then
// by Y, as cmp.Compare does. It applies no tolerance, so points Equal holds for can order apart,
// and a NaN coordinate sorts first.
func (p Point[T]) Compare(point Point[T]) int {
	if c := cmp.Compare(p.X, point.X); c != 0 {
		return c
	}

	return cmp.Compare(p.Y, point.Y)
}

// IsZero reports whether the point is the origin.
func (p Point[T]) IsZero() bool {
	return p.Equal(Point[T]{})
}

// IsFinite reports whether both coordinates are neither NaN nor infinite, which Cast and Int need
// for an integer type. An integer point always is.
func (p Point[T]) IsFinite() bool {
	return IsFinite(p.X) && IsFinite(p.Y)
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

// String returns the point as its coordinates, with no type name: (x,y).
func (p Point[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = p.appendString(buf)

	return string(buf)
}

// appendString appends the point as String prints it to buf and returns the extended buffer.
func (p Point[T]) appendString(buf []byte) []byte {
	buf = append(buf, '(')
	buf = appendString(buf, p.X)
	buf = append(buf, ',')
	buf = appendString(buf, p.Y)
	buf = append(buf, ')')

	return buf
}

// minMaxOf returns the corners of the Bounds of the vertices, and two zero points for none.
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

// supportOf returns the vertex farthest along the direction: the first of several that reach alike,
// and the zero point for no vertices or a NaN direction.
func supportOf[T Number](vertices []Point[T], direction Vector[float64]) Point[float64] {
	var support Point[float64]

	reach := math.Inf(-1)
	for _, vertex := range vertices {
		point := vertex.Float()
		if along := point.Vector().Dot(direction); along > reach {
			support, reach = point, along
		}
	}

	return support
}

// overlaps reports whether the boxes from a1 to b1 and from a2 to b2 share a point, within the
// tolerance, each a being the lesser corner on both axes.
func overlaps[T Number](a1, b1, a2, b2 Point[T]) bool {
	epsilon := epsilonAt[T](max(a1.magnitude(), b1.magnitude(), a2.magnitude(), b2.magnitude()))

	return LessOrEqualDelta(a1.X, b2.X, epsilon) && LessOrEqualDelta(a2.X, b1.X, epsilon) &&
		LessOrEqualDelta(a1.Y, b2.Y, epsilon) && LessOrEqualDelta(a2.Y, b1.Y, epsilon)
}

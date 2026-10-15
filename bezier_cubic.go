package geom

import (
	"math"
	"slices"
)

// flattenDepth bounds how many times a curve is halved while flattened, so a curve no tolerance can
// be met on, NaN included, still ends.
const flattenDepth = 10

// bezierDepth bounds how many times a curve is halved while a segment is cast at it, far more than
// a piece needs to be straight within the tolerance of its magnitude.
const bezierDepth = 32

// CubicBezier is a curve from Start to End, leaving Start toward Control1 and reaching End from
// Control2. It lies within the hull of the four points and passes through neither control point in
// general.
//
// The curve is no Shape: Polyline flattens it to the path that measures it, and segments, rays and
// lines are cast at it as a wall. It has no center, so it moves by Translate and turns and scales
// by Transform alone.
type CubicBezier[T Number] struct {
	Start    Point[T] `json:"s"`
	Control1 Point[T] `json:"c"`
	Control2 Point[T] `json:"d"`
	End      Point[T] `json:"e"`
}

// CubBez is shorthand for CubicBezier{start, control1, control2, end}.
func CubBez[T Number](start, control1, control2, end Point[T]) CubicBezier[T] {
	return CubicBezier[T]{start, control1, control2, end}
}

// Bounds returns the axis-aligned bounding box of the curve, not of its control points. For integer
// T it is rounded outward, so it still holds the curve.
func (b CubicBezier[T]) Bounds() Box[T] {
	points := [6]Point[float64]{b.Start.Float(), b.End.Float()}

	n := 2
	for _, t := range b.extrema() {
		if t > 0 && t < 1 {
			_, _, points[n] = b.casteljau(t)
			n++
		}
	}

	a1, b1 := minMaxOf(points[:n])

	return Box[T]{Point[T]{-castUp[T](-a1.X), -castUp[T](-a1.Y)}, Point[T]{castUp[T](b1.X), castUp[T](b1.Y)}}
}

// extrema returns the parameters at which the curve turns back along X and along Y, two per axis,
// NaN where there is none and unclamped to the curve.
func (b CubicBezier[T]) extrema() [4]float64 {
	first, second, third := b.Control1.Float().Subtract(b.Start.Float()), b.Control2.Float().Subtract(b.Control1.Float()), b.End.Float().Subtract(b.Control2.Float())

	x1, x2 := b.roots(first.X, second.X, third.X)
	y1, y2 := b.roots(first.Y, second.Y, third.Y)

	return [4]float64{x1, x2, y1, y2}
}

// roots returns the two roots of one coordinate of the derivative, given the differences of the
// points in it, NaN or infinite where there are fewer. Dividing the sum that does not cancel keeps
// the one root of a nearly linear derivative, as a raised quadratic has.
func (CubicBezier[T]) roots(first, second, third float64) (float64, float64) {
	a, b := first-float64(2*second)+third, 2*(second-first)
	sum := -(b + math.Copysign(math.Sqrt(float64(b*b)-float64(4*a*first)), b)) / 2

	return sum / a, first / sum
}

// casteljau returns the points of de Casteljau's construction at the parameter: the three between
// neighbouring points, the two between those, and the point of the curve between the two.
func (b CubicBezier[T]) casteljau(t float64) ([3]Point[float64], [2]Point[float64], Point[float64]) {
	start, control1, control2, end := b.Start.Float(), b.Control1.Float(), b.Control2.Float(), b.End.Float()

	first := [3]Point[float64]{start.Lerp(control1, t), control1.Lerp(control2, t), control2.Lerp(end, t)}
	second := [2]Point[float64]{first[0].Lerp(first[1], t), first[1].Lerp(first[2], t)}

	return first, second, second[0].Lerp(second[1], t)
}

// Translate returns the curve moved by the vector.
func (b CubicBezier[T]) Translate(vector Vector[T]) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Add(vector), b.Control1.Add(vector), b.Control2.Add(vector), b.End.Add(vector)}
}

// Transform returns the curve mapped by the matrix, exactly, by mapping its four points.
func (b CubicBezier[T]) Transform[M Float](matrix Matrix[M]) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Transform(matrix), b.Control1.Transform(matrix), b.Control2.Transform(matrix), b.End.Transform(matrix)}
}

// Lerp returns the curve interpolated toward the given one, its four points together, and
// extrapolated outside [0, 1].
func (b CubicBezier[T]) Lerp(curve CubicBezier[T], t float64) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Lerp(curve.Start, t), b.Control1.Lerp(curve.Control1, t), b.Control2.Lerp(curve.Control2, t), b.End.Lerp(curve.End, t)}
}

// Reverse returns the curve run the other way, from End to Start.
func (b CubicBezier[T]) Reverse() CubicBezier[T] {
	return CubicBezier[T]{b.End, b.Control2, b.Control1, b.Start}
}

// PointAt returns the point of the curve at the parameter, Start at zero and End at one,
// extrapolated outside [0, 1]. The parameter is no fraction of the length: equal steps of it bunch
// up near the control points.
func (b CubicBezier[T]) PointAt(t float64) Point[T] {
	_, _, point := b.casteljau(t)

	return point.Cast[T]()
}

// TangentAt returns the derivative of the curve at the parameter, zero at a cusp and at an end that
// its control point lies on.
func (b CubicBezier[T]) TangentAt(t float64) Vector[float64] {
	_, second, _ := b.casteljau(t)

	return second[1].Subtract(second[0]).Multiply(3)
}

// Split returns the curve cut at the parameter, from Start to PointAt(t) and from there to End. For
// integer T the new points are rounded, so the halves only approximate the curve.
func (b CubicBezier[T]) Split(t float64) (CubicBezier[T], CubicBezier[T]) {
	first, second, point := b.casteljau(t)
	middle := point.Cast[T]()

	return CubicBezier[T]{b.Start, first[0].Cast[T](), second[0].Cast[T](), middle}, CubicBezier[T]{middle, second[1].Cast[T](), first[2].Cast[T](), b.End}
}

// IntersectionSegment returns the points where the segment crosses the curve, from Start to End.
func (b CubicBezier[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionCubicBezier(b)
}

// IntersectionRay returns the points where the ray crosses the curve, from its Origin on.
func (b CubicBezier[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionCubicBezier(b)
}

// IntersectionLine returns the points where the line crosses the curve, in order along its
// Vector.
func (b CubicBezier[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionCubicBezier(b)
}

// appendIntersectionSegment appends the points where the segment crosses the curve to dst, rounded
// once and a point two pieces share given once; rays and lines call it on their reach and span.
func (b CubicBezier[T]) appendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	cast := bezierCast{segment: segment.Float(), curve: b.Float()}
	cast.step(cast.curve, 0, 1, bezierDepth)

	held := len(dst)
	for _, hit := range cast.hits[:cast.count] {
		if rounded := hit.Point.Cast[T](); !(len(dst) > held && rounded.Equal(dst[len(dst)-1])) {
			dst = append(dst, rounded)
		}
	}

	return dst
}

// Equal reports whether the curves have equal points in the same places, so a curve and its Reverse
// differ.
func (b CubicBezier[T]) Equal(curve CubicBezier[T]) bool {
	return b.Start.Equal(curve.Start) && b.Control1.Equal(curve.Control1) && b.Control2.Equal(curve.Control2) && b.End.Equal(curve.End)
}

// IsZero reports whether the curve is the zero value.
func (b CubicBezier[T]) IsZero() bool {
	return b.Equal(CubicBezier[T]{})
}

// Polyline returns the curve flattened to a path within the tolerance of it, from Start to End,
// with more points where it bends, one edge where it is straight and one point where it has no
// length. A negative tolerance is taken absolute, and the curve is halved at most ten times, so the
// path has at most 1024 edges.
func (b CubicBezier[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{b.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of Polyline(tolerance) to dst and returns the extended slice.
func (b CubicBezier[T]) AppendPolyline(dst []Point[T], tolerance float64) []Point[T] {
	return b.Float().appendPolyline(dst, tolerance)
}

// Cast converts the curve to a CubicBezier of another number type, rounding as Cast does.
func (b CubicBezier[T]) Cast[R Number]() CubicBezier[R] {
	return CubicBezier[R]{b.Start.Cast[R](), b.Control1.Cast[R](), b.Control2.Cast[R](), b.End.Cast[R]()}
}

// Int converts the curve to a CubicBezier[int].
func (b CubicBezier[T]) Int() CubicBezier[int] {
	return CubicBezier[int]{b.Start.Int(), b.Control1.Int(), b.Control2.Int(), b.End.Int()}
}

// Float converts the curve to a CubicBezier[float64].
func (b CubicBezier[T]) Float() CubicBezier[float64] {
	return CubicBezier[float64]{b.Start.Float(), b.Control1.Float(), b.Control2.Float(), b.End.Float()}
}

// isFlatWithin reports whether both control points lie within the tolerance of the chord, which
// holds the curve within it too, since the curve lies in their hull.
func (b CubicBezier[T]) isFlatWithin(tolerance float64) bool {
	chord := Segment[T]{b.Start, b.End}

	return max(chord.distanceSquaredTo(b.Control1), chord.distanceSquaredTo(b.Control2)) <= tolerance*tolerance
}

// pieces returns how many chords flatten gives at most, so appendPolyline grows dst once.
func (b CubicBezier[T]) pieces(tolerance float64, depth int) int {
	if depth == 0 || b.isFlatWithin(tolerance) {
		return 1
	}

	first, second := b.Split(0.5)

	return first.pieces(tolerance, depth-1) + second.pieces(tolerance, depth-1)
}

// appendPolyline appends Start and the ends flatten gives to dst, rounded into its type. It is
// called on the float64 curve, where Split is exact, so an integer curve is halved as it is, not as
// its rounded halves.
func (b CubicBezier[T]) appendPolyline[R Number](dst []Point[R], tolerance float64) []Point[R] {
	tolerance = math.Abs(tolerance)
	dst = append(slices.Grow(dst, b.pieces(tolerance, flattenDepth)+1), b.Start.Cast[R]())

	return b.flatten(dst, tolerance, flattenDepth)
}

// flatten appends the ends of the chords replacing the curve, less Start: its End where it is flat
// or the depth is spent, the ends of its halves otherwise. An end rounding onto the one before is
// left out, so no chord has zero length.
func (b CubicBezier[T]) flatten[R Number](dst []Point[R], tolerance float64, depth int) []Point[R] {
	if depth == 0 || b.isFlatWithin(tolerance) {
		if end := b.End.Cast[R](); end != dst[len(dst)-1] {
			dst = append(dst, end)
		}

		return dst
	}

	first, second := b.Split(0.5)

	return second.flatten(first.flatten[R](dst, tolerance, depth-1), tolerance, depth-1)
}

// String returns the curve in the form of its constructor: CubBez((x,y);(x,y);(x,y);(x,y)).
func (b CubicBezier[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "CubBez("...)
	buf = b.Start.appendString(buf)
	buf = append(buf, ';')
	buf = b.Control1.appendString(buf)
	buf = append(buf, ';')
	buf = b.Control2.appendString(buf)
	buf = append(buf, ';')
	buf = b.End.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

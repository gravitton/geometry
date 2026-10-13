package geom

import (
	"math"
	"slices"
)

// flattenDepth is how many times a curve is halved at most while it is flattened, so a curve
// no tolerance can be met on, one with a NaN among its points, still ends: a Polyline of a
// curve has at most two to that power of edges.
const flattenDepth = 10

// bezierDepth is how many times a curve is halved at most while a segment is cast at it: far
// more than a piece needs to become straight to within Delta, so only a curve of no sensible
// size is crossed as chords that are not yet straight.
const bezierDepth = 32

// CubicBezier is a 2D cubic Bezier curve: it runs from Start to End, leaving Start toward
// Control1 and reaching End from Control2, the curve of a path in a vector drawing and of an
// eased motion. It touches neither control point in general and lies within the hull of the
// four.
//
// The curve is no Shape, and has no Length or Nearest of its own, neither of which has a closed
// form: Polyline flattens it to a path within a tolerance, and that path measures it. A
// segment, a ray or a line is cast at it as a wall, by their Intersection and Hit methods
// naming it. It has no center either, so it moves by Translate and turns and scales by
// Transform alone, which maps the curve by mapping its four points.
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

// Bounds returns the axis-aligned bounding box of the curve itself, not of its control points:
// the box around its two ends and the points where it turns back along an axis, the roots of
// its derivative. For integer T the box is rounded outward, so it still holds the curve.
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

// extrema returns the fractions at which the curve turns back along X and along Y, two for
// each axis: the roots of the derivative, a quadratic in each coordinate, wherever they fall,
// and NaN where there is none. Bounds keeps those within the curve.
func (b CubicBezier[T]) extrema() [4]float64 {
	first, second, third := b.Control1.Float().Subtract(b.Start.Float()), b.Control2.Float().Subtract(b.Control1.Float()), b.End.Float().Subtract(b.Control2.Float())

	x1, x2 := b.roots(first.X, second.X, third.X)
	y1, y2 := b.roots(first.Y, second.Y, third.Y)

	return [4]float64{x1, x2, y1, y2}
}

// roots returns the two fractions at which one coordinate of the derivative is zero, given the
// three differences between the points in that coordinate: the roots of the quadratic they are
// the Bernstein coefficients of, and a NaN or an infinity where it has fewer. The larger root
// in magnitude is divided out of the sum that does not cancel and the other out of their
// product, so a derivative that is linear, or a rounding error from it as that of a quadratic
// curve raised to a cubic is, still gives its one root. It reads no field of the curve, so the
// receiver is unnamed.
func (CubicBezier[T]) roots(first, second, third float64) (float64, float64) {
	a, b := first-float64(2*second)+third, 2*(second-first)
	sum := -(b + math.Copysign(math.Sqrt(float64(b*b)-float64(4*a*first)), b)) / 2

	return sum / a, first / sum
}

// casteljau returns the points de Casteljau's construction passes through at the fraction t,
// in float64: the three between neighbouring points of the curve, the two between those, and
// the point between the two, which is the point of the curve. PointAt, TangentAt and Split are
// all read off it.
func (b CubicBezier[T]) casteljau(t float64) ([3]Point[float64], [2]Point[float64], Point[float64]) {
	start, control1, control2, end := b.Start.Float(), b.Control1.Float(), b.Control2.Float(), b.End.Float()

	first := [3]Point[float64]{start.Lerp(control1, t), control1.Lerp(control2, t), control2.Lerp(end, t)}
	second := [2]Point[float64]{first[0].Lerp(first[1], t), first[1].Lerp(first[2], t)}

	return first, second, second[0].Lerp(second[1], t)
}

// Translate creates a new CubicBezier translated by the given vector.
func (b CubicBezier[T]) Translate(vector Vector[T]) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Add(vector), b.Control1.Add(vector), b.Control2.Add(vector), b.End.Add(vector)}
}

// Transform creates a new CubicBezier by applying the given matrix to its four points, like
// Point.Transform, which is the curve itself transformed: a Bezier curve keeps its form under
// every affine map.
func (b CubicBezier[T]) Transform[M Float](matrix Matrix[M]) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Transform(matrix), b.Control1.Transform(matrix), b.Control2.Transform(matrix), b.End.Transform(matrix)}
}

// Lerp creates a new CubicBezier in linear interpolation towards the given curve, moving the
// four points together like Point.Lerp, and extrapolating outside [0, 1]: each point of the
// curve moves straight towards the point of the other at the same parameter.
func (b CubicBezier[T]) Lerp(curve CubicBezier[T], t float64) CubicBezier[T] {
	return CubicBezier[T]{b.Start.Lerp(curve.Start, t), b.Control1.Lerp(curve.Control1, t), b.Control2.Lerp(curve.Control2, t), b.End.Lerp(curve.End, t)}
}

// Reverse creates a new CubicBezier over the same points run the other way, from End to Start.
func (b CubicBezier[T]) Reverse() CubicBezier[T] {
	return CubicBezier[T]{b.End, b.Control2, b.Control1, b.Start}
}

// PointAt returns the point of the curve at the parameter t, Start at zero and End at one,
// extrapolating the curve outside [0, 1]. The parameter is not a fraction of the length: equal
// steps of it are closer together where the control points are. For integer T the point is
// rounded once.
func (b CubicBezier[T]) PointAt(t float64) Point[T] {
	_, _, point := b.casteljau(t)

	return point.Cast[T]()
}

// TangentAt returns the derivative of the curve at the parameter t, the direction it runs in
// there with the speed of the parameter as its length, in float64 whatever T is. It is the
// zero vector where the curve has no direction: at an end whose control point lies on it, and
// at a cusp.
func (b CubicBezier[T]) TangentAt(t float64) Vector[float64] {
	_, second, _ := b.casteljau(t)

	return second[1].Subtract(second[0]).Multiply(3)
}

// Split returns the two curves the curve is cut into at the parameter t, the first from Start
// to PointAt(t) and the second from there to End, which together trace the same points. For
// integer T the new points are rounded, so the halves only approximate the curve.
func (b CubicBezier[T]) Split(t float64) (CubicBezier[T], CubicBezier[T]) {
	first, second, point := b.casteljau(t)
	middle := point.Cast[T]()

	return CubicBezier[T]{b.Start, first[0].Cast[T](), second[0].Cast[T](), middle}, CubicBezier[T]{middle, second[1].Cast[T](), first[2].Cast[T](), b.End}
}

// IntersectionSegment returns the points where the segment crosses the curve, from its Start to
// its End, three at most. They are closed on, not solved for: the curve is halved where the
// segment passes its control points, until each piece is straight to within Delta and crossed
// as its chord, so a point is within that of the curve, a segment touching the curve within the
// tolerance meets it, and no tolerance is asked of the caller as Polyline asks one. A segment
// running along a straight curve shares a stretch with it and returns an end of the stretch or
// none, and one that touches a bend within the tolerance over a stretch, without crossing it,
// can meet it at two points of that stretch where there is one touch. For integer T the points
// are rounded once.
func (b CubicBezier[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return b.AppendIntersectionSegment(nil, segment)
}

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns
// the extended slice, so a caller reusing dst allocates nothing once it has room for three.
func (b CubicBezier[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
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

// IntersectionRay returns the points where the ray crosses the curve, from its Origin on, as
// IntersectionSegment finds them on the reach of the ray past the curve.
func (b CubicBezier[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return b.AppendIntersectionRay(nil, ray)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice.
func (b CubicBezier[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	bounds := b.Bounds()

	return b.AppendIntersectionSegment(dst, ray.reach(bounds.Min, bounds.Max))
}

// IntersectionLine returns the points where the line crosses the curve, in order along its
// Direction, as IntersectionSegment finds them on the span of the line past the curve.
func (b CubicBezier[T]) IntersectionLine(line Line[T]) []Point[T] {
	return b.AppendIntersectionLine(nil, line)
}

// AppendIntersectionLine appends the points IntersectionLine returns to dst and returns the
// extended slice.
func (b CubicBezier[T]) AppendIntersectionLine(dst []Point[T], line Line[T]) []Point[T] {
	bounds := b.Bounds()

	return b.AppendIntersectionSegment(dst, line.span(bounds.Min, bounds.Max))
}

// Equal checks if the four points of the curves are equal, in their places, so a curve and its
// Reverse are not.
func (b CubicBezier[T]) Equal(curve CubicBezier[T]) bool {
	return b.Start.Equal(curve.Start) && b.Control1.Equal(curve.Control1) && b.Control2.Equal(curve.Control2) && b.End.Equal(curve.End)
}

// IsZero checks if all four points are zero.
func (b CubicBezier[T]) IsZero() bool {
	return b.Equal(CubicBezier[T]{})
}

// Polyline converts the curve to a Polyline within the given tolerance of it, from Start to
// End: the curve is halved until each piece lies within the tolerance of its chord, its control
// points no farther from it than that, so the path has more points where the curve bends and
// two where it is straight. It is what measures and tests the curve: its Length, Nearest and
// DistanceTo are the curve's to within the tolerance. A negative tolerance is taken absolute,
// and a curve is halved ten times at most, so a tolerance of zero gives a path of a thousand
// edges. For integer T the points are rounded and one rounding onto the point before it is
// left out.
func (b CubicBezier[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{b.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of the path Polyline returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them. The
// curve is walked twice, once to count the points, so dst grows once.
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

// isFlatWithin reports whether the curve lies within the tolerance of the chord from Start to
// End: both control points do, by the distance Segment.DistanceSquaredTo measures, and the
// curve lies within their hull.
func (b CubicBezier[T]) isFlatWithin(tolerance float64) bool {
	chord := Segment[T]{b.Start, b.End}

	return max(chord.distanceSquaredTo(b.Control1), chord.distanceSquaredTo(b.Control2)) <= tolerance*tolerance
}

// pieces returns how many chords flatten replaces the curve by at most: one where it is flat
// within the tolerance or the depth is spent, and those of its two halves otherwise.
func (b CubicBezier[T]) pieces(tolerance float64, depth int) int {
	if depth == 0 || b.isFlatWithin(tolerance) {
		return 1
	}

	first, second := b.Split(0.5)

	return first.pieces(tolerance, depth-1) + second.pieces(tolerance, depth-1)
}

// appendPolyline appends the points of the path the curve is flattened to within the
// tolerance, rounded into the type of dst, which grows once to hold them: Start, and the ends
// flatten gives. It is called on the curve in float64, where Split is exact, so the curve of an
// integer T is halved as it is, not as its rounded halves.
func (b CubicBezier[T]) appendPolyline[R Number](dst []Point[R], tolerance float64) []Point[R] {
	tolerance = math.Abs(tolerance)
	dst = append(slices.Grow(dst, b.pieces(tolerance, flattenDepth)+1), b.Start.Cast[R]())

	return b.flatten(dst, tolerance, flattenDepth)
}

// flatten appends the ends of the chords the curve is replaced by, less its Start, which the
// caller has appended: its End where it is flat within the tolerance or the depth is spent,
// and the ends of its two halves otherwise. An end rounding onto the point before it is left
// out, so no chord has zero length. Each end is rounded once into the type of dst.
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

// bezierCast gathers where a segment crosses a Bezier curve, in order from its Start: the
// hits of the segment at the chords of the pieces the curve is halved into, each with the
// perpendicular of the curve at the parameter the point lies at as the normal, a point two
// pieces share held once. A cubic is crossed three times at most, and the room beyond that is
// for a crossing found by both of the pieces it lies between.
type bezierCast struct {
	segment Segment[float64]
	curve   CubicBezier[float64]
	hits    [6]Hit[float64]
	count   int
}

// step casts the segment at the piece of the curve between the two parameters: nothing where
// it does not reach the box around the control points, which holds the piece, the crossing of
// its chord where the piece is straight to within Delta or the depth is spent, and its two
// halves otherwise.
func (c *bezierCast) step(piece CubicBezier[float64], from, to float64, depth int) {
	a, b := minMaxOf([]Point[float64]{piece.Start, piece.Control1, piece.Control2, piece.End})
	if !(Box[float64]{a, b}).IntersectsSegment(c.segment) {
		return
	}

	if depth == 0 || piece.isFlatWithin(Delta) {
		chord := Segment[float64]{piece.Start, piece.End}
		if point, ok := c.segment.IntersectionSegment(chord); ok {
			t := from + float64((to-from)*Clamp(chord.fraction(point), 0, 1))

			c.add(c.segment.hitAt(point, c.curve.TangentAt(t).Normal()))
		}

		return
	}

	first, second := piece.Split(0.5)
	middle := (from + to) / 2

	c.step(first, from, middle, depth-1)
	c.step(second, middle, to, depth-1)
}

// add holds the hit in its place by Fraction, unless its point is already held. Where there is
// no room left the farthest gives way, which copy drops off the end.
func (c *bezierCast) add(hit Hit[float64]) {
	at := c.count
	for i, held := range c.hits[:c.count] {
		if held.Point.Equal(hit.Point) {
			return
		}

		if hit.Fraction < held.Fraction && at == c.count {
			at = i
		}
	}

	at = min(at, len(c.hits)-1)

	copy(c.hits[at+1:], c.hits[at:c.count])
	c.hits[at] = hit
	c.count = min(c.count+1, len(c.hits))
}

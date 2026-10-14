package geom

// QuadraticBezier is a 2D quadratic Bezier curve: it runs from Start to End, leaving Start
// toward Control and reaching End from it, the curve of a rounded corner and of a font outline.
// It does not touch the control point in general and lies within the triangle of the three.
//
// The curve is no Shape, and has no Length or Nearest of its own: Polyline flattens it to a path
// within a tolerance, and that path measures it. A segment, a ray or a line is cast at it as a
// wall, by their Intersection and Hit methods naming it. It has no center either, so it moves
// by Translate and turns and scales by Transform alone, which maps the curve by mapping its
// three points.
type QuadraticBezier[T Number] struct {
	Start   Point[T] `json:"s"`
	Control Point[T] `json:"c"`
	End     Point[T] `json:"e"`
}

// QuadBez is shorthand for QuadraticBezier{start, control, end}.
func QuadBez[T Number](start, control, end Point[T]) QuadraticBezier[T] {
	return QuadraticBezier[T]{start, control, end}
}

// Bounds returns the axis-aligned bounding box of the curve itself, not of its control point:
// the box around its two ends and the points where it turns back along an axis, the roots of
// its derivative. For integer T the box is rounded outward, so it still holds the curve.
func (b QuadraticBezier[T]) Bounds() Box[T] {
	first, second := b.Control.Float().Subtract(b.Start.Float()), b.End.Float().Subtract(b.Control.Float())
	points := [4]Point[float64]{b.Start.Float(), b.End.Float()}

	n := 2
	for _, t := range [2]float64{first.X / (first.X - second.X), first.Y / (first.Y - second.Y)} {
		if t > 0 && t < 1 {
			_, points[n] = b.casteljau(t)
			n++
		}
	}

	a1, b1 := minMaxOf(points[:n])

	return Box[T]{Point[T]{-castUp[T](-a1.X), -castUp[T](-a1.Y)}, Point[T]{castUp[T](b1.X), castUp[T](b1.Y)}}
}

// casteljau returns the points de Casteljau's construction passes through at the fraction t,
// in float64: the two between neighbouring points of the curve, and the point between those,
// which is the point of the curve. PointAt, TangentAt and Split are all read off it.
func (b QuadraticBezier[T]) casteljau(t float64) ([2]Point[float64], Point[float64]) {
	start, control, end := b.Start.Float(), b.Control.Float(), b.End.Float()
	first := [2]Point[float64]{start.Lerp(control, t), control.Lerp(end, t)}

	return first, first[0].Lerp(first[1], t)
}

// Translate creates a new QuadraticBezier translated by the given vector.
func (b QuadraticBezier[T]) Translate(vector Vector[T]) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Add(vector), b.Control.Add(vector), b.End.Add(vector)}
}

// Transform creates a new QuadraticBezier by applying the given matrix to its three points,
// like Point.Transform, which is the curve itself transformed: a Bezier curve keeps its form
// under every affine map.
func (b QuadraticBezier[T]) Transform[M Float](matrix Matrix[M]) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Transform(matrix), b.Control.Transform(matrix), b.End.Transform(matrix)}
}

// Lerp creates a new QuadraticBezier in linear interpolation towards the given curve, moving
// the three points together like Point.Lerp, and extrapolating outside [0, 1]: each point of
// the curve moves straight towards the point of the other at the same parameter.
func (b QuadraticBezier[T]) Lerp(curve QuadraticBezier[T], t float64) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Lerp(curve.Start, t), b.Control.Lerp(curve.Control, t), b.End.Lerp(curve.End, t)}
}

// Reverse creates a new QuadraticBezier over the same points run the other way, from End to
// Start.
func (b QuadraticBezier[T]) Reverse() QuadraticBezier[T] {
	return QuadraticBezier[T]{b.End, b.Control, b.Start}
}

// PointAt returns the point of the curve at the parameter t, Start at zero and End at one,
// extrapolating the curve outside [0, 1]. The parameter is not a fraction of the length: equal
// steps of it are closer together where the control point is. For integer T the point is
// rounded once.
func (b QuadraticBezier[T]) PointAt(t float64) Point[T] {
	_, point := b.casteljau(t)

	return point.Cast[T]()
}

// TangentAt returns the derivative of the curve at the parameter t, the direction it runs in
// there with the speed of the parameter as its length, in float64 whatever T is. It is the
// zero vector where the curve has no direction, at an end the control point lies on.
func (b QuadraticBezier[T]) TangentAt(t float64) Vector[float64] {
	first, _ := b.casteljau(t)

	return first[1].Subtract(first[0]).Multiply(2)
}

// Split returns the two curves the curve is cut into at the parameter t, the first from Start
// to PointAt(t) and the second from there to End, which together trace the same points. For
// integer T the new points are rounded, so the halves only approximate the curve.
func (b QuadraticBezier[T]) Split(t float64) (QuadraticBezier[T], QuadraticBezier[T]) {
	first, point := b.casteljau(t)
	middle := point.Cast[T]()

	return QuadraticBezier[T]{b.Start, first[0].Cast[T](), middle}, QuadraticBezier[T]{middle, first[1].Cast[T](), b.End}
}

// IntersectionSegment returns the points where the segment crosses the curve, from its Start to
// its End, two at most, as CubicBezier.IntersectionSegment finds them on the curve raised to a
// cubic.
//
// The points are a Pair, which allocates nothing.
func (b QuadraticBezier[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	cast := bezierCast{segment: segment.Float(), curve: b.Float().CubicBezier()}
	cast.step(cast.curve, 0, 1, bezierDepth)

	var points Pair[Point[T]]
	for _, hit := range cast.hits[:cast.count] {
		if rounded := hit.Point.Cast[T](); points.Count < 2 && !(points.Count > 0 && rounded.Equal(points.Items[points.Count-1])) {
			points.Items[points.Count] = rounded
			points.Count++
		}
	}

	return points
}

// IntersectionRay returns the points where the ray crosses the curve, from its Origin on, as
// IntersectionSegment finds them on the reach of the ray past the curve.
func (b QuadraticBezier[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	bounds := b.Bounds()

	return b.IntersectionSegment(ray.reach(bounds.Min, bounds.Max))
}

// IntersectionLine returns the points where the line crosses the curve, in order along its
// Direction, as IntersectionSegment finds them on the span of the line past the curve.
func (b QuadraticBezier[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	bounds := b.Bounds()

	return b.IntersectionSegment(line.span(bounds.Min, bounds.Max))
}

// Equal checks if the three points of the curves are equal, in their places, so a curve and
// its Reverse are not.
func (b QuadraticBezier[T]) Equal(curve QuadraticBezier[T]) bool {
	return b.Start.Equal(curve.Start) && b.Control.Equal(curve.Control) && b.End.Equal(curve.End)
}

// IsZero checks if all three points are zero.
func (b QuadraticBezier[T]) IsZero() bool {
	return b.Equal(QuadraticBezier[T]{})
}

// CubicBezier converts the curve to the CubicBezier over the same points, its two control points two
// thirds of the way from each end to Control. For integer T they are rounded, so the cubic only
// approximates the curve.
func (b QuadraticBezier[T]) CubicBezier() CubicBezier[T] {
	return CubicBezier[T]{b.Start, b.Start.Lerp(b.Control, 2.0/3), b.End.Lerp(b.Control, 2.0/3), b.End}
}

// Polyline converts the curve to a Polyline within the given tolerance of it, from Start to
// End, as CubicBezier.Polyline flattens the cubic over the same points: with more points where
// the curve bends and two where it is straight. It is what measures and tests the curve: its
// Length, Nearest and DistanceTo are the curve's to within the tolerance. A negative tolerance
// is taken absolute, and for integer T the points are rounded.
func (b QuadraticBezier[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{b.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of the path Polyline returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them.
func (b QuadraticBezier[T]) AppendPolyline(dst []Point[T], tolerance float64) []Point[T] {
	return b.Float().CubicBezier().appendPolyline(dst, tolerance)
}

// Cast converts the curve to a QuadraticBezier of another number type, rounding as Cast does.
func (b QuadraticBezier[T]) Cast[R Number]() QuadraticBezier[R] {
	return QuadraticBezier[R]{b.Start.Cast[R](), b.Control.Cast[R](), b.End.Cast[R]()}
}

// Int converts the curve to a QuadraticBezier[int].
func (b QuadraticBezier[T]) Int() QuadraticBezier[int] {
	return QuadraticBezier[int]{b.Start.Int(), b.Control.Int(), b.End.Int()}
}

// Float converts the curve to a QuadraticBezier[float64].
func (b QuadraticBezier[T]) Float() QuadraticBezier[float64] {
	return QuadraticBezier[float64]{b.Start.Float(), b.Control.Float(), b.End.Float()}
}

// String returns the curve in the form of its constructor: QuadBez((x,y);(x,y);(x,y)).
func (b QuadraticBezier[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "QuadBez("...)
	buf = b.Start.appendString(buf)
	buf = append(buf, ';')
	buf = b.Control.appendString(buf)
	buf = append(buf, ';')
	buf = b.End.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

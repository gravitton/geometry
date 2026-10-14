package geom

// QuadraticBezier is a curve from Start to End, leaving Start toward Control and reaching End from
// it. It lies within the triangle of the three points and does not pass through Control in general.
//
// The curve is no Shape: Polyline flattens it to the path that measures it, and segments, rays and
// lines are cast at it as a wall. It has no center, so it moves by Translate and turns and scales
// by Transform alone.
type QuadraticBezier[T Number] struct {
	Start   Point[T] `json:"s"`
	Control Point[T] `json:"c"`
	End     Point[T] `json:"e"`
}

// QuadBez is shorthand for QuadraticBezier{start, control, end}.
func QuadBez[T Number](start, control, end Point[T]) QuadraticBezier[T] {
	return QuadraticBezier[T]{start, control, end}
}

// Bounds returns the axis-aligned bounding box of the curve, not of its control point. For integer
// T it is rounded outward, so it still holds the curve.
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

// casteljau returns the points of de Casteljau's construction at the parameter: the two between
// neighbouring points, and the point of the curve between those.
func (b QuadraticBezier[T]) casteljau(t float64) ([2]Point[float64], Point[float64]) {
	start, control, end := b.Start.Float(), b.Control.Float(), b.End.Float()
	first := [2]Point[float64]{start.Lerp(control, t), control.Lerp(end, t)}

	return first, first[0].Lerp(first[1], t)
}

// Translate returns the curve moved by the vector.
func (b QuadraticBezier[T]) Translate(vector Vector[T]) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Add(vector), b.Control.Add(vector), b.End.Add(vector)}
}

// Transform returns the curve mapped by the matrix, exactly, by mapping its three points.
func (b QuadraticBezier[T]) Transform[M Float](matrix Matrix[M]) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Transform(matrix), b.Control.Transform(matrix), b.End.Transform(matrix)}
}

// Lerp returns the curve interpolated toward the given one, its three points together, and
// extrapolated outside [0, 1].
func (b QuadraticBezier[T]) Lerp(curve QuadraticBezier[T], t float64) QuadraticBezier[T] {
	return QuadraticBezier[T]{b.Start.Lerp(curve.Start, t), b.Control.Lerp(curve.Control, t), b.End.Lerp(curve.End, t)}
}

// Reverse returns the curve run the other way, from End to Start.
func (b QuadraticBezier[T]) Reverse() QuadraticBezier[T] {
	return QuadraticBezier[T]{b.End, b.Control, b.Start}
}

// PointAt returns the point of the curve at the parameter, Start at zero and End at one,
// extrapolated outside [0, 1]. The parameter is no fraction of the length: equal steps of it bunch
// up near the control point.
func (b QuadraticBezier[T]) PointAt(t float64) Point[T] {
	_, point := b.casteljau(t)

	return point.Cast[T]()
}

// TangentAt returns the derivative of the curve at the parameter, zero at an end that the control
// point lies on.
func (b QuadraticBezier[T]) TangentAt(t float64) Vector[float64] {
	first, _ := b.casteljau(t)

	return first[1].Subtract(first[0]).Multiply(2)
}

// Split returns the curve cut at the parameter, from Start to PointAt(t) and from there to End. For
// integer T the new points are rounded, so the halves only approximate the curve.
func (b QuadraticBezier[T]) Split(t float64) (QuadraticBezier[T], QuadraticBezier[T]) {
	first, point := b.casteljau(t)
	middle := point.Cast[T]()

	return QuadraticBezier[T]{b.Start, first[0].Cast[T](), middle}, QuadraticBezier[T]{middle, first[1].Cast[T](), b.End}
}

// IntersectionSegment returns the points where the segment crosses the curve, from Start to End.
func (b QuadraticBezier[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	return segment.IntersectionQuadraticBezier(b)
}

// IntersectionRay returns the points where the ray crosses the curve, from its Origin on.
func (b QuadraticBezier[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return ray.IntersectionQuadraticBezier(b)
}

// IntersectionLine returns the points where the line crosses the curve, in order along its
// Vector.
func (b QuadraticBezier[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return line.IntersectionQuadraticBezier(b)
}

// Equal reports whether the curves have equal points in the same places, so a curve and its Reverse
// differ.
func (b QuadraticBezier[T]) Equal(curve QuadraticBezier[T]) bool {
	return b.Start.Equal(curve.Start) && b.Control.Equal(curve.Control) && b.End.Equal(curve.End)
}

// IsZero reports whether the curve is the zero value.
func (b QuadraticBezier[T]) IsZero() bool {
	return b.Equal(QuadraticBezier[T]{})
}

// Polyline returns the curve flattened to a path within the tolerance of it, from Start to End, as
// CubicBezier.Polyline flattens the cubic of the same curve. A negative tolerance is taken
// absolute.
func (b QuadraticBezier[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{b.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of Polyline(tolerance) to dst and returns the extended slice.
func (b QuadraticBezier[T]) AppendPolyline(dst []Point[T], tolerance float64) []Point[T] {
	return b.Float().CubicBezier().appendPolyline(dst, tolerance)
}

// CubicBezier returns the curve as a cubic, its control points two thirds of the way from each end
// to Control. For integer T they are rounded, so the cubic only approximates the curve.
func (b QuadraticBezier[T]) CubicBezier() CubicBezier[T] {
	return CubicBezier[T]{b.Start, b.Start.Lerp(b.Control, 2.0/3), b.End.Lerp(b.Control, 2.0/3), b.End}
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

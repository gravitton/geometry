package geom

import (
	"cmp"
	"math"
)

// Line is an infinite line through Origin along Vector, running without end both ways.
//
// Only the direction of Vector matters to the points of the line; its length is the unit PointAt
// measures in, and its sense names the sides Side tells apart. A zero Vector gives a line that
// is its origin alone. A line has no Bounds, so it is no Shape.
type Line[T Number] struct {
	Origin Point[T]  `json:"o"`
	Vector Vector[T] `json:"d"`
}

// Ln is shorthand for Line{origin, vector}.
func Ln[T Number](origin Point[T], vector Vector[T]) Line[T] {
	return Line[T]{origin, vector}
}

// LineThrough returns the line through the two points, from the first toward the second, a line of
// zero Vector where they are equal.
func LineThrough[T Number](origin, point Point[T]) Line[T] {
	return Line[T]{origin, point.Subtract(origin)}
}

// Angle returns the angle of Vector in radians, zero for a zero Vector.
func (l Line[T]) Angle() float64 {
	return l.Vector.Angle()
}

// Translate returns the line with its origin moved by the vector.
func (l Line[T]) Translate(vector Vector[T]) Line[T] {
	return Line[T]{l.Origin.Add(vector), l.Vector}
}

// MoveTo returns the line with its origin at the point.
func (l Line[T]) MoveTo(point Point[T]) Line[T] {
	return Line[T]{point, l.Vector}
}

// Rotate returns the line turned by the angle in radians about its origin. For integer T the
// direction is rounded, exact only for quarter turns.
func (l Line[T]) Rotate(angle float64) Line[T] {
	return Line[T]{l.Origin, l.Vector.Rotate(angle)}
}

// Scale returns the line with its Vector scaled by the factor, so PointAt steps the scaled
// length. A negative factor swaps the sides Side names, and a zero factor collapses the line onto
// its origin.
func (l Line[T]) Scale(factor float64) Line[T] {
	return l.ScaleXY(factor, factor)
}

// ScaleXY returns the line with its Vector scaled by the factors along X and Y, which can turn
// it where they differ.
func (l Line[T]) ScaleXY(factorX, factorY float64) Line[T] {
	return Line[T]{l.Origin, l.Vector.MultiplyXY(factorX, factorY)}
}

// Unscale returns the line with its Vector divided by the factor. It panics for a zero factor.
func (l Line[T]) Unscale(factor float64) Line[T] {
	return l.UnscaleXY(factor, factor)
}

// UnscaleXY returns the line with its Vector divided by the factors along X and Y. It panics for
// a zero factor.
func (l Line[T]) UnscaleXY(factorX, factorY float64) Line[T] {
	return Line[T]{l.Origin, l.Vector.DivideXY(factorX, factorY)}
}

// Transform returns the line with the matrix applied to Origin as a point and to Vector as a
// vector, so a translation moves the origin alone.
func (l Line[T]) Transform[M Float](matrix Matrix[M]) Line[T] {
	return Line[T]{l.Origin.Transform(matrix), l.Vector.Transform(matrix)}
}

// PointAt returns the point t times Vector from Origin, on the other side of it for a negative
// t.
func (l Line[T]) PointAt(t float64) Point[T] {
	point := l.Origin.Float().Add(l.Vector.Float().Multiply(t))

	return point.Cast[T]()
}

// Side returns the side of the line the point lies on: 1 on the side the Normal of Vector
// points to, the right of Vector on screen, -1 on the other, and 0 on the line. The sign is
// exact, with no tolerance, so a point Contains admits can still have a side; ask Contains first
// where the two must agree. A zero Vector and a NaN coordinate give 0.
func (l Line[T]) Side(point Point[T]) int {
	return int(Sign(l.side(point)))
}

// DistanceTo returns the distance from the point to the line, zero on it.
func (l Line[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(l.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (l Line[T]) DistanceSquaredTo(point Point[T]) float64 {
	return l.span(point, point).DistanceSquaredTo(point)
}

// Nearest returns the point of the line nearest to the given point, the point itself on the line.
func (l Line[T]) Nearest(point Point[T]) Point[T] {
	return l.span(point, point).Nearest(point)
}

// NearestFraction returns the t at which PointAt gives the point of the line nearest to the given
// point, negative on the other side of the origin, and zero for a zero Vector.
func (l Line[T]) NearestFraction(point Point[T]) float64 {
	direction := l.Vector.Float()
	offset := point.Float().Subtract(l.Origin.Float())

	return ratio(offset.Dot(direction), direction.LengthSquared())
}

// Contains reports whether the point lies on the line.
func (l Line[T]) Contains(point Point[T]) bool {
	return l.DistanceSquaredTo(point) == 0
}

// IntersectsCircle reports whether the line and the circle share a point.
func (l Line[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsLine(l)
}

// IntersectsCapsule reports whether the line and the capsule share a point.
func (l Line[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsLine(l)
}

// IntersectsSector reports whether the line and the sector share a point.
func (l Line[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsLine(l)
}

// IntersectsSegment reports whether the line and the segment share a point.
func (l Line[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsLine(l)
}

// IntersectsRay reports whether the line and the ray share a point.
func (l Line[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsLine(l)
}

// IntersectsLine reports whether the lines share a point.
func (l Line[T]) IntersectsLine(line Line[T]) bool {
	_, ok := l.crossing(line)

	return ok || l.Contains(line.Origin) || line.Contains(l.Origin)
}

// IntersectsPolygon reports whether the line and the polygon share a point.
func (l Line[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return l.span(polygon.minMax()).IntersectsPolygon(polygon)
}

// IntersectsRectangle reports whether the line and the rectangle share a point.
func (l Line[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return l.span(rectangle.MinMax()).IntersectsRectangle(rectangle)
}

// IntersectsRegularPolygon reports whether the line and the regular polygon share a point.
func (l Line[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return l.span(polygon.minMax()).IntersectsRegularPolygon(polygon)
}

// IntersectsBox reports whether the line and the box share a point.
func (l Line[T]) IntersectsBox(box Box[T]) bool {
	return l.span(box.Min, box.Max).IntersectsBox(box)
}

// IntersectsEllipse reports whether the line and the ellipse share a point.
func (l Line[T]) IntersectsEllipse(ellipse Ellipse[T]) bool {
	return l.span(ellipse.minMax()).IntersectsEllipse(ellipse)
}

// IntersectionCircle returns the points where the line crosses the boundary of the circle, in order
// along Vector.
func (l Line[T]) IntersectionCircle(circle Circle[T]) Pair[Point[T]] {
	return circle.IntersectionLine(l)
}

// IntersectionCapsule returns the points where the line crosses the boundary of the capsule, in
// order along Vector.
func (l Line[T]) IntersectionCapsule(capsule Capsule[T]) Pair[Point[T]] {
	return capsule.IntersectionLine(l)
}

// IntersectionSector returns the points where the line crosses the boundary of the sector, in order
// along Vector.
func (l Line[T]) IntersectionSector(sector Sector[T]) []Point[T] {
	return l.AppendIntersectionSector(nil, sector)
}

// AppendIntersectionSector appends the points of IntersectionSector to dst and returns the extended
// slice.
func (l Line[T]) AppendIntersectionSector(dst []Point[T], sector Sector[T]) []Point[T] {
	return l.span(sector.minMax()).AppendIntersectionSector(dst, sector)
}

// IntersectionArc returns the points where the line crosses the arc, in order along Vector.
func (l Line[T]) IntersectionArc(arc Arc[T]) Pair[Point[T]] {
	return arc.IntersectionLine(l)
}

// IntersectionSegment returns the point where the line and the segment cross, and false where they
// do not or are parallel.
func (l Line[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	return segment.IntersectionLine(l)
}

// IntersectionRay returns the point where the line and the ray cross, and false where they do not
// or are parallel.
func (l Line[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	return ray.IntersectionLine(l)
}

// IntersectionLine returns the point where the lines cross, the same whichever line asks, and false
// where they do not. Parallel lines give false even where they run along each other, though
// IntersectsLine holds. A line of zero Vector is parallel to nothing and is the point wherever
// it lies on the other.
func (l Line[T]) IntersectionLine(line Line[T]) (Point[T], bool) {
	if point, ok := l.crossing(line); ok {
		return point.Cast[T](), true
	}

	if l.Vector.hasDirection() && line.Vector.hasDirection() {
		return Point[T]{}, false
	}

	return l.touch(line)
}

// IntersectionPolyline returns the points where the line crosses the path, in order along
// Vector.
func (l Line[T]) IntersectionPolyline(polyline Polyline[T]) []Point[T] {
	return l.AppendIntersectionPolyline(nil, polyline)
}

// AppendIntersectionPolyline appends the points of IntersectionPolyline to dst and returns the
// extended slice.
func (l Line[T]) AppendIntersectionPolyline(dst []Point[T], polyline Polyline[T]) []Point[T] {
	bounds := polyline.Bounds()

	return l.span(bounds.Min, bounds.Max).AppendIntersectionPolyline(dst, polyline)
}

// IntersectionQuadraticBezier returns the points where the line crosses the curve, in order along
// Vector.
func (l Line[T]) IntersectionQuadraticBezier(curve QuadraticBezier[T]) Pair[Point[T]] {
	bounds := curve.Bounds()

	return l.span(bounds.Min, bounds.Max).IntersectionQuadraticBezier(curve)
}

// IntersectionCubicBezier returns the points where the line crosses the curve, in order along
// Vector.
func (l Line[T]) IntersectionCubicBezier(curve CubicBezier[T]) []Point[T] {
	return l.AppendIntersectionCubicBezier(nil, curve)
}

// AppendIntersectionCubicBezier appends the points of IntersectionCubicBezier to dst and returns
// the extended slice.
func (l Line[T]) AppendIntersectionCubicBezier(dst []Point[T], curve CubicBezier[T]) []Point[T] {
	bounds := curve.Bounds()

	return l.span(bounds.Min, bounds.Max).AppendIntersectionCubicBezier(dst, curve)
}

// IntersectionPolygon returns the points where the line crosses the boundary of the polygon, in
// order along Vector.
func (l Line[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return l.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points of IntersectionPolygon to dst and returns the
// extended slice.
func (l Line[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	return l.span(polygon.minMax()).AppendIntersectionPolygon(dst, polygon)
}

// IntersectionRectangle returns the points where the line crosses the boundary of the rectangle, in
// order along Vector.
func (l Line[T]) IntersectionRectangle(rectangle Rectangle[T]) Pair[Point[T]] {
	return l.span(rectangle.MinMax()).IntersectionRectangle(rectangle)
}

// IntersectionRegularPolygon returns the points where the line crosses the boundary of the regular
// polygon, in order along Vector.
func (l Line[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) Pair[Point[T]] {
	return l.span(polygon.minMax()).IntersectionRegularPolygon(polygon)
}

// IntersectionBox returns the points where the line crosses the boundary of the box, in order along
// Vector.
func (l Line[T]) IntersectionBox(box Box[T]) Pair[Point[T]] {
	return l.span(box.Min, box.Max).IntersectionBox(box)
}

// IntersectionEllipse returns the points where the line crosses the boundary of the ellipse, in
// order along Vector.
func (l Line[T]) IntersectionEllipse(ellipse Ellipse[T]) Pair[Point[T]] {
	return l.span(ellipse.minMax()).IntersectionEllipse(ellipse)
}

// ClipCircle returns the chord of the line in the circle, running along Vector, and false where
// they share no point. A tangent gives its one point.
func (l Line[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	return l.span(circle.minMax()).ClipCircle(circle)
}

// ClipCapsule returns the part of the line in the capsule, and false where they share no point.
func (l Line[T]) ClipCapsule(capsule Capsule[T]) (Segment[T], bool) {
	return l.span(capsule.minMax()).ClipCapsule(capsule)
}

// ClipSector returns the parts of the line in the sector, in order along Vector.
func (l Line[T]) ClipSector(sector Sector[T]) []Segment[T] {
	return l.AppendClipSector(nil, sector)
}

// AppendClipSector appends the parts of ClipSector to dst and returns the extended slice.
func (l Line[T]) AppendClipSector(dst []Segment[T], sector Sector[T]) []Segment[T] {
	return l.span(sector.minMax()).AppendClipSector(dst, sector)
}

// ClipPolygon returns the parts of the line in the polygon, in order and each running along
// Vector. A line along a flat polygon has none.
func (l Line[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return l.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts of ClipPolygon to dst and returns the extended slice.
func (l Line[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	return l.span(polygon.minMax()).AppendClipPolygon(dst, polygon)
}

// ClipRectangle returns the part of the line in the rectangle, running along Vector, and false
// where they share no point.
func (l Line[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	return l.span(rectangle.MinMax()).ClipRectangle(rectangle)
}

// ClipRegularPolygon returns the part of the line in the regular polygon, running along Vector,
// and false where they share no point, for an empty polygon, and for a line along a flat one.
func (l Line[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	return l.span(polygon.minMax()).ClipRegularPolygon(polygon)
}

// ClipBox returns the part of the line in the box, running along Vector, and false where they
// share no point.
func (l Line[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	return l.span(box.Min, box.Max).ClipBox(box)
}

// ClipEllipse returns the chord of the line in the ellipse, running along Vector, and false where
// they share no point. A tangent gives its one point.
func (l Line[T]) ClipEllipse(ellipse Ellipse[T]) (Segment[T], bool) {
	return l.span(ellipse.minMax()).ClipEllipse(ellipse)
}

// Reflect returns the point mirrored across the line, and through the origin for a zero Vector.
// For integer T the image is exact across an axis-parallel or diagonal line.
func (l Line[T]) Reflect(point Point[T]) Point[T] {
	origin, direction := l.Origin.Float(), l.Vector.Float()
	offset := point.Float().Subtract(origin)
	twice, lengthSquared := 2*offset.Dot(direction), direction.LengthSquared()

	x := origin.X - offset.X + ratio(twice*direction.X, lengthSquared)
	y := origin.Y - offset.Y + ratio(twice*direction.Y, lengthSquared)

	return Point[T]{Cast[T](x), Cast[T](y)}
}

// side returns the cross product of Vector and the offset of the point from Origin, in float64:
// the value Side takes the sign of and a cut along the line places a crossing by.
func (l Line[T]) side(point Point[T]) float64 {
	offset := point.Float().Subtract(l.Origin.Float())

	return l.Vector.Float().Cross(offset)
}

// halves returns the rays from Origin along Vector and against it.
func (l Line[T]) halves() (Ray[T], Ray[T]) {
	return Ray[T](l), Ray[T]{l.Origin, l.Vector.Negate()}
}

// span returns the segment between the ends of the reaches of both halves past the extent a, b, so
// it shares every point of the shape the line does and pairs read the segment's tests. Past a
// single point it runs along the reach of the ray toward that point.
func (l Line[T]) span(a, b Point[T]) Segment[T] {
	ahead, behind := l.halves()

	return Segment[T]{behind.reach(a, b).End, ahead.reach(a, b).End}
}

// crossing returns where the lines meet, in float64, and false where the cross of their directions
// is zero.
//
// The point is placed by each direction component times the numerator, divided once, so an integer
// crossing on a half unit is exact while the product fits a float64 mantissa. Past that it rounds,
// so the lines are ordered by Point.Compare of Origin then Vector to give the same bits either
// way.
func (l Line[T]) crossing(line Line[T]) (Point[float64], bool) {
	if cmp.Or(line.Origin.Compare(l.Origin), line.Vector.Point().Compare(l.Vector.Point())) < 0 {
		l, line = line, l
	}

	a, b := l.Vector.Float(), line.Vector.Float()
	origin := l.Origin.Float()
	offset := line.Origin.Float().Subtract(origin)
	denominator, along := a.Cross(b), offset.Cross(b)

	if denominator == 0 {
		return Point[float64]{}, false
	}

	return Point[float64]{origin.X + float64(along*a.X)/denominator, origin.Y + float64(along*a.Y)/denominator}, true
}

// touch returns the origin of a line of zero Vector that lies on the other, the one point it can
// share, and false where there is none. Two of zero Vector give the lesser origin by
// Point.Compare, so the answer is the same whichever line asks.
func (l Line[T]) touch(line Line[T]) (Point[T], bool) {
	switch {
	case l.Vector.hasDirection() && l.Contains(line.Origin):
		return line.Origin, true
	case line.Vector.hasDirection() && line.Contains(l.Origin):
		return l.Origin, true
	case !l.Vector.hasDirection() && !line.Vector.hasDirection() && l.Contains(line.Origin):
		if line.Origin.Compare(l.Origin) < 0 {
			return line.Origin, true
		}

		return l.Origin, true
	default:
		return Point[T]{}, false
	}
}

// Equal reports whether the lines have equal origin and direction. Lines of the same points with
// another origin or a direction of another length are not equal.
func (l Line[T]) Equal(line Line[T]) bool {
	return l.Origin.Equal(line.Origin) && l.Vector.Equal(line.Vector)
}

// IsZero reports whether the line is the zero value.
func (l Line[T]) IsZero() bool {
	return l.Equal(Line[T]{})
}

// Cast converts the line to a Line of another number type, rounding as Cast does, so a short
// direction can round to zero for an integer R.
func (l Line[T]) Cast[R Number]() Line[R] {
	return Line[R]{l.Origin.Cast[R](), l.Vector.Cast[R]()}
}

// Int converts the line to a Line[int].
func (l Line[T]) Int() Line[int] {
	return Line[int]{l.Origin.Int(), l.Vector.Int()}
}

// Float converts the line to a Line[float64].
func (l Line[T]) Float() Line[float64] {
	return Line[float64]{l.Origin.Float(), l.Vector.Float()}
}

// String returns the line in the form of its constructor: Ln((x,y);⟨x,y⟩).
func (l Line[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Ln("...)
	buf = l.Origin.appendString(buf)
	buf = append(buf, ';')
	buf = l.Vector.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

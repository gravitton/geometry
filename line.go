package geom

import (
	"cmp"
	"math"
)

// Line is a 2D infinite line, passing through Origin and running without end both ways along
// Direction. Only the direction of the vector matters to the points of the line; its length is
// the unit PointAt measures in, and its sense names the two sides Side tells apart. A zero
// Direction has no way to run, and the line is its origin alone.
//
// A line has no Bounds, so it is no Shape, but it answers Contains, DistanceTo and Nearest like
// one and is a Collider. Every pair of a line with a bounded shape is decided on its span: the
// segment between the points where the two rays from Origin, along Direction and against it,
// pass the extent of the shape, which shares every point of the shape the line does, so the line
// reads the tests of Segment and agrees with them and with a Ray along it.
type Line[T Number] struct {
	Origin    Point[T]  `json:"o"`
	Direction Vector[T] `json:"d"`
}

// Ln is shorthand for Line{origin, direction}, the line through the origin running along the
// direction. A zero direction gives a line that is its origin alone.
func Ln[T Number](origin Point[T], direction Vector[T]) Line[T] {
	return Line[T]{origin, direction}
}

// LineThrough creates a Line through the two points, with the first as its Origin and running
// toward the second. Two equal points give a zero Direction, a line that is its origin alone.
func LineThrough[T Number](origin, point Point[T]) Line[T] {
	return Line[T]{origin, point.Subtract(origin)}
}

// Angle returns the angle of the line in radians, the angle of its Direction. A line with a zero
// Direction runs nowhere and gives 0, the angle Vector.Angle gives it.
func (l Line[T]) Angle() float64 {
	return l.Direction.Angle()
}

// Translate creates a new Line with its origin moved by the given vector and the same direction.
func (l Line[T]) Translate(vector Vector[T]) Line[T] {
	return Line[T]{l.Origin.Add(vector), l.Direction}
}

// MoveTo creates a new Line through the point with the same direction: the origin is the point
// a line is placed by, as the center is on a closed shape.
func (l Line[T]) MoveTo(point Point[T]) Line[T] {
	return Line[T]{point, l.Direction}
}

// Rotate creates a new Line turned by the given angle (in radians) about its origin, in the same
// sense as Vector.Rotate. For integer T the direction is rounded; only multiples of 90° keep it
// exactly.
func (l Line[T]) Rotate(angle float64) Line[T] {
	return Line[T]{l.Origin, l.Direction.Rotate(angle)}
}

// Scale creates a new Line with its direction scaled by the factor about the origin: the points
// of the line stay, and PointAt steps the scaled length. A negative factor turns the direction
// about, which swaps the sides Side names, and a zero factor collapses the line onto its origin.
// For integer T the direction is rounded.
func (l Line[T]) Scale(factor float64) Line[T] {
	return l.ScaleXY(factor, factor)
}

// ScaleXY creates a new Line with its direction scaled by the factors along X and Y about the
// origin, which changes the direction unless the factors are equal.
func (l Line[T]) ScaleXY(factorX, factorY float64) Line[T] {
	return Line[T]{l.Origin, l.Direction.MultiplyXY(factorX, factorY)}
}

// Unscale creates a new Line with its direction scaled by the inverse factor about the origin,
// the inverse of Scale. Like Divide it panics for a zero factor.
func (l Line[T]) Unscale(factor float64) Line[T] {
	return l.UnscaleXY(factor, factor)
}

// UnscaleXY creates a new Line with its direction scaled by the inverse of the given factors
// about the origin, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (l Line[T]) UnscaleXY(factorX, factorY float64) Line[T] {
	return Line[T]{l.Origin, l.Direction.DivideXY(factorX, factorY)}
}

// Transform creates a new Line by applying the given matrix to the origin as a point and to the
// direction as a vector, like Point.Transform and Vector.Transform, so a translation moves the
// origin alone.
func (l Line[T]) Transform[M Float](matrix Matrix[M]) Line[T] {
	return Line[T]{l.Origin.Transform(matrix), l.Direction.Transform(matrix)}
}

// PointAt returns the point t times Direction away from Origin, on the other side of the origin
// for a negative t. The point is computed in float64 and rounded once for an integer T.
func (l Line[T]) PointAt(t float64) Point[T] {
	point := l.Origin.Float().Add(l.Direction.Float().Multiply(t))

	return point.Cast[T]()
}

// Side returns the side of the line the given point lies on, as the sign of the cross product of
// Direction and the vector from Origin to the point: 1 on the side Vector.Normal of the
// direction points to, the right of the direction as drawn on a screen with Y pointing down, -1
// on the other, and 0 on the line. The sign is exact, with no tolerance, so for a float T a
// point Contains admits within the tolerance still has a side; ask Contains first where the two
// must agree. A line with a zero Direction has no sides and a NaN coordinate has no sign, and
// both give 0.
func (l Line[T]) Side(point Point[T]) int {
	return int(Sign(l.side(point)))
}

// DistanceTo returns the distance from the given point to the nearest point of the line, the
// length of the perpendicular. It is zero exactly where Contains holds.
func (l Line[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(l.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, a float64 even for an integer T. It is Segment.DistanceSquaredTo on the span of the
// line past the point, the segment every pair of the line is decided on, which is the reach of
// the ray from Origin toward the point, so the line contains every point a ray along it does. A
// distance within the tolerance is snapped to zero, which is where Contains reads it.
func (l Line[T]) DistanceSquaredTo(point Point[T]) float64 {
	return l.span(point, point).DistanceSquaredTo(point)
}

// Nearest returns the point of the line nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the foot of the perpendicular, as Segment.Nearest finds it
// on the span of the line past the point. For integer T the foot is rounded once and can land
// off the line, where Contains rejects it.
func (l Line[T]) Nearest(point Point[T]) Point[T] {
	return l.span(point, point).Nearest(point)
}

// Reflect returns the given point mirrored across the line: as far beyond the foot of its
// perpendicular as it lies before it. A point on the line is its own mirror image, and a line
// with a zero Direction is its origin alone, which the point is mirrored through. For integer T
// the image is rounded once; across an axis-parallel or a diagonal line it is a lattice point
// and exact.
func (l Line[T]) Reflect(point Point[T]) Point[T] {
	origin, direction := l.Origin.Float(), l.Direction.Float()
	offset := point.Float().Subtract(origin)
	twice, lengthSquared := 2*offset.Dot(direction), direction.LengthSquared()

	x := origin.X - offset.X + ratio(twice*direction.X, lengthSquared)
	y := origin.Y - offset.Y + ratio(twice*direction.Y, lengthSquared)

	return Point[T]{Cast[T](x), Cast[T](y)}
}

// Contains reports whether the given point lies on the line, within the tolerance: it holds
// exactly where DistanceTo is zero.
func (l Line[T]) Contains(point Point[T]) bool {
	return l.DistanceSquaredTo(point) == 0
}

// IntersectsCircle reports whether the line and the circle share a point, as
// Circle.IntersectsLine does.
func (l Line[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsLine(l)
}

// IntersectsSegment reports whether the line and the segment share a point, as
// Segment.IntersectsLine does.
func (l Line[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsLine(l)
}

// IntersectsRay reports whether the line and the ray share a point, as Ray.IntersectsLine does.
func (l Line[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsLine(l)
}

// IntersectsLine reports whether the lines share a point: they are not parallel, or the origin
// of one lies on the other within the tolerance, which covers lines running along each other and
// a line with a zero Direction, its origin alone.
func (l Line[T]) IntersectsLine(line Line[T]) bool {
	_, ok := l.crossing(line)

	return ok || l.Contains(line.Origin) || line.Contains(l.Origin)
}

// IntersectsPolygon reports whether the line and the polygon share a point, as
// Segment.IntersectsPolygon decides it on the span of the line past the polygon: the line
// crosses one of its edges. Touching shapes intersect, within the tolerance, and an empty
// polygon intersects nothing.
func (l Line[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return l.span(polygon.minMax()).IntersectsPolygon(polygon)
}

// IntersectsRectangle reports whether the line and the rectangle share a point, as
// Segment.IntersectsRectangle decides it on the span of the line past the rectangle: the line
// crosses one of its edges. Touching shapes intersect, within the tolerance.
func (l Line[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return l.span(rectangle.MinMax()).IntersectsRectangle(rectangle)
}

// IntersectsRegularPolygon reports whether the line and the regular polygon share a point, as
// Segment.IntersectsRegularPolygon decides it on the span of the line past the polygon: the line
// crosses one of its edges. Touching shapes intersect, within the tolerance, and an empty
// polygon intersects nothing.
func (l Line[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return l.span(polygon.minMax()).IntersectsRegularPolygon(polygon)
}

// IntersectsBox reports whether the line and the box share a point, as Segment.IntersectsBox
// decides it on the span of the line past the box: the line crosses one of its edges. Touching
// shapes intersect, within the tolerance.
func (l Line[T]) IntersectsBox(box Box[T]) bool {
	return l.span(box.Min, box.Max).IntersectsBox(box)
}

// IntersectionCircle returns the points where the line crosses the circle boundary, as
// Circle.IntersectionLine does.
func (l Line[T]) IntersectionCircle(circle Circle[T]) Pair[Point[T]] {
	return circle.IntersectionLine(l)
}

// IntersectionSegment returns the point where the line and the segment cross, and false when
// they do not, as Segment.IntersectionLine does.
func (l Line[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	return segment.IntersectionLine(l)
}

// IntersectionRay returns the point where the line and the ray cross, and false when they do
// not, as Ray.IntersectionLine does.
func (l Line[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	return ray.IntersectionLine(l)
}

// IntersectionLine returns the point where the lines cross, and false when they do not. It
// answers exactly where IntersectsLine holds, less the parallel case: parallel lines share no
// point or every one, and return false even where they run along each other, which
// IntersectsLine still reports. A line with a zero Direction is its origin, parallel to nothing,
// and is the answer wherever it lies on the other line. The point is the same whichever line
// asks. For integer T it is rounded like every other result stored into T, and for a float T
// the crossing of nearly parallel lines lies far along both and is as accurate as their cross
// product.
func (l Line[T]) IntersectionLine(line Line[T]) (Point[T], bool) {
	if point, ok := l.crossing(line); ok {
		return point.Cast[T](), true
	}

	if l.Direction.hasDirection() && line.Direction.hasDirection() {
		return Point[T]{}, false
	}

	return l.touch(line)
}

// IntersectionPolygon returns the points where the line crosses the polygon boundary, in order
// along Direction, as Segment.IntersectionPolygon finds them on the span of the line past the
// polygon.
func (l Line[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return l.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points IntersectionPolygon returns to dst and returns the
// extended slice, as Segment.AppendIntersectionPolygon does on the span of the line past the
// polygon, so a caller reusing dst allocates nothing once it has room.
func (l Line[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	return l.span(polygon.minMax()).AppendIntersectionPolygon(dst, polygon)
}

// IntersectionRectangle returns the points where the line crosses the rectangle boundary, in
// order along Direction, as Segment.IntersectionRectangle finds them on the span of the line
// past the rectangle.
func (l Line[T]) IntersectionRectangle(rectangle Rectangle[T]) Pair[Point[T]] {
	return l.span(rectangle.MinMax()).IntersectionRectangle(rectangle)
}

// IntersectionRegularPolygon returns the points where the line crosses the regular polygon
// boundary, in order along Direction, as Segment.IntersectionRegularPolygon finds them on the
// span of the line past the polygon.
func (l Line[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) Pair[Point[T]] {
	return l.span(polygon.minMax()).IntersectionRegularPolygon(polygon)
}

// IntersectionBox returns the points where the line crosses the box boundary, in order along
// Direction, as Segment.IntersectionBox finds them on the span of the line past the box.
func (l Line[T]) IntersectionBox(box Box[T]) Pair[Point[T]] {
	return l.span(box.Min, box.Max).IntersectionBox(box)
}

// ClipCircle returns the part of the line inside the circle, boundary included within the
// tolerance, and false where they share no point, as Segment.ClipCircle clips the span of the
// line past the circle: the chord, running along Direction, or the one point of a tangent. It
// allocates nothing.
func (l Line[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	return l.span(circle.minMax()).ClipCircle(circle)
}

// ClipPolygon returns the parts of the line inside the polygon, boundary included within the
// tolerance, in order and each running along Direction, as Segment.ClipPolygon clips the span of
// the line past the polygon. A line along the line of a flat polygon has no part, as
// Segment.ClipPolygon says. The result is the one allocation, made on the first part.
func (l Line[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return l.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts ClipPolygon returns to dst and returns the extended slice,
// as Segment.AppendClipPolygon does on the span of the line past the polygon, so a caller
// reusing dst allocates nothing once it has room.
func (l Line[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	return l.span(polygon.minMax()).AppendClipPolygon(dst, polygon)
}

// ClipRectangle returns the part of the line inside the rectangle, boundary included within the
// tolerance, whatever its angle, and false where they share no point, as Segment.ClipRectangle
// clips the span of the line past the rectangle. The part runs along Direction. It allocates
// nothing.
func (l Line[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	return l.span(rectangle.MinMax()).ClipRectangle(rectangle)
}

// ClipRegularPolygon returns the part of the line inside the regular polygon, boundary included
// within the tolerance, and false where they share no point, as Segment.ClipRegularPolygon clips
// the span of the line past the polygon. The part runs along Direction, and an empty polygon
// clips everything away. A line along the line of a flat polygon has no part, as
// Segment.ClipRegularPolygon says. It allocates nothing.
func (l Line[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	return l.span(polygon.minMax()).ClipRegularPolygon(polygon)
}

// ClipBox returns the part of the line inside the box, boundary included within the tolerance,
// and false where they share no point, as Segment.ClipBox clips the span of the line past the
// box: the line as a viewport shows it, running along Direction. It allocates nothing.
func (l Line[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	return l.span(box.Min, box.Max).ClipBox(box)
}

// side returns the cross product of Direction and the vector from Origin to the point, in
// float64: the value Side takes the sign of, which a cut along the line places a crossing by.
func (l Line[T]) side(point Point[T]) float64 {
	offset := point.Float().Subtract(l.Origin.Float())

	return l.Direction.Float().Cross(offset)
}

// halves returns the two rays from Origin that make up the line, the one along Direction and the
// one against it.
func (l Line[T]) halves() (Ray[T], Ray[T]) {
	return Ray[T](l), Ray[T]{l.Origin, l.Direction.Negate()}
}

// span returns the part of the line past the extent a, b of a shape both ways: from the end of
// the reach of the ray against Direction to the end of the reach of the ray along it, each
// Origin where the whole extent lies on the other side. No point of the extent projects past
// either end, so the segment shares every point of the shape the line does, and a pair of the
// line is decided by the segment's own tests. Past a single point it is the reach of the ray
// toward that point, so the line answers a point as that ray does. A line with no direction
// spans its origin alone.
func (l Line[T]) span(a, b Point[T]) Segment[T] {
	ahead, behind := l.halves()

	return Segment[T]{behind.reach(a, b).End, ahead.reach(a, b).End}
}

// crossing returns the point where two lines meet, in float64, and false where the cross product
// of their directions is zero: parallel lines, or a line with no direction.
//
// The point is placed from Origin by each component of the direction times the numerator of the
// fraction, divided by its denominator once, rather than by the rounded fraction: for an integer
// T the product and the cross products are exact, so a crossing on a half unit is exactly there,
// while that product stays within the integers float64 holds exactly. Past it the product
// rounds, so the pair is taken in one order whichever line asks, the lesser by Point.Compare of
// Origin and then Direction first.
func (l Line[T]) crossing(line Line[T]) (Point[float64], bool) {
	if cmp.Or(line.Origin.Compare(l.Origin), line.Direction.Point().Compare(l.Direction.Point())) < 0 {
		l, line = line, l
	}

	a, b := l.Direction.Float(), line.Direction.Float()
	origin := l.Origin.Float()
	offset := line.Origin.Float().Subtract(origin)
	denominator, along := a.Cross(b), offset.Cross(b)

	if denominator == 0 {
		return Point[float64]{}, false
	}

	return Point[float64]{origin.X + float64(along*a.X)/denominator, origin.Y + float64(along*a.Y)/denominator}, true
}

// touch returns the origin of either line that lies on the other, within the tolerance as
// Contains judges it, and false when there is none: the one point a line with no direction can
// share with another.
func (l Line[T]) touch(line Line[T]) (Point[T], bool) {
	switch {
	case l.Contains(line.Origin):
		return line.Origin, true
	case line.Contains(l.Origin):
		return l.Origin, true
	default:
		return Point[T]{}, false
	}
}

// Equal checks if the origins and the directions of the lines are equal. Two lines through the
// same points with another origin, or with directions of different lengths, are different
// values, as a vector is.
func (l Line[T]) Equal(line Line[T]) bool {
	return l.Origin.Equal(line.Origin) && l.Direction.Equal(line.Direction)
}

// IsZero checks if the origin and the direction are zero.
func (l Line[T]) IsZero() bool {
	return l.Equal(Line[T]{})
}

// Cast converts the line to a Line of another number type, rounding as Cast does, so a direction
// shorter than a unit can round to zero for an integer R.
func (l Line[T]) Cast[R Number]() Line[R] {
	return Line[R]{l.Origin.Cast[R](), l.Direction.Cast[R]()}
}

// Int converts the line to a Line[int].
func (l Line[T]) Int() Line[int] {
	return Line[int]{l.Origin.Int(), l.Direction.Int()}
}

// Float converts the line to a Line[float64].
func (l Line[T]) Float() Line[float64] {
	return Line[float64]{l.Origin.Float(), l.Direction.Float()}
}

// String returns the line by its origin and direction: Line((x,y);⟨x,y⟩).
func (l Line[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Line("...)
	buf = l.Origin.appendString(buf)
	buf = append(buf, ';')
	buf = l.Direction.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

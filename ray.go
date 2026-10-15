package geom

import (
	"cmp"
	"math"
)

// Ray is a half-line from Origin along Vector without end.
//
// Only the direction of Vector matters to the points of the ray; its length is the unit PointAt
// and every Fraction measure in. A zero Vector gives a ray that is its origin alone. A ray has
// no Bounds, so it is no Shape.
type Ray[T Number] struct {
	Origin Point[T]  `json:"o"`
	Vector Vector[T] `json:"d"`
}

// Ry is shorthand for Ray{origin, vector}.
func Ry[T Number](origin Point[T], vector Vector[T]) Ray[T] {
	return Ray[T]{origin, vector}
}

// RayThrough returns the ray from the origin through the point, a ray of zero Vector where they
// are equal.
func RayThrough[T Number](origin, point Point[T]) Ray[T] {
	return Ray[T]{origin, point.Subtract(origin)}
}

// Angle returns the angle of Vector in radians, zero for a zero Vector.
func (r Ray[T]) Angle() float64 {
	return r.Vector.Angle()
}

// Direction returns the direction nearest to Vector, or DirectionNone for a zero Vector.
func (r Ray[T]) Direction() Direction {
	return r.Vector.Direction()
}

// Normal returns Vector turned a quarter turn in the sense of Vector.Rotate, of the same length. A
// zero Vector gives the zero vector.
func (r Ray[T]) Normal() Vector[T] {
	return r.Vector.Normal()
}

// Translate returns the ray with its origin moved by the vector.
func (r Ray[T]) Translate(vector Vector[T]) Ray[T] {
	return Ray[T]{r.Origin.Add(vector), r.Vector}
}

// MoveTo returns the ray with its origin at the point.
func (r Ray[T]) MoveTo(point Point[T]) Ray[T] {
	return Ray[T]{point, r.Vector}
}

// Rotate returns the ray turned by the angle in radians about its origin. For integer T the
// direction is rounded, exact only for quarter turns.
func (r Ray[T]) Rotate(angle float64) Ray[T] {
	return Ray[T]{r.Origin, r.Vector.Rotate(angle)}
}

// Scale returns the ray with its Vector scaled by the factor, so PointAt steps the scaled
// length. A negative factor turns the ray to the other side of its origin, and a zero factor
// collapses it onto its origin.
func (r Ray[T]) Scale(factor float64) Ray[T] {
	return r.ScaleXY(factor, factor)
}

// ScaleXY returns the ray with its Vector scaled by the factors along X and Y, which can turn it
// where they differ.
func (r Ray[T]) ScaleXY(factorX, factorY float64) Ray[T] {
	return Ray[T]{r.Origin, r.Vector.MultiplyXY(factorX, factorY)}
}

// Unscale returns the ray with its Vector divided by the factor. It panics for a zero factor.
func (r Ray[T]) Unscale(factor float64) Ray[T] {
	return r.UnscaleXY(factor, factor)
}

// UnscaleXY returns the ray with its Vector divided by the factors along X and Y. It panics for
// a zero factor.
func (r Ray[T]) UnscaleXY(factorX, factorY float64) Ray[T] {
	return Ray[T]{r.Origin, r.Vector.DivideXY(factorX, factorY)}
}

// Transform returns the ray with the matrix applied to Origin as a point and to Vector as a
// vector, so a translation moves the origin alone.
func (r Ray[T]) Transform[M Float](matrix Matrix[M]) Ray[T] {
	return Ray[T]{r.Origin.Transform(matrix), r.Vector.Transform(matrix)}
}

// Lerp returns the ray interpolated toward the given one, Origin and Vector together, and
// extrapolated outside [0, 1].
func (r Ray[T]) Lerp(ray Ray[T], t float64) Ray[T] {
	return Ray[T]{r.Origin.Lerp(ray.Origin, t), r.Vector.Lerp(ray.Vector, t)}
}

// PointAt returns the point t times Vector from Origin, behind it for a negative t.
func (r Ray[T]) PointAt(t float64) Point[T] {
	point := r.Origin.Float().Add(r.Vector.Float().Multiply(t))

	return point.Cast[T]()
}

// DistanceTo returns the distance from the point to the ray, zero on it.
func (r Ray[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(r.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (r Ray[T]) DistanceSquaredTo(point Point[T]) float64 {
	return r.reach(point, point).DistanceSquaredTo(point)
}

// Nearest returns the point of the ray nearest to the given point, the point itself on the ray and
// the origin for a point behind it.
func (r Ray[T]) Nearest(point Point[T]) Point[T] {
	return r.reach(point, point).Nearest(point)
}

// NearestFraction returns the t at which PointAt gives the point of the ray nearest to the given
// point, zero for a point behind the origin and for a zero Vector.
func (r Ray[T]) NearestFraction(point Point[T]) float64 {
	direction := r.Vector.Float()
	offset := point.Float().Subtract(r.Origin.Float())

	return max(ratio(offset.Dot(direction), direction.LengthSquared()), 0)
}

// Contains reports whether the point lies on the ray.
func (r Ray[T]) Contains(point Point[T]) bool {
	return r.DistanceSquaredTo(point) == 0
}

// IntersectsCircle reports whether the ray and the circle share a point.
func (r Ray[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRay(r)
}

// IntersectsCapsule reports whether the ray and the capsule share a point.
func (r Ray[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsRay(r)
}

// IntersectsSector reports whether the ray and the sector share a point.
func (r Ray[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsRay(r)
}

// IntersectsSegment reports whether the ray and the segment share a point.
func (r Ray[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsRay(r)
}

// IntersectsRay reports whether the rays share a point.
func (r Ray[T]) IntersectsRay(ray Ray[T]) bool {
	return r.crosses(ray) || r.Contains(ray.Origin) || ray.Contains(r.Origin)
}

// IntersectsLine reports whether the ray and the line share a point.
func (r Ray[T]) IntersectsLine(line Line[T]) bool {
	ahead, behind := line.halves()

	return r.IntersectsRay(ahead) || r.IntersectsRay(behind)
}

// IntersectsPolygon reports whether the ray and the polygon share a point.
func (r Ray[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return r.reach(polygon.minMax()).IntersectsPolygon(polygon)
}

// IntersectsRectangle reports whether the ray and the rectangle share a point.
func (r Ray[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return r.reach(rectangle.minMax()).IntersectsRectangle(rectangle)
}

// IntersectsRegularPolygon reports whether the ray and the regular polygon share a point.
func (r Ray[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return r.reach(polygon.minMax()).IntersectsRegularPolygon(polygon)
}

// IntersectsBox reports whether the ray and the box share a point.
func (r Ray[T]) IntersectsBox(box Box[T]) bool {
	return r.reach(box.Min, box.Max).IntersectsBox(box)
}

// IntersectsEllipse reports whether the ray and the ellipse share a point.
func (r Ray[T]) IntersectsEllipse(ellipse Ellipse[T]) bool {
	return r.reach(ellipse.minMax()).IntersectsEllipse(ellipse)
}

// IntersectionCircle returns the points where the ray crosses the boundary of the circle, from
// Origin on.
func (r Ray[T]) IntersectionCircle(circle Circle[T]) Pair[Point[T]] {
	return circle.IntersectionRay(r)
}

// IntersectionCapsule returns the points where the ray crosses the boundary of the capsule, from
// Origin on.
func (r Ray[T]) IntersectionCapsule(capsule Capsule[T]) Pair[Point[T]] {
	return capsule.IntersectionRay(r)
}

// IntersectionSector returns the points where the ray crosses the boundary of the sector, from
// Origin on.
func (r Ray[T]) IntersectionSector(sector Sector[T]) []Point[T] {
	return r.AppendIntersectionSector(nil, sector)
}

// AppendIntersectionSector appends the points of IntersectionSector to dst and returns the extended
// slice.
func (r Ray[T]) AppendIntersectionSector(dst []Point[T], sector Sector[T]) []Point[T] {
	return r.reach(sector.minMax()).AppendIntersectionSector(dst, sector)
}

// IntersectionArc returns the points where the ray crosses the arc, from Origin on.
func (r Ray[T]) IntersectionArc(arc Arc[T]) Pair[Point[T]] {
	return arc.IntersectionRay(r)
}

// IntersectionSegment returns the point where the ray and the segment cross, and false where they
// do not or are parallel.
func (r Ray[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	return segment.IntersectionRay(r)
}

// IntersectionRay returns the point where the rays cross, and false where they do not. Parallel
// rays give false even where they overlap, though IntersectsRay holds. A ray of zero Vector is
// parallel to nothing and is the point wherever it lies on the other. Rays nearly along each other
// can share more than one point; the given ray's origin is then taken before this one's, so the two
// orders can differ.
func (r Ray[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	if point, ok := r.crossing(ray); ok {
		return point.Cast[T](), true
	}

	if r.parallel(ray) {
		return Point[T]{}, false
	}

	return r.touch(ray)
}

// IntersectionLine returns the point where the ray and the line cross, and false where they do not
// or are parallel.
func (r Ray[T]) IntersectionLine(line Line[T]) (Point[T], bool) {
	ahead, behind := line.halves()
	if point, ok := r.IntersectionRay(ahead); ok {
		return point, true
	}

	return r.IntersectionRay(behind)
}

// IntersectionPolyline returns the points where the ray crosses the path, from Origin on.
func (r Ray[T]) IntersectionPolyline(polyline Polyline[T]) []Point[T] {
	return r.AppendIntersectionPolyline(nil, polyline)
}

// AppendIntersectionPolyline appends the points of IntersectionPolyline to dst and returns the
// extended slice.
func (r Ray[T]) AppendIntersectionPolyline(dst []Point[T], polyline Polyline[T]) []Point[T] {
	bounds := polyline.Bounds()

	return r.reach(bounds.Min, bounds.Max).AppendIntersectionPolyline(dst, polyline)
}

// IntersectionQuadraticBezier returns the points where the ray crosses the curve, from Origin on.
func (r Ray[T]) IntersectionQuadraticBezier(curve QuadraticBezier[T]) Pair[Point[T]] {
	bounds := curve.Bounds()

	return r.reach(bounds.Min, bounds.Max).IntersectionQuadraticBezier(curve)
}

// IntersectionCubicBezier returns the points where the ray crosses the curve, from Origin on.
func (r Ray[T]) IntersectionCubicBezier(curve CubicBezier[T]) []Point[T] {
	return r.AppendIntersectionCubicBezier(nil, curve)
}

// AppendIntersectionCubicBezier appends the points of IntersectionCubicBezier to dst and returns
// the extended slice.
func (r Ray[T]) AppendIntersectionCubicBezier(dst []Point[T], curve CubicBezier[T]) []Point[T] {
	bounds := curve.Bounds()

	return r.reach(bounds.Min, bounds.Max).AppendIntersectionCubicBezier(dst, curve)
}

// IntersectionPolygon returns the points where the ray crosses the boundary of the polygon, from
// Origin on.
func (r Ray[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return r.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points of IntersectionPolygon to dst and returns the
// extended slice.
func (r Ray[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	return r.reach(polygon.minMax()).AppendIntersectionPolygon(dst, polygon)
}

// IntersectionRectangle returns the points where the ray crosses the boundary of the rectangle,
// from Origin on.
func (r Ray[T]) IntersectionRectangle(rectangle Rectangle[T]) Pair[Point[T]] {
	return r.reach(rectangle.minMax()).IntersectionRectangle(rectangle)
}

// IntersectionRegularPolygon returns the points where the ray crosses the boundary of the regular
// polygon, from Origin on.
func (r Ray[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) Pair[Point[T]] {
	return r.reach(polygon.minMax()).IntersectionRegularPolygon(polygon)
}

// IntersectionBox returns the points where the ray crosses the boundary of the box, from Origin on.
func (r Ray[T]) IntersectionBox(box Box[T]) Pair[Point[T]] {
	return r.reach(box.Min, box.Max).IntersectionBox(box)
}

// IntersectionEllipse returns the points where the ray crosses the boundary of the ellipse, from
// Origin on.
func (r Ray[T]) IntersectionEllipse(ellipse Ellipse[T]) Pair[Point[T]] {
	return r.reach(ellipse.minMax()).IntersectionEllipse(ellipse)
}

// ClipCircle returns the part of the ray in the circle, and false where they share no point. It
// starts at Origin where the circle contains it, else at the first point of the circle the ray
// reaches.
func (r Ray[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	return r.reach(circle.minMax()).ClipCircle(circle)
}

// ClipCapsule returns the part of the ray in the capsule, and false where they share no point.
func (r Ray[T]) ClipCapsule(capsule Capsule[T]) (Segment[T], bool) {
	return r.reach(capsule.minMax()).ClipCapsule(capsule)
}

// ClipSector returns the parts of the ray in the sector, from Origin on.
func (r Ray[T]) ClipSector(sector Sector[T]) []Segment[T] {
	return r.AppendClipSector(nil, sector)
}

// AppendClipSector appends the parts of ClipSector to dst and returns the extended slice.
func (r Ray[T]) AppendClipSector(dst []Segment[T], sector Sector[T]) []Segment[T] {
	return r.reach(sector.minMax()).AppendClipSector(dst, sector)
}

// ClipPolygon returns the parts of the ray in the polygon, from Origin on. A ray along the line of
// a flat polygon has none.
func (r Ray[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return r.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts of ClipPolygon to dst and returns the extended slice.
func (r Ray[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	return r.reach(polygon.minMax()).AppendClipPolygon(dst, polygon)
}

// ClipRectangle returns the part of the ray in the rectangle, and false where they share no point.
func (r Ray[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	return r.reach(rectangle.minMax()).ClipRectangle(rectangle)
}

// ClipRegularPolygon returns the part of the ray in the regular polygon, and false where they share
// no point, for an empty polygon, and for a ray along the line of a flat one.
func (r Ray[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	return r.reach(polygon.minMax()).ClipRegularPolygon(polygon)
}

// ClipBox returns the part of the ray in the box, and false where they share no point.
func (r Ray[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	return r.reach(box.Min, box.Max).ClipBox(box)
}

// ClipEllipse returns the part of the ray in the ellipse, and false where they share no point. It
// starts at Origin where the ellipse contains it, else at the first point of the ellipse the ray
// reaches.
func (r Ray[T]) ClipEllipse(ellipse Ellipse[T]) (Segment[T], bool) {
	return r.reach(ellipse.minMax()).ClipEllipse(ellipse)
}

// HitCircle returns where the ray first reaches the circle from outside, with the normal pointing
// from the center, and false where it does not or Origin is inside.
func (r Ray[T]) HitCircle(circle Circle[T]) (Hit[T], bool) {
	return r.rescale(r.reach(circle.minMax()).HitCircle(circle))
}

// HitCapsule returns where the ray first reaches the capsule from outside, and false where it does
// not or Origin is inside.
func (r Ray[T]) HitCapsule(capsule Capsule[T]) (Hit[T], bool) {
	return r.rescale(r.reach(capsule.minMax()).HitCapsule(capsule))
}

// HitSector returns where the ray first reaches the sector from outside, and false where it does
// not or Origin is inside.
func (r Ray[T]) HitSector(sector Sector[T]) (Hit[T], bool) {
	return r.rescale(r.reach(sector.minMax()).HitSector(sector))
}

// HitArc returns where the ray first reaches the arc, and false where it does not or Origin lies on
// the arc.
func (r Ray[T]) HitArc(arc Arc[T]) (Hit[T], bool) {
	return r.rescale(r.reach(arc.minMax()).HitArc(arc))
}

// HitSegment returns where the ray first reaches the segment as a wall, and false where it does not
// or Origin lies on the segment.
func (r Ray[T]) HitSegment(segment Segment[T]) (Hit[T], bool) {
	return r.rescale(r.reach(segment.minMax()).HitSegment(segment))
}

// HitRay returns where the ray first reaches the given one as a wall, and false where it does not
// or Origin lies on the wall. A wall running along the ray is met at its origin with the normal
// back along the ray, as is a wall of zero Vector.
func (r Ray[T]) HitRay(ray Ray[T]) (Hit[T], bool) {
	if ray.Contains(r.Origin) {
		return Hit[T]{}, false
	}

	point, crossed := r.IntersectionRay(ray)

	switch {
	case crossed:
		return r.hitAt(point, ray.Vector.Float().Normal()), true
	case r.Contains(ray.Origin):
		return r.hitAt(ray.Origin, Vector[float64]{}), true
	default:
		return Hit[T]{}, false
	}
}

// HitLine returns where the ray first reaches the line as a wall, and false where it does not or
// Origin lies on the line.
func (r Ray[T]) HitLine(line Line[T]) (Hit[T], bool) {
	if line.Contains(r.Origin) {
		return Hit[T]{}, false
	}

	ahead, behind := line.halves()
	if hit, ok := r.HitRay(ahead); ok {
		return hit, true
	}

	return r.HitRay(behind)
}

// HitPolyline returns where the ray first reaches the path as a wall, and false where it does not
// or Origin lies on the path.
func (r Ray[T]) HitPolyline(polyline Polyline[T]) (Hit[T], bool) {
	bounds := polyline.Bounds()

	return r.rescale(r.reach(bounds.Min, bounds.Max).HitPolyline(polyline))
}

// HitQuadraticBezier returns where the ray first reaches the curve as a wall, and false where it
// does not or Origin lies on it.
func (r Ray[T]) HitQuadraticBezier(curve QuadraticBezier[T]) (Hit[T], bool) {
	bounds := curve.Bounds()

	return r.rescale(r.reach(bounds.Min, bounds.Max).HitQuadraticBezier(curve))
}

// HitCubicBezier returns where the ray first reaches the curve as a wall, and false where it does
// not or Origin lies on it.
func (r Ray[T]) HitCubicBezier(curve CubicBezier[T]) (Hit[T], bool) {
	bounds := curve.Bounds()

	return r.rescale(r.reach(bounds.Min, bounds.Max).HitCubicBezier(curve))
}

// HitPolygon returns where the ray first reaches the polygon from outside, and false where it does
// not or Origin is inside.
func (r Ray[T]) HitPolygon(polygon Polygon[T]) (Hit[T], bool) {
	return r.rescale(r.reach(polygon.minMax()).HitPolygon(polygon))
}

// HitRectangle returns where the ray first reaches the rectangle from outside, and false where it
// does not or Origin is inside.
func (r Ray[T]) HitRectangle(rectangle Rectangle[T]) (Hit[T], bool) {
	return r.rescale(r.reach(rectangle.minMax()).HitRectangle(rectangle))
}

// HitRegularPolygon returns where the ray first reaches the regular polygon from outside, and false
// where it does not, Origin is inside, or the polygon is empty.
func (r Ray[T]) HitRegularPolygon(polygon RegularPolygon[T]) (Hit[T], bool) {
	return r.rescale(r.reach(polygon.minMax()).HitRegularPolygon(polygon))
}

// HitBox returns where the ray first reaches the box from outside, and false where it does not or
// Origin is inside.
func (r Ray[T]) HitBox(box Box[T]) (Hit[T], bool) {
	return r.rescale(r.reach(box.Min, box.Max).HitBox(box))
}

// HitEllipse returns where the ray first reaches the ellipse from outside, and false where it does
// not or Origin is inside.
func (r Ray[T]) HitEllipse(ellipse Ellipse[T]) (Hit[T], bool) {
	return r.rescale(r.reach(ellipse.minMax()).HitEllipse(ellipse))
}

// reach returns the segment from Origin to twice the fraction the farthest corner of the extent a,
// b projects to, or Origin alone where the extent lies behind or there is no direction. It shares
// every point of the shape the ray does, so pairs read the segment's tests; ending well past the
// extent, it is never taken for a touch beside a vertex. For integer T the fraction is rounded up,
// so the end is a lattice point on the ray.
func (r Ray[T]) reach(a, b Point[T]) Segment[T] {
	if !r.Vector.hasDirection() {
		return Segment[T]{r.Origin, r.Origin}
	}

	direction := r.Vector.Float()
	corner := a
	if direction.X > 0 {
		corner.X = b.X
	}
	if direction.Y > 0 {
		corner.Y = b.Y
	}

	t := max(2*corner.Float().Subtract(r.Origin.Float()).Dot(direction)/direction.LengthSquared(), 0)
	if isInt[T]() {
		t = math.Ceil(t)
	}

	return Segment[T]{r.Origin, r.PointAt(t)}
}

// rescale returns a hit of the reach as a hit of the ray, its Fraction remeasured in units of
// Vector.
func (r Ray[T]) rescale(hit Hit[T], ok bool) (Hit[T], bool) {
	if !ok {
		return Hit[T]{}, false
	}

	direction := r.Vector.Float()
	fraction := ratio(hit.Point.Float().Subtract(r.Origin.Float()).Dot(direction), direction.LengthSquared())

	return Hit[T]{hit.Point, hit.Normal, fraction}, true
}

// hitAt returns the Hit at the point of a boundary of the given normal, built on the first step of
// the ray so the Fraction is in units of Vector.
func (r Ray[T]) hitAt(point Point[T], normal Vector[float64]) Hit[T] {
	return Segment[T]{r.Origin, r.Origin.Add(r.Vector)}.hitAt(point, normal)
}

// crosses reports whether the rays meet strictly ahead of both origins. Touching and parallel rays
// are left to the origin tests.
func (r Ray[T]) crosses(ray Ray[T]) bool {
	_, ok := r.crossing(ray)

	return ok
}

// crossing returns where the lines through two properly crossing rays meet, in float64, and false
// otherwise. The fractions are read from signs of cross products with no division, so NaN crosses
// nothing. Rounding can give nearly collinear float rays the signs of a crossing at a point
// anywhere along this ray, so a point behind the other origin is rejected and left to the origin
// tests.
//
// The point is placed by each direction component times the numerator, divided once, so an integer
// crossing on a half unit is exact while the product fits a float64 mantissa. Past that it rounds,
// so the rays are ordered by Point.Compare of Origin then Vector to give the same bits either
// way.
func (r Ray[T]) crossing(ray Ray[T]) (Point[float64], bool) {
	if cmp.Or(ray.Origin.Compare(r.Origin), ray.Vector.Point().Compare(r.Vector.Point())) < 0 {
		r, ray = ray, r
	}

	a, b := r.Vector.Float(), ray.Vector.Float()
	origin := r.Origin.Float()
	offset := ray.Origin.Float().Subtract(origin)
	denominator, along, other := a.Cross(b), offset.Cross(b), offset.Cross(a)

	ahead := (denominator > 0 && along > 0 && other > 0) || (denominator < 0 && along < 0 && other < 0)
	if !ahead {
		return Point[float64]{}, false
	}

	point := Point[float64]{origin.X + float64(along*a.X)/denominator, origin.Y + float64(along*a.Y)/denominator}

	return point, point.Subtract(ray.Origin.Float()).Dot(b) >= 0
}

// parallel reports whether the rays have the same line direction, either way. A zero Vector is
// parallel to nothing, so its origin can still be found on the other.
func (r Ray[T]) parallel(ray Ray[T]) bool {
	a, b := r.Vector.Float(), ray.Vector.Float()

	return a.hasDirection() && b.hasDirection() && a.Cross(b) == 0
}

// touch returns the origin of either ray that lies on the other, the given ray's first, and false
// where there is none.
func (r Ray[T]) touch(ray Ray[T]) (Point[T], bool) {
	switch {
	case r.Contains(ray.Origin):
		return ray.Origin, true
	case ray.Contains(r.Origin):
		return r.Origin, true
	default:
		return Point[T]{}, false
	}
}

// Equal reports whether the rays have equal origin and direction. Rays of the same points with
// directions of different lengths are not equal.
func (r Ray[T]) Equal(ray Ray[T]) bool {
	return r.Origin.Equal(ray.Origin) && r.Vector.Equal(ray.Vector)
}

// IsZero reports whether the ray is the zero value.
func (r Ray[T]) IsZero() bool {
	return r.Equal(Ray[T]{})
}

// Line returns the line through the ray, with the same Origin and Vector.
func (r Ray[T]) Line() Line[T] {
	return Line[T](r)
}

// Cast converts the ray to a Ray of another number type, rounding as Cast does, so a short
// direction can round to zero for an integer R.
func (r Ray[T]) Cast[R Number]() Ray[R] {
	return Ray[R]{r.Origin.Cast[R](), r.Vector.Cast[R]()}
}

// Int converts the ray to a Ray[int].
func (r Ray[T]) Int() Ray[int] {
	return Ray[int]{r.Origin.Int(), r.Vector.Int()}
}

// Float converts the ray to a Ray[float64].
func (r Ray[T]) Float() Ray[float64] {
	return Ray[float64]{r.Origin.Float(), r.Vector.Float()}
}

// String returns the ray by its origin and direction: Ray((x,y);⟨x,y⟩).
func (r Ray[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Ray("...)
	buf = r.Origin.appendString(buf)
	buf = append(buf, ';')
	buf = r.Vector.appendString(buf)
	buf = append(buf, ')')

	return string(buf)
}

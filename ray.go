package geom

import (
	"cmp"
	"fmt"
	"math"
)

// Ray is a 2D half-line, starting at Origin and running without end along Direction. Only the
// direction of the vector matters to the points of the ray; its length is the unit PointAt
// measures in, so a float ray with a normalized Direction measures distances along it. A zero
// Direction has no way to run, and the ray is its origin alone.
//
// A ray has no Bounds, so it is no Shape, but it answers Contains, DistanceTo and Nearest like
// one and is a Collider. Every pair of a ray with a bounded shape is decided on its reach: the
// segment from Origin to where the ray passes the extent of the shape, which shares every point
// of the shape the ray does, so the ray reads the tests of Segment and agrees with them.
type Ray[T Number] struct {
	Origin    Point[T]  `json:"o"`
	Direction Vector[T] `json:"d"`
}

// Ry is shorthand for Ray{origin, direction}, the ray starting at the origin and running along
// the direction. A zero direction gives a ray that is its origin alone.
func Ry[T Number](origin Point[T], direction Vector[T]) Ray[T] {
	return Ray[T]{origin, direction}
}

// RayThrough creates a Ray starting at the origin and passing through the point, the ray a
// picking routine casts from a camera toward a cursor. A point equal to the origin gives a zero
// Direction, a ray that is its origin alone.
func RayThrough[T Number](origin, point Point[T]) Ray[T] {
	return Ray[T]{origin, point.Subtract(origin)}
}

// Angle returns the angle of the ray in radians, the angle of its Direction. A ray with a zero
// Direction runs nowhere and gives 0, the angle Vector.Angle gives it.
func (r Ray[T]) Angle() float64 {
	return r.Direction.Angle()
}

// Translate creates a new Ray with its origin moved by the given vector and the same direction.
func (r Ray[T]) Translate(vector Vector[T]) Ray[T] {
	return Ray[T]{r.Origin.Add(vector), r.Direction}
}

// MoveTo creates a new Ray starting at the point with the same direction: the origin is the
// point a ray is placed by, as the center is on a closed shape.
func (r Ray[T]) MoveTo(point Point[T]) Ray[T] {
	return Ray[T]{point, r.Direction}
}

// Scale creates a new Ray with its direction scaled by the factor about the origin: the points of
// the ray stay, and PointAt steps the scaled length. A negative factor turns the ray to the other
// side of its origin, and a zero factor collapses the ray onto its origin. For integer T the
// direction is rounded.
func (r Ray[T]) Scale(factor float64) Ray[T] {
	return r.ScaleXY(factor, factor)
}

// ScaleXY creates a new Ray with its direction scaled by the factors along X and Y about the
// origin, which changes the direction unless the factors are equal.
func (r Ray[T]) ScaleXY(factorX, factorY float64) Ray[T] {
	return Ray[T]{r.Origin, r.Direction.MultiplyXY(factorX, factorY)}
}

// Unscale creates a new Ray with its direction scaled by the inverse factor about the origin,
// the inverse of Scale. Like Divide it panics for a zero factor.
func (r Ray[T]) Unscale(factor float64) Ray[T] {
	return r.UnscaleXY(factor, factor)
}

// UnscaleXY creates a new Ray with its direction scaled by the inverse of the given factors
// about the origin, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (r Ray[T]) UnscaleXY(factorX, factorY float64) Ray[T] {
	return Ray[T]{r.Origin, r.Direction.DivideXY(factorX, factorY)}
}

// PointAt returns the point t times Direction away from Origin, behind the origin for a
// negative t. The point is computed in float64 and rounded once for an integer T.
func (r Ray[T]) PointAt(t float64) Point[T] {
	point := r.Origin.Float().Add(r.Direction.Float().Multiply(t))

	return point.Cast[T]()
}

// Transform creates a new Ray by applying the given matrix to the origin as a point and to the
// direction as a vector, like Point.Transform and Vector.Transform, so a translation moves the
// origin alone.
func (r Ray[T]) Transform[M Float](matrix Matrix[M]) Ray[T] {
	return Ray[T]{r.Origin.Transform(matrix), r.Direction.Transform(matrix)}
}

// Rotate creates a new Ray turned by the given angle (in radians) about its origin, in the same
// sense as Vector.Rotate. For integer T the direction is rounded; only multiples of 90° keep it
// exactly.
func (r Ray[T]) Rotate(angle float64) Ray[T] {
	return Ray[T]{r.Origin, r.Direction.Rotate(angle)}
}

// Contains reports whether the given point lies on the ray, within the tolerance: it holds
// exactly where DistanceTo is zero.
func (r Ray[T]) Contains(point Point[T]) bool {
	return r.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the ray: the
// distance to the line through it for a point ahead of the origin, and to the origin for a point
// behind it. It is zero exactly where Contains holds.
func (r Ray[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(r.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, a float64 even for an integer T. It is Segment.DistanceSquaredTo on the reach of
// the ray past the point, the segment every pair of the ray is decided on, so a point the ray
// contains is one its pairs find on it. A distance within the tolerance is snapped to zero,
// which is where Contains reads it.
func (r Ray[T]) DistanceSquaredTo(point Point[T]) float64 {
	return r.reach(point, point).DistanceSquaredTo(point)
}

// Nearest returns the point of the ray nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the foot of the perpendicular, or the origin where the foot
// falls behind it, as Segment.Nearest finds it on the reach of the ray past the point. For
// integer T the foot is rounded once and can land off the ray, where Contains rejects it.
func (r Ray[T]) Nearest(point Point[T]) Point[T] {
	return r.reach(point, point).Nearest(point)
}

// IntersectsCircle reports whether the ray and the circle share a point, as
// Circle.IntersectsRay does.
func (r Ray[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRay(r)
}

// IntersectionCircle returns the points where the ray crosses the circle boundary, as
// Circle.IntersectionRay does.
func (r Ray[T]) IntersectionCircle(circle Circle[T]) []Point[T] {
	return circle.IntersectionRay(r)
}

// AppendIntersectionCircle appends the points IntersectionCircle returns to dst and returns the
// extended slice, as Circle.AppendIntersectionRay does.
func (r Ray[T]) AppendIntersectionCircle(dst []Point[T], circle Circle[T]) []Point[T] {
	return circle.AppendIntersectionRay(dst, r)
}

// IntersectsSegment reports whether the ray and the segment share a point, as
// Segment.IntersectsRay does.
func (r Ray[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsRay(r)
}

// IntersectionSegment returns the point where the ray and the segment cross, and false when they
// do not, as Segment.IntersectionRay does.
func (r Ray[T]) IntersectionSegment(segment Segment[T]) (Point[T], bool) {
	return segment.IntersectionRay(r)
}

// IntersectsRay reports whether the rays share a point, within the tolerance, the same closed
// convention as Contains: the lines through them meet ahead of both origins, decided on the
// exact signs of cross products, or the origin of one lies on the other, which covers rays
// running along each other and rays meeting at an origin.
func (r Ray[T]) IntersectsRay(ray Ray[T]) bool {
	return r.crosses(ray) || r.Contains(ray.Origin) || ray.Contains(r.Origin)
}

// IntersectionRay returns the point where the rays cross, and false when they do not. It answers
// exactly where IntersectsRay holds, less the parallel case: parallel rays have no single
// crossing point and return false even where they overlap, which IntersectsRay still reports. A
// ray with a zero Direction is its origin, parallel to nothing, and is the answer wherever it
// lies on the other ray. For integer T the crossing is rounded like every other result stored
// into T. Rays running along each other within the tolerance, short of parallel, share more than
// one point, and the origin returned is the given ray's before this one's, so the two orders can
// name different origins.
func (r Ray[T]) IntersectionRay(ray Ray[T]) (Point[T], bool) {
	if point, ok := r.crossing(ray); ok {
		return point.Cast[T](), true
	}

	if r.parallel(ray) {
		return Point[T]{}, false
	}

	return r.touch(ray)
}

// IntersectsPolygon reports whether the ray and the polygon share a point, as
// Segment.IntersectsPolygon decides it on the reach of the ray past the polygon: the origin
// lies within the polygon, or the ray crosses one of its edges. Touching shapes intersect,
// within the tolerance, and an empty polygon intersects nothing.
func (r Ray[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return r.reach(polygon.minMax()).IntersectsPolygon(polygon)
}

// IntersectionPolygon returns the points where the ray crosses the polygon boundary, from Origin
// on, as Segment.IntersectionPolygon finds them on the reach of the ray past the polygon.
func (r Ray[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return r.AppendIntersectionPolygon(nil, polygon)
}

// AppendIntersectionPolygon appends the points IntersectionPolygon returns to dst and returns the
// extended slice, as Segment.AppendIntersectionPolygon does on the reach of the ray past the
// polygon, so a caller reusing dst allocates nothing once it has room.
func (r Ray[T]) AppendIntersectionPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	return r.reach(polygon.minMax()).AppendIntersectionPolygon(dst, polygon)
}

// IntersectsRectangle reports whether the ray and the rectangle share a point, as
// Segment.IntersectsRectangle decides it on the reach of the ray past the rectangle: the origin
// lies within the rectangle, or the ray crosses one of its edges. Touching shapes intersect,
// within the tolerance.
func (r Ray[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return r.reach(rectangle.MinMax()).IntersectsRectangle(rectangle)
}

// IntersectionRectangle returns the points where the ray crosses the rectangle boundary, from
// Origin on, as Segment.IntersectionRectangle finds them on the reach of the ray past the
// rectangle.
func (r Ray[T]) IntersectionRectangle(rectangle Rectangle[T]) []Point[T] {
	return r.AppendIntersectionRectangle(nil, rectangle)
}

// AppendIntersectionRectangle appends the points IntersectionRectangle returns to dst and returns
// the extended slice, as Segment.AppendIntersectionRectangle does on the reach of the ray past the
// rectangle, so a caller reusing dst allocates nothing once it has room.
func (r Ray[T]) AppendIntersectionRectangle(dst []Point[T], rectangle Rectangle[T]) []Point[T] {
	return r.reach(rectangle.MinMax()).AppendIntersectionRectangle(dst, rectangle)
}

// IntersectsRegularPolygon reports whether the ray and the regular polygon share a point, as
// Segment.IntersectsRegularPolygon decides it on the reach of the ray past the polygon: the
// origin lies within the polygon, or the ray crosses one of its edges. Touching shapes
// intersect, within the tolerance, and an empty polygon intersects nothing.
func (r Ray[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return r.reach(polygon.minMax()).IntersectsRegularPolygon(polygon)
}

// IntersectionRegularPolygon returns the points where the ray crosses the regular polygon
// boundary, from Origin on, as Segment.IntersectionRegularPolygon finds them on the reach of the
// ray past the polygon.
func (r Ray[T]) IntersectionRegularPolygon(polygon RegularPolygon[T]) []Point[T] {
	return r.AppendIntersectionRegularPolygon(nil, polygon)
}

// AppendIntersectionRegularPolygon appends the points IntersectionRegularPolygon returns to dst
// and returns the extended slice, as Segment.AppendIntersectionRegularPolygon does on the reach of
// the ray past the regular polygon, so a caller reusing dst allocates nothing once it has room.
func (r Ray[T]) AppendIntersectionRegularPolygon(dst []Point[T], polygon RegularPolygon[T]) []Point[T] {
	return r.reach(polygon.minMax()).AppendIntersectionRegularPolygon(dst, polygon)
}

// IntersectsBox reports whether the ray and the box share a point, as Segment.IntersectsBox
// decides it on the reach of the ray past the box: the origin lies within the box, or the ray
// crosses one of its edges. Touching shapes intersect, within the tolerance.
func (r Ray[T]) IntersectsBox(box Box[T]) bool {
	return r.reach(box.Min, box.Max).IntersectsBox(box)
}

// IntersectionBox returns the points where the ray crosses the box boundary, from Origin on, as
// Segment.IntersectionBox finds them on the reach of the ray past the box.
func (r Ray[T]) IntersectionBox(box Box[T]) []Point[T] {
	return r.AppendIntersectionBox(nil, box)
}

// AppendIntersectionBox appends the points IntersectionBox returns to dst and returns the
// extended slice, as Segment.AppendIntersectionBox does on the reach of the ray past the box, so
// a caller reusing dst allocates nothing once it has room.
func (r Ray[T]) AppendIntersectionBox(dst []Point[T], box Box[T]) []Point[T] {
	return r.reach(box.Min, box.Max).AppendIntersectionBox(dst, box)
}

// ClipCircle returns the part of the ray inside the circle, boundary included within the tolerance,
// and false where they share no point, as Segment.ClipCircle clips the reach of the ray past the
// circle: from Origin where the circle contains it, or else from the first point where the ray
// meets the circle, to the point where it leaves. Its Start is the cast of the ray, the first point
// of the circle it reaches. It allocates nothing.
func (r Ray[T]) ClipCircle(circle Circle[T]) (Segment[T], bool) {
	return r.reach(circle.minMax()).ClipCircle(circle)
}

// ClipPolygon returns the parts of the ray inside the polygon, boundary included within the
// tolerance, from Origin on, as Segment.ClipPolygon clips the reach of the ray past the polygon.
// The Start of the first part is the cast of the ray, the first point of the polygon it reaches. A
// ray along the line of a flat polygon has no part, as Segment.ClipPolygon says. The result is the
// one allocation, made on the first part.
func (r Ray[T]) ClipPolygon(polygon Polygon[T]) []Segment[T] {
	return r.AppendClipPolygon(nil, polygon)
}

// AppendClipPolygon appends the parts ClipPolygon returns to dst and returns the extended slice,
// as Segment.AppendClipPolygon does on the reach of the ray past the polygon, so a caller
// reusing dst allocates nothing once it has room.
func (r Ray[T]) AppendClipPolygon(dst []Segment[T], polygon Polygon[T]) []Segment[T] {
	return r.reach(polygon.minMax()).AppendClipPolygon(dst, polygon)
}

// ClipRectangle returns the part of the ray inside the rectangle, boundary included within
// the tolerance, whatever its angle, and false where they share no point, as
// Segment.ClipRectangle clips the reach of the ray past the rectangle. Its Start is the cast of
// the ray, the first point of the rectangle it reaches. It allocates nothing.
func (r Ray[T]) ClipRectangle(rectangle Rectangle[T]) (Segment[T], bool) {
	return r.reach(rectangle.MinMax()).ClipRectangle(rectangle)
}

// ClipRegularPolygon returns the part of the ray inside the regular polygon, boundary included
// within the tolerance, and false where they share no point, as Segment.ClipRegularPolygon clips
// the reach of the ray past the polygon. Its Start is the cast of the ray, the first point of the
// polygon it reaches, and an empty polygon clips everything away. A ray along the line of a
// flat polygon has no part, as Segment.ClipRegularPolygon says. It allocates nothing.
func (r Ray[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Segment[T], bool) {
	return r.reach(polygon.minMax()).ClipRegularPolygon(polygon)
}

// ClipBox returns the part of the ray inside the box, boundary included within the tolerance,
// and false where they share no point, as Segment.ClipBox clips the reach of the ray past the
// box. Its Start is the cast of the ray, the first point of the box it reaches. It allocates
// nothing.
func (r Ray[T]) ClipBox(box Box[T]) (Segment[T], bool) {
	return r.reach(box.Min, box.Max).ClipBox(box)
}

// reach returns the part of the ray from Origin to where it passes the extent a, b of a shape:
// to the point the farthest corner of the extent projects to, or Origin alone where the whole
// extent lies behind it or the ray has no direction. No point of the extent projects past the
// end, so the segment shares every point of the shape the ray does, and a pair of the ray is
// decided by the segment's own tests. For integer T the fraction is rounded up to a whole step
// of Direction, so the end is a lattice point on the ray rather than a rounded point beside it.
func (r Ray[T]) reach(a, b Point[T]) Segment[T] {
	if !r.Direction.hasDirection() {
		return Segment[T]{r.Origin, r.Origin}
	}

	direction := r.Direction.Float()
	corner := a
	if direction.X > 0 {
		corner.X = b.X
	}
	if direction.Y > 0 {
		corner.Y = b.Y
	}

	t := max(corner.Float().Subtract(r.Origin.Float()).Dot(direction)/direction.LengthSquared(), 0)
	if isInt[T]() {
		t = math.Ceil(t)
	}

	return Segment[T]{r.Origin, r.PointAt(t)}
}

// crosses reports whether the rays properly cross, as crossing decides it: the lines through
// them meet strictly ahead of both origins. Touching and parallel rays do not cross and are left
// to the origin distances.
func (r Ray[T]) crosses(ray Ray[T]) bool {
	_, ok := r.crossing(ray)

	return ok
}

// crossing returns the point where the lines through two properly crossing rays meet, in
// float64, and false where they do not properly cross. Each fraction is read from the sign of a
// cross product against the sign of the cross of the directions, with no division, so a NaN,
// which has neither sign, crosses nothing. For nearly collinear float rays rounding can give
// both fractions the sign of a crossing, and the point, divided out of cross products near zero,
// lands anywhere along this ray; a point whose projection onto the other falls behind its
// origin is therefore no crossing, and the pair is left to the origins, on the comparison foot
// makes.
//
// The point is placed from Origin by each component of the direction times the numerator of the
// fraction, divided by its denominator once, rather than by the rounded fraction: for an integer
// T the product and the cross products are exact, so a crossing on a half unit is exactly there
// and rounds the same way whichever ray asks, while that product stays within the integers
// float64 holds exactly. Past it the product rounds, so the pair is taken in one order whichever
// ray asks, the lesser by Point.Compare of Origin and then Direction first.
func (r Ray[T]) crossing(ray Ray[T]) (Point[float64], bool) {
	if cmp.Or(ray.Origin.Compare(r.Origin), ray.Direction.Point().Compare(r.Direction.Point())) < 0 {
		r, ray = ray, r
	}

	a, b := r.Direction.Float(), ray.Direction.Float()
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

// parallel reports whether the rays run along the same line direction, either way. A ray with a
// zero Direction has none and is parallel to nothing, so its origin can still be found on the
// other.
func (r Ray[T]) parallel(ray Ray[T]) bool {
	a, b := r.Direction.Float(), ray.Direction.Float()

	return a.hasDirection() && b.hasDirection() && a.Cross(b) == 0
}

// touch returns the origin of either ray that lies on the other, within the tolerance as Contains
// judges it, and false when there is none. Non-parallel rays that do not properly cross can
// share a point only this way.
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

// Equal checks if the origins and the directions of the rays are equal. Two rays through the
// same points with directions of different lengths are different values, as a vector is.
func (r Ray[T]) Equal(ray Ray[T]) bool {
	return r.Origin.Equal(ray.Origin) && r.Direction.Equal(ray.Direction)
}

// IsZero checks if the origin and the direction are zero.
func (r Ray[T]) IsZero() bool {
	return r.Equal(Ray[T]{})
}

// Cast converts the ray to a Ray of another number type, rounding as Cast does, so a direction
// shorter than a unit can round to zero for an integer R.
func (r Ray[T]) Cast[R Number]() Ray[R] {
	return Ray[R]{r.Origin.Cast[R](), r.Direction.Cast[R]()}
}

// Int converts the ray to a Ray[int].
func (r Ray[T]) Int() Ray[int] {
	return Ray[int]{r.Origin.Int(), r.Direction.Int()}
}

// Float converts the ray to a Ray[float64].
func (r Ray[T]) Float() Ray[float64] {
	return Ray[float64]{r.Origin.Float(), r.Direction.Float()}
}

// String returns the ray by its origin and direction: Ray((x,y);⟨x,y⟩).
func (r Ray[T]) String() string {
	return fmt.Sprintf("Ray(%s;%s)", r.Origin.String(), r.Direction.String())
}

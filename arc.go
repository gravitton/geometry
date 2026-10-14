package geom

import (
	"math"
	"slices"
)

// Arc is the part of the circle about Center from Angle through Sweep, both in radians, in the
// sense of Vector.Rotate for a positive sweep, so the direction it is drawn in is part of the
// value. A sweep of a full turn or more either way is the whole circle, and a sweep of zero the
// single point at Angle. It is a curve with no inside.
//
// The radius is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative radius written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it.
//
// The points of an arc are no lattice points, so for integer T every distance to it is measured in
// float64 and compared at no less than the float64 tolerance.
//
// An arc is no Shape, since its Support would be that of its hull; segments, rays and lines are
// cast at it.
type Arc[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius T        `json:"r"`
	Angle  float64  `json:"a,omitzero"`
	Sweep  float64  `json:"s"`
}

// Ac is shorthand for Arc{center, radius, angle, sweep}, with the radius taken absolute.
func Ac[T Number](center Point[T], radius T, angle, sweep float64) Arc[T] {
	return Arc[T]{center, Abs(radius), angle, sweep}
}

// ArcThrough returns the arc from a through b to c, and false where they lie on one line, two or
// all of them being the same point included.
func ArcThrough[T Number](a, b, c Point[T]) (Arc[T], bool) {
	circle, ok := CircleThrough(a, b, c)
	if !ok {
		return Arc[T]{}, false
	}

	center := circle.Center.Float()
	angle := a.Float().Subtract(center).Angle()
	via := NormalizeAngle(b.Float().Subtract(center).Angle() - angle)

	sweep := NormalizeAngle(c.Float().Subtract(center).Angle() - angle)
	if via > sweep {
		sweep -= 2 * math.Pi
	}

	return Arc[T]{circle.Center, circle.Radius, NormalizeAngle(angle), sweep}, true
}

// Start returns PointAt(0), the point at Angle.
func (a Arc[T]) Start() Point[T] {
	return a.PointAt(0)
}

// End returns PointAt(1), Equal to Start for the whole circle.
func (a Arc[T]) End() Point[T] {
	return a.PointAt(1)
}

// Midpoint returns PointAt(0.5).
func (a Arc[T]) Midpoint() Point[T] {
	return a.PointAt(0.5)
}

// Chord returns the segment from Start to End, its ends Equal for the whole circle.
func (a Arc[T]) Chord() Segment[T] {
	return Segment[T]{a.PointAt(0), a.PointAt(1)}
}

// Length returns the length of the arc (r·|sweep|), the sweep held to a full turn.
func (a Arc[T]) Length() float64 {
	return float64(a.Radius) * math.Abs(a.turn())
}

// Bounds returns the axis-aligned bounding box. For integer T the ends are rounded outward, except
// one within the tolerance of a whole coordinate, which takes it.
func (a Arc[T]) Bounds() Box[T] {
	a1, b1 := a.minMax()

	return Box[T]{a1, b1}
}

// minMax returns the corners of Bounds, from the ends and the axis points the sweep passes. For
// integer T they are rounded outward past the tolerance, since a sine or a cosine is that far off.
func (a Arc[T]) minMax() (Point[T], Point[T]) {
	center, radius := a.Center.Float(), float64(a.Radius)
	points := [6]Point[float64]{a.pointAt(0), a.pointAt(1)}

	n := 2
	for i, offset := range [4]Vector[float64]{{radius, 0}, {0, radius}, {-radius, 0}, {0, -radius}} {
		if a.covers(float64(float64(i) * math.Pi / 2)) {
			points[n] = center.Add(offset)
			n++
		}
	}

	a1, b1 := minMaxOf(points[:n])

	epsilon := 0.0
	if isInt[T]() {
		epsilon = a.epsilon(0)
	}

	return Point[T]{-castUp[T](-a1.X - epsilon), -castUp[T](-a1.Y - epsilon)}, Point[T]{castUp[T](b1.X - epsilon), castUp[T](b1.Y - epsilon)}
}

// magnitude returns the largest absolute coordinate of the circle of the arc, the magnitude the
// tolerance of a comparison with the arc is taken from.
func (a Arc[T]) magnitude() float64 {
	return a.Center.magnitude() + float64(a.Radius)
}

// turn returns the sweep held to a full turn either way.
func (a Arc[T]) turn() float64 {
	return Clamp(a.Sweep, -2*math.Pi, 2*math.Pi)
}

// pointAt returns PointAt in float64 before rounding, the point every distance to an end is
// measured to.
func (a Arc[T]) pointAt(t float64) Point[float64] {
	return a.Center.Float().Add(VectorFromAngle(a.Angle+float64(t*a.turn()), float64(a.Radius)))
}

// Translate returns the arc moved by the vector.
func (a Arc[T]) Translate(vector Vector[T]) Arc[T] {
	return Arc[T]{a.Center.Add(vector), a.Radius, a.Angle, a.Sweep}
}

// MoveTo returns the arc with its center at the point.
func (a Arc[T]) MoveTo(point Point[T]) Arc[T] {
	return Arc[T]{point, a.Radius, a.Angle, a.Sweep}
}

// Rotate returns the arc turned by the angle in radians about its center, the sweep unchanged and
// the angle normalized to [0, 2π).
func (a Arc[T]) Rotate(angle float64) Arc[T] {
	return Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + angle), a.Sweep}
}

// Scale returns the arc with its radius scaled by the factor, taken absolute.
func (a Arc[T]) Scale(factor float64) Arc[T] {
	return Arc[T]{a.Center, Abs(Multiply(a.Radius, factor)), a.Angle, a.Sweep}
}

// Unscale returns the arc with its radius divided by the factor, taken absolute. It panics for a
// zero factor.
func (a Arc[T]) Unscale(factor float64) Arc[T] {
	return Arc[T]{a.Center, Abs(Divide(a.Radius, factor)), a.Angle, a.Sweep}
}

// Resize returns the arc with the radius, taken absolute.
func (a Arc[T]) Resize(radius T) Arc[T] {
	return Arc[T]{a.Center, Abs(radius), a.Angle, a.Sweep}
}

// Grow returns the arc with its radius increased by the amount, clamped at zero.
func (a Arc[T]) Grow(amount T) Arc[T] {
	return Arc[T]{a.Center, max(a.Radius+amount, 0), a.Angle, a.Sweep}
}

// Shrink returns the arc with its radius decreased by the amount, clamped at zero.
func (a Arc[T]) Shrink(amount T) Arc[T] {
	return Arc[T]{a.Center, max(a.Radius-amount, 0), a.Angle, a.Sweep}
}

// Lerp returns the arc interpolated toward the given one, center, radius and sweep together and the
// angle along the shorter way, normalized to [0, 2π), and extrapolated outside [0, 1]. The radius
// is taken absolute.
func (a Arc[T]) Lerp(arc Arc[T], t float64) Arc[T] {
	return Arc[T]{a.Center.Lerp(arc.Center, t), Abs(Lerp(a.Radius, arc.Radius, t)), NormalizeAngle(LerpAngle(a.Angle, arc.Angle, t)), Lerp(a.Sweep, arc.Sweep, t)}
}

// Canonical returns the arc with its radius taken absolute, its angle normalized to [0, 2π) and its
// sweep held to a full turn. An angle within Delta of zero or a full turn becomes zero.
func (a Arc[T]) Canonical() Arc[T] {
	return Arc[T]{a.Center, Abs(a.Radius), snapAngle(a.Angle), a.turn()}
}

// Reverse returns the arc over the same points drawn the other way, from End back to Start.
func (a Arc[T]) Reverse() Arc[T] {
	return Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + a.turn()), -a.turn()}
}

// PointAt returns the point at the fraction t along the arc, Start at zero and End at one,
// continuing around the circle outside [0, 1].
func (a Arc[T]) PointAt(t float64) Point[T] {
	return a.pointAt(t).Cast[T]()
}

// TangentAt returns the derivative at the fraction t: the direction the arc is drawn in there, with
// the Length of the arc as its length, and zero for an arc of zero length.
func (a Arc[T]) TangentAt(t float64) Vector[float64] {
	turn := a.turn()

	return VectorFromAngle(a.Angle+float64(t*turn), float64(a.Radius)*turn).Normal()
}

// Split returns the arcs from Start to PointAt(t) and from there to End, the sweep first held to a
// full turn. Outside [0, 1] one runs past an end and the other back over it.
func (a Arc[T]) Split(t float64) (Arc[T], Arc[T]) {
	turn := a.turn()
	first := float64(t * turn)

	return Arc[T]{a.Center, a.Radius, a.Angle, first}, Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + first), turn - first}
}

// DistanceTo returns the distance from the point to the arc, zero on it.
func (a Arc[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(a.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (a Arc[T]) DistanceSquaredTo(point Point[T]) float64 {
	_, distanceSquared := a.nearest(point.Float())
	if a.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return distanceSquared
}

// Nearest returns the point of the arc nearest to the given point, the point itself on the arc or
// with a NaN coordinate.
func (a Arc[T]) Nearest(point Point[T]) Point[T] {
	nearest, distanceSquared := a.nearest(point.Float())
	if a.containsSquared(distanceSquared, point.magnitude()) || math.IsNaN(distanceSquared) {
		return point
	}

	return nearest.Cast[T]()
}

// NearestFraction returns the fraction along the arc of Nearest, from zero at Start to one at End,
// and zero for an arc of zero sweep.
func (a Arc[T]) NearestFraction(point Point[T]) float64 {
	return a.fraction(point.Float())
}

// Contains reports whether the point lies on the arc. The tolerance is on the distance, never on an
// angle, so it does not widen far from the center.
func (a Arc[T]) Contains(point Point[T]) bool {
	return a.DistanceSquaredTo(point) == 0
}

// IntersectionSegment returns the points where the segment crosses the arc, from Start to End.
func (a Arc[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	fractions, n := a.crossings(segment)

	return segment.pointsAt(fractions[:n]...)
}

// IntersectionRay returns the points where the ray crosses the arc, from its Origin on.
func (a Arc[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return a.IntersectionSegment(ray.reach(a.minMax()))
}

// IntersectionLine returns the points where the line crosses the arc, in order along its Vector.
func (a Arc[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return a.IntersectionSegment(line.span(a.minMax()))
}

// crossings returns the fractions along the segment, ascending and clamped to it, at which it
// crosses the circle on the arc, and how many there are. IntersectionSegment and Segment.HitArc
// place points from them.
func (a Arc[T]) crossings(segment Segment[T]) ([2]float64, int) {
	cast := segment.Float()
	crossed, count := a.Circle().fractions(segment)

	var fractions [2]float64

	n := 0
	for _, t := range crossed[:count] {
		if !segment.containsAt(t) {
			continue
		}

		t = Clamp(t, 0, 1)
		if a.holds(cast.PointAt(t)) {
			fractions[n] = t
			n++
		}
	}

	return fractions, n
}

// holds reports whether the float64 point lies on the arc, the test Contains makes before rounding
// into T.
func (a Arc[T]) holds(point Point[float64]) bool {
	_, distanceSquared := a.nearest(point)

	return a.containsSquared(distanceSquared, point.magnitude())
}

// nearest returns the point of the arc nearest to the given point and the squared distance to it,
// with no tolerance, so off the arc Nearest is exactly PointAt of NearestFraction.
func (a Arc[T]) nearest(point Point[float64]) (Point[float64], float64) {
	nearest := a.pointAt(a.fraction(point))

	return nearest, nearest.DistanceSquaredTo(point)
}

// fraction returns the fraction along the arc nearest to the point, with no tolerance: the swept
// angle where the sweep covers the direction, else the nearer end, Start where both are exactly as
// near.
func (a Arc[T]) fraction(point Point[float64]) float64 {
	turn := math.Abs(a.turn())

	if offset := point.Subtract(a.Center.Float()); offset.hasDirection() {
		if swept := a.swept(offset.Angle()); swept <= turn {
			return ratio(swept, turn)
		}
	}

	if a.pointAt(1).DistanceSquaredTo(point) < a.pointAt(0).DistanceSquaredTo(point) {
		return 1
	}

	return 0
}

// covers reports whether the sweep passes the angle, with no tolerance, since every caller measures
// a distance that is the same either side of an end.
func (a Arc[T]) covers(angle float64) bool {
	return a.swept(angle) <= math.Abs(a.turn())
}

// swept returns how far the angle lies from Angle the way the arc is drawn, in [0, 2π).
func (a Arc[T]) swept(angle float64) float64 {
	if a.turn() < 0 {
		return NormalizeAngle(a.Angle - angle)
	}

	return NormalizeAngle(angle - a.Angle)
}

// containsSquared reports whether a point at the squared distance from the arc lies on it, within
// the tolerance.
func (a Arc[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, 0, a.epsilon(magnitude))
}

// epsilon returns the tolerance of a comparison between the arc and geometry of the given
// magnitude, never less than that of float64, since the points of the arc are no lattice points.
func (a Arc[T]) epsilon(magnitude float64) float64 {
	magnitude = max(a.magnitude(), magnitude)

	return max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude))
}

// Equal reports whether the arcs have equal center, radius, angle and sweep. The angles are
// compared with EqualAngle, the sweeps as numbers, so the same points drawn opposite ways are not
// equal.
func (a Arc[T]) Equal(arc Arc[T]) bool {
	return a.Center.Equal(arc.Center) && Equal(a.Radius, arc.Radius) && EqualAngle(a.Angle, arc.Angle) && Equal(a.Sweep, arc.Sweep)
}

// IsZero reports whether the arc is the zero value.
func (a Arc[T]) IsZero() bool {
	return a.Equal(Arc[T]{})
}

// Circle returns the circle the arc is part of.
func (a Arc[T]) Circle() Circle[T] {
	return Circle[T]{a.Center, a.Radius}
}

// Sector returns the sector the arc bounds with the radii to its ends.
func (a Arc[T]) Sector() Sector[T] {
	return Sector[T](a)
}

// Polyline returns the arc as a polyline within the tolerance of it, from Start to End: as few
// equal chords as keep within it, at most 1024. A negative tolerance is taken absolute.
func (a Arc[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{a.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of Polyline(tolerance) to dst and returns the extended slice.
func (a Arc[T]) AppendPolyline(dst []Point[T], tolerance float64) []Point[T] {
	return a.appendPieces(dst, a.pieces(tolerance))
}

// Cast converts the arc to an Arc of another number type, rounding as Cast does and keeping the
// angle and the sweep.
func (a Arc[T]) Cast[R Number]() Arc[R] {
	return Arc[R]{a.Center.Cast[R](), Cast[R](float64(a.Radius)), a.Angle, a.Sweep}
}

// Int converts the arc to an Arc[int].
func (a Arc[T]) Int() Arc[int] {
	return Arc[int]{a.Center.Int(), Int(a.Radius), a.Angle, a.Sweep}
}

// Float converts the arc to an Arc[float64].
func (a Arc[T]) Float() Arc[float64] {
	return Arc[float64]{a.Center.Float(), float64(a.Radius), a.Angle, a.Sweep}
}

// appendPieces appends the ends of the given number of equal chords from Start to End, leaving out
// a point that rounds onto the one before. Sector.AppendPolygon raises the count above what the
// tolerance gives.
func (a Arc[T]) appendPieces(dst []Point[T], pieces float64) []Point[T] {
	dst = append(slices.Grow(dst, int(pieces)+1), a.PointAt(0))
	for i := 1.0; i <= pieces; i++ {
		if point := a.PointAt(i / pieces); point != dst[len(dst)-1] {
			dst = append(dst, point)
		}
	}

	return dst
}

// pieces returns how many equal chords keep the arc within the tolerance, at least one and at most
// two to the flattenDepth, which a zero or NaN tolerance gives.
func (a Arc[T]) pieces(tolerance float64) float64 {
	chord := 2 * math.Acos(1-min(ratio(math.Abs(tolerance), float64(a.Radius)), 1))

	pieces := max(math.Ceil(math.Abs(a.turn())/chord), 1)
	if !(pieces <= 1<<flattenDepth) {
		pieces = 1 << flattenDepth
	}

	return pieces
}

// String returns the arc by its center, radius, angle and sweep: Arc((x,y);r;a;s).
func (a Arc[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Arc("...)
	buf = a.Center.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, a.Radius)
	buf = append(buf, ';')
	buf = appendString(buf, a.Angle)
	buf = append(buf, ';')
	buf = appendString(buf, a.Sweep)
	buf = append(buf, ')')

	return string(buf)
}

package geom

import (
	"math"
	"slices"
)

// Arc is a 2D circular arc: the part of the circle about Center with the given Radius that
// starts at Angle and runs through Sweep, both in radians, in the sense of Vector.Rotate for a
// positive sweep and against it for a negative one, so the direction it is drawn in is part of
// the value. A sweep of a full turn or more, either way, is the whole circle, and a sweep of
// zero the single point at Angle. It is a curve with no inside; Sector is the area it bounds
// with the two radii to its ends.
//
// The radius is never negative: Ac, Resize and Lerp take it absolute, Scale takes a negative
// factor absolute, and Grow and Shrink clamp at zero. A negative radius can only be written as a struct literal or decoded from JSON, and
// Contains, DistanceTo and Bounds give no meaningful answer for it; Canonical repairs it.
//
// The angles place the points of the arc off the lattice, so for an integer T every distance to
// the arc is measured in float64 and compared at the tolerance of float64, where the other
// shapes compare an integer exactly.
//
// The arc is no Shape: like a Polyline it has no Support that would not hand the functions over
// two shapes its hull. It is a curved wall a segment, a ray or a line is cast at, by their
// Intersection and Hit methods naming it.
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

// ArcThrough creates the Arc from the first point through the second to the third, a part of
// the circle CircleThrough gives for them, swept the way that passes the second on its way,
// and reports false where there is none: the points lie on one line within the tolerance, two
// or all of them being the same point included. For integer T the center and the radius are
// rounded, so the points can lie off the arc.
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

// Length returns the length of the arc, the radius times the absolute sweep, at most the
// circumference.
func (a Arc[T]) Length() float64 {
	return float64(a.Radius) * math.Abs(a.turn())
}

// Start returns the point the arc starts at, PointAt(0): the point of the circle at Angle.
func (a Arc[T]) Start() Point[T] {
	return a.PointAt(0)
}

// End returns the point the arc ends at, PointAt(1): the point of the circle at Angle plus
// Sweep, which is Start again for the whole circle.
func (a Arc[T]) End() Point[T] {
	return a.PointAt(1)
}

// Midpoint returns the point halfway along the arc, PointAt(0.5).
func (a Arc[T]) Midpoint() Point[T] {
	return a.PointAt(0.5)
}

// Chord returns the segment from Start to End, the straight line the arc bows away from. The
// whole circle has a chord of no length.
func (a Arc[T]) Chord() Segment[T] {
	return Segment[T]{a.PointAt(0), a.PointAt(1)}
}

// Bounds returns the axis-aligned bounding box: the box around the two ends and each of the
// four points of the circle along the axes that the sweep passes. For integer T the ends are
// rounded outward, so the box holds the arc, less a rounding of the angle within the tolerance.
func (a Arc[T]) Bounds() Box[T] {
	a1, b1 := a.minMax()

	return Box[T]{a1, b1}
}

// minMax returns the minimum and maximum corner of the arc, the corners of Bounds: those of its
// two ends and of each point of the circle along an axis the sweep passes, measured in float64
// and rounded outward for an integer T, an end within the tolerance of a whole coordinate
// taking it, since a sine or a cosine is that far off.
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

// magnitude returns the largest absolute coordinate of the circle the arc is part of: the size
// epsilonAt widens the tolerance by for a comparison that reads the arc.
func (a Arc[T]) magnitude() float64 {
	return a.Center.magnitude() + float64(a.Radius)
}

// turn returns the sweep the arc is drawn through, held to a full turn either way: the whole
// circle has no more of itself to run over.
func (a Arc[T]) turn() float64 {
	return Clamp(a.Sweep, -2*math.Pi, 2*math.Pi)
}

// pointAt returns the point at the fraction t of the way along the arc in float64, before it
// is rounded for an integer T: the point every distance to an end is measured to.
func (a Arc[T]) pointAt(t float64) Point[float64] {
	return a.Center.Float().Add(VectorFromAngle(a.Angle+float64(t*a.turn()), float64(a.Radius)))
}

// Translate creates a new Arc translated by the given vector.
func (a Arc[T]) Translate(vector Vector[T]) Arc[T] {
	return Arc[T]{a.Center.Add(vector), a.Radius, a.Angle, a.Sweep}
}

// MoveTo creates a new Arc with the same radius, angle and sweep and the center set to point.
func (a Arc[T]) MoveTo(point Point[T]) Arc[T] {
	return Arc[T]{point, a.Radius, a.Angle, a.Sweep}
}

// Rotate creates a new Arc turned by the given angle (in radians) about its center, in the
// same sense as Vector.Rotate: the angle it starts at moves and the sweep stays. The stored
// angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (a Arc[T]) Rotate(angle float64) Arc[T] {
	return Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + angle), a.Sweep}
}

// Scale creates a new Arc with radius scaled by the given factor, about its center. A negative
// factor scales by its absolute value.
func (a Arc[T]) Scale(factor float64) Arc[T] {
	return Arc[T]{a.Center, Abs(Multiply(a.Radius, factor)), a.Angle, a.Sweep}
}

// Unscale creates a new Arc with radius scaled by the inverse factor, the inverse of Scale, a
// negative factor by its absolute value like it. Like Divide it panics for a zero factor.
func (a Arc[T]) Unscale(factor float64) Arc[T] {
	return Arc[T]{a.Center, Abs(Divide(a.Radius, factor)), a.Angle, a.Sweep}
}

// Resize creates a new Arc with the given radius, taken absolute like Ac.
func (a Arc[T]) Resize(radius T) Arc[T] {
	return Arc[T]{a.Center, Abs(radius), a.Angle, a.Sweep}
}

// Grow creates a new Arc with radius increased by amount, clamped to zero: the curve running
// beside this one the amount farther from the center, over the same angles.
func (a Arc[T]) Grow(amount T) Arc[T] {
	return Arc[T]{a.Center, max(a.Radius+amount, 0), a.Angle, a.Sweep}
}

// Shrink creates a new Arc with radius decreased by amount, clamped to zero.
func (a Arc[T]) Shrink(amount T) Arc[T] {
	return Arc[T]{a.Center, max(a.Radius-amount, 0), a.Angle, a.Sweep}
}

// Lerp creates a new Arc in linear interpolation towards the given arc, moving the center, the
// radius and the sweep together and turning the angle along the shorter arc with LerpAngle,
// normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like Point.Lerp, with the
// radius taken absolute.
func (a Arc[T]) Lerp(arc Arc[T], t float64) Arc[T] {
	return Arc[T]{a.Center.Lerp(arc.Center, t), Abs(Lerp(a.Radius, arc.Radius, t)), NormalizeAngle(LerpAngle(a.Angle, arc.Angle, t)), Lerp(a.Sweep, arc.Sweep, t)}
}

// Canonical creates a new Arc in the form Ac and Rotate build, with the radius taken absolute,
// the angle normalized to [0, 2π) and a sweep past a full turn held to it: a well-formed arc is
// returned as it is, up to the full turns Equal already ignores in the angle. An angle within
// Delta of zero or of a full turn becomes exactly zero.
func (a Arc[T]) Canonical() Arc[T] {
	return Arc[T]{a.Center, Abs(a.Radius), snapAngle(a.Angle), a.turn()}
}

// Reverse creates a new Arc over the same points drawn the other way: it starts where this one
// ends and sweeps back.
func (a Arc[T]) Reverse() Arc[T] {
	return Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + a.turn()), -a.turn()}
}

// PointAt returns the point at the fraction t of the way along the arc, from Start at zero to
// End at one, extrapolating around the circle outside [0, 1] like Segment.PointAt along its
// line. For integer T the point is rounded once and can land off the arc, where Contains
// rejects it.
func (a Arc[T]) PointAt(t float64) Point[T] {
	return a.pointAt(t).Cast[T]()
}

// TangentAt returns the derivative of the arc at the fraction t, the direction it runs in there
// with the speed of the fraction as its length, the Length of the arc, in float64 whatever T
// is: a quarter turn from the radius to that point, the way the arc is drawn. It is the zero
// vector for an arc of no length.
func (a Arc[T]) TangentAt(t float64) Vector[float64] {
	turn := a.turn()

	return VectorFromAngle(a.Angle+float64(t*turn), float64(a.Radius)*turn).Normal()
}

// Split returns the two arcs the arc is cut into at the fraction t, the first from Start to
// PointAt(t) and the second from there to End, which together trace the same points the same
// way. A sweep past a full turn is held to it first. Outside [0, 1] one of the two runs past an
// end and the other back over it, each held to a full turn in its own right.
func (a Arc[T]) Split(t float64) (Arc[T], Arc[T]) {
	turn := a.turn()
	first := float64(t * turn)

	return Arc[T]{a.Center, a.Radius, a.Angle, first}, Arc[T]{a.Center, a.Radius, NormalizeAngle(a.Angle + first), turn - first}
}

// DistanceTo returns the distance from the given point to the nearest point of the arc: zero
// exactly where Contains holds, and otherwise the distance to the circle where the direction
// from the center to the point lies within the sweep, and to the nearer end where it does not.
func (a Arc[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(a.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. A distance within the tolerance is snapped to zero, which is where Contains
// reads it. It is a float64 even for an integer T, since the nearest point of an arc is not a
// lattice point in general.
func (a Arc[T]) DistanceSquaredTo(point Point[T]) float64 {
	_, distanceSquared := a.nearest(point.Float())
	if a.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return distanceSquared
}

// Nearest returns the point of the arc nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the point of the circle toward it from the center where
// that lies within the sweep, and the nearer end, Start of two as near, where it does not. For
// integer T it is rounded once and can land off the arc, where Contains rejects it. A point
// with a NaN coordinate has no nearest point and is returned as it is.
func (a Arc[T]) Nearest(point Point[T]) Point[T] {
	nearest, distanceSquared := a.nearest(point.Float())
	if a.containsSquared(distanceSquared, point.magnitude()) || math.IsNaN(distanceSquared) {
		return point
	}

	return nearest.Cast[T]()
}

// NearestFraction returns the fraction of the way along the arc at which the point of it
// nearest to the given point lies, the one PointAt takes back to it: from zero at Start to one
// at End, as Nearest decides between the arc and its ends. An arc of no length gives zero.
func (a Arc[T]) NearestFraction(point Point[T]) float64 {
	return a.fraction(point.Float())
}

// Contains reports whether the given point lies on the arc, within the tolerance: it holds
// exactly where DistanceTo is zero. The tolerance is on the distance to the arc, never on an
// angle, so a point a rounding error past an end is still on it.
func (a Arc[T]) Contains(point Point[T]) bool {
	return a.DistanceSquaredTo(point) == 0
}

// IntersectionSegment returns the points where the segment crosses the arc, from its Start to
// its End: the points it crosses the circle at that lie on the arc, within the tolerance, so
// two, one or none, placed as Circle.IntersectionSegment places its own. They are found in
// float64 and rounded once for an integer T.
//
// The points are a Pair, which allocates nothing.
func (a Arc[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	fractions, n := a.crossings(segment)

	return segment.pointsAt(fractions[:n]...)
}

// IntersectionRay returns the points where the ray crosses the arc, from its Origin on, as
// IntersectionSegment finds them on the reach of the ray past the arc.
func (a Arc[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return a.IntersectionSegment(ray.reach(a.minMax()))
}

// IntersectionLine returns the points where the line crosses the arc, in order along its
// Direction, as IntersectionSegment finds them on the span of the line past the arc.
func (a Arc[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return a.IntersectionSegment(line.span(a.minMax()))
}

// crossings returns the fractions along the given segment at which it crosses the arc, in
// increasing order and held to the segment, and how many of the two there are: those
// Circle.fractions gives for the circle of the arc, at the tolerance of T, that lie on the arc
// within it. IntersectionSegment and Segment.HitArc place their points from them.
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

// holds reports whether the given point in float64 lies on the arc within the tolerance: the
// comparison Contains makes, on a point not yet rounded into T.
func (a Arc[T]) holds(point Point[float64]) bool {
	_, distanceSquared := a.nearest(point)

	return a.containsSquared(distanceSquared, point.magnitude())
}

// nearest returns the point of the arc nearest to the given point, in float64, with the squared
// distance to it and no tolerance applied: the point at the fraction the given one is nearest
// at, so Nearest is PointAt of NearestFraction to the last bit.
func (a Arc[T]) nearest(point Point[float64]) (Point[float64], float64) {
	nearest := a.pointAt(a.fraction(point))

	return nearest, nearest.DistanceSquaredTo(point)
}

// fraction returns the fraction of the way along the arc at which it is nearest to the given
// point, with no tolerance applied: how far the direction from the center to the point has
// swept from Angle where the sweep covers it, and the nearer end otherwise, Start of two as
// near. The two agree where the direction is that of an end, so the sweep is judged on the
// angle itself. The center has no direction and is as far from every point of the arc, so it
// is given Start.
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

// covers reports whether the sweep passes the given angle: the angle has swept no farther from
// Angle than the sweep. It applies no tolerance, since every caller measures a distance that is
// the same either side of an end.
func (a Arc[T]) covers(angle float64) bool {
	return a.swept(angle) <= math.Abs(a.turn())
}

// swept returns how far the given angle lies from Angle the way the arc is drawn, in [0, 2π).
func (a Arc[T]) swept(angle float64) float64 {
	if a.turn() < 0 {
		return NormalizeAngle(a.Angle - angle)
	}

	return NormalizeAngle(angle - a.Angle)
}

// containsSquared reports whether a point at the given squared distance from the arc lies on
// it within the tolerance of the arc and of geometry whose largest absolute coordinate is the
// given magnitude: the snap DistanceSquaredTo and Nearest apply.
func (a Arc[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, 0, a.epsilon(magnitude))
}

// epsilon returns the tolerance of a comparison that reads the arc and geometry whose largest
// absolute coordinate is the given magnitude: that of T, and no less than that of float64,
// which an integer T is held to as well, since the points of the arc are no lattice points and
// every distance to it is measured in float64.
func (a Arc[T]) epsilon(magnitude float64) float64 {
	magnitude = max(a.magnitude(), magnitude)

	return max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude))
}

// Equal checks for equal center, radius, angle and sweep with given arc. The angles are
// compared with EqualAngle, so a full turn does not matter, and the sweeps as numbers, so two
// arcs over the same points drawn opposite ways are not equal.
func (a Arc[T]) Equal(arc Arc[T]) bool {
	return a.Center.Equal(arc.Center) && Equal(a.Radius, arc.Radius) && EqualAngle(a.Angle, arc.Angle) && Equal(a.Sweep, arc.Sweep)
}

// IsZero checks if center, radius, angle and sweep are zero.
func (a Arc[T]) IsZero() bool {
	return a.Equal(Arc[T]{})
}

// Circle converts the arc to the Circle it is part of.
func (a Arc[T]) Circle() Circle[T] {
	return Circle[T]{a.Center, a.Radius}
}

// Polyline converts the arc to a Polyline within the given tolerance of it, from Start to End:
// equal chords, as few as keep the arc within the tolerance of each, the path that draws the
// arc. A negative tolerance is taken absolute, and the path has at most a thousand edges, which
// a tolerance of zero gives. For integer T the points are rounded and one rounding onto the
// point before it is left out.
func (a Arc[T]) Polyline(tolerance float64) Polyline[T] {
	return Polyline[T]{a.AppendPolyline(nil, tolerance)}
}

// AppendPolyline appends the points of the path Polyline returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them.
func (a Arc[T]) AppendPolyline(dst []Point[T], tolerance float64) []Point[T] {
	return a.appendPieces(dst, a.pieces(tolerance))
}

// Sector converts the arc to the Sector it bounds with the radii to its two ends.
func (a Arc[T]) Sector() Sector[T] {
	return Sector[T](a)
}

// Cast converts the arc to an Arc of another number type, rounding as Cast does and keeping
// the angle and the sweep.
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

// appendPieces appends the ends of the given number of equal chords of the arc, from Start to
// End, a point rounding onto the one before it left out: AppendPolyline with the count already
// taken, which a Sector raises where the tolerance alone would leave it flat.
func (a Arc[T]) appendPieces(dst []Point[T], pieces float64) []Point[T] {
	dst = append(slices.Grow(dst, int(pieces)+1), a.PointAt(0))
	for i := 1.0; i <= pieces; i++ {
		if point := a.PointAt(i / pieces); point != dst[len(dst)-1] {
			dst = append(dst, point)
		}
	}

	return dst
}

// pieces returns how many equal chords the arc is flattened to within the given tolerance: as
// few as keep it within the tolerance of each, at least one and at most two to the
// flattenDepth, which a tolerance of zero and one that is no number give.
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

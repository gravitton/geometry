package geom

import (
	"math"
	"slices"
)

// Capsule is a 2D capsule: every point within Radius of Segment, a rectangle along the segment
// closed by half a disc at each end, the rounded collider of a character. A segment of no
// length makes it the circle about that point, and a radius of zero the segment itself.
//
// The radius is never negative: Cap and Lerp take it absolute, Scale takes a negative factor
// absolute, and Grow and Shrink clamp at zero. A negative radius can only be written as a struct literal
// or decoded from JSON, and Contains, DistanceTo and Bounds give no meaningful answer for it;
// Canonical repairs it.
//
// The capsule is a Shape and a Body, and one of the kinds the package names: it has the
// Intersects and Contains methods for every other, each decided on the distance to its segment,
// and a segment, a ray and a line are cast at it by their Hit, Intersection and Clip methods
// naming it. It is parted and moved by its Support, through Collide and Sweep. It has no
// vertices, and no Transform, since a matrix that scales one axis more than the other does not
// take it to a capsule.
type Capsule[T Number] struct {
	Segment Segment[T] `json:",embed"`
	Radius  T          `json:"r"`
}

// Cap is shorthand for Capsule{Segment{start, end}, radius}, with the radius taken absolute.
func Cap[T Number](start, end Point[T], radius T) Capsule[T] {
	return Capsule[T]{Segment[T]{start, end}, Abs(radius)}
}

// Centroid returns the center of the enclosed area, the Midpoint of the segment, rounded as it
// is for an integer T.
func (c Capsule[T]) Centroid() Point[T] {
	return c.Segment.Midpoint()
}

// Area returns the enclosed area: the rectangle along the segment, as wide as the diameter, and
// the two half discs that close it. It is a float64 even for an integer T, like the area of a
// Circle.
func (c Capsule[T]) Area() float64 {
	radius := float64(c.Radius)

	return radius * (float64(Pi*radius) + 2*c.Segment.Length())
}

// Perimeter returns the length of the boundary: the two sides along the segment and the two
// half circles around its ends.
func (c Capsule[T]) Perimeter() float64 {
	radius := float64(c.Radius)

	return 2 * (float64(Pi*radius) + c.Segment.Length())
}

// Inertia returns the polar second moment of area about the Centroid, the rotational inertia of
// the enclosed area at unit density: that of the rectangle along the segment and of the two half
// discs, each moved out to its end. It is that of a Circle for a segment of no length, and zero
// for a radius of zero, which encloses no area.
func (c Capsule[T]) Inertia() float64 {
	radius, length := float64(c.Radius), c.Segment.Length()
	rectangle := float64(length*length*length/6) + float64(2*radius*radius*length)
	discs := float64(Pi*radius*radius*radius/2) + float64(Pi*radius*length*length/4)

	return radius * (rectangle + discs)
}

// Bounds returns the axis-aligned bounding box: the Bounds of the segment, grown by the radius
// on every side.
func (c Capsule[T]) Bounds() Box[T] {
	a, b := c.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the smallest circle around the capsule: the BoundingCircle of the
// segment with the radius added, which stays the same size as the capsule turns. For integer T
// the midpoint is rounded and the radius rounded up from it, so the circle still contains the
// capsule.
func (c Capsule[T]) BoundingCircle() Circle[T] {
	circle := c.Segment.BoundingCircle()

	return Circle[T]{circle.Center, circle.Radius + c.Radius}
}

// minMax returns the minimum and maximum corner of the capsule, the corners of Bounds: those of
// the segment moved out by the radius.
func (c Capsule[T]) minMax() (Point[T], Point[T]) {
	a, b := c.Segment.minMax()

	return a.AddXY(-c.Radius, -c.Radius), b.AddXY(c.Radius, c.Radius)
}

// magnitude returns the largest absolute coordinate of the capsule, that of a corner of its
// Bounds: the size epsilonAt widens the tolerance by for a comparison that reads the capsule.
func (c Capsule[T]) magnitude() float64 {
	return c.Segment.magnitude() + float64(c.Radius)
}

// Translate creates a new Capsule translated by the given vector.
func (c Capsule[T]) Translate(vector Vector[T]) Capsule[T] {
	return Capsule[T]{c.Segment.Translate(vector), c.Radius}
}

// MoveTo creates a new Capsule with the midpoint of its segment moved to point and the same
// length, direction and radius, as Segment.MoveTo moves the segment.
func (c Capsule[T]) MoveTo(point Point[T]) Capsule[T] {
	return Capsule[T]{c.Segment.MoveTo(point), c.Radius}
}

// Rotate creates a new Capsule turned by the given angle (in radians) about the midpoint of its
// segment, in the same sense as Vector.Rotate, as Segment.Rotate turns the segment. For integer
// T the ends are rounded; only multiples of 90° keep the length exactly.
func (c Capsule[T]) Rotate(angle float64) Capsule[T] {
	return Capsule[T]{c.Segment.Rotate(angle), c.Radius}
}

// Scale creates a new Capsule uniformly scaled about the midpoint of its segment by the factor,
// the length and the radius alike. A negative factor swaps the ends of the segment and scales
// the radius by its absolute value, which is the same capsule as that of the positive factor.
func (c Capsule[T]) Scale(factor float64) Capsule[T] {
	return Capsule[T]{c.Segment.Scale(factor), Abs(Multiply(c.Radius, factor))}
}

// Unscale creates a new Capsule scaled about the midpoint of its segment by the inverse factor,
// the inverse of Scale, a negative factor by its absolute value like it. Like Divide it panics
// for a zero factor.
func (c Capsule[T]) Unscale(factor float64) Capsule[T] {
	return Capsule[T]{c.Segment.Unscale(factor), Abs(Divide(c.Radius, factor))}
}

// Grow creates a new Capsule with radius increased by amount, clamped to zero: the boundary
// moved out by the amount all around, the segment staying as it is.
func (c Capsule[T]) Grow(amount T) Capsule[T] {
	return Capsule[T]{c.Segment, max(c.Radius+amount, 0)}
}

// Shrink creates a new Capsule with radius decreased by amount, clamped to zero, where the
// capsule is its segment.
func (c Capsule[T]) Shrink(amount T) Capsule[T] {
	return Capsule[T]{c.Segment, max(c.Radius-amount, 0)}
}

// Lerp creates a new Capsule in linear interpolation towards the given capsule, moving the two
// ends and the radius together, and extrapolating outside [0, 1] like Point.Lerp, with the
// radius taken absolute like Cap: an extrapolation past a zero radius grows the capsule again.
func (c Capsule[T]) Lerp(capsule Capsule[T], t float64) Capsule[T] {
	return Capsule[T]{c.Segment.Lerp(capsule.Segment, t), Abs(Lerp(c.Radius, capsule.Radius, t))}
}

// Canonical creates a new Capsule in the form Cap builds, with the radius taken absolute: a
// well-formed capsule is returned as it is. It repairs a negative radius written as a struct
// literal or decoded from JSON before Contains, DistanceTo or Bounds read it.
func (c Capsule[T]) Canonical() Capsule[T] {
	return Capsule[T]{c.Segment, Abs(c.Radius)}
}

// DistanceTo returns the distance from the given point to the nearest point of the capsule:
// zero exactly where Contains holds, so a point within the tolerance of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the distance
// to the segment less the radius.
func (c Capsule[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.Segment.DistanceSquaredTo(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return math.Sqrt(distanceSquared) - float64(c.Radius)
}

// DistanceSquaredTo returns the square of DistanceTo, so every shape offers the same pair. It
// is a float64 even for an integer T and saves no square root, since the distance to a capsule
// already needs the root of the distance to its segment.
func (c Capsule[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := c.DistanceTo(point)

	return distance * distance
}

// DistanceToCapsule returns the distance between the nearest points of the two capsules: zero
// exactly where IntersectsCapsule holds, by the comparison it makes, and otherwise the distance
// between the two segments less the two radii.
func (c Capsule[T]) DistanceToCapsule(capsule Capsule[T]) float64 {
	if c.IntersectsCapsule(capsule) {
		return 0
	}

	return c.Segment.DistanceToSegment(capsule.Segment) - float64(c.Radius) - float64(capsule.Radius)
}

// DistanceSquaredToCapsule returns the square of DistanceToCapsule, so the pair is offered as
// every distance is. It saves no square root, since the distance between two capsules already
// needs the root of the distance between their segments.
func (c Capsule[T]) DistanceSquaredToCapsule(capsule Capsule[T]) float64 {
	distance := c.DistanceToCapsule(capsule)

	return distance * distance
}

// Nearest returns the point of the capsule nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the point of the boundary toward it from the nearest
// point of the segment, on the squared distance Contains compares. For integer T it is rounded
// once and can land off the boundary, where Contains rejects it.
func (c Capsule[T]) Nearest(point Point[T]) Point[T] {
	distanceSquared := c.Segment.DistanceSquaredTo(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return point
	}

	foot := c.Segment.Float().foot(point.Float())
	offset := point.Float().Subtract(foot).Multiply(float64(c.Radius) / math.Sqrt(distanceSquared))

	return foot.Add(offset).Cast[T]()
}

// Support returns the point of the capsule farthest along the given direction: the Support of
// the segment moved by the radius that way, whatever the length of the direction, and along X
// from Start for a direction of no length. It is a float64 point even for an integer T, since
// it is no lattice point in general.
func (c Capsule[T]) Support(direction Vector[float64]) Point[float64] {
	return c.Segment.Support(direction).Add(direction.Normalize().Multiply(float64(c.Radius)))
}

// Contains reports whether the given point lies within the capsule, boundary included within
// the tolerance: its distance to the segment is at most the radius, by the comparison a Circle
// makes against a segment, on the squared distance.
func (c Capsule[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.Segment.DistanceSquaredTo(point), point.magnitude())
}

// ContainsCircle reports whether the given circle lies within the capsule: its far point, the
// point of it farthest from the segment, is contained within the tolerance by the comparison
// Contains makes, so a circle touching the boundary from inside is enclosed.
func (c Capsule[T]) ContainsCircle(circle Circle[T]) bool {
	far := c.Segment.DistanceTo(circle.Center) + float64(circle.Radius)

	return c.containsSquared(far*far, circle.magnitude())
}

// ContainsCapsule reports whether the given capsule lies within the capsule: the circles of its
// radius about both ends of its segment do, as ContainsCircle decides it, and a capsule holds
// every point between points it contains.
func (c Capsule[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return c.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && c.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the given sector lies within the capsule: its center, the two
// ends of its arc and the points of it farthest from the segment are contained, within the
// tolerance, those reaching farthest from either end and to either side of it, and a capsule
// holds every point between points it contains.
func (c Capsule[T]) ContainsSector(sector Sector[T]) bool {
	outer, inner := c.Float(), sector.Float()
	for _, corner := range inner.corners() {
		if !outer.Contains(corner) {
			return false
		}
	}

	normal := outer.Segment.Normal()
	for _, direction := range [4]Vector[float64]{inner.Center.Subtract(outer.Segment.Start), inner.Center.Subtract(outer.Segment.End), normal, normal.Negate()} {
		if !outer.Contains(inner.Support(direction)) {
			return false
		}
	}

	return true
}

// ContainsSegment reports whether the segment lies within the capsule: both endpoints are
// contained, within the tolerance, and a capsule holds every point between two it contains.
func (c Capsule[T]) ContainsSegment(segment Segment[T]) bool {
	return c.Contains(segment.Start) && c.Contains(segment.End)
}

// ContainsPolygon reports whether the polygon lies within the capsule: every vertex is
// contained, within the tolerance, and a capsule holds every point between points it contains.
// An empty polygon is enclosed by nothing.
func (c Capsule[T]) ContainsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRectangle reports whether the rectangle lies within the capsule: every corner is
// contained, within the tolerance, whatever the rectangle's angle.
func (c Capsule[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies within the capsule: every
// vertex is contained, within the tolerance. An empty polygon is enclosed by nothing.
func (c Capsule[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsBox reports whether the box lies within the capsule: every corner is contained,
// within the tolerance.
func (c Capsule[T]) ContainsBox(box Box[T]) bool {
	for _, corner := range box.corners() {
		if !c.Contains(corner) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the capsule and the circle share a point, as
// Circle.IntersectsCapsule does.
func (c Capsule[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsCapsule(c)
}

// IntersectsCapsule reports whether the capsules overlap: their segments pass within the sum of
// the two radii of each other. Touching capsules intersect, within the tolerance, by the
// comparison Circle.IntersectsCircle makes, on the squared distance.
func (c Capsule[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return lessOrEqualSquared(c.Segment.DistanceSquaredToSegment(capsule.Segment), float64(c.Radius)+float64(capsule.Radius), epsilonAt[T](max(c.magnitude(), capsule.magnitude())))
}

// IntersectsSector reports whether the capsule and the sector share a point, as
// Sector.IntersectsCapsule does.
func (c Capsule[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsCapsule(c)
}

// IntersectsSegment reports whether the capsule and the segment share a point: the segment
// passes within the radius of the capsule's own. Touching shapes intersect, within the
// tolerance, by the same comparison Contains makes, on the squared distance.
func (c Capsule[T]) IntersectsSegment(segment Segment[T]) bool {
	return c.containsSquared(c.Segment.DistanceSquaredToSegment(segment), segment.magnitude())
}

// IntersectsRay reports whether the capsule and the ray share a point, as IntersectsSegment
// decides it on the reach of the ray past the capsule.
func (c Capsule[T]) IntersectsRay(ray Ray[T]) bool {
	return c.IntersectsSegment(ray.reach(c.minMax()))
}

// IntersectsLine reports whether the capsule and the line share a point, as IntersectsSegment
// decides it on the span of the line past the capsule.
func (c Capsule[T]) IntersectsLine(line Line[T]) bool {
	return c.IntersectsSegment(line.span(c.minMax()))
}

// IntersectsPolygon reports whether the capsule and the polygon share a point: the segment of
// the capsule meets the polygon, lying within it included, or an edge of the polygon passes
// within the radius of that segment, as IntersectsSegment decides it. An empty polygon
// intersects nothing.
func (c Capsule[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	if polygon.IsEmpty() || !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.IntersectsSegment(c.Segment) {
		return true
	}

	for edge := range polygon.Edges() {
		if c.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsRectangle reports whether the capsule and the rectangle overlap: the segment of the
// capsule meets the rectangle, lying within it included, or an edge of the rectangle passes
// within the radius of that segment, as IntersectsSegment decides it.
func (c Capsule[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a1, b1 := c.minMax()
	a2, b2 := rectangle.MinMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if rectangle.IntersectsSegment(c.Segment) {
		return true
	}

	for edge := range rectangle.Edges() {
		if c.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsRegularPolygon reports whether the capsule and the regular polygon share a point:
// the segment of the capsule meets the polygon, lying within it included, or an edge of the
// polygon passes within the radius of that segment, as IntersectsSegment decides it. An empty
// polygon intersects nothing.
func (c Capsule[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	a1, b1 := c.minMax()
	a2, b2 := polygon.minMax()

	if polygon.IsEmpty() || !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.IntersectsSegment(c.Segment) {
		return true
	}

	for edge := range polygon.Edges() {
		if c.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectsBox reports whether the capsule and the box overlap: the segment of the capsule
// meets the box, lying within it included, or an edge of the box passes within the radius of
// that segment, as IntersectsSegment decides it.
func (c Capsule[T]) IntersectsBox(box Box[T]) bool {
	a, b := c.minMax()
	if !overlaps(a, b, box.Min, box.Max) {
		return false
	}

	if box.IntersectsSegment(c.Segment) {
		return true
	}

	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if c.IntersectsSegment(edge) {
			return true
		}
	}

	return false
}

// IntersectionSegment returns the points where the segment crosses the capsule boundary, from
// its Start to its End: the first and the last of the points it meets the two round ends and
// the two straight sides at, so two where it passes through, one where it is tangent or has one
// end inside, and none where it misses or lies inside, which IntersectsSegment still reports.
// A point admitted by the tolerance is placed on the boundary: an end of the segment within it
// of the boundary is a crossing, and a segment passing within it of a round end touches it. They
// are found in float64 and rounded once for an integer T. A capsule of no length is a circle
// and one of no radius its segment, and each is crossed as that, at the point
// Segment.IntersectionSegment gives; a segment running along a capsule of no radius shares a
// stretch with it and gives the ends of that stretch.
//
// The points are a Pair, which allocates nothing.
func (c Capsule[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	if c.Segment.IsPoint() {
		return Circle[T]{c.Segment.Start, c.Radius}.IntersectionSegment(segment)
	}

	if point, ok := segment.IntersectionSegment(c.Segment); ok && c.Radius == 0 {
		return Pair[Point[T]]{[2]Point[T]{point}, 1}
	}

	fractions, n := c.fractions(segment)

	return segment.pointsAt(fractions[:n]...)
}

// IntersectionRay returns the points where the ray crosses the capsule boundary, from its Origin
// on, as IntersectionSegment finds them on the reach of the ray past the capsule.
func (c Capsule[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return c.IntersectionSegment(ray.reach(c.minMax()))
}

// IntersectionLine returns the points where the line crosses the capsule boundary, in order
// along its Direction, as IntersectionSegment finds them on the span of the line past the
// capsule.
func (c Capsule[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return c.IntersectionSegment(line.span(c.minMax()))
}

// fractions returns the fractions along the given segment at which it crosses the capsule
// boundary, in increasing order and held to the segment, and how many of the two there are: the
// first and the last of the places it can cross at that lie at the radius from the segment of
// the capsule, which leaves out the half of each end circle within it. The places are its own
// two ends, so an end on the boundary within the tolerance is a crossing, the fractions
// Circle.fractions gives for the circles about the two ends of the capsule, at the tolerance of
// T, a tangent included, and its crossings of the two straight sides. The capsule is convex, so
// a segment crosses its boundary no more than twice, and Segment.HitCapsule reads the first of
// the same fractions.
func (c Capsule[T]) fractions(segment Segment[T]) ([2]float64, int) {
	body, cast := c.Float(), segment.Float()
	offset := body.Segment.Normal().Resize(body.Radius)

	places, n := [8]float64{0, 1}, 2
	for _, center := range [2]Point[T]{c.Segment.Start, c.Segment.End} {
		crossed, count := Circle[T]{center, c.Radius}.fractions(segment)
		n += copy(places[n:], crossed[:count])
	}
	for _, side := range [2]Segment[float64]{body.Segment.Translate(offset), body.Segment.Translate(offset.Negate())} {
		if point, ok := cast.IntersectionSegment(side); ok {
			places[n] = cast.fraction(point)
			n++
		}
	}

	first, last := math.Inf(1), math.Inf(-1)
	for _, t := range places[:n] {
		if !segment.containsAt(t) {
			continue
		}

		t = Clamp(t, 0, 1)
		if point := cast.PointAt(t); c.touchesSquared(body.Segment.distanceSquaredTo(point), point.magnitude()) {
			first, last = min(first, t), max(last, t)
		}
	}

	switch {
	case first > last:
		return [2]float64{}, 0
	case first == last:
		return [2]float64{first}, 1
	default:
		return [2]float64{first, last}, 2
	}
}

// containsSquared reports whether a point at the given squared distance from the segment lies
// within the capsule, boundary included within the tolerance: lessOrEqualSquared on the radius,
// the comparison Contains, DistanceTo, Nearest and ContainsCircle make, at the tolerance of the
// capsule and of geometry whose largest absolute coordinate is the given magnitude.
func (c Capsule[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, float64(c.Radius), epsilonAt[T](max(c.magnitude(), magnitude)))
}

// touchesSquared reports whether a point at the given squared distance from the segment lies on
// the boundary of the capsule within the tolerance: equalSquared on the radius, the test of a
// place a cast crosses the boundary at. The place is computed in float64 and is no lattice
// point, so the tolerance is that of T and no less than that of float64, which an integer T is
// held to as well.
func (c Capsule[T]) touchesSquared(distanceSquared, magnitude float64) bool {
	magnitude = max(c.magnitude(), magnitude)

	return equalSquared(distanceSquared, float64(c.Radius), max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude)))
}

// Equal checks for equal segment and radius with given capsule. The ends are compared in their
// places, so a capsule and the one of the reversed segment, the same points, are not equal.
func (c Capsule[T]) Equal(capsule Capsule[T]) bool {
	return c.Segment.Equal(capsule.Segment) && Equal(c.Radius, capsule.Radius)
}

// IsZero checks if the segment and the radius are zero.
func (c Capsule[T]) IsZero() bool {
	return c.Equal(Capsule[T]{})
}

// Polygon converts the capsule to a Polygon within the given tolerance of it, the outline that
// draws or walks one: the half circle about each end as equal chords, as few as keep it within
// the tolerance of each, joined by the two straight sides, with its vertices on the boundary
// and wound as the corners of a Rectangle are. A negative tolerance is taken absolute, and each
// half circle has at most a thousand chords, which a tolerance of zero gives; a tolerance of
// the radius or more leaves the four corners of the sides. A capsule of no radius gives its two
// ends, and one of no length the circle about its point, never fewer than four chords around
// it. For integer T the vertices are rounded
// and one rounding onto the vertex before it is left out.
func (c Capsule[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{c.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of the polygon Polygon returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for them.
func (c Capsule[T]) AppendPolygon(dst []Point[T], tolerance float64) []Point[T] {
	angle := c.Segment.Angle()
	ends := [2]Arc[T]{{c.Segment.End, c.Radius, angle - math.Pi/2, math.Pi}, {c.Segment.Start, c.Radius, angle + math.Pi/2, math.Pi}}
	pieces := ends[0].pieces(tolerance)
	if c.Segment.IsPoint() {
		pieces = max(pieces, 2)
	}

	from := len(dst)
	dst = slices.Grow(dst, 2*(int(pieces)+1))

	for _, end := range ends {
		for i := 0.0; i <= pieces; i++ {
			if point := end.PointAt(i / pieces); len(dst) == from || point != dst[len(dst)-1] {
				dst = append(dst, point)
			}
		}
	}

	if last := len(dst) - 1; last > from && dst[last].Equal(dst[from]) {
		dst = dst[:last]
	}

	return dst
}

// Cast converts the capsule to a Capsule of another number type, rounding as Cast does.
func (c Capsule[T]) Cast[R Number]() Capsule[R] {
	return Capsule[R]{c.Segment.Cast[R](), Cast[R](float64(c.Radius))}
}

// Int converts the capsule to a Capsule[int].
func (c Capsule[T]) Int() Capsule[int] {
	return Capsule[int]{c.Segment.Int(), Int(c.Radius)}
}

// Float converts the capsule to a Capsule[float64].
func (c Capsule[T]) Float() Capsule[float64] {
	return Capsule[float64]{c.Segment.Float(), float64(c.Radius)}
}

// String returns the capsule in the form of its constructor: Cap((x,y);(x,y);r).
func (c Capsule[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Cap("...)
	buf = c.Segment.Start.appendString(buf)
	buf = append(buf, ';')
	buf = c.Segment.End.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, c.Radius)
	buf = append(buf, ')')

	return string(buf)
}

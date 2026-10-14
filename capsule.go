package geom

import (
	"math"
	"slices"
)

// Capsule is every point within a radius of a segment: a circle for a segment of zero length, the
// segment itself for a zero radius.
//
// The radius is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative radius written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it. A capsule has no Transform, since a matrix that scales one axis more
// than the other does not take it to a capsule.
type Capsule[T Number] struct {
	Segment Segment[T] `json:",embed"`
	Radius  T          `json:"r"`
}

// Cap is shorthand for Capsule{Segment{start, end}, radius}, with the radius taken absolute.
func Cap[T Number](start, end Point[T], radius T) Capsule[T] {
	return Capsule[T]{Segment[T]{start, end}, Abs(radius)}
}

// Centroid returns the midpoint of the segment.
func (c Capsule[T]) Centroid() Point[T] {
	return c.Segment.Midpoint()
}

// Area returns the area (π·r² + 2r·length).
func (c Capsule[T]) Area() float64 {
	radius := float64(c.Radius)

	return radius * (float64(Pi*radius) + 2*c.Segment.Length())
}

// Perimeter returns the length of the boundary (2π·r + 2·length).
func (c Capsule[T]) Perimeter() float64 {
	radius := float64(c.Radius)

	return 2 * (float64(Pi*radius) + c.Segment.Length())
}

// Inertia returns the polar second moment of area about the Centroid, the rotational inertia of the
// capsule at unit density, and zero for a zero radius.
func (c Capsule[T]) Inertia() float64 {
	radius, length := float64(c.Radius), c.Segment.Length()
	rectangle := float64(length*length*length/6) + float64(2*radius*radius*length)
	discs := float64(Pi*radius*radius*radius/2) + float64(Pi*radius*length*length/4)

	return radius * (rectangle + discs)
}

// Bounds returns the axis-aligned bounding box.
func (c Capsule[T]) Bounds() Box[T] {
	a, b := c.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the smallest circle enclosing the capsule. For integer T the center is
// rounded and the radius rounded up, so the circle still encloses the capsule.
func (c Capsule[T]) BoundingCircle() Circle[T] {
	circle := c.Segment.BoundingCircle()

	return Circle[T]{circle.Center, circle.Radius + c.Radius}
}

// minMax returns the corners of Bounds, the prefilter of the intersection tests.
func (c Capsule[T]) minMax() (Point[T], Point[T]) {
	a, b := c.Segment.minMax()

	return a.AddXY(-c.Radius, -c.Radius), b.AddXY(c.Radius, c.Radius)
}

// magnitude returns the largest absolute coordinate of Bounds, the magnitude the tolerance of a
// comparison with the capsule is taken from.
func (c Capsule[T]) magnitude() float64 {
	return c.Segment.magnitude() + float64(c.Radius)
}

// Translate returns the capsule moved by the vector.
func (c Capsule[T]) Translate(vector Vector[T]) Capsule[T] {
	return Capsule[T]{c.Segment.Translate(vector), c.Radius}
}

// MoveTo returns the capsule moved so that the midpoint of its segment lies at the point.
func (c Capsule[T]) MoveTo(point Point[T]) Capsule[T] {
	return Capsule[T]{c.Segment.MoveTo(point), c.Radius}
}

// Rotate returns the capsule turned by the angle in radians about the midpoint of its segment. For
// integer T the ends are rounded, and only multiples of a quarter turn keep the length exactly.
func (c Capsule[T]) Rotate(angle float64) Capsule[T] {
	return Capsule[T]{c.Segment.Rotate(angle), c.Radius}
}

// Scale returns the capsule scaled by the factor about the midpoint of its segment, length and
// radius alike. A negative factor swaps the ends, which gives the same points.
func (c Capsule[T]) Scale(factor float64) Capsule[T] {
	return Capsule[T]{c.Segment.Scale(factor), Abs(Multiply(c.Radius, factor))}
}

// Unscale returns the capsule divided by the factor about the midpoint of its segment, as Scale
// treats a negative one. It panics for a zero factor.
func (c Capsule[T]) Unscale(factor float64) Capsule[T] {
	return Capsule[T]{c.Segment.Unscale(factor), Abs(Divide(c.Radius, factor))}
}

// Grow returns the capsule with its radius increased by the amount, clamped at zero.
func (c Capsule[T]) Grow(amount T) Capsule[T] {
	return Capsule[T]{c.Segment, max(c.Radius+amount, 0)}
}

// Shrink returns the capsule with its radius decreased by the amount, clamped at zero, where the
// capsule is its segment.
func (c Capsule[T]) Shrink(amount T) Capsule[T] {
	return Capsule[T]{c.Segment, max(c.Radius-amount, 0)}
}

// Lerp returns the capsule interpolated toward the given one, both ends and the radius together,
// and extrapolated outside [0, 1]. The radius is taken absolute, so extrapolating past zero grows
// the capsule again.
func (c Capsule[T]) Lerp(capsule Capsule[T], t float64) Capsule[T] {
	return Capsule[T]{c.Segment.Lerp(capsule.Segment, t), Abs(Lerp(c.Radius, capsule.Radius, t))}
}

// Canonical returns the capsule with its radius taken absolute, the receiver when it already is.
func (c Capsule[T]) Canonical() Capsule[T] {
	return Capsule[T]{c.Segment, Abs(c.Radius)}
}

// DistanceTo returns the distance from the point to the capsule, zero inside.
func (c Capsule[T]) DistanceTo(point Point[T]) float64 {
	distanceSquared := c.Segment.DistanceSquaredTo(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return 0
	}

	return math.Sqrt(distanceSquared) - float64(c.Radius)
}

// DistanceSquaredTo returns the square of DistanceTo.
func (c Capsule[T]) DistanceSquaredTo(point Point[T]) float64 {
	distance := c.DistanceTo(point)

	return distance * distance
}

// DistanceToCapsule returns the distance between the nearest points of the capsules, zero where
// they share a point.
func (c Capsule[T]) DistanceToCapsule(capsule Capsule[T]) float64 {
	if c.IntersectsCapsule(capsule) {
		return 0
	}

	return c.Segment.DistanceToSegment(capsule.Segment) - float64(c.Radius) - float64(capsule.Radius)
}

// DistanceSquaredToCapsule returns the square of DistanceToCapsule.
func (c Capsule[T]) DistanceSquaredToCapsule(capsule Capsule[T]) float64 {
	distance := c.DistanceToCapsule(capsule)

	return distance * distance
}

// Nearest returns the point of the capsule nearest to the given point, the point itself inside.
func (c Capsule[T]) Nearest(point Point[T]) Point[T] {
	distanceSquared := c.Segment.DistanceSquaredTo(point)
	if c.containsSquared(distanceSquared, point.magnitude()) {
		return point
	}

	foot := c.Segment.Float().foot(point.Float())
	offset := point.Float().Subtract(foot).Multiply(float64(c.Radius) / math.Sqrt(distanceSquared))

	return foot.Add(offset).Cast[T]()
}

// Support returns the point of the capsule farthest along the direction, whatever its length, and
// along X from Start for a zero direction.
func (c Capsule[T]) Support(direction Vector[float64]) Point[float64] {
	return c.Segment.Support(direction).Add(direction.Normalize().Multiply(float64(c.Radius)))
}

// Contains reports whether the point lies in the capsule.
func (c Capsule[T]) Contains(point Point[T]) bool {
	return c.containsSquared(c.Segment.DistanceSquaredTo(point), point.magnitude())
}

// ContainsCircle reports whether the circle lies entirely in the capsule.
func (c Capsule[T]) ContainsCircle(circle Circle[T]) bool {
	far := c.Segment.DistanceTo(circle.Center) + float64(circle.Radius)

	return c.containsSquared(far*far, circle.magnitude())
}

// ContainsCapsule reports whether the given capsule lies entirely in this one.
func (c Capsule[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return c.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && c.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the sector lies entirely in the capsule.
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

// ContainsSegment reports whether the segment lies entirely in the capsule.
func (c Capsule[T]) ContainsSegment(segment Segment[T]) bool {
	return c.Contains(segment.Start) && c.Contains(segment.End)
}

// ContainsPolygon reports whether the polygon lies entirely in the capsule.
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

// ContainsRectangle reports whether the rectangle lies entirely in the capsule.
func (c Capsule[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !c.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the capsule.
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

// ContainsBox reports whether the box lies entirely in the capsule.
func (c Capsule[T]) ContainsBox(box Box[T]) bool {
	for _, corner := range box.corners() {
		if !c.Contains(corner) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the capsule and the circle share a point.
func (c Capsule[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsCapsule(c)
}

// IntersectsCapsule reports whether the capsules share a point.
func (c Capsule[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return lessOrEqualSquared(c.Segment.DistanceSquaredToSegment(capsule.Segment), float64(c.Radius)+float64(capsule.Radius), epsilonAt[T](max(c.magnitude(), capsule.magnitude())))
}

// IntersectsSector reports whether the capsule and the sector share a point.
func (c Capsule[T]) IntersectsSector(sector Sector[T]) bool {
	return lessOrEqualSquared(sector.distanceSquaredToSegment(c.Segment.Float()), float64(c.Radius), sector.epsilon(c.magnitude()))
}

// IntersectsSegment reports whether the capsule and the segment share a point.
func (c Capsule[T]) IntersectsSegment(segment Segment[T]) bool {
	return c.containsSquared(c.Segment.DistanceSquaredToSegment(segment), segment.magnitude())
}

// IntersectsRay reports whether the capsule and the ray share a point.
func (c Capsule[T]) IntersectsRay(ray Ray[T]) bool {
	return c.IntersectsSegment(ray.reach(c.minMax()))
}

// IntersectsLine reports whether the capsule and the line share a point.
func (c Capsule[T]) IntersectsLine(line Line[T]) bool {
	return c.IntersectsSegment(line.span(c.minMax()))
}

// IntersectsPolygon reports whether the capsule and the polygon share a point.
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

// IntersectsRectangle reports whether the capsule and the rectangle share a point.
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

// IntersectsRegularPolygon reports whether the capsule and the regular polygon share a point.
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

// IntersectsBox reports whether the capsule and the box share a point.
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

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End:
// two where it passes through, one where it is tangent or has one end inside, and none where it
// misses. A segment entirely inside crosses no boundary, so it gives none while IntersectsSegment
// holds. A capsule of zero length is crossed as its circle, and one of zero radius as its segment;
// a segment running along that one gives the ends of the shared stretch.
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

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (c Capsule[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return c.IntersectionSegment(ray.reach(c.minMax()))
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (c Capsule[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return c.IntersectionSegment(line.span(c.minMax()))
}

// fractions returns the fractions along the segment, ascending and clamped to it, at which it
// crosses the boundary, and how many there are: the first and last of its own ends, its crossings
// of the end circles and of the two sides that lie on the boundary. Segment.HitCapsule reads the
// first of the same values.
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

// containsSquared reports whether a point at the squared distance from the segment lies in the
// capsule, within the tolerance.
func (c Capsule[T]) containsSquared(distanceSquared, magnitude float64) bool {
	return lessOrEqualSquared(distanceSquared, float64(c.Radius), epsilonAt[T](max(c.magnitude(), magnitude)))
}

// touchesSquared reports whether a point at the squared distance from the segment lies on the
// boundary, within the tolerance. The point is computed in float64 and is no lattice point, so the
// tolerance is never less than that of float64.
func (c Capsule[T]) touchesSquared(distanceSquared, magnitude float64) bool {
	magnitude = max(c.magnitude(), magnitude)

	return equalSquared(distanceSquared, float64(c.Radius), max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude)))
}

// Equal reports whether the capsules have equal segment and radius. The ends are compared in place,
// so a capsule and the one of the reversed segment are not equal.
func (c Capsule[T]) Equal(capsule Capsule[T]) bool {
	return c.Segment.Equal(capsule.Segment) && Equal(c.Radius, capsule.Radius)
}

// IsZero reports whether the capsule is the zero value.
func (c Capsule[T]) IsZero() bool {
	return c.Equal(Capsule[T]{})
}

// Polygon returns the capsule as a polygon within the tolerance of it: each half circle as few
// equal chords as keep within it, at most 1024, joined by the two sides and wound like
// Rectangle.Vertices. A negative tolerance is taken absolute, and a tolerance of the radius or more
// leaves the four corners of the sides. A capsule of zero radius gives its two ends, and one of
// zero length its circle in at least four chords.
func (c Capsule[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{c.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of Polygon(tolerance) to dst and returns the extended slice.
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

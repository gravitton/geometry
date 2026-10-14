package geom

import (
	"cmp"
	"iter"
	"slices"
)

// edgesOf returns an iterator over the edges joining the vertices in order, the last closing back
// to the first. The slice is never retained, so a caller may pass a slice of a local array.
func edgesOf[T Number](vertices []Point[T]) iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		n := len(vertices)
		for i, vertex := range vertices {
			if !yield(Segment[T]{vertex, vertices[(i+1)%n]}) {
				return
			}
		}
	}
}

// edgeWalk folds the edges of a closed outline into the walk of a shape: the squared distance to
// the nearest edge, starting infinite, that edge, unless the point lies on one, which stops the
// walk at zero, and the even-odd parity of a ray from the point. A shape starts it as a literal and
// steps its own Edges, since an iterator passed in would allocate.
type edgeWalk[T Number] struct {
	inside   bool
	distance float64
	edge     Segment[T]
}

// step folds one edge in and reports whether the point lies on it, which ends the walk at distance
// zero whatever the remaining edges say.
func (w *edgeWalk[T]) step(edge Segment[T], point Point[T]) bool {
	distance := edge.distanceSquaredTo(point)
	if lessOrEqualSquared(distance, 0, epsilonAt[T](max(edge.magnitude(), point.magnitude()))) {
		w.distance = 0

		return true
	}

	if distance < w.distance {
		w.edge = edge
	}

	w.distance = min(w.distance, distance)

	if edge.crossesRay(point) {
		w.inside = !w.inside
	}

	return false
}

// result returns the squared distance from the point to the outline: zero inside, infinite where
// there was no edge.
func (w edgeWalk[T]) result() float64 {
	if w.inside {
		return 0
	}

	return w.distance
}

// nearest returns the point itself where result is zero, and otherwise its foot on the nearest
// edge, the zero point where there was no edge.
func (w edgeWalk[T]) nearest(point Point[T]) Point[T] {
	if w.result() == 0 {
		return point
	}

	return w.edge.foot(point)
}

// clears reports whether the circle about the point lies within the outline: the point inside and
// the nearest edge no nearer than the radius.
func (w edgeWalk[T]) clears(circle Circle[T]) bool {
	return w.result() == 0 && greaterOrEqualSquared(w.distance, float64(circle.Radius), circle.epsilonWith(w.edge.magnitude()))
}

// reaches reports whether the circle about the point shares a point with the outline, at the
// tolerance Circle.IntersectsSegment applies to the nearest edge.
func (w edgeWalk[T]) reaches(circle Circle[T]) bool {
	return circle.containsSquared(w.result(), w.edge.magnitude())
}

// edgeIntersections collects the points where a segment crosses the edges of an outline, giving a
// vertex two edges cross once. Points before from are neither compared nor sorted; a nil slice is
// allocated on the first crossing with room for the two of a convex outline.
type edgeIntersections[T Number] struct {
	segment Segment[T]
	points  []Point[T]
	from    int
}

// add folds one edge in, keeping its crossing with the segment unless an earlier edge gave an Equal
// point.
func (e *edgeIntersections[T]) add(edge Segment[T]) {
	point, ok := e.segment.IntersectionSegment(edge)
	if !ok || slices.ContainsFunc(e.points[e.from:], point.Equal) {
		return
	}

	if e.points == nil {
		e.points = make([]Point[T], 0, 2)
	}

	e.points = append(e.points, point)
}

// sorted returns the slice with the crossings ordered from the segment's Start after the points it
// was started with.
func (e *edgeIntersections[T]) sorted() []Point[T] {
	slices.SortFunc(e.points[e.from:], func(a, b Point[T]) int {
		return e.segment.compareDistance(a, b)
	})

	return e.points
}

// edgeSpan keeps only the first and the last crossing of a segment with the edges of an outline,
// and the edge of the first, so clipping, Intersection and Hit of a convex shape agree and allocate
// nothing. The first crossing is the first of any outline, so a concave Hit reads it too.
type edgeSpan[T Number] struct {
	segment Segment[T]
	ends    [2]Point[T]
	edge    Segment[T]
	crossed bool
}

// add folds one edge in, keeping its crossing where it comes before the first or after the last. A
// vertex crossed by two edges goes to the edge that faces Start, the earlier one where both or
// neither do.
func (e *edgeSpan[T]) add(edge Segment[T]) {
	point, ok := e.segment.IntersectionSegment(edge)

	switch {
	case !ok:
	case !e.crossed:
		e.ends, e.edge, e.crossed = [2]Point[T]{point, point}, edge, true
	case e.segment.compareDistance(point, e.ends[0]) < 0:
		e.ends[0], e.edge = point, edge
	case point == e.ends[0] && edge.faces(e.segment.Start, e.edge) && !e.edge.faces(e.segment.Start, edge):
		e.edge = edge
	case e.segment.compareDistance(point, e.ends[1]) > 0:
		e.ends[1] = point
	}
}

// crossings returns the first and the last crossing, the same point twice where there was one, and
// nil where there was none.
func (e *edgeSpan[T]) crossings() []Point[T] {
	if !e.crossed {
		return nil
	}

	return e.ends[:]
}

// pair returns the first and the last crossing, one point where they are Equal and none where there
// was none.
func (e *edgeSpan[T]) pair() Pair[Point[T]] {
	switch {
	case !e.crossed:
		return Pair[Point[T]]{}
	case e.ends[0].Equal(e.ends[1]):
		return Pair[Point[T]]{[2]Point[T]{e.ends[0]}, 1}
	default:
		return Pair[Point[T]]{e.ends, 2}
	}
}

// hit returns the Hit of the segment at the first crossing, with the normal of its edge, and false
// where there was none.
func (e *edgeSpan[T]) hit() (Hit[T], bool) {
	if !e.crossed {
		return Hit[T]{}, false
	}

	return e.segment.hitAt(e.ends[0], e.edge.Float().Normal()), true
}

// edgePass keeps, in order from Start, the next crossings of a segment with the edges of an outline
// beyond a reached point, as many as its array holds, Point.Compare breaking ties. A shape passes
// again from the last one only where the array is spent, so it visits every crossing without a
// slice.
type edgePass[T Number] struct {
	segment Segment[T]
	reached Point[T]
	points  [16]Point[T]
	n, read int
	touched bool
}

// add folds one edge in, keeping its crossing where it lies beyond the reached point and noting a
// touch where it is Equal to it.
func (e *edgePass[T]) add(edge Segment[T]) {
	point, ok := e.segment.IntersectionSegment(edge)

	switch {
	case !ok:
	case point.Equal(e.reached):
		e.touched = true
	case e.precedes(e.reached, point):
		e.insert(point)
	}
}

// insert places the point in order, dropping the farthest where the array is full.
func (e *edgePass[T]) insert(point Point[T]) {
	at := e.n
	for at > 0 && e.precedes(point, e.points[at-1]) {
		at--
	}

	if at == len(e.points) {
		return
	}

	e.n = min(e.n+1, len(e.points))
	copy(e.points[at+1:e.n], e.points[at:])
	e.points[at] = point
}

// next returns the first unread point beyond the given one, skipping those Equal to it, and false
// where none is left.
func (e *edgePass[T]) next(point Point[T]) (Point[T], bool) {
	for e.read < e.n && e.points[e.read].Equal(point) {
		e.read++
	}

	if e.read == e.n {
		return Point[T]{}, false
	}

	e.read++

	return e.points[e.read-1], true
}

// spent reports whether the array came back full and has all been read, so more crossings may lie
// beyond the last.
func (e *edgePass[T]) spent() bool {
	return e.read == len(e.points)
}

// touchedAt reports whether the point is the reached one and a crossing was Equal to it.
func (e *edgePass[T]) touchedAt(point Point[T]) bool {
	return e.touched && point == e.reached
}

// precedes reports whether a comes before b from the segment's Start: nearer, or as near and first
// by Point.Compare.
func (e *edgePass[T]) precedes(a, b Point[T]) bool {
	return cmp.Or(e.segment.compareDistance(a, b), a.Compare(b)) < 0
}

// edgeProbe tests the edges of an outline, of extent a, b, against one aimed segment at a time,
// with an extent prefilter on each. A shape ranges the Edges of both sides itself, since an
// iterator passed in would allocate.
type edgeProbe[T Number] struct {
	a, b    Point[T]
	segment Segment[T]
	c, d    Point[T]
}

// aim points the probe at the segment and reports whether its extent overlaps that of the outline.
func (p *edgeProbe[T]) aim(segment Segment[T]) bool {
	p.segment = segment
	p.c, p.d = segment.minMax()

	return overlaps(p.a, p.b, p.c, p.d)
}

// meets reports whether the edge shares a point with the aimed segment.
func (p *edgeProbe[T]) meets(edge Segment[T]) bool {
	e, f := edge.minMax()

	return overlaps(p.c, p.d, e, f) && p.segment.IntersectsSegment(edge)
}

// edgeConvexity folds the edge directions of an outline into the test IsConvex makes: the sense of
// the turns so far and how often the X component has changed sign, which an outline going around
// once does exactly twice.
type edgeConvexity struct {
	first, previous Vector[float64]
	turned          bool
	clockwise       bool
	x               float64
	reversals       int
}

// step folds one direction in and reports whether the outline can still be convex. A zero
// direction, from a repeated vertex, is skipped; a NaN turn has no sense, so it never settles one.
func (c *edgeConvexity) step(direction Vector[float64]) bool {
	if !direction.hasDirection() {
		return true
	}

	if !c.first.hasDirection() {
		c.first = direction
	}

	if c.previous.hasDirection() {
		switch turn := c.previous.Cross(direction); {
		case turn > 0 || turn < 0:
			if c.turned && c.clockwise != (turn > 0) {
				return false
			}

			c.turned, c.clockwise = true, turn > 0
		case c.previous.Dot(direction) < 0:
			return false
		}
	}

	if direction.X != 0 {
		if c.x != 0 && (c.x < 0) != (direction.X < 0) {
			c.reversals++
		}

		c.x = direction.X
	}

	c.previous = direction

	return true
}

// result closes the outline back to the first direction and reports whether it is convex.
func (c *edgeConvexity) result() bool {
	return c.step(c.first) && c.turned && c.reversals <= 2
}

// edgeClip cuts an outline to the inner side of one line after another, Sutherland–Hodgman style.
// The vertices are read and never written; each cut appends its outline and moves it down to from,
// so dst needs room for two outlines. A region wound against Rectangle.Vertices is reversed.
type edgeClip[T Number] struct {
	dst      []Point[T]
	from     int
	points   []Point[T]
	reversed bool
}

// add cuts the outline to the side of the line through the edge that Segment.Normal points to on an
// edge of a Rectangle.
func (c *edgeClip[T]) add(edge Segment[T]) {
	if c.reversed {
		edge = edge.Reverse()
	}

	c.cut(edge.Line(), 1)
}

// cut keeps the part of the outline on the side the sense signs, the line included. A vertex on the
// line is kept and no crossing is placed beside it, and both senses read the same sides, so the two
// cuts of one line agree. A line with no direction cuts nothing, and a clip of no real cut leaves
// dst without the outline.
func (c *edgeClip[T]) cut(line Line[T], sense float64) {
	if !line.Direction.hasDirection() || len(c.points) == 0 {
		return
	}

	start := len(c.dst)
	previous := c.points[len(c.points)-1]
	previousSide := sense * c.side(line, previous)
	for _, point := range c.points {
		side := sense * c.side(line, point)
		if (side < 0 && previousSide > 0) || (side > 0 && previousSide < 0) {
			c.keep(start, c.crossing(previous, point, previousSide, side))
		}

		if side >= 0 {
			c.keep(start, point)
		}

		previous, previousSide = point, side
	}

	if n := len(c.dst); n-start > 1 && c.dst[n-1].Equal(c.dst[start]) {
		c.dst = c.dst[:n-1]
	}

	c.dst = append(c.dst[:c.from], c.dst[start:]...)
	c.points = c.dst[c.from:]
}

// keep appends the point to the outline begun at start unless it is Equal to the last, making room
// on the first for one vertex more than the outline had.
func (c *edgeClip[T]) keep(start int, point Point[T]) {
	n := len(c.dst)
	if n == start {
		c.dst = slices.Grow(c.dst, len(c.points)+1)
	}

	if n > start && c.dst[n-1].Equal(point) {
		return
	}

	c.dst = append(c.dst, point)
}

// side returns the cross product Line.Side reads, and zero for a point the line contains.
func (*edgeClip[T]) side(line Line[T], point Point[T]) float64 {
	if line.Contains(point) {
		return 0
	}

	return line.side(point)
}

// crossing returns where the edge from a to b crosses the line, given their sides of opposite
// signs, divided once so a crossing on a half unit is exact before rounding.
func (*edgeClip[T]) crossing(a, b Point[T], sideA, sideB float64) Point[T] {
	start := a.Float()
	along := b.Float().Subtract(start)
	denominator := sideA - sideB

	return Point[T]{Cast[T](start.X + float64(sideA*along.X)/denominator), Cast[T](start.Y + float64(sideA*along.Y)/denominator)}
}

// result returns the slice with the outline the cuts left appended to it.
func (c *edgeClip[T]) result() []Point[T] {
	return c.dst
}

// unionOutline gathers the outline of the union of two convex polygons, appended to dst after from:
// the vertices of each outside the other and the crossings of their edges. It also averages the
// points both share, a point inside both about which the outline is sorted by angle.
type unionOutline[T Number] struct {
	dst            []Point[T]
	from           int
	sumX, sumY     float64
	count          int
	first, other   Point[T]
	paired, spread bool
}

// add takes a vertex of one polygon: onto the outline outside the other, among the shared points
// inside it.
func (u *unionOutline[T]) add(point Point[T], inside bool) {
	if inside {
		u.share(point)
	} else {
		u.dst = append(u.dst, point)
	}
}

// cross takes a crossing of the edges, on the outline and shared by both.
func (u *unionOutline[T]) cross(point Point[T]) {
	u.dst = append(u.dst, point)
	u.share(point)
}

// share adds a point both polygons hold to the average, and notes on exact signs whether the shared
// points so far enclose an area.
func (u *unionOutline[T]) share(point Point[T]) {
	u.sumX += float64(point.X)
	u.sumY += float64(point.Y)
	u.count++

	switch {
	case u.count == 1:
		u.first = point
	case !u.paired && !point.Equal(u.first):
		u.other, u.paired = point, true
	case u.paired && !u.spread:
		u.spread = u.other.Subtract(u.first).Cross(point.Subtract(u.first)) != 0
	}
}

// result sorts the outline by angle about the shared average in the winding of the first polygon,
// drops a point Equal to the one before, and starts it at the first point appended. Shared points
// that enclose no area leave dst as it was given.
func (u *unionOutline[T]) result(reversed bool) []Point[T] {
	if !u.spread || len(u.dst) == u.from {
		return u.dst[:u.from]
	}

	start := u.dst[u.from]
	center := Point[float64]{u.sumX / float64(u.count), u.sumY / float64(u.count)}
	outline := u.dst[u.from:]
	slices.SortFunc(outline, func(a, b Point[T]) int {
		return compareAngles(a.Float().Subtract(center), b.Float().Subtract(center))
	})
	if reversed {
		slices.Reverse(outline)
	}

	outline = slices.CompactFunc(outline, Point[T].Equal)
	if i := slices.IndexFunc(outline, start.Equal); i > 0 {
		slices.Reverse(outline[:i])
		slices.Reverse(outline[i:])
		slices.Reverse(outline)
	}

	return u.dst[:u.from+len(outline)]
}

// compareAngles orders two vectors by Vector.Angle without trigonometry: by the half plane each
// lies in, then by the sign of their cross product.
func compareAngles(a, b Vector[float64]) int {
	if halfA, halfB := angleHalf(a), angleHalf(b); halfA != halfB {
		return halfA - halfB
	}

	return -int(Sign(a.Cross(b)))
}

// angleHalf returns 0 below the X axis, 1 above or along it, and 2 against it or for the zero
// vector.
func angleHalf(v Vector[float64]) int {
	switch {
	case v.Y < 0:
		return 0
	case v.Y > 0 || v.X > 0:
		return 1
	default:
		return 2
	}
}

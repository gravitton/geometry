package geom

import (
	"cmp"
	"iter"
	"slices"
)

// edgesOf iterates the edges joining the vertices in order, the last one closing back to the
// first: the one iterator Polygon.Edges and Rectangle.Edges walk, so the edges a segment
// crosses are the edges the walk reads. The slice is read as the edges are
// yielded and never retained, so a caller may pass a slice of a local array.
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

// edgeWalk folds the edges of a closed outline one at a time into the even-odd walk every
// shape's walk makes: the squared distance to the nearest edge so far, starting infinitely
// far, the edge it was measured to, and whether a ray from the point has crossed an odd number
// of edges. A shape starts it as a literal, ranges its own Edges and calls step on each, so the
// walk shares its rule without an iterator crossing a function boundary, which would allocate
// it. The distance and the nearest point are both read off the one walk, so Nearest returns
// the point itself exactly where DistanceSquaredTo is zero. The readers take the walk by
// value, so a shape reads them straight off the walk it returns.
type edgeWalk[T Number] struct {
	inside   bool
	distance float64
	edge     Segment[T]
}

// step folds one edge in and reports whether the point lies on it within the tolerance, as
// Segment.DistanceSquaredTo snaps it, where the walk is over: the point is on the boundary at
// distance zero whatever the remaining edges say.
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

// result returns the squared distance from the point to the outline after every edge: zero
// inside by the even-odd rule, the distance to the nearest edge otherwise, and infinity where
// there was no edge.
func (w edgeWalk[T]) result() float64 {
	if w.inside {
		return 0
	}

	return w.distance
}

// nearest returns the point of the outline nearest to the point after every edge: the point
// itself where result is zero, and the foot on the nearest edge otherwise, which is the zero
// point where there was no edge.
func (w edgeWalk[T]) nearest(point Point[T]) Point[T] {
	if w.result() == 0 {
		return point
	}

	return w.edge.foot(point)
}

// clears reports whether a circle of the given radius about the point lies within the outline
// after every edge: the point inside or on it, where result is zero, and the nearest edge no
// nearer than the radius within the tolerance, so a circle touching an edge from inside is
// enclosed. A point on the boundary is at distance zero and clears only a radius within the
// tolerance.
func (w edgeWalk[T]) clears(circle Circle[T]) bool {
	return w.result() == 0 && greaterOrEqualSquared(w.distance, float64(circle.Radius), circle.epsilonWith(w.edge.magnitude()))
}

// reaches reports whether a circle about the point shares a point with the outline after every
// edge: the point inside or on it, where result is zero, or the nearest edge within the radius,
// by the comparison Circle.Contains makes, at the tolerance of the circle and that edge, the
// one Circle.IntersectsSegment applies to it.
func (w edgeWalk[T]) reaches(circle Circle[T]) bool {
	return circle.containsSquared(w.result(), w.edge.magnitude())
}

// edgeIntersections collects the points where a segment crosses the edges of a polygon, fed
// one at a time by the polygon ranging its own Edges, as edgeWalk folds a walk: each crossing by
// Segment.IntersectionSegment, a vertex hit by two edges counted once. It is started with the
// slice the crossings are appended to and the length it had, so the points already there are
// neither compared nor sorted; a nil slice is allocated on the first crossing with room for the
// two a convex outline can have. Only a concave outline grows it, and the convex shapes keep
// their two outermost crossings in edgeSpan instead.
type edgeIntersections[T Number] struct {
	segment Segment[T]
	points  []Point[T]
	from    int
}

// add folds one edge in, keeping its crossing with the segment unless a previous edge gave
// a point that compares Equal to it.
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

// sorted returns the slice with the crossings ordered from the segment's Start after the points
// it was started with, and the slice as it was where there was none.
func (e *edgeIntersections[T]) sorted() []Point[T] {
	slices.SortFunc(e.points[e.from:], func(a, b Point[T]) int {
		return e.segment.compareDistance(a, b)
	})

	return e.points
}

// edgeSpan folds the points where a segment crosses the edges of a convex outline into the
// first and the last from the segment's Start, as edgeIntersections collects them all: the ends
// of the part inside, which is all Segment.clipConvex reads, and the boundary points the
// Intersection of a convex shape returns, so the two agree and the crossings are never more than
// two. It holds them in an array of its own rather than a slice, so a shape that only clips
// allocates nothing. The edge the first crossing lies on is kept beside it, the boundary a Hit
// reads its normal from, and the first crossing is the first of any outline, so the Hit of a
// concave polygon reads the span too.
type edgeSpan[T Number] struct {
	segment Segment[T]
	ends    [2]Point[T]
	edge    Segment[T]
	crossed bool
}

// add folds one edge in, keeping its crossing with the segment where it comes before the first
// or after the last so far. A vertex two edges cross at goes to the one that faces Start, and
// stays with the edge added first where both do or neither does.
func (e *edgeSpan[T]) add(edge Segment[T]) {
	point, ok := e.segment.IntersectionSegment(edge)

	switch {
	case !ok:
	case !e.crossed:
		e.ends, e.edge, e.crossed = [2]Point[T]{point, point}, edge, true
	case e.segment.compareDistance(point, e.ends[0]) < 0:
		e.ends[0], e.edge = point, edge
	case point == e.ends[0] && e.faces(edge, e.edge) && !e.faces(e.edge, edge):
		e.edge = edge
	case e.segment.compareDistance(point, e.ends[1]) > 0:
		e.ends[1] = point
	}
}

// faces reports whether the edge is turned toward the segment's Start at the vertex it shares
// with the neighbour: Start lies strictly on the side of its line the neighbour does not run
// to, on the exact signs Line.side gives, so no winding is asked for. The neighbour is read by
// both its ends, the shared one adding nothing. At a corner pointing out, that is the outer
// side of the edge; at a notch neither edge faces, both being seen from within it.
func (e *edgeSpan[T]) faces(edge, neighbour Segment[T]) bool {
	line := edge.Line()

	return float64(line.side(e.segment.Start)*(line.side(neighbour.Start)+line.side(neighbour.End))) < 0
}

// crossings returns the first and the last crossing from Start, the same point twice where
// there was one, and nil where there was none.
func (e *edgeSpan[T]) crossings() []Point[T] {
	if !e.crossed {
		return nil
	}

	return e.ends[:]
}

// pair returns the first and the last crossing from Start, the one point where they compare
// Equal and none where there was none: the boundary points of a convex outline, whose ends
// clipConvex reads.
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

// hit returns the first crossing from Start as the Hit of the segment on the edge it lies on,
// and false where there was none.
func (e *edgeSpan[T]) hit() (Hit[T], bool) {
	if !e.crossed {
		return Hit[T]{}, false
	}

	return e.segment.hitAt(e.ends[0], e.edge.Float().Normal()), true
}

// edgePass gathers, among the points where a segment crosses the edges of an outline, the next
// ones after a point already reached on the way from Start to End: the nearest beyond it, in the
// order edgeIntersections sorts them with Point.Compare deciding a tie, as many as its array
// holds. A point that compares Equal to the reached one is that point again, and the pass
// notes that the reached point touches the boundary. A shape passes over its Edges once and reads
// the points it kept with next, and passes again from the last of them only where the array
// came back full and is spent, so the crossings are visited in order without a slice to hold
// them.
type edgePass[T Number] struct {
	segment Segment[T]
	reached Point[T]
	points  [16]Point[T]
	n, read int
	touched bool
}

// add folds one edge in, keeping its crossing with the segment where it lies beyond the reached
// point.
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

// insert places the point among those kept, in their order from the reached point, dropping
// the farthest where the array is full, the point itself where it is the farthest.
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

// next returns the first point kept that has not been read and lies beyond the given one, and
// false where none is left: a point that compares Equal to the given one is that point again
// and is passed over.
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

// spent reports whether the array came back full and every point of it has been read, so more
// crossings may lie beyond the last and the shape passes again from it.
func (e *edgePass[T]) spent() bool {
	return e.read == len(e.points)
}

// touchedAt reports whether the outline touches the given point: the one the pass started
// from, where a crossing compared Equal to it.
func (e *edgePass[T]) touchedAt(point Point[T]) bool {
	return e.touched && point == e.reached
}

// precedes reports whether a comes before b on the way from the segment's Start: nearer to
// it, or as near and first by Point.Compare.
func (e *edgePass[T]) precedes(a, b Point[T]) bool {
	return cmp.Or(e.segment.compareDistance(a, b), a.Compare(b)) < 0
}

// edgeProbe tests the edges of two outlines against each other one pair at a time, the inner
// loop of every test of an outline against an outline: it is started as a literal with the
// extent a, b of the outline being probed, aimed at a segment, and asked whether each edge of
// that outline meets the segment. A shape ranges the Edges of both sides itself, so the loop
// shares its two prefilters without an iterator crossing a function boundary, which would
// allocate it.
type edgeProbe[T Number] struct {
	a, b    Point[T]
	segment Segment[T]
	c, d    Point[T]
}

// aim points the probe at the segment and reports whether the segment's extent overlaps the
// outline's, so a segment that cannot reach the outline is rejected before any edge is examined.
func (p *edgeProbe[T]) aim(segment Segment[T]) bool {
	p.segment = segment
	p.c, p.d = segment.minMax()

	return overlaps(p.a, p.b, p.c, p.d)
}

// meets reports whether the edge shares a point with the aimed segment, by
// Segment.IntersectsSegment, run only where the extents of the two can share a point.
func (p *edgeProbe[T]) meets(edge Segment[T]) bool {
	e, f := edge.minMax()

	return overlaps(p.c, p.d, e, f) && p.segment.IntersectsSegment(edge)
}

// edgeConvexity folds the edge directions of an outline one at a time into the test IsConvex
// makes, as edgeWalk folds a walk: the first and the previous direction, the sense of the
// turns so far, and how often the X component of the direction has changed sign. The turns of
// an outline that goes around once in one sense change it exactly twice.
type edgeConvexity struct {
	first, previous Vector[float64]
	turned          bool
	clockwise       bool
	x               float64
	reversals       int
}

// step folds one edge direction in and reports whether the outline can still be convex: false
// where it turns against the turns before it or doubles back. A zero direction, the edge of a
// repeated vertex, is skipped. A NaN turn is neither sense, so it never settles one.
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

// result closes the outline with the turn from the last direction back to the first and
// reports whether it is convex: it turned at all, and the X component changed sign no more
// than twice around it.
func (c *edgeConvexity) result() bool {
	return c.step(c.first) && c.turned && c.reversals <= 2
}

// edgeClip cuts an outline to the inner side of one line after another, the steps of
// Sutherland–Hodgman clipping: the lines through the edges of a convex region, which leave the
// part of the outline inside it, or a single line. It is started as a literal with the slice
// to append to, its length as from, and the vertices of the outline, which it reads and never
// writes; each cut appends the outline it leaves after the one it read and moves it down to
// from, so the slice needs room for two outlines and holds the last one at the end. A region
// wound against Rectangle.Vertices is reversed, so the inner side is the same for both.
type edgeClip[T Number] struct {
	dst      []Point[T]
	from     int
	points   []Point[T]
	reversed bool
}

// add cuts the outline to the inner side of the line through one edge of the region, the side
// Segment.Normal points to on an edge of a Rectangle.
func (c *edgeClip[T]) add(edge Segment[T]) {
	if c.reversed {
		edge = edge.Reverse()
	}

	c.cut(edge.Line())
}

// cut keeps the part of the outline on the side of the line Line.Side gives 1 for, the line
// itself included within the tolerance: every vertex there, and between a vertex strictly on
// one side and the next strictly on the other the point where the edge crosses the line. A
// vertex on the line is kept as it is and no crossing is placed beside it, so a point the
// tolerance admits stays where it was. A line with no direction has no sides and cuts nothing.
func (c *edgeClip[T]) cut(line Line[T]) {
	if !line.Direction.hasDirection() || len(c.points) == 0 {
		return
	}

	start := len(c.dst)
	previous := c.points[len(c.points)-1]
	previousSide := c.side(line, previous)
	for _, point := range c.points {
		side := c.side(line, point)
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

// keep appends a vertex to the outline a cut is leaving, the one that began at start, unless it
// compares Equal to the vertex before it. The first vertex kept makes room for the outline a
// convex one leaves, one vertex more than it had.
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

// side returns the side of the line the point lies on as Line.Side reads it, the cross product
// itself, and zero for a point the line contains within the tolerance. A NaN has no side.
func (*edgeClip[T]) side(line Line[T], point Point[T]) float64 {
	if line.Contains(point) {
		return 0
	}

	return line.side(point)
}

// crossing returns the point where the edge from a to b crosses the line, given the sides of
// the two, of opposite signs: each component of the edge times the side of a, divided once by
// the difference of the sides, as Segment.crossing places a point, so for an integer T a
// crossing on a half unit is exactly there before it is rounded.
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

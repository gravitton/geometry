package geom

import (
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
// the point itself exactly where DistanceSquaredTo is zero. The three readers take the walk by
// value, so a shape reads them straight off the walk it returns.
type edgeWalk[T Number] struct {
	inside   bool
	distance float64
	edge     Segment[T]
}

// step folds one edge in and reports whether the point lies on it within Epsilon of T, as
// Segment.DistanceSquaredTo snaps it, where the walk is over: the point is on the boundary at
// distance zero whatever the remaining edges say.
func (w *edgeWalk[T]) step(edge Segment[T], point Point[T]) bool {
	distance := edge.distanceSquaredTo(point)
	if lessOrEqualSquared[T](distance, 0) {
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
// nearer than the radius within Epsilon of T, so a circle touching an edge from inside is
// enclosed. A point on the boundary is at distance zero and clears only a radius within the
// tolerance.
func (w edgeWalk[T]) clears(radius float64) bool {
	return w.result() == 0 && greaterOrEqualSquared[T](w.distance, radius)
}

// edgeIntersections collects the points where a segment crosses the edges of an outline, fed
// one at a time by the shape ranging its own Edges, as edgeWalk folds a walk: each crossing by Segment.IntersectionSegment, a vertex hit
// by two edges counted once, allocated on the first crossing with room for the two a convex
// outline can have. Only a concave outline grows it.
type edgeIntersections[T Number] struct {
	segment Segment[T]
	points  []Point[T]
}

// add folds one edge in, keeping its crossing with the segment unless a previous edge gave
// a point that compares Equal to it.
func (e *edgeIntersections[T]) add(edge Segment[T]) {
	point, ok := e.segment.IntersectionSegment(edge)
	if !ok || slices.ContainsFunc(e.points, point.Equal) {
		return
	}

	if e.points == nil {
		e.points = make([]Point[T], 0, 2)
	}

	e.points = append(e.points, point)
}

// sorted returns the crossings from the segment's Start, and nil where there was none.
func (e *edgeIntersections[T]) sorted() []Point[T] {
	slices.SortFunc(e.points, func(a, b Point[T]) int {
		return e.segment.compareDistance(a, b)
	})

	return e.points
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

// meets reports whether the edge shares a point with the aimed segment, by Segment.IntersectsSegment,
// run only where the extents of the two can share a point.
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

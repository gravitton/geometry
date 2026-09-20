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
// far, and whether a ray from the point has crossed an odd number of edges. A shape starts it
// as a literal, ranges its own Edges and calls step on each, so the walk shares its rule
// without an iterator crossing a function boundary, which would allocate it.
type edgeWalk[T Number] struct {
	inside   bool
	distance float64
}

// step folds one edge in and reports whether the point lies on it within Epsilon of T, as
// Segment.DistanceSquaredTo snaps it, where the walk is over: the point is on the boundary at
// distance zero whatever the remaining edges say.
func (w *edgeWalk[T]) step(edge Segment[T], point Point[T]) bool {
	w.distance = min(w.distance, edge.distanceSquaredTo(point))
	if lessOrEqualSquared[T](w.distance, 0) {
		return true
	}

	if edge.crossesRay(point) {
		w.inside = !w.inside
	}

	return false
}

// result returns the squared distance from the point to the outline after every edge: zero
// inside by the even-odd rule, the distance to the nearest edge otherwise, and infinity where
// there was no edge.
func (w *edgeWalk[T]) result() float64 {
	if w.inside {
		return 0
	}

	return w.distance
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

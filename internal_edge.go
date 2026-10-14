package geom

import (
	"iter"
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

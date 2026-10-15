package geom

import (
	"cmp"
	"slices"
)

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

// bezierCast gathers the hits of a segment at a Bezier curve in order from its Start, each normal
// the perpendicular of the curve where it is hit, a point two pieces share held once. A cubic is
// crossed three times at most; the rest of the room holds a crossing both neighbouring pieces find,
// and the last slot takes the farthest hit when there is no room, outside count.
//
// Pieces met one after another along the curve are a run. Where the chords of a run only touch the
// segment, both ends on one side of it, the curve passes it within the tolerance, so the run holds
// its first touch, and none once it holds a crossing.
type bezierCast struct {
	segment  Segment[float64]
	curve    CubicBezier[float64]
	hits     [7]Hit[float64]
	count    int
	reached  float64
	crossed  bool
	touching bool
	touch    Point[float64]
}

// step casts the segment at the piece between the two parameters: nothing where it misses the box
// of the control points, the crossing of the chord where the piece is straight within the
// tolerance of its magnitude or the depth is spent, and its two halves otherwise.
func (c *bezierCast) step(piece CubicBezier[float64], from, to float64, depth int) {
	a, b := minMaxOf([]Point[float64]{piece.Start, piece.Control1, piece.Control2, piece.End})
	box := Box[float64]{a, b}
	if !box.IntersectsSegment(c.segment) {
		return
	}

	if depth == 0 || piece.isFlatWithin(epsilonAt[float64](box.magnitude())) {
		chord := Segment[float64]{piece.Start, piece.End}
		if point, ok := c.segment.IntersectionSegment(chord); ok {
			t := from + float64((to-from)*Clamp(chord.fraction(point), 0, 1))

			c.meet(c.segment.hitAt(point, c.curve.TangentAt(t).Normal()), chord, from, to)
		}

		return
	}

	first, second := piece.Split(0.5)
	middle := (from + to) / 2

	c.step(first, from, middle, depth-1)
	c.step(second, middle, to, depth-1)
}

// meet holds the hit of the chord of the piece between the two parameters, a crossing always and a
// touch as one of its run.
func (c *bezierCast) meet(hit Hit[float64], chord Segment[float64], from, to float64) {
	if from != c.reached {
		c.crossed, c.touching = false, false
	}

	c.reached = to

	line := c.segment.Line()
	start, end := line.side(chord.Start), line.side(chord.End)

	if start < 0 && end > 0 || start > 0 && end < 0 {
		c.release()
		c.crossed = true
		c.add(hit)

		return
	}

	if c.crossed || c.touching {
		return
	}

	c.touching, c.touch = true, hit.Point
	c.add(hit)
}

// release lets the touch of the run go, where it is still held.
func (c *bezierCast) release() {
	if !c.touching {
		return
	}

	c.touching = false

	for i, held := range c.hits[:c.count] {
		if held.Point == c.touch {
			copy(c.hits[i:], c.hits[i+1:c.count])
			c.count--

			return
		}
	}
}

// add holds the hit in its place by Fraction unless its point is already held. Where there is no
// room the farthest of the held hits and the new one gives way.
func (c *bezierCast) add(hit Hit[float64]) {
	at := c.count
	for i, held := range c.hits[:c.count] {
		if held.Point.Equal(hit.Point) {
			return
		}

		if hit.Fraction < held.Fraction && at == c.count {
			at = i
		}
	}

	copy(c.hits[at+1:], c.hits[at:c.count])
	c.hits[at] = hit
	c.count = min(c.count+1, len(c.hits)-1)
}

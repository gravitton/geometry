package geom

import (
	"slices"
)

// unionOutline gathers the outline of the union of two convex polygons, appended to dst after from:
// the vertices of each outside the other and the crossings of their edges. It first averages the
// points both share, a point inside both about which the boundary of each runs in the order of the
// angle, then walks the first polygon, its vertices and the crossings along each edge, and puts each
// vertex of the second where its angle falls. The angles are compared on the crossings before they
// are rounded, so an integer crossing cannot pass a vertex beside it.
type unionOutline[T Number] struct {
	dst            []Point[T]
	from           int
	a, b           Polygon[T]
	start          Point[T]
	started        bool
	sumX, sumY     float64
	count          int
	first, other   Point[float64]
	paired, spread bool
	center         Point[float64]
	reference      Vector[float64]
	mirrored       bool
	next, step     int
	left           int
}

// gather averages the points both polygons share and notes the point the outline starts at: the
// first vertex of the first polygon outside the second, else the first of the second outside the
// first, else the first crossing.
func (u *unionOutline[T]) gather() {
	for _, point := range u.a.Points {
		u.take(point, u.b.Contains(point))
	}
	for _, point := range u.b.Points {
		u.take(point, u.a.Contains(point))
	}
	for edge := range u.a.Edges() {
		for other := range u.b.Edges() {
			if point, ok := edge.intersectionSegment(other); ok {
				u.share(point)
				u.begin(point.Cast[T]())
			}
		}
	}
}

// take adds a vertex among the shared points inside the other polygon, or as a start outside it.
func (u *unionOutline[T]) take(point Point[T], inside bool) {
	if inside {
		u.share(point.Float())
	} else {
		u.begin(point)
	}
}

// begin notes the point as the start where none is noted yet.
func (u *unionOutline[T]) begin(point Point[T]) {
	if !u.started {
		u.start, u.started = point, true
	}
}

// share adds a point both polygons hold to the average, and notes on exact signs whether the shared
// points so far enclose an area.
func (u *unionOutline[T]) share(point Point[float64]) {
	u.sumX += point.X
	u.sumY += point.Y
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

// walk appends the outline in the winding of the first polygon, starting at the noted point.
// Shared points that enclose no area leave dst as it was given.
func (u *unionOutline[T]) walk() []Point[T] {
	if !u.spread {
		return u.dst[:u.from]
	}

	u.orient()
	for i, point := range u.a.Points {
		if !u.b.Contains(point) {
			u.add(point.Float(), point)
		}

		u.cross(Segment[T]{point, u.a.Points[(i+1)%len(u.a.Points)]})
	}
	for u.left > 0 {
		u.pass()
	}

	return u.result()
}

// orient sets the center the angles are measured about, the direction they are measured from, the
// first vertex of the first polygon, and the vertex of the second the walk of it starts at, the
// first from that direction, stepping the way the first polygon winds. A repeated vertex is walked
// from its first copy, so none is left behind at the end.
func (u *unionOutline[T]) orient() {
	u.center = Point[float64]{u.sumX / float64(u.count), u.sumY / float64(u.count)}
	u.mirrored = u.a.twiceArea() < 0
	u.reference = u.offset(u.a.Points[0].Float())

	u.step, u.left = 1, len(u.b.Points)
	if (u.b.twiceArea() < 0) != u.mirrored {
		u.step = -1
	}
	for i, point := range u.b.Points {
		if u.precedes(point.Float(), u.b.Points[u.next].Float()) {
			u.next = i
		}
	}
	for range len(u.b.Points) - 1 {
		previous := (u.next - u.step + len(u.b.Points)) % len(u.b.Points)
		if u.precedes(u.b.Points[u.next].Float(), u.b.Points[previous].Float()) {
			return
		}

		u.next = previous
	}
}

// cross adds the crossings along the edge of the first polygon in order from its start, each found
// as the nearest past the one before. Its end is left to the edge that starts there.
func (u *unionOutline[T]) cross(edge Segment[T]) {
	start, along := edge.Start.Float(), edge.Float().Vector()
	reach := along.LengthSquared()

	last := -1.0
	for {
		nearest, found := Point[float64]{}, reach
		for other := range u.b.Edges() {
			if point, ok := edge.intersectionSegment(other); ok {
				if at := point.Subtract(start).Dot(along); last < at && at < found {
					nearest, found = point, at
				}
			}
		}
		if found == reach {
			return
		}

		u.add(nearest, nearest.Cast[T]())
		last = found
	}
}

// add appends a point of the walk of the first polygon, after the vertices of the second whose
// angle comes first.
func (u *unionOutline[T]) add(exact Point[float64], point Point[T]) {
	for u.left > 0 && u.precedes(u.b.Points[u.next].Float(), exact) {
		u.pass()
	}

	u.append(point)
}

// pass appends the next vertex of the second polygon where it lies outside the first, and moves on
// to the one after it.
func (u *unionOutline[T]) pass() {
	if vertex := u.b.Points[u.next]; !u.a.Contains(vertex) {
		u.append(vertex)
	}

	u.next = (u.next + u.step + len(u.b.Points)) % len(u.b.Points)
	u.left--
}

// append adds the point unless it is Equal to the one before.
func (u *unionOutline[T]) append(point Point[T]) {
	if len(u.dst) == u.from || !point.Equal(u.dst[len(u.dst)-1]) {
		u.dst = append(u.dst, point)
	}
}

// precedes reports whether the angle of a comes before that of b, both measured from the reference
// the way the first polygon winds.
func (u *unionOutline[T]) precedes(a, b Point[float64]) bool {
	va, vb := u.offset(a), u.offset(b)
	if wrapA, wrapB := compareAngles(va, u.reference) < 0, compareAngles(vb, u.reference) < 0; wrapA != wrapB {
		return wrapB
	}

	return compareAngles(va, vb) < 0
}

// offset returns the point less the center, mirrored across X for a polygon wound clockwise, so its
// walk runs the way angles grow.
func (u *unionOutline[T]) offset(point Point[float64]) Vector[float64] {
	offset := point.Subtract(u.center)
	if u.mirrored {
		return Vector[float64]{offset.X, -offset.Y}
	}

	return offset
}

// result drops a last point Equal to the first and starts the outline at the noted point.
func (u *unionOutline[T]) result() []Point[T] {
	outline := u.dst[u.from:]
	if last := len(outline) - 1; last > 0 && outline[last].Equal(outline[0]) {
		outline = outline[:last]
	}

	if i := slices.IndexFunc(outline, u.start.Equal); i > 0 {
		slices.Reverse(outline[:i])
		slices.Reverse(outline[i:])
		slices.Reverse(outline)
	}

	return u.dst[:u.from+len(outline)]
}

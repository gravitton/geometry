package geom

const (
	// simplexSteps bounds the supports nearest asks for. Two outlines end long before it, on a
	// vertex met twice, and two round shapes are within the tolerance of their distance by then.
	simplexSteps = 64

	// sweepSteps bounds the advances Sweep makes. An advance toward a flat edge lands on it at
	// once, and one toward a round shape closes most of what is left each time.
	sweepSteps = 64

	// facingAngle is how far to either side of a normal the ends of the edge facing it are asked
	// for: an edge within it of facing the normal counts as flat against it, as a box a little
	// tilted on the ground does, and gives two points of contact.
	facingAngle = 1e-3
)

// minkowski is the difference of two shapes taken as convex, the first moved by an offset:
// every vector from a point of the second to a point of the moved first. It holds the origin
// exactly where the two share a point, its point nearest the origin is the gap between them,
// and the way out of it from the origin is the way to part them, so Intersects, DistanceSquared,
// Nearest, Collide and Sweep of two shapes known only by their Support are read off it. Everything is
// measured in float64, whatever T is.
//
// With no first shape the difference is that of the single point at the offset, the origin of
// a ray, which Raycast sweeps without a shape to hold it.
type minkowski[T Number] struct {
	a, b   Shape[T]
	offset Vector[float64]
}

// support returns the point of the difference farthest along the given direction: the support
// of the first shape that way, or the point at the offset where there is none, less the support
// of the second the other way.
func (m minkowski[T]) support(direction Vector[float64]) minkowskiVertex {
	from := m.b.Support(direction.Negate())
	if m.a == nil {
		return minkowskiVertex{m.offset.Point().Subtract(from), from}
	}

	return minkowskiVertex{m.a.Support(direction).Add(m.offset).Subtract(from), from}
}

// magnitude returns the largest absolute coordinate of the two shapes, that of a corner of
// their Bounds, or of the point at the offset where there is no first shape: the size
// epsilonAt widens the tolerance of the pair by.
func (m minkowski[T]) magnitude() float64 {
	if m.a == nil {
		return max(m.offset.Point().magnitude(), m.b.Bounds().magnitude())
	}

	return max(m.a.Bounds().magnitude(), m.b.Bounds().magnitude())
}

// tolerance returns the distance within which the difference is taken to be closed on: the
// tolerance of T at the magnitude of the pair, and that of float64 for an integer T, whose own
// is zero, since a gap measured in float64 never closes exactly.
func (m minkowski[T]) tolerance() float64 {
	magnitude := m.magnitude()

	return max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude))
}

// intersects reports whether the two shapes share a point, a touch included within the
// tolerance of T: lessOrEqualSquared on the gap nearest finds, the comparison every other pair
// makes. For an integer T the tolerance is zero, so the gap must close exactly, which the
// outlines of two shapes with vertices do and a round one does not in general.
func (m minkowski[T]) intersects() bool {
	return lessOrEqualSquared(m.nearest().point.LengthSquared(), 0, epsilonAt[T](m.magnitude()))
}

// nearest returns the point of the difference nearest the origin, the zero vector where it
// holds the origin, by the Gilbert–Johnson–Keerthi walk: a simplex of supports is moved toward
// the origin, each new support asked for in the direction from the nearest point of the simplex
// to the origin, until the support is one it already holds or reaches no nearer than that point
// within Delta of its distance. A NaN ends it at once, since it reaches no nearer than
// anything.
func (m minkowski[T]) nearest() minkowskiVertex {
	simplex := minkowskiSimplex{}
	simplex.add(m.support(Vector[float64]{1, 0}))

	nearest := simplex.nearest()
	for range simplexSteps {
		next := m.support(nearest.point.Negate())

		distanceSquared := nearest.point.LengthSquared()
		if simplex.holds(next) || !(distanceSquared-float64(nearest.point.Dot(next.point)) > Delta*distanceSquared) {
			break
		}

		simplex.add(next)
		nearest = simplex.nearest()
	}

	return nearest
}

// depth returns the unit direction that moves the first shape out of the second the shortest
// way and how far along it, zero where they do not overlap, by the expanding polytope: a
// polygon inscribed in the difference is grown at the edge nearest the origin by the support
// along the normal of that edge, until the support reaches no farther than the edge within the
// tolerance, so the edge lies on the boundary of the difference, or the polytope is full. The
// way from the origin out along that normal as far as the support reaches leaves the difference
// whatever the polytope has grown to, so its opposite moves the first shape off the second, and
// by the least distance once the edge is on the boundary.
//
// The polygon starts from the supports along the four axes, in the order that winds as
// Vector.Normal turns, which may hold the origin or not: an edge the origin lies beyond has a
// negative distance, the least of all, so it is grown first, and the polygon holds the origin
// before any edge of the boundary is reached wherever the difference does. A difference of no
// area, two parallel segments, grows nowhere and has no depth, and that of two points has no
// direction either.
func (m minkowski[T]) depth() (Vector[float64], float64) {
	tolerance := m.tolerance()

	polytope := minkowskiPolytope{}
	for _, direction := range [4]Vector[float64]{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
		polytope.add(m.support(direction).point)
	}

	depth := 0.0

	edge, found := polytope.nearest()
	for found {
		point := m.support(edge.normal).point

		depth = point.Dot(edge.normal)
		if polytope.count == len(polytope.vertices) || !(depth-edge.distance > tolerance) {
			break
		}

		polytope.insert(edge.index+1, point)
		edge, _ = polytope.nearest()
	}

	return edge.normal.Negate(), max(depth, 0)
}

// sweep returns where the first shape first touches the second while it moves by the given
// motion from the offset, where it does not intersect it, and false where it reaches nothing,
// by conservative advancement: it is moved by the gap between the two over the speed they
// close at along it, which can never pass the other, and no farther than the whole motion,
// until the gap is within the tolerance. The normal is that of the last advance, since a closed
// gap has no direction left.
func (m minkowski[T]) sweep(motion Vector[float64]) (Hit[T], bool) {
	tolerance := m.tolerance()
	start := m.offset

	fraction := 0.0
	nearest := m.nearest()
	for range sweepSteps {
		normal := nearest.point.Normalize()

		approach := -motion.Dot(normal)
		if !(approach > 0) || fraction == 1 {
			break
		}

		fraction = min(fraction+nearest.point.Length()/approach, 1)

		m.offset = start.Add(motion.Multiply(fraction))
		nearest = m.nearest()
		if lessOrEqualSquared(nearest.point.LengthSquared(), 0, tolerance) {
			return Hit[T]{nearest.from.Cast[T](), normal, fraction}, true
		}
	}

	return Hit[T]{}, false
}

// contact returns the points of the second shape the two meet at, given the unit normal that
// parts them, from the second toward the first: the part of the edge the second turns to the
// first that the edge the first turns to it lies over, both ends of it where two flat edges
// meet and one point where a corner or a round boundary does. Edges that pass each other
// without lying over one another, as the corners of two shapes deep in each other can, meet
// halfway between them.
func (m minkowski[T]) contact(normal Vector[float64]) Pair[Point[T]] {
	tangent := normal.Normal()
	facing, faced := facingEdge(m.a, normal.Negate()), facingEdge(m.b, normal)

	a1, a2 := facing.Start.Vector().Dot(tangent), facing.End.Vector().Dot(tangent)
	b1, b2 := faced.Start.Vector().Dot(tangent), faced.End.Vector().Dot(tangent)

	low, high := max(min(a1, a2), min(b1, b2)), min(max(a1, a2), max(b1, b2))
	if low > high {
		low = (low + high) / 2
		high = low
	}

	first := faced.PointAt(Clamp(ratio(low-b1, b2-b1), 0, 1)).Cast[T]()
	second := faced.PointAt(Clamp(ratio(high-b1, b2-b1), 0, 1)).Cast[T]()
	if first.Equal(second) {
		return Pair[Point[T]]{[2]Point[T]{first}, 1}
	}

	return Pair[Point[T]]{[2]Point[T]{first, second}, 2}
}

// facingEdge returns the edge of a shape turned toward the given direction, of no length for
// a corner or a round boundary: the supports a little to either side of the direction, which
// are the two ends of a flat edge that faces it and the same vertex twice at a corner. The
// supports at half that turn tell a flat edge from a round boundary: across a flat edge the
// support jumps from one end to the other between them, farther than it moves outside them,
// where around a round boundary it moves as far in each half. A flat side between two round
// ends, that of a Capsule, is therefore an edge, and a round boundary is the one point farthest
// along the direction.
func facingEdge[T Number](shape Shape[T], direction Vector[float64]) Segment[float64] {
	start, end := shape.Support(direction.Rotate(facingAngle)), shape.Support(direction.Rotate(-facingAngle))
	from, to := shape.Support(direction.Rotate(facingAngle/2)), shape.Support(direction.Rotate(-facingAngle/2))
	if from.DistanceTo(to) <= 2*(start.DistanceTo(from)+end.DistanceTo(to)) {
		point := shape.Support(direction)

		return Segment[float64]{point, point}
	}

	return Segment[float64]{start, end}
}

// minkowskiVertex is a point of the difference with the point of the second shape it was
// measured from, the one a contact is placed at.
type minkowskiVertex struct {
	point Vector[float64]
	from  Point[float64]
}

// nearestToward returns the point nearest the origin on the edge from the vertex to the given
// one, with the simplex of the vertices it lies between: an end where the origin projects past
// it, and otherwise the foot of the perpendicular from the origin, placed along the normal of
// the edge by the cross product of its ends, so an edge through the origin is at exactly zero.
// An edge of no length is its first vertex.
func (v minkowskiVertex) nearestToward(vertex minkowskiVertex) (minkowskiVertex, minkowskiSimplex) {
	edge := vertex.point.Subtract(v.point)
	lengthSquared := edge.LengthSquared()

	t := ratio(-v.point.Dot(edge), lengthSquared)
	if !(t > 0) {
		return v, minkowskiSimplex{[3]minkowskiVertex{v}, 1}
	}

	if t >= 1 {
		return vertex, minkowskiSimplex{[3]minkowskiVertex{vertex}, 1}
	}

	normal := edge.Normal()
	foot := minkowskiVertex{normal.Multiply(v.point.Dot(normal) / lengthSquared), v.from.Lerp(vertex.from, t)}

	return foot, minkowskiSimplex{[3]minkowskiVertex{v, vertex}, 2}
}

// minkowskiSimplex is the simplex nearest moves toward the origin: a point, an edge or a
// triangle of the difference, reduced after every step to the vertices its nearest point lies
// between.
type minkowskiSimplex struct {
	vertices [3]minkowskiVertex
	count    int
}

// add appends a vertex to the simplex, which has room for it after every nearest.
func (s *minkowskiSimplex) add(vertex minkowskiVertex) {
	s.vertices[s.count] = vertex
	s.count++
}

// holds reports whether the simplex already has a vertex at the point of the given one: a
// support met again, past which the simplex moves no nearer.
func (s *minkowskiSimplex) holds(vertex minkowskiVertex) bool {
	for _, held := range s.vertices[:s.count] {
		if held.point == vertex.point {
			return true
		}
	}

	return false
}

// nearest returns the point of the simplex nearest the origin and reduces the simplex to the
// vertices it lies between: the vertex of a point, the nearest of an edge, and for a triangle
// the nearest of its three edges, moved to the origin itself where the triangle holds it,
// which is decided on the exact signs of the cross products of its vertices, a vertex or an
// edge through the origin included. The point of the second shape is then the one the vertices
// were measured from weighted as the origin lies between them, a point the two shapes share,
// which an edge through the origin already has.
func (s *minkowskiSimplex) nearest() minkowskiVertex {
	if s.count == 1 {
		return s.vertices[0]
	}

	a, b, c := s.vertices[0], s.vertices[1], s.vertices[2]

	nearest, kept := a.nearestToward(b)
	if s.count == 2 {
		*s = kept

		return nearest
	}

	for _, edge := range [2][2]minkowskiVertex{{b, c}, {c, a}} {
		if candidate, vertices := edge[0].nearestToward(edge[1]); candidate.point.LengthSquared() < nearest.point.LengthSquared() {
			nearest, kept = candidate, vertices
		}
	}
	*s = kept

	ab, bc, ca := a.point.Cross(b.point), b.point.Cross(c.point), c.point.Cross(a.point)
	if (ab >= 0 && bc >= 0 && ca >= 0) || (ab <= 0 && bc <= 0 && ca <= 0) {
		nearest.point = Vector[float64]{}
		if area := ab + bc + ca; area != 0 {
			nearest.from = a.from.Multiply(bc / area).Add(b.from.Multiply(ca / area).Vector()).Add(c.from.Multiply(ab / area).Vector())
		}
	}

	return nearest
}

// minkowskiPolytope is the polygon depth grows inside the difference: its vertices in the
// order that winds as Vector.Normal turns, so the outward side of every edge is the one its
// direction turned the other way points to.
type minkowskiPolytope struct {
	vertices [64]Vector[float64]
	count    int
}

// minkowskiEdge is an edge of the polytope by the index of its first vertex, with its outward
// unit normal and the distance of its line from the origin along that normal, negative where
// the origin lies beyond it.
type minkowskiEdge struct {
	index    int
	normal   Vector[float64]
	distance float64
}

// add appends a vertex to the polytope.
func (p *minkowskiPolytope) add(vertex Vector[float64]) {
	p.vertices[p.count] = vertex
	p.count++
}

// insert places a vertex at the given index, moving the ones from there on up by one.
func (p *minkowskiPolytope) insert(index int, vertex Vector[float64]) {
	copy(p.vertices[index+1:p.count+1], p.vertices[index:p.count])
	p.vertices[index] = vertex
	p.count++
}

// nearest returns the edge whose line is nearest the origin, the one of the least distance
// along its outward normal, by minkowskiEdge.precedes, and false where no edge has a length to
// take a normal from, the polytope of a single point.
func (p *minkowskiPolytope) nearest() (minkowskiEdge, bool) {
	var nearest minkowskiEdge

	found := false
	for i, vertex := range p.vertices[:p.count] {
		direction := p.vertices[(i+1)%p.count].Subtract(vertex)
		if !direction.hasDirection() {
			continue
		}

		normal := direction.Normal().Negate().Normalize()
		if edge := (minkowskiEdge{i, normal, normal.Dot(vertex)}); !found || edge.precedes(nearest) {
			nearest, found = edge, true
		}
	}

	return nearest, found
}

// precedes reports whether the edge is the one to keep of the two: the nearer the origin, and
// of two as near the one whose normal lies closer to the X axis, then the one whose normal
// points between the axes of one sign. Neither choice changes where the normal is turned
// about, so the difference of two shapes taken the other way round, which is this one mirrored
// through the origin, keeps the mirrored edge, and Collide answers with the opposite normal
// whichever shape is given first. Only two edges facing opposite ways are left to the order
// they come in, where a shape lies in the middle of the other and both ways out are as short.
// A NaN distance precedes nothing.
func (e minkowskiEdge) precedes(edge minkowskiEdge) bool {
	if e.distance != edge.distance {
		return e.distance < edge.distance
	}

	if x, other := Abs(e.normal.X), Abs(edge.normal.X); x != other {
		return x > other
	}

	return float64(e.normal.X*e.normal.Y) > float64(edge.normal.X*edge.normal.Y)
}

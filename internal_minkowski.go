package geom

const (
	// simplexSteps bounds the supports nearest asks for. Two outlines end long before it, and two
	// round shapes are within the tolerance by then.
	simplexSteps = 64

	// sweepSteps bounds the advances of sweep. An advance toward a flat edge lands on it at once,
	// and one toward a round shape closes most of what is left.
	sweepSteps = 64

	// facingAngle is how far from facing the normal an edge may turn and still count as flat,
	// giving two contact points.
	facingAngle = 1e-3
)

// minkowski is the difference of two shapes taken as convex, the first moved by the offset. It
// holds the origin exactly where they share a point, its point nearest the origin is the gap, and
// the way out from the origin parts them, so every pair known only by Support is decided on it, in
// float64. With no first shape it is the difference from the point at the offset, the origin of a
// ray.
type minkowski[T Number] struct {
	a, b   Shape[T]
	offset Vector[float64]
}

// support returns the point of the difference farthest along the direction.
func (m minkowski[T]) support(direction Vector[float64]) minkowskiVertex {
	from := m.b.Support(direction.Negate())
	if m.a == nil {
		return minkowskiVertex{m.offset.Point().Subtract(from), from}
	}

	return minkowskiVertex{m.a.Support(direction).Add(m.offset).Subtract(from), from}
}

// magnitude returns the largest absolute coordinate of the pair's Bounds, the magnitude the
// tolerance of the pair is taken from.
func (m minkowski[T]) magnitude() float64 {
	if m.a == nil {
		return max(m.offset.Point().magnitude(), m.b.Bounds().magnitude())
	}

	return max(m.a.Bounds().magnitude(), m.b.Bounds().magnitude())
}

// tolerance returns the distance within which the difference counts as closed on: that of T, or of
// float64 for an integer T, since a gap measured in float64 never closes exactly.
func (m minkowski[T]) tolerance() float64 {
	magnitude := m.magnitude()

	return max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude))
}

// intersects reports whether the shapes share a point within the tolerance, the one every other
// answer of the difference takes, so a round shape touching on the lattice is met from either side.
func (m minkowski[T]) intersects() bool {
	return lessOrEqualSquared(m.nearest().point.LengthSquared(), 0, m.tolerance())
}

// nearest returns the point of the difference nearest the origin, zero where it holds the origin,
// by the Gilbert–Johnson–Keerthi walk: it stops on a support already held or one reaching no nearer
// within Delta. A NaN ends it at once.
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

// depth returns the unit direction that moves the first shape out of the second the shortest way,
// and how far, by the expanding polytope grown at its edge nearest the origin until the support
// reaches no farther within the tolerance or the polytope is full.
//
// It starts from the supports along the four axes, wound as Vector.Normal turns. An edge the origin
// lies beyond has a negative distance and grows first, so the polytope holds the origin before any
// boundary edge is reached. A difference of no area grows nowhere and has no depth.
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

// sweep returns where the moving first shape first touches the second, by conservative advancement:
// each step moves it by the gap over the closing speed, which never passes the other, until the gap
// is within the tolerance. The normal is that of the last advance, since a closed gap has none.
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

// contact returns the points of the second shape the two meet at across the normal: both ends of
// the overlap where two flat edges face each other, one point at a corner or a round boundary.
// Edges that do not overlap, as deep corners can, meet halfway between them.
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

// facingEdge returns the edge of the shape facing the direction, of no length at a corner or a
// round boundary, from the supports a little to either side. Across a flat edge the support jumps
// between the halves of that turn farther than it moves outside them, which tells it from a round
// boundary, so the flat side of a Capsule is an edge.
func facingEdge[T Number](shape Shape[T], direction Vector[float64]) Segment[float64] {
	start, end := shape.Support(direction.Rotate(facingAngle)), shape.Support(direction.Rotate(-facingAngle))
	from, to := shape.Support(direction.Rotate(facingAngle/2)), shape.Support(direction.Rotate(-facingAngle/2))
	if from.DistanceTo(to) <= 2*(start.DistanceTo(from)+end.DistanceTo(to)) {
		point := shape.Support(direction)

		return Segment[float64]{point, point}
	}

	return Segment[float64]{start, end}
}

// minkowskiVertex is a point of the difference with the point of the second shape it was measured
// from, where a contact is placed.
type minkowskiVertex struct {
	point Vector[float64]
	from  Point[float64]
}

// nearestToward returns the point nearest the origin on the edge to the given vertex, with the
// simplex of the vertices it lies between. The foot is placed along the normal by the cross
// product, so an edge through the origin is at exactly zero; an edge of no length is its first
// vertex.
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

// minkowskiSimplex is the point, edge or triangle nearest walks toward the origin, reduced after
// each step to the vertices its nearest point lies between.
type minkowskiSimplex struct {
	vertices [3]minkowskiVertex
	count    int
}

// add appends a vertex, which nearest always leaves room for.
func (s *minkowskiSimplex) add(vertex minkowskiVertex) {
	s.vertices[s.count] = vertex
	s.count++
}

// holds reports whether the simplex already has a vertex at the point of the given one, a support
// met again, past which the walk gets no nearer.
func (s *minkowskiSimplex) holds(vertex minkowskiVertex) bool {
	for _, held := range s.vertices[:s.count] {
		if held.point == vertex.point {
			return true
		}
	}

	return false
}

// nearest returns the point of the simplex nearest the origin and reduces the simplex to the
// vertices it lies between. A triangle holds the origin on the exact signs of its cross products,
// its boundary included, and the shared point is then the points measured from, weighted as the
// origin lies between them.
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

// minkowskiPolytope is the polygon depth grows inside the difference, wound as Vector.Normal turns,
// so every edge's outward side is its direction turned the other way.
type minkowskiPolytope struct {
	vertices [64]Vector[float64]
	count    int
}

// minkowskiEdge is a polytope edge by the index of its first vertex, with its outward unit normal
// and the distance of its line from the origin, negative where the origin lies beyond it.
type minkowskiEdge struct {
	index    int
	normal   Vector[float64]
	distance float64
}

// add appends a vertex.
func (p *minkowskiPolytope) add(vertex Vector[float64]) {
	p.vertices[p.count] = vertex
	p.count++
}

// insert places a vertex at the index, moving the rest up by one.
func (p *minkowskiPolytope) insert(index int, vertex Vector[float64]) {
	copy(p.vertices[index+1:p.count+1], p.vertices[index:p.count])
	p.vertices[index] = vertex
	p.count++
}

// nearest returns the edge whose line is nearest the origin by precedes, and false where no edge
// has a length, the polytope of a single point.
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

// precedes reports whether the edge is the one to keep: the nearer, then the one whose normal lies
// nearer the X axis, then the one pointing between axes of one sign. Each choice is symmetric under
// mirroring through the origin, so Collide gives the opposite normal with the shapes swapped; only
// a shape in the very middle of another is left to the order. A NaN distance precedes nothing.
func (e minkowskiEdge) precedes(edge minkowskiEdge) bool {
	if e.distance != edge.distance {
		return e.distance < edge.distance
	}

	if x, other := Abs(e.normal.X), Abs(edge.normal.X); x != other {
		return x > other
	}

	return float64(e.normal.X*e.normal.Y) > float64(edge.normal.X*edge.normal.Y)
}

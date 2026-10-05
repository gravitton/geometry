package geom

import (
	"fmt"
	"iter"
	"math"
)

// RegularPolygon is a polygon with equally spaced vertices around a center.
//
// Size holds the semi-axes of the Ellipse the vertices lie on, the one Ellipse() gives back,
// so it is a radius, not an extent: the polygon inscribed in a circle of radius r has Size
// r x r and spans up to twice it. This differs from Rectangle, whose Size is the full width
// and height. Use Bounds for the extent.
//
// Phase places the vertices on the ellipse of the semi-axes, the first at that parameter and each
// next a step of 2π/n on, and Angle then turns the whole ring about the center. The two differ only
// for unequal semi-axes: the phase decides where the vertices sit against the axes the ellipse is
// stretched along, so a flat-top polygon is stretched across its edges and a pointy-top one across
// its vertices, and the turn moves that stretched shape. Where the semi-axes are equal, turning and
// stepping around the ring are the same thing and the first vertex lies at Angle + Phase. Every
// affine image of a regular polygon is one of these, which is what lets Transform keep any matrix
// exactly.
//
// The size is never negative: RegPol and RegularPolygonWithOrientation take it absolute, and
// Scale takes a negative factor absolute, since a negative semi-axis mirrors the vertices
// rather than describe a different kind of polygon. A negative size can only be written as a
// struct literal or decoded from JSON; Canonical takes it absolute.
//
// For an integer T the vertices are rounded, and where the semi-axes differ or the ring is
// turned the rounded outline need not be convex. Contains follows that outline exactly, but the
// Encloses methods test the points of the other shape alone, as for a convex container, so they
// can admit a shape passing through a dent the rounding made.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	N      int      `json:"n"`
	Angle  float64  `json:"a,omitzero"`
	Phase  float64  `json:"p,omitzero"`
}

// RegPol is shorthand for RegularPolygon{center, size, n, angle, phase}, with the size taken
// absolute.
func RegPol[T Number](center Point[T], size Size[T], n int, angle, phase float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, size.Abs(), n, angle, phase}
}

// RegularPolygonWithOrientation creates a RegularPolygon with the given orientation, its first
// vertex placed by RegularPolygonOrientationPhase and not turned, with the size taken absolute
// like RegPol: Width is the semi-axis across and Height the one up, for every n.
func RegularPolygonWithOrientation[T Number](center Point[T], size Size[T], n int, orientation Orientation) RegularPolygon[T] {
	return RegPol(center, size, n, 0, RegularPolygonOrientationPhase(n, orientation))
}

// Anchor returns the point of the boundary in the given direction from the center, or the center
// itself for DirectionNone: the point where the ray from the center leaves the polygon. The
// direction is taken in the world, where the polygon stands after its turn, so Rotate moves every
// anchor along the boundary: on a flat-top polygon the Top anchor is the midpoint of the top edge,
// on a pointy-top one its top vertex. For integer T the point of the edge between the rounded
// vertices is rounded once more, which can carry it off that edge by up to half a diagonal, where
// Contains rejects it. A polygon with fewer than three vertices encloses no area and anchors at its
// center, as does one whose semi-axis lies along the ray.
func (rp RegularPolygon[T]) Anchor(direction Direction) Point[T] {
	if direction.IsNone() || rp.N < 3 {
		return rp.Center
	}

	center := rp.Center.Float()
	ray := VectorFromAngle(direction.Angle(), 1.0)
	edge := rp.edge(rp.edgeIndex(direction.Angle() - rp.Angle)).Float()

	denominator := ray.Cross(edge.Vector())
	if denominator == 0 {
		return rp.Center
	}

	along := Clamp(edge.Start.Subtract(center).Cross(ray)/denominator, 0, 1)

	return edge.PointAt(along).Cast[T]()
}

// Vertices iterates the polygon vertices in order starting from Phase, by increasing angle —
// the same winding as Directions and Rectangle.Vertices, and clockwise as drawn on a screen
// with Y pointing down, without allocating; collect them with slices.Collect where a slice is
// needed. A polygon with N < 1 has no vertices and yields nothing.
// For integer T, each vertex component is rounded to the nearest integer, so vertices at
// non-right angles may be off by up to half a unit. Use float64 for exact positions.
func (rp RegularPolygon[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for i := range rp.N {
			if !yield(rp.vertex(i)) {
				return
			}
		}
	}
}

// Edges iterates the polygon edges in vertex order, each from a vertex to the next and the
// last one closing back to the first, the edges Polygon().Edges() iterates, without building
// the vertices. A polygon with N < 1 has no edges and yields nothing; one with a single
// vertex yields one zero-length edge.
func (rp RegularPolygon[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		if rp.IsEmpty() {
			return
		}

		start := rp.vertex(0)
		for i, previous := 1, start; i <= rp.N; i++ {
			next := start
			if i < rp.N {
				next = rp.vertex(i)
			}

			if !yield(Segment[T]{previous, next}) {
				return
			}

			previous = next
		}
	}
}

// Centroid returns the center of the enclosed area, the Center of the polygon.
func (rp RegularPolygon[T]) Centroid() Point[T] {
	return rp.Center
}

// Area returns the area enclosed by the polygon, in closed form: n/2 · w · h · sin(2π/n), the
// area of the polygon inscribed in the ellipse with the semi-axes Size holds, without building
// the vertices. A polygon with N < 3 encloses no area, as the sine of a full or a half turn
// gives it. For integer T the vertices are rounded but the area is not: it is the area of the
// exact polygon, which Polygon().Area() measures from the rounded vertices.
func (rp RegularPolygon[T]) Area() float64 {
	if rp.N < 3 {
		return 0
	}

	n := float64(rp.N)

	return n / 2 * float64(rp.Size.Width) * float64(rp.Size.Height) * math.Sin(2*Pi/n)
}

// Perimeter returns the total length of the edges without building the vertices. On a square
// Size it is the closed form 2 · n · r · sin(π/n); on an ellipse the chords differ in length
// and are summed from the angles, since the perimeter of a polygon inscribed in an ellipse has
// no closed form. A polygon with N < 2 has no edge with a length. For integer T it is the
// perimeter of the exact polygon, like Area, not of the rounded vertices.
func (rp RegularPolygon[T]) Perimeter() float64 {
	if rp.N < 2 {
		return 0
	}

	n, size := float64(rp.N), rp.Size.Float()
	if size.Width == size.Height {
		return 2 * n * size.Width * math.Sin(Pi/n)
	}

	perimeter, previous := 0.0, VectorFromAngleSize(rp.vertexAngle(0), size)
	for i := 1; i <= rp.N; i++ {
		offset := VectorFromAngleSize(rp.vertexAngle(i), size)
		perimeter, previous = perimeter+offset.Subtract(previous).Length(), offset
	}

	return perimeter
}

// Inertia returns the polar second moment of area about the center, the rotational inertia
// of the enclosed area at unit density, in closed form: the moment of the polygon inscribed
// in the unit circle, n · sin(2π/n) · (2 + cos(2π/n)) / 12, scaled by w · h · (w² + h²) / 2
// for the semi-axes Size holds, since a regular polygon has the same moment about every
// axis through its center and each axis scales by the cube of one semi-axis and the first
// power of the other. The turn and the orientation do not change it, and a polygon with
// N < 3 encloses no area. For integer T it is the moment of the exact polygon, like Area,
// which Polygon().Inertia() measures from the rounded vertices.
func (rp RegularPolygon[T]) Inertia() float64 {
	if rp.N < 3 {
		return 0
	}

	n, central := float64(rp.N), rp.centralAngle()
	w, h := rp.Size.Float().XY()

	return n * math.Sin(central) * (2 + math.Cos(central)) / 24 * w * h * (float64(w*w) + float64(h*h))
}

// Bounds returns the axis-aligned bounding box of the vertices without building them, or the
// zero box for a polygon without vertices: the box on the corners minMax finds.
func (rp RegularPolygon[T]) Bounds() Box[T] {
	a, b := rp.minMax()

	return Box[T]{a, b}
}

// centralAngle returns the angle between consecutive vertices.
func (rp RegularPolygon[T]) centralAngle() float64 {
	return 2 * Pi / float64(rp.N)
}

// vertexAngle returns the parameter of the vertex at the given index on the ellipse before the
// turn, that many central angles on from Phase: the one expression every vertex is placed
// from, so a measure taken without building the vertices reads the same angles Vertices does,
// and edgeIndex inverts it.
func (rp RegularPolygon[T]) vertexAngle(i int) float64 {
	return rp.Phase + float64(float64(i)*rp.centralAngle())
}

// worldPoint returns the point at the given offset from the center on the ellipse before the
// turn, turned by Angle, as Rectangle.worldPoint places a corner: the offset is rotated once
// and the sum rounded once for an integer T.
func (rp RegularPolygon[T]) worldPoint(offset Vector[float64]) Point[T] {
	return rp.Center.Float().Add(offset.Rotate(rp.Angle)).Cast[T]()
}

// vertex returns the vertex at the given index, the point Vertices places there: the point that
// many steps around the ellipse of the semi-axes, turned by Angle about the center. Equal
// semi-axes take the turn into the angle of the single displacement, since turning a circle and
// stepping around it are the same thing, which keeps an integer vertex where one rounding puts
// it; an ellipse is placed and then turned, and rounded once after.
func (rp RegularPolygon[T]) vertex(i int) Point[T] {
	if rp.Size.Width == rp.Size.Height {
		return rp.Center.Add(VectorFromAngleSize(rp.Angle+rp.vertexAngle(i), rp.Size))
	}

	return rp.worldPoint(VectorFromAngleSize(rp.vertexAngle(i), rp.Size.Float()))
}

// edge returns the edge at the given index, the one Edges yields there: from that vertex to the
// next, the last one closing back to the first.
func (rp RegularPolygon[T]) edge(i int) Segment[T] {
	return Segment[T]{rp.vertex(i), rp.vertex((i + 1) % rp.N)}
}

// edgeIndex returns the index of the edge the ray from the center leaves the polygon through in
// the given direction of the frame before the turn: the direction is taken onto the parameter
// of the ellipse the vertices lie on, where they are equally spaced from Phase, and the edge is
// the one between the two parameters around it.
func (rp RegularPolygon[T]) edgeIndex(direction float64) int {
	parameter := NormalizeAngle(rp.Ellipse().parametricAngle(direction) - rp.Phase)

	return int(parameter/rp.centralAngle()) % rp.N
}

// minMax returns the minimum and maximum corner of the vertices, the corners of Bounds: the pair
// the intersection tests reject shapes by before examining any edge. It ranges the vertices as
// Vertices places them, so the corners are exactly those of Polygon().Bounds(), rounded alike
// for an integer T. An empty polygon returns two zero points.
func (rp RegularPolygon[T]) minMax() (Point[T], Point[T]) {
	if rp.IsEmpty() {
		return Point[T]{}, Point[T]{}
	}

	a := rp.vertex(0)
	b := a
	for i := 1; i < rp.N; i++ {
		vertex := rp.vertex(i)
		a = Point[T]{min(a.X, vertex.X), min(a.Y, vertex.Y)}
		b = Point[T]{max(b.X, vertex.X), max(b.Y, vertex.Y)}
	}

	return a, b
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(vector), rp.Size, rp.N, rp.Angle, rp.Phase}
}

// MoveTo creates a new RegularPolygon with center at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N, rp.Angle, rp.Phase}
}

// Scale creates a new RegularPolygon with size scaled by the given factor. A negative factor
// scales by its absolute value, since Size holds semi-axes and negating both would place every
// vertex half a turn away rather than shrink the polygon; use Rotate for the half turn.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(factor).Abs(), rp.N, rp.Angle, rp.Phase}
}

// ScaleXY creates a new RegularPolygon with size scaled by the given factors, negative ones by
// their absolute value like Scale.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(factorX, factorY).Abs(), rp.N, rp.Angle, rp.Phase}
}

// Unscale creates a new RegularPolygon with size scaled by the inverse factor, the inverse of
// Scale and negative factors taken absolute like it. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) Unscale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Unscale(factor).Abs(), rp.N, rp.Angle, rp.Phase}
}

// UnscaleXY creates a new RegularPolygon with size scaled by the inverse of the given factors,
// the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) UnscaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.UnscaleXY(factorX, factorY).Abs(), rp.N, rp.Angle, rp.Phase}
}

// Resize creates a new RegularPolygon with the given semi-axes, taken absolute like RegPol.
func (rp RegularPolygon[T]) Resize(size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, size.Abs(), rp.N, rp.Angle, rp.Phase}
}

// Canonical creates a new RegularPolygon in the form RegPol and Rotate build, with the size taken
// absolute and the angle and the phase normalized to [0, 2π): a well-formed polygon is returned as
// it is, up to the full turns Equal already ignores. It takes absolute a negative semi-axis written
// as a struct literal or decoded from JSON, which mirrors the vertices, or turns them half a turn
// where both are negative, and brings a decoded angle onto the seam Rotate keeps. An angle within
// Delta of zero or of a full turn, the residue a chain of Rotate, Lerp and Transform calls can
// leave, becomes exactly zero, and so does a phase; Rotate itself never snaps. The phase is not
// reduced to a step around the polygon, which would renumber the vertices.
func (rp RegularPolygon[T]) Canonical() RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Abs(), rp.N, snapAngle(rp.Angle), snapAngle(rp.Phase)}
}

// Grow creates a new RegularPolygon with both semi-axes increased by amount, clamped to zero.
// The amount is added to each semi-axis, so each extent of Bounds grows by about twice it.
func (rp RegularPolygon[T]) Grow(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Grow(amount).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// GrowXY creates a new RegularPolygon with the semi-axes increased by the given amounts along
// its own axes, clamped to zero. Each amount is added to that semi-axis, like Grow.
func (rp RegularPolygon[T]) GrowXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.GrowXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// Shrink creates a new RegularPolygon with both semi-axes decreased by amount, clamped to zero.
func (rp RegularPolygon[T]) Shrink(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Shrink(amount).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// ShrinkXY creates a new RegularPolygon with the semi-axes decreased by the given amounts along
// its own axes, clamped to zero.
func (rp RegularPolygon[T]) ShrinkXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ShrinkXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// Lerp creates a new RegularPolygon in linear interpolation towards the given polygon, moving
// the center and the size together and turning the angle and the phase along the shorter arc
// with LerpAngle, so a tween across the zero angle takes the short way round; both are
// normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like Point.Lerp, with the
// size taken absolute like RegPol. The vertex count does not interpolate: polygons with a
// different N have no shape between them, so Lerp panics for them rather than hiding the
// mistake in an empty polygon.
func (rp RegularPolygon[T]) Lerp(polygon RegularPolygon[T], t float64) RegularPolygon[T] {
	if rp.N != polygon.N {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", rp.N, polygon.N))
	}

	return RegularPolygon[T]{
		rp.Center.Lerp(polygon.Center, t),
		rp.Size.Lerp(polygon.Size, t).Abs(),
		rp.N,
		NormalizeAngle(LerpAngle(rp.Angle, polygon.Angle, t)),
		NormalizeAngle(LerpAngle(rp.Phase, polygon.Phase, t)),
	}
}

// Transform creates a new RegularPolygon by applying the given matrix, exactly for every
// matrix: the center moves, and the turn, the semi-axes and the phase are those of the matrix
// applied to the stretched ring, taken apart again by its singular value decomposition, so a
// shear or a scale of the axes of a turned polygon gives the polygon the vertices land on
// rather than the nearest one. A reflection reverses the winding and mirrors the phase. A move,
// a turn and a scale along the axes of an unturned polygon keep the phase; the decomposition is
// not unique, and where it reads the semi-axes the other way round, the phase takes the
// quarter turn. For integer T the center and the semi-axes are each rounded once.
func (rp RegularPolygon[T]) Transform[M Float](matrix Matrix[M]) RegularPolygon[T] {
	after, stretch, before, reflected := matrix.mapping(rp.Angle, rp.Size.Float()).decomposition()

	phase := rp.Phase + before
	if reflected {
		phase = -phase
	}

	return RegularPolygon[T]{rp.Center.Transform(matrix), stretch.Size().Cast[T](), rp.N, NormalizeAngle(after), NormalizeAngle(phase)}
}

// Rotate creates a new RegularPolygon turned by the given angle (in radians) about its center,
// in the same sense as Vector.Rotate.
// The stored angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size, rp.N, NormalizeAngle(rp.Angle + angle), rp.Phase}
}

// AlignTo creates a new RegularPolygon moved so that its Anchor in the given direction lands
// on the point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (rp RegularPolygon[T]) AlignTo(direction Direction, point Point[T]) RegularPolygon[T] {
	return rp.Translate(point.Subtract(rp.Anchor(direction)))
}

// Contains reports whether the given point lies within the polygon, boundary included within
// the tolerance, the same closed convention as Polygon.Contains, and exactly what the Polygon
// of the vertices contains, for an integer T on the rounded vertices. A point outside Bounds
// is rejected before any edge is examined; an empty polygon contains nothing.
func (rp RegularPolygon[T]) Contains(point Point[T]) bool {
	if rp.IsEmpty() {
		return false
	}

	a, b := rp.minMax()

	return rp.containsWithin(point, a, b)
}

// DistanceTo returns the distance from the given point to the nearest point of the polygon:
// zero for a point within it, the same closed convention as Contains, and otherwise the
// distance to the nearest edge. An empty polygon is infinitely far from every point.
func (rp RegularPolygon[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(rp.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, in one pass over the edges Edges iterates without building the vertices, the
// edgeWalk every closed shape makes: zero for a point on an edge within the tolerance or
// inside by the even-odd rule, the squared distance to the nearest edge otherwise, and
// infinity for an empty polygon. It is a float64 even for an integer T, like
// Polygon.DistanceSquaredTo. Contains and IntersectsCircle are built on it.
func (rp RegularPolygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	return rp.walk(point).result()
}

// Nearest returns the point of the polygon nearest to the given point: the point itself
// exactly where Contains holds, and otherwise the foot on the nearest edge, read off the same
// walk DistanceSquaredTo makes. For integer T the foot is rounded once and can land off the
// boundary, where Contains rejects it. An empty polygon returns the zero point.
func (rp RegularPolygon[T]) Nearest(point Point[T]) Point[T] {
	return rp.walk(point).nearest(point)
}

// EnclosesCircle reports whether the circle lies within the regular polygon: its center is
// contained and every edge is at least the radius away, within the tolerance, read off the walk
// DistanceSquaredTo makes, so a circle touching an edge from inside is enclosed.
func (rp RegularPolygon[T]) EnclosesCircle(circle Circle[T]) bool {
	return rp.walk(circle.Center).clears(circle)
}

// EnclosesSegment reports whether the segment lies within the regular polygon: both endpoints are
// contained, within the tolerance, and a convex shape holds every point between two it contains.
func (rp RegularPolygon[T]) EnclosesSegment(segment Segment[T]) bool {
	a, b := rp.minMax()

	return rp.containsWithin(segment.Start, a, b) && rp.containsWithin(segment.End, a, b)
}

// EnclosesPolygon reports whether the polygon lies within the regular polygon: every vertex is
// contained, within the tolerance, and a convex shape holds every point between points it
// contains. An empty polygon is enclosed by nothing.
func (rp RegularPolygon[T]) EnclosesPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := rp.minMax()
	for vertex := range polygon.Vertices() {
		if !rp.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// EnclosesRectangle reports whether the rectangle lies within the regular polygon: every corner is
// contained, within the tolerance, whatever the angles.
func (rp RegularPolygon[T]) EnclosesRectangle(rectangle Rectangle[T]) bool {
	a, b := rp.minMax()
	for vertex := range rectangle.Vertices() {
		if !rp.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// EnclosesRegularPolygon reports whether the regular polygon lies within the regular polygon: every
// vertex is contained, within the tolerance. An empty polygon is enclosed by nothing.
func (rp RegularPolygon[T]) EnclosesRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := rp.minMax()
	for vertex := range polygon.Vertices() {
		if !rp.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// EnclosesBox reports whether the box lies within the regular polygon: every corner of the box is
// contained, within the tolerance.
func (rp RegularPolygon[T]) EnclosesBox(box Box[T]) bool {
	a, b := rp.minMax()
	for _, corner := range box.corners() {
		if !rp.containsWithin(corner, a, b) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the polygon and the circle share a point, as
// Circle.IntersectsRegularPolygon does.
func (rp RegularPolygon[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRegularPolygon(rp)
}

// IntersectsSegment reports whether the polygon and the segment share a point, as
// Segment.IntersectsRegularPolygon does.
func (rp RegularPolygon[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsRegularPolygon(rp)
}

// IntersectionSegment returns the points where the segment crosses the polygon boundary, as
// Segment.IntersectionRegularPolygon does.
func (rp RegularPolygon[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionRegularPolygon(rp)
}

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns the
// extended slice, as Segment.AppendIntersectionRegularPolygon does.
func (rp RegularPolygon[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	return segment.AppendIntersectionRegularPolygon(dst, rp)
}

// IntersectsRay reports whether the regular polygon and the ray share a point, as
// Ray.IntersectsRegularPolygon does.
func (rp RegularPolygon[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsRegularPolygon(rp)
}

// IntersectionRay returns the points where the ray crosses the regular polygon boundary, as
// Ray.IntersectionRegularPolygon does.
func (rp RegularPolygon[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionRegularPolygon(rp)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice, as Ray.AppendIntersectionRegularPolygon does.
func (rp RegularPolygon[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	return ray.AppendIntersectionRegularPolygon(dst, rp)
}

// IntersectsPolygon reports whether the regular polygon and the polygon share a point, as
// Polygon.IntersectsRegularPolygon does.
func (rp RegularPolygon[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsRegularPolygon(rp)
}

// IntersectsRectangle reports whether the polygon and the rectangle share a point, as
// Rectangle.IntersectsRegularPolygon does.
func (rp RegularPolygon[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.IntersectsRegularPolygon(rp)
}

// IntersectsRegularPolygon reports whether the polygons share a point: a vertex of one lies within
// the other, or an edge of one crosses an edge of the other, as Polygon.IntersectsPolygon decides
// on their Polygon forms, without building them. Touching polygons intersect, within the tolerance,
// and an empty polygon intersects nothing. Polygons whose Bounds do not overlap are rejected before
// any edge pair is examined. The vertices of the given polygon are placed again for every edge of
// this one, since no vertex slice is built: the cost is a sine and cosine per edge pair, not an
// allocation.
func (rp RegularPolygon[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if rp.IsEmpty() || polygon.IsEmpty() {
		return false
	}

	a1, b1 := rp.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(rp.vertex(0), a2, b2) || rp.containsWithin(polygon.vertex(0), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range rp.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range polygon.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectsBox reports whether the polygon and the box share a point, as
// Rectangle.IntersectsRegularPolygon decides it against a rectangle that is not rotated: a
// vertex of one lies within the other, or an edge of the box, between its own corners, crosses
// an edge of the polygon. The two Bounds reject the pair before any edge is examined, and an
// empty polygon intersects nothing. The polygon places its vertices again for every edge of the
// box, since no vertex slice is built.
func (rp RegularPolygon[T]) IntersectsBox(box Box[T]) bool {
	if rp.IsEmpty() {
		return false
	}

	a, b := rp.minMax()

	if !overlaps(box.Min, box.Max, a, b) {
		return false
	}

	if rp.containsWithin(box.Min, a, b) || box.Contains(rp.vertex(0)) {
		return true
	}

	probe := edgeProbe[T]{a: a, b: b}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if !probe.aim(edge) {
			continue
		}

		for other := range rp.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// containsWithin is Contains for a caller that already holds the corners of Bounds, so the
// intersection tests place the box once and reuse it for every point they test.
func (rp RegularPolygon[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && rp.DistanceSquaredTo(point) == 0
}

// walk folds every edge into the edgeWalk DistanceSquaredTo, Nearest and EnclosesCircle read,
// stopping at an edge the point lies on within the tolerance.
func (rp RegularPolygon[T]) walk(point Point[T]) edgeWalk[T] {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range rp.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	return w
}

// Equal checks if center point, size, number of vertices, angle and phase are equal. Angles
// are compared with EqualAngle, so a full turn or the sign of an angle does not matter. The
// values are compared, not the shapes: a phase a step around the polygon apart, or a turn
// moved into the phase of equal semi-axes, draws the same polygon with its vertices numbered
// otherwise, and is a different value.
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool {
	return rp.Center.Equal(polygon.Center) && rp.Size.Equal(polygon.Size) && rp.N == polygon.N && EqualAngle(rp.Angle, polygon.Angle) && EqualAngle(rp.Phase, polygon.Phase)
}

// IsZero checks if center point, size, number of vertices, angle and phase are zero, comparing
// the angles like Equal so that a full turn counts as zero.
func (rp RegularPolygon[T]) IsZero() bool {
	return rp.Equal(RegularPolygon[T]{})
}

// IsEmpty checks if the polygon has no vertices.
func (rp RegularPolygon[T]) IsEmpty() bool {
	return rp.N < 1
}

// IsAligned reports whether the polygon is not turned: its Angle is exactly zero, as RegPol with no
// angle and Rotate by a full turn leave it. No tolerance is applied; Canonical snaps a residue to
// zero. The phase places the vertices within the frame and does not turn it, so an aligned polygon
// may have any phase.
func (rp RegularPolygon[T]) IsAligned() bool {
	return rp.Angle == 0
}

// Polygon converts the regular polygon into a generic Polygon with the vertices Vertices
// iterates, in one allocation. A polygon with N < 1 has nil vertices, so its Polygon is zero
// like Pol(nil). A caller reusing a buffer appends them to it with slices.AppendSeq.
func (rp RegularPolygon[T]) Polygon() Polygon[T] {
	if rp.IsEmpty() {
		return Polygon[T]{}
	}

	vertices := make([]Point[T], rp.N)
	for i := range vertices {
		vertices[i] = rp.vertex(i)
	}

	return Polygon[T]{vertices}
}

// Ellipse converts the polygon into the Ellipse its vertices lie on, the one it is inscribed
// in: the same center, semi-axes and angle, so the conversion is exact and Ellipse.RegularPolygon
// is its inverse for the same vertex count and the orientation whose phase the polygon has,
// which the ellipse does not carry. A polygon with N < 1 has no vertices and still names the
// ellipse its Size and Angle describe.
func (rp RegularPolygon[T]) Ellipse() Ellipse[T] {
	return Ellipse[T]{rp.Center, rp.Size, rp.Angle}
}

// Circle converts the polygon into the circle around it, the circle around the Ellipse it is
// inscribed in: the one of the major semi-axis, which passes through the vertices of a polygon
// of equal semi-axes and contains every other. Only the circumscribed circle is offered: it
// names the RegularPolygon of the same center, size and angle, and the inscribed one does not.
func (rp RegularPolygon[T]) Circle() Circle[T] {
	return rp.Ellipse().Circle()
}

// Cast converts the polygon to a RegularPolygon of another number type, rounding as Cast does.
func (rp RegularPolygon[T]) Cast[R Number]() RegularPolygon[R] {
	return RegularPolygon[R]{rp.Center.Cast[R](), rp.Size.Cast[R](), rp.N, rp.Angle, rp.Phase}
}

// Int converts the regular polygon to a RegularPolygon[int].
func (rp RegularPolygon[T]) Int() RegularPolygon[int] {
	return RegularPolygon[int]{rp.Center.Int(), rp.Size.Int(), rp.N, rp.Angle, rp.Phase}
}

// Float converts the regular polygon to a RegularPolygon[float64].
func (rp RegularPolygon[T]) Float() RegularPolygon[float64] {
	return RegularPolygon[float64]{rp.Center.Float(), rp.Size.Float(), rp.N, rp.Angle, rp.Phase}
}

// String returns the polygon in the form of its constructor: RegPol((x,y);WxH;n), center, size
// and vertex count, with the angle appended as RegPol((x,y);WxH;n;a) for a rotated polygon and
// the phase after it as RegPol((x,y);WxH;n;a;p) where it is not zero, as the JSON carries each
// only then.
func (rp RegularPolygon[T]) String() string {
	switch {
	case rp.Phase != 0:
		return fmt.Sprintf("RegPol(%s;%s;%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N), String(rp.Angle), String(rp.Phase))
	case !rp.IsAligned():
		return fmt.Sprintf("RegPol(%s;%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N), String(rp.Angle))
	default:
		return fmt.Sprintf("RegPol(%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N))
	}
}

// RegularPolygonOrientationPhase returns the phase RegularPolygonWithOrientation places the first
// vertex of a regular polygon of n vertices at, normalized to [0, 2π) like Rotate.
// OrientationPointyTop puts the first vertex at the top (-Y, 3π/2); OrientationFlatTop puts the
// midpoint of an edge there, so the first vertex sits half a step before it at 3π/2 - π/n. The
// phase places the vertices on the ellipse before its semi-axes stretch it, so the polygon is
// stretched across Width and Height as they stand for either orientation and every n. A polygon
// with n < 1 has no edge to place, so both orientations give the top phase rather than dividing
// by n. An orientation other than OrientationFlatTop and OrientationPointyTop has no meaning and
// panics.
func RegularPolygonOrientationPhase(n int, orientation Orientation) float64 {
	top := 3 * Pi / 2

	switch orientation {
	case OrientationFlatTop:
		if n < 1 {
			return top
		}

		return NormalizeAngle(top - Pi/float64(n))
	case OrientationPointyTop:
		return top
	default:
		panic(fmt.Sprintf("geom: unknown orientation %d", orientation))
	}
}

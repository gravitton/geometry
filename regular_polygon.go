package geom

import (
	"fmt"
	"iter"
	"math"
)

// RegularPolygon is a polygon of N vertices equally spaced around a center.
//
// Radius holds the semi-axes of the Ellipse the vertices lie on, not the full extent a Rectangle
// holds; Bounds gives the extent. Phase places the first vertex on that ellipse and
// each next a step of 2π/n on, and Angle then turns the whole ring about the center. The two differ
// only for unequal semi-axes, where the phase decides how the vertices sit against the stretch.
// Every affine image of a regular polygon is one of these, so Transform keeps any matrix exactly.
//
// The radius is never negative: constructors and changes take it absolute. A negative one written as
// a struct literal or decoded from JSON mirrors the vertices until Canonical repairs it.
//
// For integer T the vertices are rounded, and the rounded outline need not be convex. Contains
// follows it exactly, but the Contains<Kind> methods test the points of the other shape alone, so
// they can admit a shape passing through a dent the rounding made.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius Size[T]  `json:",embed"`
	N      int      `json:"n"`
	Angle  float64  `json:"a,omitzero"`
	Phase  float64  `json:"p,omitzero"`
}

// RegPol is shorthand for RegularPolygon{center, radius, n, angle, phase}, with the radius taken
// absolute.
func RegPol[T Number](center Point[T], radius Size[T], n int, angle, phase float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, radius.Abs(), n, angle, phase}
}

// RegularPolygonWithOrientation returns the regular polygon with no turn and its first vertex
// placed by RegularPolygonOrientationPhase, the radius taken absolute: Width is the semi-axis across
// and Height the one up, for every n. It panics for an orientation that is neither
// OrientationFlatTop nor OrientationPointyTop.
func RegularPolygonWithOrientation[T Number](center Point[T], radius Size[T], n int, orientation Orientation) RegularPolygon[T] {
	return RegPol(center, radius, n, 0, RegularPolygonOrientationPhase(n, orientation))
}

// Anchor returns the point where the ray from the center in the direction leaves the polygon, or
// the center for DirectionNone. The direction is taken after the turn, so on a flat-top polygon Top
// is the midpoint of the top edge, on a pointy-top one the top vertex. For integer T the point is
// rounded once more and can lie off the edge between the rounded vertices. A polygon of fewer than
// three vertices, or one flattened along the ray, anchors at its center.
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

// Vertices returns an iterator over the vertices from Phase by increasing angle, clockwise as drawn
// on a screen with Y pointing down, and nothing for an empty polygon.
func (rp RegularPolygon[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for i := range rp.N {
			if !yield(rp.vertex(i)) {
				return
			}
		}
	}
}

// Edges returns an iterator over the edges in the order of Vertices, each from a vertex to the next
// and the last closing back to the first: none for an empty polygon, and one of zero length for a
// single vertex.
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

// Centroid returns the center.
func (rp RegularPolygon[T]) Centroid() Point[T] {
	return rp.Center
}

// Area returns the area (n/2·w·h·sin(2π/n)), zero for fewer than three vertices. For integer T it
// is the area of the exact polygon, not of the rounded vertices Polygon().Area() measures.
func (rp RegularPolygon[T]) Area() float64 {
	if rp.N < 3 {
		return 0
	}

	n := float64(rp.N)

	return n / 2 * float64(rp.Radius.Width) * float64(rp.Radius.Height) * math.Sin(2*Pi/n)
}

// Perimeter returns the total length of the edges, zero for fewer than two vertices. For integer T
// it is the perimeter of the exact polygon, not of the rounded vertices.
func (rp RegularPolygon[T]) Perimeter() float64 {
	if rp.N < 2 {
		return 0
	}

	n, size := float64(rp.N), rp.Radius.Float()
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

// Inertia returns the polar second moment of area about the center
// (n·sin(2π/n)·(2+cos(2π/n))·w·h·(w²+h²)/24), the rotational inertia of the polygon at unit
// density, zero for fewer than three vertices. For integer T it is the moment of the exact polygon,
// not of the rounded vertices.
func (rp RegularPolygon[T]) Inertia() float64 {
	if rp.N < 3 {
		return 0
	}

	n, central := float64(rp.N), rp.centralAngle()
	w, h := rp.Radius.Float().XY()

	return n * math.Sin(central) * (2 + math.Cos(central)) / 24 * w * h * (float64(w*w) + float64(h*h))
}

// Bounds returns the axis-aligned bounding box, the zero box for an empty polygon.
func (rp RegularPolygon[T]) Bounds() Box[T] {
	a, b := rp.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle about the center through the farthest vertex, the zero circle
// for an empty polygon. It is the smallest enclosing circle only for equal semi-axes and at least
// two vertices. For integer T the radius is rounded up, so it still contains the rounded vertices.
func (rp RegularPolygon[T]) BoundingCircle() Circle[T] {
	if rp.IsEmpty() {
		return Circle[T]{}
	}

	center := rp.Center.Float()

	var reach float64
	for i := range rp.N {
		reach = max(reach, center.DistanceSquaredTo(rp.vertex(i).Float()))
	}

	return Circle[T]{rp.Center, castUp[T](math.Sqrt(reach))}
}

// centralAngle returns the angle between consecutive vertices.
func (rp RegularPolygon[T]) centralAngle() float64 {
	return 2 * Pi / float64(rp.N)
}

// vertexAngle returns the parameter on the ellipse of the vertex at the index before the turn, the
// one expression every vertex and measure is placed from; edgeIndex inverts it.
func (rp RegularPolygon[T]) vertexAngle(i int) float64 {
	return rp.Phase + float64(float64(i)*rp.centralAngle())
}

// worldPoint returns the point at the offset from the center before the turn, rotated by Angle and
// rounded once.
func (rp RegularPolygon[T]) worldPoint(offset Vector[float64]) Point[T] {
	return rp.Center.Float().Add(offset.Rotate(rp.Angle)).Cast[T]()
}

// vertex returns the vertex at the index. Equal semi-axes fold the turn into the angle of one
// displacement, which keeps an integer vertex where one rounding puts it; unequal ones are placed,
// turned, then rounded.
func (rp RegularPolygon[T]) vertex(i int) Point[T] {
	if rp.Radius.Width == rp.Radius.Height {
		return rp.Center.Add(VectorFromAngleSize(rp.Angle+rp.vertexAngle(i), rp.Radius))
	}

	return rp.worldPoint(VectorFromAngleSize(rp.vertexAngle(i), rp.Radius.Float()))
}

// edge returns the edge at the index, as Edges yields it.
func (rp RegularPolygon[T]) edge(i int) Segment[T] {
	return Segment[T]{rp.vertex(i), rp.vertex((i + 1) % rp.N)}
}

// edgeIndex returns the index of the edge through which the ray in the direction, taken before the
// turn, leaves the polygon, found from the parameter of that direction on the ellipse.
func (rp RegularPolygon[T]) edgeIndex(direction float64) int {
	parameter := NormalizeAngle(rp.Ellipse().parametricAngle(direction) - rp.Phase)

	return int(parameter/rp.centralAngle()) % rp.N
}

// minMax returns the corners of Bounds from the vertices as Vertices places them, two zero points
// for an empty polygon.
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

// Translate returns the polygon moved by the vector.
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(vector), rp.Radius, rp.N, rp.Angle, rp.Phase}
}

// MoveTo returns the polygon with its center at the point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Radius, rp.N, rp.Angle, rp.Phase}
}

// AlignTo returns the polygon moved so that its Anchor in the direction lies at the point.
func (rp RegularPolygon[T]) AlignTo(direction Direction, point Point[T]) RegularPolygon[T] {
	return rp.Translate(point.Subtract(rp.Anchor(direction)))
}

// Rotate returns the polygon turned by the angle in radians about its center, the stored angle
// normalized to [0, 2π).
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius, rp.N, NormalizeAngle(rp.Angle + angle), rp.Phase}
}

// Scale returns the polygon with its semi-axes scaled by the factor, taken absolute. Rotate gives
// the half turn a negative factor would.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.Scale(factor).Abs(), rp.N, rp.Angle, rp.Phase}
}

// ScaleXY returns the polygon with its semi-axes scaled by the factors, taken absolute.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.ScaleXY(factorX, factorY).Abs(), rp.N, rp.Angle, rp.Phase}
}

// Unscale returns the polygon with its semi-axes divided by the factor, taken absolute. It panics
// for a zero factor.
func (rp RegularPolygon[T]) Unscale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.Unscale(factor).Abs(), rp.N, rp.Angle, rp.Phase}
}

// UnscaleXY returns the polygon with its semi-axes divided by the factors, taken absolute. It
// panics for a zero factor.
func (rp RegularPolygon[T]) UnscaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.UnscaleXY(factorX, factorY).Abs(), rp.N, rp.Angle, rp.Phase}
}

// Resize returns the polygon with the semi-axes, taken absolute.
func (rp RegularPolygon[T]) Resize(radius Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, radius.Abs(), rp.N, rp.Angle, rp.Phase}
}

// Transform returns the polygon mapped exactly by the matrix, a shear included: the center moves,
// and the angle, semi-axes and phase are taken apart from the mapped ring again. A reflection
// mirrors the phase, so the vertices still wind by increasing angle and, after the first, meet the
// mapped ones in reverse order. Where the decomposition reads the semi-axes the other way round,
// the phase takes a quarter turn.
func (rp RegularPolygon[T]) Transform[M Float](matrix Matrix[M]) RegularPolygon[T] {
	after, stretch, before, reflected := matrix.mapping(rp.Angle, rp.Radius.Float()).decomposition()

	phase := rp.Phase + before
	if reflected {
		phase = -phase
	}

	return RegularPolygon[T]{rp.Center.Transform(matrix), stretch.Size().Cast[T](), rp.N, NormalizeAngle(after), NormalizeAngle(phase)}
}

// Grow returns the polygon with both semi-axes increased by the amount, clamped at zero, so each
// extent of Bounds grows by about twice it.
func (rp RegularPolygon[T]) Grow(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.Grow(amount).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// GrowXY returns the polygon with its semi-axes increased by the amounts, clamped at zero.
func (rp RegularPolygon[T]) GrowXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.GrowXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// Shrink returns the polygon with both semi-axes decreased by the amount, clamped at zero.
func (rp RegularPolygon[T]) Shrink(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.Shrink(amount).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// ShrinkXY returns the polygon with its semi-axes decreased by the amounts, clamped at zero.
func (rp RegularPolygon[T]) ShrinkXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.ShrinkXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle, rp.Phase}
}

// Lerp returns the polygon interpolated toward the given one, center and radius together and the
// angle and phase along the shorter arc, and extrapolated outside [0, 1]. The radius is taken
// absolute and the angles normalized to [0, 2π). It panics for polygons of different N.
func (rp RegularPolygon[T]) Lerp(polygon RegularPolygon[T], t float64) RegularPolygon[T] {
	if rp.N != polygon.N {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", rp.N, polygon.N))
	}

	return RegularPolygon[T]{
		rp.Center.Lerp(polygon.Center, t),
		rp.Radius.Lerp(polygon.Radius, t).Abs(),
		rp.N,
		NormalizeAngle(LerpAngle(rp.Angle, polygon.Angle, t)),
		NormalizeAngle(LerpAngle(rp.Phase, polygon.Phase, t)),
	}
}

// Canonical returns the polygon with its radius taken absolute and its angle and phase normalized to
// [0, 2π), the receiver when they already are. An angle or phase within Delta of zero or a full
// turn becomes exactly zero; Rotate never snaps. The phase is not reduced to one step, which would
// renumber the vertices.
func (rp RegularPolygon[T]) Canonical() RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Radius.Abs(), rp.N, snapAngle(rp.Angle), snapAngle(rp.Phase)}
}

// DistanceTo returns the distance from the point to the polygon, zero inside, and infinity for an
// empty polygon.
func (rp RegularPolygon[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(rp.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (rp RegularPolygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	return rp.walk(point).result()
}

// Nearest returns the point of the polygon nearest to the given point, the point itself inside, and
// the zero point for an empty polygon.
func (rp RegularPolygon[T]) Nearest(point Point[T]) Point[T] {
	return rp.walk(point).nearest(point)
}

// Support returns the vertex farthest along the direction, the first in the order of Vertices where
// two reach alike, and the center for an empty polygon.
func (rp RegularPolygon[T]) Support(direction Vector[float64]) Point[float64] {
	support, reach := rp.Center.Float(), math.Inf(-1)
	for vertex := range rp.Vertices() {
		point := vertex.Float()
		if along := point.Vector().Dot(direction); along > reach {
			support, reach = point, along
		}
	}

	return support
}

// Contains reports whether the point lies in the polygon.
func (rp RegularPolygon[T]) Contains(point Point[T]) bool {
	if rp.IsEmpty() {
		return false
	}

	a, b := rp.minMax()

	return rp.containsWithin(point, a, b)
}

// ContainsCircle reports whether the circle lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsCircle(circle Circle[T]) bool {
	return rp.walk(circle.Center).clears(circle)
}

// ContainsCapsule reports whether the capsule lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return rp.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && rp.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the sector lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsSector(sector Sector[T]) bool {
	outer, inner := rp.Float(), sector.Float()
	for _, corner := range inner.corners() {
		if !outer.Contains(corner) {
			return false
		}
	}

	for edge := range outer.Edges() {
		if normal := edge.Normal(); !outer.Contains(inner.Support(normal)) || !outer.Contains(inner.Support(normal.Negate())) {
			return false
		}
	}

	return true
}

// ContainsSegment reports whether the segment lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsSegment(segment Segment[T]) bool {
	a, b := rp.minMax()

	return rp.containsWithin(segment.Start, a, b) && rp.containsWithin(segment.End, a, b)
}

// ContainsPolygon reports whether the polygon lies entirely in the regular polygon.
func (rp RegularPolygon[T]) ContainsPolygon(polygon Polygon[T]) bool {
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

// ContainsRectangle reports whether the rectangle lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	a, b := rp.minMax()
	for vertex := range rectangle.Vertices() {
		if !rp.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the given polygon lies entirely in this one.
func (rp RegularPolygon[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
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

// ContainsBox reports whether the box lies entirely in the polygon.
func (rp RegularPolygon[T]) ContainsBox(box Box[T]) bool {
	a, b := rp.minMax()
	for _, corner := range box.corners() {
		if !rp.containsWithin(corner, a, b) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the polygon and the circle share a point.
func (rp RegularPolygon[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRegularPolygon(rp)
}

// IntersectsCapsule reports whether the polygon and the capsule share a point.
func (rp RegularPolygon[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsRegularPolygon(rp)
}

// IntersectsSector reports whether the polygon and the sector share a point.
func (rp RegularPolygon[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsRegularPolygon(rp)
}

// IntersectsSegment reports whether the polygon and the segment share a point.
func (rp RegularPolygon[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsRegularPolygon(rp)
}

// IntersectsRay reports whether the polygon and the ray share a point.
func (rp RegularPolygon[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsRegularPolygon(rp)
}

// IntersectsLine reports whether the polygon and the line share a point.
func (rp RegularPolygon[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsRegularPolygon(rp)
}

// IntersectsPolygon reports whether the regular polygon and the polygon share a point.
func (rp RegularPolygon[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsRegularPolygon(rp)
}

// IntersectsRectangle reports whether the polygon and the rectangle share a point.
func (rp RegularPolygon[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.IntersectsRegularPolygon(rp)
}

// IntersectsRegularPolygon reports whether the polygons share a point.
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

// IntersectsBox reports whether the polygon and the box share a point.
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

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (rp RegularPolygon[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	return segment.IntersectionRegularPolygon(rp)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (rp RegularPolygon[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return ray.IntersectionRegularPolygon(rp)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (rp RegularPolygon[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return line.IntersectionRegularPolygon(rp)
}

// containsWithin is Contains for a caller that already holds the corners of Bounds.
func (rp RegularPolygon[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && rp.DistanceSquaredTo(point) == 0
}

// walk folds every edge into the edgeWalk DistanceSquaredTo, Nearest and ContainsCircle read,
// stopping at an edge the point lies on.
func (rp RegularPolygon[T]) walk(point Point[T]) edgeWalk[T] {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range rp.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	return w
}

// Equal reports whether the polygons have equal center, radius, N, angle and phase, the angles
// compared with EqualAngle. It compares values, not outlines: a phase a step apart draws the same
// polygon numbered otherwise and is not equal.
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool {
	return rp.Center.Equal(polygon.Center) && rp.Radius.Equal(polygon.Radius) && rp.N == polygon.N && EqualAngle(rp.Angle, polygon.Angle) && EqualAngle(rp.Phase, polygon.Phase)
}

// IsZero reports whether the polygon is the zero value, an angle or phase of a full turn counting
// as zero.
func (rp RegularPolygon[T]) IsZero() bool {
	return rp.Equal(RegularPolygon[T]{})
}

// IsEmpty reports whether the polygon has no vertices.
func (rp RegularPolygon[T]) IsEmpty() bool {
	return rp.N < 1
}

// IsAligned reports whether the angle is exactly zero; Canonical snaps a residue to zero. The phase
// does not turn the frame, so an aligned polygon may have any phase.
func (rp RegularPolygon[T]) IsAligned() bool {
	return rp.Angle == 0
}

// Circle returns the circle about the center of the major semi-axis, through the vertices where the
// semi-axes are equal and around them otherwise.
func (rp RegularPolygon[T]) Circle() Circle[T] {
	return rp.Ellipse().Circle()
}

// Polygon returns the regular polygon as a polygon of its Vertices, the zero Polygon for an empty
// one.
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

// Ellipse returns the ellipse the vertices lie on, of the same center, semi-axes and angle, even
// for an empty polygon.
func (rp RegularPolygon[T]) Ellipse() Ellipse[T] {
	return Ellipse[T]{rp.Center, rp.Radius, rp.Angle}
}

// Cast converts the polygon to a RegularPolygon of another number type, rounding as Cast does.
func (rp RegularPolygon[T]) Cast[R Number]() RegularPolygon[R] {
	return RegularPolygon[R]{rp.Center.Cast[R](), rp.Radius.Cast[R](), rp.N, rp.Angle, rp.Phase}
}

// Int converts the polygon to a RegularPolygon[int].
func (rp RegularPolygon[T]) Int() RegularPolygon[int] {
	return RegularPolygon[int]{rp.Center.Int(), rp.Radius.Int(), rp.N, rp.Angle, rp.Phase}
}

// Float converts the polygon to a RegularPolygon[float64].
func (rp RegularPolygon[T]) Float() RegularPolygon[float64] {
	return RegularPolygon[float64]{rp.Center.Float(), rp.Radius.Float(), rp.N, rp.Angle, rp.Phase}
}

// String returns the polygon as RegPol((x,y);WxH;n), center, radius and vertex count, with the angle
// appended as RegPol((x,y);WxH;n;a) for a turned polygon or a nonzero phase, and the phase after it
// as RegPol((x,y);WxH;n;a;p) where it is not zero.
func (rp RegularPolygon[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "RegPol("...)
	buf = rp.Center.appendString(buf)
	buf = append(buf, ';')
	buf = rp.Radius.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, rp.N)
	if rp.Phase != 0 || !rp.IsAligned() {
		buf = append(buf, ';')
		buf = appendString(buf, rp.Angle)
	}
	if rp.Phase != 0 {
		buf = append(buf, ';')
		buf = appendString(buf, rp.Phase)
	}
	buf = append(buf, ')')

	return string(buf)
}

// RegularPolygonOrientationPhase returns the phase of the first vertex of a regular polygon of n
// vertices with the orientation, normalized to [0, 2π). OrientationPointyTop puts a vertex at the
// top, OrientationFlatTop the midpoint of an edge, its first vertex half a step before. Fewer than
// one vertex gives the top for both. It panics for an orientation that is neither
// OrientationFlatTop nor OrientationPointyTop.
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

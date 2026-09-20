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
// Angle turns the polygon about its center, the same meaning it has on Rectangle: the vertices
// are placed on the ellipse of the semi-axes, the first at angle zero, and the whole ring is
// then turned. Where the semi-axes are equal, turning and stepping around the ring are the
// same thing and the first vertex lies at Angle.
//
// The size is never negative: RegPol and the orientation constructors take it absolute, and
// Scale takes a negative factor absolute, since a negative semi-axis would place
// every vertex half a turn away rather than describe a different polygon. A negative size can
// only be written as a struct literal or decoded from JSON; Canonical repairs it.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	N      int      `json:"n"`
	Angle  float64  `json:"a,omitzero"`
}

// RegPol is shorthand for RegularPolygon{center, size, n, angle}, with the size taken absolute.
func RegPol[T Number](center Point[T], size Size[T], n int, angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, size.Abs(), n, angle}
}

// RegularPolygonWithOrientation creates a RegularPolygon with the given orientation, with the
// size taken absolute like RegPol.
func RegularPolygonWithOrientation[T Number](center Point[T], size Size[T], n int, orientation Orientation) RegularPolygon[T] {
	return RegPol(center, size, n, RegularPolygonOrientationAngle(n, orientation))
}

// Triangle creates a RegularPolygon with 3 vertices.
func Triangle[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 3, orientation)
}

// Square creates a RegularPolygon with 4 vertices.
func Square[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 4, orientation)
}

// Hexagon creates a RegularPolygon with 6 vertices.
func Hexagon[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 6, orientation)
}

// RegularPolygonOrientationAngle returns the initial vertex angle for a regular polygon with n sides
// and the given orientation, normalized to [0, 2π) like Rotate. OrientationPointyTop puts the first vertex at
// the top (-Y, 3π/2); OrientationFlatTop puts the midpoint of an edge there, so the first vertex sits half a
// step before it at 3π/2 - π/n. A polygon with n < 1 has no edge to place, so both orientations
// give the top angle rather than dividing by n. An orientation other than OrientationFlatTop and OrientationPointyTop
// has no meaning and panics.
func RegularPolygonOrientationAngle(n int, orientation Orientation) float64 {
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

// Vertices iterates the polygon vertices in order starting from Angle, by increasing angle —
// the same winding as Directions and Rectangle.Vertices, and clockwise as drawn on a screen
// with Y pointing down, without allocating; collect them with slices.Collect where a slice is
// needed. A polygon with N < 1 has no vertices and yields nothing.
// For integer T, each vertex component is rounded to the nearest integer, so vertices at
// non-right angles may be off by up to half a unit. Use float64 for exact positions.
func (rp RegularPolygon[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for i := 0; i < rp.N; i++ {
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

// Anchor returns the point of the boundary in the given direction from the center, or the
// center itself for DirectionNone: the point where the ray from the center leaves the polygon,
// as Rectangle.Anchor gives a corner or an edge midpoint and Ellipse.Anchor the point of its
// boundary. The direction is taken in the world, where the polygon stands after its turn,
// unlike Rectangle.Anchor, since the orientation constructors turn the polygon to place its
// top: on a flat-top polygon the Top anchor is the midpoint of the top edge, on a pointy-top
// one its top vertex, and Rotate moves every anchor along the boundary. For integer T
// the edge is the one between the rounded vertices, so the anchor lies on the boundary
// Contains judges, rounded once. A polygon with fewer than three vertices encloses no area and
// anchors at its center, as does one whose semi-axis lies along the ray.
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

	return n * math.Sin(central) * (2 + math.Cos(central)) / 24 * w * h * (w*w + h*h)
}

// Bounds returns the axis-aligned bounding rectangle of the vertices without building them, or
// the zero rectangle for a polygon without vertices, like Polygon.Bounds: the rectangle on the
// corners minMax finds.
func (rp RegularPolygon[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(rp.minMax())
}

// centralAngle returns the angle between consecutive vertices.
func (rp RegularPolygon[T]) centralAngle() float64 {
	return 2 * Pi / float64(rp.N)
}

// vertexAngle returns the angle of the vertex at the given index on the ellipse before the
// turn, that many central angles from the first: the one expression every vertex is placed
// from, so a measure taken without building the vertices reads the same angles Vertices does,
// and nearestIndex inverts it.
func (rp RegularPolygon[T]) vertexAngle(i int) float64 {
	return float64(i) * rp.centralAngle()
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
// of the ellipse the vertices lie on, where they are equally spaced, and the edge is the one
// between the two parameters around it.
func (rp RegularPolygon[T]) edgeIndex(direction float64) int {
	parameter := NormalizeAngle(rp.Ellipse().parametricAngle(direction))

	return int(parameter/rp.centralAngle()) % rp.N
}

// nearestIndex returns the index of the vertex that reaches farthest in the given direction of
// the world, at most half a central angle from the point that reaches farthest there. The
// direction is taken into the frame before the turn, and for unequal semi-axes onto the
// ellipse, where a semi-axis stretches it away from the direction itself; equal semi-axes take
// the direction as it is, the same expression vertex places them from. With one or two
// vertices the nearest can be more than a quarter turn off, so the reach is negative exactly
// where no vertex lies on that side.
func (rp RegularPolygon[T]) nearestIndex(direction float64) int {
	local := direction - rp.Angle
	if rp.Size.Width != rp.Size.Height {
		sin, cos := math.Sincos(local)
		local = math.Atan2(float64(rp.Size.Height)*sin, float64(rp.Size.Width)*cos)
	}

	return Mod(int(math.Round(local/rp.centralAngle())), rp.N)
}

// minMax returns the minimum and maximum corner of the vertices, the corners of Bounds, exact
// for an integer T where Bounds places a center: the pair the intersection tests reject shapes
// by before examining any edge, without placing a rectangle. Each side is set by the vertex
// nearest to that direction, at most half a step away, and reads that vertex as Vertices
// places it, so the corners are exactly those of Polygon().Bounds(), rounded alike for an
// integer T. An empty polygon returns two zero points.
func (rp RegularPolygon[T]) minMax() (Point[T], Point[T]) {
	if rp.IsEmpty() {
		return Point[T]{}, Point[T]{}
	}

	a := Point[T]{rp.vertex(rp.nearestIndex(Pi)).X, rp.vertex(rp.nearestIndex(3 * Pi / 2)).Y}
	b := Point[T]{rp.vertex(rp.nearestIndex(0)).X, rp.vertex(rp.nearestIndex(Pi / 2)).Y}

	return a, b
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(vector), rp.Size, rp.N, rp.Angle}
}

// MoveTo creates a new RegularPolygon with center at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N, rp.Angle}
}

// Scale creates a new RegularPolygon with size scaled by the given factor. A negative factor
// scales by its absolute value like Rectangle.Scale, since Size holds semi-axes and a negative
// one would place every vertex half a turn away rather than shrink the polygon; use Rotate for
// the half turn.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(factor).Abs(), rp.N, rp.Angle}
}

// ScaleXY creates a new RegularPolygon with size scaled by the given factors, negative ones by
// their absolute value like Scale.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(factorX, factorY).Abs(), rp.N, rp.Angle}
}

// Unscale creates a new RegularPolygon with size scaled by the inverse factor, the inverse of
// Scale and negative factors taken absolute like it. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) Unscale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Unscale(factor).Abs(), rp.N, rp.Angle}
}

// UnscaleXY creates a new RegularPolygon with size scaled by the inverse of the given factors,
// the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) UnscaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.UnscaleXY(factorX, factorY).Abs(), rp.N, rp.Angle}
}

// Resize creates a new RegularPolygon with the given semi-axes, taken absolute like RegPol.
func (rp RegularPolygon[T]) Resize(size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, size.Abs(), rp.N, rp.Angle}
}

// Canonical creates a new RegularPolygon in the form RegPol and Rotate build, with the size
// taken absolute and the angle normalized to [0, 2π): a well-formed polygon is returned as it
// is, up to the full turns Equal already ignores. It repairs a negative semi-axis written as a
// struct literal or decoded from JSON, which would place every vertex half a turn away, and
// brings a decoded angle onto the seam Rotate keeps. An angle within Delta of zero or of a
// full turn, the residue a chain of Rotate and Lerp calls can leave, becomes exactly zero, as
// Rectangle.Canonical makes it; Rotate itself never snaps.
func (rp RegularPolygon[T]) Canonical() RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Abs(), rp.N, snapAngle(rp.Angle)}
}

// Grow creates a new RegularPolygon with both semi-axes increased by amount, clamped to zero.
// The amount is added to each semi-axis, as Ellipse.Grow adds it, so each extent of Bounds
// grows by about twice it.
func (rp RegularPolygon[T]) Grow(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Grow(amount).AtLeastZero(), rp.N, rp.Angle}
}

// GrowXY creates a new RegularPolygon with the semi-axes increased by the given amounts along
// its own axes, clamped to zero. Each amount is added to that semi-axis, like Grow.
func (rp RegularPolygon[T]) GrowXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.GrowXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle}
}

// Shrink creates a new RegularPolygon with both semi-axes decreased by amount, clamped to zero.
func (rp RegularPolygon[T]) Shrink(amount T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Shrink(amount).AtLeastZero(), rp.N, rp.Angle}
}

// ShrinkXY creates a new RegularPolygon with the semi-axes decreased by the given amounts along
// its own axes, clamped to zero.
func (rp RegularPolygon[T]) ShrinkXY(amountX, amountY T) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ShrinkXY(amountX, amountY).AtLeastZero(), rp.N, rp.Angle}
}

// Lerp creates a new RegularPolygon in linear interpolation towards the given polygon, moving
// the center and the size together and turning the angle along the shorter arc with
// LerpAngle, so a tween from 350° to 10° passes through the top rather than the long way
// round; the result is normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like
// Point.Lerp, with the size taken absolute like RegPol. The vertex count does not
// interpolate: polygons with a different N have no shape between them, so Lerp panics for
// them, the convention RegularPolygonOrientationAngle follows for an orientation with no
// meaning, rather than hiding the mistake in an empty polygon.
func (rp RegularPolygon[T]) Lerp(polygon RegularPolygon[T], t float64) RegularPolygon[T] {
	if rp.N != polygon.N {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", rp.N, polygon.N))
	}

	return RegularPolygon[T]{
		rp.Center.Lerp(polygon.Center, t),
		rp.Size.Lerp(polygon.Size, t).Abs(),
		rp.N,
		NormalizeAngle(LerpAngle(rp.Angle, polygon.Angle, t)),
	}
}

// Transform creates a new RegularPolygon by applying the given matrix: the center moves, the
// semi-axes scale by the factors the matrix applies along its axes, the angle turns by the
// angle of the matrix or is mirrored about it for a reflection, and the vertex count is kept.
// A move, a turn, a reflection and a uniform scale are exact, and so is a scale of the axes
// while the polygon is not turned, exactly as they are for a rectangle. A shear, or a scale of
// the axes of a turned polygon, maps it onto one no value of this type holds and gives the
// nearest, on the factors Matrix.Scaling reports. Take that case exactly through
// Polygon().Transform. For integer T the
// center and the semi-axes are each rounded once.
func (rp RegularPolygon[T]) Transform[M Float](matrix Matrix[M]) RegularPolygon[T] {
	m := matrix.Float()
	scaling := m.Scaling()

	return RegularPolygon[T]{rp.Center.Transform(matrix), rp.Size.ScaleXY(scaling.X, scaling.Y).Abs(), rp.N, m.turnedAngle(rp.Angle)}
}

// Rotate creates a new RegularPolygon turned by the given angle (in radians) about its center,
// in the same sense as Vector.Rotate and as Rectangle.Rotate turns a box.
// The stored angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size, rp.N, NormalizeAngle(rp.Angle + angle)}
}

// AlignTo creates a new RegularPolygon moved so that its Anchor in the given direction lands
// on the point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (rp RegularPolygon[T]) AlignTo(direction Direction, point Point[T]) RegularPolygon[T] {
	return rp.Translate(point.Subtract(rp.Anchor(direction)))
}

// Contains reports whether the given point lies within the polygon, boundary included within
// Epsilon of T, the same closed convention as Polygon.Contains, and exactly what the Polygon
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
// edgeWalk every closed shape makes: zero for a point on an edge within Epsilon of T or
// inside by the even-odd rule, the squared distance to the nearest edge otherwise, and
// infinity for an empty polygon. It is a float64 even for an integer T, like
// Polygon.DistanceSquaredTo. Contains and IntersectsCircle are built on it.
func (rp RegularPolygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range rp.Edges() {
		if w.step(edge, point) {
			return 0
		}
	}

	return w.result()
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

// IntersectsRegularPolygon reports whether the polygons share a point: a vertex of one lies within the
// other, or an edge of one crosses an edge of the other, as Polygon.IntersectsPolygon decides on
// their Polygon forms, without building them. Touching polygons intersect, within Epsilon of
// T, and an empty polygon intersects nothing. Polygons whose Bounds do not overlap are
// rejected before any edge pair is examined.
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

// containsWithin is Contains for a caller that already holds the corners of Bounds, so the
// intersection tests place the box once and reuse it for every point they test.
func (rp RegularPolygon[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && rp.DistanceSquaredTo(point) == 0
}

// Equal checks if center point, size, number of vertices and angle are equal. Angles are
// compared with EqualAngle, so a full turn or the sign of an angle does not matter.
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool {
	return rp.Center.Equal(polygon.Center) && rp.Size.Equal(polygon.Size) && rp.N == polygon.N && EqualAngle(rp.Angle, polygon.Angle)
}

// IsZero checks if center point, size, number of vertices and angle are zero, comparing
// the angle like Equal so that a full turn counts as zero.
func (rp RegularPolygon[T]) IsZero() bool {
	return rp.Equal(RegularPolygon[T]{})
}

// IsEmpty checks if the polygon has no vertices.
func (rp RegularPolygon[T]) IsEmpty() bool {
	return rp.N < 1
}

// IsAligned reports whether the polygon is not turned: its Angle is exactly zero, as RegPol
// with no angle and Rotate by a full turn leave it, the same exact test Rectangle.IsAligned
// makes. No tolerance is applied; Canonical snaps a residue to zero. An aligned polygon of
// equal semi-axes has its first vertex at the end of the width semi-axis.
func (rp RegularPolygon[T]) IsAligned() bool {
	return rp.Angle == 0
}

// Polygon converts the regular polygon into a generic Polygon with the vertices Vertices
// iterates, in one allocation. A polygon with N < 1 has nil vertices, so its Polygon is zero
// like Pol(nil).
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
// is its inverse for the same vertex count. A polygon with N < 1 has no vertices and still
// names the ellipse its Size and Angle describe.
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
	return RegularPolygon[R]{rp.Center.Cast[R](), rp.Size.Cast[R](), rp.N, rp.Angle}
}

// Int converts the regular polygon to a RegularPolygon[int].
func (rp RegularPolygon[T]) Int() RegularPolygon[int] {
	return RegularPolygon[int]{rp.Center.Int(), rp.Size.Int(), rp.N, rp.Angle}
}

// Float converts the regular polygon to a RegularPolygon[float64].
func (rp RegularPolygon[T]) Float() RegularPolygon[float64] {
	return RegularPolygon[float64]{rp.Center.Float(), rp.Size.Float(), rp.N, rp.Angle}
}

// String returns the polygon in the form of its constructor: RegPol((x,y);WxH;n), center, size
// and vertex count, with the angle appended as RegPol((x,y);WxH;n;a) for a rotated polygon, as
// the JSON carries it only then.
func (rp RegularPolygon[T]) String() string {
	if rp.IsAligned() {
		return fmt.Sprintf("RegPol(%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N))
	}

	return fmt.Sprintf("RegPol(%s;%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N), String(rp.Angle))
}

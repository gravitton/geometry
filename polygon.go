package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
)

// Polygon is a closed outline through its vertices, each joined to the next and the last back to
// the first.
//
// The vertex count is not checked: a polygon of fewer than three vertices encloses no area, and
// every method still answers for it. Pol keeps the slice it is given, and every method returning a
// Polygon allocates a new one; a nil Points stays nil through all of them.
type Polygon[T Number] struct {
	Points []Point[T]
}

// Pol is shorthand for Polygon{vertices}.
func Pol[T Number](vertices []Point[T]) Polygon[T] {
	return Polygon[T]{vertices}
}

// Vertices returns an iterator over Points in order.
func (p Polygon[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// Edges returns an iterator over the edges in order, each from a vertex to the next and the last
// closing back to the first: none for an empty polygon, and one of zero length for a single vertex.
func (p Polygon[T]) Edges() iter.Seq[Segment[T]] {
	return edgesOf(p.Points)
}

// Centroid returns the center of the enclosed area. A polygon enclosing no area, its vertices on
// one line or its lobes cancelling exactly, gives the average of its vertices instead, and an empty
// one the zero point.
func (p Polygon[T]) Centroid() Point[T] {
	if p.IsEmpty() {
		return Point[T]{}
	}

	origin := p.Points[0].Float()
	offset := origin.Vector().Negate()

	var x, y, twiceArea, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		along, cross := shifted.Vector(), shifted.cross()

		x += float64((shifted.Start.X + shifted.End.X) * cross)
		y += float64((shifted.Start.Y + shifted.End.Y) * cross)
		twiceArea += cross
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if twiceArea == 0 || p.isFlatWithin(twiceArea, length) {
		return p.mean()
	}

	centroid := origin.AddXY(x/(3*twiceArea), y/(3*twiceArea))

	return centroid.Cast[T]()
}

// Area returns the enclosed area, whatever the winding. The lobes of a self-intersecting polygon
// that wind opposite ways cancel, and vertices on one line enclose none.
func (p Polygon[T]) Area() float64 {
	return math.Abs(p.twiceArea()) / 2
}

// Perimeter returns the total length of the edges.
func (p Polygon[T]) Perimeter() float64 {
	var perimeter float64
	for edge := range p.Edges() {
		perimeter += edge.Length()
	}

	return perimeter
}

// Inertia returns the polar second moment of area about the Centroid, the rotational inertia of the
// polygon at unit density, whatever the winding, and zero for a polygon enclosing no area. It holds
// for a simple polygon only: the lobes of a self-intersecting one cancel unevenly and can give a
// negative value.
func (p Polygon[T]) Inertia() float64 {
	if p.IsEmpty() {
		return 0
	}

	offset := p.Points[0].Vector().Float().Negate()

	var x, y, twiceArea, moment, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		a, b := shifted.Start.Vector(), shifted.End.Vector()
		along, cross := shifted.Vector(), shifted.cross()

		x += float64((a.X + b.X) * cross)
		y += float64((a.Y + b.Y) * cross)
		twiceArea += cross
		moment += float64((a.Dot(a) + a.Dot(b) + b.Dot(b)) * cross)
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if twiceArea == 0 || p.isFlatWithin(twiceArea, length) {
		return 0
	}

	centroid := Vector[float64]{x / (3 * twiceArea), y / (3 * twiceArea)}

	return math.Abs(moment)/12 - float64(math.Abs(twiceArea)/2*centroid.LengthSquared())
}

// Winding returns the sense in which the vertices run around the enclosed area, and WindingNone for
// a polygon enclosing none. A self-intersecting polygon winds the way its larger lobes do, and
// lobes of equal area either way give WindingNone.
func (p Polygon[T]) Winding() Winding {
	switch twiceArea := p.twiceArea(); {
	case twiceArea > 0:
		return WindingClockwise
	case twiceArea < 0:
		return WindingCounterClockwise
	default:
		return WindingNone
	}
}

// Bounds returns the axis-aligned bounding box, the zero box for an empty polygon.
func (p Polygon[T]) Bounds() Box[T] {
	a, b := p.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle about Centroid through the farthest vertex, the zero circle for
// an empty polygon. It is not the smallest enclosing circle, but around a simple polygon its radius
// is at most twice that one's. For integer T the radius is rounded up, so it still contains every
// vertex.
func (p Polygon[T]) BoundingCircle() Circle[T] {
	center := p.Centroid()
	origin := center.Float()

	var reach float64
	for _, vertex := range p.Points {
		reach = max(reach, origin.DistanceSquaredTo(vertex.Float()))
	}

	return Circle[T]{center, castUp[T](math.Sqrt(reach))}
}

// mean returns the average of the vertices, the Centroid of a polygon enclosing no area.
func (p Polygon[T]) mean() Point[T] {
	var x, y float64
	for _, vertex := range p.Points {
		x, y = x+float64(vertex.X), y+float64(vertex.Y)
	}

	n := float64(len(p.Points))

	return Point[T]{Cast[T](x / n), Cast[T](y / n)}
}

// twiceArea returns the signed shoelace sum, twice the enclosed area, positive for WindingClockwise
// and zero for a flat or empty polygon. The origin is moved to the first vertex so the products
// stay small.
func (p Polygon[T]) twiceArea() float64 {
	if p.IsEmpty() {
		return 0
	}

	offset := p.Points[0].Vector().Float().Negate()

	var twiceArea, length float64
	for edge := range p.Edges() {
		shifted := edge.Float().Translate(offset)
		along := shifted.Vector()

		twiceArea += shifted.cross()
		length += math.Abs(along.X) + math.Abs(along.Y)
	}

	if p.isFlatWithin(twiceArea, length) {
		return 0
	}

	return twiceArea
}

// minMax returns the corners of Bounds, two zero points for an empty polygon.
func (p Polygon[T]) minMax() (Point[T], Point[T]) {
	return minMaxOf(p.Points)
}

// Translate returns the polygon moved by the vector.
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return point.Add(vector)
	})}
}

// MoveTo returns the polygon moved so that its Centroid lies at the point. For integer T a centroid
// exactly on a half can land a unit off, where the move changes its sign and so the way it rounds.
func (p Polygon[T]) MoveTo(point Point[T]) Polygon[T] {
	return p.Translate(point.Subtract(p.Centroid()))
}

// Rotate returns the polygon turned by the angle in radians about its Centroid. For integer T only
// quarter turns keep the shape exactly.
func (p Polygon[T]) Rotate(angle float64) Polygon[T] {
	pivot := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return point.RotateAround(pivot, angle)
	})}
}

// Scale returns the polygon scaled about its Centroid by the factor; a negative factor also mirrors
// it through the centroid.
func (p Polygon[T]) Scale(factor float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY returns the polygon scaled about its Centroid by the factors; a negative factor also
// mirrors it along that axis.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Unscale returns the polygon scaled about its Centroid by the inverse of the factor. It panics for
// a zero factor.
func (p Polygon[T]) Unscale(factor float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Divide(factor))
	})}
}

// UnscaleXY returns the polygon scaled about its Centroid by the inverses of the factors. It panics
// for a zero factor.
func (p Polygon[T]) UnscaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Centroid()

	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).DivideXY(factorX, factorY))
	})}
}

// Transform returns the polygon with the matrix applied to every vertex.
func (p Polygon[T]) Transform[M Float](matrix Matrix[M]) Polygon[T] {
	return Polygon[T]{p.mapPoints(func(point Point[T]) Point[T] {
		return point.Transform(matrix)
	})}
}

// Grow returns the polygon with every edge moved out along its normal by the amount, each vertex
// along the bisector of its corner, a repeated vertex given once. A corner turning straight back,
// and a polygon of one point, stay where they are. A negative amount gives Shrink. A polygon that
// is not convex by IsConvex is grown only while the amount is less than its narrowest notch, past
// which the moved edges cross each other.
func (p Polygon[T]) Grow(amount float64) Polygon[T] {
	return Polygon[T]{p.AppendGrow(nil, amount)}
}

// AppendGrow appends the vertices of Grow(amount) to dst and returns the extended slice.
func (p Polygon[T]) AppendGrow(dst []Point[T], amount float64) []Point[T] {
	reversed := p.twiceArea() < 0
	if amount < 0 && p.IsConvex() {
		return p.shrink(dst, -amount, reversed)
	}

	return p.grow(dst, amount, reversed)
}

// Shrink returns the polygon with every edge moved in along its normal by the amount, Grow by its
// negative. A convex polygon is cut to the inner side of each moved edge, so an edge the others
// close over is dropped and a polygon shrunk past its middle is empty; for integer T both the moved
// edges and the corners are rounded. A polygon that is not convex has each vertex moved along the
// bisector of its corner instead, and is shrunk only while the amount is less than its narrowest
// neck, past which the moved edges cross each other.
func (p Polygon[T]) Shrink(amount float64) Polygon[T] {
	return p.Grow(-amount)
}

// AppendShrink appends the vertices of Shrink(amount) to dst and returns the extended slice.
func (p Polygon[T]) AppendShrink(dst []Point[T], amount float64) []Point[T] {
	return p.AppendGrow(dst, -amount)
}

// Lerp returns the polygon interpolated toward the given one, each vertex toward the vertex of the
// same index, and extrapolated outside [0, 1]. Outlines that start at different vertices or wind
// opposite ways twist through each other. It panics for polygons of different vertex counts.
func (p Polygon[T]) Lerp(polygon Polygon[T], t float64) Polygon[T] {
	if len(p.Points) != len(polygon.Points) {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", len(p.Points), len(polygon.Points)))
	}

	return Polygon[T]{p.mapPointsIndexed(func(i int, point Point[T]) Point[T] {
		return point.Lerp(polygon.Points[i], t)
	})}
}

// Reverse returns the polygon with its vertices in the opposite order, the same outline wound the
// other way.
func (p Polygon[T]) Reverse() Polygon[T] {
	points := slices.Clone(p.Points)
	slices.Reverse(points)

	return Polygon[T]{points}
}

// ConvexHull returns the smallest convex polygon containing every vertex, wound clockwise on screen
// from the least vertex by Point.Compare, and an empty polygon for an empty one. Its vertices are
// vertices of the polygon, with those on an edge of the hull or repeating another dropped, so
// collinear vertices give the two ends of their line and a single point gives itself. For integer T
// and float32 the turns are exact, so the hull is convex by IsConvex and the hull of a hull is
// itself; for float64 a vertex within rounding of the line through its neighbours may be kept.
func (p Polygon[T]) ConvexHull() Polygon[T] {
	return Polygon[T]{p.AppendConvexHull(nil)}
}

// AppendConvexHull appends the vertices of ConvexHull() to dst and returns the extended slice.
func (p Polygon[T]) AppendConvexHull(dst []Point[T]) []Point[T] {
	if p.IsEmpty() {
		return dst
	}

	chord := Segment[T]{slices.MinFunc(p.Points, Point[T].Compare), slices.MaxFunc(p.Points, Point[T].Compare)}
	from := len(dst)
	dst = append(dst, p.Points...)
	points := dst[from:]
	slices.SortFunc(points, func(a, b Point[T]) int {
		return p.compareAround(chord, a, b)
	})

	split := slices.IndexFunc(points, func(point Point[T]) bool {
		return p.returns(chord, point)
	})
	if split < 0 {
		split = len(points)
	}

	n := 0
	for _, point := range points[:split] {
		n = p.extendHull(points, n, 1, point)
	}

	greatest := n
	for _, point := range points[split:] {
		n = p.extendHull(points, n, greatest, point)
	}

	for n > greatest && p.turn(points[n-2], points[n-1], points[0]) <= 0 {
		n--
	}

	return dst[:from+n]
}

// Simplify returns the polygon without the vertices it needs no more to stay within the tolerance,
// by Ramer–Douglas–Peucker: every dropped vertex lies within the tolerance of the edge that
// replaces it. At zero tolerance it drops repeated vertices and those on the straight line between
// their kept neighbours, a fold back along that line included. The kept vertices stay in order from
// the first kept one, so simplifying again changes nothing; an outline within the tolerance of its
// least vertex by Point.Compare keeps that vertex alone. A negative tolerance is taken absolute,
// and an empty polygon gives an empty one.
func (p Polygon[T]) Simplify(tolerance float64) Polygon[T] {
	return Polygon[T]{p.AppendSimplify(nil, tolerance)}
}

// AppendSimplify appends the vertices of Simplify(tolerance) to dst and returns the extended slice.
// dst must not share memory with Points.
func (p Polygon[T]) AppendSimplify(dst []Point[T], tolerance float64) []Point[T] {
	if p.IsEmpty() {
		return dst
	}

	n, least := len(p.Points), 0
	for i, vertex := range p.Points {
		if vertex.Compare(p.Points[least]) < 0 {
			least = i
		}
	}

	from := len(dst)

	dst, wrapped := p.appendKept(append(slices.Grow(dst, n), p.Points[least]), least, least+n, math.Abs(tolerance))

	kept := dst[from:]
	slices.Reverse(kept)
	slices.Reverse(kept[:wrapped])
	slices.Reverse(kept[wrapped:])

	return dst
}

// Triangulate returns the triangles that fill the polygon as three indices into Points each, wound
// as the polygon is. A vertex the outline goes straight on or turns back at, a repeated one
// included, gets no triangle, so no triangle has zero area and there may be fewer than the vertices
// less two. A polygon of fewer than three vertices or enclosing no area gives nil.
//
// The polygon is taken to be simple by IsSimple: a self-intersecting one is left partly unfilled.
// The cost grows with the square of the vertex count, and with its cube at worst.
func (p Polygon[T]) Triangulate() []int {
	return p.AppendTriangulate(nil)
}

// AppendTriangulate appends the indices of Triangulate() to dst and returns the extended slice.
func (p Polygon[T]) AppendTriangulate(dst []int) []int {
	n := len(p.Points)

	twiceArea := p.twiceArea()
	if n < 3 || twiceArea == 0 {
		return dst
	}

	sense := 1.0
	if twiceArea < 0 {
		sense = -1
	}

	from := len(dst)
	dst = slices.Grow(dst, 3*(n-2))[:from+3*(n-2)]
	triangles := dst[from:]

	left := triangles[len(triangles)-n:]
	for i := range left {
		left[i] = i
	}

	count := 0
	for i, misses := 0, 0; len(left) > 2 && misses < len(left); {
		m := len(left)
		a, b, c := left[(i+m-1)%m], left[i], left[(i+1)%m]

		turn := float64(sense * p.turn(p.Points[a], p.Points[b], p.Points[c]))
		if !(turn == 0 || turn > 0 && p.cuts(left, a, b, c, sense)) {
			i, misses = (i+1)%m, misses+1

			continue
		}

		copy(left[1:i+1], left[:i])
		left = left[1:]
		if turn > 0 {
			triangles[count], triangles[count+1], triangles[count+2] = a, b, c
			count += 3
		}

		i, misses = (i+m-2)%(m-1), 0
	}

	return dst[:from+count]
}

// mapPoints returns fn of every vertex in order, as mapPointsIndexed.
func (p Polygon[T]) mapPoints[R any](fn func(Point[T]) R) []R {
	return p.mapPointsIndexed(func(_ int, point Point[T]) R {
		return fn(point)
	})
}

// mapPointsIndexed returns fn of every vertex and its index in a new slice, the one allocation of
// every mapping, and nil for a nil Points so IsZero survives the mapping.
func (p Polygon[T]) mapPointsIndexed[R any](fn func(int, Point[T]) R) []R {
	if p.Points == nil {
		return nil
	}

	mapped := make([]R, len(p.Points))
	for i, point := range p.Points {
		mapped[i] = fn(i, point)
	}

	return mapped
}

// grow appends the vertices with every edge moved out by the amount, each vertex by the sum of the
// unit normals of its two edges scaled to lie on both. Segment.Normal points into a polygon unless
// reversed, so the amount is negated then. A vertex opening a zero-length edge is skipped, and a
// polygon with no edge of any length is appended as it is.
func (p Polygon[T]) grow(dst []Point[T], amount float64, reversed bool) []Point[T] {
	if !reversed {
		amount = -amount
	}

	var previous Vector[float64]
	for i := len(p.Points) - 1; i >= 0 && !previous.hasDirection(); i-- {
		previous = p.edge(i).Float().Normal()
	}

	if !previous.hasDirection() {
		return append(dst, p.Points...)
	}

	previous = previous.Normalize()
	dst = slices.Grow(dst, len(p.Points))
	for i, point := range p.Points {
		normal := p.edge(i).Float().Normal()
		if !normal.hasDirection() {
			continue
		}

		normal = normal.Normalize()
		offset := previous.Add(normal).Multiply(ratio(amount, 1+float64(previous.Dot(normal))))
		dst = append(dst, point.Float().Add(offset).Cast[T]())
		previous = normal
	}

	return dst
}

// shrink appends the vertices of a convex polygon cut to the inner side of each edge moved in by
// the amount, and leaves dst as given where nothing remains.
func (p Polygon[T]) shrink(dst []Point[T], amount float64, reversed bool) []Point[T] {
	if reversed {
		amount = -amount
	}

	clip := edgeClip[T]{dst: slices.Grow(dst, 2*len(p.Points)+2), from: len(dst), points: p.Points, reversed: reversed}
	for edge := range p.Edges() {
		inward := edge.Float().Normal().Normalize().Multiply(amount)
		clip.add(edge.Translate(inward.Cast[T]()))
	}

	if points := clip.result(); len(points) > len(dst) {
		return points
	}

	return dst
}

// compareAround orders two points along the hull whose least and greatest vertex the chord joins:
// those on the side reached first, by Point.Compare, then those on the side returns reports, in
// reverse.
func (p Polygon[T]) compareAround(chord Segment[T], a, b Point[T]) int {
	switch ra, rb := p.returns(chord, a), p.returns(chord, b); {
	case ra != rb && ra:
		return 1
	case ra != rb:
		return -1
	case ra:
		return b.Compare(a)
	default:
		return a.Compare(b)
	}
}

// returns reports whether the point lies on the side of the chord a clockwise walk from Start
// reaches on the way back.
func (p Polygon[T]) returns(chord Segment[T], point Point[T]) bool {
	return p.turn(chord.Start, chord.End, point) > 0
}

// extendHull places the point at index n of the hull compacted into the front of points, after
// dropping the vertices from floor on that do not turn clockwise into it, and returns the new
// count. A repeat of the last vertex is skipped. The floor keeps the greatest vertex once passed,
// so a rounding turn can never drop an end.
func (p Polygon[T]) extendHull(points []Point[T], n, floor int, point Point[T]) int {
	if n > 0 && point == points[n-1] {
		return n
	}

	for n > floor && p.turn(points[n-2], points[n-1], point) <= 0 {
		n--
	}

	points[n] = point

	return n + 1
}

// turn returns the cross product of the edge from a to b and the edge from b to c: positive where
// the outline turns clockwise on screen at b, negative the other way, zero for collinear points.
func (Polygon[T]) turn(a, b, c Point[T]) float64 {
	return b.Float().Subtract(a.Float()).Cross(c.Float().Subtract(b.Float()))
}

// appendKept appends the vertices Simplify keeps strictly between the indices from and to, counted
// on past the last vertex around the outline. It returns the extended slice and how many of the
// appended vertices wrapped past the last index.
func (p Polygon[T]) appendKept(dst []Point[T], from, to int, tolerance float64) ([]Point[T], int) {
	n := len(p.Points)
	edge := Segment[T]{p.Points[from%n], p.Points[to%n]}

	farthest, distance := from, 0.0
	for i := from + 1; i < to; i++ {
		if d := edge.DistanceSquaredTo(p.Points[i%n]); d > distance {
			farthest, distance = i, d
		}
	}

	if farthest == from || lessOrEqualSquared(distance, tolerance, epsilonAt[T](max(edge.magnitude(), p.Points[farthest%n].magnitude()))) {
		return dst, 0
	}

	dst, before := p.appendKept(dst, from, farthest, tolerance)
	dst = append(dst, p.Points[farthest%n])
	dst, after := p.appendKept(dst, farthest, to, tolerance)

	return dst, before + farthest/n + after
}

// cuts reports whether the diagonal from a to c cuts off the ear at b: no other vertex left, by
// index or position, lies within the triangle or on its sides, on exact signs.
func (p Polygon[T]) cuts(left []int, a, b, c int, sense float64) bool {
	pa, pb, pc := p.Points[a], p.Points[b], p.Points[c]
	for _, i := range left {
		point := p.Points[i]
		if i == a || i == b || i == c || point == pa || point == pb || point == pc {
			continue
		}

		if sense*p.turn(pa, pb, point) >= 0 && sense*p.turn(pb, pc, point) >= 0 && sense*p.turn(pc, pa, point) >= 0 {
			return false
		}
	}

	return true
}

// DistanceTo returns the distance from the point to the polygon, zero inside, and infinity for an
// empty polygon.
func (p Polygon[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(p.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (p Polygon[T]) DistanceSquaredTo(point Point[T]) float64 {
	return p.walk(point).result()
}

// Nearest returns the point of the polygon nearest to the given point, the point itself inside, and
// the zero point for an empty polygon.
func (p Polygon[T]) Nearest(point Point[T]) Point[T] {
	return p.walk(point).nearest(point)
}

// Support returns the vertex farthest along the direction, the first in order where several reach
// alike, and the zero point for an empty polygon. For a polygon that is not convex it is the
// support of its ConvexHull.
func (p Polygon[T]) Support(direction Vector[float64]) Point[float64] {
	return supportOf(p.Points, direction)
}

// Contains reports whether the point lies in the polygon. The inside follows the even-odd rule, so
// a self-intersecting polygon leaves out the regions it winds around twice.
func (p Polygon[T]) Contains(point Point[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a, b := p.minMax()

	return p.containsWithin(point, a, b)
}

// ContainsCircle reports whether the circle lies entirely in the polygon.
func (p Polygon[T]) ContainsCircle(circle Circle[T]) bool {
	return p.walk(circle.Center).clears(circle)
}

// ContainsCapsule reports whether the capsule lies entirely in the polygon.
func (p Polygon[T]) ContainsCapsule(capsule Capsule[T]) bool {
	if !p.ContainsSegment(capsule.Segment) {
		return false
	}

	for edge := range p.Edges() {
		if !greaterOrEqualSquared(edge.DistanceSquaredToSegment(capsule.Segment), float64(capsule.Radius), epsilonAt[T](max(edge.magnitude(), capsule.magnitude()))) {
			return false
		}
	}

	return true
}

// ContainsSector reports whether the sector lies entirely in the polygon.
func (p Polygon[T]) ContainsSector(sector Sector[T]) bool {
	epsilon := sector.epsilon(p.Bounds().magnitude())
	for _, corner := range sector.corners() {
		if !p.containsPoint(corner, epsilon) {
			return false
		}
	}

	if !p.containsPoint(sector.Float().Centroid(), epsilon) {
		return false
	}

	for edge := range p.Edges() {
		if sector.enters(edge.Float()) {
			return false
		}
	}

	return true
}

// ContainsSegment reports whether the segment lies entirely in the polygon. The polygon is taken to
// be simple: a self-intersecting one is left wherever its edges cross, so it does not contain
// itself.
func (p Polygon[T]) ContainsSegment(segment Segment[T]) bool {
	a, b := p.minMax()

	return p.containsWithin(segment.Start, a, b) && p.containsWithin(segment.End, a, b) && p.keeps(segment, p.twiceArea())
}

// ContainsPolygon reports whether the given polygon lies entirely in this one.
func (p Polygon[T]) ContainsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range polygon.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// ContainsRectangle reports whether the rectangle lies entirely in the polygon.
func (p Polygon[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range rectangle.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the polygon.
func (p Polygon[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := p.minMax()
	twiceArea := p.twiceArea()

	for edge := range polygon.Edges() {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// ContainsBox reports whether the box lies entirely in the polygon.
func (p Polygon[T]) ContainsBox(box Box[T]) bool {
	a, b := p.minMax()
	twiceArea := p.twiceArea()

	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		if !p.containsWithin(edge.Start, a, b) || !p.keeps(edge, twiceArea) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the polygon and the circle share a point.
func (p Polygon[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsPolygon(p)
}

// IntersectsCapsule reports whether the polygon and the capsule share a point.
func (p Polygon[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsPolygon(p)
}

// IntersectsSector reports whether the polygon and the sector share a point.
func (p Polygon[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsPolygon(p)
}

// IntersectsSegment reports whether the polygon and the segment share a point.
func (p Polygon[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsPolygon(p)
}

// IntersectsRay reports whether the polygon and the ray share a point.
func (p Polygon[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsPolygon(p)
}

// IntersectsLine reports whether the polygon and the line share a point.
func (p Polygon[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsPolygon(p)
}

// IntersectsPolygon reports whether the polygons share a point.
func (p Polygon[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	if p.IsEmpty() || polygon.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(p.Points[0], a2, b2) || p.containsWithin(polygon.Points[0], a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
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

// IntersectsRectangle reports whether the polygon and the rectangle share a point.
func (p Polygon[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := rectangle.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if rectangle.containsWithin(p.Points[0], a2, b2) || p.containsWithin(rectangle.TopLeft(), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range rectangle.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectsRegularPolygon reports whether the polygon and the regular polygon share a point.
func (p Polygon[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if p.IsEmpty() || polygon.IsEmpty() {
		return false
	}

	a1, b1 := p.minMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(p.Points[0], a2, b2) || p.containsWithin(polygon.vertex(0), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range p.Edges() {
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
func (p Polygon[T]) IntersectsBox(box Box[T]) bool {
	if p.IsEmpty() {
		return false
	}

	a, b := p.minMax()

	if !overlaps(a, b, box.Min, box.Max) {
		return false
	}

	if box.Contains(p.Points[0]) || p.containsWithin(box.Min, a, b) {
		return true
	}

	probe := edgeProbe[T]{a: box.Min, b: box.Max}
	corners := box.corners()
	for edge := range p.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range edgesOf(corners[:]) {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (p Polygon[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionPolygon(p)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (p Polygon[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionPolygon(p)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (p Polygon[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionPolygon(p)
}

// ClipPolygon returns the part of the polygons in common, and false where there is none or neither
// is convex by IsConvex. One that is not convex, or this one where both are, is cut to the other,
// so the result keeps the winding and the vertices of the one cut. A common part of several pieces
// comes back as one outline joined along the cutting edges, polygons that only touch share a point
// or a stretch, and a polygon enclosing no area cuts everything away.
func (p Polygon[T]) ClipPolygon(polygon Polygon[T]) (Polygon[T], bool) {
	points := p.AppendClipPolygon(nil, polygon)

	return Polygon[T]{points}, len(points) > 0
}

// AppendClipPolygon appends the vertices of ClipPolygon(polygon) to dst and returns the extended
// slice.
func (p Polygon[T]) AppendClipPolygon(dst []Point[T], polygon Polygon[T]) []Point[T] {
	outline, region := p, polygon

	switch {
	case !p.IntersectsPolygon(polygon):
		return dst
	case polygon.IsConvex():
	case p.IsConvex():
		outline, region = polygon, p
	default:
		return dst
	}

	clip := edgeClip[T]{
		dst:      slices.Grow(dst, 2*(len(outline.Points)+len(region.Points))),
		from:     len(dst),
		points:   outline.Points,
		reversed: region.twiceArea() < 0,
	}
	for edge := range region.Edges() {
		clip.add(edge)
	}

	return clip.result()
}

// ClipRectangle returns the part of the polygon in the rectangle, and false where there is none or
// the rectangle has a zero width or height. The polygon need not be convex, and the result keeps
// its winding and vertices.
func (p Polygon[T]) ClipRectangle(rectangle Rectangle[T]) (Polygon[T], bool) {
	points := p.AppendClipRectangle(nil, rectangle)

	return Polygon[T]{points}, len(points) > 0
}

// AppendClipRectangle appends the vertices of ClipRectangle(rectangle) to dst and returns the
// extended slice.
func (p Polygon[T]) AppendClipRectangle(dst []Point[T], rectangle Rectangle[T]) []Point[T] {
	if !rectangle.Size.IsPositive() || !p.IntersectsRectangle(rectangle) {
		return dst
	}

	clip := edgeClip[T]{dst: slices.Grow(dst, 2*(len(p.Points)+4)), from: len(dst), points: p.Points}
	for edge := range rectangle.Edges() {
		clip.add(edge)
	}

	return clip.result()
}

// ClipRegularPolygon returns the part of the polygon in the regular polygon, and false where there
// is none or the regular polygon encloses no area, with fewer than three vertices or a zero
// semi-axis. The polygon need not be convex, and the result keeps its winding and vertices.
func (p Polygon[T]) ClipRegularPolygon(polygon RegularPolygon[T]) (Polygon[T], bool) {
	points := p.AppendClipRegularPolygon(nil, polygon)

	return Polygon[T]{points}, len(points) > 0
}

// AppendClipRegularPolygon appends the vertices of ClipRegularPolygon(polygon) to dst and returns
// the extended slice.
func (p Polygon[T]) AppendClipRegularPolygon(dst []Point[T], polygon RegularPolygon[T]) []Point[T] {
	if polygon.N < 3 || !polygon.Radius.IsPositive() || !p.IntersectsRegularPolygon(polygon) {
		return dst
	}

	clip := edgeClip[T]{dst: slices.Grow(dst, 2*(len(p.Points)+polygon.N)), from: len(dst), points: p.Points}
	for edge := range polygon.Edges() {
		clip.add(edge)
	}

	return clip.result()
}

// ClipBox returns the part of the polygon in the box, and false where there is none or the box has
// a zero width or height. The polygon need not be convex, and the result keeps its winding and
// vertices.
func (p Polygon[T]) ClipBox(box Box[T]) (Polygon[T], bool) {
	points := p.AppendClipBox(nil, box)

	return Polygon[T]{points}, len(points) > 0
}

// AppendClipBox appends the vertices of ClipBox(box) to dst and returns the extended slice.
func (p Polygon[T]) AppendClipBox(dst []Point[T], box Box[T]) []Point[T] {
	if !box.Size().IsPositive() || !p.IntersectsBox(box) {
		return dst
	}

	clip := edgeClip[T]{dst: slices.Grow(dst, 2*(len(p.Points)+4)), from: len(dst), points: p.Points}
	corners := box.corners()
	for edge := range edgesOf(corners[:]) {
		clip.add(edge)
	}

	return clip.result()
}

// Split returns the two parts the line cuts the polygon into, the first on the side Line.Side gives
// 1 for, each keeping the winding and vertices of the polygon. A vertex on the line belongs to
// both, so the parts meet exactly along it. A polygon wholly on one side leaves the other part
// empty, or the point or stretch where it touches the line. A part of several pieces comes back as
// one outline joined along the line, and a line with no direction leaves the whole polygon in the
// first part.
func (p Polygon[T]) Split(line Line[T]) (Polygon[T], Polygon[T]) {
	if p.Points == nil {
		return p, p
	}

	n := 2 * len(p.Points)
	buffer := make([]Point[T], 0, 2*n)
	front, back := p.AppendSplit(buffer[:0:n], buffer[n:n], line)

	return Polygon[T]{front}, Polygon[T]{back}
}

// AppendSplit appends the vertices of the two parts of Split(line) to front and back and returns
// the extended slices.
func (p Polygon[T]) AppendSplit(front, back []Point[T], line Line[T]) ([]Point[T], []Point[T]) {
	if !line.Vector.hasDirection() {
		return append(front, p.Points...), back
	}

	ahead := edgeClip[T]{dst: front, from: len(front), points: p.Points}
	ahead.cut(line, 1)

	behind := edgeClip[T]{dst: back, from: len(back), points: p.Points}
	behind.cut(line, -1)

	return ahead.result(), behind.result()
}

// Union returns the outline of the area the polygons cover together, and false where either is not
// convex or they share no area, touching included. The outline keeps the winding of this polygon
// and starts at its first vertex outside the given one; where one contains the other it has that
// one's vertices, in this winding. A vertex on the boundary of the other polygon gives way to the
// crossing beside it. For an integer type the outline is placed as for float64 and its crossings
// rounded once, so where the polygons overlap by less than a unit it can touch itself.
func (p Polygon[T]) Union(polygon Polygon[T]) (Polygon[T], bool) {
	points := p.AppendUnion(nil, polygon)

	return Polygon[T]{points}, len(points) > 0
}

// AppendUnion appends the vertices of Union(polygon) to dst and returns the extended slice.
func (p Polygon[T]) AppendUnion(dst []Point[T], polygon Polygon[T]) []Point[T] {
	if !p.Bounds().IntersectsBox(polygon.Bounds()) || !p.IsConvex() || !polygon.IsConvex() {
		return dst
	}

	outline := unionOutline[T]{dst: slices.Grow(dst, 2*(len(p.Points)+len(polygon.Points))), from: len(dst), a: p, b: polygon}
	outline.gather()

	return outline.walk()
}

// containsSupport reports whether the polygon, one that is not convex, holds a shape no Contains
// method names: no edge the shape meets has it reaching past the edge's outer side by more than the
// tolerance, and the polygon contains the middle of its supports along X. A shape across a notch is
// not held, however thin, and one touching a reflex or repeated vertex may answer either way within
// the tolerance.
func (p Polygon[T]) containsSupport(shape Shape[T]) bool {
	reversed := p.twiceArea() < 0
	bounds := shape.Bounds()
	for edge := range p.Edges() {
		if !edge.Bounds().IntersectsBox(bounds) {
			continue
		}

		difference := minkowski[T]{a: edge, b: shape}
		if !difference.intersects() {
			continue
		}

		outward := edge.Float().Normal().Negate().Normalize()
		if reversed {
			outward = outward.Negate()
		}

		if shape.Support(outward).Subtract(edge.Start.Float()).Dot(outward) > difference.tolerance() {
			return false
		}
	}

	middle := shape.Support(Vector[float64]{1, 0}).Lerp(shape.Support(Vector[float64]{-1, 0}), 0.5)

	return p.Contains(middle.Cast[T]())
}

// intersectsSupport reports whether the polygon, one with a vertex, and a shape no Intersects
// method names share a point: a vertex inside the shape, a point of the shape inside the polygon,
// or an edge meeting the shape by their supports, so a notch is not filled in.
func (p Polygon[T]) intersectsSupport(shape Shape[T]) bool {
	if shape.Contains(p.Points[0]) || p.Contains(shape.Nearest(p.Points[0])) {
		return true
	}

	for edge := range p.Edges() {
		if (minkowski[T]{a: edge, b: shape}).intersects() {
			return true
		}
	}

	return false
}

// nearestSupport returns the nearest points of the polygon, one that is not convex, and a shape, by
// the first edge of the least distance, so a notch is not filled in. Where they intersect with no
// edge reaching the shape, the shape lies within the polygon and both points are its point nearest
// the first vertex.
func (p Polygon[T]) nearestSupport(shape Shape[T], intersects bool) (Point[T], Point[T]) {
	var nearest Segment[T]

	gap := math.Inf(1)
	for edge := range p.Edges() {
		if distanceSquared := DistanceSquared[T](edge, shape); distanceSquared < gap {
			nearest, gap = edge, distanceSquared
		}
	}

	if intersects && gap > 0 {
		point := shape.Nearest(p.Points[0])

		return point, point
	}

	return Nearest[T](nearest, shape)
}

// collideSupport returns how the polygon, one that is not convex, and a shape meet, edge by edge so
// a notch is not filled in, and false where the shape crosses no edge. The polygon is the second
// shape of the contact where second is set. Each edge is met as the first shape whichever is asked,
// so an edge crossed in its very middle parts the two the same way in both orders, and the contact
// is turned round at the end.
//
// The deepest contact of the edges crossed is kept, so a shape in an inside corner parts from the
// deeper wall first. An edge is a wall with no inside, so its own contact may part the shapes into
// the polygon or past its end; such an edge is asked instead for the way out through its outer
// side, as far as the shape reaches behind it. Where that is longer than the contact it replaces
// the edge is passed over, and where every edge is, the shortest way out is returned. An edge of no
// length keeps its own contact.
func (p Polygon[T]) collideSupport(shape Shape[T], second bool) (Contact[T], bool) {
	deepest, shallowest := Contact[T]{Depth: math.Inf(-1)}, Contact[T]{Depth: math.Inf(1)}

	var deepestEdge, shallowestEdge Segment[T]

	reversed := p.twiceArea() < 0
	for edge := range p.Edges() {
		contact, ok := Collide[T](edge, shape)
		if !ok {
			continue
		}

		inward := edge.Float().Normal()
		if reversed {
			inward = inward.Negate()
		}

		if inward.hasDirection() && !(contact.Normal.Dot(inward) > 0) {
			inward = inward.Normalize()

			depth := max(shape.Support(inward).Subtract(edge.Start.Float()).Dot(inward), 0)
			outward := Contact[T]{inward, depth, minkowski[T]{a: edge, b: shape}.contact(inward)}
			if depth > contact.Depth {
				if depth < shallowest.Depth {
					shallowest, shallowestEdge = outward, edge
				}

				continue
			}

			contact = outward
		}

		if contact.Depth > deepest.Depth {
			deepest, deepestEdge = contact, edge
		}
	}

	contact, edge := deepest, deepestEdge
	if !(deepest.Depth >= 0) {
		if !(shallowest.Depth < math.Inf(1)) {
			return Contact[T]{}, false
		}

		contact, edge = shallowest, shallowestEdge
	}

	if second {
		normal := contact.Normal.Negate()

		return Contact[T]{normal, contact.Depth, minkowski[T]{a: shape, b: edge}.contact(normal)}, true
	}

	return contact, true
}

// containsWithin is Contains for a caller that already holds the extent of the vertices and tests
// many points against it.
func (p Polygon[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && p.DistanceSquaredTo(point) == 0
}

// containsMidpoint is Contains at the midpoint of a and b, taken in float64 so for an integer T it
// is not rounded onto either point. Segment.ClipPolygon judges each piece between two crossings by
// it.
func (p Polygon[T]) containsMidpoint(a, b Point[T]) bool {
	return p.containsPoint(a.Float().Midpoint(b.Float()), 0)
}

// containsPoint is Contains at a point of float64, for an integer T no lattice point, at the
// tolerance of T and no less than the floor.
func (p Polygon[T]) containsPoint(point Point[float64], floor float64) bool {
	inside := false
	for edge := range p.Edges() {
		if lessOrEqualSquared(edge.Float().distanceSquaredTo(point), 0, max(epsilonAt[T](max(edge.magnitude(), point.magnitude())), floor)) {
			return true
		}

		if edge.Float().crossesRay(point) {
			inside = !inside
		}
	}

	return inside
}

// walk folds every edge into the edgeWalk DistanceSquaredTo, Nearest and ContainsCircle read,
// stopping at an edge the point lies on.
func (p Polygon[T]) walk(point Point[T]) edgeWalk[T] {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range p.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	return w
}

// keeps reports whether the segment, its endpoints already found contained, stays within the simple
// polygon between them. It leaves by properly crossing an edge, by running through a vertex into
// the outside, or by running from an endpoint inside an edge to its outer side, each on exact signs
// in the sense of twiceArea. An endpoint at a vertex is left to the vertex, and a repeated vertex
// is skipped, so a vertex is judged by its distinct neighbours, never by the zero direction of a
// repeat.
func (p Polygon[T]) keeps(segment Segment[T], twiceArea float64) bool {
	previous := p.Points[0]
	for _, vertex := range slices.Backward(p.Points) {
		if vertex != p.Points[0] {
			previous = vertex

			break
		}
	}

	for edge := range p.Edges() {
		along := edge.Float().Vector()
		if !along.hasDirection() {
			continue
		}

		if segment.crosses(edge) {
			if _, touching := segment.touch(edge); !touching {
				return false
			}
		}

		if vertex := edge.Start; segment.Contains(vertex) {
			if p.leavesThrough(vertex, previous.Float().Subtract(vertex.Float()), along, segment, twiceArea) {
				return false
			}
		}

		for _, endpoint := range [2]Point[T]{segment.Start, segment.End} {
			if !edge.Contains(endpoint) || endpoint.coincides(edge.Start) || endpoint.coincides(edge.End) {
				continue
			}

			if p.leavesThrough(endpoint, along.Negate(), along, segment, twiceArea) {
				return false
			}
		}

		previous = edge.Start
	}

	return true
}

// leavesThrough reports whether the segment, passing a boundary point where the outline turns from
// the direction toward the previous vertex to the one toward the next, runs from there into the
// outside toward either endpoint.
func (p Polygon[T]) leavesThrough(point Point[T], toPrevious, toNext Vector[float64], segment Segment[T], twiceArea float64) bool {
	if twiceArea < 0 {
		toPrevious, toNext = toNext, toPrevious
	}

	start, end, origin := segment.Start.Float(), segment.End.Float(), point.Float()
	epsilon := epsilonAt[T](max(point.magnitude(), segment.magnitude()))

	return !p.admits(toPrevious, toNext, start.Subtract(origin), epsilon) ||
		!p.admits(toPrevious, toNext, end.Subtract(origin), epsilon)
}

// admits reports whether the offset from a boundary point stays within the polygon where the
// outline turns there between the two directions: on the inner side of both edges at a convex,
// straight or spike turn, of either at a reflex one, the turn decided on the exact sign.
func (p Polygon[T]) admits(toPrevious, toNext, offset Vector[float64], epsilon float64) bool {
	insideNext, insidePrevious := p.stays(toNext.Cross(offset), toNext, epsilon), p.stays(offset.Cross(toPrevious), toPrevious, epsilon)
	if toNext.Cross(toPrevious) >= 0 {
		return insideNext && insidePrevious
	}

	return insideNext || insidePrevious
}

// stays reports whether an offset lies on the inner side of an edge, given their cross product, or
// beyond its line by no more than the tolerance, on the gap (cross/√len²)² beside a segment. An
// offset within the tolerance of the point stays on every side, so a segment ending there is not
// judged by a direction it has not got.
func (Polygon[T]) stays(cross float64, edge Vector[float64], epsilon float64) bool {
	gap := cross / math.Sqrt(edge.LengthSquared())

	return cross >= 0 || lessOrEqualSquared(gap*gap, 0, epsilon)
}

// Equal reports whether the polygons have equal vertices in the same order. A nil and an empty
// Points are equal; only IsZero tells them apart.
func (p Polygon[T]) Equal(polygon Polygon[T]) bool {
	if len(p.Points) != len(polygon.Points) {
		return false
	}

	for i, vertex := range p.Points {
		if !vertex.Equal(polygon.Points[i]) {
			return false
		}
	}

	return true
}

// IsZero reports whether Points is nil.
func (p Polygon[T]) IsZero() bool {
	return p.Points == nil
}

// IsEmpty reports whether the polygon has no vertices.
func (p Polygon[T]) IsEmpty() bool {
	return len(p.Points) == 0
}

// IsConvex reports whether the polygon is convex and simple: it encloses an area, turns the same
// way at every vertex in either winding, and goes around once, so a star turning one way throughout
// is not convex. A repeated vertex is skipped and one on the line between its neighbours is
// allowed, but an edge doubling back is not. The turns are exact signs with no tolerance, so a
// float vertex a rounding error inside that line makes the polygon concave.
func (p Polygon[T]) IsConvex() bool {
	c := edgeConvexity{}
	for edge := range p.Edges() {
		if !c.step(edge.Float().Vector()) {
			return false
		}
	}

	return c.result()
}

// IsSimple reports whether the outline goes around without touching itself: no two edges share a
// point other than neighbours at their vertex, and no edge doubles back onto the one before. A
// repeated vertex is skipped, so the edges either side of it are neighbours. An empty polygon is
// simple and one enclosing no area is not; a NaN vertex meets no edge. The cost grows with the
// square of the vertex count.
func (p Polygon[T]) IsSimple() bool {
	if p.IsEmpty() {
		return true
	}

	if p.isFlat() {
		return false
	}

	for i := range p.Points {
		a := p.edge(i)
		if a.Start == a.End {
			continue
		}

		next := p.next(i)
		if b := p.edge(next); a.Contains(b.End) || b.Contains(a.Start) {
			return false
		}

		a1, b1 := a.minMax()
		for j := i + 1; j < len(p.Points); j++ {
			b := p.edge(j)
			if j == next || b.Start == b.End || p.next(j) == i {
				continue
			}

			if a2, b2 := b.minMax(); overlaps(a1, b1, a2, b2) && a.IntersectsSegment(b) {
				return false
			}
		}
	}

	return true
}

// isFlat reports whether the polygon encloses no area: fewer than three vertices, or every vertex
// on the segment between the least and the greatest by Point.Compare. The tolerance is on each
// vertex's distance from that line, never on the area.
func (p Polygon[T]) isFlat() bool {
	if len(p.Points) < 3 {
		return true
	}

	chord := Segment[T]{slices.MinFunc(p.Points, Point[T].Compare), slices.MaxFunc(p.Points, Point[T].Compare)}
	for _, vertex := range p.Points {
		if !chord.Contains(vertex) {
			return false
		}
	}

	return true
}

// isFlatWithin is isFlat for a caller holding the shoelace sum and the axis lengths of the edges:
// an outline within ε of a line sums to at most about 2ε times that length, so only a sum within
// the bound walks the vertices. The tolerance is taken at the first vertex widened by the length,
// never below that of any vertex.
func (p Polygon[T]) isFlatWithin(twiceArea, length float64) bool {
	epsilon := epsilonAt[T](p.Points[0].magnitude() + length)
	bound := float64(4*epsilon) + float64(float64(len(p.Points))*0x1p-48*length)

	return math.Abs(twiceArea) <= float64(bound*length) && p.isFlat()
}

// next returns the index of the next edge after i that has a direction. It is asked only of a
// polygon that is not flat, so the search ends; a NaN vertex never equals itself and ends it too.
func (p Polygon[T]) next(i int) int {
	for {
		i = (i + 1) % len(p.Points)
		if edge := p.edge(i); edge.Start != edge.End {
			return i
		}
	}
}

// edge returns the edge Edges yields at index i.
func (p Polygon[T]) edge(i int) Segment[T] {
	return Segment[T]{p.Points[i], p.Points[(i+1)%len(p.Points)]}
}

// Polyline returns the outline as a closed path, from the first vertex around and back to it, so
// IsClosed holds. An empty polygon gives an empty path, a nil Points a nil one.
func (p Polygon[T]) Polyline() Polyline[T] {
	if p.IsEmpty() {
		return Polyline[T](p)
	}

	points := make([]Point[T], 0, len(p.Points)+1)

	return Polyline[T]{append(append(points, p.Points...), p.Points[0])}
}

// Cast converts the polygon to a Polygon of another number type, rounding as Cast does.
func (p Polygon[T]) Cast[R Number]() Polygon[R] {
	return Polygon[R]{p.mapPoints(Point[T].Cast[R])}
}

// Int converts the polygon to a Polygon[int].
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{p.mapPoints(Point[T].Int)}
}

// Float converts the polygon to a Polygon[float64].
func (p Polygon[T]) Float() Polygon[float64] {
	return Polygon[float64]{p.mapPoints(Point[T].Float)}
}

// String returns the polygon as Pol((x,y);(x,y);...).
func (p Polygon[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Pol("...)
	for i, point := range p.Points {
		if i > 0 {
			buf = append(buf, ';')
		}

		buf = point.appendString(buf)
	}
	buf = append(buf, ')')

	return string(buf)
}

// MarshalJSON implements json.Marshaler.
func (p Polygon[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON implements json.Unmarshaler. The vertices are decoded into a fresh slice, so a
// slice the polygon shared is left untouched.
func (p *Polygon[T]) UnmarshalJSON(bytes []byte) error {
	var vertices []Point[T]
	if err := json.Unmarshal(bytes, &vertices); err != nil {
		return err
	}

	p.Points = vertices

	return nil
}

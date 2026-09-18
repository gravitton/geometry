package geom

import (
	"cmp"
	"fmt"
	"iter"
	"math"
	"slices"
)

// Line is a 2D line.
type Line[T Number] struct {
	Start Point[T] `json:"s"`
	End   Point[T] `json:"e"`
}

// Ln is shorthand for Line{start, end}.
func Ln[T Number](start, end Point[T]) Line[T] {
	return Line[T]{start, end}
}

// Vector returns the line as a vector, from start to end.
func (l Line[T]) Vector() Vector[T] {
	return l.End.Subtract(l.Start)
}

// Length returns the length of the line.
func (l Line[T]) Length() float64 {
	return l.Vector().Length()
}

// Angle returns the angle of the segment in radians, the angle of the vector from Start to End.
// A zero-length segment has no direction and gives 0, the angle Vector.Angle gives it.
func (l Line[T]) Angle() float64 {
	return l.Vector().Angle()
}

// Direction returns the direction nearest to the segment, from Start to End, or DirectionNone
// for a zero-length segment, as Vector.Direction judges it.
func (l Line[T]) Direction() Direction {
	return l.Vector().Direction()
}

// Midpoint returns the midpoint of the line, Lerp(0.5).
func (l Line[T]) Midpoint() Point[T] {
	return l.Start.Midpoint(l.End)
}

// Vertices returns the start and end points as a slice.
func (l Line[T]) Vertices() []Point[T] {
	return []Point[T]{l.Start, l.End}
}

// MinMax returns the minimum and maximum corner of the segment, the pair Rectangle.MinMax
// returns for its Bounds, exact for an integer T where Bounds places a center.
func (l Line[T]) MinMax() (Point[T], Point[T]) {
	return Point[T]{min(l.Start.X, l.End.X), min(l.Start.Y, l.End.Y)}, Point[T]{max(l.Start.X, l.End.X), max(l.Start.Y, l.End.Y)}
}

// Bounds returns the axis-aligned bounding rectangle.
func (l Line[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(l.MinMax())
}

// wedge returns Start × End in float64, the term the shoelace formula sums per edge. Cross is
// exact for parallel vectors, so a degenerate edge contributes exactly zero.
func (l Line[T]) wedge() float64 {
	return l.Start.Vector().Float().Cross(l.End.Vector().Float())
}

// Translate creates a new Line translated by the given vector.
func (l Line[T]) Translate(vector Vector[T]) Line[T] {
	return Line[T]{l.Start.Add(vector), l.End.Add(vector)}
}

// MoveTo creates a new Line with the start point moved to point and same length and direction.
func (l Line[T]) MoveTo(point Point[T]) Line[T] {
	return Line[T]{point, l.End.Add(point.Subtract(l.Start))}
}

// Scale creates a new Line uniformly scaled about its midpoint by the factor: the midpoint and
// direction stay, the length multiplies. A zero factor collapses the line onto its midpoint.
// For integer T the midpoint and both scaled points are rounded, so an odd span scaled by one
// is not exactly the same line.
func (l Line[T]) Scale(factor float64) Line[T] {
	return l.ScaleXY(factor, factor)
}

// ScaleXY creates a new Line scaled about its midpoint by the factors along X and Y, which
// changes the direction unless the factors are equal.
func (l Line[T]) ScaleXY(factorX, factorY float64) Line[T] {
	pivot := l.Midpoint()

	return Line[T]{
		pivot.Add(l.Start.Subtract(pivot).MultiplyXY(factorX, factorY)),
		pivot.Add(l.End.Subtract(pivot).MultiplyXY(factorX, factorY)),
	}
}

// Reverse creates a new Line with the start and end points swapped.
func (l Line[T]) Reverse() Line[T] {
	return Line[T]{l.End, l.Start}
}

// Lerp returns the point at the fraction t of the way from Start to End, extrapolating along
// the line outside [0, 1] like Point.Lerp. Midpoint is Lerp(0.5).
func (l Line[T]) Lerp(t float64) Point[T] {
	return l.Start.Lerp(l.End, t)
}

// Transform creates a new Line by applying the given matrix to both points, like Point.Transform.
func (l Line[T]) Transform[M Float](matrix Matrix[M]) Line[T] {
	return Line[T]{l.Start.Transform(matrix), l.End.Transform(matrix)}
}

// Rotate creates a new Line rotated by the given angle (in radians) about its midpoint, in the
// same sense as Vector.Rotate. For integer T the midpoint and both rotated points are rounded;
// only multiples of 90° keep the length exactly.
func (l Line[T]) Rotate(angle float64) Line[T] {
	pivot := l.Midpoint()

	return Line[T]{l.Start.RotateAround(pivot, angle), l.End.RotateAround(pivot, angle)}
}

// Contains reports whether the given point lies on the segment, within Epsilon of T, the same
// closed convention as Rectangle.Contains: it holds exactly where DistanceTo is zero.
func (l Line[T]) Contains(point Point[T]) bool {
	return l.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the segment:
// zero exactly where Contains holds, so a point within Epsilon of T of the segment is at
// distance zero rather than at the rounding error that put it there. A lattice point on a
// lattice segment is at zero without any tolerance, since the perpendicular distance comes
// from a cross product that is exact in float64 rather than from a projection.
func (l Line[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(l.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. It is a float64 even for an integer T, unlike Point.DistanceSquaredTo, since the
// nearest point of a segment is not a lattice point in general. A distance within Epsilon of T
// is snapped to zero, which is where Contains reads it, so a polygon boundary is walked without
// a square root per edge.
func (l Line[T]) DistanceSquaredTo(point Point[T]) float64 {
	epsilon := Epsilon[T]()

	distance := l.distanceSquaredTo(point)
	if LessOrEqualDelta(distance, 0, epsilon*epsilon) {
		return 0
	}

	return distance
}

// distanceSquaredTo returns the squared distance to the point with no tolerance applied, which
// DistanceSquaredTo snaps to zero within Epsilon of T.
func (l Line[T]) distanceSquaredTo(point Point[T]) float64 {
	direction, offset := l.Vector().Float(), point.Subtract(l.Start).Float()

	along := offset.Dot(direction)
	if along <= 0 {
		return offset.LengthSquared()
	}

	lengthSquared := direction.LengthSquared()
	if along >= lengthSquared {
		return point.Subtract(l.End).Float().LengthSquared()
	}

	cross := offset.Cross(direction)

	return cross * cross / lengthSquared
}

// DistanceToLine returns the distance between the nearest points of the two segments: zero
// exactly where Intersects holds, and otherwise the smallest distance from an endpoint of one
// to the other.
func (l Line[T]) DistanceToLine(line Line[T]) float64 {
	return math.Sqrt(l.DistanceSquaredToLine(line))
}

// DistanceSquaredToLine returns the squared distance DistanceToLine takes the root of, faster
// for comparisons: zero where the segments properly cross or an endpoint of one lies on the
// other within Epsilon of T, as DistanceSquaredTo snaps it.
func (l Line[T]) DistanceSquaredToLine(line Line[T]) float64 {
	if l.crosses(line) {
		return 0
	}

	return min(
		l.DistanceSquaredTo(line.Start), l.DistanceSquaredTo(line.End),
		line.DistanceSquaredTo(l.Start), line.DistanceSquaredTo(l.End),
	)
}

// Intersects reports whether the segments share a point, within Epsilon of T, the same closed
// convention as Contains: segments that touch at an endpoint or overlap collinearly intersect.
// It holds exactly where DistanceToLine is zero.
func (l Line[T]) Intersects(line Line[T]) bool {
	return l.DistanceSquaredToLine(line) == 0
}

// Intersection returns the point where the segments cross, and false when they do not. It
// answers exactly where Intersects holds, less the parallel case: parallel segments have no
// single crossing point and return false even where they overlap, which Intersects still
// reports. A zero-length segment is a point, parallel to nothing, and is the answer wherever it
// lies on the other segment. For integer T the crossing is rounded like every other result
// stored into T.
//
// The segments cross where Intersects says they do: either they properly cross, decided on
// cross products that are exact for an integer T, and the point is where the lines through
// them meet; or an endpoint of one lies on the other within Epsilon of T, and that endpoint is
// the point. The touch is judged on the endpoint's distance, like Contains, rather than on the
// fraction along the segment, which for a shallow crossing can put the same endpoint far
// outside the other segment.
func (l Line[T]) Intersection(line Line[T]) (Point[T], bool) {
	if l.crosses(line) {
		return l.crossing(line), true
	}

	if l.parallel(line) {
		return Point[T]{}, false
	}

	return l.touching(line)
}

// IntersectsCircle reports whether the segment and the circle share a point: the point of the
// segment closest to the center lies within the radius. Touching shapes intersect, within
// Epsilon of T, by the same comparison Circle.Contains makes, on the squared distance.
func (l Line[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.reaches(l.DistanceSquaredTo(circle.Center))
}

// IntersectionCircle returns the points where the segment crosses the circle boundary, from
// Start to End: two where it passes through, one where it is tangent or ends inside, within
// Epsilon of T like IntersectsCircle, and none where it misses, lies entirely inside, or the
// radius is negative. A chord shorter than Epsilon of T is a tangent and gives its midpoint;
// a longer one gives both ends even where it grazes the boundary within the tolerance. A segment inside crosses no boundary, so it returns none while
// IntersectsCircle still reports it. An endpoint within Epsilon of the boundary is a crossing
// in its own right, judged by the same comparison IntersectsCircle makes, so a shallow touch is
// not lost to the fraction along the chord and the two agree to the last bit. For integer T the points are rounded like
// every other result stored into T.
func (l Line[T]) IntersectionCircle(circle Circle[T]) []Point[T] {
	if circle.Radius < 0 {
		return nil
	}

	points := l.chord(circle)
	for _, endpoint := range [2]Point[T]{l.Start, l.End} {
		if circle.touches(circle.Center.Float().DistanceSquaredTo(endpoint.Float())) && !slices.ContainsFunc(points, endpoint.Equal) {
			points = append(points, endpoint)
		}
	}

	slices.SortFunc(points, l.nearer)

	return points
}

// IntersectionRectangle returns the points where the segment crosses the rectangle boundary,
// from Start to End: the crossings with its edges by Intersection, with a corner hit by two
// edges counted once. A segment inside crosses no boundary and returns none while
// IntersectsRectangle still reports it, and a segment along an edge is parallel to it and
// crosses only the edges at its ends, if it reaches them.
func (l Line[T]) IntersectionRectangle(rectangle Rectangle[T]) []Point[T] {
	return l.crossings(rectangle.edges())
}

// IntersectsRectangle reports whether the segment and the rectangle share a point: the start
// lies within the rectangle, or the segment crosses one of its edges. Touching shapes intersect,
// within Epsilon of T. A segment whose extent lies outside the rectangle is rejected before any
// edge is examined.
func (l Line[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a, b := rectangle.MinMax()

	if c, d := l.MinMax(); !overlaps(a, b, c, d) {
		return false
	}

	if l.Start.Between(a, b) {
		return true
	}

	for edge := range rectangle.edges() {
		if l.Intersects(edge) {
			return true
		}
	}

	return false
}

// IntersectsPolygon reports whether the segment and the polygon share a point, as
// Polygon.IntersectsLine does.
func (l Line[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsLine(l)
}

// IntersectionPolygon returns the points where the segment crosses the polygon boundary, as
// Polygon.IntersectionLine does.
func (l Line[T]) IntersectionPolygon(polygon Polygon[T]) []Point[T] {
	return polygon.IntersectionLine(l)
}

// crossings returns the points where the segment crosses the given edges, from Start to End:
// each crossing by Intersection, with a vertex hit by two edges counted once. It is the loop
// IntersectionRectangle and Polygon.IntersectionLine share.
func (l Line[T]) crossings(edges iter.Seq[Line[T]]) []Point[T] {
	var points []Point[T]
	for edge := range edges {
		point, ok := l.Intersection(edge)
		if !ok || slices.ContainsFunc(points, point.Equal) {
			continue
		}

		points = append(points, point)
	}

	slices.SortFunc(points, l.nearer)

	return points
}

// chord returns the points where the segment crosses the circle boundary as the chord it cuts:
// the crossings of the line through it, kept where they fall within the segment. A zero-length
// segment cuts no chord, and an endpoint on the boundary is left to IntersectionCircle. Whether
// the line reaches the circle is decided by the same comparison IntersectsCircle makes, so the
// two agree to the last bit. A chord whose ends lie within Epsilon of T of each other is a
// tangent and gives its midpoint alone: the tolerance collapses two crossings only where they
// would compare Equal, never a chord that merely grazes the boundary within the tolerance.
func (l Line[T]) chord(circle Circle[T]) []Point[T] {
	a := l.Float()
	direction := a.Vector()
	if !direction.hasDirection() {
		return nil
	}

	lengthSquared := direction.LengthSquared()
	along := -a.Start.Subtract(circle.Center.Float()).Dot(direction) / lengthSquared
	gapSquared := a.Lerp(along).DistanceSquaredTo(circle.Center.Float())

	if !circle.reaches(gapSquared) {
		return nil
	}

	radius := float64(circle.Radius)
	halfChord := math.Sqrt(max(radius*radius-gapSquared, 0))
	if 2*halfChord <= Epsilon[T]() {
		return l.pointsAt(along)
	}

	half := halfChord / math.Sqrt(lengthSquared)

	return l.pointsAt(along-half, along+half)
}

// pointsAt returns the points at the given fractions along the segment that lie within it,
// the endpoints included within Epsilon of T scaled to the length, so the tolerance is the same
// distance the Intersects methods apply. The fractions must be in increasing order.
func (l Line[T]) pointsAt(fractions ...float64) []Point[T] {
	var points []Point[T]
	for _, t := range fractions {
		if l.covers(t) {
			point := l.Float().Lerp(Clamp(t, 0, 1))
			points = append(points, Point[T]{Cast[T](point.X), Cast[T](point.Y)})
		}
	}

	return points
}

// nearer orders two points by their distance from Start, the order every boundary crossing
// method returns its points in.
func (l Line[T]) nearer(a, b Point[T]) int {
	return cmp.Compare(l.Start.DistanceSquaredTo(a), l.Start.DistanceSquaredTo(b))
}

// covers reports whether the fraction t of the way along the segment lies within it, the
// endpoints included within Epsilon of T scaled to the length, so that the tolerance is the
// same distance Contains and Intersects apply.
func (l Line[T]) covers(t float64) bool {
	epsilon := ratio(Epsilon[T](), l.Length())

	return LessOrEqualDelta(0, t, epsilon) && LessOrEqualDelta(t, 1, epsilon)
}

// crosses reports whether the segments properly cross: each has its endpoints on opposite sides
// of the other. Touching and collinear segments do not cross and are left to the endpoint
// distances, which cover them within the tolerance of the caller.
func (l Line[T]) crosses(line Line[T]) bool {
	return l.separates(line) && line.separates(l)
}

// crossing returns the point where the lines through two properly crossing segments meet, the
// fraction of the way along this segment from the same cross products crosses decided on.
func (l Line[T]) crossing(line Line[T]) Point[T] {
	a, b := l.Float(), line.Float()

	t := b.Start.Subtract(a.Start).Cross(b.Vector()) / a.Vector().Cross(b.Vector())
	point := a.Lerp(t)

	return Point[T]{Cast[T](point.X), Cast[T](point.Y)}
}

// parallel reports whether the segments run along the same direction. A zero-length segment
// has no direction and is parallel to nothing, so its point can still be found on the other.
func (l Line[T]) parallel(line Line[T]) bool {
	a, b := l.Vector(), line.Vector()

	return a.hasDirection() && b.hasDirection() && a.Float().Cross(b.Float()) == 0
}

// touching returns the endpoint of either segment that lies on the other, within Epsilon of T
// as Contains judges it, and false when there is none. Non-parallel segments that do not
// properly cross can share a point only this way.
func (l Line[T]) touching(line Line[T]) (Point[T], bool) {
	switch {
	case l.Contains(line.Start):
		return line.Start, true
	case l.Contains(line.End):
		return line.End, true
	case line.Contains(l.Start):
		return l.Start, true
	case line.Contains(l.End):
		return l.End, true
	default:
		return Point[T]{}, false
	}
}

// separates reports whether the endpoints of the given segment lie strictly on opposite sides
// of the line through this one.
func (l Line[T]) separates(line Line[T]) bool {
	direction := l.Vector().Float()
	start := direction.Cross(line.Start.Subtract(l.Start).Float())
	end := direction.Cross(line.End.Subtract(l.Start).Float())

	return (start > 0 && end < 0) || (start < 0 && end > 0)
}

// crossesRay reports whether a ray cast from the point along +X crosses the segment, counting
// an endpoint on the ray only when it is the lower one, so a ray through a vertex is counted
// once by the two edges that share it. It is the step of the even-odd rule Polygon.Contains uses.
//
// The ray crosses when the point lies on the side of the segment facing -X: the left side of a
// segment running toward +Y, the right side of one running toward -Y. The side comes from the
// sign of a cross product, which needs no division by the segment's Y span.
func (l Line[T]) crossesRay(point Point[T]) bool {
	start, end, p := l.Start.Float(), l.End.Float(), point.Float()

	if (start.Y > p.Y) == (end.Y > p.Y) {
		return false
	}

	upward := end.Y > start.Y
	left := end.Subtract(start).Cross(p.Subtract(start)) > 0

	return left == upward
}

// Equal checks if the start and end points of the lines are equal.
func (l Line[T]) Equal(line Line[T]) bool {
	return l.Start.Equal(line.Start) && l.End.Equal(line.End)
}

// IsZero checks if start and end points are zero.
func (l Line[T]) IsZero() bool {
	return l.Equal(Line[T]{})
}

// Int converts the line to a Line[int].
func (l Line[T]) Int() Line[int] {
	return Line[int]{l.Start.Int(), l.End.Int()}
}

// Float converts the line to a Line[float64].
func (l Line[T]) Float() Line[float64] {
	return Line[float64]{l.Start.Float(), l.End.Float()}
}

// String returns the line in the form of its constructor: Ln((x,y);(x,y)).
func (l Line[T]) String() string {
	return fmt.Sprintf("Ln(%s;%s)", l.Start.String(), l.End.String())
}

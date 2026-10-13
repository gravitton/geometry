package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
)

// Polyline is a 2D open path given by its points, each joined to the next by a straight edge:
// a road, a patrol route or a stroke drawn by a pointer, where a Polygon closes back to its
// first vertex and encloses an area. The point count is not checked: a path of one point is
// that point and an empty one has no points at all, and every method still answers for both.
//
// Points is shared, not copied: PolLn keeps the slice it is given and every method that returns
// a Polyline allocates a new one. The path is immutable as long as its caller does not write
// into that slice. A nil Points stays nil through every method, so IsZero holds after
// Translate, Transform, Int or Float.
//
// The path is no Shape: it has no inside and no Support that would not hand the functions over
// two shapes its hull. It is a wall of several edges a segment, a ray or a line is cast at, by
// their Intersection and Hit methods naming it. It has no center either, so it moves by
// Translate and turns and scales by Transform alone.
type Polyline[T Number] struct {
	Points []Point[T]
}

// PolLn is shorthand for Polyline{points}.
func PolLn[T Number](points []Point[T]) Polyline[T] {
	return Polyline[T]{points}
}

// Vertices iterates Points in order, the form every shape with an outline offers. Index Points
// directly where a position is needed.
func (p Polyline[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// Edges iterates the edges of the path in order, each from a point to the next, without
// allocating; collect them with slices.Collect where a slice is needed. The path is open, so
// no edge returns to the first point and there is one edge fewer than there are points. A
// single point yields one zero-length edge, so it is measured to as the point it is, and an
// empty path yields nothing.
func (p Polyline[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		if len(p.Points) == 1 {
			yield(Segment[T]{p.Points[0], p.Points[0]})

			return
		}

		for i := 1; i < len(p.Points); i++ {
			if !yield(Segment[T]{p.Points[i-1], p.Points[i]}) {
				return
			}
		}
	}
}

// Length returns the length of the path, the sum of the lengths of its edges.
func (p Polyline[T]) Length() float64 {
	var length float64
	for edge := range p.Edges() {
		length += edge.Length()
	}

	return length
}

// Bounds returns the axis-aligned bounding box of the points, or the zero box for a path
// without points.
func (p Polyline[T]) Bounds() Box[T] {
	a, b := minMaxOf(p.Points)

	return Box[T]{a, b}
}

// Translate creates a new Polyline translated by the given vector (applied to all points).
func (p Polyline[T]) Translate(vector Vector[T]) Polyline[T] {
	return Polyline[T]{Polygon[T](p).mapPoints(func(point Point[T]) Point[T] {
		return point.Add(vector)
	})}
}

// Transform creates a new Polyline by applying the given matrix to every point, like
// Point.Transform: the way a path is turned or scaled, about the pivot the matrix is built
// around, since a path has no center of its own.
func (p Polyline[T]) Transform[M Float](matrix Matrix[M]) Polyline[T] {
	return Polyline[T]{Polygon[T](p).mapPoints(func(point Point[T]) Point[T] {
		return point.Transform(matrix)
	})}
}

// Lerp creates a new Polyline in linear interpolation towards the given path, each point moving
// towards the point of the same index like Point.Lerp, extrapolating outside [0, 1]. Paths with
// a different point count have no path between them, so Lerp panics for them.
func (p Polyline[T]) Lerp(polyline Polyline[T], t float64) Polyline[T] {
	if len(p.Points) != len(polyline.Points) {
		panic(fmt.Sprintf("geom: lerp between polylines of %d and %d points", len(p.Points), len(polyline.Points)))
	}

	return Polyline[T]{Polygon[T](p).mapPointsIndexed(func(i int, point Point[T]) Point[T] {
		return point.Lerp(polyline.Points[i], t)
	})}
}

// Reverse creates a new Polyline with the points in the opposite order, the same path walked
// from its end.
func (p Polyline[T]) Reverse() Polyline[T] {
	points := slices.Clone(p.Points)
	slices.Reverse(points)

	return Polyline[T]{points}
}

// PointAt returns the point at the fraction t of the Length along the path, PointAt(0) its
// first point and PointAt(1) its last. Outside [0, 1] it extrapolates along the first edge and
// the last that have a length, like Segment.PointAt, so a path of two points answers as the
// segment between them and a repeated point at either end is passed over. An empty path returns
// the zero point.
//
// The edges are walked from the start on every call. A caller placing many points along one
// path ranges Every, which walks it once.
func (p Polyline[T]) PointAt(t float64) Point[T] {
	target := float64(t * p.Length())

	var edge Segment[T]
	var before, walked float64
	for next := range p.Edges() {
		length := next.Length()
		if length == 0 && walked > 0 {
			continue
		}

		edge, before = next, walked
		walked += length

		if walked > before && target <= walked {
			break
		}
	}

	return edge.PointAt(ratio(target-before, walked-before))
}

// Every iterates the points a step apart along the path, measured along its length: the first
// point, then one every step from it, in one walk over the edges and without allocating. The
// last point of the path is among them only where a whole number of steps reaches it. A step
// that is not positive has nowhere to go and yields the first point alone, and an empty path
// yields nothing. Each point is placed at its count of steps from the first, never at a sum of
// them, so the rounding of one does not carry to the next. The edges are read from Points
// rather than ranged from Edges, since an iterator ranged inside another allocates.
func (p Polyline[T]) Every(step float64) iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		if p.IsEmpty() {
			return
		}

		if len(p.Points) == 1 || !(step > 0) {
			yield(p.Points[0])

			return
		}

		var walked, next float64

		count := 0
		for i := 1; i < len(p.Points); i++ {
			edge := Segment[T]{p.Points[i-1], p.Points[i]}
			length := edge.Length()

			for next <= walked+length {
				if !yield(edge.PointAt(ratio(next-walked, length))) {
					return
				}

				count++
				next = float64(float64(count) * step)
			}

			walked += length
		}
	}
}

// Simplify creates a new Polyline with the points that lie within the tolerance of the path
// through the others removed, by the Ramer–Douglas–Peucker algorithm, as Polygon.Simplify
// thins a closed outline: the first and the last point are kept, and between two kept points
// the one farthest from the edge joining them is kept where it lies beyond the tolerance, and
// every point between them dropped where it does not. The kept points are points of the path
// in its order, so a path with nothing to drop is returned equal to itself and simplifying
// again changes nothing. A negative tolerance is taken absolute.
func (p Polyline[T]) Simplify(tolerance float64) Polyline[T] {
	return Polyline[T]{p.AppendSimplify(nil, tolerance)}
}

// AppendSimplify appends the points of the path Simplify returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room for every point
// of the path. An empty path appends nothing. The points are read while dst is written, so dst
// must not share memory with Points.
func (p Polyline[T]) AppendSimplify(dst []Point[T], tolerance float64) []Point[T] {
	n := len(p.Points)
	if n == 0 {
		return dst
	}

	dst = append(slices.Grow(dst, n), p.Points[0])
	if n == 1 {
		return dst
	}

	dst, _ = Polygon[T](p).appendKept(dst, 0, n-1, math.Abs(tolerance))

	return append(dst, p.Points[n-1])
}

// DistanceTo returns the distance from the given point to the nearest point of the path: zero
// exactly where Contains holds, and otherwise the distance to the nearest edge. An empty path
// is infinitely far from every point.
func (p Polyline[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(p.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, in one pass over the edges: zero for a point on an edge within the tolerance,
// snapped the way Segment.DistanceSquaredTo snaps it, the squared distance to the nearest edge
// otherwise, and infinity for an empty path. It is a float64 even for an integer T, since the
// nearest point of an edge is not a lattice point in general.
func (p Polyline[T]) DistanceSquaredTo(point Point[T]) float64 {
	return p.walk(point).result()
}

// Nearest returns the point of the path nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the foot on the nearest edge, the first of them where
// several are as near, read off the same walk DistanceSquaredTo makes. For integer T the foot
// is rounded once and can land off the path, where Contains rejects it. An empty path returns
// the zero point.
func (p Polyline[T]) Nearest(point Point[T]) Point[T] {
	return p.walk(point).nearest(point)
}

// NearestFraction returns the fraction of the Length at which the point of the path nearest to
// the given point lies, the one PointAt takes back to it: how far along a route a follower is,
// from zero at the first point to one at the last. Where several edges are as near it is the
// first of them, as Nearest gives it, an edge the point lies on within the tolerance before any
// other, so a path that retraces itself answers on the way out. An empty path and a path of no
// length give zero, and a NaN point is nearest to nothing and gives NaN.
func (p Polyline[T]) NearestFraction(point Point[T]) float64 {
	var walked float64

	along, nearest := math.NaN(), math.Inf(1)
	for edge := range p.Edges() {
		length := edge.Length()
		if distance := edge.DistanceSquaredTo(point); distance < nearest {
			nearest = distance
			along = walked + float64(edge.NearestFraction(point)*length)
		}

		walked += length
	}

	return ratio(along, walked)
}

// Contains reports whether the given point lies on the path, within the tolerance: it holds
// exactly where DistanceTo is zero.
func (p Polyline[T]) Contains(point Point[T]) bool {
	return p.DistanceSquaredTo(point) == 0
}

// IntersectionSegment returns the points where the segment crosses the path, as
// Segment.IntersectionPolyline does, which has the Append form.
func (p Polyline[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionPolyline(p)
}

// IntersectionRay returns the points where the ray crosses the path, as
// Ray.IntersectionPolyline does, which has the Append form.
func (p Polyline[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionPolyline(p)
}

// IntersectionLine returns the points where the line crosses the path, as
// Line.IntersectionPolyline does, which has the Append form.
func (p Polyline[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionPolyline(p)
}

// walk folds every edge into the edgeWalk DistanceSquaredTo and Nearest read, stopping at an
// edge the point lies on within the tolerance. A path has no inside, so the crossings the walk
// counts decide nothing and are dropped.
func (p Polyline[T]) walk(point Point[T]) edgeWalk[T] {
	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range p.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	w.inside = false

	return w
}

// Equal checks if two paths have the same points in the same order. A nil and an empty Points
// are equal, having the same none; only IsZero tells them apart.
func (p Polyline[T]) Equal(polyline Polyline[T]) bool {
	return Polygon[T](p).Equal(Polygon[T](polyline))
}

// IsZero checks if the points slice is nil.
func (p Polyline[T]) IsZero() bool {
	return p.Points == nil
}

// IsEmpty checks if number of points is zero.
func (p Polyline[T]) IsEmpty() bool {
	return len(p.Points) == 0
}

// IsClosed reports whether the path ends where it starts, its last point equal to its first:
// the path Polygon.Polyline gives, and the one Polygon takes the repeated point off. A path of
// fewer than two points has no end to return from and is not closed.
func (p Polyline[T]) IsClosed() bool {
	n := len(p.Points)

	return n > 1 && p.Points[0].Equal(p.Points[n-1])
}

// Polygon converts the path to the Polygon of its points, closed back from the last to the
// first, the inverse of Polygon.Polyline: a path that IsClosed already ends at its first point,
// which is left out so the polygon has no edge of zero length. Points is shared, not copied, and
// the polygon of a closed path has no room past its last vertex, so appending to its Points
// never writes over the point left out.
func (p Polyline[T]) Polygon() Polygon[T] {
	if p.IsClosed() {
		n := len(p.Points) - 1

		return Polygon[T]{p.Points[:n:n]}
	}

	return Polygon[T](p)
}

// Cast converts the path to a Polyline of another number type, rounding as Cast does.
func (p Polyline[T]) Cast[R Number]() Polyline[R] {
	return Polyline[R]{Polygon[T](p).mapPoints(Point[T].Cast[R])}
}

// Int converts the path to a Polyline[int].
func (p Polyline[T]) Int() Polyline[int] {
	return Polyline[int]{Polygon[T](p).mapPoints(Point[T].Int)}
}

// Float converts the path to a Polyline[float64].
func (p Polyline[T]) Float() Polyline[float64] {
	return Polyline[float64]{Polygon[T](p).mapPoints(Point[T].Float)}
}

// String returns the path in the form of its constructor: PolLn((x,y);(x,y);...), the points
// separated as every other shape separates its fields.
func (p Polyline[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "PolLn("...)
	for i, point := range p.Points {
		if i > 0 {
			buf = append(buf, ';')
		}

		buf = point.appendString(buf)
	}
	buf = append(buf, ')')

	return string(buf)
}

// MarshalJSON implements json.Marshaler: the array of the points, as a Polygon is the array of
// its vertices.
func (p Polyline[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON implements json.Unmarshaler. The points are decoded into a fresh slice, so a
// slice the path shared before decoding is left untouched.
func (p *Polyline[T]) UnmarshalJSON(bytes []byte) error {
	var points []Point[T]
	if err := json.Unmarshal(bytes, &points); err != nil {
		return err
	}

	p.Points = points

	return nil
}

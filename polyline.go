package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
)

// Polyline is an open path through its points, each joined to the next by a straight edge.
//
// The point count is not checked: a path of one point is that point, and an empty one has none.
// PolLn keeps the slice it is given, and every method returning a Polyline allocates a new one; a
// nil Points stays nil through all of them. The path has no inside, so it is no Shape, and no
// center, so it turns and scales by Transform alone.
type Polyline[T Number] struct {
	Points []Point[T]
}

// PolLn is shorthand for Polyline{points}.
func PolLn[T Number](points []Point[T]) Polyline[T] {
	return Polyline[T]{points}
}

// Vertices returns an iterator over Points in order.
func (p Polyline[T]) Vertices() iter.Seq[Point[T]] {
	return slices.Values(p.Points)
}

// Edges returns an iterator over the edges in order, each from a point to the next, one fewer than
// the points and none back to the first. A single point yields one edge of zero length, and an
// empty path nothing.
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

// Length returns the sum of the lengths of the edges.
func (p Polyline[T]) Length() float64 {
	var length float64
	for edge := range p.Edges() {
		length += edge.Length()
	}

	return length
}

// Bounds returns the axis-aligned bounding box, the zero box for an empty path.
func (p Polyline[T]) Bounds() Box[T] {
	a, b := minMaxOf(p.Points)

	return Box[T]{a, b}
}

// Translate returns the path moved by the vector.
func (p Polyline[T]) Translate(vector Vector[T]) Polyline[T] {
	return Polyline[T]{Polygon[T](p).mapPoints(func(point Point[T]) Point[T] {
		return point.Add(vector)
	})}
}

// Transform returns the path with the matrix applied to every point.
func (p Polyline[T]) Transform[M Float](matrix Matrix[M]) Polyline[T] {
	return Polyline[T]{Polygon[T](p).mapPoints(func(point Point[T]) Point[T] {
		return point.Transform(matrix)
	})}
}

// Lerp returns the path interpolated toward the given one, each point toward the point of the same
// index, and extrapolated outside [0, 1]. It panics for paths of different point counts.
func (p Polyline[T]) Lerp(polyline Polyline[T], t float64) Polyline[T] {
	if len(p.Points) != len(polyline.Points) {
		panic(fmt.Sprintf("geom: lerp between polylines of %d and %d points", len(p.Points), len(polyline.Points)))
	}

	return Polyline[T]{Polygon[T](p).mapPointsIndexed(func(i int, point Point[T]) Point[T] {
		return point.Lerp(polyline.Points[i], t)
	})}
}

// Reverse returns the path with its points in the opposite order.
func (p Polyline[T]) Reverse() Polyline[T] {
	points := slices.Clone(p.Points)
	slices.Reverse(points)

	return Polyline[T]{points}
}

// PointAt returns the point at the fraction t of the Length, the first point at zero and the last
// at one. Outside [0, 1] it extrapolates along the first or last edge of nonzero length. An empty
// path gives the zero point.
//
// The edges are walked on every call; to place many points, range Every.
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

// Every returns an iterator over the points a step apart along the path, from the first point on,
// in one walk. The last point is among them only where a whole number of steps reaches it. A step
// that is not positive yields the first point alone, and an empty path nothing. Each point is
// placed at its count of steps, so rounding does not accumulate.
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

// Simplify returns the path with the points within the tolerance of the path through the others
// removed, by Ramer–Douglas–Peucker, keeping the first and the last. The kept points stay in order,
// so simplifying again changes nothing. A negative tolerance is taken absolute.
func (p Polyline[T]) Simplify(tolerance float64) Polyline[T] {
	return Polyline[T]{p.AppendSimplify(nil, tolerance)}
}

// AppendSimplify appends the points of Simplify(tolerance) to dst and returns the extended slice.
// dst must not share memory with Points.
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

// DistanceTo returns the distance from the point to the path, zero on it, and infinity for an empty
// path.
func (p Polyline[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(p.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (p Polyline[T]) DistanceSquaredTo(point Point[T]) float64 {
	return p.walk(point).result()
}

// Nearest returns the point of the path nearest to the given point, the point itself on the path,
// the first edge winning a tie, and the zero point for an empty path.
func (p Polyline[T]) Nearest(point Point[T]) Point[T] {
	return p.walk(point).nearest(point)
}

// NearestFraction returns the fraction of the Length at which PointAt gives Nearest of the point,
// from zero at the first point to one at the last. The first nearest edge wins a tie, so a path
// that retraces itself answers on the way out. An empty path and one of no length give zero, and a
// NaN point gives NaN.
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

// Contains reports whether the point lies on the path.
func (p Polyline[T]) Contains(point Point[T]) bool {
	return p.DistanceSquaredTo(point) == 0
}

// IntersectionSegment returns the points where the segment crosses the path, from Start to End.
func (p Polyline[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionPolyline(p)
}

// IntersectionRay returns the points where the ray crosses the path, from its Origin on.
func (p Polyline[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionPolyline(p)
}

// IntersectionLine returns the points where the line crosses the path, in order along its
// Direction.
func (p Polyline[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionPolyline(p)
}

// walk folds the edges into an edgeWalk, stopping at one the point lies on. A path has no inside,
// so the crossings it counts are dropped.
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

// Equal reports whether the paths have equal points in the same order. A nil and an empty Points
// are equal; only IsZero tells them apart.
func (p Polyline[T]) Equal(polyline Polyline[T]) bool {
	return Polygon[T](p).Equal(Polygon[T](polyline))
}

// IsZero reports whether Points is nil.
func (p Polyline[T]) IsZero() bool {
	return p.Points == nil
}

// IsEmpty reports whether the path has no points.
func (p Polyline[T]) IsEmpty() bool {
	return len(p.Points) == 0
}

// IsClosed reports whether the path has at least two points and ends at its first.
func (p Polyline[T]) IsClosed() bool {
	n := len(p.Points)

	return n > 1 && p.Points[0].Equal(p.Points[n-1])
}

// Polygon returns the polygon of the points, the inverse of Polygon.Polyline: a closed path drops
// its repeated last point. Points is shared, capped so appending to the polygon never overwrites
// the point dropped.
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

// String returns the path in the form of its constructor: PolLn((x,y);(x,y);...).
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

// MarshalJSON implements json.Marshaler: the array of the points.
func (p Polyline[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON implements json.Unmarshaler. The points are decoded into a new slice, leaving one
// the path shared untouched.
func (p *Polyline[T]) UnmarshalJSON(bytes []byte) error {
	var points []Point[T]
	if err := json.Unmarshal(bytes, &points); err != nil {
		return err
	}

	p.Points = points

	return nil
}

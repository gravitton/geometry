package geom_test

import (
	"fmt"
	"iter"
	"math"
	"slices"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

var (
	_ Shape[int]       = Circle[int]{}
	_ Shape[int]       = Segment[int]{}
	_ Shape[int]       = Polygon[int]{}
	_ Shape[int]       = Rectangle[int]{}
	_ Shape[int]       = RegularPolygon[int]{}
	_ Shape[int]       = Box[int]{}
	_ Shape[int]       = Ellipse[int]{}
	_ intersector[int] = Circle[int]{}
	_ intersector[int] = Segment[int]{}
	_ intersector[int] = Ray[int]{}
	_ intersector[int] = Line[int]{}
	_ intersector[int] = Polygon[int]{}
	_ intersector[int] = Rectangle[int]{}
	_ intersector[int] = RegularPolygon[int]{}
	_ intersector[int] = Box[int]{}
	_ Body[int]        = Circle[int]{}
	_ Body[int]        = Polygon[int]{}
	_ Body[int]        = Rectangle[int]{}
	_ Body[int]        = RegularPolygon[int]{}
	_ Body[int]        = Ellipse[int]{}
)

// polyline is a shape whose boundary is a chain of straight edges, each starting where the
// previous one ends, the last closing back to the first on a closed shape.
type polyline[T Number] interface {
	Vertices() iter.Seq[Point[T]]
	Edges() iter.Seq[Segment[T]]
}

// caster is a segment, a ray or a line as a cast: the points it holds and the point a fraction
// of the way along it.
type caster[T Number] interface {
	Contains(point Point[T]) bool
	PointAt(t float64) Point[T]
}

// obstacle is a segment, a ray or a line as the wall a cast is made at: the points it holds.
type obstacle[T Number] interface {
	Contains(point Point[T]) bool
}

// container is a shape other shapes lie within, with the six Contains<Kind> methods: Circle,
// Polygon, Rectangle, RegularPolygon and Box.
type container[T Number] interface {
	ContainsCircle(circle Circle[T]) bool
	ContainsSegment(segment Segment[T]) bool
	ContainsPolygon(polygon Polygon[T]) bool
	ContainsRectangle(rectangle Rectangle[T]) bool
	ContainsRegularPolygon(polygon RegularPolygon[T]) bool
	ContainsBox(box Box[T]) bool
}

// transformable is a shape that moves, turns and scales in its own type, the constraint a
// generic tween would take.
type transformable[T Number, S any] interface {
	Translate(vector Vector[T]) S
	MoveTo(point Point[T]) S
	Rotate(angle float64) S
	Scale(factor float64) S
	Unscale(factor float64) S
}

// bounded is a shape with a circle around it, beside the box Shape gives: every Shape of the
// package, which the interface itself does not list.
type bounded[T Number] interface {
	Shape[T]
	BoundingCircle() Circle[T]
}

// solid is a shape with mass properties, the pair the Body properties are checked through.
type solid[T Number] interface {
	Shape[T]
	Body[T]
}

// intersector is a shape, a ray or a line with the Intersects method for each of the eight
// kinds: Circle, Segment, Ray, Line, Polygon, Rectangle, RegularPolygon and Box. The interface
// of the same name in the package lists the six over shapes, all Intersects reads; the tests
// hold the rays and the lines too.
type intersector[T Number] interface {
	IntersectsCircle(circle Circle[T]) bool
	IntersectsSegment(segment Segment[T]) bool
	IntersectsRay(ray Ray[T]) bool
	IntersectsLine(line Line[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
	IntersectsBox(box Box[T]) bool
}

// intersectsKind tests a against b by the method of a naming the kind of b, the rays and the
// lines included, which Intersects does not take.
func intersectsKind[T Number](a, b intersector[T]) bool {
	switch kind := b.(type) {
	case Circle[T]:
		return a.IntersectsCircle(kind)
	case Segment[T]:
		return a.IntersectsSegment(kind)
	case Ray[T]:
		return a.IntersectsRay(kind)
	case Line[T]:
		return a.IntersectsLine(kind)
	case Polygon[T]:
		return a.IntersectsPolygon(kind)
	case Rectangle[T]:
		return a.IntersectsRectangle(kind)
	case RegularPolygon[T]:
		return a.IntersectsRegularPolygon(kind)
	case Box[T]:
		return a.IntersectsBox(kind)
	default:
		panic(fmt.Sprintf("intersectsKind: no method for %T", b))
	}
}

// containsKind tests the shape against the container by the method naming its kind, as
// Intersects dispatches. A ray and a line run without end, so no container encloses one.
func containsKind[T Number](outer container[T], shape intersector[T]) bool {
	switch kind := shape.(type) {
	case Circle[T]:
		return outer.ContainsCircle(kind)
	case Segment[T]:
		return outer.ContainsSegment(kind)
	case Ray[T], Line[T]:
		return false
	case Polygon[T]:
		return outer.ContainsPolygon(kind)
	case Rectangle[T]:
		return outer.ContainsRectangle(kind)
	case RegularPolygon[T]:
		return outer.ContainsRegularPolygon(kind)
	case Box[T]:
		return outer.ContainsBox(kind)
	default:
		panic(fmt.Sprintf("containsKind: no method for %T", shape))
	}
}

func TestShape(t *testing.T) {
	shapes := []Shape[float64]{
		Circ(Pt(1.0, 1.0), 2.0),
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Pol(triangleVertices()),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop),
		BoxFromMinMax(Pt(-1.0, 0.5), Pt(2.0, 3.0)),
		Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/6),
	}

	t.Run("distance is zero exactly where contains holds", func(t *testing.T) {
		for _, shape := range shapes {
			for _, p := range pointFixtures {
				assert.Equal(t, shape.DistanceTo(p) == 0, shape.Contains(p), fmt.Sprintf("%s → %s: ", shape, p))
				assert.True(t, shape.DistanceSquaredTo(p) >= 0, fmt.Sprintf("%s → %s: ", shape, p))
			}
		}
	})
	t.Run("bounds contain every contained point", func(t *testing.T) {
		for _, shape := range shapes {
			for _, p := range pointFixtures {
				assert.True(t, !shape.Contains(p) || shape.Bounds().Contains(p), fmt.Sprintf("%s → %s: ", shape, p))
			}
		}
	})
	t.Run("the bounding circle contains every contained point", func(t *testing.T) {
		for _, shape := range shapes {
			circle := shape.(bounded[float64]).BoundingCircle()

			for _, p := range pointFixtures {
				assert.True(t, !shape.Contains(p) || circle.Contains(p), fmt.Sprintf("%s → %s: ", shape, p))
			}
		}
	})
	t.Run("far from the origin the float32 bounding circle contains every contained point", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, shape := range farShapes(offset) {
				circle := shape.(bounded[float32]).BoundingCircle()

				for _, p := range pointFixtures {
					point := p.Cast[float32]().Add(offset)

					assert.True(t, !shape.Contains(point) || circle.Contains(point), fmt.Sprintf("%s → %s: ", shape, point))
				}
			}
		}
	})
	t.Run("a NaN coordinate is contained by nothing and at NaN distance", func(t *testing.T) {
		for _, shape := range shapes {
			for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN()), Pt(math.Inf(1), 0.0)} {
				assert.False(t, shape.Contains(p), fmt.Sprintf("%s → %s: ", shape, p))
				assert.True(t, math.IsNaN(shape.DistanceTo(p)) || math.IsInf(shape.DistanceTo(p), 1), fmt.Sprintf("%s → %s: ", shape, p))
			}
		}
	})
	t.Run("the nearest point of a NaN coordinate is NaN", func(t *testing.T) {
		for _, shape := range shapes {
			for _, p := range []Point[float64]{Pt(math.NaN(), 0.0), Pt(1.0, math.NaN())} {
				nearest := shape.Nearest(p)

				assert.True(t, math.IsNaN(nearest.X) || math.IsNaN(nearest.Y), fmt.Sprintf("%s → %s: ", shape, p))
			}
		}
	})
	t.Run("the nearest point is contained and at the distance", func(t *testing.T) {
		for _, shape := range shapes {
			for _, p := range pointFixtures {
				assertNearest(t, shape, p)
			}
		}
	})
	t.Run("far from the origin the float32 nearest point is contained", func(t *testing.T) {
		for _, offset := range farOffsets {
			for _, shape := range farShapes(offset) {
				for _, p := range pointFixtures {
					point := p.Cast[float32]().Add(offset)

					assert.True(t, shape.Contains(shape.Nearest(point)), fmt.Sprintf("%s → %s: ", shape, point))
				}
			}
		}
	})
	t.Run("the support lies on the shape and no contained point reaches farther", func(t *testing.T) {
		for _, shape := range convexFixtures() {
			for _, direction := range vectorFixtures {
				support := shape.Support(direction)
				message := fmt.Sprintf("%v → %s: ", shape, direction)

				assert.True(t, shape.Contains(support), message)
				for _, p := range pointFixtures {
					assert.True(t, !shape.Contains(p) || LessOrEqualDelta(p.Vector().Dot(direction), support.Vector().Dot(direction), Delta*(1+direction.Length())), message+p.String())
				}
			}
		}
	})
}

// assertNearest asserts the properties Nearest has on every shape: the point itself exactly
// where Contains holds, contained itself, its own nearest point, and at DistanceSquaredTo
// from the point.
func assertNearest[T Number](t *testing.T, shape Shape[T], point Point[T]) {
	t.Helper()

	nearest := shape.Nearest(point)
	message := fmt.Sprintf("%v → %v: ", shape, point)

	assert.Equal(t, nearest.Equal(point), shape.Contains(point), message)
	assert.True(t, shape.Contains(nearest), message)
	assert.True(t, shape.Nearest(nearest).Equal(nearest), message)
	geomtest.AssertNumber(t, float64(point.DistanceSquaredTo(nearest)), shape.DistanceSquaredTo(point), message)
}

// assertBoundingCircle asserts that a bounding circle contains every vertex of the shape it
// was asked of, by the comparison Contains makes.
func assertBoundingCircle[T Number](t *testing.T, circle Circle[T], vertices iter.Seq[Point[T]], message string) {
	t.Helper()

	for vertex := range vertices {
		assert.True(t, circle.Contains(vertex), fmt.Sprintf("%s → %v ∌ %v: ", message, circle, vertex))
	}
}

// along returns how far the point lies along the cast from where it starts, the order its
// crossings and its clipped parts follow: the distance from Start or Origin in the direction the
// cast runs, negative behind the origin of a line.
func along[T Number](cast caster[T], point Point[T]) float64 {
	start := cast.PointAt(0).Float()

	return point.Float().Subtract(start).Dot(cast.PointAt(1).Float().Subtract(start).Normalize())
}

// assertCrossings asserts the properties the boundary crossings of a segment, a ray or a line
// have on every outline, given its edges: each point lies on the cast and on an edge by
// Contains, they follow one another along the cast, and there is one only where the pair
// intersects.
func assertCrossings[T Number](t *testing.T, cast caster[T], points []Point[T], intersects bool, edges func() iter.Seq[Segment[T]], message string) {
	t.Helper()

	reached := math.Inf(-1)
	for _, point := range points {
		distance := along(cast, point)

		assert.True(t, cast.Contains(point), message+point.String()+" on the cast: ")
		assert.True(t, slices.ContainsFunc(slices.Collect(edges()), func(edge Segment[T]) bool {
			return edge.Contains(point)
		}), message+point.String()+" on the boundary: ")
		assert.True(t, LessOrEqual(reached, distance), message+point.String()+" in order: ")

		reached = distance
	}
	assert.True(t, len(points) == 0 || intersects, message)
}

// assertClipped asserts the properties the parts a Clip method returns for a segment, a ray or
// a line have on every shape: they exist exactly where the cast intersects the shape, each lies
// on the cast and within the shape, and they follow one another along the cast without
// overlapping.
func assertClipped[T Number](t *testing.T, cast caster[T], parts []Segment[T], intersects bool, inside func(Segment[T]) bool, message string) {
	t.Helper()

	assert.Equal(t, len(parts) > 0, intersects, message)

	reached := math.Inf(-1)
	for _, part := range parts {
		from, to := along(cast, part.Start), along(cast, part.End)

		assert.True(t, cast.Contains(part.Start) && cast.Contains(part.End), message+part.String()+" on the cast: ")
		assert.True(t, inside(part), message+part.String()+" within the shape: ")
		assert.True(t, LessOrEqual(reached, from) && LessOrEqual(from, to), message+part.String()+" in order: ")

		reached = to
	}
}

// assertShapeHit asserts the properties the Hit a method returns for a segment or a ray has on
// every shape with an inside: there is one exactly where the shape does not contain the start
// of the cast and the cast crosses its boundary, at the first of the crossings, facing the cast.
func assertShapeHit[T Number](t *testing.T, cast caster[T], hit Hit[T], ok, contained bool, crossings []Point[T], message string) {
	t.Helper()

	assert.Equal(t, ok, !contained && len(crossings) > 0, message)
	if !ok {
		assert.Equal(t, hit, Hit[T]{}, message)

		return
	}

	assert.True(t, hit.Point.Equal(crossings[0]), message+hit.Point.String()+" is the first crossing: ")
	assertFacing(t, cast, hit, message)
}

// assertWallHit asserts the properties the Hit a method returns for a segment or a ray has
// against a wall, a segment, a ray or a line: there is one exactly where the wall does not
// contain the start of the cast and the two intersect, at a point of both, facing the cast.
func assertWallHit[T Number](t *testing.T, cast caster[T], wall obstacle[T], hit Hit[T], ok, intersects bool, message string) {
	t.Helper()

	assert.Equal(t, ok, !wall.Contains(cast.PointAt(0)) && intersects, message)
	if !ok {
		assert.Equal(t, hit, Hit[T]{}, message)

		return
	}

	assert.True(t, cast.Contains(hit.Point) && wall.Contains(hit.Point), message+hit.Point.String()+" lies on both: ")
	assertFacing(t, cast, hit, message)
}

// assertFacing asserts what every Hit holds of the cast it was found on: a unit normal that
// does not run along the cast, and a fraction ahead of its start that PointAt takes back to the
// point.
func assertFacing[T Number](t *testing.T, cast caster[T], hit Hit[T], message string) {
	t.Helper()

	direction := cast.PointAt(1).Float().Subtract(cast.PointAt(0).Float())

	assert.True(t, hit.Normal.IsNormalized(), message+hit.Normal.String()+" is a unit: ")
	assert.True(t, LessOrEqual(hit.Normal.Dot(direction.Normalize()), 0), message+hit.Normal.String()+" faces the cast: ")
	assert.True(t, hit.Fraction >= 0 && cast.PointAt(hit.Fraction).Equal(hit.Point), message+"PointAt of the fraction: ")
}

// assertReversed asserts that a segment and its Reverse cross a boundary at the same points,
// given the crossings of each: those of the one are those of the other in the opposite order.
func assertReversed[T Number](t *testing.T, points, reversed []Point[T], message string) {
	t.Helper()

	ordered := slices.Clone(reversed)
	slices.Reverse(ordered)

	geomtest.AssertVertices(t, ordered, points, message)
}

func TestPolyline(t *testing.T) {
	polylines := []polyline[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Pol(triangleVertices()),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop),
	}

	t.Run("every edge starts at a vertex and ends at the next", func(t *testing.T) {
		for _, shape := range polylines {
			vertices, edges := slices.Collect(shape.Vertices()), slices.Collect(shape.Edges())

			for i, edge := range edges {
				assert.True(t, edge.Start.Equal(vertices[i]), fmt.Sprintf("%s #%d: ", shape, i))
				assert.True(t, edge.End.Equal(vertices[(i+1)%len(vertices)]), fmt.Sprintf("%s #%d: ", shape, i))
			}
		}
	})
	t.Run("only the segment is open", func(t *testing.T) {
		assert.False(t, closed[float64](Seg(Pt(0.0, 0.0), Pt(3.0, 4.0))))
		assert.True(t, closed[float64](Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi/6)))
		assert.True(t, closed[float64](Pol(triangleVertices())))
		assert.True(t, closed[float64](RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop)))
	})
}

func TestIntersector(t *testing.T) {
	shapes := []intersector[float64]{
		Circ(Pt(1.0, 1.0), 2.0),
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		RayThrough(Pt(-1.0, -1.0), Pt(3.0, 4.0)),
		LineThrough(Pt(-1.0, -1.0), Pt(3.0, 4.0)),
		Pol(triangleVertices()),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop),
		BoxFromMinMax(Pt(-1.0, 0.5), Pt(2.0, 3.0)),
	}
	for _, c := range circleFixtures {
		shapes = append(shapes, c)
	}
	for _, s := range segmentFixtures {
		shapes = append(shapes, s)
	}
	for _, r := range rayFixtures {
		shapes = append(shapes, r)
	}
	for _, l := range lineFixtures {
		shapes = append(shapes, l)
	}
	for _, p := range polygonFixtures() {
		shapes = append(shapes, p)
	}
	for _, r := range rectFixtures {
		shapes = append(shapes, r)
	}
	for _, b := range boxFixtures {
		shapes = append(shapes, b)
	}

	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, a := range shapes {
			for _, b := range shapes {
				assert.Equal(t, intersectsKind(a, b), intersectsKind(b, a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("a shape enclosed by another intersects it", func(t *testing.T) {
		for _, a := range shapes {
			if outer, ok := a.(container[float64]); ok {
				for _, b := range shapes {
					assert.True(t, !containsKind(outer, b) || intersectsKind(a, b), fmt.Sprintf("%s → %s: ", a, b))
				}
			}
		}
	})
	t.Run("every shape with an area encloses itself", func(t *testing.T) {
		for _, a := range shapes {
			if outer, ok := a.(container[float64]); ok {
				assert.True(t, containsKind(outer, a), fmt.Sprintf("%s: ", a))
			}
		}
	})
	t.Run("far from the origin float32 shapes answer alike from either side and enclose themselves", func(t *testing.T) {
		for _, offset := range farOffsets {
			far := farIntersectors(offset)

			for _, a := range far {
				for _, b := range far {
					assert.Equal(t, intersectsKind(a, b), intersectsKind(b, a), fmt.Sprintf("%s → %s: ", a, b))
				}
				if outer, ok := a.(container[float32]); ok {
					assert.True(t, containsKind(outer, a), fmt.Sprintf("%s: ", a))
					for _, b := range far {
						assert.True(t, !containsKind(outer, b) || intersectsKind(a, b), fmt.Sprintf("%s → %s: ", a, b))
					}
				}
			}
		}
	})
}

// foreignShape is a shape of no kind this package names, as one of another package is: the
// shape it holds, known to the package by the methods of Shape alone.
type foreignShape[T Number] struct {
	Shape[T]
}

// convexFixtures returns the fixtures of every kind of shape that are convex: the polygons that
// are, and no empty one.
func convexFixtures() []Shape[float64] {
	var shapes []Shape[float64]
	for _, c := range circleFixtures {
		shapes = append(shapes, c)
	}
	for _, s := range segmentFixtures {
		shapes = append(shapes, s)
	}
	for _, p := range polygonFixtures() {
		if p.IsConvex() {
			shapes = append(shapes, p)
		}
	}
	for _, r := range rectFixtures {
		shapes = append(shapes, r)
	}
	for _, rp := range regularPolygonFixtures {
		shapes = append(shapes, rp)
	}
	for _, b := range boxFixtures {
		shapes = append(shapes, b)
	}
	for _, e := range ellipseFixtures {
		shapes = append(shapes, e)
	}

	return shapes
}

// shapeFixtures returns the fixtures of every kind of shape: the convex ones, the polygon
// fixtures that are not convex, and the notched polygon, the one with an outside between its
// vertices.
func shapeFixtures() []Shape[float64] {
	shapes := convexFixtures()
	for _, p := range polygonFixtures() {
		if !p.IsConvex() {
			shapes = append(shapes, p)
		}
	}

	return append(shapes, Pol(notchedVertices()).Float())
}

// translated returns the convex shape moved by the vector, in its own kind.
func translated(shape Shape[float64], vector Vector[float64]) Shape[float64] {
	switch kind := shape.(type) {
	case Circle[float64]:
		return kind.Translate(vector)
	case Segment[float64]:
		return kind.Translate(vector)
	case Polygon[float64]:
		return kind.Translate(vector)
	case Rectangle[float64]:
		return kind.Translate(vector)
	case RegularPolygon[float64]:
		return kind.Translate(vector)
	case Box[float64]:
		return kind.Translate(vector)
	case Ellipse[float64]:
		return kind.Translate(vector)
	default:
		panic(fmt.Sprintf("translated: no method for %T", shape))
	}
}

// farOffsets move the float32 fixtures to where an ulp of float32 exceeds Delta32, so a point
// rounded into T stays on the boundary it was computed on only by the tolerance epsilonAt
// widens with the coordinates.
var farOffsets = []Vector[float32]{Vec[float32](1e4, -5e3), Vec[float32](1e5, -5e4)}

// farShapes returns the shape fixtures of every kind cast to float32 and moved by the offset.
func farShapes(offset Vector[float32]) []Shape[float32] {
	var shapes []Shape[float32]
	for _, c := range circleFixtures {
		shapes = append(shapes, c.Cast[float32]().Translate(offset))
	}
	for _, s := range segmentFixtures {
		shapes = append(shapes, s.Cast[float32]().Translate(offset))
	}
	for _, p := range polygonFixtures() {
		shapes = append(shapes, p.Cast[float32]().Translate(offset))
	}
	for _, r := range rectFixtures {
		shapes = append(shapes, r.Cast[float32]().Translate(offset))
	}
	for _, rp := range regularPolygonFixtures {
		shapes = append(shapes, rp.Cast[float32]().Translate(offset))
	}
	for _, b := range boxFixtures {
		shapes = append(shapes, b.Cast[float32]().Translate(offset))
	}
	for _, e := range ellipseFixtures {
		shapes = append(shapes, e.Cast[float32]().Translate(offset))
	}

	return shapes
}

// farIntersectors returns the fixtures of every kind with the Intersects methods, cast to float32
// and moved by the offset: the far shapes but the ellipse, the rays and the lines.
func farIntersectors(offset Vector[float32]) []intersector[float32] {
	var shapes []intersector[float32]
	for _, shape := range farShapes(offset) {
		if kind, ok := shape.(intersector[float32]); ok {
			shapes = append(shapes, kind)
		}
	}
	for _, r := range rayFixtures {
		shapes = append(shapes, r.Cast[float32]().Translate(offset))
	}
	for _, l := range lineFixtures {
		shapes = append(shapes, l.Cast[float32]().Translate(offset))
	}

	return shapes
}

// narrowMatrix stretches the fixtures before they are cast to int16, past the square root of its
// range, so a squared distance or a difference taken in T rather than float64 overflows where
// every coordinate still fits.
var narrowMatrix = ScaleMatrix(40.0, 40.0)

// prefixPoint is the point the Append methods find in dst over the fixtures, so every result is
// checked to follow it.
var prefixPoint = Pt(-7.5, 3.25)

// prefixSegment is the segment the Append methods of segments find in dst over the fixtures.
var prefixSegment = Seg(Pt(-7.5, 3.25), Pt(4.5, -1.75))

// bufferWith returns a slice holding the element with room for more, the dst the Append methods
// are checked with over the fixtures, so their results are appended in place after it.
func bufferWith[E any](element E) []E {
	return append(make([]E, 0, 8), element)
}

// closed reports whether the last edge returns to the first vertex, through the interface as
// a caller would.
func closed[T Number](shape polyline[T]) bool {
	var first, last Point[T]
	for vertex := range shape.Vertices() {
		first = vertex

		break
	}
	for edge := range shape.Edges() {
		last = edge.End
	}

	return first.Equal(last)
}

func TestBody(t *testing.T) {
	solids := []solid[float64]{
		Circ(Pt(1.0, 1.0), 2.0),
		Pol(triangleVertices()),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		RegularPolygonWithOrientation(Pt(0.0, 0.0), SzU(10.0), 6, OrientationFlatTop),
		Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/6),
	}

	t.Run("the centroid lies within the shape", func(t *testing.T) {
		for _, s := range solids {
			assert.True(t, s.Contains(s.Centroid()), fmt.Sprintf("%v: ", s))
		}
	})
	t.Run("the inertia is positive and below the area at the farthest corner", func(t *testing.T) {
		for _, s := range solids {
			bounds := s.Bounds()
			reach := max(bounds.Min.DistanceSquaredTo(s.Centroid()), bounds.Max.DistanceSquaredTo(s.Centroid()))

			assert.True(t, s.Inertia() > 0, fmt.Sprintf("%v: ", s))
			assert.True(t, s.Inertia() < s.Area()*reach, fmt.Sprintf("%v: ", s))
		}
	})
}

func TestIntersects(t *testing.T) {
	t.Run("a circle is tested against any shape by its radius", func(t *testing.T) {
		shape := foreignShape[float64]{Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))}

		assert.True(t, Intersects[float64](shape, Circ(Pt(2.0, 0.0), 1.0)))
		assert.True(t, Intersects[float64](Circ(Pt(2.0, 0.0), 1.0), shape))
		assert.False(t, Intersects[float64](Circ(Pt(2.5, 0.0), 1.0), shape))
	})
	t.Run("convex shapes known by their support answer as the named methods do", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				if _, ok := a.(Ellipse[float64]); ok {
					continue
				}
				if _, ok := b.(Ellipse[float64]); ok {
					continue
				}

				message := fmt.Sprintf("%v → %v: ", a, b)

				assert.Equal(t, Intersects[float64](foreignShape[float64]{a}, foreignShape[float64]{b}), Intersects[float64](a, b), message)
				assert.Equal(t, Intersects[float64](foreignShape[float64]{a}, b), Intersects[float64](a, b), message)
			}
		}
	})
	t.Run("an ellipse touching a shape intersects it, and not a tolerance apart", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		assert.True(t, Intersects[float64](ellipse, Ell(Pt(6.0, 0.0), Sz(2.0, 3.0), 0)))
		assert.False(t, Intersects[float64](ellipse, Ell(Pt(6.0+1e-4, 0.0), Sz(2.0, 3.0), 0)))
		assert.True(t, Intersects[float64](ellipse, BoxFromMinMax(Pt(4.0, -1.0), Pt(6.0, 1.0))))
		assert.False(t, Intersects[float64](ellipse, BoxFromMinMax(Pt(4.0+1e-4, -1.0), Pt(6.0, 1.0))))
		assert.True(t, Intersects[float64](ellipse, Seg(Pt(-5.0, 2.0), Pt(5.0, 2.0))))
		assert.True(t, Intersects[float64](ellipse, Circ(Pt(0.0, 3.0), 1.0)))
		assert.False(t, Intersects[float64](ellipse, Circ(Pt(0.0, 3.0+1e-4), 1.0)))
	})
	t.Run("an int ellipse touches an outline on the lattice", func(t *testing.T) {
		ellipse := Ell(Pt(0, 0), Sz(4, 2), 0)

		assert.True(t, Intersects[int](ellipse, BoxFromMinMax(Pt(4, -1), Pt(6, 1))))
		assert.False(t, Intersects[int](ellipse, BoxFromMinMax(Pt(5, -1), Pt(6, 1))))
		assert.True(t, Intersects[int](ellipse, Ell(Pt(1, 1), Sz(1, 1), 0)))
	})
	t.Run("an ellipse lies between the polygon within it and the polygon about it", func(t *testing.T) {
		shapes := append(shapeFixtures(), Polygon[float64]{}, RegularPolygon[float64]{})

		for _, ellipse := range ellipseFixtures {
			inner := ellipse.RegularPolygon(64, OrientationPointyTop)
			outer := ellipse.Scale(1/math.Cos(Pi/64)).RegularPolygon(64, OrientationPointyTop)

			for _, other := range shapes {
				message := fmt.Sprintf("%v → %v: ", ellipse, other)
				intersects := Intersects[float64](ellipse, other)

				assert.Equal(t, Intersects[float64](other, ellipse), intersects, message)
				if !ellipse.Size.IsPositive() {
					continue
				}

				assert.True(t, !Intersects[float64](inner, other) || intersects, message+"within: ")
				assert.True(t, !intersects || Intersects[float64](outer, other), message+"about: ")
			}
		}
	})
	t.Run("a notch of a polygon is not filled in", func(t *testing.T) {
		polygon := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(6.0, 0.0), Pt(6.0, 6.0), Pt(4.0, 6.0), Pt(4.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 6.0), Pt(0.0, 6.0)})

		assert.False(t, Intersects[float64](polygon, Ell(Pt(3.0, 5.0), Sz(0.5, 1.0), 0)))
		assert.True(t, Intersects[float64](polygon, Ell(Pt(3.0, 5.0), Sz(1.5, 1.0), 0)))
		assert.True(t, Intersects[float64](Ell(Pt(1.0, 1.0), Sz(0.5, 0.25), 1.0), polygon))
		assert.True(t, Intersects[float64](Ell(Pt(3.0, 3.0), Sz(10.0, 10.0), 0), polygon))
	})
	t.Run("an empty polygon intersects no ellipse", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		assert.False(t, Intersects[float64](ellipse, Polygon[float64]{}))
		assert.False(t, Intersects[float64](RegularPolygon[float64]{}, ellipse))
	})
	t.Run("a NaN coordinate intersects nothing", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		assert.False(t, Intersects[float64](ellipse, Ell(Pt(math.NaN(), 0.0), Sz(4.0, 2.0), 0)))
		assert.False(t, Intersects[float64](ellipse, Seg(Pt(math.NaN(), 0.0), Pt(1.0, 1.0))))
	})
	t.Run("far from the origin float32 shapes answer alike from either side", func(t *testing.T) {
		for _, offset := range farOffsets {
			far := farShapes(offset)

			for _, a := range far {
				for _, b := range far {
					assert.Equal(t, Intersects(a, b), Intersects(b, a), fmt.Sprintf("%v → %v: ", a, b))
				}
			}
		}
	})
	t.Run("int shapes far from the origin answer as they do near it", func(t *testing.T) {
		far := Vec(1<<25, -(1 << 25))
		ellipse := Ell(Pt(0, 0), Sz(4, 2), 0).Translate(far)

		assert.True(t, Intersects[int](ellipse, BoxFromMinMax(Pt(4, -1), Pt(6, 1)).Translate(far)))
		assert.False(t, Intersects[int](ellipse, BoxFromMinMax(Pt(5, -1), Pt(6, 1)).Translate(far)))
	})
	t.Run("allocates nothing for shapes held as Shape", func(t *testing.T) {
		var rectangle, other, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Rect(Pt(1.0, 0.0), Sz(3.0, 2.0)), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = Intersects(rectangle, other)
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = Intersects(rectangle, ellipse)
		}), 0)
	})
}

func BenchmarkIntersects(b *testing.B) {
	var rectangle, other, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Rect(Pt(1.0, 0.0), Sz(3.0, 2.0)), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)

	b.Run("by the method naming the kind", func(b *testing.B) {
		for b.Loop() {
			sinkBool = Intersects(rectangle, other)
		}
	})
	b.Run("by supports", func(b *testing.B) {
		for b.Loop() {
			sinkBool = Intersects(rectangle, ellipse)
		}
	})
}

func TestDistance(t *testing.T) {
	t.Run("the gap between two shapes apart", func(t *testing.T) {
		geomtest.AssertNumber(t, Distance[float64](BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)), BoxFromMinMax(Pt(4.0, 5.0), Pt(6.0, 6.0))), 5.0)
		geomtest.AssertNumber(t, Distance[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(5.0, 0.0), 1.5)), 2.5)
		geomtest.AssertNumber(t, Distance[float64](Circ(Pt(0.0, 3.0), 1.0), Rect(Pt(0.0, 0.0), Sz(4.0, 2.0))), 1.0)
		geomtest.AssertNumber(t, Distance[float64](Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0), Ell(Pt(9.0, 0.0), Sz(2.0, 3.0), 0)), 3.0)
		geomtest.AssertNumber(t, Distance[float64](Seg(Pt(0.0, 0.0), Pt(4.0, 0.0)), Seg(Pt(0.0, 2.0), Pt(4.0, 3.0))), 2.0)
	})
	t.Run("int measures in float64", func(t *testing.T) {
		geomtest.AssertNumber(t, Distance[int](BoxFromMinMax(Pt(0, 0), Pt(1, 1)), BoxFromMinMax(Pt(2, 2), Pt(3, 3))), math.Sqrt2)
	})
	t.Run("two boxes are as far apart as the gap IntersectsBox reads", func(t *testing.T) {
		for _, a := range boxFixtures {
			for _, b := range boxFixtures {
				message := fmt.Sprintf("%v → %v: ", a, b)
				distance := DistanceSquared[float64](a, b)

				assert.Equal(t, distance == 0, a.IntersectsBox(b), message)
				assert.True(t, EqualDelta(distance, DistanceSquared[float64](a.Rectangle(), b.Rectangle()), 1e-9*(1+distance)), message)
			}
		}
	})
	t.Run("two circles are as far apart as the circle says", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.Equal(t, DistanceSquared[float64](a, b), a.DistanceSquaredToCircle(b), fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
	t.Run("two segments are as far apart as the segment says", func(t *testing.T) {
		for _, a := range segmentFixtures {
			for _, b := range segmentFixtures {
				assert.Equal(t, DistanceSquared[float64](a, b), a.DistanceSquaredToSegment(b), fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
	t.Run("a shape in a notch is measured to the notch", func(t *testing.T) {
		polygon := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(6.0, 0.0), Pt(6.0, 6.0), Pt(4.0, 6.0), Pt(4.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 6.0), Pt(0.0, 6.0)})

		geomtest.AssertNumber(t, Distance[float64](polygon, Ell(Pt(3.0, 5.0), Sz(0.5, 1.0), 0)), 0.5)
		geomtest.AssertNumber(t, Distance[float64](BoxFromMinMax(Pt(2.5, 4.0), Pt(3.5, 5.0)), polygon), 0.5)
		geomtest.AssertNumber(t, Distance[float64](polygon, Pol([]Point[float64]{Pt(2.5, 4.0), Pt(3.5, 4.0), Pt(3.0, 3.0)})), 0.5)
	})
	t.Run("an empty polygon is infinitely far", func(t *testing.T) {
		box := BoxFromMinMax(Pt(4.0, 2.0), Pt(6.0, 3.0))

		assert.True(t, math.IsInf(Distance[float64](Polygon[float64]{}, box), 1))
		assert.True(t, math.IsInf(Distance[float64](box, RegularPolygon[float64]{}), 1))
		assert.True(t, math.IsInf(Distance[float64](Circ(Pt(0.0, 0.0), 1.0), RegularPolygon[float64]{}), 1))
		assert.True(t, math.IsInf(Distance[float64](RegPol(Pt(0.0, 0.0), Sz(2.0, 2.0), 4, 0, 0), RegularPolygon[float64]{}), 1))
		assert.True(t, math.IsInf(Distance[float64](Ell(Pt(0.0, 0.0), Sz(2.0, 1.0), 0), RegularPolygon[float64]{}), 1))
	})
	t.Run("zero exactly where the shapes intersect, and the same from either side", func(t *testing.T) {
		shapes := shapeFixtures()

		for _, a := range shapes {
			for _, b := range shapes {
				message := fmt.Sprintf("%v → %v: ", a, b)
				distance := Distance(a, b)

				assert.Equal(t, distance == 0, Intersects(a, b), message)
				assert.True(t, EqualDelta(distance, Distance(b, a), 1e-5*(1+distance)), message)
			}
		}
	})
	t.Run("a point as a segment is at the distance the shape gives it", func(t *testing.T) {
		for _, a := range shapeFixtures() {
			for _, p := range pointFixtures {
				expected := a.DistanceTo(p)
				if _, ok := a.(RegularPolygon[float64]); ok && math.IsInf(expected, 1) {
					continue
				}

				distance := Distance[float64](a, Seg(p, p))

				assert.True(t, EqualDelta(distance, expected, 1e-5*(1+expected)), fmt.Sprintf("%v → %s: %v, %v", a, p, distance, expected))
			}
		}
	})
	t.Run("shapes known by their support are as far apart as the kinds they hold", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				expected := Distance[float64](a, b)
				distance := Distance[float64](foreignShape[float64]{a}, foreignShape[float64]{b})

				assert.True(t, EqualDelta(distance, expected, 1e-5*(1+expected)), fmt.Sprintf("%v → %v: %v, %v", a, b, distance, expected))
			}
		}
	})
	t.Run("far from the origin float32 shapes are at no distance exactly where they intersect", func(t *testing.T) {
		for _, offset := range farOffsets {
			far := farShapes(offset)

			for _, a := range far {
				for _, b := range far {
					assert.Equal(t, Distance(a, b) == 0, Intersects(a, b), fmt.Sprintf("%v → %v: ", a, b))
				}
			}
		}
	})
	t.Run("a NaN coordinate is at no distance from nothing", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		assert.False(t, Distance[float64](ellipse, Ell(Pt(math.NaN(), 0.0), Sz(4.0, 2.0), 0)) == 0)
		assert.False(t, Distance[float64](Circ(Pt(math.NaN(), 0.0), 1.0), ellipse) == 0)
	})
	t.Run("allocates nothing for shapes held as Shape", func(t *testing.T) {
		var rectangle, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Ell(Pt(9.0, 0.0), Sz(2.0, 1.0), 0.2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			sinkBool = Distance(rectangle, ellipse) > 0
		}), 0)
	})
}

func TestDistanceSquared(t *testing.T) {
	t.Run("the square of the distance", func(t *testing.T) {
		geomtest.AssertNumber(t, DistanceSquared[int](BoxFromMinMax(Pt(0, 0), Pt(1, 1)), BoxFromMinMax(Pt(2, 2), Pt(3, 3))), 2.0)
		geomtest.AssertNumber(t, DistanceSquared[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(5.0, 0.0), 1.5)), 6.25)
		geomtest.AssertNumber(t, DistanceSquared[float64](Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0), Ell(Pt(9.0, 0.0), Sz(2.0, 3.0), 0)), 9.0)
	})
	t.Run("zero for shapes that intersect", func(t *testing.T) {
		geomtest.AssertNumber(t, DistanceSquared[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(1.5, 0.0), 1.5)), 0.0)
	})
	t.Run("over the fixtures it is the square of Distance", func(t *testing.T) {
		for _, a := range shapeFixtures() {
			for _, b := range shapeFixtures() {
				distance := Distance(a, b)

				geomtest.AssertNumber(t, DistanceSquared(a, b), distance*distance, fmt.Sprintf("%v → %v: ", a, b))
			}
		}
	})
}

func BenchmarkDistanceSquared(b *testing.B) {
	var box, other, ellipse Shape[float64] = BoxFromMinMax(Pt(0.0, 0.0), Pt(4.0, 2.0)), BoxFromMinMax(Pt(7.0, 5.0), Pt(9.0, 6.0)), Ell(Pt(9.0, 0.0), Sz(2.0, 1.0), 0.2)

	b.Run("by the method naming the kind", func(b *testing.B) {
		for b.Loop() {
			sinkBool = DistanceSquared(box, other) > 0
		}
	})
	b.Run("by supports", func(b *testing.B) {
		for b.Loop() {
			sinkBool = DistanceSquared(box, ellipse) > 0
		}
	})
}

func TestCollide(t *testing.T) {
	ground := BoxFromMinMax(Pt(-10.0, 1.5), Pt(10.0, 5.0))

	t.Run("flat sides meet at the two ends of the stretch they share", func(t *testing.T) {
		contact, ok := Collide[float64](BoxFromMinMax(Pt(0.0, 0.0), Pt(4.0, 2.0)), ground)

		assert.True(t, ok)
		geomtest.AssertVector(t, contact.Normal, Vec(0.0, -1.0))
		geomtest.AssertNumber(t, contact.Depth, 0.5)
		assertContactPoints(t, contact, Pt(0.0, 1.5), Pt(4.0, 1.5))
	})
	t.Run("a side a little tilted still meets at two points", func(t *testing.T) {
		contact, ok := Collide[float64](Rect(Pt(2.0, 1.0), Sz(4.0, 2.0)).Rotate(2e-4), ground)

		assert.True(t, ok)
		assert.Equal(t, contact.Points.Count, 2)
	})
	t.Run("a corner meets at one point", func(t *testing.T) {
		contact, ok := Collide[float64](Rect(Pt(2.0, 0.5), Sz(2.0, 2.0)).Rotate(Pi/4), ground)

		assert.True(t, ok)
		geomtest.AssertVector(t, contact.Normal, Vec(0.0, -1.0))
		geomtest.AssertNumber(t, contact.Depth, math.Sqrt2-1)
		assertContactPoints(t, contact, Pt(2.0, 1.5))
	})
	t.Run("a round side meets at one point", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(2.0, 1.0), 1.0), ground)

		assert.True(t, ok)
		assert.True(t, contact.Normal.Subtract(Vec(0.0, -1.0)).LessOrEqual(1e-3), fmt.Sprint(contact))
		assert.True(t, EqualDelta(contact.Depth, 0.5, 1e-5), fmt.Sprint(contact))
		assert.Equal(t, contact.Points.Count, 1)
		assert.True(t, contact.Points.Items[0].DistanceTo(Pt(2.0, 1.5)) < 1e-2, fmt.Sprint(contact))

		contact, ok = Collide[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(1.5, 0.0), 1.0))

		assert.True(t, ok)
		assert.Equal(t, contact.Points.Count, 1)
		assert.True(t, contact.Points.Items[0].DistanceTo(Pt(0.5, 0.0)) < 1e-2, fmt.Sprint(contact))
	})
	t.Run("shapes that only touch have a normal and no depth", func(t *testing.T) {
		contact, ok := Collide[float64](BoxFromMinMax(Pt(0.0, 0.0), Pt(4.0, 1.5)), ground)

		assert.True(t, ok)
		geomtest.AssertVector(t, contact.Normal, Vec(0.0, -1.0))
		geomtest.AssertNumber(t, contact.Depth, 0.0)
		assertContactPoints(t, contact, Pt(0.0, 1.5), Pt(4.0, 1.5))
	})
	t.Run("int rounds the points", func(t *testing.T) {
		contact, ok := Collide[int](BoxFromMinMax(Pt(0, 0), Pt(4, 2)), BoxFromMinMax(Pt(-10, 1), Pt(10, 5)))

		assert.True(t, ok)
		geomtest.AssertVector(t, contact.Normal, Vec(0.0, -1.0))
		geomtest.AssertNumber(t, contact.Depth, 1.0)
		assertContactPoints(t, contact, Pt(0, 1), Pt(4, 1))
	})
	t.Run("two points have no normal", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(1.0, 1.0), 0.0), Circ(Pt(1.0, 1.0), 0.0))

		assert.True(t, ok)
		assert.True(t, contact.Normal.IsZero())
		assertContactPoints(t, contact, Pt(1.0, 1.0))
	})
	t.Run("shapes apart have none", func(t *testing.T) {
		contact, ok := Collide[float64](Circ(Pt(0.0, 0.0), 1.0), Circ(Pt(3.0, 0.0), 1.0))

		assert.False(t, ok)
		assert.Equal(t, contact, Contact[float64]{})
	})
	t.Run("the normal is a unit vector and the points lie on the second shape", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				message := fmt.Sprintf("%v → %v: ", a, b)

				contact, ok := Collide(a, b)

				assert.Equal(t, ok, Intersects[float64](a, b), message)
				if !ok {
					continue
				}

				assert.True(t, contact.Normal.IsNormalized() || contact.Normal.IsZero(), message)
				assert.True(t, contact.Points.Count >= 1, message)
				for point := range contact.Points.All() {
					assert.True(t, b.Contains(point), message+point.String())
					assert.True(t, a.DistanceTo(point) <= contact.Depth+1e-3, message+point.String())
				}
			}
		}
	})
	t.Run("given the other way round the normal is the opposite, or the same where both ways out are as short", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				contact, _ := Collide(a, b)
				reverse, _ := Collide(b, a)
				message := fmt.Sprintf("%v → %v: %v, %v", a, b, contact.Normal, reverse.Normal)

				geomtest.AssertNumber(t, contact.Depth, reverse.Depth, message)
				assert.True(t, contact.Normal == reverse.Normal.Negate() || contact.Normal == reverse.Normal, message)
			}
		}
	})
	t.Run("boxes overlapping alike on both axes part along X whichever is given first", func(t *testing.T) {
		a, b := BoxFromMinMax(Pt(0.0, 0.0), Pt(2.0, 2.0)), BoxFromMinMax(Pt(1.0, 1.0), Pt(3.0, 3.0))

		contact, _ := Collide[float64](a, b)
		reverse, _ := Collide[float64](b, a)

		geomtest.AssertVector(t, contact.Normal, Vec(-1.0, 0.0))
		geomtest.AssertVector(t, reverse.Normal, Vec(1.0, 0.0))
	})
	t.Run("far from the origin float32 shapes meet where they intersect, at points of the second", func(t *testing.T) {
		for _, offset := range farOffsets {
			far := farShapes(offset)

			for _, a := range far {
				for _, b := range far {
					message := fmt.Sprintf("%v → %v: ", a, b)

					contact, ok := Collide(a, b)
					assert.Equal(t, ok, Intersects(a, b), message)
					assert.True(t, !ok || contact.Depth >= 0, message)
					for point := range contact.Points.All() {
						assert.True(t, b.DistanceTo(point) < 0.5, message+point.String())
					}
				}
			}
		}
	})
	t.Run("int shapes far from the origin meet as they do near it", func(t *testing.T) {
		far := Vec(1<<25, -(1 << 25))

		contact, ok := Collide[int](BoxFromMinMax(Pt(0, 0), Pt(4, 2)).Translate(far), BoxFromMinMax(Pt(-10, 1), Pt(10, 5)).Translate(far))

		assert.True(t, ok)
		geomtest.AssertVector(t, contact.Normal, Vec(0.0, -1.0))
		geomtest.AssertNumber(t, contact.Depth, 1.0)
		assertContactPoints(t, contact, Pt(0, 1).Add(far), Pt(4, 1).Add(far))
	})
	t.Run("allocates nothing for shapes held as Shape", func(t *testing.T) {
		var rectangle, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = Collide(rectangle, ellipse)
		}), 0)
	})
}

// assertContactPoints asserts that a contact holds the given points, in either order.
func assertContactPoints[T Number](t *testing.T, contact Contact[T], expected ...Point[T]) {
	t.Helper()

	points := slices.Clone(contact.Points.Items[:contact.Points.Count])
	slices.SortFunc(points, Point[T].Compare)
	slices.SortFunc(expected, Point[T].Compare)

	geomtest.AssertVertices(t, points, expected)
}

func BenchmarkCollide(b *testing.B) {
	var rectangle, other, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Rect(Pt(1.0, 0.0), Sz(3.0, 2.0)), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)

	b.Run("two outlines", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Collide(rectangle, other)
		}
	})
	b.Run("a round shape", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Collide(rectangle, ellipse)
		}
	})
}

func FuzzCollide(f *testing.F) {
	f.Add(0.0, 0.0, 4.0, 2.0, 0.3, 3.0, 0.0, 2.0, 1.0, 0.2)
	f.Add(0.0, 0.0, 4.0, 2.0, 0.0, 6.0, 0.0, 2.0, 3.0, 0.0)
	f.Add(0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0)
	f.Add(1.0, 2.0, 0.0, 3.0, 1.0, 1.5, 2.0, 2.0, 0.0, 2.0)

	f.Fuzz(func(t *testing.T, x1, y1, w1, h1, angle1, x2, y2, w2, h2, angle2 float64) {
		for _, v := range []float64{x1, y1, w1, h1, angle1, x2, y2, w2, h2, angle2} {
			if math.IsNaN(v) || math.Abs(v) > 1e3 {
				t.Skip()
			}
		}

		var a, b Shape[float64] = Ell(Pt(x1, y1), Sz(w1, h1), angle1), Rect(Pt(x2, y2), Sz(w2, h2)).Rotate(angle2)
		message := fmt.Sprintf("%v → %v: ", a, b)

		intersects := Intersects(a, b)
		contact, ok := Collide(a, b)

		assert.Equal(t, Intersects(b, a), intersects, message)
		assert.Equal(t, ok, intersects, message)
		assert.Equal(t, Distance(a, b) == 0, intersects, message)
		assert.True(t, contact.Depth >= 0 && contact.Points.Count <= 2, message)
		assert.True(t, !ok || contact.Points.Count >= 1, message)
	})
}

func TestSweep(t *testing.T) {
	t.Run("a box stops at the side of a box", func(t *testing.T) {
		hit, ok := Sweep[float64](BoxFromMinMax(Pt(0.0, 0.0), Pt(1.0, 1.0)), Vec(10.0, 0.0), BoxFromMinMax(Pt(5.0, -1.0), Pt(6.0, 3.0)))

		assert.True(t, ok)
		geomtest.AssertNumber(t, hit.Fraction, 0.4)
		geomtest.AssertVector(t, hit.Normal, Vec(-1.0, 0.0))
		geomtest.AssertNumber(t, hit.Point.X, 5.0)
	})
	t.Run("a circle stops against a circle", func(t *testing.T) {
		hit, ok := Sweep[float64](Circ(Pt(0.0, 0.0), 1.0), Vec(10.0, 0.0), Circ(Pt(5.0, 1.0), 1.0))

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, (5-math.Sqrt(3))/10, 1e-5), fmt.Sprint(hit))
		assert.True(t, EqualDelta(hit.Normal.X, -math.Sqrt(3)/2, 1e-3) && EqualDelta(hit.Normal.Y, -0.5, 1e-3), fmt.Sprint(hit))
		assert.True(t, EqualDelta(hit.Point.X, 5-math.Sqrt(3)/2, 1e-3) && EqualDelta(hit.Point.Y, 0.5, 1e-3), fmt.Sprint(hit))
	})
	t.Run("int rounds the point", func(t *testing.T) {
		hit, ok := Sweep[int](BoxFromMinMax(Pt(0, 0), Pt(1, 1)), Vec(10, 0), BoxFromMinMax(Pt(5, -1), Pt(6, 3)))

		assert.True(t, ok)
		assert.Equal(t, hit.Point.X, 5)
		geomtest.AssertNumber(t, hit.Fraction, 0.4)
	})
	t.Run("a motion that reaches nothing hits nothing", func(t *testing.T) {
		circle, wall := Circ(Pt(0.0, 0.0), 1.0), BoxFromMinMax(Pt(5.0, -1.0), Pt(6.0, 3.0))

		for _, motion := range []Vector[float64]{Vec(3.0, 0.0), Vec(-10.0, 0.0), Vec(0.0, 10.0), Vec(0.0, 0.0), Vec(10.0, 40.0), Vec(math.NaN(), 0.0)} {
			_, ok := Sweep[float64](circle, motion, wall)

			assert.False(t, ok, motion.String())
		}
	})
	t.Run("a motion ending on the shape hits it at the end", func(t *testing.T) {
		hit, ok := Sweep[float64](Circ(Pt(0.0, 0.0), 1.0), Vec(4.0, 0.0), BoxFromMinMax(Pt(5.0, -1.0), Pt(6.0, 3.0)))

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, 1, 1e-5), fmt.Sprint(hit))
	})
	t.Run("shapes that intersect before the move hit nothing", func(t *testing.T) {
		_, ok := Sweep[float64](Circ(Pt(0.0, 0.0), 1.0), Vec(10.0, 0.0), Circ(Pt(1.5, 0.0), 1.0))

		assert.False(t, ok)
	})
	t.Run("a shape moves into the notch of a polygon and stops at the edge it reaches there", func(t *testing.T) {
		polygon := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(6.0, 0.0), Pt(6.0, 6.0), Pt(4.0, 6.0), Pt(4.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 6.0), Pt(0.0, 6.0)})
		circle := Circ(Pt(3.0, 8.0), 0.5)

		hit, ok := Sweep[float64](circle, Vec(0.0, -10.0), polygon)

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, 0.55, 1e-5), fmt.Sprint(hit))
		geomtest.AssertVector(t, hit.Normal, Vec(0.0, 1.0))

		hit, ok = Sweep[float64](polygon, Vec(0.0, 10.0), circle)

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, 0.55, 1e-5), fmt.Sprint(hit))
		geomtest.AssertVector(t, hit.Normal, Vec(0.0, -1.0))

		_, ok = Sweep[float64](circle, Vec(0.0, -5.0), polygon)

		assert.False(t, ok)
	})
	t.Run("two polygons that are not convex are swept edge against edge", func(t *testing.T) {
		polygon := Pol([]Point[float64]{Pt(0.0, 0.0), Pt(6.0, 0.0), Pt(6.0, 6.0), Pt(4.0, 6.0), Pt(4.0, 2.0), Pt(2.0, 2.0), Pt(2.0, 6.0), Pt(0.0, 6.0)})
		arrow := Pol([]Point[float64]{Pt(3.0, 7.0), Pt(3.5, 9.0), Pt(3.0, 8.5), Pt(2.5, 9.0)})

		hit, ok := Sweep[float64](arrow, Vec(0.0, -10.0), polygon)

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, 0.5, 1e-5), fmt.Sprint(hit))
		geomtest.AssertPoint(t, hit.Point, Pt(3.0, 2.0))
	})
	t.Run("an empty polygon reaches nothing", func(t *testing.T) {
		_, ok := Sweep[float64](Polygon[float64]{}, Vec(10.0, 0.0), Circ(Pt(5.0, 0.0), 1.0))

		assert.False(t, ok)
	})
	t.Run("a point swept is the cast of a segment", func(t *testing.T) {
		for _, start := range pointFixtures {
			for _, motion := range vectorFixtures {
				segment := Seg(start, start.Add(motion))

				for _, c := range circleFixtures {
					hit, ok := segment.HitCircle(c)
					assertPointSweep(t, start, motion, c, hit, ok)
				}
				for _, r := range rectFixtures {
					hit, ok := segment.HitRectangle(r)
					assertPointSweep(t, start, motion, r, hit, ok)
				}
				for _, rp := range regularPolygonFixtures {
					hit, ok := segment.HitRegularPolygon(rp)
					assertPointSweep(t, start, motion, rp, hit, ok)
				}
				for _, b := range boxFixtures {
					hit, ok := segment.HitBox(b)
					assertPointSweep(t, start, motion, b, hit, ok)
				}
			}
		}
	})
	t.Run("moved by the fraction the shapes touch", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				for _, motion := range vectorFixtures {
					hit, ok := Sweep(a, motion, b)
					if !ok {
						continue
					}

					message := fmt.Sprintf("%v + %s → %v: %v", a, motion, b, hit)
					step := motion.Resize(1e-3)

					assert.True(t, hit.Fraction >= 0 && hit.Fraction <= 1, message)
					assert.True(t, hit.Normal.IsNormalized(), message)
					assert.True(t, hit.Normal.Dot(motion) < 0, message)
					assert.True(t, b.DistanceTo(hit.Point) < 1e-3, message)
					assert.False(t, Intersects[float64](translated(a, motion.Multiply(hit.Fraction).Subtract(step)), b), message)
				}
			}
		}
	})
	t.Run("moving the second the other way reaches the first at the same fraction", func(t *testing.T) {
		for _, a := range convexFixtures() {
			for _, b := range convexFixtures() {
				for _, motion := range vectorFixtures {
					hit, ok := Sweep(a, motion, b)
					reverse, found := Sweep(b, motion.Negate(), a)
					message := fmt.Sprintf("%v + %s → %v: %v, %v", a, motion, b, hit, reverse)

					assert.Equal(t, ok, found, message)
					assert.True(t, EqualDelta(hit.Fraction, reverse.Fraction, 1e-5), message)
				}
			}
		}
	})
	t.Run("far from the origin a float32 shape stops at the second", func(t *testing.T) {
		for _, offset := range farOffsets {
			circle := Circ(Pt[float32](0, 0), 1).Translate(offset)
			wall := BoxFromMinMax(Pt[float32](5, -1), Pt[float32](6, 3)).Translate(offset)

			hit, ok := Sweep[float32](circle, Vec[float32](10, 0), wall)

			assert.True(t, ok)
			assert.True(t, EqualDelta(hit.Fraction, 0.4, 1e-2), fmt.Sprint(hit))
			assert.True(t, wall.DistanceTo(hit.Point) < 0.1, fmt.Sprint(hit))
		}
	})
	t.Run("int shapes far from the origin stop as they do near it", func(t *testing.T) {
		far := Vec(1<<25, -(1 << 25))

		hit, ok := Sweep[int](BoxFromMinMax(Pt(0, 0), Pt(1, 1)).Translate(far), Vec(10, 0), BoxFromMinMax(Pt(5, -1), Pt(6, 3)).Translate(far))

		assert.True(t, ok)
		geomtest.AssertNumber(t, hit.Fraction, 0.4)
		assert.Equal(t, hit.Point.X, 5+far.X)
	})
	t.Run("allocates nothing for shapes held as Shape", func(t *testing.T) {
		var circle, wall Shape[float64] = Circ(Pt(0.0, 0.0), 1.0), Rect(Pt(6.0, 0.0), Sz(2.0, 8.0)).Rotate(0.2)

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = Sweep(circle, Vec(10.0, 0.0), wall)
		}), 0)
	})
}

// assertPointSweep asserts that a point swept by the motion reaches the shape where the cast of
// the segment it runs along does, given the hit of that cast: the same fraction, point and
// normal within what the sweep closes on. A cast that only grazes the shape is left out,
// since the sweep closes on a touch in passing without reaching it.
func assertPointSweep(t *testing.T, start Point[float64], motion Vector[float64], shape Shape[float64], expected Hit[float64], ok bool) {
	t.Helper()

	message := fmt.Sprintf("%s + %s → %v: ", start, motion, shape)

	hit, swept := Sweep[float64](Circ(start, 0.0), motion, shape)
	if ok && math.Abs(expected.Normal.Dot(motion.Normalize())) < 1e-3 {
		return
	}

	assert.Equal(t, swept, ok, message)
	if swept && ok {
		assert.True(t, EqualDelta(hit.Fraction, expected.Fraction, 1e-4), message+fmt.Sprint(hit, expected))
		assert.True(t, hit.Point.DistanceTo(expected.Point) < 1e-3, message+fmt.Sprint(hit, expected))
	}
}

func BenchmarkSweep(b *testing.B) {
	var circle, wall, ellipse Shape[float64] = Circ(Pt(0.0, 0.0), 1.0), Rect(Pt(6.0, 0.0), Sz(2.0, 8.0)).Rotate(0.2), Ell(Pt(7.0, 0.5), Sz(2.0, 1.0), 0.2)
	motion := Vec(10.0, 0.0)

	b.Run("onto a flat edge", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Sweep(circle, motion, wall)
		}
	})
	b.Run("onto a round shape", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Sweep(circle, motion, ellipse)
		}
	})
}

func TestRaycast(t *testing.T) {
	t.Run("every kind answers as the Hit method naming it", func(t *testing.T) {
		for _, ray := range rayFixtures {
			for _, c := range circleFixtures {
				expected, found := ray.HitCircle(c)
				hit, ok := Raycast[float64](ray, c)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, c))
			}
			for _, s := range segmentFixtures {
				expected, found := ray.HitSegment(s)
				hit, ok := Raycast[float64](ray, s)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, s))
			}
			for _, p := range polygonFixtures() {
				expected, found := ray.HitPolygon(p)
				hit, ok := Raycast[float64](ray, p)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, p))
			}
			for _, r := range rectFixtures {
				expected, found := ray.HitRectangle(r)
				hit, ok := Raycast[float64](ray, r)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, r))
			}
			for _, rp := range regularPolygonFixtures {
				expected, found := ray.HitRegularPolygon(rp)
				hit, ok := Raycast[float64](ray, rp)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, rp))
			}
			for _, b := range boxFixtures {
				expected, found := ray.HitBox(b)
				hit, ok := Raycast[float64](ray, b)
				assert.True(t, hit == expected && ok == found, fmt.Sprintf("%v → %v: ", ray, b))
			}
		}
	})
	t.Run("an ellipse is reached at its boundary, the fraction in units of Direction", func(t *testing.T) {
		hit, ok := Raycast[float64](Ry(Pt(-10.0, 0.0), Vec(2.0, 0.0)), Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0))

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, 3, 1e-5), fmt.Sprint(hit))
		assert.True(t, EqualDelta(hit.Point.X, -4, 1e-5) && EqualDelta(hit.Point.Y, 0, 1e-3), fmt.Sprint(hit))
		geomtest.AssertVector(t, hit.Normal, Vec(-1.0, 0.0))
	})
	t.Run("an ellipse behind the ray, beside it or around its origin is not hit", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		for _, ray := range []Ray[float64]{Ry(Pt(-10.0, 0.0), Vec(-2.0, 0.0)), Ry(Pt(-10.0, 5.0), Vec(2.0, 0.0)), Ry(Pt(1.0, 0.0), Vec(2.0, 0.0)), Ry(Pt(-10.0, 0.0), Vec(0.0, 0.0))} {
			_, ok := Raycast[float64](ray, ellipse)

			assert.False(t, ok, ray.String())
		}
	})
	t.Run("a shape known by its support is reached where its kind is", func(t *testing.T) {
		rectangle := Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(Pi / 6)
		ray := RayThrough(Pt(-10.0, 1.0), Pt(0.0, 0.0))

		expected, _ := ray.HitRectangle(rectangle)
		hit, ok := Raycast[float64](ray, foreignShape[float64]{rectangle})

		assert.True(t, ok)
		assert.True(t, EqualDelta(hit.Fraction, expected.Fraction, 1e-5), fmt.Sprint(hit, expected))
		assert.True(t, hit.Point.DistanceTo(expected.Point) < 1e-4, fmt.Sprint(hit, expected))
		assert.True(t, hit.Normal.Subtract(expected.Normal).LessOrEqual(1e-4), fmt.Sprint(hit, expected))
	})
	t.Run("an ellipse is reached between the polygon within it and the polygon about it", func(t *testing.T) {
		for _, ellipse := range ellipseFixtures {
			if !ellipse.Size.IsPositive() {
				continue
			}

			inner := ellipse.RegularPolygon(64, OrientationPointyTop)
			outer := ellipse.Scale(1/math.Cos(Pi/64)).RegularPolygon(64, OrientationPointyTop)

			for _, ray := range rayFixtures {
				within, reached := ray.HitRegularPolygon(inner)
				about, _ := ray.HitRegularPolygon(outer)
				if !reached || outer.Contains(ray.Origin) {
					continue
				}

				message := fmt.Sprintf("%v → %v: ", ray, ellipse)

				hit, ok := Raycast[float64](ray, ellipse)
				assert.True(t, ok, message)
				assert.True(t, about.Fraction-1e-5 <= hit.Fraction && hit.Fraction <= within.Fraction+1e-5, message+fmt.Sprint(about.Fraction, hit.Fraction, within.Fraction))
				assert.True(t, hit.Normal.IsNormalized(), message)
				assert.True(t, ellipse.DistanceTo(hit.Point) < 1e-3, message)
			}
		}
	})
	t.Run("a NaN ray hits nothing", func(t *testing.T) {
		ellipse := Ell(Pt(0.0, 0.0), Sz(4.0, 2.0), 0)

		for _, ray := range []Ray[float64]{Ry(Pt(math.NaN(), 0.0), Vec(2.0, 0.0)), Ry(Pt(-10.0, 0.0), Vec(math.NaN(), 0.0))} {
			_, ok := Raycast[float64](ray, ellipse)

			assert.False(t, ok, ray.String())
		}
	})
	t.Run("allocates nothing for a shape held as Shape", func(t *testing.T) {
		var rectangle, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)
		ray := Ry(Pt(-10.0, 0.0), Vec(1.0, 0.0))

		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = Raycast(ray, rectangle)
		}), 0)
		geomtest.AssertNumber(t, testing.AllocsPerRun(100, func() {
			_, sinkBool = Raycast(ray, ellipse)
		}), 0)
	})
}

func BenchmarkRaycast(b *testing.B) {
	var rectangle, ellipse Shape[float64] = Rect(Pt(0.0, 0.0), Sz(4.0, 2.0)).Rotate(0.3), Ell(Pt(3.0, 0.0), Sz(2.0, 1.0), 0.2)
	ray := Ry(Pt(-10.0, 0.0), Vec(1.0, 0.0))

	b.Run("by the method naming the kind", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Raycast(ray, rectangle)
		}
	})
	b.Run("by supports", func(b *testing.B) {
		for b.Loop() {
			_, sinkBool = Raycast(ray, ellipse)
		}
	})
}

func TestTransformable(t *testing.T) {
	t.Run("every shape moves and scales in its own type", func(t *testing.T) {
		geomtest.AssertSegment(t, moved(Seg(Pt(0, 0), Pt(2, 2)), Vec(1, 1)), Seg(Pt(1, 1), Pt(3, 3)))
		geomtest.AssertRay(t, moved(Ry(Pt(0, 0), Vec(2, 2)), Vec(1, 1)), Ry(Pt(1, 1), Vec(2, 2)))
		geomtest.AssertLine(t, moved(Ln(Pt(0, 0), Vec(2, 2)), Vec(1, 1)), Ln(Pt(1, 1), Vec(2, 2)))
		geomtest.AssertRectangle(t, moved(Rect(Pt(0, 0), Sz(2, 2)), Vec(1, 1)), Rect(Pt(1, 1), Sz(2, 2)))
		geomtest.AssertCircle(t, moved(Circ(Pt(0, 0), 2), Vec(1, 1)), Circ(Pt(1, 1), 2))
		geomtest.AssertPolygon(t, moved(Pol(squareVertices()), Vec(1, 1)), Pol([]Point[int]{Pt(1, 1), Pt(3, 1), Pt(3, 3), Pt(1, 3)}))
		geomtest.AssertRegularPolygon(t, moved(RegPol(Pt(0, 0), Sz(2, 2), 6, 0, 0), Vec(1, 1)), RegPol(Pt(1, 1), Sz(2, 2), 6, 0, 0))
	})
}

// moved translates a shape and turns and scales it back and forth through the constraint, the
// call a generic tween makes.
func moved[T Number, S transformable[T, S]](shape S, vector Vector[T]) S {
	return shape.Translate(vector).Rotate(2 * Pi).Scale(2).Unscale(2)
}

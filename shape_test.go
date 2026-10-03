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
	_ Shape[int]    = Segment[int]{}
	_ Shape[int]    = Rectangle[int]{}
	_ Shape[int]    = Circle[int]{}
	_ Shape[int]    = Ellipse[int]{}
	_ Shape[int]    = Polygon[int]{}
	_ Shape[int]    = RegularPolygon[int]{}
	_ Shape[int]    = Box[int]{}
	_ Collider[int] = Segment[int]{}
	_ Collider[int] = Ray[int]{}
	_ Collider[int] = Rectangle[int]{}
	_ Collider[int] = Circle[int]{}
	_ Collider[int] = Polygon[int]{}
	_ Collider[int] = RegularPolygon[int]{}
	_ Collider[int] = Box[int]{}
	_ Body[int]     = Rectangle[int]{}
	_ Body[int]     = Circle[int]{}
	_ Body[int]     = Ellipse[int]{}
	_ Body[int]     = Polygon[int]{}
	_ Body[int]     = RegularPolygon[int]{}
)

// polyline is a shape whose boundary is a chain of straight edges, each starting where the
// previous one ends, the last closing back to the first on a closed shape.
type polyline[T Number] interface {
	Vertices() iter.Seq[Point[T]]
	Edges() iter.Seq[Segment[T]]
}

// encloser is a shape other shapes lie within, with the six Encloses methods: Circle, Polygon,
// Rectangle, RegularPolygon and Box.
type encloser[T Number] interface {
	EnclosesCircle(circle Circle[T]) bool
	EnclosesSegment(segment Segment[T]) bool
	EnclosesPolygon(polygon Polygon[T]) bool
	EnclosesRectangle(rectangle Rectangle[T]) bool
	EnclosesRegularPolygon(polygon RegularPolygon[T]) bool
	EnclosesBox(box Box[T]) bool
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

// solid is a shape with mass properties, the pair the Body properties are checked through.
type solid[T Number] interface {
	Shape[T]
	Body[T]
}

// encloses tests the shape against the container by the method naming its kind, as
// Intersects dispatches. A ray runs without end, so no container encloses one.
func encloses[T Number](container encloser[T], shape Collider[T]) bool {
	switch shape := shape.(type) {
	case Circle[T]:
		return container.EnclosesCircle(shape)
	case Segment[T]:
		return container.EnclosesSegment(shape)
	case Ray[T]:
		return false
	case Polygon[T]:
		return container.EnclosesPolygon(shape)
	case Rectangle[T]:
		return container.EnclosesRectangle(shape)
	case RegularPolygon[T]:
		return container.EnclosesRegularPolygon(shape)
	default:
		return container.EnclosesBox(shape.(Box[T]))
	}
}

func TestShape(t *testing.T) {
	shapes := []Shape[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Circ(Pt(1.0, 1.0), 2.0),
		Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/6),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop),
		BoxFromMinMax(Pt(-1.0, 0.5), Pt(2.0, 3.0)),
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

// assertCrossings asserts the properties the boundary crossings of a segment have on every
// outline: each point lies on the segment and on an edge by Contains, they follow one another
// from Start, and there is one only where the pair intersects.
func assertCrossings[T Number](t *testing.T, segment Segment[T], shape polyline[T], points []Point[T], intersects bool) {
	t.Helper()

	message := fmt.Sprintf("%v → %v: ", segment, shape)

	reached := 0.0
	for _, p := range points {
		distance := segment.Start.DistanceTo(p)

		assert.True(t, segment.Contains(p), message+p.String()+" on the segment: ")
		assert.True(t, slices.ContainsFunc(slices.Collect(shape.Edges()), func(edge Segment[T]) bool {
			return edge.Contains(p)
		}), message+p.String()+" on the boundary: ")
		assert.True(t, LessOrEqual(reached, distance), message+p.String()+" in order: ")

		reached = distance
	}
	assert.True(t, len(points) == 0 || intersects, message)
}

func TestPolyline(t *testing.T) {
	polylines := []polyline[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop),
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
		assert.True(t, closed[float64](Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop)))
	})
}

func TestCollider(t *testing.T) {
	colliders := []Collider[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		RayThrough(Pt(-1.0, -1.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Circ(Pt(1.0, 1.0), 2.0),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop),
		BoxFromMinMax(Pt(-1.0, 0.5), Pt(2.0, 3.0)),
	}
	for _, s := range segmentFixtures {
		colliders = append(colliders, s)
	}
	for _, r := range rayFixtures {
		colliders = append(colliders, r)
	}
	for _, r := range rectFixtures {
		colliders = append(colliders, r)
	}
	for _, c := range circleFixtures {
		colliders = append(colliders, c)
	}
	for _, p := range polygonFixtures() {
		colliders = append(colliders, p)
	}
	for _, b := range boxFixtures {
		colliders = append(colliders, b)
	}

	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, a := range colliders {
			for _, b := range colliders {
				assert.Equal(t, Intersects(a, b), Intersects(b, a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
	t.Run("a shape enclosed by another intersects it", func(t *testing.T) {
		for _, a := range colliders {
			if container, ok := a.(encloser[float64]); ok {
				for _, b := range colliders {
					assert.True(t, !encloses(container, b) || Intersects(a, b), fmt.Sprintf("%s → %s: ", a, b))
				}
			}
		}
	})
	t.Run("every shape with an area encloses itself", func(t *testing.T) {
		for _, a := range colliders {
			if container, ok := a.(encloser[float64]); ok {
				assert.True(t, encloses(container, a), fmt.Sprintf("%s: ", a))
			}
		}
	})
	t.Run("far from the origin float32 colliders answer alike from either side and enclose themselves", func(t *testing.T) {
		for _, offset := range farOffsets {
			far := farColliders(offset)

			for _, a := range far {
				for _, b := range far {
					assert.Equal(t, Intersects(a, b), Intersects(b, a), fmt.Sprintf("%s → %s: ", a, b))
				}
				if container, ok := a.(encloser[float32]); ok {
					assert.True(t, encloses(container, a), fmt.Sprintf("%s: ", a))
					for _, b := range far {
						assert.True(t, !encloses(container, b) || Intersects(a, b), fmt.Sprintf("%s → %s: ", a, b))
					}
				}
			}
		}
	})
}

// stubCollider is a Collider of no kind this package knows, for the paths Intersects takes
// when it cannot name one.
type stubCollider struct {
	result bool
}

func (s stubCollider) IntersectsCircle(Circle[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsSegment(Segment[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsRay(Ray[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsPolygon(Polygon[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsRectangle(Rectangle[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsRegularPolygon(RegularPolygon[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsBox(Box[float64]) bool {
	return s.result
}

// farOffsets move the float32 fixtures to where an ulp of float32 exceeds Delta32, so a point
// rounded into T stays on the boundary it was computed on only by the tolerance epsilonAt
// widens with the coordinates.
var farOffsets = []Vector[float32]{Vec[float32](1e4, -5e3), Vec[float32](1e5, -5e4)}

// farShapes returns the shape fixtures of every kind cast to float32 and moved by the offset.
func farShapes(offset Vector[float32]) []Shape[float32] {
	var shapes []Shape[float32]
	for _, s := range segmentFixtures {
		shapes = append(shapes, s.Cast[float32]().Translate(offset))
	}
	for _, r := range rectFixtures {
		shapes = append(shapes, r.Cast[float32]().Translate(offset))
	}
	for _, c := range circleFixtures {
		shapes = append(shapes, c.Cast[float32]().Translate(offset))
	}
	for _, e := range ellipseFixtures {
		shapes = append(shapes, e.Cast[float32]().Translate(offset))
	}
	for _, p := range polygonFixtures() {
		shapes = append(shapes, p.Cast[float32]().Translate(offset))
	}
	for _, rp := range regularPolygonFixtures {
		shapes = append(shapes, rp.Cast[float32]().Translate(offset))
	}
	for _, b := range boxFixtures {
		shapes = append(shapes, b.Cast[float32]().Translate(offset))
	}

	return shapes
}

// farColliders returns the collider fixtures of every kind cast to float32 and moved by the
// offset: the far shapes but the ellipse, and the rays.
func farColliders(offset Vector[float32]) []Collider[float32] {
	var colliders []Collider[float32]
	for _, shape := range farShapes(offset) {
		if collider, ok := shape.(Collider[float32]); ok {
			colliders = append(colliders, collider)
		}
	}
	for _, r := range rayFixtures {
		colliders = append(colliders, r.Cast[float32]().Translate(offset))
	}

	return colliders
}

// narrowMatrix stretches the fixtures before they are cast to int16, past the square root of its
// range, so a squared distance or a difference taken in T rather than float64 overflows where
// every coordinate still fits.
var narrowMatrix = ScaleMatrix(40.0, 40.0)

// prefixPoint is the point the Append methods find in dst over the fixtures, so every result is
// checked to follow it.
var prefixPoint = Pt(-7.5, 3.25)

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
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Circ(Pt(1.0, 1.0), 2.0),
		Ell(Pt(1.0, 2.0), Sz(4.0, 2.0), Pi/6),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), OrientationFlatTop),
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
	t.Run("a collider of another kind is tested from the side that has one", func(t *testing.T) {
		rectangle := Rect(Pt(0.0, 0.0), Sz(2.0, 2.0))

		assert.True(t, Intersects[float64](rectangle, stubCollider{result: true}))
		assert.False(t, Intersects[float64](stubCollider{result: false}, rectangle))
	})
	t.Run("two colliders of another kind panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Intersects[float64](stubCollider{}, stubCollider{})
		})
	})
}

func TestTransformable(t *testing.T) {
	t.Run("every shape moves and scales in its own type", func(t *testing.T) {
		geomtest.AssertSegment(t, moved(Seg(Pt(0, 0), Pt(2, 2)), Vec(1, 1)), Seg(Pt(1, 1), Pt(3, 3)))
		geomtest.AssertRay(t, moved(Ry(Pt(0, 0), Vec(2, 2)), Vec(1, 1)), Ry(Pt(1, 1), Vec(2, 2)))
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

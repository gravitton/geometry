package geom

import (
	"fmt"
	"slices"
	"testing"

	"github.com/gravitton/assert"
)

var (
	_ Shape[int]        = Segment[int]{}
	_ Shape[int]        = Rectangle[int]{}
	_ Shape[int]        = Circle[int]{}
	_ Shape[int]        = Polygon[int]{}
	_ Shape[int]        = RegularPolygon[int]{}
	_ Outline[int]      = Segment[int]{}
	_ Outline[int]      = Rectangle[int]{}
	_ Outline[int]      = Polygon[int]{}
	_ Outline[int]      = RegularPolygon[int]{}
	_ Collider[int]     = Segment[int]{}
	_ Collider[int]     = Rectangle[int]{}
	_ Collider[int]     = Circle[int]{}
	_ Collider[int]     = Polygon[int]{}
	_ Collider[int]     = RegularPolygon[int]{}
	_ Measured[int]     = Size[int]{}
	_ Measured[int]     = Rectangle[int]{}
	_ Measured[float64] = Circle[int]{}
	_ Measured[float64] = Polygon[int]{}
	_ Measured[float64] = RegularPolygon[int]{}
)

func TestShape(t *testing.T) {
	shapes := []Shape[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Circ(Pt(1.0, 1.0), 2.0),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop),
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
}

func TestOutline(t *testing.T) {
	outlines := []Outline[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop),
	}

	t.Run("every edge starts at a vertex and ends at the next", func(t *testing.T) {
		for _, outline := range outlines {
			vertices, edges := slices.Collect(outline.Vertices()), slices.Collect(outline.Edges())

			for i, edge := range edges {
				assert.True(t, edge.Start.Equal(vertices[i]), fmt.Sprintf("%s #%d: ", outline, i))
				assert.True(t, edge.End.Equal(vertices[(i+1)%len(vertices)]), fmt.Sprintf("%s #%d: ", outline, i))
			}
		}
	})
	t.Run("only the segment is open", func(t *testing.T) {
		assert.False(t, closed[float64](Seg(Pt(0.0, 0.0), Pt(3.0, 4.0))))
		assert.True(t, closed[float64](Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi/6)))
		assert.True(t, closed[float64](Pol(triangleVertices())))
		assert.True(t, closed[float64](Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop)))
	})
}

func TestCollider(t *testing.T) {
	colliders := []Collider[float64]{
		Seg(Pt(0.0, 0.0), Pt(3.0, 4.0)),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)).Rotate(Pi / 6),
		Circ(Pt(1.0, 1.0), 2.0),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop),
	}
	for _, s := range segmentFixtures {
		colliders = append(colliders, s)
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

	t.Run("either side gives the same answer", func(t *testing.T) {
		for _, a := range colliders {
			for _, b := range colliders {
				assert.Equal(t, Intersects(a, b), Intersects(b, a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
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

func (s stubCollider) IntersectsPolygon(Polygon[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsRectangle(Rectangle[float64]) bool {
	return s.result
}

func (s stubCollider) IntersectsRegularPolygon(RegularPolygon[float64]) bool {
	return s.result
}

func TestMeasured(t *testing.T) {
	measured := []Measured[float64]{
		Sz(3.0, 4.0),
		Rect(Pt(1.0, 2.0), Sz(4.0, 2.0)),
		Circ(Pt(1.0, 1.0), 2.0),
		Pol(triangleVertices()),
		Hexagon(Pt(0.0, 0.0), SzU(10.0), FlatTop),
	}

	t.Run("area and perimeter are non-negative", func(t *testing.T) {
		for _, m := range measured {
			assert.True(t, m.Area() >= 0, fmt.Sprintf("%v: ", m))
			assert.True(t, m.Perimeter() >= 0, fmt.Sprintf("%v: ", m))
		}
	})
}

// closed reports whether the last edge returns to the first vertex, through the interface as
// a caller would.
func closed[T Number](outline Outline[T]) bool {
	var first, last Point[T]
	for vertex := range outline.Vertices() {
		first = vertex

		break
	}
	for edge := range outline.Edges() {
		last = edge.End
	}

	return first.Equal(last)
}

func TestMovable(t *testing.T) {
	t.Run("every shape moves and scales in its own type", func(t *testing.T) {
		AssertSegment(t, moved(Seg(Pt(0, 0), Pt(2, 2)), Vec(1, 1)), Seg(Pt(1, 1), Pt(3, 3)))
		AssertRectangle(t, moved(Rect(Pt(0, 0), Sz(2, 2)), Vec(1, 1)), Rect(Pt(1, 1), Sz(2, 2)))
		AssertCircle(t, moved(Circ(Pt(0, 0), 2), Vec(1, 1)), Circ(Pt(1, 1), 2))
		AssertPolygon(t, moved(Pol(squareVertices()), Vec(1, 1)), Pol([]Point[int]{Pt(1, 1), Pt(3, 1), Pt(3, 3), Pt(1, 3)}))
		AssertRegularPolygon(t, moved(RegPol(Pt(0, 0), Sz(2, 2), 6, 0), Vec(1, 1)), RegPol(Pt(1, 1), Sz(2, 2), 6, 0))
	})
}

// moved translates a shape and turns and scales it back and forth through the constraint, the
// call a generic tween makes.
func moved[T Number, S Movable[T, S]](shape S, vector Vector[T]) S {
	return shape.Translate(vector).Rotate(2 * Pi).Scale(2).Unscale(2)
}

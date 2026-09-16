package geom

import (
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

func TestCollisionRectangles(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, CollisionRectangles(rectangle, Rect(Pt(100.0, -50.0), Sz(200.0, 50.0))))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, CollisionRectangles(rectangle, Rect(Pt(100.0, 350.0), Sz(200.0, 450.0))))
	})
	t.Run("a shared edge counts as a collision", func(t *testing.T) {
		assert.True(t, CollisionRectangles(rectangle, Rect(Pt(200.0, 0.0), Sz(200.0, 100.0))))
		assert.False(t, CollisionRectangles(rectangle, Rect(Pt(201.0, 0.0), Sz(200.0, 100.0))))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, CollisionRectangles(rectangle, Rect(Pt(0.0, 0.0), Sz(50.0, 50.0))))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range rectFixtures {
			for _, b := range rectFixtures {
				assert.Equal(t, CollisionRectangles(a, b), CollisionRectangles(b, a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestCollisionCircles(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 100.0)

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, CollisionCircles(circle, Circ(Pt(199.0, 0.0), 100.0)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, CollisionCircles(circle, Circ(Pt(210.0, 0.0), 100.0)))
	})
	t.Run("exactly touching does not count", func(t *testing.T) {
		assert.False(t, CollisionCircles(circle, Circ(Pt(200.0, 0.0), 100.0)))
	})
	t.Run("one contained in the other", func(t *testing.T) {
		assert.True(t, CollisionCircles(circle, Circ(Pt(0.0, 0.0), 50.0)))
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, a := range circleFixtures {
			for _, b := range circleFixtures {
				assert.Equal(t, CollisionCircles(a, b), CollisionCircles(b, a), fmt.Sprintf("%s → %s: ", a, b))
			}
		}
	})
}

func TestCollisionRectangleCircle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	t.Run("overlapping", func(t *testing.T) {
		assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 60.0)))
		assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(110.0, 80.0), 60.0)))
	})
	t.Run("apart", func(t *testing.T) {
		assert.False(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 40.0)))
	})
	t.Run("the circle center on the boundary counts", func(t *testing.T) {
		assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(100.0, 0.0), 1.0)))
	})
	t.Run("touching the edge from outside counts", func(t *testing.T) {
		assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(200.0, 0.0), 100.0)))
		assert.False(t, CollisionRectangleCircle(rectangle, Circ(Pt(201.0, 0.0), 100.0)))
	})
	t.Run("circle fully inside the rectangle", func(t *testing.T) {
		assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(0.0, 0.0), 10.0)))
	})
	t.Run("a circle collides with its own bounds", func(t *testing.T) {
		for _, c := range circleFixtures {
			if c.Radius == 0 {
				continue // a degenerate circle touches nothing
			}

			assert.True(t, CollisionRectangleCircle(c.Bounds(), c), fmt.Sprintf("%s: ", c))
		}
	})
}

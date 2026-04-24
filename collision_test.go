package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestCollisionRectangles(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	assert.True(t, CollisionRectangles(rectangle, Rect(Pt(100.0, -50.0), Sz(200.0, 50.0))))
	assert.False(t, CollisionRectangles(rectangle, Rect(Pt(100.0, 350.0), Sz(200.0, 450.0))))

	// Rectangles sharing exactly one edge count as colliding.
	assert.True(t, CollisionRectangles(rectangle, Rect(Pt(200.0, 0.0), Sz(200.0, 100.0))))
	// One rectangle just beyond the edge does not collide.
	assert.False(t, CollisionRectangles(rectangle, Rect(Pt(201.0, 0.0), Sz(200.0, 100.0))))
	// One rectangle fully contained inside the other.
	assert.True(t, CollisionRectangles(rectangle, Rect(Pt(0.0, 0.0), Sz(50.0, 50.0))))
}

func TestCollisionCircles(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 100.0)

	assert.True(t, CollisionCircles(circle, Circ(Pt(199.0, 0.0), 100.0)))
	assert.False(t, CollisionCircles(circle, Circ(Pt(210.0, 0.0), 100.0)))

	// Circles exactly touching (distance == sum of radii) do not count as colliding.
	assert.False(t, CollisionCircles(circle, Circ(Pt(200.0, 0.0), 100.0)))
	// One circle fully inside the other.
	assert.True(t, CollisionCircles(circle, Circ(Pt(0.0, 0.0), 50.0)))
}

func TestCollisionRectangleCircle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 60.0)))
	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(110.0, 80.0), 60.0)))
	assert.False(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 40.0)))

	// Circle center exactly on rectangle boundary counts as colliding.
	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(100.0, 0.0), 1.0)))
	// Circle touching rectangle edge from outside counts as colliding.
	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(200.0, 0.0), 100.0)))
	// Circle just beyond rectangle edge does not collide.
	assert.False(t, CollisionRectangleCircle(rectangle, Circ(Pt(201.0, 0.0), 100.0)))
	// Circle fully inside rectangle.
	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(0.0, 0.0), 10.0)))
}

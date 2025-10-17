package geom

import (
	"testing"

	"github.com/gravitton/assert"
)

func TestCollisionRectangles(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	assert.True(t, CollisionRectangles(rectangle, Rect(Pt(100.0, -50.0), Sz(200.0, 50.0))))
	assert.False(t, CollisionRectangles(rectangle, Rect(Pt(100.0, 350.0), Sz(200.0, 450.0))))
}

func TestCollisionCircles(t *testing.T) {
	circle := Circ(Pt(0.0, 0.0), 100.0)

	assert.True(t, CollisionCircles(circle, Circ(Pt(199.0, 0.0), 100.0)))
	assert.False(t, CollisionCircles(circle, Circ(Pt(210.0, 0.0), 100.0)))
}

func TestCollisionRectangleCircle(t *testing.T) {
	rectangle := Rect(Pt(0.0, 0.0), Sz(200.0, 100.0))

	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 60.0)))
	assert.True(t, CollisionRectangleCircle(rectangle, Circ(Pt(110.0, 80.0), 60.0)))
	assert.False(t, CollisionRectangleCircle(rectangle, Circ(Pt(150.0, 0.0), 40.0)))
}

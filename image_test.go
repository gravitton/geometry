package geom

import (
	"image"
	"testing"

	"github.com/gravitton/assert"
)

func TestPoint_ConstructorFromImage(t *testing.T) {
	AssertPoint(t, PointFromImage[int](image.Pt(1, 2)), Pt(1, 2))
	AssertPoint(t, PointFromImage[float64](image.Pt(1, 2)), Pt(1.0, 2.0))
}

func TestPoint_Point(t *testing.T) {
	assert.Equal(t, Pt(1, 2).Point(), image.Pt(1, 2))
	assert.Equal(t, Pt(0.6, -0.25).Point(), image.Pt(1, 0))
}

func TestRectangle_Rectangle(t *testing.T) {
	assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Rectangle(), image.Rect(0, 1, 2, 4))
	assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Rectangle(), image.Rect(0, -2, 1, 2))
}

func TestSize_ConstructorFromImage(t *testing.T) {
	AssertSize(t, SizeFromImage[int](image.Rect(0, 10, 55, 70)), Sz(55, 60))
	AssertSize(t, SizeFromImage[float64](image.Rect(0, 10, 55, 70)), Sz(55.0, 60.0))
}

func TestRectangle_ConstructorFromImage(t *testing.T) {
	AssertRect(t, RectangleFromImage[int](image.Rect(0, 10, 55, 70)), Rect(Pt(27, 40), Sz(55, 60)))
	AssertRect(t, RectangleFromImage[float64](image.Rect(0, 10, 55, 70)), Rect(Pt(27.5, 40), Sz(55.0, 60.0)))
}

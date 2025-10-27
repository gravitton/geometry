package geom

import (
	"image"
	"testing"

	"github.com/gravitton/assert"
)

func TestPoint_NewFromImage(t *testing.T) {
	AssertPoint(t, PointFromImage[int](image.Pt(1, 2)), 1, 2)
	AssertPoint(t, PointFromImage[float64](image.Pt(1, 2)), 1.0, 2.0)
}

func TestPoint_Point(t *testing.T) {
	assert.Equal(t, pointInt.Point(), image.Pt(1, 2))
	assert.Equal(t, pointFloat.Point(), image.Pt(1, 0))
}

func TestRectangle_Rectangle(t *testing.T) {
	assert.Equal(t, rectInt.Rectangle(), image.Rect(0, 1, 2, 4))
	assert.Equal(t, rectFloat.Rectangle(), image.Rect(0, -2, 1, 2))
}

func TestSize_NewFromImage(t *testing.T) {
	AssertSize(t, SizeFromImage[int](image.Rect(0, 10, 55, 70)), 55, 60)
	AssertSize(t, SizeFromImage[float64](image.Rect(0, 10, 55, 70)), 55.0, 60.0)
}

func TestRectangle_NewFromImage(t *testing.T) {
	AssertRect(t, RectFromImage[int](image.Rect(0, 10, 55, 70)), 27, 40, 55, 60)
	AssertRect(t, RectFromImage[float64](image.Rect(0, 10, 55, 70)), 27.5, 40, 55, 60)
}

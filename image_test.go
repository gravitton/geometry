package geom

import (
	"fmt"
	"image"
	"testing"

	"github.com/gravitton/assert"
)

func TestPointFromImage(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertPoint(t, PointFromImage[int](image.Pt(1, 2)), Pt(1, 2))
	})
	t.Run("float", func(t *testing.T) {
		AssertPoint(t, PointFromImage[float64](image.Pt(1, 2)), Pt(1.0, 2.0))
	})
}

func TestSizeFromImage(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertSize(t, SizeFromImage[int](image.Rect(0, 10, 55, 70)), Sz(55, 60))
	})
	t.Run("float", func(t *testing.T) {
		AssertSize(t, SizeFromImage[float64](image.Rect(0, 10, 55, 70)), Sz(55.0, 60.0))
	})
}

func TestRectangleFromImage(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		AssertRectangle(t, RectangleFromImage[int](image.Rect(0, 10, 55, 70)), Rect(Pt(27, 40), Sz(55, 60)))
	})
	t.Run("float", func(t *testing.T) {
		AssertRectangle(t, RectangleFromImage[float64](image.Rect(0, 10, 55, 70)), Rect(Pt(27.5, 40.0), Sz(55.0, 60.0)))
	})
}

func TestPoint_Image(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Pt(1, 2).Point(), image.Pt(1, 2))
	})
	t.Run("float rounds", func(t *testing.T) {
		assert.Equal(t, Pt(0.6, -0.25).Point(), image.Pt(1, 0))
	})
	t.Run("is the image point of Int", func(t *testing.T) {
		for _, p := range pointFixtures {
			assert.Equal(t, p.Point(), p.Int().Point(), p.String())
		}
	})
}

func TestRectangle_Image(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(1, 2), Sz(2, 3)).Rectangle(), image.Rect(0, 1, 2, 4))
	})
	t.Run("float rounds like Int", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.6, -0.25), Sz(1.2, 3.6)).Rectangle(), image.Rect(1, -2, 2, 2))
	})
	t.Run("float spans exactly the rounded size", func(t *testing.T) {
		assert.Equal(t, Rect(Pt(0.5, 0.5), SzU(16.0)).Rectangle(), image.Rect(-7, -7, 9, 9))
		assert.Equal(t, Rect(Pt(0.4, 0.4), SzU(0.2)).Rectangle(), image.Rect(0, 0, 0, 0))
	})
	t.Run("is the image rectangle of Int", func(t *testing.T) {
		for _, r := range rectFixtures {
			assert.Equal(t, r.Rectangle(), r.Int().Rectangle(), r.String())
		}
	})
	t.Run("round-trips through the image package", func(t *testing.T) {
		// only rectangles on the integer lattice with even extents survive: the center is
		// stored, so an odd extent cannot be halved exactly
		for _, r := range []Rectangle[int]{
			RectangleFromMin(Pt(0, 0), Sz(4, 2)),
			RectangleFromMin(Pt(-6, 10), Sz(20, 60)),
			RectangleFromMinMax(Pt(-2, -4), Pt(6, 2)),
		} {
			AssertRectangle(t, RectangleFromImage[int](r.Rectangle()), r, fmt.Sprintf("%s: ", r))
		}
	})
}

package geom

import (
	"image"
)

// PointFromImage converts a Point from an image.Point.
func PointFromImage[T Number](p image.Point) Point[T] {
	return Point[T]{T(p.X), T(p.Y)}
}

// SizeFromImage converts a Size from an image.Rectangle.
func SizeFromImage[T Number](r image.Rectangle) Size[T] {
	return Size[T]{T(r.Dx()), T(r.Dy())}
}

// RectangleFromImage converts a Rectangle from an image.Rectangle.
func RectangleFromImage[T Number](r image.Rectangle) Rectangle[T] {
	return RectangleFromMin(PointFromImage[T](r.Min), SizeFromImage[T](r))
}

// Point converts the point to an image.Point, rounding a float T.
func (p Point[T]) Point() image.Point {
	return image.Point(p.Int())
}

// Rectangle converts the rectangle to a half-open image.Rectangle, the Min and Max of its Int:
// it spans exactly Size.Int pixels for any T, so a sprite keeps its width as it moves through
// sub-pixel positions, where rounding the two corners on their own would have it flicker by a
// pixel. The box is placed where Int places it. A rotated rectangle gives the box around its
// rounded vertices, since an image.Rectangle has no angle.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	rounded := r.Int()

	return image.Rectangle{rounded.Min().Point(), rounded.Max().Point()}
}

package geom

import (
	"image"
)

// PointFromImage returns the point of an image.Point.
func PointFromImage[T Number](p image.Point) Point[T] {
	return Point[T]{T(p.X), T(p.Y)}
}

// SizeFromImage returns the size of an image.Rectangle.
func SizeFromImage[T Number](r image.Rectangle) Size[T] {
	return Size[T]{T(r.Dx()), T(r.Dy())}
}

// RectangleFromImage returns the rectangle of an image.Rectangle. For an integer T it inverts
// Rectangle.Rectangle; a float T centers an odd extent on a half, which Rectangle rounds away from
// zero, so the box it gives back can lie a pixel off.
func RectangleFromImage[T Number](r image.Rectangle) Rectangle[T] {
	return RectangleFromMin(PointFromImage[T](r.Min), SizeFromImage[T](r))
}

// BoxFromImage returns the box of an image.Rectangle, between its Min and Max. For an integer T it
// keeps every extent, odd ones included, where RectangleFromImage halves it about a center.
func BoxFromImage[T Number](r image.Rectangle) Box[T] {
	return BoxFromMinMax(PointFromImage[T](r.Min), PointFromImage[T](r.Max))
}

// Point converts the point to an image.Point, rounding a float T.
func (p Point[T]) Point() image.Point {
	return image.Point(p.Int())
}

// Rectangle converts the rectangle to the half-open image.Rectangle of the Bounds of its Int, so it
// spans exactly Size.Int pixels and keeps its width through sub-pixel positions. A rotated rectangle
// gives the box around its rounded vertices.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	bounds := r.Int().Bounds()

	return image.Rectangle{bounds.Min.Point(), bounds.Max.Point()}
}

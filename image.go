package geom

import (
	"image"
)

// PointFromImage returns the point of an image.Point.
func PointFromImage[T Number](point image.Point) Point[T] {
	return Point[T]{T(point.X), T(point.Y)}
}

// SizeFromImage returns the size of an image.Rectangle.
func SizeFromImage[T Number](rectangle image.Rectangle) Size[T] {
	return Size[T]{T(rectangle.Dx()), T(rectangle.Dy())}
}

// RectangleFromImage returns the rectangle of an image.Rectangle. For an integer T it inverts
// Rectangle.Rectangle; a float T centers an odd extent on a half, which Rectangle rounds away from
// zero, so the box it gives back can lie a pixel off.
func RectangleFromImage[T Number](rectangle image.Rectangle) Rectangle[T] {
	return RectangleFromMin(PointFromImage[T](rectangle.Min), SizeFromImage[T](rectangle))
}

// BoxFromImage returns the box of an image.Rectangle, between its Min and Max, for every T, where a
// float RectangleFromImage can come back a pixel off.
func BoxFromImage[T Number](rectangle image.Rectangle) Box[T] {
	return BoxFromMinMax(PointFromImage[T](rectangle.Min), PointFromImage[T](rectangle.Max))
}

// Point converts the point to an image.Point, rounding a float T.
func (p Point[T]) Point() image.Point {
	return image.Point(p.Int())
}

// Rectangle converts the rectangle to the half-open image.Rectangle of the Bounds of its Int, so it
// spans exactly Size.Int pixels and keeps its width through sub-pixel positions. A turned
// rectangle gives the box around its rounded vertices.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	bounds := r.Int().Bounds()

	return image.Rectangle{bounds.Min.Point(), bounds.Max.Point()}
}

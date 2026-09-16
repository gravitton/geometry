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

// Rectangle converts the rectangle to a half-open image.Rectangle spanning Width by Height pixels.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	return image.Rectangle{r.Min().Point(), r.Max().Point()}
}

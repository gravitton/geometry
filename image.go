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

// RectFromImage converts a Rectangle from an image.Rectangle.
func RectFromImage[T Number](r image.Rectangle) Rectangle[T] {
	return RectFromMin(PointFromImage[T](r.Min), SizeFromImage[T](r))
}

// Point converts a Point to an image.Point.
func (p Point[T]) Point() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

// Rectangle converts a Rectangle to an image.Rectangle.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	return image.Rectangle{r.Min().Point(), r.Max().Point()}
}

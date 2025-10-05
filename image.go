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

// ToImagePoint converts a Point to an image.Point.
func (p Point[T]) ToImagePoint() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

// ToImageRectangle converts a Size to an image.Rectangle.
func (s Size[T]) ToImageRectangle() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{int(s.Width), int(s.Height)}}
}

// ToImageRectangle converts a Rectangle to an image.Rectangle.
func (r Rectangle[T]) ToImageRectangle() image.Rectangle {
	return image.Rectangle{r.Min().ToImagePoint(), r.Max().ToImagePoint()}
}

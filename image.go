package geom

import (
	"image"
)

// PtFromImage converts a Point from an image.Point.
func PtFromImage[T Number](p image.Point) Point[T] {
	return Point[T]{T(p.X), T(p.Y)}
}

// SzFromImage converts a Size from an image.Rectangle.
func SzFromImage[T Number](r image.Rectangle) Size[T] {
	return Size[T]{T(r.Dx()), T(r.Dy())}
}

// RectFromImage converts a Rectangle from an image.Rectangle.
func RectFromImage[T Number](r image.Rectangle) Rectangle[T] {
	return RectFromMin(PtFromImage[T](r.Min), SzFromImage[T](r))
}

// Point converts a Point to an image.Point.
func (p Point[T]) Point() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

// Rectangle converts a Rectangle to an image.Rectangle.
func (r Rectangle[T]) Rectangle() image.Rectangle {
	return image.Rectangle{r.Min().Point(), r.Max().Point()}
}

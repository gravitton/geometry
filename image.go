package geom

import (
	"image"
)

func PointFromImage[T Number](p image.Point) Point[T] {
	return Point[T]{T(p.X), T(p.Y)}
}

func SizeFromImage[T Number](r image.Rectangle) Size[T] {
	return Size[T]{T(r.Dx()), T(r.Dy())}
}

func RectangleFromImage[T Number](r image.Rectangle) Rectangle[T] {
	return RectangleFromMin(PointFromImage[T](r.Min), SizeFromImage[T](r))
}

func (p Point[T]) ToImagePoint() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

func (s Size[T]) ToImageRectangle() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{int(s.Width), int(s.Height)}}
}

func (r Rectangle[T]) ToImageRectangle() image.Rectangle {
	return image.Rectangle{r.Min().ToImagePoint(), r.Max().ToImagePoint()}
}

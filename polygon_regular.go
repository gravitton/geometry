package geom

import (
	"math"
)

// RegularPolygon is a polygon with equally spaced Transform around a center.
type RegularPolygon[T Number] struct {
	Center Point[T]
	Size   Size[T]
	N      int
}

// RP is shorthand for RegularPolygon{center, size, n}.
func RP[T Number](center Point[T], size Size[T], n int) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, n}
}

// Triangle creates a RegularPolygon with 3 Transform.
func Triangle[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 3}
}

// Square creates a RegularPolygon with 4 Transform.
func Square[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 4}
}

// Hexagon creates a RegularPolygon with 6 Transform.
func Hexagon[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 6}
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(change Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(change), rp.Size, rp.N}
}

// MoveTo creates a new RegularPolygon with the same size and sides centered at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N}
}

// Multiple creates a new RegularPolygon with size scaled by the given factor.
func (rp RegularPolygon[T]) Scale(scale float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(scale), rp.N}
}

// Multiple creates a new RegularPolygon with size scaled by the given factors.
func (rp RegularPolygon[T]) ScaleXY(scaleX, scaleY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(scaleX, scaleY), rp.N}
}

// Vertices returns the polygon Transform in order starting from angle 0, counter-clockwise.
func (rp RegularPolygon[T]) Vertices() []Point[T] {
	initAngle := 0.0
	angleStep := 2 * math.Pi / float64(rp.N)

	vertices := make([]Point[T], rp.N)
	for i := 0; i < rp.N; i++ {
		vertices[i] = rp.Center.Add(VecFromAngle[T](initAngle+float64(i)*angleStep, 1).MultiplyXY(float64(rp.Size.Width), float64(rp.Size.Height)))
	}

	return vertices
}

// ToPolygon converts the regular polygon into a generic Polygon with computed vertices.
func (rp RegularPolygon[T]) ToPolygon() Polygon[T] {
	return Polygon[T]{rp.Vertices()}
}

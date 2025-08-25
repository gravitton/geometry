package geom

import (
	"math"
)

// RegularPolygon is a polygon with equally spaced vertices around a center.
type RegularPolygon[T Number] struct {
	Center Point[T]
	Size   Size[T]
	N      int
}

// RP is shorthand for RegularPolygon{center, size, n}.
func RP[T Number](center Point[T], size Size[T], n int) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, n}
}

// Triangle creates a RegularPolygon with 3 vertices.
func Triangle[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 3}
}

// Square creates a RegularPolygon with 4 vertices.
func Square[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 4}
}

// Hexagon creates a RegularPolygon with 6 vertices.
func Hexagon[T Number](center Point[T], size Size[T]) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 6}
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(change Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(change), rp.Size, rp.N}
}

// MoveTo creates a new RegularPolygon with center at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N}
}

// Multiple creates a new RegularPolygon with size scaled by the given factor.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(factor), rp.N}
}

// Multiple creates a new RegularPolygon with size scaled by the given factors.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(factorX, factorY), rp.N}
}

// Vertices returns the polygon vertices in order starting from angle 0, counter-clockwise.
func (rp RegularPolygon[T]) Vertices() []Point[T] {
	initAngle := 0.0
	angleStep := 2 * math.Pi / float64(rp.N)

	vertices := make([]Point[T], rp.N)
	for i := 0; i < rp.N; i++ {
		vertices[i] = rp.Center.Add(VecFromAngle[T](initAngle+float64(i)*angleStep, 1).MultiplyXY(float64(rp.Size.Width), float64(rp.Size.Height)))
	}

	return vertices
}

// Bounds returns the axis-aligned bounding rectangle.
func (rp RegularPolygon[T]) Bounds() Rectangle[T] {
	return Rectangle[T]{rp.Center, rp.Size}
}

// ToPolygon converts the regular polygon into a generic Polygon with computed vertices.
func (rp RegularPolygon[T]) ToPolygon() Polygon[T] {
	return Polygon[T]{rp.Vertices()}
}

// Equal checks if center point, size and number of vertices are equal.
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool {
	return rp.Center.Equal(polygon.Center) && rp.Size.Equal(polygon.Size) && rp.N == polygon.N
}

// IsZero checks if center point, size and number of vertices are zero.
func (rp RegularPolygon[T]) IsZero() bool {
	return rp.Center.IsZero() && rp.Size.IsZero() && rp.N == 0
}

// Empty checks if number of vertices is zero.
func (rp RegularPolygon[T]) Empty() bool {
	return rp.N == 0
}

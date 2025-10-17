package geom

import (
	"fmt"
	"math"
)

// RegularPolygon is a polygon with equally spaced vertices around a center.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",inline"`
	Size   Size[T]  `json:",inline"`
	N      int      `json:"n"`
	Angle  float64  `json:"a"`
}

// RegPol is shorthand for RegularPolygon{center, size, n, angle}.
func RegPol[T Number](center Point[T], size Size[T], n int, angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, n, angle}
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(change Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(change), rp.Size, rp.N, rp.Angle}
}

// MoveTo creates a new RegularPolygon with center at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N, rp.Angle}
}

// Scale creates a new RegularPolygon with size scaled by the given factor.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(factor), rp.N, rp.Angle}
}

// ScaleXY creates a new RegularPolygon with size scaled by the given factors.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(factorX, factorY), rp.N, rp.Angle}
}

// Rotate creates a new RegularPolygon rotated by the given angle (in radians).
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	// TODO: normalize angle to [0, 2*pi]
	return RegularPolygon[T]{rp.Center, rp.Size, rp.N, rp.Angle + angle}
}

// Vertices returns the polygon vertices in order starting from angle 0, counter-clockwise.
func (rp RegularPolygon[T]) Vertices() []Point[T] {
	initAngle := rp.Angle
	angleStep := (2 * math.Pi) / float64(rp.N)

	vertices := make([]Point[T], rp.N)
	for i := 0; i < rp.N; i++ {
		// TODO: Vec(1,0).Rotate(angle)
		vertices[i] = rp.Center.Add(VectorFromAngle[T](initAngle+float64(i)*angleStep, 1).MultiplyXY(float64(rp.Size.Width), float64(rp.Size.Height)))
	}

	return vertices
}

// Bounds returns the axis-aligned bounding rectangle.
func (rp RegularPolygon[T]) Bounds() Rectangle[T] {
	// TODO: calculate
	maxAbsCos, maxAbsSin := 1.0, 1.0
	return Rectangle[T]{rp.Center, rp.Size.ScaleXY(2.0*maxAbsCos, 2.0*maxAbsSin)}
}

// Polygon converts the regular polygon into a generic Polygon with computed vertices.
func (rp RegularPolygon[T]) Polygon() Polygon[T] {
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

// Int converts the regular polygon to a [int] regular polygon.
func (rp RegularPolygon[T]) Int() RegularPolygon[int] {
	return RegularPolygon[int]{rp.Center.Int(), rp.Size.Int(), rp.N, rp.Angle}
}

// Float converts the regular polygon to a [float64] regular polygon.
func (rp RegularPolygon[T]) Float() RegularPolygon[float64] {
	return RegularPolygon[float64]{rp.Center.Float(), rp.Size.Float(), rp.N, rp.Angle}
}

// String returns a string representation of the RegularPolygon.
func (rp RegularPolygon[T]) String() string {
	return fmt.Sprintf("RegPol(%s;%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N), String(rp.Angle))
}

type Orientation int

const (
	FlatTop Orientation = iota
	PointTop
)

func RegularPolygonAngle(n int, orientation Orientation) float64 {
	switch orientation {
	case FlatTop:
		// 90 - 180/n degrees
		return math.Pi * float64(n-2) / (2 * float64(n))
	case PointTop:
		// 90 degrees
		return math.Pi / 2
	default:
		return 0
	}
}

// Triangle creates a RegularPolygon with 3 vertices.
func Triangle[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 3, RegularPolygonAngle(3, orientation)}
}

// Square creates a RegularPolygon with 4 vertices.
func Square[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 4, RegularPolygonAngle(4, orientation)}
}

// Hexagon creates a RegularPolygon with 6 vertices.
func Hexagon[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, 6, RegularPolygonAngle(6, orientation)}
}

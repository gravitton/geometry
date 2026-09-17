package geom

import (
	"fmt"
)

// RegularPolygon is a polygon with equally spaced vertices around a center.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	N      int      `json:"n"`
	Angle  float64  `json:"a"`
}

// RegPol is shorthand for RegularPolygon{center, size, n, angle}.
func RegPol[T Number](center Point[T], size Size[T], n int, angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, n, angle}
}

// Orientation defines the rotational alignment of a regular polygon.
type Orientation int

const (
	// FlatTop places a flat edge at the top of the polygon.
	FlatTop Orientation = iota
	// PointyTop places a vertex at the top of the polygon.
	PointyTop
)

// RegularPolygonOrientationAngle returns the initial vertex angle for a regular polygon with n sides
// and the given orientation, normalized to [0, 2π) like Rotate. PointyTop puts the first vertex at
// the top (-Y, 3π/2); FlatTop puts the midpoint of an edge there, so the first vertex sits half a
// step before it at 3π/2 - π/n. A polygon with n < 1 has no edge to place, so both orientations
// give the top angle rather than dividing by n. An orientation other than FlatTop and PointyTop
// has no meaning and panics.
func RegularPolygonOrientationAngle(n int, orientation Orientation) float64 {
	top := 3 * Pi / 2

	switch orientation {
	case FlatTop:
		if n < 1 {
			return top
		}

		return NormalizeAngle(top - Pi/float64(n))
	case PointyTop:
		return top
	default:
		panic(fmt.Sprintf("geom: unknown orientation %d", orientation))
	}
}

// RegularPolygonWithOrientation creates a RegularPolygon with the given orientation.
func RegularPolygonWithOrientation[T Number](center Point[T], size Size[T], n int, orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{center, size, n, RegularPolygonOrientationAngle(n, orientation)}
}

// Triangle creates a RegularPolygon with 3 vertices.
func Triangle[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 3, orientation)
}

// Square creates a RegularPolygon with 4 vertices.
func Square[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 4, orientation)
}

// Hexagon creates a RegularPolygon with 6 vertices.
func Hexagon[T Number](center Point[T], size Size[T], orientation Orientation) RegularPolygon[T] {
	return RegularPolygonWithOrientation(center, size, 6, orientation)
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(vector), rp.Size, rp.N, rp.Angle}
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
// The stored angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size, rp.N, NormalizeAngle(rp.Angle + angle)}
}

// Vertices returns the polygon vertices in order starting from Angle, by increasing angle —
// the same winding as Directions and Rectangle.Vertices, and clockwise as drawn on a screen
// with Y pointing down. A polygon with N < 1 has no vertices and returns nil, so its Polygon
// is zero like Pol(nil).
// For integer T, each vertex component is rounded to the nearest integer, so vertices at
// non-right angles may be off by up to half a unit. Use float64 for exact positions.
func (rp RegularPolygon[T]) Vertices() []Point[T] {
	if rp.Empty() {
		return nil
	}

	angleStep := (2 * Pi) / float64(rp.N)

	vertices := make([]Point[T], rp.N)
	for i := range vertices {
		vertices[i] = rp.Center.Add(VectorFromAngleSize(rp.Angle+float64(i)*angleStep, rp.Size))
	}

	return vertices
}

// Bounds returns the axis-aligned bounding rectangle computed from the polygon vertices,
// or the zero rectangle for a polygon without vertices, like Polygon.Bounds.
func (rp RegularPolygon[T]) Bounds() Rectangle[T] {
	return rp.Polygon().Bounds()
}

// Polygon converts the regular polygon into a generic Polygon with computed vertices.
func (rp RegularPolygon[T]) Polygon() Polygon[T] {
	return Polygon[T]{rp.Vertices()}
}

// Equal checks if center point, size, number of vertices and angle are equal. Angles are
// compared with EqualAngle, so a full turn or the sign of an angle does not matter.
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool {
	return rp.Center.Equal(polygon.Center) && rp.Size.Equal(polygon.Size) && rp.N == polygon.N && EqualAngle(rp.Angle, polygon.Angle)
}

// IsZero checks if center point, size, number of vertices and angle are zero, comparing
// the angle like Equal so that a full turn counts as zero.
func (rp RegularPolygon[T]) IsZero() bool {
	return rp.Equal(RegularPolygon[T]{})
}

// Empty checks if the polygon has no vertices.
func (rp RegularPolygon[T]) Empty() bool {
	return rp.N < 1
}

// Int converts the regular polygon to a RegularPolygon[int].
func (rp RegularPolygon[T]) Int() RegularPolygon[int] {
	return RegularPolygon[int]{rp.Center.Int(), rp.Size.Int(), rp.N, rp.Angle}
}

// Float converts the regular polygon to a RegularPolygon[float64].
func (rp RegularPolygon[T]) Float() RegularPolygon[float64] {
	return RegularPolygon[float64]{rp.Center.Float(), rp.Size.Float(), rp.N, rp.Angle}
}

// String returns a string representation of the RegularPolygon.
func (rp RegularPolygon[T]) String() string {
	return fmt.Sprintf("RegPol(%s;%s;%s;%s)", rp.Center.String(), rp.Size.String(), String(rp.N), String(rp.Angle))
}

package geom

import (
	"fmt"
	"math"
)

// RegularPolygon is a polygon with equally spaced vertices around a center.
//
// Size holds the semi-axes of the ellipse the vertices lie on, so it is a radius, not an
// extent: a hexagon of Size 10x10 spans 17.32x20, and the polygon inscribed in a circle of
// radius r has Size r x r. This differs from Rectangle, whose Size is the full width and
// height. Use Bounds for the extent.
//
// The size is never negative: RegPol and the orientation constructors take it absolute, and
// Scale takes a negative factor absolute, since a negative semi-axis would place
// every vertex half a turn away rather than describe a different polygon. A negative size can
// only be written as a struct literal or decoded from JSON; Canonical repairs it.
type RegularPolygon[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	N      int      `json:"n"`
	Angle  float64  `json:"a"`
}

// RegPol is shorthand for RegularPolygon{center, size, n, angle}, with the size taken absolute.
func RegPol[T Number](center Point[T], size Size[T], n int, angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{center, size.Abs(), n, angle}
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

// RegularPolygonWithOrientation creates a RegularPolygon with the given orientation, with the
// size taken absolute like RegPol.
func RegularPolygonWithOrientation[T Number](center Point[T], size Size[T], n int, orientation Orientation) RegularPolygon[T] {
	return RegPol(center, size, n, RegularPolygonOrientationAngle(n, orientation))
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

	vertices := make([]Point[T], rp.N)
	for i := range vertices {
		vertices[i] = rp.vertex(i)
	}

	return vertices
}

// Area returns the area enclosed by the polygon, in closed form: n/2 · w · h · sin(2π/n), the
// area of the polygon inscribed in the ellipse with the semi-axes Size holds, without building
// the vertices. A polygon with N < 3 encloses no area, as the sine of a full or a half turn
// gives it. For integer T the vertices are rounded but the area is not: it is the area of the
// exact polygon, which Polygon().Area() measures from the rounded vertices.
func (rp RegularPolygon[T]) Area() float64 {
	if rp.N < 3 {
		return 0
	}

	n := float64(rp.N)

	return n / 2 * float64(rp.Size.Width) * float64(rp.Size.Height) * math.Sin(2*Pi/n)
}

// Perimeter returns the total length of the edges without building the vertices. On a square
// Size it is the closed form 2 · n · r · sin(π/n); on an ellipse the chords differ in length
// and are summed from the angles, since the perimeter of a polygon inscribed in an ellipse has
// no closed form. A polygon with N < 2 has no edge with a length. For integer T it is the
// perimeter of the exact polygon, like Area, not of the rounded vertices.
func (rp RegularPolygon[T]) Perimeter() float64 {
	if rp.N < 2 {
		return 0
	}

	n, size := float64(rp.N), rp.Size.Float()
	if size.Width == size.Height {
		return 2 * n * size.Width * math.Sin(Pi/n)
	}

	perimeter, previous := 0.0, VectorFromAngleSize(rp.Angle, size)
	for i := 1; i <= rp.N; i++ {
		vertex := VectorFromAngleSize(rp.Angle+float64(i)*rp.step(), size)
		perimeter, previous = perimeter+vertex.Subtract(previous).Length(), vertex
	}

	return perimeter
}

// Bounds returns the axis-aligned bounding rectangle of the vertices without building them, or
// the zero rectangle for a polygon without vertices, like Polygon.Bounds. Each side of the box
// is set by the vertex nearest to that direction, at most half a step away, and reads that
// vertex as Vertices places it, so the box is exactly the one Polygon().Bounds() would place,
// rounded alike for an integer T.
func (rp RegularPolygon[T]) Bounds() Rectangle[T] {
	if rp.Empty() {
		return Rectangle[T]{}
	}

	a := Point[T]{rp.vertex(rp.nearest(Pi)).X, rp.vertex(rp.nearest(3 * Pi / 2)).Y}
	b := Point[T]{rp.vertex(rp.nearest(0)).X, rp.vertex(rp.nearest(Pi / 2)).Y}

	return RectangleFromMinMax(a, b)
}

// vertex returns the vertex at the given index, the point Vertices places there: the center
// displaced along the ellipse at Angle plus that many steps, rounded for an integer T.
func (rp RegularPolygon[T]) vertex(i int) Point[T] {
	return rp.Center.Add(VectorFromAngleSize(rp.Angle+float64(i)*rp.step(), rp.Size))
}

// nearest returns the index of the vertex nearest to the given direction, at most half a step
// away, the one that reaches farthest that way. With one or two vertices it can be more than
// a quarter turn off, so the reach is negative exactly where no vertex lies on that side.
func (rp RegularPolygon[T]) nearest(direction float64) int {
	return Mod(int(math.Round((direction-rp.Angle)/rp.step())), rp.N)
}

// step returns the angle between consecutive vertices.
func (rp RegularPolygon[T]) step() float64 {
	return 2 * Pi / float64(rp.N)
}

// Translate creates a new RegularPolygon translated by the given vector.
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center.Add(vector), rp.Size, rp.N, rp.Angle}
}

// MoveTo creates a new RegularPolygon with center at point.
func (rp RegularPolygon[T]) MoveTo(point Point[T]) RegularPolygon[T] {
	return RegularPolygon[T]{point, rp.Size, rp.N, rp.Angle}
}

// Scale creates a new RegularPolygon with size scaled by the given factor. A negative factor
// scales by its absolute value like Rectangle.Scale, since Size holds semi-axes and a negative
// one would place every vertex half a turn away rather than shrink the polygon; use Rotate for
// the half turn.
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Scale(factor).Abs(), rp.N, rp.Angle}
}

// ScaleXY creates a new RegularPolygon with size scaled by the given factors, negative ones by
// their absolute value like Scale.
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.ScaleXY(factorX, factorY).Abs(), rp.N, rp.Angle}
}

// Unscale creates a new RegularPolygon with size scaled by the inverse factor, the inverse of
// Scale and negative factors taken absolute like it. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) Unscale(factor float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Unscale(factor).Abs(), rp.N, rp.Angle}
}

// UnscaleXY creates a new RegularPolygon with size scaled by the inverse of the given factors,
// the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (rp RegularPolygon[T]) UnscaleXY(factorX, factorY float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.UnscaleXY(factorX, factorY).Abs(), rp.N, rp.Angle}
}

// Canonical creates a new RegularPolygon in the form RegPol and Rotate build, with the size
// taken absolute and the angle normalized to [0, 2π): a well-formed polygon is returned as it
// is, up to the full turns Equal already ignores. It repairs a negative semi-axis written as a
// struct literal or decoded from JSON, which would place every vertex half a turn away, and
// brings a decoded angle onto the seam Rotate keeps.
func (rp RegularPolygon[T]) Canonical() RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size.Abs(), rp.N, NormalizeAngle(rp.Angle)}
}

// Lerp creates a new RegularPolygon in linear interpolation towards the given polygon, moving
// the center and the size together and turning the angle along the shorter arc with
// LerpAngle, so a tween from 350° to 10° passes through the top rather than the long way
// round; the result is normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like
// Point.Lerp, with the size taken absolute like RegPol. The vertex count does not
// interpolate: polygons with a different N have no shape between them, so Lerp panics for
// them, the convention RegularPolygonOrientationAngle follows for an orientation with no
// meaning, rather than hiding the mistake in an empty polygon.
func (rp RegularPolygon[T]) Lerp(polygon RegularPolygon[T], t float64) RegularPolygon[T] {
	if rp.N != polygon.N {
		panic(fmt.Sprintf("geom: lerp between polygons of %d and %d vertices", rp.N, polygon.N))
	}

	return RegularPolygon[T]{
		rp.Center.Lerp(polygon.Center, t),
		rp.Size.Lerp(polygon.Size, t).Abs(),
		rp.N,
		NormalizeAngle(LerpAngle(rp.Angle, polygon.Angle, t)),
	}
}

// Rotate creates a new RegularPolygon rotated by the given angle (in radians).
// The stored angle is normalized to [0, 2π) to prevent drift from repeated rotations.
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T] {
	return RegularPolygon[T]{rp.Center, rp.Size, rp.N, NormalizeAngle(rp.Angle + angle)}
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

// Polygon converts the regular polygon into a generic Polygon with computed vertices.
func (rp RegularPolygon[T]) Polygon() Polygon[T] {
	return Polygon[T]{rp.Vertices()}
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

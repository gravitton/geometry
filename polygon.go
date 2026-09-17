package geom

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"

	xslices "github.com/gravitton/x/slices"
)

// Polygon is a 2D polygon given by its vertices. The vertex count is not checked: a polygon with
// fewer than three vertices is degenerate but every method still answers for it.
//
// Vertices is shared, not copied: Pol keeps the slice it is given and every method that
// returns a Polygon allocates a new one. The polygon is immutable as long as its caller does
// not write into that slice. A nil Vertices stays nil through every method, so IsZero holds
// after Translate, Scale, Int or Float.
type Polygon[T Number] struct {
	Vertices []Point[T]
}

// Pol is shorthand for Polygon{vertices}.
func Pol[T Number](vertices []Point[T]) Polygon[T] {
	return Polygon[T]{vertices}
}

// Translate creates a new Polygon translated by the given vector (applied to all vertices).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{xslices.Map(p.Vertices, func(e Point[T]) Point[T] {
		return e.Add(vector)
	})}
}

// MoveTo creates a new Polygon whose centroid is moved to point, preserving shape.
// For integer T the centroid is rounded, and a translation by a whole number of units preserves
// the fractional part of the centroid, so the moved centroid lands on point except when it sits
// exactly on a half and rounding away from zero flips side as the sign changes.
func (p Polygon[T]) MoveTo(point Point[T]) Polygon[T] {
	return p.Translate(point.Subtract(p.Center()))
}

// Scale creates a new Polygon uniformly scaled about its centroid by the factor.
func (p Polygon[T]) Scale(factor float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Center()

	return Polygon[T]{xslices.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Center returns the polygon centroid: the center of the enclosed area, so a vertex added in
// the middle of an edge does not move it. A polygon that encloses no area, with fewer than
// three vertices or all of them collinear, has no such center and falls back to the average of
// its vertices; an empty polygon returns the zero point.
// For integer T the centroid is rounded like every other result stored into T; use float64
// when centroid accuracy matters.
//
// The sum is taken with the origin moved to the first vertex, so the products stay small and
// every edge at that vertex contributes exactly zero: a polygon of one or two vertices has an
// exact zero area whatever its coordinates.
func (p Polygon[T]) Center() Point[T] {
	if p.Empty() {
		return Point[T]{}
	}

	origin := p.Vertices[0].Float()
	offset := origin.Vector().Negate()

	var x, y, twiceArea float64
	for edge := range p.edges() {
		shifted := edge.Float().Translate(offset)
		wedge := shifted.wedge()

		x += (shifted.Start.X + shifted.End.X) * wedge
		y += (shifted.Start.Y + shifted.End.Y) * wedge
		twiceArea += wedge
	}

	if twiceArea == 0 {
		return p.mean()
	}

	centroid := origin.AddXY(x/(3*twiceArea), y/(3*twiceArea))

	return Point[T]{Cast[T](centroid.X), Cast[T](centroid.Y)}
}

// mean returns the average of the vertices, which Center falls back to when the
// polygon encloses no area.
func (p Polygon[T]) mean() Point[T] {
	var x, y float64
	for _, v := range p.Vertices {
		x, y = x+float64(v.X), y+float64(v.Y)
	}

	n := float64(len(p.Vertices))

	return Point[T]{Cast[T](x / n), Cast[T](y / n)}
}

// Edges returns the polygon edges in vertex order, each from a vertex to the next and the
// last one closing back to the first. A single vertex yields one zero-length edge, and a nil
// Vertices maps to nil edges like every other mapping.
func (p Polygon[T]) Edges() []Line[T] {
	if p.IsZero() {
		return nil
	}

	return slices.AppendSeq(make([]Line[T], 0, len(p.Vertices)), p.edges())
}

// edges iterates the edges Edges returns without allocating them, for the methods that only
// need to walk them once.
func (p Polygon[T]) edges() iter.Seq[Line[T]] {
	return func(yield func(Line[T]) bool) {
		n := len(p.Vertices)
		for i, vertex := range p.Vertices {
			if !yield(Line[T]{vertex, p.Vertices[(i+1)%n]}) {
				return
			}
		}
	}
}

// Area returns the area enclosed by the polygon, by the shoelace formula, regardless of winding.
// It is a float64 even for an integer T, since a lattice polygon can enclose half a unit;
// a self-intersecting polygon has its lobes cancel where they wind the opposite way.
func (p Polygon[T]) Area() float64 {
	var twiceArea float64
	for edge := range p.edges() {
		twiceArea += edge.wedge()
	}

	return math.Abs(twiceArea) / 2
}

// Perimeter returns the total length of the edges.
func (p Polygon[T]) Perimeter() float64 {
	var perimeter float64
	for edge := range p.edges() {
		perimeter += edge.Length()
	}

	return perimeter
}

// Bounds returns the axis-aligned bounding rectangle of the vertices, or the zero rectangle
// for a polygon without vertices.
func (p Polygon[T]) Bounds() Rectangle[T] {
	if p.Empty() {
		return Rectangle[T]{}
	}

	return RectangleFromMinMax(p.extent())
}

// extent returns the minimum and maximum corner of the vertices, which must not be empty.
// Bounds rounds them into a Rectangle; Contains tests them as they are.
func (p Polygon[T]) extent() (Point[T], Point[T]) {
	minPoint, maxPoint := p.Vertices[0], p.Vertices[0]
	for _, v := range p.Vertices[1:] {
		minPoint = Point[T]{min(minPoint.X, v.X), min(minPoint.Y, v.Y)}
		maxPoint = Point[T]{max(maxPoint.X, v.X), max(maxPoint.Y, v.Y)}
	}

	return minPoint, maxPoint
}

// Equal checks if two polygons have the same vertices.
func (p Polygon[T]) Equal(polygon Polygon[T]) bool {
	if len(p.Vertices) != len(polygon.Vertices) {
		return false
	}

	for i, v := range p.Vertices {
		if !v.Equal(polygon.Vertices[i]) {
			return false
		}
	}

	return true
}

// IsZero checks if the vertices slice is nil.
func (p Polygon[T]) IsZero() bool {
	return p.Vertices == nil
}

// Empty checks if number of vertices is zero.
func (p Polygon[T]) Empty() bool {
	return len(p.Vertices) == 0
}

// Contains reports whether the given point lies within the polygon, boundary included within
// Epsilon of T, the same closed convention as Rectangle.Contains. The interior follows the
// even-odd rule, so a self-intersecting polygon excludes the regions it winds around twice.
// A point outside the extent of the vertices is rejected before any edge is examined.
func (p Polygon[T]) Contains(point Point[T]) bool {
	if p.Empty() {
		return false
	}

	if minPoint, maxPoint := p.extent(); !between(point, minPoint, maxPoint) {
		return false
	}

	inside := false
	for edge := range p.edges() {
		if edge.Contains(point) {
			return true
		}
		if edge.crossesRay(point) {
			inside = !inside
		}
	}

	return inside
}

// Int converts the polygon to a Polygon[int].
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{xslices.Map(p.Vertices, Point[T].Int)}
}

// Float converts the polygon to a Polygon[float64].
func (p Polygon[T]) Float() Polygon[float64] {
	return Polygon[float64]{xslices.Map(p.Vertices, Point[T].Float)}
}

// String returns a string representation of the Polygon.
func (p Polygon[T]) String() string {
	return fmt.Sprintf("Pol(%s)", strings.Join(xslices.Map(p.Vertices, Point[T].String), ", "))
}

// MarshalJSON implements json.Marshaler.
func (p Polygon[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Vertices)
}

// UnmarshalJSON implements json.Unmarshaler. The vertices are decoded into a fresh slice, so a
// slice the polygon shared before decoding is left untouched.
func (p *Polygon[T]) UnmarshalJSON(bytes []byte) error {
	var vertices []Point[T]
	if err := json.Unmarshal(bytes, &vertices); err != nil {
		return err
	}

	p.Vertices = vertices

	return nil
}

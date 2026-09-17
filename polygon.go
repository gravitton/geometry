package geom

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitton/x/slices"
)

// Polygon is a 2D polygon with 3+ vertices.
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

// Center returns the polygon centroid computed as the average of its vertices,
// or the zero point for a polygon without vertices.
// For integer T, the coordinate sums are divided using integer division and the
// result is truncated; use float64 when centroid accuracy matters.
func (p Polygon[T]) Center() Point[T] {
	if p.Empty() {
		return Point[T]{}
	}

	var x, y T
	for _, v := range p.Vertices {
		x, y = x+v.X, y+v.Y
	}

	l := T(len(p.Vertices))

	return Point[T]{x / l, y / l}
}

// Translate creates a new Polygon translated by the given vector (applied to all vertices).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{slices.Map(p.Vertices, func(e Point[T]) Point[T] {
		return e.Add(vector)
	})}
}

// MoveTo creates a new Polygon whose centroid is moved to point, preserving shape.
// For integer T this is a known exception to exactness: Center truncates toward zero, which does
// not commute with translation, so the moved centroid can miss the point by one unit when a
// coordinate sum crosses zero. Floor division would be exact but is deliberately not used, to
// keep a single integer rounding rule; use float64 when the centroid must land exactly.
func (p Polygon[T]) MoveTo(point Point[T]) Polygon[T] {
	return p.Translate(point.Subtract(p.Center()))
}

// Scale creates a new Polygon uniformly scaled about its centroid by the factor.
func (p Polygon[T]) Scale(factor float64) Polygon[T] {
	center := p.Center()
	return Polygon[T]{slices.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Center()
	return Polygon[T]{slices.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Bounds returns the axis-aligned bounding rectangle of the vertices, or the zero rectangle
// for a polygon without vertices.
func (p Polygon[T]) Bounds() Rectangle[T] {
	if p.Empty() {
		return Rectangle[T]{}
	}

	minPoint, maxPoint := p.Vertices[0], p.Vertices[0]
	for _, v := range p.Vertices[1:] {
		minPoint = Point[T]{min(minPoint.X, v.X), min(minPoint.Y, v.Y)}
		maxPoint = Point[T]{max(maxPoint.X, v.X), max(maxPoint.Y, v.Y)}
	}

	return RectangleFromMinMax(minPoint, maxPoint)
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

// Int converts the polygon to a Polygon[int].
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{slices.Map(p.Vertices, Point[T].Int)}
}

// Float converts the polygon to a Polygon[float64].
func (p Polygon[T]) Float() Polygon[float64] {
	return Polygon[float64]{slices.Map(p.Vertices, Point[T].Float)}
}

// String returns a string representation of the Polygon.
func (p Polygon[T]) String() string {
	return fmt.Sprintf("Pol(%s)", strings.Join(slices.Map(p.Vertices, Point[T].String), ", "))
}

// MarshalJSON implements json.Marshaler.
func (p Polygon[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Vertices)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Polygon[T]) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.Vertices)
}

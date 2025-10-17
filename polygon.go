package geom

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitton/x/slices"
)

// Polygon is a 2D polygon with 3+ vertices.
type Polygon[T Number] struct {
	Vertices []Point[T]
}

// Pol is shorthand for Polygon{vertices}.
func Pol[T Number](Vertices []Point[T]) Polygon[T] {
	return Polygon[T]{Vertices}
}

// Center returns the polygon centroid computed as the average of its vertices.
func (p Polygon[T]) Center() Point[T] {
	var x, y T
	l := T(len(p.Vertices))
	for _, v := range p.Vertices {
		x, y = x+v.X, y+v.Y
	}

	return Point[T]{x / l, y / l}
}

// Translate creates a new Polygon translated by the given vector (applied to all vertices).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{slices.Map(p.Vertices, func(e Point[T]) Point[T] {
		return e.Add(vector)
	})}
}

// MoveTo creates a new Polygon whose centroid is moved to point, preserving shape.
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

// Empty checks if number of vertices is zero.
func (p Polygon[T]) IsZero() bool {
	return p.Vertices == nil
}

// Empty checks if number of vertices is zero.
func (p Polygon[T]) Empty() bool {
	return len(p.Vertices) == 0
}

// Int converts the polygon to a [int] polygon.
func (p Polygon[T]) Int() Polygon[int] {
	return Polygon[int]{slices.Map(p.Vertices, Point[T].Int)}
}

// Float converts the polygon to a [float64] polygon.
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
func (p Polygon[T]) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.Vertices)
}

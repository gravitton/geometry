package geom

import "github.com/gravitton/geometry/internal"

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
	return Polygon[T]{internal.Map(p.Vertices, func(e Point[T]) Point[T] {
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
	return Polygon[T]{internal.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(factor))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T] {
	center := p.Center()
	return Polygon[T]{internal.Map(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(factorX, factorY))
	})}
}

// Empty checks if number of vertices is zero.
func (p Polygon[T]) Empty() bool {
	return len(p.Vertices) == 0
}

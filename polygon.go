package geom

// Polygon is a 2D polygon with 3+ transform.
type Polygon[T Number] struct {
	Vertices []Point[T]
}

// Center returns the polygon centroid computed as the average of its transform.
func (p Polygon[T]) Center() Point[T] {
	var x, y T
	l := T(len(p.Vertices))
	for _, v := range p.Vertices {
		x, y = x+v.X, y+v.Y
	}

	return Point[T]{x / l, y / l}
}

// Translate creates a new Polygon translated by the given vector (applied to all transform).
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T] {
	return Polygon[T]{transform(p.Vertices, func(e Point[T]) Point[T] {
		return e.Add(vector)
	})}
}

// MoveTo creates a new Polygon whose centroid is moved to point, preserving shape.
func (p Polygon[T]) MoveTo(point Point[T]) Polygon[T] {
	return p.Translate(point.Subtract(p.Center()))
}

// Scale creates a new Polygon uniformly scaled about its centroid by the factor.
func (p Polygon[T]) Scale(scale float64) Polygon[T] {
	center := p.Center()
	return Polygon[T]{transform(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).Multiply(scale))
	})}
}

// ScaleXY creates a new Polygon scaled about its centroid by the factors.
func (p Polygon[T]) ScaleXY(scaleX, scaleY float64) Polygon[T] {
	center := p.Center()
	return Polygon[T]{transform(p.Vertices, func(point Point[T]) Point[T] {
		return center.Add(point.Subtract(center).MultiplyXY(scaleX, scaleY))
	})}
}

func transform[S ~[]E, E any, T any](input S, fn func(E) T) []T {
	output := make([]T, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}

	return output
}

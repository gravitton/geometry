package geom

import "iter"

// Shape is what every shape answers about a point: Line, Rectangle, Circle, Polygon and
// RegularPolygon. Bounds is the axis-aligned box around it, Contains includes the boundary
// within Epsilon of T, DistanceTo is zero exactly where Contains holds and DistanceSquaredTo is
// the value it takes the root of. A spatial index or a picking routine holds a Shape and never
// needs to know which one.
type Shape[T Number] interface {
	Bounds() Rectangle[T]
	Contains(point Point[T]) bool
	DistanceTo(point Point[T]) float64
	DistanceSquaredTo(point Point[T]) float64
}

// Outline is a shape whose boundary is a chain of straight edges: Line, Rectangle, Polygon and
// RegularPolygon. Vertices iterates the corners in order and Edges the segments joining them,
// each edge starting where the previous one ends, the last one closing back to the first on a
// closed shape. Both iterate without allocating, so one loop draws or measures any of the four.
// A Circle has no vertices; convert it through a RegularPolygon of the wanted resolution.
//
// The iterators are free only on a concrete shape, where the compiler inlines them into the
// loop. Called through an Outline value, or through a type parameter constrained by it, the
// iterator and the loop body are heap allocated, since neither call is inlined. A hot loop
// keeps the concrete type and switches on it; Outline is for the code around it.
type Outline[T Number] interface {
	Vertices() iter.Seq[Point[T]]
	Edges() iter.Seq[Line[T]]
}

// Collider is a shape tested against every shape: Line, Rectangle, Circle, Polygon and
// RegularPolygon, each with the five Intersects methods, its own kind included. Every test is
// symmetric and includes a touch within Epsilon of T, so a broad collision pass calls the
// method for the other side's kind and gets the same answer from either.
type Collider[T Number] interface {
	IntersectsLine(line Line[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsCircle(circle Circle[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
}

// Measured is a shape with an area and a perimeter, in the number type M that measures them:
// Size and Rectangle measure in their own T, since a box of an integer size has an exact
// integer area and perimeter, while Circle, Polygon and RegularPolygon measure in float64,
// since a curve or a turned edge encloses what no integer expresses. Every float64 shape is a
// Measured[float64]; an integer Rectangle or Size is a Measured[int].
type Measured[M Number] interface {
	Area() M
	Perimeter() M
}

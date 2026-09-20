package geom

import "iter"

// Shape is what every shape answers about a point: Segment, Rectangle, Circle, Ellipse,
// Polygon and RegularPolygon. Bounds is the axis-aligned box around it, Contains includes the boundary
// within Epsilon of T, DistanceTo is zero exactly where Contains holds and DistanceSquaredTo is
// the value it takes the root of. A spatial index or a picking routine holds a Shape and never
// needs to know which one.
type Shape[T Number] interface {
	Bounds() Rectangle[T]
	Contains(point Point[T]) bool
	DistanceTo(point Point[T]) float64
	DistanceSquaredTo(point Point[T]) float64
}

// Outline is a shape whose boundary is a chain of straight edges: Segment, Rectangle, Polygon and
// RegularPolygon. Vertices iterates the corners in order and Edges the segments joining them,
// each edge starting where the previous one ends, the last one closing back to the first on a
// closed shape. Both iterate without allocating, so one loop draws or measures any of the four.
// A Circle and an Ellipse have no vertices; convert one through a RegularPolygon of the
// wanted resolution.
//
// The iterators are free only on a concrete shape, where the compiler inlines them into the
// loop. Called through an Outline value, or through a type parameter constrained by it, the
// iterator and the loop body are heap allocated, since neither call is inlined. A hot loop
// keeps the concrete type and switches on it; Outline is for the code around it.
type Outline[T Number] interface {
	Vertices() iter.Seq[Point[T]]
	Edges() iter.Seq[Segment[T]]
}

// Collider is a shape tested against every shape: Segment, Rectangle, Circle, Polygon and
// RegularPolygon, each with the five Intersects methods, its own kind included. Every test is
// symmetric and includes a touch within Epsilon of T, so a broad collision pass calls the
// method for the other side's kind and gets the same answer from either.
//
// Ellipse is deliberately not one: two ellipses meet at the roots of a quartic, which none of
// the closed forms the circle pairs are built on reaches. Test an ellipse through the
// RegularPolygon of the wanted resolution until the pairs land.
type Collider[T Number] interface {
	IntersectsSegment(segment Segment[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsCircle(circle Circle[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
}

// Measured is a shape with an area and a perimeter, in the number type M that measures them:
// Size and Rectangle measure in their own T, since a box of an integer size has an exact
// integer area and perimeter, while Circle, Ellipse, Polygon and RegularPolygon measure in float64,
// since a curve or a turned edge encloses what no integer expresses. Every float64 shape is a
// Measured[float64]; an integer Rectangle or Size is a Measured[int].
type Measured[M Number] interface {
	Area() M
	Perimeter() M
}

// Movable is a shape that can be moved, turned and scaled, in its own type: Segment, Rectangle, Circle,
// Ellipse, Polygon and RegularPolygon, each returning itself rather than a common type, so the parameter
// S stands for the shape and every method returns it. It is a constraint, not a value type,
// and is written self-referentially at the call site:
//
//	func Tween[T Number, S Movable[T, S]](shape S, to Point[T], t float64) S
//
// Every shape turns about its own center, Circle.Rotate giving the circle back. Lerp is
// deliberately absent: Segment.Lerp is the point a fraction along the segment rather than
// a step toward another segment, and Polygon has none, so the six shapes do not share it.
// A call through a type parameter constrained by an interface allocates, so this is for the
// code around a hot loop, never inside one.
type Movable[T Number, S any] interface {
	Translate(vector Vector[T]) S
	MoveTo(point Point[T]) S
	Rotate(angle float64) S
	Scale(factor float64) S
	Unscale(factor float64) S
}

// Intersects reports whether two shapes held as Collider share a point, by the method of the
// first naming the kind of the second: the one call a broad collision pass makes over a mixed
// list, where the concrete types are not known until it runs. It is symmetric and includes a
// touch within Epsilon of T, like every method it dispatches to.
//
// Both arguments escape to the heap, so a loop that knows its types calls the named method
// directly and allocates nothing. A shape of another package satisfying Collider works on
// either side, since the other side names its kind; two of them together have no method this
// package can reach and panic.
func Intersects[T Number](a, b Collider[T]) bool {
	if result, ok := intersectsKind(a, b); ok {
		return result
	}

	if result, ok := intersectsKind(b, a); ok {
		return result
	}

	panic("geom: intersects needs a shape of this package on one side")
}

// intersectsKind tests a against b by the method of a that names the kind of b, and false where
// b is none of the five shapes, which leaves the pair to the call with the arguments swapped.
func intersectsKind[T Number](a, b Collider[T]) (bool, bool) {
	switch shape := b.(type) {
	case Circle[T]:
		return a.IntersectsCircle(shape), true
	case Segment[T]:
		return a.IntersectsSegment(shape), true
	case Polygon[T]:
		return a.IntersectsPolygon(shape), true
	case Rectangle[T]:
		return a.IntersectsRectangle(shape), true
	case RegularPolygon[T]:
		return a.IntersectsRegularPolygon(shape), true
	}

	return false, false
}

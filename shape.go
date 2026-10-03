package geom

// Shape is what every shape answers about a point: Segment, Rectangle, Circle, Ellipse,
// Polygon, RegularPolygon and Box. Bounds is the Box around it, Contains includes the boundary
// within the tolerance, DistanceTo is zero exactly where Contains holds, DistanceSquaredTo is
// the value it takes the root of, and Nearest is the point it is measured to, the point itself
// exactly where Contains holds. A spatial index or a picking routine holds a Shape and never
// needs to know which one.
type Shape[T Number] interface {
	Bounds() Box[T]
	Contains(point Point[T]) bool
	DistanceTo(point Point[T]) float64
	DistanceSquaredTo(point Point[T]) float64
	Nearest(point Point[T]) Point[T]
}

// Collider is a shape tested against every shape: Circle, Segment, Ray, Polygon, Rectangle,
// RegularPolygon and Box, each with the seven Intersects methods, its own kind included. Every
// test is symmetric and includes a touch within the tolerance, so a broad collision pass calls
// the method for the other side's kind and gets the same answer from either. Ray is the one
// Collider that is no Shape, since it has no Bounds.
//
// Ellipse is deliberately not one: two ellipses meet at the roots of a quartic, which none of
// the closed forms the circle pairs are built on reaches. Test an ellipse as its
// RegularPolygon of the wanted resolution.
type Collider[T Number] interface {
	IntersectsCircle(circle Circle[T]) bool
	IntersectsSegment(segment Segment[T]) bool
	IntersectsRay(ray Ray[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
	IntersectsBox(box Box[T]) bool
}

// Body is a shape with an area, the mass properties a physics engine takes from it: Rectangle,
// Circle, Ellipse, Polygon and RegularPolygon. Area is the mass at unit density, Centroid the
// center of mass and Inertia the polar second moment of area about it, the rotational inertia
// at unit density; a physics body multiplies Area and Inertia by its density and never asks
// which shape it holds. Segment is not a Body, since it encloses no area, and Size is not a shape.
type Body[T Number] interface {
	Area() float64
	Centroid() Point[T]
	Inertia() float64
}

// Intersects reports whether two shapes held as Collider share a point, by the method of the
// first naming the kind of the second: the one call a broad collision pass makes over a mixed
// list, where the concrete types are not known until it runs. It is symmetric and includes a
// touch within the tolerance, like every method it dispatches to.
//
// Both arguments escape to the heap, so a loop that knows its types calls the named method
// directly and allocates nothing. The shapes of this package are named by value, so a pointer
// to one is a shape of another package here. A shape of another package satisfying Collider
// works on either side, since the other side names its kind; two of them together have no
// method this package can reach and panic.
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
// b is none of the seven shapes, which leaves the pair to the call with the arguments swapped.
func intersectsKind[T Number](a, b Collider[T]) (bool, bool) {
	switch shape := b.(type) {
	case Circle[T]:
		return a.IntersectsCircle(shape), true
	case Segment[T]:
		return a.IntersectsSegment(shape), true
	case Ray[T]:
		return a.IntersectsRay(shape), true
	case Polygon[T]:
		return a.IntersectsPolygon(shape), true
	case Rectangle[T]:
		return a.IntersectsRectangle(shape), true
	case RegularPolygon[T]:
		return a.IntersectsRegularPolygon(shape), true
	case Box[T]:
		return a.IntersectsBox(shape), true
	}

	return false, false
}

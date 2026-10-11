package geom

import (
	"math"
	"slices"
)

// Shape is what every shape answers: Circle, Segment, Polygon, Rectangle, RegularPolygon, Box
// and Ellipse. Bounds is the Box around it, Contains includes the boundary within
// the tolerance, DistanceTo is zero exactly where Contains holds, DistanceSquaredTo is the
// value it takes the root of, and Nearest is the point it is measured to, the point itself
// exactly where Contains holds. Support is the point of the shape farthest along a direction
// of any length, in float64, any one of them where several reach alike.
//
// A spatial index, a picking routine or a collision pass holds a Shape and never needs to know
// which one: Intersects, Distance, Nearest, Collide, Sweep and Raycast take any two. Two kinds of this
// package with a method for each other are tested by that method. Any other pair, an Ellipse
// or a shape of another package among them, is known by its Support alone, which describes a
// shape with no notch: such a pair is measured in float64 and closed on within the tolerance,
// so two of them touching within it may answer either way, and for an integer T only outlines
// meeting on the lattice touch at all.
//
// Ray and Line run without end, so they have no Bounds and no Support and are no Shape: they
// are what a shape is tested by, through their own methods naming its kind, and Raycast.
type Shape[T Number] interface {
	Bounds() Box[T]
	DistanceTo(point Point[T]) float64
	DistanceSquaredTo(point Point[T]) float64
	Nearest(point Point[T]) Point[T]
	Contains(point Point[T]) bool
	Support(direction Vector[float64]) Point[float64]
}

// Body is a shape with an area, the mass properties a physics engine takes from it: Circle,
// Polygon, Rectangle, RegularPolygon and Ellipse. Area is the mass at unit density, Centroid the
// center of mass and Inertia the polar second moment of area about it, the rotational inertia
// at unit density; a physics body multiplies Area and Inertia by its density and never asks
// which shape it holds. Segment is not a Body, since it encloses no area, and Size is not a
// shape.
type Body[T Number] interface {
	Area() float64
	Centroid() Point[T]
	Inertia() float64
}

// Intersects reports whether two shapes share a point: the one call a broad collision pass
// makes over a mixed list, where the concrete types are not known until it runs. It is
// symmetric and includes a touch within the tolerance.
//
// Two kinds with an Intersects method for each other are tested by that method, the one of the
// first naming the kind of the second, so the answer is the one the method gives. Any other
// pair is tested from what the two are: an empty polygon intersects nothing; a Circle is tested
// by its radius on the distance of the other to its center, the comparison every circle pair
// makes; a Polygon that is not convex edge by edge, or one within the other, so a notch is not
// filled in; and two others by their supports. That is how an Ellipse, and a shape of another
// package, is tested.
//
// Both arguments escape to the heap, so a loop that knows its types calls the named method
// directly and allocates nothing. The shapes of this package are named by value, so a pointer
// to one is a shape of another package here.
func Intersects[T Number](a, b Shape[T]) bool {
	if shape, ok := a.(intersector[T]); ok {
		switch kind := b.(type) {
		case Circle[T]:
			return shape.IntersectsCircle(kind)
		case Segment[T]:
			return shape.IntersectsSegment(kind)
		case Polygon[T]:
			return shape.IntersectsPolygon(kind)
		case Rectangle[T]:
			return shape.IntersectsRectangle(kind)
		case RegularPolygon[T]:
			return shape.IntersectsRegularPolygon(kind)
		case Box[T]:
			return shape.IntersectsBox(kind)
		}
	}

	if isEmpty(a) || isEmpty(b) {
		return false
	}

	if circle, other, ok := eitherCircle(a, b); ok {
		return circle.containsSquared(other.DistanceSquaredTo(circle.Center), other.Bounds().magnitude())
	}

	if polygon, other, ok := eitherNonConvex(a, b); ok {
		return polygon.intersectsSupport(other)
	}

	return minkowski[T]{a: a, b: b}.intersects()
}

// Distance returns the distance between the nearest points of two shapes: zero exactly where
// Intersects holds, a touch within the tolerance included, and otherwise the gap between the
// two.
func Distance[T Number](a, b Shape[T]) float64 {
	return math.Sqrt(DistanceSquared(a, b))
}

// DistanceSquared returns the squared distance Distance takes the root of, for comparisons:
// zero exactly where Intersects holds, and otherwise the squared gap between the two, measured
// by what they are.
//
// Two kinds with a distance of their own are measured by it, which is zero where they
// intersect: two circles by Circle.DistanceSquaredToCircle, two segments by
// Segment.DistanceSquaredToSegment, and two boxes by Box.DistanceSquaredToBox. Any other pair is
// asked Intersects first. A Circle is then the distance of the other to its center less its
// radius. A Polygon that is not convex is the least over its edges, so a shape in a notch is
// measured to the notch, and an empty polygon is infinitely far from everything. Any other two
// are measured by their supports, in float64: two outlines to the rounding of the arithmetic,
// and a round shape to within Delta of the distance.
//
// Both arguments escape to the heap, as in Intersects, and a polygon that is not convex
// allocates for every edge it is measured by.
func DistanceSquared[T Number](a, b Shape[T]) float64 {
	switch kind := a.(type) {
	case Circle[T]:
		if other, ok := b.(Circle[T]); ok {
			return kind.DistanceSquaredToCircle(other)
		}
	case Segment[T]:
		if other, ok := b.(Segment[T]); ok {
			return kind.DistanceSquaredToSegment(other)
		}
	case Box[T]:
		if other, ok := b.(Box[T]); ok {
			return kind.DistanceSquaredToBox(other)
		}
	}

	if Intersects(a, b) {
		return 0
	}

	if isEmpty(a) || isEmpty(b) {
		return math.Inf(1)
	}

	if circle, other, ok := eitherCircle(a, b); ok {
		gap := other.DistanceTo(circle.Center) - float64(circle.Radius)

		return gap * gap
	}

	if polygon, other, ok := eitherNonConvex(a, b); ok {
		gap := math.Inf(1)
		for edge := range polygon.Edges() {
			gap = min(gap, DistanceSquared[T](edge, other))
		}

		return gap
	}

	return minkowski[T]{a: a, b: b}.nearest().point.LengthSquared()
}

// Nearest returns the two points Distance is measured between: the point of the first shape
// nearest the second and the point of the second nearest the first. Where the two intersect,
// both are one point they share, any one of them. Where the nearest points are not single,
// two flat sides facing each other, one pair of them is returned, which may be another with
// the shapes given the other way round.
//
// The pair is found as DistanceSquared measures it. A Circle gives the point of the other
// nearest its center and its own point nearest that one. A Polygon that is not convex gives
// the pair of its nearest edge, so a shape in a notch is measured to the notch. Any other two are
// measured by their supports, in float64: two outlines to the rounding of the arithmetic, and
// a round shape to within the tolerance of the distance. For integer T the points are rounded
// once and can land off the shapes. An empty polygon has no point, and both are the zero
// point.
//
// Both arguments escape to the heap, as in Intersects, and a polygon that is not convex
// allocates for every edge it is measured by.
func Nearest[T Number](a, b Shape[T]) (Point[T], Point[T]) {
	if isEmpty(a) || isEmpty(b) {
		return Point[T]{}, Point[T]{}
	}

	intersects := Intersects(a, b)

	if circle, other, ok := eitherCircle(a, b); ok {
		point := other.Nearest(circle.Center)

		own := point
		if !intersects {
			own = circle.Nearest(point)
		}

		if _, first := a.(Circle[T]); first {
			return own, point
		}

		return point, own
	}

	if polygon, ok := nonConvex(a); ok {
		return polygon.nearestSupport(b, intersects)
	}

	if polygon, ok := nonConvex(b); ok {
		second, first := polygon.nearestSupport(a, intersects)

		return first, second
	}

	nearest := minkowski[T]{a: a, b: b}.nearest()
	if intersects {
		point := nearest.from.Cast[T]()

		return point, point
	}

	return nearest.from.Add(nearest.point).Cast[T](), nearest.from.Cast[T]()
}

// Collide returns how two overlapping shapes meet, and false where they do not intersect: the
// unit normal that parts them, from the second shape toward the first, the depth along it, and
// the points of the second shape they meet at, which a collision response moves the two apart
// by and a physics solver applies its impulses at. It is false exactly where Intersects is, and
// two shapes that only touch have a normal and no depth.
//
// The normal and the depth are measured on the supports of the two, in float64, to within the
// tolerance. Round shapes are closed on rather than solved, and two shapes one of which lies
// deep within the other have many ways out of nearly one length, of which one is returned.
// The points are read off the sides the two turn to each other across the normal: where both
// are flat, a box resting on the ground, the two ends of the stretch they share, so the box
// does not rock, and one point where either is a corner or round. A side within a thousandth
// of a radian of facing the normal counts as flat against it.
//
// A Polygon that is not convex meets edge by edge, so a shape in a notch is parted from the
// sides of the notch and not from the hull: the contact is the deepest of those with the edges
// the other shape crosses, each taken on the outer side of its edge, so one contact parts the
// two from one edge and may leave them overlapping at another. A shape pressed into an inside
// corner parts from the deeper wall, and colliding again once it is moved gives the other. A
// shape reaching through a thin part of the polygon leaves by the side it is mostly on, one
// mostly within the polygon by the edge it has the least to pass, and one wholly within it,
// crossing no edge, parts from the hull. Of two such polygons the one whose vertices come
// first by Point.Compare is the one walked, and each edge crossed allocates.
//
// Given the other way round, the two answer with the opposite normal and the same depth. Where
// several ways out are as short, a box overlapping the corner of another alike on both axes,
// the one nearest the X axis is taken; a shape in the very middle of another has two opposite
// ways out as short, and either may be returned.
func Collide[T Number](a, b Shape[T]) (Contact[T], bool) {
	if !Intersects(a, b) {
		return Contact[T]{}, false
	}

	first, walkFirst := nonConvex(a)
	second, walkSecond := nonConvex(b)
	if walkFirst && walkSecond {
		walkFirst = slices.CompareFunc(first.Points, second.Points, Point[T].Compare) <= 0
	}

	var contact Contact[T]

	found := false
	switch {
	case walkFirst:
		contact, found = first.collideSupport(b, false)
	case walkSecond:
		contact, found = second.collideSupport(a, true)
	}

	if found {
		return contact, true
	}

	difference := minkowski[T]{a: a, b: b}
	normal, depth := difference.depth()

	return Contact[T]{normal, depth, difference.contact(normal)}, true
}

// Sweep returns where the first shape first touches the second while it moves by the given
// motion, and false where it reaches nothing on the way: the Point of the second shape it
// touches, the unit Normal there, turned to face the moving shape, and the Fraction of the
// motion made by then, in [0, 1]. A body moving farther in a step than it is wide passes
// through a wall between two Intersects tests, and moving it by the fraction of the motion
// instead stops it at the wall. Two shapes that intersect before the move hit nothing, as a
// cast from inside a shape does; Collide parts them. Neither shape turns on the way.
//
// It is found by conservative advancement: the moving shape is moved by the gap between the two
// over the speed they close at, which can never pass the other, and no farther than the whole
// motion, until the gap is within the tolerance, so the fraction is that close short of the
// touch, and a sweep along a side that never closes on it reaches nothing. Where two sides
// meet flat, the point is one point of the contact. For integer T the point is rounded like
// every result and the tolerance is that of float64. A Polygon that is not convex is swept
// edge by edge, so a shape moves into a notch and stops at the first edge it reaches there,
// and an empty polygon reaches nothing and is reached by nothing.
func Sweep[T Number](a Shape[T], motion Vector[T], b Shape[T]) (Hit[T], bool) {
	var first Hit[T]

	found := false
	if Intersects(a, b) || isEmpty(a) || isEmpty(b) {
		return first, false
	}

	if polygon, ok := nonConvex(a); ok {
		for edge := range polygon.Edges() {
			if hit, ok := Sweep[T](edge, motion, b); ok && (!found || hit.Fraction < first.Fraction) {
				first, found = hit, true
			}
		}

		return first, found
	}

	if polygon, ok := nonConvex(b); ok {
		for edge := range polygon.Edges() {
			if hit, ok := Sweep[T](a, motion, edge); ok && (!found || hit.Fraction < first.Fraction) {
				first, found = hit, true
			}
		}

		return first, found
	}

	return minkowski[T]{a: a, b: b}.sweep(motion.Float())
}

// Raycast returns where the ray first reaches the shape from outside, and false where it does
// not, by the Hit method of the ray naming the kind of the shape: the one call a ray makes
// against a mixed list, as Intersects is for a test. The Fraction is in units of Direction, so
// a ray along the Vector of a segment reaches no farther than its End where the fraction is at
// most one. An Origin the shape contains hits nothing.
//
// An Ellipse, and a shape of another package, has no such method and is reached as Sweep
// reaches a shape, by the origin moved along the ray past the Bounds of the shape, to within
// the tolerance and with a touch in passing left out. Neither way allocates.
func Raycast[T Number](ray Ray[T], shape Shape[T]) (Hit[T], bool) {
	switch kind := shape.(type) {
	case Circle[T]:
		return ray.HitCircle(kind)
	case Segment[T]:
		return ray.HitSegment(kind)
	case Polygon[T]:
		return ray.HitPolygon(kind)
	case Rectangle[T]:
		return ray.HitRectangle(kind)
	case RegularPolygon[T]:
		return ray.HitRegularPolygon(kind)
	case Box[T]:
		return ray.HitBox(kind)
	}

	if shape.Contains(ray.Origin) {
		return Hit[T]{}, false
	}

	bounds := shape.Bounds()
	origin := minkowski[T]{b: shape, offset: ray.Origin.Float().Vector()}

	return ray.rescale(origin.sweep(ray.reach(bounds.Min, bounds.Max).Vector().Float()))
}

// intersector is a shape with the Intersects method for each of the six kinds of shape the
// package names: Circle, Segment, Polygon, Rectangle, RegularPolygon and Box, each of which is
// one. Intersects tests two shapes through it where the first is one; Ellipse is not, and is
// tested by its Support.
type intersector[T Number] interface {
	IntersectsCircle(circle Circle[T]) bool
	IntersectsSegment(segment Segment[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
	IntersectsBox(box Box[T]) bool
}

// isEmpty reports whether the shape is a polygon or a regular polygon of no vertices, which
// is nowhere: it intersects nothing and is infinitely far from everything, where its Support,
// the zero point or its center, would stand in for it.
func isEmpty[T Number](shape Shape[T]) bool {
	switch kind := shape.(type) {
	case Polygon[T]:
		return kind.IsEmpty()
	case RegularPolygon[T]:
		return kind.IsEmpty()
	}

	return false
}

// eitherCircle returns whichever of two shapes is a circle, the first where both are, with the
// other one, and false where neither is: the shape tested by its own radius.
func eitherCircle[T Number](a, b Shape[T]) (Circle[T], Shape[T], bool) {
	if circle, ok := a.(Circle[T]); ok {
		return circle, b, true
	}

	circle, ok := b.(Circle[T])

	return circle, a, ok
}

// eitherNonConvex returns whichever of two shapes is a polygon that is not convex, the first
// where both are, with the other one, and false where neither is.
func eitherNonConvex[T Number](a, b Shape[T]) (Polygon[T], Shape[T], bool) {
	if polygon, ok := nonConvex(a); ok {
		return polygon, b, true
	}

	polygon, ok := nonConvex(b)

	return polygon, a, ok
}

// nonConvex returns the shape as a polygon that is not convex by IsConvex, one with a notch or
// with too few vertices to enclose an area, and false where it is anything else: the shape
// taken edge by edge, where its Support, that of its convex hull, would fill it in.
func nonConvex[T Number](shape Shape[T]) (Polygon[T], bool) {
	polygon, ok := shape.(Polygon[T])

	return polygon, ok && !polygon.IsConvex()
}

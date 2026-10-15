package geom

import (
	"math"
	"slices"
)

// Shape is a bounded region: Circle, Capsule, Sector, Segment, Polygon, Rectangle, RegularPolygon,
// Box and Ellipse. Bounds is the axis-aligned bounding box, Contains is closed, DistanceTo is zero
// exactly where Contains holds and DistanceSquaredTo is its square, Nearest is the point DistanceTo
// measures to, the point itself inside, and Support is a point farthest along a direction of any
// length.
//
// Contains, Intersects, Distance, Nearest, Collide and Sweep take any two shapes, and Raycast a ray
// and one. A pair of this package's kinds with a named method is decided by it, and Intersects,
// Distance and Nearest take a Circle against any other shape by the distance of that shape to its
// center. Any other pair, an Ellipse or a shape of another package among them, is known by its
// Support, which sees a shape as convex, so a non-convex Polygon is taken edge by edge and a Sector
// past half a turn as its two halves. Such a pair is measured in float64, two shapes touching
// within the tolerance may answer either way, and for an integer T only outlines meeting on the
// lattice touch.
//
// Ray and Line have no Bounds and no Support, so they are no Shape: they test shapes by their own
// methods and through Raycast.
type Shape[T Number] interface {
	Bounds() Box[T]
	DistanceTo(point Point[T]) float64
	DistanceSquaredTo(point Point[T]) float64
	Nearest(point Point[T]) Point[T]
	Contains(point Point[T]) bool
	Support(direction Vector[float64]) Point[float64]
}

// Distance returns the distance between the nearest points of the shapes, zero where Intersects
// holds.
func Distance[T Number](a, b Shape[T]) float64 {
	return math.Sqrt(DistanceSquared(a, b))
}

// DistanceSquared returns the square of Distance.
//
// Two of a kind with a distance of their own are measured by it: Circle, Capsule, Sector, Segment
// and Box. A Circle against any other shape is the distance from its center less its radius, a
// non-convex Polygon the least over its edges, so a shape in a notch is measured to the notch, and
// a Sector past half a turn the nearer of its halves. An empty polygon is infinitely far from
// everything. The other pairs are measured by their supports: two outlines to the rounding of the
// arithmetic, and a round shape to within the tolerance.
//
// Both arguments escape to the heap, and a non-convex polygon allocates for every edge it is
// measured by.
func DistanceSquared[T Number](a, b Shape[T]) float64 {
	switch kind := a.(type) {
	case Circle[T]:
		if other, ok := b.(Circle[T]); ok {
			return kind.DistanceSquaredToCircle(other)
		}
	case Capsule[T]:
		if other, ok := b.(Capsule[T]); ok {
			return kind.DistanceSquaredToCapsule(other)
		}
	case Sector[T]:
		if other, ok := b.(Sector[T]); ok {
			return kind.DistanceSquaredToSector(other)
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
		distanceSquared, radius := other.DistanceSquaredTo(circle.Center), float64(circle.Radius)
		gap := (distanceSquared - float64(radius*radius)) / (math.Sqrt(distanceSquared) + radius)

		return gap * gap
	}

	if polygon, other, ok := eitherNonConvex(a, b); ok {
		gap := math.Inf(1)
		for edge := range polygon.Edges() {
			gap = min(gap, DistanceSquared[T](edge, other))
		}

		return gap
	}

	if sector, other, ok := eitherReflex(a, b); ok {
		halves := sector.halves()

		return min(DistanceSquared[T](halves[0], other), DistanceSquared[T](halves[1], other))
	}

	difference := minkowski[T]{a: a, b: b}

	return difference.nearest(difference.precision()).point.LengthSquared()
}

// Nearest returns the points Distance is measured between: the point of the first shape nearest the
// second, then the point of the second nearest the first. Shapes that intersect give one shared
// point twice, and where the nearest points are not single, as with two parallel sides, one pair is
// returned, which may differ with the shapes swapped. An empty polygon gives two zero points.
//
// Apart from a Circle, the pair is found by supports to within the tolerance, so it can differ by
// that much from a Distance measured in closed form. Both arguments escape to the heap, and a
// non-convex polygon allocates for every edge it is measured by.
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

	if sector, ok := reflex(a); ok {
		halves := sector.halves()
		if DistanceSquared[T](halves[1], b) < DistanceSquared[T](halves[0], b) {
			return Nearest[T](halves[1], b)
		}

		return Nearest[T](halves[0], b)
	}

	if _, ok := reflex(b); ok {
		second, first := Nearest(b, a)

		return first, second
	}

	difference := minkowski[T]{a: a, b: b}
	nearest := difference.nearest(difference.precision())
	if intersects {
		point := nearest.from.Cast[T]()

		return point, point
	}

	return nearest.from.Add(nearest.point).Cast[T](), nearest.from.Cast[T]()
}

// Contains reports whether the inner shape lies entirely in the outer one, and false where either
// is an empty polygon.
//
// Where the outer shape is a kind of this package with a Contains method naming the kind of the
// inner one, it decides. Any other pair, an Ellipse or a shape of another package on either side,
// is decided by the supports, so a shape touching the other from inside may answer either way. A round shape that nears the
// boundary of the other all the way round is taken not to lie in it: an Ellipse holds one Equal to
// it, but not one a little smaller or the same ellipse turned by a half turn, and shapes clear of
// each other by a small part of their size are decided correctly. A shape in the notch of a
// non-convex Polygon or of a Sector past half a turn is not held.
//
// Both arguments escape to the heap.
func Contains[T Number](outer, inner Shape[T]) bool {
	if shape, ok := outer.(container[T]); ok && isKind(outer) {
		switch kind := inner.(type) {
		case Circle[T]:
			return shape.ContainsCircle(kind)
		case Capsule[T]:
			return shape.ContainsCapsule(kind)
		case Sector[T]:
			return shape.ContainsSector(kind)
		case Segment[T]:
			return shape.ContainsSegment(kind)
		case Polygon[T]:
			return shape.ContainsPolygon(kind)
		case Rectangle[T]:
			return shape.ContainsRectangle(kind)
		case RegularPolygon[T]:
			return shape.ContainsRegularPolygon(kind)
		case Box[T]:
			return shape.ContainsBox(kind)
		}
	}

	if isEmpty(outer) || isEmpty(inner) {
		return false
	}

	if ellipse, ok := outer.(Ellipse[T]); ok {
		if other, ok := inner.(Ellipse[T]); ok && ellipse.Equal(other) {
			return true
		}
	}

	if polygon, ok := nonConvex(outer); ok {
		return polygon.containsSupport(inner)
	}

	if sector, ok := reflex(outer); ok {
		return sector.containsSupport(inner)
	}

	return containsSupport(outer, inner)
}

// Intersects reports whether the shapes share a point, the same whichever is given first.
//
// Two kinds of this package with an Intersects method for each other are decided by it. Any other
// pair is decided from what the two are: an empty polygon intersects nothing, a Circle by its
// radius on the distance of the other to its center, a non-convex Polygon edge by edge or one
// within the other, so a notch is not filled in, a Sector past half a turn as its two halves, and
// any other two by their supports.
//
// Both arguments escape to the heap, so a loop that knows its types calls the named method instead.
// A pointer to a shape of this package counts as a shape of another package.
func Intersects[T Number](a, b Shape[T]) bool {
	if shape, ok := a.(intersector[T]); ok && isKind(a) {
		switch kind := b.(type) {
		case Circle[T]:
			return shape.IntersectsCircle(kind)
		case Capsule[T]:
			return shape.IntersectsCapsule(kind)
		case Sector[T]:
			return shape.IntersectsSector(kind)
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

	if segment, ellipse, ok := eitherSegmentEllipse(a, b); ok {
		return segment.IntersectsEllipse(ellipse)
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

	if sector, other, ok := eitherReflex(a, b); ok {
		halves := sector.halves()

		return Intersects[T](halves[0], other) || Intersects[T](halves[1], other)
	}

	return minkowski[T]{a: a, b: b}.intersects()
}

// Collide returns how two intersecting shapes meet, and false where Intersects does not hold: the
// unit normal from the second shape toward the first, the depth along it, and the points of the
// second shape they meet at. Shapes that only touch have a normal and no depth, and two single
// points have neither. The points are the two ends of the shared stretch where both sides are flat,
// a side within a thousandth of a radian of facing the normal counting as flat, and one point
// otherwise. Swapping the shapes gives the opposite normal and the same depth.
//
// The normal and depth are measured on the supports to within the tolerance. A shape deep within
// another has many ways out of nearly one length, and one is returned: of equally short ways the
// one nearest the X axis, and in the very middle of a shape either of two opposite ones.
//
// A non-convex Polygon meets edge by edge, so a shape in a notch is parted from the sides of the
// notch: the contact is the deepest with the edges the other crosses, which may leave the two
// overlapping at another edge, and a shape inside crossing no edge parts from the hull. Of two such
// polygons the one first by Point.Compare is walked, and each edge crossed allocates. A Sector past
// half a turn meets as its two halves, the deeper of the two, so a shape lying deep across the cut
// between them is parted from one half into the other. Both shapes escape to the heap.
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

	if sector, ok := reflex(a); ok {
		for _, half := range sector.halves() {
			if deeper, ok := Collide[T](half, b); ok && (!found || deeper.Depth > contact.Depth) {
				contact, found = deeper, true
			}
		}
	} else if sector, ok := reflex(b); ok {
		for _, half := range sector.halves() {
			if deeper, ok := Collide[T](a, half); ok && (!found || deeper.Depth > contact.Depth) {
				contact, found = deeper, true
			}
		}
	}

	if found {
		return contact, true
	}

	difference := minkowski[T]{a: a, b: b}
	normal, depth := difference.depth()

	return Contact[T]{normal, depth, difference.contact(normal)}, true
}

// Sweep returns where the first shape first touches the second while it moves by the motion, and
// false where it reaches nothing: the Point of the second shape touched, the unit Normal there
// facing the moving shape, and the Fraction of the motion made, in [0, 1]. Shapes that intersect
// before the move give false, and neither shape turns.
//
// The fraction stops within the tolerance short of the touch, which for an integer T is that of
// float64, and a motion along a side that never closes on it reaches nothing. Where two sides meet
// flat, the point is one point of the contact. A non-convex Polygon is swept edge by edge, a Sector
// past half a turn as its two halves, and an empty polygon reaches and is reached by nothing. Both
// shapes escape to the heap.
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

	if sector, ok := reflex(a); ok {
		for _, half := range sector.halves() {
			if hit, ok := Sweep[T](half, motion, b); ok && (!found || hit.Fraction < first.Fraction) {
				first, found = hit, true
			}
		}

		return first, found
	}

	if sector, ok := reflex(b); ok {
		for _, half := range sector.halves() {
			if hit, ok := Sweep[T](a, motion, half); ok && (!found || hit.Fraction < first.Fraction) {
				first, found = hit, true
			}
		}

		return first, found
	}

	return minkowski[T]{a: a, b: b}.sweep(motion.Float())
}

// Raycast returns where the ray first reaches the shape from outside, and false where it does not,
// by the Hit method of the ray naming the kind. The Fraction is in units of Vector, and an
// Origin inside the shape hits nothing. A shape of another package is reached as Sweep reaches it,
// within the tolerance, a touch in passing left out. The shape escapes to the heap.
func Raycast[T Number](ray Ray[T], shape Shape[T]) (Hit[T], bool) {
	switch kind := shape.(type) {
	case Circle[T]:
		return ray.HitCircle(kind)
	case Capsule[T]:
		return ray.HitCapsule(kind)
	case Sector[T]:
		return ray.HitSector(kind)
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
	case Ellipse[T]:
		return ray.HitEllipse(kind)
	}

	if shape.Contains(ray.Origin) {
		return Hit[T]{}, false
	}

	bounds := shape.Bounds()
	origin := minkowski[T]{b: shape, offset: ray.Origin.Float().Vector()}

	return ray.rescale(origin.sweep(ray.reach(bounds.Min, bounds.Max).Vector().Float()))
}

// intersector is a shape with an Intersects method for each of the eight kinds Intersects
// dispatches on, which every shape but Ellipse has. Intersects reads it only on a kind isKind
// names, so a shape of another package with the same methods is not taken for one.
type intersector[T Number] interface {
	IntersectsCircle(circle Circle[T]) bool
	IntersectsCapsule(capsule Capsule[T]) bool
	IntersectsSector(sector Sector[T]) bool
	IntersectsSegment(segment Segment[T]) bool
	IntersectsPolygon(polygon Polygon[T]) bool
	IntersectsRectangle(rectangle Rectangle[T]) bool
	IntersectsRegularPolygon(polygon RegularPolygon[T]) bool
	IntersectsBox(box Box[T]) bool
}

// container is a shape with a Contains method for each of the eight kinds Contains dispatches on,
// which every shape with an area but Ellipse has. Contains reads it only on a kind isKind names, so a
// shape of another package with the same methods is not taken for one.
type container[T Number] interface {
	ContainsCircle(circle Circle[T]) bool
	ContainsCapsule(capsule Capsule[T]) bool
	ContainsSector(sector Sector[T]) bool
	ContainsSegment(segment Segment[T]) bool
	ContainsPolygon(polygon Polygon[T]) bool
	ContainsRectangle(rectangle Rectangle[T]) bool
	ContainsRegularPolygon(polygon RegularPolygon[T]) bool
	ContainsBox(box Box[T]) bool
}

// isKind reports whether the shape is one of the kinds this package names, the only shapes
// Intersects and Contains read named methods of.
func isKind[T Number](shape Shape[T]) bool {
	switch shape.(type) {
	case Circle[T]:
		return true
	case Capsule[T]:
		return true
	case Sector[T]:
		return true
	case Segment[T]:
		return true
	case Polygon[T]:
		return true
	case Rectangle[T]:
		return true
	case RegularPolygon[T]:
		return true
	case Box[T]:
		return true
	case Ellipse[T]:
		return true
	default:
		return false
	}
}

// isEmpty reports whether the shape is a polygon or regular polygon of no vertices, which is
// nowhere though its Support still gives a point.
func isEmpty[T Number](shape Shape[T]) bool {
	switch kind := shape.(type) {
	case Polygon[T]:
		return kind.IsEmpty()
	case RegularPolygon[T]:
		return kind.IsEmpty()
	}

	return false
}

// eitherSegmentEllipse returns the segment and the ellipse where the shapes are one of each, the
// one pair of an Ellipse with a named Intersects method.
func eitherSegmentEllipse[T Number](a, b Shape[T]) (Segment[T], Ellipse[T], bool) {
	if segment, ok := a.(Segment[T]); ok {
		ellipse, ok := b.(Ellipse[T])

		return segment, ellipse, ok
	}

	segment, ok := b.(Segment[T])
	ellipse, isEllipse := a.(Ellipse[T])

	return segment, ellipse, ok && isEllipse
}

// eitherCircle returns whichever shape is a circle, the first where both are, with the other one.
func eitherCircle[T Number](a, b Shape[T]) (Circle[T], Shape[T], bool) {
	if circle, ok := a.(Circle[T]); ok {
		return circle, b, true
	}

	circle, ok := b.(Circle[T])

	return circle, a, ok
}

// eitherNonConvex returns whichever shape is a non-convex polygon, the first where both are, with
// the other one.
func eitherNonConvex[T Number](a, b Shape[T]) (Polygon[T], Shape[T], bool) {
	if polygon, ok := nonConvex(a); ok {
		return polygon, b, true
	}

	polygon, ok := nonConvex(b)

	return polygon, a, ok
}

// nonConvex returns the shape as a polygon that is not IsConvex, notched or too small to enclose an
// area, whose Support, that of its hull, would fill it in.
func nonConvex[T Number](shape Shape[T]) (Polygon[T], bool) {
	polygon, ok := shape.(Polygon[T])

	return polygon, ok && !polygon.IsConvex()
}

// eitherReflex returns whichever shape is a sector past half a turn, the first where both are, with
// the other one.
func eitherReflex[T Number](a, b Shape[T]) (Sector[T], Shape[T], bool) {
	if sector, ok := reflex(a); ok {
		return sector, b, true
	}

	sector, ok := reflex(b)

	return sector, a, ok
}

// reflex returns the shape as a sector past half a turn, whose Support, that of its hull, would
// fill in the notch between its radii.
func reflex[T Number](shape Shape[T]) (Sector[T], bool) {
	sector, ok := shape.(Sector[T])

	return sector, ok && sector.isReflex()
}

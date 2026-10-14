package geom

import (
	"iter"
	"math"
)

// Box is an axis-aligned box between a minimum and a maximum corner, the box Bounds returns on
// every shape. It never turns; Rectangle is the shape that can.
//
// Min is never past Max on either axis: the constructors order the corners. A box with Min past Max
// written as a struct literal or decoded from JSON gives no meaningful answers until Canonical
// repairs it. For integer T a box of width w spans w+1 lattice columns from Min to Max inclusive,
// while the image.Rectangle of its Rectangle is half-open and spans w pixels.
type Box[T Number] struct {
	Min Point[T] `json:"a"`
	Max Point[T] `json:"b"`
}

// Bx is shorthand for Box{min, max} between the opposite corners a and b, in either order, each
// axis ordered.
func Bx[T Number](a, b Point[T]) Box[T] {
	return Box[T]{Point[T]{min(a.X, b.X), min(a.Y, b.Y)}, Point[T]{max(a.X, b.X), max(a.Y, b.Y)}}
}

// BoxFromMin returns the box with its min corner at a and the size. A negative extent measures the
// other way, so a is then a corner but not the minimum one.
func BoxFromMin[T Number](a Point[T], size Size[T]) Box[T] {
	return BoxFromMinMax(a, a.Add(size.Vector()))
}

// BoxFromMinMax returns the box between the opposite corners a and b, in either order.
func BoxFromMinMax[T Number](a, b Point[T]) Box[T] {
	return Bx(a, b)
}

// BoxFromSize returns the box with its min corner at the origin and the size. A negative extent
// measures the other way, so the origin is then a corner but not the minimum one.
func BoxFromSize[T Number](size Size[T]) Box[T] {
	w, h := size.XY()

	return Box[T]{Point[T]{min(w, 0), min(h, 0)}, Point[T]{max(w, 0), max(h, 0)}}
}

// Width returns the extent from Min to Max along X.
func (b Box[T]) Width() T {
	return b.Max.X - b.Min.X
}

// Height returns the extent from Min to Max along Y.
func (b Box[T]) Height() T {
	return b.Max.Y - b.Min.Y
}

// Size returns the width and height.
func (b Box[T]) Size() Size[T] {
	return Size[T]{b.Width(), b.Height()}
}

// Vertices returns an iterator over the corners from Min, clockwise as drawn on a screen with Y
// pointing down.
func (b Box[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for _, corner := range b.corners() {
			if !yield(corner) {
				return
			}
		}
	}
}

// Edges returns an iterator over the edges in the order of Vertices, each from a corner to the next
// and the last closing back to Min.
func (b Box[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		corners := b.corners()

		edgesOf(corners[:])(yield)
	}
}

// Center returns the point halfway from Min to Max. For integer T half the size truncates toward
// Min, so the Rectangle of the box keeps its corners.
func (b Box[T]) Center() Point[T] {
	return Point[T]{b.Min.X + (b.Max.X-b.Min.X)/2, b.Min.Y + (b.Max.Y-b.Min.Y)/2}
}

// Centroid returns the Center.
func (b Box[T]) Centroid() Point[T] {
	return b.Center()
}

// Area returns the area (w·h).
func (b Box[T]) Area() float64 {
	return float64(b.Width()) * float64(b.Height())
}

// Perimeter returns the perimeter (2·(w+h)).
func (b Box[T]) Perimeter() T {
	return b.Size().Perimeter()
}

// Inertia returns the polar second moment of area about the center (w·h·(w²+h²)/12), the rotational
// inertia of the box at unit density.
func (b Box[T]) Inertia() float64 {
	w, h := b.Size().Float().XY()

	return w * h * (float64(w*w) + float64(h*h)) / 12
}

// Bounds returns the box itself.
func (b Box[T]) Bounds() Box[T] {
	return b
}

// BoundingCircle returns the circle about Center through the farthest corner, the smallest
// enclosing the box. For integer T the radius is rounded up from the truncated center, so it still
// contains all four corners.
func (b Box[T]) BoundingCircle() Circle[T] {
	center := b.Center()
	midpoint := center.Float()
	reach := max(midpoint.DistanceSquaredTo(b.Min.Float()), midpoint.DistanceSquaredTo(b.Max.Float()))

	return Circle[T]{center, castUp[T](math.Sqrt(reach))}
}

// magnitude returns the largest absolute coordinate of the corners, the magnitude the tolerance of
// a comparison with the box is taken from.
func (b Box[T]) magnitude() float64 {
	return max(b.Min.magnitude(), b.Max.magnitude())
}

// corners returns the vertices as an array in the order of Vertices: the exact corners every
// outline pair of the box reads, where its Rectangle would place them an ulp off for float T.
func (b Box[T]) corners() [4]Point[T] {
	return [4]Point[T]{b.Min, {b.Max.X, b.Min.Y}, b.Max, {b.Min.X, b.Max.Y}}
}

// Translate returns the box moved by the vector.
func (b Box[T]) Translate(vector Vector[T]) Box[T] {
	return Box[T]{b.Min.Add(vector), b.Max.Add(vector)}
}

// MoveTo returns the box with its Center at the point. For integer T an odd extent keeps the center
// nearer Min.
func (b Box[T]) MoveTo(point Point[T]) Box[T] {
	return b.Translate(point.Subtract(b.Center()))
}

// Clamp returns the box moved by the least distance that brings it within the given one, the
// receiver when it already lies exactly within. On an axis along which it is the larger, it is
// centered on the given box instead.
func (b Box[T]) Clamp(box Box[T]) Box[T] {
	return b.Translate(b.clampOffset(box))
}

// Inset returns the box with each edge moved in by its padding. An edge pushed past its opposite
// stops there, so too much Left collapses the box onto its right edge; where both paddings of an
// axis over-run, Left and Top win. A negative padding moves the edge out.
func (b Box[T]) Inset(padding Padding[T]) Box[T] {
	a := Point[T]{min(b.Min.X+padding.Left, b.Max.X), min(b.Min.Y+padding.Top, b.Max.Y)}

	return Box[T]{a, Point[T]{max(b.Max.X-padding.Right, a.X), max(b.Max.Y-padding.Bottom, a.Y)}}
}

// Outset returns the box with each edge moved out by its padding, the inverse of Inset while
// nothing was clamped.
func (b Box[T]) Outset(padding Padding[T]) Box[T] {
	return b.Inset(padding.Negate())
}

// Lerp returns the box interpolated toward the given one, both corners together, and extrapolated
// outside [0, 1]. The corners are ordered again, so an extrapolation never turns the box inside
// out.
func (b Box[T]) Lerp(box Box[T], t float64) Box[T] {
	return Box[T]{b.Min.Lerp(box.Min, t), b.Max.Lerp(box.Max, t)}.Canonical()
}

// Canonical returns the box with its corners ordered on each axis, the receiver when they already
// are.
func (b Box[T]) Canonical() Box[T] {
	return Box[T]{Point[T]{min(b.Min.X, b.Max.X), min(b.Min.Y, b.Max.Y)}, Point[T]{max(b.Min.X, b.Max.X), max(b.Min.Y, b.Max.Y)}}
}

// clampOffset returns the move Clamp makes, which Rectangle.Clamp also takes in its own frame.
func (b Box[T]) clampOffset(box Box[T]) Vector[T] {
	c1, c2 := b.Center(), box.Center()

	return Vector[T]{
		b.clampAxis(b.Min.X, b.Max.X, c1.X, box.Min.X, box.Max.X, c2.X),
		b.clampAxis(b.Min.Y, b.Max.Y, c1.Y, box.Min.Y, box.Max.Y, c2.Y),
	}
}

// clampAxis returns the move along one axis from the extent a1 to b1 with center c1 into the extent
// a2 to b2 with center c2: the difference of the centers where the first is the larger, the least
// move otherwise. It reads only its arguments, so the receiver is unnamed.
func (Box[T]) clampAxis(a1, b1, c1, a2, b2, c2 T) T {
	switch {
	case b1-a1 > b2-a2:
		return c2 - c1
	case a1 < a2:
		return a2 - a1
	case b1 > b2:
		return b2 - b1
	default:
		return 0
	}
}

// DistanceTo returns the distance from the point to the box, zero inside.
func (b Box[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(b.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (b Box[T]) DistanceSquaredTo(point Point[T]) float64 {
	return b.DistanceSquaredToBox(Box[T]{point, point})
}

// DistanceToBox returns the distance between the nearest points of the boxes, zero where they share
// a point.
func (b Box[T]) DistanceToBox(box Box[T]) float64 {
	return math.Sqrt(b.DistanceSquaredToBox(box))
}

// DistanceSquaredToBox returns the square of DistanceToBox.
func (b Box[T]) DistanceSquaredToBox(box Box[T]) float64 {
	a1, b1, a2, b2 := b.Min.Float(), b.Max.Float(), box.Min.Float(), box.Max.Float()
	dx := max(a2.X-b1.X, a1.X-b2.X, 0)
	dy := max(a2.Y-b1.Y, a1.Y-b2.Y, 0)

	distance := float64(dx*dx) + float64(dy*dy)
	if lessOrEqualSquared(distance, 0, epsilonAt[T](max(b.magnitude(), box.magnitude()))) {
		return 0
	}

	return distance
}

// Nearest returns the point of the box nearest to the given point, the point itself inside and
// otherwise the point clamped onto the box on each axis.
func (b Box[T]) Nearest(point Point[T]) Point[T] {
	if b.Contains(point) {
		return point
	}

	return Point[T]{Clamp(point.X, b.Min.X, b.Max.X), Clamp(point.Y, b.Min.Y, b.Max.Y)}
}

// Support returns the corner farthest along the direction, the first in the order of Vertices where
// two reach alike.
func (b Box[T]) Support(direction Vector[float64]) Point[float64] {
	corners := b.corners()

	return supportOf(corners[:], direction)
}

// Contains reports whether the point lies in the box.
func (b Box[T]) Contains(point Point[T]) bool {
	return b.DistanceSquaredTo(point) == 0
}

// ContainsCircle reports whether the circle lies entirely in the box.
func (b Box[T]) ContainsCircle(circle Circle[T]) bool {
	if !b.Contains(circle.Center) {
		return false
	}

	a, c, center := b.Min.Float(), b.Max.Float(), circle.Center.Float()
	gap := max(min(center.X-a.X, c.X-center.X, center.Y-a.Y, c.Y-center.Y), 0)

	return greaterOrEqualSquared(float64(gap*gap), float64(circle.Radius), circle.epsilonWith(b.magnitude()))
}

// ContainsCapsule reports whether the capsule lies entirely in the box.
func (b Box[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return b.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && b.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the sector lies entirely in the box.
func (b Box[T]) ContainsSector(sector Sector[T]) bool {
	outer, inner := b.Float(), sector.Float()
	for _, corner := range inner.corners() {
		if !outer.Contains(corner) {
			return false
		}
	}

	for _, direction := range CardinalOffsets[float64]() {
		if !outer.Contains(inner.Support(direction)) {
			return false
		}
	}

	return true
}

// ContainsSegment reports whether the segment lies entirely in the box.
func (b Box[T]) ContainsSegment(segment Segment[T]) bool {
	return b.Contains(segment.Start) && b.Contains(segment.End)
}

// ContainsPolygon reports whether the polygon lies entirely in the box.
func (b Box[T]) ContainsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !b.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRectangle reports whether the rectangle lies entirely in the box.
func (b Box[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	for vertex := range rectangle.Vertices() {
		if !b.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the box.
func (b Box[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	for vertex := range polygon.Vertices() {
		if !b.Contains(vertex) {
			return false
		}
	}

	return true
}

// ContainsBox reports whether the given box lies entirely in this one.
func (b Box[T]) ContainsBox(box Box[T]) bool {
	return b.Contains(box.Min) && b.Contains(box.Max)
}

// IntersectsCircle reports whether the box and the circle share a point.
func (b Box[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsBox(b)
}

// IntersectsCapsule reports whether the box and the capsule share a point.
func (b Box[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsBox(b)
}

// IntersectsSector reports whether the box and the sector share a point.
func (b Box[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsBox(b)
}

// IntersectsSegment reports whether the box and the segment share a point.
func (b Box[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsBox(b)
}

// IntersectsRay reports whether the box and the ray share a point.
func (b Box[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsBox(b)
}

// IntersectsLine reports whether the box and the line share a point.
func (b Box[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsBox(b)
}

// IntersectsPolygon reports whether the box and the polygon share a point.
func (b Box[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsBox(b)
}

// IntersectsRectangle reports whether the box and the rectangle share a point.
func (b Box[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.IntersectsBox(b)
}

// IntersectsRegularPolygon reports whether the box and the regular polygon share a point.
func (b Box[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return polygon.IntersectsBox(b)
}

// IntersectsBox reports whether the boxes share a point.
func (b Box[T]) IntersectsBox(box Box[T]) bool {
	return b.DistanceSquaredToBox(box) == 0
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (b Box[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	return segment.IntersectionBox(b)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (b Box[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return ray.IntersectionBox(b)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Direction.
func (b Box[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return line.IntersectionBox(b)
}

// ClipBox returns the box common to both, and false where they share no point. Boxes touching at an
// edge clip to a zero width or height, and at a corner to the zero-sized box there.
func (b Box[T]) ClipBox(box Box[T]) (Box[T], bool) {
	if !b.IntersectsBox(box) {
		return Box[T]{}, false
	}

	a := Point[T]{max(b.Min.X, box.Min.X), max(b.Min.Y, box.Min.Y)}

	return Box[T]{a, Point[T]{max(min(b.Max.X, box.Max.X), a.X), max(min(b.Max.Y, box.Max.Y), a.Y)}}, true
}

// Union returns the smallest box containing both.
func (b Box[T]) Union(box Box[T]) Box[T] {
	return Box[T]{
		Point[T]{min(b.Min.X, box.Min.X), min(b.Min.Y, box.Min.Y)},
		Point[T]{max(b.Max.X, box.Max.X), max(b.Max.Y, box.Max.Y)},
	}
}

// Equal reports whether the boxes have equal corners.
func (b Box[T]) Equal(box Box[T]) bool {
	return b.Min.Equal(box.Min) && b.Max.Equal(box.Max)
}

// IsZero reports whether the box is the zero value.
func (b Box[T]) IsZero() bool {
	return b.Equal(Box[T]{})
}

// Polygon returns the box as a polygon of its Vertices.
func (b Box[T]) Polygon() Polygon[T] {
	corners := b.corners()

	return Polygon[T]{corners[:]}
}

// Rectangle returns the box as a rectangle with no angle, the inverse of Rectangle.Bounds for a
// rectangle that is not rotated. For float T its corners, placed again from the center and the
// size, can move by an ulp.
func (b Box[T]) Rectangle() Rectangle[T] {
	return Rectangle[T]{b.Center(), b.Size(), 0}
}

// Cast converts the box to a Box of another number type, rounding each corner as Cast does.
func (b Box[T]) Cast[R Number]() Box[R] {
	return Box[R]{b.Min.Cast[R](), b.Max.Cast[R]()}
}

// Int converts the box to a Box[int], rounding each corner on its own, so the size can change by a
// unit. Round the Rectangle of the box instead where the size must be kept.
func (b Box[T]) Int() Box[int] {
	return Box[int]{b.Min.Int(), b.Max.Int()}
}

// Float converts the box to a Box[float64].
func (b Box[T]) Float() Box[float64] {
	return Box[float64]{b.Min.Float(), b.Max.Float()}
}

// String returns the box as (x,y)-(x,y), Min then Max, the form image.Rectangle prints.
func (b Box[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = b.Min.appendString(buf)
	buf = append(buf, '-')
	buf = b.Max.appendString(buf)

	return string(buf)
}

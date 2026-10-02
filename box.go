package geom

import (
	"fmt"
	"math"
)

// Box is a 2D axis-aligned box represented by its minimum and maximum corners, the box Bounds
// returns on every shape. It carries no angle, so none of its methods turns it or branches on a
// turn: it is the form for clip regions, atlas regions, viewports, layout and culling, where
// Rectangle is the shape that can be turned. Rectangle converts a box into one, and
// Rectangle.Bounds gives the box back.
//
// Min is never past Max on either axis: the constructors order the corners they are given, so a
// box with Min past Max can only be written as a struct literal or decoded from JSON, and
// Contains, Nearest and IntersectsBox give no meaningful answer for it; Canonical repairs it.
//
// The box is closed: Contains and IntersectsBox include the boundary within the tolerance, one
// comparison on the squared distance, so a box contains every point of the shape it bounds. For
// integer T a box of width w spans w+1 lattice columns from Min to Max inclusive; the
// image.Rectangle of its Rectangle is half-open as the image package requires, and therefore
// spans exactly w pixels.
type Box[T Number] struct {
	Min Point[T] `json:"a"`
	Max Point[T] `json:"b"`
}

// Bx is shorthand for Box{min, max} between two opposite corners, given in either order: each
// axis is ordered, so Min is never past Max.
func Bx[T Number](a, b Point[T]) Box[T] {
	return Box[T]{Point[T]{min(a.X, b.X), min(a.Y, b.Y)}, Point[T]{max(a.X, b.X), max(a.Y, b.Y)}}
}

// BoxFromMin creates a Box from its min corner a and size. A negative extent measures the other
// way, so the given point is a corner but no longer the minimum one.
func BoxFromMin[T Number](a Point[T], size Size[T]) Box[T] {
	return BoxFromMinMax(a, a.Add(size.Vector()))
}

// BoxFromMinMax creates a Box from two opposite corners, given in either order, as Bx does.
func BoxFromMinMax[T Number](a, b Point[T]) Box[T] {
	return Bx(a, b)
}

// BoxFromSize creates a Box from zero point and size. A negative extent measures the other
// way, so the zero point is a corner but no longer the minimum one.
func BoxFromSize[T Number](size Size[T]) Box[T] {
	w, h := size.XY()

	return Box[T]{Point[T]{min(w, 0), min(h, 0)}, Point[T]{max(w, 0), max(h, 0)}}
}

// Width returns the box width, the extent from Min to Max along X.
func (b Box[T]) Width() T {
	return b.Max.X - b.Min.X
}

// Height returns the box height, the extent from Min to Max along Y.
func (b Box[T]) Height() T {
	return b.Max.Y - b.Min.Y
}

// Size returns the width and height of the box.
func (b Box[T]) Size() Size[T] {
	return Size[T]{b.Width(), b.Height()}
}

// Center returns the point halfway from Min to Max. For integer T half the size is truncated
// toward Min, so the center of an odd span lies nearer Min and the Rectangle of the box keeps
// its corners.
func (b Box[T]) Center() Point[T] {
	w, h := b.Size().XY()

	return b.Min.AddXY(w/2, h/2)
}

// Bounds returns the box itself, the box around it.
func (b Box[T]) Bounds() Box[T] {
	return b
}

// magnitude returns the largest absolute coordinate of the box, that of a corner: the size
// epsilonAt widens the tolerance by for a comparison that reads the box.
func (b Box[T]) magnitude() float64 {
	return max(b.Min.magnitude(), b.Max.magnitude())
}

// Translate creates a new Box translated by the given vector.
func (b Box[T]) Translate(vector Vector[T]) Box[T] {
	return Box[T]{b.Min.Add(vector), b.Max.Add(vector)}
}

// Canonical creates a new Box with the corners ordered, Min the lesser on each axis, as the
// constructors build it: the box itself where it is well formed. It repairs a struct literal or
// decoded JSON with Min past Max before Contains, Nearest or IntersectsBox read it.
func (b Box[T]) Canonical() Box[T] {
	return Box[T]{Point[T]{min(b.Min.X, b.Max.X), min(b.Min.Y, b.Max.Y)}, Point[T]{max(b.Min.X, b.Max.X), max(b.Min.Y, b.Max.Y)}}
}

// Inset creates a new Box inset by the given padding amounts, each edge moved in by its own. An
// edge pushed past its opposite stops there, so a padding larger than the box collapses it to a
// zero extent on that axis, at the opposite edge rather than at the one that over-ran: too much
// Left collapses it onto the right edge, too much Right onto the left one. Where both paddings on
// an axis over-run, Left and Top win, since Min moves first and Max is then stopped at it.
// A negative padding outsets the box, so Outset undoes Inset as long as nothing was clamped.
func (b Box[T]) Inset(padding Padding[T]) Box[T] {
	a := Point[T]{min(b.Min.X+padding.Left, b.Max.X), min(b.Min.Y+padding.Top, b.Max.Y)}

	return Box[T]{a, Point[T]{max(b.Max.X-padding.Right, a.X), max(b.Max.Y-padding.Bottom, a.Y)}}
}

// Outset creates a new Box expanded by the given padding amounts, the inverse of Inset.
func (b Box[T]) Outset(padding Padding[T]) Box[T] {
	return b.Inset(padding.Negate())
}

// Clamp creates a new Box moved so that it lies within the given one, keeping its size: the box
// itself where it already lies within, and otherwise moved by the shortest distance that brings
// it in. On an axis along which it is the larger it is centered on the other instead.
func (b Box[T]) Clamp(box Box[T]) Box[T] {
	return b.Translate(b.clampOffset(box))
}

// Contains reports whether the given point lies within the box, boundary included within
// the tolerance: DistanceSquaredTo is zero there.
func (b Box[T]) Contains(point Point[T]) bool {
	return b.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the box: zero
// exactly where Contains holds, so a point within the tolerance of the boundary is at distance
// zero rather than at the rounding error that put it there.
func (b Box[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(b.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons: zero for a point inside, or within the tolerance of the boundary, and otherwise
// the sum of the squared gaps beyond the box on the two axes, with no edge to walk. It is a
// float64 even for an integer T, so the box answers as every Shape does.
func (b Box[T]) DistanceSquaredTo(point Point[T]) float64 {
	return b.distanceSquaredToBox(Box[T]{point, point})
}

// Nearest returns the point of the box nearest to the given point: the point itself exactly
// where Contains holds, and otherwise the point clamped onto the box on each axis, a lattice
// point for an integer T.
func (b Box[T]) Nearest(point Point[T]) Point[T] {
	if b.Contains(point) {
		return point
	}

	return Point[T]{Clamp(point.X, b.Min.X, b.Max.X), Clamp(point.Y, b.Min.Y, b.Max.Y)}
}

// EnclosesCircle reports whether the circle lies within the box, as Rectangle.EnclosesCircle
// decides on the box's Rectangle.
func (b Box[T]) EnclosesCircle(circle Circle[T]) bool {
	return b.Rectangle().EnclosesCircle(circle)
}

// EnclosesSegment reports whether the segment lies within the box, as
// Rectangle.EnclosesSegment decides on the box's Rectangle.
func (b Box[T]) EnclosesSegment(segment Segment[T]) bool {
	return b.Rectangle().EnclosesSegment(segment)
}

// EnclosesPolygon reports whether the polygon lies within the box, as
// Rectangle.EnclosesPolygon decides on the box's Rectangle.
func (b Box[T]) EnclosesPolygon(polygon Polygon[T]) bool {
	return b.Rectangle().EnclosesPolygon(polygon)
}

// EnclosesRectangle reports whether the rectangle lies within the box, as
// Rectangle.EnclosesRectangle decides on the box's Rectangle.
func (b Box[T]) EnclosesRectangle(rectangle Rectangle[T]) bool {
	return b.Rectangle().EnclosesRectangle(rectangle)
}

// EnclosesRegularPolygon reports whether the regular polygon lies within the box, as
// Rectangle.EnclosesRegularPolygon decides on the box's Rectangle.
func (b Box[T]) EnclosesRegularPolygon(polygon RegularPolygon[T]) bool {
	return b.Rectangle().EnclosesRegularPolygon(polygon)
}

// EnclosesBox reports whether the given box lies within this one, as
// Rectangle.EnclosesRectangle decides on the Rectangle of each.
func (b Box[T]) EnclosesBox(box Box[T]) bool {
	return b.Rectangle().EnclosesRectangle(box.Rectangle())
}

// IntersectsCircle reports whether the box and the circle overlap, as Circle.IntersectsBox does.
func (b Box[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsBox(b)
}

// IntersectsSegment reports whether the box and the segment share a point, as
// Segment.IntersectsBox does.
func (b Box[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsBox(b)
}

// IntersectionSegment returns the points where the segment crosses the box boundary, as
// Segment.IntersectionBox does.
func (b Box[T]) IntersectionSegment(segment Segment[T]) []Point[T] {
	return segment.IntersectionBox(b)
}

// AppendIntersectionSegment appends the points IntersectionSegment returns to dst and returns the
// extended slice, as Segment.AppendIntersectionBox does.
func (b Box[T]) AppendIntersectionSegment(dst []Point[T], segment Segment[T]) []Point[T] {
	return segment.AppendIntersectionBox(dst, b)
}

// IntersectsRay reports whether the box and the ray share a point, as
// Ray.IntersectsBox does.
func (b Box[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsBox(b)
}

// IntersectionRay returns the points where the ray crosses the box boundary, as
// Ray.IntersectionBox does.
func (b Box[T]) IntersectionRay(ray Ray[T]) []Point[T] {
	return ray.IntersectionBox(b)
}

// AppendIntersectionRay appends the points IntersectionRay returns to dst and returns the
// extended slice, as Ray.AppendIntersectionBox does.
func (b Box[T]) AppendIntersectionRay(dst []Point[T], ray Ray[T]) []Point[T] {
	return ray.AppendIntersectionBox(dst, b)
}

// IntersectsPolygon reports whether the box and the polygon share a point, as
// Polygon.IntersectsBox does.
func (b Box[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsBox(b)
}

// IntersectsRectangle reports whether the box and the rectangle share a point, as
// Rectangle.IntersectsBox does.
func (b Box[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	return rectangle.IntersectsBox(b)
}

// IntersectsRegularPolygon reports whether the box and the regular polygon share a point, as
// RegularPolygon.IntersectsBox does.
func (b Box[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	return polygon.IntersectsBox(b)
}

// IntersectsBox reports whether the boxes share a point: the gap between them on the two axes
// is zero within the tolerance, the one comparison DistanceSquaredTo makes, so boxes that touch
// at an edge or a corner intersect, the same closed convention as Contains.
func (b Box[T]) IntersectsBox(box Box[T]) bool {
	return b.distanceSquaredToBox(box) == 0
}

// IntersectionBox returns the box common to both, and false when they do not intersect.
// Touching boxes intersect in a box of zero width or height, and boxes meeting at a corner in
// the zero-sized box at that corner, exactly where IntersectsBox holds, which is asked first:
// a corner admitted by the tolerance is placed on the boundary of the other box, never beyond
// it.
func (b Box[T]) IntersectionBox(box Box[T]) (Box[T], bool) {
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

// clampOffset returns the move Clamp makes, axis by axis, the one Rectangle.Clamp takes for its
// boxes.
func (b Box[T]) clampOffset(box Box[T]) Vector[T] {
	c1, c2 := b.Center(), box.Center()

	return Vector[T]{
		b.clampAxis(b.Min.X, b.Max.X, c1.X, box.Min.X, box.Max.X, c2.X),
		b.clampAxis(b.Min.Y, b.Max.Y, c1.Y, box.Min.Y, box.Max.Y, c2.Y),
	}
}

// clampAxis returns the move along one axis that Clamp makes, from the extent a1 to b1 and the
// center c1 of the box being moved and those of the one it is clamped within: the difference of
// the centers where the first is the larger, the least move that brings it within otherwise,
// and zero where it already lies within. It reads no field of the box, only the extents it is
// given, so the receiver is unnamed.
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

// distanceSquaredToBox returns the squared distance between the nearest points of the two
// boxes, from the gap between them on each axis: zero where they overlap or the gap is within
// the tolerance, the one comparison Contains and IntersectsBox both read. A NaN coordinate leaves
// a NaN gap, which the comparison never admits. The gaps are taken in float64, so a narrow
// integer T cannot overflow them.
func (b Box[T]) distanceSquaredToBox(box Box[T]) float64 {
	a1, b1, a2, b2 := b.Min.Float(), b.Max.Float(), box.Min.Float(), box.Max.Float()
	dx := max(a2.X-b1.X, a1.X-b2.X, 0)
	dy := max(a2.Y-b1.Y, a1.Y-b2.Y, 0)

	distance := float64(dx*dx) + float64(dy*dy)
	if lessOrEqualSquared(distance, 0, epsilonAt[T](max(b.magnitude(), box.magnitude()))) {
		return 0
	}

	return distance
}

// Equal checks for equal corners using tolerant numeric comparison.
func (b Box[T]) Equal(box Box[T]) bool {
	return b.Min.Equal(box.Min) && b.Max.Equal(box.Max)
}

// IsZero checks if both corners are zero.
func (b Box[T]) IsZero() bool {
	return b.Equal(Box[T]{})
}

// Rectangle converts the box into the Rectangle with the same corners, not rotated: the shape
// the outline pairs of the box are decided on, and the inverse of Rectangle.Bounds for a
// rectangle that is not rotated.
func (b Box[T]) Rectangle() Rectangle[T] {
	return Rectangle[T]{b.Center(), b.Size(), 0}
}

// Cast converts the box to a Box of another number type, rounding each corner as Cast does.
func (b Box[T]) Cast[R Number]() Box[R] {
	return Box[R]{b.Min.Cast[R](), b.Max.Cast[R]()}
}

// Int converts the box to a Box[int], rounding each corner on its own, so the size of a float
// box can change by a unit as it moves through positions no lattice expresses. Round the
// Rectangle of the box instead where the size must be kept.
func (b Box[T]) Int() Box[int] {
	return Box[int]{b.Min.Int(), b.Max.Int()}
}

// Float converts the box to a Box[float64].
func (b Box[T]) Float() Box[float64] {
	return Box[float64]{b.Min.Float(), b.Max.Float()}
}

// String returns the box by its corners: (x,y)-(x,y), Min then Max, the form image.Rectangle
// prints.
func (b Box[T]) String() string {
	return fmt.Sprintf("%s-%s", b.Min.String(), b.Max.String())
}

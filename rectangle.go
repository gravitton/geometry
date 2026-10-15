package geom

import (
	"cmp"
	"iter"
	"math"
)

// Rectangle is a rectangle of a center and a size, turned by Angle radians about the center.
//
// The size is never negative: constructors and changes take it absolute. A negative extent written
// as a struct literal or decoded from JSON gives no meaningful answers until Canonical repairs it.
//
// Width, Height and the named corners are those of the rectangle's own frame before the turn, so
// TopLeft is the corner that was top-left wherever the turn has put it. Min, Max and Bounds are the
// axis-aligned extent of the vertices, which coincide with the corners only while the rectangle is
// not rotated. For integer T a rectangle of width w spans w+1 lattice columns from Min to Max
// inclusive, while the image.Rectangle that Rectangle returns is half-open and spans w pixels. A
// rotated rectangle has each corner rounded onto the lattice, and its edges join those corners.
type Rectangle[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	Angle  float64  `json:"a,omitzero"`
}

// Rect is shorthand for Rectangle{center, size, angle}, with the size taken absolute.
func Rect[T Number](center Point[T], size Size[T], angle float64) Rectangle[T] {
	return Rectangle[T]{center, size.Abs(), angle}
}

// RectangleFromCenter returns the aligned rectangle with its center at the point and the size, taken
// absolute.
func RectangleFromCenter[T Number](center Point[T], size Size[T]) Rectangle[T] {
	return Rectangle[T]{center, size.Abs(), 0}
}

// RectangleFromMin returns the rectangle with its min corner at a and the size. A negative extent
// measures the other way, so a is then a corner but not the minimum one.
func RectangleFromMin[T Number](a Point[T], size Size[T]) Rectangle[T] {
	return RectangleFromMinMax(a, a.Add(size.Vector()))
}

// RectangleFromMax returns the rectangle with its max corner at b and the size. A negative extent
// measures the other way, so b is then a corner but not the maximum one.
func RectangleFromMax[T Number](b Point[T], size Size[T]) Rectangle[T] {
	return RectangleFromMinMax(b.Add(size.Vector().Negate()), b)
}

// RectangleFromMinMax returns the rectangle between the opposite corners a and b, in either order.
// For integer T the center truncates toward Min, so Max - Min is exactly the size.
func RectangleFromMinMax[T Number](a, b Point[T]) Rectangle[T] {
	a, b = Point[T]{min(a.X, b.X), min(a.Y, b.Y)}, Point[T]{max(a.X, b.X), max(a.Y, b.Y)}
	w, h := b.Subtract(a).XY()

	return Rectangle[T]{a.AddXY(w/2, h/2), Size[T]{w, h}, 0}
}

// RectangleFromSize returns the rectangle with its min corner at the origin and the size. A
// negative extent measures the other way, so the origin is then a corner but not the minimum one.
func RectangleFromSize[T Number](size Size[T]) Rectangle[T] {
	return RectangleFromMin(Pt[T](0, 0), size)
}

// Width returns the extent along the rectangle's own X axis.
func (r Rectangle[T]) Width() T {
	return r.Size.Width
}

// Height returns the extent along the rectangle's own Y axis.
func (r Rectangle[T]) Height() T {
	return r.Size.Height
}

// Min returns the minimum corner of the axis-aligned extent: TopLeft while the rectangle is not
// rotated, the minimum of its vertices once it is.
func (r Rectangle[T]) Min() Point[T] {
	a, _ := r.MinMax()

	return a
}

// Max returns the maximum corner of the axis-aligned extent: BottomRight while the rectangle is not
// rotated, the maximum of its vertices once it is.
func (r Rectangle[T]) Max() Point[T] {
	_, b := r.MinMax()

	return b
}

// MinMax returns Min and Max together, the corners of Bounds.
func (r Rectangle[T]) MinMax() (Point[T], Point[T]) {
	if r.IsAligned() {
		a, b := r.localMinMax()

		return r.Center.Add(a), r.Center.Add(b)
	}

	corners := r.corners()

	return minMaxOf(corners[:])
}

// TopLeft returns the top-left corner in the frame before the turn.
func (r Rectangle[T]) TopLeft() Point[T] {
	a, _ := r.localMinMax()

	return r.worldPoint(a)
}

// TopRight returns the top-right corner in the frame before the turn.
func (r Rectangle[T]) TopRight() Point[T] {
	a, b := r.localMinMax()

	return r.worldPoint(Vector[T]{b.X, a.Y})
}

// BottomRight returns the bottom-right corner in the frame before the turn.
func (r Rectangle[T]) BottomRight() Point[T] {
	_, b := r.localMinMax()

	return r.worldPoint(b)
}

// BottomLeft returns the bottom-left corner in the frame before the turn.
func (r Rectangle[T]) BottomLeft() Point[T] {
	a, b := r.localMinMax()

	return r.worldPoint(Vector[T]{a.X, b.Y})
}

// Anchor returns the point of the boundary in the direction from the center, in the frame before
// the turn: a corner for a diagonal, the midpoint of an edge for a cardinal, and the center for
// DirectionNone. For integer T an edge midpoint of a rotated rectangle is rounded once and can lie
// just off the edge between the rounded corners.
func (r Rectangle[T]) Anchor(direction Direction) Point[T] {
	a, b := r.localMinMax()

	switch direction.normalize() {
	case TopLeft:
		return r.worldPoint(a)
	case Top:
		return r.worldPoint(Vector[T]{0, a.Y})
	case TopRight:
		return r.worldPoint(Vector[T]{b.X, a.Y})
	case Right:
		return r.worldPoint(Vector[T]{b.X, 0})
	case BottomRight:
		return r.worldPoint(b)
	case Bottom:
		return r.worldPoint(Vector[T]{0, b.Y})
	case BottomLeft:
		return r.worldPoint(Vector[T]{a.X, b.Y})
	case Left:
		return r.worldPoint(Vector[T]{a.X, 0})
	default:
		return r.Center
	}
}

// Vertices returns an iterator over the corners from TopLeft, clockwise as drawn on a screen with Y
// pointing down.
func (r Rectangle[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for _, corner := range r.corners() {
			if !yield(corner) {
				return
			}
		}
	}
}

// Edges returns an iterator over the edges in the order of Vertices, each from a corner to the next
// and the last closing back to TopLeft.
func (r Rectangle[T]) Edges() iter.Seq[Segment[T]] {
	return func(yield func(Segment[T]) bool) {
		corners := r.corners()

		edgesOf(corners[:])(yield)
	}
}

// Centroid returns the center.
func (r Rectangle[T]) Centroid() Point[T] {
	return r.Center
}

// Area returns the area (w·h).
func (r Rectangle[T]) Area() float64 {
	return float64(r.Size.Width) * float64(r.Size.Height)
}

// Perimeter returns the perimeter (2·(w+h)).
func (r Rectangle[T]) Perimeter() T {
	return r.Size.Perimeter()
}

// Inertia returns the polar second moment of area about the center (w·h·(w²+h²)/12), the rotational
// inertia of the rectangle at unit density.
func (r Rectangle[T]) Inertia() float64 {
	w, h := r.Size.Float().XY()

	return w * h * (float64(w*w) + float64(h*h)) / 12
}

// AspectRatio returns the width divided by the height, and zero for a zero height.
func (r Rectangle[T]) AspectRatio() float64 {
	return r.Size.AspectRatio()
}

// Bounds returns the axis-aligned bounding box.
func (r Rectangle[T]) Bounds() Box[T] {
	a, b := r.MinMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the circle about the center through the farthest corner, the smallest
// enclosing the rectangle. For integer T the radius is rounded up, so it still contains the rounded
// corners.
func (r Rectangle[T]) BoundingCircle() Circle[T] {
	center := r.Center.Float()

	var reach float64
	for _, corner := range r.corners() {
		reach = max(reach, center.DistanceSquaredTo(corner.Float()))
	}

	return Circle[T]{r.Center, castUp[T](math.Sqrt(reach))}
}

// localMinMax returns the offsets of the min and max corners from the center in the frame before
// the turn. For integer T the min offset truncates toward the center and the max takes the rest, so
// they differ by exactly the size.
func (r Rectangle[T]) localMinMax() (Vector[T], Vector[T]) {
	w, h := r.Size.XY()

	return Vector[T]{-w / 2, -h / 2}, Vector[T]{w - w/2, h - h/2}
}

// worldPoint returns the point at the offset from the center in the frame before the turn. An
// aligned rectangle adds the offset alone, so integer corners stay exact; otherwise the offset is
// rotated and rounded once.
func (r Rectangle[T]) worldPoint(offset Vector[T]) Point[T] {
	if r.IsAligned() {
		return r.Center.Add(offset)
	}

	return r.Center.Add(offset.Rotate(r.Angle))
}

// localOffset returns the offset of the point from the center in the frame before the turn, in
// float64, the inverse of worldPoint.
func (r Rectangle[T]) localOffset(point Point[T]) Vector[float64] {
	return point.Float().Subtract(r.Center.Float()).Rotate(-r.Angle)
}

// localRectangle returns the given rectangle in the frame of this one, centered at its offset and
// not turned. A rectangle half a turn from this one lies mirrored about its center there, so its
// center moves by the unit an odd integer extent leaves between its corner offsets. For integer T
// the offset is rounded.
func (r Rectangle[T]) localRectangle(rectangle Rectangle[T]) Rectangle[T] {
	offset := r.localOffset(rectangle.Center).Cast[T]()
	if AngleDistance(r.Angle, rectangle.Angle) > Pi/2 {
		a, b := rectangle.localMinMax()
		offset = offset.Subtract(a.Add(b))
	}

	return Rectangle[T]{offset.Point(), rectangle.Size, 0}
}

// sharedFrame returns the rectangle whose frame a parallel pair is measured in: the one with the
// lesser center by Point.Compare, then the lesser angle. Both sides pick the same one, so the pair
// answers alike in either order.
func (r Rectangle[T]) sharedFrame(rectangle Rectangle[T]) Rectangle[T] {
	if cmp.Or(rectangle.Center.Compare(r.Center), cmp.Compare(rectangle.Angle, r.Angle)) < 0 {
		return rectangle
	}

	return r
}

// worldRectangle returns a rectangle found by localRectangle in this frame turned back into the
// world, with this angle.
func (r Rectangle[T]) worldRectangle(rectangle Rectangle[T]) Rectangle[T] {
	return Rectangle[T]{r.worldPoint(rectangle.Center.Vector()), rectangle.Size, r.Angle}
}

// corners returns the vertices as an array, placed once for Vertices, Edges and Polygon.
func (r Rectangle[T]) corners() [4]Point[T] {
	a, b := r.localMinMax()

	return [4]Point[T]{r.worldPoint(a), r.worldPoint(Vector[T]{b.X, a.Y}), r.worldPoint(b), r.worldPoint(Vector[T]{a.X, b.Y})}
}

// Translate returns the rectangle moved by the vector.
func (r Rectangle[T]) Translate(vector Vector[T]) Rectangle[T] {
	return Rectangle[T]{r.Center.Add(vector), r.Size, r.Angle}
}

// MoveTo returns the rectangle with its center at the point.
func (r Rectangle[T]) MoveTo(point Point[T]) Rectangle[T] {
	return Rectangle[T]{point, r.Size, r.Angle}
}

// AlignTo returns the rectangle moved so that its Anchor in the direction lies at the point.
func (r Rectangle[T]) AlignTo(direction Direction, point Point[T]) Rectangle[T] {
	return r.Translate(point.Subtract(r.Anchor(direction)))
}

// Clamp returns the rectangle moved by the least distance that brings it within the given one, the
// receiver when it already lies within. On an axis of the given rectangle along which it is the
// larger, it is centered on it instead. For integer T the move of a rotated pair is rounded once,
// so the corners can land a unit past the edges.
func (r Rectangle[T]) Clamp(rectangle Rectangle[T]) Rectangle[T] {
	if r.IsAligned() && rectangle.IsAligned() {
		return r.Translate(r.Bounds().clampOffset(rectangle.Bounds()))
	}

	var local [4]Point[float64]
	for i, corner := range r.corners() {
		local[i] = rectangle.localOffset(corner).Point()
	}

	a, b := rectangle.localMinMax()
	extent := BoxFromMinMax(minMaxOf(local[:]))
	frame := Box[float64]{a.Float().Point(), b.Float().Point()}
	move := extent.clampOffset(frame)

	return r.Translate(move.Rotate(rectangle.Angle).Cast[T]())
}

// Rotate returns the rectangle turned by the angle in radians about its center. The stored angle is
// normalized to [0, 2π), so turning back by its own angle leaves it aligned with exact corners.
func (r Rectangle[T]) Rotate(angle float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size, NormalizeAngle(r.Angle + angle)}
}

// Scale returns the rectangle with its size scaled by the factor, taken absolute.
func (r Rectangle[T]) Scale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Scale(factor).Abs(), r.Angle}
}

// ScaleXY returns the rectangle with its size scaled by the factors along its own axes, taken
// absolute.
func (r Rectangle[T]) ScaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ScaleXY(factorX, factorY).Abs(), r.Angle}
}

// Unscale returns the rectangle with its size divided by the factor, taken absolute. It panics for
// a zero factor.
func (r Rectangle[T]) Unscale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Unscale(factor).Abs(), r.Angle}
}

// UnscaleXY returns the rectangle with its size divided by the factors along its own axes, taken
// absolute. It panics for a zero factor.
func (r Rectangle[T]) UnscaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.UnscaleXY(factorX, factorY).Abs(), r.Angle}
}

// Resize returns the rectangle with the size, taken absolute.
func (r Rectangle[T]) Resize(size Size[T]) Rectangle[T] {
	return Rectangle[T]{r.Center, size.Abs(), r.Angle}
}

// Transform returns the rectangle mapped by the matrix: the center moves, the size scales along the
// axes of the matrix, and the angle turns with it or mirrors for a reflection. A shear, or a scale
// of the axes of a turned rectangle, gives a parallelogram, and the result is the nearest rectangle
// by Matrix.Scaling and Matrix.Angle; Polygon().Transform gives the exact one. A matrix that
// collapses an axis gives a zero extent on it.
func (r Rectangle[T]) Transform[M Float](matrix Matrix[M]) Rectangle[T] {
	m := matrix.Float()
	scaling := m.Scaling()

	return Rectangle[T]{r.Center.Transform(matrix), r.Size.ScaleXY(scaling.X, scaling.Y).Abs(), m.turnedAngle(r.Angle, scaling)}
}

// Grow returns the rectangle with each extent increased by the amount, clamped at zero, so each
// side moves out by half of it. For integer T an odd amount moves one side a unit farther than the
// other; Outset moves a chosen side by a whole amount.
func (r Rectangle[T]) Grow(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Grow(amount).AtLeastZero(), r.Angle}
}

// GrowXY returns the rectangle with its extents increased by the amounts along its own axes,
// clamped at zero. For integer T an odd amount moves one side a unit farther than the other.
func (r Rectangle[T]) GrowXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.GrowXY(amountX, amountY).AtLeastZero(), r.Angle}
}

// Shrink returns the rectangle with each extent decreased by the amount, clamped at zero, so each
// side moves in by half of it. For integer T an odd amount moves one side a unit farther than the
// other; Inset moves a chosen side by a whole amount.
func (r Rectangle[T]) Shrink(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Shrink(amount).AtLeastZero(), r.Angle}
}

// ShrinkXY returns the rectangle with its extents decreased by the amounts along its own axes,
// clamped at zero. For integer T an odd amount moves one side a unit farther than the other.
func (r Rectangle[T]) ShrinkXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ShrinkXY(amountX, amountY).AtLeastZero(), r.Angle}
}

// Inset returns the rectangle with each edge moved in by its padding, in the frame before the turn.
// An edge pushed past its opposite stops there, so too much Left collapses the rectangle onto its
// right edge; where both paddings of an axis over-run, Left and Top win. A negative padding moves
// the edge out.
func (r Rectangle[T]) Inset(padding Padding[T]) Rectangle[T] {
	a, b := r.localMinMax()
	inset := Box[T]{a.Point(), b.Point()}.Inset(padding)

	return r.worldRectangle(inset.Rectangle())
}

// Outset returns the rectangle with each edge moved out by its padding, the inverse of Inset while
// nothing was clamped.
func (r Rectangle[T]) Outset(padding Padding[T]) Rectangle[T] {
	return r.Inset(padding.Negate())
}

// Lerp returns the rectangle interpolated toward the given one, center and size together and the
// angle along the shorter arc, and extrapolated outside [0, 1]. The size is taken absolute and the
// angle normalized to [0, 2π).
func (r Rectangle[T]) Lerp(rectangle Rectangle[T], t float64) Rectangle[T] {
	return Rectangle[T]{r.Center.Lerp(rectangle.Center, t), r.Size.Lerp(rectangle.Size, t).Abs(), NormalizeAngle(LerpAngle(r.Angle, rectangle.Angle, t))}
}

// Canonical returns the rectangle with its size taken absolute and its angle normalized to [0, 2π),
// the receiver when it already is. An angle within Delta of zero or a full turn becomes exactly
// zero, so the rectangle is aligned again; Rotate never snaps.
func (r Rectangle[T]) Canonical() Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Abs(), snapAngle(r.Angle)}
}

// DistanceTo returns the distance from the point to the rectangle, zero inside.
func (r Rectangle[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(r.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (r Rectangle[T]) DistanceSquaredTo(point Point[T]) float64 {
	return r.walk(point).result()
}

// Nearest returns the point of the rectangle nearest to the given point, the point itself inside.
func (r Rectangle[T]) Nearest(point Point[T]) Point[T] {
	return r.walk(point).nearest(point)
}

// Support returns the corner farthest along the direction, the first in the order of Vertices where
// two reach alike.
func (r Rectangle[T]) Support(direction Vector[float64]) Point[float64] {
	corners := r.corners()

	return supportOf(corners[:], direction)
}

// Contains reports whether the point lies in the rectangle.
func (r Rectangle[T]) Contains(point Point[T]) bool {
	return r.DistanceSquaredTo(point) == 0
}

// ContainsCircle reports whether the circle lies entirely in the rectangle.
func (r Rectangle[T]) ContainsCircle(circle Circle[T]) bool {
	if r.IsAligned() {
		return r.Bounds().ContainsCircle(circle)
	}

	return r.walk(circle.Center).clears(circle)
}

// ContainsCapsule reports whether the capsule lies entirely in the rectangle.
func (r Rectangle[T]) ContainsCapsule(capsule Capsule[T]) bool {
	return r.ContainsCircle(Circle[T]{capsule.Segment.Start, capsule.Radius}) && r.ContainsCircle(Circle[T]{capsule.Segment.End, capsule.Radius})
}

// ContainsSector reports whether the sector lies entirely in the rectangle.
func (r Rectangle[T]) ContainsSector(sector Sector[T]) bool {
	if !r.Bounds().ContainsBox(sector.Bounds()) {
		return false
	}

	for edge := range r.Edges() {
		if !sector.liesLeftOf(edge.Float()) {
			return false
		}
	}

	return true
}

// ContainsSegment reports whether the segment lies entirely in the rectangle.
func (r Rectangle[T]) ContainsSegment(segment Segment[T]) bool {
	a, b := r.MinMax()

	return r.containsWithin(segment.Start, a, b) && r.containsWithin(segment.End, a, b)
}

// ContainsPolygon reports whether the polygon lies entirely in the rectangle.
func (r Rectangle[T]) ContainsPolygon(polygon Polygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := r.MinMax()
	for vertex := range polygon.Vertices() {
		if !r.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// ContainsRectangle reports whether the given rectangle lies entirely in this one.
func (r Rectangle[T]) ContainsRectangle(rectangle Rectangle[T]) bool {
	a, b := r.MinMax()
	for vertex := range rectangle.Vertices() {
		if !r.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// ContainsRegularPolygon reports whether the regular polygon lies entirely in the rectangle.
func (r Rectangle[T]) ContainsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a, b := r.MinMax()
	for vertex := range polygon.Vertices() {
		if !r.containsWithin(vertex, a, b) {
			return false
		}
	}

	return true
}

// ContainsBox reports whether the box lies entirely in the rectangle.
func (r Rectangle[T]) ContainsBox(box Box[T]) bool {
	a, b := r.MinMax()
	for _, corner := range box.corners() {
		if !r.containsWithin(corner, a, b) {
			return false
		}
	}

	return true
}

// IntersectsCircle reports whether the rectangle and the circle share a point.
func (r Rectangle[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRectangle(r)
}

// IntersectsCapsule reports whether the rectangle and the capsule share a point.
func (r Rectangle[T]) IntersectsCapsule(capsule Capsule[T]) bool {
	return capsule.IntersectsRectangle(r)
}

// IntersectsSector reports whether the rectangle and the sector share a point.
func (r Rectangle[T]) IntersectsSector(sector Sector[T]) bool {
	return sector.IntersectsRectangle(r)
}

// IntersectsSegment reports whether the rectangle and the segment share a point.
func (r Rectangle[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsRectangle(r)
}

// IntersectsRay reports whether the rectangle and the ray share a point.
func (r Rectangle[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsRectangle(r)
}

// IntersectsLine reports whether the rectangle and the line share a point.
func (r Rectangle[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsRectangle(r)
}

// IntersectsPolygon reports whether the rectangle and the polygon share a point.
func (r Rectangle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsRectangle(r)
}

// IntersectsRectangle reports whether the rectangles share a point. Two rectangles of the same
// angle, or half a turn apart, are measured in a frame they share, as ClipRectangle and Union are,
// so they can differ from IntersectsPolygon on their Polygon forms: by a unit for integer T, within
// the tolerance for float T. Far from the origin that frame can miss a touch ContainsRectangle
// admits.
func (r Rectangle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	if r.IsAligned() && rectangle.IsAligned() {
		return r.Bounds().IntersectsBox(rectangle.Bounds())
	}

	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	switch {
	case !overlaps(a1, b1, a2, b2):
		return false
	case r.parallel(rectangle):
		frame := r.sharedFrame(rectangle)

		return frame.localRectangle(r).IntersectsRectangle(frame.localRectangle(rectangle))
	default:
		return r.meetsWithin(rectangle, a1, b1, a2, b2)
	}
}

// IntersectsRegularPolygon reports whether the rectangle and the regular polygon share a point.
func (r Rectangle[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.IsEmpty() {
		return false
	}

	a1, b1 := r.MinMax()
	a2, b2 := polygon.minMax()

	if !overlaps(a1, b1, a2, b2) {
		return false
	}

	if polygon.containsWithin(r.TopLeft(), a2, b2) || r.containsWithin(polygon.vertex(0), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range r.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range polygon.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// IntersectsBox reports whether the rectangle and the box share a point.
func (r Rectangle[T]) IntersectsBox(box Box[T]) bool {
	if r.IsAligned() {
		return r.Bounds().IntersectsBox(box)
	}

	a, b := r.MinMax()
	if !overlaps(a, b, box.Min, box.Max) {
		return false
	}

	corners, others := r.corners(), box.corners()
	if box.Contains(corners[0]) || r.containsWithin(others[0], a, b) {
		return true
	}

	return r.edgesMeet(corners, others, box.Min, box.Max)
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (r Rectangle[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	return segment.IntersectionRectangle(r)
}

// IntersectionRay returns the points where the ray crosses the boundary, from its Origin on.
func (r Rectangle[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return ray.IntersectionRectangle(r)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (r Rectangle[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return line.IntersectionRectangle(r)
}

// ClipRectangle returns the rectangle common to both, and false where they share no point. Touching
// rectangles clip to a zero width or height. Only rectangles of the same angle, or half a turn
// apart, overlap in a rectangle; others return false even where IntersectsRectangle holds. For
// integer T the result of a rotated pair is rounded through the frame they share.
func (r Rectangle[T]) ClipRectangle(rectangle Rectangle[T]) (Rectangle[T], bool) {
	if !r.IntersectsRectangle(rectangle) {
		return Rectangle[T]{}, false
	}

	switch {
	case r.IsAligned() && rectangle.IsAligned():
		overlap, _ := r.Bounds().ClipBox(rectangle.Bounds())

		return overlap.Rectangle(), true
	case r.parallel(rectangle):
		frame := r.sharedFrame(rectangle)
		overlap, _ := frame.localRectangle(r).ClipRectangle(frame.localRectangle(rectangle))

		return frame.worldRectangle(overlap), true
	default:
		return Rectangle[T]{}, false
	}
}

// Union returns the smallest rectangle containing both, of their angle for two of the same angle or
// half a turn apart, and otherwise the aligned rectangle of the union of their Bounds. For integer
// T the union of a rotated pair is rounded through the frame they share and need not enclose their
// rounded corners.
func (r Rectangle[T]) Union(rectangle Rectangle[T]) Rectangle[T] {
	if !(r.IsAligned() && rectangle.IsAligned()) && r.parallel(rectangle) {
		frame := r.sharedFrame(rectangle)

		return frame.worldRectangle(frame.localRectangle(r).Union(frame.localRectangle(rectangle)))
	}

	return r.Bounds().Union(rectangle.Bounds()).Rectangle()
}

// walk folds every edge into the edgeWalk DistanceSquaredTo, Nearest and ContainsCircle read,
// stopping at an edge the point lies on. An aligned rectangle answers a point within its extent at
// once, with no edge measured, which is why ContainsCircle asks its Bounds instead.
func (r Rectangle[T]) walk(point Point[T]) edgeWalk[T] {
	if r.IsAligned() {
		if a, b := r.MinMax(); a.X <= point.X && point.X <= b.X && a.Y <= point.Y && point.Y <= b.Y {
			return edgeWalk[T]{inside: true}
		}
	}

	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range r.Edges() {
		if w.step(edge, point) {
			break
		}
	}

	return w
}

// containsWithin is Contains for a caller that already holds the corners of MinMax.
func (r Rectangle[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && r.DistanceSquaredTo(point) == 0
}

// meetsWithin reports whether two rectangles whose extents overlap share a point: a corner of one
// lies in the other, or their edges meet. Each side places its corners once here rather than for
// every edge of the other.
func (r Rectangle[T]) meetsWithin(rectangle Rectangle[T], a1, b1, a2, b2 Point[T]) bool {
	corners, others := r.corners(), rectangle.corners()

	if rectangle.containsWithin(corners[0], a2, b2) || r.containsWithin(others[0], a1, b1) {
		return true
	}

	return r.edgesMeet(corners, others, a2, b2)
}

// edgesMeet reports whether an edge between the corners meets an edge between the others, whose
// extent is a, b. It reads only its arguments, so the receiver is unnamed.
func (Rectangle[T]) edgesMeet(corners, others [4]Point[T], a, b Point[T]) bool {
	probe := edgeProbe[T]{a: a, b: b}
	for edge := range edgesOf(corners[:]) {
		if !probe.aim(edge) {
			continue
		}

		for other := range edgesOf(others[:]) {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// parallel reports whether the given rectangle is turned by the same angle as this one, up to a
// half turn: a corner turned by the difference moves no farther than the tolerance at its
// coordinates. The tolerance is on that distance, not the angle, since any fixed angle moves the
// far corner of a long rectangle off its edge; for integer T only an exactly equal or opposite
// angle is parallel.
func (r Rectangle[T]) parallel(rectangle Rectangle[T]) bool {
	w1, h1 := r.Size.Float().XY()
	w2, h2 := rectangle.Size.Float().XY()
	reach := float64(max(math.Hypot(w1, h1), math.Hypot(w2, h2)) / 2)
	magnitude := max(r.Center.magnitude(), rectangle.Center.magnitude()) + reach
	turn := AngleDistance(r.Angle, rectangle.Angle)

	return min(turn, math.Pi-turn)*reach <= epsilonAt[T](magnitude)
}

// Equal reports whether the rectangles have equal center, size and angle, the angles compared with
// EqualAngle.
func (r Rectangle[T]) Equal(rectangle Rectangle[T]) bool {
	return r.Center.Equal(rectangle.Center) && r.Size.Equal(rectangle.Size) && EqualAngle(r.Angle, rectangle.Angle)
}

// IsZero reports whether the rectangle is the zero value, an angle of a full turn counting as zero.
func (r Rectangle[T]) IsZero() bool {
	return r.Equal(Rectangle[T]{})
}

// IsAligned reports whether the angle is exactly zero. No tolerance is applied, since a residue
// still turns the corners off the lattice; Canonical snaps it to zero.
func (r Rectangle[T]) IsAligned() bool {
	return r.Angle == 0
}

// Polygon returns the rectangle as a polygon of its Vertices.
func (r Rectangle[T]) Polygon() Polygon[T] {
	corners := r.corners()

	return Polygon[T]{corners[:]}
}

// Cast converts the rectangle to a Rectangle of another number type, rounding as Cast does and
// keeping the angle.
func (r Rectangle[T]) Cast[R Number]() Rectangle[R] {
	return Rectangle[R]{r.Center.Cast[R](), r.Size.Cast[R](), r.Angle}
}

// Int converts the rectangle to a Rectangle[int], rounding the center and the size separately, so
// the size is kept but the corners can land a unit from where they would round on their own.
func (r Rectangle[T]) Int() Rectangle[int] {
	return Rectangle[int]{r.Center.Int(), r.Size.Int(), r.Angle}
}

// Float converts the rectangle to a Rectangle[float64].
func (r Rectangle[T]) Float() Rectangle[float64] {
	return Rectangle[float64]{r.Center.Float(), r.Size.Float(), r.Angle}
}

// String returns the rectangle as Rect((x,y);WxH), center then size, with the angle appended as
// Rect((x,y);WxH;a) for a rotated rectangle.
func (r Rectangle[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Rect("...)
	buf = r.Center.appendString(buf)
	buf = append(buf, ';')
	buf = r.Size.appendString(buf)
	if !r.IsAligned() {
		buf = append(buf, ';')
		buf = appendString(buf, r.Angle)
	}
	buf = append(buf, ')')

	return string(buf)
}

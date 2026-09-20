package geom

import (
	"fmt"
	"iter"
	"math"
)

// Rectangle is a 2D rectangle represented by its center, size and angle: the axis-aligned box
// of that size about the center, turned by Angle radians about it in the same sense as
// Vector.Rotate. Rect builds one with no turn, and Rotate turns it.
//
// The size is never negative: Rect and Resize take it absolute, the corner constructors reorder
// the corners they are given, and Scale flips a negative factor into a positive one, so a
// rectangle with a negative extent can only be written as a struct literal or decoded from JSON,
// and Contains, Clamp and the Intersects methods give no meaningful answer for it; Canonical repairs it.
//
// Width, Height, Area and the corners are those of the rectangle itself, named in its own frame
// before the turn: TopLeft is the corner that was top-left before the rectangle was rotated,
// wherever it now lies, and Inset moves the edges in that frame. Min, Max and Bounds are the
// axis-aligned extent of the vertices, the box Bounds gives on every shape, and coincide with
// the corners only while the rectangle is not rotated.
//
// The rectangle is closed: Contains and the Intersects methods include the boundary, within the
// Epsilon that Equal applies, so a float rectangle contains the corners it was built from even
// where Min is recomputed from Center with a rounding error; Clamp applies no tolerance.
// For integer T the corners are lattice points on that boundary, so a rectangle of width w
// spans w+1 lattice columns from Min to Max inclusive, and a rotated rectangle has each corner
// rounded onto the lattice, so its edges are the segments between those rounded corners. The
// image.Rectangle returned by Rectangle is half-open as the image package requires, and
// therefore spans exactly w pixels.
type Rectangle[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	Angle  float64  `json:"a,omitzero"`
}

// Rect is shorthand for Rectangle{center, size}, with the size taken absolute and no turn.
func Rect[T Number](center Point[T], size Size[T]) Rectangle[T] {
	return Rectangle[T]{center, size.Abs(), 0}
}

// RectangleFromMin creates a Rectangle from min point and size. A negative extent measures the
// other way, so the given point is a corner but no longer the minimum one.
func RectangleFromMin[T Number](min Point[T], size Size[T]) Rectangle[T] {
	return RectangleFromMinMax(min, min.Add(size.Vector()))
}

// RectangleFromMax creates a Rectangle from max point and size. A negative extent measures the
// other way, so the given point is a corner but no longer the maximum one.
func RectangleFromMax[T Number](max Point[T], size Size[T]) Rectangle[T] {
	return RectangleFromMinMax(max.Add(size.Vector().Negate()), max)
}

// RectangleFromMinMax creates a Rectangle from two opposite corners, given in either order.
// For integer T the center truncates toward Min, so Max-Min is exactly the size.
func RectangleFromMinMax[T Number](a, b Point[T]) Rectangle[T] {
	a, b = Point[T]{min(a.X, b.X), min(a.Y, b.Y)}, Point[T]{max(a.X, b.X), max(a.Y, b.Y)}
	w, h := b.Subtract(a).XY()

	return Rectangle[T]{a.AddXY(w/2, h/2), Size[T]{w, h}, 0}
}

// RectangleFromSize creates a Rectangle from zero point and size.
func RectangleFromSize[T Number](size Size[T]) Rectangle[T] {
	return RectangleFromMin(Pt[T](0, 0), size)
}

// Width returns the rectangle width, the extent along its own X axis whatever its angle.
func (r Rectangle[T]) Width() T {
	return r.Size.Width
}

// Height returns the rectangle height, the extent along its own Y axis whatever its angle.
func (r Rectangle[T]) Height() T {
	return r.Size.Height
}

// Min returns the minimum corner of the axis-aligned extent of the rectangle: the top-left
// corner while it is not rotated, and the minimum of its four vertices once it is.
// RectangleFromMin is its inverse for a rectangle that is not rotated.
func (r Rectangle[T]) Min() Point[T] {
	a, _ := r.MinMax()

	return a
}

// Max returns the maximum corner of the axis-aligned extent of the rectangle: the bottom-right
// corner while it is not rotated, and the maximum of its four vertices once it is.
// RectangleFromMax is its inverse for a rectangle that is not rotated.
func (r Rectangle[T]) Max() Point[T] {
	_, b := r.MinMax()

	return b
}

// MinMax returns the minimum and maximum corners of the axis-aligned extent together, the
// corners of Bounds, so a rotated rectangle is rejected against another shape by the same box
// test as any other outline.
func (r Rectangle[T]) MinMax() (Point[T], Point[T]) {
	if r.IsAligned() {
		a, b := r.localMinMax()

		return r.Center.Add(a), r.Center.Add(b)
	}

	corners := r.corners()

	return minMaxOf(corners[:])
}

// TopLeft returns the top-left corner, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) TopLeft() Point[T] {
	a, _ := r.localMinMax()

	return r.worldPoint(a)
}

// BottomLeft returns the bottom-left corner, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) BottomLeft() Point[T] {
	a, b := r.localMinMax()

	return r.worldPoint(Vector[T]{a.X, b.Y})
}

// BottomRight returns the bottom-right corner, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) BottomRight() Point[T] {
	_, b := r.localMinMax()

	return r.worldPoint(b)
}

// TopRight returns the top-right corner, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) TopRight() Point[T] {
	a, b := r.localMinMax()

	return r.worldPoint(Vector[T]{b.X, a.Y})
}

// Top returns the midpoint of the top edge, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) Top() Point[T] {
	a, _ := r.localMinMax()

	return r.worldPoint(Vector[T]{0, a.Y})
}

// Bottom returns the midpoint of the bottom edge, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) Bottom() Point[T] {
	_, b := r.localMinMax()

	return r.worldPoint(Vector[T]{0, b.Y})
}

// Left returns the midpoint of the left edge, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) Left() Point[T] {
	a, _ := r.localMinMax()

	return r.worldPoint(Vector[T]{a.X, 0})
}

// Right returns the midpoint of the right edge, named in the frame of the rectangle before its turn.
func (r Rectangle[T]) Right() Point[T] {
	_, b := r.localMinMax()

	return r.worldPoint(Vector[T]{b.X, 0})
}

// Anchor returns the point on the rectangle in the given direction from its center, in the
// frame of the rectangle before its turn: a corner for diagonals and the midpoint of an edge
// for cardinals, so the Top anchor of a rotated rectangle is the midpoint of the edge that was
// on top before the turn.
func (r Rectangle[T]) Anchor(direction Direction) Point[T] {
	switch direction.normalize() {
	case TopLeft:
		return r.TopLeft()
	case Top:
		return r.Top()
	case TopRight:
		return r.TopRight()
	case Right:
		return r.Right()
	case BottomRight:
		return r.BottomRight()
	case Bottom:
		return r.Bottom()
	case BottomLeft:
		return r.BottomLeft()
	case Left:
		return r.Left()
	default:
		return r.Center
	}
}

// TopEdge returns the top edge, from the top-left to the top-right corner.
func (r Rectangle[T]) TopEdge() Line[T] {
	return Ln(r.TopLeft(), r.TopRight())
}

// RightEdge returns the right edge, from the top-right to the bottom-right corner.
func (r Rectangle[T]) RightEdge() Line[T] {
	return Ln(r.TopRight(), r.BottomRight())
}

// BottomEdge returns the bottom edge, from the bottom-right to the bottom-left corner.
func (r Rectangle[T]) BottomEdge() Line[T] {
	return Ln(r.BottomRight(), r.BottomLeft())
}

// LeftEdge returns the left edge, from the bottom-left to the top-left corner.
func (r Rectangle[T]) LeftEdge() Line[T] {
	return Ln(r.BottomLeft(), r.TopLeft())
}

// Edges iterates the rectangle edges in order starting at the top-left corner, by increasing
// angle — the same winding as Directions and RegularPolygon.Points, and clockwise as drawn on
// a screen with Y pointing down. Each edge starts where the previous one ends. The edges are
// built from the corners as they are walked, without allocating; collect them with
// slices.Collect where a slice is needed. Every walk over the outline, containment, distance
// and the crossings of a segment, reads these edges, so the boundary they join is the one
// every test agrees on.
func (r Rectangle[T]) Edges() iter.Seq[Line[T]] {
	return func(yield func(Line[T]) bool) {
		corners := r.corners()

		edgesOf(corners[:])(yield)
	}
}

// Vertices iterates the rectangle vertices in order starting at the top-left corner, by
// increasing angle — the same winding as Directions and RegularPolygon.Points, and clockwise
// as drawn on a screen with Y pointing down, without allocating; collect them with
// slices.Collect where a slice is needed. For integer T each vertex of a rotated rectangle is
// rounded onto the lattice.
func (r Rectangle[T]) Vertices() iter.Seq[Point[T]] {
	return func(yield func(Point[T]) bool) {
		for _, corner := range r.corners() {
			if !yield(corner) {
				return
			}
		}
	}
}

// Area returns the rectangle area.
func (r Rectangle[T]) Area() T {
	return r.Size.Area()
}

// Perimeter returns the rectangle perimeter.
func (r Rectangle[T]) Perimeter() T {
	return r.Size.Perimeter()
}

// AspectRatio returns width/height.
func (r Rectangle[T]) AspectRatio() float64 {
	return r.Size.AspectRatio()
}

// Bounds returns the axis-aligned bounding rectangle: the rectangle itself while it is not
// rotated, and the box from Min to Max once it is.
func (r Rectangle[T]) Bounds() Rectangle[T] {
	if r.IsAligned() {
		return r
	}

	return RectangleFromMinMax(r.MinMax())
}

// localMinMax returns the offsets of the minimum and maximum corner from the center in the
// local frame of the rectangle, the frame before its turn, the pair MinMax gives in the world.
// For integer T the minimum offset truncates toward the center and the maximum takes the rest,
// so the two differ by exactly the size.
func (r Rectangle[T]) localMinMax() (Vector[T], Vector[T]) {
	w, h := r.Size.XY()

	return Vector[T]{-w / 2, -h / 2}, Vector[T]{w - w/2, h - h/2}
}

// worldPoint returns the point at the given offset from the center in the frame before the turn,
// turned by Angle: the sum alone while the rectangle is not rotated, so an integer rectangle
// keeps its exact corners, and the offset rotated once and rounded otherwise.
func (r Rectangle[T]) worldPoint(offset Vector[T]) Point[T] {
	if r.IsAligned() {
		return r.Center.Add(offset)
	}

	return r.Center.Add(offset.Rotate(r.Angle))
}

// localOffset returns the offset of the point from the center in the frame of the rectangle before
// its turn, in float64, the inverse of world, so that Clamp can clamp a rotated rectangle as
// an aligned one.
func (r Rectangle[T]) localOffset(point Point[T]) Vector[float64] {
	return point.Subtract(r.Center).Float().Rotate(-r.Angle)
}

// localRectangle returns the given rectangle as it lies in the local frame of this one, the
// frame before its turn: the offset between the centers turned back by Angle, about the origin,
// with no angle of its own. Two rectangles of the same angle taken into the local frame of the
// same rectangle, itself included, are aligned about the origin and overlap, intersect and
// unite as aligned rectangles do; worldRectangle turns the result back. For integer T the
// turned offset is rounded.
func (r Rectangle[T]) localRectangle(rectangle Rectangle[T]) Rectangle[T] {
	offset := rectangle.Center.Subtract(r.Center).Rotate(-r.Angle)

	return Rectangle[T]{offset.Point(), rectangle.Size, 0}
}

// worldRectangle returns a rectangle found in the local frame of this one, as localRectangle
// gives it, turned back into the world: its center turned by Angle about this center, with
// this angle.
func (r Rectangle[T]) worldRectangle(rectangle Rectangle[T]) Rectangle[T] {
	return Rectangle[T]{r.worldPoint(rectangle.Center.Vector()), rectangle.Size, r.Angle}
}

// corners returns the vertices as an array from a single frame, for Vertices, Edges and
// Polygon, so none recomputes the corners it shares.
func (r Rectangle[T]) corners() [4]Point[T] {
	a, b := r.localMinMax()

	return [4]Point[T]{r.worldPoint(a), r.worldPoint(Vector[T]{b.X, a.Y}), r.worldPoint(b), r.worldPoint(Vector[T]{a.X, b.Y})}
}

// Translate creates a new Rectangle translated by the given vector.
func (r Rectangle[T]) Translate(vector Vector[T]) Rectangle[T] {
	return Rectangle[T]{r.Center.Add(vector), r.Size, r.Angle}
}

// MoveTo creates a new Rectangle with the same size and angle centered at point.
func (r Rectangle[T]) MoveTo(point Point[T]) Rectangle[T] {
	return Rectangle[T]{point, r.Size, r.Angle}
}

// Scale creates a new Rectangle with size uniformly scaled by the factor. A negative factor
// scales by its absolute value, since a rectangle mirrored about its center is the same rectangle.
func (r Rectangle[T]) Scale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Scale(factor).Abs(), r.Angle}
}

// ScaleXY creates a new Rectangle with size scaled by the given factors along its own axes,
// negative ones by their absolute value like Scale.
func (r Rectangle[T]) ScaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ScaleXY(factorX, factorY).Abs(), r.Angle}
}

// Unscale creates a new Rectangle with size uniformly scaled by the inverse factor, the
// inverse of Scale and negative factors taken absolute like it. Like Divide it panics for a
// zero factor.
func (r Rectangle[T]) Unscale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Unscale(factor).Abs(), r.Angle}
}

// UnscaleXY creates a new Rectangle with size scaled by the inverse of the given factors, the
// inverse of ScaleXY. Like Divide it panics for a zero factor.
func (r Rectangle[T]) UnscaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.UnscaleXY(factorX, factorY).Abs(), r.Angle}
}

// Resize creates a new Rectangle with the given size, taken absolute like Rect.
func (r Rectangle[T]) Resize(size Size[T]) Rectangle[T] {
	return Rectangle[T]{r.Center, size.Abs(), r.Angle}
}

// Canonical creates a new Rectangle in the form Rect and Rotate build, with the size taken
// absolute and the angle normalized to [0, 2π): a rectangle mirrored about its own center is
// the same rectangle, so the center is untouched and a well-formed rectangle is returned as it
// is, up to the full turns Equal already ignores. It repairs a struct literal or decoded JSON
// with a negative extent before Contains, Clamp or the Intersects methods read it, and brings
// a decoded angle onto the seam Rotate keeps. An angle within Delta of zero or of a full turn,
// the residue a chain of Rotate and Lerp calls can leave, becomes exactly zero, so the
// rectangle IsAligned again and its corners are exact; Rotate itself never snaps, since a turn
// that small still moves a far corner of a large rectangle by more than Epsilon.
func (r Rectangle[T]) Canonical() Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Abs(), snapAngle(r.Angle)}
}

// Grow creates a new Rectangle with size expanded by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, so each side moves out by half of it, where
// Outset moves every side by the full padding.
// For integer T an odd amount cannot be split evenly and lands on one side only: on Min when
// the new extent is even, on Max when it is odd, as Min and Max place the center. Use Outset
// to move a chosen side by a whole amount.
func (r Rectangle[T]) Grow(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Grow(amount).AtLeast(Size[T]{}), r.Angle}
}

// GrowXY creates a new Rectangle with size expanded by the given amounts along its own axes, clamped to zero.
// Each amount is the total change of that extent and an odd integer one lands on one side only, like Grow.
func (r Rectangle[T]) GrowXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.GrowXY(amountX, amountY).AtLeast(Size[T]{}), r.Angle}
}

// Shrink creates a new Rectangle with size reduced by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, so each side moves in by half of it, where
// Inset moves every side by the full padding. An odd integer amount comes off one side only, like Grow.
func (r Rectangle[T]) Shrink(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Shrink(amount).AtLeast(Size[T]{}), r.Angle}
}

// ShrinkXY creates a new Rectangle with size reduced by the given amounts along its own axes, clamped to zero.
// Each amount is the total change of that extent and an odd integer one comes off one side only, like Shrink.
func (r Rectangle[T]) ShrinkXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ShrinkXY(amountX, amountY).AtLeast(Size[T]{}), r.Angle}
}

// Inset creates a new Rectangle inset by the given padding amounts, each edge moved in the
// frame of the rectangle before its turn: Left moves the edge that was on the left, wherever
// the turn has put it, and the center follows along the turned axes. An edge pushed past its
// opposite stops there, so a padding larger than the rectangle collapses it to a zero extent on
// that axis, at the opposite edge rather than at the one that over-ran: too much Left collapses
// it onto the right edge, too much Right onto the left one. Where both paddings on an axis
// over-run, Left and Top win, since the minimum corner moves first and the maximum is then
// stopped at it.
// A negative padding outsets the rectangle, so Outset undoes Inset as long as nothing was clamped.
// For integer T the moved center of a rotated rectangle is rounded once, as worldRectangle places it.
func (r Rectangle[T]) Inset(padding Padding[T]) Rectangle[T] {
	a, b := r.localMinMax()

	a = Vector[T]{min(a.X+padding.Left, b.X), min(a.Y+padding.Top, b.Y)}
	b = Vector[T]{max(b.X-padding.Right, a.X), max(b.Y-padding.Bottom, a.Y)}
	inset := RectangleFromMinMax(a.Point(), b.Point())

	return r.worldRectangle(inset)
}

// Outset creates a new Rectangle expanded by the given padding amounts, the inverse of Inset.
func (r Rectangle[T]) Outset(padding Padding[T]) Rectangle[T] {
	return r.Inset(padding.Negate())
}

// Lerp creates a new Rectangle in linear interpolation towards the given rectangle, moving the
// center and the size together and turning the angle along the shorter arc with LerpAngle,
// normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like Point.Lerp, with the
// size taken absolute, so an extrapolation that would pass through a mirrored rectangle stays
// a rectangle.
func (r Rectangle[T]) Lerp(rectangle Rectangle[T], t float64) Rectangle[T] {
	return Rectangle[T]{r.Center.Lerp(rectangle.Center, t), r.Size.Lerp(rectangle.Size, t).Abs(), NormalizeAngle(LerpAngle(r.Angle, rectangle.Angle, t))}
}

// Rotate creates a new Rectangle turned by the given angle (in radians) about its center, in
// the same sense as Vector.Rotate. The stored angle is normalized to [0, 2π) to prevent drift
// from repeated rotations, so a rectangle turned back by its own angle is not rotated at all
// and its corners are exact again.
func (r Rectangle[T]) Rotate(angle float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size, NormalizeAngle(r.Angle + angle)}
}

// AlignTo creates a new Rectangle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (r Rectangle[T]) AlignTo(direction Direction, point Point[T]) Rectangle[T] {
	return r.Translate(point.Subtract(r.Anchor(direction)))
}

// Clamp returns the given Point clamped to the rectangle: the point itself inside, and the
// nearest point of the rectangle outside, with no tolerance. A rotated rectangle clamps in its
// own frame, so the result lies on the edge nearest to the point rather than in the box from
// Min to Max; for integer T it is rounded once, after the turn back.
func (r Rectangle[T]) Clamp(point Point[T]) Point[T] {
	if r.IsAligned() {
		a, b := r.MinMax()

		return Point[T]{Clamp(point.X, a.X, b.X), Clamp(point.Y, a.Y, b.Y)}
	}

	a, b := r.localMinMax()
	local := r.localOffset(point)
	clamped := Vector[float64]{Clamp(local.X, float64(a.X), float64(b.X)), Clamp(local.Y, float64(a.Y), float64(b.Y))}
	placed := r.Center.Float().Add(clamped.Rotate(r.Angle))

	return placed.Cast[T]()
}

// Contains reports whether the given point lies within the rectangle, boundary included within
// Epsilon of T: strictly inside, or on an edge as Line.Contains judges it, so the rectangle
// contains exactly the points its edges contain and the points between them.
func (r Rectangle[T]) Contains(point Point[T]) bool {
	return r.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the rectangle:
// zero exactly where Contains holds, so a point within Epsilon of T of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the nearest edge.
func (r Rectangle[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(r.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons, in one pass over the edges, the edgeWalk every closed shape makes: zero for a
// point inside by the even-odd rule, or on an edge within Epsilon of T as
// Line.DistanceSquaredTo snaps it, and the squared distance to the nearest edge otherwise. A
// rectangle that is not rotated answers a point Clamp leaves where it is before any edge is
// examined. It is a float64 even for an integer T, since the nearest point of a rotated
// rectangle is a foot on a turned edge, not a lattice point in general; only
// Point.DistanceSquaredTo stays in T. Contains and IntersectsCircle are built on it, and the
// edges are the ones Line.IntersectionLine and Polygon read, so containment, the boundary
// crossings of a segment and the polygon of the rectangle agree by construction, for a
// rotated integer rectangle on the rounded corners its edges join.
func (r Rectangle[T]) DistanceSquaredTo(point Point[T]) float64 {
	if r.IsAligned() && r.Clamp(point) == point {
		return 0
	}

	w := edgeWalk[T]{distance: math.Inf(1)}
	for edge := range r.Edges() {
		if w.step(edge, point) {
			return 0
		}
	}

	return w.result()
}

// IntersectsCircle reports whether the rectangle and the circle overlap, as
// Circle.IntersectsRectangle does.
func (r Rectangle[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.IntersectsRectangle(r)
}

// IntersectsLine reports whether the rectangle and the segment share a point, as
// Line.IntersectsRectangle does.
func (r Rectangle[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsRectangle(r)
}

// IntersectionLine returns the points where the segment crosses the rectangle boundary, as
// Line.IntersectionRectangle does.
func (r Rectangle[T]) IntersectionLine(line Line[T]) []Point[T] {
	return line.IntersectionRectangle(r)
}

// IntersectsPolygon reports whether the rectangle and the polygon share a point, as
// Polygon.IntersectsRectangle does.
func (r Rectangle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsRectangle(r)
}

// IntersectsRectangle reports whether the rectangles share a point: a corner of one lies within the
// other, or an edge of one meets an edge of the other, as Polygon.IntersectsPolygon decides and by
// the same tolerance on the distance, so touching rectangles intersect within Epsilon of T,
// the same closed convention as Contains. Rectangles whose extents do not overlap are rejected
// before any edge is examined, and two rectangles that are not rotated whose extents overlap
// exactly are decided there, since an exact overlap of two aligned boxes always shares a
// corner or a crossing; only a gap within the tolerance goes to the edges. Two rectangles of
// the same angle are tested in their shared frame, as Intersection finds their overlap.
func (r Rectangle[T]) IntersectsRectangle(rectangle Rectangle[T]) bool {
	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	switch {
	case !overlaps(a1, b1, a2, b2):
		return false
	case r.IsAligned() && rectangle.IsAligned() && a1.X <= b2.X && a2.X <= b1.X && a1.Y <= b2.Y && a2.Y <= b1.Y:
		return true
	case r.IsAligned() && rectangle.IsAligned():
		return r.meets(rectangle)
	case r.parallel(rectangle):
		return r.localRectangle(r).IntersectsRectangle(r.localRectangle(rectangle))
	default:
		return r.meets(rectangle)
	}
}

// IntersectionRectangle returns the rectangle common to both, and false when they do not intersect.
// Touching rectangles intersect in a rectangle of zero width or height, within Epsilon of T,
// exactly where Intersects holds: a corner admitted by the tolerance is placed on the boundary
// of the other rectangle, never beyond it. Two rectangles of the same angle overlap in a
// rectangle of that angle, found in their shared frame; rectangles of different angles overlap
// in a polygon that is not a rectangle and return false even where Intersects holds, as
// parallel segments do for Line.IntersectionLine. For integer T the offset between the centers of
// two rotated rectangles is rounded into the shared frame and the result rounded back.
func (r Rectangle[T]) IntersectionRectangle(rectangle Rectangle[T]) (Rectangle[T], bool) {
	switch {
	case r.IsAligned() && rectangle.IsAligned():
		if !r.IntersectsRectangle(rectangle) {
			return Rectangle[T]{}, false
		}

		a1, b1 := r.MinMax()
		a2, b2 := rectangle.MinMax()

		a := Point[T]{max(a1.X, a2.X), max(a1.Y, a2.Y)}
		b := Point[T]{max(min(b1.X, b2.X), a.X), max(min(b1.Y, b2.Y), a.Y)}

		return RectangleFromMinMax(a, b), true
	case r.parallel(rectangle):
		overlap, ok := r.localRectangle(r).IntersectionRectangle(r.localRectangle(rectangle))
		if !ok {
			return Rectangle[T]{}, false
		}

		return r.worldRectangle(overlap), true
	default:
		return Rectangle[T]{}, false
	}
}

// Union returns the smallest rectangle containing both: of their shared angle for two
// rectangles of the same angle, found in their shared frame, and the axis-aligned box around
// the Bounds of both for rectangles of different angles, which no rectangle of either angle
// bounds tightly. For integer T the offset between the centers of two rotated rectangles is
// rounded into the shared frame and the result rounded back.
func (r Rectangle[T]) Union(rectangle Rectangle[T]) Rectangle[T] {
	switch {
	case r.IsAligned() && rectangle.IsAligned():
		a1, b1 := r.MinMax()
		a2, b2 := rectangle.MinMax()

		return RectangleFromMinMax(
			Point[T]{min(a1.X, a2.X), min(a1.Y, a2.Y)},
			Point[T]{max(b1.X, b2.X), max(b1.Y, b2.Y)},
		)
	case r.parallel(rectangle):
		return r.worldRectangle(r.localRectangle(r).Union(r.localRectangle(rectangle)))
	default:
		return r.Bounds().Union(rectangle.Bounds())
	}
}

// IntersectsRegularPolygon reports whether the rectangle and the regular polygon share a
// point, the answer IntersectsPolygon gives on the polygon's Polygon form, without building
// it: a corner of one lies within the other, or an edge of the polygon crosses an edge of the
// rectangle. The two Bounds reject the pair before any edge is examined, whatever the
// rectangle's angle, and an empty polygon intersects nothing.
func (r Rectangle[T]) IntersectsRegularPolygon(polygon RegularPolygon[T]) bool {
	if polygon.Empty() {
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

// containsWithin is Contains for a caller that already holds the corners of MinMax, so the
// intersection tests read the extent once and reuse it for every point they test.
func (r Rectangle[T]) containsWithin(point, a, b Point[T]) bool {
	return point.Between(a, b) && r.DistanceSquaredTo(point) == 0
}

// meets reports whether two rectangles whose extents overlap share a point, the way
// Polygon.IntersectsPolygon decides it: a corner of one lies within the other, or an edge of one
// meets an edge of the other by Line.IntersectsLine.
func (r Rectangle[T]) meets(rectangle Rectangle[T]) bool {
	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	if rectangle.containsWithin(r.TopLeft(), a2, b2) || r.containsWithin(rectangle.TopLeft(), a1, b1) {
		return true
	}

	probe := edgeProbe[T]{a: a2, b: b2}
	for edge := range r.Edges() {
		if !probe.aim(edge) {
			continue
		}

		for other := range rectangle.Edges() {
			if probe.meets(other) {
				return true
			}
		}
	}

	return false
}

// parallel reports whether the given rectangle is turned by the same angle as this one, up to
// a full turn as EqualAngle judges it, so the two share a frame.
func (r Rectangle[T]) parallel(rectangle Rectangle[T]) bool {
	return EqualAngle(r.Angle, rectangle.Angle)
}

// Equal checks for equal center, size and angle values using tolerant numeric comparison.
// Angles are compared with EqualAngle, so a full turn or the sign of an angle does not matter.
func (r Rectangle[T]) Equal(rectangle Rectangle[T]) bool {
	return r.Center.Equal(rectangle.Center) && r.Size.Equal(rectangle.Size) && EqualAngle(r.Angle, rectangle.Angle)
}

// IsZero checks if center point, size and angle are zero, comparing the angle like Equal so
// that a full turn counts as zero.
func (r Rectangle[T]) IsZero() bool {
	return r.Equal(Rectangle[T]{})
}

// IsAligned reports whether the rectangle is axis-aligned: its Angle is exactly zero, as Rect,
// the corner constructors, Bounds and Rotate by a full turn leave it. No tolerance is applied,
// since an angle a rounding error from zero still turns the corners off the lattice; Canonical
// snaps such a residue to zero. An aligned rectangle has exact corners, coincides with its
// Bounds, and takes the fast path of every test.
func (r Rectangle[T]) IsAligned() bool {
	return r.Angle == 0
}

// Polygon converts the rectangle into a generic Polygon with the vertices Vertices iterates,
// in one allocation.
func (r Rectangle[T]) Polygon() Polygon[T] {
	corners := r.corners()

	return Polygon[T]{corners[:]}
}

// Cast converts the rectangle to a Rectangle of another number type, rounding as Cast does and
// keeping the angle.
func (r Rectangle[T]) Cast[R Number]() Rectangle[R] {
	return Rectangle[R]{r.Center.Cast[R](), r.Size.Cast[R](), r.Angle}
}

// Int converts the rectangle to a Rectangle[int], rounding the center and the size on their
// own like every other Int, so the size is exact and a rectangle keeps its extent as it moves
// through positions no lattice expresses. The box moves instead: a center on a half rounds
// away from zero and the int center then truncates toward Min, so the box can land a whole
// unit from where the corners would round. Rectangle builds the image.Rectangle from it, so a
// float rectangle covers the same pixels its Int spans. The angle is kept.
func (r Rectangle[T]) Int() Rectangle[int] {
	return Rectangle[int]{r.Center.Int(), r.Size.Int(), r.Angle}
}

// Float converts the rectangle to a Rectangle[float64].
func (r Rectangle[T]) Float() Rectangle[float64] {
	return Rectangle[float64]{r.Center.Float(), r.Size.Float(), r.Angle}
}

// String returns the rectangle in the form of its constructor: Rect((x,y);WxH), center then
// size, with the angle appended as Rect((x,y);WxH;a) for a rotated rectangle, as the JSON
// carries it only then. See MinMaxString for the corners.
func (r Rectangle[T]) String() string {
	if r.IsAligned() {
		return fmt.Sprintf("Rect(%s;%s)", r.Center.String(), r.Size.String())
	}

	return fmt.Sprintf("Rect(%s;%s;%s)", r.Center.String(), r.Size.String(), String(r.Angle))
}

// MinMaxString returns the rectangle by its extent: (x,y)-(x,y), Min then Max. It is the form
// to read positions from, where String mirrors how the rectangle was built.
func (r Rectangle[T]) MinMaxString() string {
	return fmt.Sprintf("%s-%s", r.Min().String(), r.Max().String())
}

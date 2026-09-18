package geom

import (
	"fmt"
	"iter"
	"math"
	"slices"
)

// Rectangle is a 2D axis-aligned rectangle represented by its center and size.
//
// The size is never negative: Rect and Resize take it absolute, the corner constructors reorder
// the corners they are given, and Scale flips a negative factor into a positive one, so a
// rectangle with Min beyond Max can only be written as a struct literal or decoded from JSON,
// and Contains, Clamp and the Intersects methods give no meaningful answer for it.
//
// The rectangle is closed: Contains, Clamp and the Intersects methods include the boundary,
// within the Epsilon that Equal applies, so a float rectangle contains the corners it was built
// from even where Min is recomputed from Center with a rounding error.
// For integer T the corners are lattice points on that boundary, so a rectangle of width w
// spans w+1 lattice columns from Min to Max inclusive. The image.Rectangle returned by Rectangle
// is half-open as the image package requires, and therefore spans exactly w pixels.
type Rectangle[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
}

// Rect is shorthand for Rectangle{center, size}, with the size taken absolute.
func Rect[T Number](center Point[T], size Size[T]) Rectangle[T] {
	return Rectangle[T]{center, size.Abs()}
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

	return Rectangle[T]{a.AddXY(w/2, h/2), Size[T]{w, h}}
}

// RectangleFromSize creates a Rectangle from zero point and size.
func RectangleFromSize[T Number](size Size[T]) Rectangle[T] {
	return RectangleFromMin(Pt[T](0, 0), size)
}

// Width returns the rectangle width.
func (r Rectangle[T]) Width() T {
	return r.Size.Width
}

// Height returns the rectangle height.
func (r Rectangle[T]) Height() T {
	return r.Size.Height
}

// Min returns the minimum corner point of the rectangle. RectangleFromMin is its inverse.
func (r Rectangle[T]) Min() Point[T] {
	w, h := r.Size.XY()

	return r.Center.AddXY(-w/2, -h/2)
}

// Max returns the maximum corner point of the rectangle.
// For integer types with odd Width or Height, w-w/2 != w/2 due to truncation;
// using (w-w/2) here keeps Max-Min equal to exactly w (the extra unit goes to Max).
func (r Rectangle[T]) Max() Point[T] {
	w, h := r.Size.XY()

	return r.Center.AddXY(w-w/2, h-h/2)
}

// MinMax returns the minimum and maximum corner points of the rectangle together.
func (r Rectangle[T]) MinMax() (Point[T], Point[T]) {
	return r.Min(), r.Max()
}

// TopLeft returns the top-left corner.
func (r Rectangle[T]) TopLeft() Point[T] {
	return r.Min()
}

// BottomLeft returns the bottom-left corner.
func (r Rectangle[T]) BottomLeft() Point[T] {
	return Point[T]{r.Min().X, r.Max().Y}
}

// BottomRight returns the bottom-right corner.
func (r Rectangle[T]) BottomRight() Point[T] {
	return r.Max()
}

// TopRight returns the top-right corner.
func (r Rectangle[T]) TopRight() Point[T] {
	return Point[T]{r.Max().X, r.Min().Y}
}

// Top returns the midpoint of the top edge.
func (r Rectangle[T]) Top() Point[T] {
	return Point[T]{r.Center.X, r.Min().Y}
}

// Bottom returns the midpoint of the bottom edge.
func (r Rectangle[T]) Bottom() Point[T] {
	return Point[T]{r.Center.X, r.Max().Y}
}

// Left returns the midpoint of the left edge.
func (r Rectangle[T]) Left() Point[T] {
	return Point[T]{r.Min().X, r.Center.Y}
}

// Right returns the midpoint of the right edge.
func (r Rectangle[T]) Right() Point[T] {
	return Point[T]{r.Max().X, r.Center.Y}
}

// Anchor returns the point on the rectangle in the given direction from its center:
// a corner for diagonals and the midpoint of an edge for cardinals.
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

// Edges returns the rectangle edges as lines in order starting Min point, by increasing angle —
// the same winding as Directions and RegularPolygon.Vertices, and clockwise as drawn on a screen
// with Y pointing down. Each edge starts where the previous one ends.
func (r Rectangle[T]) Edges() []Line[T] {
	return slices.AppendSeq(make([]Line[T], 0, 4), r.edges())
}

// edges iterates the edges Edges returns without allocating them, for the methods that only
// need to walk them once.
func (r Rectangle[T]) edges() iter.Seq[Line[T]] {
	return func(yield func(Line[T]) bool) {
		corners := r.corners()
		for i, corner := range corners {
			if !yield(Line[T]{corner, corners[(i+1)%4]}) {
				return
			}
		}
	}
}

// Vertices returns the rectangle vertices in order starting Min point, by increasing angle —
// the same winding as Directions and RegularPolygon.Vertices, and clockwise as drawn on a screen
// with Y pointing down.
func (r Rectangle[T]) Vertices() []Point[T] {
	corners := r.corners()

	return corners[:]
}

// corners returns the vertices as an array from a single Min and Max, for Vertices and the
// edge walk, so neither recomputes the corners it shares.
func (r Rectangle[T]) corners() [4]Point[T] {
	a, b := r.MinMax()

	return [4]Point[T]{a, {b.X, a.Y}, b, {a.X, b.Y}}
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

// Bounds returns the axis-aligned bounding rectangle.
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}

// Translate creates a new Rectangle translated by the given vector.
func (r Rectangle[T]) Translate(vector Vector[T]) Rectangle[T] {
	return Rectangle[T]{r.Center.Add(vector), r.Size}
}

// MoveTo creates a new Rectangle with the same size centered at point.
func (r Rectangle[T]) MoveTo(point Point[T]) Rectangle[T] {
	return Rectangle[T]{point, r.Size}
}

// Scale creates a new Rectangle with size uniformly scaled by the factor. A negative factor
// scales by its absolute value, since a rectangle mirrored about its center is the same rectangle.
func (r Rectangle[T]) Scale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Scale(factor).Abs()}
}

// ScaleXY creates a new Rectangle with size scaled by the given factors, negative ones by their
// absolute value like Scale.
func (r Rectangle[T]) ScaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ScaleXY(factorX, factorY).Abs()}
}

// Resize creates a new Rectangle with the given size, taken absolute like Rect.
func (r Rectangle[T]) Resize(size Size[T]) Rectangle[T] {
	return Rectangle[T]{r.Center, size.Abs()}
}

// Grow creates a new Rectangle with size expanded by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, so each side moves out by half of it, where
// Outset moves every side by the full padding.
// For integer T an odd amount cannot be split evenly and lands on one side only: on Min when
// the new extent is even, on Max when it is odd, as Min and Max place the center. Use Outset
// to move a chosen side by a whole amount.
func (r Rectangle[T]) Grow(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Grow(amount)}
}

// GrowXY creates a new Rectangle with size expanded by the given amounts along X and Y, clamped to zero.
// Each amount is the total change of that extent and an odd integer one lands on one side only, like Grow.
func (r Rectangle[T]) GrowXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.GrowXY(amountX, amountY)}
}

// Shrink creates a new Rectangle with size reduced by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, so each side moves in by half of it, where
// Inset moves every side by the full padding. An odd integer amount comes off one side only, like Grow.
func (r Rectangle[T]) Shrink(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Shrink(amount)}
}

// ShrinkXY creates a new Rectangle with size reduced by the given amounts along X and Y, clamped to zero.
// Each amount is the total change of that extent and an odd integer one comes off one side only, like Shrink.
func (r Rectangle[T]) ShrinkXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ShrinkXY(amountX, amountY)}
}

// Inset creates a new Rectangle inset by the given padding amounts. An edge pushed past its
// opposite stops there, so a padding larger than the rectangle collapses it to a zero extent
// that still lies within the original bounds, at the last edge to move.
// A negative padding outsets the rectangle, so Outset undoes Inset as long as nothing was clamped.
func (r Rectangle[T]) Inset(padding Padding[T]) Rectangle[T] {
	a, b := r.MinMax()

	a = Point[T]{min(a.X+padding.Left, b.X), min(a.Y+padding.Top, b.Y)}
	b = Point[T]{max(b.X-padding.Right, a.X), max(b.Y-padding.Bottom, a.Y)}

	return RectangleFromMinMax(a, b)
}

// Outset creates a new Rectangle expanded by the given padding amounts, the inverse of Inset.
func (r Rectangle[T]) Outset(padding Padding[T]) Rectangle[T] {
	return r.Inset(padding.Negate())
}

// AlignTo creates a new Rectangle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (r Rectangle[T]) AlignTo(direction Direction, point Point[T]) Rectangle[T] {
	return r.Translate(point.Subtract(r.Anchor(direction)))
}

// Clamp returns the given Point clamped to the rectangle bounds.
func (r Rectangle[T]) Clamp(point Point[T]) Point[T] {
	a, b := r.MinMax()

	return Point[T]{Clamp(point.X, a.X, b.X), Clamp(point.Y, a.Y, b.Y)}
}

// Contains reports whether the given point lies within the rectangle, boundary included within
// Epsilon of T.
func (r Rectangle[T]) Contains(point Point[T]) bool {
	return point.Between(r.MinMax())
}

// DistanceTo returns the distance from the given point to the nearest point of the rectangle:
// zero exactly where Contains holds, so a point within Epsilon of T of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the clamped point.
func (r Rectangle[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(float64(r.DistanceSquaredTo(point)))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons. It stays in T like Point.DistanceSquaredTo, since the nearest point is the
// clamped point and a lattice point for an integer T.
func (r Rectangle[T]) DistanceSquaredTo(point Point[T]) T {
	if r.Contains(point) {
		return 0
	}

	return r.Clamp(point).DistanceSquaredTo(point)
}

// Intersects reports whether the rectangles overlap. Touching rectangles intersect, within
// Epsilon of T, the same closed convention as Contains.
func (r Rectangle[T]) Intersects(rectangle Rectangle[T]) bool {
	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	return overlaps(a1, b1, a2, b2)
}

// Intersection returns the rectangle common to both, and false when they do not intersect.
// Touching rectangles intersect in a rectangle of zero width or height, within Epsilon of T,
// the same closed convention as Intersects.
func (r Rectangle[T]) Intersection(rectangle Rectangle[T]) (Rectangle[T], bool) {
	if !r.Intersects(rectangle) {
		return Rectangle[T]{}, false
	}

	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	a := Point[T]{max(a1.X, a2.X), max(a1.Y, a2.Y)}
	b := Point[T]{max(min(b1.X, b2.X), a.X), max(min(b1.Y, b2.Y), a.Y)}

	return RectangleFromMinMax(a, b), true
}

// Union returns the smallest rectangle containing both.
func (r Rectangle[T]) Union(rectangle Rectangle[T]) Rectangle[T] {
	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	return RectangleFromMinMax(
		Point[T]{min(a1.X, a2.X), min(a1.Y, a2.Y)},
		Point[T]{max(b1.X, b2.X), max(b1.Y, b2.Y)},
	)
}

// IntersectsCircle reports whether the rectangle and the circle overlap: the point of the
// rectangle closest to the circle center lies within the radius. Touching shapes intersect,
// within Epsilon of T, by the same comparison Circle.Contains makes, and the rectangle bounds
// are the same Min and Max that Contains uses.
func (r Rectangle[T]) IntersectsCircle(circle Circle[T]) bool {
	return circle.reaches(circle.Center.Float().DistanceSquaredTo(r.Clamp(circle.Center).Float()))
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

// Equal checks for equal center and size values using tolerant numeric comparison.
func (r Rectangle[T]) Equal(rectangle Rectangle[T]) bool {
	return r.Center.Equal(rectangle.Center) && r.Size.Equal(rectangle.Size)
}

// IsZero checks if center point and size are zero.
func (r Rectangle[T]) IsZero() bool {
	return r.Equal(Rectangle[T]{})
}

// Polygon converts the rectangle into a generic Polygon with computed vertices.
func (r Rectangle[T]) Polygon() Polygon[T] {
	return Polygon[T]{r.Vertices()}
}

// Int converts the rectangle to a Rectangle[int].
func (r Rectangle[T]) Int() Rectangle[int] {
	return Rectangle[int]{r.Center.Int(), r.Size.Int()}
}

// Float converts the rectangle to a Rectangle[float64].
func (r Rectangle[T]) Float() Rectangle[float64] {
	return Rectangle[float64]{r.Center.Float(), r.Size.Float()}
}

// String returns the rectangle in the form of its constructor: Rect((x,y);WxH), center then
// size. See MinMaxString for the corners.
func (r Rectangle[T]) String() string {
	return fmt.Sprintf("Rect(%s;%s)", r.Center.String(), r.Size.String())
}

// MinMaxString returns the rectangle by its corners: (x,y)-(x,y), Min then Max. It is the form
// to read positions from, where String mirrors how the rectangle was built.
func (r Rectangle[T]) MinMaxString() string {
	return fmt.Sprintf("%s-%s", r.Min().String(), r.Max().String())
}

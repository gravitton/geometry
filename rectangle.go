package geom

import (
	"fmt"
	"iter"
	"slices"
)

// Rectangle is a 2D axis-aligned rectangle represented by its center and size.
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

// Rect is shorthand for Rectangle{center, size}.
func Rect[T Number](center Point[T], size Size[T]) Rectangle[T] {
	return Rectangle[T]{center, size}
}

// RectangleFromMin creates a Rectangle from min point and size.
func RectangleFromMin[T Number](min Point[T], size Size[T]) Rectangle[T] {
	w, h := size.XY()

	return Rectangle[T]{min.AddXY(w/2, h/2), size}
}

// RectangleFromMax creates a Rectangle from max point and size.
func RectangleFromMax[T Number](max Point[T], size Size[T]) Rectangle[T] {
	w, h := size.XY()

	return Rectangle[T]{max.AddXY(-w+w/2, -h+h/2), size}
}

// RectangleFromMinMax creates a Rectangle from min and max points.
func RectangleFromMinMax[T Number](min, max Point[T]) Rectangle[T] {
	return RectangleFromMin(min, Sz(max.Subtract(min).XY()))
}

// RectangleFromSize creates a Rectangle from zero point and size.
func RectangleFromSize[T Number](size Size[T]) Rectangle[T] {
	return RectangleFromMin(Pt[T](0, 0), size)
}

// Translate creates a new Rectangle translated by the given vector.
func (r Rectangle[T]) Translate(vector Vector[T]) Rectangle[T] {
	return Rectangle[T]{r.Center.Add(vector), r.Size}
}

// MoveTo creates a new Rectangle with the same size centered at point.
func (r Rectangle[T]) MoveTo(point Point[T]) Rectangle[T] {
	return Rectangle[T]{point, r.Size}
}

// Scale creates a new Rectangle with size uniformly scaled by the factor.
func (r Rectangle[T]) Scale(factor float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Scale(factor)}
}

// ScaleXY creates a new Rectangle with size scaled by the given factors.
func (r Rectangle[T]) ScaleXY(factorX, factorY float64) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ScaleXY(factorX, factorY)}
}

// Resize creates a new Rectangle with the given size.
func (r Rectangle[T]) Resize(size Size[T]) Rectangle[T] {
	return Rectangle[T]{r.Center, size}
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

// AlignTo creates a new Rectangle moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (r Rectangle[T]) AlignTo(direction Direction, point Point[T]) Rectangle[T] {
	return r.Translate(point.Subtract(r.Anchor(direction)))
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
		for _, edge := range [4]Line[T]{r.TopEdge(), r.RightEdge(), r.BottomEdge(), r.LeftEdge()} {
			if !yield(edge) {
				return
			}
		}
	}
}

// Vertices returns the rectangle vertices in order starting Min point, by increasing angle —
// the same winding as Directions and RegularPolygon.Vertices, and clockwise as drawn on a screen
// with Y pointing down.
func (r Rectangle[T]) Vertices() []Point[T] {
	return []Point[T]{
		r.TopLeft(),
		r.TopRight(),
		r.BottomRight(),
		r.BottomLeft(),
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

// Bounds returns the axis-aligned bounding rectangle.
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}

// Clamp returns the given Point clamped to the rectangle bounds.
func (r Rectangle[T]) Clamp(point Point[T]) Point[T] {
	a, b := r.MinMax()

	return Point[T]{Clamp(point.X, a.X, b.X), Clamp(point.Y, a.Y, b.Y)}
}

// Equal checks for equal center and size values using tolerant numeric comparison.
func (r Rectangle[T]) Equal(rectangle Rectangle[T]) bool {
	return r.Center.Equal(rectangle.Center) && r.Size.Equal(rectangle.Size)
}

// IsZero checks if center point and size are zero.
func (r Rectangle[T]) IsZero() bool {
	return r.Center.IsZero() && r.Size.IsZero()
}

// Contains reports whether the given point lies within the rectangle, boundary included within
// Epsilon of T.
func (r Rectangle[T]) Contains(point Point[T]) bool {
	return point.Between(r.MinMax())
}

// DistanceTo returns the distance from the given point to the nearest point of the rectangle:
// zero for a point within it, the same closed convention as Contains.
func (r Rectangle[T]) DistanceTo(point Point[T]) float64 {
	return r.Clamp(point).DistanceTo(point)
}

// DistanceSquaredTo returns the squared distance from the given point to the nearest point of
// the rectangle, faster for comparisons. It stays in T like Point.DistanceSquaredTo, since the
// nearest point is the clamped point and a lattice point for an integer T.
func (r Rectangle[T]) DistanceSquaredTo(point Point[T]) T {
	return r.Clamp(point).DistanceSquaredTo(point)
}

// Intersects reports whether the rectangles overlap. Touching rectangles intersect, within
// Epsilon of T, the same closed convention as Contains.
func (r Rectangle[T]) Intersects(rectangle Rectangle[T]) bool {
	a1, b1 := r.MinMax()
	a2, b2 := rectangle.MinMax()

	return LessOrEqual(a1.X, b2.X) && LessOrEqual(a2.X, b1.X) &&
		LessOrEqual(a1.Y, b2.Y) && LessOrEqual(a2.Y, b1.Y)
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
// within Epsilon of T, and the rectangle bounds are the same Min and Max that Contains uses.
func (r Rectangle[T]) IntersectsCircle(circle Circle[T]) bool {
	closest := r.Clamp(circle.Center)

	return closest.Subtract(circle.Center).LessOrEqual(circle.Radius)
}

// IntersectsLine reports whether the rectangle and the segment share a point, as
// Line.IntersectsRectangle does.
func (r Rectangle[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsRectangle(r)
}

// IntersectsPolygon reports whether the rectangle and the polygon share a point, as
// Polygon.IntersectsRectangle does.
func (r Rectangle[T]) IntersectsPolygon(polygon Polygon[T]) bool {
	return polygon.IntersectsRectangle(r)
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

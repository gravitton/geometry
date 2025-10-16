package geom

import (
	"fmt"
)

// Rectangle is a 2D axis-aligned rectangle represented by its center and size.
type Rectangle[T Number] struct {
	Center Point[T] `json:",inline"`
	Size   Size[T]  `json:",inline"`
}

// Rect is shorthand for Rectangle{center, size}.
func Rect[T Number](center Point[T], size Size[T]) Rectangle[T] {
	return Rectangle[T]{center, size}
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

// Grow creates a new Rectangle with size expanded by the same amount in both dimensions.
func (r Rectangle[T]) Grow(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Grow(amount)}
}

// GrowXY creates a new Rectangle with size expanded by the given amounts along X and Y.
func (r Rectangle[T]) GrowXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.GrowXY(amountX, amountY)}
}

// Shrink creates a new Rectangle with size reduced by the same amount in both dimensions.
func (r Rectangle[T]) Shrink(amount T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.Shrink(amount)}
}

// ShrinkXY creates a new Rectangle with size reduced by the given amounts along X and Y.
func (r Rectangle[T]) ShrinkXY(amountX, amountY T) Rectangle[T] {
	return Rectangle[T]{r.Center, r.Size.ShrinkXY(amountX, amountY)}
}

// Shrink creates a new Rectangle with size reduced by the same amount in both dimensions.
func (r Rectangle[T]) Inset(padding Padding[T]) Rectangle[T] {
	return Rectangle[T]{
		r.Center.AddXY(
			Divide(padding.Left-padding.Right, 2),
			Divide(padding.Bottom-padding.Top, 2),
		),
		r.Size.ShrinkXY(
			padding.Left+padding.Right,
			padding.Top+padding.Bottom,
		),
	}
}

// Width returns the rectangle width.
func (r Rectangle[T]) Width() T {
	return r.Size.Width
}

// Height returns the rectangle height.
func (r Rectangle[T]) Height() T {
	return r.Size.Height
}

// Min returns the minimum corner point of the rectangle.
func (r Rectangle[T]) Min() Point[T] {
	w, h := r.Size.XY()

	return r.Center.AddXY(-w/2, -h/2)
}

// Max returns the maximum corner point of the rectangle.
func (r Rectangle[T]) Max() Point[T] {
	w, h := r.Size.XY()

	return r.Center.AddXY(w-w/2, h-h/2)
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

// Edges returns the rectangle edges as lines in clockwise order starting from BottomLeft.
func (r Rectangle[T]) Edges() []Line[T] {
	tl := r.TopLeft()
	bl := r.BottomLeft()
	br := r.BottomRight()
	tr := r.TopRight()

	return []Line[T]{
		Ln(tl, bl),
		Ln(bl, br),
		Ln(br, tr),
		Ln(tr, tl),
	}
}

// Vertices returns the polygon vertices in order starting Min point, counter-clockwise.
func (r Rectangle[T]) Vertices() []Point[T] {
	return []Point[T]{
		r.TopLeft(),
		r.BottomLeft(),
		r.BottomRight(),
		r.TopRight(),
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

// Clamp creates a new point clamped to the rectangle.
func (r Rectangle[T]) Clamp(point Point[T]) Point[T] {
	minPoint, maxPoint := r.Min(), r.Max()

	return Point[T]{Clamp(point.X, minPoint.X, maxPoint.X), Clamp(point.Y, minPoint.Y, maxPoint.Y)}
}

// Equal checks for equal center and size values using tolerant numeric comparison.
func (r Rectangle[T]) Equal(other Rectangle[T]) bool {
	return r.Center.Equal(other.Center) && r.Size.Equal(other.Size)
}

// IsZero checks if center point and size are zero.
func (r Rectangle[T]) IsZero() bool {
	return r.Center.IsZero() && r.Size.IsZero()
}

// Contains reports whether the given point lies within or on the rectangle bounds.
func (r Rectangle[T]) Contains(point Point[T]) bool {
	minPoint, maxPoint := r.Min(), r.Max()

	return minPoint.X <= point.X && point.X <= maxPoint.X && minPoint.Y <= point.Y && point.Y <= maxPoint.Y
}

// Polygon converts the rectangle into a generic Polygon with computed vertices.
func (r Rectangle[T]) Polygon() Polygon[T] {
	return Polygon[T]{r.Vertices()}
}

// Int converts the rectangle to a [int] rectangle.
func (r Rectangle[T]) Int() Rectangle[int] {
	return Rectangle[int]{r.Center.Int(), r.Size.Int()}
}

// Float converts the rectangle to a [float64] rectangle.
func (r Rectangle[T]) Float() Rectangle[float64] {
	return Rectangle[float64]{r.Center.Float(), r.Size.Float()}
}

// String returns a string representation of the Rectangle using min and max.
func (r Rectangle[T]) String() string {
	return fmt.Sprintf("%s-%s", r.Min().String(), r.Max().String())
}

// RectFromMin creates a Rectangle from min point and size.
func RectFromMin[T Number](min Point[T], size Size[T]) Rectangle[T] {
	// must be same calculation (in reverse) as in Min() method
	w, h := size.XY()
	return Rectangle[T]{min.AddXY(w/2, h/2), size}
}

// RectFromMinMax creates a Rectangle from min and max points.
func RectFromMinMax[T Number](min, max Point[T]) Rectangle[T] {
	return RectFromMin(min, Sz(max.Subtract(min).XY()))
}

// RectFromSize creates a Rectangle from zero point and size.
func RectFromSize[T Number](size Size[T]) Rectangle[T] {
	return RectFromMin(Pt[T](0, 0), size)
}

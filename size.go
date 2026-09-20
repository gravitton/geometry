package geom

import (
	"fmt"
	"strings"
)

// Size is a 2D size: a signed pair of extents, not a shape. A negative extent is meaningful
// where a size measures a displacement, as RectangleFromMin reads it, and arithmetic such as
// Scale, Lerp and Vector.Size can produce one. A shape that stores a size as its own extent
// takes it absolute where it enters the shape, so a rectangle or regular polygon never holds a
// negative one; Abs is the same operation on the size alone. The size itself never clamps or
// takes an extent absolute, so a struct literal, Grow and Shrink carry a negative extent as it is.
type Size[T Number] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// Sz is shorthand for Size{width, height}.
func Sz[T Number](width, height T) Size[T] {
	return Size[T]{width, height}
}

// SzU is shorthand for Size{size, size}.
func SzU[T Number](size T) Size[T] {
	return Size[T]{size, size}
}

// ParseSize parses a size string in the form "WxH" (e.g. "16x16" or "23.0x12.1").
// For integer T, only integer strings parse; a fractional value is an error, not a rounded size.
func ParseSize[T Number](s string) (Size[T], error) {
	width, height, ok := strings.Cut(s, "x")
	if !ok {
		return Size[T]{}, fmt.Errorf("invalid size format: %s", s)
	}

	x, err := Parse[T](width)
	if err != nil {
		return Size[T]{}, fmt.Errorf("invalid width value: %w", err)
	}
	y, err := Parse[T](height)
	if err != nil {
		return Size[T]{}, fmt.Errorf("invalid height value: %w", err)
	}

	return Size[T]{x, y}, nil
}

// XY returns the size width, height values in standard order.
func (s Size[T]) XY() (T, T) {
	return s.Width, s.Height
}

// Area returns the size's area (width * height).
func (s Size[T]) Area() T {
	return Cast[T](float64(s.Width) * float64(s.Height))
}

// Perimeter returns the size's perimeter (2 * (width + height)).
func (s Size[T]) Perimeter() T {
	return Cast[T](2 * (float64(s.Width) + float64(s.Height)))
}

// AspectRatio returns (width / height). A zero Height has no usable ratio and returns 0 rather
// than an Inf or NaN that would poison later arithmetic; the same 0 a zero Width gives, since
// neither degenerate size has a ratio worth distinguishing.
func (s Size[T]) AspectRatio() float64 {
	return ratio(s.Width, s.Height)
}

// Scale creates a new Size scaled by the given factor in both dimensions.
func (s Size[T]) Scale(factor float64) Size[T] {
	return Size[T]{Multiply(s.Width, factor), Multiply(s.Height, factor)}
}

// ScaleXY creates a new Size scaled by the given factors along X and Y.
func (s Size[T]) ScaleXY(factorX, factorY float64) Size[T] {
	return Size[T]{Multiply(s.Width, factorX), Multiply(s.Height, factorY)}
}

// Unscale creates a new Size scaled with inverse factor in both dimensions.
func (s Size[T]) Unscale(factor float64) Size[T] {
	return Size[T]{Divide(s.Width, factor), Divide(s.Height, factor)}
}

// UnscaleXY creates a new Size scaled with inverse factors along X and Y.
func (s Size[T]) UnscaleXY(factorX, factorY float64) Size[T] {
	return Size[T]{Divide(s.Width, factorX), Divide(s.Height, factorY)}
}

// Abs creates a new Size with absolute width and height, the size a Rectangle stores.
func (s Size[T]) Abs() Size[T] {
	return Size[T]{Abs(s.Width), Abs(s.Height)}
}

// Transpose creates a new Size with the width and height swapped, the size of the same extent
// turned a quarter turn, as a rotated tile or sprite has. Axis.Size builds the same pair from
// an axis in hand.
func (s Size[T]) Transpose() Size[T] {
	return Size[T]{s.Height, s.Width}
}

// Round creates a new Size by rounding width and height to the nearest integer.
func (s Size[T]) Round() Size[T] {
	return Size[T]{Round(s.Width), Round(s.Height)}
}

// Floor creates a new Size by rounding width and height down to the nearest integer.
func (s Size[T]) Floor() Size[T] {
	return Size[T]{Floor(s.Width), Floor(s.Height)}
}

// Ceil creates a new Size by rounding width and height up to the nearest integer.
func (s Size[T]) Ceil() Size[T] {
	return Size[T]{Ceil(s.Width), Ceil(s.Height)}
}

// Lerp creates a new Size in linear interpolation towards the given size.
func (s Size[T]) Lerp(size Size[T], t float64) Size[T] {
	return Size[T]{Lerp(s.Width, size.Width, t), Lerp(s.Height, size.Height, t)}
}

// Grow creates a new Size expanded by the same amount in both dimensions. The amount is the
// total change of each extent, not an amount per side. A size is signed, so nothing is clamped:
// a shape that stores a size clamps its own extent at zero, as Rectangle.Grow does.
func (s Size[T]) Grow(amount T) Size[T] {
	return Size[T]{s.Width + amount, s.Height + amount}
}

// GrowXY creates a new Size expanded by the given amounts along X and Y.
// Each amount is the total change of that extent, like Grow.
func (s Size[T]) GrowXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width + amountX, s.Height + amountY}
}

// Shrink creates a new Size reduced by the same amount in both dimensions, the inverse of Grow.
// The amount is the total change of each extent, not an amount per side, and a reduction past
// zero gives a negative size rather than clamping, like Grow.
func (s Size[T]) Shrink(amount T) Size[T] {
	return Size[T]{s.Width - amount, s.Height - amount}
}

// ShrinkXY creates a new Size reduced by the given amounts along X and Y.
// Each amount is the total change of that extent, like Shrink.
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width - amountX, s.Height - amountY}
}

// Fit creates a new Size scaled uniformly to the largest that fits within the given size, keeping
// the aspect ratio: one extent matches the given size and the other is at most it. A size with a
// zero width or height has no ratio to keep and fits as the zero size. A negative extent
// measures the other way and fits to a negative one, which no longer bounds anything; take Abs
// first where a size may carry one.
func (s Size[T]) Fit(size Size[T]) Size[T] {
	return s.Scale(min(s.ratios(size)))
}

// Fill creates a new Size scaled uniformly to the smallest that covers the given size, keeping
// the aspect ratio: one extent matches the given size and the other is at least it. A size with a
// zero width or height has no ratio to keep and fills as the zero size. A negative extent
// measures the other way and fills to a negative one, which no longer covers anything; take Abs
// first where a size may carry one.
func (s Size[T]) Fill(size Size[T]) Size[T] {
	return s.Scale(max(s.ratios(size)))
}

// AtLeast creates a new Size with at least the given width and height values.
func (s Size[T]) AtLeast(size Size[T]) Size[T] {
	return Size[T]{max(s.Width, size.Width), max(s.Height, size.Height)}
}

// AtMost creates a new Size with at most the given width and height values.
func (s Size[T]) AtMost(size Size[T]) Size[T] {
	return Size[T]{min(s.Width, size.Width), min(s.Height, size.Height)}
}

// ratios returns the factors that scale the width and the height onto the given size, both
// zero when either extent of s is zero, so that Fit and Fill agree on the zero size there.
func (s Size[T]) ratios(size Size[T]) (float64, float64) {
	if s.Width == 0 || s.Height == 0 {
		return 0, 0
	}

	return ratio(size.Width, s.Width), ratio(size.Height, s.Height)
}

// Equal checks for equal width and height values with given size.
func (s Size[T]) Equal(size Size[T]) bool {
	return Equal(s.Width, size.Width) && Equal(s.Height, size.Height)
}

// IsZero checks if width and height values are zero.
func (s Size[T]) IsZero() bool {
	return s.Equal(Size[T]{})
}

// Vector converts the size to a Vector.
func (s Size[T]) Vector() Vector[T] {
	return Vector[T]{s.Width, s.Height}
}

// Cast converts the size to a Size of another number type, rounding as Cast does.
func (s Size[T]) Cast[R Number]() Size[R] {
	return Size[R]{Cast[R](float64(s.Width)), Cast[R](float64(s.Height))}
}

// Int converts the size to a Size[int].
func (s Size[T]) Int() Size[int] {
	return Size[int]{Int(s.Width), Int(s.Height)}
}

// Float converts the size to a Size[float64].
func (s Size[T]) Float() Size[float64] {
	return Size[float64]{float64(s.Width), float64(s.Height)}
}

// String returns a string in the form "WxH" using the underlying number formatting.
func (s Size[T]) String() string {
	return fmt.Sprintf("%sx%s", String(s.Width), String(s.Height))
}

package geom

import (
	"fmt"
	"strings"
)

// Size is a 2D size. Width and Height are expected to be non-negative: the constructors do not
// check it, but every Rectangle constructor takes the size absolute or reorders its corners, so
// a rectangle with Min beyond Max can only be written as a struct literal or decoded from JSON.
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

// Grow creates a new Size expanded by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, not an amount per side.
func (s Size[T]) Grow(amount T) Size[T] {
	return Size[T]{max(s.Width+amount, 0), max(s.Height+amount, 0)}
}

// GrowXY creates a new Size expanded by the given amounts along X and Y, clamped to zero.
// Each amount is the total change of that extent, like Grow.
func (s Size[T]) GrowXY(amountX, amountY T) Size[T] {
	return Size[T]{max(s.Width+amountX, 0), max(s.Height+amountY, 0)}
}

// Shrink creates a new Size reduced by the same amount in both dimensions, clamped to zero.
// The amount is the total change of each extent, not an amount per side.
func (s Size[T]) Shrink(amount T) Size[T] {
	return Size[T]{max(s.Width-amount, 0), max(s.Height-amount, 0)}
}

// ShrinkXY creates a new Size reduced by the given amounts along X and Y, clamped to zero.
// Each amount is the total change of that extent, like Shrink.
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T] {
	return Size[T]{max(s.Width-amountX, 0), max(s.Height-amountY, 0)}
}

// Fit creates a new Size scaled uniformly to the largest that fits within the given size, keeping
// the aspect ratio: one extent matches the given size and the other is at most it. A size with a
// zero width or height has no ratio to keep and fits as the zero size.
func (s Size[T]) Fit(size Size[T]) Size[T] {
	return s.Scale(min(s.ratios(size)))
}

// Fill creates a new Size scaled uniformly to the smallest that covers the given size, keeping
// the aspect ratio: one extent matches the given size and the other is at least it. A size with a
// zero width or height has no ratio to keep and fills as the zero size.
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

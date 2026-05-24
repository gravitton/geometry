package geom

import (
	"fmt"
	"strings"
)

// Size is a 2D size.
type Size[T Number] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// Sz is shorthand for Size{width, height}.
func Sz[T Number](width, height T) Size[T] {
	return Size[T]{width, height}
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

// Grow creates a new Size expanded by the same delta in both dimensions.
func (s Size[T]) Grow(amount T) Size[T] {
	return Size[T]{s.Width + amount, s.Height + amount}
}

// GrowXY creates a new Size expanded by the given amounts along X and Y.
func (s Size[T]) GrowXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width + amountX, s.Height + amountY}
}

// Shrink creates a new Size reduced by the same delta in both dimensions, clamped to zero.
func (s Size[T]) Shrink(amount T) Size[T] {
	return Size[T]{max(s.Width-amount, 0), max(s.Height-amount, 0)}
}

// ShrinkXY creates a new Size reduced by the given amounts along X and Y, clamped to zero.
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T] {
	return Size[T]{max(s.Width-amountX, 0), max(s.Height-amountY, 0)}
}

// Area returns the size's area (width * height).
func (s Size[T]) Area() T {
	return s.Width * s.Height
}

// Perimeter returns the size's perimeter (2 * (width + height)).
func (s Size[T]) Perimeter() T {
	return 2 * (s.Width + s.Height)
}

// AspectRatio returns (width / height). Returns 0 when Height is zero.
func (s Size[T]) AspectRatio() float64 {
	if s.Height == 0 {
		return 0
	}
	return float64(s.Width) / float64(s.Height)
}

// AtLeast creates a new Size with at least the given width and height values.
func (s Size[T]) AtLeast(size Size[T]) Size[T] {
	return Size[T]{max(s.Width, size.Width), max(s.Height, size.Height)}
}

// AtMost creates a new Size with at most the given width and height values.
func (s Size[T]) AtMost(size Size[T]) Size[T] {
	return Size[T]{min(s.Width, size.Width), min(s.Height, size.Height)}
}

// Equal checks for equal width and height values with given size.
func (s Size[T]) Equal(size Size[T]) bool {
	return Equal(s.Width, size.Width) && Equal(s.Height, size.Height)
}

// IsZero checks if width and height values are zero.
func (s Size[T]) IsZero() bool {
	return s.Equal(Size[T]{})
}

// XY returns the size width, height values in standard order.
func (s Size[T]) XY() (T, T) {
	return s.Width, s.Height
}

// Vector converts the size to a Vector.
func (s Size[T]) Vector() Vector[T] {
	return Vector[T]{s.Width, s.Height}
}

// Int converts the size to a [int] size.
func (s Size[T]) Int() Size[int] {
	return Size[int]{Cast[int](float64(s.Width)), Cast[int](float64(s.Height))}
}

// Float converts the size to a [float64] size.
func (s Size[T]) Float() Size[float64] {
	return Size[float64]{float64(s.Width), float64(s.Height)}
}

// String returns a string in the form "WxH" using the underlying number formatting.
func (s Size[T]) String() string {
	return fmt.Sprintf("%sx%s", String(s.Width), String(s.Height))
}

// SzU is shorthand for Size{size, size}.
func SzU[T Number](size T) Size[T] {
	return Size[T]{size, size}
}

// ParseSize parses a size string in the form "WxH" (e.g. "16x16" or "23.0x12.1").
// For integer T, float values are rounded to the nearest integer via Cast.
func ParseSize[T Number](s string) (Size[T], error) {
	parts := strings.SplitN(s, "x", 2)
	if len(parts) != 2 {
		return Size[T]{}, fmt.Errorf("invalid size format: %s", s)
	}

	x, err := Parse[T](parts[0])
	if err != nil {
		return Size[T]{}, fmt.Errorf("invalid width value: %s", parts[0])
	}
	y, err := Parse[T](parts[1])
	if err != nil {
		return Size[T]{}, fmt.Errorf("invalid height value: %s", parts[1])
	}

	return Size[T]{x, y}, nil
}

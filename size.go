package geom

import (
	"fmt"
)

// Size is a 2D size.
type Size[T Number] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// S is shorthand for Size{width, height}.
func S[T Number](width, height T) Size[T] {
	return Size[T]{width, height}
}

// Scale creates a new Size scaled by the given factor in both dimensions.
func (s Size[T]) Scale(factor float64) Size[T] {
	return Size[T]{Multiple(s.Width, factor), Multiple(s.Height, factor)}
}

// ScaleXY creates a new Size scaled by the given factors along X and Y.
func (s Size[T]) ScaleXY(factorX, factorY float64) Size[T] {
	return Size[T]{Multiple(s.Width, factorX), Multiple(s.Height, factorY)}
}

// Grow creates a new Size expanded by the same delta in both dimensions.
func (s Size[T]) Grow(amount T) Size[T] {
	return Size[T]{s.Width + amount, s.Height + amount}
}

// GrowXY creates a new Size expanded by the given amounts along X and Y.
func (s Size[T]) GrowXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width + amountX, s.Height + amountY}
}

// Shrink creates a new Size reduced by the same delta in both dimensions.
func (s Size[T]) Shrink(amount T) Size[T] {
	return Size[T]{s.Width - amount, s.Height - amount}
}

// ShrinkXY creates a new Size reduced by the given amounts along X and Y.
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width - amountX, s.Height - amountY}
}

// Area returns the size's area (width * height).
func (s Size[T]) Area() T {
	return s.Width * s.Height
}

// Perimeter returns the size's perimeter (2 * (width + height)).
func (s Size[T]) Perimeter() T {
	return 2 * (s.Width + s.Height)
}

// AspectRatio returns (width / height).
func (s Size[T]) AspectRatio() float64 {
	return float64(s.Width) / float64(s.Height)
}

// Equal checks for equal width and height values with given size.
func (s Size[T]) Equal(other Size[T]) bool {
	return Equal(s.Width, other.Width) && Equal(s.Height, other.Height)
}

// IsZero checks if width and height values are zero.
func (s Size[T]) IsZero() bool {
	return s.Equal(Size[T]{})
}

// XY returns the size width, height values in standard order.
func (s Size[T]) XY() (T, T) {
	return s.Width, s.Height
}

// ToVector converts the size to a Vector.
func (s Size[T]) ToVector() Vector[T] {
	return Vector[T]{s.Width, s.Height}
}

// String returns a string in the form "WxH" using the underlying number formatting.
func (s Size[T]) String() string {
	return fmt.Sprintf("%sx%s", ToString(s.Width), ToString(s.Height))
}

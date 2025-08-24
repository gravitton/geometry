package geom

import "fmt"

// Size is a 2D size.
type Size[T Number] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// S is shorthand for Size[T]{width, height}.
func S[T Number](width, height T) Size[T] {
	return Size[T]{width, height}
}

// Scale creates a new Size scaled by the given factor in both dimensions.
func (s Size[T]) Scale(scale float64) Size[T] {
	return Size[T]{Multiple(s.Width, scale), Multiple(s.Height, scale)}
}

// ScaleXY creates a new Size scaled by the given factors along X and Y.
func (s Size[T]) ScaleXY(scaleX, scaleY float64) Size[T] {
	return Size[T]{Multiple(s.Width, scaleX), Multiple(s.Height, scaleY)}
}

// Expand creates a new Size expanded by the same delta in both dimensions.
func (s Size[T]) Expand(delta T) Size[T] {
	return Size[T]{s.Width + delta, s.Height + delta}
}

// ExpandXY creates a new Size expanded by the given deltas along X and Y.
func (s Size[T]) ExpandXY(deltaWidth, deltaHeight T) Size[T] {
	return Size[T]{s.Width + deltaWidth, s.Height + deltaHeight}
}

// Shrunk creates a new Size reduced by the same delta in both dimensions.
func (s Size[T]) Shrunk(delta T) Size[T] {
	return Size[T]{s.Width - delta, s.Height - delta}
}

// ShrunkXY creates a new Size reduced by the given deltas along X and Y.
func (s Size[T]) ShrunkXY(deltaWidth, deltaHeight T) Size[T] {
	return Size[T]{s.Width - deltaWidth, s.Height - deltaHeight}
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

// Zero checks if width and height values are 0.
func (s Size[T]) Zero() bool {
	return s.Equal(Size[T]{})
}

// XY returns the size width, height values in standard order.
func (s Size[T]) XY() (T, T) {
	return s.Width, s.Height
}

// String returns a string in the form "WxH" using the underlying number formatting.
func (s Size[T]) String() string {
	return fmt.Sprintf("%sx%s", ToString(s.Width), ToString(s.Height))
}

package geom

import (
	"fmt"
	"strings"
)

// Size is a pair of extents, not a shape.
//
// A size is signed and never clamps or takes an extent absolute: a negative extent measures a
// displacement the other way. A shape that stores a size takes it absolute, and AtLeastZero and Abs
// are the two repairs.
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

// ZeroSize returns the size of zero extents.
func ZeroSize[T Number]() Size[T] {
	return Size[T]{}
}

// ParseSize returns the size in the form String prints, split at its first x, and an error for any
// other string. Each side is a number Parse accepts, so a hexadecimal float can stand for the
// height but not the width; for integer T a fractional side is an error, not rounded.
func ParseSize[T Number](s string) (Size[T], error) {
	width, height, ok := strings.Cut(s, "x")
	if !ok {
		return Size[T]{}, fmt.Errorf("geom: invalid size format %q", s)
	}

	w, err := Parse[T](width)
	if err != nil {
		return Size[T]{}, fmt.Errorf("geom: invalid width value: %w", err)
	}
	h, err := Parse[T](height)
	if err != nil {
		return Size[T]{}, fmt.Errorf("geom: invalid height value: %w", err)
	}

	return Size[T]{w, h}, nil
}

// XY returns the width and the height.
func (s Size[T]) XY() (T, T) {
	return s.Width, s.Height
}

// Area returns the area (width·height).
func (s Size[T]) Area() T {
	return Cast[T](float64(s.Width) * float64(s.Height))
}

// Perimeter returns the perimeter (2·(width+height)).
func (s Size[T]) Perimeter() T {
	return Cast[T](2 * (float64(s.Width) + float64(s.Height)))
}

// AspectRatio returns the width divided by the height, and zero where either is zero.
func (s Size[T]) AspectRatio() float64 {
	return ratio(s.Width, s.Height)
}

// Scale returns the size with both extents multiplied by the factor.
func (s Size[T]) Scale(factor float64) Size[T] {
	return Size[T]{Multiply(s.Width, factor), Multiply(s.Height, factor)}
}

// ScaleXY returns the size with each extent multiplied by its own factor.
func (s Size[T]) ScaleXY(factorX, factorY float64) Size[T] {
	return Size[T]{Multiply(s.Width, factorX), Multiply(s.Height, factorY)}
}

// Unscale returns the size with both extents divided by the factor. It panics for a zero factor.
func (s Size[T]) Unscale(factor float64) Size[T] {
	return Size[T]{Divide(s.Width, factor), Divide(s.Height, factor)}
}

// UnscaleXY returns the size with each extent divided by its own factor. It panics for a zero
// factor.
func (s Size[T]) UnscaleXY(factorX, factorY float64) Size[T] {
	return Size[T]{Divide(s.Width, factorX), Divide(s.Height, factorY)}
}

// Grow returns the size with both extents increased by the amount, the total change of each extent
// rather than one per side.
func (s Size[T]) Grow(amount T) Size[T] {
	return Size[T]{s.Width + amount, s.Height + amount}
}

// GrowXY returns the size with each extent increased by its own amount.
func (s Size[T]) GrowXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width + amountX, s.Height + amountY}
}

// Shrink returns the size with both extents decreased by the amount, the total change of each
// extent rather than one per side, and negative past zero.
func (s Size[T]) Shrink(amount T) Size[T] {
	return Size[T]{s.Width - amount, s.Height - amount}
}

// ShrinkXY returns the size with each extent decreased by its own amount.
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T] {
	return Size[T]{s.Width - amountX, s.Height - amountY}
}

// Lerp returns the size interpolated toward the given one, extrapolated outside [0, 1].
func (s Size[T]) Lerp(size Size[T], t float64) Size[T] {
	return Size[T]{Lerp(s.Width, size.Width, t), Lerp(s.Height, size.Height, t)}
}

// Abs returns the size with absolute extents.
func (s Size[T]) Abs() Size[T] {
	return Size[T]{Abs(s.Width), Abs(s.Height)}
}

// Round returns the size with extents rounded to the nearest integer.
func (s Size[T]) Round() Size[T] {
	return Size[T]{Round(s.Width), Round(s.Height)}
}

// Floor returns the size with extents rounded down.
func (s Size[T]) Floor() Size[T] {
	return Size[T]{Floor(s.Width), Floor(s.Height)}
}

// Ceil returns the size with extents rounded up.
func (s Size[T]) Ceil() Size[T] {
	return Size[T]{Ceil(s.Width), Ceil(s.Height)}
}

// Transpose returns the size with the width and the height swapped, the extent turned a quarter
// turn.
func (s Size[T]) Transpose() Size[T] {
	return Size[T]{s.Height, s.Width}
}

// Fit returns the largest size of the same aspect ratio that fits in the given one. A size with a
// zero extent gives the zero size, and a negative extent gives a negative one: take Abs first.
func (s Size[T]) Fit(size Size[T]) Size[T] {
	return s.Scale(min(s.ratios(size)))
}

// Fill returns the smallest size of the same aspect ratio that covers the given one. A size with a
// zero extent gives the zero size, and a negative extent gives a negative one: take Abs first.
func (s Size[T]) Fill(size Size[T]) Size[T] {
	return s.Scale(max(s.ratios(size)))
}

// AtLeast returns the size with each extent raised to at least the one of the given size.
func (s Size[T]) AtLeast(size Size[T]) Size[T] {
	return Size[T]{max(s.Width, size.Width), max(s.Height, size.Height)}
}

// AtMost returns the size with each extent lowered to at most the one of the given size.
func (s Size[T]) AtMost(size Size[T]) Size[T] {
	return Size[T]{min(s.Width, size.Width), min(s.Height, size.Height)}
}

// AtLeastZero returns the size with each negative extent raised to zero, where Abs mirrors it.
func (s Size[T]) AtLeastZero() Size[T] {
	return Size[T]{max(s.Width, 0), max(s.Height, 0)}
}

// ratios returns the factors that scale the width and the height onto the given size, both zero
// where either extent is zero, so Fit and Fill agree on the zero size.
func (s Size[T]) ratios(size Size[T]) (float64, float64) {
	if s.Width == 0 || s.Height == 0 {
		return 0, 0
	}

	return ratio(size.Width, s.Width), ratio(size.Height, s.Height)
}

// Equal reports whether the sizes have equal width and height.
func (s Size[T]) Equal(size Size[T]) bool {
	return Equal(s.Width, size.Width) && Equal(s.Height, size.Height)
}

// IsZero reports whether both extents are zero.
func (s Size[T]) IsZero() bool {
	return s.Equal(Size[T]{})
}

// IsPositive reports whether both extents are greater than zero, exactly and with no tolerance. A
// NaN extent is not positive and an infinite one is.
func (s Size[T]) IsPositive() bool {
	return s.Width > 0 && s.Height > 0
}

// IsFinite reports whether both extents are neither NaN nor infinite, which Cast and Int need for
// an integer type. An integer size always is.
func (s Size[T]) IsFinite() bool {
	return IsFinite(s.Width) && IsFinite(s.Height)
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

// String returns the size in the form WxH.
func (s Size[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = s.appendString(buf)

	return string(buf)
}

// appendString appends the size as String prints it to buf and returns the extended buffer.
func (s Size[T]) appendString(buf []byte) []byte {
	buf = appendString(buf, s.Width)
	buf = append(buf, 'x')
	buf = appendString(buf, s.Height)

	return buf
}

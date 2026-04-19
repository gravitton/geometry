package geom

import "fmt"

// Padding represents a 2D padding.
type Padding[T Number] struct {
	Top    T `json:"t"`
	Right  T `json:"r"`
	Bottom T `json:"b"`
	Left   T `json:"l"`
}

// Pad is shorthand for Padding{top, right, bottom, left}.
func Pad[T Number](top, right, bottom, left T) Padding[T] {
	return Padding[T]{top, right, bottom, left}
}

// Width returns the width of the padding.
func (p Padding[T]) Width() T {
	return p.Left + p.Right
}

// Height returns the height of the padding.
func (p Padding[T]) Height() T {
	return p.Top + p.Bottom
}

// XY returns the width and height of the padding.
func (p Padding[T]) XY() (T, T) {
	return p.Width(), p.Height()
}

// Size converts the padding to a Size.
func (p Padding[T]) Size() Size[T] {
	return Size[T]{p.Width(), p.Height()}
}

// Int converts the padding to a [int] padding.
func (p Padding[T]) Int() Padding[int] {
	return Padding[int]{Cast[int](float64(p.Top)), Cast[int](float64(p.Right)), Cast[int](float64(p.Bottom)), Cast[int](float64(p.Left))}
}

// Float converts the padding to a [float64] padding.
func (p Padding[T]) Float() Padding[float64] {
	return Padding[float64]{float64(p.Top), float64(p.Right), float64(p.Bottom), float64(p.Left)}
}

// String returns a string representation of the Padding.
func (p Padding[T]) String() string {
	return fmt.Sprintf("Pad(%s;%s;%s;%s)", String(p.Top), String(p.Right), String(p.Bottom), String(p.Left))
}

// PadU is shorthand for Padding{padding, padding, padding, padding}.
func PadU[T Number](padding T) Padding[T] {
	return Padding[T]{padding, padding, padding, padding}
}

// PadXY is shorthand for Padding{topBottom, leftRight, topBottom, leftRight}.
func PadXY[T Number](topBottom, leftRight T) Padding[T] {
	return Padding[T]{topBottom, leftRight, topBottom, leftRight}
}

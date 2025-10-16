package geom

// Padding represents a 2D padding.
type Padding[T Number] struct {
	Top, Right, Bottom, Left T
}

// Pad is shorthand for Padding{top, right, bottom, left}.
func Pad[T Number](top, right, bottom, left T) Padding[T] {
	return Padding[T]{top, right, bottom, left}
}

// PadU is shorthand for Padding{padding, padding, padding, padding}.
func PadU[T Number](padding T) Padding[T] {
	return Padding[T]{padding, padding, padding, padding}
}

// PadXY is shorthand for Padding{topBottom, leftRight, topBottom, leftRight}.
func PadXY[T Number](topBottom, leftRight T) Padding[T] {
	return Padding[T]{topBottom, leftRight, topBottom, leftRight}
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

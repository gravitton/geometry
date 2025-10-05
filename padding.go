package geom

// Padding represents a 2D padding.
type Padding[T Number] struct {
	Top, Right, Bottom, Left T
}

// Pd is shorthand for Padding{top, right, bottom, left}.
func Pd[T Number](top, right, bottom, left T) Padding[T] {
	return Padding[T]{top, right, bottom, left}
}

// PdU is shorthand for Padding{padding, padding, padding, padding}.
func PdU[T Number](padding T) Padding[T] {
	return Padding[T]{padding, padding, padding, padding}
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

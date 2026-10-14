package geom

// Padding is the space kept inside the four edges of a box, X first: Left, Right, Top, Bottom. A
// negative amount moves an edge out.
type Padding[T Number] struct {
	Left   T `json:"l"`
	Right  T `json:"r"`
	Top    T `json:"t"`
	Bottom T `json:"b"`
}

// Pad is shorthand for Padding{left, right, top, bottom}.
func Pad[T Number](left, right, top, bottom T) Padding[T] {
	return Padding[T]{left, right, top, bottom}
}

// PadXY is shorthand for Padding{leftRight, leftRight, topBottom, topBottom}.
func PadXY[T Number](leftRight, topBottom T) Padding[T] {
	return Padding[T]{leftRight, leftRight, topBottom, topBottom}
}

// PadU is shorthand for Padding{padding, padding, padding, padding}.
func PadU[T Number](padding T) Padding[T] {
	return Padding[T]{padding, padding, padding, padding}
}

// XY returns the Width and the Height.
func (p Padding[T]) XY() (T, T) {
	return p.Width(), p.Height()
}

// Width returns the sum of the left and right amounts.
func (p Padding[T]) Width() T {
	return p.Left + p.Right
}

// Height returns the sum of the top and bottom amounts.
func (p Padding[T]) Height() T {
	return p.Top + p.Bottom
}

// Add returns the sum of the paddings, edge by edge.
func (p Padding[T]) Add(padding Padding[T]) Padding[T] {
	return Padding[T]{p.Left + padding.Left, p.Right + padding.Right, p.Top + padding.Top, p.Bottom + padding.Bottom}
}

// Subtract returns the difference of the paddings, edge by edge.
func (p Padding[T]) Subtract(padding Padding[T]) Padding[T] {
	return Padding[T]{p.Left - padding.Left, p.Right - padding.Right, p.Top - padding.Top, p.Bottom - padding.Bottom}
}

// Negate returns the padding with every amount negated.
func (p Padding[T]) Negate() Padding[T] {
	return Padding[T]{-p.Left, -p.Right, -p.Top, -p.Bottom}
}

// Scale returns the padding with every amount multiplied by the factor.
func (p Padding[T]) Scale(factor float64) Padding[T] {
	return Padding[T]{Multiply(p.Left, factor), Multiply(p.Right, factor), Multiply(p.Top, factor), Multiply(p.Bottom, factor)}
}

// Unscale returns the padding with every amount divided by the factor. It panics for a zero factor.
func (p Padding[T]) Unscale(factor float64) Padding[T] {
	return Padding[T]{Divide(p.Left, factor), Divide(p.Right, factor), Divide(p.Top, factor), Divide(p.Bottom, factor)}
}

// Lerp returns the padding interpolated toward the given one, edge by edge, extrapolated outside
// [0, 1].
func (p Padding[T]) Lerp(padding Padding[T], t float64) Padding[T] {
	return Padding[T]{Lerp(p.Left, padding.Left, t), Lerp(p.Right, padding.Right, t), Lerp(p.Top, padding.Top, t), Lerp(p.Bottom, padding.Bottom, t)}
}

// Equal reports whether the paddings have equal amounts.
func (p Padding[T]) Equal(padding Padding[T]) bool {
	return Equal(p.Left, padding.Left) && Equal(p.Right, padding.Right) && Equal(p.Top, padding.Top) && Equal(p.Bottom, padding.Bottom)
}

// IsZero reports whether every amount is zero.
func (p Padding[T]) IsZero() bool {
	return p.Equal(Padding[T]{})
}

// Size returns the Width and the Height as a Size.
func (p Padding[T]) Size() Size[T] {
	return Size[T]{p.Width(), p.Height()}
}

// Cast converts the padding to a Padding of another number type, rounding as Cast does.
func (p Padding[T]) Cast[R Number]() Padding[R] {
	return Padding[R]{Cast[R](float64(p.Left)), Cast[R](float64(p.Right)), Cast[R](float64(p.Top)), Cast[R](float64(p.Bottom))}
}

// Int converts the padding to a Padding[int].
func (p Padding[T]) Int() Padding[int] {
	return Padding[int]{Int(p.Left), Int(p.Right), Int(p.Top), Int(p.Bottom)}
}

// Float converts the padding to a Padding[float64].
func (p Padding[T]) Float() Padding[float64] {
	return Padding[float64]{float64(p.Left), float64(p.Right), float64(p.Top), float64(p.Bottom)}
}

// String returns the padding in the form of its constructor: Pad(l;r;t;b).
func (p Padding[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Pad("...)
	buf = appendString(buf, p.Left)
	buf = append(buf, ';')
	buf = appendString(buf, p.Right)
	buf = append(buf, ';')
	buf = appendString(buf, p.Top)
	buf = append(buf, ';')
	buf = appendString(buf, p.Bottom)
	buf = append(buf, ')')

	return string(buf)
}

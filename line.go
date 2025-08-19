package geom

import "fmt"

// Line is a 2D line from point A to point B.
type Line[T Number] struct {
	A Point[T] `json:"a"`
	B Point[T] `json:"b"`
}

func L[T Number](x1, y1, x2, y2 T) Line[T] {
	return Line[T]{Point[T]{x1, y1}, Point[T]{x2, y2}}
}

func (l Line[T]) Translate(change Vector[T]) Line[T] {
	return Line[T]{l.A.Add(change), l.B.Add(change)}
}

func (l Line[T]) Reversed() Line[T] {
	return Line[T]{l.B, l.A}
}

func (l Line[T]) Midpoint() Point[T] {
	return l.A.Midpoint(l.B)
}

func (l Line[T]) Direction() Vector[T] {
	return l.B.Sub(l.A)
}

func (l Line[T]) Length() float64 {
	return l.Direction().Length()
}

func (l Line[T]) Equal(other Line[T]) bool {
	return l.A.Equal(other.A) && l.B.Equal(other.B)
}

func (l Line[T]) String() string {
	return fmt.Sprintf("L(%s;%s)", l.A.String(), l.B.String())
}

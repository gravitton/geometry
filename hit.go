package geom

// Hit is where a cast first reaches the boundary of a shape from outside, the result of the Hit
// methods of Segment and Ray, Sweep and Raycast: the Point reached, the unit Normal of the boundary
// there, facing the cast, and the Fraction of the cast at which the point lies, the argument
// PointAt takes, or of the motion for Sweep.
//
// For integer T the point is rounded while the normal is not, and the fraction is that of the
// rounded point, so PointAt of it can give a neighboring lattice point.
type Hit[T Number] struct {
	Point    Point[T]        `json:"p"`
	Normal   Vector[float64] `json:"n"`
	Fraction float64         `json:"f"`
}

// Equal reports whether the hits have equal point, normal and fraction.
func (h Hit[T]) Equal(hit Hit[T]) bool {
	return h.Point.Equal(hit.Point) && h.Normal.Equal(hit.Normal) && Equal(h.Fraction, hit.Fraction)
}

// IsZero reports whether the hit is the zero value, the one a cast that reaches nothing returns.
func (h Hit[T]) IsZero() bool {
	return h.Equal(Hit[T]{})
}

// String returns the hit as Hit((x,y);⟨x,y⟩;f).
func (h Hit[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Hit("...)
	buf = h.Point.appendString(buf)
	buf = append(buf, ';')
	buf = h.Normal.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, h.Fraction)
	buf = append(buf, ')')

	return string(buf)
}

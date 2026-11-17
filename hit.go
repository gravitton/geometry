package geom

// Hit is where a cast first reaches the boundary of a shape from outside, the result of the Hit
// methods of Segment and Ray, Sweep and Raycast.
//
// For integer T the point is rounded while the normal is not, and on a Segment or Ray the fraction
// is that of the rounded point, so PointAt of it can give a neighboring lattice point.
type Hit[T Number] struct {
	// Point is the point of the boundary reached.
	Point Point[T] `json:"p"`
	// Normal is the unit normal of the boundary at the point, facing the cast.
	Normal Vector[float64] `json:"n"`
	// Fraction is how far along the cast the point lies, the argument PointAt takes, or along the
	// motion for Sweep.
	Fraction float64 `json:"f"`
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

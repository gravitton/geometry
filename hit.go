package geom

// Hit is where a cast first reaches the boundary of a shape from outside, the answer of the
// Hit<Kind> methods of Segment and Ray: the Point reached, the unit Normal of the boundary there,
// turned to the side the cast comes from, and the Fraction of the cast at which the point lies,
// of the way from Start to End for a segment and of Direction for a ray, the argument PointAt
// takes. Sweep answers with one for a moving shape, the Fraction that of its motion, and Raycast
// for a ray against a shape of any kind.
//
// The normal and the fraction are float64 whatever T is. For integer T the point is rounded like
// every result while the normal is not, and the fraction is the one the rounded point projects
// to, so PointAt of it can give a neighbouring lattice point.
type Hit[T Number] struct {
	Point    Point[T]        `json:"p"`
	Normal   Vector[float64] `json:"n"`
	Fraction float64         `json:"f"`
}

// Equal checks for equal point, normal and fraction with given hit, the point exactly for an
// integer T and all three within the tolerance otherwise.
func (h Hit[T]) Equal(hit Hit[T]) bool {
	return h.Point.Equal(hit.Point) && h.Normal.Equal(hit.Normal) && Equal(h.Fraction, hit.Fraction)
}

// IsZero checks if the point, the normal and the fraction are zero, the hit a cast that
// reaches nothing returns beside its false.
func (h Hit[T]) IsZero() bool {
	return h.Equal(Hit[T]{})
}

// String returns the hit by its point, normal and fraction: Hit((x,y);⟨x,y⟩;f).
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

package geom

// Contact is how two overlapping shapes meet, the result of Collide. A Polygon that is not convex,
// or a Sector past half a turn, meets by its deepest edge or half alone, as Collide describes.
type Contact[T Number] struct {
	// Normal is the unit normal that parts the shapes, from the second toward the first, and none
	// between two single points.
	Normal Vector[float64] `json:"n"`
	// Depth is how far the shapes overlap along the normal, zero for shapes that only touch.
	Depth float64 `json:"d"`
	// Points are the points of the second shape where the shapes meet, two where flat sides lie
	// against each other and one otherwise.
	Points Pair[Point[T]] `json:"p"`
}

// Separation returns the shortest vector that moves the first shape out of the second until they
// only touch (Normal·Depth); its opposite moves the second out of the first. Against a Polygon that
// is not convex, or a Sector past half a turn, it parts them from one edge or half, so they may
// still overlap at another.
func (c Contact[T]) Separation() Vector[float64] {
	return c.Normal.Multiply(c.Depth)
}

// Equal reports whether the contacts have equal normal, depth and points, the points as Pair.Equal
// compares them.
func (c Contact[T]) Equal(contact Contact[T]) bool {
	return c.Normal.Equal(contact.Normal) && Equal(c.Depth, contact.Depth) && c.Points.Equal(contact.Points)
}

// IsZero reports whether the contact is the zero value, the one Collide returns beside false.
func (c Contact[T]) IsZero() bool {
	return c.Equal(Contact[T]{})
}

// String returns the contact as Contact(⟨x,y⟩;d) followed by each of its points, as
// Contact(⟨x,y⟩;d;(x,y);(x,y)).
func (c Contact[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Contact("...)
	buf = c.Normal.appendString(buf)
	buf = append(buf, ';')
	buf = appendString(buf, c.Depth)

	for i := range c.Points.count() {
		buf = append(buf, ';')
		buf = c.Points.Items[i].appendString(buf)
	}

	buf = append(buf, ')')

	return string(buf)
}

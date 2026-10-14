package geom

// Contact is how two overlapping shapes meet, the result of Collide: the unit Normal that parts
// them, from the second shape toward the first, the Depth they overlap by along it, and the Points
// of the second shape they meet at, two where flat sides lie against each other and one otherwise.
//
// Shapes that only touch have a depth of zero, and two single points have no normal. A Polygon that
// is not convex meets by its deepest edge alone, as Collide describes.
type Contact[T Number] struct {
	Normal Vector[float64] `json:"n"`
	Depth  float64         `json:"d"`
	Points Pair[Point[T]]  `json:"p"`
}

// Separation returns the shortest vector that moves the first shape out of the second until they
// only touch (Normal·Depth); its opposite moves the second out of the first. Against a Polygon that
// is not convex it parts them from one edge, so they may still overlap at another.
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

	for i := range c.Points.Count {
		buf = append(buf, ';')
		buf = c.Points.Items[i].appendString(buf)
	}

	buf = append(buf, ')')

	return string(buf)
}

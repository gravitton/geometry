package geom

// Contact is how two overlapping shapes meet, the answer of Collide, what a physics solver
// resolves a collision from: the unit Normal that parts them, from the second shape toward the
// first, the Depth they overlap by along it, and the Points of the second shape they meet at,
// two where flat sides lie against each other and one at a corner or a round side.
//
// The normal and the depth are float64 whatever T is, and the points are rounded like every
// result for an integer T. Two shapes that only touch have a depth of zero, and two that are
// both a single point have no normal. With a Polygon that is not convex the contact is that
// with one of its edges, the deepest, as Collide says.
type Contact[T Number] struct {
	Normal Vector[float64] `json:"n"`
	Depth  float64         `json:"d"`
	Points Pair[Point[T]]  `json:"p"`
}

// Separation returns the shortest vector that moves the first shape out of the second, so the
// two only touch: the Normal as far as the Depth. It is the zero vector for two shapes that
// only touch, and moving the second shape by its opposite parts them as well. A Polygon that
// is not convex is left by one edge at a time, so the two may still overlap at another edge
// once moved, and are parted by colliding again.
func (c Contact[T]) Separation() Vector[float64] {
	return c.Normal.Multiply(c.Depth)
}

package geom

import (
	"slices"
)

// edgeClip cuts an outline to the inner side of one line after another, Sutherland–Hodgman style.
// The vertices are read and never written; each cut appends its outline and moves it down to from,
// so dst needs room for two outlines. A region wound against Rectangle.Vertices is reversed.
type edgeClip[T Number] struct {
	dst      []Point[T]
	from     int
	points   []Point[T]
	reversed bool
}

// add cuts the outline to the side of the line through the edge that Segment.Normal points to on an
// edge of a Rectangle.
func (c *edgeClip[T]) add(edge Segment[T]) {
	if c.reversed {
		edge = edge.Reverse()
	}

	c.cut(edge.Line(), 1)
}

// cut keeps the part of the outline on the side the sense signs, the line included. A vertex on the
// line is kept and no crossing is placed beside it, and both senses read the same sides, so the two
// cuts of one line agree. A line with no direction cuts nothing, and a clip of no real cut leaves
// dst without the outline.
func (c *edgeClip[T]) cut(line Line[T], sense float64) {
	if !line.Vector.hasDirection() || len(c.points) == 0 {
		return
	}

	start := len(c.dst)
	previous := c.points[len(c.points)-1]
	previousSide := sense * c.side(line, previous)
	for _, point := range c.points {
		side := sense * c.side(line, point)
		if (side < 0 && previousSide > 0) || (side > 0 && previousSide < 0) {
			c.keep(start, c.crossing(previous, point, previousSide, side))
		}

		if side >= 0 {
			c.keep(start, point)
		}

		previous, previousSide = point, side
	}

	if n := len(c.dst); n-start > 1 && c.dst[n-1].Equal(c.dst[start]) {
		c.dst = c.dst[:n-1]
	}

	c.dst = append(c.dst[:c.from], c.dst[start:]...)
	c.points = c.dst[c.from:]
}

// keep appends the point to the outline begun at start unless it is Equal to the last, making room
// on the first for one vertex more than the outline had.
func (c *edgeClip[T]) keep(start int, point Point[T]) {
	n := len(c.dst)
	if n == start {
		c.dst = slices.Grow(c.dst, len(c.points)+1)
	}

	if n > start && c.dst[n-1].Equal(point) {
		return
	}

	c.dst = append(c.dst, point)
}

// side returns the cross product Line.Side reads, and zero for a point the line contains.
func (*edgeClip[T]) side(line Line[T], point Point[T]) float64 {
	if line.Contains(point) {
		return 0
	}

	return line.side(point)
}

// crossing returns where the edge from a to b crosses the line, given their sides of opposite
// signs, divided once so a crossing on a half unit is exact before rounding.
func (*edgeClip[T]) crossing(a, b Point[T], sideA, sideB float64) Point[T] {
	start := a.Float()
	along := b.Float().Subtract(start)
	denominator := sideA - sideB

	return Point[T]{Cast[T](start.X + float64(sideA*along.X)/denominator), Cast[T](start.Y + float64(sideA*along.Y)/denominator)}
}

// result returns the slice with the outline the cuts left appended to it.
func (c *edgeClip[T]) result() []Point[T] {
	return c.dst
}

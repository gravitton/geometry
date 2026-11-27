package geom

import (
	"slices"
	"sync"
)

// exactOutlines holds the float64 outlines an integer clip cuts, so a clip into a buffer with room
// allocates nothing once one is warm.
var exactOutlines = sync.Pool{New: func() any {
	return new([]Point[float64])
}}

// edgeClip cuts an outline to the inner side of one line after another, Sutherland–Hodgman style.
// The vertices are read and never written; each cut appends its outline and moves it down to from,
// so dst needs room for two outlines. A region wound against Rectangle.Vertices is reversed.
//
// For an integer T every cut is made in float64 on outline, kept in exact past a copy of the
// vertices with room for two outlines more, and result rounds what they leave once, since a cut
// reading the outline the last one rounded drifts a unit off at a sharp corner and parts the two
// orders of a pair. exact is taken from exactOutlines and given back by result.
type edgeClip[T Number] struct {
	dst      []Point[T]
	from     int
	points   []Point[T]
	reversed bool
	exact    []Point[float64]
	outline  []Point[float64]
	scratch  *[]Point[float64]
}

// edgeClipInto returns a clip of the points that appends to dst, with room for the given number of
// vertices more, and for an integer T the outline in float64 it cuts.
func edgeClipInto[T Number](dst, points []Point[T], reversed bool, room int) edgeClip[T] {
	clip := edgeClip[T]{dst: slices.Grow(dst, room), from: len(dst), points: points, reversed: reversed}
	if isInt[T]() {
		clip.scratch = exactOutlines.Get().(*[]Point[float64])
		clip.exact = slices.Grow((*clip.scratch)[:0], len(points)+room)
		for _, point := range points {
			clip.exact = append(clip.exact, point.Float())
		}

		clip.outline = clip.exact
	}

	return clip
}

// add cuts the outline to the side of the line through the edge that Segment.Normal points to on an
// edge of a Rectangle.
func (c *edgeClip[T]) add(edge Segment[T]) {
	if c.exact != nil {
		c.cutExact(edge.Float())

		return
	}

	if c.reversed {
		edge = edge.Reverse()
	}

	c.cut(edge.Line(), 1)
}

// shift cuts the outline as add does by the edge moved by the offset, which is rounded only where
// the cut is made in T.
func (c *edgeClip[T]) shift(edge Segment[T], offset Vector[float64]) {
	if c.exact != nil {
		c.cutExact(edge.Float().Translate(offset))

		return
	}

	c.add(edge.Translate(offset.Cast[T]()))
}

// cutExact cuts outline, the outline in float64, by the edge, as a clip appending to exact past the
// copy of the vertices.
func (c *edgeClip[T]) cutExact(edge Segment[float64]) {
	clip := edgeClip[float64]{dst: c.exact, from: len(c.points), points: c.outline, reversed: c.reversed}
	clip.add(edge)
	c.exact, c.outline = clip.dst, clip.points
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

// result returns the slice with the outline the cuts left appended to it, for an integer T the one
// left on exact, rounded, a vertex rounding onto the one before it given once.
func (c *edgeClip[T]) result() []Point[T] {
	if c.exact == nil {
		return c.dst
	}

	start := len(c.dst)
	for _, point := range c.exact[len(c.points):] {
		c.keep(start, point.Cast[T]())
	}

	if n := len(c.dst); n-start > 1 && c.dst[n-1].Equal(c.dst[start]) {
		c.dst = c.dst[:n-1]
	}

	*c.scratch = c.exact[:0]
	exactOutlines.Put(c.scratch)

	return c.dst
}

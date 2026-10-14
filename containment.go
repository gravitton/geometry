package geom

import "math"

const (
	// containHalvings bounds the arcs containsSupport halves. Shapes meeting at a few points need a
	// few dozen; only a shape close to the other all the way round needs more, and is taken not to
	// lie in it.
	containHalvings = 1024

	// containDepth bounds how often one arc is halved, far below the angle at which the corner
	// loses precision, and sizes the stack of arcs left to decide.
	containDepth = 32
)

// containsSupport reports whether the outer shape, taken as convex, holds the inner one, known by
// their supports alone: the inner one reaches no farther in any direction, within the tolerance.
// The four axes start four arcs, and an arc the inner shape may leave on is halved until it holds,
// a support of the inner one reaches past the outer one, or both are settled over it, which answers
// false, as do the bounds on halving and a NaN.
func containsSupport[T Number](outer, inner Shape[T]) bool {
	tolerance := minkowski[T]{a: inner, b: outer}.tolerance()

	var ends [4]arcEnd
	for i, direction := range [4]Vector[float64]{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
		end, ok := reachAlong(outer, inner, direction, tolerance)
		if !ok {
			return false
		}

		ends[i] = end
	}

	var stack [4 + containDepth]supportArc
	for i := range ends {
		stack[i] = supportArc{ends[i], ends[(i+1)%4], 0}
	}

	count, halvings := len(ends), 0
	for count > 0 {
		count--
		arc := stack[count]
		if !arc.exceeds(tolerance) {
			continue
		}

		if halvings == containHalvings || arc.depth == containDepth || arc.settled(tolerance) {
			return false
		}

		middle, ok := reachAlong(outer, inner, arc.start.direction.Add(arc.end.direction).Normalize(), tolerance)
		if !ok {
			return false
		}

		stack[count] = supportArc{arc.start, middle, arc.depth + 1}
		stack[count+1] = supportArc{middle, arc.end, arc.depth + 1}
		count += 2
		halvings++
	}

	return true
}

// reachAlong returns the supports of both shapes along the unit direction, and false where the
// inner one reaches farther than the outer one by more than the tolerance, or either is NaN.
func reachAlong[T Number](outer, inner Shape[T], direction Vector[float64], tolerance float64) (arcEnd, bool) {
	end := arcEnd{direction, inner.Support(direction), outer.Support(direction)}

	return end, end.inner.Subtract(end.outer).Dot(direction) <= tolerance
}

// arcEnd is a unit direction with the supports of the inner and the outer shape along it.
type arcEnd struct {
	direction    Vector[float64]
	inner, outer Point[float64]
}

// supportArc is the arc of directions between two ends, turning as Vector.Normal turns and less
// than a half turn, with how often it was halved from a quarter turn.
type supportArc struct {
	start, end arcEnd
	depth      int
}

// exceeds reports whether some direction of the arc may take the inner shape past the outer one by
// more than the tolerance. Over the arc the inner shape lies within the corner of its tangent lines
// at the ends and the outer one reaches at least its two supports, so the arc holds unless a
// direction takes that corner past both, decided by angles from the middle of the arc.
func (a supportArc) exceeds(tolerance float64) bool {
	corner := a.corner(a.start.inner, a.end.inner)
	middle := a.start.direction.Add(a.end.direction).Normalize()
	half := math.Atan2(a.start.direction.Cross(a.end.direction), a.start.direction.Dot(a.end.direction)) / 2

	low, high := -half, half
	for _, support := range [2]Point[float64]{a.start.outer, a.end.outer} {
		way := corner.Subtract(support)

		length := way.Length()
		if !(length > tolerance) {
			return false
		}

		angle, spread := math.Atan2(middle.Cross(way), middle.Dot(way)), math.Acos(tolerance/length)
		low, high = max(low, angle-spread), min(high, angle+spread)
	}

	return low < high
}

// settled reports whether both shapes are known over the arc within the tolerance, each corner of
// tangent lines lying within it of the chord between the supports, so halving further tells
// nothing.
func (a supportArc) settled(tolerance float64) bool {
	inner := Segment[float64]{a.start.inner, a.end.inner}
	outer := Segment[float64]{a.start.outer, a.end.outer}

	return lessOrEqualSquared(inner.distanceSquaredTo(a.corner(inner.Start, inner.End)), 0, tolerance) &&
		lessOrEqualSquared(outer.distanceSquaredTo(a.corner(outer.Start, outer.End)), 0, tolerance)
}

// corner returns where the lines through the points square to the two ends of the arc cross, each
// numerator divided once by the cross product of the directions, which the bound on halving keeps
// from vanishing.
func (a supportArc) corner(from, to Point[float64]) Point[float64] {
	start, end := a.start.direction, a.end.direction
	cross := start.Cross(end)
	first, second := from.Vector().Dot(start), to.Vector().Dot(end)

	return Point[float64]{
		(float64(first*end.Y) - float64(second*start.Y)) / cross,
		(float64(second*start.X) - float64(first*end.X)) / cross,
	}
}

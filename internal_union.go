package geom

import (
	"slices"
)

// unionOutline gathers the outline of the union of two convex polygons, appended to dst after from:
// the vertices of each outside the other and the crossings of their edges. It also averages the
// points both share, a point inside both about which the outline is sorted by angle.
type unionOutline[T Number] struct {
	dst            []Point[T]
	from           int
	sumX, sumY     float64
	count          int
	first, other   Point[T]
	paired, spread bool
}

// add takes a vertex of one polygon: onto the outline outside the other, among the shared points
// inside it.
func (u *unionOutline[T]) add(point Point[T], inside bool) {
	if inside {
		u.share(point)
	} else {
		u.dst = append(u.dst, point)
	}
}

// cross takes a crossing of the edges, on the outline and shared by both.
func (u *unionOutline[T]) cross(point Point[T]) {
	u.dst = append(u.dst, point)
	u.share(point)
}

// share adds a point both polygons hold to the average, and notes on exact signs whether the shared
// points so far enclose an area.
func (u *unionOutline[T]) share(point Point[T]) {
	u.sumX += float64(point.X)
	u.sumY += float64(point.Y)
	u.count++

	switch {
	case u.count == 1:
		u.first = point
	case !u.paired && !point.Equal(u.first):
		u.other, u.paired = point, true
	case u.paired && !u.spread:
		u.spread = u.other.Subtract(u.first).Cross(point.Subtract(u.first)) != 0
	}
}

// result sorts the outline by angle about the shared average in the winding of the first polygon,
// drops a point Equal to the one before, and starts it at the first point appended. Shared points
// that enclose no area leave dst as it was given.
func (u *unionOutline[T]) result(reversed bool) []Point[T] {
	if !u.spread || len(u.dst) == u.from {
		return u.dst[:u.from]
	}

	start := u.dst[u.from]
	center := Point[float64]{u.sumX / float64(u.count), u.sumY / float64(u.count)}
	outline := u.dst[u.from:]
	slices.SortFunc(outline, func(a, b Point[T]) int {
		return compareAngles(a.Float().Subtract(center), b.Float().Subtract(center))
	})
	if reversed {
		slices.Reverse(outline)
	}

	outline = slices.CompactFunc(outline, Point[T].Equal)
	if i := slices.IndexFunc(outline, start.Equal); i > 0 {
		slices.Reverse(outline[:i])
		slices.Reverse(outline[i:])
		slices.Reverse(outline)
	}

	return u.dst[:u.from+len(outline)]
}

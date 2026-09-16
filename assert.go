package geom

import (
	"fmt"

	"github.com/gravitton/assert"
)

// Testing is the subset of *testing.T the helpers need. It is an interface
// rather than testing.TB so a test can pass its own recorder.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// AssertNumber asserts that actual equals expected: exactly for an integer T, and for a
// float T within [EpsilonRelative], so the tolerance holds at any magnitude. Assertions use
// the scaled comparison rather than the absolute one [Equal] applies in hot paths.
func AssertNumber[T Number](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if isIntType[T]() {
		return assert.Equal(t, actual, expected, messages...)
	}

	return assert.EqualDelta(t, float64(actual), float64(expected), EpsilonRelative(actual, expected), messages...)
}

// AssertPoint asserts that p has the given X and Y values (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPoint[T Number](t Testing, p Point[T], x, y T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, p.X, x, append(messages, "X: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Y, y, append(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertVector asserts that v has the given X and Y values (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertVector[T Number](t Testing, p Vector[T], x, y T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, p.X, x, append(messages, "X: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Y, y, append(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertSize asserts that s has the given Width and Height values (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertSize[T Number](t Testing, s Size[T], w, h T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, s.Width, w, append(messages, "Width: ")...) {
		ok = false
	}
	if !AssertNumber(t, s.Height, h, append(messages, "Height: ")...) {
		ok = false
	}

	return ok
}

// AssertCircle asserts that c has the given center (x, y) and radius (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertCircle[T Number](t Testing, c Circle[T], x, y, radius T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, c.Center, x, y, append(messages, "Center.")...) {
		ok = false
	}
	if !AssertNumber(t, c.Radius, radius, append(messages, "Radius: ")...) {
		ok = false
	}

	return ok
}

// AssertLine asserts that l has the given start (sx, sy) and end (ex, ey) points (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertLine[T Number](t Testing, l Line[T], sx, sy, ex, ey T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, l.Start, sx, sy, append(messages, "Start.")...) {
		ok = false
	}
	if !AssertPoint(t, l.End, ex, ey, append(messages, "End.")...) {
		ok = false
	}

	return ok
}

// AssertRect asserts that r has the given center (cx, cy) and size (w, h) (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertRect[T Number](t Testing, r Rectangle[T], cx, cy, w, h T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, r.Center, cx, cy, append(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, r.Size, w, h, append(messages, "Size.")...) {
		ok = false
	}

	return ok
}

// AssertPolygon asserts that p has the given vertices (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPolygon[T Number](t Testing, p Polygon[T], vertices []Point[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, p.Vertices, vertices, messages...)
}

// AssertVertices asserts that vertices matches points element-by-element (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertVertices[T Number](t Testing, vertices []Point[T], points []Point[T], messages ...string) bool {
	t.Helper()

	if !assert.Equal(t, len(vertices), len(points), append(messages, "Length: ")...) {
		return false
	}

	ok := true
	for i := 0; i < len(vertices); i++ {
		if !AssertPoint(t, vertices[i], points[i].X, points[i].Y, append(messages, fmt.Sprintf("#%d.", i))...) {
			ok = false
		}
	}

	return ok
}

// AssertRegularPolygon asserts that p has the given center (x, y), size (w, h), vertex count n, and angle (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertRegularPolygon[T Number](t Testing, p RegularPolygon[T], x, y, w, h T, n int, angle float64, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, p.Center, x, y, append(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, p.Size, w, h, append(messages, "Size.")...) {
		ok = false
	}
	if !assert.Equal(t, p.N, n, append(messages, "N: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Angle, angle, append(messages, "Angle: ")...) {
		ok = false
	}

	return ok
}

// AssertPadding asserts that p has the given Top, Right, Bottom, and Left values (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPadding[T Number](t Testing, p Padding[T], top, right, bottom, left T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, p.Top, top, append(messages, "Top: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Right, right, append(messages, "Right: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Bottom, bottom, append(messages, "Bottom: ")...) {
		ok = false
	}
	if !AssertNumber(t, p.Left, left, append(messages, "Left: ")...) {
		ok = false
	}

	return ok
}

// AssertMatrix asserts that m is equal to expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertMatrix[T Number](t Testing, m, expected Matrix[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, m.A, expected.A, append(messages, "A: ")...) {
		ok = false
	}

	if !AssertNumber(t, m.B, expected.B, append(messages, "B: ")...) {
		ok = false
	}

	if !AssertNumber(t, m.C, expected.C, append(messages, "C: ")...) {
		ok = false
	}

	if !AssertNumber(t, m.D, expected.D, append(messages, "D: ")...) {
		ok = false
	}

	if !AssertNumber(t, m.E, expected.E, append(messages, "E: ")...) {
		ok = false
	}

	if !AssertNumber(t, m.F, expected.F, append(messages, "F: ")...) {
		ok = false
	}

	return ok
}

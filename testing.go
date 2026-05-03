package geom

import (
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

// AssertPoint asserts that p has the given X and Y values within Delta.
func AssertPoint[T Number](t *testing.T, p Point[T], x, y T, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(p.X), float64(x), Delta, append(messages, "X: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Y), float64(y), Delta, append(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertVector asserts that v has the given X and Y values within Delta.
func AssertVector[T Number](t *testing.T, p Vector[T], x, y T, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(p.X), float64(x), Delta, append(messages, "X: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Y), float64(y), Delta, append(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertSize asserts that s has the given Width and Height values within Delta.
func AssertSize[T Number](t *testing.T, s Size[T], w, h T, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(s.Width), float64(w), Delta, append(messages, "Width: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(s.Height), float64(h), Delta, append(messages, "Height: ")...) {
		ok = false
	}

	return ok
}

// AssertCircle asserts that c has the given center (x, y) and radius within Delta.
func AssertCircle[T Number](t *testing.T, c Circle[T], x, y, radius T, messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, c.Center, x, y, append(messages, "Center.")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(c.Radius), float64(radius), Delta, append(messages, "Radius: ")...) {
		ok = false
	}

	return ok
}

// AssertLine asserts that l has the given start (sx, sy) and end (ex, ey) points within Delta.
func AssertLine[T Number](t *testing.T, l Line[T], sx, sy, ex, ey T, messages ...string) bool {
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

// AssertRect asserts that r has the given center (cx, cy) and size (w, h) within Delta.
func AssertRect[T Number](t *testing.T, r Rectangle[T], cx, cy, w, h T, messages ...string) bool {
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

// AssertPolygon asserts that p has the given vertices within Delta.
func AssertPolygon[T Number](t *testing.T, p Polygon[T], vertices []Point[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, p.Vertices, vertices, messages...)
}

// AssertVertices asserts that vertices matches points element-by-element within Delta.
func AssertVertices[T Number](t *testing.T, vertices []Point[T], points []Point[T], messages ...string) bool {
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

// AssertRegularPolygon asserts that p has the given center (x, y), size (w, h), vertex count n, and angle within Delta.
func AssertRegularPolygon[T Number](t *testing.T, p RegularPolygon[T], x, y, w, h T, n int, angle float64, messages ...string) bool {
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
	if !assert.EqualDelta(t, p.Angle, angle, Delta, append(messages, "Angle: ")...) {
		ok = false
	}

	return ok
}

// AssertPadding asserts that p has the given Top, Right, Bottom, and Left values within Delta.
func AssertPadding[T Number](t *testing.T, p Padding[T], top, right, bottom, left T, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(p.Top), float64(top), Delta, append(messages, "Top: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Right), float64(right), Delta, append(messages, "Right: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Bottom), float64(bottom), Delta, append(messages, "Bottom: ")...) {
		ok = false
	}
	if !assert.EqualDelta(t, float64(p.Left), float64(left), Delta, append(messages, "Left: ")...) {
		ok = false
	}

	return ok
}

// AssertMatrix asserts that m is equal to expected within Delta.
func AssertMatrix[T Number](t *testing.T, m, expected Matrix[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.EqualDelta(t, float64(m.A), float64(expected.A), Delta, append(messages, "A: ")...) {
		ok = false
	}

	if !assert.EqualDelta(t, float64(m.B), float64(expected.B), Delta, append(messages, "B: ")...) {
		ok = false
	}

	if !assert.EqualDelta(t, float64(m.C), float64(expected.C), Delta, append(messages, "C: ")...) {
		ok = false
	}

	if !assert.EqualDelta(t, float64(m.D), float64(expected.D), Delta, append(messages, "D: ")...) {
		ok = false
	}

	if !assert.EqualDelta(t, float64(m.E), float64(expected.E), Delta, append(messages, "E: ")...) {
		ok = false
	}

	if !assert.EqualDelta(t, float64(m.F), float64(expected.F), Delta, append(messages, "F: ")...) {
		ok = false
	}

	return ok
}

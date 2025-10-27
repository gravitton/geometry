package geom

import (
	"fmt"
	"testing"

	"github.com/gravitton/assert"
)

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

func AssertPolygon[T Number](t *testing.T, p Polygon[T], vertices []Point[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, p.Vertices, vertices, messages...)
}

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

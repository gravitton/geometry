package geom

import (
	"fmt"
	"slices"

	"github.com/gravitton/assert"
)

// Testing is the subset of *testing.T the helpers need.
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

// AssertPoint asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPoint[T Number](t Testing, actual, expected Point[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.X, expected.X, prefixed(messages, "X: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Y, expected.Y, prefixed(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertVector asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertVector[T Number](t Testing, actual, expected Vector[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.X, expected.X, prefixed(messages, "X: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Y, expected.Y, prefixed(messages, "Y: ")...) {
		ok = false
	}

	return ok
}

// AssertSize asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertSize[T Number](t Testing, actual, expected Size[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.Width, expected.Width, prefixed(messages, "Width: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Height, expected.Height, prefixed(messages, "Height: ")...) {
		ok = false
	}

	return ok
}

// AssertCircle asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertCircle[T Number](t Testing, actual, expected Circle[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Radius, expected.Radius, prefixed(messages, "Radius: ")...) {
		ok = false
	}

	return ok
}

// AssertLine asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertLine[T Number](t Testing, actual, expected Line[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Start, expected.Start, prefixed(messages, "Start.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.End, expected.End, prefixed(messages, "End.")...) {
		ok = false
	}

	return ok
}

// AssertRectangle asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
// Angles are compared with EqualAngle, like Rectangle.Equal.
func AssertRectangle[T Number](t Testing, actual, expected Rectangle[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Size, expected.Size, prefixed(messages, "Size.")...) {
		ok = false
	}
	if !assert.True(t, EqualAngle(actual.Angle, expected.Angle), prefixed(messages, fmt.Sprintf("Angle: %v should equal %v modulo 2π: ", actual.Angle, expected.Angle))...) {
		ok = false
	}

	return ok
}

// AssertPolygon asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPolygon[T Number](t Testing, actual, expected Polygon[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, actual.Points, expected.Points, messages...)
}

// AssertVertices asserts that actual matches expected element-by-element (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertVertices[T Number](t Testing, actual, expected []Point[T], messages ...string) bool {
	t.Helper()

	if !assert.Equal(t, len(actual), len(expected), prefixed(messages, "Length: ")...) {
		return false
	}

	ok := true
	for i := range actual {
		if !AssertPoint(t, actual[i], expected[i], prefixed(messages, fmt.Sprintf("#%d.", i))...) {
			ok = false
		}
	}

	return ok
}

// AssertRegularPolygon asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
// Angles are compared with EqualAngle, like RegularPolygon.Equal.
func AssertRegularPolygon[T Number](t Testing, actual, expected RegularPolygon[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Size, expected.Size, prefixed(messages, "Size.")...) {
		ok = false
	}
	if !assert.Equal(t, actual.N, expected.N, prefixed(messages, "N: ")...) {
		ok = false
	}
	if !assert.True(t, EqualAngle(actual.Angle, expected.Angle), prefixed(messages, fmt.Sprintf("Angle: %v should equal %v modulo 2π: ", actual.Angle, expected.Angle))...) {
		ok = false
	}

	return ok
}

// AssertPadding asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertPadding[T Number](t Testing, actual, expected Padding[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.Top, expected.Top, prefixed(messages, "Top: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Right, expected.Right, prefixed(messages, "Right: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Bottom, expected.Bottom, prefixed(messages, "Bottom: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Left, expected.Left, prefixed(messages, "Left: ")...) {
		ok = false
	}

	return ok
}

// AssertMatrix asserts that actual equals expected (exactly for an integer T, within [EpsilonRelative] for a float T).
func AssertMatrix[T Number](t Testing, actual, expected Matrix[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.A, expected.A, prefixed(messages, "A: ")...) {
		ok = false
	}

	if !AssertNumber(t, actual.B, expected.B, prefixed(messages, "B: ")...) {
		ok = false
	}

	if !AssertNumber(t, actual.C, expected.C, prefixed(messages, "C: ")...) {
		ok = false
	}

	if !AssertNumber(t, actual.D, expected.D, prefixed(messages, "D: ")...) {
		ok = false
	}

	if !AssertNumber(t, actual.E, expected.E, prefixed(messages, "E: ")...) {
		ok = false
	}

	if !AssertNumber(t, actual.F, expected.F, prefixed(messages, "F: ")...) {
		ok = false
	}

	return ok
}

// prefixed returns messages followed by prefix in a fresh slice, so nested helpers never
// write into a backing array the caller still owns.
func prefixed(messages []string, prefix string) []string {
	return slices.Concat(messages, []string{prefix})
}

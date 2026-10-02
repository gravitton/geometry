// Package geomtest provides assertions for tests of code built on package geom, comparing a
// value field by field with the tolerance of its type and naming the field that differs.
package geomtest

import (
	"fmt"
	"slices"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

// AssertNumber asserts that actual equals expected: exactly for an integer T, and for a
// float T within [geom.EpsilonRelative], so the tolerance holds at any magnitude. Assertions use
// the scaled comparison rather than the absolute one [geom.Equal] applies in hot paths.
func AssertNumber[T geom.Number](t assert.Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if geom.Epsilon[T]() == 0 {
		return assert.Equal(t, actual, expected, messages...)
	}

	return assert.EqualDelta(t, float64(actual), float64(expected), geom.EpsilonRelative(actual, expected), messages...)
}

// AssertAngle asserts that actual equals expected as an angle in radians: within [geom.Delta] and
// modulo a full turn, as [geom.EqualAngle] compares them, the comparison every shape carrying an
// Angle makes. An angle is a float64 wherever the package stores one, so it is not generic.
func AssertAngle(t assert.Testing, actual, expected float64, messages ...string) bool {
	t.Helper()

	return assert.True(t, geom.EqualAngle(actual, expected), prefixed(messages, fmt.Sprintf("Angle: %v should equal %v modulo 2π: ", actual, expected))...)
}

// AssertPoint asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertPoint[T geom.Number](t assert.Testing, actual, expected geom.Point[T], messages ...string) bool {
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

// AssertVector asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertVector[T geom.Number](t assert.Testing, actual, expected geom.Vector[T], messages ...string) bool {
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

// AssertSize asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertSize[T geom.Number](t assert.Testing, actual, expected geom.Size[T], messages ...string) bool {
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

// AssertCircle asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertCircle[T geom.Number](t assert.Testing, actual, expected geom.Circle[T], messages ...string) bool {
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

// AssertEllipse asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
// The angle goes through AssertAngle, so a full turn does not matter, like Ellipse.Equal.
func AssertEllipse[T geom.Number](t assert.Testing, actual, expected geom.Ellipse[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Size, expected.Size, prefixed(messages, "Size.")...) {
		ok = false
	}
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}

	return ok
}

// AssertSegment asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertSegment[T geom.Number](t assert.Testing, actual, expected geom.Segment[T], messages ...string) bool {
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

// AssertSegments asserts that actual holds the expected segments in order, each by AssertSegment,
// the parts a Clip method returns.
func AssertSegments[T geom.Number](t assert.Testing, actual, expected []geom.Segment[T], messages ...string) bool {
	t.Helper()

	if !assert.Equal(t, len(actual), len(expected), prefixed(messages, "Length: ")...) {
		return false
	}

	ok := true
	for i := range actual {
		if !AssertSegment(t, actual[i], expected[i], prefixed(messages, fmt.Sprintf("#%d.", i))...) {
			ok = false
		}
	}

	return ok
}

// AssertRay asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertRay[T geom.Number](t assert.Testing, actual, expected geom.Ray[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Origin, expected.Origin, prefixed(messages, "Origin.")...) {
		ok = false
	}
	if !AssertVector(t, actual.Direction, expected.Direction, prefixed(messages, "Direction.")...) {
		ok = false
	}

	return ok
}

// AssertRectangle asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
// The angle goes through AssertAngle, so a full turn does not matter, like Rectangle.Equal.
func AssertRectangle[T geom.Number](t assert.Testing, actual, expected geom.Rectangle[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Size, expected.Size, prefixed(messages, "Size.")...) {
		ok = false
	}
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}

	return ok
}

// AssertBox asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertBox[T geom.Number](t assert.Testing, actual, expected geom.Box[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Min, expected.Min, prefixed(messages, "Min.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.Max, expected.Max, prefixed(messages, "Max.")...) {
		ok = false
	}

	return ok
}

// AssertPolygon asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertPolygon[T geom.Number](t assert.Testing, actual, expected geom.Polygon[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, actual.Points, expected.Points, messages...)
}

// AssertVertices asserts that actual matches expected element-by-element (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertVertices[T geom.Number](t assert.Testing, actual, expected []geom.Point[T], messages ...string) bool {
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

// AssertRegularPolygon asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
// The angle and the phase go through AssertAngle, so a full turn does not matter, like RegularPolygon.Equal.
func AssertRegularPolygon[T geom.Number](t assert.Testing, actual, expected geom.RegularPolygon[T], messages ...string) bool {
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
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}
	if !AssertAngle(t, actual.Phase, expected.Phase, prefixed(messages, "Phase.")...) {
		ok = false
	}

	return ok
}

// AssertPadding asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertPadding[T geom.Number](t assert.Testing, actual, expected geom.Padding[T], messages ...string) bool {
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

// AssertMatrix asserts that actual equals expected (exactly for an integer T, within [geom.EpsilonRelative] for a float T).
func AssertMatrix[T geom.Number](t assert.Testing, actual, expected geom.Matrix[T], messages ...string) bool {
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

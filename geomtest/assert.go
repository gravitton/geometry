// Package geomtest provides test assertions for code built on package geom. Each compares the
// values field by field, exactly for an integer T and within [geom.EpsilonRelative] for a float T,
// reports every field that differs by its name, and returns whether the values are equal.
package geomtest

import (
	"fmt"
	"slices"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

// Testing is the part of testing.TB the assertions report through, so a *testing.T, a *testing.B or
// a recorder of the caller's own can be passed.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// AssertNumber asserts that actual equals expected.
func AssertNumber[T geom.Number](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if geom.Epsilon[T]() == 0 {
		return assert.Equal(t, actual, expected, messages...)
	}

	return assert.EqualDelta(t, float64(actual), float64(expected), geom.EpsilonRelative(actual, expected), messages...)
}

// AssertAngle asserts that actual equals expected as an angle in radians, by [geom.EqualAngle].
func AssertAngle(t Testing, actual, expected float64, messages ...string) bool {
	t.Helper()

	return assert.True(t, geom.EqualAngle(actual, expected), prefixed(messages, fmt.Sprintf("Angle: %v should equal %v modulo 2π: ", actual, expected))...)
}

// AssertPoint asserts that actual equals expected.
func AssertPoint[T geom.Number](t Testing, actual, expected geom.Point[T], messages ...string) bool {
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

// AssertVector asserts that actual equals expected.
func AssertVector[T geom.Number](t Testing, actual, expected geom.Vector[T], messages ...string) bool {
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

// AssertSize asserts that actual equals expected.
func AssertSize[T geom.Number](t Testing, actual, expected geom.Size[T], messages ...string) bool {
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

// AssertPadding asserts that actual equals expected.
func AssertPadding[T geom.Number](t Testing, actual, expected geom.Padding[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertNumber(t, actual.Left, expected.Left, prefixed(messages, "Left: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Right, expected.Right, prefixed(messages, "Right: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Top, expected.Top, prefixed(messages, "Top: ")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Bottom, expected.Bottom, prefixed(messages, "Bottom: ")...) {
		ok = false
	}

	return ok
}

// AssertMatrix asserts that actual equals expected.
func AssertMatrix[T geom.Number](t Testing, actual, expected geom.Matrix[T], messages ...string) bool {
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

// AssertCircle asserts that actual equals expected.
func AssertCircle[T geom.Number](t Testing, actual, expected geom.Circle[T], messages ...string) bool {
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

// AssertCapsule asserts that actual equals expected.
func AssertCapsule[T geom.Number](t Testing, actual, expected geom.Capsule[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertSegment(t, actual.Segment, expected.Segment, prefixed(messages, "Segment.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Radius, expected.Radius, prefixed(messages, "Radius: ")...) {
		ok = false
	}

	return ok
}

// AssertSector asserts that actual equals expected, the angle by AssertAngle and the sweep as a
// number.
func AssertSector[T geom.Number](t Testing, actual, expected geom.Sector[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Radius, expected.Radius, prefixed(messages, "Radius: ")...) {
		ok = false
	}
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}
	if !AssertNumber(t, actual.Sweep, expected.Sweep, prefixed(messages, "Sweep: ")...) {
		ok = false
	}

	return ok
}

// AssertArc asserts that actual equals expected, the angle by AssertAngle and the sweep as a
// number.
func AssertArc[T geom.Number](t Testing, actual, expected geom.Arc[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Radius, expected.Radius, prefixed(messages, "Radius: ")...) {
		ok = false
	}
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}
	if !AssertNumber(t, actual.Sweep, expected.Sweep, prefixed(messages, "Sweep: ")...) {
		ok = false
	}

	return ok
}

// AssertSegment asserts that actual equals expected.
func AssertSegment[T geom.Number](t Testing, actual, expected geom.Segment[T], messages ...string) bool {
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

// AssertSegments asserts that actual holds the expected segments in order, each by AssertSegment.
func AssertSegments[T geom.Number](t Testing, actual, expected []geom.Segment[T], messages ...string) bool {
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

// AssertRay asserts that actual equals expected.
func AssertRay[T geom.Number](t Testing, actual, expected geom.Ray[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Origin, expected.Origin, prefixed(messages, "Origin.")...) {
		ok = false
	}
	if !AssertVector(t, actual.Vector, expected.Vector, prefixed(messages, "Vector.")...) {
		ok = false
	}

	return ok
}

// AssertLine asserts that actual equals expected.
func AssertLine[T geom.Number](t Testing, actual, expected geom.Line[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Origin, expected.Origin, prefixed(messages, "Origin.")...) {
		ok = false
	}
	if !AssertVector(t, actual.Vector, expected.Vector, prefixed(messages, "Vector.")...) {
		ok = false
	}

	return ok
}

// AssertPolyline asserts that actual equals expected point by point.
func AssertPolyline[T geom.Number](t Testing, actual, expected geom.Polyline[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, actual.Points, expected.Points, messages...)
}

// AssertQuadraticBezier asserts that actual equals expected point by point.
func AssertQuadraticBezier[T geom.Number](t Testing, actual, expected geom.QuadraticBezier[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Start, expected.Start, prefixed(messages, "Start.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.Control, expected.Control, prefixed(messages, "Control.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.End, expected.End, prefixed(messages, "End.")...) {
		ok = false
	}

	return ok
}

// AssertCubicBezier asserts that actual equals expected point by point.
func AssertCubicBezier[T geom.Number](t Testing, actual, expected geom.CubicBezier[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Start, expected.Start, prefixed(messages, "Start.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.Control1, expected.Control1, prefixed(messages, "Control1.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.Control2, expected.Control2, prefixed(messages, "Control2.")...) {
		ok = false
	}
	if !AssertPoint(t, actual.End, expected.End, prefixed(messages, "End.")...) {
		ok = false
	}

	return ok
}

// AssertPolygon asserts that actual equals expected point by point.
func AssertPolygon[T geom.Number](t Testing, actual, expected geom.Polygon[T], messages ...string) bool {
	t.Helper()

	return AssertVertices(t, actual.Points, expected.Points, messages...)
}

// AssertVertices asserts that actual holds the expected points in order. A different count fails
// without comparing the points.
func AssertVertices[T geom.Number](t Testing, actual, expected []geom.Point[T], messages ...string) bool {
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

// AssertRectangle asserts that actual equals expected, the angle by AssertAngle.
func AssertRectangle[T geom.Number](t Testing, actual, expected geom.Rectangle[T], messages ...string) bool {
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

// AssertRegularPolygon asserts that actual equals expected, the angle and the phase by AssertAngle.
func AssertRegularPolygon[T geom.Number](t Testing, actual, expected geom.RegularPolygon[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Radius, expected.Radius, prefixed(messages, "Radius.")...) {
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

// AssertBox asserts that actual equals expected.
func AssertBox[T geom.Number](t Testing, actual, expected geom.Box[T], messages ...string) bool {
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

// AssertEllipse asserts that actual equals expected, the angle by AssertAngle.
func AssertEllipse[T geom.Number](t Testing, actual, expected geom.Ellipse[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Center, expected.Center, prefixed(messages, "Center.")...) {
		ok = false
	}
	if !AssertSize(t, actual.Radius, expected.Radius, prefixed(messages, "Radius.")...) {
		ok = false
	}
	if !AssertAngle(t, actual.Angle, expected.Angle, messages...) {
		ok = false
	}

	return ok
}

// AssertHit asserts that actual equals expected, the normal and the fraction within
// [geom.EpsilonRelative] whatever T is.
func AssertHit[T geom.Number](t Testing, actual, expected geom.Hit[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertPoint(t, actual.Point, expected.Point, prefixed(messages, "Point.")...) {
		ok = false
	}
	if !AssertVector(t, actual.Normal, expected.Normal, prefixed(messages, "Normal.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Fraction, expected.Fraction, prefixed(messages, "Fraction: ")...) {
		ok = false
	}

	return ok
}

// AssertContact asserts that actual equals expected, the normal and the depth within
// [geom.EpsilonRelative] whatever T is, and the points as many and each in its place.
func AssertContact[T geom.Number](t Testing, actual, expected geom.Contact[T], messages ...string) bool {
	t.Helper()

	ok := true

	if !AssertVector(t, actual.Normal, expected.Normal, prefixed(messages, "Normal.")...) {
		ok = false
	}
	if !AssertNumber(t, actual.Depth, expected.Depth, prefixed(messages, "Depth: ")...) {
		ok = false
	}
	if !AssertVertices(t, actual.Points.Items[:actual.Points.Count], expected.Points.Items[:expected.Points.Count], prefixed(messages, "Points.")...) {
		ok = false
	}

	return ok
}

// prefixed returns messages followed by prefix in a fresh slice, so nested helpers never write into
// a backing array the caller still owns.
func prefixed(messages []string, prefix string) []string {
	return slices.Concat(messages, []string{prefix})
}

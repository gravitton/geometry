// Package geom provides a small, generic, and immutable 2D geometry toolkit.
//
// Design goals
//   - Generic: All core types are parameterized by a Number constraint so the
//     same API works with ints and floats.
//   - Immutable: Methods do not mutate receivers; they return new values.
//   - Practical: Focused on game/graphics use-cases with clear, minimal API.
//
// # Integer rounding
//
// A float result stored into an integer T goes through Cast, which rounds half away from
// zero: Lerp, Midpoint, Multiply, Divide, Int and every method built on them follow it, so
// Pt(0, 0).Midpoint(Pt(5, 5)) is (3,3). The exception is Rectangle, whose center is placed
// by truncating half the size toward Min so that Max-Min stays exactly the size:
// RectangleFromMinMax(Pt(0, 0), Pt(5, 5)).Center is (2,2), and Line.Bounds().Center can
// therefore differ from Line.Midpoint() by one unit on an odd span.
package geom

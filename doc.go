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
//
// # Panics
//
// Divide, and every method built on it (Point.Divide, Padding.Unscale, and the Unscale of every
// shape and of Matrix), panics
// for a zero factor, and Matrix.Inverse panics for a singular matrix, the same way the integer
// / operator and Mod do. Check IsInvertible before inverting a matrix that may be singular.
// Cast panics when a NaN or ±Inf would be stored into an integer T: Multiply, Lerp, Rotate,
// Transform and Int on an integer shape all go through it. A float T carries NaN and ±Inf
// through unchanged. A finite value outside the range of an integer T is not checked and
// stores a platform-dependent value, as Cast documents. RegularPolygonOrientationAngle panics
// for an Orientation that is neither FlatTop nor PointyTop, OrientationNone included: the
// absence of an alignment has no angle to give. RegularPolygon.Lerp panics for a polygon with
// a different vertex count, which has no shape between.
// These are the only panics: every other guard returns a value the type can express, such as
// DirectionNone, an empty polygon, or the 0 that Size.AspectRatio gives for a zero height.
//
// # Arithmetic
//
// Products, distances and interpolations are computed in float64 and stored back through Cast,
// whatever T is, so a narrow integer T never overflows mid-computation and every integer result
// follows the one rounding rule above. Sums and differences of two values stay in T.
//
// # Boundaries
//
// Rectangle.Contains, Circle.Contains, Line.Contains, Polygon.Contains, Vector.LessOrEqual and
// the Intersects methods are closed and tolerant: a point within Epsilon of the boundary counts
// as on it, so a float rectangle contains the corners it was built from even where Min is
// recomputed with a rounding error. Every such test is one comparison on a squared distance,
// so containment, the distance methods and the intersection tests round alike at the boundary.
// Vector.Less is the strict counterpart and applies no tolerance.
// DistanceTo is zero exactly where Contains holds, and Intersection returns a point or a
// rectangle exactly where Intersects holds, less the parallel and coincident cases that have
// no single answer. The boundary crossings of a segment, IntersectionCircle, IntersectionRectangle
// and IntersectionPolygon, follow the same rule with one more exception: a segment entirely
// inside a shape crosses no boundary and returns none while Intersects reports it.
package geom

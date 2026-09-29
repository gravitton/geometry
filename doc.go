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
// the midpoint of an odd span rounds up. The exception is Rectangle and Box, whose center is
// placed by truncating half the size toward Min so that Max-Min stays exactly the size, and
// Segment.Bounds().Center() can therefore differ from Segment.Midpoint() by one unit on an odd span.
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
// for an Orientation that is neither OrientationFlatTop nor OrientationPointyTop, OrientationNone included: the
// absence of an alignment has no angle to give. Polygon.Lerp and RegularPolygon.Lerp panic for
// a polygon with a different vertex count, which has no shape between. Intersects panics for two Colliders
// of another package, which have no method this package can reach.
// These are the only panics: every other guard returns a value the type can express, such as
// DirectionNone, an empty polygon, or the 0 that Size.AspectRatio gives for a zero height.
//
// # Arithmetic
//
// Products, distances and interpolations are computed in float64 and stored back through Cast,
// whatever T is, so a narrow integer T never overflows mid-computation and every integer result
// follows the one rounding rule above. Sums and differences of two values stay in T.
//
// # Reproducibility
//
// Go lets the compiler fuse a product and a sum into one multiply-add, which rounds once where
// the two steps round twice, and gc does so on arm64 and on amd64 built for GOAMD64=v3. Every
// product the package adds to another value is rounded on its own first, so the same inputs give
// the same bits on every architecture, and an integer T rounds the same way on all of them. The
// exception is what the math package computes: Sincos, Atan2, Hypot and the rest promise no
// identical bits across architectures, so a method placing a point from an angle, such as Rotate,
// VectorFromAngle or the vertices of a RegularPolygon, can differ in the last bit from one
// architecture to another, and for an integer T round the other way where the value lies that
// close to a half.
//
// # Inlining
//
// The methods that store two coordinates through Cast, such as Vector.Multiply and
// Point.Transform, cost more than the compiler's inlining budget, and on amd64 Cast itself does,
// since math.Round is no intrinsic there. A profile-guided build (go build -pgo) inlines them at
// the call sites the profile marks hot, which is where it matters.
//
// # Boundaries
//
// Rectangle.Contains, Box.Contains, Circle.Contains, Ellipse.Contains, Segment.Contains, Ray.Contains,
// Polygon.Contains, RegularPolygon.Contains, Vector.LessOrEqual and the Intersects methods are closed and tolerant: a point within Epsilon of the boundary counts
// as on it, so a float rectangle contains the corners it was built from even where Min is
// recomputed with a rounding error. Every such test is one comparison on a squared distance,
// never on a coordinate, so containment, the distance methods and the intersection tests
// round alike at the boundary. Vector.Less is the strict counterpart and applies no tolerance.
// DistanceTo is zero exactly where Contains holds, and IntersectionSegment and
// IntersectionRectangle return a point or a rectangle exactly where the Intersects methods hold,
// less the parallel and coincident segments and the rectangles of different angles that have
// no single answer. The boundary crossings of a segment, IntersectionCircle, IntersectionRectangle
// and IntersectionPolygon, follow the same rule with one more exception: a segment entirely
// inside a shape crosses no boundary and returns none while Intersects reports it.
package geom

// Package geom provides a small, generic and immutable 2D geometry toolkit for games and UI.
//
// # Design goals
//
//   - Generic: every type is parameterized by the [Number] constraint, so the same API works
//     with integers and floats, named types included.
//   - Immutable: no method mutates its receiver; each returns a new value.
//   - Practical: focused on game and graphics use cases, in screen space with the origin at the
//     top left and +Y pointing down.
//
// # Types
//
// The values are [Point], [Vector], [Size], [Padding] and the affine [Matrix]. The shapes are
// [Rectangle], which may be turned by an angle, the axis-aligned [Box], [Circle], [Ellipse],
// [Segment], [Polygon] and [RegularPolygon], and the half-line [Ray]. [Shape], [Collider] and
// [Body] name what the shapes share, and [Intersects] tests two colliders of any kind.
// [Direction], [Axis], [Orientation] and [Winding] are enums. Package geomtest holds test
// assertions, and packages ints and floats alias the int and float64 instantiations.
//
// Every value type has Equal, Cast, Int, Float and String; every shape adds Contains, every one
// but Ray adds Bounds, and every one but Ellipse adds the Intersects methods. The enums are compared with == and
// marshal as their names.
//
// # Coordinates
//
// The origin is top-left and +Y points down. Angles follow the mathematical convention, so a
// positive angle is counterclockwise in math coordinates and appears clockwise on screen.
// Direction order and polygon winding follow the same rule: [Directions], [Rectangle.Vertices]
// and [RegularPolygon.Vertices] all wind by increasing angle, clockwise as drawn. Each Direction
// has three names: canonical ([DirectionDownRight]), compass ([SouthEast]) and rectangle corner
// ([BottomRight]).
//
// # Numbers
//
// Products, distances and interpolations are computed in float64 and stored back through [Cast],
// whatever T is. Cast rounds half away from zero: Lerp, Midpoint, Multiply, Divide, Int and
// every method built on them follow it, so the midpoint of an odd span rounds up. The exception is [Rectangle] and [Box], whose center is placed
// by truncating half the size toward Min so that Max-Min stays exactly the size, and
// Segment.Bounds().Center() can therefore differ from [Segment.Midpoint] by one unit on an odd
// span. Sums and differences of two values stay in T. The signs that decide a crossing, a turn or
// a winding are cross products of coordinate differences in float64, exact for an integer T while
// the differences stay within 2^26, the square root of the integers float64 holds exactly; beyond
// it a sign can round to zero.
//
// # Supported range
//
// The package is tested for float64 at any finite coordinates, for float32 within 1e5 of the
// origin, and for int, int32 and int64 while coordinate differences stay within 2^26. Within
// that range no method gives an answer wrong by more than the tolerance, but where two methods
// place the same point by different roundings, which only a turned shape or a ray does, they
// can disagree about a point within the tolerance of a boundary, and float32 far from the origin
// is where that shows. A narrow integer T such as int8 or int16 compiles and never panics or
// hangs, but a sum or difference near the end of its range may wrap. The extent of a shape
// must itself fit T, a circle's center plus its radius or the point a ray reaches a shape at,
// and the square of a length must stay finite, so a segment or a ray direction past about 1e150
// measures nothing.
//
// # Equality
//
// Equal compares an integer T exactly and a float T within [Epsilon] of T, a tolerance matched to
// float32 or float64. [EqualRelative] scales it for values far from zero, and [EqualAngle]
// compares modulo a full turn.
//
// # Boundaries
//
// The Contains method of every shape, [Point.Between], [Vector.LessOrEqual] and the Intersects
// methods are closed and tolerant: a point within the tolerance of the boundary counts as on it,
// so a float rectangle contains the corners it was built from. The tolerance is [Epsilon] of T
// near the origin, widens to two ulps of T far from it, and is zero for an integer T. Every such
// test is one comparison on a squared distance, so DistanceTo is zero exactly where Contains
// holds, and each Intersection method answers exactly where its Intersects holds, apart from
// parallel and coincident segments, rectangles of different angles, and a segment entirely
// inside a shape, which crosses no boundary. The logic of each pair is written once and the other
// side delegates to it, so both sides give the same answer. [Vector.Less] is strict.
//
// # Matrices
//
// An integer [Matrix] composes lattice transforms exactly: translation, integer scale,
// reflection and quarter turns. Anything else rounds into a different matrix; use a float Matrix
// there. Transform takes a float matrix, so convert an integer one with Float at the call.
//
// # Panics
//
// [Divide], and every method built on it ([Point.Divide], [Padding.Unscale], and the Unscale of
// every shape but [Box], which has none, and of [Matrix]), panics for a zero factor, and
// [Matrix.Inverse] panics for a singular matrix, the same way the integer / operator and [Mod] do.
// Check [Matrix.IsInvertible] before inverting a matrix that may be singular. [Cast] panics when a
// NaN or ±Inf would be stored into an integer T: Multiply, Lerp, Rotate, Transform and Int on an
// integer shape all go through it. A float T carries NaN and ±Inf through unchanged. A finite
// value outside the range of an integer T is not checked and stores a platform-dependent value, as
// Cast documents. [RegularPolygonOrientationPhase] panics for an [Orientation] that is neither
// [OrientationFlatTop] nor [OrientationPointyTop], [OrientationNone] included: the absence of an
// alignment has no phase to give. [Polygon.Lerp] and [RegularPolygon.Lerp] panic for a polygon
// with a different vertex count, which has no shape between. [Intersects] panics for two colliders
// of another package, which have no method this package can reach, a pointer to a shape of this
// package included. An integer [RegularPolygon] with a NaN or infinite Angle or Phase panics
// wherever it places a vertex, through Cast. These are the only panics:
// every other guard returns a value the type can express, such as [DirectionNone], an empty
// polygon, or the 0 that [Size.AspectRatio] gives for a zero height.
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
// # JSON
//
// The flat JSON of [Rectangle], [Circle] and [RegularPolygon] relies on the embed struct tag of
// the v2-backed encoding/json, the default since Go 1.27; under GOEXPERIMENT=nojsonv2 the nested
// fields are emitted as objects.
package geom

// Package geom provides generic, immutable 2D geometry for games and UI.
//
// Every type is parameterized by the [Number] constraint, so one API serves integers and floats,
// named types included. No method changes its receiver: each returns a new value, or the receiver
// when nothing changes.
//
// # Types
//
//   - Values: [Point], [Vector], [Size], [Padding] and the affine [Matrix].
//   - Shapes: [Circle], [Capsule], [Sector], [Segment], [Polygon], [Rectangle] (which may be turned
//     by an angle), [RegularPolygon], the axis-aligned [Box] and [Ellipse]. Each is a [Shape].
//   - Paths without an inside: [Ray], [Line], [Polyline], [Arc] and the curves [QuadraticBezier]
//     and [CubicBezier], which are measured through the Polyline they flatten to.
//   - Results: [Hit], [Contact] and [Pair].
//   - Enums: [Direction], [Axis], [Orientation] and [Winding], compared with == and marshaled as
//     their names.
//
// Shapes are tested against each other by their Intersects<Kind>, Contains<Kind> and
// Intersection<Kind> methods, and against a shape of any kind by [Intersects], [Contains],
// [Distance], [Nearest], [Collide], [Sweep] and [Raycast]. Every type but the enums has Equal,
// IsZero and String; geometry also has Cast, Int and Float, results do not. Package geomtest holds
// test assertions, and packages ints and floats alias the int and float64 instantiations.
//
// # Naming
//
// [Point] and [Vector] are math objects and take the words of vector algebra: Multiply, Divide,
// [Vector.Dot], [Vector.Cross], [Vector.Normal]. Shapes take the words of geometry: Scale,
// Unscale, Rotate, Centroid. A word means one thing on each side, so [Vector.Normal] is the
// perpendicular of the same length while [Hit] holds a Normal of unit length.
//
// # Coordinates
//
// The origin is top-left and +Y points down. Angles are in radians and positive counterclockwise
// in math coordinates, which appears clockwise on screen. [Directions], [Rectangle.Vertices] and
// [RegularPolygon.Vertices] wind by increasing angle. Each Direction has a canonical, a compass and
// a corner name ([DirectionDownRight], [SouthEast], [BottomRight]).
//
// # Boundaries
//
// Every Contains and Intersects test is closed and tolerant: a point within the tolerance of a
// boundary is on it, so a float rectangle contains the corners it was built from. The tolerance
// is [Epsilon] of T near the origin, widens to two ulps of T far from it, and is zero for an
// integer T. Related methods decide by the same comparison, so DistanceTo is zero exactly where
// Contains holds, a pair gives the same answer from either side, and Intersection<Kind> gives
// points exactly where Intersects<Kind> holds, except where no boundary is crossed: a shape
// entirely inside another, and parallel or coincident segments. An empty polygon contains,
// intersects and is contained by nothing. [Vector.Less] is strict.
//
// Equal compares an integer T exactly and a float T within [Epsilon]; [EqualRelative] scales the
// tolerance with the values and [EqualAngle] compares modulo a full turn.
//
// # Numbers
//
// Products, distances and interpolations are computed in float64 and stored once through [Cast],
// which rounds half away from zero; sums and differences stay in T. For an integer T a constructed
// point is therefore rounded and may lie just off the boundary it was placed on. The center of an
// integer [Rectangle] or [Box] truncates toward Min instead, so Max-Min stays exactly the size. A
// result that rounds onto one it already holds is given once.
//
// Results hold at most what the geometry allows and come in a defined order. Up to two are a
// [Pair], which allocates nothing; more are a slice, allocated once, with an Append form that
// reuses the caller's slice where no allocation-free alternative exists.
//
// # Supported range
//
// float64 is supported at any finite coordinates, float32 within 1e5 of the origin, and int,
// int32 and int64 while coordinate differences stay within 2^26, where the signs of crossings
// and turns are exact. Within that range no answer is wrong by more than the tolerance, though
// two methods placing the same point by different roundings, as a turned shape, a ray or a line
// does, can disagree within the tolerance of a boundary. int8 and int16 compile and never panic
// or hang, but carry no promise of precision. The extent of a shape must fit T, and the square of
// a length must stay finite, so a segment or direction past about 1e150 measures nothing.
//
// # Matrices
//
// An integer [Matrix] composes lattice transforms exactly: translation, integer scale, reflection
// and quarter turns; anything else needs a float Matrix. Transform takes a float matrix, so
// convert an integer one with Float.
//
// # Panics
//
// The only panics are for arguments with no meaning:
//
//   - a zero factor to [Divide] and every method built on it, Unscale included;
//   - a singular matrix to [Matrix.Inverse] (check [Matrix.IsInvertible]);
//   - NaN or ±Inf stored into an integer T by [Cast], so through Multiply, Lerp, Rotate, Transform
//     or Int of an integer value, and an integer [RegularPolygon] with a non-finite Angle or Phase;
//   - an [Orientation] other than [OrientationFlatTop] or [OrientationPointyTop] to
//     [RegularPolygonOrientationPhase] and every function built on it;
//   - a differing vertex count to [Polygon.Lerp], [RegularPolygon.Lerp] and [Polyline.Lerp].
//
// A float T carries NaN and ±Inf through. A finite value outside the range of an integer T stores
// a platform-dependent value. Every other input returns a value, such as [DirectionNone] or an
// empty polygon.
//
// # Reproducibility
//
// Every product is rounded before it is added, so no fused multiply-add changes the result and the
// same inputs give the same bits on every architecture. The exception is what the math package
// computes: a point placed from an angle, by Rotate or the vertices of a RegularPolygon, can
// differ in the last bit between architectures.
//
// # Performance
//
// Methods that store two coordinates through Cast exceed the inlining budget; a profile-guided
// build (go build -pgo) inlines them where the profile marks them hot.
//
// # JSON
//
// Keys are one character. Shapes with a center or a segment among their fields embed it flat,
// which relies on the embed tag of the v2-backed encoding/json, the default since Go 1.27; under
// GOEXPERIMENT=nojsonv2 those fields are nested objects.
package geom

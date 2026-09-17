# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.12.0...main)
### Added
- `LessOrEqual(a, b)` and `LessOrEqualDelta(a, b, delta)` – ordered comparison within `Epsilon[T]()` or a given delta, the `<=` counterparts of `Equal` and `EqualDelta`; `Rectangle.Contains`, `CollisionRectangles` and `Vector.LessOrEqual` are built on them
- `Padding.Add`, `Padding.Negate` and `Padding.Scale`, so a padding can follow a scaled rectangle
- `Rectangle.Outset(padding)` – expands the rectangle by the padding, the inverse of `Inset`
- `Int[T Number](value T) int` – converts a `Number` to `int`, an integer `T` directly and a float `T` rounded through `Cast`; the `Int()` methods are built on it
- `Delta32` – the equality tolerance for `float32`, `1e-4`. A `float32` ulp reaches 6e-5 at magnitude 1e3, so `Delta` sat below one ulp across the coordinate range `float32` is normally chosen for and asked for bit-exact equality there
- `Epsilon[T]()` – the tolerance for `T`: zero for an integer `T`, `Delta32` for `float32`, `Delta` for `float64`. It depends only on the type, never on the values compared, so `Equal` costs what it did before
- `EqualRelative(a, b)` and `EpsilonRelative(a, b)` – equality within a tolerance that scales with magnitude, for values far from zero where an absolute tolerance falls below one ulp. Measured at ~1.1 ns against ~0.6 ns for `Equal`, which is why it is a separate function rather than a change to `Equal`
- `AssertNumber(t, actual, expected)` – asserts a bare coordinate with the same integer/float rule as the shape helpers
- `Vector.LessOrEqual(length)` – reports whether the vector is at most the given length, the closed counterpart of `Less`
- `Padding.Equal` and `Padding.IsZero`, so `Padding` has the same equality pair as every other type
- `Polygon.Bounds` – the axis-aligned bounding rectangle of the vertices, the zero rectangle for an empty polygon; every shape now has `Bounds()`
- `Matrix.IsInvertible` – reports whether the determinant is non-zero, so a caller can tell a real inverse from the unchanged matrix `Inverse` returns for a singular one
- `AngleDistance(a, b)` – the shortest angular distance between two angles, in `[0, π]`
- `EqualAngle(a, b)` – reports whether two angles are equal modulo a full turn within `Delta`, holding across the `0`/`2π` seam where comparing normalized angles does not

### Changed
- `Rectangle.Contains`, `Circle.Contains`, `Vector.LessOrEqual`, `CollisionRectangles`, `CollisionCircles` and `CollisionRectangleCircle` include the boundary within `Epsilon[T]()` instead of comparing exactly, so a float rectangle contains the corners it was built from even where `Min` is recomputed from `Center` with a rounding error, and a circle contains its anchors; `Vector.Less` stays strict (**breaking**)
- Every product, distance and interpolation is computed in `float64` and stored back through `Cast`: `Vector.Dot`, `Cross`, `LengthSquared`, `Less`, `LessOrEqual`, `IsNormalized`, `Point.ManhattanDistanceTo`, `ChebyshevDistanceTo`, `OctileDistanceTo`, `Size.Area`, `Perimeter`, `Polygon.Center`, `Matrix.Multiply`, `Determinant`, `Inverse` and `IsInvertible`, so a narrow `int8` or `int16` no longer overflows mid-computation and an `int64` beyond 2^53 loses precision instead (**breaking** for such values)
- `Polygon.Center` rounds the vertex average half away from zero like every other integer result instead of truncating toward zero, so an integer `MoveTo` now lands on the point except when a half average changes sign (**breaking**)
- `Directions`, `CardinalDirections`, `DiagonalDirections` and `Axes` are functions returning a fresh array instead of package-level variables an importer could write into (**breaking**)
- `RegularPolygonOrientationAngle` panics for an `Orientation` other than `FlatTop` and `PointyTop` instead of returning `0`
- `RegularPolygon.Vertices` returns nil for `N < 1` instead of an empty slice, so `RegularPolygon.Polygon()` of an empty polygon is zero like `Pol(nil)`
- `Cast` documents that a finite value outside the range of an integer `T` is not checked and stores a platform-dependent value, `Int` that an `int64` truncates on a 32-bit target, and `Rectangle.Inset` that `Outset` undoes it
- `Line` marshals `Start` and `End` under the JSON keys `s` and `e` instead of `a` and `b`, the initials of the fields like every other key in the package (**breaking**)
- `Cast` panics when a `NaN` or `±Inf` would be stored into an integer `T`, the same convention `Divide` follows for a zero scale, so `Multiply`, `Lerp`, `Vector.Rotate`, `Transform` and `Int()` on an integer shape panic on a non-finite input instead of storing the platform-dependent value Go's conversion gives; a float `T` keeps it (**breaking**)
- `RegularPolygonOrientationAngle` returns the top angle for `n < 1` instead of dividing by `n`, which stored a `NaN` angle in a `FlatTop` polygon without vertices and left it unequal to itself
- `Point.Transform` and `Vector.Transform` convert the matrix with `Matrix.Float` once and name the coordinates instead of spelling out every conversion inline; `ParseSize` splits with `strings.Cut`
- `Rectangle` and `Circle.Contains` document that containment compares exactly, without the `Epsilon` that `Equal` applies; `Direction.Vector` documents that an integer diagonal keeps the lattice step `(1,1)` where `Unit` snaps to `(1,0)`, and `Circle.Anchor` that such an anchor can fall outside a small integer circle; `RegularPolygon.Vertices` documents that it starts from `Angle`, `Polygon` that the vertex count is not checked, `Rectangle.Shrink` and `ShrinkXY` that they clamp to zero, and the `ints` and `floats` aliases have doc comments
- `Divide` panics for a zero scale instead of returning the value unchanged, like the integer `/` operator and `Mod`, so `Point.Divide(0)`, `Vector.DivideXY(0, 2)`, `Size.Unscale(0)` and `Matrix.Unscale(0, 3)` panic instead of silently skipping the axis; `Matrix.Inverse` panics for a singular matrix instead of returning it, with `IsInvertible` as the check to run first (**breaking**)
- `Matrix.Scale`, `PreScale` and `Unscale` take `float64` factors like every other `Scale` in the package, so an integer matrix can be scaled by `0.5`; each scaled component is rounded, following `Multiply`. `ScaleMatrix` keeps its `T` factors like `TranslationMatrix` (**breaking** for a `Matrix[int]` passing typed `int` factors)
- `Triangle`, `Square` and `Hexagon` delegate to `RegularPolygonWithOrientation` instead of repeating its body; `Axis` methods return from each `switch` case directly, like `Direction`
- `Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, so a `float32` now gets a tolerance matched to its precision instead of the `float64` one. The comparison stays absolute, and therefore stays a subtract and a compare for hot paths
- `AssertNumber` and every `Assert*` helper compare an integer `T` exactly and a float `T` within `EpsilonRelative`, instead of within `Delta` for both. A tolerance only has meaning where rounding error can occur, so an integer assertion no longer accepts a neighbouring value, and a float assertion now holds at any magnitude – beyond the ~1e3 (`float32`) and ~1e10 (`float64`) limits of the absolute comparison
- The `Assert*` helpers take a `Testing` interface declared by `geom` instead of `assert.Testing`, so no third-party type appears in the public API. `*testing.T` satisfies it unchanged, and unlike `testing.TB` a test can still pass its own recorder
- `AssertPoint`, `AssertVector`, `AssertSize`, `AssertCircle`, `AssertLine`, `AssertRect`, `AssertPolygon`, `AssertVertices`, `AssertRegularPolygon`, and `AssertPadding` take the expected value as the type they assert on – `AssertPoint(t, p, Pt(1, 2))` instead of `AssertPoint(t, p, 1, 2)` – matching `AssertNumber` and `AssertMatrix`, which already read `(t, actual, expected)`. A caller no longer has to remember the coordinate order of a four- or six-argument helper, and the expected value can be built with the same constructor as the actual one (**breaking**)
- `Point.Transform` and `Vector.Transform` accept a `Matrix[M]` of any `Float` type instead of only `Matrix[float64]`, using a method type parameter (Go 1.27). Existing calls with a `Matrix[float64]` are unaffected; an integer matrix is storage and converts with `Matrix.Float` at the call, the same way an angle is `float64`
- `String` (and therefore every `String()` on a geometry type) formats by `T` rather than by value: a float `8.0` now prints as `8.00`, like `8.1`, instead of `8`. Previously a single value could mix both forms – `Pt(100, -34.0000115).String()` gave `(100,-34.00)` (**breaking**)
- `Circle.Contains`, `CollisionCircles`, and `CollisionRectangleCircle` include the boundary, the same closed convention `Rectangle.Contains` and `CollisionRectangles` already followed; touching shapes now collide (**breaking**)
- `CollisionRectangleCircle` tests the point of the rectangle closest to the circle center, so it shares the exact `Min`/`Max` bounds of the other rectangle predicates instead of rounding the half extents of an odd integer size
- `RegularPolygon.Equal`, `RegularPolygon.IsZero` and `AssertRegularPolygon` compare the angle as well, with `EqualAngle`, so two polygons that produce different vertices no longer compare equal, while a full turn, the sign of an angle, or the `0`/`2π` seam does not matter: a polygon equals itself after `Rotate(0)`, `Rotate(2π)` or `Rotate(-1e-9)` (**breaking**)
- `Vector.Resize` on the zero vector returns `(length,0)` instead of NaN, the same +X convention `Normalize` uses; `Vector.Less` is false for a non-positive length instead of squaring the sign away
- `Polygon.Center` returns the zero point for a polygon without vertices instead of dividing by zero, and `RegularPolygon.Vertices` returns no vertices for `N < 1` instead of panicking in `make`; `RegularPolygon.Empty` is true for any `N < 1`
- `DirectionFromAngle` returns `DirectionNone` for `NaN` and `±Inf` instead of `DirectionRight` or a platform-dependent direction, so `Vector.Direction` on a NaN vector is `DirectionNone`; it also normalizes the angle before rounding it to a step, so a very large angle no longer converts an out-of-range float to an integer with platform-dependent results
- `Rectangle` documents that it is closed and that an integer rectangle of width `w` spans `w+1` lattice columns from `Min` to `Max` inclusive, while `Rectangle()` yields the half-open pixel rectangle of exactly `w` pixels for an integer `T`, and rounds the two corners independently of `Size.Int` for a float one
- `Vector.Resize`, `Normalize` and `Direction` treat only the exact zero vector as directionless. Previously they used the tolerant `IsZero`, so any vector shorter than `Delta` (`1e-6`) or `Delta32` (`1e-4`) snapped to `(1,0)` or `DirectionNone` and lost its direction
- `RegularPolygonOrientationAngle` returns `3π/2` instead of `-π/2` for `PointyTop`, the same normalized form `Rotate` stores
- `Axis.ScaleAlong` on `AxisNone` returns the size unchanged instead of a zero size
- `Size.Grow`, `Size.GrowXY`, `Rectangle.Grow`, `Rectangle.GrowXY`, and `Circle.Grow` clamp to zero for a negative amount, the same way `Shrink` already did, so no method can produce the negative size `Size` documents as unsupported (**breaking**)
- `Mod` documents that it panics for `m == 0`, like the `%` operator
- `ParseSize` wraps the underlying parse error with `%w`, so `errors.Is(err, strconv.ErrRange)` and `ErrSyntax` work
- `Integer` and `Matrix` document that `int8` and `int16` are storage-only: products such as `LengthSquared`, `Less`, `Circle.Contains`, `Size.Area` and `Matrix.Multiply` compute in `T` and overflow there. `Size` documents that a negative width or height is unsupported, and `Axis.IsNone` that every value outside the two axes counts as none
- `Direction.Angle` returns `NaN` for `DirectionNone` instead of `0`, which was indistinguishable from `DirectionRight`; `DirectionFromAngle(NaN)` is `DirectionNone`, so the two round-trip for every direction (**breaking**)
- `Axis.Size` stores the values as given, like `Sz`, instead of routing through `Vector.Size` and taking their absolute value; `Axis.ScaleAlong` with a negative factor now flips the sign the way `Size.ScaleXY` does
- `Polygon.MoveTo` documents that on an integer polygon the moved centroid can miss the point by one unit when a coordinate sum crosses zero, since the truncating `Center` does not commute with translation
- `Direction.Vector` documents that an integer diagonal is rounded and only approximates the requested length: `DirectionDownRight.Vector(5)` is `(4,4)`
- `Polygon` documents that `Pol` shares the vertices slice it is given, `Parse` that the `NaN` and `Inf` literals parse into a float `T`, and every `Direction` constant has a doc comment; the package doc and README state the integer rounding rule and its `Rectangle` exception, and the README notes that the flat JSON shape needs the v2-backed `encoding/json`, the default since Go 1.27

### Fixed
- `Point.Int`, `Vector.Int`, `Size.Int`, `Circle.Int`, `Padding.Int`, `Matrix.Int` and the shapes built on them convert an integer `T` directly instead of through `float64`, so an `int64` beyond 2^53 stays exact, like `Abs`, `Round`, `Floor` and `Ceil` already did
- `Polygon.Translate`, `Scale`, `ScaleXY`, `Int` and `Float` keep a nil `Vertices` nil, so `IsZero` survives every mapping instead of flipping to false on the first one (requires `gravitton/x` v1.2.1, where `slices.Map` maps nil to nil)
- `String` prints a float negative zero as `0.00` instead of `-0.00`, which `Matrix.Inverse` produced for every zero component it negated
- `Circle.Bounds` returns a square of side `Diameter` instead of `Radius`, so the rectangle actually bounds the circle. Previously it was half the required size and clipped the circle at every anchor, and `CollisionRectangleCircle(c.Bounds(), c)` could miss (**breaking**)
- `Vector.Normalize` snaps an integer vector to the longer axis and keeps its sign, so the result is always one of the four axis-aligned unit vectors, as documented. Previously it rounded each component of the resized vector and only corrected the result when both landed on the same value, so `Vec(-10, -16)` gave `(1,0)` – pointing the opposite way – and a mixed-sign vector like `Vec(-10, 16)` kept a diagonal `(-1,1)` of length √2. A tie between equal magnitudes now resolves to the X axis (**breaking**)
- `Direction.Unit` follows the same rule, so an integer diagonal collapses onto an axis – `DirectionUpRight.Unit[int]()` is `(1,0)` instead of the `(1,-1)` it returned before, which had length √2 rather than 1. `Direction.Offset` still gives the lattice step `(±1,±1)` for callers that need the diagonal (**breaking**)
- `Rectangle.Inset` moves the center down for a larger top padding instead of up: it used `Bottom-Top` where the X axis used `Left-Right`, so each vertical edge moved by the other edge's padding
- `Axis.Project` keeps the sign of the component, so projecting `(-3, 2)` onto the horizontal axis gives `-3` instead of `3`
- `EqualDelta` and `Lerp` take the difference in `float64`, so a narrow integer `T` no longer wraps: `EqualDelta[int8](127, -128, 1)` was true and `Lerp[int8](-100, 100, 0.5)` gave `-128`
- `Abs` stays in `T` and `Round`, `Floor`, `Ceil` return an integer `T` unchanged, so an `int64` beyond 2^53 is no longer rounded through `float64`
- `Circle.Area` and `Matrix.Inverse` square and cross-multiply in `float64`, so a large integer radius or translation no longer overflows before the conversion
- `Matrix.Unscale` divides each column by its own factor, following `Divide`, instead of multiplying by a rounded inverse scale. A zero factor leaves that axis unchanged: previously the guard only triggered when *both* factors were zero, so `ScaleMatrix(2, 3).Unscale(0, 3)` produced `[[+Inf, 0, 0], [NaN, 1, 0]]`. An integer matrix now unscales exactly when its components divide: `ScaleMatrix(4, 6).Unscale(2, 3)` is `ScaleMatrix(2, 2)` instead of the singular `ScaleMatrix(4, 0)` the rounded inverse `1/3 → 0` gave (**breaking** for integer `T`)
- `Parse` rejected defined types over an integer or float (`Parse[Direction]("3")`, `ParseSize[Coord]`) with "unsupported number type". It now parses in `int64`/`float64` and narrows to `T`, checking the result is in range, so every type the `~` constraints admit parses like its underlying type. A `strconv` error now returns the zero value instead of the out-of-range value strconv reports alongside it
- `ParseSize` documented that "float values are rounded to the nearest integer via Cast" for an integer `T`; it has always rejected them – `ParseSize[int]("23.5x12.4")` is an error
- `isIntType` reported `false` for defined types over an integer (`type Coord int`), so `Cast[Coord]` truncated instead of rounding. It now tests integer division directly, which works for every type the `~` constraints admit
- `Matrix`, `RotationMatrix`, `Rotate`, `PreRotate`, and `Inverse` document what an integer `T` cannot represent: only quarter turns are rotations and only `|det| = 1` inverts. An integer matrix is for storing and composing lattice transforms
- `Rectangle.Inset` on an integer rectangle with an odd asymmetric padding no longer leaks outside the original: the center shift rounded `0.5` up while `Min` truncated, so `Rect(Pt(0,0), Sz(10,10)).Inset(Pad(0,0,0,1))` spanned `-3..6` instead of `-4..5`. The inset is now derived from the padded `Min` corner
- `RegularPolygonOrientationAngle` with `FlatTop` placed a flat edge at the *bottom* (`π/2 - π/n`), which only coincides with a flat top for an even `n`. A flat-top triangle or pentagon therefore had a vertex at the top and was indistinguishable from `PointyTop`. It now puts an edge midpoint at the top for every `n`: `3π/2 - π/n`, half a step before the `PointyTop` vertex. For an even `n` the shape is unchanged but the first vertex moves by 180°, so `Vertices()` starts on the opposite side and `Equal` against a hard-coded angle changes (**breaking**)
- `NormalizeAngle` could return exactly `2π` for a tiny negative angle, violating its `[0, 2π)` contract, because adding `2π` to `-1e-17` rounds to `2π`; it now returns `0` there, and `NaN` for `NaN` and `±Inf` input


## [v1.12.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.11.0...v1.12.0)
### Added
- `Point.Compare(point) int` – orders points by X and then by Y in the `cmp.Compare` convention, for `slices.SortFunc` and `slices.BinarySearchFunc`. The comparison is exact and does not apply the `Delta` tolerance `Point.Equal` uses, since a tolerant comparison is not transitive and would leave a sort no valid ordering


## [v1.11.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.10.0...v1.11.0)
### Changed
- `Direction` constants are renumbered to run by increasing angle — `Right`, `DownRight`, `Down`, `DownLeft`, `Left`, `UpLeft`, `Up`, `UpRight` — so a positive `Rotate` step turns the same way as a positive `Vector.Rotate` angle: counterclockwise in math coordinates, clockwise as drawn on a screen with Y pointing down. Previously the order ran the opposite way and `DirectionFromAngle` negated its input to compensate (**breaking**)
- `Directions`, `CardinalDirections`, and `DiagonalDirections` follow the new order (**breaking**)
- `DirectionFromAngle` no longer negates its argument, now that direction numbering follows the angle directly
- `Rectangle.Vertices` and `Rectangle.Edges` now wind by increasing angle, matching `Directions` and `RegularPolygon.Vertices` — clockwise as drawn on a screen with Y pointing down. `Rectangle.Polygon()` winds the same way as a result (**breaking**)
- `Rectangle.TopEdge`, `RightEdge`, `BottomEdge`, and `LeftEdge` are reversed so the edges still form a closed chain in the new winding (**breaking**)

### Fixed
- `Direction.Rotate` documentation had the screen sense inverted — it claimed positive steps appear clockwise while the old ordering made them appear counterclockwise
- `RegularPolygon.Vertices` documentation claimed counter-clockwise winding; it winds by increasing angle, which is clockwise as drawn
- `Direction.Angle`, `Vector.Rotate`, `VectorFromAngle`, `Rectangle.Vertices`, and `Rectangle.Edges` now state which convention their angles and winding use


## [v1.10.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.9.0...v1.10.0)
### Added
- `Integer` and `Float` constraints – split out of `Number`, which is now `Integer | Float`
- `Mod[T Integer](n, m T) T` – wraps `n` into `[0, m)`, correctly for negative `n`
- `Axis` – names one of the two coordinate axes, or `AxisNone`, with methods for writing orientation-agnostic code in terms of "along" and "across"
- `Direction` – one of the eight neighbor directions on a square lattice, or `DirectionNone`, ordered counterclockwise from `DirectionRight`
- `East`/`North`/`West`/`South` (and diagonals) and `Top`/`Bottom`/`Left`/`Right`/`TopLeft`/`TopRight`/`BottomLeft`/`BottomRight` – aliases for lattice code and edge/corner aliases for rectangle anchors
- `Directions`, `CardinalDirections`, `DiagonalDirections`, `Axes` – the direction and axis sets in counterclockwise order
- `DirectionFromAngle(angle float64) Direction` – returns the direction nearest to an angle
- `Vector.Direction() Direction` – returns the direction nearest to a vector, or `DirectionNone` for the zero vector
- `Rectangle.Anchor(direction Direction) Point[T]` – returns a corner for diagonals and an edge midpoint for cardinals, unifying `TopLeft`, `TopRight`, `BottomLeft`, and `BottomRight`
- `Rectangle.Top()`, `Rectangle.Bottom()`, `Rectangle.Left()`, `Rectangle.Right()` – edge-midpoint accessors completing the corner accessor family
- `Rectangle.TopEdge()`, `Rectangle.BottomEdge()`, `Rectangle.LeftEdge()`, `Rectangle.RightEdge()` – dedicated edge accessors, which `Edges()` now composes
- `Circle.Anchor(direction Direction) Point[T]` – returns the boundary point in the given direction from the center, or the center for `DirectionNone`
- `VectorFromAngleSize[T](angle float64, size Size[T]) Vector[T]` – returns the point at the given angle on an ellipse with the given semi-axes, `(width*cos, height*sin)`; `VectorFromAngle` is now the square-size case of it

### Changed
- Require Go 1.27
- `Direction[T](x T) T` renamed to `Sign[T](x T) T`, freeing the name for the new `Direction` type
- `Line.Direction() Vector[T]` renamed to `Line.Vector()`, since it returns the start-to-end vector rather than a `Direction`
- `Line.Reversed()` renamed to `Line.Reverse()`, matching the imperative naming of every other immutable method
- `DirectionFromAxes(up, down, left, right bool)` now returns `Direction` instead of `Vector[float64]` (**breaking**)
- `RectFromMin`, `RectFromMax`, `RectFromMinMax`, `RectFromSize`, `RectFromImage` renamed to `RectangleFrom…`, so every `From` constructor uses its full type name like `PointFromImage` and `VectorFromAngle` (**breaking**)
- `Assert*` helpers take `assert.Testing` instead of `*testing.T`, matching the interface `gravitton/assert` already uses; existing calls passing `*testing.T` are unaffected

### Removed
- `UpVector`, `DownVector`, `LeftVector`, `RightVector`, `UpLeftVector`, `UpRightVector`, `DownLeftVector`, `DownRightVector` – superseded by `DirectionX.Unit[T]()` (**breaking**)

### Fixed
- `Matrix.Int()` now rounds to the nearest integer via `Cast` instead of truncating toward zero, matching every other `Int()` conversion
- `RegularPolygon.Vertices` now rounds after scaling by the size instead of rounding a unit vector first, so integer vertices land within half a unit of the exact position (previously up to 1.5 units off at non-right angles)


## [v1.9.0 (2026-05-29)](https://github.com/gravitton/geometry/compare/v1.8.1...v1.9.0)
### Added
- `Parse[T Number](s string) (T, error)` – parses a string into any `Number` type using `strconv.ParseInt` (with correct bit size) for integer types and `strconv.ParseFloat` for float types
- `ParseSize[T Number](s string) (Size[T], error)` – parses a size string in `"WxH"` form into `Size[T]`


## [v1.8.1 (2026-05-22)](https://github.com/gravitton/geometry/compare/v1.8.0...v1.8.1)
### Fixed
- `RegularPolygonOrientationAngle` with `PointyTop` now returns `-π/2` instead of `π/2`, correctly placing the first vertex at the visual top in the library's +Y-down coordinate system


## [v1.8.0 (2026-05-12)](https://github.com/gravitton/geometry/compare/v1.7.0...v1.8.0)
### Added
- `DirectionFromAxes(up, down, left, right bool) Vector[float64]` – returns the unit direction vector for the given axis inputs; cancels opposite directions, normalizes diagonals
- `Line.Vertices() []Point[T]` – returns `[Start, End]` as a slice


## [v1.7.0 (2026-05-09)](https://github.com/gravitton/geometry/compare/v1.6.0...v1.7.0)
### Added
- `RectFromMax` – constructor for `Rectangle` from max (bottom-right) point and size
- `Direction[T](x T) T` – returns the sign of a number: `1`, `-1`, or `0`
- `Vector.Round() Vector[T]` – returns a new vector with each component rounded to the nearest integer
- `Vector.Floor() Vector[T]` – returns a new vector with each component rounded down
- `Vector.Ceil() Vector[T]` – returns a new vector with each component rounded up


## [v1.6.0 (2026-05-03)](https://github.com/gravitton/geometry/compare/v1.5.0...v1.6.0)
### Added
- `RegPolWithOrientation` – constructor for `RegularPolygon` with orientation

### Changed
- `Matrix[T Number]` — `Matrix` is now generic over the full `Number` constraint (was restricted to `float64`)


## [v1.5.0 (2026-04-25)](https://github.com/gravitton/geometry/compare/v1.4.0...v1.5.0)
### Added
- `Point.ChebyshevDistanceTo(point Point[T]) T` — returns the Chebyshev (chessboard) distance: `max(|dx|, |dy|)`
- `Point.OctileDistanceTo(point Point[T]) float64` — returns the Octile distance for grid pathfinding where diagonals cost √2: `max(dx, dy) + (√2-1) * min(dx, dy)`
- `Round[T]`, `Floor[T]`, `Ceil[T]` generic free functions in `math.go` (same pattern as `Abs[T]`)

### Changed
- `Point.Round`, `Point.Floor`, `Point.Ceil` now delegate to the new `Round`, `Floor`, `Ceil` free functions


## [v1.4.0 (2026-04-25)](https://github.com/gravitton/geometry/compare/v1.3.0...v1.4.0)
### Added
- `Point.ManhattanDistanceTo(point Point[T]) T` — returns the Manhattan (taxicab) distance between two points


## [v1.3.0 (2026-04-25)](https://github.com/gravitton/geometry/compare/v1.2.1...v1.3.0)
### Added
- `NormalizeAngle(angle float64) float64` utility function — normalizes any angle to `[0, 2π)`
- `GOEXPERIMENT=jsonv2` requirement documented in `README.md`; `json:",inline"` struct tags are intentional and valid under jsonv2

### Fixed
- `RegularPolygon.Bounds` now computes the exact bounding rectangle from actual vertices instead of always returning the full bounding circle
- `RegularPolygon.Rotate` normalizes the stored angle to `[0, 2π)` after each rotation to prevent floating-point drift
- `Size.AspectRatio` returns `0` when `Height` is zero (previously returned `+Inf`)

### Changed
- Methods that produce imprecise results for integer `T` are now documented with a note in their comments: `Vector.Normalize`, `Vector.Rotate`, `Vector.Resize`, `VectorFromAngle`, `Point.Transform`, `Vector.Transform`, `RegularPolygon.Vertices`, `Polygon.Center`
- `Rectangle.Max` comment clarifies the `w-w/2` formula used for correct pixel span with odd integer sizes


## [v1.2.0 (2026-04-19)](https://github.com/gravitton/geometry/compare/v1.1.1...v1.2.0)
### Added
- `Point.Abs()`, `Point.Round()`, `Point.Floor()`, `Point.Ceil()` methods
- `Size.Unscale(factor)`, `Size.UnscaleXY(factorX, factorY)`, `Size.AtLeast(size)`, `Size.AtMost(size)` methods
- `RegularPolygon.Translate(vector)` method
- `Pol[T]` shorthand constructor for `Polygon`
- `Sqrt3` and `OneOverSqrt2` constants
- `Polygon.UnmarshalJSON` support

### Changed
- `Shrink` and `ShrinkXY` now clamp to `0` — negative dimensions are no longer possible
- Renamed `PointTop` → `PointyTop` in `RegularPolygon` orientation constants (**breaking**)
- Updated minimum Go version to 1.26

### Fixed
- `PadXY` parameter order corrected
- Logic errors and typos in comments across multiple files


## [v1.1.1 (2025-10-27)](https://github.com/gravitton/geometry/compare/v1.1.0...v1.1.1)
### Fixed
- Move Assert helper methods to `geom` package


## [v1.1.0 (2025-10-27)](https://github.com/gravitton/geometry/compare/v1.0.0...v1.1.0)
### Updated
- Make Assert helper methods public


## v1.0.0 (2025-10-17)
### Added
- Added new generic geometry primitives:
  - `Point`
  - `Vector`
  - `Size`
  - `Line`
  - `Circle`
  - `Rectangle`
  - `Polygon`
  - `RegularPolygon`
  - `Matrix`
  - `Padding`
- Added type alias packages `ints` for `Int` and `floats` as `float64`
- Added collision detection functions

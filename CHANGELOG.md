# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.12.0...master)
### Added
- `Delta32` – the equality tolerance for `float32`, `1e-4`. A `float32` ulp reaches 6e-5 at magnitude 1e3, so `Delta` sat below one ulp across the coordinate range `float32` is normally chosen for and asked for bit-exact equality there
- `Epsilon[T]()` – the tolerance for `T`: zero for an integer `T`, `Delta32` for `float32`, `Delta` for `float64`. It is a property of the type, so the compiler folds it to a constant at each instantiation and `Equal` costs what it did before
- `EqualRelative(a, b)` and `EpsilonRelative(a, b)` – equality within a tolerance that scales with magnitude, for values far from zero where an absolute tolerance falls below one ulp. Measured at ~1.1 ns against ~0.6 ns for `Equal`, which is why it is a separate function rather than a change to `Equal`
- `AssertNumber(t, actual, expected)` – asserts a bare coordinate with the same integer/float rule as the shape helpers

### Changed
- `Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, so a `float32` now gets a tolerance matched to its precision instead of the `float64` one. The comparison stays absolute, and therefore stays a subtract and a compare for hot paths
- `AssertNumber` and every `Assert*` helper compare an integer `T` exactly and a float `T` within `EpsilonRelative`, instead of within `Delta` for both. A tolerance only has meaning where rounding error can occur, so an integer assertion no longer accepts a neighbouring value, and a float assertion now holds at any magnitude – beyond the ~1e3 (`float32`) and ~1e10 (`float64`) limits of the absolute comparison
- The `Assert*` helpers take a `Testing` interface declared by `geom` instead of `assert.Testing`, so no third-party type appears in the public API. `*testing.T` satisfies it unchanged, and unlike `testing.TB` a test can still pass its own recorder
- `AssertPoint`, `AssertVector`, `AssertSize`, `AssertCircle`, `AssertLine`, `AssertRect`, `AssertPolygon`, `AssertVertices`, `AssertRegularPolygon`, and `AssertPadding` take the expected value as the type they assert on – `AssertPoint(t, p, Pt(1, 2))` instead of `AssertPoint(t, p, 1, 2)` – matching `AssertNumber` and `AssertMatrix`, which already read `(t, actual, expected)`. A caller no longer has to remember the coordinate order of a four- or six-argument helper, and the expected value can be built with the same constructor as the actual one (**breaking**)
- `Point.Transform` and `Vector.Transform` accept a `Matrix[M]` of any `Float` type instead of only `Matrix[float64]`, using a method type parameter (Go 1.27). Existing calls with a `Matrix[float64]` are unaffected; an integer matrix is storage and converts with `Matrix.Float` at the call, the same way an angle is `float64`
- `String` (and therefore every `String()` on a geometry type) formats by `T` rather than by value: a float `8.0` now prints as `8.00`, like `8.1`, instead of `8`. Previously a single value could mix both forms – `Pt(100, -34.0000115).String()` gave `(100,-34.00)` (**breaking**)
- The test suite follows one shape across every file: each test is named `Test<Type>_<Method>` (or `Test<Func>`), ordered to match the source, and split into named `t.Run` cases; marshal and unmarshal are merged into a single `_JSON` test with a round-trip case, each type has a `_Properties` test covering its invariants over shared fixtures, and scalars are asserted with `AssertNumber` rather than `assert.Equal`/`assert.EqualDelta`

### Fixed
- `Circle.Bounds` returns a square of side `Diameter` instead of `Radius`, so the rectangle actually bounds the circle. Previously it was half the required size and clipped the circle at every anchor, and `CollisionRectangleCircle(c.Bounds(), c)` could miss (**breaking**)
- `Vector.Normalize` snaps an integer vector to the longer axis and keeps its sign, so the result is always one of the four axis-aligned unit vectors, as documented. Previously it rounded each component of the resized vector and only corrected the result when both landed on the same value, so `Vec(-10, -16)` gave `(1,0)` – pointing the opposite way – and a mixed-sign vector like `Vec(-10, 16)` kept a diagonal `(-1,1)` of length √2. A tie between equal magnitudes now resolves to the X axis (**breaking**)
- `Direction.Unit` follows the same rule, so an integer diagonal collapses onto an axis – `DirectionUpRight.Unit[int]()` is `(1,0)` instead of the `(1,-1)` it returned before, which had length √2 rather than 1. `Direction.Offset` still gives the lattice step `(±1,±1)` for callers that need the diagonal (**breaking**)
- `Matrix.Unscale` inverts each axis on its own, following `Divide`: a zero factor leaves that axis unchanged instead of the whole matrix. Previously the guard only triggered when *both* factors were zero, so `ScaleMatrix(2, 3).Unscale(0, 3)` divided by zero and produced `[[+Inf, 0, 0], [NaN, 1, 0]]`
- `Parse` rejected defined types over an integer or float (`Parse[Direction]("3")`, `ParseSize[Coord]`) with "unsupported number type". It now parses in `int64`/`float64` and narrows to `T`, checking the result is in range, so every type the `~` constraints admit parses like its underlying type. A `strconv` error now returns the zero value instead of the out-of-range value strconv reports alongside it
- `ParseSize` documented that "float values are rounded to the nearest integer via Cast" for an integer `T`; it has always rejected them – `ParseSize[int]("23.5x12.4")` is an error
- `isIntType` reported `false` for defined types over an integer (`type Coord int`), so `Cast[Coord]` truncated instead of rounding. It now tests integer division directly, which works for every type the `~` constraints admit
- `Matrix`, `RotationMatrix`, `Rotate`, `PreRotate`, `Inverse`, and `Unscale` document what an integer `T` cannot represent: only quarter turns are rotations, only `|det| = 1` inverts, and only a factor of ±1 unscales. An integer matrix is for storing and composing lattice transforms


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

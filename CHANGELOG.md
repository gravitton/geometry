# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.13.0...main)

The main change is the oriented rectangle: `Rectangle` gains an `Angle`, every rectangle method answers for a turned one, and rectangle against rectangle is decided on the edges like every other pair, so no boundary in the package is judged on a coordinate any more.

### Added
- `Rectangle.Angle` – the box of `Size` about `Center`, turned by the angle in the sense of `Vector.Rotate`. `Rect` and the corner constructors leave it zero; `Rotate` turns the rectangle about its center, normalized to `[0, 2π)`; `Lerp` turns along the shorter arc; `Equal`, `IsZero` and `AssertRectangle` compare it modulo a full turn. It marshals as `a` with `omitzero`, so the JSON and `String` of a rectangle that is not rotated are unchanged, and a rotated one prints as `Rect((x,y);WxH;a)`
- `Rectangle.Rotate(angle)` and `IsAligned` – the turn, and the exact test for a rectangle with no turn, which has exact corners and coincides with its `Bounds`
- `Rectangle.Canonical` and `RegularPolygon.Canonical` snap an angle within `Delta` of zero or of a full turn to exactly zero, the residue a chain of turns can leave; `Rotate` and `Lerp` never snap, since a turn that small still moves a far corner of a large shape by more than `Epsilon`
- `Orientation.String`, `MarshalText`, `UnmarshalText`, `ParseOrientation`, `Orientations` and `OrientationNone` with `IsNone`, like `Direction` and `Axis`

### Changed
- A turned rectangle keeps the names of its corners, edges, anchors and `Inset` paddings from the frame before the turn; `Min`, `Max`, `MinMax` and `Bounds` become the axis-aligned extent of its vertices, and `Clamp` snaps to the nearest turned edge. `Contains`, `DistanceTo` and `IntersectsCircle` walk the edges as `Polygon` does, so a rotated integer rectangle contains exactly what the polygon of its rounded corners contains, and `DistanceSquaredTo` returns `float64` like `Line.DistanceSquaredTo`, since the nearest point is no longer a lattice point (**breaking**)
- `Rectangle.Intersects` is decided on the edges, a corner of one within the other or an edge meeting an edge, with the extent test as a prefilter and an exact overlap of two aligned extents decided at once. Rectangles of the same angle intersect and unite in a rectangle of that angle, found in their shared frame; rectangles of different angles have no rectangle in common, so `Intersection` returns false for them, as `Line.Intersection` does for parallel segments, and `Union` gives the box around both `Bounds` (**breaking**)
- Every test against a radius or a segment is one squared-distance comparison.
- `Line.IntersectionCircle` never returns more than two points and counts points that compare `Equal` once; `Circle.Intersection` places a tangent point halfway between the two boundaries, within half the tolerance of both
- `Line.MoveTo` places the midpoint on the point, the center every other shape places with `MoveTo` and the pivot its `Scale`, `Resize` and `Rotate` turn about (**breaking**)
- `Rectangle.Rectangle` builds the `image.Rectangle` from `Int`, so it spans exactly `Size.Int` pixels and a sprite keeps its width at sub-pixel positions (**breaking**)
- `Line.MinMax` and `Polygon.MinMax` are unexported: they are the corners of `Bounds`, so `Bounds().MinMax()` gives the same pair at the same cost; `Rectangle` alone keeps `Min`, `Max` and `MinMax` public (**breaking**)
- `Size.Grow`, `GrowXY`, `Shrink` and `ShrinkXY` no longer clamp at zero, since a size is signed; `Rectangle.Grow` and `Shrink` clamp their own extent as before (**breaking**)
- `Polygon.String` separates the vertices with `;` like every other shape: `Pol((0,0);(2,0);(2,2))`; `RegularPolygon` omits a zero angle from its JSON and `String`, as `Rectangle` does: `RegPol((1,2);2x2;4)` (**breaking**)
- `Polygon.Edges` moved before `Vertices` and `Center`, in the method order every shape follows


## [v1.13.0 (2026-09-18)](https://github.com/gravitton/geometry/compare/v1.12.0...v1.13.0)

Every entry, with breaking changes marked, is in [docs/releases/v1.13.0.md](docs/releases/v1.13.0.md).

### Added
- `Intersects` between every pair of shapes, with `Intersection` for lines, rectangles and circles and the boundary crossings of a segment
- `Contains`, `DistanceTo` and `DistanceSquaredTo` on every shape
- `Polygon.Area`, `Perimeter`, `Contains`, `Edges`, `Bounds`, `Rotate` and `Transform`
- `Unscale`, `Lerp`, `Canonical` and `AlignTo` across the shapes; `Line.Scale`, `Resize`, `Rotate`, `Normal`; `Point.RotateAround`; `Vector.Project`, `Reflect`, `AtMost`; `Size.Fit`, `Fill`, `Transpose`
- `ShearMatrix`, `ReflectionMatrix`, and `Matrix.Angle`, `Scaling`, `IsInvertible`
- `Epsilon`, `EqualRelative`, `LessOrEqual`, `AngleDistance`, `LerpAngle`, `ParseDirection`, `ParseAxis`

### Changed
- Boundaries are closed and tolerant within `Epsilon[T]()`, and `DistanceTo` is zero exactly where `Contains` holds
- No shape stores a negative size or radius; the empty circle of a negative radius is gone
- Arithmetic computes in `float64` and rounds back into `T`, so narrow integers never overflow mid-computation
- `Cast`, `Divide`, `Unscale` and `Matrix.Inverse` panic on a non-finite, zero or singular input
- `Polygon.Center` is the area centroid; `RegularPolygon` measures without building its vertices
- `Collision*` functions replaced by methods, `AssertRect` renamed, `Line` JSON keys `s`/`e`, `Direction` and `Axis` marshal as names

### Fixed
- `Circle.Bounds`, `Rectangle.Inset`, `Axis.Project`, `Vector.Normalize`, `NormalizeAngle`, `Matrix.Unscale`, `Polygon.UnmarshalJSON` and `Cast` on defined integer types


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

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.8.0...master)


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

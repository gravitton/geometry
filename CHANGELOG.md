# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.4.0...master)


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

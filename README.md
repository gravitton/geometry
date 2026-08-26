# Geometry

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Generic immutable 2D geometry library for game development.

> Uses a top-left origin with +Y down. This only affects directional getters (`Top`, `Bottom`, `Up`, `Down`) and rotation. Angles follow the standard mathematical convention, so a positive angle or step is counterclockwise in math coordinates and appears clockwise on screen; a negative one appears counterclockwise on screen. Direction order and polygon winding follow the same rule.

## Installation

```bash
go get github.com/gravitton/geometry
```

## Usage

```go
import (
	geom "github.com/gravitton/geometry"
)

p1 := geom.Pt(1, 2)
p2 := geom.Pt(3, 4)

v := p2.Subtract(p1) // Vector
if v.Equal(geom.Vec(2, 2)) { ... }

r := geom.Rect(p1, geom.Sz(10, 5))
if r.Contains(p2) { ... }

hex := geom.Hexagon(p1, geom.SzU(20), geom.FlatTop)
for _, vertex := range hex.Vertices() { ... }
```

Matrix transforms:

```go
m := geom.IdentityMatrix[float64]().Rotate(math.Pi / 4).Scale(2, 2)
p := geom.Pt(1.0, 0.0).Transform(m)
```

Directions:

```go
dir := geom.DirectionDownRight
v := dir.Vector(5.0)

p := geom.RectangleFromMinMax(p1, p2).Anchor(geom.Bottom)

axis := geom.AxisVertical
length := axis.Along(size)          // Height, because the axis is vertical
bar := axis.Size(length, thickness) // Size{thickness, length}
```

Type aliases for common numeric types ([`ints`](./types/ints/types.go), [`floats`](./types/floats/types.go)):

```go
import (
    "github.com/gravitton/geometry/types/floats"
    "github.com/gravitton/geometry/types/ints"
)

type Grid struct {
    Size     ints.Size
    CellSize floats.Size
}
```

## API

All types are generic over the `Number` constraint:

```go
type Integer interface { ~int | ~int8 | ~int16 | ~int32 | ~int64 }
type Float   interface { ~float32 | ~float64 }
type Number  interface { Integer | Float }
```

Full documentation is available at [pkg.go.dev/github.com/gravitton/geometry][link-go-dev-reference].

### Types

| Type                | Constructor                                                                                              | Description                            |
|---------------------|----------------------------------------------------------------------------------------------------------|----------------------------------------|
| `Point[T]`          | `Pt(x, y)`                                                                                               | 2D position                            |
| `Vector[T]`         | `Vec(x, y)`, `VectorFromAngle`, `VectorFromAngleSize`                                                    | 2D displacement                        |
| `Size[T]`           | `Sz(w, h)`, `SzU(n)`, `ParseSize("WxH")`                                                                 | Width and height                       |
| `Rectangle[T]`      | `Rect(center, size)`, `RectangleFromMin`, `RectangleFromMax`, `RectangleFromMinMax`, `RectangleFromSize` | Axis-aligned rectangle (center + size) |
| `Circle[T]`         | `Circ(center, r)`                                                                                        | Circle (center + radius)               |
| `Line[T]`           | `Ln(start, end)`                                                                                         | Line segment                           |
| `Polygon[T]`        | `Pol(vertices)`                                                                                          | Arbitrary polygon                      |
| `RegularPolygon[T]` | `RegPol(center, size, n, angle)`, `Triangle`, `Square`, `Hexagon`                                        | Regular polygon                        |
| `Padding[T]`        | `Pad(t,r,b,l)`, `PadU(n)`, `PadXY(tb, lr)`                                                               | Top/Right/Bottom/Left padding          |
| `Matrix[T]`         | `Mat(a,b,c,d,e,f)`, `IdentityMatrix`, `TranslationMatrix`, `RotationMatrix`, `ScaleMatrix`               | 2D affine matrix                       |

### Direction

`Direction` names one of the eight neighbor directions on a square lattice, or `DirectionNone`. Constants are ordered by
increasing angle from `DirectionRight`, so `Rotate(1)` and a `Vector.Rotate` of `+π/4` turn the same way. Any other
integer wraps into range, so `Rotate` accepts any step count.

Three names exist for each direction — the canonical one, a compass alias for lattice and map code, and an edge/corner
alias for rectangle anchors:

| Index | Angle   | Canonical            | Compass     | Rectangle     |
|-------|---------|----------------------|-------------|---------------|
| 0     | `0°`    | `DirectionRight`     | `East`      | `Right`       |
| 1     | `45°`   | `DirectionDownRight` | `SouthEast` | `BottomRight` |
| 2     | `90°`   | `DirectionDown`      | `South`     | `Bottom`      |
| 3     | `135°`  | `DirectionDownLeft`  | `SouthWest` | `BottomLeft`  |
| 4     | `180°`  | `DirectionLeft`      | `West`      | `Left`        |
| 5     | `-135°` | `DirectionUpLeft`    | `NorthWest` | `TopLeft`     |
| 6     | `-90°`  | `DirectionUp`        | `North`     | `Top`         |
| 7     | `-45°`  | `DirectionUpRight`   | `NorthEast` | `TopRight`    |

`Angle()` reports the angle in radians over `atan2`'s `(-π, π]` range, which is why the last three rows are negative
rather than `225°`–`315°`. `DirectionFromAngle` is its inverse and the two round-trip for every direction.

### Axis

`Axis` is `AxisHorizontal`, `AxisVertical`, or `AxisNone`. It exists to write orientation-agnostic code in terms of "along" and "across".

### Conventions

All methods return new values — no mutation.

Every shape type exposes `.Int()`, `.Float()`, and implements `String()`, `Equal()`, `IsZero()`. 

The `Direction` and `Axis` enums implement `String()` and `IsNone()` instead.

Types with spatial extent also implement `Bounds() Rectangle[T]`.

`Line`, `Polygon`, and `RegularPolygon` also expose `Vertices() []Point[T]`.

`Rectangle.Vertices`, `Rectangle.Edges`, and `RegularPolygon.Vertices` all wind by increasing angle, the same order as
`Directions` — clockwise as drawn.

### Collision

```go
CollisionRectangles[T](a, b Rectangle[T]) bool
CollisionCircles[T](a, b Circle[T]) bool
CollisionRectangleCircle[T](r Rectangle[T], c Circle[T]) bool
```

### Image interop

```go
PointFromImage[T](p image.Point) Point[T]
SizeFromImage[T](r image.Rectangle) Size[T]
RectangleFromImage[T](r image.Rectangle) Rectangle[T]

(Point[T]).Point() image.Point
(Rectangle[T]).Rectangle() image.Rectangle
```

### Math utilities

```go
Lerp[T](a, b T, t float64) T
Clamp[T](v, min, max T) T
Equal[T](a, b T) bool // within Delta (1e-6)
EqualDelta[T](a, b T, d float64) bool
Midpoint[T](a, b T) T
Abs[T](a T) T
Round[T](a T) T
Floor[T](a T) T
Ceil[T](a T) T
Sign[T](x T) T // 1, -1, or 0
Mod[T Integer](n, m T) T // wraps into [0, m), correct for negative n
Multiply[T](a T, factor float64) T
Divide[T](a T, factor float64) T
Parse[T](s string) (T, error)
ToRadians(deg float64) float64
ToDegrees(rad float64) float64
NormalizeAngle(angle float64) float64 // into [0, 2π)
```

## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]

## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/geometry.svg?style=flat-square&colorB=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/geometry/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/geometry?style=flat-square&colorB=blue
[ico-go-dev-reference]:     https://img.shields.io/badge/go.dev-reference-blue?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/geometry?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/geometry/releases
[link-contributors]:        https://github.com/gravitton/geometry/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/geometry/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/geometry
[link-coverage]:            https://coveralls.io/github/gravitton/geometry

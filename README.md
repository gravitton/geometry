# Geometry

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Generic immutable 2D geometry library for game development.

> Uses a top-left origin with +Y down. This only affects directional getters (`Top`, `Bottom`, `Up`, `Down`).

## Installation

```bash
go get github.com/gravitton/geometry
```

## Requirements

Requires **`GOEXPERIMENT=jsonv2`** (Go 1.26+). This enables the JSON v2 experiment, which supports the `json:",inline"` struct tag used by several types to produce flat JSON objects instead of nested ones.

Set the environment variable when running or testing:

```bash
GOEXPERIMENT=jsonv2 go test ./...
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
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}
```

Full documentation is available at [pkg.go.dev/github.com/gravitton/geometry][link-go-dev-reference].

### Types

| Type | Constructor | Description |
|---|---|---|
| `Point[T]` | `Pt(x, y)` | 2D position |
| `Vector[T]` | `Vec(x, y)` | 2D displacement |
| `Size[T]` | `Sz(w, h)`, `SzU(n)` | Width and height |
| `Rectangle[T]` | `Rect(center, size)`, `RectFromMin`, `RectFromMax`, `RectFromMinMax`, `RectFromSize` | Axis-aligned rectangle (center + size) |
| `Circle[T]` | `Circ(center, r)` | Circle (center + radius) |
| `Line[T]` | `Ln(start, end)` | Line segment |
| `Polygon[T]` | `Pol(vertices)` | Arbitrary polygon |
| `RegularPolygon[T]` | `RegPol(center, size, n, angle)`, `Triangle`, `Square`, `Hexagon` | Regular polygon |
| `Padding[T]` | `Pad(t,r,b,l)`, `PadU(n)`, `PadXY(tb, lr)` | Top/Right/Bottom/Left padding |
| `Matrix[T]` | `Mat(a,b,c,d,e,f)`, `IdentityMatrix[T]()`, `TranslationMatrix`, `RotationMatrix[T]`, `ScaleMatrix` | 2D affine matrix |

### Conventions

All methods return new values — no mutation.

Every type exposes `.Int()`, `.Float()`, and implements `String()`, `Equal()`, `IsZero()`.

Types with spatial extent also implement `Bounds() Rectangle[T]`.

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
RectFromImage[T](r image.Rectangle) Rectangle[T]

(Point[T]).Point() image.Point
(Rectangle[T]).Rectangle() image.Rectangle
```

### Math utilities

```go
Lerp[T](a, b T, t float64) T
Clamp[T](v, min, max T) T
Equal[T](a, b T) bool          // within Delta (1e-6)
EqualDelta[T](a, b T, d float64) bool
Midpoint[T](a, b T) T
Abs[T](a T) T
Round[T](a T) T
Floor[T](a T) T
Ceil[T](a T) T
Direction[T](x T) T
Multiply[T](a T, factor float64) T
Divide[T](a T, factor float64) T
ToRadians(deg float64) float64
ToDegrees(rad float64) float64
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
[ico-go-report-card]:       https://goreportcard.com/badge/github.com/gravitton/geometry?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/geometry?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/geometry/releases
[link-contributors]:        https://github.com/gravitton/geometry/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/geometry/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/geometry
[link-go-report-card]:      https://goreportcard.com/report/github.com/gravitton/geometry
[link-coverage]:            https://coveralls.io/github/gravitton/geometry

<div align="center" width="100%">

<a href="https://github.com/gravitton">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/logo-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/logo-light.svg">
  <img alt="Gravitton geometry" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/logo-light.svg" width="300">
</picture>
</a>

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Generic, immutable 2D geometry library for game development

<hr>

</div>


## Features

- **Generic** over all integer and float types, named ones included.
- **Immutable** – every method returns a new value.
- **Shapes** – point, vector, size, rectangle, circle, line, polygon, regular polygon, padding, affine matrix.
- **Directions and axes** as enums, with compass and rectangle-anchor aliases.
- **Screen space** – top-left origin, `+Y` down, one winding order everywhere.
- **Extras** – `image` interop, JSON tags, string parsing, numeric helpers, test assertions.

## Installation

```shell
go get github.com/gravitton/geometry
```

## Usage

```go
import geom "github.com/gravitton/geometry"
```

Points and vectors:

```go
p := geom.Pt(1, 2)
v := geom.Pt(4, 6).Subtract(p) // Vector{3, 4}

v.Length()          // 5
v.Normal()          // Vector{-4, 3}, perpendicular
p.Add(v.Resize(10)) // Point{7, 10}

v.Project(geom.Vec(1, 0))       // Vector{3, 0}, the component along X
v.Reflect(geom.Vec(0, 1))       // Vector{3, -4}, bounced off a horizontal wall
v.AngleBetween(geom.Vec(1, 0))  // 0.93, unsigned

p.RotateAround(geom.Pt(0, 0), math.Pi/2)     // Point{-2, 1}
p.Lerp(geom.Pt(9, 10), 0.25)                 // Point{3, 4}
p.ManhattanDistanceTo(geom.Pt(4, 5))         // 6, for grid pathfinding
geom.Pt(0.0, 0.0).AngleTo(geom.Pt(1.0, 1.0)) // π/4
```

Sizes and padding:

```go
s := geom.Sz(1920, 1080)
s.Scale(0.5)                // Size{960, 540}
s.AtMost(geom.SzU(800))     // Size{800, 800}, clamped per axis
s.Fit(geom.SzU(800))        // Size{800, 450}, largest with the same ratio inside
s.Fill(geom.SzU(800))       // Size{1422, 800}, smallest with the same ratio around
s.AspectRatio()             // 1.78

geom.PadXY(4, 8).Size()     // Size{16, 8}, horizontal and vertical total
```

Rectangles:

```go
r := geom.Rect(geom.Pt(50, 50), geom.Sz(20, 10)) // center + size

r.Contains(geom.Pt(55, 52))                 // true
r.Inset(geom.PadU(2)).Anchor(geom.TopRight) // Point{58, 47}
r.Outset(geom.PadXY(1, 2))                  // Rectangle (38,44)-(62,56) by MinMaxString
r.Clamp(geom.Pt(80, 0))                     // Point{60, 45}, nearest point inside
r.AlignTo(geom.TopLeft, geom.Pt(0, 0))      // Rectangle (0,0)-(20,10) by MinMaxString

b := geom.RectangleFromMinMax(geom.Pt(0, 0), geom.Pt(8, 6))
b.Scale(2)     // Rectangle (-4,-3)-(12,9) by MinMaxString, scaled around the center
b.Edges()[0]   // Line (0,0)-(8,0), the top edge
b.Vertices()   // clockwise from the top-left corner
```

Circles, lines, polygons:

```go
c := geom.Circ(geom.Pt(0.0, 0.0), 5.0)
c.Anchor(geom.Bottom) // Point{0, 5}

l := geom.Ln(geom.Pt(0, 0), geom.Pt(3, 4))
l.Length()                  // 5
l.DistanceTo(geom.Pt(3, 0)) // 2.4, to the nearest point of the segment
l.Contains(geom.Pt(6, 8))   // false, the segment ends at (3,4)

p := geom.Pol([]geom.Point[int]{{0, 0}, {4, 0}, {4, 4}, {2, 1}, {0, 4}})
p.Area()                  // 10, by the shoelace formula
p.Perimeter()             // 12 + 2√13
p.Contains(geom.Pt(2, 3)) // false, inside the notch
p.Edges()[4]              // Line (0,4)-(0,0), closing back to the first vertex

hex := geom.Hexagon(geom.Pt(0, 0), geom.SzU(20), geom.FlatTop)
hex.Bounds() // Rectangle (-20,-17)-(20,17) by MinMaxString
for _, vertex := range hex.Vertices() { ... }
```

Intersections:

```go
a.Intersects(b)        // Rectangle, Circle, Line and Polygon with their own kind
r.IntersectsCircle(c)  // and IntersectsLine, IntersectsPolygon
c.IntersectsRectangle(r)
l.IntersectsPolygon(p)

point, ok := l.Intersection(m)  // where two segments cross
c.Intersection(d)               // zero, one or two points where two circles cross
box, ok := a.Intersection(b)    // the overlap of two rectangles
a.Union(b)                      // the smallest rectangle around both
```

Directions and axes:

```go
dir := geom.DirectionUp // marshals as "Up"; ParseDirection reads it back
dir.Rotate(2)   // DirectionRight, two 45° steps
dir.Vector(5.0) // Vector{0, -5}

geom.DirectionFromAxes(up, down, left, right) // keyboard input to an 8-way direction
geom.Vec(3, -7).Direction()                   // DirectionUpRight, nearest of the eight

axis := geom.AxisVertical
axis.Along(size)             // Height, because the axis is vertical
axis.Size(length, thickness) // Size{thickness, length}
```

Matrix transforms:

```go
m := geom.IdentityMatrix[float64]().Rotate(math.Pi / 4).Scale(2, 2)

geom.Pt(1.0, 0.0).Transform(m) // Point{1.41, 1.41}
m.Angle()                      // π/4, read back from the matrix
m.Scaling()                    // Vector{2, 2}
m.Translation()                // Vector{0, 0}

geom.ShearMatrix(0.5, 0.0)                    // x' = x + 0.5y
geom.ReflectionMatrix[float64](geom.AxisHorizontal) // flips Y
```

Interop with `image`:

```go
geom.RectangleFromImage[int](img.Bounds())
geom.Pt(3, 4).Point() // image.Point
geom.RectangleFromMin(geom.Pt(0, 0), geom.Sz(4, 2)).Rectangle() // image.Rectangle
```

Any type satisfying `Number` works, including named types. `Int()` and `Float()` convert between them:

```go
type Tile int32

geom.Pt[Tile](3, 4).Add(geom.DirectionRight.Unit[Tile]()) // Point[Tile]{4, 4}
geom.Circ(geom.Pt[float32](0, 0), 5).Contains(geom.Pt[float32](3, 4))

geom.Pt(1.4, 2.6).Int()     // Point[int]{1, 3}, rounded
geom.Sz(4, 2).Float()       // Size[float64]{4, 2}
geom.Vec(3, 0).Transform(m) // an int vector through a float matrix, rounded
```

The `ints` and `floats` packages alias the two common instantiations and add constructors that round or widen
from any `Number`:

```go
import (
	"github.com/gravitton/geometry/types/floats"
	"github.com/gravitton/geometry/types/ints"
)

type Grid struct {
	Size     ints.Size   // geom.Size[int]
	CellSize floats.Size // geom.Size[float64]
}

ints.Pt(1.4, 2.6)   // Point[int]{1, 3}
floats.Sz(4, 2)     // Size[float64]{4, 2}
ints.IdentityMatrix()
```

JSON and string parsing:

```go
json.Marshal(geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4))) // {"x":1,"y":2,"w":3,"h":4}
json.Marshal(geom.Circ(geom.Pt(1, 2), 3))             // {"x":1,"y":2,"r":3}

geom.ParseSize[int]("4x2") // Size{4, 2}
```

The flat shape of `Rectangle`, `Circle`, and `RegularPolygon` relies on the `embed` struct tag of the v2-backed
`encoding/json`, the default since Go 1.27; under `GOEXPERIMENT=nojsonv2` the nested fields are emitted as objects.

Test assertions, one per shape, comparing with the tolerance of the asserted type:

```go
geom.AssertPoint(t, got, geom.Pt(1.0, 2.0))
geom.AssertRect(t, got, want, "after inset")
geom.AssertVertices(t, hex.Vertices(), want)
```

Full reference: [pkg.go.dev][link-go-dev-reference].

## Conventions

**Screen space:** The origin is top-left and `+Y` points down. Only the directional getters (`Top`, `Bottom`, `Up`,
`Down`) and rotation depend on it.

**Angles:** Angles follow the mathematical convention. A positive angle or step is counterclockwise in math
coordinates, which appears clockwise on screen. Direction order and polygon winding follow the same rule:
`Rectangle.Vertices`, `Rectangle.Edges`, and `RegularPolygon.Vertices` wind by increasing angle, clockwise as drawn.

**Directions:** `Direction` covers the eight neighbours on a square lattice plus `DirectionNone`. They are ordered by
increasing angle from `DirectionRight`. Each one has three names: canonical (`DirectionDownRight`), compass
(`SouthEast`), and rectangle corner (`BottomRight`). `Angle()` reports in `atan2`'s `(-π, π]` range and
`DirectionFromAngle` inverts it.

**Matrices:** An integer `Matrix` composes lattice transforms exactly: translation, integer scale, reflection, quarter
turns. Rotation by any other angle, `Inverse`, and `Unscale` are not closed over the integers and round into a
different matrix. Use a float `Matrix` there. `Transform` takes a `float32` or `float64` matrix, so convert an integer
one with `Float()` at the call, the same way an angle is always `float64`.

**Equality:** `Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, a tolerance matched to
`float32` or `float64`. `EqualRelative` scales that tolerance for values far from zero. `EqualAngle` compares angles
modulo a full turn, across the `0`/`2π` seam, and is what `RegularPolygon.Equal` uses.

**Integer rounding:** A float result stored into an integer `T` rounds half away from zero, so
`Pt(0, 0).Midpoint(Pt(5, 5))` is `(3,3)` and `Polygon.Center` rounds the vertex average. `Rectangle` is the exception:
its center truncates half the size toward `Min` so that `Max-Min` stays exactly the size, so
`RectangleFromMinMax(Pt(0, 0), Pt(5, 5)).Center` is `(2,2)`. A finite value outside the range of an integer `T` is not
checked and stores a platform-dependent value.

**Division by zero:** `Divide`, `Unscale`, and `Matrix.Inverse` on a singular matrix panic, like the integer `/`
operator and `Mod`. Check `IsInvertible` first when a matrix may be singular.

**Arithmetic:** Products, distances and interpolations are computed in `float64` and rounded back into `T`, so a
narrow `int8` or `int16` never overflows mid-computation and only a result outside its range is lost. An `int64`
beyond 2^53 loses precision on the way through `float64`.

**Boundaries:** every `Contains`, `Intersects` and `Vector.LessOrEqual` are closed and tolerant: a point
within `Epsilon[T]()` of the boundary counts as on it, so a float rectangle contains the corners it was built from, a
circle contains its anchors and a polygon contains its vertices. `Vector.Less` and `LessOrEqual`'s strict counterparts apply no tolerance.
`DistanceTo` is zero exactly where `Contains` holds, and `Intersection` answers exactly where `Intersects` holds, less
the parallel and coincident cases that have no single answer.

**Layout:** every shape file lists its methods in the same order, and the tests follow it: constructors; properties
read from the value alone (`XY`, `Length`, `Width`, `Min`, `Area`, `Bounds`); arithmetic on the components
(`Add`, `Multiply`, `Lerp`, `Translate`, `Scale`, `Grow`, `Inset`); geometry treating the value as a position or
direction (`Transform`, `Rotate`, `Normalize`, `Project`, `AlignTo`); relations with other geometry (`Dot`,
`Contains`, `DistanceTo`, `Intersects`, `Intersection`, `Union`); equality and state (`Equal`, `IsZero`, `Is*`);
conversions leaving the type (`Vector`, `Polygon`, `Int`, `Float`); and `String` with JSON last.

**Common API:** Every shape exposes `Int()`, `Float()`, `String()`, `Equal()`, and `IsZero()`. Shapes with spatial
extent add `Bounds()`, `Contains(point)`, `DistanceTo(point)`, `DistanceSquaredTo(point)` and an `Intersects` method for every
other shape.
`Line`, `Polygon`, and `RegularPolygon` add `Vertices()`; `Rectangle` and `Polygon` add `Edges()`, `Area()` and
`Perimeter()`. `Point`, `Vector`, `Line` and `Polygon` take a float `Matrix` in `Transform()`.

## Planned

- **`Rectangle.Angle`** – an oriented rectangle. `Contains`, `Clamp`, `Intersects`, `Intersection` and `Union` assume
  axis alignment through `Min` and `Max` today; the edge-based tests already work for any orientation.
- **`Polygon.Winding`, `IsConvex` and `ConvexHull`** – the vertex order as a sign, since `Center` already computes
  the signed area that `Area` discards; a convexity test built on it; and the hull of a point set. Convexity also
  unlocks a separating-axis test for `Intersects`: two polygons apart but with overlapping bounds currently compare
  every edge pair, the slow case in `BenchmarkPolygon_Intersects`.
- **`Encloses`** – shape-in-shape containment such as `Rectangle.Encloses(rectangle)` and `Circle.Encloses(circle)`
  for culling, distinct from `Contains`, which takes a point.
- **`Ray`** – a half-line with origin and direction, for casts against every shape.
- **`Ellipse`** – `RegularPolygon` already takes semi-axes; the continuous shape has no type.

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

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

- **Generic** over every integer and float type, named types included.
- **Immutable** – every method returns a new value.
- **Shapes** – point, vector, size, padding, rectangle, circle, line, polygon, regular polygon, affine matrix.
- **Directions and axes** as enums, with compass and rectangle-anchor aliases.
- **Screen space** – top-left origin, `+Y` down, one winding order everywhere.
- **Extras** – `image` interop, JSON, string parsing, numeric helpers, test assertions.

## Installation

```shell
go get github.com/gravitton/geometry
```

## Usage

```go
import geom "github.com/gravitton/geometry"
```

### Points and vectors

```go
p := geom.Pt(1, 2)
v := geom.Pt(4, 6).Subtract(p) // Vector{3, 4}

v.Length()          // 5
v.Normal()          // Vector{-4, 3}, perpendicular
p.Add(v.Resize(10)) // Point{7, 10}
v.AtMost(2)         // Vector{1.2, 1.6}, a velocity capped at a top speed; AtLeast is the floor

v.Project(geom.Vec(1, 0))      // Vector{3, 0}, the component along X
v.Reflect(geom.Vec(0, 1))      // Vector{3, -4}, bounced off a horizontal wall
v.AngleBetween(geom.Vec(1, 0)) // 0.93, unsigned

p.RotateAround(geom.Pt(0, 0), math.Pi/2) // Point{-2, 1}
p.Lerp(geom.Pt(9, 10), 0.25)             // Point{3, 4}
p.ManhattanDistanceTo(geom.Pt(4, 5))     // 6, for grid pathfinding
```

### Sizes and padding

```go
s := geom.Sz(1920, 1080)
s.Scale(0.5)            // Size{960, 540}
s.Unscale(2)            // Size{960, 540}, divided rather than multiplied
s.Transpose()           // Size{1080, 1920}, turned a quarter turn
s.AtMost(geom.SzU(800)) // Size{800, 800}, clamped per axis
s.Fit(geom.SzU(800))    // Size{800, 450}, largest with the same ratio inside
s.Fill(geom.SzU(800))   // Size{1422, 800}, smallest with the same ratio around

geom.Sz(1.6, 2.4).Round()                 // Size{2, 2}, and Floor, Ceil, Lerp
geom.PadXY(4, 8).Size()                   // Size{16, 8}, horizontal and vertical total
geom.PadU(4).Subtract(geom.PadXY(1, 2))   // Padding{3, 2, 3, 2}
```

### Rectangles

```go
r := geom.Rect(geom.Pt(50, 50), geom.Sz(20, 10)) // center and size

r.Contains(geom.Pt(55, 52))                 // true
r.Clamp(geom.Pt(80, 0))                     // Point{60, 45}, nearest point inside
r.Inset(geom.PadU(2)).Anchor(geom.TopRight) // Point{58, 47}
r.AlignTo(geom.TopLeft, geom.Pt(0, 0))      // Rectangle (0,0)-(20,10) by MinMaxString

b := geom.RectangleFromMinMax(geom.Pt(0, 0), geom.Pt(8, 6))
b.Scale(2)         // Rectangle (-4,-3)-(12,9), scaled around the center
b.Unscale(2)       // Rectangle (2,2)-(6,5), the inverse; every shape has the pair
b.Canonical()      // the form Rect builds, repairing a literal or decoded negative extent
b.Lerp(c, 0.5)     // center and size together, for tweening
b.Edges()[0]       // Line (0,0)-(8,0), the top edge
b.Vertices()       // clockwise from the top-left corner
```

### Circles, lines and polygons

```go
c := geom.Circ(geom.Pt(0.0, 0.0), 5.0)
c.Anchor(geom.Bottom) // Point{0, 5}
c.Scale(-2)           // Circ((0,0);10), a negative factor scales by its absolute value, as Circ and Canonical take the radius
c.AlignTo(geom.Top, geom.Pt(0.0, 0.0)) // Circ((0,5);5), the anchor moved onto the point

l := geom.Ln(geom.Pt(0, 0), geom.Pt(3, 4))
l.Length()                  // 5
l.Resize(10)                // Line (-1,-2)-(5,6), the same midpoint and direction at length 10
l.Direction()               // DirectionDownRight, nearest of the eight, and Angle in radians
l.Normal()                  // Vector{-4, 3}, a quarter turn of the segment, inward on a rectangle edge
l.DistanceTo(geom.Pt(3, 0)) // 2.4, to the nearest point of the segment
l.Contains(geom.Pt(6, 8))   // false, the segment ends at (3,4)

p := geom.Pol([]geom.Point[int]{{0, 0}, {4, 0}, {4, 4}, {2, 1}, {0, 4}})
p.Area()                  // 10, by the shoelace formula
p.Contains(geom.Pt(2, 3)) // false, inside the notch
p.Edges()[4]              // Line (0,4)-(0,0), closing back to the first vertex

hex := geom.Hexagon(geom.Pt(0, 0), geom.SzU(20), geom.FlatTop)
hex.Bounds() // Rectangle (-20,-17)-(20,17)
hex.Area()   // 1020, in closed form, and Perimeter
hex.Lerp(hex.Rotate(1).Scale(2), 0.5) // halfway in size and angle, along the shorter arc
```

### Intersections

Every pair of shapes has a test on both sides, and the derived result where one exists:

```go
a.Intersects(b)          // any shape with its own kind
r.IntersectsCircle(c)    // and IntersectsLine, IntersectsPolygon, on every shape
c.IntersectsRectangle(r)

point, ok := l.Intersection(m) // where two segments cross
box, ok := a.Intersection(b)   // the overlap of two rectangles
c.Intersection(d)              // zero, one or two points where two circles cross
l.IntersectionCircle(c)        // where a segment crosses the boundary of a circle
l.IntersectionRectangle(r)     // or of a rectangle, or of a polygon with IntersectionPolygon, from Start to End
a.Union(b)                     // the smallest rectangle around both
```

### Directions and axes

```go
dir := geom.DirectionUp
dir.Rotate(2)   // DirectionRight, two 45° steps
dir.Vector(5.0) // Vector{0, -5}

geom.LerpAngle(geom.ToRadians(350), geom.ToRadians(10), 0.5) // 0, along the shorter arc

geom.DirectionFromAxes(up, down, left, right) // keyboard input to an 8-way direction
geom.Vec(3, -7).Direction()                   // DirectionUpRight, nearest of the eight

axis := geom.AxisVertical
axis.Along(size)               // Height, because the axis is vertical
axis.Size(length, thickness)   // Size{thickness, length}
axis.ScaleAcross(size, 0.5)    // halves Width, the cross extent; ScaleAlong halves Height
```

### Matrices

```go
m := geom.IdentityMatrix[float64]().Rotate(math.Pi / 4).Scale(2, 2)

geom.Pt(1.0, 0.0).Transform(m) // Point{1.41, 1.41}
m.Angle()                      // π/4, read back from the matrix
m.Scaling()                    // Vector{2, 2}
m.IsIdentity()                 // false

geom.ShearMatrix(0.5, 0.0)                          // x' = x + 0.5y
geom.ReflectionMatrix[float64](geom.AxisHorizontal) // flips Y

m.Shear(0.5, 0.0)              // compose either one, and PreShear, PreReflect on the other side
m.Reflect(geom.AxisHorizontal)
```

### Number types

Any type satisfying `Number` works, including named types:

```go
type Tile int32

geom.Pt[Tile](3, 4).Add(geom.DirectionRight.Unit[Tile]()) // Point[Tile]{4, 4}
geom.Vec(3, 0).Transform(m)                               // an int vector through a float matrix, rounded
```

`Int()` and `Float()` convert any value to its `int` and `float64` instantiation. The `ints` and `floats` packages
alias those two and add constructors that round or widen from any `Number`:

```go
import (
	"github.com/gravitton/geometry/types/floats"
	"github.com/gravitton/geometry/types/ints"
)

type Grid struct {
	Size     ints.Size   // geom.Size[int]
	CellSize floats.Size // geom.Size[float64]
}

geom.Pt(1.4, 2.6).Int() // Point[int]{1, 3}, rounded
geom.Sz(4, 2).Float()   // Size[float64]{4, 2}
ints.Pt(1.4, 2.6)       // Point[int]{1, 3}
floats.Sz(4, 2)         // Size[float64]{4, 2}
```

### Interop

```go
geom.RectangleFromImage[int](img.Bounds())
geom.Pt(3, 4).Point()                                            // image.Point
geom.RectangleFromMin(geom.Pt(0, 0), geom.Sz(4, 2)).Rectangle() // image.Rectangle

json.Marshal(geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4))) // {"x":1,"y":2,"w":3,"h":4}
json.Marshal(geom.DirectionUp)                        // "Up"
geom.ParseSize[int]("4x2")                            // Size{4, 2}
```

The flat JSON of `Rectangle`, `Circle` and `RegularPolygon` relies on the `embed` struct tag of the v2-backed
`encoding/json`, the default since Go 1.27; under `GOEXPERIMENT=nojsonv2` the nested fields are emitted as objects.

### Testing

One assertion per shape, comparing with the tolerance of the asserted type:

```go
geom.AssertPoint(t, got, geom.Pt(1.0, 2.0))
geom.AssertRectangle(t, got, want, "after inset")
geom.AssertVertices(t, hex.Vertices(), want)
```

Full reference: [pkg.go.dev][link-go-dev-reference].

## Conventions

**Coordinates.** The origin is top-left and `+Y` points down. Angles follow the mathematical convention, so a positive
angle is counterclockwise in math coordinates and appears clockwise on screen. Direction order and polygon winding
follow the same rule: `Directions`, `Rectangle.Vertices` and `RegularPolygon.Vertices` all wind by increasing angle,
clockwise as drawn. Each `Direction` has three names: canonical (`DirectionDownRight`), compass (`SouthEast`) and
rectangle corner (`BottomRight`).

**Numbers.** Products, distances and interpolations are computed in `float64` and rounded back into `T`, so a narrow
`int8` never overflows mid-computation. A float result stored into an integer `T` rounds half away from zero;
`Rectangle` is the exception, truncating half its size toward `Min` so that `Max - Min` stays exactly the size.
`Divide`, `Unscale` and `Matrix.Inverse` on a singular matrix panic, like the integer `/` operator; check
`IsInvertible` first when a matrix may be singular. Every other degenerate input returns a value the type can express.

**Equality.** `Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, a tolerance matched to
`float32` or `float64`. `EqualRelative` scales it for values far from zero, and `EqualAngle` compares modulo a full
turn.

**Boundaries.** `Contains`, `Intersects` and `Vector.LessOrEqual` are closed and tolerant: a point within
`Epsilon[T]()` of the boundary counts as on it, so a float rectangle contains the corners it was built from and a
polygon contains its vertices. `DistanceTo` is zero exactly where `Contains` holds, and `Intersection` answers exactly
where `Intersects` holds, apart from parallel and coincident cases that have no single answer. `Vector.Less` is strict.

**Matrices.** An integer `Matrix` composes lattice transforms exactly: translation, integer scale, reflection, quarter
turns. Anything else rounds into a different matrix; use a float `Matrix` there. `Transform` takes a float matrix, so
convert an integer one with `Float()` at the call.

**Methods.** Every type has `Equal`, `Int` and `String`; every shape adds `Bounds`, `Contains` and `Intersects`.
Each file lists its methods in the same order, and the tests follow it: constructors, properties (`Width`, `Area`,
`Bounds`), arithmetic (`Add`, `Scale`, `Inset`), geometry (`Transform`, `Rotate`, `Project`), relations (`Contains`,
`DistanceTo`, `Intersects`), equality and state (`Equal`, `IsZero`), conversions (`Int`, `Float`), and `String` with
JSON last.

## Planned

- **`Nearest(point)`** – the closest point of a shape to a point, on every shape.
- **`Encloses`** – shape-in-shape containment for culling, distinct from `Contains`, which takes a point.
- **`Rectangle.Clamp(rectangle)`** – moves a rectangle so it lies within another.
- **`Orientation` as a full enum** – `String`, `MarshalText`, `UnmarshalText`, `ParseOrientation` and an `Orientations()` list.
- **Vertex and edge iterators** – public `iter.Seq` forms of `Vertices` and `Edges`, forward and backward as
  `slices.Backward` spells it.
- **`Circle.RegularPolygon(n, orientation)` and `RegularPolygon.Circle()`** – the conversion between the two shapes,
  with two options for where the polygon meets the circle.
- **`Line.Clip()`** – the part of a segment inside a shape.
- **More fuzz targets**
- **`Rectangle.Angle`** – an oriented rectangle. `Contains`, `Clamp`, `Intersects`, `Intersection` and `Union` assume
  axis alignment today.
- **`Polygon.Winding`, `IsConvex` and `ConvexHull`** – convexity also unlocks a separating-axis `Intersects`, the slow
  case in `BenchmarkPolygon_Intersects` today.
- **`Polygon.Simplify(tolerance)`** – drops every vertex within the tolerance of the edge between its neighbours.
- **`Vector.Slerp(vector, t)`** – interpolation of the direction along the shorter arc, on `LerpAngle`-
- **`Polygon.Lerp(polygon, t)`** – vertex-by-vertex interpolation for shape morphing, left out of the `Lerp` pass because
- **`Ray`** – a half-line with origin and direction, for casts against every shape.
- **`Ellipse`** – `RegularPolygon` already takes semi-axes; the continuous shape has no type.
- **`Transform` on every shape** – a rotated rectangle needs `Rectangle.Angle`, a non-uniformly scaled circle is an `Ellipse`.

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

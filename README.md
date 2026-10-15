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

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/hero-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/hero-light.svg">
  <img alt="A tile platformer level: a RegularPolygon sentry with a short Sector showing its facing lights the level with a view built from Raycasts that stop at the Box walls and the Polygon spikes in the pit, its facing ray reflected off the floor at the hit normal; a Capsule hero hides in a pillar's shadow; Circle coins are filled where the light Contains them and outlined in the dark" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/hero-light.svg">
</picture>
</p>

Every picture in this README is drawn from the package's own results: each point is what the code beside it returns.

## Features

- **Generic** – every type works with any integer or float type, named types included.
- **Immutable** – every method returns a new value.
- **Shapes** – circle, capsule, sector, polygon, rectangle, regular polygon, box and ellipse.
- **Paths** – segment, ray, line, polyline, arc and Bezier curves.
- **Collisions** – intersections, containment, raycasts, nearest points, contacts and sweeps.
- **Interfaces** – `Shape` holds a shape without knowing which one.
- **Screen space** – top-left origin, `+Y` down, one winding order everywhere.
- **Closed boundaries** – the edge counts as inside, within a tolerance matched to `T`.
- **Reproducible** – the same bits on amd64 and arm64, trigonometry aside.
- **Extras** – direction enums, `image` interop, JSON, parsing and `geomtest` assertions.

The rules behind them are in the [package documentation][link-go-dev-reference].

## Installation

```shell
go get github.com/gravitton/geometry
```

```go
import geom "github.com/gravitton/geometry"
```

## Usage

The level at the top, call by call:

```go
ray := geom.Ry(view.Center, geom.VectorFromAngle(angle, 1.0)) // one ray of the sentry's view
exit, _ := ray.ClipBox(scene)                                  // it leaves the level here
hit, ok := geom.Raycast(ray, solid)                            // or stops at a wall or the spikes: Point, Normal
bounce := ray.Vector.Reflect(hit.Normal)                       // the facing ray reflected off the floor
light := geom.Pol(outline)                                     // the rays' ends outline what the sentry sees
light.Contains(coin.Center)                                    // a coin in sight
```

Every figure below is drawn on a grid of one unit per cell, from the values in the code beside it.

### Shapes

Every kind has a shorthand constructor in its own letters, and named constructors that read as phrases
(`CircleThrough`, `RectangleFromMinMax`, `RegularPolygonWithOrientation`).

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-light.svg">
  <img alt="The nine shapes on a unit grid: Circle, Capsule, Sector, Segment, Polygon, Rectangle, RegularPolygon, Box and Ellipse" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-light.svg">
</picture>
</p>

```go
geom.Circ(geom.Pt(4.5, 4.0), 3.0)                                                        // center, radius
geom.Cap(geom.Pt(11.5, 5.5), geom.Pt(15.5, 2.5), 1.5)                                    // every point within the radius of a segment
geom.Sect(geom.Pt(22.5, 6.5), 5.0, -3*geom.Pi/4, geom.Pi/2)                              // center, radius, angle and sweep in radians
geom.Seg(geom.Pt(28.5, 6.0), geom.Pt(34.5, 1.0))                                         // from Start to End
geom.Pol([]geom.Point[float64]{{37.5, 1}, {43.5, 1}, {43.5, 6}, {40.5, 3.5}, {37.5, 6}}) // closed, convex or not
geom.Rect(geom.Pt(9.0, 11.0), geom.Sz(6.0, 3.0), 0.4)                                    // center, size and angle
geom.RegularPolygonWithOrientation(geom.Pt(18.0, 11.0), geom.SzU(2.5), 6, geom.OrientationFlatTop)
geom.Bx(geom.Pt(24.0, 9.0), geom.Pt(30.0, 13.0))      // axis-aligned, the corners in either order
geom.Ell(geom.Pt(36.0, 11.0), geom.Sz(3.5, 1.8), 0.5) // center, semi-axes and angle
```

### Paths

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-light.svg">
  <img alt="The paths on a unit grid: Ray, Line, Polyline, Arc, QuadraticBezier and CubicBezier" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-light.svg">
</picture>
</p>

```go
geom.Ry(geom.Pt(1.0, 6.0), geom.Vec(2.0, -1.0))                       // from Origin along a direction, without end
geom.Ln(geom.Pt(11.0, 4.0), geom.Vec(2.0, 1.0))                       // without end both ways
geom.PolLn([]geom.Point[float64]{{16, 7}, {18, 1}, {20, 5}, {22, 1}}) // an open path
geom.Ac(geom.Pt(27.0, 4.0), 3.0, -geom.Pi, 1.25*geom.Pi)              // a part of a circle, the sweep signed
geom.QuadBez(geom.Pt(31.0, 7.0), geom.Pt(34.0, -1.0), geom.Pt(37.0, 7.0))
geom.CubBez(geom.Pt(38.0, 6.0), geom.Pt(39.0, 0.0), geom.Pt(43.0, 8.0), geom.Pt(44.0, 2.0))
```

### Points and vectors

A `Point` is a position and a `Vector` a displacement; subtracting two points gives the vector between them.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/vectors-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/vectors-light.svg">
  <img alt="Points a and b, the vector v between them, its normal, its reflection off a horizontal wall, a point a quarter of the way from a to b, and a turned a quarter turn about b" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/vectors-light.svg">
</picture>
</p>

```go
a, b := geom.Pt(2.0, 8.0), geom.Pt(8.0, 5.0)
v := b.Subtract(a) // ⟨6,-3⟩

v.Length()                    // 6.71
v.Normal()                    // ⟨3,6⟩, perpendicular, of the same length
v.Reflect(geom.Vec(0.0, 1.0)) // ⟨6,3⟩, bounced off a horizontal wall
a.Lerp(b, 0.25)               // (3.5,7.25)
a.RotateAround(b, -geom.Pi/2) // (11,11)
```

`Lerp` moves by a fraction; a game loop moves by a speed times the frame time, a step that never passes the target:

```go
a.MoveToward(b, 2.5)                          // 2.5 along the way to b, and b itself once the step reaches it
v.RotateToward(geom.Pi, 0.1)                  // turned by 0.1 along the shorter arc, the length kept
geom.MoveTowardAngle(angle, target, speed*dt) // the same step on a bare angle
geom.MoveToward(volume, 1.0, fade*dt)         // and on any number
```

### Sizes and padding

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/sizes-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/sizes-light.svg">
  <img alt="A 16x9 screen size fitted into and filled around a 6x6 slot" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/sizes-light.svg">
</picture>
</p>

```go
screen, slot := geom.Sz(16.0, 9.0), geom.SzU(6.0)
screen.Fit(slot)  // 6x3.38, the largest of the same ratio inside
screen.Fill(slot) // 10.67x6, the smallest of the same ratio around

geom.PadXY(1.0, 0.5).Size() // 2x1, the horizontal and vertical total
```

### Rectangles

A `Rectangle` is a center, a size and an angle. Turned, it keeps its corner names, and `Bounds` becomes the box
around its vertices.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangles-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangles-light.svg">
  <img alt="A turned rectangle r with its corner names, Bounds, BoundingCircle and top anchor; a window w clamped into a screen" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangles-light.svg">
</picture>
</p>

```go
r := geom.Rect(geom.Pt(6.0, 5.5), geom.Sz(7.0, 3.5), 0.4)
r.TopLeft()        // (3.46,2.53), the corner that was top-left before the turn
r.Anchor(geom.Top) // (6.68,3.89), the midpoint of that edge
r.Bounds()         // (2.09,2.53)-(9.91,8.47)
r.BoundingCircle() // Circ((6,5.5);3.91)

screen := geom.RectangleFromMinMax(geom.Pt(15.0, 1.0), geom.Pt(27.0, 10.0))
w := geom.RectangleFromCenter(geom.Pt(26.0, 1.5), geom.Sz(5.0, 3.0))
w.Clamp(screen) // Rect((24.5,2.5);5x3), moved by the least that brings it inside
```

### Boxes

A `Box` is the axis-aligned extent, for clipping, viewports, layout and culling, and the `Bounds` of every shape.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/boxes-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/boxes-light.svg">
  <img alt="A view box with a HUD pinned to its top-right corner, a tooltip clamped inside it, and the part of the tooltip on screen" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/boxes-light.svg">
</picture>
</p>

```go
view := geom.BoxFromSize(geom.Sz(20.0, 12.0))
view.Inset(geom.PadU(1.0)) // (1,1)-(19,11)

hud := geom.Bx(geom.Pt(6.0, 2.0), geom.Pt(0.0, 0.0))   // (0,0)-(6,2), the corners in either order
hud.AlignTo(geom.TopRight, view.Anchor(geom.TopRight)) // (14,0)-(20,2), pinned to the corner of the view

tip := geom.BoxFromMin(geom.Pt(17.0, 8.0), geom.Sz(7.0, 3.0))
tip.Clamp(view)   // (13,8)-(20,11)
view.ClipBox(tip) // (17,8)-(20,11), true: the part on screen
```

### Circles

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/circles-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/circles-light.svg">
  <img alt="The tangents from a point to circle c, the outer and inner tangents of c and a wheel, and the smallest circle around both" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/circles-light.svg">
</picture>
</p>

```go
c, wheel := geom.Circ(geom.Pt(3.0, 5.0), 2.0), geom.Circ(geom.Pt(8.0, 6.0), 1.0)

c.Tangents(geom.Pt(8.5, 1.5)) // a Pair of the points the lines from there touch
c.OuterTangents(wheel)        // a belt around both, each a Segment from c to wheel
c.InnerTangents(wheel)        // a crossed belt; none once the circles touch
c.Union(wheel)                // Circ((5.01,5.4);4.05)
```

### Ellipses and capsules

A circle has no `Transform` of its own, since a matrix takes it to an ellipse. A `Capsule` is a segment with a
radius, the rounded collider of a character.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/ellipses-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/ellipses-light.svg">
  <img alt="A circle transformed into an ellipse with its foci, a point inside it and the distance to a point outside; a capsule and the distance to a point beside it" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/ellipses-light.svg">
</picture>
</p>

```go
o := geom.Circ(geom.Pt(4.0, 5.0), 2.5)
e := o.Ellipse().Transform(geom.ScaleMatrix(2.0, 1.0)) // Ell((8,5);5x2.5), exact for every matrix
e.Contains(geom.Pt(6.0, 4.0))                          // true
e.DistanceTo(geom.Pt(10.0, 8.5))                       // 1.18, to the nearest point of the boundary
e.Foci()                                               // (3.67,5), (12.33,5)

body := geom.Cap(geom.Pt(4.0, 2.0), geom.Pt(4.0, 8.0), 1.5)
body.DistanceTo(geom.Pt(8.5, 5.0)) // 3, to the side
```

### Polygons

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-light.svg">
  <img alt="A notched polygon p with its ConvexHull, its triangles, its Grow and Shrink outlines, and its two parts after a Split" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-light.svg">
</picture>
</p>

```go
p := geom.Pol([]geom.Point[float64]{{1, 1}, {9, 1}, {9, 9}, {5, 4}, {1, 9}})
p.Area()        // 44
p.IsConvex()    // false
p.ConvexHull()  // Pol((1,1);(9,1);(9,9);(1,9)), the notch dropped
p.Triangulate() // [1 2 3 0 1 3 4 0 3]: three indices into Points per triangle
p.Grow(0.5)     // every edge moved out along its normal; Shrink moves them in

front, back := p.Split(geom.Ln(geom.Pt(3.0, 0.0), geom.Vec(1.0, 2.0))) // meeting exactly along the line

for vertex := range p.Vertices() { // the same loop walks a Rectangle, RegularPolygon or Box
}
```

### Regular polygons

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-light.svg">
  <img alt="Flat-top and pointy-top hexagons with their vertices numbered and their top anchor, a stretched hexagon on its ellipse, and a pentagon with a phase or an angle" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-light.svg">
</picture>
</p>

```go
hex := geom.RegularPolygonWithOrientation(geom.Pt(5.0, 5.0), geom.SzU(4.0), 6, geom.OrientationFlatTop)
hex.Anchor(geom.Top) // (5,1.54), the midpoint of the top edge; a pointy-top one gives its top vertex
hex.Bounds()         // (1,1.54)-(9,8.46)
hex.Area()           // 41.57

geom.RegularPolygonWithOrientation(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 6, geom.OrientationFlatTop) // on an ellipse
geom.RegPol(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 5, 0, 0.6)                                         // the phase slides the vertices along the ellipse
geom.RegPol(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 5, 0.6, 0)                                         // the angle turns the ellipse
```

### Intersections

Every shape has `Intersects<Kind>`, `Contains<Kind>` and `Intersection<Kind>` for each kind, symmetric and counting a
touch. `Intersection<Kind>` gives the boundary points, a `Pair` where there are at most two; `Clip<Kind>` gives the
part in common.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-light.svg">
  <img alt="The crossing points of two circles, the parts of a segment inside a notched polygon, the two parts of a segment through the notch of a wide sector, and the rectangle two turned rectangles share" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-light.svg">
</picture>
</p>

```go
a, b := geom.Circ(geom.Pt(4.0, 5.0), 3.0), geom.Circ(geom.Pt(7.0, 6.0), 2.5)
a.IntersectionCircle(b) // Pair((5.21,7.75);(6.62,3.53)), which allocates nothing

s := geom.Seg(geom.Pt(0.0, 6.0), geom.Pt(10.0, 6.0))
s.ClipPolygon(p) // [Seg((1,6);(3.4,6)) Seg((6.6,6);(9,6))], from Start to End

cone := geom.Sect(geom.Pt(5.0, 5.0), 4.0, -geom.Pi/4, 1.5*geom.Pi)
geom.Seg(geom.Pt(0.0, 3.5), geom.Pt(10.0, 3.5)).ClipSector(cone) // two parts through the notch

r := geom.Rect(geom.Pt(4.0, 4.5), geom.Sz(5.0, 3.5), 0.3)
r.ClipRectangle(geom.Rect(geom.Pt(6.5, 6.0), geom.Sz(4.0, 3.0), 0.3)) // the overlap, a Rectangle of the same angle

geom.Intersects[float64](a, r) // the same answers for two shapes of any kind
```

### Casts

`Hit<Kind>` on a segment or a ray is the first point of the boundary it reaches from outside, with the unit normal
there and the fraction of the cast; `Sweep` moves a whole shape.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-light.svg">
  <img alt="A ray hitting a turned wall, with the normal and the bounce; a ball swept into a block, stopped at the fraction of its motion" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-light.svg">
</picture>
</p>

```go
wall := geom.Rect(geom.Pt(10.0, 5.0), geom.Sz(2.0, 6.0), 0.35)
ray := geom.Ry(geom.Pt(1.0, 7.0), geom.Vec(3.0, -1.0))
hit, ok := ray.HitRectangle(wall) // Hit((9.2,4.27);⟨-0.94,-0.34⟩;2.73), true
ray.Vector.Reflect(hit.Normal)    // ⟨-1.65,-2.7⟩, the bounce

ball := geom.Circ(geom.Pt(2.0, 2.0), 1.0)
block := geom.Bx(geom.Pt(8.0, 4.0), geom.Pt(11.0, 9.0))
hit, ok = geom.Sweep[float64](ball, geom.Vec(11.0, 5.0), block) // Fraction 0.45: where the ball first touches
geom.Raycast[float64](ray, wall)                                // a Hit on a shape of any kind
```

### Lines

A `Line` runs without end both ways and tells its two sides apart.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/lines-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/lines-light.svg">
  <img alt="A line through two points, a point a on one side of it, its nearest point on the line and its reflection" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/lines-light.svg">
</picture>
</p>

```go
wall := geom.LineThrough(geom.Pt(1.0, 8.0), geom.Pt(13.0, 2.0))
a := geom.Pt(10.0, 8.0)
wall.Side(a)       // 1, the side the normal of the direction points to
wall.Nearest(a)    // (8.2,4.4)
wall.Reflect(a)    // (6.4,0.8)
wall.ClipBox(view) // the part a viewport shows
```

### Collisions

Two shapes of any kind are parted, measured and moved against each other by package functions.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-light.svg">
  <img alt="A ball and a crate overlapping a floor with their contact points, normals and separated position, and the nearest points of two shapes with the distance between them" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-light.svg">
</picture>
</p>

```go
floor := geom.Bx(geom.Pt(1.0, 6.0), geom.Pt(9.0, 8.0))
ball := geom.Circ(geom.Pt(4.0, 5.0), 1.5)
contact, ok := geom.Collide[float64](ball, floor) // Contact(⟨0,-1⟩;0.5;(4,6)): Normal, Depth, Points
ball.Translate(contact.Separation())              // the shortest move out of the floor

crate := geom.RectangleFromCenter(geom.Pt(5.0, 5.0), geom.Sz(3.0, 2.5))
geom.Collide[float64](crate, floor) // two Points where the sides lie flat

on, at := geom.Nearest[float64](turret, target) // the point of each nearest the other
geom.Distance[float64](turret, target)          // 2.21, zero exactly where they intersect
```

`Shape` holds a shape without knowing which one, for a spatial index or a broad pass:

```go
var s geom.Shape[float64] = c // Bounds, Contains, DistanceTo, DistanceSquaredTo, Nearest, Support
```

The functions take the type argument where a concrete shape is passed and infer it from a value held as a `Shape`. A
call through an interface allocates where the same named method does not, so a hot loop that knows its types calls
the method.

### Arcs, sectors and curves

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-light.svg">
  <img alt="An arc as a progress ring, points inside and outside a vision cone, a step along a polyline, and a cubic Bezier flattened to a polyline" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-light.svg">
</picture>
</p>

```go
ring := geom.Ac(geom.Pt(5.0, 5.0), 4.0, -geom.Pi/2, 1.3*geom.Pi) // a full turn is the circle
ring.PointAt(0.5)                                                // (8.56,6.82), also Start, End and Length

cone := geom.Sect(geom.Pt(1.0, 5.0), 8.5, -geom.Pi/6, geom.Pi/3)
cone.Contains(geom.Pt(5.0, 8.5)) // false; the tolerance is on the distance, never an angle

route := geom.PolLn([]geom.Point[float64]{{1, 8}, {4, 2}, {6, 7}, {9, 2}})
for step := range route.Every(1) { // a point every unit of length, in one walk
}

curve := geom.CubBez(geom.Pt(1.0, 9.0), geom.Pt(1.0, 1.0), geom.Pt(9.0, 1.0), geom.Pt(9.0, 9.0))
curve.PointAt(0.5)  // (5,3); the parameter is not a fraction of the length
curve.Polyline(0.2) // a path within the tolerance, to measure the curve by
```

### Coordinates and directions

The origin is top-left and `+Y` points down. Angles grow counterclockwise in math terms, so clockwise on screen, and
every outline winds that way.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-light.svg">
  <img alt="The screen axes with +Y down and the sense of a positive angle, the eight directions, and the clockwise winding of a rectangle's vertices" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-light.svg">
</picture>
</p>

```go
dir := geom.DirectionUp
dir.Turn(2)                 // DirectionRight, two 45° steps
dir.Vector(5.0)             // ⟨0,-5⟩
geom.Vec(3, -7).Direction() // DirectionUpRight, the nearest of the eight

geom.DirectionFromAxes(up, down, left, right) // keyboard input to an 8-way direction
```

### Matrices

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/matrices-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/matrices-light.svg">
  <img alt="An arrow turned and scaled about a pivot and sheared, its original outline dashed" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/matrices-light.svg">
</picture>
</p>

```go
geom.RotationMatrixAround(5.0, 5.0, geom.Pi/4) // turn about a pivot
geom.ScaleMatrixAround(5.0, 5.0, 1.5, 0.5)     // scale about a pivot, which stays where it is
geom.ShearMatrix(0.5, 0.0)                     // x moved by half of y

m := geom.IdentityMatrix[float64]().Rotate(geom.Pi / 4).Scale(2, 2)
arrow.Transform(m) // every shape, point and vector takes a matrix
m.Angle()          // π/4, read back from the matrix
```

`Translate`, `Rotate`, `RotateAround`, `Scale`, `ScaleAround`, `Shear` and `Reflect` compose with their constructor on
the right, applied first, and their `Pre` forms on the left, applied last.

### Numbers

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/numbers-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/numbers-light.svg">
  <img alt="A stick axis from -1 to 1 remapped onto a steering angle from -30 to 30" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/numbers-light.svg">
</picture>
</p>

```go
geom.Remap(0.5, -1.0, 1.0, -30.0, 30.0)                      // 15, a stick axis as a steering angle
geom.Lerp(10.0, 20.0, 0.25)                                  // 12.5
geom.InverseLerp(10.0, 20.0, 12.5)                           // 0.25, the ratio Lerp takes
geom.Clamp(15.0, 0.0, 10.0)                                  // 10
geom.Snap(37, 16)                                            // 32, the nearest multiple of the step
geom.Pt(37, 21).Snap(16)                                     // (32,16), and on points and vectors
geom.LerpAngle(geom.ToRadians(350), geom.ToRadians(10), 0.5) // across the zero angle, the short way
```

### Number types

Any type satisfying `Number` works, named types included. `Int()` rounds half away from zero, and the `ints` and
`floats` packages alias the two common instantiations.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rounding-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rounding-light.svg">
  <img alt="Five float points and the lattice points they round to" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rounding-light.svg">
</picture>
</p>

```go
geom.Pt(1.4, 2.6).Int() // (1,3)
geom.Pt(3.5, 1.5).Int() // (4,2), a half rounds away from zero
geom.Pt(6.6, 3.2).Int() // (7,3)
geom.Pt(2.5, 7.5).Int() // (3,8)
geom.Pt(7.2, 6.7).Int() // (7,7)

type Tile int32
geom.Pt[Tile](3, 4).Add(geom.DirectionRight.Unit[Tile]()) // (4,4)
ints.Pt(1.4, 2.6)                                         // (1,3), from github.com/gravitton/geometry/types/ints
```

### Interop

`Int()` rounds a float rectangle onto the lattice. An integer rectangle is closed: its corners are lattice points on
the boundary, so a width of 6 spans 7 columns of them. The `image.Rectangle` it converts to is half-open and covers
exactly 6 pixels across.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/pixels-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/pixels-light.svg">
  <img alt="A float rectangle f rounded by Int to a 6 by 6 integer rectangle with its 7 by 7 lattice points, over the 6 by 6 pixels of its image.Rectangle" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/pixels-light.svg">
</picture>
</p>

```go
f := geom.RectangleFromMin(geom.Pt(1.6, 2.3), geom.Sz(5.8, 5.6))
r := f.Int()  // Rect((5,5);6x6), the center and size rounded
r.Rectangle() // (2,2)-(8,8), an image.Rectangle of 6x6 pixels
geom.RectangleFromImage[int](img.Bounds())
geom.BoxFromImage[int](img.Bounds()) // Min and Max as they are

json.Marshal(geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4), 0)) // {"x":1,"y":2,"w":3,"h":4}
json.Marshal(geom.DirectionUp)                           // "Up"
geom.ParseSize[int]("4x2")                               // 4x2
geom.ParsePoint[int]("(1,2)")                            // (1,2)
```

### Testing

The `geomtest` package holds one assertion per type, with the tolerance of the asserted type, so the main package
imports no test library:

```go
geomtest.AssertPoint(t, got, geom.Pt(1.0, 2.0))
geomtest.AssertRectangle(t, got, want, "after inset")
```

Full reference: [pkg.go.dev][link-go-dev-reference].

## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]

## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/geometry.svg?style=flat-square&color=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/geometry/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/geometry?style=flat-square&color=blue
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

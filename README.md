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

- **Generic** – every type works with any integer or float type, named types included.
- **Immutable** – every method returns a new value.
- **Shapes** – rectangle, box, circle, ellipse, segment, ray, polygon and regular polygon.
- **Collisions** – intersections, containment, ray casts, nearest points, separation and sweeps.
- **Interfaces** – `Shape` and `Body` hold a shape without knowing which one.
- **Screen space** – top-left origin, `+Y` down, one winding order everywhere.
- **Closed boundaries** – the edge counts as inside, within a tolerance matched to `T`.
- **Reproducible** – the same bits on amd64 and arm64.
- **Extras** – direction enums, `image` interop, JSON, parsing and `geomtest` assertions.

The rules behind them are in the [package documentation][link-go-dev-reference].

## Installation

```shell
go get github.com/gravitton/geometry
```

## Usage

```go
import geom "github.com/gravitton/geometry"
```

```go
ship := geom.Circ(geom.Pt(0.0, 0.0), 10.0)
rock := geom.RegularPolygonWithOrientation(geom.Pt(50.0, 0.0), geom.SzU(20.0), 6, geom.OrientationFlatTop).Rotate(0.3)

ship.IntersectsRegularPolygon(rock)                                // false
ship.Translate(geom.Vec(25.0, 0.0)).IntersectsRegularPolygon(rock) // true, after the move

laser := geom.Ry(ship.Center, geom.Vec(1.0, 0.0))
hit, ok := laser.ClipRegularPolygon(rock) // Seg((32.24,0.00);(67.76,0.00)), true; hit.Start is the first point hit
```

### Points and vectors

```go
p := geom.Pt(1, 2)
v := geom.Pt(4, 6).Subtract(p) // Vector{3, 4}

v.Length()          // 5
v.Normal()          // Vector{-4, 3}, perpendicular
p.Add(v.Resize(10)) // Point{7, 10}

v.Reflect(geom.Vec(0, 1))                 // Vector{3, -4}, bounced off a horizontal wall
p.RotateAround(geom.Pt(0, 0), math.Pi/2)  // Point{-2, 1}
p.Lerp(geom.Pt(9, 10), 0.25)              // Point{3, 4}
geom.Vec(8, 0).Slerp(geom.Vec(0, 8), 0.5) // Vector{6, 6}, turned along the arc, where Lerp gives Vector{4, 4}
```

`Lerp` and `Slerp` move by a fraction; a game loop moves by a speed times the frame time, a step that never passes
the target:

```go
p.MoveToward(geom.Pt(4, 6), 2.5)              // Point{3, 4}, and Point{4, 6} once the step reaches it
geom.Vec(2.0, 0).RotateToward(math.Pi, 0.1)   // turned by 0.1 along the shorter arc, the length kept
velocity.MoveToward(wanted, accel*dt)         // a velocity taken toward another by an acceleration
geom.MoveTowardAngle(angle, target, speed*dt) // the same step on a bare angle
geom.MoveToward(volume, 1.0, fade*dt)         // and on any number
```

Vectors also have `Project`, `Reject`, `AngleBetween`, and `AtMost` and `AtLeast` to cap or floor a length.

### Sizes and padding

```go
s := geom.Sz(1920, 1080)
s.Scale(0.5)                 // Size{960, 540}
s.AtMost(geom.SzU(800))      // Size{800, 800}, clamped per axis
s.Shrink(2000).AtLeastZero() // Size{0, 0}, the clamp every shape applies to its own size
s.Fit(geom.SzU(800))         // Size{800, 450}, largest with the same ratio inside
s.Fill(geom.SzU(800))        // Size{1422, 800}, smallest with the same ratio around

geom.PadXY(4, 8).Size() // Size{16, 8}, horizontal and vertical total
```

### Rectangles

```go
r := geom.Rect(geom.Pt(50, 50), geom.Sz(20, 10)) // center and size

r.Contains(geom.Pt(55, 52))                 // true
r.Nearest(geom.Pt(80, 0))                   // Point{60, 45}, the nearest point, on every shape
r.Inset(geom.PadU(2)).Anchor(geom.TopRight) // Point{58, 47}
r.AlignTo(geom.TopLeft, geom.Pt(0, 0))      // Rectangle (0,0)-(20,10)

screen := geom.RectangleFromMinMax(geom.Pt(0, 0), geom.Pt(55, 100))
r.Clamp(screen) // Rectangle (35,45)-(55,55), moved by the least that brings it inside

b := geom.RectangleFromMinMax(geom.Pt(0, 0), geom.Pt(8, 6))
b.Scale(2)                    // Rectangle (-4,-3)-(12,9), scaled around the center
for edge := range b.Edges() { // clockwise from the top edge, Segment (0,0)-(8,0), without allocating
	edge.Midpoint()
}
slices.Collect(b.Vertices()) // the four corners as a slice, clockwise from the top-left

d := geom.Rect(geom.Pt(0.0, 0.0), geom.Sz(2.0, 2.0)).Rotate(geom.Pi / 4) // a diamond
d.Contains(geom.Pt(0.9, 0.9))                                            // false, outside the turned edges
d.Bounds()                                                               // the Box around it
d.BoundingCircle()                                                       // the Circle around it, the same at any angle
d.TopLeft()                                                              // the corner that was top-left before the turn

d.Transform(geom.RotationMatrix[float64](geom.Pi / 2)) // a turn is exact, and so is a move, a reflection or a scale
d.Transform(geom.ShearMatrix(1.0, 0.0))                // the nearest rectangle; Polygon() holds the parallelogram
```

A rectangle turned by `Rotate` keeps its center, size and corner names; `Min`, `Max` and `Bounds` become the box
around its vertices. Two rectangles of the same angle intersect and unite in a rectangle of that angle.

### Boxes

```go
view := geom.BoxFromSize(geom.Sz(800, 600)) // Min and Max, no angle: (0,0)-(800,600)
view.Contains(geom.Pt(800, 600))            // true, closed like every shape
view.Inset(geom.PadU(16))                   // (16,16)-(784,584)

tip := geom.BoxFromMin(geom.Pt(750, 20), geom.Sz(120, 40))
hud := geom.Bx(geom.Pt(200, 40), geom.Pt(0, 0)) // (0,0)-(200,40), the corners in either order

tip.Clamp(view)   // (680,20)-(800,60), moved by the least that brings it inside
view.ClipBox(tip) // (750,20)-(800,60), true: the part on screen

r.Bounds().IntersectsBox(view) // every shape's Bounds is a Box, for the broad pass
view.IntersectsRectangle(r)    // a Box is a Shape, tested against every other
```

`Rectangle` is the shape that turns; `Box` is the axis-aligned extent, for clipping, viewports, layout and culling.

### Circles and ellipses

```go
c := geom.Circ(geom.Pt(0.0, 0.0), 5.0)
c.Anchor(geom.Bottom)          // Point{0, 5}
c.Tangents(geom.Pt(13.0, 0.0)) // a Pair: the two points a line from there touches the circle at; one on it, none inside

wheel := geom.Circ(geom.Pt(20.0, 0.0), 2.0)
c.OuterTangents(wheel) // the two straight runs of a belt around both, each a Segment from c to wheel
c.InnerTangents(wheel) // the two of a crossed belt; none once the circles touch

geom.CircleBetween(geom.Pt(0.0, 0.0), geom.Pt(6.0, 8.0))                   // Circ((3,4);5), the two across a diameter
geom.CircleThrough(geom.Pt(0.0, 0.0), geom.Pt(6.0, 0.0), geom.Pt(0.0, 8.0)) // Circ((3,4);5), true; false on one line

// a circle has no Transform of its own: an affine matrix takes it to an ellipse
e := c.Ellipse().Transform(geom.ScaleMatrix(2.0, 1.0)) // Ellipse, semi-axes 10x5, exact for every matrix
e.Contains(geom.Pt(9.0, 0.0))                          // true
e.DistanceTo(geom.Pt(0.0, 9.0))                        // 4, to the nearest point of the boundary
e.Rotate(geom.Pi / 2).Bounds()                         // the box around the turned ellipse, 10x20
e.Foci()                                               // the two focal points, on the major axis
e.Circle()                                             // the circle around it, of the major semi-axis
```

Neither has vertices: the `RegularPolygon` of the wanted resolution is what draws or
walks one, and it converts back exactly.

```go
e.RegularPolygon(64, geom.OrientationPointyTop) // every vertex on the boundary, and Ellipse() converts back
c.RegularPolygon(6, geom.OrientationFlatTop)     // the orientation places the first vertex before the stretch

slices.Collect(c.RegularPolygon(64, geom.OrientationFlatTop).Vertices()) // the points that draw the circle
```

### Segments and polygons

```go
s := geom.Seg(geom.Pt(0, 0), geom.Pt(3, 4))
s.Length()                  // 5
s.DistanceTo(geom.Pt(3, 0)) // 2.4, to the nearest point of the segment
s.Contains(geom.Pt(6, 8))   // false, the segment ends at (3,4)

p := geom.Pol([]geom.Point[int]{{0, 0}, {4, 0}, {4, 4}, {2, 1}, {0, 4}})
p.Area()                  // 10
p.Contains(geom.Pt(2, 3)) // false, inside the notch
p.Points[3]               // Point{2, 1}, the notch
p.Winding()               // WindingClockwise, the winding of a Rectangle's corners
p.IsConvex()              // false, the notch turns the other way
p.IsSimple()              // true, no two edges meet but neighbours at their vertex
p.ConvexHull()            // Pol((0,0);(4,0);(4,4);(0,4)), the notch dropped
p.BoundingCircle()        // Circ((2,1);4), about the rounded Centroid to the farthest vertex, the radius rounded up

edged := geom.Pol([]geom.Point[int]{{0, 0}, {2, 0}, {4, 0}, {4, 4}, {0, 4}})
edged.Simplify(0) // Pol((0,0);(4,0);(4,4);(0,4)), the vertex on an edge dropped; a tolerance drops the ones near it
edged.Grow(1)     // Pol((-1,-1);(2,-1);(5,-1);(5,5);(-1,5)): every edge moved out; Shrink moves them in

for vertex := range p.Vertices() { // the same loop draws a Segment, Rectangle or RegularPolygon
	vertex.Float()
}
```

### Regular polygons

```go
hex := geom.RegularPolygonWithOrientation(geom.Pt(0, 0), geom.SzU(20), 6, geom.OrientationFlatTop)
hex.Anchor(geom.Top)                 // Point{0, -17}, the midpoint of the top edge; a pointy-top one gives its top vertex
hex.AlignTo(geom.Top, geom.Pt(0, 0)) // the hexagon moved so that anchor is at the origin
hex.Bounds()                         // Box (-20,-17)-(20,17)
hex.Area()                           // 1039, 3√3/2 · r²
hex.Contains(geom.Pt(10, 5))         // true, walked on the edges without building the vertices
hex.Ellipse()                        // the ellipse its vertices lie on, exactly

geom.RegularPolygonWithOrientation(geom.Pt(0, 0), geom.Sz(20, 10), 4, geom.OrientationFlatTop) // stretched across its edges: a 28x14 rectangle
geom.RegPol(geom.Pt(0, 0), geom.Sz(20, 10), 8, 0, geom.Pi/8)                                   // n, angle and phase: the phase places the vertices, the angle turns the shape
hex.Transform(geom.ShearMatrix(1.0, 0.0))                                                      // exact for every matrix, shears included
```

### Interfaces

Every shape has `Translate`, `Bounds`, `Contains`, `DistanceTo` and `Nearest`, and all but `Box` also `MoveTo`,
`Scale`, `Unscale` and `Rotate`. Two interfaces name what they share, so a spatial index or a broad collision pass
holds a shape without knowing which one:

```go
var s geom.Shape[float64] = c // Bounds, Contains, DistanceTo, DistanceSquaredTo, Nearest, Support — every shape
s.DistanceTo(geom.Pt(10.0, 5.0))
s.Support(geom.Vec(1.0, 0.0)) // the point of it farthest along a direction

geom.Intersects(s, other) // two shapes of any kind, by the method of one naming the other where there is one
geom.Distance(s, other)   // the gap between them, zero exactly where they intersect; DistanceSquared to compare
geom.Raycast(ray, s)      // the Hit of a ray on a shape of any kind

var b geom.Body[float64] = c // Area, Centroid, Inertia — the mass properties at unit density,
b.Inertia()                  // for a physics body to scale by its density; every shape but Segment and Box
```

The functions take the type argument where a concrete shape is passed, `geom.Intersects[float64](c, r)`, and infer
it from a value already held as a `Shape`. The interfaces are for the code around a hot loop: a call through one allocates, where the same call on a concrete
shape does not. `Ellipse` has no `Intersects<Kind>` methods: the functions test it by its `Support`, as they do a
shape of another package. `Ray` and `Line` run without end and are no `Shape`: they test a shape by their own methods.

Two shapes are parted and moved against each other whatever kind they are:

```go
on, at := geom.Nearest(turret, target)    // the point of each nearest the other, one point where they intersect
contact, ok := geom.Collide(crate, floor) // how they meet: Normal, Depth and one or two Points
crate.Translate(contact.Separation())     // the shortest move out of the floor, the normal as far as the depth
hit, ok := geom.Sweep(bullet, step, wall)     // where the bullet first touches the wall while it moves by step:
bullet.Translate(step.Multiply(hit.Fraction)) // the point, the normal and the fraction of the step made by then
```

### Intersections

Every shape has one `Intersects<Kind>` per kind, symmetric and counting a touch. `Intersection<Kind>` gives where the
two meet: one point with a bool between segments, rays and lines, a `Pair` of at most two points on a circle or a
convex shape, and a slice on a polygon, which may be concave. `Clip<Kind>` gives the part the two have in common, in
the type of the receiver: a segment of a segment, a rectangle of a rectangle:

```go
a.IntersectsRectangle(b) // IntersectsSegment, IntersectsRay, IntersectsLine, IntersectsRectangle, IntersectsCircle,
r.IntersectsCircle(c)    // IntersectsPolygon, IntersectsRegularPolygon and IntersectsBox, on every shape but Ellipse
geom.Intersects(a, c)    // the same answer without knowing either type

point, ok := s.IntersectionSegment(other) // where two segments cross
crossings := c.IntersectionCircle(other)  // a Pair: zero, one or two points where two circles cross
for p := range crossings.All() {          // or crossings.Items[:crossings.Count], or AppendTo(buffer); none allocates
}
s.IntersectionCircle(c)                   // where a segment crosses a boundary; also of a rectangle or a box
s.IntersectionPolygon(p)                  // every crossing of a polygon, as a slice
s.ClipBox(view)                           // the part of s inside, and false where there is none
s.ClipPolygon(p)                          // the parts inside a concave polygon, from Start to End
part, ok := p.ClipBox(view)               // the polygon as a viewport shows it; also ClipRectangle, ClipPolygon
                                          // and ClipRegularPolygon
overlap, ok := a.ClipRectangle(b)         // the overlap of two rectangles of one angle; ClipBox of two boxes
hit, ok := s.HitCircle(c)                 // the first point s reaches from outside, with the normal there
hit, ok = s.HitSegment(wall)              // or of a wall: HitSegment, HitRay and HitLine

buffer = s.AppendIntersectionPolygon(buffer[:0], p) // a slice result has an Append form, as strconv does, on the
hull = p.AppendConvexHull(hull[:0])                 // side that holds it, so a reused buffer allocates nothing
```

A `Ray` is a half-line from `Origin` along `Direction`. Its `Clip<Kind>` is the part inside a shape, and its
`Hit<Kind>` the first point of the boundary it reaches from outside, with the normal there:

```go
ray := geom.RayThrough(camera, cursor) // from the camera through the cursor; Ry takes a direction
ray.IntersectsRectangle(r)             // IntersectsRay on the shape gives the same answer
ray.ClipPolygon(p)                     // the parts inside a concave polygon, from Origin on
ray.IntersectionBox(box)               // the boundary crossings, from Origin on

if part, ok := ray.ClipCircle(c); ok { // allocation-free on every convex shape
	part.Start // the first point of c on the ray, and part.End where the ray leaves it
}

if hit, ok := ray.HitRectangle(r); ok { // false from an Origin inside r; Segment has the same family
	hit.Point                           // where the ray reaches r, ray.PointAt(hit.Fraction)
	ray.Direction.Reflect(hit.Normal)   // the unit normal there faces the ray, so this is the bounce
}
```

A `Line` runs without end both ways through `Origin` along `Direction`. It tells its two sides apart, and tests
a shape as the ray does, with the same `Intersection<Kind>` and `Clip<Kind>` in order along its direction:

```go
wall := geom.LineThrough(a, b) // through two points; Ln takes a direction, and s.Line() and ray.Line() convert
wall.Side(p)                   // 1 on the side the normal of the direction points to, -1 on the other, 0 on the line
wall.Nearest(p)                // the foot of the perpendicular, with DistanceTo and Contains
wall.Reflect(p)                // p mirrored across the line
wall.IntersectionLine(other)   // where two lines cross, and false for parallel ones
wall.IntersectsCircle(c)       // IntersectsLine on the shape gives the same answer
wall.ClipBox(view)             // the part of the line a viewport shows, and false where it passes by
```

### Containment

`Contains` takes a point; `Contains<Kind>` takes a whole shape, on every shape with an area but `Ellipse`:

```go
view.ContainsCircle(c) // ContainsCircle, ContainsSegment, ContainsPolygon, ContainsRectangle,
r.ContainsPolygon(p)   // ContainsRegularPolygon and ContainsBox
p.ContainsSegment(s)   // a concave polygon also checks that s does not leave it between its ends
```

### Directions, axes and orientations

```go
dir := geom.DirectionUp
dir.Turn(2)     // DirectionRight, two 45° steps
dir.Vector(5.0) // Vector{0, -5}

geom.DirectionFromAxes(up, down, left, right) // keyboard input to an 8-way direction
geom.Vec(3, -7).Direction()                   // DirectionUpRight, nearest of the eight
geom.CardinalDirections()                     // Right, Down, Left, Up
geom.CardinalOffsets[int]()                   // their lattice steps, in the same order

axis := geom.AxisVertical
axis.Along(size)             // Height, because the axis is vertical
axis.Size(length, thickness) // Size{thickness, length}

geom.RegularPolygonWithOrientation(center, size, 6, geom.OrientationPointyTop) // orientation places a vertex or an edge at the top
geom.ParseOrientation("FlatTop")                                               // the name back to the constant, "None" to OrientationNone
```

### Matrices

```go
m := geom.IdentityMatrix[float64]().Rotate(math.Pi/4).Scale(2, 2)

geom.Pt(1.0, 0.0).Transform(m) // Point{1.41, 1.41}
m.Angle()                      // π/4, read back from the matrix
m.Scaling()                    // Vector{2, 2}

geom.TRSMatrix(10.0, 20.0, math.Pi/2, 2.0, 3.0) // scale, then turn, then move: TranslationMatrix(10, 20).Rotate(π/2).Scale(2, 3)
```

`TranslationMatrix`, `RotationMatrix`, `ScaleMatrix`, `ShearMatrix` and `ReflectionMatrix` each have a composing method
on both sides.

### Numbers

The scalar helpers are exported too, generic like everything else:

```go
geom.Lerp(10.0, 20.0, 0.25)                                  // 12.5
geom.InverseLerp(10.0, 20.0, 12.5)                           // 0.25, the ratio Lerp takes
geom.Remap(0.5, -1.0, 1.0, -30.0, 30.0)                      // 15, a stick axis as a steering angle
geom.Clamp(15.0, 0.0, 10.0)                                  // 10
geom.LerpAngle(geom.ToRadians(350), geom.ToRadians(10), 0.5) // across the zero angle, the short way
```

### Number types

Any type satisfying `Number` works, including named types, and `Int()` and `Float()` convert between instantiations.
The `ints` and `floats` packages alias the two common ones:

```go
type Tile int32

geom.Pt[Tile](3, 4).Add(geom.DirectionRight.Unit[Tile]()) // Point[Tile]{4, 4}
geom.Pt(1.4, 2.6).Int()                                   // Point[int]{1, 3}, rounded
ints.Pt(1.4, 2.6)                                         // the same, from github.com/gravitton/geometry/types/ints
```

### Interop

```go
geom.RectangleFromImage[int](img.Bounds())
geom.RectangleFromMin(geom.Pt(0, 0), geom.Sz(4, 2)).Rectangle() // image.Rectangle

json.Marshal(geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4))) // {"x":1,"y":2,"w":3,"h":4}
json.Marshal(geom.DirectionUp)                        // "Up"
json.Marshal(geom.OrientationPointyTop)               // "PointyTop"
geom.ParseSize[int]("4x2")                            // Size{4, 2}
geom.ParsePoint[int]("(1,2)")                         // Point{1, 2}
```

### Testing

The `geomtest` package holds one assertion per shape, comparing with the tolerance of the asserted type, so the
main package imports no test library:

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

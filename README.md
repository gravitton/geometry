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
- **Shapes** – circle, capsule, sector, polygon, rectangle, regular polygon, box and ellipse.
- **Paths** – segment, ray, line, polyline, arc and Bezier curves.
- **Collisions** – intersections, containment, raycasts, nearest points, contacts and sweeps.
- **Interfaces** – `Shape` and `Body` hold a shape without knowing which one.
- **Screen space** – top-left origin, `+Y` down, one winding order everywhere.
- **Closed boundaries** – the edge counts as inside, within a tolerance matched to `T`.
- **Reproducible** – the same bits on amd64 and arm64, trigonometry aside.
- **Extras** – direction enums, `image` interop, JSON, parsing and `geomtest` assertions.

The rules behind them are in the [package documentation][link-go-dev-reference]. Every figure below is drawn from
the results of the calls it shows.

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
part, ok := laser.ClipRegularPolygon(rock) // Seg((32.24,0.00);(67.76,0.00)), true: part.Start is the first point hit
```

### Shapes and paths

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-light.svg">
  <img alt="The nine shapes, each labelled with its shorthand constructor: Circ, Cap, Sect, Seg, Pol, Rect, RegPol, Bx and Ell" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/shapes-light.svg">
</picture>

Every kind has a shorthand constructor in its own letters, and named constructors that read as phrases
(`CircleThrough`, `RectangleFromMinMax`, `RegularPolygonWithOrientation`):

```go
geom.Circ(center, radius)
geom.Cap(start, end, radius)            // every point within the radius of a segment
geom.Sect(center, radius, angle, sweep) // a slice of a disc, angles in radians
geom.Seg(start, end)
geom.Pol(vertices)                         // closed, convex or not
geom.Rect(center, size, angle)             // turned about its center; RectangleFromCenter has no angle
geom.RegPol(center, size, n, angle, phase) // n vertices on the ellipse of the semi-axes
geom.Bx(a, b)                              // axis-aligned, the corners in either order
geom.Ell(center, size, angle)              // the semi-axes, turned about the center
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-light.svg">
  <img alt="The paths, each labelled with its shorthand constructor: Ry, Ln, PolLn, Ac, QuadBez and CubBez" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/paths-light.svg">
</picture>

```go
geom.Ry(origin, direction)            // from Origin without end
geom.Ln(origin, direction)            // without end both ways
geom.PolLn(points)                    // an open path
geom.Ac(center, radius, angle, sweep) // a part of a circle, the sweep signed
geom.QuadBez(start, control, end)     // Bezier curves, flattened by Polyline(tolerance)
geom.CubBez(start, control1, control2, end)
```

### Points and vectors

```go
p := geom.Pt(1, 2)
v := geom.Pt(4, 6).Subtract(p) // ⟨3,4⟩

v.Length()          // 5
v.Normal()          // ⟨-4,3⟩, perpendicular
p.Add(v.Resize(10)) // (7,10)

v.Reflect(geom.Vec(0, 1))                 // ⟨3,-4⟩, bounced off a horizontal wall
p.RotateAround(geom.Pt(0, 0), math.Pi/2)  // (-2,1)
p.Lerp(geom.Pt(9, 10), 0.25)              // (3,4)
geom.Vec(8, 0).Slerp(geom.Vec(0, 8), 0.5) // ⟨6,6⟩, turned along the arc, where Lerp gives ⟨4,4⟩
```

`Lerp` and `Slerp` move by a fraction; a game loop moves by a speed times the frame time, a step that never passes
the target:

```go
p.MoveToward(geom.Pt(4, 6), 2.5)              // (3,4), and (4,6) once the step reaches it
geom.Vec(2.0, 0).RotateToward(math.Pi, 0.1)   // turned by 0.1 along the shorter arc, the length kept
velocity.MoveToward(wanted, accel*dt)         // a velocity taken toward another by an acceleration
geom.MoveTowardAngle(angle, target, speed*dt) // the same step on a bare angle
geom.MoveToward(volume, 1.0, fade*dt)         // and on any number
```

### Sizes and padding

```go
s := geom.Sz(1920, 1080)
s.Scale(0.5)                 // 960x540
s.AtMost(geom.SzU(800))      // 800x800, clamped per axis
s.Shrink(2000).AtLeastZero() // 0x0, the clamp every shape applies to its own size
s.Fit(geom.SzU(800))         // 800x450, largest with the same ratio inside
s.Fill(geom.SzU(800))        // 1422x800, smallest with the same ratio around

geom.PadXY(8, 4).Size() // 16x8, horizontal and vertical total
```

### Rectangles and boxes

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangle-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangle-light.svg">
  <img alt="A turned rectangle with its corner names, its Bounds and its BoundingCircle; the anchor of every direction on a rectangle; a rectangle clamped into another" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/rectangle-light.svg">
</picture>

A `Rectangle` is a center, a size and an angle. Turned, it keeps its center, size and corner names, and `Min`, `Max`
and `Bounds` become the box around its vertices:

```go
r := geom.RectangleFromCenter(geom.Pt(50, 50), geom.Sz(20, 10))

r.Contains(geom.Pt(55, 52))                 // true
r.Nearest(geom.Pt(80, 0))                   // (60,45), the nearest point, on every shape
r.Inset(geom.PadU(2)).Anchor(geom.TopRight) // (58,47)
r.AlignTo(geom.TopLeft, geom.Pt(0, 0))      // Rect((10,5);20x10), its Min at the origin
r.Clamp(screen)                             // moved by the least that brings it inside

d := geom.Rect(geom.Pt(0.0, 0.0), geom.Sz(2.0, 2.0), geom.Pi/4) // a diamond
d.Contains(geom.Pt(0.9, 0.9))                                   // false, outside the turned edges
d.TopLeft()                                                     // the corner that was top-left before the turn
d.Transform(geom.ShearMatrix(1.0, 0.0))                         // the nearest rectangle; Polygon() holds the parallelogram

for edge := range d.Edges() { // clockwise from TopLeft, without allocating
	edge.Midpoint()
}
```

A `Box` is the axis-aligned extent, for clipping, viewports, layout and culling, and the `Bounds` of every shape:

```go
view := geom.BoxFromSize(geom.Sz(800, 600)) // (0,0)-(800,600)
view.Contains(geom.Pt(800, 600))            // true, closed like every shape
view.Inset(geom.PadU(16))                   // (16,16)-(784,584)

hud := geom.Bx(geom.Pt(200, 40), geom.Pt(0, 0))        // (0,0)-(200,40), the corners in either order
hud.AlignTo(geom.TopRight, view.Anchor(geom.TopRight)) // (600,0)-(800,40), pinned to the corner of the view

tip := geom.BoxFromMin(geom.Pt(750, 20), geom.Sz(120, 40))
tip.Clamp(view)   // (680,20)-(800,60)
view.ClipBox(tip) // (750,20)-(800,60), true: the part on screen

r.Bounds().IntersectsBox(view) // the broad pass
view.IntersectsRectangle(r)    // and a Box is a Shape like every other
```

### Circles, capsules and ellipses

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/tangents-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/tangents-light.svg">
  <img alt="The tangents from a point to a circle, the outer and the inner tangents of two circles, and the smallest circle around both" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/tangents-light.svg">
</picture>

```go
c := geom.Circ(geom.Pt(0.0, 0.0), 5.0)
c.Anchor(geom.Bottom)          // (0,5)
c.Tangents(geom.Pt(13.0, 0.0)) // a Pair of the points the lines from there touch

wheel := geom.Circ(geom.Pt(20.0, 0.0), 2.0)
c.OuterTangents(wheel) // a belt around both, each a Segment from c to wheel
c.InnerTangents(wheel) // a crossed belt; none once the circles touch
c.Union(wheel)         // Circ((8.5,0);13.5), the smallest circle around both

geom.CircleThrough(geom.Pt(0.0, 0.0), geom.Pt(6.0, 0.0), geom.Pt(0.0, 8.0)) // Circ((3,4);5), true; false on one line
```

A circle has no `Transform` of its own, since a matrix takes it to an ellipse:

```go
e := c.Ellipse().Transform(geom.ScaleMatrix(2.0, 1.0)) // Ell((0,0);10x5), exact for every matrix
e.Contains(geom.Pt(9.0, 0.0))                          // true
e.DistanceTo(geom.Pt(0.0, 9.0))                        // 4, to the nearest point of the boundary
e.Foci()                                               // the two focal points, on the major axis
e.RegularPolygon(64, geom.OrientationPointyTop)        // the outline that draws it, every vertex on the boundary
```

A `Capsule` is the rounded collider of a character:

```go
body := geom.Cap(geom.Pt(0.0, 0.0), geom.Pt(0.0, 10.0), 2.0)
body.Contains(geom.Pt(1.5, 11.0))  // true, in the round end
body.DistanceTo(geom.Pt(6.0, 5.0)) // 4, to the side
body.Polygon(0.25)                 // the outline that draws it, within a tolerance
```

### Polygons

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-light.svg">
  <img alt="A notched polygon with its ConvexHull, its Triangulate triangles, its Grow and Shrink outlines, and its two parts after a Split" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/polygons-light.svg">
</picture>

```go
p := geom.Pol([]geom.Point[int]{{0, 0}, {4, 0}, {4, 4}, {2, 1}, {0, 4}})
p.Area()                  // 10
p.Contains(geom.Pt(2, 3)) // false, inside the notch
p.Winding()               // WindingClockwise, the winding of a Rectangle's corners
p.IsConvex()              // false
p.ConvexHull()            // Pol((0,0);(4,0);(4,4);(0,4)), the notch dropped
p.Triangulate()           // [1 2 3 0 1 3 4 0 3]: three indices into Points per triangle, an index buffer
p.Grow(1)                 // every edge moved out along its normal; Shrink moves them in

front, back := p.Split(geom.Ln(cut, dir)) // the two parts either side of a line, meeting exactly along it
both, ok := p.Union(other)                // the outline two overlapping convex polygons cover together

for vertex := range p.Vertices() { // the same loop walks a Rectangle, RegularPolygon or Box
	vertex.Float()
}
```

### Regular polygons

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-light.svg">
  <img alt="Flat-top and pointy-top hexagons with their vertices numbered and their top anchor; a stretched hexagon on its ellipse; a phase that slides the vertices and an angle that turns the ellipse" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/regular-polygons-light.svg">
</picture>

```go
hex := geom.RegularPolygonWithOrientation(geom.Pt(0, 0), geom.SzU(20), 6, geom.OrientationFlatTop)
hex.Anchor(geom.Top)         // (0,-17), the midpoint of the top edge; a pointy-top one gives its top vertex
hex.Bounds()                 // (-20,-17)-(20,17)
hex.Area()                   // 1039.23, 3√3/2 · r²
hex.Contains(geom.Pt(10, 5)) // true, walked on the edges without building the vertices
hex.Ellipse()                // the ellipse its vertices lie on

geom.RegPol(geom.Pt(0, 0), geom.Sz(20, 10), 8, 0, geom.Pi/8) // n, angle and phase
hex.Transform(geom.ShearMatrix(1.0, 0.0))                    // exact for every matrix, shears included
```

### Intersections

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-light.svg">
  <img alt="The crossing points of two circles; the parts of a segment inside a notched polygon; the two parts of a segment through the notch of a wide sector; the rectangle two turned rectangles share" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/intersections-light.svg">
</picture>

Every shape has `Intersects<Kind>`, `Contains<Kind>` and `Intersection<Kind>` for each kind, symmetric and counting a
touch. `Intersection<Kind>` gives the boundary points: one with a bool between segments, rays and lines, a `Pair` of
at most two on a convex shape, and a slice where there can be more. `Clip<Kind>` gives the part in common, in the
type of the receiver:

```go
a.IntersectsRectangle(b) // on every kind of a, for every kind of b
view.ContainsCircle(c)   // the whole circle inside
geom.Intersects(a, c)    // the same answers for two shapes of any kind

c.IntersectionCircle(other)       // a Pair, which allocates nothing
s.IntersectionPolygon(p)          // every crossing, as a slice
s.ClipPolygon(p)                  // the parts inside a concave polygon, from Start to End
beam.ClipSector(cone)             // two parts through the notch of a sector past half a turn
overlap, ok := a.ClipRectangle(b) // the overlap of two rectangles of one angle
part, ok := p.ClipBox(view)       // the polygon as a viewport shows it

buffer = s.AppendIntersectionPolygon(buffer[:0], p) // a slice result has an Append form for a reused buffer
```

### Casts

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-light.svg">
  <img alt="A ray reaching a turned rectangle, with the normal and the bounce; a circle swept into a box, stopped at the fraction of the motion; a segment hitting a circle" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/casts-light.svg">
</picture>

`Hit<Kind>` on a segment or a ray is the first point of the boundary reached from outside, with the unit normal there
and the fraction of the cast; `Raycast` and `Sweep` take a shape of any kind:

```go
ray := geom.RayThrough(camera, cursor) // from the camera through the cursor; Ry takes a direction

if hit, ok := ray.HitRectangle(r); ok { // false from an Origin inside r
	hit.Point                              // ray.PointAt(hit.Fraction)
	ray.Vector.Reflect(hit.Normal)         // the normal faces the ray, so this is the bounce
}

hit, ok := geom.Raycast(ray, shape)                  // the same for a Shape of any kind
hit, ok = geom.Sweep(ball, motion, wall)             // where ball first touches wall while it moves:
ball = ball.Translate(motion.Multiply(hit.Fraction)) // the point, the normal and the fraction of the motion
```

A `Line` runs both ways and tells its sides apart:

```go
wall := geom.LineThrough(a, b) // Ln takes a direction; s.Line() and ray.Line() convert
wall.Side(p)                   // 1 on the side the normal points to, -1 on the other, 0 on the line
wall.Reflect(p)                // p mirrored across the line
wall.ClipBox(view)             // the part a viewport shows
```

### Collisions

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-light.svg">
  <img alt="A circle and a crate overlapping a floor, with their contact points, normals and the separated position; the nearest points of two shapes and the distance between them" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/collisions-light.svg">
</picture>

```go
contact, ok := geom.Collide(crate, floor) // Normal, Depth and one or two Points
crate.Translate(contact.Separation())     // the shortest move out of the floor
on, at := geom.Nearest(turret, target)    // the point of each nearest the other
geom.Distance(turret, target)             // zero exactly where they intersect
```

`Shape` and `Body` hold a shape without knowing which one, for a spatial index or a broad pass:

```go
var s geom.Shape[float64] = c // Bounds, Contains, DistanceTo, DistanceSquaredTo, Nearest, Support
s.Support(geom.Vec(1.0, 0.0)) // the point farthest along a direction, what the functions part shapes by

var b geom.Body[float64] = c // Area, Centroid and Inertia at unit density; every shape but Segment
b.Inertia()
```

The functions take the type argument where a concrete shape is passed, `geom.Collide[float64](c, r)`, and infer it
from a value held as a `Shape`. A call through an interface allocates where the same named method does not, so a hot
loop that knows its types calls the method.

### Arcs, sectors and curves

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-light.svg">
  <img alt="An arc as a progress ring; points inside and outside a sector; a cubic Bezier flattened to a polyline; points at an even step along a polyline" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/curves-light.svg">
</picture>

```go
ring := geom.Ac(geom.Pt(0.0, 0.0), 5.0, -geom.Pi/2, 2*geom.Pi*progress) // a full turn is the circle
ring.PointAt(0.5)                                                       // also Start, End, Midpoint and Length

cone := geom.Sect(guard, 8.0, facing-geom.Pi/6, geom.Pi/3) // a vision cone
cone.Contains(player)                                      // the tolerance is on the distance, never an angle

route := geom.PolLn([]geom.Point[float64]{geom.Pt(0.0, 0.0), geom.Pt(4.0, 0.0), geom.Pt(4.0, 3.0)})
route.Length()                    // 7
route.PointAt(0.5)                // (3.5,0), half the length along
for p := range route.Every(0.5) { // a point every half unit of length, in one walk
}

curve := geom.CubBez(geom.Pt(0.0, 0.0), geom.Pt(0.0, 4.0), geom.Pt(8.0, 4.0), geom.Pt(8.0, 0.0))
curve.PointAt(0.5)   // (4,3); the parameter is not a fraction of the length
curve.Polyline(0.05) // a path within the tolerance, denser where it bends, to measure the curve by
```

### Coordinates and directions

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-light.svg">
  <img alt="The screen axes with +Y down and the sense of a positive angle; the eight directions; the clockwise winding of a rectangle's vertices" src="https://raw.githubusercontent.com/gravitton/geometry/refs/heads/main/docs/images/coordinates-light.svg">
</picture>

```go
dir := geom.DirectionUp
dir.Turn(2)                 // DirectionRight, two 45° steps
dir.Vector(5.0)             // ⟨0,-5⟩
geom.Vec(3, -7).Direction() // DirectionUpRight, nearest of the eight

geom.DirectionFromAxes(up, down, left, right) // keyboard input to an 8-way direction
geom.CardinalOffsets[int]()                   // the lattice steps of Right, Down, Left and Up

axis := geom.AxisVertical
axis.Along(size) // Height, because the axis is vertical
```

### Matrices

```go
m := geom.IdentityMatrix[float64]().Rotate(math.Pi/4).Scale(2, 2)

geom.Pt(1.0, 0.0).Transform(m) // (1.41,1.41)
m.Angle()                      // π/4, read back from the matrix
m.Scaling()                    // ⟨2,2⟩

geom.TRSMatrix(10.0, 20.0, math.Pi/2, 2.0, 3.0)  // scale, then turn, then move
geom.RotationMatrixAround(10.0, 20.0, math.Pi/2) // turn about a pivot
geom.ScaleMatrixAround(10.0, 20.0, 2.0, 3.0)     // scale about a pivot, which stays where it is
```

`Translate`, `Rotate`, `RotateAround`, `Scale`, `ScaleAround`, `Shear` and `Reflect` compose with their constructor on
the right, applied first, and their `Pre` forms on the left, applied last.

### Numbers

```go
geom.Lerp(10.0, 20.0, 0.25)                                  // 12.5
geom.InverseLerp(10.0, 20.0, 12.5)                           // 0.25, the ratio Lerp takes
geom.Remap(0.5, -1.0, 1.0, -30.0, 30.0)                      // 15, a stick axis as a steering angle
geom.Clamp(15.0, 0.0, 10.0)                                  // 10
geom.LerpAngle(geom.ToRadians(350), geom.ToRadians(10), 0.5) // across the zero angle, the short way
```

Any type satisfying `Number` works, named types included, and `Int()` and `Float()` convert between them. The
`ints` and `floats` packages alias the two common ones:

```go
type Tile int32

geom.Pt[Tile](3, 4).Add(geom.DirectionRight.Unit[Tile]()) // (4,4)
geom.Pt(1.4, 2.6).Int()                                   // (1,3), rounded
ints.Pt(1.4, 2.6)                                         // the same, from github.com/gravitton/geometry/types/ints
```

### Interop

```go
geom.RectangleFromImage[int](img.Bounds())
geom.BoxFromImage[int](img.Bounds())                            // Min and Max as they are, an odd extent kept
geom.RectangleFromMin(geom.Pt(0, 0), geom.Sz(4, 2)).Rectangle() // image.Rectangle

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

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project follows [Semantic Versioning](http://semver.org/spec/v2.0.0.html) with one deliberate exception:

**a breaking change may ship in a minor release.** Breaking changes are avoided where the cost is reasonable, but a
Go major version changes the import path for every user, and for a package with a small user base that is often the
more expensive of the two. A breaking change is marked **breaking** in its entry, and every release that has one
lists them under **Breaking** at the top of its section. Renames land as a rename; no deprecated alias is kept.


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.16.0...main)

### Added
- `Line` – the infinite line through `Origin` along `Direction`, with `Ln`, `LineThrough`, `Segment.Line` and `Ray.Line`, no `Shape` since it has no `Bounds`: `Side` gives the exact sign of the side a point lies on, `Reflect` mirrors a point across it, `Contains`, `DistanceTo`, `DistanceSquaredTo` and `Nearest` answer as a ray from its origin toward the point does, and it moves, turns and scales as `Ray` does; it has `Intersects<Kind>`, `Intersection<Kind>` and `Clip<Kind>` for every kind, with `AppendIntersectionPolygon` and `AppendClipPolygon`, decided on the segment tests like the ray's and returned in order along the direction, and `Circle`, `Segment`, `Ray`, `Polygon`, `Rectangle`, `RegularPolygon` and `Box` have `IntersectsLine` and `IntersectionLine`; two lines cross at one point, the same whichever asks and false for parallel lines
- `Pair` – a result of at most two items held in the value, `Count` of them in `Items`, ranged with `All` or appended to a slice with `AppendTo`, with `IsZero`, `IsEmpty`, an `Equal` that compares the items by their own, a `String` in the form `Pair(item;item)` and JSON as the array of its items
- `Polygon.ClipPolygon`, `ClipRectangle`, `ClipRegularPolygon` and `ClipBox`, each with its `Append` form – the part of the polygon inside a convex region, and false where there is none: the polygon is cut to the inner side of each edge of the region in turn, so it may be concave, a part of several pieces coming back as one outline joined along the cuts; `ClipPolygon` takes whichever of the two is convex as the region and is false for two concave polygons; a region of no area cuts everything away, and for an integer type the crossings are rounded at each cut
- `Polygon.Grow` and `Polygon.Shrink`, with `AppendGrow` and `AppendShrink` – the polygon with every edge moved out or in along its normal by a distance, the corners staying sharp; shrinking a convex polygon drops an edge the others close over and leaves an empty polygon past the middle, and a polygon that is not convex has each vertex moved along its corner, its edges crossing once the distance passes its narrowest notch or neck
- `Support` on `Circle`, `Segment`, `Polygon`, `Rectangle`, `RegularPolygon`, `Box` and `Ellipse` – the point farthest along a direction of any length, in `float64`: a vertex as `Vertices` places it, and the point of the boundary on a circle and an ellipse
- `Distance` and `DistanceSquared` – the distance between the nearest points of two shapes of any kind, and its square for comparisons, zero exactly where `Intersects` holds and measured as it tests the pair, so a shape in a notch is measured to the notch; two circles and two boxes have it as methods too, `Circle.DistanceToCircle`, `DistanceSquaredToCircle`, `Box.DistanceToBox` and `DistanceSquaredToBox`, the distance between the centers less the two radii and the gap on each axis, which `Distance` reads as it does `Segment.DistanceSquaredToSegment`
- `Nearest` over two shapes – the two points `Distance` is measured between, the point of each shape nearest the other, found as `DistanceSquared` measures the pair, so a shape in a notch is measured to the notch; two shapes that intersect give one point they share, and an empty polygon the zero point
- `Hit` and `HitCircle`, `HitSegment`, `HitRay`, `HitLine`, `HitPolygon`, `HitRectangle`, `HitRegularPolygon` and `HitBox` on `Segment` and `Ray` – where a cast first reaches the boundary of a shape from outside: the point, the `Start` of the part `Clip<Kind>` gives a cast from outside, the unit normal of the boundary there turned to face the cast, in `float64`, and the fraction of the segment or of `Direction` at which `PointAt` reaches it; a `Start` or `Origin` the shape contains hits nothing, and none allocates; a segment, a ray or a line is a wall with no inside, hit wherever the cast intersects it from a start off it, and a wall along the cast is met end on, the normal back along the cast; a cast that touches a corner in passing takes the normal of the edge facing its start, and a corner met head on that of the edge first in `Edges`
- `Raycast` – the `Hit<Kind>` of a ray against a shape of any kind, the one call a ray makes against a mixed list; an `Ellipse` and a shape of another package are reached as `Sweep` reaches a shape, to within the tolerance
- `Sweep` – where a shape first touches another while it moves by a vector, as a `Hit` whose fraction is of the motion, so a body moving farther in a step than it is wide stops at a wall; shapes that intersect before the move hit nothing, and a polygon that is not convex is swept edge by edge, into its notches
- `Collide` and `Contact` – how two overlapping shapes meet, what a collision response and a physics solver resolve a collision from, and false exactly where `Intersects` is: the unit normal that parts them, the depth along it, and the points of the second shape they meet at, two where flat sides lie against each other and one at a corner or a round side; `Contact.Separation` is the shortest vector that moves the first out of the second; shapes that only touch still have a normal, two shapes given the other way round answer with the opposite normal, and a polygon that is not convex meets edge by edge, the deepest contact with an edge the other shape crosses, so a shape in a notch parts from a side of the notch and one in an inside corner from one wall at a time
- `MoveToward`, `Point.MoveToward`, `Vector.MoveToward`, `Vector.RotateToward` and `MoveTowardAngle` – a number, a point or a vector moved, a vector turned and an angle turned towards a target by at most a step, the move of a game loop by a speed times the frame time, where `Lerp` and `Slerp` move by a fraction: none passes the target, a step that reaches it gives the target itself, the turns take the shorter arc, by increasing angle where the two are half a turn apart, and a negative step moves away without being clamped
- `InverseLerp` and `Remap` – the ratio at which a value lies between two others, the one `Lerp` takes to give it back, 0 for an empty range, and a value taken from one range into another at that ratio, either range running up or down, rounded once for an integer type; neither clamps
- `TRSMatrix` – the matrix that scales, then rotates, then translates, built in one call from two deltas, an angle and two factors: the product `TranslationMatrix(...).Rotate(...).Scale(...)` composes, with the order written once, and for an integer type each entry rounded once
- `Segment.IsPoint` – whether the start and the end are equal, a segment of no length, the one point a `Clip<Kind>` gives where a cast touches a shape in passing
- `Size.IsPositive` – whether the width and the height are both greater than zero, a size that encloses an area, exactly and with no tolerance, so a size too small for `IsZero` to tell from zero still is and one with a negative or a NaN extent is not
- `IsFinite`, `Point.IsFinite`, `Vector.IsFinite` and `Size.IsFinite` – whether a number, or both components, is neither NaN nor infinite, the test to make before `Int` or `Cast` to an integer type, which panic on any other value
- `BoundingCircle` on `Segment`, `Rectangle`, `Circle`, `Ellipse`, `Polygon`, `RegularPolygon` and `Box` – the circle around the shape, a broad pass by one squared distance that stays the same as the shape turns: about the shape's own center, the `Centroid` of a `Polygon`, to its farthest vertex; for an integer type the radius is rounded up, so the circle contains every vertex
- `CircleBetween` and `CircleThrough` – the circle with two points at the ends of a diameter, and the circle through three points, false where they lie on one line within the tolerance
- `Circle.Tangents`, `Circle.OuterTangents` and `Circle.InnerTangents`, each a `Pair` – the points where the lines through a point touch the circle, two for a point outside, mirrored across the line to the center, the point itself on the boundary within the tolerance, and none inside; and the segments of the lines that touch two circles, from the same side or crossing between them, each from its point on the circle to its point on the given one, two, one between two circles of zero radius, and none where one circle lies within the other (outer) or the circles intersect (inner)
- `Line`, `Ln`, `Hit` and `Contact` in `types/ints` and `types/floats`
- `geomtest.AssertLine`, `geomtest.AssertHit` and `geomtest.AssertContact` – a `Line`, a `Hit` and a `Contact` compared field by field, the point and the points with the tolerance of their type and the normal, the fraction and the depth as the `float64` they are

### Changed
- `Multiply` and `Divide` on a number are inlined on arm64, so the scaling of a `Point`, a `Vector` or a `Size` is one call where it was three, about half the time of a chain of them; `Cast` tests for a non-finite value with one comparison
- **breaking** `Shape` requires `Support`, and `Intersects` takes two `Shape` and never panics: a pair with no `Intersects<Kind>` method for each other, an `Ellipse` or a shape of another package among them, is tested by its supports, a circle by its radius and a polygon that is not convex edge by edge, so a notch is not filled in; a `Ray` and a `Line` are no longer passed to it and test a shape by their own `Intersects<Kind>`
- **breaking** Every `Intersection<Kind>` that returned a slice of at most two points returns a `Pair` of them and allocates nothing: `IntersectionCircle` on `Circle`, `Segment` and `Ray`, `IntersectionRectangle`, `IntersectionRegularPolygon` and `IntersectionBox` on `Segment` and `Ray`, and `IntersectionSegment` and `IntersectionRay` on `Circle`, `Rectangle`, `RegularPolygon` and `Box`; `pair.Items[:pair.Count]` is the slice they returned, and `IntersectionPolygon` stays a slice, since a concave outline is crossed any number of times
- **breaking** `Rectangle.IntersectionRectangle` renamed to `Rectangle.ClipRectangle` and `Box.IntersectionBox` to `Box.ClipBox`: an `Intersection<Kind>` is the points where two boundaries meet on every shape, and a `Clip<Kind>` the part the two have in common, in the type of the receiver
- **breaking** `EnclosesCircle`, `EnclosesSegment`, `EnclosesPolygon`, `EnclosesRectangle`, `EnclosesRegularPolygon` and `EnclosesBox` on `Circle`, `Polygon`, `Rectangle`, `RegularPolygon` and `Box` renamed to `Contains<Kind>`, so the test of a whole shape reads as `Contains` of a point does

### Removed
- **breaking** `Collider`, for the `Shape` that `Intersects` takes
- **breaking** Every `AppendIntersection<Kind>` but `AppendIntersectionPolygon` on `Segment` and `Ray`, which appends what the others delegated to it: `AppendIntersectionCircle`, `AppendIntersectionSegment` and `AppendIntersectionRay` on `Circle`, `AppendIntersectionCircle`, `AppendIntersectionRectangle`, `AppendIntersectionRegularPolygon` and `AppendIntersectionBox` on `Segment` and `Ray`, and `AppendIntersectionSegment` and `AppendIntersectionRay` on `Polygon`, `Rectangle`, `RegularPolygon` and `Box`

## [v1.16.0](https://github.com/gravitton/geometry/compare/v1.15.0...v1.16.0) (2026-10-07)

### Added
- `RoundHalfUp` – the nearest integer with a tie rounded toward positive infinity, where `Round` rounds it away from zero, so rounding commutes with a shift by an integer: the rounding a lattice needs for every cell to own the same edges
- `DirectionOffsets`, `CardinalOffsets` and `DiagonalOffsets` – the lattice steps of `Directions`, `CardinalDirections` and `DiagonalDirections` in the same order, as `Vector[T]` for any `T`, so a step and its direction share an index
- `Polygon.IsSimple` – whether the outline goes around without touching itself: no two edges meet but neighbours at their shared vertex, and no edge doubles back

### Changed
- The `geomtest` helpers take `geomtest.Testing`, declared in the package with the `Helper` and `Errorf` methods they call, rather than `assert.Testing`, so no signature names the assertion library; a `*testing.T` and any `assert.Testing` still pass as they are
- `String` on every type and the package `String` append into a buffer on the stack rather than going through `fmt.Sprintf`, so each allocates the string it returns and nothing else, a polygon of many vertices its growing buffer besides; the text printed is unchanged
### Removed
- **breaking** `Triangle`, `Square` and `Hexagon` – `RegularPolygonWithOrientation` with 3, 4 or 6 vertices; `Square` of unequal semi-axes drew a rhombus or a rectangle rather than a square, and the three took generic names at package level

### Fixed
- `Circle.IntersectionSegment`, `Segment.ClipCircle` and the ray forms gave a segment and its `Reverse` different crossings where an endpoint on the boundary lay inside the chord
- `Segment.ClipPolygon` and `Ray.ClipPolygon` lost the inside part where the entry crossing rounded onto a `Start` just outside the polygon
- `Matrix.Scaling` and `Matrix.Angle` gave zero for a matrix collapsing only the X axis, so `Rectangle.Transform` shrank to a point
- `Ellipse.Foci` and `Ellipse.Nearest` panicked for an integer ellipse with a negative semi-axis
- `Ellipse.Contains` rejected an integer point lying exactly on the boundary where dividing by the semi-axes rounded the quadratic form above one
- For an integer T past 2^17, a segment and its `Reverse` rounded a crossing on a half unit to different neighbours in every `Intersection` and `Clip` form
- `Segment.IntersectionRectangle`, `IntersectionRegularPolygon` and `IntersectionBox` and their ray forms could list three points where an endpoint lay within the band of a long edge just past a corner; they now return the ends their `Clip` form clips to, never more than two
- A ray could end its reach beside a vertex and list a touch there that no longer segment along it gave; the reach now ends past the shape


## [v1.15.0](https://github.com/gravitton/geometry/compare/v1.14.0...v1.15.0) (2026-10-04)

Two new kinds and the queries between shapes. `Box` is the axis-aligned box every `Bounds` now returns, and `Ray` the half-line cast through every collider; shapes enclose, clip and find their nearest point, and every slice result has an `Append` form.
The boundary tolerance widens with the coordinates, the same inputs give the same bits on every architecture, and a regular polygon takes a phase that keeps its orientation and every transform exact.

Every entry, with breaking changes marked, is in [docs/releases/v1.15.0.md](docs/releases/v1.15.0.md).

### Breaking
- `Bounds` returns a `Box` on every shape and in `Shape`; `Box.Rectangle` converts it back
- `Collider` requires `IntersectsBox` and `IntersectsRay`
- `Outline`, `Transformable` and `Rectangle.MinMaxString` are removed
- `Rectangle.Clamp(point)` is `Nearest`, read off the edge walk
- `RegPol` takes a phase, and the orientation constructors set it rather than the angle; `RegularPolygonOrientationAngle` is `RegularPolygonOrientationPhase`
- `Ellipse.RegularPolygon` takes an `Orientation`, as `Circle.RegularPolygon` does
- The `Assert*` helpers moved to `geomtest`, taking `assert.Testing`; `geom.Testing` is removed

### Added
- `Box` – the axis-aligned box as `Min` and `Max`, a `Shape` and a `Collider`, with `Bx`
- `Ray` – a half-line from `Origin` along `Direction`, a `Collider` deciding every pair on the `Segment` tests
- `Clip` on `Segment` and `Ray` against every bounded collider, allocation-free on the convex shapes
- `Encloses` on every container for every bounded kind but `Ellipse`
- `Nearest(point)` on every shape, and `Rectangle.Clamp(rectangle)` at any two angles
- `Polygon.Winding`, `IsConvex`, `ConvexHull`, `Simplify` and `Lerp`
- `Append` forms of every slice result, so a loop reusing its buffer allocates nothing
- `Vector.Slerp` and `ParsePoint`
- `RegularPolygon.Phase`, placing the vertices before the stretch, so a flat-top polygon of unequal semi-axes is stretched across its edges

### Changed
- The boundary tolerance widens to two ulps of `T` far from the origin, so a rounded float32 point stays on its boundary
- `Point.Transform` and `Vector.Transform` are about 7% faster, and parse errors start with `geom:`
- `github.com/gravitton/x` is dropped; `github.com/gravitton/assert` is the only dependency, at v1.6.0
- `RegularPolygon.Transform` and `Ellipse.Transform` are exact for every matrix, shears included
- The package doc states the supported range of each `T`: `int8` and `int16` never panic or hang but may wrap

### Fixed
- Every product is rounded before it is added, so arm64 and amd64 v3 no longer fuse it into a multiply-add, and CI fails on a new one
- Nearly collinear float segments no longer cross off both of them
- Float32 and narrow integer shapes convert before subtracting, so far coordinates measure true
- Narrow integer crossings are ordered from `Start`, and `Mod` no longer overflows near the end of the range
- Long rectangles turned a hair apart are no longer tested as parallel
- Integer regular polygons contain their own vertices, and polygons on one line enclose no area
- A segment along an edge or past a vertex crosses at most twice, and a crossing on a half unit rounds alike at any length and from either side
- An integer circle tangent to an edge intersects it across the supported range
- Rectangles half a turn apart share a frame, and a pair of one angle answers alike in either order


## [v1.14.0](https://github.com/gravitton/geometry/compare/v1.13.0...v1.14.0) (2026-09-20)

Oriented shapes and one normalized surface. `Rectangle` carries an `Angle` that turns it about its center, `RegularPolygon.Angle` turns the polygon rather than phasing its vertices, and `Ellipse` is the continuous shape an affine matrix makes of a circle.
Every shape pair has one name and one holder, every closed shape walks the same edges, the outlines are `iter.Seq`, and five interfaces name what the shapes share.

Every entry, with breaking changes marked, is in [docs/releases/v1.14.0.md](docs/releases/v1.14.0.md).

### Breaking
- A rotated `Rectangle` answers from its vertices: `Min`, `Max`, `MinMax` and `Bounds` are their extent, `Clamp` snaps to the nearest turned edge, and `DistanceSquaredTo` returns `float64`
- Rectangles of different angles have no common rectangle, so `IntersectionRectangle` declines and `Union` gives the box around both `Bounds`
- `RegularPolygon.Angle` turns the polygon instead of phasing its vertices along a fixed ellipse
- `Vertices` and `Edges` return an `iter.Seq` instead of a slice
- `Line` is `Segment`, the same-kind `Intersects` and `Intersection` carry the shape name, `Circle.Circumference` is `Perimeter`, and `Polygon.Vertices` is the `Points` field
- A size is signed throughout: `Size.Grow` and `Shrink` no longer clamp, and `Vector.Size` keeps the sign of each component
- `Segment.MoveTo` places the midpoint, and `Rectangle.Rectangle` spans exactly `Size.Int` pixels
- `String` and JSON change form for `Polygon` and a turned `RegularPolygon`
- `Rectangle.Area` returns `float64`, like every product, so the rectangle is a `Body`
- Names settle on one vocabulary: `Empty` is `IsEmpty`, `FlatTop` is `OrientationFlatTop`, `Segment.Lerp` is `PointAt`, `Axis.Cross` is `Perpendicular`, `Direction.Rotate` is `Turn` and `Polygon.Center` is `Centroid`
- `Rectangle` keeps its corner getters alone: the edge midpoints `Top` and its siblings and the edge getters `TopEdge` and its siblings are gone, since `Anchor(Top)` names the point on every shape and `Edges` yields the edges

### Added
- `Ellipse` – the continuous shape of two semi-axes and an angle, with the foci, the Ramanujan perimeter and a boundary distance found by bisection
- `Rectangle.Angle`, `Rotate` and `IsAligned` – the oriented box, with corners, edges, anchors and `Inset` named in the frame before the turn
- The conversions between a curved boundary and the polygon of the wanted resolution, in both directions and exact
- `Shape`, `Outline`, `Collider`, `Body` and `Transformable`, with the `Intersects(a, b)` package function over two colliders and `Centroid` with `Inertia` on every shape that encloses an area
- `Transform` on `Rectangle` and `RegularPolygon`, `Rotate` on `Circle`, and `Cast[R]()` on every type
- `Intersects` between every remaining pair of shapes, `RegularPolygon` included on both sides
- `Canonical` snapping an angle residue, the `Orientation` text encoding, and `AssertEllipse` with `AssertAngle`
- `RegularPolygon.IsAligned`, the exact-zero angle test `Rectangle` and `Ellipse` already have
- `RegularPolygon.Anchor` and `AlignTo` on the point where the ray leaves the polygon, with `Resize`, `Grow` and `Shrink` as on `Ellipse`, and `Size.AtLeastZero`, the clamp behind every `Grow`

### Changed
- Every boundary is one squared-distance comparison, and every closed shape walks its own `Edges` through the shared `edge.go` accumulators
- Every shape pair is held by one shape and delegated to from the other, in the order `Circle`, `Segment`, `Polygon`, `Rectangle`, `RegularPolygon`
- Only a method whose result is a slice allocates: the outline iterators and the intersection walks allocate nothing, and `Circle.IntersectionSegment` allocates its result once
- Every earlier release with a breaking change opens with a **Breaking** group

### Fixed
- A point with a NaN coordinate is inside no shape, where it was inside every rectangle and polygon


## [v1.13.0](https://github.com/gravitton/geometry/compare/v1.12.0...v1.13.0) (2026-09-18)

One rule for every boundary: closed within `Epsilon[T]()`, `DistanceTo` zero exactly where `Contains` holds, and collision as methods on the shapes rather than free functions.
Arithmetic computes in `float64`, no shape stores a negative extent, and an input with no meaning panics.

Every entry, with breaking changes marked, is in [docs/releases/v1.13.0.md](docs/releases/v1.13.0.md).

### Breaking
- Boundaries are closed, so touching shapes collide and an edge is inside
- No negative size or radius; the empty circle of a negative radius is gone
- `float64` arithmetic rounded back into `T` changes narrow integer and large `int64` results
- A non-finite `Cast`, a zero divisor and a singular `Inverse` panic
- `Collision*` become methods, `AssertRect` is `AssertRectangle`, and `Assert*` take a shape value
- `String`, `Line` JSON keys and the `Direction` and `Axis` encodings change form
- `Polygon.Center` is the area centroid
- `Directions` and `Axes` are functions; `Normalize` and `Direction.Unit` give an axis-aligned unit vector
- Fixed answers change: `Circle.Bounds`, `Rectangle.Inset`, `FlatTop` and `Matrix.Unscale`

### Added
- `Intersects` and `Intersection` between shapes, and `Contains`, `DistanceTo` and `DistanceSquaredTo` on every shape
- `Polygon` measurement and transforms: `Area`, `Perimeter`, `Edges`, `Bounds`, `Rotate`, `Transform`
- `Unscale`, `Lerp`, `Canonical` and `AlignTo` across the shapes, with gaps filled on `Line`, `Vector` and `Size`
- `ShearMatrix`, `ReflectionMatrix` and the matrix queries
- Tolerance, angle and parsing helpers in `math.go`

### Changed
- `RegularPolygon` measures without building its vertices
- Every shape's `String` and JSON follow one form

### Fixed
- `Circle.Bounds`, `Rectangle.Inset`, `Axis.Project`, `Vector.Normalize`, `NormalizeAngle`, `Matrix.Unscale`, `Polygon.UnmarshalJSON` and `Cast` on defined integer types


## [v1.12.0](https://github.com/gravitton/geometry/compare/v1.11.0...v1.12.0) (2026-08-26)
### Added
- `Point.Compare(point) int` – orders points by X and then by Y in the `cmp.Compare` convention, for `slices.SortFunc` and `slices.BinarySearchFunc`. The comparison is exact and does not apply the `Delta` tolerance `Point.Equal` uses, since a tolerant comparison is not transitive and would leave a sort no valid ordering


## [v1.11.0](https://github.com/gravitton/geometry/compare/v1.10.0...v1.11.0) (2026-08-26)
### Breaking
- `Direction` constants are renumbered to run by increasing angle — `Right`, `DownRight`, `Down`, `DownLeft`, `Left`, `UpLeft`, `Up`, `UpRight` — so a positive `Rotate` step turns the same way as a positive `Vector.Rotate` angle: counterclockwise in math coordinates, clockwise as drawn on a screen with Y pointing down. Previously the order ran the opposite way and `DirectionFromAngle` negated its input to compensate (**breaking**)
- `Directions`, `CardinalDirections`, and `DiagonalDirections` follow the new order (**breaking**)
- `Rectangle.Vertices` and `Rectangle.Edges` now wind by increasing angle, matching `Directions` and `RegularPolygon.Vertices` — clockwise as drawn on a screen with Y pointing down. `Rectangle.Polygon()` winds the same way as a result (**breaking**)
- `Rectangle.TopEdge`, `RightEdge`, `BottomEdge`, and `LeftEdge` are reversed so the edges still form a closed chain in the new winding (**breaking**)

### Changed
- `DirectionFromAngle` no longer negates its argument, now that direction numbering follows the angle directly

### Fixed
- `Direction.Rotate` documentation had the screen sense inverted — it claimed positive steps appear clockwise while the old ordering made them appear counterclockwise
- `RegularPolygon.Vertices` documentation claimed counter-clockwise winding; it winds by increasing angle, which is clockwise as drawn
- `Direction.Angle`, `Vector.Rotate`, `VectorFromAngle`, `Rectangle.Vertices`, and `Rectangle.Edges` now state which convention their angles and winding use


## [v1.10.0](https://github.com/gravitton/geometry/compare/v1.9.0...v1.10.0) (2026-08-26)
### Breaking
- `Direction[T](x T) T` renamed to `Sign[T](x T) T`, freeing the name for the new `Direction` type (**breaking**)
- `Line.Direction() Vector[T]` renamed to `Line.Vector()`, since it returns the start-to-end vector rather than a `Direction` (**breaking**)
- `Line.Reversed()` renamed to `Line.Reverse()`, matching the imperative naming of every other immutable method (**breaking**)
- `DirectionFromAxes(up, down, left, right bool)` now returns `Direction` instead of `Vector[float64]` (**breaking**)
- `RectFromMin`, `RectFromMax`, `RectFromMinMax`, `RectFromSize`, `RectFromImage` renamed to `RectangleFrom…`, so every `From` constructor uses its full type name like `PointFromImage` and `VectorFromAngle` (**breaking**)
- `UpVector`, `DownVector`, `LeftVector`, `RightVector`, `UpLeftVector`, `UpRightVector`, `DownLeftVector`, `DownRightVector` removed – superseded by `DirectionX.Unit[T]()` (**breaking**)

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
- `Assert*` helpers take `assert.Testing` instead of `*testing.T`, matching the interface `gravitton/assert` already uses; existing calls passing `*testing.T` are unaffected

### Fixed
- `Matrix.Int()` now rounds to the nearest integer via `Cast` instead of truncating toward zero, matching every other `Int()` conversion
- `RegularPolygon.Vertices` now rounds after scaling by the size instead of rounding a unit vector first, so integer vertices land within half a unit of the exact position (previously up to 1.5 units off at non-right angles)


## [v1.9.0](https://github.com/gravitton/geometry/compare/v1.8.1...v1.9.0) (2026-05-29)
### Added
- `Parse[T Number](s string) (T, error)` – parses a string into any `Number` type using `strconv.ParseInt` (with correct bit size) for integer types and `strconv.ParseFloat` for float types
- `ParseSize[T Number](s string) (Size[T], error)` – parses a size string in `"WxH"` form into `Size[T]`


## [v1.8.1](https://github.com/gravitton/geometry/compare/v1.8.0...v1.8.1) (2026-05-22)
### Fixed
- `RegularPolygonOrientationAngle` with `PointyTop` now returns `-π/2` instead of `π/2`, correctly placing the first vertex at the visual top in the library's +Y-down coordinate system


## [v1.8.0](https://github.com/gravitton/geometry/compare/v1.7.0...v1.8.0) (2026-05-12)
### Added
- `DirectionFromAxes(up, down, left, right bool) Vector[float64]` – returns the unit direction vector for the given axis inputs; cancels opposite directions, normalizes diagonals
- `Line.Vertices() []Point[T]` – returns `[Start, End]` as a slice


## [v1.7.0](https://github.com/gravitton/geometry/compare/v1.6.0...v1.7.0) (2026-05-09)
### Added
- `RectFromMax` – constructor for `Rectangle` from max (bottom-right) point and size
- `Direction[T](x T) T` – returns the sign of a number: `1`, `-1`, or `0`
- `Vector.Round() Vector[T]` – returns a new vector with each component rounded to the nearest integer
- `Vector.Floor() Vector[T]` – returns a new vector with each component rounded down
- `Vector.Ceil() Vector[T]` – returns a new vector with each component rounded up


## [v1.6.0](https://github.com/gravitton/geometry/compare/v1.5.0...v1.6.0) (2026-05-03)
### Breaking
- `Matrix[T Number]` — `Matrix` is now generic over the full `Number` constraint (was restricted to `float64`), so every use of the bare type name takes a type argument (**breaking**)

### Added
- `RegPolWithOrientation` – constructor for `RegularPolygon` with orientation


## [v1.5.0](https://github.com/gravitton/geometry/compare/v1.4.0...v1.5.0) (2026-04-25)
### Added
- `Point.ChebyshevDistanceTo(point Point[T]) T` — returns the Chebyshev (chessboard) distance: `max(|dx|, |dy|)`
- `Point.OctileDistanceTo(point Point[T]) float64` — returns the Octile distance for grid pathfinding where diagonals cost √2: `max(dx, dy) + (√2-1) * min(dx, dy)`
- `Round[T]`, `Floor[T]`, `Ceil[T]` generic free functions in `math.go` (same pattern as `Abs[T]`)

### Changed
- `Point.Round`, `Point.Floor`, `Point.Ceil` now delegate to the new `Round`, `Floor`, `Ceil` free functions


## [v1.4.0](https://github.com/gravitton/geometry/compare/v1.3.0...v1.4.0) (2026-04-25)
### Added
- `Point.ManhattanDistanceTo(point Point[T]) T` — returns the Manhattan (taxicab) distance between two points


## [v1.3.0](https://github.com/gravitton/geometry/compare/v1.2.0...v1.3.0) (2026-04-25)
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


## [v1.2.0](https://github.com/gravitton/geometry/compare/v1.1.1...v1.2.0) (2026-04-19)
### Breaking
- Renamed `PointTop` → `PointyTop` in `RegularPolygon` orientation constants (**breaking**)
- `Shrink` and `ShrinkXY` now clamp to `0` — negative dimensions are no longer possible (**breaking**)

### Added
- `Point.Abs()`, `Point.Round()`, `Point.Floor()`, `Point.Ceil()` methods
- `Size.Unscale(factor)`, `Size.UnscaleXY(factorX, factorY)`, `Size.AtLeast(size)`, `Size.AtMost(size)` methods
- `RegularPolygon.Translate(vector)` method
- `Pol[T]` shorthand constructor for `Polygon`
- `Sqrt3` and `OneOverSqrt2` constants
- `Polygon.UnmarshalJSON` support

### Changed
- Updated minimum Go version to 1.26

### Fixed
- `PadXY` parameter order corrected
- Logic errors and typos in comments across multiple files


## [v1.1.1](https://github.com/gravitton/geometry/compare/v1.1.0...v1.1.1) (2025-10-27)
### Fixed
- Move Assert helper methods to `geom` package


## [v1.1.0](https://github.com/gravitton/geometry/compare/v1.0.0...v1.1.0) (2025-10-27)
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

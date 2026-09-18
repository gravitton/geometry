# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/geometry/compare/v1.12.0...main)
### Added
- `Int[T Number](value T) int` – converts a `Number` to `int`, an integer `T` by plain conversion and a float `T` rounded through `Cast`; the `Int()` methods are built on it
- `Delta32` – the equality tolerance for `float32`, `1e-4`; `Delta` sat below one `float32` ulp above magnitude 1e3 and asked for bit-exact equality there
- `Epsilon[T]()` – the tolerance for `T`: zero for an integer `T`, `Delta32` for `float32`, `Delta` for `float64`; it depends only on the type, so `Equal` costs what it did before
- `EqualRelative(a, b)` and `EpsilonRelative(a, b)` – equality within a tolerance that scales with magnitude, for values where an absolute tolerance falls below one ulp; ~1.1 ns against ~0.6 ns for `Equal`, which is why it is a separate function
- `LessOrEqual(a, b)` and `LessOrEqualDelta(a, b, delta)` – the `<=` counterparts of `Equal` and `EqualDelta`; every closed boundary check is built on them
- `Sum(values)` – adds a slice of `Number`, accumulating in `float64` and storing the total through `Cast`
- `AngleDistance(a, b)` – the shortest angular distance between two angles, in `[0, π]`
- `EqualAngle(a, b)` – angle equality modulo a full turn within `Delta`, holding across the `0`/`2π` seam
- `AssertNumber(t, actual, expected)` – asserts a bare coordinate with the same integer/float rule as the shape helpers
- `Vector.LessOrEqual(length)` – the closed counterpart of `Less`
- `Padding.Equal`, `IsZero`, `Add`, `Negate`, `Scale` and `Unscale`, so `Padding` has the same equality pair and arithmetic as the other types
- `Rectangle.Outset(padding)` – expands the rectangle by the padding, the inverse of `Inset`
- `Rectangle.MinMax()`, `Line.MinMax()` and `Polygon.MinMax()` – both corners of the bounding box at once, exact for an integer `T` where `Bounds` places a center; the intersection tests reject shapes by them before examining any edge
- `Point.Between(a, b)` – reports whether the point lies within the box spanned by the corners, closed within `Epsilon[T]()`; the check `Rectangle.Contains` and `Polygon.Contains` make
- `Line.Intersects(line)`, `IntersectsCircle`, `IntersectsRectangle` and `IntersectsPolygon` – segment intersection, closed within `Epsilon[T]()` like `Contains`; touching endpoints, tangents and collinear overlaps intersect. `IntersectsCircle` compares squared distances and `IntersectsRectangle` rejects a segment by its extent before walking the edges, so a segment apart from the rectangle costs 6 ns instead of 95 ns
- `Polygon.Intersects(polygon)`, `IntersectsLine`, `IntersectsRectangle` and `IntersectsCircle` – a vertex of one inside the other or crossing edges, by the even-odd rule `Contains` uses; an empty polygon intersects nothing. Every test rejects the other shape by its extent first and each edge pair by their extents, and `IntersectsRectangle` walks the rectangle edges without building its polygon: two apart polygons within a shared extent went from 81 µs to 2 µs
- `Rectangle.IntersectsLine`, `IntersectsPolygon`, `Circle.IntersectsLine` and `IntersectsPolygon` – the mirrors of the above, so every pair of shapes has an intersection test on both sides
- `Line.Intersection(line)` – the point where two segments cross and whether they do, answering exactly where `Intersects` holds less the parallel case: a proper crossing is decided on cross products that are exact for an integer `T`, and a touch on the endpoint distance like `Contains`, so a shallow touch is not lost to the fraction along the segment; parallel segments have no single point and report false even where they overlap, while a zero-length segment is parallel to nothing and is the point where it lies on the other
- `Rectangle.Intersection(rectangle)` – the common rectangle and whether there is one; touching rectangles give a zero-width or zero-height one
- `Rectangle.Union(rectangle)` – the smallest rectangle containing both
- `Circle.DistanceTo(point)`, `Rectangle.DistanceTo(point)` and `Polygon.DistanceTo(point)` – the distance to the nearest point of the shape, zero exactly where `Contains` holds, the tolerance band included, so every shape measures distance the way `Line` does; an empty polygon is infinitely far
- `Line.Transform(matrix)` and `Polygon.Transform(matrix)` – apply a float matrix to every point, like `Point.Transform`
- `Line.Scale(factor)` and `ScaleXY(factorX, factorY)` – scaling about the midpoint, so every shape scales about its own center; an integer `T` rounds the midpoint first like `Rotate`
- `Line.Rotate(angle)` and `Polygon.Rotate(angle)` – rotation about the midpoint and the centroid in the sense of `Vector.Rotate`; for an integer `T` the pivot is rounded first, so a half turn of an odd span is not exactly `Reverse`
- `Line.IntersectionCircle(circle)`, `IntersectionRectangle(rectangle)` and `IntersectionPolygon(polygon)`, with `Circle.IntersectionLine`, `Rectangle.IntersectionLine` and `Polygon.IntersectionLine` as mirrors – the points where a segment crosses the boundary, from `Start` to `End`: two through a circle or rectangle, one tangent or ending inside, none apart or entirely inside, and as many as a concave polygon is crossed; a vertex hit by two edges is counted once
- `Size.Abs()` – the absolute size, which every `Rectangle` constructor stores
- `Circle.Intersection(circle)` – the points where two circles cross: two for overlapping, one for tangent within `Epsilon[T]()`, none for apart, nested, concentric, coincident within `Epsilon[T]()` or a negative radius
- `Size.Fit(size)` and `Fill(size)` – uniform scaling to the largest size inside or the smallest around the given one, keeping the aspect ratio; a zero extent has no ratio and gives the zero size
- `Point.RotateAround(pivot, angle)` – rotation about a pivot in the sense of `Vector.Rotate`; `Line.Rotate` and `Polygon.Rotate` are built on it
- `Rectangle.AlignTo(direction, point)` – moves the rectangle so its `Anchor` in that direction lands on the point, the inverse of `Anchor`; `DirectionNone` aligns the center like `MoveTo`
- `Matrix.Translation()`, `Angle()` and `Scaling()` – read a transform back: the `C`/`F` components, the angle of the transformed X axis, and the signed scale factors along the rotated axes, negative in Y for a reflection; a sheared matrix gives the nearest rotation and scale
- `ShearMatrix(shearX, shearY)` and `ReflectionMatrix(axis)` – the two affine constructors that were missing next to translation, rotation and scale; reflecting across `AxisNone` is the identity
- `ParseDirection(name)` and `ParseAxis(name)` – the inverses of `String`
- `Vector.Project(vector)`, `Reject(vector)`, `Reflect(normal)` and `AngleBetween(vector)` – the component along and perpendicular to a vector, the mirror across a surface with a normal of any length, and the unsigned angle in `[0, π]`; the zero vector projects to zero, reflects nothing and is at angle 0 to everything; an integer `Reflect` rounds once, so `Vec(1, 1).Reflect(Vec(1, 2))` is `(0,-1)`, not the `(-1,-1)` of doubling a rounded projection
- `Polygon.Bounds` – the axis-aligned bounding rectangle of the vertices, the zero rectangle for an empty polygon; every shape now has `Bounds()`
- `Matrix.IsInvertible` – reports whether the determinant is non-zero, the check to run before `Inverse`
- `Line.Lerp(t)` – the point at a fraction of the way from `Start` to `End`, extrapolating outside `[0, 1]` like `Point.Lerp`; `Midpoint` is `Lerp(0.5)`
- `Line.DistanceToLine(line)` and `DistanceSquaredToLine(line)` – the distance between the nearest points of two segments and its square for comparisons, zero exactly where `Intersects` holds
- `Rectangle.DistanceSquaredTo(point)`, `Circle.DistanceSquaredTo(point)` and `Polygon.DistanceSquaredTo(point)` – the squared distances, so every shape has the same pair as `Point` and `Line`: in `T` for the rectangle since its nearest point is the clamped lattice point, `float64` for the circle and the polygon; the circle's saves no root and exists for symmetry
- `Line.DistanceSquaredTo(point)` – the squared distance for comparisons, a `float64` even for an integer `T` since the nearest point of a segment is not a lattice point; `Contains` compares it against the squared tolerance, so `Polygon.Contains` walks the boundary without a square root per edge
- `Line.DistanceTo(point)` – the distance to the nearest point of the segment, zero exactly where `Contains` holds: a point within `Epsilon[T]()` is snapped to zero, and a lattice point on a lattice segment is at zero without any tolerance since the perpendicular distance comes from a cross product rather than a projection
- `Line.Contains(point)` – reports whether a point lies on the segment, closed within `Epsilon[T]()` like every other `Contains`
- `Polygon.Edges` – the edges in vertex order, the last closing back to the first vertex; a nil `Vertices` maps to nil like every other mapping
- `Polygon.Area` – the enclosed area by the shoelace formula regardless of winding; a `float64` even for an integer `T`, since a lattice polygon can enclose half a unit
- `Polygon.Perimeter` – the total edge length
- `Polygon.Contains(point)` – point-in-polygon by the even-odd rule, boundary included within `Epsilon[T]()`
- `Rectangle.MinMaxString` – the rectangle by its corners, `(0,1)-(2,4)`, the form `String` printed before

### Changed
- `Circle.Bounds` returns the zero-size rectangle at the center for a negative radius, which contains nothing, instead of a rectangle with a negative size
- `Rectangle` never stores a negative size: `Rect`, `Resize`, `Scale` and `ScaleXY` take the size absolute, `RectangleFromMin` and `RectangleFromMax` measure a negative extent the other way from the given corner, and `RectangleFromMinMax` accepts its two corners in either order, so `Contains`, `Clamp` and the `Intersects` methods always have a `Min` below `Max` unless a struct literal or JSON says otherwise (**breaking** for a caller relying on a mirrored size)
- `Circle.Anchor` returns the center for a negative radius, which contains and intersects nothing and so has no boundary to anchor on; it placed the anchor on the far side of the center, where a negative length put it
- `AssertRect` is renamed to `AssertRectangle`, matching `AssertCircle`, `AssertPolygon` and `AssertRegularPolygon`, none of which abbreviate (**breaking**)
- `Rectangle.Vertices` and its edge walk compute `Min` and `Max` once instead of six times, and `Polygon.DistanceSquaredTo` walks the edges once for both the boundary test and the nearest edge, so `Contains` and `DistanceSquaredTo` are built on the same pass and agree by construction
- `Point.Between` documents that `a` must be the lesser corner on each axis, as `Rectangle.Contains` passes it, and that corners out of order contain nothing
- `Mod` documents that a negative modulus wraps into `(m, 0]`, taking the sign of `m` like a floored division
- `CollisionRectangles`, `CollisionCircles` and `CollisionRectangleCircle` are replaced by the methods `Rectangle.Intersects`, `Circle.Intersects`, `Rectangle.IntersectsCircle` and `Circle.IntersectsRectangle`, so intersection sits next to `Contains` on the shape (**breaking**)
- `Rectangle.Inset` stops an edge at its opposite edge, so a padding larger than the rectangle collapses it to a zero extent inside the original bounds; `Rect(Pt(0, 0), Sz(10, 10)).Inset(Pad(0, 0, 0, 20))` sits at `x = 5` instead of `x = 15`, outside the rectangle it was inset from (**breaking**)
- `Rectangle.Grow`, `GrowXY`, `Shrink` and `ShrinkXY` document where an odd integer amount lands: on the `Min` side when the new extent is even and on the `Max` side when it is odd, so two `Grow(1)` calls move each side once
- `Rectangle.Contains`, `Circle.Contains`, `Vector.LessOrEqual`, `Rectangle.Intersects`, `Circle.Intersects` and `Rectangle.IntersectsCircle` are closed within `Epsilon[T]()`: touching shapes collide, a float rectangle contains the corners it was built from even where `Min` is recomputed from `Center` with a rounding error, and a circle contains its anchors; a negative radius intersects nothing, on `Circle.Intersects` as on every other shape. `Vector.Less` stays strict, and `Rectangle.IntersectsCircle` tests the point of the rectangle closest to the circle center instead of rounding the half extents of an odd integer size (**breaking**)
- Every product, distance and interpolation is computed in `float64` and stored back through `Cast`: `Vector.Dot`, `Cross`, `LengthSquared`, `Less`, `LessOrEqual`, `IsNormalized`, `Point.ManhattanDistanceTo`, `ChebyshevDistanceTo`, `OctileDistanceTo`, `Size.Area`, `Perimeter`, `Circle.Area`, `Polygon.Center`, `Matrix.Multiply`, `Determinant`, `Inverse`, `IsInvertible`, `Lerp` and `EqualDelta`. A narrow `int8` or `int16` no longer overflows mid-computation – `EqualDelta[int8](127, -128, 1)` was true and `Lerp[int8](-100, 100, 0.5)` gave `-128` – while an `int64` beyond 2^53 loses precision instead (**breaking** for such values)
- `Direction` and `Axis` implement `encoding.TextMarshaler` and `TextUnmarshaler`, so JSON stores `"UpRight"` and `"Horizontal"` instead of `7` and `0`, and both work as map keys; an unknown name fails to decode (**breaking** for stored data)
- `Polygon.Intersects` rejects polygons whose `Bounds` do not intersect before examining any edge pair, so two polygons apart cost one box test instead of every edge against every edge
- `Polygon.Center` returns the area centroid by the shoelace formula instead of the vertex average, so a vertex added in the middle of an edge no longer moves it and `Scale` and `MoveTo` pivot on the true center; a polygon enclosing no area, with fewer than three vertices or all collinear, falls back to the vertex average. The sum runs with the origin at the first vertex and rounds every product on its own, so a degenerate edge contributes exactly zero even where the compiler fuses multiply and add (**breaking**)
- `Polygon.Center` rounds half away from zero like every other integer result instead of truncating toward zero, so an integer `MoveTo` lands on the point except when a half changes sign, and returns the zero point for a polygon without vertices instead of dividing by zero (**breaking**)
- `Polygon.Area`, `Perimeter` and `Contains` walk the edges with an iterator instead of allocating an edge slice, and `Edges` builds its slice with an index loop instead of a counter threaded through a closure
- `Vector.Dot`, `Cross` and the `Matrix` determinant round their two products separately, so a fused multiply-add cannot turn the dot product of perpendicular vectors, the cross product of parallel ones or the determinant of a singular matrix into a rounding error: `Vec(13.5, 1.9).Cross(Vec(13.5, 1.9))` is exactly `0` instead of `-2.2e-16`, and `Mat(0.1, 0.3, 0, 0.1, 0.3, 0).IsInvertible()` is false on arm64 as on amd64, where `Inverse` returned components of `1e17` instead of panicking; `Polygon.Area` and `Center` rely on it for degenerate edges
- `Polygon.Contains` rejects a point outside the extent of the vertices before examining any edge, and decides which side of an edge the point is on by the sign of a cross product instead of dividing by the edge's Y span
- `Size.Grow`, `Shrink`, `GrowXY`, `ShrinkXY` and their `Rectangle` counterparts document that the amount is the total change of the extent, not an amount per side, where `Inset` and `Outset` move every side by the full padding
- `Divide` names its parameter `factor` like `Multiply` and every other scaling function; the docs follow
- `Circle.IsZero`, `Rectangle.IsZero` and `Line.IsZero` are `Equal` against the zero value, the form every other type already uses; `Vector.LessOrEqual`, `Polygon.Equal` and `Direction` document that a negative length never counts, that nil and empty vertices are equal, and that the numbering follows the normalized angle
- `Point.Compare` documents that it applies no tolerance, unlike `Equal`, so two float points `Equal` considers the same can still order apart
- `RegularPolygon` documents that `Size` holds the semi-axes of the ellipse the vertices lie on, a radius rather than the full extent `Rectangle.Size` holds, and `Circle.Diameter` documents that it stays in `T` and only overflows for a result outside the range of `T`
- `Cast` panics when a `NaN` or `±Inf` would be stored into an integer `T`, the same convention `Divide` follows for a zero factor, so `Multiply`, `Lerp`, `Vector.Rotate`, `Transform` and `Int()` on an integer shape panic on a non-finite input; a float `T` keeps it. A finite value outside the range of `T` stays unchecked and platform-dependent, as documented (**breaking**)
- `Divide` panics for a zero factor instead of returning the value unchanged, like the integer `/` operator and `Mod`, so `Point.Divide(0)`, `Vector.DivideXY(0, 2)`, `Size.Unscale(0)` and `Matrix.Unscale(0, 3)` panic instead of silently skipping the axis; `Matrix.Inverse` panics for a singular matrix instead of returning it (**breaking**)
- `Matrix.Scale`, `PreScale` and `Unscale` take `float64` factors like every other `Scale` in the package, so an integer matrix can be scaled by `0.5`; each scaled component is rounded, following `Multiply`. `ScaleMatrix` keeps its `T` factors like `TranslationMatrix` (**breaking** for a `Matrix[int]` passing typed `int` factors)
- `Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, so a `float32` gets a tolerance matched to its precision instead of the `float64` one. The comparison stays absolute, and therefore a subtract and a compare for hot paths
- `AssertNumber` and every `Assert*` helper compare an integer `T` exactly and a float `T` within `EpsilonRelative`, instead of within `Delta` for both, so an integer assertion no longer accepts a neighbouring value and a float assertion holds at any magnitude
- The `Assert*` helpers take a `Testing` interface declared by `geom` instead of `assert.Testing`, so no third-party type appears in the public API; `*testing.T` satisfies it unchanged
- `AssertPoint`, `AssertVector`, `AssertSize`, `AssertCircle`, `AssertLine`, `AssertRect`, `AssertPolygon`, `AssertVertices`, `AssertRegularPolygon` and `AssertPadding` take the expected value as the type they assert on – `AssertPoint(t, p, Pt(1, 2))` instead of `AssertPoint(t, p, 1, 2)` – matching `AssertNumber` and `AssertMatrix` (**breaking**)
- `Point.Transform` and `Vector.Transform` accept a `Matrix[M]` of any `Float` type instead of only `Matrix[float64]`, using a method type parameter (Go 1.27); an integer matrix converts with `Matrix.Float` at the call, the same way an angle is `float64`
- `String` (and therefore every `String()` on a geometry type) formats by `T` rather than by value: a float `8.0` prints as `8.00`, like `8.1`, instead of `8`, so `Pt(100, -34.0000115).String()` no longer mixes both forms (**breaking**)
- `Circle.String`, `Line.String` and `Rectangle.String` print in the form of their constructor, so every shape follows one rule: `Circ((10,16);5)` instead of `C(...)`, `Ln((10,16);(1,2))` instead of `L(...)`, and `Rect((1,2);2x3)` with center and size instead of the bare corners, which moved to `MinMaxString` (**breaking**)
- `Line` marshals `Start` and `End` under the JSON keys `s` and `e` instead of `a` and `b`, the initials of the fields like every other key in the package (**breaking**)
- `RegularPolygon.Equal`, `RegularPolygon.IsZero` and `AssertRegularPolygon` compare the angle as well, with `EqualAngle`, so two polygons that produce different vertices no longer compare equal while a full turn, the sign of an angle or the `0`/`2π` seam does not matter (**breaking**)
- `RegularPolygonOrientationAngle` returns `3π/2` instead of `-π/2` for `PointyTop`, the same normalized form `Rotate` stores; returns the top angle for `n < 1` instead of dividing by `n` and storing a `NaN` angle; and panics for an `Orientation` other than `FlatTop` and `PointyTop` instead of returning `0`
- `RegularPolygon.Vertices` returns nil for `N < 1` instead of panicking in `make`, so `RegularPolygon.Polygon()` of an empty polygon is zero like `Pol(nil)`; `RegularPolygon.Empty` is true for any `N < 1`
- `RegularPolygon.Bounds` returns the zero rectangle for `N < 1`, the same answer `Polygon.Bounds` gives for no vertices, instead of a zero-size rectangle at the center (**breaking**)
- `Directions`, `CardinalDirections`, `DiagonalDirections` and `Axes` are functions returning a fresh array instead of package-level variables an importer could write into (**breaking**)
- `Vector.Resize` on the zero vector returns `(length,0)` instead of NaN, the same +X convention `Normalize` uses; `Vector.Less` is false for a non-positive length instead of squaring the sign away; `Resize`, `Normalize` and `Direction` treat only the exact zero vector as directionless, where the tolerant `IsZero` snapped any vector shorter than `Delta` to `(1,0)` or `DirectionNone`
- `DirectionFromAngle` returns `DirectionNone` for `NaN` and `±Inf` instead of `DirectionRight` or a platform-dependent direction, and normalizes the angle before rounding it to a step; `Direction.Angle` returns `NaN` for `DirectionNone` instead of `0`, which was indistinguishable from `DirectionRight`, so the two round-trip for every direction (**breaking**)
- `Axis.ScaleAlong` on `AxisNone` returns the size unchanged instead of a zero size, and with a negative factor flips the sign the way `Size.ScaleXY` does; `Axis.Size` stores the values as given, like `Sz`, instead of taking their absolute value
- `Size.Grow`, `Size.GrowXY`, `Rectangle.Grow`, `Rectangle.GrowXY` and `Circle.Grow` clamp to zero for a negative amount, the same way `Shrink` already did, so no method can produce the negative size `Size` documents as unsupported (**breaking**)
- `ParseSize` wraps the underlying parse error with `%w`, so `errors.Is(err, strconv.ErrRange)` and `ErrSyntax` work
- Documentation: `Rectangle` is closed, an integer rectangle of width `w` spans `w+1` lattice columns while `Rectangle()` yields the half-open pixel rectangle of exactly `w` pixels; `Cast` leaves a finite out-of-range value unchecked and `Int` truncates an `int64` on a 32-bit target; `Mod` panics for `m == 0`; `Direction.Vector` keeps an integer diagonal as the lattice step `(1,1)` where `Unit` snaps to `(1,0)`, and `Circle.Anchor` can fall outside a small integer circle; `RegularPolygon.Vertices` starts from `Angle`; `Polygon` shares the slice `Pol` is given and does not check the vertex count; `Parse` accepts the `NaN` and `Inf` literals for a float `T`; `Size` treats a negative width or height as unsupported; `Axis.IsNone` counts every value outside the two axes; `Matrix` states what an integer `T` cannot represent (only quarter turns rotate, only `|det| = 1` inverts); every `Direction` constant and the `ints` and `floats` aliases have doc comments; the package doc and README state the rounding, arithmetic and boundary conventions and note that the flat JSON shape needs the v2-backed `encoding/json`, the default since Go 1.27
- Internal: `Point.Transform` and `Vector.Transform` convert the matrix with `Matrix.Float` once; `ParseSize` splits with `strings.Cut`; `Triangle`, `Square` and `Hexagon` delegate to `RegularPolygonWithOrientation`; `Axis` methods return from each `switch` case directly

### Fixed
- `Polygon.IntersectsCircle` reports false for a circle with a negative radius whose center lies inside the polygon, as `Rectangle.IntersectsCircle` and `Line.IntersectsCircle` already did
- `Polygon.UnmarshalJSON` decodes into a fresh slice instead of the one the polygon holds, so decoding into `Pol(shared)` no longer writes the new vertices into `shared`
- `String` prints a float that rounds to zero as `0.00` without a sign; `String(-0.004)` gave `-0.00` while a negative zero already printed unsigned
- `LessOrEqual` compares an integer `T` in `T` instead of through `float64`, so `Rectangle[int64].Contains` and `Intersects` stay exact beyond 2^53, where `LessOrEqual[int64](1<<53+1, 1<<53)` was true
- `Circle.Intersects` sums the radii in `float64` instead of `T`, so two `Circle[int8]` of radius 100 no longer wrap the threshold to `-56` and report no collision while overlapping
- `Vector.Resize` and `Normalize` divide each component by the current length before scaling it instead of multiplying by a ratio, so a subnormal vector such as `Vec(5e-324, 0)` resizes to `(length,0)` instead of `(+Inf,NaN)`
- `Point.Int`, `Vector.Int`, `Size.Int`, `Circle.Int`, `Padding.Int`, `Matrix.Int` and the shapes built on them convert an integer `T` directly instead of through `float64`, and `Abs`, `Round`, `Floor` and `Ceil` stay in `T`, so an `int64` beyond 2^53 stays exact
- `Polygon.Translate`, `Scale`, `ScaleXY`, `Int` and `Float` keep a nil `Vertices` nil, so `IsZero` survives every mapping (requires `gravitton/x` v1.2.1, where `slices.Map` maps nil to nil)
- `String` prints a float negative zero as `0.00` instead of `-0.00`, which `Matrix.Inverse` produced for every zero component it negated
- `Circle.Bounds` returns a square of side `Diameter` instead of `Radius`, so the rectangle actually bounds the circle; previously `c.Bounds().IntersectsCircle(c)` could miss (**breaking**)
- `Vector.Normalize` snaps an integer vector to the longer axis and keeps its sign, so the result is always one of the four axis-aligned unit vectors as documented; previously `Vec(-10, -16)` gave `(1,0)`, pointing the opposite way, and `Vec(-10, 16)` kept a diagonal of length √2. A tie resolves to the X axis. `Direction.Unit` follows the same rule, so `DirectionUpRight.Unit[int]()` is `(1,0)` instead of `(1,-1)`; `Direction.Offset` still gives the lattice step (**breaking**)
- `Rectangle.Inset` moves the center down for a larger top padding instead of up, and on an integer rectangle with an odd asymmetric padding no longer leaks outside the original: the inset is derived from the padded `Min` corner, so `Rect(Pt(0,0), Sz(10,10)).Inset(Pad(0,0,0,1))` spans `-4..5` instead of `-3..6`
- `Axis.Project` keeps the sign of the component, so projecting `(-3, 2)` onto the horizontal axis gives `-3` instead of `3`
- `Matrix.Unscale` divides each column by its own factor, following `Divide`, instead of multiplying by a rounded inverse scale; a zero factor panics instead of producing `[[+Inf, 0, 0], [NaN, 1, 0]]`, and an integer matrix unscales exactly when its components divide: `ScaleMatrix(4, 6).Unscale(2, 3)` is `ScaleMatrix(2, 2)` instead of the singular `ScaleMatrix(4, 0)` (**breaking** for integer `T`)
- `Parse` accepts defined types over an integer or float (`Parse[Direction]("3")`, `ParseSize[Coord]`) instead of failing with "unsupported number type": it parses in `int64`/`float64` and narrows to `T`, checking the range, and returns the zero value on error
- `ParseSize` documented that float values are rounded for an integer `T`; it has always rejected them, and the doc now says so
- `isIntType` reported `false` for defined types over an integer (`type Coord int`), so `Cast[Coord]` truncated instead of rounding
- `RegularPolygonOrientationAngle` with `FlatTop` placed a flat edge at the *bottom* (`π/2 - π/n`), which only coincides with a flat top for an even `n`, so a flat-top triangle or pentagon was indistinguishable from `PointyTop`. It now puts an edge midpoint at the top for every `n`, `3π/2 - π/n`; for an even `n` the shape is unchanged but `Vertices()` starts on the opposite side (**breaking**)
- `NormalizeAngle` could return exactly `2π` for a tiny negative angle, violating its `[0, 2π)` contract; it now returns `0` there, and `NaN` for `NaN` and `±Inf` input

## [v1.12.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.11.0...v1.12.0)
### Added
- `Point.Compare(point) int` – orders points by X and then by Y in the `cmp.Compare` convention, for `slices.SortFunc` and `slices.BinarySearchFunc`. The comparison is exact and does not apply the `Delta` tolerance `Point.Equal` uses, since a tolerant comparison is not transitive and would leave a sort no valid ordering


## [v1.11.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.10.0...v1.11.0)
### Changed
- `Direction` constants are renumbered to run by increasing angle — `Right`, `DownRight`, `Down`, `DownLeft`, `Left`, `UpLeft`, `Up`, `UpRight` — so a positive `Rotate` step turns the same way as a positive `Vector.Rotate` angle: counterclockwise in math coordinates, clockwise as drawn on a screen with Y pointing down. Previously the order ran the opposite way and `DirectionFromAngle` negated its input to compensate (**breaking**)
- `Directions`, `CardinalDirections`, and `DiagonalDirections` follow the new order (**breaking**)
- `DirectionFromAngle` no longer negates its argument, now that direction numbering follows the angle directly
- `Rectangle.Vertices` and `Rectangle.Edges` now wind by increasing angle, matching `Directions` and `RegularPolygon.Vertices` — clockwise as drawn on a screen with Y pointing down. `Rectangle.Polygon()` winds the same way as a result (**breaking**)
- `Rectangle.TopEdge`, `RightEdge`, `BottomEdge`, and `LeftEdge` are reversed so the edges still form a closed chain in the new winding (**breaking**)

### Fixed
- `Direction.Rotate` documentation had the screen sense inverted — it claimed positive steps appear clockwise while the old ordering made them appear counterclockwise
- `RegularPolygon.Vertices` documentation claimed counter-clockwise winding; it winds by increasing angle, which is clockwise as drawn
- `Direction.Angle`, `Vector.Rotate`, `VectorFromAngle`, `Rectangle.Vertices`, and `Rectangle.Edges` now state which convention their angles and winding use


## [v1.10.0 (2026-08-26)](https://github.com/gravitton/geometry/compare/v1.9.0...v1.10.0)
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
- `Direction[T](x T) T` renamed to `Sign[T](x T) T`, freeing the name for the new `Direction` type
- `Line.Direction() Vector[T]` renamed to `Line.Vector()`, since it returns the start-to-end vector rather than a `Direction`
- `Line.Reversed()` renamed to `Line.Reverse()`, matching the imperative naming of every other immutable method
- `DirectionFromAxes(up, down, left, right bool)` now returns `Direction` instead of `Vector[float64]` (**breaking**)
- `RectFromMin`, `RectFromMax`, `RectFromMinMax`, `RectFromSize`, `RectFromImage` renamed to `RectangleFrom…`, so every `From` constructor uses its full type name like `PointFromImage` and `VectorFromAngle` (**breaking**)
- `Assert*` helpers take `assert.Testing` instead of `*testing.T`, matching the interface `gravitton/assert` already uses; existing calls passing `*testing.T` are unaffected

### Removed
- `UpVector`, `DownVector`, `LeftVector`, `RightVector`, `UpLeftVector`, `UpRightVector`, `DownLeftVector`, `DownRightVector` – superseded by `DirectionX.Unit[T]()` (**breaking**)

### Fixed
- `Matrix.Int()` now rounds to the nearest integer via `Cast` instead of truncating toward zero, matching every other `Int()` conversion
- `RegularPolygon.Vertices` now rounds after scaling by the size instead of rounding a unit vector first, so integer vertices land within half a unit of the exact position (previously up to 1.5 units off at non-right angles)


## [v1.9.0 (2026-05-29)](https://github.com/gravitton/geometry/compare/v1.8.1...v1.9.0)
### Added
- `Parse[T Number](s string) (T, error)` – parses a string into any `Number` type using `strconv.ParseInt` (with correct bit size) for integer types and `strconv.ParseFloat` for float types
- `ParseSize[T Number](s string) (Size[T], error)` – parses a size string in `"WxH"` form into `Size[T]`


## [v1.8.1 (2026-05-22)](https://github.com/gravitton/geometry/compare/v1.8.0...v1.8.1)
### Fixed
- `RegularPolygonOrientationAngle` with `PointyTop` now returns `-π/2` instead of `π/2`, correctly placing the first vertex at the visual top in the library's +Y-down coordinate system


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

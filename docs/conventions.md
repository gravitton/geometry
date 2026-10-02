# Conventions

The rules every type in `geom` follows. The full reference is the godoc on
[pkg.go.dev](https://pkg.go.dev/github.com/gravitton/geometry).

## Coordinates

The origin is top-left and `+Y` points down. Angles follow the mathematical convention, so a positive angle is
counterclockwise in math coordinates and appears clockwise on screen. Direction order and polygon winding follow the
same rule: `Directions`, `Rectangle.Vertices` and `RegularPolygon.Vertices` all wind by increasing angle, clockwise as
drawn. Each `Direction` has three names: canonical (`DirectionDownRight`), compass (`SouthEast`) and rectangle corner
(`BottomRight`).

## Numbers

Products, distances and interpolations are computed in `float64` and rounded back into `T`, so a narrow `int8` never
overflows mid-computation. A float result stored into an integer `T` rounds half away from zero; `Rectangle` and `Box`
are the exception, truncating half the size toward `Min` so that `Max - Min` stays exactly the size.

The signs that decide a crossing, a turn or a winding are cross products of coordinate differences, computed in
`float64`: exact for an integer `T` while the differences stay within 2^26, beyond which a sign can round to zero.

`Divide`, `Unscale` and `Matrix.Inverse` on a singular matrix panic, like the integer `/` operator; check
`IsInvertible` first when a matrix may be singular. Every other degenerate input returns a value the type can express.

## Equality

`Equal` compares an integer `T` exactly and a float `T` within `Epsilon[T]()`, a tolerance matched to `float32` or
`float64`. `EqualRelative` scales it for values far from zero, and `EqualAngle` compares modulo a full turn.

## Boundaries

`Contains`, `Intersects` and `Vector.LessOrEqual` are closed and tolerant: a point within the tolerance of the
boundary counts as on it, so a float rectangle contains the corners it was built from and a polygon contains its
vertices. The tolerance is `Epsilon[T]()` near the origin and widens to two ulps of `T` at the largest coordinate
compared, so a `Nearest` point or an `Intersection` rounded into `T` far from the origin stays on the boundary;
`float32` widens beyond a few hundred units. `Vector.Less` is strict.

Every boundary is judged on a distance, never on a coordinate, so two rectangles that meet corner to corner intersect
exactly where their polygons do. `DistanceTo` is zero exactly where `Contains` holds, and `Intersection` answers
exactly where `Intersects` holds, apart from parallel and coincident segments and rectangles of different angles,
which have no single answer.

## Intersections

The logic of each pair is written once, on the earlier shape of `Circle`, `Segment`, `Ray`, `Polygon`, `Rectangle`,
`RegularPolygon`, `Box`, and the other side delegates to it, so `a.IntersectsCircle(c)` and `c.IntersectsRectangle(a)`
always give the same answer.

## Matrices

An integer `Matrix` composes lattice transforms exactly: translation, integer scale, reflection, quarter turns.
Anything else rounds into a different matrix; use a float `Matrix` there. `Transform` takes a float matrix, so convert
an integer one with `Float()` at the call.

## Reproducibility

Every product is rounded before it is added, so the same inputs give the same bits on amd64 and arm64, where the
compiler would otherwise fuse the two into one multiply-add. Only what `math` computes from an angle (`Sincos`,
`Atan2`, `Hypot`) may differ in the last bit between architectures.

## Methods

Every value type has `Equal`, `Cast`, `Int`, `Float` and `String`; every shape adds `Bounds` and `Contains`, and every
shape but `Ellipse` adds `Intersects`. The enums (`Direction`, `Axis`, `Orientation`, `Winding`) are compared with `==`
and marshal as their names.

## JSON

The flat JSON of `Rectangle`, `Circle` and `RegularPolygon` relies on the `embed` struct tag of the v2-backed
`encoding/json`, the default since Go 1.27; under `GOEXPERIMENT=nojsonv2` the nested fields are emitted as objects.

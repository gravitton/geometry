# Geometry

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Generic immutable 2D geometry library for game development.

> This library use top-left origin with `+Y` down and `-Y` up. 
> 
> It only affect human-readable getters (like `Up`, `Down`, `Top`, `Bottom`).


## Installation

```bash
go get github.com/gravitton/geometry
```


## Usage

```go
package main

import (
	geom "github.com/gravitton/geometry"
)

type HexLayout struct {
	Origin    geom.Point[float64]
	Size      geom.Size[float64]
	FromPixel geom.Matrix
}

func (l HexLayout) FromPixel(pixel geom.Point[float64]) (Q, R float64) {
	return pixel.Subtract(l.Origin).DivideXY(l.Size.XY()).Transform(l.FromPixel).XY()
}
```

It also provide packages for type aliases ([`ints`](./types/ints/types.go) for `int` and [`floats`](./types/floats/types.go) for `float64`).

```go
package main

import (
    "github.com/gravitton/geometry/types/floats"
    "github.com/gravitton/geometry/types/ints"
)

type Grid struct {
    Size     ints.Size
    CellSize floats.Size
}

func Bounds(grid *Grid, x, y int) floats.Size {
    return geom.RectFromMin(floats.Pt(x, y), grid.CellSize.ScaleXY(grid.Size.Float().XY()))
}
```

## API

All types and methods are generic and can be used with any numeric type from `Number` type constraint.

```go
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}
```

### Constructors

All types have exported fields and can be created directly (`geom.Point[int]{X: 1, Y: 2}`), but there are also shorthand constructors (`geom.Pt(1, 2)`).

```go
func Pt[T Number](x, y T) Point[T]
func Vec[T Number](x, y T) Vector[T]
func Sz[T Number](w, h T) Size[T]
func Circ[T Number](center Point[T], radius T) Circle[T]
func Rect[T Number](center Point[T], size Size[T]) Rectangle[T]
func Ln[T Number](start, end Point[T]) Line[T]
func Pol[T Number](vertices []Point[T]) Polygon[T]
func RegPol[T Number](center Point[T], size Size[T], n int, angle float64) RegularPolygon[T]
func Mat(a, b, c, d, e, f float64) Matrix
func Pad[T Number](top, right, bottom, left T) Padding[T]
```

### Point

```go
type Point[T Number] struct{ 
	X, Y T 
}

// Properties
func (p Point[T]) XY() (T, T)

// Transformations
func (p Point[T]) Transform(matrix Matrix) Point[T]

// Mathematical operations
func (p Point[T]) Add(vector Vector[T]) Point[T]
func (p Point[T]) AddXY(deltaX, deltaY T) Point[T]
func (p Point[T]) Subtract(point Point[T]) Vector[T] 
func (p Point[T]) Multiply(factor float64) Point[T]
func (p Point[T]) MultiplyXY(factorX, factorY float64) Point[T]
func (p Point[T]) Divide(factor float64) Point[T]
func (p Point[T]) DivideXY(factorX, factorY float64) Point[T]

// Geometric operations
func (p Point[T]) DistanceTo(point Point[T]) float64
func (p Point[T]) DistanceSquaredTo(point Point[T]) T
func (p Point[T]) Midpoint(point Point[T]) Point[T]
func (p Point[T]) Lerp(point Point[T], t float64) Point[T]
func (p Point[T]) AngleTo(point Point[T]) float64

// Utilities
func (p Point[T]) Equal(point Point[T]) bool
func (p Point[T]) IsZero() bool
func (p Point[T]) Vector() Vector[T]
func (p Point[T]) Point() image.Point
func (p Point[T]) Int() Point[int]
func (p Point[T]) Float() Point[float64]
func (p Point[T]) String() string
```

### Vector

```go
type Vector[T Number] struct { 
	X, Y T 
}

// Properties
func (v Vector[T]) XY() (T, T)

// Transformations
func (v Vector[T]) Transform(matrix Matrix) Vector[T]

// Mathematical operations
func (v Vector[T]) Add(vector Vector[T]) Vector[T]
func (v Vector[T]) AddXY(deltaX, deltaY T) Vector[T]
func (v Vector[T]) Subtract(vector Vector[T]) Vector[T]
func (v Vector[T]) SubtractXY(deltaX, deltaY T) Vector[T]
func (v Vector[T]) Multiply(factor float64) Vector[T]
func (v Vector[T]) MultiplyXY(factorX, factorY float64) Vector[T]
func (v Vector[T]) Divide(factor float64) Vector[T]
func (v Vector[T]) DivideXY(factorX, factorY float64) Vector[T]
func (v Vector[T]) Negate() Vector[T]

// Vector operations
func (v Vector[T]) Dot(vector Vector[T]) T
func (v Vector[T]) Cross(vector Vector[T]) T
func (v Vector[T]) Length() float64
func (v Vector[T]) LengthSquared() T
func (v Vector[T]) Angle() float64
func (v Vector[T]) Lerp(vector Vector[T], t float64) Vector[T]

// Transformations
func (v Vector[T]) Rotate(angle float64) Vector[T]
func (v Vector[T]) Normal() Vector[T]
func (v Vector[T]) Resize(length float64) Vector[T]
func (v Vector[T]) Normalize() Vector[T]
func (v Vector[T]) Abs() Vector[T]

// Utilities
func (v Vector[T]) Equal(vector Vector[T]) bool
func (v Vector[T]) IsZero() bool
func (v Vector[T]) IsOne() bool
func (v Vector[T]) IsNormalized() bool
func (v Vector[T]) Less(value T) bool
func (v Vector[T]) Point() Point[T] 
func (v Vector[T]) Size() Size[T]
func (v Vector[T]) Int() Vector[int]
func (v Vector[T]) Float() Vector[float64]
func (v Vector[T]) String() string
```

### Matrix

```go
type Matrix struct {
	A, B, C float64 // scale X, shear Y, translate X
	D, E, F float64 // shear X, scale Y, translate Y
	// [0 0 1] implicit third row
}

// Matrix operations
func (m Matrix) Multiply(n Matrix) Matrix
func (m Matrix) Inverse() Matrix
func (m Matrix) Determinant() float64

// Transformations
func (m Matrix) Translate(deltaX, deltaY float64) Matrix
func (m Matrix) Untranslate(deltaX, deltaY float64) Matrix
func (m Matrix) PreTranslate(deltaX, deltaY float64)
func (m Matrix) Rotate(angle float64) Matrix
func (m Matrix) PreRotate(angle float64) Matrix
func (m Matrix) Scale(factorX, factorY float64) Matrix
func (m Matrix) Unscale(factorX, factorY float64) Matrix

// Utilities
func (m Matrix) Equal(matrix Matrix) bool
func (m Matrix) IsZero() bool
func (m Matrix) String() string
```

### Size

```go
type Size[T Number] struct { 
	Width, Height T
}

// Properties
func (s Size[T]) XY() (T, T)
func (s Size[T]) Area() T
func (s Size[T]) Perimeter() T
func (s Size[T]) AspectRatio() float64

// Dimension operations
func (s Size[T]) Scale(factor float64) Size[T]
func (s Size[T]) ScaleXY(factorX, factorY float64) Size[T]
func (s Size[T]) Grow(amount T) Size[T]
func (s Size[T]) GrowXY(amountX, amountY T) Size[T]
func (s Size[T]) Shrink(amount T) Size[T]
func (s Size[T]) ShrinkXY(amountX, amountY T) Size[T]

// Utilities
func (s Size[T]) Equal(size Size[T]) bool
func (s Size[T]) IsZero() bool
func (s Size[T]) Vector() Vector[T]
func (s Size[T]) Int() Size[int]
func (s Size[T]) Float() Size[float64]
func (s Size[T]) String() string
```

### Circle

```go
type Circle[T Number] struct { 
	Center Point[T]
	Radius T
}

// Properties
func (c Circle[T]) Area() float64
func (c Circle[T]) Circumference() float64
func (c Circle[T]) Diameter() T

// Transformations
func (c Circle[T]) Translate(vector Vector[T]) Circle[T]
func (c Circle[T]) MoveTo(center Point[T]) Circle[T]
func (c Circle[T]) Scale(factor float64) Circle[T]

// Size operations
func (c Circle[T]) Resize(radius T) Circle[T]
func (c Circle[T]) Grow(amount T) Circle[T]
func (c Circle[T]) Shrink(amount T) Circle[T]

// Geometric queries
func (c Circle[T]) Contains(point Point[T]) bool

// Utilities
func (c Circle[T]) Equal(circle Circle[T]) bool
func (c Circle[T]) IsZero() bool
func (c Circle[T]) Bounds() Rectangle[T]
func (c Circle[T]) Int() Circle[int]
func (c Circle[T]) Float() Circle[float64]
func (c Circle[T]) String() string
```

### Rectangle

```go
type Rectangle[T Number] struct{
	Center Point[T]
	Size   Size[T]
}

// Properties
func (r Rectangle[T]) Width() T
func (r Rectangle[T]) Height() T
func (r Rectangle[T]) Min() Point[T]
func (r Rectangle[T]) Max() Point[T]
func (r Rectangle[T]) BottomLeft() Point[T]
func (r Rectangle[T]) BottomRight() Point[T]
func (r Rectangle[T]) TopLeft() Point[T]
func (r Rectangle[T]) TopRight() Point[T]
func (r Rectangle[T]) Edges() []Line[T]
func (r Rectangle[T]) Vertices() []Point[T]
func (r Rectangle[T]) Area() T
func (r Rectangle[T]) Perimeter() T
func (r Rectangle[T]) AspectRatio() float64

// Transformations
func (r Rectangle[T]) Translate(vector Vector[T]) Rectangle[T]
func (r Rectangle[T]) MoveTo(center Point[T]) Rectangle[T]
func (r Rectangle[T]) Scale(factor float64) Rectangle[T]
func (r Rectangle[T]) ScaleXY(factorX, factorY float64) Rectangle[T]

// Size operations
func (r Rectangle[T]) Resize(size Size[T]) Rectangle[T]
func (r Rectangle[T]) Grow(amount T) Rectangle[T]
func (r Rectangle[T]) GrowXY(amountX, amountY T) Rectangle[T]
func (r Rectangle[T]) Shrink(amount T) Rectangle[T]
func (r Rectangle[T]) ShrinkXY(amountX, amountY T) Rectangle[T]
func (r Rectangle[T]) Inset(padding Padding[T]) Rectangle[T]

// Geometric queries
func (r Rectangle[T]) Contains(point Point[T]) bool

// Utilities
func (r Rectangle[T]) Equal(rectangle Rectangle[T]) bool
func (r Rectangle[T]) IsZero() bool
func (r Rectangle[T]) Bounds() Rectangle[T]
func (r Rectangle[T]) Polygon() Polygon[T]
func (r Rectangle[T]) Rectangle() image.Rectangle
func (r Rectangle[T]) Int() Rectangle[int]
func (r Rectangle[T]) Float() Rectangle[float64]
func (r Rectangle[T]) String() string
```

### Line

```go
type Line[T Number] struct{
	Start, End Point[T]
}

// Transformations
func (l Line[T]) Translate(vector Vector[T]) Line[T]
func (l Line[T]) MoveTo(point Point[T]) Line[T]
func (l Line[T]) Reversed() Line[T]

// Measurements
func (l Line[T]) Midpoint() Point[T]
func (l Line[T]) Direction() Vector[T]
func (l Line[T]) Length() float64

// Utilities
func (l Line[T]) Equal(line Line[T]) bool
func (l Line[T]) IsZero() bool
func (l Line[T]) Bounds() Rectangle[T]
func (l Line[T]) Int() Line[int]
func (l Line[T]) Float() Line[float64]
func (l Line[T]) String() string
```

### Polygon

```go
type Polygon[T Number] struct{
	Vertices []Point[T]
}

// Properties
func (p Polygon[T]) Center() Point[T]

// Transformations
func (p Polygon[T]) Translate(vector Vector[T]) Polygon[T]
func (p Polygon[T]) MoveTo(center Point[T]) Polygon[T]
func (p Polygon[T]) Scale(factor float64) Polygon[T]
func (p Polygon[T]) ScaleXY(factorX, factorY float64) Polygon[T]

// Utilities
func (p Polygon[T]) Equal(polygon Polygon[T]) bool
func (p Polygon[T]) IsZero() bool
func (p Polygon[T]) Empty() bool
func (p Polygon[T]) Int() Polygon[int]
func (p Polygon[T]) Float() Polygon[float64]
func (p Polygon[T]) String() string
```

### Regular Polygon

```go
type RegularPolygon[T Number] struct {
	Center Point[T]
	Size   Size[T]
	N      int
	Angle  float64
}

// Properties
func (rp RegularPolygon[T]) Vertices() []Point[T]

// Transformations
func (rp RegularPolygon[T]) Translate(vector Vector[T]) RegularPolygon[T]
func (rp RegularPolygon[T]) MoveTo(center Point[T]) RegularPolygon[T]
func (rp RegularPolygon[T]) Scale(factor float64) RegularPolygon[T]
func (rp RegularPolygon[T]) ScaleXY(factorX, factorY float64) RegularPolygon[T]
func (rp RegularPolygon[T]) Rotate(angle float64) RegularPolygon[T]

// Utilities
func (rp RegularPolygon[T]) Equal(polygon RegularPolygon[T]) bool
func (rp RegularPolygon[T]) IsZero() bool
func (rp RegularPolygon[T]) Empty() bool
func (rp RegularPolygon[T]) Bounds() Rectangle[T]
func (rp RegularPolygon[T]) Polygon() Polygon[T]
func (rp RegularPolygon[T]) Int() RegularPolygon[int]
func (rp RegularPolygon[T]) Float() RegularPolygon[float64]
func (rp RegularPolygon[T]) String() string
```

### Padding

```go
type Padding[T Number] struct{
	Top, Right, Bottom, Left T
}

// Properties
func (p Padding[T]) XY() (T, T)
func (p Padding[T]) Width() T
func (p Padding[T]) Height() T

// Utilities
func (p Padding[T]) Size() Size[T]
func (p Padding[T]) Int() Padding[int]
func (p Padding[T]) Float() Padding[float64]
func (p Padding[T]) String() string
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

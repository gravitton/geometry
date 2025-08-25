# Geometry

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Generic immutable 2D geometry library for game development.


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
```

## API

All types and methods are generic and can be used with any numeric type.

```go
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}
```

### Point

```go
type Point[T Number] struct{ 
	X, Y T 
}

// Properties
func (p Point[T]) XY() (T, T)

// Mathematical operations
func (p Point[T]) Add(vector Vector[T]) Point[T]
func (p Point[T]) AddXY(deltaX, deltaY T) Point[T]
func (p Point[T]) Subtract(other Point[T]) Vector[T] 
func (p Point[T]) Multiply(scale float64) Point[T]
func (p Point[T]) MultiplyXY(scaleX, scaleY float64) Point[T]
func (p Point[T]) Divide(scale float64) Point[T]
func (p Point[T]) DivideXY(scaleX, scaleY float64) Point[T]

// Geometric operations
func (p Point[T]) DistanceTo(point Point[T]) float64
func (p Point[T]) DistanceSquaredTo(point Point[T]) T
func (p Point[T]) Midpoint(point Point[T]) Point[T]
func (p Point[T]) Lerp(point Point[T], t float64) Point[T]
func (p Point[T]) AngleTo(point Point[T]) float64

// Utilities
func (p Point[T]) Equal(point Point[T]) bool
func (p Point[T]) Zero() bool
```

### Vector

```go
type Vector[T Number] struct { 
	X, Y T 
}

// Properties
func (v Vector[T]) XY() (T, T)

// Mathematical operations
func (v Vector[T]) Add(vector Vector[T]) Vector[T]
func (v Vector[T]) AddXY(deltaX, deltaY T) Vector[T]
func (v Vector[T]) Sub(vector Vector[T]) Vector[T]
func (v Vector[T]) SubXY(deltaX, deltaY T) Vector[T]
func (v Vector[T]) Multiply(scale float64) Vector[T]
func (v Vector[T]) MultiplyXY(scaleX, scaleY float64) Vector[T]
func (v Vector[T]) Divide(scale float64) Vector[T]
func (v Vector[T]) DivideXY(scaleX, scaleY float64) Vector[T]
func (v Vector[T]) Negate() Vector[T]

// Vector operations
func (v Vector[T]) Dot(vector Vector[T]) float64
func (v Vector[T]) Cross(vector Vector[T]) float64
func (v Vector[T]) Length() float64
func (v Vector[T]) LengthSquared() T
func (v Vector[T]) Normalize() Vector[T]
func (v Vector[T]) Angle() float64

// Transformations
func (v Vector[T]) Rotate(angle float64) Vector[T]
func (v Vector[T]) Normal() Vector[T]

// Utilities
func (v Vector[T]) Equal(vector Vector) bool
func (v Vector[T]) Zero() bool
```

### Size

```go
type Size[T Number] struct { 
	Width, Height T
}

// Properties
func (s Size[T]) XY() (T, T)
func (s Size[T]) Area() float64
func (s Size[T]) Perimeter() T
func (s Size[T]) AspectRatio() float64


// Dimension operations
func (s Size[T]) Scale(factor float64) Size[T]
func (s Size[T]) ScaleXY(factorX, factorY float64) Size[T]
func (s Size[T]) Grow(amount T) Size[T]
func (s Size[T]) GrowXY(amountX, amountY T) Size[T]
func (s Size[T]) Shrink(amountX, amountY T) Size[T]

// Utilities
func (s Size[T]) Equal(other Size[T]) bool
func (s Size[T]) Zero() bool
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
func (c Circle[T]) Bounds() Rectangle[T]

// Transformations
func (c Circle[T]) Translate(vector Vector[T]) Circle[T]
func (c Circle[T]) MoveTo(center Point[T]) Circle[T]
func (c Circle[T]) Scale(factor float64) Circle[T]

// Size operations
func (c Circle[T]) Resize(radius T) Circle[T]
func (c Circle[T]) Expand(amount T) Circle[T]
func (c Circle[T]) Shrunk(amount T) Circle[T]

// Geometric queries
func (c Circle[T]) Contains(point Point[T]) bool

// Utilities
func (c Circle[T]) Equal(circle Circle) bool
```

### Rectangle

```go
type Rectangle[T Number] struct{
    Center Point[T]
    Size   Size[T]
}
```

### Line

```go
type Line[T Number] struct{
    Start Point[T]
    End   Point[T]
}
```

### Polygon

```go
type Polygon[T Number] struct{
    Vertices []Point[T]
}
```

### Regular Polygon

```go
type RegularPolygon[T Number] struct {
    Center Point[T]
    Size   Size[T]
    N      int
}
```

### Constructor Functions

Short constructors for convenience

```go
func P[T Number](x, y T) Point[T]
func V[T Number](x, y T) Vector[T]
func S[T Number](w, h T) Size[T]
func C[T Number](center Point[T], radius T) Circle[T]
func R[T Number](center Point[T], size Size[T]) Rectangle[T]
func L[T Number](start, end Point[T]) Line[T]
func RP[T Number](center Point[T], size Size[T], n int) RegularPolygon[T]
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
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/assert?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/geometry/releases
[link-contributors]:        https://github.com/gravitton/geometry/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/geometry/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/geometry
[link-go-report-card]:      https://goreportcard.com/report/github.com/gravitton/geometry
[link-coverage]:            https://coveralls.io/github/gravitton/geometry

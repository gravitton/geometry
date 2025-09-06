package ints

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/internal"
)

type Point = geom.Point[int]
type Vector = geom.Vector[int]
type Size = geom.Size[int]
type Circle = geom.Circle[int]
type Line = geom.Line[int]
type Rectangle = geom.Rectangle[int]
type Polygon = geom.Polygon[int]
type RegularPolygon = geom.RegularPolygon[int]

func P[T geom.Number](x, y T) Point {
	return Point{X: geom.Cast[int](float64(x)), Y: geom.Cast[int](float64(y))}
}

func ToPoint[T geom.Number](point geom.Point[T]) Point {
	return Point{X: geom.Cast[int](float64(point.X)), Y: geom.Cast[int](float64(point.Y))}
}

func V[T geom.Number](x, y T) Vector {
	return Vector{X: geom.Cast[int](float64(x)), Y: geom.Cast[int](float64(y))}
}

func ToVector[T geom.Number](vector geom.Vector[T]) Vector {
	return Vector{X: geom.Cast[int](float64(vector.X)), Y: geom.Cast[int](float64(vector.Y))}
}

func S[T geom.Number](width, height T) Size {
	return Size{Width: geom.Cast[int](float64(width)), Height: geom.Cast[int](float64(height))}
}

func ToSize[T geom.Number](size geom.Size[T]) Size {
	return Size{Width: geom.Cast[int](float64(size.Width)), Height: geom.Cast[int](float64(size.Height))}
}

func C[T geom.Number](center geom.Point[T], radius T) Circle {
	return Circle{Center: ToPoint(center), Radius: geom.Cast[int](float64(radius))}
}

func ToCircle[T geom.Number](circle geom.Circle[T]) Circle {
	return Circle{Center: ToPoint(circle.Center), Radius: geom.Cast[int](float64(circle.Radius))}
}

func L[T geom.Number](start, end geom.Point[T]) Line {
	return Line{Start: ToPoint(start), End: ToPoint(end)}
}

func ToLine[T geom.Number](line geom.Line[T]) Line {
	return Line{Start: ToPoint(line.Start), End: ToPoint(line.End)}
}

func R[T geom.Number](center geom.Point[T], size geom.Size[T]) Rectangle {
	return Rectangle{Center: ToPoint(center), Size: ToSize(size)}
}

func ToRectangle[T geom.Number](rectangle geom.Rectangle[T]) Rectangle {
	return Rectangle{Center: ToPoint(rectangle.Center), Size: ToSize(rectangle.Size)}
}

func ToPolygon[T geom.Number](polygon geom.Polygon[T]) Polygon {
	return Polygon{Vertices: internal.Map(polygon.Vertices, func(point geom.Point[T]) Point {
		return ToPoint(point)
	})}
}

func ToRegularPolygon[T geom.Number](polygon geom.RegularPolygon[T]) RegularPolygon {
	return RegularPolygon{Center: ToPoint(polygon.Center), Size: ToSize(polygon.Size), N: polygon.N}
}

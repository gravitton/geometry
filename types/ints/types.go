// Package ints provides type aliases and constructors for geometry objects with [int] type argument.
package ints

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/x/slices"
)

type Point = geom.Point[int]
type Vector = geom.Vector[int]
type Size = geom.Size[int]
type Circle = geom.Circle[int]
type Line = geom.Line[int]
type Rectangle = geom.Rectangle[int]
type Polygon = geom.Polygon[int]
type RegularPolygon = geom.RegularPolygon[int]
type Padding = geom.Padding[int]

func Pt[T geom.Number](x, y T) Point {
	return Point{X: geom.Cast[int](float64(x)), Y: geom.Cast[int](float64(y))}
}

func ToPoint[T geom.Number](point geom.Point[T]) Point {
	return Pt(point.X, point.Y)
}

func Vec[T geom.Number](x, y T) Vector {
	return Vector{X: geom.Cast[int](float64(x)), Y: geom.Cast[int](float64(y))}
}

func ToVector[T geom.Number](vector geom.Vector[T]) Vector {
	return Vec(vector.X, vector.Y)
}

func Sz[T geom.Number](width, height T) Size {
	return Size{Width: geom.Cast[int](float64(width)), Height: geom.Cast[int](float64(height))}
}

func ToSize[T geom.Number](size geom.Size[T]) Size {
	return Sz(size.Width, size.Height)
}

func Circ[T geom.Number](center geom.Point[T], radius T) Circle {
	return Circle{Center: ToPoint(center), Radius: geom.Cast[int](float64(radius))}
}

func ToCircle[T geom.Number](circle geom.Circle[T]) Circle {
	return Circ(circle.Center, circle.Radius)
}

func Ln[T geom.Number](start, end geom.Point[T]) Line {
	return Line{Start: ToPoint(start), End: ToPoint(end)}
}

func ToLine[T geom.Number](line geom.Line[T]) Line {
	return Ln(line.Start, line.End)
}

func Rect[T geom.Number](center geom.Point[T], size geom.Size[T]) Rectangle {
	return Rectangle{Center: ToPoint(center), Size: ToSize(size)}
}

func ToRectangle[T geom.Number](rectangle geom.Rectangle[T]) Rectangle {
	return Rect(rectangle.Center, rectangle.Size)
}

func Pol[T geom.Number](Vertices []geom.Point[T]) Polygon {
	return Polygon{Vertices: slices.Map(Vertices, func(point geom.Point[T]) Point {
		return ToPoint(point)
	})}
}

func ToPolygon[T geom.Number](polygon geom.Polygon[T]) Polygon {
	return Pol(polygon.Vertices)
}

func RegPol[T geom.Number](center geom.Point[T], size geom.Size[T], n int, angle float64) RegularPolygon {
	return RegularPolygon{Center: ToPoint(center), Size: ToSize(size), N: n, Angle: angle}
}

func ToRegularPolygon[T geom.Number](polygon geom.RegularPolygon[T]) RegularPolygon {
	return RegPol(polygon.Center, polygon.Size, polygon.N, polygon.Angle)
}

func Pd[T geom.Number](top, right, bottom, left T) Padding {
	return Padding{Top: geom.Cast[int](float64(top)), Right: geom.Cast[int](float64(right)), Bottom: geom.Cast[int](float64(bottom)), Left: geom.Cast[int](float64(left))}
}

func ToPadding[T geom.Number](padding geom.Padding[T]) Padding {
	return Pd(padding.Top, padding.Right, padding.Bottom, padding.Left)
}

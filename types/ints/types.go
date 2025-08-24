package ints

import geom "github.com/gravitton/geometry"

type Point = geom.Point[int]
type Vector = geom.Vector[int]
type Size = geom.Size[int]
type Circle = geom.Circle[int]
type Line = geom.Line[int]
type Rectangle = geom.Rectangle[int]
type Polygon = geom.Polygon[int]
type RegularPolygon = geom.RegularPolygon[int]

func ToPoint[T geom.Number](point geom.Point[T]) Point {
	return Point{X: geom.Cast[int](float64(point.X)), Y: geom.Cast[int](float64(point.Y))}
}

func ToVector[T geom.Number](vector geom.Vector[T]) Vector {
	return Vector{X: geom.Cast[int](float64(vector.X)), Y: geom.Cast[int](float64(vector.Y))}
}

func ToSize[T geom.Number](size geom.Size[T]) Size {
	return Size{Width: geom.Cast[int](float64(size.Width)), Height: geom.Cast[int](float64(size.Height))}
}

func ToCircle[T geom.Number](circle geom.Circle[T]) Circle {
	return Circle{Center: ToPoint(circle.Center), Radius: geom.Cast[int](float64(circle.Radius))}
}

func ToLine[T geom.Number](line geom.Line[T]) Line {
	return Line{Start: ToPoint(line.Start), End: ToPoint(line.End)}
}

func ToRectangle[T geom.Number](rectangle geom.Rectangle[T]) Rectangle {
	return Rectangle{Center: ToPoint(rectangle.Center), Size: ToSize(rectangle.Size)}
}

func ToPolygon[T geom.Number](polygon geom.Polygon[T]) Polygon {
	return Polygon{Vertices: geom.Transform(polygon.Vertices, func(point geom.Point[T]) Point {
		return ToPoint(point)
	})}
}

func ToRegularPolygon[T geom.Number](polygon geom.RegularPolygon[T]) RegularPolygon {
	return RegularPolygon{Center: ToPoint(polygon.Center), Size: ToSize(polygon.Size), N: polygon.N}
}

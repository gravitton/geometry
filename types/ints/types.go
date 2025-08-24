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
	return Point{geom.Cast[int](float64(point.X)), geom.Cast[int](float64(point.Y))}
}

func ToVector[T geom.Number](vector geom.Vector[T]) Vector {
	return Vector{geom.Cast[int](float64(vector.X)), geom.Cast[int](float64(vector.Y))}
}

func ToSize[T geom.Number](size geom.Size[T]) Size {
	return Size{geom.Cast[int](float64(size.Width)), geom.Cast[int](float64(size.Height))}
}

func ToCircle[T geom.Number](circle geom.Circle[T]) Circle {
	return Circle{ToPoint(circle.Center), geom.Cast[int](float64(circle.Radius))}
}

func ToLine[T geom.Number](line geom.Line[T]) Line {
	return Line{ToPoint(line.Start), ToPoint(line.End)}
}

func ToRectangle[T geom.Number](rectangle geom.Rectangle[T]) Rectangle {
	return Rectangle{ToPoint(rectangle.Center), ToSize(rectangle.Size)}
}

func ToPolygon[T geom.Number](polygon geom.Polygon[T]) Polygon {
	return Polygon{geom.Transform(polygon.Vertices, func(point geom.Point[T]) Point {
		return ToPoint(point)
	})}
}

func ToRegularPolygon[T geom.Number](polygon geom.RegularPolygon[T]) RegularPolygon {
	return RegularPolygon{ToPoint(polygon.Center), ToSize(polygon.Size), polygon.N}
}

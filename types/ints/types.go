// Package ints provides type aliases and constructors for geometry objects with [int] type argument.
package ints

import (
	geom "github.com/gravitton/geometry"
)

// Point is geom.Point[int].
type Point = geom.Point[int]

// Vector is geom.Vector[int].
type Vector = geom.Vector[int]

// Size is geom.Size[int].
type Size = geom.Size[int]

// Circle is geom.Circle[int].
type Circle = geom.Circle[int]

// Line is geom.Line[int].
type Line = geom.Line[int]

// Rectangle is geom.Rectangle[int].
type Rectangle = geom.Rectangle[int]

// Polygon is geom.Polygon[int].
type Polygon = geom.Polygon[int]

// RegularPolygon is geom.RegularPolygon[int].
type RegularPolygon = geom.RegularPolygon[int]

// Padding is geom.Padding[int].
type Padding = geom.Padding[int]

// Matrix is geom.Matrix[int].
type Matrix = geom.Matrix[int]

// Pt is shorthand for geom.Pt(x, y).Int()
func Pt[T geom.Number](x, y T) Point {
	return geom.Pt(x, y).Int()
}

// Vec is shorthand for geom.Vec(x, y).Int()
func Vec[T geom.Number](x, y T) Vector {
	return geom.Vec(x, y).Int()
}

// Sz is shorthand for geom.Sz(width, height).Int()
func Sz[T geom.Number](width, height T) Size {
	return geom.Sz(width, height).Int()
}

// Circ is shorthand for geom.Circ(center, radius).Int()
func Circ[T geom.Number](center geom.Point[T], radius T) Circle {
	return geom.Circ(center, radius).Int()
}

// Ln is shorthand for geom.Ln(start, end).Int()
func Ln[T geom.Number](start, end geom.Point[T]) Line {
	return geom.Ln(start, end).Int()
}

// Rect is shorthand for geom.Rect(center, size).Int()
func Rect[T geom.Number](center geom.Point[T], size geom.Size[T]) Rectangle {
	return geom.Rect(center, size).Int()
}

// Pol is shorthand for geom.Pol(vertices).Int()
func Pol[T geom.Number](vertices []geom.Point[T]) Polygon {
	return geom.Pol(vertices).Int()
}

// RegPol is shorthand for geom.RegPol(center, size, n, angle).Int()
func RegPol[T geom.Number](center geom.Point[T], size geom.Size[T], n int, angle float64) RegularPolygon {
	return geom.RegPol(center, size, n, angle).Int()
}

// Pad is shorthand for geom.Pad(top, right, bottom, left).Int()
func Pad[T geom.Number](top, right, bottom, left T) Padding {
	return geom.Pad(top, right, bottom, left).Int()
}

// Mat is shorthand for geom.Mat(a, b, c, d, e, f).Int()
func Mat[T geom.Number](a, b, c, d, e, f T) Matrix {
	return geom.Mat(a, b, c, d, e, f).Int()
}

// IdentityMatrix is shorthand for geom.IdentityMatrix[int]()
func IdentityMatrix() Matrix {
	return geom.IdentityMatrix[int]()
}

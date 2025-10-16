// Package ints provides type aliases and constructors for geometry objects with [float64] type argument.
package floats

import (
	geom "github.com/gravitton/geometry"
)

type Point = geom.Point[float64]
type Vector = geom.Vector[float64]
type Size = geom.Size[float64]
type Circle = geom.Circle[float64]
type Line = geom.Line[float64]
type Rectangle = geom.Rectangle[float64]
type Polygon = geom.Polygon[float64]
type RegularPolygon = geom.RegularPolygon[float64]
type Padding = geom.Padding[float64]

// Pt is shorthand for geom.Pt(x, y).Float()
func Pt[T geom.Number](x, y T) Point {
	return geom.Pt(x, y).Float()
}

// Vec is shorthand for geom.Vec(x, y).Float()
func Vec[T geom.Number](x, y T) Vector {
	return geom.Vec(x, y).Float()
}

// Sz is shorthand for geom.Sz(width, height).Float()
func Sz[T geom.Number](width, height T) Size {
	return geom.Sz(width, height).Float()
}

// Circ is shorthand for geom.Circ(center, radius).Float()
func Circ[T geom.Number](center geom.Point[T], radius T) Circle {
	return geom.Circ(center, radius).Float()
}

// Ln is shorthand for geom.Ln(start, end).Float()
func Ln[T geom.Number](start, end geom.Point[T]) Line {
	return geom.Ln(start, end).Float()
}

// Rect is shorthand for geom.Rect(center, size).Float()
func Rect[T geom.Number](center geom.Point[T], size geom.Size[T]) Rectangle {
	return geom.Rect(center, size).Float()
}

// Pol is shorthand for geom.Pol(vertices).Float()
func Pol[T geom.Number](vertices []geom.Point[T]) Polygon {
	return geom.Pol(vertices).Float()
}

// RegPol is shorthand for geom.RegPol(center, size, n, angle).Float()
func RegPol[T geom.Number](center geom.Point[T], size geom.Size[T], n int, angle float64) RegularPolygon {
	return geom.RegPol(center, size, n, angle).Float()
}

// Pad is shorthand for geom.Pad(top, right, bottom, left).Float()
func Pad[T geom.Number](top, right, bottom, left T) Padding {
	return geom.Pad(top, right, bottom, left).Float()
}

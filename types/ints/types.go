// Package ints provides aliases of the geom types with an [int] type argument, and shorthand
// constructors that take any Number and convert through Int, which rounds and panics on a
// non-finite value.
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

// Padding is geom.Padding[int].
type Padding = geom.Padding[int]

// Matrix is geom.Matrix[int].
type Matrix = geom.Matrix[int]

// Circle is geom.Circle[int].
type Circle = geom.Circle[int]

// Capsule is geom.Capsule[int].
type Capsule = geom.Capsule[int]

// Sector is geom.Sector[int].
type Sector = geom.Sector[int]

// Arc is geom.Arc[int].
type Arc = geom.Arc[int]

// Segment is geom.Segment[int].
type Segment = geom.Segment[int]

// Ray is geom.Ray[int].
type Ray = geom.Ray[int]

// Line is geom.Line[int].
type Line = geom.Line[int]

// Polyline is geom.Polyline[int].
type Polyline = geom.Polyline[int]

// QuadraticBezier is geom.QuadraticBezier[int].
type QuadraticBezier = geom.QuadraticBezier[int]

// CubicBezier is geom.CubicBezier[int].
type CubicBezier = geom.CubicBezier[int]

// Polygon is geom.Polygon[int].
type Polygon = geom.Polygon[int]

// Rectangle is geom.Rectangle[int].
type Rectangle = geom.Rectangle[int]

// RegularPolygon is geom.RegularPolygon[int].
type RegularPolygon = geom.RegularPolygon[int]

// Box is geom.Box[int].
type Box = geom.Box[int]

// Ellipse is geom.Ellipse[int].
type Ellipse = geom.Ellipse[int]

// Hit is geom.Hit[int].
type Hit = geom.Hit[int]

// Contact is geom.Contact[int].
type Contact = geom.Contact[int]

// Pt is shorthand for geom.Pt(x, y).Int().
func Pt[T geom.Number](x, y T) Point {
	return geom.Pt(x, y).Int()
}

// Vec is shorthand for geom.Vec(x, y).Int().
func Vec[T geom.Number](x, y T) Vector {
	return geom.Vec(x, y).Int()
}

// Sz is shorthand for geom.Sz(width, height).Int().
func Sz[T geom.Number](width, height T) Size {
	return geom.Sz(width, height).Int()
}

// Pad is shorthand for geom.Pad(top, right, bottom, left).Int().
func Pad[T geom.Number](top, right, bottom, left T) Padding {
	return geom.Pad(top, right, bottom, left).Int()
}

// Mat is shorthand for geom.Mat(a, b, c, d, e, f).Int().
func Mat[T geom.Number](a, b, c, d, e, f T) Matrix {
	return geom.Mat(a, b, c, d, e, f).Int()
}

// IdentityMatrix is shorthand for geom.IdentityMatrix[int]().
func IdentityMatrix() Matrix {
	return geom.IdentityMatrix[int]()
}

// Circ is shorthand for geom.Circ(center, radius).Int().
func Circ[T geom.Number](center geom.Point[T], radius T) Circle {
	return geom.Circ(center, radius).Int()
}

// Cap is shorthand for geom.Cap(start, end, radius).Int().
func Cap[T geom.Number](start, end geom.Point[T], radius T) Capsule {
	return geom.Cap(start, end, radius).Int()
}

// Sect is shorthand for geom.Sect(center, radius, angle, sweep).Int().
func Sect[T geom.Number](center geom.Point[T], radius T, angle, sweep float64) Sector {
	return geom.Sect(center, radius, angle, sweep).Int()
}

// Ac is shorthand for geom.Ac(center, radius, angle, sweep).Int().
func Ac[T geom.Number](center geom.Point[T], radius T, angle, sweep float64) Arc {
	return geom.Ac(center, radius, angle, sweep).Int()
}

// Seg is shorthand for geom.Seg(start, end).Int().
func Seg[T geom.Number](start, end geom.Point[T]) Segment {
	return geom.Seg(start, end).Int()
}

// Ry is shorthand for geom.Ry(origin, vector).Int().
func Ry[T geom.Number](origin geom.Point[T], vector geom.Vector[T]) Ray {
	return geom.Ry(origin, vector).Int()
}

// Ln is shorthand for geom.Ln(origin, vector).Int().
func Ln[T geom.Number](origin geom.Point[T], vector geom.Vector[T]) Line {
	return geom.Ln(origin, vector).Int()
}

// PolLn is shorthand for geom.PolLn(points).Int().
func PolLn[T geom.Number](points []geom.Point[T]) Polyline {
	return geom.PolLn(points).Int()
}

// QuadBez is shorthand for geom.QuadBez(start, control, end).Int().
func QuadBez[T geom.Number](start, control, end geom.Point[T]) QuadraticBezier {
	return geom.QuadBez(start, control, end).Int()
}

// CubBez is shorthand for geom.CubBez(start, control1, control2, end).Int().
func CubBez[T geom.Number](start, control1, control2, end geom.Point[T]) CubicBezier {
	return geom.CubBez(start, control1, control2, end).Int()
}

// Pol is shorthand for geom.Pol(vertices).Int().
func Pol[T geom.Number](vertices []geom.Point[T]) Polygon {
	return geom.Pol(vertices).Int()
}

// Rect is shorthand for geom.Rect(center, size, angle).Int().
func Rect[T geom.Number](center geom.Point[T], size geom.Size[T], angle float64) Rectangle {
	return geom.Rect(center, size, angle).Int()
}

// RegPol is shorthand for geom.RegPol(center, radius, n, angle, phase).Int().
func RegPol[T geom.Number](center geom.Point[T], radius geom.Size[T], n int, angle, phase float64) RegularPolygon {
	return geom.RegPol(center, radius, n, angle, phase).Int()
}

// Bx is shorthand for geom.Bx(a, b).Int().
func Bx[T geom.Number](a, b geom.Point[T]) Box {
	return geom.Bx(a, b).Int()
}

// Ell is shorthand for geom.Ell(center, radius, angle).Int().
func Ell[T geom.Number](center geom.Point[T], radius geom.Size[T], angle float64) Ellipse {
	return geom.Ell(center, radius, angle).Int()
}

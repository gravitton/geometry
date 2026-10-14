// Package floats provides aliases of the geom types with a [float64] type argument, and shorthand
// constructors that take any Number and convert through Float.
package floats

import (
	geom "github.com/gravitton/geometry"
)

// Point is geom.Point[float64].
type Point = geom.Point[float64]

// Vector is geom.Vector[float64].
type Vector = geom.Vector[float64]

// Size is geom.Size[float64].
type Size = geom.Size[float64]

// Padding is geom.Padding[float64].
type Padding = geom.Padding[float64]

// Matrix is geom.Matrix[float64].
type Matrix = geom.Matrix[float64]

// Circle is geom.Circle[float64].
type Circle = geom.Circle[float64]

// Capsule is geom.Capsule[float64].
type Capsule = geom.Capsule[float64]

// Sector is geom.Sector[float64].
type Sector = geom.Sector[float64]

// Arc is geom.Arc[float64].
type Arc = geom.Arc[float64]

// Segment is geom.Segment[float64].
type Segment = geom.Segment[float64]

// Ray is geom.Ray[float64].
type Ray = geom.Ray[float64]

// Line is geom.Line[float64].
type Line = geom.Line[float64]

// Polyline is geom.Polyline[float64].
type Polyline = geom.Polyline[float64]

// QuadraticBezier is geom.QuadraticBezier[float64].
type QuadraticBezier = geom.QuadraticBezier[float64]

// CubicBezier is geom.CubicBezier[float64].
type CubicBezier = geom.CubicBezier[float64]

// Polygon is geom.Polygon[float64].
type Polygon = geom.Polygon[float64]

// Rectangle is geom.Rectangle[float64].
type Rectangle = geom.Rectangle[float64]

// RegularPolygon is geom.RegularPolygon[float64].
type RegularPolygon = geom.RegularPolygon[float64]

// Box is geom.Box[float64].
type Box = geom.Box[float64]

// Ellipse is geom.Ellipse[float64].
type Ellipse = geom.Ellipse[float64]

// Hit is geom.Hit[float64].
type Hit = geom.Hit[float64]

// Contact is geom.Contact[float64].
type Contact = geom.Contact[float64]

// Pt is shorthand for geom.Pt(x, y).Float().
func Pt[T geom.Number](x, y T) Point {
	return geom.Pt(x, y).Float()
}

// Vec is shorthand for geom.Vec(x, y).Float().
func Vec[T geom.Number](x, y T) Vector {
	return geom.Vec(x, y).Float()
}

// Sz is shorthand for geom.Sz(width, height).Float().
func Sz[T geom.Number](width, height T) Size {
	return geom.Sz(width, height).Float()
}

// Pad is shorthand for geom.Pad(top, right, bottom, left).Float().
func Pad[T geom.Number](top, right, bottom, left T) Padding {
	return geom.Pad(top, right, bottom, left).Float()
}

// Mat is shorthand for geom.Mat(a, b, c, d, e, f).Float().
func Mat[T geom.Number](a, b, c, d, e, f T) Matrix {
	return geom.Mat(a, b, c, d, e, f).Float()
}

// IdentityMatrix is shorthand for geom.IdentityMatrix[float64]().
func IdentityMatrix() Matrix {
	return geom.IdentityMatrix[float64]()
}

// Circ is shorthand for geom.Circ(center, radius).Float().
func Circ[T geom.Number](center geom.Point[T], radius T) Circle {
	return geom.Circ(center, radius).Float()
}

// Cap is shorthand for geom.Cap(start, end, radius).Float().
func Cap[T geom.Number](start, end geom.Point[T], radius T) Capsule {
	return geom.Cap(start, end, radius).Float()
}

// Sect is shorthand for geom.Sect(center, radius, angle, sweep).Float().
func Sect[T geom.Number](center geom.Point[T], radius T, angle, sweep float64) Sector {
	return geom.Sect(center, radius, angle, sweep).Float()
}

// Ac is shorthand for geom.Ac(center, radius, angle, sweep).Float().
func Ac[T geom.Number](center geom.Point[T], radius T, angle, sweep float64) Arc {
	return geom.Ac(center, radius, angle, sweep).Float()
}

// Seg is shorthand for geom.Seg(start, end).Float().
func Seg[T geom.Number](start, end geom.Point[T]) Segment {
	return geom.Seg(start, end).Float()
}

// Ry is shorthand for geom.Ry(origin, direction).Float().
func Ry[T geom.Number](origin geom.Point[T], direction geom.Vector[T]) Ray {
	return geom.Ry(origin, direction).Float()
}

// Ln is shorthand for geom.Ln(origin, direction).Float().
func Ln[T geom.Number](origin geom.Point[T], direction geom.Vector[T]) Line {
	return geom.Ln(origin, direction).Float()
}

// PolLn is shorthand for geom.PolLn(points).Float().
func PolLn[T geom.Number](points []geom.Point[T]) Polyline {
	return geom.PolLn(points).Float()
}

// QuadBez is shorthand for geom.QuadBez(start, control, end).Float().
func QuadBez[T geom.Number](start, control, end geom.Point[T]) QuadraticBezier {
	return geom.QuadBez(start, control, end).Float()
}

// CubBez is shorthand for geom.CubBez(start, control1, control2, end).Float().
func CubBez[T geom.Number](start, control1, control2, end geom.Point[T]) CubicBezier {
	return geom.CubBez(start, control1, control2, end).Float()
}

// Pol is shorthand for geom.Pol(vertices).Float().
func Pol[T geom.Number](vertices []geom.Point[T]) Polygon {
	return geom.Pol(vertices).Float()
}

// Rect is shorthand for geom.Rect(center, size).Float().
func Rect[T geom.Number](center geom.Point[T], size geom.Size[T]) Rectangle {
	return geom.Rect(center, size).Float()
}

// RegPol is shorthand for geom.RegPol(center, size, n, angle, phase).Float().
func RegPol[T geom.Number](center geom.Point[T], size geom.Size[T], n int, angle, phase float64) RegularPolygon {
	return geom.RegPol(center, size, n, angle, phase).Float()
}

// Bx is shorthand for geom.Bx(a, b).Float().
func Bx[T geom.Number](a, b geom.Point[T]) Box {
	return geom.Bx(a, b).Float()
}

// Ell is shorthand for geom.Ell(center, size, angle).Float().
func Ell[T geom.Number](center geom.Point[T], size geom.Size[T], angle float64) Ellipse {
	return geom.Ell(center, size, angle).Float()
}

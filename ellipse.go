package geom

import (
	"math"
	"slices"
)

// Ellipse is an ellipse of a center and semi-axes, turned by Angle radians about the center.
//
// Size holds the semi-axes, not the extent, named in the frame before the turn: Size.Width is the
// semi-axis that pointed along X. Equal compares the angle even where equal semi-axes make it a
// circle that no angle moves.
//
// The size is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative semi-axis written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it. A zero semi-axis is not repaired: the ellipse is then the segment the
// other axis spans, and answers as that segment.
//
// The ellipse has no Intersects methods and no vertices. Intersects, Distance, Collide and Sweep
// test it by its Support, or by its distance to the center of a Circle, Segment.HitEllipse and
// Ray.HitEllipse cast against it, and RegularPolygon and Polygon give its outline.
type Ellipse[T Number] struct {
	Center Point[T] `json:",embed"`
	Size   Size[T]  `json:",embed"`
	Angle  float64  `json:"a,omitzero"`
}

// Ell is shorthand for Ellipse{center, size, angle}, with the size taken absolute.
func Ell[T Number](center Point[T], size Size[T], angle float64) Ellipse[T] {
	return Ellipse[T]{center, size.Abs(), angle}
}

// SemiMajor returns the longer semi-axis, the one the foci lie on.
func (e Ellipse[T]) SemiMajor() T {
	return max(e.Size.Width, e.Size.Height)
}

// SemiMinor returns the shorter semi-axis.
func (e Ellipse[T]) SemiMinor() T {
	return min(e.Size.Width, e.Size.Height)
}

// Eccentricity returns how far the ellipse departs from a circle, in [0, 1]: zero for a circle or a
// point, one where the shorter semi-axis is zero.
func (e Ellipse[T]) Eccentricity() float64 {
	a, b := float64(e.SemiMajor()), float64(e.SemiMinor())
	if a == 0 {
		return 0
	}

	return math.Sqrt(1 - b*b/(a*a))
}

// Foci returns the two focal points on the major axis, the one on its negative end before the turn
// first. A circle has both at the center.
func (e Ellipse[T]) Foci() (Point[T], Point[T]) {
	w, h := e.Size.Float().Abs().XY()
	a, b := max(w, h), min(w, h)

	focal := Vector[float64]{math.Sqrt(float64(a*a) - float64(b*b)), 0}
	if h > w {
		focal = Vector[float64]{0, focal.X}
	}

	return e.worldPoint(focal.Negate()), e.worldPoint(focal)
}

// Anchor returns the point of the boundary in the direction from the center, named in the frame
// before the turn, or the center for DirectionNone. A diagonal anchor lies where the diagonal
// leaves the ellipse, and a degenerate ellipse answers with the point of its segment.
func (e Ellipse[T]) Anchor(direction Direction) Point[T] {
	if direction.IsNone() {
		return e.Center
	}

	return e.worldPoint(VectorFromAngleSize(e.parametricAngle(direction.Angle()), e.Size.Float()))
}

// Centroid returns the center.
func (e Ellipse[T]) Centroid() Point[T] {
	return e.Center
}

// Area returns the area (π·w·h).
func (e Ellipse[T]) Area() float64 {
	return Pi * float64(e.Size.Width) * float64(e.Size.Height)
}

// Perimeter returns the length of the boundary by Ramanujan's second approximation, exact for a
// circle and within a part in a thousand for a degenerate ellipse.
func (e Ellipse[T]) Perimeter() float64 {
	w, h := e.Size.Float().XY()

	sum := w + h
	if sum == 0 {
		return 0
	}

	r := (w - h) / sum
	t := float64(3 * r * r)

	return Pi * sum * (1 + t/(10+math.Sqrt(4-t)))
}

// Inertia returns the polar second moment of area about the center (π·w·h·(w²+h²)/4), the
// rotational inertia of the ellipse at unit density.
func (e Ellipse[T]) Inertia() float64 {
	w, h := e.Size.Float().XY()

	return Pi * w * h * (float64(w*w) + float64(h*h)) / 4
}

// Bounds returns the axis-aligned bounding box.
func (e Ellipse[T]) Bounds() Box[T] {
	a, b := e.minMax()

	return Box[T]{a, b}
}

// BoundingCircle returns the smallest circle enclosing the ellipse, the one Circle returns.
func (e Ellipse[T]) BoundingCircle() Circle[T] {
	return e.Circle()
}

// worldPoint returns the point at the offset from the center in the frame before the turn, rotated
// once and rounded once.
func (e Ellipse[T]) worldPoint(offset Vector[float64]) Point[T] {
	return e.Center.Float().Add(offset.Rotate(e.Angle)).Cast[T]()
}

// localOffset returns the offset of the point from the center in the frame before the turn, the
// inverse of worldPoint and the frame every distance is measured in.
func (e Ellipse[T]) localOffset(point Point[T]) Vector[float64] {
	return point.Float().Subtract(e.Center.Float()).Rotate(-e.Angle)
}

// parametricAngle returns the parameter of the boundary point in the direction from the center,
// which differs from the direction unless the ellipse is a circle. A degenerate ellipse keeps the
// direction, so the point lands on its segment.
func (e Ellipse[T]) parametricAngle(direction float64) float64 {
	w, h := e.Size.Float().XY()
	if w == 0 || h == 0 {
		return direction
	}

	sin, cos := math.Sincos(direction)

	return math.Atan2(w*sin, h*cos)
}

// extent returns the half extent of the turned ellipse along each world axis, in closed form, since
// the boundary touches the box where its tangent is axis-aligned.
func (e Ellipse[T]) extent() Vector[float64] {
	w, h := e.Size.Float().XY()
	sin, cos := math.Sincos(e.Angle)

	return Vector[float64]{math.Hypot(w*cos, h*sin), math.Hypot(w*sin, h*cos)}
}

// minMax returns the corners of Bounds, the center less and plus its extent.
func (e Ellipse[T]) minMax() (Point[T], Point[T]) {
	center, extent := e.Center.Float(), e.extent()

	return center.Add(extent.Negate()).Cast[T](), center.Add(extent).Cast[T]()
}

// Translate returns the ellipse moved by the vector.
func (e Ellipse[T]) Translate(vector Vector[T]) Ellipse[T] {
	return Ellipse[T]{e.Center.Add(vector), e.Size, e.Angle}
}

// MoveTo returns the ellipse with its center at the point.
func (e Ellipse[T]) MoveTo(point Point[T]) Ellipse[T] {
	return Ellipse[T]{point, e.Size, e.Angle}
}

// AlignTo returns the ellipse moved so that its Anchor in the direction lies at the point.
func (e Ellipse[T]) AlignTo(direction Direction, point Point[T]) Ellipse[T] {
	return e.Translate(point.Subtract(e.Anchor(direction)))
}

// Rotate returns the ellipse turned by the angle in radians about its center, the stored angle
// normalized to [0, 2π).
func (e Ellipse[T]) Rotate(angle float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size, NormalizeAngle(e.Angle + angle)}
}

// Scale returns the ellipse with both semi-axes scaled by the factor, taken absolute.
func (e Ellipse[T]) Scale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Scale(factor).Abs(), e.Angle}
}

// ScaleXY returns the ellipse with its semi-axes scaled by the factors along its own axes, taken
// absolute.
func (e Ellipse[T]) ScaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.ScaleXY(factorX, factorY).Abs(), e.Angle}
}

// Unscale returns the ellipse with both semi-axes divided by the factor, taken absolute. It panics
// for a zero factor.
func (e Ellipse[T]) Unscale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Unscale(factor).Abs(), e.Angle}
}

// UnscaleXY returns the ellipse with its semi-axes divided by the factors along its own axes, taken
// absolute. It panics for a zero factor.
func (e Ellipse[T]) UnscaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.UnscaleXY(factorX, factorY).Abs(), e.Angle}
}

// Resize returns the ellipse with the semi-axes, taken absolute.
func (e Ellipse[T]) Resize(size Size[T]) Ellipse[T] {
	return Ellipse[T]{e.Center, size.Abs(), e.Angle}
}

// Transform returns the ellipse mapped by the matrix, exactly, since an affine map takes an ellipse
// to an ellipse. A reflection mirrors the angle, and where the semi-axes come out swapped the angle
// takes a quarter turn.
func (e Ellipse[T]) Transform[M Float](matrix Matrix[M]) Ellipse[T] {
	after, stretch, _, _ := matrix.mapping(e.Angle, e.Size.Float()).decomposition()

	return Ellipse[T]{e.Center.Transform(matrix), stretch.Size().Cast[T](), NormalizeAngle(after)}
}

// Grow returns the ellipse with both semi-axes increased by the amount, clamped at zero.
func (e Ellipse[T]) Grow(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Grow(amount).AtLeastZero(), e.Angle}
}

// GrowXY returns the ellipse with its semi-axes increased by the amounts along its own axes,
// clamped at zero.
func (e Ellipse[T]) GrowXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.GrowXY(amountX, amountY).AtLeastZero(), e.Angle}
}

// Shrink returns the ellipse with both semi-axes decreased by the amount, clamped at zero.
func (e Ellipse[T]) Shrink(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Shrink(amount).AtLeastZero(), e.Angle}
}

// ShrinkXY returns the ellipse with its semi-axes decreased by the amounts along its own axes,
// clamped at zero.
func (e Ellipse[T]) ShrinkXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.ShrinkXY(amountX, amountY).AtLeastZero(), e.Angle}
}

// Lerp returns the ellipse interpolated toward the given one, center and semi-axes together and the
// angle along the shorter arc by LerpAngle, normalized to [0, 2π), and extrapolated outside [0, 1].
// The size is taken absolute, so extrapolating past a zero semi-axis grows the ellipse again.
func (e Ellipse[T]) Lerp(ellipse Ellipse[T], t float64) Ellipse[T] {
	return Ellipse[T]{e.Center.Lerp(ellipse.Center, t), e.Size.Lerp(ellipse.Size, t).Abs(), NormalizeAngle(LerpAngle(e.Angle, ellipse.Angle, t))}
}

// Canonical returns the ellipse with its size taken absolute and its angle normalized to [0, 2π),
// an angle within Delta of zero or a full turn becoming exactly zero. Rotate and Lerp never snap.
func (e Ellipse[T]) Canonical() Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Abs(), snapAngle(e.Angle)}
}

// DistanceTo returns the distance from the point to the ellipse, zero inside.
func (e Ellipse[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(e.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the square of DistanceTo.
func (e Ellipse[T]) DistanceSquaredTo(point Point[T]) float64 {
	_, distanceSquared := e.nearestLocal(point)

	return distanceSquared
}

// Nearest returns the point of the ellipse nearest to the given point, the point itself inside.
func (e Ellipse[T]) Nearest(point Point[T]) Point[T] {
	foot, distanceSquared := e.nearestLocal(point)
	if distanceSquared == 0 {
		return point
	}

	return e.worldPoint(foot)
}

// Support returns the point of the ellipse farthest along the direction, whatever its length. A
// degenerate ellipse answers with an end of its segment, or its center for a direction across it,
// and a zero direction with the end of the width semi-axis.
func (e Ellipse[T]) Support(direction Vector[float64]) Point[float64] {
	w, h := e.Size.Float().XY()
	local := direction.Rotate(-e.Angle)
	unit := Vector[float64]{w * local.X, h * local.Y}.Normalize()

	return e.Center.Float().Add(Vector[float64]{w * unit.X, h * unit.Y}.Rotate(e.Angle))
}

// Contains reports whether the point lies in the ellipse. A degenerate ellipse contains the points
// of its segment.
func (e Ellipse[T]) Contains(point Point[T]) bool {
	return e.DistanceSquaredTo(point) == 0
}

// nearestLocal returns the nearest point of the ellipse as an offset in the local frame, with the
// squared distance to it: the offset itself at zero inside, so DistanceSquaredTo and Nearest agree.
func (e Ellipse[T]) nearestLocal(point Point[T]) (Vector[float64], float64) {
	local := e.localOffset(point)
	if e.containsOffset(local) {
		return local, 0
	}

	foot := e.nearestOffset(local)
	distanceSquared := local.Subtract(foot).LengthSquared()
	if lessOrEqualSquared(distanceSquared, 0, epsilonAt[T](max(e.Center.magnitude(), point.magnitude()))) {
		return local, 0
	}

	return foot, distanceSquared
}

// containsOffset reports whether the local offset lies inside, on the quadratic form multiplied
// through by the squared semi-axes so no division rounds it. A degenerate ellipse has no inside and
// leaves every point, NaN included, to the distance.
func (e Ellipse[T]) containsOffset(local Vector[float64]) bool {
	w, h := e.Size.Float().XY()
	if w == 0 || h == 0 {
		return false
	}

	x, y, extent := local.X*h, local.Y*w, w*h

	return float64(x*x)+float64(y*y) <= extent*extent
}

// nearestOffset returns the foot of the perpendicular from a local offset outside the ellipse. For
// a degenerate ellipse it is the offset clamped to the segment on each axis.
func (e Ellipse[T]) nearestOffset(local Vector[float64]) Vector[float64] {
	w, h := e.Size.Float().Abs().XY()
	if w == 0 || h == 0 {
		return Vector[float64]{Clamp(local.X, -w, w), Clamp(local.Y, -h, h)}
	}

	x, y := math.Abs(local.X), math.Abs(local.Y)

	var px, py float64
	if w >= h {
		px, py = e.foot(w, h, x, y)
	} else {
		py, px = e.foot(h, w, y, x)
	}

	return Vector[float64]{math.Copysign(px, local.X), math.Copysign(py, local.Y)}
}

// foot returns the foot of the perpendicular from (x, y) outside the ellipse of semi-axes a ≥ b, in
// the first quadrant, which nearestOffset reflects every point into. A point on an axis is answered
// directly; elsewhere the foot is placed from the root.
func (e Ellipse[T]) foot(a, b, x, y float64) (float64, float64) {
	if y == 0 {
		return a, 0
	}

	if x == 0 {
		return 0, b
	}

	aspect := float64((a / b) * (a / b))
	s := e.root(aspect, x/a, y/b)

	return aspect * x / (s + aspect), y / (s + 1)
}

// root returns the parameter at which the perpendicular from (x, y), in units of the semi-axes,
// meets the boundary, aspect being the squared ratio of the semi-axes. It bisects the bracket the
// point gives until no float lies between its ends, so no iteration count is chosen, and a NaN ends
// the loop at once.
func (Ellipse[T]) root(aspect, x, y float64) float64 {
	gradient := aspect * x

	s0, s1 := y-1, math.Hypot(gradient, y)-1

	for {
		s := (s0 + s1) / 2
		if s == s0 || s == s1 || math.IsNaN(s) {
			return s
		}

		if px, py := gradient/(s+aspect), y/(s+1); float64(px*px)+float64(py*py) > 1 {
			s0 = s
		} else {
			s1 = s
		}
	}
}

// Equal reports whether the ellipses have equal center and size and an angle equal by EqualAngle.
func (e Ellipse[T]) Equal(ellipse Ellipse[T]) bool {
	return e.Center.Equal(ellipse.Center) && e.Size.Equal(ellipse.Size) && EqualAngle(e.Angle, ellipse.Angle)
}

// IsZero reports whether the ellipse is the zero value, an angle of a full turn included.
func (e Ellipse[T]) IsZero() bool {
	return e.Equal(Ellipse[T]{})
}

// IsAligned reports whether the angle is exactly zero. Canonical snaps a residue to zero.
func (e Ellipse[T]) IsAligned() bool {
	return e.Angle == 0
}

// IsCircle reports whether the semi-axes are equal, where Circle returns the same shape.
func (e Ellipse[T]) IsCircle() bool {
	return Equal(e.Size.Width, e.Size.Height)
}

// Circle returns the circle of the major semi-axis about the center, the smallest enclosing the
// ellipse. It is the same shape only where IsCircle holds, and there the inverse of Circle.Ellipse.
func (e Ellipse[T]) Circle() Circle[T] {
	return Circle[T]{e.Center, e.SemiMajor()}
}

// RegularPolygon returns the regular polygon of n vertices inscribed in the ellipse, of its center,
// semi-axes and angle, its first vertex placed by RegularPolygonOrientationPhase before the
// stretch. It is empty for n below one, and the only polygon that RegularPolygon.Ellipse converts
// back exactly. It panics for an orientation that is neither OrientationFlatTop nor
// OrientationPointyTop.
//
// For the polygon whose edges touch the ellipse, scale it first by the circumradius over the
// apothem:
//
//	e.Scale(1 / math.Cos(Pi / float64(n))).RegularPolygon(n, orientation)
func (e Ellipse[T]) RegularPolygon(n int, orientation Orientation) RegularPolygon[T] {
	return RegularPolygon[T]{e.Center, e.Size, n, e.Angle, RegularPolygonOrientationPhase(n, orientation)}
}

// Polygon returns the ellipse as a polygon within the tolerance of it: equal steps of the parameter
// from the end of the width semi-axis, as many as the circle of the major semi-axis needs, at least
// three and at most 1024, wound like Rectangle.Vertices. A negative tolerance is taken
// absolute, and an ellipse of no extent gives its center alone. RegularPolygon chooses the vertex
// count instead.
func (e Ellipse[T]) Polygon(tolerance float64) Polygon[T] {
	return Polygon[T]{e.AppendPolygon(nil, tolerance)}
}

// AppendPolygon appends the vertices of Polygon(tolerance) to dst and returns the extended slice.
func (e Ellipse[T]) AppendPolygon(dst []Point[T], tolerance float64) []Point[T] {
	size := e.Size.Float()
	pieces := max(Arc[float64]{Radius: max(size.Width, size.Height), Sweep: 2 * math.Pi}.pieces(tolerance), 3)

	from := len(dst)
	dst = slices.Grow(dst, int(pieces))

	for i := 0.0; i < pieces; i++ {
		if point := e.worldPoint(VectorFromAngleSize(2*math.Pi*i/pieces, size)); len(dst) == from || point != dst[len(dst)-1] {
			dst = append(dst, point)
		}
	}

	if last := len(dst) - 1; last > from && dst[last].Equal(dst[from]) {
		dst = dst[:last]
	}

	return dst
}

// Cast converts the ellipse to an Ellipse of another number type, rounding as Cast does.
func (e Ellipse[T]) Cast[R Number]() Ellipse[R] {
	return Ellipse[R]{e.Center.Cast[R](), e.Size.Cast[R](), e.Angle}
}

// Int converts the ellipse to an Ellipse[int].
func (e Ellipse[T]) Int() Ellipse[int] {
	return Ellipse[int]{e.Center.Int(), e.Size.Int(), e.Angle}
}

// Float converts the ellipse to an Ellipse[float64].
func (e Ellipse[T]) Float() Ellipse[float64] {
	return Ellipse[float64]{e.Center.Float(), e.Size.Float(), e.Angle}
}

// String returns the ellipse in the form of its constructor: Ell((x,y);WxH), with ;a before the
// closing parenthesis for a turned ellipse.
func (e Ellipse[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Ell("...)
	buf = e.Center.appendString(buf)
	buf = append(buf, ';')
	buf = e.Size.appendString(buf)
	if !e.IsAligned() {
		buf = append(buf, ';')
		buf = appendString(buf, e.Angle)
	}
	buf = append(buf, ')')

	return string(buf)
}

package geom

import (
	"fmt"
	"math"
)

// Ellipse is a 2D ellipse represented by its center, semi-axes and angle: the ellipse of those
// semi-axes about the center, turned by Angle radians about it in the same sense as
// Vector.Rotate. Ell builds one, and Rotate turns it.
//
// Size holds the semi-axes, so it is a radius, not an extent, the meaning it has on
// RegularPolygon and the one Circle.Radius has: an ellipse of Size 10x4 spans 20x8, and the
// circle of radius r is the ellipse of Size r x r. Use Bounds for the extent.
//
// Angle turns the ellipse about its center, the same meaning it has on Rectangle and
// RegularPolygon: the semi-axes are named in the frame before the turn, so Size.Width is the
// semi-axis that pointed along X before the ellipse was turned, wherever it now lies. An
// ellipse of equal semi-axes is a circle and every angle leaves it where it is, though Equal
// still compares the angle, as it does on RegularPolygon.
//
// The size is never negative: Ell and Resize take it absolute, Scale takes a negative factor
// absolute and Grow and Shrink clamp at zero, since an ellipse mirrored about its center is
// the same ellipse. A negative semi-axis can only be written as a struct literal or decoded
// from JSON, and Contains, DistanceTo and Bounds give no meaningful answer for one; Canonical
// repairs it.
//
// A zero semi-axis is not repaired: it is the degenerate ellipse, the segment the other axis
// spans, and Contains, DistanceTo and Bounds answer for it as for that segment.
//
// The ellipse is a Shape and a Transformable, and deliberately not an Outline or a
// Collider: it has no vertices, and two ellipses meet at the roots of a quartic rather than at
// anything the circle pairs are built on. Both are left to RegularPolygon, which holds the
// same center, semi-axes and angle: RegularPolygon(n) is the polygon of the wanted resolution
// inscribed in the ellipse, and RegularPolygon.Ellipse converts back.
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

// SemiMinor returns the shorter semi-axis, the one perpendicular to the foci.
func (e Ellipse[T]) SemiMinor() T {
	return min(e.Size.Width, e.Size.Height)
}

// Eccentricity returns how far the ellipse departs from a circle, in [0, 1]: zero for a circle,
// where the two semi-axes are equal, and one for a degenerate ellipse, where the shorter one is
// zero and the ellipse is a segment. An ellipse of no extent at all is a point, which is a
// circle, and gives zero.
func (e Ellipse[T]) Eccentricity() float64 {
	a, b := float64(e.SemiMajor()), float64(e.SemiMinor())
	if a == 0 {
		return 0
	}

	return math.Sqrt(1 - b*b/(a*a))
}

// Foci returns the two focal points, the pair whose distances to a point of the boundary sum to
// twice the major semi-axis: they lie on the major axis, either side of the center and turned
// by Angle with it, the nearer one first as minMax gives its corners. A circle has both at the
// center. For integer T each focus is rounded like every other point placed from an angle.
func (e Ellipse[T]) Foci() (Point[T], Point[T]) {
	a, b := float64(e.SemiMajor()), float64(e.SemiMinor())

	focal := Vector[float64]{math.Sqrt(a*a - b*b), 0}
	if e.Size.Height > e.Size.Width {
		focal = Vector[float64]{0, focal.X}
	}

	return e.worldPoint(focal.Negate()), e.worldPoint(focal)
}

// Anchor returns the point on the boundary in the given direction from the center, named in
// the frame of the ellipse before its turn, or the center itself for DirectionNone: the end of
// a semi-axis for a cardinal direction, and the point of the boundary at an eighth of a turn
// for a diagonal, which is the corner of Bounds scaled onto the ellipse rather than the corner
// itself. Unlike Circle.Anchor it lies on the boundary at every angle, since the point is
// placed on the ellipse rather than rounded from a lattice step.
func (e Ellipse[T]) Anchor(direction Direction) Point[T] {
	if direction.IsNone() {
		return e.Center
	}

	return e.worldPoint(VectorFromAngleSize(direction.Angle(), e.Size.Float()))
}

// Centroid returns the center of the enclosed area, the Center of the ellipse.
func (e Ellipse[T]) Centroid() Point[T] {
	return e.Center
}

// Area returns the area enclosed by the boundary (π * width * height). It is a float64 even
// for an integer T, since the factor π leaves no pair of semi-axes with an area T could
// express, as Circle.Area is.
func (e Ellipse[T]) Area() float64 {
	return Pi * float64(e.Size.Width) * float64(e.Size.Height)
}

// Perimeter returns the length of the boundary, by Ramanujan's second approximation: the
// circumference of an ellipse has no closed form, and the approximation is exact for a circle,
// within a rounding error of the true value for the sizes a shape of a screen has, and
// loosest, still within a part in a thousand, for the degenerate ellipse, where the true
// perimeter is four times the major semi-axis. An ellipse of no extent has a perimeter of zero.
func (e Ellipse[T]) Perimeter() float64 {
	w, h := e.Size.Float().XY()

	sum := w + h
	if sum == 0 {
		return 0
	}

	ratio := (w - h) / sum
	t := 3 * ratio * ratio

	return Pi * sum * (1 + t/(10+math.Sqrt(4-t)))
}

// Inertia returns the polar second moment of area about the center, the rotational inertia
// of the enclosed area at unit density: π * w * h * (w^2 + h^2) / 4 for the semi-axes w and h,
// the sum of the moments about each axis, which the turn does not change.
func (e Ellipse[T]) Inertia() float64 {
	w, h := e.Size.Float().XY()

	return Pi * w * h * (w*w + h*h) / 4
}

// Bounds returns the axis-aligned bounding rectangle: the box on the corners minMax finds,
// which touches the boundary at four points whatever the angle.
func (e Ellipse[T]) Bounds() Rectangle[T] {
	return RectangleFromMinMax(e.minMax())
}

// worldPoint returns the point at the given offset from the center in the frame before the
// turn, turned by Angle, as RegularPolygon.worldPoint places a vertex: the offset is rotated
// once and the sum rounded once for an integer T.
func (e Ellipse[T]) worldPoint(offset Vector[float64]) Point[T] {
	return e.Center.Float().Add(offset.Rotate(e.Angle)).Cast[T]()
}

// localOffset returns the offset of the point from the center in the frame of the ellipse
// before its turn, in float64, the inverse of worldPoint: the frame the boundary is a plain
// ellipse of the semi-axes in, where every distance is measured.
func (e Ellipse[T]) localOffset(point Point[T]) Vector[float64] {
	return point.Subtract(e.Center).Float().Rotate(-e.Angle)
}

// extent returns half the axis-aligned extent of the ellipse, in float64: the reach of the
// turned ellipse along each axis of the world, in closed form, since the boundary touches the
// box where its tangent is axis-aligned rather than at a vertex. An ellipse that is not turned
// reaches its own semi-axes.
func (e Ellipse[T]) extent() Vector[float64] {
	w, h := e.Size.Float().XY()
	sin, cos := math.Sincos(e.Angle)

	return Vector[float64]{math.Hypot(w*cos, h*sin), math.Hypot(w*sin, h*cos)}
}

// minMax returns the minimum and maximum corner of the ellipse, the corners of Bounds: the
// center less and plus its extent, each rounded once for an integer T.
func (e Ellipse[T]) minMax() (Point[T], Point[T]) {
	center, extent := e.Center.Float(), e.extent()

	return center.Add(extent.Negate()).Cast[T](), center.Add(extent).Cast[T]()
}

// Translate creates a new Ellipse translated by the given vector.
func (e Ellipse[T]) Translate(vector Vector[T]) Ellipse[T] {
	return Ellipse[T]{e.Center.Add(vector), e.Size, e.Angle}
}

// MoveTo creates a new Ellipse with the same size and angle centered at point.
func (e Ellipse[T]) MoveTo(point Point[T]) Ellipse[T] {
	return Ellipse[T]{point, e.Size, e.Angle}
}

// Scale creates a new Ellipse with both semi-axes scaled by the given factor. A negative
// factor scales by its absolute value, since an ellipse mirrored about its center is the same
// ellipse.
func (e Ellipse[T]) Scale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Scale(factor).Abs(), e.Angle}
}

// ScaleXY creates a new Ellipse with the semi-axes scaled by the given factors along its own
// axes, negative ones by their absolute value like Scale.
func (e Ellipse[T]) ScaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.ScaleXY(factorX, factorY).Abs(), e.Angle}
}

// Unscale creates a new Ellipse with both semi-axes scaled by the inverse factor, the inverse
// of Scale and negative factors taken absolute like it. Like Divide it panics for a zero factor.
func (e Ellipse[T]) Unscale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Unscale(factor).Abs(), e.Angle}
}

// UnscaleXY creates a new Ellipse with the semi-axes scaled by the inverse of the given
// factors, the inverse of ScaleXY. Like Divide it panics for a zero factor.
func (e Ellipse[T]) UnscaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.UnscaleXY(factorX, factorY).Abs(), e.Angle}
}

// Resize creates a new Ellipse with the given semi-axes, taken absolute like Ell.
func (e Ellipse[T]) Resize(size Size[T]) Ellipse[T] {
	return Ellipse[T]{e.Center, size.Abs(), e.Angle}
}

// Canonical creates a new Ellipse in the form Ell and Rotate build, with the size taken
// absolute and the angle normalized to [0, 2π): a well-formed ellipse is returned as it is, up
// to the full turns Equal already ignores. It repairs a negative semi-axis written as a struct
// literal or decoded from JSON, which would place the boundary half a turn away, and brings a
// decoded angle onto the seam Rotate keeps. An angle within Delta of zero or of a full turn,
// the residue a chain of Rotate and Lerp calls can leave, becomes exactly zero, as
// Rectangle.Canonical makes it; Rotate itself never snaps.
func (e Ellipse[T]) Canonical() Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Abs(), snapAngle(e.Angle)}
}

// Grow creates a new Ellipse with both semi-axes increased by amount, clamped to zero. The
// amount is added to each radius, as Circle.Grow adds it to the one radius, so each extent of
// Bounds grows by twice it.
func (e Ellipse[T]) Grow(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Grow(amount).AtLeast(Size[T]{}), e.Angle}
}

// GrowXY creates a new Ellipse with the semi-axes increased by the given amounts along its own
// axes, clamped to zero. Each amount is added to that radius, like Grow.
func (e Ellipse[T]) GrowXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.GrowXY(amountX, amountY).AtLeast(Size[T]{}), e.Angle}
}

// Shrink creates a new Ellipse with both semi-axes decreased by amount, clamped to zero.
func (e Ellipse[T]) Shrink(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.Shrink(amount).AtLeast(Size[T]{}), e.Angle}
}

// ShrinkXY creates a new Ellipse with the semi-axes decreased by the given amounts along its
// own axes, clamped to zero.
func (e Ellipse[T]) ShrinkXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size.ShrinkXY(amountX, amountY).AtLeast(Size[T]{}), e.Angle}
}

// Lerp creates a new Ellipse in linear interpolation towards the given ellipse, moving the
// center and the semi-axes together and turning the angle along the shorter arc with
// LerpAngle, normalized to [0, 2π) like Rotate. It extrapolates outside [0, 1] like
// Point.Lerp, with the size taken absolute, so an extrapolation past a zero semi-axis grows
// the ellipse again.
func (e Ellipse[T]) Lerp(ellipse Ellipse[T], t float64) Ellipse[T] {
	return Ellipse[T]{e.Center.Lerp(ellipse.Center, t), e.Size.Lerp(ellipse.Size, t).Abs(), NormalizeAngle(LerpAngle(e.Angle, ellipse.Angle, t))}
}

// Transform creates a new Ellipse by applying the given matrix: the center moves, the
// semi-axes scale by the factors the matrix applies along its axes, and the angle turns by the
// angle of the matrix or is mirrored about it for a reflection, exactly as Rectangle.Transform
// and RegularPolygon.Transform place theirs. A move, a turn, a reflection and a uniform scale
// are exact, and so is a scale of the axes while the ellipse is not turned. A shear, or a scale
// of the axes of a turned ellipse, maps it onto an ellipse of another angle, which this
// method gives as the nearest ellipse of the angle the matrix names rather than the exact one,
// on the factors Matrix.Scaling reports. For integer T the center and the semi-axes are each
// rounded once.
func (e Ellipse[T]) Transform[M Float](matrix Matrix[M]) Ellipse[T] {
	m := matrix.Float()
	scaling := m.Scaling()

	return Ellipse[T]{e.Center.Transform(matrix), e.Size.ScaleXY(scaling.X, scaling.Y).Abs(), m.turnedAngle(e.Angle)}
}

// Rotate creates a new Ellipse turned by the given angle (in radians) about its center, in the
// same sense as Vector.Rotate and as Rectangle.Rotate turns a box. The stored angle is
// normalized to [0, 2π) to prevent drift from repeated rotations. An ellipse of equal
// semi-axes is a circle and no angle moves it, though the angle is stored all the same.
func (e Ellipse[T]) Rotate(angle float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Size, NormalizeAngle(e.Angle + angle)}
}

// AlignTo creates a new Ellipse moved so that its Anchor in the given direction lands on the
// point, the inverse of Anchor: DirectionNone aligns the center, like MoveTo.
func (e Ellipse[T]) AlignTo(direction Direction, point Point[T]) Ellipse[T] {
	return e.Translate(point.Subtract(e.Anchor(direction)))
}

// Contains reports whether the given point lies within the ellipse, boundary included within
// Epsilon of T, the same closed convention as Circle.Contains: it holds exactly where
// DistanceTo is zero. A degenerate ellipse has no interior and contains the points of the
// segment it is.
func (e Ellipse[T]) Contains(point Point[T]) bool {
	return e.DistanceSquaredTo(point) == 0
}

// DistanceTo returns the distance from the given point to the nearest point of the ellipse:
// zero exactly where Contains holds, so a point within Epsilon of T of the boundary is at
// distance zero rather than at the rounding error that put it there, and otherwise the
// distance to the nearest point of the boundary.
func (e Ellipse[T]) DistanceTo(point Point[T]) float64 {
	return math.Sqrt(e.DistanceSquaredTo(point))
}

// DistanceSquaredTo returns the squared distance DistanceTo takes the root of, faster for
// comparisons: zero for a point within the boundary, or within Epsilon of T of it, and the
// squared distance to the foot of the perpendicular otherwise. It is a float64 even for an
// integer T, since that foot is not a lattice point in general, like every shape with an edge.
// Contains is built on it.
//
// The interior is decided on the quadratic form, which is exact, and only a point outside it
// pays for the foot, which no closed form gives and nearestOffset finds by bisection.
func (e Ellipse[T]) DistanceSquaredTo(point Point[T]) float64 {
	local := e.localOffset(point)
	if e.form(local) <= 1 {
		return 0
	}

	distanceSquared := local.Subtract(e.nearestOffset(local)).LengthSquared()
	if lessOrEqualSquared[T](distanceSquared, 0) {
		return 0
	}

	return distanceSquared
}

// form returns the quadratic form of the ellipse at the given offset in its local frame,
// (x/width)² + (y/height)²: below one inside the boundary, exactly one on it and above it
// outside, the one exact test of the interior. A zero semi-axis divides by zero and gives an
// infinity, or a NaN on the axis itself, and neither is at most one, so a degenerate ellipse
// has no interior and every point is left to the distance to its boundary.
func (e Ellipse[T]) form(local Vector[float64]) float64 {
	w, h := e.Size.Float().XY()
	x, y := local.X/w, local.Y/h

	return x*x + y*y
}

// nearestOffset returns the offset of the point of the boundary nearest to the given offset,
// both in the local frame of the ellipse: the foot of the perpendicular from it, which is the
// nearest point inside the boundary as well as outside. A degenerate ellipse is the segment
// its other axis spans, and the foot is the offset clamped to it on each axis, the zero
// semi-axis clamping its own to zero.
func (e Ellipse[T]) nearestOffset(local Vector[float64]) Vector[float64] {
	w, h := e.Size.Float().XY()
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

// foot returns the foot of the perpendicular from (x, y) to the ellipse of semi-axes a and b,
// in the quarter the symmetry of the ellipse leaves: both semi-axes positive with the longer
// first, and both coordinates non-negative. Every other point is one of the four reflections
// of a point of that quarter, which nearestOffset takes it back to.
//
// A point on an axis is answered directly, since the perpendicular from it either runs along
// that axis or meets the boundary where the evolute, the curve of the centers of curvature,
// still reaches: within the evolute the foot leaves the axis, and the two cases meet where the
// evolute crosses it. Anywhere else the foot is placed from the parameter root bisects for,
// which no closed form gives.
func (e Ellipse[T]) foot(a, b, x, y float64) (float64, float64) {
	if y == 0 {
		if evolute, center := a*x, a*a-b*b; evolute < center {
			ratio := evolute / center

			return a * ratio, b * math.Sqrt(1-ratio*ratio)
		}

		return a, 0
	}

	if x == 0 {
		return 0, b
	}

	aspect := (a / b) * (a / b)
	s := e.root(aspect, x/a, y/b)

	return aspect * x / (s + aspect), y / (s + 1)
}

// root returns the parameter of the pencil of ellipses at which the boundary meets the
// perpendicular from the point (x, y), which foot reads the foot off: the point is given in
// units of the semi-axes, off both axes, and aspect is the squared ratio of the longer
// semi-axis to the shorter. The parameter is the sole root of a function falling through zero
// over the bracket the point itself gives, negative from inside the boundary and positive from
// outside.
//
// It is found by bisection, which halves that bracket until the midpoint is one of its ends
// and no float lies between them, so it ends in the precision of a float64 and no iteration
// count has to be chosen; a midpoint that lands exactly on the root is kept as the far end and
// the halving runs down to it. It reads no field of the ellipse, only the point in the frame
// of its semi-axes, so the receiver is unnamed.
func (Ellipse[T]) root(aspect, x, y float64) float64 {
	gradient := aspect * x

	s0, s1 := y-1, math.Hypot(gradient, y)-1
	if x*x+y*y < 1 {
		s1 = 0
	}

	for {
		s := (s0 + s1) / 2
		if s == s0 || s == s1 {
			return s
		}

		if px, py := gradient/(s+aspect), y/(s+1); px*px+py*py > 1 {
			s0 = s
		} else {
			s1 = s
		}
	}
}

// Equal checks for equal center, size and angle values using tolerant numeric comparison.
// Angles are compared with EqualAngle, so a full turn or the sign of an angle does not matter.
// The angle is compared even for two circles, which no angle moves: they are different values,
// not the same ellipse, as RegularPolygon.Equal compares the angle of two empty polygons.
func (e Ellipse[T]) Equal(ellipse Ellipse[T]) bool {
	return e.Center.Equal(ellipse.Center) && e.Size.Equal(ellipse.Size) && EqualAngle(e.Angle, ellipse.Angle)
}

// IsZero checks if center point, size and angle are zero, comparing the angle like Equal so
// that a full turn counts as zero.
func (e Ellipse[T]) IsZero() bool {
	return e.Equal(Ellipse[T]{})
}

// IsAligned reports whether the ellipse is axis-aligned: its Angle is exactly zero, as Ell
// with no angle and Rotate by a full turn leave it, the same exact test Rectangle.IsAligned
// makes. No tolerance is applied; Canonical snaps a residue to zero.
func (e Ellipse[T]) IsAligned() bool {
	return e.Angle == 0
}

// IsCircle reports whether the two semi-axes are equal, within Epsilon of T, so that the
// ellipse is a circle and no angle moves it, which is where Circle gives the same shape back
// rather than the circle around it.
func (e Ellipse[T]) IsCircle() bool {
	return Equal(e.Size.Width, e.Size.Height)
}

// Circle converts the ellipse into the circle around it, the one of the major semi-axis: the
// smallest circle containing the boundary, which touches it at the two ends of that axis
// whatever the angle, as Bounds is the smallest box around it. An ellipse of equal semi-axes
// is that circle, so the conversion is exact exactly where IsCircle holds and is the inverse
// of Circle.Ellipse there; anywhere else it is the round hull, not the same shape.
func (e Ellipse[T]) Circle() Circle[T] {
	return Circle[T]{e.Center, e.SemiMajor()}
}

// RegularPolygon converts the ellipse into the RegularPolygon of n vertices inscribed in it,
// on the same center, semi-axes and angle: every vertex lies on the boundary, the first at the
// end of the width semi-axis, and the polygon is the outline a circle and an ellipse do not
// have, so its Vertices and Edges are what draws or walks one. The resolution is the caller's:
// the polygon is the ellipse only in the limit, and it is always inside it.
//
// It takes no Orientation, where Circle.RegularPolygon does: an orientation is the phase of the
// first vertex around the ring, and only a circle can be turned to place it, since turning the
// ring and stepping around it are the same thing there. A polygon with fewer than one vertex
// is empty, as RegPol builds it.
//
// It is the polygon within the ellipse, the only one that converts back exactly. For the
// polygon about it, whose edges touch the boundary at their midpoints, scale the ellipse by
// the ratio of the circumradius of a regular polygon to its apothem first:
// e.Scale(1 / math.Cos(Pi / float64(n))).RegularPolygon(n).
func (e Ellipse[T]) RegularPolygon(n int) RegularPolygon[T] {
	return RegularPolygon[T]{e.Center, e.Size, n, e.Angle}
}

// Cast converts the ellipse to an Ellipse of another number type, rounding as Cast does and
// keeping the angle.
func (e Ellipse[T]) Cast[R Number]() Ellipse[R] {
	return Ellipse[R]{e.Center.Cast[R](), e.Size.Cast[R](), e.Angle}
}

// Int converts the ellipse to an Ellipse[int], rounding the center and the semi-axes on their
// own like every other Int, and keeping the angle.
func (e Ellipse[T]) Int() Ellipse[int] {
	return Ellipse[int]{e.Center.Int(), e.Size.Int(), e.Angle}
}

// Float converts the ellipse to an Ellipse[float64].
func (e Ellipse[T]) Float() Ellipse[float64] {
	return Ellipse[float64]{e.Center.Float(), e.Size.Float(), e.Angle}
}

// String returns the ellipse in the form of its constructor: Ell((x,y);WxH), center then
// semi-axes, with the angle appended as Ell((x,y);WxH;a) for a turned ellipse, as the JSON
// carries it only then.
func (e Ellipse[T]) String() string {
	if e.IsAligned() {
		return fmt.Sprintf("Ell(%s;%s)", e.Center.String(), e.Size.String())
	}

	return fmt.Sprintf("Ell(%s;%s;%s)", e.Center.String(), e.Size.String(), String(e.Angle))
}

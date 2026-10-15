package geom

import (
	"math"
	"slices"
)

// Ellipse is an ellipse of a center and semi-axes, turned by Angle radians about the center.
//
// Radius holds the semi-axes, named in the frame before the turn: Radius.Width is the semi-axis
// that pointed along X. Equal compares the angle even where equal semi-axes make it a
// circle that no angle moves.
//
// The radius is never negative: constructors and changes take it absolute or clamp it at zero. A
// negative semi-axis written as a struct literal or decoded from JSON gives no meaningful answers
// until Canonical repairs it. A zero semi-axis is not repaired: the ellipse is then the segment the
// other axis spans, and answers as that segment.
//
// For integer T a turned ellipse is measured and compared at no less than the tolerance of float64,
// since its frame places no lattice points.
//
// The ellipse has no vertices, and no Intersects method for a shape but Segment. Intersects,
// Distance, Collide and Sweep test it against any other by its Support, or by its distance to the
// center of a Circle; Segment, Ray and Line cross, clip and hit it in closed form, and
// RegularPolygon and Polygon give its outline.
type Ellipse[T Number] struct {
	Center Point[T] `json:",embed"`
	Radius Size[T]  `json:",embed"`
	Angle  float64  `json:"a,omitzero"`
}

// Ell is shorthand for Ellipse{center, radius, angle}, with the radius taken absolute.
func Ell[T Number](center Point[T], radius Size[T], angle float64) Ellipse[T] {
	return Ellipse[T]{center, radius.Abs(), angle}
}

// SemiMajor returns the longer semi-axis, the one the foci lie on.
func (e Ellipse[T]) SemiMajor() T {
	return max(e.Radius.Width, e.Radius.Height)
}

// SemiMinor returns the shorter semi-axis.
func (e Ellipse[T]) SemiMinor() T {
	return min(e.Radius.Width, e.Radius.Height)
}

// Foci returns the two focal points on the major axis, the one on its negative end before the turn
// first. A circle has both at the center.
func (e Ellipse[T]) Foci() (Point[T], Point[T]) {
	w, h := e.Radius.Float().Abs().XY()
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

	return e.worldPoint(VectorFromAngleSize(e.parametricAngle(direction.Angle()), e.Radius.Float()))
}

// Centroid returns the center.
func (e Ellipse[T]) Centroid() Point[T] {
	return e.Center
}

// Area returns the area (π·w·h).
func (e Ellipse[T]) Area() float64 {
	return Pi * float64(e.Radius.Width) * float64(e.Radius.Height)
}

// Perimeter returns the length of the boundary by Ramanujan's second approximation, exact for a
// circle and within a part in a thousand for a degenerate ellipse.
func (e Ellipse[T]) Perimeter() float64 {
	w, h := e.Radius.Float().XY()

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
	w, h := e.Radius.Float().XY()

	return Pi * w * h * (float64(w*w) + float64(h*h)) / 4
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

// Bounds returns the axis-aligned bounding box. For integer T it is rounded outward, except a side
// within the tolerance of a whole coordinate, which takes it.
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
	w, h := e.Radius.Float().XY()
	if w == 0 || h == 0 {
		return direction
	}

	sin, cos := math.Sincos(direction)

	return math.Atan2(w*sin, h*cos)
}

// extent returns the half extent of the turned ellipse along each world axis, in closed form, since
// the boundary touches the box where its tangent is axis-aligned.
func (e Ellipse[T]) extent() Vector[float64] {
	w, h := e.Radius.Float().XY()
	sin, cos := math.Sincos(e.Angle)

	return Vector[float64]{math.Hypot(w*cos, h*sin), math.Hypot(w*sin, h*cos)}
}

// minMax returns the corners of Bounds, the center less and plus its extent, for integer T rounded
// outward past the tolerance.
func (e Ellipse[T]) minMax() (Point[T], Point[T]) {
	center, extent := e.Center.Float(), e.extent()
	a, b := center.Add(extent.Negate()), center.Add(extent)

	epsilon := 0.0
	if isInt[T]() {
		epsilon = e.epsilon(0)
	}

	return Point[T]{-castUp[T](-a.X - epsilon), -castUp[T](-a.Y - epsilon)}, Point[T]{castUp[T](b.X - epsilon), castUp[T](b.Y - epsilon)}
}

// Translate returns the ellipse moved by the vector.
func (e Ellipse[T]) Translate(vector Vector[T]) Ellipse[T] {
	return Ellipse[T]{e.Center.Add(vector), e.Radius, e.Angle}
}

// MoveTo returns the ellipse with its center at the point.
func (e Ellipse[T]) MoveTo(point Point[T]) Ellipse[T] {
	return Ellipse[T]{point, e.Radius, e.Angle}
}

// AlignTo returns the ellipse moved so that its Anchor in the direction lies at the point.
func (e Ellipse[T]) AlignTo(direction Direction, point Point[T]) Ellipse[T] {
	return e.Translate(point.Subtract(e.Anchor(direction)))
}

// Rotate returns the ellipse turned by the angle in radians about its center, the stored angle
// normalized to [0, 2π).
func (e Ellipse[T]) Rotate(angle float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius, NormalizeAngle(e.Angle + angle)}
}

// Scale returns the ellipse with both semi-axes scaled by the factor, taken absolute.
func (e Ellipse[T]) Scale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.Scale(factor).Abs(), e.Angle}
}

// ScaleXY returns the ellipse with its semi-axes scaled by the factors along its own axes, taken
// absolute.
func (e Ellipse[T]) ScaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.ScaleXY(factorX, factorY).Abs(), e.Angle}
}

// Unscale returns the ellipse with both semi-axes divided by the factor, taken absolute. It panics
// for a zero factor.
func (e Ellipse[T]) Unscale(factor float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.Unscale(factor).Abs(), e.Angle}
}

// UnscaleXY returns the ellipse with its semi-axes divided by the factors along its own axes, taken
// absolute. It panics for a zero factor.
func (e Ellipse[T]) UnscaleXY(factorX, factorY float64) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.UnscaleXY(factorX, factorY).Abs(), e.Angle}
}

// Resize returns the ellipse with the semi-axes, taken absolute.
func (e Ellipse[T]) Resize(radius Size[T]) Ellipse[T] {
	return Ellipse[T]{e.Center, radius.Abs(), e.Angle}
}

// Transform returns the ellipse mapped by the matrix, exactly, since an affine map takes an ellipse
// to an ellipse. A reflection mirrors the angle, and where the semi-axes come out swapped the angle
// takes a quarter turn.
func (e Ellipse[T]) Transform[M Float](matrix Matrix[M]) Ellipse[T] {
	after, stretch, _, _ := matrix.mapping(e.Angle, e.Radius.Float()).decomposition()

	return Ellipse[T]{e.Center.Transform(matrix), stretch.Size().Cast[T](), NormalizeAngle(after)}
}

// Grow returns the ellipse with both semi-axes increased by the amount, clamped at zero.
func (e Ellipse[T]) Grow(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.Grow(amount).AtLeastZero(), e.Angle}
}

// GrowXY returns the ellipse with its semi-axes increased by the amounts along its own axes,
// clamped at zero.
func (e Ellipse[T]) GrowXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.GrowXY(amountX, amountY).AtLeastZero(), e.Angle}
}

// Shrink returns the ellipse with both semi-axes decreased by the amount, clamped at zero.
func (e Ellipse[T]) Shrink(amount T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.Shrink(amount).AtLeastZero(), e.Angle}
}

// ShrinkXY returns the ellipse with its semi-axes decreased by the amounts along its own axes,
// clamped at zero.
func (e Ellipse[T]) ShrinkXY(amountX, amountY T) Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.ShrinkXY(amountX, amountY).AtLeastZero(), e.Angle}
}

// Lerp returns the ellipse interpolated toward the given one, center and semi-axes together and the
// angle along the shorter arc by LerpAngle, normalized to [0, 2π), and extrapolated outside [0, 1].
// The radius is taken absolute, so extrapolating past a zero semi-axis grows the ellipse again.
func (e Ellipse[T]) Lerp(ellipse Ellipse[T], t float64) Ellipse[T] {
	return Ellipse[T]{e.Center.Lerp(ellipse.Center, t), e.Radius.Lerp(ellipse.Radius, t).Abs(), NormalizeAngle(LerpAngle(e.Angle, ellipse.Angle, t))}
}

// Canonical returns the ellipse with its radius taken absolute and its angle normalized to [0, 2π),
// an angle within Delta of zero or a full turn becoming exactly zero. Rotate and Lerp never snap.
func (e Ellipse[T]) Canonical() Ellipse[T] {
	return Ellipse[T]{e.Center, e.Radius.Abs(), snapAngle(e.Angle)}
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
	w, h := e.Radius.Float().XY()
	local := direction.Rotate(-e.Angle)
	unit := Vector[float64]{w * local.X, h * local.Y}.Normalize()

	return e.Center.Float().Add(Vector[float64]{w * unit.X, h * unit.Y}.Rotate(e.Angle))
}

// Contains reports whether the point lies in the ellipse. A degenerate ellipse contains the points
// of its segment.
func (e Ellipse[T]) Contains(point Point[T]) bool {
	return e.DistanceSquaredTo(point) == 0
}

// IntersectsSegment reports whether the ellipse and the segment share a point.
func (e Ellipse[T]) IntersectsSegment(segment Segment[T]) bool {
	return segment.IntersectsEllipse(e)
}

// IntersectsRay reports whether the ellipse and the ray share a point.
func (e Ellipse[T]) IntersectsRay(ray Ray[T]) bool {
	return ray.IntersectsEllipse(e)
}

// IntersectsLine reports whether the ellipse and the line share a point.
func (e Ellipse[T]) IntersectsLine(line Line[T]) bool {
	return line.IntersectsEllipse(e)
}

// IntersectionSegment returns the points where the segment crosses the boundary, from Start to End.
func (e Ellipse[T]) IntersectionSegment(segment Segment[T]) Pair[Point[T]] {
	return segment.IntersectionEllipse(e)
}

// IntersectionRay returns the points where the ray crosses the boundary, from Origin on.
func (e Ellipse[T]) IntersectionRay(ray Ray[T]) Pair[Point[T]] {
	return ray.IntersectionEllipse(e)
}

// IntersectionLine returns the points where the line crosses the boundary, in order along its
// Vector.
func (e Ellipse[T]) IntersectionLine(line Line[T]) Pair[Point[T]] {
	return line.IntersectionEllipse(e)
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
	if lessOrEqualSquared(distanceSquared, 0, e.epsilon(point.magnitude())) {
		return local, 0
	}

	return foot, distanceSquared
}

// containsOffset reports whether the local offset lies inside, on the quadratic form multiplied
// through by the squared semi-axes so no division rounds it. A degenerate ellipse has no inside and
// leaves every point, NaN included, to the distance.
func (e Ellipse[T]) containsOffset(local Vector[float64]) bool {
	w, h := e.Radius.Float().XY()
	if w == 0 || h == 0 {
		return false
	}

	x, y, extent := local.X*h, local.Y*w, w*h

	return float64(x*x)+float64(y*y) <= extent*extent
}

// nearestOffset returns the foot of the perpendicular from a local offset outside the ellipse. For
// a degenerate ellipse it is the offset clamped to the segment on each axis.
func (e Ellipse[T]) nearestOffset(local Vector[float64]) Vector[float64] {
	w, h := e.Radius.Float().Abs().XY()
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

// fractions returns the fractions along the segment, ascending and not yet clamped to it, at which
// it crosses the boundary, and how many there are. An endpoint on the boundary stands in for the
// crossing it is, so Segment.IntersectionEllipse and Segment.HitEllipse place points from the same
// values.
func (e Ellipse[T]) fractions(segment Segment[T]) ([2]float64, int) {
	entry, exit, ok := e.chord(segment)
	magnitude := max(e.Center.magnitude(), segment.magnitude())
	start, end := e.touches(segment.Start, magnitude), e.touches(segment.End, magnitude)

	switch {
	case start && end && segment.Vector().hasDirection():
		return [2]float64{0, 1}, 2
	case ok && entry < exit:
		if start {
			entry, exit = segment.snapToEndpoint(entry, exit, 0)
		}
		if end {
			entry, exit = segment.snapToEndpoint(entry, exit, 1)
		}

		return [2]float64{entry, exit}, 2
	case start:
		return [2]float64{0}, 1
	case end:
		return [2]float64{1}, 1
	case ok:
		return [2]float64{entry}, 1
	default:
		return [2]float64{}, 0
	}
}

// chord returns the fractions along the line of the segment where it enters and leaves the ellipse,
// not clamped to the segment, and false where it misses or the segment has no direction. The line
// misses where its distance from the center exceeds the half width of the ellipse across it, the
// reach of Support along its normal, so the miss compares two lengths as it does for a circle. A
// chord whose ends would compare Equal is a tangent and gives one fraction twice: the middle of the
// chord, or for a line passing outside within the tolerance, its point nearest the boundary. A line
// along a zero semi-axis runs along the segment of the other.
func (e Ellipse[T]) chord(segment Segment[T]) (float64, float64, bool) {
	cast := segment.Float()
	direction := cast.Vector().Rotate(-e.Angle)
	if !direction.hasDirection() {
		return 0, 0, false
	}

	w, h := e.Radius.Float().XY()
	offset := e.localOffset(segment.Start)
	epsilon := e.epsilon(segment.magnitude())

	lengthSquared := direction.LengthSquared()
	length := math.Sqrt(lengthSquared)
	gap := offset.Cross(direction) / length
	gapSquared := float64(gap * gap)

	x, y := float64(direction.X*h), float64(direction.Y*w)
	across := float64(x*x) + float64(y*y)
	widthSquared := across / lengthSquared
	if !lessOrEqualSquared(gapSquared, math.Sqrt(across)/length, epsilon) {
		return 0, 0, false
	}

	along := -offset.Dot(direction) / lengthSquared
	if across == 0 {
		reach := math.Hypot(float64(direction.X*w), float64(direction.Y*h)) / lengthSquared
		if lessOrEqualSquared(4*float64(reach*reach)*lengthSquared, 0, epsilon) {
			return along, along, true
		}

		return along - reach, along + reach, true
	}

	ox, oy := float64(offset.X*h), float64(offset.Y*w)
	middle := -(float64(ox*x) + float64(oy*y)) / across

	area := float64(w * h)
	halfSquared := float64(area*area) * max(widthSquared-gapSquared, 0) / float64(widthSquared*widthSquared)
	if lessOrEqualSquared(4*halfSquared, 0, epsilon) {
		if gapSquared <= widthSquared {
			return middle, middle, true
		}

		support := e.Support(cast.PointAt(along).Subtract(e.Center.Float()))
		touch := support.Subtract(cast.Start).Dot(cast.Vector()) / lengthSquared

		return touch, touch, true
	}

	half := math.Sqrt(halfSquared) / length

	return middle - half, middle + half, true
}

// touches reports whether the point lies on the boundary within the tolerance, from either side: it
// lies in the ellipse, and its level on the quadratic form over the gradient there, the distance to
// the boundary to first order, is within the tolerance of zero. Every point of a degenerate ellipse
// is on its boundary.
func (e Ellipse[T]) touches(point Point[T], magnitude float64) bool {
	if !e.Contains(point) {
		return false
	}

	w, h := e.Radius.Float().XY()
	if w == 0 || h == 0 {
		return true
	}

	local := e.localOffset(point)
	x, y, extent := float64(local.X*h), float64(local.Y*w), float64(w*h)
	level := float64(x*x) + float64(y*y) - float64(extent*extent)
	gap := level / (2 * math.Hypot(float64(x*h), float64(y*w)))

	return lessOrEqualSquared(float64(gap*gap), 0, e.epsilon(magnitude))
}

// epsilon returns the tolerance of a comparison between the ellipse and geometry of the given
// magnitude, for a turned ellipse never less than that of float64, since its frame places no
// lattice points.
func (e Ellipse[T]) epsilon(magnitude float64) float64 {
	magnitude = max(e.Center.magnitude(), magnitude)
	if e.IsAligned() {
		return epsilonAt[T](magnitude)
	}

	return max(epsilonAt[T](magnitude), epsilonAt[float64](magnitude))
}

// normal returns the perpendicular of the boundary at the point, the gradient of the quadratic form,
// of any length.
func (e Ellipse[T]) normal(point Point[float64]) Vector[float64] {
	w, h := e.Radius.Float().XY()
	local := point.Subtract(e.Center.Float()).Rotate(-e.Angle)

	return Vector[float64]{float64(local.X * float64(h*h)), float64(local.Y * float64(w*w))}.Rotate(e.Angle)
}

// Equal reports whether the ellipses have equal center and radius and an angle equal by EqualAngle.
func (e Ellipse[T]) Equal(ellipse Ellipse[T]) bool {
	return e.Center.Equal(ellipse.Center) && e.Radius.Equal(ellipse.Radius) && EqualAngle(e.Angle, ellipse.Angle)
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
	return Equal(e.Radius.Width, e.Radius.Height)
}

// Circle returns the circle of the major semi-axis about the center, the smallest enclosing the
// ellipse. It is the same shape only where IsCircle holds, and there the inverse of Circle.Ellipse.
func (e Ellipse[T]) Circle() Circle[T] {
	return Circle[T]{e.Center, e.SemiMajor()}
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
	size := e.Radius.Float()
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
	return RegularPolygon[T]{e.Center, e.Radius, n, e.Angle, RegularPolygonOrientationPhase(n, orientation)}
}

// Cast converts the ellipse to an Ellipse of another number type, rounding as Cast does.
func (e Ellipse[T]) Cast[R Number]() Ellipse[R] {
	return Ellipse[R]{e.Center.Cast[R](), e.Radius.Cast[R](), e.Angle}
}

// Int converts the ellipse to an Ellipse[int].
func (e Ellipse[T]) Int() Ellipse[int] {
	return Ellipse[int]{e.Center.Int(), e.Radius.Int(), e.Angle}
}

// Float converts the ellipse to an Ellipse[float64].
func (e Ellipse[T]) Float() Ellipse[float64] {
	return Ellipse[float64]{e.Center.Float(), e.Radius.Float(), e.Angle}
}

// String returns the ellipse in the form of its constructor: Ell((x,y);WxH), with ;a before the
// closing parenthesis for a turned ellipse.
func (e Ellipse[T]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Ell("...)
	buf = e.Center.appendString(buf)
	buf = append(buf, ';')
	buf = e.Radius.appendString(buf)
	if !e.IsAligned() {
		buf = append(buf, ';')
		buf = appendString(buf, e.Angle)
	}
	buf = append(buf, ')')

	return string(buf)
}

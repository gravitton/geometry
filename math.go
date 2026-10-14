package geom

import (
	"math"
)

const (
	// Pi is the ratio of the circumference of a circle to its diameter.
	Pi = math.Pi
	// RadToDeg converts radians to degrees (180/π).
	RadToDeg float64 = 180.0 / math.Pi
	// DegToRad converts degrees to radians (π/180).
	DegToRad float64 = math.Pi / 180.0

	// Delta is the equality tolerance of float64.
	Delta float64 = 1e-6
	// Delta32 is the equality tolerance of float32.
	Delta32 float64 = 1e-4

	// Sqrt2 is the square root of 2.
	Sqrt2 = math.Sqrt2
	// Sqrt3 is the square root of 3.
	Sqrt3 = 1.732050807568877293527446341505872367
	// OneOverSqrt2 is 1/√2, each component of a unit diagonal vector.
	OneOverSqrt2 = 1 / math.Sqrt2
)

// Multiply returns x scaled by the factor.
func Multiply[T Number](x T, factor float64) T {
	return Cast[T](float64(x) * factor)
}

// Divide returns x divided by the factor. It panics for a zero factor.
func Divide[T Number](x T, factor float64) T {
	if factor == 0 {
		panic("geom: division by zero")
	}

	return Cast[T](float64(x) / factor)
}

// ratio returns a divided by b in float64, or 0 for a zero b, where an Inf or NaN would poison
// later arithmetic.
func ratio[T Number](a, b T) float64 {
	if b == 0 {
		return 0
	}

	return float64(a) / float64(b)
}

// Abs returns the absolute value, computed in T, so integers stay exact; the most negative integer
// of T has no positive counterpart and is returned as it is.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}

	return x
}

// Sign returns 1 for a positive x, -1 for a negative one and 0 for zero.
func Sign[T Number](x T) T {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// Round returns x rounded to the nearest integer, half away from zero; an integer T is unchanged.
func Round[T Number](x T) T {
	if isInt[T]() {
		return x
	}

	return T(math.Round(float64(x)))
}

// RoundHalfUp returns x rounded to the nearest integer, half toward positive infinity, so rounding
// a value shifted by an integer is the rounded value shifted by it, the rounding of a lattice. An
// integer T is unchanged.
func RoundHalfUp[T Number](x T) T {
	if isInt[T]() {
		return x
	}

	value := float64(x)
	rounded := math.Floor(value)
	if value-rounded >= 0.5 {
		rounded++
	}

	return T(rounded)
}

// Floor returns the largest integer at most x; an integer T is unchanged.
func Floor[T Number](x T) T {
	if isInt[T]() {
		return x
	}

	return T(math.Floor(float64(x)))
}

// Ceil returns the smallest integer at least x; an integer T is unchanged.
func Ceil[T Number](x T) T {
	if isInt[T]() {
		return x
	}

	return T(math.Ceil(float64(x)))
}

// Mod returns n wrapped into [0, m), or (m, 0] for a negative m, as a floored division gives it. It
// panics for a zero m.
func Mod[T Integer](n, m T) T {
	remainder := n % m
	if remainder != 0 && (remainder < 0) != (m < 0) {
		remainder += m
	}

	return remainder
}

// IsFinite reports whether x is neither NaN nor infinite, the values Cast and Int take into an
// integer T. An integer always is.
func IsFinite[T Number](x T) bool {
	return math.Abs(float64(x)) <= math.MaxFloat64
}

// Clamp returns the value limited to [min, max].
func Clamp[T Number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

// Sum returns the sum of the values, accumulated in float64 and rounded once into T, and 0 for
// none.
func Sum[T Number](values []T) T {
	var total float64
	for _, value := range values {
		total += float64(value)
	}

	return Cast[T](total)
}

// Lerp returns the linear interpolation between a and b at t, extrapolated outside [0, 1].
func Lerp[T Number](a, b T, t float64) T {
	return Cast[T](float64(a) + float64((float64(b)-float64(a))*t))
}

// InverseLerp returns the t at which Lerp of a and b gives the value, unclamped, and 0 for an empty
// range where a equals b.
func InverseLerp[T Number](a, b, value T) float64 {
	from := float64(a)

	return ratio(float64(value)-from, float64(b)-from)
}

// Remap returns the value moved from the range a1 to b1 to the same ratio of the range a2 to b2,
// unclamped, and a2 for an empty first range. Either range may run downward.
func Remap[T Number](value, a1, b1, a2, b2 T) T {
	return Lerp(a2, b2, InverseLerp(a1, b1, value))
}

// Midpoint returns the value halfway between a and b.
func Midpoint[T Number](a, b T) T {
	return Lerp(a, b, 0.5)
}

// MoveToward returns the value moved toward the target by at most step, the target itself exactly
// once it is within reach. A negative step moves away and is not clamped, and a value at the target
// stays there. For integer T a step below half a unit leaves the value where it is.
func MoveToward[T Number](value, target T, step float64) T {
	from := float64(value)
	delta := float64(target) - from

	switch {
	case step >= math.Abs(delta):
		return target
	case delta > 0:
		return Cast[T](from + step)
	case delta < 0:
		return Cast[T](from - step)
	}

	return target
}

// Equal reports whether a and b are equal within Epsilon of T, exactly for an integer T.
//
// The tolerance is absolute, so far from zero it falls below one ulp and becomes exact equality.
// EqualRelative holds at any magnitude.
func Equal[T Number](a, b T) bool {
	if isInt[T]() {
		return a == b
	}

	return EqualDelta(a, b, Epsilon[T]())
}

// EqualDelta reports whether a and b are equal within the delta, or are the same value, so an
// infinity equals itself.
func EqualDelta[T Number](a, b T, delta float64) bool {
	return a == b || math.Abs(float64(a)-float64(b)) <= delta
}

// EqualRelative reports whether a and b are equal within EpsilonRelative, a tolerance that scales
// with their magnitude, exactly for an integer T.
func EqualRelative[T Number](a, b T) bool {
	if isInt[T]() {
		return a == b
	}

	return EqualDelta(a, b, EpsilonRelative(a, b))
}

// LessOrEqual reports whether a is at most b within Epsilon of T, exactly for an integer T.
func LessOrEqual[T Number](a, b T) bool {
	if isInt[T]() {
		return a <= b
	}

	return LessOrEqualDelta(a, b, Epsilon[T]())
}

// LessOrEqualDelta reports whether a is at most b within the delta.
func LessOrEqualDelta[T Number](a, b T, delta float64) bool {
	return float64(a) <= float64(b)+delta
}

// lessOrEqualSquared reports whether a distance, given as its square a2, is at most b within the
// epsilon. It is the one comparison against a radius or a segment: every Contains, DistanceTo and
// Intersects makes it at the epsilonAt of what it compares, so every test rounds alike at the
// boundary and no distance is rooted to be compared. A zero b is the tolerance alone.
func lessOrEqualSquared(a2, b, epsilon float64) bool {
	reach := b + epsilon

	return a2 <= reach*reach
}

// greaterOrEqualSquared reports whether a distance, given as its square a2, is at least b within
// the epsilon.
func greaterOrEqualSquared(a2, b, epsilon float64) bool {
	reach := max(b-epsilon, 0)

	return a2 >= reach*reach
}

// equalSquared reports whether a distance, given as its square a2, equals b within the epsilon, the
// squared form of EqualDelta.
func equalSquared(a2, b, epsilon float64) bool {
	return lessOrEqualSquared(a2, b, epsilon) && greaterOrEqualSquared(a2, b, epsilon)
}

// Epsilon returns the equality tolerance of T: zero for an integer T, Delta32 for float32 and Delta
// for float64, whatever the values compared.
func Epsilon[T Number]() float64 {
	if isInt[T]() {
		return 0
	}

	if isFloat32[T]() {
		return Delta32
	}

	return Delta
}

// EpsilonRelative returns the tolerance of EqualRelative: Epsilon of T scaled by the larger
// magnitude of a and b, and never less than Epsilon of T.
func EpsilonRelative[T Number](a, b T) float64 {
	return Epsilon[T]() * max(1, math.Abs(float64(a)), math.Abs(float64(b)))
}

// epsilonAt returns the boundary tolerance at the magnitude, the largest absolute coordinate a
// comparison reads: Epsilon of T, or two ulps of T at that magnitude where wider, so a point
// rounded into T far from the origin stays on the boundary it was computed on. It is zero for an
// integer T, and NaN, which admits nothing, for the NaN magnitude of an infinite coordinate.
func epsilonAt[T Number](magnitude float64) float64 {
	if isInt[T]() {
		return 0
	}

	if isFloat32[T]() {
		return max(Delta32, magnitude*0x1p-22)
	}

	return max(Delta, magnitude*0x1p-51)
}

// ToRadians converts degrees to radians.
func ToRadians(degrees float64) float64 {
	return degrees * DegToRad
}

// ToDegrees converts radians to degrees.
func ToDegrees(radians float64) float64 {
	return radians * RadToDeg
}

// NormalizeAngle returns the angle wrapped into [0, 2π), and NaN for a non-finite angle.
func NormalizeAngle(angle float64) float64 {
	a := math.Mod(angle, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	if a >= 2*math.Pi {
		a = 0
	}

	return a
}

// snapAngle returns the angle normalized, and zero within Delta of zero or a full turn: the repair
// Canonical applies to an angle, so a shape turned back is aligned again. Rotate and Lerp skip it,
// since a turn below Delta still moves a far corner of a large shape by more than Epsilon.
func snapAngle(angle float64) float64 {
	angle = NormalizeAngle(angle)
	if angle <= Delta || 2*math.Pi-angle <= Delta {
		return 0
	}

	return angle
}

// AngleDistance returns the shortest angular distance between a and b, in [0, π].
func AngleDistance(a, b float64) float64 {
	d := NormalizeAngle(a - b)

	return min(d, 2*math.Pi-d)
}

// LerpAngle returns the angle interpolated from a toward b at t along the shorter arc, extrapolated
// outside [0, 1] and measured from a, not normalized. Angles half a turn apart turn toward
// increasing angle.
func LerpAngle(a, b, t float64) float64 {
	delta := NormalizeAngle(b - a)
	if delta > math.Pi {
		delta -= 2 * math.Pi
	}

	return a + float64(delta*t)
}

// MoveTowardAngle returns the angle turned toward the target by at most step along the shorter arc,
// the target itself, as given, once it is within reach, and otherwise measured from the angle, not
// normalized. Angles half a turn apart turn toward increasing angle. A negative step turns away and
// is not clamped, and an angle at the target stays there.
func MoveTowardAngle(angle, target, step float64) float64 {
	delta := NormalizeAngle(target - angle)
	if delta > math.Pi {
		delta -= 2 * math.Pi
	}

	turn := MoveToward(0, delta, step)
	if turn == delta {
		return target
	}

	return angle + turn
}

// EqualAngle reports whether a and b are the same angle within Delta, whole turns aside, across the
// seam at zero.
func EqualAngle(a, b float64) bool {
	return EqualDelta(AngleDistance(a, b), 0, Delta)
}

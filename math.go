package geom

import (
	"math"
	"strconv"
)

const (
	// Pi is the ratio of the circumference of a circle to its diameter.
	Pi = math.Pi
	// RadToDeg is the multiplier to convert radians to degrees (180/π).
	RadToDeg float64 = 180.0 / math.Pi
	// DegToRad is the multiplier to convert degrees to radians (π/180).
	DegToRad float64 = math.Pi / 180.0

	// Delta is the tolerance used for float64 equality comparisons.
	Delta float64 = 1e-6
	// Delta32 is the tolerance used for float32 equality comparisons.
	Delta32 float64 = 1e-4

	// Sqrt2 is the square root of 2.
	Sqrt2 = math.Sqrt2
	// Sqrt3 is the square root of 3.
	Sqrt3 = 1.732050807568877293527446341505872367
	// OneOverSqrt2 is 1/√2, the length of a unit diagonal vector component.
	OneOverSqrt2 = 1 / math.Sqrt2
)

// NormalizeAngle returns angle normalized to [0, 2π). NaN and ±Inf are returned as NaN.
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

// AngleDistance returns the shortest angular distance between a and b, in [0, π].
func AngleDistance(a, b float64) float64 {
	d := NormalizeAngle(a - b)

	return min(d, 2*math.Pi-d)
}

// EqualAngle reports whether a and b are the same angle within Delta, a full turn or the
// sign of an angle aside. Unlike comparing normalized angles it holds across the 0/2π seam.
func EqualAngle(a, b float64) bool {
	return EqualDelta(AngleDistance(a, b), 0, Delta)
}

// ToRadians converts degrees to radians.
func ToRadians(degrees float64) float64 {
	return degrees * DegToRad
}

// ToDegrees converts radians to degrees.
func ToDegrees(radians float64) float64 {
	return radians * RadToDeg
}

// Multiply multiplies a number by a scale factor.
func Multiply[T Number](x T, factor float64) T {
	return Cast[T](float64(x) * factor)
}

// Divide divides a number by a scale factor. Like the / operator on integers it panics for a
// zero factor, rather than returning an infinity that an integer T could not hold.
func Divide[T Number](x T, factor float64) T {
	if factor == 0 {
		panic("geom: division by zero")
	}

	return Cast[T](float64(x) / factor)
}

// Abs returns the absolute value. It never leaves T, so integers of any width stay exact.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}

	return x
}

// Round returns x rounded to the nearest integer. An integer T is returned unchanged.
func Round[T Number](x T) T {
	if isIntType[T]() {
		return x
	}

	return T(math.Round(float64(x)))
}

// Floor returns the largest integer value less than or equal to x. An integer T is returned unchanged.
func Floor[T Number](x T) T {
	if isIntType[T]() {
		return x
	}

	return T(math.Floor(float64(x)))
}

// Ceil returns the smallest integer value greater than or equal to x. An integer T is returned unchanged.
func Ceil[T Number](x T) T {
	if isIntType[T]() {
		return x
	}

	return T(math.Ceil(float64(x)))
}

// Mod wraps n into [0, m), correctly for negative n. Like the % operator it panics for m == 0.
func Mod[T Integer](n, m T) T {
	return ((n % m) + m) % m
}

// Lerp calculates the linear interpolation between a and b at a ratio t.
// The difference is taken in float64, so it cannot overflow a narrow integer T.
func Lerp[T Number](a, b T, t float64) T {
	return Cast[T](float64(a) + (float64(b)-float64(a))*t)
}

// Midpoint calculates the midpoint between two values. Equivalent to Lerp(a, b, 0.5).
func Midpoint[T Number](a, b T) T {
	return Lerp(a, b, 0.5)
}

// Sum adds the values, accumulating in float64 and storing the total back through Cast, so a
// narrow integer T cannot overflow mid-sum; only a total outside its range is affected, as Cast
// documents. An empty or nil slice sums to 0.
func Sum[T Number](values []T) T {
	var total float64
	for _, value := range values {
		total += float64(value)
	}

	return Cast[T](total)
}

// Clamp adjusts the given value to be between the given minimum and maximum value.
func Clamp[T Number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

// Sign returns the sign of x: 1 if x > 0, -1 if x < 0, or 0 if x == 0.
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

// Equal reports whether a and b are equal within Epsilon of T: exactly for an integer T,
// within Delta32 for float32, within Delta for float64.
//
// The tolerance is absolute, which keeps the comparison to a subtraction and a compare
// for hot paths. It therefore falls below one ulp for values far from zero — beyond ~1e3
// for float32 and ~1e10 for float64 — where it degenerates into exact equality.
// Use EqualRelative when the magnitude is large or unknown.
func Equal[T Number](a, b T) bool {
	if isIntType[T]() {
		return a == b
	}

	return EqualDelta(a, b, Epsilon[T]())
}

// LessOrEqual reports whether a is at most b within Epsilon of T: exactly for an integer T,
// and with the same tolerance as Equal for a float T, so that a value a rounding error past
// a boundary still counts as on it.
func LessOrEqual[T Number](a, b T) bool {
	if isIntType[T]() {
		return a <= b
	}

	return LessOrEqualDelta(a, b, Epsilon[T]())
}

// LessOrEqualDelta reports whether a is at most b within the given delta.
// The comparison is made in float64, so it cannot overflow a narrow integer T.
func LessOrEqualDelta[T Number](a, b T, delta float64) bool {
	return float64(a) <= float64(b)+delta
}

// EqualDelta reports whether a and b are equal within the given delta.
// The difference is taken in float64, so it cannot overflow a narrow integer T.
func EqualDelta[T Number](a, b T, delta float64) bool {
	return math.Abs(float64(a)-float64(b)) <= delta
}

// EqualRelative reports whether a and b are equal within a tolerance that scales with
// their magnitude: Epsilon of T near zero, and Epsilon of T times the larger magnitude
// above it. Unlike Equal it holds at any scale, at the cost of the extra arithmetic.
func EqualRelative[T Number](a, b T) bool {
	if isIntType[T]() {
		return a == b
	}

	return EqualDelta(a, b, EpsilonRelative(a, b))
}

// Epsilon returns the equality tolerance for T: zero for an integer T, which is
// compared exactly, Delta32 for float32, and Delta for float64. It depends only on T,
// never on the values compared.
func Epsilon[T Number]() float64 {
	if isIntType[T]() {
		return 0
	}

	if isFloat32[T]() {
		return Delta32
	}

	return Delta
}

// EpsilonRelative returns the tolerance EqualRelative applies to a and b: Epsilon of T
// scaled by the larger magnitude, and never less than Epsilon of T itself.
func EpsilonRelative[T Number](a, b T) float64 {
	return Epsilon[T]() * max(1, math.Abs(float64(a)), math.Abs(float64(b)))
}

// Parse parses s into T: an integer T with strconv.ParseInt, a float T with strconv.ParseFloat.
// Parsing is done in int64 or float64 and narrowed to T, so defined types over those kinds
// (type Coord int) parse like their underlying type. A value outside the range of T is an error,
// never a wrapped integer or an overflowed infinity. The literals strconv.ParseFloat accepts,
// "NaN" and "Inf" among them, parse into a float T as they would into float64.
func Parse[T Number](s string) (T, error) {
	if isIntType[T]() {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}

		if int64(T(v)) != v {
			return 0, rangeError(s)
		}

		return T(v), nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	if math.IsInf(float64(T(v)), 0) && !math.IsInf(v, 0) {
		return 0, rangeError(s)
	}

	return T(v), nil
}

func rangeError(s string) error {
	return &strconv.NumError{Func: "Parse", Num: s, Err: strconv.ErrRange}
}

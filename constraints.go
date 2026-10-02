package geom

import (
	"fmt"
	"math"
)

// Integer is a generic integer type, supporting operations like modulo that floats don't.
//
// Every product, distance and interpolation is computed in float64 and stored back through
// Cast, so a narrow T such as int8 never overflows mid-computation; only a result outside its
// range is affected, as Cast documents. A value of an int64 T beyond 2^53 loses precision on
// the way through float64, and the cross products that decide a crossing, a turn or a winding
// are exact only while the coordinate differences stay within 2^26.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Float is a generic floating-point type.
type Float interface {
	~float32 | ~float64
}

// Number is a generic number type supported by all types and functions in this package.
type Number interface {
	Integer | Float
}

// Cast converts a float64 to T, rounding half away from zero for an integer T. NaN and ±Inf
// have no integer form, so for an integer T Cast panics on them, the same convention Divide
// follows for a zero factor. A float T keeps them.
//
// A finite value outside the range of an integer T is not checked: Go leaves that conversion
// platform-dependent, so a value beyond the range of T stores an arbitrary one.
// Keep results within the range of T, or use a wider T.
func Cast[T Number](a float64) T {
	if isInt[T]() {
		if math.IsNaN(a) || math.IsInf(a, 0) {
			panic("geom: cast of a non-finite value to an integer")
		}

		return T(math.Round(a))
	}

	return T(a)
}

// Int converts a Number to int: an integer T by plain conversion, exact wherever int is at
// least as wide as T and truncated otherwise (int64 on a 32-bit target), and a float T rounded
// through Cast.
func Int[T Number](value T) int {
	if isInt[T]() {
		return int(value)
	}

	return Cast[int](float64(value))
}

// String formats a Number as a numeric string: integer types without a decimal point,
// float types with two decimals. The formatting follows T, not the value. A value that rounds
// to zero prints as 0.00 without a sign, whether it is a negative zero or a small negative.
func String[T Number](value T) string {
	if isInt[T]() {
		return fmt.Sprintf("%d", int64(value))
	}

	s := fmt.Sprintf("%.2f", float64(value))
	if s == "-0.00" {
		return "0.00"
	}

	return s
}

// isInt reports whether T is an integer type.
func isInt[T Number]() bool {
	return T(1)/T(2) == 0
}

// isFloat32 reports whether T is a 32-bit float. 1e-10 lies between the float32 epsilon
// (1.2e-7) and the float64 one (2.2e-16), so adding it to one is a no-op for float32 alone.
// It is built by division because T(1e-10) does not compile when T may be an integer, and a
// constant above 127 overflows int8 at compile time even where the line never runs.
//
// It is asked only for a float T, after isInt, as Epsilon and epsilonAt ask it: for an integer
// T the division truncates to zero and the sum is one again.
//
// The arithmetic detection is deliberate: unsafe.Sizeof would be more direct, but the package
// stays free of the unsafe import.
func isFloat32[T Number]() bool {
	const step = 100

	one := T(1)

	return one+one/step/step/step/step/step == one
}

// castClamped is Cast with the value first clamped into the range of an integer T, so a value
// past either end stores that end rather than the arbitrary one Cast leaves it: the extent of a
// shape reaching past the range, which every shape it is tested against lies within. A float T
// is cast as it is.
func castClamped[T Number](a float64) T {
	if isInt[T]() {
		lowest, highest := intRange[T]()
		a = Clamp(a, lowest, highest)
	}

	return Cast[T](a)
}

// intRange returns the least and the greatest value of an integer T, as the float64 values
// nearest within it: the range of the narrowest width that holds the next value past it, found
// by converting that value and reading it back, as Parse checks a parsed one, so a defined type
// over any width is measured by its own. The greatest int64 is not a float64, and the one just
// below it stands in.
func intRange[T Number]() (float64, float64) {
	for _, bits := range [...]uint{7, 15, 31} {
		if past := int64(1) << bits; int64(T(past)) != past {
			return float64(-past), float64(past - 1)
		}
	}

	return -0x1p63, math.Nextafter(0x1p63, 0)
}

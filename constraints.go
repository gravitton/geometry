package geom

import (
	"fmt"
	"math"
)

// Integer is a generic integer type, supporting operations like modulo that floats don't.
//
// Every product, distance and interpolation is computed in float64 and stored back through
// Cast, so a value of an int64 T beyond 2^53 loses precision on the way. The package doc states
// the range each T is supported in.
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

// isFloat32 reports whether T is a 32-bit float: a third taken in T is the float64 third only
// where T is float64. It is asked only for a float T, after isInt, as Epsilon and epsilonAt ask
// it, since an integer third is zero.
//
// The arithmetic detection is deliberate: unsafe.Sizeof would be more direct, but the package
// stays free of the unsafe import.
func isFloat32[T Number]() bool {
	return float64(T(1)/T(3)) != 1.0/3
}

package geom

import (
	"fmt"
	"math"
)

// Integer is a generic integer type, supporting operations like modulo that floats don't.
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

// Cast number to type, round integer values.
func Cast[T Number](a float64) T {
	if isIntType[T]() {
		return T(math.Round(float64(a)))
	}

	return T(a)
}

// String formats a Number as a numeric string: integer types without a decimal point,
// float types with two decimals. The formatting follows T, not the value.
func String[T Number](value T) string {
	if isIntType[T]() {
		return fmt.Sprintf("%d", int64(value))
	}

	return fmt.Sprintf("%.2f", float64(value))
}

// isIntType reports whether T is an integer type.
func isIntType[T Number]() bool {
	return T(1)/T(2) == 0
}

// isFloat32 reports whether T is a 32-bit float. 1e-10 lies between the float32 epsilon
// (1.2e-7) and the float64 one (2.2e-16), so adding it to one is a no-op for float32 alone.
// It is built by division because T(1e-10) does not compile when T may be an integer, and a
// constant above 127 overflows int8 at compile time even where the line never runs.
func isFloat32[T Number]() bool {
	if isIntType[T]() {
		return false
	}

	const step = 100

	one := T(1)

	return one+one/step/step/step/step/step == one
}

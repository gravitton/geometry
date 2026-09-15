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
// float types with two decimals. The formatting follows T, not the value, so a float 8.0
// and 8.1 print alike.
func String[T Number](value T) string {
	if isIntType[T]() {
		return fmt.Sprintf("%d", int64(value))
	}

	return fmt.Sprintf("%.2f", float64(value))
}

func isIntType[T Number]() bool {
	return T(1)/T(2) == 0
}

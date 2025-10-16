package geom

import (
	"fmt"
	"math"
)

// Number is a generic number type supported by all types and functions in this package.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Cast number to type, round integer values.
func Cast[T Number](a float64) T {
	if isIntType[T]() {
		return T(math.Round(float64(a)))
	}

	return T(a)
}

// String format Number as numeric string
func String[T Number](value T) string {
	if isIntValue(value) {
		return fmt.Sprintf("%d", int64(value))
	} else {
		return fmt.Sprintf("%.2f", float64(value))
	}
}

func isIntType[T Number]() bool {
	var zero T
	switch any(zero).(type) {
	case float64, float32:
		return false
	case int, int64, int32, int16, int8:
		return true
	default:
		return false
	}
}

func isIntValue[T Number](value T) bool {
	return float64(value) == float64(int64(value))
}

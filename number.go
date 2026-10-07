package geom

import (
	"math"
	"strconv"
)

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
	buf := make([]byte, 0, stringSize)
	buf = appendString(buf, value)

	return string(buf)
}

// stringSize is the capacity every String starts its buffer at: it holds each type of fixed
// size at coordinates of a few digits, so the buffer stays on the stack and the string
// returned is the one allocation. A longer text grows the buffer like any append.
const stringSize = 128

// appendString appends the number to buf and returns the extended buffer, as String prints it.
func appendString[T Number](buf []byte, value T) []byte {
	if isInt[T]() {
		return strconv.AppendInt(buf, int64(value), 10)
	}

	start := len(buf)
	buf = strconv.AppendFloat(buf, float64(value), 'f', 2, 64)
	if string(buf[start:]) == "-0.00" {
		return append(buf[:start], "0.00"...)
	}

	return buf
}

// Parse parses s into T: an integer T with strconv.ParseInt, a float T with strconv.ParseFloat.
// Parsing is done in int64, float32 or float64 and narrowed to T, so defined types over those
// kinds (type Coord int) parse like their underlying type, and a float32 T is rounded once,
// correctly, rather than through float64. A value outside the range of T is an error,
// never a wrapped integer or an overflowed infinity. The literals strconv.ParseFloat accepts,
// "NaN" and "Inf" among them, parse into a float T as they would into float64.
func Parse[T Number](s string) (T, error) {
	if isInt[T]() {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}

		if int64(T(v)) != v {
			return 0, rangeError(s)
		}

		return T(v), nil
	}

	bitSize := 64
	if isFloat32[T]() {
		bitSize = 32
	}

	v, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

func rangeError(s string) error {
	return &strconv.NumError{Func: "Parse", Num: s, Err: strconv.ErrRange}
}

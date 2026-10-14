package geom

import (
	"math"
	"strconv"
)

// Cast converts a float64 to T, rounding half away from zero for an integer T. It panics for a
// non-finite value and an integer T; a float T keeps it. A finite value outside the range of an
// integer T is not checked and stores an arbitrary value.
func Cast[T Number](a float64) T {
	if isInt[T]() {
		if !IsFinite(a) {
			panic("geom: cast of a non-finite value to an integer")
		}

		return T(math.Round(a))
	}

	return T(a)
}

// castUp converts a float64 to T as Cast does, but rounds up for an integer T: the cast of an
// extent that must still reach what it was measured to, the radius of a BoundingCircle.
func castUp[T Number](a float64) T {
	if isInt[T]() {
		return Cast[T](math.Ceil(a))
	}

	return Cast[T](a)
}

// Int converts a Number to int: an integer T by plain conversion, truncated where int is narrower
// than T, and a float T through Cast.
func Int[T Number](value T) int {
	if isInt[T]() {
		return int(value)
	}

	return Cast[int](float64(value))
}

// String returns the value formatted: an integer T without a decimal point, a float T with two
// decimals. A value that rounds to zero prints without a sign.
func String[T Number](value T) string {
	buf := make([]byte, 0, stringSize)
	buf = appendString(buf, value)

	return string(buf)
}

// stringSize is the capacity every String starts its buffer at, enough for each fixed-size type at
// coordinates of a few digits, so the buffer stays on the stack and the result is the one
// allocation.
const stringSize = 128

// appendString appends the number to buf as String prints it and returns the extended buffer.
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

// Parse returns s as T, and an error where it is no number T can hold. An integer T is read with
// strconv.ParseInt and a float T with strconv.ParseFloat at the width of T, so a float32 is rounded
// once. Defined types parse like their underlying type, and the literals strconv.ParseFloat
// accepts, NaN and Inf among them, parse into a float T.
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

// rangeError returns the error strconv reports for a number out of range, which Parse gives for a
// value that does not fit T.
func rangeError(s string) error {
	return &strconv.NumError{Func: "Parse", Num: s, Err: strconv.ErrRange}
}

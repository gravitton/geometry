package geom

import (
	"fmt"
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

	// Delta is the tolerance used for floating-point equality comparisons.
	Delta float64 = 1e-6

	// Sqrt2 is the square root of 2.
	Sqrt2 = math.Sqrt2
	// Sqrt3 is the square root of 3.
	Sqrt3 = 1.732050807568877293527446341505872367
	// OneOverSqrt2 is 1/√2, the length of a unit diagonal vector component.
	OneOverSqrt2 = 1 / math.Sqrt2
)

// NormalizeAngle returns angle normalized to [0, 2π).
func NormalizeAngle(angle float64) float64 {
	a := math.Mod(angle, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
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

// Divide divides a number by a scale factor; if scale is zero, returns the original value.
func Divide[T Number](x T, scale float64) T {
	if scale == 0 {
		return x
	}

	return Cast[T](float64(x) / scale)
}

// Abs returns the absolute value.
func Abs[T Number](x T) T {
	return Cast[T](math.Abs(float64(x)))
}

// Round returns x rounded to the nearest integer.
func Round[T Number](x T) T {
	return Cast[T](math.Round(float64(x)))
}

// Floor returns the largest integer value less than or equal to x.
func Floor[T Number](x T) T {
	return Cast[T](math.Floor(float64(x)))
}

// Ceil returns the smallest integer value greater than or equal to x.
func Ceil[T Number](x T) T {
	return Cast[T](math.Ceil(float64(x)))
}

// Mod wraps n into [0, m), correctly for negative n.
func Mod[T Integer](n, m T) T {
	return ((n % m) + m) % m
}

// Lerp calculates the linear interpolation between a and b at a ratio t.
func Lerp[T Number](a, b T, t float64) T {
	return Cast[T](float64(a) + float64(b-a)*t)
}

// Midpoint calculates the midpoint between two values. Equivalent to Lerp(a, b, 0.5).
func Midpoint[T Number](a, b T) T {
	return Lerp(a, b, 0.5)
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

// Equal reports whether a and b are equal within Delta.
func Equal[T Number](a, b T) bool {
	return EqualDelta(a, b, Delta)
}

// EqualDelta reports whether a and b are equal within the given delta.
func EqualDelta[T Number](a, b T, delta float64) bool {
	return equalDelta(float64(a), float64(b), delta)
}

func equalDelta(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

// Parse parses s into T using the parser and bit size appropriate for Number type.
func Parse[T Number](s string) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int8:
		v, err := strconv.ParseInt(s, 10, 8)
		return T(v), err
	case int16:
		v, err := strconv.ParseInt(s, 10, 16)
		return T(v), err
	case int32:
		v, err := strconv.ParseInt(s, 10, 32)
		return T(v), err
	case int64, int:
		v, err := strconv.ParseInt(s, 10, 64)
		return T(v), err
	case float32:
		v, err := strconv.ParseFloat(s, 32)
		return T(v), err
	case float64:
		v, err := strconv.ParseFloat(s, 64)
		return T(v), err
	default:
		return 0, fmt.Errorf("unsupported number type %T", zero)
	}
}

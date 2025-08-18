package geom

import (
	"math"
)

const (
	RadToDeg float64 = 180.0 / math.Pi
	DegToRad float64 = math.Pi / 180.0

	Delta float64 = 1e-6
)

// ToRadians converts degrees to radians
func ToRadians(degrees float64) float64 {
	return degrees * DegToRad
}

// ToDegrees converts radians to degrees
func ToDegrees(radians float64) float64 {
	return radians * RadToDeg
}

// Scale multiple number by scale factor
func Scale[T Number](a T, scale float64) T {
	return Cast[T](float64(a) * scale)
}

// Abs returns absolute value
func Abs[T Number](a T) T {
	return Cast[T](math.Abs(float64(a)))
}

// Midpoint calculate midpoint (point exactly halfway between two points)
// Shorthand for `lerp(a, b, 0.5)`
// TODO: fix rounding for int numbers
func Midpoint[T Number](a, b T) T {
	// return Cast(lerp(float64(a), float64(b), 0.5))
	return Cast[T](midpoint(float64(a), float64(b)))
}

func midpoint(a, b float64) float64 {
	// optimized `a + (b-a)/2.0`
	// return lerp(a, b, 0.5)
	return (a + b) / 2.0
}

// Lerp calculate linear interpolation (point along a line between two points based on a given ratio)
// TODO: fix rounding for int numbers
func Lerp[T Number](a, b T, t float64) T {
	return Cast[T](lerp(float64(a), float64(b), t))
}

func lerp(a, b, t float64) float64 {
	// optimized `a + (b-a)*t`
	return a*(1-t) + b*t
}

// Equal checks for nearly-equal values in Delta
func Equal[T Number](a, b T) bool {
	return EqualDelta(a, b, Delta)
}

// Equal checks for nearly-equal values in given delta
func EqualDelta[T Number](a, b T, delta float64) bool {
	return equalDelta(float64(a), float64(b), delta)
}

func equalDelta(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

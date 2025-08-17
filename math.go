package geom

import "math"

const (
	Delta float64 = 1e-6
)

// Midpoint calculate midpoint (point exactly halfway between two points)
// Shorthand for `lerp(a, b, 0.5)`
// TODO: fix rounding for int numbers
func Midpoint[T Number](a, b T) T {
	return T(midpoint(float64(a), float64(b)))
}

func midpoint(a, b float64) float64 {
	// optimized `a + (b-a)/2.0`
	// return lerp(a, b, 0.5)
	return (a + b) / 2.0
}

// Lerp calculate linear interpolation (point along a line between two points based on a given ratio)
// TODO: fix rounding for int numbers
func Lerp[T Number](a, b T, t float64) T {
	return T(lerp(float64(a), float64(b), t))
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

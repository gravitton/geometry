package geom

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

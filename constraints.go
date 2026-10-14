package geom

// Integer is the integer types the package supports.
//
// Products, distances and interpolations are computed in float64, so an int64 T beyond 2^53 loses
// precision.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Float is the floating-point types the package supports.
type Float interface {
	~float32 | ~float64
}

// Number is every type the package supports.
type Number interface {
	Integer | Float
}

// isInt reports whether T is an integer type.
func isInt[T Number]() bool {
	return T(1)/T(2) == 0
}

// isFloat32 reports whether a float T is 32-bit: a third taken in T is the float64 third only in
// float64. It is asked only after isInt, since an integer third is zero, and avoids unsafe.Sizeof
// to keep the package free of the unsafe import.
func isFloat32[T Number]() bool {
	return float64(T(1)/T(3)) != 1.0/3
}

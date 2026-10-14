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

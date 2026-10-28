package geom

import (
	"errors"
)

var (
	// ErrDivisionByZero is the panic value of a division by a zero factor.
	ErrDivisionByZero = errors.New("geom: division by zero")
	// ErrSingularMatrix is the panic value of inverting a matrix with a zero determinant.
	ErrSingularMatrix = errors.New("geom: inverse of a singular matrix")
	// ErrNonFinite is the panic value of storing NaN or ±Inf into an integer T.
	ErrNonFinite = errors.New("geom: cast of a non-finite value to an integer")
	// ErrCountMismatch is wrapped by the panic value of interpolating between outlines of
	// different vertex or point counts.
	ErrCountMismatch = errors.New("geom: lerp between different counts")
	// ErrUnknownOrientation is wrapped by the panic value of an orientation that is neither
	// OrientationFlatTop nor OrientationPointyTop.
	ErrUnknownOrientation = errors.New("geom: unknown orientation")
)

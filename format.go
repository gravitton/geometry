package geom

import (
	"fmt"
)

// IsInt checks if Number is integer value
func IsInt[T Number](value T) bool {
	return float64(value) == float64(int64(value))
}

// ToString format Number as numeric string
func ToString[T Number](value T) string {
	if IsInt(value) {
		return fmt.Sprintf("%+d", int64(value))
	} else {
		return fmt.Sprintf("%+.2f", float64(value))
	}
}

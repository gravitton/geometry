package internal

// Map applies a function to each element in a slice and returns a new slice.
func Map[S ~[]E, E any, T any](input S, fn func(E) T) []T {
	output := make([]T, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}

	return output
}

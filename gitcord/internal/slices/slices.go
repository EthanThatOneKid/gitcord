// Package slices contains convenient generic functions for slices.
package slices

// Find returns the pointer to the first element that satisfies the predicate.
func Find[T any](slice []T, predicate func(*T) bool) *T {
	for i := range slice {
		if predicate(&slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

// Filter returns a new slice containing all elements that satisfy the
// predicate.
func Filter[T any](slice []T, predicate func(*T) bool) []T {
	filtered := make([]T, 0, len(slice))
	for i, v := range slice {
		if predicate(&slice[i]) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// FilterReuse is like Filter, but the returned slice will reuse the same array
// buffer as the given slice. Use this function only if you know that the old
// slice will never be used again.
func FilterReuse[T any](slice []T, predicate func(*T) bool) []T {
	filtered := slice[:0]
	for i, v := range slice {
		if predicate(&slice[i]) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

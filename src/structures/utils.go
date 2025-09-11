package structures

import (
	"cmp"
	"slices"
)

// SliceReduce returns a copy of original with no duplicate, sorted.
// Due to the order, it applies only to cmp.Ordered
func SliceReduce[T cmp.Ordered](original []T) []T {
	elements := make(map[T]bool)
	for _, value := range original {
		elements[value] = true
	}

	var result []T
	for k := range elements {
		result = append(result, k)
	}

	if len(result) == 0 {
		return []T{}
	}

	slices.Sort(result)
	return result
}

// SliceDeduplicate returns the slice content with one value only from original slice
func SliceDeduplicate[T comparable](original []T) []T {
	var result []T
	seen := make(map[T]bool)
	for _, v := range original {
		seen[v] = true
	}

	for k := range seen {
		result = append(result, k)
	}

	return result
}

// SliceDeduplicateFunc returns a slice containing the same elements, just once
func SliceDeduplicateFunc[T any](original []T, equals func(a, b T) bool) []T {
	var result []T
	for _, source := range original {
		if !slices.ContainsFunc(result, func(value T) bool { return equals(source, value) }) {
			result = append(result, source)
		}
	}

	return result
}

// SliceCommonElement returns true if there is a common element in the slices
func SliceCommonElement[T comparable](first, second []T) bool {
	values := make(map[T]bool)
	for _, k := range first {
		values[k] = true
	}

	for _, v := range second {
		if values[v] {
			return true
		}
	}

	return false
}

// SliceCommonElement returns true if there is a common element in the slices based on a equals test
func SliceCommonElementFunc[T any](first, second []T, equalsFunc func(a, b T) bool) bool {
	for _, source := range first {
		if slices.ContainsFunc(second, func(val T) bool { return equalsFunc(val, source) }) {
			return true
		}
	}

	return false
}

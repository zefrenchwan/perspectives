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

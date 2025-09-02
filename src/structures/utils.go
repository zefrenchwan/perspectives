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

package commons

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
)

// NewId builds a new unique id.
// Two different calls should return two different values.
func NewId() string {
	return uuid.NewString()
}

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

// SlicesEqualsAsSetsFunc returns true if slices have the same elements based on a func
func SlicesEqualsAsSetsFunc[T any](first, second []T, equalsFunc func(a, b T) bool) bool {
	for _, base := range first {
		if !slices.ContainsFunc(second, func(element T) bool { return equalsFunc(element, base) }) {
			return false
		}
	}

	for _, base := range second {
		if !slices.ContainsFunc(first, func(element T) bool { return equalsFunc(element, base) }) {
			return false
		}
	}

	return true
}

// MapsReverseFind returns all the key containing the value in values
func MapsReverseFind[T comparable](mapping map[T][]T, value T) []T {
	var result []T
	for key, values := range mapping {
		if slices.Contains(values, value) {
			result = append(result, key)
		}
	}

	return result
}

// SlicesContainsAll returns true if other is included in base based on an equals function.
// In other words, it returns true if base contains other as a set based on equality.
func SlicesContainsAllFunc[T any](base []T, other []T, equals func(a, b T) bool) bool {
	if len(other) == 0 {
		return true
	} else if len(base) == 0 {
		return false
	}

	for _, value := range other {
		if !slices.ContainsFunc(base, func(e T) bool { return equals(e, value) }) {
			return false
		}
	}

	return true
}

// SlicesFilter returns a new slice containing only  elements that match the predicate
func SlicesFilter[T any](base []T, keepPredicate func(T) bool) []T {
	var result []T
	for _, element := range base {
		if keepPredicate == nil || keepPredicate(element) {
			result = append(result, element)
		}
	}

	return result
}

package values

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// EnsureValuesMappingInvariant goes through a mapping and ensures that all values have the same type.
// It returns true if all values have the same expected type and false otherwise.
func EnsureValuesMappingInvariant[V Value](rawMapping periods.DynamicMapping[V]) bool {
	var expectedType string
	for _, value := range rawMapping.Range() {
		if expectedType == "" {
			expectedType = value.Datatype()
		} else if value.Datatype() != expectedType {
			return false
		}
	}

	if expectedType == "" {
		// no element, so ok so far
		return true
	}

	// ALl elements have the same type and it is what we expect ? OK
	return expectedType == rawMapping.DataType()
}

// ImmutableValuesMapping is an immutable mapping of periods to values (reference or primitive values).
// It is used to represent a mapping of periods to values that cannot be modified after creation.
type ImmutableValuesMapping[V Value] interface {
	// Hashable is an interface that provides a hash function for the mapping.
	// Because mapping is immutable, it will not change after creation.
	commons.Hashable
	// IsEmpty returns true if the mapping is empty and false otherwise.
	IsEmpty() bool
	// Range returns an iterator over periods and values in the mapping.
	Range() iter.Seq2[periods.Period, V]
	// ValuesType returns the type of values in the mapping.
	// For instance, on a primitive mapping, it will return "int", "string", etc.
	ValuesType() string
}

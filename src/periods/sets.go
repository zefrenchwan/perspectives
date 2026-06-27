package periods

import (
	"iter"
	"time"
)

// dynamicSet is a dynamic set : multiple values per period.
type dynamicSet[T any] struct {
	// underlying values handler to regroup code as much as possible.
	*valuesHandler[T]
}

// isPartition returns false by definition
func (ds *dynamicSet[T]) isPartition() bool {
	return false
}

// At returns all values at a given moment.
// Result is an iterator over the values to reduce slices copies and preserve GC.
func (ds *dynamicSet[T]) At(moment time.Time) (iter.Seq[T], bool) {
	return ds.all(moment)
}

// Copy returns a copy of the dynamic set.
func (ds *dynamicSet[T]) Copy() DynamicSet[T] {
	return &dynamicSet[T]{
		valuesHandler: ds.clone(),
	}
}

// NewDynamicSet creates a new dynamic set.
// Parameters are the data type and an equality function.
// There is NO test if parameters make sense, it should be done by the caller.
func NewDynamicSet[T any](dataType string, equals func(T, T) bool) DynamicSet[T] {
	handler := &valuesHandler[T]{
		values:      make([]valueNode[T], 0),
		storedType:  dataType,
		isPartition: false,
		equals:      equals,
	}

	return &dynamicSet[T]{
		valuesHandler: handler,
	}
}

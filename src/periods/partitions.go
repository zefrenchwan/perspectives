package periods

import "time"

// dynamicPartition is a dynamic partition per period : one value per period.
type dynamicPartition[T any] struct {
	*valuesHandler[T]
}

// isPartition returns true by definition
func (dp *dynamicPartition[T]) isPartition() bool {
	return true
}

// At returns the unique element (if any) matching the given moment.
func (dp *dynamicPartition[T]) At(moment time.Time) (T, bool) {
	return dp.first(moment)
}

// Copy returns a copy of the dynamic partition.
func (dp *dynamicPartition[T]) Copy() DynamicPartition[T] {
	return &dynamicPartition[T]{
		valuesHandler: dp.clone(),
	}
}

// NewDynamicPartition creates a new dynamic partition per periods.
// Parameters are the data type and an equality function.
// There is NO test if parameters make sense, it should be done by the caller.
func NewDynamicPartition[T any](dataType string, equals func(T, T) bool) DynamicPartition[T] {
	handler := &valuesHandler[T]{
		values:      make([]valueNode[T], 0),
		storedType:  dataType,
		isPartition: true,
		equals:      equals,
	}

	return &dynamicPartition[T]{
		valuesHandler: handler,
	}
}

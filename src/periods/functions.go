package periods

import (
	"time"
)

// DynamicFunction defines a dynamic mapping with ONE value maximum at a given period.
// It is a mapping of values, each value being valid during a specific period.
type DynamicFunction[T any] interface {
	// DynamicMapping to regroup the common methods.
	DynamicMapping[T]
	// At returns the unique element (if any) matching the given moment.
	At(moment time.Time) (T, bool)
	// Copy returns a copy of the dynamic function.
	Copy() DynamicFunction[T]
}

// HashDynamicFunction returns a hash of the given dynamic function.
func HashDynamicFunction[T any](f DynamicFunction[T]) string {
	return HashDynamicMapping(f, true)
}

// timeFunction is a dynamic partition per period : one value per period.
// Then, as a function, it picks the unique (if any) element matching the given moment.
type timeFunction[T any] struct {
	*valuesHandler[T]
}

// At returns the unique element (if any) matching the given moment.
func (dp *timeFunction[T]) At(moment time.Time) (T, bool) {
	return dp.first(moment)
}

// Copy returns a copy of the current function.
func (dp *timeFunction[T]) Copy() DynamicFunction[T] {
	return &timeFunction[T]{
		valuesHandler: dp.clone(),
	}
}

// NewTimeFunction creates a new time dependent function of T as a dynamic partition per periods.
// Parameters are the data type and an equality function.
// There is NO test if parameters make sense, it should be done by the caller.
func NewTimeFunction[T any](dataType string, equals func(T, T) bool) DynamicFunction[T] {
	handler := &valuesHandler[T]{
		values:     make([]valueNode[T], 0),
		storedType: dataType,
		isFunction: true,
		equals:     equals,
	}

	return &timeFunction[T]{
		valuesHandler: handler,
	}
}

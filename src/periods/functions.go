package periods

import "time"

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

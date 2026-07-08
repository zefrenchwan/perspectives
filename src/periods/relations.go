package periods

import (
	"iter"
	"time"
)

// timeRelation is a time-dependent relation : multiple values per period.
type timeRelation[T any] struct {
	// underlying values handler to regroup code as much as possible.
	*valuesHandler[T]
}

// At returns all values at a given moment.
// Result is an iterator over the values to reduce slices copies and preserve GC.
func (tr *timeRelation[T]) At(moment time.Time) (iter.Seq[T], bool) {
	return tr.all(moment)
}

// Copy returns a copy of the current relation.
func (tr *timeRelation[T]) Copy() DynamicRelation[T] {
	return &timeRelation[T]{
		valuesHandler: tr.clone(),
	}
}

// NewTimeRelation creates a new time dependent relation with multiple values per period.
// Parameters are the data type and an equality function.
// There is NO test if parameters make sense, it should be done by the caller.
func NewTimeRelation[T any](dataType string, equals func(T, T) bool) DynamicRelation[T] {
	handler := &valuesHandler[T]{
		values:     make([]valueNode[T], 0),
		storedType: dataType,
		isFunction: false,
		equals:     equals,
	}

	return &timeRelation[T]{
		valuesHandler: handler,
	}
}

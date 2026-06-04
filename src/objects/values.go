package objects

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/configuration"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// TemporalValues represents a collection of values with associated time periods.
// It depends on a type T assumed to be the primitive type of the values stored per period.
type TemporalValues[T any] interface {
	// IsEmpty checks if the TemporalValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// Add a value for a given period
	Add(period periods.Period, value T)
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (T, bool)
	// Clear removes all values from the TemporalValues collection
	Clear()
	// Remove removes all values for a given period
	Remove(period periods.Period)
	// Cut returns a new TemporalValues collection containing only values within the specified period
	Cut(period periods.Period) TemporalValues[T]
	// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, T) bool)
}

// valueNode stores a value set during a specific matchingPeriod
// T is the type of the value stored in the node
type valueNode[T any] struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod periods.Period
	// value is the actual value stored in the node
	value T
}

// valuesHandler manages the full history of values with their respective matching periods
type valuesHandler[T any] struct {
	// values : one value per matching period
	values []valueNode[T]
	// equals : function to compare two values of type T
	equals func(T, T) bool
}

// IsEmpty checks if the valuesHandler contains any values
func (vh *valuesHandler[T]) IsEmpty() bool {
	return len(vh.values) == 0
}

// Add adds a new value with a specific matchingPeriod to the valuesHandler
func (vh *valuesHandler[T]) Add(p periods.Period, v T) {
	matchingPeriodValue := p
	for _, element := range vh.values {
		if vh.equals(element.value, v) {
			matchingPeriodValue = matchingPeriodValue.Union(element.matchingPeriod)
		}
	}

	result := make([]valueNode[T], 0, len(vh.values)+1)
	for _, element := range vh.values {
		if !vh.equals(element.value, v) {
			remaining := element.matchingPeriod.Remove(matchingPeriodValue)
			if !remaining.IsEmpty() {
				result = append(result, valueNode[T]{matchingPeriod: remaining, value: element.value})
			}
		}
	}

	if !matchingPeriodValue.IsEmpty() {
		result = append(result, valueNode[T]{matchingPeriod: matchingPeriodValue, value: v})
	}

	vh.values = result
}

// At returns the value at the given moment in time, or the zero value and false if no value is found.
func (vh *valuesHandler[T]) At(moment time.Time) (T, bool) {
	var empty T
	for _, element := range vh.values {
		if element.matchingPeriod.Contains(moment) {
			return element.value, true
		}
	}
	return empty, false
}

// Remove removes the given period from the values handler, if the period is empty or the handler is empty, it does nothing.
func (vh *valuesHandler[T]) Remove(period periods.Period) {
	if period.IsEmpty() || len(vh.values) == 0 {
		return
	}

	result := make([]valueNode[T], 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode[T]{matchingPeriod: remaining, value: element.value})
		}
	}

	vh.values = result
}

// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
func (vh *valuesHandler[T]) Range(yield func(periods.Period, T) bool) {
	for _, element := range vh.values {
		if !yield(element.matchingPeriod, element.value) {
			break
		}
	}
}

// Cut returns a copy with same values, restricted to given period
func (vh *valuesHandler[T]) Cut(period periods.Period) TemporalValues[T] {
	remainingValues := make([]valueNode[T], 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode[T]{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler[T]{values: remainingValues, equals: vh.equals}
}

// Clear removes all the values
func (vh *valuesHandler[T]) Clear() {
	vh.values = nil
}

// =================================================
// EQUALITY HANDLERS
// =================================================

// intEquals is just a common func definition to reuse for int values
func intEquals(a, b int) bool { return a == b }

// stringEquals is just a common func definition to reuse for string values
func stringEquals(a, b string) bool { return a == b }

// floatEquals is just a common func definition to reuse for float values.
// It uses configuration.LONG_EPSILON to determine if two float values are considered equal.
func floatEquals(a, b float64) bool {
	diff := a - b
	if diff < 0.0 {
		diff = -diff
	}

	return diff < configuration.LONG_EPSILON
}

// =====================================================
// BASIC TYPES : MANAGE INT, FLOATS AND STRINGS
// =====================================================

// stringValues is a dedicated type for time-dependent string values
type stringValues = valuesHandler[string]

// intValues is a dedicated type for time-dependent int values
type intValues = valuesHandler[int]

// floatValues is a dedicated type for time-dependent float values
type floatValues = valuesHandler[float64]

// newIntHandler creates a new empty intValues instance with the default equality handler
func newIntHandler() *intValues {
	return &intValues{equals: intEquals}
}

// newStringHandler creates a new empty stringValues instance with the default equality handler
func newStringHandler() *stringValues {
	return &valuesHandler[string]{equals: stringEquals}
}

// newFloatHandler creates a new empty floatValues instance with the default equality handler
func newFloatHandler() *floatValues {
	return &floatValues{equals: floatEquals}
}

// NewTemporalStringValues creates a new temporal values manager for string values
func NewTemporalStringValues() TemporalValues[string] {
	return newStringHandler()
}

// NewTemporalIntValues creates a new temporal values manager for int values
func NewTemporalIntValues() TemporalValues[int] {
	return newIntHandler()
}

// NewTemporalFloatValues creates a new temporal values manager for float values
func NewTemporalFloatValues() TemporalValues[float64] {
	return newFloatHandler()
}

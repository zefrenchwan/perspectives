package objects

import (
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// TemporalValues represents a collection of values with associated time periods.
// It uses "any" to store any type of values per period.
type TemporalValues interface {
	// IsEmpty checks if the TemporalValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// Add a value for a given period
	Add(period periods.Period, value any)
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Clear removes all values from the TemporalValues collection
	Clear()
	// Remove removes all values for a given period
	Remove(period periods.Period)
	// Cut returns a new TemporalValues collection containing only values within the specified period
	Cut(period periods.Period) TemporalValues
	// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, any) bool)
}

// =========================================================================
// TEMPORAL VALUES IMPLEMENTATION
// =========================================================================

// valueNode stores a value set during a specific matchingPeriod
// value is the actual value (of type any) stored in the node
type valueNode struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod periods.Period
	// value is the actual value stored in the node
	value any
}

// valuesHandler manages the full history of values with their respective matching periods
type valuesHandler struct {
	// values have one value per matching period
	values []valueNode
}

// IsEmpty checks if the valuesHandler contains any values
func (vh *valuesHandler) IsEmpty() bool {
	return len(vh.values) == 0
}

// Add adds a new value with a specific matchingPeriod to the valuesHandler
func (vh *valuesHandler) Add(p periods.Period, v any) {
	matchingPeriodValue := p
	for _, element := range vh.values {
		if reflect.DeepEqual(element.value, v) {
			matchingPeriodValue = matchingPeriodValue.Union(element.matchingPeriod)
		}
	}

	result := make([]valueNode, 0, len(vh.values)+1)
	for _, element := range vh.values {
		if !reflect.DeepEqual(element.value, v) {
			remaining := element.matchingPeriod.Remove(matchingPeriodValue)
			if !remaining.IsEmpty() {
				result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
			}
		}
	}

	if !matchingPeriodValue.IsEmpty() {
		result = append(result, valueNode{matchingPeriod: matchingPeriodValue, value: v})
	}

	vh.values = result
}

// At returns the value at the given moment in time, or nil and false if no value is found.
func (vh *valuesHandler) At(moment time.Time) (any, bool) {
	for _, element := range vh.values {
		if element.matchingPeriod.Contains(moment) {
			return element.value, true
		}
	}
	return nil, false
}

// Remove removes the given period from the values handler, if the period is empty or the handler is empty, it does nothing.
func (vh *valuesHandler) Remove(period periods.Period) {
	if period.IsEmpty() || len(vh.values) == 0 {
		return
	}

	result := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	vh.values = result
}

// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
func (vh *valuesHandler) Range(yield func(periods.Period, any) bool) {
	for _, element := range vh.values {
		if !yield(element.matchingPeriod, element.value) {
			break
		}
	}
}

// Cut returns a copy with same values, restricted to given period
func (vh *valuesHandler) Cut(period periods.Period) TemporalValues {
	remainingValues := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: remainingValues}
}

// Clear removes all the values
func (vh *valuesHandler) Clear() {
	vh.values = nil
}

func NewTemporalValues() TemporalValues {
	return &valuesHandler{}
}

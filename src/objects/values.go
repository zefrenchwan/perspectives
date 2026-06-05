package objects

import (
	"fmt"
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// TemporalValues represents a collection of values with associated time periods.
// It uses "any" to store any type of values per period.
type TemporalValues interface {
	// Same returns true if content is the same as another TemporalValues
	Same(other TemporalValues) bool
	// IsEmpty checks if the TemporalValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// Add a value for a given period
	Add(period periods.Period, value any) TemporalValues
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Remove removes all values for a given period
	Remove(period periods.Period) TemporalValues
	// Cut returns a new TemporalValues collection containing only values within the specified period
	Cut(period periods.Period) TemporalValues
	// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, any) bool)
	// DataType returns the type of values stored in the TemporalValues collection.
	// It looks for the most common type among all values, or returns "any" if types are diverse.
	// For instance, if all values are integers, it will return "int". If there are both integers and strings, it will return "any".
	// Special case for empty collection: returns ""
	DataType() string
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

func (vh *valuesHandler) Same(other TemporalValues) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if vh.IsEmpty() {
		return true
	}

	// TODO : implement this

	return true
}

// IsEmpty checks if the valuesHandler contains any values
func (vh *valuesHandler) IsEmpty() bool {
	return len(vh.values) == 0
}

// Add adds a new value with a specific matchingPeriod to the valuesHandler
func (vh *valuesHandler) Add(p periods.Period, v any) TemporalValues {
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

	return &valuesHandler{values: result}
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
func (vh *valuesHandler) Remove(period periods.Period) TemporalValues {
	if period.IsEmpty() || len(vh.values) == 0 {
		return &valuesHandler{}
	}

	result := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: result}
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

// DataType returns the string representation of the common type of all stored values or "any" if types differ.
func (vh *valuesHandler) DataType() string {
	if vh == nil || len(vh.values) == 0 {
		return ""
	}

	var commonType string
	isFirst := true

	for _, element := range vh.values {
		currentType := fmt.Sprintf("%T", element.value)

		if isFirst {
			commonType = currentType
			isFirst = false
			continue
		}

		if currentType != commonType {
			return "any"
		}
	}

	return commonType
}

func NewTemporalValues() TemporalValues {
	return &valuesHandler{}
}

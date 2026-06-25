package entities

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// AttributeDetails represents the metadata details of the attribute.
// It contains information about the attribute's name, type, validity, and instance activity.
type AttributeDetails struct {
	// AttributeName is the actual name of the attribute
	AttributeName string
	// AttributeType is the actual type of the attribute
	AttributeType string
	// AttributeValidity is the validity period of the attribute
	AttributeValidity periods.Period
	// InstanceActivity is the activity period of the instance
	InstanceActivity periods.Period
}

// DynamicValues represents a value that depends on time.
// It is basically equivalent to a map of disjoined time intervals linked to primitive values.
// Implementations have to ensure that value accepts only PrimitiveValue types.
type DynamicValues interface {
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Same returns true if instance is the same as another DynamicValues.
	// It means : same periods, same values, same underlying type
	Same(other DynamicValues) bool
	// IsEmpty checks if the TimeDependentValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Range iterates over all values in the TimeDependentValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, any) bool)
	// DataType returns the type name of the stored values.
	// By design, it should be the same at all times
	DataType() string
}

// =========================================================================
// TIME DEPENDENT VALUES IMPLEMENTATION : in memory, no storage
// =========================================================================

// valueNode stores a value set during a specific matchingPeriod
// value is the actual value (of type any) stored in the node.
type valueNode struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod periods.Period
	// value is the actual value stored in the node
	value any
}

// valuesHandler manages the full history of values with their respective matching periods.
// Its purpose is to provide a way to store and retrieve values over time.
// KEY INVARIANT : storedType is the actual type (should be primitive) and should be unique over time.
// There is NO LOCK at all, because it is immutable by design.
type valuesHandler struct {
	// values have one value per matching period
	values []valueNode
	// storedType is the actual type name of the content (should be primitive)
	storedType string
	// equality function
	equals func(a, b any) bool
}

// Same returns true if the two temporal values have the same values at the same periods, and same type
func (vh *valuesHandler) Same(other DynamicValues) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if vh.IsEmpty() {
		return true
	} else if vh.storedType != other.DataType() {
		return false
	}

	counter := 0
	for period, value := range other.Range {
		counter++
		found := false
		// find matching element if any
		for _, matching := range vh.values {
			if period.Equals(matching.matchingPeriod) {
				found = true
				if !vh.equals(matching.value, value) {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return counter == len(vh.values)
}

// IsEmpty checks if the valuesHandler contains any values
func (vh *valuesHandler) IsEmpty() bool {
	return vh == nil || len(vh.values) == 0
}

// Validity returns the union of periods for which values are set
func (vh *valuesHandler) Validity() periods.Period {
	if vh == nil || len(vh.values) == 0 {
		return periods.NewEmptyPeriod()
	}

	validity := periods.NewEmptyPeriod()
	for _, element := range vh.values {
		validity = validity.Union(element.matchingPeriod)
	}

	return validity
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

// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
func (vh *valuesHandler) Range(yield func(periods.Period, any) bool) {
	for _, element := range vh.values {
		if !yield(element.matchingPeriod, element.value) {
			break
		}
	}
}

// DataType returns the string representation of the type.
// It is UNIQUE by design : this is a key invariant.
func (vh *valuesHandler) DataType() string {
	return vh.storedType
}

// Copy returns a copy of the valuesHandler with the same values and type
func (vh *valuesHandler) Copy() *valuesHandler {
	result := make([]valueNode, len(vh.values))
	copy(result, vh.values)
	return &valuesHandler{values: result, storedType: vh.storedType, equals: vh.equals}
}

// withoutValidity returns a copy without values for the given period.
// If the period is empty or the handler is empty, it does nothing.
func (vh *valuesHandler) withoutValidity(period periods.Period) *valuesHandler {
	if len(vh.values) == 0 {
		return &valuesHandler{storedType: vh.storedType}
	} else if period.IsEmpty() {
		return vh
	}

	result := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: result, storedType: vh.storedType, equals: vh.equals}
}

// cut returns a copy with same values, restricted to a given period
func (vh *valuesHandler) cut(period periods.Period) *valuesHandler {
	remainingValues := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: remainingValues, storedType: vh.storedType, equals: vh.equals}
}

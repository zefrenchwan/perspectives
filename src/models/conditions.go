package models

import (
	"slices"
	"strings"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// LocalCondition is a boolean condition to accept a single entity
type LocalCondition interface {
	// Matches returns true if the entity matches the condition
	Matches(ModelEntity) bool
}

// ValueOperator applies to values (such as string) for conditions definition
type ValueOperator int

// ValuesEqual test if values are the same (including case)
const ValuesEqual ValueOperator = 1

// ValuesEqualIgnoreCase test if values are the same no matter the case
const ValuesEqualIgnoreCase ValueOperator = 2

// valuesOperatorRun applies the ValueOperator on source (parameter) and reference (from condition)
func valuesOperatorRun(source, reference string, operator ValueOperator) bool {
	switch operator {
	case ValuesEqual:
		return source == reference
	case ValuesEqualIgnoreCase:
		return strings.EqualFold(source, reference)
	default:
		return false
	}
}

// PeriodOperator is the type to define a condition on periods
type PeriodOperator int

// NonDisjoinPeriods tests if periods have a common point
const NonDisjoinPeriods PeriodOperator = 1

// SamePeriods tests if periods are equals
const SamePeriods PeriodOperator = 2

// AcceptsAllOperator returns true no matter the condition
const AcceptsAllOperator PeriodOperator = 3

// periodOperatorRun returns true if operator applied to current and reference matches
func periodOperatorRun(current, reference structures.Period, operator PeriodOperator) bool {
	switch operator {
	case AcceptsAllOperator:
		return true
	case NonDisjoinPeriods:
		return !current.Intersection(reference).IsEmpty()
	case SamePeriods:
		return current.Equals(reference)
	default:
		return false
	}
}

// LocalTypeCondition accepts entities if they match a type.
// For instance, a condition may be "only links" or "links or objects".
type LocalTypeCondition struct {
	// values are the types to accept.
	values []EntityType
}

// NewTypeCondition returns a condition that accepts only a set of types.
func NewTypeCondition(types ...EntityType) LocalCondition {
	return LocalTypeCondition{values: types}
}

// Matches for a LocalTypeCondition is true if the entity's type is in the list of accepted types.
// Note that nil value is rejected (because we cannot tell its type for sure)
func (l LocalTypeCondition) Matches(e ModelEntity) bool {
	// No type accepted => false
	if len(l.values) == 0 {
		return false
	} else if e == nil {
		// Nil => false
		return false
	} else {
		return slices.Contains(l.values, e.GetType())
	}
}

// LocalMatchingAttributeCondition is a condition for an object attribute on a given period.
// For instance, nationality = "french" during a period
type LocalMatchingAttributeCondition struct {
	AttributeName     string            // Name of the attribute to find in the object
	AttributeValue    string            // Value of the attribute to compare with
	AttributeOperator ValueOperator     // Operator for value (such as equals)
	ReferencePeriod   structures.Period // Period to match attribute during
	PeriodOoperator   PeriodOperator    // Operator for the period (such as with at least a common point)
}

// Matches returns true if all of those conditions apply :
// The parameter is indeed an object and has that attribute
// The condition on attribute compared to value matches
// The period of matching is acceptable regarding the period condition
func (l LocalMatchingAttributeCondition) Matches(e ModelEntity) bool {
	if e == nil {
		return false
	} else if e.GetType() != EntityTypeObject {
		return false
	}

	object, _ := e.AsObject()
	if matchingValues, found := object.GetValue(l.AttributeName); !found {
		return false
	} else {
		for value, period := range matchingValues {
			if valuesOperatorRun(value, l.AttributeValue, l.AttributeOperator) {
				if periodOperatorRun(period, l.ReferencePeriod, l.PeriodOoperator) {
					return true
				}
			}
		}
	}

	return false
}

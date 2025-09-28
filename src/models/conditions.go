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

// LocalTemporalCondition accepts a temporal entity for a given period.
// By convention, if period is empty, then entity is refused.
// For instance, a condition may be to be active since a given period.
type LocalTemporalCondition interface {
	// MatchesDuring may accept an entity for a given period
	MatchesDuring(TemporalEntity) structures.Period
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

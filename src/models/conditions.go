package models

import (
	"regexp"
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

// ValuesEqualIgnoreCase tests if values are the same no matter the case
const ValuesEqualIgnoreCase ValueOperator = 2

// ValueMatchesRegexp tests if a value matches a given regexp as a reference
const ValueMatchesRegexp ValueOperator = 3

// valuesOperatorRun applies the ValueOperator on source (parameter) and reference (from condition)
func valuesOperatorRun(source, reference string, operator ValueOperator) bool {
	switch operator {
	case ValuesEqual:
		return source == reference
	case ValuesEqualIgnoreCase:
		return strings.EqualFold(source, reference)
	case ValueMatchesRegexp:
		if validator := regexp.MustCompile(reference); validator == nil {
			return false
		} else {
			return validator.MatchString(source)
		}
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

// AcceptsAllPeriods returns true no matter the condition
const AcceptsAllPeriods PeriodOperator = 3

// IsIncludedInReferencePeriod tests if current period is in reference period
const IsIncludedInReferencePeriod PeriodOperator = 4

// periodOperatorRun returns true if operator applied to current and reference matches
func periodOperatorRun(current, reference structures.Period, operator PeriodOperator) bool {
	switch operator {
	case IsIncludedInReferencePeriod:
		return reference.Intersection(current).Equals(reference)
	case AcceptsAllPeriods:
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

// LocalActiveCondition means that a temporal entity is active during that reference period
type LocalActiveCondition struct {
	ReferencePeriod structures.Period
}

// Matches returns true if m is a temporal entity with a period included in reference period
func (l LocalActiveCondition) Matches(m ModelEntity) bool {
	if m == nil {
		return false
	} else if t, ok := m.(TemporalEntity); !ok {
		return false
	} else {
		return periodOperatorRun(t.ActivePeriod(), l.ReferencePeriod, IsIncludedInReferencePeriod)
	}
}

// NewActiveCondition builds a condition to be active during a provided period.
// Condition is true if object has an active period that contains reference (or is equals to)
func NewActiveCondition(reference structures.Period) LocalCondition {
	return LocalActiveCondition{ReferencePeriod: reference}
}

// LocalMatchingAttributeCondition is a condition for an object attribute on a given period.
// For instance, nationality = "french" during a period
type LocalMatchingAttributeCondition struct {
	AttributeName     string            // Name of the attribute to find in the object
	AttributeValue    string            // Value of the attribute to compare with
	AttributeOperator ValueOperator     // Operator for value (such as equals)
	ReferencePeriod   structures.Period // Period to match attribute during
	PeriodOperator    PeriodOperator    // Operator for the period (such as with at least a common point)
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
				if periodOperatorRun(period, l.ReferencePeriod, l.PeriodOperator) {
					return true
				}
			}
		}
	}

	return false
}

// NewAttributeValueCondition builds a time independant condition on attribute and value.
// For instance, attribute = value no matter the time.
func NewAttributeValueCondition(attribute, value string, operator ValueOperator) LocalCondition {
	return LocalMatchingAttributeCondition{
		AttributeName:     attribute,
		AttributeValue:    value,
		AttributeOperator: operator,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOperator:    AcceptsAllPeriods,
	}
}

// NewAttributeRegexpCondition builds a condition for an attribute value matching a regexp no matter the period.
// The regexp definition may not be valid, so we return an error if so.
func NewAttributeRegexpCondition(attribute, regexpDefinition string) (LocalCondition, error) {
	if _, err := regexp.Compile(regexpDefinition); err != nil {
		return nil, err
	}

	return LocalMatchingAttributeCondition{
		AttributeName:     attribute,
		AttributeValue:    regexpDefinition,
		AttributeOperator: ValueMatchesRegexp,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOperator:    AcceptsAllPeriods,
	}, nil
}

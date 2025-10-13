package commons

import (
	"regexp"
	"strings"
)

// StringOperator is a binary operator to apply on a value compared to a reference.
// Value is the first operand, Reference is the second operand.
type StringOperator int

// StringEquals tests if value equals reference
const StringEquals StringOperator = 1

// StringEqualsIgnoreCase tests if value equals reference no matter the case
const StringEqualsIgnoreCase StringOperator = 2

// StringContains tests if REFERENCE contains value
const StringContains StringOperator = 3

// StringMatchesRegexp tests if value matches reference as a regexp
const StringMatchesRegexp StringOperator = 4

// Accepts returns true if operator applied to (base, reference) returns true
func (o StringOperator) Accepts(base, reference string) bool {
	switch o {
	case StringEquals:
		return base == reference
	case StringEqualsIgnoreCase:
		return strings.EqualFold(base, reference)
	case StringContains:
		return strings.Contains(reference, base)
	case StringMatchesRegexp:
		if expression, err := regexp.Compile(reference); err != nil {
			return false
		} else if expression == nil {
			return false
		} else {
			return expression.MatchString(base)
		}
	default:
		return false
	}
}

// TemporalOperator defines binay operator working on a period compared to a reference period.
// Operators do NOT commute in general.
// It means that first operand HAS TO BE the current period to test whereas second operand is the REFERENCE period.
type TemporalOperator int

// TemporalEquals tests if current equals reference
const TemporalEquals TemporalOperator = 1

// TemporalCommonPoint tests if current and reference have at least a common point
const TemporalCommonPoint TemporalOperator = 2

// TemporalAlwaysAccept always accepts no matter current period
const TemporalAlwaysAccept TemporalOperator = 3

// TemporalAlwaysRefuse always refuses no matter current period
const TemporalAlwaysRefuse TemporalOperator = 4

// TemporalReferenceContains tests if current is included in reference
const TemporalReferenceContains TemporalOperator = 5

// Accepts executes the operator on current and reference (in that order)
func (t TemporalOperator) Accepts(current Period, reference Period) bool {
	switch t {
	case TemporalAlwaysAccept:
		return true
	case TemporalAlwaysRefuse:
		return false
	case TemporalCommonPoint:
		return !current.Intersection(reference).IsEmpty()
	case TemporalEquals:
		return current.Equals(reference)
	case TemporalReferenceContains:
		return current.IsIncludedIn(reference)
	default:
		return false
	}
}

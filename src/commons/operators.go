package commons

import (
	"regexp"
	"strings"
)

// LocalOperand defines types an operator applies to.
// Its purpose is to restrict types an operator applies to.
// So far, it is any, letting the ability to add later further restrictions
type LocalOperand any

// LocalOperator defines an operation between an operand and its reference
type LocalOperator[T LocalOperand] interface {
	// Accepts returns true if operator applied to operand and reference returns true
	Accepts(operand, reference T) bool
}

// IntOperator defines operations on ints.
// They apply to a value as first operand and reference as second operand
type IntOperator uint8

// IntEquals accepts when base == reference
const IntEquals IntOperator = 0x1

// IntNotEquals accepts when base == reference
const IntNotEquals IntOperator = 0x2

// IntStrictLess accepts when base < reference
const IntStrictLess IntOperator = 0x3

// IntStrictGreater accepts when base > reference
const IntStrictGreater IntOperator = 0x4

// IntLessOrEquals accepts when base <= reference
const IntLessOrEquals IntOperator = 0x5

// IntGreaterOrEquals accepts when base >= reference
const IntGreaterOrEquals IntOperator = 0x6

// Accepts applies the operator to operand and reference (in that order)
func (i IntOperator) Accepts(operand, reference int) bool {
	switch i {
	case IntEquals:
		return operand == reference
	case IntNotEquals:
		return operand != reference
	case IntStrictLess:
		return operand < reference
	case IntLessOrEquals:
		return operand <= reference
	case IntStrictGreater:
		return operand > reference
	case IntGreaterOrEquals:
		return operand >= reference
	default:
		return false
	}
}

// FloatOperator defines operators on float64
type FloatOperator uint8

// FloatEquals uses == (no epsilon check)
const FloatEquals FloatOperator = 0x1

// FloatNotEquals uses != (no epsilon check)
const FloatNotEquals FloatOperator = 0x2

// FloatStrictLess applies < operator
const FloatStrictLess FloatOperator = 0x3

// FloatStrictGreater applies > operator
const FloatStrictGreater FloatOperator = 0x4

// FloatLessOrEquals applies <= operator
const FloatLessOrEquals FloatOperator = 0x5

// FloatGreaterOrEquals applies >= operator
const FloatGreaterOrEquals FloatOperator = 0x6

// Accepts applies the operator to that operand and reference (in that order)
func (f FloatOperator) Accepts(operand, reference float64) bool {
	switch f {
	case FloatEquals:
		return operand == reference
	case FloatNotEquals:
		return operand != reference
	case FloatStrictLess:
		return operand < reference
	case FloatLessOrEquals:
		return operand <= reference
	case FloatStrictGreater:
		return operand > reference
	case FloatGreaterOrEquals:
		return operand >= reference
	default:
		return false
	}
}

// StringOperator is a binary operator to apply on a value compared to a reference.
// Value is the first operand, Reference is the second operand.
type StringOperator uint8

// StringEquals tests if value equals reference
const StringEquals StringOperator = 0x1

// StringEqualsIgnoreCase tests if value equals reference no matter the case
const StringEqualsIgnoreCase StringOperator = 0x2

// StringContains tests if REFERENCE contains value
const StringContains StringOperator = 0x3

// StringMatchesRegexp tests if value matches reference as a regexp
const StringMatchesRegexp StringOperator = 0x4

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
type TemporalOperator uint8

// TemporalEquals tests if current equals reference
const TemporalEquals TemporalOperator = 0x1

// TemporalCommonPoint tests if current and reference have at least a common point
const TemporalCommonPoint TemporalOperator = 0x2

// TemporalAlwaysAccept always accepts no matter current period
const TemporalAlwaysAccept TemporalOperator = 0x3

// TemporalAlwaysRefuse always refuses no matter current period
const TemporalAlwaysRefuse TemporalOperator = 0x4

// TemporalReferenceContains tests if current is included in reference
const TemporalReferenceContains TemporalOperator = 0x5

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

// LocalSetOperator applies operand to a set of values (slice as a set)
type LocalSetOperator[T LocalOperand] interface {
	// Accepts returns true if operator applied to operand and reference returns true
	Accepts(operand T, reference []T) bool
}

// LocalSetMapper defines set operators from local operators.
// For instance, given a set of string ["a","b","c"],
// the mapper MatchesOneInSetOperator applied to equals would
// accept "b" because "b"  matches one string in the set for equals.
type LocalSetMapper uint8

// MatchesOneInSetOperator accepts if at least one value matches within the reference set
const MatchesOneInSetOperator LocalSetMapper = 0x1

// MatchesAllInSetOperator accepts if value matches all elements within the reference set
const MatchesAllInSetOperator LocalSetMapper = 0x2

// MatchesNoneInSetOperator refuses if value matches one element within the reference set, accepts otherwise
const MatchesNoneInSetOperator LocalSetMapper = 0x3

// localOperatorMapper links a mapper to a local operator
type localOperatorMapper[T LocalOperand] struct {
	// sliceOperator is the loop end condition: accepts once, for no match, for all matches ?
	sliceOperator LocalSetMapper
	// operator is the operator to apply to each element within the reference set
	operator LocalOperator[T]
}

// Accepts runs through the reference set and applies local operator to accept
func (l localOperatorMapper[T]) Accepts(operand T, reference []T) bool {
	switch l.sliceOperator {
	case MatchesAllInSetOperator:
		for _, value := range reference {
			if !l.operator.Accepts(operand, value) {
				return false
			}
		}

		return true
	case MatchesNoneInSetOperator:
		for _, value := range reference {
			if l.operator.Accepts(operand, value) {
				return false
			}
		}

		return true
	case MatchesOneInSetOperator:
		for _, value := range reference {
			if l.operator.Accepts(operand, value) {
				return true
			}
		}

		return false
	default:
		return false
	}
}

// NewLocalSetOperator returns a local set operator built with
// a set operator (iteration method)
// and local operator (applied to each element)
func NewLocalSetOperator[T LocalOperand](setOperator LocalSetMapper, operator LocalOperator[T]) LocalSetOperator[T] {
	if operator == nil {
		return nil
	}

	return localOperatorMapper[T]{sliceOperator: setOperator, operator: operator}
}

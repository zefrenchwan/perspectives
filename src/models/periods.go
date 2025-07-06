package models

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"
)

// TIME_FORMAT defines how to serialize and deserialize time data
const TIME_FORMAT = time.RFC3339

// TIME_PRECISION is the accepted thresold to define when two times are the same
const TIME_PRECISION = time.Second

const INTERVAL_BOUNDARY_LEFT = "]"
const INTERVAL_BOUNDARY_RIGHT = "["
const INTERVAL_PARTS_SEPARATOR = ";"
const INTERVAL_EMPTY = "]["
const INTERVAL_VALUE_LEFT_INFINITY = "-oo"
const INTERVAL_VALUE_RIGHT_INFINITY = "+oo"

// interval is an interval of time
type interval struct {
	// empty is true for empty sets, overrides every other info
	empty bool
	// leftFinite is true if left boundary is finite, valid only for non empty
	leftFinite bool
	//  leftIncluded is true if left boundary is included in the interval, valid only for non empty
	leftIncluded bool
	// rightFinite is true when right bound is not infinite, false otherwise
	rightFinite bool
	// rightIncluded is true if right value is included in the interval
	rightIncluded bool
	// left finite border
	leftMoment time.Time
	// right finite border
	rightMoment time.Time
}

// intervalEquals tests if two periods are the same
func intervalEquals(a, b interval) bool {
	if a.empty != b.empty {
		return false
	} else if a.empty {
		return true
	} else if a.leftFinite != b.leftFinite {
		return false
	} else if a.rightFinite != b.rightFinite {
		return false
	}

	if a.leftFinite {
		if a.leftIncluded != b.leftIncluded {
			return false
		} else if a.leftMoment.Compare(b.leftMoment) != 0 {
			return false
		}
	}

	if a.rightFinite {
		if a.rightIncluded != b.rightIncluded {
			return false
		} else if !a.rightMoment.Equal(b.rightMoment) {
			return false
		}
	}

	return true
}

// isFull returns true if interval is the full space
func (i interval) isFull() bool {
	return !i.empty && !i.leftFinite && !i.rightFinite
}

// compareByLexicographicOrder tests if intervals are less or equa
func (i interval) compareByLexicographicOrder(other interval) int {
	if i.empty && other.empty {
		return 0
	} else if i.empty {
		return -1
	} else if other.empty {
		return 1
	}

	left, right := compareIntervalsBounds(i, other)
	switch {
	case left == relativePositionComparison(left_position_strict):
		return -1
	case left == relativePositionComparison(left_position_same_value_different_bounds):
		return -1
	case left == relativePositionComparison(right_position_same_value_different_bounds):
		return 1
	case left == relativePositionComparison(right_position_same_value_different_bounds):
		return 1
	case right == relativePositionComparison(left_position_strict):
		return -1
	case right == relativePositionComparison(left_position_same_value_different_bounds):
		return -1
	case right == relativePositionComparison(right_position_same_value_different_bounds):
		return 1
	case right == relativePositionComparison(right_position_same_value_different_bounds):
		return 1
	default:
		return 0
	}
}

// intervalFromString parses an interval of time
func intervalFromString(rawData string) (interval, error) {
	if rawData == INTERVAL_EMPTY {
		return interval{empty: true}, nil
	}

	var empty interval
	parts := strings.Split(rawData, INTERVAL_PARTS_SEPARATOR)
	letters := strings.Split(rawData, "")
	size := len(rawData)

	if len(parts) != 2 {
		return empty, errors.New("malformed interval around separator")
	} else if len(parts[0]) <= 2 {
		return empty, errors.New("invalid left part, insufficient size")
	} else if len(parts[1]) <= 2 {
		return empty, errors.New("invalid right part, insufficient size")
	}

	leftBound, rightBound := letters[0], letters[size-1]
	leftPart, _ := strings.CutPrefix(parts[0], leftBound)
	rightPart, _ := strings.CutSuffix(parts[1], rightBound)

	// parse parts and raise error if intervals are malformed
	if leftBound != INTERVAL_BOUNDARY_LEFT && leftBound != INTERVAL_BOUNDARY_RIGHT {
		return empty, errors.New("invalid interval boundaries")
	} else if rightBound != INTERVAL_BOUNDARY_LEFT && rightBound != INTERVAL_BOUNDARY_RIGHT {
		return empty, errors.New("invalid interval boundaries")
	} else if leftPart == INTERVAL_VALUE_RIGHT_INFINITY {
		return empty, errors.New("invalid infinite left part")
	} else if rightPart == INTERVAL_VALUE_LEFT_INFINITY {
		return empty, errors.New("invalid infinite right part")
	} else if leftPart == INTERVAL_VALUE_LEFT_INFINITY && leftBound != INTERVAL_BOUNDARY_LEFT {
		return empty, errors.New("invalid infinite left boundary")
	} else if rightPart == INTERVAL_VALUE_RIGHT_INFINITY && rightBound != INTERVAL_BOUNDARY_RIGHT {
		return empty, errors.New("invalid infinite right boundary")
	}

	leftInfinite, rightInfinite := leftPart == INTERVAL_VALUE_LEFT_INFINITY, rightPart == INTERVAL_VALUE_RIGHT_INFINITY
	if leftInfinite && rightInfinite {
		return interval{empty: false, leftFinite: false, rightFinite: false}, nil
	}

	leftIn, rightIn := leftBound == INTERVAL_BOUNDARY_RIGHT, rightBound == INTERVAL_BOUNDARY_LEFT
	var leftVal, rightVal time.Time

	// parse values if any
	if !leftInfinite {
		value, errLV := time.Parse(TIME_FORMAT, leftPart)
		if errLV != nil {
			return empty, errLV
		} else {
			leftVal = value
		}
	}

	if !rightInfinite {
		value, errRV := time.Parse(TIME_FORMAT, rightPart)
		if errRV != nil {
			return empty, errRV
		} else {
			rightVal = value
		}
	}

	// and (finally) make the interval
	if leftInfinite {
		return interval{empty: false, leftFinite: false, rightFinite: true, rightIncluded: rightIn, rightMoment: rightVal}, nil
	} else if rightInfinite {
		return interval{empty: false, leftFinite: true, rightFinite: false, leftIncluded: leftIn, leftMoment: leftVal}, nil
	}

	comparison := leftVal.Compare(rightVal)
	if comparison > 0 {
		return empty, errors.New("min value is more than max value")
	} else if comparison == 0 && (!leftIn || !rightIn) {
		return empty, errors.New("min value equals max value but boundaries are not included")
	}

	// finite interval build
	return interval{empty: false, leftFinite: true, rightFinite: true, leftIncluded: leftIn, rightIncluded: rightIn, leftMoment: leftVal, rightMoment: rightVal}, nil
}

// toString returns the interval as a string
func (i interval) toString() string {
	var result string
	if i.empty {
		return INTERVAL_EMPTY
	}

	if i.leftFinite {
		if i.leftIncluded {
			result = INTERVAL_BOUNDARY_RIGHT
		} else {
			result = INTERVAL_BOUNDARY_LEFT
		}

		result = result + i.leftMoment.Format(TIME_FORMAT)
	} else {
		result = INTERVAL_BOUNDARY_LEFT + INTERVAL_VALUE_LEFT_INFINITY
	}

	result = result + INTERVAL_PARTS_SEPARATOR

	if i.rightFinite {
		result = result + i.rightMoment.Format(TIME_FORMAT)
		if i.rightIncluded {
			result = result + INTERVAL_BOUNDARY_LEFT
		} else {
			result = result + INTERVAL_BOUNDARY_RIGHT
		}
	} else {
		result = result + INTERVAL_VALUE_RIGHT_INFINITY + INTERVAL_BOUNDARY_RIGHT
	}

	return result
}

// toRawString returns value as raw data
func (i interval) toRawString() string {
	return fmt.Sprintf("Period: [ empty %t finite: %t %t included: %t %t values: %s %s ]",
		i.empty,
		i.leftFinite, i.rightFinite,
		i.leftIncluded, i.rightIncluded,
		i.leftMoment.Format(TIME_FORMAT), i.rightMoment.Format(TIME_FORMAT),
	)
}

// complement returns the complement of the interval.
// It may be a full interval, empty, one infinite interval, or two.
// i union its complements should be full interval
func (i interval) complement() []interval {
	if i.empty {
		return []interval{{empty: false, leftFinite: false, rightFinite: false}}
	} else if !i.leftFinite && !i.rightFinite {
		return []interval{{empty: true}}
	} else if !i.leftFinite {
		return []interval{{
			empty: false, leftFinite: true, rightFinite: false,
			leftIncluded: !i.rightIncluded, leftMoment: i.rightMoment},
		}
	} else if !i.rightFinite {
		return []interval{{
			empty: false, leftFinite: false, rightFinite: true,
			rightIncluded: !i.leftIncluded, rightMoment: i.leftMoment},
		}
	} else {
		return []interval{
			{
				empty: false, leftFinite: false, rightFinite: true,
				rightIncluded: !i.leftIncluded, rightMoment: i.leftMoment,
			},
			{
				empty: false, rightFinite: false, leftFinite: true,
				leftIncluded: !i.rightIncluded, leftMoment: i.rightMoment,
			},
		}
	}
}

// relativePositionComparison is a tehnical type to deal with intervals comparison
type relativePositionComparison uint8

const left_position_strict = 0x1
const left_position_same_value_different_bounds uint8 = 0x2
const equals_infinite_position uint8 = 0x3
const equals_finite_position uint8 = 0x4
const right_position_same_value_different_bounds uint8 = 0x5
const right_position_strict uint8 = 0x6
const both_values_empty uint8 = 0x7
const empty_with_value uint8 = 0x8
const value_with_empty uint8 = 0x9

// isLeft returns true if ref is lower than other
func (r relativePositionComparison) isLeft() bool {
	return r == relativePositionComparison(left_position_strict) || r == relativePositionComparison(left_position_same_value_different_bounds)
}

// isRight returns true if ref is more than other
func (r relativePositionComparison) isRight() bool {
	return r == relativePositionComparison(right_position_strict) || r == relativePositionComparison(right_position_same_value_different_bounds)
}

// isEqual returns true if ref equals than other
func (r relativePositionComparison) isEqual() bool {
	return r == relativePositionComparison(equals_finite_position) || r == relativePositionComparison(equals_infinite_position)
}

// compareIntervalsBounds compares two intervals boundaries.
// It returns two values: relative value of left and right values of the interval
func compareIntervalsBounds(ref, other interval) (relativePositionComparison, relativePositionComparison) {
	if ref.empty && other.empty {
		return relativePositionComparison(both_values_empty), relativePositionComparison(both_values_empty)
	} else if ref.empty {
		return relativePositionComparison(empty_with_value), relativePositionComparison(empty_with_value)
	} else if other.empty {
		return relativePositionComparison(value_with_empty), relativePositionComparison(value_with_empty)
	}

	var leftResult relativePositionComparison
	var rightResult relativePositionComparison

	if !ref.leftFinite && !other.leftFinite {
		leftResult = relativePositionComparison(equals_infinite_position)
	} else if ref.leftFinite && !other.leftFinite {
		leftResult = relativePositionComparison(right_position_strict)
	} else if !ref.leftFinite && other.leftFinite {
		leftResult = left_position_strict
	} else {
		// both values are finite
		leftComparison := ref.leftMoment.Compare(other.leftMoment)
		if leftComparison < 0 {
			leftResult = relativePositionComparison(left_position_strict)
		} else if leftComparison > 0 {
			leftResult = relativePositionComparison(right_position_strict)
		} else if ref.leftIncluded && !other.leftIncluded {
			leftResult = relativePositionComparison(right_position_same_value_different_bounds)
		} else if !ref.leftIncluded && other.leftIncluded {
			leftResult = relativePositionComparison(left_position_same_value_different_bounds)
		} else {
			leftResult = relativePositionComparison(equals_finite_position)
		}
	}

	// same on the right side
	if !ref.rightFinite && !other.rightFinite {
		rightResult = relativePositionComparison(equals_infinite_position)
	} else if !ref.rightFinite {
		rightResult = relativePositionComparison(right_position_strict)
	} else if !other.rightFinite {
		rightResult = relativePositionComparison(left_position_strict)
	} else {
		rightComparison := ref.leftMoment.Compare(other.rightMoment)
		if rightComparison < 0 {
			rightResult = relativePositionComparison(left_position_strict)
		} else if rightComparison > 0 {
			rightResult = relativePositionComparison(right_position_strict)
		} else if ref.rightIncluded && !other.rightIncluded {
			rightResult = relativePositionComparison(right_position_same_value_different_bounds)
		} else if !ref.rightIncluded && other.rightIncluded {
			rightResult = relativePositionComparison(left_position_same_value_different_bounds)
		} else {
			rightResult = relativePositionComparison(equals_finite_position)
		}
	}

	return leftResult, rightResult
}

// intervalsIntersection returns the intersection of all parameters
func intervalsIntersection(intervals []interval) interval {
	var remaining []interval
	var empty bool
	for _, interval := range intervals {
		if !interval.empty {
			remaining = append(remaining, interval)
		} else {
			empty = true
		}
	}

	if len(remaining) == 0 || empty {
		return interval{empty: true}
	}

	var intersection interval
	for index, value := range remaining {
		if index == 0 {
			intersection = value
			continue
		}

		// calculate the actual intersection between var intersection and value
		leftC, rightC := compareIntervalsBounds(intersection, value)
		switch leftC {
		case relativePositionComparison(left_position_strict), relativePositionComparison(left_position_same_value_different_bounds):
			intersection.leftFinite = value.leftFinite
			intersection.leftIncluded = value.leftIncluded
			intersection.leftMoment = value.leftMoment
		}
		switch rightC {
		case relativePositionComparison(right_position_strict), relativePositionComparison(right_position_same_value_different_bounds):
			intersection.rightFinite = value.rightFinite
			intersection.rightIncluded = value.rightIncluded
			intersection.rightMoment = value.rightMoment
		}
	}

	// then, test if intersection is empty or not
	if intersection.leftFinite && intersection.rightFinite {
		comparison := intersection.leftMoment.Compare(intersection.rightMoment)
		switch {
		case comparison == 0 && !(intersection.leftIncluded && intersection.rightIncluded):
			return interval{empty: true}
		case comparison > 0:
			return interval{empty: true}
		}
	}

	return intersection
}

// Period is a given period of time.
// Formally, a finite union of intervals
type Period struct {
	intervals []interval
}

// NewFullPeriod returns a period equivalent to ]-oo,+oo[
func NewFullPeriod() Period {
	value := interval{empty: false, leftFinite: false, rightFinite: false}
	return Period{intervals: []interval{value}}
}

// NewEmptyPeriod builds an empty period
func NewEmptyPeriod() Period {
	return Period{}
}

// NewFinitePeriod builds a period equivalent to a new finite interval (min, max)
// SPECIAL CASES: it may return an empty period according to mathematical definition
func NewFinitePeriod(min, max time.Time, minIncluded, maxIncluded bool) Period {
	comparison := min.Compare(max)
	if comparison == 0 && !(minIncluded && maxIncluded) {
		return Period{}
	} else if comparison > 0 {
		return Period{}
	}

	content := interval{
		empty:         false,
		leftFinite:    true,
		rightFinite:   true,
		leftIncluded:  minIncluded,
		rightIncluded: maxIncluded,
		leftMoment:    min.Truncate(TIME_PRECISION),
		rightMoment:   max.Truncate(TIME_PRECISION),
	}

	return Period{intervals: []interval{content}}
}

// NewPeriodSince builds a period equivalent to (leftLimit, +oo[
func NewPeriodSince(leftLimit time.Time, leftIn bool) Period {
	content := interval{
		empty:        false,
		rightFinite:  false,
		leftFinite:   true,
		leftIncluded: leftIn,
		leftMoment:   leftLimit.Truncate(TIME_PRECISION),
	}

	return Period{intervals: []interval{content}}
}

// NewPeriodUntil builds a period equivalent to ]-oo,rightLimit)
func NewPeriodUntil(rightLimit time.Time, rightIn bool) Period {
	content := interval{
		empty:         false,
		leftFinite:    false,
		rightFinite:   true,
		rightIncluded: rightIn,
		rightMoment:   rightLimit.Truncate(TIME_PRECISION),
	}

	return Period{intervals: []interval{content}}
}

// Intersection returns the set intersection between p and other as intervals
func (p Period) Intersection(other Period) Period {
	if len(p.intervals) == 0 || len(other.intervals) == 0 {
		return Period{}
	}

	var result []interval
	for _, sourceInterval := range p.intervals {
		for _, otherInterval := range other.intervals {
			value := intervalsIntersection([]interval{sourceInterval, otherInterval})
			if !value.empty {
				result = append(result, value)
			}
		}
	}

	return Period{result}
}

// Equals returns true if periods have the same content
func (p Period) Equals(other Period) bool {
	if len(p.intervals) != len(other.intervals) {
		return false
	}

	for _, value := range p.intervals {
		if !slices.ContainsFunc(other.intervals, func(a interval) bool { return intervalEquals(a, value) }) {
			return false
		}
	}

	return true
}

// AsRawString returns the period as a string, concatenation of underlying intervals
func (p Period) AsRawString() string {
	var values []string
	for _, val := range p.intervals {
		values = append(values, val.toString())
	}

	sort.Strings(values)
	return "Period [" + strings.Join(values, ",") + "]"
}

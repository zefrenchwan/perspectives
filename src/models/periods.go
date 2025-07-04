package models

import (
	"errors"
	"strings"
	"time"
)

const TIME_FORMAT = time.RFC3339
const INTERVAL_BOUNDARY_LEFT = "]"
const INTERVAL_BOUNDARY_RIGHT = "["
const INTERVAL_PARTS_SEPARATOR = ";"
const INTERVAL_EMPTY = "]["
const INTERVAL_VALUE_LEFT_INFINITY = "-oo"
const INTERVAL_VALUE_RIGHT_INFINITY = "+oo"

// Period is an interval of time
type Period struct {
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

// PeriodFromString parses an interval of time
func PeriodFromString(interval string) (Period, error) {
	if interval == INTERVAL_EMPTY {
		return Period{empty: true}, nil
	}

	var empty Period
	parts := strings.Split(interval, INTERVAL_PARTS_SEPARATOR)
	letters := strings.Split(interval, "")
	size := len(interval)

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
		return Period{empty: false, leftFinite: false, rightFinite: false}, nil
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
		return Period{empty: false, leftFinite: false, rightFinite: true, rightIncluded: rightIn, rightMoment: rightVal}, nil
	} else if rightInfinite {
		return Period{empty: false, leftFinite: true, rightFinite: false, leftIncluded: leftIn, leftMoment: leftVal}, nil
	}

	comparison := leftVal.UTC().Compare(rightVal.UTC())
	if comparison > 0 {
		return empty, errors.New("min value is more than max value")
	} else if comparison == 0 && (!leftIn || !rightIn) {
		return empty, errors.New("min value equals max value but boundaries are not included")
	}

	// finite interval build
	return Period{empty: false, leftFinite: true, rightFinite: true, leftIncluded: leftIn, rightIncluded: rightIn, leftMoment: leftVal, rightMoment: rightVal}, nil
}

// AsString returns the interval as a string
func (p *Period) AsString() string {
	var result string
	if p.empty {
		return INTERVAL_EMPTY
	}

	if p.leftFinite {
		if p.leftIncluded {
			result = INTERVAL_BOUNDARY_RIGHT
		} else {
			result = INTERVAL_BOUNDARY_LEFT
		}

		result = result + p.leftMoment.Format(TIME_FORMAT)
	} else {
		result = INTERVAL_BOUNDARY_LEFT + INTERVAL_VALUE_LEFT_INFINITY
	}

	result = result + INTERVAL_PARTS_SEPARATOR

	if p.rightFinite {
		result = result + p.rightMoment.Format(TIME_FORMAT)
		if p.rightIncluded {
			result = result + INTERVAL_BOUNDARY_LEFT
		} else {
			result = result + INTERVAL_BOUNDARY_RIGHT
		}
	} else {
		result = result + INTERVAL_VALUE_RIGHT_INFINITY + INTERVAL_BOUNDARY_RIGHT
	}

	return result
}

// RelativePositionComparison is a tehnical type to deal with intervals comparison
type RelativePositionComparison uint8

const LEFT_POSITION_STRICT = 0x1
const LEFT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS uint8 = 0x2
const EQUALS_INFINITE_POSITION uint8 = 0x3
const EQUALS_FINITE_POSITION uint8 = 0x4
const RIGHT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS uint8 = 0x5
const RIGHT_POSITION_STRICT uint8 = 0x6
const BOTH_VALUES_EMPTY uint8 = 0x7
const EMPTY_WITH_VALUE uint8 = 0x8
const VALUE_WITH_EMPTY uint8 = 0x9

// comparePeriodsBounds compares two intervals boundaries.
// It returns two values: relative value of left and right values of the interval
func comparePeriodsBounds(ref, other Period) (RelativePositionComparison, RelativePositionComparison) {
	if ref.empty && other.empty {
		return RelativePositionComparison(BOTH_VALUES_EMPTY), RelativePositionComparison(BOTH_VALUES_EMPTY)
	} else if ref.empty {
		return RelativePositionComparison(EMPTY_WITH_VALUE), RelativePositionComparison(EMPTY_WITH_VALUE)
	} else if other.empty {
		return RelativePositionComparison(VALUE_WITH_EMPTY), RelativePositionComparison(VALUE_WITH_EMPTY)
	}

	var leftResult RelativePositionComparison
	var rightResult RelativePositionComparison

	if !ref.leftFinite && !other.leftFinite {
		leftResult = RelativePositionComparison(EQUALS_INFINITE_POSITION)
	} else if ref.leftFinite && !other.leftFinite {
		leftResult = RelativePositionComparison(RIGHT_POSITION_STRICT)
	} else if !ref.leftFinite && other.leftFinite {
		leftResult = LEFT_POSITION_STRICT
	} else {
		// both values are finite
		leftComparison := ref.leftMoment.UTC().Compare(other.leftMoment.UTC())
		if leftComparison < 0 {
			leftResult = RelativePositionComparison(LEFT_POSITION_STRICT)
		} else if leftComparison > 0 {
			leftResult = RelativePositionComparison(RIGHT_POSITION_STRICT)
		} else if ref.leftIncluded && !other.leftIncluded {
			leftResult = RelativePositionComparison(RIGHT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS)
		} else if !ref.leftIncluded && other.leftIncluded {
			leftResult = RelativePositionComparison(LEFT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS)
		} else {
			leftResult = RelativePositionComparison(EQUALS_FINITE_POSITION)
		}
	}

	// same on the right side
	if !ref.rightFinite && !other.rightFinite {
		rightResult = RelativePositionComparison(EQUALS_INFINITE_POSITION)
	} else if !ref.rightFinite {
		rightResult = RelativePositionComparison(RIGHT_POSITION_STRICT)
	} else if !other.rightFinite {
		rightResult = RelativePositionComparison(LEFT_POSITION_STRICT)
	} else {
		rightComparison := ref.leftMoment.UTC().Compare(other.rightMoment.UTC())
		if rightComparison < 0 {
			rightResult = RelativePositionComparison(LEFT_POSITION_STRICT)
		} else if rightComparison > 0 {
			rightResult = RelativePositionComparison(RIGHT_POSITION_STRICT)
		} else if ref.rightIncluded && !other.rightIncluded {
			rightResult = RelativePositionComparison(RIGHT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS)
		} else if !ref.rightIncluded && other.rightIncluded {
			rightResult = RelativePositionComparison(LEFT_POSITION_SAME_VALUE_DIFFERENT_BOUNDS)
		} else {
			rightResult = RelativePositionComparison(EQUALS_FINITE_POSITION)
		}
	}

	return leftResult, rightResult
}

// NewFullPeriod returns ]-oo,+oo[
func NewFullPeriod() Period {
	return Period{empty: false, leftFinite: false, rightFinite: false}
}

// NewEmptyPeriod builds an empty period
func NewEmptyPeriod() Period {
	return Period{empty: true}
}

// NewFinitePeriod builds a new finite interval
// SPECIAL CASES: it may return an empty interval according to mathematical definition
func NewFinitePeriod(min, max time.Time, minIncluded, maxIncluded bool) Period {
	comparison := min.UTC().Compare(max.UTC())
	if comparison == 0 && !(minIncluded && maxIncluded) {
		return Period{empty: true}
	} else if comparison > 0 {
		return Period{empty: true}
	}

	return Period{
		empty:         false,
		leftFinite:    true,
		rightFinite:   true,
		leftIncluded:  minIncluded,
		rightIncluded: maxIncluded,
		leftMoment:    min,
		rightMoment:   max,
	}
}

// NewPeriodSince builds (leftLimit, +oo[
func NewPeriodSince(leftLimit time.Time, leftIn bool) Period {
	return Period{
		empty:        false,
		rightFinite:  false,
		leftFinite:   true,
		leftIncluded: leftIn,
		leftMoment:   leftLimit,
	}
}

// NewPeriodUntil builds ]-oo,rightLimit)
func NewPeriodUntil(rightLimit time.Time, rightIn bool) Period {
	return Period{
		empty:         false,
		leftFinite:    false,
		rightFinite:   true,
		rightIncluded: rightIn,
		rightMoment:   rightLimit,
	}
}

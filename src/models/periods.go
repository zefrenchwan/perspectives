package models

import "time"

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

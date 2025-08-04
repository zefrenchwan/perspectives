package models

import (
	"errors"
	"slices"
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

// interval is an interval of time.
// It uses:
// empty that overrides the rest
// left and right boundaries represented by: finite or not, included or not (if finite), bounds as time.Time (if finite )
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

// buildInterval returns an interval built from values (may be empty if values lead to it)
func buildInterval(empty, minFinite, maxFinite bool, min, max time.Time, minIn, maxIn bool) interval {
	if empty {
		return interval{empty: true}
	}

	var left, right time.Time
	var leftIn, rightIn bool
	if minFinite {
		left = min.Truncate(TIME_PRECISION)
		leftIn = minIn
	} else {
		leftIn = false
	}

	if maxFinite {
		right = max.Truncate(TIME_PRECISION)
		rightIn = maxIn
	} else {
		rightIn = false
	}

	if minFinite && maxFinite {
		comparison := left.Compare(right)
		switch {
		case comparison > 0:
			return interval{empty: true}
		case comparison == 0 && !(maxIn && minIn):
			return interval{empty: true}
		}
	}

	return interval{empty: false,
		leftFinite: minFinite, leftMoment: left, leftIncluded: leftIn,
		rightFinite: maxFinite, rightMoment: right, rightIncluded: rightIn,
	}
}

// newFullInterval builds a new interval equals to full space
func newFullInterval() interval {
	return interval{empty: false, leftFinite: false, rightFinite: false}
}

// newIntervalSince builds the interval (leftLimit, +oo[
func newIntervalSince(leftLimit time.Time, leftIn bool) interval {
	return interval{
		empty:        false,
		rightFinite:  false,
		leftFinite:   true,
		leftIncluded: leftIn,
		leftMoment:   leftLimit.Truncate(TIME_PRECISION),
	}

}

// newIntervalUntil builds the interval ]-oo, rightLimit)
func newIntervalUntil(rightLimit time.Time, rightIn bool) interval {
	return interval{
		empty:         false,
		leftFinite:    false,
		rightFinite:   true,
		rightIncluded: rightIn,
		rightMoment:   rightLimit.Truncate(TIME_PRECISION),
	}
}

// newIntervalDuring returns the interval (min,max) or empty when result is mathematically empty.
// If min > max, for instance, result is mathematically empty and so is result of the function
func newIntervalDuring(min, max time.Time, minIncluded, maxIncluded bool) interval {
	left := min.Truncate(TIME_PRECISION)
	right := max.Truncate(TIME_PRECISION)
	comparison := left.Compare(right)
	switch {
	case comparison > 0:
		return interval{empty: true}
	case comparison == 0 && !(minIncluded && maxIncluded):
		return interval{empty: true}
	default:
		return interval{
			empty:         false,
			leftFinite:    true,
			rightFinite:   true,
			leftIncluded:  minIncluded,
			rightIncluded: maxIncluded,
			leftMoment:    min.Truncate(TIME_PRECISION),
			rightMoment:   max.Truncate(TIME_PRECISION),
		}
	}
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

// intervalCompare compares two intervals by empty, then left border, then right border
func intervalCompare(a, b interval) int {
	if a.empty && b.empty {
		return 0
	} else if a.empty {
		return -1
	} else if b.empty {
		return 1
	}

	if a.leftFinite && !b.leftFinite {
		return 1
	} else if !a.leftFinite && b.leftFinite {
		return -1
	} else if a.leftFinite && b.leftFinite {
		comparison := a.leftMoment.Compare(b.leftMoment)
		if comparison != 0 {
			return comparison
		} else if a.leftIncluded != b.leftIncluded {
			if a.leftIncluded {
				return 1
			} else {
				return -1
			}
		}
	}

	// same values on the left side, do the same on the right side
	if a.rightFinite && !b.rightFinite {
		return -1
	} else if !a.rightFinite && b.rightFinite {
		return 1
	} else if a.rightFinite && b.rightFinite {
		comparison := a.rightMoment.Compare(b.rightMoment)
		if comparison != 0 {
			return comparison
		} else if a.rightIncluded == b.rightIncluded {
			return 0
		} else if a.rightIncluded {
			return 1
		} else {
			return -1
		}
	}

	return 0
}

// isFull returns true if interval is the full space
func (i interval) isFull() bool {
	return !i.empty && !i.leftFinite && !i.rightFinite
}

// contains returns true if point is in the interval (as in set theory)
func (i interval) contains(point time.Time) bool {
	if i.empty {
		return false
	}

	// assume i is (L,R), L for left value, R for right value

	// test if point is less than L and, in that case, return false for sure
	if i.leftFinite {
		comparison := i.leftMoment.Compare(point)
		switch {
		case comparison > 0:
			// L > point, ie. point < L, so false for sure
			return false
		case comparison == 0:
			// L == point, so depends if L is in or not
			return i.leftIncluded
		}
	}

	// test is point is more than R and, in that case, return false for sure
	if i.rightFinite {
		comparison := i.rightMoment.Compare(point)
		switch {
		case comparison < 0:
			// R < point, so return false for sure
			return false
		case comparison == 0:
			// R == point, so depends if R is in or not
			return i.rightIncluded
		}
	}

	return true
}

// isIncludedIn returns true if other contains i
func (i interval) isIncludedIn(other interval) bool {
	if i.empty {
		return true
	} else if other.empty {
		return false
	}

	var isIncludedLeft bool
	// i is in other if
	// (i.left bound) is more than (other.left bound)
	// AND
	// (i.right bound) is less than (other.right bound)
	if other.leftFinite && !i.leftFinite {
		// no need to go deeper
		return false
	} else if !i.leftFinite && !other.leftFinite {
		isIncludedLeft = true
	} else if !other.leftFinite && i.leftFinite {
		isIncludedLeft = true
	} else {
		comparison := i.leftMoment.Compare(other.leftMoment)
		switch {
		case comparison < 0:
			return false
		case comparison > 0:
			isIncludedLeft = true
		case comparison == 0:
			if i.leftIncluded {
				isIncludedLeft = other.leftIncluded
			} else {
				isIncludedLeft = true
			}
		}
	}

	// no need to go deeper
	if !isIncludedLeft {
		return false
	}

	// same on the right side
	if other.rightFinite && !i.rightFinite {
		return false
	} else if !other.rightFinite {
		return isIncludedLeft
	} else {
		comparison := i.rightMoment.Compare(other.rightMoment)
		switch {
		case comparison < 0:
			return isIncludedLeft
		case comparison > 0:
			return false
		case other.rightIncluded:
			return isIncludedLeft
		case !other.rightIncluded:
			return isIncludedLeft && !i.rightIncluded
		default:
			return false
		}
	}
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

		// calculate the actual intersection between intersection and value
		// Get the max value for left bounds
		if value.leftFinite {
			if !intersection.leftFinite {
				intersection.leftFinite = value.leftFinite
				intersection.leftIncluded = value.leftIncluded
				intersection.leftMoment = value.leftMoment
			} else {
				comparison := intersection.leftMoment.Compare(value.leftMoment)
				if comparison < 0 || (comparison == 0 && value.leftIncluded) {
					intersection.leftFinite = value.leftFinite
					intersection.leftIncluded = value.leftIncluded
					intersection.leftMoment = value.leftMoment
				}
			}
		}
		// Get the min value for right bounds
		if value.rightFinite {
			if !intersection.rightFinite {
				intersection.rightFinite = value.rightFinite
				intersection.rightIncluded = value.rightIncluded
				intersection.rightMoment = value.rightMoment
			} else {
				comparison := intersection.rightMoment.Compare(value.rightMoment)
				if comparison > 0 || (comparison == 0 && !value.rightIncluded) {
					intersection.rightFinite = value.rightFinite
					intersection.rightIncluded = value.rightIncluded
					intersection.rightMoment = value.rightMoment
				}
			}
		}
	}

	// then, test if intersection is empty or not.
	// Interval is built but may be empty
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

// intersection is the "object" version of intervalsIntersection
func (i interval) intersection(other interval) interval {
	return intervalsIntersection([]interval{i, other})
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

// union calculates the union of intervals
func (i interval) union(other interval) []interval {
	if i.empty || other.isFull() {
		return []interval{other}
	} else if other.empty || i.isFull() {
		return []interval{i}
	}

	var joinable bool
	var comparison int

	switch {
	case !i.leftFinite && !other.leftFinite:
		joinable = true
	case !i.rightFinite && !other.rightFinite:
		joinable = true
	case !i.leftFinite:
		comparison = i.rightMoment.Compare(other.leftMoment)
		if comparison < 0 {
			joinable = false
		} else if comparison == 0 {
			joinable = i.rightIncluded || other.leftIncluded
		} else {
			joinable = true
		}
	case !other.leftFinite:
		comparison = other.rightMoment.Compare(i.leftMoment)
		if comparison < 0 {
			joinable = false
		} else if comparison == 0 {
			joinable = other.rightIncluded || i.leftIncluded
		} else {
			joinable = true
		}
	case !i.rightFinite:
		comparison = i.leftMoment.Compare(other.rightMoment)
		if comparison > 0 {
			joinable = false
		} else if comparison == 0 {
			joinable = i.leftIncluded || other.rightIncluded
		} else {
			joinable = true
		}
	case !other.rightFinite:
		comparison = i.rightMoment.Compare(other.leftMoment)
		if comparison < 0 {
			joinable = false
		} else if comparison == 0 {
			joinable = i.leftIncluded || other.rightIncluded
		} else {
			joinable = true
		}
	default:
		comparison = i.rightMoment.Compare(other.leftMoment)
		if comparison < 0 {
			joinable = false
		} else if comparison == 0 {
			joinable = i.rightIncluded || other.leftIncluded
		} else {
			comparison = i.leftMoment.Compare(other.rightMoment)
			if comparison > 0 {
				joinable = false
			} else if comparison == 0 {
				joinable = i.leftIncluded || other.rightIncluded
			} else {
				joinable = true
			}
		}
	}

	if !joinable {
		return []interval{i, other}
	}

	// build the result getting the most extreme values
	var minFinite, maxFinite, minIncluded, maxIncluded bool
	var minValue, maxValue time.Time
	// left bound: pick the less the values
	minFinite = i.leftFinite && other.leftFinite
	if minFinite {
		comparison = i.leftMoment.Compare(other.rightMoment)
		switch {
		case comparison < 0:
			minIncluded, minValue = i.leftIncluded, i.leftMoment
		case comparison > 0:
			minIncluded, minValue = other.leftIncluded, other.leftMoment
		default:
			minIncluded, minValue = i.leftIncluded || other.leftIncluded, i.leftMoment
		}
	}
	// right bound: pick the more the values
	maxFinite = i.rightFinite && other.rightFinite
	if maxFinite {
		comparison = i.rightMoment.Compare(other.rightMoment)
		switch {
		case comparison < 0:
			maxIncluded, maxValue = other.rightIncluded, other.rightMoment
		case comparison > 0:
			maxIncluded, maxValue = i.rightIncluded, i.rightMoment
		default:
			maxIncluded, maxValue = i.rightIncluded || other.rightIncluded, i.rightMoment
		}
	}

	// and finally, return
	return []interval{{
		empty:      false,
		leftFinite: minFinite, leftIncluded: minIncluded, leftMoment: minValue,
		rightFinite: maxFinite, rightIncluded: maxIncluded, rightMoment: maxValue,
	}}
}

// intervalsUnionAll returns the union of all values
func intervalsUnionAll(intervals []interval) []interval {
	size := len(intervals)
	if size <= 1 {
		return intervals
	}

	// initialize for loop
	var unions []interval
	currents := make([]interval, size)
	copy(currents, intervals)

	// make as many unions as possible
	for {
		sizeBefore := len(currents)
		for index, current := range currents {
			if current.empty {
				continue
			}

			for otherIndex, otherCurrrent := range currents {
				if otherCurrrent.empty {
					continue
				}

				if index < otherIndex {
					localUnion := current.union(otherCurrrent)
					for _, value := range localUnion {
						if !slices.ContainsFunc(unions, func(i interval) bool { return intervalEquals(i, value) }) {
							unions = append(unions, value)
						}
					}
				}
			}
		}

		sizeAfter := len(unions)
		if sizeBefore == sizeAfter {
			return unions
		} else {
			currents = unions
		}
	}
}

// remove removes other from i
// A - B = A inter (FULL - other)
// And if FULL - OTHER = X U Y then
// A inter (X U Y) = (A inter X) U (A inter Y)
func (i interval) remove(other interval) []interval {
	if other.empty {
		return nil
	} else if i.empty {
		return nil
	}

	complements := other.complement()
	size := len(complements)
	switch {
	case size == 0:
		return []interval{}
	case size == 1 && complements[0].empty:
		return []interval{i}
	case size == 1 && !complements[0].empty:
		result := i.intersection(complements[0])
		return []interval{result}
	default:
		leftOperand := i.intersection(complements[0])
		rightOperand := i.intersection(complements[1])
		var result []interval
		for _, r := range leftOperand.union(rightOperand) {
			if !r.empty {
				result = append(result, r)
			}
		}

		return result
	}
}

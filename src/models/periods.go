package models

import (
	"slices"
	"sort"
	"strings"
	"time"
)

// sortIntervals copies and sorts values by intervalCompare order
func sortIntervals(values []interval) []interval {
	size := len(values)
	if size == 0 {
		return []interval{}
	}

	newValues := make([]interval, size)
	copy(newValues, values)

	slices.SortStableFunc(newValues, intervalCompare)
	return newValues
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

// Union builds the union of two periods
func (p Period) Union(other Period) Period {
	var params []interval
	params = append(params, p.intervals...)
	params = append(params, other.intervals...)
	var result []interval
	for _, value := range intervalsUnionAll(params) {
		if !value.empty {
			result = append(result, value)
		}
	}

	return Period{intervals: result}
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

// Contains returns true if point is in the interval (as in set theory)
func (p Period) Contains(point time.Time) bool {
	for _, partition := range p.intervals {
		if partition.contains(point) {
			return true
		}
	}

	return false
}

// IsIncludedIn returns true if p is included in other
func (p Period) IsIncludedIn(other Period) bool {
	if len(p.intervals) == 0 {
		return true
	} else if len(other.intervals) == 0 {
		return false
	}

	for _, source := range p.intervals {
		contained := false
		for _, other := range other.intervals {
			if source.isIncludedIn(other) {
				contained = true
				break
			}
		}

		if !contained {
			return false
		}
	}

	return true
}

// Complement returns the complement of the period,
// that is the other period that forms a partition of full space with others
func (p Period) Complement() Period {
	if len(p.intervals) == 0 {
		return NewFullPeriod()
	}

	var result []interval
	var previousValue time.Time
	var previousFinite, previousIncluded bool
	// using the "completing hole" method: find all intervals so that the union would make full
	for index, value := range p.intervals {
		if value.isFull() {
			return Period{}
		}

		if index == 0 {
			// may complete left
			if value.leftFinite {
				// left completion
				completion := interval{
					empty: false, leftFinite: false,
					rightFinite: true, rightIncluded: !value.leftIncluded, rightMoment: value.leftMoment,
				}

				result = append(result, completion)
			}
		} else {
			// complete from previous to value
			completion := interval{
				empty: false, leftFinite: true, leftIncluded: !previousIncluded, leftMoment: previousValue,
				rightFinite: true, rightIncluded: !value.leftIncluded, rightMoment: value.leftMoment,
			}

			result = append(result, completion)
		}

		previousFinite, previousIncluded = value.rightFinite, value.rightIncluded
		previousValue = value.rightMoment
	}

	if previousFinite {
		// complete to reach +oo
		completion := interval{
			empty: false, rightFinite: false,
			leftFinite: true, leftIncluded: !previousIncluded, leftMoment: previousValue,
		}

		result = append(result, completion)
	}

	// result contains the partition that completes the initial period
	return Period{intervals: result}
}

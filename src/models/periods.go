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
	value := newFullInterval()
	return Period{intervals: []interval{value}}
}

// NewEmptyPeriod builds an empty period
func NewEmptyPeriod() Period {
	return Period{}
}

// NewFinitePeriod builds a period equivalent to a new finite interval (min, max)
// SPECIAL CASES: it may return an empty period according to mathematical definition
func NewFinitePeriod(min, max time.Time, minIncluded, maxIncluded bool) Period {
	content := newIntervalDuring(min, max, minIncluded, maxIncluded)
	if content.empty {
		return Period{}
	} else {
		return Period{intervals: []interval{content}}
	}
}

// NewPeriodSince builds a period equivalent to (leftLimit, +oo[
func NewPeriodSince(leftLimit time.Time, leftIn bool) Period {
	return Period{intervals: []interval{newIntervalSince(leftLimit, leftIn)}}
}

// NewPeriodUntil builds a period equivalent to ]-oo,rightLimit)
func NewPeriodUntil(rightLimit time.Time, rightIn bool) Period {
	content := newIntervalUntil(rightLimit, rightIn)
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
			value := sourceInterval.intersection(otherInterval)
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
				completion := newIntervalUntil(value.leftMoment, !value.leftIncluded)
				result = append(result, completion)
			}
		} else {
			// complete from previous to value
			completion := newIntervalDuring(previousValue, value.leftMoment, !previousIncluded, !value.leftIncluded)
			if !completion.empty {
				result = append(result, completion)
			}
		}

		previousFinite, previousIncluded = value.rightFinite, value.rightIncluded
		previousValue = value.rightMoment
	}

	if previousFinite {
		// complete to reach +oo
		completion := newIntervalSince(previousValue, !previousIncluded)
		result = append(result, completion)
	}

	// result contains the partition that completes the initial period
	return Period{intervals: result}
}

// Remove a period from another
func (p Period) Remove(other Period) Period {
	if len(p.intervals) == 0 {
		return other
	} else if len(other.intervals) == 0 {
		return p
	}

	size := len(p.intervals)
	currents := make([]interval, size)
	copy(currents, p.intervals)

	for _, value := range other.intervals {
		var remainings []interval
		for _, current := range currents {
			for _, rem := range current.remove(value) {
				if !rem.empty {
					remainings = append(remainings, rem)
				}
			}
		}

		currents = remainings
	}

	return Period{intervals: currents}
}

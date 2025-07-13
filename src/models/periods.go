package models

import (
	"slices"
	"sort"
	"strings"
	"time"
)

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
	var unions []interval
	unions = append(unions, p.intervals...)
	unions = append(unions, other.intervals...)
	return Period{intervals: intervalsUnionAll(unions)}
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

package structures

import (
	"errors"
	"maps"
	"slices"
	"time"
)

// Mapping defines a value that changes over time.
// Keys are serialized intervals, vaues are the values over time
type Mapping[T interface{ comparable }] map[interval]T

// NewValue returns element for the full period
func NewValue[T interface{ comparable }](element T) Mapping[T] {
	result := make(Mapping[T])
	full := newFullInterval()
	result[full] = element
	return result
}

// NewValueSince sets the value for period (leftMoment, +oo[
func NewValueSince[T interface{ comparable }](element T, leftMoment time.Time, leftIn bool) Mapping[T] {
	result := make(Mapping[T])
	moment := newIntervalSince(leftMoment, leftIn)
	result[moment] = element
	return result
}

// NewValueUntil sets the value for period ]-oo, rightMoment)
func NewValueUntil[T interface{ comparable }](element T, rightMoment time.Time, rightIn bool) Mapping[T] {
	result := make(Mapping[T])
	moment := newIntervalUntil(rightMoment, rightIn)
	result[moment] = element
	return result
}

// NewValueDuring sets the value for period (leftMoment, rightMoment)
func NewValueDuring[T interface{ comparable }](element T, leftMoment, rightMoment time.Time, leftIn, rightIn bool) Mapping[T] {
	result := make(Mapping[T])
	moment := newIntervalDuring(leftMoment, rightMoment, leftIn, rightIn)
	if !moment.empty {
		result[moment] = element
	}
	return result
}

// NewValueDuringPeriod sets the value during a given period
func NewValueDuringPeriod[T interface{ comparable }](element T, period Period) Mapping[T] {
	result := make(Mapping[T])
	var nonEmptyIntervals []interval
	for _, interval := range period.intervals {
		if !interval.empty {
			nonEmptyIntervals = append(nonEmptyIntervals, interval)
		}
	}

	// get first value to build the result
	result[nonEmptyIntervals[0]] = element

	// add other values
	nonEmptyIntervals = nonEmptyIntervals[1:]
	for len(nonEmptyIntervals) != 0 {
		current := nonEmptyIntervals[0]
		nonEmptyIntervals = nonEmptyIntervals[1:]
		result.addValue(element, current)
	}

	return result
}

// addValue sets value for given interval
func (m Mapping[T]) addValue(value T, i interval) {
	if i.empty {
		return
	}

	// we cannot change current map while reading it, so we create a local copy
	newValues := make(map[interval]T)

	// algorithm is:
	// for each period, current value
	//   if same values, then we want to regroup intervals
	//   else, we remove i to the current period, remaining part keeps its value
	// then we regroup the rest into a common interval with value as value
	var intervalsWithSameValue []interval
	for existingPeriod, existingValue := range m {
		if existingValue == value {
			intervalsWithSameValue = append(intervalsWithSameValue, existingPeriod)
		} else {
			remainings := existingPeriod.remove(i)
			for _, remaining := range remainings {
				if !remaining.empty {
					newValues[remaining] = existingValue
				}
			}
		}
	}

	// We take the union of all the intervals linked to value
	intervalsWithSameValue = append(intervalsWithSameValue, i)
	newCommonIntervals := intervalsUnionAll(intervalsWithSameValue)
	// for each period, link it to value
	for _, period := range newCommonIntervals {
		if !period.empty {
			newValues[period] = value
		}
	}

	// then perform the replacement:
	// delete current map
	// replace it all with the new values
	for k := range m {
		delete(m, k)
	}

	maps.Copy(m, newValues)
}

// Set value for all time
func (m Mapping[T]) Set(value T) {
	m.addValue(value, newFullInterval())
}

// SetSince sets value during (leftBound, +oo[
func (m Mapping[T]) SetSince(value T, leftBound time.Time, leftIn bool) {
	m.addValue(value, newIntervalSince(leftBound, leftIn))
}

// SetUntil sets value during period ]-oo, rightBound)
func (m Mapping[T]) SetUntil(value T, rightBound time.Time, rightIn bool) {
	m.addValue(value, newIntervalUntil(rightBound, rightIn))
}

// SetDuring sets value during period (leftBound, rightBound)
func (m Mapping[T]) SetDuring(value T, leftBound, rightBound time.Time, leftIn, rightIn bool) {
	period := newIntervalDuring(leftBound, rightBound, leftIn, rightIn)
	if !period.empty {
		m.addValue(value, period)
	}
}

// SetDuringPeriod sets the value during a given period
func (m Mapping[T]) SetDuringPeriod(value T, period Period) {
	for _, interval := range period.intervals {
		if !interval.empty {
			m.addValue(value, interval)
		}
	}
}

// GetValues returns all the values set
func (m Mapping[T]) GetValues() []T {
	var result []T
	for _, v := range m {
		if result == nil || !slices.Contains(result, v) {
			result = append(result, v)
		}
	}

	return result
}

// GetValue returns the value at a given time (value, true) or empty, false if not found
func (m Mapping[T]) GetValue(moment time.Time) (T, bool) {
	var empty T
	for i, v := range m {
		if i.contains(moment) {
			return v, true
		}
	}

	return empty, false
}

// Get returns the values and periods the values are set during.
// For instance, if content is ]-oo, now[ => "a", [after, +oo[ => "a"
// then result would be "a" => Period(]-oo, now[ union [after, +oo[)
func (m Mapping[T]) Get() map[T]Period {
	// first, group intervals per value.
	// One value => all intervals value is linked to
	values := make(map[T][]interval)
	for period, value := range m {
		if period.empty {
			continue
		} else if content, found := values[value]; !found {
			values[value] = []interval{period}
		} else {
			content = append(content, period)
			values[value] = content
		}
	}

	// for each value, just make the union
	result := make(map[T]Period)
	for value, intervals := range values {
		unioned := intervalsUnionAll(intervals)
		result[value] = Period{intervals: unioned}
	}

	return result
}

// Cut returns the mapping restricted to other.
// For instance,
// if mapping contains ]-oo,+oo[ => "a" and other is ]-oo, now],
// then result would be ]-oo, now] => "a"
func (m Mapping[T]) Cut(other Period) Mapping[T] {
	result := make(Mapping[T])
	for current, currentValue := range m {
		for _, i := range other.intervals {
			if inter := i.intersection(current); !inter.empty {
				result[inter] = currentValue
			}
		}
	}

	return result
}

// BuildCompactMapOfValues maps intervals to their string representation and returns the corresponding map
func (m Mapping[T]) BuildCompactMapOfValues() map[string]T {
	result := make(map[string]T)
	for interval, value := range m {
		result[interval.toString()] = value
	}

	return result
}

// LoadValuesFromCompactMap maps keys to intervals and returns corresponding map
func LoadValuesFromCompactMap[T interface{ comparable }](content map[string]T) (Mapping[T], error) {
	if content == nil {
		return nil, nil
	}

	result := make(Mapping[T])
	var globalErrors error
	for k, v := range content {
		if i, err := intervalFromString(k); err != nil {
			globalErrors = errors.Join(globalErrors, err)
		} else {
			result[i] = v
		}
	}

	return result, globalErrors
}

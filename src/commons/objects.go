package commons

import (
	"errors"
	"time"
)

// StateObject is basically an object with a map of attributes and values
type StateObject[T StateValue] struct {
	// id of the object
	id string
	// values is a map with keys as attributes, and related values
	values map[string]T
}

// GetType returns an object type
func (s *StateObject[T]) GetType() ModelableType {
	return TypeObject
}

// Id returns the id of the object
func (s *StateObject[T]) Id() string {
	return s.id
}

// SetValue changes an attribute value by name
func (s *StateObject[T]) SetValue(name string, value T) {
	if s.values == nil {
		s.values = make(map[string]T)
	}

	s.values[name] = value
}

// GetValue returns, if any, current value for that attribute
func (s *StateObject[T]) GetValue(name string) (T, bool) {
	var empty T
	if s == nil || s.values == nil {
		return empty, false
	} else if value, found := s.values[name]; !found {
		return empty, false
	} else {
		return value, true
	}
}

// Attributes returns the available attributes
func (s *StateObject[T]) Attributes() []string {
	if s == nil {
		return nil
	}

	var result []string
	for name := range s.values {
		result = append(result, name)
	}

	return result
}

// NewStateObject makes a new empty state object
func NewStateObject[T StateValue]() *StateObject[T] {
	result := new(StateObject[T])
	result.id = NewId()
	result.values = make(map[string]T)
	return result
}

// TimedStateObject defines an object with a state that changes over time.
// State is defined as a map of key => values depending over time.
// Keys are the attributes name, values depend on time.
// The object itself may be active at a given time.
// So a timed state object also implements TemporalObject.
type TimedStateObject[T StateValue] struct {
	// id of the object
	id string
	// lifetime of the object, that is its the activation period
	lifetime Period
	// attributes represent the state of the object as a time varying map
	attributes map[string]TimeDependentValues[T]
}

// NewTimedStateObject returns a forever lasting object with no attribute
func NewTimedStateObject[T StateValue]() *TimedStateObject[T] {
	result := new(TimedStateObject[T])
	result.id = NewId()
	result.lifetime = NewFullPeriod()
	result.attributes = make(map[string]TimeDependentValues[T])
	return result
}

// NewTimedStateObjectSince returns a object with a varying state, valid since creationTime
func NewTimedStateObjectSince[T StateValue](creationTime time.Time) *TimedStateObject[T] {
	base := NewTimedStateObject[T]()
	base.lifetime = NewPeriodSince(creationTime, true)
	return base
}

// NewTimedStateObjectDuring returns an object with a varying state, valid during a period.
// It may raise an error if endTime is before startTime
func NewTimedStateObjectDuring[T StateValue](traits []string, startTime, endTime time.Time) (*TimedStateObject[T], error) {
	if endTime.Before(startTime) {
		return nil, errors.New("cannot make an object with an end date before its start date")
	}

	base := NewTimedStateObject[T]()
	base.lifetime = NewFinitePeriod(startTime, endTime, true, true)
	return base, nil
}

// Id returns an unique id for that object.
// It is constant for that object, and globally unique.
func (t *TimedStateObject[T]) Id() string {
	return t.id
}

// GetType returns an object type
func (t *TimedStateObject[T]) GetType() ModelableType {
	return TypeObject
}

// ActivePeriod returns the lifetime of that object
func (t *TimedStateObject[T]) ActivePeriod() Period {
	return t.lifetime
}

// SetActivePeriod forces the object lifetime
func (t *TimedStateObject[T]) SetActivePeriod(p Period) {
	t.lifetime = p
}

// Attributes return the attributes of the object
func (o *TimedStateObject[T]) Attributes() []string {
	var result []string
	for name := range o.attributes {
		result = append(result, name)
	}

	if len(result) == 0 {
		result = make([]string, 0)
		return result
	} else {
		return SliceReduce(result)
	}
}

// SetValueDuringPeriod changes that attribute to set value during period.
// If object is nil or period is empty, no action.
// Else value changes during that period no matter the object's lifetime
func (o *TimedStateObject[T]) SetValueDuringPeriod(attribute string, value T, period Period) {
	if o == nil {
		return
	} else if period.IsEmpty() {
		return
	}

	if o.attributes == nil {
		o.attributes = make(map[string]TimeDependentValues[T])
	}

	if attr, found := o.attributes[attribute]; !found {
		o.attributes[attribute] = NewValueDuringPeriod(value, period)
	} else {
		attr.SetDuringPeriod(value, period)
		o.attributes[attribute] = attr
	}
}

// SetValue sets a value for that attribute
func (o *TimedStateObject[T]) SetValue(attribute string, value T) {
	if o == nil {
		return
	}

	o.SetValueDuringPeriod(attribute, value, NewFullPeriod())
}

// SetValueSince sets the value for that attribute since startingTime
func (o *TimedStateObject[T]) SetValueSince(attribute string, value T, startingTime time.Time, includeStartingTime bool) {
	if o == nil {
		return
	}

	period := NewPeriodSince(startingTime, includeStartingTime)
	o.SetValueDuringPeriod(attribute, value, period)
}

// SetValueUntil sets the value for that attribute until endingTime
func (o *TimedStateObject[T]) SetValueUntil(attribute string, value T, endingTime time.Time, includeEndingTime bool) {
	if o == nil {
		return
	}

	period := NewPeriodUntil(endingTime, includeEndingTime)
	o.SetValueDuringPeriod(attribute, value, period)
}

// SetValueDuring sets value for that attribute during the interval [startingTime, endingTime] (both included)
func (o *TimedStateObject[T]) SetValueDuring(attribute string, value T, startingTime, endingTime time.Time) {
	if o == nil {
		return
	}

	period := NewFinitePeriod(startingTime, endingTime, true, true)
	o.SetValueDuringPeriod(attribute, value, period)
}

// GetAllValues returns all the values for all attributes (including the ones with no value)
// Two options:
// Either reduceToObjectLifetime is true and we get values only during object lifetime
// Or reduceToObjectLifetime is false and we get all values
func (o *TimedStateObject[T]) GetAllValues(reduceToObjectLifetime bool) map[string][]T {
	if o == nil {
		return nil
	}

	result := make(map[string][]T)

	// for each attribute
	for name, attr := range o.attributes {
		// values contain all the possible values
		var values []T
		// for each value and then period for that value
		for value, period := range attr.Get() {
			if reduceToObjectLifetime {
				if !period.IsEmpty() && !period.Intersection(o.lifetime).IsEmpty() {
					values = append(values, value)
				}
			} else {
				values = append(values, value)
			}
		}

		// we made the values, so set for that attribute
		result[name] = SliceDeduplicate(values)
	}

	return result
}

// GetValue returns the value for an attribute (by name) if any.
// Result (if any) is then the mapping value -> validity, true or nil, false for no match.
// Depending on reduceToObjectLifetime:
// Either it is true and then validity is the intersection of the object lifetime and the attribute validity
// Or we keep values and matching period as is
func (o *TimedStateObject[T]) GetValue(attribute string, reduceToObjectLifetime bool) (map[T]Period, bool) {
	if o == nil {
		return nil, false
	}

	// values are the values from the attribute.
	var values map[T]Period
	if attr, found := o.attributes[attribute]; !found {
		return nil, false
	} else {
		values = attr.Get()
	}

	// result contains the intersection with the object's lifetime
	result := make(map[T]Period)
	for key, period := range values {
		if reduceToObjectLifetime {
			inter := period.Intersection(o.lifetime)
			if !inter.IsEmpty() {
				result[key] = inter
			}
		} else {
			result[key] = period
		}
	}

	return result, true
}

// GetValues returns the values and their activity (during the object lifetime).
// If no value was set for that attribute, return nil
func (o *TimedStateObject[T]) GetValues(attribute string) map[T]Period {
	if result, found := o.GetValue(attribute, true); found {
		return result
	} else {
		return nil
	}
}

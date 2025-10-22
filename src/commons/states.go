package commons

import (
	"maps"
	"time"
)

// StateValue is the definition of accepted types
type StateValue interface{ string | int | float64 | bool }

// StateDescription describes a state with no memory of previous values
type StateDescription[T StateValue] interface {
	// Id returns the id of the state, not the object
	Id() string
	// Values returns the current state as a map of attributes and related values
	Values() map[string]T
}

// TemporalStateDescription is historicized content
type TemporalStateDescription[T StateValue] interface {
	// Id of the state description, not the object
	Id() string
	// ActivePeriod returns the period the state makes sense (is active)
	ActivePeriod() Period
	// Values returns, for each attribute, the values and related periods.
	// For a given attribute, pay attention that period for value is NOT linked to activity.
	// For instance, when using SetValue, it changes the value for the full period.
	Values() map[string]map[T]Period
}

// StateReader declares the ability to deal with a time independent state.
// It is NOT linked to objects only.
// For instance, some structures may describe themselves.
// Some constraints with parameters may also describe themselves.
type StateReader[T StateValue] interface {
	// Read returns the current state of this element
	Read() StateDescription[T]
}

// StateHandler is reading and updating a state.
type StateHandler[T StateValue] interface {
	// Handler needs ability to read
	StateReader[T]
	// SetValue sets value for that attribute
	SetValue(name string, value T)
	// SetValues sets values for a group of attributes
	SetValues(values map[string]T)
	// Remove excludes an attribute (if present).
	// It returns true if name was found, false otherwise
	Remove(name string) bool
}

// TemporalStateReader reads the temporal state of a temporal source.
// Source may be itself for elements able to self describe
type TemporalStateReader[T StateValue] interface {
	// Read() returns current state, state is time dependent
	Read() TemporalStateDescription[T]
	// ReadAtTime() returns state at that time as a constant content
	ReadAtTime(moment time.Time) StateDescription[T]
}

// TemporalStateHandler declares a state that varies over time.
// It means the ability to list attributes,
// for each attribute, be able to get values and related periods,
// and change those values during a given period.
// There is no SetValues(map[string]map[T]Period because periods may overlap and map does not guarantee order.
// For instance, given an attribute a,
// 10 => Full Period , 100 => [now, +oo[ would create two different states whether 10 or 100 is picked first
type TemporalStateHandler[T StateValue] interface {
	// TemporalStateReader is necessary to change state over time
	TemporalStateReader[T]
	// TemporalHandler to deal with state lifetime
	TemporalHandler
	// SetValueDuringPeriod sets value for that attribute during a given period
	SetValueDuringPeriod(name string, value T, period Period)
	// Remove removes an attribute by name (if any) and returns if it was present before removal
	Remove(name string) bool
}

// StateRepresentation is basically as state as a map of attributes and values
type StateRepresentation[T StateValue] struct {
	// id of source
	id string
	// values is a map with keys as attributes, and related values
	values map[string]T
}

// Id returns the id of the state.
// NOT THE source, the id of the state.
// Reason is that we may want to save it separated from the object
func (s *StateRepresentation[T]) Id() string {
	return s.id
}

// Values returns current state as a map
func (s *StateRepresentation[T]) Values() map[string]T {
	if s == nil {
		return nil
	}

	result := make(map[string]T)
	if s.values != nil {
		maps.Copy(result, s.values)
	}

	return result
}

// SetValue changes an attribute value by name
func (s *StateRepresentation[T]) SetValue(name string, value T) {
	if s.values == nil {
		s.values = make(map[string]T)
	}

	s.values[name] = value
}

// SetValues set values from values, does not affect other values
func (s *StateRepresentation[T]) SetValues(values map[string]T) {
	if s.values == nil {
		s.values = make(map[string]T)
	}

	maps.Copy(s.values, values)
}

// Remove excludes an attribute and returns true, or does nothing and return false if attribute was not set
func (s *StateRepresentation[T]) Remove(name string) bool {
	if s == nil || s.values == nil {
		return false
	}

	_, found := s.values[name]
	if !found {
		return false
	} else {
		delete(s.values, name)
		return true
	}
}

// GetValue returns, if any, current value for that attribute
func (s *StateRepresentation[T]) GetValue(name string) (T, bool) {
	var empty T
	if s == nil {
		return empty, false
	}

	if s.values == nil {
		return empty, false
	} else if value, found := s.values[name]; !found {
		return empty, false
	} else {
		return value, true
	}
}

// Attributes returns the available attributes
func (s *StateRepresentation[T]) Attributes() []string {
	if s == nil {
		return nil
	}

	var result []string
	for name := range s.values {
		result = append(result, name)
	}

	return result
}

// NewStateRepresentation makes a new empty state
func NewStateRepresentation[T StateValue]() *StateRepresentation[T] {
	result := new(StateRepresentation[T])
	result.id = NewId()
	result.values = make(map[string]T)
	return result
}

// NewStateRepresentationFrom builds a new state, wih preset values
func NewStateRepresentationFrom[T StateValue](values map[string]T) *StateRepresentation[T] {
	result := NewStateRepresentation[T]()
	for attr, value := range values {
		result.SetValue(attr, value)
	}

	return result
}

// TimedStateRepresentation defines a state that changes over time.
// State is defined as a map of key => values depending over time.
// Keys are the attributes name, values depend on time.
type TimedStateRepresentation[T StateValue] struct {
	// id of the state (not the object)
	id string
	// lifetime of the source, that is its the activation period
	lifetime Period
	// attributes represent the state as a time varying map
	attributes map[string]TimeDependentValues[T]
}

// Id returns an unique id for that state.
func (t *TimedStateRepresentation[T]) Id() string {
	return t.id
}

// ActivePeriod returns the lifetime of that state
func (t *TimedStateRepresentation[T]) ActivePeriod() Period {
	return t.lifetime
}

// SetActivePeriod forces the state lifetime
func (t *TimedStateRepresentation[T]) SetActivePeriod(p Period) {
	t.lifetime = p
}

// Attributes return the attributes of the state
func (t *TimedStateRepresentation[T]) Attributes() []string {
	var result []string
	for name := range t.attributes {
		result = append(result, name)
	}

	if len(result) == 0 {
		result = make([]string, 0)
		return result
	} else {
		return SliceReduce(result)
	}
}

// Remove removes an attribute by name and returns if the attribute was set already (no action if not)
func (t *TimedStateRepresentation[T]) Remove(name string) bool {
	if t == nil || t.attributes == nil {
		return false
	}

	_, found := t.attributes[name]
	if found {
		delete(t.attributes, name)
	}

	return found
}

// SetValueDuringPeriod changes that attribute to set value during period.
// If state is nil or period is empty, no action.
// Else value changes during that period no matter the state's lifetime
func (t *TimedStateRepresentation[T]) SetValueDuringPeriod(attribute string, value T, period Period) {
	if t == nil {
		return
	} else if period.IsEmpty() {
		return
	}

	if t.attributes == nil {
		t.attributes = make(map[string]TimeDependentValues[T])
	}

	if attr, found := t.attributes[attribute]; !found {
		t.attributes[attribute] = NewValueDuringPeriod(value, period)
	} else {
		attr.SetDuringPeriod(value, period)
		t.attributes[attribute] = attr
	}
}

// SetValue sets a value for that attribute
func (t *TimedStateRepresentation[T]) SetValue(attribute string, value T) {
	if t == nil {
		return
	}

	t.SetValueDuringPeriod(attribute, value, NewFullPeriod())
}

// SetValueSince sets the value for that attribute since startingTime
func (t *TimedStateRepresentation[T]) SetValueSince(attribute string, value T, startingTime time.Time, includeStartingTime bool) {
	if t == nil {
		return
	}

	period := NewPeriodSince(startingTime, includeStartingTime)
	t.SetValueDuringPeriod(attribute, value, period)
}

// SetValueUntil sets the value for that attribute until endingTime
func (t *TimedStateRepresentation[T]) SetValueUntil(attribute string, value T, endingTime time.Time, includeEndingTime bool) {
	if t == nil {
		return
	}

	period := NewPeriodUntil(endingTime, includeEndingTime)
	t.SetValueDuringPeriod(attribute, value, period)
}

// SetValueDuring sets value for that attribute during the interval [startingTime, endingTime] (both included)
func (t *TimedStateRepresentation[T]) SetValueDuring(attribute string, value T, startingTime, endingTime time.Time) {
	if t == nil {
		return
	}

	period := NewFinitePeriod(startingTime, endingTime, true, true)
	t.SetValueDuringPeriod(attribute, value, period)
}

// GetAllValues returns all the values for all attributes (including the ones with no value)
// Two options:
// Either reduceToObjectLifetime is true and we get values only during state lifetime
// Or reduceToObjectLifetime is false and we get all values
func (t *TimedStateRepresentation[T]) GetAllValues(reduceToObjectLifetime bool) map[string][]T {
	if t == nil {
		return nil
	}

	result := make(map[string][]T)

	// for each attribute
	for name, attr := range t.attributes {
		// values contain all the possible values
		var values []T
		// for each value and then period for that value
		for value, period := range attr.Get() {
			if reduceToObjectLifetime {
				if !period.IsEmpty() && !period.Intersection(t.lifetime).IsEmpty() {
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
// Either it is true and then validity is the intersection of the state lifetime and the attribute validity
// Or we keep values and matching period as is
func (t *TimedStateRepresentation[T]) GetValue(attribute string, reduceToObjectLifetime bool) (map[T]Period, bool) {
	if t == nil {
		return nil, false
	}

	// values are the values from the attribute.
	var values map[T]Period
	if attr, found := t.attributes[attribute]; !found {
		return nil, false
	} else {
		values = attr.Get()
	}

	// result contains the intersection with the state lifetime
	result := make(map[T]Period)
	for key, period := range values {
		if reduceToObjectLifetime {
			inter := period.Intersection(t.lifetime)
			if !inter.IsEmpty() {
				result[key] = inter
			}
		} else {
			result[key] = period
		}
	}

	return result, true
}

// GetValues returns the values and their activity (during the state lifetime).
// If no value was set for that attribute, return nil
func (t *TimedStateRepresentation[T]) GetValues(attribute string) map[T]Period {
	if result, found := t.GetValue(attribute, true); found {
		return result
	} else {
		return nil
	}
}

// Values returns the full state as a map of attributes and time dependent values
func (t *TimedStateRepresentation[T]) Values() map[string]map[T]Period {
	if t == nil {
		return nil
	}

	result := make(map[string]map[T]Period)
	for name, attr := range t.attributes {
		current := make(map[T]Period)
		maps.Copy(current, attr.Get())
		result[name] = current
	}

	return result
}

// Snapshot returns matching state with values fixed at given time.
// If moment is NOT in lifetime, it returns empty map.
func (t *TimedStateRepresentation[T]) Snapshot(moment time.Time) *StateRepresentation[T] {
	if t == nil {
		return nil
	}

	result := make(map[string]T)
	if t.lifetime.Contains(moment) {
		for name, attr := range t.attributes {
			for value, period := range attr.Get() {
				if period.Contains(moment) {
					result[name] = value
				}
			}
		}
	}

	return NewStateRepresentationFrom(result)
}

// NewSimpleTimedStateRepresentation returns a time dependent state for a constantly living element
func NewSimpleTimedStateRepresentation[T StateValue]() *TimedStateRepresentation[T] {
	return NewTimedStateRepresentation[T](NewFullPeriod())
}

// NewTimedStateRepresentation returns a state with no attribute, during a given period
func NewTimedStateRepresentation[T StateValue](period Period) *TimedStateRepresentation[T] {
	result := new(TimedStateRepresentation[T])
	result.id = NewId()
	result.lifetime = period
	result.attributes = make(map[string]TimeDependentValues[T])
	return result
}

// NewTimedStateRepresentationFrom builds a new state, valid during period, wih preset values
func NewTimedStateRepresentationFrom[T StateValue](period Period, values map[string]T) *TimedStateRepresentation[T] {
	result := NewTimedStateRepresentation[T](period)
	for attr, value := range values {
		result.SetValue(attr, value)
	}

	return result
}

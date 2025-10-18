package commons

import "time"

// StateValue is the definition of accepted types
type StateValue interface{ string | int | float64 | bool }

// StateReader declares the ability to deal with a time independent state.
// It is NOT linked to objects only.
// For instance, some structures may describe themselves.
// Some constraints with parameters may also describe themselves.
type StateReader[T StateValue] interface {
	// Values returns the values
	Read() StateDescription[T]
}

// StateDescription describes a state
type StateDescription[T StateValue] interface {
	// Id returns the id if any, true if source has an id
	Id() (string, bool)
	// Values returns the values of the source when asked for description
	Values() map[string]T
}

// constantState is a static state description
type constantState[T StateValue] struct {
	// id is the id of the source, if any
	id string
	// hasId is true when the source had an id
	hasId bool
	// values are the values from the source
	values map[string]T
}

// Id returns the id of the source, if any, and true if found
func (c constantState[T]) Id() (string, bool) {
	return c.id, c.hasId
}

// Values returns current content
func (c constantState[T]) Values() map[string]T {
	return c.values
}

// NewStateDescription builds a new state based on values
func NewStateDescription[T StateValue](values map[string]T) StateDescription[T] {
	return constantState[T]{id: "", hasId: false, values: values}
}

// NewStateDescriptionWithId builds a new state based on values and a given id
func NewStateDescriptionWithId[T StateValue](id string, values map[string]T) StateDescription[T] {
	return constantState[T]{id: id, hasId: true, values: values}
}

// StateHandler is reading and updating a state.
type StateHandler[T StateValue] interface {
	// Handler needs ability to read
	StateReader[T]
	// SetValue sets value for that attribute
	SetValue(name string, value T)
}

// TemporalStateReader reads the temporal state of a temporal source
type TemporalStateReader[T StateValue] interface {
	// Read() returns current state, state is time dependent
	Read() TemporalStateDescription[T]
	// ReadAtTime() returns state at that time as a constant content
	ReadAtTime(moment time.Time) StateDescription[T]
}

// TemporalStateDescription describes a state that varies over time
type TemporalStateDescription[T StateValue] interface {
	Id() (string, bool)
	// ActivePeriod returns the period of activity from the source.
	// If source does not implement temporal, it returns nil, false.
	// Else, it returns the active period, true
	ActivePeriod() (Period, bool)
	// Values returns, for each attribute, the values over time.
	// For instance, let us consider a source with unique attribute named attr, and values a => [now, +oo[
	// then result should be attr => a => [now, +oo[
	Values() map[string]map[T]Period
}

// temporalStateContainer defines basic implementation of a TemporalStateDescription
type temporalStateContainer[T StateValue] struct {
	// id of the source
	id string
	// hasId is true if the source has an id
	hasId bool
	// period is the activity period of the source
	period Period
	// hasPeriod is true if source implements temporal (and then has an activity period)
	hasPeriod bool
	// values links attributes with their time dependent values
	values map[string]map[T]Period
}

// Id returns the id if any, or "", false for no id
func (t temporalStateContainer[T]) Id() (string, bool) {
	return t.id, t.hasId
}

// ActivePeriod returns the active period, if any, and nil, false if not found
func (t temporalStateContainer[T]) ActivePeriod() (Period, bool) {
	return t.period, t.hasPeriod
}

// Values returns the attributes names and their time dependent values
func (t temporalStateContainer[T]) Values() map[string]map[T]Period {
	return t.values
}

// snapshot returns the state at a given moment
func (t temporalStateContainer[T]) snapshot(moment time.Time) StateDescription[T] {
	snapshotValues := make(map[string]T)
	for attr, values := range t.values {
		for value, period := range values {
			if period.Contains(moment) {
				snapshotValues[attr] = value
			}
		}
	}

	var result constantState[T]
	result.id = t.id
	result.hasId = t.hasId
	result.values = snapshotValues
	return result
}

// TemporalStateHandler declares a state that varies over time.
// It means the ability to list attributes,
// for each attribute, be able to get values and related periods,
// and change those values during a given period.
type TemporalStateHandler[T StateValue] interface {
	// TemporalStateReader is necessary to change state over time
	TemporalStateReader[T]
	// SetValueDuringPeriod sets value for that attribute during a given period
	SetValueDuringPeriod(name string, value T, period Period)
}

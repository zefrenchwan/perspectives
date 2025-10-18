package commons

import (
	"time"
)

// StateObject is an object with a state.
// This object may change over time, but it has no lifetime.
// Typical use case would be a particule to simulate.
type StateObject[T StateValue] interface {
	// StateObject is a model object
	ModelObject
	// By definition, we may read the state
	StateReader[T]
	// Handler returns a state handler to modify the object.
	// An object may modify its state (for example, we may move an arm).
	// But a handler is a tool to modity an object with a direct access to its state.
	// For instance, a system is subject to gravity no matter what.
	Handler() StateHandler[T]
}

// TemporalStateObject is an object with a lifetime and historicized state.
// Its state is then a TemporalStateHandler and it has an active period.
type TemporalStateObject[T StateValue] interface {
	// Temporal object is an object with an active period
	TemporalObject
	// TemporalStateReader to get a state with historicized content
	TemporalStateReader[T]
	// Handler returns an object to be able to change the state of the object
	Handler() TemporalStateHandler[T]
}

// ModelStateObject implements StateObject by decorating a state.
type ModelStateObject[T StateValue] struct {
	// id of the object
	id string
	// State is a shared (pointer) implementation of a state
	State *StateRepresentation[T]
}

// Id returns the object id
func (m ModelStateObject[T]) Id() string {
	return m.id
}

// GetType returns TypeObject for sure
func (m ModelStateObject[T]) GetType() ModelableType {
	return TypeObject
}

// Attributes returns set attributes
func (m ModelStateObject[T]) Attributes() []string {
	return m.State.Attributes()
}

// GetValue returns the value for a given attribute if any
func (m ModelStateObject[T]) GetValue(attribute string) (T, bool) {
	return m.State.GetValue(attribute)
}

// Read returns the state description
func (m ModelStateObject[T]) Read() StateDescription[T] {
	return m.State
}

// SetValue sets the value for that attribute
func (m ModelStateObject[T]) SetValue(name string, value T) {
	m.State.SetValue(name, value)
}

// Handler returns a state handler to modify the state
func (m ModelStateObject[T]) Handler() StateHandler[T] {
	return m
}

// NewModelStateObject returns a new empty ModelStateObject
func NewModelStateObject[T StateValue]() ModelStateObject[T] {
	return ModelStateObject[T]{State: NewStateRepresentation[T]()}
}

// TemporalModelStateObject is a model object to deal with historicized state
type TemporalModelStateObject[T StateValue] struct {
	// id of the object
	id string
	// State is the shared historiziced state
	State *TimedStateRepresentation[T]
}

// ActivePeriod gets the active period of the object
func (m TemporalModelStateObject[T]) ActivePeriod() Period {
	return m.State.ActivePeriod()
}

// SetActivePeriod sets active period of the object
func (m TemporalModelStateObject[T]) SetActivePeriod(period Period) {
	m.State.SetActivePeriod(period)
}

// GetType returns TypeObject because this component is an object
func (m TemporalModelStateObject[T]) GetType() ModelableType {
	return TypeObject
}

// Id returns the id of the object
func (m TemporalModelStateObject[T]) Id() string {
	return m.id
}

// GetValue returns the historicized content for an attribute
func (m TemporalModelStateObject[T]) GetValue(name string) (map[T]Period, bool) {
	return m.State.GetValue(name, true)
}

// Read returns the historicized state for that object.
func (m TemporalModelStateObject[T]) Read() TemporalStateDescription[T] {
	return m.State
}

// ReadAtTime builds the state at a given moment
func (m TemporalModelStateObject[T]) ReadAtTime(moment time.Time) StateDescription[T] {
	result := m.State.Snapshot(moment)
	return result
}

// SetValueDuringPeriod changes value for a given attribute during a given period
func (m TemporalModelStateObject[T]) SetValueDuringPeriod(name string, value T, period Period) {
	m.State.SetValueDuringPeriod(name, value, period)
}

// Handler returns an handler to modify historicized state
func (m TemporalModelStateObject[T]) Handler() TemporalStateHandler[T] {
	return m
}

// SetValue changes the value for that attribute
func (m TemporalModelStateObject[T]) SetValue(name string, value T) {
	m.State.SetValue(name, value)
}

// NewTemporalModelStateObject returns a new empty TemporalModelStateObject
func NewTemporalModelStateObject[T StateValue](lifetime Period) TemporalModelStateObject[T] {
	var result TemporalModelStateObject[T]
	result.id = NewId()
	result.State = NewTimedStateRepresentation[T](lifetime)
	return result
}

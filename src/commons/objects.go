package commons

import "time"

// ModelObject is the component that runs in the structure.
type ModelObject interface {
	// Linkable to put in links
	Linkable
	// An object is a component of a model
	ModelComponent
}

// StateObject is an object, with a lifetime, and a state
type StateObject[T StateValue] struct {
	// id returns the id of the object
	id string
	// state allows to read and change current state
	StateRepresentation[T]
	// activity to deal with time.
	// Current object becomes a temporal reader, then
	activity Period
}

// Id returns the id of the object
func (s *StateObject[T]) Id() string {
	return s.id
}

// GetType returns TypeObject
func (s *StateObject[T]) GetType() ModelableType {
	return TypeObject
}

// ActivePeriod returns current period of activity
func (s *StateObject[T]) ActivePeriod() Period {
	return s.activity
}

// SetActivePeriod forces current period of activity
func (s *StateObject[T]) SetActivePeriod(newPeriod Period) {
	s.activity = newPeriod
}

// NewStateObject returns an empty state object living forever
func NewStateObject[T StateValue]() *StateObject[T] {
	result := new(StateObject[T])
	result.activity = NewFullPeriod()
	result.id = NewId()
	result.StateRepresentation = NewStateRepresentation[T]()
	return result
}

// NewStateObjectSince creates a state object active since its creation
func NewStateObjectSince[T StateValue](creation time.Time) *StateObject[T] {
	result := NewStateObject[T]()
	result.activity = NewPeriodSince(creation, true)
	return result
}

// TemporalStateObject is an object with an activity and time dependent values
type TemporalStateObject[T StateValue] struct {
	// id of the object
	id string
	// state deals with time dependent values AND activity
	*TimedStateRepresentation[T]
}

// Id returns the object id
func (t *TemporalStateObject[T]) Id() string {
	return t.id
}

// GetType flags this object as an object
func (t *TemporalStateObject[T]) GetType() ModelableType {
	return TypeObject
}

// NewTemporalStateObject creates a new time dependent object active during period
func NewTemporalStateObject[T StateValue](period Period) *TemporalStateObject[T] {
	result := new(TemporalStateObject[T])
	result.id = NewId()
	result.TimedStateRepresentation = NewTimedStateRepresentation[T](period)
	return result
}

package commons

import "time"

// ModelObject is the component that runs in the structure.
type ModelObject interface {
	// Linkable to put in links
	Linkable
	// An object is a component of a model
	ModelComponent
}

// baseObject implements object the simplest way
type baseObject struct {
	// id of the object
	id string
}

// Id returns the id of the object
func (b baseObject) Id() string {
	return b.id
}

// GetType returns TypeObject to flag element as an object
func (b baseObject) GetType() ModelableType {
	return TypeObject
}

// NewModelObject returns a new object
func NewModelObject() ModelObject {
	return baseObject{id: NewId()}
}

// StateObject is an object, with a lifetime, and a state
type StateObject[T StateValue] struct {
	// state object is an object
	ModelObject
	// state allows to read and change current state
	StateRepresentation[T]
	// activity to deal with time.
	// Current object becomes a temporal reader, then
	activity Period
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
	result.ModelObject = NewModelObject()
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
	// a temporal state object is an object
	ModelObject
	// state deals with time dependent values AND activity
	*TimedStateRepresentation[T]
}

// NewTemporalStateObject creates a new time dependent object active during period
func NewTemporalStateObject[T StateValue](period Period) *TemporalStateObject[T] {
	result := new(TemporalStateObject[T])
	result.ModelObject = NewModelObject()
	result.TimedStateRepresentation = NewTimedStateRepresentation[T](period)
	return result
}

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
	state StateRepresentation[T]
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

// GetValue returns the value for that attribute (if any) or empty, false
func (s *StateObject[T]) GetValue(name string) (T, bool) {
	return s.state.GetValue(name)
}

// SetValue forces value for that attribute (by name)
func (s *StateObject[T]) SetValue(name string, value T) {
	s.state.SetValue(name, value)
}

// Read returns the current state of this element
func (s *StateObject[T]) Read() StateDescription[T] {
	return s.state.Read()
}

// SetValues sets values for a group of attributes
func (s *StateObject[T]) SetValues(values map[string]T) {
	s.state.SetValues(values)
}

// Remove excludes an attribute (if present).
func (s *StateObject[T]) Remove(name string) bool {
	return s.state.Remove(name)
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
	result.state = NewStateRepresentation[T]()
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
	state TemporalStateHandler[T]
}

// Id returns the object id
func (t *TemporalStateObject[T]) Id() string {
	return t.id
}

// GetType flags this object as an object
func (t *TemporalStateObject[T]) GetType() ModelableType {
	return TypeObject
}

// Read() returns current state, state is time dependent
func (t *TemporalStateObject[T]) Read() TemporalStateDescription[T] {
	return t.state.Read()
}

// ReadAtTime() returns state at that time as a constant content
func (t *TemporalStateObject[T]) ReadAtTime(moment time.Time) StateDescription[T] {
	return t.state.ReadAtTime(moment)
}

// SetValueDuringPeriod sets value for that attribute during a given period
func (t *TemporalStateObject[T]) SetValueDuringPeriod(name string, value T, period Period) {
	t.state.SetValueDuringPeriod(name, value, period)
}

// Remove removes an attribute by name (if any) and returns if it was present before removal
func (t *TemporalStateObject[T]) Remove(name string) bool {
	return t.state.Remove(name)
}

// ActivePeriod returns object active period
func (t *TemporalStateObject[T]) ActivePeriod() Period {
	return t.state.ActivePeriod()
}

// SetActivePeriod changes the period for that object
func (t *TemporalStateObject[T]) SetActivePeriod(newPeriod Period) {
	t.state.SetActivePeriod(newPeriod)
}

// NewTemporalStateObject creates a new time dependent object active during period
func NewTemporalStateObject[T StateValue](period Period) *TemporalStateObject[T] {
	result := new(TemporalStateObject[T])
	result.id = NewId()
	result.state = NewTimedStateRepresentation[T](period)
	return result
}

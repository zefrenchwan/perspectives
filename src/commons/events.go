package commons

import "time"

// Event is the general definition of an event:
// messages between agents, structures triggering a change of state, etc.
type Event interface {
	// each event has a unique id
	Identifiable
	// Source returns the unique source of the event
	Source() ModelComponent
}

// EventProcessor processes events by reacting to events.
type EventProcessor interface {
	// Process the notified events, may emit some events or raise an error
	Process(notified []Event) ([]Event, error)
}

// ObjectEventProcessor is an object able to deal with events
type ObjectEventProcessor interface {
	// ObjectEventProcessor is a model object (then may be included in a structure)
	ModelObject
	// ObjectEventProcessor is able to process events
	EventProcessor
}

// EventTick notifies an event processor to run one step further
type EventTick struct {
	// id of the event
	id string
	// source is the structure emitting the tick
	source ModelStructure
}

// Id returns the id of the event
func (t EventTick) Id() string {
	return t.id
}

// Source returns the model structure as a component
func (t EventTick) Source() ModelComponent {
	return t.source
}

// NewEventTick returns a new tick for that moment
func NewEventTick(source ModelStructure) EventTick {
	return EventTick{id: NewId(), source: source}
}

// EventLifetimeEnd ends lifetime of temporal values at end time
type EventLifetimeEnd struct {
	// id is the event id
	id string
	// source is the structure emitting the event
	source ModelStructure
	// end is the moment a temporal ends
	end time.Time
}

// Id returns that event id
func (l EventLifetimeEnd) Id() string {
	return l.id
}

// Source returns the structure source
func (l EventLifetimeEnd) Source() ModelComponent {
	return l.source
}

// ProcessingTime returns the time to end the period
func (l EventLifetimeEnd) ProcessingTime() time.Time {
	return l.end
}

// NewEventLifetimeEnd builds a new event to end a lifetime at given time from that structure
func NewEventLifetimeEnd(source ModelStructure, end time.Time) EventLifetimeEnd {
	return EventLifetimeEnd{id: NewId(), source: source, end: end.Truncate(TIME_PRECISION)}
}

// EventStateChanges notifies a state handler that it should set those values for those attributes
type EventStateChanges[T StateValue] interface {
	// this is an event
	Event
	// Changes are the changes to perform as key values
	Changes() map[string]T
	// ProcessingTime returns the moment to change values.
	// For temporal values, it means that we end previous values at that date.
	// For simple state values, it is just ignored
	ProcessingTime() time.Time
}

// timedEventStateChange is a simple EventStateChanges
type timedEventStateChange[T StateValue] struct {
	// id returns the id of the event
	id string
	// source is the structure that emitted change event
	source ModelStructure
	// moment is shared with all attributes and values.
	// It does not apply for state handlers (they just store current state)
	moment time.Time
	// values are the values to change.
	// It contains new values to set
	values map[string]T
}

// Id returns that event id
func (t timedEventStateChange[T]) Id() string {
	return t.id
}

// Source returns the source that created the event
func (t timedEventStateChange[T]) Source() ModelComponent {
	return t.source
}

// ProcessingTime returns the time to apply changes
func (t timedEventStateChange[T]) ProcessingTime() time.Time {
	return t.moment
}

// Changes returns the changes to force on the processor
func (t timedEventStateChange[T]) Changes() map[string]T {
	return t.values
}

// NewEventStateChanges defines a source setting values since given moment
func NewEventStateChanges[T StateValue](source ModelStructure, moment time.Time, values map[string]T) EventStateChanges[T] {
	return timedEventStateChange[T]{id: NewId(), source: source, moment: moment, values: values}
}

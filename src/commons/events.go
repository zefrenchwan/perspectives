package commons

import (
	"time"
)

// Event is the general definition of an event:
// messages between agents, structures triggering a change of state, etc.
type Event interface {
	// each event has a unique id
	Identifiable
	// ProcessingTime returns the moment to consider the event should be processed
	ProcessingTime() time.Time
}

// simpleEvent is the most basic event implementation
type simpleEvent struct {
	// id of the event
	id string
	// processingTime is usually the moment to process the event.
	// For instance, on an element creation, it means the creation date of that element
	processingTime time.Time
}

// Id returns the event id
func (s simpleEvent) Id() string {
	return s.id
}

// ProcessingTime returns the time to consider as the time to process the event
func (s simpleEvent) ProcessingTime() time.Time {
	return s.processingTime
}

// newSimpleEvent builds a new simple content for a given processing time
func newSimpleEvent(moment time.Time) simpleEvent {
	return simpleEvent{id: NewId(), processingTime: moment}
}

// eventContent encapsulates a content
type eventContent[C any] struct {
	// base is a simple event, we just add a content
	simpleEvent
	// content is the content to provide
	content C
}

// NewEventTick returns a new tick at now (truncated according to configuration)
func NewEventTick() Event {
	return simpleEvent{id: NewId(), processingTime: time.Now().Truncate(TIME_PRECISION)}
}

// NewEventTickTime returns a new tick at moment
func NewEventTickTime(moment time.Time) Event {
	return simpleEvent{id: NewId(), processingTime: moment}
}

// EventLifetimeEnd defines, for active elements, when to end their lifetime
type EventLifetimeEnd interface {
	// a lifetime end event is an event
	Event
	// End returns the moment to end the lifetime
	End() time.Time
}

// eventEnd uses a simple event with end lifetime = processing time
type eventEnd struct {
	// eventEnd is a simple event with a different use of its processing time
	simpleEvent
}

// End returns the moment to end the lifetime
func (e eventEnd) End() time.Time {
	return e.processingTime
}

// NewEventLifetimeEnd builds a new event to end a lifetime at given time
func NewEventLifetimeEnd(end time.Time) EventLifetimeEnd {
	return eventEnd{simpleEvent: newSimpleEvent(end)}
}

// EventStateChanges notifies a state handler that it should set those values for those attributes.
// For this particular kind of events, the processing time returns the moment to change values.
// For temporal values, it means that we end previous values at that date.
// For simple state values, it is just ignored
type EventStateChanges[T StateValue] interface {
	// this is an event
	Event
	// Changes are the changes to perform as key values
	Changes() map[string]T
}

// timedEventStateChange is a simple EventStateChanges
type simpleEventStateChange[T StateValue] eventContent[map[string]T]

// Changes returns the changes to force on the processor
func (t simpleEventStateChange[T]) Changes() map[string]T {
	return t.content
}

// NewEventStateChanges defines an event to set values since given moment
func NewEventStateChanges[T StateValue](moment time.Time, values map[string]T) EventStateChanges[T] {
	var result simpleEventStateChange[T]
	result.simpleEvent = newSimpleEvent(moment)
	result.content = values
	return result
}

// EventCreation defines an event to notify that content exists since processing time.
// Some event processors may not pay attention to processing time
type EventCreation[T Identifiable] interface {
	// creating elements is an event
	Event
	// Creation is the new content to create
	Creation() T
}

// simpleEventCreation reuses event containers
type simpleEventCreation[T Identifiable] eventContent[T]

// Creation returns the element to create
func (s simpleEventCreation[T]) Creation() T {
	return s.content
}

// NewEventCreation builds an event to create content
func NewEventCreation[T Identifiable](content T) EventCreation[T] {
	return simpleEventCreation[T]{content: content}
}

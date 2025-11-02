package commons

import (
	"time"
)

// Event is the general definition of an event:
// messages between agents, structures triggering a change of state, etc.
type Event interface {
	// each event has a unique id
	Identifiable
}

// EventProcessor defines the ability to process events
type EventProcessor interface {
	// OnEvent is called when processor receives an event
	OnEvent(Event) error
}

// simpleEvent is the most basic event implementation
type simpleEvent struct {
	// id of the event
	id string
}

// Id returns the event id
func (s simpleEvent) Id() string {
	return s.id
}

// newSimpleEvent builds a new simple content
func newSimpleEvent() simpleEvent {
	return simpleEvent{id: NewId()}
}

// EventLifetimeEnd defines, for active elements, when to end their lifetime
type EventLifetimeEnd interface {
	// a lifetime end event is an event
	Event
	// End returns the moment to end the lifetime
	End() time.Time
}

// eventEnd notifies to end a lifetime
type eventEnd struct {
	// eventEnd is a simple event with a different use of its processing time
	simpleEvent
	// moment to end a lifetime
	end time.Time
}

// End returns the moment to end the lifetime
func (e eventEnd) End() time.Time {
	return e.end
}

// NewEventLifetimeEnd builds a new event to end a lifetime at given time
func NewEventLifetimeEnd(end time.Time) EventLifetimeEnd {
	return eventEnd{simpleEvent: newSimpleEvent(), end: end}
}

// eventContent encapsulates a content
type eventContent[C any] struct {
	// base is a simple event, we just add a content
	simpleEvent
	// processingTime is the time to deal with the content
	processingTime time.Time
	// content is the content to provide
	content C
}

// ProcessingTime returns the processing time for that event
func (e eventContent[C]) ProcessingTime() time.Time {
	return e.processingTime
}

// newEventContent returns a new
func newEventContent[C any](moment time.Time, value C) eventContent[C] {
	var result eventContent[C]
	result.simpleEvent = newSimpleEvent()
	result.processingTime = moment
	result.content = value
	return result
}

// EventStateChanges notifies a state handler that it should set those values for those attributes.
// For this particular kind of events, the processing time returns the moment to change values.
// For temporal values, it means that we end previous values at that date.
// For simple state values, it is just ignored
type EventStateChanges[T StateValue] interface {
	// this is an event
	Event
	// ProcessingTime returns the moment to apply the changes
	ProcessingTime() time.Time
	// Changes are the changes to perform as key values
	Changes() map[string]T
}

// timedEventStateChange is a simple EventStateChanges
type simpleEventStateChange[T StateValue] eventContent[map[string]T]

// Changes returns the changes to force on the processor
func (t simpleEventStateChange[T]) Changes() map[string]T {
	return t.content
}

// ProcessingTime (re)implementation
func (t simpleEventStateChange[T]) ProcessingTime() time.Time {
	return t.processingTime
}

// NewEventStateChanges defines an event to set values since given moment
func NewEventStateChanges[T StateValue](moment time.Time, values map[string]T) EventStateChanges[T] {
	result := newEventContent(moment, values)
	return simpleEventStateChange[T](result)
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

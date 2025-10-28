package commons

import (
	"time"
)

// Event is the general definition of an event:
// messages between agents, structures triggering a change of state, etc.
type Event interface {
	// each event has a unique id
	Identifiable
	// Source returns the unique source of the event
	Source() ModelComponent
	// ProcessingTime returns the moment to consider the event should be processed
	ProcessingTime() time.Time
}

// Message is a directed event: a source emits it and expects it to reach a destination
type Message interface {
	// a message is an event
	Event
	// Destination returns the destination of the message
	Destination() EventProcessor
	// Content returns the event to process for destination
	Content() Event
}

// simpleEvent is the most basic event implementation
type simpleEvent struct {
	// id of the event
	id string
	// source is the component the event comes from
	source ModelComponent
	// processingTime is usually the moment to process the event.
	// For instance, on an element creation, it means the creation date of that element
	processingTime time.Time
}

// Id returns the event id
func (s simpleEvent) Id() string {
	return s.id
}

// Source returns the component asking for the creation
func (s simpleEvent) Source() ModelComponent {
	return s.source
}

// ProcessingTime returns the time to consider as the time to process the event
func (s simpleEvent) ProcessingTime() time.Time {
	return s.processingTime
}

// newSimpleEvent builds a new simple content for a given processing time and coming from a given source
func newSimpleEvent(moment time.Time, source ModelComponent) simpleEvent {
	return simpleEvent{id: NewId(), processingTime: moment, source: source}
}

// simpleMessage decorates an event to add a recipient, and then form a message
type simpleMessage struct {
	// a message is an event
	Event
	// recipient for that message
	recipient EventProcessor
	// payload is the content to process by recipient
	payload Event
}

// Destination returns the recipient for that message
func (m simpleMessage) Destination() EventProcessor {
	return m.recipient
}

// Content returns the payload as the event to process by destination
func (m simpleMessage) Content() Event {
	return m.payload
}

// NewMessage builds a new message for an event to reach its destination.
// Result has a different id from base (it is an event per se),
// but the same source and same processing time as base.
func NewMessage(base Event, destination EventProcessor) Message {
	if base == nil {
		return nil
	}

	var result simpleMessage
	result.Event = newSimpleEvent(base.ProcessingTime(), base.Source())
	result.recipient = destination
	result.payload = base
	return result
}

// eventContent encapsulates a content
type eventContent[C any] struct {
	// base is a simple event, we just add a content
	simpleEvent
	// content is the content to provide
	content C
}

// NewEventTick returns a new tick at now (truncated according to configuration)
func NewEventTick(source ModelStructure) Event {
	return simpleEvent{id: NewId(), source: source, processingTime: time.Now().Truncate(TIME_PRECISION)}
}

// NewEventTickTime returns a new tick at moment
func NewEventTickTime(source ModelStructure, moment time.Time) Event {
	return simpleEvent{id: NewId(), source: source, processingTime: moment}
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

// NewEventLifetimeEnd builds a new event to end a lifetime at given time from that structure
func NewEventLifetimeEnd(source ModelStructure, end time.Time) EventLifetimeEnd {
	return eventEnd{simpleEvent: newSimpleEvent(end, source)}
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

// NewEventStateChanges defines a source setting values since given moment
func NewEventStateChanges[T StateValue](source ModelStructure, moment time.Time, values map[string]T) EventStateChanges[T] {
	var result simpleEventStateChange[T]
	result.simpleEvent = newSimpleEvent(moment, source)
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

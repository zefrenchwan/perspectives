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

// IsEventComingFromStructure returns true if source of e is a structure
func IsEventComingFromStructure(e Event) bool {
	if e == nil {
		return false
	} else if e.Source() == nil {
		return false
	} else {
		return e.Source().GetType() == TypeStructure
	}
}

// EventProcessor processes events each time an event is received.
type EventProcessor interface {
	// Process the notified event, may emit some events or raise an error
	Process(notified Event) ([]Event, error)
}

// functionalEventProcessor is the tool to convert a function to an event processor
type functionalEventProcessor func(Event) ([]Event, error)

// Process just uses inner function to process events
func (f functionalEventProcessor) Process(event Event) ([]Event, error) {
	return f(event)
}

// NewEventProcessor builds a new event processor based on that function
func NewEventProcessor(processFn func(Event) ([]Event, error)) EventProcessor {
	if processFn == nil {
		return nil
	}

	return functionalEventProcessor(processFn)
}

// EventObserver is notified once events are received and processed from the source it listens.
// Although interface is permissive, the idea is to read events, no act on the source itself.
type EventObserver interface {
	// OnIncomingEvents is called as soon as an event is received from source.
	OnIncomingEvent(Event)
	// OnProcessingEvents is called as soon as events are processed by the source
	OnProcessingEvents([]Event, error)
}

// EventObservableProcessor is an event processer that notifies observers when it processes events
type EventObservableProcessor interface {
	// an observable processor is a processor
	EventProcessor
	// AddObserver registers a new observer to be notified
	AddObserver(EventObserver)
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

// EventCreation defines an event to notify that content exists since processing time.
// Some event processors may not pay attention to processing time
type EventCreation[T Identifiable] interface {
	// creating elements is an event
	Event
	// Content is the new content to create
	Content() T
	// CreationTime is processing time, the "birth date" of that content
	CreationTime() time.Time
}

// simpleEventCreation implements an event creation by storing fields
type simpleEventCreation[T Identifiable] struct {
	id           string
	source       ModelComponent
	content      T
	creationTime time.Time
}

// Id returns the event id
func (s simpleEventCreation[T]) Id() string {
	return s.id
}

// Source returns the component asking for the creation
func (s simpleEventCreation[T]) Source() ModelComponent {
	return s.source
}

// Content returns the content to create
func (s simpleEventCreation[T]) Content() T {
	return s.content
}

// CreationTime returns the time to consider as the content creation time
func (s simpleEventCreation[T]) CreationTime() time.Time {
	return s.creationTime
}

// NewEventCreation returns a new creation event from that source, to create content at creation time
func NewEventCreation[T Identifiable](source ModelComponent, content T, creationTime time.Time) EventCreation[T] {
	return simpleEventCreation[T]{id: NewId(), source: source, content: content, creationTime: creationTime}
}

// EventCreateLink is a link creation event
type EventCreateLink = EventCreation[Link]

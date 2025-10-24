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
	// OnEventProcessing is received by observers as soon as source processes the message.
	// Parameters are:
	// source as the event observable processor,
	// incoming as the event received by observable,
	// outgoings as the outgoing events (if any),
	// err as the raised error if any
	OnEventProcessing(source Identifiable, incoming Event, outgoings []Event, err error)
}

// EventObservableProcessor is an event processer that notifies observers when it processes events
type EventObservableProcessor interface {
	// Identifiable to know who emitted the message
	Identifiable
	// an observable processor is a processor
	EventProcessor
	// AddObserver registers a new observer to be notified
	AddObserver(EventObserver)
}

// eventObserverDecorator decorates a processor to deal with observers
type eventObserverDecorator struct {
	// id of the decorator, to implement Identifiable
	id string
	// observers are the observers to notify when a message is received or emitted
	observers []EventObserver
	// processor is the actual event processor
	processor EventProcessor
}

// Id returns the id of the processor (because it defines Process)
func (e *eventObserverDecorator) Id() string {
	return e.id
}

// AddObserver adds an observer (if not nil)
func (e *eventObserverDecorator) AddObserver(observer EventObserver) {
	if e == nil {
		return
	} else if observer != nil {
		existing := e.observers
		existing = append(existing, observer)
		existing = SliceDeduplicate(existing)
		e.observers = existing
	}
}

// Process notifies observers, actually processes the event, and notifies observers with result
func (e *eventObserverDecorator) Process(event Event) ([]Event, error) {
	if e == nil {
		return nil, nil
	}

	result, errProcessing := e.processor.Process(event)
	for _, observer := range e.observers {
		observer.OnEventProcessing(e, event, result, errProcessing)
	}

	return result, errProcessing
}

// NewEventObservableProcessor decorates a processor to become able to notify others
func NewEventObservableProcessor(processor EventProcessor) EventObservableProcessor {
	if processor == nil {
		return nil
	}

	result := new(eventObserverDecorator)
	result.observers = make([]EventObserver, 0)
	result.processor = processor
	return result
}

// eventInterceptor catches some incoming events from processor to redirect them to a catcher
type eventInterceptor struct {
	// catcher may process some events originally from processor depending on a condition
	catcher EventProcessor
	// processor is the event processor if catcher did not process it first
	processor EventProcessor
	// catchingPredicate may redirect the event to catcher (if true) or let it to processor (if false)
	catchingPredicate func(Event) bool
}

// Process picks if event will be processed by:
// catcher (if any and if predicate said yes),
// or processor (if no predicate or if predicate redirected to the processor)
func (e eventInterceptor) Process(event Event) ([]Event, error) {
	if e.catchingPredicate == nil {
		return e.processor.Process(event)
	} else if e.catchingPredicate(event) && e.catcher != nil {
		return e.catcher.Process(event)
	} else {
		return e.processor.Process(event)
	}
}

// NewEventRediction redirects an event from original to catcher if catcherAcceptance for that evebt is true.
// The event will be processed anyway, but, under given conditions, by the catcher.
// The idea is that once a state object receives an "end your lifetime" event, it cannot ignore it.
// To avoid any object implementation breaking structure invariants, structure may catch events.
// Result is then:
// catcher for no original,
// original for no catcher or no predicate,
// expected predicate for nil value,
// nil if the three of them are nil.
func NewEventRediction(catcher, original EventProcessor, catcherAcceptance func(Event) bool) EventProcessor {
	if catcher == nil || catcherAcceptance == nil {
		return original
	} else if original == nil {
		return catcher
	}

	return eventInterceptor{catcher: catcher, processor: original, catchingPredicate: catcherAcceptance}
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

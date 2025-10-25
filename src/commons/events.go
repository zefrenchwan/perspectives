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

// newSimpleContent builds a new simple content for a given processing time and coming from a given source
func newSimpleContent(moment time.Time, source ModelComponent) simpleEvent {
	return simpleEvent{id: NewId(), processingTime: moment, source: source}
}

// eventContent encapsulates a content
type eventContent[C any] struct {
	// base is a simple event, we just add a content
	simpleEvent
	// content is the content to provide
	content C
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

// EventInterceptor is the interface to implement for event interception.
// Assume a processor P expectes event E, then an interceptor will be notified
// and will execute OnRecipientProcessing(E, P) INSTEAD OF P.
// Result will be sent INSTEAD OF P.Process(E).
// Why do we do this ?
// Assume a structure that notifies an object of an "end lifetime" event.
// Code for that object may not accept or process that event.
// So, to avoid it, we regroup all "classical" event processing within an interceptor,
// and interceptor will deal with special events itself, letting object unable to act
// on its states or activity changes
type EventInterceptor interface {
	// OnRecipientProcessing catches event from recipient and returns a result.
	// Note that it is possible to call recipient.Process(event) in this function
	OnRecipientProcessing(event Event, recipient EventProcessor) ([]Event, error)
}

// eventFunctionalInterceptor implements EventInterceptor as a function call
type eventFunctionalInterceptor func(Event, EventProcessor) ([]Event, error)

// OnRecipientProcessing just calls itself
func (f eventFunctionalInterceptor) OnRecipientProcessing(event Event, recipient EventProcessor) ([]Event, error) {
	return f(event, recipient)
}

// NewEventInterceptor builds a new event interceptor decorating replacer
func NewEventInterceptor(replacer func(Event, EventProcessor) ([]Event, error)) EventInterceptor {
	if replacer == nil {
		return nil
	}

	return eventFunctionalInterceptor(replacer)
}

// NewEventRedirection redirectes events from catcher to processor based on catcherAcceptance.
// If catcherAcceptance is true for an event, then processing goes to catcher, otherwise, it goes to processor.
func NewEventRedirection(catcher, processor EventProcessor, catcherAcceptance func(e Event) bool) EventProcessor {
	if catcherAcceptance == nil || catcher == nil {
		return processor
	} else if processor == nil {
		return catcher
	}

	result := func(e Event) ([]Event, error) {
		if catcherAcceptance(e) {
			return catcher.Process(e)
		} else {
			return processor.Process(e)
		}
	}

	return NewEventProcessor(result)
}

// NewEventInterception returns a new processor built from interceptor replacing original
func NewEventInterception(original EventProcessor, interceptor EventInterceptor) EventProcessor {
	if interceptor == nil {
		return original
	} else {
		return NewEventProcessor(func(e Event) ([]Event, error) {
			return interceptor.OnRecipientProcessing(e, original)
		})
	}
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
	return eventEnd{simpleEvent: newSimpleContent(end, source)}
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
	result.simpleEvent = newSimpleContent(moment, source)
	result.content = values
	return result
}

// EventCreation defines an event to notify that content exists since processing time.
// Some event processors may not pay attention to processing time
type EventCreation[T Identifiable] interface {
	// creating elements is an event
	Event
	// Content is the new content to create
	Content() T
}

// EventCreateLink is a link creation event
type EventCreateLink = EventCreation[Link]

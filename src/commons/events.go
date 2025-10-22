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

// functionalEventProcessor is the tool to convert a function to an event processor
type functionalEventProcessor func([]Event) ([]Event, error)

// Process just uses inner function to process events
func (f functionalEventProcessor) Process(events []Event) ([]Event, error) {
	return f(events)
}

// NewEventProcessor builds a new event processor based on that function
func NewEventProcessor(processFn func([]Event) ([]Event, error)) EventProcessor {
	if processFn == nil {
		return nil
	}

	return functionalEventProcessor(processFn)
}

// EventObserver is notified once events are received and processed from the source it listens.
// Although interface is permissive, the idea is to read events, no act on the source itself.
type EventObserver interface {
	// OnIncomingEvents is called as soon as events are received from source.
	OnIncomingEvents([]Event)
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

// ObjectEventObservableProcessor is an object able to deal with events : process and be observed
type ObjectEventObservableProcessor interface {
	// ObjectEventProcessor is a model object (then may be included in a structure)
	ModelObject
	// an object event processor is able to process events and deal with observers
	EventObservableProcessor
}

// simpleEventObservableProcessor is an event processor that notified observers before and after execution
type simpleObjectEventObservableProcessor struct {
	// id returns the id of the processor
	id string
	// observers to notify (deduplicated)
	observers []EventObserver
	// decorated event processor
	decorated EventProcessor
}

// Id returns the id of an object
func (s *simpleObjectEventObservableProcessor) Id() string {
	return s.id
}

// GetType returns the object type
func (s *simpleObjectEventObservableProcessor) GetType() ModelableType {
	return TypeObject
}

// AddObserver adds a new observer to notify
func (s *simpleObjectEventObservableProcessor) AddObserver(observer EventObserver) {
	if s == nil {
		return
	} else if observer != nil {
		existing := s.observers
		existing = append(existing, observer)
		existing = SliceDeduplicate(existing)
		s.observers = existing
	}
}

// Process starts by notyfing observers, processes the events, and notifies with result.
// Performance question was raised: one loop to notify once or notify first, do and then notifies for result.
// Answer is: follow the most logical implementation and notify inputs before processing
func (s *simpleObjectEventObservableProcessor) Process(events []Event) ([]Event, error) {
	if s == nil {
		return nil, nil
	}

	for _, observer := range s.observers {
		observer.OnIncomingEvents(events)
	}

	result, errProcessing := s.decorated.Process(events)
	for _, observer := range s.observers {
		observer.OnProcessingEvents(result, errProcessing)
	}

	return result, errProcessing
}

// NewEventObservableProcessorFromProcessor decorates a processor to include observers mechanism
func NewEventObservableProcessorFromProcessor(decorated EventProcessor) ObjectEventObservableProcessor {
	if decorated == nil {
		return nil
	}

	result := new(simpleObjectEventObservableProcessor)
	result.decorated = decorated
	return result
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

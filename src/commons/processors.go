package commons

// EventObserver is notified once events are received and processed from the source it listens.
// Although interface is permissive, the idea is to read events, no act on the source itself.
type EventObserver interface {
	// an observer should have an id too to distinguish one from another
	Identifiable
	// OnEventProcessing is received by observers as soon as source processes the message.
	// Parameters are:
	// source as the event observable processor,
	// incoming as the event received by observable,
	// outgoings as the outgoing events (if any),
	// err as the raised error if any
	OnEventProcessing(receiver EventProcessor, incoming Event, outgoings []Event, err error)
}

// functionalEventObserver encapsulates a listener to implement EventObserver
type functionalEventObserver struct {
	// id of the observer
	id string
	// listener is called when source processes an event
	listener func(source EventProcessor, incoming Event, outgoings []Event, processingError error)
}

// Id returns the id of that observer
func (f functionalEventObserver) Id() string {
	return f.id
}

// OnEventProcessing is called when receiver has processed its event.
// Observer is then notified with the receiver, its incoming event, what it answered and a possible processing error.
// All elements should be read only, although it is possible to deal with them
func (f functionalEventObserver) OnEventProcessing(receiver EventProcessor, incoming Event, outgoings []Event, err error) {
	if f.listener != nil {
		f.listener(receiver, incoming, outgoings, err)
	}
}

// NewEventObserver builds an event observer from a listening function (same signature as OnEventProcessing)
func NewEventObserver(listenerFn func(source EventProcessor, incoming Event, outgoings []Event, processingError error)) EventObserver {
	if listenerFn == nil {
		return nil
	}

	return functionalEventObserver{id: NewId(), listener: listenerFn}
}

// EventProcessor processes events each time an event is received.
// It provides the opportunity to add listeners too
type EventProcessor interface {
	// A processor deals with events, we want to know who made that
	Identifiable
	// Process the notified event, may emit some events or raise an error
	Process(notified Event) ([]Event, error)
	// AddObserver registers a new observer to be notified
	AddObserver(EventObserver)
	// RemoveObserver removes obsevers matching a given predicate
	RemoveObservers(func(EventObserver) bool)
	// Observers returns the current observers
	Observers() []EventObserver
}

// EventObjectProcessor just makes an event processor AND an object
type EventObjectProcessor interface {
	// it is an object
	ModelObject
	// it is an event processor
	EventProcessor
}

// functionalEventProcessor is the tool to convert a function to an event processor
type functionalEventProcessor struct {
	// id of the current functional processor
	id string
	// processorFn is the function to use for events processing
	processorFn func(Event) ([]Event, error)
	// observers are the observers to notify when processing an event
	observers []EventObserver
}

// Id returns current processor id
func (f *functionalEventProcessor) Id() string {
	return f.id
}

// Process just uses inner function to process events
func (f *functionalEventProcessor) Process(event Event) ([]Event, error) {
	if f == nil || f.processorFn == nil {
		return nil, nil
	}

	result, errProcessing := f.processorFn(event)
	for _, observer := range f.observers {
		observer.OnEventProcessing(f, event, result, errProcessing)
	}

	return result, errProcessing
}

// AddObserver registers a new observer to be notified
func (f *functionalEventProcessor) AddObserver(observer EventObserver) {
	if f == nil {
		return
	}

	newValues := append(f.observers, observer)
	f.observers = SliceDeduplicateFunc(newValues, func(a, b EventObserver) bool { return a.Id() == b.Id() })
}

// RemoveObserver removes obsevers matching a given predicate
func (f *functionalEventProcessor) RemoveObservers(predicate func(EventObserver) bool) {
	if predicate != nil {
		f.observers = SlicesFilter(f.observers, func(o EventObserver) bool { return !predicate(o) })
	}
}

// Observers returns the current observers
func (f *functionalEventProcessor) Observers() []EventObserver {
	if f == nil {
		return nil
	}

	result := make([]EventObserver, len(f.observers))
	copy(result, f.observers)
	return result
}

// NewEventProcessor builds a new event processor based on that function
func NewEventProcessor(processFn func(Event) ([]Event, error)) EventProcessor {
	if processFn == nil {
		return nil
	}

	result := new(functionalEventProcessor)
	result.id = NewId()
	result.processorFn = processFn
	return result
}

// NewEventRedirection redirects events from catcher to processor based on catcherAcceptance.
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

// StateBasedProcessor is a function that changes its state and calculates events to produce when receiving event.
// Parameters are: state as the current state, event as the event to process.
// Result is: events to produce (may be nil), an error if any
type StateBasedProcessor[T StateValue] func(state StateHandler[T], event Event) ([]Event, error)

// StateBasedEventProcessor uses a state to process events.
// It is an object able to process events based on source and its own state.
// It may change its state depending on the events, using that object as a handler.
type StateBasedEventProcessor[T StateValue] struct {
	// current state representation
	*StateObject[T]
	// calculator calculates the next events and state changes based on current state and incoming event
	calculator StateBasedProcessor[T]
}

// Process is using the calculator based on current state
func (s *StateBasedEventProcessor[T]) Process(event Event) ([]Event, error) {
	return s.calculator(s.StateObject, event)
}

// NewStateBasedEventProcessor builds a new state based event processor with an initial state and a state event processor.
// Result is a state handler, an object, and an event processor.
func NewStateBasedEventProcessor[T StateValue](initialState map[string]T, transformer StateBasedProcessor[T]) *StateBasedEventProcessor[T] {
	if transformer == nil {
		return nil
	}

	object := NewStateObject[T]()
	object.SetValues(initialState)
	result := new(StateBasedEventProcessor[T])
	result.StateObject = object
	result.calculator = transformer
	return result
}

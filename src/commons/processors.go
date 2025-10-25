package commons

// EventProcessor processes events each time an event is received.
type EventProcessor interface {
	// A processor deals with events, we want to know who made that
	Identifiable
	// Process the notified event, may emit some events or raise an error
	Process(notified Event) ([]Event, error)
}

// functionalEventProcessor is the tool to convert a function to an event processor
type functionalEventProcessor struct {
	// id of the current functional processor
	id string
	// processorFn is the function to use for events processing
	processorFn func(Event) ([]Event, error)
}

// Id returns current processor id
func (f functionalEventProcessor) Id() string {
	return f.id
}

// Process just uses inner function to process events
func (f functionalEventProcessor) Process(event Event) ([]Event, error) {
	return f.processorFn(event)
}

// NewEventProcessor builds a new event processor based on that function
func NewEventProcessor(processFn func(Event) ([]Event, error)) EventProcessor {
	if processFn == nil {
		return nil
	}

	return functionalEventProcessor{id: NewId(), processorFn: processFn}
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
	// an observable processor is a processor
	EventProcessor
	// AddObserver registers a new observer to be notified
	AddObserver(EventObserver)
}

// eventObserverDecorator decorates a processor to deal with observers.
// EventProcessor is identifiable, so we use the same id as the original processor
type eventObserverDecorator struct {
	// EventProcess by definition
	EventProcessor
	// observers are the observers to notify when a message is received or emitted
	observers []EventObserver
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

	result, errProcessing := e.EventProcessor.Process(event)
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
	result.EventProcessor = processor
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

// NewEventInterception returns a new processor built from interceptor replacing original.
// What it does is creating an event processor that will redirects it all to the interceptor.
// NOTE THAT original will keep working as is.
// When a message arrives to the result, it goes to interceptor and only interceptor.
// Then, interceptor may decide to apply it, make some changes on original, and maybe let original
// get it back once interceptor is done.
func NewEventInterception(original EventProcessor, interceptor EventInterceptor) EventProcessor {
	if interceptor == nil {
		return original
	} else {
		return NewEventProcessor(func(e Event) ([]Event, error) {
			return interceptor.OnRecipientProcessing(e, original)
		})
	}
}

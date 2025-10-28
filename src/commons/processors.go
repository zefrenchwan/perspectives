package commons

import (
	"iter"
)

// EventProcessor processes events each time an event is received.
// It provides the opportunity to add listeners too
type EventProcessor interface {
	// A processor deals with events, we want to know who made that
	Identifiable
	// Process the notified event, may emit some events or raise an error
	Process(notified Event) ([]Event, error)
}

// EventsField defines which processors should process a given event.
// For instance, to deal with propagations, we use neighborhood approach.
// For instance, on a centralized system such as a mail, deal with messages to send to recipients.
type EventsField interface {
	// Recipients provides processors that may process that event.
	Recipients(Event) (iter.Seq[EventProcessor], bool)
}

// functionalEventProcessor is the tool to convert a function to an event processor
type functionalEventProcessor struct {
	// id of the current functional processor
	id string
	// processorFn is the function to use for events processing
	processorFn func(Event) ([]Event, error)
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

	return f.processorFn(event)
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

package commons

// Event is the general definition of an event:
// messages between agents, structures triggering a change of state, etc.
type Event interface {
	// each event has a unique id
	Identifiable
	// Source returns the unique source of the event
	Source() ModelComponent
}

// EventProcessor processes received events and may emit some events itself.
type EventProcessor interface {
	// Process the notified events, may emit some events or raise an error
	Process(notified []Event) ([]Event, error)
}

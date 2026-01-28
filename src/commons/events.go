package commons

// Event defines any impulse that would change a state
type Event interface {
}

// EventMapper processes events.
// Special cases: event tick to act, no event to emit (just change state)
type EventMapper interface {
	// OnEvents forces the processor to act, it may emit events
	OnEvents(events []Event) []Event
}

// functionaEventMapper implements EventMapper as a basic function
type functionaEventMapper func(events []Event) []Event

// OnEvents is indeed a function call
func (f functionaEventMapper) OnEvents(events []Event) []Event {
	if f == nil {
		return nil
	}

	return f(events)
}

// NewEventMapper return a new EventMapper decoring a function
func NewEventMapper(mapper func([]Event) []Event) EventMapper {
	return functionaEventMapper(mapper)
}

// NewEventIdMapper returns the id mapper
func NewEventIdMapper() EventMapper {
	return functionaEventMapper(func(events []Event) []Event {
		return events
	})
}

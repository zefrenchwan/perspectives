package commons

// Event defines any impulse that would change a state
type Event interface {
}

// EventMapper processes events.
// Special cases: event tick to act, no event to emit (just change state)
type EventMapper interface {
	// Identifiable to ensure that we may distinguish a mapper from another
	Identifiable
	// OnEvents forces the processor to act, it may emit events
	OnEvents(events []Event) []Event
}

// functionaEventMapper implements EventMapper as a basic function
type functionaEventMapper struct {
	// id of the mapper
	id string
	// processor is the decorated function
	processor func(events []Event) []Event
}

// Id to implement Identifiable
func (em *functionaEventMapper) Id() string {
	return em.id
}

// OnEvents is indeed a function call
func (em *functionaEventMapper) OnEvents(events []Event) []Event {
	if em == nil || em.processor == nil {
		return nil
	}

	return em.processor(events)
}

// NewEventMapper return a new EventMapper decoring a function
func NewEventMapper(mapper func([]Event) []Event) EventMapper {
	result := new(functionaEventMapper)
	result.id = NewId()
	result.processor = mapper
	return result
}

// NewEventIdMapper returns the id mapper
func NewEventIdMapper() EventMapper {
	return NewEventMapper(func(events []Event) []Event {
		return events
	})
}

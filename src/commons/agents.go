package commons

// StateBasedProcessor is a function that changes its state and calculates events to produce when receiving event.
// Parameters are: state as the current state, event as the event to process.
// Result is: events to produce (may be nil), an error if any
type StateBasedProcessor[T StateValue] func(state StateHandler[T], event Event) ([]Event, error)

// StateBasedAgentProcessor uses a state to process events.
// It is an object able to process events based on source and its own state.
// It may change its state depending on the events, using that object as a handler.
type StateBasedAgentProcessor[T StateValue] struct {
	// current state representation
	*StateObject[T]
	// calculator calculates the next events and state changes based on current state and incoming event
	calculator StateBasedProcessor[T]
}

// Process is using the calculator based on current state
func (s *StateBasedAgentProcessor[T]) Process(event Event) ([]Event, error) {
	return s.calculator(s.StateObject, event)
}

// NewStateBasedAgentProcessor builds a new state based agent processor with an initial state and a state event processor.
// Result is a state handler, an object, and an event processor.
func NewStateBasedAgentProcessor[T StateValue](initialState map[string]T, transformer StateBasedProcessor[T]) *StateBasedAgentProcessor[T] {
	if transformer == nil {
		return nil
	}

	object := NewStateObject[T]()
	object.SetValues(initialState)
	result := new(StateBasedAgentProcessor[T])
	result.StateObject = object
	result.calculator = transformer
	return result
}

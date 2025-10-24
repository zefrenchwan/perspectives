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
	object *StateObject[T]
	// calculator calculates the next events and state changes based on current state and incoming event
	calculator StateBasedProcessor[T]
}

// Process is using the calculator based on current state
func (s *StateBasedAgentProcessor[T]) Process(event Event) ([]Event, error) {
	return s.calculator(s.object, event)
}

// Id returns the id of the object
func (s *StateBasedAgentProcessor[T]) Id() string {
	return s.object.id
}

// GetType returns TypeObject
func (s *StateBasedAgentProcessor[T]) GetType() ModelableType {
	return TypeObject
}

// GetValue returns the value for that attribute (if any) or empty, false
func (s *StateBasedAgentProcessor[T]) GetValue(name string) (T, bool) {
	return s.object.GetValue(name)
}

// SetValue forces value for that attribute (by name)
func (s *StateBasedAgentProcessor[T]) SetValue(name string, value T) {
	s.object.SetValue(name, value)
}

// Read returns the current state of this element
func (s *StateBasedAgentProcessor[T]) Read() StateDescription[T] {
	return s.object.Read()
}

// SetValues sets values for a group of attributes
func (s *StateBasedAgentProcessor[T]) SetValues(values map[string]T) {
	s.object.SetValues(values)
}

// Remove excludes an attribute (if present).
func (s *StateBasedAgentProcessor[T]) Remove(name string) bool {
	return s.object.Remove(name)
}

// ActivePeriod returns current period of activity
func (s *StateBasedAgentProcessor[T]) ActivePeriod() Period {
	return s.object.activity
}

// SetActivePeriod forces current period of activity
func (s *StateBasedAgentProcessor[T]) SetActivePeriod(newPeriod Period) {
	s.object.SetActivePeriod(newPeriod)
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
	result.object = object
	result.calculator = transformer
	return result
}

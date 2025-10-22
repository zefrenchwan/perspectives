package commons

import "time"

// ActiveObjectHandler is a temporal state object able to deal with states.
type ActiveObjectHandler[T StateValue] interface {
	// to deal with activity and lifetime events
	TemporalHandler
	// It is a model object (then may be included in a structure)
	ModelObject
	// it is able to process events and deal with observers
	EventObservableProcessor
	// it should deal with its state
	StateHandler[T]
}

// simpleActiveObjectHandler is an active object handler decorating its content
type simpleActiveObjectHandler[T StateValue] struct {
	// id returns the id of the processor
	id string
	// activity defines the period of activity for that object
	activity Period
	// observers to notify (deduplicated)
	observers []EventObserver
	// decorated event processor
	decorated EventProcessor
	// handler to deal with state
	state StateHandler[T]
}

// Id returns the id of an object
func (s *simpleActiveObjectHandler[T]) Id() string {
	return s.id
}

// GetType returns the object type
func (s *simpleActiveObjectHandler[T]) GetType() ModelableType {
	return TypeObject
}

// ActivePeriod returns the actvity period of the object
func (s *simpleActiveObjectHandler[T]) ActivePeriod() Period {
	return s.activity
}

// SetActivePeriod changes current period
func (s *simpleActiveObjectHandler[T]) SetActivePeriod(newPeriod Period) {
	if s != nil {
		s.activity = newPeriod
	}
}

// AddObserver adds a new observer to notify
func (s *simpleActiveObjectHandler[T]) AddObserver(observer EventObserver) {
	if s == nil {
		return
	} else if observer != nil {
		existing := s.observers
		existing = append(existing, observer)
		existing = SliceDeduplicate(existing)
		s.observers = existing
	}
}

// Process starts by notyfing observers, processes the event, and notifies with result.
// Performance question was raised: one loop to notify once or notify first, do and then notifies for result.
// Answer is: follow the most logical implementation and notify inputs before processing
func (s *simpleActiveObjectHandler[T]) Process(event Event) ([]Event, error) {
	if s == nil {
		return nil, nil
	}

	for _, observer := range s.observers {
		observer.OnIncomingEvent(event)
	}

	result, errProcessing := s.decorated.Process(event)
	for _, observer := range s.observers {
		observer.OnProcessingEvents(result, errProcessing)
	}

	return result, errProcessing
}

// Read returns current state
func (s *simpleActiveObjectHandler[T]) Read() StateDescription[T] {
	if s == nil {
		return nil
	} else if s.state == nil {
		return nil
	} else {
		return s.state.Read()
	}
}

// SetValue sets value for that attribute (by name)
func (s *simpleActiveObjectHandler[T]) SetValue(name string, value T) {
	if s != nil && s.state != nil {
		s.state.SetValue(name, value)
	}
}

// SetValues sets values for a group of attributes
func (s *simpleActiveObjectHandler[T]) SetValues(values map[string]T) {
	if s != nil && s.state != nil {
		s.state.SetValues(values)
	}
}

// Remove excludes an attribute and returns if it was there
func (s *simpleActiveObjectHandler[T]) Remove(name string) bool {
	return s != nil && s.state != nil && s.state.Remove(name)
}

// NewActiveObjectHandler decorates a processor to include observers mechanism
func NewActiveObjectHandler[T StateValue](creationTime time.Time, initialState map[string]T, decorated EventProcessor) ActiveObjectHandler[T] {
	if decorated == nil {
		return nil
	}

	result := new(simpleActiveObjectHandler[T])
	result.activity = NewPeriodSince(creationTime, true)
	result.decorated = decorated
	result.id = NewId()
	result.state = NewModelStateObject[T]()
	result.state.SetValues(initialState)
	return result
}

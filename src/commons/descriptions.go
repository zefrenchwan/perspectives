package commons

import "time"

// describeState asks for a value to describe its current state
// It reads value in the content and, if possible, reads its current state.
type describeState[T StateValue] struct {
	// variable is the name of the variable to read from content
	variable string
}

// NewRequestDescription returns a RequestDescription read from a given variable
func NewRequestDescription[T StateValue](variable string) RequestDescription[T] {
	var result describeState[T]
	result.variable = variable
	return result
}

// Signature asks for that variable
func (s describeState[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{s.variable})
}

// Describe gets the state description from the content
func (s describeState[T]) Describe(c Content) StateDescription[T] {
	var value Modelable
	if c == nil {
		return nil
	} else if v, found := c.GetByName(s.variable); !found {
		return nil
	} else if v == nil {
		return nil
	} else {
		value = v
	}

	// two options:
	// either it is a state reader and we read the state
	// or it is a temporal reader and we read state now
	if r, ok := value.(StateReader[T]); ok {
		if r == nil {
			return nil
		} else {
			return r.Read()
		}
	} else if tr, ok := value.(TemporalStateReader[T]); ok {
		if tr == nil {
			return nil
		} else {
			return tr.ReadAtTime(time.Now())
		}
	}

	return nil
}

// describeTemporalState reads a value by variable and asks for its temporal state
type describeTemporalState[T StateValue] struct {
	// variable is the name of the variable to read from content
	variable string
}

// Signature asks for that variable
func (s describeTemporalState[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{s.variable})
}

// Describe gets the temporal description from the content
func (s describeTemporalState[T]) Describe(c Content) TemporalStateDescription[T] {
	if c == nil {
		return nil
	} else if v, found := c.GetByName(s.variable); !found {
		return nil
	} else if v == nil {
		return nil
	} else if d, ok := v.(TemporalStateReader[T]); !ok {
		return nil
	} else if d == nil {
		return nil
	} else {
		return d.Read()
	}
}

// NewRequestTemporalDescription builds a new request for a temporal description reading a given variable
func NewRequestTemporalDescription[T StateValue](variable string) RequestTemporalDescription[T] {
	var result describeTemporalState[T]
	result.variable = variable
	return result
}

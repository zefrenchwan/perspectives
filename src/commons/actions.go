package commons

import "maps"

// SetStateAction changes a state by forcing value to an attribute
type SetStateAction[T StateValue] struct {
	// variable is the name of expected variable
	variable string
	// changes contain the attributes and values to force
	changes map[string]T
}

// Signature returns the expected variable
func (s SetStateAction[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{s.variable})
}

// Execute runs the action by setting value for that attribute
func (s SetStateAction[T]) Execute(c Content) error {
	if c == nil {
		return nil
	} else if len(s.changes) == 0 {
		return nil
	} else if value, found := c.GetByName(s.variable); !found {
		return nil
	} else if value == nil {
		return nil
	} else if h, ok := value.(StateHandler[T]); !ok {
		return nil
	} else if h == nil {
		return nil
	} else {
		for attr, newValue := range s.changes {
			h.SetValue(attr, newValue)
		}
		return nil
	}
}

// NewSetStateAction builds an action for a single change: variable.attribute = value
func NewSetStateAction[T StateValue](variable, attribute string, value T) SetStateAction[T] {
	changes := make(map[string]T)
	changes[attribute] = value
	return SetStateAction[T]{variable: variable, changes: changes}
}

// NewSetStateActionFrom returns a change for multiple attributes.
// Parameters are the name of variable and a map of attributes and related new values
func NewSetStateActionFrom[T StateValue](variable string, changes map[string]T) SetStateAction[T] {
	result := SetStateAction[T]{}
	result.variable = variable
	result.changes = make(map[string]T)
	if len(changes) != 0 {
		maps.Copy(result.changes, changes)
	}

	return result
}

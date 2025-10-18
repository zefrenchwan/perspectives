package commons

import (
	"errors"
	"fmt"
)

// StateSetValueAction changes value for variable.attribute
type StateSetValueAction[T StateValue] struct {
	// Variable is the name of the Variable to read from content
	Variable string
	// Attribute is the name of the Attribute
	Attribute string
	// NewValue is the new value to set
	NewValue T
}

// Signature returns expected parameters
func (c StateSetValueAction[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{c.Variable})
}

// Execute runs the action on a content:
// it reads the value based on the variable, set the value for that attribute
func (s StateSetValueAction[T]) Execute(c Content) error {
	if c == nil {
		return errors.New("no content")
	} else if value, found := c.GetByName(s.Variable); !found {
		return fmt.Errorf("no variable %s in content", s.Variable)
	} else if value == nil {
		return errors.New("nil value")
	} else if handler, ok := value.(StateHandler[T]); !ok {
		return errors.New("nil value")
	} else if handler == nil {
		return errors.New("nil handler")
	} else {
		handler.SetValue(s.Attribute, s.NewValue)
		return nil
	}
}

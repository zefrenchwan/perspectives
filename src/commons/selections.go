package commons

// FilterById defines a condition for a value to match a given id.
type FilterById struct {
	id               string           // id to match
	expectedVariable string           // name of the variable to read
	parameters       FormalParameters // parameters to accept that variable at least
}

// Signature returns the expected formal parameters
func (f FilterById) Signature() FormalParameters {
	return f.parameters
}

// Matches returns true if content is identifiable with that id
func (f FilterById) Matches(c Content) (bool, error) {
	if value, found := c.GetByName(f.expectedVariable); !found {
		return false, nil
	} else if value == nil {
		return false, nil
	} else if idValue, ok := value.(Identifiable); !ok {
		return false, nil
	} else {
		return idValue.Id() == f.id, nil
	}
}

// NewFilterById returns a new condition for a variable to have a given id.
// If variable = x and expected id is "id", then condition is x.id == "id".
func NewFilterById(variable string, expectedId string) Condition {
	return FilterById{
		id:               expectedId,
		expectedVariable: variable,
		parameters:       NewNamedFormalParameters([]string{variable}),
	}
}

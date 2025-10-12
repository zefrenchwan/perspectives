package models

// Parameters define parameters of a condition or an action.
// A condition does not depend on a single entity.
// For instance, a join condition uses two entities.
// So, to regroup all cases into a general form, we use Parameters.
type Parameters interface {
	// AppendAsVariable adds a new parameter as a variable
	AppendAsVariable(name string, value ModelElement)
	// Append adds an element at the end
	Append(ModelElement)
	// Size returns the number of positional elements for that parameter.
	// It means the number of values ONLY
	Size() int
	// Get returns the positional argument for a given index.
	// If there is no value for that index, return nil
	Get(int) ModelElement
	// Variables returns the names of variables set
	Variables() []string
	// GetVariable returns the value for that variable, nil for no match
	GetVariable(name string) ModelElement
	// IsEmpty returns true if parameters are empty and should be neglected
	IsEmpty() bool
	// SelectVariables picks variables by name to make a new parameter.
	// Result contains no positional value, and variables if list if any
	SelectVariables([]string) Parameters
	// Select picks values at given indexes to make a new parameter.
	// Result contains only positional values, with selected indexes (if any)
	Select([]int) Parameters
	// PositionalParameters returns the positional parameters as a slice
	PositionalParameters() []ModelElement
	// NamedParameters returns the named parameters as a map
	NamedParameters() map[string]ModelElement
}

// genericParameters defines a basic implementation
// as an array for positional elements
// as a map for named elements
type genericParameters struct {
	// positionals contain the positional arguments
	positionals []ModelElement
	// named contains named arguments
	named map[string]ModelElement
}

// Size returns the number of values in that parameter
func (a *genericParameters) Size() int {
	if a == nil {
		return 0
	}

	return len(a.positionals)
}

// IsEmpty is true for no element
func (a *genericParameters) IsEmpty() bool {
	return a == nil || (len(a.positionals) == 0 && len(a.named) == 0)
}

// Append adds a positional parameter
func (a *genericParameters) Append(element ModelElement) {
	if a != nil {
		a.positionals = append(a.positionals, element)
	}
}

// AppendAsVariable adds a new named value as a variable
func (a *genericParameters) AppendAsVariable(name string, value ModelElement) {
	if a != nil {
		if a.named == nil {
			a.named = make(map[string]ModelElement)
		}

		a.named[name] = value
	}
}

// Get returns the value at a given position, or nil if index does not match
func (a *genericParameters) Get(index int) ModelElement {
	if a == nil || index < 0 || index >= len(a.positionals) {
		return nil
	}

	return a.positionals[index]
}

// Variables returns the name of variables set for that parameter
func (a *genericParameters) Variables() []string {
	if a == nil {
		return nil
	}

	var result []string
	for name := range a.named {
		result = append(result, name)
	}

	return result
}

// GetVariable returns the value (if any) for that name
func (a *genericParameters) GetVariable(name string) ModelElement {
	if a == nil {
		return nil
	}

	return a.named[name]
}

// SelectVariables reduces this parameters to only matching variables
func (a *genericParameters) SelectVariables(names []string) Parameters {
	if a == nil {
		return nil
	}

	result := new(genericParameters)
	if len(a.named) == 0 {
		return result
	}

	result.named = make(map[string]ModelElement)
	for _, name := range names {
		if value, found := a.named[name]; found {
			result.named[name] = value
		}
	}

	return result
}

// Select reduces this parameters to only matching indexes
func (a *genericParameters) Select(indexes []int) Parameters {
	if a == nil {
		return nil
	}

	result := new(genericParameters)
	size := len(a.positionals)
	if size == 0 {
		return result
	}

	for _, index := range indexes {
		if index >= 0 && index < size {
			result.positionals = append(result.positionals, a.positionals[index])
		}
	}

	return result
}

// PositionalParameters returns the positional parameters as a slice
func (a *genericParameters) PositionalParameters() []ModelElement {
	var result []ModelElement
	if a == nil {
		return result
	}

	result = append(result, a.positionals...)

	return result
}

// NamedParameters returns the named parameters as a map
func (a *genericParameters) NamedParameters() map[string]ModelElement {
	if a == nil {
		return nil
	}

	result := make(map[string]ModelElement)
	for name, value := range a.named {
		result[name] = value
	}

	return result
}

// NewParameter returns a new parameter for a single element
func NewParameter(element ModelElement) Parameters {
	result := new(genericParameters)
	result.positionals = append(result.positionals, element)
	return result
}

// NewNamedParameter returns a new parameter for a single named element
func NewNamedParameter(name string, element ModelElement) Parameters {
	result := new(genericParameters)
	result.named = make(map[string]ModelElement)
	result.named[name] = element
	return result
}

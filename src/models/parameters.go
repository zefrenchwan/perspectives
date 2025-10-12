package models

// Parameters define parameters of a condition or an action.
// A condition does not depend on a single entity.
// For instance, a join condition uses two entities.
// So, to regroup all cases into a general form, we use Parameters.
type Parameters[T ModelElement] interface {
	// AppendAsVariable adds a new parameter as a variable
	AppendAsVariable(name string, value T)
	// Append adds an element at the end
	Append(T)
	// Size returns the number of positional elements for that parameter.
	// It means the number of values ONLY
	Size() int
	// Get returns the positional argument for a given index.
	// If there is no value for that index, return nil
	Get(int) T
	// Variables returns the names of variables set
	Variables() []string
	// GetVariable returns the value for that variable, nil for no match
	GetVariable(name string) T
	// IsEmpty returns true if parameters are empty and should be neglected
	IsEmpty() bool
	// SelectVariables picks variables by name to make a new parameter.
	// Result contains no positional value, and variables if list if any
	SelectVariables([]string) Parameters[T]
	// Select picks values at given indexes to make a new parameter.
	// Result contains only positional values, with selected indexes (if any)
	Select([]int) Parameters[T]
	// PositionalParameters returns the positional parameters as a slice
	PositionalParameters() []T
	// NamedParameters returns the named parameters as a map
	NamedParameters() map[string]T
}

// genericParameters defines a basic implementation
// as an array for positional elements
// as a map for named elements
type genericParameters[T ModelElement] struct {
	// positionals contain the positional arguments
	positionals []T
	// named contains named arguments
	named map[string]T
}

// Size returns the number of values in that parameter
func (a *genericParameters[T]) Size() int {
	if a == nil {
		return 0
	}

	return len(a.positionals)
}

// IsEmpty is true for no element
func (a *genericParameters[T]) IsEmpty() bool {
	return a == nil || (len(a.positionals) == 0 && len(a.named) == 0)
}

// Append adds a positional parameter
func (a *genericParameters[T]) Append(element T) {
	if a != nil {
		a.positionals = append(a.positionals, element)
	}
}

// AppendAsVariable adds a new named value as a variable
func (a *genericParameters[T]) AppendAsVariable(name string, value T) {
	if a != nil {
		if a.named == nil {
			a.named = make(map[string]T)
		}

		a.named[name] = value
	}
}

// Get returns the value at a given position, or nil if index does not match
func (a *genericParameters[T]) Get(index int) T {
	var empty T
	if a == nil || index < 0 || index >= len(a.positionals) {
		return empty
	}

	return a.positionals[index]
}

// Variables returns the name of variables set for that parameter
func (a *genericParameters[T]) Variables() []string {
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
func (a *genericParameters[T]) GetVariable(name string) T {
	var empty T
	if a == nil {
		return empty
	}

	return a.named[name]
}

// SelectVariables reduces this parameters to only matching variables
func (a *genericParameters[T]) SelectVariables(names []string) Parameters[T] {
	if a == nil {
		return nil
	}

	result := new(genericParameters[T])
	if len(a.named) == 0 {
		return result
	}

	result.named = make(map[string]T)
	for _, name := range names {
		if value, found := a.named[name]; found {
			result.named[name] = value
		}
	}

	return result
}

// Select reduces this parameters to only matching indexes
func (a *genericParameters[T]) Select(indexes []int) Parameters[T] {
	if a == nil {
		return nil
	}

	result := new(genericParameters[T])
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
func (a *genericParameters[T]) PositionalParameters() []T {
	var result []T
	if a == nil {
		return result
	}

	result = append(result, a.positionals...)

	return result
}

// NamedParameters returns the named parameters as a map
func (a *genericParameters[T]) NamedParameters() map[string]T {
	if a == nil {
		return nil
	}

	result := make(map[string]T)
	for name, value := range a.named {
		result[name] = value
	}

	return result
}

// NewParameter returns a new parameter for a single element
func NewParameter[T ModelElement](element T) Parameters[T] {
	result := new(genericParameters[T])
	result.positionals = append(result.positionals, element)
	return result
}

// NewNamedParameter returns a new parameter for a single named element
func NewNamedParameter[T ModelElement](name string, element T) Parameters[T] {
	result := new(genericParameters[T])
	result.named = make(map[string]T)
	result.named[name] = element
	return result
}

package commons

// Content defines any values grouped together.
// Typically, it may be used as parameters to run conditions.
// A condition does not depend on a single entity.
// For instance, a join condition uses two entities.
// So, to regroup all cases into a general form, we use Content.
type Content interface {
	// AppendAsVariable adds a new value as a variable
	AppendAsVariable(name string, value ModelComponent)
	// Append adds an element at the end
	Append(ModelComponent)
	// Size returns the number of positional elements for that content.
	// It means the number of positional values, no matter the variable content
	Size() int
	// Get returns the positional argument for a given index.
	// If there is no value for that index, return nil
	Get(int) ModelComponent
	// Variables returns the names of variables set
	Variables() []string
	// GetVariable returns the value for that variable, nil for no match
	GetVariable(name string) ModelComponent
	// IsEmpty returns true if content is empty and should be neglected
	IsEmpty() bool
	// SelectVariables picks variables by name to make a new content.
	// Result contains no positional value, and variables if list if any
	SelectVariables([]string) Content
	// Select picks values at given indexes to make a new content.
	// Result contains only positional values, with selected indexes (if any)
	Select([]int) Content
	// PositionalContent returns the positional content as a slice
	PositionalContent() []ModelComponent
	// NamedContent returns the named content as a map
	NamedContent() map[string]ModelComponent
	// Unique picks first element, if content contains EITHER one positional value, OR one single named value.
	// It returns nil, false if there are too many elements or if content is empty
	Unique() (ModelComponent, bool)
}

// simpleContainer defines a basic implementation
// as an array for positional elements
// as a map for named elements
type simpleContainer struct {
	// positionals contain the positional values
	positionals []ModelComponent
	// named contains named values
	named map[string]ModelComponent
}

// Size returns the number of positional values in that content
func (a *simpleContainer) Size() int {
	if a == nil {
		return 0
	}

	return len(a.positionals)
}

// IsEmpty is true for no element
func (a *simpleContainer) IsEmpty() bool {
	return a == nil || (len(a.positionals) == 0 && len(a.named) == 0)
}

// Append adds a positional value in that content
func (a *simpleContainer) Append(element ModelComponent) {
	if a != nil {
		a.positionals = append(a.positionals, element)
	}
}

// AppendAsVariable adds a new named value as a variable
func (a *simpleContainer) AppendAsVariable(name string, value ModelComponent) {
	if a != nil {
		if a.named == nil {
			a.named = make(map[string]ModelComponent)
		}

		a.named[name] = value
	}
}

// Get returns the value at a given position, or nil if index does not match
func (a *simpleContainer) Get(index int) ModelComponent {
	if a == nil || index < 0 || index >= len(a.positionals) {
		return nil
	}

	return a.positionals[index]
}

// Variables returns the name of variables set for that content
func (a *simpleContainer) Variables() []string {
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
func (a *simpleContainer) GetVariable(name string) ModelComponent {
	if a == nil {
		return nil
	}

	return a.named[name]
}

// SelectVariables reduces this content to only matching variables
func (a *simpleContainer) SelectVariables(names []string) Content {
	if a == nil {
		return nil
	}

	result := new(simpleContainer)
	if len(a.named) == 0 {
		return result
	}

	result.named = make(map[string]ModelComponent)
	for _, name := range names {
		if value, found := a.named[name]; found {
			result.named[name] = value
		}
	}

	return result
}

// Select reduces this content to only matching indexes
func (a *simpleContainer) Select(indexes []int) Content {
	if a == nil {
		return nil
	}

	result := new(simpleContainer)
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

// PositionalContent returns the positional content as a slice
func (a *simpleContainer) PositionalContent() []ModelComponent {
	var result []ModelComponent
	if a == nil {
		return result
	}

	result = append(result, a.positionals...)

	return result
}

// NamedContent returns the named content as a map
func (a *simpleContainer) NamedContent() map[string]ModelComponent {
	if a == nil {
		return nil
	}

	result := make(map[string]ModelComponent)
	for name, value := range a.named {
		result[name] = value
	}

	return result
}

// Unique picks the only value if any, or nil false
func (a *simpleContainer) Unique() (ModelComponent, bool) {
	if a == nil {
		return nil, false
	}

	if len(a.positionals) == 1 {
		if len(a.named) == 0 {
			return a.positionals[0], true
		} else {
			return nil, false
		}
	} else if len(a.positionals) == 0 {
		if len(a.named) == 1 {
			for _, value := range a.named {
				return value, true
			}
		}
	}

	return nil, false

}

// NewContent returns a new content for a single element
func NewContent(element ModelComponent) Content {
	result := new(simpleContainer)
	result.positionals = append(result.positionals, element)
	return result
}

// NewNamedContent returns a new content for a single named element
func NewNamedContent(name string, element ModelComponent) Content {
	result := new(simpleContainer)
	result.named = make(map[string]ModelComponent)
	result.named[name] = element
	return result
}

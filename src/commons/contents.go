package commons

import "maps"

// GenericContent defines any values grouped together.
// Typically, it may be used as parameters to run conditions.
// A condition does not depend on a single entity:
// for instance, a join condition uses two entities.
// So, to regroup all cases into a general form, we use GenericContent.
// There are two types of values:
// POSITIONAL values: no name, just values one after another
// NAMED values: a map of names and values
type GenericContent[T any] interface {
	// Matches tests if content matches expected parameters:
	// Enough total values, enough positional values, and enough variables.
	// Failure to comply to any condition returns false
	Matches(fp FormalParameters) bool
	// AppendAs adds a new value with a name
	AppendAs(name string, value T)
	// Append adds an element at the end
	Append(T)
	// Size returns the number of positional elements for that content.
	// It means the number of positional values, no matter the named content
	Size() int
	// Disjoin returns true if there is no common allocated place.
	// That is, two contents are disjoin when:
	// names a different from one another (no matter the value)
	// AND if one has positional values, the other does not
	Disjoin(GenericContent[T]) bool
	// Get returns the positional argument for a given index.
	// If there is no value for that index, return nil, false
	Get(int) (T, bool)
	// Names returns the names set for some values.
	// For intance, if content is x => v1, y => v2, then result is x,y
	Names() []string
	// GetByName returns the value for that name, nil, false for no match.
	// Result is then the match (if any), and a bool indicating if there was a match
	GetByName(name string) (T, bool)
	// IsEmpty returns true if content is empty and should be neglected
	IsEmpty() bool
	// PickByNames picks values by name to make a new content.
	// Result contains no positional value, and named values if list if any
	PickByNames([]string) GenericContent[T]
	// PickByIndexes picks values at given indexes to make a new content.
	// Result contains only positional values, with selected indexes (if any)
	PickByIndexes([]int) GenericContent[T]
	// PositionalContent returns the positional content as a slice
	PositionalContent() []T
	// MapPositionalsToNamed gets names in given order and gets positional parameters to make named content.
	// For instance, if values are [a,b,c] and we want []{"x","y"}, then result is "x" => a, "y" => b.
	// If there is not enough positional values, then return empty, false.
	MapPositionalsToNamed(names []string) (GenericContent[T], bool)
	// MapNamedToPositionals returns positional content in order of named values.
	// If there is not enough positional values, it returns nil, false.
	// For instance, if named values are "x" => a, "y" => b, and names is "y","x"
	// then result should be (b,a)
	MapNamedToPositionals(names []string) (GenericContent[T], bool)
	// NamedContent returns the named content as a map
	NamedContent() map[string]T
	// Unique picks first element, if content contains EITHER one positional value, OR one single named value.
	// It returns nil, false if there are too many elements or if content is empty
	Unique() (T, bool)
}

type Content GenericContent[Modelable]

// simpleContainer defines a basic implementation
// as an array for positional elements
// as a map for named elements
type simpleContainer[T any] struct {
	// positionals contain the positional values
	positionals []T
	// named contains named values
	named map[string]T
}

// Matches tests if content matches expected parameters
func (a *simpleContainer[T]) Matches(fp FormalParameters) bool {
	size := 0
	var variables []string

	if a != nil {
		size = a.Size()
		variables = a.Names()
	}

	if fp.minimalPositionalSize > size {
		return false
	} else if len(fp.expectedVariables) > len(variables) {
		return false
	} else {
		return SlicesContainsAllFunc(variables, fp.expectedVariables, func(a, b string) bool { return a == b })
	}
}

// Size returns the number of positional values in that content
func (a *simpleContainer[T]) Size() int {
	if a == nil {
		return 0
	}

	return len(a.positionals)
}

// IsEmpty is true for no element
func (a *simpleContainer[T]) IsEmpty() bool {
	return a == nil || (len(a.positionals) == 0 && len(a.named) == 0)
}

// Disjoin returns true if
// other and a have no variable in common
// and if one has positional values, the other does not
func (a *simpleContainer[T]) Disjoin(other GenericContent[T]) bool {
	if a == nil {
		return true
	} else if other == nil {
		return true
	}

	if a.Size() != 0 && other.Size() != 0 {
		return false
	}

	return !SliceCommonElement(a.Names(), other.Names())
}

// Append adds a positional value in that content
func (a *simpleContainer[T]) Append(element T) {
	if a != nil {
		a.positionals = append(a.positionals, element)
	}
}

// AppendAs adds a new named value as a variable
func (a *simpleContainer[T]) AppendAs(name string, value T) {
	if a != nil {
		if a.named == nil {
			a.named = make(map[string]T)
		}

		a.named[name] = value
	}
}

// Get returns the value at a given position, or nil if index does not match
func (a *simpleContainer[T]) Get(index int) (T, bool) {
	var empty T
	if a == nil || index < 0 || index >= len(a.positionals) {
		return empty, false
	}

	return a.positionals[index], true
}

// Names returns the name of named values set for that content
func (a *simpleContainer[T]) Names() []string {
	if a == nil {
		return nil
	}

	var result []string
	for name := range a.named {
		result = append(result, name)
	}

	return result
}

// GetByName returns the value (if any) for that name
func (a *simpleContainer[T]) GetByName(name string) (T, bool) {
	var empty T
	if a == nil {
		return empty, false
	} else if value, found := a.named[name]; !found {
		return empty, false
	} else {
		return value, true
	}
}

// PickByNames reduces this content to only matching named values by names
func (a *simpleContainer[T]) PickByNames(names []string) GenericContent[T] {
	if a == nil {
		return nil
	}

	result := new(simpleContainer[T])
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

// PickByIndexes reduces this content to only matching indexes
func (a *simpleContainer[T]) PickByIndexes(indexes []int) GenericContent[T] {
	if a == nil {
		return nil
	}

	result := new(simpleContainer[T])
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
func (a *simpleContainer[T]) PositionalContent() []T {
	var result []T
	if a == nil {
		return result
	}

	result = append(result, a.positionals...)

	return result
}

// NamedContent returns the named content as a map
func (a *simpleContainer[T]) NamedContent() map[string]T {
	if a == nil {
		return nil
	}

	result := make(map[string]T)
	for name, value := range a.named {
		result[name] = value
	}

	return result
}

// Unique picks the only value if any, or nil false
func (a *simpleContainer[T]) Unique() (T, bool) {
	var empty T
	if a == nil {
		return empty, false
	}

	if len(a.positionals) == 1 {
		if len(a.named) == 0 {
			return a.positionals[0], true
		} else {
			return empty, false
		}
	} else if len(a.positionals) == 0 {
		if len(a.named) == 1 {
			for _, value := range a.named {
				return value, true
			}
		}
	}

	return empty, false

}

// MapPositionalsToNamed gets names in given order and gets positional parameters to make named content.
func (a *simpleContainer[T]) MapPositionalsToNamed(names []string) (GenericContent[T], bool) {
	result := new(simpleContainer[T])
	result.named = make(map[string]T)

	if a == nil {
		return nil, false
	} else if len(names) == 0 {
		return result, true
	} else if len(names) > a.Size() {
		return nil, false
	}

	// sizes fit, so just copy.
	// Because a.Size >= len(names) >= 1, then positionals is not nil
	for index := 0; index < len(names); index++ {
		name := names[index]
		value := a.positionals[index]
		result.named[name] = value
	}

	return result, true
}

// MapNamedToPositionals returns positional content in order of named values.
func (a *simpleContainer[T]) MapNamedToPositionals(names []string) (GenericContent[T], bool) {
	result := new(simpleContainer[T])
	result.named = make(map[string]T)

	if a == nil {
		return nil, false
	} else if len(names) == 0 {
		return result, true
	} else if len(names) > len(a.named) {
		return nil, false
	}

	result.positionals = make([]T, 0)
	for _, name := range names {
		if value, found := a.named[name]; !found {
			return nil, false
		} else {
			result.positionals = append(result.positionals, value)
		}
	}

	return result, true
}

// NewContent returns a new content for a single element
func NewContent[T any](element T) GenericContent[T] {
	result := new(simpleContainer[T])
	result.positionals = append(result.positionals, element)
	return result
}

// NewNamedContent returns a new content for a single named element
func NewNamedContent[T any](name string, element T) GenericContent[T] {
	result := new(simpleContainer[T])
	result.named = make(map[string]T)
	result.named[name] = element
	return result
}

// NewNamedContentFromMap reads a map as named content
func NewNamedContentFromMap[T any](values map[string]T) GenericContent[T] {
	result := new(simpleContainer[T])
	result.named = make(map[string]T)
	if values != nil {
		maps.Copy(result.named, values)
	}

	return result
}

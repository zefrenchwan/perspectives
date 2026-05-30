package commons

import "slices"

// Variable is used as a placeholder in a Link, with type constraints.
type Variable struct {
	name         string  // name of the variable (same name => same variable)
	allowedTypes []Class // allowedTypes are possible types for replacement
}

// NewVariable creates an immutable variable with replacement constraints.
func NewVariable(name string, allowedTypes ...Class) Variable {
	constraints := make([]Class, len(allowedTypes))
	copy(constraints, allowedTypes)

	return Variable{
		name:         name,
		allowedTypes: constraints,
	}
}

// Name of the variable
func (v Variable) Name() string {
	return v.name
}

// AllowedTypes returns possible types for replacement
func (v Variable) AllowedTypes() []Class {
	return slices.Clone(v.allowedTypes)
}

// CanBeReplacedBy checks if an element can replace this variable
func (v Variable) CanBeReplacedBy(element Element) bool {
	if element == nil {
		return false
	}

	if len(v.allowedTypes) == 0 {
		return true
	}
	return slices.Contains(v.allowedTypes, element.DeclaringClass())
}

// Same returns true if other is a variable with the same name
func (v Variable) Same(other Element) bool {
	if other == nil {
		return false
	}
	if other.DeclaringClass() != CLASS_VARIABLE {
		return false
	}

	if otherVar, ok := other.(Variable); !ok {
		return false
	} else {
		return v.name == otherVar.name
	}

	return false
}

// DeclaringClass returns that this is a variable
func (v Variable) DeclaringClass() Class {
	return CLASS_VARIABLE
}

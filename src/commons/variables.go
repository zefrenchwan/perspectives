package commons

import "slices"

// Variable is used as a placeholder in a Link, with type constraints.
type Variable struct {
	name         string  // name of the variable (same name => same variable)
	allowedTypes []Class // allowedTypes are possible types for replacement
}

// NewVariable creates an immutable variable with replacement constraints.
func NewVariable(name string, allowedTypes ...Class) *Variable {
	constraints := make([]Class, len(allowedTypes))
	copy(constraints, allowedTypes)

	return &Variable{
		name:         name,
		allowedTypes: constraints,
	}
}

// Id is the name
func (v *Variable) Id() string {
	if v == nil {
		return ""
	}
	return v.name
}

// Name of the variable
func (v *Variable) Name() string {
	if v == nil {
		return ""
	}
	return v.name
}

// AllowedTypes returns possible types for replacement
func (v *Variable) AllowedTypes() []Class {
	if v == nil {
		return nil
	}
	return slices.Clone(v.allowedTypes)
}

// CanBeReplacedBy checks if an element can replace this variable
func (v *Variable) CanBeReplacedBy(element Element) bool {
	if v == nil || element == nil {
		return false
	}

	if len(v.allowedTypes) == 0 {
		return true
	}

	elementClasses := element.DeclaringClasses()
	for _, allowed := range v.allowedTypes {
		if slices.Contains(elementClasses, allowed) {
			return true
		}
	}

	return false
}

// Same returns true if other is a variable with the same name
func (v *Variable) Same(other Element) bool {
	if v == nil && other == nil {
		return true
	}
	if v == nil || other == nil {
		return false
	}
	return v.Id() == other.Id()
}

// DeclaringClasses returns that this is a variable
func (v *Variable) DeclaringClasses() []Class {
	return []Class{CLASS_VARIABLE}
}

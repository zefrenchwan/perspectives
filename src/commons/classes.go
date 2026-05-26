package commons

import "slices"

// Class is the general definition of elements within the system.
// It applies to any element that can be declared with a specific class.
// Class is NOT a declaration type, but rather a categorization of elements.
// To use declaration types, see the traits mechanism.
type Class string

const CLASS_TRAIT Class = "trait"
const CLASS_LINK Class = "link"
const CLASS_INSTANCE Class = "instance"
const CLASS_VARIABLE Class = "variable"

// Element is a system entity.
// For instance, traits, graphs, links, etc.
type Element interface {
	Id() string                // Id returns the unique identifier of the element
	Same(other Element) bool   // Same checks if two elements are functionally equivalent
	DeclaringClasses() []Class // DeclaringClasses returns the classes for that element. If a class is declared, casting should be possible
}

// IsElementDeclaredInstance checks if an element is declared with a specific class.
// For instance, a link should declare CLASS_LINK
func IsElementDeclaredInstance(element Element, c Class) bool {
	if element == nil {
		return false
	}
	return slices.Contains(element.DeclaringClasses(), c)
}

package objects

// Class is the general definition of elements within the system.
// It applies to any element that can be declared with a specific class.
// Class is NOT a declaration type, but rather a categorization of elements.
// To use declaration types, see the traits mechanism.
type Class string

const CLASS_TRAIT Class = "trait"
const CLASS_LINK Class = "link"
const CLASS_INSTANCE Class = "instance"
const CLASS_VARIABLE Class = "variable"

// Element is a system entity : for instance, traits, instances, links, etc.
type Element interface {
	Same(other Element) bool // Same checks if two elements are functionally equivalent
	DeclaringClass() Class   // DeclaringClass returns the class for that element.
}

// IdentifiableElement is an element that may be distinguished by an identifier.
// It is unique, meaning that one may decide whether two elements are the same based on their identifier.
// When implementing IdentifiableElement, ensure that the Same function is consistent with the identifier.
// Typical implementations should then be same if id and classes match
type IdentifiableElement interface {
	Element     // Element to ensure that classes logic applies
	Id() string // Id returns the unique identifier for the element
}

// IsElementDeclaredInstance checks if an element is declared with a specific class.
// For instance, a link should declare CLASS_LINK
func IsElementDeclaredInstance(element Element, c Class) bool {
	if element == nil {
		return false
	}
	return element.DeclaringClass() == c
}

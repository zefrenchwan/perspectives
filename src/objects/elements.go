package objects

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

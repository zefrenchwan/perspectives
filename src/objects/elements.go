package objects

// Element is a system entity : for instance, traits, instances, links, etc.
// It has an identifier and Same implementation has to be consistent with the identifier.
type Element interface {
	// Id returns the unique identifier for the element
	Id() string
	// Same checks if two elements are functionally equivalent
	Same(other Element) bool
	// DeclaringClass returns the class for that element.
	DeclaringClass() Class
}

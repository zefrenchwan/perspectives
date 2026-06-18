package objects

// Element is a system entity : for instance, instances, links, etc.
// It has an identifier and Same implementation has to be consistent with the identifier.
// It should be :
// EITHER A LINK (link composition)
// OR AN INSTANCE (instance as operand) OR A SET OF INSTANCES (as a collection).
// To do so, we use the sealed interface :
// we include a private function to force that no other type can implement it.
// This way, element are linkable types and can only be implemented within this package.
// VERY IMPORTANT : Element should be immutable.
type Element interface {
	// Id returns the unique identifier for the element
	Id() string
	// Same checks if two elements are functionally equivalent
	Same(other Element) bool
	// DeclaringClass returns the class for that element.
	DeclaringClass() Class
	// we use the sealed interface to manage linkable types
	isLinkable() bool
	// we use toHashString to manage sames without full link walks
	toHashString() string
}

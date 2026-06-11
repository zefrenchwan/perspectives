package objects

import "time"

// Element is a system entity : for instance, traits, instances, links, etc.
type Element interface {
	// Same checks if two elements are functionally equivalent
	Same(other Element) bool
	// DeclaringClass returns the class for that element.
	DeclaringClass() Class
}

// Observable is an element that may be distinguished by an identifier and we may observe at a given time.
// It is unique, meaning that one may decide whether two elements are the same based on their identifier.
// When implementing Observable, ensure that the Same function is consistent with the identifier.
// Typical implementations should then be same if id and classes match
type Observable[T any] interface {
	Element // Element to ensure that classes logic applies
	// Id returns the unique identifier for the element
	Id() string
	// Observe returns the value of the observable at a given time
	Observe(time.Time) T
}

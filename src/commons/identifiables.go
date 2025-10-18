package commons

// Identifiable defines anything that has an id.
// An id should be globally unique no matter the type:
// for instance, no link should share the same id as an object.
// A model element has an id if any observer may distinguish it from another.
type Identifiable interface {
	// Id returns the id of that element.
	Id() string
}

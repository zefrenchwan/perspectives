package models

// Field is the definition of the structural part of a model.
// "Field" defines a field for entities
type Field interface {
	// A field has an id
	IdentifiableElement
	// Apply an operation right now
	Apply(Operation) (bool, error)
}

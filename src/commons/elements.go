package commons

import "github.com/google/uuid"

// ModelElement is the most general element within a model
type ModelElement interface{}

// IdentifiableElement defines anything that has an id.
// An id should be globally unique : no link should have the same id as an object.
// A model element has an id if any observer may distinguish it from another.
type IdentifiableElement interface {
	Id() string // Id returns the id of that entity.
}

// NewId builds a new unique id.
// Two different calls should return two different values.
func NewId() string {
	return uuid.NewString()
}

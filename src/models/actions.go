package models

// Action is the general possible action.
// An action may apply to:
// An entity (for instance, end lifetime)
type Action interface{}

// EntityTransformation applies to an entity and transforms it.
// It is indeed a one to one mapping.
// Result is then:
// New version of the entity
// A boolean to be considered first: did the content change ?
// An error just in case
type EntityTransformation interface {
	// A transformation is an action
	Action
	// Apply maps an entity to another if there is a change.
	// Result is the new version (if there was a change), change bool, and an error if any
	Apply(ModelEntity) (ModelEntity, bool, error)
}

// EntityCreation builds a new instance
type EntityCreation struct {
	// Element is the new element
	Element ModelEntity
}

// EntityDeletion defines the deletion of an entity.
// Definition is just its type.
// Note that a deletion is NOT ending the lifetime of a temporal entity.
// It is really getting rid of an entity within a field.
type EntityDeletion struct {
}

package models

import "time"

// EventAction is the most generic definition of an event to change a system
type EventAction interface{}

// EventCreation defines the creation of a new entity
type EventCreation struct {
	EventAction             // A creation is an action
	Value       ModelEntity // the new value to add
}

// EventDeletion defines the deletion of an entity.
// It means that the entity is deleted, not and end of its lifetime.
// To change a lifetime, use an EventChange.
type EventDeletion interface {
	EventAction               // A deletion is an event
	Matches(ModelEntity) bool // Matches returns true if the entity should be deleted
}

// EventChange defines a modification of an entity.
// It applies to the entity only.
type EventChange interface {
	Matches(ModelEntity) bool      // Matches returns true if change applies to this entity
	Apply(ModelEntity) ModelEntity // Apply returns the new version of the entity
}

// Event is the container for an action.
// It contains an id (unique), metadata and the action per se.
type Event struct {
	Id           string            // unique id of the event
	CreationTime time.Time         // event creation time
	SourceId     string            // id of the source
	Metadata     map[string]string // general metadata definition
	Action       EventAction       // the action to execute
}

// EventDeletionById defines an entity to delete by id.
// Definition with an explicit, public id is the best solution for performance issues.
// Reason is we really want to avoid a full scan of a system to test a match.
// No matter the actual implementation of a storage system, id access should be super fast.
// So, using an id as an explicit public value means going directly to the entity.
// It uses extensively the fact that an id is globally unique.
type EventDeletionById struct {
	Id string // Id of the object to delete
}

// Matches implements the event deletion interface
func (e EventDeletionById) Matches(me ModelEntity) bool {
	if identifiable, ok := me.(IdentifiableEntity); !ok {
		return false
	} else {
		return identifiable.Id() == e.Id
	}
}

package models

import (
	"maps"
	"time"
)

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
	Metadata     map[string]string // general metadata definition
	Action       EventAction       // the action to execute
}

// IsEmpty returns true for an empty event, and should be ignored
func (e Event) IsEmpty() bool {
	return e.Id == "" || e.Action == nil
}

// NewEvent builds an event for an action with provided metadata.
// If there is no action (set to nil) then the result is an empty event, with no id
func NewEvent(metadata map[string]string, action EventAction) Event {
	if action == nil {
		return Event{}
	}

	var md map[string]string
	if len(metadata) != 0 {
		md = make(map[string]string)
		maps.Copy(md, metadata)
	}

	return Event{
		Id:           NewId(),
		CreationTime: time.Now(),
		Metadata:     md,
		Action:       action,
	}
}

// NewSimpleEvent builds the simplest event based on an action
func NewSimpleEvent(action EventAction) Event {
	return NewEvent(nil, action)
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

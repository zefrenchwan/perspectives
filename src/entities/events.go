package entities

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Event is any change of state in the system.
// It is an immutable object.
type Event interface {
	// Identifiable to get the event id
	commons.Identifiable
	// CreationDate to get the date when the event was created
	CreationDate() time.Time
}

// baseEvent is the event skeleton to manage common event properties :
// event id and creation date
type baseEvent struct {
	// eventId is the unique identifier of the event
	eventId string
	// creationDate is the date when the event was created
	creationDate time.Time
}

// Id returns the immutable id of the event
func (e baseEvent) Id() string {
	return e.eventId
}

// CreationDate returns the date when the event was created
func (e baseEvent) CreationDate() time.Time {
	return e.creationDate
}

// EventNewEntity is an event that represents the creation of a new entity.
// We know the id of the entity we want to create
type EventNewEntity struct {
	// baseEvent to implement the Event interface
	baseEvent
	// entityId is the unique identifier of the entity (immutable)
	entityId string
}

// EntityId returns the immutable id of the entity
func (e EventNewEntity) EntityId() string {
	return e.entityId
}

// BuildEventNewEntity creates a new event that represents the creation of a new entity.
// Parameter entityId is the unique identifier of the entity (immutable)
func BuildEventNewEntity(entityId string) EventNewEntity {
	return EventNewEntity{
		baseEvent: baseEvent{
			eventId:      commons.NewId(),
			creationDate: time.Now(),
		},
		entityId: entityId,
	}
}

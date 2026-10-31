package entities

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
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
	// entityId to get the id of the entity to change
	entityId string
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

// EntityId returns the id of the entity that event applies to.
// The event id is for the event, the entity id is for the entity to change.
func (e baseEvent) EntityId() string {
	return e.entityId
}

// EventChangeActivity changes the activity of an entity.
// For a new entity, it creates the entity with that activity
type EventChangeActivity struct {
	// baseEvent to share main elements
	baseEvent
	// activity is the new activity of the entity
	activity periods.Period
}

// Activity returns the new activity of the entity.
// Once event is processed, the entity will have this activity
func (e EventChangeActivity) Activity() periods.Period {
	return e.activity
}

// BuildEventChangeActivity builds a new event to change the activity of an entity.
func BuildEventChangeActivity(entityId string, newActivity periods.Period) EventChangeActivity {
	return EventChangeActivity{
		baseEvent: baseEvent{
			eventId:      commons.NewId(),
			entityId:     entityId,
			creationDate: time.Now(),
		},
		activity: newActivity,
	}
}

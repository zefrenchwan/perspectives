package entities

import (
	"strings"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Event is the mechanism to change some content.
// Events are immutable.
type Event interface {
	// Identifiable for an event to have an id
	commons.Identifiable
	// Hashable for an event to have a hash string (immutable)
	commons.Hashable
	// RecordDate is the creation date of the event
	RecordDate() time.Time
	// isEvent to seal the interface
	isEvent()
}

// genericEvent is the base implementation of an event.
// Use it for any implementation of an event, to not manage basic code about core fields.
type genericEvent struct {
	// id of the event
	id string
	// hashString is the hash of the event.
	// DO NOT include it in the hash string computation
	hashString string
	// recordDate is the creation date of the event
	recordDate time.Time
}

// isEvent to seal the interface : only this package can implement it
func (e genericEvent) isEvent() {}

// Id is the id of the event
func (e genericEvent) Id() string {
	return e.id
}

// ToHashString is the hash string of the event.
// Because events are immutable, the hash string is computed from the id and the record date
func (e genericEvent) ToHashString() string {
	return e.hashString
}

// RecordDate is the creation date of the event
func (e genericEvent) RecordDate() time.Time {
	return e.recordDate
}

// LocalEvent defines events that apply to ONE entity only
type LocalEvent interface {
	// Event is of course parent interface : a local event is an event
	Event
	// EntityId is the id of the entity concerned by the event
	EntityId() string
}

// genericLocalEvent is the generic implementation of LocalEvent
type genericLocalEvent struct {
	// genericEvent to automatically reuse the event part
	genericEvent
	// entityId is the id of the entity concerned by the event
	entityId string
}

// EntityId is the id of the entity concerned by the event
func (e genericLocalEvent) EntityId() string {
	return e.entityId
}

// EntityCreationEvent defines an event that creates an entity
type EntityCreationEvent struct {
	// genericLocalEvent to automatically reuse the local event part
	genericLocalEvent
	// initialState is the initial state of the entity to create
	initialState State
	// initialPeriod is the initial activity of the entity to create
	initialPeriod periods.Period
}

// InitialState is the state of the new entity to create
func (e EntityCreationEvent) InitialState() State {
	return e.initialState
}

// InitialPeriod is the activity of the new entity to create
func (e EntityCreationEvent) InitialPeriod() periods.Period {
	return e.initialPeriod
}

// CreateEntityEvent creates an event that creates an entity.
func CreateEntityEvent(
	eventId, // id of the event
	entityId string, // id of the entity to create
	recordDate time.Time, // record date of the event
	initialLifetime periods.Period, // initial activity of the entity to create
	initialState State, // initial state of the entity to create
) EntityCreationEvent {
	base := genericEvent{
		id:         eventId,
		recordDate: recordDate,
	}

	localEvent := genericLocalEvent{
		genericEvent: base,
		entityId:     entityId,
	}

	result := EntityCreationEvent{
		genericLocalEvent: localEvent,
		initialState:      initialState,
		initialPeriod:     initialLifetime,
	}

	var contentWriter strings.Builder
	contentWriter.WriteString("create entity event. id = \n")
	contentWriter.WriteString(eventId)
	contentWriter.WriteString(" record date = ")
	contentWriter.WriteString(recordDate.UTC().Format(time.RFC3339))
	contentWriter.WriteString("\n ENTITY INFORMATION \nentity id = ")
	contentWriter.WriteString(entityId)
	contentWriter.WriteString(" with initial state ")
	contentWriter.WriteString(initialState.ToHashString())
	contentWriter.WriteString(" and initial period ")
	contentWriter.WriteString(initialLifetime.AsRawString())

	result.hashString = commons.HashString(contentWriter.String())

	return result
}

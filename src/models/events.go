package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
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

// String is the string representation of a genericLocalEvent
func (e genericLocalEvent) String() string {
	var contentWriter strings.Builder
	contentWriter.WriteString("GenericLocalEvent\nid = \n")
	contentWriter.WriteString(strconv.Itoa(len(e.id)) + ":" + e.id)
	contentWriter.WriteString(" record date = ")
	contentWriter.WriteString(e.recordDate.UTC().Format(time.RFC3339))
	contentWriter.WriteString("\n ENTITY INFORMATION \nentity id = ")
	contentWriter.WriteString(strconv.Itoa(len(e.entityId)) + ":" + e.entityId)
	return contentWriter.String()
}

// EntityCreationEvent defines an event that creates an entity.
// It contains TWO information :
// 1) the values to set to start the entity
// 2) the TYPE OF VALUES to set to start the entity (function ? relation ? )
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
	contentWriter.WriteString("EntityCreationEvent\n")
	contentWriter.WriteString(result.genericLocalEvent.String())
	contentWriter.WriteString("\n with initial state ")
	contentWriter.WriteString(initialState.ToHashString())
	contentWriter.WriteString(" and initial period ")
	contentWriter.WriteString(initialLifetime.AsRawString())

	result.hashString = commons.HashString(contentWriter.String())

	return result
}

// ActivityChangeEvent means that the activity of an entity changed to a new value
type ActivityChangeEvent struct {
	// genericLocalEvent to apply to an entity
	genericLocalEvent
	// activity is the new activity of the entity
	activity periods.Period
}

// Activity returns the new activity of the entity
func (e ActivityChangeEvent) Activity() periods.Period {
	return e.activity
}

// UpdateActivityEvent creates an event that updates the activity of an entity
func UpdateActivityEvent(
	id string, // id of the event
	entityId string, // id of the entity to update
	recordDate time.Time, // record date of the event
	activity periods.Period, // new activity to set
) ActivityChangeEvent {
	base := genericEvent{
		id:         id,
		recordDate: recordDate,
	}

	entityEvent := genericLocalEvent{
		genericEvent: base,
		entityId:     entityId,
	}

	result := ActivityChangeEvent{
		genericLocalEvent: entityEvent,
		activity:          activity,
	}

	var baseWriter strings.Builder
	baseWriter.WriteString("ActivityChangeEvent\n")
	baseWriter.WriteString(result.genericLocalEvent.String())
	baseWriter.WriteString(", activity: " + activity.AsRawString())
	result.hashString = commons.HashString(baseWriter.String())

	return result
}

// AttributeChangeEvent creates an event that updates the attribute of an entity
type AttributeChangeEvent struct {
	// genericLocalEvent to apply to an entity
	genericLocalEvent
	// name of the attribute to update
	name string
	// value is the new value to set
	value values.PrimitiveValue
	// period is the validity period of the attribute change
	period periods.Period
}

// Name of the attribute to update
func (e AttributeChangeEvent) Name() string {
	return e.name
}

// Validity is the period to set value for
func (e AttributeChangeEvent) Validity() periods.Period {
	return e.period
}

// Value is the new value to set during that period
func (e AttributeChangeEvent) Value() values.PrimitiveValue {
	return e.value
}

// RoleChangeEvent upserts the value of a role for a given period
type RoleChangeEvent struct {
	// genericLocalEvent for entity event
	genericLocalEvent
	// name of the role to upsert
	name string
	// value of the role to upsert
	value values.ReferenceValue
	// period new validity of the role
	period periods.Period
}

// Name of the role to update
func (e RoleChangeEvent) Name() string {
	return e.name
}

// Validity of the role to upsert
func (e RoleChangeEvent) Validity() periods.Period {
	return e.period
}

// Value of the role to upsert
func (e RoleChangeEvent) Value() values.ReferenceValue {
	return e.value
}

// UpdateAttributeEvent creates an event to upsert an attribute value for a given period
func UpdateAttributeEvent(
	id string, // id of the event
	entityId string, // id of the entity to update
	recordDate time.Time, // record date of the event
	name string, // name of the attribute to upsert
	value values.PrimitiveValue, // value of the attribute to upsert
	validity periods.Period, // period is the new validity of the attribute
) AttributeChangeEvent {

	base := genericEvent{
		id:         id,
		recordDate: recordDate,
	}

	entityEvent := genericLocalEvent{
		genericEvent: base,
		entityId:     entityId,
	}

	result := AttributeChangeEvent{
		genericLocalEvent: entityEvent,
		name:              name,
		value:             value,
		period:            validity,
	}

	var baseWriter strings.Builder
	baseWriter.WriteString(result.genericLocalEvent.String())
	baseWriter.WriteString("\n")
	baseWriter.WriteString(strconv.Itoa(len(name)) + ":" + result.name)
	baseWriter.WriteString("\nvalue => " + value.ToHashString())
	baseWriter.WriteString("\nvalidity => " + validity.AsRawString())

	result.hashString = commons.HashString(baseWriter.String())

	return result
}

// UpdateRoleEvent creates an event to upsert a role value for a given period
func UpdateRoleEvent(
	id string, // id of the event
	entityId string, // id of the entity to update
	recordDate time.Time, // record date of the event
	name string, // name of the role to upsert
	value values.ReferenceValue, // value of the role to upsert
	validity periods.Period, // period is the new validity of the role
) RoleChangeEvent {
	base := genericEvent{
		id:         id,
		recordDate: recordDate,
	}

	entityEvent := genericLocalEvent{
		genericEvent: base,
		entityId:     entityId,
	}

	result := RoleChangeEvent{
		genericLocalEvent: entityEvent,
		name:              name,
		value:             value,
		period:            validity,
	}

	var baseWriter strings.Builder
	baseWriter.WriteString(result.genericLocalEvent.String())
	baseWriter.WriteString("\n")
	baseWriter.WriteString(strconv.Itoa(len(name)) + ":" + result.name)
	baseWriter.WriteString("\nvalue => " + value.ToHashString())
	baseWriter.WriteString("\nvalidity => " + validity.AsRawString())

	result.hashString = commons.HashString(baseWriter.String())

	return result
}

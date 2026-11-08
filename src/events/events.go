package events

import (
	"iter"
	"slices"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Event is the container for changes grouped together
type Event struct {
	// eventId is the unique id of the event
	eventId string
	// eventCreationDate is the creation date of the event
	eventCreationDate time.Time
	// changes is all the changes to apply for that event
	changes []Change
}

// Id is the id of the event, should be unique
func (e Event) Id() string {
	return e.eventId
}

// RecordDate is the date the event was created at
func (e Event) RecordDate() time.Time {
	return e.eventCreationDate
}

// Changes is the list of all the changes to apply
func (e Event) Changes() iter.Seq[Change] {
	return slices.Values(e.changes)
}

// CreateBaseEvent creates a new event that regroups all the changes
func CreateBaseEvent(changes []Change) Event {
	result := Event{eventId: commons.NewId(), eventCreationDate: time.Now()}
	result.changes = append(result.changes, changes...)
	return result
}

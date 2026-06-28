package entities

import (
	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Entity is the unique concept of the system. It is the base of all the other concepts.
// It is ABSOLUTELY MANDATORY to make it as an immutable object.
type Entity interface {
	commons.Identifiable // Identifiable to get the id of the entity.
	periods.TimeBounded  // TimeBounded to define an activity
	Stateful             // Stateful to define state as a map of attributes linked to its time-dependent values
	Linkable             // Linkable to link an entity with others (and then building a reified graph)
}

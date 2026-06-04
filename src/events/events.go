package events

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// Event is the general interface for any event that change elements of the system
type Event interface {
	// RecordDate returns the date when the event was recorded
	RecordDate() time.Time
	// SourceId returns the id of the source that generated the event
	SourceId() string
	// Id returns the unique identifier of the event
	Id() string
}

// ==================================================================
// TRAITS EVENTS : add trait, change trait, remove traits
// ==================================================================

// EventTrait is the general interface for any trait event
type EventTrait interface {
	Event
	Name() string
}

// EventCreateTrait is the general interface for any trait event that creates a new trait.
type EventCreateTrait interface {
	EventTrait
	// Description returns the attributes and related types of the trait
	Description() map[string]string
}

// EventAddTraitAttributes is the general interface for any trait event that adds attributes to a trait.
// Note that if there is already this attribute, it will be overwritten
type EventAddTraitAttributes interface {
	EventTrait
	// Attributes returns the attributes and related types of the trait.
	// For instance, on a trait of type "Person", the attributes could be "name", "age", etc.
	Attributes() map[string]string
}

// EventRemoveTraitAttributes represents an event for removing specific attributes from a trait.
// It extends EventTrait and provides the Attributes method to list the attributes being removed.
type EventRemoveTraitAttributes interface {
	EventTrait
	// Attributes returns the attributes to be removed from the trait.
	Attributes() []string
}

// EventRemoveTrait represents a specialized event indicating the removal of that trait by name.
type EventRemoveTrait interface {
	EventTrait
}

// =================================================================================
// EVENTS : create instance, change instance lifetime, change instance state
// =================================================================================

// EventInstance is the general event definition to act on instances
type EventInstance interface {
	// InstanceId returns the unique identifier of the instance.
	InstanceId() string
}

// EventCreateInstance represents an event for creating a new instance.
type EventCreateInstance interface {
	EventInstance
	// Lifetime returns the period during which the instance is active.
	Lifetime() periods.Period
}

// EventChangeInstanceLifetime represents an event for changes in the lifecycle of an instance.
type EventChangeInstanceLifetime interface {
	EventInstance
	// Validity returns the time at which the instance's lifecycle changes
	Validity() time.Time
	// Active indicates whether the instance should be active or inactive during the specified period.
	Active() bool
}

// EventChangeInstanceState represents an event for changing the state of an instance.
type EventChangeInstanceState interface {
	EventInstance
	// Validity returns the time period during which the instance's state is changed.
	Validity() periods.Period
	// Attribute returns the attribute name that is being changed.
	Attribute() string
	// Value returns the new value for the attribute.
	Value() any
}

package entities

import (
	"github.com/zefrenchwan/perspectives.git/periods"
)

// EntityBuilder is an interface for building entities.
// Entities are immutable entities that represent entities in the system.
// So to act on their state, a builder allows entities creation and entities copies from others for modification.
// Rules are:
// => Name of the attribute must not be blank, and it is case sensitive.
// => Same for roles.
// => Attributes types should be primitive as defined within the project.
// => Periods of attributes are not related to entity's activity.
// Reason is that lifetime change would make information loss.
// => Removing a period larger than base duration deletes the attribute or role (no empty rule).
// It means that we delete anything that would leave a role or attribute empty (no period, no element).
// => To manage chaining, return the same builder.
type EntityBuilder interface {
	// WithActivity changes the current activity of the entity to build
	WithActivity(period periods.Period) EntityBuilder
	// WithAttributeDuring adds an attribute to the entity during a specific period.
	// If attribute already exists, value will be overwritten for given period, rest remaining unchanged.
	WithAttributeDuring(attribute string, period periods.Period, value any) EntityBuilder
	// WithoutAttributeDuring removes period for that attribute.
	WithoutAttributeDuring(attribute string, period periods.Period) EntityBuilder
	// Cut transforms all periods (activity and attributes) as current value intersection with period.
	// If attributes or lifetime periods are not intersecting with given period, they are removed.
	// In particular, entity may become empty.
	Cut(period periods.Period) EntityBuilder
	// WithLinkDuring adds a link from the entity to another entity during a specific period.
	// Parameters are role (name of the link), period (time span of the link), and operand (the entity linked).
	WithLinkDuring(role string, period periods.Period, operand Entity) EntityBuilder
	// WithoutLinkDuring removes the link for that role and that entity during a given period.
	WithoutLinkDuring(role string, period periods.Period, element Entity) EntityBuilder
	// Errors returns all errors encountered during building.
	Errors() error
	// Build returns the built entity and any errors encountered during building.
	Build() (Entity, error)
}

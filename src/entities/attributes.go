package entities

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// AttributeDetails represents the metadata details of the attribute.
// It contains information about the attribute's name, type, validity, and entity activity.
type AttributeDetails struct {
	// AttributeName is the actual name of the attribute
	AttributeName string
	// AttributeType is the actual type of the attribute
	AttributeType string
	// AttributeValidity is the validity period of the attribute
	AttributeValidity periods.Period
	// EntityActivity is the activity period of the entity
	EntityActivity periods.Period
}

// Stateful is an interface that represents an entity with attributes and values.
// Those values are dynamic, meaning they can change over time.
type Stateful interface {
	// Hashable allows to calculate a hash of the stateful entity.
	commons.Hashable

	// SameState compares two states and tests whether they are equal.
	// Remember : SameState on an entity does not mean same id.
	SameState(other Stateful) bool

	// Attributes allows iteration over the attributes of the entity, by name.
	// It does not return the values of the attributes in a slice, to avoid multiple allocations.
	Attributes() iter.Seq[string]
	// Attribute returns the details of the attribute with the given name.
	Attribute(attribute string) (AttributeDetails, bool)

	// Values allows an iteration over the attributes (by name) and values (for that name)
	Values() iter.Seq2[string, DynamicValues]
	// Value returns the values of the attribute with the given name (if it exists).
	Value(attribute string) (DynamicValues, bool)
	// ValuesAt returns the values of the entity at the given moment.
	ValuesAt(moment time.Time) (iter.Seq2[string, any], bool)
	// ValueAt returns the value of the attribute with the given name at the given moment (if it exists).
	ValueAt(attribute string, moment time.Time) (any, bool)
}

package states

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// Attribute defines a name, a type, and time dependent values.BAM
// Formally, it is a period-dependent function of values, by name.
// In terms of implementation, many options appear:
// 1. In memory implementation for POC / testing / R&D
// 2. Storage-based implementation for production and long historical content.
type Attribute interface {
	// Hashable to get the hash of the attribute.
	commons.Hashable
	// Name of the attribute.
	Name() string
	// Values of the attribute : key is period for the value, value as a primitive value.
	Values() iter.Seq2[periods.Period, values.PrimitiveValue]
	// ValueAt returns the value at that time.
	ValueAt(time.Time) values.PrimitiveValue
	// Datatype is the primitive type of the attribute.
	Datatype() string
	// Domain is the period over which the attribute is defined.
	Domain() periods.Period
}

// State is current state of an entity.
// State should be immutable for time travel, especially audit purpose.
// If not, ensure hash keeps being consistent with the state.
type State interface {
	// Identifiable to get the identifier of the state.
	commons.Identifiable
	// TimeBounded because a state is defined over a period.
	periods.TimeBounded
	// Hashable to get the hash of the state.
	commons.Hashable
	// Attributes returns the attributes of the state by name and related value.
	Attributes() iter.Seq2[string, Attribute]
	// Attribute returns the attribute by name, if any
	Attribute(string) (Attribute, bool)
}

// baseImmutableState is the base implementation of State.
// Because attributes are interface based, we just use them in a map with string as key.
type baseImmutableState struct {
	// id is the stable id of the map
	id string
	// attributes is the map of attributes by name
	attributes map[string]Attribute
	// activity is the period over which the state is defined
	activity periods.Period
	// hashString is the hash of the state, calculated once, because it is immutable
	hashString string
}

// Id returns the identifier of the state.
func (b *baseImmutableState) Id() string {
	return b.id
}

// Activity returns the period over which the state is defined.
func (b *baseImmutableState) Activity() periods.Period {
	return b.activity
}

// ToHashString returns the hash of the state.
func (b *baseImmutableState) ToHashString() string {
	return b.hashString
}

// Attributes returns the attributes of the state by name and related value.
func (b *baseImmutableState) Attributes() iter.Seq2[string, Attribute] {
	return func(yield func(string, Attribute) bool) {
		for name, attr := range b.attributes {
			if !yield(name, attr) {
				return
			}
		}
	}
}

// Attribute returns the attribute by name, if any
func (b *baseImmutableState) Attribute(name string) (Attribute, bool) {
	attr, ok := b.attributes[name]
	return attr, ok
}

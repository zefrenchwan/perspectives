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
type Attribute interface {
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
type State interface {
	// Identifiable to get the identifier of the state.
	commons.Identifiable
	// TimeBounded because a state is defined over a period.
	periods.TimeBounded
	// Attributes returns the attributes of the state by name and related value.
	Attributes() iter.Seq2[string, Attribute]
	// Attribute returns the attribute by name, if any
	Attribute(string) (Attribute, bool)
}

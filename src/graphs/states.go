package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// State is the immutable description of an entity at a given time.
type State interface {
	// Identifiable to define a unique identifier for the state.
	commons.Identifiable
	// Hashable to get the hash of the state.
	// States are immutables, so hash sums up the current state.
	commons.Hashable
	// TimeBounded to define a time period during which the entity exists.
	// It may vary, because, for instance, X is alive so far, until death (and then end of period).
	periods.TimeBounded
	// Attributes describe the state of an element.
	// Keys are names, and values are attributes (basically a map[period]primitives)
	Attributes() iter.Seq2[string, values.ImmutableValuesMapping[values.PrimitiveValue]]
	// Roles describe the relationships between elements.
	// Keys are names, and values are roles (basically a map[period]references)
	Roles() iter.Seq2[string, values.ImmutableValuesMapping[values.ReferenceValue]]
}

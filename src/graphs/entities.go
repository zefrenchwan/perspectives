package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Entity is an immutable graph element.
type Entity interface {
	// Identifiable provides a unique identifier for the entity.
	commons.Identifiable
	// TimeBounded to define a time period during which the entity exists.
	periods.TimeBounded
	// Hashable provides a hash value for the entity, assumed to be immutable.
	commons.Hashable

	// Attributes describe the state of an element.
	// Keys are names, and values are attributes (basically a map[period]primitives)
	Attributes() iter.Seq2[string, Attribute]
	// Roles describe the relationships between elements.
	// Keys are names, and values are roles (basically a map[period]references)
	Roles() iter.Seq2[string, Role]
}

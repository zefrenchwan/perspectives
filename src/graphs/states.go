package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/values"
)

type State interface {
	commons.Identifiable
	commons.Hashable

	// Attributes describe the state of an element.
	// Keys are names, and values are attributes (basically a map[period]primitives)
	Attributes() iter.Seq2[string, values.ImmutableValuesMapping[values.PrimitiveValue]]
	// Roles describe the relationships between elements.
	// Keys are names, and values are roles (basically a map[period]references)
	Roles() iter.Seq2[string, values.ImmutableValuesMapping[values.ReferenceValue]]
}

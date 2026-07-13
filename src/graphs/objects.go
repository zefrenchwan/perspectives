package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

type GraphObject interface {
	commons.Identifiable
	periods.TimeBounded
	commons.Hashable

	Attributes() iter.Seq2[string, Attribute]
	Roles() iter.Seq2[string, Role]
}

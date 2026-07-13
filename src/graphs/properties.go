package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

type Property[V values.Value] interface {
	commons.Hashable
	Name() string
	Values() iter.Seq2[periods.Period, V]
}

type Attribute Property[values.PrimitiveValue]

type Role Property[values.ReferenceValue]

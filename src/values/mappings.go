package values

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func EnsureValuesMappingInvariant[V Value](rawMapping periods.DynamicMapping[V]) bool {
	var expectedType string
	for _, value := range rawMapping.Range() {
		if expectedType == "" {
			expectedType = value.Datatype()
		} else if value.Datatype() != expectedType {
			return false
		}
	}

	if expectedType == "" {
		// no element, so ok so far
		return true
	}

	// ALl elements have the same type and it is what we expect ? OK
	return expectedType == rawMapping.DataType()
}

type ImmutableValuesMapping[V Value] interface {
	commons.Hashable
	IsEmpty() bool
	Range() iter.Seq2[periods.Period, V]
	ValuesType() string
}

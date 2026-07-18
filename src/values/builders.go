package values

import "github.com/zefrenchwan/perspectives.git/periods"

type ReferenceMappingBuilder interface {
	Add(reference string, period periods.Period) error
	Remove(periods.Period)
	Build() (ImmutableValuesMapping[ReferenceValue], error)
}

type PrimitiveMappingBuilder interface {
	ValuesType() string
	Add(value any, period periods.Period) error
	Remove(periods.Period)
	Build() (ImmutableValuesMapping[PrimitiveValue], error)
}

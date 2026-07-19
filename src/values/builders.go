package values

import "github.com/zefrenchwan/perspectives.git/periods"

// ReferenceMappingBuilder is a toolbox to build a mapping of periods to reference values.
type ReferenceMappingBuilder interface {
	// Add adds a reference value to the mapping for the given period.
	// It may raise an error, such as when the reference is empty
	Add(reference string, period periods.Period) error
	// Remove removes a reference value from the mapping for the given period.
	Remove(periods.Period)
	// Build builds a mapping of periods to reference values.
	// It may raise an error, if values are inconsistent.
	Build() (ImmutableValuesMapping[ReferenceValue], error)
}

// PrimitiveMappingBuilder is a toolbox to build a mapping of periods to primitive values.
// It decorates a mapping of periods to primitive values to allow access to core values (not primitive decorator).
type PrimitiveMappingBuilder interface {
	// ValuesType returns the type of values that this builder can build.
	ValuesType() string
	// Add adds raw values to the mapping for the given period.
	Add(value any, period periods.Period) error
	// Remove removes a raw value from the mapping for the given period.
	Remove(periods.Period)
	// Build builds a mapping of periods to primitive values.
	// It may raise an error, if values types are inconsistent.
	Build() (ImmutableValuesMapping[PrimitiveValue], error)
}

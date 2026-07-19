package values

import (
	"errors"
	"fmt"
	"iter"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// ReferenceMappingBuilder is a toolbox to build a mapping of periods to reference values.
type ReferenceMappingBuilder interface {
	// Add adds a reference value to the mapping for the given period.
	// It may raise an error, such as when the reference is empty
	Add(reference string, period periods.Period) error
	// Load appends an existing mapping of references
	Load(other ImmutableValuesMapping[ReferenceValue]) error
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
	// Load appends an existing mapping of values.
	// It may raise an error, if other values types are inconsistent.
	Load(other ImmutableValuesMapping[PrimitiveValue]) error
	// Remove removes a raw value from the mapping for the given period.
	Remove(periods.Period)
	// Build builds a mapping of periods to primitive values.
	// It may raise an error, if values types are inconsistent.
	Build() (ImmutableValuesMapping[PrimitiveValue], error)
}

// genericValuesMapping is a generic implementation of ImmutableValuesMapping.
// It just wraps a DynamicMapping and provides a hash function.
type genericValuesMapping[V Value] struct {
	// wrapped is the actual mapping
	wrapped periods.DynamicMapping[V]
	// hash, calculated once due to immutability
	hash string
}

// IsEmpty returns true if the mapping is empty and false otherwise.
func (m *genericValuesMapping[V]) IsEmpty() bool {
	return m.wrapped == nil || m.wrapped.IsEmpty()
}

// Range returns an iterator over periods and values in the mapping.
func (m *genericValuesMapping[V]) Range() iter.Seq2[periods.Period, V] {
	if m.wrapped == nil {
		return func(yield func(periods.Period, V) bool) {}
	}
	return m.wrapped.Range()
}

// ValuesType is the decorated type of values.
func (m *genericValuesMapping[V]) ValuesType() string {
	return m.wrapped.DataType()
}

// ToHashString returns the hash string of the mapping.
// Due to immutability, the hash string is computed only once and stored in the mapping.
func (m *genericValuesMapping[V]) ToHashString() string {
	return m.hash
}

// mappingDecoratorReference decorates a reference mapping
type mappingDecoratorReference struct {
	// decorated is the original mapping that manages the reference values.
	decorated periods.DynamicMapping[ReferenceValue]
}

// NewReferenceMappingBuilder decorates an original mapping for reference values.
// It may be originally empty.
func NewReferenceMappingBuilder(originalMapping periods.DynamicMapping[ReferenceValue]) ReferenceMappingBuilder {
	if originalMapping == nil {
		return &mappingDecoratorReference{
			decorated: nil,
		}
	}

	mappingCopy := periods.DynamicMappingCopy(originalMapping)
	return &mappingDecoratorReference{
		decorated: mappingCopy,
	}
}

// ValuesType returns the reference constant (because we deal with references)
func (r *mappingDecoratorReference) ValuesType() string {
	return r.decorated.DataType()
}

// Add a reference to the builder for a given period
func (r *mappingDecoratorReference) Add(reference string, period periods.Period) error {
	if reference == "" {
		return errors.New("reference cannot be empty")
	}

	matchedValue := NewReference(reference)
	r.decorated.Add(matchedValue, period)
	return nil
}

// Remove clears all the reference values from the mapping during the given period
func (r *mappingDecoratorReference) Remove(period periods.Period) {
	r.decorated.Remove(period)
}

// Load appends an existing mapping of references
func (r *mappingDecoratorReference) Load(other ImmutableValuesMapping[ReferenceValue]) error {
	for period, refValue := range other.Range() {
		r.decorated.Add(refValue, period)
	}
	return nil
}

// Build returns the values mapping as an immutable content, or an error if any
func (r *mappingDecoratorReference) Build() (ImmutableValuesMapping[ReferenceValue], error) {
	if r.decorated == nil {
		return nil, errors.New("cannot build a mapping : invalid source values")
	}

	hashValue := periods.HashDynamicMapping(r.decorated)

	// Build the immutable mapping
	return &genericValuesMapping[ReferenceValue]{
		wrapped: r.decorated,
		hash:    hashValue,
	}, nil
}

// mappingDecoratorPrimitive decorates a dynamic mapping of periods to primitive values.
// The reason is that the dynamic function may be implemeted in memory or database or...
// So we just decorate the original mapping to ensure invariants
type mappingDecoratorPrimitive struct {
	// decorated is the original mapping that we decorate.
	decorated periods.DynamicMapping[PrimitiveValue]
}

// NewPrimitiveMappingBuilder decorates an original mapping for primitive values.
// It may be originally empty.
func NewPrimitiveMappingBuilder(originalMapping periods.DynamicMapping[PrimitiveValue]) PrimitiveMappingBuilder {
	if originalMapping == nil {
		return &mappingDecoratorPrimitive{
			decorated: nil,
		}
	}

	mappingCopy := periods.DynamicMappingCopy(originalMapping)
	return &mappingDecoratorPrimitive{
		decorated: mappingCopy,
	}
}

// ValuesType returns the type of values that this builder can build.
// For instance, int, string, etc.
func (p *mappingDecoratorPrimitive) ValuesType() string {
	return p.decorated.DataType()
}

// Add adds a value to the mapping for a given period.
// Note that the value must be a primitive value directly.
func (p *mappingDecoratorPrimitive) Add(value any, period periods.Period) error {
	// Get the primitive value from the given value, if possible
	matchedValue, err := BuildPrimitiveValue(value)
	if err != nil {
		return err
	} else if period.IsEmpty() {
		return nil
	}

	expectedType := p.decorated.DataType()
	realType := matchedValue.Datatype()
	if realType != expectedType {
		return fmt.Errorf("cannot add a value of type %s to a mapping of type %s", matchedValue.Datatype(), p.decorated.DataType())
	}

	p.decorated.Add(matchedValue, period)
	return nil
}

// Load loads the values from another mapping into this one.
// It may raise an error, if other values types are inconsistent.
func (p *mappingDecoratorPrimitive) Load(other ImmutableValuesMapping[PrimitiveValue]) error {
	if other == nil {
		return nil
	}

	if other.ValuesType() != p.decorated.DataType() {
		return fmt.Errorf("cannot load a mapping of type %s into a mapping of type %s", other.ValuesType(), p.decorated.DataType())
	}

	for period, value := range other.Range() {
		p.decorated.Add(value, period)
	}

	// NO NEED TO REDO A GLOBAL CHECK : invariant is ok by design
	//if !EnsureValuesMappingInvariant(p.decorated) {
	//	return errors.New("Invariant mapping break)")
	//}

	return nil
}

// Remove clears any value from the mapping for a given period.
func (p *mappingDecoratorPrimitive) Remove(period periods.Period) {
	if p.decorated != nil {
		p.decorated.Remove(period)
	}
}

// Build returns a read-only mapping of primitive values.
// It may raise an error if the mapping is invalid (errors are cumulative).
func (p *mappingDecoratorPrimitive) Build() (ImmutableValuesMapping[PrimitiveValue], error) {
	if p.decorated == nil {
		return nil, errors.New("cannot build a mapping : invalid source values")
	}

	// Check that the mapping is globally valid
	if !EnsureValuesMappingInvariant(p.decorated) {
		return nil, errors.New("l'invariant du mapping n'est pas respecté (types mixtes ou invalides détectés)")
	}

	// Once values are OK, calculate the hash for the mapping
	hashValue := periods.HashDynamicMapping(p.decorated)

	// Build the immutable mapping
	return &genericValuesMapping[PrimitiveValue]{
		wrapped: p.decorated,
		hash:    hashValue,
	}, nil
}

package values

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// EnsureValuesMappingInvariant goes through a mapping and ensures that all values have the same type.
// It returns true if all values have the same expected type and false otherwise.
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

// ImmutableValuesMapping is an immutable mapping of periods to values (reference or primitive values).
// It is used to represent a mapping of periods to values that cannot be modified after creation.
type ImmutableValuesMapping[V Value] interface {
	// Hashable is an interface that provides a hash function for the mapping.
	// Because mapping is immutable, it will not change after creation.
	commons.Hashable
	// IsEmpty returns true if the mapping is empty and false otherwise.
	IsEmpty() bool
	// Range returns an iterator over periods and values in the mapping.
	Range() iter.Seq2[periods.Period, V]
	// ValuesType returns the type of values in the mapping.
	// For instance, on a primitive mapping, it will return "int", "string", etc.
	ValuesType() string
}

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

////////////////////////////////////
// MAPPINGS LOCAL IMPLEMENTATIONS //
////////////////////////////////////

// localNode represents a placeholder for a period and value pair.
// It will be used to store as a node in a localMapping.
type localNode[V Value] struct {
	// duration represents the period for which the value is valid.
	duration periods.Period
	// value represents the value associated with the period.
	value V
}

// String will be used to generate a hash string representation of the localNode.
// So, it is injective and should be kept as is.
// No need to add type, thought, because it is coming from the local mapping.
func (n localNode[V]) String() string {
	return "localNode : duration = " + n.duration.AsRawString() + " value = " + n.value.Serialize()
}

// localMapping represents a mapping of periods to values.
// It is immutable and used to store an in-memory implementation of an immutable values mapping.
type localMapping[V Value] struct {
	// dataType represents the type of values stored in the mapping.
	dataType string
	// nodes represents the actual nodes in the mapping.
	// Remember that we cannot use a map, so we see it as a list of nodes.
	nodes []localNode[V]
	// hashString as the CONSTANT hash string representation of the local mapping.
	// It works because it is based on the sorted string representation of the nodes
	// and nodes are immutable.
	hashString string
}

// localMappingHash calculates the hash string representation of a local mapping.
// Mapping is constant and based on the sorted string representation of the nodes.
func localMappingHash[V Value](l *localMapping[V]) string {
	if len(l.nodes) == 0 {
		return commons.HashString("localMapping : empty for type " + l.dataType)
	}

	size := len(l.nodes)
	sortedContent := make([]string, len(l.nodes))
	for index, value := range l.nodes {
		sortedContent[index] = value.String()
	}

	slices.Sort(sortedContent)
	stringValues := "size = " + strconv.Itoa(size) + " values = " + strings.Join(sortedContent, "|")

	return commons.HashString("localMapping : type = " + l.dataType + " content = " + stringValues)
}

// ToHashString returns the hash string representation of the local mapping (to avoid recalculation).
func (l *localMapping[V]) ToHashString() string {
	return l.hashString
}

// IsEmpty returns true if the local mapping contains no value (empty values are not stored).
func (l *localMapping[V]) IsEmpty() bool {
	return len(l.nodes) == 0
}

// ValuesType returns the type of values stored in the local mapping.
func (l *localMapping[V]) ValuesType() string {
	return l.dataType
}

// Range returns an iterator over the local mapping.
func (l *localMapping[V]) Range() iter.Seq2[periods.Period, V] {
	return func(yield func(periods.Period, V) bool) {
		for _, value := range l.nodes {
			if !yield(value.duration, value.value) {
				return
			}
		}
	}
}

// newLocalMapping is the factory function to build a new immutable local mapping for each kind of values.
// Note that values matching empty periods are not stored in the local mapping.
// It is KEY TO REMEMBER that the local mapping performs NO CALCULATION ON PERIONDS.
// VALUES ARE STORED AS IS (except empty periods).
// So if you want to build a function, you need to perform the necessary calculations before, on values.
func newLocalMapping[V Value, P comparable](
	dataType string, // dataType of value (string, int, reference, etc)
	values map[P]periods.Period, // values are the raw values to map to related V instances
	mapper func(P) V, // mapper is the function to map raw values to V instances
) ImmutableValuesMapping[V] {
	result := new(localMapping[V])
	result.dataType = dataType
	result.nodes = make([]localNode[V], 0)

	for rawContent, matchingPeriod := range values {
		if !matchingPeriod.IsEmpty() {
			mappedValue := mapper(rawContent)
			result.nodes = append(result.nodes, localNode[V]{duration: matchingPeriod, value: mappedValue})
		}
	}

	// hash may now be calculated
	result.hashString = localMappingHash(result)

	return result
}

// NewStringLocalMapping builds a new immutable local mapping for string values linked to periods.
// Note that values matching empty periods are not stored in the local mapping.
func NewStringLocalMapping(values map[string]periods.Period) ImmutableValuesMapping[PrimitiveValue] {
	return newLocalMapping[PrimitiveValue, string](PRIMITIVE_TYPE_STRING, values, func(value string) PrimitiveValue {
		return NewString(value)
	})
}

// NewReferenceLocalMapping builds a new immutable local mapping for references linked to periods
// Note that values matching empty periods are not stored in the local mapping.
func NewReferenceLocalMapping(values map[string]periods.Period) ImmutableValuesMapping[ReferenceValue] {
	return newLocalMapping[ReferenceValue, string](REFERENCE_TYPE, values, func(value string) ReferenceValue {
		return NewReference(value)
	})
}

//////////////////////////////////////////////////////////////////////
// BUILDERS LOCAL IMPLEMENTATIONS TO BUILD IMMUTABLE LOCAL MAPPINGS //
//////////////////////////////////////////////////////////////////////

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

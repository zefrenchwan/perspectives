package values

import (
	"fmt"
	"iter"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// valuesMapping decorated a mapper of V (for instance, functions on primitive values).
// V is the type of the values mapped by the mapping (primitive or reference).
// M is the type of the mapping (for instance, a function or a relation).
type valuesMapping[V Value, M periods.DynamicMapping[V]] struct {
	// mapper is the actual mapping (as a function or relation from time to V).
	mapper M
}

// Datatype returns the data type of the underlying content of V.
// For instance, primitive values for V of type int => result is "int".
func (vm *valuesMapping[V, M]) Datatype() string {
	return vm.mapper.DataType()
}

// Domain returns the union of periods with values set
func (vm *valuesMapping[V, M]) Domain() periods.Period {
	return vm.mapper.Domain()
}

// Equals tests equality based on mapping.
func (vm *valuesMapping[V, M]) Equals(other valuesMapping[V, M]) bool {
	return vm.mapper.Equals(other.mapper)
}

// IsEmpty tests whether the mapping is empty (no element)
func (vm *valuesMapping[V, M]) IsEmpty() bool {
	return vm.mapper.IsEmpty()
}

// Range returns the values per period.
// For relations, periods may not be disjoint.
// For functions, they are disjoint for sure.
func (vm *valuesMapping[V, M]) Range() iter.Seq2[periods.Period, V] {
	return vm.mapper.Range()
}

// Values returns the real content of the mapping.
// No Value is returned, only the period and the actual content.
func (vm *valuesMapping[V, M]) Values() iter.Seq2[periods.Period, any] {
	// work because even references have a conten
	return func(yield func(periods.Period, any) bool) {
		for period, primitiveValue := range vm.mapper.Range() {
			if !yield(period, primitiveValue.Content()) {
				break
			}
		}
	}
}

// Remove removes all the values for the given period.
func (vm *valuesMapping[V, M]) Remove(period periods.Period) {
	vm.mapper.Remove(period)
}

// primitiveMapping is a mapping that only contains primitive values.
type primitiveMapping[M periods.DynamicMapping[PrimitiveValue]] struct {
	valuesMapping[PrimitiveValue, M]
}

// Equals returns true if the two mappings are equal.
func (pm *primitiveMapping[M]) Equals(other primitiveMapping[M]) bool {
	return pm.valuesMapping.Equals(other.valuesMapping)
}

// Add adds a value for that given period: it creates the primitive value and adds it to the mapping.
func (pm *primitiveMapping[M]) Add(value any, period periods.Period) error {
	expectedType := pm.mapper.DataType()
	newValue, errBuild := BuildPrimitiveValue(value)
	if errBuild != nil {
		return errBuild
	} else if realType := newValue.Datatype(); realType != expectedType {
		return fmt.Errorf("value type %s does not match expected type %s", realType, expectedType)
	}

	pm.valuesMapping.mapper.Add(newValue, period)
	return nil
}

// referenceMapping is a mapping that only contains reference values.
type referenceMapping[M periods.DynamicMapping[ReferenceValue]] struct {
	valuesMapping[ReferenceValue, M]
}

// Add adds a reference (as a string) for that given period: it creates the reference value and adds it to the mapping.
func (rm *referenceMapping[M]) Add(reference string, period periods.Period) {
	referenceValue := NewReference(reference)
	rm.valuesMapping.mapper.Add(referenceValue, period)
}

// Equals returns true if the two mappings are equal.
func (rm *referenceMapping[M]) Equals(other referenceMapping[M]) bool {
	return rm.valuesMapping.Equals(other.valuesMapping)
}

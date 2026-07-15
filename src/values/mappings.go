package values

import (
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

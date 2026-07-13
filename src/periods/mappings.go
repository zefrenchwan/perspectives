package periods

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// DynamicMapping is a mapping of values that can change over time.
// For instance, given a company, CEO is a role that may be X during a given period, then Y, then...
// The CEO role would be a DYNAMIC mapping over time.
type DynamicMapping[T any] interface {
	// Domain is the union of periods linked to at least one value.
	Domain() Period
	// Equals returns true if the two mappings are equal :
	// same type, same elements, same underlying type.
	Equals(other DynamicMapping[T]) bool
	// IsEmpty returns true if the mapping is empty.
	IsEmpty() bool
	// Range iterates over the mapping as a couple of matching periods and related value.
	// Note that, for relations, there are multiple values per period.
	// For functions, there is only one value per period.
	Range() iter.Seq2[Period, T]
	// Add adds a value to the mapping.
	// General contract is defined here, but implementation is left to the concrete type.
	// For relations: just add the value to the mapping.
	// For functions: add the value to the mapping to maintain the partition invariant.
	Add(value T, period Period)
	// Remove removes the given period from the mapping and all related values.
	Remove(period Period)
	// DataType returns the type of the values stored in the mapping.
	DataType() string
	// isFunctionalMapping returns true if the mapping is a function.
	// It is true if, given a moment, there is only one value.
	// It means a sealed interface.
	isFunctionalMapping() bool
}

// DynamicRelation is a dynamic mapping with possibly multiple values per period.
// At returns an iterator over the values, there may be many values per period.
type DynamicRelation[T any] interface {
	// DynamicMapping to regroup the common methods.
	DynamicMapping[T]
	// At returns the elements at a given moment (if any).
	// If none matches, the iterator is empty and the second result is then false
	At(moment time.Time) (iter.Seq[T], bool)
	// Copy returns a copy of the dynamic relation.
	Copy() DynamicRelation[T]
}

// DynamicFunction defines a dynamic mapping with ONE value maximum at a given period.
// It is a mapping of values, each value being valid during a specific period.
type DynamicFunction[T any] interface {
	// DynamicMapping to regroup the common methods.
	DynamicMapping[T]
	// At returns the unique element (if any) matching the given moment.
	At(moment time.Time) (T, bool)
	// Copy returns a copy of the dynamic function.
	Copy() DynamicFunction[T]
}

// ===================================================
// HASHING FUNCTION TO CALCULATE EQUALS AND CHANGES ==
// ===================================================

// General assumptions apply :
// That hash system should be injective almost every time.
// Idea is to speed up link equality calculation and avoid a full walkthrough.

// HashDynamicMapping calculates a hash for a dynamic mapping.
// Parameter isFunction indicates whether the mapping is a function (to distinguish the same content, different type)
func HashDynamicMapping[T any](dv DynamicMapping[T], isFunction bool) string {
	if dv == nil || dv.IsEmpty() {
		value := fmt.Sprintf("Dynamic mapping of %s with functional %t", dv.DataType(), isFunction)
		return commons.HashString(value)
	}

	valueType := dv.DataType()

	// We don't know the exact number of periods in advance when using the range iterator,
	// so we start with an empty slice.
	elements := make([]string, 0)

	for period, value := range dv.Range() {
		valueString := fmt.Sprintf("%v", value)
		sizeString := strconv.Itoa(len(valueString))

		// Use strict formatting with length prefixing to prevent delimiter injection.
		// Format: [Period]->Type(Length):Value
		mappedString := fmt.Sprintf("[%s]->%s(%s):%s", period.AsRawString(), valueType, sizeString, valueString)
		elements = append(elements, mappedString)
	}

	// Sort ONLY the dynamic elements to ensure a deterministic hash regardless of iteration order.
	slices.Sort(elements)

	var builder strings.Builder
	builder.WriteString("Dynamic mapping of ")
	builder.WriteString(valueType)
	builder.WriteString(" with functional ")
	builder.WriteString(fmt.Sprintf("%t", isFunction))
	builder.WriteString("\n\n")
	builder.WriteString(strings.Join(elements, "|"))

	return commons.HashString(builder.String())
}

// HashDynamicFunction returns a hash of the given dynamic function.
func HashDynamicFunction[T any](f DynamicFunction[T]) string {
	return HashDynamicMapping(f, true)
}

// HashDynamicRelation returns a hash of the given dynamic relation.
func HashDynamicRelation[T any](r DynamicRelation[T]) string {
	return HashDynamicMapping(r, false)
}

// =========================================================================
// DYNAMIC COLLECTIONS IMPLEMENTATION : in memory, no storage
// =========================================================================

// valueNode stores a value set during a specific matchingPeriod
// value is the actual value (of type T) stored in the node.
type valueNode[T any] struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod Period
	// value is the actual value stored in the node
	value T
}

// valuesHandler stores a set of values, each value being valid during a specific period.
// For relations, values are not disjoint and can overlap.
// For functions, values are disjoint and cannot overlap.
type valuesHandler[T any] struct {
	// values have one value per matching period
	values []valueNode[T]
	// storedType is the actual type name of the content
	storedType string
	// equality function
	equals func(a, b T) bool
	// isFunction is true for partition, false for set
	isFunction bool
}

// Equals returns true if the two dynamic mappings share a same mapping type,
// and have the same values at the same periods, with the exact same type.
func (vh *valuesHandler[T]) Equals(other DynamicMapping[T]) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if other.isFunctionalMapping() != vh.isFunctionalMapping() {
		return false
	} else if other.IsEmpty() && vh.IsEmpty() {
		return true
	}

	// basically, it is checking that two lists are equals with iterators...
	counter := 0 // how many match, to compare to len(vh.values)
	for _, content := range vh.values {
		referencePeriod := content.matchingPeriod
		referenceItem := content.value
		// find it in the other iterator
		found := false
		for otherPeriod, otherValue := range other.Range() {
			if referencePeriod.Equals(otherPeriod) {
				if vh.equals(referenceItem, otherValue) {
					found = true
					counter++
					break
				}
			}
		}

		// No local match => exclude directly
		if !found {
			return false
		}
	}

	return counter == len(vh.values)
}

// IsEmpty checks if the valuesHandler contains any values
func (vh *valuesHandler[T]) IsEmpty() bool {
	return vh == nil || len(vh.values) == 0
}

// Domain returns the union of periods for which values are set
func (vh *valuesHandler[T]) Domain() Period {
	if vh == nil || len(vh.values) == 0 {
		return NewEmptyPeriod()
	}

	validity := NewEmptyPeriod()
	for _, element := range vh.values {
		validity = validity.Union(element.matchingPeriod)
	}

	return validity
}

// first returns the first value at the given moment in time, or nil and false if no value is found.
func (vh *valuesHandler[T]) first(moment time.Time) (T, bool) {
	var empty T
	for _, element := range vh.values {
		if element.matchingPeriod.Contains(moment) {
			return element.value, true
		}
	}
	return empty, false
}

// all returns all the elements matching that moment.
func (vh *valuesHandler[T]) all(moment time.Time) (iter.Seq[T], bool) {
	hasElements := false

	// 1. Check RIGHT NOW whether there are elements
	//Remember that we want to check whether there are matching elements FIRST.
	// Compiler will return directly, we CANNOT do the check within the yield loop.
	for _, value := range vh.values {
		if value.matchingPeriod.Contains(moment) {
			hasElements = true
			break
		}
	}

	// 2. Then, perform the yield loop
	seq := func(yield func(T) bool) {
		for _, value := range vh.values {
			if value.matchingPeriod.Contains(moment) {
				// yield only if the period matches
				if !yield(value.value) {
					return
				}
			}
		}
	}

	return seq, hasElements
}

// Range iterates over all values in the TemporalValues mapping
func (vh *valuesHandler[T]) Range() iter.Seq2[Period, T] {
	return func(yield func(Period, T) bool) {
		for _, element := range vh.values {
			if !yield(element.matchingPeriod, element.value) {
				break
			}
		}
	}
}

// DataType returns the string representation of the type.
// It is UNIQUE by design : this is a key invariant.
func (vh *valuesHandler[T]) DataType() string {
	return vh.storedType
}

// clone returns a copy of the valuesHandler with the same values and type
func (vh *valuesHandler[T]) clone() *valuesHandler[T] {
	result := make([]valueNode[T], len(vh.values))
	copy(result, vh.values)
	return &valuesHandler[T]{
		values:     result,
		storedType: vh.storedType,
		equals:     vh.equals,
		isFunction: vh.isFunction,
	}
}

// Remove removes the given period from the valuesHandler and all related values
func (vh *valuesHandler[T]) Remove(period Period) {
	if len(vh.values) == 0 {
		return
	} else if period.IsEmpty() {
		return
	}

	result := make([]valueNode[T], 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode[T]{matchingPeriod: remaining, value: element.value})
		}
	}

	vh.values = result
}

// cut returns a copy with same values, restricted to a given period
func (vh *valuesHandler[T]) cut(period Period) *valuesHandler[T] {
	remainingValues := make([]valueNode[T], 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode[T]{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler[T]{values: remainingValues, storedType: vh.storedType, equals: vh.equals}
}

// Add appends an element in that handler.
// Depending on the partition flag, the element is added as
// *  another element (set) : just add it
// * another disjoin element (partition) : removes period from existing elements
func (vh *valuesHandler[T]) Add(value T, period Period) {
	if vh == nil || period.IsEmpty() {
		return
	} else if period.IsEmpty() {
		return
	} else if !vh.isFunction {
		vh.values = append(vh.values, valueNode[T]{matchingPeriod: period, value: value})
		return
	}

	// at this point, it is a partition.
	// first, find union of all matching periods (same value)
	commonPeriod := period
	for _, element := range vh.values {
		currentValue := element.value
		currentPeriod := element.matchingPeriod
		if vh.equals(currentValue, value) {
			commonPeriod = commonPeriod.Union(currentPeriod)
		}
	}

	// Then, for each element, remove common period from non-matching values
	var remainingValues []valueNode[T]
	for _, element := range vh.values {
		currentValue := element.value
		currentPeriod := element.matchingPeriod
		if !vh.equals(currentValue, value) {
			remainingPeriod := currentPeriod.Remove(commonPeriod)
			if !remainingPeriod.IsEmpty() {
				remainingValues = append(remainingValues, valueNode[T]{value: currentValue, matchingPeriod: remainingPeriod})
			}
		}
	}

	// We may now add the new value
	vh.values = append(remainingValues, valueNode[T]{matchingPeriod: commonPeriod, value: value})
}

// isFunctionalMapping returns true if the mapping is functional or false for relational
func (vh *valuesHandler[T]) isFunctionalMapping() bool {
	return vh.isFunction
}

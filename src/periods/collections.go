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

// DynamicCollection is a collection of values that can change over time.
// For instance, given a company, CEO is a role that changes over time.
// The CEO role would be a DYNAMIC collection.
type DynamicCollection[T any] interface {
	// Domain is the union of periods linked to at least one value.
	Domain() Period
	// Equals returns true if the two collections are equal :
	// same type, same elements, same underlying type.
	Equals(other DynamicCollection[T]) bool
	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool
	// Range iterates over the collection as a couple of matching periods and related value.
	// Note that, for sets, there are multiple values per period.
	// For partitions, there is only one value per period.
	Range() iter.Seq2[Period, T]
	// Add adds a value to the collection within the collection.
	// General contract is defined here, but implementation is left to the concrete type.
	// For sets : just add the value to the collection.
	// For partitions : add the value to the collection to maintain the partition invariant.
	Add(value T, period Period)
	// Remove removes the given period from the collection and all related values.
	Remove(period Period)
	// DataType returns the type of the values stored in the collection.
	DataType() string
	// isPartition returns true if the collection is a partition.
	// It means a sealed interface.
	isPartition() bool
}

// DynamicSet is a dynamic collection with possibly multiple values per period.
// At returns an iterator over the values, there may be many values per period.
type DynamicSet[T any] interface {
	// DynamicCollection[T] to regroup the common methods.
	DynamicCollection[T]
	// At returns the elements at a given moment (if any).
	// If none matches, the iterator is empty and second result is then false
	At(moment time.Time) (iter.Seq[T], bool)
	// Copy returns a copy of the dynamic set.
	Copy() DynamicSet[T]
}

// DynamicPartition defines a dynamic collection with ONE value maximum at a given period.
// It is a collection of values, each value being valid during a specific period.
type DynamicPartition[T any] interface {
	// DynamicCollection[T] to regroup the common methods.
	DynamicCollection[T]
	// At returns the unique element (if any) matching the given moment.
	At(moment time.Time) (T, bool)
	// Copy returns a copy of the dynamic partition.
	Copy() DynamicPartition[T]
}

// ===================================================
// HASHING FUNCTION TO CALCULATE EQUALS AND CHANGES ==
// ===================================================

// General assumptions apply :
// Hash system should be injective almost every time.
// Idea is to speed up link equality calculation and avoid full walkthrough.

// HashDynamicCollection calculates a hash for a dynamic collection.
// Partition indicates whether the collection is partitioned.
func HashDynamicCollection[T any](dv DynamicCollection[T], partition bool) string {
	if dv == nil || dv.IsEmpty() {
		value := fmt.Sprintf("Dynamic collection of %s with partition %t", dv.DataType(), partition)
		return commons.HashString(value)
	}

	valueType := dv.DataType()

	// We don't know the exact number of periods in advance when using the range iterator,
	// so we start with an empty slice.
	elements := make([]string, 0)

	// Range over the time-dependent values using Go 1.22+ iterator pattern
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
	builder.WriteString("Dynamic collection of ")
	builder.WriteString(valueType)
	builder.WriteString(" with partition ")
	builder.WriteString(fmt.Sprintf("%t", partition))
	builder.WriteString("\n\n")
	builder.WriteString(strings.Join(elements, "|"))

	return commons.HashString(builder.String())
}

// HashDynamicPartition returns a hash of the given dynamic partition.
func HashDynamicPartition[T any](p DynamicPartition[T]) string {
	return HashDynamicCollection(p, true)
}

// HashDynamicSet returns a hash of the given dynamic set.
func HashDynamicSet[T any](p DynamicSet[T]) string {
	return HashDynamicCollection(p, false)
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
// For sets, values are not disjoint, and can overlap.
// For partitions, values are disjoint, and cannot overlap.
type valuesHandler[T any] struct {
	// values have one value per matching period
	values []valueNode[T]
	// storedType is the actual type name of the content
	storedType string
	// equality function
	equals func(a, b T) bool
	// isPartition is true for partition, false for set
	isPartition bool
}

// Equals returns true if the two temporal values are both partitions or sets (same collection type),
// and have the same values at the same periods, with the exact same type.
func (vh *valuesHandler[T]) Equals(other DynamicCollection[T]) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if other.isPartition() != vh.isPartition {
		return false
	} else if other.IsEmpty() && vh.IsEmpty() {
		return true
	}

	counter := 0
	for period, value := range other.Range() {
		counter++
		found := false
		for _, matching := range vh.values {
			if period.Equals(matching.matchingPeriod) {
				found = true
				if !vh.equals(matching.value, value) {
					return false
				}
			}
		}

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

// Range iterates over all values in the TemporalValues collection
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
	return &valuesHandler[T]{values: result, storedType: vh.storedType, equals: vh.equals}
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
	} else if !vh.isPartition {
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

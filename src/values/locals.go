package values

import (
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

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

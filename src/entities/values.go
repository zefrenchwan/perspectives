package entities

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// DynamicValues represents a value that depends on time.
// It is basically equivalent to a map of disjoined time intervals linked to primitive values.
// Implementations have to ensure that value accepts only PrimitiveValue types.
// Implementation (like the full content of entities) should be immutable.
// Among the list of all advantages for this use case, it allows calculating the hash once and reusing it.
type DynamicValues interface {
	commons.Hashable // Hashable because values are immutable so it makes sense to calculate it once.
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Equals returns true if current dynamic value is the same as another DynamicValues.
	// It means : same periods, same values, same underlying type
	Equals(other DynamicValues) bool
	// IsEmpty checks if the current content is empty (no value on a non empty period)
	IsEmpty() bool
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Range iterates over all values in the TimeDependentValues collection
	Range() iter.Seq2[periods.Period, any]
	// DataType returns the type name of the stored values.
	// By design, it should be the same at all times
	DataType() string
}

// hashDynamicValues returns a collision-resistant hash string for the given DynamicValues.
func hashDynamicValues(dv DynamicValues) string {
	if dv == nil || dv.IsEmpty() {
		return commons.HashString("DynamicValues:empty")
	}

	valueType := dv.DataType()

	// We don't know the exact number of periods in advance when using the range iterator,
	// so we start with an empty slice.
	elements := make([]string, 0)

	// Range over the time-dependent values using Go 1.22+ iterator pattern
	for period, value := range dv.Range() {
		valueString := primitiveValueToString(value)
		sizeString := strconv.Itoa(len(valueString))

		// Use strict formatting with length prefixing to prevent delimiter injection.
		// Format: [Period]->Type(Length):Value
		mappedString := fmt.Sprintf("[%s]->%s(%s):%s", period.AsRawString(), valueType, sizeString, valueString)
		elements = append(elements, mappedString)
	}

	// Sort ONLY the dynamic elements to ensure a deterministic hash regardless of iteration order.
	slices.Sort(elements)

	var builder strings.Builder
	builder.WriteString("DynamicValues|")
	builder.WriteString(strings.Join(elements, "|"))

	return commons.HashString(builder.String())
}

// DynamicValuesAdapter is a wrapper around a DynamicPartition that provides a hash string representation.
// Decorating a partition allows to use it as a DynamicValues.
type DynamicValuesAdapter struct {
	// hashString is the hash string representation of the dynamic values.
	// Because the values are immutable, hash should be stable.
	hashString string
	// partition is the decorated dynamic partition that holds the actual values.
	partition periods.DynamicPartition[any]
}

// Equals checks whether two dynamic values have the same content (equals state).
// It uses the hash string of the dynamic values.
func (dv DynamicValuesAdapter) Equals(other DynamicValues) bool {
	return dv.hashString == other.ToHashString()
}

// IsEmpty checks whether the dynamic values are empty.
func (dv DynamicValuesAdapter) IsEmpty() bool {
	return dv.partition.IsEmpty()
}

// ToHashString returns the hash string representation of the dynamic values.
// Because the values are immutable, hash should be stable.
func (dv DynamicValuesAdapter) ToHashString() string {
	return dv.hashString
}

// At returns the value at the given moment.
// It is unique (if any) because the underlying structure is a partition.
func (dv DynamicValuesAdapter) At(moment time.Time) (any, bool) {
	return dv.partition.At(moment)
}

// Range returns the sequence of periods and related values.
func (dv DynamicValuesAdapter) Range() iter.Seq2[periods.Period, any] {
	return dv.partition.Range()
}

// Validity returns the validity of the dynamic values, that is
// the union of all periods for which the value is defined.
func (dv DynamicValuesAdapter) Validity() periods.Period {
	return dv.partition.Domain()
}

// DataType returns the data type of the dynamic values.
func (dv DynamicValuesAdapter) DataType() string {
	return dv.partition.DataType()
}

// NewDynamicValuesFromPartition creates a new dynamic values from the given partition.
// It decorates the original one, but does not modify it.
func NewDynamicValuesFromPartition(partition periods.DynamicPartition[any]) DynamicValues {
	hashString := periods.HashDynamicPartition(partition)
	return DynamicValuesAdapter{hashString: hashString, partition: partition}
}

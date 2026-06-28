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

// AttributeDetails represents the metadata details of the attribute.
// It contains information about the attribute's name, type, validity, and instance activity.
type AttributeDetails struct {
	// AttributeName is the actual name of the attribute
	AttributeName string
	// AttributeType is the actual type of the attribute
	AttributeType string
	// AttributeValidity is the validity period of the attribute
	AttributeValidity periods.Period
	// InstanceActivity is the activity period of the instance
	InstanceActivity periods.Period
}

// DynamicValues represents a value that depends on time.
// It is basically equivalent to a map of disjoined time intervals linked to primitive values.
// Implementations have to ensure that value accepts only PrimitiveValue types.
type DynamicValues interface {
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Equals returns true if instance is the same as another DynamicValues.
	// It means : same periods, same values, same underlying type
	Equals(other DynamicValues) bool
	// IsEmpty checks if the TimeDependentValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Range iterates over all values in the TimeDependentValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, any) bool)
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
	for period, value := range dv.Range {
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

// Stateful is an interface that represents an entity with attributes and values.
// Those values are dynamic, meaning they can change over time.
type Stateful interface {
	// Hashable allows to calculate a hash of the stateful entity.
	commons.Hashable

	// SameState compares two states and test whether they are equal, no matter the id.
	SameState(other Stateful) bool

	// Attributes allows iteration over the attributes of the entity, by name.
	// It does not return the values of the attributes in a slice, to avoid multiple allocations.
	Attributes() iter.Seq[string]
	// Attribute returns the details of the attribute with the given name.
	Attribute(attribute string) (AttributeDetails, bool)

	// Values allows an iteration over the attributes (by name) and values (for that name)
	Values() iter.Seq2[string, DynamicValues]
	// Value returns the values of the attribute with the given name (if it exists).
	Value(attribute string) (DynamicValues, bool)
	// ValueAt returns the value of the attribute with the given name at the given moment (if it exists).
	ValueAt(attribute string, moment time.Time) (any, bool)
	// ValuesAt returns the values of the entity at the given moment.
	ValuesAt(moment time.Time) (iter.Seq2[string, any], bool)
}

package objects

import (
	"math"
	"time"

	"github.com/zefrenchwan/perspectives.git/configuration"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// PrimitiveValue represents a strictly basic data type.
// Custom types (aliases) are explicitly rejected by design to ensure
// seamless serialization and strict Trait matching.
// Except time.Time, which is a special useful case, we want to restrict to basic values.
// No pointer types are allowed, as they are not suitable for serde and distributed systems.
// No structs (except time.Time), as they would allow bad design (use content instead)
type PrimitiveValue interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string |
		bool |
		time.Time
}

// primitiveTypeName returns the string representation of allowed primitive types.
// To ensure that the type is correctly identified and handled, it works with the kind and not the raw name.
func primitiveTypeName(v any) string {
	if v == nil {
		// changing this means changing the behavior of IsPrimitiveValue
		return ""
	}

	// accept time.Time.
	// In general, put in here any additional types that should be considered primitive.
	if _, okTime := v.(time.Time); okTime {
		return "time.Time"
	}

	switch v.(type) {
	case bool:
		return "bool"
	case int:
		return "int"
	case int8:
		return "int8"
	case int16:
		return "int16"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case uint:
		return "uint"
	case uint8:
		return "uint8"
	case uint16:
		return "uint16"
	case uint32:
		return "uint32"
	case uint64:
		return "uint64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case string:
		return "string"
	case time.Time:
		return "time.Time"
	default:
		return ""
	}

}

// equalsTime tests two time.Time values for equality.
func equalsTime(a, b any) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	t1, ok1 := a.(time.Time)
	t2, ok2 := b.(time.Time)

	if !ok1 || !ok2 {
		return false
	}

	return t1.Equal(t2)
}

// defaultEquals tests two values for equality, applying the == operator.
func defaultEquals(a, b any) bool {
	return a == b
}

// equalsFloat tests two floats with an epsilon
func equalsFloat(a, b any) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	switch v1 := a.(type) {
	case float64:
		v2, ok := b.(float64)
		if !ok {
			return false
		}

		return math.Abs(v1-v2) < configuration.LONG_EPSILON

	case float32:
		v2, ok := b.(float32)
		if !ok {
			return false
		}

		diff := v1 - v2
		if diff < 0 {
			diff = -diff
		}
		return diff < configuration.SHORT_EPSILON

	default:
		return false
	}
}

// primitiveTypeEqualsFunc returns a function that tests two values for equality, based on the type name.
// IMPORTANT : it assumes that the values are primitive as defined in PrimitiveValue.
func primitiveTypeEqualsFunc(typeName string) func(any, any) bool {
	switch typeName {
	case "time.Time":
		return equalsTime
	case "float32", "float64":
		return equalsFloat
	default:
		return defaultEquals
	}
}

// IsPrimitiveValue checks if the given value is a PrimitiveValue.
// In contents implementation, it is used to ensure that only primitive values are stored.
func IsPrimitiveValue(v any) bool {
	// note that it depends on the implementation of primitiveTypeName
	return primitiveTypeName(v) != ""
}

// DynamicValues represents a value that depends on time.
// It is basically equivalent to a map of disjoined time intervals linked to primitive values.
// Implementations have to ensure that value accepts only PrimitiveValue types.
type DynamicValues interface {
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Same returns true if content is the same as another TimeDependentValues.
	// It means : same periods, same values, same underlying type
	Same(other DynamicValues) bool
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

// DynamicContent defines attributes and time dependent values over time, with a global activity period.
// For instance, a person has a global life time,
// and that person's name, age, and address can be considered as attributes, while their values can change over time.
type DynamicContent interface {
	// Activity returns the global activity period this content lasts
	Activity() periods.Period
	// Description returns a map of attribute names to their data types.
	// They cannot change over time : it is impossible to change the type of an attribute once it is defined.
	Description() map[string]string
	// Values returns the attributes and their values at a given moment in time.
	// The map keys are attribute names, and the values are the values of those attributes over time
	Values() map[string]DynamicValues
	// Value returns, if any, the values over time for that given attribute
	Value(attribute string) (DynamicValues, bool)
	// At returns, if any, the values of all attributes at a given moment in time.
	// Because it is a snapshot of the content at a specific point in time,
	// result is a map of attribute names to their values at that moment.
	At(moment time.Time) (map[string]any, bool)
	// Matches returns, if any, the period during which this content matches the given trait.
	// For instance, a person may have a student identity, and a student trait may match that identity during a specific period.
	// The returned period indicates the time frame during which the content's attributes and values align with the trait's requirements.
	// If that given trait is incompatible with the content, the result will be empty, false
	Matches(trait Trait) (periods.Period, bool)
	// Same returns whether the content is the same (same values, same attributes) as the other content.
	Same(other DynamicContent) bool
}

// ContentBuilder manages the changes to apply on a given content.
// Typical use is to implement a load from existing content, perform changes and build a new content.
// Conventionally, it returns itself to allow method chaining.
type ContentBuilder interface {
	// WithActivity changes the content's activity to that specific period.
	WithActivity(period periods.Period) ContentBuilder
	// WithAttributeDuring sets the attribute to the given value during the specified period.
	// Types for value are defined in PrimitiveValue.
	// If there is a type change, it should raise an error.
	// For instance, an age that contains 10 and "twenty" should raise an error.
	// Reasons are : storage, type safety, consistency
	WithAttributeDuring(attribute string, period periods.Period, value any) ContentBuilder
	// WithoutAttributeDuring removes the attribute during the specified period.
	// If period covers all the content, the attribute is removed entirely.
	WithoutAttributeDuring(attribute string, period periods.Period) ContentBuilder
	// Cut reduces the content to a given period.
	// Typical use is to restrict attributes values to global content activity.
	Cut(period periods.Period) ContentBuilder
	// Errors returns, if any, current errors so far.
	// Errors are cumulative
	Errors() error
	// Build creates a new content with the applied changes.
	// It returns the new content and an error if any occurred during the build process.
	// It resets the builder to its initial state, ready for new content modifications.
	// But the recommended use would be to create a new content with a new builder.
	Build() (DynamicContent, error)
}

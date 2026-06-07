package objects

import (
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// PrimitiveValue represents a basic data type that can be used as a value in a time-dependent context.
// It does not include pointer types, as they are not suitable for time-dependent values.
// Neither does it include interfaces, as they are not concrete types and cannot be stored directly.
// In particular, it excludes to pass nil as a value, just do not store value instead of nil.
// NOTE : if you add a type in this interface, make sure to review and test in deep the implementations.
// For instance, for the "in memory" implementation, you need to check whether you want to use == or reflect.DeepEqual.
type PrimitiveValue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string | ~bool
}

// primitiveTypeName returns the string representation of allowed primitive types.
// To ensure that the type is correctly identified and handled, it works with the kind and not the raw name.
func primitiveTypeName(v any) string {
	if v == nil {
		// changing this means changing the behavior of IsPrimitiveValue
		return ""
	}

	valueKind := reflect.TypeOf(v).Kind()
	switch valueKind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return valueKind.String()
	default:
		// changing this means changing the behavior of IsPrimitiveValue
		return ""
	}
}

// IsPrimitiveValue checks if the given value is a PrimitiveValue.
// In contents implementation, it is used to ensure that only primitive values are stored.
// Otherwise, it panics if the value is not primitive.
func IsPrimitiveValue(v any) bool {
	// note that it depends on the implementation of primitiveTypeName
	return primitiveTypeName(v) != ""
}

// TimeDependentValue represents a value that depends on time.
// It is basically equivalent to a map of disjoined time intervals linked to primitive values.
// Implementations have to ensure that value accepts only PrimitiveValue types.
type TimeDependentValue interface {
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Same returns true if content is the same as another TimeDependentValues.
	// It means : same periods, same values, same underlying type
	Same(other TimeDependentValue) bool
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

// TimeDependentContent defines attributes and time dependent values over time, with a global activity period.
// For instance, a person has a global life time,
// and that person's name, age, and address can be considered as attributes, while their values can change over time.
type TimeDependentContent interface {
	// Activity returns the global activity period this content lasts
	Activity() periods.Period
	// Description returns a map of attribute names to their data types.
	// They cannot change over time : it is impossible to change the type of an attribute once it is defined.
	Description() map[string]string
	// Values returns the attributes and their values at a given moment in time.
	// The map keys are attribute names, and the values are the values of those attributes over time
	Values() map[string]TimeDependentValue
	// Value returns, if any, the values over time for that given attribute
	Value(attribute string) (TimeDependentValue, bool)
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
	Same(other TimeDependentContent) bool
}

// ContentBuilder manages the changes to apply on a given content.
// Typical use is to implement a load from existing content, perform changes and build a new content.
// Conventionally, it returns itself to allow method chaining.
type ContentBuilder interface {
	// WithActivity changes the content's activity to that specific period.
	WithActivity(period periods.Period) ContentBuilder
	// WithAttributeDuring sets the attribute to the given value during the specified period.
	// Types for value are defined in PrimitiveValue
	WithAttributeDuring(attribute string, period periods.Period, value any) ContentBuilder
	// WithoutAttributeDuring removes the attribute during the specified period.
	// If period covers all the content, the attribute is removed entirely.
	WithoutAttributeDuring(attribute string, period periods.Period) ContentBuilder
	// Cut reduces the content to a given period.
	// Typical use is to restrict attributes values to global content actvity.
	Cut(period periods.Period) ContentBuilder
	// Errors returns, if any, current errors so far.
	// Errors are cumulative
	Errors() error
	// Build creates a new content with the applied changes.
	// It returns the new content and an error if any occurred during the build process.
	// It resets the builder to its initial state, ready for new content modifications.
	// But the recommended use would be to create a new content with a new builder.
	Build() (TimeDependentContent, error)
}

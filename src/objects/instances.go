package objects

import (
	"errors"
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
// No structs (except time.Time), as they would allow bad design (use instance instead)
// Note that this code is actually dead : not used elsewhere.
// But it documents the allowed primitive types.
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
// In instances implementation, it is used to ensure that only primitive values are stored.
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
	// Same returns true if instance is the same as another TimeDependentValues.
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

// Instance defines attributes and time dependent values over time, with a global activity period.
// For instance, a person has a global life time,
// and that person's name, age, and address can be considered as attributes, while their values can change over time.
type Instance interface {
	Linkable // Linkable to become an operand of a link
	Element  // Element because it is a system entity
	// Activity returns the global activity period this instance lasts
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
	// Because it is a snapshot of the instance at a specific point in time,
	// result is a map of attribute names to their values at that moment.
	At(moment time.Time) (map[string]any, bool)
	// Matches returns, if any, the period during which this instance matches the given trait.
	// For instance, a person may have a student identity, and a student trait may match that identity during a specific period.
	// The returned period indicates the time frame during which the instance's attributes and values align with the trait's requirements.
	// If that given trait is incompatible with the instance, the result will be empty, false
	Matches(trait Trait) (periods.Period, bool)
}

// InstanceBuilder manages the changes to apply on a given instance.
// Typical use is to implement a load from existing instance, perform changes and build a new instance.
// Conventionally, it returns itself to allow method chaining.
type InstanceBuilder interface {
	// WithActivity changes the instance's activity to that specific period.
	WithActivity(period periods.Period) InstanceBuilder
	// WithAttributeDuring sets the attribute to the given value during the specified period.
	// Types for value are defined in PrimitiveValue.
	// If there is a type change, it should raise an error.
	// For instance, an age that contains 10 and "twenty" should raise an error.
	// Reasons are : storage, type safety, consistency
	WithAttributeDuring(attribute string, period periods.Period, value any) InstanceBuilder
	// WithoutAttributeDuring removes the attribute during the specified period.
	// If period covers all the instance, the attribute is removed entirely.
	WithoutAttributeDuring(attribute string, period periods.Period) InstanceBuilder
	// Cut reduces the instance to a given period.
	// Typical use is to restrict attributes values to global instance activity.
	Cut(period periods.Period) InstanceBuilder
	// Errors returns, if any, current errors so far.
	// Errors are cumulative
	Errors() error
	// Build creates a new instance with the applied changes.
	// It returns the new instance and an error if any occurred during the build process.
	// It resets the builder to its initial state, ready for new instance modifications.
	// But the recommended use would be to create a new instance with a new builder.
	Build() (Instance, error)
}

// =========================================================================
// TIME DEPENDENT VALUES IMPLEMENTATION : in memory, no storage
// =========================================================================

// valueNode stores a value set during a specific matchingPeriod
// value is the actual value (of type any) stored in the node.
type valueNode struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod periods.Period
	// value is the actual value stored in the node
	value any
}

// valuesHandler manages the full history of values with their respective matching periods.
// Its purpose is to provide a way to store and retrieve values over time.
// KEY INVARIANT : storedType is the actual type (should be primitive) and should be unique over time.
// There is NO LOCK at all, because it is immutable by design.
type valuesHandler struct {
	// values have one value per matching period
	values []valueNode
	// storedType is the actual type name of the content (should be primitive)
	storedType string
	// equality function
	equals func(a, b any) bool
}

// Same returns true if the two temporal values have the same values at the same periods, and same type
func (vh *valuesHandler) Same(other DynamicValues) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if vh.IsEmpty() {
		return true
	} else if vh.storedType != other.DataType() {
		return false
	}

	counter := 0
	for period, value := range other.Range {
		counter++
		found := false
		// find matching element if any
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
func (vh *valuesHandler) IsEmpty() bool {
	return vh == nil || len(vh.values) == 0
}

// Validity returns the union of periods for which values are set
func (vh *valuesHandler) Validity() periods.Period {
	if vh == nil || len(vh.values) == 0 {
		return periods.NewEmptyPeriod()
	}

	validity := periods.NewEmptyPeriod()
	for _, element := range vh.values {
		validity = validity.Union(element.matchingPeriod)
	}

	return validity
}

// At returns the value at the given moment in time, or nil and false if no value is found.
func (vh *valuesHandler) At(moment time.Time) (any, bool) {
	for _, element := range vh.values {
		if element.matchingPeriod.Contains(moment) {
			return element.value, true
		}
	}
	return nil, false
}

// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
func (vh *valuesHandler) Range(yield func(periods.Period, any) bool) {
	for _, element := range vh.values {
		if !yield(element.matchingPeriod, element.value) {
			break
		}
	}
}

// DataType returns the string representation of the type.
// It is UNIQUE by design : this is a key invariant.
func (vh *valuesHandler) DataType() string {
	return vh.storedType
}

// Copy returns a copy of the valuesHandler with the same values and type
func (vh *valuesHandler) Copy() *valuesHandler {
	result := make([]valueNode, len(vh.values))
	copy(result, vh.values)
	return &valuesHandler{values: result, storedType: vh.storedType, equals: vh.equals}
}

// withValueDuring adds a new value to a copy during a given period
func (vh *valuesHandler) withValueDuring(p periods.Period, v any) *valuesHandler {
	if !IsPrimitiveValue(v) {
		panic("cannot add value of incompatible type to valuesHandler")
	} else if primitiveTypeName(v) != vh.storedType {
		panic("cannot add value of incompatible type to valuesHandler")
	}

	matchingPeriodValue := p
	for _, element := range vh.values {
		if vh.equals(element.value, v) {
			matchingPeriodValue = matchingPeriodValue.Union(element.matchingPeriod)
		}
	}

	result := make([]valueNode, 0, len(vh.values)+1)
	for _, element := range vh.values {
		if !vh.equals(element.value, v) {
			remaining := element.matchingPeriod.Remove(matchingPeriodValue)
			if !remaining.IsEmpty() {
				result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
			}
		}
	}

	if !matchingPeriodValue.IsEmpty() {
		result = append(result, valueNode{matchingPeriod: matchingPeriodValue, value: v})
	}

	return &valuesHandler{values: result, storedType: vh.storedType, equals: vh.equals}
}

// withoutValidity returns a copy without values for the given period.
// If the period is empty or the handler is empty, it does nothing.
func (vh *valuesHandler) withoutValidity(period periods.Period) *valuesHandler {
	if len(vh.values) == 0 {
		return &valuesHandler{storedType: vh.storedType}
	} else if period.IsEmpty() {
		return vh
	}

	result := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Remove(period)
		if !remaining.IsEmpty() {
			result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: result, storedType: vh.storedType, equals: vh.equals}
}

// cut returns a copy with same values, restricted to a given period
func (vh *valuesHandler) cut(period periods.Period) *valuesHandler {
	remainingValues := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: remainingValues, storedType: vh.storedType, equals: vh.equals}
}

// newValuesHandler creates a new TemporalValues instance with a single value for the given period.
func newValuesHandler(period periods.Period, value any) *valuesHandler {
	if !IsPrimitiveValue(value) {
		panic("cannot create valuesHandler with non-primitive value")
	}
	typeName := primitiveTypeName(value)
	equalsForValue := primitiveTypeEqualsFunc(typeName)
	return &valuesHandler{
		equals:     equalsForValue,
		storedType: typeName,
		values:     []valueNode{{matchingPeriod: period, value: value}},
	}
}

// valuesHandlerLoad creates a new TemporalValues instance from another DynamicValues.
func valuesHandlerLoad(other DynamicValues) *valuesHandler {
	result := new(valuesHandler)
	result.storedType = other.DataType()
	result.equals = primitiveTypeEqualsFunc(result.storedType)
	for period, value := range other.Range {
		if !IsPrimitiveValue(value) {
			panic("cannot create valuesHandler with non-primitive value")
		}

		result.values = append(result.values, valueNode{matchingPeriod: period, value: value})
	}

	return result
}

// =========================================================================
// INSTANCE IMPLEMENTATION
// =========================================================================

// baseInstance is the in memory representation of an instance
type baseInstance struct {
	// id of the
	id string
	// activity defines when current instance is active (its current lifetime)
	activity periods.Period
	// values are the temporal values associated with their attributes names
	values map[string]*valuesHandler
}

// isLinkable is a SEALED INTERFACE pattern implementation.
// It allows instances to be linked to other elements.
func (b *baseInstance) isLinkable() bool {
	return true
}

// Same returns true if the instance is the same as the other element : same class, same id, same period, same values
func (b *baseInstance) Same(other Element) bool {
	if b == nil && other == nil {
		return true
	} else if b == nil || other == nil {
		return false
	} else if other.DeclaringClass() != CLASS_INSTANCE {
		return false
	}

	otherInstance, okInstance := other.(Instance)
	if !okInstance {
		return false
	} else if otherInstance.Id() != b.Id() {
		return false
	}

	if !b.activity.Equals(otherInstance.Activity()) {
		return false
	}

	counter := 0
	for name, content := range otherInstance.Values() {
		counter++
		if matching, found := b.values[name]; !found {
			return false
		} else if !matching.Same(content) {
			return false
		}
	}

	if len(b.values) != counter {
		return false
	}

	return true
}

// Activity returns the period during which the instance is valid
func (b *baseInstance) Activity() periods.Period {
	return b.activity
}

// Id returns the id of the instance
func (b *baseInstance) Id() string {
	return b.id
}

// DeclaringClass returns CLASS_INSTANCE to allow dynamic class discovery
func (b *baseInstance) DeclaringClass() Class {
	return CLASS_INSTANCE
}

// At returns the content at a given time, as a map of attributes and values.
// If instance is not active at moment, then it returns nil, false.
func (b *baseInstance) At(moment time.Time) (map[string]any, bool) {
	if b == nil {
		return nil, false
	} else if !b.activity.Contains(moment) {
		return nil, false
	}

	result := make(map[string]any)
	for attribute, content := range b.values {
		if value, exists := content.At(moment); exists {
			result[attribute] = value
		}
	}

	return result, true
}

// Matches tests if instance matches a given trait and returns the matching period.
// For instance, an instance has a name, an age and, during 5 years, a student id.
// Trait student may match on name and student id, but during 5 years only.
func (b *baseInstance) Matches(trait Trait) (periods.Period, bool) {
	if b == nil {
		return periods.NewEmptyPeriod(), false
	}

	matchingPeriod := b.activity
	for attribute, attributeType := range trait.Attributes() {
		// early test : leave when no match
		if matchingPeriod.IsEmpty() {
			return periods.NewEmptyPeriod(), false
		} else if matchingAttribute, exists := b.values[attribute]; !exists {
			return periods.NewEmptyPeriod(), false
		} else if attributeType != matchingAttribute.DataType() {
			return periods.NewEmptyPeriod(), false
		} else {
			matchingPeriod = matchingPeriod.Intersection(matchingAttribute.Validity())
		}
	}

	return matchingPeriod, !matchingPeriod.IsEmpty()
}

// Description returns a map of attribute names to their data types
func (b *baseInstance) Description() map[string]string {
	result := make(map[string]string)
	for attribute, content := range b.values {
		result[attribute] = content.DataType()
	}
	return result
}

// Values returns a copy of the temporal values associated with their attribute names
func (b *baseInstance) Values() map[string]DynamicValues {
	result := make(map[string]DynamicValues)
	for attribute, content := range b.values {
		result[attribute] = content
	}
	return result
}

// Value returns the temporal values associated with the given attribute name, if it exists
func (b *baseInstance) Value(attribute string) (DynamicValues, bool) {
	value, found := b.values[attribute]
	return value, found
}

// newBaseInstance returns an empty baseInstance
func newBaseInstance(id string) *baseInstance {
	return &baseInstance{
		id:       id,
		activity: periods.NewEmptyPeriod(),
		values:   make(map[string]*valuesHandler),
	}
}

// baseInstanceLoad creates a new baseInstance instance from a Instance.
// It performs a full copy : it imports the activity period and initializes the values map with loaded content
func baseInstanceLoad(other Instance) *baseInstance {
	result := new(baseInstance)
	result.id = other.Id()
	result.activity = other.Activity()
	result.values = make(map[string]*valuesHandler)
	for attribute, content := range other.Values() {
		result.values[attribute] = valuesHandlerLoad(content)
	}
	return result
}

// ===============================================================
// LOCAL BUILDER MANAGES IN MEMORY IMPLEMENTATION
// ===============================================================

// LocalInstanceBuilder manages in-memory instance creation and modification.
// It applies any changes to the instance and prepares it for building.
// It resets the builder to its initial state after building, ready for new instance modifications.
// Note the globalErrors field is used to accumulate errors during instance building.
type LocalInstanceBuilder struct {
	// element is the decorated instance to build
	element *baseInstance
	// globalErrors contain the global errors during instance building.
	// It accumulates errors that occur during instance building, allowing for comprehensive error handling.
	globalErrors error
}

// LocalInstanceBuilderLoad allows to read any instance and get ready for an in-memory instance rebuild.
// Two main use cases: modify existing instance or create a "in memory" instance from another implementation.
func LocalInstanceBuilderLoad(element Instance) InstanceBuilder {
	return &LocalInstanceBuilder{
		element: baseInstanceLoad(element),
	}
}

// NewLocalInstanceBuilder creates a new local builder with an empty base instance.
// Typical use case: create a new instance from scratch.
func NewLocalInstanceBuilder(id string) InstanceBuilder {
	return &LocalInstanceBuilder{
		element: newBaseInstance(id),
	}
}

// WithActivity sets the activity period for the instance being built.
// Although it makes no sense, it accepts empty periods.
// It returns the builder for method chaining.
func (b *LocalInstanceBuilder) WithActivity(period periods.Period) InstanceBuilder {
	b.element.activity = period
	return b
}

// WithAttributeDuring sets a value during a period for a given attribute.
// It validates the attribute and value types, and handles errors gracefully.
// It will add an error if the value is incompatible.
// It returns the builder for method chaining.
func (b *LocalInstanceBuilder) WithAttributeDuring(attribute string, period periods.Period, value any) InstanceBuilder {
	if value == nil || !IsPrimitiveValue(value) {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("attribute value cannot be nil or non-primitive"))
		return b
	} else if values, exists := b.element.values[attribute]; !exists {
		values = newValuesHandler(period, value)
		b.element.values[attribute] = values
	} else if primitiveTypeName(value) != values.storedType {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("cannot add value of incompatible type to valuesHandler"))
		return b
	} else {
		b.element.values[attribute] = values.withValueDuring(period, value)
	}

	return b
}

// WithoutAttributeDuring changes decorated instance to remove all values within that given period for that attribute.
// It returns the builder for method chaining.
func (b *LocalInstanceBuilder) WithoutAttributeDuring(attribute string, period periods.Period) InstanceBuilder {
	values, exists := b.element.values[attribute]
	if !exists {
		return b
	} else if period.IsEmpty() {
		return b
	}

	newValue := values.withoutValidity(period)
	if !newValue.IsEmpty() {
		b.element.values[attribute] = newValue
	} else {
		delete(b.element.values, attribute)
	}

	return b
}

// Cut reduces the whole instance (activity and values) to given period.
// It returns the builder for method chaining.
func (b *LocalInstanceBuilder) Cut(period periods.Period) InstanceBuilder {
	empty := &baseInstance{activity: periods.NewEmptyPeriod(), values: make(map[string]*valuesHandler)}
	if period.IsEmpty() {
		b.element = empty
		return b
	}

	remainingActivity := period.Intersection(b.element.activity)
	if remainingActivity.IsEmpty() {
		b.element = empty
		return b
	}

	valuesMap := make(map[string]*valuesHandler)
	for attribute, value := range b.element.values {
		newValue := value.cut(remainingActivity)
		if !newValue.IsEmpty() {
			valuesMap[attribute] = newValue
		}
	}

	b.element.values = nil
	b.element.values = valuesMap
	b.element.activity = remainingActivity
	return b
}

// Errors returns, if any, current errors so far.
// Errors are cumulative.
func (b *LocalInstanceBuilder) Errors() error { return b.globalErrors }

// Build returns the built instance and resets the builder for future use.
// It returns the builder for method chaining.
func (b *LocalInstanceBuilder) Build() (Instance, error) {
	result := b.element
	resultErr := b.globalErrors
	resultId := b.element.id
	b.element = newBaseInstance(resultId)
	b.globalErrors = nil
	if resultErr != nil {
		return nil, resultErr
	}

	return result, resultErr
}

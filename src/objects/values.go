package objects

import (
	"errors"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

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
func (vh *valuesHandler) Same(other TimeDependentValue) bool {
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

// valuesHandlerLoad creates a new TemporalValues instance from another TimeDependentValue.
func valuesHandlerLoad(other TimeDependentValue) *valuesHandler {
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
// CONTENT IMPLEMENTATION
// =========================================================================

// baseContent is the in memory representation of a content
type baseContent struct {
	// activity defines when content is valid and related instance was alive / active.
	activity periods.Period
	// values are the temporal values associated with their attributes names
	values map[string]*valuesHandler
}

// Same returns true if the content is the same as the other content : same period, same values
func (b *baseContent) Same(other TimeDependentContent) bool {
	if b == nil && other == nil {
		return true
	} else if b == nil || other == nil {
		return false
	}

	if !b.activity.Equals(other.Activity()) {
		return false
	}

	counter := 0
	for name, content := range other.Values() {
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

// Activity returns the period during which the content is valid
func (b *baseContent) Activity() periods.Period {
	return b.activity
}

// At returns the content at a given time, as a map of attributes and values.
// If content is not active at moment, then it returns nil, false.
func (b *baseContent) At(moment time.Time) (map[string]any, bool) {
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

// Matches tests if content matches a given trait and returns the matching period.
// For instance, a content has a name, an age and, during 5 years, a student id.
// Trait student may match on name and student id, but during 5 years only.
func (b *baseContent) Matches(trait Trait) (periods.Period, bool) {
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
func (b *baseContent) Description() map[string]string {
	result := make(map[string]string)
	for attribute, content := range b.values {
		result[attribute] = content.DataType()
	}
	return result
}

// Values returns a copy of the temporal values associated with their attributes names
func (b *baseContent) Values() map[string]TimeDependentValue {
	result := make(map[string]TimeDependentValue)
	for attribute, content := range b.values {
		result[attribute] = content
	}
	return result
}

// Value returns the temporal values associated with the given attribute name, if it exists
func (b *baseContent) Value(attribute string) (TimeDependentValue, bool) {
	value, found := b.values[attribute]
	return value, found
}

// newBaseContent returns an empty baseContent
func newBaseContent() *baseContent {
	return &baseContent{
		activity: periods.NewEmptyPeriod(),
		values:   make(map[string]*valuesHandler),
	}
}

// baseContentLoad creates a new baseContent instance from a TimeDependentContent.
// It performs a full copy : it imports the activity period and initializes the values map with loaded content handlers.
func baseContentLoad(other TimeDependentContent) *baseContent {
	result := new(baseContent)
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

// LocalContentBuilder manages in-memory content creation and modification.
// It applies any changes to the content and prepares it for building.
// It resets the builder to its initial state after building, ready for new content modifications.
// Note the globalErrors field is used to accumulate errors during content building.
type LocalContentBuilder struct {
	// element is the decorated content to build
	element *baseContent
	// globalErrors contain the global errors during content building.
	// It accumulates errors that occur during content building, allowing for comprehensive error handling.
	globalErrors error
}

// LocalContentBuilderLoad allows to read any content and get ready for an in-memory content rebuild.
// Two main use cases: modify existing content or create a "in memory" content from another implementation.
func LocalContentBuilderLoad(element TimeDependentContent) ContentBuilder {
	return &LocalContentBuilder{
		element: baseContentLoad(element),
	}
}

// NewLocalContentBuilder creates a new local content builder with an empty base content.
// Typical use case: create a new content from scratch.
func NewLocalContentBuilder() ContentBuilder {
	return &LocalContentBuilder{
		element: newBaseContent(),
	}
}

// WithActivity sets the activity period for the content being built.
// Although it makes no sense, it accepts empty periods.
// It returns the builder for method chaining.
func (b *LocalContentBuilder) WithActivity(period periods.Period) ContentBuilder {
	b.element.activity = period
	return b
}

// WithAttributeDuring sets a value during a period for a given attribute.
// It validates the attribute and value types, and handles errors gracefully.
// It will add an error if the value is incompatible.
// It returns the builder for method chaining.
func (b *LocalContentBuilder) WithAttributeDuring(attribute string, period periods.Period, value any) ContentBuilder {
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

// WithoutAttributeDuring changes decorated content to remove all values within that given period for that attribute.
// It returns the builder for method chaining.
func (b *LocalContentBuilder) WithoutAttributeDuring(attribute string, period periods.Period) ContentBuilder {
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

// Cut reduces the whole content (activity and values) to given period.
// It returns the builder for method chaining.
func (b *LocalContentBuilder) Cut(period periods.Period) ContentBuilder {
	empty := &baseContent{activity: periods.NewEmptyPeriod(), values: make(map[string]*valuesHandler)}
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
func (b *LocalContentBuilder) Errors() error { return b.globalErrors }

// Build returns the built content and resets the builder for future use.
// It returns the builder for method chaining.
func (b *LocalContentBuilder) Build() (TimeDependentContent, error) {
	result := b.element
	resultErr := b.globalErrors
	b.element = newBaseContent()
	b.globalErrors = nil
	if resultErr != nil {
		return b.element, resultErr
	}

	return result, resultErr
}

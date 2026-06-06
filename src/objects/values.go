package objects

import (
	"maps"
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// TemporalValues represents a collection of values with associated time periods.
// It uses "any" to store any type of values per period.
type TemporalValues interface {
	// Validity returns the period the values are set for.
	// Basically, it is empty for nil or empty, the union of periods for values otherwise
	Validity() periods.Period
	// Same returns true if content is the same as another TemporalValues.
	// It means : same periods, same values
	Same(other TemporalValues) bool
	// IsEmpty checks if the TemporalValues collection is empty (no value on a non empty period)
	IsEmpty() bool
	// Add a value for a given period
	Add(period periods.Period, value any) TemporalValues
	// At retrieves the value at a specific moment in time, if any
	At(moment time.Time) (any, bool)
	// Remove removes all values for a given period
	Remove(period periods.Period) TemporalValues
	// Cut returns a new TemporalValues collection containing only values within the specified period
	Cut(period periods.Period) TemporalValues
	// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
	Range(yield func(periods.Period, any) bool)
	// DataType returns the type of values stored in the TemporalValues collection.
	// It looks for the most common type among all values, or returns "any" if types are diverse.
	// For instance, if all values are integers, it will return "int". If there are both integers and strings, it will return "any".
	// Special case for empty collection: returns ""
	DataType() string
}

// Content is the historized state of an instance.
type Content interface {
	// Same tests if other is the same as the current content (same values, same activity, same description)
	Same(other Content) bool
	// Activity returns the period during which the content is valid
	Activity() periods.Period
	// Add a value for an attribute during a given period
	Add(attribute string, period periods.Period, value any) Content
	// Description returns the metadata of the content : attributes and their types
	Description() map[string]string
	// Values returns the collection of TemporalValues, which are the values with their respective matching periods
	Values() map[string]TemporalValues
	// Value returns the TemporalValues associated with the given attribute, or false if not found
	Value(string) (TemporalValues, bool)
	// Remove removes all values for an attribute during a given period.
	// If period completely covers all existing values, the attribute is removed
	Remove(attribute string, period periods.Period) Content
	// Cut returns a new Content containing only values within the specified period.
	// Note that it cuts the full content : active period and attributes !
	// If the period does not overlap with any existing values, returns an empty content
	Cut(period periods.Period) Content
	// At returns the content at a given time, as a map of attributes and values.
	// Only values with a matching period containing the given time are included.
	// If content does not exist at the given time, returns an empty map and false
	At(time time.Time) (map[string]any, bool)
}

// =========================================================================
// TEMPORAL VALUES IMPLEMENTATION
// =========================================================================

// valueNode stores a value set during a specific matchingPeriod
// value is the actual value (of type any) stored in the node
type valueNode struct {
	// matchingPeriod is the period during which the value is valid
	matchingPeriod periods.Period
	// value is the actual value stored in the node
	value any
}

// valuesHandler manages the full history of values with their respective matching periods
type valuesHandler struct {
	// values have one value per matching period
	values []valueNode
}

// Same returns true if the two temporal values have the same values at the same periods and same id
func (vh *valuesHandler) Same(other TemporalValues) bool {
	if vh == nil && other == nil {
		return true
	} else if vh == nil || other == nil {
		return false
	} else if vh.IsEmpty() != other.IsEmpty() {
		return false
	} else if vh.IsEmpty() {
		return true
	}

	counter := 0
	for period, value := range other.Range {
		counter++
		found := false
		// find matching element if any
		for _, matching := range vh.values {
			if period.Equals(matching.matchingPeriod) {
				found = true
				if !reflect.DeepEqual(matching.value, value) {
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

// Add adds a new value with a specific matchingPeriod to the valuesHandler
func (vh *valuesHandler) Add(p periods.Period, v any) TemporalValues {
	matchingPeriodValue := p
	for _, element := range vh.values {
		if reflect.DeepEqual(element.value, v) {
			matchingPeriodValue = matchingPeriodValue.Union(element.matchingPeriod)
		}
	}

	result := make([]valueNode, 0, len(vh.values)+1)
	for _, element := range vh.values {
		if !reflect.DeepEqual(element.value, v) {
			remaining := element.matchingPeriod.Remove(matchingPeriodValue)
			if !remaining.IsEmpty() {
				result = append(result, valueNode{matchingPeriod: remaining, value: element.value})
			}
		}
	}

	if !matchingPeriodValue.IsEmpty() {
		result = append(result, valueNode{matchingPeriod: matchingPeriodValue, value: v})
	}

	return &valuesHandler{values: result}
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

// Remove removes the given period from the values handler, if the period is empty or the handler is empty, it does nothing.
func (vh *valuesHandler) Remove(period periods.Period) TemporalValues {
	if len(vh.values) == 0 {
		return &valuesHandler{}
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

	return &valuesHandler{values: result}
}

// Range iterates over all values in the TemporalValues collection, yielding each period and value to a provided function
func (vh *valuesHandler) Range(yield func(periods.Period, any) bool) {
	for _, element := range vh.values {
		if !yield(element.matchingPeriod, element.value) {
			break
		}
	}
}

// Cut returns a copy with same values, restricted to given period
func (vh *valuesHandler) Cut(period periods.Period) TemporalValues {
	remainingValues := make([]valueNode, 0, len(vh.values))
	for _, element := range vh.values {
		remaining := element.matchingPeriod.Intersection(period)
		if !remaining.IsEmpty() {
			remainingValues = append(remainingValues, valueNode{matchingPeriod: remaining, value: element.value})
		}
	}

	return &valuesHandler{values: remainingValues}
}

// DataType returns the string representation of the common type of all stored values or "any" if types differ.
func (vh *valuesHandler) DataType() string {
	if vh == nil || len(vh.values) == 0 {
		return ""
	}

	var commonType string
	isFirst := true

	for _, element := range vh.values {
		currentType := reflect.TypeOf(element.value).String()

		if isFirst {
			commonType = currentType
			isFirst = false
			continue
		}

		if currentType != commonType {
			return "any"
		}
	}

	return commonType
}

// buildTemporalValues creates a new TemporalValues instance with a single value for the given period.
func buildTemporalValues(period periods.Period, value interface{}) TemporalValues {
	return &valuesHandler{
		values: []valueNode{{matchingPeriod: period, value: value}},
	}
}

// NewTemporalValues creates a new TemporalValues instance with no values.
func NewTemporalValues() TemporalValues {
	return &valuesHandler{}
}

// =========================================================================
// CONTENT IMPLEMENTATION
// =========================================================================

// baseContent is the in memory representation of a content
type baseContent struct {
	// activity defines when content is valid and related instance was alive / active.
	activity periods.Period
	// values are the temporal values associated with their attributes names
	values map[string]TemporalValues
}

// Same returns true if the content is the same as the other content : same period, same values
func (b *baseContent) Same(other Content) bool {
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

// Description returns a map of attribute names to their data types
func (b *baseContent) Description() map[string]string {
	result := make(map[string]string)
	for attribute, content := range b.values {
		result[attribute] = content.DataType()
	}
	return result
}

// Values returns a copy of the temporal values associated with their attributes names
func (b *baseContent) Values() map[string]TemporalValues {
	result := make(map[string]TemporalValues)
	maps.Copy(result, b.values)
	return result
}

// Value returns the temporal values associated with the given attribute name, if it exists
func (b *baseContent) Value(attribute string) (TemporalValues, bool) {
	value, found := b.values[attribute]
	return value, found
}

// Add adds a new temporal value to the content for the given attribute and period, returning a new Content instance.
// It basically looks for the previous values for the given attribute, and uses the temporal values handler to add the new value.
func (b *baseContent) Add(attribute string, period periods.Period, value any) Content {
	if b == nil {
		return nil
	}

	valuesMap := make(map[string]TemporalValues)
	maps.Copy(valuesMap, b.values)

	if values, exists := valuesMap[attribute]; !exists {
		values = buildTemporalValues(period, value)
		valuesMap[attribute] = values
	} else {
		valuesMap[attribute] = values.Add(period, value)
	}

	return &baseContent{
		activity: b.activity,
		values:   valuesMap,
	}
}

// Remove reduces the attribute values to exclude a period.
// If all values are excluded, the attribute itself is removed.
func (b *baseContent) Remove(attribute string, period periods.Period) Content {
	if b == nil {
		return nil
	}

	values, exists := b.values[attribute]
	if !exists {
		return b
	} else {
		newValues := make(map[string]TemporalValues)
		maps.Copy(newValues, b.values)
		newValue := values.Remove(period)
		if !newValue.IsEmpty() {
			newValues[attribute] = newValue
		} else {
			delete(newValues, attribute)
		}

		return &baseContent{
			activity: b.activity,
			values:   newValues,
		}
	}

}

// Cut reduces the content to only include values within the specified period.
// It means reducing the content's lifetime, and the attributes values to those that are active within the given period.
// If content is not active at all during that period, it returns an empty content
func (b *baseContent) Cut(period periods.Period) Content {
	if b == nil {
		// cut on empty => no possible match
		return nil
	}
	if period.IsEmpty() {
		return &baseContent{activity: periods.NewEmptyPeriod()}
	}

	remainingActivity := period.Intersection(b.activity)
	if remainingActivity.IsEmpty() {
		return &baseContent{activity: periods.NewEmptyPeriod()}
	}

	valuesMap := make(map[string]TemporalValues)
	for attribute, value := range b.values {
		newValue := value.Cut(remainingActivity)
		if !newValue.IsEmpty() {
			valuesMap[attribute] = newValue
		}
	}
	return &baseContent{
		activity: remainingActivity,
		values:   valuesMap,
	}
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

// NewContent returns a new empty content. Default lifetime is full period.
func NewContent() Content {
	return &baseContent{
		activity: periods.NewFullPeriod(),
		values:   make(map[string]TemporalValues),
	}
}

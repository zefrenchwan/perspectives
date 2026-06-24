package objects

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zefrenchwan/perspectives.git/periods"
)

type EntityBuilder interface {
	WithActivity(period periods.Period) EntityBuilder

	WithAttributeDuring(attribute string, period periods.Period, value any) EntityBuilder

	WithoutAttributeDuring(attribute string, period periods.Period) EntityBuilder

	Cut(period periods.Period) EntityBuilder

	WithOperand(role string, operand Entity) EntityBuilder

	WithoutOperand(role string) EntityBuilder

	Errors() error

	Build() (Entity, error)
}

type localEntityBuilder struct {
	element      *localEntity
	globalErrors error
}

func LocalEntityBuilderLoad(element Entity) EntityBuilder {
	return &localEntityBuilder{
		element: localEntityLoad(element),
	}
}

func NewLocalEntityBuilder(id string) EntityBuilder {
	var globalErrors error
	if id == "" {
		globalErrors = errors.New("entity id cannot be empty")
	}

	return &localEntityBuilder{
		globalErrors: globalErrors,
		element:      newLocalEntity(id),
	}
}

// WithActivity sets the activity period for the entity being built.
// Although it makes no sense, it accepts empty periods.
// It returns the builder for method chaining.
func (b *localEntityBuilder) WithActivity(period periods.Period) EntityBuilder {
	b.element.activity = period
	return b
}

// WithAttributeDuring sets a value during a period for a given attribute.
// It validates the attribute and value types, and handles errors gracefully.
// It will add an error if the value is incompatible.
// It returns the builder for method chaining.
func (b *localEntityBuilder) WithAttributeDuring(attribute string, period periods.Period, value any) EntityBuilder {
	if value == nil || !IsPrimitiveValue(value) {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("attribute value cannot be nil or non-primitive"))
		return b
	} else if existingHandler, exists := b.element.values[attribute]; !exists {
		typeName := primitiveTypeName(value)
		equalsForValue := primitiveTypeEqualsFunc(typeName)
		existingHandler = &valuesHandler{
			equals:     equalsForValue,
			storedType: typeName,
			values:     []valueNode{{matchingPeriod: period, value: value}},
		}

		b.element.values[attribute] = existingHandler
	} else if primitiveTypeName(value) != existingHandler.storedType {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("cannot add value of incompatible type to valuesHandler"))
		return b
	} else {
		// value is OK, we already have a matching attribute and related mapping.
		// At this point: values is the valuesHandler for the attribute
		matchingPeriodValue := period
		for existingPeriod, existingValue := range existingHandler.Range {
			if existingHandler.equals(existingValue, value) {
				matchingPeriodValue = matchingPeriodValue.Union(existingPeriod)
			}
		}

		result := make([]valueNode, 0)
		for existingPeriod, existingValue := range existingHandler.Range {
			if !existingHandler.equals(existingValue, value) {
				remaining := existingPeriod.Remove(matchingPeriodValue)
				if !remaining.IsEmpty() {
					result = append(result, valueNode{matchingPeriod: remaining, value: existingValue})
				}
			}
		}

		if !matchingPeriodValue.IsEmpty() {
			result = append(result, valueNode{matchingPeriod: matchingPeriodValue, value: value})
		}

		existingHandler.values = result
		b.element.values[attribute] = existingHandler
	}

	return b
}

// WithoutAttributeDuring changes decorated entity to remove all values within that given period for that attribute.
// It returns the builder for method chaining.
func (b *localEntityBuilder) WithoutAttributeDuring(attribute string, period periods.Period) EntityBuilder {
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

// Cut reduces the whole entity (activity and values) to given period.
// It does NOT change the other elements on links.
// It returns the builder for method chaining.
func (b *localEntityBuilder) Cut(period periods.Period) EntityBuilder {
	empty := &localEntity{id: b.element.id, activity: periods.NewEmptyPeriod(), values: make(map[string]*valuesHandler)}
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

// WithOperand adds an operand to the entity to build.
// Role is the key of the operand to add, operand is the actual value to add
func (l *localEntityBuilder) WithOperand(role string, operand Entity) EntityBuilder {
	if l.element.roles == nil {
		l.element.roles = make(map[string]Entity)
	}

	if operand == nil {
		l.globalErrors = errors.Join(l.globalErrors, fmt.Errorf("operand cannot be nil for %s", role))
		return l
	} else if strings.TrimSpace(role) == "" {
		l.globalErrors = errors.Join(l.globalErrors, fmt.Errorf("role cannot be empty"))
		return l
	}

	l.element.roles[role] = operand
	return l
}

// WithoutOperand removes the given role and related content
func (l *localEntityBuilder) WithoutOperand(role string) EntityBuilder {
	// role may be empty, no problem.
	// We may raise an error, but operation with empty role does not create an error
	if l.element.roles != nil {
		delete(l.element.roles, role)
	}
	return l
}

// Errors returns, if any, current errors so far.
// Errors are cumulative.
func (b *localEntityBuilder) Errors() error { return b.globalErrors }

// Build returns the built entity and resets the builder for future use.
// It returns the builder for method chaining.
func (b *localEntityBuilder) Build() (Entity, error) {
	// no check on attributes (might not have) or roles (might not have)
	// But id is mandatory
	if b.element.id == "" {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("no id defined for link"))
	}

	result := b.element
	resultErr := b.globalErrors
	resultId := b.element.id
	b.element = newLocalEntity(resultId)
	b.globalErrors = nil
	if resultErr != nil {
		return nil, resultErr
	}

	result.hashString = hashEntity(result)
	return result, resultErr
}

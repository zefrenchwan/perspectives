package objects

import (
	"maps"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

type Instance interface {
	IdentifiableElement
	Activity() periods.Period
	Description() map[string]string
	Values() map[string]TemporalValues
	Value(string) (TemporalValues, bool)
}

// =========================================================================
// ENTITY IMPLEMENTATION
// =========================================================================

type baseInstance struct {
	id       string
	activity periods.Period
	values   map[string]TemporalValues
}

func (b *baseInstance) Id() string {
	return b.id
}

func (b *baseInstance) Same(other Element) bool {
	if b == nil && other == nil {
		return true
	} else if b == nil || other == nil {
		return false
	}

	otherInstance, ok := other.(Instance)
	if !ok {
		return false
	} else if otherInstance == nil {
		return false
	} else if b.id != otherInstance.Id() {
		return false
	}

	// Same instance, may not be same content
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

func (b *baseInstance) DeclaringClass() Class {
	return CLASS_INSTANCE
}

func (b *baseInstance) Activity() periods.Period {
	return b.activity
}

func (b *baseInstance) Description() map[string]string {
	result := make(map[string]string)
	for attribute, content := range b.values {
		result[attribute] = content.DataType()
	}
	return result
}

func (b *baseInstance) Values() map[string]TemporalValues {
	result := make(map[string]TemporalValues)
	maps.Copy(result, b.values)
	return result
}

func (b *baseInstance) Value(attribute string) (TemporalValues, bool) {
	value, found := b.values[attribute]
	return value, found
}

func NewInstance() Instance {
	return &baseInstance{
		id:       commons.NewId(),
		activity: periods.NewFullPeriod(),
		values:   make(map[string]TemporalValues),
	}
}

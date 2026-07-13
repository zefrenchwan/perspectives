package graphs

import (
	"iter"
	"strconv"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// Property defines a characteristic of an element.
// It might be its age, the links it is related to, or any other attribute or role.
type Property[V values.Value] interface {
	// Property is immutable, so we have a hash on it.
	commons.Hashable
	// Name of the property : a role, or an attribute ("height", "age", "name", etc.)
	Name() string
	// Values of the property, as a sequence of periods and values.
	// It allows multi values and single value at a given time.
	Values() iter.Seq2[periods.Period, V]
}

// Attribute describes a state for an element via a primitive value.
type Attribute Property[values.PrimitiveValue]

// Role describes a relationship between two elements.
type Role Property[values.ReferenceValue]

// valuesProperty describes a property in the most generic way : name, mapping (and calculated hash).
type valuesProperty[V values.Value] struct {
	// name of the property : a role, or an attribute ("height", "age", "name", etc.)
	name string
	// hashString is the hash of the property : calculated from its name and mapping
	hashString string
	// mapping is the dynamic mapping of the property : time dependent values on codomain V
	mapping periods.DynamicMapping[V]
}

// Name of the property
func (p valuesProperty[V]) Name() string {
	return p.name
}

// Values of the property.
// It allows multiple values per period within the iterator
func (p valuesProperty[V]) Values() iter.Seq2[periods.Period, V] {
	return p.mapping.Range()
}

// ToHashString returns the hash string of the property (precalculated)
func (p valuesProperty[V]) ToHashString() string {
	return p.hashString
}

// NewSingleValueProperty creates a property with zero or one value at a time.
// Given a time t, it returns a single value if any.
func NewSingleValueProperty[V values.Value](name string, function periods.DynamicFunction[V]) Property[V] {
	// defensive copy to avoid side effects on that function
	newFunction := function.Copy()
	value := strconv.Itoa(len(name)) + "|" + name + "|" + periods.HashDynamicFunction(newFunction)
	hashString := commons.HashString(value)
	return valuesProperty[V]{
		name:       name,
		hashString: hashString,
		mapping:    newFunction,
	}
}

// NewMultiValuesProperty creates a property allowing multiple values per period.
func NewMultiValuesProperty[V values.Value](name string, relation periods.DynamicRelation[V]) Property[V] {
	// defensive copy to avoid side effects on that relation
	newRelation := relation.Copy()
	value := strconv.Itoa(len(name)) + "|" + name + "|" + periods.HashDynamicRelation(newRelation)
	hashString := commons.HashString(value)
	return valuesProperty[V]{
		name:       name,
		hashString: hashString,
		mapping:    newRelation,
	}
}

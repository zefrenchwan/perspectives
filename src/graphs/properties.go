package graphs

import (
	"iter"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// Property defines a characteristic of an element.
// It might be its age, the links it is related to, or any other attribute or role.
type Property[V values.Value] interface {
	// Hashable because Property is immutable, so we have a hash on it.
	commons.Hashable
	// Name of the property: a role, or an attribute ("height", "age", "name", etc.)
	Name() string
	// Values of the property, as a sequence of periods and values.
	// It allows multi values and single value at a given time.
	Values() iter.Seq2[periods.Period, V]
}

// Attribute describes a state for an element via a primitive value.
type Attribute Property[values.PrimitiveValue]

// Role describes a relationship between two elements.
type Role Property[values.ReferenceValue]

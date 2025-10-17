package commons

import (
	"errors"
	"maps"
)

type Linkable any

// Link is a constant relation over instances of linkables.
// Link is also Linkable, so it may be used in links.
type Link[T Linkable] interface {
	// Links are information about elements, they are then part of a model
	Modelable
	// links are composable due to this: a link is linkable
	Linkable
	// IsEmpty returns true if link should be empty and then not being used
	IsEmpty() bool
	// Name returns the name of the link.
	// It is usually a verb or a noun.
	// For instance, knows, couple, etc.
	Name() string
	// Roles returns all the roles set for that link
	Roles() []string
	// Operands returns the roles and related values for that link
	Operands() map[string]T
	// Get returns, if any, the elements for that name.
	// First result is the values if any, second is true if value was found
	Get(role string) (T, bool)
}

// TemporalLink decorates a link over a period.
// It should implement both link and Temporal.
type TemporalLink[T Linkable] struct {
	// activity of the link
	period Period
	// value is link decoration
	value Link[T]
}

// NewTemporalLink decorates a link true for given duration
func NewTemporalLink[T Linkable](duration Period, value Link[T]) *TemporalLink[T] {
	result := new(TemporalLink[T])
	result.period = duration
	result.value = value
	return result
}

// ActivePeriod is the duration which the link is true
func (t *TemporalLink[T]) ActivePeriod() Period {
	if t == nil {
		return NewEmptyPeriod()
	}

	return t.period
}

// SetActivePeriod forces active period
func (t *TemporalLink[T]) SetActivePeriod(period Period) {
	if t != nil {
		t.period = period
	}
}

// Name returns the name of the link.
func (t *TemporalLink[T]) Name() string {
	var empty string
	if t == nil {
		return empty
	}

	return t.Name()
}

// Roles returns all the roles set for that link
func (t *TemporalLink[T]) Roles() []string {
	if t == nil {
		return nil
	}

	return t.Roles()
}

// Size returns the number of elements per role.
func (t *TemporalLink[T]) Size(role string) int {
	if t == nil {
		return 0
	}

	return t.Size(role)
}

// Operands returns the roles and related values for that link
func (t *TemporalLink[T]) Operands() map[string]T {
	if t == nil {
		return nil
	}

	return t.Operands()
}

// Get returns, if any, the element for that name.
func (t *TemporalLink[T]) Get(role string) (T, bool) {
	var empty T
	if t == nil {
		return empty, false
	}

	return t.Get(role)
}

// First returns the first value, if any, for that role
func (t *TemporalLink[T]) First(role string) (T, bool) {
	var empty T
	if t == nil {
		return empty, false
	}

	return t.First(role)
}

// simpleLink implements links as its canonical implementation
type simpleLink[T Linkable] struct {
	name   string
	values map[string]T
}

// IsEmpty tests if link is empty (no name or no value)
func (s simpleLink[T]) IsEmpty() bool {
	return len(s.values) == 0 || s.name == ""
}

// GetType acts the fact that a link is a model link
func (s simpleLink[T]) GetType() ModelableType {
	return TypeLink
}

// Name returns the name of the link.
func (s simpleLink[T]) Name() string {
	return s.name
}

// Roles returns all the roles set for that link
func (s simpleLink[T]) Roles() []string {
	var result []string
	for role := range s.values {
		result = append(result, role)
	}

	return result
}

// Operands returns the roles and related values for that link.
// To avoid side effects, we return a copy (not direct access to values)
func (s simpleLink[T]) Operands() map[string]T {
	result := make(map[string]T)
	maps.Copy(result, s.values)
	return result
}

// Get returns, if any, the elements for that name.
func (s simpleLink[T]) Get(role string) (T, bool) {
	var empty T
	if result, found := s.values[role]; found {
		return result, true
	}

	return empty, false
}

// NewLink builds a new link, or raises an error if link would be malformed.
// A valid link is not empty: non empty name and at least one value.
// Of course, a "creative" user may create a link with " " name, but it is discouraged
func NewLink[T Linkable](name string, values map[string]T) (Link[T], error) {
	var empty string
	if len(values) == 0 {
		return nil, errors.New("no value for roles")
	} else if name == empty {
		return nil, errors.New("no name for link")
	}

	var result simpleLink[T]
	result.name = name
	result.values = make(map[string]T)
	maps.Copy(result.values, values)

	return result, nil
}

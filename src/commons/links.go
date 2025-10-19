package commons

import (
	"errors"
)

// RoleSubject specifies the subject role
const RoleSubject = "subject"

// RoleObject specifies the object role
const RoleObject = "object"

// Linkable should be as simple as possible.
type Linkable interface{}

// Link is a constant relation over instances of linkables.
// Link is also Linkable, so it may be used in links.
type Link interface {
	// Links have an id, new each time we build one link
	Identifiable
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
	Operands() map[string]Linkable
	// Get returns, if any, the elements for that name.
	// First result is the values if any, second is true if value was found
	Get(role string) (Linkable, bool)
}

// TemporalLink decorates a link over a period.
// It should implement both link and Temporal.
type TemporalLink struct {
	// id: not the same as the underlying link
	id string
	// activity of the link
	period Period
	// value is link decoration
	value Link
}

// NewTemporalLink decorates a link true for given duration
func NewTemporalLink(duration Period, value Link) *TemporalLink {
	result := new(TemporalLink)
	result.id = NewId()
	result.period = duration
	result.value = value
	return result
}

// Id() returns the id of the temporal link, not the same as underlying
func (t *TemporalLink) Id() string {
	return t.id
}

// GetType flags temporal link as a link
func (t *TemporalLink) GetType() ModelableType {
	return TypeLink
}

// ActivePeriod is the duration which the link is true
func (t *TemporalLink) ActivePeriod() Period {
	if t == nil {
		return NewEmptyPeriod()
	}

	return t.period
}

// SetActivePeriod forces active period
func (t *TemporalLink) SetActivePeriod(period Period) {
	if t != nil {
		t.period = period
	}
}

// Name returns the name of the link.
func (t *TemporalLink) Name() string {
	var empty string
	if t == nil {
		return empty
	} else if t.value == nil {
		return empty
	}

	return t.value.Name()
}

// Roles returns all the roles set for that link
func (t *TemporalLink) Roles() []string {
	if t == nil {
		return nil
	} else if t.value == nil {
		return nil
	}

	return t.value.Roles()
}

// Operands returns the roles and related values for that link
func (t *TemporalLink) Operands() map[string]Linkable {
	if t == nil {
		return nil
	} else if t.value == nil {
		return nil
	}

	return t.value.Operands()
}

// Get returns, if any, the element for that name.
func (t *TemporalLink) Get(role string) (Linkable, bool) {
	if t == nil || t.value == nil {
		return nil, false
	}

	return t.value.Get(role)
}

// simpleLinkNode decorates a value to ensure that it has an unique id (even for similar values)
type simpleLinkNode struct {
	// id should be unique
	id string
	// node is the decorated value
	node Linkable
}

// simpleLink implements links as its canonical implementation
type simpleLink struct {
	// id of the link
	id string
	// name of the link
	name string
	// values map roles to decorated linkable values
	values map[string]simpleLinkNode
}

// Id returns the link unique id
func (s simpleLink) Id() string {
	return s.id
}

// IsEmpty tests if link is empty (no name or no value)
func (s simpleLink) IsEmpty() bool {
	return len(s.values) == 0 || s.name == ""
}

// GetType acts the fact that a link is a model link
func (s simpleLink) GetType() ModelableType {
	return TypeLink
}

// Name returns the name of the link.
func (s simpleLink) Name() string {
	return s.name
}

// Roles returns all the roles set for that link
func (s simpleLink) Roles() []string {
	var result []string
	for role := range s.values {
		result = append(result, role)
	}

	return result
}

// Operands returns the roles and related values for that link.
// To avoid side effects, we return a copy (not direct access to values)
func (s simpleLink) Operands() map[string]Linkable {
	result := make(map[string]Linkable)
	for role, container := range s.values {
		result[role] = container.node
	}

	return result
}

// Get returns, if any, the elements for that name.
func (s simpleLink) Get(role string) (Linkable, bool) {
	if result, found := s.values[role]; found {
		return result.node, true
	}

	return nil, false
}

// NewLink builds a new link, or raises an error if link would be malformed.
// A valid link is not empty: non empty name and at least one value.
// Of course, a "creative" user may create a link with " " name, but it is discouraged
func NewLink(name string, values map[string]Linkable) (Link, error) {
	var empty string
	if len(values) == 0 {
		return nil, errors.New("no value for roles")
	} else if name == empty {
		return nil, errors.New("no name for link")
	}

	var result simpleLink
	result.id = NewId()
	result.name = name
	result.values = make(map[string]simpleLinkNode)
	for role, value := range values {
		decorated := simpleLinkNode{id: NewId(), node: value}
		result.values[role] = decorated
	}

	return result, nil
}

// NewSimpleLink is just creating a link of given name with content equals to "subject" => subject, "object": object
func NewSimpleLink(name string, subject Linkable, object Linkable) (Link, error) {
	return NewLink(name, map[string]Linkable{RoleSubject: subject, RoleObject: object})
}

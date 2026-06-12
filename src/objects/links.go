package objects

import (
	"errors"
	"fmt"
	"slices"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// Linkable defines any link operand.
// It should be :
// EITHER A LINK (link composition)
// OR AN INSTANCE (instance as operand),
// OR A TRAIT (trait as operand)
// OR A VARIABLE (for pattern matching).
// To do so, we use the sealed interface :
// we include a private function to force that no other type can implement it.
// This way, linkable types can only be implemented within this package.
// VERY IMPORTANT : Linkable should be immutable.
type Linkable interface {
	// isLinkable is a private function to force that no other type can implement it.
	// It is used as an implementation of the SEALED INTERFACE go pattern.
	isLinkable() bool
}

// sameLinkable checks if two linkables are the same.
func sameLinkable(a, b Linkable) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	elemA, okA := a.(Element)
	elemB, okB := b.(Element)

	if okA && okB {
		return elemA.Same(elemB)
	}

	return false
}

// Link relates elements together during a given period.
// For instance, FriendOf(subject=Marie,object=Paul) since now() - 3 years is a link.
type Link interface {
	Linkable // Linkable to use a link as a link operand (compositions)
	Element  // Element to use links as base components of the system
	// Name of the link, it defines its semantic
	Name() string
	// Roles associated to the link, to define how the elements are related together.
	// Although it is not mandatory, it is recommended to sort result.
	Roles() []string
	// Role returns the element associated to the given role, if any
	Role(string) (Linkable, bool)
	// Activity returns the period during which the link is active
	Activity() periods.Period
	// Range iterates over the roles and their associated elements
	Range(func(string, Linkable) bool)
}

// LinkBuilder is used to build a link.
// Principle is to declare a new builder, then fill it, then build it.
// Note that errors are cumulative :
// if there is ONE error once you build the link, all errors are returned.
type LinkBuilder interface {
	// WithOperand adds an operand with the given role
	WithOperand(role string, operand Linkable) LinkBuilder
	// WithoutRole removes the given role and related content
	WithoutRole(role string) LinkBuilder
	// WithActivity changes the activity period for the link
	WithActivity(period periods.Period) LinkBuilder
	// WithName changes the name of the link
	WithName(name string) LinkBuilder
	// Errors returns, if any, current errors so far.
	// Errors are cumulative
	Errors() error
	// Build the link or raise an error
	Build() (Link, error)
}

// localLink defines an immutable link with in memory information
type localLink struct {
	// id of the link
	id string
	// name of the link
	name string
	// roles of the link : for each name, linkable
	roles map[string]Linkable
	// activity period of the link
	activity periods.Period
}

// isLinkable uses sealed pattern to ensure that localLink instances can satisfy the Linkable interface requirements.
func (l *localLink) isLinkable() bool {
	return true
}

// Id returns the id of the link
func (l *localLink) Id() string { return l.id }

// Name returns the name of the link.
// Consider it to be its semantic
func (l *localLink) Name() string { return l.name }

// DeclaringClass returns the declaring class of the link : a CLASS_LINK
func (l *localLink) DeclaringClass() Class {
	return CLASS_LINK
}

// Same returns true if the link is the same as the other element
// It means same id, name, activity period and roles
func (l *localLink) Same(other Element) bool {
	if l == nil && other == nil {
		return true
	} else if l == nil || other == nil {
		return false
	}

	if otherLink, ok := other.(Link); ok {
		if l.Id() != otherLink.Id() {
			return false
		} else if otherLink.Name() != l.Name() {
			return false
		} else if !otherLink.Activity().Equals(l.Activity()) {
			return false
		}

		counter := 0
		for otherRole, otherLinkable := range otherLink.Range {
			if value, found := l.roles[otherRole]; !found {
				return false
			} else if !sameLinkable(value, otherLinkable) {
				return false
			}

			counter++
		}

		if counter != len(l.roles) {
			return false
		}

		return true
	}

	return false
}

// Activity returns the activity period of the link
func (l *localLink) Activity() periods.Period {
	return l.activity
}

// Roles returns the name of each role of the link
func (l *localLink) Roles() []string {
	result := make([]string, len(l.roles))
	index := 0
	for role := range l.roles {
		result[index] = role
		index++
	}

	slices.Sort(result)
	return result
}

// Role returns, if any, the linkable value associated to the role
func (l *localLink) Role(role string) (Linkable, bool) {
	linkable, found := l.roles[role]
	return linkable, found
}

// Range iterates over the roles of the link and yields each role and its associated linkable value
func (l *localLink) Range(yield func(string, Linkable) bool) {
	for role, linkable := range l.roles {
		if !yield(role, linkable) {
			return
		}
	}
}

// localLinkBuilder is a builder for links that stay in memory
type localLinkBuilder struct {
	// id of the link
	id string
	// name of the link
	name string
	// roles mapping from role names to linkable value
	roles map[string]Linkable
	// activity period for the link
	activity periods.Period
	// globalErrors when building (cumulative)
	globalErrors error
}

// NewLocalLinkBuilder creates a new local link builder with the given id
// Id will not change over time.
func NewLocalLinkBuilder(id string) LinkBuilder {
	return &localLinkBuilder{
		id:    id,
		roles: make(map[string]Linkable),
	}
}

// LocalLinkBuilderLoad creates a new local link builder with the given link
// It copies the full content roles from the original link.
// Due to the immutable nature of links, the roles are copied to a new map and the rest is passed as is.
func LocalLinkBuilderLoad(original Link) LinkBuilder {
	if original == nil {
		return &localLinkBuilder{
			id: "", name: "",
			roles:        nil,
			activity:     periods.NewEmptyPeriod(),
			globalErrors: errors.New("nil value"),
		}
	}

	newRoles := make(map[string]Linkable)
	for role, linkable := range original.Range {
		newRoles[role] = linkable
	}

	return &localLinkBuilder{
		id:       original.Id(),
		name:     original.Name(),
		activity: original.Activity(),
		roles:    newRoles,
	}
}

// WithOperand adds an operand with the given role
func (l *localLinkBuilder) WithOperand(role string, operand Linkable) LinkBuilder {
	if l.roles == nil {
		l.roles = make(map[string]Linkable)
	}

	if operand == nil {
		l.globalErrors = errors.Join(l.globalErrors, fmt.Errorf("operand cannot be nil for %s", role))
		return l
	}

	l.roles[role] = operand
	return l
}

// WithoutRole removes the given role and related content
func (l *localLinkBuilder) WithoutRole(role string) LinkBuilder {
	if l.roles != nil {
		delete(l.roles, role)
	}
	return l
}

// WithActivity changes the activity period for the link
func (l *localLinkBuilder) WithActivity(period periods.Period) LinkBuilder {
	l.activity = period
	return l
}

// WithName changes the name of the link
func (l *localLinkBuilder) WithName(name string) LinkBuilder {
	l.name = name
	return l
}

// Errors returns, if any, current errors so far.
// Errors are cumulative
func (l *localLinkBuilder) Errors() error {
	return l.globalErrors
}

// Build the link or raise an error.
// Once used, it resets all values except the id.
func (l *localLinkBuilder) Build() (Link, error) {
	rolesCopy := make(map[string]Linkable, len(l.roles))
	for role, linkable := range l.roles {
		rolesCopy[role] = linkable
	}

	result := &localLink{
		id:       l.id,
		name:     l.name,
		roles:    rolesCopy,
		activity: l.activity,
	}

	resultErr := l.globalErrors

	// result the builder
	l.roles = nil
	l.activity = periods.NewEmptyPeriod()
	l.name = ""
	l.globalErrors = nil

	return result, resultErr
}

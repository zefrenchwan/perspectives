package objects

import "github.com/zefrenchwan/perspectives.git/periods"

// Linkable defines any link operand.
// Golang does not allow type constraint over interfaces, but it should be
// EITHER A LINK (link composition)
// OR AN INSTANCE (instance as operand),
// OR A TRAIT (trait as operand)
// OR A VARIABLE (for pattern matching).
// To do so, we use the sealed interface :
// we include a private function to force that no other type can implement it.
// This way, linkable types can only be implemented within this package.
type Linkable interface {
	// isLinkable is a private function to force that no other type can implement it.
	// It is used as an implementation of the SEALED INTERFACE go pattern.
	isLinkable() bool
}

// Link relates elements together during a given period.
// For instance, FriendOf(subject=Marie,object=Paul) since now() - 3 years is a link.
type Link interface {
	Linkable // Linkable to use a link as a link operand (compositions)
	Element  // Element to use links as base components of the system
	// Name of the link, it defines its semantic
	Name() string
	// Roles associated to the link, to define how the elements are related together
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
type LinkBuilder interface {
	// Add an operand with the given role
	Add(role string, operand Linkable) LinkBuilder
	// Remove the given role and related content
	Remove(role string) LinkBuilder
	// SetActivity changes the activity period for the link
	SetActivity(period periods.Period) LinkBuilder
	// SetName changes the name of the link
	SetName(name string) LinkBuilder
	// Build the link or raise an error
	Build() (Link, error)
}

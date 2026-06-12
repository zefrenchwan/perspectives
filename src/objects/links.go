package objects

import "github.com/zefrenchwan/perspectives.git/periods"

// Link relates elements together during a given period.
// For instance, FriendOf(subject=Marie,object=Paul) since now() - 3 years is a link.
type Link interface {
	Element // Element to use a link as a link operand (compositions)
	// Name of the link, it defines its semantic
	Name() string
	// Roles associated to the link, to define how the elements are related together
	Roles() []string
	// Role returns the element associated to the given role, if any
	Role(string) (Element, bool)
	// Activity returns the period during which the link is active
	Activity() periods.Period
	// Range iterates over the roles and their associated elements
	Range(func(string, Element) bool)
}

// Relation is a view of a link.
// It is used to observe changes in the link.
type Relation interface {
	Observable[Link] // Used to observe changes in the link
}

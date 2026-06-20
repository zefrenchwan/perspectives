package objects

import (
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Definition is the concrete meaning of a term at a given time.
// Term = word to describe a concept.
// Definition = concrete meaning of a term at a given time.
type Definition interface {
	// Element to link predicates to instances via links
	Element
	// Matches returns the period during which the definition applies to that instance (if any).
	Matches(Instance) (periods.Period, bool)
}

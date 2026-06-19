package objects

import (
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Predicate is an interface for unary predicates.
// Other predicates are indeed links.
// Predicates are concrete ways to match a definition.
type Predicate interface {
	// Element to link predicates to instances via links
	Element
	// Matches returns the period during which the predicate matches the instance (if any).
	Matches(Instance) (periods.Period, bool)
}

package objects

import (
	"github.com/zefrenchwan/perspectives.git/periods"
)

type Entity interface {
	IdentifiableElement
	Activity() periods.Period
	Attributes() []string
	Values(string) map[string]TemporalValues
}

// =========================================================================
// ENTITY IMPLEMENTATION
// =========================================================================

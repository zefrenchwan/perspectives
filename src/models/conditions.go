package models

import "slices"

// Condition tests if parameters match a given test.
// For instance, if an entity as a given type.
//
// A condition may not only be boolean.
// For instance:
// temporal condition may return a period
// probabilistic condition may return a number between 0 and 1
// fuzzy condition may return a map of sets and values
type Condition interface {
	// Triggers returns true if the condition parameters matches the condition
	Triggers(Parameters) bool
}

// ConditionOnEntityType accepts a non nil entity which type is in a set
type ConditionOnEntityType struct {
	// MatchingTypes define the types of accepted entities
	MatchingTypes []EntityType
}

// Triggers accepts parameters if there is one value and its matches an entity definition for that type
func (et ConditionOnEntityType) Triggers(p Parameters) bool {
	if len(et.MatchingTypes) == 0 {
		return false
	} else if p == nil {
		return false
	} else if p.Size() != 1 {
		return false
	} else if len(p.Variables()) != 0 {
		return false
	} else {
		value := p.Get(0)
		if value == nil {
			return false
		} else if entity, ok := value.(ModelEntity); !ok {
			return false
		} else {
			return slices.Contains(et.MatchingTypes, entity.GetType())
		}
	}
}

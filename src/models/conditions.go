package models

// Condition is the most abstract defintion of a condition to match
type Condition interface {
	// Matches returns true if a condition accepts parameters
	Matches(Parameters) bool
}

// IdBasedCondition is a condition to match a given id.
// It matches if parameters has one unique identifiable and ids match between identifiable and Id.
// We use this struct with something in mind.
// It makes no sense to perform a full scan to match an id.
// A clever implementation would use a massive index and then find the matching element with a direct access.
type IdBasedCondition struct {
	// Id to match
	Id string
}

// Matches returns true for an unique identifiable object with that id, false otherwise
func (i IdBasedCondition) Matches(p Parameters) bool {
	if p == nil {
		return false
	} else if value, matches := p.Unique(); !matches {
		return false
	} else if value == nil {
		return false
	} else if identifiable, ok := value.(IdentifiableElement); !ok {
		return false
	} else if identifiable.Id() == i.Id {
		return true
	}

	return false
}

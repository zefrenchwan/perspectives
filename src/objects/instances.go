package objects

import "time"

// Instance is an element that may change over time, but we may identify as unique even for a state change.
// We use a BItemporal model, so it means that there are two time dimensions:
// the time of the instance itself (what changed during its lifetime) and the time of the observation.
type Instance interface {
	// IdentifiableElement to provide a unique identifier for the instance.
	IdentifiableElement
	// Observe returns the state of the instance at the given time.
	// State itself contains the full content over time as far as we know at that obervation time. 
	Observe(time.Time) TimeDependentContent
}

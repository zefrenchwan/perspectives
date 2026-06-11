package objects

// Instance is an element that may change over time, but we may identify as unique even for a state change.
// We use a BItemporal model, so it means that there are two time dimensions:
// the time of the instance itself (what changed during its lifetime) and the time of the observation.
type Instance interface {
	// Observable to provide
	// a unique identifier for the instance (distinguishing it from other instances)
	// an observation as a TimeDependentContent (historical data)
	Observable[TimeDependentContent]
}

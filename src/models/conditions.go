package models

// Condition is the most abstract defintion of a condition to match
type Condition interface {
	// Matches returns true if a condition accepts parameters
	Matches(Parameters) bool
}

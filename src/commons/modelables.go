package commons

// ModelableType defines the type of element we use.
// It does not deal with inheritance, though.
// It is an int, convention is:
// do not use 0,
// use negative values for tests,
// use postive values otherwise
type ModelableType int

// ModelableTypesEquals defines equals for modelable types
func ModelableTypesEquals(a, b ModelableType) bool {
	return a == b
}

// TypeUnmanaged is there just in case
const TypeUnmanaged ModelableType = 1

// TypeStructure defines structures
const TypeStructure ModelableType = 2

// TypeConstraint defines constraints
const TypeConstraint ModelableType = 3

// TypeObject defines objects
const TypeObject ModelableType = 4

// TypeLink defines links
const TypeLink ModelableType = 5

// Modelable is anything that appears in a model
type Modelable interface {
	// GetType returns the most precise possible type of a modelable
	GetType() ModelableType
}

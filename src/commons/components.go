package commons

// ModelComponent is the most general component within a model.
// There are three types of components:
// * objects : what we observe
// * structures: what contains the objects
// * constraints: what objects can and cannot do (and how it changes their structure)
type ModelComponent any

// ModelObject is an object in the model.
// An object operates in a structure and faces constraints
type ModelObject interface {
	// An object is a component of a model
	ModelComponent
}

// ModelStructure defines a structure (where or when objects live in)
type ModelStructure interface {
	// A structure is a component  of a model
	ModelComponent
}

// ModelConstraint defines a constraint (what components can and cannot do)
type ModelConstraint interface {
	// a constraint is a component of a model
	ModelComponent
}

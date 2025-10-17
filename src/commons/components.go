package commons

// ModelComponent is the most general component within a model.
// There are three types of components:
// * entities : what we observe
// * structures: what contains the objects
// * constraints: what objects can and cannot do (and how it changes their structure)
type ModelComponent interface {
	// A component of a model is modelable
	Modelable
}

// ModelConstraint defines a constraint (what components can and cannot do)
type ModelConstraint interface {
	// a constraint is a component of a model
	ModelComponent
}

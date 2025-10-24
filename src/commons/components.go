package commons

// ModelComponent is the most general component within a model.
// There are three types of components:
// * objects : what we observe
// * structures: what contains the objects
// * constraints: what objects can and cannot do (and how it changes their structure)
type ModelComponent interface {
	// We may distinguish a component from the other
	Identifiable
	// A component of a model is modelable
	Modelable
}

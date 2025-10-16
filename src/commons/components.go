package commons

// ModelComponentType describes the type of the component
type ModelComponentType uint8

// ModelUnmanagedType means that this is not something defined in the model
const ModelUnmanagedType ModelComponentType = 0x1

// ModelObjectType is for objects only
const ModelObjectType ModelComponentType = 0x10

// ModelConstraintType is for constraints only
const ModelConstraintType ModelComponentType = 0x11

// ModelStructureType is for structures only
const ModelStructureType ModelComponentType = 0x12

// ModelComponent is the most general component within a model.
// There are three types of components:
// * objects : what we observe
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

// ModelType returns the current type of c within that model.
// Getting the type as a function attached to model component
// such as ( type ModelComponent interface { GetType() int })
// Would make an inheritance issue.
// This method is a type predicate.
// This way, we may complete it with other type management.
func ModelType(c any) ModelComponentType {
	if c == nil {
		return ModelUnmanagedType
	} else if _, ok := c.(ModelConstraint); ok {
		return ModelConstraintType
	} else if _, ok := c.(ModelObject); ok {
		return ModelObjectType
	} else if _, ok := c.(ModelStructure); ok {
		return ModelStructureType
	} else {
		return ModelUnmanagedType
	}
}

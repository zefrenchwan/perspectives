package commons

// ModelStructure defines a structure (where or when objects live in)
type ModelStructure interface {
	// A structure is a component  of a model
	ModelComponent
	// Register adds a non existing object in the structure.
	// It returns true if object was added (not existed before) or raises an error.
	// It means it is not an upsert, just an insert
	Register(ModelObject) (bool, error)
}

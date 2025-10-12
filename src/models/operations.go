package models

// Operation is the most general definition of a way to change a state
type Operation interface {
}

// Operations groups operations.
// There is no specific order to execute each operation
type Operations struct {
	// Content are the operations to perform
	Content []Operation
}

// OperationToUpsertObjects uperts objets
type OperationToUpsertObjects struct {
	// Content is the set of objects to upsert
	Content []Object
}

// OperationToUpsertLinks upserts links.
// Leafs of links may not be concerned by the change.
// Note that non existing leafs will be created no matter the value of UpdateLeafs.
// UpdateLeafs is just about UPDATING leafs.
type OperationToUpsertLinks struct {
	// Content defines the links to create or update
	Content []Link
	// UpdateLeafs true means updating leafs too
	UpdateLeafs bool
}

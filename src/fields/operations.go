package fields

import "github.com/zefrenchwan/perspectives.git/models"

// OperationToUpsertObjects uperts objets
type OperationToUpsertObjects struct {
	// Content is the set of objects to upsert
	Content []models.Object
}

// OperationToUpsertLinks upserts links.
// Leafs of links may not be concerned by the change.
// Note that non existing leafs will be created no matter the value of UpdateLeafs.
// UpdateLeafs is just about UPDATING leafs.
type OperationToUpsertLinks struct {
	// Content defines the links to create or update
	Content []models.Link
	// UpdateLeafs true means updating leafs too
	UpdateLeafs bool
}

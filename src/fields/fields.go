package fields

import "github.com/zefrenchwan/perspectives.git/commons"

// Field is the definition of the structural part of a model.
// "Field" defines a field for entities
type Field interface {
	// A field has an id
	commons.IdentifiableElement
}

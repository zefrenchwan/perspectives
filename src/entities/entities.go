package entities

import "github.com/zefrenchwan/perspectives.git/commons"

// Entity is an object within the system that exists for sure, and that we may identify.
type Entity interface {
	// Identifiable to get an immutable id for that entity.
	commons.Identifiable
}

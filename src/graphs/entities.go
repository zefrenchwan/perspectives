package graphs

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Entity is a graph element.
// It exists for sure but its state might change over time.
type Entity interface {
	// Identifiable provides a unique identifier for the entity.
	commons.Identifiable
	// CreationDate returns the moment the entity was created.
	CreationDate() time.Time
	// AsOf returns the state of the entity at the given time.
	AsOf(time time.Time) (State, bool)
}

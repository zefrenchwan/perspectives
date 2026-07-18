package graphs

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Entity is an immutable graph element.
type Entity interface {
	// Identifiable provides a unique identifier for the entity.
	commons.Identifiable
	// CreationDate returns the moment the entity was created.
	CreationDate() time.Time
	// Asof returns the state of the entity at the given time.
	Asof(time time.Time) (State, bool)
}

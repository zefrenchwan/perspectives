package graphs

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// Entity is an immutable graph element.
type Entity interface {
	// Identifiable provides a unique identifier for the entity.
	commons.Identifiable
	// TimeBounded to define a time period during which the entity exists.
	periods.TimeBounded

	Asof(time time.Time) (State, bool)
}

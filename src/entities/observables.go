package entities

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Observable represents an element which state may depend over time, but keeps having the same identity.
type Observable interface {
	commons.Identifiable // Identifiable to represent an element with a unique time-independent identifier.
	At(time.Time) Entity // At a given moment in time, returns the state of the observable.
}

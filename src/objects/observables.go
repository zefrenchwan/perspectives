package objects

import "time"

type Observable interface {
	// Id returns a unique identifier for the observable.
	Id() string
	// Observe returns the entity at a given time.
	Observe(at time.Time) Entity
	// Change updates the entity at a given time.
	Change(moment time.Time, newValue Entity) string
}

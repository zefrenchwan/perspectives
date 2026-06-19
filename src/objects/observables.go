package objects

import "time"

// Observable manages the bitemporal model by defining a content fixed at a given time.
// That time is the moment we analyze the environment, not the moment the environment changed.
type Observable[T Element] interface {
	// Id returns a unique identifier for the observable.
	Id() string
	// Observe returns the content of the observable at a given time.
	Observe(at time.Time) T
	// Change updates the content of the observable at a given time.
	Change(moment time.Time, newValue T) string
}

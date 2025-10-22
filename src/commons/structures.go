package commons

import (
	"iter"
	"time"
)

// ModelStructure defines a structure (where or when objects live in)
type ModelStructure interface {
	// A structure is a component  of a model
	ModelComponent
}

// SpreadingStructure defines a graph structure with a time dependent topology.
// Time is the same for any object, events and the structure itself
type SpreadingStructure interface {
	// Register adds a new processor and returns if structure accepted that object.
	// If result is true, structure contains that object since creationTime
	Register(EventProcessor, creationTime time.Time) bool
	// Neighbors returns, for a given event and time, the neighbors of the object.
	// It depends on the object, the event (especially the kind of events),
	// and the moment the object asks for.
	Neighbors(EventProcessor, Event, time.Time) iter.Seq[EventProcessor]
	// Notify notifies ALL objects to process those events.
	// For instance, an EventTick should apply to all objects in that structure
	Notify([]Event)
}

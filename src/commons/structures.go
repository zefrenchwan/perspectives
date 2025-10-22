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

// DecentralizedStructure defines a graph structure with a time dependent topology.
// Time is the same for any object, events and the structure itself
type DecentralizedStructure interface {
	// DecentralizedStructure implements event processor.
	// Usually, accepted events are:
	// content creation such as adding objects in the structure,
	// ticks to run the next step
	EventProcessor
	// Neighbors returns, for a given event and time, the neighbors of the object.
	// It depends on the object, the event (especially the kind of events),
	// and the moment the object asks for.
	Neighbors(EventProcessor, Event, time.Time) iter.Seq[EventProcessor]
}

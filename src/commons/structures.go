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

type SpreadingStructure interface {
	EventProcessor
	Neighbors(moment time.Time, source EventProcessor) iter.Seq[EventProcessor]
}

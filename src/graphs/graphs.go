package graphs

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/entities"
	"github.com/zefrenchwan/perspectives.git/events"
)

// DynamicGraph represents a graph that may change over time.
// Changes to the graph are applied through events.
type DynamicGraph interface {
	// Apply applies an event to the graph, that may change its content.
	Apply(event events.Event)
	// At returns the snapshot of the graph at the given time.
	At(time.Time) GraphSnapshot
	// Content returns the content as an iterator of observables.
	Content() iter.Seq[entities.Observable]
}

package graphs

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/entities"
)

// GraphSnapshot represents a snapshot of a graph at a given time.
// Content is then all the entities and links in the graph.
type GraphSnapshot interface {
	// SnapshotTime returns the time at which the snapshot was taken.
	SnapshotTime() time.Time
	// Entities returns an iterator over entities in the graph, fixed at that time.
	Entities() iter.Seq[entities.Entity]
	// Links returns an iterator over links in the graph from that entity, fixed at that time.
	Links(entities.Entity) iter.Seq[entities.DynamicLink]
}

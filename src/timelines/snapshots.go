package timelines

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Snapshot represents a frozen record of the timeline's state at a specific point in time.
type Snapshot[T any] interface {
	// Identifiable provides Id() method for identifying the snapshot.
	commons.Identifiable
	// RecordDate returns the moment when the snapshot was taken.
	RecordDate() time.Time
	// Content returns the content of the snapshot.
	Content() T
}

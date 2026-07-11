package timelines

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Timeline defines what we know about a system during a long period of time.
// An analogy would be a sort of diary : at that date, we know this snapshot about the system.
type Timeline[T any] interface {
	// Identifiable to load or use the current timeline
	commons.Identifiable
	// AsOf returns what we knew about the system at the given date, as a snapshot.
	// When no data exists, just return false.
	AsOf(time.Time) (Snapshot[T], bool)
	// Snapshots returns all the snapshots of the timeline, in chronological order
	Snapshots() iter.Seq2[time.Time, Snapshot[T]]
}

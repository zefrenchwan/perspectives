package timelines

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// Layer is a node with a parent forming a hierarchy.
// When combined with other layers, it forms a full information snapshot.
type Layer[T any] interface {
	// Identifiable provides Id() method for identifying the layer.
	commons.Identifiable
	// Parent returns the parent layer, assuming to be unique.
	Parent() Layer[T]
	// Import loads the content of another layer into this layer.
	Import(layer Layer[T])
	// Flatten returns the flattened content of the full hierarchy of layers.
	Flatten() T
}

// SnapshotLayer builds a snapshot as a combination of layers.
type SnapshotLayer[T any] struct {
	// id of the snapshot, different from the layer.
	id string
	// recordDate is the time to consider as the record date of the snapshot.
	recordDate time.Time
	// content is the layer that contains, with its hierarchy, all the information.
	content Layer[T]
}

// NewSnapshotLayer builds a snapshot as a combination of layers.
// Parameters are:
// id of the snapshot (an uuid is enough)
// asof is the time to consider as the record date of the snapshot.
// content is the layer that contains, with its hierarchy, all the information.
func NewSnapshotLayer[T any](id string, asof time.Time, content Layer[T]) *SnapshotLayer[T] {
	return &SnapshotLayer[T]{
		id:         id,
		recordDate: asof,
		content:    content,
	}
}

// Id returns the id of the snapshot.
// It should NOT be the id of the layer.
func (l *SnapshotLayer[T]) Id() string {
	return l.id
}

// RecordDate returns the record date of the snapshot.
func (l *SnapshotLayer[T]) RecordDate() time.Time {
	return l.recordDate
}

// Content returns the content of the snapshot.
func (l *SnapshotLayer[T]) Content() T {
	return l.content.Flatten()
}

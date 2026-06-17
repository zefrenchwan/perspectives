package objects

// Relation represents a link between two entities that changes over time.
// At a given observation time, the relation is a link between entities.
// The relation defines the core essence of what links have in common.
type Relation interface {
	// Observable to read the content of the relation at a given time (a link) or change it.
	Observable[Link]
}

package objects

// Entity represents a thing that can be observed.
// When observed, its concrete representation is an instance with a given history.
type Entity interface {
	// Observable means that, at a given observation time, the entity definition is an instance.
	Observable[Instance]
}

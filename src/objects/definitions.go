package objects

// Definition defines concepts with a name (its id).
// Predicates are the concrete way to test at a given time (because words change their meaning over time).
type Definition interface {
	// Observable of a definition is the actual way to test it, hence, a predicate.
	Observable[Predicate]
}

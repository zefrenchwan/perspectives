package objects

// Term defines a general notion, something we know about, no matter its formal definition at a given time.
// Definitions are the concrete implementation of a term at a given time (because words change their meaning over time).
type Term interface {
	// Observable of an abstraction is the actual way to test it, hence, a definition.
	Observable[Definition]
}

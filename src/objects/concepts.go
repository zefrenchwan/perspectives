package objects

// Concept represents a definition of something (a trait) that depends on time.
// Concept is the general idea, whereas concrete criterias may change.
type Concept interface {
	// Observable means that, at a given observation time, the concept definition is a trait.
	Observable[Trait]
}

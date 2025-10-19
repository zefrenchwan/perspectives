package commons

// ModelConstraint defines a constraint (what components can and cannot do)
type ModelConstraint interface {
	// a constraint is a component of a model
	ModelComponent
	// a constraint is identifiable for sure
	Identifiable
	// AcceptanceCriteria returns the condition that accepts elements
	AcceptanceCriteria() Condition
}

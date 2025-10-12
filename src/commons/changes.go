package commons

// Change is the most general definition of a state variation
type Change interface {
	Action() Operation
}

// TriggeredChange is a conditional change
type TriggeredChange struct {
	// Trigger defines the condition to apply the operation
	Trigger Condition
	// Response is the operation to apply if condition matches
	Response Operation
}

// Action returns the operation to perform if trigger matches
func (t TriggeredChange) Action() Operation {
	return t.Response
}

package commons

// Change is the most general definition of a state variation
type Change interface {
	// Operation returns the action to execute
	Operation() Action
}

// TriggeredChange is a conditional change
type TriggeredChange struct {
	// Trigger defines the condition to apply the action
	Trigger Condition
	// Response is the action to apply if condition matches
	Response Action
}

// Operation returns the action to perform if trigger matches
func (t TriggeredChange) Operation() Action {
	return t.Response
}

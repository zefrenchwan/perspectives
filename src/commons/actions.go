package commons

// Action is the most general definition of a way to change a state
type Action interface {
	// Execute an action over parameters
	Execute(Parameters) error
}

// Actions groups actions.
// There is no specific order to execute each action
type Actions struct {
	// Content are the actions to perform
	Content []Action
}

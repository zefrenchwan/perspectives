package commons

// Action is the most general definition of a way to change a state
type Action interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
	// Execute an action over content
	Execute(Content) error
}

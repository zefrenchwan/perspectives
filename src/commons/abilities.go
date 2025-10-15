package commons

// Ability is the general defintion of any call to a component.
// It means:
// Actions to execute (no return value)
// Descriptions to get informations
type Ability interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
}

// Action is the most general definition of a way to change a state
type Action interface {
	// An action is an ability to do something
	Ability
	// Execute an action over content
	Execute(Content) error
}

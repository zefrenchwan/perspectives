package commons

// Ability is the general defintion of any call to a component.
// It means:
// Actions to execute (no return value)
// Descriptions to get informations
type Ability interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
}

package commons

// Ability is the general defintion of any call to a component.
// It means:
// Actions to execute (no return value)
// Descriptions to get informations
type Ability interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
}

// Action executes code on a content
type ExecuteAction interface {
	// Ability is necessary to get signature
	Ability
	// Execute runs the action
	Execute(c Content) error
}

// RequestDescription is the ability to describe a state
type RequestDescription[T StateValue] interface {
	// Ability is necessary to get signature
	Ability
	// Describe gets the description from the content
	Describe(c Content) StateDescription[T]
}

// RequestTemporalDescription is the ability to describe a state that changes over time
type RequestTemporalDescription[T StateValue] interface {
	// Ability is necessary to get signature
	Ability
	// Describe gets the description from the content
	Describe(c Content) TemporalStateDescription[T]
}

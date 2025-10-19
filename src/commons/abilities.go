package commons

import "errors"

// Ability is the general defintion of any call to a component.
// It means:
// Actions to execute (no return value)
// Descriptions to get informations
type Ability interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
}

// Action executes code on a content.
// Convention is usually to raise an error if the action failed, not if parameters were missing.
// For instance, if action is x.Price = 0 then it is fine if content did not contain x.
// Why ? Because signature checks this part.
type ExecuteAction interface {
	// Ability is necessary to get signature
	Ability
	// Execute runs the action
	Execute(c Content) error
}

// SequentialActions runs actions one after another
type SequentialActions struct {
	// ExpectedParameters is the global signature
	ExpectedParameters FormalParameters
	// Actions are the actions to execute
	Actions []ExecuteAction
	// StopAtFirstError is true if we end execution at first error, false for a full run
	StopAtFirstError bool
}

// Signature returns the signature of the actions
func (s SequentialActions) Signature() FormalParameters {
	return s.ExpectedParameters
}

// Execute runs each action, and returns errors if any.
// If StopAtFirstError, run until an error and return first.
// Else, keep running and join errors during the execution
func (s SequentialActions) Execute(c Content) error {
	var currentError error
	for _, action := range s.Actions {
		if err := action.Execute(c); err != nil {
			currentError = errors.Join(currentError, err)
			if s.StopAtFirstError {
				return err
			}
		}
	}

	return currentError
}

// NewEmptyAction returns an action that accepts everything and do nothing
func NewEmptyAction() ExecuteAction {
	return SequentialActions{}
}

// NewSequentialActions returns an action expecting the max of signatures and running actions in same order.
// When stopAtError, then actions chain stops at first error
func NewSequentialActions(actions []ExecuteAction, stopAtError bool) ExecuteAction {
	if len(actions) == 0 {
		return NewEmptyAction()
	}

	var signature FormalParameters
	var order []ExecuteAction
	for index, action := range actions {
		if index == 0 {
			signature = action.Signature()
		} else {
			signature = signature.Max(action.Signature())
		}

		order = append(order, action)
	}

	return SequentialActions{
		ExpectedParameters: signature,
		Actions:            order,
		StopAtFirstError:   stopAtError,
	}
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

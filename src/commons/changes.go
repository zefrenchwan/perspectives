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

// Run triggers the condition and, if it matches, execute the action on p
func (t TriggeredChange) Run(p Parameters) error {
	// no need to launch code if condition or action is nil.
	// Otherwise,
	// check trigger
	// launch action accordingly
	if condition := t.Trigger; condition == nil {
		return nil
	} else if response := t.Response; response == nil {
		return nil
	} else if triggered, err := condition.Matches(p); err != nil {
		return err
	} else if !triggered {
		return nil
	} else {
		return response.Execute(p)
	}
}

package commons

// ModelConstraint defines a constraint (what components can and cannot do)
type ModelConstraint interface {
	// a constraint is a component of a model
	ModelComponent
}

// applyConstraintLifetimeEvent returns true if event was applied to object
func applyConstraintLifetimeEvent(event Event, object ModelObject) bool {
	if event == nil || object == nil {
		return false
	}

	if end, ok := event.(EventLifetimeEnd); ok {
		endTime := end.ProcessingTime()
		if thandler, okHandler := object.(TemporalHandler); okHandler && thandler != nil {
			current := thandler.ActivePeriod()
			remaining := current.Remove(NewPeriodSince(endTime, true))
			thandler.SetActivePeriod(remaining)
			return true
		}
	}

	return false
}

// applyConstraintStateEventForType returns true if event was a state change and object was relevant
func applyConstraintStateEventForType[T StateValue](event Event, object ModelObject) bool {
	if change, ok := event.(EventStateChanges[T]); ok {
		if stateHandler, okHandler := object.(StateHandler[T]); okHandler {
			if change != nil && len(change.Changes()) != 0 {
				stateHandler.SetValues(change.Changes())
			}

			return true
		} else if stateHandler, okHandler := object.(TemporalStateHandler[T]); okHandler {
			if change != nil && len(change.Changes()) != 0 {
				period := NewPeriodSince(change.ProcessingTime(), true)
				for attr, value := range change.Changes() {
					stateHandler.SetValueDuringPeriod(attr, value, period)
				}
			}
			return true
		}
	}

	return false
}

// applyConstraintStateEvent applies state events for all types within StateValue.
// Consider it as a technical implementation but it ensures we do not miss a change.
func applyConstraintStateEvent(event Event, object ModelObject) bool {
	counter := 0
	if applyConstraintStateEventForType[int](event, object) {
		counter++
	}

	if applyConstraintStateEventForType[bool](event, object) {
		counter++
	}

	if applyConstraintStateEventForType[string](event, object) {
		counter++
	}

	if applyConstraintStateEventForType[float64](event, object) {
		counter++
	}

	return counter != 0
}

// ApplyStateActivityConstraintsOnEvent applies event on object to change its state or lifetime if possible.
// It is a constraint: an object has to follow some events without being able to negociate them.
// Result is true if event should be processed by object too, false otherwise.
// For instance, a ending lifetime shuts down the object no matter what, so no need to transmit.
// This is default implementation.
func ApplyStateActivityConstraintsOnEvent(event Event, object ModelObject) bool {
	if event == nil {
		return false
	} else if object == nil {
		return false
	}

	counter := 0
	if applyConstraintStateEvent(event, object) {
		counter++
	}

	if applyConstraintLifetimeEvent(event, object) {
		counter++
	}

	// if counter > 0, event was applied and does not need object processing
	return counter == 0
}

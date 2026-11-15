package events

import (
	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// Change is any state change within the system
type Change interface {
	// Identifiable to know what change was applied
	commons.Identifiable
	// TargetId is the id of the element that changed
	TargetId() string
}

// baseChange contains the main functions for a base change
type baseChange struct {
	// changeId is the id of the change
	changeId string
	// targetId is the id of the target object
	targetId string
}

// Id is the id of the change
func (c baseChange) Id() string {
	return c.changeId
}

// TargetId is the id of the element to change
func (c baseChange) TargetId() string {
	return c.targetId
}

// ActivityChange defines a change for an activity.
// Use it for instance to create an element
type ActivityChange struct {
	// baseChange to reuse the core changes
	baseChange
	// activity is the new activity to set
	activity periods.Period
}

// Activity returns the new activity to set for that element
func (a ActivityChange) Activity() periods.Period {
	return a.activity
}

// CreateActivityChange returns a new activity change for the element with that reference.
// target id is the id of the element to change
// activity is the activity to set for that element.
// It sets the activity, no matter the previous value.
func CreateActivityChange(targetId string, activity periods.Period) ActivityChange {
	return ActivityChange{
		baseChange: baseChange{
			changeId: commons.NewId(),
			targetId: targetId,
		},
		activity: activity,
	}
}

// PropertyChange changes the value of that property for that target.
// For instance, set name of element ... to be "John", valid since 2020
type PropertyChange struct {
	// baseChange to get id and target
	baseChange
	// propertyName is the name of the property to change
	propertyName string
	// propertyValue is the value to set
	propertyValue values.PrimitiveValue
	// validity is the validity of the change
	validity periods.Period
}

// PropertyName returns the name of the property to set
func (p PropertyChange) PropertyName() string {
	return p.propertyName
}

// PropertyValue returns the value to set
func (p PropertyChange) PropertyValue() values.PrimitiveValue {
	return p.propertyValue
}

// Validity returns the validity of the mapping name => value
func (p PropertyChange) Validity() periods.Period {
	return p.validity
}

// CreatePropertyChange creates a new change to set a property value (by name), valid during a period.
func CreatePropertyChange(targetId string, propertyName string, propertyValue values.PrimitiveValue, validity periods.Period) PropertyChange {
	return PropertyChange{
		baseChange: baseChange{
			changeId: commons.NewId(),
			targetId: targetId,
		},
		propertyName:  propertyName,
		propertyValue: propertyValue,
		validity:      validity,
	}
}

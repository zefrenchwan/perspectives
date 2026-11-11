package events

import (
	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
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

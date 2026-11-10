package events

import "github.com/zefrenchwan/perspectives.git/commons"

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

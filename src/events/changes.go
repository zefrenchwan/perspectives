package events

import "time"

// Change is any state change within the system
type Change interface {
	// ActualDate is the moment the change applies
	ActualDate() time.Time
	// TargetId is the id of the object that changed
	TargetId() string
}

// baseChange is the most basic change to regroup all base code
type baseChange struct {
	// targetId is the id of the element to change
	targetId string
	// actualDate is the moment the element changed
	actualDate time.Time
	// actualDateIncluded
	actualDateIncluded bool
}

// ActualDate returns the time the event applies
func (b baseChange) ActualDate() time.Time {
	return b.actualDate
}

// ActualDateIncluded returns true if actual date is included or not
func (b baseChange) ActualDateIncluded() bool {
	return b.actualDateIncluded
}

// TargetId returns the element of the system to apply change to
func (b baseChange) TargetId() string {
	return b.targetId
}

// ActivityChange means a change for an activity
type ActivityChange struct {
	baseChange
	active bool
}

func (ac ActivityChange) Active() bool {
	return ac.active
}

func CreateActivityChange(targetId string, active bool, actualDate time.Time, included bool) ActivityChange {
	return ActivityChange{
		active: active,
		baseChange: baseChange{
			targetId:           targetId,
			actualDate:         actualDate,
			actualDateIncluded: included,
		},
	}
}

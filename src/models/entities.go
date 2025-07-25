package models

import "time"

// Entity is base container for any element
type Entity struct {
	Id       string // id of the entity
	lifetime Period // period during the entity is still "alive"
}

// NewEntity builds an entity starting at a given time (included)
func NewEntity(id string, creationTime time.Time) Entity {
	return Entity{Id: id, lifetime: NewPeriodSince(creationTime, true)}
}

// End closes the entity lifetime.
// After that, remaining is [startTime, endTime[
func (e *Entity) End(endTime time.Time) {
	if e != nil {
		remaining := e.lifetime.Remove(NewPeriodSince(endTime, true))
		e.lifetime = remaining
	}
}

// LifetimeDuringPeriod returns the intersection of current lifetime with reference period
func (e *Entity) LifetimeDuringPeriod(reference Period) Period {
	return reference.Intersection(e.lifetime)
}

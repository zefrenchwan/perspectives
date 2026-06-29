package events

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

type Event interface {
	TransactionTime() time.Time
	isEvent() bool
}

type EntityCreation struct {
	EntityId       string
	ActionTime     time.Time
	CreationTime   time.Time
	OriginalValues map[string]any
}

func (e EntityCreation) TransactionTime() time.Time {
	return e.ActionTime
}

func (e EntityCreation) isEvent() bool {
	return true
}

type EntitySetValue struct {
	EntityId   string
	Attribute  string
	Value      any
	ValueType  string
	Period     periods.Period
	ActionTime time.Time
}

func (e EntitySetValue) TransactionTime() time.Time {
	return e.ActionTime
}

func (e EntitySetValue) isEvent() bool {
	return true
}

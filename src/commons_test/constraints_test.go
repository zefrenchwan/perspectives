package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestConstraintOnStateObject(t *testing.T) {
	obj := commons.NewStateObject[int]()

	// test event change
	obj.SetValue("age", 10)
	event := commons.NewEventStateChanges(time.Now(), map[string]int{"age": 100})
	if propagate := commons.OnEventApplyConstraintsToObject(event, obj); propagate {
		t.Fail()
	} else if v, found := obj.GetValue("age"); !found || v != 100 {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Fail()
	}

	// test period change
	now := time.Now().Truncate(time.Second)
	end := commons.NewEventLifetimeEnd(now)
	if propagate := commons.OnEventApplyConstraintsToObject(end, obj); propagate {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewPeriodUntil(now, false)) {
		t.Fail()
	}
}

func TestConstraintOnTemporalStateObject(t *testing.T) {
	obj := commons.NewTemporalStateObject[int](commons.NewFullPeriod())

	// test event change
	obj.SetValueDuringPeriod("age", 10, commons.NewFullPeriod())
	now := time.Now().Truncate(time.Second)
	event := commons.NewEventStateChanges(now, map[string]int{"age": 100})
	if propagate := commons.OnEventApplyConstraintsToObject(event, obj); propagate {
		t.Fail()
	} else if values, found := obj.GetValue("age", true); !found || len(values) != 2 {
		t.Fail()
	} else if !values[10].Equals(commons.NewPeriodUntil(now, false)) {
		t.Fail()
	} else if !values[100].Equals(commons.NewPeriodSince(now, true)) {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Fail()
	}

	// test period change
	obj = commons.NewTemporalStateObject[int](commons.NewFullPeriod())
	end := commons.NewEventLifetimeEnd(now)
	if propagate := commons.OnEventApplyConstraintsToObject(end, obj); propagate {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewPeriodUntil(now, false)) {
		t.Fail()
	}
}

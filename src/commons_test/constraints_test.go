package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestConstraintOnStateEvent(t *testing.T) {
	structure := DummyStructure{}
	obj := commons.NewStateObject[int]()

	// test event change
	obj.SetValue("age", 10)
	event := commons.NewEventStateChanges(structure, time.Now(), map[string]int{"age": 100})
	if propagate := commons.ApplyStateActivityConstraintsOnEvent(event, obj); propagate {
		t.Fail()
	} else if v, found := obj.GetValue("age"); !found || v != 100 {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Fail()
	}

	// test period change
	now := time.Now().Truncate(time.Second)
	end := commons.NewEventLifetimeEnd(structure, now)
	if propagate := commons.ApplyStateActivityConstraintsOnEvent(end, obj); propagate {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewPeriodUntil(now, false)) {
		t.Fail()
	}
}

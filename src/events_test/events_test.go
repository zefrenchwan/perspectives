package events_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/events"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestEventCreationAt(t *testing.T) {
	moment := time.Now().AddDate(-1, 0, 0).Truncate(time.Hour)
	changes := []events.Change{
		events.CreateActivityChange("popo", periods.NewFullPeriod()),
	}
	e := events.CreateBaseEventAt(moment, changes)

	if e.Id() == "" {
		t.Error("Id should not be empty")
		t.Fail()
	} else if e.RecordDate() != moment {
		t.Error("Record date should be moment")
		t.Fail()
	} else if foundChanges := slices.Collect(e.Changes()); len(foundChanges) != len(changes) {
		t.Error("Found changes should be the same")
		t.Fail()
	} else if !slices.EqualFunc(changes, foundChanges, func(a, b events.Change) bool { return a.Id() == b.Id() }) {
		t.Error("Found changes should be the same")
		t.Fail()
	}
}

func TestEventCreation(t *testing.T) {
	now := time.Now()
	changes := []events.Change{
		events.CreateActivityChange("popo", periods.NewFullPeriod()),
	}
	e := events.CreateBaseEvent(changes)

	if e.Id() == "" {
		t.Error("Id should not be empty")
		t.Fail()
	} else if e.RecordDate().Sub(now).Abs() >= time.Minute {
		t.Error("base event should be almost now")
		t.Fail()
	} else if foundChanges := slices.Collect(e.Changes()); len(foundChanges) != len(changes) {
		t.Error("Found changes should be the same")
		t.Fail()
	} else if !slices.EqualFunc(changes, foundChanges, func(a, b events.Change) bool { return a.Id() == b.Id() }) {
		t.Error("Found changes should be the same")
		t.Fail()
	}
}

package events_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/events"
)

func TestEvents(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	change := events.CreateActivityChange("test", true, now, false)
	event := events.CreateBaseEvent([]events.Change{change})

	if event.Id() == "" {
		t.Error("Expected an event id")
		t.Fail()
	} else if diff := event.RecordDate().Sub(now); diff >= 10*time.Second {
		t.Errorf("Record date should be now")
		t.Fail()
	} else if values := slices.Collect(event.Changes()); len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
		t.Fail()
	} else if values[0] != change {
		t.Errorf("Expected %v, got %v", change, values[0])
		t.Fail()
	}
}

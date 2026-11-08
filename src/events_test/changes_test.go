package events_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/events"
)

func TestActivityChange(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	activity := events.CreateActivityChange("other id", true, now, false)

	if activity.TargetId() != "other id" {
		t.Error("Activity target id should not be set")
		t.Fail()
	} else if !activity.Active() {
		t.Errorf("CreateActvityChange was incorrect")
		t.Fail()
	} else if activity.ActualDateIncluded() {
		t.Errorf("CreateActvityChange was incorrect")
		t.Fail()
	} else if !activity.ActualDate().Equal(now) {
		t.Errorf("CreateActvityChange was incorrect")
		t.Fail()
	}
}

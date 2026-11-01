package entities_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/entities"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestBuildEventActivity(t *testing.T) {
	now := time.Now()
	epsilon := time.Minute
	event := entities.BuildEventChangeActivity("id", periods.NewFullPeriod())

	if event.Id() == "" {
		t.Error("Id should not be empty")
		t.Fail()
	} else if event.CreationDate().Sub(now) >= epsilon {
		t.Log("creation date failure, should be almost now")
		t.Fail()
	} else if event.EntityId() != "id" {
		t.Log("entity id failure")
		t.Fail()
	} else if !event.Activity().Equals(periods.NewFullPeriod()) {
		t.Log("activity failure, should be full")
		t.Fail()
	}
}

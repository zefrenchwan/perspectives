package events_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/events"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestActivityChange(t *testing.T) {
	ac := events.CreateActivityChange("test", periods.NewFullPeriod())
	if ac.Id() == "" {
		t.Error("id should not be empty")
		t.Fail()
	} else if ac.TargetId() != "test" {
		t.Error("id should be test")
		t.Fail()
	} else if !ac.Activity().Equals(periods.NewFullPeriod()) {
		t.Error("activity should be full")
		t.Fail()
	}
}

func TestPropertyChange(t *testing.T) {
	v := values.NewBool(true)
	pc := events.CreatePropertyChange("test", "prop", v, periods.NewFullPeriod())
	if pc.Id() == "" {
		t.Error("id should not be empty")
		t.Fail()
	} else if pc.TargetId() != "test" {
		t.Error("id should be test")
		t.Fail()
	} else if !pc.PropertyValue().Equals(v) {
		t.Error("property value should be equal")
		t.Fail()
	} else if !pc.Validity().Equals(periods.NewFullPeriod()) {
		t.Error("activity should be full")
		t.Fail()
	}
}

func TestLinkChange(t *testing.T) {
	lc := events.CreateLinkChange("linkId", "subject", "destinationId", periods.NewFullPeriod())
	if lc.Id() == "" {
		t.Error("id should not be empty")
		t.Fail()
	} else if lc.TargetId() != "linkId" {
		t.Error("id should be linkId")
		t.Fail()
	} else if lc.RoleName() != "subject" {
		t.Error("role name should be subject")
		t.Fail()
	} else if lc.LinkDestination().Content() != "destinationId" {
		t.Error("link destination should be destinationId")
		t.Fail()
	} else if !lc.Validity().Equals(periods.NewFullPeriod()) {
		t.Error("activity should be full")
		t.Fail()
	}
}

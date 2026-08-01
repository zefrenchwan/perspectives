package models_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestActivityChangeEvent(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	full := periods.NewFullPeriod()
	event := models.UpdateActivityEvent(
		"id",
		"entity id",
		now,
		full,
	)

	if event.Id() != "id" {
		t.Errorf("Expected event id to be id")
	} else if event.EntityId() != "entity id" {
		t.Errorf("Expected event entity id to be entity id")
	} else if !event.Activity().Equals(full) {
		t.Errorf("Expected event activity to be full period")
	} else if !event.RecordDate().Equal(now) {
		t.Errorf("Expected event record date to be %v", now)
	} else if event.ToHashString() == "" {
		t.Errorf("Expected event hash string to be non empty")
	}
}

func TestUpdateAttributeEvent(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	full := periods.NewFullPeriod()
	event := models.UpdateAttributeEvent("id", "entity id", now, "name", values.NewString("value"), full)

	if event.Id() != "id" {
		t.Errorf("Expected event id to be id")
	} else if event.EntityId() != "entity id" {
		t.Errorf("entity id not set")
	} else if !event.RecordDate().Equal(now) {
		t.Errorf("Expected event record date to be %v", now)
	} else if event.Name() != "name" {
		t.Errorf("Expected event attribute name to be name")
	} else if event.Value().Content() != "value" {
		t.Errorf("Expected event attribute value to be value")
	} else if !event.Validity().Equals(full) {
		t.Errorf("Expected event attribute period to be full period")
	}
}

func TestUpdateRoleEvent(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	full := periods.NewFullPeriod()
	event := models.UpdateRoleEvent("id", "entity id", now,
		"subject", values.NewReference("id of an entity"), full)

	if event.Id() != "id" {
		t.Errorf("Expected event id to be id")
	} else if event.EntityId() != "entity id" {
		t.Errorf("entity id not set")
	} else if !event.RecordDate().Equal(now) {
		t.Errorf("Expected event record date to be %v", now)
	} else if event.Name() != "subject" {
		t.Errorf("Expected event role name to be subject")
	} else if event.Value().Content() != "id of an entity" {
		t.Errorf("Expected event role value to be the one set before")
	} else if !event.Validity().Equals(full) {
		t.Errorf("Expected event role period to be full period")
	}
}

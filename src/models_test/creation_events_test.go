package models_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestCreateEntityEvent(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	full := periods.NewFullPeriod()

	name := values.NewStringLocalMapping(map[string]periods.Period{"Name of the entity": periods.NewFullPeriod()})
	attributes := map[string]values.ImmutableValuesMapping[values.PrimitiveValue]{"name": name}
	state := models.NewLocalState("state id", periods.NewFullPeriod(), attributes, nil)

	event := models.CreateEntityEvent(
		"event1",
		"entity1",
		now,
		periods.NewFullPeriod(),
		state)

	if event.Id() != "event1" {
		t.Errorf("Expected event id to be 'event1', but got '%s'", event.Id())
	} else if event.EntityId() != "entity1" {
		t.Errorf("Expected entity id to be 'entity1', but got '%s'", event.EntityId())
	} else if !event.RecordDate().Equal(now) {
		t.Errorf("Expected created at to be '%s', but got '%s'", now, event.RecordDate())
	} else if !event.InitialPeriod().Equals(full) {
		t.Errorf("Expected initial period to be '%v', but got '%v'", full, event.InitialPeriod())
	} else if event.InitialState().Id() != "state id" {
		t.Errorf("Expected initial state id to be 'state id', but got '%s'", event.InitialState().Id())
	} else if event.ToHashString() == "" {
		t.Errorf("Expected event hash string to be non-empty, but got '%s'", event.ToHashString())
	}
}

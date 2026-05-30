package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestInstanceBuilder(t *testing.T) {
	if instance, err := commons.NewInstanceBuilder().SetId("id").SetValue("name", commons.NewFullPeriod(), "john doe").Build(); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if instance.Id() != "id" {
		t.Errorf("expected id to be 'id', got '%s'", instance.Id())
	} else if value, found := instance.Attribute("name").At(time.Now()); !found || value != "john doe" {
		t.Errorf("expected name value to be 'john doe', got '%s'", value)
	}

	if instance, err := commons.NewInstanceBuilder().
		SetId("id").
		SetValue("age", commons.NewFullPeriod(), 30).
		SetValue("age", commons.NewFullPeriod(), "i am not a good value").
		Build(); err == nil {
		t.Errorf("Expected build error. Got no error: %v", err)
	} else if instance != nil {
		t.Errorf("Expected nil instance. Got non-nil instance: %v", instance)
	}
}

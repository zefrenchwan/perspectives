package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestInstance(t *testing.T) {
	instance := commons.NewTemporalInstance()
	instance.SetAttribute("name", commons.NewFullPeriod(), "john doe")
	if description := instance.Description(); len(description) == 0 {
		t.Log("expected non-empty description, got empty")
		t.Fail()
	} else if len(description) != 1 {
		t.Log("expected only name => string")
		t.Fail()
	} else if description["name"] != "string" {
		t.Log("expected name to be a string")
		t.Fail()
	}

	instance.SetAttribute("size", commons.NewFullPeriod(), 175)
	if description := instance.Description(); len(description) != 2 {
		t.Log("expected name and size")
		t.Fail()
	} else if description["name"] != "string" {
		t.Log("expected name to be a string")
		t.Fail()
	} else if description["size"] != "int" {
		t.Log("expected size to be an int")
		t.Fail()
	}
}

func TestInstanceAttributesMismatch(t *testing.T) {
	instance := commons.NewTemporalInstance()
	instance.SetAttribute("name", commons.NewFullPeriod(), "john doe")
	if err := instance.SetAttribute("name", commons.NewFullPeriod(), 175); err == nil {
		t.Log("expected error setting mismatched attribute type")
		t.Fail()
	}

	if err := instance.SetAttribute("name", commons.NewFullPeriod(), "jane doe"); err != nil {
		t.Log("unexpected error setting attribute name to string")
		t.Fail()
	}

	if value, found := instance.Attribute("name").At(time.Now()); !found || value != "jane doe" {
		t.Log("unexpected attribute value for name")
		t.Fail()
	}
}

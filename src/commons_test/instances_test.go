package commons_test

import (
	"testing"

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

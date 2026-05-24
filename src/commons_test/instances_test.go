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
	} else if description["name"] != commons.StringType.Name() {
		t.Log("expected name to be a string")
		t.Fail()
	}
}

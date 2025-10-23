package commons

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestState(t *testing.T) {
	state := commons.NewStateRepresentation[int]()
	state.SetValue("attr", 10)
	if v, f := state.GetValue("attr"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	} else if attrs := state.Attributes(); slices.Compare([]string{"attr"}, attrs) != 0 {
		t.Fail()
	}

	state.SetValue("other", 100)
	if v, f := state.GetValue("other"); !f {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	} else if v, f := state.GetValue("attr"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	}

	if status := state.Read(); status == nil {
		t.Fail()
	} else if values := status.Values(); len(values) != 2 {
		t.Fail()
	} else if values["attr"] != 10 {
		t.Fail()
	} else if values["other"] != 100 {
		t.Fail()
	}

	// remove non existing field
	if found := state.Remove("cheat"); found {
		t.Fail()
	} else if _, found := state.GetValue("cheat"); found {
		t.Fail()
	}

	// remove existing field
	if found := state.Remove("other"); !found {
		t.Fail()
	} else if _, found := state.GetValue("other"); found {
		t.Fail()
	}

	// test multiple values
	state = commons.NewStateRepresentation[int]()
	newValues := map[string]int{"a": 10, "b": 100}
	state.SetValues(newValues)
	if v, f := state.GetValue("a"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	} else if v, f := state.GetValue("b"); !f {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	}

}

func TestTemporalState(t *testing.T) {
	state := commons.NewTimedStateRepresentation[string](commons.NewFullPeriod())

	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-1, 0, 0)

	period := commons.NewPeriodSince(now, true)
	state.SetValue("name", "paul")
	state.SetActivePeriod(period)

	if !state.Remove("name") {
		t.Fail()
	} else {
		state.SetValue("name", "mickael")
	}

	if _, found := state.GetValue("nope", true); found {
		t.Fail()
	} else if value, found := state.GetValue("name", true); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["mickael"]; !p.Equals(period) {
		t.Fail()
	}

	state.SetValueDuringPeriod("age", "young", period)
	if value, found := state.GetValue("name", true); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["mickael"]; !p.Equals(period) {
		t.Fail()
	} else if value, found := state.GetValue("age", true); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["young"]; !p.Equals(period) {
		t.Fail()
	}

	if desc := state.ReadAtTime(before); desc != nil {
		t.Fail()
	}

	if desc := state.ReadAtTime(now.AddDate(10, 0, 0)); desc == nil {
		t.Fail()
	} else if values := desc.Values(); len(values) != 2 {
		t.Fail()
	} else if values["name"] != "mickael" {
		t.Fail()
	} else if values["age"] != "young" {
		t.Fail()
	}

	if value := state.Read(); !value.ActivePeriod().Equals(period) {
		t.Fail()
	} else if value == nil {
		t.Fail()
	} else if values := value.Values(); len(values) != 2 {
		t.Fail()
	} else if name := values["name"]; len(name) != 1 {
		t.Fail()
	} else if !name["mickael"].Equals(commons.NewFullPeriod()) {
		t.Fail()
	}
}

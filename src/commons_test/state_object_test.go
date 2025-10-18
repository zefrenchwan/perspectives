package commons

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestStateObject(t *testing.T) {
	obj := commons.NewModelStateObject[int]()
	obj.SetValue("attr", 10)
	if v, f := obj.GetValue("attr"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	} else if attrs := obj.Attributes(); slices.Compare([]string{"attr"}, attrs) != 0 {
		t.Fail()
	} else if obj.GetType() != commons.TypeObject {
		t.Fail()
	}

	obj.State.SetValue("other", 100)
	if v, f := obj.GetValue("other"); !f {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	} else if v, f := obj.GetValue("attr"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	}

	obj.Handler().SetValue("final", 1000)
	if v, f := obj.GetValue("other"); !f {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	} else if v, f := obj.GetValue("attr"); !f {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	} else if v, f := obj.GetValue("final"); !f {
		t.Fail()
	} else if v != 1000 {
		t.Fail()
	}

	other := obj
	other.SetValue("cheat", 10000)
	if v, f := obj.GetValue("cheat"); !f {
		t.Fail()
	} else if v != 10000 {
		t.Fail()
	}

	if status := obj.Read(); status.Id() == obj.Id() {
		t.Fail()
	} else if values := status.Values(); len(values) != 4 {
		t.Fail()
	} else if values["attr"] != 10 {
		t.Fail()
	}
}

func TestTemporalStateObject(t *testing.T) {
	obj := commons.NewTemporalModelStateObject[string](commons.NewFullPeriod())
	if obj.GetType() != commons.TypeObject {
		t.Fail()
	}

	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-1, 0, 0)

	period := commons.NewPeriodSince(now, true)
	obj.SetValue("name", "mickael")
	obj.SetActivePeriod(period)

	if _, found := obj.GetValue("nope"); found {
		t.Fail()
	} else if value, found := obj.GetValue("name"); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["mickael"]; !p.Equals(period) {
		t.Fail()
	}

	obj.Handler().SetValueDuringPeriod("age", "young", period)
	if value, found := obj.GetValue("name"); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["mickael"]; !p.Equals(period) {
		t.Fail()
	} else if value, found := obj.GetValue("age"); !found {
		t.Fail()
	} else if len(value) != 1 {
		t.Fail()
	} else if p := value["young"]; !p.Equals(period) {
		t.Fail()
	}

	if desc := obj.ReadAtTime(before); len(desc.Values()) != 0 {
		t.Fail()
	}

	if desc := obj.ReadAtTime(now.AddDate(10, 0, 0)); desc.Id() == obj.Id() {
		t.Fail()
	} else if values := desc.Values(); len(values) != 2 {
		t.Fail()
	} else if values["name"] != "mickael" {
		t.Fail()
	} else if values["age"] != "young" {
		t.Fail()
	}

	if value := obj.Read(); !value.ActivePeriod().Equals(period) {
		t.Fail()
	} else if value.Id() == obj.Id() {
		t.Fail()
	} else if values := value.Values(); len(values) != 2 {
		t.Fail()
	} else if name := values["name"]; len(name) != 1 {
		t.Fail()
	} else if !name["mickael"].Equals(commons.NewFullPeriod()) {
		t.Fail()
	}
}

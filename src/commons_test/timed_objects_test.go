package commons

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestObjectAttributes(t *testing.T) {
	obj := commons.NewTimedStateObject[string]()
	if len(obj.Attributes()) != 0 {
		t.Fail()
	}

	obj.SetValue("test", "value")
	if value := obj.Attributes(); len(value) != 1 {
		t.Fail()
	} else if value[0] != "test" {
		t.Fail()
	}
}

func TestObjectGetValueWithLifetime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	o := commons.NewTimedStateObjectSince[string](after)
	o.SetValueDuring("last name", "Doe", before, now)
	o.SetValueSince("last name", "Dodo", now, true)
	// values are then
	// [before, now[ => Doe
	// [now, +oo[ => Dodo
	matching := commons.NewPeriodSince(after, true)

	if values, found := o.GetValue("last name", true); !found {
		t.Log("expected last name to be present")
		t.Fail()
	} else if len(values) != 1 {
		t.Log("missing values")
		t.Fail()
	} else if period, found := values["Dodo"]; !found {
		t.Log("missing value")
		t.Fail()
	} else if !period.Equals(matching) {
		t.Log("bad period")
		t.Fail()
	}

	if values, found := o.GetValue("last name", false); !found {
		t.Log("expected last name to be present")
		t.Fail()
	} else if len(values) != 2 {
		t.Log(values)
		t.Log("missing values")
		t.Fail()
	} else if period, found := values["Dodo"]; !found {
		t.Log("missing value")
		t.Fail()
	} else if !period.Equals(commons.NewPeriodSince(now, true)) {
		t.Log("bad period")
		t.Fail()
	} else if period, found := values["Doe"]; !found {
		t.Log("missing value")
		t.Fail()
	} else if !period.Equals(commons.NewFinitePeriod(before, now, true, false)) {
		t.Log(period.AsRawString())
		t.Log("bad period")
		t.Fail()
	}
}

func TestObjectGetValuesWithLifetime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	o := commons.NewTimedStateObjectSince[string](after)
	o.SetValueDuring("last name", "Doe", before, now)

	if values := o.GetAllValues(true); len(values) != 1 {
		t.Log("missing fields")
		t.Fail()
	} else if value, found := values["last name"]; !found || len(value) != 0 {
		t.Log("bad field value")
		t.Fail()
	}

	if values := o.GetAllValues(false); len(values) != 1 {
		t.Log("missing fields")
		t.Fail()
	} else if value, found := values["last name"]; !found || len(value) != 1 {
		t.Log("bad field value")
		t.Fail()
	} else if value[0] != "Doe" {
		t.Log("bad content")
		t.Fail()
	}
}

func TestObjectGetValue(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-1, 0, 0)
	object := commons.NewTimedStateObjectSince[string](before)
	object.SetValue("name", "John Doe")
	period := commons.NewPeriodSince(before, true)

	if _, found := object.GetValue("non existing", true); found {
		t.Log("found non existing attribute")
		t.Fail()
	} else if values, found := object.GetValue("name", true); !found {
		t.Log("should find attribute")
		t.Fail()
	} else if len(values) != 1 {
		t.Log("bad values")
		t.Fail()
	} else if p := values["John Doe"]; !p.Equals(period) {
		t.Log("no lifetime intersection")
		t.Fail()
	}
}

func TestObjectTemporalFeatures(t *testing.T) {
	mario := commons.NewTimedStateObject[string]()
	period := commons.NewPeriodSince(time.Now().AddDate(-30, 0, 0), true)

	if !mario.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Log("default value for lifetime is full")
		t.Fail()
	}

	mario.SetActivePeriod(period)
	if !mario.ActivePeriod().Equals(period) {
		t.Log("no impact when changing period")
		t.Fail()
	}
}

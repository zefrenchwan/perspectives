package models_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestObjectTraits(t *testing.T) {
	o := models.NewObject([]string{"Human"})
	if slices.Compare([]string{"Human"}, o.DeclaringTraits()) != 0 {
		t.Fail()
	}
}

func TestObjectAttributes(t *testing.T) {
	o := models.NewObject([]string{"Human"})
	o.AddSemanticForAttribute("first name", "name")
	o.AddSemanticForAttribute("last name", "name")
	o.SetValue("first name", "Jane")
	o.SetValue("last name", "Doe")
	o.AddSemanticForAttribute("other", "a semantic")

	attributes := o.Attributes()
	if len(attributes) != 3 {
		t.Log("missing attributes")
		t.Fail()
	} else if !slices.Contains(attributes, "first name") {
		t.Log("missing attribute")
		t.Fail()
	} else if !slices.Contains(attributes, "last name") {
		t.Log("missing attribute")
		t.Fail()
	} else if !slices.Contains(attributes, "other") {
		t.Log("missing attribute")
		t.Fail()
	}

	if s, found := o.GetSemanticForAttribute("first name"); !found {
		t.Log("missing info for field")
		t.Fail()
	} else if slices.Compare(s, []string{"name"}) != 0 {
		t.Log("failed to read field metadata")
		t.Fail()
	}

	if values := o.GetAllValues(true); len(values) != 3 {
		t.Log(values)
		t.Log("missing fields")
		t.Fail()
	} else if slices.Compare([]string{"Jane"}, values["first name"]) != 0 {
		t.Log("missing content for attribute")
		t.Fail()
	} else if slices.Compare([]string{"Doe"}, values["last name"]) != 0 {
		t.Log("missing content for attribute")
		t.Fail()
	} else if matching, found := values["other"]; !found || len(matching) != 0 {
		t.Log("expecting other => empty")
		t.Fail()
	}
}

func TestObjectGetValueWithLifetime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	o := models.NewObjectSince([]string{"Human"}, after)
	o.SetValueDuring("last name", "Doe", before, now)
	o.SetValueSince("last name", "Dodo", now, true)
	// values are then
	// [before, now[ => Doe
	// [now, +oo[ => Dodo
	matching := structures.NewPeriodSince(after, true)

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
	} else if !period.Equals(structures.NewPeriodSince(now, true)) {
		t.Log("bad period")
		t.Fail()
	} else if period, found := values["Doe"]; !found {
		t.Log("missing value")
		t.Fail()
	} else if !period.Equals(structures.NewFinitePeriod(before, now, true, false)) {
		t.Log(period.AsRawString())
		t.Log("bad period")
		t.Fail()
	}
}

func TestObjectGetValuesWithLifetime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	o := models.NewObjectSince([]string{"Human"}, after)
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
	object := models.NewObjectSince([]string{"Human"}, before)
	object.SetValue("name", "John Doe")
	period := structures.NewPeriodSince(before, true)

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

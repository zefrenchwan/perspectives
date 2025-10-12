package models_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/models"
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
	mario := models.NewObject([]string{"Human"})
	period := commons.NewPeriodSince(time.Now().AddDate(-30, 0, 0), true)

	if !mario.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Log("default value for lifetime is full")
		t.Fail()
	}

	mario.SetActivity(period)
	if !mario.ActivePeriod().Equals(period) {
		t.Log("no impact when changing period")
		t.Fail()
	}
}

func TestObjectDescription(t *testing.T) {
	obj := models.NewObject([]string{"Human"})

	if desc := obj.Describe(); desc.IdObject != obj.Id() {
		t.Log("failed to reference object")
		t.Fail()
	} else if len(desc.Attributes) != 0 {
		t.Log("wrong attributes for object")
		t.Fail()
	} else if slices.Compare(desc.Traits, []string{"Human"}) != 0 {
		t.Log("wrong traits")
		t.Fail()
	}

	obj.SetValue("name", "Cesar")
	obj.SetValue("email", "cesar@rome.it")
	obj.AddSemanticForAttribute("email", "email account")
	desc := obj.Describe()

	if desc.IdObject != obj.Id() {
		t.Log("object id failed")
		t.Fail()
	} else if slices.Compare(desc.Traits, []string{"Human"}) != 0 {
		t.Log("wrong traits")
		t.Fail()
	} else if len(desc.Attributes) != 2 {
		t.Log("wrong attributes")
	}

	for name, semantics := range desc.Attributes {
		switch name {
		case "name":
			if len(semantics) != 0 {
				t.Fail()
			}
		case "email":
			if len(semantics) != 1 {
				t.Fail()
			} else if value := semantics[0]; value != "email account" {
				t.Fail()
			}
		default:
			t.Fail()
		}
	}
}

func TestEmptyObjectBuildFromDescription(t *testing.T) {
	description := models.ObjectDescription{
		Id:       "id",
		IdObject: "id object",
		Traits:   []string{"Human"},
		Attributes: map[string][]string{
			"name":    nil,
			"account": {"email account"},
		},
	}

	object := description.BuildEmptyObjectFromDescription("other id")

	if object.Id() != "other id" {
		t.Log("wrong id")
		t.Fail()
	} else if !object.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Log("should be full")
		t.Fail()
	}

	attributes := object.Attributes()
	if len(attributes) != 2 {
		t.Log("missing attributes")
		t.Fail()
	} else if !slices.Contains(attributes, "name") {
		t.Log("missing name")
		t.Fail()
	} else if !slices.Contains(attributes, "account") {
		t.Log("missing account")
		t.Fail()
	}

	for _, attr := range attributes {
		switch attr {
		case "name":
			if value, found := object.GetSemanticForAttribute(attr); !found || len(value) != 0 {
				t.Fail()
			}

		case "account":
			if value, found := object.GetSemanticForAttribute(attr); !found || len(value) != 1 {
				t.Fail()
			} else if value[0] != "email account" {
				t.Fail()
			}

		default:
			t.Logf("no attr for %s", attr)
		}
	}
}

func TestObjectBuildFromDescription(t *testing.T) {
	base := models.NewObject([]string{"Human"})
	base.SetValue("test", "value")
	base.AddSemanticForAttribute("account", "email account")
	base.SetValue("other field", "value")

	description := base.Describe()
	values := map[string]string{"test": "other value", "account": "dev@dev.com"}

	period := commons.NewPeriodSince(time.Now().Truncate(time.Second), true)
	object := description.BuildObjectFromDescription("id object", period, values)

	if object.Id() != "id object" {
		t.Log("wrong id")
		t.Fail()
	} else if !object.ActivePeriod().Equals(period) {
		t.Log("wrong period")
		t.Fail()
	}

	traits := object.DeclaringTraits()
	if len(traits) != 1 {
		t.Log("missing traits")
		t.Fail()
	} else if !slices.Contains(traits, "Human") {
		t.Log("wrong traits")
		t.Fail()
	}

	attributes := object.Attributes()
	if len(attributes) != 3 {
		t.Log(attributes)
		t.Log("missing attributes")
		t.Fail()
	} else if !slices.Contains(attributes, "test") {
		t.Fail()
	} else if !slices.Contains(attributes, "account") {
		t.Fail()
	} else if !slices.Contains(attributes, "other field") {
		t.Fail()
	}

	allValues := object.GetAllValues(true)
	if len(allValues) != 3 {
		t.Log(allValues)
		t.Log("missing values")
		t.Fail()
	} else if !slices.Equal(allValues["test"], []string{"other value"}) {
		t.Fail()
	} else if !slices.Equal(allValues["other field"], nil) {
		t.Fail()
	} else if !slices.Equal(allValues["account"], []string{"dev@dev.com"}) {
		t.Fail()
	}

}

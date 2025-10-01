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

	if s, found := o.GetSemanticForAttribute("first name"); !found {
		t.Log("missing info for field")
		t.Fail()
	} else if slices.Compare(s, []string{"name"}) != 0 {
		t.Log("failed to read field metadata")
		t.Fail()
	}

	if values := o.GetAllValues(); len(values) != 2 {
		t.Log("missing fields")
		t.Fail()
	} else if slices.Compare([]string{"Jane"}, values["first name"]) != 0 {
		t.Log("missing content for attribute")
		t.Fail()
	} else if slices.Compare([]string{"Doe"}, values["last name"]) != 0 {
		t.Log("missing content for attribute")
		t.Fail()
	}
}

func TestObjectAttributesPartiallyFilled(t *testing.T) {
	o := models.NewObject([]string{"Human"})
	o.AddSemanticForAttribute("first name", "name")
	o.AddSemanticForAttribute("last name", "name")
	o.SetValue("last name", "Doe")

	if values := o.GetAllValues(); len(values) != 1 {
		t.Log("missing fields")
		t.Fail()
	} else if slices.Compare([]string{"Doe"}, values["last name"]) != 0 {
		t.Log("missing content for attribute")
		t.Fail()
	}
}

func TestObjectGetValue(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-1, 0, 0)
	object := models.NewObjectSince([]string{"Human"}, before)
	object.SetValue("name", "John Doe")
	period := structures.NewPeriodSince(before, true)

	if _, found := object.GetValue("non existing"); found {
		t.Log("found non existing attribute")
		t.Fail()
	} else if values, found := object.GetValue("name"); !found {
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

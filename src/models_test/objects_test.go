package models_test

import (
	"slices"
	"testing"

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

func TestObjectDescription(t *testing.T) {
	obj := models.NewObject([]string{"Human"})

	if desc := obj.Describe(); desc.IdObject != obj.Id {
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
	if desc := obj.Describe(); slices.Compare(desc.Attributes, []string{"name"}) != 0 {
		t.Log("wrong attributes for name")
		t.Fail()
	} else if slices.Compare(desc.Traits, []string{"Human"}) != 0 {
		t.Log("wrong traits")
		t.Fail()
	}
}

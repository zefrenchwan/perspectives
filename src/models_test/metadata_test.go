package models_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

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
	obj.SetValue("email", "cesar@rome.it")
	obj.AddSemanticForAttribute("email", "email account")
	desc := obj.Describe()

	if desc.IdObject != obj.Id {
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

	if object.Id != "other id" {
		t.Log("wrong id")
		t.Fail()
	} else if !object.ActivePeriod().Equals(structures.NewFullPeriod()) {
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

	period := structures.NewPeriodSince(time.Now().Truncate(time.Second), true)
	object := description.BuildObjectFromDescription("id object", period, values)

	if object.Id != "id object" {
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

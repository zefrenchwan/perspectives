package structures_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestLoadDependencies(t *testing.T) {
	wValue := "class Worker extends humans.Human"
	hValue := "class Human"
	h := structures.NewDependencies[string]()
	h.SetValue("humans", hValue)
	h.SetValue("workers", wValue)

	if _, f := h.GetValue("none"); f {
		t.Log("unexpected value for none (was not set)")
		t.Fail()
	} else if v, f := h.GetValue("humans"); !f {
		t.Log("should have found humans")
		t.Fail()
	} else if v != hValue {
		t.Log("unexpected value for humans")
		t.Fail()
	}

	h.AddDependency("workers", "humans")

	dependencies := h.LoadWithDependencies("workers")
	if len(dependencies) != 2 {
		t.Log("failed to load dependencies")
		t.Fail()
	} else if value, found := dependencies["workers"]; !found {
		t.Log("failed to load workers")
		t.Fail()
	} else if wValue != value {
		t.Log("failed to load workers value")
		t.Fail()
	} else if value, found := dependencies["humans"]; !found {
		t.Log("failed to load humans")
		t.Fail()
	} else if hValue != value {
		t.Log("failed to load humans value")
		t.Fail()
	}
}

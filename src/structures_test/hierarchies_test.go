package structures_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestLoadDependencies(t *testing.T) {
	wValue := "class Worker extends humans.Human"
	hValue := "class Human"
	h := structures.NewHierarchy[string]()
	h.Set("humans", hValue)
	h.Set("workers", wValue)
	h.LinkToParent("workers", "humans")

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

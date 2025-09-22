package structures_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestHierarchyAncesors(t *testing.T) {
	hierarchy := structures.NewHierarchy[string]()
	hierarchy.SetValue("humans", "people")
	hierarchy.SetValue("men", "male people")
	hierarchy.SetValue("women", "female people")

	if err := hierarchy.AddChildInPartition("men", "humans"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := hierarchy.AddChildInPartition("women", "humans"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if parents, found := hierarchy.Ancestors("women"); !found {
		t.Log("failed to load women -> humans")
		t.Fail()
	} else if len(parents) != 2 {
		t.Log("missing values")
		t.Fail()
	} else if !slices.Contains(parents, "women") {
		t.Log("missing base")
		t.Fail()
	} else if !slices.Contains(parents, "humans") {
		t.Log("missing destination")
		t.Fail()
	}

}

func TestHierarchyChilds(t *testing.T) {
	hierarchy := structures.NewHierarchy[string]()
	hierarchy.SetValue("humans", "people")
	hierarchy.SetValue("men", "male people")
	hierarchy.SetValue("women", "female people")
	hierarchy.AddChildInPartition("men", "humans")
	hierarchy.AddChildInPartition("women", "humans")

	if values, _ := hierarchy.Childs("pizza"); len(values) != 0 {
		t.Fail()
	}

	if values, exclusive := hierarchy.Childs("humans"); len(values) != 2 {
		t.Log("values mismatch")
		t.Fail()
	} else if !slices.Contains(values, "men") {
		t.Log("missing men")
		t.Fail()
	} else if !slices.Contains(values, "women") {
		t.Log("missing women")
		t.Fail()
	} else if !exclusive {
		t.Log("exclusive failed")
		t.Fail()
	}
}

func TestMismatch(t *testing.T) {
	hierarchy := structures.NewHierarchy[string]()
	hierarchy.SetValue("humans", "people")
	hierarchy.SetValue("men", "male people")
	hierarchy.SetValue("women", "female people")
	hierarchy.AddChildInPartition("men", "humans")
	if err := hierarchy.AddChildToParent("women", "humans"); err == nil {
		t.Log("missing exception")
		t.Fail()
	}
}

func TestElements(t *testing.T) {
	hierarchy := structures.NewHierarchy[string]()
	hierarchy.SetValue("humans", "people")
	hierarchy.SetValue("men", "male people")
	hierarchy.SetValue("women", "female people")
	hierarchy.AddChildInPartition("men", "humans")

	got := hierarchy.Elements()
	expected := map[string]string{
		"humans": "people",
		"men":    "male people",
		"women":  "female people",
	}

	if len(got) != len(expected) {
		t.Fail()
	} else {
		for k, v := range expected {
			if value, found := got[k]; !found {
				t.Fail()
			} else if value != v {
				t.Fail()
			}
		}
	}
}

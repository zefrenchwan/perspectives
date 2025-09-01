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

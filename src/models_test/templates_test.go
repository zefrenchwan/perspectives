package models

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestInheritance(t *testing.T) {
	dogs := models.NewFormalClass("dogs")
	animals := models.NewFormalClass("animals")
	hierarchy := models.NewFormalHierarchy()
	hierarchy.SetClass(dogs)
	hierarchy.SetClass(animals)

	// register link
	if err := hierarchy.AddChildClass(dogs.Name, animals.Name, true); err != nil {
		t.Log(err)
		t.Fail()
	}

	// test inheritance

	// case 0: not present, no value
	if values := hierarchy.GetClassHierarchy("no value"); len(values) != 0 {
		t.Log("no class expected")
		t.Fail()
	}

	// case 1: should get animals only because no parent registered
	if values := hierarchy.GetClassHierarchy(animals.Name); len(values) != 1 {
		t.Log("missing class in class tree")
		t.Fail()
	} else if values[0].Id != animals.Id {
		t.Log("wrong class")
		t.Fail()
	}

	// case 2: get the full tree
	if values := hierarchy.GetClassHierarchy(dogs.Name); len(values) != 2 {
		t.Log("missing class in class tree")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(value models.FormalClass) bool { return value.Id == dogs.Id }) {
		t.Log("missing base class")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(value models.FormalClass) bool { return value.Id == animals.Id }) {
		t.Log("missing base class")
		t.Fail()
	}
}

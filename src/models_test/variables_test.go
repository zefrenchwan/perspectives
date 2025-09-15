package models_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestMapToObject(t *testing.T) {
	x := models.NewVariableForObject("x", []string{"Human"})
	tiramisu := models.NewObject([]string{"dessert"})
	if _, err := x.MapAs(tiramisu); err == nil {
		t.Log("traits mismatch")
		t.Fail()
	}

	worker := models.NewObject([]string{"Human"})
	if v, err := x.MapAs(worker); err != nil {
		t.Log("traits match but raised error")
		t.Fail()
	} else if v.GetType() != models.EntityTypeObject {
		t.Log("wrong type")
		t.Fail()
	} else if o, err := v.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if o.Id != worker.Id {
		t.Log("wrong object")
		t.Fail()
	}
}

func TestMapToGroup(t *testing.T) {
	x := models.NewVariableForGroup("x", []string{"Human"})
	jane := models.NewObject([]string{"Human"})
	lara := models.NewObject([]string{"Human"})
	cherry := models.NewObject([]string{"Food"})

	validGroup := []models.Object{jane, lara}
	invalidGroup := []models.Object{lara, cherry}

	if g, err := x.MapAs(validGroup); err != nil {
		t.Log(err)
		t.Fail()
	} else if _, err := x.MapAs(invalidGroup); err == nil {
		t.Log("failed to detect non matching element")
		t.Fail()
	} else if group, err := g.AsGroup(); err != nil {
		t.Log(err)
		t.Fail()
	} else if !slices.ContainsFunc(group, func(v models.Object) bool { return v.Id == jane.Id }) {
		t.Log("missing element")
		t.Fail()
	} else if !slices.ContainsFunc(group, func(v models.Object) bool { return v.Id == lara.Id }) {
		t.Log("missing element")
		t.Fail()
	} else if len(group) != 2 {
		t.Log("wrong elements")
		t.Fail()
	}

}

func TestMapToTrait(t *testing.T) {
	x := models.NewVariableForTrait("x")
	dogs := models.NewTrait("dogs")
	obj := models.NewObject([]string{"any"})

	if d, err := x.MapAs(dogs); err != nil {
		t.Log(err)
		t.Fail()
	} else if trait, err := d.AsTrait(); err != nil {
		t.Log(err)
		t.Fail()
	} else if !trait.Equals(dogs) {
		t.Log("wrong value")
		t.Fail()
	}

	if _, err := x.MapAs(obj); err == nil {
		t.Log("wrong type for mapping")
		t.Fail()
	}
}

func TestMapToSpecificTrait(t *testing.T) {
	x := models.NewVariableForSpecificTraits("x", []string{"dogs", "cats"})
	dogs := models.NewTrait("dogs")
	cheese := models.NewTrait("cheese")

	if d, err := x.MapAs(dogs); err != nil {
		t.Log(err)
		t.Fail()
	} else if trait, err := d.AsTrait(); err != nil {
		t.Log(err)
		t.Fail()
	} else if !trait.Equals(dogs) {
		t.Log("wrong value")
		t.Fail()
	}

	if _, err := x.MapAs(cheese); err == nil {
		t.Log("impossible map because trait not in accepted traits")
		t.Fail()
	}
}

func TestMapToLink(t *testing.T) {
	maria := models.NewObject([]string{"Human"})
	spain := models.NewObject([]string{"Country"})
	likes, _ := models.NewSimpleLink("loves", maria, spain)

	x := models.NewVariableForLink("x")

	if _, err := x.MapAs(maria); err == nil {
		t.Log("wrong mapping, object cannot replace link")
		t.Fail()
	}

	if l, err := x.MapAs(likes); err != nil {
		t.Log(err)
		t.Fail()
	} else if link, err := l.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Id() != likes.Id() {
		t.Log("wrong match")
		t.Fail()
	}
}

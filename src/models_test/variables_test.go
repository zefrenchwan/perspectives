package models_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func testIfEntityIsObjectWithId(e models.Entity, objectId string) bool {
	if e == nil {
		return false
	} else if e.GetType() != models.EntityTypeObject {
		return false
	} else if o, err := models.AsObject(e); err != nil {
		return false
	} else {
		return o.Id() == objectId
	}
}

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
	} else if o, err := models.AsObject(v); err != nil {
		t.Log(err)
		t.Fail()
	} else if o.Id() != worker.Id() {
		t.Log("wrong object")
		t.Fail()
	}
}

func TestMapToGroup(t *testing.T) {
	x := models.NewVariableForGroup("x", []string{"Human"})
	jane := models.NewObject([]string{"Human"})
	lara := models.NewObject([]string{"Human"})
	cherry := models.NewObject([]string{"Food"})

	validGroup := []*models.Object{jane, lara}
	invalidGroup := []*models.Object{lara, cherry}

	if g, err := x.MapAs(validGroup); err != nil {
		t.Log(err)
		t.Fail()
	} else if _, err := x.MapAs(invalidGroup); err == nil {
		t.Log("failed to detect non matching element")
		t.Fail()
	} else if group, err := models.AsGroup(g); err != nil {
		t.Log(err)
		t.Fail()
	} else if !slices.ContainsFunc(group, func(e models.Entity) bool { return testIfEntityIsObjectWithId(e, jane.Id()) }) {
		t.Log("missing element")
		t.Fail()
	} else if !slices.ContainsFunc(group, func(e models.Entity) bool { return testIfEntityIsObjectWithId(e, lara.Id()) }) {
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
	} else if trait, err := models.AsTrait(d); err != nil {
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
	} else if trait, err := models.AsTrait(d); err != nil {
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
	} else if link, err := models.AsLink(l); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Id() != likes.Id() {
		t.Log("wrong match")
		t.Fail()
	}
}

func TestMapToLinksGroupShouldFail(t *testing.T) {
	x := models.NewVariableForGroup("x", []string{"Human"})
	marie := models.NewObject([]string{"Human"})
	anna := models.NewObject([]string{"Human"})
	am, _ := models.NewSimpleLink("knows", anna, marie)
	ma, _ := models.NewSimpleLink("knows", marie, anna)

	objectsGroup, _ := models.NewGroupOfObjects(anna, marie)
	mixedGroup, _ := models.NewLinksGroup([]*models.Link{am, ma})
	// objects should map
	if !x.Matches(objectsGroup) {
		t.Log("objects should match")
		t.Fail()
	}

	// links as groups should NOT
	if x.Matches(mixedGroup) {
		t.Log("links should fail")
		t.Fail()
	}
}

func TestMatchesObjectsOrGroups(t *testing.T) {
	variable := models.NewVariableForObject("x", []string{"Human"})
	dog := models.NewObject([]string{"Dog"})
	human := models.NewObject([]string{"Human"})
	link, _ := models.NewSimpleLink("owns", human, dog)

	if variable.Matches(link) {
		t.Log("different type, should refuse")
		t.Fail()
	} else if variable.Matches(dog) {
		t.Log("different accepted traits, should refuse")
		t.Fail()
	} else if !variable.Matches(human) {
		t.Log("same traits, should accept")
		t.Fail()
	}

	gVar := models.NewVariableForGroup("y", []string{"Human", "Monkey"})
	values, _ := models.NewObjectsGroup([]*models.Object{human, dog})
	if gVar.Matches(values) {
		t.Log("dog is neither human nor monkey, should stop")
		t.Fail()
	}

	values, _ = models.NewObjectsGroup([]*models.Object{human})
	if !gVar.Matches(values) {
		t.Log("human accepted")
		t.Fail()
	}
}

func TestMatchesLink(t *testing.T) {
	lindsley := models.NewObject([]string{"Musician"})
	iceStorm := models.NewObject([]string{"Song"})
	link, _ := models.NewSimpleLink("wrote", lindsley, iceStorm)

	variable := models.NewVariableForLink("x")
	if !variable.Matches(link) {
		t.Fail()
	} else if variable.Matches(lindsley) {
		t.Log("link variable cannot match an object")
		t.Fail()
	}
}

func TestMatchesVariables(t *testing.T) {
	base := models.NewVariableForTrait("x")
	otherTrait := models.NewVariableForSpecificTraits("y", []string{"Human"})
	otherGenericTrait := models.NewVariableForTrait("z")

	if base.Matches(otherTrait) {
		t.Fail()
	} else if !base.Matches(otherGenericTrait) {
		t.Fail()
	}

	base = models.NewVariableForLink("z")
	if base.Matches(otherTrait) {
		t.Fail()
	}
}

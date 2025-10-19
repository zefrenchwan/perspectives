package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestSetStateAction(t *testing.T) {
	action := commons.NewSetStateAction("x", "price", 10)
	object := commons.NewModelStateObject[int]()
	object.SetValue("price", 1000)

	// bad type
	other := DummyIdBasedImplementation{id: "sure"}
	content := commons.NewNamedContent[commons.Modelable]("x", other)
	if err := action.Execute(content); err != nil {
		t.Fail()
	}

	// bad variable
	content = commons.NewNamedContent[commons.Modelable]("x", object)
	action = commons.NewSetStateAction("y", "other attr", 10)
	if err := action.Execute(content); err != nil {
		t.Fail()
	} else if _, found := object.GetValue("other attr"); found {
		t.Fail()
	}

	// matching single value
	object = commons.NewModelStateObject[int]()
	object.SetValue("price", 1000)
	content = commons.NewNamedContent[commons.Modelable]("x", object)
	action = commons.NewSetStateAction("x", "price", 10)
	if err := action.Execute(content); err != nil {
		t.Fail()
	} else if v, found := object.GetValue("price"); !found {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	}

	// matching multiple values
	object = commons.NewModelStateObject[int]()
	object.SetValue("price", 1000)
	action = commons.NewSetStateActionFrom("x", map[string]int{"price": 5, "status": 170})
	content = commons.NewNamedContent[commons.Modelable]("x", object)
	action.Execute(content)
	if v, found := object.GetValue("price"); !found {
		t.Fail()
	} else if v != 5 {
		t.Fail()
	} else if v, found := object.GetValue("status"); !found {
		t.Fail()
	} else if v != 170 {
		t.Fail()
	}
}

func TestSequentialActions(t *testing.T) {
	object := commons.NewModelStateObject[int]()
	object.SetValue("status", 1000)
	content := commons.NewNamedContent[commons.Modelable]("x", object)
	content.AppendAs("y", object)
	action := commons.NewSetStateAction("x", "price", 10)
	other := commons.NewSetStateAction("y", "price", 100)

	// first, test two actions
	composite := commons.NewSequentialActions([]commons.ExecuteAction{action, other}, true)
	if variables := composite.Signature().Variables(); len(variables) != 2 {
		t.Fail()
	} else if !slices.Contains(variables, "x") {
		t.Fail()
	} else if !slices.Contains(variables, "y") {
		t.Fail()
	} else if err := composite.Execute(content); err != nil {
		t.Fail()
	} else if v, found := object.GetValue("price"); !found {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	} else if v, found := object.GetValue("price"); !found {
		t.Fail()
	} else if v != 100 {
		t.Fail()
	}

	// no action
	composite = commons.NewSequentialActions([]commons.ExecuteAction{}, true)
	if variables := composite.Signature().Variables(); len(variables) != 0 {
		t.Fail()
	} else if err := composite.Execute(content); err != nil {
		t.Fail()
	}
}

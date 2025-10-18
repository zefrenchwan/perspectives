package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestStateObject(t *testing.T) {
	obj := commons.NewStateObject[int]()

	if len(obj.Attributes()) != 0 {
		t.Fail()
	} else if _, found := obj.GetValue("test"); found {
		t.Fail()
	}

	obj.SetValue("key", 10)
	if value, found := obj.GetValue("key"); !found {
		t.Fail()
	} else if value != 10 {
		t.Fail()
	} else if attr := obj.Attributes(); len(attr) != 1 {
		t.Fail()
	} else if attr[0] != "key" {
		t.Fail()
	}
}

func TestStateSetValueAction(t *testing.T) {
	obj := commons.NewStateObject[int]()
	action := commons.StateSetValueAction[int]{
		Variable:  "x",
		Attribute: "attr",
		NewValue:  10,
	}

	other := DummyIdBasedImplementation{}

	// test variable mismatch
	content := commons.NewNamedContent[commons.Modelable]("y", obj)
	if err := action.Execute(content); err == nil {
		t.Fail()
	}

	content = commons.NewNamedContent[commons.Modelable]("x", other)
	if err := action.Execute(content); err == nil {
		t.Fail()
	}

	content = commons.NewNamedContent[commons.Modelable]("x", obj)
	if err := action.Execute(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if v, found := obj.GetValue("attr"); !found {
		t.Fail()
	} else if v != action.NewValue {
		t.Fail()
	}
}

func TestReadStatusFromStatedObjects(t *testing.T) {
	obj := commons.NewStateObject[string]()
	obj.SetValue("attr", "test")

	if status := obj.Read(); status == nil {
		t.Fail()
	} else if id, found := status.Id(); !found {
		t.Fail()
	} else if id != obj.Id() {
		t.Fail()
	} else if values := status.Values(); len(values) != 1 {
		t.Fail()
	} else if values["attr"] != "test" {
		t.Fail()
	}
}

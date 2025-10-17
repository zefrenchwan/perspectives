package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestFilterById(t *testing.T) {
	condition := commons.NewFilterById("x", "id")

	// test variables check
	otherVariable := commons.NewNamedContent[commons.Modelable]("y", DummyIdBasedImplementation{id: "id"})
	if value, err := condition.Matches(otherVariable); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}

	// test values condition
	matching := commons.NewNamedContent[commons.Modelable]("x", DummyIdBasedImplementation{id: "id"})
	notMatching := commons.NewNamedContent[commons.Modelable]("x", DummyIdBasedImplementation{id: "nope"})
	notId := commons.NewNamedContent[commons.Modelable]("x", DummyComponentImplementation{})
	if value, err := condition.Matches(notMatching); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	} else if value, err := condition.Matches(matching); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	} else if value, err := condition.Matches(notId); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}

}

func TestTypesMatches(t *testing.T) {
	a := DummyComponentImplementation{}
	content := commons.NewNamedContent[commons.Modelable]("x", a)

	condition := commons.NewFilterByTypes("y", []commons.ModelableType{DummyTestingType})
	if value, err := condition.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("bad variable")
		t.Fail()
	}

	condition = commons.NewFilterByTypes("x", []commons.ModelableType{DummyTestingType})
	if value, err := condition.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Log("should match: same variable, good type")
		t.Fail()
	}

	condition = commons.NewFilterByTypes("y", []commons.ModelableType{commons.TypeConstraint})
	if value, err := condition.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("bad type")
		t.Fail()
	}
}

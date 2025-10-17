package commons_test

import (
	"testing"
	"time"

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

func TestFilterByTypes(t *testing.T) {
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

func TestCompareActivePeriod(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-10, 0, 0)
	activity := commons.NewPeriodUntil(before, true)

	obj := commons.NewStateObject[string]()
	tobj := commons.NewTimedStateObjectSince[string](now)
	matching := commons.NewTimedStateObject[string]()

	condition := commons.NewFilterActivePeriod("x", commons.TemporalCommonPoint, activity)

	content := commons.NewNamedContent[commons.Modelable]("x", obj)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if matches {
		t.Fail()
	}

	content = commons.NewNamedContent[commons.Modelable]("x", tobj)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if matches {
		t.Fail()
	}

	content = commons.NewNamedContent[commons.Modelable]("x", matching)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if !matches {
		t.Fail()
	}
}

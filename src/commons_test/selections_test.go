package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestFilterById(t *testing.T) {
	condition := commons.NewFilterById("x", "id")

	// test variables check
	otherVariable := commons.NewNamedContent("y", DummyIdBasedImplementation{id: "id"})
	if value, err := condition.Matches(otherVariable); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}

	// test values condition
	matching := commons.NewNamedContent("x", DummyIdBasedImplementation{id: "id"})
	notMatching := commons.NewNamedContent("x", DummyIdBasedImplementation{id: "nope"})
	notId := commons.NewNamedContent("x", DummyComponentImplementation{})
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
	content := commons.NewNamedContent("x", a)

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

	obj := DummyObject{id: "i have no time"}
	nonMatching := commons.NewStateObjectSince[string](now)
	matching := commons.NewTemporalStateObject[string](commons.NewFullPeriod())

	condition := commons.NewFilterActivePeriod("x", commons.TemporalCommonPoint, activity)

	content := commons.NewNamedContent("x", obj)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if matches {
		t.Fail()
	}

	content = commons.NewNamedContent("x", nonMatching)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if matches {
		t.Log("content starts now, expected ends before")
		t.Fail()
	}

	content = commons.NewNamedContent("x", matching)
	if matches, err := condition.Matches(content); err != nil {
		t.Fail()
	} else if !matches {
		t.Log("full period in content matches activity")
		t.Fail()
	}
}

func TestFilterByStateAttribute(t *testing.T) {
	operator := commons.NewFilterByStateOperator("x", "name", commons.StringEqualsIgnoreCase, "Oriane")

	if variables := operator.Signature().Variables(); len(variables) != 1 {
		t.Fail()
	} else if variables[0] != "x" {
		t.Fail()
	}

	// bad variable
	content := commons.NewNamedContent("y", DummyObject{id: "whatever"})
	if m, err := operator.Matches(content); err != nil {
		t.Fail()
	} else if m {
		t.Log("bad variable")
		t.Fail()
	}

	// bad type for content
	content = commons.NewNamedContent("x", DummyObject{id: "whatever"})
	if m, err := operator.Matches(content); err != nil {
		t.Fail()
	} else if m {
		t.Log("bad type")
		t.Fail()
	}

	// no match due to missing value
	value := commons.NewStateObject[string]()
	content = commons.NewNamedContent("x", value)
	if m, err := operator.Matches(content); err != nil {
		t.Fail()
	} else if m {
		t.Log("no value")
		t.Fail()
	}

	// no value match
	value = commons.NewStateObject[string]()
	value.SetValue("name", "Thomas")
	content = commons.NewNamedContent("x", value)
	if m, err := operator.Matches(content); err != nil {
		t.Fail()
	} else if m {
		t.Log("no match")
		t.Fail()
	}

	// match
	value = commons.NewStateObject[string]()
	value.SetValue("name", "Oriane")
	content = commons.NewNamedContent("x", value)
	if m, err := operator.Matches(content); err != nil {
		t.Fail()
	} else if !m {
		t.Log("should match")
		t.Fail()
	}
}

func TestFilterByStateSetAttribute(t *testing.T) {
	setOperator := commons.NewLocalSetOperator(commons.MatchesOneInSetOperator, commons.IntEquals)
	operator := commons.NewFilterByStateSetOperator("x", "age", setOperator, []int{10, 20, 30})

	if variables := operator.Signature().Variables(); len(variables) != 1 {
		t.Fail()
	} else if variables[0] != "x" {
		t.Fail()
	}

	value := commons.NewStateObject[int]()
	value.SetValue("age", 20)
	content := commons.NewNamedContent("x", value)

	// test match
	if m, err := operator.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if !m {
		t.Fail()
	}

	// test bad variable
	value.SetValue("age", 20)
	content = commons.NewNamedContent("y", value)

	if m, err := operator.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if m {
		t.Fail()
	}

	// test bad value
	value.SetValue("age", 25)
	content = commons.NewNamedContent("x", value)

	if m, err := operator.Matches(content); err != nil {
		t.Log(err)
		t.Fail()
	} else if m {
		t.Fail()
	}

}

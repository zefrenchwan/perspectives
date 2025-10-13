package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestIdBasedCondition(t *testing.T) {
	accepting := DummyIdBasedImplementation{id: "id"}
	refusing := DummyIdBasedImplementation{id: "refused"}
	condition := commons.IdBasedCondition{Id: "id"}

	p := commons.NewNamedParameter("x", accepting)
	p.AppendAsVariable("y", refusing)

	if condition.Matches(p) {
		t.Log("multiple values should not match")
		t.Fail()
	}

	p = commons.NewNamedParameter("y", refusing)
	if condition.Matches(p) {
		t.Log("bad id matching")
		t.Fail()
	}

	p = commons.NewParameter(accepting)
	if !condition.Matches(p) {
		t.Log("id should match")
		t.Fail()
	}

	// test not implementing
	var empty DummyBasicModelElementImplementation
	p = commons.NewParameter(empty)
	if condition.Matches(p) {
		t.Log("should refuse non id based")
		t.Fail()
	}

	// test for nil
	p = commons.NewParameter(nil)
	if condition.Matches(p) {
		t.Log("nil cannot match")
		t.Fail()
	}
}

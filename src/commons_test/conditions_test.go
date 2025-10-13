package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyIdBasedImplementation struct {
	id string
}

func (d DummyIdBasedImplementation) Id() string {
	return d.id
}

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
}

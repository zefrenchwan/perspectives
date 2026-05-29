package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestLinks(t *testing.T) {
	john := commons.NewTemporalInstance()
	pizza := commons.NewTemporalInstance()
	link := commons.NewLink("loves")
	link.Add("subject", john)
	link.Add("object", pizza)

	if link.Name() != "loves" {
		t.Errorf("Expected link name to be 'loves', got '%s'", link.Name())
		t.Fail()
	}

	if values, found := link.Has("subject"); !found {
		t.Fail()
	} else if !slices.EqualFunc(values, []commons.Linkable{john}, func(a, b commons.Linkable) bool { return a.Same(b) }) {
		t.Fail()
	} else if values, found := link.Has("object"); !found {
		t.Fail()
	} else if !slices.EqualFunc(values, []commons.Linkable{pizza}, func(a, b commons.Linkable) bool { return a.Same(b) }) {
		t.Fail()
	} else if _, found := link.Has("other"); found {
		t.Fail()
	}

	if operands := link.Operands(); len(operands) != 2 {
		t.Fail()
	} else if svalue, sfound := operands["subject"]; !sfound {
		t.Fail()
	} else if !slices.EqualFunc(svalue, []commons.Linkable{john}, func(a, b commons.Linkable) bool { return a.Same(b) }) {
		t.Fail()
	} else if ovalue, ofound := operands["object"]; !ofound {
		t.Fail()
	} else if !slices.EqualFunc(ovalue, []commons.Linkable{pizza}, func(a, b commons.Linkable) bool { return a.Same(b) }) {
		t.Fail()
	}

	link.Remove("subject", func(linkable commons.Linkable) bool {
		return linkable.Same(john)
	})
	if _, found := link.Has("subject"); found {
		t.Fail()
	}
}

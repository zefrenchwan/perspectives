package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestVariableReplacement(t *testing.T) {
	instance := commons.NewTemporalInstance()
	trait := commons.NewTrait("trait")

	if variable := commons.NewVariable("x"); !variable.CanBeReplacedBy(instance) {
		t.Errorf("Variable %s with no constraint should be replaced by instance", variable.Name())
	} else if !variable.CanBeReplacedBy(trait) {
		t.Errorf("Variable %s with no constraint should be replaced by trait", variable.Name())
	}

	if variable := commons.NewVariable("y", commons.CLASS_INSTANCE); !variable.CanBeReplacedBy(instance) {
		t.Errorf("Variable %s with constraint should be replaced by instance", variable.Name())
	} else if variable.CanBeReplacedBy(trait) {
		t.Errorf("Variable %s with constraint should not be replaced by trait", variable.Name())
	}
}

func TestVariableMatching(t *testing.T) {
	a, b := commons.NewTemporalInstance(), commons.NewTemporalInstance()
	varA, varB := commons.NewVariable("varA"), commons.NewVariable("varB")

	link := commons.NewLink("same", commons.NewFullPeriod()).
		WithOperand("subject", []commons.Element{a}).
		WithOperand("object", []commons.Element{b})

	if !link.Same(link) {
		t.Errorf("Link should be same as itself")
	} else if sub, ok := commons.Match(link, link); !ok {
		t.Errorf("Link should match itself")
	} else if len(sub) != 0 {
		t.Errorf("Link should match itself with no substitutions")
	}

	notMatchingName := commons.NewLink("not same name", commons.NewFullPeriod()).
		WithOperand("subject", []commons.Element{varA}).
		WithOperand("object", []commons.Element{varB})
	if _, ok := commons.Match(notMatchingName, link); ok {
		t.Errorf("Link with different name should not match")
	}

	notMatchingPeriod := commons.NewLink("same", commons.NewPeriodSince(time.Now(), true)).
		WithOperand("subject", []commons.Element{varA}).
		WithOperand("object", []commons.Element{varB})
	if _, ok := commons.Match(notMatchingPeriod, link); ok {
		t.Errorf("Link with different period should not match")
	}

	notMatchingStructure := commons.NewLink("same", commons.NewFullPeriod()).
		WithOperand("subject", []commons.Element{varA})
	if _, ok := commons.Match(notMatchingStructure, link); ok {
		t.Errorf("Link with different structure should not match")
	}

	matching := commons.NewLink("same", commons.NewFullPeriod()).
		WithOperand("subject", []commons.Element{varA}).
		WithOperand("object", []commons.Element{varB})
	if sub, ok := commons.Match(matching, link); !ok {
		t.Errorf("Link with same structure should match")
	} else if len(sub) != 2 {
		t.Errorf("substitution should be varA => a, varB => b")
	} else if !sub["varA"].Same(a) {
		t.Errorf("substitution should be varA => a, varB => b")
	} else if !sub["varB"].Same(b) {
		t.Errorf("substitution should be varA => a, varB => b")
	}
}

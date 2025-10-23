package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestConstantCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewContent(DummyComponentImplementation{})

	if !p.Matches(ctrue.Signature()) {
		t.Fail()
	} else if value, err := ctrue.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	}

	if !p.Matches(cfalse.Signature()) {
		t.Fail()
	} else if value, err := cfalse.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}
}

func TestNotCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewContent(DummyComponentImplementation{})

	if not := commons.NewConditionNot(ctrue); !p.Matches(not.Signature()) {
		t.Log("failed to accept signature")
		t.Fail()
	} else if value, err := not.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("unexpected result: not true is false")
		// not true is false
		t.Fail()
	}

	if not := commons.NewConditionNot(cfalse); !p.Matches(not.Signature()) {
		t.Log("failed to accept signature")
		t.Fail()
	} else if value, err := not.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Log("unexpected result: not false is true")
		// not false is true
		t.Fail()
	}

	// special case: nil
	if not := commons.NewConditionNot(nil); !p.Matches(not.Signature()) {
		t.Log("failed to accept nil")
		t.Fail()
	} else if value, err := not.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("should refuse when applied to nil")
	}
}

func TestOrCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewContent(DummyComponentImplementation{})

	if or := commons.NewConditionOr([]commons.Condition{ctrue, cfalse}); !p.Matches(or.Signature()) {
		t.Fail()
	} else if value, err := or.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	} else if value, err := commons.NewConditionOr([]commons.Condition{ctrue, ctrue}).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	} else if value, err := commons.NewConditionOr([]commons.Condition{cfalse, cfalse}).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}

	// empty or nil should accept p but return false
	if or := commons.NewConditionOr([]commons.Condition{}); !p.Matches(or.Signature()) {
		t.Fail()
	} else if value, err := or.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("empty or nil should return false")
		t.Fail()
	} else if value, err := commons.NewConditionOr(nil).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("should refuse when applied to nil")
	}
}

func TestAndCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewContent(DummyComponentImplementation{})

	if and := commons.NewConditionAnd([]commons.Condition{ctrue, cfalse}); !p.Matches(and.Signature()) {
		t.Fail()
	} else if value, err := and.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	} else if value, err := commons.NewConditionAnd([]commons.Condition{ctrue, ctrue}).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	} else if value, err := commons.NewConditionAnd([]commons.Condition{cfalse, cfalse}).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}

	// empty or nil should return false
	if value, err := commons.NewConditionAnd([]commons.Condition{}).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("empty or nil should return false")
		t.Fail()
	} else if value, err := commons.NewConditionAnd(nil).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("should refuse when applied to nil")
	}
}

func TestCompositeCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	// or is true all the times because it contains ctrue
	or := commons.NewConditionOr([]commons.Condition{
		ctrue, commons.NewConditionNot(cfalse), commons.NewConditionNot(commons.NewConditionNot(cfalse)),
	})

	// condition is true because all operands are true
	condition := commons.NewConditionAnd([]commons.Condition{or, ctrue})

	p := commons.NewContent(DummyComponentImplementation{})
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Log("condition is true")
		t.Fail()
	}
}

func TestFormalParametersTree(t *testing.T) {
	xCondition := commons.NewFilterById("x", "id")
	yCondition := commons.NewFilterById("y", "id")
	or := commons.NewConditionOr([]commons.Condition{xCondition, yCondition})

	if variables := or.Signature().Variables(); len(variables) != 2 {
		t.Fail()
	} else if !slices.Contains(variables, "x") {
		t.Fail()
	} else if !slices.Contains(variables, "y") {
		t.Fail()
	}

	content := commons.NewNamedContent("x", DummyComponentImplementation{})
	if content.Matches(or.Signature()) {
		t.Fail()
	}

	content.AppendAs("y", DummyComponentImplementation{})
	if !content.Matches(or.Signature()) {
		t.Fail()
	}
}

package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestConstantCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewParameter(DummyBasicModelElementImplementation{})

	if value, err := ctrue.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	}

	if value, err := cfalse.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Fail()
	}
}

func TestNotCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewParameter(DummyBasicModelElementImplementation{})

	if value, err := commons.NewConditionNot(ctrue).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		// not true is false
		t.Fail()
	}

	if value, err := commons.NewConditionNot(cfalse).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		// not false is true
		t.Fail()
	}

	// special case: nil
	if value, err := commons.NewConditionNot(nil).Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("should refuse when applied to nil")
	}
}

func TestOrCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewParameter(DummyBasicModelElementImplementation{})

	if value, err := commons.NewConditionOr([]commons.Condition{ctrue, cfalse}).Matches(p); err != nil {
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

	// empty or nil should return false
	if value, err := commons.NewConditionOr([]commons.Condition{}).Matches(p); err != nil {
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

	p := commons.NewParameter(DummyBasicModelElementImplementation{})

	if value, err := commons.NewConditionAnd([]commons.Condition{ctrue, cfalse}).Matches(p); err != nil {
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

	p := commons.NewParameter(DummyBasicModelElementImplementation{})
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Log("condition is true")
		t.Fail()
	}
}

func TestIdBasedCondition(t *testing.T) {
	accepting := DummyIdBasedImplementation{id: "id"}
	refusing := DummyIdBasedImplementation{id: "refused"}
	condition := commons.IdBasedCondition{Id: "id"}

	p := commons.NewNamedParameter("x", accepting)
	p.AppendAsVariable("y", refusing)

	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("multiple values should not match")
		t.Fail()
	}

	p = commons.NewNamedParameter("y", refusing)
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("bad id matching")
		t.Fail()
	}

	p = commons.NewParameter(accepting)
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Log("id should match")
		t.Fail()
	}

	// test not implementing
	var empty DummyBasicModelElementImplementation
	p = commons.NewParameter(empty)
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("should refuse non id based")
		t.Fail()
	}

	// test for nil
	p = commons.NewParameter(nil)
	if value, err := condition.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		t.Log("nil cannot match")
		t.Fail()
	}
}

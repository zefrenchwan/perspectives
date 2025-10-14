package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestConstantCondition(t *testing.T) {
	ctrue := commons.NewConditionConstant(true)
	cfalse := commons.NewConditionConstant(false)

	p := commons.NewContent(DummyComponentImplementation{})

	if !ctrue.Signature().Accepts(p) {
		t.Fail()
	} else if value, err := ctrue.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		t.Fail()
	}

	if !cfalse.Signature().Accepts(p) {
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

	if not := commons.NewConditionNot(ctrue); !not.Signature().Accepts(p) {
		t.Fail()
	} else if value, err := not.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if value {
		// not true is false
		t.Fail()
	}

	if not := commons.NewConditionNot(cfalse); !not.Signature().Accepts(p) {
		t.Fail()
	} else if value, err := not.Matches(p); err != nil {
		t.Log(err)
		t.Fail()
	} else if !value {
		// not false is true
		t.Fail()
	}

	// special case: nil
	if not := commons.NewConditionNot(nil); !not.Signature().Accepts(p) {
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

	if or := commons.NewConditionOr([]commons.Condition{ctrue, cfalse}); !or.Signature().Accepts(p) {
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
	if or := commons.NewConditionOr([]commons.Condition{}); !or.Signature().Accepts(p) {
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

	if and := commons.NewConditionAnd([]commons.Condition{ctrue, cfalse}); !and.Signature().Accepts(p) {
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

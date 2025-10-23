package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestTemporalOperatorEquals(t *testing.T) {
	reference := commons.NewFullPeriod()
	other := commons.NewPeriodSince(time.Now(), true)

	if commons.TemporalEquals.Accepts(other, reference) {
		t.Fail()
	} else if !commons.TemporalEquals.Accepts(reference, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorConstants(t *testing.T) {
	reference := commons.NewFullPeriod()
	other := commons.NewPeriodSince(time.Now(), true)

	if commons.TemporalAlwaysRefuse.Accepts(reference, reference) {
		t.Fail()
	} else if !commons.TemporalAlwaysAccept.Accepts(other, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorCommonPoint(t *testing.T) {
	all := commons.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := commons.NewPeriodSince(after, true)
	other := commons.NewPeriodUntil(before, true)

	if commons.TemporalCommonPoint.Accepts(base, other) {
		t.Fail()
	} else if !commons.TemporalCommonPoint.Accepts(all, base) {
		t.Fail()
	}
}

func TestTemporalOperatorInclusion(t *testing.T) {
	all := commons.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := commons.NewPeriodSince(after, true)
	other := commons.NewPeriodUntil(before, true)

	if commons.TemporalReferenceContains.Accepts(base, other) {
		t.Fail()
	} else if !commons.TemporalCommonPoint.Accepts(base, all) {
		t.Fail()
	}
}

func TestEqualsOperator(t *testing.T) {
	if !commons.StringEquals.Accepts("a", "a") {
		t.Fail()
	} else if commons.StringEquals.Accepts("a", "b") {
		t.Fail()
	}
}

func TestEqualsIgnoreCaseOperator(t *testing.T) {
	if !commons.StringEqualsIgnoreCase.Accepts("a", "a") {
		t.Fail()
	} else if !commons.StringEqualsIgnoreCase.Accepts("a", "A") {
		t.Fail()
	} else if commons.StringEqualsIgnoreCase.Accepts("a", "B") {
		t.Fail()
	}
}

func TestContainsOperator(t *testing.T) {
	if commons.StringContains.Accepts("aa", "a") {
		t.Log("expected semantic is that reference contains value, not the other way around")
		t.Fail()
	} else if commons.StringContains.Accepts("a", "b") {
		t.Fail()
	} else if !commons.StringContains.Accepts("a", "aaa") {
		t.Fail()
	}
}

func TestRegexpOperator(t *testing.T) {
	if !commons.StringMatchesRegexp.Accepts("aa", "a+") {
		t.Fail()
	} else if commons.StringMatchesRegexp.Accepts("a", "b+") {
		t.Fail()
	} else if commons.StringMatchesRegexp.Accepts("a", "\\a\\b\\c\\d\\e\\f\\g") {
		t.Log("invalid regexp should fail")
		t.Fail()
	}
}

func TestIntOperator(t *testing.T) {
	if commons.IntEquals.Accepts(0, 10) {
		t.Fail()
	} else if !commons.IntEquals.Accepts(10, 10) {
		t.Fail()
	} else if commons.IntNotEquals.Accepts(20, 20) {
		t.Fail()
	} else if !commons.IntNotEquals.Accepts(0, 20) {
		t.Fail()
	}

	if commons.IntStrictLess.Accepts(30, 30) {
		t.Fail()
	} else if !commons.IntStrictLess.Accepts(0, 30) {
		t.Fail()
	} else if commons.IntLessOrEquals.Accepts(300, 30) {
		t.Fail()
	} else if !commons.IntLessOrEquals.Accepts(0, 30) {
		t.Fail()
	} else if commons.IntStrictGreater.Accepts(30, 30) {
		t.Fail()
	} else if !commons.IntStrictGreater.Accepts(300, 30) {
		t.Fail()
	} else if commons.IntGreaterOrEquals.Accepts(30, 3000) {
		t.Fail()
	} else if !commons.IntGreaterOrEquals.Accepts(300, 30) {
		t.Fail()
	}
}

func TestFloatOperator(t *testing.T) {
	if commons.FloatEquals.Accepts(0.0, 10.0) {
		t.Fail()
	} else if !commons.FloatEquals.Accepts(10.0, 10.0) {
		t.Fail()
	} else if commons.FloatNotEquals.Accepts(20.0, 20.0) {
		t.Fail()
	} else if !commons.FloatNotEquals.Accepts(0.0, 20.0) {
		t.Fail()
	}

	if commons.FloatStrictLess.Accepts(30.0, 30.0) {
		t.Fail()
	} else if !commons.FloatStrictLess.Accepts(0.0, 30.0) {
		t.Fail()
	} else if commons.FloatLessOrEquals.Accepts(300.0, 30.0) {
		t.Fail()
	} else if !commons.FloatLessOrEquals.Accepts(0.0, 30.0) {
		t.Fail()
	} else if commons.FloatStrictGreater.Accepts(30.0, 30.0) {
		t.Fail()
	} else if !commons.FloatStrictGreater.Accepts(300.0, 30.0) {
		t.Fail()
	} else if commons.FloatGreaterOrEquals.Accepts(30.0, 3000.0) {
		t.Fail()
	} else if !commons.FloatGreaterOrEquals.Accepts(300.0, 30.0) {
		t.Fail()
	}
}

func TestLocalSetOperators(t *testing.T) {
	values := []int{0, 10, 100}

	// All values in set should be different from operand
	operator := commons.NewLocalSetOperator(commons.MatchesAllInSetOperator, commons.IntNotEquals)
	if !operator.Accepts(20, values) {
		t.Fail()
	} else if operator.Accepts(10, values) {
		t.Fail()
	}

	// No value should match equals, equivalent to previous
	operator = commons.NewLocalSetOperator(commons.MatchesNoneInSetOperator, commons.IntEquals)
	if !operator.Accepts(20, values) {
		t.Fail()
	} else if operator.Accepts(10, values) {
		t.Fail()
	}

	// one element in set should be equal to value, so it means basically to contain the value
	operator = commons.NewLocalSetOperator(commons.MatchesOneInSetOperator, commons.IntEquals)
	if operator.Accepts(20, values) {
		t.Fail()
	} else if !operator.Accepts(10, values) {
		t.Fail()
	}
}

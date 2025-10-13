package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

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

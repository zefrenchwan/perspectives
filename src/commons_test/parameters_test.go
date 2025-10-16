package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestFormalParametersAccept(t *testing.T) {
	varOnly := commons.NewNamedFormalParameters([]string{"x", "y"})
	positionalOnly := commons.NewPositionalFormalParameters(1)
	yes := commons.NewMostPermissiveFormalParameters()

	if commons.NewNamedContent("x", DummyComponentImplementation{}).Matches(varOnly) {
		t.Fail()
	} else if commons.NewNamedContent("x", DummyComponentImplementation{}).Matches(positionalOnly) {
		t.Fail()
	} else if !commons.NewNamedContent("x", DummyComponentImplementation{}).Matches(yes) {
		t.Fail()
	}

	// one value waiting for one at least
	content := commons.NewContent(DummyComponentImplementation{})
	if !content.Matches(positionalOnly) {
		t.Fail()
	}

	// two values waiting for one at least
	content.Append(DummyComponentImplementation{})
	if !content.Matches(positionalOnly) {
		t.Fail()
	}

	// get necessary variables
	content = commons.NewNamedContent("x", DummyComponentImplementation{})
	content.AppendAs("y", DummyComponentImplementation{})
	if !content.Matches(varOnly) {
		t.Fail()
	}

	// add extra variable
	content.AppendAs("z", DummyComponentImplementation{})
	if !content.Matches(varOnly) {
		t.Fail()
	}
}

func TestMaxParameters(t *testing.T) {
	varOnly := commons.NewNamedFormalParameters([]string{"x", "y"})
	positionalOnly := commons.NewPositionalFormalParameters(1)
	maxParameters := varOnly.Max(positionalOnly)

	content := commons.NewContent(DummyComponentImplementation{})
	if content.Matches(maxParameters) {
		t.Fail()
	}

	content.AppendAs("x", DummyComponentImplementation{})
	if content.Matches(maxParameters) {
		t.Fail()
	}

	content.AppendAs("y", DummyComponentImplementation{})
	if !content.Matches(maxParameters) {
		t.Fail()
	}

	// retest with variables first
	content = commons.NewNamedContent("x", DummyComponentImplementation{})
	if content.Matches(maxParameters) {
		t.Fail()
	}

	content.AppendAs("y", DummyComponentImplementation{})
	if content.Matches(maxParameters) {
		t.Fail()
	}

	content.Append(DummyComponentImplementation{})
	if !content.Matches(maxParameters) {
		t.Fail()
	}
}

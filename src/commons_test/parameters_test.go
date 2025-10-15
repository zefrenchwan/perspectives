package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestFormalParametersAccept(t *testing.T) {
	varOnly := commons.NewNamedFormalParameters([]string{"x", "y"})
	positionalOnly := commons.NewPositionalFormalParameters(1)
	yes := commons.NewMostPermissiveFormalParameters()

	if varOnly.Accepts(commons.NewNamedContent("x", DummyComponentImplementation{})) {
		t.Fail()
	} else if positionalOnly.Accepts(commons.NewNamedContent("x", DummyComponentImplementation{})) {
		t.Fail()
	} else if !yes.Accepts(nil) {
		t.Fail()
	} else if !yes.Accepts(commons.NewNamedContent("x", DummyComponentImplementation{})) {
		t.Fail()
	}

	// one value waiting for one at least
	content := commons.NewContent(DummyComponentImplementation{})
	if !positionalOnly.Accepts(content) {
		t.Fail()
	}

	// two values waiting for one at least
	content.Append(DummyComponentImplementation{})
	if !positionalOnly.Accepts(content) {
		t.Fail()
	}

	// get necessary variables
	content = commons.NewNamedContent("x", DummyComponentImplementation{})
	content.AppendAs("y", DummyComponentImplementation{})
	if !varOnly.Accepts(content) {
		t.Fail()
	}

	// add extra variable
	content.AppendAs("z", DummyComponentImplementation{})
	if !varOnly.Accepts(content) {
		t.Fail()
	}
}

func TestMaxParameters(t *testing.T) {
	varOnly := commons.NewNamedFormalParameters([]string{"x", "y"})
	positionalOnly := commons.NewPositionalFormalParameters(1)
	maxParameters := varOnly.Max(positionalOnly)

	content := commons.NewContent(DummyComponentImplementation{})
	if maxParameters.Accepts(content) {
		t.Fail()
	}

	content.AppendAs("x", DummyComponentImplementation{})
	if maxParameters.Accepts(content) {
		t.Fail()
	}

	content.AppendAs("y", DummyComponentImplementation{})
	if !maxParameters.Accepts(content) {
		t.Fail()
	}

	// retest with variables first
	content = commons.NewNamedContent("x", DummyComponentImplementation{})
	if maxParameters.Accepts(content) {
		t.Fail()
	}

	content.AppendAs("y", DummyComponentImplementation{})
	if maxParameters.Accepts(content) {
		t.Fail()
	}

	content.Append(DummyComponentImplementation{})
	if !maxParameters.Accepts(content) {
		t.Fail()
	}
}

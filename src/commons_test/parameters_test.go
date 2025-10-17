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

func TestContentMappers(t *testing.T) {
	first := DummyComponentImplementation{}
	second := DummyIdBasedImplementation{}

	variables := commons.NewNamedContent[commons.Modelable]("x", first)
	variables.AppendAs("y", second)

	if _, ok := variables.MapNamedToPositionals([]string{"x", "y", "z"}); ok {
		t.Fail()
	} else if _, ok := variables.MapNamedToPositionals([]string{"a"}); ok {
		t.Fail()
	} else if _, ok := variables.MapPositionalsToNamed([]string{"x"}); ok {
		t.Fail()
	} else if value, ok := variables.MapNamedToPositionals([]string{"y", "x"}); !ok {
		t.Fail()
	} else if content := value.PositionalContent(); len(content) != 2 {
		t.Fail()
	} else if content[0] != second {
		t.Fail()
	} else if content[1] != first {
		t.Fail()
	}

	positionals := commons.NewContent[commons.Modelable](first)
	positionals.Append(second)
	if _, ok := positionals.MapPositionalsToNamed([]string{"x", "y", "z"}); ok {
		t.Fail()
	} else if value, ok := positionals.MapPositionalsToNamed([]string{"y", "x"}); !ok {
		t.Fail()
	} else if content := value.NamedContent(); len(content) != 2 {
		t.Fail()
	} else if content["x"] != second {
		t.Fail()
	} else if content["y"] != first {
		t.Fail()
	}
}

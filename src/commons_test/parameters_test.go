package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestPermissiveParameters(t *testing.T) {
	fp := commons.NewMostPermissiveFormalParameters()

	if !fp.Accepts(nil) {
		t.Fail()
	}

	content := commons.NewNamedContent("x", DummyComponentImplementation{})
	if !fp.Accepts(content) {
		t.Fail()
	}

	content.Append(DummyComponentImplementation{})
	if !fp.Accepts(content) {
		t.Fail()
	}
}

func TestVariableFormalParameters(t *testing.T) {

	fp := commons.NewVariablesFormalParameters(nil)
	if !fp.Accepts(nil) {
		t.Fail()
	}

	dummy := commons.NewContent(DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

	fp = commons.NewVariablesFormalParameters([]string{"x", "y"})

	if fp.Accepts(nil) {
		t.Fail()
	}

	content := commons.NewNamedContent("x", DummyComponentImplementation{})
	if fp.Accepts(content) {
		t.Fail()
	}

	content.AppendAsVariable("y", DummyComponentImplementation{})
	if !fp.Accepts(content) {
		t.Fail()
	}

	content.AppendAsVariable("z", DummyComponentImplementation{})
	if !fp.Accepts(content) {
		t.Fail()
	}
}

func TestMinimalSizeFormalParameters(t *testing.T) {
	fp := commons.NewPositionalFormalParameters(0)

	if !fp.Accepts(nil) {
		t.Fail()
	}

	dummy := commons.NewContent(DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

	fp = commons.NewPositionalFormalParameters(2)
	// lower values count
	if fp.Accepts(dummy) {
		t.Fail()
	}

	// same size exactly
	dummy.Append(DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

	// more values
	dummy.Append(DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

}

func TestUniqueFormalParameters(t *testing.T) {

	fp := commons.NewUniqueFormalParameters()
	if fp.Accepts(nil) {
		t.Fail()
	}

	// test one named
	dummy := commons.NewNamedContent("x", DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

	// test 2 > 1
	dummy.Append(DummyComponentImplementation{})
	if fp.Accepts(dummy) {
		t.Fail()
	}

	// test one positional
	dummy = commons.NewContent(DummyComponentImplementation{})
	if !fp.Accepts(dummy) {
		t.Fail()
	}

	// test 2 > 1
	dummy.Append(DummyComponentImplementation{})
	if fp.Accepts(dummy) {
		t.Fail()
	}
}

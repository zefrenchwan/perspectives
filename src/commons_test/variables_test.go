package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestLinkVariableReplace(t *testing.T) {
	variable := commons.NewVariable("X", commons.CLASS_TRAIT)
	trait := commons.NewTrait("Trait")
	instance := commons.NewTemporalInstance()

	if !variable.CanBeReplacedBy(trait) {
		t.Errorf("Variable cannot be replaced by trait")
		return
	} else if variable.CanBeReplacedBy(instance) {
		t.Errorf("Variable cannot be replaced by instance")
		return
	}

	rule := commons.NewLink("equals", commons.NewFullPeriod()).
		WithOperand("subject", []commons.Element{variable}).
		WithOperand("object", []commons.Element{variable})
	link := rule.ReplaceVariable(variable, trait)
	if link == nil {
		t.Errorf("Expected non-nil link after replacement, got nil")
	}

	if values, found := link.Operand("subject"); !found {
		t.Errorf("Expected 'subject' operand to be present after replacement")
	} else if len(values) != 1 {
		t.Errorf("Expected 'subject' operand to have one value, got %d", len(values))
	} else if value := values[0]; !value.Same(trait) {
		t.Errorf("Expected 'subject' operand to be replaced with trait, got %v", value)
	}

	if values, found := link.Operand("object"); !found {
		t.Errorf("Expected 'object' operand to be present after replacement")
	} else if len(values) != 1 {
		t.Errorf("Expected 'objectt' operand to have one value, got %d", len(values))
	} else if value := values[0]; !value.Same(trait) {
		t.Errorf("Expected 'object' operand to be replaced with trait, got %v", value)
	}
}

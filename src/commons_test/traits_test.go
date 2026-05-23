package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestTraits(t *testing.T) {
	trait := commons.NewTrait("person")
	if trait.Name() != "person" {
		t.Errorf("Expected trait name to be 'person', got '%s'", trait.Name())
	} else if trait.Id() == "" {
		t.Errorf("Expected trait ID to be non-empty, got empty string")
	} else if !commons.IsElementDeclaredInstance(trait, commons.CLASS_TRAIT) {
		t.Errorf("Expected trait to be declared as CLASS_TRAIT, got undeclared")
	}
}

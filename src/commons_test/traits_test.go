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

	trait.WithAttribute("age", "int")
	if trait.Attributes()["age"] != "int" {
		t.Errorf("Expected attribute 'age' to have type 'int', got '%s'", trait.Attributes()["age"])
	}

	trait.RemoveAttribute("age")
	if _, ok := trait.Attributes()["age"]; ok {
		t.Errorf("Expected attribute 'age' to be removed, but it still exists")
	}

	if len(trait.Attributes()) != 0 {
		t.Errorf("Expected no attributes after removal, but found %d attributes", len(trait.Attributes()))
	}
}

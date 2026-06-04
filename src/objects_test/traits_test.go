package objects_test_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestTraits(t *testing.T) {
	trait := objects.NewTrait("person")
	if trait.Name() != "person" {
		t.Errorf("Expected trait name to be 'person', got '%s'", trait.Name())
	} else if !objects.IsElementDeclaredInstance(trait, objects.CLASS_TRAIT) {
		t.Errorf("Expected trait to be declared as CLASS_TRAIT, got undeclared")
	}

	trait = trait.WithAttribute("age", "int")
	if trait.Attributes()["age"] != "int" {
		t.Errorf("Expected attribute 'age' to have type 'int', got '%s'", trait.Attributes()["age"])
	}

	trait = trait.WithoutAttribute("age")
	if _, ok := trait.Attributes()["age"]; ok {
		t.Errorf("Expected attribute 'age' to be removed, but it still exists")
	}

	if len(trait.Attributes()) != 0 {
		t.Errorf("Expected no attributes after removal, but found %d attributes", len(trait.Attributes()))
	}
}

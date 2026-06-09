package objects_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestTraits(t *testing.T) {
	trait := objects.NewTrait("person")
	if trait.Name() != "person" {
		t.Errorf("Expected trait name to be 'person', got '%s'", trait.Name())
	} else if !objects.IsInstanceOfClass(trait, objects.CLASS_TRAIT) {
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

func TestTraitSame(t *testing.T) {
	t1 := objects.NewTrait("person")
	t2 := objects.NewTrait("person")
	t3 := objects.NewTrait("animal")

	if !t1.Same(t2) {
		t.Errorf("Expected trait %v and %v to be same", t1, t2)
	}

	if t1.Same(t3) {
		t.Errorf("Expected trait %v and %v to be different", t1, t3)
	}

	if t1.Same(nil) {
		t.Errorf("Expected trait %v and nil to be different", t1)
	}

	t4 := objects.NewTrait("person").WithAttribute("age", "int")
	t5 := objects.NewTrait("person").WithAttribute("age", "int")
	t6 := objects.NewTrait("person").WithAttribute("age", "string")

	if !t4.Same(t5) {
		t.Errorf("Expected traits with same name and attributes to be same")
	}

	if t4.Same(t6) {
		t.Errorf("Expected traits with same name but different attributes to be different")
	}

	if t1.Same(t4) {
		t.Errorf("Expected trait without attribute and trait with attribute to be different")
	}
}

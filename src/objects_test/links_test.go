package objects_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestLinkSame(t *testing.T) {
	instance, _ := objects.NewLocalInstanceBuilder("id").WithActivity(periods.NewFullPeriod()).Build()
	other, _ := objects.NewLocalInstanceBuilder("other").WithActivity(periods.NewFullPeriod()).Build()
	reference, _ := objects.NewLocalLinkBuilder("links:is").
		WithName("is").
		WithOperand("subject", instance).
		WithOperand("object", other).
		Build()

	noMatchName, _ := objects.NewLocalLinkBuilder("links:other").
		WithName("is not").
		WithOperand("subject", instance).
		WithOperand("object", other).
		Build()

	noMatchOperand, _ := objects.NewLocalLinkBuilder("links:is").
		WithName("is").
		WithOperand("subject", other).
		WithOperand("object", instance).
		Build()

	if !reference.Same(reference) {
		t.Errorf("Expected link to be the same as itself")
	} else if reference.Same(noMatchName) {
		t.Errorf("Expected link to be different from another link with different name")
	} else if reference.Same(noMatchOperand) {
		t.Errorf("Expected link to be different from another link with different operand")
	}
}

func TestBuildLinks(t *testing.T) {
	instance, _ := objects.NewLocalInstanceBuilder("id").WithActivity(periods.NewFullPeriod()).Build()
	one, _ := objects.NewLocalInstanceBuilder("one").Build()
	other, _ := objects.NewLocalInstanceBuilder("other").Build()

	is, errIs := objects.NewLocalLinkBuilder("links:is").
		WithName("is").
		WithActivity(periods.NewFullPeriod()).
		WithOperand("subject", instance).
		WithOperand("object", one).
		WithOperand("other", other).
		WithoutOperand("other").
		Build()

	if errIs != nil {
		t.Errorf("Error building link: %v", errIs)
	} else if is.Name() != "is" {
		t.Errorf("Expected link name to be 'is', got '%s'", is.Name())
	} else if is.Id() != "links:is" {
		t.Errorf("Expected link ID to be 'links:is', got '%s'", is.Id())
	} else if subject, hasSubject := is.Role("subject"); !hasSubject {
		t.Errorf("Expected link to have a 'subject' role")
	} else if subject != instance {
		t.Errorf("Expected link subject to be '%s', got '%v'", instance.Id(), subject)
	} else if object, hasObject := is.Role("object"); !hasObject {
		t.Errorf("Expected link to have an 'object' role")
	} else if object != one {
		t.Errorf("Expected link object to be '%s', got '%v'", one.Id(), object)
	} else if !is.Activity().Equals(periods.NewFullPeriod()) {
		t.Errorf("Expected link activity to be full period, got '%v'", is.Activity())
	}
}

func TestBuildLinksErrors(t *testing.T) {
	instance, _ := objects.NewLocalInstanceBuilder("id").WithActivity(periods.NewFullPeriod()).Build()

	var builder objects.LinkBuilder
	builder = objects.NewLocalLinkBuilder("links:is")
	if _, err := builder.Build(); err == nil {
		t.Errorf("no name, no operand, should fail")
	}

	builder = objects.NewLocalLinkBuilder("links:is").WithName("is")
	if _, err := builder.Build(); err == nil {
		t.Errorf("no operand, should fail")
	}

	builder = objects.NewLocalLinkBuilder("links:is").WithName("is").WithOperand("no op", nil)
	if _, err := builder.Build(); err == nil {
		t.Errorf("nil, should fail")
	}

	builder = objects.NewLocalLinkBuilder("links:is").WithName("is").WithOperand("", instance)
	if _, err := builder.Build(); err == nil {
		t.Errorf("empty role, should fail")
	}
}

func TestReBuildLinks(t *testing.T) {
	instance, _ := objects.NewLocalInstanceBuilder("id").WithActivity(periods.NewFullPeriod()).Build()
	animals, _ := objects.NewLocalInstanceBuilder("Animals").Build()
	is, errIs := objects.NewLocalLinkBuilder("links:is").
		WithName("is").
		WithOperand("subject", instance).
		WithOperand("object", animals).
		Build()

	if errIs != nil {
		t.Errorf("Error building link: %v", errIs)
	}

	if isCopy, err := objects.LocalLinkBuilderLoad(is).Build(); err != nil {
		t.Errorf("Error rebuilding link: %v", err)
	} else if isCopy.Id() != is.Id() {
		t.Errorf("Expected link ID to be '%s', got '%s'", is.Id(), isCopy.Id())
	} else if !is.Same(isCopy) {
		t.Errorf("Expected link to be the same as the copy")
	}
}

func TestLinkRange(t *testing.T) {
	instance, _ := objects.NewLocalInstanceBuilder("id").WithActivity(periods.NewFullPeriod()).Build()
	animals, _ := objects.NewLocalInstanceBuilder("Animals").Build()
	is, _ := objects.NewLocalLinkBuilder("links:is").
		WithName("is").
		WithOperand("subject", instance).
		WithOperand("object", animals).
		Build()

	if roles := is.Roles(); len(roles) != 2 {
		t.Errorf("Expected link to have 2 roles, got %d", len(roles))
	} else if slices.Compare(roles, []string{"object", "subject"}) != 0 {
		t.Errorf("Expected link roles to be ['object', 'subject'], got %v", roles)
	}

	for name, value := range is.Range {
		switch name {
		case "subject":
			if value != instance {
				t.Errorf("Expected subject to be '%s', got '%v'", instance.Id(), value)
			}
		case "object":
			if value != animals {
				t.Errorf("Expected object to be '%s', got '%v'", animals.Id(), value)
			}
		default:
			t.Errorf("Unexpected role '%s'", name)
		}
	}
}

package objects

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestVariableBuild(t *testing.T) {
	if _, err := objects.NewVariableBuilder("").WithAcceptedClass(objects.CLASS_VARIABLE).Build(); err == nil {
		t.Errorf("Expected variable with no name to fail")
	}

	if _,err := objects.NewVariableBuilder("x").Build(); err == nil {
		t.Errorf("Expected variable with no accepted classes to fail")
	}

	if v, err := objects.NewVariableBuilder("x").WithAcceptedClass(objects.CLASS_VARIABLE).Build(); err != nil {
		t.Errorf("Expected variable with accepted class to succeed")
	} else if v.Name() != "x" {
		t.Errorf("Expected variable name to be 'x', got '%s'", v.Name())
	} else if len(v.Accepts()) != 1 || v.Accepts()[0] != objects.CLASS_VARIABLE {
		t.Errorf("Expected variable to accept only VARIABLE class, got %v", v.Accepts())
	}
}

func TestVariableAccepts(t *testing.T) {
	variable, _ := objects.NewVariableBuilder("x").WithAcceptedClass(objects.CLASS_LINK).Build()
	if variable.Name() != "x" {
		t.Errorf("Expected variable name to be 'x', got '%s'", variable.Name())
	} else if slices.Compare(variable.Accepts(), []objects.Class{objects.CLASS_LINK}) != 0 {
		t.Errorf("Expected variable to accept only LINK class, got %v", variable.Accepts())
	} else if variable.AcceptsOneOf(objects.CLASS_VARIABLE, objects.CLASS_INSTANCE) {
		t.Errorf("Expected variable to accept neither VARIABLE nor INSTANCE class, got %v", variable.Accepts())
	} else if !variable.AcceptsOneOf(objects.CLASS_LINK, objects.CLASS_INSTANCE) {
		t.Errorf("Expected variable to accept either LINK or INSTANCE class, got %v", variable.Accepts())
	}
}

func TestVariableSame(t *testing.T) {
	variable, _ := objects.NewVariableBuilder("x").WithAcceptedClass(objects.CLASS_LINK).Build()
	if !variable.Same(variable) {
		t.Errorf("Expected variable to be the same as itself")
	}

	noMatchName, _ := objects.NewVariableBuilder("y").WithAcceptedClass(objects.CLASS_LINK).Build()
	if variable.Same(noMatchName) {
		t.Errorf("Expected variable to not be the same as another variable with a different name")
	}

	noMatchClass, _ := objects.NewVariableBuilder("x").WithAcceptedClass(objects.CLASS_VARIABLE).Build()
	if variable.Same(noMatchClass) {
		t.Errorf("Expected variable to not be the same as another variable with a different class")
	}

	noMatchMoreClasses, _ := objects.NewVariableBuilder("x").
		WithAcceptedClass(objects.CLASS_LINK).
		WithAcceptedClass(objects.CLASS_INSTANCE).
		Build()
	if variable.Same(noMatchMoreClasses) {
		t.Errorf("Expected variable to not be the same as another variable with more classes")
	}
}

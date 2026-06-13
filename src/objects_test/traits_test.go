package objects_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestTraitBuilder(t *testing.T) {
	if trait, err := objects.NewTraitBuilder().WithName("Animals").Build(); err != nil {
		t.Errorf("trait with just name should be OK")
	} else if trait == nil {
		t.Errorf("trait with just name should be OK")
	} else if trait.Name() != "Animals" {
		t.Errorf("error when passing name")
	}

	if trait, err := objects.NewTraitBuilder().WithName("Humans").
		WithAttribute("age", "int").
		Build(); err != nil {
		t.Errorf("trait with name and attribute should be OK")
	} else if trait == nil {
		t.Errorf("trait with name and attribute should be OK")
	} else if trait.Name() != "Humans" {
		t.Errorf("error when passing name")
	} else if attrs := trait.Attributes(); len(attrs) != 1 {
		t.Errorf("error when passing attribute")
	} else if attrs["age"] != "int" {
		t.Errorf("error when reading attribute type")
	}

	if trait, err := objects.NewTraitBuilder().WithName("Humans").
		WithAttribute("age", "string").
		WithAttribute("age", "int").
		WithAttribute("test", "string").
		WithoutAttribute("test").
		Build(); err != nil {
		t.Errorf("trait with name and attribute should be OK")
	} else if trait == nil {
		t.Errorf("trait with name and attribute should be OK")
	} else if trait.Name() != "Humans" {
		t.Errorf("error when passing name")
	} else if attrs := trait.Attributes(); len(attrs) != 1 {
		t.Errorf("error when passing attribute")
	} else if attrs["age"] != "int" {
		t.Errorf("error when reading attribute type")
	}
}

func TestTraitBuilderErrors(t *testing.T) {
	if _, err := objects.NewTraitBuilder().WithName("").Build(); err == nil {
		t.Errorf("error expected when passing empty name")
	} else if _, err := objects.NewTraitBuilder().WithAttribute("", "int").Build(); err == nil {
		t.Errorf("error expected when passing attribute without name")
	}

	if _, err := objects.NewTraitBuilder().
		WithName("Humans").
		WithAttribute("age", "").
		Build(); err == nil {
		t.Errorf("error expected when passing empty attribute type")
	} else if _, err := objects.NewTraitBuilder().
		WithName("Humans").
		WithAttribute("", "int").
		Build(); err == nil {
		t.Errorf("error expected when passing empty attribute name")
	}
}

func TestTraitSame(t *testing.T) {
	trait, _ := objects.NewTraitBuilder().WithName("Humans").
		WithAttribute("age", "int").
		Build()

	if mismatchName, _ := objects.NewTraitBuilder().WithName("Dogs").
		WithAttribute("age", "int").
		Build(); trait.Same(mismatchName) {
		t.Errorf("names mismatch")
	}

	if mismatchAttr, _ := objects.NewTraitBuilder().WithName("Humans").
		WithAttribute("age", "string").
		Build(); trait.Same(mismatchAttr) {
		t.Errorf("attributes mismatch")
	}

	if mismatchAttr, _ := objects.NewTraitBuilder().
		WithName("Humans").
		WithAttribute("age", "int").
		WithAttribute("name", "string").
		Build(); trait.Same(mismatchAttr) {
		t.Errorf("attributes count mismatch")
	}
}

package objects_test

import (
	"maps"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestTraitBuildFromScratch(t *testing.T) {
	builder := objects.NewTraitBuilder()
	if trait, err := builder.WithName("test").
		WithAttribute("attr", "string").
		WithAttribute("attr2", "int").
		WithoutAttribute("attr2").
		Build(); err != nil {
		t.Error("failed to build basic trait", err)
	} else if trait.Name() != "test" {
		t.Error("trait name is not test")
	} else if trait.DeclaringClass() != objects.CLASS_DEFINITION {
		t.Error("trait declaring class is not class definition")
	} else if attrs := maps.Collect(trait.Attributes); len(attrs) != 1 {
		t.Error("missing trait attributes")
	} else if attrs["attr"] != "string" {
		t.Error("missing trait attribute name")
	}
}

func TestTraitBuildErrors(t *testing.T) {
	builder := objects.NewTraitBuilder()
	if _, err := builder.Build(); err == nil {
		t.Error("expected error when building trait without name")
	} else if _, err := builder.WithName("     ").Build(); err == nil {
		t.Error("expected error when building trait with spaces only")
	} else if _, err := builder.WithName("").Build(); err == nil {
		t.Error("expected error when building trait with empty name")
	}

	builder = objects.NewTraitBuilder().WithName("test")
	if builder.WithAttribute("", "string").Errors() == nil {
		t.Error("expected errors when building trait with empty attribute")
	}

	builder = objects.NewTraitBuilder().WithName("test")
	if builder.WithAttribute("attr", "not a good name").Errors() == nil {
		t.Error("expected errors when building trait with a bad type for attribute")
	}
}

func TestTraitMatchingAccepts(t *testing.T) {
	birthDate := time.Now().Truncate(24*time.Hour).AddDate(-20, 0, 0)
	studentDate := birthDate.AddDate(18, 0, 0)
	studentPeriod := periods.NewPeriodSince(studentDate, true)
	student, _ := objects.NewLocalInstanceBuilder("john").
		WithActivity(periods.NewPeriodSince(birthDate, true)).
		WithAttributeDuring("student id", studentPeriod, 178).
		Build()

	if studentTrait, err := objects.NewTraitBuilder().
		WithName("student").
		WithAttribute("student id", "int").
		Build(); err != nil {
		t.Error("failed to build student trait")
	} else if matchingPeriod, matching := studentTrait.Matches(student); !matching || !matchingPeriod.Equals(studentPeriod) {
		t.Error("student trait does not match student definition")
	}
}

func TestTraitSame(t *testing.T) {
	// Matching traits
	t1, _ := objects.NewTraitBuilder().
		WithName("student").
		WithAttribute("id", "int").
		Build()
	t2, _ := objects.NewTraitBuilder().
		WithName("student").
		WithAttribute("id", "int").
		Build()

	if !t1.Same(t2) {
		t.Error("identical traits should be the same")
	}

	// Not a trait (should fail)
	instance, _ := objects.NewLocalInstanceBuilder("john").Build()
	if t1.Same(instance) {
		t.Error("trait should not be the same as an instance")
	}

	// Case for nil
	if t1.Same(nil) {
		t.Error("trait should not be the same as nil")
	}

	// Different name, same attributes
	t3, _ := objects.NewTraitBuilder().
		WithName("pupil").
		WithAttribute("id", "int").
		Build()
	if t1.Same(t3) {
		t.Error("traits with different names should not be the same")
	}

	// Same name, different attributes
	t4, _ := objects.NewTraitBuilder().
		WithName("student").
		WithAttribute("id", "string").
		Build()
	if t1.Same(t4) {
		t.Error("traits with different attributes should not be the same")
	}
}

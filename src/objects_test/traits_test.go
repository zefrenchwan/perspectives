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
	if trait, err := builder.WithName("test").WithAttribute("attr", "string").Build(); err != nil {
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

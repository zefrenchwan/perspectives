package objects_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestBuildFromScratch(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if content, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if content == nil {
		t.Error("content should not be nil")
	} else if values, exists := content.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := content.Description(); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}

func TestBuildFromOther(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if content, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if content == nil {
		t.Error("content should not be nil")
	} else if other, errOther := objects.LocalContentBuilderLoad(content).Build(); errOther != nil {
		t.Error(errOther)
	} else if other == nil {
		t.Error("other should not be nil")
	} else if values, exists := other.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := other.Description(); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}

func TestBuildError(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if _, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		WithAttributeDuring("name", periods.NewFullPeriod(), 10).
		Build(); err == nil {
		t.Error("expected error for invalid attribute that changed its type")
	}
}

func TestContentAt(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	after := now.Add(time.Hour * 24)
	timmy, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(now, true)).
		WithAttributeDuring("name", periods.NewPeriodSince(now, true), "Timmy").
		WithAttributeDuring("age", periods.NewPeriodSince(after, true), 25).
		Build()
	if err != nil {
		t.Error(err)
	} else if timmy == nil {
		t.Error("expected content to be non-nil after creation")
	}

	if _, beforeFound := timmy.At(before); beforeFound {
		t.Error("expected content NOT to exist (was created now)")
	}

	if values, afterFound := timmy.At(after); !afterFound {
		t.Error("expected content to exist at after (was created now)")
	} else if values == nil {
		t.Error("expected content values to be non-nil at after (was created now)")
	} else if len(values) != 2 {
		t.Error("expected 2 values at after")
	} else if values["name"] != "Timmy" {
		t.Error("expected name value to be 'Timmy' at after")
	} else if values["age"] != 25 {
		t.Error("expected age value to be 25 at after")
	}

	if value, found := timmy.Value("name"); !found {
		t.Error("expected 'name' value to exist at after")
	} else if _, hasBefore := value.At(before); hasBefore {
		t.Error("no value expected at before")
	} else if vnow, hasNow := value.At(now); !hasNow {
		t.Error("value expected at now")
	} else if vnow != "Timmy" {
		t.Error("expected name value to be 'Timmy' at now")
	}

	if value, found := timmy.Value("age"); !found {
		t.Error("expected 'age' value to exist at after")
	} else if _, hasBefore := value.At(before); hasBefore {
		t.Error("no value expected at before")
	} else if _, hasNow := value.At(now); hasNow {
		t.Error("no value expected at now")
	}
}

func TestContentCut(t *testing.T) {
	now := time.Now()
	before := time.Now().AddDate(-25, 0, 0)
	lifetime := periods.NewPeriodSince(before, true)
	lara, err := objects.NewLocalContentBuilder().
		WithActivity(lifetime).
		WithAttributeDuring("name", periods.NewFullPeriod(), "Lara").
		WithAttributeDuring("age", periods.NewFullPeriod(), 25).
		Build()
	if err != nil {
		t.Error(err)
	}

	if cutLara, err := objects.LocalContentBuilderLoad(lara).Cut(lara.Activity()).Build(); err != nil {
		t.Error(err)
	} else if cutLara == nil {
		t.Error("cutLara should not be nil")
	} else if nameValue, has := cutLara.Value("name"); !has {
		t.Error("name attribute should exist in cutLara")
	} else if nameValue.DataType() != "string" {
		t.Errorf("expected 'string', got '%s'", nameValue.DataType())
	} else if !nameValue.Validity().Equals(lifetime) {
		t.Errorf("expected attribute periode reduction, got '%v'", nameValue.Validity())
	} else if ageValue, has := cutLara.Value("age"); !has {
		t.Error("age attribute should exist in cutLara")
	} else if ageValue.DataType() != "int" {
		t.Errorf("expected 'int', got '%s'", ageValue.DataType())
	} else if !ageValue.Validity().Equals(lifetime) {
		t.Errorf("expected attribute periode reduction, got '%v'", ageValue.Validity())
	} else if value, foundAge := ageValue.At(now); !foundAge || value != 25 {
		t.Errorf("expected age to be 25, got '%v'", value)
	}
}

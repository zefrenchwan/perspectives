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

// =========================================================================
// PRIMITIVE TYPES TESTS
// =========================================================================

func TestIsPrimitiveValue(t *testing.T) {
	// Valid primitives
	if !objects.IsPrimitiveValue(42) {
		t.Error("expected int to be a valid primitive value")
	} else if !objects.IsPrimitiveValue(3.14) {
		t.Error("expected float64 to be a valid primitive value")
	} else if !objects.IsPrimitiveValue("hello") {
		t.Error("expected string to be a valid primitive value")
	} else if !objects.IsPrimitiveValue(true) {
		t.Error("expected bool to be a valid primitive value")
	}

	// Invalid primitives
	if objects.IsPrimitiveValue([]int{1, 2}) {
		t.Error("expected slice NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(map[string]int{"a": 1}) {
		t.Error("expected map NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(nil) {
		t.Error("expected nil NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(struct{ name string }{name: "test"}) {
		t.Error("expected struct NOT to be a primitive value")
	}
}

// =========================================================================
// BUILDER EDGE CASES TESTS
// =========================================================================

func TestBuilderErrorAccumulation(t *testing.T) {
	lifetime := periods.NewFullPeriod()
	builder := objects.NewLocalContentBuilder().WithActivity(lifetime)

	// Inducing multiple errors
	builder.WithAttributeDuring("age", lifetime, nil)               // Error 1: nil value
	builder.WithAttributeDuring("name", lifetime, []string{"John"}) // Error 2: non-primitive

	if err := builder.Errors(); err == nil {
		t.Error("expected errors to be accumulated, got nil")
	} else if content, buildErr := builder.Build(); buildErr == nil {
		t.Error("expected build to fail and return accumulated errors")
	} else if content != nil && !content.Activity().IsEmpty() {
		t.Error("expected an empty content on error build")
	}
}

func TestBuilderWithoutAttributeDuring(t *testing.T) {
	now := time.Now()
	lifetime := periods.NewPeriodSince(now, true)

	builder := objects.NewLocalContentBuilder().
		WithActivity(lifetime).
		WithAttributeDuring("status", lifetime, "active")

	// 1. Remove on an empty period (should do nothing)
	builder.WithoutAttributeDuring("status", periods.NewEmptyPeriod())

	content, err := builder.Build()
	if err != nil {
		t.Error(err)
	} else if _, has := content.Value("status"); !has {
		t.Error("expected 'status' to still exist after removing an empty period")
	}

	// 2. Remove entirely
	// Now content is well-defined in the function scope
	builder = objects.LocalContentBuilderLoad(content)
	builder.WithoutAttributeDuring("status", lifetime)

	if content2, err := builder.Build(); err != nil {
		t.Error(err)
	} else if _, has := content2.Value("status"); has {
		t.Error("expected 'status' attribute to be completely removed")
	}
}

// =========================================================================
// TEMPORAL LOGIC & VALUES HANDLER (TESTED VIA BUILDER)
// =========================================================================

func TestValueOverlapsAndScissions(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 48)
	after := now.Add(time.Hour * 48)

	fullPeriod := periods.NewFinitePeriod(before, after, true, true)
	midPeriod := periods.NewFinitePeriod(now.Add(-time.Hour), now.Add(time.Hour), true, true)

	// We add a value on the full period, then overwrite the middle with a DIFFERENT value.
	// This should split the original period into two distinct nodes.
	content, err := objects.NewLocalContentBuilder().
		WithActivity(fullPeriod).
		WithAttributeDuring("status", fullPeriod, "offline").
		WithAttributeDuring("status", midPeriod, "online").
		Build()

	if err != nil {
		t.Error(err)
	} else if statusValue, found := content.Value("status"); !found {
		t.Error("expected 'status' value to exist")
	} else {
		// Check middle period (should have been overwritten)
		if val, ok := statusValue.At(now); !ok {
			t.Error("expected a value at the middle period")
		} else if val != "online" {
			t.Errorf("expected 'online' in mid period, got '%v'", val)
		}

		// Check before middle period (should still be offline)
		if val, ok := statusValue.At(before.Add(time.Hour)); !ok {
			t.Error("expected a value before the middle period")
		} else if val != "offline" {
			t.Errorf("expected 'offline' before mid period, got '%v'", val)
		}

		// Check after middle period (should still be offline)
		if val, ok := statusValue.At(after.Add(-time.Hour)); !ok {
			t.Error("expected a value after the middle period")
		} else if val != "offline" {
			t.Errorf("expected 'offline' after mid period, got '%v'", val)
		}
	}
}

// =========================================================================
// CONTENT METHODS (MATCHES & SAME)
// =========================================================================

func TestContentMatches(t *testing.T) {
	now := time.Now()
	activePeriod := periods.NewPeriodSince(now, true)

	content, err := objects.NewLocalContentBuilder().
		WithActivity(activePeriod).
		WithAttributeDuring("role", activePeriod, "admin").
		WithAttributeDuring("level", activePeriod, 5).
		Build()
	if err != nil {
		t.Error(err)
	}

	// 1. Perfect match
	validTrait := objects.NewTrait("valid").
		WithAttribute("role", "string").
		WithAttribute("level", "int")
	if matchPeriod, matches := content.Matches(validTrait); !matches {
		t.Error("expected content to match valid trait")
	} else if !matchPeriod.Equals(activePeriod) {
		t.Error("expected matching period to equal the content's active period")
	}

	// 2. Wrong type
	invalidTypeTrait := objects.NewTrait("int role").WithAttribute("role", "int")
	// Role is string in content

	if _, matches := content.Matches(invalidTypeTrait); matches {
		t.Error("expected content NOT to match trait due to incorrect type")
	}

	// 3. Missing attribute
	missingAttrTrait := objects.NewTrait("missing fields").WithAttribute("unknown_field", "string")
	if _, matches := content.Matches(missingAttrTrait); matches {
		t.Error("expected content NOT to match trait due to missing attribute")
	}
}

func TestContentSame(t *testing.T) {
	p := periods.NewFullPeriod()

	c1, _ := objects.NewLocalContentBuilder().
		WithActivity(p).
		WithAttributeDuring("key", p, "value").
		Build()

	c2, _ := objects.NewLocalContentBuilder().
		WithActivity(p).
		WithAttributeDuring("key", p, "value").
		Build()

	c3, _ := objects.NewLocalContentBuilder().
		WithActivity(p).
		WithAttributeDuring("key", p, "different_value").
		Build()

	c4, _ := objects.NewLocalContentBuilder().
		WithActivity(periods.NewEmptyPeriod()). // Different activity
		WithAttributeDuring("key", p, "value").
		Build()

	if !c1.Same(c2) {
		t.Error("expected identical contents to be evaluated as same")
	} else if c1.Same(c3) {
		t.Error("expected contents with different values NOT to be evaluated as same")
	} else if c1.Same(c4) {
		t.Error("expected contents with different activities NOT to be evaluated as same")
	}
}

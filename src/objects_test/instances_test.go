package objects_test

import (
	"maps"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// ========================================================================
// INSTANCES BASIC BEHAVIOR TESTS
// ========================================================================

func TestBuildFromScratch(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if instance, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if instance == nil {
		t.Error("content should not be nil")
	} else if instance.Id() != "id" {
		t.Errorf("expected 'id', got '%s'", instance.Id())
	} else if values, exists := instance.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := maps.Collect(instance.Description); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}

func TestBuildFromOther(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if content, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if content == nil {
		t.Error("content should not be nil")
	} else if other, errOther := objects.LocalInstanceBuilderLoad(content).Build(); errOther != nil {
		t.Error(errOther)
	} else if other == nil {
		t.Error("other should not be nil")
	} else if other.Id() != "id" {
		t.Errorf("expected 'id', got '%s'", other.Id())
	} else if values, exists := other.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := maps.Collect(other.Description); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}

func TestBuildError(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if _, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		WithAttributeDuring("name", periods.NewFullPeriod(), 10).
		Build(); err == nil {
		t.Error("expected error for invalid attribute that changed its type")
	}
}

func TestInstanceAt(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	after := now.Add(time.Hour * 24)
	timmy, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(periods.NewPeriodSince(now, true)).
		WithAttributeDuring("name", periods.NewPeriodSince(now, true), "Timmy").
		WithAttributeDuring("age", periods.NewPeriodSince(after, true), 25).
		Build()
	if err != nil {
		t.Error(err)
	} else if timmy == nil {
		t.Error("expected instance to be non-nil after creation")
	}

	if _, beforeFound := timmy.At(before); beforeFound {
		t.Error("expected instance NOT to exist (was created now)")
	}

	if values, afterFound := timmy.At(after); !afterFound {
		t.Error("expected instance to exist at after (was created now)")
	} else if values == nil {
		t.Error("expected instance values to be non-nil at after (was created now)")
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

func TestInstanceCut(t *testing.T) {
	now := time.Now()
	before := time.Now().AddDate(-25, 0, 0)
	lifetime := periods.NewPeriodSince(before, true)
	if lara, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(lifetime).
		WithAttributeDuring("name", periods.NewFullPeriod(), "Lara").
		WithAttributeDuring("age", periods.NewFullPeriod(), 25).
		Build(); err != nil {
		t.Error(err)
	} else if cutLara, err := objects.LocalInstanceBuilderLoad(lara).Cut(lara.Activity()).Build(); err != nil {
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
// BUILDER EDGE CASES TESTS
// =========================================================================

func TestBuilderNoId(t *testing.T) {
	builder := objects.NewLocalInstanceBuilder("").WithActivity(periods.NewFullPeriod())

	if _, err := builder.Build(); err == nil {
		t.Errorf("cannot make instance with no id")
	}
}

func TestBuilderErrorAccumulation(t *testing.T) {
	lifetime := periods.NewFullPeriod()
	builder := objects.NewLocalInstanceBuilder("id").WithActivity(lifetime)

	// Inducing multiple errors
	builder.WithAttributeDuring("age", lifetime, nil)               // Error 1: nil value
	builder.WithAttributeDuring("name", lifetime, []string{"John"}) // Error 2: non-primitive

	if err := builder.Errors(); err == nil {
		t.Error("expected errors to be accumulated, got nil")
	} else if content, buildErr := builder.Build(); buildErr == nil {
		t.Error("expected build to fail and return accumulated errors")
	} else if content != nil && !content.Activity().IsEmpty() {
		t.Error("expected an empty instance on error build")
	}
}

func TestBuilderWithoutAttributeDuring(t *testing.T) {
	now := time.Now()
	lifetime := periods.NewPeriodSince(now, true)

	builder := objects.NewLocalInstanceBuilder("id").
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
	// Now instance is well-defined in the function scope
	builder = objects.LocalInstanceBuilderLoad(content)
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
	content, err := objects.NewLocalInstanceBuilder("id").
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
// INSTANCE METHODS (MATCHES & SAME)
// =========================================================================

func TestInstanceMatches(t *testing.T) {
	now := time.Now()
	activePeriod := periods.NewPeriodSince(now, true)

	content, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(activePeriod).
		WithAttributeDuring("role", activePeriod, "admin").
		WithAttributeDuring("level", activePeriod, 5).
		Build()
	if err != nil {
		t.Error(err)
	}

	// 1. Perfect match
	validTrait, _ := objects.NewTraitBuilder().WithName("valid").
		WithAttribute("role", "string").
		WithAttribute("level", "int").
		Build()
	if matchPeriod, matches := content.Matches(validTrait); !matches {
		t.Error("expected instance to match valid trait")
	} else if !matchPeriod.Equals(activePeriod) {
		t.Error("expected matching period to equal the instance's active period")
	}

	// 2. Wrong type
	invalidTypeTrait, _ := objects.NewTraitBuilder().WithName("int role").
		WithAttribute("role", "int").
		Build()
	// Role is string in content

	if _, matches := content.Matches(invalidTypeTrait); matches {
		t.Error("expected instance NOT to match trait due to incorrect type")
	}

	// 3. Missing attribute
	missingAttrTrait, _ := objects.NewTraitBuilder().WithName("missing fields").
		WithAttribute("unknown_field", "string").
		Build()
	if _, matches := content.Matches(missingAttrTrait); matches {
		t.Error("expected instance NOT to match trait due to missing attribute")
	}
}

func TestInstanceSame(t *testing.T) {
	p := periods.NewFullPeriod()

	c1, _ := objects.NewLocalInstanceBuilder("1").
		WithActivity(p).
		WithAttributeDuring("key", p, "value").
		Build()

	c1Copy, _ := objects.NewLocalInstanceBuilder("1").
		WithActivity(p).
		WithAttributeDuring("key", p, "value").
		Build()

	c2, _ := objects.NewLocalInstanceBuilder("2").
		WithActivity(p).
		WithAttributeDuring("key", p, "value").
		Build()

	c3, _ := objects.NewLocalInstanceBuilder("3").
		WithActivity(p).
		WithAttributeDuring("key", p, "different_value").
		Build()

	c4, _ := objects.NewLocalInstanceBuilder("4").
		WithActivity(periods.NewEmptyPeriod()). // Different activity
		WithAttributeDuring("key", p, "value").
		Build()

	if !c1.Same(c1Copy) {
		t.Error("expected identical instances to be evaluated as same")
	} else if c1.Same(c2) {
		t.Error("expected different is but identical contents NOT to be evaluated as same")
	} else if c1.Same(c3) {
		t.Error("expected contents with different values NOT to be evaluated as same")
	} else if c1.Same(c4) {
		t.Error("expected contents with different activities NOT to be evaluated as same")
	}
}

// =========================================================================
// DEEP COVERAGE TESTS (HOLES, MERGES, RANGE EXITS, NIL COMPARISONS)
// =========================================================================

func TestWithoutAttributeDuringPartialDeletion(t *testing.T) {
	now := time.Now()
	day1 := now.Add(time.Hour * 24)
	day2 := now.Add(time.Hour * 48)
	day3 := now.Add(time.Hour * 72)
	day4 := now.Add(time.Hour * 96)

	fullPeriod := periods.NewFinitePeriod(now, day4, true, true)
	// The hole is right in the middle (day2 to day3)
	holePeriod := periods.NewFinitePeriod(day2, day3, true, true)

	instance, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(fullPeriod).
		WithAttributeDuring("status", fullPeriod, "active").
		WithoutAttributeDuring("status", holePeriod).
		Build()

	if err != nil {
		t.Fatalf("unexpected error during build: %v", err)
	}

	value, found := instance.Value("status")
	if !found {
		t.Fatalf("expected 'status' attribute to exist")
	}

	// 1. Verify value exists BEFORE the hole
	if _, ok := value.At(day1); !ok {
		t.Error("expected value to exist before the hole (day 1)")
	}

	// 2. Verify value is ABSENT DURING the hole
	midHole := day2.Add(time.Hour * 12)
	if _, ok := value.At(midHole); ok {
		t.Error("expected no value inside the removed hole period")
	}

	// 3. Verify value exists AFTER the hole
	midAfter := day3.Add(time.Hour * 12)
	if _, ok := value.At(midAfter); !ok {
		t.Error("expected value to exist after the hole (day 3+)")
	}
}

func TestAdjacentPeriodsUnion(t *testing.T) {
	now := time.Now()
	mid := now.Add(time.Hour * 24)
	end := now.Add(time.Hour * 48)

	// Two adjacent periods: [now, mid) and [mid, end]
	p1 := periods.NewFinitePeriod(now, mid, true, false)
	p2 := periods.NewFinitePeriod(mid, end, true, true)
	fullActivity := p1.Union(p2)

	// Inserting the SAME value on adjacent periods should merge them
	instance, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(fullActivity).
		WithAttributeDuring("state", p1, "on").
		WithAttributeDuring("state", p2, "on").
		Build()

	if err != nil {
		t.Fatalf("unexpected error during build: %v", err)
	}

	value, found := instance.Value("state")
	if !found {
		t.Fatalf("expected 'state' attribute to exist")
	}

	iterations := 0
	value.Range(func(p periods.Period, v any) bool {
		iterations++
		return true
	})

	// The internal valuesHandler.withValueDuring should merge overlapping/adjacent periods
	// sharing the exact same value.
	if iterations != 1 {
		t.Errorf("expected adjacent periods with same value to be merged into 1 node, got %d", iterations)
	}
}

func TestRangeEarlyExit(t *testing.T) {
	now := time.Now()
	p1 := periods.NewFinitePeriod(now, now.Add(time.Hour), true, false)
	p2 := periods.NewFinitePeriod(now.Add(time.Hour), now.Add(time.Hour*2), true, false)
	p3 := periods.NewFinitePeriod(now.Add(time.Hour*2), now.Add(time.Hour*3), true, true)
	fullActivity := p1.Union(p2).Union(p3)

	// 3 distinct periods with 3 distinct values
	instance, err := objects.NewLocalInstanceBuilder("id").
		WithActivity(fullActivity).
		WithAttributeDuring("phase", p1, "A").
		WithAttributeDuring("phase", p2, "B").
		WithAttributeDuring("phase", p3, "C").
		Build()

	if err != nil {
		t.Fatalf("unexpected error during build: %v", err)
	}

	value, found := instance.Value("phase")
	if !found {
		t.Fatalf("expected 'phase' attribute to exist")
	}

	iterations := 0
	value.Range(func(p periods.Period, v any) bool {
		iterations++
		// Instruct Range to break the loop after 2 iterations
		return iterations < 2
	})

	if iterations != 2 {
		t.Errorf("expected Range to exit after exactly 2 iterations, got %d", iterations)
	}
}

func TestInstanceSameEdgeCases(t *testing.T) {
	emptyPeriod := periods.NewEmptyPeriod()

	// Create two completely empty instances
	c1, err1 := objects.NewLocalInstanceBuilder("id").WithActivity(emptyPeriod).Build()

	if err1 != nil {
		t.Fatalf("failed to build empty instance")
	}

	// 1. Compare with nil (safe interface handling)
	if c1.Same(nil) {
		t.Error("expected instantiated content NOT to be same as nil")
	}
}

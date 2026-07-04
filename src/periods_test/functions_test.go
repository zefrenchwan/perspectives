package periods_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestFunctionEmpty(t *testing.T) {
	function := periods.NewTimeFunction[int]("int", func(a int, b int) bool { return a == b })
	emptyFunction := periods.NewTimeFunction[int]("int", func(a int, b int) bool { return a == b })
	if !function.IsEmpty() {
		t.Errorf("Expected function to be empty")
	} else if function.DataType() != "int" {
		t.Errorf("Expected function to have data type int")
	} else if !function.Domain().Equals(periods.NewEmptyPeriod()) {
		t.Errorf("function domain should be empty")
	} else if !function.Equals(emptyFunction) {
		t.Errorf("function should not be equal to empty set")
	} else if !function.Equals(function) {
		t.Errorf("function should be equal to itself")
	} else if _, has := function.At(time.Now()); has {
		t.Errorf("Expected function.At(time.Now()) to be nil")
	}

	function.Add(1, periods.NewEmptyPeriod())
	if !function.IsEmpty() {
		t.Errorf("Expected function to be empty after adding empty period")
	}
}

func TestFunctionAt(t *testing.T) {
	function := periods.NewTimeFunction[int]("int", func(a int, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	if _, has := function.At(now); has {
		t.Errorf("Expected function.At(time.Now()) to be nil")
	}

	function.Add(1, periods.NewPeriodSince(now, true))
	function.Add(2, periods.NewPeriodSince(after, true))

	if matching, has := function.At(now); !has {
		t.Errorf("Expected function.At(time.Now()) to be non-nil")
	} else if matching != 1 {
		t.Errorf("Expected function.At(time.Now()) to be 1")
	}

	if _, has := function.At(before); has {
		t.Errorf("Expected no element at before")
	}

	if matching, has := function.At(after); !has {
		t.Errorf("Expected function.At(after) to be non-nil")
	} else if matching != 2 {
		t.Errorf("Expected function.At(after)) to be 2")
	}
}

func TestFunctionRemove(t *testing.T) {
	function := periods.NewTimeFunction[int]("int", func(a int, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)

	function.Add(1, periods.NewPeriodSince(now, true))
	function.Add(2, periods.NewPeriodSince(after, true))
	function.Remove(periods.NewPeriodSince(after, true))

	if matching, has := function.At(now); !has {
		t.Errorf("Expected function.At(time.Now()) to be non-nil")
	} else if matching != 1 {
		t.Errorf("Expected function.At(time.Now()) to be 1")
	}

	if _, has := function.At(before); has {
		t.Errorf("Expected no element at before")
	}

	if _, has := function.At(after); has {
		t.Errorf("Expected function.At(after) to be empty due to removal")
	}
}

func TestFunctionAdd(t *testing.T) {
	function := periods.NewTimeFunction[string]("string", func(a, b string) bool { return a == b })
	now := time.Now().Truncate(time.Hour)
	before := now.Add(-10 * time.Hour)
	after := now.Add(10 * time.Hour)

	// Add Alice for [before, after]
	function.Add("Alice", periods.NewFinitePeriod(before, after, true, true))

	// Add Bob for [now, after]
	// Because it's a function, this should truncate Alice's period to [before, now[
	function.Add("Bob", periods.NewFinitePeriod(now, after, true, true))

	var alicePeriod periods.Period
	var bobPeriod periods.Period
	var aliceCount, bobCount int

	for p, v := range function.Range() {
		if v == "Alice" {
			alicePeriod = p
			aliceCount++
		} else if v == "Bob" {
			bobPeriod = p
			bobCount++
		}
	}

	if aliceCount != 1 {
		t.Errorf("Expected Alice to have exactly 1 matching period, got %d", aliceCount)
	}
	if bobCount != 1 {
		t.Errorf("Expected Bob to have exactly 1 matching period, got %d", bobCount)
	}

	expectedAlicePeriod := periods.NewFinitePeriod(before, now, true, false)
	if !alicePeriod.Equals(expectedAlicePeriod) {
		t.Errorf("Expected Alice's period to be truncated to %s, got %s", expectedAlicePeriod.AsRawString(), alicePeriod.AsRawString())
	}

	expectedBobPeriod := periods.NewFinitePeriod(now, after, true, true)
	if !bobPeriod.Equals(expectedBobPeriod) {
		t.Errorf("Expected Bob's period to be %s, got %s", expectedBobPeriod.AsRawString(), bobPeriod.AsRawString())
	}
}

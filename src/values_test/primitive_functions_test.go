package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestPrimitiveFunctionsTypes(t *testing.T) {
	acceptedTypes := []string{
		values.PRIMITIVE_TYPE_STRING,
		values.PRIMITIVE_TYPE_INT,
		values.PRIMITIVE_TYPE_BOOL,
		values.PRIMITIVE_TYPE_FLOAT,
		values.PRIMITIVE_TYPE_TIME,
	}

	refusedTypes := []string{
		"",
		"not a real type",
	}

	for _, a := range acceptedTypes {
		if function, err := values.NewPrimitiveTimeFunction(a); err != nil {
			t.Errorf("Expected function to be created")
		} else if !function.IsEmpty() {
			t.Errorf("Expected function to be empty")
		} else if function.Datatype() != a {
			t.Errorf("Expected function datatype to be %s", a)
		}
	}

	for _, r := range refusedTypes {
		if _, err := values.NewPrimitiveTimeFunction(r); err == nil {
			t.Errorf("Expected function to be refused because type is not valid")
		}
	}
}

func TestIntPrimitiveFunction(t *testing.T) {
	function := values.NewIntTimeFunction()
	function.Add(10, periods.NewFullPeriod())

	if v, has := function.Value(time.Now()); !has {
		t.Errorf("Expected value to be present")
	} else if v != 10 {
		t.Errorf("Expected value to be 10")
	} else if p, has := function.At(time.Now()); !has {
		t.Errorf("Expected period to be present")
	} else if p.Content() != 10 {
		t.Errorf("Expected period content to be 10")
	}

	// test value override
	function.Add(20, periods.NewFullPeriod())
	if v, has := function.Value(time.Now()); !has {
		t.Errorf("Expected value to be present")
	} else if v != 20 {
		t.Errorf("Expected value to be 20")
	} else if p, has := function.At(time.Now()); !has {
		t.Errorf("Expected period to be present")
	} else if p.Content() != 20 {
		t.Errorf("Expected period content to be 10")
	}

	function.Remove(periods.NewFullPeriod())
	if !function.IsEmpty() {
		t.Errorf("Expected function to be empty")
	}
}

func TestPrimitiveFunctionRange(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.Add(-time.Hour)
	after := now.Add(time.Hour)
	beforePeriod := periods.NewPeriodUntil(before, true)
	currentPeriod := periods.NewFinitePeriod(before, after, false, false)
	afterPeriod := periods.NewPeriodSince(after, true)

	function := values.NewStringTimeFunction()
	function.Add("before", beforePeriod)
	function.Add("current", currentPeriod)
	function.Add("after", afterPeriod)

	if v, has := function.Value(before); !has {
		t.Errorf("Expected value to be present")
	} else if v != "before" {
		t.Errorf("Expected value to be 'before'")
	}

	if v, has := function.Value(now); !has {
		t.Errorf("Expected current period to be present")
	} else if v != "current" {
		t.Errorf("Expected period content to be 'current'")
	}

	if v, has := function.Value(after); !has {
		t.Errorf("Expected after period to be present")
	} else if v != "after" {
		t.Errorf("Expected period content to be 'after'")
	}

	for period, pvalue := range function.Range() {
		if period.Equals(beforePeriod) {
			if value := pvalue.Content(); value != "before" {
				t.Errorf("Expected value to be 'before'")
			}
		} else if period.Equals(currentPeriod) {
			if value := pvalue.Content(); value != "current" {
				t.Errorf("Expected value to be 'current'")
			}
		} else if period.Equals(afterPeriod) {
			if value := pvalue.Content(); value != "after" {
				t.Errorf("Expected value to be 'after'")
			}
		} else {
			t.Errorf("Unexpected period")
		}
	}

	for period, value := range function.Values() {
		if period.Equals(beforePeriod) {
			if value != "before" {
				t.Errorf("Expected value to be 'before'")
			}
		} else if period.Equals(currentPeriod) {
			if value != "current" {
				t.Errorf("Expected value to be 'current'")
			}
		} else if period.Equals(afterPeriod) {
			if value != "after" {
				t.Errorf("Expected value to be 'after'")
			}
		} else {
			t.Errorf("Unexpected period")
		}
	}
}

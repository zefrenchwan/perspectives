package objects_test_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestAdd(t *testing.T) {
	values := objects.NewTemporalIntValues()
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}

	values.Add(periods.NewFullPeriod(), 10)
	if res, found := values.At(time.Now()); !found || res != 10 {
		t.Errorf("Expected value at current time to be 10, got %v", res)
	}

	values.Add(periods.NewFullPeriod(), 20)
	if res, found := values.At(time.Now()); !found || res != 20 {
		t.Errorf("Expected value at current time to be 20, got %v", res)
	}

	values.Clear()
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}
}

func TestRemove(t *testing.T) {
	values := objects.NewTemporalIntValues()
	values.Add(periods.NewFullPeriod(), 10)
	values.Remove(periods.NewFullPeriod())
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}

	values.Add(periods.NewFullPeriod(), 50)
	values.Remove(periods.NewPeriodUntil(time.Now().Add(24*time.Hour), false))
	if _, found := values.At(time.Now()); found {
		t.Errorf("values without period should start in 24 hours, cannot have value now")
	}

}

func TestCut(t *testing.T) {
	values := objects.NewTemporalIntValues()
	values.Add(periods.NewFullPeriod(), 10)
	if res, found := values.At(time.Now()); !found || res != 10 {
		t.Errorf("Expected value at current time to be 10, got %v", res)
	}

	nextPeriod := periods.NewPeriodSince(time.Now().Add(24*time.Hour), true)
	cutValues := values.Cut(nextPeriod)
	if cutValues.IsEmpty() {
		t.Errorf("Expected cut values to not be empty, got %v", cutValues)
	} else if _, found := cutValues.At(time.Now()); found {
		t.Errorf("cutValues should start in 24 hours, cannot have value now")
	}
}

func TestRange(t *testing.T) {
	values := objects.NewTemporalIntValues()
	values.Add(periods.NewFullPeriod(), 10)
	if res, found := values.At(time.Now()); !found || res != 10 {
		t.Errorf("Expected value at current time to be 10, got %v", res)
	}

	for period, value := range values.Range {
		if value != 10 {
			t.Errorf("Expected value for period %v to be 10, got %v", period, value)
		} else if !period.Equals(periods.NewFullPeriod()) {
			t.Errorf("Expected period to be full period, got %v", period)
		}
	}
}

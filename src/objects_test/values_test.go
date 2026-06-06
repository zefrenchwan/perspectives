package objects_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestAdd(t *testing.T) {
	values := objects.NewTemporalValues()
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}

	values = values.Add(periods.NewFullPeriod(), 10)
	if res, found := values.At(time.Now()); !found || res != 10 {
		t.Errorf("Expected value at current time to be 10, got %v", res)
	}

	values = values.Add(periods.NewFullPeriod(), 20)
	if res, found := values.At(time.Now()); !found || res != 20 {
		t.Errorf("Expected value at current time to be 20, got %v", res)
	}

	values = values.Remove(periods.NewFullPeriod())
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}
}

func TestRemove(t *testing.T) {
	values := objects.NewTemporalValues()
	values = values.Add(periods.NewFullPeriod(), 10)
	values = values.Remove(periods.NewFullPeriod())
	if !values.IsEmpty() {
		t.Errorf("Expected values to be empty, got %v", values)
	}

	values = values.Add(periods.NewFullPeriod(), 50)
	values = values.Remove(periods.NewPeriodUntil(time.Now().Add(24*time.Hour), false))
	if _, found := values.At(time.Now()); found {
		t.Errorf("values without period should start in 24 hours, cannot have value now")
	} else if value, found := values.At(time.Now().Add(48 * time.Hour)); !found || value != 50 {
		t.Errorf("Expected value at 48 hours to be 50, got %v", value)
	}

}

func TestCut(t *testing.T) {
	values := objects.NewTemporalValues()
	values = values.Add(periods.NewFullPeriod(), 10)
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
	values := objects.NewTemporalValues()
	values = values.Add(periods.NewFullPeriod(), 10)
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

func TestRangeEmpty(t *testing.T) {
	values := objects.NewTemporalValues()
	for period, value := range values.Range {
		t.Errorf("Expected no values in empty range, got period %v with value %v", period, value)
	}
}

func TestDataTypes(t *testing.T) {
	values := objects.NewTemporalValues()
	if values.DataType() != "" {
		t.Errorf("Expected new values to have no data type, got %v", values.DataType())
	}

	values = values.Add(periods.NewFullPeriod(), 10)
	if values.DataType() != "int" {
		t.Errorf("Expected data type to be int after adding int value, got %v", values.DataType())
	}

	values = values.Add(periods.NewFullPeriod(), "twenty ! Happy birthday my friend !!!")
	if values.DataType() != "string" {
		t.Errorf("Expected data type to be string (not any) because full period changed it all, got %v", values.DataType())
	}

	values = values.Add(periods.NewPeriodSince(time.Now(), false), 50)
	if values.DataType() != "any" {
		t.Errorf("Expected data type to be any because int and string coexist, got %v", values.DataType())
	}
}

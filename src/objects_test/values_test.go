package objects_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestTemporalValuesAdd(t *testing.T) {
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

func TestTemporalValuesRemove(t *testing.T) {
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

func TestTemporalValuesCut(t *testing.T) {
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

func TestTemporalValuesRange(t *testing.T) {
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

func TestTemporalValuesRangeEmpty(t *testing.T) {
	values := objects.NewTemporalValues()
	for period, value := range values.Range {
		t.Errorf("Expected no values in empty range, got period %v with value %v", period, value)
	}
}

func TestTemporalValuesDataTypes(t *testing.T) {
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

func TestContentAdd(t *testing.T) {
	content := objects.NewContent()
	content = content.Add("age", periods.NewFullPeriod(), 20)
	if _, found := content.Value("not existing"); found {
		t.Errorf("Expected 'not existing' value to not be found")
	} else if temporalValue, found := content.Value("age"); !found {
		t.Errorf("age should be present in values")
	} else if v, has := temporalValue.At(time.Now()); !has || v != 20 {
		t.Errorf("expecting age to be filled, got %v for value and %v for found", v, has)
	}
}

func TestContentAt(t *testing.T) {
	content := objects.NewContent()
	if !content.Activity().Equals(periods.NewFullPeriod()) {
		t.Log("Expected activity to be full period, got different")
		t.Fail()
	}

	now := time.Now()
	before := now.Add(-time.Hour)
	after := now.Add(time.Hour)
	content = content.Add("age", periods.NewPeriodSince(now, false), 20)
	content = content.Add("name", periods.NewFullPeriod(), "John Doe")
	content = content.Add("address", periods.NewPeriodUntil(now, true), "123 Main St")

	if valuesBefore, hasBefore := content.At(before); !hasBefore {
		t.Errorf("Expected 'before' period to have values, got none")
	} else if len(valuesBefore) != 2 {
		t.Errorf("Expected 2 values before, got %d", len(valuesBefore))
	} else if valuesBefore["name"] != "John Doe" {
		t.Errorf("Expected 'name' value to be 'John Doe', got %v", valuesBefore["name"])
	} else if valuesBefore["address"] != "123 Main St" {
		t.Errorf("Expected 'address' value to be '123 Main St', got %v", valuesBefore["address"])
	}

	if valuesAfter, hasAfter := content.At(after); !hasAfter {
		t.Errorf("Expected 'after' period to have values, got none")
	} else if len(valuesAfter) != 2 {
		t.Errorf("Expected 2 values after, got %d", len(valuesAfter))
	} else if valuesAfter["age"] != 20 {
		t.Errorf("Expected 'age' value to be 20, got %v", valuesAfter["age"])
	} else if valuesAfter["name"] != "John Doe" {
		t.Errorf("Expected 'name' value to be 'John Doe', got %v", valuesAfter["name"])
	}
}

func TestContentRemove(t *testing.T) {
	content := objects.NewContent()

	now := time.Now()
	before := now.Add(-time.Hour)
	after := now.Add(time.Hour)

	content = content.Add("age", periods.NewPeriodSince(before, false), 20)
	content = content.Add("name", periods.NewFullPeriod(), "John Doe")
	content = content.Add("address", periods.NewPeriodSince(after, true), "123 Main St")

	content = content.Remove("age", periods.NewPeriodUntil(now, false))
	if values, has := content.Value("age"); !has {
		t.Errorf("Expected 'age' value to be valid between before and now, but was removed entirely")
	} else if values.Validity().Equals(periods.NewFinitePeriod(before, after, false, true)) {
		t.Errorf("Expected 'age' value to be valid between before and now, got %v", values.Validity())
	}
}

func TestContentCut(t *testing.T) {
	content := objects.NewContent()
	now := time.Now()
	before := now.Add(-time.Hour)
	after := now.Add(time.Hour)

	content = content.Add("age", periods.NewFullPeriod(), 20)
	content = content.Add("name", periods.NewFullPeriod(), "John Doe")

	cutPeriod := periods.NewFinitePeriod(before, after, true, true)
	cutContent := content.Cut(cutPeriod)

	if cutContent == nil {
		t.Fatal("Expected cut content to not be nil")
	}

	if !cutContent.Activity().Equals(cutPeriod) {
		t.Errorf("Expected activity to be %v, got %v", cutPeriod, cutContent.Activity())
	}

	for _, attr := range []string{"age", "name"} {
		if tv, found := cutContent.Value(attr); !found {
			t.Errorf("Expected attribute %s to be found", attr)
		} else if !tv.Validity().Equals(cutPeriod) {
			t.Errorf("Expected attribute %s validity to be %v, got %v", attr, cutPeriod, tv.Validity())
		}
	}
}

func TestContentDescription(t *testing.T) {
	content := objects.NewContent()
	content = content.Add("age", periods.NewFullPeriod(), 20)
	content = content.Add("name", periods.NewFullPeriod(), "John Doe")
	content = content.Add("weight", periods.NewFullPeriod(), 75.5)

	desc := content.Description()
	expected := map[string]string{
		"age":    "int",
		"name":   "string",
		"weight": "float64",
	}

	if len(desc) != len(expected) {
		t.Errorf("Expected description length %d, got %d", len(expected), len(desc))
	}

	for attr, dataType := range expected {
		if desc[attr] != dataType {
			t.Errorf("Expected attribute %s to have type %s, got %s", attr, dataType, desc[attr])
		}
	}
}

func TestContentSame(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour)
	after := now.Add(time.Hour)

	c1 := objects.NewContent().Add("age", periods.NewFullPeriod(), 25)
	c2 := objects.NewContent().Add("age", periods.NewFullPeriod(), 25)

	if !c1.Same(c2) {
		t.Errorf("Expected identical content to be the same")
	}

	c3 := c1.Cut(periods.NewFinitePeriod(before, after, true, true))
	if c1.Same(c3) {
		t.Errorf("Expected content with different activity periods to not be the same")
	}

	c4 := c1.Add("name", periods.NewFullPeriod(), "Alice")
	if c1.Same(c4) {
		t.Errorf("Expected content with different attributes to not be the same")
	}
}

package models

import (
	"maps"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestValueSet(t *testing.T) {
	value := models.NewValue(50)
	values := value.GetValues()
	if len(values) != 1 || values[0] != 50 {
		t.Log("Fail to read values")
		t.Fail()
	}

	value.Set(40)
	if len(values) != 1 || values[0] != 50 {
		t.Log("Fail to update values")
		t.Fail()
	}

	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(10, 0, 0)
	value.SetUntil(30, now, true)

	if v, found := value.GetValue(before); !found {
		t.Log("Value should be set because interval is ]-oo, now]")
		t.Log(value)
		t.Fail()
	} else if v != 30 {
		t.Log("Expecting value set before now, got ", v)
		t.Log(value)
		t.Fail()
	}

	if v, found := value.GetValue(after); !found {
		t.Log("Value should be unchanged outside ]-oo, now]")
		t.Fail()
	} else if v != 40 {
		t.Log("Expecting value set after now")
		t.Fail()
	}
}

func TestValueGet(t *testing.T) {
	value := models.NewValue(50)
	expected := map[int]models.Period{
		50: models.NewFullPeriod(),
	}

	if !maps.EqualFunc(expected, value.Get(), func(a, b models.Period) bool { return a.Equals(b) }) {
		t.Log("Failed to flatten content")
		t.Fail()
	}

	now := time.Now().Truncate(time.Hour)
	value.SetUntil(30, now, true)
	// value is now
	// ]-oo, now] => 30
	// ]now, +oo[ => 50
	expected = map[int]models.Period{
		30: models.NewPeriodUntil(now, true),
		50: models.NewPeriodSince(now, false),
	}

	values := value.Get()

	if !maps.EqualFunc(expected, values, func(a, b models.Period) bool { return a.Equals(b) }) {
		t.Log("Failed to flatten content with multiple values")
		t.Log(values)
		t.Log(expected)
		t.Fail()
	}
}

package structures_test

import (
	"maps"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestValueSet(t *testing.T) {
	value := structures.NewValue(50)
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
	value := structures.NewValue(50)
	expected := map[int]structures.Period{
		50: structures.NewFullPeriod(),
	}

	if !maps.EqualFunc(expected, value.Get(), func(a, b structures.Period) bool { return a.Equals(b) }) {
		t.Log("Failed to flatten content")
		t.Fail()
	}

	now := time.Now().Truncate(time.Hour)
	value.SetUntil(30, now, true)
	// value is now
	// ]-oo, now] => 30
	// ]now, +oo[ => 50
	expected = map[int]structures.Period{
		30: structures.NewPeriodUntil(now, true),
		50: structures.NewPeriodSince(now, false),
	}

	values := value.Get()

	if !maps.EqualFunc(expected, values, func(a, b structures.Period) bool { return a.Equals(b) }) {
		t.Log("Failed to flatten content with multiple values")
		t.Log(values)
		t.Log(expected)
		t.Fail()
	}
}

func TestSerde(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	expected := structures.NewValueSince(44, now, true)
	result := expected.BuildCompactMapOfValues()
	if value, err := structures.LoadValuesFromCompactMap(result); err != nil {
		t.Log(err)
		t.Fail()
	} else if len(value) != 1 {
		t.Log("failed to load content: sizes differ")
		t.Fail()
	} else {
		for k, v := range expected {
			if other, found := value[k]; !found {
				t.Log("values not found")
				t.Fail()
			} else if other != v {
				t.Log("values mismatch")
				t.Fail()
			}
		}
	}

}

func TestValuesCut(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	after := now.AddDate(10, 0, 0)
	reference := structures.NewPeriodUntil(now, true)
	value := structures.NewValue("a")
	result := value.Cut(reference).Get()

	if len(result) != 1 {
		t.Fail()
	} else if period := result["a"]; !period.Equals(reference) {
		t.Fail()
	}

	// test empty result
	reference = structures.NewPeriodSince(after, true)
	value = structures.NewValueUntil("a", now, true)
	result = value.Cut(reference).Get()

	if len(result) != 0 {
		t.Log("empty expected")
		t.Log(result)
		t.Fail()
	}
}

func TestValuesSetDuringPeriod(t *testing.T) {
	value := structures.NewValue("test")
	now := time.Now().Truncate(time.Hour)
	after := now.AddDate(1, 0, 0)
	before := now.AddDate(-1, 0, 0)
	value.SetUntil("other", now, true)
	// value is
	// ]-oo,now[ => "other"
	// [now, +oo[ => test

	if result, found := value.GetValue(after); !found || result != "test" {
		t.Fail()
	} else if result, found := value.GetValue(before); !found || result != "other" {
		t.Fail()
	}

	// value changes back forever.
	// ]-oo,+oo[ => "final"
	value.Set("final")
	if result, found := value.GetValue(after); !found || result != "final" {
		t.Fail()
	} else if result, found := value.GetValue(before); !found || result != "final" {
		t.Fail()
	}
}

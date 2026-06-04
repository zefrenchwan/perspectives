package periods_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestPeriodComplements(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	complement := value.Complement()
	expected := periods.NewPeriodUntil(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewPeriodUntil(now, true)
	complement = value.Complement()
	expected = periods.NewPeriodSince(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewEmptyPeriod()
	complement = value.Complement()
	expected = periods.NewFullPeriod()
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewFullPeriod()
	complement = value.Complement()
	expected = periods.NewEmptyPeriod()
	if !expected.Equals(complement) {
		t.Log("Complement failed to reverse full to empty")
		t.Fail()
	}
}

func TestIntersectionWithFull(t *testing.T) {
	value := periods.NewFullPeriod()
	now := time.Now()
	otherValue := periods.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(otherValue) {
		t.Log("intersection with full failed")
		t.Fail()
	}
}

func TestIntersectionWithEmpty(t *testing.T) {
	value := periods.NewEmptyPeriod()
	now := time.Now()
	otherValue := periods.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(value) {
		t.Log("intersection with empty failed")
		t.Fail()
	}
}

func TestIntersectionWithOther(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour)
	value := periods.NewPeriodSince(before, true)
	otherValue := periods.NewPeriodUntil(now, true)
	expected := periods.NewFinitePeriod(before, now, true, true)
	result := otherValue.Intersection(value)
	if !result.Equals(expected) {
		t.Logf("intersection with other failed: got %s BUT EXPECTED %s", result.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestUnionWithEmpty(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	expected := value
	result := periods.NewEmptyPeriod().Union(value)
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}

	result = value.Union(periods.NewEmptyPeriod())
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestUnionWithSame(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	expected := value
	result := value.Union(value)
	if !result.Equals(expected) {
		t.Logf("union with same failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestInfiniteUnion(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	first := periods.NewPeriodSince(before, true)
	second := periods.NewPeriodUntil(now, true)
	expected := periods.NewFullPeriod()
	result := first.Union(second)
	if !result.Equals(expected) {
		t.Logf("union to make full failed: got %s", result.AsRawString())
		t.Fail()
	}

	result = second.Union(first)
	if !result.Equals(expected) {
		t.Logf("union to make full failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestPeriodContains(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	// period is [before, +oo[
	period := periods.NewPeriodSince(before, true)
	// before is in period
	if !period.Contains(before) {
		t.Log("Failed to test when value is boundary included")
		t.Fail()
	}

	// interval is [before, +oo[ and before < now, it should belong
	if !period.Contains(now) {
		t.Log("Failed to test outside value")
		t.Fail()
	}

	// period is ]before, +oo[
	period = periods.NewPeriodSince(before, false)
	// before is not in interval before it is excluded
	if period.Contains(before) {
		t.Log("Failed to test when value is boundary excluded")
		t.Fail()
	}
	// interval is ]before,+oo[ and before < now
	if !period.Contains(now) {
		t.Log("Failed to test outside value")
		t.Fail()
	}

	// period is ]-oo, now]
	period = periods.NewPeriodUntil(now, true)
	// before < now, so expecting period to contain in
	if !period.Contains(before) {
		t.Log("Failed to test when value is strictly included")
		t.Fail()
	}

	// now is included, so expecting to belong
	if !period.Contains(now) {
		t.Log("Failed to test bound value included")
		t.Fail()
	}

	// period is ]-oo, now[
	period = periods.NewPeriodUntil(now, false)
	// period should not contain now, because now is excluded
	if period.Contains(now) {
		t.Log("Failed to test bound value excluded")
		t.Fail()
	}
}

func TestPeriodRemove(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-10, 0, 0)
	after := now.AddDate(10, 0, 0)
	period := periods.NewPeriodUntil(after, true)
	toRemove := periods.NewFinitePeriod(before, now, true, false)
	remaining := period.Remove(toRemove)
	expected := periods.NewPeriodUntil(before, false).Union(periods.NewFinitePeriod(now, after, true, true))
	if !remaining.Equals(expected) {
		t.Logf("Failed to remove period, got %s but expected %s", remaining.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestPeriodRemoveLargerPeriod(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-10, 0, 0)
	after := now.AddDate(10, 0, 0)
	period := periods.NewFinitePeriod(before, now, true, true)
	toRemove := periods.NewPeriodUntil(after, true)
	remaining := period.Remove(toRemove)
	// remaining should be empty because [before, now] is in ]-oo, after]
	if !remaining.IsEmpty() {
		t.Log("remaining should be empty because toRemove contains each point of period")
		t.Fail()
	}

	toRemove = periods.NewFullPeriod()
	remaining = period.Remove(toRemove)
	// remaining should be empty because [before, now] is in the full period
	if !remaining.IsEmpty() {
		t.Log("remaining should be empty because toRemove contains period")
		t.Fail()
	}

}

func TestPeriodSerde(t *testing.T) {
	tested := periods.NewEmptyPeriod()
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser empty")
		t.Fail()
	}

	tested = periods.NewFullPeriod()
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser full")
		t.Fail()
	}

	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(10, 0, 0)
	tested = periods.NewFinitePeriod(before, now, true, false)
	tested = tested.Union(periods.NewPeriodSince(after, true))
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser union of intervals")
		t.Fail()
	}
}

func TestPeriodInfiniteBoundaries(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	a := periods.NewPeriodSince(after, true)
	b := periods.NewPeriodUntil(before, true)
	res := a.Union(b).BoundingPeriod()
	expected := periods.NewFullPeriod()
	if !expected.Equals(res) {
		t.Logf("failed to find full as boundaries, got %s", res.AsRawString())
		t.Fail()
	}
}

func TestPeriodFiniteBoundaries(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	evenAfter := after.AddDate(10, 0, 0)
	a := periods.NewFinitePeriod(before, now, false, true)
	b := periods.NewFinitePeriod(after, evenAfter, true, true)
	res := a.Union(b).BoundingPeriod()
	expected := periods.NewFinitePeriod(before, evenAfter, false, true)
	if !expected.Equals(res) {
		t.Logf("failed to find finite intervals as boundaries, got %s", res.AsRawString())
		t.Fail()
	}
}

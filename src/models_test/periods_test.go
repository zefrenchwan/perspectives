package models_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestPeriodComplements(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := models.NewPeriodSince(now, true)
	complement := value.Complement()
	expected := models.NewPeriodUntil(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = models.NewPeriodUntil(now, true)
	complement = value.Complement()
	expected = models.NewPeriodSince(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = models.NewEmptyPeriod()
	complement = value.Complement()
	expected = models.NewFullPeriod()
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = models.NewFullPeriod()
	complement = value.Complement()
	expected = models.NewEmptyPeriod()
	if !expected.Equals(complement) {
		t.Log("Complement failed to reverse full to empty")
		t.Fail()
	}
}

func TestIntersectionWithFull(t *testing.T) {
	value := models.NewFullPeriod()
	now := time.Now()
	otherValue := models.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(otherValue) {
		t.Log("intersection with full failed")
		t.Fail()
	}
}

func TestIntersectionWithEmpty(t *testing.T) {
	value := models.NewEmptyPeriod()
	now := time.Now()
	otherValue := models.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(value) {
		t.Log("intersection with empty failed")
		t.Fail()
	}
}

func TestIntersectionWithOther(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour)
	value := models.NewPeriodSince(before, true)
	otherValue := models.NewPeriodUntil(now, true)
	expected := models.NewFinitePeriod(before, now, true, true)
	result := otherValue.Intersection(value)
	if !result.Equals(expected) {
		t.Logf("intersection with other failed: got %s BUT EXPECTED %s", result.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestUnionWithEmpty(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := models.NewPeriodSince(now, true)
	expected := value
	result := models.NewEmptyPeriod().Union(value)
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}

	result = value.Union(models.NewEmptyPeriod())
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestUnionWithSame(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := models.NewPeriodSince(now, true)
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
	first := models.NewPeriodSince(before, true)
	second := models.NewPeriodUntil(now, true)
	expected := models.NewFullPeriod()
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
	period := models.NewPeriodSince(before, true)
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
	period = models.NewPeriodSince(before, false)
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
	period = models.NewPeriodUntil(now, true)
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
	period = models.NewPeriodUntil(now, false)
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
	period := models.NewPeriodUntil(after, true)
	toRemove := models.NewFinitePeriod(before, now, true, false)
	remaining := period.Remove(toRemove)
	expected := models.NewPeriodUntil(before, false).Union(models.NewFinitePeriod(now, after, true, true))
	if !remaining.Equals(expected) {
		t.Logf("Failed to remove period, got %s but expected %s", remaining.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestPeriodSerde(t *testing.T) {
	tested := models.NewEmptyPeriod()
	if res, err := models.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser empty")
		t.Fail()
	}

	tested = models.NewFullPeriod()
	if res, err := models.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser full")
		t.Fail()
	}

	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(10, 0, 0)
	tested = models.NewFinitePeriod(before, now, true, false)
	tested = tested.Union(models.NewPeriodSince(after, true))
	if res, err := models.PeriodLoad(tested.AsStrings()); err != nil {
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
	a := models.NewPeriodSince(after, true)
	b := models.NewPeriodUntil(before, true)
	res := a.Union(b).BoundingPeriod()
	expected := models.NewFullPeriod()
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
	a := models.NewFinitePeriod(before, now, false, true)
	b := models.NewFinitePeriod(after, evenAfter, true, true)
	res := a.Union(b).BoundingPeriod()
	expected := models.NewFinitePeriod(before, evenAfter, false, true)
	if !expected.Equals(res) {
		t.Logf("failed to find finite intervals as boundaries, got %s", res.AsRawString())
		t.Fail()
	}
}
